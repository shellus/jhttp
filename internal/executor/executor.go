package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shellus/jhttp/internal/models"
	"github.com/shellus/jhttp/internal/parser"
)

// Executor HTTP请求执行器
type Executor struct {
	client  *http.Client
	verbose bool
}

// NewExecutor 创建一个新的执行器
func NewExecutor(verbose bool) *Executor {
	return &Executor{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 允许最多5次重定向
				if len(via) >= 5 {
					return fmt.Errorf("达到最大重定向次数")
				}
				return nil
			},
		},
		verbose: verbose,
	}
}

// SetTimeout 设置HTTP请求超时时间
func (e *Executor) SetTimeout(timeout time.Duration) {
	e.client.Timeout = timeout
}

// Execute 执行单个HTTP请求
func (e *Executor) Execute(httpFile *models.HTTPFile, request *models.HTTPRequest, env string) (*models.HTTPResponse, error) {
	// 解析请求中的变量
	resolvedReq, err := parser.ResolveVariables(httpFile, request, env)
	if err != nil {
		return nil, fmt.Errorf("解析变量失败: %w", err)
	}

	// 创建HTTP请求
	req, err := e.createHTTPRequest(resolvedReq)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 如果启用详细模式，打印请求信息
	if e.verbose {
		fmt.Printf("> %s %s\n", req.Method, req.URL.String())
		for key, values := range req.Header {
			for _, value := range values {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		if req.Body != nil && resolvedReq.Body != "" {
			fmt.Println(">")
			fmt.Println(resolvedReq.Body)
		}
		fmt.Println()
	}

	// 添加一个有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), e.client.Timeout)
	defer cancel()
	req = req.WithContext(ctx)

	// 执行请求
	startTime := time.Now()
	resp, err := e.client.Do(req)
	duration := time.Since(startTime)

	// 处理请求错误
	if err != nil {
		errorMessage := err.Error()
		errorType := "网络错误"

		// 检查不同类型的错误
		switch {
		case strings.Contains(errorMessage, "context deadline exceeded") ||
			strings.Contains(errorMessage, "timeout") ||
			strings.Contains(errorMessage, "timed out"):
			errorType = "请求超时"
			errorMessage = "请求超时 - 服务器在规定时间内没有响应"

		case strings.Contains(errorMessage, "no such host"):
			errorType = "域名解析错误"
			errorMessage = fmt.Sprintf("无法解析主机名 '%s'", req.URL.Host)

		case strings.Contains(errorMessage, "connection refused"):
			errorType = "连接被拒绝"
			errorMessage = fmt.Sprintf("连接被拒绝 - 服务器 '%s' 拒绝了连接请求", req.URL.Host)

		case strings.Contains(errorMessage, "certificate"):
			errorType = "SSL/TLS 错误"
			errorMessage = "SSL/TLS 证书验证失败"

		case strings.Contains(errorMessage, "no route to host"):
			errorType = "路由错误"
			errorMessage = fmt.Sprintf("无法连接到主机 '%s' - 网络不可达", req.URL.Host)

		case strings.Contains(errorMessage, "i/o timeout"):
			errorType = "I/O 超时"
			errorMessage = "读取/写入操作超时 - 可能是网络问题或服务器响应缓慢"
		}

		// 打印详细的错误信息
		if e.verbose {
			fmt.Printf("\n请求失败: [%s] %s\n", errorType, errorMessage)
			fmt.Printf("原始错误: %v\n", err)
			fmt.Printf("请求耗时: %d ms\n", duration.Milliseconds())
		}

		return &models.HTTPResponse{
			Request: resolvedReq,
			Error:   fmt.Errorf("%s: %s", errorType, errorMessage),
			Time:    duration.Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 创建响应对象
	response := &models.HTTPResponse{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header.Clone(),
		Body:       body,
		BodyString: string(body),
		Time:       duration.Milliseconds(),
		Request:    resolvedReq,
	}

	// 如果启用详细模式，打印响应信息
	if e.verbose {
		fmt.Printf("< HTTP/1.1 %s\n", resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}

		if len(response.Body) > 0 {
			fmt.Println("<")

			// 检查Content-Type以决定如何格式化输出
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// 尝试美化JSON输出（这是简化版，实际中可以使用json包格式化）
				fmt.Println(response.BodyString)
			} else {
				fmt.Println(response.BodyString)
			}
		}

		fmt.Printf("\n请求耗时: %d ms\n", response.Time)
	}

	return response, nil
}

// ExecuteFile 执行HTTP文件中的所有请求
func (e *Executor) ExecuteFile(httpFile *models.HTTPFile, requestName string, env string) ([]*models.HTTPResponse, error) {
	responses := make([]*models.HTTPResponse, 0)

	// 如果指定了请求名称，只执行该请求
	if requestName != "" {
		req := httpFile.FindRequestByName(requestName)
		if req == nil {
			return nil, fmt.Errorf("未找到名为 '%s' 的请求", requestName)
		}

		resp, err := e.Execute(httpFile, req, env)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
		return responses, nil
	}

	// 否则，执行所有请求
	for _, req := range httpFile.Requests {
		if e.verbose {
			fmt.Printf("\n===== 执行请求: %s =====\n", req.Name)
		}

		resp, err := e.Execute(httpFile, req, env)
		if err != nil {
			return nil, fmt.Errorf("执行请求 '%s' 失败: %w", req.Name, err)
		}
		responses = append(responses, resp)

		// 请求之间添加一些延迟，避免过快请求服务器
		time.Sleep(200 * time.Millisecond)
	}

	return responses, nil
}

// createHTTPRequest 创建HTTP请求
func (e *Executor) createHTTPRequest(req *models.HTTPRequest) (*http.Request, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	// 验证URL
	if req.URL == nil {
		return nil, fmt.Errorf("无效的URL: 为空")
	}

	// 检查URL是否有效
	urlStr := req.URL.String()
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return nil, fmt.Errorf("无效的URL: %s（必须以http://或https://开头）", urlStr)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest(req.Method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for name, values := range req.Headers {
		for _, value := range values {
			httpReq.Header.Add(name, value)
		}
	}

	// 如果没有设置User-Agent，设置一个默认的
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", "jhttp/0.1.0")
	}

	return httpReq, nil
}

// PrintResponse 打印响应结果
func PrintResponse(resp *models.HTTPResponse) {
	if resp.Error != nil {
		fmt.Printf("请求失败: %v\n", resp.Error)
		fmt.Printf("请求耗时: %d ms\n", resp.Time)
		return
	}

	// 打印状态行
	fmt.Printf("HTTP/1.1 %s\n", resp.Status)

	// 打印响应头
	for name, values := range resp.Headers {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// 打印空行和响应体
	if len(resp.Body) > 0 {
		fmt.Println()

		// 检查Content-Type以决定如何格式化输出
		contentType := resp.Headers.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// 尝试美化JSON输出（这是简化版）
			fmt.Println(resp.BodyString)
		} else {
			fmt.Println(resp.BodyString)
		}
	}

	// 打印请求耗时
	fmt.Printf("\n请求耗时: %d ms\n", resp.Time)
}

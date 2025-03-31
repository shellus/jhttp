package parser

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/shellus/jhttp/internal/models"
)

// 正则表达式定义
var (
	// 请求行正则表达式：方法 + URL
	requestLineRegex = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|TRACE)\s+(.+)$`)

	// 请求名称正则表达式：以###开头的行
	requestNameRegex = regexp.MustCompile(`^###\s*(.+)$`)

	// 变量定义正则表达式：@变量名 = 变量值
	variableRegex = regexp.MustCompile(`^@(\w+)\s*=\s*(.+)$`)

	// 请求头正则表达式：头名称: 头值
	headerRegex = regexp.MustCompile(`^([^:]+):\s*(.+)$`)

	// 变量引用正则表达式：{{变量名}}
	variableRefRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)
)

// ParseFile 解析HTTP文件
func ParseFile(filePath string) (*models.HTTPFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	httpFile := models.NewHTTPFile(filePath)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	var currentRequest *models.HTTPRequest
	var isReadingBody bool
	var bodyBuilder strings.Builder
	var currentName string
	var currentDescription string
	var readingRequestComment bool      // 用于标记是否正在读取请求注释
	var foundEmptyLineAfterHeaders bool // 用于标记是否找到了请求头之后的空行

	// 逐行解析文件
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 跳过空行
		if line == "" {
			if currentRequest != nil && currentRequest.Method != "" {
				// 如果已经有了请求方法，那么这个空行可能是请求头和请求体的分隔
				if !isReadingBody && !foundEmptyLineAfterHeaders {
					isReadingBody = true
					foundEmptyLineAfterHeaders = true
					bodyBuilder.Reset()
				} else if isReadingBody {
					// 如果已经在读取请求体，那么空行也是请求体的一部分
					bodyBuilder.WriteString("\n")
				}
			}
			continue
		}

		// 处理请求名称行
		if matches := requestNameRegex.FindStringSubmatch(line); len(matches) > 1 {
			// 如果上一个请求还在处理中，保存其请求体
			if isReadingBody && currentRequest != nil {
				currentRequest.Body = strings.TrimSpace(bodyBuilder.String())
				isReadingBody = false
				bodyBuilder.Reset()
			}

			// 保存当前请求名（不再包含注释）
			currentName = matches[1]
			currentDescription = ""            // 重置描述
			readingRequestComment = true       // 标记正在读取请求注释
			foundEmptyLineAfterHeaders = false // 重置标记
			continue
		}

		// 处理注释行，现在注释内容不会添加到请求名中
		if strings.HasPrefix(line, "#") {
			if readingRequestComment {
				// 收集注释作为请求的描述（而不是名称的一部分）
				commentText := strings.TrimSpace(line[1:])
				if commentText != "" {
					if currentDescription != "" {
						currentDescription += "\n"
					}
					currentDescription += commentText
				}
			}
			continue
		}

		// 处理变量定义
		if matches := variableRegex.FindStringSubmatch(line); len(matches) > 2 {
			name, value := matches[1], matches[2]
			httpFile.GlobalVars[name] = value
			readingRequestComment = false // 变量定义不是请求注释
			continue
		}

		// 处理请求行（方法+URL）
		if matches := requestLineRegex.FindStringSubmatch(line); len(matches) > 2 {
			method, rawURL := matches[1], matches[2]
			parsedURL, err := url.Parse(rawURL)
			if err != nil {
				return nil, fmt.Errorf("行 %d: 无效的URL: %w", lineNum, err)
			}

			// 创建新请求
			currentRequest = &models.HTTPRequest{
				Name:        currentName,
				Description: currentDescription, // 保存收集的注释内容
				Method:      method,
				URL:         parsedURL,
				Headers:     make(http.Header),
				Variables:   make(map[string]string),
				LineNumber:  lineNum,
			}
			httpFile.AddRequest(currentRequest)

			// 重置状态
			currentName = ""
			currentDescription = ""
			readingRequestComment = false
			isReadingBody = false
			foundEmptyLineAfterHeaders = false
			continue
		}

		// 处理请求头
		if matches := headerRegex.FindStringSubmatch(line); len(matches) > 2 && currentRequest != nil && !isReadingBody {
			name, value := matches[1], matches[2]
			currentRequest.Headers.Add(name, value)
			readingRequestComment = false // 请求头不是请求注释
			continue
		}

		// 如果正在读取请求体
		if isReadingBody && currentRequest != nil {
			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
			readingRequestComment = false // 请求体不是请求注释
			continue
		}

		// 如果是JSON请求体的开始（{），那么进入请求体模式
		if currentRequest != nil && !isReadingBody &&
			(strings.HasPrefix(line, "{") || strings.HasPrefix(line, "[")) {
			isReadingBody = true
			foundEmptyLineAfterHeaders = true
			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
			continue
		}

		// 遇到未处理的行，重置请求注释状态
		readingRequestComment = false
	}

	// 处理最后一个请求的请求体
	if isReadingBody && currentRequest != nil {
		currentRequest.Body = strings.TrimSpace(bodyBuilder.String())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时发生错误: %w", err)
	}

	return httpFile, nil
}

// ResolveVariables 解析HTTP请求中的变量引用
func ResolveVariables(httpFile *models.HTTPFile, request *models.HTTPRequest, env string) (*models.HTTPRequest, error) {
	// 创建请求的副本
	resolvedReq := &models.HTTPRequest{
		Name:           request.Name,
		Method:         request.Method,
		Headers:        make(http.Header),
		Body:           request.Body,
		FormParameters: request.FormParameters,
		Variables:      make(map[string]string),
		LineNumber:     request.LineNumber,
	}

	// 复制并解析URL
	if request.URL != nil {
		// 获取原始URL字符串
		rawURLStr := request.URL.String()

		// URL解码，处理可能被编码的变量
		decodedURLStr, err := url.QueryUnescape(rawURLStr)
		if err != nil {
			// 如果解码失败，使用原始字符串
			decodedURLStr = rawURLStr
		}

		// 解析变量
		resolvedURLStr := resolveVariablesInString(httpFile, decodedURLStr, env)

		// 解析新的URL
		parsedURL, err := url.Parse(resolvedURLStr)
		if err != nil {
			return nil, fmt.Errorf("解析URL时发生错误: %w", err)
		}
		resolvedReq.URL = parsedURL
	}

	// 解析请求头中的变量
	for name, values := range request.Headers {
		for _, value := range values {
			resolvedValue := resolveVariablesInString(httpFile, value, env)
			resolvedReq.Headers.Add(name, resolvedValue)
		}
	}

	// 解析请求体中的变量
	resolvedReq.Body = resolveVariablesInString(httpFile, request.Body, env)

	return resolvedReq, nil
}

// resolveVariablesInString 解析字符串中的变量引用
func resolveVariablesInString(httpFile *models.HTTPFile, input string, env string) string {
	if input == "" {
		return input
	}

	// 使用正则表达式替换所有变量引用
	result := variableRefRegex.ReplaceAllStringFunc(input, func(match string) string {
		// 提取变量名
		varName := match[2 : len(match)-2]

		// 解析变量
		if value, found := httpFile.ResolveVariable(varName, env); found {
			return value
		}

		// 未找到变量时保留原样
		return match
	})

	return result
}

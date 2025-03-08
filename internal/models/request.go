package models

import (
	"net/http"
	"net/url"
)

// HTTPRequest 表示一个HTTP请求
type HTTPRequest struct {
	Name           string            // 请求名称
	Method         string            // HTTP方法 (GET, POST, PUT等)
	URL            *url.URL          // 请求URL
	Headers        http.Header       // 请求头
	Body           string            // 请求体内容
	FormParameters url.Values        // 表单参数
	Variables      map[string]string // 请求中定义的变量
	LineNumber     int               // 文件中的行号
}

// HTTPFile 表示解析后的HTTP文件
type HTTPFile struct {
	Path           string                 // 文件路径
	Requests       []*HTTPRequest         // 文件中的请求列表
	GlobalVars     map[string]string      // 全局变量
	EnvironmentVars map[string]map[string]string // 环境变量 [环境名][变量名]值
}

// HTTPResponse 表示HTTP响应
type HTTPResponse struct {
	StatusCode int               // 状态码
	Status     string            // 状态文本
	Headers    http.Header       // 响应头
	Body       []byte            // 响应体
	BodyString string            // 响应体字符串形式
	Time       int64             // 请求耗时(毫秒)
	Request    *HTTPRequest      // 原始请求
	Error      error             // 错误(如果有)
}

// NewHTTPFile 创建一个新的HTTP文件结构
func NewHTTPFile(path string) *HTTPFile {
	return &HTTPFile{
		Path:           path,
		Requests:       make([]*HTTPRequest, 0),
		GlobalVars:     make(map[string]string),
		EnvironmentVars: make(map[string]map[string]string),
	}
}

// AddRequest 向HTTP文件添加一个请求
func (f *HTTPFile) AddRequest(req *HTTPRequest) {
	f.Requests = append(f.Requests, req)
}

// FindRequestByName 通过名称查找请求
func (f *HTTPFile) FindRequestByName(name string) *HTTPRequest {
	for _, req := range f.Requests {
		if req.Name == name {
			return req
		}
	}
	return nil
}

// ResolveVariable 解析变量，支持环境变量替换
func (f *HTTPFile) ResolveVariable(name string, env string) (string, bool) {
	// 首先查找环境变量
	if env != "" {
		if envVars, ok := f.EnvironmentVars[env]; ok {
			if val, ok := envVars[name]; ok {
				return val, true
			}
		}
	}
	
	// 然后查找全局变量
	if val, ok := f.GlobalVars[name]; ok {
		return val, true
	}
	
	return "", false
} 
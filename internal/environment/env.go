package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// 环境文件可能的名称
var envFileNames = []string{
	"http-client.private.env.json",
	"http-client.env.json",
}

// LoadEnvFile 从环境文件中加载环境变量
func LoadEnvFile(filePath string, envName string) (map[string]string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取环境文件: %w", err)
	}

	// 解析JSON
	var environments map[string]map[string]string
	if err := json.Unmarshal(data, &environments); err != nil {
		return nil, fmt.Errorf("无法解析环境文件: %w", err)
	}

	// 查找指定的环境
	env, ok := environments[envName]
	if !ok {
		return nil, fmt.Errorf("环境 '%s' 不存在", envName)
	}

	return env, nil
}

// FindEnvFile 查找环境变量文件
// httpFilePath: HTTP文件的路径
// 返回值: 环境文件路径, 是否找到, 是否需要警告(只有在上级目录找到时才为true)
func FindEnvFile(httpFilePath string) (string, bool, bool) {
	// 获取HTTP文件所在的目录
	httpDir := filepath.Dir(httpFilePath)

	// 首先检查HTTP文件同目录下的环境文件
	for _, name := range envFileNames {
		envFilePath := filepath.Join(httpDir, name)
		if fileExists(envFilePath) {
			// 在当前目录找到，不需要警告
			return envFilePath, true, false
		}
	}

	// 如果在当前目录没有找到，向上查找
	current := httpDir
	for {
		// 向上一级目录
		parent := filepath.Dir(current)
		if parent == current {
			// 已经到达根目录，停止查找
			break
		}
		current = parent

		// 检查这两个文件名
		for _, name := range envFileNames {
			envFilePath := filepath.Join(current, name)
			if fileExists(envFilePath) {
				// 在上级目录找到，需要警告
				return envFilePath, true, true
			}
		}
	}

	// 没有找到环境文件
	return "", false, false
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ListEnvironments 列出环境文件中所有的环境名称
func ListEnvironments(filePath string) ([]string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取环境文件: %w", err)
	}

	// 解析JSON
	var environments map[string]map[string]string
	if err := json.Unmarshal(data, &environments); err != nil {
		return nil, fmt.Errorf("无法解析环境文件: %w", err)
	}

	// 提取所有环境名称
	envNames := make([]string, 0, len(environments))
	for name := range environments {
		envNames = append(envNames, name)
	}

	return envNames, nil
}

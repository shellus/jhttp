package environment

import (
	"encoding/json"
	"fmt"
	"os"
)

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

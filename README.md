# JHTTP - HTTP文件命令行执行工具

## 项目介绍

JHTTP是一个命令行工具，用于执行IntelliJ IDEA支持的.http文件。这些文件通常用于测试HTTP API，包含HTTP请求的详细信息。本工具允许用户在命令行环境中执行这些请求，无需启动IDE。

## 功能特点

- 支持解析和执行IntelliJ IDEA格式的.http文件
- 命令行界面，易于集成到自动化脚本和CI/CD流程中
- 支持各种HTTP方法（GET, POST, PUT, DELETE等）
- 支持设置请求头、请求体和查询参数
- 支持环境变量和请求之间的引用
- 提供丰富的输出格式选项

## 技术选择

本项目使用Go语言实现，具有以下优势：

- 性能优秀，低内存占用
- 编译为单个二进制文件，无需依赖关系
- 跨平台支持良好（Windows、macOS、Linux）
- 丰富的标准库，包括HTTP客户端和解析工具

## 安装方法

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/jhttp/jhttp.git
cd jhttp

# 构建项目
go build -o jhttp cmd/jhttp/main.go

# 将可执行文件移动到PATH目录
sudo mv jhttp /usr/local/bin/
```

### 使用Go工具安装

```bash
go install github.com/jhttp/jhttp/cmd/jhttp@latest
```

## 使用方法

```bash
# 基本用法
jhttp example.http

# 使用特定环境变量执行
jhttp --env-file http-client.private.env.json --env 开发环境 example.http

# 只执行特定请求
jhttp --request "Get User" example.http

# 保存响应到文件
jhttp --output response.json example.http

# 详细输出模式
jhttp --verbose example.http

# 列出文件中的所有请求
jhttp --list example.http
```

## 环境变量配置

本工具支持使用环境变量文件来简化请求中的参数配置和管理敏感信息。

### 环境变量文件格式

```json
{
    "开发环境": {
        "urlPrefix": "https://api.example.com",
        "username": "your_username",
        "password": "your_password",
        "Token": "your_access_token"
    },
    "测试环境": {
        "urlPrefix": "https://test-api.example.com",
        "username": "test_username",
        "password": "test_password",
        "Token": "test_access_token"
    }
}
```

### 安全注意事项

- 项目包含两个环境文件：
  - `http-client.env.example.json`: 示例文件，可以提交到版本控制
  - `http-client.private.env.json`: 包含实际敏感信息，**不应提交到版本控制**

- 请确保已将敏感环境文件添加到 `.gitignore` 中
- 详细说明请参考 [环境变量说明.md](环境变量说明.md)

## .http文件格式示例

```
### Get all users
GET {{urlPrefix}}/users
Accept: application/json

### Create new user
POST {{urlPrefix}}/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

### Get user by ID
GET {{urlPrefix}}/users/{{userId}}
```

## 项目结构

```
jhttp/
├── .gitignore                         # Git忽略文件
├── README.md                          # 项目说明文档
├── CONTRIBUTING.md                    # 贡献指南
├── 环境变量说明.md                      # 环境变量使用说明
├── http-client.env.example.json       # 环境变量示例文件
├── go.mod                             # Go模块定义
├── go.sum                             # 依赖版本锁定
├── example.http                       # 示例HTTP文件
├── cmd/
│   └── jhttp/
│       └── main.go                    # 应用入口点
├── internal/
│   ├── cli/
│   │   └── cli.go                     # 命令行参数处理
│   ├── parser/
│   │   └── parser.go                  # .http文件解析器
│   ├── executor/
│   │   └── executor.go                # 请求执行器
│   └── models/
│       └── request.go                 # 数据模型
```

## 贡献指南

1. Fork本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启Pull Request

详细的贡献指南请参考 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可证

MIT

## 联系方式

项目GitHub地址: https://github.com/jhttp/jhttp 
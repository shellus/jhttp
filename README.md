# JHTTP - HTTP文件命令行执行工具

## 项目介绍

JHTTP是一个命令行工具，用于执行IntelliJ IDEA支持的.http文件。这些文件通常用于测试HTTP API，包含HTTP请求的详细信息。本工具允许用户在命令行环境中执行这些请求，无需启动IDE。

此工具特别适合以下场景：
- 自动化测试API端点
- CI/CD流程中的API验证
- 脚本化执行HTTP请求
- 在无GUI环境中测试API

## 功能特点

- 支持解析和执行IntelliJ IDEA格式的.http文件
- 命令行界面，易于集成到自动化脚本和CI/CD流程中
- 支持各种HTTP方法（GET, POST, PUT, DELETE等）
- 支持设置请求头、请求体和查询参数
- 支持环境变量和请求之间的引用
- 提供详细的输出和错误信息
- 支持导出响应到文件
- 可指定执行单个请求或整个文件

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

### 基本命令

```bash
# 基本用法
jhttp example.http

# 使用特定环境变量执行
jhttp --env-file http-client.private.env.json --env 开发环境 example.http

# 只执行特定请求
jhttp --request "获取用户信息" example.http

# 保存响应到文件
jhttp --output response.json example.http

# 详细输出模式
jhttp --verbose example.http

# 列出文件中的所有请求
jhttp --list example.http
```

### 命令参数说明

| 参数 | 说明 |
|------|------|
| `--env-file <file>` | 指定环境变量文件路径 |
| `--env <name>` | 指定使用的环境名称 |
| `--request <name>` | 指定要执行的请求名称 |
| `--output <file>` | 指定响应输出文件 |
| `--verbose` | 输出详细信息 |
| `--version` | 显示版本信息 |
| `--help` | 显示帮助信息 |
| `--list` | 列出所有请求名称 |

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
  - `http-client.env.json`: 示例文件，可以提交到版本控制
  - `http-client.private.env.json`: 包含实际敏感信息，**不应提交到版本控制**

- 请确保已将敏感环境文件添加到 `.gitignore` 中

## .http文件格式示例

```
### 获取用户列表
GET {{urlPrefix}}/users
Accept: application/json

### 创建新用户
POST {{urlPrefix}}/users
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 30
}

### 获取特定用户
GET {{urlPrefix}}/users/{{userId}}
Accept: application/json
```

### 支持的语法特性

- 请求名称 (以`###`开头)
- 变量引用 (使用`{{变量名}}`格式)
- 环境变量 (以`@变量名 = 值`格式定义)
- 多行请求体
- 文件上传
- JSON, XML, 表单数据等多种内容类型

## 错误处理

JHTTP 提供了详细的错误信息，帮助用户快速定位问题：

- 请求错误（网络问题）
- 域名解析错误
- 连接被拒绝
- 请求超时
- SSL/TLS 证书验证错误
- HTTP 状态错误
- 解析错误

## 项目结构

```
jhttp/
├── .gitignore                         # Git忽略文件
├── README.md                          # 项目说明文档
├── http-client.env.json               # 环境变量示例文件
├── go.mod                             # Go模块定义
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
│   ├── environment/
│   │   └── env.go                     # 环境变量管理
│   └── models/
│       └── request.go                 # 数据模型
```

## 联系方式

项目GitHub地址: https://github.com/jhttp/jhttp 
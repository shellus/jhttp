#!/bin/bash

# 设置错误时退出
set -e

echo "开始安装 jhttp..."

# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 编译器，请先安装 Go"
    exit 1
fi

# 检查 Go 版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
MIN_VERSION="1.16"
if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo "警告: 当前 Go 版本 ($GO_VERSION) 可能过旧，建议使用 Go 1.16 或更高版本"
fi

# 编译程序
echo "正在编译 jhttp..."
if ! go build -o jhttp cmd/jhttp/main.go; then
    echo "错误: 编译失败"
    exit 1
fi

# 设置执行权限
echo "正在设置执行权限..."
chmod +x jhttp

# 检查 /usr/local/bin 目录是否存在
if [ ! -d "/usr/local/bin" ]; then
    echo "错误: /usr/local/bin 目录不存在"
    exit 1
fi

# 检查是否有写入权限
if [ ! -w "/usr/local/bin" ]; then
    echo "错误: 没有写入 /usr/local/bin 的权限，请使用 sudo 运行此脚本"
    exit 1
fi

# 移动文件
echo "正在安装到 /usr/local/bin..."
cp jhttp /usr/local/bin/

# 验证安装
if command -v jhttp &> /dev/null; then
    echo "安装成功！"
    echo "jhttp 版本: $(jhttp -version)"
else
    echo "警告: 安装可能未成功，请检查 PATH 环境变量"
fi 
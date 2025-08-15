#!/bin/bash

# clipboard-manager 构建脚本
# 用于编译 macOS 可执行文件

set -e

echo "🚀 开始构建 Clipboard Manager..."

# 项目信息
APP_NAME="clipboard-manager"
VERSION=$(git describe --tags --always 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GO_VERSION=$(go version | cut -d' ' -f3)

# 清理之前的构建
echo "🧹 清理旧的构建文件..."
rm -f ${APP_NAME}

# 安装依赖
echo "📦 安装依赖..."
go mod tidy

# 编译程序 - 现在编译所有Go文件
echo "🔨 编译程序..."
go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}'" -o ${APP_NAME} *.go

# 检查构建结果
if [ -f "${APP_NAME}" ]; then
    echo "✅ 构建成功！"
    echo "📁 可执行文件: ./${APP_NAME}"
    echo "📊 文件大小: $(du -h ${APP_NAME} | cut -f1)"
    echo ""
    echo "🎯 使用方法:"
    echo "   ./${APP_NAME}           # 运行程序"
    echo "   ./${APP_NAME} --version # 查看版本信息"
    echo "   sudo cp ${APP_NAME} /usr/local/bin/  # 安装到系统路径"
    echo ""
else
    echo "❌ 构建失败！"
    exit 1
fi

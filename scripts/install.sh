#!/bin/bash

# clipboard-manager 安装脚本
# 用于将程序安装到系统路径

set -e

APP_NAME="clipboard-manager"
INSTALL_PATH="/usr/local/bin"

echo "📋 Clipboard Manager 安装脚本"
echo "==============================="

# 检查可执行文件是否存在
if [ ! -f "${APP_NAME}" ]; then
    echo "❌ 错误: 找不到可执行文件 '${APP_NAME}'"
    echo "请先运行构建脚本: ./scripts/build.sh"
    exit 1
fi

# 检查安装权限
if [ ! -w "${INSTALL_PATH}" ]; then
    echo "🔐 需要管理员权限来安装到 ${INSTALL_PATH}"
    echo "请输入密码..."
    SUDO_CMD="sudo"
else
    SUDO_CMD=""
fi

# 安装文件
echo "📦 正在安装 ${APP_NAME} 到 ${INSTALL_PATH}..."
${SUDO_CMD} cp "${APP_NAME}" "${INSTALL_PATH}/"
${SUDO_CMD} chmod +x "${INSTALL_PATH}/${APP_NAME}"

# 验证安装
if command -v "${APP_NAME}" >/dev/null 2>&1; then
    echo "✅ 安装成功！"
    echo ""
    echo "🎯 现在你可以在任何地方使用以下命令:"
    echo "   ${APP_NAME}"
    echo ""
    echo "💡 提示:"
    echo "   - 添加别名: alias cb='${APP_NAME}'"
    echo "   - 查看帮助: ${APP_NAME} --help (如果实现了)"
    echo ""
else
    echo "❌ 安装可能失败，请检查 PATH 环境变量"
    exit 1
fi

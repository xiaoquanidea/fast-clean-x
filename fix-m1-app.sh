#!/bin/bash

# Fast Clean X - M1 Mac 修复脚本
# 用于解决 M1 Mac 上"应用程序无法打开"的问题

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Fast Clean X - M1 Mac 修复工具${NC}"
echo ""

# 检查是否提供了应用路径
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}用法: ./fix-m1-app.sh <应用路径>${NC}"
    echo ""
    echo "示例:"
    echo "  ./fix-m1-app.sh build/bin/Fast-Clean-X-arm64.app"
    echo "  ./fix-m1-app.sh ~/Downloads/Fast-Clean-X-arm64.app"
    echo ""
    exit 1
fi

APP_PATH="$1"

# 检查应用是否存在
if [ ! -d "$APP_PATH" ]; then
    echo -e "${RED}❌ 错误: 找不到应用 '$APP_PATH'${NC}"
    exit 1
fi

echo -e "${BLUE}📱 应用路径: $APP_PATH${NC}"
echo ""

# 步骤 1: 移除隔离属性
echo -e "${GREEN}步骤 1/3: 移除隔离属性 (quarantine)...${NC}"
xattr -cr "$APP_PATH"
echo -e "${GREEN}✅ 隔离属性已移除${NC}"
echo ""

# 步骤 2: 检查架构
echo -e "${GREEN}步骤 2/3: 检查应用架构...${NC}"
BINARY_PATH="$APP_PATH/Contents/MacOS/fast-clean-x"
if [ -f "$BINARY_PATH" ]; then
    ARCH=$(file "$BINARY_PATH" | grep -o "arm64\|x86_64")
    if [ "$ARCH" = "arm64" ]; then
        echo -e "${GREEN}✅ 应用架构: ARM64 (Apple Silicon)${NC}"
    elif [ "$ARCH" = "x86_64" ]; then
        echo -e "${YELLOW}⚠️  应用架构: x86_64 (Intel)${NC}"
        echo -e "${YELLOW}   在 M1 Mac 上运行需要 Rosetta 2${NC}"
        
        # 检查是否安装了 Rosetta 2
        if /usr/bin/pgrep -q oahd; then
            echo -e "${GREEN}✅ Rosetta 2 已安装${NC}"
        else
            echo -e "${YELLOW}⚠️  Rosetta 2 未安装，正在安装...${NC}"
            softwareupdate --install-rosetta --agree-to-license
            echo -e "${GREEN}✅ Rosetta 2 安装完成${NC}"
        fi
    else
        echo -e "${RED}❌ 无法识别应用架构${NC}"
    fi
else
    echo -e "${RED}❌ 找不到可执行文件${NC}"
fi
echo ""

# 步骤 3: 设置执行权限
echo -e "${GREEN}步骤 3/3: 设置执行权限...${NC}"
if [ -f "$BINARY_PATH" ]; then
    chmod +x "$BINARY_PATH"
    echo -e "${GREEN}✅ 执行权限已设置${NC}"
else
    echo -e "${RED}❌ 找不到可执行文件，跳过${NC}"
fi
echo ""

# 显示修复结果
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 修复完成！${NC}"
echo ""
echo -e "${BLUE}现在可以尝试以下方式打开应用：${NC}"
echo ""
echo -e "${YELLOW}方法 1: 双击打开${NC}"
echo "  直接双击应用图标"
echo ""
echo -e "${YELLOW}方法 2: 右键打开${NC}"
echo "  右键点击应用 → 打开 → 仍要打开"
echo ""
echo -e "${YELLOW}方法 3: 命令行打开${NC}"
echo "  open \"$APP_PATH\""
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 询问是否立即打开
read -p "是否立即打开应用？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}正在打开应用...${NC}"
    open "$APP_PATH"
fi


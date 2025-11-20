#!/bin/bash

# Fast Clean X 构建脚本
# 用于构建 macOS ARM64、AMD64 和 Windows AMD64 版本

set -e

echo "🚀 开始构建 Fast Clean X..."
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查是否安装了 MinGW-w64（用于构建 Windows 版本）
HAS_MINGW=false
if command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    HAS_MINGW=true
fi

# 清理旧的构建产物
echo -e "${BLUE}📦 清理旧的构建产物...${NC}"
rm -rf build/bin/*
mkdir -p build/bin

# 构建 macOS ARM64 版本 (Apple Silicon)
echo ""
echo -e "${GREEN}🍎 构建 macOS ARM64 版本 (Apple Silicon M1/M2/M3)...${NC}"
wails build -platform darwin/arm64

# 重命名输出文件
if [ -d "build/bin/Fast Clean X.app" ]; then
    mv "build/bin/Fast Clean X.app" build/bin/Fast-Clean-X-arm64.app
    echo -e "${GREEN}✅ ARM64 版本构建完成: build/bin/Fast-Clean-X-arm64.app${NC}"
fi

# 构建 macOS AMD64 版本 (Intel)
echo ""
echo -e "${GREEN}🍎 构建 macOS AMD64 版本 (Intel x86_64)...${NC}"
wails build -platform darwin/amd64

# 重命名输出文件
if [ -d "build/bin/Fast Clean X.app" ]; then
    mv "build/bin/Fast Clean X.app" build/bin/Fast-Clean-X-amd64.app
    echo -e "${GREEN}✅ AMD64 版本构建完成: build/bin/Fast-Clean-X-amd64.app${NC}"
fi

# 构建 Windows AMD64 版本
echo ""
if [ "$HAS_MINGW" = true ]; then
    echo -e "${GREEN}🪟 构建 Windows AMD64 版本 (Intel/AMD x86_64)...${NC}"
    wails build -platform windows/amd64

    # 重命名输出文件
    if [ -f "build/bin/fast-clean-x.exe" ]; then
        mv build/bin/fast-clean-x.exe build/bin/Fast-Clean-X-amd64.exe
        echo -e "${GREEN}✅ Windows AMD64 版本构建完成: build/bin/Fast-Clean-X-amd64.exe${NC}"
    else
        echo -e "${RED}❌ Windows 版本构建失败${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  跳过 Windows 版本构建（未安装 MinGW-w64）${NC}"
    echo -e "${YELLOW}    如需构建 Windows 版本，请先安装 MinGW-w64:${NC}"
    echo -e "${YELLOW}    brew install mingw-w64${NC}"
fi

# 显示构建结果
echo ""
echo -e "${BLUE}📊 构建完成！${NC}"
echo ""
echo "构建产物："
ls -lh build/bin/

# 计算文件大小
echo ""
if [ -d "build/bin/Fast-Clean-X-arm64.app" ]; then
    ARM64_SIZE=$(du -sh build/bin/Fast-Clean-X-arm64.app | cut -f1)
    echo -e "${GREEN}🍎 macOS ARM64: ${ARM64_SIZE}${NC}"
fi

if [ -d "build/bin/Fast-Clean-X-amd64.app" ]; then
    AMD64_SIZE=$(du -sh build/bin/Fast-Clean-X-amd64.app | cut -f1)
    echo -e "${GREEN}🍎 macOS AMD64: ${AMD64_SIZE}${NC}"
fi

if [ -f "build/bin/Fast-Clean-X-amd64.exe" ]; then
    WIN_SIZE=$(du -sh build/bin/Fast-Clean-X-amd64.exe | cut -f1)
    echo -e "${GREEN}🪟 Windows AMD64: ${WIN_SIZE}${NC}"
fi

echo ""
echo -e "${YELLOW}💡 提示：${NC}"
echo "  - macOS ARM64 版本适用于 Apple Silicon (M1/M2/M3) Mac"
echo "  - macOS AMD64 版本适用于 Intel 芯片 Mac"
if [ "$HAS_MINGW" = true ]; then
    echo "  - Windows AMD64 版本适用于 Windows 10/11 (64位)"
else
    echo "  - Windows 版本未构建（需要安装 MinGW-w64）"
fi
echo ""
echo -e "${GREEN}🎉 构建完成！${NC}"


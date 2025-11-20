#!/bin/bash

# Fast Clean X - 应用诊断工具
# 用于诊断 macOS 应用无法打开的问题

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Fast Clean X - 应用诊断工具${NC}"
echo ""

# 检查是否提供了应用路径
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}用法: ./diagnose-app.sh <应用路径>${NC}"
    echo ""
    echo "示例:"
    echo "  ./diagnose-app.sh build/bin/Fast-Clean-X-arm64.app"
    echo "  ./diagnose-app.sh ~/Downloads/Fast-Clean-X-arm64.app"
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

# 1. 检查系统架构
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}1. 系统信息${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
SYSTEM_ARCH=$(uname -m)
if [ "$SYSTEM_ARCH" = "arm64" ]; then
    echo -e "系统架构: ${GREEN}ARM64 (Apple Silicon)${NC}"
elif [ "$SYSTEM_ARCH" = "x86_64" ]; then
    echo -e "系统架构: ${YELLOW}x86_64 (Intel)${NC}"
else
    echo -e "系统架构: ${RED}$SYSTEM_ARCH (未知)${NC}"
fi

MACOS_VERSION=$(sw_vers -productVersion)
echo -e "macOS 版本: $MACOS_VERSION"
echo ""

# 2. 检查应用架构
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}2. 应用信息${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
BINARY_PATH="$APP_PATH/Contents/MacOS/fast-clean-x"

if [ -f "$BINARY_PATH" ]; then
    APP_ARCH=$(file "$BINARY_PATH" | grep -o "arm64\|x86_64")
    if [ "$APP_ARCH" = "arm64" ]; then
        echo -e "应用架构: ${GREEN}ARM64 (Apple Silicon)${NC}"
        if [ "$SYSTEM_ARCH" = "x86_64" ]; then
            echo -e "${RED}⚠️  警告: 应用是 ARM64，但系统是 Intel${NC}"
            echo -e "${YELLOW}   建议下载 Intel 版本 (Fast-Clean-X-amd64.app)${NC}"
        fi
    elif [ "$APP_ARCH" = "x86_64" ]; then
        echo -e "应用架构: ${YELLOW}x86_64 (Intel)${NC}"
        if [ "$SYSTEM_ARCH" = "arm64" ]; then
            echo -e "${YELLOW}⚠️  提示: 应用是 Intel 版本，在 M1 Mac 上需要 Rosetta 2${NC}"
            
            # 检查 Rosetta 2
            if /usr/bin/pgrep -q oahd; then
                echo -e "${GREEN}✅ Rosetta 2 已安装${NC}"
            else
                echo -e "${RED}❌ Rosetta 2 未安装${NC}"
                echo -e "${YELLOW}   运行: softwareupdate --install-rosetta${NC}"
            fi
        fi
    else
        echo -e "应用架构: ${RED}未知${NC}"
    fi
    
    # 检查可执行权限
    if [ -x "$BINARY_PATH" ]; then
        echo -e "可执行权限: ${GREEN}✅ 已设置${NC}"
    else
        echo -e "可执行权限: ${RED}❌ 未设置${NC}"
        echo -e "${YELLOW}   运行: chmod +x \"$BINARY_PATH\"${NC}"
    fi
else
    echo -e "${RED}❌ 找不到可执行文件: $BINARY_PATH${NC}"
fi
echo ""

# 3. 检查隔离属性
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}3. 安全属性${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
QUARANTINE=$(xattr -l "$APP_PATH" 2>/dev/null | grep "com.apple.quarantine" || echo "")
if [ -n "$QUARANTINE" ]; then
    echo -e "隔离属性: ${RED}❌ 存在 (这会导致无法打开)${NC}"
    echo -e "${YELLOW}   运行: xattr -cr \"$APP_PATH\"${NC}"
else
    echo -e "隔离属性: ${GREEN}✅ 不存在${NC}"
fi

# 检查其他扩展属性
OTHER_ATTRS=$(xattr -l "$APP_PATH" 2>/dev/null | grep -v "com.apple.quarantine" || echo "")
if [ -n "$OTHER_ATTRS" ]; then
    echo -e "其他属性: ${YELLOW}存在${NC}"
    echo "$OTHER_ATTRS" | sed 's/^/  /'
else
    echo -e "其他属性: ${GREEN}无${NC}"
fi
echo ""

# 4. 检查代码签名
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}4. 代码签名${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
CODESIGN_INFO=$(codesign -dv "$APP_PATH" 2>&1 || echo "未签名")
if echo "$CODESIGN_INFO" | grep -q "Signature=adhoc"; then
    echo -e "签名状态: ${YELLOW}临时签名 (adhoc)${NC}"
    echo -e "${YELLOW}   这是正常的开发构建，但可能被 Gatekeeper 阻止${NC}"
elif echo "$CODESIGN_INFO" | grep -q "Authority="; then
    echo -e "签名状态: ${GREEN}已签名${NC}"
    echo "$CODESIGN_INFO" | grep "Authority=" | sed 's/^/  /'
else
    echo -e "签名状态: ${RED}未签名${NC}"
fi
echo ""

# 5. 诊断总结
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}5. 诊断总结${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

ISSUES=0

# 检查架构匹配
if [ "$SYSTEM_ARCH" = "arm64" ] && [ "$APP_ARCH" = "x86_64" ]; then
    echo -e "${YELLOW}⚠️  架构不匹配: 建议下载 ARM64 版本${NC}"
    ISSUES=$((ISSUES + 1))
fi

# 检查隔离属性
if [ -n "$QUARANTINE" ]; then
    echo -e "${RED}❌ 隔离属性存在: 需要移除${NC}"
    ISSUES=$((ISSUES + 1))
fi

# 检查可执行权限
if [ -f "$BINARY_PATH" ] && [ ! -x "$BINARY_PATH" ]; then
    echo -e "${RED}❌ 缺少可执行权限${NC}"
    ISSUES=$((ISSUES + 1))
fi

# 检查 Rosetta 2
if [ "$SYSTEM_ARCH" = "arm64" ] && [ "$APP_ARCH" = "x86_64" ]; then
    if ! /usr/bin/pgrep -q oahd; then
        echo -e "${RED}❌ 需要安装 Rosetta 2${NC}"
        ISSUES=$((ISSUES + 1))
    fi
fi

if [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✅ 未发现问题，应用应该可以正常打开${NC}"
else
    echo -e "${YELLOW}发现 $ISSUES 个问题${NC}"
fi
echo ""

# 6. 修复建议
if [ $ISSUES -gt 0 ]; then
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}6. 修复建议${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}运行修复脚本:${NC}"
    echo -e "  ./fix-m1-app.sh \"$APP_PATH\""
    echo ""
    echo -e "${YELLOW}或手动修复:${NC}"
    if [ -n "$QUARANTINE" ]; then
        echo -e "  xattr -cr \"$APP_PATH\""
    fi
    if [ -f "$BINARY_PATH" ] && [ ! -x "$BINARY_PATH" ]; then
        echo -e "  chmod +x \"$BINARY_PATH\""
    fi
    if [ "$SYSTEM_ARCH" = "arm64" ] && [ "$APP_ARCH" = "x86_64" ]; then
        if ! /usr/bin/pgrep -q oahd; then
            echo -e "  softwareupdate --install-rosetta"
        fi
    fi
    echo ""
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"


#!/bin/bash

# Fast Clean X - 快速发布脚本
# 用于创建新版本并触发自动构建

set -e

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Fast Clean X - 发布工具${NC}"
echo ""

# 检查是否有未提交的更改
if [[ -n $(git status -s) ]]; then
    echo -e "${RED}❌ 错误: 有未提交的更改${NC}"
    echo ""
    git status -s
    echo ""
    echo -e "${YELLOW}请先提交所有更改：${NC}"
    echo "  git add ."
    echo "  git commit -m \"your message\""
    exit 1
fi

# 获取当前版本
CURRENT_VERSION=$(grep -o '"productVersion": "[^"]*"' wails.json | cut -d'"' -f4)
echo -e "${BLUE}当前版本: ${GREEN}v$CURRENT_VERSION${NC}"
echo ""

# 询问新版本号
echo -e "${YELLOW}请输入新版本号（如 1.0.0）:${NC}"
read -r NEW_VERSION

# 验证版本号格式
if ! [[ $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ 错误: 版本号格式不正确${NC}"
    echo -e "${YELLOW}   正确格式: 1.0.0${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}准备发布版本: v$NEW_VERSION${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 确认
echo -e "${YELLOW}是否继续？(y/n)${NC}"
read -r -n 1 CONFIRM
echo ""

if [[ ! $CONFIRM =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

echo ""
echo -e "${GREEN}步骤 1/5: 更新版本号...${NC}"

# 更新 wails.json
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/\"productVersion\": \".*\"/\"productVersion\": \"$NEW_VERSION\"/" wails.json
else
    # Linux
    sed -i "s/\"productVersion\": \".*\"/\"productVersion\": \"$NEW_VERSION\"/" wails.json
fi

echo -e "${GREEN}✅ 版本号已更新为: v$NEW_VERSION${NC}"
echo ""

# 提交更改
echo -e "${GREEN}步骤 2/5: 提交更改...${NC}"
git add wails.json
git commit -m "chore: bump version to $NEW_VERSION"
echo -e "${GREEN}✅ 更改已提交${NC}"
echo ""

# 创建标签
echo -e "${GREEN}步骤 3/5: 创建标签...${NC}"
git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"
echo -e "${GREEN}✅ 标签已创建: v$NEW_VERSION${NC}"
echo ""

# 推送到远程
echo -e "${GREEN}步骤 4/5: 推送到 GitHub...${NC}"
git push
git push origin "v$NEW_VERSION"
echo -e "${GREEN}✅ 已推送到 GitHub${NC}"
echo ""

# 完成
echo -e "${GREEN}步骤 5/5: 完成！${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 发布流程已启动！${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}GitHub Actions 正在自动构建...${NC}"
echo ""
echo -e "${BLUE}📦 查看构建进度:${NC}"
echo "   https://github.com/hutiquan/fast-clean-x/actions"
echo ""
echo -e "${BLUE}📦 查看 Release:${NC}"
echo "   https://github.com/hutiquan/fast-clean-x/releases/tag/v$NEW_VERSION"
echo ""
echo -e "${YELLOW}💡 提示:${NC}"
echo "   - 构建通常需要 5-10 分钟"
echo "   - 构建完成后会自动创建 Release"
echo "   - Release 包含 macOS 和 Windows 版本"
echo ""


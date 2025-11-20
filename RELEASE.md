# 发布指南

本文档说明如何发布 Fast Clean X 的新版本。

## 🚀 自动发布流程

项目使用 GitHub Actions 自动构建和发布。

### 发布步骤

#### 1. 更新版本号

编辑 `wails.json`，更新版本号：

```json
{
  "info": {
    "productVersion": "1.0.0"  // 修改这里
  }
}
```

#### 2. 提交更改

```bash
git add wails.json
git commit -m "chore: bump version to 1.0.0"
git push
```

#### 3. 创建并推送标签

```bash
# 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送标签到 GitHub
git push origin v1.0.0
```

#### 4. 自动构建

推送标签后，GitHub Actions 会自动：
1. ✅ 构建 macOS ARM64 版本
2. ✅ 构建 macOS AMD64 版本
3. ✅ 构建 Windows AMD64 版本
4. ✅ 创建 GitHub Release
5. ✅ 上传所有构建产物

#### 5. 查看 Release

访问 `https://github.com/你的用户名/fast-clean-x/releases` 查看发布。

## 📦 构建产物

每次发布会生成以下文件：

| 文件名 | 平台 | 架构 | 大小 |
|--------|------|------|------|
| `Fast-Clean-X-macOS-arm64.zip` | macOS | ARM64 | ~9 MB |
| `Fast-Clean-X-macOS-amd64.zip` | macOS | AMD64 | ~10 MB |
| `Fast-Clean-X-Windows-amd64.zip` | Windows | AMD64 | ~8 MB |

## 🔧 手动构建（可选）

如果需要手动构建：

```bash
# 构建所有平台
./build.sh

# 打包
cd build/bin
zip -r Fast-Clean-X-macOS-arm64.zip Fast-Clean-X-arm64.app
zip -r Fast-Clean-X-macOS-amd64.zip Fast-Clean-X-amd64.app
zip Fast-Clean-X-Windows-amd64.zip Fast-Clean-X-amd64.exe
```

## 📝 版本号规范

使用语义化版本号（Semantic Versioning）：

- **主版本号**（Major）：不兼容的 API 修改
- **次版本号**（Minor）：向下兼容的功能性新增
- **修订号**（Patch）：向下兼容的问题修正

示例：
- `v1.0.0` - 首个正式版本
- `v1.1.0` - 新增功能
- `v1.1.1` - 修复 bug
- `v2.0.0` - 重大更新

## 🏷️ 标签命名

- 正式版本：`v1.0.0`
- 预发布版本：`v1.0.0-beta.1`
- 候选版本：`v1.0.0-rc.1`

## 📋 发布检查清单

发布前确认：

- [ ] 所有测试通过
- [ ] 更新了版本号（`wails.json`）
- [ ] 更新了 `README.md`（如有必要）
- [ ] 更新了 `BUILD.md`（如有必要）
- [ ] 本地构建成功
- [ ] 应用可以正常启动
- [ ] 核心功能正常工作

## 🐛 故障排除

### Q: GitHub Actions 构建失败

**检查**：
1. 访问 `https://github.com/你的用户名/fast-clean-x/actions`
2. 查看失败的工作流
3. 检查错误日志

**常见问题**：
- Go 版本不匹配
- Node.js 依赖安装失败
- Wails 构建错误

### Q: Release 没有自动创建

**检查**：
1. 确认标签格式正确（`v*`）
2. 确认标签已推送到 GitHub
3. 检查 GitHub Actions 权限

### Q: 构建产物缺失

**检查**：
1. 查看 Actions 日志
2. 确认所有构建任务都成功
3. 检查 artifact 上传步骤

## 🔐 GitHub Token 权限

GitHub Actions 需要以下权限：

1. 访问仓库设置：`Settings` → `Actions` → `General`
2. 确保 `Workflow permissions` 设置为：
   - ✅ Read and write permissions
   - ✅ Allow GitHub Actions to create and approve pull requests

## 📚 相关文档

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [GitHub Releases 文档](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [语义化版本](https://semver.org/lang/zh-CN/)

## 🎯 快速发布命令

```bash
# 一键发布脚本
VERSION="1.0.0"

# 更新版本号（手动编辑 wails.json）
# vim wails.json

# 提交并打标签
git add wails.json
git commit -m "chore: bump version to $VERSION"
git push

git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

# 等待 GitHub Actions 自动构建和发布
echo "✅ 标签已推送，GitHub Actions 正在构建..."
echo "📦 访问 https://github.com/你的用户名/fast-clean-x/actions 查看进度"
```

## 🎉 发布后

1. 检查 Release 页面
2. 测试下载链接
3. 验证安装包可用
4. 更新文档（如有必要）
5. 通知用户（社交媒体、邮件等）

---

**祝发布顺利！** 🚀


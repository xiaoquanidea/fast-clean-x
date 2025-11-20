# 构建指南

## 📦 支持的平台

| 平台 | 架构 | 适用设备 | 状态 |
|------|------|----------|------|
| macOS | ARM64 | Apple Silicon (M1/M2/M3/M4) | ✅ 完全支持 |
| macOS | AMD64 | Intel 芯片 | ✅ 完全支持 |
| Windows | AMD64 | Intel/AMD 64位 | ✅ 支持（需测试） |

## 🚀 快速构建

### 一键构建（推荐）

```bash
./build.sh
```

**自动检测**：
- ✅ 始终构建 macOS ARM64 和 AMD64 版本
- ✅ 如果安装了 MinGW-w64，自动构建 Windows 版本
- ✅ 如果未安装 MinGW-w64，跳过 Windows 版本并提示

**生成文件**：
- `build/bin/Fast-Clean-X-arm64.app` - macOS Apple Silicon 版本
- `build/bin/Fast-Clean-X-amd64.app` - macOS Intel 版本
- `build/bin/Fast-Clean-X-amd64.exe` - Windows 64位 版本（如果安装了 MinGW-w64）

**构建 Windows 版本**（可选）：
```bash
# 安装 MinGW-w64
brew install mingw-w64

# 重新运行构建
./build.sh
```

### 手动构建

```bash
# macOS ARM64 版本（Apple Silicon）
wails build -platform darwin/arm64 -clean

# macOS AMD64 版本（Intel）
wails build -platform darwin/amd64 -clean

# Windows AMD64 版本
wails build -platform windows/amd64 -clean

# 当前平台版本
wails build
```

## 🔧 构建选项

| 命令 | 说明 |
|------|------|
| `wails build` | 构建当前平台版本 |
| `wails build -clean` | 清理后构建 |
| `wails build -upx` | 压缩可执行文件（减小 50% 体积） |
| `wails build -debug` | 调试模式构建 |
| `wails build -v` | 详细输出 |

## 🎯 环境要求

### 必需软件

- **Go**: 1.23+
- **Node.js**: 18+
- **Wails**: v2.11.0
- **Xcode Command Line Tools**: `xcode-select --install`

### 环境检查

```bash
# 检查 Wails 环境
wails doctor

# 检查版本
go version
node --version
wails version
```

## 🔍 交叉编译

### 在 macOS 上构建

#### ARM64 Mac (Apple Silicon)

```bash
# 本地构建（快）
wails build -platform darwin/arm64

# 交叉编译 Intel 版本（稍慢）
wails build -platform darwin/amd64

# 交叉编译 Windows 版本（需要 MinGW-w64）
wails build -platform windows/amd64
```

#### Intel Mac

```bash
# 本地构建（快）
wails build -platform darwin/amd64

# 交叉编译 ARM64 版本（稍慢）
wails build -platform darwin/arm64

# 交叉编译 Windows 版本（需要 MinGW-w64）
wails build -platform windows/amd64
```

### 在 Windows 上构建

```bash
# 本地构建
wails build

# 或指定平台
wails build -platform windows/amd64
```

**注意**: 在 Windows 上无法交叉编译 macOS 版本。

### 交叉编译 Windows 版本的要求

在 macOS 上构建 Windows 版本需要安装 MinGW-w64：

```bash
# 使用 Homebrew 安装
brew install mingw-w64
```

## 📦 打包分发

### 创建 ZIP

```bash
cd build/bin

# macOS 版本
zip -r Fast-Clean-X-macOS-arm64.zip Fast-Clean-X-arm64.app
zip -r Fast-Clean-X-macOS-amd64.zip Fast-Clean-X-amd64.app

# Windows 版本
zip Fast-Clean-X-Windows-amd64.zip Fast-Clean-X-amd64.exe
```

### 创建 DMG

```bash
# 安装工具
brew install create-dmg

# 创建 DMG
create-dmg \
  --volname "Fast Clean X" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --app-drop-link 600 185 \
  "Fast-Clean-X-arm64.dmg" \
  "build/bin/Fast-Clean-X-arm64.app"
```

## 🔐 代码签名（可选）

### 签名应用

```bash
# 查看证书
security find-identity -v -p codesigning

# 签名
codesign --deep --force --verify --verbose \
  --sign "Developer ID Application: Your Name (TEAM_ID)" \
  build/bin/Fast-Clean-X-arm64.app

# 验证
codesign --verify --deep --strict --verbose=2 \
  build/bin/Fast-Clean-X-arm64.app
```

### 公证（Notarization）

```bash
# 创建 ZIP
ditto -c -k --keepParent \
  build/bin/Fast-Clean-X-arm64.app \
  Fast-Clean-X-arm64.zip

# 提交公证
xcrun notarytool submit Fast-Clean-X-arm64.zip \
  --apple-id "your@email.com" \
  --team-id "TEAM_ID" \
  --password "app-specific-password" \
  --wait

# 装订票据
xcrun stapler staple build/bin/Fast-Clean-X-arm64.app
```

## 🐛 常见问题

### Q: "command not found: wails"

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Q: 前端构建失败

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..
wails build -clean
```

### Q: M1 Mac 提示"应用程序无法打开"

这是 macOS Gatekeeper 安全机制导致的，有以下解决方案：

**方法 1: 使用修复脚本（推荐）**
```bash
./fix-m1-app.sh build/bin/Fast-Clean-X-arm64.app
```

**方法 2: 手动修复**
```bash
# 1. 移除隔离属性
xattr -cr build/bin/Fast-Clean-X-arm64.app

# 2. 设置执行权限
chmod +x build/bin/Fast-Clean-X-arm64.app/Contents/MacOS/fast-clean-x

# 3. 验证修复
xattr -l build/bin/Fast-Clean-X-arm64.app
```

**方法 3: 右键打开**
- 右键点击应用 → 打开 → 仍要打开
- 首次打开后，后续可以正常双击打开

**方法 4: 系统设置**
- 系统设置 → 隐私与安全性 → 安全性
- 找到被阻止的应用 → 点击"仍要打开"

### Q: "无法验证开发者"

```bash
# 移除隔离属性
xattr -cr build/bin/Fast-Clean-X-arm64.app
```

### Q: 交叉编译失败

```bash
# ARM64 Mac 构建 Intel 版本需要 Rosetta 2
softwareupdate --install-rosetta

# Intel Mac 构建 ARM64 版本需要 Xcode Command Line Tools
xcode-select --install

# macOS 构建 Windows 版本需要 MinGW-w64
brew install mingw-w64
```

### Q: Windows 版本构建失败

**在 macOS 上**:
```bash
# 安装 MinGW-w64
brew install mingw-w64

# 验证安装
x86_64-w64-mingw32-gcc --version
```

**在 Windows 上**:
- 确保安装了 Go 1.23+
- 确保安装了 Node.js 18+
- 确保安装了 Wails CLI
- 运行 `wails doctor` 检查环境

## 📝 构建检查清单

**构建前**:
- [ ] Go >= 1.23
- [ ] Node.js >= 18
- [ ] Wails = v2.11.0
- [ ] 前端依赖已安装
- [ ] `wails doctor` 通过

**构建后**:
- [ ] 应用可以启动
- [ ] 扫描功能正常
- [ ] 清理功能正常
- [ ] 配置持久化正常
- [ ] 文件大小 < 50MB

## 🚀 自动化构建

### GitHub Actions

创建 `.github/workflows/build.yml`:

```yaml
name: Build

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Build
        run: |
          chmod +x build.sh
          ./build.sh

      - name: Upload Artifacts
        uses: actions/upload-artifact@v3
        with:
          name: Fast-Clean-X
          path: build/bin/*.app
```

## 📚 相关资源

- [Wails 官方文档](https://wails.io/docs/introduction)
- [Wails 构建指南](https://wails.io/docs/guides/building)
- [macOS 代码签名](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)

---

**使用 `./build.sh` 一键构建两个平台版本！** 🚀


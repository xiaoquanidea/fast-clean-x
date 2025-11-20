# GitHub Actions 详解

## 📚 GitHub Actions 基础概念

### 什么是 GitHub Actions？

GitHub Actions 是 GitHub 提供的 **CI/CD（持续集成/持续部署）自动化平台**，可以让你：
- 自动化构建、测试、部署流程
- 在代码推送、PR、发布等事件时自动执行任务
- 使用 GitHub 提供的云服务器（runners）运行任务
- 免费使用（公开仓库无限制，私有仓库有配额）

### 核心概念

1. **Workflow（工作流）**: 自动化流程的配置文件（YAML 格式）
2. **Event（事件）**: 触发 workflow 的条件（如 push、pull_request、tag）
3. **Job（任务）**: workflow 中的一组步骤，可以并行或串行执行
4. **Step（步骤）**: job 中的单个操作
5. **Action（动作）**: 可复用的步骤模块
6. **Runner（运行器）**: 执行 workflow 的服务器

---

## 📄 项目 Workflows 详解

本项目有两个 workflow 文件：
- `build.yml` - 测试构建工作流
- `release.yml` - 发布工作流

---

## 1️⃣ build.yml - 测试构建工作流

这是一个**手动触发的测试构建**，用于验证项目能否正常构建。

### 完整配置

```yaml
name: Build

on:
  workflow_dispatch:  # 仅允许手动触发

permissions:
  contents: read

jobs:
  test-build:
    name: Test Build
    runs-on: macos-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Install frontend dependencies
        run: |
          cd frontend
          npm install

      - name: Run Wails doctor
        run: wails doctor

      - name: Build for current platform
        run: wails build

      - name: Check build output
        run: |
          ls -lh build/bin/
          du -sh build/bin/*
```

### 逐行解析

#### 基本配置

```yaml
name: Build
```
- **作用**: 工作流的名称，会显示在 GitHub Actions 页面

```yaml
on:
  workflow_dispatch:  # 仅允许手动触发
```
- **作用**: 定义触发条件
- `workflow_dispatch`: 只能通过 GitHub 网页手动点击 "Run workflow" 按钮触发

```yaml
permissions:
  contents: read
```
- **作用**: 定义 workflow 的权限
- `contents: read`: 只读权限，可以读取代码但不能修改

#### Job 配置

```yaml
jobs:
  test-build:
    name: Test Build
    runs-on: macos-latest
```
- **jobs**: 定义所有任务
- **test-build**: 任务的 ID
- **name**: 任务显示名称
- **runs-on**: 指定运行环境，`macos-latest` 表示使用最新版 macOS 虚拟机

#### 步骤详解

**步骤 1: Checkout code**
```yaml
- name: Checkout code
  uses: actions/checkout@v4
```
- 检出代码到 runner
- 使用官方的 checkout action（v4 版本）
- 将仓库代码下载到虚拟机

**步骤 2: Setup Go**
```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'
```
- 安装 Go 环境
- 使用官方的 Go 安装 action（v5 版本）
- 指定 Go 版本为 1.23

**步骤 3: Setup Node.js**
```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: '18'
```
- 安装 Node.js 环境（前端需要）
- 指定 Node.js 版本为 18

**步骤 4: Install Wails**
```yaml
- name: Install Wails
  run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
- 安装 Wails CLI 工具
- `run`: 执行 shell 命令
- 使用 `go install` 安装最新版 Wails

**步骤 5: Install frontend dependencies**
```yaml
- name: Install frontend dependencies
  run: |
    cd frontend
    npm install
```
- 安装前端依赖
- `|`: 多行命令
- 进入 frontend 目录并执行 `npm install`

**步骤 6: Run Wails doctor**
```yaml
- name: Run Wails doctor
  run: wails doctor
```
- 检查 Wails 环境配置
- 诊断工具，确保所有依赖都正确安装

**步骤 7: Build for current platform**
```yaml
- name: Build for current platform
  run: wails build
```
- 构建应用
- 为当前平台（macOS）构建应用

**步骤 8: Check build output**
```yaml
- name: Check build output
  run: |
    ls -lh build/bin/
    du -sh build/bin/*
```
- 检查构建产物
- 列出构建文件并显示大小

---

## 2️⃣ release.yml - 发布工作流

这是**自动发布流程**，当推送版本标签时自动构建多平台版本并创建 GitHub Release。

### 基本配置

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'  # 当推送 v* 标签时触发（如 v1.0.0）
  workflow_dispatch:  # 也允许手动触发

permissions:
  contents: write  # 允许创建 Release 和上传文件
```

#### 触发条件
- **自动触发**: 当推送以 `v` 开头的标签时（如 `v1.0.0`、`v2.1.3`）
- **手动触发**: 也可以通过网页手动触发

#### 权限
- `contents: write`: 需要写权限来创建 Release 和上传文件

---

### Job 1: build-macos - 构建 macOS 版本

```yaml
jobs:
  build-macos:
    name: Build macOS
    runs-on: macos-latest
```
- 在 macOS 环境中构建（因为需要 macOS SDK）

#### 步骤详解

**步骤 1-5**: 与 build.yml 相同
- Checkout code
- Setup Go (1.23)
- Setup Node.js (18)
- Install Wails
- Install frontend dependencies

**步骤 6: Build macOS ARM64**
```yaml
- name: Build macOS ARM64
  run: wails build -platform darwin/arm64
```
- 构建 ARM64 版本
- 适用于 Apple Silicon（M1/M2/M3/M4）芯片

**步骤 7: Build macOS AMD64**
```yaml
- name: Build macOS AMD64
  run: wails build -platform darwin/amd64
```
- 构建 AMD64 版本
- 适用于 Intel 芯片

**步骤 8: Rename and package**
```yaml
- name: Rename and package
  run: |
    # 重命名 ARM64
    mv "build/bin/Fast Clean X.app" build/bin/Fast-Clean-X-arm64.app
    
    # 创建 ARM64 ZIP
    cd build/bin
    zip -r Fast-Clean-X-macOS-arm64.zip Fast-Clean-X-arm64.app
    
    # 构建 AMD64
    cd ../..
    wails build -platform darwin/amd64
    mv "build/bin/Fast Clean X.app" build/bin/Fast-Clean-X-amd64.app
    
    # 创建 AMD64 ZIP
    cd build/bin
    zip -r Fast-Clean-X-macOS-amd64.zip Fast-Clean-X-amd64.app
```
- 重命名 .app 文件（区分架构）
- 创建 ZIP 压缩包
- 对两个架构分别处理

**步骤 9-10: Upload artifacts**
```yaml
- name: Upload ARM64 artifact
  uses: actions/upload-artifact@v4
  with:
    name: macos-arm64
    path: build/bin/Fast-Clean-X-macOS-arm64.zip

- name: Upload AMD64 artifact
  uses: actions/upload-artifact@v4
  with:
    name: macos-amd64
    path: build/bin/Fast-Clean-X-macOS-amd64.zip
```
- 上传构建产物为 artifacts（临时存储）
- 后续的 release job 会下载这些 artifacts

---

### Job 2: build-windows - 构建 Windows 版本

```yaml
build-windows:
  name: Build Windows
  runs-on: ubuntu-latest
```
- 在 Ubuntu 上构建
- 使用交叉编译构建 Windows 版本（更快、更便宜）

#### 步骤详解

**步骤 1-3**: 基础环境
- Checkout code
- Setup Go (1.23)
- Setup Node.js (18)

**步骤 4: Install dependencies**
```yaml
- name: Install dependencies
  run: |
    sudo apt-get update
    sudo apt-get install -y gcc-mingw-w64-x86-64
```
- 安装交叉编译工具
- `gcc-mingw-w64`: MinGW 编译器，用于在 Linux 上编译 Windows 程序

**步骤 5-6**: 安装 Wails 和前端依赖
- Install Wails
- Install frontend dependencies

**步骤 7: Build Windows AMD64**
```yaml
- name: Build Windows AMD64
  run: wails build -platform windows/amd64
```
- 构建 Windows 64位版本

**步骤 8: Package Windows**
```yaml
- name: Package Windows
  run: |
    cd build/bin
    mv fast-clean-x.exe Fast-Clean-X-amd64.exe
    zip Fast-Clean-X-Windows-amd64.zip Fast-Clean-X-amd64.exe
```
- 重命名 exe 文件
- 创建 ZIP 压缩包

**步骤 9: Upload Windows artifact**
```yaml
- name: Upload Windows artifact
  uses: actions/upload-artifact@v4
  with:
    name: windows-amd64
    path: build/bin/Fast-Clean-X-Windows-amd64.zip
```
- 上传 Windows 构建产物

---

### Job 3: release - 创建 GitHub Release

```yaml
release:
  name: Create Release
  needs: [build-macos, build-windows]
  runs-on: ubuntu-latest
```
- **needs**: 依赖关系，必须等 build-macos 和 build-windows 都完成后才执行
- 实现**串行执行**（前两个 job 是并行的）

#### 步骤详解

**步骤 1: Checkout code**
```yaml
- name: Checkout code
  uses: actions/checkout@v4
```
- 检出代码（获取 README 等文件）

**步骤 2: Download all artifacts**
```yaml
- name: Download all artifacts
  uses: actions/download-artifact@v4
  with:
    path: artifacts
```
- 下载前面两个 job 上传的所有构建产物
- 保存到 `artifacts` 目录

**步骤 3: Display structure**
```yaml
- name: Display structure
  run: ls -R artifacts
```
- 显示目录结构
- 调试用，查看下载的文件

**步骤 4: Create Release**
```yaml
- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    files: |
      artifacts/macos-arm64/Fast-Clean-X-macOS-arm64.zip
      artifacts/macos-amd64/Fast-Clean-X-macOS-amd64.zip
      artifacts/windows-amd64/Fast-Clean-X-Windows-amd64.zip
    draft: false
    prerelease: false
    generate_release_notes: true
    body: |
      ## 🎉 Fast Clean X ${{ github.ref_name }}
      
      快速清理项目构建文件的桌面应用。
      
      ### 📦 下载
      
      #### macOS
      - **Apple Silicon (M1/M2/M3/M4)**: `Fast-Clean-X-macOS-arm64.zip`
      - **Intel**: `Fast-Clean-X-macOS-amd64.zip`
      
      #### Windows
      - **64位**: `Fast-Clean-X-Windows-amd64.zip`
      
      ### 🚀 安装说明
      
      #### macOS
      1. 下载对应版本的 ZIP 文件
      2. 解压得到 `.app` 文件
      3. 如遇到"无法打开"提示，运行：
         ```bash
         xattr -cr /path/to/Fast-Clean-X-arm64.app
         ```
      4. 双击打开应用
      
      #### Windows
      1. 下载 ZIP 文件
      2. 解压得到 `.exe` 文件
      3. 双击运行
      
      ### ✨ 主要功能
      - 🔍 智能扫描多种项目类型
      - 📁 项目分组显示
      - 🎯 选择性清理
      - 📊 空间统计
      - ⚡ 高性能并发扫描
      
      ### 📝 更新日志
      
      查看下方自动生成的更新日志。
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**参数说明**:
- `files`: 要上传的文件列表（3个平台的 ZIP 文件）
- `draft: false`: 不是草稿，直接发布
- `prerelease: false`: 不是预发布版本
- `generate_release_notes: true`: 自动生成更新日志（基于 commits）
- `body`: Release 说明（支持 Markdown）
- `${{ github.ref_name }}`: 变量，获取标签名（如 v1.0.0）
- `GITHUB_TOKEN`: GitHub 自动提供的认证令牌

---

## 🔄 完整发布流程示例

### 发布新版本的完整流程

#### 1. 开发者本地操作
```bash
# 创建版本标签
git tag v1.0.0

# 推送标签到 GitHub
git push origin v1.0.0
```

#### 2. GitHub Actions 自动执行

```
触发 release.yml workflow
    ↓
┌───────────────────────────────────┐
│   并行执行（同时进行）              │
├─────────────────┬─────────────────┤
│  build-macos    │  build-windows  │
│  ├─ ARM64       │  └─ AMD64       │
│  └─ AMD64       │                 │
└─────────────────┴─────────────────┘
    ↓
上传 artifacts (3个文件)
    ↓
release job 开始
    ↓
下载所有 artifacts
    ↓
创建 GitHub Release
    ↓
上传 3 个 ZIP 文件
    ↓
生成发布说明
    ↓
✅ 完成！
```

#### 3. 用户可以
- 在 GitHub Releases 页面下载对应平台的安装包
- 查看自动生成的更新日志
- 阅读安装说明

---

## 💡 关键特性

### 1. 并行构建
```yaml
jobs:
  build-macos:    # 并行执行
  build-windows:  # 并行执行
  release:
    needs: [build-macos, build-windows]  # 等待前两个完成
```
- macOS 和 Windows 构建同时进行，节省时间
- release 等待两者都完成后再执行

### 2. 跨平台编译
- **macOS**: 在 macOS runner 上构建（需要 macOS SDK）
- **Windows**: 在 Ubuntu 上交叉编译（更快、更便宜）

### 3. Artifacts 传递机制
```
build-macos    → upload-artifact → 
build-windows  → upload-artifact → release → download-artifact
```
- 构建产物通过 artifacts 在 jobs 之间传递
- artifacts 会在 workflow 完成后保留一段时间

### 4. 自动化发布
- 推送标签 → 自动构建 → 自动发布
- 无需手动操作
- 支持多平台、多架构

### 5. 版本管理
- 使用 Git 标签管理版本
- 标签名自动用于 Release 标题
- 自动生成更新日志（基于 commits）

---

## 🎯 使用指南

### 测试构建
1. 进入 GitHub 仓库
2. 点击 **Actions** 标签
3. 选择 **Build** workflow
4. 点击 **Run workflow** 按钮
5. 等待构建完成

### 发布新版本
```bash
# 1. 确保代码已提交
git add .
git commit -m "Release v1.0.0"
git push

# 2. 创建并推送标签
git tag v1.0.0
git push origin v1.0.0

# 3. 等待 GitHub Actions 自动完成
# 4. 在 Releases 页面查看发布结果
```

### 查看构建进度
1. GitHub 仓库 → **Actions** 标签
2. 点击对应的 workflow run
3. 查看每个 job 的执行状态
4. 点击 job 查看详细日志

### 下载发布文件
1. GitHub 仓库 → **Releases** 标签
2. 选择对应版本
3. 下载对应平台的 ZIP 文件

---

## 📊 Workflow 对比

| 特性 | build.yml | release.yml |
|------|-----------|-------------|
| **触发方式** | 仅手动触发 | 标签推送 + 手动 |
| **目的** | 测试构建 | 正式发布 |
| **平台** | 仅 macOS | macOS + Windows |
| **架构** | 当前平台 | ARM64 + AMD64 |
| **产物** | 不保存 | 上传到 Release |
| **权限** | 只读 | 读写 |
| **执行时间** | ~5-10分钟 | ~15-20分钟 |

---

## 🔧 常见问题

### Q1: 为什么 macOS 要构建两个版本？
**A**: 因为 Apple 有两种芯片架构：
- **ARM64**: Apple Silicon (M1/M2/M3/M4)
- **AMD64**: Intel 芯片

### Q2: 为什么 Windows 在 Ubuntu 上构建？
**A**: 
- 使用交叉编译，速度更快
- Ubuntu runner 比 Windows runner 便宜
- MinGW 工具链成熟稳定

### Q3: 什么是 artifacts？
**A**: 
- 临时存储的构建产物
- 在 jobs 之间传递文件
- 默认保留 90 天

### Q4: 如何修改版本号？
**A**: 
- 版本号由 Git 标签决定
- 推送 `v1.2.3` 标签，版本号就是 `v1.2.3`

### Q5: 构建失败怎么办？
**A**: 
1. 查看 Actions 页面的错误日志
2. 检查依赖是否正确安装
3. 本地运行 `wails doctor` 诊断
4. 使用 `build.yml` 测试构建

---

## 📚 相关资源

- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [Wails 官方文档](https://wails.io/)
- [YAML 语法参考](https://yaml.org/)
- [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)

---

## 🎓 总结

本项目的 GitHub Actions 配置实现了：

✅ **自动化构建**: 推送标签即可自动构建多平台版本  
✅ **并行执行**: macOS 和 Windows 同时构建，节省时间  
✅ **跨平台支持**: 支持 macOS (ARM64/AMD64) 和 Windows (AMD64)  
✅ **自动发布**: 自动创建 Release 并上传安装包  
✅ **版本管理**: 基于 Git 标签的版本控制  
✅ **用户友好**: 详细的安装说明和自动生成的更新日志  

这是一个非常专业和完善的 CI/CD 配置！🚀


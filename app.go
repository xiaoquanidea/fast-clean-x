package main

import (
	"context"
	"fast-clean-x/backend/cleaner"
	"fast-clean-x/backend/config"
	"fast-clean-x/backend/models"
	"fast-clean-x/backend/scanner"
	"fmt"
	"os/exec"
	"runtime"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	configManager  *config.Manager
	currentScanner *scanner.Scanner
	currentCleaner *cleaner.Cleaner
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		configManager: config.GetManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetConfig 获取配置
func (a *App) GetConfig() *models.Config {
	return a.configManager.GetConfig()
}

// UpdateConfig 更新配置
func (a *App) UpdateConfig(cfg *models.Config) error {
	return a.configManager.UpdateConfig(cfg)
}

// AddScanPath 添加扫描路径
func (a *App) AddScanPath(path string) error {
	return a.configManager.AddScanPath(path)
}

// RemoveScanPath 移除扫描路径
func (a *App) RemoveScanPath(path string) error {
	return a.configManager.RemoveScanPath(path)
}

// AddIgnorePattern 添加忽略模式
func (a *App) AddIgnorePattern(pattern string) error {
	return a.configManager.AddIgnorePattern(pattern)
}

// RemoveIgnorePattern 移除忽略模式
func (a *App) RemoveIgnorePattern(pattern string) error {
	return a.configManager.RemoveIgnorePattern(pattern)
}

// UpdateScanRule 更新扫描规则
func (a *App) UpdateScanRule(ruleName string, enabled bool) error {
	return a.configManager.UpdateScanRule(ruleName, enabled)
}

// StartScan 开始扫描
func (a *App) StartScan() (*models.ScanResult, error) {
	cfg := a.configManager.GetConfig()

	// 创建扫描器
	a.currentScanner = scanner.New(
		a.configManager.GetEnabledRules(),
		cfg.IgnorePatterns,
		cfg.GlobalPathExcludes,
	)

	// 启动进度监听
	go a.listenScanProgress()

	// 执行扫描
	result, err := a.currentScanner.Scan(cfg.ScanPaths)

	// 关闭扫描器
	a.currentScanner.Close()
	a.currentScanner = nil

	return result, err
}

// listenScanProgress 监听扫描进度
func (a *App) listenScanProgress() {
	if a.currentScanner == nil {
		return
	}

	for progress := range a.currentScanner.GetProgressChan() {
		wailsRuntime.EventsEmit(a.ctx, "scan:progress", progress)
	}
}

// CancelScan 取消扫描
func (a *App) CancelScan() {
	if a.currentScanner != nil {
		a.currentScanner.Cancel()
	}
}

// StartClean 开始清理
func (a *App) StartClean(items []models.ScanItem) error {
	// 创建清理器
	a.currentCleaner = cleaner.New()

	// 启动进度监听
	go a.listenCleanProgress()

	// 执行清理
	err := a.currentCleaner.Clean(items)

	// 关闭清理器
	a.currentCleaner.Close()
	a.currentCleaner = nil

	return err
}

// listenCleanProgress 监听清理进度
func (a *App) listenCleanProgress() {
	if a.currentCleaner == nil {
		return
	}

	for progress := range a.currentCleaner.GetProgressChan() {
		wailsRuntime.EventsEmit(a.ctx, "clean:progress", progress)
	}
}

// CancelClean 取消清理
func (a *App) CancelClean() {
	if a.currentCleaner != nil {
		a.currentCleaner.Cancel()
	}
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择扫描目录",
	})
}

// OpenFolder 在文件管理器中打开文件夹
func (a *App) OpenFolder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

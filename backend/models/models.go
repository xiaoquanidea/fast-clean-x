package models

import "time"

// ScanRule 扫描规则
type ScanRule struct {
	Name        string   `json:"name"`        // 规则名称，如 "Maven"
	Description string   `json:"description"` // 规则描述
	TargetDirs  []string `json:"targetDirs"`  // 要扫描的目录名，如 ["target"]
	Enabled     bool     `json:"enabled"`     // 是否启用
}

// Config 应用配置
type Config struct {
	ScanPaths      []string   `json:"scanPaths"`      // 扫描路径列表
	IgnorePatterns []string   `json:"ignorePatterns"` // 忽略的项目路径模式
	ScanRules      []ScanRule `json:"scanRules"`      // 扫描规则
	LastScanTime   time.Time  `json:"lastScanTime"`   // 上次扫描时间
}

// ScanItem 扫描到的单个项目
type ScanItem struct {
	Path         string    `json:"path"`         // 完整路径
	ProjectPath  string    `json:"projectPath"`  // 项目根路径
	ProjectName  string    `json:"projectName"`  // 项目名称
	Type         string    `json:"type"`         // 类型，如 "maven", "gradle", "node"
	Size         int64     `json:"size"`         // 大小（字节）
	SizeReadable string    `json:"sizeReadable"` // 可读的大小，如 "1.2 GB"
	FileCount    int       `json:"fileCount"`    // 文件数量
	LastModified time.Time `json:"lastModified"` // 最后修改时间
	Selected     bool      `json:"selected"`     // 是否选中（用于删除）
}

// ScanResult 扫描结果
type ScanResult struct {
	Items      []ScanItem `json:"items"`      // 扫描到的项目列表
	TotalSize  int64      `json:"totalSize"`  // 总大小
	TotalCount int        `json:"totalCount"` // 总数量
	ScanTime   time.Time  `json:"scanTime"`   // 扫描时间
}

// ScanProgress 扫描进度
type ScanProgress struct {
	CurrentPath  string `json:"currentPath"`  // 当前扫描路径
	ScannedCount int    `json:"scannedCount"` // 已扫描数量
	TotalSize    int64  `json:"totalSize"`    // 已发现的总大小
	IsScanning   bool   `json:"isScanning"`   // 是否正在扫描
	Progress     int    `json:"progress"`     // 进度百分比 (0-100)
}

// CleanProgress 清理进度
type CleanProgress struct {
	CurrentPath   string `json:"currentPath"`   // 当前清理路径
	CleanedCount  int    `json:"cleanedCount"`  // 已清理数量
	TotalCount    int    `json:"totalCount"`    // 总数量
	CleanedSize   int64  `json:"cleanedSize"`   // 已清理大小
	IsCleaning    bool   `json:"isCleaning"`    // 是否正在清理
	Progress      int    `json:"progress"`      // 进度百分比 (0-100)
	FailedItems   []string `json:"failedItems"` // 清理失败的项目
}

// DefaultScanRules 返回默认的扫描规则
func DefaultScanRules() []ScanRule {
	return []ScanRule{
		{
			Name:        "Maven",
			Description: "Java Maven 构建目录",
			TargetDirs:  []string{"target"},
			Enabled:     true,
		},
		{
			Name:        "Gradle",
			Description: "Java Gradle 构建目录",
			TargetDirs:  []string{"build", ".gradle"},
			Enabled:     true,
		},
		{
			Name:        "Java IDE",
			Description: "Java IDE 输出目录",
			TargetDirs:  []string{"out"},
			Enabled:     true,
		},
		{
			Name:        "Node.js",
			Description: "Node.js 依赖和构建目录",
			TargetDirs:  []string{
				"node_modules",  // 依赖目录
				"dist",          // 通用构建输出
				"build",         // Create React App 等
				".next",         // Next.js
				".nuxt",         // Nuxt.js
				"out",           // Next.js 静态导出
				".output",       // Nuxt 3
				".vite",         // Vite 缓存
				".turbo",        // Turborepo 缓存
				".cache",        // Gatsby, Parcel 等
				".parcel-cache", // Parcel
				"coverage",      // 测试覆盖率
				".nyc_output",   // NYC 覆盖率
			},
			Enabled:     true,
		},
		{
			Name:        "Python",
			Description: "Python 缓存和虚拟环境",
			TargetDirs:  []string{"__pycache__", ".venv", "venv", ".pytest_cache", ".mypy_cache"},
			Enabled:     true,
		},
		{
			Name:        "Rust",
			Description: "Rust 构建目录",
			TargetDirs:  []string{"target"},
			Enabled:     true,
		},
		{
			Name:        "Go",
			Description: "Go vendor 目录",
			TargetDirs:  []string{"vendor"},
			Enabled:     false, // 默认不启用
		},
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ScanPaths:      []string{},
		IgnorePatterns: []string{},
		ScanRules:      DefaultScanRules(),
		LastScanTime:   time.Time{},
	}
}


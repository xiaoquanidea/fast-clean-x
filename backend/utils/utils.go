package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FormatSize 格式化文件大小为可读格式
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetHomeDir 获取用户主目录
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// GetConfigDir 获取配置目录
func GetConfigDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".fast-clean-x")
	
	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	
	return configDir, nil
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// IsHiddenFile 判断是否为隐藏文件或目录
func IsHiddenFile(path string) bool {
	name := filepath.Base(path)
	return strings.HasPrefix(name, ".")
}

// ShouldSkipDir 判断是否应该跳过该目录
func ShouldSkipDir(path string) bool {
	// 跳过隐藏目录
	if IsHiddenFile(path) {
		return true
	}
	
	// 跳过系统目录
	systemDirs := []string{
		"System Volume Information",
		"$RECYCLE.BIN",
		"Windows",
		"Program Files",
		"Program Files (x86)",
		"ProgramData",
		"Library",
		"System",
		"Applications",
	}
	
	name := filepath.Base(path)
	for _, sysDir := range systemDirs {
		if strings.EqualFold(name, sysDir) {
			return true
		}
	}
	
	return false
}

// MatchPattern 检查路径是否匹配模式
func MatchPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		// 简单的通配符匹配
		if strings.Contains(path, pattern) {
			return true
		}
		// 使用 filepath.Match 进行更精确的匹配
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}
	}
	return false
}

// FindProjectRoot 从构建目录向上查找项目根目录
func FindProjectRoot(buildPath string) string {
	// 从构建目录的父目录开始向上查找
	current := filepath.Dir(buildPath)
	var projectRoots []string

	// 向上查找所有包含项目标识的目录
	for {
		parent := filepath.Dir(current)
		if parent == current {
			// 已到达根目录
			break
		}

		// 检查当前目录是否有项目标识
		if hasProjectMarker(current) {
			projectRoots = append(projectRoots, current)
		}

		current = parent
	}

	// 如果找到多个项目根，返回最顶层的（优先选择包含 .git 的）
	if len(projectRoots) > 0 {
		// 从最顶层开始检查
		for i := len(projectRoots) - 1; i >= 0; i-- {
			root := projectRoots[i]
			// 优先返回包含 .git 的目录
			gitPath := filepath.Join(root, ".git")
			if _, err := os.Stat(gitPath); err == nil {
				return root
			}
		}
		// 如果没有 .git，返回最顶层的项目根
		return projectRoots[len(projectRoots)-1]
	}

	// 如果没找到任何项目标识，返回构建目录的父目录
	return filepath.Dir(buildPath)
}

// GetProjectName 从路径中提取项目名称
func GetProjectName(projectPath string) string {
	return filepath.Base(projectPath)
}

// hasProjectMarker 检查目录是否包含项目标识
func hasProjectMarker(dir string) bool {
	markers := []string{
		".git",
		".svn",
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
		"settings.gradle",
		"settings.gradle.kts",
		"package.json",
		"Cargo.toml",
		"go.mod",
	}

	for _, marker := range markers {
		markerPath := filepath.Join(dir, marker)
		if _, err := os.Stat(markerPath); err == nil {
			return true
		}
	}

	return false
}

// CalculateDirSize 计算目录大小
func CalculateDirSize(path string) (int64, int, error) {
	var size int64
	var count int
	
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			// 忽略权限错误等
			return nil
		}
		
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		
		return nil
	})
	
	return size, count, err
}

// PathExists 检查路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory 检查路径是否为目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}


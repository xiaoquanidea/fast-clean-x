package scanner

import (
	"context"
	"fast-clean-x/backend/models"
	"fast-clean-x/backend/utils"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Scanner 扫描器
type Scanner struct {
	rules          []models.ScanRule
	ignorePatterns []string
	progressChan   chan models.ScanProgress
	mu             sync.Mutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// New 创建新的扫描器
func New(rules []models.ScanRule, ignorePatterns []string) *Scanner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scanner{
		rules:          rules,
		ignorePatterns: ignorePatterns,
		progressChan:   make(chan models.ScanProgress, 100),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Scan 扫描指定路径
func (s *Scanner) Scan(paths []string) (*models.ScanResult, error) {
	result := &models.ScanResult{
		Items:      make([]models.ScanItem, 0),
		TotalSize:  0,
		TotalCount: 0,
		ScanTime:   time.Now(),
	}

	var wg sync.WaitGroup
	itemsChan := make(chan models.ScanItem, 100)
	
	// 启动结果收集器
	go func() {
		for item := range itemsChan {
			s.mu.Lock()
			result.Items = append(result.Items, item)
			result.TotalSize += item.Size
			result.TotalCount++
			s.mu.Unlock()
		}
	}()

	// 扫描每个路径
	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			s.scanPath(p, itemsChan)
		}(path)
	}

	// 等待所有扫描完成
	wg.Wait()
	close(itemsChan)

	return result, nil
}

// scanPath 扫描单个路径
func (s *Scanner) scanPath(rootPath string, itemsChan chan<- models.ScanItem) {
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// 检查是否取消
		select {
		case <-s.ctx.Done():
			return filepath.SkipDir
		default:
		}

		if err != nil {
			return nil // 忽略错误，继续扫描
		}

		if !info.IsDir() {
			return nil
		}

		// 跳过系统目录和隐藏目录
		if utils.ShouldSkipDir(path) {
			return filepath.SkipDir
		}

		// 检查是否匹配忽略模式
		if utils.MatchPattern(path, s.ignorePatterns) {
			return filepath.SkipDir
		}

		// 检查是否匹配扫描规则
		dirName := filepath.Base(path)
		for _, rule := range s.rules {
			for _, targetDir := range rule.TargetDirs {
				if dirName == targetDir {
					// 找到匹配的目录
					item := s.createScanItem(path, rule.Name)
					if item != nil {
						itemsChan <- *item
						
						// 发送进度更新
						s.sendProgress(path)
					}
					return filepath.SkipDir
				}
			}
		}

		return nil
	})
}

// createScanItem 创建扫描项
func (s *Scanner) createScanItem(path string, ruleType string) *models.ScanItem {
	// 计算目录大小
	size, fileCount, err := utils.CalculateDirSize(path)
	if err != nil {
		return nil
	}

	// 获取最后修改时间
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	// 查找项目根目录
	projectPath := utils.FindProjectRoot(path)
	projectName := utils.GetProjectName(projectPath)

	return &models.ScanItem{
		Path:         path,
		ProjectPath:  projectPath,
		ProjectName:  projectName,
		Type:         ruleType,
		Size:         size,
		SizeReadable: utils.FormatSize(size),
		FileCount:    fileCount,
		LastModified: info.ModTime(),
		Selected:     true, // 默认选中
	}
}

// sendProgress 发送进度更新
func (s *Scanner) sendProgress(currentPath string) {
	select {
	case s.progressChan <- models.ScanProgress{
		CurrentPath: currentPath,
		IsScanning:  true,
	}:
	default:
		// 如果通道满了，跳过这次更新
	}
}

// GetProgressChan 获取进度通道
func (s *Scanner) GetProgressChan() <-chan models.ScanProgress {
	return s.progressChan
}

// Cancel 取消扫描
func (s *Scanner) Cancel() {
	s.cancel()
}

// Close 关闭扫描器
func (s *Scanner) Close() {
	s.cancel()
	close(s.progressChan)
}


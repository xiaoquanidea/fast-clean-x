package cleaner

import (
	"context"
	"fast-clean-x/backend/models"
	"os"
	"sync"
)

// Cleaner 清理器
type Cleaner struct {
	progressChan chan models.CleanProgress
	mu           sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
}

// New 创建新的清理器
func New() *Cleaner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cleaner{
		progressChan: make(chan models.CleanProgress, 100),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Clean 清理指定的项目
func (c *Cleaner) Clean(items []models.ScanItem) error {
	totalCount := len(items)
	cleanedCount := 0
	var cleanedSize int64
	failedItems := make([]string, 0)

	for i, item := range items {
		// 检查是否取消
		select {
		case <-c.ctx.Done():
			return context.Canceled
		default:
		}

		// 只清理选中的项目
		if !item.Selected {
			continue
		}

		// 发送进度更新
		c.sendProgress(models.CleanProgress{
			CurrentPath:  item.Path,
			CleanedCount: cleanedCount,
			TotalCount:   totalCount,
			CleanedSize:  cleanedSize,
			IsCleaning:   true,
			Progress:     (i * 100) / totalCount,
			FailedItems:  failedItems,
		})

		// 删除目录
		err := os.RemoveAll(item.Path)
		if err != nil {
			failedItems = append(failedItems, item.Path)
		} else {
			cleanedCount++
			cleanedSize += item.Size
		}
	}

	// 发送最终进度
	c.sendProgress(models.CleanProgress{
		CurrentPath:  "",
		CleanedCount: cleanedCount,
		TotalCount:   totalCount,
		CleanedSize:  cleanedSize,
		IsCleaning:   false,
		Progress:     100,
		FailedItems:  failedItems,
	})

	return nil
}

// sendProgress 发送进度更新
func (c *Cleaner) sendProgress(progress models.CleanProgress) {
	select {
	case c.progressChan <- progress:
	default:
		// 如果通道满了，跳过这次更新
	}
}

// GetProgressChan 获取进度通道
func (c *Cleaner) GetProgressChan() <-chan models.CleanProgress {
	return c.progressChan
}

// Cancel 取消清理
func (c *Cleaner) Cancel() {
	c.cancel()
}

// Close 关闭清理器
func (c *Cleaner) Close() {
	c.cancel()
	close(c.progressChan)
}


package config

import (
	"encoding/json"
	"fast-clean-x/backend/models"
	"fast-clean-x/backend/utils"
	"os"
	"sync"
)

// Manager 配置管理器
type Manager struct {
	config *models.Config
	mu     sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取配置管理器单例
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			config: models.DefaultConfig(),
		}
		// 尝试加载配置
		_ = instance.Load()
	})
	return instance
}

// Load 从文件加载配置
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}

	// 如果配置文件不存在，使用默认配置
	if !utils.PathExists(configPath) {
		m.config = models.DefaultConfig()
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	m.config = &config
	return nil
}

// Save 保存配置到文件
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveInternal()
}

// saveInternal 内部保存方法，不加锁
func (m *Manager) saveInternal() error {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *models.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本以避免并发修改
	configCopy := *m.config
	return &configCopy
}

// UpdateConfig 更新配置
func (m *Manager) UpdateConfig(config *models.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config
	return m.saveInternal()
}

// AddScanPath 添加扫描路径
func (m *Manager) AddScanPath(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	for _, p := range m.config.ScanPaths {
		if p == path {
			return nil
		}
	}

	m.config.ScanPaths = append(m.config.ScanPaths, path)
	return m.saveInternal()
}

// RemoveScanPath 移除扫描路径
func (m *Manager) RemoveScanPath(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	newPaths := make([]string, 0)
	for _, p := range m.config.ScanPaths {
		if p != path {
			newPaths = append(newPaths, p)
		}
	}

	m.config.ScanPaths = newPaths
	return m.saveInternal()
}

// AddIgnorePattern 添加忽略模式
func (m *Manager) AddIgnorePattern(pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	for _, p := range m.config.IgnorePatterns {
		if p == pattern {
			return nil
		}
	}

	m.config.IgnorePatterns = append(m.config.IgnorePatterns, pattern)
	return m.saveInternal()
}

// RemoveIgnorePattern 移除忽略模式
func (m *Manager) RemoveIgnorePattern(pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	newPatterns := make([]string, 0)
	for _, p := range m.config.IgnorePatterns {
		if p != pattern {
			newPatterns = append(newPatterns, p)
		}
	}

	m.config.IgnorePatterns = newPatterns
	return m.saveInternal()
}

// UpdateScanRule 更新扫描规则
func (m *Manager) UpdateScanRule(ruleName string, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.config.ScanRules {
		if m.config.ScanRules[i].Name == ruleName {
			m.config.ScanRules[i].Enabled = enabled
			return m.saveInternal()
		}
	}

	return nil
}

// GetEnabledRules 获取启用的扫描规则
func (m *Manager) GetEnabledRules() []models.ScanRule {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rules := make([]models.ScanRule, 0)
	for _, rule := range m.config.ScanRules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}

	return rules
}

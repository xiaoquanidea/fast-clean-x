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

	// 合并默认配置，确保新字段有默认值
	m.config = m.mergeWithDefaults(&config)
	return nil
}

// mergeWithDefaults 将加载的配置与默认配置合并
// 这样可以确保：
// 1. 旧配置文件缺少的新字段会使用默认值
// 2. 用户的数据（扫描路径、规则启用状态）会被保留
// 3. 用户不需要手动删除配置文件
func (m *Manager) mergeWithDefaults(loaded *models.Config) *models.Config {
	defaults := models.DefaultConfig()

	// 保留用户配置的路径和时间
	defaults.ScanPaths = loaded.ScanPaths
	defaults.IgnorePatterns = loaded.IgnorePatterns
	defaults.LastScanTime = loaded.LastScanTime

	// 如果旧配置有 GlobalPathExcludes，保留它；否则使用默认值
	if len(loaded.GlobalPathExcludes) > 0 {
		defaults.GlobalPathExcludes = loaded.GlobalPathExcludes
	}

	// 合并规则：保留用户的 enabled 状态，但使用默认的其他字段
	ruleEnabledMap := make(map[string]bool)
	for _, rule := range loaded.ScanRules {
		ruleEnabledMap[rule.Name] = rule.Enabled
	}

	for i := range defaults.ScanRules {
		if enabled, exists := ruleEnabledMap[defaults.ScanRules[i].Name]; exists {
			defaults.ScanRules[i].Enabled = enabled
		}
	}

	return defaults
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

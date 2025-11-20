package scanner

import (
	"fast-clean-x/backend/models"
	"testing"
)

func TestSelectBestRule(t *testing.T) {
	s := &Scanner{}

	tests := []struct {
		name         string
		path         string
		rules        []models.ScanRule
		expectedName string
	}{
		{
			name: "优先级高的规则优先（都没有验证）",
			path: "/Users/xiao/project/build",
			rules: []models.ScanRule{
				{Name: "Gradle", Priority: 80, RequireMarkers: false},
				{Name: "Node.js", Priority: 100, RequireMarkers: false},
			},
			expectedName: "Node.js",
		},
		{
			name: "有验证的规则优先于无验证的规则",
			path: "/Users/xiao/project/build",
			rules: []models.ScanRule{
				{Name: "Node.js", Priority: 100, RequireMarkers: false},
				{Name: "Gradle", Priority: 80, RequireMarkers: true, ProjectMarkers: []string{"build.gradle"}},
			},
			expectedName: "Gradle", // 有验证的优先
		},
		{
			name: "优先级相同时，都有验证，选择第一个",
			path: "/Users/xiao/project/target",
			rules: []models.ScanRule{
				{Name: "Maven", Priority: 80, RequireMarkers: true, ProjectMarkers: []string{"pom.xml"}},
				{Name: "Rust", Priority: 80, RequireMarkers: true, ProjectMarkers: []string{"Cargo.toml"}},
			},
			expectedName: "Maven", // 优先级相同，选择第一个
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.selectBestRule(tt.path, tt.rules)
			if result.Name != tt.expectedName {
				t.Errorf("selectBestRule() = %v, want %v", result.Name, tt.expectedName)
			}
		})
	}
}

func TestShouldExcludeByPath(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		rule               models.ScanRule
		globalPathExcludes []string
		expected           bool
	}{
		{
			name: "全局排除规则生效",
			path: "/Users/xiao/project/node_modules/package/build",
			rule: models.ScanRule{
				Name:              "Gradle",
				ExcludeFromGlobal: false,
			},
			globalPathExcludes: []string{"node_modules"},
			expected:           true,
		},
		{
			name: "路径不包含任何排除项",
			path: "/Users/xiao/project/build",
			rule: models.ScanRule{
				Name:              "Gradle",
				ExcludeFromGlobal: false,
			},
			globalPathExcludes: []string{"node_modules"},
			expected:           false,
		},
		{
			name: "智能上下文检测：node_modules 下的 __pycache__",
			path: "/Users/xiao/project/node_modules/flatted/python/__pycache__",
			rule: models.ScanRule{
				Name:              "Python",
				TargetDirs:        []string{"__pycache__", ".venv", "venv"},
				ExcludeFromGlobal: true,
			},
			globalPathExcludes: []string{"node_modules", ".venv", "venv"},
			expected:           true, // node_modules 不是自己的目标目录，仍然被排除
		},
		{
			name: "智能上下文检测：项目根目录的 __pycache__",
			path: "/Users/xiao/project/__pycache__",
			rule: models.ScanRule{
				Name:              "Python",
				TargetDirs:        []string{"__pycache__", ".venv", "venv"},
				ExcludeFromGlobal: true,
			},
			globalPathExcludes: []string{"node_modules", ".venv", "venv"},
			expected:           false, // 不在 node_modules 下，保留
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				globalPathExcludes: tt.globalPathExcludes,
			}
			result := s.shouldExcludeByPath(tt.path, tt.rule)
			if result != tt.expected {
				t.Errorf("shouldExcludeByPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

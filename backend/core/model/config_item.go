package model

import (
	"time"
)

// ConfigCategory 配置分类
type ConfigCategory string

const (
	ConfigCategoryPlatform   ConfigCategory = "platform"   // 平台设置
	ConfigCategoryAI         ConfigCategory = "ai"         // AI智能
	ConfigCategoryMonitoring ConfigCategory = "monitoring" // 监控告警
	ConfigCategoryIntegration ConfigCategory = "integration" // 集成
)

// ConfigItem 配置项
type ConfigItem struct {
	ID          int64     `json:"id"`
	Category    string    `json:"category"`    // 一级分类: platform/ai/monitoring/integration
	SubCategory string    `json:"sub_category"` // 二级分类
	Key         string    `json:"key"`         // 配置键
	Value       string    `json:"value"`       // 配置值
	Type        string    `json:"type"`        // 类型: string/number/boolean/json/password
	Options     string    `json:"options"`     // 可选项(JSON数组)
	DefaultVal  string    `json:"default_val"` // 默认值
	Required    bool      `json:"required"`    // 是否必填
	Sensitive   bool      `json:"sensitive"`   // 是否敏感(密码等)
	Description string    `json:"description"` // 描述
	SortOrder   int       `json:"sort_order"`  // 排序
	Icon        string    `json:"icon"`        // 图标
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConfigGroup 配置分组(用于前端展示)
type ConfigGroup struct {
	Category    string        `json:"category"`
	SubCategory string        `json:"sub_category"`
	Label       string        `json:"label"`
	Icon        string        `json:"icon"`
	Description string        `json:"description"`
	Items       []ConfigItem  `json:"items"`
}

// ConfigValue 配置值(用户设置的实际值)
type ConfigValue struct {
	ID        int64     `json:"id"`
	ConfigID  int64     `json:"config_id"`
	Value     string    `json:"value"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigTree 配置树形结构
type ConfigTree struct {
	Key         string        `json:"key"`
	Label       string        `json:"label"`
	Icon        string        `json:"icon"`
	Description string        `json:"description"`
	Children    []ConfigTree  `json:"children,omitempty"`
}

package repository

import (
	"core/model"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

// ConfigRepository 配置仓库
type ConfigRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewConfigRepository 创建配置仓库
func NewConfigRepository(db *sql.DB, logger *zap.Logger) *ConfigRepository {
	return &ConfigRepository{
		db:     db,
		logger: logger,
	}
}

// InitConfigItems 初始化配置项定义
func (r *ConfigRepository) InitConfigItems() error {
	// 检查是否已初始化
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM config_items").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Info("Config items already initialized", zap.Int("count", count))
		return nil
	}

	// 定义所有配置项
	items := []model.ConfigItem{
		// ========== Platform Settings ==========
		// 基础配置 - 只保留系统名称
		{Category: "platform", SubCategory: "basic", Key: "app.name", Value: "SRE AI Platform", Type: "string", Required: true, Description: "系统名称", SortOrder: 1, Icon: "AppstoreOutlined"},

		// 系统配置
		{Category: "platform", SubCategory: "system", Key: "system.session_timeout", Value: "1440", Type: "number", Required: true, Description: "会话超时时间(分钟)", SortOrder: 1, Icon: "FieldTimeOutlined"},
		{Category: "platform", SubCategory: "system", Key: "system.cache_ttl", Value: "300", Type: "number", Required: true, Description: "缓存TTL(秒)", SortOrder: 2, Icon: "DatabaseOutlined"},
		{Category: "platform", SubCategory: "system", Key: "system.log_level", Value: "info", Type: "string", Options: `["debug","info","warn","error"]`, Description: "日志级别", SortOrder: 3, Icon: "FileTextOutlined"},
		{Category: "platform", SubCategory: "system", Key: "system.max_upload_size", Value: "10485760", Type: "number", Description: "最大上传文件大小(字节)", SortOrder: 4, Icon: "UploadOutlined"},
		{Category: "platform", SubCategory: "system", Key: "system.enable_registration", Value: "false", Type: "boolean", Description: "允许用户注册", SortOrder: 5, Icon: "UserAddOutlined"},

		// 通知配置
		{Category: "platform", SubCategory: "notification", Key: "notification.email_enabled", Value: "false", Type: "boolean", Description: "启用邮件通知", SortOrder: 1, Icon: "MailOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.email_host", Value: "", Type: "string", Description: "SMTP服务器地址", SortOrder: 2, Icon: "MailOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.email_port", Value: "587", Type: "number", Description: "SMTP端口", SortOrder: 3, Icon: "MailOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.email_user", Value: "", Type: "string", Description: "SMTP用户名", SortOrder: 4, Icon: "UserOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.email_password", Value: "", Type: "password", Sensitive: true, Description: "SMTP密码", SortOrder: 5, Icon: "LockOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.webhook_enabled", Value: "false", Type: "boolean", Description: "启用Webhook通知", SortOrder: 6, Icon: "LinkOutlined"},
		{Category: "platform", SubCategory: "notification", Key: "notification.webhook_url", Value: "", Type: "string", Description: "Webhook URL", SortOrder: 7, Icon: "LinkOutlined"},

		// ========== AI & Intelligence ==========
		// 分析策略
		{Category: "ai", SubCategory: "strategy", Key: "ai.analysis_timeout", Value: "120", Type: "number", Required: true, Description: "AI分析超时时间(秒)", SortOrder: 1, Icon: "ClockCircleOutlined"},
		{Category: "ai", SubCategory: "strategy", Key: "ai.auto_analysis", Value: "true", Type: "boolean", Description: "告警自动触发AI分析", SortOrder: 2, Icon: "ThunderboltOutlined"},
		{Category: "ai", SubCategory: "strategy", Key: "ai.confidence_threshold", Value: "70", Type: "number", Description: "置信度阈值(%)", SortOrder: 3, Icon: "SafetyOutlined"},
		{Category: "ai", SubCategory: "strategy", Key: "ai.include_metrics", Value: "true", Type: "boolean", Description: "分析时包含指标数据", SortOrder: 4, Icon: "LineChartOutlined"},
		{Category: "ai", SubCategory: "strategy", Key: "ai.include_logs", Value: "false", Type: "boolean", Description: "分析时包含日志数据", SortOrder: 5, Icon: "FileTextOutlined"},

		// 对话设置
		{Category: "ai", SubCategory: "chat", Key: "ai.chat_context_length", Value: "10", Type: "number", Description: "对话上下文长度(轮数)", SortOrder: 1, Icon: "MessageOutlined"},
		{Category: "ai", SubCategory: "chat", Key: "ai.chat_enable_memory", Value: "true", Type: "boolean", Description: "启用对话记忆", SortOrder: 2, Icon: "SaveOutlined"},
		{Category: "ai", SubCategory: "chat", Key: "ai.system_prompt", Value: "You are an expert SRE AI assistant.", Type: "string", Description: "系统提示词", SortOrder: 3, Icon: "RobotOutlined"},
		{Category: "ai", SubCategory: "chat", Key: "ai.welcome_message", Value: "Hello! I am your SRE AI assistant. How can I help you today?", Type: "string", Description: "欢迎消息", SortOrder: 4, Icon: "SmileOutlined"},

		// ========== Monitoring ==========
		// Prometheus设置
		{Category: "monitoring", SubCategory: "prometheus", Key: "prometheus.server_url", Value: "http://prometheus:9090", Type: "string", Required: true, Description: "Prometheus服务器地址", SortOrder: 0, Icon: "GlobalOutlined"},
		{Category: "monitoring", SubCategory: "prometheus", Key: "prometheus.refresh_interval", Value: "30", Type: "number", Required: true, Description: "数据刷新间隔(秒)", SortOrder: 1, Icon: "ReloadOutlined"},
		{Category: "monitoring", SubCategory: "prometheus", Key: "prometheus.query_timeout", Value: "10", Type: "number", Description: "查询超时(秒)", SortOrder: 2, Icon: "ClockCircleOutlined"},
		{Category: "monitoring", SubCategory: "prometheus", Key: "prometheus.max_data_points", Value: "1000", Type: "number", Description: "最大数据点数量", SortOrder: 3, Icon: "DatabaseOutlined"},

		// 告警默认配置
		{Category: "monitoring", SubCategory: "alert", Key: "alert.default_severity", Value: "warning", Type: "string", Options: `["critical","warning","info"]`, Description: "默认告警级别", SortOrder: 1, Icon: "WarningOutlined"},
		{Category: "monitoring", SubCategory: "alert", Key: "alert.auto_resolve", Value: "true", Type: "boolean", Description: "告警自动恢复", SortOrder: 2, Icon: "CheckCircleOutlined"},
		{Category: "monitoring", SubCategory: "alert", Key: "alert.grouping_enabled", Value: "true", Type: "boolean", Description: "启用告警分组", SortOrder: 3, Icon: "GroupOutlined"},
		{Category: "monitoring", SubCategory: "alert", Key: "alert.silence_duration", Value: "30", Type: "number", Description: "默认静默时长(分钟)", SortOrder: 4, Icon: "BellOutlined"},

		// ========== Integration ==========
		// SSO设置
		{Category: "integration", SubCategory: "sso", Key: "sso.enabled", Value: "false", Type: "boolean", Description: "启用SSO登录", SortOrder: 1, Icon: "LoginOutlined"},
		{Category: "integration", SubCategory: "sso", Key: "sso.provider", Value: "", Type: "string", Options: `["ldap","oauth2","saml"]`, Description: "SSO提供商", SortOrder: 2, Icon: "SafetyOutlined"},
		{Category: "integration", SubCategory: "sso", Key: "sso.ldap_server", Value: "", Type: "string", Description: "LDAP服务器地址", SortOrder: 3, Icon: "GlobalOutlined"},
		{Category: "integration", SubCategory: "sso", Key: "sso.ldap_bind_dn", Value: "", Type: "string", Description: "LDAP绑定DN", SortOrder: 4, Icon: "UserOutlined"},
		{Category: "integration", SubCategory: "sso", Key: "sso.ldap_bind_password", Value: "", Type: "password", Sensitive: true, Description: "LDAP绑定密码", SortOrder: 5, Icon: "LockOutlined"},
		{Category: "integration", SubCategory: "sso", Key: "sso.ldap_base_dn", Value: "", Type: "string", Description: "LDAP基础DN", SortOrder: 6, Icon: "FolderOutlined"},

		// API设置
		{Category: "integration", SubCategory: "api", Key: "api.rate_limit", Value: "1000", Type: "number", Description: "API速率限制(请求/小时)", SortOrder: 1, Icon: "ApiOutlined"},
		{Category: "integration", SubCategory: "api", Key: "api.enable_cors", Value: "true", Type: "boolean", Description: "启用CORS", SortOrder: 2, Icon: "SwapOutlined"},
		{Category: "integration", SubCategory: "api", Key: "api.allowed_origins", Value: "*", Type: "string", Description: "允许的跨域来源", SortOrder: 3, Icon: "GlobalOutlined"},
	}

	// 插入配置项
	for _, item := range items {
		query := `INSERT INTO config_items
			(category, sub_category, key_name, value, type, options, default_val, required, ` + "`sensitive`" + `, description, sort_order, icon, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
		_, err := r.db.Exec(query, item.Category, item.SubCategory, item.Key, item.Value, item.Type,
			item.Options, item.DefaultVal, item.Required, item.Sensitive, item.Description,
			item.SortOrder, item.Icon)
		if err != nil {
			r.logger.Error("Failed to insert config item", zap.String("key", item.Key), zap.Error(err))
		}
	}

	r.logger.Info("Config items initialized", zap.Int("count", len(items)))
	return nil
}

// GetConfigTree 获取配置树
func (r *ConfigRepository) GetConfigTree() ([]model.ConfigTree, error) {
	// 定义树形结构
	tree := []model.ConfigTree{
		{
			Key:         "platform",
			Label:       "平台设置",
			Icon:        "SettingOutlined",
			Description: "平台基础配置",
			Children: []model.ConfigTree{
				{Key: "basic", Label: "基础配置", Icon: "InfoCircleOutlined", Description: "应用名称、版本、时区等"},
				{Key: "system", Label: "系统配置", Icon: "ToolOutlined", Description: "会话、缓存、日志等"},
				{Key: "notification", Label: "通知配置", Icon: "BellOutlined", Description: "邮件、Webhook通知"},
			},
		},
		{
			Key:         "ai",
			Label:       "AI智能",
			Icon:        "RobotOutlined",
			Description: "AI模型和分析配置",
			Children: []model.ConfigTree{
				{Key: "models", Label: "模型配置", Icon: "ThunderboltOutlined", Description: "AI模型连接配置"},
				{Key: "strategy", Label: "分析策略", Icon: "LineChartOutlined", Description: "AI分析行为配置"},
				{Key: "chat", Label: "对话设置", Icon: "MessageOutlined", Description: "AI对话相关配置"},
			},
		},
		{
			Key:         "monitoring",
			Label:       "监控告警",
			Icon:        "DashboardOutlined",
			Description: "监控和告警配置",
			Children: []model.ConfigTree{
				{Key: "prometheus", Label: "Prometheus", Icon: "DatabaseOutlined", Description: "Prometheus连接配置"},
				{Key: "alert", Label: "告警配置", Icon: "WarningOutlined", Description: "告警行为和通知"},
			},
		},
		{
			Key:         "integration",
			Label:       "集成",
			Icon:        "ApiOutlined",
			Description: "第三方系统集成",
			Children: []model.ConfigTree{
				{Key: "sso", Label: "单点登录", Icon: "SafetyOutlined", Description: "LDAP/OAuth2/SAML"},
				{Key: "api", Label: "API设置", Icon: "KeyOutlined", Description: "API密钥和限流"},
				{Key: "webhook", Label: "Webhook", Icon: "LinkOutlined", Description: "外部Webhook配置"},
			},
		},
	}

	return tree, nil
}

// GetConfigItemsByCategory 获取分类下的配置项
func (r *ConfigRepository) GetConfigItemsByCategory(category, subCategory string) ([]model.ConfigItem, error) {
	query := "SELECT id, category, sub_category, key_name, value, type, options, default_val, " +
		"required, `sensitive`, description, sort_order, icon, created_at, updated_at " +
		"FROM config_items WHERE category = ?"
	args := []interface{}{category}

	if subCategory != "" {
		query += " AND sub_category = ?"
		args = append(args, subCategory)
	}

	query += " ORDER BY sort_order ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.ConfigItem
	for rows.Next() {
		var item model.ConfigItem
		var options sql.NullString
		err := rows.Scan(&item.ID, &item.Category, &item.SubCategory, &item.Key, &item.Value,
			&item.Type, &options, &item.DefaultVal, &item.Required, &item.Sensitive,
			&item.Description, &item.SortOrder, &item.Icon, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			r.logger.Error("Failed to scan config item", zap.Error(err))
			continue
		}
		if options.Valid {
			item.Options = options.String
		}
		items = append(items, item)
	}

	return items, nil
}

// GetConfigItem 获取单个配置项
func (r *ConfigRepository) GetConfigItem(key string) (*model.ConfigItem, error) {
	query := "SELECT id, category, sub_category, key_name, value, type, options, default_val, " +
		"required, `sensitive`, description, sort_order, icon, created_at, updated_at " +
		"FROM config_items WHERE key_name = ?"

	var item model.ConfigItem
	var options sql.NullString
	err := r.db.QueryRow(query, key).Scan(&item.ID, &item.Category, &item.SubCategory, &item.Key,
		&item.Value, &item.Type, &options, &item.DefaultVal, &item.Required, &item.Sensitive,
		&item.Description, &item.SortOrder, &item.Icon, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if options.Valid {
		item.Options = options.String
	}
	return &item, nil
}

// UpdateConfigValue 更新配置值
func (r *ConfigRepository) UpdateConfigValue(key, value, user string) error {
	query := "UPDATE config_items SET value = ?, updated_at = NOW() WHERE key_name = ?"
	result, err := r.db.Exec(query, value, key)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("config key not found: %s", key)
	}

	return nil
}

// GetConfigValue 获取配置值
func (r *ConfigRepository) GetConfigValue(key string) (string, error) {
	query := "SELECT value FROM config_items WHERE key_name = ?"
	var value string
	err := r.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetAllConfigValues 获取所有配置值(返回map)
func (r *ConfigRepository) GetAllConfigValues() (map[string]string, error) {
	query := "SELECT key_name, value FROM config_items"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err == nil {
			values[key] = value
		}
	}

	return values, nil
}

// GetConfigItemsByKeys 批量获取配置项
func (r *ConfigRepository) GetConfigItemsByKeys(keys []string) ([]model.ConfigItem, error) {
	if len(keys) == 0 {
		return []model.ConfigItem{}, nil
	}

	// 构建IN查询
	query := "SELECT id, category, sub_category, key_name, value, type, options, default_val, " +
		"required, `sensitive`, description, sort_order, icon, created_at, updated_at " +
		"FROM config_items WHERE key_name IN ("
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = key
	}
	query += ") ORDER BY sort_order ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.ConfigItem
	for rows.Next() {
		var item model.ConfigItem
		var options sql.NullString
		err := rows.Scan(&item.ID, &item.Category, &item.SubCategory, &item.Key, &item.Value,
			&item.Type, &options, &item.DefaultVal, &item.Required, &item.Sensitive,
			&item.Description, &item.SortOrder, &item.Icon, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			r.logger.Error("Failed to scan config item", zap.Error(err))
			continue
		}
		if options.Valid {
			item.Options = options.String
		}
		items = append(items, item)
	}

	return items, nil
}

// ResetConfigToDefault 重置配置为默认值
func (r *ConfigRepository) ResetConfigToDefault(key string) error {
	query := "UPDATE config_items SET value = default_val, updated_at = NOW() WHERE key_name = ?"
	_, err := r.db.Exec(query, key)
	return err
}

// ExportConfig 导出所有配置
func (r *ConfigRepository) ExportConfig() ([]model.ConfigItem, error) {
	return r.GetConfigItemsByCategory("", "")
}

// ImportConfig 导入配置
func (r *ConfigRepository) ImportConfig(items []model.ConfigItem, user string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		query := `UPDATE config_items SET value = ?, updated_at = NOW() WHERE key_name = ?`
		_, err := tx.Exec(query, item.Value, item.Key)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetAIConfig 获取AI配置(用于服务层)
func (r *ConfigRepository) GetAIConfig() (map[string]interface{}, error) {
	// 获取默认AI模型配置
	var aiConfig model.AIModelConfig
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description FROM ai_model_configs
		WHERE is_default = TRUE AND is_enabled = TRUE LIMIT 1`

	var baseURL sql.NullString
	var desc sql.NullString
	err := r.db.QueryRow(query).Scan(&aiConfig.ID, &aiConfig.Name, &aiConfig.Provider, &aiConfig.Model,
		&aiConfig.APIKey, &baseURL, &aiConfig.MaxTokens, &aiConfig.Temperature, &aiConfig.Timeout,
		&aiConfig.IsDefault, &aiConfig.IsEnabled, &desc)

	if err != nil {
		if err == sql.ErrNoRows {
			// 返回空配置
			return map[string]interface{}{
				"enabled": false,
				"error":   "no active AI model config found",
			}, nil
		}
		return nil, err
	}

	if baseURL.Valid {
		aiConfig.BaseURL = baseURL.String
	}
	if desc.Valid {
		aiConfig.Description = desc.String
	}

	return map[string]interface{}{
		"enabled":     aiConfig.IsEnabled,
		"provider":    aiConfig.Provider,
		"model":       aiConfig.Model,
		"api_key":     aiConfig.APIKey,
		"base_url":    aiConfig.BaseURL,
		"max_tokens":  aiConfig.MaxTokens,
		"temperature": aiConfig.Temperature,
		"timeout":     aiConfig.Timeout,
	}, nil
}

// SaveAIConfig 保存AI配置
func (r *ConfigRepository) SaveAIConfig(config *model.AIModelConfig) error {
	query := `INSERT INTO ai_model_configs
		(name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		 is_default, is_enabled, description, created_by, updated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
		name = VALUES(name), provider = VALUES(provider), model = VALUES(model),
		api_key = VALUES(api_key), base_url = VALUES(base_url), max_tokens = VALUES(max_tokens),
		temperature = VALUES(temperature), timeout = VALUES(timeout),
		is_default = VALUES(is_default), is_enabled = VALUES(is_enabled),
		description = VALUES(description), updated_by = VALUES(updated_by), updated_at = NOW()`

	_, err := r.db.Exec(query, config.Name, config.Provider, config.Model, config.APIKey,
		config.BaseURL, config.MaxTokens, config.Temperature, config.Timeout,
		config.IsDefault, config.IsEnabled, config.Description, config.CreatedBy, config.UpdatedBy)
	return err
}

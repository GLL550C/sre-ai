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
	return &ConfigRepository{db: db, logger: logger}
}

// GetByCategory 按分类获取配置
func (r *ConfigRepository) GetByCategory(category string) ([]model.SystemConfig, error) {
	query := `SELECT id, category, sub_category, config_key, config_value, value_type, is_sensitive, description
		FROM system_configs WHERE category = ? ORDER BY id`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.SystemConfig
	for rows.Next() {
		var c model.SystemConfig
		err := rows.Scan(&c.ID, &c.Category, &c.SubCategory, &c.Key, &c.Value, &c.ValueType, &c.IsSensitive, &c.Description)
		if err != nil {
			continue
		}
		configs = append(configs, c)
	}
	return configs, nil
}

// GetByKey 根据key获取配置
func (r *ConfigRepository) GetByKey(key string) (*model.SystemConfig, error) {
	query := `SELECT id, category, sub_category, config_key, config_value, value_type, is_sensitive, description
		FROM system_configs WHERE config_key = ?`

	var c model.SystemConfig
	err := r.db.QueryRow(query, key).Scan(&c.ID, &c.Category, &c.SubCategory, &c.Key, &c.Value, &c.ValueType, &c.IsSensitive, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetValue 获取配置值
func (r *ConfigRepository) GetValue(key string) (string, error) {
	query := `SELECT config_value FROM system_configs WHERE config_key = ?`
	var value string
	err := r.db.QueryRow(query, key).Scan(&value)
	return value, err
}

// Update 更新配置
func (r *ConfigRepository) Update(key, value, user string) error {
	query := `UPDATE system_configs SET config_value = ?, updated_by = ?, updated_at = NOW()
		WHERE config_key = ?`
	result, err := r.db.Exec(query, value, user, key)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("config not found: %s", key)
	}
	return nil
}

// GetAll 获取所有配置
func (r *ConfigRepository) GetAll() ([]model.SystemConfig, error) {
	query := `SELECT id, category, sub_category, config_key, config_value, value_type, is_sensitive, description
		FROM system_configs ORDER BY category, id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.SystemConfig
	for rows.Next() {
		var c model.SystemConfig
		err := rows.Scan(&c.ID, &c.Category, &c.SubCategory, &c.Key, &c.Value, &c.ValueType, &c.IsSensitive, &c.Description)
		if err != nil {
			continue
		}
		configs = append(configs, c)
	}
	return configs, nil
}

// GetByCategoryAndSubCategory 按分类和子分类获取配置
func (r *ConfigRepository) GetByCategoryAndSubCategory(category, subCategory string) ([]model.SystemConfig, error) {
	query := `SELECT id, category, sub_category, config_key, config_value, value_type, is_sensitive, description
		FROM system_configs WHERE category = ? AND (sub_category = ? OR sub_category IS NULL) ORDER BY id`

	rows, err := r.db.Query(query, category, subCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.SystemConfig
	for rows.Next() {
		var c model.SystemConfig
		err := rows.Scan(&c.ID, &c.Category, &c.SubCategory, &c.Key, &c.Value, &c.ValueType, &c.IsSensitive, &c.Description)
		if err != nil {
			continue
		}
		configs = append(configs, c)
	}
	return configs, nil
}

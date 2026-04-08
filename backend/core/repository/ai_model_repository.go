package repository

import (
	"core/model"
	"database/sql"

	"go.uber.org/zap"
)

// AIModelRepository AI模型数据访问
type AIModelRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAIModelRepository 创建仓库
func NewAIModelRepository(db *sql.DB, logger *zap.Logger) *AIModelRepository {
	return &AIModelRepository{db: db, logger: logger}
}

// GetAll 获取所有模型配置
func (r *AIModelRepository) GetAll() ([]model.AIModel, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_models ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []model.AIModel
	for rows.Next() {
		var m model.AIModel
		var updatedBy sql.NullString
		err := rows.Scan(&m.ID, &m.Name, &m.Provider, &m.Model, &m.APIKey, &m.BaseURL, &m.MaxTokens,
			&m.Temperature, &m.Timeout, &m.IsDefault, &m.IsEnabled, &m.Description,
			&m.CreatedBy, &updatedBy, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			r.logger.Error("扫描AI模型失败", zap.Error(err))
			continue
		}
		if updatedBy.Valid {
			m.UpdatedBy = updatedBy.String
		}
		models = append(models, m)
	}
	return models, nil
}

// GetByID 根据ID获取
func (r *AIModelRepository) GetByID(id int64) (*model.AIModel, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_models WHERE id = ?`

	var m model.AIModel
	var updatedBy sql.NullString
	err := r.db.QueryRow(query, id).Scan(&m.ID, &m.Name, &m.Provider, &m.Model, &m.APIKey, &m.BaseURL, &m.MaxTokens,
		&m.Temperature, &m.Timeout, &m.IsDefault, &m.IsEnabled, &m.Description,
		&m.CreatedBy, &updatedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if updatedBy.Valid {
		m.UpdatedBy = updatedBy.String
	}
	return &m, nil
}

// GetDefault 获取默认配置
func (r *AIModelRepository) GetDefault() (*model.AIModel, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_models WHERE is_default = TRUE AND is_enabled = TRUE LIMIT 1`

	var m model.AIModel
	var updatedBy sql.NullString
	err := r.db.QueryRow(query).Scan(&m.ID, &m.Name, &m.Provider, &m.Model, &m.APIKey, &m.BaseURL, &m.MaxTokens,
		&m.Temperature, &m.Timeout, &m.IsDefault, &m.IsEnabled, &m.Description,
		&m.CreatedBy, &updatedBy, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if updatedBy.Valid {
		m.UpdatedBy = updatedBy.String
	}
	return &m, nil
}

// Create 创建配置
func (r *AIModelRepository) Create(m *model.AIModel) error {
	query := `INSERT INTO ai_models (name, provider, model, api_key, base_url, max_tokens,
		temperature, timeout, is_default, is_enabled, description, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if m.IsDefault {
		r.clearDefault()
	}

	result, err := r.db.Exec(query, m.Name, m.Provider, m.Model, m.APIKey, m.BaseURL, m.MaxTokens,
		m.Temperature, m.Timeout, m.IsDefault, m.IsEnabled, m.Description, m.CreatedBy, m.UpdatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	m.ID = id
	return nil
}

// Update 更新配置
func (r *AIModelRepository) Update(m *model.AIModel) error {
	query := `UPDATE ai_models SET name = ?, provider = ?, model = ?, api_key = ?, base_url = ?,
		max_tokens = ?, temperature = ?, timeout = ?, is_default = ?, is_enabled = ?, description = ?,
		updated_by = ? WHERE id = ?`

	if m.IsDefault {
		r.clearDefault()
	}

	_, err := r.db.Exec(query, m.Name, m.Provider, m.Model, m.APIKey, m.BaseURL, m.MaxTokens,
		m.Temperature, m.Timeout, m.IsDefault, m.IsEnabled, m.Description, m.UpdatedBy, m.ID)
	return err
}

// Delete 删除配置
func (r *AIModelRepository) Delete(id int64) error {
	query := `DELETE FROM ai_models WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// clearDefault 清除默认标记
func (r *AIModelRepository) clearDefault() error {
	_, err := r.db.Exec(`UPDATE ai_models SET is_default = FALSE WHERE is_default = TRUE`)
	return err
}

// SetDefault 设置默认
func (r *AIModelRepository) SetDefault(id int64) error {
	if err := r.clearDefault(); err != nil {
		return err
	}
	_, err := r.db.Exec(`UPDATE ai_models SET is_default = TRUE WHERE id = ?`, id)
	return err
}

package repository

import (
	"core/model"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// AIModelRepository handles AI model config data access
type AIModelRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAIModelRepository creates a new AI model repository
func NewAIModelRepository(db *sql.DB, logger *zap.Logger) *AIModelRepository {
	return &AIModelRepository{
		db:     db,
		logger: logger,
	}
}

// GetAllConfigs retrieves all AI model configs
func (r *AIModelRepository) GetAllConfigs() ([]model.AIModelConfig, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_model_configs ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.AIModelConfig
	for rows.Next() {
		var c model.AIModelConfig
		var updatedBy sql.NullString
		err := rows.Scan(
			&c.ID, &c.Name, &c.Provider, &c.Model, &c.APIKey, &c.BaseURL, &c.MaxTokens,
			&c.Temperature, &c.Timeout, &c.IsDefault, &c.IsEnabled, &c.Description,
			&c.CreatedBy, &updatedBy, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan AI config", zap.Error(err))
			continue
		}
		if updatedBy.Valid {
			c.UpdatedBy = updatedBy.String
		}
		configs = append(configs, c)
	}

	return configs, nil
}

// GetConfigByID retrieves an AI model config by ID
func (r *AIModelRepository) GetConfigByID(id int64) (*model.AIModelConfig, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_model_configs WHERE id = ?`

	var c model.AIModelConfig
	var updatedBy sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Provider, &c.Model, &c.APIKey, &c.BaseURL, &c.MaxTokens,
		&c.Temperature, &c.Timeout, &c.IsDefault, &c.IsEnabled, &c.Description,
		&c.CreatedBy, &updatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if updatedBy.Valid {
		c.UpdatedBy = updatedBy.String
	}
	return &c, nil
}

// GetDefaultConfig retrieves the default AI model config
func (r *AIModelRepository) GetDefaultConfig() (*model.AIModelConfig, error) {
	query := `SELECT id, name, provider, model, api_key, base_url, max_tokens, temperature, timeout,
		is_default, is_enabled, description, created_by, updated_by, created_at, updated_at
		FROM ai_model_configs WHERE is_default = TRUE AND is_enabled = TRUE LIMIT 1`

	var c model.AIModelConfig
	var updatedBy sql.NullString
	err := r.db.QueryRow(query).Scan(
		&c.ID, &c.Name, &c.Provider, &c.Model, &c.APIKey, &c.BaseURL, &c.MaxTokens,
		&c.Temperature, &c.Timeout, &c.IsDefault, &c.IsEnabled, &c.Description,
		&c.CreatedBy, &updatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if updatedBy.Valid {
		c.UpdatedBy = updatedBy.String
	}
	return &c, nil
}

// CreateConfig creates a new AI model config
func (r *AIModelRepository) CreateConfig(c *model.AIModelConfig) error {
	query := `INSERT INTO ai_model_configs (name, provider, model, api_key, base_url, max_tokens,
		temperature, timeout, is_default, is_enabled, description, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// If this is set as default, clear other defaults
	if c.IsDefault {
		r.clearDefaultConfig()
	}

	result, err := r.db.Exec(query,
		c.Name, c.Provider, c.Model, c.APIKey, c.BaseURL, c.MaxTokens,
		c.Temperature, c.Timeout, c.IsDefault, c.IsEnabled, c.Description,
		c.CreatedBy, time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	c.ID = id
	return nil
}

// UpdateConfig updates an AI model config
func (r *AIModelRepository) UpdateConfig(c *model.AIModelConfig) error {
	query := `UPDATE ai_model_configs SET name = ?, provider = ?, model = ?, api_key = ?, base_url = ?,
		max_tokens = ?, temperature = ?, timeout = ?, is_default = ?, is_enabled = ?, description = ?,
		updated_by = ?, updated_at = ? WHERE id = ?`

	// If this is set as default, clear other defaults
	if c.IsDefault {
		r.clearDefaultConfig()
	}

	_, err := r.db.Exec(query,
		c.Name, c.Provider, c.Model, c.APIKey, c.BaseURL, c.MaxTokens,
		c.Temperature, c.Timeout, c.IsDefault, c.IsEnabled, c.Description,
		c.UpdatedBy, time.Now(), c.ID,
	)
	return err
}

// DeleteConfig deletes an AI model config
func (r *AIModelRepository) DeleteConfig(id int64) error {
	query := "DELETE FROM ai_model_configs WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// TestConfig tests an AI model config by ID
func (r *AIModelRepository) TestConfig(id int64) (bool, string, error) {
	config, err := r.GetConfigByID(id)
	if err != nil {
		return false, "", err
	}

	if !config.IsEnabled {
		return false, "Config is disabled", nil
	}

	// Return success - actual test will be done by service layer
	return true, fmt.Sprintf("Config '%s' is valid", config.Name), nil
}

// clearDefaultConfig clears the default flag from all configs
func (r *AIModelRepository) clearDefaultConfig() error {
	_, err := r.db.Exec("UPDATE ai_model_configs SET is_default = FALSE WHERE is_default = TRUE")
	return err
}

// SetDefaultConfig sets a config as default
func (r *AIModelRepository) SetDefaultConfig(id int64) error {
	if err := r.clearDefaultConfig(); err != nil {
		return err
	}
	_, err := r.db.Exec("UPDATE ai_model_configs SET is_default = TRUE WHERE id = ?", id)
	return err
}

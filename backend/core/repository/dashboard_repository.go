package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

// DashboardRepository handles dashboard configuration data access
type DashboardRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewDashboardRepository creates a new dashboard repository
func NewDashboardRepository(db *sql.DB, logger *zap.Logger) *DashboardRepository {
	return &DashboardRepository{
		db:     db,
		logger: logger,
	}
}

// GetDashboardConfigs retrieves dashboard configurations
func (r *DashboardRepository) GetDashboardConfigs(tenantID *int64) ([]model.DashboardConfig, error) {
	query := `SELECT id, name, tenant_id, layout, widgets, is_default, created_at, updated_at
		FROM dashboard_configs WHERE 1=1`
	args := []interface{}{}

	if tenantID != nil {
		query += " AND tenant_id = ?"
		args = append(args, *tenantID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []model.DashboardConfig
	for rows.Next() {
		var config model.DashboardConfig
		var layout, widgets []byte
		err := rows.Scan(
			&config.ID, &config.Name, &config.TenantID, &layout, &widgets,
			&config.IsDefault, &config.CreatedAt, &config.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan dashboard config", zap.Error(err))
			continue
		}
		if layout != nil {
			config.Layout = layout
		}
		if widgets != nil {
			config.Widgets = widgets
		}
		configs = append(configs, config)
	}

	return configs, nil
}

// GetDefaultDashboard retrieves the default dashboard
func (r *DashboardRepository) GetDefaultDashboard(tenantID *int64) (*model.DashboardConfig, error) {
	query := `SELECT id, name, tenant_id, layout, widgets, is_default, created_at, updated_at
		FROM dashboard_configs WHERE is_default = true`
	args := []interface{}{}

	if tenantID != nil {
		query += " AND (tenant_id = ? OR tenant_id IS NULL)"
		args = append(args, *tenantID)
	}

	query += " LIMIT 1"

	var config model.DashboardConfig
	var layout, widgets []byte
	err := r.db.QueryRow(query, args...).Scan(
		&config.ID, &config.Name, &config.TenantID, &layout, &widgets,
		&config.IsDefault, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if layout != nil {
		config.Layout = layout
	}
	if widgets != nil {
		config.Widgets = widgets
	}

	return &config, nil
}

// CreateDashboardConfig creates a new dashboard configuration
func (r *DashboardRepository) CreateDashboardConfig(config *model.DashboardConfig) error {
	query := `INSERT INTO dashboard_configs (name, tenant_id, layout, widgets, is_default)
		VALUES (?, ?, ?, ?, ?)`

	layoutJSON, _ := json.Marshal(config.Layout)
	widgetsJSON, _ := json.Marshal(config.Widgets)

	result, err := r.db.Exec(query,
		config.Name, config.TenantID, layoutJSON, widgetsJSON, config.IsDefault,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	config.ID = id
	return nil
}

// UpdateDashboardConfig updates a dashboard configuration
func (r *DashboardRepository) UpdateDashboardConfig(config *model.DashboardConfig) error {
	query := `UPDATE dashboard_configs SET name = ?, tenant_id = ?, layout = ?, widgets = ?, is_default = ? WHERE id = ?`

	layoutJSON, _ := json.Marshal(config.Layout)
	widgetsJSON, _ := json.Marshal(config.Widgets)

	_, err := r.db.Exec(query,
		config.Name, config.TenantID, layoutJSON, widgetsJSON, config.IsDefault, config.ID,
	)
	return err
}

// DeleteDashboardConfig deletes a dashboard configuration
func (r *DashboardRepository) DeleteDashboardConfig(id int64) error {
	query := "DELETE FROM dashboard_configs WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

package repository

import (
	"database/sql"
	"encoding/json"
	"tenant/model"

	"go.uber.org/zap"
)

// TenantRepository handles tenant data access
type TenantRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *sql.DB, logger *zap.Logger) *TenantRepository {
	return &TenantRepository{
		db:     db,
		logger: logger,
	}
}

// GetTenants retrieves all tenants
func (r *TenantRepository) GetTenants(status int) ([]model.Tenant, error) {
	query := `SELECT id, name, code, description, config, clusters, status, created_at, updated_at
		FROM tenants WHERE status = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []model.Tenant
	for rows.Next() {
		var tenant model.Tenant
		var config, clusters []byte
		err := rows.Scan(
			&tenant.ID, &tenant.Name, &tenant.Code, &tenant.Description,
			&config, &clusters, &tenant.Status, &tenant.CreatedAt, &tenant.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan tenant", zap.Error(err))
			continue
		}
		if config != nil {
			tenant.Config = config
		}
		if clusters != nil {
			tenant.Clusters = clusters
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// GetTenantByID retrieves a tenant by ID
func (r *TenantRepository) GetTenantByID(id int64) (*model.Tenant, error) {
	query := `SELECT id, name, code, description, config, clusters, status, created_at, updated_at
		FROM tenants WHERE id = ?`

	var tenant model.Tenant
	var config, clusters []byte
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID, &tenant.Name, &tenant.Code, &tenant.Description,
		&config, &clusters, &tenant.Status, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if config != nil {
		tenant.Config = config
	}
	if clusters != nil {
		tenant.Clusters = clusters
	}

	return &tenant, nil
}

// GetTenantByCode retrieves a tenant by code
func (r *TenantRepository) GetTenantByCode(code string) (*model.Tenant, error) {
	query := `SELECT id, name, code, description, config, clusters, status, created_at, updated_at
		FROM tenants WHERE code = ?`

	var tenant model.Tenant
	var config, clusters []byte
	err := r.db.QueryRow(query, code).Scan(
		&tenant.ID, &tenant.Name, &tenant.Code, &tenant.Description,
		&config, &clusters, &tenant.Status, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if config != nil {
		tenant.Config = config
	}
	if clusters != nil {
		tenant.Clusters = clusters
	}

	return &tenant, nil
}

// CreateTenant creates a new tenant
func (r *TenantRepository) CreateTenant(tenant *model.Tenant) error {
	query := `INSERT INTO tenants (name, code, description, config, clusters, status)
		VALUES (?, ?, ?, ?, ?, ?)`

	configJSON, _ := json.Marshal(tenant.Config)
	clustersJSON, _ := json.Marshal(tenant.Clusters)

	result, err := r.db.Exec(query,
		tenant.Name, tenant.Code, tenant.Description, configJSON, clustersJSON, tenant.Status,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	tenant.ID = id
	return nil
}

// UpdateTenant updates a tenant
func (r *TenantRepository) UpdateTenant(tenant *model.Tenant) error {
	query := `UPDATE tenants SET name = ?, code = ?, description = ?, config = ?, clusters = ?, status = ? WHERE id = ?`

	configJSON, _ := json.Marshal(tenant.Config)
	clustersJSON, _ := json.Marshal(tenant.Clusters)

	_, err := r.db.Exec(query,
		tenant.Name, tenant.Code, tenant.Description, configJSON, clustersJSON, tenant.Status, tenant.ID,
	)
	return err
}

// DeleteTenant soft deletes a tenant
func (r *TenantRepository) DeleteTenant(id int64) error {
	query := "UPDATE tenants SET status = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

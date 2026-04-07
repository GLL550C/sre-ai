package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

// ClusterRepository handles Prometheus cluster data access
type ClusterRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewClusterRepository creates a new cluster repository
func NewClusterRepository(db *sql.DB, logger *zap.Logger) *ClusterRepository {
	return &ClusterRepository{
		db:     db,
		logger: logger,
	}
}

// GetClusters retrieves all Prometheus clusters
func (r *ClusterRepository) GetClusters(status int) ([]model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters WHERE status = ? ORDER BY is_default DESC, created_at DESC`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []model.PrometheusCluster
	for rows.Next() {
		var cluster model.PrometheusCluster
		var configJSON []byte
		err := rows.Scan(
			&cluster.ID, &cluster.Name, &cluster.URL, &cluster.Status, &cluster.IsDefault,
			&configJSON, &cluster.CreatedAt, &cluster.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan cluster", zap.Error(err))
			continue
		}
		if configJSON != nil {
			cluster.ConfigJSON = configJSON
		}
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// GetClusterByID retrieves a cluster by ID
func (r *ClusterRepository) GetClusterByID(id int64) (*model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters WHERE id = ?`

	var cluster model.PrometheusCluster
	var configJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&cluster.ID, &cluster.Name, &cluster.URL, &cluster.Status, &cluster.IsDefault,
		&configJSON, &cluster.CreatedAt, &cluster.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if configJSON != nil {
		cluster.ConfigJSON = configJSON
	}

	return &cluster, nil
}

// GetDefaultCluster retrieves the default cluster
func (r *ClusterRepository) GetDefaultCluster() (*model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters WHERE is_default = TRUE AND status = 1 LIMIT 1`

	var cluster model.PrometheusCluster
	var configJSON []byte
	err := r.db.QueryRow(query).Scan(
		&cluster.ID, &cluster.Name, &cluster.URL, &cluster.Status, &cluster.IsDefault,
		&configJSON, &cluster.CreatedAt, &cluster.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if configJSON != nil {
		cluster.ConfigJSON = configJSON
	}

	return &cluster, nil
}

// GetActiveClusters retrieves all active clusters
func (r *ClusterRepository) GetActiveClusters() ([]model.PrometheusCluster, error) {
	return r.GetClusters(1)
}

// CreateCluster creates a new Prometheus cluster
func (r *ClusterRepository) CreateCluster(cluster *model.PrometheusCluster) error {
	query := `INSERT INTO prometheus_clusters (name, url, status, is_default, config_json)
		VALUES (?, ?, ?, ?, ?)`

	configJSON, _ := json.Marshal(cluster.ConfigJSON)

	result, err := r.db.Exec(query, cluster.Name, cluster.URL, cluster.Status, cluster.IsDefault, configJSON)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	cluster.ID = id
	return nil
}

// UpdateCluster updates a Prometheus cluster
func (r *ClusterRepository) UpdateCluster(cluster *model.PrometheusCluster) error {
	query := `UPDATE prometheus_clusters SET name = ?, url = ?, status = ?, is_default = ?, config_json = ? WHERE id = ?`

	configJSON, _ := json.Marshal(cluster.ConfigJSON)

	_, err := r.db.Exec(query, cluster.Name, cluster.URL, cluster.Status, cluster.IsDefault, configJSON, cluster.ID)
	return err
}

// ClearDefaultCluster clears the default flag from all clusters
func (r *ClusterRepository) ClearDefaultCluster() error {
	query := `UPDATE prometheus_clusters SET is_default = FALSE WHERE is_default = TRUE`
	_, err := r.db.Exec(query)
	return err
}

// SetDefaultCluster sets a cluster as default
func (r *ClusterRepository) SetDefaultCluster(id int64) error {
	query := `UPDATE prometheus_clusters SET is_default = TRUE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// DeleteCluster soft deletes a cluster
func (r *ClusterRepository) DeleteCluster(id int64) error {
	query := "UPDATE prometheus_clusters SET status = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

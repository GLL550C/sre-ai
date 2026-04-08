package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

// ClusterRepository Prometheus集群数据访问
type ClusterRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewClusterRepository 创建仓库
func NewClusterRepository(db *sql.DB, logger *zap.Logger) *ClusterRepository {
	return &ClusterRepository{db: db, logger: logger}
}

// GetAll 获取所有集群
func (r *ClusterRepository) GetAll() ([]model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters ORDER BY is_default DESC, id ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []model.PrometheusCluster
	for rows.Next() {
		var c model.PrometheusCluster
		var cfg []byte
		err := rows.Scan(&c.ID, &c.Name, &c.URL, &c.Status, &c.IsDefault, &cfg, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			r.logger.Error("扫描集群失败", zap.Error(err))
			continue
		}
		if cfg != nil {
			c.ConfigJSON = cfg
		}
		clusters = append(clusters, c)
	}
	return clusters, nil
}

// GetByID 根据ID获取
func (r *ClusterRepository) GetByID(id int64) (*model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters WHERE id = ?`

	var c model.PrometheusCluster
	var cfg []byte
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.URL, &c.Status, &c.IsDefault, &cfg, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		c.ConfigJSON = cfg
	}
	return &c, nil
}

// GetClusterByID 根据ID获取(兼容旧代码)
func (r *ClusterRepository) GetClusterByID(id int64) (*model.PrometheusCluster, error) {
	return r.GetByID(id)
}

// GetDefault 获取默认集群
func (r *ClusterRepository) GetDefault() (*model.PrometheusCluster, error) {
	query := `SELECT id, name, url, status, is_default, config_json, created_at, updated_at
		FROM prometheus_clusters WHERE is_default = TRUE AND status = 1 LIMIT 1`

	var c model.PrometheusCluster
	var cfg []byte
	err := r.db.QueryRow(query).Scan(&c.ID, &c.Name, &c.URL, &c.Status, &c.IsDefault, &cfg, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		c.ConfigJSON = cfg
	}
	return &c, nil
}

// Create 创建集群
func (r *ClusterRepository) Create(c *model.PrometheusCluster) error {
	query := `INSERT INTO prometheus_clusters (name, url, status, is_default, config_json, created_by)
		VALUES (?, ?, ?, ?, ?, ?)`

	cfg, _ := json.Marshal(c.ConfigJSON)
	result, err := r.db.Exec(query, c.Name, c.URL, c.Status, c.IsDefault, cfg, c.CreatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	c.ID = id
	return nil
}

// Update 更新集群
func (r *ClusterRepository) Update(c *model.PrometheusCluster) error {
	query := `UPDATE prometheus_clusters
		SET name = ?, url = ?, status = ?, is_default = ?, config_json = ?, updated_by = ?
		WHERE id = ?`

	cfg, _ := json.Marshal(c.ConfigJSON)
	_, err := r.db.Exec(query, c.Name, c.URL, c.Status, c.IsDefault, cfg, c.UpdatedBy, c.ID)
	return err
}

// Delete 删除集群(软删除)
func (r *ClusterRepository) Delete(id int64) error {
	query := `UPDATE prometheus_clusters SET status = 0 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// ClearDefault 清除默认标记
func (r *ClusterRepository) ClearDefault() error {
	query := `UPDATE prometheus_clusters SET is_default = FALSE WHERE is_default = TRUE`
	_, err := r.db.Exec(query)
	return err
}

// SetDefault 设置默认集群
func (r *ClusterRepository) SetDefault(id int64) error {
	query := `UPDATE prometheus_clusters SET is_default = TRUE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

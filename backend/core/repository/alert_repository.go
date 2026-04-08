package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// AlertRepository 告警数据访问
type AlertRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAlertRepository 创建仓库
func NewAlertRepository(db *sql.DB, logger *zap.Logger) *AlertRepository {
	return &AlertRepository{db: db, logger: logger}
}

// GetAlerts 获取告警列表
func (r *AlertRepository) GetAlerts(status, severity string, limit, offset int) ([]model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, created_at, updated_at
		FROM alerts WHERE 1=1`
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if severity != "" {
		query += " AND severity = ?"
		args = append(args, severity)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.Alert
	for rows.Next() {
		var a model.Alert
		var labels []byte
		err := rows.Scan(&a.ID, &a.RuleID, &a.Fingerprint, &a.Status, &a.Severity,
			&a.Summary, &a.Description, &labels, &a.StartsAt, &a.EndsAt,
			&a.AcknowledgedBy, &a.AcknowledgedAt, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			r.logger.Error("扫描告警失败", zap.Error(err))
			continue
		}
		if labels != nil {
			a.Labels = labels
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

// GetByID 根据ID获取
func (r *AlertRepository) GetByID(id int64) (*model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, created_at, updated_at
		FROM alerts WHERE id = ?`

	var a model.Alert
	var labels []byte
	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.RuleID, &a.Fingerprint, &a.Status, &a.Severity,
		&a.Summary, &a.Description, &labels, &a.StartsAt, &a.EndsAt,
		&a.AcknowledgedBy, &a.AcknowledgedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if labels != nil {
		a.Labels = labels
	}
	return &a, nil
}

// GetByFingerprint 根据指纹获取
func (r *AlertRepository) GetByFingerprint(fingerprint string) (*model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, created_at, updated_at
		FROM alerts WHERE fingerprint = ?`

	var a model.Alert
	var labels []byte
	err := r.db.QueryRow(query, fingerprint).Scan(&a.ID, &a.RuleID, &a.Fingerprint, &a.Status, &a.Severity,
		&a.Summary, &a.Description, &labels, &a.StartsAt, &a.EndsAt,
		&a.AcknowledgedBy, &a.AcknowledgedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if labels != nil {
		a.Labels = labels
	}
	return &a, nil
}

// Create 创建告警
func (r *AlertRepository) Create(a *model.Alert) error {
	query := `INSERT INTO alerts (rule_id, fingerprint, status, severity, summary, description, labels, starts_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	labelsJSON, _ := json.Marshal(a.Labels)
	result, err := r.db.Exec(query, a.RuleID, a.Fingerprint, a.Status, a.Severity,
		a.Summary, a.Description, labelsJSON, a.StartsAt)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	a.ID = id
	return nil
}

// UpdateStatus 更新状态
func (r *AlertRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE alerts SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}

// Acknowledge 确认告警
func (r *AlertRepository) Acknowledge(id int64, user string) error {
	query := `UPDATE alerts SET status = 'acknowledged', acknowledged_by = ?, acknowledged_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, user, id)
	return err
}

// Resolve 解决告警
func (r *AlertRepository) Resolve(fingerprint string) error {
	query := `UPDATE alerts SET status = 'resolved', ends_at = NOW() WHERE fingerprint = ? AND status = 'firing'`
	_, err := r.db.Exec(query, fingerprint)
	return err
}

// GetCount 获取数量
func (r *AlertRepository) GetCount(status, severity string) (int, error) {
	query := `SELECT COUNT(*) FROM alerts WHERE 1=1`
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if severity != "" {
		query += " AND severity = ?"
		args = append(args, severity)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

// GetRecent 获取最近告警
func (r *AlertRepository) GetRecent(duration time.Duration) ([]model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, created_at, updated_at
		FROM alerts WHERE created_at > DATE_SUB(NOW(), INTERVAL ? HOUR) ORDER BY created_at DESC`

	hours := int(duration.Hours())
	if hours < 1 {
		hours = 24
	}

	rows, err := r.db.Query(query, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.Alert
	for rows.Next() {
		var a model.Alert
		var labels []byte
		err := rows.Scan(&a.ID, &a.RuleID, &a.Fingerprint, &a.Status, &a.Severity,
			&a.Summary, &a.Description, &labels, &a.StartsAt, &a.EndsAt,
			&a.AcknowledgedBy, &a.AcknowledgedAt, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			continue
		}
		if labels != nil {
			a.Labels = labels
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

// GetRecentAlertsByCluster 获取集群最近告警(兼容旧代码)
func (r *AlertRepository) GetRecentAlertsByCluster(clusterID int64, duration time.Duration) ([]model.Alert, error) {
	// 新表结构没有cluster_id字段，返回所有最近告警
	return r.GetRecent(duration)
}

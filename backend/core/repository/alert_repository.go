package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// AlertRepository handles alert data access
type AlertRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(db *sql.DB, logger *zap.Logger) *AlertRepository {
	return &AlertRepository{
		db:     db,
		logger: logger,
	}
}

// GetAlerts retrieves alerts with optional filters
func (r *AlertRepository) GetAlerts(status, severity string, limit, offset int) ([]model.Alert, error) {
	query := "SELECT id, rule_id, fingerprint, status, severity, summary, description, labels, starts_at, ends_at, acknowledged_by, acknowledged_at, cluster_id, created_at, updated_at FROM alerts WHERE 1=1"
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
		var alert model.Alert
		var labels []byte
		err := rows.Scan(
			&alert.ID, &alert.RuleID, &alert.Fingerprint, &alert.Status, &alert.Severity,
			&alert.Summary, &alert.Description, &labels, &alert.StartsAt, &alert.EndsAt,
			&alert.AcknowledgedBy, &alert.AcknowledgedAt, &alert.ClusterID, &alert.CreatedAt, &alert.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan alert", zap.Error(err))
			continue
		}
		if labels != nil {
			alert.Labels = labels
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetAlertByID retrieves an alert by ID
func (r *AlertRepository) GetAlertByID(id int64) (*model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, cluster_id, created_at, updated_at
		FROM alerts WHERE id = ?`

	var alert model.Alert
	var labels []byte
	err := r.db.QueryRow(query, id).Scan(
		&alert.ID, &alert.RuleID, &alert.Fingerprint, &alert.Status, &alert.Severity,
		&alert.Summary, &alert.Description, &labels, &alert.StartsAt, &alert.EndsAt,
		&alert.AcknowledgedBy, &alert.AcknowledgedAt, &alert.ClusterID, &alert.CreatedAt, &alert.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if labels != nil {
		alert.Labels = labels
	}

	return &alert, nil
}

// GetAlertByFingerprint retrieves an alert by fingerprint
func (r *AlertRepository) GetAlertByFingerprint(fingerprint string) (*model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, cluster_id, created_at, updated_at
		FROM alerts WHERE fingerprint = ?`

	var alert model.Alert
	var labels []byte
	err := r.db.QueryRow(query, fingerprint).Scan(
		&alert.ID, &alert.RuleID, &alert.Fingerprint, &alert.Status, &alert.Severity,
		&alert.Summary, &alert.Description, &labels, &alert.StartsAt, &alert.EndsAt,
		&alert.AcknowledgedBy, &alert.AcknowledgedAt, &alert.ClusterID, &alert.CreatedAt, &alert.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if labels != nil {
		alert.Labels = labels
	}

	return &alert, nil
}

// GetRecentAlertsByCluster retrieves recent alerts for a cluster
func (r *AlertRepository) GetRecentAlertsByCluster(clusterID int64, duration time.Duration) ([]model.Alert, error) {
	query := `SELECT id, rule_id, fingerprint, status, severity, summary, description, labels,
		starts_at, ends_at, acknowledged_by, acknowledged_at, cluster_id, created_at, updated_at
		FROM alerts WHERE cluster_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL ? HOUR)
		ORDER BY created_at DESC`

	hours := int(duration.Hours())
	if hours < 1 {
		hours = 24
	}

	rows, err := r.db.Query(query, clusterID, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.Alert
	for rows.Next() {
		var alert model.Alert
		var labels []byte
		err := rows.Scan(
			&alert.ID, &alert.RuleID, &alert.Fingerprint, &alert.Status, &alert.Severity,
			&alert.Summary, &alert.Description, &labels, &alert.StartsAt, &alert.EndsAt,
			&alert.AcknowledgedBy, &alert.AcknowledgedAt, &alert.ClusterID, &alert.CreatedAt, &alert.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan alert", zap.Error(err))
			continue
		}
		if labels != nil {
			alert.Labels = labels
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// CreateAlert creates a new alert
func (r *AlertRepository) CreateAlert(alert *model.Alert) error {
	query := `INSERT INTO alerts (rule_id, fingerprint, status, severity, summary, description, labels, starts_at, cluster_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	labelsJSON, _ := json.Marshal(alert.Labels)
	result, err := r.db.Exec(query,
		alert.RuleID, alert.Fingerprint, alert.Status, alert.Severity,
		alert.Summary, alert.Description, labelsJSON, alert.StartsAt, alert.ClusterID,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	alert.ID = id
	return nil
}

// UpdateAlertStatus updates alert status
func (r *AlertRepository) UpdateAlertStatus(id int64, status string) error {
	query := "UPDATE alerts SET status = ? WHERE id = ?"
	_, err := r.db.Exec(query, status, id)
	return err
}

// AcknowledgeAlert acknowledges an alert
func (r *AlertRepository) AcknowledgeAlert(id int64, user string) error {
	query := "UPDATE alerts SET status = 'acknowledged', acknowledged_by = ?, acknowledged_at = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, user, id)
	return err
}

// ResolveAlert resolves an alert by fingerprint
func (r *AlertRepository) ResolveAlert(fingerprint string) error {
	query := "UPDATE alerts SET status = 'resolved', ends_at = NOW() WHERE fingerprint = ? AND status = 'firing'"
	_, err := r.db.Exec(query, fingerprint)
	return err
}

// GetAlertCount gets total alert count with filters
func (r *AlertRepository) GetAlertCount(status, severity string) (int, error) {
	query := "SELECT COUNT(*) FROM alerts WHERE 1=1"
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

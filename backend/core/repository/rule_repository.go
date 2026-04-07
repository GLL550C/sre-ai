package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

// RuleRepository handles alert rule data access
type RuleRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRuleRepository creates a new rule repository
func NewRuleRepository(db *sql.DB, logger *zap.Logger) *RuleRepository {
	return &RuleRepository{
		db:     db,
		logger: logger,
	}
}

// GetRules retrieves all alert rules
func (r *RuleRepository) GetRules(status int) ([]model.AlertRule, error) {
	query := `SELECT id, name, description, expr, duration, severity, labels, annotations, status, cluster_id, created_at, updated_at
		FROM alert_rules WHERE status = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []model.AlertRule
	for rows.Next() {
		var rule model.AlertRule
		var labels, annotations []byte
		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Expr, &rule.Duration, &rule.Severity,
			&labels, &annotations, &rule.Status, &rule.ClusterID, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan rule", zap.Error(err))
			continue
		}
		if labels != nil {
			rule.Labels = labels
		}
		if annotations != nil {
			rule.Annotations = annotations
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetRuleByID retrieves a rule by ID
func (r *RuleRepository) GetRuleByID(id int64) (*model.AlertRule, error) {
	query := `SELECT id, name, description, expr, duration, severity, labels, annotations, status, cluster_id, created_at, updated_at
		FROM alert_rules WHERE id = ?`

	var rule model.AlertRule
	var labels, annotations []byte
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Expr, &rule.Duration, &rule.Severity,
		&labels, &annotations, &rule.Status, &rule.ClusterID, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if labels != nil {
		rule.Labels = labels
	}
	if annotations != nil {
		rule.Annotations = annotations
	}

	return &rule, nil
}

// CreateRule creates a new alert rule
func (r *RuleRepository) CreateRule(rule *model.AlertRule) error {
	query := `INSERT INTO alert_rules (name, description, expr, duration, severity, labels, annotations, status, cluster_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	labelsJSON, _ := json.Marshal(rule.Labels)
	annotationsJSON, _ := json.Marshal(rule.Annotations)

	result, err := r.db.Exec(query,
		rule.Name, rule.Description, rule.Expr, rule.Duration, rule.Severity,
		labelsJSON, annotationsJSON, rule.Status, rule.ClusterID,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	rule.ID = id
	return nil
}

// UpdateRule updates an alert rule
func (r *RuleRepository) UpdateRule(rule *model.AlertRule) error {
	query := `UPDATE alert_rules SET name = ?, description = ?, expr = ?, duration = ?, severity = ?,
		labels = ?, annotations = ?, status = ?, cluster_id = ? WHERE id = ?`

	labelsJSON, _ := json.Marshal(rule.Labels)
	annotationsJSON, _ := json.Marshal(rule.Annotations)

	_, err := r.db.Exec(query,
		rule.Name, rule.Description, rule.Expr, rule.Duration, rule.Severity,
		labelsJSON, annotationsJSON, rule.Status, rule.ClusterID, rule.ID,
	)
	return err
}

// DeleteRule soft deletes a rule
func (r *RuleRepository) DeleteRule(id int64) error {
	query := "UPDATE alert_rules SET status = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

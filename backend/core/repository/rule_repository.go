package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"

	"go.uber.org/zap"
)

// RuleRepository 告警规则数据访问
type RuleRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRuleRepository 创建仓库
func NewRuleRepository(db *sql.DB, logger *zap.Logger) *RuleRepository {
	return &RuleRepository{db: db, logger: logger}
}

// GetAll 获取所有规则
func (r *RuleRepository) GetAll(status int) ([]model.AlertRule, error) {
	query := `SELECT id, name, description, expr, duration, severity, labels, annotations, status, created_at, updated_at
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
		err := rows.Scan(&rule.ID, &rule.Name, &rule.Description, &rule.Expr, &rule.Duration, &rule.Severity,
			&labels, &annotations, &rule.Status, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			r.logger.Error("扫描规则失败", zap.Error(err))
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

// GetByID 根据ID获取
func (r *RuleRepository) GetByID(id int64) (*model.AlertRule, error) {
	query := `SELECT id, name, description, expr, duration, severity, labels, annotations, status, created_at, updated_at
		FROM alert_rules WHERE id = ?`

	var rule model.AlertRule
	var labels, annotations []byte
	err := r.db.QueryRow(query, id).Scan(&rule.ID, &rule.Name, &rule.Description, &rule.Expr, &rule.Duration, &rule.Severity,
		&labels, &annotations, &rule.Status, &rule.CreatedAt, &rule.UpdatedAt)
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

// Create 创建规则
func (r *RuleRepository) Create(rule *model.AlertRule) error {
	query := `INSERT INTO alert_rules (name, description, expr, duration, severity, labels, annotations, status, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	labelsJSON, _ := json.Marshal(rule.Labels)
	annotationsJSON, _ := json.Marshal(rule.Annotations)

	result, err := r.db.Exec(query, rule.Name, rule.Description, rule.Expr, rule.Duration, rule.Severity,
		labelsJSON, annotationsJSON, rule.Status, rule.CreatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	rule.ID = id
	return nil
}

// Update 更新规则
func (r *RuleRepository) Update(rule *model.AlertRule) error {
	query := `UPDATE alert_rules SET name = ?, description = ?, expr = ?, duration = ?, severity = ?,
		labels = ?, annotations = ?, status = ?, updated_by = ? WHERE id = ?`

	labelsJSON, _ := json.Marshal(rule.Labels)
	annotationsJSON, _ := json.Marshal(rule.Annotations)

	_, err := r.db.Exec(query, rule.Name, rule.Description, rule.Expr, rule.Duration, rule.Severity,
		labelsJSON, annotationsJSON, rule.Status, rule.UpdatedBy, rule.ID)
	return err
}

// Delete 删除规则(软删除)
func (r *RuleRepository) Delete(id int64) error {
	query := `UPDATE alert_rules SET status = 0 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

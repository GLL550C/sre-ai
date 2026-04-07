package repository

import (
	"database/sql"
	"encoding/json"
	"runbook/model"

	"go.uber.org/zap"
)

// RunbookRepository handles runbook data access
type RunbookRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRunbookRepository creates a new runbook repository
func NewRunbookRepository(db *sql.DB, logger *zap.Logger) *RunbookRepository {
	return &RunbookRepository{
		db:     db,
		logger: logger,
	}
}

// GetRunbooks retrieves runbooks with filters
func (r *RunbookRepository) GetRunbooks(alertName, severity string, status int, limit, offset int) ([]model.Runbook, error) {
	query := `SELECT id, title, alert_name, severity, content, steps, related_alerts, created_by, updated_by, view_count, status, created_at, updated_at
		FROM runbooks WHERE status = ?`
	args := []interface{}{status}

	if alertName != "" {
		query += " AND alert_name LIKE ?"
		args = append(args, "%"+alertName+"%")
	}
	if severity != "" {
		query += " AND severity = ?"
		args = append(args, severity)
	}

	query += " ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runbooks []model.Runbook
	for rows.Next() {
		var runbook model.Runbook
		var steps, relatedAlerts []byte
		err := rows.Scan(
			&runbook.ID, &runbook.Title, &runbook.AlertName, &runbook.Severity, &runbook.Content,
			&steps, &relatedAlerts, &runbook.CreatedBy, &runbook.UpdatedBy, &runbook.ViewCount,
			&runbook.Status, &runbook.CreatedAt, &runbook.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan runbook", zap.Error(err))
			continue
		}
		if steps != nil {
			runbook.Steps = steps
		}
		if relatedAlerts != nil {
			runbook.RelatedAlerts = relatedAlerts
		}
		runbooks = append(runbooks, runbook)
	}

	return runbooks, nil
}

// GetRunbookByID retrieves a runbook by ID
func (r *RunbookRepository) GetRunbookByID(id int64) (*model.Runbook, error) {
	query := `SELECT id, title, alert_name, severity, content, steps, related_alerts, created_by, updated_by, view_count, status, created_at, updated_at
		FROM runbooks WHERE id = ?`

	var runbook model.Runbook
	var steps, relatedAlerts []byte
	err := r.db.QueryRow(query, id).Scan(
		&runbook.ID, &runbook.Title, &runbook.AlertName, &runbook.Severity, &runbook.Content,
		&steps, &relatedAlerts, &runbook.CreatedBy, &runbook.UpdatedBy, &runbook.ViewCount,
		&runbook.Status, &runbook.CreatedAt, &runbook.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if steps != nil {
		runbook.Steps = steps
	}
	if relatedAlerts != nil {
		runbook.RelatedAlerts = relatedAlerts
	}

	return &runbook, nil
}

// CreateRunbook creates a new runbook
func (r *RunbookRepository) CreateRunbook(runbook *model.Runbook) error {
	query := `INSERT INTO runbooks (title, alert_name, severity, content, steps, related_alerts, created_by, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	stepsJSON, _ := json.Marshal(runbook.Steps)
	relatedAlertsJSON, _ := json.Marshal(runbook.RelatedAlerts)

	result, err := r.db.Exec(query,
		runbook.Title, runbook.AlertName, runbook.Severity, runbook.Content,
		stepsJSON, relatedAlertsJSON, runbook.CreatedBy, runbook.Status,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	runbook.ID = id
	return nil
}

// UpdateRunbook updates a runbook
func (r *RunbookRepository) UpdateRunbook(runbook *model.Runbook) error {
	query := `UPDATE runbooks SET title = ?, alert_name = ?, severity = ?, content = ?, steps = ?, related_alerts = ?, updated_by = ?, status = ? WHERE id = ?`

	stepsJSON, _ := json.Marshal(runbook.Steps)
	relatedAlertsJSON, _ := json.Marshal(runbook.RelatedAlerts)

	_, err := r.db.Exec(query,
		runbook.Title, runbook.AlertName, runbook.Severity, runbook.Content,
		stepsJSON, relatedAlertsJSON, runbook.UpdatedBy, runbook.Status, runbook.ID,
	)
	return err
}

// DeleteRunbook soft deletes a runbook
func (r *RunbookRepository) DeleteRunbook(id int64) error {
	query := "UPDATE runbooks SET status = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// IncrementViewCount increments the view count
func (r *RunbookRepository) IncrementViewCount(id int64) error {
	query := "UPDATE runbooks SET view_count = view_count + 1 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// SearchRunbooks searches runbooks by keyword
func (r *RunbookRepository) SearchRunbooks(keyword string, limit, offset int) ([]model.Runbook, int, error) {
	countQuery := "SELECT COUNT(*) FROM runbooks WHERE status = 1 AND (title LIKE ? OR content LIKE ? OR alert_name LIKE ?)"
	searchPattern := "%" + keyword + "%"

	var total int
	err := r.db.QueryRow(countQuery, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, title, alert_name, severity, content, steps, related_alerts, created_by, updated_by, view_count, status, created_at, updated_at
		FROM runbooks WHERE status = 1 AND (title LIKE ? OR content LIKE ? OR alert_name LIKE ?)
		ORDER BY updated_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, searchPattern, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var runbooks []model.Runbook
	for rows.Next() {
		var runbook model.Runbook
		var steps, relatedAlerts []byte
		err := rows.Scan(
			&runbook.ID, &runbook.Title, &runbook.AlertName, &runbook.Severity, &runbook.Content,
			&steps, &relatedAlerts, &runbook.CreatedBy, &runbook.UpdatedBy, &runbook.ViewCount,
			&runbook.Status, &runbook.CreatedAt, &runbook.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan runbook", zap.Error(err))
			continue
		}
		if steps != nil {
			runbook.Steps = steps
		}
		if relatedAlerts != nil {
			runbook.RelatedAlerts = relatedAlerts
		}
		runbooks = append(runbooks, runbook)
	}

	return runbooks, total, nil
}

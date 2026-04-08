package repository

import (
	"core/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// AnalysisRepository handles AI analysis data access
type AnalysisRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAnalysisRepository creates a new analysis repository
func NewAnalysisRepository(db *sql.DB, logger *zap.Logger) *AnalysisRepository {
	return &AnalysisRepository{
		db:     db,
		logger: logger,
	}
}

// GetAnalysisWithFilters retrieves AI analysis results with filters
func (r *AnalysisRepository) GetAnalysisWithFilters(clusterID int64, analysisType string, status int, limit, offset int) ([]model.AIAnalysis, error) {
	query := `SELECT id, alert_id, analysis_type,
		input_data, result, confidence, status, created_by, created_at, updated_at
		FROM ai_analysis WHERE status = ?`
	args := []interface{}{status}

	if clusterID > 0 {
		// cluster_id 字段不存在，暂时跳过过滤
		// query += " AND cluster_id = ?"
		// args = append(args, clusterID)
	}
	if analysisType != "" {
		query += " AND analysis_type = ?"
		args = append(args, analysisType)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []model.AIAnalysis
	for rows.Next() {
		var analysis model.AIAnalysis
		var inputData []byte
		err := rows.Scan(
			&analysis.ID, &analysis.AlertID,
			&analysis.AnalysisType, &inputData, &analysis.Result,
			&analysis.Confidence, &analysis.Status, &analysis.CreatedBy,
			&analysis.CreatedAt, &analysis.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan analysis", zap.Error(err))
			continue
		}
		if inputData != nil {
			analysis.InputData = inputData
		}
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetAnalysis retrieves AI analysis results (backward compatible)
func (r *AnalysisRepository) GetAnalysis(alertFingerprint, analysisType string, limit, offset int) ([]model.AIAnalysis, error) {
	return r.GetAnalysisWithFilters(0, analysisType, 1, limit, offset)
}

// GetAnalysisByID retrieves an analysis by ID
func (r *AnalysisRepository) GetAnalysisByID(id int64) (*model.AIAnalysis, error) {
	query := `SELECT id, alert_id, analysis_type,
		input_data, result, confidence, status, created_by, created_at, updated_at
		FROM ai_analysis WHERE id = ?`

	var analysis model.AIAnalysis
	var inputData []byte
	err := r.db.QueryRow(query, id).Scan(
		&analysis.ID, &analysis.AlertID,
		&analysis.AnalysisType, &inputData, &analysis.Result,
		&analysis.Confidence, &analysis.Status, &analysis.CreatedBy,
		&analysis.CreatedAt, &analysis.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if inputData != nil {
		analysis.InputData = inputData
	}

	return &analysis, nil
}

// CreateAnalysis creates a new AI analysis
func (r *AnalysisRepository) CreateAnalysis(analysis *model.AIAnalysis) error {
	query := `INSERT INTO ai_analysis (alert_id, analysis_type,
		input_data, result, confidence, status, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	inputDataJSON, _ := json.Marshal(analysis.InputData)

	result, err := r.db.Exec(query,
		analysis.AlertID, analysis.AnalysisType, inputDataJSON,
		analysis.Result, analysis.Confidence, analysis.Status,
		analysis.CreatedBy, analysis.CreatedAt, analysis.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	analysis.ID = id
	return nil
}

// UpdateAnalysisStatus updates analysis status
func (r *AnalysisRepository) UpdateAnalysisStatus(id int64, status int) error {
	query := "UPDATE ai_analysis SET status = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// DeleteAnalysis soft deletes an analysis
func (r *AnalysisRepository) DeleteAnalysis(id int64) error {
	return r.UpdateAnalysisStatus(id, 0)
}

// GetAnalysisStats returns analysis statistics
func (r *AnalysisRepository) GetAnalysisStats(clusterID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total count
	query := "SELECT COUNT(*) FROM ai_analysis WHERE status = 1"
	args := []interface{}{}

	if clusterID > 0 {
		query += " AND cluster_id = ?"
		args = append(args, clusterID)
	}

	var total int
	if err := r.db.QueryRow(query, args...).Scan(&total); err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by type
	typeQuery := "SELECT analysis_type, COUNT(*) FROM ai_analysis WHERE status = 1"
	if clusterID > 0 {
		typeQuery += " AND cluster_id = ?"
	}
	typeQuery += " GROUP BY analysis_type"

	var typeArgs []interface{}
	if clusterID > 0 {
		typeArgs = append(typeArgs, clusterID)
	}

	rows, err := r.db.Query(typeQuery, typeArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byType := make(map[string]int)
	for rows.Next() {
		var t string
		var c int
		if err := rows.Scan(&t, &c); err == nil {
			byType[t] = c
		}
	}
	stats["by_type"] = byType

	// Average confidence
	confQuery := "SELECT AVG(confidence) FROM ai_analysis WHERE status = 1 AND confidence IS NOT NULL"
	if clusterID > 0 {
		confQuery += " AND cluster_id = ?"
	}

	var avgConfidence float64
	if clusterID > 0 {
		r.db.QueryRow(confQuery, clusterID).Scan(&avgConfidence)
	} else {
		r.db.QueryRow(confQuery).Scan(&avgConfidence)
	}
	stats["avg_confidence"] = fmt.Sprintf("%.1f%%", avgConfidence)

	return stats, nil
}

// GetRecentAnalyses gets recent analyses
func (r *AnalysisRepository) GetRecentAnalyses(limit int) ([]model.AIAnalysis, error) {
	return r.GetAnalysisWithFilters(0, "", 1, limit, 0)
}

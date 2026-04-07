package model

import (
	"encoding/json"
	"time"
)

// Runbook represents an operational runbook
type Runbook struct {
	ID            int64           `json:"id"`
	Title         string          `json:"title"`
	AlertName     string          `json:"alert_name"`
	Severity      string          `json:"severity"`
	Content       string          `json:"content"`
	Steps         json.RawMessage `json:"steps,omitempty"`
	RelatedAlerts json.RawMessage `json:"related_alerts,omitempty"`
	CreatedBy     string          `json:"created_by"`
	UpdatedBy     string          `json:"updated_by"`
	ViewCount     int             `json:"view_count"`
	Status        int             `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// RunbookSearchResult represents search results
type RunbookSearchResult struct {
	Runbooks []Runbook `json:"runbooks"`
	Total    int       `json:"total"`
}

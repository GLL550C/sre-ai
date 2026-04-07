package model

import (
	"encoding/json"
	"time"
)

// Tenant represents a tenant/organization
type Tenant struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Config      json.RawMessage `json:"config,omitempty"`
	Clusters    json.RawMessage `json:"clusters,omitempty"`
	Status      int             `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

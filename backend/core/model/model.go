package model

import (
	"encoding/json"
	"time"
)

// Alert represents an alert record
type Alert struct {
	ID              int64           `json:"id"`
	RuleID          *int64          `json:"rule_id,omitempty"`
	Fingerprint     string          `json:"fingerprint"`
	Status          string          `json:"status"`
	Severity        string          `json:"severity"`
	Summary         string          `json:"summary"`
	Description     string          `json:"description"`
	Labels          json.RawMessage `json:"labels,omitempty"`
	StartsAt        time.Time       `json:"starts_at"`
	EndsAt          *time.Time      `json:"ends_at,omitempty"`
	AcknowledgedBy  *string         `json:"acknowledged_by,omitempty"`
	AcknowledgedAt  *time.Time      `json:"acknowledged_at,omitempty"`
	ClusterID       *int64          `json:"cluster_id,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// AlertRule represents an alert rule
type AlertRule struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Expr        string          `json:"expr"`
	Duration    string          `json:"duration"`
	Severity    string          `json:"severity"`
	Labels      json.RawMessage `json:"labels,omitempty"`
	Annotations json.RawMessage `json:"annotations,omitempty"`
	Status      int             `json:"status"`
	ClusterID   *int64          `json:"cluster_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// PrometheusCluster represents a Prometheus cluster
type PrometheusCluster struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	URL        string          `json:"url"`
	Status     int             `json:"status"`
	IsDefault  bool            `json:"is_default"`
	ConfigJSON json.RawMessage `json:"config_json,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// AIAnalysis represents an AI analysis result
type AIAnalysis struct {
	ID               int64           `json:"id"`
	AlertID          *int64          `json:"alert_id,omitempty"`
	ClusterID        *int64          `json:"cluster_id,omitempty"`
	AlertFingerprint string          `json:"alert_fingerprint"`
	AnalysisType     string          `json:"analysis_type"`
	AnalysisMode     string          `json:"analysis_mode"` // realtime, historical, predictive
	InputData        json.RawMessage `json:"input_data,omitempty"`
	Result           string          `json:"result"`
	RootCause        string          `json:"root_cause,omitempty"`
	Suggestions      []string        `json:"suggestions,omitempty"`
	RelatedAlerts    []int64         `json:"related_alerts,omitempty"`
	Confidence       *float64        `json:"confidence,omitempty"`
	ModelVersion     *string         `json:"model_version,omitempty"`
	Status           int             `json:"status"` // 0=deleted, 1=active, 2=archived
	CreatedBy        string          `json:"created_by,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at,omitempty"`
}

// AnalysisTask represents an analysis task
type AnalysisTask struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	ClusterID    int64      `json:"cluster_id"`
	Query        string     `json:"query"`
	AnalysisType string     `json:"analysis_type"`
	Schedule     string     `json:"schedule,omitempty"` // cron expression for scheduled analysis
	Status       int        `json:"status"`             // 0=disabled, 1=enabled
	LastRunAt    *time.Time `json:"last_run_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AnalysisReport represents a comprehensive analysis report
type AnalysisReport struct {
	ID              int64           `json:"id"`
	TaskID          *int64          `json:"task_id,omitempty"`
	ClusterID       int64           `json:"cluster_id"`
	ClusterName     string          `json:"cluster_name"`
	ReportType      string          `json:"report_type"` // summary, detailed, anomaly
	TimeRange       string          `json:"time_range"`
	Metrics         json.RawMessage `json:"metrics,omitempty"`
	Findings        []Finding       `json:"findings"`
	Recommendations []string        `json:"recommendations"`
	RiskLevel       string          `json:"risk_level"` // low, medium, high, critical
	CreatedAt       time.Time       `json:"created_at"`
}

// Finding represents a single analysis finding
type Finding struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Metric      string  `json:"metric,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Threshold   float64 `json:"threshold,omitempty"`
}

// DashboardConfig represents a dashboard configuration
type DashboardConfig struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	TenantID  *int64          `json:"tenant_id,omitempty"`
	Layout    json.RawMessage `json:"layout,omitempty"`
	Widgets   json.RawMessage `json:"widgets,omitempty"`
	IsDefault bool            `json:"is_default"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// PlatformConfig represents platform configuration
type PlatformConfig struct {
	ID          int64     `json:"id"`
	ConfigKey   string    `json:"config_key"`
	ConfigValue string    `json:"config_value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WebhookPayload represents an alertmanager webhook payload
type WebhookPayload struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []WebhookAlert    `json:"alerts"`
}

// WebhookAlert represents a single alert in webhook payload
type WebhookAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// PrometheusQueryResult represents Prometheus query result
type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value,omitempty"`
			Values [][]interface{}   `json:"values,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// User represents a system user
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"` // Never expose password in JSON
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"` // admin, operator, viewer
	Status       int       `json:"status"` // 1:active, 0:inactive
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP  string    `json:"last_login_ip,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	CaptchaID string `json:"captcha_id"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// Captcha represents a captcha challenge
type Captcha struct {
	ID        string    `json:"id"`
	Image     string    `json:"image"` // base64 encoded image
	ExpiresAt time.Time `json:"expires_at"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

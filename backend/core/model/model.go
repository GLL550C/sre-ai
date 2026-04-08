package model

import (
	"encoding/json"
	"time"
)

// ============================================
// 核心模型
// ============================================

// Tenant 租户
type Tenant struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Status      int8            `json:"status"`
	Settings    json.RawMessage `json:"settings,omitempty"`
	CreatedBy   string          `json:"created_by,omitempty"`
	UpdatedBy   string          `json:"updated_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// User 用户
type User struct {
	ID          int64      `json:"id"`
	TenantID    int64      `json:"tenant_id"`
	Username    string     `json:"username"`
	Password    string     `json:"-"` // 不序列化
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Role        string     `json:"role"`
	Status      int8       `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP string     `json:"last_login_ip,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID          int64     `json:"id"`
	Category    string    `json:"category"`
	SubCategory string    `json:"sub_category,omitempty"`
	Key         string    `json:"config_key"`
	Value       string    `json:"config_value"`
	ValueType   string    `json:"value_type"`
	IsSensitive bool      `json:"is_sensitive"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ============================================
// AI模型
// ============================================

type AIModel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	APIKey      string    `json:"api_key"`
	BaseURL     string    `json:"base_url,omitempty"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Timeout     int       `json:"timeout"`
	IsDefault   bool      `json:"is_default"`
	IsEnabled   bool      `json:"is_enabled"`
	Description string    `json:"description,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AIAnalysis AI分析记录
type AIAnalysis struct {
	ID               int64           `json:"id"`
	AlertID          *int64          `json:"alert_id,omitempty"`
	ClusterID        *int64          `json:"cluster_id,omitempty"`
	AlertFingerprint string          `json:"alert_fingerprint,omitempty"`
	AnalysisType     string          `json:"analysis_type"`
	AnalysisMode     string          `json:"analysis_mode,omitempty"`
	InputData        json.RawMessage `json:"input_data,omitempty"`
	Result           string          `json:"result"`
	RootCause        string          `json:"root_cause,omitempty"`
	Suggestions      []string        `json:"suggestions,omitempty"`
	RelatedAlerts    []int64         `json:"related_alerts,omitempty"`
	Confidence       *float64        `json:"confidence,omitempty"`
	ModelVersion     *string         `json:"model_version,omitempty"`
	Status           int8            `json:"status"`
	CreatedBy        string          `json:"created_by,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// ============================================
// 监控告警
// ============================================

// PrometheusCluster Prometheus集群
type PrometheusCluster struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	URL        string          `json:"url"`
	Status     int8            `json:"status"`
	IsDefault  bool            `json:"is_default"`
	ConfigJSON json.RawMessage `json:"config_json,omitempty"`
	CreatedBy  string          `json:"created_by,omitempty"`
	UpdatedBy  string          `json:"updated_by,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// AlertRule 告警规则
type AlertRule struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Expr        string          `json:"expr"`
	Duration    string          `json:"duration"`
	Severity    string          `json:"severity"`
	Labels      json.RawMessage `json:"labels,omitempty"`
	Annotations json.RawMessage `json:"annotations,omitempty"`
	Status      int8            `json:"status"`
	CreatedBy   string          `json:"created_by,omitempty"`
	UpdatedBy   string          `json:"updated_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Alert 告警记录
type Alert struct {
	ID              int64           `json:"id"`
	RuleID          *int64          `json:"rule_id,omitempty"`
	Fingerprint     string          `json:"fingerprint"`
	Status          string          `json:"status"`
	Severity        string          `json:"severity"`
	Summary         string          `json:"summary,omitempty"`
	Description     string          `json:"description,omitempty"`
	Labels          json.RawMessage `json:"labels,omitempty"`
	StartsAt        time.Time       `json:"starts_at"`
	EndsAt          *time.Time      `json:"ends_at,omitempty"`
	AcknowledgedBy  *string         `json:"acknowledged_by,omitempty"`
	AcknowledgedAt  *time.Time      `json:"acknowledged_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// ============================================
// 运维手册
// ============================================

// Runbook 运维手册
type Runbook struct {
	ID         int64           `json:"id"`
	Title      string          `json:"title"`
	AlertName  string          `json:"alert_name,omitempty"`
	Severity   string          `json:"severity"`
	Content    string          `json:"content"`
	Steps      json.RawMessage `json:"steps,omitempty"`
	ViewCount  int             `json:"view_count"`
	Status     int8            `json:"status"`
	CreatedBy  string          `json:"created_by,omitempty"`
	UpdatedBy  string          `json:"updated_by,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// ============================================
// 请求/响应模型
// ============================================

// LoginRequest 登录请求
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaID string `json:"captcha_id"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// WebhookPayload AlertManager webhook负载
type WebhookPayload struct {
	Version     string       `json:"version"`
	Status      string       `json:"status"`
	Alerts      []WebhookAlert `json:"alerts"`
	ExternalURL string       `json:"externalURL"`
}

// WebhookAlert webhook中的单个告警
type WebhookAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       *time.Time        `json:"endsAt,omitempty"`
	Fingerprint  string            `json:"fingerprint"`
}

// AnalysisReport 分析报告
type AnalysisReport struct {
	ID              int64           `json:"id"`
	ClusterID       int64           `json:"cluster_id"`
	ClusterName     string          `json:"cluster_name"`
	ReportType      string          `json:"report_type"`
	TimeRange       string          `json:"time_range"`
	Metrics         json.RawMessage `json:"metrics,omitempty"`
	Findings        []Finding       `json:"findings"`
	Recommendations []string        `json:"recommendations"`
	RiskLevel       string          `json:"risk_level"`
	CreatedAt       time.Time       `json:"created_at"`
}

// Finding 分析发现
type Finding struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Metric      string  `json:"metric,omitempty"`
	Value       float64 `json:"value,omitempty"`
	Threshold   float64 `json:"threshold,omitempty"`
}

// AnalysisTask 分析任务
type AnalysisTask struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	ClusterID   int64      `json:"cluster_id"`
	Query       string     `json:"query"`
	AnalysisType string    `json:"analysis_type"`
	Schedule    string     `json:"schedule,omitempty"`
	Status      int        `json:"status"`
	LastRunAt   *time.Time `json:"last_run_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Captcha 验证码
type Captcha struct {
	ID        string    `json:"id"`
	Image     string    `json:"image"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

// PrometheusQueryResult Prometheus查询结果
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

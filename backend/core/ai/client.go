package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
		Delta        *struct {
			Content string `json:"content,omitempty"`
		} `json:"delta,omitempty"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// StreamResponse represents a streaming response
type StreamResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// Client interface for AI providers
type Client interface {
	// Chat sends a chat completion request
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatStream sends a streaming chat completion request
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamResponse, error)

	// Health checks if the client is healthy
	Health() error
}

// AnalysisPrompt represents a prompt for analysis
type AnalysisPrompt struct {
	AnalysisType string                 `json:"analysis_type"`
	ClusterName  string                 `json:"cluster_name"`
	ClusterURL   string                 `json:"cluster_url"`
	Metrics      map[string]interface{} `json:"metrics"`
	Alerts       []AlertInfo            `json:"alerts"`
	TimeRange    string                 `json:"time_range"`
	Context      string                 `json:"context"`
}

// AlertInfo represents alert information
type AlertInfo struct {
	Name        string `json:"name"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	Description string `json:"description"`
	StartedAt   string `json:"started_at"`
}

// AnalysisResult represents the AI analysis result
type AnalysisResult struct {
	Summary       string   `json:"summary"`
	RootCause     string   `json:"root_cause"`
	Findings      []string `json:"findings"`
	Suggestions   []string `json:"suggestions"`
	RiskLevel     string   `json:"risk_level"` // low, medium, high, critical
	Confidence    float64  `json:"confidence"`
	RelatedAlerts []string `json:"related_alerts"`
}

// PromptTemplate manages prompt templates
type PromptTemplate struct {
	SystemPrompt string
	UserTemplate string
}

// GetAnalysisPromptTemplate returns the prompt template for analysis
func GetAnalysisPromptTemplate(analysisType string) *PromptTemplate {
	switch analysisType {
	case "root_cause":
		return &PromptTemplate{
			SystemPrompt: `You are an expert SRE (Site Reliability Engineer) AI assistant specializing in root cause analysis.
Your task is to analyze system metrics and alerts to identify the root cause of issues.
Provide clear, actionable insights with specific technical details.`,
			UserTemplate: `Please perform a root cause analysis for the following system state:

Cluster: {{.ClusterName}} ({{.ClusterURL}})
Time Range: {{.TimeRange}}

Active Alerts:
{{range .Alerts}}
- [{{.Severity}}] {{.Name}}: {{.Description}} (Started: {{.StartedAt}})
{{end}}

Metrics Data:
{{.Metrics}}

Additional Context:
{{.Context}}

Please provide your analysis in the following JSON format:
{
  "summary": "Brief summary of the situation",
  "root_cause": "The identified root cause with technical details",
  "findings": ["Key finding 1", "Key finding 2", ...],
  "suggestions": ["Actionable suggestion 1", "Actionable suggestion 2", ...],
  "risk_level": "low|medium|high|critical",
  "confidence": 0.85,
  "related_alerts": ["alert names that are related"]
}`,
		}
	case "trend":
		return &PromptTemplate{
			SystemPrompt: `You are an expert SRE AI assistant specializing in trend analysis and capacity planning.
Analyze metric trends to predict future issues and provide proactive recommendations.`,
			UserTemplate: `Please analyze trends for the following cluster:

Cluster: {{.ClusterName}} ({{.ClusterURL}})
Time Range: {{.TimeRange}}

Historical Metrics:
{{.Metrics}}

Please provide your trend analysis in JSON format with:
- summary: Overview of trends
- findings: List of observed trends
- suggestions: Proactive recommendations
- risk_level: Current risk assessment
- confidence: Confidence score (0-1)`,
		}
	case "anomaly":
		return &PromptTemplate{
			SystemPrompt: `You are an expert SRE AI assistant specializing in anomaly detection.
Identify unusual patterns in metrics that may indicate issues or security concerns.`,
			UserTemplate: `Please detect anomalies for the following cluster:

Cluster: {{.ClusterName}} ({{.ClusterURL}})
Time Range: {{.TimeRange}}

Metrics Data:
{{.Metrics}}

Alerts:
{{range .Alerts}}
- [{{.Severity}}] {{.Name}}: {{.Description}}
{{end}}

Please provide anomaly detection results in JSON format.`,
		}
	case "capacity":
		return &PromptTemplate{
			SystemPrompt: `You are an expert SRE AI assistant specializing in capacity planning.
Analyze resource usage patterns and provide scaling recommendations.`,
			UserTemplate: `Please perform capacity planning analysis for:

Cluster: {{.ClusterName}} ({{.ClusterURL}})
Time Range: {{.TimeRange}}

Resource Metrics:
{{.Metrics}}

Please provide capacity analysis in JSON format with:
- Current utilization assessment
- Growth projections
- Scaling recommendations
- Risk assessment`,
		}
	default:
		return &PromptTemplate{
			SystemPrompt: `You are an expert SRE AI assistant. Analyze the provided data and provide insights.`,
			UserTemplate: `Please analyze the following:

Cluster: {{.ClusterName}}
Time Range: {{.TimeRange}}

Data:
{{.Metrics}}

Alerts:
{{range .Alerts}}
- {{.Name}}: {{.Description}}
{{end}}

Please provide analysis in JSON format.`,
		}
	}
}

// FormatMetrics formats metrics data for prompt
func FormatMetrics(metrics map[string]interface{}) string {
	data, _ := json.MarshalIndent(metrics, "", "  ")
	return string(data)
}

// FormatTimeRange formats time range
func FormatTimeRange(start, end time.Time) string {
	if end.IsZero() || end.Equal(start) {
		return fmt.Sprintf("Last %s", time.Since(start).Round(time.Minute))
	}
	return fmt.Sprintf("%s to %s", start.Format(time.RFC3339), end.Format(time.RFC3339))
}

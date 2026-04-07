package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"core/config"
	"go.uber.org/zap"
)

// AnalysisService provides AI-powered analysis capabilities
type AnalysisService struct {
	client Client
	config *config.AIConfig
	logger *zap.Logger
}

// NewAnalysisService creates a new AI analysis service
func NewAnalysisService(cfg *config.AIConfig, logger *zap.Logger) (*AnalysisService, error) {
	// Convert config to client config
	clientCfg := &ClientConfig{
		Provider:    cfg.Provider,
		APIKey:      cfg.APIKey,
		BaseURL:     cfg.BaseURL,
		Model:       cfg.Model,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Timeout:     cfg.Timeout,
		Enabled:     cfg.Enabled,
	}

	factory := NewClientFactory(clientCfg)
	client, err := factory.CreateClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	return &AnalysisService{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// IsEnabled returns whether AI analysis is enabled
func (s *AnalysisService) IsEnabled() bool {
	return s.config.Enabled && s.config.APIKey != ""
}

// Health checks AI service health
func (s *AnalysisService) Health() error {
	if !s.IsEnabled() {
		return fmt.Errorf("AI analysis is disabled")
	}
	return s.client.Health()
}

// Analyze performs AI analysis on the given prompt
func (s *AnalysisService) Analyze(ctx context.Context, prompt *AnalysisPrompt) (*AnalysisResult, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("AI analysis is disabled or not configured")
	}

	s.logger.Info("Starting AI analysis",
		zap.String("type", prompt.AnalysisType),
		zap.String("cluster", prompt.ClusterName),
	)

	// Get prompt template
	tmpl := GetAnalysisPromptTemplate(prompt.AnalysisType)

	// Build messages
	messages := []Message{
		{Role: "system", Content: tmpl.SystemPrompt},
	}

	// Render user prompt
	userPrompt, err := s.renderPrompt(tmpl.UserTemplate, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to render prompt: %w", err)
	}

	messages = append(messages, Message{Role: "user", Content: userPrompt})

	// Send request to AI
	req := &ChatRequest{
		Model:       s.config.Model,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
	}

	resp, err := s.client.Chat(ctx, req)
	if err != nil {
		s.logger.Error("AI request failed", zap.Error(err))
		return nil, fmt.Errorf("AI request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Parse response
	content := resp.Choices[0].Message.Content
	result, err := s.parseAnalysisResult(content)
	if err != nil {
		s.logger.Warn("Failed to parse structured result, using raw content",
			zap.Error(err),
		)
		// Fallback to raw content
		result = &AnalysisResult{
			Summary:    content,
			Confidence: 0.5,
			RiskLevel:  "unknown",
		}
	}

	s.logger.Info("AI analysis completed",
		zap.String("risk_level", result.RiskLevel),
		zap.Float64("confidence", result.Confidence),
	)

	return result, nil
}

// AnalyzeStream performs streaming AI analysis
func (s *AnalysisService) AnalyzeStream(ctx context.Context, prompt *AnalysisPrompt) (<-chan StreamResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("AI analysis is disabled or not configured")
	}

	tmpl := GetAnalysisPromptTemplate(prompt.AnalysisType)

	messages := []Message{
		{Role: "system", Content: tmpl.SystemPrompt},
	}

	userPrompt, err := s.renderPrompt(tmpl.UserTemplate, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to render prompt: %w", err)
	}

	messages = append(messages, Message{Role: "user", Content: userPrompt})

	req := &ChatRequest{
		Model:       s.config.Model,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Stream:      true,
	}

	return s.client.ChatStream(ctx, req)
}

// Chat performs a conversational chat with context
func (s *AnalysisService) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("AI analysis is disabled or not configured")
	}

	req := &ChatRequest{
		Model:       s.config.Model,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
	}

	return s.client.Chat(ctx, req)
}

// ChatStream performs streaming conversational chat
func (s *AnalysisService) ChatStream(ctx context.Context, messages []Message) (<-chan StreamResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("AI analysis is disabled or not configured")
	}

	req := &ChatRequest{
		Model:       s.config.Model,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Stream:      true,
	}

	return s.client.ChatStream(ctx, req)
}

// renderPrompt renders the prompt template with data
func (s *AnalysisService) renderPrompt(templateStr string, data *AnalysisPrompt) (string, error) {
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// parseAnalysisResult parses the AI response into structured result
func (s *AnalysisService) parseAnalysisResult(content string) (*AnalysisResult, error) {
	// Try to extract JSON from the response
	// The AI might wrap JSON in markdown code blocks
	jsonStr := content

	// Remove markdown code blocks if present
	if idx := strings.Index(content, "```json"); idx != -1 {
		start := idx + 7
		if end := strings.Index(content[start:], "```"); end != -1 {
			jsonStr = strings.TrimSpace(content[start : start+end])
		}
	} else if idx := strings.Index(content, "```"); idx != -1 {
		start := idx + 3
		if end := strings.Index(content[start:], "```"); end != -1 {
			jsonStr = strings.TrimSpace(content[start : start+end])
		}
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON result: %w", err)
	}

	// Validate required fields
	if result.Summary == "" {
		result.Summary = "Analysis completed"
	}
	if result.RiskLevel == "" {
		result.RiskLevel = "unknown"
	}
	if result.Confidence == 0 {
		result.Confidence = 0.5
	}

	return &result, nil
}

// BuildAnalysisPrompt builds analysis prompt from metrics and alerts
func (s *AnalysisService) BuildAnalysisPrompt(
	analysisType string,
	clusterName string,
	clusterURL string,
	metrics map[string]interface{},
	alerts []AlertInfo,
	timeRange string,
	context string,
) *AnalysisPrompt {
	return &AnalysisPrompt{
		AnalysisType: analysisType,
		ClusterName:  clusterName,
		ClusterURL:   clusterURL,
		Metrics:      metrics,
		Alerts:       alerts,
		TimeRange:    timeRange,
		Context:      context,
	}
}

// QuickAnalyze performs a quick analysis with minimal parameters
func (s *AnalysisService) QuickAnalyze(ctx context.Context, analysisType, clusterName, query string) (*AnalysisResult, error) {
	prompt := &AnalysisPrompt{
		AnalysisType: analysisType,
		ClusterName:  clusterName,
		TimeRange:    "current",
		Context:      query,
	}

	return s.Analyze(ctx, prompt)
}

// GetModelInfo returns current model information
func (s *AnalysisService) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"provider":    s.config.Provider,
		"model":       s.config.Model,
		"enabled":     s.config.Enabled,
		"max_tokens":  s.config.MaxTokens,
		"temperature": s.config.Temperature,
	}
}

// ConversationManager manages conversation context
type ConversationManager struct {
	sessions map[string][]Message
	timeout  time.Duration
}

// NewConversationManager creates a new conversation manager
func NewConversationManager(timeout time.Duration) *ConversationManager {
	if timeout == 0 {
		timeout = 30 * time.Minute
	}

	return &ConversationManager{
		sessions: make(map[string][]Message),
		timeout:  timeout,
	}
}

// GetSession gets or creates a conversation session
func (cm *ConversationManager) GetSession(sessionID string) []Message {
	if messages, ok := cm.sessions[sessionID]; ok {
		return messages
	}
	return []Message{}
}

// AddMessage adds a message to a session
func (cm *ConversationManager) AddMessage(sessionID string, msg Message) {
	cm.sessions[sessionID] = append(cm.sessions[sessionID], msg)
}

// ClearSession clears a conversation session
func (cm *ConversationManager) ClearSession(sessionID string) {
	delete(cm.sessions, sessionID)
}

// GetSystemPrompt returns the system prompt for SRE assistant
func GetSystemPrompt() string {
	return `You are an expert Site Reliability Engineer (SRE) AI assistant with deep knowledge of:
- Prometheus monitoring and metrics analysis
- Kubernetes and container orchestration
- Distributed systems troubleshooting
- Incident response and root cause analysis
- Performance optimization
- Capacity planning

Your responsibilities:
1. Analyze system metrics and alerts to identify issues
2. Provide root cause analysis with specific technical details
3. Suggest actionable remediation steps
4. Predict potential future issues based on trends
5. Correlate related alerts to find patterns

Guidelines:
- Always provide specific, actionable recommendations
- Include confidence levels for your analysis
- Cite specific metrics or patterns you observed
- Consider both immediate fixes and long-term improvements
- Be concise but thorough in your analysis

Respond in the requested format (usually JSON) with clear structure.`
}

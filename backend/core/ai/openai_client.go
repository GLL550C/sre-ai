package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAIClient implements Client interface for OpenAI API
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey, baseURL, model string, timeout int) *OpenAIClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4"
	}
	if timeout == 0 {
		timeout = 60
	}

	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		model: model,
	}
}

// Chat sends a chat completion request
func (c *OpenAIClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = c.model
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ChatStream sends a streaming chat completion request
func (c *OpenAIClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamResponse, error) {
	if req.Model == "" {
		req.Model = c.model
	}
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	stream := make(chan StreamResponse)
	go func() {
		defer close(stream)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					stream <- StreamResponse{Error: err.Error(), Done: true}
				} else {
					stream <- StreamResponse{Done: true}
				}
				return
			}

			line = string(bytes.TrimSpace([]byte(line)))
			if len(line) == 0 {
				continue
			}

			// Skip "data: [DONE]"
			if string(line) == "data: [DONE]" {
				stream <- StreamResponse{Done: true}
				return
			}

			// Parse SSE data
			if strings.HasPrefix(line, "data: ") {
				data := []byte(strings.TrimPrefix(line, "data: "))

				var chunk ChatResponse
				if err := json.Unmarshal(data, &chunk); err != nil {
					continue
				}

				if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
					content := chunk.Choices[0].Delta.Content
					stream <- StreamResponse{Content: content}
				}
			}
		}
	}()

	return stream, nil
}

// Health checks if the client is healthy
func (c *OpenAIClient) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &ChatRequest{
		Model:    c.model,
		Messages: []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 5,
	}

	_, err := c.Chat(ctx, req)
	return err
}

// ClaudeClient implements Client interface for Anthropic Claude API
type ClaudeClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient(apiKey, baseURL, model string, timeout int) *ClaudeClient {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	if model == "" {
		model = "claude-3-opus-20240229"
	}
	if timeout == 0 {
		timeout = 60
	}

	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		model: model,
	}
}

// ClaudeRequest represents Claude API request
type ClaudeRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	System      string    `json:"system,omitempty"`
}

// ClaudeResponse represents Claude API response
type ClaudeResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Role         string `json:"role"`
	Content      []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// Chat sends a chat completion request
func (c *ClaudeClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	claudeReq := ClaudeRequest{
		Model:       c.model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stream:      false,
	}

	// Extract system message if present
	for i, msg := range claudeReq.Messages {
		if msg.Role == "system" {
			claudeReq.System = msg.Content
			claudeReq.Messages = append(claudeReq.Messages[:i], claudeReq.Messages[i+1:]...)
			break
		}
	}

	body, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to standard format
	result := &ChatResponse{
		ID:      claudeResp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   claudeResp.Model,
	}

	if len(claudeResp.Content) > 0 {
		result.Choices = []struct {
			Index        int     `json:"index"`
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
			Delta        *struct {
				Content string `json:"content,omitempty"`
			} `json:"delta,omitempty"`
		}{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: claudeResp.Content[0].Text,
				},
				FinishReason: claudeResp.StopReason,
			},
		}
		result.Usage.PromptTokens = claudeResp.Usage.InputTokens
		result.Usage.CompletionTokens = claudeResp.Usage.OutputTokens
		result.Usage.TotalTokens = claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens
	}

	return result, nil
}

// ChatStream sends a streaming chat completion request
func (c *ClaudeClient) ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamResponse, error) {
	// Claude streaming implementation similar to OpenAI
	// For simplicity, fallback to non-streaming
	stream := make(chan StreamResponse)
	go func() {
		defer close(stream)

		resp, err := c.Chat(ctx, req)
		if err != nil {
			stream <- StreamResponse{Error: err.Error(), Done: true}
			return
		}

		if len(resp.Choices) > 0 {
			stream <- StreamResponse{Content: resp.Choices[0].Message.Content}
		}
		stream <- StreamResponse{Done: true}
	}()

	return stream, nil
}

// Health checks if the client is healthy
func (c *ClaudeClient) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &ChatRequest{
		Messages: []Message{{Role: "user", Content: "Hi"}},
		MaxTokens: 5,
	}

	_, err := c.Chat(ctx, req)
	return err
}

// ClientFactory creates AI clients
type ClientFactory struct {
	config *ClientConfig
}

// ClientConfig represents AI client configuration
type ClientConfig struct {
	Provider    string
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
	Timeout     int
	Enabled     bool
}

// NewClientFactory creates a new client factory
func NewClientFactory(config *ClientConfig) *ClientFactory {
	return &ClientFactory{config: config}
}

// CreateClient creates an AI client based on configuration
func (f *ClientFactory) CreateClient() (Client, error) {
	switch f.config.Provider {
	case "openai", "azure", "custom":
		return NewOpenAIClient(f.config.APIKey, f.config.BaseURL, f.config.Model, f.config.Timeout), nil
	case "claude", "anthropic":
		return NewClaudeClient(f.config.APIKey, f.config.BaseURL, f.config.Model, f.config.Timeout), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", f.config.Provider)
	}
}

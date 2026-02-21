package server

import "encoding/json"

// ServerConfig holds server-level configuration.
type ServerConfig struct {
	APIKey       string
	CLIPath      string
	WorkDir      string
	DefaultModel string
	MaxBudget    float64
	MaxTurns     int
}

// --- Anthropic Messages API Request Types ---

// MessagesRequest is the request body for POST /v1/messages.
type MessagesRequest struct {
	Model         string            `json:"model"`
	MaxTokens     int               `json:"max_tokens"`
	Messages      []Message         `json:"messages"`
	System        json.RawMessage   `json:"system,omitempty"`
	Stream        bool              `json:"stream,omitempty"`
	Temperature   *float64          `json:"temperature,omitempty"`
	TopP          *float64          `json:"top_p,omitempty"`
	TopK          *int              `json:"top_k,omitempty"`
	StopSequences []string          `json:"stop_sequences,omitempty"`
	Metadata      *RequestMetadata  `json:"metadata,omitempty"`
	Tools         []Tool            `json:"tools,omitempty"`
	ToolChoice    json.RawMessage   `json:"tool_choice,omitempty"`
}

// Message represents a single message in the conversation.
type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// ContentBlock represents a content block in a message.
type ContentBlock struct {
	Type   string          `json:"type"`
	Text   string          `json:"text,omitempty"`
	Source json.RawMessage `json:"source,omitempty"`
	ID     string          `json:"id,omitempty"`
	Name   string          `json:"name,omitempty"`
	Input  json.RawMessage `json:"input,omitempty"`
}

// SystemBlock represents a system prompt block.
type SystemBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// RequestMetadata holds optional request metadata.
type RequestMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

// Tool represents a tool definition.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
}

// --- Anthropic Messages API Response Types ---

// MessagesResponse is the response body for POST /v1/messages (non-streaming).
type MessagesResponse struct {
	ID           string              `json:"id"`
	Type         string              `json:"type"`
	Role         string              `json:"role"`
	Content      []ResponseContent   `json:"content"`
	Model        string              `json:"model"`
	StopReason   *string             `json:"stop_reason"`
	StopSequence *string             `json:"stop_sequence"`
	Usage        MessagesUsage       `json:"usage"`
}

// ResponseContent represents a content block in the response.
type ResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// MessagesUsage holds token usage for the response.
type MessagesUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// --- Error Types ---

// ErrorResponse is the Anthropic API error format.
type ErrorResponse struct {
	Type  string      `json:"type"`
	Error ErrorDetail `json:"error"`
}

// ErrorDetail holds error details.
type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// --- Health ---

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status string `json:"status"`
}

// --- Quiz API Types ---

// QuizType represents the type of quiz question.
type QuizType string

const (
	QuizTypeMultipleChoice QuizType = "multiple_choice"
	QuizTypeShortAnswer    QuizType = "short_answer"
	QuizTypeEssay          QuizType = "essay"
)

// QuizRequest is the request body for POST /v1/quiz.
type QuizRequest struct {
	Model    string   `json:"model,omitempty"`
	Question string   `json:"question"`
	Type     QuizType `json:"type"`
	Options  []string `json:"options,omitempty"`
	Answer   string   `json:"answer"`
}

// QuizResult holds the grading result from Claude.
type QuizResult struct {
	Correct     bool   `json:"correct"`
	Score       int    `json:"score"`
	Feedback    string `json:"feedback"`
	ModelAnswer string `json:"model_answer"`
}

// QuizResponse is the response body for POST /v1/quiz.
type QuizResponse struct {
	ID     string        `json:"id"`
	Result QuizResult    `json:"result"`
	Model  string        `json:"model"`
	Usage  MessagesUsage `json:"usage"`
}

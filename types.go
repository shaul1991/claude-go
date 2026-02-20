package claude

import "encoding/json"

// OutputFormat specifies the output format for the claude CLI.
type OutputFormat string

const (
	FormatText       OutputFormat = "text"
	FormatJSON       OutputFormat = "json"
	FormatStreamJSON OutputFormat = "stream-json"
)

// Response represents the parsed JSON response from claude -p --output-format json.
type Response struct {
	SessionID string `json:"session_id"`
	Result    string `json:"result"`
	Cost      Cost   `json:"cost_usd"`
	Usage     Usage  `json:"usage"`
	Model     string `json:"model"`
	Duration  int    `json:"duration_ms"`
}

// Cost holds token cost information.
type Cost struct {
	InputTokens         int     `json:"input_tokens"`
	OutputTokens        int     `json:"output_tokens"`
	CacheReadTokens     int     `json:"cache_read_input_tokens"`
	CacheCreationTokens int     `json:"cache_creation_input_tokens"`
	TotalUSD            float64 `json:"total_usd,omitempty"`
}

// Usage holds token usage counters.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// StreamEvent represents a single event from claude -p --output-format stream-json.
type StreamEvent struct {
	Type  string          `json:"type"`
	Event json.RawMessage `json:"event,omitempty"`
}

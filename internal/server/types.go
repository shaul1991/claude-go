package server

// ClientOptions maps to the claude.Option functional options.
type ClientOptions struct {
	Model              string   `json:"model,omitempty"`
	SystemPrompt       string   `json:"system_prompt,omitempty"`
	AppendSystemPrompt string   `json:"append_system_prompt,omitempty"`
	AllowedTools       []string `json:"allowed_tools,omitempty"`
	MaxTurns           int      `json:"max_turns,omitempty"`
	MaxBudgetUSD       float64  `json:"max_budget_usd,omitempty"`
	WorkDir            string   `json:"work_dir,omitempty"`
	CLIPath            string   `json:"cli_path,omitempty"`
}

// AskRequest is the request body for POST /api/v1/ask.
type AskRequest struct {
	Prompt  string         `json:"prompt"`
	Options *ClientOptions `json:"options,omitempty"`
}

// AskJSONRequest is the request body for POST /api/v1/ask-json.
type AskJSONRequest struct {
	Prompt  string         `json:"prompt"`
	Options *ClientOptions `json:"options,omitempty"`
}

// AskWithSchemaRequest is the request body for POST /api/v1/ask-with-schema.
type AskWithSchemaRequest struct {
	Prompt  string         `json:"prompt"`
	Schema  string         `json:"schema"`
	Options *ClientOptions `json:"options,omitempty"`
}

// ResumeRequest is the request body for POST /api/v1/resume.
type ResumeRequest struct {
	SessionID string         `json:"session_id"`
	Prompt    string         `json:"prompt"`
	Options   *ClientOptions `json:"options,omitempty"`
}

// ContinueRequest is the request body for POST /api/v1/continue.
type ContinueRequest struct {
	Prompt  string         `json:"prompt"`
	Options *ClientOptions `json:"options,omitempty"`
}

// StreamRequest is the request body for POST /api/v1/stream.
type StreamRequest struct {
	Prompt  string         `json:"prompt"`
	Options *ClientOptions `json:"options,omitempty"`
}

// TextResponse wraps a plain-text result.
type TextResponse struct {
	Result string `json:"result"`
}

// ErrorResponse is returned on errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status string `json:"status"`
}

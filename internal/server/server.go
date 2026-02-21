package server

import (
	"encoding/json"
	"net/http"

	claude "github.com/shaul1991/claude-go"
)

// Server is the HTTP API server implementing the Anthropic Messages API.
type Server struct {
	mux    *http.ServeMux
	config ServerConfig
}

// NewServer creates a new Server with the given config and registers routes.
func NewServer(config ServerConfig) *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		config: config,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /v1/messages", s.handleMessages)
	s.mux.HandleFunc("POST /v1/quiz", s.handleQuiz)
}

// ServeHTTP implements http.Handler with middleware chain.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := loggingMiddleware(corsMiddleware(authMiddleware(s.config.APIKey, s.mux)))
	handler.ServeHTTP(w, r)
}

// buildClient creates a claude.Client from the request and server config.
func (s *Server) buildClient(req *MessagesRequest, systemPrompt string) *claude.Client {
	var opts []claude.Option

	if req.Model != "" {
		opts = append(opts, claude.WithModel(req.Model))
	}
	if systemPrompt != "" {
		opts = append(opts, claude.WithSystemPrompt(systemPrompt))
	}
	if s.config.CLIPath != "" {
		opts = append(opts, claude.WithCLIPath(s.config.CLIPath))
	}
	if s.config.WorkDir != "" {
		opts = append(opts, claude.WithWorkDir(s.config.WorkDir))
	}
	if s.config.MaxBudget > 0 {
		opts = append(opts, claude.WithMaxBudget(s.config.MaxBudget))
	}
	if s.config.MaxTurns > 0 {
		opts = append(opts, claude.WithMaxTurns(s.config.MaxTurns))
	}

	return claude.NewClient(opts...)
}

// buildQuizClient creates a claude.Client configured for quiz grading.
func (s *Server) buildQuizClient(model, systemPrompt string) *claude.Client {
	var opts []claude.Option

	if model != "" {
		opts = append(opts, claude.WithModel(model))
	}
	if systemPrompt != "" {
		opts = append(opts, claude.WithSystemPrompt(systemPrompt))
	}
	if s.config.CLIPath != "" {
		opts = append(opts, claude.WithCLIPath(s.config.CLIPath))
	}
	if s.config.WorkDir != "" {
		opts = append(opts, claude.WithWorkDir(s.config.WorkDir))
	}
	if s.config.MaxBudget > 0 {
		opts = append(opts, claude.WithMaxBudget(s.config.MaxBudget))
	}
	if s.config.MaxTurns > 0 {
		opts = append(opts, claude.WithMaxTurns(s.config.MaxTurns))
	}

	return claude.NewClient(opts...)
}

// respondError writes an Anthropic-format error response.
func respondError(w http.ResponseWriter, status int, errType string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Type: "error",
		Error: ErrorDetail{
			Type:    errType,
			Message: message,
		},
	})
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

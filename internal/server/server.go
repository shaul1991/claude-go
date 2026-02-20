package server

import (
	"encoding/json"
	"net/http"

	claude "github.com/shaul1991/claude-go"
)

// Server is the HTTP API server wrapping the claude-go library.
type Server struct {
	mux *http.ServeMux
}

// NewServer creates a new Server with all routes registered.
func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /api/v1/ask", s.handleAsk)
	s.mux.HandleFunc("POST /api/v1/ask-json", s.handleAskJSON)
	s.mux.HandleFunc("POST /api/v1/ask-with-schema", s.handleAskWithSchema)
	s.mux.HandleFunc("POST /api/v1/resume", s.handleResume)
	s.mux.HandleFunc("POST /api/v1/continue", s.handleContinue)
	s.mux.HandleFunc("POST /api/v1/stream", s.handleStream)
}

// ServeHTTP implements http.Handler with middleware chain.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := loggingMiddleware(corsMiddleware(s.mux))
	handler.ServeHTTP(w, r)
}

// buildClient creates a claude.Client from ClientOptions.
func buildClient(opts *ClientOptions) *claude.Client {
	if opts == nil {
		return claude.NewClient()
	}

	var options []claude.Option

	if opts.Model != "" {
		options = append(options, claude.WithModel(opts.Model))
	}
	if opts.SystemPrompt != "" {
		options = append(options, claude.WithSystemPrompt(opts.SystemPrompt))
	}
	if opts.AppendSystemPrompt != "" {
		options = append(options, claude.WithAppendSystemPrompt(opts.AppendSystemPrompt))
	}
	if len(opts.AllowedTools) > 0 {
		options = append(options, claude.WithAllowedTools(opts.AllowedTools...))
	}
	if opts.MaxTurns > 0 {
		options = append(options, claude.WithMaxTurns(opts.MaxTurns))
	}
	if opts.MaxBudgetUSD > 0 {
		options = append(options, claude.WithMaxBudget(opts.MaxBudgetUSD))
	}
	if opts.WorkDir != "" {
		options = append(options, claude.WithWorkDir(opts.WorkDir))
	}
	if opts.CLIPath != "" {
		options = append(options, claude.WithCLIPath(opts.CLIPath))
	}

	return claude.NewClient(options...)
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes a JSON error response.
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{Error: msg})
}

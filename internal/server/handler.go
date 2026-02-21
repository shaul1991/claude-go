package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	var req MessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON in request body")
		return
	}

	if err := validateRequest(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	systemPrompt := extractSystemPrompt(req.System)
	prompt := extractPrompt(req.Messages)

	if req.Stream {
		s.handleStream(w, r, &req, systemPrompt, prompt)
	} else {
		s.handleNonStream(w, r, &req, systemPrompt, prompt)
	}
}

func (s *Server) handleNonStream(w http.ResponseWriter, r *http.Request, req *MessagesRequest, systemPrompt, prompt string) {
	client := s.buildClient(req, systemPrompt)
	resp, err := client.AskJSON(r.Context(), prompt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "api_error", err.Error())
		return
	}

	stopReason := "end_turn"
	respondJSON(w, http.StatusOK, MessagesResponse{
		ID:         generateMsgID(resp.SessionID),
		Type:       "message",
		Role:       "assistant",
		Content:    []ResponseContent{{Type: "text", Text: resp.Result}},
		Model:      resp.Model,
		StopReason: &stopReason,
		Usage: MessagesUsage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
		},
	})
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request, req *MessagesRequest, systemPrompt, prompt string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, "api_error", "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := s.buildClient(req, systemPrompt)
	events, errc := client.AskStream(r.Context(), prompt)

	for ev := range events {
		if ev.Type != "stream_event" || ev.Event == nil {
			continue
		}

		var inner struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(ev.Event, &inner); err != nil {
			continue
		}

		writeSSE(w, flusher, inner.Type, ev.Event)
	}

	if err := <-errc; err != nil {
		errData, _ := json.Marshal(ErrorResponse{
			Type: "error",
			Error: ErrorDetail{
				Type:    "api_error",
				Message: err.Error(),
			},
		})
		writeSSE(w, flusher, "error", errData)
	}
}

// validateRequest checks required fields in the Messages API request.
func validateRequest(req *MessagesRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if req.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens is required and must be > 0")
	}
	if len(req.Messages) == 0 {
		return fmt.Errorf("messages is required and must not be empty")
	}
	return nil
}

// extractPrompt converts the messages array into a single prompt string for the CLI.
func extractPrompt(messages []Message) string {
	if len(messages) == 1 {
		return contentToText(messages[0].Content)
	}

	var parts []string
	for _, m := range messages {
		text := contentToText(m.Content)
		switch m.Role {
		case "user":
			parts = append(parts, "Human: "+text)
		case "assistant":
			parts = append(parts, "Assistant: "+text)
		default:
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n\n")
}

// contentToText extracts text from a message's content field.
// Content can be a plain string or an array of content blocks.
func contentToText(raw json.RawMessage) string {
	// Try string first
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	// Try array of content blocks
	var blocks []ContentBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		var texts []string
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				texts = append(texts, b.Text)
			}
		}
		return strings.Join(texts, "\n")
	}

	return string(raw)
}

// extractSystemPrompt parses the system field which can be a string or []SystemBlock.
func extractSystemPrompt(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	// Try string first
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	// Try array of system blocks
	var blocks []SystemBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		var texts []string
		for _, b := range blocks {
			if b.Text != "" {
				texts = append(texts, b.Text)
			}
		}
		return strings.Join(texts, "\n")
	}

	return ""
}

// generateMsgID creates a message ID with msg_ prefix.
func generateMsgID(sessionID string) string {
	if sessionID != "" {
		return "msg_" + sessionID
	}
	b := make([]byte, 16)
	rand.Read(b)
	return "msg_" + hex.EncodeToString(b)
}

// writeSSE writes a single SSE event and flushes.
func writeSSE(w http.ResponseWriter, flusher http.Flusher, event string, data []byte) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
	flusher.Flush()
}

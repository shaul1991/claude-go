package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (s *Server) handleAsk(w http.ResponseWriter, r *http.Request) {
	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	client := buildClient(req.Options)
	result, err := client.Ask(r.Context(), req.Prompt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, TextResponse{Result: result})
}

func (s *Server) handleAskJSON(w http.ResponseWriter, r *http.Request) {
	var req AskJSONRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	client := buildClient(req.Options)
	resp, err := client.AskJSON(r.Context(), req.Prompt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAskWithSchema(w http.ResponseWriter, r *http.Request) {
	var req AskWithSchemaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}
	if req.Schema == "" {
		respondError(w, http.StatusBadRequest, "schema is required")
		return
	}

	client := buildClient(req.Options)
	resp, err := client.AskWithSchema(r.Context(), req.Prompt, req.Schema)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) handleResume(w http.ResponseWriter, r *http.Request) {
	var req ResumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.SessionID == "" {
		respondError(w, http.StatusBadRequest, "session_id is required")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	client := buildClient(req.Options)
	resp, err := client.Resume(r.Context(), req.SessionID, req.Prompt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) handleContinue(w http.ResponseWriter, r *http.Request) {
	var req ContinueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	client := buildClient(req.Options)
	resp, err := client.Continue(r.Context(), req.Prompt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	var req StreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Prompt == "" {
		respondError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := buildClient(req.Options)
	events, errc := client.AskStream(r.Context(), req.Prompt)

	for ev := range events {
		data, err := json.Marshal(ev)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\":%q}\n\n", err.Error())
			flusher.Flush()
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	if err := <-errc; err != nil {
		fmt.Fprintf(w, "data: {\"error\":%q}\n\n", err.Error())
		flusher.Flush()
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

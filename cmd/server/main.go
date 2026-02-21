package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shaul1991/claude-go/internal/server"
)

func main() {
	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}

	port := flag.String("port", defaultPort, "server port")
	host := flag.String("host", "0.0.0.0", "server host")
	apiKey := flag.String("api-key", os.Getenv("API_KEY"), "API key for authentication (empty = no auth)")
	cliPath := flag.String("cli-path", os.Getenv("CLAUDE_CLI_PATH"), "path to claude CLI binary")
	workDir := flag.String("work-dir", os.Getenv("CLAUDE_WORK_DIR"), "working directory for claude CLI")
	defaultModel := flag.String("model", envOrDefault("CLAUDE_MODEL", "opus"), "default model when not specified in request")
	maxBudget := flag.Float64("max-budget", 0, "max budget in USD per request")
	maxTurns := flag.Int("max-turns", 0, "max turns per request")
	flag.Parse()

	config := server.ServerConfig{
		APIKey:       *apiKey,
		CLIPath:      *cliPath,
		WorkDir:      *workDir,
		DefaultModel: *defaultModel,
		MaxBudget:    *maxBudget,
		MaxTurns:     *maxTurns,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		Handler:      server.NewServer(config),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("server stopped")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	claude "github.com/jihoonkim/claude-go"
)

func main() {
	client := claude.NewClient(
		claude.WithMaxTurns(3),
	)

	ctx := context.Background()

	// 1. Simple text question
	fmt.Println("=== Ask (text) ===")
	answer, err := client.Ask(ctx, "What is 2+2? Answer with just the number.")
	if err != nil {
		log.Fatalf("Ask failed: %v", err)
	}
	fmt.Println(answer)

	// 2. JSON response
	fmt.Println("\n=== AskJSON ===")
	resp, err := client.AskJSON(ctx, "Say hello in one word.")
	if err != nil {
		log.Fatalf("AskJSON failed: %v", err)
	}
	fmt.Printf("Result: %s\nSession: %s\nUsage: %+v\n", resp.Result, resp.SessionID, resp.Usage)

	// 3. Resume session
	if resp.SessionID != "" {
		fmt.Println("\n=== Resume ===")
		resp2, err := client.Resume(ctx, resp.SessionID, "What did I just ask you?")
		if err != nil {
			log.Fatalf("Resume failed: %v", err)
		}
		fmt.Println(resp2.Result)
	}

	// 4. Pipe
	fmt.Println("\n=== Pipe ===")
	reader := strings.NewReader("error: connection refused at line 42\npanic: nil pointer dereference")
	answer, err = client.Pipe(ctx, reader, "Summarize these errors in one sentence.")
	if err != nil {
		log.Fatalf("Pipe failed: %v", err)
	}
	fmt.Println(answer)

	// 5. Streaming
	fmt.Println("\n=== AskStream ===")
	events, errc := client.AskStream(ctx, "Count from 1 to 3.")
	for ev := range events {
		fmt.Printf("[%s] %s\n", ev.Type, string(ev.Event))
	}
	if err := <-errc; err != nil {
		fmt.Fprintf(os.Stderr, "stream error: %v\n", err)
	}
}

package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AskStream runs the prompt with stream-json output and returns channels for
// events and errors. The events channel is closed when the stream ends.
// The error channel receives at most one error, then is closed.
// Cancelling the context will kill the underlying process.
func (c *Client) AskStream(ctx context.Context, prompt string) (<-chan StreamEvent, <-chan error) {
	events := make(chan StreamEvent)
	errc := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errc)

		args := c.buildArgs(prompt, FormatStreamJSON)
		cmd := c.newCmd(ctx, args)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errc <- fmt.Errorf("claude: stdout pipe: %w", err)
			return
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			errc <- fmt.Errorf("claude: start: %w", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var ev StreamEvent
			if err := json.Unmarshal(line, &ev); err != nil {
				errc <- fmt.Errorf("claude: parse stream event: %w", err)
				return
			}

			select {
			case events <- ev:
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errc <- fmt.Errorf("claude: read stream: %w", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			msg := strings.TrimSpace(stderr.String())
			if msg != "" {
				errc <- fmt.Errorf("claude: %s", msg)
			} else {
				errc <- wrapExecError(err)
			}
		}
	}()

	return events, errc
}

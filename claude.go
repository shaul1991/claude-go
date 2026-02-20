package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

// Client wraps the claude CLI.
type Client struct {
	cliPath      string
	model        string
	systemPrompt string
	appendPrompt string
	allowedTools []string
	maxTurns     int
	maxBudget    float64
	workDir      string
}

// NewClient creates a new Client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		cliPath: "claude",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// buildArgs assembles the CLI arguments for a given prompt and output format.
// Extra flags (e.g. --resume, --continue) can be appended via extra.
func (c *Client) buildArgs(prompt string, format OutputFormat, extra ...string) []string {
	args := []string{"-p", prompt, "--output-format", string(format)}

	if format == FormatStreamJSON {
		args = append(args, "--verbose")
	}

	if c.model != "" {
		args = append(args, "--model", c.model)
	}
	if c.systemPrompt != "" {
		args = append(args, "--system-prompt", c.systemPrompt)
	}
	if c.appendPrompt != "" {
		args = append(args, "--append-system-prompt", c.appendPrompt)
	}
	for _, tool := range c.allowedTools {
		args = append(args, "--allowedTools", tool)
	}
	if c.maxTurns > 0 {
		args = append(args, "--max-turns", strconv.Itoa(c.maxTurns))
	}
	if c.maxBudget > 0 {
		args = append(args, "--max-budget-usd", strconv.FormatFloat(c.maxBudget, 'f', -1, 64))
	}
	args = append(args, extra...)
	return args
}

// newCmd creates an *exec.Cmd with working directory set if configured.
func (c *Client) newCmd(ctx context.Context, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.cliPath, args...)
	if c.workDir != "" {
		cmd.Dir = c.workDir
	}
	return cmd
}

// Ask runs the prompt and returns the plain-text response.
func (c *Client) Ask(ctx context.Context, prompt string) (string, error) {
	args := c.buildArgs(prompt, FormatText)
	cmd := c.newCmd(ctx, args)

	out, err := cmd.Output()
	if err != nil {
		return "", wrapExecError(err)
	}
	return string(bytes.TrimSpace(out)), nil
}

// AskJSON runs the prompt with JSON output and returns a parsed Response.
func (c *Client) AskJSON(ctx context.Context, prompt string) (*Response, error) {
	args := c.buildArgs(prompt, FormatJSON)
	cmd := c.newCmd(ctx, args)

	out, err := cmd.Output()
	if err != nil {
		return nil, wrapExecError(err)
	}
	var resp Response
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("claude: failed to parse JSON response: %w", err)
	}
	return &resp, nil
}

// AskWithSchema runs the prompt with a JSON schema constraint (--output-format json --output-schema).
func (c *Client) AskWithSchema(ctx context.Context, prompt string, schema string) (*Response, error) {
	args := c.buildArgs(prompt, FormatJSON, "--output-schema", schema)
	cmd := c.newCmd(ctx, args)

	out, err := cmd.Output()
	if err != nil {
		return nil, wrapExecError(err)
	}
	var resp Response
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("claude: failed to parse JSON response: %w", err)
	}
	return &resp, nil
}

// Resume continues a previous session identified by sessionID.
func (c *Client) Resume(ctx context.Context, sessionID string, prompt string) (*Response, error) {
	args := c.buildArgs(prompt, FormatJSON, "--resume", sessionID)
	cmd := c.newCmd(ctx, args)

	out, err := cmd.Output()
	if err != nil {
		return nil, wrapExecError(err)
	}
	var resp Response
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("claude: failed to parse JSON response: %w", err)
	}
	return &resp, nil
}

// Continue resumes the most recent session.
func (c *Client) Continue(ctx context.Context, prompt string) (*Response, error) {
	args := c.buildArgs(prompt, FormatJSON, "--continue")
	cmd := c.newCmd(ctx, args)

	out, err := cmd.Output()
	if err != nil {
		return nil, wrapExecError(err)
	}
	var resp Response
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("claude: failed to parse JSON response: %w", err)
	}
	return &resp, nil
}

// Pipe sends input from an io.Reader as stdin to the claude process alongside the prompt.
func (c *Client) Pipe(ctx context.Context, input io.Reader, prompt string) (string, error) {
	args := c.buildArgs(prompt, FormatText)
	cmd := c.newCmd(ctx, args)
	cmd.Stdin = input

	out, err := cmd.Output()
	if err != nil {
		return "", wrapExecError(err)
	}
	return string(bytes.TrimSpace(out)), nil
}

// wrapExecError extracts stderr from *exec.ExitError if available.
func wrapExecError(err error) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr := string(bytes.TrimSpace(exitErr.Stderr))
		if stderr != "" {
			return fmt.Errorf("claude: process exited with code %d: %s", exitErr.ExitCode(), stderr)
		}
		return fmt.Errorf("claude: process exited with code %d", exitErr.ExitCode())
	}
	return fmt.Errorf("claude: %w", err)
}

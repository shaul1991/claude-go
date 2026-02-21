package claude

import (
	"testing"
)

func TestNewClientDefaults(t *testing.T) {
	c := NewClient()
	if c.cliPath != "claude" {
		t.Errorf("expected default cliPath 'claude', got %q", c.cliPath)
	}
	if c.model != "" {
		t.Errorf("expected empty model, got %q", c.model)
	}
	if c.maxTurns != 0 {
		t.Errorf("expected maxTurns 0, got %d", c.maxTurns)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	c := NewClient(
		WithCLIPath("/usr/local/bin/claude"),
		WithModel("opus"),
		WithSystemPrompt("You are helpful"),
		WithAppendSystemPrompt("Be concise"),
		WithAllowedTools("bash", "read"),
		WithMaxTurns(5),
		WithMaxBudget(1.5),
		WithWorkDir("/tmp"),
	)

	if c.cliPath != "/usr/local/bin/claude" {
		t.Errorf("cliPath = %q", c.cliPath)
	}
	if c.model != "opus" {
		t.Errorf("model = %q", c.model)
	}
	if c.systemPrompt != "You are helpful" {
		t.Errorf("systemPrompt = %q", c.systemPrompt)
	}
	if c.appendPrompt != "Be concise" {
		t.Errorf("appendPrompt = %q", c.appendPrompt)
	}
	if len(c.allowedTools) != 2 || c.allowedTools[0] != "bash" || c.allowedTools[1] != "read" {
		t.Errorf("allowedTools = %v", c.allowedTools)
	}
	if c.maxTurns != 5 {
		t.Errorf("maxTurns = %d", c.maxTurns)
	}
	if c.maxBudget != 1.5 {
		t.Errorf("maxBudget = %f", c.maxBudget)
	}
	if c.workDir != "/tmp" {
		t.Errorf("workDir = %q", c.workDir)
	}
}

func TestBuildArgsMinimal(t *testing.T) {
	c := NewClient()
	args := c.buildArgs("hello", FormatText)

	expected := []string{"-p", "hello", "--output-format", "text"}
	assertArgs(t, expected, args)
}

func TestBuildArgsAllOptions(t *testing.T) {
	c := NewClient(
		WithModel("sonnet"),
		WithSystemPrompt("sys"),
		WithAppendSystemPrompt("append"),
		WithAllowedTools("bash", "read"),
		WithMaxTurns(3),
		WithMaxBudget(2.0),
	)
	args := c.buildArgs("hello", FormatJSON, "--resume", "sess123")

	expected := []string{
		"-p", "hello",
		"--output-format", "json",
		"--model", "sonnet",
		"--system-prompt", "sys",
		"--append-system-prompt", "append",
		"--allowedTools", "bash",
		"--allowedTools", "read",
		"--max-turns", "3",
		"--max-budget-usd", "2",
		"--resume", "sess123",
	}
	assertArgs(t, expected, args)
}

func TestBuildArgsStreamJSON(t *testing.T) {
	c := NewClient()
	args := c.buildArgs("test", FormatStreamJSON)

	expected := []string{"-p", "test", "--output-format", "stream-json", "--verbose", "--include-partial-messages"}
	assertArgs(t, expected, args)
}

func TestBuildArgsContinue(t *testing.T) {
	c := NewClient()
	args := c.buildArgs("more", FormatJSON, "--continue")

	expected := []string{"-p", "more", "--output-format", "json", "--continue"}
	assertArgs(t, expected, args)
}

func TestBuildArgsOutputSchema(t *testing.T) {
	c := NewClient()
	schema := `{"type":"object","properties":{"name":{"type":"string"}}}`
	args := c.buildArgs("test", FormatJSON, "--output-schema", schema)

	expected := []string{"-p", "test", "--output-format", "json", "--output-schema", schema}
	assertArgs(t, expected, args)
}

func assertArgs(t *testing.T, expected, actual []string) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("args length mismatch: expected %d, got %d\nexpected: %v\nactual:   %v", len(expected), len(actual), expected, actual)
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("args[%d]: expected %q, got %q", i, expected[i], actual[i])
		}
	}
}

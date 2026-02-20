package claude

// Option configures a Client.
type Option func(*Client)

// WithModel sets the --model flag.
func WithModel(model string) Option {
	return func(c *Client) {
		c.model = model
	}
}

// WithSystemPrompt sets the --system-prompt flag.
func WithSystemPrompt(prompt string) Option {
	return func(c *Client) {
		c.systemPrompt = prompt
	}
}

// WithAppendSystemPrompt sets the --append-system-prompt flag.
func WithAppendSystemPrompt(prompt string) Option {
	return func(c *Client) {
		c.appendPrompt = prompt
	}
}

// WithAllowedTools sets the --allowedTools flag.
func WithAllowedTools(tools ...string) Option {
	return func(c *Client) {
		c.allowedTools = tools
	}
}

// WithMaxTurns sets the --max-turns flag.
func WithMaxTurns(n int) Option {
	return func(c *Client) {
		c.maxTurns = n
	}
}

// WithMaxBudget sets the --max-budget-usd flag.
func WithMaxBudget(usd float64) Option {
	return func(c *Client) {
		c.maxBudget = usd
	}
}

// WithWorkDir sets the working directory for the claude process.
func WithWorkDir(dir string) Option {
	return func(c *Client) {
		c.workDir = dir
	}
}

// WithCLIPath overrides the default "claude" binary path.
func WithCLIPath(path string) Option {
	return func(c *Client) {
		c.cliPath = path
	}
}

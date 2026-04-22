package config

// Config represents the root configuration structure
type Config struct {
	Routes []Route `mapstructure:"routes"`
}

// Route defines a single webhook endpoint and its associated execution logic
type Route struct {
	Path                string   `mapstructure:"path"`
	Method              string   `mapstructure:"method"`
	APIKey              string   `mapstructure:"api_key"`
	GithubWebhookSecret string   `mapstructure:"github_webhook_secret"`
	Headers             []Header `mapstructure:"headers"`
	Rules               []Rule   `mapstructure:"rules"`
	Command             Command  `mapstructure:"command"`
}

// Header represents a key-value pair for HTTP header validation
type Header struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

// Rule defines a conditional check against the incoming webhook payload
type Rule struct {
	Field    string      `mapstructure:"field"`
	Operator string      `mapstructure:"operator"`
	Value    interface{} `mapstructure:"value"`
}

// Command describes the shell command to be executed upon a successful match
type Command struct {
	Execute string `mapstructure:"execute"`
	Async   bool   `mapstructure:"async"`
}

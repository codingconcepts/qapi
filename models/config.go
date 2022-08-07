package models

// Config holds the top-level configuration for the application.
type Config struct {
	Environment Environment       `yaml:"environment"`
	Variables   map[string]string `yaml:"variables"`
	Requests    []Request         `yaml:"requests"`
}

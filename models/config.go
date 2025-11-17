package models

import "time"

// Config holds the top-level configuration for the application.
type Config struct {
	Environment      Environment    `yaml:"environment"`
	Variables        map[string]any `yaml:"variables"`
	SetupRequests    []Request      `yaml:"setup_requests"`
	RequestFrequency time.Duration  `yaml:"request_frequency"`
	Requests         []Request      `yaml:"requests"`
}

package models

// Environment contains known, named values that configuration the application.
type Environment struct {
	BaseURL string `yaml:"base_url"`
}

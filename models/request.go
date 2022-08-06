package models

// Request holds the configuration that will be used when making a request
// to an API.
type Request struct {
	Name       string            `yaml:"name"`
	Headers    map[string]string `yaml:"headers"`
	Method     string            `yaml:"method"`
	Path       string            `yaml:"path"`
	Body       string            `yaml:"body"`
	Extractors []Extractor       `yaml:"extractors"`
}

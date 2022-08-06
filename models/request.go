package models

type Request struct {
	Name       string            `yaml:"name"`
	Headers    map[string]string `yaml:"headers"`
	Method     string            `yaml:"method"`
	Path       string            `yaml:"path"`
	Body       string            `yaml:"body"`
	Extractors []Extractor       `yaml:"extractors"`
}

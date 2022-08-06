package models

// Extractor allows values from responses to be captured and used in
// subsequent requests.
type Extractor struct {
	Type      string            `yaml:"type"`
	Selectors map[string]string `yaml:"selectors"`
}

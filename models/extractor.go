package models

type Extractor struct {
	Type      string            `yaml:"type"`
	Selectors map[string]string `yaml:"selectors"`
}

package models

// Assertion allows assertions to be performed against variables.
type Assertion struct {
	Type     string `json:"type"`
	Variable string `json:"variable"`
}

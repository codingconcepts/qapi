package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Validator func(key string, variables map[string]string) error

func fetchAndValidate(key string, variables map[string]string) (string, error) {
	value, ok := variables[key]
	if !ok {
		return "", fmt.Errorf("variable missing: %q", key)
	}

	return value, nil
}

func ValidateIsNotNull(key string, variables map[string]string) error {
	value, err := fetchAndValidate(key, variables)
	if err != nil {
		return err
	}

	if value == "" {
		return fmt.Errorf("variable null: %q", key)
	}

	return nil
}

func ValidateIsUUID(key string, variables map[string]string) error {
	value, err := fetchAndValidate(key, variables)
	if err != nil {
		return err
	}

	_, err = uuid.Parse(value)
	return err
}

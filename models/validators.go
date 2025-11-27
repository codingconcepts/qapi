package models

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type Validator func(key string, variables map[string]any) error

func fetchAndValidate(key string, variables map[string]any) (any, error) {
	value, ok := variables[key]
	if !ok {
		return "", fmt.Errorf("variable missing: %q", key)
	}

	return value, nil
}

func ValidateIsNotNull(key string, variables map[string]any) error {
	value, err := fetchAndValidate(key, variables)
	if err != nil {
		return err
	}

	if value == "" {
		return fmt.Errorf("variable null: %q", key)
	}

	return nil
}

func ValidateIsNotEmpty(key string, variables map[string]any) error {
	value, err := fetchAndValidate(key, variables)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)

	if val.Len() == 0 {
		return fmt.Errorf("variable empty: %q", key)
	}

	return nil
}

func ValidateIsUUID(key string, variables map[string]any) error {
	value, err := fetchAndValidate(key, variables)
	if err != nil {
		return err
	}

	valueString := fmt.Sprintf("%v", value)

	_, err = uuid.Parse(valueString)
	return err
}

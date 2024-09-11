package util

import "fmt"

func GetMapNullableValue[T any](m map[string]any, key string) *T {
	value, exists := m[key]
	if !exists {
		return nil
	}

	converted, ok := value.(T)
	if !ok {
		return nil
	}

	return &converted
}

func GetMapValue[T any](m map[string]any, key string, defaultValue T) T {
	value, exists := m[key]
	if !exists {
		return defaultValue
	}

	converted, ok := value.(T)
	if !ok {
		return defaultValue
	}

	return converted
}

func GetMapValueAsString(m map[string]any, key string, defaultValue string) string {
	value, exists := m[key]
	if !exists {
		return defaultValue
	}

	if value == nil {
		return defaultValue
	}

	return fmt.Sprintf("%s", value)
}

func GetMapValueString[T any](m map[string]any, key string, defaultValue T) string {
	value, exists := m[key]
	if !exists {
		value = defaultValue
	}

	converted, ok := value.(T)
	if !ok {
		value = defaultValue
	}

	return fmt.Sprintf("%v", converted)
}

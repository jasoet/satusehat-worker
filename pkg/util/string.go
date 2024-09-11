package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func IsSameType(err error, expected error) bool {
	sameType := fmt.Sprintf("%T", err) == fmt.Sprintf("%T", expected)
	return sameType
}

func StringNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

func StringNotNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IntToString(i *int) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%v", *i)
}

func JsonNumber(s string) *json.Number {
	s = strings.Replace(s, ",", ".", 1)
	s = strings.TrimRight(s, ".")
	if s == "" {
		return nil
	}

	value := json.Number(s)

	return &value
}

func StrPtr(s string) *string { return &s }

func StrPtrFmt(format string, a ...any) *string {
	s := fmt.Sprintf(format, a...)
	return &s
}

func NotEmpty(s *string) bool {
	if s == nil {
		return false
	}

	if strings.TrimSpace(*s) == "" {
		return false
	}
	return true
}

package util

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestStringNotEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"NotEmpty", "test string", true},
		{"JustSpaces", "    ", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringNotEmpty(tt.s); got != tt.want {
				t.Errorf("StringNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonNumber(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want *json.Number
	}{
		{"ValidNumber", "123", jsonNumber("123")},
		{"FloatNumber", "123.456", jsonNumber("123.456")},
		{"NegativeNumber", "-123", jsonNumber("-123")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonNumber(tt.s); *got != *tt.want {
				t.Errorf("JsonNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func jsonNumber(s string) *json.Number {
	value := json.Number(s)
	return &value
}

func TestStrPtr(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want *string
	}{
		{"ValidString", "test string", StrPtr("test string")},
		{"EmptyString", "", StrPtr("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrPtr(tt.s); *got != *tt.want {
				t.Errorf("StrPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrPtrFmt(t *testing.T) {
	tests := []struct {
		name   string
		format string
		a      []any
		want   *string
	}{
		{"OneParam", "test %s", []any{"123"}, StrPtr("test 123")},
		{"MultipleParams", "test %s %d", []any{"abc", 123}, StrPtr("test abc 123")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrPtrFmt(tt.format, tt.a...); *got != *tt.want {
				t.Errorf("StrPtrFmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameType(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected error
		want     bool
	}{
		{
			name:     "BothErrorsAreNil",
			err:      nil,
			expected: nil,
			want:     true,
		},
		{
			name:     "BothErrorsAreCustom",
			err:      errors.New("custom error"),
			expected: errors.New("another custom error"),
			want:     true,
		},
		{
			name:     "ErrorAndNil",
			err:      errors.New("custom error"),
			expected: nil,
			want:     false,
		},
		{
			name:     "NilAndError",
			err:      nil,
			expected: errors.New("custom error"),
			want:     false,
		},
		{
			name:     "DifferentTypes",
			err:      errors.New("custom error"),
			expected: &MyError{},
			want:     false,
		},
		{
			name:     "SampeType",
			err:      &MyError{},
			expected: &MyError{},
			want:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsSameType(tc.err, tc.expected)
			if got != tc.want {
				t.Errorf("IsSameType(%v, %v) = %v; want %v", tc.err, tc.expected, got, tc.want)
			}
		})
	}
}

type MyError struct{}

func (e *MyError) Error() string {
	return "my error"
}

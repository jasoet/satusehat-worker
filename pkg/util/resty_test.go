package util

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"testing"
)

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name string
		res  *resty.Response
		want bool
	}{
		{"SuccessStatus", &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusOK}}, false},
		{"ServerErrorStatus", &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusInternalServerError}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServerError(tt.res); got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		res  *resty.Response
		want bool
	}{
		{"SuccessStatus", &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusOK}}, false},
		{"UnauthorizedStatus", &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusUnauthorized}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.res); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

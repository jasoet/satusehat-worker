package util

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

func IsServerError(response *resty.Response) bool {
	return response.StatusCode() >= 500
}

func IsUnauthorized(response *resty.Response) bool {
	return response.StatusCode() == http.StatusUnauthorized
}

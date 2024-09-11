package satusehat

import "fmt"

type UnauthorizedError struct {
	StatusCode int
	Msg        string
	RespBody   string
}

func (e *UnauthorizedError) Error() string { return e.Msg }
func NewUnauthorizedError(statusCode int, msg string, respBody string) *UnauthorizedError {
	return &UnauthorizedError{
		StatusCode: statusCode,
		Msg:        msg,
		RespBody:   respBody,
	}
}

type ExecutionError struct {
	Msg string
	Err error
}

func (e *ExecutionError) Error() string { return e.Msg }
func (e *ExecutionError) Unwrap() error { return e.Err }
func NewExecutionError(msg string, err error) *ExecutionError {
	return &ExecutionError{
		Msg: msg,
		Err: err,
	}
}

type ServerError struct {
	StatusCode int
	Msg        string
	RespBody   string
}

func (e *ServerError) Error() string { return fmt.Sprintf("%s: %s", e.Msg, e.RespBody) }
func NewServerError(statusCode int, msg string, respBody string) *ServerError {
	return &ServerError{
		StatusCode: statusCode,
		Msg:        msg,
		RespBody:   respBody,
	}
}

type ResponseError struct {
	StatusCode int
	Msg        string
	RespBody   string
}

func (e *ResponseError) Error() string { return fmt.Sprintf("%s: %s", e.Msg, e.RespBody) }
func NewResponseError(statusCode int, msg string, respBody string) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		Msg:        msg,
		RespBody:   respBody,
	}
}

type ResourceNotFoundError struct {
	StatusCode int
	Msg        string
	RespBody   string
}

func (e *ResourceNotFoundError) Error() string { return fmt.Sprintf("%s: %s", e.Msg, e.RespBody) }
func NewResourceNotFoundError(statusCode int, msg string, respBody string) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		StatusCode: statusCode,
		Msg:        msg,
		RespBody:   respBody,
	}
}

package internal

import "fmt"

// ErrorType represents the type of error
type ErrorType string

const (
	ValidationError ErrorType = "ValidationError"
	APIError        ErrorType = "APIError"
	TimeoutError    ErrorType = "TimeoutError"
	ParseError      ErrorType = "ParseError"
	NotFoundError   ErrorType = "NotFoundError"
	ServerError     ErrorType = "ServerError"
)

// AppError represents an application error with type and code
type AppError struct {
	Type    ErrorType
	Message string
	Code    string
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: message,
		Code:    "INVALID_INPUT",
	}
}

// NewAPIError creates an API error
func NewAPIError(statusCode int, message string) *AppError {
	return &AppError{
		Type:    APIError,
		Message: fmt.Sprintf("API returned %d: %s", statusCode, message),
		Code:    fmt.Sprintf("API_%d", statusCode),
	}
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(message string) *AppError {
	return &AppError{
		Type:    TimeoutError,
		Message: message,
		Code:    "TIMEOUT",
	}
}

// NewParseError creates a parse error
func NewParseError(message string) *AppError {
	return &AppError{
		Type:    ParseError,
		Message: message,
		Code:    "PARSE_ERROR",
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Message: message,
		Code:    "NOT_FOUND",
	}
}

// NewServerError creates a server error
func NewServerError(message string) *AppError {
	return &AppError{
		Type:    ServerError,
		Message: message,
		Code:    "SERVER_ERROR",
	}
}

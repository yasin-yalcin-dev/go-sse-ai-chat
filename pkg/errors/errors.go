/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// AppError is the application-specific error type
type AppError struct {
	Err        error
	Message    string
	StatusCode int
	Code       string
	Stack      string
	Context    map[string]interface{}
}

// Standard error codes
const (
	CodeBadRequest           = "BAD_REQUEST"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeForbidden            = "FORBIDDEN"
	CodeNotFound             = "NOT_FOUND"
	CodeConflict             = "CONFLICT"
	CodeInternalServerError  = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeTimeout              = "TIMEOUT"
	CodeValidationError      = "VALIDATION_ERROR"
	CodeDatabaseError        = "DATABASE_ERROR"
	CodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
)

// New creates a new error with a message
func New(message string) error {
	return errors.New(message)
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Err != nil {
		return e.Err.Error()
	}

	return "unknown error"
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithStack adds stack trace information to the error
func (e *AppError) WithStack() *AppError {
	if e.Stack != "" {
		return e
	}

	// Capture stack trace
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stackBuilder strings.Builder
	for {
		frame, more := frames.Next()
		// Skip stdlib and runtime frames
		if !strings.Contains(frame.File, "runtime/") {
			fmt.Fprintf(&stackBuilder, "%s:%d - %s\n", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}

	e.Stack = stackBuilder.String()
	return e
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithStatusCode adds HTTP status code to the error
func (e *AppError) WithStatusCode(statusCode int) *AppError {
	e.StatusCode = statusCode
	return e
}

// WithCode adds an error code to the error
func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

// GetStatusCode returns the HTTP status code for the error
func (e *AppError) GetStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}

	// Default error mapping
	switch e.Code {
	case CodeBadRequest, CodeValidationError:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case CodeTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// Format provides custom formatting for %+v
func (e *AppError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// Detailed error with stack trace
			fmt.Fprintf(s, "%s\n", e.Error())
			if e.Code != "" {
				fmt.Fprintf(s, "Code: %s\n", e.Code)
			}
			if e.StatusCode != 0 {
				fmt.Fprintf(s, "Status: %d\n", e.StatusCode)
			}
			if len(e.Context) > 0 {
				fmt.Fprintf(s, "Context: %+v\n", e.Context)
			}
			if e.Stack != "" {
				fmt.Fprintf(s, "Stack:\n%s\n", e.Stack)
			}
			if e.Err != nil {
				fmt.Fprintf(s, "Cause: %+v\n", e.Err)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", e.Error())
	}
}

// ToResponse creates a standardized error response for API
func (e *AppError) ToResponse() map[string]interface{} {
	resp := map[string]interface{}{
		"error": map[string]interface{}{
			"message": e.Error(),
		},
	}

	if e.Code != "" {
		resp["error"].(map[string]interface{})["code"] = e.Code
	}

	// Only include detailed information in non-production environments
	// This would typically be controlled by a configuration flag
	if len(e.Context) > 0 {
		resp["error"].(map[string]interface{})["context"] = e.Context
	}

	return resp
}

// NewAppError creates a new application error
func NewAppError(message string, err error) *AppError {
	appErr := &AppError{
		Message: message,
		Err:     err,
	}
	return appErr.WithStack()
}

// Wrap wraps an existing error with a message
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	appErr := &AppError{
		Message: message,
		Err:     err,
	}
	return appErr.WithStack()
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeBadRequest).
		WithStatusCode(http.StatusBadRequest)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeUnauthorized).
		WithStatusCode(http.StatusUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeForbidden).
		WithStatusCode(http.StatusForbidden)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeNotFound).
		WithStatusCode(http.StatusNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeConflict).
		WithStatusCode(http.StatusConflict)
}

// NewInternalError creates an internal server error
func NewInternalError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeInternalServerError).
		WithStatusCode(http.StatusInternalServerError)
}

// NewServiceUnavailableError creates a service unavailable error
func NewServiceUnavailableError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeServiceUnavailable).
		WithStatusCode(http.StatusServiceUnavailable)
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeTimeout).
		WithStatusCode(http.StatusGatewayTimeout)
}

// NewValidationError creates a validation error
func NewValidationError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeValidationError).
		WithStatusCode(http.StatusBadRequest)
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeDatabaseError).
		WithStatusCode(http.StatusInternalServerError)
}

// NewExternalServiceError creates an external service error
func NewExternalServiceError(message string, err error) *AppError {
	return NewAppError(message, err).
		WithCode(CodeExternalServiceError).
		WithStatusCode(http.StatusBadGateway)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// AsAppError tries to convert an error to AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

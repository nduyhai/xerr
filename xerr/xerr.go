// Package xerr provides structured error handling with support for gRPC and HTTP.
package xerr

import (
	"fmt"
	"google.golang.org/grpc/codes"
)

// StructuredError represents a rich error with code, message, and metadata.
// It can be converted to/from gRPC status and HTTP responses.
type StructuredError struct {
	Code     string            // Machine-readable app error code, e.g. "INVALID_AMOUNT"
	Message  string            // Developer-facing message
	Reason   string            // Optional user-facing message (for UI/i18n)
	GRPCCode codes.Code        // gRPC status code
	HTTPCode int               // HTTP status code
	Metadata map[string]string // Optional context (trace ID, field, etc.)
}

// Error implements the error interface.
func (e *StructuredError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New creates a new StructuredError with the given code and message.
func New(code string, message string) *StructuredError {
	return &StructuredError{
		Code:     code,
		Message:  message,
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
	}
}

// WithReason adds a user-facing reason to the error.
func (e *StructuredError) WithReason(reason string) *StructuredError {
	e.Reason = reason
	return e
}

// WithGRPCCode sets the gRPC status code.
func (e *StructuredError) WithGRPCCode(code codes.Code) *StructuredError {
	e.GRPCCode = code
	return e
}

// WithHTTPCode sets the HTTP status code.
func (e *StructuredError) WithHTTPCode(code int) *StructuredError {
	e.HTTPCode = code
	return e
}

// WithMetadata adds metadata to the error.
func (e *StructuredError) WithMetadata(key string, value string) *StructuredError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// Wrap wraps an existing error with a structured error.
func Wrap(err error, code string) *StructuredError {
	if err == nil {
		return nil
	}

	// If it's already a StructuredError, just update the code
	if se, ok := err.(*StructuredError); ok {
		se.Code = code
		return se
	}

	return &StructuredError{
		Code:     code,
		Message:  err.Error(),
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
	}
}

// Is The implements the errors.Is interface for error comparison.
func (e *StructuredError) Is(target error) bool {
	if se, ok := target.(*StructuredError); ok {
		return e.Code == se.Code
	}
	return false
}

// NewWithHTTPAndGRPC creates a new StructuredError with the given code, message, HTTP code, and gRPC code.
func NewWithHTTPAndGRPC(code string, message string, httpCode int, grpcCode codes.Code) *StructuredError {
	return &StructuredError{
		Code:     code,
		Message:  message,
		GRPCCode: grpcCode,
		HTTPCode: httpCode,
	}
}

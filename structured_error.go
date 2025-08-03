// Package xerr provides structured error handling with support for gRPC and HTTP.
package xerr

import (
	"fmt"
	"google.golang.org/grpc/codes"
)

// StructuredError represents a rich error with code, message, and metadata.
// It implements the Error interface and can be converted to/from gRPC status and HTTP responses.
// This is the concrete implementation that is returned by the factory functions.
type StructuredError struct {
	reason   Reason            // Reason interface implementation
	GRPCCode codes.Code        // gRPC status code
	HTTPCode int               // HTTP status code
	Metadata map[string]string // Optional context (trace ID, field, etc.)
	Domain   string            // Domain for gRPC ErrorInfo
	Cause    error             // Original error that caused this error
}

// Accessor methods for StructuredError

// GetReason returns the Reason interface.
func (e *StructuredError) GetReason() Reason {
	return e.reason
}

// GetCode returns the error code from the Reason.
func (e *StructuredError) GetCode() string {
	if e.reason == nil {
		return ""
	}
	return e.reason.Code()
}

// GetMessage returns the error message from the Reason.
func (e *StructuredError) GetMessage() string {
	if e.reason == nil {
		return ""
	}
	return e.reason.Message()
}

// GetUserReason returns the user-facing reason from the Reason.
func (e *StructuredError) GetUserReason() string {
	if e.reason == nil {
		return ""
	}
	return e.reason.Reason()
}

// GetGRPCCode returns the gRPC status code.
func (e *StructuredError) GetGRPCCode() codes.Code {
	return e.GRPCCode
}

// GetHTTPCode returns the HTTP status code.
func (e *StructuredError) GetHTTPCode() int {
	return e.HTTPCode
}

// GetMetadata returns the error metadata.
func (e *StructuredError) GetMetadata() map[string]string {
	return e.Metadata
}

// GetCause returns the underlying cause of the error.
func (e *StructuredError) GetCause() error {
	return e.Cause
}

// Error implements the error interface.
func (e *StructuredError) Error() string {
	if e.reason == nil {
		return "unknown error"
	}
	return fmt.Sprintf("[%s] %s", e.reason.Code(), e.reason.Message())
}

// New creates a new Error with the given code and message.
// It returns an Error interface that can be used with all the methods defined in the interface.
func New(code string, message string) Error {
	return NewWithHTTPAndGRPC(code, message, 500, codes.Unknown)
}

// WithReason adds a user-facing reason to the error.
func (e *StructuredError) WithReason(reason string) Error {
	if defaultReason, ok := e.reason.(*DefaultReason); ok {
		defaultReason.WithReason(reason)
	} else {
		// If the reason is not a DefaultReason, create a new one with the same code and message
		newReason := NewDefaultReason(e.GetCode(), e.GetMessage()).WithReason(reason)
		e.reason = newReason
	}
	return e
}

// WithCustomReason sets a custom implementation of the Reason interface.
// This allows for more flexible error reason handling.
func (e *StructuredError) WithCustomReason(reason Reason) Error {
	e.reason = reason
	return e
}

// WithGRPCCode sets the gRPC status code.
func (e *StructuredError) WithGRPCCode(code codes.Code) Error {
	e.GRPCCode = code
	return e
}

// WithHTTPCode sets the HTTP status code.
func (e *StructuredError) WithHTTPCode(code int) Error {
	e.HTTPCode = code
	return e
}

// WithMetadata adds metadata to the error.
func (e *StructuredError) WithMetadata(key string, value string) Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// Is implements the errors.Is interface for error comparison.
func (e *StructuredError) Is(target error) bool {
	if se, ok := target.(*StructuredError); ok {
		return e.GetCode() == se.GetCode()
	}
	return false
}

// Unwrap implements the errors.Unwrap interface to return the underlying cause.
func (e *StructuredError) Unwrap() error {
	return e.Cause
}

// NewWithHTTPAndGRPC creates a new Error with the given code, message, HTTP code, and gRPC code.
// It returns an Error interface that can be used with all the methods defined in the interface.
func NewWithHTTPAndGRPC(code string, message string, httpCode int, grpcCode codes.Code) Error {
	return &StructuredError{
		reason:   NewDefaultReason(code, message),
		GRPCCode: grpcCode,
		HTTPCode: httpCode,
	}
}

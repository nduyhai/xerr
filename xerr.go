// Package xerr provides structured error handling with support for gRPC and HTTP.
package xerr

import (
	"fmt"
	"google.golang.org/grpc/codes"
)

// Error is the interface that wraps the basic error functionality.
// It extends the standard error interface with additional methods for
// structured error handling. This interface allows for loose coupling
// between error producers and consumers, making the code more testable
// and maintainable.
type Error interface {
	// Error returns the error message.
	error
	
	// Core accessor methods
	
	// GetCode returns the error code.
	GetCode() string
	
	// GetMessage returns the error message.
	GetMessage() string
	
	// GetReason returns the user-facing reason.
	GetReason() string
	
	// GetGRPCCode returns the gRPC status code.
	GetGRPCCode() codes.Code
	
	// GetHTTPCode returns the HTTP status code.
	GetHTTPCode() int
	
	// GetMetadata returns the error metadata.
	GetMetadata() map[string]string
	
	// GetCause returns the underlying cause of the error.
	GetCause() error
	
	// Core modifier methods
	
	// WithReason adds a user-facing reason to the error.
	WithReason(reason string) Error
	
	// WithGRPCCode sets the gRPC status code.
	WithGRPCCode(code codes.Code) Error
	
	// WithHTTPCode sets the HTTP status code.
	WithHTTPCode(code int) Error
	
	// WithMetadata adds metadata to the error.
	WithMetadata(key string, value string) Error
	
	// Standard error interface methods
	
	// Is implements the errors.Is interface for error comparison.
	Is(target error) bool
	
	// Unwrap implements the errors.Unwrap interface to return the underlying cause.
	Unwrap() error
}

// StructuredError represents a rich error with code, message, and metadata.
// It implements the Error interface and can be converted to/from gRPC status and HTTP responses.
// This is the concrete implementation that is returned by the factory functions.
type StructuredError struct {
	Code     string            // Machine-readable app error code, e.g. "INVALID_AMOUNT"
	Message  string            // Developer-facing message
	Reason   string            // Optional user-facing message (for UI/i18n)
	GRPCCode codes.Code        // gRPC status code
	HTTPCode int               // HTTP status code
	Metadata map[string]string // Optional context (trace ID, field, etc.)
	Cause    error             // Original error that caused this error
}

// Accessor methods for StructuredError

// GetCode returns the error code.
func (e *StructuredError) GetCode() string {
	return e.Code
}

// GetMessage returns the error message.
func (e *StructuredError) GetMessage() string {
	return e.Message
}

// GetReason returns the user-facing reason.
func (e *StructuredError) GetReason() string {
	return e.Reason
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
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New creates a new Error with the given code and message.
// It returns an Error interface that can be used with all the methods defined in the interface.
func New(code string, message string) Error {
	return &StructuredError{
		Code:     code,
		Message:  message,
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
	}
}

// WithReason adds a user-facing reason to the error.
func (e *StructuredError) WithReason(reason string) Error {
	e.Reason = reason
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
		return e.Code == se.Code
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
		Code:     code,
		Message:  message,
		GRPCCode: grpcCode,
		HTTPCode: httpCode,
	}
}

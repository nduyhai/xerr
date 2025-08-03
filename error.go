package xerr

import "google.golang.org/grpc/codes"

// Error is the interface that wraps the basic error functionality.
// It extends the standard error interface with additional methods for
// structured error handling. This interface allows for loose coupling
// between error producers and consumers, making the code more testable
// and maintainable.
type Error interface {
	// Error returns the error message.
	error

	// Core accessor methods

	// GetReason returns the Reason interface.
	GetReason() Reason

	// GetGRPCCode returns the gRPC status code.
	GetGRPCCode() codes.Code

	// GetHTTPCode returns the HTTP status code.
	GetHTTPCode() int

	// GetMetadata returns the error metadata.
	GetMetadata() map[string]string

	// GetCause returns the underlying cause of the error.
	GetCause() error

	// Convenience accessor methods

	// GetCode returns the error code from the Reason.
	GetCode() string

	// GetMessage returns the error message from the Reason.
	GetMessage() string

	// GetUserReason returns the user-facing reason from the Reason.
	GetUserReason() string

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

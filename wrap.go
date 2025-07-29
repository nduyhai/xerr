package xerr

import (
	"errors"
	"google.golang.org/grpc/codes"
)

// Wrap wraps an existing error with a structured error.
// It returns an Error interface that can be used with all the methods defined in the interface.
func Wrap(err error, code string) Error {
	if err == nil {
		return nil
	}

	// If it's already a StructuredError, just update the reason with the new code
	var se *StructuredError
	if errors.As(err, &se) {
		// Create a new DefaultReason with the new code and the existing message
		se.reason = NewDefaultReason(code, se.GetMessage())
		return se
	}
	return &StructuredError{
		reason:   NewDefaultReason(code, err.Error()),
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
		Cause:    err,
	}
}

// WrapDefault wraps an existing error with a structured error using the default UNKNOWN code.
// It returns an Error interface that can be used with all the methods defined in the interface.
func WrapDefault(err error) Error {
	return Wrap(err, UNKNOWN)
}

package xerr

import (
	"errors"
	"google.golang.org/grpc/codes"
)

// WrapWithReason wraps an existing error with a structured error using the provided Reason.
// It returns an Error interface that can be used with all the methods defined in the interface.
func WrapWithReason(err error, reason Reason) Error {
	if err == nil {
		return nil
	}

	// If it's already a StructuredError, just update the reason
	var se *StructuredError
	if errors.As(err, &se) {
		se.reason = reason
		return se
	}
	return &StructuredError{
		reason:   reason,
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
		Cause:    err,
	}
}

// WrapDefault wraps an existing error with a structured error using the default UNKNOWN code.
// It returns an Error interface that can be used with all the methods defined in the interface.
func WrapDefault(err error) Error {
	if err == nil {
		return nil
	}
	reason := NewDefaultReason("UNKNOWN", err.Error())
	return WrapWithReason(err, reason)
}

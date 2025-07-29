package xerr

import (
	"errors"
	"google.golang.org/grpc/codes"
)

// Wrap wraps an existing error with a structured error.
func Wrap(err error, code string) *StructuredError {
	if err == nil {
		return nil
	}

	// If it's already a StructuredError, just update the code
	var se *StructuredError
	if errors.As(err, &se) {
		se.Code = code
		return se
	}
	return &StructuredError{
		Code:     code,
		Message:  err.Error(),
		GRPCCode: codes.Unknown,
		HTTPCode: 500,
		Cause:    err,
	}
}

// WrapDefault wraps an existing error with a structured error using the default UNKNOWN code.
func WrapDefault(err error) *StructuredError {
	return Wrap(err, UNKNOWN)
}

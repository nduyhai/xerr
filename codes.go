package xerr

import (
	"google.golang.org/grpc/codes"
)

// Standard error codes for application errors.
// These codes are designed to be used with StructuredError.
const (
	// General errors
	UNKNOWN     = "UNKNOWN"     // Unknown error
	INTERNAL    = "INTERNAL"    // Internal server error
	UNAVAILABLE = "UNAVAILABLE" // Service unavailable
	TIMEOUT     = "TIMEOUT"     // Request timeout
	CANCELLED   = "CANCELLED"   // Request cancelled

	// Client errors
	INVALID_ARGUMENT    = "INVALID_ARGUMENT"    // Invalid argument
	FAILED_PRECONDITION = "FAILED_PRECONDITION" // Failed precondition
	OUT_OF_RANGE        = "OUT_OF_RANGE"        // Value out of range
	UNAUTHENTICATED     = "UNAUTHENTICATED"     // Unauthenticated request
	PERMISSION_DENIED   = "PERMISSION_DENIED"   // Permission denied
	NOT_FOUND           = "NOT_FOUND"           // Resource not found
	ALREADY_EXISTS      = "ALREADY_EXISTS"      // Resource already exists
	RESOURCE_EXHAUSTED  = "RESOURCE_EXHAUSTED"  // Resource quota exceeded
	ABORTED             = "ABORTED"             // Operation aborted

	// Data errors
	DATA_LOSS       = "DATA_LOSS"       // Unrecoverable data loss or corruption
	DATA_VALIDATION = "DATA_VALIDATION" // Data validation error

	// Business logic errors
	BUSINESS_RULE = "BUSINESS_RULE" // Business rule violation
	CONFLICT      = "CONFLICT"      // Conflict with current state
)

// StandardErrorMapping maps standard error codes to their corresponding gRPC and HTTP codes.
// This can be used to create errors with consistent codes across different protocols.
var StandardErrorMapping = map[string]struct {
	GRPCCode codes.Code
	HTTPCode int
}{
	UNKNOWN:     {GRPCCode: codes.Unknown, HTTPCode: 500},          // HTTP 500 Internal Server Error
	INTERNAL:    {GRPCCode: codes.Internal, HTTPCode: 500},         // HTTP 500 Internal Server Error
	UNAVAILABLE: {GRPCCode: codes.Unavailable, HTTPCode: 503},      // HTTP 503 Service Unavailable
	TIMEOUT:     {GRPCCode: codes.DeadlineExceeded, HTTPCode: 504}, // HTTP 504 Gateway Timeout
	CANCELLED:   {GRPCCode: codes.Canceled, HTTPCode: 499},         // HTTP 499 Client Closed Request

	INVALID_ARGUMENT:    {GRPCCode: codes.InvalidArgument, HTTPCode: 400},    // HTTP 400 Bad Request
	FAILED_PRECONDITION: {GRPCCode: codes.FailedPrecondition, HTTPCode: 400}, // HTTP 400 Bad Request
	OUT_OF_RANGE:        {GRPCCode: codes.OutOfRange, HTTPCode: 400},         // HTTP 400 Bad Request
	UNAUTHENTICATED:     {GRPCCode: codes.Unauthenticated, HTTPCode: 401},    // HTTP 401 Unauthorized
	PERMISSION_DENIED:   {GRPCCode: codes.PermissionDenied, HTTPCode: 403},   // HTTP 403 Forbidden
	NOT_FOUND:           {GRPCCode: codes.NotFound, HTTPCode: 404},           // HTTP 404 Not Found
	ALREADY_EXISTS:      {GRPCCode: codes.AlreadyExists, HTTPCode: 409},      // HTTP 409 Conflict
	RESOURCE_EXHAUSTED:  {GRPCCode: codes.ResourceExhausted, HTTPCode: 429},  // HTTP 429 Too Many Requests
	ABORTED:             {GRPCCode: codes.Aborted, HTTPCode: 409},            // HTTP 409 Conflict

	DATA_LOSS:       {GRPCCode: codes.DataLoss, HTTPCode: 500},        // HTTP 500 Internal Server Error
	DATA_VALIDATION: {GRPCCode: codes.InvalidArgument, HTTPCode: 422}, // HTTP 422 Unprocessable Entity

	BUSINESS_RULE: {GRPCCode: codes.FailedPrecondition, HTTPCode: 422}, // HTTP 422 Unprocessable Entity
	CONFLICT:      {GRPCCode: codes.Aborted, HTTPCode: 409},            // HTTP 409 Conflict
}

// NewStandardError creates a new Error with standard error code mapping.
// It automatically sets the appropriate gRPC and HTTP codes based on the error code.
// It returns an Error interface that can be used with all the methods defined in the interface.
func NewStandardError(code string, message string) Error {
	mapping, exists := StandardErrorMapping[code]
	if !exists {
		// Default to UNKNOWN if the code is not recognized
		mapping = StandardErrorMapping[UNKNOWN]
	}

	return &StructuredError{
		Code:     code,
		Message:  message,
		GRPCCode: mapping.GRPCCode,
		HTTPCode: mapping.HTTPCode,
	}
}

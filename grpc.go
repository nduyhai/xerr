package xerr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCStatus converts a StructuredError to a gRPC status.Status.
// It includes error details if available.
func (e *StructuredError) ToGRPCStatus() *status.Status {
	st := status.New(e.GRPCCode, e.GetMessage())

	// If we have additional details, add them to the status
	if len(e.Metadata) > 0 {
		errorInfo := &errdetails.ErrorInfo{
			Reason:   e.GetCode(),
			Domain:   "github.com/nduyhai/xerr",
			Metadata: e.Metadata,
		}

		// Add ErrorInfo with metadata
		var err error
		st, err = st.WithDetails(errorInfo)
		if err != nil {
			// If we can't add details, just return the status without details
			return st
		}
	}

	// Add localized message if available
	userReason := e.GetUserReason()
	if userReason != "" {
		localizedMsg := &errdetails.LocalizedMessage{
			Locale:  "en-US",
			Message: userReason,
		}

		// Add localized message
		st, _ = st.WithDetails(localizedMsg)
	}

	return st
}

// FromGRPCStatus converts a gRPC status.Status to an Error.
// It extracts error details if available and returns an Error interface
// that can be used with all the methods defined in the interface.
func FromGRPCStatus(st *status.Status) Error {
	if st == nil {
		return nil
	}

	// Default values
	code := "UNKNOWN"
	message := st.Message()
	userReason := ""
	metadata := make(map[string]string)

	// Extract details from the status
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			// Use the reason as the error code
			code = d.Reason

			// Copy metadata
			for k, v := range d.Metadata {
				metadata[k] = v
			}

		case *errdetails.LocalizedMessage:
			// Use the localized message as the user reason
			userReason = d.Message
		}
	}

	// Create the error with the extracted information
	reason := NewDefaultReason(code, message)
	if userReason != "" {
		reason.WithReason(userReason)
	}

	return &StructuredError{
		reason:   reason,
		GRPCCode: st.Code(),
		HTTPCode: grpcToHTTPCode(st.Code()),
		Metadata: metadata,
	}
}

// ToGRPCStatusProto converts a StructuredError to a google.rpc.Status proto.
func (e *StructuredError) ToGRPCStatusProto() *status.Status {
	return e.ToGRPCStatus()
}

// FromGRPCStatusProto converts a google.rpc.Status proto to an Error.
// It returns an Error interface that can be used with all the methods defined in the interface.
func FromGRPCStatusProto(st *status.Status) Error {
	return FromGRPCStatus(st)
}

// grpcToHTTPCode maps gRPC status codes to HTTP status codes.
func grpcToHTTPCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return 200
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return 500
	case codes.InvalidArgument:
		return 400
	case codes.DeadlineExceeded:
		return 504
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.PermissionDenied:
		return 403
	case codes.ResourceExhausted:
		return 429
	case codes.FailedPrecondition:
		return 400
	case codes.Aborted:
		return 409
	case codes.OutOfRange:
		return 400
	case codes.Unimplemented:
		return 501
	case codes.Internal:
		return 500
	case codes.Unavailable:
		return 503
	case codes.DataLoss:
		return 500
	case codes.Unauthenticated:
		return 401
	default:
		return 500
	}
}

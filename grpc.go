package xerr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCStatus converts a StructuredError to a gRPC status.Status.
// It includes error details if available.
func (e *StructuredError) ToGRPCStatus() *status.Status {
	st := status.New(e.GRPCCode, e.Message)

	// If we have additional details, add them to the status
	if len(e.Metadata) > 0 {
		errorInfo := &errdetails.ErrorInfo{
			Reason:   e.Code,
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
	if e.Reason != "" {
		localizedMsg := &errdetails.LocalizedMessage{
			Locale:  "en-US",
			Message: e.Reason,
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

	e := &StructuredError{
		Code:     "UNKNOWN", // Default code
		Message:  st.Message(),
		GRPCCode: st.Code(),
		HTTPCode: grpcToHTTPCode(st.Code()),
		Metadata: make(map[string]string),
	}

	// Extract details from the status
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			// Use the reason as the error code
			e.Code = d.Reason

			// Copy metadata
			for k, v := range d.Metadata {
				e.Metadata[k] = v
			}

		case *errdetails.LocalizedMessage:
			// Use the localized message as the reason
			e.Reason = d.Message
		}
	}

	return e
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

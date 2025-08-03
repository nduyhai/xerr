package xerr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	_ "google.golang.org/grpc/codes" // Used for GRPCCode field type (codes.Code)
	"google.golang.org/grpc/status"
)

// ToGRPCStatus converts a StructuredError to a gRPC status.Status.
// It includes error details if available.
func (e *StructuredError) ToGRPCStatus() *status.Status {
	st := status.New(e.GRPCCode, e.GetMessage())

	// If we have additional details, add them to the status
	if len(e.Metadata) > 0 || e.Domain != "" {
		domain := e.Domain
		if domain == "" {
			domain = "github.com/nduyhai/xerr"
		}
		errorInfo := &errdetails.ErrorInfo{
			Reason:   e.GetCode(),
			Domain:   domain,
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
	domain := ""
	metadata := make(map[string]string)

	// Extract details from the status
	for _, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			// Use the reason as the error code
			code = d.Reason
			domain = d.Domain

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
		HTTPCode: DefaultConverter.GRPCToHTTP(st.Code()),
		Metadata: metadata,
		Domain:   domain,
	}
}


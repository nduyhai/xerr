package xerr

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
)

// HTTPError represents the JSON structure for HTTP error responses.
type HTTPError struct {
	Code     string            `json:"code"`               // Machine-readable error code
	Message  string            `json:"message"`            // Developer-facing error message
	Reason   string            `json:"reason,omitempty"`   // User-facing error message
	Metadata map[string]string `json:"metadata,omitempty"` // Additional error context
}

// ToHTTP converts a StructuredError to an HTTP response.
// It writes the error as JSON to the http.ResponseWriter with the appropriate status code.
func (e *StructuredError) ToHTTP(w http.ResponseWriter) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(e.HTTPCode)

	// Create HTTP error response
	httpErr := HTTPError{
		Code:     e.Code,
		Message:  e.Message,
		Reason:   e.Reason,
		Metadata: e.Metadata,
	}

	// Write JSON response
	_ = json.NewEncoder(w).Encode(httpErr)
}

// ToHTTPJSON converts a StructuredError to an HTTP JSON error response.
// It returns the JSON bytes and the HTTP status code.
func (e *StructuredError) ToHTTPJSON() ([]byte, int) {
	httpErr := HTTPError{
		Code:     e.Code,
		Message:  e.Message,
		Reason:   e.Reason,
		Metadata: e.Metadata,
	}

	jsonBytes, _ := json.Marshal(httpErr)
	return jsonBytes, e.HTTPCode
}

// FromHTTPJSON converts an HTTP JSON error response to a StructuredError.
func FromHTTPJSON(jsonBytes []byte, statusCode int) (*StructuredError, error) {
	var httpErr HTTPError
	if err := json.Unmarshal(jsonBytes, &httpErr); err != nil {
		return nil, err
	}

	return &StructuredError{
		Code:     httpErr.Code,
		Message:  httpErr.Message,
		Reason:   httpErr.Reason,
		GRPCCode: httpToGRPCCode(statusCode),
		HTTPCode: statusCode,
		Metadata: httpErr.Metadata,
	}, nil
}

// httpToGRPCCode maps HTTP status codes to gRPC status codes.
func httpToGRPCCode(httpCode int) codes.Code {
	switch {
	case httpCode >= 200 && httpCode < 300:
		return codes.OK
	case httpCode == 400:
		return codes.InvalidArgument
	case httpCode == 401:
		return codes.Unauthenticated
	case httpCode == 403:
		return codes.PermissionDenied
	case httpCode == 404:
		return codes.NotFound
	case httpCode == 409:
		return codes.Aborted
	case httpCode == 422:
		return codes.FailedPrecondition
	case httpCode == 429:
		return codes.ResourceExhausted
	case httpCode == 499:
		return codes.Canceled
	case httpCode == 500:
		return codes.Internal
	case httpCode == 501:
		return codes.Unimplemented
	case httpCode == 503:
		return codes.Unavailable
	case httpCode == 504:
		return codes.DeadlineExceeded
	default:
		if httpCode >= 400 && httpCode < 500 {
			return codes.InvalidArgument
		}
		return codes.Unknown
	}
}

// WriteHTTPError writes a structured error to an HTTP response.
// This is a convenience function for creating and writing an error in one step.
func WriteHTTPError(w http.ResponseWriter, code string, message string, httpCode int) {
	err := &StructuredError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
	err.ToHTTP(w)
}

// WriteStandardHTTPError writes a standard error to an HTTP response.
// It uses the standard error code mapping to determine the appropriate HTTP status code.
func WriteStandardHTTPError(w http.ResponseWriter, code string, message string) {
	err := NewStandardError(code, message)
	err.ToHTTP(w)
}

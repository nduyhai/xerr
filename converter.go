package xerr

import (
	"google.golang.org/grpc/codes"
)

// CodeConverter is an interface for converting between HTTP and gRPC status codes.
// This interface allows for different implementations of the conversion logic.
//
// The CodeConverter interface is used by the following functions:
// - FromHTTPJSON: Converts HTTP status codes to gRPC codes when creating an Error from HTTP JSON
// - WriteHTTPError: Converts HTTP status codes to gRPC codes when writing an HTTP error
// - FromGRPCStatus: Converts gRPC codes to HTTP status codes when creating an Error from a gRPC status
//
// By default, the DefaultConverter is used, but you can implement your own converter
// and use it in your application by replacing the DefaultConverter variable.
type CodeConverter interface {
	// HTTPToGRPC converts an HTTP status code to a gRPC status code.
	HTTPToGRPC(httpCode int) codes.Code

	// GRPCToHTTP converts a gRPC status code to an HTTP status code.
	GRPCToHTTP(code codes.Code) int
}

// DefaultCodeConverter is the default implementation of the CodeConverter interface.
// It uses the standard mapping between HTTP and gRPC status codes as defined in the
// gRPC and HTTP specifications.
//
// The DefaultCodeConverter implements the following mappings:
// - HTTP 2xx -> gRPC OK
// - HTTP 400 -> gRPC InvalidArgument
// - HTTP 401 -> gRPC Unauthenticated
// - HTTP 403 -> gRPC PermissionDenied
// - HTTP 404 -> gRPC NotFound
// - HTTP 409 -> gRPC Aborted
// - HTTP 422 -> gRPC FailedPrecondition
// - HTTP 429 -> gRPC ResourceExhausted
// - HTTP 499 -> gRPC Canceled
// - HTTP 500 -> gRPC Internal
// - HTTP 501 -> gRPC Unimplemented
// - HTTP 503 -> gRPC Unavailable
// - HTTP 504 -> gRPC DeadlineExceeded
//
// And the reverse mappings for gRPC to HTTP.
type DefaultCodeConverter struct{}

// HTTPToGRPC converts an HTTP status code to a gRPC status code.
func (c *DefaultCodeConverter) HTTPToGRPC(httpCode int) codes.Code {
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

// GRPCToHTTP converts a gRPC status code to an HTTP status code.
func (c *DefaultCodeConverter) GRPCToHTTP(code codes.Code) int {
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

// DefaultConverter is the default instance of the CodeConverter interface.
// It is used by the package functions that need to convert between HTTP and gRPC status codes.
// You can replace this variable with your own implementation of the CodeConverter interface
// to customize the conversion logic in your application.
//
// Example:
//
//	// Create a custom converter
//	type MyCodeConverter struct{}
//
//	func (c *MyCodeConverter) HTTPToGRPC(httpCode int) codes.Code {
//		// Custom implementation
//	}
//
//	func (c *MyCodeConverter) GRPCToHTTP(code codes.Code) int {
//		// Custom implementation
//	}
//
//	// Replace the default converter
//	func init() {
//		xerr.DefaultConverter = &MyCodeConverter{}
//	}
var DefaultConverter CodeConverter = &DefaultCodeConverter{}

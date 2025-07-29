# xerr

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/xerr)](LICENSE)

A structured error handling package for Go with seamless integration for gRPC and HTTP services.

## Features

- ✅ **Structured Errors** - Rich error objects with code, message, and metadata
- ✅ **Protocol Integration** - Seamless conversion between errors and gRPC/HTTP responses
- ✅ **Standard Error Codes** - Predefined error codes aligned with gRPC and HTTP standards
- ✅ **Error Details** - Support for gRPC error details (ErrorInfo, BadRequest, PreconditionFailure)
- ✅ **Fluent API** - Builder pattern for creating and customizing errors
- ✅ **Error Wrapping** - Wrap existing errors with structured information
- ✅ **Default Error Wrapping** - Wrap errors with default error code
- ✅ **Error Cause Tracking** - Track and retrieve the original cause of errors
- ✅ **Error Unwrapping** - Standard Go error unwrapping support

## Installation

```bash
go get github.com/nduyhai/xerr
```

## Quick Start

```go
package main

import (
	"fmt"
	"net/http"
	
	"github.com/nduyhai/xerr"
)

func main() {
	// Create a simple error
	err := xerr.New("INVALID_INPUT", "The input is invalid")
	fmt.Println(err) // [INVALID_INPUT] The input is invalid
	
	// Create an error with standard code mapping
	err = xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Invalid email format")
	fmt.Printf("HTTP Code: %d, gRPC Code: %d\n", err.HTTPCode, err.GRPCCode)
	
	// Use the fluent API to customize an error
	err = xerr.New("PAYMENT_FAILED", "Payment processing failed")
		.WithReason("Your payment could not be processed. Please try again.")
		.WithHTTPCode(http.StatusBadGateway)
		.WithMetadata("transaction_id", "tx_123456")
	
	// Handle HTTP requests
	http.HandleFunc("/api/example", func(w http.ResponseWriter, r *http.Request) {
		// Write a standard error directly to the response
		xerr.WriteStandardHTTPError(w, xerr.NOT_FOUND, "Resource not found")
	})
}
```

## Usage Examples

### Creating Errors

```go
// Basic error
err := xerr.New("INVALID_INPUT", "The input is invalid")

// Error with standard code mapping
err := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Invalid email format")

// Error with custom HTTP and gRPC codes
err := xerr.NewWithHTTPAndGRPC("RATE_LIMITED", "Too many requests", 429, codes.ResourceExhausted)
```

### Customizing Errors

```go
// Add a user-facing reason
err := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	.WithReason("Your payment could not be processed. Please try again.")

// Add metadata
err := xerr.New("SERVER_ERROR", "Internal server error")
	.WithMetadata("request_id", "req_123456")
	.WithMetadata("server", "api-west-1")

// Add field violations for validation errors
validationErr := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Validation failed")
	.WithBadRequest(map[string]string{
		"email": "Invalid email format",
		"age": "Must be at least 18",
	})
```

### Error Cause Tracking and Unwrapping

```go
// Wrap an error with a structured error
originalErr := errors.New("database connection failed")
wrappedErr := xerr.Wrap(originalErr, xerr.UNAVAILABLE)

// Unwrap to get the original error
unwrappedErr := errors.Unwrap(wrappedErr) // Returns originalErr

// Get the root cause of a deeply nested error using standard Go unwrapping
deeplyNestedErr := fmt.Errorf("operation failed: %w", wrappedErr)
rootCause := errors.Unwrap(deeplyNestedErr) // Returns wrappedErr
rootCause = errors.Unwrap(rootCause) // Returns originalErr
```

### HTTP Integration

```go
// Convert error to HTTP response
func handleRequest(w http.ResponseWriter, r *http.Request) {
	err := validateInput(r)
	if err != nil {
		// If it's already a StructuredError, use it directly
		if se, ok := err.(*xerr.StructuredError); ok {
			se.ToHTTP(w)
			return
		}
		
		// Otherwise, wrap it with a specific code
		xerr.Wrap(err, xerr.INVALID_ARGUMENT).ToHTTP(w)
		return
		
		// Or wrap with default UNKNOWN code
		// xerr.WrapDefault(err).ToHTTP(w)
		// return
	}
	
	// Process request...
}

// Convenience function for writing errors
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	xerr.WriteStandardHTTPError(w, xerr.NOT_FOUND, "Resource not found")
}
```

### gRPC Integration

```go
import (
	"google.golang.org/grpc/status"
	"github.com/nduyhai/xerr"
)

// Server-side: Convert error to gRPC status
func (s *server) MyGRPCMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	err := processRequest(req)
	if err != nil {
		// If it's already a StructuredError, convert to gRPC status
		if se, ok := err.(*xerr.StructuredError); ok {
			return nil, se.ToGRPCStatus().Err()
		}
		
		// Otherwise, wrap it with a specific code
		return nil, xerr.Wrap(err, xerr.INTERNAL).ToGRPCStatus().Err()
		
		// Or wrap with default UNKNOWN code
		// return nil, xerr.WrapDefault(err).ToGRPCStatus().Err()
	}
	
	// Process request...
}

// Client-side: Convert gRPC status to StructuredError
func handleGRPCError(err error) {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			se := xerr.FromGRPCStatus(st)
			fmt.Printf("Error: %s, Code: %s, HTTP Code: %d\n", 
				se.Message, se.Code, se.HTTPCode)
		}
	}
}
```

## Error Codes

The package provides standard error codes that align with both gRPC and HTTP standards:

### General Errors
- `UNKNOWN` - Unknown error
- `INTERNAL` - Internal server error
- `UNAVAILABLE` - Service unavailable
- `TIMEOUT` - Request timeout
- `CANCELLED` - Request cancelled

### Client Errors
- `INVALID_ARGUMENT` - Invalid argument
- `FAILED_PRECONDITION` - Failed precondition
- `OUT_OF_RANGE` - Value out of range
- `UNAUTHENTICATED` - Unauthenticated request
- `PERMISSION_DENIED` - Permission denied
- `NOT_FOUND` - Resource not found
- `ALREADY_EXISTS` - Resource already exists
- `RESOURCE_EXHAUSTED` - Resource quota exceeded
- `ABORTED` - Operation aborted

### Data Errors
- `DATA_LOSS` - Unrecoverable data loss or corruption
- `DATA_VALIDATION` - Data validation error

### Business Logic Errors
- `BUSINESS_RULE` - Business rule violation
- `CONFLICT` - Conflict with current state

## API Documentation

For detailed API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/nduyhai/xerr).

## Sample Code

The repository includes a `samples` directory with working examples for each use case:

- **basic**: Basic error creation and handling
- **standard_codes**: Using standard error codes
- **customization**: Customizing errors with metadata and reasons
- **wrapping**: Error wrapping and unwrapping
- **http**: HTTP error integration
- **grpc**: gRPC error integration
- **details**: Error details usage (BadRequest, ErrorInfo, etc.)

Each sample contains a `main.go` file that demonstrates the specific functionality.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nduyhai/xerr"
)

func main() {
	fmt.Println("Error Wrapping and Unwrapping Example")
	fmt.Println("===================================")

	// Simulate a standard error
	originalErr := simulateStandardError()
	fmt.Printf("Original error: %v\n", originalErr)

	// Wrap with a specific code
	wrappedErr := xerr.Wrap(originalErr, xerr.UNAVAILABLE)
	fmt.Printf("\nWrapped error: %v\n", wrappedErr)
	fmt.Printf("Wrapped error type: %T\n", wrappedErr)
	fmt.Printf("Wrapped error code: %s\n", wrappedErr.GetCode())
	fmt.Printf("Wrapped error HTTP code: %d\n", wrappedErr.GetHTTPCode())
	fmt.Printf("Wrapped error gRPC code: %d\n", wrappedErr.GetGRPCCode())

	// Wrap with default code (UNKNOWN)
	defaultWrappedErr := xerr.WrapDefault(originalErr)
	fmt.Printf("\nDefault wrapped error: %v\n", defaultWrappedErr)
	fmt.Printf("Default wrapped error code: %s\n", defaultWrappedErr.GetCode())

	// Unwrap to get the original error
	unwrappedErr := errors.Unwrap(wrappedErr)
	fmt.Printf("\nUnwrapped error: %v\n", unwrappedErr)
	fmt.Printf("Is unwrapped error the same as original? %v\n", unwrappedErr == originalErr)

	// Demonstrate errors.Is functionality
	fmt.Printf("\nDoes wrapped error satisfy errors.Is(wrappedErr, originalErr)? %v\n",
		errors.Is(wrappedErr, originalErr))

	// Demonstrate deeply nested errors
	nestedErr := simulateNestedError()
	fmt.Printf("\nDeeply nested error: %v\n", nestedErr)

	// Check if the nested error contains the original error
	fmt.Printf("Does nested error satisfy errors.Is(nestedErr, originalErr)? %v\n",
		errors.Is(nestedErr, originalErr))

	// Demonstrate wrapping an already structured error
	structuredErr := xerr.New("ALREADY_STRUCTURED", "This is already a structured error")
	rewrappedErr := xerr.Wrap(structuredErr, xerr.INTERNAL)
	fmt.Printf("\nRe-wrapped structured error: %v\n", rewrappedErr)
	fmt.Printf("Re-wrapped error code (should be updated): %s\n", rewrappedErr.GetCode())
}

// simulateStandardError simulates a standard error from a common operation
func simulateStandardError() error {
	// Try to open a non-existent file
	_, err := os.Open("non_existent_file.txt")
	return err
}

// simulateNestedError creates a deeply nested error
func simulateNestedError() error {
	// Start with a standard error
	baseErr := io.EOF

	// Wrap with xerr
	wrappedErr := xerr.Wrap(baseErr, xerr.DATA_LOSS)

	// Wrap with fmt.Errorf
	return fmt.Errorf("operation failed: %w", wrappedErr)
}

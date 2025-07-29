package main

import (
	"fmt"

	"github.com/nduyhai/xerr"
)

func main() {
	fmt.Println("Standard Error Codes Example")
	fmt.Println("===========================")

	// Create errors with standard error codes
	fmt.Println("General Errors:")
	printErrorDetails(xerr.NewStandardError(xerr.UNKNOWN, "Unknown error occurred"))
	printErrorDetails(xerr.NewStandardError(xerr.INTERNAL, "Internal server error"))
	printErrorDetails(xerr.NewStandardError(xerr.UNAVAILABLE, "Service is currently unavailable"))
	printErrorDetails(xerr.NewStandardError(xerr.TIMEOUT, "Request timed out"))
	printErrorDetails(xerr.NewStandardError(xerr.CANCELLED, "Request was cancelled"))

	fmt.Println("\nClient Errors:")
	printErrorDetails(xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Invalid argument provided"))
	printErrorDetails(xerr.NewStandardError(xerr.FAILED_PRECONDITION, "Failed precondition"))
	printErrorDetails(xerr.NewStandardError(xerr.OUT_OF_RANGE, "Value out of range"))
	printErrorDetails(xerr.NewStandardError(xerr.UNAUTHENTICATED, "Unauthenticated request"))
	printErrorDetails(xerr.NewStandardError(xerr.PERMISSION_DENIED, "Permission denied"))
	printErrorDetails(xerr.NewStandardError(xerr.NOT_FOUND, "Resource not found"))
	printErrorDetails(xerr.NewStandardError(xerr.ALREADY_EXISTS, "Resource already exists"))
	printErrorDetails(xerr.NewStandardError(xerr.RESOURCE_EXHAUSTED, "Resource quota exceeded"))
	printErrorDetails(xerr.NewStandardError(xerr.ABORTED, "Operation aborted"))

	fmt.Println("\nData Errors:")
	printErrorDetails(xerr.NewStandardError(xerr.DATA_LOSS, "Unrecoverable data loss"))
	printErrorDetails(xerr.NewStandardError(xerr.DATA_VALIDATION, "Data validation error"))

	fmt.Println("\nBusiness Logic Errors:")
	printErrorDetails(xerr.NewStandardError(xerr.BUSINESS_RULE, "Business rule violation"))
	printErrorDetails(xerr.NewStandardError(xerr.CONFLICT, "Conflict with current state"))

	// Demonstrate using a non-standard code (falls back to UNKNOWN)
	fmt.Println("\nNon-standard code (falls back to UNKNOWN):")
	printErrorDetails(xerr.NewStandardError("CUSTOM_CODE", "Custom error code"))
}

// Helper function to print error details
func printErrorDetails(err xerr.Error) {
	fmt.Printf("Code: %-20s | Message: %-30s | HTTP: %3d | gRPC: %d\n",
		err.GetMessage(), err.GetMessage(), err.GetHTTPCode(), err.GetGRPCCode())
}

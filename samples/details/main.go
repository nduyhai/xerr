package main

import (
	"fmt"
	
	"github.com/nduyhai/xerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
	fmt.Println("Error Details Example")
	fmt.Println("====================")
	
	// Example 1: Using ErrorInfo
	fmt.Println("Example 1: Using ErrorInfo")
	
	// Create metadata for ErrorInfo
	metadata := map[string]string{
		"request_id": "req_123456",
		"server":     "api-west-1",
		"version":    "1.0.0",
	}
	
	// Create error with ErrorInfo
	errorInfoErr := xerr.NewStandardError(xerr.INTERNAL, "Internal server error")
	errorInfoErr = errorInfoErr.WithErrorInfo("github.com/nduyhai/xerr/samples", metadata)
	
	// Print error details
	fmt.Printf("Error: %v\n", errorInfoErr)
	fmt.Printf("Error code: %s\n", errorInfoErr.Code)
	fmt.Printf("Error message: %s\n", errorInfoErr.Message)
	
	fmt.Println("Metadata:")
	for k, v := range errorInfoErr.Metadata {
		fmt.Printf("  %s: %s\n", k, v)
	}
	
	// Get ErrorInfo for gRPC status
	errorInfo := errorInfoErr.GetErrorInfo()
	fmt.Printf("\nErrorInfo for gRPC:\n")
	fmt.Printf("  Reason: %s\n", errorInfo.Reason)
	fmt.Printf("  Domain: %s\n", errorInfo.Domain)
	fmt.Printf("  Metadata count: %d\n", len(errorInfo.Metadata))
	
	// Example 2: Using BadRequest for validation errors
	fmt.Println("\nExample 2: Using BadRequest for validation errors")
	
	// Create field violations for BadRequest
	fieldViolations := map[string]string{
		"email":    "Invalid email format",
		"password": "Password must be at least 8 characters",
		"username": "Username cannot contain special characters",
	}
	
	// Create error with BadRequest
	badRequestErr := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Validation failed")
	badRequestErr = badRequestErr.WithBadRequest(fieldViolations)
	
	// Print error details
	fmt.Printf("Error: %v\n", badRequestErr)
	fmt.Printf("Error code: %s\n", badRequestErr.Code)
	fmt.Printf("Error message: %s\n", badRequestErr.Message)
	
	fmt.Println("Field violations (from metadata):")
	for k, v := range badRequestErr.Metadata {
		if len(k) > 6 && k[:6] == "field:" {
			field := k[6:] // Remove "field:" prefix
			fmt.Printf("  %s: %s\n", field, v)
		}
	}
	
	// Get BadRequest for gRPC status
	badRequest := badRequestErr.GetBadRequest()
	fmt.Printf("\nBadRequest for gRPC:\n")
	fmt.Printf("  Field violations count: %d\n", len(badRequest.FieldViolations))
	
	for _, violation := range badRequest.FieldViolations {
		fmt.Printf("  Field: %s, Description: %s\n", violation.Field, violation.Description)
	}
	
	// Example 3: Using PreconditionFailure
	fmt.Println("\nExample 3: Using PreconditionFailure")
	
	// Create precondition violations
	preconditionViolations := map[string]string{
		"account_status": "Account must be active",
		"balance":        "Insufficient balance",
		"permissions":    "User does not have required permissions",
	}
	
	// Create error with PreconditionFailure
	preconditionErr := xerr.NewStandardError(xerr.FAILED_PRECONDITION, "Precondition check failed")
	preconditionErr = preconditionErr.WithPreconditionFailure(preconditionViolations)
	
	// Print error details
	fmt.Printf("Error: %v\n", preconditionErr)
	fmt.Printf("Error code: %s\n", preconditionErr.Code)
	fmt.Printf("Error message: %s\n", preconditionErr.Message)
	
	fmt.Println("Precondition violations (from metadata):")
	for k, v := range preconditionErr.Metadata {
		if len(k) > 13 && k[:13] == "precondition:" {
			condition := k[13:] // Remove "precondition:" prefix
			fmt.Printf("  %s: %s\n", condition, v)
		}
	}
	
	// Get PreconditionFailure for gRPC status
	preconditionFailure := preconditionErr.GetPreconditionFailure()
	fmt.Printf("\nPreconditionFailure for gRPC:\n")
	fmt.Printf("  Violations count: %d\n", len(preconditionFailure.Violations))
	
	for _, violation := range preconditionFailure.Violations {
		fmt.Printf("  Type: %s, Subject: %s, Description: %s\n", 
			violation.Type, violation.Subject, violation.Description)
	}
	
	// Example 4: Combining multiple error details
	fmt.Println("\nExample 4: Combining multiple error details")
	
	// Create a complex error with multiple details
	complexErr := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Multiple validation issues")
	complexErr = complexErr.WithReason("Please correct the issues and try again")
	complexErr = complexErr.WithBadRequest(map[string]string{
		"email": "Invalid email format",
		"age":   "Must be at least 18",
	})
	complexErr = complexErr.WithPreconditionFailure(map[string]string{
		"account_status": "Account must be verified",
	})
	complexErr = complexErr.WithMetadata("request_id", "req_789012")
	
	// Print error details
	fmt.Printf("Complex error: %v\n", complexErr)
	fmt.Printf("User-facing reason: %s\n", complexErr.Reason)
	
	fmt.Println("All metadata:")
	for k, v := range complexErr.Metadata {
		fmt.Printf("  %s: %s\n", k, v)
	}
	
	// Convert to gRPC status to see all details
	st := complexErr.ToGRPCStatus()
	fmt.Printf("\ngRPC Status with all details:\n")
	fmt.Printf("  Code: %d\n", st.Code())
	fmt.Printf("  Message: %s\n", st.Message())
	fmt.Printf("  Details count: %d\n", len(st.Details()))
	
	// Print all details
	fmt.Println("gRPC Status Details:")
	for i, detail := range st.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			fmt.Printf("  Detail %d: ErrorInfo - Reason: %s, Domain: %s\n", 
				i+1, d.Reason, d.Domain)
		case *errdetails.BadRequest:
			fmt.Printf("  Detail %d: BadRequest - %d field violations\n", 
				i+1, len(d.FieldViolations))
		case *errdetails.PreconditionFailure:
			fmt.Printf("  Detail %d: PreconditionFailure - %d violations\n", 
				i+1, len(d.Violations))
		case *errdetails.LocalizedMessage:
			fmt.Printf("  Detail %d: LocalizedMessage - %s\n", 
				i+1, d.Message)
		default:
			fmt.Printf("  Detail %d: %T\n", i+1, d)
		}
	}
}
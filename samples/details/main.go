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
 // Type assertion to *StructuredError to access detail-specific methods
 structuredErr, ok := errorInfoErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }
 structuredErr = structuredErr.WithErrorInfo("github.com/nduyhai/xerr/samples", metadata).(*xerr.StructuredError)

 // Print error details
 fmt.Printf("Error: %v\n", structuredErr)
 fmt.Printf("Error code: %s\n", structuredErr.GetCode())
 fmt.Printf("Error message: %s\n", structuredErr.GetMessage())

 fmt.Println("Metadata:")
 for k, v := range structuredErr.GetMetadata() {
 	fmt.Printf("  %s: %s\n", k, v)
 }

 // Get ErrorInfo for gRPC status
 errorInfo := structuredErr.GetErrorInfo()
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
 // Type assertion to *StructuredError to access detail-specific methods
 badRequestStructErr, ok := badRequestErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }
 badRequestStructErr = badRequestStructErr.WithBadRequest(fieldViolations).(*xerr.StructuredError)

 // Print error details
 fmt.Printf("Error: %v\n", badRequestStructErr)
 fmt.Printf("Error code: %s\n", badRequestStructErr.GetCode())
 fmt.Printf("Error message: %s\n", badRequestStructErr.GetMessage())

 fmt.Println("Field violations (from metadata):")
 for k, v := range badRequestStructErr.GetMetadata() {
 	if len(k) > 6 && k[:6] == "field:" {
 		field := k[6:] // Remove "field:" prefix
 		fmt.Printf("  %s: %s\n", field, v)
 	}
 }

 // Get BadRequest for gRPC status
 badRequest := badRequestStructErr.GetBadRequest()
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
 // Type assertion to *StructuredError to access detail-specific methods
 preconditionStructErr, ok := preconditionErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }
 preconditionStructErr = preconditionStructErr.WithPreconditionFailure(preconditionViolations).(*xerr.StructuredError)

 // Print error details
 fmt.Printf("Error: %v\n", preconditionStructErr)
 fmt.Printf("Error code: %s\n", preconditionStructErr.GetCode())
 fmt.Printf("Error message: %s\n", preconditionStructErr.GetMessage())

 fmt.Println("Precondition violations (from metadata):")
 for k, v := range preconditionStructErr.GetMetadata() {
 	if len(k) > 13 && k[:13] == "precondition:" {
 		condition := k[13:] // Remove "precondition:" prefix
 		fmt.Printf("  %s: %s\n", condition, v)
 	}
 }

 // Get PreconditionFailure for gRPC status
 preconditionFailure := preconditionStructErr.GetPreconditionFailure()
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
 // Type assertion to *StructuredError to access detail-specific methods
 complexStructErr, ok := complexErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }
 complexStructErr = complexStructErr.WithReason("Please correct the issues and try again").(*xerr.StructuredError)
 complexStructErr = complexStructErr.WithBadRequest(map[string]string{
 	"email": "Invalid email format",
 	"age":   "Must be at least 18",
 }).(*xerr.StructuredError)
 complexStructErr = complexStructErr.WithPreconditionFailure(map[string]string{
 	"account_status": "Account must be verified",
 }).(*xerr.StructuredError)
 complexStructErr = complexStructErr.WithMetadata("request_id", "req_789012").(*xerr.StructuredError)

 // Print error details
 fmt.Printf("Complex error: %v\n", complexStructErr)
 fmt.Printf("User-facing reason: %s\n", complexStructErr.GetReason())

 fmt.Println("All metadata:")
 for k, v := range complexStructErr.GetMetadata() {
 	fmt.Printf("  %s: %s\n", k, v)
 }

 // Convert to gRPC status to see all details
 st := complexStructErr.ToGRPCStatus()
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
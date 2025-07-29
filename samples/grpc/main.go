package main

import (
	"fmt"
	
	"github.com/nduyhai/xerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("gRPC Error Integration Example")
	fmt.Println("============================")
	
 // Example 1: Convert a structured error to a gRPC status
 fmt.Println("Example 1: Convert a structured error to a gRPC status")
 err := xerr.NewStandardError(xerr.NOT_FOUND, "User not found")
 err = err.WithMetadata("user_id", "12345")
 err = err.WithReason("The requested user could not be found in the database")

 // Type assertion to *StructuredError to access ToGRPCStatus method
 structuredErr, ok := err.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 // Convert to gRPC status
 st := structuredErr.ToGRPCStatus()
	
	fmt.Printf("gRPC Status Code: %d\n", st.Code())
	fmt.Printf("gRPC Status Message: %s\n", st.Message())
	fmt.Printf("gRPC Status Proto: %v\n", st.Proto())
	fmt.Printf("gRPC Status Details Count: %d\n", len(st.Details()))
	
	// Print details
	fmt.Println("gRPC Status Details:")
	for i, detail := range st.Details() {
		fmt.Printf("  Detail %d: %T - %v\n", i+1, detail, detail)
	}
	
 // Example 2: Convert a gRPC status back to a structured error
 fmt.Println("\nExample 2: Convert a gRPC status back to a structured error")
 recoveredErr := xerr.FromGRPCStatus(st)

 // Type assertion to *StructuredError to access fields
 recoveredStructErr, ok := recoveredErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 fmt.Printf("Recovered Error Code: %s\n", recoveredStructErr.GetCode())
 fmt.Printf("Recovered Error Message: %s\n", recoveredStructErr.GetMessage())
 fmt.Printf("Recovered Error Reason: %s\n", recoveredStructErr.GetReason())
 fmt.Printf("Recovered Error HTTP Code: %d\n", recoveredStructErr.GetHTTPCode())
 fmt.Printf("Recovered Error gRPC Code: %d\n", recoveredStructErr.GetGRPCCode())

 fmt.Println("Recovered Error Metadata:")
 for k, v := range recoveredStructErr.GetMetadata() {
 	fmt.Printf("  %s: %s\n", k, v)
 }
	
 // Example 3: Create a gRPC status and convert it to a structured error
 fmt.Println("\nExample 3: Create a gRPC status and convert it to a structured error")

 // Create a gRPC status directly
 grpcStatus := status.New(codes.FailedPrecondition, "Precondition failed")

 // Convert to structured error
 convertedErr := xerr.FromGRPCStatus(grpcStatus)

 // Type assertion to *StructuredError to access fields
 convertedStructErr, ok := convertedErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 fmt.Printf("Converted Error Code: %s\n", convertedStructErr.GetCode())
 fmt.Printf("Converted Error Message: %s\n", convertedStructErr.GetMessage())
 fmt.Printf("Converted Error HTTP Code: %d\n", convertedStructErr.GetHTTPCode())
 fmt.Printf("Converted Error gRPC Code: %d\n", convertedStructErr.GetGRPCCode())
	
	// Example 4: Demonstrate server-side error handling in gRPC
	fmt.Println("\nExample 4: Demonstrate server-side error handling in gRPC")
	
	// Simulate a gRPC service method that returns an error
	_, grpcErr := simulateGRPCMethod("12345")
	
	fmt.Printf("gRPC Error: %v\n", grpcErr)
	
	// Example 5: Demonstrate client-side error handling in gRPC
	fmt.Println("\nExample 5: Demonstrate client-side error handling in gRPC")
	
	// Simulate handling a gRPC error on the client side
	handleGRPCError(grpcErr)
}

// simulateGRPCMethod simulates a gRPC service method that returns an error
func simulateGRPCMethod(userID string) (interface{}, error) {
	// Simulate a not found error
	err := xerr.NewStandardError(xerr.NOT_FOUND, fmt.Sprintf("User %s not found", userID))
	err = err.WithMetadata("user_id", userID)
	err = err.WithReason("The requested user could not be found in the database")
	
	// Type assertion to *StructuredError to access ToGRPCStatus method
	structuredErr, ok := err.(*xerr.StructuredError)
	if !ok {
		return nil, fmt.Errorf("failed to convert to *StructuredError")
	}
	
	// Convert to gRPC status and return as error
	return nil, structuredErr.ToGRPCStatus().Err()
}

// handleGRPCError simulates handling a gRPC error on the client side
func handleGRPCError(err error) {
	if err == nil {
		fmt.Println("No error to handle")
		return
	}
	
	// Extract gRPC status from error
	st, ok := status.FromError(err)
	if !ok {
		fmt.Printf("Not a gRPC status error: %v\n", err)
		return
	}
	
 // Convert to structured error
 se := xerr.FromGRPCStatus(st)

 // Type assertion to *StructuredError to access fields
 seStruct, ok := se.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 fmt.Printf("Client received structured error:\n")
 fmt.Printf("  Code: %s\n", seStruct.GetCode())
 fmt.Printf("  Message: %s\n", seStruct.GetMessage())
 fmt.Printf("  Reason: %s\n", seStruct.GetReason())
 fmt.Printf("  HTTP Code: %d\n", seStruct.GetHTTPCode())
 fmt.Printf("  gRPC Code: %d\n", seStruct.GetGRPCCode())

 // Handle based on error code
 switch seStruct.GetCode() {
 case xerr.NOT_FOUND:
 	fmt.Println("  Handling not found error: Show 'not found' UI")
 case xerr.PERMISSION_DENIED:
 	fmt.Println("  Handling permission denied: Prompt for authentication")
 case xerr.INVALID_ARGUMENT:
 	fmt.Println("  Handling invalid argument: Show validation errors")
 default:
 	fmt.Println("  Handling generic error: Show error message")
 }
}
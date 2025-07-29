package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/nduyhai/xerr"
)

func main() {
	fmt.Println("HTTP Error Integration Example")
	fmt.Println("============================")

 // Example 1: Convert a structured error to an HTTP response
 fmt.Println("Example 1: Convert a structured error to an HTTP response")
 err := xerr.NewStandardError(xerr.NOT_FOUND, "User not found")
 err = err.WithMetadata("user_id", "12345")

 // Type assertion to *StructuredError to access ToHTTP method
 structuredErr, ok := err.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 // Create a test HTTP response recorder
 w := httptest.NewRecorder()

 // Write the error to the HTTP response
 structuredErr.ToHTTP(w)

	// Print the HTTP response
	fmt.Printf("HTTP Status Code: %d\n", w.Code)
	fmt.Printf("HTTP Response Body: %s\n", w.Body.String())

	// Example 2: Using WriteHTTPError convenience function
	fmt.Println("\nExample 2: Using WriteHTTPError convenience function")
	w2 := httptest.NewRecorder()
	xerr.WriteHTTPError(w2, "RATE_LIMITED", "Too many requests", 429)

	fmt.Printf("HTTP Status Code: %d\n", w2.Code)
	fmt.Printf("HTTP Response Body: %s\n", w2.Body.String())

	// Example 3: Using WriteStandardHTTPError convenience function
	fmt.Println("\nExample 3: Using WriteStandardHTTPError convenience function")
	w3 := httptest.NewRecorder()
	xerr.WriteStandardHTTPError(w3, xerr.PERMISSION_DENIED, "Access denied")

	fmt.Printf("HTTP Status Code: %d\n", w3.Code)
	fmt.Printf("HTTP Response Body: %s\n", w3.Body.String())

 // Example 4: Convert a structured error to JSON bytes
 fmt.Println("\nExample 4: Convert a structured error to JSON bytes")
 validationErr := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Validation failed")
 validationErr = validationErr.WithMetadata("field:email", "Invalid email format")
 validationErr = validationErr.WithMetadata("field:age", "Must be at least 18")

 // Type assertion to *StructuredError to access ToHTTPJSON method
 validationStructErr, ok := validationErr.(*xerr.StructuredError)
 if !ok {
 	fmt.Println("Error: Failed to convert to *StructuredError")
 	return
 }

 jsonBytes, statusCode := validationStructErr.ToHTTPJSON()
	fmt.Printf("HTTP Status Code: %d\n", statusCode)
	fmt.Printf("JSON Bytes: %s\n", string(jsonBytes))

	// Example 5: Convert JSON bytes back to a structured error
	fmt.Println("\nExample 5: Convert JSON bytes back to a structured error")
	recoveredErr, parseErr := xerr.FromHTTPJSON(jsonBytes, statusCode)
	if parseErr != nil {
		log.Fatalf("Failed to parse JSON: %v", parseErr)
	}

	fmt.Printf("Recovered Error Code: %s\n", recoveredErr.GetCode())
	fmt.Printf("Recovered Error Message: %s\n", recoveredErr.GetMessage())
	fmt.Printf("Recovered Error HTTP Code: %d\n", recoveredErr.GetHTTPCode())
	fmt.Printf("Recovered Error gRPC Code: %d\n", recoveredErr.GetGRPCCode())
	fmt.Println("Recovered Error Metadata:")
	for k, v := range recoveredErr.GetMetadata() {
		fmt.Printf("  %s: %s\n", k, v)
	}

	// Example 6: Demonstrate a complete HTTP handler
	fmt.Println("\nExample 6: Demonstrate a complete HTTP handler")

	// Create a test request
	req := httptest.NewRequest("GET", "/api/users/12345", nil)
	w6 := httptest.NewRecorder()

	// Call our example handler
	handleUserRequest(w6, req)

	fmt.Printf("HTTP Status Code: %d\n", w6.Code)
	fmt.Printf("HTTP Response Body: %s\n", w6.Body.String())
}

// handleUserRequest is an example HTTP handler that demonstrates error handling
func handleUserRequest(w http.ResponseWriter, r *http.Request) {
	// Simulate a validation error
	if err := validateRequest(r); err != nil {
		// If it's already a StructuredError, use it directly
		if se, ok := err.(*xerr.StructuredError); ok {
			se.ToHTTP(w)
			return
		}

		// Otherwise, wrap it with a specific code
		wrappedErr := xerr.Wrap(err, xerr.INVALID_ARGUMENT)
		
		// Type assertion to *StructuredError to access ToHTTP method
		wrappedStructErr, ok := wrappedErr.(*xerr.StructuredError)
		if !ok {
			fmt.Fprintf(w, "Error: Failed to convert to *StructuredError")
			return
		}
		
		wrappedStructErr.ToHTTP(w)
		return
	}

	// Simulate a not found error
	userID := "12345" // Extract from request path in a real application
	user, err := findUser(userID)
	if err != nil {
		xerr.WriteStandardHTTPError(w, xerr.NOT_FOUND, fmt.Sprintf("User %s not found", userID))
		return
	}

	// Success case - would normally return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	})
}

// validateRequest simulates request validation
func validateRequest(r *http.Request) error {
	// Simulate a validation error
	validationErr := xerr.NewStandardError(xerr.INVALID_ARGUMENT, "Invalid request parameters")
	validationErr = validationErr.WithMetadata("field:query", "Missing required parameter")
	return validationErr
}

// User represents a simple user model
type User struct {
	ID   string
	Name string
}

// findUser simulates a database lookup
func findUser(id string) (*User, error) {
	// Simulate a not found error
	return nil, fmt.Errorf("user not found in database")
}

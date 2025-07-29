package main

import (
	"errors"
	"fmt"
	
	"github.com/nduyhai/xerr"
)

func main() {
	fmt.Println("Basic Error Creation and Handling Example")
	fmt.Println("========================================")
	
	// Create a simple error with a custom code and message
	err := xerr.New("INVALID_INPUT", "The input is invalid")
	fmt.Println("Simple error:", err)
	
	// Access the error's properties
	fmt.Printf("Error code: %s\n", err.GetCode())
	fmt.Printf("Error message: %s\n", err.GetMessage())
	fmt.Printf("Default HTTP code: %d\n", err.GetHTTPCode())
	fmt.Printf("Default gRPC code: %d\n", err.GetGRPCCode())
	
	// Create an error with custom HTTP and gRPC codes
	customErr := xerr.NewWithHTTPAndGRPC(
		"RATE_LIMITED", 
		"Too many requests", 
		429, // HTTP 429 Too Many Requests
		8,   // gRPC ResourceExhausted
	)
	fmt.Println("\nCustom error with specific codes:", customErr)
	fmt.Printf("HTTP code: %d\n", customErr.GetHTTPCode())
	fmt.Printf("gRPC code: %d\n", customErr.GetGRPCCode())
	
	// Demonstrate error comparison using errors.Is
	fmt.Println("\nError comparison:")
	sameCodeErr := xerr.New("INVALID_INPUT", "Different message but same code")
	fmt.Printf("Original error: %v\n", err)
	fmt.Printf("Same code error: %v\n", sameCodeErr)
	fmt.Printf("Are they considered the same error? %v\n", errors.Is(err, sameCodeErr))
	
	differentCodeErr := xerr.New("DIFFERENT_CODE", "Different code")
	fmt.Printf("Different code error: %v\n", differentCodeErr)
	fmt.Printf("Are they considered the same error? %v\n", errors.Is(err, differentCodeErr))
}
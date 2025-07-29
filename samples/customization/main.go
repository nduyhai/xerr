package main

import (
	"fmt"
	"net/http"

	"github.com/nduyhai/xerr"
	"google.golang.org/grpc/codes"
)

func main() {
	fmt.Println("Error Customization Example")
	fmt.Println("==========================")

	// Basic error
	basicErr := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	fmt.Println("Basic error:", basicErr)

	// Add a user-facing reason
	withReason := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	withReason = withReason.WithReason("Your payment could not be processed. Please try again or use a different payment method.")
	fmt.Println("\nWith reason:")
	fmt.Printf("Developer message: %s\n", withReason.GetMessage())
	fmt.Printf("User-facing reason: %s\n", withReason.GetReason())

	// Customize HTTP and gRPC codes
	withCodes := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	withCodes = withCodes.WithHTTPCode(http.StatusBadGateway) // 502
	withCodes = withCodes.WithGRPCCode(codes.Unavailable)     // 14
	fmt.Println("\nWith custom codes:")
	fmt.Printf("HTTP code: %d\n", withCodes.GetHTTPCode())
	fmt.Printf("gRPC code: %d\n", withCodes.GetGRPCCode())

	// Add metadata
	withMetadata := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	withMetadata = withMetadata.WithMetadata("transaction_id", "tx_123456")
	withMetadata = withMetadata.WithMetadata("payment_provider", "stripe")
	withMetadata = withMetadata.WithMetadata("amount", "99.99")
	fmt.Println("\nWith metadata:")
	fmt.Printf("Error: %s\n", withMetadata)
	fmt.Println("Metadata:")
	for key, value := range withMetadata.GetMetadata() {
		fmt.Printf("  %s: %s\n", key, value)
	}

	// Combine all customizations
	fullyCustomized := xerr.New("PAYMENT_FAILED", "Payment processing failed")
	fullyCustomized = fullyCustomized.WithReason("Your payment could not be processed. Please try again or use a different payment method.")
	fullyCustomized = fullyCustomized.WithHTTPCode(http.StatusBadGateway)
	fullyCustomized = fullyCustomized.WithGRPCCode(codes.Unavailable)
	fullyCustomized = fullyCustomized.WithMetadata("transaction_id", "tx_123456")
	fullyCustomized = fullyCustomized.WithMetadata("payment_provider", "stripe")
	fullyCustomized = fullyCustomized.WithMetadata("amount", "99.99")

	fmt.Println("\nFully customized error:")
	fmt.Printf("Error: %s\n", fullyCustomized)
	fmt.Printf("Reason: %s\n", fullyCustomized.GetReason())
	fmt.Printf("HTTP code: %d\n", fullyCustomized.GetHTTPCode())
	fmt.Printf("gRPC code: %d\n", fullyCustomized.GetGRPCCode())
	fmt.Println("Metadata:")
	for key, value := range fullyCustomized.GetMetadata() {
		fmt.Printf("  %s: %s\n", key, value)
	}
}

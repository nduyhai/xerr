package xerr

// Reason interface defines the contract for error reason implementations.
// This interface allows for more flexible error reason handling.
type Reason interface {
	// Code returns the machine-readable error code.
	// e.g. "AUTH.USER.INVALID_PASSWORD"
	Code() string

	// Message returns the developer-facing error message.
	Message() string

	// Reason returns an optional user-friendly or localized error message.
	Reason() string
}

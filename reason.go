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

// DefaultReason is the default implementation of the Reason interface.
// It provides a simple struct-based implementation of the Reason interface.
type DefaultReason struct {
	code    string
	message string
	reason  string
}

// NewDefaultReason creates a new DefaultReason with the given code and message.
func NewDefaultReason(code string, message string) *DefaultReason {
	return &DefaultReason{
		code:    code,
		message: message,
	}
}

// Code returns the machine-readable error code.
func (r *DefaultReason) Code() string {
	return r.code
}

// Message returns the developer-facing error message.
func (r *DefaultReason) Message() string {
	return r.message
}

// Reason returns the user-friendly or localized error message.
func (r *DefaultReason) Reason() string {
	return r.reason
}

// WithReason adds a user-friendly reason to the DefaultReason.
func (r *DefaultReason) WithReason(reason string) *DefaultReason {
	r.reason = reason
	return r
}
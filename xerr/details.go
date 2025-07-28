package xerr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// WithErrorInfo adds ErrorInfo detail to the structured error.
// ErrorInfo is a standard gRPC error detail that provides structured error information.
func (e *StructuredError) WithErrorInfo(domain string, metadata map[string]string) *StructuredError {
	// Store metadata in the error
	if metadata != nil {
		if e.Metadata == nil {
			e.Metadata = make(map[string]string)
		}
		for k, v := range metadata {
			e.Metadata[k] = v
		}
	}

	return e
}

// WithBadRequest adds field violations to the error.
// This is useful for validation errors where multiple fields have issues.
func (e *StructuredError) WithBadRequest(fieldViolations map[string]string) *StructuredError {
	// Store field violations in metadata with a special prefix
	if fieldViolations != nil {
		if e.Metadata == nil {
			e.Metadata = make(map[string]string)
		}
		for field, description := range fieldViolations {
			e.Metadata["field:"+field] = description
		}
	}

	return e
}

// GetErrorInfo extracts ErrorInfo from the structured error.
// This is used when converting to gRPC status.
func (e *StructuredError) GetErrorInfo() *errdetails.ErrorInfo {
	return &errdetails.ErrorInfo{
		Reason:   e.Code,
		Domain:   "github.com/nduyhai/xerr",
		Metadata: e.Metadata,
	}
}

// GetBadRequest extracts BadRequest field violations from the structured error.
// This is used when converting to gRPC status.
func (e *StructuredError) GetBadRequest() *errdetails.BadRequest {
	if e.Metadata == nil {
		return nil
	}

	var fieldViolations []*errdetails.BadRequest_FieldViolation

	// Extract field violations from metadata
	for k, v := range e.Metadata {
		if len(k) > 6 && k[:6] == "field:" {
			field := k[6:] // Remove "field:" prefix
			fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: v,
			})
		}
	}

	if len(fieldViolations) == 0 {
		return nil
	}

	return &errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}
}

// WithPreconditionFailure adds precondition failures to the error.
// This is useful for errors where certain preconditions were not met.
func (e *StructuredError) WithPreconditionFailure(violations map[string]string) *StructuredError {
	// Store precondition violations in metadata with a special prefix
	if violations != nil {
		if e.Metadata == nil {
			e.Metadata = make(map[string]string)
		}
		for condition, description := range violations {
			e.Metadata["precondition:"+condition] = description
		}
	}

	return e
}

// GetPreconditionFailure extracts PreconditionFailure from the structured error.
// This is used when converting to gRPC status.
func (e *StructuredError) GetPreconditionFailure() *errdetails.PreconditionFailure {
	if e.Metadata == nil {
		return nil
	}

	var violations []*errdetails.PreconditionFailure_Violation

	// Extract precondition violations from metadata
	for k, v := range e.Metadata {
		if len(k) > 13 && k[:13] == "precondition:" {
			condition := k[13:] // Remove "precondition:" prefix
			violations = append(violations, &errdetails.PreconditionFailure_Violation{
				Type:        "PRECONDITION_FAILURE",
				Subject:     condition,
				Description: v,
			})
		}
	}

	if len(violations) == 0 {
		return nil
	}

	return &errdetails.PreconditionFailure{
		Violations: violations,
	}
}

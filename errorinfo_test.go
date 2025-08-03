package xerr

import "testing"

func TestWithErrorInfoDomain(t *testing.T) {
	err := New("EXAMPLE", "example message")
	se, ok := err.(*StructuredError)
	if !ok {
		t.Fatalf("expected *StructuredError, got %T", err)
	}
	se.WithErrorInfo("example.com", map[string]string{"k": "v"})

	info := se.GetErrorInfo()
	if info.Domain != "example.com" {
		t.Fatalf("expected domain example.com, got %s", info.Domain)
	}
	if info.Metadata["k"] != "v" {
		t.Fatalf("expected metadata k=v, got %v", info.Metadata)
	}
}

func TestFromGRPCStatusDomain(t *testing.T) {
	err := New("CODE", "msg")
	se := err.(*StructuredError)
	se.WithErrorInfo("service.domain", nil)
	st := se.ToGRPCStatus()
	converted := FromGRPCStatus(st)
	se2, ok := converted.(*StructuredError)
	if !ok {
		t.Fatalf("expected *StructuredError, got %T", converted)
	}
	if se2.Domain != "service.domain" {
		t.Fatalf("expected domain service.domain, got %s", se2.Domain)
	}
}

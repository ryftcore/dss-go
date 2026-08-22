package enumerations

import "testing"

func TestQCIdentMethodEnumFields(t *testing.T) {
	tests := []struct {
		v           QCIdentMethodEnum
		description string
		oid         string
	}{
		{QCIdentMethodEnumQCTEIDAS2ACD, "qc-ident-method-eIDAS2-acd", "0.4.0.1862.1.8.3"},
		{QCIdentMethodEnumQCTEIDAS2B, "qc-ident-method-eIDAS2-acd", "0.4.0.1862.1.8.4"},
	}
	for _, tt := range tests {
		if got := tt.v.OID(); got != tt.oid {
			t.Errorf("%v.OID() = %q, want %q", tt.v, got, tt.oid)
		}
		if got := tt.v.Description(); got != tt.description {
			t.Errorf("%v.Description() = %q, want %q", tt.v, got, tt.description)
		}
	}
}

func TestQCIdentMethodEnumForLabel(t *testing.T) {
	// Both constants share the same description upstream; forLabel resolves
	// to the first declared match.
	if got := QCIdentMethodEnumForLabel("qc-ident-method-eIDAS2-acd"); got != QCIdentMethodEnumQCTEIDAS2ACD {
		t.Errorf("QCIdentMethodEnumForLabel(shared) = %q, want %q", got, QCIdentMethodEnumQCTEIDAS2ACD)
	}
	if got := QCIdentMethodEnumForLabel("bogus"); got != "" {
		t.Errorf("QCIdentMethodEnumForLabel(bogus) = %q, want empty", got)
	}
}

func TestQCIdentMethodEnumForLabelPanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty description, got none")
		}
	}()
	QCIdentMethodEnumForLabel("")
}

func TestQCIdentMethodFromOIDResolvesEnum(t *testing.T) {
	// Sanity-check the externally-defined QCIdentMethodFromOID resolves
	// against this file's QCIdentMethodEnumValues().
	got := QCIdentMethodFromOID("0.4.0.1862.1.8.3")
	if got.OID() != "0.4.0.1862.1.8.3" {
		t.Errorf("QCIdentMethodFromOID known oid = %q, want %q", got.OID(), "0.4.0.1862.1.8.3")
	}
	unknown := QCIdentMethodFromOID("9.9.9")
	if unknown.Description() != QCIdentMethodUnknownMethod {
		t.Errorf("QCIdentMethodFromOID unknown description = %q, want %q", unknown.Description(), QCIdentMethodUnknownMethod)
	}
}

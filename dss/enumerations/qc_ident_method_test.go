// Ported from dss-enumerations/.../QCIdentMethod.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestQCIdentMethod_UNKNOWN_METHOD_Value(t *testing.T) {
	if QCIdentMethodUnknownMethod != "qc-identification-method-unknown" {
		t.Errorf("QCIdentMethodUnknownMethod = %q, want %q", QCIdentMethodUnknownMethod, "qc-identification-method-unknown")
	}
}

func TestQCIdentMethodUnknownFallback(t *testing.T) {
	// Exercises the fallback value shape directly (qcIdentMethodUnknown),
	// since QCIdentMethodFromOID's "no match" path additionally depends on
	// QCIdentMethodEnumValues(), which is defined outside this manifest and
	// is not available to test in isolation — see PORTER_BRIEF notes.
	fallback := &qcIdentMethodUnknown{oid: "1.2.3.4"}
	if fallback.OID() != "1.2.3.4" {
		t.Errorf("OID() = %q, want %q", fallback.OID(), "1.2.3.4")
	}
	if fallback.Description() != QCIdentMethodUnknownMethod {
		t.Errorf("Description() = %q, want %q", fallback.Description(), QCIdentMethodUnknownMethod)
	}
}

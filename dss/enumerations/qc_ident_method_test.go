// Ported from dss-enumerations/.../QCIdentMethod.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestQCIdentMethod_UNKNOWN_METHOD_Value(t *testing.T) {
	if QCIdentMethod_UNKNOWN_METHOD != "qc-identification-method-unknown" {
		t.Errorf("QCIdentMethod_UNKNOWN_METHOD = %q, want %q", QCIdentMethod_UNKNOWN_METHOD, "qc-identification-method-unknown")
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
	if fallback.Description() != QCIdentMethod_UNKNOWN_METHOD {
		t.Errorf("Description() = %q, want %q", fallback.Description(), QCIdentMethod_UNKNOWN_METHOD)
	}
}

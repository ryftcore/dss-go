// Ported from dss-enumerations/.../QCType.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestQCType_UNKNOWN_TYPE_Value(t *testing.T) {
	if QCTypeUnknownType != "type-unknown" {
		t.Errorf("QCTypeUnknownType = %q, want %q", QCTypeUnknownType, "type-unknown")
	}
}

func TestQCTypeUnknownFallback(t *testing.T) {
	// Exercises the fallback value shape directly (qcTypeUnknown), since
	// QCTypeFromOID's "no match" path additionally depends on
	// QCTypeEnumValues(), which is defined outside this manifest and is not
	// available to test in isolation.
	fallback := &qcTypeUnknown{oid: "1.2.3.4"}
	if fallback.OID() != "1.2.3.4" {
		t.Errorf("OID() = %q, want %q", fallback.OID(), "1.2.3.4")
	}
	if fallback.Description() != QCTypeUnknownType {
		t.Errorf("Description() = %q, want %q", fallback.Description(), QCTypeUnknownType)
	}
}

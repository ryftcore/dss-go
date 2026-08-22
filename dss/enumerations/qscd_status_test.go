// Ported from dss-enumerations/.../QSCDStatus.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestQSCDStatus(t *testing.T) {
	if len(QSCDStatusValues()) != 2 {
		t.Fatalf("expected 2 values, got %d", len(QSCDStatusValues()))
	}
	for _, v := range QSCDStatusValues() {
		got, err := QSCDStatusValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("QSCDStatusValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := QSCDStatusValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
	if !QSCDStatusIsQSCD(QSCDStatusQSCD) {
		t.Error("expected QSCDStatusIsQSCD(QSCD) = true")
	}
	if QSCDStatusIsQSCD(QSCDStatusNotQSCD) {
		t.Error("expected QSCDStatusIsQSCD(NOT_QSCD) = false")
	}
}

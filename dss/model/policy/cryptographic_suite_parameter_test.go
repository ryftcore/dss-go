// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteParameter.java (DSS 6.5.RC1).
package policy

import "testing"

func TestCryptographicSuiteParameter_RoundTrip(t *testing.T) {
	p := NewCryptographicSuiteParameter()
	p.SetName("modulusLength")
	min := 2048
	max := 4096
	p.SetMin(&min)
	p.SetMax(&max)

	if p.Name() != "modulusLength" {
		t.Fatalf("Name() = %q, want modulusLength", p.Name())
	}
	if p.Min() == nil || *p.Min() != 2048 {
		t.Fatalf("Min() = %v, want 2048", p.Min())
	}
	if p.Max() == nil || *p.Max() != 4096 {
		t.Fatalf("Max() = %v, want 4096", p.Max())
	}

	cp := CryptographicSuiteParameterCopy(p)
	if !p.Equals(cp) {
		t.Fatalf("copy not equal to original")
	}
	// mutating the copy's pointee must not affect the original: copy
	// shares the *int per Java's field-by-field copy (Integer is
	// immutable in Java), so replacing the copy's pointer is the correct
	// mutation check.
	otherMax := 8192
	cp.SetMax(&otherMax)
	if p.Equals(cp) {
		t.Fatalf("mutated copy still equal to original")
	}
	if p.Max() == nil || *p.Max() != 4096 {
		t.Fatalf("original Max() mutated via copy: got %v", p.Max())
	}

	if CryptographicSuiteParameterCopy(nil) != nil {
		t.Fatalf("copy of nil must be nil")
	}
}

package model

import "testing"

func TestX500PrincipalHelperDelegatesTheNameForms(t *testing.T) {
	principal, err := NewX500Principal(mustHex(t, "3025310b30090603550406130246523116301406035504030c0d4a6f73c3a920416d706c69c3a9"))
	if err != nil {
		t.Fatal(err)
	}
	helper := NewX500PrincipalHelper(principal)

	if helper.Principal() != principal {
		t.Error("Principal() must return the wrapped principal")
	}
	if got := helper.RFC2253(); got != "CN=José Amplié,C=FR" {
		t.Errorf("RFC2253() = %q", got)
	}
	if got := helper.Canonical(); got != "cn=josé amplié,c=fr" {
		t.Errorf("Canonical() = %q", got)
	}
	if got := string(helper.Encoded()); got != string(principal.Encoded()) {
		t.Error("Encoded() must return the principal's DER")
	}
	// The pretty print substitutes the X520Attributes descriptions for the keywords.
	pretty, err := helper.PrettyPrintRFC2253()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := pretty, "commonName=José Amplié,countryName=FR"; got != want {
		t.Errorf("PrettyPrintRFC2253() = %q, want %q", got, want)
	}
}

func TestX500PrincipalHelperEqualsComparesEncodings(t *testing.T) {
	// The same name encoded as PrintableString and as UTF8String: X500Principal#equals sees
	// them as equal (canonical forms match), but X500PrincipalHelper#equals compares the raw
	// DER and must therefore say no.
	printable, err := NewX500Principal(mustHex(t, "3014311230100603550403130954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	utf8, err := NewX500Principal(mustHex(t, "30143112301006035504030c0954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	if !printable.Equals(utf8) {
		t.Fatal("the two principals should have equal canonical forms")
	}
	if NewX500PrincipalHelper(printable).Equals(NewX500PrincipalHelper(utf8)) {
		t.Error("the helper compares encodings, so these must not be equal")
	}
	if !NewX500PrincipalHelper(printable).Equals(NewX500PrincipalHelper(printable)) {
		t.Error("the helper must equal itself")
	}
	if NewX500PrincipalHelper(printable).Equals(nil) {
		t.Error("the helper must not equal nil")
	}
}

func TestX500PrincipalHelperPanicsOnNilPrincipal(t *testing.T) {
	defer func() {
		if r := recover(); r != "X500Principal cannot be null!" {
			t.Errorf("recover() = %v, want the Java requireNonNull message", r)
		}
	}()
	NewX500PrincipalHelper(nil)
	t.Error("a nil principal must panic")
}

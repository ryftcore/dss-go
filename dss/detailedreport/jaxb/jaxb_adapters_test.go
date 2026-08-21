package jaxb

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TestAdaptersRoundTrip is the exhaustive table test PORTING.md requires for
// registry-like tables: every enumeration constant an adapter binds must
// marshal to its Java lexical form and parse back to the same constant.
func TestAdaptersRoundTrip(t *testing.T) {
	t.Run("Indication", func(t *testing.T) {
		for _, c := range enumerations.IndicationValues() {
			roundTrip(t, IndicationValue(c), func(v IndicationValue) enumerations.Indication { return v.Indication() }, string(c))
		}
	})
	t.Run("SubIndication", func(t *testing.T) {
		for _, c := range enumerations.SubIndicationValues() {
			roundTrip(t, SubIndicationValue(c), func(v SubIndicationValue) enumerations.SubIndication { return v.SubIndication() }, string(c))
		}
	})
	t.Run("CertificateSourceType", func(t *testing.T) {
		for _, c := range enumerations.CertificateSourceTypeValues() {
			roundTrip(t, CertificateSourceTypeValue(c), func(v CertificateSourceTypeValue) enumerations.CertificateSourceType {
				return v.CertificateSourceType()
			}, string(c))
		}
	})
	t.Run("RevocationReason", func(t *testing.T) {
		for _, c := range enumerations.RevocationReasonValues() {
			roundTrip(t, RevocationReasonValue(c), func(v RevocationReasonValue) enumerations.RevocationReason { return v.RevocationReason() }, c.ShortName())
		}
	})
	t.Run("CertificateQualification", func(t *testing.T) {
		for _, c := range enumerations.CertificateQualificationValues() {
			roundTrip(t, CertificateQualificationValue(c), func(v CertificateQualificationValue) enumerations.CertificateQualification {
				return v.CertificateQualification()
			}, c.Readable())
		}
	})
	t.Run("SignatureQualification", func(t *testing.T) {
		for _, c := range enumerations.SignatureQualificationValues() {
			roundTrip(t, SignatureQualificationValue(c), func(v SignatureQualificationValue) enumerations.SignatureQualification {
				return v.SignatureQualification()
			}, c.Readable())
		}
	})
	t.Run("TimestampQualification", func(t *testing.T) {
		for _, c := range enumerations.TimestampQualificationValues() {
			roundTrip(t, TimestampQualificationValue(c), func(v TimestampQualificationValue) enumerations.TimestampQualification {
				return v.TimestampQualification()
			}, c.Readable())
		}
	})
	t.Run("EAAQualification", func(t *testing.T) {
		for _, c := range enumerations.EAAQualificationValues() {
			roundTrip(t, EAAQualificationValue(c), func(v EAAQualificationValue) enumerations.EAAQualification { return v.EAAQualification() }, c.Readable())
		}
	})
	t.Run("Context", func(t *testing.T) {
		for _, c := range enumerations.ContextValues() {
			roundTrip(t, ContextValue(c), func(v ContextValue) enumerations.Context { return v.Context() }, string(c))
		}
	})
	t.Run("ValidationTime", func(t *testing.T) {
		for _, c := range enumerations.ValidationTimeValues() {
			roundTrip(t, ValidationTimeValue(c), func(v ValidationTimeValue) enumerations.ValidationTime { return v.ValidationTime() }, string(c))
		}
	})
	t.Run("QWACProfile", func(t *testing.T) {
		for _, c := range enumerations.QWACProfileValues() {
			roundTrip(t, QWACProfileValue(c), func(v QWACProfileValue) enumerations.QWACProfile { return v.QWACProfile() }, c.Readable())
		}
	})
}

// roundTrip is generic over every "*Value" adapter type: it marshals v,
// checks the lexical form against wantText, then unmarshals that text back
// into a fresh value and checks it resolves to the same underlying constant.
func roundTrip[V interface {
	MarshalText() ([]byte, error)
}, E comparable](t *testing.T, v V, get func(V) E, wantText string) {
	t.Helper()
	want := get(v)
	text, err := v.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText(%v): %v", want, err)
	}
	if string(text) != wantText {
		t.Errorf("MarshalText(%v) = %q, want %q", want, text, wantText)
	}
	var parsed V
	pv, ok := any(&parsed).(interface{ UnmarshalText([]byte) error })
	if !ok {
		t.Fatalf("%T does not implement UnmarshalText", parsed)
	}
	if err := pv.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText(%q): %v", text, err)
	}
	if got := get(parsed); got != want {
		t.Errorf("UnmarshalText(%q) round-tripped to %v, want %v", text, got, want)
	}
}

// TestAdaptersRejectUnknown checks that an unrecognised lexical form is an
// error, not a silently-produced zero value, for every adapter.
func TestAdaptersRejectUnknown(t *testing.T) {
	const bogus = "NOT_A_REAL_VALUE"
	cases := []interface{ UnmarshalText([]byte) error }{
		new(IndicationValue),
		new(SubIndicationValue),
		new(CertificateSourceTypeValue),
		new(RevocationReasonValue),
		new(CertificateQualificationValue),
		new(SignatureQualificationValue),
		new(TimestampQualificationValue),
		new(EAAQualificationValue),
		new(ContextValue),
		new(ValidationTimeValue),
		new(QWACProfileValue),
	}
	for _, c := range cases {
		if err := c.UnmarshalText([]byte(bogus)); err == nil {
			t.Errorf("%T.UnmarshalText(%q) succeeded, want an error", c, bogus)
		}
	}
}

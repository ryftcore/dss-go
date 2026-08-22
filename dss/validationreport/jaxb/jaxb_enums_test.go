package jaxb

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
)

// TestNamespaceDSigProvenance pins namespaceDSig (jaxb_crossns.go) against
// dss/internal/xmldsig.NamespaceDSig, the constant already ported for the
// CAdES/XAdES signing stack (see jaxb_crossns.go's header: "reusing xml/
// common definitions where they exist ONLY for constants").
func TestNamespaceDSigProvenance(t *testing.T) {
	if namespaceDSig != xmldsig.NamespaceDSig {
		t.Fatalf("namespaceDSig = %q, want dss/internal/xmldsig.NamespaceDSig = %q", namespaceDSig, xmldsig.NamespaceDSig)
	}
}

// TestUriBasedEnumParserRoundTrip is an exhaustive table test (PORTING.md's
// "every registry-like table... gets an exhaustive table test") of
// UriBasedEnumParser: every value of every family it registers must parse
// back from its own URI, and an unknown URI must parse to the zero value
// for every family (mirroring UriBasedEnumParser.parse returning null).
func TestUriBasedEnumParserRoundTrip(t *testing.T) {
	for _, v := range enumerations.IndicationValues() {
		if got := ParseMainIndication(v.URI()); got != v {
			t.Errorf("ParseMainIndication(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range enumerations.SubIndicationValues() {
		if got := ParseSubIndication(v.URI()); got != v {
			t.Errorf("ParseSubIndication(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range enumerations.RevocationReasonValues() {
		if got := ParseRevocationReason(v.URI()); got != v {
			t.Errorf("ParseRevocationReason(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range ObjectTypeValues() {
		if got := ParseObjectType(v.URI()); got != v {
			t.Errorf("ParseObjectType(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range SignatureValidationProcessIDValues() {
		if got := ParseSignatureValidationProcessID(v.URI()); got != v {
			t.Errorf("ParseSignatureValidationProcessID(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range TypeOfProofValues() {
		if got := ParseTypeOfProof(v.URI()); got != v {
			t.Errorf("ParseTypeOfProof(%q) = %v, want %v", v.URI(), got, v)
		}
	}
	for _, v := range ConstraintStatusValues() {
		if got := ParseConstraintStatus(v.URI()); got != v {
			t.Errorf("ParseConstraintStatus(%q) = %v, want %v", v.URI(), got, v)
		}
	}

	// An unmatched URI parses to the zero value, mirroring parse() returning
	// null rather than throwing.
	if got := ParseMainIndication("urn:not-a-real-uri"); got != "" {
		t.Errorf("ParseMainIndication(unmatched) = %v, want zero value", got)
	}
	if got := ParseObjectType("urn:not-a-real-uri"); got != "" {
		t.Errorf("ParseObjectType(unmatched) = %v, want zero value", got)
	}

	// Print mirrors UriBasedEnumParser.print: the URI of a non-nil value, ""
	// for a nil interface.
	if got := Print(enumerations.IndicationTotalPassed); got != enumerations.IndicationTotalPassed.URI() {
		t.Errorf("Print(TOTAL_PASSED) = %q, want %q", got, enumerations.IndicationTotalPassed.URI())
	}
	if got := Print(nil); got != "" {
		t.Errorf("Print(nil) = %q, want \"\"", got)
	}
}

// TestValueEndorsementTypeRoundTrip is the exhaustive table test for the
// EndorsementTypeParser adapter (value(), not URI, per jaxb_enums.go's
// header).
func TestValueEndorsementTypeRoundTrip(t *testing.T) {
	for _, v := range enumerations.EndorsementTypeValues() {
		w := ValueEndorsementType(v)
		text, err := w.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v): %v", v, err)
		}
		if string(text) != v.Value() {
			t.Errorf("MarshalText(%v) = %q, want %q", v, text, v.Value())
		}
		var back ValueEndorsementType
		if err := back.UnmarshalText(text); err != nil {
			t.Fatalf("UnmarshalText(%q): %v", text, err)
		}
		if back.EndorsementType() != v {
			t.Errorf("round-trip %v -> %q -> %v", v, text, back.EndorsementType())
		}
	}
}

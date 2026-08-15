package jaxb

import (
	"encoding/xml"
	"strings"
	"testing"
)

// TestDSPrefixOverridesFieldTag pins the encoding/xml behaviour
// jaxb_crossns.go's header relies on: a field value's own MarshalXML fully
// determines the element it writes, including its name, regardless of what
// name its enclosing struct's field tag would otherwise have produced. This
// is what lets DigestMethodType/DSDigestValue/SignatureValueType hard-code
// their "ns2:" element name once, rather than needing every type that
// embeds one of them (SACertIDType, SACRLIDType, DigestAlgAndValueType,
// SignatureIdentifierType) to marshal itself by hand too.
func TestDSPrefixOverridesFieldTag(t *testing.T) {
	v := DigestAlgAndValueType{
		DigestMethod: DigestMethodType{Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256"},
		DigestValue:  DSDigestValue{0x01, 0x02},
	}
	b, err := xml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, "<ns2:DigestMethod ") {
		t.Errorf("marshalled output does not use ns2:DigestMethod: %s", got)
	}
	if !strings.Contains(got, "<ns2:DigestValue>") {
		t.Errorf("marshalled output does not use ns2:DigestValue: %s", got)
	}
	if strings.Contains(got, `xmlns="http://www.w3.org/2000/09/xmldsig#"`) {
		t.Errorf("marshalled output re-declares the xmldsig namespace locally instead of using the ns2: prefix: %s", got)
	}
}

// KAT test for xades111: every XAdES111Attribute_* attribute name below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdES111Attribute_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ID", XAdES111AttributeID.AttributeName(), "Id"},
		{"OBJECT_REFERENCE", XAdES111AttributeObjectReference.AttributeName(), "ObjectReference"},
		{"QUALIFIER", XAdES111AttributeQualifier.AttributeName(), "Qualifier"},
		{"TARGET", XAdES111AttributeTarget.AttributeName(), "Target"},
		{"URI", XAdES111AttributeURI.AttributeName(), "uri"},
		{"URI2", XAdES111AttributeURI2.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

// KAT test for xades122: every XAdES122Attribute_* attribute name below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdES122Attribute_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ID", XAdES122AttributeID.AttributeName(), "Id"},
		{"OBJECT_REFERENCE", XAdES122AttributeObjectReference.AttributeName(), "ObjectReference"},
		{"QUALIFIER", XAdES122AttributeQualifier.AttributeName(), "Qualifier"},
		{"REFERENCED_DATA", XAdES122AttributeReferencedData.AttributeName(), "referencedData"},
		{"TARGET", XAdES122AttributeTarget.AttributeName(), "Target"},
		{"URI", XAdES122AttributeURI.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

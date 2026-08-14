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
		{"ID", XAdES122Attribute_ID.AttributeName(), "Id"},
		{"OBJECT_REFERENCE", XAdES122Attribute_OBJECT_REFERENCE.AttributeName(), "ObjectReference"},
		{"QUALIFIER", XAdES122Attribute_QUALIFIER.AttributeName(), "Qualifier"},
		{"REFERENCED_DATA", XAdES122Attribute_REFERENCED_DATA.AttributeName(), "referencedData"},
		{"TARGET", XAdES122Attribute_TARGET.AttributeName(), "Target"},
		{"URI", XAdES122Attribute_URI.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

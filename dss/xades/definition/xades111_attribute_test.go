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
		{"ID", XAdES111Attribute_ID.AttributeName(), "Id"},
		{"OBJECT_REFERENCE", XAdES111Attribute_OBJECT_REFERENCE.AttributeName(), "ObjectReference"},
		{"QUALIFIER", XAdES111Attribute_QUALIFIER.AttributeName(), "Qualifier"},
		{"TARGET", XAdES111Attribute_TARGET.AttributeName(), "Target"},
		{"URI", XAdES111Attribute_URI.AttributeName(), "uri"},
		{"URI2", XAdES111Attribute_URI2.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

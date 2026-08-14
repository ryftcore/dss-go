// KAT test for xades132: every XAdES132Attribute_* attribute name below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdES132Attribute_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ENCODING", XAdES132Attribute_ENCODING.AttributeName(), "Encoding"},
		{"ID", XAdES132Attribute_ID.AttributeName(), "Id"},
		{"OBJECT_REFERENCE", XAdES132Attribute_OBJECT_REFERENCE.AttributeName(), "ObjectReference"},
		{"QUALIFIER", XAdES132Attribute_QUALIFIER.AttributeName(), "Qualifier"},
		{"REFERENCED_DATA", XAdES132Attribute_REFERENCED_DATA.AttributeName(), "referencedData"},
		{"TARGET", XAdES132Attribute_TARGET.AttributeName(), "Target"},
		{"URI", XAdES132Attribute_URI.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

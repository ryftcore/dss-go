// KAT test for xades141: every XAdES141Attribute_* attribute name below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdES141Attribute_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ID", XAdES141AttributeID.AttributeName(), "Id"},
		{"ORDER", XAdES141AttributeOrder.AttributeName(), "Order"},
		{"URI", XAdES141AttributeURI.AttributeName(), "URI"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.AttributeName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

// Ported from java.text.MessageFormat semantics as exercised by
// dss-messages.properties (DSS 6.5.RC1). See message_format.go's file
// header for scope and the apostrophe-quoting quirk this pins down.
package i18n

import "testing"

func TestMessageFormatPlaceholders(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		args    []any
		want    string
	}{
		{"no placeholders", "Status", nil, "Status"},
		{"single placeholder", "Status : {0}", []any{"granted"}, "Status : granted"},
		{"multiple placeholders", "{0}-{1}-{2}", []any{"a", "b", "c"}, "a-b-c"},
		{"missing arg left verbatim", "Status : {0}", nil, "Status : {0}"},
		{"doubled quote escapes to literal", "It''s fine", nil, "It's fine"},
		{"quoted literal braces", "a'{0}'b", []any{"X"}, "a{0}b"},
		{"bare apostrophe swallows itself", "Responder's path", nil, "Responders path"},
		{"repeated index", "{0} and {0}", []any{"x"}, "x and x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := messageFormat(tc.pattern, tc.args); got != tc.want {
				t.Errorf("messageFormat(%q, %v) = %q, want %q", tc.pattern, tc.args, got, tc.want)
			}
		})
	}
}

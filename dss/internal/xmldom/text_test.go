package xmldom

import "testing"

// TestTextPreservation covers the rules of the design note's whitespace, CDATA and
// entity section: nothing is stripped, nothing is coalesced, and a character
// reference and the character it denotes are indistinguishable in a text node.
func TestTextPreservation(t *testing.T) {
	tests := []struct {
		name, src string
		want      []struct {
			kind Kind
			data string
		}
	}{{
		name: "adjacent text and CDATA are separate nodes",
		src:  `<r>a<![CDATA[b]]>c</r>`,
		want: []struct {
			kind Kind
			data string
		}{{Text, "a"}, {CDATA, "b"}, {Text, "c"}},
	}, {
		name: "CDATA content is stored verbatim",
		src:  `<r><![CDATA[a < b & c > d "e" 'f' &amp;]]></r>`,
		want: []struct {
			kind Kind
			data string
		}{{CDATA, `a < b & c > d "e" 'f' &amp;`}},
	}, {
		name: "a split CDATA end marker rejoins as one run per section",
		src:  `<r><![CDATA[x]]]]><![CDATA[>y]]></r>`,
		want: []struct {
			kind Kind
			data string
		}{{CDATA, "x]]"}, {CDATA, ">y"}},
	}, {
		name: "entities are decoded and not split",
		src:  `<r>a&amp;b&#65;c</r>`,
		want: []struct {
			kind Kind
			data string
		}{{Text, "a&bAc"}},
	}, {
		name: "a CR reference survives clause 2.11",
		src:  `<r>a&#13;b</r>`,
		want: []struct {
			kind Kind
			data string
		}{{Text, "a\rb"}},
	}, {
		name: "literal line endings are normalized",
		src:  "<r>a\r\nb\rc\nd</r>",
		want: []struct {
			kind Kind
			data string
		}{{Text, "a\nb\nc\nd"}},
	}, {
		name: "whitespace-only text is kept",
		src:  "<r>\n  \t</r>",
		want: []struct {
			kind Kind
			data string
		}{{Text, "\n  \t"}},
	}, {
		name: "text is not coalesced across a comment",
		src:  `<r>a<!--x-->b</r>`,
		want: []struct {
			kind Kind
			data string
		}{{Text, "a"}, {Comment, "x"}, {Text, "b"}},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			children := mustParseRoot(t, tc.src).Children()
			if len(children) != len(tc.want) {
				t.Fatalf("got %d children, want %d: %s", len(children), len(tc.want), dump(mustParseRoot(t, tc.src)))
			}
			for i, w := range tc.want {
				if children[i].Kind != w.kind || children[i].Value != w.data {
					t.Errorf("child %d = (%s, %q), want (%s, %q)", i, children[i].Kind, children[i].Value, w.kind, w.data)
				}
			}
		})
	}
}

// TestXMLSpaceIsAnOrdinaryAttribute pins that xml:space is stored and never acted on.
func TestXMLSpaceIsAnOrdinaryAttribute(t *testing.T) {
	el := mustParseRoot(t, `<r xml:space="preserve">  a  </r>`)
	if got := el.AttrValue(XMLNamespace, "space"); got != "preserve" {
		t.Errorf("xml:space = %q, want %q", got, "preserve")
	}
	if got := el.TextContent(); got != "  a  " {
		t.Errorf("text = %q, want %q; xml:space must not be interpreted", got, "  a  ")
	}
	// The same document without xml:space must keep its whitespace too.
	if got := mustParseRoot(t, `<r>  a  </r>`).TextContent(); got != "  a  " {
		t.Errorf("text = %q, want %q", got, "  a  ")
	}
}

func TestCommentAndProcInstData(t *testing.T) {
	// Comments and PIs have no markup or reference recognition, so &#13; is six
	// literal characters, not a CR.
	root := mustParseRoot(t, `<r><!--a&#13;b--><?t d&#13;e?></r>`)
	kids := root.Children()
	if len(kids) != 2 {
		t.Fatalf("got %d children, want 2", len(kids))
	}
	if kids[0].Kind != Comment || kids[0].Value != "a&#13;b" {
		t.Errorf("comment = (%s, %q), want (Comment, %q)", kids[0].Kind, kids[0].Value, "a&#13;b")
	}
	if kids[1].Kind != ProcInst || kids[1].Name.Local != "t" || kids[1].Value != "d&#13;e" {
		t.Errorf("PI = (%s, target %q, %q), want (ProcInst, %q, %q)",
			kids[1].Kind, kids[1].Name.Local, kids[1].Value, "t", "d&#13;e")
	}
}

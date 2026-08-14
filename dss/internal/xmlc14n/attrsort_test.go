package xmlc14n

import (
	"bufio"
	"bytes"
	"testing"
)

// emit renders a set the way the emitters do, so the expectations below can be read straight
// off the checked-in goldens.
func emit(t *testing.T, s *attrSet) string {
	t.Helper()
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	s.writeTo(w)
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	return buf.String()
}

func plainAttr(space, local, qname, value string) outAttr {
	return outAttr{space: space, local: local, qname: qname, value: value}
}

func TestAttrSortOrder(t *testing.T) {
	// corpus/attr-order.xml, in document order:
	//   <r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d" Id="r"/>
	// golden/attr-order.c14n10.doc, produced by Santuario:
	//   <r xmlns="urn:d" xmlns:a="urn:a" xmlns:b="urn:b" Id="r" a="4" z="3" a:z="2" b:z="1">
	var set attrSet
	set.add(nsAttr("b", "urn:b"))
	set.add(nsAttr("a", "urn:a"))
	set.add(plainAttr("urn:b", "z", "b:z", "1"))
	set.add(plainAttr("urn:a", "z", "a:z", "2"))
	set.add(plainAttr("", "z", "z", "3"))
	set.add(plainAttr("", "a", "a", "4"))
	set.add(nsAttr(xmlnsPrefix, "urn:d"))
	set.add(plainAttr("", "Id", "Id", "r"))

	want := ` xmlns="urn:d" xmlns:a="urn:a" xmlns:b="urn:b" Id="r" a="4" z="3" a:z="2" b:z="1"`
	if got := emit(t, &set); got != want {
		t.Errorf("emitted\n  %q\nwant\n  %q", got, want)
	}
}

func TestAttrSortNamespacedAfterNoNamespace(t *testing.T) {
	// An xml:* attribute is namespaced, so it sorts after every no-namespace attribute:
	// golden/xmlattrs-lang-space.c14n10.apex-t is
	//   <t Id="t" a="1" xml:lang="en" xml:space="preserve"></t>
	var set attrSet
	set.add(plainAttr(xmlNamespace, "space", "xml:space", "preserve"))
	set.add(plainAttr(xmlNamespace, "lang", "xml:lang", "en"))
	set.add(plainAttr("", "a", "a", "1"))
	set.add(plainAttr("", "Id", "Id", "t"))

	want := ` Id="t" a="1" xml:lang="en" xml:space="preserve"`
	if got := emit(t, &set); got != want {
		t.Errorf("emitted\n  %q\nwant\n  %q", got, want)
	}
}

func TestAttrSetIsASet(t *testing.T) {
	// Santuario collects into a TreeSet, so an element that compares equal to one already
	// present is discarded and the first one wins.
	var set attrSet
	set.add(plainAttr("", "a", "a", "first"))
	set.add(plainAttr("", "a", "a", "second"))
	if set.len() != 1 {
		t.Fatalf("len = %d, want 1", set.len())
	}
	if got := emit(t, &set); got != ` a="first"` {
		t.Errorf("emitted %q, want the first insertion to win", got)
	}
}

func TestCompareUTF16(t *testing.T) {
	// java.lang.String.compareTo compares UTF-16 code units, so a code point above U+FFFF -
	// a surrogate pair starting at U+D800 - sorts BELOW U+E000..U+FFFD, the opposite of Go's
	// byte-wise order. golden/sort-astral.c14n10.doc pins it:
	//   <r xmlns:a="urn:xﷰ" xmlns:b="urn:x\U00010000" Id="r" b:z="2" a:z="1"></r>
	astral, bmp := "urn:x\U00010000", "urn:xﷰ"
	if compareUTF16(astral, bmp) >= 0 {
		t.Errorf("compareUTF16(astral, U+FDF0) = %d, want negative (Java order)", compareUTF16(astral, bmp))
	}
	if astral <= bmp {
		t.Fatal("test premise broken: Go byte order should put the astral URI second")
	}

	var set attrSet
	set.add(nsAttr("a", bmp))
	set.add(nsAttr("b", astral))
	set.add(plainAttr(bmp, "z", "a:z", "1"))
	set.add(plainAttr(astral, "z", "b:z", "2"))
	set.add(plainAttr("", "Id", "Id", "r"))
	want := ` xmlns:a="urn:x` + "ﷰ" + `" xmlns:b="urn:x` + "\U00010000" + `" Id="r" b:z="2" a:z="1"`
	if got := emit(t, &set); got != want {
		t.Errorf("emitted\n  %q\nwant\n  %q", got, want)
	}

	for _, tc := range []struct {
		a, b string
		want int // sign
	}{
		{"", "", 0},
		{"a", "a", 0},
		{"a", "b", -1},
		{"b", "a", 1},
		{"a", "ab", -1},
		{"ab", "a", 1},
		{"café", "café", 0},
		{"caf", "café", -1},
		{"�", "\U00010000", 1}, // UTF-16: FFFD > D800
		{"́", "ﷰ", -1},
		{"x\U00010000", "x\U00010001", -1},
	} {
		got := compareUTF16(tc.a, tc.b)
		if (got < 0) != (tc.want < 0) || (got > 0) != (tc.want > 0) {
			t.Errorf("compareUTF16(%q, %q) = %d, want sign %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestNSAttrQName(t *testing.T) {
	if got := nsAttr(xmlnsPrefix, "urn:d"); got.qname != "xmlns" || got.local != "xmlns" || got.space != xmlnsNamespace {
		t.Errorf("nsAttr(xmlns) = %+v", got)
	}
	if got := nsAttr("p", "urn:1"); got.qname != "xmlns:p" || got.local != "p" {
		t.Errorf("nsAttr(p) = %+v", got)
	}
}

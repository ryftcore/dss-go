// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against
// dss-utils-apache-commons/.../ApacheCommonsUtilsTest.java and known
// commons-lang3 StringUtils semantics.

package utils

import "testing"

func TestIsStringEmpty(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"nowina", false},
		{"   ", false}, // blank is NOT empty
	}
	for _, c := range cases {
		if got := IsStringEmpty(c.in); got != c.want {
			t.Errorf("IsStringEmpty(%q) = %v, want %v", c.in, got, c.want)
		}
		if got := IsStringNotEmpty(c.in); got != !c.want {
			t.Errorf("IsStringNotEmpty(%q) = %v, want %v", c.in, got, !c.want)
		}
	}
}

func TestAreAllStringsEmptyAndAtLeastOne(t *testing.T) {
	if !AreAllStringsEmpty() {
		t.Error("AreAllStringsEmpty() with no args should be true")
	}
	if !AreAllStringsEmpty("", "") {
		t.Error("AreAllStringsEmpty(\"\", \"\") should be true")
	}
	if AreAllStringsEmpty("", "x") {
		t.Error("AreAllStringsEmpty(\"\", \"x\") should be false")
	}
	if IsAtLeastOneStringNotEmpty("", "") {
		t.Error("IsAtLeastOneStringNotEmpty(\"\", \"\") should be false")
	}
	if !IsAtLeastOneStringNotEmpty("", "x") {
		t.Error("IsAtLeastOneStringNotEmpty(\"\", \"x\") should be true")
	}
}

func TestIsStringBlank(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n\r", true},
		{"nowina", false},
		{" a ", false},
		{"\u00a0", false}, // NBSP: Java's Character.isWhitespace excludes it
		{"\u2003", true},  // EM SPACE (Zs): Java DOES consider it whitespace
	}
	for _, c := range cases {
		if got := IsStringBlank(c.in); got != c.want {
			t.Errorf("IsStringBlank(%q) = %v, want %v", c.in, got, c.want)
		}
		if got := IsStringNotBlank(c.in); got != !c.want {
			t.Errorf("IsStringNotBlank(%q) = %v, want %v", c.in, got, !c.want)
		}
	}
}

func TestAreStringsEqual(t *testing.T) {
	if !AreStringsEqual("nowina", "nowina") {
		t.Error("expected equal")
	}
	if AreStringsEqual("nowina", "Nowina") {
		t.Error("expected not equal (case-sensitive)")
	}
	if !AreStringsEqual("", "") {
		t.Error("empty strings should be equal")
	}
}

func TestAreStringsEqualIgnoreCase(t *testing.T) {
	if !AreStringsEqualIgnoreCase("nowina", "Nowina") {
		t.Error("expected equal ignoring case")
	}
	if AreStringsEqualIgnoreCase("water", "fire") {
		t.Error("expected not equal")
	}
}

func TestIsStringDigits(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"123", true},
		{"1a2b", false},
		{"", false},
		{"-123", false},
		{"12.3", false},
	}
	for _, c := range cases {
		if got := IsStringDigits(c.in); got != c.want {
			t.Errorf("IsStringDigits(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestTrim(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"   123 ", "123"},
		{"", ""},
		{"   ", ""},
		{"\tabc\n", "abc"},
		{"no-trim-needed", "no-trim-needed"},
	}
	for _, c := range cases {
		if got := Trim(c.in); got != c.want {
			t.Errorf("Trim(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestJoinStrings(t *testing.T) {
	if got := JoinStrings([]string{"Nowina", "123"}, ","); got != "Nowina,123" {
		t.Errorf("got %q", got)
	}
	if got := JoinStrings(nil, ","); got != "" {
		t.Errorf("JoinStrings(nil, \",\") = %q, want \"\"", got)
	}
	if got := JoinStrings([]string{"a"}, ","); got != "a" {
		t.Errorf("got %q", got)
	}
}

func TestSubstringAfter(t *testing.T) {
	cases := []struct {
		text, after, want string
	}{
		{"aaaaa?bbb", "?", "bbb"},
		{"", "?", ""},
		{"aaaaa", "?", ""},
		{"aaaaa?bbb?ccc", "?", "bbb?ccc"},
		{"abc", "", "abc"},
	}
	for _, c := range cases {
		if got := SubstringAfter(c.text, c.after); got != c.want {
			t.Errorf("SubstringAfter(%q, %q) = %q, want %q", c.text, c.after, got, c.want)
		}
	}
}

func TestEndsWithIgnoreCase(t *testing.T) {
	cases := []struct {
		text, expected string
		want           bool
	}{
		{"hello", "LO", true},
		{"hello", "a", false},
		{"hello", "", true},
		{"hi", "hello", false},
	}
	for _, c := range cases {
		if got := EndsWithIgnoreCase(c.text, c.expected); got != c.want {
			t.Errorf("EndsWithIgnoreCase(%q, %q) = %v, want %v", c.text, c.expected, got, c.want)
		}
	}
}

func TestLowerUpperCase(t *testing.T) {
	if got := LowerCase("Nowina"); got != "nowina" {
		t.Errorf("got %q", got)
	}
	if got := UpperCase("Nowina"); got != "NOWINA" {
		t.Errorf("got %q", got)
	}
}

func TestGetFileNameExtension(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"file.xml", "xml"},
		{"document.pdf", "pdf"},
		{"noext", ""},
		{"archive.tar.gz", "gz"},
		{".hidden", "hidden"},
		{"dir.with.dot/file", ""},
		{"dir.with.dot/file.txt", "txt"},
		{"", ""},
	}
	for _, c := range cases {
		if got := GetFileNameExtension(c.in); got != c.want {
			t.Errorf("GetFileNameExtension(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsTrue(t *testing.T) {
	tr := true
	fa := false
	if IsTrue(nil) {
		t.Error("IsTrue(nil) should be false")
	}
	if !IsTrue(&tr) {
		t.Error("IsTrue(&true) should be true")
	}
	if IsTrue(&fa) {
		t.Error("IsTrue(&false) should be false")
	}
}

// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1) test
// vectors, cross-checked against known org.apache.commons.codec.binary.Hex
// / Base64 semantics.

package utils

import (
	"reflect"
	"testing"
)

func TestIsHexEncoded(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"deadbeef", true},
		{"DEADBEEF", true},
		{"aBcD", true},
		{"abc", false}, // odd length
		{"zz", false},  // invalid digit
		{"g0", false},
	}
	for _, c := range cases {
		if got := IsHexEncoded(c.in); got != c.want {
			t.Errorf("IsHexEncoded(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToHexFromHex(t *testing.T) {
	if got := ToHex([]byte{0xDE, 0xAD, 0xBE, 0xEF}); got != "deadbeef" {
		t.Errorf("got %q", got)
	}
	if got := ToHex(nil); got != "" {
		t.Errorf("ToHex(nil) = %q, want \"\"", got)
	}
	if got := ToHex([]byte{}); got != "" {
		t.Errorf("ToHex([]) = %q, want \"\"", got)
	}

	got, err := FromHex("deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, []byte{0xDE, 0xAD, 0xBE, 0xEF}) {
		t.Errorf("got %v", got)
	}

	if _, err := FromHex("abc"); err == nil {
		t.Error("FromHex(odd length) should error")
	}
	if _, err := FromHex("zz"); err == nil {
		t.Error("FromHex(invalid digit) should error")
	}
}

func TestIsBase64Encoded(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"aGVsbG8=", true},
		{"aGVsbG8", true},       // isBase64 is a charset check only, no padding validation
		{"aGVsbG8= ", true},     // trailing whitespace is ignorable
		{"hello world!", false}, // '!' and possibly space handling: '!' invalid
		{"not_base64", false},   // '_' not in standard alphabet
	}
	for _, c := range cases {
		if got := IsBase64Encoded(c.in); got != c.want {
			t.Errorf("IsBase64Encoded(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToBase64FromBase64(t *testing.T) {
	if got := ToBase64([]byte("hello")); got != "aGVsbG8=" {
		t.Errorf("got %q", got)
	}
	if got := ToBase64(nil); got != "" {
		t.Errorf("ToBase64(nil) = %q, want \"\"", got)
	}

	if got := FromBase64("aGVsbG8="); string(got) != "hello" {
		t.Errorf("got %q", got)
	}
	if got := FromBase64(""); len(got) != 0 {
		t.Errorf("FromBase64(\"\") = %v, want empty", got)
	}
	// Lenient decoding: embedded whitespace is ignored, not an error.
	if got := FromBase64("aGVs\nbG8="); string(got) != "hello" {
		t.Errorf("got %q, want \"hello\"", got)
	}
	// Lenient decoding: unpadded input still decodes.
	if got := FromBase64("aGVsbG8"); string(got) != "hello" {
		t.Errorf("got %q, want \"hello\"", got)
	}
}

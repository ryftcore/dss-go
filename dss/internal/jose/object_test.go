package jose

import (
	"encoding/base64"
	"strings"
	"testing"
)

// TestObjectKeepsInsertionOrder is the property the whole package rests on: a Go map cannot
// provide it, so if this ever regresses every signature this library produces changes.
func TestObjectKeepsInsertionOrder(t *testing.T) {
	o := NewObject()
	// Reverse alphabetical, so a sorted or randomised iteration would be obvious.
	for _, k := range []string{"zeta", "yankee", "x-ray", "whiskey", "victor", "alpha"} {
		o.Put(k, k)
	}
	want := "zeta,yankee,x-ray,whiskey,victor,alpha"
	if got := strings.Join(o.Keys(), ","); got != want {
		t.Errorf("keys = %s, want %s", got, want)
	}
	// Repeat it enough times that Go's map randomisation would have shown up.
	for i := 0; i < 50; i++ {
		if got := strings.Join(o.Keys(), ","); got != want {
			t.Fatalf("iteration %d: keys = %s, want %s", i, got, want)
		}
	}
}

// TestObjectRePutKeepsPosition covers LinkedHashMap's insertion-order (not access-order)
// contract, which decides where a header parameter that is set twice ends up.
func TestObjectRePutKeepsPosition(t *testing.T) {
	o := NewObject()
	o.Put("a", 1)
	o.Put("b", 2)
	o.Put("c", 3)
	previous := o.Put("a", 9)
	if previous != 1 {
		t.Errorf("Put returned %v, want the previous value 1", previous)
	}
	if got := strings.Join(o.Keys(), ","); got != "a,b,c" {
		t.Errorf("keys = %s, want a,b,c", got)
	}
}

// TestObjectRemoveAndReplace covers the two mutators DSS uses on a parsed header:
// JWSConverter replaces the 'etsiU' array in place, and EtsiUHeader removes components.
func TestObjectRemoveAndReplace(t *testing.T) {
	o := NewObjectFromPairs("a", 1, "b", 2, "c", 3)

	if got := o.Remove("b"); got != 2 {
		t.Errorf("Remove returned %v, want 2", got)
	}
	if got := strings.Join(o.Keys(), ","); got != "a,c" {
		t.Errorf("after Remove, keys = %s, want a,c", got)
	}
	if got := o.Remove("missing"); got != nil {
		t.Errorf("removing an absent key returned %v", got)
	}

	if got := o.Replace("a", 9); got != 1 {
		t.Errorf("Replace returned %v, want 1", got)
	}
	if got := o.Value("a"); got != 9 {
		t.Errorf("after Replace, a = %v, want 9", got)
	}
	// Replace on an absent key must not add it - that is what distinguishes it from Put.
	if got := o.Replace("absent", 1); got != nil {
		t.Errorf("replacing an absent key returned %v", got)
	}
	if o.ContainsKey("absent") {
		t.Error("Replace added an absent key")
	}
}

// TestNilObjectReadsAsEmpty keeps the read paths total, the way a Java null map reference is
// not - DSS calls getUnprotected() and reads straight through the result in several places.
func TestNilObjectReadsAsEmpty(t *testing.T) {
	var o *Object
	if o.Size() != 0 || !o.IsEmpty() || o.ContainsKey("x") || o.Value("x") != nil || o.Keys() != nil {
		t.Error("a nil *Object did not read as empty")
	}
}

// TestObjectCloneIsShallowAndIndependent covers the copy JWS plumbing relies on.
func TestObjectCloneIsShallowAndIndependent(t *testing.T) {
	original := NewObjectFromPairs("a", 1, "b", 2)
	clone := original.Clone()
	clone.Put("c", 3)
	if original.ContainsKey("c") {
		t.Error("mutating the clone changed the original")
	}
	if got := strings.Join(clone.Keys(), ","); got != "a,b,c" {
		t.Errorf("clone keys = %s", got)
	}
	if NewHashObject().Clone().IsHashOrdered() != true {
		t.Error("Clone lost the ordering mode")
	}
}

// TestBase64URLEncodeMatchesStdlib checks the hand-written encoder against
// encoding/base64.RawURLEncoding, which is the one place the two are expected to agree exactly.
func TestBase64URLEncodeMatchesStdlib(t *testing.T) {
	for length := 0; length < 64; length++ {
		input := make([]byte, length)
		for i := range input {
			input[i] = byte(i*7 + length)
		}
		if got, want := Base64URLEncode(input), base64.RawURLEncoding.EncodeToString(input); got != want {
			t.Fatalf("length %d: got %q, want %q", length, got, want)
		}
	}
}

// TestBase64URLDecodeIsLenient states the leniency explicitly, because
// DSSJsonUtils.isBase64UrlEncoded is built on top of it and a strict decoder would change which
// 'etsiU' components DSS treats as base64url-encoded.
func TestBase64URLDecodeIsLenient(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"padding is optional", "Zm9v", "foo"},
		{"padding is accepted", "Zm8=", "fo"},
		{"the standard alphabet is accepted", "+/8=", "\xfb\xff"},
		{"the url-safe alphabet is accepted", "-_8", "\xfb\xff"},
		{"a space is skipped", "Zm 9v", "foo"},
		{"a newline is skipped", "Zm9v\n", "foo"},
		{"a period is skipped", "Zm9v.", "foo"},
		{"non-alphabet characters decode to nothing", "!!!!", ""},
		{"a single character has no whole byte", "A", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := string(Base64URLDecode(tc.input)); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// TestASCIIBytesReplacesNonASCII covers Java's US-ASCII encoder, whose '?' substitution decides
// what gets signed if a header ever contains a non-ASCII character.
func TestASCIIBytesReplacesNonASCII(t *testing.T) {
	if got := string(ASCIIBytes("aéb")); got != "a?b" {
		t.Errorf("got %q, want %q", got, "a?b")
	}
	if got := string(ASCIIBytes("a中b")); got != "a?b" {
		t.Errorf("got %q, want %q", got, "a?b")
	}
	if got := string(ASCIIBytes("plain.ASCII-123_")); got != "plain.ASCII-123_" {
		t.Errorf("ASCII was not passed through: %q", got)
	}
}

// TestUTF8StringReplacesMalformedSequences covers Java's new String(bytes, UTF_8), which is how
// a corrupt protected header reaches the parser.
func TestUTF8StringReplacesMalformedSequences(t *testing.T) {
	if got := UTF8String([]byte{0x66, 0xff, 0x6f}); got != "f�o" {
		t.Errorf("got %q, want %q", got, "f�o")
	}
	if got := UTF8String([]byte("café")); got != "café" {
		t.Errorf("valid UTF-8 was altered: %q", got)
	}
}

// TestJavaStringHashCode pins the hash function the bucket order depends on, against values
// taken from the Java specification's own worked examples.
func TestJavaStringHashCode(t *testing.T) {
	tests := map[string]int32{
		"":         0,
		"a":        97,
		"hello":    99162322,
		"Aa":       2112,
		"BB":       2112, // the classic collision
		"alg":      96668,
		"crlVals":  1034673391,
		"ocspVals": 1245415491,
		"x5t#S256": 1832724684,
	}
	for input, want := range tests {
		if got := JavaStringHashCode(input); got != want {
			t.Errorf("JavaStringHashCode(%q) = %d, want %d", input, got, want)
		}
	}
}

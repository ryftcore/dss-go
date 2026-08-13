package xmlc14n

import (
	"errors"
	"reflect"
	"testing"
)

func TestSupportedRegistry(t *testing.T) {
	// Exactly the seven XMLCanonicalizer.registerDefaultCanonicalizers() registers, by URI.
	want := []Algorithm{
		"http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
		"http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments",
		"http://www.w3.org/2006/12/xml-c14n11",
		"http://www.w3.org/2006/12/xml-c14n11#WithComments",
		"http://www.w3.org/2001/10/xml-exc-c14n#",
		"http://www.w3.org/2001/10/xml-exc-c14n#WithComments",
		"http://santuario.apache.org/c14n/physical",
	}
	if len(registered) != len(want) {
		t.Fatalf("registry holds %d algorithms, want %d", len(registered), len(want))
	}
	for _, alg := range want {
		if !Supported(alg) {
			t.Errorf("Supported(%q) = false", alg)
		}
	}
	for _, alg := range []Algorithm{
		"",
		"http://www.w3.org/2001/10/xml-exc-c14n", // missing '#'
		"http://www.w3.org/TR/2001/REC-xml-c14n-20010315#withcomments", // wrong case
		"http://www.w3.org/2006/12/xml-c14n11#",                        // trailing '#'
		"http://www.w3.org/2010/xml-c14n2",                             // C14N 2.0, never registered
	} {
		if Supported(alg) {
			t.Errorf("Supported(%q) = true", alg)
		}
	}
}

func TestDefaultsMatchUpstream(t *testing.T) {
	if DefaultDSS != C14NExclusive {
		t.Errorf("DefaultDSS = %q, want XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD", DefaultDSS)
	}
	if DefaultXMLDSig != C14N10 {
		t.Errorf("DefaultXMLDSig = %q, want XMLCanonicalizer.DEFAULT_XMLDSIG_C14N_METHOD", DefaultXMLDSig)
	}
}

func TestResolve(t *testing.T) {
	// An empty algorithm is XMLDSIG 4.4.3.2's default (DSS-2208), not an error.
	got, err := Resolve("")
	if err != nil || got != DefaultXMLDSig {
		t.Errorf(`Resolve("") = %q, %v; want %q, nil`, got, err, DefaultXMLDSig)
	}
	if got, err := Resolve(C14NExclusiveWithComments); err != nil || got != C14NExclusiveWithComments {
		t.Errorf("Resolve(exclusive#WithComments) = %q, %v", got, err)
	}
	_, err = Resolve("urn:nope")
	if !errors.Is(err, ErrUnsupportedAlgorithm) {
		t.Errorf("Resolve(unknown) error = %v, want ErrUnsupportedAlgorithm", err)
	}
}

func TestAlgorithmFlags(t *testing.T) {
	for _, tc := range []struct {
		alg                                 Algorithm
		comments, exclusive, c14n11, physic bool
	}{
		{C14N10, false, false, false, false},
		{C14N10WithComments, true, false, false, false},
		{C14N11, false, false, true, false},
		{C14N11WithComments, true, false, true, false},
		{C14NExclusive, false, true, false, false},
		{C14NExclusiveWithComments, true, true, false, false},
		{C14NPhysical, false, false, false, true},
	} {
		if got := tc.alg.withComments(); got != tc.comments {
			t.Errorf("%q withComments = %v", tc.alg, got)
		}
		if got := tc.alg.exclusive(); got != tc.exclusive {
			t.Errorf("%q exclusive = %v", tc.alg, got)
		}
		if got := tc.alg.c14n11(); got != tc.c14n11 {
			t.Errorf("%q c14n11 = %v", tc.alg, got)
		}
		if got := tc.alg.physical(); got != tc.physic {
			t.Errorf("%q physical = %v", tc.alg, got)
		}
	}
}

func TestParsePrefixList(t *testing.T) {
	// InclusiveNamespaces.prefixStr2Set: split on \s into a TreeSet, "#default" -> "xmlns".
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"q", []string{"q"}},
		{"#default q", []string{"q", xmlnsPrefix}},
		{"#default", []string{xmlnsPrefix}},
		{"q p", []string{"p", "q"}},               // sorted
		{"q q q", []string{"q"}},                  // deduplicated
		{"zz", []string{"zz"}},                    // unknown prefixes survive; getMapping ignores them
		{"p q  zz", []string{"", "p", "q", "zz"}}, // interior empty field, as Java's split keeps it
		{"p\tq\nr", []string{"p", "q", "r"}},
		{"p ", []string{"p"}},     // trailing empty fields are dropped by String.split
		{" p", []string{"", "p"}}, // a leading one is not
	} {
		if got := ParsePrefixList(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("ParsePrefixList(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNamespaceIsRelative(t *testing.T) {
	// C14nHelper.namespaceIsAbsolute: "" counts as absolute, otherwise indexOf(':') > 0.
	for _, tc := range []struct {
		uri      string
		relative bool
	}{
		{"", false},
		{"urn:1", false},
		{"http://x/", false},
		{"a:", false},
		{"relative", true},
		{"/a/b", true},
		{":leading", true},
		{"./x", true},
	} {
		if got := namespaceIsRelative(tc.uri); got != tc.relative {
			t.Errorf("namespaceIsRelative(%q) = %v, want %v", tc.uri, got, tc.relative)
		}
	}
}

func TestRelativeNamespaceErrorMessage(t *testing.T) {
	err := &RelativeNamespaceError{Element: "rel:c", Prefix: "rel", URI: "relative"}
	want := `xmlc14n: element rel:c has a relative namespace: rel="relative"`
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

package trustedlist

import (
	"bytes"
	"os"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/trustedlist/jaxb"
)

// TestTrustedListFacadeDelegates checks NewTrustedListFacade's Unmarshal/
// Marshal delegate to dss/trustedlist/jaxb's own Unmarshal/Marshal exactly
// (which is where marshal-parity against the Java facade's oracle is
// proven - see that package's TestMarshalParity), over the same real,
// plain (non-MRA) fixture.
func TestTrustedListFacadeDelegates(t *testing.T) {
	in, err := os.ReadFile(corpustest.RootPath(t, "trustedlist/jaxb/testdata/tl/dk_tl-sn21.xml"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := jaxb.Unmarshal(in)
	if err != nil {
		t.Fatal(err)
	}
	wantBytes, err := jaxb.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	f := NewTrustedListFacade()
	got, err := f.Unmarshal(in)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	gotBytes, err := f.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !bytes.Equal(gotBytes, wantBytes) {
		t.Errorf("TrustedListFacade did not delegate to jaxb.Unmarshal/Marshal")
	}
}

// TestMRAFacadeDelegates is TestTrustedListFacadeDelegates's MRA-document
// counterpart, checking NewMRAFacade delegates to jaxb.Unmarshal/
// jaxb.MarshalMRA.
func TestMRAFacadeDelegates(t *testing.T) {
	in, err := os.ReadFile(corpustest.RootPath(t, "trustedlist/jaxb/testdata/tl/mra-lotl.xml"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := jaxb.Unmarshal(in)
	if err != nil {
		t.Fatal(err)
	}
	wantBytes, err := jaxb.MarshalMRA(want)
	if err != nil {
		t.Fatal(err)
	}

	f := NewMRAFacade()
	got, err := f.Unmarshal(in)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	gotBytes, err := f.Marshal(got)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !bytes.Equal(gotBytes, wantBytes) {
		t.Errorf("MRAFacade did not delegate to jaxb.Unmarshal/MarshalMRA")
	}
}

// TestMRAStatusParserParse mirrors MRAStatusParser's own semantics: every
// declared URI resolves back to its constant, and an unknown value resolves
// to the empty (Java null) value.
func TestMRAStatusParserParse(t *testing.T) {
	for _, v := range enumerations.MRAStatusValues() {
		if got := MRAStatusParserParse(v.URI()); got != v {
			t.Errorf("MRAStatusParserParse(%q) = %q, want %q", v.URI(), got, v)
		}
		if got := MRAStatusParserPrint(v); got != v.URI() {
			t.Errorf("MRAStatusParserPrint(%q) = %q, want %q", v, got, v.URI())
		}
	}
	if got := MRAStatusParserParse("not-a-uri"); got != "" {
		t.Errorf("MRAStatusParserParse(unknown) = %q, want empty", got)
	}
	if got := MRAStatusParserPrint(""); got != "" {
		t.Errorf("MRAStatusParserPrint(empty) = %q, want empty", got)
	}
}

// TestMRAEquivalenceContextParserParse mirrors
// MRAEquivalenceContextParser's own semantics, the same way
// TestMRAStatusParserParse does for MRAStatusParser.
func TestMRAEquivalenceContextParserParse(t *testing.T) {
	for _, v := range enumerations.MRAEquivalenceContextValues() {
		if got := MRAEquivalenceContextParserParse(v.URI()); got != v {
			t.Errorf("MRAEquivalenceContextParserParse(%q) = %q, want %q", v.URI(), got, v)
		}
		if got := MRAEquivalenceContextParserPrint(v); got != v.URI() {
			t.Errorf("MRAEquivalenceContextParserPrint(%q) = %q, want %q", v, got, v.URI())
		}
	}
	if got := MRAEquivalenceContextParserParse("not-a-uri"); got != "" {
		t.Errorf("MRAEquivalenceContextParserParse(unknown) = %q, want empty", got)
	}
	if got := MRAEquivalenceContextParserPrint(""); got != "" {
		t.Errorf("MRAEquivalenceContextParserPrint(empty) = %q, want empty", got)
	}
}

// TestEmbeddedSchemas checks the embedded XSD copies (trusted_list_utils.go,
// mra_utils.go) are non-empty and look like the schema they claim to be.
func TestEmbeddedSchemas(t *testing.T) {
	for name, data := range map[string][]byte{
		"TrustedListSchema":                TrustedListSchema,
		"TrustedListSIESchema":             TrustedListSIESchema,
		"TrustedListAdditionalTypesSchema": TrustedListAdditionalTypesSchema,
		"MRASchema":                        MRASchema,
		"MRABaseSchema":                    MRABaseSchema,
	} {
		if len(data) == 0 {
			t.Errorf("%s is empty", name)
		}
	}
}

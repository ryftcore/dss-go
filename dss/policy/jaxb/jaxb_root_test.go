// Marshal-parity KATs for the policy validation model (DSS 6.5.RC1), per
// PORTING.md/S8A_BRIEF.md's "POLICY round-trip KATs" requirement.
//
// testdata/oracle/*.remarshal.xml are byte-exact JAXB RI output: each is
// `new JAXBContext(ObjectFactory.class)` unmarshalling the corresponding
// testdata/policy/*.xml through upstream's real ConstraintsParameters/
// ObjectFactory classes and re-marshalling with
// Marshaller.JAXB_FORMATTED_OUTPUT=true - i.e. exactly the pipeline
// ValidationPolicyFacade drives. The oracle program (PolicyOracle.java) was
// compiled and run against upstream's dss-policy-jaxb target/classes with the
// project's own Maven dependencies (offline, already-cached ~/.m2) and is not
// checked into this repository; see doc.go.
package jaxb

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// oracleCases pairs each source policy resource with its oracle remarshal.
var oracleCases = []string{
	"constraint",
	"certificate-constraint",
	"eaa-constraint",
	"qwac-constraint",
}

// TestUnmarshalMarshalMatchesOracle is the byte-exact marshal-parity KAT:
// Unmarshal(source) |> Marshal must equal the real JAXB RI's own
// remarshalling of that same source.
func TestUnmarshalMarshalMatchesOracle(t *testing.T) {
	for _, name := range oracleCases {
		t.Run(name, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join("testdata", "policy", name+".xml"))
			if err != nil {
				t.Fatalf("read source: %v", err)
			}
			want, err := os.ReadFile(filepath.Join("testdata", "oracle", name+".remarshal.xml"))
			if err != nil {
				t.Fatalf("read oracle: %v", err)
			}

			cp, err := Unmarshal(source)
			if err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			got, err := Marshal(cp)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}

			if !bytes.Equal(got, want) {
				reportXMLDiff(t, got, want)
			}
		})
	}
}

// TestUnmarshalMarshalUnmarshalRoundTrip strengthens the byte-exact KAT above
// with a semantic round-trip over the oracle output itself: Unmarshal the
// oracle's own remarshal and confirm marshalling it again reproduces the same
// bytes. This additionally exercises every field actually populated by these
// four upstream policies a second, independent way (unmarshalling
// JAXB-produced bytes rather than the hand-authored source), per the porter
// brief's "every element/attribute accounted for" requirement.
func TestUnmarshalMarshalUnmarshalRoundTrip(t *testing.T) {
	for _, name := range oracleCases {
		t.Run(name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join("testdata", "oracle", name+".remarshal.xml"))
			if err != nil {
				t.Fatalf("read oracle: %v", err)
			}
			cp, err := Unmarshal(want)
			if err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			got, err := Marshal(cp)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if !bytes.Equal(got, want) {
				reportXMLDiff(t, got, want)
			}
		})
	}
}

// reportXMLDiff fails t with the first differing line of got vs want, to
// keep failures readable against multi-hundred-line policy documents.
func reportXMLDiff(t *testing.T, got, want []byte) {
	t.Helper()
	gotLines := bytes.Split(got, []byte("\n"))
	wantLines := bytes.Split(want, []byte("\n"))
	for i := 0; i < len(gotLines) && i < len(wantLines); i++ {
		if !bytes.Equal(gotLines[i], wantLines[i]) {
			t.Fatalf("byte mismatch at line %d:\n got: %s\nwant: %s", i+1, gotLines[i], wantLines[i])
		}
	}
	t.Fatalf("byte mismatch: got %d lines, want %d lines", len(gotLines), len(wantLines))
}

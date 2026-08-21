package jaxb

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestUnmarshalMarshalIdempotent checks that Unmarshal->Marshal reaches a
// fixed point: marshalling the ALREADY-MARSHALLED bytes a second time
// produces the identical output. This is a structural self-consistency
// check independent of the Java oracle (unlike TestMarshalParity) and so
// also exercises every real fixture in the corpus, MRA or not, without
// needing a pre-computed oracle pair for it.
func TestUnmarshalMarshalIdempotent(t *testing.T) {
	dir := testdataDir(t)
	files, err := filepath.Glob(filepath.Join(dir, "*.xml"))
	if err != nil {
		t.Fatal(err)
	}
	var cases int
	for _, file := range files {
		if strings.HasSuffix(file, ".oracle.xml") {
			continue
		}
		cases++
		t.Run(filepath.Base(file), func(t *testing.T) {
			in, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			marshalFn := Marshal
			if strings.Contains(file, "mra") {
				marshalFn = MarshalMRA
			}
			tsl1, err := Unmarshal(in)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			out1, err := marshalFn(tsl1)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			tsl2, err := Unmarshal(out1)
			if err != nil {
				t.Fatalf("re-unmarshal: %v", err)
			}
			out2, err := marshalFn(tsl2)
			if err != nil {
				t.Fatalf("re-marshal: %v", err)
			}
			if !bytes.Equal(out1, out2) {
				t.Errorf("not a fixed point:\n%s", firstDiff(out1, out2))
			}
		})
	}
	if cases == 0 {
		t.Fatalf("no fixtures under %s", dir)
	}
}

// TestUnmarshalRejectsGarbage checks that Unmarshal reports an error rather
// than panicking or silently succeeding on non-XML/malformed input - the
// structural-model counterpart to TSLCORE's parseable/not-parseable
// classification (which additionally covers well-formed-but-schema-invalid
// documents, outside this package's scope - see doc.go's header).
func TestUnmarshalRejectsGarbage(t *testing.T) {
	for name, data := range map[string]string{
		"empty":      "",
		"not-xml":    "this is not XML at all",
		"truncated":  `<TrustServiceStatusList xmlns="http://uri.etsi.org/02231/v2#"><SchemeInformation>`,
		"wrong-root": `<NotATrustedList xmlns="http://example.org/other#"/>`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Unmarshal([]byte(data))
			if err == nil && name != "wrong-root" {
				t.Errorf("Unmarshal(%q) succeeded, want an error", data)
			}
		})
	}
}

// TestUnmarshalMinimalDocument checks a trusted list with no
// TrustServiceProviderList/Signature at all (a legitimate, if unusual,
// minimal SchemeInformation-only document) parses and round-trips without
// error - the Go equivalent of upstream's tl-empty-with-identifier.xml
// fixture, inlined here since that fixture (and its sibling tl-empty.xml,
// missing even the required SchemeInformation) predate SchemeInformation
// gaining enough required children to survive a real JAXB schema-validated
// remarshal, so no Java oracle exists for them.
func TestUnmarshalMinimalDocument(t *testing.T) {
	src := `<?xml version="1.0" encoding="utf-8" standalone="no"?><TrustServiceStatusList xmlns="http://uri.etsi.org/02231/v2#" Id="id_for_enveloped_signing_of_the_entire_list" TSLTag="http://uri.etsi.org/19612/TSLTag">
    <SchemeInformation>
        <TSLVersionIdentifier>6</TSLVersionIdentifier>
    </SchemeInformation>
</TrustServiceStatusList>`
	tsl, err := Unmarshal([]byte(src))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if tsl.SchemeInformation == nil || tsl.SchemeInformation.TSLVersionIdentifier == nil {
		t.Fatalf("SchemeInformation.TSLVersionIdentifier not populated: %+v", tsl)
	}
	if got := tsl.SchemeInformation.TSLVersionIdentifier.String(); got != "6" {
		t.Errorf("TSLVersionIdentifier = %q, want \"6\"", got)
	}
	if _, err := Marshal(tsl); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

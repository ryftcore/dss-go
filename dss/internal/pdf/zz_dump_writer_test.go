// Scratch harness for the manual KAT-C (Go -> Java) cross-check of DESIGN.md
// §6.3. It is env-gated and never runs in CI; Go tests never invoke Java.
//
// Reproduce (the run that produced the evidence quoted in the phase notes):
//
//	PDFOUT=/tmp/pdfwout go test ./internal/pdf/ -run TestDumpWriterOutputs
//	M2=$HOME/.m2/repository
//	CP=$M2/org/apache/pdfbox/pdfbox/3.0.7/pdfbox-3.0.7.jar\
//	:$M2/org/apache/pdfbox/pdfbox-io/3.0.7/pdfbox-io-3.0.7.jar\
//	:$M2/org/apache/pdfbox/fontbox/3.0.7/fontbox-3.0.7.jar\
//	:$M2/commons-logging/commons-logging/1.3.5/commons-logging-1.3.5.jar
//	javac -nowarn -cp "$CP" -d /tmp/wcheck WCheck.java && java -cp "$CP:/tmp/wcheck" WCheck /tmp/pdfwout
//
// pdfbox 3.0.7 loads all 14 outputs, and for each reports the xref style, the
// signature inventory and the /ByteRange this package reports.

package pdf

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDumpWriterOutputs writes the writer scenarios to $PDFOUT for a manual
// pdfbox cross-check. Scratch only.
func TestDumpWriterOutputs(t *testing.T) {
	dir := os.Getenv("PDFOUT")
	if dir == "" {
		t.Skip("no PDFOUT")
	}
	write := func(name string, b []byte) {
		if err := os.WriteFile(filepath.Join(dir, name), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	_, rawTable := openFixture(t, simpleSpec())
	write("00-orig-table.pdf", rawTable)
	specStream := simpleSpec()
	specStream.Style = XRefStream
	_, rawStream := openFixture(t, specStream)
	write("01-orig-stream.pdf", rawStream)
	specHybrid := simpleSpec()
	specHybrid.Hybrid = true
	_, rawHybrid := openFixture(t, specHybrid)
	write("02-orig-hybrid.pdf", rawHybrid)
	_, rawField := openFixture(t, emptyFieldSpec())
	write("03-orig-emptyfield.pdf", rawField)
	_, rawTwoPage := openFixture(t, twoPageSpec())
	write("04-orig-twopage.pdf", rawTwoPage)

	cms := []byte{0x30, 0x82, 0x01, 0x00, 0xAA, 0xBB, 0xCC}
	sig := SignatureOptions{
		SubFilter: "ETSI.CAdES.detached", SignerName: "Alice",
		Reason: "test", Location: "Brussels", AppName: "DSS",
		SigningTime: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), ContentSize: 512,
	}
	outT, _ := signIncrement(t, rawTable, sig, cms)
	write("10-signed-table.pdf", outT)
	outS, _ := signIncrement(t, rawStream, sig, cms)
	write("11-signed-stream.pdf", outS)
	outH, _ := signIncrement(t, rawHybrid, sig, cms)
	write("12-signed-hybrid.pdf", outH)

	fill := sig
	fill.FieldID = "Signature1"
	outF, _ := signIncrement(t, rawField, fill, cms)
	write("13-signed-fill.pdf", outF)

	vis := sig
	vis.Page = 2
	vis.Rect = Rect{MinX: 10, MinY: 20, MaxX: 210, MaxY: 70}
	vis.Appearance = NewStream(DictOf(
		Name("Type"), Name("XObject"), Name("Subtype"), Name("Form"),
		Name("BBox"), Array{Integer(0), Integer(0), Integer(200), Integer(50)},
	), []byte("q 1 0 0 RG 0 0 200 50 re S Q"))
	outV, _ := signIncrement(t, rawTwoPage, vis, cms)
	write("14-signed-visible.pdf", outV)

	ts := SignatureOptions{Type: "DocTimeStamp", SubFilter: "ETSI.RFC3161", ContentSize: 512}
	outTS, _ := signIncrement(t, outT, ts, cms)
	write("15-doctimestamp.pdf", outTS)

	out2, _ := signIncrement(t, outT, sig, cms)
	write("16-second-signature.pdf", out2)

	mdp := sig
	mdp.DocMDP = 2
	outM, _ := signIncrement(t, rawTable, mdp, cms)
	write("17-docmdp.pdf", outM)

	// /DSS + /VRI on top of a signed document.
	doc, err := OpenBytes(outT, nil)
	if err != nil {
		t.Fatal(err)
	}
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := u.SetDSSDictionary(DSSDictionary{
		Certs: []TokenRef{{Data: []byte{0x30, 0x03, 1, 2, 3}}, {Data: []byte{0x30, 0x02, 4, 5}}},
		CRLs:  []TokenRef{{Data: []byte{0x30, 0x01, 9}}},
		VRI: []VRIEntry{{
			Name:  "ABCDEF0123456789ABCDEF0123456789ABCDEF01",
			Certs: []TokenRef{{Data: []byte{0x30, 0x03, 1, 2, 3}}},
			TU:    time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC),
			TS:    []byte{0x30, 0x04, 7, 7},
		}},
	}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	write("18-dss-vri.pdf", res.Bytes)
}

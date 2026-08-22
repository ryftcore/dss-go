package enumerations

import (
	"reflect"
	"testing"
)

type mimeTypeEnumCase struct {
	v              MimeTypeEnum
	mimeTypeString string
	extensions     []string
}

func mimeTypeEnumCases() []mimeTypeEnumCase {
	return []mimeTypeEnumCase{
		{MimeTypeEnumBinary, "application/octet-stream", nil},
		{MimeTypeEnumText, "text/plain", []string{"txt"}},
		{MimeTypeEnumXML, "text/xml", []string{"xml"}},
		{MimeTypeEnumHTML, "text/html", []string{"html"}},
		{MimeTypeEnumPDF, "application/pdf", []string{"pdf"}},
		{MimeTypeEnumJSON, "application/json", []string{"json"}},
		{MimeTypeEnumJOSE, "application/jose", []string{"jose"}},
		{MimeTypeEnumJOSEJSON, "application/jose+json", []string{"json"}},
		{MimeTypeEnumSDJWTVC, "application/dc+sd-jwt", []string{"json"}},
		{MimeTypeEnumKBJWT, "application/kb+jwt", []string{"json"}},
		{MimeTypeEnumPKCS7, "application/pkcs7-signature", []string{"pkcs7", "p7m", "p7s"}},
		{MimeTypeEnumCBOR, "application/cbor", []string{"cbor"}},
		{MimeTypeEnumCose, "application/cose", []string{"cose"}},
		{MimeTypeEnumTST, "application/vnd.etsi.timestamp-token", []string{"tst"}},
		{MimeTypeEnumCRL, "application/pkix-crl", []string{"crl"}},
		{MimeTypeEnumCER, "application/pkix-cert", []string{"cer", "crt"}},
		{MimeTypeEnumZIP, "application/zip", []string{"zip"}},
		{MimeTypeEnumASiCS, "application/vnd.etsi.asic-s+zip", []string{"scs", "asics"}},
		{MimeTypeEnumASiCE, "application/vnd.etsi.asic-e+zip", []string{"sce", "asice", "bdoc"}},
		{MimeTypeEnumODT, "application/vnd.oasis.opendocument.text", []string{"odt"}},
		{MimeTypeEnumODS, "application/vnd.oasis.opendocument.spreadsheet", []string{"ods"}},
		{MimeTypeEnumODP, "application/vnd.oasis.opendocument.presentation", []string{"odp"}},
		{MimeTypeEnumODG, "application/vnd.oasis.opendocument.graphics", []string{"odg"}},
		{MimeTypeEnumPNG, "image/png", []string{"png"}},
		{MimeTypeEnumJPEG, "image/jpeg", []string{"jpg", "jpeg"}},
		{MimeTypeEnumSVG, "image/svg+xml", []string{"svg"}},
	}
}

func TestMimeTypeEnumFields(t *testing.T) {
	cases := mimeTypeEnumCases()
	if len(MimeTypeEnumValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(MimeTypeEnumValues()))
	}
	for _, c := range cases {
		if got := c.v.MimeTypeString(); got != c.mimeTypeString {
			t.Errorf("%v.MimeTypeString() = %q, want %q", c.v, got, c.mimeTypeString)
		}
		if got := c.v.Extensions(); !reflect.DeepEqual(got, c.extensions) {
			t.Errorf("%v.Extensions() = %v, want %v", c.v, got, c.extensions)
		}
		wantExt := ""
		if len(c.extensions) > 0 {
			wantExt = c.extensions[0]
		}
		if got := c.v.Extension(); got != wantExt {
			t.Errorf("%v.Extension() = %q, want %q", c.v, got, wantExt)
		}
	}
}

func TestMimeTypeEnumImplementsMimeType(t *testing.T) {
	var _ MimeType = MimeTypeEnumPDF
}

func TestMimeTypeEnumRegisteredAsLoader(t *testing.T) {
	// The package init() registers MimeTypeEnumLoader; verify resolution
	// works end-to-end through the MimeType lookup functions.
	if got := MimeTypeFromMimeTypeString("application/pdf"); got != MimeTypeEnumPDF {
		t.Errorf("MimeTypeFromMimeTypeString(pdf) = %v, want %v", got, MimeTypeEnumPDF)
	}
	if got := MimeTypeFromFileExtension("xml"); got != MimeTypeEnumXML {
		t.Errorf("MimeTypeFromFileExtension(xml) = %v, want %v", got, MimeTypeEnumXML)
	}
	if got := MimeTypeFromMimeTypeString("does/not-exist"); got != MimeTypeEnumBinary {
		t.Errorf("MimeTypeFromMimeTypeString(unknown) = %v, want %v", got, MimeTypeEnumBinary)
	}
}

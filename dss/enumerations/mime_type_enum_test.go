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
		{MimeTypeEnum_BINARY, "application/octet-stream", nil},
		{MimeTypeEnum_TEXT, "text/plain", []string{"txt"}},
		{MimeTypeEnum_XML, "text/xml", []string{"xml"}},
		{MimeTypeEnum_HTML, "text/html", []string{"html"}},
		{MimeTypeEnum_PDF, "application/pdf", []string{"pdf"}},
		{MimeTypeEnum_JSON, "application/json", []string{"json"}},
		{MimeTypeEnum_JOSE, "application/jose", []string{"jose"}},
		{MimeTypeEnum_JOSE_JSON, "application/jose+json", []string{"json"}},
		{MimeTypeEnum_SD_JWT_VC, "application/dc+sd-jwt", []string{"json"}},
		{MimeTypeEnum_KB_JWT, "application/kb+jwt", []string{"json"}},
		{MimeTypeEnum_PKCS7, "application/pkcs7-signature", []string{"pkcs7", "p7m", "p7s"}},
		{MimeTypeEnum_CBOR, "application/cbor", []string{"cbor"}},
		{MimeTypeEnum_COSE, "application/cose", []string{"cose"}},
		{MimeTypeEnum_TST, "application/vnd.etsi.timestamp-token", []string{"tst"}},
		{MimeTypeEnum_CRL, "application/pkix-crl", []string{"crl"}},
		{MimeTypeEnum_CER, "application/pkix-cert", []string{"cer", "crt"}},
		{MimeTypeEnum_ZIP, "application/zip", []string{"zip"}},
		{MimeTypeEnum_ASICS, "application/vnd.etsi.asic-s+zip", []string{"scs", "asics"}},
		{MimeTypeEnum_ASICE, "application/vnd.etsi.asic-e+zip", []string{"sce", "asice", "bdoc"}},
		{MimeTypeEnum_ODT, "application/vnd.oasis.opendocument.text", []string{"odt"}},
		{MimeTypeEnum_ODS, "application/vnd.oasis.opendocument.spreadsheet", []string{"ods"}},
		{MimeTypeEnum_ODP, "application/vnd.oasis.opendocument.presentation", []string{"odp"}},
		{MimeTypeEnum_ODG, "application/vnd.oasis.opendocument.graphics", []string{"odg"}},
		{MimeTypeEnum_PNG, "image/png", []string{"png"}},
		{MimeTypeEnum_JPEG, "image/jpeg", []string{"jpg", "jpeg"}},
		{MimeTypeEnum_SVG, "image/svg+xml", []string{"svg"}},
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
	var _ MimeType = MimeTypeEnum_PDF
}

func TestMimeTypeEnumRegisteredAsLoader(t *testing.T) {
	// The package init() registers MimeTypeEnumLoader; verify resolution
	// works end-to-end through the MimeType lookup functions.
	if got := MimeTypeFromMimeTypeString("application/pdf"); got != MimeTypeEnum_PDF {
		t.Errorf("MimeTypeFromMimeTypeString(pdf) = %v, want %v", got, MimeTypeEnum_PDF)
	}
	if got := MimeTypeFromFileExtension("xml"); got != MimeTypeEnum_XML {
		t.Errorf("MimeTypeFromFileExtension(xml) = %v, want %v", got, MimeTypeEnum_XML)
	}
	if got := MimeTypeFromMimeTypeString("does/not-exist"); got != MimeTypeEnum_BINARY {
		t.Errorf("MimeTypeFromMimeTypeString(unknown) = %v, want %v", got, MimeTypeEnum_BINARY)
	}
}

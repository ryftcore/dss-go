// Ported from dss-enumerations/.../MimeTypeEnum.java (DSS 6.5.RC1).
package enumerations

// MimeTypeEnum contains default MimeType enumerations. Implements MimeType.
type MimeTypeEnum string

const (
	// MimeTypeEnum_BINARY is octet-stream.
	MimeTypeEnum_BINARY MimeTypeEnum = "BINARY"
	// MimeTypeEnum_TEXT is plain text.
	MimeTypeEnum_TEXT MimeTypeEnum = "TEXT"
	// MimeTypeEnum_XML is xml.
	MimeTypeEnum_XML MimeTypeEnum = "XML"
	// MimeTypeEnum_HTML is html.
	MimeTypeEnum_HTML MimeTypeEnum = "HTML"
	// MimeTypeEnum_PDF is pdf.
	MimeTypeEnum_PDF MimeTypeEnum = "PDF"
	// MimeTypeEnum_JSON is json.
	MimeTypeEnum_JSON MimeTypeEnum = "JSON"
	// MimeTypeEnum_JOSE is jose.
	MimeTypeEnum_JOSE MimeTypeEnum = "JOSE"
	// MimeTypeEnum_JOSE_JSON is jose+json.
	MimeTypeEnum_JOSE_JSON MimeTypeEnum = "JOSE_JSON"
	// MimeTypeEnum_SD_JWT_VC is dc+sd-jwt.
	MimeTypeEnum_SD_JWT_VC MimeTypeEnum = "SD_JWT_VC"
	// MimeTypeEnum_KB_JWT is kb+jwt.
	MimeTypeEnum_KB_JWT MimeTypeEnum = "KB_JWT"
	// MimeTypeEnum_PKCS7 is pkcs7-signature.
	MimeTypeEnum_PKCS7 MimeTypeEnum = "PKCS7"
	// MimeTypeEnum_CBOR is cbor.
	MimeTypeEnum_CBOR MimeTypeEnum = "CBOR"
	// MimeTypeEnum_COSE is cose-sign or cose-sign1.
	MimeTypeEnum_COSE MimeTypeEnum = "COSE"
	// MimeTypeEnum_TST is timestamp-token.
	MimeTypeEnum_TST MimeTypeEnum = "TST"
	// MimeTypeEnum_CRL is crl.
	MimeTypeEnum_CRL MimeTypeEnum = "CRL"
	// MimeTypeEnum_CER is certificate.
	MimeTypeEnum_CER MimeTypeEnum = "CER"
	// MimeTypeEnum_ZIP is zip.
	MimeTypeEnum_ZIP MimeTypeEnum = "ZIP"
	// MimeTypeEnum_ASICS is asic-s.
	MimeTypeEnum_ASICS MimeTypeEnum = "ASICS"
	// MimeTypeEnum_ASICE is asic-e.
	MimeTypeEnum_ASICE MimeTypeEnum = "ASICE"
	// MimeTypeEnum_ODT is opendocument text.
	MimeTypeEnum_ODT MimeTypeEnum = "ODT"
	// MimeTypeEnum_ODS is opendocument spreadsheet.
	MimeTypeEnum_ODS MimeTypeEnum = "ODS"
	// MimeTypeEnum_ODP is opendocument presentation.
	MimeTypeEnum_ODP MimeTypeEnum = "ODP"
	// MimeTypeEnum_ODG is opendocument graphics.
	MimeTypeEnum_ODG MimeTypeEnum = "ODG"
	// MimeTypeEnum_PNG is png.
	MimeTypeEnum_PNG MimeTypeEnum = "PNG"
	// MimeTypeEnum_JPEG is jpeg.
	MimeTypeEnum_JPEG MimeTypeEnum = "JPEG"
	// MimeTypeEnum_SVG is svg.
	MimeTypeEnum_SVG MimeTypeEnum = "SVG"
)

type mimeTypeEnumFields struct {
	mimeTypeString string
	extensions     []string
}

// mimeTypeEnumData holds the (mimeTypeString, extensions) tuple for each
// constant, copied verbatim from the Java enum constructors.
var mimeTypeEnumData = map[MimeTypeEnum]mimeTypeEnumFields{
	MimeTypeEnum_BINARY:    {"application/octet-stream", nil},
	MimeTypeEnum_TEXT:      {"text/plain", []string{"txt"}},
	MimeTypeEnum_XML:       {"text/xml", []string{"xml"}},
	MimeTypeEnum_HTML:      {"text/html", []string{"html"}},
	MimeTypeEnum_PDF:       {"application/pdf", []string{"pdf"}},
	MimeTypeEnum_JSON:      {"application/json", []string{"json"}},
	MimeTypeEnum_JOSE:      {"application/jose", []string{"jose"}},
	MimeTypeEnum_JOSE_JSON: {"application/jose+json", []string{"json"}},
	MimeTypeEnum_SD_JWT_VC: {"application/dc+sd-jwt", []string{"json"}},
	MimeTypeEnum_KB_JWT:    {"application/kb+jwt", []string{"json"}},
	MimeTypeEnum_PKCS7:     {"application/pkcs7-signature", []string{"pkcs7", "p7m", "p7s"}},
	MimeTypeEnum_CBOR:      {"application/cbor", []string{"cbor"}},
	MimeTypeEnum_COSE:      {"application/cose", []string{"cose"}},
	MimeTypeEnum_TST:       {"application/vnd.etsi.timestamp-token", []string{"tst"}},
	MimeTypeEnum_CRL:       {"application/pkix-crl", []string{"crl"}},
	MimeTypeEnum_CER:       {"application/pkix-cert", []string{"cer", "crt"}},
	MimeTypeEnum_ZIP:       {"application/zip", []string{"zip"}},
	MimeTypeEnum_ASICS:     {"application/vnd.etsi.asic-s+zip", []string{"scs", "asics"}},
	MimeTypeEnum_ASICE:     {"application/vnd.etsi.asic-e+zip", []string{"sce", "asice", "bdoc"}},
	MimeTypeEnum_ODT:       {"application/vnd.oasis.opendocument.text", []string{"odt"}},
	MimeTypeEnum_ODS:       {"application/vnd.oasis.opendocument.spreadsheet", []string{"ods"}},
	MimeTypeEnum_ODP:       {"application/vnd.oasis.opendocument.presentation", []string{"odp"}},
	MimeTypeEnum_ODG:       {"application/vnd.oasis.opendocument.graphics", []string{"odg"}},
	MimeTypeEnum_PNG:       {"image/png", []string{"png"}},
	MimeTypeEnum_JPEG:      {"image/jpeg", []string{"jpg", "jpeg"}},
	MimeTypeEnum_SVG:       {"image/svg+xml", []string{"svg"}},
}

// MimeTypeEnumValues returns all constants in declaration order.
func MimeTypeEnumValues() []MimeTypeEnum {
	return []MimeTypeEnum{
		MimeTypeEnum_BINARY,
		MimeTypeEnum_TEXT,
		MimeTypeEnum_XML,
		MimeTypeEnum_HTML,
		MimeTypeEnum_PDF,
		MimeTypeEnum_JSON,
		MimeTypeEnum_JOSE,
		MimeTypeEnum_JOSE_JSON,
		MimeTypeEnum_SD_JWT_VC,
		MimeTypeEnum_KB_JWT,
		MimeTypeEnum_PKCS7,
		MimeTypeEnum_CBOR,
		MimeTypeEnum_COSE,
		MimeTypeEnum_TST,
		MimeTypeEnum_CRL,
		MimeTypeEnum_CER,
		MimeTypeEnum_ZIP,
		MimeTypeEnum_ASICS,
		MimeTypeEnum_ASICE,
		MimeTypeEnum_ODT,
		MimeTypeEnum_ODS,
		MimeTypeEnum_ODP,
		MimeTypeEnum_ODG,
		MimeTypeEnum_PNG,
		MimeTypeEnum_JPEG,
		MimeTypeEnum_SVG,
	}
}

// MimeTypeString returns the String identifying the MimeType. Implements
// MimeType.
func (m MimeTypeEnum) MimeTypeString() string {
	return mimeTypeEnumData[m].mimeTypeString
}

// Extension returns the first file extension corresponding to the
// MimeType, or "" if none is defined. Implements MimeType.
func (m MimeTypeEnum) Extension() string {
	extensions := mimeTypeEnumData[m].extensions
	if len(extensions) > 0 {
		return extensions[0]
	}
	return ""
}

// Extensions returns the full list of file extensions corresponding to the
// MimeType (Java's package-private `extensions` array field), as opposed to
// the single first extension returned by Extension().
func (m MimeTypeEnum) Extensions() []string {
	return mimeTypeEnumData[m].extensions
}

// compile-time interface assertion.
var _ MimeType = MimeTypeEnum("")

func init() {
	RegisterMimeTypeLoader(NewMimeTypeEnumLoader())
}

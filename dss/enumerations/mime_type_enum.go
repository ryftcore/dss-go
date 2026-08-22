// Ported from dss-enumerations/.../MimeTypeEnum.java (DSS 6.5.RC1).
package enumerations

// MimeTypeEnum contains default MimeType enumerations. Implements MimeType.
type MimeTypeEnum string

const (
	// MimeTypeEnumBinary is octet-stream.
	MimeTypeEnumBinary MimeTypeEnum = "BINARY"
	// MimeTypeEnumText is plain text.
	MimeTypeEnumText MimeTypeEnum = "TEXT"
	// MimeTypeEnumXML is xml.
	MimeTypeEnumXML MimeTypeEnum = "XML"
	// MimeTypeEnumHTML is html.
	MimeTypeEnumHTML MimeTypeEnum = "HTML"
	// MimeTypeEnumPDF is pdf.
	MimeTypeEnumPDF MimeTypeEnum = "PDF"
	// MimeTypeEnumJSON is json.
	MimeTypeEnumJSON MimeTypeEnum = "JSON"
	// MimeTypeEnumJOSE is jose.
	MimeTypeEnumJOSE MimeTypeEnum = "JOSE"
	// MimeTypeEnumJOSEJSON is jose+json.
	MimeTypeEnumJOSEJSON MimeTypeEnum = "JOSE_JSON"
	// MimeTypeEnumSDJWTVC is dc+sd-jwt.
	MimeTypeEnumSDJWTVC MimeTypeEnum = "SD_JWT_VC"
	// MimeTypeEnumKBJWT is kb+jwt.
	MimeTypeEnumKBJWT MimeTypeEnum = "KB_JWT"
	// MimeTypeEnumPKCS7 is pkcs7-signature.
	MimeTypeEnumPKCS7 MimeTypeEnum = "PKCS7"
	// MimeTypeEnumCBOR is cbor.
	MimeTypeEnumCBOR MimeTypeEnum = "CBOR"
	// MimeTypeEnumCose is cose-sign or cose-sign1.
	MimeTypeEnumCose MimeTypeEnum = "COSE"
	// MimeTypeEnumTST is timestamp-token.
	MimeTypeEnumTST MimeTypeEnum = "TST"
	// MimeTypeEnumCRL is crl.
	MimeTypeEnumCRL MimeTypeEnum = "CRL"
	// MimeTypeEnumCER is certificate.
	MimeTypeEnumCER MimeTypeEnum = "CER"
	// MimeTypeEnumZIP is zip.
	MimeTypeEnumZIP MimeTypeEnum = "ZIP"
	// MimeTypeEnumASiCS is asic-s.
	MimeTypeEnumASiCS MimeTypeEnum = "ASICS"
	// MimeTypeEnumASiCE is asic-e.
	MimeTypeEnumASiCE MimeTypeEnum = "ASICE"
	// MimeTypeEnumODT is opendocument text.
	MimeTypeEnumODT MimeTypeEnum = "ODT"
	// MimeTypeEnumODS is opendocument spreadsheet.
	MimeTypeEnumODS MimeTypeEnum = "ODS"
	// MimeTypeEnumODP is opendocument presentation.
	MimeTypeEnumODP MimeTypeEnum = "ODP"
	// MimeTypeEnumODG is opendocument graphics.
	MimeTypeEnumODG MimeTypeEnum = "ODG"
	// MimeTypeEnumPNG is png.
	MimeTypeEnumPNG MimeTypeEnum = "PNG"
	// MimeTypeEnumJPEG is jpeg.
	MimeTypeEnumJPEG MimeTypeEnum = "JPEG"
	// MimeTypeEnumSVG is svg.
	MimeTypeEnumSVG MimeTypeEnum = "SVG"
)

type mimeTypeEnumFields struct {
	mimeTypeString string
	extensions     []string
}

// mimeTypeEnumData holds the (mimeTypeString, extensions) tuple for each
// constant, copied verbatim from the Java enum constructors.
var mimeTypeEnumData = map[MimeTypeEnum]mimeTypeEnumFields{
	MimeTypeEnumBinary:   {"application/octet-stream", nil},
	MimeTypeEnumText:     {"text/plain", []string{"txt"}},
	MimeTypeEnumXML:      {"text/xml", []string{"xml"}},
	MimeTypeEnumHTML:     {"text/html", []string{"html"}},
	MimeTypeEnumPDF:      {"application/pdf", []string{"pdf"}},
	MimeTypeEnumJSON:     {"application/json", []string{"json"}},
	MimeTypeEnumJOSE:     {"application/jose", []string{"jose"}},
	MimeTypeEnumJOSEJSON: {"application/jose+json", []string{"json"}},
	MimeTypeEnumSDJWTVC:  {"application/dc+sd-jwt", []string{"json"}},
	MimeTypeEnumKBJWT:    {"application/kb+jwt", []string{"json"}},
	MimeTypeEnumPKCS7:    {"application/pkcs7-signature", []string{"pkcs7", "p7m", "p7s"}},
	MimeTypeEnumCBOR:     {"application/cbor", []string{"cbor"}},
	MimeTypeEnumCose:     {"application/cose", []string{"cose"}},
	MimeTypeEnumTST:      {"application/vnd.etsi.timestamp-token", []string{"tst"}},
	MimeTypeEnumCRL:      {"application/pkix-crl", []string{"crl"}},
	MimeTypeEnumCER:      {"application/pkix-cert", []string{"cer", "crt"}},
	MimeTypeEnumZIP:      {"application/zip", []string{"zip"}},
	MimeTypeEnumASiCS:    {"application/vnd.etsi.asic-s+zip", []string{"scs", "asics"}},
	MimeTypeEnumASiCE:    {"application/vnd.etsi.asic-e+zip", []string{"sce", "asice", "bdoc"}},
	MimeTypeEnumODT:      {"application/vnd.oasis.opendocument.text", []string{"odt"}},
	MimeTypeEnumODS:      {"application/vnd.oasis.opendocument.spreadsheet", []string{"ods"}},
	MimeTypeEnumODP:      {"application/vnd.oasis.opendocument.presentation", []string{"odp"}},
	MimeTypeEnumODG:      {"application/vnd.oasis.opendocument.graphics", []string{"odg"}},
	MimeTypeEnumPNG:      {"image/png", []string{"png"}},
	MimeTypeEnumJPEG:     {"image/jpeg", []string{"jpg", "jpeg"}},
	MimeTypeEnumSVG:      {"image/svg+xml", []string{"svg"}},
}

// MimeTypeEnumValues returns all constants in declaration order.
func MimeTypeEnumValues() []MimeTypeEnum {
	return []MimeTypeEnum{
		MimeTypeEnumBinary,
		MimeTypeEnumText,
		MimeTypeEnumXML,
		MimeTypeEnumHTML,
		MimeTypeEnumPDF,
		MimeTypeEnumJSON,
		MimeTypeEnumJOSE,
		MimeTypeEnumJOSEJSON,
		MimeTypeEnumSDJWTVC,
		MimeTypeEnumKBJWT,
		MimeTypeEnumPKCS7,
		MimeTypeEnumCBOR,
		MimeTypeEnumCose,
		MimeTypeEnumTST,
		MimeTypeEnumCRL,
		MimeTypeEnumCER,
		MimeTypeEnumZIP,
		MimeTypeEnumASiCS,
		MimeTypeEnumASiCE,
		MimeTypeEnumODT,
		MimeTypeEnumODS,
		MimeTypeEnumODP,
		MimeTypeEnumODG,
		MimeTypeEnumPNG,
		MimeTypeEnumJPEG,
		MimeTypeEnumSVG,
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

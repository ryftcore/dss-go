// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PDFServiceMode.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). This enum's four constants are exactly what
// native_pdf_obj_factory.go / native_pdf_signature_service.go were already written against
// (PDFServiceMode_SIGNATURE, _CONTENT_TIMESTAMP, _SIGNATURE_TIMESTAMP, _ARCHIVE_TIMESTAMP -
// Java_ENUM_CASE kept verbatim after the type name, as is this port's convention for Go-ported
// Java enum constants). A string base (matching the enumerations package's convention, rather
// than an int/iota), because native_pdf_signature_service.go's already-landed constructor
// zero-checks it with `serviceMode == ""`.
package pades

// PDFServiceMode defines the executing PDF service modes. Port of the PDFServiceMode enum.
type PDFServiceMode string

const (
	// PDFServiceMode_CONTENT_TIMESTAMP is used for independent timestamp creation.
	PDFServiceMode_CONTENT_TIMESTAMP PDFServiceMode = "CONTENT_TIMESTAMP"

	// PDFServiceMode_SIGNATURE is used for signature creation.
	PDFServiceMode_SIGNATURE PDFServiceMode = "SIGNATURE"

	// PDFServiceMode_SIGNATURE_TIMESTAMP is used for signature timestamp creation.
	PDFServiceMode_SIGNATURE_TIMESTAMP PDFServiceMode = "SIGNATURE_TIMESTAMP"

	// PDFServiceMode_ARCHIVE_TIMESTAMP is used for document timestamp creation.
	PDFServiceMode_ARCHIVE_TIMESTAMP PDFServiceMode = "ARCHIVE_TIMESTAMP"
)

// String ports the implicit Enum#name/toString.
func (m PDFServiceMode) String() string {
	switch m {
	case PDFServiceMode_CONTENT_TIMESTAMP:
		return "CONTENT_TIMESTAMP"
	case PDFServiceMode_SIGNATURE:
		return "SIGNATURE"
	case PDFServiceMode_SIGNATURE_TIMESTAMP:
		return "SIGNATURE_TIMESTAMP"
	case PDFServiceMode_ARCHIVE_TIMESTAMP:
		return "ARCHIVE_TIMESTAMP"
	default:
		return "UNKNOWN"
	}
}

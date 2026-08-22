// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PDFServiceMode.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). This enum's four constants are exactly what
// native_pdf_obj_factory.go / native_pdf_signature_service.go were already written against
// (PDFServiceModeSignature, PDFServiceModeContentTimestamp, PDFServiceModeSignatureTimestamp,
// PDFServiceModeArchiveTimestamp - the type name prefixes the Java constant name rendered in Go
// MixedCaps, as is this port's convention for Go-ported Java enum constants; the string VALUES
// stay Java's name() verbatim). A string base (matching the enumerations package's convention, rather
// than an int/iota), because native_pdf_signature_service.go's already-landed constructor
// zero-checks it with `serviceMode == ""`.
package pades

// PDFServiceMode defines the executing PDF service modes. Port of the PDFServiceMode enum.
type PDFServiceMode string

const (
	// PDFServiceModeContentTimestamp is used for independent timestamp creation.
	PDFServiceModeContentTimestamp PDFServiceMode = "CONTENT_TIMESTAMP"

	// PDFServiceModeSignature is used for signature creation.
	PDFServiceModeSignature PDFServiceMode = "SIGNATURE"

	// PDFServiceModeSignatureTimestamp is used for signature timestamp creation.
	PDFServiceModeSignatureTimestamp PDFServiceMode = "SIGNATURE_TIMESTAMP"

	// PDFServiceModeArchiveTimestamp is used for document timestamp creation.
	PDFServiceModeArchiveTimestamp PDFServiceMode = "ARCHIVE_TIMESTAMP"
)

// String ports the implicit Enum#name/toString.
func (m PDFServiceMode) String() string {
	switch m {
	case PDFServiceModeContentTimestamp:
		return "CONTENT_TIMESTAMP"
	case PDFServiceModeSignature:
		return "SIGNATURE"
	case PDFServiceModeSignatureTimestamp:
		return "SIGNATURE_TIMESTAMP"
	case PDFServiceModeArchiveTimestamp:
		return "ARCHIVE_TIMESTAMP"
	default:
		return "UNKNOWN"
	}
}

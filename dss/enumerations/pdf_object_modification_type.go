// Ported from dss-enumerations/.../PdfObjectModificationType.java (DSS 6.5.RC1).
package enumerations

// PdfObjectModificationType specifies a modification origin kind.
type PdfObjectModificationType string

const (
	// PdfObjectModificationType_CREATION represents an object addition to a
	// final revision.
	PdfObjectModificationType_CREATION PdfObjectModificationType = "CREATION"
	// PdfObjectModificationType_DELETION represents an object deletion from
	// a final revision.
	PdfObjectModificationType_DELETION PdfObjectModificationType = "DELETION"
	// PdfObjectModificationType_MODIFICATION represents an object change in
	// a final revision.
	PdfObjectModificationType_MODIFICATION PdfObjectModificationType = "MODIFICATION"
)

// PdfObjectModificationTypeValues returns all constants in declaration order.
func PdfObjectModificationTypeValues() []PdfObjectModificationType {
	return []PdfObjectModificationType{
		PdfObjectModificationType_CREATION,
		PdfObjectModificationType_DELETION,
		PdfObjectModificationType_MODIFICATION,
	}
}

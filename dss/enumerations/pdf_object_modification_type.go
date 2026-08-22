// Ported from dss-enumerations/.../PdfObjectModificationType.java (DSS 6.5.RC1).
package enumerations

// PdfObjectModificationType specifies a modification origin kind.
type PdfObjectModificationType string

const (
	// PdfObjectModificationTypeCreation represents an object addition to a
	// final revision.
	PdfObjectModificationTypeCreation PdfObjectModificationType = "CREATION"
	// PdfObjectModificationTypeDeletion represents an object deletion from
	// a final revision.
	PdfObjectModificationTypeDeletion PdfObjectModificationType = "DELETION"
	// PdfObjectModificationTypeModification represents an object change in
	// a final revision.
	PdfObjectModificationTypeModification PdfObjectModificationType = "MODIFICATION"
)

// PdfObjectModificationTypeValues returns all constants in declaration order.
func PdfObjectModificationTypeValues() []PdfObjectModificationType {
	return []PdfObjectModificationType{
		PdfObjectModificationTypeCreation,
		PdfObjectModificationTypeDeletion,
		PdfObjectModificationTypeModification,
	}
}

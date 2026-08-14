// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfModification.java and
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/CommonPdfModification.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header for the sibling eu.europa.esig.dss.pdf package;
// this sub-package shares the same gap). java.io.Serializable is dropped, as elsewhere.
package pades

// PdfModification contains information about a modification that occurred in a PDF.
// Port of the PdfModification interface.
type PdfModification interface {
	// Page returns the page where the modification occurs. Port of getPage().
	Page() int
}

// CommonPdfModification is the default PDF Modification object. Port of the
// CommonPdfModification class.
type CommonPdfModification struct {
	// page defines the page of the found modification.
	page int
}

// NewCommonPdfModification builds a CommonPdfModification for the given page.
// Port of the CommonPdfModification(int) constructor.
func NewCommonPdfModification(page int) *CommonPdfModification {
	return &CommonPdfModification{page: page}
}

// Page returns the modified page. Port of #getPage.
func (m *CommonPdfModification) Page() int { return m.page }

// Compile-time assertion standing in for Java's "implements PdfModification".
var _ PdfModification = (*CommonPdfModification)(nil)

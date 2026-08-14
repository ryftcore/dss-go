// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfObjectModificationsFinder.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header).
package pades

// PdfObjectModificationsFinder is used to find and return all object modifications occurred
// between two PDF document revisions. Port of the PdfObjectModificationsFinder interface.
type PdfObjectModificationsFinder interface {
	// Find returns found and categorized object modifications occurred between
	// originalRevisionReader and finalRevisionReader. Port of find(PdfDocumentReader, PdfDocumentReader).
	Find(originalRevisionReader, finalRevisionReader PdfDocumentReader) PdfObjectModifications
}

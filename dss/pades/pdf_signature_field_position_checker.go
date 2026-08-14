// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureFieldPositionChecker.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header).
package pades

import (
	"fmt"

	"github.com/utain/esig/dss/alert"
)

// PdfSignatureFieldPositionChecker is used to verify the correctness of a new signature field
// placement within a PDF document. Port of the PdfSignatureFieldPositionChecker class.
type PdfSignatureFieldPositionChecker struct {
	// alertOnSignatureFieldOverlap sets the behavior to follow in case of overlapping a new
	// signature field with existing annotations. Default: ExceptionOnStatusAlert.
	alertOnSignatureFieldOverlap alert.StatusAlert

	// alertOnSignatureFieldOutsidePageDimensions sets the behavior when a new signature field is
	// created outside the page dimensions. Default: ExceptionOnStatusAlert.
	alertOnSignatureFieldOutsidePageDimensions alert.StatusAlert

	// alertOnDocumentReadException sets the behavior when an error is thrown on an attempt to
	// read the document properties. Default: ExceptionOnStatusAlert.
	alertOnDocumentReadException alert.StatusAlert
}

// NewPdfSignatureFieldPositionChecker instantiates the checker with the default configuration.
// Port of the default constructor.
func NewPdfSignatureFieldPositionChecker() *PdfSignatureFieldPositionChecker {
	return &PdfSignatureFieldPositionChecker{
		alertOnSignatureFieldOverlap:               alert.NewExceptionOnStatusAlert(),
		alertOnSignatureFieldOutsidePageDimensions: alert.NewExceptionOnStatusAlert(),
		alertOnDocumentReadException:               alert.NewExceptionOnStatusAlert(),
	}
}

// SetAlertOnSignatureFieldOverlap sets the alert for a signature field overlap with existing
// fields or/and annotations. Default: ExceptionOnStatusAlert. Port of #setAlertOnSignatureFieldOverlap.
func (c *PdfSignatureFieldPositionChecker) SetAlertOnSignatureFieldOverlap(alertOnSignatureFieldOverlap alert.StatusAlert) {
	if alertOnSignatureFieldOverlap == nil {
		panic("StatusAlert cannot be null!")
	}
	c.alertOnSignatureFieldOverlap = alertOnSignatureFieldOverlap
}

// SetAlertOnSignatureFieldOutsidePageDimensions sets the behavior to follow when a new signature
// field is created outside the page's dimensions. Default: ExceptionOnStatusAlert.
// Port of #setAlertOnSignatureFieldOutsidePageDimensions.
func (c *PdfSignatureFieldPositionChecker) SetAlertOnSignatureFieldOutsidePageDimensions(alertOnSignatureFieldOutsidePageDimensions alert.StatusAlert) {
	if alertOnSignatureFieldOutsidePageDimensions == nil {
		panic("StatusAlert cannot be null!")
	}
	c.alertOnSignatureFieldOutsidePageDimensions = alertOnSignatureFieldOutsidePageDimensions
}

// SetAlertOnDocumentReadException sets the behavior to follow when an error is thrown on an
// attempt to read document properties. Default: ExceptionOnStatusAlert.
// Port of #setAlertOnDocumentReadException.
func (c *PdfSignatureFieldPositionChecker) SetAlertOnDocumentReadException(alertOnDocumentReadException alert.StatusAlert) {
	c.alertOnDocumentReadException = alertOnDocumentReadException
}

// AssertSignatureFieldPositionValid verifies whether annotationBox can be placed within
// documentReader on the given pageNumber. Port of #assertSignatureFieldPositionValid.
func (c *PdfSignatureFieldPositionChecker) AssertSignatureFieldPositionValid(documentReader PdfDocumentReader,
	annotationBox AnnotationBox, pageNumber int) {
	pageBox := documentReader.PageBox(pageNumber)
	c.CheckSignatureFieldAgainstPageDimensions(annotationBox, pageBox)
	pdfAnnotations := c.annotations(documentReader, pageNumber)
	c.CheckSignatureFieldBoxOverlap(annotationBox, pdfAnnotations)
}

// CheckSignatureFieldBoxOverlap verifies whether signatureFieldBox overlaps with one of the
// extracted pdfAnnotations. Port of #checkSignatureFieldBoxOverlap.
func (c *PdfSignatureFieldPositionChecker) CheckSignatureFieldBoxOverlap(signatureFieldBox AnnotationBox, pdfAnnotations []*PdfAnnotation) {
	pdfDifferencesFinder := NewDefaultPdfDifferencesFinder()
	if pdfDifferencesFinder.IsAnnotationBoxOverlapping(signatureFieldBox, pdfAnnotations) {
		c.alertOnSignatureFieldOverlapMessage()
	}
}

// alertOnSignatureFieldOverlapMessage executes the alertOnSignatureFieldOverlap alert.
// Port of the private #alertOnSignatureFieldOverlap.
func (c *PdfSignatureFieldPositionChecker) alertOnSignatureFieldOverlapMessage() {
	status := alert.NewMessageStatus()
	status.SetMessage("The new signature field position overlaps with an existing annotation!")
	if err := c.alertOnSignatureFieldOverlap.Alert(status); err != nil {
		panic(err)
	}
}

// CheckSignatureFieldAgainstPageDimensions verifies whether signatureFieldBox is within pageBox.
// Port of #checkSignatureFieldAgainstPageDimensions.
func (c *PdfSignatureFieldPositionChecker) CheckSignatureFieldAgainstPageDimensions(signatureFieldBox, pageBox AnnotationBox) {
	if signatureFieldBox.MinX() < pageBox.MinX() || signatureFieldBox.MaxX() > pageBox.MaxX() ||
		signatureFieldBox.MinY() < pageBox.MinY() || signatureFieldBox.MaxY() > pageBox.MaxY() {
		c.alertOnSignatureFieldOutsidePageDimensionsMessage(signatureFieldBox, pageBox)
	}
}

// alertOnSignatureFieldOutsidePageDimensionsMessage executes the
// alertOnSignatureFieldOutsidePageDimensions alert. Port of the private
// #alertOnSignatureFieldOutsidePageDimensions.
func (c *PdfSignatureFieldPositionChecker) alertOnSignatureFieldOutsidePageDimensionsMessage(signatureFieldBox, pageBox AnnotationBox) {
	status := alert.NewMessageStatus()
	status.SetMessage(fmt.Sprintf("The new signature field position is outside the page dimensions! "+
		"Signature Field : [minX=%v, maxX=%v, minY=%v, maxY=%v], "+
		"Page : [minX=%v, maxX=%v, minY=%v, maxY=%v]",
		signatureFieldBox.MinX(), signatureFieldBox.MaxX(), signatureFieldBox.MinY(), signatureFieldBox.MaxY(),
		pageBox.MinX(), pageBox.MaxX(), pageBox.MinY(), pageBox.MaxY()))
	if err := c.alertOnSignatureFieldOutsidePageDimensions.Alert(status); err != nil {
		panic(err)
	}
}

// annotations ports the private getAnnotations(PdfDocumentReader, int).
func (c *PdfSignatureFieldPositionChecker) annotations(documentReader PdfDocumentReader, pageNumber int) []*PdfAnnotation {
	annotations, err := documentReader.PdfAnnotations(pageNumber)
	if err != nil {
		c.alertOnDocumentReadExceptionMessage(err)
		return nil
	}
	return annotations
}

// alertOnDocumentReadExceptionMessage executes the alertOnDocumentReadException alert.
// Port of the private #alertOnDocumentReadException.
func (c *PdfSignatureFieldPositionChecker) alertOnDocumentReadExceptionMessage(err error) {
	status := alert.NewMessageStatus()
	status.SetMessage(fmt.Sprintf("An error occurred while reading PDF document! Reason : %s", err.Error()))
	if alertErr := c.alertOnDocumentReadException.Alert(status); alertErr != nil {
		panic(alertErr)
	}
}

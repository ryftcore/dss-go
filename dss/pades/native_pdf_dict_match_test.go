// Regression cover for the PdfDict#match null rule, which nothing else exercises.
//
// Upstream PdfBoxDict#match reads each expected entry through
// COSDictionary#getDictionaryObject, which folds an absent entry AND a COSNull entry into
// null, and then skips the comparison for a null:
//
//	for (COSName key : pdfBoxDict.wrapped.keySet()) {
//	    COSBase targetObject = pdfBoxDict.wrapped.getDictionaryObject(key);
//	    COSBase currentObject = wrapped.getDictionaryObject(key);
//	    if (targetObject != null && !targetObject.equals(currentObject)) {
//	        return false;
//	    }
//	}
//	return true;
//
// internal/pdf's Resolve returns pdf.Null{} rather than a nil interface for exactly those
// cases, so the port has to test for Null to reproduce the skip.
package pades

import (
	"testing"

	"github.com/utain/esig/dss/internal/pdf"
)

func TestNativePdfDictMatchSkipsNullExpectedEntries(t *testing.T) {
	document := &pdf.Document{}

	// A null expected entry is Java's null target object: the comparison is skipped and the
	// candidate's own value for that key is irrelevant.
	expected := pdf.NewDict()
	expected.Set("BaseVersion", pdf.Name("1.7"))
	expected.Set("ExtensionLevel", pdf.Null{})

	candidate := pdf.NewDict()
	candidate.Set("BaseVersion", pdf.Name("1.7"))
	candidate.Set("ExtensionLevel", pdf.Integer(8))

	if !nativePDFSignatureServiceMatch(document, candidate, expected) {
		t.Error("a null expected entry must be skipped, as getDictionaryObject maps COSNull to null")
	}

	// An absent expected entry behaves the same way, and a present-but-different one still
	// fails the match.
	absent := pdf.NewDict()
	absent.Set("BaseVersion", pdf.Name("1.7"))
	if !nativePDFSignatureServiceMatch(document, candidate, absent) {
		t.Error("an expected dictionary that is a subset of the candidate must match")
	}

	different := pdf.NewDict()
	different.Set("BaseVersion", pdf.Name("2.0"))
	if nativePDFSignatureServiceMatch(document, candidate, different) {
		t.Error("a differing non-null expected entry must fail the match")
	}
}

func TestNativePdfDictMatchSkipsNullExpectedEntriesThroughPdfDict(t *testing.T) {
	document := &pdf.Document{}

	expectedDict := pdf.NewDict()
	expectedDict.Set("Type", pdf.Name("Sig"))
	expectedDict.Set("Contents", pdf.Null{})

	candidateDict := pdf.NewDict()
	candidateDict.Set("Type", pdf.Name("Sig"))
	candidateDict.Set("Contents", pdf.String{Bytes: []byte{0x01, 0x02}})

	candidate := newNativePdfDict(document, candidateDict, pdf.ObjectKey{})
	expected := newNativePdfDict(document, expectedDict, pdf.ObjectKey{})

	if !candidate.Match(expected) {
		t.Error("PdfDict#match must skip an expected entry whose value is the PDF null object")
	}
}

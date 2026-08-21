// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfDocTimestampRevision.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). slf4j is dropped, per PORTING.md.
//
// Java's `timestampToken = new PdfTimestampToken(this)` reads the not-yet-fully-constructed
// PdfDocTimestampRevision through `this` (legal in Java: the super() call has already run, so
// signatureDictionary etc. are set). This port reproduces that construction order explicitly:
// the *PdfDocTimestampRevision is allocated with its base fields populated and timestampToken
// still nil, handed to NewPdfTimestampToken (which only reads PdfSigDictInfo() - available
// already), and the resulting token is then stored back and matched against the signed data -
// mirroring pdf_timestamp_token_identifier_builder.go's precedent for the same "this passed
// mid-construction" shape.
package pades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// PdfDocTimestampRevision is the signature timestamp representation. This class is only used in
// case of a Document Timestamp (not a signature-timestamp from CAdES/CMS).
// Port of the PdfDocTimestampRevision class, extending PdfCMSRevision.
type PdfDocTimestampRevision struct {
	pdfCMSRevisionBase

	// timestampToken is the document timestamp token from the revision.
	timestampToken *PdfTimestampToken
}

// NewPdfDocTimestampRevision is the default constructor to create a PdfDocTimestampRevision.
// Port of the PdfDocTimestampRevision(PdfSignatureDictionary, List<PdfSignatureField>,
// DSSDocument, DSSDocument, boolean) constructor.
//
// Panics (Java's DSSException, thrown unchecked from this constructor) when the timestamp token
// cannot be built or matched against the signed data.
func NewPdfDocTimestampRevision(signatureDictionary *PdfSignatureDictionary, timestampFields []*PdfSignatureField,
	signedContent, previousRevision model.DSSDocument, coverCompleteRevision bool) *PdfDocTimestampRevision {
	revision := &PdfDocTimestampRevision{
		pdfCMSRevisionBase: newPdfCMSRevisionBase(signatureDictionary, timestampFields, signedContent, previousRevision, coverCompleteRevision),
	}

	timestampToken, err := NewPdfTimestampToken(revision)
	if err != nil {
		panic(fmt.Sprintf("Unable to create a PdfDocTimestampRevision : %s", err.Error()))
	}
	if _, err := timestampToken.MatchDataDocument(revision.SignedData()); err != nil {
		panic(fmt.Sprintf("Unable to create a PdfDocTimestampRevision : %s", err.Error()))
	}
	revision.timestampToken = timestampToken
	// Upstream logs "Created PdfDocTimestampInfo : {}" at debug level.

	return revision
}

// SigningDate gets the claimed signing time. Port of the overridden #getSigningDate.
func (r *PdfDocTimestampRevision) SigningDate() time.Time {
	return r.timestampToken.GenerationTime()
}

// TimestampToken returns the corresponding PdfTimestampToken. Port of #getTimestampToken.
func (r *PdfDocTimestampRevision) TimestampToken() *PdfTimestampToken { return r.timestampToken }

// Compile-time assertion standing in for Java's "extends PdfCMSRevision".
var _ PdfCMSRevision = (*PdfDocTimestampRevision)(nil)

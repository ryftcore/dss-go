// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/timestamp/PdfRevisionTimestampSource.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES:
//   - PdfDocDssRevision (eu.europa.esig.dss.pdf.PdfDocDssRevision) - pdf_document_analyzer.go's
//     header already assumes CertificateSource() *PdfDssDictCertificateSource,
//     CRLSource() *PdfDssDictCRLSource, OCSPSource() *PdfDssDictOCSPSource,
//     DssDictionary() PdfDssDict; it implements PdfRevision directly (not PdfCMSRevision, per
//     upstream's `class PdfDocDssRevision implements PdfRevision`).
//   - PdfDocTimestampRevision (eu.europa.esig.dss.pdf.PdfDocTimestampRevision) -
//     pdf_document_analyzer.go's header already assumes TimestampToken() *PdfTimestampToken.
//     PdfTimestampToken is landed by this same chunk (pdf_timestamp_token.go).
package pades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/timestamp"
)

// PdfRevisionTimestampSource extracts a timestamp from a single PdfRevision. Port of the class
// PdfRevisionTimestampSource, extending timestamp.AbstractTimestampSource.
type PdfRevisionTimestampSource struct {
	timestamp.AbstractTimestampSource

	// pdfRevision is the PdfRevision to extract references from.
	pdfRevision PdfRevision

	// certificateSource is the merged CertificateSource to find certificate binaries from.
	certificateSource *spi.ListCertificateSource

	// crlSource is the merged CRL source.
	crlSource *spi.ListRevocationSource[revocation.CRL]

	// ocspSource is the merged OCSP source.
	ocspSource *spi.ListRevocationSource[revocation.OCSP]
}

// NewPdfRevisionTimestampSource is the default constructor.
func NewPdfRevisionTimestampSource(pdfRevision PdfRevision, certificateSource *spi.ListCertificateSource,
	crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) *PdfRevisionTimestampSource {
	return &PdfRevisionTimestampSource{
		pdfRevision:       pdfRevision,
		certificateSource: certificateSource,
		crlSource:         crlSource,
		ocspSource:        ocspSource,
	}
}

// IncorporatedReferences returns incorporated references for the revision. Port of
// getIncorporatedReferences().
//
// Java's getReferencesFromTimestamp(TimestampToken, ...) never declares a checked exception
// (any failure surfaces as an unchecked DSSException); this method's own signature has no error
// return either (per PORTING.md, matching Java's `List<TimestampedReference>
// getIncorporatedReferences()`), so a genuine failure here panics with a model.DSSError, exactly
// as the unchecked Java exception would propagate.
func (s *PdfRevisionTimestampSource) IncorporatedReferences() []*validation.TimestampedReference {
	switch pdfRevision := s.pdfRevision.(type) {
	case *PdfDocTimestampRevision:
		timestampToken := pdfRevision.TimestampToken()
		references, err := timestamp.ReferencesFromTimestamp(timestampToken.TimestampToken, s.certificateSource, s.crlSource, s.ocspSource)
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		return references

	case *PdfDocDssRevision:
		references := []*validation.TimestampedReference{}

		dssCertificateSource := pdfRevision.CertificateSource()
		pdfRevisionTSAddReferences(&references, timestamp.CreateReferencesForCertificates(dssCertificateSource.Certificates()))

		dssCRLSource := pdfRevision.CRLSource()
		pdfRevisionTSAddReferences(&references, timestamp.CreateReferencesForCRLBinaries(dssCRLSource.DSSDictionaryBinaries()))
		pdfRevisionTSAddReferences(&references, timestamp.CreateReferencesForCRLBinaries(dssCRLSource.VRIDictionaryBinaries()))

		dssOCSPSource := pdfRevision.OCSPSource()
		dssBinaryReferences, err := timestamp.CreateReferencesForOCSPBinaries(dssOCSPSource.DSSDictionaryBinaries(), s.certificateSource)
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		pdfRevisionTSAddReferences(&references, dssBinaryReferences)
		vriBinaryReferences, err := timestamp.CreateReferencesForOCSPBinaries(dssOCSPSource.VRIDictionaryBinaries(), s.certificateSource)
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		pdfRevisionTSAddReferences(&references, vriBinaryReferences)

		return references

	default:
		return []*validation.TimestampedReference{}
	}
}

// VRITimestampToken returns a timestamp token extracted from the VRI dictionary with the given
// key. Port of getVRITimestampToken(String); vriKey is the sha-1 of the corresponding signature
// value.
func (s *PdfRevisionTimestampSource) VRITimestampToken(vriKey string) *validation.TimestampToken {
	if pdfRevision, ok := s.pdfRevision.(*PdfDocDssRevision); ok {
		pdfVriDictTimestampSource := NewPdfVriDictSource(pdfRevision.DssDictionary(), vriKey)
		return pdfVriDictTimestampSource.TimestampToken()
	}
	return nil
}

// -----------------------------------------------------------------------------
// Small reference-list helper, mirroring the frozen spi/validation/timestamp package's own
// unexported addReferences (abstract_timestamp_source.go), which this file cannot call directly
// (different package). Prefixed distinctly (pdfRevisionTS) to avoid colliding with any sibling
// file in this package needing the same helper - see cades_timestamp_source.go's cadesTS-
// prefixed precedent for the same situation.
// -----------------------------------------------------------------------------

// pdfRevisionTSAddReferences adds referencesToAdd to *referenceList without duplicates (by
// TimestampedReference.Equals). Port of the AbstractTimestampSource addReferences(List, List)
// helper.
func pdfRevisionTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
	for _, candidate := range referencesToAdd {
		found := false
		for _, existing := range *referenceList {
			if existing.Equals(candidate) {
				found = true
				break
			}
		}
		if !found {
			*referenceList = append(*referenceList, candidate)
		}
	}
}

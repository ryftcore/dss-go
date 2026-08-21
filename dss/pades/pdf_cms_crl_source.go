// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfCmsCRLSource.java
// (DSS 6.5.RC1).
//
// BouncyCastle replacements used here (see PORTING.md), following the precedent
// pdf_cms_ocsp_source.go established for the sibling PdfCmsOCSPSource:
//   - org.bouncycastle.asn1.cms.AttributeTable -> internal/cmscore.Attributes.
//   - Attribute#getAttributeValues() -> Attribute.Values ([]*asn1ber.Element).
//   - DSSASN1Utils.getAsn1Attributes(AttributeTable, ASN1ObjectIdentifier) ->
//     spi.DSSASN1UtilsAsn1Attributes.
//   - CRLUtils.buildCRLBinary(byte[]) -> crlparser.CRLUtilsBuildCRLBinary([]byte); Java's thrown
//     Exception on a malformed CRL becomes the returned error, matched with a skip-and-continue
//     the same way every other malformed-revocation degradation in this port does (slf4j
//     dropped per PORTING.md).
//
// FORWARD DEPENDENCY (not in this chunk's manifest, shared with pdf_cms_ocsp_source.go):
//   - PAdESUtilsRevocationInfoArchival(attrValue *asn1ber.Element) *RevocationInfoArchival - the
//     flattened static PAdESUtils.getRevocationInfoArchival(ASN1Encodable).
//
// AllRevocationBinariesWithOrigins()/AllRevocationTokensWithOrigins() are declared by the Java
// OfflineRevocationSource base (upstream, PdfCmsCRLSource never overrides them, inheriting the
// base's default map-field getter) but were not carried over into spi.OfflineRevocationSourceBase
// (see pdf_dss_dict_crl_source.go's header for the same gap, and pades_crl_source.go's header
// for the landed sibling chunk that already assumes this file defines them). Every binary this
// source ever adds carries the single RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL origin
// (extractRevocationInfoArchival below is the only AddBinary call site), so the two methods
// below reproduce the base's default behaviour by pairing spi.OfflineRevocationSourceBase's
// already-ported origin-filtered accessors (ADBERevocationValuesBinaries/ADBERevocationValuesTokens,
// promoted from spi.OfflineCRLSourceBase) with that one constant origin, rather than needing
// access to the base's private per-binary/per-token origin bookkeeping.
package pades

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PdfCmsCRLSourceBinaryOriginsEntry pairs a CRL binary with the origins it has been found with,
// mirroring PdfDssDictCRLSourceBinaryOriginsEntry's {Binary, Origins} shape
// (pdf_dss_dict_crl_source.go), the naming pades_crl_source.go's header already assumes.
type PdfCmsCRLSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]
	Origins []enumerations.RevocationOrigin
}

// PdfCmsCRLSource represents a source of CRL tokens extracted from a PDF's CMS.
type PdfCmsCRLSource struct {
	spi.OfflineCRLSourceBase
}

// NewPdfCmsCRLSource is the default constructor. Port of the constructor
// PdfCmsCRLSource(AttributeTable).
func NewPdfCmsCRLSource(signedAttributes cmscore.Attributes) *PdfCmsCRLSource {
	source := &PdfCmsCRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's inherited RevocationTokens implementation.
	source.InitOfflineRevocationSource(source)
	source.extractCRLArchivalValues(signedAttributes)
	return source
}

// extractCRLArchivalValues ports the private extractCRLArchivalValues(AttributeTable).
func (s *PdfCmsCRLSource) extractCRLArchivalValues(signedAttributes cmscore.Attributes) {
	if signedAttributes != nil {
		attributes := spi.DSSASN1UtilsAsn1Attributes(signedAttributes, spi.OID_adbe_revocationInfoArchival)
		for _, attribute := range attributes {
			for _, attrValue := range attribute.Values {
				s.extractRevocationInfoArchival(attrValue)
			}
		}
	}
}

// extractRevocationInfoArchival ports the private extractRevocationInfoArchival(ASN1Encodable).
//
// DEVIATION: upstream logs a warning ("Could not convert CertificateList to CRLBinary : {}")
// and continues when CRLUtils.buildCRLBinary raises an exception; slf4j logging is dropped per
// PORTING.md, so the entry is simply skipped here, same as every other malformed-revocation
// degradation elsewhere in this port.
func (s *PdfCmsCRLSource) extractRevocationInfoArchival(attrValue *asn1ber.Element) {
	revValues := PAdESUtilsRevocationInfoArchival(attrValue)
	if revValues != nil {
		for _, revValue := range revValues.CrlVals() {
			crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(revValue)
			if err != nil {
				continue
			}
			s.AddBinary(crlBinary, enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL)
		}
	}
}

// AllRevocationBinariesWithOrigins returns all the CRL binaries of this source with their
// origins. Port of the inherited getAllRevocationBinariesWithOrigins(); see the file header for
// why this is a local reproduction rather than an override.
func (s *PdfCmsCRLSource) AllRevocationBinariesWithOrigins() []PdfCmsCRLSourceBinaryOriginsEntry {
	binaries := s.ADBERevocationValuesBinaries()
	result := make([]PdfCmsCRLSourceBinaryOriginsEntry, 0, len(binaries))
	for _, binary := range binaries {
		result = append(result, PdfCmsCRLSourceBinaryOriginsEntry{
			Binary:  binary,
			Origins: []enumerations.RevocationOrigin{enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL},
		})
	}
	return result
}

// AllRevocationTokensWithOrigins returns all the CRL tokens of this source with their origins.
// Port of the inherited getAllRevocationTokensWithOrigins(); see the file header for why this is
// a local reproduction rather than an override.
func (s *PdfCmsCRLSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.CRL] {
	tokens := s.ADBERevocationValuesTokens()
	result := make([]spi.RevocationTokenOriginsEntry[revocation.CRL], 0, len(tokens))
	for _, token := range tokens {
		result = append(result, spi.RevocationTokenOriginsEntry[revocation.CRL]{
			Token:   token,
			Origins: []enumerations.RevocationOrigin{enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL},
		})
	}
	return result
}

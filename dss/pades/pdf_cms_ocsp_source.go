// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfCmsOCSPSource.java
// (DSS 6.5.RC1).
//
// BouncyCastle replacements used here (see PORTING.md), following the precedent
// spi/cms_ocsp_source.go established for the sibling CMSOCSPSource:
//   - org.bouncycastle.asn1.cms.AttributeTable -> internal/cmscore.Attributes.
//   - Attribute#getAttributeValues() -> Attribute.Values ([]*asn1ber.Element).
//   - DSSASN1Utils.getAsn1Attributes(AttributeTable, ASN1ObjectIdentifier) ->
//     spi.DSSASN1UtilsAsn1Attributes.
//   - DSSASN1Utils.toBasicOCSPResp(OCSPResponse) has no direct Go port (dss_asn1_utils.go does
//     not carry it); it is reproduced here as the composition spi.DSSRevocationUtilsOcspResp
//     (parses the raw OCSPResponse DER bytes into an OCSPResp) followed by
//     spi.DSSRevocationUtilsFromRespToBasic (extracts the embedded BasicOCSPResponse), the same
//     two calls spi/cms_ocsp_source.go's addBasicOcspRespFromIDRIOcspResponse chains for the
//     structurally identical "full OCSPResponse, not a bare BasicOCSPResponse" case.
//
// FORWARD DEPENDENCIES (not in this chunk's manifest):
//   - RevocationInfoArchival (eu.europa.esig.dss.pades.validation.RevocationInfoArchival) - a
//     struct with OcspVals() [][]byte, each element being the DER encoding of one ASN.1
//     OCSPResponse (RFC 6960) found in the RevocationInfoArchival ASN.1 SEQUENCE's [1] member.
//   - PAdESUtilsRevocationInfoArchival(attrValue *asn1ber.Element) *RevocationInfoArchival - the
//     flattened static PAdESUtils.getRevocationInfoArchival(ASN1Encodable), following the
//     PAdESUtilsVRIsWithName precedent (pdf_composite_dss_dict_certificate_source.go) for how a
//     flattened PAdESUtils static method is named in this port; nil stands for Java's null
//     (RevocationInfoArchival.getInstance returns null for an unrecognised shape).
package pades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PdfCmsOCSPSource represents a source of OCSP tokens extracted from a PDF's CMS.
type PdfCmsOCSPSource struct {
	spi.OfflineOCSPSourceBase
}

// NewPdfCmsOCSPSource is the default constructor. Port of the constructor
// PdfCmsOCSPSource(AttributeTable).
func NewPdfCmsOCSPSource(signedAttributes cmscore.Attributes) *PdfCmsOCSPSource {
	source := &PdfCmsOCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's inherited RevocationTokens implementation.
	source.InitOfflineRevocationSource(source)
	source.extractOCSPArchivalValues(signedAttributes)
	return source
}

// extractOCSPArchivalValues ports the private extractOCSPArchivalValues(AttributeTable).
func (s *PdfCmsOCSPSource) extractOCSPArchivalValues(signedAttributes cmscore.Attributes) {
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
// DEVIATION: upstream logs a warning ("Error while extracting OCSPResponse from Revocation Info
// Archivals (ADBE) : {}") and continues when DSSASN1Utils.toBasicOCSPResp raises an
// OCSPException; slf4j logging is dropped per PORTING.md, so the entry is simply skipped here,
// same as every malformed-revocation degradation elsewhere in this port.
func (s *PdfCmsOCSPSource) extractRevocationInfoArchival(attrValue *asn1ber.Element) {
	revocationArchival := PAdESUtilsRevocationInfoArchival(attrValue)
	if revocationArchival != nil {
		for _, encodedOcspResponse := range revocationArchival.OcspVals() {
			ocspResp := spi.DSSRevocationUtilsOcspResp(encodedOcspResponse)
			if ocspResp == nil {
				continue
			}
			basicOCSPResp := spi.DSSRevocationUtilsFromRespToBasic(ocspResp)
			if basicOCSPResp == nil {
				continue
			}
			ocspResponseIdentifier, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
			if err != nil {
				continue
			}
			s.AddBinary(ocspResponseIdentifier, enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL)
		}
	}
}

// PdfCmsOCSPSourceBinaryOriginsEntry pairs an OCSP binary with the origins it has been found
// with, standing in for one entry of the inherited
// Map<EncapsulatedRevocationTokenIdentifier<OCSP>, Set<RevocationOrigin>>.
type PdfCmsOCSPSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]
	Origins []enumerations.RevocationOrigin
}

// AllRevocationBinariesWithOrigins returns all the OCSP binaries with their origins.
//
// Port of the inherited OfflineRevocationSource#getAllRevocationBinariesWithOrigins(), which
// upstream PdfCmsOCSPSource does not override: Java's abstract base answers its internal
// binaryOrigins bookkeeping map directly, but that generic bookkeeping was not carried over into
// spi.OfflineRevocationSourceBase (see pdf_dss_dict_crl_source.go's header for the same gap and
// pades_crl_source.go/pades_ocsp_source.go for the composing callers). This class only ever
// calls AddBinary with a single origin (RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL, see
// extractRevocationInfoArchival above), so the exported ADBERevocationValuesBinaries() - which
// spi.OfflineRevocationSourceBase does carry, filtered on exactly that origin - answers the
// identical binary set the generic bookkeeping map would, each paired with that single origin.
func (s *PdfCmsOCSPSource) AllRevocationBinariesWithOrigins() []PdfCmsOCSPSourceBinaryOriginsEntry {
	binaries := s.ADBERevocationValuesBinaries()
	result := make([]PdfCmsOCSPSourceBinaryOriginsEntry, 0, len(binaries))
	for _, binary := range binaries {
		result = append(result, PdfCmsOCSPSourceBinaryOriginsEntry{
			Binary:  binary,
			Origins: []enumerations.RevocationOrigin{enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL},
		})
	}
	return result
}

// AllRevocationTokensWithOrigins returns all the OCSP tokens with their origins. Port of the
// inherited OfflineRevocationSource#getAllRevocationTokensWithOrigins(); see
// AllRevocationBinariesWithOrigins's doc comment for why ADBERevocationValuesTokens() (also
// carried over into spi.OfflineRevocationSourceBase) answers the equivalent result for this
// single-origin source.
func (s *PdfCmsOCSPSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.OCSP] {
	tokens := s.ADBERevocationValuesTokens()
	result := make([]spi.RevocationTokenOriginsEntry[revocation.OCSP], 0, len(tokens))
	for _, token := range tokens {
		result = append(result, spi.RevocationTokenOriginsEntry[revocation.OCSP]{
			Token:   token,
			Origins: []enumerations.RevocationOrigin{enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL},
		})
	}
	return result
}

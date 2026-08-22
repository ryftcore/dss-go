// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESOCSPSource.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart). Java's Map<..., Set<RevocationOrigin>>
// return values become slices of pairs, as in pades_crl_source.go's
// CRLSourceBinaryOriginsEntry for the sibling CRLSource (see that file's header for
// why getAllRevocationBinariesWithOrigins()/getAllRevocationTokensWithOrigins() are each
// source's own method rather than an override of a base one in this port).
package pades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// OCSPSource is an OCSPSource that retrieves the OCSPResp from a PAdES Signature.
type OCSPSource struct {
	spi.OfflineOCSPSourceBase

	// cmsOCSPSource is the CMS OCSP source.
	cmsOCSPSource *PdfCmsOCSPSource

	// dssDictOCSPSource is the DSS dictionary OCSP source.
	dssDictOCSPSource *PdfDssDictOCSPSource
}

// NewPAdESOCSPSource is the default constructor. Port of the constructor
// OCSPSource(PdfSignatureRevision, String, AttributeTable).
//
// Panics with the Java message when vriDictionaryName is empty (Objects.requireNonNull; the
// empty string means no VRI-name filter, see pdf_dss_dict_crl_source.go).
func NewPAdESOCSPSource(pdfSignatureRevision *PdfSignatureRevision, vriDictionaryName string,
	signedAttributes cmscore.Attributes) *OCSPSource {
	if vriDictionaryName == "" {
		panic("vriDictionaryName cannot be null!")
	}

	source := &OCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
		cmsOCSPSource:         NewPdfCmsOCSPSource(signedAttributes),
		dssDictOCSPSource: NewPdfDssDictOCSPSourceWithVRIName(
			pdfSignatureRevision.CompositeDssDictionary().OcspSource(),
			pdfSignatureRevision.DssDictionary(), vriDictionaryName),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// RevocationTokens returns the OCSP tokens found by both the CMS source and the DSS dictionary
// source. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *OCSPSource) RevocationTokens(certificateToken, issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.OCSP], error) {
	revocationTokens := make([]spi.RevocationToken[revocation.OCSP], 0)
	cmsTokens, err := s.cmsOCSPSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = append(revocationTokens, cmsTokens...)
	dssDictTokens, err := s.dssDictOCSPSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = append(revocationTokens, dssDictTokens...)
	return revocationTokens, nil
}

// AllRevocationBinaries retrieves all found revocation binaries. Port of the inherited
// getAllRevocationBinaries(); see pades_crl_source.go's CRLSource.AllRevocationBinaries doc
// comment - the identical CRL/OCSP asymmetry applies here (Java: virtual dispatch to THIS type's
// own getAllRevocationBinariesWithOrigins() override below, which the promoted
// spi.OfflineRevocationSourceBase.AllRevocationBinaries this type would otherwise inherit cannot
// reach).
func (s *OCSPSource) AllRevocationBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	entries := s.AllRevocationBinariesWithOrigins()
	result := make([]spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP], 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Binary)
	}
	return result
}

// AllRevocationTokens retrieves a slice of all found RevocationTokens. Port of the inherited
// getAllRevocationTokens(); see AllRevocationBinaries's doc comment.
func (s *OCSPSource) AllRevocationTokens() []spi.RevocationToken[revocation.OCSP] {
	entries := s.AllRevocationTokensWithOrigins()
	result := make([]spi.RevocationToken[revocation.OCSP], 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Token)
	}
	return result
}

// IsEmpty checks if the current source is empty. Port of the inherited isEmpty(); see
// pades_crl_source.go's CRLSource.IsEmpty doc comment.
func (s *OCSPSource) IsEmpty() bool {
	return len(s.AllRevocationBinariesWithOrigins()) == 0 &&
		len(s.AllRevocationTokensWithOrigins()) == 0 &&
		len(s.AllRevocationReferences()) == 0
}

// OcspMap returns a map of all OCSP entries contained in DSS dictionary or into nested VRI
// dictionaries. Port of getOcspMap().
func (s *OCSPSource) OcspMap() map[PdfObjectKey]*spi.OCSPResponseBinary {
	return s.dssDictOCSPSource.OcspMap()
}

// DSSDictionaryBinaries returns the OCSP binaries of the /DSS dictionary. Port of the
// getDSSDictionaryBinaries() override.
func (s *OCSPSource) DSSDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	return s.dssDictOCSPSource.DSSDictionaryBinaries()
}

// DSSDictionaryTokens returns the OCSP tokens of the /DSS dictionary. Port of the
// getDSSDictionaryTokens() override.
func (s *OCSPSource) DSSDictionaryTokens() []spi.RevocationToken[revocation.OCSP] {
	return s.dssDictOCSPSource.DSSDictionaryTokens()
}

// VRIDictionaryBinaries returns the OCSP binaries of the /VRI dictionaries. Port of the
// getVRIDictionaryBinaries() override.
func (s *OCSPSource) VRIDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	return s.dssDictOCSPSource.VRIDictionaryBinaries()
}

// VRIDictionaryTokens returns the OCSP tokens of the /VRI dictionaries. Port of the
// getVRIDictionaryTokens() override.
func (s *OCSPSource) VRIDictionaryTokens() []spi.RevocationToken[revocation.OCSP] {
	return s.dssDictOCSPSource.VRIDictionaryTokens()
}

// ADBERevocationValuesBinaries returns the OCSP binaries found in the ADBE revocation info
// archival CMS attribute. Port of the getADBERevocationValuesBinaries() override.
//
// The base OfflineOCSPSourceBase already implements this by filtering on
// RevocationOriginAdbeRevocationInfoArchival, which is the only origin this class's
// constructor ever adds to; this override is the promoted method inherited from
// spi.OfflineOCSPSourceBase (embedded via cmsOCSPSource), reproducing the Java override
// unchanged.
func (s *OCSPSource) ADBERevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	return s.cmsOCSPSource.ADBERevocationValuesBinaries()
}

// ADBERevocationValuesTokens returns the OCSP tokens found in the ADBE revocation info archival
// CMS attribute. Port of the getADBERevocationValuesTokens() override.
func (s *OCSPSource) ADBERevocationValuesTokens() []spi.RevocationToken[revocation.OCSP] {
	return s.dssDictOCSPSource.ADBERevocationValuesTokens()
}

// OCSPSourceBinaryOriginsEntry pairs an OCSP binary with the origins it has been found
// with, standing in for one entry of Java's
// Map<EncapsulatedRevocationTokenIdentifier<OCSP>, Set<RevocationOrigin>>.
type OCSPSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]
	Origins []enumerations.RevocationOrigin
}

// AllRevocationBinariesWithOrigins returns a map of all revocation binaries with the
// corresponding origins. Port of the getAllRevocationBinariesWithOrigins() override, together
// with the private populateMapWithSet(Map, Map) helper it calls twice.
func (s *OCSPSource) AllRevocationBinariesWithOrigins() []OCSPSourceBinaryOriginsEntry {
	result := make([]OCSPSourceBinaryOriginsEntry, 0)
	for _, entry := range s.cmsOCSPSource.AllRevocationBinariesWithOrigins() {
		result = padesOCSPSourceMergeBinaryOrigins(result, entry.Binary, entry.Origins)
	}
	for _, entry := range s.dssDictOCSPSource.AllRevocationBinariesWithOrigins() {
		result = padesOCSPSourceMergeBinaryOrigins(result, entry.Binary, entry.Origins)
	}
	return result
}

// AllRevocationTokensWithOrigins returns a map of all revocation tokens with the corresponding
// origins. Port of the getAllRevocationTokensWithOrigins() override, together with the private
// populateMapWithSet(Map, Map) helper it calls twice.
func (s *OCSPSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.OCSP] {
	result := make([]spi.RevocationTokenOriginsEntry[revocation.OCSP], 0)
	for _, entry := range s.cmsOCSPSource.AllRevocationTokensWithOrigins() {
		result = padesOCSPSourceMergeTokenOrigins(result, entry.Token, entry.Origins)
	}
	for _, entry := range s.dssDictOCSPSource.AllRevocationTokensWithOrigins() {
		result = padesOCSPSourceMergeTokenOrigins(result, entry.Token, entry.Origins)
	}
	return result
}

// padesOCSPSourceMergeBinaryOrigins merges one (binary, origins) pair into result, unioning the
// origins of an already-present entry for the same binary (matched on its DSS identifier,
// AsXmlID()) rather than duplicating it - the Set<RevocationOrigin> semantics of the private
// populateMapWithSet(Map, Map) helper's Map<..., Set<RevocationOrigin>> value.
func padesOCSPSourceMergeBinaryOrigins(result []OCSPSourceBinaryOriginsEntry,
	binary spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP],
	origins []enumerations.RevocationOrigin) []OCSPSourceBinaryOriginsEntry {
	for i := range result {
		if result[i].Binary.AsXmlID() == binary.AsXmlID() {
			result[i].Origins = padesOCSPSourceUnionOrigins(result[i].Origins, origins)
			return result
		}
	}
	merged := append([]enumerations.RevocationOrigin{}, origins...)
	return append(result, OCSPSourceBinaryOriginsEntry{Binary: binary, Origins: merged})
}

// padesOCSPSourceMergeTokenOrigins merges one (token, origins) pair into result, unioning the
// origins of an already-present entry for the same token (matched on its DSS identifier,
// DSSIDAsString()) rather than duplicating it. Same Set<RevocationOrigin> semantics as
// padesOCSPSourceMergeBinaryOrigins.
func padesOCSPSourceMergeTokenOrigins(result []spi.RevocationTokenOriginsEntry[revocation.OCSP],
	token spi.RevocationToken[revocation.OCSP],
	origins []enumerations.RevocationOrigin) []spi.RevocationTokenOriginsEntry[revocation.OCSP] {
	for i := range result {
		if result[i].Token.DSSIDAsString() == token.DSSIDAsString() {
			result[i].Origins = padesOCSPSourceUnionOrigins(result[i].Origins, origins)
			return result
		}
	}
	merged := append([]enumerations.RevocationOrigin{}, origins...)
	return append(result, spi.RevocationTokenOriginsEntry[revocation.OCSP]{Token: token, Origins: merged})
}

// padesOCSPSourceUnionOrigins appends every element of toAdd not already present in existing,
// reproducing Set<RevocationOrigin>#addAll semantics.
func padesOCSPSourceUnionOrigins(existing, toAdd []enumerations.RevocationOrigin) []enumerations.RevocationOrigin {
	for _, origin := range toAdd {
		found := false
		for _, e := range existing {
			if e == origin {
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, origin)
		}
	}
	return existing
}

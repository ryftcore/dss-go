// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESCRLSource.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart). Java's Map<..., Set<RevocationOrigin>>
// return values become slices of pairs, the convention pdf_dss_dict_crl_source.go established
// for the sibling PdfDssDictCRLSource with PdfDssDictCRLSourceBinaryOriginsEntry. Note (per that
// same file's header) that getAllRevocationBinariesWithOrigins()/getAllRevocationTokensWithOrigins()
// are declared by the Java OfflineRevocationSource base but were not carried over into
// spi.OfflineRevocationSourceBase, so every concrete source that needs them (this one included)
// defines its own copy of the merge, rather than overriding an inherited one.
//
// FORWARD DEPENDENCIES (not in this chunk's manifest):
//
//   - PdfSignatureRevision (eu.europa.esig.dss.pdf.PdfSignatureRevision); see
//     pades_certificate_source.go's header for its full assumed shape - this file additionally
//     uses CompositeDssDictionary().CrlSource() and DssDictionary().
//
//   - PdfCmsCRLSource (eu.europa.esig.dss.pades.validation.PdfCmsCRLSource), the CRL-flavoured
//     sibling of this chunk's PdfCmsOCSPSource (pdf_cms_ocsp_source.go), extending
//     spi.OfflineCRLSourceBase the same way. Assumed shape, mirroring PdfCmsOCSPSource and the
//     RevocationOrigin.ADBE_REVOCATION_INFO_ARCHIVAL-only content this class ever carries:
//
//     func NewPdfCmsCRLSource(signedAttributes cmscore.Attributes) *PdfCmsCRLSource
//     func (s *PdfCmsCRLSource) RevocationTokens(certificateToken, issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.CRL], error)
//     func (s *PdfCmsCRLSource) ADBERevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]
//     func (s *PdfCmsCRLSource) AllRevocationBinariesWithOrigins() []PdfCmsCRLSourceBinaryOriginsEntry
//     func (s *PdfCmsCRLSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.CRL]
//
//     where PdfCmsCRLSourceBinaryOriginsEntry mirrors PdfDssDictCRLSourceBinaryOriginsEntry's
//     {Binary, Origins} shape (pdf_dss_dict_crl_source.go).
package pades

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PAdESCRLSource is a CRLSource that will retrieve the CRL from a PAdES Signature.
type PAdESCRLSource struct {
	spi.OfflineCRLSourceBase

	// cmsCrlSource is the CMS CRL source.
	cmsCrlSource *PdfCmsCRLSource

	// dssDictCrlSource is the DSS dictionary CRL source.
	dssDictCrlSource *PdfDssDictCRLSource
}

// NewPAdESCRLSource is the default constructor. Port of the constructor
// PAdESCRLSource(PdfSignatureRevision, String, AttributeTable).
//
// Panics with the Java message when vriDictionaryName is empty (Objects.requireNonNull; the
// empty string stands for Java's null throughout this port, see pdf_dss_dict_crl_source.go).
func NewPAdESCRLSource(pdfSignatureRevision *PdfSignatureRevision, vriDictionaryName string,
	signedAttributes cmscore.Attributes) *PAdESCRLSource {
	if vriDictionaryName == "" {
		panic("vriDictionaryName cannot be null!")
	}

	source := &PAdESCRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		cmsCrlSource:         NewPdfCmsCRLSource(signedAttributes),
		dssDictCrlSource: NewPdfDssDictCRLSourceWithVRIName(
			pdfSignatureRevision.CompositeDssDictionary().CrlSource(),
			pdfSignatureRevision.DssDictionary(), vriDictionaryName),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// RevocationTokens returns the CRL tokens found by both the CMS source and the DSS dictionary
// source. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *PAdESCRLSource) RevocationTokens(certificateToken, issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.CRL], error) {
	revocationTokens := make([]spi.RevocationToken[revocation.CRL], 0)
	cmsTokens, err := s.cmsCrlSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = append(revocationTokens, cmsTokens...)
	dssDictTokens, err := s.dssDictCrlSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = append(revocationTokens, dssDictTokens...)
	return revocationTokens, nil
}

// AllRevocationBinaries retrieves all found revocation binaries. Port of the inherited
// getAllRevocationBinaries() (OfflineRevocationSource base): Java implements it as
// getAllRevocationBinariesWithOrigins().keySet(), which virtual dispatch resolves to THIS type's
// own getAllRevocationBinariesWithOrigins() override above. The promoted
// spi.OfflineRevocationSourceBase.AllRevocationBinaries this type would otherwise inherit reads
// its own binaryOrigins field directly - populated only by RevocationTokens(cert, issuer) calls,
// which nothing makes for every DSS-dictionary-embedded CRL up front - so it never sees the DSS
// dictionary/VRI-sourced binaries AllRevocationBinariesWithOrigins already exposes
// unconditionally. Shadowing here (Go method redefinition standing in for Java's virtual
// dispatch; see PORTING.md's "Virtual dispatch" precedent, also used by
// pades_certificate_source.go's DSSDictionaryCertValues) is required for every caller reaching
// this type through the spi.OfflineRevocationSource[R] interface (e.g.
// PAdESSignature.CompleteCRLSource(), and thus BaselineRequirementsChecker.MinimalLTRequirement's
// LT-level revocation-presence check) to see the DSS dictionary's CRLs at all.
func (s *PAdESCRLSource) AllRevocationBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	entries := s.AllRevocationBinariesWithOrigins()
	result := make([]spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL], 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Binary)
	}
	return result
}

// AllRevocationTokens retrieves a slice of all found RevocationTokens. Port of the inherited
// getAllRevocationTokens(); see AllRevocationBinaries's doc comment for why this needs the same
// shadowing treatment (Java: getAllRevocationTokensWithOrigins().keySet()).
func (s *PAdESCRLSource) AllRevocationTokens() []spi.RevocationToken[revocation.CRL] {
	entries := s.AllRevocationTokensWithOrigins()
	result := make([]spi.RevocationToken[revocation.CRL], 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Token)
	}
	return result
}

// IsEmpty checks if the current source is empty. Port of the inherited isEmpty(); see
// AllRevocationBinaries's doc comment for why this needs the same shadowing treatment (Java:
// Utils.isMapEmpty(getAllRevocationBinariesWithOrigins()) &&
// Utils.isMapEmpty(getAllRevocationTokensWithOrigins()) &&
// Utils.isMapEmpty(getRevocationReferencesWithOrigins()) - the last of which this port's base
// still answers correctly, since PAdES neither overrides it nor ever populates references).
func (s *PAdESCRLSource) IsEmpty() bool {
	return len(s.AllRevocationBinariesWithOrigins()) == 0 &&
		len(s.AllRevocationTokensWithOrigins()) == 0 &&
		len(s.AllRevocationReferences()) == 0
}

// CrlMap returns a map of all CRL entries contained in DSS dictionary or into nested VRI
// dictionaries. Port of getCrlMap().
func (s *PAdESCRLSource) CrlMap() map[PdfObjectKey]*crlparser.CRLBinary {
	return s.dssDictCrlSource.CrlMap()
}

// DSSDictionaryBinaries returns the CRL binaries of the /DSS dictionary. Port of the
// getDSSDictionaryBinaries() override.
func (s *PAdESCRLSource) DSSDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	return s.dssDictCrlSource.DSSDictionaryBinaries()
}

// DSSDictionaryTokens returns the CRL tokens of the /DSS dictionary. Port of the
// getDSSDictionaryTokens() override.
func (s *PAdESCRLSource) DSSDictionaryTokens() []spi.RevocationToken[revocation.CRL] {
	return s.dssDictCrlSource.DSSDictionaryTokens()
}

// VRIDictionaryBinaries returns the CRL binaries of the /VRI dictionaries. Port of the
// getVRIDictionaryBinaries() override.
func (s *PAdESCRLSource) VRIDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	return s.dssDictCrlSource.VRIDictionaryBinaries()
}

// VRIDictionaryTokens returns the CRL tokens of the /VRI dictionaries. Port of the
// getVRIDictionaryTokens() override.
func (s *PAdESCRLSource) VRIDictionaryTokens() []spi.RevocationToken[revocation.CRL] {
	return s.dssDictCrlSource.VRIDictionaryTokens()
}

// ADBERevocationValuesBinaries returns the CRL binaries found in the ADBE revocation info
// archival CMS attribute. Port of the getADBERevocationValuesBinaries() override.
func (s *PAdESCRLSource) ADBERevocationValuesBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	return s.cmsCrlSource.ADBERevocationValuesBinaries()
}

// ADBERevocationValuesTokens returns the CRL tokens found in the ADBE revocation info archival
// CMS attribute. Port of the getADBERevocationValuesTokens() override.
func (s *PAdESCRLSource) ADBERevocationValuesTokens() []spi.RevocationToken[revocation.CRL] {
	return s.dssDictCrlSource.ADBERevocationValuesTokens()
}

// PAdESCRLSourceBinaryOriginsEntry pairs a CRL binary with the origins it has been found with,
// standing in for one entry of Java's
// Map<EncapsulatedRevocationTokenIdentifier<CRL>, Set<RevocationOrigin>>.
type PAdESCRLSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]
	Origins []enumerations.RevocationOrigin
}

// AllRevocationBinariesWithOrigins returns a map of all revocation binaries with the
// corresponding origins. Port of the getAllRevocationBinariesWithOrigins() override, together
// with the private populateMapWithSet(Map, Map) helper it calls twice.
func (s *PAdESCRLSource) AllRevocationBinariesWithOrigins() []PAdESCRLSourceBinaryOriginsEntry {
	result := make([]PAdESCRLSourceBinaryOriginsEntry, 0)
	for _, entry := range s.cmsCrlSource.AllRevocationBinariesWithOrigins() {
		result = padesCRLSourceMergeBinaryOrigins(result, entry.Binary, entry.Origins)
	}
	for _, entry := range s.dssDictCrlSource.AllRevocationBinariesWithOrigins() {
		result = padesCRLSourceMergeBinaryOrigins(result, entry.Binary, entry.Origins)
	}
	return result
}

// AllRevocationTokensWithOrigins returns a map of all revocation tokens with the corresponding
// origins. Port of the getAllRevocationTokensWithOrigins() override, together with the private
// populateMapWithSet(Map, Map) helper it calls twice.
func (s *PAdESCRLSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.CRL] {
	result := make([]spi.RevocationTokenOriginsEntry[revocation.CRL], 0)
	for _, entry := range s.cmsCrlSource.AllRevocationTokensWithOrigins() {
		result = padesCRLSourceMergeTokenOrigins(result, entry.Token, entry.Origins)
	}
	for _, entry := range s.dssDictCrlSource.AllRevocationTokensWithOrigins() {
		result = padesCRLSourceMergeTokenOrigins(result, entry.Token, entry.Origins)
	}
	return result
}

// padesCRLSourceMergeBinaryOrigins merges one (binary, origins) pair into result, unioning the
// origins of an already-present entry for the same binary (matched on its DSS identifier,
// AsXmlID()) rather than duplicating it - the Set<RevocationOrigin> semantics of the private
// populateMapWithSet(Map, Map) helper's Map<..., Set<RevocationOrigin>> value.
func padesCRLSourceMergeBinaryOrigins(result []PAdESCRLSourceBinaryOriginsEntry,
	binary spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL],
	origins []enumerations.RevocationOrigin) []PAdESCRLSourceBinaryOriginsEntry {
	for i := range result {
		if result[i].Binary.AsXmlID() == binary.AsXmlID() {
			result[i].Origins = padesCRLSourceUnionOrigins(result[i].Origins, origins)
			return result
		}
	}
	merged := append([]enumerations.RevocationOrigin{}, origins...)
	return append(result, PAdESCRLSourceBinaryOriginsEntry{Binary: binary, Origins: merged})
}

// padesCRLSourceMergeTokenOrigins merges one (token, origins) pair into result, unioning the
// origins of an already-present entry for the same token (matched on its DSS identifier,
// DSSIDAsString()) rather than duplicating it. Same Set<RevocationOrigin> semantics as
// padesCRLSourceMergeBinaryOrigins.
func padesCRLSourceMergeTokenOrigins(result []spi.RevocationTokenOriginsEntry[revocation.CRL],
	token spi.RevocationToken[revocation.CRL],
	origins []enumerations.RevocationOrigin) []spi.RevocationTokenOriginsEntry[revocation.CRL] {
	for i := range result {
		if result[i].Token.DSSIDAsString() == token.DSSIDAsString() {
			result[i].Origins = padesCRLSourceUnionOrigins(result[i].Origins, origins)
			return result
		}
	}
	merged := append([]enumerations.RevocationOrigin{}, origins...)
	return append(result, spi.RevocationTokenOriginsEntry[revocation.CRL]{Token: token, Origins: merged})
}

// padesCRLSourceUnionOrigins appends every element of toAdd not already present in existing,
// reproducing Set<RevocationOrigin>#addAll semantics.
func padesCRLSourceUnionOrigins(existing, toAdd []enumerations.RevocationOrigin) []enumerations.RevocationOrigin {
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

// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfDssDictCRLSource.java (DSS 6.5.RC1).
//
// The forward dependencies this file shares with the rest of the chunk (PdfObjectKey,
// PdfDssDict, PdfVriDict, PAdESUtilsVRIsWithName) and the determinism rule applied to
// PDF-object-keyed map walks are documented in pdf_composite_dss_dict_certificate_source.go.
//
// Java's Map<..., Set<RevocationOrigin>> return values become slices of pairs, the convention
// spi.OfflineRevocationSourceBase established with RevocationTokenOriginsEntry. Note that
// getAllRevocationBinariesWithOrigins() and getAllRevocationTokensWithOrigins() are declared by
// the Java OfflineRevocationSource base but were not carried over into
// spi.OfflineRevocationSourceBase, so in Go they are new methods of this type rather than
// overrides; nothing dispatches to them through the base.
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PdfDssDictCRLSourceBinaryOriginsEntry pairs a CRL binary with the origins it has been found
// with, standing in for one entry of Java's
// Map<EncapsulatedRevocationTokenIdentifier<CRL>, Set<RevocationOrigin>>.
type PdfDssDictCRLSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]
	Origins []enumerations.RevocationOrigin
}

// PdfDssDictCRLSource is the CRL source extracted from a DSS dictionary.
type PdfDssDictCRLSource struct {
	spi.OfflineCRLSourceBase

	// compositeCRLSource is the merged CRL source combined from all /DSS revisions.
	compositeCRLSource *PdfCompositeDssDictCRLSource

	// dssDictionary is the DSS dictionary.
	dssDictionary PdfDssDict

	// relatedVRIDictionaryName is the name of the signature's VRI dictionary, when applicable;
	// the empty string stands for Java's null, i.e. "every VRI dictionary".
	relatedVRIDictionaryName string

	// crlMap is the cached CRL map, nil until CrlMap() computes it.
	crlMap map[PdfObjectKey]*crlparser.CRLBinary
}

// NewPdfDssDictCRLSource is the port of the
// PdfDssDictCRLSource(PdfCompositeDssDictCRLSource, PdfDssDict) constructor.
func NewPdfDssDictCRLSource(compositeCRLSource *PdfCompositeDssDictCRLSource,
	dssDictionary PdfDssDict) *PdfDssDictCRLSource {
	return NewPdfDssDictCRLSourceWithVRIName(compositeCRLSource, dssDictionary, "")
}

// NewPdfDssDictCRLSourceWithVRIName is the port of the
// PdfDssDictCRLSource(PdfCompositeDssDictCRLSource, PdfDssDict, String) constructor, to be used
// for a signature source. vriDictionaryName is the SHA-1 of the signature name; the empty
// string stands for Java's null.
func NewPdfDssDictCRLSourceWithVRIName(compositeCRLSource *PdfCompositeDssDictCRLSource,
	dssDictionary PdfDssDict, vriDictionaryName string) *PdfDssDictCRLSource {
	source := &PdfDssDictCRLSource{
		OfflineCRLSourceBase:     spi.NewOfflineCRLSourceBase(),
		compositeCRLSource:       compositeCRLSource,
		dssDictionary:            dssDictionary,
		relatedVRIDictionaryName: vriDictionaryName,
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// CrlMap returns a map of all the CRL entries contained in the DSS dictionary or in the nested
// VRI dictionaries, with their object ids. Port of getCrlMap().
func (s *PdfDssDictCRLSource) CrlMap() map[PdfObjectKey]*crlparser.CRLBinary {
	if s.crlMap == nil {
		s.crlMap = make(map[PdfObjectKey]*crlparser.CRLBinary)
		if s.dssDictionary != nil {
			for key, crlBinary := range s.dssDictionary.CRLs() {
				s.crlMap[key] = crlBinary
			}
			vriDicts := PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
			for _, vriDict := range vriDicts {
				for key, crlBinary := range vriDict.CRLs() {
					s.crlMap[key] = crlBinary
				}
			}
		}
	}
	return s.crlMap
}

// RevocationTokens returns the CRL tokens of the composite source that belong to this
// dictionary, followed by the ones the offline base builds from this source's own binaries.
// Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *PdfDssDictCRLSource) RevocationTokens(certificateToken *model.CertificateToken,
	issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.CRL], error) {
	revocationTokens, err := s.compositeCRLSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = s.filterTokensFromCrlMap(revocationTokens)
	superTokens, err := s.OfflineCRLSourceBase.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	revocationTokens = append(revocationTokens, superTokens...)
	return revocationTokens, nil
}

// DSSDictionaryBinaries returns the CRL binaries of the /DSS dictionary.
// Port of the getDSSDictionaryBinaries() override.
func (s *PdfDssDictCRLSource) DSSDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	if s.dssDictionary != nil {
		return s.filterBinariesFromKeys(s.compositeCRLSource.DSSDictionaryBinaries(),
			pdfDssDictCRLSourceSortedKeys(s.dssDictionary.CRLs()))
	}
	return []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]{}
}

// DSSDictionaryTokens returns the CRL tokens of the /DSS dictionary.
// Port of the getDSSDictionaryTokens() override.
func (s *PdfDssDictCRLSource) DSSDictionaryTokens() []spi.RevocationToken[revocation.CRL] {
	if s.dssDictionary != nil {
		return s.filterTokensFromKeys(s.compositeCRLSource.DSSDictionaryTokens(),
			pdfDssDictCRLSourceSortedKeys(s.dssDictionary.CRLs()))
	}
	return []spi.RevocationToken[revocation.CRL]{}
}

// VRIDictionaryBinaries returns the CRL binaries of the /VRI dictionaries.
// Port of the getVRIDictionaryBinaries() override.
func (s *PdfDssDictCRLSource) VRIDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	if s.dssDictionary != nil {
		return s.filterBinariesFromKeys(s.compositeCRLSource.VRIDictionaryBinaries(), s.keySetFromVRIDictionaries())
	}
	return []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]{}
}

// VRIDictionaryTokens returns the CRL tokens of the /VRI dictionaries.
// Port of the getVRIDictionaryTokens() override.
func (s *PdfDssDictCRLSource) VRIDictionaryTokens() []spi.RevocationToken[revocation.CRL] {
	if s.dssDictionary != nil {
		return s.filterTokensFromKeys(s.compositeCRLSource.VRIDictionaryTokens(), s.keySetFromVRIDictionaries())
	}
	return []spi.RevocationToken[revocation.CRL]{}
}

// keySetFromVRIDictionaries ports the private getKeySetFromVRIDictionaries().
func (s *PdfDssDictCRLSource) keySetFromVRIDictionaries() []PdfObjectKey {
	if s.dssDictionary != nil {
		result := make([]PdfObjectKey, 0)
		vris := PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
		for _, vriDict := range vris {
			for _, key := range pdfDssDictCRLSourceSortedKeys(vriDict.CRLs()) {
				if !pdfDssDictCRLSourceContainsKey(result, key) {
					result = append(result, key)
				}
			}
		}
		return result
	}
	return []PdfObjectKey{}
}

// filterBinariesFromKeys ports the private filterBinariesFromKeys(Collection, Collection).
func (s *PdfDssDictCRLSource) filterBinariesFromKeys(
	crlBinaries []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL],
	keySet []PdfObjectKey) []spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL] {
	result := make([]spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL], 0)
	for _, crlBinary := range crlBinaries {
		objectIDs := s.compositeCRLSource.TokenBinaryObjectIDs(crlBinary)
		if utils.ContainsAny(keySet, objectIDs) {
			result = append(result, crlBinary)
		}
	}
	return result
}

// filterTokensFromCrlMap ports the private filterTokensFromCrlMap(List).
func (s *PdfDssDictCRLSource) filterTokensFromCrlMap(revocationTokens []spi.RevocationToken[revocation.CRL]) []spi.RevocationToken[revocation.CRL] {
	return s.filterTokensFromKeys(revocationTokens, pdfDssDictCRLSourceSortedKeys(s.CrlMap()))
}

// filterTokensFromKeys ports the private filterTokensFromKeys(Collection, Collection).
func (s *PdfDssDictCRLSource) filterTokensFromKeys(revocationTokens []spi.RevocationToken[revocation.CRL],
	keySet []PdfObjectKey) []spi.RevocationToken[revocation.CRL] {
	result := make([]spi.RevocationToken[revocation.CRL], 0)
	for _, crlToken := range revocationTokens {
		objectIDs := s.compositeCRLSource.RevocationTokenIDs(crlToken)
		if utils.ContainsAny(keySet, objectIDs) {
			result = append(result, crlToken)
		}
	}
	return result
}

// AllRevocationBinariesWithOrigins returns all the CRL binaries of this dictionary with their
// origins. Port of the getAllRevocationBinariesWithOrigins() override.
func (s *PdfDssDictCRLSource) AllRevocationBinariesWithOrigins() []PdfDssDictCRLSourceBinaryOriginsEntry {
	result := make([]PdfDssDictCRLSourceBinaryOriginsEntry, 0)

	binaries := s.compositeCRLSource.AllRevocationBinaries()
	filteredBinaries := s.filterBinariesFromKeys(binaries, pdfDssDictCRLSourceSortedKeys(s.CrlMap()))
	for _, crlBinary := range filteredBinaries {
		result = append(result, PdfDssDictCRLSourceBinaryOriginsEntry{
			Binary:  crlBinary,
			Origins: s.revocationDataOriginsForBinary(crlBinary),
		})
	}
	return result
}

// revocationDataOriginsForBinary ports the private
// getRevocationDataOrigins(EncapsulatedRevocationTokenIdentifier<CRL>).
//
// As upstream, dssDictionary is dereferenced without a nil check here.
func (s *PdfDssDictCRLSource) revocationDataOriginsForBinary(
	crlBinary spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]) []enumerations.RevocationOrigin {
	result := make([]enumerations.RevocationOrigin, 0)
	tokenBinaryObjectIDs := s.compositeCRLSource.TokenBinaryObjectIDs(crlBinary)
	if utils.ContainsAny(pdfDssDictCRLSourceSortedKeys(s.dssDictionary.CRLs()), tokenBinaryObjectIDs) {
		result = append(result, enumerations.RevocationOriginDSSDictionary)
	}
	for _, vriDict := range PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName) {
		if utils.ContainsAny(pdfDssDictCRLSourceSortedKeys(vriDict.CRLs()), tokenBinaryObjectIDs) {
			if !pdfDssDictCRLSourceContainsOrigin(result, enumerations.RevocationOriginVRIDictionary) {
				result = append(result, enumerations.RevocationOriginVRIDictionary)
			}
		}
	}
	return result
}

// AllRevocationTokensWithOrigins returns all the CRL tokens of this dictionary with their
// origins. Port of the getAllRevocationTokensWithOrigins() override.
func (s *PdfDssDictCRLSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.CRL] {
	result := make([]spi.RevocationTokenOriginsEntry[revocation.CRL], 0)

	tokens := s.compositeCRLSource.AllRevocationTokens()
	filteredTokens := s.filterTokensFromKeys(tokens, pdfDssDictCRLSourceSortedKeys(s.CrlMap()))
	for _, crlToken := range filteredTokens {
		result = append(result, spi.RevocationTokenOriginsEntry[revocation.CRL]{
			Token:   crlToken,
			Origins: s.revocationDataOriginsForToken(crlToken),
		})
	}
	return result
}

// revocationDataOriginsForToken ports the private getRevocationDataOrigins(RevocationToken<CRL>).
//
// As upstream, dssDictionary is dereferenced without a nil check here.
func (s *PdfDssDictCRLSource) revocationDataOriginsForToken(
	crlToken spi.RevocationToken[revocation.CRL]) []enumerations.RevocationOrigin {
	result := make([]enumerations.RevocationOrigin, 0)
	tokenObjectIDs := s.compositeCRLSource.RevocationTokenIDs(crlToken)
	if utils.ContainsAny(pdfDssDictCRLSourceSortedKeys(s.dssDictionary.CRLs()), tokenObjectIDs) {
		result = append(result, enumerations.RevocationOriginDSSDictionary)
	}
	for _, vriDict := range PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName) {
		if utils.ContainsAny(pdfDssDictCRLSourceSortedKeys(vriDict.CRLs()), tokenObjectIDs) {
			if !pdfDssDictCRLSourceContainsOrigin(result, enumerations.RevocationOriginVRIDictionary) {
				result = append(result, enumerations.RevocationOriginVRIDictionary)
			}
		}
	}
	return result
}

// pdfDssDictCRLSourceContainsOrigin reports whether origins already contains origin, standing in
// for the Set<RevocationOrigin> semantics of the Java result.
func pdfDssDictCRLSourceContainsOrigin(origins []enumerations.RevocationOrigin, origin enumerations.RevocationOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}

// pdfDssDictCRLSourceContainsKey reports whether objectIDs already contains objectID.
func pdfDssDictCRLSourceContainsKey(objectIDs []PdfObjectKey, objectID PdfObjectKey) bool {
	for _, id := range objectIDs {
		if id == objectID {
			return true
		}
	}
	return false
}

// pdfDssDictCRLSourceSortedKeys returns the keys of a PDF-object-keyed map in ascending
// (object number, generation) order; see the determinism note in
// pdf_composite_dss_dict_certificate_source.go.
func pdfDssDictCRLSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
	keys := make([]PdfObjectKey, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Number() != keys[j].Number() {
			return keys[i].Number() < keys[j].Number()
		}
		return keys[i].Generation() < keys[j].Generation()
	})
	return keys
}

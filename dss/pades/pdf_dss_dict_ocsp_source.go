// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfDssDictOCSPSource.java (DSS 6.5.RC1).
//
// The forward dependencies this file shares with the rest of the chunk (PdfObjectKey,
// PdfDssDict, PdfVriDict, PAdESUtilsVRIsWithName) and the determinism rule applied to
// PDF-object-keyed map walks are documented in pdf_composite_dss_dict_certificate_source.go.
//
// As in PdfDssDictCRLSource, the Java Map<..., Set<RevocationOrigin>> return values become
// slices of pairs, and AllRevocationBinariesWithOrigins/AllRevocationTokensWithOrigins are new
// methods of this type rather than overrides, since spi.OfflineRevocationSourceBase does not
// carry those two Java getters.
//
// UPSTREAM ASYMMETRY PRESERVED: unlike PdfDssDictCRLSource#getRevocationTokens, this class's
// override does NOT append super.getRevocationTokens(...) to its result.
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PdfDssDictOCSPSourceBinaryOriginsEntry pairs an OCSP binary with the origins it has been found
// with, standing in for one entry of Java's
// Map<EncapsulatedRevocationTokenIdentifier<OCSP>, Set<RevocationOrigin>>.
type PdfDssDictOCSPSourceBinaryOriginsEntry struct {
	Binary  spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]
	Origins []enumerations.RevocationOrigin
}

// PdfDssDictOCSPSource is the OCSP source extracted from a DSS dictionary.
type PdfDssDictOCSPSource struct {
	spi.OfflineOCSPSourceBase

	// compositeOCSPSource is the merged OCSP source combined from all /DSS revisions.
	compositeOCSPSource *PdfCompositeDssDictOCSPSource

	// dssDictionary is the DSS dictionary.
	dssDictionary PdfDssDict

	// relatedVRIDictionaryName is the name of the signature's VRI dictionary, when applicable;
	// the empty string stands for Java's null, i.e. "every VRI dictionary".
	relatedVRIDictionaryName string

	// ocspMap is the cached OCSP map, nil until OcspMap() computes it.
	ocspMap map[PdfObjectKey]*spi.OCSPResponseBinary
}

// NewPdfDssDictOCSPSource is the port of the
// PdfDssDictOCSPSource(PdfCompositeDssDictOCSPSource, PdfDssDict) constructor.
func NewPdfDssDictOCSPSource(compositeOCSPSource *PdfCompositeDssDictOCSPSource,
	dssDictionary PdfDssDict) *PdfDssDictOCSPSource {
	return NewPdfDssDictOCSPSourceWithVRIName(compositeOCSPSource, dssDictionary, "")
}

// NewPdfDssDictOCSPSourceWithVRIName is the port of the
// PdfDssDictOCSPSource(PdfCompositeDssDictOCSPSource, PdfDssDict, String) constructor, to be
// used for a signature source. vriDictionaryName is the SHA-1 of the signature name; the empty
// string stands for Java's null.
func NewPdfDssDictOCSPSourceWithVRIName(compositeOCSPSource *PdfCompositeDssDictOCSPSource,
	dssDictionary PdfDssDict, vriDictionaryName string) *PdfDssDictOCSPSource {
	source := &PdfDssDictOCSPSource{
		OfflineOCSPSourceBase:    spi.NewOfflineOCSPSourceBase(),
		compositeOCSPSource:      compositeOCSPSource,
		dssDictionary:            dssDictionary,
		relatedVRIDictionaryName: vriDictionaryName,
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// OcspMap returns a map of all the OCSP entries contained in the DSS dictionary or in the
// nested VRI dictionaries, with their object ids. Port of getOcspMap().
func (s *PdfDssDictOCSPSource) OcspMap() map[PdfObjectKey]*spi.OCSPResponseBinary {
	if s.ocspMap == nil {
		s.ocspMap = make(map[PdfObjectKey]*spi.OCSPResponseBinary)
		if s.dssDictionary != nil {
			for key, ocspBinary := range s.dssDictionary.OCSPs() {
				s.ocspMap[key] = ocspBinary
			}
			vriDicts := PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
			for _, vriDict := range vriDicts {
				for key, ocspBinary := range vriDict.OCSPs() {
					s.ocspMap[key] = ocspBinary
				}
			}
		}
	}
	return s.ocspMap
}

// RevocationTokens returns the OCSP tokens of the composite source that belong to this
// dictionary. Port of the getRevocationTokens(CertificateToken, CertificateToken) override.
func (s *PdfDssDictOCSPSource) RevocationTokens(certificateToken *model.CertificateToken,
	issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.OCSP], error) {
	revocationTokens, err := s.compositeOCSPSource.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	return s.filterTokensFromOcspMap(revocationTokens), nil
}

// DSSDictionaryBinaries returns the OCSP binaries of the /DSS dictionary.
// Port of the getDSSDictionaryBinaries() override.
func (s *PdfDssDictOCSPSource) DSSDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	if s.dssDictionary != nil {
		return s.filterBinariesFromKeys(s.compositeOCSPSource.DSSDictionaryBinaries(),
			pdfDssDictOCSPSourceSortedKeys(s.dssDictionary.OCSPs()))
	}
	return []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]{}
}

// DSSDictionaryTokens returns the OCSP tokens of the /DSS dictionary.
// Port of the getDSSDictionaryTokens() override.
func (s *PdfDssDictOCSPSource) DSSDictionaryTokens() []spi.RevocationToken[revocation.OCSP] {
	if s.dssDictionary != nil {
		return s.filterTokensFromKeys(s.compositeOCSPSource.DSSDictionaryTokens(),
			pdfDssDictOCSPSourceSortedKeys(s.dssDictionary.OCSPs()))
	}
	return []spi.RevocationToken[revocation.OCSP]{}
}

// VRIDictionaryBinaries returns the OCSP binaries of the /VRI dictionaries.
// Port of the getVRIDictionaryBinaries() override.
func (s *PdfDssDictOCSPSource) VRIDictionaryBinaries() []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	if s.dssDictionary != nil {
		return s.filterBinariesFromKeys(s.compositeOCSPSource.VRIDictionaryBinaries(), s.keySetFromVRIDictionaries())
	}
	return []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]{}
}

// VRIDictionaryTokens returns the OCSP tokens of the /VRI dictionaries.
// Port of the getVRIDictionaryTokens() override.
func (s *PdfDssDictOCSPSource) VRIDictionaryTokens() []spi.RevocationToken[revocation.OCSP] {
	if s.dssDictionary != nil {
		return s.filterTokensFromKeys(s.compositeOCSPSource.VRIDictionaryTokens(), s.keySetFromVRIDictionaries())
	}
	return []spi.RevocationToken[revocation.OCSP]{}
}

// keySetFromVRIDictionaries ports the private getKeySetFromVRIDictionaries().
func (s *PdfDssDictOCSPSource) keySetFromVRIDictionaries() []PdfObjectKey {
	if s.dssDictionary != nil {
		result := make([]PdfObjectKey, 0)
		vris := PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
		for _, vriDict := range vris {
			for _, key := range pdfDssDictOCSPSourceSortedKeys(vriDict.OCSPs()) {
				if !pdfDssDictOCSPSourceContainsKey(result, key) {
					result = append(result, key)
				}
			}
		}
		return result
	}
	return []PdfObjectKey{}
}

// filterBinariesFromKeys ports the private filterBinariesFromKeys(Collection, Collection).
func (s *PdfDssDictOCSPSource) filterBinariesFromKeys(
	ocspBinaries []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP],
	keySet []PdfObjectKey) []spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP] {
	result := make([]spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP], 0)
	for _, ocspBinary := range ocspBinaries {
		objectIDs := s.compositeOCSPSource.TokenBinaryObjectIDs(ocspBinary)
		if utils.ContainsAny(keySet, objectIDs) {
			result = append(result, ocspBinary)
		}
	}
	return result
}

// filterTokensFromOcspMap ports the private filterTokensFromOcspMap(List).
func (s *PdfDssDictOCSPSource) filterTokensFromOcspMap(revocationTokens []spi.RevocationToken[revocation.OCSP]) []spi.RevocationToken[revocation.OCSP] {
	return s.filterTokensFromKeys(revocationTokens, pdfDssDictOCSPSourceSortedKeys(s.OcspMap()))
}

// filterTokensFromKeys ports the private filterTokensFromKeys(Collection, Collection).
func (s *PdfDssDictOCSPSource) filterTokensFromKeys(revocationTokens []spi.RevocationToken[revocation.OCSP],
	keySet []PdfObjectKey) []spi.RevocationToken[revocation.OCSP] {
	result := make([]spi.RevocationToken[revocation.OCSP], 0)
	for _, ocspToken := range revocationTokens {
		objectIDs := s.compositeOCSPSource.RevocationTokenIDs(ocspToken)
		if utils.ContainsAny(keySet, objectIDs) {
			result = append(result, ocspToken)
		}
	}
	return result
}

// AllRevocationBinariesWithOrigins returns all the OCSP binaries of this dictionary with their
// origins. Port of the getAllRevocationBinariesWithOrigins() override.
func (s *PdfDssDictOCSPSource) AllRevocationBinariesWithOrigins() []PdfDssDictOCSPSourceBinaryOriginsEntry {
	result := make([]PdfDssDictOCSPSourceBinaryOriginsEntry, 0)

	binaries := s.compositeOCSPSource.AllRevocationBinaries()
	filteredBinaries := s.filterBinariesFromKeys(binaries, pdfDssDictOCSPSourceSortedKeys(s.OcspMap()))
	for _, ocspBinary := range filteredBinaries {
		result = append(result, PdfDssDictOCSPSourceBinaryOriginsEntry{
			Binary:  ocspBinary,
			Origins: s.revocationDataOriginsForBinary(ocspBinary),
		})
	}
	return result
}

// revocationDataOriginsForBinary ports the private
// getRevocationDataOrigins(EncapsulatedRevocationTokenIdentifier<OCSP>).
//
// As upstream, dssDictionary is dereferenced without a nil check here.
func (s *PdfDssDictOCSPSource) revocationDataOriginsForBinary(
	ocspBinary spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]) []enumerations.RevocationOrigin {
	result := make([]enumerations.RevocationOrigin, 0)
	tokenBinaryObjectIDs := s.compositeOCSPSource.TokenBinaryObjectIDs(ocspBinary)
	if utils.ContainsAny(pdfDssDictOCSPSourceSortedKeys(s.dssDictionary.OCSPs()), tokenBinaryObjectIDs) {
		result = append(result, enumerations.RevocationOrigin_DSS_DICTIONARY)
	}
	for _, vriDict := range PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName) {
		if utils.ContainsAny(pdfDssDictOCSPSourceSortedKeys(vriDict.OCSPs()), tokenBinaryObjectIDs) {
			if !pdfDssDictOCSPSourceContainsOrigin(result, enumerations.RevocationOrigin_VRI_DICTIONARY) {
				result = append(result, enumerations.RevocationOrigin_VRI_DICTIONARY)
			}
		}
	}
	return result
}

// AllRevocationTokensWithOrigins returns all the OCSP tokens of this dictionary with their
// origins. Port of the getAllRevocationTokensWithOrigins() override.
func (s *PdfDssDictOCSPSource) AllRevocationTokensWithOrigins() []spi.RevocationTokenOriginsEntry[revocation.OCSP] {
	result := make([]spi.RevocationTokenOriginsEntry[revocation.OCSP], 0)

	tokens := s.compositeOCSPSource.AllRevocationTokens()
	filteredTokens := s.filterTokensFromKeys(tokens, pdfDssDictOCSPSourceSortedKeys(s.OcspMap()))
	for _, ocspToken := range filteredTokens {
		result = append(result, spi.RevocationTokenOriginsEntry[revocation.OCSP]{
			Token:   ocspToken,
			Origins: s.revocationDataOriginsForToken(ocspToken),
		})
	}
	return result
}

// revocationDataOriginsForToken ports the private getRevocationDataOrigins(RevocationToken<OCSP>).
//
// As upstream, dssDictionary is dereferenced without a nil check here.
func (s *PdfDssDictOCSPSource) revocationDataOriginsForToken(
	ocspToken spi.RevocationToken[revocation.OCSP]) []enumerations.RevocationOrigin {
	result := make([]enumerations.RevocationOrigin, 0)
	tokenObjectIDs := s.compositeOCSPSource.RevocationTokenIDs(ocspToken)
	if utils.ContainsAny(pdfDssDictOCSPSourceSortedKeys(s.dssDictionary.OCSPs()), tokenObjectIDs) {
		result = append(result, enumerations.RevocationOrigin_DSS_DICTIONARY)
	}
	for _, vriDict := range PAdESUtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName) {
		if utils.ContainsAny(pdfDssDictOCSPSourceSortedKeys(vriDict.OCSPs()), tokenObjectIDs) {
			if !pdfDssDictOCSPSourceContainsOrigin(result, enumerations.RevocationOrigin_VRI_DICTIONARY) {
				result = append(result, enumerations.RevocationOrigin_VRI_DICTIONARY)
			}
		}
	}
	return result
}

// pdfDssDictOCSPSourceContainsOrigin reports whether origins already contains origin, standing
// in for the Set<RevocationOrigin> semantics of the Java result.
func pdfDssDictOCSPSourceContainsOrigin(origins []enumerations.RevocationOrigin, origin enumerations.RevocationOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}

// pdfDssDictOCSPSourceContainsKey reports whether objectIDs already contains objectID.
func pdfDssDictOCSPSourceContainsKey(objectIDs []PdfObjectKey, objectID PdfObjectKey) bool {
	for _, id := range objectIDs {
		if id == objectID {
			return true
		}
	}
	return false
}

// pdfDssDictOCSPSourceSortedKeys returns the keys of a PDF-object-keyed map in ascending
// (object number, generation) order; see the determinism note in
// pdf_composite_dss_dict_certificate_source.go.
func pdfDssDictOCSPSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
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

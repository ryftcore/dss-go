// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfCompositeDssDictCRLSource.java (DSS 6.5.RC1).
//
// The forward dependencies this file shares with the rest of the chunk (PdfObjectKey,
// PdfDssDict, PdfVriDict, PAdESUtilsVRIsWithName) and the determinism rule applied to
// PDF-object-keyed map walks are documented in pdf_composite_dss_dict_certificate_source.go.
//
// VIRTUAL DISPATCH: Java's override of addRevocation(RevocationToken,
// EncapsulatedRevocationTokenIdentifier) is reached by virtual dispatch from
// OfflineCRLSource#getRevocationTokens, which is where every token this class caches in
// crlTokenMap is actually built. Go's embedding does not dispatch a base method call back into
// the embedder, so RevocationTokens is overridden here as well and re-applies the bookkeeping
// over the tokens spi.OfflineCRLSourceBase has just built (see RevocationTokens below). Without
// it RevocationTokenIDs would answer nil for every token and PdfDssDictCRLSource would filter
// all of them away.
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// pdfCompositeDssDictCRLSourceTokenEntry pairs a built CRL token with the PDF object ids of the
// binary it was built from, standing in for one entry of Java's
// Map<RevocationToken<CRL>, Set<PdfObjectKey>>: RevocationToken has no Go-usable map-key
// representation (its equals() also compares the related certificate), so the map is an
// insertion-ordered slice of pairs searched with RevocationToken.Equals, the convention
// spi.OfflineRevocationSourceBase already established for the same Java shape.
type pdfCompositeDssDictCRLSourceTokenEntry struct {
	token     spi.RevocationToken[revocation.CRL]
	objectIDs []PdfObjectKey
}

// PdfCompositeDssDictCRLSource represents a merged result of the CRL binaries extracted from
// the different /DSS revisions of a PDF document.
type PdfCompositeDssDictCRLSource struct {
	spi.OfflineCRLSourceBase

	// crlBinaryByIDMap is the composite map of CRL tokens extracted from different /DSS
	// revisions, by id. Java's Map<PdfObjectKey, Set<CRLBinary>>, whose Set value is kept here
	// as an insertion-ordered slice deduplicated on the binary's DSS identifier. NOTE: upstream
	// only ever writes this field; it is carried for parity.
	crlBinaryByIDMap map[PdfObjectKey][]*crlparser.CRLBinary

	// crlBinaryByObjectMap is the composite map of CRL tokens extracted from different /DSS
	// revisions, by encoded object binaries. Java keys it on
	// EncapsulatedRevocationTokenIdentifier<CRL>, whose equality is its digest; the Go port
	// keys it on the identifier's AsXmlID(), which is that same digest.
	crlBinaryByObjectMap map[string][]PdfObjectKey

	// crlTokenMap is the cached map of created CRLTokens and the corresponding PDF object ids.
	crlTokenMap []pdfCompositeDssDictCRLSourceTokenEntry
}

// NewPdfCompositeDssDictCRLSource instantiates an object with an empty map of CRL token
// objects. Port of the default constructor.
func NewPdfCompositeDssDictCRLSource() *PdfCompositeDssDictCRLSource {
	source := &PdfCompositeDssDictCRLSource{
		OfflineCRLSourceBase: spi.NewOfflineCRLSourceBase(),
		crlBinaryByIDMap:     make(map[PdfObjectKey][]*crlparser.CRLBinary),
		crlBinaryByObjectMap: make(map[string][]PdfObjectKey),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// PopulateFromDssDictionary adds the CRL tokens extracted from a /DSS revision.
// Port of populateFromDssDictionary(PdfDssDict).
func (s *PdfCompositeDssDictCRLSource) PopulateFromDssDictionary(dssDictionary PdfDssDict) {
	s.ExtractDSSCRLs(dssDictionary)
	s.ExtractVRICRLs(dssDictionary)
}

// ExtractDSSCRLs extracts the CRLs from the DSS dictionary.
// Port of the protected extractDSSCRLs(PdfDssDict).
//
// As upstream, a nil dssDictionary is not tolerated here (Java raises a NullPointerException on
// dssDictionary.getCRLs(), Go panics on the nil-interface call); the only caller,
// PdfCompositeDssDictionary, guards the nil case itself.
func (s *PdfCompositeDssDictCRLSource) ExtractDSSCRLs(dssDictionary PdfDssDict) {
	dssCrlMap := dssDictionary.CRLs()
	s.populateObjectsMap(dssCrlMap)
	for _, key := range pdfCompositeDssDictCRLSourceSortedKeys(dssCrlMap) {
		s.AddBinary(dssCrlMap[key], enumerations.RevocationOriginDSSDictionary)
	}
}

// ExtractVRICRLs extracts the CRLs from all the embedded VRI dictionaries.
// Port of the protected extractVRICRLs(PdfDssDict).
func (s *PdfCompositeDssDictCRLSource) ExtractVRICRLs(dssDictionary PdfDssDict) {
	if dssDictionary != nil {
		vriDictList := dssDictionary.VRIs()
		for _, vriDict := range vriDictList {
			s.populateObjectsMap(vriDict.CRLs())
			s.ExtractVRICRLsFromVRIDictionary(vriDict)
		}
	}
}

// ExtractVRICRLsFromVRIDictionary extracts the CRLs from the VRI dictionary.
// Port of the protected extractVRICRLs(PdfVriDict) overload, which Go's lack of overloading
// gives a distinct name from ExtractVRICRLs (*PdfVriDict also satisfies PdfDssDict, so the two
// Java signatures would otherwise collide).
func (s *PdfCompositeDssDictCRLSource) ExtractVRICRLsFromVRIDictionary(vriDictionary *PdfVriDict) {
	if vriDictionary != nil {
		crlMap := vriDictionary.CRLs()
		for _, key := range pdfCompositeDssDictCRLSourceSortedKeys(crlMap) {
			s.AddBinary(crlMap[key], enumerations.RevocationOriginVRIDictionary)
		}
	}
}

// populateObjectsMap ports the private populateObjectsMap(Map<PdfObjectKey, CRLBinary>).
func (s *PdfCompositeDssDictCRLSource) populateObjectsMap(crlMap map[PdfObjectKey]*crlparser.CRLBinary) {
	for _, key := range pdfCompositeDssDictCRLSourceSortedKeys(crlMap) {
		s.populateMapByID(key, crlMap[key])
		s.populateMapByObject(key, crlMap[key])
	}
}

// populateMapByID ports the private populateMapById(PdfObjectKey, CRLBinary).
func (s *PdfCompositeDssDictCRLSource) populateMapByID(objectID PdfObjectKey, crlBinary *crlparser.CRLBinary) {
	crlBinaries := s.crlBinaryByIDMap[objectID]
	found := false
	for _, binary := range crlBinaries {
		if binary.AsXmlID() == crlBinary.AsXmlID() {
			found = true
			break
		}
	}
	if !found {
		crlBinaries = append(crlBinaries, crlBinary)
	}
	s.crlBinaryByIDMap[objectID] = crlBinaries
}

// populateMapByObject ports the private populateMapByObject(PdfObjectKey, CRLBinary).
func (s *PdfCompositeDssDictCRLSource) populateMapByObject(objectID PdfObjectKey, crlBinary *crlparser.CRLBinary) {
	objectIDs := s.crlBinaryByObjectMap[crlBinary.AsXmlID()]
	if !pdfCompositeDssDictCRLSourceContainsKey(objectIDs, objectID) {
		objectIDs = append(objectIDs, objectID)
	}
	s.crlBinaryByObjectMap[crlBinary.AsXmlID()] = objectIDs
}

// RevocationTokenIDs returns the PDF object identifiers of the extracted revocation token, nil
// when the token is unknown (Java returns null).
// Port of the protected getRevocationTokenIds(RevocationToken<CRL>).
func (s *PdfCompositeDssDictCRLSource) RevocationTokenIDs(crlToken spi.RevocationToken[revocation.CRL]) []PdfObjectKey {
	for _, entry := range s.crlTokenMap {
		if entry.token.Equals(crlToken) {
			return entry.objectIDs
		}
	}
	return nil
}

// AddRevocationWithBinary adds a RevocationToken built from binary and caches the PDF object
// ids the binary was found at. Port of the addRevocation(RevocationToken,
// EncapsulatedRevocationTokenIdentifier) override.
func (s *PdfCompositeDssDictCRLSource) AddRevocationWithBinary(token spi.RevocationToken[revocation.CRL],
	binary spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]) {
	s.OfflineCRLSourceBase.AddRevocationWithBinary(token, binary)

	tokenBinaryObjectIDs := s.TokenBinaryObjectIDs(binary)
	s.putRevocationTokenIDs(token, tokenBinaryObjectIDs)
}

// TokenBinaryObjectIDs returns the PDF object identifiers of the provided binary, nil when the
// binary is unknown (Java returns null).
// Port of the protected getTokenBinaryObjectIds(EncapsulatedRevocationTokenIdentifier<CRL>).
func (s *PdfCompositeDssDictCRLSource) TokenBinaryObjectIDs(binary spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL]) []PdfObjectKey {
	return s.crlBinaryByObjectMap[binary.AsXmlID()]
}

// RevocationTokens builds a CRL token per contained CRL valid for the certificate's issuer,
// exactly as spi.OfflineCRLSourceBase does, and then records the PDF object ids of the binary
// each built token comes from.
//
// This second half is what Java gets for free: OfflineCRLSource#getRevocationTokens calls
// addRevocation(RevocationToken, EncapsulatedRevocationTokenIdentifier), which virtual dispatch
// routes to this class's override. Go's embedded base calls its own AddRevocationWithBinary
// instead, so the override is re-applied here over the tokens it returned. The binary a CRL
// token was built from is recovered from the token itself (CRLToken#getCrlValidity().getCrlBinary()
// is the very binary spi.OfflineCRLSourceBase passed to AddRevocationWithBinary).
func (s *PdfCompositeDssDictCRLSource) RevocationTokens(certificateToken *model.CertificateToken,
	issuerToken *model.CertificateToken) ([]spi.RevocationToken[revocation.CRL], error) {
	revocationTokens, err := s.OfflineCRLSourceBase.RevocationTokens(certificateToken, issuerToken)
	if err != nil {
		return nil, err
	}
	for _, revocationToken := range revocationTokens {
		crlToken, ok := revocationToken.(*spi.CRLToken)
		if !ok {
			continue
		}
		crlValidity := crlToken.CrlValidity()
		if crlValidity == nil {
			continue
		}
		s.putRevocationTokenIDs(revocationToken, s.TokenBinaryObjectIDs(crlValidity.CrlBinary()))
	}
	return revocationTokens, nil
}

// putRevocationTokenIDs is the Map#put half of the addRevocation override: it replaces the
// entry of token, adding it when absent. Java stores a null value for an unknown binary too,
// which the nil objectIDs slice reproduces.
func (s *PdfCompositeDssDictCRLSource) putRevocationTokenIDs(token spi.RevocationToken[revocation.CRL], objectIDs []PdfObjectKey) {
	for i := range s.crlTokenMap {
		if s.crlTokenMap[i].token.Equals(token) {
			s.crlTokenMap[i].objectIDs = objectIDs
			return
		}
	}
	s.crlTokenMap = append(s.crlTokenMap, pdfCompositeDssDictCRLSourceTokenEntry{token: token, objectIDs: objectIDs})
}

// pdfCompositeDssDictCRLSourceContainsKey reports whether objectIDs already contains objectID.
func pdfCompositeDssDictCRLSourceContainsKey(objectIDs []PdfObjectKey, objectID PdfObjectKey) bool {
	for _, id := range objectIDs {
		if id == objectID {
			return true
		}
	}
	return false
}

// pdfCompositeDssDictCRLSourceSortedKeys returns the keys of a PDF-object-keyed map in ascending
// (object number, generation) order; see the determinism note in
// pdf_composite_dss_dict_certificate_source.go.
func pdfCompositeDssDictCRLSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
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

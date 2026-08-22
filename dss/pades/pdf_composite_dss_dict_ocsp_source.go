// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfCompositeDssDictOCSPSource.java (DSS 6.5.RC1).
//
// The forward dependencies this file shares with the rest of the chunk (PdfObjectKey,
// PdfDssDict, PdfVriDict, PAdESUtilsVRIsWithName) and the determinism rule applied to
// PDF-object-keyed map walks are documented in pdf_composite_dss_dict_certificate_source.go.
//
// VIRTUAL DISPATCH: as in PdfCompositeDssDictCRLSource, Java's override of
// addRevocation(RevocationToken, EncapsulatedRevocationTokenIdentifier) is reached by virtual
// dispatch from OfflineOCSPSource#getRevocationTokens, which is where every token cached in
// ocspTokenMap is built. Go's embedding cannot dispatch that call back into the embedder, so
// RevocationTokens is overridden here too and re-applies the bookkeeping (see RevocationTokens).
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// pdfCompositeDssDictOCSPSourceTokenEntry pairs a built OCSP token with the PDF object ids of
// the binary it was built from, standing in for one entry of Java's
// Map<RevocationToken<OCSP>, Set<PdfObjectKey>>; see the same type in
// pdf_composite_dss_dict_crl_source.go for why the map is a slice of pairs.
type pdfCompositeDssDictOCSPSourceTokenEntry struct {
	token     spi.RevocationToken[revocation.OCSP]
	objectIDs []PdfObjectKey
}

// PdfCompositeDssDictOCSPSource represents a merged result of the OCSP binaries extracted from
// the different /DSS revisions of a PDF document.
type PdfCompositeDssDictOCSPSource struct {
	spi.OfflineOCSPSourceBase

	// ocspBinaryByIDMap is the composite map of OCSP tokens extracted from different /DSS
	// revisions, by object id. Java's Map<PdfObjectKey, Set<OCSPResponseBinary>>, whose Set
	// value is kept here as an insertion-ordered slice deduplicated on the binary's DSS
	// identifier. NOTE: upstream only ever writes this field; it is carried for parity.
	ocspBinaryByIDMap map[PdfObjectKey][]*spi.OCSPResponseBinary

	// ocspBinaryByObjectMap is the composite map of OCSP tokens extracted from different /DSS
	// revisions, by encoded object binaries, keyed on the identifier's AsXmlID() - the digest
	// Java's EncapsulatedRevocationTokenIdentifier<OCSP> key compares on.
	ocspBinaryByObjectMap map[string][]PdfObjectKey

	// ocspTokenMap is the cached map of created OCSPTokens and the corresponding PDF object ids.
	ocspTokenMap []pdfCompositeDssDictOCSPSourceTokenEntry
}

// NewPdfCompositeDssDictOCSPSource instantiates an object with an empty map of OCSP token
// objects. Port of the default constructor.
func NewPdfCompositeDssDictOCSPSource() *PdfCompositeDssDictOCSPSource {
	source := &PdfCompositeDssDictOCSPSource{
		OfflineOCSPSourceBase: spi.NewOfflineOCSPSourceBase(),
		ocspBinaryByIDMap:     make(map[PdfObjectKey][]*spi.OCSPResponseBinary),
		ocspBinaryByObjectMap: make(map[string][]PdfObjectKey),
	}
	// The outermost concrete source registers itself, so that the RevocationToken dispatch of
	// OfflineRevocationSourceBase reaches this type's RevocationTokens override.
	source.InitOfflineRevocationSource(source)
	return source
}

// PopulateFromDssDictionary adds the OCSP tokens extracted from a /DSS revision.
// Port of populateFromDssDictionary(PdfDssDict).
func (s *PdfCompositeDssDictOCSPSource) PopulateFromDssDictionary(dssDictionary PdfDssDict) {
	s.ExtractDSSOCSPs(dssDictionary)
	s.ExtractVRIOCSPs(dssDictionary)
}

// ExtractDSSOCSPs extracts the OCSPs from the DSS dictionary.
// Port of the protected extractDSSOCSPs(PdfDssDict).
//
// As upstream, a nil dssDictionary is not tolerated here (Java raises a NullPointerException on
// dssDictionary.getOCSPs(), Go panics on the nil-interface call); the only caller,
// PdfCompositeDssDictionary, guards the nil case itself.
func (s *PdfCompositeDssDictOCSPSource) ExtractDSSOCSPs(dssDictionary PdfDssDict) {
	dssOCSPMap := dssDictionary.OCSPs()
	s.populateObjectsMap(dssOCSPMap)
	for _, key := range pdfCompositeDssDictOCSPSourceSortedKeys(dssOCSPMap) {
		s.AddBinary(dssOCSPMap[key], enumerations.RevocationOriginDSSDictionary)
	}
}

// ExtractVRIOCSPs extracts the OCSPs from all the embedded VRI dictionaries.
// Port of the protected extractVRIOCSPs(PdfDssDict).
func (s *PdfCompositeDssDictOCSPSource) ExtractVRIOCSPs(dssDictionary PdfDssDict) {
	if dssDictionary != nil {
		vriDictList := dssDictionary.VRIs()
		for _, vriDict := range vriDictList {
			s.populateObjectsMap(vriDict.OCSPs())
			s.ExtractVRIOCSPsFromVRIDictionary(vriDict)
		}
	}
}

// ExtractVRIOCSPsFromVRIDictionary extracts the OCSPs from the VRI dictionary.
// Port of the protected extractVRIOCSPs(PdfVriDict) overload, which Go's lack of overloading
// gives a distinct name from ExtractVRIOCSPs (*PdfVriDict also satisfies PdfDssDict, so the two
// Java signatures would otherwise collide).
func (s *PdfCompositeDssDictOCSPSource) ExtractVRIOCSPsFromVRIDictionary(vriDictionary *PdfVriDict) {
	if vriDictionary != nil {
		ocspMap := vriDictionary.OCSPs()
		for _, key := range pdfCompositeDssDictOCSPSourceSortedKeys(ocspMap) {
			s.AddBinary(ocspMap[key], enumerations.RevocationOriginVRIDictionary)
		}
	}
}

// populateObjectsMap ports the private populateObjectsMap(Map<PdfObjectKey, OCSPResponseBinary>).
func (s *PdfCompositeDssDictOCSPSource) populateObjectsMap(ocspMap map[PdfObjectKey]*spi.OCSPResponseBinary) {
	for _, key := range pdfCompositeDssDictOCSPSourceSortedKeys(ocspMap) {
		s.populateMapByID(key, ocspMap[key])
		s.populateMapByObject(key, ocspMap[key])
	}
}

// populateMapByID ports the private populateMapById(PdfObjectKey, OCSPResponseBinary).
func (s *PdfCompositeDssDictOCSPSource) populateMapByID(objectID PdfObjectKey, ocspBinary *spi.OCSPResponseBinary) {
	ocspBinaries := s.ocspBinaryByIDMap[objectID]
	found := false
	for _, binary := range ocspBinaries {
		if binary.AsXmlID() == ocspBinary.AsXmlID() {
			found = true
			break
		}
	}
	if !found {
		ocspBinaries = append(ocspBinaries, ocspBinary)
	}
	s.ocspBinaryByIDMap[objectID] = ocspBinaries
}

// populateMapByObject ports the private populateMapByObject(PdfObjectKey, OCSPResponseBinary).
func (s *PdfCompositeDssDictOCSPSource) populateMapByObject(objectID PdfObjectKey, ocspBinary *spi.OCSPResponseBinary) {
	objectIDs := s.ocspBinaryByObjectMap[ocspBinary.AsXmlID()]
	if !pdfCompositeDssDictOCSPSourceContainsKey(objectIDs, objectID) {
		objectIDs = append(objectIDs, objectID)
	}
	s.ocspBinaryByObjectMap[ocspBinary.AsXmlID()] = objectIDs
}

// RevocationTokenIDs returns the PDF object identifiers of the extracted revocation token, nil
// when the token is unknown (Java returns null).
// Port of the protected getRevocationTokenIds(RevocationToken<OCSP>).
func (s *PdfCompositeDssDictOCSPSource) RevocationTokenIDs(ocspToken spi.RevocationToken[revocation.OCSP]) []PdfObjectKey {
	for _, entry := range s.ocspTokenMap {
		if entry.token.Equals(ocspToken) {
			return entry.objectIDs
		}
	}
	return nil
}

// AddRevocationWithBinary adds a RevocationToken built from binary and caches the PDF object
// ids the binary was found at. Port of the addRevocation(RevocationToken,
// EncapsulatedRevocationTokenIdentifier) override.
func (s *PdfCompositeDssDictOCSPSource) AddRevocationWithBinary(token spi.RevocationToken[revocation.OCSP],
	binary spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]) {
	s.OfflineOCSPSourceBase.AddRevocationWithBinary(token, binary)

	tokenBinaryObjectIDs := s.TokenBinaryObjectIDs(binary)
	s.putRevocationTokenIDs(token, tokenBinaryObjectIDs)
}

// TokenBinaryObjectIDs returns the PDF object identifiers of the provided binary, nil when the
// binary is unknown (Java returns null).
// Port of the protected getTokenBinaryObjectIds(EncapsulatedRevocationTokenIdentifier<OCSP>).
func (s *PdfCompositeDssDictOCSPSource) TokenBinaryObjectIDs(binary spi.EncapsulatedRevocationTokenIdentifier[revocation.OCSP]) []PdfObjectKey {
	return s.ocspBinaryByObjectMap[binary.AsXmlID()]
}

// RevocationTokens returns the OCSP tokens concerning the certificate exactly as
// spi.OfflineOCSPSourceBase does, and then records the PDF object ids of the binary each token
// comes from - the bookkeeping Java's addRevocation(RevocationToken,
// EncapsulatedRevocationTokenIdentifier) override performs when OfflineOCSPSource builds them;
// see the file header.
//
// The binary a token was built from is recovered by identity of the wrapped basic OCSP
// response: spi.OfflineOCSPSourceBase hands OCSPResponseBinary#BasicOCSPResp() straight to
// NewOCSPToken, which keeps that very pointer.
func (s *PdfCompositeDssDictOCSPSource) RevocationTokens(certificate *model.CertificateToken,
	issuer *model.CertificateToken) ([]spi.RevocationToken[revocation.OCSP], error) {
	revocationTokens, err := s.OfflineOCSPSourceBase.RevocationTokens(certificate, issuer)
	if err != nil {
		return nil, err
	}
	for _, revocationToken := range revocationTokens {
		ocspToken, ok := revocationToken.(*spi.OCSPToken)
		if !ok {
			continue
		}
		for _, binary := range s.AllRevocationBinaries() {
			ocspBinary, ok := binary.(*spi.OCSPResponseBinary)
			if ok && ocspBinary.BasicOCSPResp() == ocspToken.BasicOCSPResp() {
				s.putRevocationTokenIDs(revocationToken, s.TokenBinaryObjectIDs(ocspBinary))
				break
			}
		}
	}
	return revocationTokens, nil
}

// putRevocationTokenIDs is the Map#put half of the addRevocation override: it replaces the
// entry of token, adding it when absent. Java stores a null value for an unknown binary too,
// which the nil objectIDs slice reproduces.
func (s *PdfCompositeDssDictOCSPSource) putRevocationTokenIDs(token spi.RevocationToken[revocation.OCSP], objectIDs []PdfObjectKey) {
	for i := range s.ocspTokenMap {
		if s.ocspTokenMap[i].token.Equals(token) {
			s.ocspTokenMap[i].objectIDs = objectIDs
			return
		}
	}
	s.ocspTokenMap = append(s.ocspTokenMap, pdfCompositeDssDictOCSPSourceTokenEntry{token: token, objectIDs: objectIDs})
}

// pdfCompositeDssDictOCSPSourceContainsKey reports whether objectIDs already contains objectID.
func pdfCompositeDssDictOCSPSourceContainsKey(objectIDs []PdfObjectKey, objectID PdfObjectKey) bool {
	for _, id := range objectIDs {
		if id == objectID {
			return true
		}
	}
	return false
}

// pdfCompositeDssDictOCSPSourceSortedKeys returns the keys of a PDF-object-keyed map in
// ascending (object number, generation) order; see the determinism note in
// pdf_composite_dss_dict_certificate_source.go.
func pdfCompositeDssDictOCSPSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
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

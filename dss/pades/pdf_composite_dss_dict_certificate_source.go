// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfCompositeDssDictCertificateSource.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES (not in this chunk's manifest, ported by sibling chunks into this same
// flattened pades package). Every file of this chunk was written against these shapes:
//
//   - PdfObjectKey (eu.europa.esig.dss.pades.validation.PdfObjectKey) - an interface with
//     Value() any, Number() int64, Generation() int. Java uses it as a HashMap key, so its
//     implementations carry equals/hashCode over the (object number, generation) pair; the Go
//     port therefore uses it directly as a Go map key, which requires the implementing type to
//     be comparable (it is: a PDF object key is a number plus a generation).
//   - PdfDssDict (eu.europa.esig.dss.pdf.PdfDssDict) - an interface with
//     CRLs() map[PdfObjectKey]*crlparser.CRLBinary, OCSPs() map[PdfObjectKey]*spi.OCSPResponseBinary,
//     CERTs() map[PdfObjectKey]*model.CertificateToken and VRIs() []*PdfVriDict.
//   - PdfVriDict (eu.europa.esig.dss.pdf.PdfVriDict) - a concrete struct implementing PdfDssDict
//     and additionally exposing Name() string, TUTime() *time.Time and TSStream() []byte.
//   - PAdESUtilsVRIsWithName(pdfDssDict PdfDssDict, vriName string) []*PdfVriDict - the
//     flattened static PAdESUtils.getVRIsWithName(PdfDssDict, String); Java's null vriName
//     ("every VRI dictionary") is the empty string here, which is unambiguous since a VRI name
//     is the base-16 SHA-1 of a signature and never empty.
//
// DETERMINISM: upstream walks java.util.HashMaps keyed by PdfObjectKey whenever it materialises
// a List of tokens; Go's map iteration order is randomised rather than merely unspecified, so
// every such walk in this chunk is performed in ascending (object number, generation) order.
// This is the determinism normalisation the port applies elsewhere for the same reason, and it
// only fixes an order Java leaves arbitrary.
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PdfCompositeDssDictCertificateSource represents a merged result of extracted certificate
// sources from the /DSS revisions of a PDF document.
type PdfCompositeDssDictCertificateSource struct {
	spi.TokenCertificateSource

	// certMap is the composite map of certificate tokens extracted from different /DSS
	// revisions. Java's Map<PdfObjectKey, Set<CertificateToken>>; the Set<CertificateToken>
	// value keeps this port's convention of a map keyed by DSSIDAsString().
	certMap map[PdfObjectKey]map[string]*model.CertificateToken
}

// NewPdfCompositeDssDictCertificateSource instantiates an object with an empty map of
// certificate token objects. Port of the default constructor.
func NewPdfCompositeDssDictCertificateSource() *PdfCompositeDssDictCertificateSource {
	return &PdfCompositeDssDictCertificateSource{
		TokenCertificateSource: spi.NewTokenCertificateSource(),
		certMap:                make(map[PdfObjectKey]map[string]*model.CertificateToken),
	}
}

// PopulateFromDssDictionary adds the certificates extracted from a /DSS revision.
// Port of populateFromDssDictionary(PdfDssDict).
func (s *PdfCompositeDssDictCertificateSource) PopulateFromDssDictionary(dssDictionary PdfDssDict) {
	for _, certToken := range s.dssDictionaryCertValues(dssDictionary) {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginDSSDictionary)
	}
	for _, certToken := range s.vriDictionaryCertValues(dssDictionary) {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginVRIDictionary)
	}
}

// dssDictionaryCertValues gets the list of the DSS dictionary certificate tokens.
// Port of the private getDSSDictionaryCertValues(PdfDssDict).
func (s *PdfCompositeDssDictCertificateSource) dssDictionaryCertValues(dssDictionary PdfDssDict) []*model.CertificateToken {
	if dssDictionary != nil {
		dssCerts := dssDictionary.CERTs()
		s.populateObjectsMap(dssCerts)
		result := make([]*model.CertificateToken, 0, len(dssCerts))
		for _, key := range pdfCompositeDssDictCertificateSourceSortedKeys(dssCerts) {
			result = append(result, dssCerts[key])
		}
		return result
	}
	return []*model.CertificateToken{}
}

// vriDictionaryCertValues gets the list of the certificate tokens extracted from all VRI
// dictionaries. Port of the private getVRIDictionaryCertValues(PdfDssDict).
func (s *PdfCompositeDssDictCertificateSource) vriDictionaryCertValues(dssDictionary PdfDssDict) []*model.CertificateToken {
	if dssDictionary != nil {
		vriCerts := make(map[PdfObjectKey]*model.CertificateToken)
		vris := dssDictionary.VRIs()
		if vris != nil {
			for _, vri := range vris {
				for key, certToken := range vri.CERTs() {
					vriCerts[key] = certToken
				}
			}
		}
		s.populateObjectsMap(vriCerts)
		result := make([]*model.CertificateToken, 0, len(vriCerts))
		for _, key := range pdfCompositeDssDictCertificateSourceSortedKeys(vriCerts) {
			result = append(result, vriCerts[key])
		}
		return result
	}
	return []*model.CertificateToken{}
}

// populateObjectsMap ports the private populateObjectsMap(Map<PdfObjectKey, CertificateToken>).
func (s *PdfCompositeDssDictCertificateSource) populateObjectsMap(certificateTokenMap map[PdfObjectKey]*model.CertificateToken) {
	for _, key := range pdfCompositeDssDictCertificateSourceSortedKeys(certificateTokenMap) {
		certificateTokens := s.certMap[key]
		if certificateTokens == nil {
			certificateTokens = make(map[string]*model.CertificateToken)
		}
		certificateToken := certificateTokenMap[key]
		certificateTokens[certificateToken.DSSIDAsString()] = certificateToken
		s.certMap[key] = certificateTokens
	}
}

// CertificateTokensByObjectID returns the set of CertificateTokens carrying the given PDF
// object id, nil when the id is unknown (Java returns null).
// Port of the protected getCertificateTokensByObjectId(PdfObjectKey).
func (s *PdfCompositeDssDictCertificateSource) CertificateTokensByObjectID(objectID PdfObjectKey) map[string]*model.CertificateToken {
	return s.certMap[objectID]
}

// pdfCompositeDssDictCertificateSourceSortedKeys returns the keys of a PDF-object-keyed map in
// ascending (object number, generation) order; see the determinism note in the file header.
func pdfCompositeDssDictCertificateSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
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

// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfDssDictCertificateSource.java (DSS 6.5.RC1).
package pades

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PdfDssDictCertificateSource is the certificate source extracted from a DSS dictionary.
type PdfDssDictCertificateSource struct {
	spi.TokenCertificateSource

	// compositeCertificateSource is the merged certificate source combined from all /DSS
	// revisions.
	compositeCertificateSource *PdfCompositeDssDictCertificateSource

	// dssDictionary is the DSS dictionary.
	dssDictionary PdfDssDict

	// relatedVRIDictionaryName is the name of the signature's VRI dictionary, when applicable;
	// the empty string stands for Java's null, i.e. "every VRI dictionary".
	relatedVRIDictionaryName string
}

// NewPdfDssDictCertificateSource is the port of the
// PdfDssDictCertificateSource(PdfCompositeDssDictCertificateSource, PdfDssDict) constructor.
func NewPdfDssDictCertificateSource(compositeCertificateSource *PdfCompositeDssDictCertificateSource,
	dssDictionary PdfDssDict) *PdfDssDictCertificateSource {
	return NewPdfDssDictCertificateSourceWithVRIName(compositeCertificateSource, dssDictionary, "")
}

// NewPdfDssDictCertificateSourceWithVRIName is the port of the
// PdfDssDictCertificateSource(PdfCompositeDssDictCertificateSource, PdfDssDict, String)
// constructor, to be used for a signature. vriDictionaryName is the SHA-1 of the signature
// name; the empty string stands for Java's null.
func NewPdfDssDictCertificateSourceWithVRIName(compositeCertificateSource *PdfCompositeDssDictCertificateSource,
	dssDictionary PdfDssDict, vriDictionaryName string) *PdfDssDictCertificateSource {
	source := &PdfDssDictCertificateSource{
		TokenCertificateSource:     spi.NewTokenCertificateSource(),
		compositeCertificateSource: compositeCertificateSource,
		dssDictionary:              dssDictionary,
		relatedVRIDictionaryName:   vriDictionaryName,
	}
	source.extractFromDssDictSource()
	return source
}

// extractFromDssDictSource ports the private extractFromDssDictSource().
func (s *PdfDssDictCertificateSource) extractFromDssDictSource() {
	for _, certToken := range s.DSSDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginDSSDictionary)
	}
	for _, certToken := range s.VRIDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginVRIDictionary)
	}
}

// CertificateMap gets a map of PDF object ids and the corresponding certificate tokens.
// Port of getCertificateMap().
//
// UPSTREAM BEHAVIOUR PRESERVED: Java puts the VRI entries into the very map
// PdfDssDict#getCERTs() returned - AbstractPdfDssDict hands out its internal field, so the
// dictionary itself gains the VRI certificates as a side effect of this call. The Go port keeps
// the aliasing rather than defensively copying, since the returned map is the same object the
// sibling PdfDssDict port hands out.
func (s *PdfDssDictCertificateSource) CertificateMap() map[PdfObjectKey]*model.CertificateToken {
	if s.dssDictionary != nil {
		dssCerts := s.dssDictionary.CERTs()
		vriDicts := UtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
		for _, vriDict := range vriDicts {
			for key, certToken := range vriDict.CERTs() {
				dssCerts[key] = certToken
			}
		}
		return dssCerts
	}
	return map[PdfObjectKey]*model.CertificateToken{}
}

// DSSDictionaryCertValues gets the list of the DSS dictionary certificate tokens.
// Port of getDSSDictionaryCertValues().
func (s *PdfDssDictCertificateSource) DSSDictionaryCertValues() []*model.CertificateToken {
	if s.dssDictionary != nil {
		return s.certificatesByKeys(pdfDssDictCertificateSourceSortedKeys(s.dssDictionary.CERTs()))
	}
	return []*model.CertificateToken{}
}

// VRIDictionaryCertValues gets the list of the certificate tokens extracted from all the VRI
// dictionaries. Port of getVRIDictionaryCertValues().
func (s *PdfDssDictCertificateSource) VRIDictionaryCertValues() []*model.CertificateToken {
	if s.dssDictionary != nil {
		certKeys := make([]PdfObjectKey, 0)
		vris := UtilsVRIsWithName(s.dssDictionary, s.relatedVRIDictionaryName)
		for _, vri := range vris {
			for _, key := range pdfDssDictCertificateSourceSortedKeys(vri.CERTs()) {
				if !pdfDssDictCertificateSourceContainsKey(certKeys, key) {
					certKeys = append(certKeys, key)
				}
			}
		}
		return s.certificatesByKeys(certKeys)
	}
	return []*model.CertificateToken{}
}

// certificatesByKeys ports the private getCertificatesByKeys(Collection<PdfObjectKey>).
//
// DEVIATION: Java calls addAll on the Set the composite source answers, which throws a
// NullPointerException for an object id the composite source has never seen; ranging over the
// nil map simply contributes nothing here. The tokens of one object id are appended in
// ascending DSS identifier order, since Java's HashSet leaves that order arbitrary.
func (s *PdfDssDictCertificateSource) certificatesByKeys(objectIDs []PdfObjectKey) []*model.CertificateToken {
	certificateTokens := make([]*model.CertificateToken, 0)
	for _, objectID := range objectIDs {
		tokens := s.compositeCertificateSource.CertificateTokensByObjectID(objectID)
		ids := make([]string, 0, len(tokens))
		for id := range tokens {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			certificateTokens = append(certificateTokens, tokens[id])
		}
	}
	return certificateTokens
}

// pdfDssDictCertificateSourceContainsKey reports whether objectIDs already contains objectID.
func pdfDssDictCertificateSourceContainsKey(objectIDs []PdfObjectKey, objectID PdfObjectKey) bool {
	for _, id := range objectIDs {
		if id == objectID {
			return true
		}
	}
	return false
}

// pdfDssDictCertificateSourceSortedKeys returns the keys of a PDF-object-keyed map in ascending
// (object number, generation) order; see the determinism note in
// pdf_composite_dss_dict_certificate_source.go.
func pdfDssDictCertificateSourceSortedKeys[V any](m map[PdfObjectKey]V) []PdfObjectKey {
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

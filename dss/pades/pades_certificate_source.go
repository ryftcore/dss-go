// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESCertificateSource.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart).
//
// cades.CertificateSource's constructor takes the CMS and the SignerInformation
// directly (see cades/cades_certificate_source.go), matching
// PdfSignatureRevision.getCMS()/PdfSignatureRevision itself not carrying a SignerInformation
// of its own - the signerInformation constructor parameter is threaded through unchanged.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CertificateSource is a CertificateSource that retrieves the certificate from a PAdES
// Signature.
type CertificateSource struct {
	*cades.CertificateSource

	// dssDictionaryCertificateSource is the certificate source of the DSS dictionary.
	dssDictionaryCertificateSource *PdfDssDictCertificateSource
}

// NewPAdESCertificateSource is the default constructor for CertificateSource.
// Port of the constructor PAdESCertificateSource(PdfSignatureRevision, String,
// SignerInformation).
//
// Panics with the Java message when vriDictionaryName is empty (Objects.requireNonNull; the
// empty string means no VRI-name filter, see pdf_dss_dict_certificate_source.go).
func NewPAdESCertificateSource(pdfSignatureRevision *PdfSignatureRevision, vriDictionaryName string,
	signerInformation *cmscore.SignerInfo) (*CertificateSource, error) {
	if vriDictionaryName == "" {
		panic("vriDictionaryName cannot be null!")
	}

	base, err := cades.NewCAdESCertificateSource(pdfSignatureRevision.CMS(), signerInformation)
	if err != nil {
		return nil, err
	}

	s := &CertificateSource{
		CertificateSource: base,
		dssDictionaryCertificateSource: NewPdfDssDictCertificateSourceWithVRIName(
			pdfSignatureRevision.CompositeDssDictionary().CertificateSource(),
			pdfSignatureRevision.DssDictionary(), vriDictionaryName),
	}
	s.extractFromDssDictSource()
	return s, nil
}

// extractFromDssDictSource ports the private extractFromDssDictSource().
func (s *CertificateSource) extractFromDssDictSource() {
	for _, certToken := range s.DSSDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginDSSDictionary)
	}
	for _, certToken := range s.VRIDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOriginVRIDictionary)
	}
}

// CertificateMap gets the map of certificate PDF object ids and the certificateTokens.
// Port of getCertificateMap().
func (s *CertificateSource) CertificateMap() map[PdfObjectKey]*model.CertificateToken {
	return s.dssDictionaryCertificateSource.CertificateMap()
}

// CertificateValues is not applicable for PAdES. Port of the getCertificateValues() override.
func (s *CertificateSource) CertificateValues() []*model.CertificateToken {
	return []*model.CertificateToken{}
}

// CompleteCertificateRefs is not applicable for PAdES. Port of the getCompleteCertificateRefs()
// override.
func (s *CertificateSource) CompleteCertificateRefs() []*spi.CertificateRef {
	return []*spi.CertificateRef{}
}

// AttributeCertificateRefs is not applicable for PAdES. Port of the
// getAttributeCertificateRefs() override.
func (s *CertificateSource) AttributeCertificateRefs() []*spi.CertificateRef {
	return []*spi.CertificateRef{}
}

// DSSDictionaryCertValues gets the list of the DSS dictionary certificate tokens.
// Port of the getDSSDictionaryCertValues() override. Shadows the embedded base's method of the
// same name (spi.SignatureCertificateSource.DSSDictionaryCertValues): see PORTING.md's "Virtual
// dispatch" precedent (analyzer/default_document_analyzer.go) - callers reaching this behaviour
// through the generic *spi.SignatureCertificateSource the AdvancedSignature interface's
// CertificateSource() method returns get the base's origin-tagged-map lookup instead, which
// extractFromDssDictSource above populates with the very same tokens this override answers, so
// the two stay in agreement for every caller either way.
func (s *CertificateSource) DSSDictionaryCertValues() []*model.CertificateToken {
	return s.dssDictionaryCertificateSource.DSSDictionaryCertValues()
}

// VRIDictionaryCertValues gets the list of the certificate tokens extracted from the VRI
// dictionary. Port of the getVRIDictionaryCertValues() override. Shadows the embedded base's
// method of the same name; see DSSDictionaryCertValues's doc comment.
func (s *CertificateSource) VRIDictionaryCertValues() []*model.CertificateToken {
	return s.dssDictionaryCertificateSource.VRIDictionaryCertValues()
}

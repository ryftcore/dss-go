// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESCertificateSource.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart).
//
// FORWARD DEPENDENCIES (not in this chunk's manifest, ported by sibling chunks):
//
//   - PdfSignatureRevision (eu.europa.esig.dss.pdf.PdfSignatureRevision), which flattens into
//     this same package per the phase 5b layout. Every file of this chunk that needs it was
//     written against this assumed shape:
//
//     type PdfSignatureRevision struct { /* embeds PdfRevision */ }
//     func (r *PdfSignatureRevision) CMS() *cms.CMS
//     func (r *PdfSignatureRevision) CompositeDssDictionary() *PdfCompositeDssDictionary
//     func (r *PdfSignatureRevision) DssDictionary() PdfDssDict
//     func (r *PdfSignatureRevision) Fields() []*PdfSignatureField
//
//   - cades.CAdESCertificateSource's constructor takes the CMS and the SignerInformation
//     directly (see cades/cades_certificate_source.go), matching
//     PdfSignatureRevision.getCMS()/PdfSignatureRevision itself not carrying a SignerInformation
//     of its own - the signerInformation constructor parameter is threaded through unchanged.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PAdESCertificateSource is a CertificateSource that retrieves the certificate from a PAdES
// Signature.
type PAdESCertificateSource struct {
	*cades.CAdESCertificateSource

	// dssDictionaryCertificateSource is the certificate source of the DSS dictionary.
	dssDictionaryCertificateSource *PdfDssDictCertificateSource
}

// NewPAdESCertificateSource is the default constructor for PAdESCertificateSource.
// Port of the constructor PAdESCertificateSource(PdfSignatureRevision, String,
// SignerInformation).
//
// Panics with the Java message when vriDictionaryName is empty (Objects.requireNonNull; the
// empty string stands for Java's null throughout this port, see
// pdf_dss_dict_certificate_source.go).
func NewPAdESCertificateSource(pdfSignatureRevision *PdfSignatureRevision, vriDictionaryName string,
	signerInformation *cmscore.SignerInfo) (*PAdESCertificateSource, error) {
	if vriDictionaryName == "" {
		panic("vriDictionaryName cannot be null!")
	}

	base, err := cades.NewCAdESCertificateSource(pdfSignatureRevision.CMS(), signerInformation)
	if err != nil {
		return nil, err
	}

	s := &PAdESCertificateSource{
		CAdESCertificateSource: base,
		dssDictionaryCertificateSource: NewPdfDssDictCertificateSourceWithVRIName(
			pdfSignatureRevision.CompositeDssDictionary().CertificateSource(),
			pdfSignatureRevision.DssDictionary(), vriDictionaryName),
	}
	s.extractFromDssDictSource()
	return s, nil
}

// extractFromDssDictSource ports the private extractFromDssDictSource().
func (s *PAdESCertificateSource) extractFromDssDictSource() {
	for _, certToken := range s.DSSDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOrigin_DSS_DICTIONARY)
	}
	for _, certToken := range s.VRIDictionaryCertValues() {
		s.AddCertificateWithOrigin(certToken, enumerations.CertificateOrigin_VRI_DICTIONARY)
	}
}

// CertificateMap gets the map of certificate PDF object ids and the certificateTokens.
// Port of getCertificateMap().
func (s *PAdESCertificateSource) CertificateMap() map[PdfObjectKey]*model.CertificateToken {
	return s.dssDictionaryCertificateSource.CertificateMap()
}

// CertificateValues is not applicable for PAdES. Port of the getCertificateValues() override.
func (s *PAdESCertificateSource) CertificateValues() []*model.CertificateToken {
	return []*model.CertificateToken{}
}

// CompleteCertificateRefs is not applicable for PAdES. Port of the getCompleteCertificateRefs()
// override.
func (s *PAdESCertificateSource) CompleteCertificateRefs() []*spi.CertificateRef {
	return []*spi.CertificateRef{}
}

// AttributeCertificateRefs is not applicable for PAdES. Port of the
// getAttributeCertificateRefs() override.
func (s *PAdESCertificateSource) AttributeCertificateRefs() []*spi.CertificateRef {
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
func (s *PAdESCertificateSource) DSSDictionaryCertValues() []*model.CertificateToken {
	return s.dssDictionaryCertificateSource.DSSDictionaryCertValues()
}

// VRIDictionaryCertValues gets the list of the certificate tokens extracted from the VRI
// dictionary. Port of the getVRIDictionaryCertValues() override. Shadows the embedded base's
// method of the same name; see DSSDictionaryCertValues's doc comment.
func (s *PAdESCertificateSource) VRIDictionaryCertValues() []*model.CertificateToken {
	return s.dssDictionaryCertificateSource.VRIDictionaryCertValues()
}

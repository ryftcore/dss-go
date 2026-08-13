// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSBuilder.java (DSS 6.5.RC1).
//
// org.bouncycastle.cert.jcajce.JcaCertStore is dropped: it merely re-encodes a
// java.security.cert.X509Certificate collection into a Store<X509CertificateHolder>, which here
// is already the [][]byte convention every certificate store uses (see doc.go). Deduplication
// ("if (!certificates.contains(certificateToken))") becomes de-duplication by
// CertificateToken.Equals, since a Go byte slice used as a de-dup key would be comparing raw
// certificate bytes rather than the CertificateToken identity Java compares - the two agree for
// well-formed input.
package cms

import (
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// CMSBuilder builds a CMS. Port of the CMSBuilder class.
type CMSBuilder struct {
	// signingCertificate is the signing-certificate to generate CMS with.
	signingCertificate *model.CertificateToken
	// certificateChain is the certificate chain to be incorporated within
	// SignedData.certificates field.
	certificateChain []*model.CertificateToken
	// generateWithoutCertificates defines whether a CMS should be generated without
	// certificates inside.
	generateWithoutCertificates bool
	// trustedCertificateSource contains a list of trusted certificate sources (see
	// trustAnchorBPPolicy).
	trustedCertificateSource spi.CertificateSource
	// trustAnchorBPPolicy indicates whether a trust anchor policy should be used.
	trustAnchorBPPolicy bool
	// originalCMS is the original CMS to be used on creation of a new CMS in a way that all
	// original field values will be copied to a new CMS.
	originalCMS *CMS
	// encapsulate defines whether a signer content shall be encapsulated to a CMS.
	encapsulate bool
}

// NewCMSBuilder is the default constructor for CMSBuilder.
func NewCMSBuilder() *CMSBuilder {
	return &CMSBuilder{trustAnchorBPPolicy: true, encapsulate: true}
}

// SetSigningCertificate sets a signing-certificate to be used for CMS generation. Port of
// #setSigningCertificate.
func (b *CMSBuilder) SetSigningCertificate(signingCertificate *model.CertificateToken) *CMSBuilder {
	b.signingCertificate = signingCertificate
	return b
}

// SetCertificateChain sets a collection of certificates to be incorporated within
// SignedData.certificates field. Port of #setCertificateChain.
func (b *CMSBuilder) SetCertificateChain(certificateChain []*model.CertificateToken) *CMSBuilder {
	b.certificateChain = certificateChain
	return b
}

// SetGenerateWithoutCertificates sets whether CMS is to be generated without certificates
// inside. Default: false (an attempt to generate without certificates will result in an
// error). Port of #setGenerateWithoutCertificates.
func (b *CMSBuilder) SetGenerateWithoutCertificates(generateWithoutCertificates bool) *CMSBuilder {
	b.generateWithoutCertificates = generateWithoutCertificates
	return b
}

// SetTrustedCertificateSource sets a trusted certificate source. See trustAnchorBPPolicy for
// more details. Port of #setTrustedCertificateSource.
func (b *CMSBuilder) SetTrustedCertificateSource(trustedCertificateSource spi.CertificateSource) *CMSBuilder {
	b.trustedCertificateSource = trustedCertificateSource
	return b
}

// SetTrustAnchorBPPolicy sets whether a B-level trust anchor policy should be used. When
// enabled, the trust anchor is not included in the generated certificate chain. Otherwise, the
// chain is generated up to a trust anchor, including the trust anchor itself. Default: true.
// Port of #setTrustAnchorBPPolicy.
func (b *CMSBuilder) SetTrustAnchorBPPolicy(trustAnchorBPPolicy bool) *CMSBuilder {
	b.trustAnchorBPPolicy = trustAnchorBPPolicy
	return b
}

// SetOriginalCMS sets the original CMS, whose internal field values will be copied to a new
// CMS. Port of #setOriginalCMS.
func (b *CMSBuilder) SetOriginalCMS(originalCMS *CMS) *CMSBuilder {
	b.originalCMS = originalCMS
	return b
}

// SetEncapsulate sets whether a signer content shall be encapsulated to the CMS. When enabled
// creates an enveloping signature, otherwise creates a detached signature. Default: true. Port
// of #setEncapsulate.
func (b *CMSBuilder) SetEncapsulate(encapsulate bool) *CMSBuilder {
	b.encapsulate = encapsulate
	return b
}

// CreateCMS builds a CMS. Port of #createCMS(SignerInfoGenerator, DSSDocument).
func (b *CMSBuilder) CreateCMS(signerInfoGenerator *SignerInfoGenerator, toSignDocument model.DSSDocument) (*CMS, error) {
	generator := NewAbstractCMSGenerator()

	generator.SetSignerInfoGenerator(signerInfoGenerator)
	certificateStore, err := b.certificateStore()
	if err != nil {
		return nil, err
	}
	generator.SetCertificates(certificateStore)
	generator.SetDigestAlgorithmIDs(b.digestAlgorithmIDs(signerInfoGenerator))

	generator.SetToBeSignedDocument(toSignDocument)
	generator.SetEncapsulate(b.encapsulate)

	if b.originalCMS != nil {
		generator.SetSigners(b.originalCMS.SignerInfos())
		generator.SetAttributeCertificates(b.originalCMS.AttributeCertificates())
		generator.SetCRLs(b.originalCMS.CRLs())
		generator.SetOcspBasicStore(b.originalCMS.OcspBasicStore())
		generator.SetOcspResponsesStore(b.originalCMS.OcspResponseStore())
	}

	return generator.Generate()
}

// certificateStore returns the certificate store: the original CMS's own certificates (when
// one is set), deduplicated with the ones getJcaCertStore selects. Port of the private
// #getCertificateStore.
func (b *CMSBuilder) certificateStore() ([][]byte, error) {
	var certificates []*model.CertificateToken
	if b.originalCMS != nil {
		for _, encoded := range b.originalCMS.Certificates() {
			token, err := spi.DSSUtilsLoadCertificateFromBinary(encoded)
			if err != nil {
				return nil, err
			}
			if !containsCertificateToken(certificates, token) {
				certificates = append(certificates, token)
			}
		}
	}
	certificates, err := b.addJcaCertStoreCertificates(certificates)
	if err != nil {
		return nil, err
	}
	encoded := make([][]byte, len(certificates))
	for index, certificate := range certificates {
		encoded[index] = certificate.Encoded()
	}
	return encoded, nil
}

// addJcaCertStoreCertificates appends the certificate chain for the B-level signature creation
// to certificates, skipping any already present. Port of the private #getJcaCertStore, minus
// the JcaCertStore wrapping (see the file header).
//
// "The order of the certificates is important, the fist one must be the signing certificate."
func (b *CMSBuilder) addJcaCertStoreCertificates(certificates []*model.CertificateToken) ([]*model.CertificateToken, error) {
	var certificatesToAdd []*model.CertificateToken
	if b.signingCertificate == nil && b.generateWithoutCertificates {
		certificatesToAdd = nil
	} else {
		selector := spi.NewBaselineBCertificateSelector(b.signingCertificate, b.certificateChain).
			SetTrustedCertificateSource(b.trustedCertificateSource).
			SetTrustAnchorBPPolicy(b.trustAnchorBPPolicy)
		selected, err := selector.Certificates()
		if err != nil {
			return nil, err
		}
		certificatesToAdd = selected
	}

	for _, certificateToken := range certificatesToAdd {
		if !containsCertificateToken(certificates, certificateToken) {
			certificates = append(certificates, certificateToken)
		}
	}
	return certificates, nil
}

// containsCertificateToken is Collection<CertificateToken>#contains.
func containsCertificateToken(certificates []*model.CertificateToken, candidate *model.CertificateToken) bool {
	for _, certificate := range certificates {
		if certificate.Equals(candidate) {
			return true
		}
	}
	return false
}

// digestAlgorithmIDs ports the private #getDigestAlgorithmIDs(SignerInfoGenerator): the
// original CMS's digest algorithm identifiers, plus the new signer's own.
func (b *CMSBuilder) digestAlgorithmIDs(signerInfoGenerator *SignerInfoGenerator) []*asn1ber.AlgorithmIdentifier {
	var digestAlgorithmIDs []*asn1ber.AlgorithmIdentifier
	if b.originalCMS != nil {
		digestAlgorithmIDs = mergeAlgorithmIdentifiers(digestAlgorithmIDs, b.originalCMS.DigestAlgorithmIDs())
	}
	digestAlgorithmIDs = mergeAlgorithmIdentifiers(digestAlgorithmIDs, []*asn1ber.AlgorithmIdentifier{signerInfoGenerator.DigestAlgorithm()})
	return digestAlgorithmIDs
}

// ExtendCMSSignedData extends the provided originalCMS with the required validation data. Port
// of #extendCMSSignedData(Collection<CertificateToken>, Collection<CRLToken>, Collection<OCSPToken>).
//
// Panics with the Java message when originalCMS was not set (Java's NullPointerException raised
// directly, not a DSSException, so the port matches with a panic rather than an error).
func (b *CMSBuilder) ExtendCMSSignedData(certificateTokens []*model.CertificateToken, crlTokens []*spi.CRLToken,
	ocspTokens []*spi.OCSPToken) (*CMS, error) {
	if b.originalCMS == nil {
		panic("Original CMSSignedData shall be provided! Use #setOriginalCMSSignedData(CMSSignedData) method.")
	}

	certificates := newUniqueByteSet(b.originalCMS.Certificates())
	for _, certificateToken := range certificateTokens {
		certificates.add(certificateToken.Encoded())
	}

	crls := newUniqueByteSet(b.originalCMS.CRLs())
	for _, crlToken := range crlTokens {
		crls.add(crlToken.Encoded())
	}

	ocspResponses := newUniqueByteSet(b.originalCMS.OcspResponseStore())
	for _, ocspToken := range ocspTokens {
		ocspResponses.add(ocspToken.Encoded())
	}

	return CMSUtilsReplaceCertificatesAndCRLs(b.originalCMS, certificates.values, b.originalCMS.AttributeCertificates(),
		crls.values, ocspResponses.values, b.originalCMS.OcspBasicStore())
}

// uniqueByteSet is a []byte counterpart of Java's LinkedHashSet<Encodable>: it preserves
// insertion order and skips a member already present, comparing by content.
type uniqueByteSet struct {
	values [][]byte
}

// newUniqueByteSet seeds the set with a defensive copy of initial, mirroring
// `new LinkedHashSet<>(store.getMatches(null))`.
func newUniqueByteSet(initial [][]byte) *uniqueByteSet {
	values := make([][]byte, len(initial))
	copy(values, initial)
	return &uniqueByteSet{values: values}
}

// add appends value when it is not already present.
func (s *uniqueByteSet) add(value []byte) {
	for _, member := range s.values {
		if string(member) == string(value) {
			return
		}
	}
	s.values = append(s.values, value)
}

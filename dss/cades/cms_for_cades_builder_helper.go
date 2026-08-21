// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CMSForCAdESBuilderHelper.java (DSS 6.5.RC1).
//
// The Java class is a thin assembly layer over the dss-cms builders, so the port is a thin layer
// over the Go cms package: CMSBuilder, CMSSignerInfoGeneratorBuilder, SignerInfoGenerator and
// ContentSigner keep their Java names there.
//
// Java's protected methods exist so that a subclass (dss-asic-cades has one) can override the
// profile or the builder configuration; Go has no method overriding across embedding, so they
// are exported methods here and a specialised helper embeds this one and re-implements what it
// needs, calling the embedded implementation for the rest.
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CMSForCAdESBuilderHelper is used to build an instance of cms.CMS for a CAdES Baseline B
// creation.
type CMSForCAdESBuilderHelper struct {
	// DocumentToSign is the document to be signed by the CAdES signature.
	DocumentToSign model.DSSDocument

	// SignatureParameters are the signature parameters used on the signature creation.
	SignatureParameters *CAdESSignatureParameters

	// ContentSigner is the content signer used for the signature creation.
	ContentSigner cms.ContentSigner

	// originalCMS is the original CMS, when available.
	originalCMS *cms.CMS

	// trustedCertificateSource is the certificate source containing trust anchors.
	trustedCertificateSource spi.CertificateSource

	// includeUnsignedAttributes defines whether the unsigned attributes should be included to a
	// generated SignerInfoGenerator.
	includeUnsignedAttributes bool

	// counterSignature marks the parameters as counter-signature parameters; see
	// CAdESLevelBaselineB.counterSignature, which is where upstream's
	// "parameters instanceof CAdESCounterSignatureParameters" test lives.
	counterSignature bool

	// cadesProfile is the cached instance of a CAdES profile.
	cadesProfile *CAdESLevelBaselineB
}

// NewCMSForCAdESBuilderHelper is the default constructor. Port of
// CMSForCAdESBuilderHelper(DSSDocument, CAdESSignatureParameters, ContentSigner); a nil argument
// panics with Java's Objects.requireNonNull message.
func NewCMSForCAdESBuilderHelper(documentToSign model.DSSDocument, signatureParameters *CAdESSignatureParameters,
	contentSigner cms.ContentSigner) *CMSForCAdESBuilderHelper {
	if documentToSign == nil {
		panic("documentToSign cannot be null!")
	}
	if signatureParameters == nil {
		panic("signatureParameters cannot be null!")
	}
	if contentSigner == nil {
		panic("contentSigner cannot be null!")
	}
	return &CMSForCAdESBuilderHelper{
		DocumentToSign:      documentToSign,
		SignatureParameters: signatureParameters,
		ContentSigner:       contentSigner,
	}
}

// SetOriginalCMS sets the original CMS, when available. Port of #setOriginalCMS.
func (h *CMSForCAdESBuilderHelper) SetOriginalCMS(originalCMS *cms.CMS) *CMSForCAdESBuilderHelper {
	h.originalCMS = originalCMS
	return h
}

// SetTrustedCertificateSource sets the trusted certificate source.
// Port of #setTrustedCertificateSource.
func (h *CMSForCAdESBuilderHelper) SetTrustedCertificateSource(trustedCertificateSource spi.CertificateSource) *CMSForCAdESBuilderHelper {
	h.trustedCertificateSource = trustedCertificateSource
	return h
}

// SetIncludeUnsignedAttributes sets whether the unsigned attributes should be included into the
// generated SignerInfoGenerator. Port of #setIncludeUnsignedAttributes.
func (h *CMSForCAdESBuilderHelper) SetIncludeUnsignedAttributes(includeUnsignedAttributes bool) *CMSForCAdESBuilderHelper {
	h.includeUnsignedAttributes = includeUnsignedAttributes
	return h
}

// SetCounterSignature marks the signature parameters as counter-signature parameters, which
// suppresses the mime-type signed attribute. It has no upstream counterpart: Java tests
// "signatureParameters instanceof CAdESCounterSignatureParameters" inside
// CAdESLevelBaselineB#addMimeType, and Go's embedding does not preserve that dynamic type across
// the *CAdESSignatureParameters this helper is given. CAdESCounterSignatureBuilder is the sole
// caller.
func (h *CMSForCAdESBuilderHelper) SetCounterSignature(counterSignature bool) *CMSForCAdESBuilderHelper {
	h.counterSignature = counterSignature
	return h
}

// CreateCMS creates a CMS using the ContentSigner. Port of #createCMS.
func (h *CMSForCAdESBuilderHelper) CreateCMS() (*cms.CMS, error) {
	if h.ContentSigner == nil {
		panic("contentSigner cannot be null!")
	}
	cmsBuilder := h.InitCMSBuilder()
	signerInfoGenerator, err := h.CreateSignerInfoGenerator()
	if err != nil {
		return nil, err
	}
	return cmsBuilder.CreateCMS(signerInfoGenerator, h.DocumentToSign)
}

// CreateSignerInfoGenerator creates a SignerInfoGenerator for a CAdES creation.
// Port of #createSignerInfoGenerator.
func (h *CMSForCAdESBuilderHelper) CreateSignerInfoGenerator() (*cms.SignerInfoGenerator, error) {
	if err := h.AssertSignatureParametersValid(); err != nil {
		return nil, err
	}

	signedAttributes, err := h.InitSignedAttributesTable()
	if err != nil {
		return nil, err
	}
	unsignedAttributes := h.InitUnsignedAttributesTable()

	signerInfoGeneratorBuilder := h.CreateCMSSignerInfoGeneratorBuilder(signedAttributes, unsignedAttributes)
	return signerInfoGeneratorBuilder.Build(h.DocumentToSign, h.ContentSigner)
}

// InitSignedAttributesTable creates a signed attributes table for the CAdES Baseline B creation.
// Port of the protected #initSignedAttributesTable.
func (h *CMSForCAdESBuilderHelper) InitSignedAttributesTable() (cmscore.Attributes, error) {
	cadesProfile := h.CAdESProfile()
	return cadesProfile.SignedAttributes(h.SignatureParameters)
}

// InitUnsignedAttributesTable creates an unsigned attributes table for the CAdES Baseline B
// creation, nil when the unsigned attributes are not to be included.
// Port of the protected #initUnsignedAttributesTable.
func (h *CMSForCAdESBuilderHelper) InitUnsignedAttributesTable() cmscore.Attributes {
	if h.includeUnsignedAttributes {
		cadesProfile := h.CAdESProfile()
		return cadesProfile.UnsignedAttributes()
	}
	return nil
}

// CAdESProfile gets the CAdESLevelBaselineB used for the signed and unsigned attributes table
// creation. Port of the protected #getCAdESProfile.
func (h *CMSForCAdESBuilderHelper) CAdESProfile() *CAdESLevelBaselineB {
	if h.cadesProfile == nil {
		h.cadesProfile = h.InitCAdESProfile()
	}
	return h.cadesProfile
}

// InitCAdESProfile instantiates a new CAdESLevelBaselineB.
// Port of the protected #initCAdESProfile.
func (h *CMSForCAdESBuilderHelper) InitCAdESProfile() *CAdESLevelBaselineB {
	profile := NewCAdESLevelBaselineBWithDocument(h.DocumentToSign)
	profile.SetCounterSignature(h.counterSignature)
	return profile
}

// CreateCMSSignerInfoGeneratorBuilder creates and configures a CMSSignerInfoGeneratorBuilder to
// be used for a SignerInfo creation.
// Port of the protected #createCMSSignerInfoGeneratorBuilder.
func (h *CMSForCAdESBuilderHelper) CreateCMSSignerInfoGeneratorBuilder(signedAttributes,
	unsignedAttributes cmscore.Attributes) *cms.CMSSignerInfoGeneratorBuilder {
	return h.InitCMSSignerInfoGeneratorBuilder().
		SetSigningCertificate(h.SignatureParameters.SigningCertificate()).
		SetDigestAlgorithm(h.SignatureParameters.ReferenceDigestAlgorithm()).
		SetSignedAttributes(signedAttributes).
		SetUnsignedAttributes(unsignedAttributes)
}

// InitCMSSignerInfoGeneratorBuilder creates a new instance of CMSSignerInfoGeneratorBuilder.
// Port of the protected #initCMSSignerInfoGeneratorBuilder.
func (h *CMSForCAdESBuilderHelper) InitCMSSignerInfoGeneratorBuilder() *cms.CMSSignerInfoGeneratorBuilder {
	return cms.NewCMSSignerInfoGeneratorBuilder()
}

// AssertSignatureParametersValid verifies the validity of the signature parameters
// configuration. Port of the protected #assertSignatureParametersValid; Java's unchecked
// IllegalArgumentException is a returned error here, since the whole creation path already has
// an error channel.
func (h *CMSForCAdESBuilderHelper) AssertSignatureParametersValid() error {
	if h.SignatureParameters.SigningCertificate() == nil && !h.SignatureParameters.GenerateTBSWithoutCertificate() {
		return model.NewDSSError("Signing-certificate is not provided! " +
			"Use #setGenerateWithoutCertificates(true) method.")
	}
	return nil
}

// InitCMSBuilder instantiates a CMSBuilder for the CMS creation.
// Port of the protected #initCMSBuilder.
func (h *CMSForCAdESBuilderHelper) InitCMSBuilder() *cms.CMSBuilder {
	return cms.NewCMSBuilder().
		SetSigningCertificate(h.SignatureParameters.SigningCertificate()).
		SetCertificateChain(h.SignatureParameters.CertificateChain()).
		SetGenerateWithoutCertificates(h.SignatureParameters.GenerateTBSWithoutCertificate()).
		SetTrustAnchorBPPolicy(h.SignatureParameters.BLevel().IsTrustAnchorBPPolicy()).
		SetTrustedCertificateSource(h.trustedCertificateSource).
		SetEncapsulate(h.IsEncapsulateSignerData()).
		SetOriginalCMS(h.originalCMS)
}

// IsEncapsulateSignerData reports whether the signed data shall be encapsulated.
// Port of the protected #isEncapsulateSignerData.
func (h *CMSForCAdESBuilderHelper) IsEncapsulateSignerData() bool {
	return enumerations.SignaturePackaging_DETACHED != h.SignatureParameters.SignaturePackaging()
}

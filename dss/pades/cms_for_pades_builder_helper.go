// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/CMSForPAdESBuilderHelper.java (DSS 6.5.RC1).
//
// Java extends cades.CMSForCAdESBuilderHelper and overrides initCAdESProfile (to return a
// LevelBaselineB) and isEncapsulateSignerData (always false), plus three covariant-return
// setters. Go has no method overriding across embedding: the embedded helper's CreateCMS and
// CreateSignerInfoGenerator would call the base's own InitCAdESProfile / IsEncapsulateSignerData.
// As cades/cms_for_cades_builder_helper.go's header prescribes for a specialised helper, the
// methods that sit on the path from CreateCMS down to the two overridden hooks are re-declared
// here and delegate to the embedded implementation for everything else.
//
// The base keeps originalCMS, trustedCertificateSource and includeUnsignedAttributes unexported,
// so the setters record them here as well: this type's InitCMSBuilder / InitUnsignedAttributesTable
// need to read them back, which Go cannot do through the embedded struct.
package pades

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// CMSForPAdESBuilderHelper creates a new CMS for a PAdES signature creation.
type CMSForPAdESBuilderHelper struct {
	cades.CMSForCAdESBuilderHelper

	// messageDigest is the message-digest computed on the PAdES revision to be signed.
	messageDigest model.DSSMessageDigest

	// signatureParameters is the PAdES view of the parameters handed to the constructor. The
	// embedded helper only keeps the CAdES projection of them.
	signatureParameters *SignatureParameters

	// originalCMS mirrors the base field of the same name; see the file header.
	originalCMS *cms.CMS

	// trustedCertificateSource mirrors the base field of the same name.
	trustedCertificateSource spi.CertificateSource

	// includeUnsignedAttributes mirrors the base field of the same name.
	includeUnsignedAttributes bool

	// padesProfile is the cached instance of the PAdES profile, i.e. the base's cadesProfile.
	padesProfile *LevelBaselineB
}

// NewCMSForPAdESBuilderHelper is the default constructor.
// Port of CMSForPAdESBuilderHelper(DSSMessageDigest, PAdESSignatureParameters, ContentSigner);
// Java passes DSSUtils.toDigestDocument(messageDigest) as the document to sign.
func NewCMSForPAdESBuilderHelper(messageDigest model.DSSMessageDigest,
	signatureParameters *SignatureParameters, contentSigner cms.ContentSigner) *CMSForPAdESBuilderHelper {
	digestDocument := spi.DSSUtilsToDigestDocument(messageDigest.Digest)
	return &CMSForPAdESBuilderHelper{
		CMSForCAdESBuilderHelper: *cades.NewCMSForCAdESBuilderHelper(digestDocument,
			&signatureParameters.SignatureParameters, contentSigner),
		messageDigest:       messageDigest,
		signatureParameters: signatureParameters,
	}
}

// SetOriginalCMS sets the original CMS, when available. Port of the covariant #setOriginalCMS.
func (h *CMSForPAdESBuilderHelper) SetOriginalCMS(originalCMS *cms.CMS) *CMSForPAdESBuilderHelper {
	h.CMSForCAdESBuilderHelper.SetOriginalCMS(originalCMS)
	h.originalCMS = originalCMS
	return h
}

// SetTrustedCertificateSource sets the trusted certificate source.
// Port of the covariant #setTrustedCertificateSource.
func (h *CMSForPAdESBuilderHelper) SetTrustedCertificateSource(trustedCertificateSource spi.CertificateSource) *CMSForPAdESBuilderHelper {
	h.CMSForCAdESBuilderHelper.SetTrustedCertificateSource(trustedCertificateSource)
	h.trustedCertificateSource = trustedCertificateSource
	return h
}

// SetIncludeUnsignedAttributes sets whether the unsigned attributes should be included into the
// generated SignerInfoGenerator. Port of the covariant #setIncludeUnsignedAttributes.
func (h *CMSForPAdESBuilderHelper) SetIncludeUnsignedAttributes(includeUnsignedAttributes bool) *CMSForPAdESBuilderHelper {
	h.CMSForCAdESBuilderHelper.SetIncludeUnsignedAttributes(includeUnsignedAttributes)
	h.includeUnsignedAttributes = includeUnsignedAttributes
	return h
}

// CreateCMS creates a CMS using the ContentSigner. Port of the inherited #createCMS; re-declared
// so that the two hooks this type overrides are reached.
func (h *CMSForPAdESBuilderHelper) CreateCMS() (*cms.CMS, error) {
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

// CreateSignerInfoGenerator creates a SignerInfoGenerator for the PAdES creation.
// Port of the inherited #createSignerInfoGenerator.
func (h *CMSForPAdESBuilderHelper) CreateSignerInfoGenerator() (*cms.SignerInfoGenerator, error) {
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

// InitSignedAttributesTable creates the signed attributes table using the PAdES profile.
// Port of the inherited protected #initSignedAttributesTable.
func (h *CMSForPAdESBuilderHelper) InitSignedAttributesTable() (cmscore.Attributes, error) {
	return h.PAdESProfile().SignedAttributes(h.SignatureParameters)
}

// InitUnsignedAttributesTable creates an unsigned attributes table, nil when the unsigned
// attributes are not to be included. Port of the inherited protected #initUnsignedAttributesTable.
func (h *CMSForPAdESBuilderHelper) InitUnsignedAttributesTable() cmscore.Attributes {
	if h.includeUnsignedAttributes {
		return h.PAdESProfile().UnsignedAttributes()
	}
	return nil
}

// PAdESProfile gets the cached LevelBaselineB used for the attribute tables.
// Port of the inherited protected #getCAdESProfile with this type's #initCAdESProfile.
func (h *CMSForPAdESBuilderHelper) PAdESProfile() *LevelBaselineB {
	if h.padesProfile == nil {
		h.padesProfile = h.InitCAdESProfile()
	}
	return h.padesProfile
}

// InitCAdESProfile instantiates the PAdES Baseline B profile carrying the revision's
// message-digest. Port of the protected #initCAdESProfile override.
func (h *CMSForPAdESBuilderHelper) InitCAdESProfile() *LevelBaselineB {
	return NewLevelBaselineB(h.messageDigest)
}

// InitCMSBuilder instantiates a Builder for the CMS creation.
// Port of the inherited protected #initCMSBuilder; re-declared so that IsEncapsulateSignerData
// resolves to this type's override.
func (h *CMSForPAdESBuilderHelper) InitCMSBuilder() *cms.Builder {
	return cms.NewBuilder().
		SetSigningCertificate(h.SignatureParameters.SigningCertificate()).
		SetCertificateChain(h.SignatureParameters.CertificateChain()).
		SetGenerateWithoutCertificates(h.SignatureParameters.GenerateTBSWithoutCertificate()).
		SetTrustAnchorBPPolicy(h.SignatureParameters.BLevel().IsTrustAnchorBPPolicy()).
		SetTrustedCertificateSource(h.trustedCertificateSource).
		SetEncapsulate(h.IsEncapsulateSignerData()).
		SetOriginalCMS(h.originalCMS)
}

// IsEncapsulateSignerData reports whether the signed content shall be encapsulated: never for
// PAdES, whose CMS is always detached from the PDF revision.
// Port of the protected #isEncapsulateSignerData override.
func (h *CMSForPAdESBuilderHelper) IsEncapsulateSignerData() bool {
	return false
}

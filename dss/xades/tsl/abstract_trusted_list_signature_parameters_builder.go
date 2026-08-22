// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/tsl/AbstractTrustedListSignatureParametersBuilder.java (DSS 6.5.RC1).
//
// Java's abstract isEn319132()/getTargetTLVersion() become the
// AbstractTrustedListSignatureParametersBuilderOverrides interface, dispatched through the
// InitAbstractTrustedListSignatureParametersBuilder(self, ...) registration - the
// TokenBase.InitToken(self) convention of PORTING.md - since Go has no method overriding across
// embedding. Java's build() override (mutating the AbstractSignatureParametersBuilder<SP> result
// before returning it) is likewise not an override in Go: it is this type's own Build() method,
// which first calls the embedded document.AbstractSignatureParametersBuilder's Build().
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// defaultCanonicalization is the EXCLUSIVE canonicalization method used for the enveloped
// signature reference. See TS 119 612 "B.1 The Signature element".
const defaultCanonicalization = xmlutils.XMLCanonicalizerDefaultDSSC14NMethod

// defaultReferencePrefix is the default prefix for an enveloped signature reference id.
const defaultReferencePrefix = "ref-enveloped-signature"

// AbstractTrustedListSignatureParametersBuilderOverrides captures the two protected abstract
// methods of Java's AbstractTrustedListSignatureParametersBuilder that concrete subclasses must
// supply.
type AbstractTrustedListSignatureParametersBuilderOverrides interface {
	// IsEn319132 reports whether the created XAdES signature shall be conformant to ETSI EN 319
	// 132 standard (new XAdES). Port of the protected abstract isEn319132().
	IsEn319132() bool

	// TargetTLVersion returns the target XML Trusted List version to be signed. Port of the
	// protected abstract getTargetTLVersion().
	TargetTLVersion() int
}

// AbstractTrustedListSignatureParametersBuilder contains common methods for signature
// parameters creation for an XML Trusted List signature.
type AbstractTrustedListSignatureParametersBuilder struct {
	document.AbstractSignatureParametersBuilder[*xades.SignatureParameters]

	// overrides points back at the concrete builder; see
	// InitAbstractTrustedListSignatureParametersBuilder.
	overrides AbstractTrustedListSignatureParametersBuilderOverrides

	// tlXmlDocument is the XML Trusted List document.
	tlXmlDocument model.DSSDocument

	// referenceId is the Enveloped reference Id to use.
	referenceId string

	// referenceDigestAlgorithm is the DigestAlgorithm to be used for an Enveloped reference.
	referenceDigestAlgorithm enumerations.DigestAlgorithm
}

// InitAbstractTrustedListSignatureParametersBuilder is the port of the protected constructor
// AbstractTrustedListSignatureParametersBuilder(CertificateToken, DSSDocument).
//
// Panics when tlXmlDocument is nil (Java's Objects.requireNonNull).
func (b *AbstractTrustedListSignatureParametersBuilder) InitAbstractTrustedListSignatureParametersBuilder(
	overrides AbstractTrustedListSignatureParametersBuilderOverrides, signingCertificate *model.CertificateToken,
	tlXmlDocument model.DSSDocument) {
	if tlXmlDocument == nil {
		panic("XML Trusted List document cannot be null!")
	}
	b.AbstractSignatureParametersBuilder = *document.NewAbstractSignatureParametersBuilder[*xades.SignatureParameters](signingCertificate)
	b.AbstractSignatureParametersBuilder.InitParameters = xades.NewXAdESSignatureParameters
	b.overrides = overrides
	b.tlXmlDocument = tlXmlDocument
	b.referenceDigestAlgorithm = enumerations.DigestAlgorithmSHA512
}

// SetReferenceId sets an Enveloped Reference Id to use.
// Default: "ref-enveloped-signature". Port of #setReferenceId, chainable.
func (b *AbstractTrustedListSignatureParametersBuilder) SetReferenceId(referenceId string) *AbstractTrustedListSignatureParametersBuilder {
	b.referenceId = referenceId
	return b
}

// SetReferenceDigestAlgorithm sets an Enveloped Reference DigestAlgorithm to use. Port of
// #setReferenceDigestAlgorithm, chainable.
func (b *AbstractTrustedListSignatureParametersBuilder) SetReferenceDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *AbstractTrustedListSignatureParametersBuilder {
	b.referenceDigestAlgorithm = digestAlgorithm
	return b
}

// Build ports the overridden build(). Note: unlike the base
// document.AbstractSignatureParametersBuilder.Build (a plain, non-overridden method in Go),
// this type's Build is its own method, not a Go embedding override - callers on this type reach
// it directly, and V5/V6 subtypes do not need to redeclare it since they never change the
// return type.
func (b *AbstractTrustedListSignatureParametersBuilder) Build() *xades.SignatureParameters {
	signatureParameters := b.AbstractSignatureParametersBuilder.Build()

	signatureParameters.SetSignaturePackaging(enumerations.SignaturePackagingEnveloped)
	signatureParameters.SetSignatureLevel(enumerations.SignatureLevelXAdESBaselineB)
	signatureParameters.SetEn319132(b.overrides.IsEn319132())

	references := b.GetReferences()
	signatureParameters.SetReferences(references)

	return signatureParameters
}

// GetReferences returns a list of ds:References to be incorporated within the signature. Port
// of the protected #getReferences().
func (b *AbstractTrustedListSignatureParametersBuilder) GetReferences() []*xades.DSSReference {
	envelopedSignatureReference := b.GetEnvelopedSignatureReference()
	return []*xades.DSSReference{envelopedSignatureReference}
}

// GetEnvelopedSignatureReference creates the enveloped-signature ds:Reference. Port of the
// protected #getEnvelopedSignatureReference().
func (b *AbstractTrustedListSignatureParametersBuilder) GetEnvelopedSignatureReference() *xades.DSSReference {
	dssReference := xades.NewDSSReference()
	if b.referenceId != "" {
		dssReference.SetId(b.referenceId)
	} else {
		dssReference.SetId(defaultReferencePrefix)
	}
	dssReference.SetUri("")
	dssReference.SetContents(b.tlXmlDocument)
	dssReference.SetDigestMethodAlgorithm(b.referenceDigestAlgorithm)

	transforms := []xades.DSSTransform{
		xades.NewEnvelopedSignatureTransform(),
		xades.NewCanonicalizationTransform(defaultCanonicalization),
	}
	dssReference.SetTransforms(transforms)
	return dssReference
}

// AssertConfigurationIsValid verifies whether the chosen signature parameters builder is
// applicable to the given document. It verifies whether the provided document representing the
// XML Trusted List is conformant to the definition and the target version.
// NOTE: this method requires 'dss-validation' module.
// Port of #assertConfigurationIsValid().
func (b *AbstractTrustedListSignatureParametersBuilder) AssertConfigurationIsValid() error {
	errors, err := XAdESTrustedListUtilsValidateUnsignedTrustedList(b.tlXmlDocument, b.overrides.TargetTLVersion())
	if err != nil {
		return fmt.Errorf("an error occurred on XML Trusted List validation : %s", err.Error())
	}
	if utils.IsCollectionNotEmpty(errors) {
		return exception.NewIllegalInputException(fmt.Sprintf(
			"XML Trusted List failed the validation : %s", utils.JoinStrings(errors, "; ")))
	}
	return nil
}

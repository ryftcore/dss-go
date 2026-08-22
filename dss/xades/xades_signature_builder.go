// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESSignatureBuilder.java (DSS 6.5.RC1).
//
// # The abstract base, and how its virtual methods survive the port
//
// Java's XAdESSignatureBuilder is abstract: build() drives a fixed sequence of steps, six of
// which its packaging subclasses (Enveloped, Enveloping, Detached, InternallyDetached and the
// intermediate XPathPlacementSignatureBuilder) override. Go has no method overriding across
// embedding, so - exactly as xades_builder.go, document/abstract_document_extender.go and
// cades/cades_signature_extension.go already do, per the TokenBase.InitToken(self) convention of
// PORTING.md - those six are collected in SignatureBuilderOverrides, every concrete builder
// registers itself with InitXAdESSignatureBuilder, and this base always reaches them through
// b.overrides. The base's own implementations stay as methods, so a subclass that does not
// override a hook inherits it by embedding and still satisfies the interface.
//
// A grep over dss-xades gives the closed set of overridden members:
//
//	assertSignaturePossible          Enveloped, InternallyDetached
//	incorporateFiles                 InternallyDetached
//	buildRootDocumentDom             Enveloped, Detached, InternallyDetached
//	getParentNodeOfSignature         Detached, XPathPlacement
//	incorporateSignatureDom(Node)    XPathPlacement
//	incorporateSignedObjects         Enveloping
//	alignNodes                       (BuilderOverrides; implemented here)
//
// Java's two incorporateSignatureDom overloads cannot share one Go name: the public no-argument
// one keeps IncorporateSignatureDom, the protected Node-taking one - the overridden one - becomes
// IncorporateSignatureDomToParent.
//
// # Errors
//
// build() and every incorporate* step it drives return errors where Java throws DSSException,
// IllegalArgumentException or UnsupportedOperationException (PORTING.md: throw -> (T, error)).
// Objects.requireNonNull keeps panicking with the Java message. Serialization, canonicalization
// and digesting, which Java performs behind DSSException-free APIs, can fail in Go and are
// likewise reported. slf4j logging is dropped.
//
// # instanceof XmlPolicyWithTransforms
//
// incorporatePolicy() branches on `signaturePolicy instanceof XmlPolicyWithTransforms`, but
// model.BLevelParameters stores the policy as a concrete *model.Policy and XmlPolicyWithTransforms
// embeds model.Policy by value, so Go cannot recover the outer value from the pointer. The branch
// is therefore resolved through the registry below; see
// SignatureBuilderRegisterPolicyTransforms. FLAGGED FOR THE INTEGRATOR: the structurally
// clean fix is for model.BLevelParameters to hold the signature policy behind an interface, which
// is a change to a frozen package.
package xades

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

const (
	// SignatureBuilderReferencePrefix is the Id-prefix for the ds:Reference element.
	// Port of the protected REFERENCE_PREFIX.
	SignatureBuilderReferencePrefix = "r-"

	// AbstractSignatureBuilderKeyInfoPrefix is the Id-prefix for the ds:KeyInfo element.
	// Port of the protected KEYINFO_PREFIX.
	AbstractSignatureBuilderKeyInfoPrefix = "keyInfo-"

	// AbstractSignatureBuilderValuePrefix is the Id-prefix for the ds:SignatureValue element.
	// Port of the protected VALUE_PREFIX.
	AbstractSignatureBuilderValuePrefix = "value-"

	// AbstractSignatureBuilderXAdESPrefix is the Id-prefix for the xades:SignedProperties element.
	// Port of the protected XADES_PREFIX.
	AbstractSignatureBuilderXAdESPrefix = "xades-"
)

// XAdESSignatureBuilderOverrides declares the members Java's abstract XAdESSignatureBuilder lets
// its packaging subclasses override and that build() calls back into. Every concrete builder
// satisfies this - inheriting the base implementation by embedding for the hooks it does not
// override - and registers itself through InitXAdESSignatureBuilder.
type SignatureBuilderOverrides interface {
	BuilderOverrides

	// AssertSignaturePossible verifies whether the provided documents allow signature creation
	// with the given signature format. Port of the protected #assertSignaturePossible.
	AssertSignaturePossible() error

	// IncorporateFiles incorporates the provided documents within the final file.
	// Port of the protected #incorporateFiles.
	IncorporateFiles() error

	// BuildRootDocumentDom builds the root Document the signature is created in.
	// Port of the protected #buildRootDocumentDom.
	BuildRootDocumentDom() *xmldom.Node

	// ParentNodeOfSignature returns the parent node of the signature.
	// Port of the protected #getParentNodeOfSignature.
	ParentNodeOfSignature() *xmldom.Node

	// IncorporateSignatureDomToParent incorporates the signature element into the parent node.
	// Port of the protected #incorporateSignatureDom(Node).
	IncorporateSignatureDomToParent(parentNodeOfSignature *xmldom.Node)

	// IncorporateSignedObjects incorporates the list of signed ds:Object elements (used for
	// Enveloping packaging). Port of the protected #incorporateSignedObjects.
	IncorporateSignedObjects() error
}

// AbstractSignatureBuilder implements all the necessary mechanisms to build each form of the XML
// signature. It is the abstract base of the four packaging builders.
type AbstractSignatureBuilder struct {
	Builder

	// Built indicates if the signature was already built (two-steps building).
	// Port of the protected built.
	Built bool

	// Documents is the list of original documents to be signed.
	// Port of the protected documents.
	Documents []model.DSSDocument

	// KeyInfoCanonicalizationMethod is the canonicalization method used for KeyInfo signing.
	KeyInfoCanonicalizationMethod string

	// SignedInfoCanonicalizationMethod is the canonicalization method used for SignedInfo signing.
	SignedInfoCanonicalizationMethod string

	// SignedPropertiesCanonicalizationMethod is the canonicalization method used for
	// SignedProperties signing.
	SignedPropertiesCanonicalizationMethod string

	// DeterministicId is the deterministic Id used for elements creation.
	DeterministicId string

	// SignatureDom is the current DOM signature object.
	SignatureDom *xmldom.Node

	// KeyInfoDom is the cached ds:KeyInfo element.
	KeyInfoDom *xmldom.Node

	// SignedInfoDom is the cached ds:SignedInfo element.
	SignedInfoDom *xmldom.Node

	// SignatureValueDom is the cached ds:SignatureValue element.
	SignatureValueDom *xmldom.Node

	// QualifyingPropertiesDom is the cached xades:QualifyingProperties element.
	QualifyingPropertiesDom *xmldom.Node

	// SignedPropertiesDom is the cached xades:SignedProperties element.
	SignedPropertiesDom *xmldom.Node

	// SignedSignaturePropertiesDom is the cached xades:SignedSignatureProperties element.
	SignedSignaturePropertiesDom *xmldom.Node

	// SignedDataObjectPropertiesDom is the cached xades:SignedDataObjectProperties element.
	SignedDataObjectPropertiesDom *xmldom.Node

	// UnsignedSignaturePropertiesDom is the cached xades:UnsignedSignatureProperties element.
	UnsignedSignaturePropertiesDom *xmldom.Node

	// overrides points back at the concrete builder; see InitXAdESSignatureBuilder.
	overrides SignatureBuilderOverrides
}

// SignatureBuilderGetSignatureBuilder creates the signature builder according to the
// packaging, for signing a single document.
// Port of the static #getSignatureBuilder(XAdESSignatureParameters, DSSDocument, CertificateVerifier).
func SignatureBuilderGetSignatureBuilder(params *SignatureParameters,
	document model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) (SignatureBuilderRef, error) {
	return SignatureBuilderGetSignatureBuilderForDocuments(params,
		[]model.DSSDocument{document}, certificateVerifier)
}

// SignatureBuilderGetSignatureBuilderForDocuments creates the signature builder according to
// the packaging, for signing a list of documents. Java's Objects.requireNonNull panics with its
// message; the unsupported-packaging DSSException becomes a returned error.
// Port of the static #getSignatureBuilder(XAdESSignatureParameters, List<DSSDocument>, CertificateVerifier).
func SignatureBuilderGetSignatureBuilderForDocuments(params *SignatureParameters,
	documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) (SignatureBuilderRef, error) {
	if params.SignaturePackaging() == "" {
		panic("Cannot create a SignatureBuilder. SignaturePackaging shall be defined!")
	}

	switch params.SignaturePackaging() {
	case enumerations.SignaturePackagingEnveloped:
		return NewEnvelopedSignatureBuilderForDocuments(params, documents, certificateVerifier), nil
	case enumerations.SignaturePackagingEnveloping:
		return NewEnvelopingSignatureBuilderForDocuments(params, documents, certificateVerifier), nil
	case enumerations.SignaturePackagingDetached:
		return NewDetachedSignatureBuilderForDocuments(params, documents, certificateVerifier), nil
	case enumerations.SignaturePackagingInternallyDetached:
		return NewInternallyDetachedSignatureBuilderForDocuments(params, documents, certificateVerifier), nil
	default:
		return nil, model.NewDSSError("Unsupported packaging " + string(params.SignaturePackaging()))
	}
}

// InitXAdESSignatureBuilder registers the concrete builder with its base and applies the two
// protected AbstractSignatureBuilder constructors' body: it stores the parameters, the documents and
// the deterministic Id, and reads the three canonicalization methods off the parameters.
func (b *AbstractSignatureBuilder) InitXAdESSignatureBuilder(self SignatureBuilderOverrides,
	params *SignatureParameters, documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) {
	b.InitXAdESBuilderWithVerifier(self, certificateVerifier)
	b.overrides = self

	b.Params = params
	b.Documents = documents
	b.DeterministicId = params.GetDeterministicId()

	b.setCanonicalizationMethods(params)
}

// setCanonicalizationMethods ports the private setCanonicalizationMethods.
func (b *AbstractSignatureBuilder) setCanonicalizationMethods(params *SignatureParameters) {
	b.KeyInfoCanonicalizationMethod = params.KeyInfoCanonicalizationMethod()
	b.SignedInfoCanonicalizationMethod = params.SignedInfoCanonicalizationMethod()
	b.SignedPropertiesCanonicalizationMethod = params.SignedPropertiesCanonicalizationMethod()
}

// Build is the main method which is called to build the XML signature. It returns the
// canonicalized ds:SignedInfo segment of the signature, the data used to define the
// ds:SignatureValue element. Port of #build().
func (b *AbstractSignatureBuilder) Build() ([]byte, error) {
	if err := b.overrides.AssertSignaturePossible(); err != nil {
		return nil, err
	}

	if err := b.ensureConfigurationValidity(); err != nil {
		return nil, err
	}

	xadesPath, err := b.CurrentXAdESPath()
	if err != nil {
		return nil, err
	}
	b.XadesPath = xadesPath

	b.initRootDocumentDom()

	if err := b.overrides.IncorporateFiles(); err != nil {
		return nil, err
	}

	b.IncorporateSignatureDom()

	if err := b.IncorporateSignedInfo(); err != nil {
		return nil, err
	}

	b.IncorporateSignatureValue()

	if err := b.IncorporateKeyInfo(); err != nil {
		return nil, err
	}

	if err := b.IncorporateObjects(); err != nil {
		return nil, err
	}

	// We create the ds:Reference segments only now, because we need first to define the
	// SignedProperties segment to calculate the digest of references.
	if utils.IsArrayEmpty(b.Params.SignedData()) {
		if err := b.incorporateReferences(); err != nil {
			return nil, err
		}
		if err := b.IncorporateReferenceSignedProperties(); err != nil {
			return nil, err
		}
		if err := b.IncorporateReferenceKeyInfo(); err != nil {
			return nil, err
		}
	}

	// Preparation of SignedInfo
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(b.SignedInfoCanonicalizationMethod)
	if err != nil {
		return nil, err
	}
	nodeToCanonicalize, err := b.NodeToCanonicalize(b.SignedInfoDom)
	if err != nil {
		return nil, err
	}
	canonicalizedSignedInfo, err := canonicalizer.CanonicalizeNode(nodeToCanonicalize)
	if err != nil {
		return nil, err
	}
	// Upstream traces the canonicalized SignedInfo and its SHA-256 here.
	b.Built = true
	return canonicalizedSignedInfo, nil
}

// AssertSignaturePossible verifies whether the provided documents allow signature creation with
// the given signature format. Port of the protected #assertSignaturePossible.
func (b *AbstractSignatureBuilder) AssertSignaturePossible() error {
	if utils.CollectionSize(b.Documents) == 0 {
		return errors.New("No documents have been provided to the signature creation!")
	}
	return nil
}

// ensureConfigurationValidity ports the private ensureConfigurationValidity.
func (b *AbstractSignatureBuilder) ensureConfigurationValidity() error {
	if err := b.checkSignaturePackagingValidity(); err != nil {
		return err
	}

	references := b.Params.References()
	if utils.IsCollectionEmpty(references) {
		referenceBuilder := b.initReferenceBuilder()
		defaultReferences, err := referenceBuilder.Build()
		if err != nil {
			return err
		}
		// The SignatureParameters object is updated with the default references in order to
		// ensure validity on next steps
		b.Params.GetContext().SetReferences(defaultReferences)
	} else {
		referenceVerifier := NewReferenceVerifier(b.Params)
		if err := referenceVerifier.CheckReferencesValidity(); err != nil {
			return err
		}
	}
	return nil
}

// initReferenceBuilder ports the private initReferenceBuilder.
func (b *AbstractSignatureBuilder) initReferenceBuilder() *ReferenceBuilder {
	detachedContent := b.Documents
	referenceIdProvider := NewReferenceIdProvider()
	referenceIdProvider.SetSignatureParameters(b.Params)
	return NewReferenceBuilder(detachedContent, b.Params, referenceIdProvider)
}

// checkSignaturePackagingValidity ports the private checkSignaturePackagingValidity.
func (b *AbstractSignatureBuilder) checkSignaturePackagingValidity() error {
	if enumerations.SignaturePackagingEnveloping != b.Params.SignaturePackaging() {
		if b.Params.IsManifestSignature() {
			return fmt.Errorf(
				"The signature packaging %s is not compatible with manifestSignature(true) configuration!",
				b.Params.SignaturePackaging())
		}
		if b.Params.IsEmbedXML() {
			return fmt.Errorf(
				"The signature packaging %s is not compatible with embedXML(true) configuration!",
				b.Params.SignaturePackaging())
		}
	}
	return nil
}

// IncorporateFiles incorporates the provided documents within the final file. Not implemented by
// default. Port of the protected #incorporateFiles.
func (b *AbstractSignatureBuilder) IncorporateFiles() error {
	// not implemented by default
	return nil
}

// initRootDocumentDom instantiates a root Document DOM when needed.
// Port of the protected #initRootDocumentDom.
func (b *AbstractSignatureBuilder) initRootDocumentDom() {
	if b.DocumentDom == nil {
		b.DocumentDom = b.overrides.BuildRootDocumentDom()
	}
}

// BuildRootDocumentDom builds an empty Document. Port of the protected #buildRootDocumentDom.
func (b *AbstractSignatureBuilder) BuildRootDocumentDom() *xmldom.Node {
	return xmlutils.DomUtilsBuildDOMEmpty()
}

// IncorporateSignatureDom creates a new instance of the ds:Signature element and incorporates it
// into its parent node. Port of the public #incorporateSignatureDom().
func (b *AbstractSignatureBuilder) IncorporateSignatureDom() {
	b.SignatureDom = xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementSignature)
	xmlutils.DomUtilsAddNamespaceAttribute(b.SignatureDom, b.overrides.XmldsigNamespace())
	b.SignatureDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, b.DeterministicId)

	parentNodeOfSignature := b.overrides.ParentNodeOfSignature()
	b.overrides.IncorporateSignatureDomToParent(parentNodeOfSignature)
}

// ParentNodeOfSignature returns the parent node of the signature.
// Port of the protected #getParentNodeOfSignature.
func (b *AbstractSignatureBuilder) ParentNodeOfSignature() *xmldom.Node {
	return b.DocumentDom
}

// IncorporateSignatureDomToParent incorporates the signature element into the parent node.
// Port of the protected #incorporateSignatureDom(Node).
func (b *AbstractSignatureBuilder) IncorporateSignatureDomToParent(parentNodeOfSignature *xmldom.Node) {
	parentNodeOfSignature.AppendChild(b.SignatureDom)
}

// IncorporateSignedInfo incorporates the ds:SignedInfo tag:
//
//	<ds:SignedInfo>
//	    <ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
//	    <ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/>
//	    ...
//	</ds:SignedInfo>
//
// Port of the public #incorporateSignedInfo.
func (b *AbstractSignatureBuilder) IncorporateSignedInfo() error {
	if utils.IsArrayNotEmpty(b.Params.SignedData()) {
		// Upstream logs "Using explicit SignedInfo from parameter" here.
		parsed, err := xmlutils.DomUtilsBuildDOMFromBytes(b.Params.SignedData())
		if err != nil {
			return err
		}
		b.SignedInfoDom = b.DocumentDom.Import(parsed.DocumentElement(), true)
		b.SignatureDom.AppendChild(b.SignedInfoDom)
		return nil
	}

	b.SignedInfoDom = xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementSignedInfo)
	b.SignatureDom.AppendChild(b.SignedInfoDom)
	b.incorporateCanonicalizationMethod(b.SignedInfoDom, b.SignedInfoCanonicalizationMethod)

	signatureMethod := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementSignatureMethod)
	b.SignedInfoDom.AppendChild(signatureMethod)
	signatureAlgorithm := b.Params.SignatureAlgorithm()
	signatureAlgorithmXMLId := signatureAlgorithm.URI()
	if utils.IsStringBlank(signatureAlgorithmXMLId) {
		return errors.New("Unsupported signature algorithm " + string(signatureAlgorithm))
	}
	signatureMethod.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()},
		signatureAlgorithmXMLId)
	return nil
}

// incorporateCanonicalizationMethod creates the ds:CanonicalizationMethod tag:
//
//	<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
//
// Port of the private incorporateCanonicalizationMethod.
func (b *AbstractSignatureBuilder) incorporateCanonicalizationMethod(parentDom *xmldom.Node,
	signedInfoCanonicalizationMethod string) {
	canonicalizationMethodDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom,
		b.overrides.XmldsigNamespace(), common.XMLDSigElementCanonicalizationMethod)
	parentDom.AppendChild(canonicalizationMethodDom)
	canonicalizationMethodDom.SetAttr(
		xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()},
		signedInfoCanonicalizationMethod)
}

// incorporateReferences creates the ds:Reference elements in the signature.
// Port of the private incorporateReferences.
func (b *AbstractSignatureBuilder) incorporateReferences() error {
	referenceProcessor := NewReferenceProcessor(b.Params)
	return referenceProcessor.IncorporateReferences(b.SignedInfoDom, b.Params.References(),
		b.overrides.XmldsigNamespace())
}

// IncorporateKeyInfo creates the ds:KeyInfo tag.
//
// NOTE: when the trust anchor baseline profile policy is defined, only the certificates previous
// to the trust anchor are included.
//
//	<ds:KeyInfo>
//	    <ds:X509Data>
//	        <ds:X509Certificate>MIIB....</ds:X509Certificate>
//	        <ds:X509Certificate>MIIB+...</ds:X509Certificate>
//	    </ds:X509Data>
//	</ds:KeyInfo>
//
// Port of the protected #incorporateKeyInfo.
func (b *AbstractSignatureBuilder) IncorporateKeyInfo() error {
	if b.Params.SigningCertificate() == nil && b.Params.GenerateTBSWithoutCertificate() {
		// Upstream logs "Signing certificate not available and must be added to signature DOM later".
		return nil
	}

	// <ds:KeyInfo>
	keyInfoElement := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementKeyInfo)
	b.SignatureDom.AppendChild(keyInfoElement)
	if b.Params.IsSignKeyInfo() {
		keyInfoElement.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
			AbstractSignatureBuilderKeyInfoPrefix+b.DeterministicId)
	}
	certificates, err := spi.NewBaselineBCertificateSelector(b.Params.SigningCertificate(),
		b.Params.CertificateChain()).
		SetTrustedCertificateSource(b.CertificateVerifier.TrustedCertSources()).
		SetTrustAnchorBPPolicy(b.Params.BLevel().IsTrustAnchorBPPolicy()).
		Certificates()
	if err != nil {
		return err
	}

	if b.Params.IsAddX509SubjectName() {
		for _, token := range certificates {
			// <ds:X509Data>
			x509DataDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
				common.XMLDSigElementX509Data)
			keyInfoElement.AppendChild(x509DataDom)
			b.addSubjectAndCertificate(x509DataDom, token)
		}
	} else {
		// <ds:X509Data>
		x509DataDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
			common.XMLDSigElementX509Data)
		keyInfoElement.AppendChild(x509DataDom)
		for _, token := range certificates {
			b.addCertificate(x509DataDom, token)
		}
	}

	b.KeyInfoDom = keyInfoElement
	return nil
}

// addSubjectAndCertificate creates the ds:X509SubjectName (optional) and ds:X509Certificate
// (mandatory) tags. Port of the private addSubjectAndCertificate.
func (b *AbstractSignatureBuilder) addSubjectAndCertificate(x509DataDom *xmldom.Node,
	token *model.CertificateToken) {
	xmlutils.DomUtilsAddTextElement(b.DocumentDom, x509DataDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementX509SubjectName, token.Subject().RFC2253())
	b.addCertificate(x509DataDom, token)
}

// addCertificate creates the mandatory ds:X509Certificate tag.
// Port of the private addCertificate.
func (b *AbstractSignatureBuilder) addCertificate(x509DataDom *xmldom.Node, token *model.CertificateToken) {
	xmlutils.DomUtilsAddTextElement(b.DocumentDom, x509DataDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementX509Certificate, utils.ToBase64(token.Encoded()))
}

// IncorporateObjects incorporates the ds:Object tags.
// Port of the protected #incorporateObjects.
func (b *AbstractSignatureBuilder) IncorporateObjects() error {
	if err := b.IncorporateQualifyingProperties(); err != nil {
		return err
	}
	if err := b.overrides.IncorporateSignedObjects(); err != nil {
		return err
	}
	return b.IncorporateCustomObjects()
}

// IncorporateQualifyingProperties incorporates the ds:Object holding the
// xades:QualifyingProperties element:
//
//	<ds:Object>
//	    <xades:QualifyingProperties>
//	        <xades:SignedProperties>...</xades:SignedProperties>
//	    </xades:QualifyingProperties>
//	</ds:Object>
//
// Port of the protected #incorporateQualifyingProperties.
func (b *AbstractSignatureBuilder) IncorporateQualifyingProperties() error {
	if utils.IsArrayNotEmpty(b.Params.SignedAdESObject()) {
		// Upstream logs "Incorporating signed XAdES Object from parameter" here.
		if xmlutils.DomUtilsIsDOMBytes(b.Params.SignedAdESObject()) {
			parsed, err := xmlutils.DomUtilsBuildDOMFromBytes(b.Params.SignedAdESObject())
			if err != nil {
				return err
			}
			signedObjectDom := b.DocumentDom.Import(parsed.DocumentElement(), true)
			b.SignatureDom.AppendChild(signedObjectDom)
		} else {
			return errors.New("The signed AdES Object shall represent an XML!")
		}
		return nil
	}

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	objectDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementObject)
	b.SignatureDom.AppendChild(objectDom)

	b.QualifyingPropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom, objectDom,
		b.overrides.XadesNamespace(), currentElements.ElementQualifyingProperties())
	xmlutils.DomUtilsAddNamespaceAttribute(b.QualifyingPropertiesDom, b.overrides.XadesNamespace())
	b.QualifyingPropertiesDom.SetAttr(xmldom.Name{Local: BuilderTarget},
		xmlutils.DomUtilsToElementReference(b.DeterministicId))

	return b.IncorporateSignedProperties()
}

// IncorporateSignedObjects incorporates the list of signed ds:Object elements (used for
// Enveloping packaging). By default only the objects enforced through the references are
// processed. Port of the protected #incorporateSignedObjects.
func (b *AbstractSignatureBuilder) IncorporateSignedObjects() error {
	// process only for enforced objects by default
	references := b.Params.References()
	for _, reference := range references {
		if reference.Object() != nil {
			if err := b.IncorporateObject(reference.Object()); err != nil {
				return err
			}
		}
	}
	return nil
}

// IncorporateCustomObjects incorporates a list of custom ds:Object elements within the
// ds:Signature element. Port of the protected #incorporateCustomObjects.
func (b *AbstractSignatureBuilder) IncorporateCustomObjects() error {
	if utils.IsCollectionNotEmpty(b.Params.Objects()) {
		for _, object := range b.Params.Objects() {
			if err := b.IncorporateObject(object); err != nil {
				return err
			}
		}
	}
	return nil
}

// IncorporateObject incorporates the given object within the ds:Signature.
// Port of the protected #incorporateObject.
func (b *AbstractSignatureBuilder) IncorporateObject(object *DSSObject) error {
	if object.Content() == nil {
		return errors.New("The content shall be defined inside DSSObject element! " +
			"Incorporation is not possible.")
	}

	// incorporate ds:Object dom
	objectDom := xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignatureDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementObject)

	// incorporate content
	if xmlutils.DomUtilsIsDOM(object.Content()) {
		contentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(object.Content())
		if err != nil {
			return err
		}
		xmlutils.DomUtilsAdoptChildren(objectDom, contentDom)
	} else {
		content, err := spi.DSSUtilsToByteArrayOfDocument(object.Content())
		if err != nil {
			return err
		}
		objectDom.AppendChild(xmldom.NewText(string(content)))
	}

	// incorporate Id attribute
	if utils.IsStringNotBlank(object.Id()) {
		objectDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, object.Id())
	}

	// incorporate MimeType attribute
	if utils.IsStringNotBlank(object.MimeType()) {
		objectDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeMIMEType.AttributeName()},
			object.MimeType())
	}

	// incorporate Encoding attribute
	if utils.IsStringNotBlank(object.EncodingAlgorithm()) {
		objectDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeEncoding.AttributeName()},
			object.EncodingAlgorithm())
	}

	return nil
}

// IncorporateReferenceSignedProperties incorporates the ds:Reference covering the
// xades:SignedProperties element:
//
//	<ds:Reference Type="http://uri.etsi.org/01903#SignedProperties" URI="#xades-id-A430...">
//	    <ds:Transforms>
//	        <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
//	    </ds:Transforms>
//	    <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//	    <ds:DigestValue>uijX/nvuu8g10ZVEklEnYatvFe8=</ds:DigestValue>
//	</ds:Reference>
//
// Port of the protected #incorporateReferenceSignedProperties.
func (b *AbstractSignatureBuilder) IncorporateReferenceSignedProperties() error {
	reference := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementReference)
	b.SignedInfoDom.AppendChild(reference)
	reference.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeType.AttributeName()},
		b.XadesPath.SignedPropertiesUri())
	reference.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeURI.AttributeName()},
		xmlutils.DomUtilsToElementReference(AbstractSignatureBuilderXAdESPrefix+b.DeterministicId))

	transforms := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementTransforms)
	reference.AppendChild(transforms)
	transform := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementTransform)
	transforms.AppendChild(transform)
	transform.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()},
		b.SignedPropertiesCanonicalizationMethod)

	digestAlgorithm := DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault(b.Params)
	DSSXMLUtilsIncorporateDigestMethod(reference, digestAlgorithm, b.overrides.XmldsigNamespace())

	// This is a workaround to ensure the canonicalization is performed successfully when the same
	// namespace prefix is used within the element
	if b.overrides.XmldsigNamespace().Prefix() != "" &&
		b.overrides.XmldsigNamespace().Prefix() == b.overrides.XadesNamespace().Prefix() {
		signedPropertiesDom, err := DSSXMLUtilsEnsureNamespacesDefinedWithQuery(b.DocumentDom, b.DeterministicId,
			b.XadesPath.SignedPropertiesPath())
		if err != nil {
			return err
		}
		b.SignedPropertiesDom = signedPropertiesDom
	}

	nodeToCanonicalize, err := b.NodeToCanonicalize(b.SignedPropertiesDom)
	if err != nil {
		return err
	}
	messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedNode(nodeToCanonicalize, digestAlgorithm,
		b.SignedPropertiesCanonicalizationMethod)
	if err != nil {
		return err
	}
	// Upstream traces the canonicalization method and the digest here.
	b.incorporateDigestValueOfReference(reference, messageDigest.Value())
	return nil
}

// IncorporateReferenceKeyInfo incorporates the ds:Reference covering the ds:KeyInfo element:
//
//	<ds:Reference URI="#keyInfo-id-A43023AFEB149830C242377CC941360F">
//	    <ds:Transforms>
//	        <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>
//	    </ds:Transforms>
//	    <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//	    <ds:DigestValue>uijX/nvuu2g10ZVEklEnYatvFe4=</ds:DigestValue>
//	</ds:Reference>
//
// Port of the protected #incorporateReferenceKeyInfo.
func (b *AbstractSignatureBuilder) IncorporateReferenceKeyInfo() error {
	if !b.Params.IsSignKeyInfo() {
		return nil
	}

	reference := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementReference)
	b.SignedInfoDom.AppendChild(reference)
	reference.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
		SignatureBuilderReferencePrefix+AbstractSignatureBuilderKeyInfoPrefix+b.DeterministicId)
	reference.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeURI.AttributeName()},
		xmlutils.DomUtilsToElementReference(AbstractSignatureBuilderKeyInfoPrefix+b.DeterministicId))

	transforms := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementTransforms)
	reference.AppendChild(transforms)
	transform := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementTransform)
	transforms.AppendChild(transform)
	transform.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()},
		b.KeyInfoCanonicalizationMethod)

	digestAlgorithm := DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault(b.Params)
	DSSXMLUtilsIncorporateDigestMethod(reference, digestAlgorithm, b.overrides.XmldsigNamespace())

	nodeToCanonicalize, err := b.NodeToCanonicalize(b.KeyInfoDom)
	if err != nil {
		return err
	}
	messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedNode(nodeToCanonicalize, digestAlgorithm,
		b.KeyInfoCanonicalizationMethod)
	if err != nil {
		return err
	}
	// Upstream traces the canonicalization method and the digest here.
	b.incorporateDigestValueOfReference(reference, messageDigest.Value())
	return nil
}

// incorporateDigestValueOfReference creates the ds:DigestValue DOM object for the given digest
// value computed on a canonicalized content.
// Port of the private incorporateDigestValueOfReference.
func (b *AbstractSignatureBuilder) incorporateDigestValueOfReference(referenceDom *xmldom.Node,
	digestValue []byte) {
	digestValueDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementDigestValue)
	base64EncodedDigestBytes := utils.ToBase64(digestValue)
	textNode := xmldom.NewText(base64EncodedDigestBytes)
	digestValueDom.AppendChild(textNode)
	referenceDom.AppendChild(digestValueDom)
}

// IncorporateSignatureValue incorporates the ds:SignatureValue element.
// Port of the protected #incorporateSignatureValue.
func (b *AbstractSignatureBuilder) IncorporateSignatureValue() {
	b.SignatureValueDom = xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
		common.XMLDSigElementSignatureValue)
	b.SignatureDom.AppendChild(b.SignatureValueDom)
	b.SignatureValueDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
		AbstractSignatureBuilderValuePrefix+b.DeterministicId)
}

// IncorporateSignedProperties creates the xades:SignedProperties DOM element:
//
//	<SignedProperties Id="xades-ide5c549340079fe19f3f90f03354a5965">
//
// Port of the protected #incorporateSignedProperties.
func (b *AbstractSignatureBuilder) IncorporateSignedProperties() error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	b.SignedPropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.QualifyingPropertiesDom,
		b.overrides.XadesNamespace(), currentElements.ElementSignedProperties())
	b.SignedPropertiesDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
		AbstractSignatureBuilderXAdESPrefix+b.DeterministicId)

	if err := b.IncorporateSignedSignatureProperties(); err != nil {
		return err
	}

	return b.incorporateSignedDataObjectProperties()
}

// IncorporateSignedSignatureProperties creates the xades:SignedSignatureProperties DOM element.
// Port of the protected #incorporateSignedSignatureProperties.
func (b *AbstractSignatureBuilder) IncorporateSignedSignatureProperties() error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	b.SignedSignaturePropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedPropertiesDom,
		b.overrides.XadesNamespace(), currentElements.ElementSignedSignatureProperties())

	if err := b.incorporateSigningTime(); err != nil {
		return err
	}

	if err := b.incorporateSigningCertificate(); err != nil {
		return err
	}

	if err := b.incorporatePolicy(); err != nil {
		return err
	}

	if err := b.incorporateSignatureProductionPlace(); err != nil {
		return err
	}

	return b.incorporateSignerRole()
}

// incorporatePolicy creates the xades:SignaturePolicyIdentifier DOM object.
// Port of the private incorporatePolicy.
func (b *AbstractSignatureBuilder) incorporatePolicy() error {
	signaturePolicy := b.Params.BLevel().SignaturePolicy()
	if signaturePolicy == nil {
		return nil
	}

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	signaturePolicyIdentifierDom := xmlutils.DomUtilsAddElement(b.DocumentDom,
		b.SignedSignaturePropertiesDom, b.overrides.XadesNamespace(),
		currentElements.ElementSignaturePolicyIdentifier())

	signaturePolicyId := signaturePolicy.Id()
	if utils.IsStringEmpty(signaturePolicyId) { // implicit
		xmlutils.DomUtilsAddElement(b.DocumentDom, signaturePolicyIdentifierDom,
			b.overrides.XadesNamespace(), currentElements.ElementSignaturePolicyImplied())
		return nil
	}

	// explicit
	signaturePolicyIdDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signaturePolicyIdentifierDom,
		b.overrides.XadesNamespace(), currentElements.ElementSignaturePolicyId())
	sigPolicyIdDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signaturePolicyIdDom,
		b.overrides.XadesNamespace(), currentElements.ElementSigPolicyId())

	identifierDom := xmlutils.DomUtilsAddTextElement(b.DocumentDom, sigPolicyIdDom,
		b.overrides.XadesNamespace(), currentElements.ElementIdentifier(), signaturePolicyId)

	qualifier := signaturePolicy.Qualifier()
	if qualifier != "" {
		identifierDom.SetAttr(
			xmldom.Name{Local: definition.XAdES132AttributeQualifier.AttributeName()}, qualifier.Value())
	}

	description := signaturePolicy.Description()
	if utils.IsStringNotEmpty(description) {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, sigPolicyIdDom, b.overrides.XadesNamespace(),
			currentElements.ElementDescription(), description)
	}

	documentationReferences := signaturePolicy.DocumentationReferences()
	if utils.IsArrayNotEmpty(documentationReferences) {
		if err := b.incorporateDocumentationReferences(sigPolicyIdDom, documentationReferences); err != nil {
			return err
		}
	}

	// Java: `if (signaturePolicy instanceof XmlPolicyWithTransforms)`; see the file header.
	if transforms, ok := SignatureBuilderPolicyTransforms(signaturePolicy); ok {
		DSSXMLUtilsIncorporateTransforms(signaturePolicyIdDom, transforms, b.overrides.XmldsigNamespace())
	}

	if signaturePolicy.DigestAlgorithm() != "" && signaturePolicy.DigestValue() != nil {
		sigPolicyHashDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signaturePolicyIdDom,
			b.overrides.XadesNamespace(), currentElements.ElementSigPolicyHash())

		digestAlgorithm := signaturePolicy.DigestAlgorithm()
		b.IncorporateDigestMethod(sigPolicyHashDom, digestAlgorithm)

		hashValue := signaturePolicy.DigestValue()
		base64EncodedHashValue := utils.ToBase64(hashValue)
		b.IncorporateDigestValue(sigPolicyHashDom, base64EncodedHashValue)
	}

	if signaturePolicy.IsSPQualifierPresent() {
		if err := b.incorporateSigPolicyQualifiers(signaturePolicyIdDom, signaturePolicy); err != nil {
			return err
		}
	}
	return nil
}

// incorporateSigPolicyQualifiers creates the xades:SigPolicyQualifiers DOM object.
// Port of the private incorporateSigPolicyQualifiers.
func (b *AbstractSignatureBuilder) incorporateSigPolicyQualifiers(signaturePolicyIdDom *xmldom.Node,
	signaturePolicy *model.Policy) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	sigPolicyQualifiers := xmlutils.DomUtilsAddElement(b.DocumentDom, signaturePolicyIdDom,
		b.overrides.XadesNamespace(), currentElements.ElementSigPolicyQualifiers())

	spUri := signaturePolicy.Spuri()
	if utils.IsStringNotEmpty(spUri) {
		sigPolicyQualifier := xmlutils.DomUtilsAddElement(b.DocumentDom, sigPolicyQualifiers,
			b.overrides.XadesNamespace(), currentElements.ElementSigPolicyQualifier())
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, sigPolicyQualifier, b.overrides.XadesNamespace(),
			currentElements.ElementSPURI(), spUri)
	}

	userNotice := signaturePolicy.UserNotice()
	if userNotice != nil && !userNotice.IsEmpty() {
		if err := spi.DSSUtilsAssertSPUserNoticeConfigurationValid(userNotice); err != nil {
			return err
		}

		sigPolicyQualifier := xmlutils.DomUtilsAddElement(b.DocumentDom, sigPolicyQualifiers,
			b.overrides.XadesNamespace(), currentElements.ElementSigPolicyQualifier())
		spUserNotice := xmlutils.DomUtilsAddElement(b.DocumentDom, sigPolicyQualifier,
			b.overrides.XadesNamespace(), currentElements.ElementSPUserNotice())

		organization := userNotice.Organization()
		noticeNumbers := userNotice.NoticeNumbers()
		if utils.IsStringNotEmpty(organization) && noticeNumbers != nil && len(noticeNumbers) > 0 {
			noticeRef := xmlutils.DomUtilsAddElement(b.DocumentDom, spUserNotice,
				b.overrides.XadesNamespace(), currentElements.ElementNoticeRef())
			xmlutils.DomUtilsAddTextElement(b.DocumentDom, noticeRef, b.overrides.XadesNamespace(),
				currentElements.ElementOrganization(), organization)

			noticeNumbersElement := xmlutils.DomUtilsAddElement(b.DocumentDom, noticeRef,
				b.overrides.XadesNamespace(), currentElements.ElementNoticeNumbers())
			for _, number := range noticeNumbers {
				xmlutils.DomUtilsAddTextElement(b.DocumentDom, noticeNumbersElement,
					b.overrides.XadesNamespace(), currentElements.Elementint(),
					fmt.Sprintf("%d", number))
			}
		}

		explicitText := userNotice.ExplicitText()
		if utils.IsStringNotEmpty(explicitText) {
			xmlutils.DomUtilsAddTextElement(b.DocumentDom, spUserNotice, b.overrides.XadesNamespace(),
				currentElements.ElementExplicitText(), explicitText)
		}
	}

	spDocSpecification := signaturePolicy.SpDocSpecification()
	if spDocSpecification != nil && utils.IsStringNotEmpty(spDocSpecification.Id()) {
		sigPolicyQualifier := xmlutils.DomUtilsAddElement(b.DocumentDom, sigPolicyQualifiers,
			b.overrides.XadesNamespace(), currentElements.ElementSigPolicyQualifier())
		if err := b.IncorporateSPDocSpecification(sigPolicyQualifier, spDocSpecification); err != nil {
			return err
		}
	}
	return nil
}

// incorporateSigningTime creates the xades:SigningTime DOM object element:
//
//	<SigningTime>2013-11-23T11:22:52Z</SigningTime>
//
// Port of the private incorporateSigningTime. DomUtils.createXMLGregorianCalendar(date) is
// immediately followed by toXMLFormat() upstream; the Go DomUtils folds the two into one call
// that returns the lexical form directly.
func (b *AbstractSignatureBuilder) incorporateSigningTime() error {
	var signingDate time.Time
	if sd := b.Params.BLevel().SigningDate(); sd != nil {
		signingDate = *sd
	}
	xmlSigningTime := xmlutils.DomUtilsCreateXMLGregorianCalendar(signingDate)

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	signingTimeDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XadesNamespace(),
		currentElements.ElementSigningTime())
	b.SignedSignaturePropertiesDom.AppendChild(signingTimeDom)
	textNode := xmldom.NewText(xmlSigningTime)
	signingTimeDom.AppendChild(textNode)
	return nil
}

// incorporateSigningCertificate creates the xades:SigningCertificate(V2) building block DOM
// object. Port of the private incorporateSigningCertificate.
func (b *AbstractSignatureBuilder) incorporateSigningCertificate() error {
	if b.Params.SigningCertificate() == nil && b.Params.GenerateTBSWithoutCertificate() {
		return nil
	}

	// Java collects the single signing certificate into a HashSet; a one-element slice is the
	// order-stable Go equivalent (the set never holds more than the signing certificate).
	certificates := []*model.CertificateToken{b.Params.SigningCertificate()}

	if b.Params.IsEn319132() {
		return b.incorporateSigningCertificateV2(certificates)
	}
	return b.incorporateSigningCertificateV1(certificates)
}

// incorporateSigningCertificateV1 ports the private incorporateSigningCertificateV1.
func (b *AbstractSignatureBuilder) incorporateSigningCertificateV1(certificates []*model.CertificateToken) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	signingCertificateDom := xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedSignaturePropertiesDom,
		b.overrides.XadesNamespace(), currentElements.ElementSigningCertificate())

	signingCertificateDigestMethod := b.Params.SigningCertificateDigestMethod()
	for _, certificate := range certificates {
		if _, err := b.IncorporateCert(signingCertificateDom, certificate, signingCertificateDigestMethod); err != nil {
			return err
		}
	}
	return nil
}

// incorporateSigningCertificateV2 ports the private incorporateSigningCertificateV2.
func (b *AbstractSignatureBuilder) incorporateSigningCertificateV2(certificates []*model.CertificateToken) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	signingCertificateDom := xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedSignaturePropertiesDom,
		b.overrides.XadesNamespace(), currentElements.ElementSigningCertificateV2())

	signingCertificateDigestMethod := b.Params.SigningCertificateDigestMethod()
	for _, certificate := range certificates {
		if _, err := b.IncorporateCert(signingCertificateDom, certificate, signingCertificateDigestMethod); err != nil {
			return err
		}
	}
	return nil
}

// incorporateSignedDataObjectProperties incorporates the xades:SignedDataObjectProperties DOM
// element. Port of the private incorporateSignedDataObjectProperties.
func (b *AbstractSignatureBuilder) incorporateSignedDataObjectProperties() error {
	if err := b.incorporateDataObjectFormat(); err != nil {
		return err
	}
	if err := b.incorporateCommitmentTypeIndications(); err != nil {
		return err
	}
	return b.incorporateContentTimestamps()
}

// signedDataObjectPropertiesDom returns (creating it on first use) the
// xades:SignedDataObjectProperties element. EN 319 132-1 clause 4.3.5: a XAdES signature shall
// not incorporate an empty SignedDataObjectProperties element.
// Port of the private getSignedDataObjectPropertiesDom.
func (b *AbstractSignatureBuilder) signedDataObjectPropertiesDom() (*xmldom.Node, error) {
	if b.SignedDataObjectPropertiesDom == nil {
		currentElements, err := b.CurrentXAdESElements()
		if err != nil {
			return nil, err
		}
		b.SignedDataObjectPropertiesDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedPropertiesDom,
			b.overrides.XadesNamespace(), currentElements.ElementSignedDataObjectProperties())
	}
	return b.SignedDataObjectPropertiesDom, nil
}

// incorporateDataObjectFormat incorporates the xades:DataObjectFormat DOM elements:
//
//	<DataObjectFormat ObjectReference="#detached-ref-id">
//	    <MimeType>text/plain</MimeType>
//	    ...
//	</DataObjectFormat>
//
// Port of the private incorporateDataObjectFormat.
func (b *AbstractSignatureBuilder) incorporateDataObjectFormat() error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	dataObjectFormats := b.Params.DataObjectFormatList()
	if dataObjectFormats == nil {
		dataObjectFormats = NewDataObjectFormatBuilder().SetReferences(b.Params.References()).Build()
		if b.Params.IsSignKeyInfo() {
			keyInfoDataObjectFormat := b.keyInfoDataObjectFormat()
			dataObjectFormats = append(dataObjectFormats, keyInfoDataObjectFormat)
		}
	}
	for _, dataObjectFormat := range dataObjectFormats {
		if err := xadesSignatureBuilderAssertDataObjectFormatValid(dataObjectFormat); err != nil {
			return err
		}

		signedDataObjectPropertiesDom, err := b.signedDataObjectPropertiesDom()
		if err != nil {
			return err
		}

		// create xades:DataObjectFormat
		dataObjectFormatDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signedDataObjectPropertiesDom,
			b.overrides.XadesNamespace(), currentElements.ElementDataObjectFormat())

		// add xades:DataObjectFormat/xades:Description
		if dataObjectFormat.Description() != "" {
			descriptionDom := xmlutils.DomUtilsAddElement(b.DocumentDom, dataObjectFormatDom,
				b.overrides.XadesNamespace(), currentElements.ElementDescription())
			xmlutils.DomUtilsSetTextNode(b.DocumentDom, descriptionDom, dataObjectFormat.Description())
		}

		// add xades:DataObjectFormat/xades:ObjectIdentifier
		if dataObjectFormat.ObjectIdentifier() != nil {
			objectIdentifierDom := xmlutils.DomUtilsAddElement(b.DocumentDom, dataObjectFormatDom,
				b.overrides.XadesNamespace(), currentElements.ElementObjectIdentifier())
			if err := b.incorporateObjectIdentifier(objectIdentifierDom, dataObjectFormat.ObjectIdentifier()); err != nil {
				return err
			}
		}

		// add xades:DataObjectFormat/xades:MimeType
		if dataObjectFormat.MimeType() != "" {
			mimeTypeDom := xmlutils.DomUtilsAddElement(b.DocumentDom, dataObjectFormatDom,
				b.overrides.XadesNamespace(), currentElements.ElementMimeType())
			xmlutils.DomUtilsSetTextNode(b.DocumentDom, mimeTypeDom, dataObjectFormat.MimeType())
		}

		// add xades:DataObjectFormat/xades:Encoding
		if dataObjectFormat.Encoding() != "" {
			encodingDom := xmlutils.DomUtilsAddElement(b.DocumentDom, dataObjectFormatDom,
				b.overrides.XadesNamespace(), currentElements.ElementEncoding())
			xmlutils.DomUtilsSetTextNode(b.DocumentDom, encodingDom, dataObjectFormat.Encoding())
		}

		// add xades:DataObjectFormat@ObjectReference
		if dataObjectFormat.ObjectReference() != "" {
			dataObjectFormatDom.SetAttr(
				xmldom.Name{Local: definition.XAdES132AttributeObjectReference.AttributeName()},
				dataObjectFormat.ObjectReference())
		}
	}
	return nil
}

// keyInfoDataObjectFormat ports the private getKeyInfoDataObjectFormat.
func (b *AbstractSignatureBuilder) keyInfoDataObjectFormat() *DSSDataObjectFormat {
	keyInfoDataObjectFormat := NewDSSDataObjectFormat()
	keyInfoDataObjectFormat.SetObjectReference(xmlutils.DomUtilsToElementReference(
		SignatureBuilderReferencePrefix + AbstractSignatureBuilderKeyInfoPrefix + b.DeterministicId))
	keyInfoDataObjectFormat.SetMimeType(enumerations.MimeTypeEnumXML.MimeTypeString())
	return keyInfoDataObjectFormat
}

// xadesSignatureBuilderAssertDataObjectFormatValid verifies that the DataObjectFormat conforms to
// the XAdES specification. Port of the private assertDataObjectFormatValid.
func xadesSignatureBuilderAssertDataObjectFormatValid(dataObjectFormat *DSSDataObjectFormat) error {
	if dataObjectFormat == nil {
		panic("DataObjectFormat cannot be null!")
	}
	if dataObjectFormat.Description() == "" && dataObjectFormat.ObjectIdentifier() == nil &&
		dataObjectFormat.MimeType() == "" {
		return errors.New("At least one of the Description, ObjectIdentifier or MimeType " +
			"shall be defined for a DataObjectFormat object!")
	}
	if dataObjectFormat.ObjectReference() == "" {
		return errors.New("ObjectReference attribute of DataObjectFormat shall be present!")
	}
	if !xmlutils.DomUtilsIsElementReference(dataObjectFormat.ObjectReference()) {
		return errors.New("ObjectReference attribute of DataObjectFormat " +
			"shall define a reference to an element within signature (i.e. shall begin with '#')!")
	}
	return nil
}

// incorporateContentTimestamps incorporates the content-timestamps within the signature being
// created. Port of the private incorporateContentTimestamps.
func (b *AbstractSignatureBuilder) incorporateContentTimestamps() error {
	contentTimestamps := b.Params.ContentTimestamps()
	if contentTimestamps == nil {
		return nil
	}

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	for _, contentTimestamp := range contentTimestamps {
		timeStampType := contentTimestamp.TimeStampType()
		var timestampDom *xmldom.Node
		switch {
		case enumerations.TimestampTypeAllDataObjectsTimestamp == timeStampType:
			signedDataObjectPropertiesDom, err := b.signedDataObjectPropertiesDom()
			if err != nil {
				return err
			}
			timestampDom = xmlutils.DomUtilsAddElement(b.DocumentDom, signedDataObjectPropertiesDom,
				b.overrides.XadesNamespace(), currentElements.ElementAllDataObjectsTimeStamp())
		case enumerations.TimestampTypeIndividualDataObjectsTimestamp == timeStampType:
			signedDataObjectPropertiesDom, err := b.signedDataObjectPropertiesDom()
			if err != nil {
				return err
			}
			timestampDom = xmlutils.DomUtilsAddElement(b.DocumentDom, signedDataObjectPropertiesDom,
				b.overrides.XadesNamespace(), currentElements.ElementIndividualDataObjectsTimeStamp())
		default:
			return errors.New(
				"Only types ALL_DATA_OBJECTS_TIMESTAMP and INDIVIDUAL_DATA_OBJECTS_TIMESTAMP are allowed")
		}
		if err := b.AddContentTimestamp(timestampDom, contentTimestamp); err != nil {
			return err
		}
	}
	return nil
}

// incorporateSignerRole incorporates the signer claimed roleType into the signed signature
// properties. Port of the private incorporateSignerRole.
func (b *AbstractSignatureBuilder) incorporateSignerRole() error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	claimedSignerRoles := b.Params.BLevel().ClaimedSignerRoles()
	signedAssertions := b.Params.BLevel().SignedAssertions()

	var signerRoleDom *xmldom.Node
	if claimedSignerRoles != nil {
		if b.Params.IsEn319132() {
			signerRoleDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedSignaturePropertiesDom,
				b.overrides.XadesNamespace(), currentElements.ElementSignerRoleV2())
		} else {
			signerRoleDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedSignaturePropertiesDom,
				b.overrides.XadesNamespace(), currentElements.ElementSignerRole())
		}

		if utils.IsCollectionNotEmpty(claimedSignerRoles) {
			claimedRolesDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signerRoleDom,
				b.overrides.XadesNamespace(), currentElements.ElementClaimedRoles())
			b.addRoles(claimedSignerRoles, claimedRolesDom, currentElements.ElementClaimedRole())
		}
	}
	if signedAssertions != nil && b.Params.IsEn319132() {
		if signerRoleDom == nil {
			signerRoleDom = xmlutils.DomUtilsAddElement(b.DocumentDom, b.SignedSignaturePropertiesDom,
				b.overrides.XadesNamespace(), currentElements.ElementSignerRoleV2())
		}

		if utils.IsCollectionNotEmpty(signedAssertions) {
			signedAssertionsDom := xmlutils.DomUtilsAddElement(b.DocumentDom, signerRoleDom,
				b.overrides.XadesNamespace(), currentElements.ElementSignedAssertions())
			if err := b.addAssertions(signedAssertions, signedAssertionsDom); err != nil {
				return err
			}
		}
	}
	return nil
}

// addRoles ports the private addRoles.
func (b *AbstractSignatureBuilder) addRoles(signerRoles []string, rolesDom *xmldom.Node,
	roleType common.DSSElement) {
	for _, signerRole := range signerRoles {
		roleDom := xmlutils.DomUtilsAddElement(b.DocumentDom, rolesDom, b.overrides.XadesNamespace(), roleType)
		xmlutils.DomUtilsSetTextNode(b.DocumentDom, roleDom, signerRole)
	}
}

// incorporateSignatureProductionPlace ports the private incorporateSignatureProductionPlace.
func (b *AbstractSignatureBuilder) incorporateSignatureProductionPlace() error {
	signatureProductionPlace := b.Params.BLevel().SignerLocation()
	if signatureProductionPlace == nil || signatureProductionPlace.IsEmpty() {
		return nil
	}

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	var signatureProductionPlaceDom *xmldom.Node
	if b.Params.IsEn319132() {
		signatureProductionPlaceDom = xmlutils.DomUtilsAddElement(b.DocumentDom,
			b.SignedSignaturePropertiesDom, b.overrides.XadesNamespace(),
			currentElements.ElementSignatureProductionPlaceV2())
	} else {
		signatureProductionPlaceDom = xmlutils.DomUtilsAddElement(b.DocumentDom,
			b.SignedSignaturePropertiesDom, b.overrides.XadesNamespace(),
			currentElements.ElementSignatureProductionPlace())
	}

	city := signatureProductionPlace.Locality()
	if city != "" {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, signatureProductionPlaceDom,
			b.overrides.XadesNamespace(), currentElements.ElementCity(), city)
	}

	if b.Params.IsEn319132() {
		streetAddress := signatureProductionPlace.StreetAddress()
		if streetAddress != "" {
			xmlutils.DomUtilsAddTextElement(b.DocumentDom, signatureProductionPlaceDom,
				b.overrides.XadesNamespace(), currentElements.ElementStreetAddress(), streetAddress)
		}
	}

	stateOrProvince := signatureProductionPlace.StateOrProvince()
	if stateOrProvince != "" {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, signatureProductionPlaceDom,
			b.overrides.XadesNamespace(), currentElements.ElementStateOrProvince(), stateOrProvince)
	}

	postalCode := signatureProductionPlace.PostalCode()
	if postalCode != "" {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, signatureProductionPlaceDom,
			b.overrides.XadesNamespace(), currentElements.ElementPostalCode(), postalCode)
	}

	country := signatureProductionPlace.Country()
	if country != "" {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, signatureProductionPlaceDom,
			b.overrides.XadesNamespace(), currentElements.ElementCountryName(), country)
	}
	return nil
}

// incorporateCommitmentTypeIndications creates the xades:CommitmentTypeIndication elements.
//
//	<xsd:element name="CommitmentTypeIndication" type="CommitmentTypeIndicationType"/>
//	<xsd:complexType name="CommitmentTypeIndicationType">
//	    <xsd:sequence>
//	        <xsd:element name="CommitmentTypeId" type="ObjectIdentifierType"/>
//	        <xsd:choice>
//	            <xsd:element name="ObjectReference" type="xsd:anyURI" maxOccurs="unbounded"/>
//	            <xsd:element name="AllSignedDataObjects"/>
//	        </xsd:choice>
//	        <xsd:element name="CommitmentTypeQualifiers" type="CommitmentTypeQualifiersListType" minOccurs="0"/>
//	    </xsd:sequence>
//	</xsd:complexType>
//
// Port of the private incorporateCommitmentTypeIndications.
func (b *AbstractSignatureBuilder) incorporateCommitmentTypeIndications() error {
	commitmentTypeIndications := b.Params.BLevel().CommitmentTypeIndications()
	if !utils.IsCollectionNotEmpty(commitmentTypeIndications) {
		return nil
	}

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	for _, commitmentTypeIndication := range commitmentTypeIndications {
		if err := xadesSignatureBuilderAssertCommitmentTypeNotNull(commitmentTypeIndication); err != nil {
			return err
		}

		signedDataObjectPropertiesDom, err := b.signedDataObjectPropertiesDom()
		if err != nil {
			return err
		}

		commitmentTypeIndicationDom := xmlutils.DomUtilsAddElement(b.DocumentDom,
			signedDataObjectPropertiesDom, b.overrides.XadesNamespace(),
			currentElements.ElementCommitmentTypeIndication())

		commitmentTypeIdDom := xmlutils.DomUtilsAddElement(b.DocumentDom, commitmentTypeIndicationDom,
			b.overrides.XadesNamespace(), currentElements.ElementCommitmentTypeId())
		if err := b.incorporateObjectIdentifier(commitmentTypeIdDom, commitmentTypeIndication); err != nil {
			return err
		}

		var signedDataObjects []string
		var commitmentTypeQualifiers []*model.CommitmentQualifier
		if commonCommitmentType, ok := commitmentTypeIndication.(*model.CommonCommitmentType); ok {
			signedDataObjects = commonCommitmentType.SignedDataObjects()
			commitmentTypeQualifiers = commonCommitmentType.CommitmentTypeQualifiers()
		}

		if utils.IsArrayNotEmpty(signedDataObjects) {
			// add xades:ObjectReference
			for _, signedDataObjectUri := range signedDataObjects {
				if utils.IsStringBlank(signedDataObjectUri) {
					return errors.New("SignedDataObject URI cannot be null!")
				}
				signedDataObjectUri = xmlutils.DomUtilsToElementReference(signedDataObjectUri)
				xmlutils.DomUtilsAddTextElement(b.DocumentDom, commitmentTypeIndicationDom,
					b.overrides.XadesNamespace(), currentElements.ElementObjectReference(),
					signedDataObjectUri)
			}
		} else {
			// add xades:AllSignedDataObjects
			xmlutils.DomUtilsAddElement(b.DocumentDom, commitmentTypeIndicationDom,
				b.overrides.XadesNamespace(), currentElements.ElementAllSignedDataObjects())
		}

		// add xades:CommitmentTypeQualifiers
		if utils.IsArrayNotEmpty(commitmentTypeQualifiers) {
			commitmentTypeQualifiersElement := xmlutils.DomUtilsAddElement(b.DocumentDom,
				commitmentTypeIndicationDom, b.overrides.XadesNamespace(),
				currentElements.ElementCommitmentTypeQualifiers())

			for _, commitmentQualifier := range commitmentTypeQualifiers {
				if commitmentQualifier == nil {
					panic("CommitmentTypeQualifier cannot be null!")
				}
				content := commitmentQualifier.Content()
				if content == nil {
					return errors.New("CommitmentTypeQualifier content cannot be null!")
				}

				commitmentTypeQualifierElement := xmlutils.DomUtilsAddElement(b.DocumentDom,
					commitmentTypeQualifiersElement, b.overrides.XadesNamespace(),
					currentElements.ElementCommitmentTypeQualifier())
				// incorporate content
				var objectContentDom *xmldom.Node
				if xmlutils.DomUtilsIsDOM(content) {
					parsed, err := xmlutils.DomUtilsBuildDOMFromDocument(content)
					if err != nil {
						return err
					}
					objectContentDom = b.DocumentDom.Import(parsed.DocumentElement(), true)
				} else {
					// Upstream logs "None XML encoded CommitmentTypeQualifier has been provided.
					// Incorporate as text node." here.
					raw, err := spi.DSSUtilsToByteArrayOfDocument(content)
					if err != nil {
						return err
					}
					objectContentDom = xmldom.NewText(string(raw))
				}
				commitmentTypeQualifierElement.AppendChild(objectContentDom)
			}
		}
	}
	return nil
}

// xadesSignatureBuilderAssertCommitmentTypeNotNull verifies that the commitment type is not nil
// and contains at least one of the mandatory elements: URI or OID.
// Port of the private assertCommitmentTypeNotNull.
func xadesSignatureBuilderAssertCommitmentTypeNotNull(commitmentType enumerations.CommitmentType) error {
	if commitmentType == nil {
		panic("CommitmentType cannot be null!")
	}
	if commitmentType.URI() == "" && commitmentType.OID() == "" {
		return errors.New(
			"The URI or OID must be defined for commitmentTypeIndication for XAdES creation!")
	}
	return nil
}

// incorporateObjectIdentifier creates the xades:ObjectIdentifierType DOM object:
//
//	<xsd:complexType name="ObjectIdentifierType">
//	    <xsd:sequence>
//	        <xsd:element name="Identifier" type="IdentifierType"/>
//	        <xsd:element name="Description" type="xsd:string" minOccurs="0"/>
//	        <xsd:element name="DocumentationReferences" type="DocumentationReferencesType" minOccurs="0"/>
//	    </xsd:sequence>
//	</xsd:complexType>
//
// Port of the private incorporateObjectIdentifier.
func (b *AbstractSignatureBuilder) incorporateObjectIdentifier(parentDom *xmldom.Node,
	objectIdentifier enumerations.ObjectIdentifier) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	// add xades:Identifier
	if err := b.incorporateIdentifier(parentDom, objectIdentifier); err != nil {
		return err
	}

	// add xades:Description
	description := objectIdentifier.Description()
	if description != "" {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, parentDom, b.overrides.XadesNamespace(),
			currentElements.ElementDescription(), description)
	}

	// add xades:DocumentationReferences
	documentationReferences := objectIdentifier.DocumentationReferences()
	if utils.IsArrayNotEmpty(documentationReferences) {
		return b.incorporateDocumentationReferences(parentDom, documentationReferences)
	}
	return nil
}

// incorporateIdentifier creates the xades:Identifier DOM object of an xades:ObjectIdentifierType.
// Port of the private incorporateIdentifier.
func (b *AbstractSignatureBuilder) incorporateIdentifier(parentDom *xmldom.Node,
	objectIdentifier enumerations.ObjectIdentifier) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	uri := objectIdentifier.URI()
	oid := objectIdentifier.OID()
	qualifier := objectIdentifier.Qualifier()
	if utils.IsStringEmpty(uri) {
		if utils.IsStringEmpty(oid) {
			return errors.New(
				"The URI or OID must be defined for XAdES IdentifierType element!")
		}
		if qualifier == "" {
			return errors.New("When using OID as object identifier in XAdES, " +
				"a Qualifier shall be provided! See EN 319 132-1 for more details.")
		}

		switch qualifier {
		case enumerations.ObjectIdentifierQualifierOIDAsURI:
			if spi.DSSUtilsIsUrnOid(oid) {
				return fmt.Errorf(
					"Qualifier '%s' shall not be used for URN encoded OID! "+
						"See EN 319 132-1 for more details.", qualifier)
			}

		case enumerations.ObjectIdentifierQualifierOIDAsURN:
			if !spi.DSSUtilsIsUrnOid(oid) {
				oid = spi.DSSUtilsToUrnOid(oid)
			}

		default:
			return fmt.Errorf(
				"The Qualifier '%s' is not supported!", qualifier)
		}
		uri = oid

	} else {
		if qualifier != "" {
			return errors.New("When using URI as object identifier in XAdES, " +
				"a Qualifier shall not be present! See EN 319 132-1 for more details.")
		}
	}

	// add xades:Identifier
	identifierDom := xmlutils.DomUtilsAddTextElement(b.DocumentDom, parentDom,
		b.overrides.XadesNamespace(), currentElements.ElementIdentifier(), uri)

	// add xades:Identifier@Qualifier
	if qualifier != "" {
		identifierDom.SetAttr(
			xmldom.Name{Local: definition.XAdES132AttributeQualifier.AttributeName()}, qualifier.Value())
	}
	return nil
}

// incorporateDocumentationReferences ports the private incorporateDocumentationReferences.
func (b *AbstractSignatureBuilder) incorporateDocumentationReferences(parentElement *xmldom.Node,
	documentationReferences []string) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	documentReferencesDom := xmlutils.DomUtilsAddElement(b.DocumentDom, parentElement,
		b.overrides.XadesNamespace(), currentElements.ElementDocumentationReferences())
	for _, ref := range documentationReferences {
		xmlutils.DomUtilsAddTextElement(b.DocumentDom, documentReferencesDom, b.overrides.XadesNamespace(),
			currentElements.ElementDocumentationReference(), ref)
	}
	return nil
}

// SignDocument adds the signature value to the signature and returns the XML signature.
// Port of the overridden #signDocument(byte[]).
func (b *AbstractSignatureBuilder) SignDocument(signatureValue []byte) (model.DSSDocument, error) {
	if !b.Built {
		if _, err := b.Build(); err != nil {
			return nil, err
		}
	}

	encryptionAlgorithm := b.Params.EncryptionAlgorithm()
	signatureValueBytes, err := spi.DSSASN1UtilsEnsurePlainSignatureValue(encryptionAlgorithm, signatureValue)
	if err != nil {
		return nil, err
	}
	signatureValueBase64Encoded := utils.ToBase64(signatureValueBytes)
	signatureValueNode := xmldom.NewText(signatureValueBase64Encoded)
	b.SignatureValueDom.AppendChild(signatureValueNode)
	return b.CreateXmlDocument()
}

// AddContentTimestamp adds the content of a timestamp into a given timestamp element.
// Port of the protected #addContentTimestamp.
func (b *AbstractSignatureBuilder) AddContentTimestamp(timestampElement *xmldom.Node,
	token *validation.TimestampToken) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}

	// add includes: URI + referencedData = "true"
	// add canonicalizationMethod: Algorithm
	// add encapsulatedTimestamp: Encoding, Id, while its textContent is the base64 encoding of
	// the data to digest
	includes := token.TimestampIncludes()
	if includes != nil {
		for _, include := range includes {
			timestampIncludeElement := xmlutils.DomUtilsCreateElementNS(b.DocumentDom,
				b.overrides.XadesNamespace(), currentElements.ElementInclude())
			uri := xmlutils.DomUtilsToElementReference(include.URI())
			timestampIncludeElement.SetAttr(xmldom.Name{Local: BuilderURI}, uri)
			timestampIncludeElement.SetAttr(xmldom.Name{Local: BuilderReferencedData}, "true")
			timestampElement.AppendChild(timestampIncludeElement)
		}
	}

	canonicalizationMethod := token.CanonicalizationMethod()
	if utils.IsStringNotEmpty(canonicalizationMethod) {
		canonicalizationMethodElement := xmlutils.DomUtilsCreateElementNS(b.DocumentDom,
			b.overrides.XmldsigNamespace(), common.XMLDSigElementCanonicalizationMethod)
		canonicalizationMethodElement.SetAttr(
			xmldom.Name{Local: common.XMLDSigAttributeAlgorithm.AttributeName()}, canonicalizationMethod)
		timestampElement.AppendChild(canonicalizationMethodElement)
	} else {
		return errors.New("Unable to create a timestamp with empty canonicalization method. " +
			"See EN 319 132-1: 4.5 Managing canonicalization of XML nodesets.")
	}

	encapsulatedTimestampElement := xmlutils.DomUtilsCreateElementNS(b.DocumentDom,
		b.overrides.XadesNamespace(), currentElements.ElementEncapsulatedTimeStamp())
	encapsulatedTimestampElement.SetTextContent(utils.ToBase64(token.Encoded()))
	timestampElement.AppendChild(encapsulatedTimestampElement)

	// Build Id after time-stamp incorporation to ensure timestampElement contains a new time-stamp
	attributeIdentifier := AttributeIdentifierBuild(timestampElement)
	xmlIdentifier, err := b.ToXmlIdentifier(attributeIdentifier)
	if err != nil {
		return err
	}
	timestampId := xadesBuilderTimestampPrefix + xmlIdentifier
	timestampElement.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, timestampId)
	encapsulatedTimestampElement.SetAttr(
		xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
		xadesBuilderEncapsulatedTimestampPrefix+timestampId)
	return nil
}

// NodeToCanonicalize returns the node to be canonicalized, applying indents if required.
// Port of the protected #getNodeToCanonicalize.
func (b *AbstractSignatureBuilder) NodeToCanonicalize(node *xmldom.Node) (*xmldom.Node, error) {
	if b.Params.IsPrettyPrint() {
		return DSSXMLUtilsGetIndentedNode(b.DocumentDom, node)
	}
	return node, nil
}

// AlignNodes aligns children indents. Port of the overridden protected #alignNodes.
func (b *AbstractSignatureBuilder) AlignNodes() {
	if b.UnsignedSignaturePropertiesDom != nil {
		DSSXMLUtilsAlignChildrenIndents(b.UnsignedSignaturePropertiesDom)
	}
	if b.QualifyingPropertiesDom != nil {
		DSSXMLUtilsAlignChildrenIndents(b.QualifyingPropertiesDom)
	}
}

// addAssertions ports the private addAssertions.
func (b *AbstractSignatureBuilder) addAssertions(signedAssertions []string, rolesDom *xmldom.Node) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	for _, signedAssertion := range signedAssertions {
		roleDom := xmlutils.DomUtilsAddElement(b.DocumentDom, rolesDom, b.overrides.XadesNamespace(),
			currentElements.ElementSignedAssertion())
		samlAssertion, err := xmlutils.DomUtilsBuildDOMFromString(signedAssertion)
		if err != nil {
			return err
		}
		docEl := samlAssertion.DocumentElement()
		node := b.DocumentDom.Import(docEl, true)
		roleDom.AppendChild(node)
	}
	return nil
}

// xadesSignatureBuilderPolicyTransformsRegistry backs
// SignatureBuilderRegisterPolicyTransforms / SignatureBuilderPolicyTransforms. See the
// file header: it stands in for Java's `instanceof XmlPolicyWithTransforms`, which Go cannot
// express because model.BLevelParameters hands out a *model.Policy and XmlPolicyWithTransforms
// embeds model.Policy by value.
var xadesSignatureBuilderPolicyTransformsRegistry sync.Map // map[*model.Policy][]DSSTransform

// SignatureBuilderRegisterPolicyTransforms records policy as an XmlPolicyWithTransforms and
// returns the *model.Policy to hand to BLevelParameters.SetSignaturePolicy, so that
// incorporatePolicy writes its ds:Transforms exactly where Java's
// `signaturePolicy instanceof XmlPolicyWithTransforms` branch does.
//
//	xmlPolicy := xades.NewXmlPolicyWithTransforms()
//	xmlPolicy.SetTransforms(transforms)
//	params.BLevel().SetSignaturePolicy(xades.SignatureBuilderRegisterPolicyTransforms(xmlPolicy))
func SignatureBuilderRegisterPolicyTransforms(policy *XmlPolicyWithTransforms) *model.Policy {
	xadesSignatureBuilderPolicyTransformsRegistry.Store(&policy.Policy, policy.Transforms())
	return &policy.Policy
}

// SignatureBuilderPolicyTransforms resolves the ds:Transforms registered for policy, and
// reports whether policy is the model.Policy of an XmlPolicyWithTransforms.
func SignatureBuilderPolicyTransforms(policy *model.Policy) ([]DSSTransform, bool) {
	value, ok := xadesSignatureBuilderPolicyTransformsRegistry.Load(policy)
	if !ok {
		return nil, false
	}
	transforms, ok := value.([]DSSTransform)
	return transforms, ok
}

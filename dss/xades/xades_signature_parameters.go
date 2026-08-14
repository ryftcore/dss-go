// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/XAdESSignatureParameters.java (DSS 6.5.RC1).
//
// # GetContext / GetDeterministicId / DetachedContents / Reinit shadowing (DEVIATION, flagged
// # for the integrator)
//
// Java's getContext() override lazily instantiates and returns an XAdESSignatureParameters
// through the SAME protected `context` field AbstractSignatureParameters declares (statically
// typed ProfileParameters there, but the concrete runtime object is the XAdESProfileParameters
// this override created and assigned) - so every base method that reads `context` (getContext()
// itself, getDeterministicId(), getDetachedContents(), reinit()) transparently sees the XAdES
// object once this override has run.
//
// Go has neither field-level covariance nor virtual dispatch across embedding: the private
// `context` field document.AbstractSignatureParameters declares belongs to package document and
// cannot be written from here even if its static type were changed. This type therefore keeps
// its own separate `context *XAdESProfileParameters` field (shadowing the embedded field's NAME
// only, not its storage) and shadows every base method whose Java behaviour depends on that
// field being the XAdES one: GetContext (returns *XAdESProfileParameters, the type every XAdES
// call site in this package already assumes per S4D_BRIEF.md's forward-dependency comments),
// GetDeterministicId (re-implemented against this shadowed GetContext so
// XAdESCounterSignatureParameters's own override - see xades_counter_signature_parameters.go's
// "Deterministic Id priming" section - and every base-typed reader agree on one cached value),
// DetachedContents (checks the shadowed context first, exactly like the base, then falls back to
// the promoted AbstractSignatureParameters.DetachedContents() for the base's own instance-field
// fallback path - a path nothing in this package ever populates through the base's own context,
// since every XAdES call site sets detached contents via GetContext().SetDetachedContents, per
// xades_service.go), and Reinit (clears this field in addition to delegating to the promoted
// Reinit for symmetry, though nothing reads the base's own separate context afterwards).
//
// java.io.Serializable and serialVersionUID are dropped; hashCode() has no Go counterpart.
package xades

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XPathElementPlacement defines ways to embed a signature (used for ENVELOPED format), the
// XAdESSignatureParameters.XPathElementPlacement Java enum.
type XPathElementPlacement string

const (
	// XPathElementPlacement_XPathAfter inserts the signature after the element referenced by
	// XPath.
	XPathElementPlacement_XPathAfter XPathElementPlacement = "XPathAfter"

	// XPathElementPlacement_XPathFirstChildOf inserts the signature as the first child of the
	// element referenced by XPath.
	XPathElementPlacement_XPathFirstChildOf XPathElementPlacement = "XPathFirstChildOf"
)

// XAdESSignatureParameters defines SignatureParameters to deal with XAdES signature
// creation/extension.
type XAdESSignatureParameters struct {
	document.AbstractSignatureParameters[*XAdESTimestampParameters]

	// context is the XAdES-specific signature creation context; see the file header for why
	// this shadows, rather than reuses, the embedded base's own private context field.
	context *XAdESProfileParameters

	// addX509SubjectName allows adding an optional X509SubjectName in the tag X509Data.
	addX509SubjectName bool

	// dssReferences is a list of references to incorporate.
	dssReferences []*DSSReference

	// embedXML: in case of ENVELOPING signature, allows including the complete XML and not its
	// base64 encoded value. NOTE: not compatible with Base64Transform.
	embedXML bool

	// en319132 defines if the signature shall be created according to ETSI EN 319 132.
	// Default: true.
	en319132 bool

	// keyInfoCanonicalizationMethod: ds:CanonicalizationMethod indicates the canonicalization
	// algorithm: Algorithm="..." for KeyInfo. The EXCLUSIVE canonicalization is used by
	// default.
	keyInfoCanonicalizationMethod string

	// signedInfoCanonicalizationMethod: ds:CanonicalizationMethod indicates the
	// canonicalization algorithm: Algorithm="..." for SignedInfo. The EXCLUSIVE
	// canonicalization is used by default.
	signedInfoCanonicalizationMethod string

	// signedPropertiesCanonicalizationMethod: ds:CanonicalizationMethod indicates the
	// canonicalization algorithm: Algorithm="..." for SignedProperties. The EXCLUSIVE
	// canonicalization is used by default.
	signedPropertiesCanonicalizationMethod string

	// manifestSignature allows producing a Manifest signature
	// (https://www.w3.org/TR/xmldsig-core/#sec-o-Manifest).
	manifestSignature bool

	// rootDocument defines the root element of the file to create signature in (used in
	// INTERNALLY_DETACHED).
	rootDocument *xmldom.Node

	// signedAdESObject is an optional parameter that contains the canonicalized XML of the
	// XAdES object that was digested, referenced from the SigningInfo, and indirectly signed
	// when the signature value was created. If this parameter is specified it will be used in
	// the signed XML document.
	signedAdESObject []byte

	// signingCertificateDigestMethod is the digest method used to create the digest of the
	// signer's certificate.
	signingCertificateDigestMethod enumerations.DigestAlgorithm

	// signKeyInfo optionally defines whether the "KeyInfo" element should be signed. If the
	// value of parameter is TRUE, reference of the "KeyInfo" element will be added to
	// "SignedInfo". FALSE by default.
	signKeyInfo bool

	// xPathLocationString defines the signature position xpath string (used for ENVELOPED
	// format).
	xPathLocationString string

	// xPathElementPlacement defines the signature placement relatively to the element defined
	// in xPathLocationString (used for ENVELOPED format).
	xPathElementPlacement XPathElementPlacement

	// prettyPrint: if true, prints each signature's tag to a new line with a relevant indent.
	prettyPrint bool

	// xmldsigNamespace is the XMLDSig definition.
	xmldsigNamespace *common.DSSNamespace

	// xadesNamespace is the XAdES 1.1.1, 1.2.2 or 1.3.2 definition.
	xadesNamespace *common.DSSNamespace

	// xades141Namespace is the XAdES 1.4.1 definition.
	xades141Namespace *common.DSSNamespace

	// objects is a list of custom ds:Object elements to be incorporated inside the signature.
	objects []*DSSObject

	// tokenReferencesDigestAlgorithm is the DigestAlgorithm used to incorporate
	// CompleteCertificateRefs/CompleteRevocationRefs on -C level. Default: SHA512.
	tokenReferencesDigestAlgorithm enumerations.DigestAlgorithm

	// dataObjectFormatList is a list of custom xades:DataObjectFormat elements incorporated
	// within xades:SignedDataObjectProperties element of the signature.
	dataObjectFormatList []*DSSDataObjectFormat
}

// NewXAdESSignatureParameters instantiates the object with null values, including the
// en319132, keyInfoCanonicalizationMethod, signedInfoCanonicalizationMethod,
// signedPropertiesCanonicalizationMethod, signingCertificateDigestMethod, xmldsigNamespace,
// xadesNamespace, xades141Namespace and tokenReferencesDigestAlgorithm field initializers. Port
// of the default constructor.
func NewXAdESSignatureParameters() *XAdESSignatureParameters {
	return &XAdESSignatureParameters{
		AbstractSignatureParameters:            document.NewAbstractSignatureParameters[*XAdESTimestampParameters](),
		en319132:                               true,
		keyInfoCanonicalizationMethod:          xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
		signedInfoCanonicalizationMethod:       xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
		signedPropertiesCanonicalizationMethod: xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
		signingCertificateDigestMethod:         enumerations.DigestAlgorithm_SHA512,
		xmldsigNamespace:                       common.XMLDSigNS,
		xadesNamespace:                         common.NewDSSNamespace(definition.XAdESNamespace_XADES_132.Uri(), "xades"),
		xades141Namespace:                      definition.XAdESNamespace_XADES_141,
		tokenReferencesDigestAlgorithm:         enumerations.DigestAlgorithm_SHA512,
	}
}

// SetSignatureLevel overrides AbstractSerializableSignatureParameters#setSignatureLevel,
// restricting the value to the XAdES form. Panics with the Java message when signatureLevel is
// empty or not a XAdES level (IllegalArgumentException upstream).
func (p *XAdESSignatureParameters) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) {
	form, err := signatureLevel.SignatureForm()
	if signatureLevel == "" || err != nil || enumerations.SignatureForm_XAdES != form {
		panic("Only XAdES form is allowed !")
	}
	p.AbstractSignatureParameters.SetSignatureLevel(signatureLevel)
}

// SigningCertificateDigestMethod sees SetSigningCertificateDigestMethod. Ports
// getSigningCertificateDigestMethod().
func (p *XAdESSignatureParameters) SigningCertificateDigestMethod() enumerations.DigestAlgorithm {
	return p.signingCertificateDigestMethod
}

// SetSigningCertificateDigestMethod is a part of the standard: 7.2.2 The SigningCertificate
// element (101 903 V1.4.2 (2010-12) XAdES). The digest method indicates the digest algorithm to
// be used to calculate the CertDigest element that contains the digest for each certificate
// referenced in the sequence. Default: SHA512 (DigestAlgorithm.SHA512). Panics with the Java
// message when signingCertificateDigestMethod is empty (Objects.requireNonNull upstream). Ports
// setSigningCertificateDigestMethod(DigestAlgorithm).
func (p *XAdESSignatureParameters) SetSigningCertificateDigestMethod(signingCertificateDigestMethod enumerations.DigestAlgorithm) {
	if signingCertificateDigestMethod == "" {
		panic("SigningCertificateDigestMethod cannot be null!")
	}
	p.signingCertificateDigestMethod = signingCertificateDigestMethod
}

// SignedInfoCanonicalizationMethod gets the SignedInfo canonicalization algorithm. Ports
// getSignedInfoCanonicalizationMethod().
func (p *XAdESSignatureParameters) SignedInfoCanonicalizationMethod() string {
	return p.signedInfoCanonicalizationMethod
}

// SetSignedInfoCanonicalizationMethod sets the canonicalization algorithm to be used when
// dealing with SignedInfo. Ports setSignedInfoCanonicalizationMethod(String).
func (p *XAdESSignatureParameters) SetSignedInfoCanonicalizationMethod(signedInfoCanonicalizationMethod string) {
	xadesSignatureParametersAssertCanonicalizationNotEmpty(signedInfoCanonicalizationMethod)
	p.signedInfoCanonicalizationMethod = signedInfoCanonicalizationMethod
}

// SignedPropertiesCanonicalizationMethod gets the SignedProperties canonicalization algorithm.
// Ports getSignedPropertiesCanonicalizationMethod().
func (p *XAdESSignatureParameters) SignedPropertiesCanonicalizationMethod() string {
	return p.signedPropertiesCanonicalizationMethod
}

// SetSignedPropertiesCanonicalizationMethod sets the canonicalization algorithm to be used when
// dealing with SignedProperties. Ports setSignedPropertiesCanonicalizationMethod(String).
func (p *XAdESSignatureParameters) SetSignedPropertiesCanonicalizationMethod(signedPropertiesCanonicalizationMethod string) {
	xadesSignatureParametersAssertCanonicalizationNotEmpty(signedPropertiesCanonicalizationMethod)
	p.signedPropertiesCanonicalizationMethod = signedPropertiesCanonicalizationMethod
}

// KeyInfoCanonicalizationMethod returns the canonicalization algorithm used for dealing with
// KeyInfo. Ports getKeyInfoCanonicalizationMethod().
func (p *XAdESSignatureParameters) KeyInfoCanonicalizationMethod() string {
	return p.keyInfoCanonicalizationMethod
}

// SetKeyInfoCanonicalizationMethod sets the canonicalization algorithm used for dealing with
// KeyInfo. Ports setKeyInfoCanonicalizationMethod(String).
func (p *XAdESSignatureParameters) SetKeyInfoCanonicalizationMethod(keyInfoCanonicalizationMethod string) {
	xadesSignatureParametersAssertCanonicalizationNotEmpty(keyInfoCanonicalizationMethod)
	p.keyInfoCanonicalizationMethod = keyInfoCanonicalizationMethod
}

// xadesSignatureParametersAssertCanonicalizationNotEmpty ports the private static
// assertCanonicalizationNotEmpty(String). Panics with the Java message when
// canonicalizationMethod is empty (IllegalArgumentException upstream).
func xadesSignatureParametersAssertCanonicalizationNotEmpty(canonicalizationMethod string) {
	if utils.IsStringEmpty(canonicalizationMethod) {
		panic("Canonicalization cannot be empty! See EN 319 132-1: 3.1.2 Signature Generation.")
	}
}

// IsSignKeyInfo returns whether the "KeyInfo" element must be signed. Ports isSignKeyInfo().
func (p *XAdESSignatureParameters) IsSignKeyInfo() bool {
	return p.signKeyInfo
}

// SetSignKeyInfo sets the parameter SignKeyInfo defining if the "KeyInfo" element must be
// signed and its reference must be included to "SignedInfo" element. The value is FALSE by
// default. Ports setSignKeyInfo(boolean).
func (p *XAdESSignatureParameters) SetSignKeyInfo(signKeyInfo bool) {
	p.signKeyInfo = signKeyInfo
}

// References returns a list of references to be incorporated to the signature. Ports
// getReferences().
func (p *XAdESSignatureParameters) References() []*DSSReference {
	context := p.GetContext()
	if context != nil && utils.IsCollectionNotEmpty(context.References()) {
		return context.References()
	}
	if utils.IsCollectionNotEmpty(p.dssReferences) {
		return p.dssReferences
	}
	return nil
}

// SetReferences sets a list of references to be incorporated into the signature. NOTE: This
// method overwrites a default behavior on ds:Reference's creation. It should be used only by
// experienced users. Ports setReferences(List<DSSReference>).
func (p *XAdESSignatureParameters) SetReferences(references []*DSSReference) {
	p.dssReferences = references
}

// XPathLocationString gets the xPath signature location string (ENVELOPED only). Ports
// getXPathLocationString().
func (p *XAdESSignatureParameters) XPathLocationString() string {
	return p.xPathLocationString
}

// SetXPathLocationString defines the position where the signature will be added (XAdES
// Enveloped). Ports setXPathLocationString(String).
func (p *XAdESSignatureParameters) SetXPathLocationString(xPathLocationString string) {
	p.xPathLocationString = xPathLocationString
}

// XPathElementPlacement returns the XPath element placement for Enveloped signature creation.
// Ports getXPathElementPlacement().
func (p *XAdESSignatureParameters) XPathElementPlacement() XPathElementPlacement {
	return p.xPathElementPlacement
}

// SetXPathElementPlacement defines the relation to the element referenced by the XPath where
// the signature will be added (XAdES Enveloped). Ports setXPathElementPlacement(XPathElementPlacement).
func (p *XAdESSignatureParameters) SetXPathElementPlacement(xPathElementPlacement XPathElementPlacement) {
	p.xPathElementPlacement = xPathElementPlacement
}

// RootDocument returns the root document for INTERNALLY_DETACHED signature creation. Ports
// getRootDocument().
func (p *XAdESSignatureParameters) RootDocument() *xmldom.Node {
	return p.rootDocument
}

// SetRootDocument sets the root document for INTERNALLY_DETACHED signature creation. Ports the
// Document overload of setRootDocument.
func (p *XAdESSignatureParameters) SetRootDocument(rootDocument *xmldom.Node) {
	p.rootDocument = rootDocument
}

// SetRootDocumentFromDSSDocument sets the root XML document for a signature creation. This
// method expects rootDocument to be represented by a valid XML document.
//
// DEVIATION preserved verbatim: Java's overload does not return after handling a nil
// rootDocument (it calls setRootDocument((Document) null) and then falls through into the
// isDOM/buildDOM checks with the same nil value), so a nil rootDocument here still panics with
// the "valid XML document" message below rather than clearing rootDocument and returning -
// exactly as DomUtils.isDOM(null) safely answers false in Java (broad catch) and this port's
// DomUtilsIsDOM does too (recover-based). Ports the DSSDocument overload of setRootDocument.
func (p *XAdESSignatureParameters) SetRootDocumentFromDSSDocument(rootDocument model.DSSDocument) {
	if rootDocument == nil {
		p.SetRootDocument(nil)
	}
	if !xmlutils.DomUtilsIsDOM(rootDocument) {
		panic("The rootDocument shall be represented by a valid XML document!")
	}
	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(rootDocument)
	if err != nil {
		panic(err)
	}
	p.SetRootDocument(dom)
}

// GetContext gets the signature creation context (internal variable). Ports the overridden
// getContext(). See the file header for why this shadows, rather than reuses, the promoted
// AbstractSignatureParameters.GetContext.
func (p *XAdESSignatureParameters) GetContext() *XAdESProfileParameters {
	if p.context == nil {
		p.context = NewXAdESProfileParameters()
	}
	return p.context
}

// GetDeterministicId returns the deterministic identifier used for unique identification of a
// created signature, built through the shadowed GetContext so every reader of it (including
// XAdESCounterSignatureParameters's own override) observes the same cached value. See the file
// header. Ports the inherited AbstractSignatureParameters#getDeterministicId as reached through
// Java's virtual getContext() dispatch.
func (p *XAdESSignatureParameters) GetDeterministicId() string {
	deterministicId := p.GetContext().DeterministicId()
	if deterministicId == "" {
		var identifier *model.TokenIdentifier
		if signingCertificate := p.SigningCertificate(); signingCertificate != nil {
			identifier = signingCertificate.BuildTokenIdentifier()
		}
		var signingDate time.Time
		if sd := p.BLevel().SigningDate(); sd != nil {
			signingDate = *sd
		}
		newDeterministicId, err := spi.DSSUtilsDeterministicID(signingDate, identifier)
		if err != nil {
			panic(err)
		}
		deterministicId = newDeterministicId
		p.GetContext().SetDeterministicId(deterministicId)
	}
	return deterministicId
}

// DetachedContents returns the documents to sign, checking the shadowed XAdES context first and
// falling back to the promoted base behaviour. See the file header. Ports the inherited
// AbstractSignatureParameters#getDetachedContents as reached through Java's virtual getContext()
// dispatch.
func (p *XAdESSignatureParameters) DetachedContents() []model.DSSDocument {
	if detachedContents := p.GetContext().DetachedContents(); len(detachedContents) > 0 {
		return detachedContents
	}
	return p.AbstractSignatureParameters.DetachedContents()
}

// IsEn319132 gets if the signature shall be created according to ETSI EN 319 132. Ports
// isEn319132().
func (p *XAdESSignatureParameters) IsEn319132() bool {
	return p.en319132
}

// SetEn319132 sets if the signature shall be created according to ETSI EN 319 132-1. Default:
// true. Ports setEn319132(boolean).
func (p *XAdESSignatureParameters) SetEn319132(en319132 bool) {
	p.en319132 = en319132
}

// IsEmbedXML gets if the signed content shall be incorporated as XML (used for ENVELOPING).
// Ports isEmbedXML().
func (p *XAdESSignatureParameters) IsEmbedXML() bool {
	return p.embedXML
}

// SetEmbedXML sets if the signed content shall be incorporated as XML (used for ENVELOPING). If
// false, incorporates the document content in its base64 encoded representation. Default:
// false. Ports setEmbedXML(boolean).
func (p *XAdESSignatureParameters) SetEmbedXML(embedXML bool) {
	p.embedXML = embedXML
}

// IsManifestSignature gets if the signature signs a manifest. Ports isManifestSignature().
func (p *XAdESSignatureParameters) IsManifestSignature() bool {
	return p.manifestSignature
}

// SetManifestSignature sets if the signature signs a manifest. Ports
// setManifestSignature(boolean).
func (p *XAdESSignatureParameters) SetManifestSignature(manifestSignature bool) {
	p.manifestSignature = manifestSignature
}

// IsAddX509SubjectName gets if the ds:X509Data element shall be added. Ports
// isAddX509SubjectName().
func (p *XAdESSignatureParameters) IsAddX509SubjectName() bool {
	return p.addX509SubjectName
}

// SetAddX509SubjectName sets if the ds:X509Data element shall be added. Default: false. Ports
// setAddX509SubjectName(boolean).
func (p *XAdESSignatureParameters) SetAddX509SubjectName(addX509SubjectName bool) {
	p.addX509SubjectName = addX509SubjectName
}

// SignedAdESObject gets a custom XAdES Object content. Ports getSignedAdESObject().
func (p *XAdESSignatureParameters) SignedAdESObject() []byte {
	return p.signedAdESObject
}

// SetSignedAdESObject sets a custom XAdES Object content to incorporate into the signature.
// Ports setSignedAdESObject(byte[]).
func (p *XAdESSignatureParameters) SetSignedAdESObject(signedAdESObject []byte) {
	p.signedAdESObject = signedAdESObject
}

// IsPrettyPrint gets if the signature shall be pretty-printed. Ports isPrettyPrint().
func (p *XAdESSignatureParameters) IsPrettyPrint() bool {
	return p.prettyPrint
}

// SetPrettyPrint sets if the signature shall be pretty-printed. Default: false. Ports
// setPrettyPrint(boolean).
func (p *XAdESSignatureParameters) SetPrettyPrint(prettyPrint bool) {
	p.prettyPrint = prettyPrint
}

// XmldsigNamespace returns the current used XMLDSig namespace. Never returns nil. Ports
// getXmldsigNamespace().
func (p *XAdESSignatureParameters) XmldsigNamespace() *common.DSSNamespace {
	return p.xmldsigNamespace
}

// SetXmldsigNamespace sets the XMLDSIG namespace. Default: ds:http://www.w3.org/2000/09/xmldsig#.
// Panics with the Java message when xmldsigNamespace is nil (Objects.requireNonNull upstream) or
// its URI is not accepted (IllegalArgumentException upstream). Ports
// setXmldsigNamespace(DSSNamespace).
func (p *XAdESSignatureParameters) SetXmldsigNamespace(xmldsigNamespace *common.DSSNamespace) {
	if xmldsigNamespace == nil {
		panic("xmldsigNamespace must not be null")
	}
	uri := xmldsigNamespace.Uri()
	if !common.XMLDSigNS.IsSameUri(uri) {
		panic("Not accepted URI")
	}
	p.xmldsigNamespace = xmldsigNamespace
}

// XadesNamespace returns the current used XAdES namespace. Never returns nil. Ports
// getXadesNamespace().
func (p *XAdESSignatureParameters) XadesNamespace() *common.DSSNamespace {
	return p.xadesNamespace
}

// SetXadesNamespace sets the XAdES namespace. Default: xades:http://uri.etsi.org/01903/v1.3.2#.
// Panics with the Java message when xadesNamespace is nil (Objects.requireNonNull upstream) or
// its URI is not accepted (IllegalArgumentException upstream). Ports
// setXadesNamespace(DSSNamespace).
func (p *XAdESSignatureParameters) SetXadesNamespace(xadesNamespace *common.DSSNamespace) {
	if xadesNamespace == nil {
		panic("xadesNamespace must not be null")
	}
	uri := xadesNamespace.Uri()
	if definition.XAdESNamespace_XADES_111.IsSameUri(uri) ||
		definition.XAdESNamespace_XADES_122.IsSameUri(uri) ||
		definition.XAdESNamespace_XADES_132.IsSameUri(uri) {
		p.xadesNamespace = xadesNamespace
	} else {
		panic("Not accepted URI")
	}
}

// Xades141Namespace returns the current used XAdES 1.4.1 namespace. Never returns nil. Ports
// getXades141Namespace().
func (p *XAdESSignatureParameters) Xades141Namespace() *common.DSSNamespace {
	return p.xades141Namespace
}

// SetXades141Namespace sets the XAdES 1.4.1 namespace. Default:
// xades141:http://uri.etsi.org/01903/v1.4.1#. Panics with the Java message when
// xades141Namespace is nil (Objects.requireNonNull upstream) or its URI is not accepted
// (IllegalArgumentException upstream). Ports setXades141Namespace(DSSNamespace).
func (p *XAdESSignatureParameters) SetXades141Namespace(xades141Namespace *common.DSSNamespace) {
	if xades141Namespace == nil {
		panic("xades141Namespace must not be null")
	}
	uri := xades141Namespace.Uri()
	if !definition.XAdESNamespace_XADES_141.IsSameUri(uri) {
		panic("Not accepted URI")
	}
	p.xades141Namespace = xades141Namespace
}

// Objects gets the list of custom ds:Object elements. Ports getObjects().
func (p *XAdESSignatureParameters) Objects() []*DSSObject {
	return p.objects
}

// SetObjects sets the list of custom ds:Object elements to be incorporated within the
// ds:Signature. Ports setObjects(List<DSSObject>).
func (p *XAdESSignatureParameters) SetObjects(objects []*DSSObject) {
	p.objects = objects
}

// TokenReferencesDigestAlgorithm gets a DigestAlgorithm to create
// CompleteCertificateRefs/CompleteRevocationRefs with. Ports
// getTokenReferencesDigestAlgorithm().
func (p *XAdESSignatureParameters) TokenReferencesDigestAlgorithm() enumerations.DigestAlgorithm {
	return p.tokenReferencesDigestAlgorithm
}

// SetTokenReferencesDigestAlgorithm sets a DigestAlgorithm to create
// CompleteCertificateRefs/CompleteRevocationRefs for -C level. Default: SHA512. Panics with the
// Java message when tokenReferencesDigestAlgorithm is empty (Objects.requireNonNull upstream).
// Ports setTokenReferencesDigestAlgorithm(DigestAlgorithm).
func (p *XAdESSignatureParameters) SetTokenReferencesDigestAlgorithm(tokenReferencesDigestAlgorithm enumerations.DigestAlgorithm) {
	if tokenReferencesDigestAlgorithm == "" {
		panic("TokenReferencesDigestAlgorithm cannot be null!")
	}
	p.tokenReferencesDigestAlgorithm = tokenReferencesDigestAlgorithm
}

// DataObjectFormatList gets a list of custom xades:DataObjectFormat elements. Ports
// getDataObjectFormatList().
func (p *XAdESSignatureParameters) DataObjectFormatList() []*DSSDataObjectFormat {
	return p.dataObjectFormatList
}

// SetDataObjectFormatList sets a list of custom xades:DataObjectFormat elements to be
// incorporated within xades:SignedDataObjectProperties element of the signature. NOTE: this
// method overwrites default behavior on xades:DataObjectFormat creation. It should be used only
// by experienced users. Ports setDataObjectFormatList(List<DSSDataObjectFormat>).
func (p *XAdESSignatureParameters) SetDataObjectFormatList(dataObjectFormatList []*DSSDataObjectFormat) {
	p.dataObjectFormatList = dataObjectFormatList
}

// GetContentTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating XAdESTimestampParameters. Ports the overridden #getContentTimestampParameters.
func (p *XAdESSignatureParameters) GetContentTimestampParameters() *XAdESTimestampParameters {
	if p.ContentTimestampParameters == nil {
		p.ContentTimestampParameters = NewXAdESTimestampParameters()
	}
	return p.ContentTimestampParameters
}

// GetSignatureTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating XAdESTimestampParameters. Ports the overridden #getSignatureTimestampParameters.
func (p *XAdESSignatureParameters) GetSignatureTimestampParameters() *XAdESTimestampParameters {
	if p.SignatureTimestampParameters == nil {
		p.SignatureTimestampParameters = NewXAdESTimestampParameters()
	}
	return p.SignatureTimestampParameters
}

// GetArchiveTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating XAdESTimestampParameters. Ports the overridden #getArchiveTimestampParameters.
func (p *XAdESSignatureParameters) GetArchiveTimestampParameters() *XAdESTimestampParameters {
	if p.ArchiveTimestampParameters == nil {
		p.ArchiveTimestampParameters = NewXAdESTimestampParameters()
	}
	return p.ArchiveTimestampParameters
}

// Reinit ports the overridden #reinit, delegating to the promoted base Reinit (harmless: it
// clears the base's own separate, otherwise-unused context field) and clearing the shadowed
// XAdES context. See the file header.
func (p *XAdESSignatureParameters) Reinit() {
	p.AbstractSignatureParameters.Reinit()
	p.context = nil
}

// String ports toString().
func (p *XAdESSignatureParameters) String() string {
	return fmt.Sprintf("XAdESSignatureParameters [addX509SubjectName=%v, dssReferences=%v, embedXML=%v, "+
		"en319132=%v, keyInfoCanonicalizationMethod='%s', signedInfoCanonicalizationMethod='%s', "+
		"signedPropertiesCanonicalizationMethod='%s', manifestSignature=%v, rootDocument=%v, "+
		"signedAdESObject=%v, signingCertificateDigestMethod=%v, signKeyInfo=%v, xPathLocationString='%s', "+
		"xPathElementPlacement=%v, prettyPrint=%v, xmldsigNamespace=%v, xadesNamespace=%v, "+
		"xades141Namespace=%v, objects=%v, tokenReferencesDigestAlgorithm=%v, dataObjectFormatList=%v] %s",
		p.addX509SubjectName, p.dssReferences, p.embedXML, p.en319132, p.keyInfoCanonicalizationMethod,
		p.signedInfoCanonicalizationMethod, p.signedPropertiesCanonicalizationMethod, p.manifestSignature,
		p.rootDocument, p.signedAdESObject, p.signingCertificateDigestMethod, p.signKeyInfo,
		p.xPathLocationString, p.xPathElementPlacement, p.prettyPrint, p.xmldsigNamespace, p.xadesNamespace,
		p.xades141Namespace, p.objects, p.tokenReferencesDigestAlgorithm, p.dataObjectFormatList,
		p.AbstractSignatureParameters.String())
}

// Equals ports equals(Object).
func (p *XAdESSignatureParameters) Equals(other *XAdESSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.AbstractSignatureParameters.Equals(&other.AbstractSignatureParameters) {
		return false
	}
	if p.addX509SubjectName != other.addX509SubjectName ||
		p.embedXML != other.embedXML ||
		p.en319132 != other.en319132 ||
		p.manifestSignature != other.manifestSignature ||
		p.signKeyInfo != other.signKeyInfo ||
		p.prettyPrint != other.prettyPrint {
		return false
	}
	if !xadesSignatureParametersDSSReferencesEqual(p.dssReferences, other.dssReferences) {
		return false
	}
	if p.keyInfoCanonicalizationMethod != other.keyInfoCanonicalizationMethod ||
		p.signedInfoCanonicalizationMethod != other.signedInfoCanonicalizationMethod ||
		p.signedPropertiesCanonicalizationMethod != other.signedPropertiesCanonicalizationMethod {
		return false
	}
	if !xadesSignatureParametersRootDocumentEqual(p.rootDocument, other.rootDocument) {
		return false
	}
	if string(p.signedAdESObject) != string(other.signedAdESObject) {
		return false
	}
	if p.signingCertificateDigestMethod != other.signingCertificateDigestMethod {
		return false
	}
	if p.xPathLocationString != other.xPathLocationString || p.xPathElementPlacement != other.xPathElementPlacement {
		return false
	}
	if !xadesSignatureParametersNamespaceEqual(p.xmldsigNamespace, other.xmldsigNamespace) ||
		!xadesSignatureParametersNamespaceEqual(p.xadesNamespace, other.xadesNamespace) ||
		!xadesSignatureParametersNamespaceEqual(p.xades141Namespace, other.xades141Namespace) {
		return false
	}
	if !xadesSignatureParametersObjectsEqual(p.objects, other.objects) {
		return false
	}
	if p.tokenReferencesDigestAlgorithm != other.tokenReferencesDigestAlgorithm {
		return false
	}
	return xadesSignatureParametersDataObjectFormatListEqual(p.dataObjectFormatList, other.dataObjectFormatList)
}

// xadesSignatureParametersDSSReferencesEqual ports Objects.equals(dssReferences, ...).
func xadesSignatureParametersDSSReferencesEqual(a, b []*DSSReference) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// xadesSignatureParametersRootDocumentEqual ports Objects.equals(rootDocument, ...): a
// pointer-identity comparison, since *xmldom.Node has no value-Equals method and Java's
// org.w3c.dom.Document does not override equals() either (identity semantics).
func xadesSignatureParametersRootDocumentEqual(a, b *xmldom.Node) bool {
	return a == b
}

// xadesSignatureParametersNamespaceEqual ports Objects.equals(...Namespace, ...).
func xadesSignatureParametersNamespaceEqual(a, b *common.DSSNamespace) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Uri() == b.Uri() && a.Prefix() == b.Prefix()
}

// xadesSignatureParametersObjectsEqual ports Objects.equals(objects, ...).
func xadesSignatureParametersObjectsEqual(a, b []*DSSObject) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// xadesSignatureParametersDataObjectFormatListEqual ports Objects.equals(dataObjectFormatList,
// ...).
func xadesSignatureParametersDataObjectFormatListEqual(a, b []*DSSDataObjectFormat) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

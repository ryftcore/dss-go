// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESBuilder.java (DSS 6.5.RC1).
//
// # The abstract base, and how its virtual methods survive the port
//
// Java's XAdESBuilder is abstract: it declares alignNodes() abstract and lets ExtensionBuilder
// override getXmldsigNamespace()/getXadesNamespace(), while its own concrete methods
// (incorporateCert, incorporateCertDigest, incorporateSPDocSpecification, createXmlDocument)
// call back into those three through Java's virtual dispatch. Go has no method overriding
// across embedding, so - exactly as document/abstract_document_extender.go and
// cades/cades_signature_extension.go already do, per the TokenBase.InitToken(self) convention
// of PORTING.md - the three are collected in BuilderOverrides, every concrete builder
// registers itself with InitXAdESBuilder/InitXAdESBuilderWithVerifier, and this base always
// reaches them through b.overrides. Builder keeps its own XmldsigNamespace/XadesNamespace
// as the "super" implementations an override may delegate back to.
//
// A grep over dss-xades confirms the closed set of overridden members: alignNodes
// (ExtensionBuilder, AbstractSignatureBuilder), getXmldsigNamespace and getXadesNamespace
// (ExtensionBuilder). getXades141Namespace, getCurrentXAdESElements, getCurrentXAdESPath and
// createXmlDocument are never overridden and therefore stay plain methods here - they still
// consult b.overrides for the namespace they are built on, so an ExtensionBuilder sees the
// namespace taken from the signature being extended, as in Java.
//
// # Errors
//
// Java's unchecked IllegalArgumentException on an unsupported XAdES URI becomes a returned
// error (PORTING.md: throw -> (T, error)); the token digest and the DOM serialization, which
// Java performs behind DSSException-free APIs, can fail in Go and are likewise reported.
// slf4j logging and java.io.Serializable are dropped.
package xades

import (
	"fmt"

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

// BuilderReferencedData is the attribute used for timestamp includes.
// Port of the public constant REFERENCED_DATA.
const BuilderReferencedData = "referencedData"

// BuilderTarget is the qualifying properties target.
// Port of the public constant TARGET.
const BuilderTarget = "Target"

// XAdESBuilderURI is the URI attribute. Port of the public constant URI.
const BuilderURI = "URI"

const (
	// xadesBuilderIDPrefix is the Xml Id prefix. Port of the protected ID_PREFIX.
	xadesBuilderIDPrefix = "id-"

	// xadesBuilderTimestampPrefix is the Id-prefix for the TimeStamp element.
	// Port of the protected TIMESTAMP_PREFIX.
	xadesBuilderTimestampPrefix = "ts-"

	// xadesBuilderEncapsulatedTimestampPrefix is the Id-prefix for the EncapsulatedTimeStamp
	// element. Port of the protected ENCAPSULATED_TIMESTAMP_PREFIX.
	xadesBuilderEncapsulatedTimestampPrefix = "ets-"

	// xadesBuilderTstVdPrefix is the Id-prefix for the TimeStampValidationData element.
	// Port of the protected TST_VD_PREFIX.
	xadesBuilderTstVdPrefix = "tsvd-"

	// xadesBuilderAnyVdPrefix is the Id-prefix for the AnyValidationData element.
	// Port of the protected ANY_VD_PREFIX.
	xadesBuilderAnyVdPrefix = "avd-"
)

// XAdESBuilderOverrides declares the members Java's abstract XAdESBuilder either leaves
// abstract or expects a subclass to override, and that the base implementation itself calls
// back into. Every concrete builder satisfies this and registers itself through
// InitXAdESBuilder / InitXAdESBuilderWithVerifier.
type BuilderOverrides interface {
	// AlignNodes aligns children indents. Port of the protected abstract #alignNodes.
	AlignNodes()

	// XmldsigNamespace returns the currently used XMLDSig namespace.
	// Port of the protected #getXmldsigNamespace (overridden by ExtensionBuilder).
	XmldsigNamespace() *common.DSSNamespace

	// XadesNamespace returns the currently used XAdES namespace.
	// Port of the protected #getXadesNamespace (overridden by ExtensionBuilder).
	XadesNamespace() *common.DSSNamespace
}

// Builder builds a XAdES signature. It is the abstract base of every XAdES builder:
// AbstractSignatureBuilder for signature creation, ExtensionBuilder (and through it every
// -T/-C/-X/-XL/-A/-LT/-LTA level and the SignaturePolicyStoreBuilder) for extension.
type Builder struct {
	// XadesPath holds the XAdESPath implementation which contains all constants and queries
	// needed to cope with the default signature schema. Port of the protected xadesPath.
	XadesPath definition.XAdESPath

	// Params is the set of parameters relating to the structure and process of the creation or
	// extension of the electronic signature. Port of the protected params.
	Params *SignatureParameters

	// DocumentDom is the root XML document root (with signature).
	// Port of the protected documentDom.
	DocumentDom *xmldom.Node

	// CertificateVerifier is the object in charge of certificates validation.
	// Port of the protected certificateVerifier.
	CertificateVerifier validation.CertificateVerifier

	// overrides points back at the concrete builder; see InitXAdESBuilder. Left nil (bare zero
	// value) panics on first use, matching the TokenBase.InitToken(self) convention.
	overrides BuilderOverrides
}

// InitXAdESBuilder registers the concrete builder with its base.
// Port of the protected empty XAdESBuilder() constructor.
func (b *Builder) InitXAdESBuilder(self BuilderOverrides) {
	b.overrides = self
}

// InitXAdESBuilderWithVerifier registers the concrete builder with its base and stores the
// CertificateVerifier. Port of the protected XAdESBuilder(CertificateVerifier) constructor.
func (b *Builder) InitXAdESBuilderWithVerifier(self BuilderOverrides,
	certificateVerifier validation.CertificateVerifier) {
	b.overrides = self
	b.CertificateVerifier = certificateVerifier
}

// IncorporateCertDigest creates the xades:CertDigest DOM object.
//
//	<CertDigest>
//	   <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//	   <ds:DigestValue>fj8SJujSXU4fi342bdtiKVbglA0=</ds:DigestValue>
//	</CertDigest>
//
// Port of the protected #incorporateCertDigest.
func (b *Builder) IncorporateCertDigest(parentDom *xmldom.Node,
	digestAlgorithm enumerations.DigestAlgorithm, token model.Token) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	certDigestDom := xmlutils.DomUtilsAddElement(b.DocumentDom, parentDom, b.overrides.XadesNamespace(),
		currentElements.ElementCertDigest())
	b.IncorporateDigestMethod(certDigestDom, digestAlgorithm)
	return b.IncorporateDigestValueOfToken(certDigestDom, digestAlgorithm, token)
}

// IncorporateDigestMethod creates the ds:DigestMethod DOM object.
//
//	<ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//
// Port of the protected #incorporateDigestMethod.
func (b *Builder) IncorporateDigestMethod(parentDom *xmldom.Node,
	digestAlgorithm enumerations.DigestAlgorithm) {
	namespace := b.digestAlgAndValueNamespace()
	DSSXMLUtilsIncorporateDigestMethod(parentDom, digestAlgorithm, namespace)
}

// IncorporateDigestValueOfToken creates the ds:DigestValue DOM object for a given Token.
//
//	<ds:DigestValue>fj8SJujSXU4fi342bdtiKVbglA0=</ds:DigestValue>
//
// Port of the protected #incorporateDigestValue(Element, DigestAlgorithm, Token); Java's two
// incorporateDigestValue overloads cannot share one Go name, so the Token-taking one carries
// the OfToken suffix and the String-taking one keeps the plain name.
func (b *Builder) IncorporateDigestValueOfToken(parentDom *xmldom.Node,
	digestAlgorithm enumerations.DigestAlgorithm, token model.Token) error {
	digest, err := token.Digest(digestAlgorithm)
	if err != nil {
		return err
	}
	base64EncodedDigestBytes := utils.ToBase64(digest)
	// Upstream traces "Digest value for the token with Id [{}] --> {}" here.
	b.IncorporateDigestValue(parentDom, base64EncodedDigestBytes)
	return nil
}

// IncorporateDigestValue creates the ds:DigestValue DOM object.
//
//	<ds:DigestValue>fj8SJujSXU4fi342bdtiKVbglA0=</ds:DigestValue>
//
// Port of the protected #incorporateDigestValue(Element, String).
func (b *Builder) IncorporateDigestValue(parentDom *xmldom.Node, base64EncodedDigestBytes string) {
	namespace := b.digestAlgAndValueNamespace()
	DSSXMLUtilsIncorporateDigestValue(parentDom, base64EncodedDigestBytes, namespace)
}

// digestAlgAndValueNamespace ports the private getDigestAlgAndValueNamespace.
func (b *Builder) digestAlgAndValueNamespace() *common.DSSNamespace {
	if definition.XAdESNamespaceXAdES111.IsSameUri(b.overrides.XadesNamespace().Uri()) {
		return b.overrides.XadesNamespace()
	}
	return b.overrides.XmldsigNamespace()
}

// IncorporateCert creates the xades:Cert DOM object:
//
//	<Cert>
//	    <CertDigest>
//	        <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//	        <ds:DigestValue>fj8SJujSXU4fi342bdtiKVbglA0=</ds:DigestValue>
//	    </CertDigest>
//	    <IssuerSerial>
//	        <ds:X509IssuerName>CN=ICA A,O=DSS,C=AA</ds:X509IssuerName>
//	        <ds:X509SerialNumber>4</ds:X509SerialNumber>
//	    </IssuerSerial>
//	</Cert>
//
// Port of the protected #incorporateCert.
func (b *Builder) IncorporateCert(parentDom *xmldom.Node, certificate *model.CertificateToken,
	digestAlgorithm enumerations.DigestAlgorithm) (*xmldom.Node, error) {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return nil, err
	}
	certDom := xmlutils.DomUtilsAddElement(b.DocumentDom, parentDom, b.overrides.XadesNamespace(),
		currentElements.ElementCert())

	if err := b.IncorporateCertDigest(certDom, digestAlgorithm, certificate); err != nil {
		return nil, err
	}

	if b.Params.IsEn319132() {
		if err := b.IncorporateIssuerV2(certDom, certificate); err != nil {
			return nil, err
		}
	} else {
		if err := b.IncorporateIssuerV1(certDom, certificate); err != nil {
			return nil, err
		}
	}

	return certDom, nil
}

// IncorporateIssuerV1 incorporates the xades:IssuerSerial element.
// Port of the protected #incorporateIssuerV1.
func (b *Builder) IncorporateIssuerV1(parentDom *xmldom.Node, certificate *model.CertificateToken) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	issuerSerialDom := xmlutils.DomUtilsAddElement(b.DocumentDom, parentDom, b.overrides.XadesNamespace(),
		currentElements.ElementIssuerSerial())
	x509IssuerNameDom := xmlutils.DomUtilsAddElement(b.DocumentDom, issuerSerialDom,
		b.overrides.XmldsigNamespace(), common.XMLDSigElementX509IssuerName)

	// X500Principal#getName() with no format argument is RFC 2253.
	issuerX500PrincipalName := certificate.IssuerX500Principal().RFC2253Name()
	xmlutils.DomUtilsSetTextNode(b.DocumentDom, x509IssuerNameDom, issuerX500PrincipalName)

	x509SerialNumberDom := xmlutils.DomUtilsAddElement(b.DocumentDom, issuerSerialDom,
		b.overrides.XmldsigNamespace(), common.XMLDSigElementX509SerialNumber)

	serialNumber := certificate.SerialNumber()
	serialNumberString := serialNumber.String()
	xmlutils.DomUtilsSetTextNode(b.DocumentDom, x509SerialNumberDom, serialNumberString)
	return nil
}

// IncorporateIssuerV2 incorporates the xades:IssuerSerialV2 element.
// Port of the protected #incorporateIssuerV2.
func (b *Builder) IncorporateIssuerV2(parentDom *xmldom.Node, certificate *model.CertificateToken) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	issuerSerialDom := xmlutils.DomUtilsAddElement(b.DocumentDom, parentDom, b.overrides.XadesNamespace(),
		currentElements.ElementIssuerSerialV2())

	issuerSerial := spi.DSSASN1UtilsIssuerSerialForCertificate(certificate)
	// DSSASN1Utils.getDEREncoded(IssuerSerial) is the structure's own DER, which
	// internal/asn1ber's IssuerSerial produces directly - no re-encoding through a second
	// ASN.1 path, as PORTING.md requires.
	issuerBase64 := utils.ToBase64(issuerSerial.DER())
	xmlutils.DomUtilsSetTextNode(b.DocumentDom, issuerSerialDom, issuerBase64)
	return nil
}

// notIndentedObjectIds returns the list of object ids that must not be indented in any case.
// Port of the private getNotIndentedObjectIds.
func (b *Builder) notIndentedObjectIds() []string {
	ids := make([]string, 0)
	dssReferences := b.Params.References()
	for _, reference := range dssReferences {
		// do not change external objects
		if DSSXMLUtilsIsObjectReferenceType(reference.Type()) {
			ids = append(ids, xmlutils.DomUtilsGetId(reference.Uri()))
		}
	}
	return ids
}

// CreateXmlDocument creates a DSSDocument from the current documentDom.
// Port of the protected #createXmlDocument.
func (b *Builder) CreateXmlDocument() (model.DSSDocument, error) {
	var bytes []byte
	var err error
	if enumerations.SigningOperationSign == b.Params.GetContext().OperationKind() && b.Params.IsPrettyPrint() {
		b.overrides.AlignNodes()
		indented, indentErr := DSSXMLUtilsGetDocWithIndentedSignature(b.DocumentDom,
			b.Params.GetDeterministicId(), b.notIndentedObjectIds())
		if indentErr != nil {
			return nil, indentErr
		}
		bytes, err = xmlutils.DomUtilsSerializeNode(indented)
	} else {
		bytes, err = xmlutils.DomUtilsSerializeNode(b.DocumentDom)
	}
	if err != nil {
		return nil, err
	}
	inMemoryDocument := model.NewInMemoryDocument(bytes)
	inMemoryDocument.SetMimeType(enumerations.MimeTypeEnumXML)
	return inMemoryDocument, nil
}

// XmldsigNamespace returns the currently used XMLDSig namespace.
// Port of the protected #getXmldsigNamespace; ExtensionBuilder overrides it and delegates back
// here when the signature under extension carries no namespace of its own.
func (b *Builder) XmldsigNamespace() *common.DSSNamespace {
	return b.Params.XmldsigNamespace()
}

// XadesNamespace returns the currently used XAdES namespace.
// Port of the protected #getXadesNamespace; see XmldsigNamespace for the override note.
func (b *Builder) XadesNamespace() *common.DSSNamespace {
	return b.Params.XadesNamespace()
}

// Xades141Namespace returns the currently used XAdES 1.4.1 namespace.
// Port of the protected #getXades141Namespace.
func (b *Builder) Xades141Namespace() *common.DSSNamespace {
	return b.Params.Xades141Namespace()
}

// CurrentXAdESElements gets the relevant XAdESElement implementation for the namespace in use.
// Port of the protected #getCurrentXAdESElements; Java's XAdES132Element.values()[0] merely
// needs any enum instance to reach the interface's instance methods (which ignore the
// receiver), so the Go port names that same first constant. Java's IllegalArgumentException
// becomes a returned error.
func (b *Builder) CurrentXAdESElements() (definition.XAdESElement, error) {
	xadesURI := b.overrides.XadesNamespace().Uri()
	switch {
	case definition.XAdESNamespaceXAdES132.Uri() == xadesURI:
		return definition.XAdES132ElementAllDataObjectsTimestamp, nil
	case definition.XAdESNamespaceXAdES122.Uri() == xadesURI:
		return definition.XAdES122ElementAllDataObjectsTimestamp, nil
	case definition.XAdESNamespaceXAdES111.Uri() == xadesURI:
		return definition.XAdES111ElementAllDataObjectsTimestamp, nil
	}
	return nil, fmt.Errorf("Unsupported URI : %s", xadesURI)
}

// CurrentXAdESPath gets the relevant XAdESPath implementation for the namespace in use.
// Port of the protected #getCurrentXAdESPath.
func (b *Builder) CurrentXAdESPath() (definition.XAdESPath, error) {
	xadesURI := b.overrides.XadesNamespace().Uri()
	switch {
	case utils.AreStringsEqual(definition.XAdESNamespaceXAdES132.Uri(), xadesURI):
		return definition.NewXAdES132Path(), nil
	case utils.AreStringsEqual(definition.XAdESNamespaceXAdES122.Uri(), xadesURI):
		return definition.NewXAdES122Path(), nil
	case utils.AreStringsEqual(definition.XAdESNamespaceXAdES111.Uri(), xadesURI):
		return definition.NewXAdES111Path(), nil
	}
	return nil, fmt.Errorf("Unsupported URI : %s", xadesURI)
}

// IncorporateSPDocSpecification creates the SpDocSpecification DOM object:
//
//	<xades141:SPDocSpecification xmlns:xades141="http://uri.etsi.org/01903/v1.4.1#">
//	    <xades:Identifier>1.2.3.4.5</xades:Identifier>
//	</xades141:SPDocSpecification>
//
// Port of the protected #incorporateSPDocSpecification.
func (b *Builder) IncorporateSPDocSpecification(parentElement *xmldom.Node,
	spDocSpecification *model.SpDocSpecification) error {
	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	spDocSpecElement := xmlutils.DomUtilsAddElement(b.DocumentDom, parentElement,
		b.Xades141Namespace(), definition.XAdES141ElementSPDocSpecification)
	xmlutils.DomUtilsAddNamespaceAttribute(spDocSpecElement, b.Xades141Namespace())

	identifierElement := xmlutils.DomUtilsAddElement(b.DocumentDom, spDocSpecElement,
		b.overrides.XadesNamespace(), currentElements.ElementIdentifier())
	if spDocSpecification.Qualifier() != "" {
		identifierElement.SetAttr(
			xmldom.Name{Local: definition.XAdES132AttributeQualifier.AttributeName()},
			spDocSpecification.Qualifier().Value())
	}
	xmlutils.DomUtilsSetTextNode(b.DocumentDom, identifierElement, spDocSpecification.Id())

	if utils.IsStringNotEmpty(spDocSpecification.Description()) {
		descriptionElement := xmlutils.DomUtilsAddElement(b.DocumentDom, spDocSpecElement,
			b.overrides.XadesNamespace(), currentElements.ElementDescription())
		xmlutils.DomUtilsSetTextNode(b.DocumentDom, descriptionElement, spDocSpecification.Description())
	}

	if utils.IsArrayNotEmpty(spDocSpecification.DocumentationReferences()) {
		documentReferencesElement := xmlutils.DomUtilsAddElement(b.DocumentDom, spDocSpecElement,
			b.overrides.XadesNamespace(), currentElements.ElementDocumentationReferences())

		for _, docRef := range spDocSpecification.DocumentationReferences() {
			documentReferenceElement := xmlutils.DomUtilsAddElement(b.DocumentDom,
				documentReferencesElement, b.overrides.XadesNamespace(),
				currentElements.ElementDocumentationReference())
			xmlutils.DomUtilsSetTextNode(b.DocumentDom, documentReferenceElement, docRef)
		}
	}
	return nil
}

// ToXmlIdentifier transforms a DSS Identifier to an XML Id type.
// Port of the protected #toXmlIdentifier.
func (b *Builder) ToXmlIdentifier(identifier model.Identifier) (string, error) {
	digest, err := spi.DSSUtilsSHA1Digest(identifier.AsXmlID())
	if err != nil {
		return "", err
	}
	return xadesBuilderIDPrefix + digest, nil
}

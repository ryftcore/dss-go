// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignature.java
// (DSS 6.5.RC1).
//
// Signature implements spi/validation.AdvancedSignature (all interface methods, real
// logic), the XAdES counterpart of cades.Signature. It embeds
// validation.DefaultAdvancedSignature and registers itself via InitDefaultAdvancedSignature,
// satisfying the DefaultAdvancedSignatureOverrides interface with real per-format logic; the
// remaining AdvancedSignature methods this base does not cover (SigningTime,
// ProductionPlace, SignaturePolicyStore, ContentType, MimeType, SignatureType,
// DAIdentifier, DataFoundUpToLevel, SignatureValue, ReferenceValidations,
// DataToBeSignedRepresentation, AddExternalTimestamp, CommitmentTypeIndications) are implemented
// directly, mirroring Java's own class exactly as CAdESSignature does for its set.
//
// # Santuario replacement
//
// org.apache.xml.security.signature.{XMLSignature,SignedInfo,Reference,Manifest} are replaced
// wholesale by internal/xmldsig (frozen). getSantuarioSignature() below builds an
// *xmldsig.XMLSignature the same way Java's private method builds an
// org.apache.xml.security.signature.XMLSignature: registering the document's ID attributes first
// (xmldom's RegisterIDs via DOMDocument.RecursiveIdBrowse - required before
// internal/xmldsig/signature.go's NewXMLSignature is called), then wiring detached-content and
// counter-signature resolvers when DetachedContents is non-empty. internal/xmldsig has no
// JCE-style registry to initialize (RSA-RIPEMD160 is dispatched directly in its sigalg.go), and
// EnforcedResolverFragment/ResolverXPointer are already covered by internal/xmldsig's own
// DefaultResolvers(). Only DSSXMLUtils.registerXAdESNamespaces() has an observable Go counterpart
// (namespace-prefix registration for XPath evaluation), reused here via init().
//
// # Deviations
//
//   - checkSignatureIntegrity's Java catch(Exception) branch appends a synthetic
//     "/ Signature/Line number/<n>" suffix derived from walking the JVM stack trace for the
//     first frame inside this class. Go has no equivalent stack-introspection API in use
//     elsewhere in this port, so checkSignatureIntegrityUnsafe's defer/recover (which reproduces
//     the broad try/catch itself, since GetReferences/ReferenceValidations/getSantuarioSignature
//     all panic on the same unchecked failure Java's DSSException represents) keeps the
//     underlying error message only, dropping the line-number suffix.
//   - SignaturePolicy()/BuildSignaturePolicy(): Java's getSignaturePolicy() override
//     downcasts super.getSignaturePolicy() to SignaturePolicy for its own callers' benefit.
//     No landed caller in this package needs that narrower type (signature_policy_store_builder.go
//     only reads .Digest(), present on the plain base signature.Policy already), so this
//     file implements only BuildSignaturePolicy() (the required override-interface entry) and
//     leaves the promoted, non-covariant Policy() getter from DefaultAdvancedSignature in
//     place - matching the "Go has no covariant returns, type-assert if you need the narrower
//     type" convention advanced_signature.go's own doc comment documents for DSSID(). A caller
//     that DOES need SignaturePolicy.TransformsDescription()'s override will only see the
//     base (always-empty) implementation through the plain *signature.Policy this file
//     hands back - the same embedding-has-no-virtual-dispatch limitation
//     xades_reference_validation.go already flags for TransformationNames.
//   - GetReferenceValidations' manifest-entry dependent validations
//     (ds:Reference[@Type=.../Manifest]) are attached through
//     model.ReferenceValidation.AddDependentValidations, the additive setter this port needs
//     because Java mutates the live List its lazy getter returns
//     ("validation.getDependentValidations().addAll(...)"), which a Go slice returned by value
//     cannot express. Before that setter landed the entries were computed and dropped, and every
//     DataObjectFormat pointing into a signed ds:Manifest failed to resolve - see the setter's
//     own doc comment in dss-model for the fixture that caught it.
//   - slf4j logging is dropped; every LOG.warn/LOG.trace/LOG.debug call site is
//     called out in the surrounding comment instead, and a Java catch-and-log-and-continue
//     becomes a Go catch-and-continue (best effort), matching every other file of this port.
package xades

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// xadesSignatureDefaultCanonicalizationMethod is the default canonicalization method used in
// SignatureDigestReference computation. Port of the private static final String
// DEFAULT_CANONICALIZATION_METHOD = CanonicalizationMethod.EXCLUSIVE.
const xadesSignatureDefaultCanonicalizationMethod = xmldsig.TransformC14NExcl

func init() {
	// Port of the static initializer's DSSXMLUtils.registerXAdESNamespaces() call; see the file
	// header for why the rest of the static block (SantuarioInitializer, JCEMapper, resolver
	// registration) has no Go counterpart here.
	DSSXMLUtilsRegisterXAdESNamespaces()
}

// Signature parses an XAdES signature structure. Note that for each signature to be
// validated a new instance of this object must be created. Port of the class XAdESSignature,
// extending validation.DefaultAdvancedSignature.
//
// serialVersionUID and java.io.Serializable are dropped (no Go counterpart).
type Signature struct {
	validation.DefaultAdvancedSignature

	// signatureElement represents the signature element container.
	signatureElement *DOMElement

	// xmldSigNamespace is the XMLDSIG namespace.
	xmldSigNamespace *common.DSSNamespace

	// xadesNamespace is the current signature xades namespace.
	xadesNamespace *common.DSSNamespace

	// xadesPath is the XAdES XPath to use.
	xadesPath definition.XAdESPath

	// disableXSWProtection defines if the XSW protection shall be disabled (false by default).
	disableXSWProtection bool

	// santuarioSignature is the cached Santuario-replacement signature. "transient" (Java) has
	// no Go counterpart since this port has no serialization.
	santuarioSignature *xmldsig.XMLSignature

	// daIdentifier is a signature identifier provided by a Driving Application, computed once.
	daIdentifier         string
	daIdentifierComputed bool

	// references contains all references found within the signature, populated on first use.
	// "transient" (Java) has no Go counterpart.
	references []*xmldsig.Reference

	// cachedCryptoVerification is this port's idempotency guard for CheckSignatureIntegrity, a
	// replacement for reading the (inaccessible, cross-package-private)
	// signatureCryptographicVerification field directly the way Java's own method does for its
	// "already computed" early-return check; see cades.Signature's identical field for the
	// precedent. It is set to the same *signature.CryptographicVerification handed to
	// SetSignatureCryptographicVerification, so the two never disagree.
	cachedCryptoVerification *signature.CryptographicVerification
}

// NewXAdESSignature is used when creating the signature. The default XPathQueryHolder is set.
// Port of the public XAdESSignature(Element) constructor.
func NewXAdESSignature(signatureElement *xmldom.Node) *Signature {
	return NewXAdESSignatureWithPathHolders(signatureElement, []definition.XAdESPath{&definition.XAdES132Path{}})
}

// NewXAdESSignatureWithPathHolders is the default constructor for XAdESSignature. Port of the
// public Signature(Element, List<XAdESPath>) constructor.
//
// Panics with the Java messages when signatureElement or xadesPathHolders is missing
// (Objects.requireNonNull).
func NewXAdESSignatureWithPathHolders(signatureElement *xmldom.Node, xadesPathHolders []definition.XAdESPath) *Signature {
	if signatureElement == nil {
		panic("Signature Element cannot be null")
	}
	if xadesPathHolders == nil {
		panic("XAdES Path holders cannot be null")
	}
	ownerDocument := NewXAdESDOMDocument(signatureElement.OwnerDocument(), xadesPathHolders)
	return NewXAdESSignatureFromDOMElement(NewXAdESDOMElement(signatureElement, ownerDocument))
}

// NewXAdESSignatureFromDOMElement is used when creating the signature. Port of the public
// Signature(DOMElement) constructor.
//
// Panics with the Java message when signatureElement is missing (Objects.requireNonNull).
func NewXAdESSignatureFromDOMElement(signatureElement *DOMElement) *Signature {
	if signatureElement == nil {
		panic("Signature Element cannot be null")
	}
	s := &Signature{
		DefaultAdvancedSignature: validation.NewDefaultAdvancedSignatureBase(),
		signatureElement:         signatureElement,
	}
	s.InitDefaultAdvancedSignature(s)
	s.initialiseSettings()
	return s
}

// SetDisableXSWProtection: NOT RECOMMENDED - this parameter allows disabling protection against
// XML Signature wrapping attacks (XSW). It disables the research by XPath expression for defined
// Type attributes. Port of setDisableXSWProtection(boolean).
func (s *Signature) SetDisableXSWProtection(disableXSWProtection bool) {
	s.disableXSWProtection = disableXSWProtection
}

// initialiseSettings is called when creating a new instance of Signature with unknown
// schema. Port of the private initialiseSettings().
func (s *Signature) initialiseSettings() {
	s.RecursiveNamespaceBrowser(s.SignatureElement())
	if s.xadesPath == nil {
		// Upstream logs "There is no suitable XAdESPaths / XAdESNamespace to manage the
		// signature. The default ones will be used.".
		s.xadesPath = &definition.XAdES132Path{}
		s.xadesNamespace = definition.XAdESNamespaceXAdES132
	}
}

// RecursiveNamespaceBrowser sets the namespace which will determine the XAdESPaths to use. The
// content of the Transform element is ignored. Port of the public recursiveNamespaceBrowser(Element).
func (s *Signature) RecursiveNamespaceBrowser(element *xmldom.Node) {
	for node := element.FirstChild; node != nil; node = node.NextSibling {
		if node.Kind != xmldom.Element {
			continue
		}
		prefix := node.Name.Prefix
		namespaceURI := node.Name.Space
		localName := node.Name.Local
		if common.XMLDSigElementTransform.IsSameTagName(localName) && common.XMLDSigElementTransform.URI() == namespaceURI {
			s.xmldSigNamespace = common.NewDSSNamespace(namespaceURI, prefix)
			continue
		} else if definition.XAdES132ElementQualifyingProperties.IsSameTagName(localName) {
			s.setXAdESPathAndNamespace(prefix, namespaceURI)
			return
		}
		s.RecursiveNamespaceBrowser(node)
	}
}

// setXAdESPathAndNamespace ports the private setXAdESPathAndNamespace(String, String).
func (s *Signature) setXAdESPathAndNamespace(prefix, namespaceURI string) {
	for _, currentXAdESPaths := range s.signatureElement.XAdESPathHolders() {
		if currentXAdESPaths.Namespace().IsSameUri(namespaceURI) {
			s.xadesPath = currentXAdESPaths
			s.xadesNamespace = common.NewDSSNamespace(namespaceURI, prefix)
		}
	}
}

// XAdESPathsHolders returns a list of used XAdESPaths. Port of the deprecated
// getXAdESPathsHolders().
func (s *Signature) XAdESPathsHolders() []definition.XAdESPath {
	return s.signatureElement.XAdESPathHolders()
}

// XAdESPaths gets the current XAdESPath. Port of getXAdESPaths().
func (s *Signature) XAdESPaths() definition.XAdESPath {
	return s.xadesPath
}

// XmldSigNamespace returns the XMLDSIG namespace. Port of getXmldSigNamespace().
func (s *Signature) XmldSigNamespace() *common.DSSNamespace {
	return s.xmldSigNamespace
}

// XadesNamespace returns the XAdES namespace. Port of getXadesNamespace().
func (s *Signature) XadesNamespace() *common.DSSNamespace {
	return s.xadesNamespace
}

// SignatureElement returns the w3c.dom encapsulated signature element. Port of getSignatureElement().
func (s *Signature) SignatureElement() *xmldom.Node {
	return s.signatureElement.Element()
}

// OwnerDocument gets the owner document. Port of getOwnerDocument().
func (s *Signature) OwnerDocument() *DOMDocument {
	return s.signatureElement.OwnerDocument()
}

// SignatureForm specifies the format of the signature. Port of getSignatureForm().
func (s *Signature) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormXAdES
}

// SignatureAlgorithm retrieves the signature algorithm (or cipher) used for generating the
// signature. Port of getSignatureAlgorithm().
//
// Panics when ds:SignedInfo/ds:SignatureMethod cannot be found: Java calls .getAttribute(...)
// directly on the XPathUtils.getElement(...) result with no null check, an unchecked
// NullPointerException for a signature without one.
func (s *Signature) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	element, err := xmlutils.XPathUtilsGetElement(s.SignatureElement(), common.XMLDSigPathSignatureMethodPath)
	if err != nil || element == nil {
		panic("xades: ds:SignedInfo/ds:SignatureMethod not found")
	}
	xmlName := element.AttrValue("", common.XMLDSigAttributeAlgorithm.AttributeName())
	signatureAlgorithm := enumerations.SignatureAlgorithmForXMLDefault(xmlName, "")
	if signatureAlgorithm == "" {
		// Upstream logs "SignatureAlgorithm '{}' is not supported!".
	}
	return signatureAlgorithm
}

// CertificateSource gets a certificate source which contains ALL certificates embedded in the
// signature. Port of getCertificateSource().
func (s *Signature) CertificateSource() *spi.SignatureCertificateSource {
	if s.OfflineCertificateSource() == nil {
		xadesCertificateSource := NewXAdESCertificateSource(s.SignatureElement(), s.xadesPath)
		s.SetOfflineCertificateSource(&xadesCertificateSource.SignatureCertificateSource)
	}
	return s.OfflineCertificateSource()
}

// CRLSource gets a CRL source which contains ALL CRLs embedded in the signature. Port of getCRLSource().
func (s *Signature) CRLSource() spi.OfflineRevocationSource[revocation.CRL] {
	if s.SignatureCRLSource() == nil {
		s.SetSignatureCRLSource(NewXAdESCRLSource(s.SignatureElement(), s.xadesPath))
	}
	return s.SignatureCRLSource()
}

// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the signature.
// Port of getOCSPSource().
func (s *Signature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP] {
	if s.SignatureOCSPSource() == nil {
		s.SetSignatureOCSPSource(NewXAdESOCSPSource(s.SignatureElement(), s.xadesPath))
	}
	return s.SignatureOCSPSource()
}

// TimestampSource gets a Signature Timestamp source which contains ALL timestamps embedded in
// the signature. Port of getTimestampSource(), covariant in Java (returns XAdESTimestampSource);
// Go callers needing the concrete type type-assert instead (xades_level_baseline_t.go and
// xades_level_x.go both do this).
func (s *Signature) TimestampSource() validation.TimestampSource {
	if s.SignatureTimestampSource() == nil {
		s.SetSignatureTimestampSource(NewXAdESTimestampSource(s))
	}
	return s.SignatureTimestampSource()
}

// SigningTime returns the signing time included within the signature, or nil. Port of getSigningTime().
func (s *Signature) SigningTime() *time.Time {
	signingTimeEl, err := xmlutils.XPathUtilsGetElement(s.SignatureElement(), s.xadesPath.SigningTimePath())
	if err != nil || signingTimeEl == nil {
		return nil
	}
	signingDate := xmlutils.DomUtilsGetDate(signingTimeEl.TextContent())
	if signingDate.IsZero() {
		return nil
	}
	return &signingDate
}

// BuildSignaturePolicy extracts a signature policy from a signature and builds the object. Port
// of the protected buildSignaturePolicy().
func (s *Signature) BuildSignaturePolicy() *signature.Policy {
	policyIdentifier, err := xmlutils.XPathUtilsGetElement(s.SignatureElement(), s.xadesPath.SignaturePolicyIdentifierPath())
	if err != nil || policyIdentifier == nil {
		return nil
	}

	policyId, err := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyId())
	if err == nil && policyId != nil {
		// Explicit policy.
		var policyUrlString string

		var qualifier enumerations.ObjectIdentifierQualifier
		qualifierString := policyId.AttrValue("", definition.XAdES132AttributeQualifier.AttributeName())
		if utils.IsStringNotBlank(qualifierString) {
			qualifier = enumerations.ObjectIdentifierQualifierFromValue(qualifierString)
		}

		policyIdString := spi.DSSUtilsObjectIdentifierValueWithQualifier(policyId.TextContent(), qualifier)
		if utils.IsStringNotBlank(policyIdString) && !spi.DSSUtilsIsUrnOid(policyIdString) && !spi.DSSUtilsIsOidCode(policyIdString) {
			policyUrlString = policyIdString
		}

		xadesSignaturePolicy := NewXAdESSignaturePolicyWithIdentifier(policyIdString)

		digestElement, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyDigestAlgAndValue())
		xadesSignaturePolicy.SetDigest(DSSXMLUtilsGetDigestAndValue(digestElement))

		policyUrl, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicySPURI())
		if policyUrl != nil {
			policyUrlString = spi.DSSUtilsTrimWhitespacesAndNewlines(policyUrl.TextContent())
		}
		xadesSignaturePolicy.SetURI(policyUrlString)

		spUserNotice, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicySPUserNotice())
		if spUserNotice != nil {
			if userNotice := s.buildSPUserNotice(spUserNotice); userNotice != nil {
				xadesSignaturePolicy.SetUserNotice(userNotice)
			}
		}

		if docSpecPath := s.xadesPath.CurrentSignaturePolicySPDocSpecification(); docSpecPath != nil {
			spDocSpecification, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, docSpecPath)
			if spDocSpecification != nil {
				xadesSignaturePolicy.SetDocSpecification(s.buildSpDocSpecification(spDocSpecification))
			}
		}

		policyDescription, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyDescription())
		if policyDescription != nil && utils.IsStringNotEmpty(policyDescription.TextContent()) {
			xadesSignaturePolicy.SetDescription(policyDescription.TextContent())
		}

		docRefsNode, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyDocumentationReferences())
		if docRefsNode != nil {
			xadesSignaturePolicy.SetDocumentationReferences(s.getDocumentationReferences(docRefsNode))
		}

		transformsNode, _ := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyTransforms())
		if transformsNode != nil {
			xadesSignaturePolicy.SetTransforms(transformsNode)
			xadesSignaturePolicy.SetHashAsInTechnicalSpecification(s.isHashComputationAsInPolicySpecification(transformsNode))
		}

		return &xadesSignaturePolicy.Policy
	}

	// Implicit policy.
	signaturePolicyImplied, err := xmlutils.XPathUtilsGetElement(policyIdentifier, s.xadesPath.CurrentSignaturePolicyImplied())
	if err == nil && signaturePolicyImplied != nil {
		return &NewXAdESSignaturePolicy().Policy
	}
	return nil
}

// buildSPUserNotice ports the private buildSPUserNotice(Element). A malformed value is caught
// and dropped (Java logs "Unable to build SPUserNotice qualifier. Reason : {}" and returns nil).
func (s *Signature) buildSPUserNotice(spUserNoticeElement *xmldom.Node) (result *model.UserNotice) {
	defer func() {
		if recover() != nil {
			result = nil
		}
	}()
	userNotice := model.NewUserNotice()

	organization, _ := xmlutils.XPathUtilsGetElement(spUserNoticeElement, s.xadesPath.CurrentSPUserNoticeNoticeRefOrganization())
	if organization != nil {
		userNotice.SetOrganization(organization.TextContent())
	}
	noticeNumbers, _ := xmlutils.XPathUtilsGetElement(spUserNoticeElement, s.xadesPath.CurrentSPUserNoticeNoticeRefNoticeNumbers())
	if noticeNumbers != nil && noticeNumbers.FirstChild != nil {
		var noticeNumbersList []int
		for child := noticeNumbers.FirstChild; child != nil; child = child.NextSibling {
			if child.Kind == xmldom.Element && definition.XAdES132ElementInt.IsSameTagName(child.Name.Local) {
				n, err := strconv.Atoi(child.TextContent())
				if err != nil {
					panic(err)
				}
				noticeNumbersList = append(noticeNumbersList, n)
			}
		}
		userNotice.SetNoticeNumbers(noticeNumbersList...)
	}
	explicitText, _ := xmlutils.XPathUtilsGetElement(spUserNoticeElement, s.xadesPath.CurrentSPUserNoticeExplicitText())
	if explicitText != nil {
		userNotice.SetExplicitText(explicitText.TextContent())
	}
	return userNotice
}

// isHashComputationAsInPolicySpecification ports the private
// isHashComputationAsInPolicySpecification(Element).
func (s *Signature) isHashComputationAsInPolicySpecification(transforms *xmldom.Node) bool {
	if transforms != nil && transforms.FirstChild != nil {
		transformList, err := xmlutils.XPathUtilsGetNodeList(transforms, common.XMLDSigPathTransformPath)
		if err == nil && len(transformList) == 1 {
			algorithm := transformList[0].AttrValue("", common.XMLDSigAttributeAlgorithm.AttributeName())
			return DSSXMLUtilsSPDocDigestAsInSpecificationAlgorithmURI == algorithm
		}
	}
	return false
}

// buildSpDocSpecification ports the private buildSpDocSpecification(Element).
func (s *Signature) buildSpDocSpecification(spDocSpecificationElement *xmldom.Node) *model.SpDocSpecification {
	spDocSpec := model.NewSpDocSpecification()

	identifierElement, _ := xmlutils.XPathUtilsGetElement(spDocSpecificationElement, s.xadesPath.CurrentIdentifier())
	if identifierElement != nil {
		spDocSpecId := identifierElement.TextContent()

		var qualifier enumerations.ObjectIdentifierQualifier
		qualifierString := identifierElement.AttrValue("", definition.XAdES132AttributeQualifier.AttributeName())
		if utils.IsStringNotBlank(qualifierString) {
			qualifier = enumerations.ObjectIdentifierQualifierFromValue(qualifierString)
			spDocSpec.SetQualifier(qualifier)
		}
		spDocSpec.SetId(spi.DSSUtilsObjectIdentifierValueWithQualifier(spDocSpecId, qualifier))
	}

	description, _ := xmlutils.XPathUtilsGetValue(spDocSpecificationElement, s.xadesPath.CurrentDescription())
	if utils.IsStringNotBlank(description) {
		spDocSpec.SetDescription(description)
	}

	if docRefElementsPath := s.xadesPath.CurrentDocumentationReferenceElements(); docRefElementsPath != nil {
		documentReferenceList, err := xmlutils.XPathUtilsGetNodeList(spDocSpecificationElement, docRefElementsPath)
		if err == nil && len(documentReferenceList) > 0 {
			documentationReferences := make([]string, len(documentReferenceList))
			for i, node := range documentReferenceList {
				documentationReferences[i] = node.TextContent()
			}
			spDocSpec.SetDocumentationReferences(documentationReferences...)
		}
	}
	return spDocSpec
}

// ProductionPlace returns information about the place where the signature was
// generated. Port of getSignatureProductionPlace().
func (s *Signature) SignatureProductionPlace() *signature.ProductionPlace {
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.SignatureProductionPlacePath())
	if err != nil {
		nodeList = nil
	}
	if len(nodeList) == 0 {
		if v2Path := s.xadesPath.SignatureProductionPlaceV2Path(); v2Path != nil {
			nodeList, err = xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), v2Path)
			if err != nil {
				nodeList = nil
			}
		}
	}
	if len(nodeList) == 0 {
		return nil
	}

	signatureProductionPlace := signature.NewSignatureProductionPlace()
	for item := nodeList[0].FirstChild; item != nil; item = item.NextSibling {
		name := item.Name.Local
		nodeValue := item.TextContent()
		switch {
		case definition.XAdES132ElementCity.IsSameTagName(name):
			signatureProductionPlace.SetCity(nodeValue)
		case definition.XAdES132ElementStateOrProvince.IsSameTagName(name):
			signatureProductionPlace.SetStateOrProvince(nodeValue)
		case definition.XAdES132ElementPostalCode.IsSameTagName(name):
			signatureProductionPlace.SetPostalCode(nodeValue)
		case definition.XAdES132ElementCountryName.IsSameTagName(name):
			signatureProductionPlace.SetCountryName(nodeValue)
		case definition.XAdES132ElementStreetAddress.IsSameTagName(name):
			signatureProductionPlace.SetStreetAddress(nodeValue)
		}
	}
	return signatureProductionPlace
}

// SignaturePolicyStore returns the Signature Policy Store from the signature. Port of getSignaturePolicyStore().
func (s *Signature) SignaturePolicyStore() *model.SignaturePolicyStore {
	signaturePolicyStorePath := s.xadesPath.SignaturePolicyStorePath()
	if signaturePolicyStorePath == nil {
		return nil
	}
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), signaturePolicyStorePath)
	if err != nil || len(nodeList) == 0 {
		return nil
	}

	sps := model.NewSignaturePolicyStore()
	signaturePolicyStoreElement := nodeList[0]
	id := signaturePolicyStoreElement.AttrValue("", common.XMLDSigAttributeID.AttributeName())
	if utils.IsStringNotEmpty(id) {
		sps.SetId(id)
	}

	var spDocSpec *model.SpDocSpecification
	if currentSPDocSpecificationPath := s.xadesPath.CurrentSPDocSpecification(); currentSPDocSpecificationPath != nil {
		spDocSpecificationElement, _ := xmlutils.XPathUtilsGetElement(signaturePolicyStoreElement, currentSPDocSpecificationPath)
		if spDocSpecificationElement != nil {
			spDocSpec = s.buildSpDocSpecification(spDocSpecificationElement)
		}
	}
	sps.SetSpDocSpecification(spDocSpec)

	if currentSignaturePolicyDocumentPath := s.xadesPath.CurrentSignaturePolicyDocument(); currentSignaturePolicyDocumentPath != nil {
		spDocB64, _ := xmlutils.XPathUtilsGetValue(signaturePolicyStoreElement, currentSignaturePolicyDocumentPath)
		if utils.IsStringNotEmpty(spDocB64) && utils.IsBase64Encoded(spDocB64) {
			sps.SetSignaturePolicyContent(model.NewInMemoryDocument(utils.FromBase64(spDocB64)))
		}
	}

	if currentSigPolDocLocalURI := s.xadesPath.CurrentSigPolDocLocalURI(); currentSigPolDocLocalURI != nil {
		sigPolDocLocalURI, _ := xmlutils.XPathUtilsGetValue(signaturePolicyStoreElement, currentSigPolDocLocalURI)
		if utils.IsStringNotEmpty(sigPolDocLocalURI) {
			sps.SetSigPolDocLocalURI(sigPolDocLocalURI)
		}
	}
	return sps
}

// SignedAssertions returns the list of embedded signed assertions. Port of getSignedAssertions().
func (s *Signature) SignedAssertions() []*signature.SignerRole {
	result := []*signature.SignerRole{}
	signedAssertionPath := s.xadesPath.SignedAssertionPath()
	if signedAssertionPath != nil {
		nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), signedAssertionPath)
		if err == nil {
			for _, node := range nodeList {
				xmlString, xerr := xmlutils.DomUtilsXmlToString(node.FirstChild)
				if xerr != nil {
					panic(xerr)
				}
				result = append(result, signature.NewSignerRole(xmlString, enumerations.EndorsementTypeSigned))
			}
		}
	}
	return result
}

// ClaimedSignerRoles returns the claimed roles of the signer. Port of getClaimedSignerRoles().
func (s *Signature) ClaimedSignerRoles() []*signature.SignerRole {
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.ClaimedRolePath())
	if err != nil {
		nodeList = nil
	}
	if len(nodeList) == 0 {
		if claimedRoleV2Path := s.xadesPath.ClaimedRoleV2Path(); claimedRoleV2Path != nil {
			nodeList, err = xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), claimedRoleV2Path)
			if err != nil || len(nodeList) == 0 {
				return []*signature.SignerRole{}
			}
		}
	}
	claimedRoles := []*signature.SignerRole{}
	for _, node := range nodeList {
		claimedRoles = append(claimedRoles, signature.NewSignerRole(node.TextContent(), enumerations.EndorsementTypeClaimed))
	}
	return claimedRoles
}

// CertifiedSignerRoles returns the certified roles of the signer. Port of getCertifiedSignerRoles().
func (s *Signature) CertifiedSignerRoles() []*signature.SignerRole {
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.CertifiedRolePath())
	if err != nil {
		nodeList = nil
	}
	if len(nodeList) == 0 {
		if certifiedRoleV2Path := s.xadesPath.CertifiedRoleV2Path(); certifiedRoleV2Path != nil {
			nodeList, err = xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), certifiedRoleV2Path)
			if err != nil || len(nodeList) == 0 {
				return []*signature.SignerRole{}
			}
		}
	}
	certifiedRoles := []*signature.SignerRole{}
	for _, node := range nodeList {
		certifiedRoles = append(certifiedRoles, signature.NewSignerRole(node.TextContent(), enumerations.EndorsementTypeCertified))
	}
	return certifiedRoles
}

// ContentType returns the value of the signed attribute content-type. Port of getContentType().
func (s *Signature) ContentType() string {
	allContentTypes, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.DataObjectFormatObjectIdentifier())
	if err == nil {
		for _, node := range allContentTypes {
			if node.Kind == xmldom.Element {
				// TODO (matches upstream): returns the first one.
				return node.TextContent()
			}
		}
	}
	return ""
}

// MimeType returns the value of the signed attribute mime-type. Port of getMimeType().
func (s *Signature) MimeType() string {
	allMimeTypes, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.DataObjectFormatMimeType())
	if err == nil {
		for _, node := range allMimeTypes {
			if node.Kind == xmldom.Element {
				// TODO (matches upstream): returns the first one.
				return node.TextContent()
			}
		}
	}
	return ""
}

// SignatureType returns the value of the signature type protected header. Not supported for
// XAdES. Port of getSignatureType().
func (s *Signature) SignatureType() string {
	return ""
}

// SignatureValueBase64 returns a base64 SignatureValue. Port of the public getSignatureValueBase64().
func (s *Signature) SignatureValueBase64() string {
	signatureValueElement, err := xmlutils.XPathUtilsGetElement(s.SignatureElement(), common.XMLDSigPathSignatureValuePath)
	if err != nil || signatureValueElement == nil {
		return ""
	}
	return signatureValueElement.TextContent()
}

// SignatureValue returns the digital signature value. Port of getSignatureValue().
func (s *Signature) SignatureValue() []byte {
	signatureValueBase64 := s.SignatureValueBase64()
	if signatureValueBase64 != "" && utils.IsBase64Encoded(signatureValueBase64) {
		return utils.FromBase64(signatureValueBase64)
	}
	// Upstream logs "The signature value is not represented by a base64-encoded string! ...".
	return nil
}

// SignatureValueId returns Id of the ds:SignatureValue element. Port of getSignatureValueId().
func (s *Signature) SignatureValueId() string {
	value, err := xmlutils.XPathUtilsGetValue(s.SignatureElement(), common.XMLDSigPathSignatureValueIDPath)
	if err != nil {
		return ""
	}
	return value
}

// Objects returns the list of ds:Object elements for the current signature element. Port of the
// public getObjects().
func (s *Signature) Objects() []*xmldom.Node {
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), common.XMLDSigPathObjectPath)
	if err != nil {
		return nil
	}
	return nodeList
}

// AddExternalTimestamp allows adding an external timestamp. Not supported for XAdES. Port of
// addExternalTimestamp(TimestampToken).
func (s *Signature) AddExternalTimestamp(timestamp *validation.TimestampToken) {
	panic("The action is not supported for XAdES!")
}

// BaselineRequirementsChecker returns the cached instance of the BaselineRequirementsChecker.
// Port of the protected covariant-return override getBaselineRequirementsChecker().
func (s *Signature) BaselineRequirementsChecker() *BaselineRequirementsChecker {
	return s.DefaultAdvancedSignature.BaselineRequirementsChecker().(*BaselineRequirementsChecker)
}

// CreateBaselineRequirementsChecker instantiates a BaselineRequirementsChecker according to the
// signature format. Port of the protected createBaselineRequirementsChecker(CertificateVerifier).
func (s *Signature) CreateBaselineRequirementsChecker(certificateVerifier validation.CertificateVerifier) validation.BaselineRequirementsCheckerContract {
	return NewXAdESBaselineRequirementsChecker(s, certificateVerifier)
}

// CheckSignatureIntegrity verifies the signature integrity; checks if the signed content has not
// been tampered with. Port of checkSignatureIntegrity().
//
// The early "already computed" return reads s.cachedCryptoVerification rather than the base's
// own (inaccessible, cross-package-private) signatureCryptographicVerification field; see that
// field's doc comment.
func (s *Signature) CheckSignatureIntegrity() {
	if s.cachedCryptoVerification != nil {
		return
	}
	verification := signature.NewSignatureCryptographicVerification()
	s.cachedCryptoVerification = verification
	s.SetSignatureCryptographicVerification(verification)

	if err := s.checkSignatureIntegrityUnsafe(verification); err != nil {
		// Java appends a synthetic "/ XAdESSignature/Line number/<n>" suffix derived from the
		// stack trace; see the file header DEVIATION note for why this port keeps only the
		// message.
		verification.SetErrorMessage(err.Error())
	}
}

// checkSignatureIntegrityUnsafe runs the try body of checkSignatureIntegrity(); its defer/
// recover reproduces the surrounding catch(Exception), since GetSantuarioSignature/References/
// ReferenceValidations all panic on the same unchecked failures Java's DSSException represents.
func (s *Signature) checkSignatureIntegrityUnsafe(verification *signature.CryptographicVerification) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	currentSantuarioSignature, sigErr := s.getSantuarioSignature()
	if sigErr != nil {
		return sigErr
	}
	candidatesForSigningCertificate := s.CandidatesForSigningCertificate()

	signingCertificateValidator := NewXAdESSignatureIntegrityValidator(currentSantuarioSignature)
	certificateValidity := signingCertificateValidator.Validate(candidatesForSigningCertificate)
	if certificateValidity != nil {
		if err := candidatesForSigningCertificate.SetTheCertificateValidity(certificateValidity); err != nil {
			return err
		}
	}
	verification.SetErrorMessages(signingCertificateValidator.ErrorMessages())

	allReferenceDataFound := true
	allReferenceDataIntact := true
	for _, referenceValidation := range s.ReferenceValidations() {
		allReferenceDataFound = allReferenceDataFound && referenceValidation.IsFound()
		allReferenceDataIntact = allReferenceDataIntact && referenceValidation.IsIntact()
	}
	verification.SetReferenceDataFound(allReferenceDataFound)
	verification.SetReferenceDataIntact(allReferenceDataIntact)
	verification.SetSignatureIntact(certificateValidity != nil)
	return nil
}

// ReferenceValidations returns individual validation for each reference. Port of
// getReferenceValidations().
func (s *Signature) ReferenceValidations() []*model.ReferenceValidation {
	if s.CachedReferenceValidations() != nil {
		return s.CachedReferenceValidations()
	}

	referenceValidations := []*model.ReferenceValidation{}
	currentSantuarioSignature := s.mustGetSantuarioSignature()
	atLeastOneReferenceElementFound := false

	santuarioReferences := s.References()
	for _, reference := range santuarioReferences {
		// Named refValidation, not "validation": that identifier would shadow the imported
		// spi/validation package for the rest of this loop body.
		refValidation := NewXAdESReferenceValidation(reference)
		refValidation.SetType(enumerations.DigestMatcherTypeReference)
		referenceValidations = append(referenceValidations, &refValidation.ReferenceValidation)

		found := false
		intact := false
		func() {
			// Java wraps the whole body in try { ... } catch (Exception e) { LOG.warn(...) }: a
			// malformed reference is skipped with found/intact left at their zero values.
			defer func() { recover() }()

			digest := DSSXMLUtilsGetReferenceDigest(reference)
			refValidation.SetDigest(digest)

			found = DSSXMLUtilsIsAbleToDeReferenceContent(reference)

			uri := refValidation.Uri()
			isDuplicated := DSSXMLUtilsIsReferencedContentAmbiguous(s.OwnerDocument().Document(), uri)
			refValidation.SetDuplicated(isDuplicated)

			isElementReference := xmlutils.DomUtilsIsElementReference(uri)

			switch {
			case isElementReference && DSSXMLUtilsIsSignedProperties(reference, s.xadesPath):
				refValidation.SetType(enumerations.DigestMatcherTypeSignedProperties)
				found = found && (s.disableXSWProtection || s.findSignedPropertiesById(uri))

			case xmlutils.DomUtilsIsXPointerQuery(uri):
				refValidation.SetType(enumerations.DigestMatcherTypeXPointer)
				// found is checked in the reference validation.

			case isElementReference && DSSXMLUtilsIsKeyInfoReference(reference, currentSantuarioSignature.Element()):
				refValidation.SetType(enumerations.DigestMatcherTypeKeyInfo)
				found = true // checked prior, inside "isKeyInfoReference".

			case isElementReference && DSSXMLUtilsIsSignaturePropertiesReference(reference, currentSantuarioSignature.Element()):
				refValidation.SetType(enumerations.DigestMatcherTypeSignatureProperties)
				found = true // Id is verified inside "isSignaturePropertiesReference".

			case isElementReference && reference.TypeIsReferenceToObject():
				refValidation.SetType(enumerations.DigestMatcherTypeObject)
				found = found && (s.disableXSWProtection || s.findObjectById(uri))

			case isElementReference && reference.TypeIsReferenceToManifest():
				refValidation.SetType(enumerations.DigestMatcherTypeManifest)
				manifestElement := DSSXMLUtilsGetManifestById(s.SignatureElement(), uri)
				found = found && (s.disableXSWProtection || manifestElement != nil)
				if manifestElement != nil {
					refValidation.AddDependentValidations(s.getManifestReferences(manifestElement)...)
				}

			case DSSXMLUtilsIsCounterSignatureReference(reference, s):
				refValidation.SetType(enumerations.DigestMatcherTypeCounterSignature)
				// found is checked in the reference validation.
				if masterSignature, ok := s.MasterSignature().(*Signature); ok {
					referenceValidations = append(referenceValidations,
						s.getCounterSignatureReferenceValidation(reference, masterSignature))
				}
			}

			if found && !isDuplicated {
				intact, _ = reference.Verify()
			}
		}()

		switch refValidation.Type() {
		case enumerations.DigestMatcherTypeReference, enumerations.DigestMatcherTypeObject,
			enumerations.DigestMatcherTypeManifest, enumerations.DigestMatcherTypeXPointer,
			enumerations.DigestMatcherTypeCounterSignature:
			atLeastOneReferenceElementFound = true
		}

		refValidation.SetFound(found)
		refValidation.SetIntact(intact)
	}

	// If at least one reference is not found, add an empty referenceValidation.
	if !atLeastOneReferenceElementFound {
		referenceValidations = append(referenceValidations, s.notFound(enumerations.DigestMatcherTypeReference))
	}

	// Upstream logs "Not all references were validated!" when
	// len(referenceValidations) < len(santuarioReferences).

	if s.IsCounterSignature() {
		referenceValidations = s.ensureCounterSignatureReferenceFound(referenceValidations)
	}

	s.SetCachedReferenceValidations(referenceValidations)
	return referenceValidations
}

// XAdESReferenceValidations narrows ReferenceValidations() down to the elements that were built
// as *ReferenceValidation (via the registry xades_reference_validation.go maintains for
// exactly this purpose - Go has no instanceof/covariant-return equivalent through embedding).
// Ports the effect of Java's `getReferenceValidations().stream().filter(rv -> rv instanceof
// ReferenceValidation)...` idiom used by SignatureScopeFinder,
// TimestampScopeFinder and TimestampSource. Not every element qualifies: e.g. notFound()
// and getCounterSignatureReferenceValidation() below build plain *model.ReferenceValidation
// values that were never registered.
func (s *Signature) XAdESReferenceValidations() []*ReferenceValidation {
	referenceValidations := s.ReferenceValidations()
	result := make([]*ReferenceValidation, 0, len(referenceValidations))
	for _, rv := range referenceValidations {
		if xrv, ok := ReferenceValidationFor(rv); ok {
			result = append(result, xrv)
		}
	}
	return result
}

// getCounterSignatureReferenceValidation ports the private
// getCounterSignatureReferenceValidation(Reference, Signature).
func (s *Signature) getCounterSignatureReferenceValidation(counterSignatureReference *xmldsig.Reference, masterSignature *Signature) *model.ReferenceValidation {
	referenceValidation := model.NewReferenceValidation()
	referenceValidation.SetType(enumerations.DigestMatcherTypeCounterSignedSignatureValue)

	masterSignatureValueBase64 := masterSignature.SignatureValueBase64()
	if utils.IsStringNotEmpty(masterSignatureValueBase64) {
		referenceValidation.SetFound(true)

		func() {
			defer func() { recover() }()
			out, err := counterSignatureReference.ContentsAfterTransformation()
			if err != nil {
				return
			}
			referencedBytes, err := out.Bytes()
			if err != nil {
				return
			}
			document, err := xmlutils.DomUtilsBuildDOMFromBytes(referencedBytes)
			if err != nil || document == nil {
				return
			}
			referencedElement := document.DocumentElement()
			if referencedElement != nil && common.XMLDSigElementSignatureValue.IsSameTagName(referencedElement.Name.Local) {
				referencedSignatureValueBase64 := referencedElement.TextContent()
				intact := utils.AreStringsEqual(masterSignatureValueBase64, referencedSignatureValueBase64)
				// Upstream logs a warning when !intact ("...does not match the master
				// signature's ds:SignatureValue content!").
				referenceValidation.SetIntact(intact)
			}
			// Upstream logs "The counter signature reference does not result to a
			// ds:SignatureValue element!" otherwise.
		}()
	}
	// Upstream logs "Master signature's ds:SignatureValue element does not contain data!" otherwise.

	return referenceValidation
}

// ensureCounterSignatureReferenceFound ports the private
// ensureCounterSignatureReferenceFound(List<ReferenceValidation>).
func (s *Signature) ensureCounterSignatureReferenceFound(referenceValidations []*model.ReferenceValidation) []*model.ReferenceValidation {
	for _, r := range referenceValidations {
		if r.Type() == enumerations.DigestMatcherTypeCounterSignedSignatureValue {
			return referenceValidations
		}
	}
	// Upstream logs "No Reference covering the master signature's SignatureValue has been found!".
	referenceValidation := model.NewReferenceValidation()
	referenceValidation.SetType(enumerations.DigestMatcherTypeCounterSignedSignatureValue)
	referenceValidation.SetFound(s.MasterSignature() != nil)
	referenceValidation.SetIntact(false)
	return append(referenceValidations, referenceValidation)
}

// BuildSignatureDigestReference builds a new DigestReference according to the
// applicable signature format rules.
//
// TS 119 442 - V1.1.1 - Electronic Signatures and Infrastructures (ESI), ch. 5.1.4.2.1.3 XML
// component: in case of XAdES signatures, the input of the digest value computation shall be the
// result of applying the canonicalization algorithm identified within the
// CanonicalizationMethod child element's value to the corresponding ds:Signature element and its
// contents. The canonicalization shall be computed keeping this ds:Signature element as a
// descendant of the XML root element, without detaching it. Port of the protected
// buildSignatureDigestReference(DigestAlgorithm).
func (s *Signature) BuildSignatureDigestReference(digestAlgorithm enumerations.DigestAlgorithm) *signature.DigestReference {
	messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedNode(s.SignatureElement(), digestAlgorithm, xadesSignatureDefaultCanonicalizationMethod)
	if err != nil {
		panic(err)
	}
	return signature.NewSignatureDigestReferenceWithCanonicalization(xadesSignatureDefaultCanonicalizationMethod, messageDigest.Digest)
}

// DataToBeSignedRepresentation returns the DTBSR. Port of getDataToBeSignedRepresentation().
func (s *Signature) DataToBeSignedRepresentation() model.Digest {
	digestAlgorithm := s.DigestAlgorithm()
	if digestAlgorithm == "" {
		// Upstream logs "DigestAlgorithm is not found! Unable to compute DTBSR.".
		return model.Digest{}
	}
	signedInfoElement := s.SignedInfo()
	if signedInfoElement == nil {
		// Upstream logs "SignedInfo element is not found! Unable to compute DTBSR.".
		return model.Digest{}
	}
	canonicalizationMethod, err := xmlutils.XPathUtilsGetValue(signedInfoElement, common.XMLDSigPathCanonicalizationAlgorithmPath)
	if err != nil || utils.IsStringEmpty(canonicalizationMethod) {
		// Upstream logs "Canonicalization method is not present in SignedInfo element! Unable
		// to compute DTBSR.".
		return model.Digest{}
	}
	messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedNode(signedInfoElement, digestAlgorithm, canonicalizationMethod)
	if err != nil {
		return model.Digest{}
	}
	return messageDigest.Digest
}

// SignedInfo returns the ds:SignedInfo element. Port of the public getSignedInfo().
func (s *Signature) SignedInfo() *xmldom.Node {
	element, err := xmlutils.XPathUtilsGetElement(s.SignatureElement(), common.XMLDSigPathSignedInfoPath)
	if err != nil {
		// Upstream logs "Unable to extract ds:SignedInfo element! Reason : %s.".
		return nil
	}
	return element
}

// getManifestReferences returns a list of all references contained in the given manifest. Port
// of the private getManifestReferences(Element).
func (s *Signature) getManifestReferences(manifestElement *xmldom.Node) []*model.ReferenceValidation {
	mv := NewManifestValidator(manifestElement, s.DetachedContents())
	return mv.Validate()
}

// findSignedPropertiesById ports the private findSignedPropertiesById(String).
func (s *Signature) findSignedPropertiesById(uri string) bool {
	return s.getSignedPropertiesById(uri) != nil
}

// getSignedPropertiesById ports the private getSignedPropertiesById(String).
func (s *Signature) getSignedPropertiesById(uri string) *xmldom.Node {
	if utils.IsStringNotBlank(uri) {
		return xmlutils.XPathUtilsGetElementByIdWithQuery(s.SignatureElement(), s.xadesPath.SignedPropertiesPath(), uri)
	}
	return nil
}

// findObjectById ports the private findObjectById(String).
func (s *Signature) findObjectById(uri string) bool {
	return s.ObjectById(uri) != nil
}

// ObjectById gets ds:Object by its Id. Port of the public getObjectById(String).
func (s *Signature) ObjectById(id string) *xmldom.Node {
	if utils.IsStringNotBlank(id) {
		return xmlutils.XPathUtilsGetElementByIdWithQuery(s.SignatureElement(), common.XMLDSigPathObjectPath, id)
	}
	return nil
}

// ManifestById gets ds:Manifest by its Id. Port of the public getManifestById(String).
func (s *Signature) ManifestById(id string) *xmldom.Node {
	if utils.IsStringNotBlank(id) {
		return xmlutils.XPathUtilsGetElementByIdWithQuery(s.SignatureElement(), common.XMLDSigPathManifestPath, id)
	}
	return nil
}

// notFound ports the private notFound(DigestMatcherType).
func (s *Signature) notFound(typ enumerations.DigestMatcherType) *model.ReferenceValidation {
	refValidation := model.NewReferenceValidation()
	refValidation.SetType(typ)
	refValidation.SetFound(false)
	return refValidation
}

// getSantuarioSignature ports the private getSantuarioSignature(), lazily building and caching
// the Santuario-replacement *xmldsig.XMLSignature.
//
// The Java method has no declared throws (its XMLSecurityException catch rethrows an unchecked
// DSSException), so every direct caller in Java relies on unchecked propagation; this port
// returns the error explicitly to the one caller that needs to catch it faithfully
// (checkSignatureIntegrityUnsafe) and panics through mustGetSantuarioSignature everywhere else,
// reproducing that unchecked propagation.
func (s *Signature) getSantuarioSignature() (*xmldsig.XMLSignature, error) {
	if s.santuarioSignature != nil {
		return s.santuarioSignature, nil
	}
	// This enables Ids for cryptographic validation.
	s.OwnerDocument().RecursiveIdBrowse()

	// Secure validation disabled to support all signature algorithms.
	sig, err := xmldsig.NewXMLSignature(s.SignatureElement(), nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize Santuario XMLSignature. Reason : %s", err.Error())
	}
	s.santuarioSignature = sig

	if utils.IsCollectionNotEmpty(s.DetachedContents()) {
		s.initDetachedSignatureResolvers(s.DetachedContents())
		s.initCounterSignatureResolver(s.DetachedContents())
	}
	return s.santuarioSignature, nil
}

// mustGetSantuarioSignature calls getSantuarioSignature and panics on error, reproducing the
// unchecked DSSException every direct Java caller (with no try/catch of its own) relies on.
func (s *Signature) mustGetSantuarioSignature() *xmldsig.XMLSignature {
	sig, err := s.getSantuarioSignature()
	if err != nil {
		panic(err)
	}
	return sig
}

// initDetachedSignatureResolvers ports the private initDetachedSignatureResolvers(List<DSSDocument>).
func (s *Signature) initDetachedSignatureResolvers(detachedContents []model.DSSDocument) {
	signedInfoElement := s.SignedInfo()
	if signedInfoElement == nil {
		return
	}
	xmlSignature := s.mustGetSantuarioSignature()
	for _, digestAlgorithm := range DSSXMLUtilsGetReferenceDigestAlgos(signedInfoElement) {
		xmlSignature.AddResourceResolver(&xmldsig.DetachedSignatureResolver{
			Documents:       detachedContents,
			DigestAlgorithm: digestAlgorithm,
		})
	}
}

// initCounterSignatureResolver is used for a counter signature extension only. Port of the
// private initCounterSignatureResolver(List<DSSDocument>).
func (s *Signature) initCounterSignatureResolver(detachedContents []model.DSSDocument) {
	signedInfoElement := s.SignedInfo()
	if signedInfoElement == nil {
		return
	}
	xmlSignature := s.mustGetSantuarioSignature()
	for _, typ := range DSSXMLUtilsGetReferenceTypes(signedInfoElement) {
		if s.xadesPath.CounterSignatureUri() == typ {
			for _, document := range detachedContents {
				// Only one SignatureValue document shall be provided.
				if s.isDetachedSignatureValueDocument(document) {
					xmlSignature.AddResourceResolver(NewCounterSignatureResolver(document))
					break
				}
			}
		}
	}
}

// isDetachedSignatureValueDocument ports the private isDetachedSignatureValueDocument(DSSDocument).
func (s *Signature) isDetachedSignatureValueDocument(detachedContent model.DSSDocument) (result bool) {
	defer func() {
		if recover() != nil {
			result = false
		}
	}()
	if !xmlutils.DomUtilsIsDOM(detachedContent) {
		return false
	}
	document, err := xmlutils.DomUtilsBuildDOMFromDocument(detachedContent)
	if err != nil || document == nil {
		return false
	}
	node := document.FirstChild
	if node == nil {
		return false
	}
	return common.XMLDSigElementSignatureValue.TagName() == node.Name.Local
}

// CounterSignatures retrieves the potential countersignatures embedded in the XAdES signature
// document. From ETSI TS 101 903 v1.4.2:
//
// 7.2.4.1 Countersignature identifier in Type attribute of ds:Reference: a XAdES signature
// containing a ds:Reference element whose Type attribute has value
// "http://uri.etsi.org/01903#CountersignedSignature" will indicate that it is, in fact, a
// countersignature of the signature referenced by this element.
//
// 7.2.4.2 Enveloped countersignatures: the CounterSignature element is an unsigned property that
// qualifies the signature. A XAdES signature MAY have more than one CounterSignature properties.
// As indicated by its name, it contains one countersignature of the qualified signature.
//
// Port of getCounterSignatures().
func (s *Signature) CounterSignatures() []validation.AdvancedSignature {
	if s.CachedCounterSignatures() != nil {
		return s.CachedCounterSignatures()
	}

	counterSignatures := []validation.AdvancedSignature{}

	// See ETSI TS 101 903 V1.4.2 (2010-12) pp. 38/39/40.
	counterSignaturesElements, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.CounterSignaturePath())
	if err == nil {
		for _, el := range counterSignaturesElements {
			counterSignature := DSSXMLUtilsCreateCounterSignature(el, s)
			if counterSignature != nil {
				counterSignatures = append(counterSignatures, counterSignature)
			}
		}
	}
	s.SetCachedCounterSignatures(counterSignatures)
	return counterSignatures
}

// SignatureIdentifierBuilder returns a builder to define and build a signature Id. Port of the
// protected getSignatureIdentifierBuilder().
func (s *Signature) SignatureIdentifierBuilder() validation.SignatureIdentifierBuilder {
	return NewXAdESSignatureIdentifierBuilder(s)
}

// DAIdentifier returns an identifier provided by the Driving Application (DA). Port of getDAIdentifier().
func (s *Signature) DAIdentifier() string {
	if !s.daIdentifierComputed {
		s.daIdentifier = DSSXMLUtilsGetIDIdentifier(s.SignatureElement())
		s.daIdentifierComputed = true
	}
	return s.daIdentifier
}

// UnsignedSignatureProperties retrieves the name of each node found under the
// UnsignedSignatureProperties element. Port of the public getUnsignedSignatureProperties().
func (s *Signature) UnsignedSignatureProperties() []string {
	names, err := xmlutils.XPathUtilsGetChildrenNames(s.SignatureElement(), s.xadesPath.UnsignedSignaturePropertiesPath())
	if err != nil {
		return nil
	}
	return names
}

// SignedSignatureProperties retrieves the name of each node found under the
// SignedSignatureProperties element. Port of the public getSignedSignatureProperties().
func (s *Signature) SignedSignatureProperties() []string {
	names, err := xmlutils.XPathUtilsGetChildrenNames(s.SignatureElement(), s.xadesPath.SignedSignaturePropertiesPath())
	if err != nil {
		return nil
	}
	return names
}

// SignedProperties retrieves the name of each node found under the SignedProperties element.
// Port of the public getSignedProperties().
func (s *Signature) SignedProperties() []string {
	names, err := xmlutils.XPathUtilsGetChildrenNames(s.SignatureElement(), s.xadesPath.SignedPropertiesPath())
	if err != nil {
		return nil
	}
	return names
}

// UnsignedProperties retrieves the name of each node found under the UnsignedProperties element.
// Port of the public getUnsignedProperties().
func (s *Signature) UnsignedProperties() []string {
	names, err := xmlutils.XPathUtilsGetChildrenNames(s.SignatureElement(), s.xadesPath.UnsignedPropertiesPath())
	if err != nil {
		return nil
	}
	return names
}

// SignedDataObjectProperties retrieves the name of each node found under the
// SignedDataObjectProperties element. Port of the public getSignedDataObjectProperties().
func (s *Signature) SignedDataObjectProperties() []string {
	names, err := xmlutils.XPathUtilsGetChildrenNames(s.SignatureElement(), s.xadesPath.SignedDataObjectPropertiesPath())
	if err != nil {
		return nil
	}
	return names
}

// DataFoundUpToLevel returns the signature level. Port of getDataFoundUpToLevel().
func (s *Signature) DataFoundUpToLevel() enumerations.SignatureLevel {
	if !s.HasBESProfile() {
		return enumerations.SignatureLevelXMLNotETSI
	}

	baselineProfile := s.HasBProfile()

	if !s.HasExtendedTProfile() {
		if baselineProfile {
			return enumerations.SignatureLevelXAdESBaselineB
		} else if s.HasEPESProfile() {
			return enumerations.SignatureLevelXAdESEPES
		}
		return enumerations.SignatureLevelXAdESBES
	}

	baselineProfile = baselineProfile && s.HasTProfile()

	if baselineProfile && s.HasLTProfile() {
		if s.HasERSProfile() {
			return enumerations.SignatureLevelXAdESERS
		}
		if s.HasLTAProfile() {
			return enumerations.SignatureLevelXAdESBaselineLTA
		}
		return enumerations.SignatureLevelXAdESBaselineLT

	} else if s.HasCProfile() {
		if s.HasXLProfile() {
			if s.HasERSProfile() {
				return enumerations.SignatureLevelXAdESERS
			}
			if s.HasAProfile() {
				return enumerations.SignatureLevelXAdESA
			}
			if s.HasXProfile() {
				return enumerations.SignatureLevelXAdESXL
			}
		}
		if s.HasXProfile() {
			return enumerations.SignatureLevelXAdESX
		}
		return enumerations.SignatureLevelXAdESC

	} else if s.HasXLProfile() {
		if s.HasERSProfile() {
			return enumerations.SignatureLevelXAdESERS
		}
		if s.HasAProfile() {
			// XAdES-E-A can be built on XAdES-E-T directly.
			return enumerations.SignatureLevelXAdESA
		}
		return enumerations.SignatureLevelXAdESLT
	}

	if baselineProfile {
		return enumerations.SignatureLevelXAdESBaselineT
	}
	return enumerations.SignatureLevelXAdEST
}

// ValidateStructure processes the structure validation of the signature. Port of validateStructure().
func (s *Signature) ValidateStructure() []string {
	structureValidator := StructureValidatorFactoryGetInstance().FromXAdESSignature(s)
	structureValidator.Validate()
	return structureValidator.ValidationErrors()
}

// FindSignatureScopes finds signature scopes. Port of the protected findSignatureScopes().
func (s *Signature) FindSignatureScopes() []scope.SignatureScope {
	return NewXAdESSignatureScopeFinder().FindSignatureScope(s)
}

// CommitmentTypeIndications obtains the information concerning commitment type indication
// linked to the signature. Port of getCommitmentTypeIndications().
func (s *Signature) CommitmentTypeIndications() []*signature.CommitmentTypeIndication {
	var result []*signature.CommitmentTypeIndication
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), s.xadesPath.CommitmentTypeIndicationPath())
	if err == nil && len(nodeList) > 0 {
		result = []*signature.CommitmentTypeIndication{}
		for _, commitmentTypeIndicationNode := range nodeList {
			identifier, _ := xmlutils.XPathUtilsGetElement(commitmentTypeIndicationNode, s.xadesPath.CurrentCommitmentIdentifierPath())
			if identifier == nil {
				continue
			}
			uri := identifier.TextContent()
			if uri == "" {
				// Upstream logs "The Identifier for a CommitmentTypeIndication is not defined!
				// The CommitmentType is skipped.".
				continue
			}

			var qualifier enumerations.ObjectIdentifierQualifier
			qualifierString := identifier.AttrValue("", definition.XAdES132AttributeQualifier.AttributeName())
			if utils.IsStringNotBlank(qualifierString) {
				qualifier = enumerations.ObjectIdentifierQualifierFromValue(qualifierString)
			}
			uri = spi.DSSUtilsObjectIdentifierValueWithQualifier(uri, qualifier)

			commitmentTypeIndication := signature.NewCommitmentTypeIndication(uri)

			descriptionNode, _ := xmlutils.XPathUtilsGetElement(commitmentTypeIndicationNode, s.xadesPath.CurrentCommitmentDescriptionPath())
			if descriptionNode != nil {
				commitmentTypeIndication.SetDescription(descriptionNode.TextContent())
			}
			docRefsNode, _ := xmlutils.XPathUtilsGetElement(commitmentTypeIndicationNode, s.xadesPath.CurrentCommitmentDocumentationReferencesPath())
			if docRefsNode != nil {
				commitmentTypeIndication.SetDocumentReferences(s.getDocumentationReferences(docRefsNode))
			}

			allSignedDataObjectsNode, _ := xmlutils.XPathUtilsGetElement(commitmentTypeIndicationNode, s.xadesPath.CurrentCommitmentAllSignedDataObjectsPath())
			if allSignedDataObjectsNode != nil {
				commitmentTypeIndication.SetAllDataSignedObjects(true)
			} else {
				commitmentObjectReferencesNodeList, cerr := xmlutils.XPathUtilsGetNodeList(commitmentTypeIndicationNode, s.xadesPath.CurrentCommitmentObjectReferencesPath())
				if cerr == nil && len(commitmentObjectReferencesNodeList) > 0 {
					commitmentTypeIndication.SetObjectReferences(s.getObjectReferences(commitmentObjectReferencesNodeList))
				}
			}

			result = append(result, commitmentTypeIndication)
		}
	}
	return result
}

// getDocumentationReferences ports the private getDocumentationReferences(Element).
func (s *Signature) getDocumentationReferences(docRefsNode *xmldom.Node) []string {
	docRefsChildNodes, err := xmlutils.XPathUtilsGetNodeList(docRefsNode, s.xadesPath.CurrentDocumentationReference())
	if err == nil && len(docRefsChildNodes) > 0 {
		docRefs := make([]string, 0, len(docRefsChildNodes))
		for _, docRefNode := range docRefsChildNodes {
			docRefs = append(docRefs, docRefNode.TextContent())
		}
		return docRefs
	}
	return nil
}

// getObjectReferences ports the private getObjectReferences(NodeList).
func (s *Signature) getObjectReferences(commitmentObjectReferencesNodeList []*xmldom.Node) []string {
	signedDataObjects := make([]string, 0, len(commitmentObjectReferencesNodeList))
	for _, node := range commitmentObjectReferencesNodeList {
		signedDataObjects = append(signedDataObjects, xmlutils.DomUtilsGetId(node.TextContent()))
	}
	return signedDataObjects
}

// References gets a list of found references. Port of the public getReferences().
//
// A reference that fails to parse is skipped, matching Java's behavior (which logs "Unable to
// retrieve reference #{} : {}" and continues).
func (s *Signature) References() []*xmldsig.Reference {
	if s.references == nil {
		xmlSignature := s.mustGetSantuarioSignature()
		signedInfo := xmlSignature.SignedInfo()
		references := make([]*xmldsig.Reference, 0, signedInfo.Length())
		for i := 0; i < signedInfo.Length(); i++ {
			reference, err := signedInfo.Item(i)
			if err != nil {
				continue
			}
			references = append(references, reference)
		}
		s.references = references
	}
	return s.references
}

// SignatureObjects gets a list of found signature ds:Object elements. Port of the public
// getSignatureObjects().
func (s *Signature) SignatureObjects() []*xmldom.Node {
	list, err := xmlutils.XPathUtilsGetNodeList(s.SignatureElement(), common.XMLDSigPathObjectPath)
	if err != nil {
		return nil
	}
	objectElements := make([]*xmldom.Node, 0, len(list))
	for _, element := range list {
		if signedProperties, serr := xmlutils.XPathUtilsGetElement(element, s.xadesPath.SignedPropertiesPath()); serr == nil && signedProperties != nil {
			// Ignore signed properties.
			continue
		}
		objectElements = append(objectElements, element)
	}
	return objectElements
}

// RegisterXAdESPaths allows registering a new XAdESPaths. Port of the deprecated
// registerXAdESPaths(XAdESPath).
func (s *Signature) RegisterXAdESPaths(xadesPaths definition.XAdESPath) {
	s.OwnerDocument().AddXAdESPathHolder(xadesPaths)
}

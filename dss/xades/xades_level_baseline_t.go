// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelBaselineT.java (DSS 6.5.RC1).
//
// # The extension chain, and how Java's virtual extendSignatures survives the port
//
// Upstream builds the whole XAdES augmentation ladder out of one overridden protected method:
//
//	XAdESLevelBaselineT.extendSignatures(List<AdvancedSignature>)
//	  <- XAdESLevelBaselineLT   <- XAdESLevelBaselineLTA
//	  <- XAdESLevelC            <- XAdESLevelX  <- XAdESLevelXL  <- XAdESLevelA
//
// each override calling super.extendSignatures(signatures) first. Go has no method overriding
// across embedding, so - per the TokenBase.InitToken(self) convention of PORTING.md, the same
// one cades/cades_signature_extension.go uses - every level embeds the level below it,
// registers itself through Init<Level>(self, certificateVerifier), and the public entry point
// ExtendSignaturesDocument dispatches into the most-derived override via t.overrides. A "super"
// call is then the plain, explicit embedded-field call, e.g. lt.XAdESLevelBaselineT.ExtendSignatures.
//
// Java's two extendSignatures overloads get two Go names:
//
//	extendSignatures(List<AdvancedSignature>)                -> ExtendSignatures (the virtual one)
//	extendSignatures(DSSDocument, XAdESSignatureParameters)  -> ExtendSignaturesDocument
//
// Likewise the two incorporateCertificateValues / incorporateRevocationValues overloads: the
// indent-taking ones carry the WithIndent suffix.
//
// # Errors
//
// Objects.requireNonNull becomes a panic carrying the Java message; every other Java throw
// becomes a returned error. slf4j logging is dropped.
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/executor"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESSignatureExtensionOverrides declares the operation every XAdES extension level overrides and
// that XAdESLevelBaselineT.ExtendSignaturesDocument dispatches into, together with the XAdESBuilder
// hooks the shared base already needs. Every level satisfies both halves through its embedded
// ancestors.
type XAdESSignatureExtensionOverrides interface {
	XAdESBuilderOverrides

	// ExtendSignatures extends the given signatures to the level of the concrete
	// implementation. Port of the protected, overridden #extendSignatures(List).
	ExtendSignatures(signatures []validation.AdvancedSignature) error
}

// XAdESLevelBaselineT is the -T profile of a XAdES signature.
type XAdESLevelBaselineT struct {
	ExtensionBuilder

	// TspSource encapsulates the Time Stamp Protocol needed to create the level -T of the
	// signature. Port of the protected tspSource.
	TspSource validation.TSPSource

	// overrides points back at the most-derived level; see InitXAdESLevelBaselineT.
	overrides XAdESSignatureExtensionOverrides
}

// NewXAdESLevelBaselineT is the default constructor for XAdESLevelBaselineT.
// Port of XAdESLevelBaselineT(CertificateVerifier).
func NewXAdESLevelBaselineT(certificateVerifier validation.CertificateVerifier) *XAdESLevelBaselineT {
	extension := &XAdESLevelBaselineT{}
	extension.InitXAdESLevelBaselineT(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelBaselineT registers the concrete extension level with this base and forwards
// the CertificateVerifier to ExtensionBuilder. Port of the protected
// XAdESLevelBaselineT(CertificateVerifier) constructor's super(certificateVerifier) call.
func (t *XAdESLevelBaselineT) InitXAdESLevelBaselineT(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	t.overrides = self
	t.InitExtensionBuilderWithVerifier(self, certificateVerifier)
}

// SetTspSource sets the TSP source to be used when extending the digital signature.
// Port of #setTspSource.
func (t *XAdESLevelBaselineT) SetTspSource(tspSource validation.TSPSource) {
	t.TspSource = tspSource
}

// incorporateC14nMethod ports the private incorporateC14nMethod: it appends
// <ds:CanonicalizationMethod Algorithm="..."/> to parentDom.
func (t *XAdESLevelBaselineT) incorporateC14nMethod(parentDom *xmldom.Node, signedInfoC14nMethod string) {
	canonicalizationMethodDom := xmlutils.DomUtilsCreateElementNS(t.DocumentDom,
		t.XmldsigNamespace(), common.XMLDSigElement_CANONICALIZATION_METHOD)
	canonicalizationMethodDom.SetAttr(
		xmldom.Name{Local: common.XMLDSigAttribute_ALGORITHM.AttributeName()}, signedInfoC14nMethod)
	parentDom.AppendChild(canonicalizationMethodDom)
}

// ExtendSignaturesDocument extends every signature of the given document to the level of the concrete
// implementation. Port of the #extendSignatures(DSSDocument, XAdESSignatureParameters) overload
// declared by SignatureExtension; Java's DSSException becomes the returned error, its
// Objects.requireNonNull a panic.
func (t *XAdESLevelBaselineT) ExtendSignaturesDocument(dssDocument model.DSSDocument,
	params *XAdESSignatureParameters) (model.DSSDocument, error) {
	if dssDocument == nil {
		panic("The document cannot be null")
	}
	if t.TspSource == nil {
		panic("The TSPSource cannot be null")
	}
	t.Params = params
	context := params.GetContext()
	// Upstream logs "====> Extending: {}" here.

	documentAnalyzer, err := NewXMLDocumentAnalyzer(dssDocument)
	if err != nil {
		return nil, err
	}
	t.DocumentAnalyzer = documentAnalyzer
	t.DocumentAnalyzer.SetCertificateVerifier(t.CertificateVerifier)
	t.DocumentAnalyzer.SetDetachedContents(params.DetachedContents())
	t.DocumentAnalyzer.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)

	t.DocumentDom = t.DocumentAnalyzer.RootElement()

	signatures := t.DocumentAnalyzer.Signatures()
	if utils.IsCollectionEmpty(signatures) {
		return nil, exception.NewIllegalInputException("No signatures found to be extended!")
	}

	// In the case of the enveloped signature we have a specific treatment:
	// we will just extend the signature that is being created (during creation process)
	signaturesToExtend := signatures

	operationKind := context.OperationKind()
	if enumerations.SigningOperation_SIGN == operationKind {
		signatureID := params.GetDeterministicId()

		for _, signature := range signatures {
			if signatureID == signature.DAIdentifier() {
				signaturesToExtend = []validation.AdvancedSignature{signature}
				break
			}
		}
	}

	signaturesToExtend = xadesLevelBaselineTAssertNoEmbeddedSignaturesPresent(signaturesToExtend)

	if err := t.overrides.ExtendSignatures(signaturesToExtend); err != nil {
		return nil, err
	}

	return t.CreateXmlDocument()
}

// xadesLevelBaselineTAssertNoEmbeddedSignaturesPresent excludes signatures embedded within other
// signatures from consecutive extension. Port of the private assertNoEmbeddedSignaturesPresent.
func xadesLevelBaselineTAssertNoEmbeddedSignaturesPresent(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	result := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's cast; an XMLDocumentAnalyzer only ever yields XAdESSignatures.
			continue
		}
		signatureElement := xadesSignature.SignatureElement()
		if !xadesLevelBaselineTHasSignatureAsParent(signatureElement) {
			result = append(result, signature)
		}
		// Upstream warns "The signature with Id '{}' has a ds:Signature parent within its XML
		// tree! The signature will not be extended." for the excluded ones.
	}
	return result
}

// xadesLevelBaselineTHasSignatureAsParent ports the private hasSignatureAsParent.
func xadesLevelBaselineTHasSignatureAsParent(element *xmldom.Node) bool {
	parent := element.Parent
	for parent != nil {
		if common.XMLDSigElement_SIGNATURE.IsSameTagName(parent.Name.Local) &&
			common.XMLDSigElement_SIGNATURE.URI() == parent.Name.Space {
			return true
		}
		parent = parent.Parent
	}
	return false
}

// ExtendSignatures extends signatures to a desired level. This method is overridden by other
// profiles. For the -T profile it adds the SignatureTimeStamp element which contains a single
// HashDataInfo element that refers to the ds:SignatureValue element of the [XMLDSIG] signature;
// the timestamp token is obtained from the TSP source. Adds a <SignatureTimeStamp> segment into
// the <UnsignedSignatureProperties> element.
// Port of the protected #extendSignatures(List).
func (t *XAdESLevelBaselineT) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	signaturesToExtend := t.extendToTLevelSignatures(signatures)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	signatureRequirementsChecker := t.SignatureRequirementsChecker()
	signatureRequirementsChecker.AssertExtendToTLevelPossible(signatures)

	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertSigningCertificatesAreValid(signaturesToExtend)

	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := t.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		levelBUnsignedProperties := t.UnsignedSignaturePropertiesDom.Clone(true)

		signatureTimestampParameters := t.Params.GetSignatureTimestampParameters()
		digestAlgorithm := signatureTimestampParameters.DigestAlgorithm()
		canonicalizationMethod := signatureTimestampParameters.CanonicalizationMethod()
		timestampSource, err := xadesLevelBaselineTTimestampSource(t.XadesSignature.TimestampSource())
		if err != nil {
			return err
		}
		messageDigest := timestampSource.GetSignatureTimestampMessageDigest(digestAlgorithm,
			canonicalizationMethod)
		if err := t.CreateXAdESTimeStampType(enumerations.TimestampType_SIGNATURE_TIMESTAMP,
			canonicalizationMethod, messageDigest); err != nil {
			return err
		}

		indented, err := t.IndentIfPrettyPrint(t.UnsignedSignaturePropertiesDom, levelBUnsignedProperties)
		if err != nil {
			return err
		}
		t.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// SignatureRequirementsChecker instantiates a SignatureRequirementsChecker.
// Port of the protected #getSignatureRequirementsChecker.
func (t *XAdESLevelBaselineT) SignatureRequirementsChecker() *document.SignatureRequirementsChecker[*XAdESTimestampParameters] {
	return document.NewSignatureRequirementsChecker[*XAdESTimestampParameters](t.CertificateVerifier,
		&t.Params.AbstractSignatureParameters)
}

// extendToTLevelSignatures ports the private getExtendToTLevelSignatures.
func (t *XAdESLevelBaselineT) extendToTLevelSignatures(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if t.tLevelExtensionRequired(signature) {
			toBeExtended = append(toBeExtended, signature)
		}
	}
	return toBeExtended
}

// tLevelExtensionRequired ports the private tLevelExtensionRequired.
func (t *XAdESLevelBaselineT) tLevelExtensionRequired(signature validation.AdvancedSignature) bool {
	return enumerations.SignatureLevel_XAdES_BASELINE_T == t.Params.SignatureLevel() || !signature.HasTProfile()
}

// IncorporateCertificateValues incorporates all certificates passed as parameter:
//
//	<xades:CertificateValues>
//	    <xades:EncapsulatedX509Certificate>MIIC9TC...</xades:EncapsulatedX509Certificate>
//	    ...
//	</xades:CertificateValues>
//
// Port of the protected #incorporateCertificateValues(Element, Collection).
func (t *XAdESLevelBaselineT) IncorporateCertificateValues(parentDom *xmldom.Node,
	certificatesToBeAdded []*model.CertificateToken) (*xmldom.Node, error) {
	var certificateValuesDom *xmldom.Node
	if utils.IsCollectionNotEmpty(certificatesToBeAdded) {
		currentElements, err := t.CurrentXAdESElements()
		if err != nil {
			return nil, err
		}
		certificateValuesDom = xmlutils.DomUtilsAddElement(t.DocumentDom, parentDom,
			t.XadesNamespace(), currentElements.ElementCertificateValues())
		for _, certificateToken := range certificatesToBeAdded {
			base64EncodeCertificate := utils.ToBase64(certificateToken.Encoded())
			xmlutils.DomUtilsAddTextElement(t.DocumentDom, certificateValuesDom,
				t.XadesNamespace(), currentElements.ElementEncapsulatedX509Certificate(),
				base64EncodeCertificate)
		}
	}
	return certificateValuesDom, nil
}

// IncorporateRevocationValues incorporates revocation values.
//
//	<xades:RevocationValues>
//
// Port of the protected #incorporateRevocationValues(Element, Collection, Collection).
func (t *XAdESLevelBaselineT) IncorporateRevocationValues(parentDom *xmldom.Node,
	crlsToAdd []*spi.CRLToken, ocspsToAdd []*spi.OCSPToken) (*xmldom.Node, error) {
	var revocationValuesDom *xmldom.Node

	if utils.IsCollectionNotEmpty(crlsToAdd) || utils.IsCollectionNotEmpty(ocspsToAdd) {
		currentElements, err := t.CurrentXAdESElements()
		if err != nil {
			return nil, err
		}
		revocationValuesDom = xmlutils.DomUtilsAddElement(t.DocumentDom, parentDom,
			t.XadesNamespace(), currentElements.ElementRevocationValues())
		if err := t.incorporateCrlTokens(revocationValuesDom, crlsToAdd); err != nil {
			return nil, err
		}
		if err := t.incorporateOcspTokens(revocationValuesDom, ocspsToAdd); err != nil {
			return nil, err
		}
	}
	return revocationValuesDom, nil
}

// incorporateCrlTokens incorporates the CRLValues:
//
//	<xades:CRLValues>
//	    <xades:EncapsulatedCRLValue>...</xades:EncapsulatedCRLValue>
//	    ...
//	</xades:CRLValues>
//
// Port of the private incorporateCrlTokens.
func (t *XAdESLevelBaselineT) incorporateCrlTokens(parentDom *xmldom.Node, crlTokens []*spi.CRLToken) error {
	if len(crlTokens) == 0 {
		return nil
	}
	currentElements, err := t.CurrentXAdESElements()
	if err != nil {
		return err
	}
	crlValuesDom := xmlutils.DomUtilsAddElement(t.DocumentDom, parentDom, t.XadesNamespace(),
		currentElements.ElementCRLValues())

	for _, revocationToken := range crlTokens {
		encodedCRL := revocationToken.Encoded()
		base64EncodedCRL := utils.ToBase64(encodedCRL)
		xmlutils.DomUtilsAddTextElement(t.DocumentDom, crlValuesDom, t.XadesNamespace(),
			currentElements.ElementEncapsulatedCRLValue(), base64EncodedCRL)
	}
	return nil
}

// incorporateOcspTokens incorporates the OCSP responses:
//
//	<xades:OCSPValues>
//	    <xades:EncapsulatedOCSPValue>...</xades:EncapsulatedOCSPValue>
//	    ...
//	</xades:OCSPValues>
//
// Port of the private incorporateOcspTokens.
func (t *XAdESLevelBaselineT) incorporateOcspTokens(parentDom *xmldom.Node, ocspTokens []*spi.OCSPToken) error {
	if len(ocspTokens) == 0 {
		return nil
	}
	currentElements, err := t.CurrentXAdESElements()
	if err != nil {
		return err
	}
	ocspValuesDom := xmlutils.DomUtilsAddElement(t.DocumentDom, parentDom, t.XadesNamespace(),
		currentElements.ElementOCSPValues())

	for _, revocationToken := range ocspTokens {
		encodedOCSP := revocationToken.Encoded()
		base64EncodedOCSP := utils.ToBase64(encodedOCSP)
		xmlutils.DomUtilsAddTextElement(t.DocumentDom, ocspValuesDom, t.XadesNamespace(),
			currentElements.ElementEncapsulatedOCSPValue(), base64EncodedOCSP)
	}
	return nil
}

// IncorporateCertificateValuesWithIndent incorporates all certificates passed as parameter, as
// well as adding missing indents if the indent parameter is specified.
// Port of the protected #incorporateCertificateValues(Element, Collection, String).
func (t *XAdESLevelBaselineT) IncorporateCertificateValuesWithIndent(parentDom *xmldom.Node,
	certificatesToBeAdded []*model.CertificateToken, indent string) error {
	certificatesDom, err := t.IncorporateCertificateValues(parentDom, certificatesToBeAdded)
	if err != nil {
		return err
	}
	if certificatesDom != nil && indent != "" {
		xmlutils.DomUtilsSetTextNode(t.DocumentDom, parentDom, indent)
		if _, err := DSSXMLUtilsIndentAndReplace(t.DocumentDom, certificatesDom); err != nil {
			return err
		}
	}
	return nil
}

// IncorporateRevocationValuesWithIndent incorporates revocation values, as well as adding
// missing indents if the indent parameter is specified.
// Port of the protected #incorporateRevocationValues(Element, Collection, Collection, String).
func (t *XAdESLevelBaselineT) IncorporateRevocationValuesWithIndent(parentDom *xmldom.Node,
	crlsToAdd []*spi.CRLToken, ocspsToAdd []*spi.OCSPToken, indent string) error {
	revocationDom, err := t.IncorporateRevocationValues(parentDom, crlsToAdd, ocspsToAdd)
	if err != nil {
		return err
	}
	if revocationDom != nil && indent != "" {
		xmlutils.DomUtilsSetTextNode(t.DocumentDom, parentDom, indent)
		if _, err := DSSXMLUtilsIndentAndReplace(t.DocumentDom, revocationDom); err != nil {
			return err
		}
	}
	return nil
}

// RemoveOldCertificateValues removes old certificate values from the unsigned signature
// properties element and returns the indent of the removed node, nil when there was none.
// Port of the protected #removeOldCertificateValues. Java's String return is null-able and its
// callers test it against null (see XAdESLevelBaselineLT); ExtensionBuilder.RemoveNode answers
// "" where Java answers null, and the two agree in practice - a removed node either has a
// whitespace TEXT sibling in a pretty-printed document or none at all, never an empty one.
func (t *XAdESLevelBaselineT) RemoveOldCertificateValues() (string, error) {
	text := ""
	toRemove, err := xmlutils.XPathUtilsGetElement(t.XadesSignature.SignatureElement(),
		t.XadesPath.CertificateValuesPath())
	if err != nil {
		return "", err
	}
	if toRemove != nil {
		text = t.RemoveNode(toRemove)
		t.XadesSignature.ResetCertificateSource()
	}
	return text, nil
}

// RemoveOldRevocationValues removes old revocation values from the unsigned signature
// properties element. Port of the protected #removeOldRevocationValues.
func (t *XAdESLevelBaselineT) RemoveOldRevocationValues() error {
	toRemove, err := xmlutils.XPathUtilsGetElement(t.XadesSignature.SignatureElement(),
		t.XadesPath.RevocationValuesPath())
	if err != nil {
		return err
	}
	if toRemove != nil {
		t.RemoveNode(toRemove)
		t.XadesSignature.ResetRevocationSources()
	}
	return nil
}

// RemoveLastTimestampAndAnyValidationData removes the TimeStampValidationData and
// AnyValidationData elements appearing at the end of the unsigned properties, and returns the
// indent of the last removed element, nil when none was present.
// Port of the protected #removeLastTimestampAndAnyValidationData.
func (t *XAdESLevelBaselineT) RemoveLastTimestampAndAnyValidationData() (string, error) {
	toRemove, err := t.lastElementIfPresent(definition.XAdES141Element_TIMESTAMP_VALIDATION_DATA,
		definition.XAdES141Element_ANY_VALIDATION_DATA)
	if err != nil {
		return "", err
	}
	if toRemove == nil {
		return "", nil
	}
	intent := ""
	for toRemove != nil {
		intent = t.RemoveNode(toRemove)
		toRemove, err = t.lastElementIfPresent(definition.XAdES141Element_TIMESTAMP_VALIDATION_DATA,
			definition.XAdES141Element_ANY_VALIDATION_DATA)
		if err != nil {
			return "", err
		}
	}
	/*
	 * Certificate and revocation sources need to be reset because of
	 * the removing of TimeStampValidationData or AnyValidationData element
	 */
	t.XadesSignature.ResetCertificateSource()
	t.XadesSignature.ResetRevocationSources()
	return intent, nil
}

// lastElementIfPresent ports the private getLastElementIfPresent.
func (t *XAdESLevelBaselineT) lastElementIfPresent(xadesElements ...common.DSSElement) (*xmldom.Node, error) {
	unsignedSignatureProperties, err := xmlutils.XPathUtilsGetNode(t.XadesSignature.SignatureElement(),
		t.XadesPath.UnsignedSignaturePropertiesPath())
	if err != nil {
		return nil, err
	}
	if unsignedSignatureProperties != nil {
		childNodes, err := xmlutils.XPathUtilsGetNodeList(unsignedSignatureProperties,
			common.XPathQueryBuilderFromCurrentPosition().Build())
		if err != nil {
			return nil, err
		}
		if len(childNodes) > 0 {
			unsignedSignatureElement := childNodes[len(childNodes)-1]
			nodeName := unsignedSignatureElement.Name.Local
			for _, element := range xadesElements {
				if element.IsSameTagName(nodeName) {
					return unsignedSignatureElement, nil
				}
			}
		}
	}
	return nil, nil
}

// IncorporateTimestampValidationData incorporates the timestamp validation data in the
// signature. Port of the protected #incorporateTimestampValidationData.
func (t *XAdESLevelBaselineT) IncorporateTimestampValidationData(
	validationDataForInclusion *validation.ValidationData, indent string) error {
	return t.IncorporateValidationData(validationDataForInclusion, indent,
		definition.XAdES141Element_TIMESTAMP_VALIDATION_DATA, xadesBuilderTstVdPrefix)
}

// IncorporateAnyValidationData incorporates the AnyValidationData in the signature.
// Port of the protected #incorporateAnyValidationData.
func (t *XAdESLevelBaselineT) IncorporateAnyValidationData(
	validationDataForInclusion *validation.ValidationData, indent string) error {
	return t.IncorporateValidationData(validationDataForInclusion, indent,
		definition.XAdES141Element_ANY_VALIDATION_DATA, xadesBuilderAnyVdPrefix)
}

// IncorporateValidationData incorporates the given validation data in the signature, under the
// given XAdES 1.4.1 element and with the given Id prefix.
// Port of the protected #incorporateValidationData.
func (t *XAdESLevelBaselineT) IncorporateValidationData(
	validationDataForInclusion *validation.ValidationData, indent string,
	element common.DSSElement, prefix string) error {
	if validationDataForInclusion.IsEmpty() {
		return nil
	}

	certificateValuesToAdd := validationDataForInclusion.CertificateTokens()
	crlsToAdd := validationDataForInclusion.CrlTokens()
	ocspsToAdd := validationDataForInclusion.OcspTokens()

	timeStampValidationDataDom := xmlutils.DomUtilsAddElement(t.DocumentDom,
		t.UnsignedSignaturePropertiesDom, t.Xades141Namespace(), element)

	if err := t.IncorporateCertificateValuesWithIndent(timeStampValidationDataDom,
		certificateValuesToAdd, indent); err != nil {
		return err
	}
	if err := t.IncorporateRevocationValuesWithIndent(timeStampValidationDataDom,
		crlsToAdd, ocspsToAdd, indent); err != nil {
		return err
	}

	id := "1"
	timestamps := t.XadesSignature.AllTimestamps()
	if utils.IsCollectionNotEmpty(timestamps) {
		timestampToken := timestamps[len(timestamps)-1]
		xmlID, err := t.ToXmlIdentifier(timestampToken.DSSID())
		if err != nil {
			return err
		}
		id = xmlID
	}

	timeStampValidationDataDom.SetAttr(xmldom.Name{Local: "Id"}, prefix+id)
	if t.Params.IsPrettyPrint() {
		if _, err := DSSXMLUtilsIndentAndReplace(t.DocumentDom, timeStampValidationDataDom); err != nil {
			return err
		}
	}
	return nil
}

// IncorporateArchiveTimestamp incorporates the archive timestamp object.
// Port of the protected #incorporateArchiveTimestamp.
func (t *XAdESLevelBaselineT) IncorporateArchiveTimestamp() error {
	archiveTimestampParameters := t.Params.GetArchiveTimestampParameters()
	digestAlgorithm := archiveTimestampParameters.DigestAlgorithm()
	canonicalizationMethod := archiveTimestampParameters.CanonicalizationMethod()
	timestampSource, err := xadesLevelBaselineTTimestampSource(t.XadesSignature.TimestampSource())
	if err != nil {
		return err
	}
	messageDigest := timestampSource.GetArchiveTimestampData(digestAlgorithm, canonicalizationMethod)
	return t.CreateXAdESTimeStampType(enumerations.TimestampType_ARCHIVE_TIMESTAMP,
		canonicalizationMethod, messageDigest)
}

// CreateXAdESTimeStampType creates any XAdES TimeStamp object representation. The timestamp
// token is obtained from the TSP source. Port of the protected #createXAdESTimeStampType.
func (t *XAdESLevelBaselineT) CreateXAdESTimeStampType(timestampType enumerations.TimestampType,
	timestampC14nMethod string, messageDigest model.DSSMessageDigest) error {

	xadesURI := t.XadesNamespace().Uri()
	if (definition.XAdESNamespace_XADES_111.IsSameUri(xadesURI) ||
		definition.XAdESNamespace_XADES_122.IsSameUri(xadesURI)) &&
		enumerations.TimestampType_SIGNATURE_TIMESTAMP != timestampType {
		return fmt.Errorf("Signature Timestamp creation is only supported for XAdES 1.1.1 and 1.2.2")
	}

	currentElements, err := t.CurrentXAdESElements()
	if err != nil {
		return err
	}

	signatureTimestampParameters := t.Params.GetSignatureTimestampParameters()
	timestampDigestAlgorithm := signatureTimestampParameters.DigestAlgorithm()

	var timeStampDom *xmldom.Node
	switch timestampType {
	case enumerations.TimestampType_SIGNATURE_TIMESTAMP:
		// <xades:SignatureTimeStamp Id="time-stamp-1dee38c4-8388-40d1-8880-9eeda853fe60">
		timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
			t.XadesNamespace(), currentElements.ElementSignatureTimeStamp())
	case enumerations.TimestampType_VALIDATION_DATA_TIMESTAMP:
		// <xades:SigAndRefsTimeStamp Id="time-stamp-a762ab0e-e05c-4cc8-a804-cf2c4ffb5516">
		if t.Params.IsEn319132() {
			timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
				t.Xades141Namespace(), definition.XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2)
		} else {
			timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
				t.XadesNamespace(), currentElements.ElementSigAndRefsTimeStamp())
		}
	case enumerations.TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP:
		// <xades:RefsOnlyTimeStamp Id="time-stamp-a762ab0e-e05c-4cc8-a804-cf2c4ffb5516">
		if t.Params.IsEn319132() {
			timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
				t.Xades141Namespace(), definition.XAdES141Element_REFS_ONLY_TIMESTAMP_V2)
		} else {
			timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
				t.XadesNamespace(), currentElements.ElementRefsOnlyTimeStamp())
		}
	case enumerations.TimestampType_ARCHIVE_TIMESTAMP:
		// <xades141:ArchiveTimeStamp Id="time-stamp-a762ab0e-e05c-4cc8-a804-cf2c4ffb5516">
		timeStampDom = xmlutils.DomUtilsAddElement(t.DocumentDom, t.UnsignedSignaturePropertiesDom,
			t.Xades141Namespace(), definition.XAdES141Element_ARCHIVE_TIMESTAMP)
		timestampDigestAlgorithm = t.Params.GetArchiveTimestampParameters().DigestAlgorithm()
	default:
		// Content timestamps need to be generated before the signature itself
		return fmt.Errorf("Unsupported timestamp type : %s", timestampType)
	}

	// Upstream logs "Timestamp generation: {} / {} / {}" here.
	timeStampToken, err := t.TspSource.TimeStampResponse(timestampDigestAlgorithm, messageDigest.Value())
	if err != nil {
		return err
	}
	derEncoded, err := spi.DSSASN1UtilsDEREncodedTimestampBinary(timeStampToken)
	if err != nil {
		return err
	}
	base64EncodedTimeStampToken := utils.ToBase64(derEncoded)

	if definition.XAdESNamespace_XADES_122.IsSameUri(xadesURI) {
		t.incorporateXAdES122Include(timeStampDom)
	}
	if definition.XAdESNamespace_XADES_111.IsSameUri(xadesURI) {
		t.incorporateHashDataInfo(timeStampDom, timestampC14nMethod)
	} else {
		// <ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
		t.incorporateC14nMethod(timeStampDom, timestampC14nMethod)
	}

	// <xades:EncapsulatedTimeStamp Id="time-stamp-token-6a150419-caab-4615-9a0b-6e239596643a">MIAGCSqGSIb3DQEH
	encapsulatedTimeStampDom := xmlutils.DomUtilsAddElement(t.DocumentDom, timeStampDom,
		t.XadesNamespace(), currentElements.ElementEncapsulatedTimeStamp())
	xmlutils.DomUtilsSetTextNode(t.DocumentDom, encapsulatedTimeStampDom, base64EncodedTimeStampToken)

	// Id="..." attribute is not allowed in XAdES 1.1.1
	if !definition.XAdESNamespace_XADES_111.IsSameUri(xadesURI) {
		// Add Id after the element is constructed
		attributeIdentifier := XAdESAttributeIdentifierBuild(timeStampDom)
		timestampID, err := t.ToXmlIdentifier(attributeIdentifier)
		if err != nil {
			return err
		}
		timeStampDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ID.AttributeName()},
			xadesBuilderTimestampPrefix+timestampID)
		encapsulatedTimeStampDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ID.AttributeName()},
			xadesBuilderEncapsulatedTimestampPrefix+timestampID)
	}
	return nil
}

// incorporateHashDataInfo builds
//
//	<HashDataInfo URI="AI-NDS-HGI-32019423">
//	    <Transforms xmlns="http://www.w3.org/2000/09/xmldsig#">
//	        <Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform>
//	    </Transforms>
//	</HashDataInfo>
//
// Port of the private incorporateHashDataInfo.
func (t *XAdESLevelBaselineT) incorporateHashDataInfo(timeStampDom *xmldom.Node, timestampC14nMethod string) {
	hashDataInfoDom := xmlutils.DomUtilsAddElement(t.DocumentDom, timeStampDom,
		t.XadesNamespace(), definition.XAdES111Element_HASH_DATA_INFO)
	hashDataInfoDom.SetAttr(xmldom.Name{Local: definition.XAdES111Attribute_URI.AttributeName()},
		"#"+t.XadesSignature.ID())
	transformsDom := xmlutils.DomUtilsAddElement(t.DocumentDom, hashDataInfoDom,
		t.XadesNamespace(), definition.XAdES111Element_TRANSFORMS)
	transformDom := xmlutils.DomUtilsAddElement(t.DocumentDom, transformsDom,
		t.XmldsigNamespace(), common.XMLDSigElement_TRANSFORM)
	transformDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ALGORITHM.AttributeName()},
		timestampC14nMethod)
}

// incorporateXAdES122Include ports the private incorporateXAdES122Include.
func (t *XAdESLevelBaselineT) incorporateXAdES122Include(timeStampDom *xmldom.Node) {
	includeDom := xmlutils.DomUtilsAddElement(t.DocumentDom, timeStampDom, t.XadesNamespace(),
		definition.XAdES122Element_INCLUDE)
	includeDom.SetAttr(xmldom.Name{Local: definition.XAdES122Attribute_URI.AttributeName()},
		"#"+t.XadesSignature.SignatureValueId())
}

// AssertDetachedDocumentsContainBinaries checks that the detached content is represented by
// binary documents (used for the -LTA level extension).
// Port of the protected #assertDetachedDocumentsContainBinaries; Java's
// IllegalArgumentException becomes a returned error.
func (t *XAdESLevelBaselineT) AssertDetachedDocumentsContainBinaries() error {
	detachedContents := t.Params.DetachedContents()
	if utils.IsCollectionNotEmpty(detachedContents) {
		for _, detachedDocument := range detachedContents {
			if _, isDigestDocument := detachedDocument.(*model.DigestDocument); isDigestDocument {
				return fmt.Errorf("XAdES-LTA requires complete binaries of signed documents! " +
					"Extension with a DigestDocument is not possible.")
			}
		}
	}
	return nil
}

// xadesLevelBaselineTTimestampSource narrows the signature's timestamp source to the XAdES one.
// Java's XAdESSignature#getTimestampSource() is a covariant override returning
// XAdESTimestampSource; whether the Go port of XAdESSignature returns that concrete type or the
// validation.TimestampSource interface the AdvancedSignature contract declares (the CAdES
// precedent), an `any` parameter accepts both, and the assertion below recovers the concrete
// source either way.
func xadesLevelBaselineTTimestampSource(source any) (*XAdESTimestampSource, error) {
	if timestampSource, ok := source.(*XAdESTimestampSource); ok {
		return timestampSource, nil
	}
	return nil, fmt.Errorf("the signature timestamp source is not a XAdESTimestampSource, but %T", source)
}

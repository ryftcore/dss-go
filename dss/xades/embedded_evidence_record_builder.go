// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/evidencerecord/EmbeddedEvidenceRecordBuilder.java (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// EmbeddedEvidenceRecordBuilder is used to embed an existing evidence record to a XAdES
// signature.
type EmbeddedEvidenceRecordBuilder struct {
	ExtensionBuilder
}

// NewEmbeddedEvidenceRecordBuilder is the default constructor. Port of the
// EmbeddedEvidenceRecordBuilder(CertificateVerifier) constructor.
func NewEmbeddedEvidenceRecordBuilder(certificateVerifier validation.CertificateVerifier) *EmbeddedEvidenceRecordBuilder {
	b := &EmbeddedEvidenceRecordBuilder{}
	b.InitExtensionBuilderWithVerifier(b, validation.NewCertificateVerifierBuilder(certificateVerifier).BuildOfflineCopy())
	return b
}

// AddEvidenceRecord adds the evidence record document to a signature with the given
// signatureId (via parameters), provided the evidence record correctly applies to the
// signature. Port of #addEvidenceRecord(DSSDocument, DSSDocument,
// XAdESEvidenceRecordIncorporationParameters).
//
// Panics when signatureDocument, evidenceRecordDocument or parameters is nil (Java
// Objects.requireNonNull).
func (b *EmbeddedEvidenceRecordBuilder) AddEvidenceRecord(signatureDocument, evidenceRecordDocument model.DSSDocument,
	parameters *XAdESEvidenceRecordIncorporationParameters) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}
	if evidenceRecordDocument == nil {
		panic("Evidence record document must be provided!")
	}
	if parameters == nil {
		panic("XAdESEvidenceRecordIncorporationParameters must be provided!")
	}

	signature, err := b.getXAdESSignature(signatureDocument, parameters.SignatureId(), parameters.DetachedContents())
	if err != nil {
		return nil, err
	}
	return b.addEvidenceRecord(signature, evidenceRecordDocument, parameters)
}

// getXAdESSignature gets a signature to incorporate evidence record into. Port of the
// protected #getXAdESSignature(DSSDocument, String, List).
func (b *EmbeddedEvidenceRecordBuilder) getXAdESSignature(signatureDocument model.DSSDocument, signatureId string,
	detachedContent []model.DSSDocument) (*XAdESSignature, error) {
	documentAnalyzer, err := b.initDocumentAnalyzer(signatureDocument, detachedContent)
	if err != nil {
		return nil, err
	}

	if signatureId != "" {
		sig := documentAnalyzer.SignatureByID(signatureId)
		if sig == nil {
			return nil, fmt.Errorf("unable to find a signature with Id : %s!", signatureId)
		}
		xadesSignature, ok := sig.(*XAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", sig)
		}
		return xadesSignature, nil
	}

	signatures := documentAnalyzer.Signatures()
	switch {
	case utils.IsCollectionEmpty(signatures):
		return nil, exception.NewIllegalInputException(fmt.Sprintf("No signatures found in the document with name '%s'",
			documentAnalyzer.Document().Name()))
	case utils.CollectionSize(signatures) > 1:
		return nil, fmt.Errorf("more than one signature found in a document with name '%s'! "+
			"Please provide a signatureId within the parameters.", documentAnalyzer.Document().Name())
	}
	// if one signature
	xadesSignature, ok := signatures[0].(*XAdESSignature)
	if !ok {
		return nil, fmt.Errorf("unexpected signature type %T", signatures[0])
	}
	return xadesSignature, nil
}

// addEvidenceRecord adds evidenceRecordDocument to the signature's DOM. Port of the protected
// #addEvidenceRecord(XAdESSignature, DSSDocument, XAdESEvidenceRecordIncorporationParameters).
func (b *EmbeddedEvidenceRecordBuilder) addEvidenceRecord(xadesSignature *XAdESSignature, evidenceRecordDocument model.DSSDocument,
	parameters *XAdESEvidenceRecordIncorporationParameters) (model.DSSDocument, error) {
	xadesSignature, err := b.InitializeSignatureBuilder(xadesSignature)
	if err != nil {
		return nil, err
	}

	if err := b.EnsureUnsignedProperties(); err != nil {
		return nil, err
	}
	if err := b.EnsureUnsignedSignatureProperties(); err != nil {
		return nil, err
	}

	unsignedAttribute := b.getUnsignedAttributeToEmbed(parameters)
	evidenceRecord, err := b.getEvidenceRecord(evidenceRecordDocument, xadesSignature, unsignedAttribute, parameters)
	if err != nil {
		return nil, err
	}

	sealingEvidenceRecordElement, err := b.getSealingEvidenceRecordElement(unsignedAttribute, parameters)
	if err != nil {
		return nil, err
	}

	var evidenceRecordElement *xmldom.Node
	switch evidenceRecord.EvidenceRecordType() {
	case enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord:
		erDom, err := xmlutils.DomUtilsBuildDOMFromDocument(evidenceRecordDocument)
		if err != nil {
			return nil, err
		}
		nodes := xmlutils.DomUtilsAdoptChildren(sealingEvidenceRecordElement, erDom)
		evidenceRecordElement, err = b.getEvidenceRecordElement(nodes)
		if err != nil {
			return nil, err
		}

	case enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord:
		base64EncodedER := utils.ToBase64(evidenceRecord.Encoded())
		evidenceRecordElement = xmlutils.DomUtilsAddTextElement(b.DocumentDom, sealingEvidenceRecordElement,
			parameters.XadesERNamespace(), definition.XAdESEvidencerecordNamespaceElementASN1EvidenceRecord, base64EncodedER)

	default:
		return nil, fmt.Errorf("the Evidence Record type '%s' is not supported!", evidenceRecord.EvidenceRecordType())
	}

	/*
	 * In case of XAdES embedded ER, we need to first embed the ER within signature before its
	 * validation, to ensure correct namespace processing on canonicalization. The signature and
	 * all the related data shall be re-initialized for a proper ER validation.
	 */
	if enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord == evidenceRecord.EvidenceRecordType() {
		signatureDocument := xmlutils.NewDOMDocument(evidenceRecordElement.OwnerDocument())
		xadesSignature, err = b.getXAdESSignature(signatureDocument, parameters.SignatureId(), parameters.DetachedContents())
		if err != nil {
			return nil, err
		}
		xadesSignature, err = b.InitializeSignatureBuilder(xadesSignature)
		if err != nil {
			return nil, err
		}
		unsignedAttribute = b.getLastSealingEvidenceRecordAttribute()

		evidenceRecordDocument = xmlutils.NewDOMDocument(evidenceRecordElement)
		evidenceRecord, err = b.getEvidenceRecord(evidenceRecordDocument, xadesSignature, unsignedAttribute, parameters)
		if err != nil {
			return nil, err
		}
	}

	if err := b.assertEvidenceRecordValid(evidenceRecord, unsignedAttribute, parameters); err != nil {
		return nil, err
	}

	return b.CreateXmlDocument()
}

// getEvidenceRecordElement ports the private getEvidenceRecordElement(Collection<Node>).
//
// Java throws IllegalStateException("No EvidenceRecord element found!") on the missing element
// case; per PORTING.md's throw->(T,error) convention this is surfaced as a returned error
// instead.
func (b *EmbeddedEvidenceRecordBuilder) getEvidenceRecordElement(nodes []*xmldom.Node) (*xmldom.Node, error) {
	for _, node := range nodes {
		if xmldom.Element == node.Kind && definition.XAdESEvidencerecordNamespaceElementEvidenceRecord.IsSameTagName(node.Name.Local) {
			return node, nil
		}
	}
	return nil, fmt.Errorf("no EvidenceRecord element found!")
}

// getUnsignedAttributeToEmbed ports the private
// getUnsignedAttributeToEmbed(XAdESEvidenceRecordIncorporationParameters).
func (b *EmbeddedEvidenceRecordBuilder) getUnsignedAttributeToEmbed(parameters *XAdESEvidenceRecordIncorporationParameters) *XAdESAttribute {
	if parameters.IsParallelEvidenceRecord() {
		return b.getLastSealingEvidenceRecordAttribute()
	}
	// new XAdESAttribute to be created
	return nil
}

// getLastSealingEvidenceRecordAttribute ports the private
// getLastSealingEvidenceRecordAttribute().
func (b *EmbeddedEvidenceRecordBuilder) getLastSealingEvidenceRecordAttribute() *XAdESAttribute {
	unsignedSigProperties := NewXAdESUnsignedSigProperties(b.UnsignedSignaturePropertiesDom, b.XadesPath)
	return XAdESSignatureUtilsGetLastSealingEvidenceRecordAttribute(unsignedSigProperties)
}

// getEvidenceRecord ports the private getEvidenceRecord(DSSDocument, XAdESSignature,
// XAdESAttribute, XAdESEvidenceRecordIncorporationParameters).
func (b *EmbeddedEvidenceRecordBuilder) getEvidenceRecord(evidenceRecordDocument model.DSSDocument, signature *XAdESSignature,
	unsignedAttribute *XAdESAttribute, parameters *XAdESEvidenceRecordIncorporationParameters) (validation.EvidenceRecord, error) {
	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(evidenceRecordDocument)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("Unable to build an evidence record from the provided document. Reason : %s", err.Error()), err)
	}

	embeddedEvidenceRecordHelper := NewXAdESEmbeddedEvidenceRecordHelper(signature, unsignedAttribute)
	embeddedEvidenceRecordHelper.SetDetachedContents(parameters.DetachedContents())
	evidenceRecordAnalyzer.SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper)

	return evidenceRecordAnalyzer.EvidenceRecord(), nil
}

// assertEvidenceRecordValid ports the private #assertEvidenceRecordValid(EvidenceRecord,
// XAdESAttribute, XAdESEvidenceRecordIncorporationParameters).
func (b *EmbeddedEvidenceRecordBuilder) assertEvidenceRecordValid(evidenceRecord validation.EvidenceRecord, unsignedAttribute *XAdESAttribute,
	parameters *XAdESEvidenceRecordIncorporationParameters) error {
	if unsignedAttribute != nil {
		if err := b.assertContainsOnlySameTypeEvidenceRecords(unsignedAttribute, evidenceRecord.EvidenceRecordType()); err != nil {
			return err
		}
	}
	for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
		if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference != referenceValidation.Type() && !referenceValidation.IsIntact() {
			if utils.IsCollectionEmpty(parameters.DetachedContents()) {
				return exception.NewIllegalInputException("The digest covered by the evidence record do not correspond to " +
					"the digest computed on the signature and/or detached content! " +
					"In case of detached signature, please use #setDetachedContent method to provide original documents.")
			}
			return exception.NewIllegalInputException("The digest covered by the evidence record do not correspond to " +
				"the digest computed on the signature and/or detached content!")
		}
	}
	b.validateTimestamps(evidenceRecord)
	return nil
}

// assertContainsOnlySameTypeEvidenceRecords ports the private
// assertContainsOnlySameTypeEvidenceRecords(XAdESAttribute, EvidenceRecordTypeEnum).
func (b *EmbeddedEvidenceRecordBuilder) assertContainsOnlySameTypeEvidenceRecords(unsignedAttribute *XAdESAttribute,
	evidenceRecordType enumerations.EvidenceRecordTypeEnum) error {
	sealingEvidenceRecordElement := unsignedAttribute.Element()

	for childNode := sealingEvidenceRecordElement.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if xmldom.Element != childNode.Kind {
			continue
		}
		switch evidenceRecordType {
		case enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord:
			if !definition.XAdESEvidencerecordNamespaceElementEvidenceRecord.IsSameTagName(childNode.Name.Local) {
				return exception.NewIllegalInputException(
					"The latest signature unsigned property contains evidence records other " +
						"than ers:EvidenceRecordType type specified in IETF RFC 6283. " +
						"The incorporation of different evidence record types within " +
						"the same unsigned property is not supported.")
			}
		case enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord:
			if !definition.XAdESEvidencerecordNamespaceElementASN1EvidenceRecord.IsSameTagName(childNode.Name.Local) {
				return exception.NewIllegalInputException(
					"The latest signature unsigned property contains evidence records other " +
						"than EvidenceRecord type specified in IETF RFC 4998. " +
						"The incorporation of different evidence record types within " +
						"the same unsigned property is not supported.")
			}
		default:
			return fmt.Errorf("the evidence record type '%s' is not supported!", evidenceRecordType)
		}
	}
	return nil
}

// validateTimestamps ports the private #validateTimestamps(EvidenceRecord).
func (b *EmbeddedEvidenceRecordBuilder) validateTimestamps(evidenceRecord validation.EvidenceRecord) {
	validationContext := validation.NewSignatureValidationContext()
	validationContext.Initialize(b.CertificateVerifier)

	validationContext.AddDocumentCertificateSource(evidenceRecord.CertificateSource())
	for _, timestampToken := range evidenceRecord.Timestamps() {
		validationContext.AddTimestampTokenForVerification(timestampToken)
	}

	validationContext.Validate()

	signatureValidationAlerter := validation.NewSignatureValidationAlerter(validationContext)
	signatureValidationAlerter.SetSigningOperation(enumerations.SigningOperationAddEvidenceRecord)
	signatureValidationAlerter.AssertAllTimestampsValid()
}

// getSealingEvidenceRecordElement ports the private
// getSealingEvidenceRecordElement(XAdESAttribute, XAdESEvidenceRecordIncorporationParameters).
func (b *EmbeddedEvidenceRecordBuilder) getSealingEvidenceRecordElement(unsignedAttribute *XAdESAttribute,
	parameters *XAdESEvidenceRecordIncorporationParameters) (*xmldom.Node, error) {
	if unsignedAttribute != nil {
		// parallel evidence record
		return unsignedAttribute.Element(), nil
	}
	// new evidence record unsigned property
	return xmlutils.DomUtilsAddElement(b.DocumentDom, b.UnsignedSignaturePropertiesDom,
		parameters.XadesERNamespace(), definition.XAdESEvidencerecordNamespaceElementSealingEvidenceRecords), nil
}

// initDocumentAnalyzer ports the private #initDocumentAnalyzer(DSSDocument, List<DSSDocument>).
func (b *EmbeddedEvidenceRecordBuilder) initDocumentAnalyzer(signatureDocument model.DSSDocument, detachedContents []model.DSSDocument) (*XMLDocumentAnalyzer, error) {
	documentAnalyzer, err := b.InitDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}
	documentAnalyzer.SetDetachedContents(detachedContents)
	return documentAnalyzer, nil
}

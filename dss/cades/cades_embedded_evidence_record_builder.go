// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/evidencerecord/CAdESEmbeddedEvidenceRecordBuilder.java (DSS 6.5.RC1).
package cades

import (
	"encoding/asn1"
	"fmt"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
)

// CAdESEmbeddedEvidenceRecordBuilder is used to embed an existing ERS evidence record within a
// CAdES signature as an unsigned property.
type CAdESEmbeddedEvidenceRecordBuilder struct {
	// certificateVerifier is the CertificateVerifier to be used for timestamps validation.
	certificateVerifier validation.CertificateVerifier

	// manifestFile is a signature-signed manifest. Used for ASiC.
	manifestFile *model.ManifestFile
}

// NewCAdESEmbeddedEvidenceRecordBuilder is the default constructor.
func NewCAdESEmbeddedEvidenceRecordBuilder(certificateVerifier validation.CertificateVerifier) *CAdESEmbeddedEvidenceRecordBuilder {
	return &CAdESEmbeddedEvidenceRecordBuilder{
		certificateVerifier: validation.NewCertificateVerifierBuilder(certificateVerifier).BuildOfflineCopy(),
	}
}

// SetManifestFile sets a signed manifest file. NOTE: ASiC only. Port of #setManifestFile.
func (b *CAdESEmbeddedEvidenceRecordBuilder) SetManifestFile(manifestFile *model.ManifestFile) {
	b.manifestFile = manifestFile
}

// AddEvidenceRecord adds the evidence record document to a signature with the given signatureId
// (via parameters), provided the evidence record correctly applies to the signature. Port of
// #addEvidenceRecord(DSSDocument, DSSDocument, CAdESEvidenceRecordIncorporationParameters).
//
// Panics when signatureDocument, evidenceRecordDocument or parameters is nil (Java
// Objects.requireNonNull).
func (b *CAdESEmbeddedEvidenceRecordBuilder) AddEvidenceRecord(signatureDocument, evidenceRecordDocument model.DSSDocument,
	parameters *CAdESEvidenceRecordIncorporationParameters) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}
	if evidenceRecordDocument == nil {
		panic("Evidence record document must be provided!")
	}
	if parameters == nil {
		panic("CAdESEvidenceRecordIncorporationParameters must be provided!")
	}

	documentAnalyzer := b.initDocumentAnalyzer(signatureDocument, parameters.DetachedContents())

	sig, err := b.getCAdESSignature(documentAnalyzer, parameters.SignatureId())
	if err != nil {
		return nil, err
	}
	if err := b.assertSignatureExtensionPossible(sig, parameters); err != nil {
		return nil, err
	}

	unsignedAttribute, err := b.getUnsignedAttributeToEmbed(sig, parameters)
	if err != nil {
		return nil, err
	}
	evidenceRecord, err := b.getEvidenceRecord(evidenceRecordDocument, sig, unsignedAttribute, parameters.DetachedContents())
	if err != nil {
		return nil, err
	}
	if err := b.assertEvidenceRecordValid(evidenceRecord, parameters); err != nil {
		return nil, err
	}

	var newSignerInformationList []*cmscore.SignerInfo
	for _, currentSignature := range documentAnalyzer.Signatures() {
		cadesSignature, ok := currentSignature.(*CAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", currentSignature)
		}
		if sig.ID() == cadesSignature.ID() {
			newSignerInformation, err := b.addEvidenceRecordUnsignedProperty(cadesSignature, evidenceRecord, unsignedAttribute)
			if err != nil {
				return nil, err
			}
			newSignerInformationList = append(newSignerInformationList, newSignerInformation)
		} else {
			newSignerInformationList = append(newSignerInformationList, cadesSignature.SignerInformation())
		}
	}
	newCMS, err := cms.CMSUtilsReplaceSigners(sig.CMS(), newSignerInformationList)
	if err != nil {
		return nil, err
	}
	return model.NewInMemoryDocument(newCMS.Encoded()), nil // preserve original coding
}

// getCAdESSignature gets a signature to incorporate evidence record into. Port of the protected
// #getCAdESSignature(DefaultDocumentAnalyzer, String).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getCAdESSignature(documentAnalyzer *CMSDocumentAnalyzer, signatureId string) (*CAdESSignature, error) {
	if signatureId != "" {
		sig := documentAnalyzer.SignatureByID(signatureId)
		if sig == nil {
			return nil, fmt.Errorf("unable to find a signature with Id : %s!", signatureId)
		}
		cadesSignature, ok := sig.(*CAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", sig)
		}
		return cadesSignature, nil
	}

	signatures := documentAnalyzer.Signatures()
	switch {
	case len(signatures) == 0:
		return nil, exception.NewIllegalInputException(fmt.Sprintf("No signatures found in the document with name '%s'", documentAnalyzer.Document().Name()))
	case len(signatures) > 1:
		return nil, fmt.Errorf("more than one signature found in a document with name '%s'! "+
			"Please provide a signatureId within the parameters", documentAnalyzer.Document().Name())
	}
	// if one signature
	cadesSignature, ok := signatures[0].(*CAdESSignature)
	if !ok {
		return nil, fmt.Errorf("unexpected signature type %T", signatures[0])
	}
	return cadesSignature, nil
}

// getUnsignedAttributeToEmbed ports the private
// getUnsignedAttributeToEmbed(CAdESSignature, CAdESEvidenceRecordIncorporationParameters).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getUnsignedAttributeToEmbed(sig *CAdESSignature,
	parameters *CAdESEvidenceRecordIncorporationParameters) (*CAdESAttribute, error) {
	if parameters.IsParallelEvidenceRecord() {
		unsignedAttributes := CAdESUnsignedAttributesBuild(sig.SignerInformation())
		if unsignedAttributes.IsExist() {
			attributes := unsignedAttributes.Attributes()
			lastUnsignedAttribute := attributes[len(attributes)-1]
			if lastUnsignedAttribute.IsEvidenceRecord() {
				expectedEvidenceRecordAttributeType := b.getEvidenceRecordUnsignedPropertyOID(sig)
				if !expectedEvidenceRecordAttributeType.Equal(lastUnsignedAttribute.ASN1Oid()) {
					return nil, exception.NewIllegalInputException(fmt.Sprintf(
						"Unable to embed the parallel evidence record. Expected type '%s', obtained type '%s'.",
						lastUnsignedAttribute.ASN1Oid().String(), expectedEvidenceRecordAttributeType.String()))
				}
				return lastUnsignedAttribute, nil
			}
		}
	}
	// new CAdESAttribute to be created
	return nil, nil
}

// getEvidenceRecord ports the private getEvidenceRecord(DSSDocument, CAdESSignature,
// CAdESAttribute, List<DSSDocument>).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getEvidenceRecord(evidenceRecordDocument model.DSSDocument, sig *CAdESSignature,
	unsignedAttribute *CAdESAttribute, detachedContents []model.DSSDocument) (validation.EvidenceRecord, error) {
	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(evidenceRecordDocument)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("Unable to build an evidence record from the provided document. Reason : %s", err.Error()), err)
	}

	embeddedEvidenceRecordHelper := NewCAdESEmbeddedEvidenceRecordHelper(sig, unsignedAttribute)
	if len(detachedContents) > 0 {
		evidenceRecordAnalyzer.SetEvidenceRecordIncorporationType(enumerations.EvidenceRecordIncorporationType_EXTERNAL_EVIDENCE_RECORD)
		embeddedEvidenceRecordHelper.SetDetachedContents(detachedContents)
	}
	evidenceRecordAnalyzer.SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper)

	return evidenceRecordAnalyzer.EvidenceRecord(), nil
}

// validateTimestamps ports the private #validateTimestamps(EvidenceRecord).
func (b *CAdESEmbeddedEvidenceRecordBuilder) validateTimestamps(evidenceRecord validation.EvidenceRecord) {
	validationContext := validation.NewSignatureValidationContext()
	validationContext.Initialize(b.certificateVerifier)

	validationContext.AddDocumentCertificateSource(evidenceRecord.CertificateSource())
	for _, timestampToken := range evidenceRecord.Timestamps() {
		validationContext.AddTimestampTokenForVerification(timestampToken)
	}

	validationContext.Validate()

	signatureValidationAlerter := validation.NewSignatureValidationAlerter(validationContext)
	signatureValidationAlerter.SetSigningOperation(enumerations.SigningOperation_ADD_EVIDENCE_RECORD)
	signatureValidationAlerter.AssertAllTimestampsValid()
}

// initDocumentAnalyzer ports the private #initDocumentAnalyzer(DSSDocument, List<DSSDocument>).
func (b *CAdESEmbeddedEvidenceRecordBuilder) initDocumentAnalyzer(signatureDocument model.DSSDocument, detachedContents []model.DSSDocument) *CMSDocumentAnalyzer {
	documentAnalyzer, err := NewCMSDocumentAnalyzerFromDocument(signatureDocument)
	if err != nil {
		panic(err)
	}
	documentAnalyzer.SetManifestFile(b.manifestFile)
	documentAnalyzer.SetDetachedContents(detachedContents)
	return documentAnalyzer
}

// addEvidenceRecordUnsignedProperty ports the private
// #addEvidenceRecordUnsignedProperty(CAdESSignature, EvidenceRecord, CAdESAttribute).
func (b *CAdESEmbeddedEvidenceRecordBuilder) addEvidenceRecordUnsignedProperty(sig *CAdESSignature, evidenceRecord validation.EvidenceRecord,
	unsignedAttribute *CAdESAttribute) (*cmscore.SignerInfo, error) {
	attributeOID := b.getEvidenceRecordUnsignedPropertyOID(sig)
	evidenceRecordAttribute, err := b.getEvidenceRecordAttribute(evidenceRecord, attributeOID, unsignedAttribute)
	if err != nil {
		return nil, err
	}

	signerInformation := sig.SignerInformation()
	unsignedAttributesWithER := b.getUnsignedPropertiesTable(signerInformation, evidenceRecordAttribute, unsignedAttribute != nil)
	return cms.CMSUtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributesWithER)
}

// getEvidenceRecordAttribute ports the private #getEvidenceRecordAttribute(EvidenceRecord,
// ASN1ObjectIdentifier, CAdESAttribute).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getEvidenceRecordAttribute(evidenceRecord validation.EvidenceRecord,
	attributeOID asn1.ObjectIdentifier, unsignedAttribute *CAdESAttribute) (*cmscore.Attribute, error) {
	asn1EvidenceRecord, err := b.getASN1EvidenceRecord(evidenceRecord)
	if err != nil {
		return nil, err
	}
	if unsignedAttribute != nil {
		// existing unsigned property
		values := make([][]byte, 0, len(unsignedAttribute.AttrValues())+1)
		for _, v := range unsignedAttribute.AttrValues() {
			values = append(values, v.Encoded())
		}
		values = append(values, asn1EvidenceRecord)
		return cmscore.NewAttribute(attributeOID, values...), nil
	}
	// new unsigned property
	return cmscore.NewAttribute(attributeOID, asn1EvidenceRecord), nil
}

// getUnsignedPropertiesTable ports the private #getUnsignedPropertiesTable(SignerInformation,
// Attribute, boolean).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getUnsignedPropertiesTable(signerInformation *cmscore.SignerInfo,
	evidenceRecordAttribute *cmscore.Attribute, parallelER bool) cmscore.Attributes {
	unsignedAttributes := CAdESUtilsUnsignedAttributes(signerInformation)
	originalAttributeTableLength := len(unsignedAttributes)
	if parallelER {
		originalAttributeTableLength--
	}

	attrs := make(cmscore.Attributes, 0, originalAttributeTableLength+1)
	for i := 0; i < originalAttributeTableLength; i++ {
		attrs = append(attrs, unsignedAttributes[i])
	}
	attrs = append(attrs, evidenceRecordAttribute)
	return attrs
}

// getASN1EvidenceRecord ports the private #getASN1EvidenceRecord(EvidenceRecord).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getASN1EvidenceRecord(evidenceRecord validation.EvidenceRecord) ([]byte, error) {
	element, rest, err := asn1ber.Parse(evidenceRecord.Encoded())
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 || !element.IsUniversal(asn1ber.TagSequence) {
		return nil, fmt.Errorf("the evidence record is not an ASN1Sequence")
	}
	return element.Encoded(), nil
}

// getEvidenceRecordUnsignedPropertyOID ports the private
// #getEvidenceRecordUnsignedPropertyOID(CAdESSignature).
func (b *CAdESEmbeddedEvidenceRecordBuilder) getEvidenceRecordUnsignedPropertyOID(sig *CAdESSignature) asn1.ObjectIdentifier {
	if sig.CMS().IsDetachedSignature() {
		return spi.OID_id_aa_er_external
	}
	return spi.OID_id_aa_er_internal
}

// assertEvidenceRecordValid ports the private #assertEvidenceRecordValid(EvidenceRecord,
// CAdESEvidenceRecordIncorporationParameters).
func (b *CAdESEmbeddedEvidenceRecordBuilder) assertEvidenceRecordValid(evidenceRecord validation.EvidenceRecord,
	parameters *CAdESEvidenceRecordIncorporationParameters) error {
	if enumerations.EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD != evidenceRecord.EvidenceRecordType() {
		return exception.NewIllegalInputException(fmt.Sprintf("Only RFC 4998 ERS type of Evidence Records is allowed "+
			"for CAdES signatures! Identified type of evidence record: '%s'", evidenceRecord.EvidenceRecordType()))
	}
	for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
		if !referenceValidation.IsIntact() {
			switch referenceValidation.Type() {
			case enumerations.DigestMatcherType_EVIDENCE_RECORD_MASTER_SIGNATURE:
				return exception.NewIllegalInputException("The digest covered by the evidence record do not correspond to " +
					"the digest computed on the signature!")
			case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT:
				if len(parameters.DetachedContents()) == 0 {
					return exception.NewIllegalInputException("The digest covered by the evidence record do not correspond to " +
						"the digest computed on the detached content! " +
						"Please use #setDetachedContent method to provide original documents.")
				}
				return exception.NewIllegalInputException("The digest covered by the evidence record do not correspond to " +
					"the digest computed on the detached content!")
			case enumerations.DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE:
				// acceptable status
			default:
				panic(fmt.Sprintf("Unexpected digest matcher type '%s' does not correspond to the value present in "+
					"the evidence record!", referenceValidation.Type()))
			}
		}
	}
	b.validateTimestamps(evidenceRecord)
	return nil
}

// assertSignatureExtensionPossible ports the private
// #assertSignatureExtensionPossible(CAdESSignature, CAdESEvidenceRecordIncorporationParameters).
func (b *CAdESEmbeddedEvidenceRecordBuilder) assertSignatureExtensionPossible(sig *CAdESSignature,
	parameters *CAdESEvidenceRecordIncorporationParameters) error {
	if err := cms.CMSUtilsAssertEvidenceRecordEmbeddingSupported(); err != nil {
		return err
	}
	if err := b.assertNoEvidenceRecordsInOtherSignerInfos(sig); err != nil {
		return err
	}

	if CAdESUtilsContainsATSTv2(sig.SignerInformation()) {
		return exception.NewIllegalInputException("Cannot add evidence record to a CAdES containing an archiveTimestampV2")
	}
	if sig.CMS().IsDetachedSignature() && len(parameters.DetachedContents()) != 1 {
		return fmt.Errorf("one and only one detached document is allowed for an embedded evidence record in CAdES")
	}
	return nil
}

// assertNoEvidenceRecordsInOtherSignerInfos ports the private
// #assertNoEvidenceRecordsInOtherSignerInfos(CAdESSignature).
func (b *CAdESEmbeddedEvidenceRecordBuilder) assertNoEvidenceRecordsInOtherSignerInfos(sig *CAdESSignature) error {
	for _, signerInfo := range sig.CMS().SignerInfos() {
		if sig.SignerInformation() != signerInfo && CAdESUtilsContainsEvidenceRecord(signerInfo) {
			return exception.NewIllegalInputException("At most one of the SignerInfo instances within " +
				"the SignedData instance shall contain evidence-records attributes! " +
				"Please abolish the operation or provide another signature Id.")
		}
	}
	return nil
}

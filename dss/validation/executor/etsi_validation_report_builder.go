// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/signature/ETSIValidationReportBuilder.java
// (DSS 6.5.RC1) - the ETSI TS 119 102-2 validation report builder.
//
// Shape differences against Java, all forced:
//
//   - Java's ObjectFactory only news up the generated JAXB types; the Go
//     model has no factory, so every createXxxType() becomes a composite
//     literal.
//
//   - The 35-member SignatureAttributes choice and the two representation
//     choices are JAXBElement<?> lists in Java. The Go model records each
//     member as a name-tagged jaxb choice item; the element names below are
//     exactly the ones the createSignatureAttributesTypeXxx factory methods
//     produce, so the marshalled order and naming are unchanged.
//
//   - Java's VOReferenceType.getVOReference() holds the referenced JAXB
//     objects and the marshaller writes out their xs:ID; the Go model stores
//     the ids directly (IDREFS), so voReference takes the ids.
//
//   - getIndividualValidationConstraintReport takes an XmlConstraintsConclusion
//     (the JAXB supertype). The generated Go model has no supertype, so the
//     port passes the embedded *XmlConstraintsConclusionContent, whose nil
//     value stands for Java's null block.
//
//   - Java's List.contains(orphanToken) uses OrphanTokenWrapper.equals, which
//     compares getId(); containsOrphanId does the same.
//
//   - wrap() throws IllegalArgumentException for an unsupported time-stamp
//     type; the Go port panics with the same message, matching the ported
//     process tree's convention.

package executor

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/utain/esig/dss/detailedreport"
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
	"github.com/utain/esig/dss/validationreport/jaxb"
)

// ETSIValidationReportBuilder builds the ETSI Validation report. Port of
// ETSIValidationReportBuilder.
type ETSIValidationReportBuilder struct {
	// currentTime is the validation time.
	currentTime time.Time

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// detailedReport is the detailed report.
	detailedReport *detailedreport.DetailedReport

	// signatureIdentifierMap maps signature Ids to the respective
	// SignatureIdentifierTypes.
	signatureIdentifierMap map[string]*jaxb.SignatureIdentifierType

	// validationObjectMap maps token Ids to the respective
	// ValidationObjectTypes.
	validationObjectMap map[string]*jaxb.ValidationObjectType
}

// NewETSIValidationReportBuilder is the default constructor. Port of
// ETSIValidationReportBuilder(Date, DiagnosticData, DetailedReport).
func NewETSIValidationReportBuilder(currentTime time.Time, diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport) *ETSIValidationReportBuilder {
	return &ETSIValidationReportBuilder{
		currentTime:            currentTime,
		diagnosticData:         diagnosticData,
		detailedReport:         detailedReport,
		signatureIdentifierMap: make(map[string]*jaxb.SignatureIdentifierType),
		validationObjectMap:    make(map[string]*jaxb.ValidationObjectType),
	}
}

// Build builds the ValidationReportType. Port of build().
func (b *ETSIValidationReportBuilder) Build() *jaxb.ValidationReportType {
	result := &jaxb.ValidationReportType{}

	if len(b.diagnosticData.Signatures()) != 0 {
		// iterate over the complete list of signatures, including counter signatures
		for _, sigWrapper := range b.diagnosticData.Signatures() {
			result.SignatureValidationReport = append(result.SignatureValidationReport, b.signatureValidationReport(sigWrapper))
		}
	} else {
		result.SignatureValidationReport = append(result.SignatureValidationReport, b.noSignatureFoundReport())
	}

	signatureValidationObjects := b.signatureValidationObjects()
	if len(signatureValidationObjects.ValidationObject) != 0 {
		result.SignatureValidationObjects = signatureValidationObjects
	}

	return result
}

// signatureValidationReport is the port of the private
// getSignatureValidationReport(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signatureValidationReport(sigWrapper *diagnostic.SignatureWrapper) *jaxb.SignatureValidationReportType {
	signatureValidationReport := &jaxb.SignatureValidationReportType{}
	signatureValidationReport.SignatureIdentifier = b.signatureIdentifier(sigWrapper)
	signatureValidationReport.SignersDocument = b.signersDocument(sigWrapper)
	signatureAttributes := b.signatureAttributes(sigWrapper)
	if len(signatureAttributes.Items) != 0 {
		signatureValidationReport.SignatureAttributes = signatureAttributes
	}
	signatureValidationReport.SignerInformation = b.signerInformation(sigWrapper)
	signatureValidationReport.SignatureQuality = b.signatureQuality(sigWrapper)
	signatureValidationReport.SignatureValidationProcess = b.signatureValidationProcess(sigWrapper)
	signatureValidationReport.SignatureValidationStatus = *b.tokenValidationStatus(sigWrapper)
	signatureValidationReport.ValidationTimeInfo = b.validationTimeInfo(sigWrapper)
	signatureValidationReport.ValidationConstraintsEvaluationReport = b.validationConstraintsEvaluationReport(sigWrapper)
	return signatureValidationReport
}

// noSignatureFoundReport is the port of the private noSignatureFoundReport().
func (b *ETSIValidationReportBuilder) noSignatureFoundReport() *jaxb.SignatureValidationReportType {
	signatureValidationReport := &jaxb.SignatureValidationReportType{}
	signatureValidationReport.SignatureValidationStatus = *b.noSignatureFoundValidationStatus()
	return signatureValidationReport
}

// noSignatureFoundValidationStatus is the port of the private
// noSignatureFoundValidationStatus().
func (b *ETSIValidationReportBuilder) noSignatureFoundValidationStatus() *jaxb.ValidationStatusType {
	validationStatus := &jaxb.ValidationStatusType{}
	validationStatus.MainIndication = jaxb.URIIndication(enumerations.Indication_NO_SIGNATURE_FOUND)
	return validationStatus
}

// validationConstraintsEvaluationReport is the port of the private
// getValidationConstraintsEvaluationReport(AbstractTokenProxy).
func (b *ETSIValidationReportBuilder) validationConstraintsEvaluationReport(
	token diagnostic.TokenProxy) *jaxb.ValidationConstraintsEvaluationReportType {
	validationConstraintsEvaluationReport := &jaxb.ValidationConstraintsEvaluationReportType{}
	bbbResults := b.detailedReport.BasicBuildingBlockById(token.Id())

	// The generated Go model has no XmlConstraintsConclusion supertype; the
	// embedded content struct (nil when the block itself is absent) stands in.
	var fc, isc, vci, cv, sav, xcv, psv, pcv, vts *drjaxb.XmlConstraintsConclusionContent
	if bbbResults.FC != nil {
		fc = &bbbResults.FC.XmlConstraintsConclusionContent
	}
	if bbbResults.ISC != nil {
		isc = &bbbResults.ISC.XmlConstraintsConclusionContent
	}
	if bbbResults.VCI != nil {
		vci = &bbbResults.VCI.XmlConstraintsConclusionContent
	}
	if bbbResults.CV != nil {
		cv = &bbbResults.CV.XmlConstraintsConclusionContent
	}
	if bbbResults.SAV != nil {
		sav = &bbbResults.SAV.XmlConstraintsConclusionContent
	}
	if bbbResults.XCV != nil {
		xcv = &bbbResults.XCV.XmlConstraintsConclusionContent
	}

	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_FORMAT_CHECKING, fc)
	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE, isc)
	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_VALIDATION_CONTEXT_INITIALIZATION, vci)
	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_CRYPTOGRAPHIC_VERIFICATION, cv)
	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_SIGNATURE_ACCEPTANCE_VALIDATION, sav)
	b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_X509_CERTIFICATE_VALIDATION, xcv)
	switch token.(type) {
	case *diagnostic.SignatureWrapper, *diagnostic.TimestampWrapper:
		if bbbResults.PSV != nil {
			psv = &bbbResults.PSV.XmlConstraintsConclusionContent
		}
		if bbbResults.PCV != nil {
			pcv = &bbbResults.PCV.XmlConstraintsConclusionContent
		}
		if bbbResults.VTS != nil {
			vts = &bbbResults.VTS.XmlConstraintsConclusionContent
		}
		b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_PAST_SIGNATURE_VALIDATION, psv)
		b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_PAST_CERTIFICATE_VALIDATION, pcv)
		b.addBBB(validationConstraintsEvaluationReport, process.BasicBuildingBlockDefinition_VALIDATION_TIME_SLIDING, vts)
	}
	return validationConstraintsEvaluationReport
}

// addBBB is the port of the private
// addBBB(ValidationConstraintsEvaluationReportType, BasicBuildingBlockDefinition,
// XmlConstraintsConclusion).
func (b *ETSIValidationReportBuilder) addBBB(
	validationConstraintsEvaluationReport *jaxb.ValidationConstraintsEvaluationReportType,
	bbbUri process.BasicBuildingBlockDefinition, constraintConclusion *drjaxb.XmlConstraintsConclusionContent) {
	if constraintConclusion != nil {
		validationConstraintsEvaluationReport.ValidationConstraint = append(
			validationConstraintsEvaluationReport.ValidationConstraint,
			b.individualValidationConstraintReport(bbbUri, constraintConclusion, b.applied()))
	} else {
		validationConstraintsEvaluationReport.ValidationConstraint = append(
			validationConstraintsEvaluationReport.ValidationConstraint,
			b.individualValidationConstraintReport(bbbUri, nil, b.disabled()))
	}
}

// individualValidationConstraintReport is the port of the private
// getIndividualValidationConstraintReport(BasicBuildingBlockDefinition,
// XmlConstraintsConclusion, ConstraintStatusType).
func (b *ETSIValidationReportBuilder) individualValidationConstraintReport(bbbUri process.BasicBuildingBlockDefinition,
	constraintConclusion *drjaxb.XmlConstraintsConclusionContent,
	constraintStatusType *jaxb.ConstraintStatusType) *jaxb.IndividualValidationConstraintReportType {
	validationConstraint := &jaxb.IndividualValidationConstraintReportType{}
	validationConstraint.ValidationConstraintIdentifier = bbbUri.URI()
	validationConstraint.ConstraintStatus = *constraintStatusType
	if constraintConclusion != nil {
		validationConstraint.ValidationStatus = b.conclusionValidationStatus(constraintConclusion.Conclusion)
	}
	return validationConstraint
}

// applied is the port of the private applied().
func (b *ETSIValidationReportBuilder) applied() *jaxb.ConstraintStatusType {
	return b.constraintStatus(jaxb.ConstraintStatus_APPLIED)
}

// disabled is the port of the private disabled().
func (b *ETSIValidationReportBuilder) disabled() *jaxb.ConstraintStatusType {
	return b.constraintStatus(jaxb.ConstraintStatus_DISABLED)
}

// constraintStatus is the port of the private
// constraintStatus(ConstraintStatus).
func (b *ETSIValidationReportBuilder) constraintStatus(status jaxb.ConstraintStatus) *jaxb.ConstraintStatusType {
	constraintStatus := &jaxb.ConstraintStatusType{}
	constraintStatus.Status = status
	return constraintStatus
}

// tokenValidationReport is the port of the private
// getValidationReport(AbstractTokenProxy).
func (b *ETSIValidationReportBuilder) tokenValidationReport(token diagnostic.TokenProxy) *jaxb.SignatureValidationReportType {
	tokenBBB := b.detailedReport.BasicBuildingBlockById(token.Id())
	// return null if validation was not performed
	if tokenBBB == nil {
		return nil
	}
	signatureValidationReport := &jaxb.SignatureValidationReportType{}
	signatureValidationReport.SignerInformation = b.signerInformation(token)
	signatureValidationReport.SignatureValidationStatus = *b.tokenValidationStatus(token)
	signatureValidationReport.ValidationConstraintsEvaluationReport = b.validationConstraintsEvaluationReport(token)

	timestampQualification := b.detailedReport.TimestampQualification(token.Id())
	if timestampQualification != "" {
		signatureQualityType := &jaxb.SignatureQualityType{}
		signatureQualityType.SignatureQualityInformation = append(
			signatureQualityType.SignatureQualityInformation, timestampQualification.URI())
		signatureValidationReport.SignatureQuality = signatureQualityType
	}

	return signatureValidationReport
}

// validationTimeInfo is the port of the private
// getValidationTimeInfo(SignatureWrapper).
func (b *ETSIValidationReportBuilder) validationTimeInfo(sigWrapper *diagnostic.SignatureWrapper) *jaxb.ValidationTimeInfoType {
	validationTimeInfoType := &jaxb.ValidationTimeInfoType{}
	validationTimeInfoType.ValidationTime = jaxb.XSDateTime(b.currentTime)

	proofOfExistence := b.detailedReport.BestProofOfExistence(sigWrapper.Id())
	poeType := &jaxb.POEType{}
	poeType.POETime = jaxb.XSDateTime(time.Time(proofOfExistence.Time))
	poeType.TypeOfProof = jaxb.TypeOfProof_VALIDATION

	timestampId := ""
	if proofOfExistence.TimestampId != nil {
		timestampId = *proofOfExistence.TimestampId
	}
	if timestampId != "" {
		timestampWrapper := b.diagnosticData.TimestampById(timestampId)
		if timestampWrapper != nil {
			poeType.POEObject = b.voReferenceOfObject(b.timestampValidationObject(timestampWrapper))
		}
		evidenceRecordWrapper := b.diagnosticData.EvidenceRecordById(timestampId)
		if evidenceRecordWrapper != nil {
			poeType.POEObject = b.voReferenceOfObject(b.evidenceRecordValidationObject(evidenceRecordWrapper))
		}
	}
	validationTimeInfoType.BestSignatureTime = *poeType
	return validationTimeInfoType
}

// signerInformation is the port of the private
// getSignerInformation(AbstractTokenProxy).
func (b *ETSIValidationReportBuilder) signerInformation(token diagnostic.TokenProxy) *jaxb.SignerInformationType {
	signingCert := token.SigningCertificate()
	if signingCert == nil {
		return nil
	}
	signerInfo := &jaxb.SignerInformationType{}
	pseudoUseStatus, present := b.pseudoUseStatus(token)
	if present {
		pseudonym := b.isPseudoUse(pseudoUseStatus)
		signerInfo.Pseudonym = &pseudonym
	}
	signer := signingCert.ReadableCertificateName()
	signerInfo.Signer = &signer
	signerInfo.SignerCertificate = *b.voReferenceOfObject(b.certificateValidationObject(signingCert))
	return signerInfo
}

// pseudoUseStatus is the port of the private
// getPseudoUseStatus(AbstractTokenProxy); the second result stands for Java's
// null return, since XmlStatus is a Go string type with no null value.
func (b *ETSIValidationReportBuilder) pseudoUseStatus(token diagnostic.TokenProxy) (drjaxb.XmlStatus, bool) {
	signingCertificateXCV := b.detailedReport.SigningCertificate(token.Id())
	if signingCertificateXCV != nil {
		constraints := signingCertificateXCV.Constraint
		for _, xmlConstraint := range constraints {
			if xmlConstraint.Name != nil && xmlConstraint.Name.Key != nil && *xmlConstraint.Name.Key == "BBB_XCV_PSEUDO_USE" {
				return xmlConstraint.Status, true
			}
		}
	}
	return "", false
}

// isPseudoUse is the port of the private isPseudoUse(XmlStatus).
func (b *ETSIValidationReportBuilder) isPseudoUse(status drjaxb.XmlStatus) bool {
	return status != drjaxb.XmlStatus_OK && status != drjaxb.XmlStatus_IGNORED
}

// signatureQuality is the port of the private
// getSignatureQuality(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signatureQuality(signatureWrapper *diagnostic.SignatureWrapper) *jaxb.SignatureQualityType {
	signatureQualification := b.detailedReport.SignatureQualification(signatureWrapper.Id())
	if signatureQualification != "" {
		signatureQualityType := &jaxb.SignatureQualityType{}
		signatureQualityType.SignatureQualityInformation = append(
			signatureQualityType.SignatureQualityInformation, signatureQualification.URI())
		return signatureQualityType
	}
	return nil
}

// voReferenceOfSignature is the port of the private
// getVOReference(SignatureIdentifierType).
func (b *ETSIValidationReportBuilder) voReferenceOfSignature(signatureIdentifier *jaxb.SignatureIdentifierType) *jaxb.VOReferenceType {
	voRef := &jaxb.VOReferenceType{}
	voRef.VOReference = append(voRef.VOReference, signatureIdentifier.Id)
	return voRef
}

// voReferenceOfObject is the port of the private
// getVOReference(ValidationObjectType).
func (b *ETSIValidationReportBuilder) voReferenceOfObject(validationObject *jaxb.ValidationObjectType) *jaxb.VOReferenceType {
	return b.voReference([]*jaxb.ValidationObjectType{validationObject})
}

// voReference is the port of the private
// getVOReference(List<ValidationObjectType>).
func (b *ETSIValidationReportBuilder) voReference(validationObjects []*jaxb.ValidationObjectType) *jaxb.VOReferenceType {
	voRef := &jaxb.VOReferenceType{}
	for _, validationObject := range validationObjects {
		voRef.VOReference = append(voRef.VOReference, validationObject.Id)
	}
	return voRef
}

// signatureValidationProcess is the port of the private
// getSignatureValidationProcess(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signatureValidationProcess(sigWrapper *diagnostic.SignatureWrapper) *jaxb.SignatureValidationProcessType {
	validationProcess := &jaxb.SignatureValidationProcessType{}
	processId := b.currentProcessId(sigWrapper)
	validationProcess.SignatureValidationProcessID = &processId
	return validationProcess
}

// currentProcessId is the port of the private
// getCurrentProcessId(SignatureWrapper).
func (b *ETSIValidationReportBuilder) currentProcessId(sigWrapper *diagnostic.SignatureWrapper) jaxb.SignatureValidationProcessID {
	processId := jaxb.SignatureValidationProcessID_BASIC
	indicationLTA := b.detailedReport.ArchiveDataValidationIndication(sigWrapper.Id())
	indicationLTVM := b.detailedReport.LongTermValidationIndication(sigWrapper.Id())
	if indicationLTA != "" {
		processId = jaxb.SignatureValidationProcessID_LTA
	} else if indicationLTVM != "" {
		processId = jaxb.SignatureValidationProcessID_LTVM
	}
	return processId
}

// signatureValidationObjects is the port of the private
// getSignatureValidationObjects().
func (b *ETSIValidationReportBuilder) signatureValidationObjects() *jaxb.ValidationObjectListType {
	validationObjectListType := &jaxb.ValidationObjectListType{}

	poeExtraction := vpfswatsp.NewPOEExtraction()
	poeExtraction.Init(b.diagnosticData, b.currentTime)

	// 1. Extract POEs
	evidenceRecords := b.diagnosticData.EvidenceRecords()
	for _, evidenceRecord := range evidenceRecords {
		if enumerations.Indication_PASSED == b.detailedReport.EvidenceRecordValidationIndication(evidenceRecord.Id()) {
			poeExtraction.ExtractEvidenceRecordPOE(evidenceRecord)
		}
	}

	timestampList := b.diagnosticData.NonEvidenceRecordTimestamps()
	poeExtraction.CollectAllPOE(timestampList)

	// 2. Add validation object types
	for _, evidenceRecord := range evidenceRecords {
		evidenceRecordValidationObject := b.evidenceRecordValidationObject(evidenceRecord)
		evidenceRecordValidationObject.POE = b.poe(evidenceRecord.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, evidenceRecordValidationObject)
	}

	for _, timestamp := range b.diagnosticData.TimestampList() {
		timestampValidationObject := b.timestampValidationObject(timestamp)
		timestampValidationObject.POE = b.poe(timestamp.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, timestampValidationObject)
	}

	for _, eaa := range b.diagnosticData.EAAs() {
		eaaValidationObject := b.eaaValidationObject(eaa)
		eaaValidationObject.POE = b.poe(eaa.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, eaaValidationObject)
	}

	for _, certificate := range b.diagnosticData.UsedCertificates() {
		certificateValidationObject := b.certificateValidationObject(certificate)
		certificateValidationObject.POE = b.poe(certificate.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, certificateValidationObject)
	}

	for _, orphanCertificate := range b.diagnosticData.AllOrphanCertificateObjects() {
		orphanCertificateValidationObject := b.orphanCertificateValidationObject(orphanCertificate)
		orphanCertificateValidationObject.POE = b.poe(orphanCertificate.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, orphanCertificateValidationObject)
	}

	for _, revocationData := range JavaHashSetOrder(b.diagnosticData.AllRevocationData()) {
		revocationValidationObject := b.revocationValidationObject(revocationData)
		revocationValidationObject.POE = b.poe(revocationData.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, revocationValidationObject)
	}

	for _, orphanRevocation := range b.diagnosticData.AllOrphanRevocationObjects() {
		orphanRevocationValidationObject := b.orphanRevocationValidationObject(orphanRevocation)
		orphanRevocationValidationObject.POE = b.poe(orphanRevocation.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, orphanRevocationValidationObject)
	}

	for _, signedData := range b.diagnosticData.AllSignerDocuments() {
		signerDataValidationObject := b.signerDataValidationObject(signedData)
		signerDataValidationObject.POE = b.poe(signedData.Id(), poeExtraction)
		validationObjectListType.ValidationObject = append(validationObjectListType.ValidationObject, signerDataValidationObject)
	}

	return validationObjectListType
}

// certificateValidationObject is the port of the private
// getCertificateValidationObject(CertificateWrapper).
func (b *ETSIValidationReportBuilder) certificateValidationObject(certificate *diagnostic.CertificateWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[certificate.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[certificate.Id()] = validationObject

		validationObject.Id = certificate.Id()
		validationObject.ObjectType = jaxb.ObjectType_CERTIFICATE
		representation := jaxb.ValidationObjectRepresentationType{}
		if len(certificate.Binaries()) > 0 {
			representation.Items = append(representation.Items, base64Item(certificate.Binaries()))
		} else {
			representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(certificate.DigestAlgoAndValue())))
		}
		validationObject.ValidationObjectRepresentation = representation
	}
	return validationObject
}

// digestAlgAndValueType is the port of the private
// getDigestAlgAndValueType(XmlDigestAlgoAndValue).
func (b *ETSIValidationReportBuilder) digestAlgAndValueType(xmlDigestAlgoAndValue *diagnosticjaxb.XmlDigestAlgoAndValue) *jaxb.DigestAlgAndValueType {
	digestAlgAndValueType := &jaxb.DigestAlgAndValueType{}
	digestAlgAndValueType.DigestMethod = *b.digestMethodType(digestAlgorithmOf(xmlDigestAlgoAndValue))
	digestAlgAndValueType.DigestValue = jaxb.DSDigestValue(digestValueOf(xmlDigestAlgoAndValue))
	return digestAlgAndValueType
}

// urn is the port of the private getUrn(DigestAlgorithm).
func (b *ETSIValidationReportBuilder) urn(digestAlgorithm enumerations.DigestAlgorithm) string {
	if digestAlgorithm != "" {
		if digestAlgorithm.URI() != "" {
			return digestAlgorithm.URI()
		}
		if digestAlgorithm.OID() != "" {
			return process.ToUrnOid(digestAlgorithm.OID())
		}
	}
	return "?"
}

// poe is the port of the private getPOE(String, POEExtraction).
func (b *ETSIValidationReportBuilder) poe(tokenId string, poeExtraction *vpfswatsp.POEExtraction) *jaxb.POEType {
	poeType := &jaxb.POEType{}
	lowestPOE := poeExtraction.GetLowestPOE(tokenId)
	poeType.POETime = jaxb.XSDateTime(lowestPOE.Time())
	switch lowestPOE.(type) {
	case *vpfswatsp.TimestampPOE:
		// Java's Utils.isStringNotEmpty(getPOEProviderId()) - the ported POE
		// returns *string, whose nil is Java's null.
		if timestampId := lowestPOE.POEProviderId(); timestampId != nil && *timestampId != "" {
			timestampWrapper := b.diagnosticData.TimestampById(*timestampId)
			poeType.POEObject = b.voReferenceOfObject(b.timestampValidationObject(timestampWrapper))
		}
	case *vpfswatsp.EvidenceRecordPOE:
		if evidenceRecordId := lowestPOE.POEProviderId(); evidenceRecordId != nil && *evidenceRecordId != "" {
			evidenceRecord := b.diagnosticData.EvidenceRecordById(*evidenceRecordId)
			poeType.POEObject = b.voReferenceOfObject(b.evidenceRecordValidationObject(evidenceRecord))
		}
	}
	poeType.TypeOfProof = jaxb.TypeOfProof_VALIDATION
	return poeType
}

// evidenceRecordValidationObject is the port of the private
// getEvidenceRecordValidationObject(EvidenceRecordWrapper).
func (b *ETSIValidationReportBuilder) evidenceRecordValidationObject(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[evidenceRecord.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[evidenceRecord.Id()] = validationObject

		validationObject.Id = evidenceRecord.Id()
		validationObject.ObjectType = jaxb.ObjectType_EVIDENCE_RECORD
		representation := jaxb.ValidationObjectRepresentationType{}
		if len(evidenceRecord.Binaries()) > 0 {
			representation.Items = append(representation.Items, base64Item(evidenceRecord.Binaries()))
		} else {
			representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(evidenceRecord.DigestAlgoAndValue())))
		}
		validationObject.ValidationObjectRepresentation = representation
		validationObject.POEProvisioning = b.evidenceRecordPOEProvisioningType(evidenceRecord)
		validationObject.ValidationReport = b.evidenceRecordValidationReport(evidenceRecord)
	}
	return validationObject
}

// evidenceRecordPOEProvisioningType is the port of the private
// getPOEProvisioningType(EvidenceRecordWrapper).
func (b *ETSIValidationReportBuilder) evidenceRecordPOEProvisioningType(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.POEProvisioningType {
	poeProvisioning := &jaxb.POEProvisioningType{}
	if productionTime := evidenceRecord.FirstTimestamp().ProductionTime(); productionTime != nil {
		poeProvisioning.POETime = jaxb.XSDateTime(*productionTime)
	}

	for _, cert := range evidenceRecord.CoveredCertificates() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.certificateValidationObject(cert)))
	}

	// only created validation objects must be added (not references)
	allOrphanObjectCertificates := b.diagnosticData.AllOrphanCertificateObjects()
	for _, orphanCert := range evidenceRecord.CoveredOrphanCertificates() {
		if containsOrphanCertificateId(allOrphanObjectCertificates, orphanCert.Id()) {
			poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
				b.voReferenceOfObject(b.orphanCertificateValidationObject(orphanCert)))
		}
	}

	for _, revocation := range evidenceRecord.CoveredRevocations() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.revocationValidationObject(revocation)))
	}
	// only created validation objects must be added (not references)
	allOrphanObjectRevocations := b.diagnosticData.AllOrphanRevocationObjects()
	for _, orphanRevocation := range evidenceRecord.CoveredOrphanRevocations() {
		if containsOrphanRevocationId(allOrphanObjectRevocations, orphanRevocation.Id()) {
			poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
				b.voReferenceOfObject(b.orphanRevocationValidationObject(orphanRevocation)))
		}
	}

	for _, er := range evidenceRecord.CoveredEvidenceRecords() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.evidenceRecordValidationObject(er)))
	}

	for _, tst := range evidenceRecord.CoveredTimestamps() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.timestampValidationObject(tst)))
	}

	for _, signerData := range evidenceRecord.CoveredSignedData() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.signerDataValidationObject(signerData)))
	}

	timestampedSignatures := evidenceRecord.CoveredSignatures()
	for _, timestampedSignature := range timestampedSignatures {
		poeProvisioning.SignatureReference = append(poeProvisioning.SignatureReference, b.signatureReference(timestampedSignature))
	}

	return poeProvisioning
}

// eaaValidationObject is the port of the private
// getEAAValidationObject(EAAWrapper).
func (b *ETSIValidationReportBuilder) eaaValidationObject(eaa *diagnostic.EAAWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[eaa.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[eaa.Id()] = validationObject

		validationObject.Id = eaa.Id()
		// TODO (upstream) : no EAA specific type is available
		validationObject.ObjectType = jaxb.ObjectType_OTHER
		representation := jaxb.ValidationObjectRepresentationType{}
		// TODO (upstream) : fill representation (base64/digest) ?
		representation.Items = append(representation.Items, uriItem(b.uri(eaa)))
		validationObject.ValidationObjectRepresentation = representation
		validationObject.ValidationReport = b.eaaValidationReport(eaa)
	}
	return validationObject
}

// uri is the port of the private getURI(EAAWrapper).
func (b *ETSIValidationReportBuilder) uri(eaa *diagnostic.EAAWrapper) string {
	if eaa.Filename() != "" {
		return eaa.Filename()
	}
	return "?"
}

// evidenceRecordValidationReport is the port of the private
// getValidationReport(EvidenceRecordWrapper).
func (b *ETSIValidationReportBuilder) evidenceRecordValidationReport(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.SignatureValidationReportType {
	xmlEvidenceRecord := b.detailedReport.XmlEvidenceRecordById(evidenceRecord.Id())
	// return null if validation was not performed
	if xmlEvidenceRecord == nil {
		return nil
	}
	signatureValidationReport := &jaxb.SignatureValidationReportType{}
	signatureValidationReport.SignatureValidationStatus = *b.evidenceRecordValidationStatus(evidenceRecord)
	return signatureValidationReport
}

// evidenceRecordValidationStatus is the port of the private
// getValidationStatus(EvidenceRecordWrapper).
func (b *ETSIValidationReportBuilder) evidenceRecordValidationStatus(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.ValidationStatusType {
	validationStatus := &jaxb.ValidationStatusType{}
	b.fillIndicationSubIndication(validationStatus, evidenceRecord.Id())
	b.fillMessages(validationStatus, evidenceRecord.Id())
	b.addEvidenceRecordValidationReportData(validationStatus, evidenceRecord)
	return validationStatus
}

// addEvidenceRecordValidationReportData is the port of the private
// addValidationReportData(ValidationStatusType, EvidenceRecordWrapper).
func (b *ETSIValidationReportBuilder) addEvidenceRecordValidationReportData(validationStatus *jaxb.ValidationStatusType,
	evidenceRecord *diagnostic.EvidenceRecordWrapper) {
	validationReportData := b.associatedValidationReportData(validationStatus)
	if enumerations.Indication_PASSED != b.detailedReport.EvidenceRecordValidationIndication(evidenceRecord.Id()) {
		for _, timestampWrapper := range evidenceRecord.TimestampList() {
			if enumerations.Indication_PASSED != b.detailedReport.FinalIndication(timestampWrapper.Id()) {
				timestampValidationObject := b.timestampValidationObject(timestampWrapper)
				validationReportData.RelatedValidationObject = append(validationReportData.RelatedValidationObject,
					b.voReferenceOfObject(timestampValidationObject))
			}
		}
	}
	xmlEvidenceRecord := b.detailedReport.XmlEvidenceRecordById(evidenceRecord.Id())
	validationProcessEvidenceRecord := xmlEvidenceRecord.ValidationProcessEvidenceRecord
	xmlAOV := validationProcessEvidenceRecord.AOV
	cryptographicValidation := process.GetFinalCryptographicValidation(xmlAOV)
	if cryptographicValidation != nil {
		b.fillEvidenceRecordCryptographicInfo(validationReportData, evidenceRecord, cryptographicValidation)
	}
}

// fillEvidenceRecordCryptographicInfo is the port of the private
// fillCryptographicInfo(ValidationReportDataType, EvidenceRecordWrapper,
// XmlCryptographicValidation).
func (b *ETSIValidationReportBuilder) fillEvidenceRecordCryptographicInfo(validationReportData *jaxb.ValidationReportDataType,
	evidenceRecord *diagnostic.EvidenceRecordWrapper, cryptographicValidation *drjaxb.XmlCryptographicValidation) {
	cryptoInformationType := &jaxb.CryptoInformationType{}
	cryptoInformationType.ValidationObjectId = *b.voReferenceOfObject(b.evidenceRecordValidationObject(evidenceRecord))
	cryptoInformationType.SecureAlgorithm = enumerations.Indication_PASSED == cryptographicValidation.Conclusion.Indication.Indication()
	algorithm := cryptographicValidation.Algorithm
	if algorithm != nil {
		cryptoInformationType.Algorithm = algorithm.Uri
	}
	cryptoInformationType.NotAfter = notAfterOf(cryptographicValidation)
	validationReportData.CryptoInformation = cryptoInformationType
}

// eaaValidationReport is the port of the private
// getValidationReport(EAAWrapper).
func (b *ETSIValidationReportBuilder) eaaValidationReport(eaa *diagnostic.EAAWrapper) *jaxb.SignatureValidationReportType {
	xmlEAA := b.detailedReport.XmlEAAById(eaa.Id())
	// return null if validation was not performed
	if xmlEAA == nil {
		return nil
	}
	signatureValidationReport := &jaxb.SignatureValidationReportType{}
	signatureValidationReport.SignatureValidationStatus = *b.eaaValidationStatus(eaa)

	eaaQualifications := b.detailedReport.EAAQualifications(eaa.Id())
	if len(eaaQualifications) > 0 {
		signatureQualityType := &jaxb.SignatureQualityType{}
		for _, eaaQualification := range eaaQualifications {
			signatureQualityType.SignatureQualityInformation = append(
				signatureQualityType.SignatureQualityInformation, eaaQualification.URI())
		}
		signatureValidationReport.SignatureQuality = signatureQualityType
	}

	return signatureValidationReport
}

// eaaValidationStatus is the port of the private
// getValidationStatus(EAAWrapper).
func (b *ETSIValidationReportBuilder) eaaValidationStatus(eaa *diagnostic.EAAWrapper) *jaxb.ValidationStatusType {
	validationStatus := &jaxb.ValidationStatusType{}
	b.fillIndicationSubIndication(validationStatus, eaa.Id())
	b.fillMessages(validationStatus, eaa.Id())
	b.addTokenValidationReportData(validationStatus, eaa)
	return validationStatus
}

// timestampValidationObject is the port of the private
// getTimestampValidationObject(TimestampWrapper).
func (b *ETSIValidationReportBuilder) timestampValidationObject(timestamp *diagnostic.TimestampWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[timestamp.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[timestamp.Id()] = validationObject

		validationObject.Id = timestamp.Id()
		validationObject.ObjectType = jaxb.ObjectType_TIMESTAMP
		representation := jaxb.ValidationObjectRepresentationType{}
		if len(timestamp.Binaries()) > 0 {
			representation.Items = append(representation.Items, base64Item(timestamp.Binaries()))
		} else {
			representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(timestamp.DigestAlgoAndValue())))
		}
		validationObject.ValidationObjectRepresentation = representation
		validationObject.POEProvisioning = b.timestampPOEProvisioningType(timestamp)
		validationObject.ValidationReport = b.tokenValidationReport(timestamp)
	}
	return validationObject
}

// timestampPOEProvisioningType is the port of the private
// getPOEProvisioningType(TimestampWrapper).
func (b *ETSIValidationReportBuilder) timestampPOEProvisioningType(timestamp *diagnostic.TimestampWrapper) *jaxb.POEProvisioningType {
	poeProvisioning := &jaxb.POEProvisioningType{}
	if productionTime := timestamp.ProductionTime(); productionTime != nil {
		poeProvisioning.POETime = jaxb.XSDateTime(*productionTime)
	}

	for _, cert := range timestamp.TimestampedCertificates() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.certificateValidationObject(cert)))
	}
	// only created validation objects must be added (not references)
	allOrphanObjectCertificates := b.diagnosticData.AllOrphanCertificateObjects()
	for _, orphanCert := range timestamp.TimestampedOrphanCertificates() {
		if containsOrphanCertificateId(allOrphanObjectCertificates, orphanCert.Id()) {
			poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
				b.voReferenceOfObject(b.orphanCertificateValidationObject(orphanCert)))
		}
	}

	for _, revocation := range timestamp.TimestampedRevocations() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.revocationValidationObject(revocation)))
	}
	// only created validation objects must be added (not references)
	allOrphanObjectRevocations := b.diagnosticData.AllOrphanRevocationObjects()
	for _, orphanRevocation := range timestamp.TimestampedOrphanRevocations() {
		if containsOrphanRevocationId(allOrphanObjectRevocations, orphanRevocation.Id()) {
			poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
				b.voReferenceOfObject(b.orphanRevocationValidationObject(orphanRevocation)))
		}
	}

	for _, er := range timestamp.TimestampedEvidenceRecords() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.evidenceRecordValidationObject(er)))
	}

	for _, tst := range timestamp.TimestampedTimestamps() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.timestampValidationObject(tst)))
	}

	for _, signerData := range timestamp.TimestampedSignedData() {
		poeProvisioning.ValidationObject = append(poeProvisioning.ValidationObject,
			b.voReferenceOfObject(b.signerDataValidationObject(signerData)))
	}

	timestampedSignatures := timestamp.TimestampedSignatures()
	for _, timestampedSignature := range timestampedSignatures {
		poeProvisioning.SignatureReference = append(poeProvisioning.SignatureReference, b.signatureReference(timestampedSignature))
	}

	return poeProvisioning
}

// signatureReference is the port of the private
// getSignatureReference(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signatureReference(signature *diagnostic.SignatureWrapper) *jaxb.SignatureReferenceType {
	signatureReference := &jaxb.SignatureReferenceType{}
	signatureDigestReference := signature.SignatureDigestReference()
	if signatureDigestReference != nil {
		signatureReference.CanonicalizationMethod = signatureDigestReference.CanonicalizationMethod
		if signatureDigestReference.DigestMethod != nil {
			digestMethod := enumerations.DigestAlgorithm(*signatureDigestReference.DigestMethod).URI()
			signatureReference.DigestMethod = &digestMethod
		}
		if signatureDigestReference.DigestValue != nil {
			signatureReference.DigestValue = jaxb.Base64Binary(*signatureDigestReference.DigestValue)
		}
	} else if signature.FirstFieldName() != "" {
		fieldName := signature.FirstFieldName()
		signatureReference.PAdESFieldName = &fieldName
	}
	return signatureReference
}

// signerDataValidationObject is the port of the private
// getSignerDataValidationObject(SignerDataWrapper).
func (b *ETSIValidationReportBuilder) signerDataValidationObject(signedData *diagnostic.SignerDataWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[signedData.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[signedData.Id()] = validationObject

		validationObject.Id = signedData.Id()
		validationObject.ObjectType = jaxb.ObjectType_SIGNED_DATA
		representation := jaxb.ValidationObjectRepresentationType{}
		representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(signedData.DigestAlgoAndValue())))
		validationObject.ValidationObjectRepresentation = representation
	}
	return validationObject
}

// revocationValidationObject is the port of the private
// getRevocationValidationObject(RevocationWrapper).
func (b *ETSIValidationReportBuilder) revocationValidationObject(revocationData *diagnostic.RevocationWrapper) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[revocationData.Id()]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[revocationData.Id()] = validationObject

		validationObject.Id = revocationData.Id()
		if enumerations.RevocationType_CRL == revocationData.RevocationType() {
			validationObject.ObjectType = jaxb.ObjectType_CRL
		} else {
			validationObject.ObjectType = jaxb.ObjectType_OCSP_RESPONSE
		}
		representation := jaxb.ValidationObjectRepresentationType{}
		if len(revocationData.Binaries()) > 0 {
			representation.Items = append(representation.Items, base64Item(revocationData.Binaries()))
		} else {
			representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(revocationData.DigestAlgoAndValue())))
		}
		// Upstream keeps the "Standard says choice" sourceAddress branch
		// commented out; it is not ported.
		validationObject.ValidationObjectRepresentation = representation
		validationObject.ValidationReport = b.tokenValidationReport(revocationData)
	}
	return validationObject
}

// orphanCertificateValidationObject is the port of the private
// getOrphanCertificateValidationObject(OrphanCertificateTokenWrapper).
func (b *ETSIValidationReportBuilder) orphanCertificateValidationObject(
	orphanCertificate *diagnostic.OrphanCertificateTokenWrapper) *jaxb.ValidationObjectType {
	return b.createOrphanToken(orphanCertificate.Id(), orphanCertificate.Binaries(),
		orphanCertificate.DigestAlgoAndValue(), jaxb.ObjectType_CERTIFICATE)
}

// orphanRevocationValidationObject is the port of the private
// getOrphanRevocationValidationObject(OrphanRevocationTokenWrapper).
func (b *ETSIValidationReportBuilder) orphanRevocationValidationObject(
	orphanRevocation *diagnostic.OrphanRevocationTokenWrapper) *jaxb.ValidationObjectType {
	objectType := jaxb.ObjectType_OCSP_RESPONSE
	if enumerations.RevocationType_CRL == orphanRevocation.RevocationType() {
		objectType = jaxb.ObjectType_CRL
	}
	return b.createOrphanToken(orphanRevocation.Id(), orphanRevocation.Binaries(),
		orphanRevocation.DigestAlgoAndValue(), objectType)
}

// createOrphanToken is the port of the private
// createOrphanToken(OrphanTokenWrapper<?>, ObjectType). Java's OrphanTokenWrapper
// supertype has no Go counterpart, so the three members the method reads are
// passed directly.
func (b *ETSIValidationReportBuilder) createOrphanToken(id string, binaries []byte,
	digestAlgoAndValue *diagnosticjaxb.XmlDigestAlgoAndValue, objectType jaxb.ObjectType) *jaxb.ValidationObjectType {
	validationObject := b.validationObjectMap[id]
	if validationObject == nil {
		validationObject = &jaxb.ValidationObjectType{}
		b.validationObjectMap[id] = validationObject

		validationObject.Id = id
		validationObject.ObjectType = objectType
		representation := jaxb.ValidationObjectRepresentationType{}
		if len(binaries) > 0 {
			representation.Items = append(representation.Items, base64Item(binaries))
		} else {
			representation.Items = append(representation.Items, digestAlgAndValueItem(b.digestAlgAndValueType(digestAlgoAndValue)))
		}
		validationObject.ValidationObjectRepresentation = representation
	}
	return validationObject
}

// tokenValidationStatus is the port of the private
// getValidationStatus(AbstractTokenProxy).
func (b *ETSIValidationReportBuilder) tokenValidationStatus(token diagnostic.TokenProxy) *jaxb.ValidationStatusType {
	validationStatus := &jaxb.ValidationStatusType{}
	b.fillIndicationSubIndication(validationStatus, token.Id())
	b.fillMessages(validationStatus, token.Id())
	b.addTokenValidationReportData(validationStatus, token)
	return validationStatus
}

// fillIndicationSubIndication is the port of the private
// fillIndicationSubIndication(ValidationStatusType, String).
func (b *ETSIValidationReportBuilder) fillIndicationSubIndication(validationStatus *jaxb.ValidationStatusType, tokenId string) {
	finalIndication := b.detailedReport.FinalIndication(tokenId)
	if finalIndication != "" {
		validationStatus.MainIndication = jaxb.URIIndication(finalIndication)
	}
	finalSubIndication := b.detailedReport.FinalSubIndication(tokenId)
	if finalSubIndication != "" {
		validationStatus.SubIndication = append(validationStatus.SubIndication, jaxb.URISubIndication(finalSubIndication))
	}
}

// fillMessages is the port of the private
// fillMessages(ValidationStatusType, String).
func (b *ETSIValidationReportBuilder) fillMessages(validationStatus *jaxb.ValidationStatusType, tokenId string) {
	b.fillMessagesOfType(validationStatus, messageValues(b.detailedReport.AdESValidationErrors(tokenId)), enumerations.MessageType_ERROR)
	b.fillMessagesOfType(validationStatus, messageValues(b.detailedReport.AdESValidationWarnings(tokenId)), enumerations.MessageType_WARN)
	b.fillMessagesOfType(validationStatus, messageValues(b.detailedReport.AdESValidationInfos(tokenId)), enumerations.MessageType_INFO)
}

// conclusionValidationStatus is the port of the private
// getValidationStatus(XmlConclusion).
func (b *ETSIValidationReportBuilder) conclusionValidationStatus(conclusion *drjaxb.XmlConclusion) *jaxb.ValidationStatusType {
	validationStatus := &jaxb.ValidationStatusType{}
	b.fillIndicationSubIndicationFromConclusion(validationStatus, conclusion)
	b.fillMessagesFromConclusion(validationStatus, conclusion)
	return validationStatus
}

// fillIndicationSubIndicationFromConclusion is the port of the private
// fillIndicationSubIndication(ValidationStatusType, XmlConclusion).
func (b *ETSIValidationReportBuilder) fillIndicationSubIndicationFromConclusion(validationStatus *jaxb.ValidationStatusType,
	conclusion *drjaxb.XmlConclusion) {
	if conclusion.Indication != "" {
		validationStatus.MainIndication = jaxb.URIIndication(conclusion.Indication.Indication())
	}
	if conclusion.SubIndication != nil {
		validationStatus.SubIndication = append(validationStatus.SubIndication,
			jaxb.URISubIndication(conclusion.SubIndication.SubIndication()))
	}
}

// fillMessagesFromConclusion is the port of the private
// fillMessages(ValidationStatusType, XmlConclusion).
func (b *ETSIValidationReportBuilder) fillMessagesFromConclusion(validationStatus *jaxb.ValidationStatusType,
	conclusion *drjaxb.XmlConclusion) {
	b.fillMessagesOfType(validationStatus, xmlMessageValues(conclusion.Errors), enumerations.MessageType_ERROR)
	b.fillMessagesOfType(validationStatus, xmlMessageValues(conclusion.Warnings), enumerations.MessageType_WARN)
	b.fillMessagesOfType(validationStatus, xmlMessageValues(conclusion.Infos), enumerations.MessageType_INFO)
}

// fillMessagesOfType is the port of the private
// fillMessagesOfType(ValidationStatusType, List<String>, MessageType).
func (b *ETSIValidationReportBuilder) fillMessagesOfType(validationStatus *jaxb.ValidationStatusType, messages []string,
	level enumerations.MessageType) {
	if len(messages) > 0 {
		validationReportData := b.associatedValidationReportData(validationStatus)
		additionalValidationReportData := b.additionalValidationReportData(validationReportData)
		for _, message := range messages {
			reportData := &jaxb.TypedDataType{}
			reportData.Type = level.URI()
			reportData.Value = xsStringValue(message)
			additionalValidationReportData.ReportData = append(additionalValidationReportData.ReportData, reportData)
		}
	}
}

// associatedValidationReportData is the port of the private
// getAssociatedValidationReportData(ValidationStatusType).
func (b *ETSIValidationReportBuilder) associatedValidationReportData(
	validationStatus *jaxb.ValidationStatusType) *jaxb.ValidationReportDataType {
	if len(validationStatus.AssociatedValidationReportData) > 0 {
		// only one is used
		return validationStatus.AssociatedValidationReportData[0]
	}
	validationReportData := &jaxb.ValidationReportDataType{}
	validationStatus.AssociatedValidationReportData = append(validationStatus.AssociatedValidationReportData, validationReportData)
	return validationReportData
}

// additionalValidationReportData is the port of the private
// getAdditionalValidationReportData(ValidationReportDataType).
func (b *ETSIValidationReportBuilder) additionalValidationReportData(
	validationReportData *jaxb.ValidationReportDataType) *jaxb.AdditionalValidationReportDataType {
	additionalValidationReportData := validationReportData.AdditionalValidationReportData
	if additionalValidationReportData == nil {
		additionalValidationReportData = &jaxb.AdditionalValidationReportDataType{}
		validationReportData.AdditionalValidationReportData = additionalValidationReportData
	}
	return additionalValidationReportData
}

// addTokenValidationReportData is the port of the private
// addValidationReportData(ValidationStatusType, AbstractTokenProxy).
func (b *ETSIValidationReportBuilder) addTokenValidationReportData(validationStatus *jaxb.ValidationStatusType,
	token diagnostic.TokenProxy) {
	basicBuildingBlock := b.detailedReport.BasicBuildingBlockById(token.Id())
	signingCertificate := b.detailedReport.SigningCertificate(token.Id())

	if basicBuildingBlock != nil || signingCertificate != nil {
		validationReportData := b.associatedValidationReportData(validationStatus)
		if basicBuildingBlock != nil {
			certificateChain := basicBuildingBlock.CertificateChain
			if certificateChain != nil {
				b.fillCertificateChainAndTrustAnchor(validationReportData, certificateChain)
			}
			aov := basicBuildingBlock.AOV
			cryptographicValidation := process.GetFinalCryptographicValidation(aov)
			if cryptographicValidation != nil {
				b.fillTokenCryptographicInfo(validationReportData, token, cryptographicValidation)
			}
		}
		if signingCertificate != nil && signingCertificate.RevocationInfo != nil {
			b.fillRevocationInfo(validationReportData, signingCertificate.RevocationInfo)
		}
	}
}

// fillTokenCryptographicInfo is the port of the private
// fillCryptographicInfo(ValidationReportDataType, AbstractTokenProxy,
// XmlCryptographicValidation).
func (b *ETSIValidationReportBuilder) fillTokenCryptographicInfo(validationReportData *jaxb.ValidationReportDataType,
	token diagnostic.TokenProxy, cryptographicValidation *drjaxb.XmlCryptographicValidation) {
	cryptoInformationType := &jaxb.CryptoInformationType{}
	switch typed := token.(type) {
	case *diagnostic.SignatureWrapper:
		cryptoInformationType.ValidationObjectId = *b.voReferenceOfSignature(b.signatureIdentifier(typed))
	case *diagnostic.TimestampWrapper:
		cryptoInformationType.ValidationObjectId = *b.voReferenceOfObject(b.timestampValidationObject(typed))
	case *diagnostic.RevocationWrapper:
		cryptoInformationType.ValidationObjectId = *b.voReferenceOfObject(b.revocationValidationObject(typed))
	case *diagnostic.EAAWrapper:
		cryptoInformationType.ValidationObjectId = *b.voReferenceOfObject(b.eaaValidationObject(typed))
	default:
		panic(fmt.Sprintf("Unsupported class %T", token))
	}
	cryptoInformationType.SecureAlgorithm = enumerations.Indication_PASSED == cryptographicValidation.Conclusion.Indication.Indication()
	algorithm := cryptographicValidation.Algorithm
	if algorithm != nil {
		cryptoInformationType.Algorithm = algorithm.Uri
	}
	cryptoInformationType.NotAfter = notAfterOf(cryptographicValidation)
	validationReportData.CryptoInformation = cryptoInformationType
}

// fillRevocationInfo is the port of the private
// fillRevocationInfo(ValidationReportDataType, XmlRevocationInformation).
func (b *ETSIValidationReportBuilder) fillRevocationInfo(validationReportData *jaxb.ValidationReportDataType,
	revocationInfo *drjaxb.XmlRevocationInformation) {
	revocationStatusInformationType := &jaxb.RevocationStatusInformationType{}
	revocationStatusInformationType.RevocationTime = jaxb.XSDateTime(time.Time(revocationInfo.RevocationDate))
	revocationWrapper := b.diagnosticData.RevocationById(revocationInfo.RevocationId)
	revocationStatusInformationType.RevocationObject = b.voReferenceOfObject(b.revocationValidationObject(revocationWrapper))
	certificateWrapper := b.diagnosticData.CertificateById(revocationInfo.CertificateId)
	revocationStatusInformationType.ValidationObjectId = *b.voReferenceOfObject(b.certificateValidationObject(certificateWrapper))
	if revocationInfo.Reason != nil {
		reason := jaxb.URIRevocationReason(enumerations.RevocationReason(*revocationInfo.Reason))
		revocationStatusInformationType.RevocationReason = &reason
	}
	validationReportData.RevocationStatusInformation = revocationStatusInformationType
}

// fillCertificateChainAndTrustAnchor is the port of the private
// fillCertificateChainAndTrustAnchor(ValidationReportDataType,
// XmlCertificateChain).
func (b *ETSIValidationReportBuilder) fillCertificateChainAndTrustAnchor(validationReportData *jaxb.ValidationReportDataType,
	certificateChain *drjaxb.XmlCertificateChain) {
	chainItem := certificateChain.ChainItem
	if len(chainItem) == 0 {
		return
	}

	certificateChainType := &jaxb.CertificateChainType{}
	var signingCert *jaxb.VOReferenceType
	var trustAnchor *jaxb.VOReferenceType
	for i := 0; i < len(chainItem); i++ {
		currentChainItem := chainItem[i]
		certificateWrapper := b.diagnosticData.CertificateById(currentChainItem.Id)
		currentVORef := b.voReferenceOfObject(b.certificateValidationObject(certificateWrapper))

		isSigningCert := i == 0
		isTrustAnchor := certificateWrapper.IsTrusted()

		if isSigningCert || isTrustAnchor {
			if isSigningCert {
				signingCert = currentVORef
			}
			if isTrustAnchor {
				trustAnchor = currentVORef
				// Stops with the first found trust anchor
				break
			}
		} else {
			certificateChainType.IntermediateCertificate = append(certificateChainType.IntermediateCertificate, currentVORef)
		}
	}

	if signingCert != nil {
		certificateChainType.SigningCertificate = *signingCert
	}
	certificateChainType.TrustAnchor = trustAnchor

	validationReportData.CertificateChain = certificateChainType
	validationReportData.TrustAnchor = trustAnchor
}

// signatureIdentifier is the port of the private
// getSignatureIdentifier(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signatureIdentifier(sigWrapper *diagnostic.SignatureWrapper) *jaxb.SignatureIdentifierType {
	signatureIdentifier := b.signatureIdentifierMap[sigWrapper.Id()]
	if signatureIdentifier == nil {
		signatureIdentifier = &jaxb.SignatureIdentifierType{}
		b.signatureIdentifierMap[sigWrapper.Id()] = signatureIdentifier

		signatureIdentifier.Id = sigWrapper.Id()
		if daIdentifier := sigWrapper.DAIdentifier(); daIdentifier != "" {
			signatureIdentifier.DAIdentifier = &daIdentifier
		}
		signatureIdentifier.DocHashOnly = sigWrapper.IsDocHashOnly()
		signatureIdentifier.HashOnly = sigWrapper.IsHashOnly()
		signatureIdentifier.DigestAlgAndValue = b.dtbsrDigestAlgAndValue(sigWrapper)
		sigValue := &jaxb.SignatureValueType{}
		sigValue.Value = jaxb.Base64Binary(sigWrapper.SignatureValue())
		signatureIdentifier.SignatureValue = sigValue
	}
	return signatureIdentifier
}

// dtbsrDigestAlgAndValue is the port of the private
// getDTBSRDigestAlgAndValue(SignatureWrapper).
func (b *ETSIValidationReportBuilder) dtbsrDigestAlgAndValue(sigWrapper *diagnostic.SignatureWrapper) *jaxb.DigestAlgAndValueType {
	dtbsr := sigWrapper.DataToBeSignedRepresentation()
	if dtbsr != nil {
		return b.digestAlgAndValueType(sigWrapper.DataToBeSignedRepresentation())
	}
	return nil
}

// signersDocument is the port of the private
// getSignersDocument(SignatureWrapper).
func (b *ETSIValidationReportBuilder) signersDocument(sigWrapper *diagnostic.SignatureWrapper) *jaxb.SignersDocumentType {
	signerDocuments := b.diagnosticData.SignerDocuments(sigWrapper.Id())
	if len(signerDocuments) == 0 {
		return nil
	}

	signersDocumentType := &jaxb.SignersDocumentType{}
	if len(signerDocuments) == 1 {
		signerDocument := signerDocuments[0]
		digestAlgAndValueType := b.digestAlgAndValueType(signerDocument.DigestAlgoAndValue())
		signersDocumentType.Items = append(signersDocumentType.Items,
			jaxb.NewChoiceItem("DigestAlgAndValue", digestAlgAndValueType))
	}

	var validationObjectList []*jaxb.ValidationObjectType

	signatureScopes := sigWrapper.SignatureScopes()
	signerDataList := make([]*diagnostic.SignerDataWrapper, 0, len(signatureScopes))
	for _, s := range signatureScopes {
		signerDataList = append(signerDataList, diagnostic.NewSignerDataWrapper(s.SignerData))
	}
	for _, signerDataWrapper := range signerDataList {
		validationObjectList = append(validationObjectList, b.signerDataValidationObject(signerDataWrapper))
	}

	signersDocumentType.Items = append(signersDocumentType.Items,
		jaxb.NewChoiceItem("SignersDocumentRepresentation", b.voReference(validationObjectList)))

	return signersDocumentType
}

// signatureAttributes is the port of the private
// getSignatureAttributes(SignatureWrapper). The element names below are the
// ones the corresponding ObjectFactory.createSignatureAttributesTypeXxx
// factory methods produce.
func (b *ETSIValidationReportBuilder) signatureAttributes(sigWrapper *diagnostic.SignatureWrapper) *jaxb.SignatureAttributesType {
	sigAttributes := &jaxb.SignatureAttributesType{}
	// <element name="SigningTime" type="SASigningTimeType"/>
	b.addSigningTime(sigAttributes, sigWrapper)
	// <element name="SigningCertificate" type="SACertIDListType"/>
	b.addSigningCertificate(sigAttributes, sigWrapper)
	// <element name="DataObjectFormat" type="SADataObjectFormatType"/>
	b.addDataObjectFormat(sigAttributes, sigWrapper)
	// <element name="CommitmentTypeIndication" type="SACommitmentTypeIndicationType"/>
	b.addCommitmentTypeIndications(sigAttributes, sigWrapper)
	// <element name="AllDataObjectsTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_ALL_DATA_OBJECTS_TIMESTAMP)
	// see TS 119 102-2 - V1.2.1 A.6.3 CAdES
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_CONTENT_TIMESTAMP)
	// <element name="IndividualDataObjectsTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP)
	// <element name="SigPolicyIdentifier" type="SASigPolicyIdentifierType"/>
	b.addSigPolicyIdentifier(sigAttributes, sigWrapper)
	// <element name="SignatureProductionPlace" type="SASignatureProductionPlaceType"/>
	b.addProductionPlace(sigAttributes, sigWrapper)
	// <element name="SignerRole" type="SASignerRoleType"/>
	b.addSignerRoles(sigAttributes, sigWrapper)
	// <element name="CounterSignature" type="SACounterSignatureType"/>
	b.addCounterSignatures(sigAttributes, sigWrapper)
	// <element name="SignatureTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	// <element name="CompleteCertificateRefs" type="SACertIDListType"/>
	b.addCompleteCertificateRefs(sigAttributes, sigWrapper.FoundCertificates())
	// <element name="CompleteRevocationRefs" type="SARevIDListType"/>
	b.addCompleteRevocationRefs(sigAttributes, sigWrapper.FoundRevocations())
	// <element name="AttributeCertificateRefs" type="SACertIDListType"/>
	b.addAttributeCertificateRefs(sigAttributes, sigWrapper.FoundCertificates())
	// <element name="AttributeRevocationRefs" type="SARevIDListType"/>
	b.addAttributeRevocationRefs(sigAttributes, sigWrapper.FoundRevocations())
	// <element name="SigAndRefsTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_VALIDATION_DATA_TIMESTAMP)
	// <element name="RefsOnlyTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP)
	// <element name="CertificateValues" type="AttributeBaseType"/>
	b.addCertificateValues(sigAttributes, sigWrapper.FoundCertificates())
	// <element name="RevocationValues" type="AttributeBaseType"/>
	b.addRevocationValues(sigAttributes, sigWrapper.FoundRevocations())
	// <element name="AttrAuthoritiesCertValues" type="AttributeBaseType"/>
	b.addAttrAuthoritiesCertValues(sigAttributes, sigWrapper.FoundCertificates())
	// <element name="AttributeRevocationValues" type="AttributeBaseType"/>
	b.addAttributeRevocationValues(sigAttributes, sigWrapper.FoundRevocations())
	// <element name="TimeStampValidationData" type="AttributeBaseType"/>
	b.addTimeStampValidationData(sigAttributes, sigWrapper.FoundCertificates(), sigWrapper.FoundRevocations())
	// <element name="ArchiveTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
	// <element name="RenewedDigests" type="SAListOfIntegersType"/>
	// <element name="MessageDigest" type="SAMessageDigestType"/>
	b.addMessageDigest(sigAttributes, sigWrapper)
	// <element name="DSS" type="SADSSType"/>
	b.addDSS(sigAttributes, sigWrapper)
	// <element name="VRI" type="SAVRIType"/>
	b.addVRI(sigAttributes, sigWrapper)
	// <element name="DocTimeStamp" type="SATimestampType"/>
	b.addTimestampsByType(sigAttributes, sigWrapper, enumerations.TimestampType_DOCUMENT_TIMESTAMP)
	// <element name="Reason" type="SAReasonType"/>
	b.addReason(sigAttributes, sigWrapper)
	// <element name="Name" type="SANameType"/>
	b.addSignerName(sigAttributes, sigWrapper)
	// <element name="ContactInfo" type="SAContactInfoType"/>
	b.addContactInfo(sigAttributes, sigWrapper)
	// <element name="SubFilter" type="SASubFilterType"/>
	b.addSubFilter(sigAttributes, sigWrapper)
	// <element name="ByteRange" type="SAListOfIntegersType"/>
	b.addSignatureByteRange(sigAttributes, sigWrapper)
	// <element name="Filter" type="SAFilterType"/>
	b.addFilter(sigAttributes, sigWrapper)
	return sigAttributes
}

// addAttrAuthoritiesCertValues is the port of the private
// addAttrAuthoritiesCertValues(SignatureAttributesType, FoundCertificatesProxy).
func (b *ETSIValidationReportBuilder) addAttrAuthoritiesCertValues(sigAttributes *jaxb.SignatureAttributesType,
	foundCertificates *diagnostic.FoundCertificatesProxy) {
	var validationObjectTypes []*jaxb.ValidationObjectType
	for _, certificateWrapper := range foundCertificates.RelatedCertificatesByOrigin(enumerations.CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES) {
		validationObjectTypes = append(validationObjectTypes, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range foundCertificates.OrphanCertificatesByOrigin(enumerations.CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES) {
		validationObjectTypes = append(validationObjectTypes, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}
	if len(validationObjectTypes) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("AttrAuthoritiesCertValues", b.buildAttributeObjectList(validationObjectTypes)))
	}
}

// addTimeStampValidationData is the port of the private
// addTimeStampValidationData(SignatureAttributesType, FoundCertificatesProxy,
// FoundRevocationsProxy).
func (b *ETSIValidationReportBuilder) addTimeStampValidationData(sigAttributes *jaxb.SignatureAttributesType,
	foundCertificates *diagnostic.FoundCertificatesProxy, foundRevocations *diagnostic.FoundRevocationsProxy) {
	var validationObjectTypes []*jaxb.ValidationObjectType
	for _, certificateWrapper := range foundCertificates.RelatedCertificatesByOrigin(enumerations.CertificateOrigin_TIMESTAMP_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range foundCertificates.OrphanCertificatesByOrigin(enumerations.CertificateOrigin_TIMESTAMP_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}
	for _, revocationWrapper := range foundRevocations.RelatedRevocationsByOrigin(enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range foundRevocations.OrphanRevocationsByOrigin(enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}

	if len(validationObjectTypes) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("TimeStampValidationData", b.buildAttributeObjectList(validationObjectTypes)))
	}
}

// addCertificateValues is the port of the private
// addCertificateValues(SignatureAttributesType, FoundCertificatesProxy).
func (b *ETSIValidationReportBuilder) addCertificateValues(sigAttributes *jaxb.SignatureAttributesType,
	foundCertificates *diagnostic.FoundCertificatesProxy) {
	var validationObjectTypes []*jaxb.ValidationObjectType
	for _, certificateWrapper := range foundCertificates.RelatedCertificatesByOrigin(enumerations.CertificateOrigin_CERTIFICATE_VALUES) {
		validationObjectTypes = append(validationObjectTypes, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range foundCertificates.OrphanCertificatesByOrigin(enumerations.CertificateOrigin_CERTIFICATE_VALUES) {
		validationObjectTypes = append(validationObjectTypes, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}
	// TODO (upstream) : temporary handling for AnyValidationData -> embed in CertificateValues
	for _, certificateWrapper := range foundCertificates.RelatedCertificatesByOrigin(enumerations.CertificateOrigin_ANY_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range foundCertificates.OrphanCertificatesByOrigin(enumerations.CertificateOrigin_ANY_VALIDATION_DATA) {
		validationObjectTypes = append(validationObjectTypes, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}

	if len(validationObjectTypes) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("CertificateValues", b.buildAttributeObjectList(validationObjectTypes)))
	}
}

// addAttributeCertificateRefs is the port of the private
// addAttributeCertificateRefs(SignatureAttributesType, FoundCertificatesProxy).
func (b *ETSIValidationReportBuilder) addAttributeCertificateRefs(sigAttributes *jaxb.SignatureAttributesType,
	foundCertificates *diagnostic.FoundCertificatesProxy) {
	relatedCerts := foundCertificates.RelatedCertificatesByRefOrigin(enumerations.CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS)
	orphanCerts := foundCertificates.OrphanCertificatesByRefOrigin(enumerations.CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS)

	if len(relatedCerts) > 0 || len(orphanCerts) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("AttributeCertificateRefs", b.buildCertIDListType(relatedCerts, orphanCerts)))
	}
}

// addCompleteCertificateRefs is the port of the private
// addCompleteCertificateRefs(SignatureAttributesType, FoundCertificatesProxy).
func (b *ETSIValidationReportBuilder) addCompleteCertificateRefs(sigAttributes *jaxb.SignatureAttributesType,
	foundCertificates *diagnostic.FoundCertificatesProxy) {
	relatedCerts := foundCertificates.RelatedCertificatesByRefOrigin(enumerations.CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS)
	orphanCerts := foundCertificates.OrphanCertificatesByRefOrigin(enumerations.CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS)

	if len(relatedCerts) > 0 || len(orphanCerts) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("CompleteCertificateRefs", b.buildCertIDListType(relatedCerts, orphanCerts)))
	}
}

// addSigningCertificate is the port of the private
// addSigningCertificate(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSigningCertificate(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	foundCertificates := sigWrapper.FoundCertificates()
	relatedCerts := foundCertificates.RelatedCertificatesByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	orphanCerts := foundCertificates.OrphanCertificatesByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)

	if len(relatedCerts) > 0 || len(orphanCerts) > 0 {
		signingCertAttribute := b.buildCertIDListType(relatedCerts, orphanCerts)
		if sigWrapper.IsBLevelTechnicallyValid() {
			signed := true
			signingCertAttribute.Signed = &signed
		}
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("SigningCertificate", signingCertAttribute))
	}
}

// buildAttributeObjectList is the port of the private
// buildAttributeObjectList(List<ValidationObjectType>).
func (b *ETSIValidationReportBuilder) buildAttributeObjectList(validationObjects []*jaxb.ValidationObjectType) *jaxb.AttributeBaseType {
	attributeBaseType := &jaxb.AttributeBaseType{}
	for _, validationObjectType := range validationObjects {
		attributeBaseType.AttributeObject = append(attributeBaseType.AttributeObject,
			b.voReference([]*jaxb.ValidationObjectType{validationObjectType}))
	}
	return attributeBaseType
}

// buildCertIDListType is the port of the private
// buildCertIDListType(List<RelatedCertificateWrapper>, List<OrphanCertificateWrapper>).
func (b *ETSIValidationReportBuilder) buildCertIDListType(relatedCerts []*diagnostic.RelatedCertificateWrapper,
	orphanCerts []*diagnostic.OrphanCertificateWrapper) *jaxb.SACertIDListType {
	certIdList := &jaxb.SACertIDListType{}
	var validationObjects []*jaxb.ValidationObjectType

	if len(relatedCerts) > 0 {
		for _, cert := range relatedCerts {
			validationObjects = append(validationObjects, b.certificateValidationObject(&cert.CertificateWrapper))
		}
	}

	// The getCertID is not instantiated, because all tokens are listed in ValidationObjects
	//
	// See TS 119 102-2 ch. A.3.2 XML (SigningCertificate): for every certificate
	// referenced within the reported attribute that is not present as validation
	// object (for instance because the creator of the validation report cannot
	// gain access to it), this component shall have one CertID child.

	if len(orphanCerts) > 0 {
		allOrphanCertificates := b.diagnosticData.AllOrphanCertificateObjects()
		for _, orphanCert := range orphanCerts {
			if orphanCert != nil {
				if len(orphanCert.References()) > 0 && !containsOrphanCertificateId(allOrphanCertificates, orphanCert.Id()) {
					for _, certRef := range orphanCert.References() {
						certIdList.CertID = append(certIdList.CertID,
							b.buildCertIDType(certRef.DigestAlgoAndValue(), certRef.IssuerSerial()))
					}
				} else {
					validationObjects = append(validationObjects,
						b.orphanCertificateValidationObject(&orphanCert.OrphanCertificateTokenWrapper))
				}
			}

		}
	}

	if len(validationObjects) > 0 {
		certIdList.AttributeObject = append(certIdList.AttributeObject, b.voReference(validationObjects))
	}

	return certIdList
}

// buildCertIDType is the port of the private
// buildCertIDType(XmlDigestAlgoAndValue, byte[]).
func (b *ETSIValidationReportBuilder) buildCertIDType(digestAlgoAndValue *diagnosticjaxb.XmlDigestAlgoAndValue,
	issuerSerial []byte) *jaxb.SACertIDType {
	certIDType := &jaxb.SACertIDType{}
	certIDType.DigestMethod = *b.digestMethodType(digestAlgorithmOf(digestAlgoAndValue))
	certIDType.DigestValue = jaxb.DSDigestValue(digestValueOf(digestAlgoAndValue))
	if issuerSerial != nil {
		certIDType.X509IssuerSerial = jaxb.Base64Binary(issuerSerial)
	}
	return certIDType
}

// addCompleteRevocationRefs is the port of the private
// addCompleteRevocationRefs(SignatureAttributesType, FoundRevocationsProxy).
func (b *ETSIValidationReportBuilder) addCompleteRevocationRefs(sigAttributes *jaxb.SignatureAttributesType,
	revocationsProxy *diagnostic.FoundRevocationsProxy) {
	relatedRevs := revocationsProxy.RelatedRevocationsByRefOrigin(enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)
	orphanRevs := revocationsProxy.OrphanRevocationsByRefOrigin(enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)

	if len(relatedRevs) > 0 || len(orphanRevs) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("CompleteRevocationRefs", b.buildRevIDListType(relatedRevs, orphanRevs)))
	}
}

// addAttributeRevocationRefs is the port of the private
// addAttributeRevocationRefs(SignatureAttributesType, FoundRevocationsProxy).
func (b *ETSIValidationReportBuilder) addAttributeRevocationRefs(sigAttributes *jaxb.SignatureAttributesType,
	revocationsProxy *diagnostic.FoundRevocationsProxy) {
	relatedRevs := revocationsProxy.RelatedRevocationsByRefOrigin(enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)
	orphanRevs := revocationsProxy.OrphanRevocationsByRefOrigin(enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)

	if len(relatedRevs) > 0 || len(orphanRevs) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("AttributeRevocationRefs", b.buildRevIDListType(relatedRevs, orphanRevs)))
	}
}

// buildRevIDListType is the port of the private
// buildRevIDListType(List<RelatedRevocationWrapper>, List<OrphanRevocationWrapper>).
func (b *ETSIValidationReportBuilder) buildRevIDListType(relatedRevs []*diagnostic.RelatedRevocationWrapper,
	orphanRevs []*diagnostic.OrphanRevocationWrapper) *jaxb.SARevIDListType {
	revIDListType := &jaxb.SARevIDListType{}
	var validationObjects []*jaxb.ValidationObjectType

	if len(relatedRevs) > 0 {
		for _, revocation := range relatedRevs {
			validationObjects = append(validationObjects, b.revocationValidationObject(&revocation.RevocationWrapper))
		}
	}

	if len(orphanRevs) > 0 {
		allOrphanRevocations := b.diagnosticData.AllOrphanRevocationObjects()
		for _, orphanRev := range orphanRevs {
			if orphanRev != nil {
				if len(orphanRev.References()) > 0 && !containsOrphanRevocationId(allOrphanRevocations, orphanRev.Id()) {
					for _, revRef := range orphanRev.References() {
						var revID *jaxb.SARevIDItem
						if enumerations.RevocationType_CRL == orphanRev.RevocationType() {
							revID = &jaxb.SARevIDItem{CRLID: b.buildCRLID(revRef.DigestAlgoAndValue())}
						} else {
							if ocspID := b.buildOCSPID(revRef); ocspID != nil {
								revID = &jaxb.SARevIDItem{OCSPID: ocspID}
							}
						}
						if revID != nil {
							revIDListType.CRLIDOrOCSPID = append(revIDListType.CRLIDOrOCSPID, *revID)
						}
					}
				} else {
					validationObjects = append(validationObjects,
						b.orphanRevocationValidationObject(&orphanRev.OrphanRevocationTokenWrapper))
				}
			}

		}
	}

	if len(validationObjects) > 0 {
		revIDListType.AttributeObject = append(revIDListType.AttributeObject, b.voReference(validationObjects))
	}

	return revIDListType
}

// buildCRLID is the port of the private buildCRLID(XmlDigestAlgoAndValue).
func (b *ETSIValidationReportBuilder) buildCRLID(digestAlgoAndValue *diagnosticjaxb.XmlDigestAlgoAndValue) *jaxb.SACRLIDType {
	sacrlidType := &jaxb.SACRLIDType{}
	sacrlidType.DigestMethod = *b.digestMethodType(digestAlgorithmOf(digestAlgoAndValue))
	sacrlidType.DigestValue = jaxb.DSDigestValue(digestValueOf(digestAlgoAndValue))
	return sacrlidType
}

// buildOCSPID is the port of the private buildOCSPID(RevocationRefWrapper).
func (b *ETSIValidationReportBuilder) buildOCSPID(revRef *diagnostic.RevocationRefWrapper) *jaxb.SAOCSPIDType {
	if revRef.ResponderIdName() != "" || len(revRef.ResponderIdKey()) > 0 {
		saocspidType := &jaxb.SAOCSPIDType{}
		if productionTime := revRef.ProductionTime(); productionTime != nil {
			saocspidType.ProducedAt = jaxb.XSDateTime(*productionTime)
		}
		if revRef.ResponderIdName() != "" {
			responderIDByName := revRef.ResponderIdName()
			saocspidType.ResponderIDByName = &responderIDByName
		} else {
			saocspidType.ResponderIDByKey = jaxb.Base64Binary(revRef.ResponderIdKey())
		}
		return saocspidType
	}
	return nil
}

// digestMethodType is the port of the private
// getDigestMethodType(DigestAlgorithm).
func (b *ETSIValidationReportBuilder) digestMethodType(digestAlgorithm enumerations.DigestAlgorithm) *jaxb.DigestMethodType {
	digestMethodType := &jaxb.DigestMethodType{}
	digestMethodType.Algorithm = b.urn(digestAlgorithm)
	return digestMethodType
}

// addRevocationValues is the port of the private
// addRevocationValues(SignatureAttributesType, FoundRevocationsProxy).
func (b *ETSIValidationReportBuilder) addRevocationValues(sigAttributes *jaxb.SignatureAttributesType,
	foundRevocations *diagnostic.FoundRevocationsProxy) {
	var validationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range foundRevocations.RelatedRevocationsByOrigin(enumerations.RevocationOrigin_REVOCATION_VALUES) {
		validationObjects = append(validationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range foundRevocations.OrphanRevocationsByOrigin(enumerations.RevocationOrigin_REVOCATION_VALUES) {
		validationObjects = append(validationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}
	// TODO (upstream) : temporary handling for AnyValidationData -> embed in RevocationValues
	for _, revocationWrapper := range foundRevocations.RelatedRevocationsByOrigin(enumerations.RevocationOrigin_ANY_VALIDATION_DATA) {
		validationObjects = append(validationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range foundRevocations.OrphanRevocationsByOrigin(enumerations.RevocationOrigin_ANY_VALIDATION_DATA) {
		validationObjects = append(validationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}
	if len(validationObjects) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("RevocationValues", b.buildAttributeObjectList(validationObjects)))
	}
}

// addAttributeRevocationValues is the port of the private
// addAttributeRevocationValues(SignatureAttributesType, FoundRevocationsProxy).
func (b *ETSIValidationReportBuilder) addAttributeRevocationValues(sigAttributes *jaxb.SignatureAttributesType,
	foundRevocations *diagnostic.FoundRevocationsProxy) {
	var validationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range foundRevocations.RelatedRevocationsByOrigin(enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES) {
		validationObjects = append(validationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range foundRevocations.OrphanRevocationsByOrigin(enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES) {
		validationObjects = append(validationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}
	if len(validationObjects) > 0 {
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("AttributeRevocationValues", b.buildAttributeObjectList(validationObjects)))
	}
}

// addMessageDigest is the port of the private
// addMessageDigest(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addMessageDigest(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	messageDigest := sigWrapper.MessageDigest()
	if messageDigest != nil && messageDigest.DigestValue != nil {
		messageDigestType := &jaxb.SAMessageDigestType{}
		messageDigestType.Digest = jaxb.Base64Binary(*messageDigest.DigestValue)
		b.setSignedIfValid(sigWrapper, &messageDigestType.AttributeBaseType)
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("MessageDigest", messageDigestType))
	}
}

// addTimestampsByType is the port of the private
// addTimestampsByType(SignatureAttributesType, SignatureWrapper, TimestampType).
func (b *ETSIValidationReportBuilder) addTimestampsByType(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper, timestampType enumerations.TimestampType) {
	timestampListByType := sigWrapper.TimestampListByType(timestampType)
	isSigned := sigWrapper.IsBLevelTechnicallyValid() && timestampType.IsContentTimestamp()

	for _, timestampWrapper := range timestampListByType {
		timestamp := b.saTimestampType(timestampWrapper)
		name := b.wrap(timestampType)
		if isSigned {
			signed := isSigned
			timestamp.Signed = &signed
		}
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem(name, timestamp))
	}
}

// saTimestampType is the port of the private
// getSATimestampType(TimestampWrapper).
func (b *ETSIValidationReportBuilder) saTimestampType(timestampWrapper *diagnostic.TimestampWrapper) *jaxb.SATimestampType {
	timestamp := &jaxb.SATimestampType{}
	if productionTime := timestampWrapper.ProductionTime(); productionTime != nil {
		timestamp.TimeStampValue = jaxb.XSDateTime(*productionTime)
	}
	timestamp.AttributeObject = append(timestamp.AttributeObject,
		b.voReferenceOfObject(b.timestampValidationObject(timestampWrapper)))
	return timestamp
}

// wrap is the port of the private wrap(TimestampType, SATimestampType); Go has
// no JAXBElement, so the element name the factory method would have produced is
// returned instead.
func (b *ETSIValidationReportBuilder) wrap(timestampType enumerations.TimestampType) string {
	switch timestampType {
	case enumerations.TimestampType_SIGNATURE_TIMESTAMP:
		return "SignatureTimeStamp"
	case enumerations.TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP:
		return "IndividualDataObjectsTimeStamp"
	case enumerations.TimestampType_ALL_DATA_OBJECTS_TIMESTAMP, enumerations.TimestampType_CONTENT_TIMESTAMP:
		return "AllDataObjectsTimeStamp"
	case enumerations.TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP:
		return "RefsOnlyTimeStamp"
	case enumerations.TimestampType_VALIDATION_DATA_TIMESTAMP:
		return "SigAndRefsTimeStamp"
	case enumerations.TimestampType_ARCHIVE_TIMESTAMP:
		return "ArchiveTimeStamp"
	case enumerations.TimestampType_DOCUMENT_TIMESTAMP:
		return "DocTimeStamp"
	default:
		panic(fmt.Sprintf("Unsupported timestamp type %s", timestampType))
	}
}

// addSigPolicyIdentifier is the port of the private
// addSigPolicyIdentifier(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSigPolicyIdentifier(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	policyId := sigWrapper.PolicyId()
	// exclude empty and default values
	if policyId != "" && string(enumerations.SignaturePolicyType_IMPLICIT_POLICY) != policyId {
		saSigPolicyIdentifierType := &jaxb.SASigPolicyIdentifierType{}
		saSigPolicyIdentifierType.SigPolicyId = policyId
		b.setSignedIfValid(sigWrapper, &saSigPolicyIdentifierType.AttributeBaseType)
		sigAttributes.Items = append(sigAttributes.Items,
			jaxb.NewChoiceItem("SigPolicyIdentifier", saSigPolicyIdentifierType))
	}
}

// addDataObjectFormat is the port of the private
// addDataObjectFormat(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addDataObjectFormat(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	contentType := sigWrapper.ContentType()
	mimeType := sigWrapper.MimeType()
	if contentType != "" || mimeType != "" {
		dataObjectFormatType := &jaxb.SADataObjectFormatType{}
		if contentType != "" {
			value := contentType
			dataObjectFormatType.ContentType = &value
		}
		if mimeType != "" {
			value := mimeType
			dataObjectFormatType.MimeType = &value
		}
		b.setSignedIfValid(sigWrapper, &dataObjectFormatType.AttributeBaseType)
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("DataObjectFormat", dataObjectFormatType))
	}
}

// addCommitmentTypeIndications is the port of the private
// addCommitmentTypeIndications(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addCommitmentTypeIndications(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	commitmentTypeIndications := sigWrapper.CommitmentTypeIndications()
	if len(commitmentTypeIndications) > 0 {
		for _, commitmentTypeIndication := range commitmentTypeIndications {
			commitmentType := &jaxb.SACommitmentTypeIndicationType{}
			if commitmentTypeIndication.Identifier != nil {
				commitmentType.CommitmentTypeIdentifier = *commitmentTypeIndication.Identifier
			}
			b.setSignedIfValid(sigWrapper, &commitmentType.AttributeBaseType)
			sigAttributes.Items = append(sigAttributes.Items,
				jaxb.NewChoiceItem("CommitmentTypeIndication", commitmentType))
		}
	}
}

// addSignerRoles is the port of the private
// addSignerRoles(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSignerRoles(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	signerRoles := sigWrapper.SignerRoles()
	if len(signerRoles) > 0 {
		signerRoleType := &jaxb.SASignerRoleType{}
		for _, role := range signerRoles {
			oneSignerRole := &jaxb.SAOneSignerRoleType{}
			if role.Role != nil {
				oneSignerRole.Role = *role.Role
			}
			if role.Category != nil {
				oneSignerRole.EndorsementType = jaxb.ValueEndorsementType(*role.Category)
			}
			signerRoleType.RoleDetails = append(signerRoleType.RoleDetails, oneSignerRole)
		}
		b.setSignedIfValid(sigWrapper, &signerRoleType.AttributeBaseType)
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("SignerRole", signerRoleType))
	}
}

// addCounterSignatures is the port of the private
// addCounterSignatures(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addCounterSignatures(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	counterSignatures := JavaHashSetOrder(b.diagnosticData.AllCounterSignaturesForMasterSignature(sigWrapper))
	for _, counterSignature := range counterSignatures {
		saCounterSignatureType := &jaxb.SACounterSignatureType{}
		saCounterSignatureType.AttributeObject = append(saCounterSignatureType.AttributeObject,
			b.voReferenceOfSignature(b.signatureIdentifier(counterSignature)))
		signatureReference := b.signatureReference(counterSignature)
		saCounterSignatureType.CounterSignature = *signatureReference
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("CounterSignature", saCounterSignatureType))
	}
}

// addProductionPlace is the port of the private
// addProductionPlace(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addProductionPlace(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	sigProductionPlace := &jaxb.SASignatureProductionPlaceType{}

	/*
	 * A.9.5 PAdES
	 *
	 * For PAdES signatures as specified in ETSI EN 319 142-1 [i.3] and ETSI EN
	 * 319 142-2 [i.4], clause 5 this component shall have the contents of the
	 * Location entry in the Signature PDF dictionary.
	 */
	if sigWrapper.PDFRevision() != nil {
		location := sigWrapper.Location()
		if location != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, location)
		} else {
			return
		}
		/*
		 * XAdES, CAdES, JAdES
		 */
	} else if sigWrapper.IsSignatureProductionPlacePresent() {
		postalAddress := sigWrapper.PostalAddress()
		streetAddress := sigWrapper.StreetAddress()
		city := sigWrapper.City()
		stateOrProvince := sigWrapper.StateOrProvince()
		postOfficeBoxNumber := sigWrapper.PostOfficeBoxNumber()
		postalCode := sigWrapper.PostalCode()
		countryName := sigWrapper.CountryName()

		if len(postalAddress) == 0 && streetAddress == "" && city == "" && stateOrProvince == "" &&
			postOfficeBoxNumber == "" && postalCode == "" && countryName == "" {
			return
		}

		if countryName != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, countryName)
		}
		if stateOrProvince != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, stateOrProvince)
		}
		if city != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, city)
		}
		if streetAddress != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, streetAddress)
		}
		if len(postalAddress) > 0 {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, postalAddress...)
		}
		if postalCode != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, postalCode)
		}
		if postOfficeBoxNumber != "" {
			sigProductionPlace.AddressString = append(sigProductionPlace.AddressString, postOfficeBoxNumber)
		}

	} else {
		return

	}

	b.setSignedIfValid(sigWrapper, &sigProductionPlace.AttributeBaseType)
	sigAttributes.Items = append(sigAttributes.Items,
		jaxb.NewChoiceItem("SignatureProductionPlace", sigProductionPlace))
}

// addFilter is the port of the private
// addFilter(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addFilter(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	filter := sigWrapper.Filter()
	if filter != "" {
		filterType := &jaxb.SAFilterType{}
		filterType.Filter = filter
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("Filter", filterType))
	}
}

// addSubFilter is the port of the private
// addSubFilter(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSubFilter(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	subFilter := sigWrapper.SubFilter()
	if subFilter != "" {
		subFilterType := &jaxb.SASubFilterType{}
		subFilterType.SubFilterElement = subFilter
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("SubFilter", subFilterType))
	}
}

// addContactInfo is the port of the private
// addContactInfo(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addContactInfo(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	contactInfo := sigWrapper.ContactInfo()
	if contactInfo != "" {
		contactInfoType := &jaxb.SAContactInfoType{}
		contactInfoType.ContactInfoElement = contactInfo
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("ContactInfo", contactInfoType))
	}
}

// addSignerName is the port of the private
// addSignerName(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSignerName(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	signerName := sigWrapper.SignerName()
	if signerName != "" {
		nameType := &jaxb.SANameType{}
		nameType.NameElement = signerName
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("Name", nameType))
	}
}

// addReason is the port of the private
// addReason(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addReason(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	reason := sigWrapper.Reason()
	if reason != "" {
		reasonType := &jaxb.SAReasonType{}
		reasonType.ReasonElement = reason
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("Reason", reasonType))
	}
}

// addSignatureByteRange is the port of the private
// addSignatureByteRange(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSignatureByteRange(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	signatureByteRange := sigWrapper.SignatureByteRange()
	if len(signatureByteRange) > 0 {
		byteRange := jaxb.SAListOfIntegers(append([]*big.Int{}, signatureByteRange...))
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("ByteRange", &byteRange))
	}
}

// addDSS is the port of the private addDSS(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addDSS(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	var certificateValidationObjects []*jaxb.ValidationObjectType
	for _, certificateWrapper := range sigWrapper.FoundCertificates().RelatedCertificatesByOrigin(enumerations.CertificateOrigin_DSS_DICTIONARY) {
		certificateValidationObjects = append(certificateValidationObjects, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range sigWrapper.FoundCertificates().OrphanCertificatesByOrigin(enumerations.CertificateOrigin_DSS_DICTIONARY) {
		certificateValidationObjects = append(certificateValidationObjects, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}

	var crlValidationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range sigWrapper.FoundRevocations().RelatedRevocationsByTypeAndOrigin(enumerations.RevocationType_CRL, enumerations.RevocationOrigin_DSS_DICTIONARY) {
		crlValidationObjects = append(crlValidationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range sigWrapper.FoundRevocations().OrphanRevocationsByTypeAndOrigin(enumerations.RevocationType_CRL, enumerations.RevocationOrigin_DSS_DICTIONARY) {
		crlValidationObjects = append(crlValidationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}

	var ocspValidationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range sigWrapper.FoundRevocations().RelatedRevocationsByTypeAndOrigin(enumerations.RevocationType_OCSP, enumerations.RevocationOrigin_DSS_DICTIONARY) {
		ocspValidationObjects = append(ocspValidationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range sigWrapper.FoundRevocations().OrphanRevocationsByTypeAndOrigin(enumerations.RevocationType_OCSP, enumerations.RevocationOrigin_DSS_DICTIONARY) {
		ocspValidationObjects = append(ocspValidationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}

	if len(certificateValidationObjects) > 0 || len(crlValidationObjects) > 0 || len(ocspValidationObjects) > 0 {
		dssType := &jaxb.SADSSType{}
		if len(certificateValidationObjects) > 0 {
			dssType.Certs = b.voReference(certificateValidationObjects)
		}
		if len(crlValidationObjects) > 0 {
			dssType.CRLs = b.voReference(crlValidationObjects)
		}
		if len(ocspValidationObjects) > 0 {
			dssType.OCSPs = b.voReference(ocspValidationObjects)
		}
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("DSS", dssType))
	}
}

// addVRI is the port of the private addVRI(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addVRI(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	var certificateValidationObjects []*jaxb.ValidationObjectType
	for _, certificateWrapper := range sigWrapper.FoundCertificates().RelatedCertificatesByOrigin(enumerations.CertificateOrigin_VRI_DICTIONARY) {
		certificateValidationObjects = append(certificateValidationObjects, b.certificateValidationObject(&certificateWrapper.CertificateWrapper))
	}
	for _, orphanCertificate := range sigWrapper.FoundCertificates().OrphanCertificatesByOrigin(enumerations.CertificateOrigin_VRI_DICTIONARY) {
		certificateValidationObjects = append(certificateValidationObjects, b.orphanCertificateValidationObject(&orphanCertificate.OrphanCertificateTokenWrapper))
	}

	var crlValidationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range sigWrapper.FoundRevocations().RelatedRevocationsByTypeAndOrigin(enumerations.RevocationType_CRL, enumerations.RevocationOrigin_VRI_DICTIONARY) {
		crlValidationObjects = append(crlValidationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range sigWrapper.FoundRevocations().OrphanRevocationsByTypeAndOrigin(enumerations.RevocationType_CRL, enumerations.RevocationOrigin_VRI_DICTIONARY) {
		crlValidationObjects = append(crlValidationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}

	var ocspValidationObjects []*jaxb.ValidationObjectType
	for _, revocationWrapper := range sigWrapper.FoundRevocations().RelatedRevocationsByTypeAndOrigin(enumerations.RevocationType_OCSP, enumerations.RevocationOrigin_VRI_DICTIONARY) {
		ocspValidationObjects = append(ocspValidationObjects, b.revocationValidationObject(&revocationWrapper.RevocationWrapper))
	}
	for _, orphanRevocation := range sigWrapper.FoundRevocations().OrphanRevocationsByTypeAndOrigin(enumerations.RevocationType_OCSP, enumerations.RevocationOrigin_VRI_DICTIONARY) {
		ocspValidationObjects = append(ocspValidationObjects, b.orphanRevocationValidationObject(&orphanRevocation.OrphanRevocationTokenWrapper))
	}

	if len(certificateValidationObjects) > 0 || len(crlValidationObjects) > 0 || len(ocspValidationObjects) > 0 {
		vriType := &jaxb.SAVRIType{}
		if len(certificateValidationObjects) > 0 {
			vriType.Certs = b.voReference(certificateValidationObjects)
		}
		if len(crlValidationObjects) > 0 {
			vriType.CRLs = b.voReference(crlValidationObjects)
		}
		if len(ocspValidationObjects) > 0 {
			vriType.OCSPs = b.voReference(ocspValidationObjects)
		}
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("VRI", vriType))
	}
}

// addSigningTime is the port of the private
// addSigningTime(SignatureAttributesType, SignatureWrapper).
func (b *ETSIValidationReportBuilder) addSigningTime(sigAttributes *jaxb.SignatureAttributesType,
	sigWrapper *diagnostic.SignatureWrapper) {
	signingTime := sigWrapper.ClaimedSigningTime()
	if signingTime != nil {
		saSigningTimeType := &jaxb.SASigningTimeType{}
		saSigningTimeType.Time = jaxb.XSDateTime(*signingTime)
		b.setSignedIfValid(sigWrapper, &saSigningTimeType.AttributeBaseType)
		sigAttributes.Items = append(sigAttributes.Items, jaxb.NewChoiceItem("SigningTime", saSigningTimeType))
	}
}

// setSignedIfValid is the port of the private
// setSignedIfValid(SignatureWrapper, AttributeBaseType).
func (b *ETSIValidationReportBuilder) setSignedIfValid(sigWrapper *diagnostic.SignatureWrapper,
	attribute *jaxb.AttributeBaseType) {
	if sigWrapper.IsBLevelTechnicallyValid() {
		signed := true
		attribute.Signed = &signed
	}
}

// ------------------------------------------------------------ small helpers

// base64Item records a byte[] representation member under the "base64"
// element name JAXB binds byte[] to inside
// ValidationObjectRepresentationType's choice.
func base64Item(binaries []byte) jaxb.ChoiceItem {
	value := jaxb.Base64Binary(binaries)
	return jaxb.NewChoiceItem("base64", &value)
}

// digestAlgAndValueItem records a DigestAlgAndValueType representation member.
func digestAlgAndValueItem(value *jaxb.DigestAlgAndValueType) jaxb.ChoiceItem {
	return jaxb.NewChoiceItem("DigestAlgAndValue", value)
}

// uriItem records a String representation member, which JAXB binds to "URI".
func uriItem(value string) jaxb.ChoiceItem {
	v := value
	return jaxb.NewChoiceItem("URI", &v)
}

// digestAlgorithmOf reads the digest algorithm of a diagnostic digest pair,
// standing in for Java's XmlDigestAlgoAndValue.getDigestMethod().
func digestAlgorithmOf(v *diagnosticjaxb.XmlDigestAlgoAndValue) enumerations.DigestAlgorithm {
	if v == nil || v.DigestMethod == nil {
		return ""
	}
	return enumerations.DigestAlgorithm(*v.DigestMethod)
}

// digestValueOf reads the digest value of a diagnostic digest pair.
func digestValueOf(v *diagnosticjaxb.XmlDigestAlgoAndValue) []byte {
	if v == nil || v.DigestValue == nil {
		return nil
	}
	return []byte(*v.DigestValue)
}

// notAfterOf converts XmlCryptographicValidation.getNotAfter() into the
// validation report's date type, preserving the Java null.
func notAfterOf(cryptographicValidation *drjaxb.XmlCryptographicValidation) *jaxb.XSDateTime {
	if cryptographicValidation.NotAfter == nil {
		return nil
	}
	value := jaxb.XSDateTime(time.Time(*cryptographicValidation.NotAfter))
	return &value
}

// messageValues is the port of the ".stream().map(Message::getValue)" the
// fillMessages overloads apply.
func messageValues(messages []detailedreport.Message) []string {
	result := make([]string, 0, len(messages))
	for _, m := range messages {
		result = append(result, m.Value)
	}
	return result
}

// xmlMessageValues is the port of the ".stream().map(XmlMessage::getValue)"
// the conclusion-based fillMessages applies.
func xmlMessageValues(messages []*drjaxb.XmlMessage) []string {
	result := make([]string, 0, len(messages))
	for _, m := range messages {
		result = append(result, m.Value)
	}
	return result
}

// containsOrphanCertificateId ports Java's List.contains over
// OrphanTokenWrapper, whose equals compares getId().
func containsOrphanCertificateId(tokens []*diagnostic.OrphanCertificateTokenWrapper, id string) bool {
	for _, token := range tokens {
		if token.Id() == id {
			return true
		}
	}
	return false
}

// containsOrphanRevocationId ports Java's List.contains over
// OrphanTokenWrapper, whose equals compares getId().
func containsOrphanRevocationId(tokens []*diagnostic.OrphanRevocationTokenWrapper, id string) bool {
	for _, token := range tokens {
		if token.Id() == id {
			return true
		}
	}
	return false
}

// xsStringValue wraps a plain string as the value of an xs:anyType-typed
// element (TypedDataType.Value is declared `Object` in the generated JAXB
// class). JAXB records the runtime type of such a value as an xsi:type
// attribute, together with the two namespace declarations that make it
// resolvable, and writes the string as escaped character data - so the port
// has to place both explicitly, since the ported RawContent carries raw inner
// XML plus the element's own attributes.
func xsStringValue(value string) jaxb.RawContent {
	return jaxb.RawContent{
		Attrs: []xml.Attr{
			{Name: xml.Name{Space: "http://www.w3.org/2001/XMLSchema-instance", Local: "type"}, Value: "xs:string"},
			{Name: xml.Name{Space: "xmlns", Local: "xs"}, Value: "http://www.w3.org/2001/XMLSchema"},
			{Name: xml.Name{Space: "xmlns", Local: "xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
		},
		Content: escapeCharData(value),
	}
}

// escapeCharData escapes a string for the raw-inner-XML position RawContent
// holds, the way the JAXB RI escapes character data: the three markup
// delimiters only - an apostrophe and a quote stay literal in element content.
func escapeCharData(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\r", "&#13;")
	return replacer.Replace(value)
}

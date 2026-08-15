// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/EAAValidationBlock.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: qualification.EAAQualificationBlock and
// qualification.SignatureQualificationBlock are ported by the shared
// QCERT/QTRUST/QSIG porters into package
// github.com/utain/esig/dss/validation/process/qualification, assumed to have
// the shapes
//
//	func NewEAAQualificationBlock(i18nProvider *i18n.I18nProvider, eaa *diagnostic.EAAWrapper,
//	    eaaConclusion *jaxb.XmlConclusion, signatureMap map[string]*jaxb.XmlSignature,
//	    tlAnalysis []*jaxb.XmlTLAnalysis, loteAnalysis []*jaxb.XmlLoTEAnalysis,
//	    currentTime time.Time) *EAAQualificationBlock
//	// with an Execute() *jaxb.XmlValidationEAAQualification method
//
//	func NewSignatureQualificationBlock(i18nProvider *i18n.I18nProvider,
//	    etsi319102validation *jaxb.XmlConstraintsConclusionWithProofOfExistence,
//	    signingCertificate *diagnostic.CertificateWrapper, tlAnalysis []*jaxb.XmlTLAnalysis) *SignatureQualificationBlock
//	// with an Execute() *jaxb.XmlValidationSignatureQualification method
//
// mirroring Java's EAAQualificationBlock(I18nProvider, EAAWrapper, XmlConclusion, Map, List, List, Date)
// and SignatureQualificationBlock(I18nProvider, XmlConstraintsConclusionWithProofOfExistence, CertificateWrapper, List).
package eaa

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process/qualification"
	"github.com/utain/esig/dss/validation/process/vpfbs"
	"github.com/utain/esig/dss/validation/reports"
)

// EAAValidationBlock performs validation of the EAA.
type EAAValidationBlock struct {
	// i18nProvider is the i18n provider.
	i18nProvider *i18n.I18nProvider

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// Policy is the validation policy. Exported because Java declares the
	// field protected.
	Policy policy.ValidationPolicy

	// CurrentTime is the validation time. Exported because Java declares the
	// field protected.
	CurrentTime time.Time

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// tlAnalysis is the list of Trusted List validations.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// loteAnalysis is the list of List of Trusted Entities validations.
	loteAnalysis []*jaxb.XmlLoTEAnalysis
}

// NewEAAValidationBlock is the default constructor. Port of
// EAAValidationBlock(I18nProvider, DiagnosticData, ValidationPolicy, Date, Map, List, List).
func NewEAAValidationBlock(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	tlAnalysis []*jaxb.XmlTLAnalysis, loteAnalysis []*jaxb.XmlLoTEAnalysis) *EAAValidationBlock {
	return &EAAValidationBlock{
		i18nProvider:   i18nProvider,
		diagnosticData: diagnosticData,
		Policy:         validationPolicy,
		CurrentTime:    currentTime,
		bbbs:           bbbs,
		tlAnalysis:     tlAnalysis,
		loteAnalysis:   loteAnalysis,
	}
}

// Execute performs validation of EAA presentations. Port of execute().
func (b *EAAValidationBlock) Execute() []*jaxb.XmlEAA {
	var result []*jaxb.XmlEAA

	for _, eaaWrapper := range b.diagnosticData.EAAs() {
		eaaAnalysis := &jaxb.XmlEAA{}
		id := eaaWrapper.Id()
		eaaAnalysis.Id = &id

		signatureValidationMap := make(map[string]*jaxb.XmlSignature)

		for _, signature := range eaaWrapper.EAASignatures() {
			signatureValidation := b.getEAASignatureValidation(signature)
			eaaAnalysis.Signature = append(eaaAnalysis.Signature, signatureValidation)
			signatureValidationMap[signature.Id()] = signatureValidation
		}

		if eaaWrapper.KeyBindingSignature() != nil {
			signatureValidation := b.getEAASignatureValidation(eaaWrapper.KeyBindingSignature())
			eaaAnalysis.KeyBindingSignature = signatureValidation
			signatureValidationMap[eaaWrapper.KeyBindingSignature().Id()] = signatureValidation
		}

		eaapvp := NewEAAValidationProcess(b.i18nProvider, eaaWrapper, signatureValidationMap, b.bbbs, b.Policy)
		validationProcessEAA := eaapvp.Execute()
		eaaAnalysis.ValidationProcessEAA = validationProcessEAA

		conclusion := validationProcessEAA.Conclusion
		eaaAnalysis.Conclusion = conclusion

		if b.Policy.EIDASConstraintPresent() {

			for _, signature := range eaaWrapper.EAASignatures() {

				xmlSignature := signatureValidationMap[signature.Id()]
				validationSignatureQualification := b.getXmlValidationSignatureQualification(signature, xmlSignature)
				xmlSignature.ValidationSignatureQualification = validationSignatureQualification

			}

			qualificationBlock := qualification.NewEAAQualificationBlock(
				b.i18nProvider, eaaWrapper, conclusion, signatureValidationMap, b.tlAnalysis, b.loteAnalysis, b.CurrentTime)
			eaaAnalysis.ValidationEAAQualification = qualificationBlock.Execute()

		}

		result = append(result, eaaAnalysis)
	}

	return result
}

// getEAASignatureValidation ports the private
// getEAASignatureValidation(SignatureWrapper).
func (b *EAAValidationBlock) getEAASignatureValidation(signatureWrapper *diagnostic.SignatureWrapper) *jaxb.XmlSignature {

	xmlSignature := &jaxb.XmlSignature{}
	id := signatureWrapper.Id()
	xmlSignature.Id = &id

	validation := b.executeBasicValidation(xmlSignature, signatureWrapper)

	conclusion := validation.Conclusion
	finalIndication, err := b.getSignatureFinalIndication(conclusion.Indication.Indication())
	if err != nil {
		panic(err)
	}
	conclusion.Indication = jaxb.IndicationValue(finalIndication)
	xmlSignature.Conclusion = conclusion

	return xmlSignature
}

// executeBasicValidation ports the private
// executeBasicValidation(XmlSignature, SignatureWrapper, Map).
func (b *EAAValidationBlock) executeBasicValidation(signatureAnalysis *jaxb.XmlSignature,
	signature *diagnostic.SignatureWrapper) *jaxb.XmlValidationProcessBasicSignature {
	vpfbsProcess := vpfbs.NewBasicSignatureValidationProcess(b.i18nProvider, b.diagnosticData, signature, nil, b.bbbs)
	bs := vpfbsProcess.Execute()
	signatureAnalysis.ValidationProcessBasicSignature = bs
	return bs
}

// getXmlValidationSignatureQualification ports the private
// getXmlValidationSignatureQualification(SignatureWrapper, XmlSignature).
func (b *EAAValidationBlock) getXmlValidationSignatureQualification(signature *diagnostic.SignatureWrapper,
	xmlSignature *jaxb.XmlSignature) *jaxb.XmlValidationSignatureQualification {
	if xmlSignature == nil {
		panic(fmt.Sprintf("Signature validation is not found for Id '%s'", signature.Id()))
	}

	signatureQualificationBlock := qualification.NewSignatureQualificationBlock(
		b.i18nProvider, &xmlSignature.ValidationProcessBasicSignature.XmlConstraintsConclusionWithProofOfExistenceContent,
		signature.SigningCertificate(), b.tlAnalysis)
	return signatureQualificationBlock.Execute()
}

// getSignatureFinalIndication ports the private getSignatureFinalIndication(Indication).
func (b *EAAValidationBlock) getSignatureFinalIndication(highestIndication enumerations.Indication) (enumerations.Indication, error) {
	switch highestIndication {
	case enumerations.Indication_PASSED:
		return enumerations.Indication_TOTAL_PASSED, nil
	case enumerations.Indication_INDETERMINATE:
		return enumerations.Indication_INDETERMINATE, nil
	case enumerations.Indication_FAILED:
		return enumerations.Indication_TOTAL_FAILED, nil
	default:
		return "", reports.NewDSSReportExceptionMessage(fmt.Sprintf("The Indication '%s' is not supported!", highestIndication))
	}
}

// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReportMessageCollector.java
// (DSS 6.5.RC1).
//
// slf4j logging is dropped per PORTING.md.
package detailedreport

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// messageType mirrors eu.europa.esig.dss.enumerations.MessageType's role
// inside this collector (dispatching to Errors/Warnings/Infos); the actual
// enumerations.MessageType constants are used directly at call sites.
type messageType = enumerations.MessageType

// DetailedReportMessageCollector is used to collect all messages for a token
// validation by a defined type from a DetailedReport.
type DetailedReportMessageCollector struct {
	// detailedReport is the DetailedReport used to collect messages from.
	detailedReport *DetailedReport
}

// newDetailedReportMessageCollector is the default constructor.
func newDetailedReportMessageCollector(detailedReport *DetailedReport) *DetailedReportMessageCollector {
	if detailedReport == nil {
		panic("DetailedReport cannot be nil!")
	}
	return &DetailedReportMessageCollector{detailedReport: detailedReport}
}

// AdESValidationErrors returns a list of ETSI EN 319 102-1 AdES validation
// error messages for a token with the given id.
func (c *DetailedReportMessageCollector) AdESValidationErrors(tokenId string) []Message {
	return c.collectAdESValidationMessages(enumerations.MessageTypeError, tokenId)
}

// AdESValidationWarnings returns a list of ETSI EN 319 102-1 AdES validation
// warning messages for a token with the given id.
func (c *DetailedReportMessageCollector) AdESValidationWarnings(tokenId string) []Message {
	return c.collectAdESValidationMessages(enumerations.MessageTypeWarn, tokenId)
}

// AdESValidationInfos returns a list of ETSI EN 319 102-1 AdES validation info
// messages for a token with the given id.
func (c *DetailedReportMessageCollector) AdESValidationInfos(tokenId string) []Message {
	return c.collectAdESValidationMessages(enumerations.MessageTypeInfo, tokenId)
}

// QualificationErrors returns a list of qualification validation errors for a
// token with the given id.
func (c *DetailedReportMessageCollector) QualificationErrors(tokenId string) []Message {
	return c.collectQualificationMessages(enumerations.MessageTypeError, tokenId)
}

// QualificationWarnings returns a list of qualification validation warnings
// for a token with the given id.
func (c *DetailedReportMessageCollector) QualificationWarnings(tokenId string) []Message {
	return c.collectQualificationMessages(enumerations.MessageTypeWarn, tokenId)
}

// QualificationInfos returns a list of qualification validation infos for a
// token with the given id.
func (c *DetailedReportMessageCollector) QualificationInfos(tokenId string) []Message {
	return c.collectQualificationMessages(enumerations.MessageTypeInfo, tokenId)
}

// CertificateQualificationErrorsAtIssuanceTime returns a list of
// qualification validation errors for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationErrorsAtIssuanceTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtIssuanceTime(enumerations.MessageTypeError, certificateId)
}

// CertificateQualificationWarningsAtIssuanceTime returns a list of
// qualification validation warnings for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationWarningsAtIssuanceTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtIssuanceTime(enumerations.MessageTypeWarn, certificateId)
}

// CertificateQualificationInfosAtIssuanceTime returns a list of qualification
// validation information messages for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationInfosAtIssuanceTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtIssuanceTime(enumerations.MessageTypeInfo, certificateId)
}

// CertificateQualificationErrorsAtValidationTime returns a list of
// qualification validation errors for a certificate with the given id at
// validation time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationErrorsAtValidationTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtValidationTime(enumerations.MessageTypeError, certificateId)
}

// CertificateQualificationWarningsAtValidationTime returns a list of
// qualification validation warnings for a certificate with the given id at
// validation time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationWarningsAtValidationTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtValidationTime(enumerations.MessageTypeWarn, certificateId)
}

// CertificateQualificationInfosAtValidationTime returns a list of
// qualification validation information messages for a certificate with the
// given id at validation time.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateQualificationInfosAtValidationTime(certificateId string) []Message {
	return c.collectCertificateQualificationAtValidationTime(enumerations.MessageTypeInfo, certificateId)
}

// QWACValidationErrors returns a list of QWAC validation errors for a
// certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (c *DetailedReportMessageCollector) QWACValidationErrors(certificateId string) []Message {
	return c.collectQWACValidationDetails(enumerations.MessageTypeError, certificateId)
}

// QWACValidationWarnings returns a list of QWAC validation warnings for a
// certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (c *DetailedReportMessageCollector) QWACValidationWarnings(certificateId string) []Message {
	return c.collectQWACValidationDetails(enumerations.MessageTypeWarn, certificateId)
}

// QWACValidationInfos returns a list of QWAC validation information messages
// for a certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (c *DetailedReportMessageCollector) QWACValidationInfos(certificateId string) []Message {
	return c.collectQWACValidationDetails(enumerations.MessageTypeInfo, certificateId)
}

// CertificateApprovalStatusErrorsAtIssuanceTime returns a list of
// TS 119 602 certificate approval status validation errors for a certificate
// with the given id at certificate issuance time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusErrorsAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtIssuanceTime(enumerations.MessageTypeError, certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusWarningsAtIssuanceTime returns a list of
// TS 119 602 certificate approval status validation warnings for a
// certificate with the given id at certificate issuance time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusWarningsAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtIssuanceTime(enumerations.MessageTypeWarn, certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusInfosAtIssuanceTime returns a list of TS 119 602
// certificate approval status validation information messages for a
// certificate with the given id at certificate issuance time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusInfosAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtIssuanceTime(enumerations.MessageTypeInfo, certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusErrorsAtValidationTime returns a list of
// TS 119 602 certificate approval status validation errors for a certificate
// with the given id at validation time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusErrorsAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtValidationTime(enumerations.MessageTypeError, certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusWarningsAtValidationTime returns a list of
// TS 119 602 certificate approval status validation warnings for a
// certificate with the given id at validation time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusWarningsAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtValidationTime(enumerations.MessageTypeWarn, certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusInfosAtValidationTime returns a list of
// TS 119 602 certificate approval status validation information messages for
// a certificate with the given id at validation time and the given
// certificateApprovalStatus.
// NOTE: applicable only for certificate validation.
func (c *DetailedReportMessageCollector) CertificateApprovalStatusInfosAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtValidationTime(enumerations.MessageTypeInfo, certificateId, certificateApprovalStatus)
}

func (c *DetailedReportMessageCollector) collectAdESValidationMessages(t messageType, tokenId string) []Message {
	if signatureById := c.detailedReport.XmlSignatureById(tokenId); signatureById != nil {
		return c.collectSignatureValidation(t, signatureById)
	}
	if timestampById := c.detailedReport.XmlTimestampById(tokenId); timestampById != nil {
		return c.collectTimestampValidation(t, timestampById)
	}
	if evidenceRecordById := c.detailedReport.XmlEvidenceRecordById(tokenId); evidenceRecordById != nil {
		return c.collectEvidenceRecordValidation(t, evidenceRecordById)
	}
	if eaaById := c.detailedReport.XmlEAAById(tokenId); eaaById != nil {
		return c.collectEAAValidation(t, eaaById)
	}
	if tlAnalysisById := c.detailedReport.TLAnalysisById(tokenId); tlAnalysisById != nil {
		return c.collectTLAnalysisValidation(t, tlAnalysisById)
	}
	if bbbById := c.detailedReport.BasicBuildingBlockById(tokenId); bbbById != nil {
		return c.collectBBBValidation(t, bbbById)
	}
	// supported only for certificate validation
	if c.detailedReport.IsCertificateValidation() {
		if certXCVConclusion := c.detailedReport.CertificateXCVConclusion(tokenId); certXCVConclusion != nil {
			return c.collectXmlConclusionValidation(t, certXCVConclusion)
		}
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) collectQualificationMessages(t messageType, tokenId string) []Message {
	if signatureById := c.detailedReport.XmlSignatureById(tokenId); signatureById != nil {
		return c.collectSignatureQualification(t, signatureById)
	}
	if timestampById := c.detailedReport.XmlTimestampById(tokenId); timestampById != nil {
		return c.collectTimestampQualification(t, timestampById)
	}
	if eaaById := c.detailedReport.XmlEAAById(tokenId); eaaById != nil {
		return c.collectEAAQualification(t, eaaById)
	}
	if certificateById := c.certificateQualificationProcess(tokenId); certificateById != nil {
		return c.collectCertificateQualification(t, certificateById)
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) collectSignatureValidation(t messageType, xmlSignature *jaxb.XmlSignature) []Message {
	result := []Message{}

	highestConclusion := c.detailedReport.HighestConclusion(derefString(xmlSignature.Id))
	if enumerations.MessageTypeError != t || (xmlSignature.ValidationProcessBasicSignature != nil &&
		subIndicationOf(highestConclusion.Conclusion) == subIndicationOf(xmlSignature.ValidationProcessBasicSignature.Conclusion)) {
		if xmlSignature.ValidationProcessBasicSignature != nil {
			addMessages(&result, messages(t, xmlSignature.ValidationProcessBasicSignature.Conclusion))
		}
	}
	if enumerations.MessageTypeError != t || (xmlSignature.ValidationProcessLongTermData != nil &&
		subIndicationOf(highestConclusion.Conclusion) == subIndicationOf(xmlSignature.ValidationProcessLongTermData.Conclusion)) {
		if xmlSignature.ValidationProcessLongTermData != nil {
			addMessages(&result, messages(t, xmlSignature.ValidationProcessLongTermData.Conclusion))
		}
	}
	addMessages(&result, messages(t, highestConclusion.Conclusion))
	return result
}

func (c *DetailedReportMessageCollector) collectTimestampValidation(t messageType, xmlTimestamp *jaxb.XmlTimestamp) []Message {
	result := []Message{}

	timestampBasic := xmlTimestamp.ValidationProcessBasicTimestamp
	timestampArchivalData := xmlTimestamp.ValidationProcessArchivalDataTimestamp
	if timestampArchivalData == nil || enumerations.MessageTypeError != t ||
		enumerations.IndicationPassed != indicationOf(timestampArchivalData.Conclusion) {
		if timestampBasic != nil {
			addMessages(&result, messages(t, timestampBasic.Conclusion))
		}
	}
	if timestampArchivalData != nil {
		addMessages(&result, messages(t, timestampArchivalData.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectEvidenceRecordValidation(t messageType, xmlEvidenceRecord *jaxb.XmlEvidenceRecord) []Message {
	result := []Message{}

	validationProcessEvidenceRecord := xmlEvidenceRecord.ValidationProcessEvidenceRecord
	if validationProcessEvidenceRecord != nil {
		addMessages(&result, messages(t, validationProcessEvidenceRecord.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectEAAValidation(t messageType, xmlEAA *jaxb.XmlEAA) []Message {
	result := []Message{}

	validationProcessEAA := xmlEAA.ValidationProcessEAA
	if validationProcessEAA != nil {
		addMessages(&result, messages(t, validationProcessEAA.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectTLAnalysisValidation(t messageType, tlAnalysisById *jaxb.XmlTLAnalysis) []Message {
	result := []Message{}
	addMessages(&result, messages(t, tlAnalysisById.Conclusion))
	return result
}

func (c *DetailedReportMessageCollector) collectBBBValidation(t messageType, bbbById *jaxb.XmlBasicBuildingBlocks) []Message {
	result := []Message{}
	addMessages(&result, messages(t, bbbById.Conclusion))
	return result
}

func (c *DetailedReportMessageCollector) collectXmlConclusionValidation(t messageType, xmlConclusion *jaxb.XmlConclusion) []Message {
	result := []Message{}
	addMessages(&result, messages(t, xmlConclusion))
	return result
}

func (c *DetailedReportMessageCollector) collectSignatureQualification(t messageType, xmlSignature *jaxb.XmlSignature) []Message {
	result := []Message{}
	if v := xmlSignature.ValidationSignatureQualification; v != nil {
		addMessages(&result, messages(t, v.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectTimestampQualification(t messageType, xmlTimestamp *jaxb.XmlTimestamp) []Message {
	result := []Message{}
	if v := xmlTimestamp.ValidationTimestampQualification; v != nil {
		addMessages(&result, messages(t, v.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectEAAQualification(t messageType, xmlEAA *jaxb.XmlEAA) []Message {
	result := []Message{}
	if v := xmlEAA.ValidationEAAQualification; v != nil {
		addMessages(&result, messages(t, v.Conclusion))
	}
	return result
}

func (c *DetailedReportMessageCollector) collectCertificateQualification(t messageType, certificateQualificationProcess []*jaxb.XmlValidationCertificateQualification) []Message {
	result := []Message{}
	addMessages(&result, c.collectCertificateQualificationAtIssuanceTimeFrom(t, certificateQualificationProcess))
	addMessages(&result, c.collectCertificateQualificationAtBestSignatureTime(t, certificateQualificationProcess))
	addMessages(&result, c.collectCertificateQualificationAtValidationTimeFrom(t, certificateQualificationProcess))
	return result
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtIssuanceTimeFrom(t messageType, certificateQualificationProcess []*jaxb.XmlValidationCertificateQualification) []Message {
	return c.collectCertificateQualificationAtTime(t, certificateQualificationProcess, enumerations.ValidationTimeCertificateIssuanceTime)
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtBestSignatureTime(t messageType, certificateQualificationProcess []*jaxb.XmlValidationCertificateQualification) []Message {
	return c.collectCertificateQualificationAtTime(t, certificateQualificationProcess, enumerations.ValidationTimeBESTSignatureTime)
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtValidationTimeFrom(t messageType, certificateQualificationProcess []*jaxb.XmlValidationCertificateQualification) []Message {
	return c.collectCertificateQualificationAtTime(t, certificateQualificationProcess, enumerations.ValidationTimeValidationTime)
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtIssuanceTime(t messageType, certificateId string) []Message {
	certificateQualificationProcess := c.certificateQualificationProcess(certificateId)
	if certificateQualificationProcess != nil {
		return c.collectCertificateQualificationAtIssuanceTimeFrom(t, certificateQualificationProcess)
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) certificateQualificationProcess(certificateId string) []*jaxb.XmlValidationCertificateQualification {
	xmlCertificate := c.detailedReport.XmlCertificateById(certificateId)
	if xmlCertificate != nil && xmlCertificate.CertificateQualificationProcess != nil {
		return xmlCertificate.CertificateQualificationProcess.ValidationCertificateQualification
	}

	signatures := c.detailedReport.Signatures()
	if len(signatures) != 0 {
		for _, xmlSignature := range signatures {
			signatureQualification := xmlSignature.ValidationSignatureQualification
			if signatureQualification != nil {
				return signatureQualification.ValidationCertificateQualification
			}
		}
	}

	return []*jaxb.XmlValidationCertificateQualification{}
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtValidationTime(t messageType, certificateId string) []Message {
	certificateQualificationProcess := c.certificateQualificationProcess(certificateId)
	if certificateQualificationProcess != nil {
		return c.collectCertificateQualificationAtValidationTimeFrom(t, certificateQualificationProcess)
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) collectCertificateQualificationAtTime(t messageType, certificateQualificationProcess []*jaxb.XmlValidationCertificateQualification, validationTime enumerations.ValidationTime) []Message {
	if certificateQualificationProcess != nil {
		for _, cq := range certificateQualificationProcess {
			if validationTime == validationTimeOf(cq.ValidationTime) {
				return messages(t, cq.Conclusion)
			}
		}
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) collectQWACValidationDetails(t messageType, certificateId string) []Message {
	if certificateId == "" {
		return []Message{}
	}

	var qwacProcess *jaxb.XmlQWACProcess

	xmlCertificate := c.detailedReport.XmlCertificateById(certificateId)
	if xmlCertificate != nil {
		qwacProcess = xmlCertificate.QWACProcess
	} else {
		signatures := c.detailedReport.Signatures()
		if len(signatures) != 0 {
			for _, xmlSignature := range signatures {
				signatureQualification := xmlSignature.ValidationSignatureQualification
				if signatureQualification != nil {
					if signatureQualification.QWACProcess != nil && certificateId == signatureQualification.QWACProcess.Id {
						qwacProcess = signatureQualification.QWACProcess
					}
				}
			}
		}
	}

	if qwacProcess != nil {
		return messages(t, qwacProcess.Conclusion)
	}
	return []Message{}
}

func (c *DetailedReportMessageCollector) collectCertificateApprovalStatusAtIssuanceTime(t messageType, certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	certificateApprovalStatusProcess := c.certificateApprovalStatusProcess(certificateId, certificateApprovalStatus)
	return c.collectCertificateApprovalStatusAtIssuanceTimeFrom(t, certificateApprovalStatusProcess)
}

func (c *DetailedReportMessageCollector) collectCertificateApprovalStatusAtValidationTime(t messageType, certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	certificateApprovalStatusProcess := c.certificateApprovalStatusProcess(certificateId, certificateApprovalStatus)
	return c.collectCertificateApprovalStatusAtValidationTimeFrom(t, certificateApprovalStatusProcess)
}

func (c *DetailedReportMessageCollector) certificateApprovalStatusProcess(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []*jaxb.XmlValidationCertificateApprovalStatus {
	result := []*jaxb.XmlValidationCertificateApprovalStatus{}

	xmlCertificate := c.detailedReport.XmlCertificateById(certificateId)
	if xmlCertificate != nil && xmlCertificate.CertificateApprovalStatusProcess != nil {
		for _, vcas := range xmlCertificate.CertificateApprovalStatusProcess.ValidationCertificateApprovalStatus {
			xmlCertificateApprovalStatus := vcas.CertificateApprovalStatus
			if xmlCertificateApprovalStatus != nil &&
				certificateApprovalStatus.ListType() != nil && certificateApprovalStatus.ListType().URI() != "" &&
				xmlCertificateApprovalStatus.ListType != nil && certificateApprovalStatus.ListType().URI() == xmlCertificateApprovalStatus.ListType.URI() &&
				certificateApprovalStatus.ServiceTypeIdentifier() != nil && certificateApprovalStatus.ServiceTypeIdentifier().URI() != "" &&
				xmlCertificateApprovalStatus.ServiceTypeIdentifier != nil && certificateApprovalStatus.ServiceTypeIdentifier().URI() == xmlCertificateApprovalStatus.ServiceTypeIdentifier.URI() {
				result = append(result, vcas)
			}
		}
	}

	return result
}

func (c *DetailedReportMessageCollector) collectCertificateApprovalStatusAtIssuanceTimeFrom(t messageType, certificateApprovalStatusProcesses []*jaxb.XmlValidationCertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtTime(t, certificateApprovalStatusProcesses, enumerations.ValidationTimeCertificateIssuanceTime)
}

func (c *DetailedReportMessageCollector) collectCertificateApprovalStatusAtValidationTimeFrom(t messageType, certificateApprovalStatusProcesses []*jaxb.XmlValidationCertificateApprovalStatus) []Message {
	return c.collectCertificateApprovalStatusAtTime(t, certificateApprovalStatusProcesses, enumerations.ValidationTimeValidationTime)
}

func (c *DetailedReportMessageCollector) collectCertificateApprovalStatusAtTime(t messageType, certificateApprovalStatusProcesses []*jaxb.XmlValidationCertificateApprovalStatus, validationTime enumerations.ValidationTime) []Message {
	if certificateApprovalStatusProcesses != nil {
		for _, cas := range certificateApprovalStatusProcesses {
			if validationTime == validationTimeOf(cas.ValidationTime) {
				return messages(t, cas.Conclusion)
			}
		}
	}
	return []Message{}
}

func messages(t messageType, conclusion *jaxb.XmlConclusion) []Message {
	if conclusion != nil {
		switch t {
		case enumerations.MessageTypeError:
			return convertMessages(conclusion.Errors)
		case enumerations.MessageTypeWarn:
			return convertMessages(conclusion.Warnings)
		case enumerations.MessageTypeInfo:
			return convertMessages(conclusion.Infos)
		}
	}
	return []Message{}
}

func convertMessage(m *jaxb.XmlMessage) (Message, bool) {
	if m == nil {
		return Message{}, false
	}
	return NewMessage(derefString(m.Key), m.Value), true
}

func convertMessages(messages []*jaxb.XmlMessage) []Message {
	result := make([]Message, 0, len(messages))
	for _, m := range messages {
		if msg, ok := convertMessage(m); ok {
			result = append(result, msg)
		}
	}
	return result
}

func addMessages(result *[]Message, toAdd []Message) {
	for _, m := range toAdd {
		addMessage(result, m)
	}
}

func addMessage(result *[]Message, toAdd Message) {
	for _, existing := range *result {
		if existing == toAdd {
			return
		}
	}
	*result = append(*result, toAdd)
}

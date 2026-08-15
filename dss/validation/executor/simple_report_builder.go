// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/signature/SimpleReportBuilder.java
// (DSS 6.5.RC1).
//
// Deviations, all forced:
//
//   - finalIndications / finalSubIndications are HashSets of enums upstream,
//     so the order of the <Semantic> elements addSemantics() appends is the
//     JVM's identity-hash bucket order - not stable even across upstream runs.
//     This port emits them in enumeration declaration order. Semantics are off
//     by default (includeSemantics=false).
//
//   - getUniqueServiceNames() returns a HashSet<String> whose iteration order
//     reaches the marshalled <TrustServiceName> elements. This port keeps
//     first-seen insertion order of the deduplicated names, matching the
//     precedent set by qualification/granted_status_check.go.
//
//   - Java has three getXmlDisclosableClaim overloads; Go has none, so they
//     become xmlDisclosableClaim / xmlDisclosableClaimWithDisclosure /
//     newXmlDisclosableClaim.
//
//   - DiagnosticData.getValidationDate() is nullable in Java and handed
//     straight to POEExtraction.init(...); Go's POEExtraction.Init takes a
//     time.Time value, so a nil validation date becomes the zero time.

package executor

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/utain/esig/dss/detailedreport"
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/simplereport/jaxb"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
)

// SimpleReportBuilder builds a SimpleReport XmlDom from the diagnostic data
// and detailed validation report. Port of SimpleReportBuilder.
type SimpleReportBuilder struct {
	// i18nProvider is the i18n provider.
	i18nProvider *i18n.I18nProvider

	// includeSemantics defines if the semantics shall be included.
	includeSemantics bool

	// currentTime is the validation time.
	currentTime time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// diagnosticData is the DiagnosticData to use.
	diagnosticData *diagnostic.DiagnosticData

	// detailedReport is the detailed report.
	detailedReport *detailedreport.DetailedReport

	// totalSignatureCount is the number of processed signatures.
	totalSignatureCount int

	// validSignatureCount is the number of valid signatures.
	validSignatureCount int

	// finalIndications is the set of all used Indications (used for
	// semantics); see the file header for the ordering deviation.
	finalIndications map[enumerations.Indication]struct{}

	// finalSubIndications is the set of all used SubIndications (used for
	// semantics).
	finalSubIndications map[enumerations.SubIndication]struct{}

	// poe is the POE set.
	poe *vpfswatsp.POEExtraction
}

// NewSimpleReportBuilder is the default constructor. Port of
// SimpleReportBuilder(I18nProvider, Date, ValidationPolicy, DiagnosticData,
// DetailedReport, boolean).
func NewSimpleReportBuilder(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport, includeSemantics bool) *SimpleReportBuilder {
	return &SimpleReportBuilder{
		currentTime:         currentTime,
		policy:              validationPolicy,
		diagnosticData:      diagnosticData,
		detailedReport:      detailedReport,
		i18nProvider:        i18nProvider,
		includeSemantics:    includeSemantics,
		finalIndications:    make(map[enumerations.Indication]struct{}),
		finalSubIndications: make(map[enumerations.SubIndication]struct{}),
	}
}

// Build generates the validation simpleReport. Port of build().
func (b *SimpleReportBuilder) Build() *jaxb.XmlSimpleReport {
	b.validSignatureCount = 0
	b.totalSignatureCount = 0

	b.poe = vpfswatsp.NewPOEExtraction()
	validationDate := time.Time{}
	if b.diagnosticData.ValidationDate() != nil {
		validationDate = *b.diagnosticData.ValidationDate()
	}
	b.poe.Init(b.diagnosticData, validationDate)
	b.poe.CollectAllPOE(b.diagnosticData.TimestampList())

	simpleReport := &jaxb.XmlSimpleReport{}

	b.addPolicyNode(simpleReport)
	b.addValidationTime(simpleReport)
	b.addDocumentName(simpleReport)

	containerInfoPresent := b.diagnosticData.IsContainerInfoPresent()
	if containerInfoPresent {
		b.addContainerType(simpleReport)
	}

	attachedSignatureIds := make(map[string]struct{})
	attachedTimestampIds := make(map[string]struct{})
	attachedEvidenceRecordIds := make(map[string]struct{})
	if len(b.diagnosticData.EAAs()) > 0 {
		for _, eaa := range b.diagnosticData.EAAs() {
			for _, signatureId := range eaa.EAASignatureIds() {
				attachedSignatureIds[signatureId] = struct{}{}
			}
			if eaa.KeyBindingSignature() != nil {
				attachedSignatureIds[eaa.KeyBindingSignatureId()] = struct{}{}
			}
			simpleReport.SignatureOrTimestampOrEvidenceRecord =
				append(simpleReport.SignatureOrTimestampOrEvidenceRecord, b.eaa(eaa))
		}
	}

	for _, signature := range b.diagnosticData.Signatures() {
		if _, attached := attachedSignatureIds[signature.Id()]; attached {
			continue
		}
		for _, timestampId := range signature.TimestampIdsList() {
			attachedTimestampIds[timestampId] = struct{}{}
		}
		for _, evidenceRecordId := range signature.EvidenceRecordIdsList() {
			attachedEvidenceRecordIds[evidenceRecordId] = struct{}{}
		}
		for _, timestampId := range signature.EvidenceRecordTimestampIds() {
			attachedTimestampIds[timestampId] = struct{}{}
		}
		simpleReport.SignatureOrTimestampOrEvidenceRecord =
			append(simpleReport.SignatureOrTimestampOrEvidenceRecord, b.signature(signature, containerInfoPresent))
	}

	for _, timestamp := range b.diagnosticData.NonEvidenceRecordTimestamps() {
		if _, attached := attachedTimestampIds[timestamp.Id()]; attached {
			continue
		}
		attachedTimestampIds[timestamp.Id()] = struct{}{}
		for _, evidenceRecordId := range timestamp.EvidenceRecordIdsList() {
			attachedEvidenceRecordIds[evidenceRecordId] = struct{}{}
		}
		for _, timestampId := range timestamp.EvidenceRecordTimestampIds() {
			attachedTimestampIds[timestampId] = struct{}{}
		}
		tstValidationIndication := b.detailedReport.BasicTimestampValidationIndication(timestamp.Id())
		if tstValidationIndication != "" {
			xmlTimestamp := b.xmlTimestamp(timestamp)
			if b.isValidConclusion(timestamp.Id()) {
				// perform extension period check only for detached timestamps
				b.determineTimestampExtensionPeriod(xmlTimestamp)
			}
			simpleReport.SignatureOrTimestampOrEvidenceRecord =
				append(simpleReport.SignatureOrTimestampOrEvidenceRecord, xmlTimestamp)
		}
	}

	for _, evidenceRecord := range b.diagnosticData.EvidenceRecords() {
		if _, attached := attachedEvidenceRecordIds[evidenceRecord.Id()]; attached {
			continue
		}
		for _, timestampId := range evidenceRecord.TimestampIdsList() {
			attachedTimestampIds[timestampId] = struct{}{}
		}
		erValidationIndication := b.detailedReport.EvidenceRecordValidationIndication(evidenceRecord.Id())
		if erValidationIndication != "" {
			simpleReport.SignatureOrTimestampOrEvidenceRecord =
				append(simpleReport.SignatureOrTimestampOrEvidenceRecord, b.xmlEvidenceRecord(evidenceRecord))
		}
	}

	for _, timestamp := range b.diagnosticData.TimestampList() {
		if _, attached := attachedTimestampIds[timestamp.Id()]; attached {
			continue
		}
		tstValidationIndication := b.detailedReport.BasicTimestampValidationIndication(timestamp.Id())
		if tstValidationIndication != "" {
			simpleReport.SignatureOrTimestampOrEvidenceRecord =
				append(simpleReport.SignatureOrTimestampOrEvidenceRecord, b.xmlTimestamp(timestamp))
		}
	}

	b.addStatistics(simpleReport)

	if b.includeSemantics {
		b.addSemantics(simpleReport)
	}

	b.addPDFAProfile(simpleReport)

	return simpleReport
}

// addPolicyNode is the port of the private addPolicyNode(XmlSimpleReport).
func (b *SimpleReportBuilder) addPolicyNode(report *jaxb.XmlSimpleReport) {
	xmlPolicy := &jaxb.XmlValidationPolicy{}
	policyName := b.policy.PolicyName()
	xmlPolicy.PolicyName = &policyName
	policyDescription := b.policy.PolicyDescription()
	xmlPolicy.PolicyDescription = &policyDescription
	report.ValidationPolicy = xmlPolicy
}

// addValidationTime is the port of the private
// addValidationTime(XmlSimpleReport).
func (b *SimpleReportBuilder) addValidationTime(report *jaxb.XmlSimpleReport) {
	report.ValidationTime = jaxb.NewXSDateTime(b.currentTime)
}

// addDocumentName is the port of the private
// addDocumentName(XmlSimpleReport).
func (b *SimpleReportBuilder) addDocumentName(report *jaxb.XmlSimpleReport) {
	documentName := b.diagnosticData.DocumentName()
	report.DocumentName = &documentName
}

// addContainerType is the port of the private
// addContainerType(XmlSimpleReport).
func (b *SimpleReportBuilder) addContainerType(simpleReport *jaxb.XmlSimpleReport) {
	containerType := jaxb.ASiCContainerTypeValue(b.diagnosticData.ContainerType())
	simpleReport.ContainerType = &containerType
}

// addSemantics is the port of the private addSemantics(XmlSimpleReport). The
// enumeration declaration order replaces Java's HashSet iteration order; see
// the file header.
func (b *SimpleReportBuilder) addSemantics(simpleReport *jaxb.XmlSimpleReport) {

	for _, indication := range enumerations.IndicationValues() {
		if _, used := b.finalIndications[indication]; !used {
			continue
		}
		semantic := &jaxb.XmlSemantic{}
		semantic.Key = string(indication)
		tag, _ := i18n.MessageTagGetSemantic(string(indication))
		semantic.Value = b.i18nProvider.GetMessage(tag)
		simpleReport.Semantic = append(simpleReport.Semantic, semantic)
	}

	for _, subIndication := range enumerations.SubIndicationValues() {
		if _, used := b.finalSubIndications[subIndication]; !used {
			continue
		}
		semantic := &jaxb.XmlSemantic{}
		semantic.Key = string(subIndication)
		tag, _ := i18n.MessageTagGetSemantic(string(subIndication))
		semantic.Value = b.i18nProvider.GetMessage(tag)
		simpleReport.Semantic = append(simpleReport.Semantic, semantic)
	}

}

// addStatistics is the port of the private addStatistics(XmlSimpleReport).
func (b *SimpleReportBuilder) addStatistics(simpleReport *jaxb.XmlSimpleReport) {
	simpleReport.ValidSignaturesCount = b.validSignatureCount
	simpleReport.SignaturesCount = b.totalSignatureCount
}

// addPDFAProfile is the port of the private addPDFAProfile(XmlSimpleReport).
func (b *SimpleReportBuilder) addPDFAProfile(simpleReport *jaxb.XmlSimpleReport) {
	pdfaProfileId := b.diagnosticData.PDFAProfileId()
	if pdfaProfileId != "" {
		xmlPDFAInfo := &jaxb.XmlPDFAInfo{}
		xmlPDFAInfo.PDFAProfile = &pdfaProfileId
		valid := b.diagnosticData.IsPDFACompliant()
		xmlPDFAInfo.Valid = &valid
		if len(b.diagnosticData.PDFAValidationErrors()) > 0 {
			xmlPDFAInfo.ValidationMessages = b.toXmlValidationMessages(b.diagnosticData.PDFAValidationErrors())
		}
		simpleReport.PDFAInfo = xmlPDFAInfo
	}
}

// toXmlValidationMessages is the port of the private
// toXmlValidationMessages(Collection<String>).
func (b *SimpleReportBuilder) toXmlValidationMessages(errors []string) *jaxb.XmlValidationMessages {
	xmlValidationMessages := &jaxb.XmlValidationMessages{}
	xmlValidationMessages.Error = append(xmlValidationMessages.Error, errors...)
	return xmlValidationMessages
}

// signature builds a XmlSignature object. Port of the private
// getSignature(SignatureWrapper, boolean).
func (b *SimpleReportBuilder) signature(signature *diagnostic.SignatureWrapper, container bool) *jaxb.XmlSignature {

	b.totalSignatureCount++

	signatureId := signature.Id()
	xmlSignature := &jaxb.XmlSignature{}
	xmlSignature.Id = signatureId

	b.addCounterSignature(signature, xmlSignature)
	b.addSignatureScope(signature, xmlSignature)
	b.addSigningTime(signature, xmlSignature)
	b.addBestSignatureTime(signature, xmlSignature)
	b.addSignatureFormat(signature, xmlSignature)

	signingCertificate := signature.SigningCertificate()
	if signingCertificate != nil {
		signedBy := b.readableCertificateName(signingCertificate.Id())
		xmlSignature.SignedBy = &signedBy
	}

	validationDetails := b.adESValidationDetails(signatureId)
	if b.isNotEmpty(validationDetails) {
		xmlSignature.AdESValidationDetails = validationDetails
	}

	qualificationDetails := b.qualificationDetails(signatureId)
	if b.isNotEmpty(qualificationDetails) {
		xmlSignature.QualificationDetails = qualificationDetails
	}

	if container {
		filename := signature.Filename()
		xmlSignature.Filename = &filename
	}

	indication := b.detailedReport.FinalIndication(signatureId)
	subIndication := b.detailedReport.FinalSubIndication(signatureId)
	if indication == enumerations.Indication_TOTAL_PASSED {
		b.determineSignatureExtensionPeriod(xmlSignature)
		b.validSignatureCount++

	} else if indication == enumerations.Indication_INDETERMINATE &&
		subIndication == enumerations.SubIndication_TRY_LATER {
		// indication is temporary, execute when applicable
		b.determineSignatureExtensionPeriod(xmlSignature)
	}

	xmlSignature.Indication = jaxb.IndicationValue(indication)
	b.finalIndications[indication] = struct{}{}

	if subIndication != "" {
		value := jaxb.SubIndicationValue(subIndication)
		xmlSignature.SubIndication = &value
		b.finalSubIndications[subIndication] = struct{}{}
	}

	b.addSignatureProfile(xmlSignature)

	xmlSignature.CertificateChain = b.certChain(signatureId)

	evidenceRecordList := signature.EvidenceRecords()
	if len(evidenceRecordList) > 0 {
		xmlEvidenceRecords := &jaxb.XmlEvidenceRecords{}
		for _, evidenceRecord := range evidenceRecordList {
			erIndication := b.detailedReport.EvidenceRecordValidationIndication(evidenceRecord.Id())
			if erIndication != "" {
				xmlEvidenceRecords.EvidenceRecord = append(xmlEvidenceRecords.EvidenceRecord, b.xmlEvidenceRecord(evidenceRecord))
			}
		}
		if len(xmlEvidenceRecords.EvidenceRecord) > 0 {
			xmlSignature.EvidenceRecords = xmlEvidenceRecords
		}
	}

	timestampList := signature.TimestampList()
	if len(timestampList) > 0 {
		xmlTimestamps := &jaxb.XmlTimestamps{}
		for _, timestamp := range timestampList {
			tstValidationIndication := b.detailedReport.BasicTimestampValidationIndication(timestamp.Id())
			if tstValidationIndication != "" {
				xmlTimestamps.Timestamp = append(xmlTimestamps.Timestamp, b.xmlTimestamp(timestamp))
			}
		}
		if len(xmlTimestamps.Timestamp) > 0 {
			xmlSignature.Timestamps = xmlTimestamps
		}
	}

	return xmlSignature
}

// adESValidationDetails is the port of the private
// getAdESValidationDetails(String).
func (b *SimpleReportBuilder) adESValidationDetails(tokenId string) *jaxb.XmlDetails {
	validationDetails := &jaxb.XmlDetails{}
	validationDetails.Error = append(validationDetails.Error, b.convert(b.detailedReport.AdESValidationErrors(tokenId))...)
	validationDetails.Warning = append(validationDetails.Warning, b.convert(b.detailedReport.AdESValidationWarnings(tokenId))...)
	validationDetails.Info = append(validationDetails.Info, b.convert(b.detailedReport.AdESValidationInfos(tokenId))...)
	return validationDetails
}

// qualificationDetails is the port of the private
// getQualificationDetails(String).
func (b *SimpleReportBuilder) qualificationDetails(tokenId string) *jaxb.XmlDetails {
	qualificationDetails := &jaxb.XmlDetails{}
	qualificationDetails.Error = append(qualificationDetails.Error, b.convert(b.detailedReport.QualificationErrors(tokenId))...)
	qualificationDetails.Warning = append(qualificationDetails.Warning, b.convert(b.detailedReport.QualificationWarnings(tokenId))...)
	qualificationDetails.Info = append(qualificationDetails.Info, b.convert(b.detailedReport.QualificationInfos(tokenId))...)
	return qualificationDetails
}

// isNotEmpty is the port of the private isNotEmpty(XmlDetails).
func (b *SimpleReportBuilder) isNotEmpty(details *jaxb.XmlDetails) bool {
	return len(details.Error) > 0 || len(details.Warning) > 0 || len(details.Info) > 0
}

// convert is the port of the private convert(Collection<Message>).
func (b *SimpleReportBuilder) convert(messages []detailedreport.Message) []*jaxb.XmlMessage {
	result := make([]*jaxb.XmlMessage, 0, len(messages))
	for _, m := range messages {
		xmlMessage := &jaxb.XmlMessage{}
		key := m.Key
		xmlMessage.Key = &key
		xmlMessage.Value = m.Value
		result = append(result, xmlMessage)
	}
	return result
}

// certChain is the port of the private getCertChain(String).
func (b *SimpleReportBuilder) certChain(tokenId string) *jaxb.XmlCertificateChain {
	certIds := b.detailedReport.BasicBuildingBlocksCertChain(tokenId)
	xmlCertificateChain := &jaxb.XmlCertificateChain{}
	if len(certIds) > 0 {
		for _, certId := range certIds {
			certificate := &jaxb.XmlCertificate{}
			certificate.Id = certId
			certificate.QualifiedName = b.readableCertificateName(certId)
			if b.isTrustAnchor(certId) {
				trusted := true
				certificate.Trusted = &trusted
				if sunsetDate := b.certificateSunsetDate(certId); sunsetDate != nil {
					certificate.SunsetDate = jaxb.NewXSDateTime(*sunsetDate)
				}
				certificate.TrustAnchors = b.xmlTrustAnchors(certId)
			}
			xmlCertificateChain.Certificate = append(xmlCertificateChain.Certificate, certificate)
		}
	}
	return xmlCertificateChain
}

// certificateSunsetDate is the port of the private
// getCertificateSunsetDate(String).
func (b *SimpleReportBuilder) certificateSunsetDate(certId string) *time.Time {
	certificate := b.diagnosticData.CertificateById(certId)
	if certificate != nil {
		return certificate.TrustSunsetDate()
	}
	return nil
}

// xmlTrustAnchors is the port of the private getXmlTrustAnchors(String).
func (b *SimpleReportBuilder) xmlTrustAnchors(certId string) *jaxb.XmlTrustAnchors {
	xmlTrustServiceProviders := b.filterByCertificateId(certId)
	if len(xmlTrustServiceProviders) > 0 {
		xmlTrustAnchors := &jaxb.XmlTrustAnchors{}
		for _, trustServiceProvider := range xmlTrustServiceProviders {
			trustAnchor := &jaxb.XmlTrustAnchor{}
			if trustServiceProvider.TL != nil {
				trustAnchor.CountryCode = trustServiceProvider.TL.CountryCode
				trustAnchor.TSLType = trustServiceProvider.TL.Type
			}
			trustAnchor.TrustServiceProvider = b.enOrFirst(trustServiceProvider.TSPNames.All())
			tspRegistrationIdentifiers := trustServiceProvider.TSPRegistrationIdentifiers.All()
			if len(tspRegistrationIdentifiers) > 0 {
				trustAnchor.TrustServiceProviderRegistrationId = &tspRegistrationIdentifiers[0]
			}
			trustAnchor.TrustServiceName = append(trustAnchor.TrustServiceName, b.uniqueServiceNames(trustServiceProvider)...)
			xmlTrustAnchors.TrustAnchor = append(xmlTrustAnchors.TrustAnchor, trustAnchor)
		}
		return xmlTrustAnchors
	}
	return nil
}

// uniqueServiceNames is the port of the private
// getUniqueServiceNames(XmlTrustServiceProvider). Java returns a HashSet;
// this port deduplicates while keeping first-seen order, see the file header.
// Java's Set may contain a null when a service has no name; the corresponding
// nil is kept as the empty string, which is what JAXB marshals for a null
// element of a List<String>.
func (b *SimpleReportBuilder) uniqueServiceNames(trustServiceProvider *diagnosticjaxb.XmlTrustServiceProvider) []string {
	var result []string
	seen := make(map[string]struct{})
	for _, xmlTrustService := range trustServiceProvider.TrustServices.All() {
		name := ""
		if value := b.enOrFirst(xmlTrustService.ServiceNames.All()); value != nil {
			name = *value
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return result
}

// enOrFirst is the port of the private getEnOrFirst(List<XmlLangAndValue>).
func (b *SimpleReportBuilder) enOrFirst(langAndValues []*diagnosticjaxb.XmlLangAndValue) *string {
	if len(langAndValues) > 0 {
		for _, langAndValue := range langAndValues {
			if langAndValue.Lang != nil && strings.EqualFold(*langAndValue.Lang, "en") {
				value := langAndValue.Value
				return &value
			}
		}
		value := langAndValues[0].Value
		return &value
	}
	return nil
}

// filterByCertificateId is the port of the private
// filterByCertificateId(String).
func (b *SimpleReportBuilder) filterByCertificateId(certId string) []*diagnosticjaxb.XmlTrustServiceProvider {
	certificate := b.diagnosticData.CertificateById(certId)
	var result []*diagnosticjaxb.XmlTrustServiceProvider
	for _, xmlTrustServiceProvider := range certificate.TrustServiceProviders() {
		trustServices := xmlTrustServiceProvider.TrustServices.All()
		foundCertId := false
		for _, xmlTrustService := range trustServices {
			// Java compares with Utils.areStringsEqual, which treats two nulls
			// as equal; certId is never null at this call site, so a missing
			// service digital identifier simply does not match.
			if xmlTrustService.ServiceDigitalIdentifier != nil &&
				xmlTrustService.ServiceDigitalIdentifier.Id != nil &&
				certId == string(*xmlTrustService.ServiceDigitalIdentifier.Id) {
				foundCertId = true
				break
			}
		}
		if foundCertId {
			result = append(result, xmlTrustServiceProvider)
		}
	}
	return result
}

// addBestSignatureTime is the port of the private
// addBestSignatureTime(SignatureWrapper, XmlSignature).
func (b *SimpleReportBuilder) addBestSignatureTime(signature *diagnostic.SignatureWrapper, xmlSignature *jaxb.XmlSignature) {
	if bestSignatureTime := b.detailedReport.BestSignatureTime(signature.Id()); bestSignatureTime != nil {
		xmlSignature.BestSignatureTime = jaxb.NewXSDateTime(*bestSignatureTime)
	}
}

// addCounterSignature is the port of the private
// addCounterSignature(SignatureWrapper, XmlSignature).
func (b *SimpleReportBuilder) addCounterSignature(signature *diagnostic.SignatureWrapper, xmlSignature *jaxb.XmlSignature) {
	if signature.IsCounterSignature() {
		counterSignature := true
		xmlSignature.CounterSignature = &counterSignature
		parentId := signature.Parent().Id()
		xmlSignature.ParentId = &parentId
	}
}

// addSignatureScope is the port of the private
// addSignatureScope(SignatureWrapper, XmlSignature).
func (b *SimpleReportBuilder) addSignatureScope(signature *diagnostic.SignatureWrapper, xmlSignature *jaxb.XmlSignature) {
	signatureScopes := signature.SignatureScopes()
	if len(signatureScopes) > 0 {
		for _, signatureScope := range signatureScopes {
			xmlSignature.SignatureScope = append(xmlSignature.SignatureScope, b.xmlSignatureScope(signatureScope))
		}
	}
}

// xmlSignatureScope is the port of the private
// getXmlSignatureScope(eu.europa.esig.dss.diagnostic.jaxb.XmlSignatureScope).
func (b *SimpleReportBuilder) xmlSignatureScope(signatureScope *diagnosticjaxb.XmlSignatureScope) *jaxb.XmlSignatureScope {
	xmlSignatureScope := &jaxb.XmlSignatureScope{}
	if signatureScope.SignerData != nil && signatureScope.SignerData.Id != nil {
		xmlSignatureScope.Id = string(*signatureScope.SignerData.Id)
	}
	xmlSignatureScope.Name = signatureScope.Name
	if signatureScope.Scope != nil {
		scope := jaxb.SignatureScopeTypeValue(*signatureScope.Scope)
		xmlSignatureScope.Scope = &scope
	}
	if signatureScope.Description != nil {
		xmlSignatureScope.Value = *signatureScope.Description
	}
	return xmlSignatureScope
}

// addSigningTime is the port of the private
// addSigningTime(SignatureWrapper, XmlSignature).
func (b *SimpleReportBuilder) addSigningTime(signature *diagnostic.SignatureWrapper, xmlSignature *jaxb.XmlSignature) {
	if claimedSigningTime := signature.ClaimedSigningTime(); claimedSigningTime != nil {
		xmlSignature.SigningTime = jaxb.NewXSDateTime(*claimedSigningTime)
	}
}

// addSignatureFormat is the port of the private
// addSignatureFormat(SignatureWrapper, XmlSignature).
func (b *SimpleReportBuilder) addSignatureFormat(signature *diagnostic.SignatureWrapper, xmlSignature *jaxb.XmlSignature) {
	xmlSignature.SignatureFormat = jaxb.SignatureLevelValue(signature.SignatureFormat())
}

// readableCertificateName is the port of the private
// getReadableCertificateName(String).
func (b *SimpleReportBuilder) readableCertificateName(certId string) string {
	certificateWrapper := b.diagnosticData.UsedCertificateByIdNullSafe(certId)
	return certificateWrapper.ReadableCertificateName()
}

// isTrustAnchor is the port of the private isTrustAnchor(String).
func (b *SimpleReportBuilder) isTrustAnchor(certId string) bool {
	certificateWrapper := b.diagnosticData.UsedCertificateByIdNullSafe(certId)
	return certificateWrapper.IsTrusted()
}

// addSignatureProfile is the port of the private
// addSignatureProfile(XmlSignature).
func (b *SimpleReportBuilder) addSignatureProfile(xmlSignature *jaxb.XmlSignature) {
	qualification := b.detailedReport.SignatureQualification(xmlSignature.Id)
	if qualification != "" {
		sigLevel := &jaxb.XmlSignatureLevel{}
		sigLevel.Value = jaxb.SignatureQualificationValue(qualification)
		description := qualification.Label()
		sigLevel.Description = &description
		xmlSignature.SignatureLevel = sigLevel
	}
}

// xmlTimestamp is the port of the private getXmlTimestamp(TimestampWrapper).
func (b *SimpleReportBuilder) xmlTimestamp(timestampWrapper *diagnostic.TimestampWrapper) *jaxb.XmlTimestamp {
	xmlTimestamp := &jaxb.XmlTimestamp{}
	timestampId := timestampWrapper.Id()
	xmlTimestamp.Id = timestampId
	if productionTime := timestampWrapper.ProductionTime(); productionTime != nil {
		xmlTimestamp.ProductionTime = jaxb.NewXSDateTime(*productionTime)
	}
	producedBy := b.producedByName(timestampWrapper)
	xmlTimestamp.ProducedBy = &producedBy
	xmlTimestamp.CertificateChain = b.certChain(timestampId)
	filename := timestampWrapper.Filename()
	xmlTimestamp.Filename = &filename

	indication := b.detailedReport.FinalIndication(timestampId)
	xmlTimestamp.Indication = jaxb.IndicationValue(indication)
	b.finalIndications[indication] = struct{}{}

	subIndication := b.detailedReport.FinalSubIndication(timestampId)
	if subIndication != "" {
		value := jaxb.SubIndicationValue(subIndication)
		xmlTimestamp.SubIndication = &value
		b.finalSubIndications[subIndication] = struct{}{}
	}

	timestampQualification := b.detailedReport.TimestampQualification(timestampId)
	if timestampQualification != "" {
		xmlTimestampLevel := &jaxb.XmlTimestampLevel{}
		xmlTimestampLevel.Value = jaxb.TimestampQualificationValue(timestampQualification)
		description := timestampQualification.Label()
		xmlTimestampLevel.Description = &description
		xmlTimestamp.TimestampLevel = xmlTimestampLevel
	}

	validationDetails := b.adESValidationDetails(timestampId)
	if b.isNotEmpty(validationDetails) {
		xmlTimestamp.AdESValidationDetails = validationDetails
	}

	qualificationDetails := b.qualificationDetails(timestampId)
	if b.isNotEmpty(qualificationDetails) {
		xmlTimestamp.QualificationDetails = qualificationDetails
	}

	if len(timestampWrapper.TimestampScopes()) > 0 {
		for _, timestampScope := range timestampWrapper.TimestampScopes() {
			xmlTimestamp.TimestampScope = append(xmlTimestamp.TimestampScope, b.xmlSignatureScope(timestampScope))
		}
	}

	if len(timestampWrapper.EvidenceRecords()) > 0 {
		xmlEvidenceRecords := &jaxb.XmlEvidenceRecords{}
		for _, evidenceRecord := range timestampWrapper.EvidenceRecords() {
			erIndication := b.detailedReport.EvidenceRecordValidationIndication(evidenceRecord.Id())
			if erIndication != "" {
				xmlEvidenceRecords.EvidenceRecord = append(xmlEvidenceRecords.EvidenceRecord, b.xmlEvidenceRecord(evidenceRecord))
			}
		}
		if len(xmlEvidenceRecords.EvidenceRecord) > 0 {
			xmlTimestamp.EvidenceRecords = xmlEvidenceRecords
		}
	}

	return xmlTimestamp
}

// xmlEvidenceRecord is the port of the private
// getXmlEvidenceRecord(EvidenceRecordWrapper).
func (b *SimpleReportBuilder) xmlEvidenceRecord(evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper) *jaxb.XmlEvidenceRecord {
	xmlEvidenceRecord := &jaxb.XmlEvidenceRecord{}

	evidenceRecordId := evidenceRecordWrapper.Id()
	xmlEvidenceRecord.Id = evidenceRecordId
	filename := evidenceRecordWrapper.Filename()
	xmlEvidenceRecord.Filename = &filename

	if poeTime := b.detailedReport.EvidenceRecordLowestPOETime(evidenceRecordId); poeTime != nil {
		xmlEvidenceRecord.POETime = jaxb.NewXSDateTime(*poeTime)
	}

	indication := b.detailedReport.FinalIndication(evidenceRecordId)
	xmlEvidenceRecord.Indication = jaxb.IndicationValue(indication)
	b.finalIndications[indication] = struct{}{}

	subIndication := b.detailedReport.FinalSubIndication(evidenceRecordId)
	if subIndication != "" {
		value := jaxb.SubIndicationValue(subIndication)
		xmlEvidenceRecord.SubIndication = &value
		b.finalSubIndications[subIndication] = struct{}{}
	}

	validationDetails := b.adESValidationDetails(evidenceRecordId)
	if b.isNotEmpty(validationDetails) {
		xmlEvidenceRecord.AdESValidationDetails = validationDetails
	}

	if b.isValidConclusion(evidenceRecordId) {
		b.determineEvidenceRecordExtensionPeriod(xmlEvidenceRecord)
	}

	if len(evidenceRecordWrapper.EvidenceRecordScopes()) > 0 {
		for _, timestampScope := range evidenceRecordWrapper.EvidenceRecordScopes() {
			xmlEvidenceRecord.EvidenceRecordScope = append(xmlEvidenceRecord.EvidenceRecordScope, b.xmlSignatureScope(timestampScope))
		}
	}

	timestampList := evidenceRecordWrapper.TimestampList()
	if len(timestampList) > 0 {
		xmlTimestamps := &jaxb.XmlTimestamps{}
		for _, timestamp := range timestampList {
			tstValidationIndication := b.detailedReport.BasicTimestampValidationIndication(timestamp.Id())
			if tstValidationIndication != "" {
				xmlTimestamps.Timestamp = append(xmlTimestamps.Timestamp, b.xmlTimestamp(timestamp))
			}
		}
		if len(xmlTimestamps.Timestamp) > 0 {
			xmlEvidenceRecord.Timestamps = xmlTimestamps
		}
	}

	if evidenceRecordWrapper.IsEmbedded() {
		embedded := true
		xmlEvidenceRecord.Embedded = &embedded
		parentId := evidenceRecordWrapper.Parent().Id()
		xmlEvidenceRecord.ParentId = &parentId
	}

	return xmlEvidenceRecord
}

// producedByName is the port of the private
// getProducedByName(TimestampWrapper).
func (b *SimpleReportBuilder) producedByName(timestampWrapper *diagnostic.TimestampWrapper) string {
	signingCertificate := timestampWrapper.SigningCertificate()
	if signingCertificate != nil {
		return signingCertificate.ReadableCertificateName()
	}
	return ""
}

// determineSignatureExtensionPeriod is the port of the private
// determineExtensionPeriod(XmlSignature).
func (b *SimpleReportBuilder) determineSignatureExtensionPeriod(xmlSignature *jaxb.XmlSignature) {
	signatureWrapper := b.diagnosticData.SignatureById(xmlSignature.Id)
	timestampList := b.allTimestampsForSignature(signatureWrapper)
	if min := b.minExtensionPeriod(signatureWrapper, timestampList); min != nil {
		xmlSignature.ExtensionPeriodMin = jaxb.NewXSDateTime(*min)
	}
	if max := b.maxExtensionPeriod(signatureWrapper, timestampList); max != nil {
		xmlSignature.ExtensionPeriodMax = jaxb.NewXSDateTime(*max)
	}
}

// allTimestampsForSignature is the port of the private
// getAllTimestampsForSignature(SignatureWrapper).
func (b *SimpleReportBuilder) allTimestampsForSignature(signatureWrapper *diagnostic.SignatureWrapper) []*diagnostic.TimestampWrapper {
	timestampList := append([]*diagnostic.TimestampWrapper{}, signatureWrapper.AllTimestampsProducedAfterSignatureCreation()...)
	timestampList = append(timestampList, b.evidenceRecordTimestampsForTokenWithId(signatureWrapper.Id())...)
	return timestampList
}

// determineTimestampExtensionPeriod is the port of the private
// determineExtensionPeriod(XmlTimestamp).
func (b *SimpleReportBuilder) determineTimestampExtensionPeriod(xmlTimestamp *jaxb.XmlTimestamp) {
	timestampWrapper := b.diagnosticData.TimestampById(xmlTimestamp.Id)
	timestampList := b.allTimestampsForTimestamp(timestampWrapper)
	if min := b.minExtensionPeriod(timestampWrapper, timestampList); min != nil {
		xmlTimestamp.ExtensionPeriodMin = jaxb.NewXSDateTime(*min)
	}
	if max := b.maxExtensionPeriod(timestampWrapper, timestampList); max != nil {
		xmlTimestamp.ExtensionPeriodMax = jaxb.NewXSDateTime(*max)
	}
}

// allTimestampsForTimestamp is the port of the private
// getAllTimestampsForTimestamp(TimestampWrapper).
func (b *SimpleReportBuilder) allTimestampsForTimestamp(timestampWrapper *diagnostic.TimestampWrapper) []*diagnostic.TimestampWrapper {
	var timestampList []*diagnostic.TimestampWrapper

	for _, timestamp := range b.diagnosticData.TimestampList() {
		timestampedObjects := timestamp.TimestampedObjects()
		if b.containObjectWithId(timestampedObjects, timestampWrapper.Id()) {
			timestampList = append(timestampList, timestamp)
		}
	}

	timestampList = append(timestampList, b.evidenceRecordTimestampsForTokenWithId(timestampWrapper.Id())...)

	return timestampList
}

// determineEvidenceRecordExtensionPeriod is the port of the private
// determineExtensionPeriod(XmlEvidenceRecord).
func (b *SimpleReportBuilder) determineEvidenceRecordExtensionPeriod(xmlEvidenceRecord *jaxb.XmlEvidenceRecord) {
	evidenceRecordWrapper := b.diagnosticData.EvidenceRecordById(xmlEvidenceRecord.Id)
	timestampList := b.allTimestampsForEvidenceRecord(evidenceRecordWrapper)
	if min := b.minExtensionPeriodForTimestampList(timestampList); min != nil {
		xmlEvidenceRecord.ExtensionPeriodMin = jaxb.NewXSDateTime(*min)
	}
	if max := b.maxExtensionPeriodForTimestampList(timestampList); max != nil {
		xmlEvidenceRecord.ExtensionPeriodMax = jaxb.NewXSDateTime(*max)
	}
}

// allTimestampsForEvidenceRecord is the port of the private
// getAllTimestampsForEvidenceRecord(EvidenceRecordWrapper).
func (b *SimpleReportBuilder) allTimestampsForEvidenceRecord(
	evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper) []*diagnostic.TimestampWrapper {
	timestampList := append([]*diagnostic.TimestampWrapper{}, evidenceRecordWrapper.TimestampList()...)
	timestampList = append(timestampList, b.evidenceRecordTimestampsForTokenWithId(evidenceRecordWrapper.Id())...)
	return timestampList
}

// evidenceRecordTimestampsForTokenWithId is the port of the private
// getEvidenceRecordTimestampsForTokenWithId(String).
func (b *SimpleReportBuilder) evidenceRecordTimestampsForTokenWithId(tokenId string) []*diagnostic.TimestampWrapper {
	var timestampList []*diagnostic.TimestampWrapper

	for _, evidenceRecord := range b.diagnosticData.EvidenceRecords() {
		if !b.isValidConclusion(evidenceRecord.Id()) {
			continue
		}

		coveredObjects := evidenceRecord.CoveredObjects()
		if b.containObjectWithId(coveredObjects, tokenId) {
			timestampList = append(timestampList, evidenceRecord.TimestampList()...)
		}
	}

	return timestampList
}

// eaa is the port of the private getEAA(EAAWrapper).
func (b *SimpleReportBuilder) eaa(eaaWrapper *diagnostic.EAAWrapper) *jaxb.XmlEAA {
	xmlEAA := &jaxb.XmlEAA{}

	eaaId := eaaWrapper.Id()
	xmlEAA.Id = eaaId
	filename := eaaWrapper.Filename()
	xmlEAA.Filename = &filename

	indication := b.detailedReport.FinalIndication(eaaId)
	xmlEAA.Indication = jaxb.IndicationValue(indication)
	b.finalIndications[indication] = struct{}{}

	subIndication := b.detailedReport.FinalSubIndication(eaaId)
	if subIndication != "" {
		value := jaxb.SubIndicationValue(subIndication)
		xmlEAA.SubIndication = &value
		b.finalSubIndications[subIndication] = struct{}{}
	}

	eaaQualifications := b.detailedReport.EAAQualifications(eaaId)
	if len(eaaQualifications) > 0 {
		xmlEAA.EAALevel = append(xmlEAA.EAALevel, b.xmlEAALevels(eaaQualifications)...)
	}

	validationDetails := b.adESValidationDetails(eaaId)
	if b.isNotEmpty(validationDetails) {
		xmlEAA.AdESValidationDetails = validationDetails
	}

	qualificationDetails := b.qualificationDetails(eaaId)
	if b.isNotEmpty(qualificationDetails) {
		xmlEAA.QualificationDetails = qualificationDetails
	}

	signatures := eaaWrapper.EAASignatures()
	if len(signatures) > 0 {
		for _, signature := range signatures {
			xmlEAA.EAASignature = append(xmlEAA.EAASignature, b.signature(signature, false))
		}
	}
	keyBindingSignature := eaaWrapper.KeyBindingSignature()
	if keyBindingSignature != nil {
		xmlEAA.KeyBindingSignature = b.signature(keyBindingSignature, false)
	}

	xmlEAA.EAAPayload = b.buildXmlEAAPayload(eaaWrapper)

	return xmlEAA
}

// xmlEAALevels is the port of the private
// getXmlEAALevels(List<EAAQualification>).
func (b *SimpleReportBuilder) xmlEAALevels(eaaQualifications []enumerations.EAAQualification) []*jaxb.XmlEAALevel {
	if len(eaaQualifications) == 0 {
		return nil
	}
	var result []*jaxb.XmlEAALevel
	for _, eaaQualification := range eaaQualifications {
		xmlEAALevel := &jaxb.XmlEAALevel{}
		xmlEAALevel.Value = jaxb.EAAQualificationValue(eaaQualification)
		description := eaaQualification.Label()
		xmlEAALevel.Description = &description
		result = append(result, xmlEAALevel)
	}
	return result
}

// minExtensionPeriod is the port of the private
// getMinExtensionPeriod(AbstractTokenProxy, List<TimestampWrapper>).
func (b *SimpleReportBuilder) minExtensionPeriod(token diagnostic.TokenProxy,
	timestampList []*diagnostic.TimestampWrapper) *time.Time {
	var min *time.Time

	var chains [][]*diagnostic.CertificateWrapper
	chains = append(chains, token.CertificateChain())
	relatedRevocations := b.relatedRevocations(token)
	for _, revocation := range relatedRevocations {
		chains = append(chains, revocation.CertificateChain())
	}

	for _, certificateChain := range chains {
		certChainMin := b.minExtensionPeriodForChain(certificateChain, nil)
		if certChainMin != nil {
			if min == nil || min.Before(*certChainMin) {
				min = certChainMin
			}
		}
	}

	minExtensionPeriodForTimestampList := b.minExtensionPeriodForTimestampList(timestampList)
	if minExtensionPeriodForTimestampList != nil && (min == nil || min.Before(*minExtensionPeriodForTimestampList)) {
		min = minExtensionPeriodForTimestampList
	}

	return min
}

// relatedRevocations is the port of the private
// getRelatedRevocations(AbstractTokenProxy). It returns the revocation data
// protected by a time-stamp; other revocation data need not be processed.
func (b *SimpleReportBuilder) relatedRevocations(token diagnostic.TokenProxy) []*diagnostic.RelatedRevocationWrapper {
	var result []*diagnostic.RelatedRevocationWrapper
	for _, r := range token.FoundRevocations().RelatedRevocationData() {
		if b.poe.GetLowestPOE(r.Id()).IsTokenProvided() {
			result = append(result, r)
		}
	}
	return result
}

// minExtensionPeriodForTimestampList is the port of the private
// getMinExtensionPeriodForTimestampList(List<TimestampWrapper>).
func (b *SimpleReportBuilder) minExtensionPeriodForTimestampList(timestampList []*diagnostic.TimestampWrapper) *time.Time {
	var min *time.Time

	for _, timestampWrapper := range timestampList {
		certChainMin := b.minExtensionPeriodForChain(timestampWrapper.CertificateChain(), timestampWrapper.ProductionTime())
		if certChainMin != nil {
			if min == nil || min.Before(*certChainMin) {
				min = certChainMin
			}
		}
	}

	return min
}

// minExtensionPeriodForChain is the port of the private
// getMinExtensionPeriodForChain(List<CertificateWrapper>, Date).
func (b *SimpleReportBuilder) minExtensionPeriodForChain(certificateChain []*diagnostic.CertificateWrapper,
	usageTime *time.Time) *time.Time {
	var min *time.Time
	for _, certificateWrapper := range certificateChain {
		if certificateWrapper.IsTrusted() || certificateWrapper.IsSelfSigned() {
			break
		}

		if len(certificateWrapper.CertificateRevocationData()) > 0 {
			var lastTrustedUsage time.Time
			if usageTime != nil {
				lastTrustedUsage = *usageTime
			} else {
				lastTrustedUsage = b.poe.GetLowestPOETime(certificateWrapper.Id())
			}

			var tempMin *time.Time
			goodRevocationFound := false

			certificateRevocationData := certificateWrapper.CertificateRevocationData()
			for _, revocationData := range certificateRevocationData {
				// Revocation data shall be issued after the POE time
				thisUpdate := revocationData.ThisUpdate()
				if thisUpdate != nil && lastTrustedUsage.Before(*thisUpdate) {
					goodRevocationFound = true
					break

				} else {
					nextUpdate := revocationData.NextUpdate()
					if nextUpdate == nil {
						// last usage + 1s
						shifted := lastTrustedUsage.Add(time.Second)
						nextUpdate = &shifted
					}

					// find the minimum for the certificate across related revocations
					if tempMin == nil || tempMin.After(*nextUpdate) {
						tempMin = nextUpdate
					}
				}
			}

			if goodRevocationFound {
				continue
			}

			// find maximum across all certificates in the chain
			if tempMin != nil {
				if min == nil || min.Before(*tempMin) {
					min = tempMin
				}
			}
		}

	}
	return min
}

// maxExtensionPeriod is the port of the private
// getMaxExtensionPeriod(AbstractTokenProxy, List<TimestampWrapper>).
func (b *SimpleReportBuilder) maxExtensionPeriod(token diagnostic.TokenProxy,
	timestampList []*diagnostic.TimestampWrapper) *time.Time {
	var max *time.Time

	signingCertificate := token.SigningCertificate()
	if signingCertificate != nil {
		max = signingCertificate.NotAfter()
	} else {
		return nil
	}

	maxTimestampExtensionPeriod := b.maxExtensionPeriodForTimestampList(timestampList)
	if maxTimestampExtensionPeriod != nil && (max == nil || maxTimestampExtensionPeriod.After(*max)) {
		max = maxTimestampExtensionPeriod
	}

	return b.ensureMaxExtensionTimeIsAfterValidationTime(max)
}

// maxExtensionPeriodForTimestampList is the port of the private
// getMaxExtensionPeriodForTimestampList(List<TimestampWrapper>).
func (b *SimpleReportBuilder) maxExtensionPeriodForTimestampList(timestampList []*diagnostic.TimestampWrapper) *time.Time {
	var max *time.Time

	for _, timestampWrapper := range timestampList {
		if !b.isValidConclusion(timestampWrapper.Id()) {
			continue
		}
		timestampSigningCertificate := timestampWrapper.SigningCertificate()
		if timestampSigningCertificate != nil {
			notAfter := timestampSigningCertificate.NotAfter()
			if notAfter != nil && (max == nil || notAfter.After(*max)) {
				max = notAfter
			}
		}
	}

	return b.ensureMaxExtensionTimeIsAfterValidationTime(max)
}

// ensureMaxExtensionTimeIsAfterValidationTime is the port of the private
// ensureMaxExtensionTimeIsAfterValidationTime(Date).
func (b *SimpleReportBuilder) ensureMaxExtensionTimeIsAfterValidationTime(maxExtensionTime *time.Time) *time.Time {
	if maxExtensionTime == nil {
		return nil
	}
	if b.currentTime.Before(*maxExtensionTime) {
		return maxExtensionTime
	}
	return nil
}

// isValidConclusion is the port of the private isValidConclusion(String).
func (b *SimpleReportBuilder) isValidConclusion(tokenId string) bool {
	finalIndication := b.detailedReport.FinalIndication(tokenId)
	return enumerations.Indication_TOTAL_PASSED == finalIndication || enumerations.Indication_PASSED == finalIndication
}

// containObjectWithId is the port of the private
// containObjectWithId(List<XmlTimestampedObject>, String).
func (b *SimpleReportBuilder) containObjectWithId(timestampedObjects []*diagnosticjaxb.XmlTimestampedObject,
	tokenId string) bool {
	for _, o := range timestampedObjects {
		if o.Token != nil && tokenId == o.Token.ID {
			return true
		}
	}
	return false
}

// buildXmlEAAPayload is the port of the private
// buildXmlEAAPayload(EAAWrapper).
func (b *SimpleReportBuilder) buildXmlEAAPayload(eaaWrapper *diagnostic.EAAWrapper) *jaxb.XmlEAAPayload {
	xmlEAAPayload := &jaxb.XmlEAAPayload{}

	eaaPayloadProxy := eaaWrapper.EAAPayload()
	xmlEAAPayload.Identifier = b.xmlDisclosableClaim(eaaPayloadProxy.EAAIdentifier())
	xmlEAAPayload.Issuer = b.xmlDisclosableClaim(eaaPayloadProxy.EAAIssuer())
	xmlEAAPayload.Subject = b.xmlDisclosableClaim(eaaPayloadProxy.EAASubject())
	xmlEAAPayload.Audience = b.xmlDisclosableClaim(eaaPayloadProxy.EAAAudience())
	xmlEAAPayload.Expiration = b.xmlDisclosableClaim(eaaPayloadProxy.EAAExpiration())
	xmlEAAPayload.NotBefore = b.xmlDisclosableClaim(eaaPayloadProxy.EAANotBefore())
	xmlEAAPayload.IssuedAt = b.xmlDisclosableClaim(eaaPayloadProxy.EAAIssuedAt())
	xmlEAAPayload.UpdatedAt = b.xmlDisclosableClaim(eaaPayloadProxy.EAAUpdatedAt())
	// Java assigns UpdatedAt twice in a row; kept for statement fidelity.
	xmlEAAPayload.UpdatedAt = b.xmlDisclosableClaim(eaaPayloadProxy.EAAUpdatedAt())
	xmlEAAPayload.Category = b.xmlDisclosableClaim(eaaPayloadProxy.EAACategory())
	xmlEAAPayload.VerifiableCredentialsType = b.xmlDisclosableClaim(eaaPayloadProxy.EAAVerifiableCredentialsType())
	eaaStatus := eaaPayloadProxy.EAAStatus()
	if eaaStatus != nil {
		disclosable := eaaStatus.IsSelectivelyDisclosable()
		xmlEAAPayload.StatusIndex = b.xmlDisclosableClaimWithDisclosure(eaaStatus.Index(), &disclosable)
		xmlEAAPayload.StatusUri = b.xmlDisclosableClaimWithDisclosure(eaaStatus.Uri(), &disclosable)
		xmlEAAPayload.StatusType = b.xmlDisclosableClaimWithDisclosure(eaaStatus.Type(), &disclosable)
		xmlEAAPayload.StatusPurpose = b.xmlDisclosableClaimWithDisclosure(eaaStatus.Purpose(), &disclosable)
	}
	xmlEAAPayload.Nonce = b.xmlDisclosableClaim(eaaPayloadProxy.EAANonce())
	eaaDeviceKey := eaaPayloadProxy.EAADeviceKey()
	if eaaDeviceKey != nil && eaaDeviceKey.PublicKey() != nil {
		disclosable := eaaDeviceKey.IsSelectivelyDisclosable()
		name := eaaDeviceKey.Name()
		xmlEAAPayload.DeviceKey = b.newXmlDisclosableClaim(&name, &disclosable,
			base64.StdEncoding.EncodeToString(eaaDeviceKey.PublicKey()))
	}
	xmlEAAPayload.Version = b.xmlDisclosableClaim(eaaPayloadProxy.EAAVersion())
	xmlEAAPayload.DocType = b.xmlDisclosableClaim(eaaPayloadProxy.EAADocType())
	eaaValidityInfo := eaaPayloadProxy.EAAValidityInfo()
	if eaaValidityInfo != nil {
		disclosable := eaaValidityInfo.IsSelectivelyDisclosable()
		xmlEAAPayload.IssuedAt = b.xmlDisclosableClaimWithDisclosure(eaaValidityInfo.Signed(), &disclosable)
		xmlEAAPayload.NotBefore = b.xmlDisclosableClaimWithDisclosure(eaaValidityInfo.ValidFrom(), &disclosable)
		xmlEAAPayload.AdministrativeExpirationDate = b.xmlDisclosableClaimWithDisclosure(eaaValidityInfo.ValidUntil(), &disclosable)
		xmlEAAPayload.NextUpdate = b.xmlDisclosableClaimWithDisclosure(eaaValidityInfo.ExpectedUpdate(), &disclosable)
	}
	xmlEAAPayload.AdministrativeIssuanceDate = b.xmlDisclosableClaim(eaaPayloadProxy.AdministrativeIssuanceDate())
	xmlEAAPayload.AdministrativeExpirationDate = b.xmlDisclosableClaim(eaaPayloadProxy.AdministrativeExpirationDate())
	xmlEAAPayload.OneTimeUse = b.xmlDisclosableClaim(eaaPayloadProxy.OneTimeUse())
	xmlEAAPayload.ShortLived = b.xmlDisclosableClaim(eaaPayloadProxy.ShortLived())
	xmlEAAPayload.Evidence = b.xmlDisclosableClaim(eaaPayloadProxy.Evidence())
	if eaaPayloadProxy.AttestedAttributesSubject() != nil {
		subjectId := eaaPayloadProxy.AttestedAttributesSubject().SubjectId()
		if subjectId != nil {
			if subjectId.Text() != "" {
				xmlEAAPayload.AttestedAttributesSubjectId = b.xmlDisclosableClaim(subjectId.AsClaim())
			}
			if subjectId.FamilyName() != nil {
				xmlEAAPayload.AttestedAttributesSubjectFamilyName = b.xmlDisclosableClaim(subjectId.FamilyName())
			}
			if subjectId.GivenName() != nil {
				xmlEAAPayload.AttestedAttributesSubjectGivenName = b.xmlDisclosableClaim(subjectId.GivenName())
			}
			if subjectId.DocumentNumber() != nil {
				xmlEAAPayload.AttestedAttributesSubjectDocumentNumber = b.xmlDisclosableClaim(subjectId.DocumentNumber())
			}
		}
		xmlEAAPayload.AttestedAttributesSubjectPseudonym = b.xmlDisclosableClaim(eaaPayloadProxy.AttestedAttributesSubject().SubjectPseudonym())
		xmlEAAPayload.AttestedAttributes = b.xmlDisclosableClaim(eaaPayloadProxy.AttestedAttributesSubject().Attributes())
	}

	xmlEAAPayload.FullName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderFullName())
	xmlEAAPayload.GivenName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderGivenName())
	xmlEAAPayload.FamilyName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderFamilyName())
	xmlEAAPayload.MiddleName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderMiddleName())
	xmlEAAPayload.Nickname = b.xmlDisclosableClaim(eaaPayloadProxy.HolderNickname())
	xmlEAAPayload.ShortName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderShortName())
	xmlEAAPayload.ProfileUrl = b.xmlDisclosableClaim(eaaPayloadProxy.HolderProfileUrl())
	xmlEAAPayload.PictureUrl = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPictureUrl())
	xmlEAAPayload.WebsiteUrl = b.xmlDisclosableClaim(eaaPayloadProxy.HolderWebsiteUrl())
	xmlEAAPayload.Email = b.xmlDisclosableClaim(eaaPayloadProxy.HolderEmail())
	xmlEAAPayload.EmailVerified = b.xmlDisclosableClaim(eaaPayloadProxy.HolderEmailVerified())
	xmlEAAPayload.Gender = b.xmlDisclosableClaim(eaaPayloadProxy.HolderGender())
	if eaaPayloadProxy.HolderBirthdate() != nil {
		xmlEAAPayload.Birthdate = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthdate().Birthdate())
		xmlEAAPayload.BirthdateApproximateMask = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthdate().ApproximateMask())
	}
	xmlEAAPayload.Timezone = b.xmlDisclosableClaim(eaaPayloadProxy.HolderTimezone())
	xmlEAAPayload.Locale = b.xmlDisclosableClaim(eaaPayloadProxy.HolderLocale())
	userAddress := eaaPayloadProxy.HolderAddress()
	if userAddress != nil {
		disclosable := userAddress.IsSelectivelyDisclosable()
		xmlEAAPayload.AddressPostalAddress = b.xmlDisclosableClaimWithDisclosure(userAddress.PostalAddress(), &disclosable)
		xmlEAAPayload.AddressCity = b.xmlDisclosableClaimWithDisclosure(userAddress.City(), &disclosable)
		xmlEAAPayload.AddressCountryName = b.xmlDisclosableClaimWithDisclosure(userAddress.Country(), &disclosable)
		xmlEAAPayload.AddressPostalCode = b.xmlDisclosableClaimWithDisclosure(userAddress.PostalCode(), &disclosable)
		xmlEAAPayload.AddressStateOrProvince = b.xmlDisclosableClaimWithDisclosure(userAddress.StateOrProvince(), &disclosable)
		xmlEAAPayload.AddressStreetAddress = b.xmlDisclosableClaimWithDisclosure(userAddress.StreetAddress(), &disclosable)
	}
	xmlEAAPayload.PhoneNumber = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPhoneNumber())
	xmlEAAPayload.PhoneNumberVerified = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPhoneNumberVerified())
	userPlaceOfBirth := eaaPayloadProxy.HolderPlaceOfBirth()
	if userPlaceOfBirth != nil {
		disclosable := userPlaceOfBirth.IsSelectivelyDisclosable()
		xmlEAAPayload.PlaceOfBirth = b.xmlDisclosableClaimWithDisclosure(userPlaceOfBirth.AsClaim(), &disclosable)
		xmlEAAPayload.PlaceOfBirthCity = b.xmlDisclosableClaimWithDisclosure(userPlaceOfBirth.City(), &disclosable)
		xmlEAAPayload.PlaceOfBirthCountry = b.xmlDisclosableClaimWithDisclosure(userPlaceOfBirth.Country(), &disclosable)
		xmlEAAPayload.PlaceOfBirthRegion = b.xmlDisclosableClaimWithDisclosure(userPlaceOfBirth.Region(), &disclosable)
	}
	xmlEAAPayload.Nationalities = b.xmlDisclosableClaim(eaaPayloadProxy.HolderNationalities())
	xmlEAAPayload.BirthFamilyName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthFamilyName())
	xmlEAAPayload.BirthGivenName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthGivenName())
	xmlEAAPayload.BirthMiddleName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthMiddleName())
	xmlEAAPayload.Salutation = b.xmlDisclosableClaim(eaaPayloadProxy.HolderSalutation())
	xmlEAAPayload.Title = b.xmlDisclosableClaim(eaaPayloadProxy.HolderTitle())
	xmlEAAPayload.MobilePhoneNumber = b.xmlDisclosableClaim(eaaPayloadProxy.HolderMobilePhoneNumber())
	xmlEAAPayload.Pseudonym = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPseudonym())

	xmlEAAPayload.IssuingCountry = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentIssuingAuthorityCountry())
	xmlEAAPayload.IssuingAuthority = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentIssuingAuthority())
	xmlEAAPayload.DocumentNumber = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentNumber())
	xmlEAAPayload.Portrait = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPortrait())
	xmlEAAPayload.DrivingPrivileges = b.xmlDisclosableClaim(asClaimOrNil(eaaPayloadProxy.HolderDrivingPrivileges()))
	xmlEAAPayload.UNDistinguishingSign = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentIssuingAuthorityUNDistinguishingSign())
	xmlEAAPayload.PersonalAdministrativeNumber = b.xmlDisclosableClaim(eaaPayloadProxy.PersonalAdministrativeNumber())
	xmlEAAPayload.Height = b.xmlDisclosableClaim(eaaPayloadProxy.HolderHeight())
	xmlEAAPayload.Weight = b.xmlDisclosableClaim(eaaPayloadProxy.HolderWeight())
	xmlEAAPayload.EyeColour = b.xmlDisclosableClaim(eaaPayloadProxy.HolderEyeColour())
	xmlEAAPayload.HairColour = b.xmlDisclosableClaim(eaaPayloadProxy.HolderHairColour())
	xmlEAAPayload.ResidentPostalAddress = b.xmlDisclosableClaim(eaaPayloadProxy.ResidentPostalAddress())
	xmlEAAPayload.PortraitCaptureDate = b.xmlDisclosableClaim(eaaPayloadProxy.HolderPortraitCaptureDate())
	xmlEAAPayload.AgeInYears = b.xmlDisclosableClaim(eaaPayloadProxy.HolderAgeInYears())
	xmlEAAPayload.AgeBirthYear = b.xmlDisclosableClaim(eaaPayloadProxy.HolderAgeBirthYear())
	if eaaPayloadProxy.HolderAgeEqualOrOver() != nil {
		xmlEAAPayload.AgeOverNN = append(xmlEAAPayload.AgeOverNN,
			b.xmlAgeOverNNClaims(eaaPayloadProxy.HolderAgeEqualOrOver().AgeEqualOrOverList())...)
	}
	xmlEAAPayload.AgeOverNN = append(xmlEAAPayload.AgeOverNN, b.xmlAgeOverNNClaims(eaaPayloadProxy.HolderAgeOverList())...)
	xmlEAAPayload.IssuingJurisdiction = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentIssuingAuthorityJurisdiction())
	xmlEAAPayload.ResidentAddressCity = b.xmlDisclosableClaim(eaaPayloadProxy.HolderResidentAddressCity())
	xmlEAAPayload.ResidentAddressState = b.xmlDisclosableClaim(eaaPayloadProxy.HolderResidentAddressState())
	xmlEAAPayload.ResidentAddressPostalCode = b.xmlDisclosableClaim(eaaPayloadProxy.HolderResidentAddressPostalCode())
	xmlEAAPayload.ResidentAddressCountry = b.xmlDisclosableClaim(eaaPayloadProxy.HolderResidentAddressCountry())
	xmlEAAPayload.BiometricTemplate = append(xmlEAAPayload.BiometricTemplate,
		b.biometricTemplateXXClaims(eaaPayloadProxy.HolderBiometricTemplateList())...)
	xmlEAAPayload.SignatureUsualMark = b.xmlDisclosableClaim(eaaPayloadProxy.HolderSignatureUsualMark())
	xmlEAAPayload.Fingerprint = b.xmlDisclosableClaim(eaaPayloadProxy.HolderFingerprint())
	xmlEAAPayload.BusinessName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBusinessName())
	xmlEAAPayload.OrganizationName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderOrganizationName())
	xmlEAAPayload.BirthFullName = b.xmlDisclosableClaim(eaaPayloadProxy.HolderBirthFullName())
	xmlEAAPayload.Profession = b.xmlDisclosableClaim(eaaPayloadProxy.HolderProfession())
	xmlEAAPayload.RelationshipFather = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipFather())
	xmlEAAPayload.RelationshipMother = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipMother())
	xmlEAAPayload.RelationshipParent = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipParent())
	xmlEAAPayload.RelationshipSon = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipSon())
	xmlEAAPayload.RelationshipDaughter = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipDaughter())
	xmlEAAPayload.RelationshipBrother = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipBrother())
	xmlEAAPayload.RelationshipSister = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipSister())
	xmlEAAPayload.RelationshipSibling = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipSibling())
	xmlEAAPayload.RelationshipSpouse = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipSpouse())
	xmlEAAPayload.RelationshipFatherInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipFatherInLaw())
	xmlEAAPayload.RelationshipMotherInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipMotherInLaw())
	xmlEAAPayload.RelationshipParentInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipParentInLaw())
	xmlEAAPayload.RelationshipSonInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipSonInLaw())
	xmlEAAPayload.RelationshipDaughterInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipDaughterInLaw())
	xmlEAAPayload.RelationshipChildInLaw = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipChildInLaw())
	xmlEAAPayload.RelationshipParentalAuthority = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipParentalAuthority())
	xmlEAAPayload.RelationshipLegalRepresentative = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipLegalRepresentative())
	xmlEAAPayload.RelationshipAgent = b.xmlDisclosableClaim(eaaPayloadProxy.HolderRelationshipAgent())
	xmlEAAPayload.DocumentType = b.xmlDisclosableClaim(eaaPayloadProxy.DocumentType())

	xmlEAAPayload.IssuingAuthorityRegistrationIdentifier = b.xmlDisclosableClaim(eaaPayloadProxy.IssuingAuthorityRegistrationIdentifier())

	xmlEAAPayload.TrustAnchor = b.xmlDisclosableClaim(eaaPayloadProxy.TrustAnchor())
	xmlEAAPayload.ResidentAddressStreet = b.xmlDisclosableClaim(eaaPayloadProxy.ResidentAddressStreet())
	xmlEAAPayload.ResidentAddressHouseNumber = b.xmlDisclosableClaim(eaaPayloadProxy.ResidentAddressHouseNumber())

	otherClaims := eaaWrapper.OtherClaims()
	if len(otherClaims) > 0 {
		for _, otherClaim := range otherClaims {
			xmlEAAPayload.OtherClaim = append(xmlEAAPayload.OtherClaim, b.xmlDisclosableClaim(otherClaim))
		}
	}

	return xmlEAAPayload
}

// asClaimOrNil unwraps a typed claim wrapper into its ClaimWrapper base,
// preserving the Java null (the typed getters return null, which the Java
// call site passes straight to getXmlDisclosableClaim(ClaimWrapper)).
func asClaimOrNil(wrapper *diagnostic.DrivingPrivilegesClaimWrapper) *diagnostic.ClaimWrapper {
	if wrapper == nil {
		return nil
	}
	return wrapper.AsClaim()
}

// xmlAgeOverNNClaims is the port of the private
// getXmlAgeOverNNClaims(List<AgeOverNNClaimWrapper>).
func (b *SimpleReportBuilder) xmlAgeOverNNClaims(claimWrappers []*diagnostic.AgeOverNNClaimWrapper) []*jaxb.XmlParametrizedDisclosableClaim {
	if len(claimWrappers) == 0 {
		return nil
	}
	result := make([]*jaxb.XmlParametrizedDisclosableClaim, 0, len(claimWrappers))
	for _, claimWrapper := range claimWrappers {
		result = append(result, b.xmlAgeOverNNClaim(claimWrapper))
	}
	return result
}

// xmlAgeOverNNClaim is the port of the private
// getXmlAgeOverNNClaim(AgeOverNNClaimWrapper).
func (b *SimpleReportBuilder) xmlAgeOverNNClaim(ageOverNNClaim *diagnostic.AgeOverNNClaimWrapper) *jaxb.XmlParametrizedDisclosableClaim {
	xmlParametrizedDisclosableClaim := &jaxb.XmlParametrizedDisclosableClaim{}
	name := ageOverNNClaim.Name()
	xmlParametrizedDisclosableClaim.Name = &name
	disclosure := ageOverNNClaim.IsSelectivelyDisclosable()
	xmlParametrizedDisclosableClaim.Disclosure = &disclosure
	parameter := strconv.Itoa(ageOverNNClaim.Age())
	xmlParametrizedDisclosableClaim.Parameter = &parameter
	xmlParametrizedDisclosableClaim.Value = ageOverNNClaim.DisplayValue()
	return xmlParametrizedDisclosableClaim
}

// biometricTemplateXXClaims is the port of the private
// getBiometricTemplateXXClaims(List<BiometricTemplateXXClaimWrapper>).
func (b *SimpleReportBuilder) biometricTemplateXXClaims(
	claimWrappers []*diagnostic.BiometricTemplateXXClaimWrapper) []*jaxb.XmlParametrizedDisclosableClaim {
	if len(claimWrappers) == 0 {
		return nil
	}
	result := make([]*jaxb.XmlParametrizedDisclosableClaim, 0, len(claimWrappers))
	for _, claimWrapper := range claimWrappers {
		result = append(result, b.xmlBiometricTemplateXXClaim(claimWrapper))
	}
	return result
}

// xmlBiometricTemplateXXClaim is the port of the private
// getXmlBiometricTemplateXXClaim(BiometricTemplateXXClaimWrapper).
func (b *SimpleReportBuilder) xmlBiometricTemplateXXClaim(
	biometricTemplateXXClaim *diagnostic.BiometricTemplateXXClaimWrapper) *jaxb.XmlParametrizedDisclosableClaim {
	xmlParametrizedDisclosableClaim := &jaxb.XmlParametrizedDisclosableClaim{}
	name := biometricTemplateXXClaim.Name()
	xmlParametrizedDisclosableClaim.Name = &name
	disclosure := biometricTemplateXXClaim.IsSelectivelyDisclosable()
	xmlParametrizedDisclosableClaim.Disclosure = &disclosure
	parameter := biometricTemplateXXClaim.Type()
	xmlParametrizedDisclosableClaim.Parameter = &parameter
	xmlParametrizedDisclosableClaim.Value = biometricTemplateXXClaim.DisplayValue()
	return xmlParametrizedDisclosableClaim
}

// xmlDisclosableClaim is the port of the private
// getXmlDisclosableClaim(ClaimWrapper).
func (b *SimpleReportBuilder) xmlDisclosableClaim(claimWrapper *diagnostic.ClaimWrapper) *jaxb.XmlDisclosableClaim {
	return b.xmlDisclosableClaimWithDisclosure(claimWrapper, nil)
}

// xmlDisclosableClaimWithDisclosure is the port of the private
// getXmlDisclosableClaim(ClaimWrapper, Boolean).
func (b *SimpleReportBuilder) xmlDisclosableClaimWithDisclosure(claimWrapper *diagnostic.ClaimWrapper,
	selectivelyDisclosable *bool) *jaxb.XmlDisclosableClaim {
	if claimWrapper == nil {
		return nil
	}
	isDisclosure := selectivelyDisclosable
	if isDisclosure == nil {
		value := claimWrapper.IsSelectivelyDisclosable()
		isDisclosure = &value
	}
	name := b.namePath(claimWrapper)
	return b.newXmlDisclosableClaim(name, isDisclosure, claimWrapper.DisplayValue())
}

// newXmlDisclosableClaim is the port of the private
// getXmlDisclosableClaim(String, Boolean, String).
func (b *SimpleReportBuilder) newXmlDisclosableClaim(name *string, selectivelyDisclosable *bool,
	displayValue string) *jaxb.XmlDisclosableClaim {
	xmlDisclosableClaim := &jaxb.XmlDisclosableClaim{}
	xmlDisclosableClaim.Name = name
	if selectivelyDisclosable != nil && *selectivelyDisclosable {
		xmlDisclosableClaim.Disclosure = selectivelyDisclosable
	}
	xmlDisclosableClaim.Value = displayValue
	return xmlDisclosableClaim
}

// namePath is the port of the private getNamePath(ClaimWrapper). Java's
// getName() is a nullable String, so the port keeps a *string: a claim with
// no name and no parent yields nil, which JAXB omits the name attribute for.
func (b *SimpleReportBuilder) namePath(claimWrapper *diagnostic.ClaimWrapper) *string {
	parent := claimWrapper.Parent()
	if parent == nil {
		name := claimWrapper.Name()
		if name == "" {
			return nil
		}
		return &name
	}
	sb := claimWrapper.Name()
	for parent != nil {
		if parent.IsList() {
			position := indexOfClaim(parent.Wrapped().Entry, claimWrapper.Wrapped())
			sb = strconv.Itoa(position) + "/" + sb
		}
		parentName := parent.Name()
		if parentName != "" {
			sb = parentName + "/" + sb
		}
		parent = parent.Parent()
	}
	return &sb
}

// indexOfClaim is the port of List.indexOf(Object) over the parent's entry
// list. The generated XmlClaim has no value-equality contract, so the search
// is by pointer identity, which is what the Java call site relies on (the
// wrapped claims are the very objects held by the parent's list).
func indexOfClaim(entries []*diagnosticjaxb.XmlClaim, claim *diagnosticjaxb.XmlClaim) int {
	for i, entry := range entries {
		if entry == claim {
			return i
		}
	}
	return -1
}

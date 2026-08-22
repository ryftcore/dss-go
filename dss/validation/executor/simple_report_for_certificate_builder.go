// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/certificate/SimpleReportForCertificateBuilder.java
// (DSS 6.5.RC1).
//
// Deviations:
//
//   - getUniqueServiceNames() returns a HashSet<String>, so the order in which
//     Java appends the <trustAnchor> elements is unspecified (and unstable
//     across JVM runs for more than one name). This port keeps the names in
//     first-occurrence order, which is deterministic and identical to Java's
//     for the single-name case the corpus exercises.
//
//   - Java's null-versus-empty distinctions are byte-visible and preserved
//     verbatim: emptyToNull() nulls out an empty URL list (the wrapper element
//     is then OMITTED, not written empty), getReadable() returns null for an
//     empty OID list, setPdsUrls(null) is an explicit null, and an
//     XmlRevocation is always set even when there is no revocation data.

package executor

import (
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
)

// SimpleReportForCertificateBuilder builds a SimpleReport for a certificate
// validation. Port of SimpleReportForCertificateBuilder.
type SimpleReportForCertificateBuilder struct {
	// diagnosticData is the diagnostic data. Port of the private final
	// diagnosticData field.
	diagnosticData *diagnostic.DiagnosticData

	// detailedReport is the detailed report. Port of the private final
	// detailedReport field.
	detailedReport *detailedreport.DetailedReport

	// policy is the validation policy. Port of the private final policy field.
	policy policy.ValidationPolicy

	// currentTime is the validation time. Port of the private final
	// currentTime field.
	currentTime time.Time

	// certificateId is the id of the certificate to be validated. Port of the
	// private final certificateId field.
	certificateId string
}

// NewSimpleReportForCertificateBuilder is the default constructor. Port of
// SimpleReportForCertificateBuilder(DiagnosticData, DetailedReport,
// ValidationPolicy, Date, String).
func NewSimpleReportForCertificateBuilder(diagnosticData *diagnostic.DiagnosticData,
	detailedReport *detailedreport.DetailedReport, validationPolicy policy.ValidationPolicy,
	currentTime time.Time, certificateId string) *SimpleReportForCertificateBuilder {
	return &SimpleReportForCertificateBuilder{
		diagnosticData: diagnosticData,
		detailedReport: detailedReport,
		policy:         validationPolicy,
		currentTime:    currentTime,
		certificateId:  certificateId,
	}
}

// Build builds the XmlSimpleCertificateReport. Port of build().
func (b *SimpleReportForCertificateBuilder) Build() *jaxb.XmlSimpleCertificateReport {
	simpleReport := &jaxb.XmlSimpleCertificateReport{}

	b.addPolicyNode(simpleReport)
	b.addValidationTime(simpleReport)

	simpleReport.ValidationTime = jaxb.NewXSDateTime(b.currentTime)

	certificate := b.diagnosticData.UsedCertificateById(b.certificateId)
	targetCertificate := b.chainItem(certificate, false)
	b.addQualifications(targetCertificate, certificate)
	b.addQWACValidationDetails(targetCertificate, certificate)
	b.addCertificateApprovalStatuss(targetCertificate, certificate)
	simpleReport.Certificate = targetCertificate

	chain := make([]*jaxb.XmlChainItem, 0)
	trustAnchorReached := false
	certificateChain := certificate.CertificateChain()
	for _, cert := range certificateChain {
		trustAnchorReached = trustAnchorReached || cert.IsTrusted()
		chain = append(chain, b.chainItem(cert, trustAnchorReached))
	}
	targetCertificate.Chain = chain

	b.addConnectionDetails(simpleReport)

	return simpleReport
}

// addPolicyNode is the port of the private addPolicyNode(XmlSimpleCertificateReport).
func (b *SimpleReportForCertificateBuilder) addPolicyNode(report *jaxb.XmlSimpleCertificateReport) {
	xmlPolicy := &jaxb.XmlValidationPolicy{}
	policyName := b.policy.PolicyName()
	xmlPolicy.PolicyName = &policyName
	policyDescription := b.policy.PolicyDescription()
	xmlPolicy.PolicyDescription = &policyDescription
	report.ValidationPolicy = xmlPolicy
}

// addValidationTime is the port of the private
// addValidationTime(XmlSimpleCertificateReport).
func (b *SimpleReportForCertificateBuilder) addValidationTime(report *jaxb.XmlSimpleCertificateReport) {
	report.ValidationTime = jaxb.NewXSDateTime(b.currentTime)
}

// chainItem is the port of the private getChainItem(CertificateWrapper, boolean).
func (b *SimpleReportForCertificateBuilder) chainItem(certificate *diagnostic.CertificateWrapper,
	trustAnchorReached bool) *jaxb.XmlChainItem {
	item := &jaxb.XmlChainItem{}
	item.Id = certificate.Id()
	item.Subject = b.subject(certificate)
	signingCertificate := certificate.SigningCertificate()
	if signingCertificate != nil {
		issuerId := signingCertificate.Id()
		item.IssuerId = &issuerId
	}
	if notBefore := certificate.NotBefore(); notBefore != nil {
		item.NotBefore = jaxb.NewXSDateTime(*notBefore)
	}
	if notAfter := certificate.NotAfter(); notAfter != nil {
		item.NotAfter = jaxb.NewXSDateTime(*notAfter)
	}
	// Java calls setKeyUsages unconditionally with the (possibly empty) list
	// getKeyUsages() returns, so the @XmlElementWrapper is always written -
	// as <keyUsages/> when there is nothing in it. Unlike the URL wrappers
	// below, which emptyToNull() nulls out.
	keyUsages := certificate.KeyUsages()
	values := make([]jaxb.KeyUsageBitValue, 0, len(keyUsages))
	for _, keyUsage := range keyUsages {
		values = append(values, jaxb.KeyUsageBitValue(keyUsage))
	}
	item.KeyUsages = &jaxb.XmlKeyUsages{KeyUsage: values}
	if extendedKeyUsages := readableOIDs(certificate.ExtendedKeyUsages()); extendedKeyUsages != nil {
		item.ExtendedKeyUsages = &jaxb.XmlExtendedKeyUsages{ExtendedKeyUsage: extendedKeyUsages}
	}
	if aiaUrls := emptyToNilURLs(certificate.CAIssuersAccessUrls()); aiaUrls != nil {
		item.AiaUrls = &jaxb.XmlAiaUrls{AiaUrl: aiaUrls}
	}
	if ocspUrls := emptyToNilURLs(certificate.OCSPAccessUrls()); ocspUrls != nil {
		item.OcspUrls = &jaxb.XmlOcspUrls{OcspUrl: ocspUrls}
	}
	if crlUrls := emptyToNilURLs(certificate.CRLDistributionPoints()); crlUrls != nil {
		item.CrlUrls = &jaxb.XmlCrlUrls{CrlUrl: crlUrls}
	}
	if cpsUrls := emptyToNilURLs(certificate.CpsUrls()); cpsUrls != nil {
		item.CpsUrls = &jaxb.XmlCpsUrls{CpsUrl: cpsUrls}
	}
	// Java: item.setPdsUrls(null)
	item.PdsUrls = nil

	revocation := &jaxb.XmlRevocation{}
	revocationData := b.diagnosticData.LatestRevocationDataForCertificate(certificate)
	if revocationData != nil {
		if thisUpdate := revocationData.ThisUpdate(); thisUpdate != nil {
			revocation.ThisUpdate = jaxb.NewXSDateTime(*thisUpdate)
		}
		if revocationDate := revocationData.RevocationDate(); revocationDate != nil {
			revocation.RevocationDate = jaxb.NewXSDateTime(*revocationDate)
		}
		if reason := revocationData.Reason(); reason != "" {
			revocationReason := jaxb.RevocationReasonValue(reason)
			revocation.RevocationReason = &revocationReason
		}
	}
	item.Revocation = revocation

	if certificate.IsTrusted() {
		trustAnchors := make([]*jaxb.XmlTrustAnchor, 0)

		trustServiceProviders := filterTSPsByCertificateId(certificate.TrustServiceProviders(), certificate.Id())
		for _, xmlTrustServiceProvider := range trustServiceProviders {
			var trustServices []*diagnosticjaxb.XmlTrustService
			if xmlTrustServiceProvider.TrustServices != nil {
				trustServices = xmlTrustServiceProvider.TrustServices.Items
			}
			uniqueServiceNames := uniqueTrustServiceNames(trustServices)
			for _, serviceName := range uniqueServiceNames {
				trustAnchor := &jaxb.XmlTrustAnchor{}
				if xmlTrustServiceProvider.TL != nil {
					if countryCode := xmlTrustServiceProvider.TL.CountryCode; countryCode != nil {
						trustAnchor.CountryCode = *countryCode
					}
					trustAnchor.TslType = xmlTrustServiceProvider.TL.Type
				}
				if tspName := enOrFirstLangAndValue(tspNames(xmlTrustServiceProvider)); tspName != nil {
					trustAnchor.TrustServiceProvider = *tspName
				}
				tspRegistrationIdentifiers := tspRegistrationIdentifiers(xmlTrustServiceProvider)
				if len(tspRegistrationIdentifiers) > 0 {
					registrationId := tspRegistrationIdentifiers[0]
					trustAnchor.TrustServiceProviderRegistrationId = &registrationId
				}
				trustAnchor.TrustServiceName = serviceName
				trustAnchors = append(trustAnchors, trustAnchor)
			}
		}
		// NOTE: separate ?
		trustedEntities := filterTEsByCertificateId(certificate.TrustedEntities(), certificate.Id())
		for _, xmlTrustedEntity := range trustedEntities {
			var trustedEntityServices []*diagnosticjaxb.XmlTrustedEntityService
			if xmlTrustedEntity.TrustedEntityServices != nil {
				trustedEntityServices = xmlTrustedEntity.TrustedEntityServices.Items
			}
			uniqueServiceNames := uniqueTrustedEntityServiceNames(trustedEntityServices)
			for _, serviceName := range uniqueServiceNames {
				trustAnchor := &jaxb.XmlTrustAnchor{}
				if xmlTrustedEntity.LoTE != nil {
					if countryCode := xmlTrustedEntity.LoTE.CountryCode; countryCode != nil {
						trustAnchor.CountryCode = *countryCode
					}
					trustAnchor.TslType = xmlTrustedEntity.LoTE.Type
				}
				if name := enOrFirstLangAndValue(trustedEntityNames(xmlTrustedEntity)); name != nil {
					trustAnchor.TrustServiceProvider = *name
				}
				registrationIdentifiers := trustedEntityRegistrationIdentifiers(xmlTrustedEntity)
				if len(registrationIdentifiers) > 0 {
					registrationId := registrationIdentifiers[0]
					trustAnchor.TrustServiceProviderRegistrationId = &registrationId
				}
				trustAnchor.TrustServiceName = serviceName
				trustAnchors = append(trustAnchors, trustAnchor)
			}
		}

		item.TrustAnchors = &jaxb.XmlTrustAnchors{TrustAnchor: trustAnchors}
		if trustStartDate := certificate.TrustStartDate(); trustStartDate != nil {
			item.TrustStartDate = jaxb.NewXSDateTime(*trustStartDate)
		}
		if trustSunsetDate := certificate.TrustSunsetDate(); trustSunsetDate != nil {
			item.TrustSunsetDate = jaxb.NewXSDateTime(*trustSunsetDate)
		}

	} else {
		item.TrustAnchors = nil
	}

	conclusion := b.detailedReport.CertificateXCVConclusion(certificate.Id())
	if conclusion != nil {
		item.Indication = jaxb.IndicationValue(conclusion.Indication.Indication())
		if conclusion.SubIndication != nil {
			subIndication := jaxb.SubIndicationValue(conclusion.SubIndication.SubIndication())
			item.SubIndication = &subIndication
		}

	} else if certificate.IsTrusted() || trustAnchorReached {
		item.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)

	} else {
		// if certificate was not validated or not trusted
		item.Indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
		subIndication := jaxb.SubIndicationValue(enumerations.SubIndicationNoCertificateChainFound)
		item.SubIndication = &subIndication
	}

	validationDetails := b.validationDetails(certificate.Id())
	if detailsNotEmpty(validationDetails) {
		item.X509ValidationDetails = validationDetails
	}

	return item
}

// enOrFirstLangAndValue is the port of the private
// getEnOrFirst(List<XmlLangAndValue>).
func enOrFirstLangAndValue(langAndValues []*diagnosticjaxb.XmlLangAndValue) *string {
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

// tspNames unwraps XmlTrustServiceProvider.getTSPNames().
func tspNames(xmlTrustServiceProvider *diagnosticjaxb.XmlTrustServiceProvider) []*diagnosticjaxb.XmlLangAndValue {
	if xmlTrustServiceProvider.TSPNames == nil {
		return nil
	}
	return xmlTrustServiceProvider.TSPNames.Items
}

// tspRegistrationIdentifiers unwraps
// XmlTrustServiceProvider.getTSPRegistrationIdentifiers().
func tspRegistrationIdentifiers(xmlTrustServiceProvider *diagnosticjaxb.XmlTrustServiceProvider) []string {
	if xmlTrustServiceProvider.TSPRegistrationIdentifiers == nil {
		return nil
	}
	return xmlTrustServiceProvider.TSPRegistrationIdentifiers.Items
}

// trustedEntityNames unwraps XmlTrustedEntity.getNames().
func trustedEntityNames(xmlTrustedEntity *diagnosticjaxb.XmlTrustedEntity) []*diagnosticjaxb.XmlLangAndValue {
	if xmlTrustedEntity.Names == nil {
		return nil
	}
	return xmlTrustedEntity.Names.Items
}

// trustedEntityRegistrationIdentifiers unwraps
// XmlTrustedEntity.getRegistrationIdentifiers().
func trustedEntityRegistrationIdentifiers(xmlTrustedEntity *diagnosticjaxb.XmlTrustedEntity) []string {
	if xmlTrustedEntity.RegistrationIdentifiers == nil {
		return nil
	}
	return xmlTrustedEntity.RegistrationIdentifiers.Items
}

// filterTSPsByCertificateId is the port of the private
// filterTSPsByCertificateId(List<XmlTrustServiceProvider>, String).
func filterTSPsByCertificateId(trustServiceProviders []*diagnosticjaxb.XmlTrustServiceProvider,
	certificateId string) []*diagnosticjaxb.XmlTrustServiceProvider {
	result := make([]*diagnosticjaxb.XmlTrustServiceProvider, 0)
	for _, xmlTrustServiceProvider := range trustServiceProviders {
		var trustServices []*diagnosticjaxb.XmlTrustService
		if xmlTrustServiceProvider.TrustServices != nil {
			trustServices = xmlTrustServiceProvider.TrustServices.Items
		}
		foundCertId := false
		for _, xmlTrustService := range trustServices {
			if xmlTrustService.ServiceDigitalIdentifier != nil &&
				certificateId == serviceDigitalIdentifierId(xmlTrustService.ServiceDigitalIdentifier) {
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

// filterTEsByCertificateId is the port of the private
// filterTEsByCertificateId(List<XmlTrustedEntity>, String).
func filterTEsByCertificateId(trustedEntities []*diagnosticjaxb.XmlTrustedEntity,
	certificateId string) []*diagnosticjaxb.XmlTrustedEntity {
	result := make([]*diagnosticjaxb.XmlTrustedEntity, 0)
	for _, xmlTrustedEntity := range trustedEntities {
		var trustedEntityServices []*diagnosticjaxb.XmlTrustedEntityService
		if xmlTrustedEntity.TrustedEntityServices != nil {
			trustedEntityServices = xmlTrustedEntity.TrustedEntityServices.Items
		}
		foundCertId := false
		for _, xmlTrustedService := range trustedEntityServices {
			if xmlTrustedService.ServiceDigitalIdentifier != nil &&
				certificateId == serviceDigitalIdentifierId(xmlTrustedService.ServiceDigitalIdentifier) {
				foundCertId = true
				break
			}
		}
		if foundCertId {
			result = append(result, xmlTrustedEntity)
		}
	}
	return result
}

// serviceDigitalIdentifierId reads the Id of the certificate referenced by a
// trust service's ServiceDigitalIdentifier (Java's
// xmlTrustService.getServiceDigitalIdentifier().getId()).
func serviceDigitalIdentifierId(certificate *diagnosticjaxb.XmlCertificate) string {
	if certificate.Id == nil {
		return ""
	}
	return string(*certificate.Id)
}

// readableOIDs is the port of the private getReadable(List<XmlOID>).
func readableOIDs(oids []*diagnosticjaxb.XmlOID) []string {
	if len(oids) > 0 {
		result := make([]string, 0, len(oids))
		for _, xmlOID := range oids {
			if xmlOID.Description != nil && *xmlOID.Description != "" {
				result = append(result, *xmlOID.Description)
			} else {
				result = append(result, xmlOID.Value)
			}
		}
		return result
	}
	return nil
}

// uniqueTrustServiceNames is the port of the private
// getUniqueServiceNames(List<? extends XmlTrustedEntityService>) applied to
// XmlTrustService. See the file header on the ordering deviation.
func uniqueTrustServiceNames(trustServices []*diagnosticjaxb.XmlTrustService) []string {
	result := make([]string, 0, len(trustServices))
	seen := make(map[string]struct{}, len(trustServices))
	for _, xmlTrustService := range trustServices {
		var names []*diagnosticjaxb.XmlLangAndValue
		if xmlTrustService.ServiceNames != nil {
			names = xmlTrustService.ServiceNames.Items
		}
		name := ""
		if value := enOrFirstLangAndValue(names); value != nil {
			name = *value
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	// Java returns a HashSet<String>, whose iteration order reaches the
	// marshalled <trustServiceName> sequence; see JavaHashSetStringOrder.
	return JavaHashSetStringOrder(result)
}

// uniqueTrustedEntityServiceNames is the port of the private
// getUniqueServiceNames(List<? extends XmlTrustedEntityService>).
func uniqueTrustedEntityServiceNames(trustedEntityServices []*diagnosticjaxb.XmlTrustedEntityService) []string {
	result := make([]string, 0, len(trustedEntityServices))
	seen := make(map[string]struct{}, len(trustedEntityServices))
	for _, xmlTrustService := range trustedEntityServices {
		var names []*diagnosticjaxb.XmlLangAndValue
		if xmlTrustService.ServiceNames != nil {
			names = xmlTrustService.ServiceNames.Items
		}
		name := ""
		if value := enOrFirstLangAndValue(names); value != nil {
			name = *value
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return JavaHashSetStringOrder(result)
}

// subject is the port of the private getSubject(CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) subject(certificate *diagnostic.CertificateWrapper) *jaxb.XmlSubject {
	subject := &jaxb.XmlSubject{}
	subject.CommonName = nilIfEmpty(certificate.CommonName())
	subject.Pseudonym = nilIfEmpty(certificate.Pseudo())
	subject.Surname = nilIfEmpty(certificate.Surname())
	subject.GivenName = nilIfEmpty(certificate.GivenName())
	subject.OrganizationName = nilIfEmpty(certificate.OrganizationName())
	subject.OrganizationUnit = nilIfEmpty(certificate.OrganizationalUnit())
	subject.Email = nilIfEmpty(certificate.Email())
	subject.Locality = nilIfEmpty(certificate.Locality())
	subject.State = nilIfEmpty(certificate.State())
	subject.Country = nilIfEmpty(certificate.CountryName())
	return subject
}

// nilIfEmpty maps the Go zero string, which the diagnostic wrappers return
// where Java returns null, back onto a null field.
func nilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// emptyToNilURLs is the port of the private emptyToNull(List<String>).
func emptyToNilURLs(listUrls []string) []string {
	if len(listUrls) == 0 {
		return nil
	}
	return listUrls
}

// addQualifications is the port of the private
// addQualifications(XmlChainItem, CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) addQualifications(chainItem *jaxb.XmlChainItem,
	certificate *diagnostic.CertificateWrapper) {
	qualificationAtIssuance := jaxb.CertificateQualificationValue(
		b.detailedReport.CertificateQualificationAtIssuance(certificate.Id()))
	chainItem.QualificationAtIssuance = &qualificationAtIssuance
	qualificationAtValidation := jaxb.CertificateQualificationValue(
		b.detailedReport.CertificateQualificationAtValidation(certificate.Id()))
	chainItem.QualificationAtValidation = &qualificationAtValidation

	qualificationDetailsAtIssuanceTime := b.certificateQualificationDetailsAtIssuanceTime(certificate.Id())
	if detailsNotEmpty(qualificationDetailsAtIssuanceTime) {
		chainItem.QualificationDetailsAtIssuance = qualificationDetailsAtIssuanceTime
	}
	qualificationDetailsAtValidationTime := b.certificateQualificationDetailsAtValidationTime(certificate.Id())
	if detailsNotEmpty(qualificationDetailsAtValidationTime) {
		chainItem.QualificationDetailsAtValidation = qualificationDetailsAtValidationTime
	}

	var enactedMRA *bool
	trustServices := certificate.TrustServices()
	for _, trustServiceWrapper := range trustServices {
		if trustServiceWrapper.IsEnactedMRA() {
			enacted := true
			enactedMRA = &enacted
			break
		}
	}
	chainItem.EnactedMRA = enactedMRA
}

// addQWACValidationDetails is the port of the private
// addQWACValidationDetails(XmlChainItem, CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) addQWACValidationDetails(chainItem *jaxb.XmlChainItem,
	certificate *diagnostic.CertificateWrapper) {
	b.addQWACProfile(chainItem, certificate)
	b.addTLSBindingSignature(chainItem)
}

// addQWACProfile is the port of the private
// addQWACProfile(XmlChainItem, CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) addQWACProfile(chainItem *jaxb.XmlChainItem,
	certificate *diagnostic.CertificateWrapper) {
	// Java's getCertificateQWACProfile returns null when the certificate has no
	// QWAC process; the ported DetailedReport returns the empty QWACProfile for
	// that null, which must stay an OMITTED element, not an empty one.
	if profile := b.detailedReport.CertificateQWACProfile(certificate.Id()); profile != "" {
		qwacProfile := jaxb.QWACProfileValue(profile)
		chainItem.QwacProfile = &qwacProfile
	}
	qwacValidationDetails := b.qwacValidationDetails(certificate.Id())
	if detailsNotEmpty(qwacValidationDetails) {
		chainItem.QwacDetails = qwacValidationDetails
	}
}

// addCertificateApprovalStatuss is the port of the private
// addCertificateApprovalStatuss(XmlChainItem, CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) addCertificateApprovalStatuss(chainItem *jaxb.XmlChainItem,
	certificate *diagnostic.CertificateWrapper) {
	chainItem.CertificateApprovalStatusAtIssuanceTime = b.certificateApprovalStatusAtIssuanceTime(certificate)
	chainItem.CertificateApprovalStatusAtValidationTime = b.certificateApprovalStatusAtValidationTime(certificate)
}

// certificateApprovalStatusAtIssuanceTime is the port of the private
// getCertificateApprovalStatusAtIssuanceTime(CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) certificateApprovalStatusAtIssuanceTime(
	certificate *diagnostic.CertificateWrapper) *jaxb.XmlCertificateApprovalStatusAtIssuanceTime {
	certificateApprovalStatussAtIssuanceTime := b.detailedReport.CertificateApprovalStatussAtIssuanceTime(certificate.Id())
	if len(certificateApprovalStatussAtIssuanceTime) == 0 {
		return nil
	}

	xmlCertificateApprovalStatusAtTime := &jaxb.XmlCertificateApprovalStatusAtIssuanceTime{}
	for _, certificateApprovalStatus := range certificateApprovalStatussAtIssuanceTime {
		xmlCertificateApprovalStatus := &jaxb.XmlCertificateApprovalStatus{}
		if listType := certificateApprovalStatus.ListType(); listType != nil {
			xmlCertificateApprovalStatus.ListType = &jaxb.ListTypeValue{Value: listType}
		}
		if sti := certificateApprovalStatus.ServiceTypeIdentifier(); sti != nil {
			xmlCertificateApprovalStatus.ServiceTypeIdentifier = &jaxb.LoTEServiceTypeIdentifierValue{Value: sti}
		}
		if status := certificateApprovalStatus.ServiceStatus(); status != nil {
			xmlCertificateApprovalStatus.ServiceStatus = &jaxb.LoTEServiceStatusValue{Value: status}
		}
		xmlCertificateApprovalStatus.Label = nilIfEmpty(certificateApprovalStatus.Label())
		xmlCertificateApprovalStatus.Details = b.certificateApprovalStatusDetailsAtIssuanceTime(
			certificate.Id(), certificateApprovalStatus)

		xmlCertificateApprovalStatusAtTime.CertificateApprovalStatus =
			append(xmlCertificateApprovalStatusAtTime.CertificateApprovalStatus, xmlCertificateApprovalStatus)
	}
	return xmlCertificateApprovalStatusAtTime
}

// certificateApprovalStatusAtValidationTime is the port of the private
// getCertificateApprovalStatusAtValidationTime(CertificateWrapper).
func (b *SimpleReportForCertificateBuilder) certificateApprovalStatusAtValidationTime(
	certificate *diagnostic.CertificateWrapper) *jaxb.XmlCertificateApprovalStatusAtValidationTime {
	certificateApprovalStatussAtIssuanceTime := b.detailedReport.CertificateApprovalStatussAtValidationTime(certificate.Id())
	if len(certificateApprovalStatussAtIssuanceTime) == 0 {
		return nil
	}

	xmlCertificateApprovalStatusAtTime := &jaxb.XmlCertificateApprovalStatusAtValidationTime{}
	for _, certificateApprovalStatus := range certificateApprovalStatussAtIssuanceTime {
		xmlCertificateApprovalStatus := &jaxb.XmlCertificateApprovalStatus{}
		if listType := certificateApprovalStatus.ListType(); listType != nil {
			xmlCertificateApprovalStatus.ListType = &jaxb.ListTypeValue{Value: listType}
		}
		if sti := certificateApprovalStatus.ServiceTypeIdentifier(); sti != nil {
			xmlCertificateApprovalStatus.ServiceTypeIdentifier = &jaxb.LoTEServiceTypeIdentifierValue{Value: sti}
		}
		if status := certificateApprovalStatus.ServiceStatus(); status != nil {
			xmlCertificateApprovalStatus.ServiceStatus = &jaxb.LoTEServiceStatusValue{Value: status}
		}
		xmlCertificateApprovalStatus.Label = nilIfEmpty(certificateApprovalStatus.Label())
		xmlCertificateApprovalStatus.Details = b.certificateApprovalStatusDetailsAtValidationTime(
			certificate.Id(), certificateApprovalStatus)

		xmlCertificateApprovalStatusAtTime.CertificateApprovalStatus =
			append(xmlCertificateApprovalStatusAtTime.CertificateApprovalStatus, xmlCertificateApprovalStatus)
	}
	return xmlCertificateApprovalStatusAtTime
}

// validationDetails is the port of the private getValidationDetails(String).
func (b *SimpleReportForCertificateBuilder) validationDetails(tokenId string) *jaxb.XmlDetails {
	validationDetails := &jaxb.XmlDetails{}
	validationDetails.Error = append(validationDetails.Error, convertMessages(b.detailedReport.AdESValidationErrors(tokenId))...)
	validationDetails.Warning = append(validationDetails.Warning, convertMessages(b.detailedReport.AdESValidationWarnings(tokenId))...)
	validationDetails.Info = append(validationDetails.Info, convertMessages(b.detailedReport.AdESValidationInfos(tokenId))...)
	return validationDetails
}

// certificateQualificationDetailsAtIssuanceTime is the port of the private
// getCertificateQualificationDetailsAtIssuanceTime(String).
func (b *SimpleReportForCertificateBuilder) certificateQualificationDetailsAtIssuanceTime(tokenId string) *jaxb.XmlDetails {
	qualificationDetails := &jaxb.XmlDetails{}
	qualificationDetails.Error = append(qualificationDetails.Error, convertMessages(b.detailedReport.CertificateQualificationErrorsAtIssuanceTime(tokenId))...)
	qualificationDetails.Warning = append(qualificationDetails.Warning, convertMessages(b.detailedReport.CertificateQualificationWarningsAtIssuanceTime(tokenId))...)
	qualificationDetails.Info = append(qualificationDetails.Info, convertMessages(b.detailedReport.CertificateQualificationInfosAtIssuanceTime(tokenId))...)
	return qualificationDetails
}

// certificateQualificationDetailsAtValidationTime is the port of the private
// getCertificateQualificationDetailsAtValidationTime(String).
func (b *SimpleReportForCertificateBuilder) certificateQualificationDetailsAtValidationTime(tokenId string) *jaxb.XmlDetails {
	qualificationDetails := &jaxb.XmlDetails{}
	qualificationDetails.Error = append(qualificationDetails.Error, convertMessages(b.detailedReport.CertificateQualificationErrorsAtValidationTime(tokenId))...)
	qualificationDetails.Warning = append(qualificationDetails.Warning, convertMessages(b.detailedReport.CertificateQualificationWarningsAtValidationTime(tokenId))...)
	qualificationDetails.Info = append(qualificationDetails.Info, convertMessages(b.detailedReport.CertificateQualificationInfosAtValidationTime(tokenId))...)
	return qualificationDetails
}

// qwacValidationDetails is the port of the private
// getQWACValidationDetails(String).
func (b *SimpleReportForCertificateBuilder) qwacValidationDetails(tokenId string) *jaxb.XmlDetails {
	qualificationDetails := &jaxb.XmlDetails{}
	qualificationDetails.Error = append(qualificationDetails.Error, convertMessages(b.detailedReport.QWACValidationErrors(tokenId))...)
	qualificationDetails.Warning = append(qualificationDetails.Warning, convertMessages(b.detailedReport.QWACValidationWarnings(tokenId))...)
	qualificationDetails.Info = append(qualificationDetails.Info, convertMessages(b.detailedReport.QWACValidationInfos(tokenId))...)
	return qualificationDetails
}

// certificateApprovalStatusDetailsAtIssuanceTime is the port of the private
// getCertificateApprovalStatusDetailsAtIssuanceTime(String, CertificateApprovalStatus).
func (b *SimpleReportForCertificateBuilder) certificateApprovalStatusDetailsAtIssuanceTime(tokenId string,
	certificateApprovalStatus enumerations.CertificateApprovalStatus) *jaxb.XmlDetails {
	usageDetails := &jaxb.XmlDetails{}
	usageDetails.Error = append(usageDetails.Error, convertMessages(b.detailedReport.CertificateApprovalStatusErrorsAtIssuanceTime(tokenId, certificateApprovalStatus))...)
	usageDetails.Warning = append(usageDetails.Warning, convertMessages(b.detailedReport.CertificateApprovalStatusWarningsAtIssuanceTime(tokenId, certificateApprovalStatus))...)
	usageDetails.Info = append(usageDetails.Info, convertMessages(b.detailedReport.CertificateApprovalStatusInfosAtIssuanceTime(tokenId, certificateApprovalStatus))...)
	return usageDetails
}

// certificateApprovalStatusDetailsAtValidationTime is the port of the private
// getCertificateApprovalStatusDetailsAtValidationTime(String, CertificateApprovalStatus).
func (b *SimpleReportForCertificateBuilder) certificateApprovalStatusDetailsAtValidationTime(tokenId string,
	certificateApprovalStatus enumerations.CertificateApprovalStatus) *jaxb.XmlDetails {
	usageDetails := &jaxb.XmlDetails{}
	usageDetails.Error = append(usageDetails.Error, convertMessages(b.detailedReport.CertificateApprovalStatusErrorsAtValidationTime(tokenId, certificateApprovalStatus))...)
	usageDetails.Warning = append(usageDetails.Warning, convertMessages(b.detailedReport.CertificateApprovalStatusWarningsAtValidationTime(tokenId, certificateApprovalStatus))...)
	usageDetails.Info = append(usageDetails.Info, convertMessages(b.detailedReport.CertificateApprovalStatusInfosAtValidationTime(tokenId, certificateApprovalStatus))...)
	return usageDetails
}

// addTLSBindingSignature is the port of the private
// addTLSBindingSignature(XmlChainItem).
func (b *SimpleReportForCertificateBuilder) addTLSBindingSignature(chainItem *jaxb.XmlChainItem) {
	bindingSignature := b.diagnosticData.TLSCertificateBindingSignature()
	if bindingSignature != nil {
		xmlSignature := &jaxb.XmlSignature{}
		xmlSignature.Id = bindingSignature.Id()
		xmlSignature.Url = nilIfEmpty(b.diagnosticData.TLSCertificateBindingUrl())
		if signingTime := bindingSignature.ClaimedSigningTime(); signingTime != nil {
			xmlSignature.SigningTime = jaxb.NewXSDateTime(*signingTime)
		}
		xmlSignature.SignatureFormat = jaxb.SignatureLevelValue(bindingSignature.SignatureFormat())
		b.addSignatureScope(bindingSignature, xmlSignature)
		b.addFinalIndication(bindingSignature, xmlSignature)
		b.addAdESValidationDetails(bindingSignature, xmlSignature)
		b.addChain(bindingSignature, xmlSignature)

		chainItem.TLSBindingSignature = xmlSignature
	}
}

// addSignatureScope is the port of the private
// addSignatureScope(SignatureWrapper, XmlSignature).
func (b *SimpleReportForCertificateBuilder) addSignatureScope(signature *diagnostic.SignatureWrapper,
	xmlSignature *jaxb.XmlSignature) {
	signatureScopes := signature.SignatureScopes()
	if len(signatureScopes) > 0 {
		for _, signatureScope := range signatureScopes {
			xmlSignature.SignatureScope = append(xmlSignature.SignatureScope, xmlSignatureScope(signatureScope))
		}
	}
}

// addFinalIndication is the port of the private
// addFinalIndication(SignatureWrapper, XmlSignature).
func (b *SimpleReportForCertificateBuilder) addFinalIndication(signature *diagnostic.SignatureWrapper,
	xmlSignature *jaxb.XmlSignature) {
	xmlSignature.Indication = jaxb.IndicationValue(b.detailedReport.FinalIndication(signature.Id()))
	subIndication := b.detailedReport.FinalSubIndication(signature.Id())
	if subIndication != "" {
		value := jaxb.SubIndicationValue(subIndication)
		xmlSignature.SubIndication = &value
	}
}

// addAdESValidationDetails is the port of the private
// addAdESValidationDetails(SignatureWrapper, XmlSignature).
func (b *SimpleReportForCertificateBuilder) addAdESValidationDetails(signature *diagnostic.SignatureWrapper,
	xmlSignature *jaxb.XmlSignature) {
	validationDetails := b.validationDetails(signature.Id())
	if detailsNotEmpty(validationDetails) {
		xmlSignature.AdESValidationDetails = validationDetails
	}
}

// xmlSignatureScope is the port of the private
// getXmlSignatureScope(eu.europa.esig.dss.diagnostic.jaxb.XmlSignatureScope).
func xmlSignatureScope(signatureScope *diagnosticjaxb.XmlSignatureScope) *jaxb.XmlSignatureScope {
	scope := &jaxb.XmlSignatureScope{}
	if signatureScope.SignerData != nil && signatureScope.SignerData.Id != nil {
		scope.Id = string(*signatureScope.SignerData.Id)
	}
	scope.Name = signatureScope.Name
	if signatureScope.Scope != nil {
		value := jaxb.SignatureScopeTypeValue(signatureScope.Scope.SignatureScopeType())
		scope.Scope = &value
	}
	if signatureScope.Description != nil {
		scope.Value = *signatureScope.Description
	}
	return scope
}

// addChain is the port of the private addChain(SignatureWrapper, XmlSignature).
func (b *SimpleReportForCertificateBuilder) addChain(bindingSignature *diagnostic.SignatureWrapper,
	xmlSignature *jaxb.XmlSignature) {
	certificateChain := bindingSignature.CertificateChain()
	if len(certificateChain) == 0 {
		return
	}

	trustAnchorReached := false
	chain := make([]*jaxb.XmlChainItem, 0)
	for _, cert := range certificateChain {
		trustAnchorReached = trustAnchorReached || cert.IsTrusted()
		chainItem := b.chainItem(cert, trustAnchorReached)
		if signingCertificate := bindingSignature.SigningCertificate(); signingCertificate != nil &&
			signingCertificate.Id() == cert.Id() {
			b.addQualifications(chainItem, cert)
			b.addQWACProfile(chainItem, cert)
		}
		chain = append(chain, chainItem)
	}
	xmlSignature.Chain = chain
}

// convertMessages is the port of the private convert(Collection<Message>).
func convertMessages(messages []detailedreport.Message) []*jaxb.XmlMessage {
	result := make([]*jaxb.XmlMessage, 0, len(messages))
	for _, m := range messages {
		xmlMessage := &jaxb.XmlMessage{}
		xmlMessage.Key = nilIfEmpty(m.Key)
		xmlMessage.Value = m.Value
		result = append(result, xmlMessage)
	}
	return result
}

// addConnectionDetails is the port of the private
// addConnectionDetails(XmlSimpleCertificateReport).
func (b *SimpleReportForCertificateBuilder) addConnectionDetails(simpleReport *jaxb.XmlSimpleCertificateReport) {
	if websiteUrl := b.diagnosticData.WebsiteUrl(); websiteUrl != "" {
		xmlConnectionDetails := &jaxb.XmlConnectionDetails{}
		xmlConnectionDetails.Url = &websiteUrl
		xmlConnectionDetails.TLSCertificateBindingLink = nilIfEmpty(b.diagnosticData.TLSCertificateBindingUrl())
		simpleReport.ConnectionDetails = xmlConnectionDetails
	}
}

// detailsNotEmpty is the port of the private isNotEmpty(XmlDetails).
func detailsNotEmpty(details *jaxb.XmlDetails) bool {
	return len(details.Error) > 0 || len(details.Warning) > 0 || len(details.Info) > 0
}

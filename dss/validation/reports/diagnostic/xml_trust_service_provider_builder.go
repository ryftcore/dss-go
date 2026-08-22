// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/XmlTrustServiceProviderBuilder.java (DSS 6.5.RC1).
//
// Java's slf4j logging (LOG.trace/.debug/.info/.warn) has no Go equivalent and is not ported;
// none of it was load-bearing (no control-flow decision depended on whether a message was
// logged).
package diagnostic

import (
	"fmt"
	"sort"

	dssdiag "github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	"github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// XmlTrustServiceProviderBuilder is used to build an XmlTrustServiceProvider object instance.
type XmlTrustServiceProviderBuilder struct {
	// xmlCertsMap is the map of certificate identifiers and their corresponding XML
	// representations.
	xmlCertsMap map[string]*jaxb.XmlCertificate

	// xmlTrustedListsMap is the map of trusted lists.
	xmlTrustedListsMap map[string]*jaxb.XmlTrustedList

	// tlInfoMap is the map between trusted list identifiers and corresponding TLInfo.
	tlInfoMap map[string]*tsl.TLInfo

	// qcStatementsBuilder is the builder for QcStatements.
	qcStatementsBuilder *XmlQcStatementsBuilder
}

// NewXmlTrustServiceProviderBuilder is the port of the default constructor.
func NewXmlTrustServiceProviderBuilder(xmlCertsMap map[string]*jaxb.XmlCertificate, xmlTrustedListsMap map[string]*jaxb.XmlTrustedList,
	tlInfoMap map[string]*tsl.TLInfo) *XmlTrustServiceProviderBuilder {
	return &XmlTrustServiceProviderBuilder{
		xmlCertsMap:         xmlCertsMap,
		xmlTrustedListsMap:  xmlTrustedListsMap,
		tlInfoMap:           tlInfoMap,
		qcStatementsBuilder: NewXmlQcStatementsBuilder(),
	}
}

// Build builds a list of XmlTrustServiceProviders corresponding to the given CertificateToken.
// Port of build(CertificateToken, Map<CertificateToken, List<TrustProperties>>).
func (b *XmlTrustServiceProviderBuilder) Build(certificateToken *model.CertificateToken,
	relatedTrustServices map[*model.CertificateToken][]*tsl.TrustProperties) []*jaxb.XmlTrustServiceProvider {
	result := make([]*jaxb.XmlTrustServiceProvider, 0)
	// Java iterates relatedTrustServices.entrySet() and servicesByProviders.entrySet(),
	// both HashMaps; ranging the Go maps here would order the produced
	// XmlTrustServiceProvider list by Go's randomised map iteration, i.e. differently on
	// every run. Java's exact bucket order is not reproducible in Go (see
	// diagnostic_data_builder.go's header), so - as everywhere else in this package - the
	// port substitutes a deterministic order: trust anchors by their DSS id, and, within
	// one anchor, service providers in first-encounter order of the services list.
	for _, trustedCert := range sortedTrustAnchors(relatedTrustServices) {
		for _, trustServices := range b.classifyByServiceProvider(relatedTrustServices[trustedCert]) {
			if utils.IsCollectionNotEmpty(trustServices) {
				result = append(result, b.getXmlTrustServiceProvider(certificateToken, trustServices, trustedCert))
			}
		}
	}
	return result
}

// sortedTrustAnchors returns the map's certificate keys ordered by their DSS id, so that
// the produced list does not depend on Go's map iteration order.
func sortedTrustAnchors[V any](m map[*model.CertificateToken]V) []*model.CertificateToken {
	keys := make([]*model.CertificateToken, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].DSSIDAsString() < keys[j].DSSIDAsString()
	})
	return keys
}

// classifyByServiceProvider groups the trust properties by their trust service provider,
// returning the groups in first-encounter order. Port of
// classifyByServiceProvider(List): Java returns the HashMap itself and the caller drains
// its entrySet(), which this port replaces by the ordered group list (see Build).
func (b *XmlTrustServiceProviderBuilder) classifyByServiceProvider(trustPropertiesList []*tsl.TrustProperties) [][]*tsl.TrustProperties {
	servicesByProviders := make(map[*tsl.TrustServiceProvider][]*tsl.TrustProperties)
	var order []*tsl.TrustServiceProvider
	if utils.IsCollectionNotEmpty(trustPropertiesList) {
		for _, trustProperties := range trustPropertiesList {
			currentTrustServiceProvider := trustProperties.TrustServiceProvider()
			if _, seen := servicesByProviders[currentTrustServiceProvider]; !seen {
				order = append(order, currentTrustServiceProvider)
			}
			servicesByProviders[currentTrustServiceProvider] = append(servicesByProviders[currentTrustServiceProvider], trustProperties)
		}
	}
	groups := make([][]*tsl.TrustProperties, 0, len(order))
	for _, provider := range order {
		groups = append(groups, servicesByProviders[provider])
	}
	return groups
}

func (b *XmlTrustServiceProviderBuilder) getXmlTrustServiceProvider(certificateToken *model.CertificateToken,
	trustServices []*tsl.TrustProperties, trustAnchor *model.CertificateToken) *jaxb.XmlTrustServiceProvider {
	trustProperties := trustServices[0]

	result := &jaxb.XmlTrustServiceProvider{}

	lotlInfo := trustProperties.LOTLInfo()
	if lotlInfo != nil {
		xmlLOTL, ok := b.xmlTrustedListsMap[lotlInfo.DSSIDAsString()]
		if !ok {
			panic(fmt.Sprintf("LOTL with Id '%s' has not been found! "+
				"Please verify TrustedListsCertificateSource contains TLValidationSummary.", lotlInfo.DSSIDAsString()))
		}
		result.LOTL = xmlLOTL
	}
	tlInfo := trustProperties.TLInfo()
	if tlInfo != nil {
		xmlTL, ok := b.xmlTrustedListsMap[tlInfo.DSSIDAsString()]
		if !ok {
			panic(fmt.Sprintf("TL with Id '%s' has not been found! "+
				"Please verify TrustedListsCertificateSource contains TLValidationSummary.", tlInfo.DSSIDAsString()))
		}
		result.TL = xmlTL
	}

	tsp := trustProperties.TrustServiceProvider()
	if names := b.getLangAndValues(tsp.Names()); names != nil {
		result.TSPNames = &jaxb.TSPNamesWrapper{Items: names}
	}
	if tradeNames := b.getLangAndValues(tsp.TradeNames()); tradeNames != nil {
		result.TSPTradeNames = &jaxb.TSPTradeNamesWrapper{Items: tradeNames}
	}
	result.TSPRegistrationIdentifiers = &jaxb.TSPRegistrationIdentifiersWrapper{Items: tsp.RegistrationIdentifiers()}

	result.TrustServices = &jaxb.TrustServicesWrapper{Items: b.buildXmlTrustServicesList(certificateToken, trustServices, trustAnchor)}

	return result
}

func (b *XmlTrustServiceProviderBuilder) getLangAndValues(m map[string][]string) []*jaxb.XmlLangAndValue {
	if utils.IsMapNotEmpty(m) {
		result := make([]*jaxb.XmlLangAndValue, 0)
		langs := make([]string, 0, len(m))
		for lang := range m {
			langs = append(langs, lang)
		}
		sort.Strings(langs)
		for _, lang := range langs {
			for _, value := range m[lang] {
				l, v := lang, value
				result = append(result, &jaxb.XmlLangAndValue{Lang: &l, Value: v})
			}
		}
		return result
	}
	return nil
}

func (b *XmlTrustServiceProviderBuilder) buildXmlTrustServicesList(certToken *model.CertificateToken,
	trustServices []*tsl.TrustProperties, trustAnchor *model.CertificateToken) []*jaxb.XmlTrustService {
	result := make([]*jaxb.XmlTrustService, 0)

	for _, trustProperties := range trustServices {
		trustService := trustProperties.TrustService()
		serviceStatusAfterOfEqualsCertIssuance := trustService.After(certToken.NotBefore())
		if utils.IsCollectionNotEmpty(serviceStatusAfterOfEqualsCertIssuance) {
			for _, serviceInfoStatus := range serviceStatusAfterOfEqualsCertIssuance {
				mra := b.getMRA(trustProperties)
				if mra != nil {
					result = append(result, b.buildXmlTrustServicesWithMRA(serviceInfoStatus, certToken, trustAnchor, mra)...)
				} else {
					result = append(result, b.getXmlTrustService(serviceInfoStatus, certToken, trustAnchor))
				}
			}
		}
	}
	return result
}

func (b *XmlTrustServiceProviderBuilder) getMRA(trustProperties *tsl.TrustProperties) *tsl.MRA {
	if trustProperties.TLInfo() != nil {
		tlInfo, ok := b.tlInfoMap[trustProperties.TLInfo().DSSIDAsString()]
		if ok && tlInfo != nil && tlInfo.OtherTSLPointer() != nil {
			// may be nil when no TLValidationJob is used
			return tlInfo.OtherTSLPointer().Mra()
		}
	}
	return nil
}

func (b *XmlTrustServiceProviderBuilder) getXmlTrustService(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken, trustAnchor *model.CertificateToken) *jaxb.XmlTrustService {
	trustService := &jaxb.XmlTrustService{}

	trustService.ServiceDigitalIdentifier = b.xmlCertsMap[trustAnchor.DSSIDAsString()]
	if names := b.getLangAndValues(serviceInfoStatus.Names()); names != nil {
		trustService.ServiceNames = &jaxb.ServiceNamesWrapper{Items: names}
	}
	serviceType := serviceInfoStatus.Type()
	trustService.ServiceType = &serviceType
	status := serviceInfoStatus.Status()
	trustService.Status = &status
	if !serviceInfoStatus.StartDate().IsZero() {
		trustService.StartDate = jaxb.NewXSDateTime(serviceInfoStatus.StartDate())
	}
	if !serviceInfoStatus.EndDate().IsZero() {
		trustService.EndDate = jaxb.NewXSDateTime(serviceInfoStatus.EndDate())
	}

	qualifiers := b.getQualifiers(serviceInfoStatus, certToken)
	if utils.IsCollectionNotEmpty(qualifiers) {
		trustService.CapturedQualifiers = &jaxb.CapturedQualifiersWrapper{Items: qualifiers}
	}

	additionalServiceInfoUris := serviceInfoStatus.AdditionalServiceInfoUris()
	if utils.IsCollectionNotEmpty(additionalServiceInfoUris) {
		trustService.AdditionalServiceInfoUris = &jaxb.AdditionalServiceInfoUrisWrapper{Items: additionalServiceInfoUris}
	}

	serviceSupplyPoints := serviceInfoStatus.ServiceSupplyPoints()
	if utils.IsCollectionNotEmpty(serviceSupplyPoints) {
		trustService.ServiceSupplyPoints = &jaxb.ServiceSupplyPointsWrapper{Items: serviceSupplyPoints}
	}

	if !serviceInfoStatus.ExpiredCertsRevocationInfo().IsZero() {
		trustService.ExpiredCertsRevocationInfo = jaxb.NewXSDateTime(serviceInfoStatus.ExpiredCertsRevocationInfo())
	}

	return trustService
}

// getQualifiers retrieves all the qualifiers for which the corresponding conditionEntry is
// true. Port of the private getQualifiers(TrustServiceStatusAndInformationExtensions,
// CertificateToken).
func (b *XmlTrustServiceProviderBuilder) getQualifiers(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certificateToken *model.CertificateToken) []*jaxb.XmlQualifier {
	list := make([]*jaxb.XmlQualifier, 0)
	conditionsForQualifiers := serviceInfoStatus.ConditionsForQualifiers()
	if utils.IsCollectionNotEmpty(conditionsForQualifiers) {
		for _, conditionForQualifiers := range conditionsForQualifiers {
			condition := conditionForQualifiers.Condition()
			if condition.Check(certificateToken) {
				for _, qualifier := range conditionForQualifiers.Qualifiers() {
					list = append(list, b.getXmlQualifier(qualifier, conditionForQualifiers.IsCritical()))
				}
			}
		}
	}
	return list
}

func (b *XmlTrustServiceProviderBuilder) getXmlQualifier(value string, critical bool) *jaxb.XmlQualifier {
	xmlQualifier := &jaxb.XmlQualifier{}
	xmlQualifier.Value = value
	xmlQualifier.Critical = &critical
	return xmlQualifier
}

func (b *XmlTrustServiceProviderBuilder) buildXmlTrustServicesWithMRA(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken, trustAnchor *model.CertificateToken, mra *tsl.MRA) []*jaxb.XmlTrustService {
	if utils.IsCollectionNotEmpty(serviceInfoStatus.AdditionalServiceInfoUris()) {
		result := make([]*jaxb.XmlTrustService, 0)
		for _, aSI := range serviceInfoStatus.AdditionalServiceInfoUris() {
			serviceInfoStatusCopy := tsl.NewTrustServiceStatusAndInformationExtensions(
				tsl.NewTrustServiceStatusAndInformationExtensionsBuilderFrom(serviceInfoStatus).
					SetAdditionalServiceInfoUris([]string{aSI}))
			result = append(result, b.getXmlTrustServicesForMRA(serviceInfoStatusCopy, certToken, trustAnchor, mra)...)
		}
		return result
	}
	return b.getXmlTrustServicesForMRA(serviceInfoStatus, certToken, trustAnchor, mra)
}

func (b *XmlTrustServiceProviderBuilder) getXmlTrustServicesForMRA(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken, trustAnchor *model.CertificateToken, mra *tsl.MRA) []*jaxb.XmlTrustService {
	mraEquivalences := b.getMRAServiceEquivalences(serviceInfoStatus, certToken, mra)
	enactedMra := utils.IsCollectionNotEmpty(mraEquivalences)
	if enactedMra {
		if len(mraEquivalences) == 1 {
			serviceEquivalenceValues := mraEquivalences[0]

			result := make([]*jaxb.XmlTrustService, 0)

			serviceEquivalenceList := serviceEquivalenceValues.After(certToken.NotBefore())
			for _, serviceEquivalence := range serviceEquivalenceList {
				// shall be computed before translation
				equivalent := b.getEquivalent(serviceInfoStatus, serviceEquivalence)

				xmlTrustService := b.getXmlTrustService(equivalent, certToken, trustAnchor)
				xmlTrustService.MRATrustServiceMapping = b.getXmlMRATrustServiceMapping(serviceInfoStatus, certToken, serviceEquivalence)
				enacted := serviceEquivalence.Status().IsEnacted()
				xmlTrustService.EnactedMRA = &enacted
				result = append(result, xmlTrustService)
			}
			b.translateCertificate(certToken, serviceEquivalenceList)

			return result
		}
		// Port of LOG.warn("More than one MRA equivalence found..."): slf4j dropped.
	}

	return []*jaxb.XmlTrustService{b.getXmlTrustService(serviceInfoStatus, certToken, trustAnchor)}
}

func (b *XmlTrustServiceProviderBuilder) getMRAServiceEquivalences(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken *model.CertificateToken, mra *tsl.MRA) []*timedependent.MutableTimeDependentValues[*tsl.ServiceEquivalence] {
	equivalences := make([]*timedependent.MutableTimeDependentValues[*tsl.ServiceEquivalence], 0)
	for _, serviceEquivalenceList := range mra.ServiceEquivalence() {
		// filter TrustServices that can be potentially applied to the validation
		for _, serviceEquivalence := range serviceEquivalenceList.After(certToken.NotBefore()) {
			if b.check(serviceInfoStatus, serviceEquivalence) {
				equivalences = append(equivalences, serviceEquivalenceList)
				break
			}
		}
	}
	return equivalences
}

func (b *XmlTrustServiceProviderBuilder) check(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions, serviceEquivalence *tsl.ServiceEquivalence) bool {
	if !b.checkServiceTypeAsiEquivalence(serviceInfoStatus, serviceEquivalence.TypeAsiEquivalence()) {
		return false
	}
	statusEquivalence := serviceEquivalence.StatusEquivalence()
	return b.checkStatusEquivalence(serviceInfoStatus, statusEquivalence)
}

func (b *XmlTrustServiceProviderBuilder) checkServiceTypeAsiEquivalence(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	typeAsiEquivalenceMap map[tsl.ServiceTypeASi]tsl.ServiceTypeASi) bool {
	if utils.IsMapEmpty(typeAsiEquivalenceMap) {
		return false
	}
	for typeAsiEquivalence := range typeAsiEquivalenceMap {
		if b.checkServiceTypeASi(serviceInfoStatus, typeAsiEquivalence) {
			return true
		}
	}
	return false
}

func (b *XmlTrustServiceProviderBuilder) checkServiceTypeASi(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions, serviceTypeASi tsl.ServiceTypeASi) bool {
	return serviceInfoStatus.Type() != "" && serviceInfoStatus.Type() == serviceTypeASi.Type() &&
		(serviceTypeASi.Asi() == "" || containsString(serviceInfoStatus.AdditionalServiceInfoUris(), serviceTypeASi.Asi()))
}

func containsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func (b *XmlTrustServiceProviderBuilder) checkCertTypeAsiEquivalence(certToken *model.CertificateToken,
	typeAsiEquivalenceMap map[tsl.ServiceTypeASi]tsl.ServiceTypeASi) bool {
	xmlCertificate, ok := b.xmlCertsMap[certToken.DSSIDAsString()]
	if !ok {
		panic(fmt.Sprintf("XML certificate with Id '%s' is not yet created!", certToken.DSSIDAsString()))
	}
	if utils.IsMapEmpty(typeAsiEquivalenceMap) {
		return false
	}

	certificateWrapper := dssdiag.NewCertificateWrapper(xmlCertificate)
	qcCompliance := certificateWrapper.IsQcCompliance()
	qcTypes := certificateWrapper.QcTypes()
	for _, serviceTypeASi := range typeAsiEquivalenceMap {
		if serviceTypeASi.Asi() == "" {
			// no aSI -> accept all
			return true
		}

		if utils.IsCollectionNotEmpty(qcTypes) {
			for _, qcType := range qcTypes {
				if b.isQcTypeMatch(qcType, serviceTypeASi) {
					return true
				}
			}
		} else if qcCompliance {
			// qcCompliance + no type -> foreSign
			if b.isQcTypeMatch(enumerations.QCTypeEnumQCTESign, serviceTypeASi) {
				return true
			}
		} else {
			// no qcType -> accept all
			return true
		}
	}
	return false
}

func (b *XmlTrustServiceProviderBuilder) isQcTypeMatch(qcType enumerations.QCType, serviceTypeASi tsl.ServiceTypeASi) bool {
	asi := serviceTypeASi.Asi()
	switch qcType.OID() {
	case enumerations.QCTypeEnumQCTESign.OID():
		return enumerations.AdditionalServiceInformationIsForeSignatures(asi)
	case enumerations.QCTypeEnumQCTESeal.OID():
		return enumerations.AdditionalServiceInformationIsForeSeals(asi)
	case enumerations.QCTypeEnumQCTWeb.OID():
		return enumerations.AdditionalServiceInformationIsForWebAuth(asi)
	}
	return false
}

func (b *XmlTrustServiceProviderBuilder) checkStatusEquivalence(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	statusEquivalenceMap []tsl.StatusEquivalenceMapping) bool {
	if utils.IsCollectionEmpty(statusEquivalenceMap) {
		return false
	}
	for _, statusEquivalence := range statusEquivalenceMap {
		if containsString(statusEquivalence.PointedStatuses, serviceInfoStatus.Status()) {
			return true
		}
	}
	return false
}

func (b *XmlTrustServiceProviderBuilder) getXmlMRATrustServiceMapping(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken *model.CertificateToken, serviceEquivalence *tsl.ServiceEquivalence) *jaxb.XmlMRATrustServiceMapping {
	mraTrustServiceMapping := &jaxb.XmlMRATrustServiceMapping{}
	legalID := serviceEquivalence.LegalInfoIdentifier()
	mraTrustServiceMapping.TrustServiceLegalIdentifier = &legalID
	if !serviceEquivalence.StartDate().IsZero() {
		mraTrustServiceMapping.EquivalenceStatusStartingTime = jaxb.NewXSDateTime(serviceEquivalence.StartDate())
	}
	if !serviceEquivalence.EndDate().IsZero() {
		mraTrustServiceMapping.EquivalenceStatusEndingTime = jaxb.NewXSDateTime(serviceEquivalence.EndDate())
	}
	mraTrustServiceMapping.OriginalThirdCountryMapping = b.getXmlOriginalThirdCountryTrustServiceMapping(serviceInfoStatus, certToken)
	return mraTrustServiceMapping
}

func (b *XmlTrustServiceProviderBuilder) getXmlOriginalThirdCountryTrustServiceMapping(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	certToken *model.CertificateToken) *jaxb.XmlOriginalThirdCountryTrustServiceMapping {
	originalThirdCountryMapping := &jaxb.XmlOriginalThirdCountryTrustServiceMapping{}
	serviceType := serviceInfoStatus.Type()
	originalThirdCountryMapping.ServiceType = &serviceType
	status := serviceInfoStatus.Status()
	originalThirdCountryMapping.Status = &status

	qualifiers := b.getQualifiers(serviceInfoStatus, certToken)
	if utils.IsCollectionNotEmpty(qualifiers) {
		originalThirdCountryMapping.CapturedQualifiers = &jaxb.CapturedQualifiersWrapper{Items: qualifiers}
	}

	additionalServiceInfoUris := serviceInfoStatus.AdditionalServiceInfoUris()
	if utils.IsCollectionNotEmpty(additionalServiceInfoUris) {
		originalThirdCountryMapping.AdditionalServiceInfoUris = &jaxb.AdditionalServiceInfoUrisWrapper{Items: additionalServiceInfoUris}
	}

	return originalThirdCountryMapping
}

func (b *XmlTrustServiceProviderBuilder) translateCertificate(certToken *model.CertificateToken, serviceEquivalenceList []*tsl.ServiceEquivalence) {
	// See PRO-4.3.4-03B (apply first enacted serviceEquivalence. NOTE: list inverted)
	var qcStatements *jaxb.XmlQcStatements
	for _, serviceEquivalence := range serviceEquivalenceList {
		if serviceEquivalence.Status().IsEnacted() && b.checkEnacted(certToken, serviceEquivalence) {
			currentQcStatement := b.applyCertContentEquivalence(certToken, serviceEquivalence)
			if qcStatements == nil {
				qcStatements = currentQcStatement
			} else if !b.checkQcStatementsEquivalence(qcStatements, currentQcStatement) {
				// Port of LOG.warn("Enacted MRA equivalences ... lead to different certificate content results..."): slf4j dropped.
				return
			}
		}
		// else: Port of LOG.debug("MRA equivalence was not applied for a certificate..."): slf4j dropped.
	}

	if qcStatements != nil {
		// Port of LOG.info("MRA equivalence is applied for a certificate..."): slf4j dropped.

		// update QcStatements certificate content
		xmlCertificate := b.xmlCertsMap[certToken.DSSIDAsString()]
		b.setQcStatements(xmlCertificate, qcStatements)
	}
}

// checkEnacted checks if a certificate has been issued at or after the starting time of the
// service equivalence. Port of the private check(CertificateToken, ServiceEquivalence).
func (b *XmlTrustServiceProviderBuilder) checkEnacted(certificateToken *model.CertificateToken, serviceEquivalence *tsl.ServiceEquivalence) bool {
	certIssuance := certificateToken.NotBefore()
	if certIssuance.Before(serviceEquivalence.StartDate()) {
		return false
	}
	if !serviceEquivalence.EndDate().IsZero() && certIssuance.After(serviceEquivalence.EndDate()) {
		return false
	}
	return b.checkCertTypeAsiEquivalence(certificateToken, serviceEquivalence.TypeAsiEquivalence())
}

func (b *XmlTrustServiceProviderBuilder) getEquivalent(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	serviceEquivalence *tsl.ServiceEquivalence) *tsl.TrustServiceStatusAndInformationExtensions {
	typeASiSubstitution := b.getTypeASiSubstitution(serviceInfoStatus, serviceEquivalence)
	status := b.getStatusSubstitution(serviceInfoStatus, serviceEquivalence)
	qualifiersSubstitution := b.getQualifiersSubstitution(serviceInfoStatus, serviceEquivalence)

	builder := tsl.NewTrustServiceStatusAndInformationExtensionsBuilder()
	if typeASiSubstitution != nil {
		builder.SetType(typeASiSubstitution.Type())
		if typeASiSubstitution.Asi() != "" {
			builder.SetAdditionalServiceInfoUris([]string{typeASiSubstitution.Asi()})
		}
	}
	builder.SetStatus(status)
	builder.SetConditionsForQualifiers(qualifiersSubstitution)
	// copy
	builder.SetStartDate(serviceInfoStatus.StartDate())
	builder.SetEndDate(serviceInfoStatus.EndDate())
	builder.SetNames(serviceInfoStatus.Names())
	builder.SetExpiredCertsRevocationInfo(serviceInfoStatus.ExpiredCertsRevocationInfo())
	builder.SetServiceSupplyPoints(serviceInfoStatus.ServiceSupplyPoints())
	return tsl.NewTrustServiceStatusAndInformationExtensions(builder)
}

func (b *XmlTrustServiceProviderBuilder) getTypeASiSubstitution(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	serviceEquivalence *tsl.ServiceEquivalence) *tsl.ServiceTypeASi {
	if utils.IsMapEmpty(serviceEquivalence.TypeAsiEquivalence()) {
		return nil
	}
	for expected, pointing := range serviceEquivalence.TypeAsiEquivalence() {
		if b.checkServiceTypeASi(serviceInfoStatus, expected) {
			return b.substituteTypeASi(serviceInfoStatus, expected, pointing)
		}
	}
	return nil
}

func (b *XmlTrustServiceProviderBuilder) substituteTypeASi(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	pointed, pointing tsl.ServiceTypeASi) *tsl.ServiceTypeASi {
	serviceTypeASi := tsl.NewServiceTypeASi()
	serviceTypeASi.SetType(pointing.Type())

	var asiResult string
	if utils.IsCollectionNotEmpty(serviceInfoStatus.AdditionalServiceInfoUris()) {
		asiResult = serviceInfoStatus.AdditionalServiceInfoUris()[0]
		if asiResult == pointed.Asi() {
			asiResult = pointing.Asi()
		}
	} else {
		asiResult = pointing.Asi()
	}
	serviceTypeASi.SetAsi(asiResult)

	return serviceTypeASi
}

func (b *XmlTrustServiceProviderBuilder) getStatusSubstitution(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	serviceEquivalence *tsl.ServiceEquivalence) string {
	statusEquivalence := serviceEquivalence.StatusEquivalence()
	if utils.IsCollectionEmpty(statusEquivalence) {
		return ""
	}
	for _, equivalence := range statusEquivalence {
		if containsString(equivalence.PointedStatuses, serviceInfoStatus.Status()) {
			if len(equivalence.PointingStatuses) > 0 {
				return equivalence.PointingStatuses[0]
			}
			return ""
		}
	}
	return ""
}

func (b *XmlTrustServiceProviderBuilder) getQualifiersSubstitution(serviceInfoStatus *tsl.TrustServiceStatusAndInformationExtensions,
	serviceEquivalence *tsl.ServiceEquivalence) []*tsl.ConditionForQualifiers {
	result := make([]*tsl.ConditionForQualifiers, 0)
	qualifierEquivalence := serviceEquivalence.QualifierEquivalence()
	if utils.IsMapEmpty(qualifierEquivalence) {
		return result
	}
	for _, qualifierCondition := range serviceInfoStatus.ConditionsForQualifiers() {
		qualifiers := make([]string, 0)
		for _, qualifier := range qualifierCondition.Qualifiers() {
			pointingQualifier, ok := qualifierEquivalence[qualifier]
			if ok && utils.IsStringNotEmpty(pointingQualifier) {
				qualifier = pointingQualifier
			}
			qualifiers = append(qualifiers, qualifier)
		}
		result = append(result, tsl.NewConditionForQualifiersWithCriticality(qualifierCondition.Condition(), qualifiers, qualifierCondition.IsCritical()))
	}
	return result
}

func (b *XmlTrustServiceProviderBuilder) applyCertContentEquivalence(certToken *model.CertificateToken, serviceEquivalence *tsl.ServiceEquivalence) *jaxb.XmlQcStatements {
	certificateContentEquivalences := serviceEquivalence.CertificateContentEquivalences()
	if utils.IsCollectionEmpty(certificateContentEquivalences) {
		return nil
	}
	b.assertCertificateContentEquivalenceListIsConsistent(certificateContentEquivalences)

	xmlCertificate, ok := b.xmlCertsMap[certToken.DSSIDAsString()]
	if !ok {
		panic(fmt.Sprintf("XmlCertificate with Id '%s' is not yet created!", certToken.DSSIDAsString()))
	}

	// Overwrite with information from MRA
	qcStatements := b.getQcStatements(xmlCertificate)
	enactedMRA := true
	qcStatements.EnactedMRA = &enactedMRA

	xmlMRACertificateMapping := b.getXmlMRACertificateMapping(qcStatements, serviceEquivalence)
	qcStatements.MRACertificateMapping = xmlMRACertificateMapping

	trustServiceEquivalenceInformation := xmlMRACertificateMapping.TrustServiceEquivalenceInformation

	for _, certificateContentEquivalence := range certificateContentEquivalences {
		equivalenceContext := certificateContentEquivalence.Context()

		if equivalenceContext != "" {
			xmlCertificateContentEquivalence := &jaxb.XmlCertificateContentEquivalence{}
			uri := equivalenceContext.URI()
			xmlCertificateContentEquivalence.Uri = &uri

			condition := certificateContentEquivalence.Condition()
			if condition.Check(certToken) {
				// Port of LOG.info("MRA condition match ({})", equivalenceContext): slf4j dropped.

				contentReplacement := certificateContentEquivalence.ContentReplacement()
				switch equivalenceContext {
				case enumerations.MRAEquivalenceContextQCCompliance:
					b.replaceCompliance(qcStatements, contentReplacement)
					xmlCertificateContentEquivalence.Enacted = true
				case enumerations.MRAEquivalenceContextQCType:
					b.replaceType(qcStatements, contentReplacement)
					xmlCertificateContentEquivalence.Enacted = true
				case enumerations.MRAEquivalenceContextQCQSCD:
					b.replaceQSCD(qcStatements, contentReplacement)
					xmlCertificateContentEquivalence.Enacted = true
				default:
					// Port of LOG.warn("Unsupported equivalence context {}", equivalenceContext): slf4j dropped.
				}
			}

			trustServiceEquivalenceInformation.CertificateContentEquivalenceList = &jaxb.CertificateContentEquivalenceListWrapper{
				Items: append(trustServiceEquivalenceInformation.CertificateContentEquivalenceList.All(), xmlCertificateContentEquivalence),
			}
		}
	}

	return qcStatements
}

func (b *XmlTrustServiceProviderBuilder) assertCertificateContentEquivalenceListIsConsistent(certificateContentEquivalences []*tsl.CertificateContentEquivalence) {
	// Port of the LOG.warn-only duplicate-context check: slf4j dropped, and the check has no
	// other observable effect, so it is not reproduced.
	_ = certificateContentEquivalences
}

// getQcStatements returns a deep copy of the XmlCertificate's XmlQcStatements extension, when
// present. Empty object otherwise. Port of the private getQcStatements(XmlCertificate).
func (b *XmlTrustServiceProviderBuilder) getQcStatements(xmlCertificate *jaxb.XmlCertificate) *jaxb.XmlQcStatements {
	for _, certificateExtension := range xmlCertificate.CertificateExtensions.All() {
		if qc, ok := certificateExtension.(*jaxb.XmlQcStatements); ok && qc.ExtensionOID() != nil &&
			enumerations.CertificateExtensionEnumQCStatements.OID() == *qc.ExtensionOID() {
			return b.qcStatementsBuilder.Copy(qc)
		}
	}
	return &jaxb.XmlQcStatements{}
}

// setQcStatements sets a new XmlQcStatements certificate extension on the given XmlCertificate,
// replacing an existing one when present. Port of the private
// setQcStatements(XmlCertificate, XmlQcStatements).
func (b *XmlTrustServiceProviderBuilder) setQcStatements(xmlCertificate *jaxb.XmlCertificate, xmlQcStatements *jaxb.XmlQcStatements) {
	items := xmlCertificate.CertificateExtensions.All()
	filtered := make([]jaxb.XmlCertificateExtensionItem, 0, len(items)+1)
	for _, certificateExtension := range items {
		if qc, ok := certificateExtension.(*jaxb.XmlQcStatements); ok && qc.ExtensionOID() != nil &&
			enumerations.CertificateExtensionEnumQCStatements.OID() == *qc.ExtensionOID() {
			continue
		}
		filtered = append(filtered, certificateExtension)
	}
	filtered = append(filtered, xmlQcStatements)
	xmlCertificate.CertificateExtensions = &jaxb.CertificateExtensionsWrapper{Items: filtered}
}

func (b *XmlTrustServiceProviderBuilder) checkQcStatementsEquivalence(qcStatementsOne, qcStatementsTwo *jaxb.XmlQcStatements) bool {
	if qcStatementsOne == nil && qcStatementsTwo == nil {
		return true
	} else if qcStatementsOne == nil || qcStatementsTwo == nil {
		return false
	}

	if utils.IsTrue(qcStatementsOne.EnactedMRA) != utils.IsTrue(qcStatementsTwo.EnactedMRA) {
		return false
	}
	if (qcStatementsOne.QcCompliance != nil && qcStatementsOne.QcCompliance.Present) !=
		(qcStatementsTwo.QcCompliance != nil && qcStatementsTwo.QcCompliance.Present) {
		return false
	}
	if !equalOIDSets(qcStatementsOne.QcTypes.All(), qcStatementsTwo.QcTypes.All()) {
		return false
	}
	if (qcStatementsOne.QcSSCD != nil && qcStatementsOne.QcSSCD.Present) !=
		(qcStatementsTwo.QcSSCD != nil && qcStatementsTwo.QcSSCD.Present) {
		return false
	}
	return true
}

func equalOIDSets(a, b []*jaxb.XmlOID) bool {
	setA := map[string]bool{}
	for _, x := range a {
		setA[x.Value] = true
	}
	setB := map[string]bool{}
	for _, x := range b {
		setB[x.Value] = true
	}
	if len(setA) != len(setB) {
		return false
	}
	for k := range setA {
		if !setB[k] {
			return false
		}
	}
	return true
}

func (b *XmlTrustServiceProviderBuilder) getXmlMRACertificateMapping(qcStatements *jaxb.XmlQcStatements, serviceEquivalence *tsl.ServiceEquivalence) *jaxb.XmlMRACertificateMapping {
	xmlMRACertificateMapping := &jaxb.XmlMRACertificateMapping{}
	xmlMRACertificateMapping.TrustServiceEquivalenceInformation = b.getXmlTrustServiceEquivalenceInformation(serviceEquivalence)
	xmlMRACertificateMapping.OriginalThirdCountryMapping = b.getXmlOriginalThirdCountryQcStatementsMapping(qcStatements)
	return xmlMRACertificateMapping
}

func (b *XmlTrustServiceProviderBuilder) getXmlTrustServiceEquivalenceInformation(serviceEquivalence *tsl.ServiceEquivalence) *jaxb.XmlTrustServiceEquivalenceInformation {
	xmlTrustServiceEquivalenceInformation := &jaxb.XmlTrustServiceEquivalenceInformation{}
	legalID := serviceEquivalence.LegalInfoIdentifier()
	xmlTrustServiceEquivalenceInformation.TrustServiceLegalIdentifier = &legalID
	return xmlTrustServiceEquivalenceInformation
}

func (b *XmlTrustServiceProviderBuilder) getXmlOriginalThirdCountryQcStatementsMapping(qcStatements *jaxb.XmlQcStatements) *jaxb.XmlOriginalThirdCountryQcStatementsMapping {
	originalQcStatements := &jaxb.XmlOriginalThirdCountryQcStatementsMapping{}
	if qcStatements.QcCompliance != nil {
		originalQcStatements.QcCompliance = b.qcStatementsBuilder.BuildXmlQcCompliance(qcStatements.QcCompliance.Present)
	}
	if qcStatements.QcSSCD != nil {
		originalQcStatements.QcSSCD = b.qcStatementsBuilder.BuildXmlQcSSCD(qcStatements.QcSSCD.Present)
	}
	originalQcTypes := qcStatements.QcTypes.All()
	if utils.IsCollectionNotEmpty(originalQcTypes) {
		originalQcStatements.QcTypes = &jaxb.QcTypesWrapper{Items: append([]*jaxb.XmlOID{}, originalQcTypes...)}
	}
	qcCClegislations := qcStatements.QcCClegislation.All()
	if utils.IsCollectionNotEmpty(qcCClegislations) {
		originalQcStatements.QcCClegislation = &jaxb.QcCClegislationWrapper{Items: append([]string{}, qcCClegislations...)}
	}
	qcQSCDlegislations := qcStatements.QcQSCDlegislation.All()
	if utils.IsCollectionNotEmpty(qcQSCDlegislations) {
		originalQcStatements.QcQSCDlegislation = &jaxb.QcQSCDlegislationWrapper{Items: append([]string{}, qcQSCDlegislations...)}
	}
	otherOIDs := qcStatements.OtherOIDs.All()
	if utils.IsCollectionNotEmpty(otherOIDs) {
		originalQcStatements.OtherOIDs = &jaxb.OtherOIDsWrapper{Items: append([]*jaxb.XmlOID{}, otherOIDs...)}
	}
	return originalQcStatements
}

func (b *XmlTrustServiceProviderBuilder) replaceCompliance(qcStatements *jaxb.XmlQcStatements, contentReplacement *tsl.QCStatementOids) {
	isQcCompliance := false
	qcCClegislations := qcStatements.QcCClegislation.All()
	for _, oid := range contentReplacement.QcStatementIds() {
		if spi.QcStatementUtilsIsQcCompliance(oid) {
			isQcCompliance = true
		}
	}
	if utils.IsCollectionNotEmpty(contentReplacement.QcCClegislations()) {
		qcCClegislations = contentReplacement.QcCClegislations()
	}
	for _, oid := range contentReplacement.QcStatementIdsToRemove() {
		if spi.QcStatementUtilsIsQcCompliance(oid) {
			isQcCompliance = false
		}
		if spi.QcStatementUtilsIsQcCClegislation(oid) {
			if utils.IsCollectionNotEmpty(contentReplacement.QcCClegislationsToRemove()) {
				qcCClegislations = removeAllStrings(qcCClegislations, contentReplacement.QcCClegislationsToRemove())
			} else {
				qcCClegislations = nil
			}
		}
	}
	qcStatements.QcCompliance = b.qcStatementsBuilder.BuildXmlQcCompliance(isQcCompliance)
	qcStatements.QcCClegislation = &jaxb.QcCClegislationWrapper{Items: qcCClegislations}
}

func removeAllStrings(from, toRemove []string) []string {
	remove := map[string]bool{}
	for _, r := range toRemove {
		remove[r] = true
	}
	result := make([]string, 0, len(from))
	for _, v := range from {
		if !remove[v] {
			result = append(result, v)
		}
	}
	return result
}

func (b *XmlTrustServiceProviderBuilder) replaceType(qcStatements *jaxb.XmlQcStatements, contentReplacement *tsl.QCStatementOids) {
	originalQcTypes := qcStatements.QcTypes.All()
	qcTypesIds := make([]string, 0, len(originalQcTypes))
	for _, oid := range originalQcTypes {
		qcTypesIds = append(qcTypesIds, oid.Value)
	}
	if utils.IsCollectionNotEmpty(contentReplacement.QcTypeIds()) {
		qcTypesIds = contentReplacement.QcTypeIds()
	}
	for _, oid := range contentReplacement.QcStatementIdsToRemove() {
		if spi.QcStatementUtilsIsQcType(oid) {
			if utils.IsCollectionNotEmpty(contentReplacement.QcTypeIdsToRemove()) {
				qcTypesIds = removeAllStrings(qcTypesIds, contentReplacement.QcTypeIdsToRemove())
			} else {
				qcTypesIds = nil
			}
		}
	}
	qcTypes := spi.QcStatementUtilsQcTypesForOIDs(qcTypesIds)
	qcStatements.QcTypes = &jaxb.QcTypesWrapper{Items: b.qcStatementsBuilder.BuildXmlQcTypes(qcTypes)}
}

func (b *XmlTrustServiceProviderBuilder) replaceQSCD(qcStatements *jaxb.XmlQcStatements, contentReplacement *tsl.QCStatementOids) {
	isQcSSCD := false
	for _, oid := range contentReplacement.QcStatementIds() {
		if spi.QcStatementUtilsIsQcSSCD(oid) {
			isQcSSCD = true
		}
	}
	for _, oid := range contentReplacement.QcStatementIdsToRemove() {
		if spi.QcStatementUtilsIsQcSSCD(oid) {
			isQcSSCD = false
		}
	}
	qcStatements.QcSSCD = b.qcStatementsBuilder.BuildXmlQcSSCD(isQcSSCD)
}

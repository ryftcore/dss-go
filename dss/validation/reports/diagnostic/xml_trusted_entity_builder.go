// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/lote/XmlTrustedEntityBuilder.java (DSS 6.5.RC1).
//
// Java package eu.europa.esig.dss.validation.reports.diagnostic.lote flattens into this same
// Go package (pkg diagnostic); the type keeps its Java name unqualified.
package diagnostic

import (
	"fmt"
	"sort"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/lote"
	"github.com/ryftcore/dss-go/dss/utils"
)

// XmlTrustedEntityBuilder builds an eu.europa.esig.dss.diagnostic.jaxb.XmlTrustedEntity
// instance.
type XmlTrustedEntityBuilder struct {
	// xmlCertsMap is the map of certificates identifiers and their corresponding XML
	// representations.
	xmlCertsMap map[string]*jaxb.XmlCertificate

	// xmlTrustSourceListsMap is the map of Trust Sources.
	xmlTrustSourceListsMap map[string]*jaxb.XmlListOfTrustedEntities
}

// NewXmlTrustedEntityBuilder is the port of the default constructor.
func NewXmlTrustedEntityBuilder(xmlCertsMap map[string]*jaxb.XmlCertificate,
	xmlTrustSourceListsMap map[string]*jaxb.XmlListOfTrustedEntities) *XmlTrustedEntityBuilder {
	return &XmlTrustedEntityBuilder{
		xmlCertsMap:            xmlCertsMap,
		xmlTrustSourceListsMap: xmlTrustSourceListsMap,
	}
}

// Build builds a list of XmlTrustedEntitys corresponding to the given CertificateToken. Port
// of build(CertificateToken, Map<CertificateToken, List<TrustedProperties>>).
func (b *XmlTrustedEntityBuilder) Build(certificateToken *model.CertificateToken,
	relatedTrustedProperties map[*model.CertificateToken][]*lote.TrustedProperties) []*jaxb.XmlTrustedEntity {
	result := make([]*jaxb.XmlTrustedEntity, 0)
	// See XmlTrustServiceProviderBuilder.Build: Java drains two HashMap entrySets here,
	// and ranging the Go maps would make the produced XmlTrustedEntity list order vary
	// between runs. The port substitutes the same deterministic order - trust anchors by
	// DSS id, trusted entities in first-encounter order.
	for _, trustedCert := range sortedTrustAnchors(relatedTrustedProperties) {
		for _, trustServices := range b.classifyByServiceProvider(relatedTrustedProperties[trustedCert]) {
			if utils.IsCollectionNotEmpty(trustServices) {
				result = append(result, b.getXmlTrustedEntity(certificateToken, trustServices, trustedCert))
			}
		}
	}
	return result
}

// classifyByServiceProvider groups the trusted properties by their trusted entity,
// returning the groups in first-encounter order. Port of
// classifyByServiceProvider(List), whose HashMap entrySet the caller drains.
func (b *XmlTrustedEntityBuilder) classifyByServiceProvider(
	trustPropertiesList []*lote.TrustedProperties) [][]*lote.TrustedProperties {
	servicesByProviders := make(map[*lote.TrustedEntity][]*lote.TrustedProperties)
	var order []*lote.TrustedEntity
	if utils.IsCollectionNotEmpty(trustPropertiesList) {
		for _, trustProperties := range trustPropertiesList {
			currentTrustedEntity := trustProperties.TrustedEntity()
			if _, seen := servicesByProviders[currentTrustedEntity]; !seen {
				order = append(order, currentTrustedEntity)
			}
			servicesByProviders[currentTrustedEntity] = append(servicesByProviders[currentTrustedEntity], trustProperties)
		}
	}
	groups := make([][]*lote.TrustedProperties, 0, len(order))
	for _, entity := range order {
		groups = append(groups, servicesByProviders[entity])
	}
	return groups
}

func (b *XmlTrustedEntityBuilder) getXmlTrustedEntity(certificateToken *model.CertificateToken,
	trustServices []*lote.TrustedProperties, trustAnchor *model.CertificateToken) *jaxb.XmlTrustedEntity {
	trustProperties := trustServices[0]

	result := &jaxb.XmlTrustedEntity{}

	loloteInfo := trustProperties.LoLoTEInfo()
	if loloteInfo != nil {
		xmlLoLoTE, ok := b.xmlTrustSourceListsMap[loloteInfo.DSSIDAsString()]
		if !ok {
			panic(fmt.Sprintf("LOTL with Id '%s' has not been found! "+
				"Please verify TrustedPropertiesCertificateSource contains LoTEValidationSummary.", loloteInfo.DSSIDAsString()))
		}
		// XmlTrustedEntity.LoLoTE is a by-reference stub (Java IDREF to the abstract
		// TrustSourceList type): it carries only the Id of the already-built
		// XmlListOfTrustedEntities so marshaling writes the reference attribute.
		result.LoLoTE = &jaxb.XmlTrustSourceList{XmlTrustSourceListAttrs: jaxb.XmlTrustSourceListAttrs{Id: xmlLoLoTE.Id}}
	}
	loteInfo := trustProperties.LoTEInfo()
	if loteInfo != nil {
		xmlLoTE, ok := b.xmlTrustSourceListsMap[loteInfo.DSSIDAsString()]
		if !ok {
			panic(fmt.Sprintf("TL with Id '%s' has not been found! "+
				"Please verify TrustedPropertiesCertificateSource contains LoTEValidationSummary.", loteInfo.DSSIDAsString()))
		}
		result.LoTE = &jaxb.XmlTrustSourceList{XmlTrustSourceListAttrs: jaxb.XmlTrustSourceListAttrs{Id: xmlLoTE.Id}}
	}

	tsp := trustProperties.TrustedEntity()
	if names := b.getLangAndValues(tsp.Names()); names != nil {
		result.Names = &jaxb.NamesWrapper{Items: names}
	}
	if tradeNames := b.getLangAndValues(tsp.TradeNames()); tradeNames != nil {
		result.TradeNames = &jaxb.TradeNamesWrapper{Items: tradeNames}
	}
	result.RegistrationIdentifiers = &jaxb.RegistrationIdentifiersWrapper{Items: tsp.RegistrationIdentifiers()}

	result.TrustedEntityServices = &jaxb.TrustedEntityServicesWrapper{
		Items: b.buildXmlTrustedEntityServicesList(certificateToken, trustServices, trustAnchor),
	}

	return result
}

// getLangAndValues ports the private getLangAndValues(Map<String, List<String>>).
//
// DIVERGENCE, deliberate: Java drains the map's entrySet(), whose order is the hash order of the
// map implementation the caller supplied; ranging the Go map would make the emitted
// XmlLangAndValue order vary from run to run. The port substitutes a deterministic order -
// languages sorted, each language's values in their given order - exactly as
// XmlTrustServiceProviderBuilder.getLangAndValues does. Only the order of entries of different
// languages can differ from upstream's; the multiset of entries never does.
func (b *XmlTrustedEntityBuilder) getLangAndValues(m map[string][]string) []*jaxb.XmlLangAndValue {
	if utils.IsMapNotEmpty(m) {
		result := make([]*jaxb.XmlLangAndValue, 0)
		langs := make([]string, 0, len(m))
		for lang := range m {
			langs = append(langs, lang)
		}
		sort.Strings(langs)
		for _, lang := range langs {
			values := m[lang]
			for _, value := range values {
				l, v := lang, value
				result = append(result, &jaxb.XmlLangAndValue{Lang: &l, Value: v})
			}
		}
		return result
	}
	return nil
}

func (b *XmlTrustedEntityBuilder) buildXmlTrustedEntityServicesList(certToken *model.CertificateToken,
	trustServices []*lote.TrustedProperties, trustAnchor *model.CertificateToken) []*jaxb.XmlTrustedEntityService {
	result := make([]*jaxb.XmlTrustedEntityService, 0)

	for _, trustProperties := range trustServices {
		trustService := trustProperties.TrustedServices()
		serviceStatusAfterOfEqualsCertIssuance := trustService.After(certToken.NotBefore())
		if utils.IsCollectionNotEmpty(serviceStatusAfterOfEqualsCertIssuance) {
			for _, serviceInfoStatus := range serviceStatusAfterOfEqualsCertIssuance {
				result = append(result, b.getXmlTrustedEntityService(serviceInfoStatus, certToken, trustAnchor))
			}
		}
	}
	return result
}

func (b *XmlTrustedEntityBuilder) getXmlTrustedEntityService(serviceInfoStatus lote.ServiceStatusAndInformationExtensions,
	certToken *model.CertificateToken, trustAnchor *model.CertificateToken) *jaxb.XmlTrustedEntityService {
	trustService := &jaxb.XmlTrustedEntityService{}

	trustService.ServiceDigitalIdentifier = b.xmlCertsMap[trustAnchor.DSSIDAsString()]
	if names := b.getLangAndValues(serviceInfoStatus.Names()); names != nil {
		trustService.ServiceNames = &jaxb.ServiceNamesWrapper{Items: names}
	}
	serviceType := serviceInfoStatus.Type()
	trustService.ServiceType = &serviceType
	status := serviceInfoStatus.Status()
	trustService.Status = &status
	startDate := jaxb.NewXSDateTime(serviceInfoStatus.StartDate())
	trustService.StartDate = startDate
	if !serviceInfoStatus.EndDate().IsZero() {
		trustService.EndDate = jaxb.NewXSDateTime(serviceInfoStatus.EndDate())
	}

	serviceSupplyPoints := serviceInfoStatus.ServiceSupplyPoints()
	if utils.IsCollectionNotEmpty(serviceSupplyPoints) {
		trustService.ServiceSupplyPoints = &jaxb.ServiceSupplyPointsWrapper{Items: serviceSupplyPoints}
	}

	return trustService
}

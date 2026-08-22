// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/TrustServiceConverter.java (DSS 6.5.RC1).
package tsl

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TrustServiceConverter converts a TSPServiceType to a TrustService.
type TrustServiceConverter struct{}

// NewTrustServiceConverter is the default constructor. Port of TrustServiceConverter().
func NewTrustServiceConverter() *TrustServiceConverter {
	return &TrustServiceConverter{}
}

// Apply ports apply(TSPServiceType).
func (c *TrustServiceConverter) Apply(original *jaxb.TSPServiceType) *tslmodel.TrustService {
	trustServiceBuilder := tslmodel.NewTrustServiceBuilder()
	if original.ServiceInformation != nil {
		trustServiceBuilder.
			SetCertificates(c.extractCertificates(original.ServiceInformation)).
			SetStatusAndInformationExtensions(c.extractStatusAndHistory(original))
	}
	return trustServiceBuilder.Build()
}

func (c *TrustServiceConverter) extractCertificates(serviceInformation *jaxb.TSPServiceInformationType) []*model.CertificateToken {
	converter := NewDigitalIdentityListTypeConverter()
	return converter.Apply(serviceInformation.ServiceDigitalIdentity)
}

func (c *TrustServiceConverter) extractStatusAndHistory(original *jaxb.TSPServiceType) *timedependent.Values[*tslmodel.TrustServiceStatusAndInformationExtensions] {
	statusHistoryList := timedependent.NewMutableTimeDependentValues[*tslmodel.TrustServiceStatusAndInformationExtensions]()

	serviceInfo := original.ServiceInformation

	converter := NewInternationalNamesTypeConverter()

	statusBuilder := tslmodel.NewTrustServiceStatusAndInformationExtensionsBuilder()
	statusBuilder.SetNames(converter.Apply(serviceInfo.ServiceName))
	statusBuilder.SetType(serviceInfo.ServiceTypeIdentifier)
	statusBuilder.SetStatus(serviceInfo.ServiceStatus)
	statusBuilder.SetServiceSupplyPoints(c.serviceSupplyPoints(serviceInfo.ServiceSupplyPoints))

	c.parseExtensionsList(serviceInfo.ServiceInformationExtensions, statusBuilder)

	nextEndDate := abstractParsingTaskConvertToDate(serviceInfo.StatusStartingTime)
	statusBuilder.SetStartDate(nextEndDate)
	statusHistoryList.AddOldest(statusBuilder.Build())

	if original.ServiceHistory != nil && utils.IsCollectionNotEmpty(original.ServiceHistory.ServiceHistoryInstance) {
		for _, serviceHistory := range original.ServiceHistory.ServiceHistoryInstance {
			if serviceHistory.StatusStartingTime == nil {
				continue
			}

			statusHistoryBuilder := tslmodel.NewTrustServiceStatusAndInformationExtensionsBuilder()
			statusHistoryBuilder.SetNames(converter.Apply(serviceHistory.ServiceName))
			statusHistoryBuilder.SetType(serviceHistory.ServiceTypeIdentifier)
			statusHistoryBuilder.SetStatus(serviceHistory.ServiceStatus)

			c.parseExtensionsList(serviceHistory.ServiceInformationExtensions, statusHistoryBuilder)

			statusHistoryBuilder.SetEndDate(nextEndDate)
			nextEndDate = abstractParsingTaskConvertToDate(serviceHistory.StatusStartingTime)
			statusHistoryBuilder.SetStartDate(nextEndDate)
			statusHistoryList.AddOldest(statusHistoryBuilder.Build())
		}
	}

	return &statusHistoryList.Values
}

func (c *TrustServiceConverter) parseExtensionsList(serviceInformationExtensions *jaxb.ExtensionsListType,
	statusBuilder *tslmodel.TrustServiceStatusAndInformationExtensionsBuilder) {
	if serviceInformationExtensions != nil {
		statusBuilder.SetConditionsForQualifiers(c.extractConditionsForQualifiers(serviceInformationExtensions.Extension))
		statusBuilder.SetAdditionalServiceInfoUris(c.extractAdditionalServiceInfoUris(serviceInformationExtensions.Extension))
		statusBuilder.SetExpiredCertsRevocationInfo(c.extractExpiredCertsRevocationInfo(serviceInformationExtensions.Extension))
	}
}

func (c *TrustServiceConverter) extractConditionsForQualifiers(extensions []*jaxb.ExtensionType) []*tslmodel.ConditionForQualifiers {
	var conditionsForQualifiersList []*tslmodel.ConditionForQualifiers
	for _, extensionType := range extensions {
		for _, item := range extensionType.Items {
			if qualifications, ok := item.Elem.(*jaxb.QualificationsType); ok {
				conditionForQualifiers := c.toConditionForQualificationsType(qualifications, extensionType.Critical)
				conditionsForQualifiersList = append(conditionsForQualifiersList, conditionForQualifiers...)
			}
		}
	}
	return conditionsForQualifiersList
}

func (c *TrustServiceConverter) toConditionForQualificationsType(qt *jaxb.QualificationsType, critical bool) []*tslmodel.ConditionForQualifiers {
	var conditionForQualifiers []*tslmodel.ConditionForQualifiers
	if qt != nil && utils.IsCollectionNotEmpty(qt.QualificationElement) {
		for _, qualificationElement := range qt.QualificationElement {
			if condition := c.toConditionForQualifiers(qualificationElement, critical); condition != nil {
				conditionForQualifiers = append(conditionForQualifiers, condition)
			}
		}
	}
	return conditionForQualifiers
}

func (c *TrustServiceConverter) toConditionForQualifiers(qualificationElement *jaxb.QualificationElementType, critical bool) *tslmodel.ConditionForQualifiers {
	qualifiers := c.extractQualifiers(qualificationElement)
	if utils.IsCollectionNotEmpty(qualifiers) {
		condition := NewCriteriaListConverter().Apply(qualificationElement.CriteriaList)
		return tslmodel.NewConditionForQualifiersWithCriticality(condition, qualifiers, critical)
	}
	return nil
}

func (c *TrustServiceConverter) extractAdditionalServiceInfoUris(extensions []*jaxb.ExtensionType) []string {
	var additionalServiceInfos []string
	for _, extensionType := range extensions {
		for _, item := range extensionType.Items {
			if additionalServiceInfo, ok := item.Elem.(*jaxb.AdditionalServiceInformationType); ok {
				uri := additionalServiceInfo.URI
				if uri != nil && utils.IsStringNotBlank(uri.Value) {
					additionalServiceInfos = append(additionalServiceInfos, uri.Value)
				}
			}
		}
	}
	return additionalServiceInfos
}

func (c *TrustServiceConverter) extractExpiredCertsRevocationInfo(extensions []*jaxb.ExtensionType) time.Time {
	for _, extensionType := range extensions {
		for _, item := range extensionType.Items {
			if lexical, ok := item.Elem.(*string); ok && item.ElemName.Local == "ExpiredCertsRevocationInfo" {
				return abstractParsingTaskConvertToDate(lexical)
			}
		}
	}
	return time.Time{}
}

func (c *TrustServiceConverter) extractQualifiers(qualificationElement *jaxb.QualificationElementType) []string {
	var qualifiers []string
	qualifiersType := qualificationElement.Qualifiers
	if qualifiersType != nil && utils.IsCollectionNotEmpty(qualifiersType.Qualifier) {
		for _, qualifierType := range qualifiersType.Qualifier {
			if qualifierType.URI != nil {
				qualifiers = append(qualifiers, *qualifierType.URI)
			} else {
				qualifiers = append(qualifiers, "")
			}
		}
	}
	return qualifiers
}

func (c *TrustServiceConverter) serviceSupplyPoints(serviceSupplyPoints *jaxb.ServiceSupplyPointsType) []string {
	var result []string
	if serviceSupplyPoints != nil && utils.IsCollectionNotEmpty(serviceSupplyPoints.ServiceSupplyPoint) {
		for _, nonEmptyURI := range serviceSupplyPoints.ServiceSupplyPoint {
			result = append(result, nonEmptyURI.Value)
		}
	}
	return result
}

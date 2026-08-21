// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/TrustServiceProviderConverter.java (DSS 6.5.RC1).
package tsl

import (
	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/trustedlist/jaxb"
	"github.com/utain/esig/dss/utils"
)

// TrustServiceProviderConverter converts a TSPType to a TrustServiceProvider.
type TrustServiceProviderConverter struct {
	// territory is the country code.
	territory string
}

// NewTrustServiceProviderConverter is the default constructor with null territory country code.
// Port of TrustServiceProviderConverter().
func NewTrustServiceProviderConverter() *TrustServiceProviderConverter {
	return &TrustServiceProviderConverter{}
}

// SetTerritory sets the territory. Port of setTerritory(String).
func (c *TrustServiceProviderConverter) SetTerritory(territory string) *TrustServiceProviderConverter {
	c.territory = territory
	return c
}

// Apply ports apply(TSPType).
func (c *TrustServiceProviderConverter) Apply(original *jaxb.TSPType) *tslmodel.TrustServiceProvider {
	tspBuilder := NewTrustServiceProviderBuilder()

	c.extractTSPInfo(tspBuilder, original.TSPInformation)
	tspBuilder.SetServices(c.extractTrustServices(original.TSPServices))

	return tspBuilder.Build()
}

func (c *TrustServiceProviderConverter) extractTSPInfo(tspBuilder *TrustServiceProviderBuilder, tspInformation *jaxb.TSPInformationType) {
	if tspInformation == nil {
		return
	}
	tspBuilder.SetTerritory(c.territory)

	converter := NewInternationalNamesTypeConverter()
	tspBuilder.SetNames(converter.Apply(tspInformation.TSPName))

	converter = NewInternationalNamesTypeConverterWithPredicate(NewTradeNamePredicate()) // filter registration identifiers
	tspBuilder.SetTradeNames(converter.Apply(tspInformation.TSPTradeName))

	tspBuilder.SetRegistrationIdentifiers(c.extractRegistrationIdentifiers(tspInformation.TSPTradeName))

	tspAddress := tspInformation.TSPAddress
	if tspAddress != nil {
		tspBuilder.SetPostalAddresses(c.extractPostalAddress(tspAddress.PostalAddresses))
		tspBuilder.SetElectronicAddresses(c.extractElectronicAddress(tspAddress.ElectronicAddress))
	}

	tspBuilder.SetInformation(c.extractInformationURI(tspInformation.TSPInformationURI))
}

func (c *TrustServiceProviderConverter) extractRegistrationIdentifiers(internationalNamesType *jaxb.InternationalNamesType) []string {
	predicate := NewOfficialRegistrationIdentifierPredicate()

	var result []string
	if internationalNamesType != nil && utils.IsCollectionNotEmpty(internationalNamesType.Name) {
		for _, multiLangNormString := range internationalNamesType.Name {
			value := multiLangNormString.Value
			if predicate.Test(value) && !stringSliceContains(result, value) {
				result = append(result, value)
			}
		}
	}
	return result
}

func (c *TrustServiceProviderConverter) extractPostalAddress(postalAddressList *jaxb.PostalAddressListType) map[string]string {
	result := make(map[string]string)
	if postalAddressList != nil && utils.IsCollectionNotEmpty(postalAddressList.PostalAddress) {
		for _, postalAddress := range postalAddressList.PostalAddress {
			lang := postalAddress.Lang
			if _, exists := result[lang]; !exists {
				result[lang] = c.getPostalAddress(postalAddress)
			}
		}
	}
	return result
}

func (c *TrustServiceProviderConverter) getPostalAddress(postalAddress *jaxb.PostalAddressType) string {
	var sb string
	if utils.IsStringNotEmpty(postalAddress.StreetAddress) {
		sb += postalAddress.StreetAddress + ", "
	}
	if postalAddress.PostalCode != nil && utils.IsStringNotEmpty(*postalAddress.PostalCode) {
		sb += *postalAddress.PostalCode + ", "
	}
	if utils.IsStringNotEmpty(postalAddress.Locality) {
		sb += postalAddress.Locality + ", "
	}
	if postalAddress.StateOrProvince != nil && utils.IsStringNotEmpty(*postalAddress.StateOrProvince) {
		sb += *postalAddress.StateOrProvince + ", "
	}
	if utils.IsStringNotEmpty(postalAddress.CountryName) {
		sb += postalAddress.CountryName
	}
	return sb
}

func (c *TrustServiceProviderConverter) extractElectronicAddress(electronicAddress *jaxb.ElectronicAddressType) map[string][]string {
	result := make(map[string][]string)
	if electronicAddress != nil && utils.IsCollectionNotEmpty(electronicAddress.URI) {
		for _, uriAndLang := range electronicAddress.URI {
			c.addEntry(result, uriAndLang.Lang, uriAndLang.Value)
		}
	}
	return result
}

func (c *TrustServiceProviderConverter) extractInformationURI(tspInformationURI *jaxb.NonEmptyMultiLangURIListType) map[string]string {
	result := make(map[string]string)
	if tspInformationURI != nil && utils.IsCollectionNotEmpty(tspInformationURI.URI) {
		for _, uriAndLang := range tspInformationURI.URI {
			lang := uriAndLang.Lang
			if _, exists := result[lang]; !exists {
				result[lang] = uriAndLang.Value
			}
		}
	}
	return result
}

func (c *TrustServiceProviderConverter) addEntry(result map[string][]string, lang, value string) {
	result[lang] = append(result[lang], value)
}

func (c *TrustServiceProviderConverter) extractTrustServices(tspServicesList *jaxb.TSPServicesListType) []*tslmodel.TrustService {
	if tspServicesList == nil || !utils.IsCollectionNotEmpty(tspServicesList.TSPService) {
		return nil
	}
	converter := NewTrustServiceConverter()
	result := make([]*tslmodel.TrustService, 0, len(tspServicesList.TSPService))
	for _, tspService := range tspServicesList.TSPService {
		result = append(result, converter.Apply(tspService))
	}
	return result
}

// stringSliceContains reports whether s contains value.
func stringSliceContains(s []string, value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

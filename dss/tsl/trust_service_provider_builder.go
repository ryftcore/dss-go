// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sync/TrustServiceProviderBuilder.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
)

// TrustServiceProviderBuilder builds a TrustServiceProvider.
type TrustServiceProviderBuilder struct {
	// names is the map of names (the key is the language).
	names map[string][]string
	// tradeNames is the map of trade names.
	tradeNames map[string][]string
	// registrationIdentifiers is the list of registration identifiers.
	registrationIdentifiers []string
	// postalAddresses is the map of postal addresses.
	postalAddresses map[string]string
	// electronicAddresses is the map of electronic addresses.
	electronicAddresses map[string][]string
	// information is the map of information.
	information map[string]string
	// services is the list of trust services.
	services []*tslmodel.TrustService
	// territory is the territory (country).
	territory string
}

// NewTrustServiceProviderBuilder is the default constructor. Port of TrustServiceProviderBuilder().
func NewTrustServiceProviderBuilder() *TrustServiceProviderBuilder {
	return &TrustServiceProviderBuilder{}
}

// NewTrustServiceProviderBuilderFromOriginal copies the original object. Port of
// TrustServiceProviderBuilder(TrustServiceProvider).
//
// Panics with the Java message when original is nil (Objects.requireNonNull).
func NewTrustServiceProviderBuilderFromOriginal(original *tslmodel.TrustServiceProvider) *TrustServiceProviderBuilder {
	if original == nil {
		panic("TrustServiceProvider cannot be null!")
	}
	return &TrustServiceProviderBuilder{
		names:                   original.Names(),
		tradeNames:              original.TradeNames(),
		registrationIdentifiers: original.RegistrationIdentifiers(),
		postalAddresses:         original.PostalAddresses(),
		electronicAddresses:     original.ElectronicAddresses(),
		information:             original.Information(),
		services:                original.Services(),
		territory:               original.Territory(),
	}
}

// Build builds a TrustServiceProvider. Port of build().
func (b *TrustServiceProviderBuilder) Build() *tslmodel.TrustServiceProvider {
	trustServiceProvider := tslmodel.NewTrustServiceProvider()
	trustServiceProvider.SetNames(b.Names())
	trustServiceProvider.SetTradeNames(b.TradeNames())
	trustServiceProvider.SetRegistrationIdentifiers(b.RegistrationIdentifiers())
	trustServiceProvider.SetPostalAddresses(b.PostalAddresses())
	trustServiceProvider.SetElectronicAddresses(b.ElectronicAddresses())
	trustServiceProvider.SetInformation(b.Information())
	trustServiceProvider.SetServices(b.Services())
	trustServiceProvider.SetTerritory(b.Territory())
	return trustServiceProvider
}

// Names gets a map of names (first key is the language). Port of getNames().
func (b *TrustServiceProviderBuilder) Names() map[string][]string {
	return unmodifiableMapWithLists(b.names)
}

// SetNames sets a map of names. Port of setNames(Map).
func (b *TrustServiceProviderBuilder) SetNames(names map[string][]string) *TrustServiceProviderBuilder {
	b.names = names
	return b
}

// TradeNames gets a map of trade names. Port of getTradeNames().
func (b *TrustServiceProviderBuilder) TradeNames() map[string][]string {
	return unmodifiableMapWithLists(b.tradeNames)
}

// SetTradeNames sets a map of trade names. Port of setTradeNames(Map).
func (b *TrustServiceProviderBuilder) SetTradeNames(tradeNames map[string][]string) *TrustServiceProviderBuilder {
	b.tradeNames = tradeNames
	return b
}

// RegistrationIdentifiers gets registration identifiers. Port of getRegistrationIdentifiers().
func (b *TrustServiceProviderBuilder) RegistrationIdentifiers() []string {
	return unmodifiableStringList(b.registrationIdentifiers)
}

// SetRegistrationIdentifiers sets registration identifiers. Port of
// setRegistrationIdentifiers(List).
func (b *TrustServiceProviderBuilder) SetRegistrationIdentifiers(registrationIdentifiers []string) *TrustServiceProviderBuilder {
	b.registrationIdentifiers = registrationIdentifiers
	return b
}

// PostalAddresses gets a map of postal addresses. Port of getPostalAddresses().
func (b *TrustServiceProviderBuilder) PostalAddresses() map[string]string {
	return unmodifiableStringMap(b.postalAddresses)
}

// SetPostalAddresses sets a map of postal addresses. Port of setPostalAddresses(Map).
func (b *TrustServiceProviderBuilder) SetPostalAddresses(postalAddresses map[string]string) *TrustServiceProviderBuilder {
	b.postalAddresses = postalAddresses
	return b
}

// ElectronicAddresses gets a map of electronic addresses. Port of getElectronicAddresses().
func (b *TrustServiceProviderBuilder) ElectronicAddresses() map[string][]string {
	return unmodifiableMapWithLists(b.electronicAddresses)
}

// SetElectronicAddresses sets a map of electronic addresses. Port of setElectronicAddresses(Map).
func (b *TrustServiceProviderBuilder) SetElectronicAddresses(electronicAddresses map[string][]string) *TrustServiceProviderBuilder {
	b.electronicAddresses = electronicAddresses
	return b
}

// Information gets a map of information. Port of getInformation().
func (b *TrustServiceProviderBuilder) Information() map[string]string {
	return unmodifiableStringMap(b.information)
}

// SetInformation sets a map of information. Port of setInformation(Map).
func (b *TrustServiceProviderBuilder) SetInformation(information map[string]string) *TrustServiceProviderBuilder {
	b.information = information
	return b
}

// Services gets a list of trust services. Port of getServices().
func (b *TrustServiceProviderBuilder) Services() []*tslmodel.TrustService {
	return b.unmodifiableTrustServices(b.services)
}

// SetServices sets a list of trust services. Port of setServices(List).
func (b *TrustServiceProviderBuilder) SetServices(services []*tslmodel.TrustService) *TrustServiceProviderBuilder {
	b.services = services
	return b
}

// Territory gets territory (country). Port of getTerritory().
func (b *TrustServiceProviderBuilder) Territory() string {
	return b.territory
}

// SetTerritory sets territory (country). Port of setTerritory(String).
func (b *TrustServiceProviderBuilder) SetTerritory(territory string) *TrustServiceProviderBuilder {
	b.territory = territory
	return b
}

// unmodifiableStringList ports the private <T> getUnmodifiableList(List<T>) specialized to
// string, returning a defensive copy (never nil, mirroring Collections.unmodifiableList of a
// freshly built ArrayList).
func unmodifiableStringList(originalList []string) []string {
	newList := make([]string, 0, len(originalList))
	newList = append(newList, originalList...)
	return newList
}

// unmodifiableStringMap ports the private <T,K> getUnmodifiableMap(Map<T,K>) specialized to
// string keys/values.
func unmodifiableStringMap(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string, len(originalMap))
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}

// unmodifiableMapWithLists ports the private getUnmodifiableMapWithLists(Map).
func unmodifiableMapWithLists(originalMap map[string][]string) map[string][]string {
	copyMap := make(map[string][]string, len(originalMap))
	for k, v := range originalMap {
		copyMap[k] = unmodifiableStringList(v)
	}
	return copyMap
}

// unmodifiableTrustServices ports the private getUnmodifiableTrustServices(List<TrustService>).
func (b *TrustServiceProviderBuilder) unmodifiableTrustServices(originalTrustServices []*tslmodel.TrustService) []*tslmodel.TrustService {
	copyTrustServices := make([]*tslmodel.TrustService, 0, len(originalTrustServices))
	for _, trustService := range originalTrustServices {
		trustServiceBuilder := tslmodel.NewTrustServiceBuilder()
		copyTrustService := trustServiceBuilder.
			SetCertificates(unmodifiableCertificateTokenList(trustService.Certificates())).
			SetStatusAndInformationExtensions(b.unmodifiableTimeDependentValues(trustService.StatusAndInformationExtensions())).
			Build()
		copyTrustServices = append(copyTrustServices, copyTrustService)
	}
	return copyTrustServices
}

// unmodifiableTimeDependentValues ports the private
// getUnmodifiableTimeDependentValues(Values<TrustServiceStatusAndInformationExtensions>).
func (b *TrustServiceProviderBuilder) unmodifiableTimeDependentValues(
	timeDependentValues *timedependent.Values[*tslmodel.TrustServiceStatusAndInformationExtensions],
) *timedependent.Values[*tslmodel.TrustServiceStatusAndInformationExtensions] {
	var copyTSSAndIEs []*tslmodel.TrustServiceStatusAndInformationExtensions

	for status := range timeDependentValues.Iterator() {
		builder := tslmodel.NewTrustServiceStatusAndInformationExtensionsBuilder()
		copyStatus := builder.
			SetNames(unmodifiableMapWithLists(status.Names())).
			SetType(status.Type()).
			SetStatus(status.Status()).
			SetConditionsForQualifiers(unmodifiableConditionForQualifiersList(status.ConditionsForQualifiers())).
			SetAdditionalServiceInfoUris(unmodifiableStringList(status.AdditionalServiceInfoUris())).
			SetServiceSupplyPoints(unmodifiableStringList(status.ServiceSupplyPoints())).
			SetExpiredCertsRevocationInfo(status.ExpiredCertsRevocationInfo()).
			SetStartDate(status.StartDate()).
			SetEndDate(status.EndDate()).
			Build()
		copyTSSAndIEs = append(copyTSSAndIEs, copyStatus)
	}

	return timedependent.NewValuesFrom(copyTSSAndIEs)
}

// unmodifiableCertificateTokenList returns a defensive copy of a CertificateToken slice.
func unmodifiableCertificateTokenList(original []*model.CertificateToken) []*model.CertificateToken {
	copyList := make([]*model.CertificateToken, 0, len(original))
	copyList = append(copyList, original...)
	return copyList
}

// unmodifiableConditionForQualifiersList returns a defensive copy of a ConditionForQualifiers
// slice.
func unmodifiableConditionForQualifiersList(original []*tslmodel.ConditionForQualifiers) []*tslmodel.ConditionForQualifiers {
	copyList := make([]*tslmodel.ConditionForQualifiers, 0, len(original))
	copyList = append(copyList, original...)
	return copyList
}

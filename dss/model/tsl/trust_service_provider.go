// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustServiceProvider.java (DSS 6.5.RC1).
package tsl

// TrustServiceProvider is a DTO representation for a trust service provider.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustServiceProvider struct {
	// names maps lang -> values for that lang.
	names map[string][]string
	// tradeNames maps lang -> values for that lang.
	tradeNames map[string][]string
	// registrationIdentifiers is the list of registration identifiers.
	registrationIdentifiers []string
	// postalAddresses maps lang -> value.
	postalAddresses map[string]string
	// electronicAddresses maps lang -> values for that lang.
	electronicAddresses map[string][]string
	// information maps lang -> value.
	information map[string]string
	// services is the list of trust services.
	services []*TrustService
	// territory is the territory (country).
	territory string
}

// NewTrustServiceProvider instantiates an object with zero values. Port of the default
// constructor.
func NewTrustServiceProvider() *TrustServiceProvider {
	return &TrustServiceProvider{}
}

// Names gets a map of names.
func (t *TrustServiceProvider) Names() map[string][]string {
	return t.names
}

// SetNames sets a map of names.
func (t *TrustServiceProvider) SetNames(names map[string][]string) {
	t.names = names
}

// TradeNames gets a map of trade names.
func (t *TrustServiceProvider) TradeNames() map[string][]string {
	return t.tradeNames
}

// SetTradeNames sets a map of trade names.
func (t *TrustServiceProvider) SetTradeNames(tradeNames map[string][]string) {
	t.tradeNames = tradeNames
}

// RegistrationIdentifiers gets a list of registration identifiers.
func (t *TrustServiceProvider) RegistrationIdentifiers() []string {
	return t.registrationIdentifiers
}

// SetRegistrationIdentifiers sets a list of registration identifiers.
func (t *TrustServiceProvider) SetRegistrationIdentifiers(registrationIdentifiers []string) {
	t.registrationIdentifiers = registrationIdentifiers
}

// PostalAddresses gets a map of postal addresses.
func (t *TrustServiceProvider) PostalAddresses() map[string]string {
	return t.postalAddresses
}

// SetPostalAddresses sets a map of postal addresses.
func (t *TrustServiceProvider) SetPostalAddresses(postalAddresses map[string]string) {
	t.postalAddresses = postalAddresses
}

// ElectronicAddresses gets a map of electronic addresses.
func (t *TrustServiceProvider) ElectronicAddresses() map[string][]string {
	return t.electronicAddresses
}

// SetElectronicAddresses sets a map of electronic addresses.
func (t *TrustServiceProvider) SetElectronicAddresses(electronicAddresses map[string][]string) {
	t.electronicAddresses = electronicAddresses
}

// Information gets a map of information.
func (t *TrustServiceProvider) Information() map[string]string {
	return t.information
}

// SetInformation sets a map of information.
func (t *TrustServiceProvider) SetInformation(information map[string]string) {
	t.information = information
}

// Services gets a list of trust services.
func (t *TrustServiceProvider) Services() []*TrustService {
	return t.services
}

// SetServices sets a list of trust services.
func (t *TrustServiceProvider) SetServices(services []*TrustService) {
	t.services = services
}

// Territory gets the territory (country).
func (t *TrustServiceProvider) Territory() string {
	return t.territory
}

// SetTerritory sets the territory (country).
func (t *TrustServiceProvider) SetTerritory(territory string) {
	t.territory = territory
}

// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/TrustedEntity.java (DSS 6.5.RC1).
package lote

// TrustedEntity is a DTO containing information extracted for a trusted entity (TS 119 602).
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustedEntity struct {
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
	// information maps lang -> values for that lang.
	information map[string][]string
	// services is the list of trusted entity services.
	services []*TrustedEntityService
	// territory is the territory (country).
	territory string
}

// NewTrustedEntity instantiates an object with zero values. Port of the default constructor.
func NewTrustedEntity() *TrustedEntity {
	return &TrustedEntity{}
}

// Names gets a map of names.
func (t *TrustedEntity) Names() map[string][]string {
	return t.names
}

// SetNames sets a map of names.
func (t *TrustedEntity) SetNames(names map[string][]string) {
	t.names = names
}

// TradeNames gets a map of trade names.
func (t *TrustedEntity) TradeNames() map[string][]string {
	return t.tradeNames
}

// SetTradeNames sets a map of trade names.
func (t *TrustedEntity) SetTradeNames(tradeNames map[string][]string) {
	t.tradeNames = tradeNames
}

// RegistrationIdentifiers gets a list of registration identifiers.
func (t *TrustedEntity) RegistrationIdentifiers() []string {
	return t.registrationIdentifiers
}

// SetRegistrationIdentifiers sets a list of registration identifiers.
func (t *TrustedEntity) SetRegistrationIdentifiers(registrationIdentifiers []string) {
	t.registrationIdentifiers = registrationIdentifiers
}

// PostalAddresses gets a map of postal addresses.
func (t *TrustedEntity) PostalAddresses() map[string]string {
	return t.postalAddresses
}

// SetPostalAddresses sets a map of postal addresses.
func (t *TrustedEntity) SetPostalAddresses(postalAddresses map[string]string) {
	t.postalAddresses = postalAddresses
}

// ElectronicAddresses gets a map of electronic addresses.
func (t *TrustedEntity) ElectronicAddresses() map[string][]string {
	return t.electronicAddresses
}

// SetElectronicAddresses sets a map of electronic addresses.
func (t *TrustedEntity) SetElectronicAddresses(electronicAddresses map[string][]string) {
	t.electronicAddresses = electronicAddresses
}

// Information gets a map of information.
func (t *TrustedEntity) Information() map[string][]string {
	return t.information
}

// SetInformation sets a map of information.
func (t *TrustedEntity) SetInformation(information map[string][]string) {
	t.information = information
}

// Services gets a list of trusted entity services.
func (t *TrustedEntity) Services() []*TrustedEntityService {
	return t.services
}

// SetServices sets a list of trusted entity services.
func (t *TrustedEntity) SetServices(services []*TrustedEntityService) {
	t.services = services
}

// Territory gets territory (country).
func (t *TrustedEntity) Territory() string {
	return t.territory
}

// SetTerritory sets territory (country).
func (t *TrustedEntity) SetTerritory(territory string) {
	t.territory = territory
}

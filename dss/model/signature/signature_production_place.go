// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureProductionPlace.java (DSS 6.5.RC1).
package signature

// ProductionPlace represents the information concerning the signature production
// place.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type ProductionPlace struct {
	// city is the location (city).
	city string

	// stateOrProvince is the region (stateOrProvince).
	stateOrProvince string

	// postOfficeBoxNumber is the postOfficeBoxNumber.
	postOfficeBoxNumber string

	// postalCode is the postalCode.
	postalCode string

	// countryName is the countryName (can be 2-letters abbreviation, e.g. LU for Luxembourg).
	countryName string

	// streetAddress is the address.
	streetAddress string

	// postalAddress is the postal address (used in CAdES).
	postalAddress []string
}

// NewSignatureProductionPlace is the default constructor instantiating the object with null
// (zero) values.
func NewSignatureProductionPlace() *ProductionPlace {
	return &ProductionPlace{}
}

// City gets location (city). Port of getCity().
func (s *ProductionPlace) City() string {
	return s.city
}

// SetCity sets location (city). Port of setCity(String).
func (s *ProductionPlace) SetCity(city string) {
	s.city = city
}

// StateOrProvince gets region (stateOrProvince). Port of getStateOrProvince().
func (s *ProductionPlace) StateOrProvince() string {
	return s.stateOrProvince
}

// SetStateOrProvince sets region (stateOrProvince). Port of setStateOrProvince(String).
func (s *ProductionPlace) SetStateOrProvince(stateOrProvince string) {
	s.stateOrProvince = stateOrProvince
}

// PostOfficeBoxNumber gets postOfficeBoxNumber. Port of getPostOfficeBoxNumber().
func (s *ProductionPlace) PostOfficeBoxNumber() string {
	return s.postOfficeBoxNumber
}

// SetPostOfficeBoxNumber sets postOfficeBoxNumber. Port of
// setPostOfficeBoxNumber(String).
func (s *ProductionPlace) SetPostOfficeBoxNumber(postOfficeBoxNumber string) {
	s.postOfficeBoxNumber = postOfficeBoxNumber
}

// PostalCode gets postal code. Port of getPostalCode().
func (s *ProductionPlace) PostalCode() string {
	return s.postalCode
}

// SetPostalCode sets postal code. Port of setPostalCode(String).
func (s *ProductionPlace) SetPostalCode(postalCode string) {
	s.postalCode = postalCode
}

// CountryName gets country name. Port of getCountryName().
func (s *ProductionPlace) CountryName() string {
	return s.countryName
}

// SetCountryName sets country name (can be 2-letters abbreviation, e.g. LU for Luxembourg).
// Port of setCountryName(String).
func (s *ProductionPlace) SetCountryName(countryName string) {
	s.countryName = countryName
}

// StreetAddress gets the address. Port of getStreetAddress().
func (s *ProductionPlace) StreetAddress() string {
	return s.streetAddress
}

// SetStreetAddress sets the address. Port of setStreetAddress(String).
func (s *ProductionPlace) SetStreetAddress(streetAddress string) {
	s.streetAddress = streetAddress
}

// PostalAddress gets postal address (used in CAdES), lazily initialising the backing slice
// as Java lazily initialises its ArrayList. Port of getPostalAddress().
func (s *ProductionPlace) PostalAddress() []string {
	if s.postalAddress == nil {
		s.postalAddress = make([]string, 0)
	}
	return s.postalAddress
}

// SetPostalAddress sets postal address (used in CAdES). Port of
// setPostalAddress(List<String>).
func (s *ProductionPlace) SetPostalAddress(postalAddress []string) {
	s.postalAddress = postalAddress
}

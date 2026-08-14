// Ported from dss-model/.../SignerLocation.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"
)

// SignerLocation represents the information concerning the signature
// production place.
type SignerLocation struct {
	// postalAddress is a sequence defined a Postal Address. NOTE: used in
	// CAdES.
	postalAddress []string

	// postOfficeBoxNumber is the post office box number for PO box
	// addresses. NOTE: used in JAdES.
	postOfficeBoxNumber string

	// postalCode is the postal code (ZIP-code). For example, 94043.
	postalCode string

	// locality is the locality (city) in which the street address is, and
	// which is in the region.
	locality string

	// stateOrProvince is the state or province. The region in which the
	// locality is, and which is in the country.
	stateOrProvince string

	// country is the country. For example, USA. You can also provide the
	// two-letter ISO 3166-1 alpha-2 country code.
	country string

	// streetAddress is the street address. For example, 1600 Amphitheatre
	// Pkwy. NOTE: used in XAdES and JAdES.
	streetAddress string
}

// NewSignerLocation creates the default SignerLocation. Ports the default
// constructor.
func NewSignerLocation() *SignerLocation {
	return &SignerLocation{postalAddress: []string{}}
}

// Country gets the country.
func (s *SignerLocation) Country() string { return s.country }

// SetCountry sets the country. Can be a country name or its two-letter
// ISO 3166-1 alpha-2 country code.
func (s *SignerLocation) SetCountry(country string) { s.country = country }

// Locality gets the locality (city).
func (s *SignerLocation) Locality() string { return s.locality }

// SetLocality sets the locality (city).
func (s *SignerLocation) SetLocality(locality string) { s.locality = locality }

// PostalAddress gets the postal address.
func (s *SignerLocation) PostalAddress() []string { return s.postalAddress }

// SetPostalAddress sets the postal address. NOTE: used in CAdES.
func (s *SignerLocation) SetPostalAddress(postalAddress []string) { s.postalAddress = postalAddress }

// AddPostalAddress adds an address item to the complete address. NOTE:
// used in CAdES.
func (s *SignerLocation) AddPostalAddress(addressItem string) {
	if s.postalAddress == nil {
		s.postalAddress = []string{}
	}
	s.postalAddress = append(s.postalAddress, addressItem)
}

// PostalCode gets the postal code.
func (s *SignerLocation) PostalCode() string { return s.postalCode }

// SetPostalCode sets the postal code.
func (s *SignerLocation) SetPostalCode(postalCode string) { s.postalCode = postalCode }

// PostOfficeBoxNumber gets the post office box number.
func (s *SignerLocation) PostOfficeBoxNumber() string { return s.postOfficeBoxNumber }

// SetPostOfficeBoxNumber sets the post office box number. NOTE: used in
// JAdES.
func (s *SignerLocation) SetPostOfficeBoxNumber(postOfficeBoxNumber string) {
	s.postOfficeBoxNumber = postOfficeBoxNumber
}

// StateOrProvince gets the state or province.
func (s *SignerLocation) StateOrProvince() string { return s.stateOrProvince }

// SetStateOrProvince sets the state or province (the region where the
// locality is).
func (s *SignerLocation) SetStateOrProvince(stateOrProvince string) {
	s.stateOrProvince = stateOrProvince
}

// StreetAddress gets the street address.
func (s *SignerLocation) StreetAddress() string { return s.streetAddress }

// SetStreetAddress sets the street address. NOTE: used in XAdES and
// JAdES.
func (s *SignerLocation) SetStreetAddress(streetAddress string) { s.streetAddress = streetAddress }

// IsEmpty checks if the SignerLocation instance is empty.
func (s *SignerLocation) IsEmpty() bool {
	if len(s.postalAddress) > 0 {
		return false
	}
	if s.postalCode != "" {
		return false
	}
	if s.postOfficeBoxNumber != "" {
		return false
	}
	if s.locality != "" {
		return false
	}
	if s.stateOrProvince != "" {
		return false
	}
	if s.country != "" {
		return false
	}
	if s.streetAddress != "" {
		return false
	}
	return true
}

// Equals ports SignerLocation#equals.
func (s *SignerLocation) Equals(other *SignerLocation) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	return s.country == other.country &&
		s.locality == other.locality &&
		s.postOfficeBoxNumber == other.postOfficeBoxNumber &&
		reflect.DeepEqual(s.postalAddress, other.postalAddress) &&
		s.postalCode == other.postalCode &&
		s.stateOrProvince == other.stateOrProvince &&
		s.streetAddress == other.streetAddress
}

// String ports SignerLocation#toString.
func (s *SignerLocation) String() string {
	return fmt.Sprintf("SignerLocation [postalAddress=%v, postOfficeBoxNumber=%s, postalCode=%s, locality=%s, stateOrProvince=%s, country=%s, street=%s]",
		s.postalAddress, s.postOfficeBoxNumber, s.postalCode, s.locality, s.stateOrProvince, s.country, s.streetAddress)
}

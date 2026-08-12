// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/QCPSB.java (DSS 6.5.RC1).
package extension

// QCPSB defines a Public Sector Body's Electronic Attestation of Attributes (PSBEAA)
// provider certificate.
type QCPSB struct {
	// countryOfLegislation is:
	//
	//	countryOfLegislation PrintableString (SIZE (2))
	//	(CONSTRAINED BY { -- ISO 3166 alpha-2 country codes or 'EU' -- }),
	//	 -- this field shall contain the alpha-2 country code of the legislation framework of
	//	public sector body
	//	 -- In the case of European Union law 'EU' shall be used in place of the country code
	countryOfLegislation string

	// authSourceIdentification is:
	//
	//	authSourceIdentification UTF8String,
	//	-- this field is for the unique identification of authentic source
	authSourceIdentification string

	// legislationIdentification is:
	//
	//	legislationIdentification UTF8String
	legislationIdentification string
}

// NewQCPSB instantiates the object with null values. Ports the default constructor.
func NewQCPSB() *QCPSB {
	return &QCPSB{}
}

// CountryOfLegislation gets the country of legislation. The value shall be represented by
// a two-letter ISO 3166 alpha-2 country code.
func (q *QCPSB) CountryOfLegislation() string {
	return q.countryOfLegislation
}

// SetCountryOfLegislation sets the country of legislation. The value shall be represented
// by a two-letter ISO 3166 alpha-2 country code.
func (q *QCPSB) SetCountryOfLegislation(countryOfLegislation string) {
	q.countryOfLegislation = countryOfLegislation
}

// AuthSourceIdentification gets the authentic source identification.
func (q *QCPSB) AuthSourceIdentification() string {
	return q.authSourceIdentification
}

// SetAuthSourceIdentification sets the authentic source identification.
func (q *QCPSB) SetAuthSourceIdentification(authSourceIdentification string) {
	q.authSourceIdentification = authSourceIdentification
}

// LegislationIdentification gets the legislation identification.
func (q *QCPSB) LegislationIdentification() string {
	return q.legislationIdentification
}

// SetLegislationIdentification sets the legislation identification.
func (q *QCPSB) SetLegislationIdentification(legislationIdentification string) {
	q.legislationIdentification = legislationIdentification
}

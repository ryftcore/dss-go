// Ported from dss-enumerations/.../PKIEncoding.java (DSS 6.5.RC1).
//
// Enumeration with the possible encoding for PKI encapsulation.
// ETSI EN 319 132-1 5.1.3
package enumerations

// PKIEncoding implements UriBasedEnum.
type PKIEncoding string

const (
	// PKIEncodingDER is http://uri.etsi.org/01903/v1.2.2#DER
	PKIEncodingDER PKIEncoding = "DER"
	// PKIEncodingBER is http://uri.etsi.org/01903/v1.2.2#BER
	PKIEncodingBER PKIEncoding = "BER"
	// PKIEncodingCER is http://uri.etsi.org/01903/v1.2.2#CER
	PKIEncodingCER PKIEncoding = "CER"
	// PKIEncodingPER is http://uri.etsi.org/01903/v1.2.2#PER
	PKIEncodingPER PKIEncoding = "PER"
	// PKIEncodingXER is http://uri.etsi.org/01903/v1.2.2#XER
	PKIEncodingXER PKIEncoding = "XER"
)

// pkiEncodingURIs holds the URI for each constant.
var pkiEncodingURIs = map[PKIEncoding]string{
	PKIEncodingDER: "http://uri.etsi.org/01903/v1.2.2#DER",
	PKIEncodingBER: "http://uri.etsi.org/01903/v1.2.2#BER",
	PKIEncodingCER: "http://uri.etsi.org/01903/v1.2.2#CER",
	PKIEncodingPER: "http://uri.etsi.org/01903/v1.2.2#PER",
	PKIEncodingXER: "http://uri.etsi.org/01903/v1.2.2#XER",
}

// PKIEncodingValues returns all constants in declaration order.
func PKIEncodingValues() []PKIEncoding {
	return []PKIEncoding{
		PKIEncodingDER,
		PKIEncodingBER,
		PKIEncodingCER,
		PKIEncodingPER,
		PKIEncodingXER,
	}
}

// URI returns the encoding URI. Implements UriBasedEnum.
func (p PKIEncoding) URI() string {
	return pkiEncodingURIs[p]
}

// compile-time interface assertion.
var _ UriBasedEnum = PKIEncoding("")

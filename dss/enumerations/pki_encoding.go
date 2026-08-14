// Ported from dss-enumerations/.../PKIEncoding.java (DSS 6.5.RC1).
//
// Enumeration with the possible encoding for PKI encapsulation.
// ETSI EN 319 132-1 5.1.3
package enumerations

// PKIEncoding implements UriBasedEnum.
type PKIEncoding string

const (
	// PKIEncoding_DER is http://uri.etsi.org/01903/v1.2.2#DER
	PKIEncoding_DER PKIEncoding = "DER"
	// PKIEncoding_BER is http://uri.etsi.org/01903/v1.2.2#BER
	PKIEncoding_BER PKIEncoding = "BER"
	// PKIEncoding_CER is http://uri.etsi.org/01903/v1.2.2#CER
	PKIEncoding_CER PKIEncoding = "CER"
	// PKIEncoding_PER is http://uri.etsi.org/01903/v1.2.2#PER
	PKIEncoding_PER PKIEncoding = "PER"
	// PKIEncoding_XER is http://uri.etsi.org/01903/v1.2.2#XER
	PKIEncoding_XER PKIEncoding = "XER"
)

// pkiEncodingURIs holds the URI for each constant.
var pkiEncodingURIs = map[PKIEncoding]string{
	PKIEncoding_DER: "http://uri.etsi.org/01903/v1.2.2#DER",
	PKIEncoding_BER: "http://uri.etsi.org/01903/v1.2.2#BER",
	PKIEncoding_CER: "http://uri.etsi.org/01903/v1.2.2#CER",
	PKIEncoding_PER: "http://uri.etsi.org/01903/v1.2.2#PER",
	PKIEncoding_XER: "http://uri.etsi.org/01903/v1.2.2#XER",
}

// PKIEncodingValues returns all constants in declaration order.
func PKIEncodingValues() []PKIEncoding {
	return []PKIEncoding{
		PKIEncoding_DER,
		PKIEncoding_BER,
		PKIEncoding_CER,
		PKIEncoding_PER,
		PKIEncoding_XER,
	}
}

// URI returns the encoding URI. Implements UriBasedEnum.
func (p PKIEncoding) URI() string {
	return pkiEncodingURIs[p]
}

// compile-time interface assertion.
var _ UriBasedEnum = PKIEncoding("")

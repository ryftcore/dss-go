// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/PdsLocation.java (DSS 6.5.RC1).
package extension

// PdsLocation is:
//
//	PdsLocation::= SEQUENCE {
//	 url IA5String,
//	 language PrintableString (SIZE(2))} --ISO 639-1 language code
type PdsLocation struct {
	// url is the URL.
	url string

	// language is the language.
	language string
}

// NewPdsLocation instantiates the object with null values. Ports the default constructor.
func NewPdsLocation() *PdsLocation {
	return &PdsLocation{}
}

// Url returns the URL.
func (p *PdsLocation) Url() string {
	return p.url
}

// SetUrl sets the URL.
func (p *PdsLocation) SetUrl(url string) {
	p.url = url
}

// Language returns the language.
func (p *PdsLocation) Language() string {
	return p.language
}

// SetLanguage sets the language.
func (p *PdsLocation) SetLanguage(language string) {
	p.language = language
}

// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/TrustedListFacade.java (DSS 6.5.RC1).

package trustedlist

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// Facade performs marshalling/unmarshalling of a Trusted List
// XML, the way eu.europa.esig.trustedlist.TrustedListFacade does. Java's
// AbstractJaxbFacade base additionally validates against the XSD Schema
// TrustedListUtils.getSchema() builds; this port performs no runtime XSD
// validation (see trusted_list_utils.go's header), so Unmarshal/Marshal
// delegate to dss/trustedlist/jaxb directly.
type Facade struct{}

// NewTrustedListFacade creates a new facade, the way Facade.newFacade does.
func NewTrustedListFacade() *Facade {
	return &Facade{}
}

// Unmarshal parses a Trusted List document.
func (f *Facade) Unmarshal(data []byte) (*jaxb.TrustStatusListType, error) {
	return jaxb.Unmarshal(data)
}

// Marshal writes a Trusted List document, the way Facade's
// inherited marshall does (JAXB_FORMATTED_OUTPUT=true - see
// dss/trustedlist/jaxb's Marshal).
func (f *Facade) Marshal(tsl *jaxb.TrustStatusListType) ([]byte, error) {
	return jaxb.Marshal(tsl)
}

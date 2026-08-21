// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/mra/MRAFacade.java (DSS 6.5.RC1).

package trustedlist

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// MRAFacade performs marshalling/unmarshalling of a Trusted List XML with
// applied MRA scheme, the way eu.europa.esig.trustedlist.mra.MRAFacade
// (which extends TrustedListFacade, overriding its JAXBContext/Schema to
// additionally register mra.ObjectFactory) does. Unmarshal is identical to
// TrustedListFacade's - the model is the same regardless of which
// JAXBContext parsed it - Marshal instead writes MRAFacade's own root
// namespace set (see dss/trustedlist/jaxb's MarshalMRA).
type MRAFacade struct{}

// NewMRAFacade creates a new facade, the way MRAFacade.newFacade does.
func NewMRAFacade() *MRAFacade {
	return &MRAFacade{}
}

// Unmarshal parses a Trusted List document.
func (f *MRAFacade) Unmarshal(data []byte) (*jaxb.TrustStatusListType, error) {
	return jaxb.Unmarshal(data)
}

// Marshal writes a Trusted List document with MRAFacade's root namespace
// set.
func (f *MRAFacade) Marshal(tsl *jaxb.TrustStatusListType) ([]byte, error) {
	return jaxb.MarshalMRA(tsl)
}

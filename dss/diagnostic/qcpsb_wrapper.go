// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/QCPSBWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

// QCPSBWrapper provides a user-friendly API for dealing with jaxb.XmlQcPSB.
type QCPSBWrapper struct {
	// wrapped is the wrapped XmlQcPSB.
	wrapped *jaxb.XmlQcPSB
}

// NewQCPSBWrapper is the default constructor.
func NewQCPSBWrapper(xmlQcPSB *jaxb.XmlQcPSB) *QCPSBWrapper {
	return &QCPSBWrapper{wrapped: xmlQcPSB}
}

// CountryOfLegislation gets the two-letter code of the legislation country (ISO 3166
// alpha-2 country codes or 'EU'). Port of getCountryOfLegislation().
func (w *QCPSBWrapper) CountryOfLegislation() string {
	if w.wrapped.CountryOfLegislation != nil {
		return *w.wrapped.CountryOfLegislation
	}
	return ""
}

// AuthSourceIdentification gets the unique identification of authentic source. Port of
// getAuthSourceIdentification().
func (w *QCPSBWrapper) AuthSourceIdentification() string {
	if w.wrapped.AuthSourceIdentification != nil {
		return *w.wrapped.AuthSourceIdentification
	}
	return ""
}

// LegislationIdentification gets the legislation identification. Port of
// getLegislationIdentification().
func (w *QCPSBWrapper) LegislationIdentification() string {
	if w.wrapped.LegislationIdentification != nil {
		return *w.wrapped.LegislationIdentification
	}
	return ""
}

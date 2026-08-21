// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/TrustedEntityServiceWrapper.java (DSS 6.5.RC1).
//
// Extends TrustedSourceServiceWrapper in Java (own TrustedSourceList/ListOfTrustedSourceList
// fields on top of the inherited entityNames/tradeNames/serviceDigitalIdentifier/serviceNames/
// countryCode/status/type/startDate/endDate/capturedQualifiers/additionalServiceInfos state).
// The inherited fields are declared directly here rather than via an embedded
// TrustedSourceServiceWrapper; see the package note at the top of
// trusted_source_service_wrapper.go for why (composite-literal field-key constraints, and the
// exact field names DIAGWRAP_A's certificate_wrapper.go already relies on).
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// TrustedEntityServiceWrapper provides a wrapper for a Trusted Entity Service.
type TrustedEntityServiceWrapper struct {
	// TrustedSourceList is the corresponding Trusted Source List.
	TrustedSourceList *jaxb.XmlTrustSourceList
	// ListOfTrustedSourceList is the corresponding List of Trusted Source Lists.
	ListOfTrustedSourceList *jaxb.XmlTrustSourceList

	// EntityNames are the Trusted Entity names (inherited from TrustedSourceServiceWrapper in
	// Java).
	EntityNames []string
	// TradeNames are the Trusted Entity trade names (inherited).
	TradeNames []string
	// ServiceDigitalIdentifier is the related certificate (inherited).
	ServiceDigitalIdentifier *CertificateWrapper
	// ServiceNames are the trusted service names (inherited).
	ServiceNames []string
	// CountryCode is the country code (inherited).
	CountryCode string
	// Status is the status (inherited).
	Status string
	// Type is the Service Type Identifier URI (inherited).
	Type string
	// StartDate is the start date of validity (inherited).
	StartDate *time.Time
	// EndDate is the end date of validity (inherited).
	EndDate *time.Time
	// CapturedQualifiers are the captured qualifiers (inherited).
	CapturedQualifiers []*jaxb.XmlQualifier
	// AdditionalServiceInfos are the additional service informations (inherited).
	AdditionalServiceInfos []string
}

// NewTrustedEntityServiceWrapper is the default constructor. Port of the empty constructor.
func NewTrustedEntityServiceWrapper() *TrustedEntityServiceWrapper {
	return &TrustedEntityServiceWrapper{}
}

// CapturedQualifierUris gets captured qualifiers. Port of the inherited
// getCapturedQualifierUris() (duplicated here since Go has no method inheritance; see
// TrustedSourceServiceWrapper.CapturedQualifierUris(), which computes it identically).
func (w *TrustedEntityServiceWrapper) CapturedQualifierUris() []string {
	if w.CapturedQualifiers == nil {
		return nil
	}
	result := make([]string, 0, len(w.CapturedQualifiers))
	for _, q := range w.CapturedQualifiers {
		result = append(result, q.Value)
	}
	return result
}

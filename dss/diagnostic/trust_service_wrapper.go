// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/TrustServiceWrapper.java (DSS 6.5.RC1).
//
// Extends TrustedSourceServiceWrapper in Java. The inherited fields are declared directly here
// rather than via an embedded TrustedSourceServiceWrapper; see the package note at the top of
// trusted_source_service_wrapper.go for why. Java's getTspNames()/setTspNames() and
// getTspTradeNames()/setTspTradeNames() are trivial aliases of the inherited
// getEntityNames()/setEntityNames() and getTradeNames()/setTradeNames() pairs (both read/write
// the single inherited "entityNames"/"tradeNames" field) - ported here as the TspNames/
// TspTradeNames fields (matching the names certificate_wrapper.go already relies on), with
// EntityNames()/TradeNames()/SetEntityNames()/SetTradeNames() as thin delegating
// methods standing in for the inherited accessor names.
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// TrustServiceWrapper wraps an extracted information from a Trusted Service.
type TrustServiceWrapper struct {
	// TrustedList is the corresponding Trusted List.
	TrustedList *jaxb.XmlTrustedList
	// ListOfTrustedLists is the corresponding List of Trusted Lists.
	ListOfTrustedLists *jaxb.XmlTrustedList
	// EnactedMRA defines whether MRA has been applied for this particular Trusted Service.
	EnactedMRA *bool
	// MraTrustServiceLegalIdentifier is the name of the Trust Service defined the Mutual
	// Recognition Agreement equivalence scheme.
	MraTrustServiceLegalIdentifier string
	// MraTrustServiceEquivalenceStatusStartingTime is the date when the status for the current
	// MRA Trust Service equivalence has been started.
	MraTrustServiceEquivalenceStatusStartingTime *time.Time
	// MraTrustServiceEquivalenceStatusEndingTime is the date when the status for the current
	// MRA Trust Service equivalence has been ended (if applicable).
	MraTrustServiceEquivalenceStatusEndingTime *time.Time
	// OriginalTCStatus is the original third-country status before applied MRA.
	OriginalTCStatus string
	// OriginalTCType is the original third-country type before applied MRA.
	OriginalTCType string
	// OriginalCapturedQualifiers are the original third-country captured qualifiers before
	// applied MRA.
	OriginalCapturedQualifiers []*jaxb.XmlQualifier
	// OriginalTCAdditionalServiceInfos are the original third-country captured qualifiers
	// before applied MRA (sic, the Java Javadoc is copy-pasted from
	// OriginalCapturedQualifiers; reproduced as-is).
	OriginalTCAdditionalServiceInfos []string

	// TspNames are the Trusted Service Provider names; storage shared with the inherited
	// EntityNames()/SetEntityNames() accessors, see the file note above.
	TspNames []string
	// TspTradeNames are the Trusted Service Provider trade names; storage shared with the
	// inherited TradeNames()/SetTradeNames() accessors, see the file note above.
	TspTradeNames []string

	// ServiceDigitalIdentifier is the related certificate (inherited from
	// TrustedSourceServiceWrapper).
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

// NewTrustServiceWrapper is the default constructor.
func NewTrustServiceWrapper() *TrustServiceWrapper {
	return &TrustServiceWrapper{}
}

// EntityNames gets Trusted Service Provider names. Port of the inherited getEntityNames()
// (TrustServiceWrapper.getTspNames() is the same accessor under a different name; see the file
// note above).
func (w *TrustServiceWrapper) EntityNames() []string { return w.TspNames }

// SetEntityNames sets Trusted Service Provider names. Port of the inherited setEntityNames().
func (w *TrustServiceWrapper) SetEntityNames(entityNames []string) { w.TspNames = entityNames }

// TradeNames gets Trusted Service Provider trade names. Port of the inherited getTradeNames().
func (w *TrustServiceWrapper) TradeNames() []string { return w.TspTradeNames }

// SetTradeNames sets Trusted Service Provider trade names. Port of the inherited
// setTradeNames().
func (w *TrustServiceWrapper) SetTradeNames(tradeNames []string) { w.TspTradeNames = tradeNames }

// IsEnactedMRA gets whether MRA has been enacted for this Trusted Service. Port of
// isEnactedMRA().
func (w *TrustServiceWrapper) IsEnactedMRA() bool {
	return w.EnactedMRA != nil && *w.EnactedMRA
}

// OriginalCapturedQualifierUris gets original third-country captured qualifier URIs defined
// within Trusted List (before applied MRA). Port of getOriginalCapturedQualifierUris().
func (w *TrustServiceWrapper) OriginalCapturedQualifierUris() []string {
	if w.OriginalCapturedQualifiers == nil {
		return nil
	}
	result := make([]string, 0, len(w.OriginalCapturedQualifiers))
	for _, q := range w.OriginalCapturedQualifiers {
		result = append(result, q.Value)
	}
	return result
}

// CapturedQualifierUris gets captured qualifiers. Port of the inherited
// getCapturedQualifierUris() (duplicated here since Go has no method inheritance; see
// TrustedSourceServiceWrapper.CapturedQualifierUris(), which computes it identically).
func (w *TrustServiceWrapper) CapturedQualifierUris() []string {
	if w.CapturedQualifiers == nil {
		return nil
	}
	result := make([]string, 0, len(w.CapturedQualifiers))
	for _, q := range w.CapturedQualifiers {
		result = append(result, q.Value)
	}
	return result
}

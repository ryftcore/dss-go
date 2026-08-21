// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/TrustedSourceServiceWrapper.java (DSS 6.5.RC1).
//
// Java's TrustedSourceServiceWrapper is an abstract JavaBean: private fields with matched
// getX()/setX() pairs, no computed logic beyond CapturedQualifierUris. Every one of its
// getX()/setX() pairs is trivial field access, so - per PORTING.md's "minus get/set prefixes
// where un-idiomatic" rule - each becomes a single exported Go field rather than a method pair;
// an exported field already serves as both getter (read) and setter (assign) with no loss of
// capability, and this port's two concrete subclasses (TrustServiceWrapper,
// TrustedEntityServiceWrapper, both DIAGWRAP_B) are already relied upon elsewhere in the port
// (dss/diagnostic/certificate_wrapper.go, DIAGWRAP_A) as plain struct literals with these exact
// exported field names.
//
// Because Go composite literals can only set a field through the struct that directly declares
// it - never through an anonymous embedded field's promoted name - and CertificateWrapper's
// TrustServices()/TrustedEntityServices() build TrustServiceWrapper/TrustedEntityServiceWrapper
// with flat field keys (ServiceDigitalIdentifier: ..., Status: ..., ...), this base type is NOT
// embedded by either subclass: each subclass instead declares the inherited fields directly
// (matching Java's inherited state) and repeats the one non-trivial accessor,
// CapturedQualifierUris(). TrustedSourceServiceWrapper itself remains a complete, correct,
// independently usable port of the Java class (Java never instantiates the abstract class
// directly either - only through TrustServiceWrapper/TrustedEntityServiceWrapper), for callers
// that want to build a bare TrustedSourceServiceWrapper value directly.
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// TrustedSourceServiceWrapper mirrors the Java abstract class of the same name.
type TrustedSourceServiceWrapper struct {
	// EntityNames are the Trusted Entity names.
	EntityNames []string
	// TradeNames are the Trusted Entity trade names.
	TradeNames []string
	// ServiceDigitalIdentifier is the related certificate.
	ServiceDigitalIdentifier *CertificateWrapper
	// ServiceNames are the trusted service names.
	ServiceNames []string
	// CountryCode is the country code.
	CountryCode string
	// Status is the status.
	Status string
	// Type is the Service Type Identifier URI.
	Type string
	// StartDate is the start date of validity.
	StartDate *time.Time
	// EndDate is the end date of validity.
	EndDate *time.Time
	// CapturedQualifiers are the captured qualifiers.
	CapturedQualifiers []*jaxb.XmlQualifier
	// AdditionalServiceInfos are the additional service informations.
	AdditionalServiceInfos []string
}

// NewTrustedSourceServiceWrapper is the default constructor. Port of the protected empty
// constructor.
func NewTrustedSourceServiceWrapper() *TrustedSourceServiceWrapper {
	return &TrustedSourceServiceWrapper{}
}

// CapturedQualifierUris gets captured qualifiers. Port of getCapturedQualifierUris().
func (w *TrustedSourceServiceWrapper) CapturedQualifierUris() []string {
	if w.CapturedQualifiers == nil {
		return nil
	}
	result := make([]string, 0, len(w.CapturedQualifiers))
	for _, q := range w.CapturedQualifiers {
		result = append(result, q.Value)
	}
	return result
}

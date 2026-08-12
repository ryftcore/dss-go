// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/InhibitAnyPolicy.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// InhibitAnyPolicy is RFC 5280 4.2.1.14. Inhibit anyPolicy.
//
// The inhibit anyPolicy extension can be used in certificates issued to CAs. The inhibit
// anyPolicy extension indicates that the special anyPolicy OID, with the value
// { 2 5 29 32 0 }, is not considered an explicit match for other certificate policies
// except when it appears in an intermediate self-issued CA certificate.
type InhibitAnyPolicy struct {
	CertificateExtension

	// value indicates the number of additional non-self-issued certificates that may
	// appear in the path before anyPolicy is no longer permitted.
	value int
}

// NewInhibitAnyPolicy builds an InhibitAnyPolicy extension. Ports the Java field
// initializer of value = -1.
func NewInhibitAnyPolicy() *InhibitAnyPolicy {
	return &InhibitAnyPolicy{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_INHIBIT_ANY_POLICY),
		value:                -1,
	}
}

// Value gets the InhibitAnyPolicy constraint value: requireExplicitPolicy int value if
// present, -1 otherwise.
func (i *InhibitAnyPolicy) Value() int {
	return i.value
}

// SetValue sets the InhibitAnyPolicy constraint value.
func (i *InhibitAnyPolicy) SetValue(value int) {
	i.value = value
}

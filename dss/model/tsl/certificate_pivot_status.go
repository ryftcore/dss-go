// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/CertificatePivotStatus.java (DSS 6.5.RC1).
package tsl

// CertificatePivotStatus describes a certificate status in the current pivot.
type CertificatePivotStatus string

const (
	// CertificatePivotStatusAdded marks a certificate that has been added with a new pivot.
	CertificatePivotStatusAdded CertificatePivotStatus = "ADDED"
	// CertificatePivotStatusNotChanged marks a certificate that has not been changed.
	CertificatePivotStatusNotChanged CertificatePivotStatus = "NOT_CHANGED"
	// CertificatePivotStatusRemoved marks a certificate that has been removed with a new pivot.
	CertificatePivotStatusRemoved CertificatePivotStatus = "REMOVED"
)

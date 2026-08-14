// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/CertificatePivotStatus.java (DSS 6.5.RC1).
package tsl

// CertificatePivotStatus describes a certificate status in the current pivot.
type CertificatePivotStatus string

const (
	// CertificatePivotStatus_ADDED marks a certificate that has been added with a new pivot.
	CertificatePivotStatus_ADDED CertificatePivotStatus = "ADDED"
	// CertificatePivotStatus_NOT_CHANGED marks a certificate that has not been changed.
	CertificatePivotStatus_NOT_CHANGED CertificatePivotStatus = "NOT_CHANGED"
	// CertificatePivotStatus_REMOVED marks a certificate that has been removed with a new pivot.
	CertificatePivotStatus_REMOVED CertificatePivotStatus = "REMOVED"
)

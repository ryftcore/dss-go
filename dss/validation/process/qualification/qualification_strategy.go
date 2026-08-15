// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationStrategy.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/enumerations"

// QualificationStrategy extracts the qualification status for a certificate.
type QualificationStrategy interface {
	// QualifiedStatus gets certificate qualification status. Port of
	// getQualifiedStatus().
	QualifiedStatus() enumerations.CertificateQualifiedStatus
}

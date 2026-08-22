// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationStrategy.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/enumerations"

// Strategy extracts the qualification status for a certificate.
type Strategy interface {
	// QualifiedStatus gets certificate qualification status. Port of
	// getQualifiedStatus().
	QualifiedStatus() enumerations.CertificateQualifiedStatus
}

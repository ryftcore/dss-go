// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeStrategy.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/enumerations"

// TypeStrategy extracts certificate approval status type for a certificate.
type TypeStrategy interface {
	// Type gets certificate approval status type. Port of getType().
	Type() enumerations.CertificateType
}

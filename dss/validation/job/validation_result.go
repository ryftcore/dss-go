// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/validation/ValidationResult.java (DSS 6.5.RC1).
package job

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// ValidationResult provides an interface for extraction of information about a validation
// task result.
type ValidationResult interface {
	CachedResult

	// Indication returns the validation Indication. Port of getIndication().
	Indication() enumerations.Indication
	// SubIndication returns the validation SubIndication. Port of getSubIndication().
	SubIndication() enumerations.SubIndication
	// SigningTime returns the (claimed) signing time. The zero time.Time stands for Java's
	// null. Port of getSigningTime().
	SigningTime() time.Time
	// SigningCertificate returns the signing certificate. Port of getSigningCertificate().
	SigningCertificate() *model.CertificateToken
	// PotentialSigners returns a list of signing candidates. Port of getPotentialSigners().
	PotentialSigners() []*model.CertificateToken
}

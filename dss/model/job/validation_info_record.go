// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ValidationInfoRecord.java (DSS 6.5.RC1).
package job

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// ValidationInfoRecord is a validation result record.
type ValidationInfoRecord interface {
	InfoRecord

	// Indication gets the validation Indication. Port of getIndication().
	Indication() enumerations.Indication
	// SubIndication gets the validation SubIndication. Port of getSubIndication().
	SubIndication() enumerations.SubIndication
	// SigningTime gets the (claimed) signing time. The zero time.Time stands for Java's null.
	// Port of getSigningTime().
	SigningTime() time.Time
	// SigningCertificate gets the signing certificate. Port of getSigningCertificate().
	SigningCertificate() *model.CertificateToken
	// PotentialSigners gets a list of signing candidates. Port of getPotentialSigners().
	PotentialSigners() []*model.CertificateToken
	// IsValid gets if the signature is valid. Port of isValid().
	IsValid() bool
	// IsIndeterminate gets if the validation result is indeterminate. Port of
	// isIndeterminate().
	IsIndeterminate() bool
	// IsInvalid gets if the signature is invalid. Port of isInvalid().
	IsInvalid() bool
}

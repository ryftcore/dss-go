// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/dto/ValidationCacheDTO.java (DSS 6.5.RC1).
package job

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modeljob "github.com/ryftcore/dss-go/dss/model/job"
)

// ValidationCacheDTO is the validation record DTO. It implements
// modeljob.ValidationInfoRecord.
type ValidationCacheDTO struct {
	*AbstractCacheDTO

	// indication is the validation Indication.
	indication enumerations.Indication

	// subIndication is the validation SubIndication.
	subIndication enumerations.SubIndication

	// signingTime is the claimed signing time. The zero time.Time stands for Java's null.
	signingTime time.Time

	// signingCertificate is the signing certificate.
	signingCertificate *model.CertificateToken

	// potentialSigners holds signing candidates.
	potentialSigners []*model.CertificateToken
}

// NewValidationCacheDTO creates an empty ValidationCacheDTO. Port of the default
// constructor.
func NewValidationCacheDTO() *ValidationCacheDTO {
	return &ValidationCacheDTO{AbstractCacheDTO: NewAbstractCacheDTO()}
}

// NewValidationCacheDTOFrom copies cacheDTO. Port of the copy constructor.
func NewValidationCacheDTOFrom(cacheDTO *AbstractCacheDTO) *ValidationCacheDTO {
	return &ValidationCacheDTO{AbstractCacheDTO: NewAbstractCacheDTOFrom(cacheDTO)}
}

// Indication returns the validation Indication. Port of getIndication().
func (d *ValidationCacheDTO) Indication() enumerations.Indication {
	return d.indication
}

// SetIndication sets the validation Indication. Port of setIndication(Indication).
func (d *ValidationCacheDTO) SetIndication(indication enumerations.Indication) {
	d.indication = indication
}

// SubIndication returns the validation SubIndication. Port of getSubIndication().
func (d *ValidationCacheDTO) SubIndication() enumerations.SubIndication {
	return d.subIndication
}

// SetSubIndication sets the validation SubIndication. Port of
// setSubIndication(SubIndication).
func (d *ValidationCacheDTO) SetSubIndication(subIndication enumerations.SubIndication) {
	d.subIndication = subIndication
}

// SigningTime returns the claimed signing time. Port of getSigningTime().
func (d *ValidationCacheDTO) SigningTime() time.Time {
	return d.signingTime
}

// SetSigningTime sets the claimed signing time. Port of setSigningTime(Date).
func (d *ValidationCacheDTO) SetSigningTime(signingTime time.Time) {
	d.signingTime = signingTime
}

// SigningCertificate returns the signing certificate token. Port of
// getSigningCertificate().
func (d *ValidationCacheDTO) SigningCertificate() *model.CertificateToken {
	return d.signingCertificate
}

// SetSigningCertificate sets the signing certificate token. Port of
// setSigningCertificate(CertificateToken).
func (d *ValidationCacheDTO) SetSigningCertificate(signingCertificate *model.CertificateToken) {
	d.signingCertificate = signingCertificate
}

// PotentialSigners returns the signing candidates. Port of getPotentialSigners().
func (d *ValidationCacheDTO) PotentialSigners() []*model.CertificateToken {
	return d.potentialSigners
}

// SetPotentialSigners sets the signing candidates. Port of setPotentialSigners(List).
func (d *ValidationCacheDTO) SetPotentialSigners(potentialSigners []*model.CertificateToken) {
	d.potentialSigners = potentialSigners
}

// IsValid reports whether the indication is TOTAL_PASSED. Port of isValid().
func (d *ValidationCacheDTO) IsValid() bool {
	return enumerations.Indication_TOTAL_PASSED == d.indication
}

// IsIndeterminate reports whether the indication is INDETERMINATE. Port of
// isIndeterminate().
func (d *ValidationCacheDTO) IsIndeterminate() bool {
	return enumerations.Indication_INDETERMINATE == d.indication
}

// IsInvalid reports whether the indication is TOTAL_FAILED. Port of isInvalid().
func (d *ValidationCacheDTO) IsInvalid() bool {
	return enumerations.Indication_TOTAL_FAILED == d.indication
}

var _ modeljob.ValidationInfoRecord = (*ValidationCacheDTO)(nil)

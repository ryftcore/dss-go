// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/CertificateValidationBeforeSunsetDateWithIdCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// CertificateValidationBeforeSunsetDateWithIdCheck verifies whether a
// validation time is before certificate's trust sunset date, with a
// certificate id provided.
type CertificateValidationBeforeSunsetDateWithIdCheck[T any] struct {
	*CertificateValidationBeforeSunsetDateCheck[T]
}

// NewCertificateValidationBeforeSunsetDateWithIdCheck is the default
// constructor. Port of
// CertificateValidationBeforeSunsetDateWithIdCheck(I18nProvider, T, CertificateWrapper, Date, LevelRule).
func NewCertificateValidationBeforeSunsetDateWithIdCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, controlTime time.Time,
	constraint policy.LevelRule) *CertificateValidationBeforeSunsetDateWithIdCheck[T] {
	id := certificate.Id()
	return &CertificateValidationBeforeSunsetDateWithIdCheck[T]{
		CertificateValidationBeforeSunsetDateCheck: newCertificateValidationBeforeSunsetDateCheck(
			i18nProvider, result, certificate, controlTime, constraint, &id),
	}
}

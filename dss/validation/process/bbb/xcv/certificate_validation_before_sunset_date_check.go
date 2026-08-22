// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/CertificateValidationBeforeSunsetDateCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateValidationBeforeSunsetDateCheck verifies whether a validation
// time is before certificate's trust sunset date.
type CertificateValidationBeforeSunsetDateCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// controlTime is the validation time to check against.
	controlTime time.Time
}

// NewCertificateValidationBeforeSunsetDateCheck is the default constructor.
// Port of CertificateValidationBeforeSunsetDateCheck(I18nProvider, T, CertificateWrapper, Date, LevelRule).
func NewCertificateValidationBeforeSunsetDateCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, controlTime time.Time, constraint policy.LevelRule) *CertificateValidationBeforeSunsetDateCheck[T] {
	return newCertificateValidationBeforeSunsetDateCheck(i18nProvider, result, certificate, controlTime, constraint, nil)
}

// newCertificateValidationBeforeSunsetDateCheck is the constructor with an
// optional certificate identifier. Port of the protected constructor
// CertificateValidationBeforeSunsetDateCheck(Provider, T, CertificateWrapper, Date, LevelRule, String);
// a nil certificateId matches Java's null (no explicit BBB id).
func newCertificateValidationBeforeSunsetDateCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, controlTime time.Time, constraint policy.LevelRule,
	certificateId *string) *CertificateValidationBeforeSunsetDateCheck[T] {
	var chainItemBase *process.ChainItemBase[T]
	if certificateId != nil {
		chainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *certificateId)
	} else {
		chainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &CertificateValidationBeforeSunsetDateCheck[T]{
		ChainItemBase: chainItemBase,
		certificate:   certificate,
		controlTime:   controlTime,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeSubXCVTA
}

// Process performs the check. Port of process().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) Process() bool {
	if c.certificate.TrustSunsetDate() != nil {
		return c.controlTime.Before(*c.certificate.TrustSunsetDate())
	}
	// if no Sunset date, trust indefinitely
	return c.certificate.IsTrusted()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIVTBCTSD
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIVTBCTSDANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoCertificateChainFoundNoPOE
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CertificateValidationBeforeSunsetDateCheck[T]) BuildAdditionalInfo() *string {
	var message string
	if c.certificate.TrustSunsetDate() != nil {
		controlTime := c.controlTime
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateSunsetDate,
			process.GetFormattedDate(&controlTime), process.GetFormattedDate(c.certificate.TrustSunsetDate()))
	} else {
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateSunsetDateValid)
	}
	return &message
}

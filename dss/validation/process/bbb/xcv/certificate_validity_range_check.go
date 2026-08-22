// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateValidityRangeCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateValidityRangeCheck checks if the certificate is not expired.
type CertificateValidityRangeCheck[T any] struct {
	*process.ChainItemBase[T]

	// currentTime is the validation date.
	currentTime time.Time

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// usedCertificateRevocation is the certificate's revocation.
	usedCertificateRevocation *diagnostic.CertificateRevocationWrapper

	// revocationDataRequired defines whether revocation data is required for
	// the certificate.
	revocationDataRequired bool

	// revocationIssuerTrusted defines whether the revocation data's issuer is
	// trusted.
	revocationIssuerTrusted bool

	// revocationIssuerCheckEnforced defines whether the validation is
	// enforced for the revocation data issuer.
	revocationIssuerCheckEnforced bool
}

// NewCertificateValidityRangeCheckMinimal is the minimal (protected)
// constructor. Port of
// CertificateValidityRangeCheck(I18nProvider, T, CertificateWrapper, Date, LevelRule).
func NewCertificateValidityRangeCheckMinimal[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, constraint policy.LevelRule) *CertificateValidityRangeCheck[T] {
	return NewCertificateValidityRangeCheck(i18nProvider, result, certificate, nil, false, false, false, currentTime, constraint)
}

// NewCertificateValidityRangeCheck is the default constructor. Port of
// CertificateValidityRangeCheck(I18nProvider, T, CertificateWrapper, CertificateRevocationWrapper, boolean, boolean, boolean, Date, LevelRule).
func NewCertificateValidityRangeCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, usedCertificateRevocation *diagnostic.CertificateRevocationWrapper,
	revocationDataRequired, revocationIssuerTrusted, revocationIssuerCheckEnforced bool,
	currentTime time.Time, constraint policy.LevelRule) *CertificateValidityRangeCheck[T] {
	c := &CertificateValidityRangeCheck[T]{
		ChainItemBase:                 process.NewChainItemBase(i18nProvider, result, constraint),
		currentTime:                   currentTime,
		certificate:                   certificate,
		usedCertificateRevocation:     usedCertificateRevocation,
		revocationDataRequired:        revocationDataRequired,
		revocationIssuerTrusted:       revocationIssuerTrusted,
		revocationIssuerCheckEnforced: revocationIssuerCheckEnforced,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateValidityRangeCheck[T]) Process() bool {
	return c.isInValidityRange(c.certificate)
}

// isInValidityRange ports the private isInValidityRange(CertificateWrapper).
func (c *CertificateValidityRangeCheck[T]) isInValidityRange(certificateWrapper *diagnostic.CertificateWrapper) bool {
	if certificateWrapper != nil {
		notBefore := certificateWrapper.NotBefore()
		notAfter := certificateWrapper.NotAfter()
		return notBefore != nil && !c.currentTime.Before(*notBefore) &&
			notAfter != nil && !c.currentTime.After(*notAfter)
	}
	return false
}

// isRevocationDataValid ports the private isRevocationDataValid().
func (c *CertificateValidityRangeCheck[T]) isRevocationDataValid() bool {
	// other checks are performed before
	return c.revocationIssuerTrusted || !c.revocationIssuerCheckEnforced ||
		c.isInValidityRange(c.usedCertificateRevocation.SigningCertificate())
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CertificateValidityRangeCheck[T]) BuildAdditionalInfo() *string {
	notBeforeStr := " ? "
	if c.certificate.NotBefore() != nil {
		notBeforeStr = process.GetFormattedDate(c.certificate.NotBefore())
	}
	notAfterStr := " ? "
	if c.certificate.NotAfter() != nil {
		notAfterStr = process.GetFormattedDate(c.certificate.NotAfter())
	}
	currentTime := c.currentTime
	validationTime := process.GetFormattedDate(&currentTime)
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_VALIDITY, validationTime, notBeforeStr, notAfterStr)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateValidityRangeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICTIVRSC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateValidityRangeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICTIVRSC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateValidityRangeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateValidityRangeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	certificateIsKnownToNotBeRevoked := c.usedCertificateRevocation != nil &&
		!c.usedCertificateRevocation.IsRevoked() && c.isRevocationDataValid()
	if !c.revocationDataRequired || certificateIsKnownToNotBeRevoked {
		return enumerations.SubIndicationOutOfBoundsNotRevoked
	}
	return enumerations.SubIndicationOutOfBoundsNoPOE
}

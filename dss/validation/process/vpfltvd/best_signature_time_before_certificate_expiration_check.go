// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/BestSignatureTimeBeforeCertificateExpirationCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BestSignatureTimeBeforeCertificateExpirationCheck checks if the
// best-signature-time is before certificate's expiration.
type BestSignatureTimeBeforeCertificateExpirationCheck[T any] struct {
	*process.ChainItemBase[T]

	// bestSignatureTime is the best signature time.
	bestSignatureTime *time.Time

	// signingCertificate is the signing certificate.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewBestSignatureTimeBeforeCertificateExpirationCheck is the default
// constructor. Port of
// BestSignatureTimeBeforeCertificateExpirationCheck(I18nProvider, T, Date, CertificateWrapper, LevelRule).
func NewBestSignatureTimeBeforeCertificateExpirationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	bestSignatureTime *time.Time, signingCertificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *BestSignatureTimeBeforeCertificateExpirationCheck[T] {
	c := &BestSignatureTimeBeforeCertificateExpirationCheck[T]{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		bestSignatureTime:  bestSignatureTime,
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process(): inclusive by RFC 5280.
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) Process() bool {
	notAfter := c.signingCertificate.NotAfter()
	if c.bestSignatureTime == nil || notAfter == nil {
		return false
	}
	return !c.bestSignatureTime.After(*notAfter)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) BuildAdditionalInfo() *string {
	bestSignatureTimeStr := " ? "
	if c.bestSignatureTime != nil {
		bestSignatureTimeStr = process.GetFormattedDate(c.bestSignatureTime)
	}
	certNotAfter := " ? "
	if c.signingCertificate.NotAfter() != nil {
		certNotAfter = process.GetFormattedDate(c.signingCertificate.NotAfter())
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_AFTER, bestSignatureTimeStr, certNotAfter)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TSV_IBSTBCEC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TSV_IBSTBCEC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BestSignatureTimeBeforeCertificateExpirationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_OUT_OF_BOUNDS_NOT_REVOKED
}

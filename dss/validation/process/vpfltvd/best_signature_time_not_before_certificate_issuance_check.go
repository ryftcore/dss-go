// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/BestSignatureTimeNotBeforeCertificateIssuanceCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BestSignatureTimeNotBeforeCertificateIssuanceCheck: if best-signature-time
// is before the issuance date of the signing certificate, the process shall
// return the indication FAILED with the sub-indication NOT_YET_VALID.
// Otherwise, the process shall return the indication and sub-indication which
// was returned by previous step.
type BestSignatureTimeNotBeforeCertificateIssuanceCheck[T any] struct {
	*process.ChainItemBase[T]

	// bestSignatureTime is the best signature time.
	bestSignatureTime *time.Time

	// signingCertificate is the signing certificate.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewBestSignatureTimeNotBeforeCertificateIssuanceCheck is the default
// constructor. Port of
// BestSignatureTimeNotBeforeCertificateIssuanceCheck(I18nProvider, T, Date, CertificateWrapper, LevelRule).
func NewBestSignatureTimeNotBeforeCertificateIssuanceCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	bestSignatureTime *time.Time, signingCertificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T] {
	c := &BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		bestSignatureTime:  bestSignatureTime,
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) Process() bool {
	notBefore := c.signingCertificate.NotBefore()
	if c.bestSignatureTime == nil || notBefore == nil {
		return false
	}
	return !c.bestSignatureTime.Before(*notBefore)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) BuildAdditionalInfo() *string {
	bestSignatureTimeStr := " ? "
	if c.bestSignatureTime != nil {
		bestSignatureTimeStr = process.GetFormattedDate(c.bestSignatureTime)
	}
	certNotBefore := " ? "
	if c.signingCertificate.NotBefore() != nil {
		certNotBefore = process.GetFormattedDate(c.signingCertificate.NotBefore())
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_BEST_SIGNATURE_TIME_CERT_NOT_BEFORE, bestSignatureTimeStr, certNotBefore)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TSV_IBSTAIDOSC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TSV_IBSTAIDOSC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BestSignatureTimeNotBeforeCertificateIssuanceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NOT_YET_VALID
}

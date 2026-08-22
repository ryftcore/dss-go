// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/TrustedListReachedForCertificateChainCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustedListReachedForCertificateChainCheck checks whether a Trusted List
// has been reached for the given certificate chain.
type TrustedListReachedForCertificateChainCheck[T any] struct {
	*process.ChainItemBase[T]

	// signingCertificate is the end-entity certificate.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewTrustedListReachedForCertificateChainCheck is the default constructor.
// Port of
// TrustedListReachedForCertificateChainCheck(I18nProvider, T, CertificateWrapper, LevelRule).
func NewTrustedListReachedForCertificateChainCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	signingCertificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *TrustedListReachedForCertificateChainCheck[T] {
	c := &TrustedListReachedForCertificateChainCheck[T]{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedListReachedForCertificateChainCheck[T]) Process() bool {
	return c.signingCertificate != nil && c.signingCertificate.IsTrustedListReached()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedListReachedForCertificateChainCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedListReachedForCertificateChainCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_CERT_TRUSTED_LIST_REACHED_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedListReachedForCertificateChainCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedListReachedForCertificateChainCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

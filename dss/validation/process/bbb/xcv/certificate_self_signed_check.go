// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateSelfSignedCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// CertificateSelfSignedCheck checks if the certificate is self-signed.
type CertificateSelfSignedCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateSelfSignedCheck is the default constructor. Port of
// CertificateSelfSignedCheck(I18nProvider, T, CertificateWrapper, LevelRule).
func NewCertificateSelfSignedCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateSelfSignedCheck[T] {
	c := &CertificateSelfSignedCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateSelfSignedCheck[T]) Process() bool {
	return c.certificate.IsSelfSigned()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateSelfSignedCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISSSC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateSelfSignedCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISSSC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateSelfSignedCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateSelfSignedCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}

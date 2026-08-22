// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateSignatureValidCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateSignatureValidCheck checks if the certificate's signature is
// valid.
type CertificateSignatureValidCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateSignatureValidCheck is the default constructor. Port of
// CertificateSignatureValidCheck(I18nProvider, T, CertificateWrapper, LevelRule).
func NewCertificateSignatureValidCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateSignatureValidCheck[T] {
	c := &CertificateSignatureValidCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateSignatureValidCheck[T]) Process() bool {
	return c.certificate.IsSignatureValid()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateSignatureValidCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICSI
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateSignatureValidCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICSI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateSignatureValidCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateSignatureValidCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/ProspectiveCertificateChainCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ProspectiveCertificateChainCheck checks if the certificate chain is
// trusted.
type ProspectiveCertificateChainCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// context is the validation context.
	context enumerations.Context
}

// NewProspectiveCertificateChainCheck is the default constructor. Port of
// ProspectiveCertificateChainCheck(I18nProvider, T, CertificateWrapper, Context, LevelRule).
func NewProspectiveCertificateChainCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, context enumerations.Context,
	constraint policy.LevelRule) *ProspectiveCertificateChainCheck[T] {
	c := &ProspectiveCertificateChainCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		context:       context,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ProspectiveCertificateChainCheck[T]) Process() bool {
	return c.certificate != nil && (c.certificate.IsTrusted() || c.certificate.IsTrustedChain())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ProspectiveCertificateChainCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CCCBB
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ProspectiveCertificateChainCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch c.context {
	case enumerations.ContextSignature, enumerations.ContextCounterSignature, enumerations.ContextKeyBindingSignature:
		return i18n.MessageTag_BBB_XCV_CCCBB_SIG_ANS
	case enumerations.ContextTimestamp:
		return i18n.MessageTag_BBB_XCV_CCCBB_TSP_ANS
	case enumerations.ContextRevocation:
		return i18n.MessageTag_BBB_XCV_CCCBB_REV_ANS
	default:
		return i18n.MessageTag_BBB_XCV_CCCBB_ANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ProspectiveCertificateChainCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ProspectiveCertificateChainCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoCertificateChainFound
}

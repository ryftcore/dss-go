// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/POEExistsWithinCertificateValidityRangeCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// POEExistsWithinCertificateValidityRangeCheck verifies if the set of POEs
// contains a POE for the certificate after the issuance date and before the
// expiration date of that certificate.
type POEExistsWithinCertificateValidityRangeCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check POE.
	certificate *diagnostic.CertificateWrapper

	// poe is a collection of POEs.
	poe *POEExtraction
}

// NewPOEExistsWithinCertificateValidityRangeCheck is the default constructor.
// Port of POEExistsWithinCertificateValidityRangeCheck(I18nProvider, T, CertificateWrapper, POEExtraction, LevelRule).
func NewPOEExistsWithinCertificateValidityRangeCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, poe *POEExtraction,
	constraint policy.LevelRule) *POEExistsWithinCertificateValidityRangeCheck[T] {
	c := &POEExistsWithinCertificateValidityRangeCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		poe:           poe,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) Process() bool {
	return c.certificate != nil &&
		c.poe.IsPOEExistInRange(c.certificate.Id(), c.certificate.NotBefore(), c.certificate.NotAfter())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVIPCRIAIDBEDC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVIPCRIAIDBEDCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationRevocationOutOfBoundsNoPOE
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
func (c *POEExistsWithinCertificateValidityRangeCheck[T]) BuildAdditionalInfo() *string {
	if c.certificate != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagCertificateID, c.certificate.Id())
		return &message
	}
	return nil
}

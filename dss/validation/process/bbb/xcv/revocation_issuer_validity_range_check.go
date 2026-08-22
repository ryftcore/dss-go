// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/RevocationIssuerValidityRangeCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationIssuerValidityRangeCheck verifies if a validation time is in the
// validity range of the certificate of the issuer of the revocation
// information.
type RevocationIssuerValidityRangeCheck[T any] struct {
	*process.ChainItemBase[T]

	// currentTime is the validation date.
	currentTime time.Time

	// revocationWrapper is the revocation data to verify issuer of.
	revocationWrapper *diagnostic.RevocationWrapper
}

// NewRevocationIssuerValidityRangeCheck is the default constructor. Port of
// RevocationIssuerValidityRangeCheck(I18nProvider, T, RevocationWrapper, Date, LevelRule).
func NewRevocationIssuerValidityRangeCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	revocationWrapper *diagnostic.RevocationWrapper, currentTime time.Time, constraint policy.LevelRule) *RevocationIssuerValidityRangeCheck[T] {
	c := &RevocationIssuerValidityRangeCheck[T]{
		ChainItemBase:     process.NewChainItemBase(i18nProvider, result, constraint),
		currentTime:       currentTime,
		revocationWrapper: revocationWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationIssuerValidityRangeCheck[T]) Process() bool {
	signingCertificate := c.revocationWrapper.SigningCertificate()
	if signingCertificate != nil {
		notBefore := signingCertificate.NotBefore()
		notAfter := signingCertificate.NotAfter()
		return notBefore != nil && !c.currentTime.Before(*notBefore) &&
			notAfter != nil && !c.currentTime.After(*notAfter)
	}
	return false
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationIssuerValidityRangeCheck[T]) BuildAdditionalInfo() *string {
	certificate := c.revocationWrapper.SigningCertificate()
	if certificate != nil {
		notBeforeStr := " ? "
		if certificate.NotBefore() != nil {
			notBeforeStr = process.GetFormattedDate(certificate.NotBefore())
		}
		notAfterStr := " ? "
		if certificate.NotAfter() != nil {
			notAfterStr = process.GetFormattedDate(certificate.NotAfter())
		}
		currentTime := c.currentTime
		validationTime := process.GetFormattedDate(&currentTime)
		message := c.I18nProvider.GetMessage(i18n.MessageTagRevocationCertValidity,
			certificate.Id(), c.revocationWrapper.Id(), notBeforeStr, notAfterStr, validationTime)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationIssuerValidityRangeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVICTIVRCIRI
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationIssuerValidityRangeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVICTIVRCIRIANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationIssuerValidityRangeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationIssuerValidityRangeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationRevocationOutOfBoundsNoPOE
}

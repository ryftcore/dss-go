// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SigningTimeInCertificateValidityRangeCheck.java (DSS 6.5.RC1).
package sav

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningTimeInCertificateValidityRangeCheck checks if a claimed signing time is
// within the signing-certificate's validity range.
type SigningTimeInCertificateValidityRangeCheck[T any] struct {
	*process.ChainItemBase[T]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSigningTimeInCertificateValidityRangeCheck is the default constructor. Port
// of SigningTimeInCertificateValidityRangeCheck(Provider, T, SignatureWrapper, LevelRule).
func NewSigningTimeInCertificateValidityRangeCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SigningTimeInCertificateValidityRangeCheck[T] {
	c := &SigningTimeInCertificateValidityRangeCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SigningTimeInCertificateValidityRangeCheck[T]) Process() bool {
	claimedSigningTime := c.signature.ClaimedSigningTime()
	signingCertificate := c.signature.SigningCertificate()
	return claimedSigningTime != nil && signingCertificate != nil &&
		!claimedSigningTime.Before(*signingCertificate.NotBefore()) &&
		!claimedSigningTime.After(*signingCertificate.NotAfter())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningTimeInCertificateValidityRangeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPSTWSCVR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningTimeInCertificateValidityRangeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVISQPSTWSCVRANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningTimeInCertificateValidityRangeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SigningTimeInCertificateValidityRangeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}

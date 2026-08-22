// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/AcceptableRevocationDataAvailableCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableRevocationDataAvailableCheck checks if an acceptable revocation
// data is found.
type AcceptableRevocationDataAvailableCheck[T any] struct {
	*process.ChainItemBase[T]

	// acceptableRevocationData is the revocation data to check.
	acceptableRevocationData *diagnostic.RevocationWrapper
}

// NewAcceptableRevocationDataAvailableCheck is the default constructor. Port
// of AcceptableRevocationDataAvailableCheck(Provider, T, RevocationWrapper, LevelRule).
func NewAcceptableRevocationDataAvailableCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	acceptableRevocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *AcceptableRevocationDataAvailableCheck[T] {
	c := &AcceptableRevocationDataAvailableCheck[T]{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		acceptableRevocationData: acceptableRevocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableRevocationDataAvailableCheck[T]) Process() bool {
	return c.acceptableRevocationData != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableRevocationDataAvailableCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIARDPFC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AcceptableRevocationDataAvailableCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIARDPFCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableRevocationDataAvailableCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AcceptableRevocationDataAvailableCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *AcceptableRevocationDataAvailableCheck[T]) BuildAdditionalInfo() *string {
	if c.acceptableRevocationData != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagLastAcceptableRevocation, c.acceptableRevocationData.Id())
		return &message
	}
	return nil
}

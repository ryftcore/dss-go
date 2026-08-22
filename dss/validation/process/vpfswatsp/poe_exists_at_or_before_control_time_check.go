// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/checks/POEExistsAtOrBeforeControlTimeCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// POEExistsAtOrBeforeControlTimeCheck checks if a POE exists before the control
// time.
type POEExistsAtOrBeforeControlTimeCheck[T any] struct {
	*process.ChainItemBase[T]

	// token is the token to check.
	token diagnostic.TokenProxy

	// referenceCategory is the object's type.
	referenceCategory enumerations.TimestampedObjectType

	// controlTime is the control time to check against.
	controlTime time.Time

	// poe is the POE container.
	poe *POEExtraction
}

// NewPOEExistsAtOrBeforeControlTimeCheck is the default constructor. Port of
// POEExistsAtOrBeforeControlTimeCheck(Provider, T, TokenProxy, TimestampedObjectType, Date, POEExtraction, LevelRule).
func NewPOEExistsAtOrBeforeControlTimeCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	token diagnostic.TokenProxy, referenceCategory enumerations.TimestampedObjectType, controlTime time.Time,
	poe *POEExtraction, constraint policy.LevelRule) *POEExistsAtOrBeforeControlTimeCheck[T] {
	c := &POEExistsAtOrBeforeControlTimeCheck[T]{
		ChainItemBase:     process.NewChainItemBase(i18nProvider, result, constraint),
		token:             token,
		referenceCategory: referenceCategory,
		controlTime:       controlTime,
		poe:               poe,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) Process() bool {
	return c.poe.IsPOEExists(c.token.Id(), c.controlTime)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagControlTime, c.token.Id(),
		process.GetFormattedDate(&c.controlTime))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag(), whose
// IllegalStateException becomes a panic - the method is called from the base
// ChainItem, which cannot propagate an error.
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) MessageTag() i18n.MessageTag {
	if enumerations.TimestampedObjectTypeCertificate == c.referenceCategory {
		return i18n.MessageTagPSVITPOCOBCT
	} else if enumerations.TimestampedObjectTypeRevocation == c.referenceCategory {
		return i18n.MessageTagPSVITPORDAOBCT
	}
	panic("Problem VTS")
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVITPOOBCTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *POEExistsAtOrBeforeControlTimeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoPOE
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/checks/RevocationIssuedBeforeControlTimeCheck.java (DSS 6.5.RC1).
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

// RevocationIssuedBeforeControlTimeCheck verifies if the issuance date of the
// revocation status information is before control time.
type RevocationIssuedBeforeControlTimeCheck[T any] struct {
	*process.ChainItemBase[T]

	// revocation is the revocation data to check.
	revocation *diagnostic.RevocationWrapper

	// controlTime is the control time.
	controlTime time.Time
}

// NewRevocationIssuedBeforeControlTimeCheck is the default constructor. Port of
// RevocationIssuedBeforeControlTimeCheck(I18nProvider, T, RevocationWrapper, Date, LevelRule).
func NewRevocationIssuedBeforeControlTimeCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	revocation *diagnostic.RevocationWrapper, controlTime time.Time,
	constraint policy.LevelRule) *RevocationIssuedBeforeControlTimeCheck[T] {
	c := &RevocationIssuedBeforeControlTimeCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		revocation:    revocation,
		controlTime:   controlTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationIssuedBeforeControlTimeCheck[T]) Process() bool {
	thisUpdate := c.revocation.ThisUpdate()
	return thisUpdate != nil && thisUpdate.Before(c.controlTime)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
//
// Java hands java.text.MessageFormat a literal null for a missing thisUpdate -
// not the empty string getFormattedDate(null) would produce - which renders as
// the four characters "null"; the Go argument is that text.
func (c *RevocationIssuedBeforeControlTimeCheck[T]) BuildAdditionalInfo() *string {
	thisUpdate := c.revocation.ThisUpdate()
	thisUpdateStr := "null"
	if thisUpdate != nil {
		thisUpdateStr = process.GetFormattedDate(thisUpdate)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagRevocationThisUpdateControlTime, c.revocation.Id(),
		thisUpdateStr, process.GetFormattedDate(&c.controlTime))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationIssuedBeforeControlTimeCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVHRDBIBCT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationIssuedBeforeControlTimeCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVHRDBIBCTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion(), which returns null.
func (c *RevocationIssuedBeforeControlTimeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return ""
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), which returns null.
func (c *RevocationIssuedBeforeControlTimeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

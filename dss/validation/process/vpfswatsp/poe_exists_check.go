// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/POEExistsCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// POEExistsCheck checks if the POE exists.
type POEExistsCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// token is the token to be validated.
	token diagnostic.TokenProxy

	// controlTime is the time when the signature validity can be proved; nil is
	// Java's null (the control time of a Past Certificate Validation that did
	// not determine one).
	controlTime *time.Time

	// poe is the set of available POEs.
	poe *POEExtraction
}

// NewPOEExistsCheck is the default constructor. Port of
// POEExistsCheck(I18nProvider, XmlPSV, TokenProxy, Date, POEExtraction, LevelRule).
func NewPOEExistsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlPSV],
	token diagnostic.TokenProxy, controlTime *time.Time, poe *POEExtraction,
	constraint policy.LevelRule) *POEExistsCheck {
	c := &POEExistsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
		controlTime:   controlTime,
		poe:           poe,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *POEExistsCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypePCV
}

// Process performs the check. Port of process(), which is public upstream: the
// PastSignatureValidation chain calls it on a second instance of this check
// before executing the chain.
func (c *POEExistsCheck) Process() bool {
	return c.controlTime != nil && c.poe.IsPOEExists(c.token.Id(), *c.controlTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *POEExistsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_ITPOSVAOBCT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *POEExistsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_ITPOSVAOBCT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion(), which returns null.
func (c *POEExistsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return ""
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), which returns null.
func (c *POEExistsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *POEExistsCheck) BuildAdditionalInfo() *string {
	lowestPOETime := c.poe.GetLowestPOETime(c.token.Id())
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CONTROL_TIME_WITH_POE,
		process.GetFormattedDate(c.controlTime), process.GetFormattedDate(&lowestPOETime))
	return &message
}

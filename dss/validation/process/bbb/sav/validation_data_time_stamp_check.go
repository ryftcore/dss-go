// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ValidationDataTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationDataTimeStampCheck checks if a validation-data-time-stamp attribute
// is present.
type ValidationDataTimeStampCheck struct {
	*AbstractTimeStampTypeCheck
}

// NewValidationDataTimeStampCheck is the default constructor. Port of
// ValidationDataTimeStampCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewValidationDataTimeStampCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *ValidationDataTimeStampCheck {
	c := &ValidationDataTimeStampCheck{
		AbstractTimeStampTypeCheck: NewAbstractTimeStampTypeCheck(i18nProvider, result, signature, constraint),
	}
	c.InitAbstractTimeStampTypeCheck(c)
	c.InitChainItem(c)
	return c
}

// TimestampType returns the associated TimestampType. Port of
// getTimestampType().
func (c *ValidationDataTimeStampCheck) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampTypeValidationDataTimestamp
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidationDataTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPVDTSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ValidationDataTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPVDTSPANS
}

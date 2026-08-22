// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ValidationDataRefsOnlyTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationDataRefsOnlyTimeStampCheck checks if a
// validation-data-refs-only-time-stamp attribute is present.
type ValidationDataRefsOnlyTimeStampCheck struct {
	*AbstractTimeStampTypeCheck
}

// NewValidationDataRefsOnlyTimeStampCheck is the default constructor. Port of
// ValidationDataRefsOnlyTimeStampCheck(Provider, XmlSAV, SignatureWrapper, LevelRule).
func NewValidationDataRefsOnlyTimeStampCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *ValidationDataRefsOnlyTimeStampCheck {
	c := &ValidationDataRefsOnlyTimeStampCheck{
		AbstractTimeStampTypeCheck: NewAbstractTimeStampTypeCheck(i18nProvider, result, signature, constraint),
	}
	c.InitAbstractTimeStampTypeCheck(c)
	c.InitChainItem(c)
	return c
}

// TimestampType returns the associated TimestampType. Port of
// getTimestampType().
func (c *ValidationDataRefsOnlyTimeStampCheck) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampTypeValidationDataRefsOnlyTimestamp
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ValidationDataRefsOnlyTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPVDROTSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ValidationDataRefsOnlyTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPVDROTSPANS
}

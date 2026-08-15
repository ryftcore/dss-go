// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SignatureTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureTimeStampCheck checks if a signature-time-stamp attribute is present.
type SignatureTimeStampCheck struct {
	*AbstractTimeStampTypeCheck
}

// NewSignatureTimeStampCheck is the default constructor. Port of
// SignatureTimeStampCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewSignatureTimeStampCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignatureTimeStampCheck {
	c := &SignatureTimeStampCheck{
		AbstractTimeStampTypeCheck: NewAbstractTimeStampTypeCheck(i18nProvider, result, signature, constraint),
	}
	c.InitAbstractTimeStampTypeCheck(c)
	c.InitChainItem(c)
	return c
}

// TimestampType returns the associated TimestampType. Port of
// getTimestampType().
func (c *SignatureTimeStampCheck) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampType_SIGNATURE_TIMESTAMP
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_IUQPSTSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignatureTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_IUQPSTSP_ANS
}

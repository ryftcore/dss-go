// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SignatureTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
	return enumerations.TimestampTypeSignatureTimestamp
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPSTSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignatureTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPSTSPANS
}

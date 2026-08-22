// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/DocumentTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DocumentTimeStampCheck checks if a document-time-stamp is present.
type DocumentTimeStampCheck struct {
	*AbstractTimeStampTypeCheck
}

// NewDocumentTimeStampCheck is the default constructor. Port of
// DocumentTimeStampCheck(Provider, XmlSAV, SignatureWrapper, LevelRule).
func NewDocumentTimeStampCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *DocumentTimeStampCheck {
	c := &DocumentTimeStampCheck{
		AbstractTimeStampTypeCheck: NewAbstractTimeStampTypeCheck(i18nProvider, result, signature, constraint),
	}
	c.InitAbstractTimeStampTypeCheck(c)
	c.InitChainItem(c)
	return c
}

// TimestampType returns the associated TimestampType. Port of
// getTimestampType().
func (c *DocumentTimeStampCheck) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampTypeDocumentTimestamp
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *DocumentTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIDTSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *DocumentTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIDTSPANS
}

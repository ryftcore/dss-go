// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/ArchiveTimeStampCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ArchiveTimeStampCheck checks if an archive-time-stamp attribute is present.
type ArchiveTimeStampCheck struct {
	*AbstractTimeStampTypeCheck
}

// NewArchiveTimeStampCheck is the default constructor. Port of
// ArchiveTimeStampCheck(Provider, XmlSAV, SignatureWrapper, LevelRule).
func NewArchiveTimeStampCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *ArchiveTimeStampCheck {
	c := &ArchiveTimeStampCheck{
		AbstractTimeStampTypeCheck: NewAbstractTimeStampTypeCheck(i18nProvider, result, signature, constraint),
	}
	c.InitAbstractTimeStampTypeCheck(c)
	c.InitChainItem(c)
	return c
}

// TimestampType returns the associated TimestampType. Port of
// getTimestampType().
func (c *ArchiveTimeStampCheck) TimestampType() enumerations.TimestampType {
	return enumerations.TimestampTypeArchiveTimestamp
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ArchiveTimeStampCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPATSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ArchiveTimeStampCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBSAVIUQPATSPANS
}

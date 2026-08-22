// Ported from dss-validation/.../validation/process/bbb/fc/checks/CAdESV3HashIndexCheck.java (DSS 6.5.RC1).
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CAdESV3HashIndexCheck verifies validity of the ats-hash-index(-v3) attribute present within an
// archive-time-stamp-v3 CAdES unsigned property, per EN 319 122-1.
type CAdESV3HashIndexCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	timestamp *diagnostic.TimestampWrapper
}

// NewCAdESV3HashIndexCheck is the default constructor.
func NewCAdESV3HashIndexCheck(i18nProvider *i18n.Provider, result *process.Result[*drjaxb.XmlFC],
	timestamp *diagnostic.TimestampWrapper, constraint policy.LevelRule) *CAdESV3HashIndexCheck {
	c := &CAdESV3HashIndexCheck{timestamp: timestamp}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *CAdESV3HashIndexCheck) Process() bool {
	if c.timestamp.Type() == enumerations.TimestampTypeArchiveTimestamp &&
		c.timestamp.ArchiveTimestampType() == enumerations.ArchiveTimestampTypeCAdESV3 {
		return c.timestamp.IsAtsHashIndexValid()
	}
	// accept for other timestamp types
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *CAdESV3HashIndexCheck) MessageTag() i18n.MessageTag { return i18n.MessageTagBBBFCIAHIV }

// ErrorMessageTag returns the error message i18n key.
func (c *CAdESV3HashIndexCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBFCIAHIVANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *CAdESV3HashIndexCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *CAdESV3HashIndexCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationFormatFailure
}

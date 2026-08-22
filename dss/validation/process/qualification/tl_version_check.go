// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLVersionCheck.java (DSS 6.5.RC1).
//
// See tl_not_expired_check.go for why currentTL is typed as the shared
// XmlTrustSourceListContent rather than XmlTrustSourceList.
package qualification

import (
	"strconv"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	dssjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// TLVersionCheck checks whether the version of the Trusted List is acceptable.
type TLVersionCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustSourceListContent

	// currentTime is the validation time.
	currentTime time.Time
}

// NewTLVersionCheck is the default constructor. Port of
// TLVersionCheck(Provider, XmlTLAnalysis, XmlTrustSourceList, Date, MultiValuesRule).
func NewTLVersionCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustSourceListContent, currentTime time.Time,
	constraint policy.MultiValuesRule) *TLVersionCheck {
	c := &TLVersionCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		currentTL:                    currentTL,
		currentTime:                  currentTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLVersionCheck) Process() bool {
	if !IsPostGracePeriod(&c.currentTime) {
		return true
	}
	tlVersion := c.currentTL.Version
	if tlVersion == nil {
		// invalid
		return false
	}
	return c.ProcessValueCheck(strconv.Itoa(*tlVersion))
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLVersionCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLVersion
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLVersionCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLVersionANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLVersionCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLVersionCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

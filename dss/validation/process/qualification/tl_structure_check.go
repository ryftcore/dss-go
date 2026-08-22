// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLStructureCheck.java (DSS 6.5.RC1).
//
// See tl_not_expired_check.go for why currentTL is typed as the shared
// XmlTrustSourceListContent rather than XmlTrustSourceList.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	dssjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLStructureCheck checks whether the structure of the Trusted List is valid.
type TLStructureCheck struct {
	*process.ChainItemBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustSourceListContent
}

// NewTLStructureCheck is the default constructor. Port of
// TLStructureCheck(Provider, XmlTLAnalysis, XmlTrustSourceList, LevelRule).
func NewTLStructureCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustSourceListContent, constraint policy.LevelRule) *TLStructureCheck {
	c := &TLStructureCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		currentTL:     currentTL,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLStructureCheck) Process() bool {
	return c.currentTL.StructuralValidation == nil || c.currentTL.StructuralValidation.Valid
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLStructureCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLSV
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLStructureCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLSVANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLStructureCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLStructureCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

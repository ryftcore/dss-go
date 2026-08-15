// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAATypeIntegrityPresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAATypeIntegrityPresentCheck verifies whether the SD-JWT EAA contains the
// claim "vct#integrity".
type EAATypeIntegrityPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAATypeIntegrityPresentCheck is the default constructor.
func NewEAATypeIntegrityPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAATypeIntegrityPresentCheck {
	c := &EAATypeIntegrityPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAATypeIntegrityPresentCheck) Process() bool {
	return c.eaa.EAAVerifiableCredentialsTypeIntegrityDigestAlgorithm() != "" && c.eaa.EAAVerifiableCredentialsTypeIntegrityBytes() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAATypeIntegrityPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_SDJWT_EAA_VCT_INT_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAATypeIntegrityPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_SDJWT_EAA_VCT_INT_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAATypeIntegrityPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAATypeIntegrityPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}

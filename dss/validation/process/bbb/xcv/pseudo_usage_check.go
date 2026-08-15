// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/PseudoUsageCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// PseudoUsageCheck checks if the certificate's pseudo usage is acceptable.
type PseudoUsageCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// pseudo is the used pseudo.
	pseudo string
}

// NewPseudoUsageCheck is the default constructor. Port of
// PseudoUsageCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewPseudoUsageCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *PseudoUsageCheck {
	c := &PseudoUsageCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PseudoUsageCheck) Process() bool {
	pseudoStrategy := NewJoinedPseudoStrategy()
	c.pseudo = pseudoStrategy.GetPseudo(c.certificate)
	return utils.IsStringEmpty(c.pseudo)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *PseudoUsageCheck) BuildAdditionalInfo() *string {
	if utils.IsStringNotEmpty(c.pseudo) {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_PSEUDO, c.pseudo)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PseudoUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_PSEUDO_USE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *PseudoUsageCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_PSEUDO_USE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PseudoUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *PseudoUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}

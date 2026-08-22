// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/PastCertificateValidationAcceptableCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PastCertificateValidationAcceptableCheck checks if the Past Certificate
// Validation result is acceptable.
type PastCertificateValidationAcceptableCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// pcv is the Past Certificate Validation.
	pcv *jaxb.XmlPCV

	// currentIndication is the current indication.
	currentIndication enumerations.Indication

	// currentSubIndication is the current subIndication.
	currentSubIndication enumerations.SubIndication
}

// NewPastCertificateValidationAcceptableCheck is the default constructor. Port
// of PastCertificateValidationAcceptableCheck(I18nProvider, XmlPSV, XmlPCV, String, Indication, SubIndication, LevelRule).
func NewPastCertificateValidationAcceptableCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlPSV],
	pcv *jaxb.XmlPCV, tokenId string, currentIndication enumerations.Indication,
	currentSubIndication enumerations.SubIndication,
	constraint policy.LevelRule) *PastCertificateValidationAcceptableCheck {
	c := &PastCertificateValidationAcceptableCheck{
		ChainItemBase:        process.NewChainItemBaseWithId(i18nProvider, result, constraint, tokenId),
		pcv:                  pcv,
		currentIndication:    currentIndication,
		currentSubIndication: currentSubIndication,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *PastCertificateValidationAcceptableCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypePCV
}

// Process performs the check. Port of process().
func (c *PastCertificateValidationAcceptableCheck) Process() bool {
	if c.pcv != nil && c.pcv.Conclusion != nil {
		pcvIndication := c.pcv.Conclusion.Indication.Indication()
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		var pcvSubIndication enumerations.SubIndication
		if c.pcv.Conclusion.SubIndication != nil {
			pcvSubIndication = c.pcv.Conclusion.SubIndication.SubIndication()
		}

		// INDETERMINATE cases are treated in following steps depending on POE
		return enumerations.IndicationPassed == pcvIndication ||
			(enumerations.IndicationIndeterminate == pcvIndication &&
				(enumerations.SubIndicationRevokedNoPOE == pcvSubIndication ||
					enumerations.SubIndicationRevokedCANoPOE == pcvSubIndication ||
					enumerations.SubIndicationOutOfBoundsNoPOE == pcvSubIndication ||
					enumerations.SubIndicationCryptoConstraintsFailureNoPOE == pcvSubIndication))
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PastCertificateValidationAcceptableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPCVA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *PastCertificateValidationAcceptableCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPCVA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PastCertificateValidationAcceptableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.currentIndication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *PastCertificateValidationAcceptableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.currentSubIndication
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose else branch delegates to the base.
func (c *PastCertificateValidationAcceptableCheck) BuildAdditionalInfo() *string {
	if c.pcv != nil && c.pcv.ControlTime != nil {
		controlTime := c.pcv.ControlTime.Time()
		message := c.I18nProvider.GetMessage(i18n.MessageTag_CONTROL_TIME_ALONE,
			process.GetFormattedDate(&controlTime))
		return &message
	}
	return c.ChainItemBase.BuildAdditionalInfo()
}

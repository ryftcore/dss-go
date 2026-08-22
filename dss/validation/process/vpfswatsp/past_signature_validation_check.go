// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/PastSignatureValidationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PastSignatureValidationCheck checks if the past signature validation result
// is acceptable.
type PastSignatureValidationCheck struct {
	*AbstractPastTokenValidationCheck[*jaxb.XmlValidationProcessArchivalData]
}

// NewPastSignatureValidationCheck is the default constructor. Port of
// PastSignatureValidationCheck(I18nProvider, XmlValidationProcessArchivalData, SignatureWrapper, XmlPSV, LevelRule).
func NewPastSignatureValidationCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessArchivalData], signature *diagnostic.SignatureWrapper,
	xmlPSV *jaxb.XmlPSV, constraint policy.LevelRule) *PastSignatureValidationCheck {
	c := &PastSignatureValidationCheck{
		AbstractPastTokenValidationCheck: NewAbstractPastTokenValidationCheck(
			i18nProvider, result, signature, xmlPSV, constraint),
	}
	// Re-register with the outer type so the overridden methods below dispatch
	// correctly.
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *PastSignatureValidationCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_PSV
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PastSignatureValidationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPSVC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *PastSignatureValidationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPSVC_ANS
}

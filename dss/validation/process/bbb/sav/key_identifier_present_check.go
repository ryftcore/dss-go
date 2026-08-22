// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/KeyIdentifierPresentCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// KeyIdentifierPresentCheck verifies whether a 'kid' (key identifier) header
// parameter is present within the protected header of a signature.
type KeyIdentifierPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to verify.
	signature *diagnostic.SignatureWrapper
}

// NewKeyIdentifierPresentCheck is the default constructor. Port of
// KeyIdentifierPresentCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewKeyIdentifierPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *KeyIdentifierPresentCheck {
	c := &KeyIdentifierPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *KeyIdentifierPresentCheck) Process() bool {
	return c.signature.KeyIdentifierReference() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *KeyIdentifierPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAKIDP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *KeyIdentifierPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAKIDP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *KeyIdentifierPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *KeyIdentifierPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}

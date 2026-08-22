// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/KeyIdentifierMatchCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// KeyIdentifierMatchCheck verifies whether a value of the signed attribute
// 'kid' (key identifier), when present, matches the signing-certificate used to
// create the signature.
type KeyIdentifierMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to verify.
	signature *diagnostic.SignatureWrapper
}

// NewKeyIdentifierMatchCheck is the default constructor. Port of
// KeyIdentifierMatchCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewKeyIdentifierMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *KeyIdentifierMatchCheck {
	c := &KeyIdentifierMatchCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *KeyIdentifierMatchCheck) Process() bool {
	keyIdentifierReference := c.signature.KeyIdentifierReference()
	if keyIdentifierReference != nil {
		return keyIdentifierReference.IsIssuerSerialMatch()
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *KeyIdentifierMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSDKIDVM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *KeyIdentifierMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSDKIDVMANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *KeyIdentifierMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *KeyIdentifierMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}

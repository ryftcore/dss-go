// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/MessageDigestOrSignedPropertiesCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// MessageDigestOrSignedPropertiesCheck checks if message-digest (CAdES/PAdES) or
// SignedProperties (XAdES) is present.
type MessageDigestOrSignedPropertiesCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewMessageDigestOrSignedPropertiesCheck is the default constructor. Port of
// MessageDigestOrSignedPropertiesCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewMessageDigestOrSignedPropertiesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *MessageDigestOrSignedPropertiesCheck {
	c := &MessageDigestOrSignedPropertiesCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *MessageDigestOrSignedPropertiesCheck) Process() bool {
	signatureForm, err := c.signature.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}
	switch signatureForm {
	case enumerations.SignatureForm_XAdES:
		return c.isRequiredDigestMatcherPresent(enumerations.DigestMatcherType_SIGNED_PROPERTIES)
	case enumerations.SignatureForm_CAdES, enumerations.SignatureForm_PAdES, enumerations.SignatureForm_PKCS7:
		return c.isRequiredDigestMatcherPresent(enumerations.DigestMatcherType_MESSAGE_DIGEST)
	default:
		// JAdES/CB-AdES shall be skipped
		return false
	}
}

// isRequiredDigestMatcherPresent ports the private
// isRequiredDigestMatcherPresent(DigestMatcherType).
func (c *MessageDigestOrSignedPropertiesCheck) isRequiredDigestMatcherPresent(digestMatcherType enumerations.DigestMatcherType) bool {
	digestMatchers := c.signature.DigestMatchers()
	if utils.IsCollectionNotEmpty(digestMatchers) {
		for _, digestMatcher := range digestMatchers {
			if digestMatcherType == messageDigestMatcherType(digestMatcher) {
				return true
			}
		}
	}
	return false
}

// messageDigestMatcherType reads XmlDigestMatcher#getType(): the generated
// member is a pointer, whose nil is Java's null.
func messageDigestMatcherType(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *MessageDigestOrSignedPropertiesCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPMDOSPP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *MessageDigestOrSignedPropertiesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ISQPMDOSPP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *MessageDigestOrSignedPropertiesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *MessageDigestOrSignedPropertiesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}

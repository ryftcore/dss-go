// Ported from dss-validation/.../validation/process/bbb/fc/checks/ReferencesNotAmbiguousCheck.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// ReferencesNotAmbiguousCheck checks if the references are not ambiguous (only one document is retrieved).
type ReferencesNotAmbiguousCheck struct {
	*process.ChainItemBase[*drjaxb.XmlFC]

	signature *diagnostic.SignatureWrapper

	duplicatedReference *diagjaxb.XmlDigestMatcher
}

// NewReferencesNotAmbiguousCheck is the default constructor.
func NewReferencesNotAmbiguousCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*drjaxb.XmlFC],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *ReferencesNotAmbiguousCheck {
	c := &ReferencesNotAmbiguousCheck{signature: signature}
	c.ChainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	c.InitChainItem(c)
	return c
}

// Process performs the check.
func (c *ReferencesNotAmbiguousCheck) Process() bool {
	for _, digestMatcher := range c.signature.DigestMatchers() {
		if digestMatcher.Duplicated != nil && *digestMatcher.Duplicated {
			c.duplicatedReference = digestMatcher
			return false
		}
	}
	return true
}

// MessageTag returns the constraint message i18n key.
func (c *ReferencesNotAmbiguousCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISRIA
}

// ErrorMessageTag returns the error message i18n key.
func (c *ReferencesNotAmbiguousCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_FC_ISRIA_ANS
}

// FailedIndicationForConclusion returns the Indication on failure.
func (c *ReferencesNotAmbiguousCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion returns the SubIndication on failure.
func (c *ReferencesNotAmbiguousCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}

// BuildAdditionalInfo builds the additional info message identifying the
// duplicated reference. Port of the overridden protected String buildAdditionalInfo().
func (c *ReferencesNotAmbiguousCheck) BuildAdditionalInfo() *string {
	if c.duplicatedReference == nil {
		return nil
	}
	referenceName := ""
	if c.duplicatedReference.Id != nil {
		referenceName = *c.duplicatedReference.Id
	}
	if strings.TrimSpace(referenceName) == "" && c.duplicatedReference.Type != nil {
		referenceName = string(c.duplicatedReference.Type.DigestMatcherType())
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_REFERENCE, referenceName)
	return &message
}

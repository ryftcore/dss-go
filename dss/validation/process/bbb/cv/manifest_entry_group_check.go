// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ManifestEntryGroupCheck.java (DSS 6.5.RC1).
package cv

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// ManifestEntryGroupCheck verifies whether all the manifest entries have been
// found during the validation process.
type ManifestEntryGroupCheck struct {
	*process.ChainItemBase[*jaxb.XmlCV]

	// digestMatchers is the digest matchers to check.
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher
}

// NewManifestEntryGroupCheck is the default constructor. Port of
// ManifestEntryGroupCheck(I18nProvider, XmlCV, List, LevelRule).
func NewManifestEntryGroupCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlCV],
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ManifestEntryGroupCheck {
	c := &ManifestEntryGroupCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatchers: digestMatchers,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ManifestEntryGroupCheck) Process() bool {
	for _, d := range c.digestMatchers {
		if enumerations.DigestMatcherType_MANIFEST_ENTRY == digestMatcherType(d) && !d.DataFound {
			return false
		}
	}
	return true
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ManifestEntryGroupCheck) BuildAdditionalInfo() *string {
	var notFoundNames []string
	for _, d := range c.digestMatchers {
		if enumerations.DigestMatcherType_MANIFEST_ENTRY == digestMatcherType(d) && !d.DataFound {
			if d.Uri != nil {
				notFoundNames = append(notFoundNames, *d.Uri)
			}
		}
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_REFERENCES_WITH_NAMES,
		utils.JoinStrings(notFoundNames, ", "))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ManifestEntryGroupCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_AAMEF
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ManifestEntryGroupCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_CV_AAMEF_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ManifestEntryGroupCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ManifestEntryGroupCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIGNED_DATA_NOT_FOUND
}

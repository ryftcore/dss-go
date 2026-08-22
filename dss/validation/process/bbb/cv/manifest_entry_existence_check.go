// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ManifestEntryExistenceCheck.java (DSS 6.5.RC1).
package cv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ManifestEntryExistenceCheck checks if at least one manifest entry is present.
type ManifestEntryExistenceCheck struct {
	*process.ChainItemBase[*jaxb.XmlCV]

	// digestMatchers is the digest matchers to check.
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher
}

// NewManifestEntryExistenceCheck is the default constructor. Port of
// ManifestEntryExistenceCheck(Provider, XmlCV, List, LevelRule).
func NewManifestEntryExistenceCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlCV],
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ManifestEntryExistenceCheck {
	c := &ManifestEntryExistenceCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatchers: digestMatchers,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ManifestEntryExistenceCheck) Process() bool {
	for _, xmlDigestMatcher := range c.digestMatchers {
		if enumerations.DigestMatcherTypeManifestEntry == digestMatcherType(xmlDigestMatcher) &&
			xmlDigestMatcher.DataFound {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ManifestEntryExistenceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBCVISMEC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ManifestEntryExistenceCheck) ErrorMessageTag() i18n.MessageTag {
	var manifestEntries []*diagnosticjaxb.XmlDigestMatcher
	for _, d := range c.digestMatchers {
		if enumerations.DigestMatcherTypeManifestEntry == digestMatcherType(d) {
			manifestEntries = append(manifestEntries, d)
		}
	}
	if utils.IsCollectionNotEmpty(manifestEntries) && noneDataFound(manifestEntries) {
		return i18n.MessageTagBBBCVISMECANS2
	}
	return i18n.MessageTagBBBCVISMECANS
}

// noneDataFound ports Stream#noneMatch(XmlDigestMatcher::isDataFound).
func noneDataFound(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) bool {
	for _, d := range digestMatchers {
		if d.DataFound {
			return false
		}
	}
	return true
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ManifestEntryExistenceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ManifestEntryExistenceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignedDataNotFound
}

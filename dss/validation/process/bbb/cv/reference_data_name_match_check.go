// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ReferenceDataNameMatchCheck.java (DSS 6.5.RC1).
package cv

import (
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ReferenceDataNameMatchCheck checks if the referenced document name matches the
// reference name.
type ReferenceDataNameMatchCheck[T any] struct {
	*process.ChainItemBase[T]

	// digestMatcher is the reference DigestMatcher.
	digestMatcher *diagnosticjaxb.XmlDigestMatcher
}

// NewReferenceDataNameMatchCheck is the default constructor. Port of
// ReferenceDataNameMatchCheck(Provider, T, XmlDigestMatcher, LevelRule).
func NewReferenceDataNameMatchCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	digestMatcher *diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ReferenceDataNameMatchCheck[T] {
	c := &ReferenceDataNameMatchCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatcher: digestMatcher,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ReferenceDataNameMatchCheck[T]) Process() bool {
	return c.digestMatcher.Uri != nil && *c.digestMatcher.Uri == digestMatcherDocumentName(c.digestMatcher)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
//
// Java guards the type with Objects.requireNonNull, which throws a
// NullPointerException on a typeless digest matcher; the Go zero value simply
// takes the non-MANIFEST_ENTRY branch.
func (c *ReferenceDataNameMatchCheck[T]) MessageTag() i18n.MessageTag {
	if digestMatcherType(c.digestMatcher) == enumerations.DigestMatcherTypeManifestEntry {
		return i18n.MessageTagBBBCVDMENMND
	}
	return i18n.MessageTagBBBCVDRNMND
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataNameMatchCheck[T]) ErrorMessageTag() i18n.MessageTag {
	if digestMatcherType(c.digestMatcher) == enumerations.DigestMatcherTypeManifestEntry {
		return i18n.MessageTagBBBCVDMENMNDANS
	}
	return i18n.MessageTagBBBCVDRNMNDANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ReferenceDataNameMatchCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ReferenceDataNameMatchCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignedDataNotFound
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ReferenceDataNameMatchCheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagReferenceNameCheck,
		digestMatcherUri(c.digestMatcher), digestMatcherDocumentName(c.digestMatcher))
	return &message
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ReferenceDataExistenceCheck.java (DSS 6.5.RC1).
package cv

import (
	"fmt"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ReferenceDataExistenceCheck checks if the referenced data is found.
type ReferenceDataExistenceCheck[T any] struct {
	*process.ChainItemBase[T]

	// digestMatcher is the reference DigestMatcher.
	digestMatcher *diagnosticjaxb.XmlDigestMatcher
}

// NewReferenceDataExistenceCheck is the default constructor. Port of
// ReferenceDataExistenceCheck(I18nProvider, T, XmlDigestMatcher, LevelRule).
func NewReferenceDataExistenceCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	digestMatcher *diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ReferenceDataExistenceCheck[T] {
	c := &ReferenceDataExistenceCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatcher: digestMatcher,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ReferenceDataExistenceCheck[T]) Process() bool {
	return c.digestMatcher.DataFound
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ReferenceDataExistenceCheck[T]) MessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTagBBBCVTSPIRDOF
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTagBBBCVCSCSSVF
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTagBBBCVIMEOF
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTagBBBCVERATSRF
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTagBBBCVERATSSRF
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTagBBBCVEAASDCBF
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTagBBBCVEAANSDCBF
	default:
		return i18n.MessageTagBBBCVIRDOF
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataExistenceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTagBBBCVTSPIRDOFANS
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTagBBBCVCSCSSVFANS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTagBBBCVIMEOFANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTagBBBCVERATSRFANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTagBBBCVERATSSRFANS
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTagBBBCVEAASDCBFANS
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTagBBBCVEAANSDCBFANS
	default:
		return i18n.MessageTagBBBCVIRDOFANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ReferenceDataExistenceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ReferenceDataExistenceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignedDataNotFound
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ReferenceDataExistenceCheck[T]) BuildAdditionalInfo() *string {
	var referenceName interface{}
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint,
		enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return nil
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		referenceName = i18n.MessageTagTSTTypeRefERATST
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		referenceName = i18n.MessageTagTSTTypeRefERATSTSeq
	default:
		referenceName = c.getReferenceName(c.digestMatcher)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTagReference, referenceName)
	return &message
}

// getReferenceName ports the private getReferenceName(XmlDigestMatcher).
func (c *ReferenceDataExistenceCheck[T]) getReferenceName(digestMatcher *diagnosticjaxb.XmlDigestMatcher) string {
	if utils.IsStringNotBlank(digestMatcherId(digestMatcher)) {
		return digestMatcherId(digestMatcher)
	} else if utils.IsStringNotBlank(digestMatcherUri(digestMatcher)) {
		return digestMatcherUri(digestMatcher)
	} else if digestMatcher.DisclosableClaim != nil && digestMatcher.DisclosableClaim.Name != nil {
		claimName := *digestMatcher.DisclosableClaim.Name
		if enumerations.DigestMatcherTypeEAANestedDisclosure == digestMatcherType(digestMatcher) &&
			digestMatcher.DisclosableClaim.Value != "" {
			claimName += fmt.Sprintf(" '%s'", digestMatcher.DisclosableClaim.Value)
		}
		return claimName
	} else {
		return string(digestMatcherType(digestMatcher))
	}
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ReferenceDataIntactCheck.java (DSS 6.5.RC1).
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

// ReferenceDataIntactCheck checks if the referenced data is intact.
type ReferenceDataIntactCheck[T any] struct {
	*process.ChainItemBase[T]

	// digestMatcher is the reference DigestMatcher.
	digestMatcher *diagnosticjaxb.XmlDigestMatcher
}

// NewReferenceDataIntactCheck is the default constructor. Port of
// ReferenceDataIntactCheck(I18nProvider, T, XmlDigestMatcher, LevelRule).
func NewReferenceDataIntactCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	digestMatcher *diagnosticjaxb.XmlDigestMatcher, constraint policy.LevelRule) *ReferenceDataIntactCheck[T] {
	c := &ReferenceDataIntactCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		digestMatcher: digestMatcher,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ReferenceDataIntactCheck[T]) Process() bool {
	return c.digestMatcher.DataIntact
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ReferenceDataIntactCheck[T]) MessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTagBBBCVTSPIRDOI
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTagBBBCVCSCSPS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTagBBBCVIMEDOI
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTagBBBCVERATSRI
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTagBBBCVERATSSRI
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTagBBBCVEAASDCBI
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTagBBBCVEAANSDCBI
	default:
		return i18n.MessageTagBBBCVIRDOI
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataIntactCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTagBBBCVTSPIRDOIANS
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTagBBBCVCSCSPSANS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTagBBBCVIMEDOIANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTagBBBCVERATSRIANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTagBBBCVERATSSRIANS
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTagBBBCVEAASDCBIANS
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTagBBBCVEAANSDCBIANS
	default:
		return i18n.MessageTagBBBCVIRDOIANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ReferenceDataIntactCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ReferenceDataIntactCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationHashFailure
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ReferenceDataIntactCheck[T]) BuildAdditionalInfo() *string {
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
func (c *ReferenceDataIntactCheck[T]) getReferenceName(digestMatcher *diagnosticjaxb.XmlDigestMatcher) string {
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

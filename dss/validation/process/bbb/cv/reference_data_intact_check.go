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
		return i18n.MessageTag_BBB_CV_TSP_IRDOI
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTag_BBB_CV_CS_CSPS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTag_BBB_CV_IMEDOI
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTag_BBB_CV_ER_ATSRI
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTag_BBB_CV_ER_ATSSRI
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTag_BBB_CV_EAA_SDCBI
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBI
	default:
		return i18n.MessageTag_BBB_CV_IRDOI
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataIntactCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTag_BBB_CV_TSP_IRDOI_ANS
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTag_BBB_CV_CS_CSPS_ANS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTag_BBB_CV_IMEDOI_ANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTag_BBB_CV_ER_ATSRI_ANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTag_BBB_CV_ER_ATSSRI_ANS
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTag_BBB_CV_EAA_SDCBI_ANS
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBI_ANS
	default:
		return i18n.MessageTag_BBB_CV_IRDOI_ANS
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
		referenceName = i18n.MessageTag_TST_TYPE_REF_ER_ATST
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		referenceName = i18n.MessageTag_TST_TYPE_REF_ER_ATST_SEQ
	default:
		referenceName = c.getReferenceName(c.digestMatcher)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_REFERENCE, referenceName)
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

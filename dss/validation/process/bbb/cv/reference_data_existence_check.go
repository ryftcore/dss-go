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
		return i18n.MessageTag_BBB_CV_TSP_IRDOF
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTag_BBB_CV_CS_CSSVF
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTag_BBB_CV_IMEOF
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTag_BBB_CV_ER_ATSRF
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTag_BBB_CV_ER_ATSSRF
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTag_BBB_CV_EAA_SDCBF
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBF
	default:
		return i18n.MessageTag_BBB_CV_IRDOF
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataExistenceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherTypeMessageImprint:
		return i18n.MessageTag_BBB_CV_TSP_IRDOF_ANS
	case enumerations.DigestMatcherTypeCounterSignedSignatureValue:
		return i18n.MessageTag_BBB_CV_CS_CSSVF_ANS
	case enumerations.DigestMatcherTypeManifestEntry:
		return i18n.MessageTag_BBB_CV_IMEOF_ANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp:
		return i18n.MessageTag_BBB_CV_ER_ATSRF_ANS
	case enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
		return i18n.MessageTag_BBB_CV_ER_ATSSRF_ANS
	case enumerations.DigestMatcherTypeEAADisclosure:
		return i18n.MessageTag_BBB_CV_EAA_SDCBF_ANS
	case enumerations.DigestMatcherTypeEAANestedDisclosure:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBF_ANS
	default:
		return i18n.MessageTag_BBB_CV_IRDOF_ANS
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

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/ReferenceDataExistenceCheck.java (DSS 6.5.RC1).
package cv

import (
	"fmt"

	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
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
	case enumerations.DigestMatcherType_MESSAGE_IMPRINT:
		return i18n.MessageTag_BBB_CV_TSP_IRDOF
	case enumerations.DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE:
		return i18n.MessageTag_BBB_CV_CS_CSSVF
	case enumerations.DigestMatcherType_MANIFEST_ENTRY:
		return i18n.MessageTag_BBB_CV_IMEOF
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP:
		return i18n.MessageTag_BBB_CV_ER_ATSRF
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE:
		return i18n.MessageTag_BBB_CV_ER_ATSSRF
	case enumerations.DigestMatcherType_EAA_DISCLOSURE:
		return i18n.MessageTag_BBB_CV_EAA_SDCBF
	case enumerations.DigestMatcherType_EAA_NESTED_DISCLOSURE:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBF
	default:
		return i18n.MessageTag_BBB_CV_IRDOF
	}
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ReferenceDataExistenceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherType_MESSAGE_IMPRINT:
		return i18n.MessageTag_BBB_CV_TSP_IRDOF_ANS
	case enumerations.DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE:
		return i18n.MessageTag_BBB_CV_CS_CSSVF_ANS
	case enumerations.DigestMatcherType_MANIFEST_ENTRY:
		return i18n.MessageTag_BBB_CV_IMEOF_ANS
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP:
		return i18n.MessageTag_BBB_CV_ER_ATSRF_ANS
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE:
		return i18n.MessageTag_BBB_CV_ER_ATSSRF_ANS
	case enumerations.DigestMatcherType_EAA_DISCLOSURE:
		return i18n.MessageTag_BBB_CV_EAA_SDCBF_ANS
	case enumerations.DigestMatcherType_EAA_NESTED_DISCLOSURE:
		return i18n.MessageTag_BBB_CV_EAA_NSDCBF_ANS
	default:
		return i18n.MessageTag_BBB_CV_IRDOF_ANS
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ReferenceDataExistenceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *ReferenceDataExistenceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIGNED_DATA_NOT_FOUND
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *ReferenceDataExistenceCheck[T]) BuildAdditionalInfo() *string {
	var referenceName interface{}
	switch digestMatcherType(c.digestMatcher) {
	case enumerations.DigestMatcherType_MESSAGE_IMPRINT,
		enumerations.DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE:
		return nil
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP:
		referenceName = i18n.MessageTag_TST_TYPE_REF_ER_ATST
	case enumerations.DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE:
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
		if enumerations.DigestMatcherType_EAA_NESTED_DISCLOSURE == digestMatcherType(digestMatcher) &&
			digestMatcher.DisclosableClaim.Value != "" {
			claimName += fmt.Sprintf(" '%s'", digestMatcher.DisclosableClaim.Value)
		}
		return claimName
	} else {
		return string(digestMatcherType(digestMatcher))
	}
}

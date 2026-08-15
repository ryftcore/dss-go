// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceStatusConsistencyCheck.java (DSS 6.5.RC1).
//
// FLAGGED HASH-ORDER SITE: getApplicableStatusesSet() collects into a
// HashSet<String> in Java; sorted lexicographically here instead of
// insertion order for a deterministic (if not necessarily Java-bucket-
// identical) additional-info text on the multi-status branch. See the
// porter brief's hard rule on hash-order leaks.
package qualification

import (
	"sort"
	"strings"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// TrustedEntityServiceStatusConsistencyCheck verifies whether the Trusted
// Entity Service statuses are consistent.
type TrustedEntityServiceStatusConsistencyCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateApprovalStatus]

	// trustedServicesWithSti is the list of TrustedEntityServiceWrappers at
	// control time.
	trustedServicesWithSti []*diagnostic.TrustedEntityServiceWrapper
}

// NewTrustedEntityServiceStatusConsistencyCheck is the default constructor.
// Port of
// TrustedEntityServiceStatusConsistencyCheck(I18nProvider, XmlValidationCertificateApprovalStatus, List, LevelRule).
func NewTrustedEntityServiceStatusConsistencyCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateApprovalStatus], trustedServicesWithSti []*diagnostic.TrustedEntityServiceWrapper,
	constraint policy.LevelRule) *TrustedEntityServiceStatusConsistencyCheck {
	c := &TrustedEntityServiceStatusConsistencyCheck{
		ChainItemBase:          process.NewChainItemBase(i18nProvider, result, constraint),
		trustedServicesWithSti: trustedServicesWithSti,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedEntityServiceStatusConsistencyCheck) Process() bool {
	return utils.CollectionSize(c.getApplicableStatusesSet()) <= 1
}

// getApplicableStatusesSet ports the private getApplicableStatusesSet().
func (c *TrustedEntityServiceStatusConsistencyCheck) getApplicableStatusesSet() []string {
	seen := map[string]struct{}{}
	var ordered []string
	for _, ts := range c.trustedServicesWithSti {
		if _, ok := seen[ts.Status]; !ok {
			seen[ts.Status] = struct{}{}
			ordered = append(ordered, ts.Status)
		}
	}
	sort.Strings(ordered)
	return ordered
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceStatusConsistencyCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STATUS_CONS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceStatusConsistencyCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STATUS_CONS_ANS
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *TrustedEntityServiceStatusConsistencyCheck) BuildAdditionalInfo() *string {
	applicableStatusesSet := c.getApplicableStatusesSet()
	var message string
	if len(applicableStatusesSet) == 1 {
		message = c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_STATUS, applicableStatusesSet[0])
	} else {
		message = c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_STATUSES, "["+strings.Join(applicableStatusesSet, ", ")+"]")
	}
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceStatusConsistencyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceStatusConsistencyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

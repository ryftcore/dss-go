// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceStatusConsistencyCheck.java (DSS 6.5.RC1).
//
// HASH-ORDER (closed in phase 8f): getApplicableStatusesSet() collects into a
// HashSet<String> in Java and buildAdditionalInfo() renders that set with
// Set#toString(), so the iteration order reaches the report on the
// multi-status branch. utils.JavaHashMapStringKeyOrder reproduces it.
package qualification

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
	// Java collects into a java.util.HashSet via Collectors.toSet() and renders
	// it with Set#toString() below, so the iteration order reaches the report's
	// <AdditionalInfo>; reproduce it rather than sorting.
	return utils.JavaHashMapStringKeyOrder(ordered)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceStatusConsistencyCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageStatusCONS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceStatusConsistencyCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageStatusCONSANS
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *TrustedEntityServiceStatusConsistencyCheck) BuildAdditionalInfo() *string {
	applicableStatusesSet := c.getApplicableStatusesSet()
	rendered := make([]string, len(applicableStatusesSet))
	for i, status := range applicableStatusesSet {
		rendered[i] = renderJavaNullableStatus(status)
	}
	var message string
	if len(rendered) == 1 {
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateUsageStatus, rendered[0])
	} else {
		message = c.I18nProvider.GetMessage(i18n.MessageTagCertificateUsageStatuses, "["+strings.Join(rendered, ", ")+"]")
	}
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceStatusConsistencyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceStatusConsistencyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

// renderJavaNullableStatus renders a trusted-entity service status the way
// Java's MessageFormat does when the status is absent. TrustedEntityServiceWrapper#
// getStatus() is a nullable String upstream, and an absent <Status> element
// therefore formats into CERTIFICATE_USAGE_STATUS's {0} as the literal "null";
// this port carries an absent status as the empty string (the convention
// trusted_entity_service_by_status_filter.go already relies on), which would
// otherwise render as nothing at all. Found by the phase-8f full-corpus report
// byte-parity run on eaa-validation/diag_data_pid.xml.
func renderJavaNullableStatus(status string) string {
	if status == "" {
		return "null"
	}
	return status
}

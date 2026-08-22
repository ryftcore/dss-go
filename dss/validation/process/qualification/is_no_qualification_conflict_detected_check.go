// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/IsNoQualificationConflictDetectedCheck.java (DSS 6.5.RC1).
//
// Java collects the simulated per-TrustService
// qualification outcomes into a Set<CertificateQualification> (a HashSet);
// CertificateQualification is a plain Java enum with no overridden
// hashCode(), so its HashSet iteration order is JVM-identity-hash dependent
// and not reproducible from the Go side at all. certificateQualificationsAtTime
// here is instead the caller-supplied slice in first-seen order (the caller,
// CertQualificationAtTimeBlock, is also ported in this package); the size
// check (>1 distinct values) is order-independent and therefore exact, but
// the RESULTS additional-info text (rendered only on the >1 branch, which
// aborts the qualification chain into CertificateQualificationNA regardless
// of the exact set of values) may list the conflicting values in a different
// order than upstream.
package qualification

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// IsNoQualificationConflictDetectedCheck verifies if there is no conflict in
// certificate qualification determination result based on a use of different
// TrustServices.
type IsNoQualificationConflictDetectedCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateQualification]

	// certificateQualificationsAtTime is the set of obtained
	// CertificateQualifications from various TrustServices, in first-seen
	// order (see file header).
	certificateQualificationsAtTime []enumerations.CertificateQualification
}

// NewIsNoQualificationConflictDetectedCheck is the default constructor. Port
// of IsNoQualificationConflictDetectedCheck(Provider, XmlValidationCertificateQualification, Set, LevelRule).
func NewIsNoQualificationConflictDetectedCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlValidationCertificateQualification],
	certificateQualificationsAtTime []enumerations.CertificateQualification,
	constraint policy.LevelRule) *IsNoQualificationConflictDetectedCheck {
	c := &IsNoQualificationConflictDetectedCheck{
		ChainItemBase:                   process.NewChainItemBase(i18nProvider, result, constraint),
		certificateQualificationsAtTime: certificateQualificationsAtTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *IsNoQualificationConflictDetectedCheck) Process() bool {
	return utils.CollectionSize(c.certificateQualificationsAtTime) == 1
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *IsNoQualificationConflictDetectedCheck) BuildAdditionalInfo() *string {
	if utils.CollectionSize(c.certificateQualificationsAtTime) > 1 {
		values := make([]string, 0, len(c.certificateQualificationsAtTime))
		for _, v := range c.certificateQualificationsAtTime {
			values = append(values, string(v))
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTagResults, "["+strings.Join(values, ", ")+"]")
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *IsNoQualificationConflictDetectedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasConf
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *IsNoQualificationConflictDetectedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualHasConfANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *IsNoQualificationConflictDetectedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *IsNoQualificationConflictDetectedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

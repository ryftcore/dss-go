// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateNotOnHoldCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateNotOnHoldCheck checks if the certificate is not on hold.
type CertificateNotOnHoldCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificateRevocation is the certificate's revocation.
	certificateRevocation *diagnostic.CertificateRevocationWrapper

	// currentTime is the validation time.
	currentTime time.Time
}

// NewCertificateNotOnHoldCheck is the default constructor. Port of
// CertificateNotOnHoldCheck(Provider, XmlSubXCV, CertificateRevocationWrapper, Date, LevelRule).
func NewCertificateNotOnHoldCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificateRevocation *diagnostic.CertificateRevocationWrapper, currentTime time.Time,
	constraint policy.LevelRule) *CertificateNotOnHoldCheck {
	c := &CertificateNotOnHoldCheck{
		ChainItemBase:         process.NewChainItemBase(i18nProvider, result, constraint),
		certificateRevocation: certificateRevocation,
		currentTime:           currentTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateNotOnHoldCheck) Process() bool {
	isOnHold := c.certificateRevocation != nil && c.certificateRevocation.IsRevoked() &&
		enumerations.RevocationReasonCertificateHold == c.certificateRevocation.Reason()
	if isOnHold {
		revocationDate := c.certificateRevocation.RevocationDate()
		isOnHold = revocationDate != nil && !c.currentTime.Before(*revocationDate)
	}
	return !isOnHold
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *CertificateNotOnHoldCheck) BuildAdditionalInfo() *string {
	if c.certificateRevocation != nil && c.certificateRevocation.RevocationDate() != nil {
		revocationDateStr := process.GetFormattedDate(c.certificateRevocation.RevocationDate())
		// Java hands MessageFormat the RevocationReason ENUM, so a NULL reason
		// renders as the literal text "null"; the ported wrapper returns the
		// empty string for that null, which would render as nothing.
		reason := "null"
		if r := c.certificateRevocation.Reason(); r != "" {
			reason = string(r)
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTagRevocationReason, reason, revocationDateStr)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateNotOnHoldCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCOH
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateNotOnHoldCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCOHANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateNotOnHoldCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateNotOnHoldCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTryLater
}

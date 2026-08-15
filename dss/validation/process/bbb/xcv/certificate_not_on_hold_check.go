// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateNotOnHoldCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
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
// CertificateNotOnHoldCheck(I18nProvider, XmlSubXCV, CertificateRevocationWrapper, Date, LevelRule).
func NewCertificateNotOnHoldCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificateRevocation *diagnostic.CertificateRevocationWrapper, currentTime time.Time,
	constraint policy.LevelRule) *CertificateNotOnHoldCheck {
	c := &CertificateNotOnHoldCheck{
		ChainItemBase:          process.NewChainItemBase(i18nProvider, result, constraint),
		certificateRevocation:  certificateRevocation,
		currentTime:            currentTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateNotOnHoldCheck) Process() bool {
	isOnHold := c.certificateRevocation != nil && c.certificateRevocation.IsRevoked() &&
		enumerations.RevocationReason_CERTIFICATE_HOLD == c.certificateRevocation.Reason()
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
		message := c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_REASON, c.certificateRevocation.Reason(), revocationDateStr)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateNotOnHoldCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCOH
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateNotOnHoldCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCOH_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateNotOnHoldCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateNotOnHoldCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_TRY_LATER
}

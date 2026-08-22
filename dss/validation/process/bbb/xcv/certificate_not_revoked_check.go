// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateNotRevokedCheck.java (DSS 6.5.RC1).
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

// CertificateNotRevokedCheck checks if the certificate is not revoked.
type CertificateNotRevokedCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificateRevocation is the certificate's revocation.
	certificateRevocation *diagnostic.CertificateRevocationWrapper

	// currentTime is the validation time.
	currentTime time.Time

	// subContext is the validation subContext.
	subContext enumerations.SubContext
}

// NewCertificateNotRevokedCheck is the default constructor. Port of
// CertificateNotRevokedCheck(I18nProvider, XmlSubXCV, CertificateRevocationWrapper, Date, LevelRule, SubContext).
func NewCertificateNotRevokedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificateRevocation *diagnostic.CertificateRevocationWrapper, currentTime time.Time,
	constraint policy.LevelRule, subContext enumerations.SubContext) *CertificateNotRevokedCheck {
	c := &CertificateNotRevokedCheck{
		ChainItemBase:         process.NewChainItemBase(i18nProvider, result, constraint),
		certificateRevocation: certificateRevocation,
		currentTime:           currentTime,
		subContext:            subContext,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateNotRevokedCheck) Process() bool {
	isRevoked := c.certificateRevocation != nil && c.certificateRevocation.IsRevoked() &&
		enumerations.RevocationReason_CERTIFICATE_HOLD != c.certificateRevocation.Reason()
	if isRevoked {
		revocationDate := c.certificateRevocation.RevocationDate()
		isRevoked = revocationDate != nil && !c.currentTime.Before(*revocationDate)
	}
	return !isRevoked
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *CertificateNotRevokedCheck) BuildAdditionalInfo() *string {
	if c.certificateRevocation != nil && c.certificateRevocation.RevocationDate() != nil {
		revocationDateStr := process.GetFormattedDate(c.certificateRevocation.RevocationDate())
		// Java hands MessageFormat the RevocationReason ENUM, so a NULL reason
		// renders as the literal text "null"; the ported wrapper returns the
		// empty string for that null, which would render as nothing.
		reason := "null"
		if r := c.certificateRevocation.Reason(); r != "" {
			reason = string(r)
		}
		message := c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_REASON, reason, revocationDateStr)
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateNotRevokedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCR
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateNotRevokedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCR_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateNotRevokedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateNotRevokedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if enumerations.SubContext_SIGNING_CERT == c.subContext {
		return enumerations.SubIndication_REVOKED_NO_POE
	}
	return enumerations.SubIndication_REVOKED_CA_NO_POE
}

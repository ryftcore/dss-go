// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/BestSignatureTimeBeforeSuspensionTimeCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// BestSignatureTimeBeforeSuspensionTimeCheck checks if best-signature-time is
// before the suspension date (onHold).
type BestSignatureTimeBeforeSuspensionTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessLongTermData]

	// certificateRevocation is the certificate's revocation.
	certificateRevocation *diagnostic.CertificateRevocationWrapper

	// bestSignatureTime is the best signature time.
	bestSignatureTime *time.Time
}

// NewBestSignatureTimeBeforeSuspensionTimeCheck is the default constructor.
// Port of
// BestSignatureTimeBeforeSuspensionTimeCheck(I18nProvider, XmlValidationProcessLongTermData, CertificateRevocationWrapper, Date, LevelRule).
func NewBestSignatureTimeBeforeSuspensionTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData], certificateRevocation *diagnostic.CertificateRevocationWrapper,
	bestSignatureTime *time.Time, constraint policy.LevelRule) *BestSignatureTimeBeforeSuspensionTimeCheck {
	c := &BestSignatureTimeBeforeSuspensionTimeCheck{
		ChainItemBase:         process.NewChainItemBase(i18nProvider, result, constraint),
		certificateRevocation: certificateRevocation,
		bestSignatureTime:     bestSignatureTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) Process() bool {
	revocationDate := c.certificateRevocation.RevocationDate()
	return revocationDate != nil && c.bestSignatureTime != nil && c.bestSignatureTime.Before(*revocationDate)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) BuildAdditionalInfo() *string {
	bestSignatureTimeStr := " ? "
	if c.bestSignatureTime != nil {
		bestSignatureTimeStr = process.GetFormattedDate(c.bestSignatureTime)
	}
	revocationTime := " ? "
	if c.certificateRevocation.RevocationDate() != nil {
		revocationTime = process.GetFormattedDate(c.certificateRevocation.RevocationDate())
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_BEST_SIGNATURE_TIME_CERT_SUSPENSION, bestSignatureTimeStr, revocationTime)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ISTPTBST
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ISTPTBST_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BestSignatureTimeBeforeSuspensionTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_TRY_LATER
}

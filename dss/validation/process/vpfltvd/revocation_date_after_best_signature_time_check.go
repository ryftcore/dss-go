// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/RevocationDateAfterBestSignatureTimeCheck.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.vpfltvd.checks flattens into this
// single Go package vpfltvd (collision-checked with the vpfltvd root
// classes), following the same checks-subpackage flattening convention used
// throughout this port.
package vpfltvd

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationDateAfterBestSignatureTimeCheck checks if the revocation date is
// after best-signature-time.
type RevocationDateAfterBestSignatureTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessLongTermData]

	// certificateRevocation is the certificate's revocation.
	certificateRevocation *diagnostic.CertificateRevocationWrapper

	// bestSignatureTime is the best signature time.
	bestSignatureTime *time.Time

	// subContext is the validation SubContext.
	subContext enumerations.SubContext
}

// NewRevocationDateAfterBestSignatureTimeCheck is the default constructor.
// Port of
// RevocationDateAfterBestSignatureTimeCheck(I18nProvider, XmlValidationProcessLongTermData, CertificateRevocationWrapper, Date, LevelRule, SubContext).
func NewRevocationDateAfterBestSignatureTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData], certificateRevocation *diagnostic.CertificateRevocationWrapper,
	bestSignatureTime *time.Time, constraint policy.LevelRule, subContext enumerations.SubContext) *RevocationDateAfterBestSignatureTimeCheck {
	c := &RevocationDateAfterBestSignatureTimeCheck{
		ChainItemBase:         process.NewChainItemBase(i18nProvider, result, constraint),
		certificateRevocation: certificateRevocation,
		bestSignatureTime:     bestSignatureTime,
		subContext:            subContext,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationDateAfterBestSignatureTimeCheck) Process() bool {
	revocationDate := c.certificateRevocation.RevocationDate()
	// revocation date can be null in case of unknown status
	return revocationDate != nil && c.bestSignatureTime != nil && revocationDate.After(*c.bestSignatureTime)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationDateAfterBestSignatureTimeCheck) BuildAdditionalInfo() *string {
	bestSignatureTimeStr := " ? "
	if c.bestSignatureTime != nil {
		bestSignatureTimeStr = process.GetFormattedDate(c.bestSignatureTime)
	}
	revocationTime := " ? "
	if c.certificateRevocation.RevocationDate() != nil {
		revocationTime = process.GetFormattedDate(c.certificateRevocation.RevocationDate())
	}
	message := c.I18nProvider.GetMessage(c.getBestSignatureTimeRevocationCheckMessageTag(), bestSignatureTimeStr, revocationTime)
	return &message
}

// getBestSignatureTimeRevocationCheckMessageTag returns the MessageTag to be
// used to build the additional info. Port of
// getBestSignatureTimeRevocationCheckMessageTag().
func (c *RevocationDateAfterBestSignatureTimeCheck) getBestSignatureTimeRevocationCheckMessageTag() i18n.MessageTag {
	return i18n.MessageTagBESTSignatureTimeCertRevocation
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDateAfterBestSignatureTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTIRTPTBST
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationDateAfterBestSignatureTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTIRTPTBSTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationDateAfterBestSignatureTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationDateAfterBestSignatureTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if enumerations.SubContextSigningCert == c.subContext {
		return enumerations.SubIndicationRevokedNoPOE
	}
	return enumerations.SubIndicationRevokedCANoPOE
}

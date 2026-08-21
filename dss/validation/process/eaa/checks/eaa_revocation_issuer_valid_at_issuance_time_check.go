// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/status/EAARevocationIssuerValidAtIssuanceTimeCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationIssuerValidAtIssuanceTimeCheck checks whether the issuer
// certificate of the EAA revocation token was valid at the EAA revocation
// token issuance time.
type EAARevocationIssuerValidAtIssuanceTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaaStatusToken is the EAA revocation token to check.
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper
}

// NewEAARevocationIssuerValidAtIssuanceTimeCheck is the default constructor.
func NewEAARevocationIssuerValidAtIssuanceTimeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, constraint policy.LevelRule) *EAARevocationIssuerValidAtIssuanceTimeCheck {
	c := &EAARevocationIssuerValidAtIssuanceTimeCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaaStatusToken: eaaStatusToken,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) Process() bool {
	issuedAt := c.eaaStatusToken.IssuedAt()
	signingCertificate := c.eaaStatusToken.SigningCertificate()
	return issuedAt != nil && signingCertificate != nil &&
		!issuedAt.Before(*signingCertificate.NotBefore()) &&
		!issuedAt.After(*signingCertificate.NotAfter())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_ISS_VALID
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_ISS_VALID_ANS
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) BuildAdditionalInfo() *string {
	if c.eaaStatusToken.SigningCertificate() != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_EAA_REV_ISS_CERT,
			process.GetFormattedDate(c.eaaStatusToken.IssuedAt()),
			process.GetFormattedDate(c.eaaStatusToken.SigningCertificate().NotBefore()),
			process.GetFormattedDate(c.eaaStatusToken.SigningCertificate().NotAfter()))
		return &message
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationIssuerValidAtIssuanceTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}

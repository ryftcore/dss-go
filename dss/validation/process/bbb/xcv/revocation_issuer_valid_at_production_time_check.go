// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationIssuerValidAtProductionTimeCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationIssuerValidAtProductionTimeCheck verifies whether the revocation's
// issuer certificate has been valid at the revocation's production time.
type RevocationIssuerValidAtProductionTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationIssuerValidAtProductionTimeCheck is the default constructor. Port
// of RevocationIssuerValidAtProductionTimeCheck(I18nProvider, XmlRAC,
// RevocationWrapper, LevelRule).
func NewRevocationIssuerValidAtProductionTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlRAC], revocationData *diagnostic.RevocationWrapper,
	constraint policy.LevelRule) *RevocationIssuerValidAtProductionTimeCheck {
	c := &RevocationIssuerValidAtProductionTimeCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationIssuerValidAtProductionTimeCheck) Process() bool {
	// check performed only for OCSP certificates
	return enumerations.RevocationTypeOCSP != c.revocationData.RevocationType() ||
		c.checkOCSPResponderValidAtRevocationProductionTime()
}

// checkOCSPResponderValidAtRevocationProductionTime ports the private
// checkOCSPResponderValidAtRevocationProductionTime().
func (c *RevocationIssuerValidAtProductionTimeCheck) checkOCSPResponderValidAtRevocationProductionTime() bool {
	revocationIssuer := c.revocationData.SigningCertificate()
	producedAt := c.revocationData.ProductionDate()
	return revocationIssuer != nil &&
		!producedAt.Before(*revocationIssuer.NotBefore()) &&
		!producedAt.After(*revocationIssuer.NotAfter())
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationIssuerValidAtProductionTimeCheck) BuildAdditionalInfo() *string {
	if c.revocationData.SigningCertificate() != nil {
		var messageTag i18n.MessageTag
		if c.Process() {
			messageTag = i18n.MessageTagRevocationProducedAtCertValidity
		} else {
			messageTag = i18n.MessageTagRevocationProducedAtOutOfBounds
		}
		message := c.I18nProvider.GetMessage(messageTag,
			process.GetFormattedDate(c.revocationData.ProductionDate()),
			process.GetFormattedDate(c.revocationData.SigningCertificate().NotBefore()),
			process.GetFormattedDate(c.revocationData.SigningCertificate().NotAfter()))
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationIssuerValidAtProductionTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRevocIssuerValidAtProd
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationIssuerValidAtProductionTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVRevocIssuerValidAtProdANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationIssuerValidAtProductionTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationIssuerValidAtProductionTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}

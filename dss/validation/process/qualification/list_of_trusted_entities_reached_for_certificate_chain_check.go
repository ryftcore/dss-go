// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/ListOfTrustedEntitiesReachedForCertificateChainCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ListOfTrustedEntitiesReachedForCertificateChainCheck checks whether a
// List of Trusted Entities has been reached for the given certificate
// chain.
type ListOfTrustedEntitiesReachedForCertificateChainCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// signingCertificate is the end-entity certificate.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewListOfTrustedEntitiesReachedForCertificateChainCheck is the default
// constructor. Port of
// ListOfTrustedEntitiesReachedForCertificateChainCheck(Provider, XmlValidationPIDQualificationProcess, CertificateWrapper, LevelRule).
func NewListOfTrustedEntitiesReachedForCertificateChainCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlValidationPIDQualificationProcess], signingCertificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *ListOfTrustedEntitiesReachedForCertificateChainCheck {
	c := &ListOfTrustedEntitiesReachedForCertificateChainCheck{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ListOfTrustedEntitiesReachedForCertificateChainCheck) Process() bool {
	return c.signingCertificate != nil && c.signingCertificate.IsListOfTrustedEntitiesReached()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ListOfTrustedEntitiesReachedForCertificateChainCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAACertLoTEReached
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *ListOfTrustedEntitiesReachedForCertificateChainCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAACertLoTEReachedANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ListOfTrustedEntitiesReachedForCertificateChainCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *ListOfTrustedEntitiesReachedForCertificateChainCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

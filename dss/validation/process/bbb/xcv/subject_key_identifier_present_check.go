// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/SubjectKeyIdentifierPresentCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SubjectKeyIdentifierPresentCheck verifies whether the CA certificate
// contains the RFC 5280 "4.2.1.2. Subject Key Identifier" certificate
// extension.
type SubjectKeyIdentifierPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewSubjectKeyIdentifierPresentCheck is the default constructor. Port of
// SubjectKeyIdentifierPresentCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewSubjectKeyIdentifierPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *SubjectKeyIdentifierPresentCheck {
	c := &SubjectKeyIdentifierPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SubjectKeyIdentifierPresentCheck) Process() bool {
	return c.certificate.SubjectKeyIdentifier() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SubjectKeyIdentifierPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISKIP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *SubjectKeyIdentifierPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISKIP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SubjectKeyIdentifierPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SubjectKeyIdentifierPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}

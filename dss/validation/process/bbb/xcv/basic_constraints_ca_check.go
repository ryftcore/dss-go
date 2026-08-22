// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/BasicConstraintsCACheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicConstraintsCACheck verifies if the certificate contains
// BasicConstraint.cA attribute and its value is set to true.
type BasicConstraintsCACheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewBasicConstraintsCACheck is the default constructor. Port of
// BasicConstraintsCACheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewBasicConstraintsCACheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *BasicConstraintsCACheck {
	c := &BasicConstraintsCACheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BasicConstraintsCACheck) Process() bool {
	return c.certificate.IsCA()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BasicConstraintsCACheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICAC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *BasicConstraintsCACheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICAC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BasicConstraintsCACheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BasicConstraintsCACheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/GivenNameCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// GivenNameCheck checks if the certificate's given name are acceptable.
type GivenNameCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewGivenNameCheck is the default constructor. Port of
// GivenNameCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewGivenNameCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *GivenNameCheck {
	c := &GivenNameCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *GivenNameCheck) Process() bool {
	return c.ProcessValueCheck(c.certificate.GivenName())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *GivenNameCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGGIVEN
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *GivenNameCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGGIVEN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *GivenNameCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *GivenNameCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}

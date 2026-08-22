// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/TitleCheck.java (DSS 6.5.RC1).
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

// TitleCheck checks if the certificate's title attribute is acceptable.
type TitleCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewTitleCheck is the default constructor. Port of
// TitleCheck(Provider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewTitleCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *TitleCheck {
	c := &TitleCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TitleCheck) Process() bool {
	return c.ProcessValueCheck(c.certificate.Title())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TitleCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCGTITLE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TitleCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCGTITLEANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TitleCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *TitleCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcCCLegislationCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// CertificateQcCCLegislationCheck checks if the country code or set of
// country codes defined in QcCClegislation is supported by the policy.
type CertificateQcCCLegislationCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// constraint is the constraint.
	constraint policy.MultiValuesRule
}

// NewCertificateQcCCLegislationCheck is the default constructor. Port of
// CertificateQcCCLegislationCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateQcCCLegislationCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateQcCCLegislationCheck {
	c := &CertificateQcCCLegislationCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
		constraint:                   constraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcCCLegislationCheck) Process() bool {
	return c.ProcessValuesCheck(c.certificate.QcLegislationCountryCodes())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcCCLegislationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCDCQCCLCEC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
//
// See EN 319 412-5 ch. 4.2.1: a certificate that includes the
// esi4-qcStatement-1 statement with the aim to declare that it is an EU
// qualified certificate that is issued according to Directive 1999/93/EC
// [i.3] or the Annex I, III or IV of the Regulation (EU) No 910/2014 [i.8]
// whichever is in force at the time of issuance shall not include the
// QcCClegislation statement.
func (c *CertificateQcCCLegislationCheck) ErrorMessageTag() i18n.MessageTag {
	if utils.IsCollectionEmpty(c.constraint.Values()) {
		return i18n.MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS_EU
	}
	return i18n.MessageTag_BBB_XCV_CMDCDCQCCLCEC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcCCLegislationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcCCLegislationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateSemanticsIdentifierCheck.java (DSS 6.5.RC1).
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

// CertificateSemanticsIdentifierCheck checks the QCStatement
// SemanticsIdentifier value.
type CertificateSemanticsIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateSemanticsIdentifierCheck is the default constructor. Port of
// CertificateSemanticsIdentifierCheck(Provider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateSemanticsIdentifierCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateSemanticsIdentifierCheck {
	c := &CertificateSemanticsIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateSemanticsIdentifierCheck) Process() bool {
	var values []string
	semanticsIdentifier := c.certificate.SemanticsIdentifier()
	if semanticsIdentifier != "" {
		values = append(values, semanticsIdentifier.Name(), semanticsIdentifier.OID(), semanticsIdentifier.Description())
	}
	return c.ProcessValuesCheck(values)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateSemanticsIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCSCSIA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateSemanticsIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCSCSIAANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateSemanticsIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateSemanticsIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/CertifiedRolesCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// CertifiedRolesCheck checks if the certified roles are acceptable.
type CertifiedRolesCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewCertifiedRolesCheck is the default constructor. Port of
// CertifiedRolesCheck(I18nProvider, XmlSAV, SignatureWrapper, MultiValuesRule).
func NewCertifiedRolesCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.MultiValuesRule) *CertifiedRolesCheck {
	c := &CertifiedRolesCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		signature:                    signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertifiedRolesCheck) Process() bool {
	certifiedRoles := c.signature.SignerRoleDetails(c.signature.CertifiedRoles())
	return c.ProcessValuesCheck(certifiedRoles)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertifiedRolesCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICERRM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CertifiedRolesCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_SAV_ICERRM_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertifiedRolesCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *CertifiedRolesCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}

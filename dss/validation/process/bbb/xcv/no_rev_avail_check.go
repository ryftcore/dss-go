// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/NoRevAvailCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// NoRevAvailCheck verifies certificate's compliance to the RFC 9608.
type NoRevAvailCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewNoRevAvailCheck is the default constructor. Port of
// NoRevAvailCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewNoRevAvailCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *NoRevAvailCheck {
	c := &NoRevAvailCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// RFC 9608 "3. Other X.509 Certificate Extensions".
//
// If the noRevAvail extension is present in a certificate, then:
//
//   - The certificate MUST NOT also include the basic constraints
//     certificate extension with the cA BOOLEAN set to TRUE; see
//     Section 4.2.1.9 of [RFC5280].
//
//   - The certificate MUST NOT also include the CRL Distribution Points
//     certificate extension; see Section 4.2.1.13 of [RFC5280].
//
//   - The certificate MUST NOT also include the Freshest CRL certificate
//     extension; see Section 4.2.1.15 of [RFC5280].
//
//   - The Authority Information Access certificate extension, if
//     present, MUST NOT include an id-ad-ocsp accessMethod; see
//     Section 4.2.2.1 of [RFC5280].
func (c *NoRevAvailCheck) Process() bool {
	if c.certificate.IsNoRevAvail() {
		return !c.certificate.IsCA() && utils.IsCollectionEmpty(c.certificate.CRLDistributionPoints()) &&
			utils.IsCollectionEmpty(c.certificate.FreshestCRLUrls()) && utils.IsCollectionEmpty(c.certificate.OCSPAccessUrls())
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *NoRevAvailCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICNRAEV
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *NoRevAvailCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ICNRAEV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *NoRevAvailCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *NoRevAvailCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}

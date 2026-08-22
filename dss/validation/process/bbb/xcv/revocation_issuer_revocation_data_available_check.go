// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationIssuerRevocationDataAvailableCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationIssuerRevocationDataAvailableCheck checks if the revocation data is
// available for the revocation issuer's certificate.
type RevocationIssuerRevocationDataAvailableCheck struct {
	*RevocationDataAvailableCheck[*jaxb.XmlRAC]
}

// NewRevocationIssuerRevocationDataAvailableCheck is the default constructor.
// Port of RevocationIssuerRevocationDataAvailableCheck(I18nProvider, XmlRAC,
// CertificateWrapper, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' getMessageTag/getErrorMessageTag rather
// than the ones inherited from RevocationDataAvailableCheck.
func NewRevocationIssuerRevocationDataAvailableCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlRAC], certificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *RevocationIssuerRevocationDataAvailableCheck {
	c := &RevocationIssuerRevocationDataAvailableCheck{
		RevocationDataAvailableCheck: NewRevocationDataAvailableCheck(i18nProvider, result, certificate, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden
// getMessageTag().
func (c *RevocationIssuerRevocationDataAvailableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIRDPFRC
}

// ErrorMessageTag returns the check's error message tag. Port of the overridden
// getErrorMessageTag().
func (c *RevocationIssuerRevocationDataAvailableCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIRDPFRCANS
}

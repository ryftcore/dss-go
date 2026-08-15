// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/CertificateKnownToBeNotRevokedEnforceFailCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// CertificateKnownToBeNotRevokedEnforceFailCheck returns information whether
// the signing-certificate is known to not be revoked and revocation data is
// acceptable, but always returns the Basic Signature Validation conclusion.
type CertificateKnownToBeNotRevokedEnforceFailCheck struct {
	*CertificateKnownToBeNotRevokedCheck[*jaxb.XmlValidationProcessLongTermData]
}

// NewCertificateKnownToBeNotRevokedEnforceFailCheck is the default
// constructor. Port of
// CertificateKnownToBeNotRevokedEnforceFailCheck(I18nProvider, XmlValidationProcessLongTermData, CertificateWrapper, CertificateRevocationWrapper, boolean, Date, XmlConclusion, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' Process rather than the one inherited
// from CertificateKnownToBeNotRevokedCheck.
func NewCertificateKnownToBeNotRevokedEnforceFailCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData], certificate *diagnostic.CertificateWrapper,
	revocationData *diagnostic.CertificateRevocationWrapper, isRevocationDataIssuerTrusted bool, currentTime *time.Time,
	bsConclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *CertificateKnownToBeNotRevokedEnforceFailCheck {
	c := &CertificateKnownToBeNotRevokedEnforceFailCheck{
		CertificateKnownToBeNotRevokedCheck: NewCertificateKnownToBeNotRevokedCheck(i18nProvider, result, certificate,
			revocationData, isRevocationDataIssuerTrusted, currentTime, bsConclusion, constraint),
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateKnownToBeNotRevokedEnforceFailCheck) Process() bool {
	return false
}

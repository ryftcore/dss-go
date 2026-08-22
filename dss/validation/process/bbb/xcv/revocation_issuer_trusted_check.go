// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/RevocationIssuerTrustedCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationIssuerTrustedCheck checks if the provided certificate token is
// trusted.
type RevocationIssuerTrustedCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// currentTime is the validation time.
	currentTime time.Time

	// revocationIssuerSunsetDateConstraint identifies the constraint for the
	// revocation data issuer sunset date.
	revocationIssuerSunsetDateConstraint policy.LevelRule
}

// NewRevocationIssuerTrustedCheck is the default constructor. Port of
// RevocationIssuerTrustedCheck(Provider, T, CertificateWrapper, Date, LevelRule, LevelRule).
func NewRevocationIssuerTrustedCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, revocationIssuerSunsetDateConstraint policy.LevelRule,
	constraint policy.LevelRule) *RevocationIssuerTrustedCheck[T] {
	c := &RevocationIssuerTrustedCheck[T]{
		ChainItemBase:                        process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:                          certificate,
		currentTime:                          currentTime,
		revocationIssuerSunsetDateConstraint: revocationIssuerSunsetDateConstraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationIssuerTrustedCheck[T]) Process() bool {
	return c.certificate != nil && process.IsTrustAnchor(c.certificate, c.currentTime, c.revocationIssuerSunsetDateConstraint)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationIssuerTrustedCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVICRDIT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationIssuerTrustedCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVICRDITANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion(), whose Java body returns null.
func (c *RevocationIssuerTrustedCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return ""
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose Java body returns null.
func (c *RevocationIssuerTrustedCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationIssuerTrustedCheck[T]) BuildAdditionalInfo() *string {
	if c.certificate != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagCertificateID, c.certificate.Id())
		return &message
	}
	return nil
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/CertificateKnownToBeNotRevokedCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateKnownToBeNotRevokedCheck verifies whether the signing-certificate
// is known to not be revoked and revocation data is acceptable.
type CertificateKnownToBeNotRevokedCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to be verified.
	certificate *diagnostic.CertificateWrapper

	// revocationData is the revocation data to be verified.
	revocationData *diagnostic.CertificateRevocationWrapper

	// isRevocationDataIssuerTrusted defines whether a revocation data issuer is
	// trusted.
	isRevocationDataIssuerTrusted bool

	// currentTime is the validation time.
	currentTime *time.Time

	// bsConclusion is the conclusion of the Basic Signature Validation block.
	bsConclusion *jaxb.XmlConclusion
}

// NewCertificateKnownToBeNotRevokedCheck is the default constructor. Port of
// CertificateKnownToBeNotRevokedCheck(I18nProvider, T, CertificateWrapper, CertificateRevocationWrapper, boolean, Date, XmlConclusion, LevelRule).
func NewCertificateKnownToBeNotRevokedCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, revocationData *diagnostic.CertificateRevocationWrapper,
	isRevocationDataIssuerTrusted bool, currentTime *time.Time, bsConclusion *jaxb.XmlConclusion,
	constraint policy.LevelRule) *CertificateKnownToBeNotRevokedCheck[T] {
	c := &CertificateKnownToBeNotRevokedCheck[T]{
		ChainItemBase:                 process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:                   certificate,
		revocationData:                revocationData,
		isRevocationDataIssuerTrusted: isRevocationDataIssuerTrusted,
		currentTime:                   currentTime,
		bsConclusion:                  bsConclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateKnownToBeNotRevokedCheck[T]) Process() bool {
	return c.isKnownToBeNotRevoked()
}

// isKnownToBeNotRevoked ports the private isKnownToBeNotRevoked().
func (c *CertificateKnownToBeNotRevokedCheck[T]) isKnownToBeNotRevoked() bool {
	return c.revocationData != nil && !c.revocationData.IsRevoked() && c.revocationData.SigningCertificate() != nil &&
		c.isRevocationIssuerValid(c.revocationData.SigningCertificate())
}

// isInValidityRange ports the private isInValidityRange(CertificateWrapper).
func (c *CertificateKnownToBeNotRevokedCheck[T]) isInValidityRange(certificateWrapper *diagnostic.CertificateWrapper) bool {
	notBefore := certificateWrapper.NotBefore()
	notAfter := certificateWrapper.NotAfter()
	if c.currentTime == nil {
		return false
	}
	return notBefore != nil && !c.currentTime.Before(*notBefore) && notAfter != nil && !c.currentTime.After(*notAfter)
}

// isRevocationIssuerValid ports the private isRevocationIssuerValid(CertificateWrapper).
func (c *CertificateKnownToBeNotRevokedCheck[T]) isRevocationIssuerValid(revocationDataIssuer *diagnostic.CertificateWrapper) bool {
	return c.isRevocationDataIssuerTrusted || c.isInValidityRange(revocationDataIssuer)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CertificateKnownToBeNotRevokedCheck[T]) BuildAdditionalInfo() *string {
	if c.revocationData != nil && c.revocationData.SigningCertificate() != nil && !c.isKnownToBeNotRevoked() {
		revocationIssuer := c.revocationData.SigningCertificate()
		notBeforeStr := " ? "
		if revocationIssuer.NotBefore() != nil {
			notBeforeStr = process.GetFormattedDate(revocationIssuer.NotBefore())
		}
		notAfterStr := " ? "
		if revocationIssuer.NotAfter() != nil {
			notAfterStr = process.GetFormattedDate(revocationIssuer.NotAfter())
		}
		validationTime := process.GetFormattedDate(c.currentTime)
		message := c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CERT_VALIDITY,
			revocationIssuer.Id(), c.revocationData.Id(), notBeforeStr, notAfterStr, validationTime)
		return &message
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TOKEN_ID, c.certificate.Id())
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateKnownToBeNotRevokedCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.bsConclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateKnownToBeNotRevokedCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.bsConclusion.SubIndication == nil {
		return ""
	}
	return c.bsConclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateKnownToBeNotRevokedCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_LTV_ISCKNR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CertificateKnownToBeNotRevokedCheck[T]) ErrorMessageTag() i18n.MessageTag {
	if !c.isKnownToBeNotRevoked() {
		return i18n.MessageTag_LTV_ISCKNR_ANS1
	}
	return i18n.MessageTag_LTV_ISCKNR_ANS0
}

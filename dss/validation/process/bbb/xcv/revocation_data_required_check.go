// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/RevocationDataRequiredCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// RevocationDataRequiredCheck is used to verify whether the revocation data
// check shall be skipped for the given certificate.
type RevocationDataRequiredCheck[T any] struct {
	*bbb.AbstractCertificateCheckItem[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// currentTime is the validation time.
	currentTime time.Time

	// certificateSunsetDateConstraint is the certificate's sunset date constraint.
	certificateSunsetDateConstraint policy.LevelRule
}

// NewRevocationDataRequiredCheck is the default constructor. Port of
// RevocationDataRequiredCheck(I18nProvider, T, CertificateWrapper, Date, LevelRule, CertificateApplicabilityRule).
func NewRevocationDataRequiredCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, certificateSunsetDateConstraint policy.LevelRule,
	constraint policy.CertificateApplicabilityRule) *RevocationDataRequiredCheck[T] {
	c := &RevocationDataRequiredCheck[T]{
		AbstractCertificateCheckItem:    bbb.NewAbstractCertificateCheckItemWithId(i18nProvider, result, certificate, constraint),
		certificate:                     certificate,
		currentTime:                     currentTime,
		certificateSunsetDateConstraint: certificateSunsetDateConstraint,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *RevocationDataRequiredCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeRACSubXCV
}

// Process performs the check. Port of process().
func (c *RevocationDataRequiredCheck[T]) Process() bool {
	return !c.isTrustAnchor(c.certificate) && !c.certificate.IsSelfSigned() && !c.ProcessCertificateCheck(c.certificate)
}

// isTrustAnchor ports the private isTrustAnchor(CertificateWrapper).
func (c *RevocationDataRequiredCheck[T]) isTrustAnchor(certificate *diagnostic.CertificateWrapper) bool {
	return process.IsTrustAnchor(certificate, c.currentTime, c.certificateSunsetDateConstraint)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDataRequiredCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IRDCSFC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationDataRequiredCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IRDCSFC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationDataRequiredCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationDataRequiredCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/IdentificationOfSigningCertificateResultCheck.java (DSS 6.5.RC1).
package vpfbs

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// IdentificationOfSigningCertificateResultCheck verifies if the identification
// of the signing certificate (as per clause 5.2.3) succeeded.
type IdentificationOfSigningCertificateResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlISC is the Identification of the Signing Certificate building block
	// result.
	xmlISC *jaxb.XmlISC
}

// NewIdentificationOfSigningCertificateResultCheck is the default constructor.
// Port of
// IdentificationOfSigningCertificateResultCheck(I18nProvider, T, XmlISC, TokenProxy, LevelRule).
func NewIdentificationOfSigningCertificateResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlISC *jaxb.XmlISC, token diagnostic.TokenProxy, constraint policy.LevelRule) *IdentificationOfSigningCertificateResultCheck[T] {
	c := &IdentificationOfSigningCertificateResultCheck[T]{
		// Identification of Signing Certificate building block suffix ("-ISC"),
		// a per-class private constant in Java; inlined here since Go
		// package-level constants share one namespace.
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+"-ISC"),
		xmlISC:        xmlISC,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *IdentificationOfSigningCertificateResultCheck[T]) Process() bool {
	return c.xmlISC != nil && c.IsValid(&c.xmlISC.XmlConstraintsConclusionContent)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *IdentificationOfSigningCertificateResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *IdentificationOfSigningCertificateResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_SIGNING_CERTIFICATE_FOUND
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *IdentificationOfSigningCertificateResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IISCRC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *IdentificationOfSigningCertificateResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BSV_IISCRC_ANS
}

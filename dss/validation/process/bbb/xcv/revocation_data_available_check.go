// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/RevocationDataAvailableCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// RevocationDataAvailableCheck checks if a revocation data is available for
// the certificate.
type RevocationDataAvailableCheck[T any] struct {
	*process.ChainItemBase[T]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewRevocationDataAvailableCheck is the default constructor. Port of
// RevocationDataAvailableCheck(I18nProvider, T, CertificateWrapper, LevelRule).
func NewRevocationDataAvailableCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *RevocationDataAvailableCheck[T] {
	c := &RevocationDataAvailableCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// NewRevocationDataAvailableCheckWithId is the constructor with token id.
// Port of RevocationDataAvailableCheck(I18nProvider, T, CertificateWrapper, LevelRule, String).
func NewRevocationDataAvailableCheckWithId[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule, tokenId string) *RevocationDataAvailableCheck[T] {
	c := &RevocationDataAvailableCheck[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, tokenId),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *RevocationDataAvailableCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_LTV_SUB_XCV
}

// Process performs the check. Port of process().
func (c *RevocationDataAvailableCheck[T]) Process() bool {
	return c.certificate.IsRevocationDataAvailable()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationDataAvailableCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IRDPFC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationDataAvailableCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IRDPFC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationDataAvailableCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationDataAvailableCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}

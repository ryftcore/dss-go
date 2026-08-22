// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateRevocationSelectorResultCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateRevocationSelectorResultCheck verifies the result of a
// CertificateRevocationSelector.
type CertificateRevocationSelectorResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// crsResult is the CRS result.
	crsResult *jaxb.XmlCRS
}

// NewCertificateRevocationSelectorResultCheck is the default constructor.
// Port of CertificateRevocationSelectorResultCheck(I18nProvider, T, XmlCRS, LevelRule).
func NewCertificateRevocationSelectorResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	crsResult *jaxb.XmlCRS, constraint policy.LevelRule) *CertificateRevocationSelectorResultCheck[T] {
	var chainItemBase *process.ChainItemBase[T]
	if crsResult.Id != nil {
		chainItemBase = process.NewChainItemBaseWithId(i18nProvider, result, constraint, *crsResult.Id)
	} else {
		chainItemBase = process.NewChainItemBase(i18nProvider, result, constraint)
	}
	c := &CertificateRevocationSelectorResultCheck[T]{
		ChainItemBase: chainItemBase,
		crsResult:     crsResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *CertificateRevocationSelectorResultCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeCRS
}

// Process performs the check. Port of process().
func (c *CertificateRevocationSelectorResultCheck[T]) Process() bool {
	return c.IsValid(&c.crsResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateRevocationSelectorResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIARDPFC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateRevocationSelectorResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVIARDPFCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateRevocationSelectorResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.crsResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateRevocationSelectorResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.crsResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.crsResult.Conclusion.SubIndication.SubIndication()
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of getPreviousErrors().
func (c *CertificateRevocationSelectorResultCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	return c.crsResult.Conclusion.Errors
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CertificateRevocationSelectorResultCheck[T]) BuildAdditionalInfo() *string {
	if c.crsResult.LatestAcceptableRevocationId != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagLastAcceptableRevocation, *c.crsResult.LatestAcceptableRevocationId)
		return &message
	}
	return nil
}

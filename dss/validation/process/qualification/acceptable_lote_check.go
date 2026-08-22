// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/AcceptableLoTECheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableLoTECheck verifies if the validation of a LoTE was successful.
// AcceptableLoLoTECheck embeds this type and re-registers itself as the
// overrides target (the same InitChainItem re-registration pattern
// bbb/xcv.RevocationDataFreshCheck over AbstractRevocationFreshCheck uses)
// to change only MessageTag()/ErrorMessageTag().
type AcceptableLoTECheck[T any] struct {
	*process.ChainItemBase[T]

	// loteAnalysis is the LoTE validation result.
	loteAnalysis *jaxb.XmlLoTEAnalysis
}

// NewAcceptableLoTECheck is the default constructor. Port of
// AcceptableLoTECheck(I18nProvider, T, XmlLoTEAnalysis, LevelRule).
func NewAcceptableLoTECheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	loteAnalysis *jaxb.XmlLoTEAnalysis, constraint policy.LevelRule) *AcceptableLoTECheck[T] {
	c := &AcceptableLoTECheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		loteAnalysis:  loteAnalysis,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableLoTECheck[T]) Process() bool {
	return c.IsValidConclusion(c.loteAnalysis.Conclusion)
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *AcceptableLoTECheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_LIST_OF_TRUSTED_ENTITIES, c.loteAnalysis.URL)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableLoTECheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LOTE_ACCEPT
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AcceptableLoTECheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_LOTE_ACCEPT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableLoTECheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *AcceptableLoTECheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

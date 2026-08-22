// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/AbstractTrustedListCheck.java (DSS 6.5.RC1).
//
// AcceptableTrustedListCheck and AcceptableListOfTrustedListsCheck embed this
// type and re-register themselves as the overrides target (the same
// InitChainItem re-registration pattern AcceptableLoLoTECheck over
// AcceptableLoTECheck uses) to define only MessageTag()/ErrorMessageTag().
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractTrustedListCheck is the abstract class verifying the validity of
// the Trusted List.
type AbstractTrustedListCheck[T any] struct {
	*process.ChainItemBase[T]

	// tlAnalysis is the Trusted List validation result.
	tlAnalysis *jaxb.XmlTLAnalysis
}

// NewAbstractTrustedListCheck is the default constructor. Port of the
// protected AbstractTrustedListCheck(I18nProvider, T, XmlTLAnalysis, LevelRule).
func NewAbstractTrustedListCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	tlAnalysis *jaxb.XmlTLAnalysis, constraint policy.LevelRule) *AbstractTrustedListCheck[T] {
	c := &AbstractTrustedListCheck[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, tlAnalysis.Id),
		tlAnalysis:    tlAnalysis,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AbstractTrustedListCheck[T]) Process() bool {
	return c.IsValidConclusion(c.tlAnalysis.Conclusion)
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *AbstractTrustedListCheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagTrustedList, c.tlAnalysis.URL)
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AbstractTrustedListCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *AbstractTrustedListCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}

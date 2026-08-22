// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/cv/checks/SignatureIntactWithIdCheck.java (DSS 6.5.RC1).
package cv

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignatureIntactWithIdCheck checks if the signature is intact for the given
// token, with a difference that provides the token's Id to the additional
// information.
type SignatureIntactWithIdCheck[T any] struct {
	*SignatureIntactCheck[T]
}

// NewSignatureIntactWithIdCheck is the default constructor. Port of
// SignatureIntactWithIdCheck(I18nProvider, T, TokenProxy, Context, LevelRule).
//
// The constructor re-registers the overrides with the outer type, so that the
// base's self-calls reach this class' buildAdditionalInfo rather than the one
// inherited from SignatureIntactCheck.
func NewSignatureIntactWithIdCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	token diagnostic.TokenProxy, context enumerations.Context,
	constraint policy.LevelRule) *SignatureIntactWithIdCheck[T] {
	c := &SignatureIntactWithIdCheck[T]{
		SignatureIntactCheck: NewSignatureIntactCheck(i18nProvider, result, token, context, constraint),
	}
	c.InitChainItem(c)
	return c
}

// BuildAdditionalInfo builds an additional information. Port of the overridden
// buildAdditionalInfo().
func (c *SignatureIntactWithIdCheck[T]) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TOKEN_ID, c.Token.Id())
	return &message
}

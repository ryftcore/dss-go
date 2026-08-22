//go:build eaa

// Ported from dss-validation/.../validation/process/bbb/fc/EAARevocationFormatChecking.java (DSS 6.5.RC1).
//
// Gated behind the "eaa" build tag - see eaa_format_checking.go's header in
// this same package for the rationale.
package fc

import (
	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/eaa/checks"
)

// EAARevocationFormatChecking verifies the format of an EAA revocation token.
type EAARevocationFormatChecking struct {
	AbstractFormatChecking[*diagnostic.EAARevocationTokenWrapper]
}

// NewEAARevocationFormatChecking is the default constructor.
func NewEAARevocationFormatChecking(i18nProvider *i18n.Provider, diagnosticData *diagnostic.Data,
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, context enumerations.Context,
	pol policy.ValidationPolicy) *EAARevocationFormatChecking {
	c := &EAARevocationFormatChecking{}
	c.InitAbstractFormatChecking(i18nProvider, diagnosticData, eaaStatusToken, context, pol)
	c.InitChainBase(c)
	return c
}

// InitChain builds the constraint chain. Port of the overridden protected void initChain().
//
// TODO: JWT/CWT format checks are not implemented (upstream Java DSS doesn't implement them either).
func (c *EAARevocationFormatChecking) InitChain() {
	c.FirstItem = c.typeCheck()
}

func (c *EAARevocationFormatChecking) typeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAARevocationTokenTypeConstraint()
	return checks.NewEAARevocationTokenTypeCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

//go:build eaa

// Ported from dss-validation/.../validation/process/bbb/fc/EAARevocationFormatChecking.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): gated behind the "eaa" build tag - see
// eaa_format_checking.go's header in this same package for the rationale.
//
// UNPORTED DEPENDENCY (flagged per porter brief - do not invent): the check
// this class wires (EAARevocationTokenTypeCheck) lives in the Java package
// eu.europa.esig.dss.validation.process.eaa.status, which is NOT part of the
// phase 8c package layout and was not present anywhere in the repository at
// port time. This file assumes it lands in a sibling package
// "github.com/utain/esig/dss/validation/process/eaa" (status flattened in,
// mirroring how fc/checks flattens in here) with a constructor
// NewEAARevocationTokenTypeCheck(i18nProvider, result *process.Result[*drjaxb.XmlFC],
// token *diagnostic.EAARevocationTokenWrapper, constraint policy.MultiValuesRule)
// returning process.ChainItem[*drjaxb.XmlFC]. See porter notes: this file does not
// build until that package exists.
package fc

import (
	drjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	policy "github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/eaa"
)

// EAARevocationFormatChecking verifies the format of an EAA revocation token.
type EAARevocationFormatChecking struct {
	AbstractFormatChecking[*diagnostic.EAARevocationTokenWrapper]
}

// NewEAARevocationFormatChecking is the default constructor.
func NewEAARevocationFormatChecking(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	eaaStatusToken *diagnostic.EAARevocationTokenWrapper, context enumerations.Context,
	pol policy.ValidationPolicy) *EAARevocationFormatChecking {
	c := &EAARevocationFormatChecking{}
	c.InitAbstractFormatChecking(i18nProvider, diagnosticData, eaaStatusToken, context, pol)
	c.InitChainBase(c)
	return c
}

// InitChain builds the constraint chain. Port of the overridden protected void initChain().
//
// TODO (ported as-is from Java): JWT/CWT format checks are not yet implemented upstream either.
func (c *EAARevocationFormatChecking) InitChain() {
	c.FirstItem = c.typeCheck()
}

func (c *EAARevocationFormatChecking) typeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAARevocationTokenTypeConstraint()
	return eaa.NewEAARevocationTokenTypeCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

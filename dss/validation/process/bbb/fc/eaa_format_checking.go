//go:build eaa

// Ported from dss-validation/.../validation/process/bbb/fc/EAAFormatChecking.java (DSS 6.5.RC1).
//
// Gated behind the "eaa" build tag: eu.europa.esig.dss.validation.process.eaa
// is deferred (see docs/compatibility/known-gaps.md, "EAA"), so it gets a
// feature tag of its own. Only this file and
// eaa_revocation_format_checking.go in this package need the tag - every
// other fc file is untagged and builds today. blocks/basic_building_blocks_eaa.go
// (the only caller of NewEAAFormatChecking) carries the same eaa tag, so
// no untagged code path references this file.
//
// The eaa-tagged import lands in
// github.com/ryftcore/dss-go/dss/validation/process/eaa/checks - a dedicated
// package split out from the eaa root package (which holds
// EAAValidationBlock/EAAValidationProcess and imports qualification), so that
// this file's import doesn't close a
// bbb/fc -> eaa -> qualification -> vpfswatsp -> bbb/sav -> eaa cycle. See
// bbb/sav/eaa_acceptance_validation.go for the full rationale.
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

// EAAFormatChecking verifies the format of an Electronic Attestation of Attributes (EAA).
type EAAFormatChecking struct {
	AbstractFormatChecking[*diagnostic.EAAWrapper]
}

// NewEAAFormatChecking is the default constructor.
func NewEAAFormatChecking(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	eaaToken *diagnostic.EAAWrapper, context enumerations.Context, pol policy.ValidationPolicy) *EAAFormatChecking {
	c := &EAAFormatChecking{}
	c.InitAbstractFormatChecking(i18nProvider, diagnosticData, eaaToken, context, pol)
	c.InitChainBase(c)
	return c
}

// InitChain builds the constraint chain. Port of the overridden protected void initChain().
func (c *EAAFormatChecking) InitChain() {
	var item process.ChainItem[*drjaxb.XmlFC] = c.signatureUnicity()
	c.FirstItem = item

	item = item.SetNextItem(c.disclosurePresent())
	item = item.SetNextItem(c.disclosureListExhaustive())
	item = item.SetNextItem(c.keyBindingSignaturePresent())
}

func (c *EAAFormatChecking) signatureUnicity() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAASignatureUnicityConstraint()
	return checks.NewEAASignatureUnicityCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) disclosurePresent() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAADisclosurePresentConstraint()
	return checks.NewDisclosurePresentCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) disclosureListExhaustive() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAADisclosureListExhaustiveConstraint()
	return checks.NewDisclosureListExhaustiveCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) keyBindingSignaturePresent() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAAKeyBindingSignaturePresentConstraint()
	return checks.NewKeyBindingSignaturePresentCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

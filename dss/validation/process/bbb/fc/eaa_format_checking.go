//go:build eaa

// Ported from dss-validation/.../validation/process/bbb/fc/EAAFormatChecking.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): gated behind the "eaa" build tag, the
// same technique already used for phase8/phase8d/phase8e forward
// dependencies elsewhere in this tree. eu.europa.esig.dss.validation.process.eaa
// is explicitly listed as deferred past phase 9 in PORTING_PLAN.md ("EAA/mdoc
// modules ... revisit after Phase 9"), not merely a later numbered phase, so
// it gets its own feature tag rather than a phaseNN one. Only this file and
// eaa_revocation_format_checking.go in this package need the tag - every
// other fc file is untagged and builds today. blocks/basic_building_blocks_eaa.go
// (the only caller of NewEAAFormatChecking) carries the same eaa tag, so
// no untagged code path references this file.
//
// UNPORTED DEPENDENCY (flagged per porter brief - do not invent): the four
// checks this class wires (EAASignatureUnicityCheck, DisclosurePresentCheck,
// DisclosureListExhaustiveCheck, KeyBindingSignaturePresentCheck) live in the
// Java package eu.europa.esig.dss.validation.process.eaa.checks, which is
// NOT part of the phase 8c package layout (bbb/{isc,vci,cv,fc,sav} only) and
// was not present anywhere in the repository at port time. This file assumes
// they will land in a sibling package "github.com/utain/esig/dss/validation/process/eaa"
// with constructors NewEAASignatureUnicityCheck/NewDisclosurePresentCheck/
// NewDisclosureListExhaustiveCheck/NewKeyBindingSignaturePresentCheck taking
// (i18nProvider, result *process.Result[*drjaxb.XmlFC], token *diagnostic.EAAWrapper, constraint policy.LevelRule)
// and returning process.ChainItem[*drjaxb.XmlFC], mirroring every other check
// constructor in this package. See porter notes: this file does not build
// until that package exists.
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
	return eaa.NewEAASignatureUnicityCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) disclosurePresent() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAADisclosurePresentConstraint()
	return eaa.NewDisclosurePresentCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) disclosureListExhaustive() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAADisclosureListExhaustiveConstraint()
	return eaa.NewDisclosureListExhaustiveCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

func (c *EAAFormatChecking) keyBindingSignaturePresent() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.EAAKeyBindingSignaturePresentConstraint()
	return eaa.NewKeyBindingSignaturePresentCheck(c.I18nProvider, c.Result, c.Token, constraint)
}

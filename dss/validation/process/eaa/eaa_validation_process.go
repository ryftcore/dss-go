// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/EAAValidationProcess.java (DSS 6.5.RC1).
//
// ValidationProcess extends vpfbs.AbstractBasicValidationProcess but fully
// replaces InitChain() rather than building on the base's own 5.3 logic (the
// EAA presentation process per EN 319 102-1 Annex is a different sequence:
// format checking, per-signature and key-binding-signature validation
// conclusiveness, digest/selective-disclosure verification, then SAV - not
// ISC/VCI/XCV/CV). Its own InitChain, defined directly on *ValidationProcess,
// shadows the promoted one from AbstractBasicValidationProcess per ordinary Go
// method resolution - the same technique CertificateRevocationSelector
// documents for its own subclassing.
package eaa

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/qualification"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfbs"
)

// ValidationProcess performs validation of a presentation of Electronic
// Attestation of Attributes.
type ValidationProcess struct {
	*vpfbs.AbstractBasicValidationProcess[*jaxb.XmlValidationProcessEAA]

	// eaa is the EAA being validated.
	eaa *diagnostic.EAAWrapper

	// xmlSignatures is the map of validated signatures.
	xmlSignatures map[string]*jaxb.XmlSignature

	// policy is the validation policy used to validate evidence records.
	policy policy.ValidationPolicy
}

// NewEAAValidationProcess is the common constructor. Port of
// ValidationProcess(Provider, EAAWrapper, Map, Map, ValidationPolicy).
func NewValidationProcess(i18nProvider *i18n.Provider, eaaWrapper *diagnostic.EAAWrapper,
	xmlSignatures map[string]*jaxb.XmlSignature, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	validationPolicy policy.ValidationPolicy) *ValidationProcess {
	xmlResult := &jaxb.XmlValidationProcessEAA{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent,
		&xmlResult.XmlConstraintsConclusionAttrs)
	c := &ValidationProcess{
		// diagnosticData is nil: Java passes null, since EAAValidationProcess's
		// own InitChain never reads it (unlike AbstractBasicValidationProcess's).
		AbstractBasicValidationProcess: vpfbs.NewAbstractBasicValidationProcess[*jaxb.XmlValidationProcessEAA](i18nProvider, result, nil, eaaWrapper, bbbs),
		eaa:                            eaaWrapper,
		xmlSignatures:                  xmlSignatures,
		policy:                         validationPolicy,
	}
	c.InitAbstractBasicValidationProcess(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *ValidationProcess) Title() i18n.MessageTag {
	return i18n.MessageTagVPEAA
}

// InitChain initializes the chain, fully replacing
// AbstractBasicValidationProcess's own 5.3 process. Port of the overridden
// initChain().
func (c *ValidationProcess) InitChain() {

	eaaBBBs := c.BBBs[c.eaa.Id()]
	if eaaBBBs == nil {
		panic(fmt.Sprintf("Missing Basic Building Blocks result for token with Id '%s'", c.eaa.Id()))
	}

	var item process.ChainItem[*jaxb.XmlValidationProcessEAA]

	// 1. Format checking: calls the inherited formatChecking(XmlFC) directly by
	// constructing vpfbs.FormatCheckingResultCheck - AbstractBasicValidationProcess's
	// own formatChecking helper is unexported and, unlike Java's protected,
	// cannot be called across the vpfbs/eaa package boundary; see file header.
	xmlFC := eaaBBBs.FC
	if xmlFC != nil {
		item = vpfbs.NewFormatCheckingResultCheck(c.I18nProvider, c.Result, xmlFC, c.Token, c.FailLevelRule())
		c.FirstItem = item
	}

	// 2. Verify electronic signatures
	for _, signatureWrapper := range c.eaa.EAASignatures() {
		item = c.appendItem(item, c.signatureValidationConclusive(signatureWrapper))
	}

	// 3. Verify Key Binding signature
	if c.eaa.KeyBindingSignature() != nil {
		item = c.appendItem(item, c.keyBindingSignatureValidationConclusive(c.eaa.KeyBindingSignature()))
	}

	// 4. Digest (selective disclosures) validation
	xmlCV := eaaBBBs.CV
	if xmlCV != nil {
		item = c.appendItem(item, vpfbs.NewCryptographicVerificationResultCheck(c.I18nProvider, c.Result, xmlCV, c.eaa, c.FailLevelRule()))
	}

	// 5. EAA Acceptance Validation
	xmlSAV := eaaBBBs.SAV
	if xmlSAV != nil {
		item = c.appendItem(item, c.signatureAcceptanceValidation(xmlSAV)) //nolint:staticcheck // mirrors upstream EAAValidationProcess#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
}

// appendItem ports `item.setNextItem(next)`. Java seeds `item` with firstItem, which only the
// format-checking step (1.) ever sets: when the basic building blocks carry no FC block, the
// first following step calls setNextItem on null and throws a NullPointerException, so an EAA
// validation without a format-checking result cannot proceed. The Go form of that null
// dereference is the bare nil-interface panic this helper replaces with an explicit
// IllegalStateException-style panic of the same fail-closed effect (never a verdict built
// without the format-checking gate).
//
// The FC block is missing in exactly one situation: a build without the `eaa` build tag, where
// blocks.BasicBuildingBlocks has no EAA format-checking block to dispatch to (see
// blocks/basic_building_blocks_noeaa.go) although the diagnostic data and this presentation
// process are untagged and reach here for any diagnostic data holding an EAA. EAA validation
// is only supported with `-tags eaa`.
func (c *ValidationProcess) appendItem(item, next process.ChainItem[*jaxb.XmlValidationProcessEAA]) process.ChainItem[*jaxb.XmlValidationProcessEAA] {
	if item == nil {
		panic(fmt.Sprintf("Invalid state! No Format Checking result found for the EAA with Id '%s': "+
			"EAA validation is only supported when built with the 'eaa' build tag", c.eaa.Id()))
	}
	return item.SetNextItem(next)
}

// signatureValidationConclusive ports the private
// signatureValidationConclusive(SignatureWrapper).
func (c *ValidationProcess) signatureValidationConclusive(signatureWrapper *diagnostic.SignatureWrapper) process.ChainItem[*jaxb.XmlValidationProcessEAA] {
	constraint := c.policy.EAASignatureValidConstraint()
	return qualification.NewSignatureValidationResultCheck(c.I18nProvider, c.Result, c.getSignatureBasicProcessConclusion(signatureWrapper), constraint)
}

// getSignatureBasicProcessConclusion ports the private
// getSignatureBasicProcessConclusion(SignatureWrapper).
func (c *ValidationProcess) getSignatureBasicProcessConclusion(signatureWrapper *diagnostic.SignatureWrapper) *jaxb.XmlConclusion {
	xmlSignature := c.xmlSignatures[signatureWrapper.Id()]
	if xmlSignature == nil {
		panic(fmt.Sprintf("Invalid state! No basic signature validation process found for the signature with Id '%s'!", signatureWrapper.Id()))
	}
	return xmlSignature.ValidationProcessBasicSignature.Conclusion
}

// keyBindingSignatureValidationConclusive ports the private
// keyBindingSignatureValidationConclusive(SignatureWrapper).
func (c *ValidationProcess) keyBindingSignatureValidationConclusive(signatureWrapper *diagnostic.SignatureWrapper) process.ChainItem[*jaxb.XmlValidationProcessEAA] {
	constraint := c.policy.EAAKeyBindingSignatureValidConstraint()
	xmlSignature := c.xmlSignatures[signatureWrapper.Id()]
	return NewKeyBindingSignatureValidationResultCheck(c.I18nProvider, c.Result, xmlSignature.ValidationProcessBasicSignature.Conclusion, constraint)
}

// signatureAcceptanceValidation executes "5.2.8 Signature Acceptance
// Validation (SAV)" building block, with EAA-specific message tags. Port of
// the overridden signatureAcceptanceValidation(XmlSAV).
func (c *ValidationProcess) signatureAcceptanceValidation(xmlSAV *jaxb.XmlSAV) process.ChainItem[*jaxb.XmlValidationProcessEAA] {
	return newEAASignatureAcceptanceValidationResultCheck(c.I18nProvider, c.Result, xmlSAV, c.eaa, c.FailLevelRule())
}

// eaaSignatureAcceptanceValidationResultCheck is the Go form of the anonymous
// SignatureAcceptanceValidationResultCheck subclass returned by
// signatureAcceptanceValidation(XmlSAV): the same check reporting EAA-specific
// message tags.
type eaaSignatureAcceptanceValidationResultCheck struct {
	*vpfbs.SignatureAcceptanceValidationResultCheck[*jaxb.XmlValidationProcessEAA]
}

func newEAASignatureAcceptanceValidationResultCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlValidationProcessEAA], xmlSAV *jaxb.XmlSAV, token diagnostic.TokenProxy,
	constraint policy.LevelRule) *eaaSignatureAcceptanceValidationResultCheck {
	c := &eaaSignatureAcceptanceValidationResultCheck{
		SignatureAcceptanceValidationResultCheck: vpfbs.NewSignatureAcceptanceValidationResultCheck(i18nProvider, result, xmlSAV, token, constraint),
	}
	c.InitChainItem(c)
	return c
}

func (c *eaaSignatureAcceptanceValidationResultCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIEAAAVRC
}

func (c *eaaSignatureAcceptanceValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBSVIEAAAVRCANS
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden collectAdditionalMessages(XmlConclusion).
func (c *ValidationProcess) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	tokenBBBs := c.BBBs[c.eaa.Id()]
	if tokenBBBs != nil {
		conclusion.Errors = c.clear(conclusion.Errors)
		conclusion.Errors = append(conclusion.Errors, tokenBBBs.Conclusion.Errors...)
		conclusion.Warnings = c.clear(conclusion.Warnings)
		conclusion.Warnings = append(conclusion.Warnings, tokenBBBs.Conclusion.Warnings...)
		conclusion.Infos = c.clear(conclusion.Infos)
		conclusion.Infos = append(conclusion.Infos, tokenBBBs.Conclusion.Infos...)

		for _, constraint := range c.Result.Constraint() {
			c.CollectMessages(conclusion, constraint)
		}
	}
}

// clear ports the private clear(List<XmlMessage>): keeps only the signature
// validation ("ADEST_IBSVPSC_ANS") and key-binding ("EAA_KBRC_ANS") messages
// out of the given list, discarding the rest.
func (c *ValidationProcess) clear(messageList []*jaxb.XmlMessage) []*jaxb.XmlMessage {
	if utils.IsCollectionEmpty(messageList) {
		return nil
	}
	adestIbsvpscAnsKey := i18n.MessageTagADESTIBSVPSCANS.Id()
	eaaKbrcAnsKey := i18n.MessageTagEAAKBRCANS.Id()
	kept := make([]*jaxb.XmlMessage, 0, len(messageList))
	for _, m := range messageList {
		if m.Key != nil && (adestIbsvpscAnsKey == *m.Key || eaaKbrcAnsKey == *m.Key) {
			kept = append(kept, m)
		}
	}
	return kept
}

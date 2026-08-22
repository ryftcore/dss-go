// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/LongTermValidationCertificateRevocationSelector.java (DSS 6.5.RC1).
//
// Extends bbb/xcv's CertificateRevocationSelector, overriding
// verifyRevocationData, getRevocationAcceptanceValidationResult,
// acceptableRevocationDataAvailable and collectMessages - the same
// embed-and-re-register technique the base package documents for its own
// subclassing (see xcv.CertificateRevocationSelectorOverrides).
package vpfltvd

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
)

// LongTermValidationCertificateRevocationSelectorOverrides declares the
// overridable protected methods this class adds on top of the ones
// CertificateRevocationSelector already routes. vpfswatsp's two subclasses
// (PastSignatureValidationCertificateRevocationSelector, which reads the
// revocation BBB out of the bbbs map instead of running
// RevocationBasicValidationProcess, and ValidationTimeSlidingCertificateRevocationSelector)
// register through InitLongTermValidationCertificateRevocationSelector, so the
// self-calls below reach the subclass the way Java's virtual dispatch does.
type LongTermValidationCertificateRevocationSelectorOverrides interface {
	xcv.CertificateRevocationSelectorOverrides

	// RevocationBBBConclusion returns a conclusion of the revocation basic
	// building block execution process. Port of the protected
	// getRevocationBBBConclusion(CertificateRevocationWrapper).
	RevocationBBBConclusion(revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlConclusion
}

// LongTermValidationCertificateRevocationSelector verifies and returns the
// latest acceptable revocation data for a long-term validation process.
type LongTermValidationCertificateRevocationSelector struct {
	*xcv.CertificateRevocationSelector

	// ltvOverrides points back at the concrete selector; see
	// InitLongTermValidationCertificateRevocationSelector.
	ltvOverrides LongTermValidationCertificateRevocationSelectorOverrides

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.Data

	// BBBs is the map of BasicBuildingBlocks. Exported because Java declares
	// the field protected.
	BBBs map[string]*jaxb.XmlBasicBuildingBlocks

	// TokenId is the Id of a token being validated (e.g. signature id,
	// timestamp id). Exported because Java declares the field protected.
	TokenId string
}

// NewLongTermValidationCertificateRevocationSelector is the default
// constructor. Port of
// LongTermValidationCertificateRevocationSelector(Provider, CertificateWrapper, Date, Data, Map, String, ValidationPolicy).
func NewLongTermValidationCertificateRevocationSelector(i18nProvider *i18n.Provider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, diagnosticData *diagnostic.Data,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tokenId string,
	validationPolicy policy.ValidationPolicy) *LongTermValidationCertificateRevocationSelector {
	c := &LongTermValidationCertificateRevocationSelector{}
	c.InitLongTermValidationCertificateRevocationSelectorState(i18nProvider, certificate, currentTime,
		diagnosticData, bbbs, tokenId, validationPolicy)
	c.InitLongTermValidationCertificateRevocationSelector(c)
	return c
}

// InitLongTermValidationCertificateRevocationSelectorState wires the shared
// state, the way the Java constructor body does. A subclass calls it before
// InitLongTermValidationCertificateRevocationSelector, in place of the Java
// super(...) call.
func (c *LongTermValidationCertificateRevocationSelector) InitLongTermValidationCertificateRevocationSelectorState(
	i18nProvider *i18n.Provider, certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	diagnosticData *diagnostic.Data, bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tokenId string,
	validationPolicy policy.ValidationPolicy) {
	base := &xcv.CertificateRevocationSelector{}
	base.InitCertificateRevocationSelectorState(i18nProvider, certificate, currentTime, validationPolicy, make(map[string]struct{}))
	c.CertificateRevocationSelector = base
	c.diagnosticData = diagnosticData
	c.BBBs = bbbs
	c.TokenId = tokenId
}

// InitLongTermValidationCertificateRevocationSelector registers the concrete
// selector with its base so that the base can dispatch to the overridden
// methods, CertificateRevocationSelector's and Chain's included. It must be
// called exactly once, by the concrete selector's constructor, before Execute.
func (c *LongTermValidationCertificateRevocationSelector) InitLongTermValidationCertificateRevocationSelector(
	overrides LongTermValidationCertificateRevocationSelectorOverrides) {
	c.ltvOverrides = overrides
	c.InitCertificateRevocationSelector(overrides)
}

// ltvSelectorOverrides returns the registered overrides, panicking when the
// concrete selector forgot to call
// InitLongTermValidationCertificateRevocationSelector.
func (c *LongTermValidationCertificateRevocationSelector) ltvSelectorOverrides() LongTermValidationCertificateRevocationSelectorOverrides {
	if c.ltvOverrides == nil {
		panic("LongTermValidationCertificateRevocationSelector was not initialised: the concrete selector must call InitLongTermValidationCertificateRevocationSelector in its constructor")
	}
	return c.ltvOverrides
}

// NewLongTermValidationCertificateRevocationSelectorWithoutDiagnosticData is
// the protected constructor. Port of
// LongTermValidationCertificateRevocationSelector(Provider, CertificateWrapper, Date, Map, String, ValidationPolicy).
func NewLongTermValidationCertificateRevocationSelectorWithoutDiagnosticData(i18nProvider *i18n.Provider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	tokenId string, validationPolicy policy.ValidationPolicy) *LongTermValidationCertificateRevocationSelector {
	return NewLongTermValidationCertificateRevocationSelector(i18nProvider, certificate, currentTime, nil, bbbs, tokenId, validationPolicy)
}

// VerifyRevocationData verifies the given revocation data and returns the
// resulting ChainItem. Port of the overridden
// verifyRevocationData(ChainItem, CertificateRevocationWrapper).
func (c *LongTermValidationCertificateRevocationSelector) VerifyRevocationData(item process.ChainItem[*jaxb.XmlCRS],
	revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS] {
	revocationBBBConclusion := c.ltvSelectorOverrides().RevocationBBBConclusion(revocationWrapper)

	if revocationBBBConclusion != nil {
		if item == nil {
			item = c.revocationBasicValidationAcceptable(revocationWrapper.Id(), revocationBBBConclusion)
			c.FirstItem = item
		} else {
			item = item.SetNextItem(c.revocationBasicValidationAcceptable(revocationWrapper.Id(), revocationBBBConclusion))
		}
		if process.IsAllowedBasicRevocationDataValidation(revocationBBBConclusion) {
			item = c.CertificateRevocationSelector.VerifyRevocationData(item, revocationWrapper)
		}
	}

	allowedBBB := process.IsAllowedBasicRevocationDataValidation(revocationBBBConclusion)

	validity, ok := c.RevocationDataValidityMap[revocationWrapper.Id()]
	if !ok {
		validity = allowedBBB
	} else {
		validity = validity && allowedBBB
	}
	c.RevocationDataValidityMap[revocationWrapper.Id()] = validity

	return item
}

// RevocationBBBConclusion returns a conclusion of the revocation basic
// building block execution process. Port of
// getRevocationBBBConclusion(CertificateRevocationWrapper).
func (c *LongTermValidationCertificateRevocationSelector) RevocationBBBConclusion(
	revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlConclusion {
	rbvp := NewRevocationBasicValidationProcess(c.I18nProvider, c.diagnosticData, &revocationWrapper.RevocationWrapper, c.BBBs)
	revocationBasicValidationResult := rbvp.Execute()
	return revocationBasicValidationResult.Conclusion
}

// RevocationAcceptanceValidationResult returns a RevocationAcceptanceValidation
// result for the given revocation token. Port of the overridden
// getRevocationAcceptanceValidationResult(CertificateRevocationWrapper).
func (c *LongTermValidationCertificateRevocationSelector) RevocationAcceptanceValidationResult(
	revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlRAC {
	return c.getRevocationAcceptanceValidationResultById(revocationWrapper.Id())
}

// getRevocationAcceptanceValidationResultById ports the private
// getRevocationAcceptanceValidationResult(String).
func (c *LongTermValidationCertificateRevocationSelector) getRevocationAcceptanceValidationResultById(revocationId string) *jaxb.XmlRAC {
	tokenBBB := c.BBBs[c.TokenId]
	return process.GetRevocationAcceptanceCheckerResult(tokenBBB, c.Certificate.Id(), revocationId)
}

// revocationBasicValidationAcceptable ports the private
// revocationBasicValidationAcceptable(String, XmlConclusion).
func (c *LongTermValidationCertificateRevocationSelector) revocationBasicValidationAcceptable(revocationId string,
	revocationBBBConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlCRS] {
	return NewRevocationDataAcceptableCheck(c.I18nProvider, c.Result, revocationId, revocationBBBConclusion, c.WarnLevelRule())
}

// AcceptableRevocationDataAvailable checks whether the acceptable revocation
// data is available. Port of the overridden acceptableRevocationDataAvailable().
func (c *LongTermValidationCertificateRevocationSelector) AcceptableRevocationDataAvailable() process.ChainItem[*jaxb.XmlCRS] {
	var acceptableRevocationData *diagnostic.RevocationWrapper
	if latest := c.LatestAcceptableCertificateRevocation(); latest != nil {
		acceptableRevocationData = &latest.RevocationWrapper
	}
	return newLongTermAcceptableRevocationDataAvailableCheck(c.I18nProvider, c.Result, acceptableRevocationData, c.FailLevelRule(), c)
}

// longTermAcceptableRevocationDataAvailableCheck is the Go form of the
// anonymous AcceptableRevocationDataAvailableCheck subclass returned by
// acceptableRevocationDataAvailable(): the same check, but reporting
// INDETERMINATE/TRY_LATER instead of the base's failure indication when any
// revocation basic validation for the certificate concluded TRY_LATER.
type longTermAcceptableRevocationDataAvailableCheck struct {
	*xcv.AcceptableRevocationDataAvailableCheck[*jaxb.XmlCRS]

	selector *LongTermValidationCertificateRevocationSelector
}

// newLongTermAcceptableRevocationDataAvailableCheck builds the anonymous
// subclass and re-registers the overrides with the outer type, so that the
// base's self-calls reach the overridden method here rather than
// AcceptableRevocationDataAvailableCheck's.
func newLongTermAcceptableRevocationDataAvailableCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlCRS],
	acceptableRevocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule,
	selector *LongTermValidationCertificateRevocationSelector) *longTermAcceptableRevocationDataAvailableCheck {
	c := &longTermAcceptableRevocationDataAvailableCheck{
		AcceptableRevocationDataAvailableCheck: xcv.NewAcceptableRevocationDataAvailableCheck(i18nProvider, result, acceptableRevocationData, constraint),
		selector:                               selector,
	}
	c.InitChainItem(c)
	return c
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port
// of the overridden getFailedIndicationForConclusion().
func (c *longTermAcceptableRevocationDataAvailableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion().
func (c *longTermAcceptableRevocationDataAvailableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.selector.isTryLater() {
		return enumerations.SubIndicationTryLater
	}
	return c.AcceptableRevocationDataAvailableCheck.FailedSubIndicationForConclusion()
}

// isTryLater ports the private isTryLater().
func (c *LongTermValidationCertificateRevocationSelector) isTryLater() bool {
	for _, revocationWrapper := range c.CertificateRevocationData() {
		conclusion := c.ltvSelectorOverrides().RevocationBBBConclusion(revocationWrapper)
		if conclusion == nil {
			continue
		}
		var subIndication enumerations.SubIndication
		if conclusion.SubIndication != nil {
			subIndication = conclusion.SubIndication.SubIndication()
		}
		if enumerations.IndicationIndeterminate == conclusion.Indication.Indication() &&
			enumerations.SubIndicationTryLater == subIndication {
			return true
		}
	}
	return false
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint).
func (c *LongTermValidationCertificateRevocationSelector) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.BlockType != nil && jaxb.XmlBlockTypeRevBBB == *constraint.BlockType && !c.IsValid(&c.Result.Value.XmlConstraintsConclusionContent) {
		c.collectMessagesForBBB(conclusion, constraint)
	}
	if constraint.BlockType != nil && jaxb.XmlBlockTypeRAC == *constraint.BlockType && !c.IsValid(&c.Result.Value.XmlConstraintsConclusionContent) {
		if constraint.Id != nil {
			xmlRAC := c.getRevocationAcceptanceValidationResultById(*constraint.Id)
			if xmlRAC != nil {
				c.CollectAllMessages(conclusion, xmlRAC.Conclusion)
			}
		}
	}
	c.CertificateRevocationSelector.CollectMessages(conclusion, constraint)
}

// collectMessagesForBBB ports the private
// collectMessagesForBBB(XmlConclusion, XmlConstraint).
func (c *LongTermValidationCertificateRevocationSelector) collectMessagesForBBB(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	c.CertificateRevocationSelector.CollectMessages(conclusion, constraint)
	if constraint.Id != nil {
		xmlBasicBuildingBlocks := c.BBBs[*constraint.Id]
		if xmlBasicBuildingBlocks != nil {
			c.CollectAllMessages(conclusion, xmlBasicBuildingBlocks.Conclusion)
		}
	}
}

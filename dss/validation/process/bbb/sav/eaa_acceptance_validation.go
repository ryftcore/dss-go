//go:build eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/EAAAcceptanceValidation.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): gated behind the "eaa" build tag - the
// same technique used for fc/eaa_format_checking.go (see its header for the
// rationale: PORTING_PLAN.md lists eu.europa.esig.dss.validation.process.eaa
// as deferred past phase 9, not a numbered upcoming phase). This keeps the
// default, untagged `go build ./...`/`go vet ./...`/`go test ./...` green;
// AbstractAcceptanceValidation (this package) and every other sav file that
// does not touch EAA build and are exercised today.
//
// UNPORTED DEPENDENCY (flagged per porter brief - do not invent): every check
// this class wires (ETSI194721ConformanceCheck, EAATypeCheck,
// EAATypeIntegrityPresentCheck, EAANotBeforePresentCheck, EAAExpirationPresentCheck,
// EAANotExpiredCheck, EAAAdministrativeIssuanceDatePresentCheck,
// EAAAdministrativeExpirationDatePresentCheck, EAAAdministrativePeriodNotExpiredCheck,
// EAAIdentifierPresentCheck, EAAIssuanceDatePresentCheck, EAACategoryCheck,
// EAASubjectCheck, EAASubjectPseudonymCheck, EAAIssuingCountryCheck,
// EAAIssuingAuthorityCheck, EAAIssuingAuthorityRegistrationIdentifierCheck,
// EAARevocationPresentCheck, EAARevocationAvailableCheck, EAARevocationStatusKnownCheck,
// EAARevocationAcceptableCheck, AcceptableEAARevocationFoundCheck, EAANotRevokedCheck,
// EAANotOnHoldCheck, EAAShortLivedCheck, EAAOneTimeUseCheck, EAAPseudonymUsageCheck,
// EAAClaimsCheck, EAASupportedClaimsCheck) lives in the Java package
// eu.europa.esig.dss.validation.process.eaa.checks, which is NOT part of the
// phase 8c package layout (bbb/{isc,vci,cv,fc,sav} only) and was not present
// anywhere in the repository at port time. Per the fc chunk's identical forward
// dependency on this same Java package (see fc/eaa_format_checking.go), this
// file assumes it lands in a sibling package
// "github.com/utain/esig/dss/validation/process/eaa" (both eaa.checks and
// eaa.status flatten into this one package, matching the checks-subpackage
// flattening convention used throughout this port, and matching the import
// path fc already committed to). Constructors are assumed to mirror their Java
// signature 1:1 (i18nProvider, result *process.Result[*jaxb.XmlSAV], plus the
// Java constructor's remaining arguments in order, constraint last) and return
// process.ChainItem[*jaxb.XmlSAV] (or, for the two exceptions read locally
// below - statusPresent()/EAARevocationPresentCheck - the concrete *eaa.EAARevocationPresentCheck
// type, since its Process() result is consulted before deciding the next
// chain item, mirroring how Java keeps the concrete type in a local variable).
// See porter notes: this file does not build until that package exists.
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/eaa"
)

// EAAAcceptanceValidation performs verification of an EAA against the
// validationPolicy defined acceptance criteria.
type EAAAcceptanceValidation struct {
	*AbstractAcceptanceValidation[*diagnostic.EAAWrapper]

	// bbbs is a map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// lastAcceptableStatus is the last acceptable EAA token status.
	lastAcceptableStatus *diagnostic.EAARevocationWrapper
}

// NewEAAAcceptanceValidation is the default constructor. Port of
// EAAAcceptanceValidation(I18nProvider, Date, EAAWrapper, Map, XmlAOV, ValidationPolicy).
func NewEAAAcceptanceValidation(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	eaaWrapper *diagnostic.EAAWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *EAAAcceptanceValidation {
	c := &EAAAcceptanceValidation{
		AbstractAcceptanceValidation: NewAbstractAcceptanceValidation(i18nProvider, eaaWrapper, currentTime,
			enumerations.Context_EAA, aovResult, validationPolicy),
		bbbs: bbbs,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *EAAAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *EAAAcceptanceValidation) InitChain() {
	item := c.etsi194721Conformance()
	c.FirstItem = item

	item = item.SetNextItem(c.eaaType())
	if enumerations.EAAType_SD_JWT_VC == c.token.EAAType() {
		item = item.SetNextItem(c.typeIntegrityPresent())
	}

	if enumerations.EAAType_ISO_IEC_MDOC == c.token.EAAType() {
		item = item.SetNextItem(c.issuanceDatePresent())
	}

	item = item.SetNextItem(c.eaaIdentifierPresent())

	item = item.SetNextItem(c.notBeforePresent())

	item = item.SetNextItem(c.expirationPresent())

	if c.token.EAANotBefore() != nil && c.token.EAAExpiration() != nil {
		item = item.SetNextItem(c.notExpired())
	}

	item = item.SetNextItem(c.administrativeIssuanceDatePresent())

	item = item.SetNextItem(c.administrativeExpirationDatePresent())

	if c.token.AdministrativeIssuanceDate() != nil && c.token.AdministrativeExpirationDate() != nil {
		item = item.SetNextItem(c.administrativePeriodNotExpired())
	}

	item = item.SetNextItem(c.category())

	item = item.SetNextItem(c.subject())

	item = item.SetNextItem(c.subjectPseudonym())

	item = item.SetNextItem(c.issuingCountry())

	item = item.SetNextItem(c.issuingAuthority())

	item = item.SetNextItem(c.issuingAuthorityRegistrationIdentifier())

	if utils.IsTrue(c.token.OneTimeUse()) {
		item = item.SetNextItem(c.oneTimeUse())
	}

	if utils.IsTrue(c.token.ShortLived()) {

		item = item.SetNextItem(c.shortLived())

	} else {

		// TODO : make status check configurable ?

		eaaRevocationPresentCheck := c.statusPresent()

		item = item.SetNextItem(eaaRevocationPresentCheck)

		if eaaRevocationPresentCheck.Process() {

			item = item.SetNextItem(c.statusAvailable())

			// TODO : improve with EAA Status selector ?
			c.lastAcceptableStatus = nil
			for _, eaaRevocationWrapper := range c.token.EAARevocations() {

				eaaRevocationBBB := c.bbbs[eaaRevocationWrapper.Id()]
				if eaaRevocationBBB == nil {
					panic("No BasicBuildingBlock found for token with Id '" + eaaRevocationWrapper.Id() + "'")
				}

				item = item.SetNextItem(c.statusKnown(eaaRevocationWrapper))

				item = item.SetNextItem(c.statusAcceptable(eaaRevocationWrapper, eaaRevocationBBB.Conclusion))

				if c.IsValidConclusion(eaaRevocationBBB.Conclusion) &&
					(c.lastAcceptableStatus == nil || c.lastAcceptableStatus.IssuedAt().Before(*eaaRevocationWrapper.IssuedAt())) {
					c.lastAcceptableStatus = eaaRevocationWrapper
				}

			}

			item = item.SetNextItem(c.acceptableStatusFound(c.lastAcceptableStatus))

			if c.lastAcceptableStatus != nil {

				item = item.SetNextItem(c.notRevoked(c.lastAcceptableStatus))

				item = item.SetNextItem(c.notOnHold(c.lastAcceptableStatus))

			}

		}

	}

	if c.token.HolderPseudonym() != "" {
		item = item.SetNextItem(c.usePseudonym())
	}

	item = item.SetNextItem(c.claims())

	item = item.SetNextItem(c.supportedClaims())

	// cryptographic check
	item = c.cryptographic(item)

}

func (c *EAAAcceptanceValidation) etsi194721Conformance() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAETSI194721ConformanceConstraint()
	return eaa.NewETSI194721ConformanceCheck(c.I18nProvider, c.Result, c.token, c.currentTime, constraint)
}

func (c *EAAAcceptanceValidation) eaaType() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAATypeConstraint()
	return eaa.NewEAATypeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) typeIntegrityPresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAATypeIntegrityPresentConstraint()
	return eaa.NewEAATypeIntegrityPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) notBeforePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAANotBeforePresentConstraint()
	return eaa.NewEAANotBeforePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) expirationPresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAExpirationPresentConstraint()
	return eaa.NewEAAExpirationPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) notExpired() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAANotExpiredConstraint()
	return eaa.NewEAANotExpiredCheck(c.I18nProvider, c.Result, c.token, c.currentTime, constraint)
}

func (c *EAAAcceptanceValidation) administrativeIssuanceDatePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAAdministrativeIssuanceDatePresentConstraint()
	return eaa.NewEAAAdministrativeIssuanceDatePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) administrativeExpirationDatePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAAdministrativeExpirationDatePresentConstraint()
	return eaa.NewEAAAdministrativeExpirationDatePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) administrativePeriodNotExpired() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAAdministrativePeriodNotExpiredConstraint()
	return eaa.NewEAAAdministrativePeriodNotExpiredCheck(c.I18nProvider, c.Result, c.token, c.currentTime, constraint)
}

func (c *EAAAcceptanceValidation) eaaIdentifierPresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAIdentifierPresentConstraint()
	return eaa.NewEAAIdentifierPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) issuanceDatePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAIssuanceDatePresentConstraint()
	return eaa.NewEAAIssuanceDatePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) category() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAACategoryConstraint()
	return eaa.NewEAACategoryCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) subject() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAASubjectConstraint()
	return eaa.NewEAASubjectCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) subjectPseudonym() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAASubjectPseudonymConstraint()
	return eaa.NewEAASubjectPseudonymCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) issuingCountry() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAIssuingCountryConstraint()
	return eaa.NewEAAIssuingCountryCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) issuingAuthority() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAIssuingAuthorityConstraint()
	return eaa.NewEAAIssuingAuthorityCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) issuingAuthorityRegistrationIdentifier() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAIssuingAuthorityRegistrationIdentifierConstraint()
	return eaa.NewEAAIssuingAuthorityRegistrationIdentifierCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) statusPresent() *eaa.EAARevocationPresentCheck {
	constraint := c.validationPolicy.EAARevocationPresentConstraint()
	return eaa.NewEAARevocationPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) statusAvailable() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationAvailableConstraint()
	return eaa.NewEAARevocationAvailableCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) statusKnown(eaaRevocationWrapper *diagnostic.EAARevocationWrapper) process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationUnknownStatusConstraint()
	return eaa.NewEAARevocationStatusKnownCheck(c.I18nProvider, c.Result, eaaRevocationWrapper, constraint)
}

func (c *EAAAcceptanceValidation) statusAcceptable(eaaRevocationWrapper *diagnostic.EAARevocationWrapper,
	xmlConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlSAV] {
	return eaa.NewEAARevocationAcceptableCheck(c.I18nProvider, c.Result, eaaRevocationWrapper, xmlConclusion, c.WarnLevelRule())
}

func (c *EAAAcceptanceValidation) acceptableStatusFound(acceptableEAARevocationWrapper *diagnostic.EAARevocationWrapper) process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationAvailableConstraint()
	return eaa.NewAcceptableEAARevocationFoundCheck(c.I18nProvider, c.Result, acceptableEAARevocationWrapper, constraint)
}

func (c *EAAAcceptanceValidation) notRevoked(eaaRevocationWrapper *diagnostic.EAARevocationWrapper) process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationNotRevokedConstraint()
	return eaa.NewEAANotRevokedCheck(c.I18nProvider, c.Result, eaaRevocationWrapper, constraint)
}

func (c *EAAAcceptanceValidation) notOnHold(eaaRevocationWrapper *diagnostic.EAARevocationWrapper) process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationNotOnHoldConstraint()
	return eaa.NewEAANotOnHoldCheck(c.I18nProvider, c.Result, eaaRevocationWrapper, constraint)
}

func (c *EAAAcceptanceValidation) shortLived() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAShortLivedConstraint()
	return eaa.NewEAAShortLivedCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) oneTimeUse() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAOneTimeUseConstraint()
	return eaa.NewEAAOneTimeUseCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) usePseudonym() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAUsePseudonymConstraint()
	return eaa.NewEAAPseudonymUsageCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) claims() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAAClaimsConstraint()
	return eaa.NewEAAClaimsCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAAAcceptanceValidation) supportedClaims() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAASupportedClaimsConstraint()
	return eaa.NewEAASupportedClaimsCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of collectMessages(XmlConclusion, XmlConstraint).
func (c *EAAAcceptanceValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.Name == nil || constraint.Name.Key == nil || *constraint.Name.Key != i18n.MessageTag_EAA_REV_ACC.Id() {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of collectAdditionalMessages(XmlConclusion).
func (c *EAAAcceptanceValidation) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	c.ChainBase.CollectAdditionalMessages(conclusion)

	if c.lastAcceptableStatus != nil {
		tokenBBB := c.bbbs[c.lastAcceptableStatus.Id()]
		c.CollectAllMessages(conclusion, tokenBBB.Conclusion)
	} else {
		for _, eaaRevocationWrapper := range c.token.EAARevocations() {
			tokenBBB := c.bbbs[eaaRevocationWrapper.Id()]
			c.CollectAllMessages(conclusion, tokenBBB.Conclusion)
		}
	}
}

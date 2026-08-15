//go:build eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/EAARevocationTokenAcceptanceValidation.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): gated behind the "eaa" build tag - see
// eaa_acceptance_validation.go's header in this same package for the
// rationale.
//
// UNPORTED DEPENDENCY (flagged per porter brief - do not invent): every check
// this class wires (EAARevocationIssuanceTimeCheck, EAARevocationExpirationTimeCheck,
// EAARevocationNotExpiredCheck, EAARevocationSubjectCheck, EAARevocationSubjectMatchCheck,
// EAARevocationIssuerValidAtIssuanceTimeCheck) lives in the Java package
// eu.europa.esig.dss.validation.process.eaa.status, which is NOT part of the
// phase 8c package layout (bbb/{isc,vci,cv,fc,sav} only) and was not present
// anywhere in the repository at port time. Per eaa_acceptance_validation.go
// (this same sav chunk's identical forward dependency on the sibling
// eaa.checks package), this file assumes eaa.status flattens into the same
// forward-declared Go package "github.com/utain/esig/dss/validation/process/eaa"
// (no name collisions with the eaa.checks types documented there).
// Constructors are assumed to mirror their Java signature 1:1 (i18nProvider,
// result *process.Result[*jaxb.XmlSAV], plus the Java constructor's remaining
// arguments in order, constraint last) and return process.ChainItem[*jaxb.XmlSAV].
// See porter notes: this file does not build until that package exists.
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/eaa"
)

// EAARevocationTokenAcceptanceValidation performs verification of an EAA
// revocation token against the validationPolicy defined acceptance criteria.
type EAARevocationTokenAcceptanceValidation struct {
	*AbstractAcceptanceValidation[*diagnostic.EAARevocationTokenWrapper]
}

// NewEAARevocationTokenAcceptanceValidation is the default constructor. Port of
// EAARevocationTokenAcceptanceValidation(I18nProvider, Date, EAARevocationTokenWrapper, XmlAOV, ValidationPolicy).
func NewEAARevocationTokenAcceptanceValidation(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	eaaRevocationTokenWrapper *diagnostic.EAARevocationTokenWrapper, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *EAARevocationTokenAcceptanceValidation {
	c := &EAARevocationTokenAcceptanceValidation{
		AbstractAcceptanceValidation: NewAbstractAcceptanceValidation(i18nProvider, eaaRevocationTokenWrapper,
			currentTime, enumerations.Context_EAA_REVOCATION, aovResult, validationPolicy),
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *EAARevocationTokenAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *EAARevocationTokenAcceptanceValidation) InitChain() {

	item := c.issuanceTime()
	c.FirstItem = item

	item = item.SetNextItem(c.expirationTime())

	if c.token.ExpirationTime() != nil {
		item = item.SetNextItem(c.notExpired())
	}

	item = item.SetNextItem(c.subject())

	if c.token.Subject() != "" {
		item = item.SetNextItem(c.subjectMatches())
	}

	if c.token.IssuedAt() != nil {
		item = item.SetNextItem(c.issuerValidAtIssuanceTime())
	}

	item = c.cryptographic(item)

}

func (c *EAARevocationTokenAcceptanceValidation) issuanceTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationIssuanceTimeConstraint()
	return eaa.NewEAARevocationIssuanceTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) expirationTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationExpirationTimeConstraint()
	return eaa.NewEAARevocationExpirationTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) notExpired() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationNotExpiredConstraint()
	return eaa.NewEAARevocationNotExpiredCheck(c.I18nProvider, c.Result, c.token, c.currentTime, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) subject() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationSubjectConstraint()
	return eaa.NewEAARevocationSubjectCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) subjectMatches() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationSubjectMatchConstraint()
	return eaa.NewEAARevocationSubjectMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) issuerValidAtIssuanceTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationIssuerValidAtIssuanceTimeConstraint()
	return eaa.NewEAARevocationIssuerValidAtIssuanceTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

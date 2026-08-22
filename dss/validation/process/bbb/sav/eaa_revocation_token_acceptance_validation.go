//go:build eaa

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/EAARevocationTokenAcceptanceValidation.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): gated behind the "eaa" build tag - see
// eaa_acceptance_validation.go's header in this same package for the
// rationale.
//
// INTEGRATION UPDATE (phase 8e integration pass): see
// eaa_acceptance_validation.go's header in this same package - the checks
// this file wires live in github.com/ryftcore/dss-go/dss/validation/process/eaa/checks
// (confirmed matching), a dedicated package rather than the eaa root package,
// specifically to keep this file's import from re-closing the
// bbb/sav -> eaa -> qualification -> vpfswatsp -> bbb/sav cycle.
package sav

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/eaa/checks"
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
			currentTime, enumerations.ContextEAARevocation, aovResult, validationPolicy),
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *EAARevocationTokenAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTagSignatureAcceptanceValidation
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
	return checks.NewEAARevocationIssuanceTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) expirationTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationExpirationTimeConstraint()
	return checks.NewEAARevocationExpirationTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) notExpired() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationNotExpiredConstraint()
	return checks.NewEAARevocationNotExpiredCheck(c.I18nProvider, c.Result, c.token, c.currentTime, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) subject() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationSubjectConstraint()
	return checks.NewEAARevocationSubjectCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) subjectMatches() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationSubjectMatchConstraint()
	return checks.NewEAARevocationSubjectMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *EAARevocationTokenAcceptanceValidation) issuerValidAtIssuanceTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.EAARevocationIssuerValidAtIssuanceTimeConstraint()
	return checks.NewEAARevocationIssuerValidAtIssuanceTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

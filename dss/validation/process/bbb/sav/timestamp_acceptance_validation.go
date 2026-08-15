// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/TimestampAcceptanceValidation.java (DSS 6.5.RC1).
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TimestampAcceptanceValidation is 5.2.8 Signature acceptance validation (SAV)
// for a timestamp token. This building block covers any additional
// verification to be performed on the signature itself or on the attributes of
// the signature ETSI EN 319 132-1.
type TimestampAcceptanceValidation struct {
	*AbstractAcceptanceValidation[*diagnostic.TimestampWrapper]
}

// NewTimestampAcceptanceValidation is the default constructor. Port of
// TimestampAcceptanceValidation(I18nProvider, Date, TimestampWrapper, XmlAOV, ValidationPolicy).
func NewTimestampAcceptanceValidation(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	timestamp *diagnostic.TimestampWrapper, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *TimestampAcceptanceValidation {
	c := &TimestampAcceptanceValidation{
		AbstractAcceptanceValidation: NewAbstractAcceptanceValidation(i18nProvider, timestamp, currentTime,
			enumerations.Context_TIMESTAMP, aovResult, validationPolicy),
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *TimestampAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *TimestampAcceptanceValidation) InitChain() {

	item := c.signingCertificateAttributePresent()
	c.FirstItem = item

	// See {@code SignatureAcceptanceValidation.initChain()}
	if c.token.IsSigningCertificateReferencePresent() {

		item = item.SetNextItem(c.unicitySigningCertificateAttribute())

		item = item.SetNextItem(c.signingCertificateReferencesValidity())

		item = item.SetNextItem(c.allCertificatesInPathReferenced())

	}

	item = item.SetNextItem(c.tsaGeneralNamePresent())

	if c.token.IsTSAGeneralNamePresent() {

		item = item.SetNextItem(c.tsaGeneralNameMatch())

		item = item.SetNextItem(c.tsaGeneralNameOrderMatch())

	}

	item = c.cryptographic(item)
}

// tsaGeneralNamePresent ports the private tsaGeneralNamePresent().
func (c *TimestampAcceptanceValidation) tsaGeneralNamePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.TimestampTSAGeneralNamePresent()
	return NewTSAGeneralNameFieldPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// tsaGeneralNameMatch ports the private tsaGeneralNameMatch().
func (c *TimestampAcceptanceValidation) tsaGeneralNameMatch() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.TimestampTSAGeneralNameContentMatch()
	return NewTSAGeneralNameValueMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// tsaGeneralNameOrderMatch ports the private tsaGeneralNameOrderMatch().
func (c *TimestampAcceptanceValidation) tsaGeneralNameOrderMatch() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.TimestampTSAGeneralNameOrderMatch()
	return NewTSAGeneralNameOrderMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

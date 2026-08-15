// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/RevocationAcceptanceValidation.java (DSS 6.5.RC1).
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
)

// RevocationAcceptanceValidation is 5.2.8 Signature acceptance validation (SAV)
// for a revocation token. This building block covers any additional
// verification to be performed on the signature itself or on the attributes of
// the signature ETSI EN 319 132-1.
type RevocationAcceptanceValidation struct {
	*AbstractAcceptanceValidation[*diagnostic.RevocationWrapper]
}

// NewRevocationAcceptanceValidation is the default constructor. Port of
// RevocationAcceptanceValidation(I18nProvider, Date, RevocationWrapper, XmlAOV, ValidationPolicy).
func NewRevocationAcceptanceValidation(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	revocationWrapper *diagnostic.RevocationWrapper, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *RevocationAcceptanceValidation {
	c := &RevocationAcceptanceValidation{
		AbstractAcceptanceValidation: NewAbstractAcceptanceValidation(i18nProvider, revocationWrapper, currentTime,
			enumerations.Context_REVOCATION, aovResult, validationPolicy),
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *RevocationAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *RevocationAcceptanceValidation) InitChain() {
	c.FirstItem = c.cryptographic(c.FirstItem)
}

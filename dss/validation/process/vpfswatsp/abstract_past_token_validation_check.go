// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/AbstractPastTokenValidationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// The class is abstract with two concrete subclasses, one here
// (PastSignatureValidationCheck) and one in vpftspwatsp
// (PastTimestampValidationCheck); neither of the base's method bodies calls an
// overridable method, so plain Go embedding is enough and no
// overrides-registration interface is needed. A subclass re-registers itself
// with InitChainItem so that the methods it does override dispatch correctly -
// the same pattern vpftspwatsp.TimestampMessageImprintCheck already follows.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// psvBlockSuffix is the past signature validation suffix. Port of the private
// static PSV_BLOCK_SUFFIX.
const psvBlockSuffix = "-PSV"

// AbstractPastTokenValidationCheck contains the main logic for
// PastSignatureValidation result check.
type AbstractPastTokenValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlPSV is the past signature validation.
	xmlPSV *jaxb.XmlPSV

	// indication is the Indication.
	indication enumerations.Indication

	// subIndication is the SubIndication.
	subIndication enumerations.SubIndication
}

// NewAbstractPastTokenValidationCheck is the default constructor. Port of the
// protected AbstractPastTokenValidationCheck(Provider, T, TokenProxy, XmlPSV, LevelRule).
func NewAbstractPastTokenValidationCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
	token diagnostic.TokenProxy, xmlPSV *jaxb.XmlPSV,
	constraint policy.LevelRule) *AbstractPastTokenValidationCheck[T] {
	c := &AbstractPastTokenValidationCheck[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()+psvBlockSuffix),
		xmlPSV:        xmlPSV,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AbstractPastTokenValidationCheck[T]) Process() bool {
	if !c.IsValid(&c.xmlPSV.XmlConstraintsConclusionContent) {
		c.indication = c.xmlPSV.Conclusion.Indication.Indication()
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		if c.xmlPSV.Conclusion.SubIndication != nil {
			c.subIndication = c.xmlPSV.Conclusion.SubIndication.SubIndication()
		} else {
			c.subIndication = ""
		}
		return false
	}
	return true
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AbstractPastTokenValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.indication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AbstractPastTokenValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.subIndication
}

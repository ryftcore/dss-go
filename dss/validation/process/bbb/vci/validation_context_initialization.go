// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/ValidationContextInitialization.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.bbb.vci and its .checks subpackage
// flatten into this single Go package vci (no name collisions), so the checks
// are referenced unqualified.
//
// Every check constructor takes *process.Result[*jaxb.XmlVCI] where Java takes
// the XmlVCI itself: see chain.go for why the generated result object has to be
// bound to the JAXB base structs it embeds.
package vci

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ValidationContextInitialization is 5.2.4 Validation context initialization.
// This building block initializes the validation constraints (chain constraints,
// cryptographic constraints, signature elements constraints) and parameters
// (X.509 validation parameters including trust anchors, certificate validation
// data) that will be used to validate the signature.
type ValidationContextInitialization struct {
	*process.ChainBase[*jaxb.XmlVCI]

	// signature is the signature to validate.
	signature *diagnostic.SignatureWrapper

	// context is the validation context.
	context enumerations.Context

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy
}

// NewValidationContextInitialization is the default constructor. Port of
// ValidationContextInitialization(Provider, SignatureWrapper, Context, ValidationPolicy).
func NewValidationContextInitialization(i18nProvider *i18n.Provider, signature *diagnostic.SignatureWrapper,
	context enumerations.Context, validationPolicy policy.ValidationPolicy) *ValidationContextInitialization {
	xmlVCI := &jaxb.XmlVCI{}
	c := &ValidationContextInitialization{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlVCI,
			&xmlVCI.XmlConstraintsConclusionContent, &xmlVCI.XmlConstraintsConclusionAttrs)),
		signature:        signature,
		context:          context,
		validationPolicy: validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *ValidationContextInitialization) Title() i18n.MessageTag {
	return i18n.MessageTagValidationContextInitialization
}

// InitChain initializes the chain. Port of initChain().
func (c *ValidationContextInitialization) InitChain() {

	item := c.signaturePolicyIdentifier()
	c.FirstItem = item

	if c.signature.IsPolicyPresent() &&
		string(enumerations.SignaturePolicyTypeImplicitPolicy) != c.signature.PolicyId() {

		item = item.SetNextItem(c.signaturePolicyIdentified())

		item = item.SetNextItem(c.signaturePolicyStorePresent())

		// Compare hash only when a policy is identified
		if c.signature.IsPolicyIdentified() {

			if !c.signature.IsPolicyZeroHash() {
				item = item.SetNextItem(c.signaturePolicyHashValid()) //nolint:staticcheck // mirrors upstream ValidationContextInitialization#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
			} else {
				item = item.SetNextItem(c.signaturePolicyZeroHash()) //nolint:staticcheck // mirrors upstream ValidationContextInitialization#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
			}

		}

	}

}

// signaturePolicyIdentifier ports the private signaturePolicyIdentifier().
func (c *ValidationContextInitialization) signaturePolicyIdentifier() process.ChainItem[*jaxb.XmlVCI] {
	signaturePolicyConstraint := c.validationPolicy.SignaturePolicyConstraint(c.context)
	return NewSignaturePolicyIdentifierCheck(c.I18nProvider, c.Result, c.signature, signaturePolicyConstraint)
}

// signaturePolicyIdentified ports the private signaturePolicyIdentified().
func (c *ValidationContextInitialization) signaturePolicyIdentified() process.ChainItem[*jaxb.XmlVCI] {
	constraint := c.validationPolicy.SignaturePolicyIdentifiedConstraint(c.context)
	return NewSignaturePolicyIdentifiedCheck(c.I18nProvider, c.Result, c.signature, constraint)
}

// signaturePolicyStorePresent ports the private signaturePolicyStorePresent().
func (c *ValidationContextInitialization) signaturePolicyStorePresent() process.ChainItem[*jaxb.XmlVCI] {
	constraint := c.validationPolicy.SignaturePolicyStorePresentConstraint(c.context)
	return NewSignaturePolicyStoreCheck(c.I18nProvider, c.Result, c.signature, constraint)
}

// signaturePolicyHashValid ports the private signaturePolicyHashValid().
func (c *ValidationContextInitialization) signaturePolicyHashValid() process.ChainItem[*jaxb.XmlVCI] {
	constraint := c.validationPolicy.SignaturePolicyPolicyHashValid(c.context)
	return NewSignaturePolicyHashValidCheck(c.I18nProvider, c.Result, c.signature, constraint)
}

// signaturePolicyZeroHash ports the private signaturePolicyZeroHash().
func (c *ValidationContextInitialization) signaturePolicyZeroHash() process.ChainItem[*jaxb.XmlVCI] {
	return NewSignaturePolicyZeroHashCheck(c.I18nProvider, c.Result, c.signature, c.WarnLevelRule())
}

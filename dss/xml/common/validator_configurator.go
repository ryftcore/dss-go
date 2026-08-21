// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/ValidatorConfigurator.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/alert"

// Validator is a documented stub standing in for javax.xml.validation.Validator. XSD schema
// validation is out of scope for this phase (see SchemaFactory's doc comment in
// schema_factory_builder.go); ErrorHandler exists so Configure/PostProcess's public contract
// - which wires a DSSErrorHandler onto the validator and later reads it back - is ready for
// the phase that adds a real Validator.
type Validator struct {
	ErrorHandler *DSSErrorHandler
}

// ValidatorConfigurator configures a provided Validator.
//
// # Documented stub
//
// Validator has no Go counterpart (see its doc comment), so Configure's
// setSecurityFeature/setSecurityAttribute are no-ops that always succeed: there is no real
// object to configure, and no real object to fail configuring either. The feature/attribute
// configuration machinery inherited from AbstractConfigurator is preserved and fully
// functional, ready for a later phase to give setSecurityFeature/setSecurityAttribute real
// behaviour without changing this type's shape.
type ValidatorConfigurator struct {
	*AbstractConfigurator[*Validator]

	// errorHandlerAlert is the alert used to process the errors collected during
	// validation. Default: DSSErrorHandlerAlert - collects exceptions and builds an
	// XSDValidationException.
	errorHandlerAlert alert.Alert[*DSSErrorHandler]
}

// newValidatorConfigurator builds a secure pre-configured instance. Ports the protected
// default constructor.
//
// The configuration protects against XXE
// (https://cheatsheetseries.owasp.org/cheatsheets/XML_External_Entity_Prevention_Cheat_Sheet.html#validator).
func newValidatorConfigurator() *ValidatorConfigurator {
	c := &ValidatorConfigurator{errorHandlerAlert: NewDSSErrorHandlerAlert()}
	c.AbstractConfigurator = NewAbstractConfigurator[*Validator](c.setSecurityFeature, c.setSecurityAttribute)
	c.SetAttribute(xmlConstantsAccessExternalDTD, "")
	c.SetAttribute(xmlConstantsAccessExternalSchema, "")
	return c
}

// GetSecureValidatorConfigurator instantiates a pre-configured, secure
// ValidatorConfigurator. Ports the static factory getSecureValidatorConfigurator().
func GetSecureValidatorConfigurator() *ValidatorConfigurator {
	return newValidatorConfigurator()
}

// Configure configures the given Validator by setting the pre-defined features and
// attributes. Panics if validator is nil (Java Objects.requireNonNull(validator, "Validator
// must be provided")).
func (c *ValidatorConfigurator) Configure(validator *Validator) error {
	if validator == nil {
		panic("Validator must be provided")
	}
	if err := c.SetSecurityFeatures(validator); err != nil {
		return err
	}
	if err := c.SetSecurityAttributes(validator); err != nil {
		return err
	}
	c.setErrorHandler(validator)
	return nil
}

// EnableFeature overrides AbstractConfigurator.EnableFeature purely for the covariant
// return, exactly as Java's override does.
func (c *ValidatorConfigurator) EnableFeature(feature string) *ValidatorConfigurator {
	c.AbstractConfigurator.EnableFeature(feature)
	return c
}

// DisableFeature overrides AbstractConfigurator.DisableFeature for the covariant return.
func (c *ValidatorConfigurator) DisableFeature(feature string) *ValidatorConfigurator {
	c.AbstractConfigurator.DisableFeature(feature)
	return c
}

// SetErrorHandlerAlert sets the alert used to process the collected exceptions during XML
// file validation. Panics if errorHandlerAlert is nil (Java Objects.requireNonNull).
func (c *ValidatorConfigurator) SetErrorHandlerAlert(errorHandlerAlert alert.Alert[*DSSErrorHandler]) {
	if errorHandlerAlert == nil {
		panic("errorHandlerAlert cannot be null!")
	}
	c.errorHandlerAlert = errorHandlerAlert
}

// SetAttribute overrides AbstractConfigurator.SetAttribute for the covariant return.
func (c *ValidatorConfigurator) SetAttribute(attribute string, value any) *ValidatorConfigurator {
	c.AbstractConfigurator.SetAttribute(attribute, value)
	return c
}

// RemoveAttribute overrides AbstractConfigurator.RemoveAttribute for the covariant return.
func (c *ValidatorConfigurator) RemoveAttribute(attribute string) *ValidatorConfigurator {
	c.AbstractConfigurator.RemoveAttribute(attribute)
	return c
}

// setSecurityFeature is a documented no-op (see the type doc comment).
func (c *ValidatorConfigurator) setSecurityFeature(validator *Validator, feature string, value bool) error {
	return nil
}

// setSecurityAttribute is a documented no-op (see the type doc comment).
func (c *ValidatorConfigurator) setSecurityAttribute(validator *Validator, attribute string, value any) error {
	return nil
}

// setErrorHandler wires a fresh DSSErrorHandler onto validator, so validation errors can be
// collected and later reported by PostProcess. Ports the protected
// setErrorHandler(Validator).
func (c *ValidatorConfigurator) setErrorHandler(validator *Validator) {
	validator.ErrorHandler = NewDSSErrorHandler()
}

// PostProcess handles the validation errors occurred during an XML file validation. Ports
// postProcess(Validator).
func (c *ValidatorConfigurator) PostProcess(validator *Validator) error {
	if validator.ErrorHandler != nil {
		return c.errorHandlerAlert.Alert(validator.ErrorHandler)
	}
	return nil
}

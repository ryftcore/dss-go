// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/AbstractConfigurator.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/alert"

// SecurityFeatureSetter configures a single named security feature on a factory-like value
// F. It stands in for the abstract method
// AbstractConfigurator#setSecurityFeature(F, String, Boolean): Java expresses it by
// subclassing and overriding; Go has no subclass to override from, so the concrete builder
// supplies the behaviour as an injected function instead - the same pattern already used by
// alert.AbstractAlert's detector/handler fields.
type SecurityFeatureSetter[F any] func(factory F, feature string, value bool) error

// SecurityAttributeSetter configures a single named security attribute on a factory-like
// value F. Stands in for the abstract method
// AbstractConfigurator#setSecurityAttribute(F, String, Object).
type SecurityAttributeSetter[F any] func(factory F, attribute string, value any) error

// AbstractConfigurator holds util methods helping to configure a Factory or a Validator.
type AbstractConfigurator[F any] struct {
	// featureNames/attributeNames record insertion order so features/attributes apply
	// deterministically (Go map iteration order is randomized).
	featureNames []string
	features     map[string]bool

	attributeNames []string
	attributes     map[string]any

	// securityExceptionAlert defines the behaviour for processing a security exception.
	securityExceptionAlert alert.StatusAlert

	setFeature   SecurityFeatureSetter[F]
	setAttribute SecurityAttributeSetter[F]
}

// NewAbstractConfigurator creates an AbstractConfigurator with empty feature/attribute
// tables and the default ExceptionOnStatusAlert, given the setters that stand in for the
// Java subclass's overridden setSecurityFeature/setSecurityAttribute.
func NewAbstractConfigurator[F any](setFeature SecurityFeatureSetter[F], setAttribute SecurityAttributeSetter[F]) *AbstractConfigurator[F] {
	return &AbstractConfigurator[F]{
		features:               make(map[string]bool),
		attributes:             make(map[string]any),
		securityExceptionAlert: alert.NewExceptionOnStatusAlert(),
		setFeature:             setFeature,
		setAttribute:           setAttribute,
	}
}

// SetSecurityExceptionAlert configures a custom alert on security exception in the builder.
// Panics if securityExceptionAlert is nil (Java Objects.requireNonNull).
func (c *AbstractConfigurator[F]) SetSecurityExceptionAlert(securityExceptionAlert alert.StatusAlert) {
	if securityExceptionAlert == nil {
		panic("securityExceptionAlert")
	}
	c.securityExceptionAlert = securityExceptionAlert
}

// EnableFeature enables a custom feature. Ports enableFeature(String).
func (c *AbstractConfigurator[F]) EnableFeature(feature string) *AbstractConfigurator[F] {
	return c.setFeatureValue(feature, true)
}

// DisableFeature disables a custom feature. Ports disableFeature(String).
func (c *AbstractConfigurator[F]) DisableFeature(feature string) *AbstractConfigurator[F] {
	return c.setFeatureValue(feature, false)
}

func (c *AbstractConfigurator[F]) setFeatureValue(feature string, value bool) *AbstractConfigurator[F] {
	if _, exists := c.features[feature]; !exists {
		c.featureNames = append(c.featureNames, feature)
	}
	c.features[feature] = value
	return c
}

// SetAttribute sets a custom attribute. Ports setAttribute(String, Object).
func (c *AbstractConfigurator[F]) SetAttribute(attribute string, value any) *AbstractConfigurator[F] {
	if _, exists := c.attributes[attribute]; !exists {
		c.attributeNames = append(c.attributeNames, attribute)
	}
	c.attributes[attribute] = value
	return c
}

// RemoveAttribute removes the attribute from the list of attributes to set. Ports
// removeAttribute(String).
func (c *AbstractConfigurator[F]) RemoveAttribute(attribute string) *AbstractConfigurator[F] {
	if _, exists := c.attributes[attribute]; exists {
		delete(c.attributes, attribute)
		for i, name := range c.attributeNames {
			if name == attribute {
				c.attributeNames = append(c.attributeNames[:i], c.attributeNames[i+1:]...)
				break
			}
		}
	}
	return c
}

// SetSecurityFeatures sets all configured features to the factory, collecting failures into
// an ObjectStatus and alerting exactly as Java's void setSecurityFeatures(F) does. Java's
// alert is an unchecked throw; here it is a returned error, per PORTING.md ("dss-alert
// handlers ... in Go ... may return an error").
func (c *AbstractConfigurator[F]) SetSecurityFeatures(factory F) error {
	status := alert.NewObjectStatus()
	for _, name := range c.featureNames {
		if err := c.setFeature(factory, name, c.features[name]); err != nil {
			status.AddRelatedObjectIdentifierAndErrorMessage(name, err.Error())
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("SECURITY : unable to set feature(s)!")
		return c.securityExceptionAlert.Alert(status)
	}
	return nil
}

// SetSecurityAttributes sets all configured attributes to the factory. See
// SetSecurityFeatures for the void-to-error mapping.
func (c *AbstractConfigurator[F]) SetSecurityAttributes(factory F) error {
	status := alert.NewObjectStatus()
	for _, name := range c.attributeNames {
		if err := c.setAttribute(factory, name, c.attributes[name]); err != nil {
			status.AddRelatedObjectIdentifierAndErrorMessage(name, err.Error())
		}
	}
	if !status.IsEmpty() {
		status.SetMessage("SECURITY : unable to set attribute(s)!")
		return c.securityExceptionAlert.Alert(status)
	}
	return nil
}

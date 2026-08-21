// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/GeneralName.java (DSS 6.5.RC1).
package extension

import "github.com/ryftcore/dss-go/dss/enumerations"

// GeneralName represents a general name element (see RFC 5280).
type GeneralName struct {
	// generalNameType represents the type of the GeneralName.
	generalNameType enumerations.GeneralNameType

	// value is the string representation of the GeneralName value.
	value string
}

// NewGeneralName instantiates the object with null values. Ports the default constructor.
func NewGeneralName() *GeneralName {
	return &GeneralName{}
}

// GeneralNameType gets the type of GeneralName.
func (g *GeneralName) GeneralNameType() enumerations.GeneralNameType {
	return g.generalNameType
}

// SetGeneralNameType sets the type of the GeneralName.
func (g *GeneralName) SetGeneralNameType(generalNameType enumerations.GeneralNameType) {
	g.generalNameType = generalNameType
}

// Value gets the string representation of the GeneralName value.
func (g *GeneralName) Value() string {
	return g.value
}

// SetValue sets the string representation of the GeneralName value.
func (g *GeneralName) SetValue(value string) {
	g.value = value
}

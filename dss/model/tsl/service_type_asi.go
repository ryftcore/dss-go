// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/ServiceTypeASi.java (DSS 6.5.RC1).
package tsl

import "fmt"

// ServiceTypeASi contains information extracted from the TrustServiceTSLType element.
//
// java.io.Serializable has no Go counterpart and is dropped.
type ServiceTypeASi struct {
	// type is the ServiceTypeIdentifier value.
	typ string
	// asi is the AdditionalServiceInformation value.
	asi string
}

// NewServiceTypeASi instantiates an object with zero values. Port of the default constructor.
func NewServiceTypeASi() *ServiceTypeASi {
	return &ServiceTypeASi{}
}

// Type gets the ServiceTypeIdentifier value.
func (s *ServiceTypeASi) Type() string {
	return s.typ
}

// SetType sets the ServiceTypeIdentifier value.
func (s *ServiceTypeASi) SetType(typ string) {
	s.typ = typ
}

// Asi gets the AdditionalServiceInformation value.
func (s *ServiceTypeASi) Asi() string {
	return s.asi
}

// SetAsi sets the AdditionalServiceInformation value.
func (s *ServiceTypeASi) SetAsi(asi string) {
	s.asi = asi
}

// String returns the Java toString() form.
func (s *ServiceTypeASi) String() string {
	return fmt.Sprintf("ServiceTypeASi [type='%s', asi='%s']", s.typ, s.asi)
}

// Equals ports ServiceTypeASi#equals(Object).
func (s *ServiceTypeASi) Equals(other *ServiceTypeASi) bool {
	if other == nil {
		return false
	}
	if s == other {
		return true
	}
	return s.typ == other.typ && s.asi == other.asi
}

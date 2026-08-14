// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/MRA.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/model/timedependent"

// MRA contains information extracted from the MutualRecognitionAgreementInformation element of
// a Mutual Recognition Agreement schema.
//
// java.io.Serializable has no Go counterpart and is dropped.
type MRA struct {
	// technicalType is the value of the technicalType attribute.
	technicalType string
	// version is the value of the version attribute.
	version string
	// pointingContractingPartyLegislation references the legal documentation of the pointing party.
	pointingContractingPartyLegislation string
	// pointedContractingPartyLegislation references the legal documentation of the pointed party.
	pointedContractingPartyLegislation string
	// serviceEquivalence contains a list of equivalence schemes defined for various Trust Services.
	serviceEquivalence []*timedependent.MutableTimeDependentValues[*ServiceEquivalence]
}

// NewMRA instantiates an MRA object with zero values. Port of the default constructor.
func NewMRA() *MRA {
	return &MRA{}
}

// TechnicalType gets the technical type attribute value.
func (m *MRA) TechnicalType() string {
	return m.technicalType
}

// SetTechnicalType sets the technical type attribute value.
func (m *MRA) SetTechnicalType(technicalType string) {
	m.technicalType = technicalType
}

// Version gets the version attribute value.
func (m *MRA) Version() string {
	return m.version
}

// SetVersion sets the version attribute value.
func (m *MRA) SetVersion(version string) {
	m.version = version
}

// PointingContractingPartyLegislation gets the value defined within the
// pointingContractingPartyLegislation attribute.
func (m *MRA) PointingContractingPartyLegislation() string {
	return m.pointingContractingPartyLegislation
}

// SetPointingContractingPartyLegislation sets the value defined within the
// pointingContractingPartyLegislation attribute.
func (m *MRA) SetPointingContractingPartyLegislation(pointingContractingPartyLegislation string) {
	m.pointingContractingPartyLegislation = pointingContractingPartyLegislation
}

// PointedContractingPartyLegislation gets the value defined within the
// pointedContractingPartyLegislation attribute.
func (m *MRA) PointedContractingPartyLegislation() string {
	return m.pointedContractingPartyLegislation
}

// SetPointedContractingPartyLegislation sets the value defined within the
// pointedContractingPartyLegislation attribute.
func (m *MRA) SetPointedContractingPartyLegislation(pointedContractingPartyLegislation string) {
	m.pointedContractingPartyLegislation = pointedContractingPartyLegislation
}

// ServiceEquivalence gets the list of equivalence mapping between Trust Services.
func (m *MRA) ServiceEquivalence() []*timedependent.MutableTimeDependentValues[*ServiceEquivalence] {
	return m.serviceEquivalence
}

// SetServiceEquivalence sets the list of equivalence mapping between Trust Services.
func (m *MRA) SetServiceEquivalence(serviceEquivalence []*timedependent.MutableTimeDependentValues[*ServiceEquivalence]) {
	m.serviceEquivalence = serviceEquivalence
}

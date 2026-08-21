// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/ServiceEquivalence.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

// StatusEquivalenceMapping is a single entry of the TrustServiceTSLStatusEquivalenceList
// equivalences between pointed and pointing parties.
//
// DEVIATION: Java keys a Map<List<String>, List<String>> by structural List equality/hashCode,
// which Go maps cannot do (slice keys are not comparable). Per PORTING.md's collections rule,
// the order-sensitive Java map becomes a slice of pairs instead.
type StatusEquivalenceMapping struct {
	// PointedStatuses is the key side of the Java map entry (pointed party statuses).
	PointedStatuses []string
	// PointingStatuses is the value side of the Java map entry (pointing party statuses).
	PointingStatuses []string
}

// ServiceEquivalence is a wrapper for the TrustServiceEquivalenceInformation element from the
// MRA scheme.
//
// java.io.Serializable has no Go counterpart and is dropped.
type ServiceEquivalence struct {
	*timedependent.BaseTimeDependent

	// legalInfoIdentifier is the TrustServiceLegalIdentifier.
	legalInfoIdentifier string
	// status is the TrustServiceEquivalenceStatus.
	status enumerations.MRAStatus
	// typeAsiEquivalence holds the AdditionalServiceInformation equivalencies.
	typeAsiEquivalence map[ServiceTypeASi]ServiceTypeASi
	// statusEquivalence holds the TrustServiceTSLStatusEquivalenceList equivalencies.
	statusEquivalence []StatusEquivalenceMapping
	// certificateContentEquivalences is the CertificateContentReferencesEquivalenceList.
	certificateContentEquivalences []*CertificateContentEquivalence
	// qualifierEquivalence holds the QualifierEquivalenceList equivalencies.
	qualifierEquivalence map[string]string
}

// ServiceEquivalenceBuilder builds a ServiceEquivalence object.
type ServiceEquivalenceBuilder struct {
	// legalInfoIdentifier is the TrustServiceLegalIdentifier.
	legalInfoIdentifier string
	// status is the TrustServiceEquivalenceStatus.
	status enumerations.MRAStatus
	// startDate is the TrustServiceEquivalenceStatusStartingTime.
	startDate time.Time
	// endDate is the start date of the next TrustServiceEquivalenceHistoryInstance or
	// TrustServiceEquivalenceInformationType.
	endDate time.Time
	// typeAsiEquivalence holds the AdditionalServiceInformation equivalencies.
	typeAsiEquivalence map[ServiceTypeASi]ServiceTypeASi
	// statusEquivalence holds the TrustServiceTSLStatusEquivalenceList equivalencies.
	statusEquivalence []StatusEquivalenceMapping
	// certificateContentEquivalences is the CertificateContentReferencesEquivalenceList.
	certificateContentEquivalences []*CertificateContentEquivalence
	// qualifierEquivalence holds the QualifierEquivalenceList equivalencies.
	qualifierEquivalence map[string]string
}

// NewServiceEquivalenceBuilder instantiates a builder with zero values. Port of the default
// constructor.
func NewServiceEquivalenceBuilder() *ServiceEquivalenceBuilder {
	return &ServiceEquivalenceBuilder{}
}

// Build builds the ServiceEquivalence object.
func (b *ServiceEquivalenceBuilder) Build() *ServiceEquivalence {
	return &ServiceEquivalence{
		BaseTimeDependent:              timedependent.NewBaseTimeDependentWithDates(b.startDate, b.endDate),
		legalInfoIdentifier:            b.legalInfoIdentifier,
		status:                         b.status,
		typeAsiEquivalence:             b.typeAsiEquivalence,
		statusEquivalence:              b.statusEquivalence,
		certificateContentEquivalences: b.certificateContentEquivalences,
		qualifierEquivalence:           b.qualifierEquivalence,
	}
}

// SetLegalInfoIdentifier sets the TrustServiceLegalIdentifier value.
func (b *ServiceEquivalenceBuilder) SetLegalInfoIdentifier(legalInfoIdentifier string) *ServiceEquivalenceBuilder {
	b.legalInfoIdentifier = legalInfoIdentifier
	return b
}

// SetStatus sets the TrustServiceEquivalenceStatus value.
func (b *ServiceEquivalenceBuilder) SetStatus(status enumerations.MRAStatus) *ServiceEquivalenceBuilder {
	b.status = status
	return b
}

// SetStartDate sets the TrustServiceEquivalenceStatusStartingTime value.
func (b *ServiceEquivalenceBuilder) SetStartDate(startDate time.Time) *ServiceEquivalenceBuilder {
	b.startDate = startDate
	return b
}

// SetEndDate sets the endDate (equivalent to the starting date of the following service
// equivalence) value.
func (b *ServiceEquivalenceBuilder) SetEndDate(endDate time.Time) *ServiceEquivalenceBuilder {
	b.endDate = endDate
	return b
}

// SetTypeAsiEquivalence sets a map of AdditionalServiceInformation equivalences between pointed
// and pointing parties.
func (b *ServiceEquivalenceBuilder) SetTypeAsiEquivalence(typeAsiEquivalence map[ServiceTypeASi]ServiceTypeASi) *ServiceEquivalenceBuilder {
	b.typeAsiEquivalence = typeAsiEquivalence
	return b
}

// SetStatusEquivalence sets the TrustServiceTSLStatusEquivalenceList equivalences between
// pointed and pointing parties.
func (b *ServiceEquivalenceBuilder) SetStatusEquivalence(statusEquivalence []StatusEquivalenceMapping) *ServiceEquivalenceBuilder {
	b.statusEquivalence = statusEquivalence
	return b
}

// SetCertificateContentEquivalences sets a list of CertificateContentReferencesEquivalenceList
// equivalences.
func (b *ServiceEquivalenceBuilder) SetCertificateContentEquivalences(certificateContentEquivalences []*CertificateContentEquivalence) *ServiceEquivalenceBuilder {
	b.certificateContentEquivalences = certificateContentEquivalences
	return b
}

// SetQualifierEquivalence sets a map of QualifierEquivalenceList equivalences between pointed
// and pointing parties.
func (b *ServiceEquivalenceBuilder) SetQualifierEquivalence(qualifierEquivalence map[string]string) *ServiceEquivalenceBuilder {
	b.qualifierEquivalence = qualifierEquivalence
	return b
}

// LegalInfoIdentifier gets the TrustServiceLegalIdentifier value.
func (s *ServiceEquivalence) LegalInfoIdentifier() string {
	return s.legalInfoIdentifier
}

// SetLegalInfoIdentifier sets the TrustServiceLegalIdentifier value.
func (s *ServiceEquivalence) SetLegalInfoIdentifier(legalInfoIdentifier string) {
	s.legalInfoIdentifier = legalInfoIdentifier
}

// Status gets the TrustServiceEquivalenceStatus value.
func (s *ServiceEquivalence) Status() enumerations.MRAStatus {
	return s.status
}

// SetStatus sets the TrustServiceEquivalenceStatus value.
func (s *ServiceEquivalence) SetStatus(status enumerations.MRAStatus) {
	s.status = status
}

// TypeAsiEquivalence gets a map of AdditionalServiceInformation equivalences between pointed
// and pointing parties.
func (s *ServiceEquivalence) TypeAsiEquivalence() map[ServiceTypeASi]ServiceTypeASi {
	return s.typeAsiEquivalence
}

// SetTypeAsiEquivalence sets a map of AdditionalServiceInformation equivalences between pointed
// and pointing parties.
func (s *ServiceEquivalence) SetTypeAsiEquivalence(typeAsiEquivalence map[ServiceTypeASi]ServiceTypeASi) {
	s.typeAsiEquivalence = typeAsiEquivalence
}

// StatusEquivalence gets the TrustServiceTSLStatusEquivalenceList equivalences between pointed
// and pointing parties.
func (s *ServiceEquivalence) StatusEquivalence() []StatusEquivalenceMapping {
	return s.statusEquivalence
}

// SetStatusEquivalence sets the TrustServiceTSLStatusEquivalenceList equivalences between
// pointed and pointing parties.
func (s *ServiceEquivalence) SetStatusEquivalence(statusEquivalence []StatusEquivalenceMapping) {
	s.statusEquivalence = statusEquivalence
}

// CertificateContentEquivalences gets the CertificateContentReferencesEquivalenceList
// equivalences.
func (s *ServiceEquivalence) CertificateContentEquivalences() []*CertificateContentEquivalence {
	return s.certificateContentEquivalences
}

// SetCertificateContentEquivalences sets the CertificateContentReferencesEquivalenceList
// equivalences.
func (s *ServiceEquivalence) SetCertificateContentEquivalences(certificateContentEquivalences []*CertificateContentEquivalence) {
	s.certificateContentEquivalences = certificateContentEquivalences
}

// QualifierEquivalence gets a map of QualifierEquivalenceList equivalences between pointed and
// pointing parties.
func (s *ServiceEquivalence) QualifierEquivalence() map[string]string {
	return s.qualifierEquivalence
}

// SetQualifierEquivalence sets a map of QualifierEquivalenceList equivalences between pointed
// and pointing parties.
func (s *ServiceEquivalence) SetQualifierEquivalence(qualifierEquivalence map[string]string) {
	s.qualifierEquivalence = qualifierEquivalence
}

// String returns the Java toString() form.
func (s *ServiceEquivalence) String() string {
	return fmt.Sprintf("ServiceEquivalence [legalInfoIdentifier='%s', status=%v, typeAsiEquivalence=%v, "+
		"statusEquivalence=%v, certificateContentEquivalences=%v, qualifierEquivalence=%v] %s",
		s.legalInfoIdentifier, s.status, s.typeAsiEquivalence, s.statusEquivalence,
		s.certificateContentEquivalences, s.qualifierEquivalence, s.BaseTimeDependent.String())
}

// Equals ports ServiceEquivalence#equals(Object), including the super.equals() call against the
// embedded BaseTimeDependent (start/end date range).
func (s *ServiceEquivalence) Equals(other *ServiceEquivalence) bool {
	if other == nil {
		return false
	}
	if s == other {
		return true
	}
	if !s.BaseTimeDependent.Equals(other.BaseTimeDependent) {
		return false
	}
	return s.legalInfoIdentifier == other.legalInfoIdentifier &&
		s.status == other.status &&
		reflect.DeepEqual(s.typeAsiEquivalence, other.typeAsiEquivalence) &&
		reflect.DeepEqual(s.statusEquivalence, other.statusEquivalence) &&
		reflect.DeepEqual(s.certificateContentEquivalences, other.certificateContentEquivalences) &&
		reflect.DeepEqual(s.qualifierEquivalence, other.qualifierEquivalence)
}

// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustServiceStatusAndInformationExtensions.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

// TrustServiceStatusAndInformationExtensions defines information for a TrustService.
type TrustServiceStatusAndInformationExtensions struct {
	*timedependent.BaseTimeDependent

	// names maps lang -> values for that lang.
	names map[string][]string
	// typ is the type.
	typ string
	// status is the status.
	status string
	// conditionsForQualifiers is a list of condition for qualifiers.
	conditionsForQualifiers []*ConditionForQualifiers
	// additionalServiceInfoUris is additional service info urls.
	additionalServiceInfoUris []string
	// serviceSupplyPoints is the list of service supply points.
	serviceSupplyPoints []string
	// expiredCertsRevocationInfo is the expired certs revocation info date. The zero time.Time
	// stands in for Java's null.
	expiredCertsRevocationInfo time.Time
}

// TrustServiceStatusAndInformationExtensionsBuilder builds
// TrustServiceStatusAndInformationExtensions.
type TrustServiceStatusAndInformationExtensionsBuilder struct {
	names                      map[string][]string
	typ                        string
	status                     string
	conditionsForQualifiers    []*ConditionForQualifiers
	additionalServiceInfoUris  []string
	serviceSupplyPoints        []string
	expiredCertsRevocationInfo time.Time
	startDate                  time.Time
	endDate                    time.Time
}

// NewTrustServiceStatusAndInformationExtensionsBuilder is the default constructor.
func NewTrustServiceStatusAndInformationExtensionsBuilder() *TrustServiceStatusAndInformationExtensionsBuilder {
	return &TrustServiceStatusAndInformationExtensionsBuilder{}
}

// NewTrustServiceStatusAndInformationExtensionsBuilderFrom is the constructor with an existing
// TrustServiceStatusAndInformationExtensions, copying its fields into the builder. Port of the
// TrustServiceStatusAndInformationExtensionsBuilder(TrustServiceStatusAndInformationExtensions)
// constructor.
func NewTrustServiceStatusAndInformationExtensionsBuilderFrom(status *TrustServiceStatusAndInformationExtensions) *TrustServiceStatusAndInformationExtensionsBuilder {
	return &TrustServiceStatusAndInformationExtensionsBuilder{
		names:                      status.Names(),
		typ:                        status.Type(),
		status:                     status.Status(),
		conditionsForQualifiers:    status.ConditionsForQualifiers(),
		additionalServiceInfoUris:  status.AdditionalServiceInfoUris(),
		serviceSupplyPoints:        status.ServiceSupplyPoints(),
		expiredCertsRevocationInfo: status.ExpiredCertsRevocationInfo(),
		startDate:                  status.StartDate(),
		endDate:                    status.EndDate(),
	}
}

// Build builds the TrustServiceStatusAndInformationExtensions.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) Build() *TrustServiceStatusAndInformationExtensions {
	return NewTrustServiceStatusAndInformationExtensions(b)
}

// SetNames sets a map of names.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetNames(names map[string][]string) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.names = names
	return b
}

// SetType sets a type.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetType(typ string) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.typ = typ
	return b
}

// SetStatus sets a status.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetStatus(status string) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.status = status
	return b
}

// SetConditionsForQualifiers sets conditions for qualifiers.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetConditionsForQualifiers(conditionsForQualifiers []*ConditionForQualifiers) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.conditionsForQualifiers = conditionsForQualifiers
	return b
}

// SetAdditionalServiceInfoUris sets additional service info urls.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetAdditionalServiceInfoUris(additionalServiceInfoUris []string) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.additionalServiceInfoUris = additionalServiceInfoUris
	return b
}

// SetServiceSupplyPoints sets the service supply points.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetServiceSupplyPoints(serviceSupplyPoints []string) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.serviceSupplyPoints = serviceSupplyPoints
	return b
}

// SetExpiredCertsRevocationInfo sets the expired certs revocation info date.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetExpiredCertsRevocationInfo(expiredCertsRevocationInfo time.Time) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.expiredCertsRevocationInfo = expiredCertsRevocationInfo
	return b
}

// SetStartDate sets the start of validity date.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetStartDate(date time.Time) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.startDate = date
	return b
}

// SetEndDate sets the end of validity date.
func (b *TrustServiceStatusAndInformationExtensionsBuilder) SetEndDate(date time.Time) *TrustServiceStatusAndInformationExtensionsBuilder {
	b.endDate = date
	return b
}

// NewTrustServiceStatusAndInformationExtensions is the default constructor.
//
// Panics with "TrustServiceStatusAndInformationExtensionsBuilder cannot be null!" when builder
// is nil, mirroring Objects.requireNonNull.
func NewTrustServiceStatusAndInformationExtensions(builder *TrustServiceStatusAndInformationExtensionsBuilder) *TrustServiceStatusAndInformationExtensions {
	if builder == nil {
		panic("TrustServiceStatusAndInformationExtensionsBuilder cannot be null!")
	}
	return &TrustServiceStatusAndInformationExtensions{
		BaseTimeDependent:          timedependent.NewBaseTimeDependentWithDates(builder.startDate, builder.endDate),
		names:                      builder.names,
		typ:                        builder.typ,
		status:                     builder.status,
		conditionsForQualifiers:    builder.conditionsForQualifiers,
		additionalServiceInfoUris:  builder.additionalServiceInfoUris,
		serviceSupplyPoints:        builder.serviceSupplyPoints,
		expiredCertsRevocationInfo: builder.expiredCertsRevocationInfo,
	}
}

// Names gets a map of names.
func (t *TrustServiceStatusAndInformationExtensions) Names() map[string][]string {
	return t.names
}

// Type gets the type.
func (t *TrustServiceStatusAndInformationExtensions) Type() string {
	return t.typ
}

// Status gets the status.
func (t *TrustServiceStatusAndInformationExtensions) Status() string {
	return t.status
}

// ConditionsForQualifiers gets a list of conditions for qualifiers.
func (t *TrustServiceStatusAndInformationExtensions) ConditionsForQualifiers() []*ConditionForQualifiers {
	return t.conditionsForQualifiers
}

// AdditionalServiceInfoUris gets additional service info urls.
func (t *TrustServiceStatusAndInformationExtensions) AdditionalServiceInfoUris() []string {
	return t.additionalServiceInfoUris
}

// ServiceSupplyPoints gets service supply points.
func (t *TrustServiceStatusAndInformationExtensions) ServiceSupplyPoints() []string {
	return t.serviceSupplyPoints
}

// ExpiredCertsRevocationInfo gets the expired certs revocation info date.
func (t *TrustServiceStatusAndInformationExtensions) ExpiredCertsRevocationInfo() time.Time {
	return t.expiredCertsRevocationInfo
}

// String returns the Java toString() form.
func (t *TrustServiceStatusAndInformationExtensions) String() string {
	return fmt.Sprintf("TrustServiceStatusAndInformationExtensions [names=%v, type='%s', status='%s', "+
		"conditionsForQualifiers=%v, additionalServiceInfoUris=%v, serviceSupplyPoints=%v, "+
		"expiredCertsRevocationInfo=%v] %s",
		t.names, t.typ, t.status, t.conditionsForQualifiers, t.additionalServiceInfoUris,
		t.serviceSupplyPoints, t.expiredCertsRevocationInfo, t.BaseTimeDependent.String())
}

// Equals ports TrustServiceStatusAndInformationExtensions#equals(Object), including the
// super.equals() call against the embedded BaseTimeDependent (start/end date range).
func (t *TrustServiceStatusAndInformationExtensions) Equals(other *TrustServiceStatusAndInformationExtensions) bool {
	if other == nil {
		return false
	}
	if t == other {
		return true
	}
	if !t.BaseTimeDependent.Equals(other.BaseTimeDependent) {
		return false
	}
	return reflect.DeepEqual(t.names, other.names) &&
		t.typ == other.typ &&
		t.status == other.status &&
		reflect.DeepEqual(t.conditionsForQualifiers, other.conditionsForQualifiers) &&
		reflect.DeepEqual(t.additionalServiceInfoUris, other.additionalServiceInfoUris) &&
		reflect.DeepEqual(t.serviceSupplyPoints, other.serviceSupplyPoints) &&
		t.expiredCertsRevocationInfo.Equal(other.expiredCertsRevocationInfo)
}

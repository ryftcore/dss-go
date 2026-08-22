// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/TrustedEntityServiceStatusAndInformationExtensions.java (DSS 6.5.RC1).
package lote

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

// TrustedEntityServiceStatusAndInformationExtensions contains information about the service
// status and other information for the given time period.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustedEntityServiceStatusAndInformationExtensions struct {
	*timedependent.BaseTimeDependent

	// names maps lang -> values for that lang.
	names map[string][]string
	// typ is the type.
	typ string
	// status is the status.
	status string
	// serviceSupplyPoints is the list of service supply points.
	serviceSupplyPoints []string
}

// NewTrustedEntityServiceStatusAndInformationExtensions is the default constructor.
//
// Panics with "ServiceStatusAndInformationExtensionsBuilder cannot be null!" when builder is
// nil, mirroring Objects.requireNonNull.
func NewTrustedEntityServiceStatusAndInformationExtensions(builder *ServiceStatusAndInformationExtensionsBuilder) *TrustedEntityServiceStatusAndInformationExtensions {
	if builder == nil {
		panic("ServiceStatusAndInformationExtensionsBuilder cannot be null!")
	}
	return &TrustedEntityServiceStatusAndInformationExtensions{
		BaseTimeDependent:   timedependent.NewBaseTimeDependentWithDates(builder.startDate, builder.endDate),
		names:               builder.names,
		typ:                 builder.typ,
		status:              builder.status,
		serviceSupplyPoints: builder.serviceSupplyPoints,
	}
}

// Names gets a map of names.
func (t *TrustedEntityServiceStatusAndInformationExtensions) Names() map[string][]string {
	return t.names
}

// Type gets type.
func (t *TrustedEntityServiceStatusAndInformationExtensions) Type() string {
	return t.typ
}

// Status gets status.
func (t *TrustedEntityServiceStatusAndInformationExtensions) Status() string {
	return t.status
}

// ServiceSupplyPoints gets service supply points.
func (t *TrustedEntityServiceStatusAndInformationExtensions) ServiceSupplyPoints() []string {
	return t.serviceSupplyPoints
}

// ServiceStatusAndInformationExtensionsBuilder builds
// TrustedEntityServiceStatusAndInformationExtensions.
type ServiceStatusAndInformationExtensionsBuilder struct {
	// names maps lang -> values for that lang.
	names map[string][]string
	// typ is the type.
	typ string
	// status is the status.
	status string
	// serviceSupplyPoints is the list of service supply points.
	serviceSupplyPoints []string
	// startDate is the start of validity date.
	startDate time.Time
	// endDate is the end of validity date.
	endDate time.Time
}

// NewServiceStatusAndInformationExtensionsBuilder is the default constructor.
func NewServiceStatusAndInformationExtensionsBuilder() *ServiceStatusAndInformationExtensionsBuilder {
	return &ServiceStatusAndInformationExtensionsBuilder{}
}

// NewServiceStatusAndInformationExtensionsBuilderFrom is the constructor with an existing
// TrustedEntityServiceStatusAndInformationExtensions, copying its fields into the builder.
func NewServiceStatusAndInformationExtensionsBuilderFrom(status *TrustedEntityServiceStatusAndInformationExtensions) *ServiceStatusAndInformationExtensionsBuilder {
	return &ServiceStatusAndInformationExtensionsBuilder{
		names:               status.Names(),
		typ:                 status.Type(),
		status:              status.Status(),
		serviceSupplyPoints: status.ServiceSupplyPoints(),
		startDate:           status.StartDate(),
		endDate:             status.EndDate(),
	}
}

// Build builds TrustedEntityServiceStatusAndInformationExtensions.
func (b *ServiceStatusAndInformationExtensionsBuilder) Build() *TrustedEntityServiceStatusAndInformationExtensions {
	return NewTrustedEntityServiceStatusAndInformationExtensions(b)
}

// SetNames sets a map of names.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetNames(names map[string][]string) *ServiceStatusAndInformationExtensionsBuilder {
	b.names = names
	return b
}

// SetType sets a type.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetType(typ string) *ServiceStatusAndInformationExtensionsBuilder {
	b.typ = typ
	return b
}

// SetStatus sets a status.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetStatus(status string) *ServiceStatusAndInformationExtensionsBuilder {
	b.status = status
	return b
}

// SetServiceSupplyPoints sets the service supply points.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetServiceSupplyPoints(serviceSupplyPoints []string) *ServiceStatusAndInformationExtensionsBuilder {
	b.serviceSupplyPoints = serviceSupplyPoints
	return b
}

// SetStartDate sets the start of validity date.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetStartDate(date time.Time) *ServiceStatusAndInformationExtensionsBuilder {
	b.startDate = date
	return b
}

// SetEndDate sets the end of validity date.
func (b *ServiceStatusAndInformationExtensionsBuilder) SetEndDate(date time.Time) *ServiceStatusAndInformationExtensionsBuilder {
	b.endDate = date
	return b
}

// compile-time interface assertion.
var _ ServiceStatusAndInformationExtensions = (*TrustedEntityServiceStatusAndInformationExtensions)(nil)

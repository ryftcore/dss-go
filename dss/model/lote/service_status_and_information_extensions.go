// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/ServiceStatusAndInformationExtensions.java (DSS 6.5.RC1).
package lote

import "github.com/utain/esig/dss/model/timedependent"

// ServiceStatusAndInformationExtensions contains information about service status and
// extensions.
type ServiceStatusAndInformationExtensions interface {
	timedependent.TimeDependent

	// Names gets a map of names.
	Names() map[string][]string
	// Type gets type.
	Type() string
	// Status gets status.
	Status() string
	// ServiceSupplyPoints gets service supply points.
	ServiceSupplyPoints() []string
}

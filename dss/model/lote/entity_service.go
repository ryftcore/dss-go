// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/EntityService.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

// EntityService contains information about an entity service. S is the concrete
// ServiceStatusAndInformationExtensions implementation, mirroring Java's
// "<S extends ServiceStatusAndInformationExtensions>" type parameter.
//
// java.io.Serializable has no Go counterpart and is dropped.
type EntityService[S ServiceStatusAndInformationExtensions] interface {
	// Certificates gets a list of certificates.
	Certificates() []*model.CertificateToken
	// StatusAndInformationExtensions gets status based on time.
	StatusAndInformationExtensions() *timedependent.Values[S]
}

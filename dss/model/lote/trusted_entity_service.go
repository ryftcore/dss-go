// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/TrustedEntityService.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/timedependent"
)

// TrustedEntityService contains information about a single trusted entity's service.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustedEntityService struct {
	// certificates is a list of certificates.
	certificates []*model.CertificateToken
	// status holds statuses based on time.
	status *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]
}

// NewTrustedEntityService is the default constructor.
func NewTrustedEntityService(certificates []*model.CertificateToken,
	status *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]) *TrustedEntityService {
	return &TrustedEntityService{certificates: certificates, status: status}
}

// Certificates gets a list of certificates.
func (t *TrustedEntityService) Certificates() []*model.CertificateToken {
	return t.certificates
}

// StatusAndInformationExtensions gets status based on time.
func (t *TrustedEntityService) StatusAndInformationExtensions() *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions] {
	return t.status
}

// TrustEntityServiceBuilder builds TrustedEntityService. Java's inner class is named
// "TrustEntityServiceBuilder" (without "ed" between "Trust" and "Entity"), kept verbatim.
type TrustEntityServiceBuilder struct {
	// certificates is a list of certificates.
	certificates []*model.CertificateToken
	// status holds statuses based on time.
	status *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]
}

// NewTrustEntityServiceBuilder is the default constructor.
func NewTrustEntityServiceBuilder() *TrustEntityServiceBuilder {
	return &TrustEntityServiceBuilder{}
}

// SetCertificates sets a list of certificates.
func (b *TrustEntityServiceBuilder) SetCertificates(certificates []*model.CertificateToken) *TrustEntityServiceBuilder {
	b.certificates = certificates
	return b
}

// SetStatusAndInformationExtensions sets a status.
func (b *TrustEntityServiceBuilder) SetStatusAndInformationExtensions(
	status *timedependent.TimeDependentValues[ServiceStatusAndInformationExtensions]) *TrustEntityServiceBuilder {
	b.status = status
	return b
}

// Build builds TrustedEntityService.
func (b *TrustEntityServiceBuilder) Build() *TrustedEntityService {
	return NewTrustedEntityService(b.certificates, b.status)
}

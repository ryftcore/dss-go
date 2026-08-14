// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustService.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/timedependent"
)

// TrustService is a DTO representation for a TSL service.
//
// java.io.Serializable has no Go counterpart and is dropped.
type TrustService struct {
	// certificates is a list of certificates.
	certificates []*model.CertificateToken
	// status holds statuses based on time.
	status *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]
}

// NewTrustService is the default constructor.
func NewTrustService(certificates []*model.CertificateToken,
	status *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]) *TrustService {
	return &TrustService{certificates: certificates, status: status}
}

// Certificates gets a list of certificates.
func (t *TrustService) Certificates() []*model.CertificateToken {
	return t.certificates
}

// StatusAndInformationExtensions gets status based on time.
func (t *TrustService) StatusAndInformationExtensions() *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions] {
	return t.status
}

// TrustServiceBuilder builds a TrustService.
type TrustServiceBuilder struct {
	// certificates is a list of certificates.
	certificates []*model.CertificateToken
	// status holds statuses based on time.
	status *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]
}

// NewTrustServiceBuilder is the default constructor.
func NewTrustServiceBuilder() *TrustServiceBuilder {
	return &TrustServiceBuilder{}
}

// SetCertificates sets a list of certificates.
func (b *TrustServiceBuilder) SetCertificates(certificates []*model.CertificateToken) *TrustServiceBuilder {
	b.certificates = certificates
	return b
}

// SetStatusAndInformationExtensions sets a status.
func (b *TrustServiceBuilder) SetStatusAndInformationExtensions(status *timedependent.TimeDependentValues[*TrustServiceStatusAndInformationExtensions]) *TrustServiceBuilder {
	b.status = status
	return b
}

// Build builds the TrustService.
func (b *TrustServiceBuilder) Build() *TrustService {
	return NewTrustService(b.certificates, b.status)
}

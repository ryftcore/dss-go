// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/lote/TrustedEntitiesCertificateSource.java (DSS 6.5.RC1).
//
// Deviation: Java's protected CommonCertificateSource#reset() is called via super.reset() from
// the same inheritance chain in a different package (protected access crosses packages in
// Java). The Go port of that method (spi.CommonCertificateSource.reset) is unexported and this
// type lives in a different Go package, so it has no access to it. This port reproduces
// reset()'s effect - discarding all previously added certificates - by replacing the embedded
// CommonTrustedCertificateSource with a freshly constructed one, which starts from the same
// empty maps reset() would restore. Flagged for integrator reconciliation.
package lote

import (
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/lote"
	"github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TrustedEntitiesCertificateSource is a certificate source built based on trusted entities.
type TrustedEntitiesCertificateSource struct {
	spi.CommonTrustedCertificateSource

	// summary is the TL Validation job summary.
	summary *lote.LoTEValidationJobSummary

	// trustPropertiesByEntity is the map of trust properties by EntityIdentifier (public
	// keys), keyed by EntityIdentifier.AsXmlID().
	trustPropertiesByEntity map[string][]*lote.TrustedProperties

	// trustTimeByEntity is the map of trust time periods by EntityIdentifier, keyed by
	// EntityIdentifier.AsXmlID().
	trustTimeByEntity map[string][]*tsl.CertificateTrustTime
}

// NewTrustedEntitiesCertificateSource is the default constructor.
func NewTrustedEntitiesCertificateSource() *TrustedEntitiesCertificateSource {
	return &TrustedEntitiesCertificateSource{
		CommonTrustedCertificateSource: *spi.NewCommonTrustedCertificateSource(),
		trustPropertiesByEntity:        make(map[string][]*lote.TrustedProperties),
		trustTimeByEntity:              make(map[string][]*tsl.CertificateTrustTime),
	}
}

// Summary gets LoTE Validation job summary.
func (s *TrustedEntitiesCertificateSource) Summary() *lote.LoTEValidationJobSummary {
	return s.summary
}

// SetSummary sets LoTE Validation job summary.
func (s *TrustedEntitiesCertificateSource) SetSummary(summary *lote.LoTEValidationJobSummary) {
	s.summary = summary
}

// CertificateSourceType returns CertificateSourceType_TRUSTED_ENTITIES.
func (s *TrustedEntitiesCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_TRUSTED_ENTITIES
}

// AddCertificate is not applicable for this kind of certificate source: it panics. You should
// use SetTrustedPropertiesByCertificates.
func (s *TrustedEntitiesCertificateSource) AddCertificate(certificate *model.CertificateToken) *model.CertificateToken {
	panic("Cannot directly add certificate to a TrustedListsCertificateSource")
}

// SetTrustedPropertiesByCertificates reinitializes the source with the given certificates and
// their TrustedProperties. Panics if trustPropertiesByCerts is nil (Java
// Objects.requireNonNull).
func (s *TrustedEntitiesCertificateSource) SetTrustedPropertiesByCertificates(trustPropertiesByCerts map[*model.CertificateToken][]*lote.TrustedProperties) {
	if trustPropertiesByCerts == nil {
		panic("TrustedPropertiesByCerts cannot be null!")
	}
	s.trustPropertiesByEntity = make(map[string][]*lote.TrustedProperties)
	s.CommonTrustedCertificateSource = *spi.NewCommonTrustedCertificateSource()
	for certificateToken, trustPropertiesList := range trustPropertiesByCerts {
		s.addCertificateWithTrustedProperties(certificateToken, trustPropertiesList)
	}
}

// addCertificateWithTrustedProperties ports the private addCertificate(CertificateToken,
// List<TrustedProperties>).
func (s *TrustedEntitiesCertificateSource) addCertificateWithTrustedProperties(certificateToken *model.CertificateToken, trustPropertiesList []*lote.TrustedProperties) {
	s.CommonTrustedCertificateSource.AddCertificate(certificateToken)
	if trustPropertiesList == nil {
		panic("TrustedPropertiesList must be filled")
	}

	entityKey := certificateToken.EntityKey().AsXmlID()
	list := s.trustPropertiesByEntity[entityKey]
	for _, trustProperties := range trustPropertiesList {
		if !trustedEntitiesCertificateSourceContainsTrustedProperties(list, trustProperties) {
			list = append(list, trustProperties)
		}
	}
	s.trustPropertiesByEntity[entityKey] = list
}

// trustedEntitiesCertificateSourceContainsTrustedProperties ports List#contains for
// TrustedProperties, which relies on Java's default (reference) equality since TrustedProperties
// does not override equals().
func trustedEntitiesCertificateSourceContainsTrustedProperties(list []*lote.TrustedProperties, value *lote.TrustedProperties) bool {
	for _, entry := range list {
		if entry == value {
			return true
		}
	}
	return false
}

// TrustedProperties returns TrustedProperties for the given certificate, when applicable.
func (s *TrustedEntitiesCertificateSource) TrustedProperties(token *model.CertificateToken) []*lote.TrustedProperties {
	if currentTrustedProperties, found := s.trustPropertiesByEntity[token.EntityKey().AsXmlID()]; found {
		return currentTrustedProperties
	}
	return nil
}

// SetTrustedTimeByCertificates reinitializes the trust time map with the given certificates and
// their CertificateTrustTime entries. Panics if trustTimeByCertificate is nil (Java
// Objects.requireNonNull).
func (s *TrustedEntitiesCertificateSource) SetTrustedTimeByCertificates(trustTimeByCertificate map[*model.CertificateToken][]*tsl.CertificateTrustTime) {
	if trustTimeByCertificate == nil {
		panic("trustTimeByCertificate cannot be null!")
	}
	s.trustTimeByEntity = make(map[string][]*tsl.CertificateTrustTime)
	for certificateToken, certificateTrustTimes := range trustTimeByCertificate {
		s.addCertificateTrustTimes(certificateToken, certificateTrustTimes)
	}
}

// addCertificateTrustTimes ports the private addCertificateTrustTimes(CertificateToken,
// List<CertificateTrustTime>).
func (s *TrustedEntitiesCertificateSource) addCertificateTrustTimes(certificateToken *model.CertificateToken, certificateTrustTimes []*tsl.CertificateTrustTime) {
	s.CommonTrustedCertificateSource.AddCertificate(certificateToken)
	if certificateTrustTimes == nil {
		panic("CertificateTrustTimes must be filled")
	}

	entityKey := certificateToken.EntityKey().AsXmlID()
	list := s.trustTimeByEntity[entityKey]
	for _, trustTime := range certificateTrustTimes {
		if !trustedEntitiesCertificateSourceContainsTrustTime(list, trustTime) {
			list = append(list, trustTime)
		}
	}
	s.trustTimeByEntity[entityKey] = list
}

// trustedEntitiesCertificateSourceContainsTrustTime ports List#contains for CertificateTrustTime,
// using its Equals method (CertificateTrustTime does override equals() in Java).
func trustedEntitiesCertificateSourceContainsTrustTime(list []*tsl.CertificateTrustTime, value *tsl.CertificateTrustTime) bool {
	for _, entry := range list {
		if entry.Equals(value) {
			return true
		}
	}
	return false
}

// TrustTime returns the trust time period for the given certificate.
func (s *TrustedEntitiesCertificateSource) TrustTime(token *model.CertificateToken) *tsl.CertificateTrustTime {
	if !s.CommonTrustedCertificateSource.IsTrusted(token) {
		return tsl.NewCertificateTrustTime(false)
	}
	trustTimes := s.trustTimeByEntity[token.EntityKey().AsXmlID()]
	if utils.IsCollectionNotEmpty(trustTimes) {
		var certificateTrustTime *tsl.CertificateTrustTime
		for _, trustTime := range trustTimes {
			if certificateTrustTime == nil || !certificateTrustTime.IsTrusted() {
				certificateTrustTime = trustTime
			} else if trustTime != nil && trustTime.IsTrusted() {
				certificateTrustTime = certificateTrustTime.JointTrustTime(trustTime.StartDate(), trustTime.EndDate())
			}
		}
		return certificateTrustTime
	}
	// no trust anchor expiration time defined
	return tsl.NewCertificateTrustTime(true)
}

// IsTrustedAtTime checks if a given certificate is trusted at controlTime.
func (s *TrustedEntitiesCertificateSource) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	return s.TrustTime(certificateToken).IsTrustedAtTime(controlTime)
}

// AlternativeOCSPUrls returns the service supply points found for the "ocsp" keyword.
func (s *TrustedEntitiesCertificateSource) AlternativeOCSPUrls(trustAnchor *model.CertificateToken) []string {
	return s.getServiceSupplyPoints(trustAnchor, "ocsp")
}

// AlternativeCRLUrls returns the service supply points found for the "crl" and
// "certificateRevocationList" keywords.
func (s *TrustedEntitiesCertificateSource) AlternativeCRLUrls(trustAnchor *model.CertificateToken) []string {
	return s.getServiceSupplyPoints(trustAnchor, "crl", "certificateRevocationList")
}

// getServiceSupplyPoints ports the private getServiceSupplyPoints(CertificateToken, String...).
func (s *TrustedEntitiesCertificateSource) getServiceSupplyPoints(trustAnchor *model.CertificateToken, keywords ...string) []string {
	var urls []string
	trustPropertiesList := s.TrustedProperties(trustAnchor)
	for _, trustProperties := range trustPropertiesList {
		for statusAndInfo := range trustProperties.TrustedServices().Iterator() {
			serviceSupplyPoints := statusAndInfo.ServiceSupplyPoints()
			if utils.IsCollectionNotEmpty(serviceSupplyPoints) {
				for _, serviceSupplyPoint := range serviceSupplyPoints {
					for _, keyword := range keywords {
						if strings.Contains(serviceSupplyPoint, keyword) {
							urls = append(urls, serviceSupplyPoint)
						}
					}
				}
			}
		}
	}
	return urls
}

// IsTrusted checks if a given certificate is trusted, taking its trust time period into
// account.
func (s *TrustedEntitiesCertificateSource) IsTrusted(certificateToken *model.CertificateToken) bool {
	if s.CommonTrustedCertificateSource.IsTrusted(certificateToken) {
		trustTime := s.TrustTime(certificateToken)
		return trustTime == nil || trustTime.IsTrusted()
	}
	return false
}

// NumberOfTrustedEntityKeys gets the number of trusted entity keys (public key + subject name).
func (s *TrustedEntitiesCertificateSource) NumberOfTrustedEntityKeys() int {
	return len(s.trustPropertiesByEntity)
}

var _ lote.TrustedPropertiesCertificateSource = (*TrustedEntitiesCertificateSource)(nil)

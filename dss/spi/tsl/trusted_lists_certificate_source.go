// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/tsl/TrustedListsCertificateSource.java (DSS 6.5.RC1).
//
// Deviation: see the equivalent note in spi/lote/trusted_entities_certificate_source.go
// (this package's sibling chunk) regarding the inability to call spi.CommonCertificateSource's
// unexported reset() from a different Go package; reset()'s effect is reproduced here the same
// way, by replacing the embedded CommonTrustedCertificateSource with a fresh one.
package tsl

import (
	"sort"
	"strings"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// TrustedListsCertificateSource allows injection of trusted certificates from Trusted Lists.
type TrustedListsCertificateSource struct {
	spi.CommonTrustedCertificateSource

	// summary is the TL Validation job summary.
	summary *tsl.TLValidationJobSummary

	// trustPropertiesByEntity is the map of trust properties by EntityIdentifier (public
	// keys), keyed by EntityIdentifier.AsXmlID().
	trustPropertiesByEntity map[string][]*tsl.TrustProperties

	// trustTimeByEntity is the map of trust time periods by EntityIdentifier, keyed by
	// EntityIdentifier.AsXmlID().
	trustTimeByEntity map[string][]*tsl.CertificateTrustTime
}

// NewTrustedListsCertificateSource is the default constructor.
func NewTrustedListsCertificateSource() *TrustedListsCertificateSource {
	return &TrustedListsCertificateSource{
		CommonTrustedCertificateSource: *spi.NewCommonTrustedCertificateSource(),
		trustPropertiesByEntity:        make(map[string][]*tsl.TrustProperties),
		trustTimeByEntity:              make(map[string][]*tsl.CertificateTrustTime),
	}
}

// Summary gets TL Validation job summary.
func (s *TrustedListsCertificateSource) Summary() *tsl.TLValidationJobSummary {
	return s.summary
}

// SetSummary sets TL Validation job summary.
func (s *TrustedListsCertificateSource) SetSummary(summary *tsl.TLValidationJobSummary) {
	s.summary = summary
}

// CertificateSourceType returns CertificateSourceType_TRUSTED_LIST.
func (s *TrustedListsCertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceType_TRUSTED_LIST
}

// AddCertificate is not applicable for this kind of certificate source: it panics. You should
// use SetTrustPropertiesByCertificates.
func (s *TrustedListsCertificateSource) AddCertificate(certificate *model.CertificateToken) *model.CertificateToken {
	panic("Cannot directly add certificate to a TrustedListsCertificateSource")
}

// SetTrustPropertiesByCertificates reinitializes the source with the given certificates and
// their TrustProperties. Panics if trustPropertiesByCerts is nil (Java
// Objects.requireNonNull).
func (s *TrustedListsCertificateSource) SetTrustPropertiesByCertificates(trustPropertiesByCerts map[*model.CertificateToken][]*tsl.TrustProperties) {
	if trustPropertiesByCerts == nil {
		panic("TrustPropertiesByCerts cannot be null!")
	}
	s.trustPropertiesByEntity = make(map[string][]*tsl.TrustProperties)
	s.CommonTrustedCertificateSource = *spi.NewCommonTrustedCertificateSource()
	for _, certificateToken := range trustedListsCertificateSourceSortedTokens(trustPropertiesByCerts) {
		s.addCertificate(certificateToken, trustPropertiesByCerts[certificateToken])
	}
}

// trustedListsCertificateSourceSortedTokens returns the map's certificate keys ordered by their
// DSS id.
//
// Java iterates the Map<CertificateToken, ?> these two setters are handed through its
// entrySet(); a java.util.HashMap's order is arbitrary but, for a given set of keys, the SAME on
// every run, so upstream's resulting certificate order is stable. Go map iteration is randomised
// per range, so ranging directly would make getCertificates() - and everything downstream that
// consumes it, up to diagnostic data - come out in a different order on every run. Sorting by DSS
// id restores a stable order (the same substitute dss/validation/reports/diagnostic's
// sortedTrustAnchors already uses for the identical Java shape); it is not Java's own bucket
// order, which is a function of the certificate digest's byte-array hashCode and is not
// reproducible here, and no ported caller depends on which arbitrary order it gets.
func trustedListsCertificateSourceSortedTokens[V any](m map[*model.CertificateToken]V) []*model.CertificateToken {
	tokens := make([]*model.CertificateToken, 0, len(m))
	for token := range m {
		tokens = append(tokens, token)
	}
	sort.SliceStable(tokens, func(i, j int) bool {
		return tokens[i].DSSIDAsString() < tokens[j].DSSIDAsString()
	})
	return tokens
}

// addCertificate ports the private addCertificate(CertificateToken, List<TrustProperties>).
func (s *TrustedListsCertificateSource) addCertificate(certificateToken *model.CertificateToken, trustPropertiesList []*tsl.TrustProperties) {
	s.CommonTrustedCertificateSource.AddCertificate(certificateToken)
	if trustPropertiesList == nil {
		panic("TrustPropertiesList must be filled")
	}

	entityKey := certificateToken.EntityKey().AsXmlID()
	list := s.trustPropertiesByEntity[entityKey]
	for _, trustProperties := range trustPropertiesList {
		if !trustedListsCertificateSourceContainsTrustProperties(list, trustProperties) {
			list = append(list, trustProperties)
		}
	}
	s.trustPropertiesByEntity[entityKey] = list
}

// trustedListsCertificateSourceContainsTrustProperties ports List#contains for TrustProperties,
// which relies on Java's default (reference) equality since TrustProperties does not override
// equals().
func trustedListsCertificateSourceContainsTrustProperties(list []*tsl.TrustProperties, value *tsl.TrustProperties) bool {
	for _, entry := range list {
		if entry == value {
			return true
		}
	}
	return false
}

// TrustServices returns TrustProperties for the given certificate, when applicable.
func (s *TrustedListsCertificateSource) TrustServices(token *model.CertificateToken) []*tsl.TrustProperties {
	if currentTrustProperties, found := s.trustPropertiesByEntity[token.EntityKey().AsXmlID()]; found {
		return currentTrustProperties
	}
	return nil
}

// SetTrustTimeByCertificates reinitializes the trust time map with the given certificates and
// their CertificateTrustTime entries. Panics if trustTimeByCertificate is nil (Java
// Objects.requireNonNull).
func (s *TrustedListsCertificateSource) SetTrustTimeByCertificates(trustTimeByCertificate map[*model.CertificateToken][]*tsl.CertificateTrustTime) {
	if trustTimeByCertificate == nil {
		panic("trustTimeByCertificate cannot be null!")
	}
	s.trustTimeByEntity = make(map[string][]*tsl.CertificateTrustTime)
	for _, certificateToken := range trustedListsCertificateSourceSortedTokens(trustTimeByCertificate) {
		s.addCertificateTrustTimes(certificateToken, trustTimeByCertificate[certificateToken])
	}
}

// addCertificateTrustTimes ports the private addCertificateTrustTimes(CertificateToken,
// List<CertificateTrustTime>).
func (s *TrustedListsCertificateSource) addCertificateTrustTimes(certificateToken *model.CertificateToken, certificateTrustTimes []*tsl.CertificateTrustTime) {
	s.CommonTrustedCertificateSource.AddCertificate(certificateToken)
	if certificateTrustTimes == nil {
		panic("CertificateTrustTimes must be filled")
	}

	entityKey := certificateToken.EntityKey().AsXmlID()
	list := s.trustTimeByEntity[entityKey]
	for _, trustTime := range certificateTrustTimes {
		if !trustedListsCertificateSourceContainsTrustTime(list, trustTime) {
			list = append(list, trustTime)
		}
	}
	s.trustTimeByEntity[entityKey] = list
}

// trustedListsCertificateSourceContainsTrustTime ports List#contains for CertificateTrustTime,
// using its Equals method (CertificateTrustTime does override equals() in Java).
func trustedListsCertificateSourceContainsTrustTime(list []*tsl.CertificateTrustTime, value *tsl.CertificateTrustTime) bool {
	for _, entry := range list {
		if entry.Equals(value) {
			return true
		}
	}
	return false
}

// TrustTime returns the trust time period for the given certificate.
func (s *TrustedListsCertificateSource) TrustTime(token *model.CertificateToken) *tsl.CertificateTrustTime {
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
func (s *TrustedListsCertificateSource) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	return s.TrustTime(certificateToken).IsTrustedAtTime(controlTime)
}

// AlternativeOCSPUrls returns the service supply points found for the "ocsp" keyword.
func (s *TrustedListsCertificateSource) AlternativeOCSPUrls(trustAnchor *model.CertificateToken) []string {
	return s.getServiceSupplyPoints(trustAnchor, "ocsp")
}

// AlternativeCRLUrls returns the service supply points found for the "crl" and
// "certificateRevocationList" keywords.
func (s *TrustedListsCertificateSource) AlternativeCRLUrls(trustAnchor *model.CertificateToken) []string {
	return s.getServiceSupplyPoints(trustAnchor, "crl", "certificateRevocationList")
}

// getServiceSupplyPoints ports the private getServiceSupplyPoints(CertificateToken, String...).
func (s *TrustedListsCertificateSource) getServiceSupplyPoints(trustAnchor *model.CertificateToken, keywords ...string) []string {
	var urls []string
	trustPropertiesList := s.TrustServices(trustAnchor)
	for _, trustProperties := range trustPropertiesList {
		for statusAndInfo := range trustProperties.TrustService().Iterator() {
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
func (s *TrustedListsCertificateSource) IsTrusted(certificateToken *model.CertificateToken) bool {
	if s.CommonTrustedCertificateSource.IsTrusted(certificateToken) {
		trustTime := s.TrustTime(certificateToken)
		return trustTime == nil || trustTime.IsTrusted()
	}
	return false
}

// NumberOfTrustedEntityKeys gets the number of trusted entity keys (public key + subject name).
func (s *TrustedListsCertificateSource) NumberOfTrustedEntityKeys() int {
	return len(s.trustPropertiesByEntity)
}

var _ tsl.TrustPropertiesCertificateSource = (*TrustedListsCertificateSource)(nil)

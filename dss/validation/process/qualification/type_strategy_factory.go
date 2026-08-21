// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeStrategyFactory.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CreateTypeFromCert creates a TypeStrategy from the certificate. Port of
// createTypeFromCert(CertificateWrapper).
func CreateTypeFromCert(signingCertificate *diagnostic.CertificateWrapper) TypeStrategy {
	if IsPostEIDAS(signingCertificate.NotBefore()) {
		return newTypeByCertificatePostEIDAS(signingCertificate)
	}
	return newTypeByCertificatePreEIDAS(signingCertificate)
}

// CreateTypeFromTL creates a TypeStrategy from a Trusted Service. Port of
// createTypeFromTL(TrustServiceWrapper, CertificateQualifiedStatus, TypeStrategy).
func CreateTypeFromTL(trustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus,
	typeInCert TypeStrategy) TypeStrategy {
	return newTypeByTL(trustService, qualified, typeInCert)
}

// CreateTypeFromCertAndTL creates a TypeStrategy from a certificate and
// Trusted Service. Port of
// createTypeFromCertAndTL(CertificateWrapper, TrustServiceWrapper, CertificateQualifiedStatus).
func CreateTypeFromCertAndTL(signingCertificate *diagnostic.CertificateWrapper,
	caQcTrustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus) TypeStrategy {
	typeFromCert := CreateTypeFromCert(signingCertificate)
	return CreateTypeFromTL(caQcTrustService, qualified, typeFromCert)
}

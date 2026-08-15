// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationStrategyFactory.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// CreateQualificationFromCert creates a QualificationStrategy from the
// certificate. Port of createQualificationFromCert(CertificateWrapper).
func CreateQualificationFromCert(signingCertificate *diagnostic.CertificateWrapper) QualificationStrategy {
	if IsPostEIDAS(signingCertificate.NotBefore()) {
		return newQualificationByCertificatePostEIDAS(signingCertificate)
	}
	return newQualificationByCertificatePreEIDAS(signingCertificate)
}

// CreateQualificationFromTL creates a QualificationStrategy from the Trusted
// Service. Port of
// createQualificationFromTL(TrustServiceWrapper, QualificationStrategy).
func CreateQualificationFromTL(trustService *diagnostic.TrustServiceWrapper, qualifiedInCert QualificationStrategy) QualificationStrategy {
	return newQualificationByTL(trustService, qualifiedInCert)
}

// CreateQualificationFromCertAndTL creates a QualificationStrategy from the
// certificate and Trusted Service. Port of
// createQualificationFromCertAndTL(CertificateWrapper, TrustServiceWrapper).
func CreateQualificationFromCertAndTL(signingCertificate *diagnostic.CertificateWrapper,
	caQcTrustService *diagnostic.TrustServiceWrapper) QualificationStrategy {
	qcFromCert := CreateQualificationFromCert(signingCertificate)
	return CreateQualificationFromTL(caQcTrustService, qcFromCert)
}

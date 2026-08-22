// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationStrategyFactory.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// CreateQualificationFromCert creates a Strategy from the
// certificate. Port of createQualificationFromCert(CertificateWrapper).
func CreateQualificationFromCert(signingCertificate *diagnostic.CertificateWrapper) Strategy {
	if IsPostEIDAS(signingCertificate.NotBefore()) {
		return newQualificationByCertificatePostEIDAS(signingCertificate)
	}
	return newQualificationByCertificatePreEIDAS(signingCertificate)
}

// CreateQualificationFromTL creates a Strategy from the Trusted
// Service. Port of
// createQualificationFromTL(TrustServiceWrapper, Strategy).
func CreateQualificationFromTL(trustService *diagnostic.TrustServiceWrapper, qualifiedInCert Strategy) Strategy {
	return newQualificationByTL(trustService, qualifiedInCert)
}

// CreateQualificationFromCertAndTL creates a Strategy from the
// certificate and Trusted Service. Port of
// createQualificationFromCertAndTL(CertificateWrapper, TrustServiceWrapper).
func CreateQualificationFromCertAndTL(signingCertificate *diagnostic.CertificateWrapper,
	caQcTrustService *diagnostic.TrustServiceWrapper) Strategy {
	qcFromCert := CreateQualificationFromCert(signingCertificate)
	return CreateQualificationFromTL(caQcTrustService, qcFromCert)
}

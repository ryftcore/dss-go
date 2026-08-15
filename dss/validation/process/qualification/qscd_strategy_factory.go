// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qscd/QSCDStrategyFactory.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
)

// CreateQSCDFromCert creates a QSCD Strategy from the given certificate.
// Port of createQSCDFromCert(CertificateWrapper).
func CreateQSCDFromCert(signingCertificate *diagnostic.CertificateWrapper) QSCDStrategy {
	if IsPostEIDAS(signingCertificate.NotBefore()) {
		return newQSCDByCertificatePostEIDAS(signingCertificate)
	}
	return newQSCDByCertificatePreEIDAS(signingCertificate)
}

// CreateQSCDFromTL creates a QSCD Strategy from the TrustService. Port of
// createQSCDFromTL(TrustServiceWrapper, CertificateQualifiedStatus, QSCDStrategy).
func CreateQSCDFromTL(trustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus,
	qscdFromCertificate QSCDStrategy) QSCDStrategy {
	return newQSCDByTL(trustService, qualified, qscdFromCertificate)
}

// CreateQSCDFromCertAndTL creates a QSCD Strategy from the given certificate
// and TrustService. Port of
// createQSCDFromCertAndTL(CertificateWrapper, TrustServiceWrapper, CertificateQualifiedStatus).
func CreateQSCDFromCertAndTL(signingCertificate *diagnostic.CertificateWrapper,
	caQcTrustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus) QSCDStrategy {
	qscdFromCert := CreateQSCDFromCert(signingCertificate)
	return CreateQSCDFromTL(caQcTrustService, qualified, qscdFromCert)
}

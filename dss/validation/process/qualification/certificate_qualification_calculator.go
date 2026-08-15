// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/CertificateQualificationCalculator.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
)

// CertificateQualificationCalculator is used to determine certificate's
// qualification based on its content and the given TrustServiceWrapper.
type CertificateQualificationCalculator struct {
	// endEntityCert is the certificate to get qualification for.
	endEntityCert *diagnostic.CertificateWrapper

	// caqcTrustService is the Trusted Service to be used to help to
	// determine the qualification.
	caqcTrustService *diagnostic.TrustServiceWrapper
}

// NewCertificateQualificationCalculator is the default constructor. Port of
// CertificateQualificationCalculator(CertificateWrapper, TrustServiceWrapper).
func NewCertificateQualificationCalculator(endEntityCert *diagnostic.CertificateWrapper,
	caqcTrustService *diagnostic.TrustServiceWrapper) *CertificateQualificationCalculator {
	return &CertificateQualificationCalculator{endEntityCert: endEntityCert, caqcTrustService: caqcTrustService}
}

// Qualification returns the qualification result for the given
// CertificateWrapper. Port of getQualification().
func (c *CertificateQualificationCalculator) Qualification() enumerations.CertificateQualification {
	qcStrategy := CreateQualificationFromCertAndTL(c.endEntityCert, c.caqcTrustService)
	qualifiedStatus := qcStrategy.QualifiedStatus()

	typeStrategy := CreateTypeFromCertAndTL(c.endEntityCert, c.caqcTrustService, qualifiedStatus)
	certType := typeStrategy.Type()

	qscdStrategy := CreateQSCDFromCertAndTL(c.endEntityCert, c.caqcTrustService, qualifiedStatus)
	qscdStatus := qscdStrategy.QSCDStatus()

	return GetCertQualification(qualifiedStatus, certType, qscdStatus)
}

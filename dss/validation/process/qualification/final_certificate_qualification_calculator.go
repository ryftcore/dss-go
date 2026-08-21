// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/FinalCertificateQualificationCalculator.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/enumerations"

// FinalCertificateQualificationCalculator determines the final qualification
// of a certificate given two qualifications at issuance and
// best-signature-time.
type FinalCertificateQualificationCalculator struct {
	// certQualificationAtIssuanceTime represents a qualification of a
	// certificate at its issuance time (notBefore).
	certQualificationAtIssuanceTime enumerations.CertificateQualification

	// certQualificationAtBestSignatureTime represents a qualification of a
	// certificate at the best-signature-time.
	certQualificationAtBestSignatureTime enumerations.CertificateQualification
}

// NewFinalCertificateQualificationCalculator is the default constructor.
// Port of
// FinalCertificateQualificationCalculator(CertificateQualification, CertificateQualification).
func NewFinalCertificateQualificationCalculator(certQualificationAtIssuanceTime,
	certQualificationAtBestSignatureTime enumerations.CertificateQualification) *FinalCertificateQualificationCalculator {
	return &FinalCertificateQualificationCalculator{
		certQualificationAtIssuanceTime:      certQualificationAtIssuanceTime,
		certQualificationAtBestSignatureTime: certQualificationAtBestSignatureTime,
	}
}

// FinalQualification returns the final qualification of the certificate,
// based on the two given times. Port of getFinalQualification().
func (c *FinalCertificateQualificationCalculator) FinalQualification() enumerations.CertificateQualification {
	if c.certQualificationAtIssuanceTime == c.certQualificationAtBestSignatureTime {
		return c.certQualificationAtIssuanceTime
	}
	if enumerations.CertificateQualification_NA == c.certQualificationAtIssuanceTime ||
		enumerations.CertificateQualification_NA == c.certQualificationAtBestSignatureTime {
		return enumerations.CertificateQualification_NA
	}
	qualStatus := c.getFinalCertQualStatus()
	certType := c.getFinalCertificateType()
	qscd := c.getFinalQSCDStatus()
	return GetCertQualification(qualStatus, certType, qscd)
}

// getFinalCertQualStatus ports the private getFinalCertQualStatus().
func (c *FinalCertificateQualificationCalculator) getFinalCertQualStatus() enumerations.CertificateQualifiedStatus {
	if c.certQualificationAtIssuanceTime.IsQc() && c.certQualificationAtBestSignatureTime.IsQc() {
		return enumerations.CertificateQualifiedStatus_QC
	}
	return enumerations.CertificateQualifiedStatus_NOT_QC
}

// getFinalCertificateType ports the private getFinalCertificateType().
func (c *FinalCertificateQualificationCalculator) getFinalCertificateType() enumerations.CertificateType {
	if c.certQualificationAtIssuanceTime.Type() == c.certQualificationAtBestSignatureTime.Type() {
		return c.certQualificationAtBestSignatureTime.Type()
	}
	return enumerations.CertificateType_UNKNOWN
}

// getFinalQSCDStatus ports the private getFinalQSCDStatus().
func (c *FinalCertificateQualificationCalculator) getFinalQSCDStatus() enumerations.QSCDStatus {
	if c.certQualificationAtBestSignatureTime.IsQscd() {
		return enumerations.QSCDStatus_QSCD
	}
	return enumerations.QSCDStatus_NOT_QSCD
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/QCTypeIdentifiers.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
)

// IsQCTypeEsign checks whether the certificate contains a QC for eSignature
// qualifier (oid "0.4.0.1862.1.6.1"). Port of isQCTypeEsign(CertificateWrapper).
func IsQCTypeEsign(certificate *diagnostic.CertificateWrapper) bool {
	return hasQCTypeOID(certificate, enumerations.QCTypeEnum_QCT_ESIGN)
}

// IsQCTypeEseal checks whether the certificate contains a QC for eSeal
// qualifier (oid "0.4.0.1862.1.6.2"). Port of isQCTypeEseal(CertificateWrapper).
func IsQCTypeEseal(certificate *diagnostic.CertificateWrapper) bool {
	return hasQCTypeOID(certificate, enumerations.QCTypeEnum_QCT_ESEAL)
}

// IsQCTypeWeb checks whether the certificate contains a QC for Web
// Authentication qualifier (oid "0.4.0.1862.1.6.3"). Port of
// isQCTypeWeb(CertificateWrapper).
func IsQCTypeWeb(certificate *diagnostic.CertificateWrapper) bool {
	return hasQCTypeOID(certificate, enumerations.QCTypeEnum_QCT_WEB)
}

// hasQCTypeOID ports the private static hasQCTypeOID(CertificateWrapper, QCType).
func hasQCTypeOID(certificate *diagnostic.CertificateWrapper, qcType enumerations.QCType) bool {
	qcTypes := certificate.QcTypes()
	if utils.IsCollectionNotEmpty(qcTypes) {
		for _, t := range qcTypes {
			if t == qcType {
				return true
			}
		}
	}
	return false
}

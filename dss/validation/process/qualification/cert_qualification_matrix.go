// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/CertQualificationMatrix.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/enumerations"

// Indices into the certQualifications cube-array. Port of the private static
// final int constants NOT_QC/QC/ESIG/ESEAL/WSA/UNKNOWN/NOT_QSCD/QSCD.
const (
	certQualNotQC = 0
	certQualQC    = 1

	certQualESig    = 0
	certQualESeal   = 1
	certQualWSA     = 2
	certQualUnknown = 3

	certQualNotQSCD = 0
	certQualQSCD    = 1
)

// certQualifications is the cached cube-array containing qualification
// results for different sets of parameters. Port of the static QUALIFS
// array and its static initializer.
var certQualifications = func() [2][4][2]enumerations.CertificateQualification {
	var q [2][4][2]enumerations.CertificateQualification

	q[certQualQC][certQualESig][certQualQSCD] = enumerations.CertificateQualification_QCERT_FOR_ESIG_QSCD
	q[certQualQC][certQualESeal][certQualQSCD] = enumerations.CertificateQualification_QCERT_FOR_ESEAL_QSCD
	q[certQualQC][certQualWSA][certQualQSCD] = enumerations.CertificateQualification_QCERT_FOR_WSA
	q[certQualQC][certQualUnknown][certQualQSCD] = enumerations.CertificateQualification_QCERT_FOR_UNKNOWN_QSCD

	q[certQualQC][certQualESig][certQualNotQSCD] = enumerations.CertificateQualification_QCERT_FOR_ESIG
	q[certQualQC][certQualESeal][certQualNotQSCD] = enumerations.CertificateQualification_QCERT_FOR_ESEAL
	q[certQualQC][certQualWSA][certQualNotQSCD] = enumerations.CertificateQualification_QCERT_FOR_WSA
	q[certQualQC][certQualUnknown][certQualNotQSCD] = enumerations.CertificateQualification_QCERT_FOR_UNKNOWN

	q[certQualNotQC][certQualESig][certQualNotQSCD] = enumerations.CertificateQualification_CERT_FOR_ESIG
	q[certQualNotQC][certQualESeal][certQualNotQSCD] = enumerations.CertificateQualification_CERT_FOR_ESEAL
	q[certQualNotQC][certQualWSA][certQualNotQSCD] = enumerations.CertificateQualification_CERT_FOR_WSA
	q[certQualNotQC][certQualUnknown][certQualNotQSCD] = enumerations.CertificateQualification_CERT_FOR_UNKNOWN

	q[certQualNotQC][certQualESig][certQualQSCD] = enumerations.CertificateQualification_CERT_FOR_ESIG
	q[certQualNotQC][certQualESeal][certQualQSCD] = enumerations.CertificateQualification_CERT_FOR_ESEAL
	q[certQualNotQC][certQualWSA][certQualQSCD] = enumerations.CertificateQualification_CERT_FOR_WSA
	q[certQualNotQC][certQualUnknown][certQualQSCD] = enumerations.CertificateQualification_CERT_FOR_UNKNOWN

	return q
}()

// GetCertQualification returns the certificate's qualification status based
// on the given parameters. Port of
// getCertQualification(CertificateQualifiedStatus, CertificateType, QSCDStatus).
func GetCertQualification(qc enumerations.CertificateQualifiedStatus, certType enumerations.CertificateType,
	qscd enumerations.QSCDStatus) enumerations.CertificateQualification {
	return certQualifications[certQualInt(enumerations.CertificateQualifiedStatusIsQC(qc))][certTypeInt(certType)][certQualInt(enumerations.QSCDStatusIsQSCD(qscd))]
}

// certTypeInt ports the private static getInt(CertificateType).
func certTypeInt(certType enumerations.CertificateType) int {
	switch certType {
	case enumerations.CertificateType_ESIGN:
		return certQualESig
	case enumerations.CertificateType_ESEAL:
		return certQualESeal
	case enumerations.CertificateType_WSA:
		return certQualWSA
	default:
		return certQualUnknown
	}
}

// certQualInt ports the private static getInt(boolean).
func certQualInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

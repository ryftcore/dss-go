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

	q[certQualQC][certQualESig][certQualQSCD] = enumerations.CertificateQualificationQCERTForESigQSCD
	q[certQualQC][certQualESeal][certQualQSCD] = enumerations.CertificateQualificationQCERTForESealQSCD
	q[certQualQC][certQualWSA][certQualQSCD] = enumerations.CertificateQualificationQCERTForWSA
	q[certQualQC][certQualUnknown][certQualQSCD] = enumerations.CertificateQualificationQCERTForUnknownQSCD

	q[certQualQC][certQualESig][certQualNotQSCD] = enumerations.CertificateQualificationQCERTForESig
	q[certQualQC][certQualESeal][certQualNotQSCD] = enumerations.CertificateQualificationQCERTForESeal
	q[certQualQC][certQualWSA][certQualNotQSCD] = enumerations.CertificateQualificationQCERTForWSA
	q[certQualQC][certQualUnknown][certQualNotQSCD] = enumerations.CertificateQualificationQCERTForUnknown

	q[certQualNotQC][certQualESig][certQualNotQSCD] = enumerations.CertificateQualificationCertForESig
	q[certQualNotQC][certQualESeal][certQualNotQSCD] = enumerations.CertificateQualificationCertForESeal
	q[certQualNotQC][certQualWSA][certQualNotQSCD] = enumerations.CertificateQualificationCertForWSA
	q[certQualNotQC][certQualUnknown][certQualNotQSCD] = enumerations.CertificateQualificationCertForUnknown

	q[certQualNotQC][certQualESig][certQualQSCD] = enumerations.CertificateQualificationCertForESig
	q[certQualNotQC][certQualESeal][certQualQSCD] = enumerations.CertificateQualificationCertForESeal
	q[certQualNotQC][certQualWSA][certQualQSCD] = enumerations.CertificateQualificationCertForWSA
	q[certQualNotQC][certQualUnknown][certQualQSCD] = enumerations.CertificateQualificationCertForUnknown

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
	case enumerations.CertificateTypeESign:
		return certQualESig
	case enumerations.CertificateTypeESeal:
		return certQualESeal
	case enumerations.CertificateTypeWSA:
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

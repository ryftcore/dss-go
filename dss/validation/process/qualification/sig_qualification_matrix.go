// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/SigQualificationMatrix.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// Indexes into the qualifs matrix for the Indication (AdES) dimension.
const (
	sigQualNotAdes           = 0
	sigQualAdes              = 1
	sigQualIndeterminateAdes = 2
)

// Indexes into the qualifs matrix for the CertificateQualification dimension.
const (
	sigQualCertForEsigQscd     = 0
	sigQualCertForEsealQscd    = 1
	sigQualCertForUnknownQscd  = 2
	sigQualCertForEsig         = 3
	sigQualCertForEseal        = 4
	sigQualCertForWsa          = 5
	sigQualCertForUnknown      = 6
	sigQualCertForEsigPlain    = 7
	sigQualCertForEsealPlain   = 8
	sigQualCertForWsaPlain     = 9
	sigQualCertForUnknownPlain = 10
	sigQualNa                  = 11
)

// sigQualifs is the double array containing the relationship between
// qualification parameters and the final signature qualification. Port of
// the static QUALIFS field and its static initializer.
var sigQualifs = func() [3][12]enumerations.SignatureQualification {
	var q [3][12]enumerations.SignatureQualification

	// AdES

	q[sigQualAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualification_QESIG
	q[sigQualAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualification_QESEAL
	q[sigQualAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualification_UNKNOWN_QC_QSCD

	q[sigQualAdes][sigQualCertForEsig] = enumerations.SignatureQualification_ADESIG_QC
	q[sigQualAdes][sigQualCertForEseal] = enumerations.SignatureQualification_ADESEAL_QC
	q[sigQualAdes][sigQualCertForWsa] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualAdes][sigQualCertForUnknown] = enumerations.SignatureQualification_UNKNOWN_QC

	q[sigQualAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualification_ADESIG
	q[sigQualAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualification_ADESEAL
	q[sigQualAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualification_UNKNOWN

	q[sigQualAdes][sigQualNa] = enumerations.SignatureQualification_NA

	// Indeterminate AdES

	q[sigQualIndeterminateAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualification_INDETERMINATE_QESIG
	q[sigQualIndeterminateAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualification_INDETERMINATE_QESEAL
	q[sigQualIndeterminateAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualification_INDETERMINATE_UNKNOWN_QC_QSCD

	q[sigQualIndeterminateAdes][sigQualCertForEsig] = enumerations.SignatureQualification_INDETERMINATE_ADESIG_QC
	q[sigQualIndeterminateAdes][sigQualCertForEseal] = enumerations.SignatureQualification_INDETERMINATE_ADESEAL_QC
	q[sigQualIndeterminateAdes][sigQualCertForWsa] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualIndeterminateAdes][sigQualCertForUnknown] = enumerations.SignatureQualification_INDETERMINATE_UNKNOWN_QC

	q[sigQualIndeterminateAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualification_INDETERMINATE_ADESIG
	q[sigQualIndeterminateAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualification_INDETERMINATE_ADESEAL
	q[sigQualIndeterminateAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualIndeterminateAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualification_INDETERMINATE_UNKNOWN

	q[sigQualIndeterminateAdes][sigQualNa] = enumerations.SignatureQualification_NA

	// Not AdES

	q[sigQualNotAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualification_NOT_ADES_QC_QSCD
	q[sigQualNotAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualification_NOT_ADES_QC_QSCD
	q[sigQualNotAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualification_NOT_ADES_QC_QSCD

	q[sigQualNotAdes][sigQualCertForEsig] = enumerations.SignatureQualification_NOT_ADES_QC
	q[sigQualNotAdes][sigQualCertForEseal] = enumerations.SignatureQualification_NOT_ADES_QC
	q[sigQualNotAdes][sigQualCertForWsa] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualNotAdes][sigQualCertForUnknown] = enumerations.SignatureQualification_NOT_ADES_QC

	q[sigQualNotAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualNotAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualNotAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualification_NOT_ADES
	q[sigQualNotAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualification_NOT_ADES

	q[sigQualNotAdes][sigQualNa] = enumerations.SignatureQualification_NOT_ADES

	return q
}()

// SigQualificationMatrixGetSignatureQualification gets signature
// qualification based on the given parameters. Port of the static
// getSignatureQualification(Indication, CertificateQualification).
func SigQualificationMatrixGetSignatureQualification(ades enumerations.Indication,
	certQualification enumerations.CertificateQualification) enumerations.SignatureQualification {
	return sigQualifs[sigQualificationMatrixGetIndicationInt(ades)][sigQualificationMatrixGetCertQualificationInt(certQualification)]
}

// sigQualificationMatrixGetIndicationInt ports the private static
// getInt(Indication).
func sigQualificationMatrixGetIndicationInt(indication enumerations.Indication) int {
	switch indication {
	case enumerations.Indication_FAILED, enumerations.Indication_TOTAL_FAILED:
		return 0
	case enumerations.Indication_PASSED, enumerations.Indication_TOTAL_PASSED:
		return 1
	case enumerations.Indication_INDETERMINATE:
		return 2
	default:
		panic(fmt.Sprintf("Unsupported indication %s", indication))
	}
}

// sigQualificationMatrixGetCertQualificationInt ports the private static
// getInt(CertificateQualification).
func sigQualificationMatrixGetCertQualificationInt(certQualification enumerations.CertificateQualification) int {
	switch certQualification {
	case enumerations.CertificateQualification_QCERT_FOR_ESIG_QSCD:
		return sigQualCertForEsigQscd
	case enumerations.CertificateQualification_QCERT_FOR_ESEAL_QSCD:
		return sigQualCertForEsealQscd
	case enumerations.CertificateQualification_QCERT_FOR_UNKNOWN_QSCD:
		return sigQualCertForUnknownQscd
	case enumerations.CertificateQualification_QCERT_FOR_ESIG:
		return sigQualCertForEsig
	case enumerations.CertificateQualification_QCERT_FOR_ESEAL:
		return sigQualCertForEseal
	case enumerations.CertificateQualification_QCERT_FOR_WSA:
		return sigQualCertForWsa
	case enumerations.CertificateQualification_QCERT_FOR_UNKNOWN:
		return sigQualCertForUnknown
	case enumerations.CertificateQualification_CERT_FOR_ESIG:
		return sigQualCertForEsigPlain
	case enumerations.CertificateQualification_CERT_FOR_ESEAL:
		return sigQualCertForEsealPlain
	case enumerations.CertificateQualification_CERT_FOR_WSA:
		return sigQualCertForWsaPlain
	case enumerations.CertificateQualification_CERT_FOR_UNKNOWN:
		return sigQualCertForUnknownPlain
	case enumerations.CertificateQualification_NA:
		return sigQualNa
	default:
		panic(fmt.Sprintf("Unsupported certificate qualification %s", certQualification))
	}
}

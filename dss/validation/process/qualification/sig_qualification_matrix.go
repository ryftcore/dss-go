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

	q[sigQualAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualificationQESig
	q[sigQualAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualificationQESeal
	q[sigQualAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualificationUnknownQCQSCD

	q[sigQualAdes][sigQualCertForEsig] = enumerations.SignatureQualificationAdESigQC
	q[sigQualAdes][sigQualCertForEseal] = enumerations.SignatureQualificationAdESealQC
	q[sigQualAdes][sigQualCertForWsa] = enumerations.SignatureQualificationNotAdES
	q[sigQualAdes][sigQualCertForUnknown] = enumerations.SignatureQualificationUnknownQC

	q[sigQualAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualificationAdESig
	q[sigQualAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualificationAdESeal
	q[sigQualAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualificationNotAdES
	q[sigQualAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualificationUnknown

	q[sigQualAdes][sigQualNa] = enumerations.SignatureQualificationNA

	// Indeterminate AdES

	q[sigQualIndeterminateAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualificationIndeterminateQESig
	q[sigQualIndeterminateAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualificationIndeterminateQESeal
	q[sigQualIndeterminateAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualificationIndeterminateUnknownQCQSCD

	q[sigQualIndeterminateAdes][sigQualCertForEsig] = enumerations.SignatureQualificationIndeterminateAdESigQC
	q[sigQualIndeterminateAdes][sigQualCertForEseal] = enumerations.SignatureQualificationIndeterminateAdESealQC
	q[sigQualIndeterminateAdes][sigQualCertForWsa] = enumerations.SignatureQualificationNotAdES
	q[sigQualIndeterminateAdes][sigQualCertForUnknown] = enumerations.SignatureQualificationIndeterminateUnknownQC

	q[sigQualIndeterminateAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualificationIndeterminateAdESig
	q[sigQualIndeterminateAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualificationIndeterminateAdESeal
	q[sigQualIndeterminateAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualificationNotAdES
	q[sigQualIndeterminateAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualificationIndeterminateUnknown

	q[sigQualIndeterminateAdes][sigQualNa] = enumerations.SignatureQualificationNA

	// Not AdES

	q[sigQualNotAdes][sigQualCertForEsigQscd] = enumerations.SignatureQualificationNotAdESQCQSCD
	q[sigQualNotAdes][sigQualCertForEsealQscd] = enumerations.SignatureQualificationNotAdESQCQSCD
	q[sigQualNotAdes][sigQualCertForUnknownQscd] = enumerations.SignatureQualificationNotAdESQCQSCD

	q[sigQualNotAdes][sigQualCertForEsig] = enumerations.SignatureQualificationNotAdESQC
	q[sigQualNotAdes][sigQualCertForEseal] = enumerations.SignatureQualificationNotAdESQC
	q[sigQualNotAdes][sigQualCertForWsa] = enumerations.SignatureQualificationNotAdES
	q[sigQualNotAdes][sigQualCertForUnknown] = enumerations.SignatureQualificationNotAdESQC

	q[sigQualNotAdes][sigQualCertForEsigPlain] = enumerations.SignatureQualificationNotAdES
	q[sigQualNotAdes][sigQualCertForEsealPlain] = enumerations.SignatureQualificationNotAdES
	q[sigQualNotAdes][sigQualCertForWsaPlain] = enumerations.SignatureQualificationNotAdES
	q[sigQualNotAdes][sigQualCertForUnknownPlain] = enumerations.SignatureQualificationNotAdES

	q[sigQualNotAdes][sigQualNa] = enumerations.SignatureQualificationNotAdES

	return q
}()

// SigQualificationMatrixGetSignatureQualification gets signature
// qualification based on the given parameters. Port of the static
// getSignatureQualification(Indication, CertificateQualification).
func SigQualificationMatrixGetSignatureQualification(ades enumerations.Indication,
	certQualification enumerations.CertificateQualification) enumerations.SignatureQualification {
	return sigQualifs[sigQualificationMatrixIndicationInt(ades)][sigQualificationMatrixCertQualificationInt(certQualification)]
}

// sigQualificationMatrixIndicationInt ports the private static
// getInt(Indication).
func sigQualificationMatrixIndicationInt(indication enumerations.Indication) int {
	switch indication {
	case enumerations.IndicationFailed, enumerations.IndicationTotalFailed:
		return 0
	case enumerations.IndicationPassed, enumerations.IndicationTotalPassed:
		return 1
	case enumerations.IndicationIndeterminate:
		return 2
	default:
		panic(fmt.Sprintf("Unsupported indication %s", indication))
	}
}

// sigQualificationMatrixCertQualificationInt ports the private static
// getInt(CertificateQualification).
func sigQualificationMatrixCertQualificationInt(certQualification enumerations.CertificateQualification) int {
	switch certQualification {
	case enumerations.CertificateQualificationQCERTForESigQSCD:
		return sigQualCertForEsigQscd
	case enumerations.CertificateQualificationQCERTForESealQSCD:
		return sigQualCertForEsealQscd
	case enumerations.CertificateQualificationQCERTForUnknownQSCD:
		return sigQualCertForUnknownQscd
	case enumerations.CertificateQualificationQCERTForESig:
		return sigQualCertForEsig
	case enumerations.CertificateQualificationQCERTForESeal:
		return sigQualCertForEseal
	case enumerations.CertificateQualificationQCERTForWSA:
		return sigQualCertForWsa
	case enumerations.CertificateQualificationQCERTForUnknown:
		return sigQualCertForUnknown
	case enumerations.CertificateQualificationCertForESig:
		return sigQualCertForEsigPlain
	case enumerations.CertificateQualificationCertForESeal:
		return sigQualCertForEsealPlain
	case enumerations.CertificateQualificationCertForWSA:
		return sigQualCertForWsaPlain
	case enumerations.CertificateQualificationCertForUnknown:
		return sigQualCertForUnknownPlain
	case enumerations.CertificateQualificationNA:
		return sigQualNa
	default:
		panic(fmt.Sprintf("Unsupported certificate qualification %s", certQualification))
	}
}

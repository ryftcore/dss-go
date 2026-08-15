// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/EAAQualificationMatrix.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
)

// Indexes into the qualifs matrices for the Indication dimension.
const (
	eaaQualPassedEAA        = 0
	eaaQualIndeterminateEAA = 1
	eaaQualFailedEAA        = 2
)

// Indexes into the qualifs matrix for the claimed-EAAQualification dimension.
const (
	eaaQualQEAA       = 0
	eaaQualPUBEAA     = 1
	eaaQualEAA        = 2
	eaaQualUnknownEAA = 3
	eaaQualNotEAA     = 4
)

// Indexes into the qualifs matrix for the SignatureQualification dimension.
const (
	eaaQualQualSigSeal              = 0
	eaaQualIndeterminateQualSigSeal = 1
	eaaQualNotQualSigSeal           = 2
	eaaQualNA                       = 3
)

// Indexes into the PID qualifs matrix for the CertificateApprovalStatus
// dimension.
const (
	eaaQualCertUsagePID   = 0
	eaaQualCertUsageOther = 1
	eaaQualCertUsageNA    = 2
)

// eaaQualifs is the array containing the relationship between qualification
// parameters and the final EAA qualification. Port of the static QUALIFS
// field and its static initializer.
var eaaQualifs = func() [3][5][4]enumerations.EAAQualification {
	var q [3][5][4]enumerations.EAAQualification

	// Passed

	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_QEAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_QEAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_EAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_PUBEAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_PUBEAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_EAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualPassedEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_EAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_EAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_UNKNOWN
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_UNKNOWN
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_UNKNOWN
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualification_NA

	// Indeterminate EAA

	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_QEAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_QEAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_PUBEAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_PUBEAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_EAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_UNKNOWN
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_UNKNOWN
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_INDETERMINATE_UNKNOWN
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualNA] = enumerations.EAAQualification_NA

	// Not EAA

	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualFailedEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualification_NA

	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualNA] = enumerations.EAAQualification_NA

	return q
}()

// eaaPidQualifs is the array containing the relationship between
// qualification parameters and the final PID qualification. Port of the
// static PID_QUALIFS field and its static initializer.
var eaaPidQualifs = func() [3][3]enumerations.EAAQualification {
	var q [3][3]enumerations.EAAQualification

	q[eaaQualPassedEAA][eaaQualCertUsagePID] = enumerations.EAAQualification_PID
	q[eaaQualPassedEAA][eaaQualCertUsageOther] = enumerations.EAAQualification_UNKNOWN
	q[eaaQualPassedEAA][eaaQualCertUsageNA] = enumerations.EAAQualification_NA

	q[eaaQualIndeterminateEAA][eaaQualCertUsagePID] = enumerations.EAAQualification_INDETERMINATE_PID
	q[eaaQualIndeterminateEAA][eaaQualCertUsageOther] = enumerations.EAAQualification_INDETERMINATE_UNKNOWN
	q[eaaQualIndeterminateEAA][eaaQualCertUsageNA] = enumerations.EAAQualification_NA

	q[eaaQualFailedEAA][eaaQualCertUsagePID] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualCertUsageOther] = enumerations.EAAQualification_NOT_EAA
	q[eaaQualFailedEAA][eaaQualCertUsageNA] = enumerations.EAAQualification_NA

	return q
}()

// EAAQualificationMatrixGetEAAQualification gets EAA qualification based on
// the given parameters. Port of the static
// getEAAQualification(Indication, EAAQualification, SignatureQualification).
func EAAQualificationMatrixGetEAAQualification(indication enumerations.Indication, claimedQualification enumerations.EAAQualification,
	signatureQualification enumerations.SignatureQualification) enumerations.EAAQualification {
	return eaaQualifs[eaaQualificationMatrixGetIndicationInt(indication)][eaaQualificationMatrixGetEAAQualificationInt(claimedQualification)][eaaQualificationMatrixGetSignatureQualificationInt(signatureQualification)]
}

// EAAQualificationMatrixGetPIDQualification gets PID qualification based on
// the given parameters. Port of the static
// getPIDQualification(Indication, CertificateApprovalStatus).
func EAAQualificationMatrixGetPIDQualification(indication enumerations.Indication,
	certificateApprovalStatus enumerations.CertificateApprovalStatus) enumerations.EAAQualification {
	return eaaPidQualifs[eaaQualificationMatrixGetIndicationInt(indication)][eaaQualificationMatrixGetCertificateApprovalStatusInt(certificateApprovalStatus)]
}

// eaaQualificationMatrixGetIndicationInt ports the private static
// getInt(Indication).
func eaaQualificationMatrixGetIndicationInt(indication enumerations.Indication) int {
	switch indication {
	case enumerations.Indication_FAILED, enumerations.Indication_TOTAL_FAILED:
		return eaaQualFailedEAA
	case enumerations.Indication_PASSED, enumerations.Indication_TOTAL_PASSED:
		return eaaQualPassedEAA
	case enumerations.Indication_INDETERMINATE:
		return eaaQualIndeterminateEAA
	default:
		panic(fmt.Sprintf("Unsupported indication %s", indication))
	}
}

// eaaQualificationMatrixGetEAAQualificationInt ports the private static
// getInt(EAAQualification).
func eaaQualificationMatrixGetEAAQualificationInt(eaaQualification enumerations.EAAQualification) int {
	switch eaaQualification {
	case enumerations.EAAQualification_QEAA:
		return eaaQualQEAA
	case enumerations.EAAQualification_PUBEAA:
		return eaaQualPUBEAA
	case enumerations.EAAQualification_EAA:
		return eaaQualEAA
	case enumerations.EAAQualification_UNKNOWN:
		return eaaQualUnknownEAA
	case enumerations.EAAQualification_NOT_EAA:
		return eaaQualNotEAA
	default:
		panic(fmt.Sprintf("Unsupported EAA qualification %s", eaaQualification))
	}
}

// eaaQualificationMatrixGetSignatureQualificationInt ports the private
// static getInt(SignatureQualification).
func eaaQualificationMatrixGetSignatureQualificationInt(signatureQualification enumerations.SignatureQualification) int {
	switch signatureQualification {
	case enumerations.SignatureQualification_QESIG, enumerations.SignatureQualification_QESEAL:
		return eaaQualQualSigSeal
	case enumerations.SignatureQualification_INDETERMINATE_QESIG, enumerations.SignatureQualification_INDETERMINATE_QESEAL:
		return eaaQualIndeterminateQualSigSeal
	case enumerations.SignatureQualification_NA:
		return eaaQualNA
	default:
		return eaaQualNotQualSigSeal
	}
}

// eaaQualificationMatrixGetCertificateApprovalStatusInt ports the private
// static getInt(CertificateApprovalStatus).
func eaaQualificationMatrixGetCertificateApprovalStatusInt(certificateApprovalStatus enumerations.CertificateApprovalStatus) int {
	if enumerations.CertificateApprovalStatusEnum_PID_PROVIDER == certificateApprovalStatus {
		return eaaQualCertUsagePID
	} else if enumerations.CertificateApprovalStatusEnum_NA == certificateApprovalStatus {
		return eaaQualCertUsageNA
	}
	return eaaQualCertUsageOther
}

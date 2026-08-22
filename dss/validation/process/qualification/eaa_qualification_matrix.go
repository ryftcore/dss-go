// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/EAAQualificationMatrix.java (DSS 6.5.RC1).
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
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

	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationQEAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateQEAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationEAA
	q[eaaQualPassedEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationPubEAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminatePubEAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationEAA
	q[eaaQualPassedEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualPassedEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationEAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationEAA
	q[eaaQualPassedEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationUnknown
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateUnknown
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationUnknown
	q[eaaQualPassedEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualificationNA

	// Indeterminate EAA

	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationIndeterminateQEAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateQEAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualIndeterminateEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationIndeterminatePubEAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminatePubEAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualIndeterminateEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationIndeterminateEAA
	q[eaaQualIndeterminateEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationIndeterminateUnknown
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationIndeterminateUnknown
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationIndeterminateUnknown
	q[eaaQualIndeterminateEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualIndeterminateEAA][eaaQualNotEAA][eaaQualNA] = enumerations.EAAQualificationNA

	// Not EAA

	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualQEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualPUBEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualFailedEAA][eaaQualEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualUnknownEAA][eaaQualNA] = enumerations.EAAQualificationNA

	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualIndeterminateQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualNotQualSigSeal] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualNotEAA][eaaQualNA] = enumerations.EAAQualificationNA

	return q
}()

// eaaPidQualifs is the array containing the relationship between
// qualification parameters and the final PID qualification. Port of the
// static PID_QUALIFS field and its static initializer.
var eaaPidQualifs = func() [3][3]enumerations.EAAQualification {
	var q [3][3]enumerations.EAAQualification

	q[eaaQualPassedEAA][eaaQualCertUsagePID] = enumerations.EAAQualificationPID
	q[eaaQualPassedEAA][eaaQualCertUsageOther] = enumerations.EAAQualificationUnknown
	q[eaaQualPassedEAA][eaaQualCertUsageNA] = enumerations.EAAQualificationNA

	q[eaaQualIndeterminateEAA][eaaQualCertUsagePID] = enumerations.EAAQualificationIndeterminatePID
	q[eaaQualIndeterminateEAA][eaaQualCertUsageOther] = enumerations.EAAQualificationIndeterminateUnknown
	q[eaaQualIndeterminateEAA][eaaQualCertUsageNA] = enumerations.EAAQualificationNA

	q[eaaQualFailedEAA][eaaQualCertUsagePID] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualCertUsageOther] = enumerations.EAAQualificationNotEAA
	q[eaaQualFailedEAA][eaaQualCertUsageNA] = enumerations.EAAQualificationNA

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
	case enumerations.IndicationFailed, enumerations.IndicationTotalFailed:
		return eaaQualFailedEAA
	case enumerations.IndicationPassed, enumerations.IndicationTotalPassed:
		return eaaQualPassedEAA
	case enumerations.IndicationIndeterminate:
		return eaaQualIndeterminateEAA
	default:
		panic(fmt.Sprintf("Unsupported indication %s", indication))
	}
}

// eaaQualificationMatrixGetEAAQualificationInt ports the private static
// getInt(EAAQualification).
func eaaQualificationMatrixGetEAAQualificationInt(eaaQualification enumerations.EAAQualification) int {
	switch eaaQualification {
	case enumerations.EAAQualificationQEAA:
		return eaaQualQEAA
	case enumerations.EAAQualificationPubEAA:
		return eaaQualPUBEAA
	case enumerations.EAAQualificationEAA:
		return eaaQualEAA
	case enumerations.EAAQualificationUnknown:
		return eaaQualUnknownEAA
	case enumerations.EAAQualificationNotEAA:
		return eaaQualNotEAA
	default:
		panic(fmt.Sprintf("Unsupported EAA qualification %s", eaaQualification))
	}
}

// eaaQualificationMatrixGetSignatureQualificationInt ports the private
// static getInt(SignatureQualification).
func eaaQualificationMatrixGetSignatureQualificationInt(signatureQualification enumerations.SignatureQualification) int {
	switch signatureQualification {
	case enumerations.SignatureQualificationQESig, enumerations.SignatureQualificationQESeal:
		return eaaQualQualSigSeal
	case enumerations.SignatureQualificationIndeterminateQESig, enumerations.SignatureQualificationIndeterminateQESeal:
		return eaaQualIndeterminateQualSigSeal
	case enumerations.SignatureQualificationNA:
		return eaaQualNA
	default:
		return eaaQualNotQualSigSeal
	}
}

// eaaQualificationMatrixGetCertificateApprovalStatusInt ports the private
// static getInt(CertificateApprovalStatus).
func eaaQualificationMatrixGetCertificateApprovalStatusInt(certificateApprovalStatus enumerations.CertificateApprovalStatus) int {
	if enumerations.CertificateApprovalStatusEnumPIDProvider == certificateApprovalStatus {
		return eaaQualCertUsagePID
	} else if enumerations.CertificateApprovalStatusEnumNA == certificateApprovalStatus {
		return eaaQualCertUsageNA
	}
	return eaaQualCertUsageOther
}

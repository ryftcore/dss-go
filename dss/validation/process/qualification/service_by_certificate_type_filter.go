// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/filter/ServiceByCertificateTypeFilter.java (DSS 6.5.RC1).
//
// Allowed services are:
//   - cert type T1 = ASi T1
//   - cert type T1 = ASi T2 + QCForXXX T2 (overrule)
package qualification

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ServiceByCertificateTypeFilter filters TrustServices by certificate type.
type ServiceByCertificateTypeFilter struct {
	AbstractTrustServiceFilter

	// certificate is the certificate to be checked.
	certificate *diagnostic.CertificateWrapper
}

// NewServiceByCertificateTypeFilter is the default constructor. Port of
// ServiceByCertificateTypeFilter(CertificateWrapper).
func NewServiceByCertificateTypeFilter(certificate *diagnostic.CertificateWrapper) *ServiceByCertificateTypeFilter {
	f := &ServiceByCertificateTypeFilter{certificate: certificate}
	f.InitAbstractTrustServiceFilter(f)
	return f
}

// IsAcceptable is the port of the overridden isAcceptable(TrustServiceWrapper).
func (f *ServiceByCertificateTypeFilter) IsAcceptable(service *diagnostic.TrustServiceWrapper) bool {
	issuance := f.certificate.NotBefore()

	if IsPostEIDAS(issuance) {

		additionalServiceInfos := service.AdditionalServiceInfos
		asiEsign := enumerations.AdditionalServiceInformationIsForeSignaturesList(additionalServiceInfos)
		asiEseals := enumerations.AdditionalServiceInformationIsForeSealsList(additionalServiceInfos)
		asiWsa := enumerations.AdditionalServiceInformationIsForWebAuthList(additionalServiceInfos)

		capturedQualifiers := service.CapturedQualifierUris()
		qcForEsign := enumerations.ServiceQualificationIsQcForEsig(capturedQualifiers)
		qcForEseals := enumerations.ServiceQualificationIsQcForEseal(capturedQualifiers)
		qcForWSA := enumerations.ServiceQualificationIsQcForWSA(capturedQualifiers)

		// if QcCompliance and no types -> for eSig by default (see TS 119 615, Table 1)
		qcForEsign = qcForEsign || (!qcForEseals && !qcForWSA && f.certificate.IsQcCompliance())

		count := 0
		for _, b := range []bool{qcForEsign, qcForEseals, qcForWSA} {
			if b {
				count++
			}
		}
		onlyOneQcForXXX := count == 1

		// QCForLegalPerson is not consistent with foreSignature type (see TS 119 615, PRO-4.4.4-12)
		qcForLegalPerson := enumerations.ServiceQualificationIsQcForLegalPerson(capturedQualifiers)
		asiEsign = asiEsign && !qcForLegalPerson

		strategy := CreateTypeFromCert(f.certificate)
		certType := strategy.Type()

		overruleForEsign := asiEsign && qcForEsign && onlyOneQcForXXX
		overruleForEseals := asiEseals && qcForEseals && onlyOneQcForXXX
		overruleForWSA := asiWsa && qcForWSA && onlyOneQcForXXX

		switch certType {
		case enumerations.CertificateType_ESIGN:
			return asiEsign || overruleForEseals || overruleForWSA
		case enumerations.CertificateType_ESEAL:
			return asiEseals || overruleForEsign || overruleForWSA
		case enumerations.CertificateType_WSA:
			return asiWsa || overruleForEseals || overruleForEsign
		case enumerations.CertificateType_UNKNOWN:
			// continue to identify qualification (keeping unknown type)
			return true
		default:
			panic(fmt.Sprintf("Unsupported CertificateType : %s", certType))
		}

	}

	return true
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/CertificatePolicyIdentifiers.java (DSS 6.5.RC1).
package process

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// IsSupportedByQSCD checks if the certificate is supported by QSCD. Port of
// isSupportedByQSCD(CertificateWrapper).
func IsSupportedByQSCD(certificate *diagnostic.CertificateWrapper) bool {
	return hasPolicyIdOIDs(certificate, enumerations.CertificatePolicy_QCP_PUBLIC_WITH_SSCD,
		enumerations.CertificatePolicy_QCP_LEGAL_QSCD, enumerations.CertificatePolicy_QCP_NATURAL_QSCD)
}

// IsQCP checks if the certificate is QCP. Port of isQCP(CertificateWrapper).
func IsQCP(certificate *diagnostic.CertificateWrapper) bool {
	return hasPolicyIdOIDs(certificate, enumerations.CertificatePolicy_QCP_PUBLIC)
}

// IsQCPPlus checks if the certificate is QCP with SSCD. Port of
// isQCPPlus(CertificateWrapper).
func IsQCPPlus(certificate *diagnostic.CertificateWrapper) bool {
	return hasPolicyIdOIDs(certificate, enumerations.CertificatePolicy_QCP_PUBLIC_WITH_SSCD)
}

// IsLegal checks if the certificate is legal QCP. Port of
// isLegal(CertificateWrapper).
func IsLegal(certificate *diagnostic.CertificateWrapper) bool {
	return hasPolicyIdOIDs(certificate, enumerations.CertificatePolicy_QCP_LEGAL,
		enumerations.CertificatePolicy_QCP_LEGAL_QSCD)
}

// IsNatural checks if the certificate is natural QCP. Port of
// isNatural(CertificateWrapper).
func IsNatural(certificate *diagnostic.CertificateWrapper) bool {
	return hasPolicyIdOIDs(certificate, enumerations.CertificatePolicy_QCP_NATURAL,
		enumerations.CertificatePolicy_QCP_NATURAL_QSCD)
}

// hasPolicyIdOIDs ports the private hasPolicyIdOIDs(CertificateWrapper, CertificatePolicy...).
func hasPolicyIdOIDs(certificate *diagnostic.CertificateWrapper, certificatePolicyIds ...enumerations.CertificatePolicy) bool {
	policyIds := certificate.PolicyIds()
	if utils.IsCollectionNotEmpty(policyIds) {
		for _, policyId := range certificatePolicyIds {
			if containsString(policyIds, policyId.OID()) {
				return true
			}
		}
	}
	return false
}

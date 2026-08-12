// Ported from dss-model/.../OidRepository.java (DSS 6.5.RC1).
package model

import "github.com/utain/esig/dss/enumerations"

// oidRepository maps OIDs to their corresponding descriptions.
var oidRepository = buildOidRepository()

func buildOidRepository() map[string]string {
	repo := make(map[string]string)
	oidRepositoryAdd(repo, enumerations.CertificatePolicyValues())
	oidRepositoryAdd(repo, enumerations.QCStatementValues())
	oidRepositoryAdd(repo, enumerations.ExtendedKeyUsageValues())
	oidRepositoryAdd(repo, enumerations.CertificateExtensionEnumValues())
	return repo
}

func oidRepositoryAdd[T enumerations.OidDescription](repo map[string]string, values []T) {
	for _, v := range values {
		repo[v.OID()] = v.Description()
	}
}

// OidRepositoryGetDescription gets the description corresponding to the
// given OID. Ports OidRepository#getDescription (the class is a static
// utility in Java; the Go port drops the private constructor / no-instance
// idiom in favor of a package-level function).
func OidRepositoryGetDescription(oid string) string {
	return oidRepository[oid]
}

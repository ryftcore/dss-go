// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceChecker.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// TrustServiceCheckerIsLegalPersonConsistent checks whether the legal person identifiers
// within the TrustServiceWrapper are consistent. Port of isLegalPersonConsistent(...).
func TrustServiceCheckerIsLegalPersonConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceLegalPersonConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsQCStatementConsistent checks whether the QC statement identifiers
// within the TrustServiceWrapper are consistent. Port of isQCStatementConsistent(...).
func TrustServiceCheckerIsQCStatementConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQCStatementConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsQSCDConsistent checks whether the QSCD identifiers within the
// TrustServiceWrapper are consistent. Port of isQSCDConsistent(...).
func TrustServiceCheckerIsQSCDConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQSCDConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsQSCDStatusAsInCertConsistent checks whether the QSCD as in cert
// identifier within the TrustServiceWrapper is consistent. Port of
// isQSCDStatusAsInCertConsistent(...).
func TrustServiceCheckerIsQSCDStatusAsInCertConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQSCDStatusAsInCertConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsPostEIDASQSCDConsistent checks whether the QSCD identifiers within
// the TrustServiceWrapper are consistent for post eIDAS. Port of
// isPostEIDASQSCDConsistent(...).
func TrustServiceCheckerIsPostEIDASQSCDConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQSCDPostEIDASConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsQualifiersListKnownConsistent checks whether the found qualifiers
// are known by the application. Port of isQualifiersListKnownConsistent(...).
func TrustServiceCheckerIsQualifiersListKnownConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQualifiersKnownConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsUsageConsistent checks whether the usage type identifiers within
// the TrustServiceWrapper are consistent. Port of isUsageConsistent(...).
func TrustServiceCheckerIsUsageConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceUsageConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsPreEIDASStatusConsistent checks whether the statuses before eIDAS
// within the TrustServiceWrapper are consistent. Port of isPreEIDASStatusConsistent(...).
func TrustServiceCheckerIsPreEIDASStatusConsistent(service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceStatusPreEIDASConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent checks whether
// the qualifiers and additional service information before eIDAS within the
// TrustServiceWrapper are consistent. Port of
// isPreEIDASQualifierAndAdditionalServiceInfoConsistent(...).
func TrustServiceCheckerIsPreEIDASQualifierAndAdditionalServiceInfoConsistent(
	service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQualifierAndAdditionalServiceInfoPreEIDASConsistency()
	return condition.IsConsistent(service)
}

// TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent checks whether the
// qualifiers and additional service information are consistent. Port of
// isQualifierAndAdditionalServiceInfoConsistent(...).
func TrustServiceCheckerIsQualifierAndAdditionalServiceInfoConsistent(
	service *diagnostic.TrustServiceWrapper) bool {
	condition := newTrustServiceQualifierAndAdditionalServiceInfoConsistency()
	return condition.IsConsistent(service)
}

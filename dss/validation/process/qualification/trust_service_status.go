// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/TrustServiceStatus.java (DSS 6.5.RC1).
//
// ETSI TS 119 612 V2.2.1.
package qualification

// TrustServiceStatus is the status of a TrustService (ETSI TS 119 612 V2.2.1).
type TrustServiceStatus string

const (
	// TrustServiceStatusUnderSupervision is the before-eIDAS 'undersupervision' status.
	TrustServiceStatusUnderSupervision TrustServiceStatus = "UNDER_SUPERVISION"
	// TrustServiceStatusSupervisionOfServiceInCessation is the before-eIDAS
	// 'supervisionincessation' status.
	TrustServiceStatusSupervisionOfServiceInCessation TrustServiceStatus = "SUPERVISION_OF_SERVICE_IN_CESSATION"
	// TrustServiceStatusSupervisionCeased is the before-eIDAS 'supervisionceased' status.
	TrustServiceStatusSupervisionCeased TrustServiceStatus = "SUPERVISION_CEASED"
	// TrustServiceStatusSupervisionRevoked is the before-eIDAS 'supervisionrevoked' status.
	TrustServiceStatusSupervisionRevoked TrustServiceStatus = "SUPERVISION_REVOKED"
	// TrustServiceStatusAccredited is the before-eIDAS 'accredited' status.
	TrustServiceStatusAccredited TrustServiceStatus = "ACCREDITED"
	// TrustServiceStatusAccreditationCeased is the before-eIDAS 'accreditationceased' status.
	TrustServiceStatusAccreditationCeased TrustServiceStatus = "ACCREDITATION_CEASED"
	// TrustServiceStatusAccreditationRevoked is the before-eIDAS 'accreditationrevoked' status.
	TrustServiceStatusAccreditationRevoked TrustServiceStatus = "ACCREDITATION_REVOKED"

	// TrustServiceStatusGranted is the after-eIDAS 'granted' status.
	TrustServiceStatusGranted TrustServiceStatus = "GRANTED"
	// TrustServiceStatusWithdrawn is the after-eIDAS 'withdrawn' status.
	TrustServiceStatusWithdrawn TrustServiceStatus = "WITHDRAWN"
	// TrustServiceStatusSetByNationalLaw is the after-eIDAS 'setbynationallaw' status.
	TrustServiceStatusSetByNationalLaw TrustServiceStatus = "SET_BY_NATIONAL_LAW"
	// TrustServiceStatusReconizedAtNationalLevel is the after-eIDAS
	// 'recognisedatnationallevel' status.
	TrustServiceStatusReconizedAtNationalLevel TrustServiceStatus = "RECONIZED_AT_NATIONAL_LEVEL"
	// TrustServiceStatusDeprecatedByNationalLaw is the after-eIDAS
	// 'deprecatedbynationallaw' status.
	TrustServiceStatusDeprecatedByNationalLaw TrustServiceStatus = "DEPRECATED_BY_NATIONAL_LAW"
	// TrustServiceStatusDeprecatedAtNationalLevel is the after-eIDAS
	// 'deprecatedatnationallevel' status.
	TrustServiceStatusDeprecatedAtNationalLevel TrustServiceStatus = "DEPRECATED_AT_NATIONAL_LEVEL"
)

// trustServiceStatusFields holds the (shortName, uri, postEidas, valid) tuple for each constant.
type trustServiceStatusFields struct {
	shortName string
	uri       string
	postEidas bool
	valid     bool
}

// trustServiceStatusData holds the fields for each constant, in declaration order.
var trustServiceStatusData = map[TrustServiceStatus]trustServiceStatusFields{
	TrustServiceStatusUnderSupervision: {"under supervision",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/undersupervision", false, true},
	TrustServiceStatusSupervisionOfServiceInCessation: {"supervision in cessation",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionincessation", false, true},
	TrustServiceStatusSupervisionCeased: {"supervision ceased",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionceased", false, false},
	TrustServiceStatusSupervisionRevoked: {"supervision revoked",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionrevoked", false, false},
	TrustServiceStatusAccredited: {"accredited",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accredited", false, true},
	TrustServiceStatusAccreditationCeased: {"accreditation ceased",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accreditationceased", false, false},
	TrustServiceStatusAccreditationRevoked: {"accreditation revoked",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accreditationrevoked", false, false},
	TrustServiceStatusGranted: {"granted",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted", true, true},
	TrustServiceStatusWithdrawn: {"withdrawn",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn", true, false},
	TrustServiceStatusSetByNationalLaw: {"set by national law",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/setbynationallaw", true, false},
	TrustServiceStatusReconizedAtNationalLevel: {"recognised at national level",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/recognisedatnationallevel", true, false},
	TrustServiceStatusDeprecatedByNationalLaw: {"deprecated by national law",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/deprecatedbynationallaw", true, false},
	TrustServiceStatusDeprecatedAtNationalLevel: {"deprecated at national level",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/deprecatedatnationallevel", true, false},
}

// trustServiceStatusValues returns all constants in declaration order.
func trustServiceStatusValues() []TrustServiceStatus {
	return []TrustServiceStatus{
		TrustServiceStatusUnderSupervision,
		TrustServiceStatusSupervisionOfServiceInCessation,
		TrustServiceStatusSupervisionCeased,
		TrustServiceStatusSupervisionRevoked,
		TrustServiceStatusAccredited,
		TrustServiceStatusAccreditationCeased,
		TrustServiceStatusAccreditationRevoked,
		TrustServiceStatusGranted,
		TrustServiceStatusWithdrawn,
		TrustServiceStatusSetByNationalLaw,
		TrustServiceStatusReconizedAtNationalLevel,
		TrustServiceStatusDeprecatedByNationalLaw,
		TrustServiceStatusDeprecatedAtNationalLevel,
	}
}

// ShortName gets the user-friendly label. Port of getShortName().
func (t TrustServiceStatus) ShortName() string {
	return trustServiceStatusData[t].shortName
}

// URI gets the URI. Port of getUri().
func (t TrustServiceStatus) URI() string {
	return trustServiceStatusData[t].uri
}

// IsPreEidas reports whether the status is related to pre-eIDAS. Port of isPreEidas().
func (t TrustServiceStatus) IsPreEidas() bool {
	return !t.IsPostEidas()
}

// IsPostEidas reports whether the status is related to post-eIDAS. Port of isPostEidas().
func (t TrustServiceStatus) IsPostEidas() bool {
	return trustServiceStatusData[t].postEidas
}

// IsValid reports whether the status identifies a valid trust service. Port of isValid().
func (t TrustServiceStatus) IsValid() bool {
	return trustServiceStatusData[t].valid
}

// TrustServiceStatusIsAcceptableStatusBeforeEIDAS gets whether the given status is
// acceptable before eIDAS. Port of isAcceptableStatusBeforeEIDAS(String).
func TrustServiceStatusIsAcceptableStatusBeforeEIDAS(uri string) bool {
	tss := TrustServiceStatusFromUri(uri)
	return tss != "" && tss.IsPreEidas() && tss.IsValid()
}

// TrustServiceStatusIsAcceptableStatusAfterEIDAS gets whether the given status is
// acceptable after eIDAS. Port of isAcceptableStatusAfterEIDAS(String).
func TrustServiceStatusIsAcceptableStatusAfterEIDAS(uri string) bool {
	tss := TrustServiceStatusFromUri(uri)
	return tss != "" && tss.IsPostEidas() && tss.IsValid()
}

// TrustServiceStatusIsSetByNationalLawAfterEIDAS gets whether the given status is set
// by national law after eIDAS. Port of isSetByNationalLawAfterEIDAS(String).
func TrustServiceStatusIsSetByNationalLawAfterEIDAS(uri string) bool {
	tss := TrustServiceStatusFromUri(uri)
	return TrustServiceStatusSetByNationalLaw == tss
}

// TrustServiceStatusIsRecognizedAtNationalLevelAfterEIDAS gets whether the given status
// is recognized at national level after eIDAS. Port of
// isRecognizedAtNationalLevelAfterEIDAS(String).
func TrustServiceStatusIsRecognizedAtNationalLevelAfterEIDAS(uri string) bool {
	tss := TrustServiceStatusFromUri(uri)
	return TrustServiceStatusReconizedAtNationalLevel == tss
}

// TrustServiceStatusFromUri returns a corresponding TrustServiceStatus by the given
// uri, or "" (Java's null) if none matches. Port of fromUri(String).
func TrustServiceStatusFromUri(uri string) TrustServiceStatus {
	for _, status := range trustServiceStatusValues() {
		if trustServiceStatusData[status].uri == uri {
			return status
		}
	}
	return ""
}

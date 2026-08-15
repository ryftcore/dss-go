// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/TrustServiceStatus.java (DSS 6.5.RC1).
//
// ETSI TS 119 612 V2.2.1.
package qualification

// TrustServiceStatus is the status of a TrustService (ETSI TS 119 612 V2.2.1).
type TrustServiceStatus string

const (
	// TrustServiceStatus_UNDER_SUPERVISION is the before-eIDAS 'undersupervision' status.
	TrustServiceStatus_UNDER_SUPERVISION TrustServiceStatus = "UNDER_SUPERVISION"
	// TrustServiceStatus_SUPERVISION_OF_SERVICE_IN_CESSATION is the before-eIDAS
	// 'supervisionincessation' status.
	TrustServiceStatus_SUPERVISION_OF_SERVICE_IN_CESSATION TrustServiceStatus = "SUPERVISION_OF_SERVICE_IN_CESSATION"
	// TrustServiceStatus_SUPERVISION_CEASED is the before-eIDAS 'supervisionceased' status.
	TrustServiceStatus_SUPERVISION_CEASED TrustServiceStatus = "SUPERVISION_CEASED"
	// TrustServiceStatus_SUPERVISION_REVOKED is the before-eIDAS 'supervisionrevoked' status.
	TrustServiceStatus_SUPERVISION_REVOKED TrustServiceStatus = "SUPERVISION_REVOKED"
	// TrustServiceStatus_ACCREDITED is the before-eIDAS 'accredited' status.
	TrustServiceStatus_ACCREDITED TrustServiceStatus = "ACCREDITED"
	// TrustServiceStatus_ACCREDITATION_CEASED is the before-eIDAS 'accreditationceased' status.
	TrustServiceStatus_ACCREDITATION_CEASED TrustServiceStatus = "ACCREDITATION_CEASED"
	// TrustServiceStatus_ACCREDITATION_REVOKED is the before-eIDAS 'accreditationrevoked' status.
	TrustServiceStatus_ACCREDITATION_REVOKED TrustServiceStatus = "ACCREDITATION_REVOKED"

	// TrustServiceStatus_GRANTED is the after-eIDAS 'granted' status.
	TrustServiceStatus_GRANTED TrustServiceStatus = "GRANTED"
	// TrustServiceStatus_WITHDRAWN is the after-eIDAS 'withdrawn' status.
	TrustServiceStatus_WITHDRAWN TrustServiceStatus = "WITHDRAWN"
	// TrustServiceStatus_SET_BY_NATIONAL_LAW is the after-eIDAS 'setbynationallaw' status.
	TrustServiceStatus_SET_BY_NATIONAL_LAW TrustServiceStatus = "SET_BY_NATIONAL_LAW"
	// TrustServiceStatus_RECONIZED_AT_NATIONAL_LEVEL is the after-eIDAS
	// 'recognisedatnationallevel' status.
	TrustServiceStatus_RECONIZED_AT_NATIONAL_LEVEL TrustServiceStatus = "RECONIZED_AT_NATIONAL_LEVEL"
	// TrustServiceStatus_DEPRECATED_BY_NATIONAL_LAW is the after-eIDAS
	// 'deprecatedbynationallaw' status.
	TrustServiceStatus_DEPRECATED_BY_NATIONAL_LAW TrustServiceStatus = "DEPRECATED_BY_NATIONAL_LAW"
	// TrustServiceStatus_DEPRECATED_AT_NATIONAL_LEVEL is the after-eIDAS
	// 'deprecatedatnationallevel' status.
	TrustServiceStatus_DEPRECATED_AT_NATIONAL_LEVEL TrustServiceStatus = "DEPRECATED_AT_NATIONAL_LEVEL"
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
	TrustServiceStatus_UNDER_SUPERVISION: {"under supervision",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/undersupervision", false, true},
	TrustServiceStatus_SUPERVISION_OF_SERVICE_IN_CESSATION: {"supervision in cessation",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionincessation", false, true},
	TrustServiceStatus_SUPERVISION_CEASED: {"supervision ceased",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionceased", false, false},
	TrustServiceStatus_SUPERVISION_REVOKED: {"supervision revoked",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/supervisionrevoked", false, false},
	TrustServiceStatus_ACCREDITED: {"accredited",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accredited", false, true},
	TrustServiceStatus_ACCREDITATION_CEASED: {"accreditation ceased",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accreditationceased", false, false},
	TrustServiceStatus_ACCREDITATION_REVOKED: {"accreditation revoked",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/accreditationrevoked", false, false},
	TrustServiceStatus_GRANTED: {"granted",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted", true, true},
	TrustServiceStatus_WITHDRAWN: {"withdrawn",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/withdrawn", true, false},
	TrustServiceStatus_SET_BY_NATIONAL_LAW: {"set by national law",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/setbynationallaw", true, false},
	TrustServiceStatus_RECONIZED_AT_NATIONAL_LEVEL: {"recognised at national level",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/recognisedatnationallevel", true, false},
	TrustServiceStatus_DEPRECATED_BY_NATIONAL_LAW: {"deprecated by national law",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/deprecatedbynationallaw", true, false},
	TrustServiceStatus_DEPRECATED_AT_NATIONAL_LEVEL: {"deprecated at national level",
		"http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/deprecatedatnationallevel", true, false},
}

// trustServiceStatusValues returns all constants in declaration order.
func trustServiceStatusValues() []TrustServiceStatus {
	return []TrustServiceStatus{
		TrustServiceStatus_UNDER_SUPERVISION,
		TrustServiceStatus_SUPERVISION_OF_SERVICE_IN_CESSATION,
		TrustServiceStatus_SUPERVISION_CEASED,
		TrustServiceStatus_SUPERVISION_REVOKED,
		TrustServiceStatus_ACCREDITED,
		TrustServiceStatus_ACCREDITATION_CEASED,
		TrustServiceStatus_ACCREDITATION_REVOKED,
		TrustServiceStatus_GRANTED,
		TrustServiceStatus_WITHDRAWN,
		TrustServiceStatus_SET_BY_NATIONAL_LAW,
		TrustServiceStatus_RECONIZED_AT_NATIONAL_LEVEL,
		TrustServiceStatus_DEPRECATED_BY_NATIONAL_LAW,
		TrustServiceStatus_DEPRECATED_AT_NATIONAL_LEVEL,
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
	return TrustServiceStatus_SET_BY_NATIONAL_LAW == tss
}

// TrustServiceStatusIsRecognizedAtNationalLevelAfterEIDAS gets whether the given status
// is recognized at national level after eIDAS. Port of
// isRecognizedAtNationalLevelAfterEIDAS(String).
func TrustServiceStatusIsRecognizedAtNationalLevelAfterEIDAS(uri string) bool {
	tss := TrustServiceStatusFromUri(uri)
	return TrustServiceStatus_RECONIZED_AT_NATIONAL_LEVEL == tss
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

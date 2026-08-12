// Ported from dss-enumerations/.../RevocationReason.java (DSS 6.5.RC1).
//
// This enum is used to get the String value of CRLReason.
//
// The CRLReason enumeration:
//
//	CRLReason ::= ENUMERATED {
//	 unspecified             (0),
//	 keyCompromise           (1),
//	 cACompromise            (2),
//	 affiliationChanged      (3),
//	 superseded              (4),
//	 cessationOfOperation    (5),
//	 certificateHold         (6),
//	 removeFromCRL           (8),
//	 privilegeWithdrawn      (9),
//	 aACompromise           (10)
//	}
package enumerations

// RevocationReason represents a revocation reason. Implements UriBasedEnum.
type RevocationReason string

const (
	// RevocationReason_UNSPECIFIED is unspecified.
	RevocationReason_UNSPECIFIED RevocationReason = "UNSPECIFIED"
	// RevocationReason_KEY_COMPROMISE is keyCompromise.
	RevocationReason_KEY_COMPROMISE RevocationReason = "KEY_COMPROMISE"
	// RevocationReason_CA_COMPROMISE is cACompromise.
	RevocationReason_CA_COMPROMISE RevocationReason = "CA_COMPROMISE"
	// RevocationReason_AFFILIATION_CHANGED is affiliationChanged.
	RevocationReason_AFFILIATION_CHANGED RevocationReason = "AFFILIATION_CHANGED"
	// RevocationReason_SUPERSEDED is superseded.
	RevocationReason_SUPERSEDED RevocationReason = "SUPERSEDED"
	// RevocationReason_CESSATION_OF_OPERATION is cessationOfOperation.
	RevocationReason_CESSATION_OF_OPERATION RevocationReason = "CESSATION_OF_OPERATION"
	// RevocationReason_CERTIFICATE_HOLD is certificateHold.
	RevocationReason_CERTIFICATE_HOLD RevocationReason = "CERTIFICATE_HOLD"
	// RevocationReason_REMOVE_FROM_CRL is removeFromCRL. Missing in ETSI VR
	// standard.
	RevocationReason_REMOVE_FROM_CRL RevocationReason = "REMOVE_FROM_CRL"
	// RevocationReason_PRIVILEGE_WITHDRAWN is privilegeWithdrawn.
	RevocationReason_PRIVILEGE_WITHDRAWN RevocationReason = "PRIVILEGE_WITHDRAWN"
	// RevocationReason_AA_COMPROMISE is aACompromise. Missing in ETSI VI
	// standard.
	RevocationReason_AA_COMPROMISE RevocationReason = "AA_COMPROMISE"
)

type revocationReasonFields struct {
	shortName string
	uri       string
	value     int
}

// revocationReasonData holds the (shortName, uri, value) tuple for each constant.
var revocationReasonData = map[RevocationReason]revocationReasonFields{
	RevocationReason_UNSPECIFIED:            {"unspecified", "urn:etsi:019102:revocationReason:unspecified", 0},
	RevocationReason_KEY_COMPROMISE:         {"keyCompromise", "urn:etsi:019102:revocationReason:keyCompromise", 1},
	RevocationReason_CA_COMPROMISE:          {"cACompromise", "urn:etsi:019102:revocationReason:cACompromise", 2},
	RevocationReason_AFFILIATION_CHANGED:    {"affiliationChanged", "urn:etsi:019102:revocationReason:affiliationChanged", 3},
	RevocationReason_SUPERSEDED:             {"superseded", "urn:etsi:019102:revocationReason:superseded", 4},
	RevocationReason_CESSATION_OF_OPERATION: {"cessationOfOperation", "urn:etsi:019102:revocationReason:cessationOfOperation", 5},
	RevocationReason_CERTIFICATE_HOLD:       {"certificateHold", "urn:etsi:019102:revocationReason:certificateHold", 6},
	RevocationReason_REMOVE_FROM_CRL:        {"removeFromCRL", "urn:etsi:019102:revocationReason:removeFromCRL", 8},
	RevocationReason_PRIVILEGE_WITHDRAWN:    {"privilegeWithdrawn", "urn:etsi:019102:revocationReason:privilegeWithdrawn", 9},
	RevocationReason_AA_COMPROMISE:          {"aACompromise", "urn:etsi:019102:revocationReason:aACompromise", 10},
}

// RevocationReasonValues returns all constants in declaration order.
func RevocationReasonValues() []RevocationReason {
	return []RevocationReason{
		RevocationReason_UNSPECIFIED,
		RevocationReason_KEY_COMPROMISE,
		RevocationReason_CA_COMPROMISE,
		RevocationReason_AFFILIATION_CHANGED,
		RevocationReason_SUPERSEDED,
		RevocationReason_CESSATION_OF_OPERATION,
		RevocationReason_CERTIFICATE_HOLD,
		RevocationReason_REMOVE_FROM_CRL,
		RevocationReason_PRIVILEGE_WITHDRAWN,
		RevocationReason_AA_COMPROMISE,
	}
}

// ShortName returns the name of the RevocationReason.
func (r RevocationReason) ShortName() string {
	return revocationReasonData[r].shortName
}

// URI returns the URI within VR. Implements UriBasedEnum.
func (r RevocationReason) URI() string {
	return revocationReasonData[r].uri
}

// Value returns the value of the RevocationReason.
func (r RevocationReason) Value() int {
	return revocationReasonData[r].value
}

// RevocationReasonFromInt returns a RevocationReason based on the given
// integer value. Returns "" (the zero value) if not found, mirroring
// Java's null return.
func RevocationReasonFromInt(value int) RevocationReason {
	for _, v := range RevocationReasonValues() {
		if revocationReasonData[v].value == value {
			return v
		}
	}
	return ""
}

// RevocationReasonFromValue returns a RevocationReason based on the given
// label string value. Returns "" (the zero value) if not found, mirroring
// Java's null return.
func RevocationReasonFromValue(value string) RevocationReason {
	for _, v := range RevocationReasonValues() {
		if revocationReasonData[v].shortName == value {
			return v
		}
	}
	return ""
}

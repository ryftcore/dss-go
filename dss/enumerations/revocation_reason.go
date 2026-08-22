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
	// RevocationReasonUnspecified is unspecified.
	RevocationReasonUnspecified RevocationReason = "UNSPECIFIED"
	// RevocationReasonKeyCompromise is keyCompromise.
	RevocationReasonKeyCompromise RevocationReason = "KEY_COMPROMISE"
	// RevocationReasonCACompromise is cACompromise.
	RevocationReasonCACompromise RevocationReason = "CA_COMPROMISE"
	// RevocationReasonAffiliationChanged is affiliationChanged.
	RevocationReasonAffiliationChanged RevocationReason = "AFFILIATION_CHANGED"
	// RevocationReasonSuperseded is superseded.
	RevocationReasonSuperseded RevocationReason = "SUPERSEDED"
	// RevocationReasonCessationOfOperation is cessationOfOperation.
	RevocationReasonCessationOfOperation RevocationReason = "CESSATION_OF_OPERATION"
	// RevocationReasonCertificateHold is certificateHold.
	RevocationReasonCertificateHold RevocationReason = "CERTIFICATE_HOLD"
	// RevocationReasonRemoveFromCRL is removeFromCRL. Missing in ETSI VR
	// standard.
	RevocationReasonRemoveFromCRL RevocationReason = "REMOVE_FROM_CRL"
	// RevocationReasonPrivilegeWithdrawn is privilegeWithdrawn.
	RevocationReasonPrivilegeWithdrawn RevocationReason = "PRIVILEGE_WITHDRAWN"
	// RevocationReasonAACompromise is aACompromise. Missing in ETSI VI
	// standard.
	RevocationReasonAACompromise RevocationReason = "AA_COMPROMISE"
)

type revocationReasonFields struct {
	shortName string
	uri       string
	value     int
}

// revocationReasonData holds the (shortName, uri, value) tuple for each constant.
var revocationReasonData = map[RevocationReason]revocationReasonFields{
	RevocationReasonUnspecified:          {"unspecified", "urn:etsi:019102:revocationReason:unspecified", 0},
	RevocationReasonKeyCompromise:        {"keyCompromise", "urn:etsi:019102:revocationReason:keyCompromise", 1},
	RevocationReasonCACompromise:         {"cACompromise", "urn:etsi:019102:revocationReason:cACompromise", 2},
	RevocationReasonAffiliationChanged:   {"affiliationChanged", "urn:etsi:019102:revocationReason:affiliationChanged", 3},
	RevocationReasonSuperseded:           {"superseded", "urn:etsi:019102:revocationReason:superseded", 4},
	RevocationReasonCessationOfOperation: {"cessationOfOperation", "urn:etsi:019102:revocationReason:cessationOfOperation", 5},
	RevocationReasonCertificateHold:      {"certificateHold", "urn:etsi:019102:revocationReason:certificateHold", 6},
	RevocationReasonRemoveFromCRL:        {"removeFromCRL", "urn:etsi:019102:revocationReason:removeFromCRL", 8},
	RevocationReasonPrivilegeWithdrawn:   {"privilegeWithdrawn", "urn:etsi:019102:revocationReason:privilegeWithdrawn", 9},
	RevocationReasonAACompromise:         {"aACompromise", "urn:etsi:019102:revocationReason:aACompromise", 10},
}

// RevocationReasonValues returns all constants in declaration order.
func RevocationReasonValues() []RevocationReason {
	return []RevocationReason{
		RevocationReasonUnspecified,
		RevocationReasonKeyCompromise,
		RevocationReasonCACompromise,
		RevocationReasonAffiliationChanged,
		RevocationReasonSuperseded,
		RevocationReasonCessationOfOperation,
		RevocationReasonCertificateHold,
		RevocationReasonRemoveFromCRL,
		RevocationReasonPrivilegeWithdrawn,
		RevocationReasonAACompromise,
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

// Ported from dss-enumerations/.../CommitmentTypeEnum.java (DSS 6.5.RC1).
package enumerations

// CommitmentTypeEnum is defined in ETSI TS 119 172-1 Annex B.
type CommitmentTypeEnum string

const (
	// CommitmentTypeEnum_ProofOfOrigin indicates that the signer recognizes
	// to have created, approved and sent the signed data.
	CommitmentTypeEnum_ProofOfOrigin CommitmentTypeEnum = "ProofOfOrigin"
	// CommitmentTypeEnum_ProofOfReceipt indicates that signer recognizes to
	// have received the content of the signed data.
	CommitmentTypeEnum_ProofOfReceipt CommitmentTypeEnum = "ProofOfReceipt"
	// CommitmentTypeEnum_ProofOfDelivery indicates that the TSP providing
	// that indication has delivered a signed data in a local store
	// accessible to the recipient of the signed data.
	CommitmentTypeEnum_ProofOfDelivery CommitmentTypeEnum = "ProofOfDelivery"
	// CommitmentTypeEnum_ProofOfSender indicates that the entity providing
	// that indication has sent the signed data (but not necessarily created
	// it).
	CommitmentTypeEnum_ProofOfSender CommitmentTypeEnum = "ProofOfSender"
	// CommitmentTypeEnum_ProofOfApproval indicates that the signer has
	// approved the content of the signed data.
	CommitmentTypeEnum_ProofOfApproval CommitmentTypeEnum = "ProofOfApproval"
	// CommitmentTypeEnum_ProofOfCreation indicates that the signer has
	// created the signed data (but not necessarily approved, nor sent it).
	CommitmentTypeEnum_ProofOfCreation CommitmentTypeEnum = "ProofOfCreation"
)

type commitmentTypeEnumFields struct {
	uri string
	oid string
}

// commitmentTypeEnumData holds the (uri, oid) tuple for each constant.
// NOTE: the qualifier and documentationReferences fields present in the
// Java constructor are always null/empty for every constant in this enum
// (see the Java class comment); they are therefore ported as fixed
// zero-value returns on Qualifier() and DocumentationReferences() rather
// than per-constant table entries.
var commitmentTypeEnumData = map[CommitmentTypeEnum]commitmentTypeEnumFields{
	CommitmentTypeEnum_ProofOfOrigin:   {"http://uri.etsi.org/01903/v1.2.2#ProofOfOrigin", "1.2.840.113549.1.9.16.6.1"},
	CommitmentTypeEnum_ProofOfReceipt:  {"http://uri.etsi.org/01903/v1.2.2#ProofOfReceipt", "1.2.840.113549.1.9.16.6.2"},
	CommitmentTypeEnum_ProofOfDelivery: {"http://uri.etsi.org/01903/v1.2.2#ProofOfDelivery", "1.2.840.113549.1.9.16.6.3"},
	CommitmentTypeEnum_ProofOfSender:   {"http://uri.etsi.org/01903/v1.2.2#ProofOfSender", "1.2.840.113549.1.9.16.6.4"},
	CommitmentTypeEnum_ProofOfApproval: {"http://uri.etsi.org/01903/v1.2.2#ProofOfApproval", "1.2.840.113549.1.9.16.6.5"},
	CommitmentTypeEnum_ProofOfCreation: {"http://uri.etsi.org/01903/v1.2.2#ProofOfCreation", "1.2.840.113549.1.9.16.6.6"},
}

// CommitmentTypeEnumValues returns all constants in declaration order.
func CommitmentTypeEnumValues() []CommitmentTypeEnum {
	return []CommitmentTypeEnum{
		CommitmentTypeEnum_ProofOfOrigin,
		CommitmentTypeEnum_ProofOfReceipt,
		CommitmentTypeEnum_ProofOfDelivery,
		CommitmentTypeEnum_ProofOfSender,
		CommitmentTypeEnum_ProofOfApproval,
		CommitmentTypeEnum_ProofOfCreation,
	}
}

// URI returns the XML URI (XAdES).
func (c CommitmentTypeEnum) URI() string {
	return commitmentTypeEnumData[c].uri
}

// OID returns the Object Identifier (CAdES).
func (c CommitmentTypeEnum) OID() string {
	return commitmentTypeEnumData[c].oid
}

// Qualifier returns the Object Identifier Qualifier. Always "" (zero value)
// for this enum: the qualifier and documentationReferences are not used;
// to use them, overwrite the CommitmentType interface (see Java Javadoc).
func (c CommitmentTypeEnum) Qualifier() ObjectIdentifierQualifier {
	return ""
}

// Description returns the enum constant name.
func (c CommitmentTypeEnum) Description() string {
	return string(c)
}

// DocumentationReferences returns nil (always empty for this enum; see
// Qualifier's note).
func (c CommitmentTypeEnum) DocumentationReferences() []string {
	return nil
}

// CommitmentTypeEnumValueOf returns the constant matching the given Java
// enum name.
func CommitmentTypeEnumValueOf(name string) (CommitmentTypeEnum, error) {
	for _, v := range CommitmentTypeEnumValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &commitmentTypeEnumInvalidValueError{name}
}

type commitmentTypeEnumInvalidValueError struct {
	name string
}

func (e *commitmentTypeEnumInvalidValueError) Error() string {
	return "no enum constant CommitmentTypeEnum." + e.name
}

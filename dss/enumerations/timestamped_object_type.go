// Ported from dss-enumerations/.../TimestampedObjectType.java (DSS 6.5.RC1).
package enumerations

// TimestampedObjectType defines possible object types to be timestamped.
type TimestampedObjectType string

const (
	// TimestampedObjectType_SIGNED_DATA is the original document (signed
	// data).
	TimestampedObjectType_SIGNED_DATA TimestampedObjectType = "SIGNED_DATA"
	// TimestampedObjectType_SIGNATURE is a signature.
	TimestampedObjectType_SIGNATURE TimestampedObjectType = "SIGNATURE"
	// TimestampedObjectType_CERTIFICATE is a certificate.
	TimestampedObjectType_CERTIFICATE TimestampedObjectType = "CERTIFICATE"
	// TimestampedObjectType_REVOCATION is revocation data.
	TimestampedObjectType_REVOCATION TimestampedObjectType = "REVOCATION"
	// TimestampedObjectType_TIMESTAMP is a timestamp.
	TimestampedObjectType_TIMESTAMP TimestampedObjectType = "TIMESTAMP"
	// TimestampedObjectType_EVIDENCE_RECORD is an evidence record.
	TimestampedObjectType_EVIDENCE_RECORD TimestampedObjectType = "EVIDENCE_RECORD"
	// TimestampedObjectType_ORPHAN_CERTIFICATE is a not used certificate
	// token.
	TimestampedObjectType_ORPHAN_CERTIFICATE TimestampedObjectType = "ORPHAN_CERTIFICATE"
	// TimestampedObjectType_ORPHAN_REVOCATION is a not used revocation
	// token.
	TimestampedObjectType_ORPHAN_REVOCATION TimestampedObjectType = "ORPHAN_REVOCATION"
)

// TimestampedObjectTypeValues returns all constants in declaration order.
func TimestampedObjectTypeValues() []TimestampedObjectType {
	return []TimestampedObjectType{
		TimestampedObjectType_SIGNED_DATA,
		TimestampedObjectType_SIGNATURE,
		TimestampedObjectType_CERTIFICATE,
		TimestampedObjectType_REVOCATION,
		TimestampedObjectType_TIMESTAMP,
		TimestampedObjectType_EVIDENCE_RECORD,
		TimestampedObjectType_ORPHAN_CERTIFICATE,
		TimestampedObjectType_ORPHAN_REVOCATION,
	}
}

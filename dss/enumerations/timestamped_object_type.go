// Ported from dss-enumerations/.../TimestampedObjectType.java (DSS 6.5.RC1).
package enumerations

// TimestampedObjectType defines possible object types to be timestamped.
type TimestampedObjectType string

const (
	// TimestampedObjectTypeSignedData is the original document (signed
	// data).
	TimestampedObjectTypeSignedData TimestampedObjectType = "SIGNED_DATA"
	// TimestampedObjectTypeSignature is a signature.
	TimestampedObjectTypeSignature TimestampedObjectType = "SIGNATURE"
	// TimestampedObjectTypeCertificate is a certificate.
	TimestampedObjectTypeCertificate TimestampedObjectType = "CERTIFICATE"
	// TimestampedObjectTypeRevocation is revocation data.
	TimestampedObjectTypeRevocation TimestampedObjectType = "REVOCATION"
	// TimestampedObjectTypeTimestamp is a timestamp.
	TimestampedObjectTypeTimestamp TimestampedObjectType = "TIMESTAMP"
	// TimestampedObjectTypeEvidenceRecord is an evidence record.
	TimestampedObjectTypeEvidenceRecord TimestampedObjectType = "EVIDENCE_RECORD"
	// TimestampedObjectTypeOrphanCertificate is a not used certificate
	// token.
	TimestampedObjectTypeOrphanCertificate TimestampedObjectType = "ORPHAN_CERTIFICATE"
	// TimestampedObjectTypeOrphanRevocation is a not used revocation
	// token.
	TimestampedObjectTypeOrphanRevocation TimestampedObjectType = "ORPHAN_REVOCATION"
)

// TimestampedObjectTypeValues returns all constants in declaration order.
func TimestampedObjectTypeValues() []TimestampedObjectType {
	return []TimestampedObjectType{
		TimestampedObjectTypeSignedData,
		TimestampedObjectTypeSignature,
		TimestampedObjectTypeCertificate,
		TimestampedObjectTypeRevocation,
		TimestampedObjectTypeTimestamp,
		TimestampedObjectTypeEvidenceRecord,
		TimestampedObjectTypeOrphanCertificate,
		TimestampedObjectTypeOrphanRevocation,
	}
}

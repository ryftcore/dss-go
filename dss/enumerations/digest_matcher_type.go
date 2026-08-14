// Ported from dss-enumerations/.../DigestMatcherType.java (DSS 6.5.RC1).
package enumerations

// DigestMatcherType defines available types of DigestMatchers (signed data
// origins).
type DigestMatcherType string

const (
	// DigestMatcherType_REFERENCE: XAdES signed reference.
	DigestMatcherType_REFERENCE DigestMatcherType = "REFERENCE"
	// DigestMatcherType_OBJECT: XAdES signed reference of Object type.
	DigestMatcherType_OBJECT DigestMatcherType = "OBJECT"
	// DigestMatcherType_MANIFEST: XAdES signed manifest.
	DigestMatcherType_MANIFEST DigestMatcherType = "MANIFEST"
	// DigestMatcherType_SIGNED_PROPERTIES: XAdES SignedProperties element.
	DigestMatcherType_SIGNED_PROPERTIES DigestMatcherType = "SIGNED_PROPERTIES"
	// DigestMatcherType_KEY_INFO: XAdES KeyInfo element.
	DigestMatcherType_KEY_INFO DigestMatcherType = "KEY_INFO"
	// DigestMatcherType_SIGNATURE_PROPERTIES: XAdES SignatureProperties
	// element.
	DigestMatcherType_SIGNATURE_PROPERTIES DigestMatcherType = "SIGNATURE_PROPERTIES"
	// DigestMatcherType_XPOINTER: XAdES XPointer reference.
	DigestMatcherType_XPOINTER DigestMatcherType = "XPOINTER"
	// DigestMatcherType_MANIFEST_ENTRY: XAdES and ASiC CAdES.
	DigestMatcherType_MANIFEST_ENTRY DigestMatcherType = "MANIFEST_ENTRY"
	// DigestMatcherType_COUNTER_SIGNATURE: XAdES signed SignatureValue
	// (counter signature).
	DigestMatcherType_COUNTER_SIGNATURE DigestMatcherType = "COUNTER_SIGNATURE"
	// DigestMatcherType_MESSAGE_DIGEST: CAdES.
	DigestMatcherType_MESSAGE_DIGEST DigestMatcherType = "MESSAGE_DIGEST"
	// DigestMatcherType_CONTENT_DIGEST: digest from decrypted content
	// SignatureValue (CAdES/PAdES).
	DigestMatcherType_CONTENT_DIGEST DigestMatcherType = "CONTENT_DIGEST"
	// DigestMatcherType_JWS_SIGNING_INPUT: JAdES Digest on result of
	// concatenation
	// ASCII(BASE64URL(UTF8(JWSProtected Header)) || '.' ||
	// BASE64URL(JWS Payload)).
	DigestMatcherType_JWS_SIGNING_INPUT DigestMatcherType = "JWS_SIGNING_INPUT"
	// DigestMatcherType_SIG_D_ENTRY: JAdES or CB-AdES Detached entry.
	DigestMatcherType_SIG_D_ENTRY DigestMatcherType = "SIG_D_ENTRY"
	// DigestMatcherType_COSE_SIG_STRUCTURE: COSE Digest on result of
	// serialization of Sig_structure array.
	DigestMatcherType_COSE_SIG_STRUCTURE DigestMatcherType = "COSE_SIG_STRUCTURE"
	// DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE: defines the
	// signature value of a master signature signed by a counter signature.
	DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE DigestMatcherType = "COUNTER_SIGNED_SIGNATURE_VALUE"
	// DigestMatcherType_MESSAGE_IMPRINT: timestamp.
	DigestMatcherType_MESSAGE_IMPRINT DigestMatcherType = "MESSAGE_IMPRINT"
	// DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT: evidence record
	// archive object.
	DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_OBJECT"
	// DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE: identifies
	// evidence record archive object which has not been associated with any
	// of the provided documents.
	DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE DigestMatcherType = "EVIDENCE_RECORD_ORPHAN_REFERENCE"
	// DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP: evidence record
	// previous archive time-stamp object.
	DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_TIME_STAMP"
	// DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE:
	// evidence record previous archive time-stamp sequence.
	DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE"
	// DigestMatcherType_EVIDENCE_RECORD_MASTER_SIGNATURE: evidence record
	// embedded in a signature.
	DigestMatcherType_EVIDENCE_RECORD_MASTER_SIGNATURE DigestMatcherType = "EVIDENCE_RECORD_MASTER_SIGNATURE"
	// DigestMatcherType_EAA_DISCLOSURE: disclosure attached to a
	// presentation of EAA.
	DigestMatcherType_EAA_DISCLOSURE DigestMatcherType = "EAA_DISCLOSURE"
	// DigestMatcherType_EAA_NESTED_DISCLOSURE: disclosure nested to
	// provided disclosure to a presentation of EAA.
	DigestMatcherType_EAA_NESTED_DISCLOSURE DigestMatcherType = "EAA_NESTED_DISCLOSURE"
	// DigestMatcherType_EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM:
	// incorporated SD claim for which no matching provided disclosure has
	// been found.
	DigestMatcherType_EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM DigestMatcherType = "EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM"
	// DigestMatcherType_EAA_KEY_BINDING: input used to compute a key
	// binding signature (used in EAA).
	DigestMatcherType_EAA_KEY_BINDING DigestMatcherType = "EAA_KEY_BINDING"
)

// DigestMatcherTypeValues returns all constants in declaration order.
func DigestMatcherTypeValues() []DigestMatcherType {
	return []DigestMatcherType{
		DigestMatcherType_REFERENCE,
		DigestMatcherType_OBJECT,
		DigestMatcherType_MANIFEST,
		DigestMatcherType_SIGNED_PROPERTIES,
		DigestMatcherType_KEY_INFO,
		DigestMatcherType_SIGNATURE_PROPERTIES,
		DigestMatcherType_XPOINTER,
		DigestMatcherType_MANIFEST_ENTRY,
		DigestMatcherType_COUNTER_SIGNATURE,
		DigestMatcherType_MESSAGE_DIGEST,
		DigestMatcherType_CONTENT_DIGEST,
		DigestMatcherType_JWS_SIGNING_INPUT,
		DigestMatcherType_SIG_D_ENTRY,
		DigestMatcherType_COSE_SIG_STRUCTURE,
		DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE,
		DigestMatcherType_MESSAGE_IMPRINT,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT,
		DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE,
		DigestMatcherType_EVIDENCE_RECORD_MASTER_SIGNATURE,
		DigestMatcherType_EAA_DISCLOSURE,
		DigestMatcherType_EAA_NESTED_DISCLOSURE,
		DigestMatcherType_EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM,
		DigestMatcherType_EAA_KEY_BINDING,
	}
}

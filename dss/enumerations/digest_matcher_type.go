// Ported from dss-enumerations/.../DigestMatcherType.java (DSS 6.5.RC1).
package enumerations

// DigestMatcherType defines available types of DigestMatchers (signed data
// origins).
type DigestMatcherType string

const (
	// DigestMatcherTypeReference: XAdES signed reference.
	DigestMatcherTypeReference DigestMatcherType = "REFERENCE"
	// DigestMatcherTypeObject: XAdES signed reference of Object type.
	DigestMatcherTypeObject DigestMatcherType = "OBJECT"
	// DigestMatcherTypeManifest: XAdES signed manifest.
	DigestMatcherTypeManifest DigestMatcherType = "MANIFEST"
	// DigestMatcherTypeSignedProperties: XAdES SignedProperties element.
	DigestMatcherTypeSignedProperties DigestMatcherType = "SIGNED_PROPERTIES"
	// DigestMatcherTypeKeyInfo: XAdES KeyInfo element.
	DigestMatcherTypeKeyInfo DigestMatcherType = "KEY_INFO"
	// DigestMatcherTypeSignatureProperties: XAdES SignatureProperties
	// element.
	DigestMatcherTypeSignatureProperties DigestMatcherType = "SIGNATURE_PROPERTIES"
	// DigestMatcherTypeXPointer: XAdES XPointer reference.
	DigestMatcherTypeXPointer DigestMatcherType = "XPOINTER"
	// DigestMatcherTypeManifestEntry: XAdES and ASiC CAdES.
	DigestMatcherTypeManifestEntry DigestMatcherType = "MANIFEST_ENTRY"
	// DigestMatcherTypeCounterSignature: XAdES signed SignatureValue
	// (counter signature).
	DigestMatcherTypeCounterSignature DigestMatcherType = "COUNTER_SIGNATURE"
	// DigestMatcherTypeMessageDigest: CAdES.
	DigestMatcherTypeMessageDigest DigestMatcherType = "MESSAGE_DIGEST"
	// DigestMatcherTypeContentDigest: digest from decrypted content
	// SignatureValue (CAdES/PAdES).
	DigestMatcherTypeContentDigest DigestMatcherType = "CONTENT_DIGEST"
	// DigestMatcherTypeJWSSigningInput: JAdES Digest on result of
	// concatenation
	// ASCII(BASE64URL(UTF8(JWSProtected Header)) || '.' ||
	// BASE64URL(JWS Payload)).
	DigestMatcherTypeJWSSigningInput DigestMatcherType = "JWS_SIGNING_INPUT"
	// DigestMatcherTypeSigDEntry: JAdES or CB-AdES Detached entry.
	DigestMatcherTypeSigDEntry DigestMatcherType = "SIG_D_ENTRY"
	// DigestMatcherTypeCoseSigStructure: COSE Digest on result of
	// serialization of Sig_structure array.
	DigestMatcherTypeCoseSigStructure DigestMatcherType = "COSE_SIG_STRUCTURE"
	// DigestMatcherTypeCounterSignedSignatureValue: defines the
	// signature value of a master signature signed by a counter signature.
	DigestMatcherTypeCounterSignedSignatureValue DigestMatcherType = "COUNTER_SIGNED_SIGNATURE_VALUE"
	// DigestMatcherTypeMessageImprint: timestamp.
	DigestMatcherTypeMessageImprint DigestMatcherType = "MESSAGE_IMPRINT"
	// DigestMatcherTypeEvidenceRecordArchiveObject: evidence record
	// archive object.
	DigestMatcherTypeEvidenceRecordArchiveObject DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_OBJECT"
	// DigestMatcherTypeEvidenceRecordOrphanReference: identifies
	// evidence record archive object which has not been associated with any
	// of the provided documents.
	DigestMatcherTypeEvidenceRecordOrphanReference DigestMatcherType = "EVIDENCE_RECORD_ORPHAN_REFERENCE"
	// DigestMatcherTypeEvidenceRecordArchiveTimeStamp: evidence record
	// previous archive time-stamp object.
	DigestMatcherTypeEvidenceRecordArchiveTimeStamp DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_TIME_STAMP"
	// DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence:
	// evidence record previous archive time-stamp sequence.
	DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence DigestMatcherType = "EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE"
	// DigestMatcherTypeEvidenceRecordMasterSignature: evidence record
	// embedded in a signature.
	DigestMatcherTypeEvidenceRecordMasterSignature DigestMatcherType = "EVIDENCE_RECORD_MASTER_SIGNATURE"
	// DigestMatcherTypeEAADisclosure: disclosure attached to a
	// presentation of EAA.
	DigestMatcherTypeEAADisclosure DigestMatcherType = "EAA_DISCLOSURE"
	// DigestMatcherTypeEAANestedDisclosure: disclosure nested to
	// provided disclosure to a presentation of EAA.
	DigestMatcherTypeEAANestedDisclosure DigestMatcherType = "EAA_NESTED_DISCLOSURE"
	// DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim:
	// incorporated SD claim for which no matching provided disclosure has
	// been found.
	DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim DigestMatcherType = "EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM"
	// DigestMatcherTypeEAAKeyBinding: input used to compute a key
	// binding signature (used in EAA).
	DigestMatcherTypeEAAKeyBinding DigestMatcherType = "EAA_KEY_BINDING"
)

// DigestMatcherTypeValues returns all constants in declaration order.
func DigestMatcherTypeValues() []DigestMatcherType {
	return []DigestMatcherType{
		DigestMatcherTypeReference,
		DigestMatcherTypeObject,
		DigestMatcherTypeManifest,
		DigestMatcherTypeSignedProperties,
		DigestMatcherTypeKeyInfo,
		DigestMatcherTypeSignatureProperties,
		DigestMatcherTypeXPointer,
		DigestMatcherTypeManifestEntry,
		DigestMatcherTypeCounterSignature,
		DigestMatcherTypeMessageDigest,
		DigestMatcherTypeContentDigest,
		DigestMatcherTypeJWSSigningInput,
		DigestMatcherTypeSigDEntry,
		DigestMatcherTypeCoseSigStructure,
		DigestMatcherTypeCounterSignedSignatureValue,
		DigestMatcherTypeMessageImprint,
		DigestMatcherTypeEvidenceRecordArchiveObject,
		DigestMatcherTypeEvidenceRecordOrphanReference,
		DigestMatcherTypeEvidenceRecordArchiveTimeStamp,
		DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence,
		DigestMatcherTypeEvidenceRecordMasterSignature,
		DigestMatcherTypeEAADisclosure,
		DigestMatcherTypeEAANestedDisclosure,
		DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim,
		DigestMatcherTypeEAAKeyBinding,
	}
}

// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSConstants.java (DSS 6.5.RC1).
//
// A Java utility class of public static final String fields collapses to one Go const block, with
// the PAdESConstants naming already in the tree: <ClassName><FieldName>, the field name having
// its underscores dropped and each word title-cased.
package jades

const (
	// JWSConstantsPayload is the "payload" member: the signed document. Port of PAYLOAD.
	JWSConstantsPayload = "payload"

	// JWSConstantsSignatures is the "signatures" member: the array of signatures of a complete
	// JWS JSON Serialization. Port of SIGNATURES.
	JWSConstantsSignatures = "signatures"

	// JWSConstantsProtected is the "protected" member: the signed properties, base64url-encoded.
	// Port of PROTECTED.
	JWSConstantsProtected = "protected"

	// JWSConstantsHeader is the "header" member: the unsigned properties. Port of HEADER.
	JWSConstantsHeader = "header"

	// JWSConstantsSignature is the "signature" member: the container for a SignatureValue.
	// Port of SIGNATURE.
	JWSConstantsSignature = "signature"
)

// Ported from dss-enumerations/.../Context.java (DSS 6.5.RC1).
package enumerations

// Context defines signature validation context.
type Context string

const (
	// ContextSignature: the processing token is a signature.
	ContextSignature Context = "SIGNATURE"
	// ContextCounterSignature: the processing token is a counter signature.
	ContextCounterSignature Context = "COUNTER_SIGNATURE"
	// ContextKeyBindingSignature: the processing token is a key binding
	// signature.
	ContextKeyBindingSignature Context = "KEY_BINDING_SIGNATURE"
	// ContextTimestamp: the processing token is a timestamp.
	ContextTimestamp Context = "TIMESTAMP"
	// ContextEvidenceRecord: the processing token is an evidence record.
	ContextEvidenceRecord Context = "EVIDENCE_RECORD"
	// ContextRevocation: the processing token is a revocation.
	ContextRevocation Context = "REVOCATION"
	// ContextCertificate: the processing token is a certificate.
	ContextCertificate Context = "CERTIFICATE"
	// ContextEAA: the processing token is an electronic attestation of
	// attributes.
	ContextEAA Context = "EAA"
	// ContextEAARevocation: the processing token is an EAA revocation
	// token.
	ContextEAARevocation Context = "EAA_REVOCATION"
)

// ContextValues returns all constants in declaration order.
func ContextValues() []Context {
	return []Context{
		ContextSignature,
		ContextCounterSignature,
		ContextKeyBindingSignature,
		ContextTimestamp,
		ContextEvidenceRecord,
		ContextRevocation,
		ContextCertificate,
		ContextEAA,
		ContextEAARevocation,
	}
}

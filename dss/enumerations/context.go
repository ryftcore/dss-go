// Ported from dss-enumerations/.../Context.java (DSS 6.5.RC1).
package enumerations

// Context defines signature validation context.
type Context string

const (
	// Context_SIGNATURE: the processing token is a signature.
	Context_SIGNATURE Context = "SIGNATURE"
	// Context_COUNTER_SIGNATURE: the processing token is a counter signature.
	Context_COUNTER_SIGNATURE Context = "COUNTER_SIGNATURE"
	// Context_KEY_BINDING_SIGNATURE: the processing token is a key binding
	// signature.
	Context_KEY_BINDING_SIGNATURE Context = "KEY_BINDING_SIGNATURE"
	// Context_TIMESTAMP: the processing token is a timestamp.
	Context_TIMESTAMP Context = "TIMESTAMP"
	// Context_EVIDENCE_RECORD: the processing token is an evidence record.
	Context_EVIDENCE_RECORD Context = "EVIDENCE_RECORD"
	// Context_REVOCATION: the processing token is a revocation.
	Context_REVOCATION Context = "REVOCATION"
	// Context_CERTIFICATE: the processing token is a certificate.
	Context_CERTIFICATE Context = "CERTIFICATE"
	// Context_EAA: the processing token is an electronic attestation of
	// attributes.
	Context_EAA Context = "EAA"
	// Context_EAA_REVOCATION: the processing token is an EAA revocation
	// token.
	Context_EAA_REVOCATION Context = "EAA_REVOCATION"
)

// ContextValues returns all constants in declaration order.
func ContextValues() []Context {
	return []Context{
		Context_SIGNATURE,
		Context_COUNTER_SIGNATURE,
		Context_KEY_BINDING_SIGNATURE,
		Context_TIMESTAMP,
		Context_EVIDENCE_RECORD,
		Context_REVOCATION,
		Context_CERTIFICATE,
		Context_EAA,
		Context_EAA_REVOCATION,
	}
}

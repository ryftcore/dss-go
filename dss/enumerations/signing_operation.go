// Ported from dss-enumerations/.../SigningOperation.java (DSS 6.5.RC1).
package enumerations

// SigningOperation defines the performing signature operation.
type SigningOperation string

const (
	// SigningOperation_SIGN is signature creation.
	SigningOperation_SIGN SigningOperation = "SIGN"
	// SigningOperation_COUNTER_SIGN is counter signature creation.
	SigningOperation_COUNTER_SIGN SigningOperation = "COUNTER_SIGN"
	// SigningOperation_TIMESTAMP is timestamp creation.
	SigningOperation_TIMESTAMP SigningOperation = "TIMESTAMP"
	// SigningOperation_EXTEND is extension process (eg : signature basic
	// (B-P-B) to signature with timestamp (B-P-T)).
	SigningOperation_EXTEND SigningOperation = "EXTEND"
	// SigningOperation_ADD_SIG_POLICY_STORE is the signature policy store
	// addition.
	SigningOperation_ADD_SIG_POLICY_STORE SigningOperation = "ADD_SIG_POLICY_STORE"
	// SigningOperation_ADD_EVIDENCE_RECORD is the evidence record addition.
	SigningOperation_ADD_EVIDENCE_RECORD SigningOperation = "ADD_EVIDENCE_RECORD"
	// SigningOperation_EAA_PRESENTATION is the EAA presentation.
	SigningOperation_EAA_PRESENTATION SigningOperation = "EAA_PRESENTATION"
)

// SigningOperationValues returns all constants in declaration order.
func SigningOperationValues() []SigningOperation {
	return []SigningOperation{
		SigningOperation_SIGN,
		SigningOperation_COUNTER_SIGN,
		SigningOperation_TIMESTAMP,
		SigningOperation_EXTEND,
		SigningOperation_ADD_SIG_POLICY_STORE,
		SigningOperation_ADD_EVIDENCE_RECORD,
		SigningOperation_EAA_PRESENTATION,
	}
}

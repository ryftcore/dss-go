// Ported from dss-enumerations/.../SigningOperation.java (DSS 6.5.RC1).
package enumerations

// SigningOperation defines the performing signature operation.
type SigningOperation string

const (
	// SigningOperationSign is signature creation.
	SigningOperationSign SigningOperation = "SIGN"
	// SigningOperationCounterSign is counter signature creation.
	SigningOperationCounterSign SigningOperation = "COUNTER_SIGN"
	// SigningOperationTimestamp is timestamp creation.
	SigningOperationTimestamp SigningOperation = "TIMESTAMP"
	// SigningOperationExtend is extension process (eg : signature basic
	// (B-P-B) to signature with timestamp (B-P-T)).
	SigningOperationExtend SigningOperation = "EXTEND"
	// SigningOperationAddSigPolicyStore is the signature policy store
	// addition.
	SigningOperationAddSigPolicyStore SigningOperation = "ADD_SIG_POLICY_STORE"
	// SigningOperationAddEvidenceRecord is the evidence record addition.
	SigningOperationAddEvidenceRecord SigningOperation = "ADD_EVIDENCE_RECORD"
	// SigningOperationEAAPresentation is the EAA presentation.
	SigningOperationEAAPresentation SigningOperation = "EAA_PRESENTATION"
)

// SigningOperationValues returns all constants in declaration order.
func SigningOperationValues() []SigningOperation {
	return []SigningOperation{
		SigningOperationSign,
		SigningOperationCounterSign,
		SigningOperationTimestamp,
		SigningOperationExtend,
		SigningOperationAddSigPolicyStore,
		SigningOperationAddEvidenceRecord,
		SigningOperationEAAPresentation,
	}
}

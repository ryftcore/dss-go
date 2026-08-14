// Ported from dss-enumerations/.../QCIdentMethod.java (DSS 6.5.RC1).
package enumerations

// QCIdentMethod defines QC Identification Method OID identifiers as defined
// in the ETSI EN 319 412-5 "4.3.5 QCStatement stating the used
// eIDAS/eIDAS2 Article 24. identification method".
//
// NOTE: QCIdentMethodEnum (the concrete enum implementing this interface,
// with its QCIdentMethodEnumValues() accessor) is defined outside this
// file's manifest and is assumed to exist per the porting brief.
type QCIdentMethod interface {
	OidDescription
}

// QCIdentMethod_UNKNOWN_METHOD defines a description for an unknown method
// by the current implementation.
const QCIdentMethod_UNKNOWN_METHOD = "qc-identification-method-unknown"

// qcIdentMethodUnknown is a plain QCIdentMethod implementation backing the
// fallback branch of QCIdentMethodFromOID (Java's QCIdentMethod.fromOid
// anonymous class).
type qcIdentMethodUnknown struct {
	oid string
}

func (q *qcIdentMethodUnknown) Description() string { return QCIdentMethod_UNKNOWN_METHOD }
func (q *qcIdentMethodUnknown) OID() string         { return q.oid }

// QCIdentMethodFromOID returns a QCIdentMethod by the given OID, if it
// matches a known QCIdentMethodEnum constant. Otherwise, a QCIdentMethod
// with QCIdentMethod_UNKNOWN_METHOD as its Description and the given oid is
// returned (Java logged this case at debug level via SLF4J; that logging
// side-effect is not ported — see PORTER_BRIEF notes).
func QCIdentMethodFromOID(oid string) QCIdentMethod {
	for _, t := range QCIdentMethodEnumValues() {
		if t.OID() == oid {
			return t
		}
	}
	return &qcIdentMethodUnknown{oid: oid}
}

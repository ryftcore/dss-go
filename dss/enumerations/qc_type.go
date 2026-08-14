// Ported from dss-enumerations/.../QCType.java (DSS 6.5.RC1).
package enumerations

// QCType defines QC Type OID identifiers.
//
// NOTE: QCTypeEnum (the concrete enum implementing this interface, with its
// QCTypeEnumValues() accessor) is defined outside this file's manifest and
// is assumed to exist per the porting brief.
type QCType interface {
	OidDescription
}

// QCType_UNKNOWN_TYPE defines a description for a type unknown by the
// current implementation.
const QCType_UNKNOWN_TYPE = "type-unknown"

// qcTypeUnknown is a plain QCType implementation backing the fallback
// branch of QCTypeFromOID (Java's QCType.fromOid anonymous class).
type qcTypeUnknown struct {
	oid string
}

func (q *qcTypeUnknown) Description() string { return QCType_UNKNOWN_TYPE }
func (q *qcTypeUnknown) OID() string         { return q.oid }

// QCTypeFromOID returns a QCType by the given OID, if it matches a known
// QCTypeEnum constant. Otherwise, a QCType with QCType_UNKNOWN_TYPE as its
// Description and the given oid is returned (Java logged this case at debug
// level via SLF4J; that logging side-effect is not ported — see
// PORTER_BRIEF notes).
func QCTypeFromOID(oid string) QCType {
	for _, t := range QCTypeEnumValues() {
		if t.OID() == oid {
			return t
		}
	}
	return &qcTypeUnknown{oid: oid}
}

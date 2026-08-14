// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/QCStatementOids.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"reflect"
)

// QCStatementOids represents a collection of properties extracted from an MRA condition.
//
// java.io.Serializable has no Go counterpart and is dropped.
type QCStatementOids struct {
	// qcStatementIds is the list of QcStatement identifiers to be included.
	qcStatementIds []string
	// qcTypeIds is the list of QcType identifiers to be included.
	qcTypeIds []string
	// qcCClegislations is the list of QcCClegislation codes to be included.
	qcCClegislations []string
	// qcStatementIdsToRemove is the list of QcStatement identifiers to be removed.
	qcStatementIdsToRemove []string
	// qcTypeIdsToRemove is the list of QcType identifiers to be removed.
	qcTypeIdsToRemove []string
	// qcCClegislationsToRemove is the list of QcCClegislation codes to be removed.
	qcCClegislationsToRemove []string
}

// NewQCStatementOids instantiates an object with zero values. Port of the default constructor.
func NewQCStatementOids() *QCStatementOids {
	return &QCStatementOids{}
}

// QcStatementIds gets QcStatement identifiers to be included.
func (q *QCStatementOids) QcStatementIds() []string {
	return q.qcStatementIds
}

// SetQcStatementIds sets QcStatement identifiers to be included.
func (q *QCStatementOids) SetQcStatementIds(qcStatementIds []string) {
	q.qcStatementIds = qcStatementIds
}

// QcTypeIds gets QcType identifiers to be included.
func (q *QCStatementOids) QcTypeIds() []string {
	return q.qcTypeIds
}

// SetQcTypeIds sets QcType identifiers to be included.
func (q *QCStatementOids) SetQcTypeIds(qcTypeIds []string) {
	q.qcTypeIds = qcTypeIds
}

// QcCClegislations gets QcCClegislation codes to be included.
func (q *QCStatementOids) QcCClegislations() []string {
	return q.qcCClegislations
}

// SetQcCClegislations sets QcCClegislation codes to be included.
func (q *QCStatementOids) SetQcCClegislations(qcCClegislations []string) {
	q.qcCClegislations = qcCClegislations
}

// QcStatementIdsToRemove gets QcStatement identifiers to be removed.
func (q *QCStatementOids) QcStatementIdsToRemove() []string {
	return q.qcStatementIdsToRemove
}

// SetQcStatementIdsToRemove sets QcStatement identifiers to be removed.
func (q *QCStatementOids) SetQcStatementIdsToRemove(qcStatementIdsToRemove []string) {
	q.qcStatementIdsToRemove = qcStatementIdsToRemove
}

// QcTypeIdsToRemove gets QcType identifiers to be removed.
func (q *QCStatementOids) QcTypeIdsToRemove() []string {
	return q.qcTypeIdsToRemove
}

// SetQcTypeIdsToRemove sets QcType identifiers to be removed.
func (q *QCStatementOids) SetQcTypeIdsToRemove(qcTypeIdsToRemove []string) {
	q.qcTypeIdsToRemove = qcTypeIdsToRemove
}

// QcCClegislationsToRemove gets QcCClegislation codes to be removed.
func (q *QCStatementOids) QcCClegislationsToRemove() []string {
	return q.qcCClegislationsToRemove
}

// SetQcCClegislationsToRemove sets QcCClegislation codes to be removed.
func (q *QCStatementOids) SetQcCClegislationsToRemove(qcCClegislationsToRemove []string) {
	q.qcCClegislationsToRemove = qcCClegislationsToRemove
}

// String returns the Java toString() form.
func (q *QCStatementOids) String() string {
	return fmt.Sprintf("QCStatementOids [qcStatementIds=%s, qcTypeIds=%s, qcCClegislations=%s, "+
		"qcStatementIdsToRemove=%s, qcTypeIdsToRemove=%s, qcCClegislationsToRemove=%s]",
		qcStatementOidsJavaList(q.qcStatementIds), qcStatementOidsJavaList(q.qcTypeIds), qcStatementOidsJavaList(q.qcCClegislations),
		qcStatementOidsJavaList(q.qcStatementIdsToRemove), qcStatementOidsJavaList(q.qcTypeIdsToRemove), qcStatementOidsJavaList(q.qcCClegislationsToRemove))
}

// qcStatementOidsJavaList renders a []string the way java.util.List#toString() would:
// "null" for a nil list, "[]" for empty, "[a, b]" otherwise.
func qcStatementOidsJavaList(list []string) string {
	if list == nil {
		return "null"
	}
	out := "["
	for i, v := range list {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out + "]"
}

// Equals ports QCStatementOids#equals(Object): a field-by-field comparison against another
// QCStatementOids of the same concrete type.
func (q *QCStatementOids) Equals(other *QCStatementOids) bool {
	if other == nil {
		return false
	}
	if q == other {
		return true
	}
	return reflect.DeepEqual(q.qcStatementIds, other.qcStatementIds) &&
		reflect.DeepEqual(q.qcTypeIds, other.qcTypeIds) &&
		reflect.DeepEqual(q.qcCClegislations, other.qcCClegislations) &&
		reflect.DeepEqual(q.qcStatementIdsToRemove, other.qcStatementIdsToRemove) &&
		reflect.DeepEqual(q.qcTypeIdsToRemove, other.qcTypeIdsToRemove) &&
		reflect.DeepEqual(q.qcCClegislationsToRemove, other.qcCClegislationsToRemove)
}

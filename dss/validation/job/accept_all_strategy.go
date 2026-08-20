// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/sync/AcceptAllStrategy.java (DSS 6.5.RC1).
package job

import modeljob "github.com/utain/esig/dss/model/job"

// AcceptAllStrategy accepts all trusted lists. D, L mirror
// SynchronizationStrategy's type parameters.
type AcceptAllStrategy[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]] struct{}

// NewAcceptAllStrategy creates an AcceptAllStrategy. Port of the default constructor.
func NewAcceptAllStrategy[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]]() *AcceptAllStrategy[D, L] {
	return &AcceptAllStrategy[D, L]{}
}

// CanBeSynchronizedDocument always returns true. Port of canBeSynchronized(D).
func (a *AcceptAllStrategy[D, L]) CanBeSynchronizedDocument(document D) bool {
	return true
}

// CanBeSynchronizedDocumentList always returns true. Port of canBeSynchronized(L).
func (a *AcceptAllStrategy[D, L]) CanBeSynchronizedDocumentList(documentList L) bool {
	return true
}

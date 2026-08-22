// Ported from dss-validation-job/src/main/java/eu/europa/esig/dss/validation/job/sync/SynchronizationStrategy.java (DSS 6.5.RC1).
package job

import modeljob "github.com/ryftcore/dss-go/dss/model/job"

// SynchronizationStrategy defines a behaviour for a trusted certificate source
// synchronization. D is the current modeljob.DocumentInfo, L the parent
// modeljob.DocumentListInfo, mirroring Java's "<D extends DocumentInfo<L>, L extends
// DocumentListInfo<L, D>>".
//
// Java overloads canBeSynchronized(D) / canBeSynchronized(L); Go has no overloading, so each
// overload gets a distinct name: CanBeSynchronizedDocument / CanBeSynchronizedDocumentList.
type SynchronizationStrategy[D modeljob.DocumentInfo[L], L modeljob.DocumentListInfo[L, D]] interface {
	// CanBeSynchronizedDocument returns true if the certificates from the trusted list can
	// be synchronized. Port of canBeSynchronized(D).
	CanBeSynchronizedDocument(document D) bool

	// CanBeSynchronizedDocumentList returns true if the certificates from the list of
	// trusted lists and its trusted list can be synchronized. Port of canBeSynchronized(L).
	CanBeSynchronizedDocumentList(documentList L) bool
}

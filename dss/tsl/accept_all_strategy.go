// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sync/AcceptAllStrategy.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see xml_download_result.go's header for the wider job.* convention):
// implements job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo]. Java overloads
// canBeSynchronized(D)/canBeSynchronized(L); Go has no overloading, so the assumed interface
// (like ExpirationAndSignatureCheckStrategy.java's, see expiration_and_signature_check_strategy.go)
// spells the two methods CanBeSynchronizedDocument(D) bool / CanBeSynchronizedDocumentList(L) bool.
package tsl

import (
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// AcceptAllStrategy accepts all trusted lists.
//
// Deprecated: since DSS 6.5. Use github.com/ryftcore/dss-go/dss/validation/job.AcceptAllStrategy instead.
type AcceptAllStrategy struct{}

var _ job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo] = (*AcceptAllStrategy)(nil)

// NewAcceptAllStrategy is the default constructor.
//
// Deprecated: since DSS 6.5.
func NewAcceptAllStrategy() *AcceptAllStrategy {
	return &AcceptAllStrategy{}
}

// CanBeSynchronizedDocument ports canBeSynchronized(TLInfo).
func (s *AcceptAllStrategy) CanBeSynchronizedDocument(trustedList *tslmodel.TLInfo) bool {
	return true
}

// CanBeSynchronizedDocumentList ports canBeSynchronized(LOTLInfo).
func (s *AcceptAllStrategy) CanBeSynchronizedDocumentList(listOfTrustedList *tslmodel.LOTLInfo) bool {
	return true
}

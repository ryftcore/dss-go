// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sync/ExpirationAndSignatureCheckStrategy.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see accept_all_strategy.go's header): implements
// job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo].
package tsl

import (
	"time"

	tslmodel "github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/validation/job"
)

// ExpirationAndSignatureCheckStrategy allows skipping expired or invalid trusted lists.
type ExpirationAndSignatureCheckStrategy struct {
	// acceptExpiredTrustedList defines if expired trusted lists (next update after current
	// time) are supported.
	acceptExpiredTrustedList bool

	// acceptInvalidTrustedList defines if trusted lists with invalid or indeterminate
	// signatures are supported.
	acceptInvalidTrustedList bool

	// acceptExpiredListOfTrustedLists defines if expired list of trusted lists (next update
	// after current time) are supported.
	acceptExpiredListOfTrustedLists bool

	// acceptInvalidListOfTrustedLists defines if list of trusted lists with invalid or
	// indeterminate signatures are supported.
	acceptInvalidListOfTrustedLists bool
}

var _ job.SynchronizationStrategy[*tslmodel.TLInfo, *tslmodel.LOTLInfo] = (*ExpirationAndSignatureCheckStrategy)(nil)

// NewExpirationAndSignatureCheckStrategy is the default constructor instantiating object with
// null values.
func NewExpirationAndSignatureCheckStrategy() *ExpirationAndSignatureCheckStrategy {
	return &ExpirationAndSignatureCheckStrategy{}
}

// SetAcceptExpiredTrustedList sets if expired trusted lists are supported (next update after
// current time). Port of setAcceptExpiredTrustedList(boolean).
func (s *ExpirationAndSignatureCheckStrategy) SetAcceptExpiredTrustedList(acceptExpiredTrustedList bool) {
	s.acceptExpiredTrustedList = acceptExpiredTrustedList
}

// SetAcceptInvalidTrustedList sets if invalid trusted lists are supported (signature with
// FAILED or INDETERMINATE Indication). Port of setAcceptInvalidTrustedList(boolean).
func (s *ExpirationAndSignatureCheckStrategy) SetAcceptInvalidTrustedList(acceptInvalidTrustedList bool) {
	s.acceptInvalidTrustedList = acceptInvalidTrustedList
}

// SetAcceptExpiredListOfTrustedLists sets if expired list of trusted lists and their TLs are
// supported (next update after current time). Port of setAcceptExpiredListOfTrustedLists(boolean).
func (s *ExpirationAndSignatureCheckStrategy) SetAcceptExpiredListOfTrustedLists(acceptExpiredListOfTrustedLists bool) {
	s.acceptExpiredListOfTrustedLists = acceptExpiredListOfTrustedLists
}

// SetAcceptInvalidListOfTrustedLists sets if invalid list of trusted lists and their TLs are
// supported (signature with FAILED or INDETERMINATE Indication). Port of
// setAcceptInvalidListOfTrustedLists(boolean).
func (s *ExpirationAndSignatureCheckStrategy) SetAcceptInvalidListOfTrustedLists(acceptInvalidListOfTrustedLists bool) {
	s.acceptInvalidListOfTrustedLists = acceptInvalidListOfTrustedLists
}

// CanBeSynchronizedDocument ports canBeSynchronized(TLInfo).
func (s *ExpirationAndSignatureCheckStrategy) CanBeSynchronizedDocument(trustedList *tslmodel.TLInfo) bool {
	return s.isSyncSupported(trustedList, s.acceptExpiredTrustedList, s.acceptInvalidTrustedList)
}

// CanBeSynchronizedDocumentList ports canBeSynchronized(LOTLInfo).
func (s *ExpirationAndSignatureCheckStrategy) CanBeSynchronizedDocumentList(listOfTrustedList *tslmodel.LOTLInfo) bool {
	return s.isSyncSupported(&listOfTrustedList.TLInfo, s.acceptExpiredListOfTrustedLists, s.acceptInvalidListOfTrustedLists)
}

// isSyncSupported ports the private isSyncSupported(TLInfo, boolean, boolean).
func (s *ExpirationAndSignatureCheckStrategy) isSyncSupported(tlInfo *tslmodel.TLInfo, syncExpired, syncInvalid bool) bool {
	if !syncExpired {
		if parsingCacheInfo, ok := tlInfo.TLParsingCacheInfo(); ok && parsingCacheInfo.IsResultExist() {
			currentDate := time.Now()
			nextUpdateDate := parsingCacheInfo.NextUpdateDate()
			if nextUpdateDate.IsZero() || currentDate.After(nextUpdateDate) {
				return false
			}
		}
	}

	if !syncInvalid {
		validationCacheInfo := tlInfo.ValidationCacheInfo()
		if validationCacheInfo != nil && validationCacheInfo.IsResultExist() {
			return validationCacheInfo.IsValid()
		}
	}

	return true
}

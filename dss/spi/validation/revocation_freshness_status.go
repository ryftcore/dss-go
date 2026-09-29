// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/status/RevocationFreshnessStatus.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.validation.status.RevocationFreshnessStatus lands in this same Go
// package; the sibling status type TokenStatus (Java spi.validation.status.TokenStatus) is
// landed in this same package too.
//
// This type is load-bearing: signature_validation_context.go passes either a *TokenStatus or a
// *RevocationFreshnessStatus to its shared checking logic through a small recorder interface
// (both satisfy it, the latter through the promoted AddRelatedTokenAndErrorMessage) and
// type-asserts back to *RevocationFreshnessStatus where Java tests `status instanceof
// RevocationFreshnessStatus`; that requires RevocationFreshnessStatus to embed TokenStatus as
// an unnamed value field (exactly mirroring Java's "extends TokenStatus").
//
// Go has no virtual dispatch through embedding: alert.ObjectStatus.String() (promoted via
// TokenStatus's embedded *alert.ObjectStatus) calls its own receiver's ErrorString(), not this
// type's override. String() is therefore re-implemented here identically to
// alert.ObjectStatus.String()'s logic so it picks up this type's ErrorString() override, the
// same pattern SignatureStatus's sibling file would need if it had overridden ErrorString (it
// does not, so it relies on plain promotion).
//
// alert.ObjectStatus.objectMapToString(), the private helper Java's getErrorString() computes
// its bracketed related-objects suffix with, is unexported in package alert and therefore
// unreachable from this package even via promotion (Go unexported identifiers are
// package-scoped, promotion notwithstanding). ErrorString() below recovers that suffix by
// stripping the known "Message() + \" \"" prefix off the promoted TokenStatus.ErrorString()
// (== alert.ObjectStatus.ErrorString()) result, which is exactly that suffix by construction.
package validation

import (
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// RevocationFreshnessStatus contains information about the performed revocation freshness
// check.
type RevocationFreshnessStatus struct {
	TokenStatus

	// tokenRevocationNextUpdateMap maps tokens concerned by the revocation freshness check to
	// the nextUpdate time of the corresponding revocation data.
	tokenRevocationNextUpdateMap map[model.Token]time.Time
}

// NewRevocationFreshnessStatus is the default constructor initializing an empty map.
func NewRevocationFreshnessStatus() *RevocationFreshnessStatus {
	return &RevocationFreshnessStatus{
		TokenStatus:                  *NewTokenStatus(),
		tokenRevocationNextUpdateMap: make(map[model.Token]time.Time),
	}
}

// AddTokenAndRevocationNextUpdateTime adds a concerned token and nextUpdate time of the
// revocation data. Port of addTokenAndRevocationNextUpdateTime(Token, Date).
func (s *RevocationFreshnessStatus) AddTokenAndRevocationNextUpdateTime(token model.Token, revocationNextUpdate time.Time) {
	s.tokenRevocationNextUpdateMap[token] = revocationNextUpdate
}

// TokenRevocationNextUpdateTime returns nextUpdate time of revocation data for the given token.
//
// NOTE: returns a non-zero time only if the obtained revocation data is not fresh enough
// (otherwise returns the zero Time). Port of getTokenRevocationNextUpdateTime(Token).
func (s *RevocationFreshnessStatus) TokenRevocationNextUpdateTime(token model.Token) time.Time {
	return s.tokenRevocationNextUpdateMap[token]
}

// MinimalNextUpdateTime returns the minimal time when revocation data should be updated for all
// concerned tokens.
//
// NOTE: returns the zero Time if no suitable revocation data found or if the revocation data is
// fresh enough.
//
// Faithful port of a quirk in the Java source: despite its name, the loop keeps the LATEST
// (not earliest) nextUpdate time seen (`minimalNextUpdate == null || minimalNextUpdate.before(nextUpdate)`
// always overwrites with the later value) - reproduced verbatim, not "fixed". Port of
// getMinimalNextUpdateTime().
func (s *RevocationFreshnessStatus) MinimalNextUpdateTime() time.Time {
	var minimalNextUpdate time.Time
	for _, nextUpdate := range s.tokenRevocationNextUpdateMap {
		if minimalNextUpdate.IsZero() || minimalNextUpdate.Before(nextUpdate) {
			minimalNextUpdate = nextUpdate
		}
	}
	return minimalNextUpdate
}

// ErrorString overrides TokenStatus's (promoted) ErrorString to append the minimal nextUpdate
// time, when present. Port of the getErrorString() override.
func (s *RevocationFreshnessStatus) ErrorString() string {
	nextUpdateTime := s.MinimalNextUpdateTime()
	if nextUpdateTime.IsZero() {
		return s.TokenStatus.ErrorString()
	}
	// Recover the "[id: msg; ...]" suffix objectMapToString() would have produced - see this
	// file's header comment for why it cannot be called directly.
	objectsSuffix := strings.TrimPrefix(s.TokenStatus.ErrorString(), s.Message()+" ")
	return s.Message() + " NextUpdate time : " + spi.DSSUtilsFormatDateToRFC(nextUpdateTime) + " " + objectsSuffix
}

// String returns "Status : Valid" when empty, this type's own ErrorString() otherwise - see
// this file's header comment on why String() cannot simply be promoted.
func (s *RevocationFreshnessStatus) String() string {
	if s.IsEmpty() {
		return "Status : Valid"
	}
	return s.ErrorString()
}

// compile-time assertion: a RevocationFreshnessStatus is an alert.Status.
var _ alert.Status = (*RevocationFreshnessStatus)(nil)

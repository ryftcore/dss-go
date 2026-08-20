// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sha2/DefaultTrustedListWithSha2Predicate.java (DSS 6.5.RC1).
package tsl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// DefaultTrustedListWithSha2PredicateOverrides captures every member Test() reaches through
// virtual dispatch: the two protected hooks this class declares (the current time NextUpdate is
// compared against, and the cache-expiry rule) plus the three protected helpers it inherits from
// AbstractTrustedListWithSha2Predicate, all of which a further subclass may override.
type DefaultTrustedListWithSha2PredicateOverrides interface {
	// CurrentTime returns the current time to compare NextUpdate against. Port of the
	// protected getCurrentTime().
	CurrentTime() time.Time

	// IsCacheExpired verifies whether the cache of the document has expired. Port of the
	// protected isCacheExpired(DSSDocument).
	IsCacheExpired(document model.DSSDocument) bool

	// OriginalDocumentDigest is AbstractTrustedListWithSha2Predicate's protected
	// getOriginalDocumentDigest(DSSDocument).
	OriginalDocumentDigest(document model.DSSDocument) (model.Digest, error)

	// Sha2Digest is AbstractTrustedListWithSha2Predicate's protected getSha2Digest(DSSDocument).
	Sha2Digest(sha2Document model.DSSDocument) (model.Digest, error)

	// NextUpdate is AbstractTrustedListWithSha2Predicate's protected getNextUpdate(DSSDocument).
	NextUpdate(tlDocument model.DSSDocument) time.Time
}

// DefaultTrustedListWithSha2Predicate is the default implementation of the ETSI TS 119 612
// sha2 file processing.
type DefaultTrustedListWithSha2Predicate struct {
	AbstractTrustedListWithSha2PredicateBase

	// overrides points back at the concrete predicate; see InitDefaultTrustedListWithSha2Predicate.
	overrides DefaultTrustedListWithSha2PredicateOverrides

	// cacheExpirationTime is the cache expiration time, in milliseconds, after which the
	// document shall be downloaded again.
	cacheExpirationTime int64
}

var (
	_ TrustedListWithSha2Predicate                 = (*DefaultTrustedListWithSha2Predicate)(nil)
	_ DefaultTrustedListWithSha2PredicateOverrides = (*DefaultTrustedListWithSha2Predicate)(nil)
)

// NewDefaultTrustedListWithSha2Predicate is the default constructor. Port of the default
// constructor, whose cacheExpirationTime field initialiser is -1 (the cache does not expire).
func NewDefaultTrustedListWithSha2Predicate() *DefaultTrustedListWithSha2Predicate {
	predicate := &DefaultTrustedListWithSha2Predicate{cacheExpirationTime: -1}
	predicate.InitDefaultTrustedListWithSha2Predicate(predicate)
	return predicate
}

// InitDefaultTrustedListWithSha2Predicate registers the concrete predicate with its base so
// that Test can dispatch CurrentTime/IsCacheExpired the way Java reaches an overridden method
// through virtual dispatch. It must be called exactly once, by the concrete predicate's
// constructor, before any other method.
func (p *DefaultTrustedListWithSha2Predicate) InitDefaultTrustedListWithSha2Predicate(
	overrides DefaultTrustedListWithSha2PredicateOverrides) {
	p.overrides = overrides
}

// defaultTrustedListWithSha2PredicateOverrides returns the registered overrides, panicking when
// the concrete predicate forgot to call InitDefaultTrustedListWithSha2Predicate.
func (p *DefaultTrustedListWithSha2Predicate) defaultTrustedListWithSha2PredicateOverrides() DefaultTrustedListWithSha2PredicateOverrides {
	if p.overrides == nil {
		panic("DefaultTrustedListWithSha2Predicate was not initialised: the concrete predicate must call InitDefaultTrustedListWithSha2Predicate in its constructor")
	}
	return p.overrides
}

// SetCacheExpirationTime sets the cache expiration time after which the document shall be
// refreshed, in milliseconds. Default: -1 (the cache does not expire). Port of
// setCacheExpirationTime(long).
func (p *DefaultTrustedListWithSha2Predicate) SetCacheExpirationTime(cacheExpirationTime int64) {
	p.cacheExpirationTime = cacheExpirationTime
}

// Test evaluates whether the sha2 digest matches the original document, returning true if the
// sha2 corresponding to the document matches the digest of the cached content and no refresh is
// required, false otherwise (if refresh is required for any reason). Port of
// test(DocumentWithSha2).
//
// Panics with the Java message when documentWithSha2 is nil (Objects.requireNonNull).
func (p *DefaultTrustedListWithSha2Predicate) Test(documentWithSha2 *DocumentWithSha2) (bool, error) {
	if documentWithSha2 == nil {
		panic("Document shall be provided!")
	}

	// For example, the TLSOx's TL published at the location
	// http://www.TLSOx.xyz/TrustedList/TL.xml is accompanied by its sha2 digest file i.e. on
	// location http://www.TLSOx.xyz/TrustedList/TL.sha2. Downloaders may adopt the following
	// strategy for downloading file TL.xml:
	document := documentWithSha2.Document()
	if document == nil {
		documentWithSha2.AddErrorMessage("No cached document has been found.")
		return false, nil // refresh required
	}

	// - check whether TL.sha2 is available for download:
	//     - if TL.sha2 has been successfully downloaded, verify the digest against the cached
	//       TL.xml file. If different, download and process TL.xml;
	//     - if TL.sha2 has not been successfully downloaded, download and process TL.xml
	//       directly.
	sha2Document := documentWithSha2.Sha2Document()
	if sha2Document == nil {
		documentWithSha2.AddErrorMessage("No sha2 document has been found.")
		return false, nil

	} else {
		overrides := p.defaultTrustedListWithSha2PredicateOverrides()
		originalDocumentDigest, err := overrides.OriginalDocumentDigest(document)
		if err != nil {
			return false, err
		}
		sha2Digest, err := overrides.Sha2Digest(sha2Document)
		if err != nil {
			return false, err
		}
		if !originalDocumentDigest.Equals(sha2Digest) {
			sha2Binaries, err := spi.DSSUtilsToByteArrayOfDocument(sha2Document)
			if err != nil {
				return false, err
			}
			errorMessage := fmt.Sprintf("Digest present within sha2 file '%s' do not match digest of "+
				"the cached document '%s'.",
				string(sha2Binaries), strings.ToLower(originalDocumentDigest.HexValue()))
			documentWithSha2.AddErrorMessage(errorMessage)
			return false, nil
		}
	}

	// - TL.xml should be downloaded/processed anyway if the nextUpdate (in the cached file) has
	//   been reached.
	nextUpdate := p.defaultTrustedListWithSha2PredicateOverrides().NextUpdate(document)
	if !nextUpdate.IsZero() && !nextUpdate.After(p.defaultTrustedListWithSha2PredicateOverrides().CurrentTime()) {
		documentWithSha2.AddErrorMessage(fmt.Sprintf("NextUpdate '%s' has been reached.",
			spi.DSSUtilsFormatDateToRFC(nextUpdate)))
		return false, nil
	}
	// Optional : validate cache expiration
	if p.defaultTrustedListWithSha2PredicateOverrides().IsCacheExpired(document) {
		// Upstream logs "Cache of the document with name '{}' expired. Request refresh.".
		return false, nil
	}
	// accept the document otherwise
	return true, nil
}

// CurrentTime returns the current time to compare NextUpdate against. Port of the protected
// getCurrentTime().
func (p *DefaultTrustedListWithSha2Predicate) CurrentTime() time.Time {
	return time.Now()
}

// IsCacheExpired verifies whether the cache of the document has expired. Port of the protected
// isCacheExpired(DSSDocument).
//
// NOTE: this method supports only the default model.FileDocument implementation. Override it
// (by embedding this type and re-registering through InitDefaultTrustedListWithSha2Predicate)
// should you need any processing of other implementations.
func (p *DefaultTrustedListWithSha2Predicate) IsCacheExpired(document model.DSSDocument) bool {
	if p.cacheExpirationTime < 0 {
		return false
	}
	if fileDocument, ok := document.(*model.FileDocument); ok {
		// Java reads getFile(), which the Go port exposes as Path() (no java.io.File type).
		info, err := os.Stat(fileDocument.Path())
		if err != nil {
			// java.io.File#exists() answers false for anything it cannot stat.
			return true
		}
		// java.io.File#lastModified() and Date#getTime() are both milliseconds since the
		// epoch; a file that cannot be read answers 0 there, which this branch cannot reach.
		currentTime := time.Now().UnixMilli()
		return (currentTime - info.ModTime().UnixMilli()) >= p.cacheExpirationTime
	}
	return true
}

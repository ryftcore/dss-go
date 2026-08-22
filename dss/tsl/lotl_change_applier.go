// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/job/LOTLChangeApplier.java (DSS 6.5.RC1).
//
// Uses job.CacheKey and job.ChangesCacheAccess (see xml_download_result.go's header for the
// wider job.* convention).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/job"
	validationjob "github.com/ryftcore/dss-go/dss/validation/job"
)

// LOTLChangeApplier applies changes in the LOTL cache.
type LOTLChangeApplier struct {
	// cacheAccess accesses the TL caches.
	cacheAccess *validationjob.ChangesCacheAccess

	// oldValues holds the old cache values.
	oldValues map[validationjob.CacheKey]job.ParsingInfoRecord

	// newValues holds the new cache values.
	newValues map[validationjob.CacheKey]job.ParsingInfoRecord
}

// NewLOTLChangeApplier is the default constructor.
func NewLOTLChangeApplier(cacheAccess *validationjob.ChangesCacheAccess,
	oldValues, newValues map[validationjob.CacheKey]job.ParsingInfoRecord) *LOTLChangeApplier {
	return &LOTLChangeApplier{cacheAccess: cacheAccess, oldValues: oldValues, newValues: newValues}
}

// AnalyzeAndApply applies changes for all defined records. Port of analyzeAndApply().
func (a *LOTLChangeApplier) AnalyzeAndApply() {
	for oldKey, oldValue := range a.oldValues {
		oldUrlCerts := a.tLPointers(oldValue)
		newUrlCerts := a.tLPointers(a.newValues[oldKey])

		a.detectUrlChanges(oldUrlCerts, newUrlCerts)
		a.detectSigCertsChanges(oldUrlCerts, newUrlCerts)
	}
}

// tLPointers ports the private getTLPointers(ParsingInfoRecord). Panics with the Java message
// when the record exists but is not a *TLParsingCacheDTO.
func (a *LOTLChangeApplier) tLPointers(parsingCache job.ParsingInfoRecord) map[string][]*model.CertificateToken {
	if parsingCache == nil || !parsingCache.IsResultExist() {
		return nil
	}
	tlParsingCacheDTO, ok := parsingCache.(*TLParsingCacheDTO)
	if !ok {
		panic("Parsing cache is not a TLParsingCacheDTO")
	}
	tlOtherPointers := tlParsingCacheDTO.TlOtherPointers()
	if len(tlOtherPointers) == 0 {
		return nil
	}
	result := make(map[string][]*model.CertificateToken, len(tlOtherPointers))
	for _, pointer := range tlOtherPointers {
		result[pointer.TSLLocation()] = pointer.SdiCertificates()
	}
	return result
}

func (a *LOTLChangeApplier) detectUrlChanges(oldUrlCerts, newUrlCerts map[string][]*model.CertificateToken) {
	for oldUrl := range oldUrlCerts {
		if _, ok := newUrlCerts[oldUrl]; !ok {
			a.cacheAccess.ToBeDeleted(validationjob.NewCacheKey(oldUrl))
		}
	}
}

func (a *LOTLChangeApplier) detectSigCertsChanges(oldUrlCerts, newUrlCerts map[string][]*model.CertificateToken) {
	for newUrl, newCerts := range newUrlCerts {
		oldCerts, ok := oldUrlCerts[newUrl]
		if ok && !certificateTokenSlicesEqual(oldCerts, newCerts) {
			a.cacheAccess.ExpireSignatureValidation(validationjob.NewCacheKey(newUrl))
		}
	}
}

// certificateTokenSlicesEqual compares two CertificateToken slices the way Java's
// List#equals compares two lists: same length, same elements in order, per-element equality by
// CertificateToken's own equals() (Equals here).
func certificateTokenSlicesEqual(a, b []*model.CertificateToken) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

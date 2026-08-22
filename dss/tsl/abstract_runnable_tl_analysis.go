// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/AbstractRunnableTLAnalysis.java (DSS 6.5.RC1).
package tsl

import (
	"sync"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/validation/job"
)

// AbstractRunnableTLAnalysis is the abstract implementation for performing a Trusted List
// analysis.
type AbstractRunnableTLAnalysis struct {
	job.AbstractRunnableAnalysis
}

// InitAbstractRunnableTLAnalysis is the default constructor. tlSource's embedded
// job.DocumentSource is what the base ends up storing/comparing by CacheKey.
func InitAbstractRunnableTLAnalysis(tlSource *TLSource, cacheAccess job.CacheAccessByKey, dssFileLoader http.DSSFileLoader, latch *sync.WaitGroup) AbstractRunnableTLAnalysis {
	return AbstractRunnableTLAnalysis{
		AbstractRunnableAnalysis: job.NewAbstractRunnableAnalysis(tlSource.DocumentSource, cacheAccess, dssFileLoader, latch),
	}
}

// GetDownloadTask ports the protected getDownloadTask(DSSFileLoader, String) override.
func (a *AbstractRunnableTLAnalysis) GetDownloadTask(dssFileLoader http.DSSFileLoader, url string) job.DownloadTask {
	return NewXmlDownloadTask(dssFileLoader, url)
}

// GetValidationTask ports the protected getValidationTask(DSSDocument, CertificateSource)
// override.
func (a *AbstractRunnableTLAnalysis) GetValidationTask(document model.DSSDocument, certificateSource spi.CertificateSource) job.ValidationTask {
	return tlValidatorTaskAdapter{task: NewTLValidatorTask(document, certificateSource)}
}

// tlValidatorTaskAdapter adapts *TLValidatorTask to job.ValidationTask.
//
// FLAG (pre-existing, tsl/tl_validator_task.go - out of this manifest, ported by the validation
// chunk): its Get() returns the concrete *TLValidationResult (Java's covariant
// TLValidatorTask#get() override), not job.ValidationResult, so *TLValidatorTask itself does not
// satisfy job.ValidationTask (Go has no covariant interface-method return types). *TLValidationResult
// otherwise already exposes every job.ValidationResult accessor (Indication/SubIndication/
// SigningTime/SigningCertificate/PotentialSigners; CachedResult is the empty interface), so no
// data is missing - only the wrapping return type. Rather than editing that frozen file, this
// adapter (in this file, which IS in-manifest) narrows Get()'s second return value to the
// interface at the one call site that needs it.
type tlValidatorTaskAdapter struct {
	task *TLValidatorTask
}

func (t tlValidatorTaskAdapter) Get() (job.ValidationResult, error) {
	return t.task.Get()
}

// parsingTaskAdapter narrows a concrete Get() (R, error) closure to job.ParsingTask's Get()
// (job.ParsingResult, error). Used by TLAnalysis/LOTLAnalysis's GetParsingTask to wrap
// *TLParsingTask/*LOTLParsingTask (tsl/tl_parsing_task.go, tsl/lotl_parsing_task.go - both out of
// this manifest, ported by TSLCORE), whose Get() likewise returns the covariant concrete
// *TLParsingResult/*LOTLParsingResult rather than job.ParsingResult - the same
// no-covariant-interface-return-types gap as tlValidatorTaskAdapter above.
type parsingTaskAdapter[R job.ParsingResult] struct {
	get func() (R, error)
}

func (p parsingTaskAdapter[R]) Get() (job.ParsingResult, error) {
	return p.get()
}

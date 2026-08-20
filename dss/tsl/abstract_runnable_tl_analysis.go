// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/AbstractRunnableTLAnalysis.java (DSS 6.5.RC1).
package tsl

import (
	"sync"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/client/http"
	"github.com/utain/esig/dss/validation/job"
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
		AbstractRunnableAnalysis: job.NewAbstractRunnableAnalysis(&tlSource.DocumentSource, cacheAccess, dssFileLoader, latch),
	}
}

// GetDownloadTask ports the protected getDownloadTask(DSSFileLoader, String) override.
func (a *AbstractRunnableTLAnalysis) GetDownloadTask(dssFileLoader http.DSSFileLoader, url string) job.DownloadTask {
	return NewXmlDownloadTask(dssFileLoader, url)
}

// GetValidationTask ports the protected getValidationTask(DSSDocument, CertificateSource)
// override.
func (a *AbstractRunnableTLAnalysis) GetValidationTask(document model.DSSDocument, certificateSource spi.CertificateSource) job.ValidationTask {
	return NewTLValidatorTask(document, certificateSource)
}

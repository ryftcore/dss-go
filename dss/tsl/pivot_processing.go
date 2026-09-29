// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/runnable/PivotProcessing.java (DSS 6.5.RC1).
//
// DEVIATION: Java extends validation.job.runnable.AbstractAnalysis to reuse its
// download()/parsing()/expireCache() bodies (protected methods a Java subclass reaches through
// inheritance). dss/validation/job.AbstractAnalysis keeps those same methods
// unexported (package-private to Go's "job" package, matching their Java "protected" intent
// against every OTHER package), so a type in package tsl cannot call them through embedding.
// The small amount of logic those private methods perform - the same shape as
// AbstractAnalysis's own download()/parsing()/expireCache() bodies - is therefore inlined
// directly against the *TLCacheAccessByKey this type already receives, rather than through the
// base.
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PivotProcessing processes a pivot analysis.
type PivotProcessing struct {
	// pivotSource is the pivot source being analyzed.
	pivotSource *LOTLSource

	// pivotCacheAccess is the cache access of the current pivot to process.
	pivotCacheAccess *TLCacheAccessByKey

	// lotlCacheAccess is the cache access of the LOTL.
	lotlCacheAccess *TLCacheAccessByKey

	// preceedingPivotCacheAccessByKeyList is a list of other pivots, to be updated in case of
	// current pivot update.
	preceedingPivotCacheAccessByKeyList []*TLCacheAccessByKey

	// dssFileLoader is the file loader.
	dssFileLoader http.DSSFileLoader
}

// NewPivotProcessing is the default constructor.
func NewPivotProcessing(pivotSource *LOTLSource, pivotCacheAccess *TLCacheAccessByKey, lotlCacheAccess *TLCacheAccessByKey,
	preceedingPivotCacheAccessByKeyList []*TLCacheAccessByKey, dssFileLoader http.DSSFileLoader) *PivotProcessing {
	return &PivotProcessing{
		pivotSource:                         pivotSource,
		pivotCacheAccess:                    pivotCacheAccess,
		lotlCacheAccess:                     lotlCacheAccess,
		preceedingPivotCacheAccessByKeyList: preceedingPivotCacheAccessByKeyList,
		dssFileLoader:                       dssFileLoader,
	}
}

// Call ports call(). Java declares `throws Exception`; this returns error instead - though,
// matching Java's own behaviour (every internal error is routed to the corresponding
// cacheAccess.*Error and swallowed, never rethrown from call()), the error return here is always
// nil; a caller only needs the *PivotProcessingResult (nil when no XML LOTL pointer is found).
func (p *PivotProcessing) Call() (*PivotProcessingResult, error) {
	pivot := p.download(p.pivotSource.Url())
	if pivot != nil {
		p.parsing(pivot)

		parsingResult, _ := p.pivotCacheAccess.GetParsingReadOnlyResult().(*TLParsingCacheDTO)
		xmlLotlPointer := ParsingUtilsXMLLOTLPointer(parsingResult)
		if xmlLotlPointer != nil {
			return NewPivotProcessingResult(pivot, ParsingUtilsLOTLAnnouncedCertificateSource(xmlLotlPointer), xmlLotlPointer.TSLLocation()), nil
		}
	}
	return nil, nil
}

// catchException runs fn and returns the error it returns or, when it panics, an error built from
// the panic value: the stand-in for the catch (Exception e) of the inherited
// AbstractAnalysis.download()/parsing(), which records the exception on the cache instead of
// letting it escape (the same helper job.AbstractAnalysis uses, unexported there).
//
// This matters more here than on the main analysis path: a pivot runs on its own goroutine
// (see LOTLWithPivotsAnalysis.downloadAndParseAllPivots) with no recover of the analysis
// runnable's around it, so an unrecovered panic - e.g. the nil dereference LOTLParsingTask
// reproduces from upstream's NullPointerException for a pivot lacking <SchemeInformation> -
// would take down the whole process, where Java records a parsing error on the pivot's cache
// entry and carries on.
func catchException(fn func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = recoveredErr
			} else {
				err = model.NewDSSError(fmt.Sprint(recovered))
			}
		}
	}()
	return fn()
}

// download ports the inherited download(String), see this file's header.
func (p *PivotProcessing) download(url string) model.DSSDocument {
	var document model.DSSDocument
	if err := catchException(func() error {
		downloadTask := NewXmlDownloadTask(p.dssFileLoader, url)
		downloadResult, err := downloadTask.Get()
		if err != nil {
			return err
		}
		if !p.pivotCacheAccess.IsUpToDate(downloadResult) {
			p.pivotCacheAccess.UpdateDownloadResult(downloadResult)
			p.expireCache()
		}
		document = downloadResult.DSSDocument()
		return nil
	}); err != nil {
		p.pivotCacheAccess.DownloadError(err)
		return nil
	}
	return document
}

// parsing ports the inherited parsing(DSSDocument), see this file's header.
func (p *PivotProcessing) parsing(document model.DSSDocument) {
	if p.pivotCacheAccess.IsParsingRefreshNeeded() {
		if err := catchException(func() error {
			parsingTask := NewLOTLParsingTask(document, p.pivotSource)
			parsingResult, err := parsingTask.Get()
			if err != nil {
				return err
			}
			p.pivotCacheAccess.UpdateParsingResult(parsingResult)
			return nil
		}); err != nil {
			p.pivotCacheAccess.ParsingError(err)
		}
	}
}

// expireCache ports the protected expireCache() override.
func (p *PivotProcessing) expireCache() {
	p.pivotCacheAccess.ExpireParsing()
	p.pivotCacheAccess.ExpireValidation()
	p.lotlCacheAccess.ExpireValidation() // ensure LOTL will be updated in case of pivot refresh
	// expire Pivots before the current
	if utils.IsCollectionNotEmpty(p.preceedingPivotCacheAccessByKeyList) {
		for _, pivotCacheAccessByKey := range p.preceedingPivotCacheAccessByKeyList {
			pivotCacheAccessByKey.ExpireValidation()
		}
	}
}

// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/timestamp/PdfTimestampToken.java
// (DSS 6.5.RC1).
//
// Java's PdfTimestampToken extends validation.TimestampToken (single inheritance), adding one
// field (pdfRevision) and overriding getTimestampIdentifierBuilder(). Go has no inheritance, and
// callers elsewhere in this package need to recover a PdfTimestampToken from a bare
// *validation.TimestampToken they hold (native_pdf_signature_service.go's
// AnalyzeAllModifications/AnalyzeRevisionModifications, an "instanceof PdfTimestampToken" check
// in the Java), which a struct wrapping (rather than being-a) *validation.TimestampToken cannot
// support through a type assertion alone. This file therefore keeps a process-wide registry
// mapping the wrapped *validation.TimestampToken back to its enclosing *PdfTimestampToken,
// populated by NewPdfTimestampToken and queried by PdfTimestampTokenOf - both names
// native_pdf_signature_service.go and pades_timestamp_scope_finder.go already call, per this
// chunk's own manifest.
//
// The registry must not keep a token, and through it the whole PDF it was read from, alive for
// the life of the process, which Java's `instanceof` (no side table at all) never does. Both
// sides of the mapping are therefore weak, and the wrapper is kept alive by the very
// *validation.TimestampToken it is looked up from: the token holds its identifier builder, and the
// builder holds the wrapper (PdfTimestampTokenIdentifierBuilder.token). A registry entry thus
// lives exactly as long as its token, and runtime.AddCleanup removes it once the token is
// unreachable.
package pades

import (
	"runtime"
	"sync"
	"weak"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// PdfTimestampToken is a specific TimestampToken for a PDF document time-stamp. Port of the
// class PdfTimestampToken, extending validation.TimestampToken.
type PdfTimestampToken struct {
	*validation.TimestampToken

	// pdfRevision is the related PDF revision.
	pdfRevision *PdfDocTimestampRevision
}

// pdfTimestampTokenRegistryMu guards pdfTimestampTokenRegistry.
var pdfTimestampTokenRegistryMu sync.Mutex

// pdfTimestampTokenRegistry maps the wrapped *validation.TimestampToken back to the enclosing
// *PdfTimestampToken; see this file's header for why this registry exists and why both sides are
// weak.
var pdfTimestampTokenRegistry = map[weak.Pointer[validation.TimestampToken]]weak.Pointer[PdfTimestampToken]{}

// pdfTimestampTokenRegistryForget drops the entry of a token that has been garbage collected.
// It is a plain function, not a closure, so that it cannot keep the token it is attached to alive.
func pdfTimestampTokenRegistryForget(key weak.Pointer[validation.TimestampToken]) {
	pdfTimestampTokenRegistryMu.Lock()
	delete(pdfTimestampTokenRegistry, key)
	pdfTimestampTokenRegistryMu.Unlock()
}

// NewPdfTimestampToken is the default constructor. Port of the
// PdfTimestampToken(PdfDocTimestampRevision) constructor.
//
// TSPException/IOException/CMSException are folded into a returned error, following this
// port's usual checked-exception convention.
func NewPdfTimestampToken(pdfTimestampRevision *PdfDocTimestampRevision) (*PdfTimestampToken, error) {
	// TODO : refactor TimestampToken to init with CMS
	encoded := pdfTimestampRevision.PdfSigDictInfo().CMS().DEREncoded()
	// See this file's header: the identifier builder cannot be built from the not-yet-existing
	// PdfTimestampToken (Java passes `this`), so it is built from the two pieces of data it
	// actually reads out of one instead - see pdf_timestamp_token_identifier_builder.go.
	identifierBuilder := NewPdfTimestampTokenIdentifierBuilder(encoded, pdfTimestampRevision)
	base, err := validation.NewTimestampTokenWithIdentifierBuilder(encoded, enumerations.TimestampTypeDocumentTimestamp,
		[]*validation.TimestampedReference{}, identifierBuilder)
	if err != nil {
		return nil, err
	}
	token := &PdfTimestampToken{TimestampToken: base, pdfRevision: pdfTimestampRevision}
	// base -> its identifier builder -> token: the wrapper lives as long as the token it wraps.
	identifierBuilder.token = token

	key := weak.Make(base)
	pdfTimestampTokenRegistryMu.Lock()
	pdfTimestampTokenRegistry[key] = weak.Make(token)
	pdfTimestampTokenRegistryMu.Unlock()
	runtime.AddCleanup(base, pdfTimestampTokenRegistryForget, key)

	return token, nil
}

// PdfRevision returns the current PDF timestamp revision. Port of getPdfRevision().
func (t *PdfTimestampToken) PdfRevision() *PdfDocTimestampRevision {
	return t.pdfRevision
}

// PdfTimestampTokenOf recovers the PdfTimestampToken wrapping timestampToken, when timestampToken
// was built through NewPdfTimestampToken. It stands in for Java's
// `timestampToken instanceof PdfTimestampToken` checks; see this file's header.
func PdfTimestampTokenOf(timestampToken *validation.TimestampToken) (*PdfTimestampToken, bool) {
	if timestampToken == nil {
		return nil, false
	}
	pdfTimestampTokenRegistryMu.Lock()
	defer pdfTimestampTokenRegistryMu.Unlock()
	weakToken, ok := pdfTimestampTokenRegistry[weak.Make(timestampToken)]
	if !ok {
		return nil, false
	}
	pdfTimestampToken := weakToken.Value()
	return pdfTimestampToken, pdfTimestampToken != nil
}

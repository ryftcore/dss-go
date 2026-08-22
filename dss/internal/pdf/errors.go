// Error taxonomy for internal/pdf. See DESIGN.md §4.2 (frozen with object.go).
//
// Provenance: the sentinel set is the Go shape of the exceptions pdfbox 3.0.7
// throws out of Loader.loadPDF / COSParser and that dss-pades-pdfbox maps onto
// InvalidPasswordException / ProtectedDocumentException / IllegalInputException.

package pdf

import "errors"

var (
	// ErrNotPDF is returned when no %PDF- marker is found in the first 1024 bytes.
	ErrNotPDF = errors.New("pdf: missing %PDF- header")
	// ErrBrokenCatalog is the port of pdfbox's
	// IOException("Page tree root must be a dictionary"): the one defect pdfbox
	// does not recover from, and neither do we.
	ErrBrokenCatalog = errors.New("pdf: page tree root must be a dictionary")
	// ErrInvalidPassword maps onto upstream InvalidPasswordException.
	ErrInvalidPassword = errors.New("pdf: invalid password")
	// ErrUnsupportedSecurityHandler is returned for a security handler this
	// package does not implement: a /Filter other than /Standard, a /V outside
	// {1, 2, 4, 5}, a /V 5 paired with an /R other than 5 or 6, a crypt filter
	// a /V 4 //V 5 document selects (/StmF or /StrF, not /Identity) that does
	// not resolve through /CF to an implemented /CFM, or that disagrees with
	// the other selected filter (DESIGN.md §2.6, crypt.go's
	// checkCryptFilters), or a /UE //OE that is not exactly 32 bytes for
	// /R 5//R 6 (crypt.go's computeEncryptionKey).
	ErrUnsupportedSecurityHandler = errors.New("pdf: unsupported security handler")
	// ErrUnsupportedFilter is what a *FilterError unwraps to.
	ErrUnsupportedFilter = errors.New("pdf: unsupported stream filter")
	// ErrLimitExceeded reports a resource guard (MaxObjects, MaxDepth,
	// MaxStreamSize). These are ours, not pdfbox's; see DESIGN.md §2.7.
	ErrLimitExceeded = errors.New("pdf: resource limit exceeded")
	// ErrNoSuchObject is returned by Document.Object for an absent key. Note that
	// Resolve never returns it: a dangling reference resolves to Null{} (§2.7 O3).
	ErrNoSuchObject = errors.New("pdf: object not found")
	// ErrSignatureAlreadyAdded ports PDDocument.addSignature's IllegalStateException.
	ErrSignatureAlreadyAdded = errors.New("pdf: only one signature may be added per increment")
	// ErrContentsTooLarge is returned when the CMS does not fit the reserved space.
	ErrContentsTooLarge = errors.New("pdf: CMS does not fit the reserved /Contents space")
	// ErrByteRangeTooLarge is returned when the formatted /ByteRange exceeds the
	// 35 reserved bytes (R18).
	ErrByteRangeTooLarge = errors.New("pdf: /ByteRange does not fit the reserved space")
)

// FilterError names the filter that could not be decoded. It unwraps to
// ErrUnsupportedFilter, so errors.Is(err, ErrUnsupportedFilter) holds.
type FilterError struct{ Filter Name }

func (e *FilterError) Error() string {
	return "pdf: unsupported stream filter /" + string(e.Filter)
}

func (e *FilterError) Unwrap() error { return ErrUnsupportedFilter }

// WarningCode classifies a recovered defect. Codes are stable across releases;
// the oracle compares codes, never message text.
type WarningCode string

const (
	WarnHeaderGarbage      WarningCode = "header-garbage"
	WarnHeaderVersion      WarningCode = "header-version-default"
	WarnMissingEOF         WarningCode = "missing-eof"
	WarnXRefOffsetRepaired WarningCode = "xref-offset-repaired"
	WarnXRefBruteForce     WarningCode = "xref-brute-force"
	WarnXRefEntryInvalid   WarningCode = "xref-entry-invalid"
	WarnXRefStmSkipped     WarningCode = "xrefstm-skipped"
	WarnObjectHeaderFixed  WarningCode = "object-header-fixed"
	WarnDanglingReference  WarningCode = "dangling-reference"
	WarnStreamLengthFixed  WarningCode = "stream-length-fixed"
	WarnStreamEndFixed     WarningCode = "stream-end-fixed"
	WarnFlateTruncated     WarningCode = "flate-truncated"
	WarnPredictorTruncated WarningCode = "predictor-truncated"
	WarnObjStmBroken       WarningCode = "objstm-broken"
	WarnKidRemoved         WarningCode = "kid-removed"
	// WarnNumberClamped records an integer literal that exceeded int64 (§2.1).
	WarnNumberClamped WarningCode = "number-clamped"
	// WarnLexer records a tolerated lexical defect (bad hex digit, malformed #XX).
	WarnLexer WarningCode = "lexer"
)

// Warning is one recovered defect. Offset is the byte offset in the source file
// where the defect was noticed, or -1 when it is not tied to a position.
type Warning struct {
	Code    WarningCode
	Offset  int64
	Message string
}

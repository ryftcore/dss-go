// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampAnalyzer.java (DSS 6.5.RC1).
//
// # Virtual dispatch
//
// Java's DetachedTimestampAnalyzer extends DefaultDocumentAnalyzer and additionally declares two
// of its own protected (overridable) methods with concrete default bodies: createTimestampToken()
// and getTimestampScopes(TimestampToken). ASiCWithCAdESTimestampAnalyzer (a later phase) overrides
// both, plus DefaultDocumentAnalyzerOverrides.AddReference/IsTimestampCoveredByEvidenceRecord.
// Following the InitChainItem/Init<Base> virtual-dispatch convention this porting effort uses
// pervasively (see analyzer.DefaultDocumentAnalyzer's own file header), DetachedTimestampAnalyzerOverrides
// embeds analyzer.DefaultDocumentAnalyzerOverrides so a single leaf concrete type (this type
// itself, or a subclass embedding it from another package) satisfies both interfaces at once and
// InitDetachedTimestampAnalyzer can register the same value with both dispatch points in one call.
//
// # Protected-superclass-method gap
//
// getTimestamp()'s Java body calls two DefaultDocumentAnalyzer methods this port declared
// unexported because, at the time analyzer.DefaultDocumentAnalyzer was written, no manifest file
// then in scope needed cross-package access to them (getTimestampedReferences(List) and
// appendExternalEvidenceRecords(TimestampToken) are `protected` in Java - reachable by any
// subclass - but ported as analyzer.DefaultDocumentAnalyzer's unexported
// getTimestampedReferences/appendExternalEvidenceRecordsToTimestamp, inaccessible from this
// different package). Per PORTING.md ("no edits to frozen packages beyond the UNGATE manifests"),
// analyzer.DefaultDocumentAnalyzer is not edited; both bodies are instead reproduced verbatim
// below using only analyzer.DefaultDocumentAnalyzer's already-exported surface
// (DetachedEvidenceRecords()) plus this type's own overrides field (which - per the embedding
// above - already exposes AddReference/IsTimestampCoveredByEvidenceRecord). Flagged as a
// cross-chunk/integration note: a future pass could promote these two methods to exported
// accessors on analyzer.DefaultDocumentAnalyzer to remove the duplication.
package timestamp

import (
	"fmt"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
	modelscope "github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
	analyzertimestamp "github.com/utain/esig/dss/spi/validation/analyzer/timestamp"
	spiscope "github.com/utain/esig/dss/spi/validation/scope"
	timestampsrc "github.com/utain/esig/dss/spi/validation/timestamp"
)

// DetachedTimestampAnalyzerOverrides declares the two additional operations
// DetachedTimestampAnalyzer calls back into virtually, on top of every
// analyzer.DefaultDocumentAnalyzerOverrides operation. See the file header.
type DetachedTimestampAnalyzerOverrides interface {
	analyzer.DefaultDocumentAnalyzerOverrides

	// CreateTimestampToken creates a timestamp token from the validating document. Default: see
	// DetachedTimestampAnalyzer.CreateTimestampToken. Port of createTimestampToken().
	//
	// Java throws an unchecked DSSException on a malformed timestamp; ported as a returned error
	// (data-dependent on the validating document's bytes), consumed by Timestamp() below - whose
	// own signature, fixed by the already-frozen timestamp.TimestampAnalyzer interface this type
	// implements, has no error return, so Timestamp() panics on failure (a forced deviation; see
	// its own doc comment).
	CreateTimestampToken() (*validation.TimestampToken, error)

	// GetTimestampScopes finds timestamp scopes. Default: see
	// DetachedTimestampAnalyzer.GetTimestampScopes. Port of getTimestampScopes(TimestampToken).
	GetTimestampScopes(timestampToken *validation.TimestampToken) []modelscope.SignatureScope
}

// DetachedTimestampAnalyzer performs a processing of a detached timestamp document. Embedded by
// concrete subclasses (in this or other packages), which must call
// InitDetachedTimestampAnalyzer once (typically from their own constructor) - see the file
// header.
type DetachedTimestampAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// overrides points back at the concrete analyzer; see InitDetachedTimestampAnalyzer.
	overrides DetachedTimestampAnalyzerOverrides

	// timestampType is the type of the timestamp.
	timestampType enumerations.TimestampType

	// timestampToken is the TimestampToken to be validated, cached on first access.
	timestampToken *validation.TimestampToken
}

// compile-time interface assertions.
var (
	_ analyzer.DocumentAnalyzer                 = (*DetachedTimestampAnalyzer)(nil)
	_ analyzertimestamp.TimestampAnalyzer       = (*DetachedTimestampAnalyzer)(nil)
	_ DetachedTimestampAnalyzerOverrides        = (*DetachedTimestampAnalyzer)(nil)
	_ analyzer.DefaultDocumentAnalyzerOverrides = (*DetachedTimestampAnalyzer)(nil)
)

// NewDetachedTimestampAnalyzerBase builds the base state a subclass embeds. Port of the
// package-private empty DetachedTimestampAnalyzer() constructor; the subclass constructor must
// follow it with InitDetachedTimestampAnalyzer.
func NewDetachedTimestampAnalyzerBase() DetachedTimestampAnalyzer {
	return DetachedTimestampAnalyzer{DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase()}
}

// InitDetachedTimestampAnalyzer registers the concrete analyzer with its base (and, through it,
// with the embedded analyzer.DefaultDocumentAnalyzer) so both can dispatch virtually. Every
// concrete subclass constructor must call this once.
func (a *DetachedTimestampAnalyzer) InitDetachedTimestampAnalyzer(overrides DetachedTimestampAnalyzerOverrides) {
	a.overrides = overrides
	a.InitDefaultDocumentAnalyzer(overrides)
}

// detachedTimestampAnalyzerOverrides returns the registered overrides, panicking when the
// concrete analyzer forgot to call InitDetachedTimestampAnalyzer.
func (a *DetachedTimestampAnalyzer) detachedTimestampAnalyzerOverrides() DetachedTimestampAnalyzerOverrides {
	if a.overrides == nil {
		panic("DetachedTimestampAnalyzer was not initialised: the concrete analyzer must call InitDetachedTimestampAnalyzer in its constructor")
	}
	return a.overrides
}

// newDetachedTimestampAnalyzer wires the overrides registration shared by both exported
// constructors.
func newDetachedTimestampAnalyzer() *DetachedTimestampAnalyzer {
	a := &DetachedTimestampAnalyzer{DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase()}
	a.InitDetachedTimestampAnalyzer(a)
	return a
}

// NewDetachedTimestampAnalyzer is the default constructor. Port of
// DetachedTimestampAnalyzer(DSSDocument).
func NewDetachedTimestampAnalyzer(timestampFile model.DSSDocument) *DetachedTimestampAnalyzer {
	return NewDetachedTimestampAnalyzerWithType(timestampFile, enumerations.TimestampType_CONTENT_TIMESTAMP)
}

// NewDetachedTimestampAnalyzerWithType is the default constructor with a type. Port of
// DetachedTimestampAnalyzer(DSSDocument, TimestampType).
func NewDetachedTimestampAnalyzerWithType(timestampFile model.DSSDocument, timestampType enumerations.TimestampType) *DetachedTimestampAnalyzer {
	a := newDetachedTimestampAnalyzer()
	a.SetDocument(timestampFile)
	a.timestampType = timestampType
	return a
}

// SetTimestampType sets the TimestampType. ADDITIVE accessor for subclasses in other packages
// following NewDetachedTimestampAnalyzerBase, mirroring Java field-assignment access a same-
// package/inherited subclass has for granted (see analyzer.DefaultDocumentAnalyzer's
// DetachedContents doc comment on this pattern).
func (a *DetachedTimestampAnalyzer) SetTimestampType(timestampType enumerations.TimestampType) {
	a.timestampType = timestampType
}

// IsSupported checks if the document is supported by the current validator. Port of
// isSupported(DSSDocument).
func (a *DetachedTimestampAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	firstByte, err := spi.DSSUtilsReadFirstByte(dssDocument)
	if err != nil {
		return false
	}
	if spi.DSSASN1UtilsIsASN1SequenceTag(firstByte) {
		return detachedTimestampAnalyzerIsTimestampToken(dssDocument)
	}
	return false
}

// detachedTimestampAnalyzerIsTimestampToken reports whether document parses as a CMS SignedData
// whose encapsulated content type is id-ct-TSTInfo. Port of DSSUtils.isTimestampToken(DSSDocument).
//
// NOT PORTED IN spi/dss_utils.go: mirrors cades.cmsDocumentAnalyzerIsTimestampToken's own
// identical porter note - completed locally per package (rather than in the frozen
// spi/dss_utils.go, which cannot import cms) since this is the second package independently
// needing it.
func detachedTimestampAnalyzerIsTimestampToken(document model.DSSDocument) bool {
	parsedCMS, err := cms.CMSUtilsParseToCMS(document)
	if err != nil {
		return false
	}
	return parsedCMS.SignedContentType().Equal(cmscore.OIDCTTSTInfo)
}

// BuildDetachedTimestamps builds a list of detached TimestampTokens extracted from the document.
// Port of buildDetachedTimestamps().
func (a *DetachedTimestampAnalyzer) BuildDetachedTimestamps() []*validation.TimestampToken {
	return []*validation.TimestampToken{a.Timestamp()}
}

// Timestamp returns a single TimestampToken to be validated, building and caching it on first
// access. Port of getTimestamp().
//
// Panics on a malformed timestamp document (see DetachedTimestampAnalyzerOverrides.
// CreateTimestampToken's doc comment on why this is a forced deviation from PORTING.md's
// data-dependent-throw-to-error rule: the interface this method satisfies,
// timestampsrc.TimestampAnalyzer, is already frozen with no error return).
func (a *DetachedTimestampAnalyzer) Timestamp() *validation.TimestampToken {
	if a.timestampToken == nil {
		overrides := a.detachedTimestampAnalyzerOverrides()

		timestampToken, err := overrides.CreateTimestampToken()
		if err != nil {
			panic(err.Error())
		}

		timestampScopes := overrides.GetTimestampScopes(timestampToken)
		timestampToken.SetTimestampScopes(overrides.GetTimestampScopes(timestampToken))
		timestampToken.SetTimestampedReferences(append(timestampToken.TimestampedReferences(),
			a.getTimestampedReferences(timestampScopes)...))
		a.appendExternalEvidenceRecordsToTimestamp(timestampToken)

		a.timestampToken = timestampToken
	}
	return a.timestampToken
}

// CreateTimestampToken creates a timestamp token from the validating document. Port of
// createTimestampToken(); this is DetachedTimestampAnalyzerOverrides' default body.
//
// Panics when certificateVerifier, document or timestampType is missing (Java's
// Objects.requireNonNull calls); returns an error for a malformed timestamp document (Java's
// caught CMSException/TSPException/IOException, rethrown as a DSSException - data-dependent on
// the validating document's bytes, so ported as an error per PORTING.md).
func (a *DetachedTimestampAnalyzer) CreateTimestampToken() (*validation.TimestampToken, error) {
	if a.CertificateVerifier() == nil {
		panic("CertificateVerifier is not defined")
	}
	if !a.HasDocument() {
		panic("The timestampFile must be defined!")
	}
	if a.timestampType == "" {
		panic("The TimestampType must be defined!")
	}

	documentBytes, err := spi.DSSUtilsToByteArrayOfDocument(a.Document())
	if err != nil {
		return nil, fmt.Errorf("unable to create a TimestampToken. Reason : %s", err.Error())
	}
	newTimestampToken, err := validation.NewTimestampToken(documentBytes, a.timestampType)
	if err != nil {
		return nil, fmt.Errorf("unable to create a TimestampToken. Reason : %s", err.Error())
	}
	newTimestampToken.SetFilename(a.Document().Name())
	if _, err := newTimestampToken.MatchDataDocument(a.TimestampedData()); err != nil {
		return nil, fmt.Errorf("unable to create a TimestampToken. Reason : %s", err.Error())
	}
	return newTimestampToken, nil
}

// SetTimestampedData sets the data that has been timestamped. Port of
// setTimestampedData(DSSDocument).
//
// Panics when document is nil (Java's Objects.requireNonNull("The document is null")).
func (a *DetachedTimestampAnalyzer) SetTimestampedData(document model.DSSDocument) {
	if document == nil {
		panic("The document is null")
	}
	a.SetDetachedContents([]model.DSSDocument{document})
}

// TimestampedData returns the timestamped data. Port of getTimestampedData().
//
// Panics when more than one detached document was provided (Java's thrown
// IllegalArgumentException("Only one detached document shall be provided for a timestamp
// validation!")).
func (a *DetachedTimestampAnalyzer) TimestampedData() model.DSSDocument {
	detachedContents := a.DetachedContents()
	if len(detachedContents) == 0 {
		return nil
	} else if len(detachedContents) > 1 {
		panic("Only one detached document shall be provided for a timestamp validation!")
	}
	return detachedContents[0]
}

// GetTimestampScopes finds timestamp scopes. Port of getTimestampScopes(TimestampToken); this is
// DetachedTimestampAnalyzerOverrides' default body.
func (a *DetachedTimestampAnalyzer) GetTimestampScopes(timestampToken *validation.TimestampToken) []modelscope.SignatureScope {
	timestampScopeFinder := spiscope.NewDetachedTimestampScopeFinder()
	timestampScopeFinder.SetTimestampedData(a.TimestampedData())
	return timestampScopeFinder.FindTimestampScope(timestampToken)
}

// OriginalDocuments returns the signed document(s) without their signature(s), given a
// signature's DSS ID. Port of the getOriginalDocuments(String) override, which always throws -
// an incomplete stub in Java itself (a bare, message-less `throw new
// UnsupportedOperationException();` behind a "TODO : add extraction of original documents"
// comment, unlike the descriptive-message override two lines below it). Shadows the promoted
// analyzer.DefaultDocumentAnalyzer.OriginalDocuments (see that type's file header on shadowing
// a top-level DocumentAnalyzer method this way).
func (a *DetachedTimestampAnalyzer) OriginalDocuments(signatureId string) []model.DSSDocument {
	panic("unsupported operation")
}

// OriginalDocumentsForSignature is DetachedTimestampAnalyzerOverrides' (inherited from
// analyzer.DefaultDocumentAnalyzerOverrides) required implementation, always throwing. Port of
// the getOriginalDocuments(AdvancedSignature) override.
func (a *DetachedTimestampAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	panic("getOriginalDocuments(AdvancedSignature) is not supported for DetachedTimestampValidator!")
}

// getTimestampedReferences returns a list of timestamped references from the given list of
// SignatureScopes. Reproduces analyzer.DefaultDocumentAnalyzer's unexported, otherwise
// inaccessible getTimestampedReferences(List) - see the file header's "Protected-superclass-
// method gap" section.
func (a *DetachedTimestampAnalyzer) getTimestampedReferences(signatureScopes []modelscope.SignatureScope) []*validation.TimestampedReference {
	overrides := a.detachedTimestampAnalyzerOverrides()
	timestampedReferences := make([]*validation.TimestampedReference, 0, len(signatureScopes))
	for _, signatureScope := range signatureScopes {
		if overrides.AddReference(signatureScope) {
			timestampedReferences = append(timestampedReferences,
				validation.NewTimestampedReference(signatureScope.DSSID().AsXmlID(), enumerations.TimestampedObjectType_SIGNED_DATA))
		}
	}
	return timestampedReferences
}

// appendExternalEvidenceRecordsToTimestamp appends the detached evidence records covering
// timestampToken. Reproduces analyzer.DefaultDocumentAnalyzer's unexported, otherwise
// inaccessible appendExternalEvidenceRecords(TimestampToken) - see the file header's "Protected-
// superclass-method gap" section.
func (a *DetachedTimestampAnalyzer) appendExternalEvidenceRecordsToTimestamp(timestampToken *validation.TimestampToken) {
	overrides := a.detachedTimestampAnalyzerOverrides()
	detachedTimestampSource := timestampsrc.NewDetachedTimestampSourceWithTimestamp(timestampToken)
	for _, evidenceRecord := range a.DetachedEvidenceRecords() {
		if overrides.IsTimestampCoveredByEvidenceRecord(timestampToken, evidenceRecord) {
			timestampToken.AddDetachedEvidenceRecord(evidenceRecord)
			_ = detachedTimestampSource.AddExternalEvidenceRecord(evidenceRecord)
		}
	}
}

// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/SignatureTimestampSource.java (DSS 6.5.RC1).
//
// # FORWARD DEPENDENCY: SignatureAttribute / SignatureProperties
//
// Java's eu.europa.esig.dss.spi.validation.SignatureAttribute and .SignatureProperties<SA> are
// flattened into the sibling dss/spi/validation package by another chunk of phase 2b (the "SIG"
// chunk, which also owns AdvancedSignature - already landed - and the concrete format-specific
// SignatureAttribute implementations of later phases). Their assumed shapes, inferred from every
// call this manifest makes on them (this file and signature_timestamp_identifier_builder.go),
// are:
//
//	type SignatureAttribute interface {
//		Identifier() model.Identifier // getIdentifier(); Java returns the narrower
//		                              // SignatureAttributeIdentifier (spi.validation.identifier,
//		                              // a 1:1 sibling package outside this manifest's scope);
//		                              // only .asXmlId() is ever called on it here, so the Go
//		                              // shape is narrowed to the minimum this manifest needs.
//	}
//
//	type SignatureProperties[SA SignatureAttribute] interface {
//		IsExist() bool  // isExist()
//		Attributes() []SA // getAttributes()
//	}
//
// Java's SignatureAttribute declares no equals() override, so its default is
// java.lang.Object's - reference identity - which getAttributeOrder below relies on; the SA type
// parameter therefore additionally carries Go's `comparable` constraint (every concrete
// SignatureAttribute in this codebase is expected to be a pointer type, for which `comparable`
// is pointer identity, matching Java exactly), rather than requiring an Equals method that
// Java's own contract does not provide.
//
// # FORWARD DEPENDENCY: EncapsulatedTimestampScopeFinder
//
// eu.europa.esig.dss.spi.validation.scope.EncapsulatedTimestampScopeFinder is a 1:1 sibling
// package (dss/spi/validation/scope, imported here as validationscope) outside this manifest.
// Its assumed shape, inferred from getTimestampScopes below (its only call site in this
// manifest), is:
//
//	type EncapsulatedTimestampScopeFinder struct{ ... }
//	func NewEncapsulatedTimestampScopeFinder() *EncapsulatedTimestampScopeFinder
//	func (f *EncapsulatedTimestampScopeFinder) SetSignature(signature validation.AdvancedSignature)
//	func (f *EncapsulatedTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []scope.SignatureScope
//
// # GAP flagged for integrator: TimestampToken.SetTimestampedReferences
//
// Throughout this file and abstract_timestamp_source.go, Java mutates the
// List<TimestampedReference> a TimestampToken.getTimestampedReferences() call already returns,
// relying on Java's List reference semantics (e.g. incorporateArchiveTimestampReferences,
// processExternalTimestamp). dss/spi/validation/timestamp_token.go (already ported, frozen,
// outside this manifest) exposes TimestampedReferences() but no mutator, and Go slices returned
// by value do not alias the field they came from. This file therefore calls
// TimestampToken.SetTimestampedReferences([]*validation.TimestampedReference), via this
// package's timestampAddReferences helper (abstract_timestamp_source.go), which does not exist
// yet in timestamp_token.go. Adding it is a one-method, additive, non-breaking change mirroring
// the SetManifestFile/SetFilename/... pattern already used there:
//
//	func (t *TimestampToken) SetTimestampedReferences(timestampedReferences []*TimestampedReference) {
//		t.timestampedReferences = timestampedReferences
//	}
//
// This chunk does not add it itself (timestamp_token.go is outside this manifest); flagged
// prominently in the porter report.
package timestamp

import (
	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/identifier"
	validationscope "github.com/utain/esig/dss/spi/validation/scope"
)

// SignatureTimestampSourceOverrides is the contract a concrete format-specific timestamp source
// (CAdES/XAdES/JAdES, ported in later phases) must implement. Port of every `protected abstract`
// method of SignatureTimestampSource; dispatched the way model.TokenBase dispatches to
// model.TokenOverrides via InitToken, following the precedent of every other Overrides interface
// in this port (TimestampIdentifierBuilderOverrides, SignatureCertificateSourceOverrides, ...).
//
// None of the make*/build* methods return an error: Java declares no throws clause on any of
// them, so a concrete implementation whose construction can fail is expected to panic, exactly
// as an unchecked DSSException would propagate through these non-throws-declared Java methods;
// see the `must` helper in abstract_timestamp_source.go for how this file does the same at its
// own fallible call sites.
type SignatureTimestampSourceOverrides[AS validation.AdvancedSignature, SA interface {
	validation.SignatureAttribute
	comparable
}] interface {
	// BuildSignedSignatureProperties creates the 'signed-signature-properties' element of the
	// signature. Port of the abstract buildSignedSignatureProperties().
	BuildSignedSignatureProperties() validation.SignatureProperties[SA]
	// BuildUnsignedSignatureProperties creates the 'unsigned-signature-properties' element of
	// the signature. Port of the abstract buildUnsignedSignatureProperties().
	BuildUnsignedSignatureProperties() validation.SignatureProperties[SA]

	// IsContentTimestamp determines if signedAttribute is a "content-timestamp" element.
	// NOTE: applicable only for CAdES. Port of the abstract isContentTimestamp(SA).
	IsContentTimestamp(signedAttribute SA) bool
	// IsAllDataObjectsTimestamp determines if signedAttribute is a "data-objects-timestamp"
	// element. NOTE: applicable only for XAdES. Port of the abstract isAllDataObjectsTimestamp(SA).
	IsAllDataObjectsTimestamp(signedAttribute SA) bool
	// IsIndividualDataObjectsTimestamp determines if signedAttribute is an
	// "individual-data-objects-timestamp" element. NOTE: applicable only for XAdES.
	// Port of the abstract isIndividualDataObjectsTimestamp(SA).
	IsIndividualDataObjectsTimestamp(signedAttribute SA) bool

	// IsSignatureTimestamp determines if unsignedAttribute is a "signature-timestamp" element.
	// Port of the abstract isSignatureTimestamp(SA).
	IsSignatureTimestamp(unsignedAttribute SA) bool
	// IsCompleteCertificateRef determines if unsignedAttribute is a "complete-certificate-ref"
	// element. Port of the abstract isCompleteCertificateRef(SA).
	IsCompleteCertificateRef(unsignedAttribute SA) bool
	// IsAttributeCertificateRef determines if unsignedAttribute is an "attribute-certificate-ref"
	// element. Port of the abstract isAttributeCertificateRef(SA).
	IsAttributeCertificateRef(unsignedAttribute SA) bool
	// IsCompleteRevocationRef determines if unsignedAttribute is a "complete-revocation-ref"
	// element. Port of the abstract isCompleteRevocationRef(SA).
	IsCompleteRevocationRef(unsignedAttribute SA) bool
	// IsAttributeRevocationRef determines if unsignedAttribute is an "attribute-revocation-ref"
	// element. Port of the abstract isAttributeRevocationRef(SA).
	IsAttributeRevocationRef(unsignedAttribute SA) bool
	// IsRefsOnlyTimestamp determines if unsignedAttribute is a "refs-only-timestamp" element.
	// Port of the abstract isRefsOnlyTimestamp(SA).
	IsRefsOnlyTimestamp(unsignedAttribute SA) bool
	// IsSigAndRefsTimestamp determines if unsignedAttribute is a "sig-and-refs-timestamp"
	// element. Port of the abstract isSigAndRefsTimestamp(SA).
	IsSigAndRefsTimestamp(unsignedAttribute SA) bool
	// IsCertificateValues determines if unsignedAttribute is a "certificate-values" element.
	// Port of the abstract isCertificateValues(SA).
	IsCertificateValues(unsignedAttribute SA) bool
	// IsRevocationValues determines if unsignedAttribute is a "revocation-values" element.
	// Port of the abstract isRevocationValues(SA).
	IsRevocationValues(unsignedAttribute SA) bool
	// IsAttrAuthoritiesCertValues determines if unsignedAttribute is an
	// "AttrAuthoritiesCertValues" element. Port of the abstract isAttrAuthoritiesCertValues(SA).
	IsAttrAuthoritiesCertValues(unsignedAttribute SA) bool
	// IsAttributeRevocationValues determines if unsignedAttribute is an
	// "AttributeRevocationValues" element. Port of the abstract isAttributeRevocationValues(SA).
	IsAttributeRevocationValues(unsignedAttribute SA) bool
	// IsArchiveTimestamp determines if unsignedAttribute is an "archive-timestamp" element.
	// Port of the abstract isArchiveTimestamp(SA).
	IsArchiveTimestamp(unsignedAttribute SA) bool
	// IsTimeStampValidationData determines if unsignedAttribute is a "timestamp-validation-data"
	// element. Port of the abstract isTimeStampValidationData(SA).
	IsTimeStampValidationData(unsignedAttribute SA) bool
	// IsAnyValidationData determines if unsignedAttribute is an "any-validation-data" element.
	// Port of the abstract isAnyValidationData(SA).
	IsAnyValidationData(unsignedAttribute SA) bool
	// IsValidationDataReferences determines if unsignedAttribute is a "references" element.
	// Port of the abstract isValidationDataReferences(SA).
	IsValidationDataReferences(unsignedAttribute SA) bool
	// IsCounterSignature determines if unsignedAttribute is a "counter-signature" element.
	// Port of the abstract isCounterSignature(SA).
	IsCounterSignature(unsignedAttribute SA) bool
	// IsSignaturePolicyStore determines if unsignedAttribute is a "signature-policy-store"
	// element. Port of the abstract isSignaturePolicyStore(SA).
	IsSignaturePolicyStore(unsignedAttribute SA) bool
	// IsEvidenceRecord determines if unsignedAttribute is an "evidence-record" element.
	// Port of the abstract isEvidenceRecord(SA).
	IsEvidenceRecord(unsignedAttribute SA) bool

	// DocumentTimestamps returns a list of document timestamps. Port of the public, overridable
	// getDocumentTimestamps(), which the base answers with an empty list and only
	// pades.PAdESTimestampSource overrides (a PDF /DocTimeStamp revision is not reachable through
	// any CMS unsigned attribute, so no other format has one). Routed through this interface -
	// rather than called on the embedded base directly - because the base's own
	// AllTimestampsExceptLastArchiveTimestamp() below consults it, and Go embedding gives that
	// call static, not virtual, dispatch: without this the base would read its own hard-coded
	// empty list even for a PAdES signature whose /DocTimeStamps ARE its archive timestamps, and
	// AllTimestampsExceptLastArchiveTimestamp() would come back empty for every PAdES signature.
	// See the identical rationale on AllTimestamps below.
	DocumentTimestamps() []*validation.TimestampToken
	// AllTimestamps returns a list of all incorporated timestamps. Port of the public,
	// overridable getAllTimestamps(); pades.PAdESTimestampSource overrides it to append the
	// document and /VRI timestamps the base cannot know about. Routed through this interface for
	// the same reason DocumentTimestamps above is: the base's own TimestampCertificateSources(),
	// TimestampCRLSources(), TimestampOCSPSources(), getTimestampsCoveredByManifest() and
	// isTimestamped() all consult it, and Java reaches the concrete override from every one of
	// them through ordinary virtual dispatch.
	AllTimestamps() []*validation.TimestampToken

	// MakeTimestampToken creates a timestamp token from the provided signatureAttribute.
	// Port of the abstract makeTimestampToken(SA, TimestampType, List).
	MakeTimestampToken(signatureAttribute SA, timestampType enumerations.TimestampType,
		references []*validation.TimestampedReference) *validation.TimestampToken
	// MakeTimestampTokens creates the (possibly several) timestamp tokens carried by
	// signatureAttribute. Port of the protected, overridable makeTimestampTokens(SA, TimestampType,
	// List): most formats produce at most one token per attribute and get that behaviour for free
	// from the base's own MakeTimestampTokens below (promoted by embedding, wrapping
	// MakeTimestampToken's single result in a slice) without needing to implement this method
	// themselves; XAdES shadows it directly since one xades132:XAdESTimeStampType element can carry
	// more than one xades132:EncapsulatedTimeStamp. Routed through SignatureTimestampSourceOverrides
	// for the same virtual-dispatch reason as every other override in this interface: the base's
	// own internal populateTimestampTokens() flow calls s.overrides.MakeTimestampTokens(...), never
	// the concrete type directly, so a format-specific override is only reachable if it is wired in
	// here.
	MakeTimestampTokens(signatureAttribute SA, timestampType enumerations.TimestampType,
		references []*validation.TimestampedReference) []*validation.TimestampToken
	// MakeEvidenceRecords creates a list of evidence records from the provided signatureAttribute.
	// Port of the abstract makeEvidenceRecords(SA, List).
	MakeEvidenceRecords(signatureAttribute SA, references []*validation.TimestampedReference) []validation.EvidenceRecord

	// GetCertificateRefs returns a list of CertificateRefs from unsignedAttribute.
	// Port of the abstract getCertificateRefs(SA).
	GetCertificateRefs(unsignedAttribute SA) []*spi.CertificateRef
	// GetCRLRefs returns a list of CRL revocation refs from unsignedAttribute.
	// Port of the abstract getCRLRefs(SA).
	GetCRLRefs(unsignedAttribute SA) []*spi.CRLRef
	// GetOCSPRefs returns a list of OCSP revocation refs from unsignedAttribute.
	// Port of the abstract getOCSPRefs(SA).
	GetOCSPRefs(unsignedAttribute SA) []*spi.OCSPRef

	// GetEncapsulatedCertificateIdentifiers returns a list of Identifiers obtained from
	// unsignedAttribute. Port of the abstract getEncapsulatedCertificateIdentifiers(SA).
	GetEncapsulatedCertificateIdentifiers(unsignedAttribute SA) []model.Identifier
	// GetEncapsulatedCRLIdentifiers returns a list of CRLBinaries obtained from unsignedAttribute.
	// Port of the abstract getEncapsulatedCRLIdentifiers(SA).
	GetEncapsulatedCRLIdentifiers(unsignedAttribute SA) []*crlparser.CRLBinary
	// GetEncapsulatedOCSPIdentifiers returns a list of OCSPResponseBinaries obtained from
	// unsignedAttribute. Port of the abstract getEncapsulatedOCSPIdentifiers(SA).
	GetEncapsulatedOCSPIdentifiers(unsignedAttribute SA) []*spi.OCSPResponseBinary

	// GetCounterSignatures extracts counter signatures from unsignedAttribute.
	// Port of the abstract getCounterSignatures(SA).
	GetCounterSignatures(unsignedAttribute SA) []validation.AdvancedSignature

	// GetArchiveTimestampType returns the ArchiveTimestampType for unsignedAttribute.
	// Port of the abstract getArchiveTimestampType(SA).
	GetArchiveTimestampType(unsignedAttribute SA) enumerations.ArchiveTimestampType

	// GetTimestampMessageImprintDigestBuilderForAlgorithm returns a TimestampMessageDigestBuilder
	// to compute a message digest with the given DigestAlgorithm.
	// Port of the abstract getTimestampMessageImprintDigestBuilder(DigestAlgorithm).
	GetTimestampMessageImprintDigestBuilderForAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) TimestampMessageDigestBuilder
	// GetTimestampMessageImprintDigestBuilderForToken returns the related TimestampMessageDigestBuilder.
	// Port of the abstract getTimestampMessageImprintDigestBuilder(TimestampToken).
	GetTimestampMessageImprintDigestBuilderForToken(timestampToken *validation.TimestampToken) TimestampMessageDigestBuilder

	// IncorporateArchiveTimestampReferences incorporates all the timestamped references for the
	// given archive timestampToken. Port of the protected incorporateArchiveTimestampReferences
	// (TimestampToken, List). Java declares this method concrete, not abstract, and relies on
	// virtual dispatch for the base's own internal callers (incorporateArchiveTimestampReferencesForTokens)
	// to reach a format-specific override; routed through this interface for the same reason
	// every other override above is. SignatureTimestampSource itself provides the base's
	// original body as IncorporateArchiveTimestampReferences below, for an override with nothing
	// prior-based-format-specific to add to delegate back to.
	IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken, previousTimestamps []*validation.TimestampToken)

	// GetSignatureSignedDataReferences returns a list of all TimestampedReferences found into
	// the signature's signed data (e.g. CMS SignedData for CAdES). Port of the protected
	// getSignatureSignedDataReferences(), empty by default; concrete-not-abstract in Java, see
	// IncorporateArchiveTimestampReferences's comment for why it is routed through this
	// interface. SignatureTimestampSource provides the empty-by-default base body as
	// GetSignatureSignedDataReferences below.
	GetSignatureSignedDataReferences() []*validation.TimestampedReference

	// GetCounterSignatureReferences returns a list of references extracted from counterSignature.
	// Port of the protected getCounterSignatureReferences(AdvancedSignature); concrete-not-abstract
	// in Java, see IncorporateArchiveTimestampReferences's comment for why it is routed through
	// this interface. SignatureTimestampSource provides the base's original body as
	// GetCounterSignatureReferences below.
	GetCounterSignatureReferences(counterSignature validation.AdvancedSignature) []*validation.TimestampedReference
}

// SignatureTimestampSource is the timestamp source of a signature.
//
// AS is the AdvancedSignature implementation, SA the corresponding SignatureAttribute; port of
// the Java type parameters <AS extends AdvancedSignature, SA extends SignatureAttribute>.
type SignatureTimestampSource[AS validation.AdvancedSignature, SA interface {
	validation.SignatureAttribute
	comparable
}] struct {
	AbstractTimestampSource

	// overrides points back at the concrete signature timestamp source; see
	// InitSignatureTimestampSource.
	overrides SignatureTimestampSourceOverrides[AS, SA]

	// signature is being validated.
	signature AS

	// crlSource is the CRL revocation source containing merged data from signature and
	// timestamps.
	crlSource *spi.ListRevocationSource[revocation.CRL]

	// ocspSource is the OCSP revocation source containing merged data from signature and
	// timestamps.
	ocspSource *spi.ListRevocationSource[revocation.OCSP]

	// certificateSource contains merged data from signature and timestamps.
	certificateSource *spi.ListCertificateSource

	// contentTimestamps are the enclosed content timestamps; nil stands for "not computed yet",
	// distinguished from "computed, empty" by createAndValidate always assigning a non-nil
	// (possibly zero-length) slice.
	contentTimestamps []*validation.TimestampToken

	// signatureTimestamps are the enclosed signature timestamps.
	signatureTimestamps []*validation.TimestampToken

	// sigAndRefsTimestamps are the enclosed SignAndRefs timestamps.
	sigAndRefsTimestamps []*validation.TimestampToken

	// refsOnlyTimestamps are the enclosed RefsOnly timestamps.
	refsOnlyTimestamps []*validation.TimestampToken

	// archiveTimestamps is the list of enclosed archive signature timestamps.
	archiveTimestamps []*validation.TimestampToken

	// detachedTimestamps is the list of detached timestamp tokens (used in ASiC with CAdES).
	detachedTimestamps []*validation.TimestampToken

	// embeddedEvidenceRecords is the list of evidence records embedded to the signature document.
	embeddedEvidenceRecords []validation.EvidenceRecord

	// detachedEvidenceRecords is the list of evidence records detached from the signature
	// document.
	detachedEvidenceRecords []validation.EvidenceRecord

	// unsignedPropertiesReferences is the list of all TimestampedReferences extracted from a
	// signature.
	unsignedPropertiesReferences []*validation.TimestampedReference

	// signedSignatureProperties is a cached instance of Signed Signature Properties.
	signedSignatureProperties validation.SignatureProperties[SA]

	// unsignedSignatureProperties is a cached instance of Unsigned Signature Properties.
	unsignedSignatureProperties validation.SignatureProperties[SA]
}

// NewSignatureTimestampSourceBase builds the base state a subclass embeds. Port of the protected
// SignatureTimestampSource(AS) constructor; a concrete subclass constructor must follow it with
// InitSignatureTimestampSource.
//
// Panics with the Java message when signature is missing (Objects.requireNonNull).
func NewSignatureTimestampSourceBase[AS validation.AdvancedSignature, SA interface {
	validation.SignatureAttribute
	comparable
}](signature AS) SignatureTimestampSource[AS, SA] {
	if any(signature) == nil {
		panic("The signature cannot be null!")
	}
	return SignatureTimestampSource[AS, SA]{signature: signature}
}

// InitSignatureTimestampSource registers the concrete signature timestamp source with its base
// so that the base can dispatch to SignatureTimestampSourceOverrides. Every concrete subclass
// constructor must call this once.
func (s *SignatureTimestampSource[AS, SA]) InitSignatureTimestampSource(overrides SignatureTimestampSourceOverrides[AS, SA]) {
	s.overrides = overrides
}

// ContentTimestamps returns a list of incorporated content timestamps.
// Port of getContentTimestamps().
func (s *SignatureTimestampSource[AS, SA]) ContentTimestamps() []*validation.TimestampToken {
	if s.contentTimestamps == nil {
		s.createAndValidate()
	}
	return s.contentTimestamps
}

// SignatureTimestamps returns a list of incorporated signature timestamps.
// Port of getSignatureTimestamps().
func (s *SignatureTimestampSource[AS, SA]) SignatureTimestamps() []*validation.TimestampToken {
	if s.signatureTimestamps == nil {
		s.createAndValidate()
	}
	return s.signatureTimestamps
}

// TimestampsX1 returns a list of incorporated SigAndRefs timestamps. Port of getTimestampsX1().
func (s *SignatureTimestampSource[AS, SA]) TimestampsX1() []*validation.TimestampToken {
	if s.sigAndRefsTimestamps == nil {
		s.createAndValidate()
	}
	return s.sigAndRefsTimestamps
}

// TimestampsX2 returns a list of incorporated RefsOnly timestamps. Port of getTimestampsX2().
func (s *SignatureTimestampSource[AS, SA]) TimestampsX2() []*validation.TimestampToken {
	if s.refsOnlyTimestamps == nil {
		s.createAndValidate()
	}
	return s.refsOnlyTimestamps
}

// ArchiveTimestamps returns a list of incorporated archive timestamps.
// Port of getArchiveTimestamps().
func (s *SignatureTimestampSource[AS, SA]) ArchiveTimestamps() []*validation.TimestampToken {
	if s.archiveTimestamps == nil {
		s.createAndValidate()
	}
	return s.archiveTimestamps
}

// DocumentTimestamps returns an empty list: applicable only for PAdES.
// Port of getDocumentTimestamps().
func (s *SignatureTimestampSource[AS, SA]) DocumentTimestamps() []*validation.TimestampToken {
	return nil
}

// DetachedTimestamps returns a list of detached timestamps (used in ASiC with CAdES).
// Port of getDetachedTimestamps().
func (s *SignatureTimestampSource[AS, SA]) DetachedTimestamps() []*validation.TimestampToken {
	if s.detachedTimestamps == nil {
		s.createAndValidate()
	}
	return s.detachedTimestamps
}

// AllTimestamps returns a list of all incorporated timestamps. Port of getAllTimestamps().
func (s *SignatureTimestampSource[AS, SA]) AllTimestamps() []*validation.TimestampToken {
	var timestampTokens []*validation.TimestampToken
	timestampTokens = append(timestampTokens, s.ContentTimestamps()...)
	timestampTokens = append(timestampTokens, s.SignatureTimestamps()...)
	timestampTokens = append(timestampTokens, s.TimestampsX1()...)
	timestampTokens = append(timestampTokens, s.TimestampsX2()...)
	timestampTokens = append(timestampTokens, s.ArchiveTimestamps()...)
	timestampTokens = append(timestampTokens, s.DetachedTimestamps()...)
	return timestampTokens
}

// EmbeddedEvidenceRecords returns a list of evidence records embedded in a signature document.
// Port of getEmbeddedEvidenceRecords().
func (s *SignatureTimestampSource[AS, SA]) EmbeddedEvidenceRecords() []validation.EvidenceRecord {
	if s.embeddedEvidenceRecords == nil {
		s.createAndValidate()
	}
	return s.embeddedEvidenceRecords
}

// DetachedEvidenceRecords returns a list of evidence records detached from a signature document.
// Port of getDetachedEvidenceRecords().
func (s *SignatureTimestampSource[AS, SA]) DetachedEvidenceRecords() []validation.EvidenceRecord {
	if s.detachedEvidenceRecords == nil {
		s.createAndValidate()
	}
	return s.detachedEvidenceRecords
}

// AllEvidenceRecords returns a list of all evidence records associated with the signature.
// Port of getAllEvidenceRecords().
func (s *SignatureTimestampSource[AS, SA]) AllEvidenceRecords() []validation.EvidenceRecord {
	var evidenceRecords []validation.EvidenceRecord
	evidenceRecords = append(evidenceRecords, s.EmbeddedEvidenceRecords()...)
	evidenceRecords = append(evidenceRecords, s.DetachedEvidenceRecords()...)
	return evidenceRecords
}

// TimestampCertificateSources returns a merged ListCertificateSource of all embedded timestamp
// certificate sources. Port of getTimestampCertificateSources().
func (s *SignatureTimestampSource[AS, SA]) TimestampCertificateSources() *spi.ListCertificateSource {
	result := spi.NewListCertificateSource()
	for _, timestampToken := range s.overrides.AllTimestamps() {
		result.Add(timestampToken.CertificateSource())
	}
	return result
}

// TimestampCertificateSourcesExceptLastArchiveTimestamp returns a merged ListCertificateSource of
// all embedded timestamp certificate sources except the latest Archive Timestamp.
// Port of getTimestampCertificateSourcesExceptLastArchiveTimestamp().
func (s *SignatureTimestampSource[AS, SA]) TimestampCertificateSourcesExceptLastArchiveTimestamp() *spi.ListCertificateSource {
	result := spi.NewListCertificateSource()
	for _, timestampToken := range s.AllTimestampsExceptLastArchiveTimestamp() {
		result.Add(timestampToken.CertificateSource())
	}
	return result
}

// AllTimestampsExceptLastArchiveTimestamp returns a list of all TimestampTokens except the last
// archive timestamp. Port of getAllTimestampsExceptLastArchiveTimestamp().
func (s *SignatureTimestampSource[AS, SA]) AllTimestampsExceptLastArchiveTimestamp() []*validation.TimestampToken {
	var timestampTokens []*validation.TimestampToken
	timestampTokens = append(timestampTokens, s.ContentTimestamps()...)
	timestampTokens = append(timestampTokens, s.SignatureTimestamps()...)
	timestampTokens = append(timestampTokens, s.TimestampsX1()...)
	timestampTokens = append(timestampTokens, s.TimestampsX2()...)

	var allArchiveTimestamps []*validation.TimestampToken
	allArchiveTimestamps = append(allArchiveTimestamps, s.ArchiveTimestamps()...)
	allArchiveTimestamps = append(allArchiveTimestamps, s.overrides.DocumentTimestamps()...) // can be a document timestamp for PAdES
	allArchiveTimestamps = append(allArchiveTimestamps, s.DetachedTimestamps()...)           // can be a detached timestamp for ASiC with CAdES
	if len(allArchiveTimestamps) > 0 {
		if len(timestampTokens) > 0 || containsTimestampsCoveringOtherTimestamps(allArchiveTimestamps) {
			// exclude the last archive timestamp
			comparator := validation.NewTimestampTokenComparator()
			timestampTokenSliceSortStable(allArchiveTimestamps, comparator)
			for ii := 0; ii < len(allArchiveTimestamps)-1; ii++ {
				timestampTokens = append(timestampTokens, allArchiveTimestamps[ii])
			}
		} else {
			// add all timestamps for validation
			timestampTokens = append(timestampTokens, allArchiveTimestamps...)
		}
	}
	return timestampTokens
}

// timestampTokenSliceSortStable sorts timestampTokens ascending per comparator, in place.
// Port of allArchiveTimestamps.sort(new TimestampTokenComparator()) (Collections.sort is a
// stable merge sort in Java; a stable sort is used here for the same reason).
func timestampTokenSliceSortStable(timestampTokens []*validation.TimestampToken, comparator validation.TimestampTokenComparator) {
	for i := 1; i < len(timestampTokens); i++ {
		for j := i; j > 0 && comparator.Less(timestampTokens[j], timestampTokens[j-1]); j-- {
			timestampTokens[j-1], timestampTokens[j] = timestampTokens[j], timestampTokens[j-1]
		}
	}
}

// containsTimestampsCoveringOtherTimestamps ports the private
// containsTimestampsCoveringOtherTimestamps(List).
func containsTimestampsCoveringOtherTimestamps(timestampTokens []*validation.TimestampToken) bool {
	for _, timestampToken := range timestampTokens {
		for _, reference := range timestampToken.TimestampedReferences() {
			if reference.Category() == enumerations.TimestampedObjectType_TIMESTAMP {
				return true
			}
		}
	}
	return false
}

// TimestampCRLSources returns a merged ListRevocationSource of all embedded timestamp CRL
// sources. Port of getTimestampCRLSources().
func (s *SignatureTimestampSource[AS, SA]) TimestampCRLSources() *spi.ListRevocationSource[revocation.CRL] {
	result := spi.NewListRevocationSource[revocation.CRL]()
	for _, timestampToken := range s.overrides.AllTimestamps() {
		result.Add(timestampToken.CRLSource())
	}
	return result
}

// TimestampOCSPSources returns a merged ListRevocationSource of all embedded timestamp OCSP
// sources. Port of getTimestampOCSPSources().
func (s *SignatureTimestampSource[AS, SA]) TimestampOCSPSources() *spi.ListRevocationSource[revocation.OCSP] {
	result := spi.NewListRevocationSource[revocation.OCSP]()
	for _, timestampToken := range s.overrides.AllTimestamps() {
		result.Add(timestampToken.OCSPSource())
	}
	return result
}

// UnsignedPropertiesReferences returns a list of TimestampedReferences for all tokens embedded
// into unsigned properties of the signature. Port of getUnsignedPropertiesReferences().
func (s *SignatureTimestampSource[AS, SA]) UnsignedPropertiesReferences() []*validation.TimestampedReference {
	if s.unsignedPropertiesReferences == nil {
		s.createAndValidate()
	}
	return s.unsignedPropertiesReferences
}

// CertificateSource returns the merged certificate source built from the signature and its
// timestamps. Additive accessor (integration-time, Phase 3) for the `certificateSource` field
// Java's subclasses reach through protected field access; every format-specific
// SignatureTimestampSourceOverrides implementation needs it the same way
// SignerDataReferences/UnsignedPropertiesReferences already expose their own state.
func (s *SignatureTimestampSource[AS, SA]) CertificateSource() *spi.ListCertificateSource {
	return s.certificateSource
}

// CRLSource returns the merged CRL revocation source built from the signature and its
// timestamps. Additive accessor (integration-time, Phase 3); see CertificateSource's comment.
func (s *SignatureTimestampSource[AS, SA]) CRLSource() *spi.ListRevocationSource[revocation.CRL] {
	return s.crlSource
}

// OCSPSource returns the merged OCSP revocation source built from the signature and its
// timestamps. Additive accessor (integration-time, Phase 3); see CertificateSource's comment.
func (s *SignatureTimestampSource[AS, SA]) OCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
	return s.ocspSource
}

// GetAttributeOrder exposes getAttributeOrder(SA) to format-specific overrides. Additive
// accessor (integration-time, Phase 3); see CertificateSource's comment.
func (s *SignatureTimestampSource[AS, SA]) GetAttributeOrder(signatureAttribute SA) *int {
	return s.getAttributeOrder(signatureAttribute)
}

// createAndValidate creates and validates all timestamps. Must be called only once.
// Port of the protected createAndValidate().
func (s *SignatureTimestampSource[AS, SA]) createAndValidate() {
	s.populateTimestampTokens()
	s.validateTimestamps()
}

// AddExternalTimestamp allows adding an external timestamp. The given timestamp must be
// processed before. Port of addExternalTimestamp(TimestampToken).
func (s *SignatureTimestampSource[AS, SA]) AddExternalTimestamp(timestamp *validation.TimestampToken) {
	// if timestamp tokens not created yet
	if s.detachedTimestamps == nil {
		s.createAndValidate()
	}
	s.processExternalTimestamp(timestamp)
	s.detachedTimestamps = append(s.detachedTimestamps, timestamp)
}

// AddExternalEvidenceRecord allows adding an external evidence record covering the signature
// file. The given evidence record must be processed before.
// Port of addExternalEvidenceRecord(EvidenceRecord).
func (s *SignatureTimestampSource[AS, SA]) AddExternalEvidenceRecord(evidenceRecord validation.EvidenceRecord) {
	// if evidence records not created yet
	if s.detachedEvidenceRecords == nil {
		s.createAndValidate()
	}
	s.processExternalEvidenceRecord(evidenceRecord)
	s.detachedEvidenceRecords = append(s.detachedEvidenceRecords, evidenceRecord)
}

// populateTimestampTokens populates all the lists by data found into the signature.
// Port of the protected makeTimestampTokens() (0-arg overload); named distinctly from the
// 3-arg makeTimestampTokens(SA, TimestampType, List) below since Go has no overloading.
func (s *SignatureTimestampSource[AS, SA]) populateTimestampTokens() {
	// initialize timestamp lists
	s.contentTimestamps = []*validation.TimestampToken{}
	s.signatureTimestamps = []*validation.TimestampToken{}
	s.sigAndRefsTimestamps = []*validation.TimestampToken{}
	s.refsOnlyTimestamps = []*validation.TimestampToken{}
	s.archiveTimestamps = []*validation.TimestampToken{}
	s.detachedTimestamps = []*validation.TimestampToken{}

	s.embeddedEvidenceRecords = []validation.EvidenceRecord{}
	s.detachedEvidenceRecords = []validation.EvidenceRecord{}

	// initialize combined revocation sources
	s.crlSource = spi.NewListRevocationSourceFrom[revocation.CRL](s.signature.CRLSource())
	s.ocspSource = spi.NewListRevocationSourceFrom[revocation.OCSP](s.signature.OCSPSource())
	s.certificateSource = spi.NewListCertificateSourceFromOne(s.signature.CertificateSource())

	// a list of all embedded references
	s.unsignedPropertiesReferences = []*validation.TimestampedReference{}

	s.makeTimestampTokensFromSignedAttributes()
	s.makeTimestampTokensFromUnsignedAttributes()
}

// makeTimestampTokensFromSignedAttributes creates TimestampTokens from all instances extracted
// from signed attributes (content TSTs). Port of the protected
// makeTimestampTokensFromSignedAttributes().
func (s *SignatureTimestampSource[AS, SA]) makeTimestampTokensFromSignedAttributes() {
	extractedSignedSignatureProperties := s.getSignedSignatureProperties()
	if extractedSignedSignatureProperties == nil || !extractedSignedSignatureProperties.IsExist() {
		return
	}

	for _, signedAttribute := range extractedSignedSignatureProperties.Attributes() {

		var timestampTokens []*validation.TimestampToken

		if s.overrides.IsContentTimestamp(signedAttribute) {
			timestampTokens = s.makeTimestampTokens(signedAttribute, enumerations.TimestampType_CONTENT_TIMESTAMP, s.SignerDataReferences())

		} else if s.overrides.IsAllDataObjectsTimestamp(signedAttribute) {
			timestampTokens = s.makeTimestampTokens(signedAttribute, enumerations.TimestampType_ALL_DATA_OBJECTS_TIMESTAMP, s.SignerDataReferences())

		} else if s.overrides.IsIndividualDataObjectsTimestamp(signedAttribute) {
			timestampTokens = s.makeTimestampTokensDefault(signedAttribute, enumerations.TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP)

		} else {
			continue
		}

		if len(timestampTokens) == 0 {
			continue
		}

		s.populateSources(timestampTokens)
		s.contentTimestamps = append(s.contentTimestamps, timestampTokens...)
	}
}

// makeTimestampTokensFromUnsignedAttributes creates TimestampTokens from found instances in
// unsigned properties. Port of the protected makeTimestampTokensFromUnsignedAttributes().
func (s *SignatureTimestampSource[AS, SA]) makeTimestampTokensFromUnsignedAttributes() {
	extractedUnsignedSignatureProperties := s.getUnsignedSignatureProperties()
	if extractedUnsignedSignatureProperties == nil || !extractedUnsignedSignatureProperties.IsExist() {
		return
	}

	var allTimestamps []*validation.TimestampToken

	for _, unsignedAttribute := range extractedUnsignedSignatureProperties.Attributes() {
		var timestampTokens []*validation.TimestampToken

		switch {
		case s.overrides.IsSignatureTimestamp(unsignedAttribute):
			timestampTokens = s.makeTimestampTokens(unsignedAttribute, enumerations.TimestampType_SIGNATURE_TIMESTAMP, s.getSignatureTimestampReferences())
			if len(timestampTokens) == 0 {
				continue
			}
			s.signatureTimestamps = append(s.signatureTimestamps, timestampTokens...)

		case s.overrides.IsCompleteCertificateRef(unsignedAttribute) || s.overrides.IsAttributeCertificateRef(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getTimestampedCertificateRefs(unsignedAttribute))
			continue

		case s.overrides.IsCompleteRevocationRef(unsignedAttribute) || s.overrides.IsAttributeRevocationRef(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getTimestampedRevocationRefs(unsignedAttribute))
			continue

		case s.overrides.IsRefsOnlyTimestamp(unsignedAttribute):
			references := []*validation.TimestampedReference{}
			addReferences(&references, s.unsignedPropertiesReferences)

			timestampTokens = s.makeTimestampTokens(unsignedAttribute, enumerations.TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP, references)
			if len(timestampTokens) == 0 {
				continue
			}
			s.refsOnlyTimestamps = append(s.refsOnlyTimestamps, timestampTokens...)

		case s.overrides.IsSigAndRefsTimestamp(unsignedAttribute):
			references := []*validation.TimestampedReference{}

			processedSignatureTimestamps := filterSignatureTimestamps(allTimestamps)
			addReferences(&references, s.getEncapsulatedReferencesFromTimestamps(processedSignatureTimestamps))
			addReferences(&references, s.unsignedPropertiesReferences)

			timestampTokens = s.makeTimestampTokens(unsignedAttribute, enumerations.TimestampType_VALIDATION_DATA_TIMESTAMP, references)
			if len(timestampTokens) == 0 {
				continue
			}
			s.sigAndRefsTimestamps = append(s.sigAndRefsTimestamps, timestampTokens...)

		case s.overrides.IsCertificateValues(unsignedAttribute) || s.overrides.IsAttrAuthoritiesCertValues(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getTimestampedCertificateValues(unsignedAttribute))
			continue

		case s.overrides.IsRevocationValues(unsignedAttribute) || s.overrides.IsAttributeRevocationValues(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getTimestampedRevocationValues(unsignedAttribute))
			continue

		case s.overrides.IsArchiveTimestamp(unsignedAttribute):
			timestampTokens = s.makeTimestampTokensDefault(unsignedAttribute, enumerations.TimestampType_ARCHIVE_TIMESTAMP)
			if len(timestampTokens) == 0 {
				continue
			}
			s.setArchiveTimestampType(timestampTokens, unsignedAttribute)
			s.incorporateArchiveTimestampReferencesForTokens(timestampTokens, allTimestamps)

			s.archiveTimestamps = append(s.archiveTimestamps, timestampTokens...)

		case s.overrides.IsTimeStampValidationData(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getTimestampValidationData(unsignedAttribute))
			continue

		case s.overrides.IsAnyValidationData(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getAnyValidationData(unsignedAttribute))
			continue

		case s.overrides.IsValidationDataReferences(unsignedAttribute):
			addReferences(&s.unsignedPropertiesReferences, s.getValidationDataReferences(unsignedAttribute))
			continue

		case s.overrides.IsCounterSignature(unsignedAttribute):
			counterSignatures := s.overrides.GetCounterSignatures(unsignedAttribute)
			addReferences(&s.unsignedPropertiesReferences, s.getCounterSignaturesReferences(counterSignatures))
			continue

		case s.overrides.IsSignaturePolicyStore(unsignedAttribute):
			// not processed
			continue

		case s.overrides.IsEvidenceRecord(unsignedAttribute):
			evidenceRecords := s.overrides.MakeEvidenceRecords(unsignedAttribute, s.unsignedPropertiesReferences)
			if len(evidenceRecords) == 0 {
				continue
			}
			s.incorporateEvidenceRecordEvidenceReferences(evidenceRecords, allTimestamps)

			for _, evidenceRecord := range evidenceRecords {
				s.populateSourcesFromEvidenceRecord(evidenceRecord)
			}
			s.embeddedEvidenceRecords = append(s.embeddedEvidenceRecords, evidenceRecords...)
			continue

		default:
			// Unsupported attribute encountered during TimestampSource processing; entry is
			// skipped. Upstream logs a warning here (slf4j dropped per PORTING_PLAN.md, not
			// load-bearing).
			continue
		}

		s.populateSources(timestampTokens)
		allTimestamps = append(allTimestamps, timestampTokens...)
	}
}

// populateSourcesFromEvidenceRecord adds the timestamps of evidenceRecord to the merged
// sources. Port of the `populateSources(evidenceRecord.getTimestamps())` call inline in
// makeTimestampTokensFromUnsignedAttributes's evidence-record branch.
func (s *SignatureTimestampSource[AS, SA]) populateSourcesFromEvidenceRecord(evidenceRecord validation.EvidenceRecord) {
	s.populateSources(evidenceRecord.Timestamps())
}

// getSignedSignatureProperties returns the 'signed-signature-properties' element of the
// signature. Port of the protected getSignedSignatureProperties().
func (s *SignatureTimestampSource[AS, SA]) getSignedSignatureProperties() validation.SignatureProperties[SA] {
	if s.signedSignatureProperties == nil {
		s.signedSignatureProperties = s.overrides.BuildSignedSignatureProperties()
	}
	return s.signedSignatureProperties
}

// getUnsignedSignatureProperties returns the 'unsigned-signature-properties' element of the
// signature. Port of the protected getUnsignedSignatureProperties().
func (s *SignatureTimestampSource[AS, SA]) getUnsignedSignatureProperties() validation.SignatureProperties[SA] {
	if s.unsignedSignatureProperties == nil {
		s.unsignedSignatureProperties = s.overrides.BuildUnsignedSignatureProperties()
	}
	return s.unsignedSignatureProperties
}

// makeTimestampTokensDefault creates timestamp tokens from signatureAttribute with an empty
// initial list of references. Port of makeTimestampTokens(SA, TimestampType) (the two-arg
// overload); named distinctly from the three-arg overload below since Go has no overloading.
func (s *SignatureTimestampSource[AS, SA]) makeTimestampTokensDefault(signatureAttribute SA,
	timestampType enumerations.TimestampType) []*validation.TimestampToken {
	return s.makeTimestampTokens(signatureAttribute, timestampType, []*validation.TimestampedReference{})
}

// makeTimestampTokens creates timestamp tokens from signatureAttribute with the given list of
// TimestampedReferences. Port of makeTimestampTokens(SA, TimestampType, List). Dispatches through
// s.overrides.MakeTimestampTokens (plural) rather than wrapping MakeTimestampToken (singular)
// directly, so a format that shadows the plural method (XAdES) is actually reached - see the
// interface doc comment on MakeTimestampTokens.
func (s *SignatureTimestampSource[AS, SA]) makeTimestampTokens(signatureAttribute SA, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) []*validation.TimestampToken {
	return s.overrides.MakeTimestampTokens(signatureAttribute, timestampType, references)
}

// MakeTimestampTokens is the default SignatureTimestampSourceOverrides.MakeTimestampTokens
// implementation, promoted by embedding to every concrete format that does not shadow it with its
// own: it wraps MakeTimestampToken's (singular) result in a single-element slice, i.e. exactly the
// behaviour makeTimestampTokens (above) had before this plural override hook existed. Port of the
// base makeTimestampTokens(SA, TimestampType, List)'s own body (Java: `TimestampToken
// timestampToken = makeTimestampToken(...); return timestampToken == null ? emptyList() :
// singletonList(timestampToken);`).
func (s *SignatureTimestampSource[AS, SA]) MakeTimestampTokens(signatureAttribute SA, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) []*validation.TimestampToken {
	timestampToken := s.overrides.MakeTimestampToken(signatureAttribute, timestampType, references)
	if timestampToken != nil {
		return []*validation.TimestampToken{timestampToken}
	}
	return nil
}

// SignerDataReferences implements validation.TimestampSource.
// Port of getSignerDataReferences().
func (s *SignatureTimestampSource[AS, SA]) SignerDataReferences() []*validation.TimestampedReference {
	return SignerDataTimestampedReferences(s.signature.SignatureScopes())
}

// getSignatureTimestampReferences returns a list of TimestampedReferences for a
// "signature-timestamp" element. Port of the protected getSignatureTimestampReferences().
func (s *SignatureTimestampSource[AS, SA]) getSignatureTimestampReferences() []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	addReferences(&references, s.getEncapsulatedReferencesFromTimestamps(s.ContentTimestamps()))
	addReferences(&references, s.SignerDataReferences())
	addReference(&references, s.getSignatureReference())
	addReferences(&references, s.getSigningCertificateTimestampReferences())
	return references
}

// getSignatureReference creates a timestamped reference for the current signature.
// Port of the protected getSignatureReference().
func (s *SignatureTimestampSource[AS, SA]) getSignatureReference() *validation.TimestampedReference {
	return validation.NewTimestampedReference(s.signature.ID(), enumerations.TimestampedObjectType_SIGNATURE)
}

// getEncapsulatedReferencesFromTimestamps returns a list of TimestampedReferences for tokens
// encapsulated within timestampTokens. Port of the protected
// getEncapsulatedReferencesFromTimestamps(List).
func (s *SignatureTimestampSource[AS, SA]) getEncapsulatedReferencesFromTimestamps(timestampTokens []*validation.TimestampToken) []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	for _, timestampToken := range timestampTokens {
		addReferences(&references, must(ReferencesFromTimestamp(timestampToken, s.certificateSource, s.crlSource, s.ocspSource)))
	}
	return references
}

// getEncapsulatedReferencesFromEvidenceRecords returns a list of TimestampedReferences for
// tokens encapsulated within evidenceRecords. Port of the protected
// getEncapsulatedReferencesFromEvidenceRecords(List).
func (s *SignatureTimestampSource[AS, SA]) getEncapsulatedReferencesFromEvidenceRecords(evidenceRecords []validation.EvidenceRecord) []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	for _, er := range evidenceRecords {
		addReferences(&references, must(ReferencesFromEvidenceRecord(er, s.certificateSource, s.crlSource, s.ocspSource)))
	}
	return references
}

// getSigningCertificateTimestampReferences returns a list of TimestampedReferences created from
// signing certificates of the signature. Port of the protected getSigningCertificateTimestampReferences().
func (s *SignatureTimestampSource[AS, SA]) getSigningCertificateTimestampReferences() []*validation.TimestampedReference {
	signatureCertificateSource := s.signature.CertificateSource()
	return CreateReferencesForCertificateRefs(signatureCertificateSource.SigningCertificateRefs(), signatureCertificateSource, s.certificateSource)
}

// GetKeyInfoReferences returns references from the KeyInfo (for XAdES) encapsulated elements.
// Port of the protected getKeyInfoReferences().
func (s *SignatureTimestampSource[AS, SA]) GetKeyInfoReferences() []*validation.TimestampedReference {
	signatureCertificateSource := s.signature.CertificateSource()
	return CreateReferencesForCertificates(signatureCertificateSource.KeyInfoCertificates())
}

// getTimestampedCertificateRefs returns a list of TimestampedReference certificate refs found in
// unsignedAttribute. Port of the protected getTimestampedCertificateRefs(SA).
func (s *SignatureTimestampSource[AS, SA]) getTimestampedCertificateRefs(unsignedAttribute SA) []*validation.TimestampedReference {
	return CreateReferencesForCertificateRefs(s.overrides.GetCertificateRefs(unsignedAttribute), s.signature.CertificateSource(), s.certificateSource)
}

// getTimestampedRevocationRefs returns a list of TimestampedReference revocation refs found in
// unsignedAttribute. Port of the protected getTimestampedRevocationRefs(SA).
func (s *SignatureTimestampSource[AS, SA]) getTimestampedRevocationRefs(unsignedAttribute SA) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	timestampedReferences = append(timestampedReferences,
		CreateReferencesForCRLRefs(s.overrides.GetCRLRefs(unsignedAttribute), s.signature.CRLSource(), s.crlSource)...)
	timestampedReferences = append(timestampedReferences,
		must(CreateReferencesForOCSPRefs(s.overrides.GetOCSPRefs(unsignedAttribute), s.signature.OCSPSource(), s.certificateSource, s.ocspSource))...)
	return timestampedReferences
}

// getTimestampedCertificateValues returns a list of TimestampedReferences from unsignedAttribute
// containing certificate values. Port of the protected getTimestampedCertificateValues(SA).
func (s *SignatureTimestampSource[AS, SA]) getTimestampedCertificateValues(unsignedAttribute SA) []*validation.TimestampedReference {
	return CreateReferencesForIdentifiers(s.overrides.GetEncapsulatedCertificateIdentifiers(unsignedAttribute), enumerations.TimestampedObjectType_CERTIFICATE)
}

// getTimestampedRevocationValues returns a list of timestamped revocation references extracted
// from unsignedAttribute. Port of the protected getTimestampedRevocationValues(SA).
func (s *SignatureTimestampSource[AS, SA]) getTimestampedRevocationValues(unsignedAttribute SA) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	timestampedReferences = append(timestampedReferences,
		CreateReferencesForCRLBinaries(s.overrides.GetEncapsulatedCRLIdentifiers(unsignedAttribute))...)
	timestampedReferences = append(timestampedReferences,
		must(CreateReferencesForOCSPBinaries(s.overrides.GetEncapsulatedOCSPIdentifiers(unsignedAttribute), s.certificateSource))...)
	return timestampedReferences
}

// incorporateArchiveTimestampReferencesForTokens ports the private
// incorporateArchiveTimestampReferences(List, List).
func (s *SignatureTimestampSource[AS, SA]) incorporateArchiveTimestampReferencesForTokens(createdTimestampTokens, previousTimestamps []*validation.TimestampToken) {
	for _, timestampToken := range createdTimestampTokens {
		s.overrides.IncorporateArchiveTimestampReferences(timestampToken, previousTimestamps)
	}
}

// IncorporateArchiveTimestampReferences incorporates all the timestamped references for the
// given archive timestampToken. Port of the protected incorporateArchiveTimestampReferences
// (TimestampToken, List).
func (s *SignatureTimestampSource[AS, SA]) IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken, previousTimestamps []*validation.TimestampToken) {
	timestampAddReferences(timestampToken, s.getArchiveTimestampReferences(previousTimestamps))
}

// incorporateEvidenceRecordEvidenceReferences ports the private
// incorporateEvidenceRecordEvidenceReferences(List, List).
func (s *SignatureTimestampSource[AS, SA]) incorporateEvidenceRecordEvidenceReferences(createdEvidenceRecords []validation.EvidenceRecord, previousTimestamps []*validation.TimestampToken) {
	for _, evidenceRecord := range createdEvidenceRecords {
		evidenceRecord.SetTimestampedReferences(mergeReferences(evidenceRecord.TimestampedReferences(), s.getArchiveTimestampReferences(previousTimestamps)))
		evidenceRecord.SetTimestampedReferences(mergeReferences(evidenceRecord.TimestampedReferences(), s.getEncapsulatedReferencesFromEvidenceRecords(s.embeddedEvidenceRecords)))
		ProcessEvidenceRecordTimestamps(evidenceRecord)
	}
}

// getArchiveTimestampReferences returns a list of time-stamped references for an archival
// time-stamp. Port of the protected getArchiveTimestampReferences(List).
func (s *SignatureTimestampSource[AS, SA]) getArchiveTimestampReferences(previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	addReferences(&timestampedReferences, s.getSignatureTimestampReferences())
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(previousTimestamps))
	addReferences(&timestampedReferences, s.unsignedPropertiesReferences)
	return timestampedReferences
}

// GetSignatureSignedDataReferences returns a list of all TimestampedReferences found into CMS
// SignedData of the signature. NOTE: used only in ASiC-E CAdES. Port of the protected
// getSignatureSignedDataReferences(), empty by default.
//
// Java subclasses may override this concrete (non-abstract) method; wired into
// SignatureTimestampSourceOverrides (integration-time fix, Phase 3) so that the base's own
// internal callers (processExternalTimestamp, processExternalEvidenceRecord) reach a
// format-specific override through s.overrides rather than statically binding to this body -
// the same virtual-dispatch gap every other override above is routed around. This is the
// default (empty) body for an override with nothing format-specific to add.
func (s *SignatureTimestampSource[AS, SA]) GetSignatureSignedDataReferences() []*validation.TimestampedReference {
	return nil
}

// getTimestampValidationData returns a list of TimestampedReferences encapsulated to the
// "timestamp-validation-data" unsignedAttribute. Port of the protected
// getTimestampValidationData(SA).
func (s *SignatureTimestampSource[AS, SA]) getTimestampValidationData(unsignedAttribute SA) []*validation.TimestampedReference {
	return s.getAnyValidationData(unsignedAttribute)
}

// getAnyValidationData returns a list of TimestampedReferences encapsulated to the
// "any-validation-data" unsignedAttribute. Port of the protected getAnyValidationData(SA).
func (s *SignatureTimestampSource[AS, SA]) getAnyValidationData(unsignedAttribute SA) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	addReferences(&timestampedReferences, CreateReferencesForIdentifiers(
		s.overrides.GetEncapsulatedCertificateIdentifiers(unsignedAttribute), enumerations.TimestampedObjectType_CERTIFICATE))
	addReferences(&timestampedReferences, CreateReferencesForCRLBinaries(s.overrides.GetEncapsulatedCRLIdentifiers(unsignedAttribute)))
	addReferences(&timestampedReferences, must(CreateReferencesForOCSPBinaries(s.overrides.GetEncapsulatedOCSPIdentifiers(unsignedAttribute), s.certificateSource)))
	return timestampedReferences
}

// getValidationDataReferences returns a list of TimestampedReferences encapsulated to the
// "validation-data-references" unsignedAttribute. Port of the protected
// getValidationDataReferences(SA).
func (s *SignatureTimestampSource[AS, SA]) getValidationDataReferences(unsignedAttribute SA) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	addReferences(&timestampedReferences, CreateReferencesForCertificateRefs(
		s.overrides.GetCertificateRefs(unsignedAttribute), s.signature.CertificateSource(), s.certificateSource))
	addReferences(&timestampedReferences, CreateReferencesForCRLRefs(
		s.overrides.GetCRLRefs(unsignedAttribute), s.signature.CRLSource(), s.crlSource))
	addReferences(&timestampedReferences, must(CreateReferencesForOCSPRefs(
		s.overrides.GetOCSPRefs(unsignedAttribute), s.signature.OCSPSource(), s.certificateSource, s.ocspSource)))
	return timestampedReferences
}

// getCounterSignaturesReferences returns a list of TimestampedReferences encapsulated from the
// list of counterSignatures. Port of the protected getCounterSignaturesReferences(List).
func (s *SignatureTimestampSource[AS, SA]) getCounterSignaturesReferences(counterSignatures []validation.AdvancedSignature) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	if len(counterSignatures) > 0 {
		for _, counterSignature := range counterSignatures {
			references = append(references, s.overrides.GetCounterSignatureReferences(counterSignature)...)
		}
	}
	return references
}

// GetCounterSignatureReferences returns a list of references extracted from counterSignature.
// Port of the protected getCounterSignatureReferences(AdvancedSignature).
func (s *SignatureTimestampSource[AS, SA]) GetCounterSignatureReferences(counterSignature validation.AdvancedSignature) []*validation.TimestampedReference {
	var counterSigReferences []*validation.TimestampedReference

	counterSigReferences = append(counterSigReferences, validation.NewTimestampedReference(counterSignature.ID(), enumerations.TimestampedObjectType_SIGNATURE))

	signatureCertificateSource := counterSignature.CertificateSource()
	addReferences(&counterSigReferences, CreateReferencesForCertificates(signatureCertificateSource.Certificates()))

	counterSignatureTimestampSource := counterSignature.TimestampSource()
	addReferences(&counterSigReferences, counterSignatureTimestampSource.SignerDataReferences())
	addReferences(&counterSigReferences, counterSignatureTimestampSource.UnsignedPropertiesReferences())
	addReferences(&counterSigReferences, s.getEncapsulatedReferencesFromTimestamps(counterSignatureTimestampSource.AllTimestamps()))

	return counterSigReferences
}

// filterSignatureTimestamps ports the private filterSignatureTimestamps(List).
func filterSignatureTimestamps(previousTimestampedTimestamp []*validation.TimestampToken) []*validation.TimestampToken {
	result := []*validation.TimestampToken{}
	for _, timestampToken := range previousTimestampedTimestamp {
		if timestampToken.TimeStampType() == enumerations.TimestampType_SIGNATURE_TIMESTAMP {
			result = append(result, timestampToken)
		}
	}
	return result
}

// setArchiveTimestampType ports the private setArchiveTimestampType(List, SA).
func (s *SignatureTimestampSource[AS, SA]) setArchiveTimestampType(timestampTokens []*validation.TimestampToken, unsignedAttribute SA) {
	archiveTimestampType := s.overrides.GetArchiveTimestampType(unsignedAttribute)
	for _, timestampToken := range timestampTokens {
		timestampToken.SetArchiveTimestampType(archiveTimestampType)
	}
}

// validateTimestamps validates the list of all timestamps present in the source.
// Port of the protected validateTimestamps().
func (s *SignatureTimestampSource[AS, SA]) validateTimestamps() {

	// This validates the content-timestamp tokensToProcess present in the signature.
	for _, timestampToken := range s.ContentTimestamps() {
		messageDigest := s.getTimestampMessageImprintDigestBuilder(timestampToken).ContentTimestampMessageDigest()
		timestampToken.MatchDataMessageDigest(messageDigest)
		timestampToken.SetTimestampScopes(s.getTimestampScopes(timestampToken))
	}

	// This validates the signature timestamp tokensToProcess present in the signature.
	for _, timestampToken := range s.SignatureTimestamps() {
		messageDigest := s.getTimestampMessageImprintDigestBuilder(timestampToken).SignatureTimestampMessageDigest()
		timestampToken.MatchDataMessageDigest(messageDigest)
	}

	// This validates the SigAndRefs timestamp tokensToProcess present in the signature.
	for _, timestampToken := range s.TimestampsX1() {
		messageDigest := s.getTimestampMessageImprintDigestBuilder(timestampToken).TimestampX1MessageDigest()
		timestampToken.MatchDataMessageDigest(messageDigest)
	}

	// This validates the RefsOnly timestamp tokensToProcess present in the signature.
	for _, timestampToken := range s.TimestampsX2() {
		messageDigest := s.getTimestampMessageImprintDigestBuilder(timestampToken).TimestampX2MessageDigest()
		timestampToken.MatchDataMessageDigest(messageDigest)
	}

	// This validates the archive timestamp tokensToProcess present in the signature.
	for _, timestampToken := range s.ArchiveTimestamps() {
		if !timestampToken.IsProcessed() {
			messageDigest := s.getTimestampMessageImprintDigestBuilder(timestampToken).ArchiveTimestampMessageDigest()
			timestampToken.MatchDataMessageDigest(messageDigest)
			timestampToken.SetTimestampScopes(s.getTimestampScopes(timestampToken))
		}
	}
}

// GetTimestampMessageImprintDigestBuilder returns a TimestampMessageDigestBuilder to compute
// message digest with the provided DigestAlgorithm. Port of the protected abstract
// getTimestampMessageImprintDigestBuilder(DigestAlgorithm).
func (s *SignatureTimestampSource[AS, SA]) GetTimestampMessageImprintDigestBuilder(digestAlgorithm enumerations.DigestAlgorithm) TimestampMessageDigestBuilder {
	return s.overrides.GetTimestampMessageImprintDigestBuilderForAlgorithm(digestAlgorithm)
}

// getTimestampMessageImprintDigestBuilder returns a related TimestampMessageDigestBuilder.
// Port of the protected abstract getTimestampMessageImprintDigestBuilder(TimestampToken).
func (s *SignatureTimestampSource[AS, SA]) getTimestampMessageImprintDigestBuilder(timestampToken *validation.TimestampToken) TimestampMessageDigestBuilder {
	return s.overrides.GetTimestampMessageImprintDigestBuilderForToken(timestampToken)
}

// getTimestampScopes generates timestamp token scopes. Port of the protected
// getTimestampScopes(TimestampToken).
func (s *SignatureTimestampSource[AS, SA]) getTimestampScopes(timestampToken *validation.TimestampToken) []scope.SignatureScope {
	timestampScopeFinder := validationscope.NewEncapsulatedTimestampScopeFinder()
	timestampScopeFinder.SetSignature(s.signature)
	return timestampScopeFinder.FindTimestampScope(timestampToken)
}

// processExternalTimestamp ports the private processExternalTimestamp(TimestampToken).
func (s *SignatureTimestampSource[AS, SA]) processExternalTimestamp(externalTimestamp *validation.TimestampToken) {
	// add all validation data present in Signature CMS SignedData, because an external
	// timestamp covers a whole signature file
	timestampAddReferences(externalTimestamp, s.overrides.GetSignatureSignedDataReferences())
	// add references from previously added timestamps
	timestampAddReferences(externalTimestamp, s.getEncapsulatedReferencesFromTimestamps(
		s.getTimestampsCoveredByManifest(externalTimestamp.ManifestFile())))
	// add existing counter signatures
	timestampAddReferences(externalTimestamp, s.overrides.GetCounterSignatureReferences(s.signature))
	// populate timestamp certificate source with values present in the timestamp
	s.populateSource(externalTimestamp)
}

// getTimestampsCoveredByManifest ports the private getTimestampsCoveredByManifest(ManifestFile).
func (s *SignatureTimestampSource[AS, SA]) getTimestampsCoveredByManifest(manifestFile *model.ManifestFile) []*validation.TimestampToken {
	result := []*validation.TimestampToken{}
	for _, timestampToken := range s.overrides.AllTimestamps() {
		if timestampTokenSliceContains(s.detachedTimestamps, timestampToken) &&
			(manifestFile == nil || !manifestFile.IsDocumentCovered(timestampToken.Filename())) {
			// the detached timestamp is not covered, continue
			continue
		}
		result = append(result, timestampToken)
	}
	return result
}

// timestampTokenSliceContains reports whether timestampTokens contains candidate, by pointer
// identity - the Go counterpart of Java's List#contains(Object) here, which relies on
// TimestampToken's default (identity) equals(), since TimestampToken does not override it.
func timestampTokenSliceContains(timestampTokens []*validation.TimestampToken, candidate *validation.TimestampToken) bool {
	for _, timestampToken := range timestampTokens {
		if timestampToken == candidate {
			return true
		}
	}
	return false
}

// processExternalEvidenceRecord ports the private processExternalEvidenceRecord(EvidenceRecord).
func (s *SignatureTimestampSource[AS, SA]) processExternalEvidenceRecord(evidenceRecord validation.EvidenceRecord) {
	timestampedReferences := []*validation.TimestampedReference{}
	addReferences(&timestampedReferences, s.getSignatureTimestampReferences())
	addReferences(&timestampedReferences, s.overrides.GetSignatureSignedDataReferences())
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(s.SignatureTimestamps()))
	addReferences(&timestampedReferences, s.unsignedPropertiesReferences)
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(s.TimestampsX1()))
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(s.TimestampsX2()))
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(s.ArchiveTimestamps()))
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromTimestamps(
		s.getTimestampsCoveredByManifest(evidenceRecord.ManifestFile())))
	addReferences(&timestampedReferences, s.getEncapsulatedReferencesFromEvidenceRecords(s.EmbeddedEvidenceRecords()))

	evidenceRecord.SetTimestampedReferences(mergeReferences(evidenceRecord.TimestampedReferences(), timestampedReferences))

	ProcessEvidenceRecordTimestamps(evidenceRecord)
	ProcessEmbeddedEvidenceRecords(evidenceRecord)
	s.populateSourceFromEvidenceRecord(evidenceRecord)
}

// populateSources allows populating all merged sources with data extracted from timestampTokens.
// Port of the protected populateSources(List).
func (s *SignatureTimestampSource[AS, SA]) populateSources(timestampTokens []*validation.TimestampToken) {
	for _, timestampToken := range timestampTokens {
		s.populateSource(timestampToken)
	}
}

// populateSource allows populating all merged sources with data extracted from timestampToken.
// Port of the protected populateSources(TimestampToken); named distinctly to avoid colliding
// with the plural overload above (Go has no overloading).
func (s *SignatureTimestampSource[AS, SA]) populateSource(timestampToken *validation.TimestampToken) {
	if timestampToken != nil {
		s.certificateSource.Add(timestampToken.CertificateSource())
		s.crlSource.Add(timestampToken.CRLSource())
		s.ocspSource.Add(timestampToken.OCSPSource())
	}
}

// populateSourceFromEvidenceRecord allows populating all sources from an external evidence
// record. Port of the protected populateSources(EvidenceRecord).
func (s *SignatureTimestampSource[AS, SA]) populateSourceFromEvidenceRecord(externalEvidenceRecord validation.EvidenceRecord) {
	if externalEvidenceRecord != nil {
		s.certificateSource.Add(externalEvidenceRecord.CertificateSource())
		s.crlSource.Add(externalEvidenceRecord.CRLSource())
		s.ocspSource.Add(externalEvidenceRecord.OCSPSource())

		s.populateSources(externalEvidenceRecord.Timestamps())
	}
}

// IsTimestamped implements validation.TimestampSource. Port of isTimestamped(String,
// TimestampedObjectType).
func (s *SignatureTimestampSource[AS, SA]) IsTimestamped(tokenId string, objectType enumerations.TimestampedObjectType) bool {
	return s.isTimestamped(s.signature, tokenId, objectType)
}

// isTimestamped ports the private isTimestamped(AdvancedSignature, String, TimestampedObjectType).
func (s *SignatureTimestampSource[AS, SA]) isTimestamped(signature validation.AdvancedSignature, tokenId string,
	objectType enumerations.TimestampedObjectType) bool {
	target := validation.NewTimestampedReference(tokenId, objectType)
	for _, timestampToken := range s.overrides.AllTimestamps() {
		for _, reference := range timestampToken.TimestampedReferences() {
			if reference.Equals(target) {
				return true
			}
		}
	}
	masterSignature := signature.MasterSignature()
	if masterSignature != nil {
		return s.isTimestamped(masterSignature, tokenId, objectType)
	}
	return false
}

// getAttributeOrder gets the position of signatureAttribute either within signed or unsigned
// properties. Port of the protected getAttributeOrder(SA); returns nil, matching Java's null
// Integer, when not found.
//
// Java writes `signatureAttribute.equals(property)`, which dispatches to the CONCRETE
// attribute class's equals() override - CAdESAttribute, XAdESAttribute, JAdESAttribute and
// CBAdESAttribute all override it as `Objects.equals(getIdentifier(), that.getIdentifier())`,
// i.e. value equality on the SA-... identifier - not to Object's reference identity. Comparing
// the Go pointers instead (the file header's original assumption) can never match: every
// SignatureProperties.Attributes() implementation REBUILDS a fresh attribute list on each call
// (see cades/cades_sig_properties.go, xades/xades_sig_properties.go), so the SA handed in here
// - produced by an earlier Attributes() call inside populateTimestampTokens - is never the same
// pointer as any element of the list re-read below. The result was a silent, always-nil order,
// dropping the "-OOA-<n>" component from every embedded timestamp's position string and so from
// every encapsulated TimestampToken's T-... identifier (found by the phase-8f document-level
// harness: 30 of its 60 CAdES/XAdES/PAdES/ASiC fixtures).
func (s *SignatureTimestampSource[AS, SA]) getAttributeOrder(signatureAttribute SA) *int {
	target := signatureAttribute.Identifier()
	signedAttributes := s.getSignedSignatureProperties().Attributes()
	for i := 0; i < len(signedAttributes); i++ {
		if signatureAttributeEquals(target, signedAttributes[i]) {
			position := i
			return &position
		}
	}
	unsignedAttributes := s.getUnsignedSignatureProperties().Attributes()
	for i := 0; i < len(unsignedAttributes); i++ {
		if signatureAttributeEquals(target, unsignedAttributes[i]) {
			position := i
			return &position
		}
	}
	return nil
}

// signatureAttributeEquals reports whether candidate carries the identifier target, reproducing
// the concrete attribute classes' equals(Object) override. model.IdentifierBase.Equals already
// carries Java Identifier#equals's getClass() check (the ported Java simple class name), so two
// attributes of different formats never compare equal even on an identical digest.
func signatureAttributeEquals[SA validation.SignatureAttribute](target identifier.SignatureAttributeIdentifier, candidate SA) bool {
	candidateIdentifier := candidate.Identifier()
	return target.Equals(&candidateIdentifier)
}

// compile-time assertion: a SignatureTimestampSource embeds AbstractTimestampSource and
// satisfies validation.TimestampSource, matching Java's "extends AbstractTimestampSource
// implements TimestampSource".
var _ validation.TimestampSource = (*SignatureTimestampSource[validation.AdvancedSignature, comparableSignatureAttributeStub])(nil)

// comparableSignatureAttributeStub is used only by the compile-time assertion above, to
// instantiate SignatureTimestampSource's SA type parameter with a concrete, comparable
// validation.SignatureAttribute implementation without depending on one from a later phase.
type comparableSignatureAttributeStub struct{}

// Identifier satisfies validation.SignatureAttribute for comparableSignatureAttributeStub.
func (comparableSignatureAttributeStub) Identifier() identifier.SignatureAttributeIdentifier {
	return identifier.SignatureAttributeIdentifier{}
}

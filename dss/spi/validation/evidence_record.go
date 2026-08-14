// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/evidencerecord/EvidenceRecord.java (DSS 6.5.RC1).
//
// SCC flattening: Java spi.x509.evidencerecord.EvidenceRecord lands in this same Go package
// per S2B_BRIEF.md's package layout table ("spi.x509.evidencerecord" is one of the packages
// flattened into dss/spi/validation).
//
// This interface's shape is additionally constrained by two already-landed sibling files that
// forward-reference it opaquely (dss/spi/validation/timestamp_source.go) and by name
// (dss/spi/validation/timestamp/abstract_timestamp_source.go's "FORWARD DEPENDENCY:
// EvidenceRecord" header comment, which documents the exact subset of methods it calls:
// Id(), TimestampedReferences(), SetTimestampedReferences(...), Timestamps(),
// CertificateSource(), CRLSource(), OCSPSource(), DetachedEvidenceRecords(), ManifestFile()).
// The full interface below matches that subset verbatim while adding every other Java-declared
// method, translated 1:1 with get/is dropped.
//
// FORWARD DEPENDENCIES (flagged per S2B_BRIEF.md): SignatureAttribute (Java
// spi.validation.SignatureAttribute) and EmbeddedEvidenceRecordHelper (Java
// spi.validation.evidencerecord.EmbeddedEvidenceRecordHelper) are both SCC-flattened into this
// same package by other chunks of phase 2b and are referenced here unqualified, opaquely (this
// file never calls a method on either - it only carries them through the interface, following
// the precedent of the already-landed AdvancedSignature/TimestampToken opaque forward
// references in timestamp_source.go).
package validation

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi"
)

// EvidenceRecord is the representation of an Evidence Record.
type EvidenceRecord interface {
	model.IdentifierBasedObject

	// Filename returns a name of the evidence record document, when present. Port of getFilename().
	Filename() string

	// ReferenceValidation returns a list of archive data object validations. Port of
	// getReferenceValidation().
	ReferenceValidation() []*model.ReferenceValidation

	// DetachedContents returns detached contents: in the case of the detached signature this is
	// the list of signed contents. Port of getDetachedContents().
	DetachedContents() []model.DSSDocument

	// CertificateSource gets a certificate source which contains ALL certificates embedded in
	// the evidence record. Port of getCertificateSource().
	CertificateSource() *spi.TokenCertificateSource

	// CRLSource gets a CRL source which contains ALL CRLs embedded in the evidence record. Port
	// of getCRLSource().
	//
	// Typed as the concrete *spi.OfflineCRLSourceBase, not the spi.OfflineRevocationSource[R]
	// interface Java's getCRLSource() returns and not the bare *spi.OfflineRevocationSourceBase[R]
	// this file originally used: OfflineRevocationSourceBase[R] alone dispatches RevocationToken
	// (singular) and RevocationTokens (plural) through an Init-registered overrides interface it
	// does not itself implement (see offline_revocation_source.go), so it satisfies neither
	// spi.OfflineRevocationSource[R] (missing RevocationTokens) nor
	// spi.MultipleRevocationSource[R] on its own. OfflineCRLSourceBase embeds it AND supplies a
	// concrete RevocationTokens, so it alone satisfies the full spi.OfflineRevocationSource[R]
	// contract plus AllRevocationReferences() (which abstract_timestamp_source.go's forward
	// dependency on this method needs) via promotion. TimestampToken.CRLSource() (already landed,
	// timestamp_crl_source.go) hit the identical problem and resolved it by returning a type built
	// the same way (*TimestampCRLSource, embedding *spi.CMSCRLSource, itself embedding
	// spi.OfflineCRLSourceBase).
	CRLSource() *spi.OfflineCRLSourceBase

	// OCSPSource gets an OCSP source which contains ALL OCSP responses embedded in the evidence
	// record. Port of getOCSPSource(); see CRLSource for the concrete-type judgment call.
	OCSPSource() *spi.OfflineOCSPSourceBase

	// Timestamps returns a list of incorporated timestamp tokens. Port of getTimestamps().
	Timestamps() []*TimestampToken

	// DetachedEvidenceRecords returns a list of detached evidence records covering the current
	// evidence record. Port of getDetachedEvidenceRecords().
	DetachedEvidenceRecords() []EvidenceRecord

	// AddExternalEvidenceRecord allows adding an external evidence record covering the current
	// evidence record. Port of addExternalEvidenceRecord(EvidenceRecord).
	AddExternalEvidenceRecord(evidenceRecord EvidenceRecord)

	// EvidenceRecordScopes returns a list of covered archival data objects. Port of
	// getEvidenceRecordScopes().
	EvidenceRecordScopes() []scope.SignatureScope

	// SetEvidenceRecordScopes sets a list of covered archival data objects. Port of
	// setEvidenceRecordScopes(List).
	SetEvidenceRecordScopes(evidenceRecordScopes []scope.SignatureScope)

	// OriginalDigestAlgorithm gets the DigestAlgorithm used on the first data object group's
	// digest computation. Port of getOriginalDigestAlgorithm().
	OriginalDigestAlgorithm() enumerations.DigestAlgorithm

	// StructureValidationResult returns a message if the structure validation fails: a list of
	// error messages if validation fails, an empty list if structural validation succeeds. Port
	// of getStructureValidationResult().
	StructureValidationResult() []string

	// EvidenceRecordType returns the type of the evidence record. Port of getEvidenceRecordType().
	EvidenceRecordType() enumerations.EvidenceRecordTypeEnum

	// Origin returns the origin of the evidence record. Port of getOrigin().
	Origin() enumerations.EvidenceRecordOrigin

	// ManifestFile returns a manifest file associated with the evidence record (used in ASiC).
	// Port of getManifestFile().
	ManifestFile() *model.ManifestFile

	// TimestampedReferences returns a list of references covered by the evidence record. Port
	// of getTimestampedReferences().
	TimestampedReferences() []*TimestampedReference

	// SetTimestampedReferences sets references to objects covered by the evidence record. Port
	// of setTimestampedReferences(List).
	//
	// Java mutates the List<TimestampedReference> getTimestampedReferences() returns in place
	// (List reference semantics); Go slices returned by value do not alias the field they came
	// from, so this explicit setter is what makes such a mutation observable to later callers -
	// see abstract_timestamp_source.go's identical judgment call for the same reason.
	SetTimestampedReferences(timestampedReferences []*TimestampedReference)

	// SetEmbeddedEvidenceRecordHelper sets a helper for processing and validation of the
	// embedded evidence record type. Port of setEmbeddedEvidenceRecordHelper(EmbeddedEvidenceRecordHelper).
	SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper EmbeddedEvidenceRecordHelper)

	// IsEmbedded returns whether the evidence record is embedded in a signature. Port of isEmbedded().
	IsEmbedded() bool

	// MasterSignature gets a master signature, enveloping the current evidence record. Port of
	// getMasterSignature().
	MasterSignature() AdvancedSignature

	// IncorporationType gets the type of the unsigned attribute used for incorporation of an
	// evidence record within a signature. NOTE: applicable only for embedded evidence records
	// within CAdES. Port of getIncorporationType().
	IncorporationType() enumerations.EvidenceRecordIncorporationType

	// EmbeddedEvidenceRecordHelper returns an EmbeddedEvidenceRecordHelper in case of an
	// embedded evidence record. Port of getEmbeddedEvidenceRecordHelper().
	EmbeddedEvidenceRecordHelper() EmbeddedEvidenceRecordHelper

	// Id returns the DSS unique signature id. It allows to unambiguously identify each
	// signature. Port of getId().
	Id() string

	// Encoded returns binaries of the evidence record document. Port of getEncoded().
	Encoded() []byte
}

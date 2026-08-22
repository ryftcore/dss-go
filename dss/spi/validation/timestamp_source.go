// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampSource.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// TimestampSource is the interface for handling validation data extracted from timestamps.
//
// java.io.Serializable is dropped (no Go counterpart).
type TimestampSource interface {
	// ContentTimestamps returns a list of incorporated content timestamps.
	// Port of getContentTimestamps().
	ContentTimestamps() []*TimestampToken

	// SignatureTimestamps returns a list of incorporated signature timestamps.
	// Port of getSignatureTimestamps().
	SignatureTimestamps() []*TimestampToken

	// TimestampsX1 returns a list of incorporated SigAndRefs timestamps.
	// Port of getTimestampsX1().
	TimestampsX1() []*TimestampToken

	// TimestampsX2 returns a list of incorporated RefsOnly timestamps.
	// Port of getTimestampsX2().
	TimestampsX2() []*TimestampToken

	// ArchiveTimestamps returns a list of incorporated archive timestamps.
	// Port of getArchiveTimestamps().
	ArchiveTimestamps() []*TimestampToken

	// DocumentTimestamps returns a list of incorporated document timestamps (PAdES only).
	// Port of getDocumentTimestamps().
	DocumentTimestamps() []*TimestampToken

	// DetachedTimestamps returns a list of detached timestamps (used in ASiC with CAdES).
	// Port of getDetachedTimestamps().
	DetachedTimestamps() []*TimestampToken

	// AllTimestamps returns a list of all incorporated timestamps.
	// Port of getAllTimestamps().
	AllTimestamps() []*TimestampToken

	// AddExternalTimestamp allows adding an external timestamp. The given timestamp must be
	// processed before. Port of addExternalTimestamp(TimestampToken).
	AddExternalTimestamp(timestamp *TimestampToken)

	// EmbeddedEvidenceRecords returns a list of evidence records embedded in a signature
	// document. Port of getEmbeddedEvidenceRecords().
	EmbeddedEvidenceRecords() []EvidenceRecord

	// DetachedEvidenceRecords returns a list of evidence records detached from a signature
	// document. Port of getDetachedEvidenceRecords().
	DetachedEvidenceRecords() []EvidenceRecord

	// AllEvidenceRecords returns a list of all evidence records associated with the signature.
	// Port of getAllEvidenceRecords().
	AllEvidenceRecords() []EvidenceRecord

	// AddExternalEvidenceRecord allows adding an external evidence record covering the
	// signature file. The given evidence record must be processed before.
	// Port of addExternalEvidenceRecord(EvidenceRecord).
	AddExternalEvidenceRecord(evidenceRecord EvidenceRecord)

	// TimestampCertificateSources returns a merged ListCertificateSource of all embedded
	// timestamp certificate sources. Port of getTimestampCertificateSources().
	TimestampCertificateSources() *spi.ListCertificateSource

	// TimestampCertificateSourcesExceptLastArchiveTimestamp returns a merged
	// ListCertificateSource of all embedded timestamp certificate sources except the latest
	// Archive Timestamp.
	// Port of getTimestampCertificateSourcesExceptLastArchiveTimestamp().
	TimestampCertificateSourcesExceptLastArchiveTimestamp() *spi.ListCertificateSource

	// AllTimestampsExceptLastArchiveTimestamp returns a list of all TimestampTokens except the
	// last archive timestamp. Port of getAllTimestampsExceptLastArchiveTimestamp().
	AllTimestampsExceptLastArchiveTimestamp() []*TimestampToken

	// TimestampCRLSources returns a merged ListRevocationSource of all embedded timestamp CRL
	// sources. Port of getTimestampCRLSources().
	TimestampCRLSources() *spi.ListRevocationSource[revocation.CRL]

	// TimestampOCSPSources returns a merged ListRevocationSource of all embedded timestamp OCSP
	// sources. Port of getTimestampOCSPSources().
	TimestampOCSPSources() *spi.ListRevocationSource[revocation.OCSP]

	// UnsignedPropertiesReferences returns a list of TimestampedReferences for all tokens
	// embedded into unsigned properties of the signature.
	// Port of getUnsignedPropertiesReferences().
	UnsignedPropertiesReferences() []*TimestampedReference

	// SignerDataReferences returns a list of TimestampedReferences obtained from the
	// signatureScopes. Port of getSignerDataReferences().
	SignerDataReferences() []*TimestampedReference

	// IsTimestamped checks whether the token with the given tokenId and objectType is covered by
	// the timestamp source. Port of isTimestamped(String, TimestampedObjectType).
	IsTimestamped(tokenId string, objectType enumerations.TimestampedObjectType) bool
}

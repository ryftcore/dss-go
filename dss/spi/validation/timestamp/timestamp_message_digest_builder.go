// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/TimestampMessageDigestBuilder.java (DSS 6.5.RC1).
package timestamp

import "github.com/utain/esig/dss/model"

// TimestampMessageDigestBuilder builds message-imprint digest to be timestamped.
type TimestampMessageDigestBuilder interface {
	// ContentTimestampMessageDigest returns the content timestamp message-imprint digest
	// (timestamped or to be). Port of getContentTimestampMessageDigest().
	ContentTimestampMessageDigest() model.DSSMessageDigest

	// SignatureTimestampMessageDigest returns the message-imprint digest on data (signature
	// value) that was timestamped by the SignatureTimeStamp for the given timestamp.
	// Port of getSignatureTimestampMessageDigest().
	SignatureTimestampMessageDigest() model.DSSMessageDigest

	// TimestampX1MessageDigest returns the message-imprint digest to be time-stamped. The data
	// used to create digest contains the digital signature (XAdES example: ds:SignatureValue
	// element), the signature time-stamp(s) present in the AdES-T form, the certification path
	// references and the revocation status references.
	// Port of getTimestampX1MessageDigest().
	TimestampX1MessageDigest() model.DSSMessageDigest

	// TimestampX2MessageDigest returns the data to be time-stamped which contains the
	// concatenation of CompleteCertificateRefs and CompleteRevocationRefs elements (XAdES
	// example). Port of getTimestampX2MessageDigest().
	TimestampX2MessageDigest() model.DSSMessageDigest

	// ArchiveTimestampMessageDigest returns the message digest for an archive timestamp.
	// Archive timestamp seals the data of the signature in a specific order; the data must be
	// retrieved for each timestamp. Port of getArchiveTimestampMessageDigest().
	ArchiveTimestampMessageDigest() model.DSSMessageDigest
}

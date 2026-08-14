// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPRefIdentifier.java (DSS 6.5.RC1).
//
// DEVIATION: when the reference carries no digest the identifier is computed over, among
// other things, DataOutputStream#writeChars(responderId.getX500Principal().toString()).
// model.X500Principal#String() renders the RFC 2253 form where the JDK's toString() renders
// the RFC 1779 flavoured "DEFAULT" form (a documented dss-model deviation), so identifiers
// built from a by-name responder without a digest differ from the Java ones. References
// that carry a digest - which is the case for every reference read from a signature - are
// unaffected.
package spi

import (
	"encoding/binary"
	"unicode/utf16"

	"github.com/utain/esig/dss/model"
)

// OCSPRefIdentifier is the identifier of an OCSP token reference.
type OCSPRefIdentifier struct {
	RevocationRefIdentifier
}

// NewOCSPRefIdentifier builds the identifier of the given reference.
// Port of the protected OCSPRefIdentifier(OCSPRef) constructor.
//
// Java wraps the IOException its ByteArrayOutputStream can never raise in a
// DSSException("Cannot build DSS ID for the OCSP Ref."); the Go counterpart of that dead
// branch is a panic carrying the same message, raised only if the SHA-256 implementation
// the identifier digest needs is unavailable.
func NewOCSPRefIdentifier(ocspRef *OCSPRef) *OCSPRefIdentifier {
	return &OCSPRefIdentifier{
		NewRevocationRefIdentifierFromDigest("OCSPRefIdentifier", ocspRefIdentifierDigest(ocspRef)),
	}
}

// ocspRefIdentifierDigest ports the private static getDigest(OCSPRef).
func ocspRefIdentifierDigest(ocspRef *OCSPRef) model.Digest {
	if ocspRef.Digest().Value() != nil {
		return ocspRef.Digest()
	}

	var bytes []byte
	if producedAt := ocspRef.ProducedAt(); !producedAt.IsZero() {
		// DataOutputStream#writeLong of Date#getTime(), i.e. the big-endian milliseconds
		// since the epoch.
		bytes = binary.BigEndian.AppendUint64(bytes, uint64(producedAt.UnixMilli()))
	}
	if responderId := ocspRef.ResponderId(); responderId != nil {
		if ski := responderId.Ski(); ski != nil {
			bytes = append(bytes, ski...)
		}
		if principal := responderId.X500Principal(); principal != nil {
			// DataOutputStream#writeChars, i.e. the UTF-16BE code units of the name.
			for _, unit := range utf16.Encode([]rune(principal.String())) {
				bytes = binary.BigEndian.AppendUint16(bytes, unit)
			}
		}
	}

	value, err := DSSUtilsDigest(model.IdentifierDigestAlgorithm, bytes)
	if err != nil {
		panic("Cannot build DSS ID for the OCSP Ref.")
	}
	return model.NewDigest(model.IdentifierDigestAlgorithm, value)
}

// compile-time assertion: an OCSPRefIdentifier is an identifier.
var _ model.Identifier = (*OCSPRefIdentifier)(nil)

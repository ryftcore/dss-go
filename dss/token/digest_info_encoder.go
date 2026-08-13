// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/digest/DigestInfoEncoder.java (DSS 6.5.RC1).
package token

import (
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/model"
)

// DigestInfoEncoderEncode encodes the algorithmOid and digest combination into its ASN.1
// DigestInfo representation:
//
//	DigestInfo ::= SEQUENCE {
//	   digestAlgorithm DigestAlgorithmIdentifier,
//	   digest Digest }
//
// This class is used to encode a given digest to its ASN.1 DigestInfo representation. NOTE:
// This is used on RSA signing.
//
// Panics with the Java message if algorithmOid is "" or digest is nil (Objects.requireNonNull).
// A malformed OID - Java's IllegalArgumentException, caught and rewrapped as a DSSException - is
// returned as an error.
//
// DEVIATION: the BER/DER engine of internal/asn1ber builds the SEQUENCE/OCTET STRING/NULL/OID
// encodings instead of the manual byte-buffer arithmetic upstream performs (see PORTING.md); the
// produced bytes are identical since both write the same DER TLV structure.
func DigestInfoEncoderEncode(algorithmOid string, digest []byte) ([]byte, error) {
	if algorithmOid == "" {
		panic("Digest algorithm OID cannot be null!")
	}
	if digest == nil {
		panic("Digest cannot be null!")
	}

	oid, err := asn1ber.OIDFromString(algorithmOid)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("An error occurred on DigestInfo encoding : %s", err.Error()), err)
	}
	oidTLV := asn1ber.EncodeOID(oid)
	if oidTLV == nil {
		return nil, model.NewDSSError(
			fmt.Sprintf("An error occurred on DigestInfo encoding : the given string '%s' does not represent "+
				"a valid OID! OID have two or more parts separated by a dot.", algorithmOid))
	}

	algorithmIdentifierContent := append(append([]byte{}, oidTLV...), asn1ber.DERNull...)
	algorithmIdentifier := asn1ber.WriteSequence(algorithmIdentifierContent)

	digestInfoContent := append(algorithmIdentifier, asn1ber.WriteTLV(asn1ber.TagOctetString, digest)...)
	return asn1ber.WriteSequence(digestInfoContent), nil
}

// DigestInfoEncoderIsEncoded verifies whether data is ASN.1 DigestInfo encoded.
//
// DEVIATION: ported statement-by-statement from the Java byte-walk, including its implicit
// assumption that every length octet is a single byte (no long-form length encoding): a
// DigestInfo whose SEQUENCE, AlgorithmIdentifier or OCTET STRING length exceeds 127 bytes -
// unreachable for the digest sizes DSS supports - is reported as not encoded, exactly as
// upstream's identical arithmetic would.
func DigestInfoEncoderIsEncoded(data []byte) (result bool) {
	defer func() {
		// Java catches any Exception from the byte-walk (e.g. ArrayIndexOutOfBoundsException),
		// logs it, and returns false; a Go index panic is the equivalent failure mode.
		if r := recover(); r != nil {
			result = false
		}
	}()

	if data == nil || len(data) < 5 {
		// Minimum length check (SEQUENCE + AlgorithmIdentifier + Digest)
		return false
	}

	index := 0

	// Check for SEQUENCE (0x30)
	if data[index] != 0x30 {
		return false
	}
	index++

	// Get SEQUENCE length
	seqLength := int(data[index]) & 0xFF
	index++
	if seqLength != len(data)-2 {
		// Ensure length matches the rest of the data
		return false
	}

	// Check for AlgorithmIdentifier SEQUENCE (0x30)
	if data[index] != 0x30 {
		return false
	}
	index++

	// Get AlgorithmIdentifier length
	algorithmIdentifierLength := int(data[index]) & 0xFF
	index++
	algorithmIdentifierEnd := index + algorithmIdentifierLength

	// Check for OID (0x06)
	if data[index] != 0x06 {
		return false
	}
	index++

	// Get OID length
	oidLength := int(data[index]) & 0xFF
	index++
	if index+oidLength > algorithmIdentifierEnd {
		// Ensure OID length is valid
		return false
	}
	index += oidLength // Skip OID content

	// Check for NULL (optional, 0x05 followed by 0x00)
	if index < algorithmIdentifierEnd {
		if data[index] == 0x05 && index+1 < algorithmIdentifierEnd && data[index+1] == 0x00 {
			index += 2 // Skip NULL
		} else {
			return false
		}
	}

	// Ensure we are at the end of AlgorithmIdentifier
	if index != algorithmIdentifierEnd {
		return false
	}

	// Check for Digest (OCTET STRING, 0x04)
	if index >= len(data) || data[index] != 0x04 {
		return false
	}
	index++

	// Get Digest length
	digestLength := int(data[index]) & 0xFF
	index++
	if index+digestLength != len(data) {
		// Ensure Digest length matches the rest of the data
		return false
	}

	// Everything checks out
	return true
}

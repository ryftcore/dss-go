// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/DataIdentifier.java (DSS 6.5.RC1).
package model

// DataIdentifier is the DSS identifier for a SignedData.
type DataIdentifier struct {
	IdentifierBase
}

// NewDataIdentifier builds a DataIdentifier over the signed data binaries.
func NewDataIdentifier(data []byte) *DataIdentifier {
	return &DataIdentifier{NewIdentifierBase("DataIdentifier", "D-", data)}
}

// NewDataIdentifierForDocument builds an identifier from a document name and the document's
// content. Port of DataIdentifier(String, DSSDocument).
func NewDataIdentifierForDocument(name string, document DSSDocument) (*DataIdentifier, error) {
	data, err := dataIdentifierBuild(name, document)
	if err != nil {
		return nil, err
	}
	return NewDataIdentifier(data), nil
}

// dataIdentifierBuild builds the byte array for a document with the given name: the name
// written as UTF-16BE characters (java.io.DataOutputStream#writeChars) followed by the
// document's digest value. Port of the private static build(String, DSSDocument).
//
// Java writes nothing for a null name; the Go port writes nothing for an empty name, which
// produces the same bytes.
//
// Upstream wraps an IOException from the ByteArrayOutputStream/DataOutputStream pair in a
// DSSException reading "Unable to build a JAdESAttributeIdentifier. Reason : %s". Neither
// stream can fail, and their Go counterparts (slice appends) cannot fail either, so that
// branch has no counterpart here.
func dataIdentifierBuild(name string, document DSSDocument) ([]byte, error) {
	data := make([]byte, 0)
	// DataOutputStream#writeChars writes every char as two big-endian bytes. Go strings are
	// UTF-8, so they are converted to UTF-16 code units first to reproduce the same bytes,
	// including the surrogate pairs Java would emit for supplementary characters.
	for _, unit := range dataIdentifierUTF16(name) {
		data = append(data, byte(unit>>8), byte(unit))
	}
	dataDigest, present, err := dataIdentifierGetDigest(document)
	if err != nil {
		return nil, err
	}
	if present {
		data = append(data, dataDigest.Value()...)
	}
	return data, nil
}

// dataIdentifierUTF16 converts a Go string into the UTF-16 code units a Java String holds.
func dataIdentifierUTF16(s string) []uint16 {
	units := make([]uint16, 0, len(s))
	for _, r := range s {
		if r > 0xFFFF {
			r -= 0x10000
			units = append(units, uint16(0xD800+(r>>10)), uint16(0xDC00+(r&0x3FF)))
			continue
		}
		units = append(units, uint16(r))
	}
	return units
}

// dataIdentifierGetDigest gets the digest of the document: the digest a DigestDocument
// already carries, otherwise the SHA-256 digest of its content. Port of the private static
// getDigest(DSSDocument); the boolean reports whether upstream would have returned a Digest
// rather than null.
func dataIdentifierGetDigest(document DSSDocument) (Digest, bool, error) {
	if document == nil {
		return Digest{}, false, nil
	}
	if digestDocument, ok := document.(*DigestDocument); ok {
		digest, err := digestDocument.ExistingDigest()
		if err != nil {
			return Digest{}, false, err
		}
		return digest, true, nil
	}
	value, err := document.DigestValue(IdentifierDigestAlgorithm)
	if err != nil {
		// Propagated unchanged: upstream's build() catches IOException only, and
		// getDigestValue raises a DSSException ("Unable to compute the digest"), so the
		// "Unable to build a JAdESAttributeIdentifier" wrapper never applies to this path.
		return Digest{}, false, err
	}
	return NewDigest(IdentifierDigestAlgorithm, value), true, nil
}

// compile-time interface assertion.
var _ Identifier = (*DataIdentifier)(nil)

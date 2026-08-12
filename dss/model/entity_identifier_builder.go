// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/EntityIdentifierBuilder.java (DSS 6.5.RC1).
package model

// EntityIdentifierBuilder builds an EntityIdentifier for a public key and subject name pair.
type EntityIdentifierBuilder struct {
	// publicKey is the public key part of the entity key.
	publicKey *PublicKey
	// subjectName is the subject name part of the entity key.
	subjectName *X500Principal
}

// NewEntityIdentifierBuilder creates the builder for the given key and subject name; either
// may be nil, in which case it contributes nothing to the identifier binaries.
func NewEntityIdentifierBuilder(publicKey *PublicKey, subjectName *X500Principal) *EntityIdentifierBuilder {
	return &EntityIdentifierBuilder{publicKey: publicKey, subjectName: subjectName}
}

// Build builds the EntityIdentifier.
//
// Java declares build() on the IdentifierBuilder interface as returning an Identifier and
// narrows it here to EntityIdentifier. Go has no covariant returns, so this method keeps the
// useful concrete type and EntityIdentifierBuilder consequently does not satisfy the
// IdentifierBuilder interface.
func (b *EntityIdentifierBuilder) Build() *EntityIdentifier {
	return NewEntityIdentifier(b.BuildBinaries())
}

// BuildBinaries builds the unique binary data describing the entity key: the DER
// SubjectPublicKeyInfo of the public key followed by the DER encoding of the subject name.
// Port of the protected buildBinaries().
func (b *EntityIdentifierBuilder) BuildBinaries() []byte {
	// The ByteArrayOutputStream upstream writes to cannot fail, so its IOException branch
	// (which raises "An error occurred while building an Identifier : %s") has no Go
	// counterpart.
	binaries := make([]byte, 0)
	if b.publicKey != nil {
		binaries = append(binaries, b.publicKey.Encoded()...)
	}
	if b.subjectName != nil {
		binaries = append(binaries, b.subjectName.Encoded()...)
	}
	return binaries
}

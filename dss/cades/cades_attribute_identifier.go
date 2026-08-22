// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESAttributeIdentifier.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi/validation/identifier"
)

// AttributeIdentifier represents a unique identifier for an attribute from a CAdES
// signature. Port of the class CAdESAttributeIdentifier, extending SignatureAttributeIdentifier.
type AttributeIdentifier struct {
	identifier.SignatureAttributeIdentifier
}

// newCAdESAttributeIdentifier is the port of the package-private
// AttributeIdentifier(byte[]) constructor.
func newAttributeIdentifier(data []byte) *AttributeIdentifier {
	return &AttributeIdentifier{
		SignatureAttributeIdentifier: identifier.NewSignatureAttributeIdentifierBase("CAdESAttributeIdentifier", data),
	}
}

// AttributeIdentifierBuild builds the identifier for a CAdES attribute.
// Port of the static build(Attribute, Integer) method.
//
// Java's DataOutputStream#write(int) call on the order argument writes only its low-order
// byte, not the full int value; that is reproduced here with a single order&0xFF byte, rather
// than "fixed" to encode the whole integer. The DSSException build() raises on an IOException
// cannot occur here, since building the byte slice below cannot fail; there is accordingly no
// error return.
func AttributeIdentifierBuild(attribute *cmscore.Attribute, order *int) *AttributeIdentifier {
	var data []byte
	if attribute != nil {
		// attribute identifier + value
		data = append(data, attribute.Encoded()...)
	}
	if order != nil {
		data = append(data, byte(*order))
	}
	return newAttributeIdentifier(data)
}

// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESAttributeIdentifier.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/spi/validation/identifier"
)

// CAdESAttributeIdentifier represents a unique identifier for an attribute from a CAdES
// signature. Port of the class CAdESAttributeIdentifier, extending SignatureAttributeIdentifier.
type CAdESAttributeIdentifier struct {
	identifier.SignatureAttributeIdentifier
}

// newCAdESAttributeIdentifier is the port of the package-private
// CAdESAttributeIdentifier(byte[]) constructor.
func newCAdESAttributeIdentifier(data []byte) *CAdESAttributeIdentifier {
	return &CAdESAttributeIdentifier{
		SignatureAttributeIdentifier: identifier.NewSignatureAttributeIdentifierBase("CAdESAttributeIdentifier", data),
	}
}

// CAdESAttributeIdentifierBuild builds the identifier for a CAdES attribute.
// Port of the static build(Attribute, Integer) method.
//
// Java's DataOutputStream#write(int) call on the order argument writes only its low-order
// byte, not the full int value; that is reproduced here with a single order&0xFF byte, rather
// than "fixed" to encode the whole integer. The DSSException build() raises on an IOException
// cannot occur here, since building the byte slice below cannot fail; there is accordingly no
// error return.
func CAdESAttributeIdentifierBuild(attribute *cmscore.Attribute, order *int) *CAdESAttributeIdentifier {
	var data []byte
	if attribute != nil {
		// attribute identifier + value
		data = append(data, attribute.Encoded()...)
	}
	if order != nil {
		data = append(data, byte(*order))
	}
	return newCAdESAttributeIdentifier(data)
}

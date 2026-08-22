// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSigProperties.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/internal/cmscore"

// CAdESSigProperties represents a list of CAdESAttributes. Port of the abstract class
// SigProperties, implementing spi/validation.SignatureProperties[*Attribute].
//
// Java distinguishes a null ASN1Set (the signedAttrs/unsignedAttrs field absent) from a
// present, possibly empty one; "exists" is threaded explicitly here since a nil
// cmscore.Attributes cannot on its own distinguish the two (see cades_certificate_source.go and
// friends on the same point) - SignedAttributes.Build/UnsignedAttributes.Build pass
// SignerInfo's HasSignedAttributes()/HasUnsignedAttributes().
type SigProperties struct {
	// attributeTable is the CMS attribute table set.
	attributeTable cmscore.Attributes
	// exists records whether the field was present.
	exists bool
}

// newCAdESSigProperties is the port of the package-private CAdESSigProperties(ASN1Set)
// constructor.
func newCAdESSigProperties(attributeTable cmscore.Attributes, exists bool) SigProperties {
	return SigProperties{attributeTable: attributeTable, exists: exists}
}

// IsExist checks if "unsigned-signature-properties"/"signed-signature-properties" exists and
// can be processed. Port of isExist().
func (p *SigProperties) IsExist() bool {
	return p.exists
}

// Attributes returns a list of children contained in the element. Port of getAttributes().
func (p *SigProperties) Attributes() []*Attribute {
	attributes := make([]*Attribute, 0)
	if p.IsExist() {
		for index, attribute := range p.attributeTable {
			order := index
			attributes = append(attributes, NewCAdESAttribute(attribute, &order))
		}
	}
	return attributes
}

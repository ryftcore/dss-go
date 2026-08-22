// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignedAttributes.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/internal/cmscore"

// CAdESSignedAttributes represents the CAdES Signed attributes. Port of the class
// SignedAttributes, extending SigProperties.
type SignedAttributes struct {
	SigProperties
}

// newCAdESSignedAttributes is the port of the package-private CAdESSignedAttributes(ASN1Set)
// constructor.
func newCAdESSignedAttributes(attributeTable cmscore.Attributes, exists bool) *SignedAttributes {
	return &SignedAttributes{SigProperties: newCAdESSigProperties(attributeTable, exists)}
}

// CAdESSignedAttributesBuild builds the CAdESSignedAttributes from a SignerInfo. Port of the
// static build(SignerInformation).
func SignedAttributesBuild(signerInformation *cmscore.SignerInfo) *SignedAttributes {
	return newCAdESSignedAttributes(signerInformation.SignedAttributes, signerInformation.HasSignedAttributes())
}

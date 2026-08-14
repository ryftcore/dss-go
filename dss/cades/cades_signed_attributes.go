// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignedAttributes.java (DSS 6.5.RC1).
package cades

import "github.com/utain/esig/dss/internal/cmscore"

// CAdESSignedAttributes represents the CAdES Signed attributes. Port of the class
// CAdESSignedAttributes, extending CAdESSigProperties.
type CAdESSignedAttributes struct {
	CAdESSigProperties
}

// newCAdESSignedAttributes is the port of the package-private CAdESSignedAttributes(ASN1Set)
// constructor.
func newCAdESSignedAttributes(attributeTable cmscore.Attributes, exists bool) *CAdESSignedAttributes {
	return &CAdESSignedAttributes{CAdESSigProperties: newCAdESSigProperties(attributeTable, exists)}
}

// CAdESSignedAttributesBuild builds the CAdESSignedAttributes from a SignerInfo. Port of the
// static build(SignerInformation).
func CAdESSignedAttributesBuild(signerInformation *cmscore.SignerInfo) *CAdESSignedAttributes {
	return newCAdESSignedAttributes(signerInformation.SignedAttributes, signerInformation.HasSignedAttributes())
}

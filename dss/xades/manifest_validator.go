// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/ManifestValidator.java
// (DSS 6.5.RC1).
//
// This class validates a ds:Manifest element against external files:
//
//	<ds:Manifest Id="manifest">
//		<ds:Reference URI="l_19420170726bg.pdf">
//			<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/>
//			<ds:DigestValue>EUcwRQ....</ds:DigestValue>
//		</ds:Reference>
//		...
//	</ds:Manifest>
//
// Upstream calls org.apache.xml.security.signature.Manifest/Reference through
// DSSXMLUtils.initManifestWithDetachedContent/initManifestDetachedContent; per
// internal/xmldsig's doc.go table this is internal/xmldsig.Manifest/NewManifest plus a
// DetachedSignatureResolver registered per distinct digest algorithm found in the manifest's
// references - the same pattern XAdESSignature.initDetachedSignatureResolvers already uses for
// ds:SignedInfo (xades_signature.go).
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ManifestValidator validates a ds:Manifest element against external files. Port of the class
// ManifestValidator.
type ManifestValidator struct {
	// manifest is the wrapped internal/xmldsig Manifest.
	manifest *xmldsig.Manifest
}

// NewManifestValidator is the port of the public ManifestValidator(Element, List<DSSDocument>)
// constructor.
//
// Panics with the Java DSSException message when the manifest element cannot be wrapped (Java's
// DSSXMLUtils.initManifestWithDetachedContent declares XMLSecurityException, which upstream
// itself wraps in an unchecked DSSException).
func NewManifestValidator(manifestElement *xmldom.Node, detachedContents []model.DSSDocument) *ManifestValidator {
	manifest, err := xmldsig.NewManifest(manifestElement, nil)
	if err != nil {
		panic("Unable to instantiate a ManifestValidator. Reason : " + err.Error())
	}
	manifestValidatorInitDetachedContent(manifest, detachedContents)
	return &ManifestValidator{manifest: manifest}
}

// NewManifestValidatorFromManifest is the port of the public ManifestValidator(Manifest,
// List<DSSDocument>) constructor.
func NewManifestValidatorFromManifest(manifest *xmldsig.Manifest, detachedContents []model.DSSDocument) *ManifestValidator {
	manifestValidatorInitDetachedContent(manifest, detachedContents)
	return &ManifestValidator{manifest: manifest}
}

// manifestValidatorInitDetachedContent ports the private static initManifestDetachedContent,
// delegating to DSSXMLUtils.initManifestDetachedContent: one DetachedSignatureResolver per
// distinct digest algorithm found among the manifest's ds:Reference/ds:DigestMethod elements.
func manifestValidatorInitDetachedContent(manifest *xmldsig.Manifest, detachedContents []model.DSSDocument) {
	if utils.IsCollectionEmpty(detachedContents) {
		return
	}
	for _, digestAlgorithm := range manifestValidatorReferenceDigestAlgos(manifest) {
		manifest.AddResourceResolver(&xmldsig.DetachedSignatureResolver{
			Documents:       detachedContents,
			DigestAlgorithm: digestAlgorithm,
		})
	}
}

// manifestValidatorReferenceDigestAlgos ports DSSXMLUtilsGetReferenceDigestAlgos for a
// ds:Manifest element (already assumed for ds:SignedInfo by xades_signature.go's
// initDetachedSignatureResolvers; the same helper walks any element's ds:Reference children).
func manifestValidatorReferenceDigestAlgos(manifest *xmldsig.Manifest) []enumerations.DigestAlgorithm {
	return DSSXMLUtilsGetReferenceDigestAlgos(manifest.Element())
}

// Validate validates the manifest and returns a list of ReferenceValidations. Port of validate().
//
// LOG.warn("Unable to verify reference with Id [{}] : {}") is dropped per PORTING.md; a
// reference whose digest/found/duplicated/intact computation fails is still returned in the
// list, just without the fields that failed to compute - matching upstream's catch-and-continue.
func (v *ManifestValidator) Validate() []*model.ReferenceValidation {
	// Upstream logs "Validation of the manifest references ...".

	references, err := v.manifest.References()
	if err != nil || utils.IsCollectionEmpty(references) {
		// Upstream logs "No references found inside the ds:Manifest element!".
		return []*model.ReferenceValidation{}
	}

	referenceValidations := make([]*model.ReferenceValidation, 0, len(references))
	for _, reference := range references {
		refValidation := NewXAdESReferenceValidation(reference)
		refValidation.SetType(enumerations.DigestMatcherTypeManifestEntry)

		referenceValidations = append(referenceValidations, &refValidation.ReferenceValidation)

		refValidation.SetDigest(DSSXMLUtilsGetReferenceDigest(reference))
		refValidation.SetTransformationNames(manifestValidatorGetTransformNames(reference.Element()))

		refFound := DSSXMLUtilsIsAbleToDeReferenceContent(reference)
		refValidation.SetFound(refFound)

		isDuplicated := DSSXMLUtilsIsReferencedContentAmbiguous(v.manifest.Element().OwnerDocument(), refValidation.Uri())
		refValidation.SetDuplicated(isDuplicated)

		if refFound && !isDuplicated {
			if intact, err := reference.Verify(); err == nil {
				refValidation.SetIntact(intact)
			}
			// Upstream logs the caught Exception here on a Verify() failure.
		}
	}
	return referenceValidations
}

// manifestValidatorGetTransformNames ports the private getTransformNames(Element).
func manifestValidatorGetTransformNames(refNode *xmldom.Node) []string {
	transformNames := make([]string, 0)
	nodeList, err := xmlutils.XPathUtilsGetNodeList(refNode, common.XMLDSigPathTransformsTransformPath)
	if err != nil {
		return transformNames
	}
	for _, transformElement := range nodeList {
		algorithm := transformElement.AttrValue("", common.XMLDSigAttributeAlgorithm.AttributeName())
		if utils.IsStringNotBlank(algorithm) {
			transformNames = append(transformNames, algorithm)
		}
	}
	return transformNames
}

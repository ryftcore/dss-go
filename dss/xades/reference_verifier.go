// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceVerifier.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_BASE64_DECODE is
// xmldsig.TransformBase64Decode, per internal/xmldsig's doc.go mapping table.
//
// checkReferencesValidity gains an error return: every IllegalArgumentException it throws
// becomes a returned error (PORTING.md), and DomUtils.buildDOM can fail in Go. slf4j debug and
// warn logging is dropped - including the "points to an XML Node, while no transforms are
// defined" warning, which changes nothing the caller can observe.
package xades

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// referenceVerifierReferenceWrongMessage is the shared prefix of every rejection message. Port
// of the local referenceWrongMessage variable.
const referenceVerifierReferenceWrongMessage = "Reference setting is not correct! "

// ReferenceVerifier is used to verify the validity of the DSSReferences setup.
type ReferenceVerifier struct {
	// signatureParameters are the used XAdESSignatureParameters.
	signatureParameters *XAdESSignatureParameters
}

// NewReferenceVerifier is the default constructor for a signature references verification.
// Ports ReferenceVerifier(XAdESSignatureParameters).
func NewReferenceVerifier(signatureParameters *XAdESSignatureParameters) *ReferenceVerifier {
	return &ReferenceVerifier{signatureParameters: signatureParameters}
}

// CheckReferencesValidity verifies the compatibility of the defined signature parameters and
// the reference transformations. Ports checkReferencesValidity().
//
// NOTE: as upstream, a reference without an Id is given a deterministic one IN PLACE here, so
// this method mutates the parameters' references.
func (v *ReferenceVerifier) CheckReferencesValidity() error {
	if v.signatureParameters == nil {
		return nil
	}

	referenceIdProvider := NewReferenceIdProvider()
	referenceIdProvider.SetSignatureParameters(v.signatureParameters)
	for _, reference := range v.signatureParameters.References() {
		if utils.IsStringEmpty(reference.Id()) {
			// No Id defined for a reference. Generate a deterministic identifier...
			reference.SetId(referenceIdProvider.ReferenceId())
		}
		if reference.Object() != nil {
			// ds:Object is defined for the reference. Use the provided value.
			continue
		}
		if xmlutils.DomUtilsIsElementReference(reference.Uri()) &&
			!DSSXMLUtilsIsObjectReferenceType(reference.Type()) &&
			xmlutils.DomUtilsIsDOM(reference.Contents()) {
			document, err := xmlutils.DomUtilsBuildDOMFromDocument(reference.Contents())
			if err != nil {
				return err
			}
			id := xmlutils.DomUtilsGetId(reference.Uri())
			if xmlutils.XPathUtilsGetElementById(document, id) == nil {
				return fmt.Errorf("An element with Id '%s' has not been found in the provided content!", id)
			}
		}

		transforms := reference.Transforms()
		if utils.IsCollectionNotEmpty(transforms) {
			for _, transform := range transforms {
				if xmldsig.TransformBase64Decode == transform.Algorithm() {
					if reference.Object() == nil || !DSSXMLUtilsIsObjectReferenceType(reference.Type()) {
						switch {
						case v.signatureParameters.IsEmbedXML():
							return errors.New(referenceVerifierReferenceWrongMessage +
								"The embedXML(true) parameter is not compatible with base64 transform.")
						case v.signatureParameters.IsManifestSignature():
							return errors.New(referenceVerifierReferenceWrongMessage +
								"Manifest signature is not compatible with base64 transform.")
						case enumerations.SignaturePackagingEnveloping != v.signatureParameters.SignaturePackaging():
							return fmt.Errorf("%sBase64 transform is not compatible with %s signature format.",
								referenceVerifierReferenceWrongMessage, v.signatureParameters.SignaturePackaging())
						}
					}
					if len(transforms) > 1 {
						return errors.New(referenceVerifierReferenceWrongMessage +
							"Base64 transform cannot be used with other transformations.")
					}
				}
			}

		} else {
			uri := reference.Uri()
			if enumerations.SignaturePackagingEnveloped == v.signatureParameters.SignaturePackaging() &&
				utils.IsStringBlank(uri) {
				return errors.New(referenceVerifierReferenceWrongMessage +
					"Enveloped signature must have an enveloped transformation!")
			}
		}
	}
	return nil
}

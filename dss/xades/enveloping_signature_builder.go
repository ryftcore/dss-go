// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/EnvelopingSignatureBuilder.java (DSS 6.5.RC1).
//
// Java's class is package-private and overrides exactly one of the hooks collected in
// XAdESSignatureBuilderOverrides - incorporateSignedObjects(). Go has no method overriding across
// embedding, so this builder registers itself with InitXAdESSignatureBuilder and every other hook
// is promoted from the embedded base, the TokenBase.InitToken(self) convention of PORTING.md.
//
// DEVIATION (flagged for the integrator, and consistent with detached_signature_builder.go): the
// Java class is package-private but is exported here; see enveloped_signature_builder.go for the
// reasoning.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// EnvelopingSignatureBuilder handles the specifics of the enveloping XML signature.
type EnvelopingSignatureBuilder struct {
	XAdESSignatureBuilder
}

// NewEnvelopingSignatureBuilder is the constructor for a single-document signing. The enveloping
// signature uses by default the inclusive method of canonicalization.
// Port of EnvelopingSignatureBuilder(XAdESSignatureParameters, DSSDocument, CertificateVerifier).
func NewEnvelopingSignatureBuilder(params *XAdESSignatureParameters, document model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *EnvelopingSignatureBuilder {
	return NewEnvelopingSignatureBuilderForDocuments(params, []model.DSSDocument{document}, certificateVerifier)
}

// NewEnvelopingSignatureBuilderForDocuments is the constructor for signing multiple documents.
// The enveloping signature uses by default the inclusive method of canonicalization.
// Port of EnvelopingSignatureBuilder(XAdESSignatureParameters, List<DSSDocument>, CertificateVerifier).
func NewEnvelopingSignatureBuilderForDocuments(params *XAdESSignatureParameters,
	documents []model.DSSDocument,
	certificateVerifier validation.CertificateVerifier) *EnvelopingSignatureBuilder {
	builder := &EnvelopingSignatureBuilder{}
	builder.InitXAdESSignatureBuilder(builder, params, documents, certificateVerifier)
	return builder
}

// IncorporateSignedObjects incorporates one ds:Object per reference: the reference's own
// DSSObject when it carries one, a rebuilt ds:Manifest when manifestSignature is configured, and
// otherwise the referenced content - embedded as XML or base64-encoded.
// Port of the overridden protected #incorporateSignedObjects.
func (b *EnvelopingSignatureBuilder) IncorporateSignedObjects() error {
	references := b.Params.References()
	for _, reference := range references {
		// <ds:Object>
		if reference.Object() != nil {
			if err := b.IncorporateObject(reference.Object()); err != nil {
				return err
			}

		} else if b.Params.IsManifestSignature() {
			doc, err := xmlutils.DomUtilsBuildDOMFromDocument(reference.Contents())
			if err != nil {
				return err
			}
			root := doc.DocumentElement()
			referencesNodes := root.Children()
			idAttribute := root.AttrValue("", common.XMLDSigAttributeID.AttributeName())

			// rebuild manifest element to avoid namespace duplication
			manifestDom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
				common.XMLDSigElementManifest)
			manifestDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, idAttribute)
			for _, referenceNode := range referencesNodes {
				copyNode := b.DocumentDom.Import(referenceNode, true)
				manifestDom.AppendChild(copyNode)
			}

			dom := xmlutils.DomUtilsCreateElementNS(b.DocumentDom, b.overrides.XmldsigNamespace(),
				common.XMLDSigElementObject)
			dom.AppendChild(manifestDom)
			b.SignatureDom.AppendChild(dom)

		} else {
			object := NewDSSObject()

			var content model.DSSDocument
			if b.Params.IsEmbedXML() {
				content = reference.Contents()
			} else {
				raw, err := spi.DSSUtilsToByteArrayOfDocument(reference.Contents())
				if err != nil {
					return err
				}
				base64EncodedOriginalDocument := utils.ToBase64(raw)
				content = model.NewInMemoryDocument([]byte(base64EncodedOriginalDocument))
			}
			object.SetContent(content)
			object.SetId(reference.Uri()[1:])

			if err := b.IncorporateObject(object); err != nil {
				return err
			}
		}
	}
	return nil
}

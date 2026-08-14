// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/ManifestBuilder.java (DSS 6.5.RC1).
//
// Java offers six constructors; Go has no overloading, so each gets its own name, and the ones
// that only supply a default delegate to the fuller form exactly as Java's this(...) chains do:
//
//	ManifestBuilder(DigestAlgorithm, List<DSSDocument>)                     -> NewManifestBuilder
//	ManifestBuilder(String, DigestAlgorithm, List<DSSDocument>)             -> NewManifestBuilderWithId
//	ManifestBuilder(String, DigestAlgorithm, List<DSSDocument>, DSSNamespace)
//	                                                                       -> NewManifestBuilderWithNamespace
//	ManifestBuilder(List<DSSReference>)                                     -> NewManifestBuilderWithReferences
//	ManifestBuilder(String, List<DSSReference>)                             -> NewManifestBuilderWithIdAndReferences
//	ManifestBuilder(String, List<DSSReference>, DSSNamespace)               -> NewManifestBuilderWithReferencesAndNamespace
//
// The IllegalArgumentException raised by the empty-list checks - thrown from a constructor -
// becomes a returned error, so every constructor returns (T, error) (PORTING.md).
//
// The ds:Reference elements are written by the same ReferenceProcessor the signature builder
// uses, so a Manifest reference and a SignedInfo reference are produced by one code path, as
// upstream intends.
package xades

import (
	"errors"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// manifestBuilderDefaultManifestID defines the default id for the Manifest element when none is
// provided. Port of the private DEFAULT_MANIFEST_ID.
const manifestBuilderDefaultManifestID = "manifest"

// manifestBuilderDefaultNamespace is the namespace the Manifest is built in when none is given.
// Port of the private DEFAULT_NAMESPACE.
var manifestBuilderDefaultNamespace = common.XMLDSigNS

// ManifestBuilder builds a ds:Manifest element:
//
//	<ds:Manifest Id="manifest">
//	    <ds:Reference URI="l_19420170726bg.pdf">
//	        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/>
//	        <ds:DigestValue>EUcwRQ....</ds:DigestValue>
//	    </ds:Reference>
//	    <ds:Reference URI="l_19420170726cs.pdf">
//	        <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/>
//	        <ds:DigestValue>NQNnr+F...</ds:DigestValue>
//	    </ds:Reference>
//	    ...
//	</ds:Manifest>
type ManifestBuilder struct {
	// manifestId is the manifest id.
	manifestId string

	// references is the list of references to be incorporated into the Manifest.
	references []*DSSReference

	// xmldsigNamespace is the namespace.
	xmldsigNamespace *common.DSSNamespace
}

// NewManifestBuilder is the constructor for the builder; the Id of the Manifest tag will be
// "manifest". Port of ManifestBuilder(DigestAlgorithm, List<DSSDocument>).
func NewManifestBuilder(digestAlgorithm enumerations.DigestAlgorithm,
	documents []model.DSSDocument) (*ManifestBuilder, error) {
	return NewManifestBuilderWithId(manifestBuilderDefaultManifestID, digestAlgorithm, documents)
}

// NewManifestBuilderWithId is the constructor for the builder taking the Id of the Manifest tag.
// Port of ManifestBuilder(String, DigestAlgorithm, List<DSSDocument>).
func NewManifestBuilderWithId(manifestId string, digestAlgorithm enumerations.DigestAlgorithm,
	documents []model.DSSDocument) (*ManifestBuilder, error) {
	return NewManifestBuilderWithNamespace(manifestId, digestAlgorithm, documents,
		manifestBuilderDefaultNamespace)
}

// NewManifestBuilderWithNamespace is the constructor for the builder taking the Id of the
// Manifest tag and the xmldsig namespace definition.
// Port of ManifestBuilder(String, DigestAlgorithm, List<DSSDocument>, DSSNamespace).
func NewManifestBuilderWithNamespace(manifestId string, digestAlgorithm enumerations.DigestAlgorithm,
	documents []model.DSSDocument, xmldsigNamespace *common.DSSNamespace) (*ManifestBuilder, error) {
	references, err := manifestBuilderCreateReferences(manifestId, digestAlgorithm, documents)
	if err != nil {
		return nil, err
	}
	return NewManifestBuilderWithReferencesAndNamespace(manifestId, references, xmldsigNamespace)
}

// NewManifestBuilderWithReferences is the constructor with custom references and the default
// manifest id. Port of ManifestBuilder(List<DSSReference>).
func NewManifestBuilderWithReferences(references []*DSSReference) (*ManifestBuilder, error) {
	return NewManifestBuilderWithIdAndReferences(manifestBuilderDefaultManifestID, references)
}

// NewManifestBuilderWithIdAndReferences is the constructor with custom references and the default
// namespace. Port of ManifestBuilder(String, List<DSSReference>).
func NewManifestBuilderWithIdAndReferences(manifestId string,
	references []*DSSReference) (*ManifestBuilder, error) {
	return NewManifestBuilderWithReferencesAndNamespace(manifestId, references,
		manifestBuilderDefaultNamespace)
}

// NewManifestBuilderWithReferencesAndNamespace is the constructor with custom references and a
// custom namespace. Port of ManifestBuilder(String, List<DSSReference>, DSSNamespace).
func NewManifestBuilderWithReferencesAndNamespace(manifestId string, references []*DSSReference,
	xmldsigNamespace *common.DSSNamespace) (*ManifestBuilder, error) {
	if utils.IsCollectionEmpty(references) {
		return nil, errors.New("List of references cannot be empty!")
	}
	return &ManifestBuilder{
		manifestId:       manifestId,
		references:       references,
		xmldsigNamespace: xmldsigNamespace,
	}, nil
}

// manifestBuilderCreateReferences ports the private static createReferences.
func manifestBuilderCreateReferences(manifestId string, digestAlgorithm enumerations.DigestAlgorithm,
	documents []model.DSSDocument) ([]*DSSReference, error) {
	if utils.IsCollectionEmpty(documents) {
		return nil, errors.New("List of documents cannot be empty!")
	}
	referenceIdProvider := NewReferenceIdProvider()
	referenceIdProvider.SetReferenceIdPrefix("r-" + manifestId)
	referenceBuilder := NewReferenceBuilderWithDigestAlgorithm(documents, digestAlgorithm, referenceIdProvider)
	return referenceBuilder.Build()
}

// Build builds the Manifest. Port of #build().
func (b *ManifestBuilder) Build() (model.DSSDocument, error) {
	documentDom := xmlutils.DomUtilsBuildDOMEmpty()

	manifestDom := xmlutils.DomUtilsCreateElementNS(documentDom, b.xmldsigNamespace,
		common.XMLDSigElement_MANIFEST)
	manifestDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ID.AttributeName()}, b.manifestId)
	documentDom.AppendChild(manifestDom)

	referenceProcessor := NewReferenceProcessorEmpty()
	if err := referenceProcessor.IncorporateReferences(manifestDom, b.references, b.xmldsigNamespace); err != nil {
		return nil, err
	}

	return xmlutils.DomUtilsCreateDssDocumentFromDomDocument(documentDom, b.manifestId)
}

// ManifestReferences returns the list of DSSReferences.
// Port of #getManifestReferences.
func (b *ManifestBuilder) ManifestReferences() []*DSSReference {
	return b.references
}

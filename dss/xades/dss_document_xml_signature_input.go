// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/DSSDocumentXMLSignatureInput.java
// (DSS 6.5.RC1).
//
// org.apache.xml.security.signature.XMLSignatureInput is replaced wholesale by
// internal/xmldsig.Data per its doc.go mapping table ("Data <- signature.XMLSignatureInput").
// This Java class's actual functional role - wrapping a model.DSSDocument (or, via the
// protected 2-arg constructor DigestDocumentXMLSignatureInput calls, its digest) into an
// XMLSignatureInput for DetachedSignatureResolver#engineResolveURI, buffering the document
// once and optionally short-circuiting to a pre-calculated digest when the reference carries no
// ds:Transforms - has already been ported directly into the frozen
// internal/xmldsig/resolver_detached.go's DetachedSignatureResolver.Resolve (its own header
// names both this file and DigestDocumentXMLSignatureInput as sources), which builds
// xmldsig.Data values inline via NewOctetData/NewPreCalculatedDigestData rather than through an
// intermediate wrapper type - Data has no exported Document backreference to attach one to, and
// internal/xmldsig is frozen (PORTING.md: no edits to frozen packages).
//
// This type is still ported in full, one Go file per Java file per PORTING.md, mirroring the
// Java public surface (the document-carrying constructor, MIMEType, the
// IsPreCalculatedDigest/PreCalculatedDigest/SetPreCalculatedDigest override triple) on top of an
// embedded *xmldsig.Data built the same way the frozen resolver builds one.
//
// FLAGGED FOR INTEGRATOR: DSSXMLUtils#getDocument(Reference), assumed elsewhere in this package
// (xades_reference_validation.go's DSSXMLUtilsGetDocument call, landed by a sibling chunk) as
// the Go analogue of "type-assert the reference's XMLSignatureInput to
// DSSDocumentXMLSignatureInput and read its document", cannot be implemented that way here:
// the frozen resolver's Data values never pass through this wrapper, so there is no
// Data->DSSDocument backreference to recover post hoc. Whoever lands DSSXMLUtils.java (not in
// this manifest) needs a different strategy - e.g. having the reference-resolution call site
// re-run DetachedSignatureResolver's own candidate selection, or threading the resolved
// model.DSSDocument through xmldsig.ResolverContext/Reference some other way - since the
// type-assertion Java performs has no Go equivalent against the frozen Data type.
package xades

import (
	"fmt"
	"io"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// DSSDocumentXMLSignatureInput is an implementation of an XMLSignatureInput created on a base
// of a DSSDocument. Port of the class DSSDocumentXMLSignatureInput, extending
// org.apache.xml.security.signature.XMLSignatureInput (internal/xmldsig.Data here).
type DSSDocumentXMLSignatureInput struct {
	*xmldsig.Data

	// document is the detached document to be provided.
	document model.DSSDocument

	// preCalculatedDigest is the pre-calculated digest value of the object in base64, empty
	// when none has been set (Java: null field).
	preCalculatedDigest string
	// preCalculatedDigestSet distinguishes "no pre-calculated digest" from an (unlikely)
	// empty-string one, since Go strings have no null.
	preCalculatedDigestSet bool
}

// NewDSSDocumentXMLSignatureInput is the default constructor for an XMLSignatureInput from a
// detached document. Port of the public DSSDocumentXMLSignatureInput(DSSDocument) constructor.
func NewDSSDocumentXMLSignatureInput(document model.DSSDocument) (*DSSDocumentXMLSignatureInput, error) {
	octets, err := dssDocumentXMLSignatureInputToOctets(document)
	if err != nil {
		return nil, err
	}
	return &DSSDocumentXMLSignatureInput{
		Data:     xmldsig.NewOctetData(octets),
		document: document,
	}, nil
}

// dssDocumentXMLSignatureInputToOctets ports the private static toInputStream(DSSDocument),
// buffered eagerly the same way the frozen internal/xmldsig/resolver_detached.go's Resolve
// reads a detached document (see that file's header for why the buffered and streamed paths
// are byte-identical for XML-DSig purposes).
func dssDocumentXMLSignatureInputToOctets(document model.DSSDocument) ([]byte, error) {
	rc, err := document.OpenStream()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// newDSSDocumentXMLSignatureInputWithDigest is the port of the protected
// DSSDocumentXMLSignatureInput(DSSDocument, DigestAlgorithm) constructor, for an
// XMLSignatureInput built from a base64-encoded document digest. Used by
// DigestDocumentXMLSignatureInput below.
func newDSSDocumentXMLSignatureInputWithDigest(document model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) (*DSSDocumentXMLSignatureInput, error) {
	base64Digest, err := dssDocumentXMLSignatureInputGetBase64Digest(document, digestAlgorithm)
	if err != nil {
		return nil, err
	}
	data := xmldsig.NewPreCalculatedDigestData(base64Digest)
	return &DSSDocumentXMLSignatureInput{
		Data:     data,
		document: document,
		// preCalculatedDigest = super.getPreCalculatedDigest(): the digest provided in the
		// XMLSignatureInput(String) constructor, copied into this class's own field.
		preCalculatedDigest:    data.PreCalculatedDigest(),
		preCalculatedDigestSet: true,
	}, nil
}

// dssDocumentXMLSignatureInputGetBase64Digest ports the private static
// getBase64Digest(DSSDocument, DigestAlgorithm).
func dssDocumentXMLSignatureInputGetBase64Digest(document model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) (string, error) {
	digestValue, err := document.DigestValue(digestAlgorithm)
	if err != nil {
		return "", fmt.Errorf("xades: unable to compute digest for pre-calculated XMLSignatureInput: %w", err)
	}
	return utils.ToBase64(digestValue), nil
}

// MIMEType is the port of the getMIMEType() override.
func (in *DSSDocumentXMLSignatureInput) MIMEType() string {
	if mimeType := in.document.MimeType(); mimeType != nil {
		return mimeType.MimeTypeString()
	}
	return ""
}

// Document returns the wrapped document. Port of the public getDocument().
func (in *DSSDocumentXMLSignatureInput) Document() model.DSSDocument {
	return in.document
}

// IsPreCalculatedDigest is the port of the isPreCalculatedDigest() override, reading this
// class's own preCalculatedDigest field rather than the base XMLSignatureInput's.
func (in *DSSDocumentXMLSignatureInput) IsPreCalculatedDigest() bool {
	return in.preCalculatedDigestSet
}

// PreCalculatedDigest is the port of the getPreCalculatedDigest() override.
//
// Java closes the original InputStream first (Utils.closeQuietly(getOctetStreamReal())); this
// port's underlying octets were already read to completion by the buffered constructors above,
// so there is no open stream left to close.
func (in *DSSDocumentXMLSignatureInput) PreCalculatedDigest() string {
	return in.preCalculatedDigest
}

// SetPreCalculatedDigest sets the pre-calculated digest to avoid document streaming. Port of
// the public setPreCalculatedDigest(String).
func (in *DSSDocumentXMLSignatureInput) SetPreCalculatedDigest(preCalculatedDigest string) {
	in.preCalculatedDigest = preCalculatedDigest
	in.preCalculatedDigestSet = true
}

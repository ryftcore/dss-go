// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/DSSReference.java (DSS 6.5.RC1).
//
// java.io.Serializable and serialVersionUID are dropped.
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// DSSReference defines a ds:Reference element to be built.
type DSSReference struct {
	// id is the Id attribute of the reference.
	id string

	// uri is the URI attribute of the reference.
	uri string

	// uriSet records whether setUri has been called, because Java's null and "" are two
	// DIFFERENT states here and both are reachable: the enveloped reference ReferenceBuilder
	// creates sets the URI to the empty string (XMLDSIG 4.4.3.2's URI="", the whole containing
	// document), while a detached reference over an unnamed document leaves it null. The two
	// are told apart by ReferenceProcessor, which writes a URI attribute for the first and none
	// for the second, and by DSSXMLUtils.isSameDocumentReference, which answers true for "" and
	// false for null. A bare "" cannot carry that distinction, so the flag does.
	//
	// The id and type fields need no equivalent: nothing in DSS ever assigns them the empty
	// string, so "" and null coincide there.
	uriSet bool

	// typ is the Type attribute of the reference.
	typ string

	// digestMethod is the DigestAlgorithm of the reference used to compute the digest value.
	// Java's field initializer is DigestAlgorithm.SHA512.
	digestMethod enumerations.DigestAlgorithm

	// contents is the referenced data.
	contents model.DSSDocument

	// transforms is the list of transforms to be performed.
	transforms []DSSTransform

	// object, when set, defines the ds:Object element to be incorporated within the signature.
	object *DSSObject
}

// NewDSSReference ports the default constructor, including the digestMethod field initializer.
func NewDSSReference() *DSSReference {
	return &DSSReference{digestMethod: enumerations.DigestAlgorithm_SHA512}
}

// Id gets the Id attribute of the reference. Ports getId().
func (r *DSSReference) Id() string { return r.id }

// SetId sets the Id attribute of the reference. Ports setId(String).
func (r *DSSReference) SetId(id string) { r.id = id }

// Uri gets the URI attribute of the reference. Ports getUri(); the empty string is returned
// both for Java's null and for Java's "", which HasUri tells apart.
func (r *DSSReference) Uri() string { return r.uri }

// HasUri reports whether a URI has been set at all, i.e. Java's getUri() != null. See the
// uriSet field for why the distinction is load-bearing.
func (r *DSSReference) HasUri() bool { return r.uriSet }

// SetUri sets the URI attribute of the reference. Ports setUri(String).
func (r *DSSReference) SetUri(uri string) {
	r.uri = uri
	r.uriSet = true
}

// Type gets the Type attribute of the reference. Ports getType().
func (r *DSSReference) Type() string { return r.typ }

// SetType sets the Type attribute of the reference. Ports setType(String).
func (r *DSSReference) SetType(typ string) { r.typ = typ }

// DigestMethodAlgorithm gets the DigestAlgorithm used for digest value computation.
// Ports getDigestMethodAlgorithm().
func (r *DSSReference) DigestMethodAlgorithm() enumerations.DigestAlgorithm { return r.digestMethod }

// SetDigestMethodAlgorithm sets the DigestAlgorithm to use for digest value computation.
// Ports setDigestMethodAlgorithm(DigestAlgorithm).
func (r *DSSReference) SetDigestMethodAlgorithm(digestMethod enumerations.DigestAlgorithm) {
	r.digestMethod = digestMethod
}

// Transforms gets the list of transforms to perform. Ports getTransforms().
func (r *DSSReference) Transforms() []DSSTransform { return r.transforms }

// SetTransforms sets the list of transforms to perform. Ports setTransforms(List).
func (r *DSSReference) SetTransforms(transforms []DSSTransform) { r.transforms = transforms }

// Contents gets the original referenced document content. Ports getContents().
func (r *DSSReference) Contents() model.DSSDocument { return r.contents }

// SetContents sets the original referenced document content. Ports setContents(DSSDocument).
func (r *DSSReference) SetContents(contents model.DSSDocument) { r.contents = contents }

// Object gets the ds:Object element's structure to be incorporated within the signature.
// Ports getObject().
func (r *DSSReference) Object() *DSSObject { return r.object }

// SetObject sets the pre-defined ds:Object element to be incorporated within the signature.
//
// NOTE: if not set, the basic ds:Object creation will be proceeded, when required.
//
// Ports setObject(DSSObject).
func (r *DSSReference) SetObject(object *DSSObject) { r.object = object }

// String ports toString(). An unset URI renders as Java's null does, i.e. uri='null'.
func (r *DSSReference) String() string {
	uri := r.uri
	if !r.uriSet {
		uri = "null"
	}
	return fmt.Sprintf("DSSReference{id='%s', uri='%s', type='%s', digestMethod=%s, contents=%v, transforms=%v, object=%v}",
		r.id, uri, r.typ, r.digestMethod, r.contents, r.transforms, r.object)
}

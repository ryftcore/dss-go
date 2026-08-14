// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XmlPolicyWithTransforms.java (DSS 6.5.RC1).
//
// Java extends model.Policy; the Go port embeds it, so every Policy getter/setter stays
// available and IsEmpty/Equals/String shadow the embedded ones exactly where Java overrides
// them. Java's hashCode() has no Go counterpart and is dropped (as elsewhere in this port);
// equals() becomes the Equals method the rest of dss-model already uses, taking the concrete
// type - Java's `getClass() != obj.getClass()` check is what a *XmlPolicyWithTransforms
// parameter expresses in Go.
//
// java.io.Serializable and the serialVersionUID are dropped (no Go counterpart).
package xades

import (
	"fmt"
	"reflect"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// XmlPolicyWithTransforms is an extension of the Policy class allowing addition of a custom
// list of DSSTransforms to build the ds:Transforms element.
//
// NOTE: The digest should be computed by the user and set through SetDigestValue.
//
// Use DSSXMLUtilsApplyTransforms(document, transforms) in order to obtain policy binaries
// after transforms.
type XmlPolicyWithTransforms struct {
	model.Policy

	// transforms is the list of transforms to be applied on the XML policy before the digest
	// calculation.
	transforms []DSSTransform
}

// NewXmlPolicyWithTransforms is the default constructor.
// Port of the XmlPolicyWithTransforms() constructor.
func NewXmlPolicyWithTransforms() *XmlPolicyWithTransforms {
	return &XmlPolicyWithTransforms{Policy: *model.NewPolicy()}
}

// Transforms gets the list of Transforms to incorporate into the signature.
// Port of #getTransforms.
func (p *XmlPolicyWithTransforms) Transforms() []DSSTransform {
	return p.transforms
}

// SetTransforms sets the list of Transforms to incorporate into the signature.
// Port of #setTransforms.
func (p *XmlPolicyWithTransforms) SetTransforms(transforms []DSSTransform) {
	p.transforms = transforms
}

// IsEmpty ports the overridden #isEmpty.
func (p *XmlPolicyWithTransforms) IsEmpty() bool {
	if !p.Policy.IsEmpty() {
		return false
	}
	if utils.IsCollectionNotEmpty(p.transforms) {
		return false
	}
	return true
}

// Equals ports the overridden #equals.
func (p *XmlPolicyWithTransforms) Equals(other *XmlPolicyWithTransforms) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.Policy.Equals(&other.Policy) {
		return false
	}
	return reflect.DeepEqual(p.transforms, other.transforms)
}

// String ports the overridden #toString.
func (p *XmlPolicyWithTransforms) String() string {
	return fmt.Sprintf("XmlPolicyWithTransforms {id='%s', qualifier=%v, description='%s', "+
		"documentationReferences=%v, digestAlgorithm=%v, digestValue=%v, spUri='%s', userNotice=%v, "+
		"spDocSpecification='%v', transforms=%v}",
		p.Id(), p.Qualifier(), p.Description(), p.DocumentationReferences(), p.DigestAlgorithm(),
		p.DigestValue(), p.Spuri(), p.UserNotice(), p.SpDocSpecification(), p.transforms)
}

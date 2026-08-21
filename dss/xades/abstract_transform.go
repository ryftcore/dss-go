// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/AbstractTransform.java (DSS 6.5.RC1).
//
// # The abstract base, and how its virtual members survive the port
//
// Java's AbstractTransform is abstract and its concrete subclasses override createTransform
// (XPathTransform, XPath2FilterTransform, XsltTransform) while ComplexTransform#buildTransformObject
// calls that method through Java's virtual dispatch. Go has no method overriding across
// embedding, so - exactly as xades_builder.go, xades_signature_builder.go and
// cades/cades_signature_extension.go already do, per the TokenBase.InitToken(self) convention
// of PORTING.md - every concrete transform hands itself to InitAbstractTransform, and the base
// reaches the override through t.self.
//
// The same self also stands in for Java's getClass() in equals(): the class-identity check that
// makes a Base64Transform unequal to any other transform carrying the same algorithm URI is
// reflect.TypeOf(t.self), the pattern spi/revocation_token.go established.
//
// hashCode() has no Go counterpart; upstream needs it only to key JDK hash collections.
// java.io.Serializable and serialVersionUID are dropped.
package xades

import (
	"reflect"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// AbstractTransform is the abstract implementation of a transform.
type AbstractTransform struct {
	// algorithm is the algorithm url string.
	algorithm string

	// namespace is the namespace. Java's field initializer is XMLDSigNamespace.NS.
	namespace *common.DSSNamespace

	// self is the concrete transform this base is embedded in; see the file header.
	self DSSTransform
}

// newAbstractTransform ports the protected AbstractTransform(String algorithm), whose field
// initializer leaves the namespace at XMLDSigNamespace.NS.
func newAbstractTransform(algorithm string) AbstractTransform {
	return AbstractTransform{algorithm: algorithm, namespace: common.XMLDSigNS}
}

// newAbstractTransformWithNamespace ports the protected
// AbstractTransform(DSSNamespace, String).
func newAbstractTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace, algorithm string) AbstractTransform {
	return AbstractTransform{algorithm: algorithm, namespace: xmlDSigNamespace}
}

// InitAbstractTransform registers the concrete transform embedding this base. Every constructor
// in the hierarchy calls it; see the file header.
func (t *AbstractTransform) InitAbstractTransform(self DSSTransform) {
	t.self = self
}

// abstractTransform returns the embedded base. Promoted through embedding, it is how a
// subclass's Equals reaches the fields Java's AbstractTransform#equals compares directly,
// mirroring spi.RevocationTokenBase#revocationTokenBase.
func (t *AbstractTransform) abstractTransform() *AbstractTransform { return t }

// Algorithm returns the transformation algorithm url. Ports getAlgorithm().
func (t *AbstractTransform) Algorithm() string { return t.algorithm }

// SetNamespace specifies a namespace for the transformation elements.
// Ports setNamespace(DSSNamespace).
func (t *AbstractTransform) SetNamespace(namespace *common.DSSNamespace) {
	t.namespace = namespace
}

// CreateTransform creates a ds:Transform element and appends it to parentNode.
// Ports createTransform(Document, Element).
func (t *AbstractTransform) CreateTransform(document, parentNode *xmldom.Node) *xmldom.Node {
	transformDom := xmlutils.DomUtilsAddElement(document, parentNode, t.namespace, common.XMLDSigElement_TRANSFORM)
	transformDom.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ALGORITHM.AttributeName()}, t.algorithm)
	return transformDom
}

// Equals ports equals(Object): the algorithm url, the namespace and the concrete class must
// all match. See the file header for how getClass() maps onto reflect.TypeOf(self).
func (t *AbstractTransform) Equals(obj DSSTransform) bool {
	if obj == nil {
		return false
	}
	other := abstractTransformOf(obj)
	if other == nil {
		return false
	}
	if t == other {
		return true
	}
	if reflect.TypeOf(t.self) != reflect.TypeOf(other.self) {
		return false
	}
	if t.algorithm != other.algorithm {
		return false
	}
	return abstractTransformNamespaceEquals(t.namespace, other.namespace)
}

// String ports toString().
func (t *AbstractTransform) String() string {
	return "DSSTransform [algorithm=" + t.algorithm + ", namespace=" + abstractTransformNamespaceString(t.namespace) + "]"
}

// abstractTransformHolder is the accessor abstractTransformOf looks for; it is promoted to
// every type embedding AbstractTransform.
type abstractTransformHolder interface {
	abstractTransform() *AbstractTransform
}

// abstractTransformOf extracts the embedded AbstractTransform of any transform.
func abstractTransformOf(transform DSSTransform) *AbstractTransform {
	if holder, ok := transform.(abstractTransformHolder); ok {
		return holder.abstractTransform()
	}
	return nil
}

// abstractTransformNamespaceEquals ports Objects.equals(namespace, other.namespace) over
// DSSNamespace, whose Java equals compares the uri and the prefix.
func abstractTransformNamespaceEquals(a, b *common.DSSNamespace) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Uri() == b.Uri() && a.Prefix() == b.Prefix()
}

// abstractTransformNamespaceString renders a DSSNamespace the way Java string concatenation
// does, i.e. "null" for a null reference.
func abstractTransformNamespaceString(namespace *common.DSSNamespace) string {
	if namespace == nil {
		return "null"
	}
	return namespace.String()
}

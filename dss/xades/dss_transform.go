// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/DSSTransform.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/xml/common"
)

// DSSTransform defines a transform used for a reference transformation.
//
// Java's interface extends java.io.Serializable, which has no Go counterpart and is dropped.
//
// PerformTransform returns an error where Java throws the unchecked DSSException /
// IllegalArgumentException that ComplexTransform and SPDocDigestAsInSpecificationTransform
// raise (PORTING.md: throw -> (T, error)). CreateTransform cannot fail in either language and
// keeps its bare return.
type DSSTransform interface {
	// Algorithm returns a particular transformation algorithm name. Ports getAlgorithm().
	Algorithm() string

	// SetNamespace specifies a namespace for the transformation elements.
	// Ports setNamespace(DSSNamespace).
	SetNamespace(namespace *common.DSSNamespace)

	// PerformTransform executes a transform on the provided DSSTransformOutput.
	// Ports performTransform(DSSTransformOutput).
	PerformTransform(transformOutput *DSSTransformOutput) (*DSSTransformOutput, error)

	// CreateTransform creates a Transform element DOM and appends it to parentNode.
	// Ports createTransform(Document, Element).
	CreateTransform(document, parentNode *xmldom.Node) *xmldom.Node
}

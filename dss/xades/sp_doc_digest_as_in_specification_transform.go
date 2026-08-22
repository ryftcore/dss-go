// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/SPDocDigestAsInSpecificationTransform.java (DSS 6.5.RC1).
package xades

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/xml/common"
)

// errSPDocDigestAsInSpecificationTransformNotForReferences is the IllegalArgumentException
// performTransform throws, verbatim.
var errSPDocDigestAsInSpecificationTransformNotForReferences = errors.New(
	"The transform SPDocDigestAsInSpecificationTransform cannot be used for reference processing!")

// spDocDigestAsInSpecificationTransformAlgorithmURI is the SPDocDigestAsInSpecification
// algorithm URI, "http://uri.etsi.org/01903/v1.3.2/Policy/SPDocDigestAsInSpecification".
// Port of the private ALGORITHM_URI field, which reads
// DSSXMLUtils.SP_DOC_DIGEST_AS_IN_SPECIFICATION_ALGORITHM_URI - so the value is taken from
// there rather than repeated. A var rather than a const only so that the declaration does not
// depend on how dss_xml_utils.go chose to declare its own.
var spDocDigestAsInSpecificationTransformAlgorithmURI = DSSXMLUtilsSPDocDigestAsInSpecificationAlgorithmURI

// SPDocDigestAsInSpecificationTransform is a special transform to be used exclusively within a
// xades:SignaturePolicyId to define special digest computation rules.
// See EN 319 132-1 "5.2.9 The SignaturePolicyIdentifier qualifying property".
type SPDocDigestAsInSpecificationTransform struct {
	AbstractTransform
}

// NewSPDocDigestAsInSpecificationTransform ports the default constructor with the ds: xmldsig
// namespace.
func NewSPDocDigestAsInSpecificationTransform() *SPDocDigestAsInSpecificationTransform {
	t := &SPDocDigestAsInSpecificationTransform{
		AbstractTransform: newAbstractTransform(spDocDigestAsInSpecificationTransformAlgorithmURI),
	}
	t.InitAbstractTransform(t)
	return t
}

// newSPDocDigestAsInSpecificationTransformWithNamespace ports the protected
// SPDocDigestAsInSpecificationTransform(DSSNamespace); it is unexported because Java declares
// it protected and no subclass exists.
func newSPDocDigestAsInSpecificationTransformWithNamespace(
	xmlDSigNamespace *common.DSSNamespace) *SPDocDigestAsInSpecificationTransform {
	t := &SPDocDigestAsInSpecificationTransform{
		AbstractTransform: newAbstractTransformWithNamespace(
			xmlDSigNamespace, spDocDigestAsInSpecificationTransformAlgorithmURI),
	}
	t.InitAbstractTransform(t)
	return t
}

// PerformTransform always fails: this transform carries a digest-computation rule, it is not a
// reference-processing step. Ports performTransform(DSSTransformOutput), whose
// IllegalArgumentException becomes the returned error with the same message.
func (t *SPDocDigestAsInSpecificationTransform) PerformTransform(
	transformOutput *DSSTransformOutput) (*DSSTransformOutput, error) {
	return nil, errSPDocDigestAsInSpecificationTransformNotForReferences
}

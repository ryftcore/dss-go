// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/CanonicalizationTransform.java (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// CanonicalizationTransform performs a canonicalization transform on an XML NodeSet.
type CanonicalizationTransform struct {
	ComplexTransform
}

// NewCanonicalizationTransform ports CanonicalizationTransform(String), which delegates to the
// namespace-taking constructor with XMLDSigNamespace.NS.
//
// Java's constructor throws UnsupportedOperationException for a canonicalization method the
// framework cannot canonicalize with; that becomes a panic here rather than an (T, error)
// return, because every call site in the port - ReferenceBuilder, CounterSignatureBuilder,
// XAdESSignatureBuilder - passes a constant that XMLCanonicalizer accepts, and Java's
// UnsupportedOperationException is likewise an unchecked programming error rather than an
// input-validation failure.
func NewCanonicalizationTransform(canonicalizationAlgorithm string) *CanonicalizationTransform {
	return NewCanonicalizationTransformWithNamespace(common.XMLDSigNS, canonicalizationAlgorithm)
}

// NewCanonicalizationTransformWithNamespace ports
// CanonicalizationTransform(DSSNamespace, String). See NewCanonicalizationTransform for how the
// UnsupportedOperationException maps.
func NewCanonicalizationTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace,
	canonicalizationAlgorithm string) *CanonicalizationTransform {
	t := &CanonicalizationTransform{
		ComplexTransform: newComplexTransform(xmlDSigNamespace, canonicalizationAlgorithm),
	}
	t.InitAbstractTransform(t)
	if !xmlutils.XMLCanonicalizerCanCanonicalize(canonicalizationAlgorithm) {
		panic(fmt.Sprintf("The provided canonicalization method [%s] is not supported!",
			canonicalizationAlgorithm))
	}
	return t
}

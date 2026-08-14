// Ported from the canonicalization registry of
// dss-xml-utils/.../XMLCanonicalizer.java (DSS 6.5.RC1), plus
// org.apache.xml.security.transforms.params.InclusiveNamespaces#prefixStr2Set and
// org.apache.xml.security.c14n.helper.C14nHelper (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Algorithm is a canonicalization method URI.
type Algorithm string

// The seven algorithms XMLCanonicalizer.registerDefaultCanonicalizers() registers. Upstream
// refuses anything else with IllegalArgumentException; Resolve returns ErrUnsupportedAlgorithm.
const (
	C14N10                    Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	C14N10WithComments        Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"
	C14N11                    Algorithm = "http://www.w3.org/2006/12/xml-c14n11"
	C14N11WithComments        Algorithm = "http://www.w3.org/2006/12/xml-c14n11#WithComments"
	C14NExclusive             Algorithm = "http://www.w3.org/2001/10/xml-exc-c14n#"
	C14NExclusiveWithComments Algorithm = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"
	C14NPhysical              Algorithm = "http://santuario.apache.org/c14n/physical"

	// DefaultDSS is XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD: what DSS signs with.
	DefaultDSS = C14NExclusive
	// DefaultXMLDSig is XMLCanonicalizer.DEFAULT_XMLDSIG_C14N_METHOD: what XMLDSIG 4.4.3.2
	// prescribes when a signature omits the method (DSS-2208).
	DefaultXMLDSig = C14N10
)

// ErrUnsupportedAlgorithm is returned by Resolve and Canonicalize.
var ErrUnsupportedAlgorithm = errors.New("xmlc14n: unsupported canonicalization algorithm")

var registered = map[Algorithm]struct{}{
	C14N10:                    {},
	C14N10WithComments:        {},
	C14N11:                    {},
	C14N11WithComments:        {},
	C14NExclusive:             {},
	C14NExclusiveWithComments: {},
	C14NPhysical:              {},
}

// Supported reports whether alg is one of the seven registered algorithms.
func Supported(alg Algorithm) bool {
	_, ok := registered[alg]
	return ok
}

// Resolve maps "" to DefaultXMLDSig and validates alg.
func Resolve(alg Algorithm) (Algorithm, error) {
	if alg == "" {
		alg = DefaultXMLDSig
	}
	if !Supported(alg) {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedAlgorithm, string(alg))
	}
	return alg, nil
}

// withComments reports whether alg keeps comment nodes.
func (a Algorithm) withComments() bool {
	return a == C14N10WithComments || a == C14N11WithComments || a == C14NExclusiveWithComments
}

// exclusive reports whether alg is Exclusive XML Canonicalization.
func (a Algorithm) exclusive() bool {
	return a == C14NExclusive || a == C14NExclusiveWithComments
}

// c14n11 reports whether alg is Canonical XML 1.1.
func (a Algorithm) c14n11() bool {
	return a == C14N11 || a == C14N11WithComments
}

// physical reports whether alg is the Santuario "physical" method.
func (a Algorithm) physical() bool { return a == C14NPhysical }

// ParsePrefixList splits an InclusiveNamespaces PrefixList attribute value on whitespace and
// maps "#default" to "xmlns".
//
// It is a faithful port of InclusiveNamespaces.prefixStr2Set, which does
// inclusiveNamespaces.split("\\s") into a TreeSet: the result is sorted and deduplicated, a
// nil or empty input yields an empty list, and - because Java's split on a single-character
// pattern keeps interior empty fields while dropping trailing ones - a run of two spaces
// contributes an empty token. That empty token names no prefix and so renders nothing, which
// is why upstream never noticed; the excl-prefixlist "p q  zz" known-answer test pins it.
// Java's \s is [ \t\n\x0B\f\r], which is not Go's regexp \s (no vertical tab), so the split is
// written out by hand.
func ParsePrefixList(s string) []string {
	if s == "" {
		return nil
	}
	isJavaSpace := func(r byte) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\v' || r == '\f' || r == '\r'
	}
	var tokens []string
	start := 0
	for i := 0; i < len(s); i++ {
		if isJavaSpace(s[i]) {
			tokens = append(tokens, s[start:i])
			start = i + 1
		}
	}
	tokens = append(tokens, s[start:])
	// String.split(regex) with the default limit drops trailing empty fields.
	for len(tokens) > 0 && tokens[len(tokens)-1] == "" {
		tokens = tokens[:len(tokens)-1]
	}

	set := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		if t == "#default" {
			t = xmlnsPrefix
		}
		set[t] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// RelativeNamespaceError reports a rendered namespace declaration whose URI is relative.
// C14N 1.0, 1.1 and exclusive reject these; physical accepts them.
type RelativeNamespaceError struct {
	Element string // the element's QName
	Prefix  string // "xmlns" for the default declaration
	URI     string
}

func (e *RelativeNamespaceError) Error() string {
	return fmt.Sprintf("xmlc14n: element %s has a relative namespace: %s=%q", e.Element, e.Prefix, e.URI)
}

// namespaceIsRelative ports C14nHelper.namespaceIsAbsolute, negated: an empty value is treated
// as absolute, and anything else is absolute exactly when it contains a ':' other than at
// index 0. Note that this is not a URI parse and must not be replaced by one.
func namespaceIsRelative(value string) bool {
	if value == "" {
		return false
	}
	return strings.IndexByte(value, ':') <= 0
}

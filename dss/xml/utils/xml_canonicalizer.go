// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/XMLCanonicalizer.java (DSS 6.5.RC1).
//
// Upstream wraps a single org.apache.xml.security.c14n.Canonicalizer instance obtained once
// at construction time (Canonicalizer.getInstance(method)) and reused across every
// canonicalize(...) overload. internal/xmlc14n is deliberately stateless per call instead
// (XML_DESIGN.md's "one-canonicalization-per-call" decision, working around SANTUARIO-463),
// so XMLCanonicalizer here stores only the resolved algorithm and calls internal/xmlc14n
// fresh on every method - functionally equivalent for upstream's own call pattern, since
// every DSS call site already does XMLCanonicalizer.createInstance(m).canonicalize(...)
// inline rather than reusing an instance across canonicalizations (verified by XML_DESIGN.md
// against a grep over dss-xades).
//
// Java throw-in-constructor becomes a (T, error)-returning constructor func per PORTING.md:
// XMLCanonicalizerCreateInstance(WithMethod) returns an error instead of throwing
// IllegalArgumentException/DSSException.
package utils

import (
	"fmt"
	"io"
	"sync"

	"github.com/utain/esig/dss/internal/xmlc14n"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// XMLCanonicalizerDefaultDSSC14NMethod is the default canonicalization method used for
// production of signatures within the DSS framework. Ports DEFAULT_DSS_C14N_METHOD
// (= javax.xml.crypto.dsig.CanonicalizationMethod.EXCLUSIVE).
const XMLCanonicalizerDefaultDSSC14NMethod = string(xmlc14n.C14NExclusive)

// XMLCanonicalizerDefaultXMLDSigC14NMethod is the default canonicalization method for XMLDSIG
// used for signatures and timestamps (see XMLDSIG 4.4.3.2) when one is not defined. Ports
// DEFAULT_XMLDSIG_C14N_METHOD (= CanonicalizationMethod.INCLUSIVE).
const XMLCanonicalizerDefaultXMLDSigC14NMethod = string(xmlc14n.C14N10)

var (
	xmlCanonicalizerMu sync.RWMutex

	// xmlCanonicalizerRegistered mirrors the "canonicalizers" HashSet<String>, seeded by
	// registerDefaultCanonicalizers() with the same seven algorithm URIs internal/xmlc14n
	// implements.
	xmlCanonicalizerRegistered = map[string]struct{}{
		string(xmlc14n.C14N10):                    {},
		string(xmlc14n.C14N10WithComments):        {},
		string(xmlc14n.C14N11):                    {},
		string(xmlc14n.C14N11WithComments):        {},
		string(xmlc14n.C14NExclusive):             {},
		string(xmlc14n.C14NExclusiveWithComments): {},
		string(xmlc14n.C14NPhysical):              {},
	}
)

// XMLCanonicalizer contains a set of methods for canonicalization of *xmldom.Node.
type XMLCanonicalizer struct {
	method xmlc14n.Algorithm
}

// XMLCanonicalizerCreateInstance creates an XMLCanonicalizer to be used with the default
// XMLDSig canonicalization method "http://www.w3.org/TR/2001/REC-xml-c14n-20010315". Ports
// createInstance().
func XMLCanonicalizerCreateInstance() (*XMLCanonicalizer, error) {
	return XMLCanonicalizerCreateInstanceWithMethod("")
}

// XMLCanonicalizerCreateInstanceWithMethod creates an XMLCanonicalizer with the provided
// canonicalization method. If canonicalizationMethod is "", the default XMLDSig
// canonicalization method is used. Ports createInstance(String) and the private constructor.
func XMLCanonicalizerCreateInstanceWithMethod(canonicalizationMethod string) (*XMLCanonicalizer, error) {
	canonicalizationMethod = xmlCanonicalizerGetCanonicalizationMethod(canonicalizationMethod)
	if !XMLCanonicalizerCanCanonicalize(canonicalizationMethod) {
		return nil, fmt.Errorf("the canonicalization method '%s' is not supported! "+
			"Use XMLCanonicalizerRegisterCanonicalizer to add support of a canonicalization method", canonicalizationMethod)
	}
	alg := xmlc14n.Algorithm(canonicalizationMethod)
	if !xmlc14n.Supported(alg) {
		// Mirrors initCanonicalizer's InvalidCanonicalizerException: the registry (settable
		// via RegisterCanonicalizer) accepted the URI, but the underlying c14n engine does
		// not actually implement it.
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("The canonicalizer cannot be instantiated with canonicalization method '%s'!", canonicalizationMethod),
			xmlc14n.ErrUnsupportedAlgorithm)
	}
	return &XMLCanonicalizer{method: alg}, nil
}

// XMLCanonicalizerCanCanonicalize says if the framework can canonicalize an XML data with the
// provided method. Ports canCanonicalize(String).
func XMLCanonicalizerCanCanonicalize(canonicalizationMethod string) bool {
	xmlCanonicalizerMu.RLock()
	defer xmlCanonicalizerMu.RUnlock()
	_, ok := xmlCanonicalizerRegistered[canonicalizationMethod]
	return ok
}

// XMLCanonicalizerRegisterCanonicalizer registers a canonicalizer URI and reports whether it
// was not already registered. Ports registerCanonicalizer(String).
func XMLCanonicalizerRegisterCanonicalizer(c14nAlgorithmURI string) bool {
	xmlCanonicalizerMu.Lock()
	defer xmlCanonicalizerMu.Unlock()
	if _, exists := xmlCanonicalizerRegistered[c14nAlgorithmURI]; exists {
		return false
	}
	xmlCanonicalizerRegistered[c14nAlgorithmURI] = struct{}{}
	return true
}

// xmlCanonicalizerGetCanonicalizationMethod returns canonicalizationMethod if provided,
// otherwise XMLCanonicalizerDefaultXMLDSigC14NMethod. Ports the private
// getCanonicalizationMethod(String) (DSS-2208); LOG.warn dropped per PORTING.md.
func xmlCanonicalizerGetCanonicalizationMethod(canonicalizationMethod string) string {
	if utils.IsStringEmpty(canonicalizationMethod) {
		return XMLCanonicalizerDefaultXMLDSigC14NMethod
	}
	return canonicalizationMethod
}

// CanonicalizeStream canonicalizes the given io.Reader using the configured
// canonicalization method. NOTE: closes r after reading. Ports canonicalize(InputStream);
// this goes through DomUtilsBuildDOMFromReader then the Node path, exactly like upstream's
// canonicalize(DomUtils.buildDOM(is)) - a different code path from CanonicalizeBytes, which
// lets the c14n engine parse the bytes itself.
func (c *XMLCanonicalizer) CanonicalizeStream(r io.Reader) ([]byte, error) {
	doc, err := DomUtilsBuildDOMFromReader(r)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Error on InputStream processing", err)
	}
	return c.CanonicalizeNode(doc)
}

// CanonicalizeStreamTo canonicalizes r using the configured canonicalization method and
// writes the result into w. NOTE: closes r after reading. Ports
// canonicalize(InputStream, OutputStream).
func (c *XMLCanonicalizer) CanonicalizeStreamTo(r io.Reader, w io.Writer) error {
	doc, err := DomUtilsBuildDOMFromReader(r)
	if err != nil {
		return model.NewDSSErrorMessageCause("Cannot canonicalize the InputStream", err)
	}
	return c.CanonicalizeNodeTo(doc, w)
}

// CanonicalizeBytes canonicalizes toCanonicalizeBytes using the configured canonicalization
// method. The bytes are parsed by the c14n engine itself, not by DomUtils. Ports
// canonicalize(byte[]).
func (c *XMLCanonicalizer) CanonicalizeBytes(toCanonicalizeBytes []byte) ([]byte, error) {
	out, err := xmlc14n.CanonicalizeBytes(c.method, toCanonicalizeBytes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot canonicalize the binaries", err)
	}
	return out, nil
}

// CanonicalizeBytesTo canonicalizes toCanonicalizeBytes and writes the result into w. Ports
// canonicalize(byte[], OutputStream).
func (c *XMLCanonicalizer) CanonicalizeBytesTo(toCanonicalizeBytes []byte, w io.Writer) error {
	out, err := c.CanonicalizeBytes(toCanonicalizeBytes)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

// CanonicalizeNode canonicalizes the given Node using the configured canonicalization
// method. Ports canonicalize(Node).
func (c *XMLCanonicalizer) CanonicalizeNode(node *xmldom.Node) ([]byte, error) {
	out, err := xmlc14n.CanonicalizeNode(c.method, node)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot canonicalize the subtree", err)
	}
	return out, nil
}

// CanonicalizeNodeTo canonicalizes node and writes the result into w. Ports
// canonicalize(Node, OutputStream).
func (c *XMLCanonicalizer) CanonicalizeNodeTo(node *xmldom.Node, w io.Writer) error {
	if err := xmlc14n.Canonicalize(c.method, xmlc14n.Input{Node: node}, w); err != nil {
		return model.NewDSSErrorMessageCause("Cannot canonicalize the subtree", err)
	}
	return nil
}

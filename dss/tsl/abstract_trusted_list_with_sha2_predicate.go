// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sha2/AbstractTrustedListWithSha2Predicate.java (DSS 6.5.RC1).
package tsl

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	xadesdefinition "github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// TrustedListWithSha2Predicate is the Go form of Java's
// Predicate<DocumentWithSha2>: it decides whether a cached Trusted List document still matches
// its published .sha2 digest and therefore needs no refresh.
//
// DEVIATION: Test returns (bool, error) rather than the bare boolean of java.util.function.
// Predicate#test. Java's implementation reaches DSSDocument#getDigestValue and
// DSSUtils#toByteArray, both of which throw DSSException on an I/O failure and propagate that
// exception out of test() into Sha2FileCacheDataLoader#getDocument; the Go ports of those two
// return errors instead, so the predicate has to carry them.
type TrustedListWithSha2Predicate interface {
	Test(documentWithSha2 *DocumentWithSha2) (bool, error)
}

// sha2Algorithm is the default sha2 digest algorithm defined in ETSI TS 119 612. Port of the
// protected constant SHA2_ALGORITHM.
const sha2Algorithm = enumerations.DigestAlgorithm_SHA256

// init ports the class's static initialiser, which registers the Trusted List namespace with
// the XPath resolver so that TrustedListPath queries can be evaluated.
func init() {
	xmlutils.XPathUtilsRegisterNamespace(xadesdefinition.TrustedListNamespace_NS)
}

// AbstractTrustedListWithSha2PredicateBase carries the concrete behaviour of the Java abstract
// class AbstractTrustedListWithSha2Predicate: the common utility methods a Trusted List sha2
// validation needs. Concrete predicates (DefaultTrustedListWithSha2Predicate) embed it and
// supply Test themselves.
//
// The type holds no state, so it needs neither a constructor nor an Init/overrides registration:
// nothing in the base dispatches back into the subclass.
type AbstractTrustedListWithSha2PredicateBase struct{}

// OriginalDocumentDigest computes the Digest on the given document's content. Port of the
// protected getOriginalDocumentDigest(DSSDocument); a nil document yields the zero Digest,
// standing in for Java's null return.
func (p *AbstractTrustedListWithSha2PredicateBase) OriginalDocumentDigest(document model.DSSDocument) (model.Digest, error) {
	if document == nil {
		return model.Digest{}, nil
	}
	digest, err := document.DigestValue(sha2Algorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(sha2Algorithm, digest), nil
}

// Sha2Digest parses the sha2Document and returns the Digest it carries. Port of the protected
// getSha2Digest(DSSDocument); a nil document yields the zero Digest, standing in for Java's
// null return.
//
// NOTE: HEX encoding is not explicitly defined in the standard, but all known implementations
// publish a HEX-encoded digest, so a HEX-looking payload is decoded and anything else is taken
// as the raw digest bytes - upstream's behaviour, verbatim.
func (p *AbstractTrustedListWithSha2PredicateBase) Sha2Digest(sha2Document model.DSSDocument) (model.Digest, error) {
	if sha2Document == nil {
		return model.Digest{}, nil
	}
	sha2DocumentBinaries, err := spi.DSSUtilsToByteArrayOfDocument(sha2Document)
	if err != nil {
		return model.Digest{}, err
	}
	sha2DocumentStr := string(sha2DocumentBinaries)
	if utils.IsHexEncoded(sha2DocumentStr) {
		decoded, err := utils.FromHex(sha2DocumentStr)
		if err != nil {
			return model.Digest{}, err
		}
		sha2DocumentBinaries = decoded
	}
	return model.NewDigest(sha2Algorithm, sha2DocumentBinaries), nil
}

// NextUpdate retrieves the NextUpdate date value from the provided Trusted List document. Port
// of the protected getNextUpdate(DSSDocument); the zero time.Time stands in for Java's null.
//
// Java catches every exception raised while building the DOM or evaluating the XPath and logs a
// warning before answering null; the logging is dropped and the zero time answered instead.
func (p *AbstractTrustedListWithSha2PredicateBase) NextUpdate(tlDocument model.DSSDocument) time.Time {
	// Upstream logs "The document is not XML! Unable to extract NextUpdate." and then carries
	// on regardless - the guard is a warning, not a short-circuit.
	_ = xmlutils.DomUtilsIsDOM(tlDocument)

	documentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(tlDocument)
	if err != nil {
		return time.Time{}
	}
	documentElement := documentDom.DocumentElement()
	if documentElement == nil {
		return time.Time{}
	}
	nextUpdateElement, err := xmlutils.XPathUtilsGetElement(documentElement,
		xadesdefinition.TrustedListPath_NEXT_UPDATE_PATH)
	if err != nil {
		return time.Time{}
	}
	if nextUpdateElement != nil {
		nextUpdate := nextUpdateElement.TextContent()
		if utils.IsStringNotEmpty(nextUpdate) {
			nextUpdate = utils.Trim(nextUpdate)
			return spi.DSSUtilsParseRFCDate(nextUpdate)
		}
		// Upstream logs "NextUpdate element has an empty content." at debug level.
	}
	// Upstream logs "No NextUpdate element found!".
	return time.Time{}
}

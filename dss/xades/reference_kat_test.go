package xades

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// The KAT for the XAdES reference-and-transform core: the DSSTransform hierarchy,
// DSSTransformOutput, DSSReference, ReferenceIdProvider, ReferenceBuilder, ReferenceProcessor
// and ReferenceVerifier.
//
// Every expectation in testdata/refs.txt is upstream DSS 6.5.RC1's own output, dumped by
// testdata/gen/RefsOracle.java from the same fixed inputs this file feeds the Go port. Nothing
// is hand-derived. The keys are:
//
//	transforms   the ds:Transforms subtree a transform WRITES into a ds:Reference - element
//	             namespaces and prefixes, attribute order, and the declarations the Filter 2.0
//	             transform puts on its XPath element
//	algorithm    the transform's algorithm URI
//	output       the octets a transform chain EXECUTES to, i.e. what a reference digest covers
//	error        the message of the exception a rejected setup throws ("" when none)
//	ref<i>-*     the fields ReferenceBuilder fills into each DSSReference
//	references   the complete ds:Reference list ReferenceProcessor incorporates, digests included
//	signedinfo-raw / signedinfo-c14n / document   the end-to-end ds:SignedInfo and signature
//
// Two cases record a deliberate divergence rather than a match; both are called out where they
// are asserted, and both belong to the frozen internal/xmldsig package.

// xadesRefsOracle is the parsed testdata/refs.txt: case name -> key -> bytes.
type xadesRefsOracle map[string]map[string][]byte

func loadXAdESRefsOracle(t *testing.T) xadesRefsOracle {
	t.Helper()
	raw, err := os.ReadFile(corpustest.Path(t, "refs.txt"))
	if err != nil {
		t.Fatalf("reading the oracle: %v", err)
	}
	oracle := xadesRefsOracle{}
	current := ""
	for _, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "case ") {
			current = strings.TrimPrefix(line, "case ")
			if _, ok := oracle[current]; !ok {
				oracle[current] = map[string][]byte{}
			}
			continue
		}
		key, value, ok := strings.Cut(line, " ")
		if !ok {
			t.Fatalf("malformed oracle line: %q", line)
		}
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			t.Fatalf("case %s key %s: %v", current, key, err)
		}
		oracle[current][key] = decoded
	}
	if len(oracle) == 0 {
		t.Fatal("the oracle is empty")
	}
	return oracle
}

func (o xadesRefsOracle) expect(t *testing.T, caseName, key string) []byte {
	t.Helper()
	record, ok := o[caseName]
	if !ok {
		t.Fatalf("the oracle has no case %q", caseName)
	}
	value, ok := record[key]
	if !ok {
		t.Fatalf("the oracle case %q has no key %q", caseName, key)
	}
	return value
}

func (o xadesRefsOracle) expectString(t *testing.T, caseName, key string) string {
	t.Helper()
	return string(o.expect(t, caseName, key))
}

// xadesRefsNull is the marker RefsOracle.java writes for a Java null string, so that it cannot
// be confused with an empty one - the distinction DSSReference.HasUri exists for.
const xadesRefsNull = "\x00NULL"

func xadesRefsAssertBytes(t *testing.T, caseName, key string, want, got []byte) {
	t.Helper()
	if !bytes.Equal(want, got) {
		t.Errorf("case %s / %s: bytes differ from the Java oracle\nwant: %s\ngot:  %s",
			caseName, key, want, got)
	}
}

func xadesRefsAssertString(t *testing.T, caseName, key, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("case %s / %s: %q, want %q (Java oracle)", caseName, key, got, want)
	}
}

// xadesRefsErrorMessage strips the exception class name RefsOracle.java prefixes its recorded
// messages with; Go errors carry the message alone.
func xadesRefsErrorMessage(recorded string) string {
	_, message, ok := strings.Cut(recorded, ": ")
	if !ok {
		return recorded
	}
	return message
}

// xadesRefsAssertError compares a Go error against the oracle's recorded exception: an empty
// record means upstream did not throw.
func xadesRefsAssertError(t *testing.T, oracle xadesRefsOracle, caseName string, err error) {
	t.Helper()
	recorded := oracle.expectString(t, caseName, "error")
	if recorded == "" {
		if err != nil {
			t.Errorf("case %s: unexpected error %v, upstream succeeds", caseName, err)
		}
		return
	}
	if err == nil {
		t.Errorf("case %s: no error, upstream throws %q", caseName, recorded)
		return
	}
	if want := xadesRefsErrorMessage(recorded); err.Error() != want {
		t.Errorf("case %s: error %q, want %q (Java oracle)", caseName, err.Error(), want)
	}
}

// -------------------------------------------------------------------------------- fixtures

var (
	xadesRefsDSNamespace      = common.XMLDSigNS
	xadesRefsDefaultNamespace = common.NewDSSNamespace(common.XMLDSigNS.Uri(), "")
	xadesRefsDsigNamespace    = common.NewDSSNamespace(common.XMLDSigNS.Uri(), "dsig")
)

const (
	xadesRefsC14NExclusive         = "http://www.w3.org/2001/10/xml-exc-c14n#"
	xadesRefsC14NExclusiveComments = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"
	xadesRefsC14NInclusive         = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	xadesRefsC14NInclusiveComments = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"
	xadesRefsC14N11                = "http://www.w3.org/2006/12/xml-c14n11"
	xadesRefsC14N11Comments        = "http://www.w3.org/2006/12/xml-c14n11#WithComments"
	xadesRefsC14NPhysical          = "http://santuario.apache.org/c14n/physical"
)

const (
	xadesRefsTextContent = "Hello World!"

	xadesRefsXMLContent = `<?xml version="1.0" encoding="UTF-8"?><root xmlns="http://sample.com" Id="root-id">` +
		`<child Id="child-id">text</child></root>`

	xadesRefsXMLWithoutIDContent = `<?xml version="1.0" encoding="UTF-8"?><root xmlns="http://sample.com">` +
		`<child>text</child></root>`

	xadesRefsXMLSignedContent = `<?xml version="1.0" encoding="UTF-8"?><root xmlns="http://sample.com" Id="root-id">` +
		`<child Id="child-id">text</child><!-- a comment -->` +
		`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="sig-1">` +
		`<ds:SignedInfo><ds:Reference URI=""/></ds:SignedInfo></ds:Signature></root>`

	xadesRefsXSLTContent = `<?xml version="1.0" encoding="UTF-8"?>` +
		`<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" version="1.0">` +
		`<xsl:template match="/"><out/></xsl:template></xsl:stylesheet>`
)

func xadesRefsTextDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesRefsTextContent), "hello.txt",
		enumerations.MimeTypeEnumText)
}

func xadesRefsSecondDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte("second"), "second.txt", enumerations.MimeTypeEnumText)
}

func xadesRefsUnnamedDocument() model.DSSDocument {
	return model.NewInMemoryDocument([]byte(xadesRefsTextContent))
}

func xadesRefsXMLDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesRefsXMLContent), "sample.xml",
		enumerations.MimeTypeEnumXML)
}

func xadesRefsSignedXMLDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesRefsXMLSignedContent), "signed.xml",
		enumerations.MimeTypeEnumXML)
}

func xadesRefsXMLWithoutIDDocument() model.DSSDocument {
	return model.NewInMemoryDocumentWithMimeType([]byte(xadesRefsXMLWithoutIDContent), "no-id.xml",
		enumerations.MimeTypeEnumXML)
}

func xadesRefsManifestDocument(t *testing.T) model.DSSDocument {
	t.Helper()
	builder, err := NewManifestBuilder(enumerations.DigestAlgorithmSHA256,
		[]model.DSSDocument{xadesRefsTextDocument()})
	if err != nil {
		t.Fatalf("NewManifestBuilder: %v", err)
	}
	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("ManifestBuilder.Build: %v", err)
	}
	return manifest
}

func xadesRefsDigestDocument(t *testing.T) model.DSSDocument {
	t.Helper()
	digest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, []byte(xadesRefsTextContent))
	if err != nil {
		t.Fatalf("digesting: %v", err)
	}
	return model.NewDigestDocumentFromValueWithName(enumerations.DigestAlgorithmSHA256, digest, "digest.bin")
}

func xadesRefsSignedXMLDOM(t *testing.T) *xmldom.Node {
	t.Helper()
	document, err := xmlutils.DomUtilsBuildDOMFromString(xadesRefsXMLSignedContent)
	if err != nil {
		t.Fatalf("building the fixture DOM: %v", err)
	}
	return document
}

func xadesRefsXSLTDOM(t *testing.T) *xmldom.Node {
	t.Helper()
	document, err := xmlutils.DomUtilsBuildDOMFromString(xadesRefsXSLTContent)
	if err != nil {
		t.Fatalf("building the stylesheet DOM: %v", err)
	}
	return document
}

// xadesRefsSigningDate is the frozen signing date every oracle case uses, the same one
// gen/RefsOracle.java freezes.
var xadesRefsSigningDate = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

func xadesRefsBaseParams(t *testing.T) *SignatureParameters {
	t.Helper()
	params := NewXAdESSignatureParameters()
	signer := xadesSignABuilderSigner(t)
	params.SetSigningCertificate(signer)
	params.SetCertificateChain([]*model.CertificateToken{signer})
	signingDate := xadesRefsSigningDate
	params.BLevel().SetSigningDate(&signingDate)
	params.SetSignatureLevel(enumerations.SignatureLevelXAdESBaselineB)
	params.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	return params
}

func xadesRefsParamsFor(t *testing.T, packaging enumerations.SignaturePackaging) *SignatureParameters {
	t.Helper()
	params := xadesRefsBaseParams(t)
	params.SetSignaturePackaging(packaging)
	return params
}

func xadesRefsReadAll(t *testing.T, document model.DSSDocument) []byte {
	t.Helper()
	reader, err := document.OpenStream()
	if err != nil {
		t.Fatalf("opening the document: %v", err)
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	return content
}

// ------------------------------------------------------- 1. the ds:Transforms a transform writes

// TestDSSTransformCreateTransformAgainstJavaOracle pins what every transform in the hierarchy
// writes into a ds:Reference. The three namespace flavours - the stock "ds" prefix, a custom
// "dsig" one and the empty prefix that makes XMLDSIG the default namespace - are what make the
// XPath Filter 2.0 transform's two namespace declarations, and the expression it builds from the
// prefix, observable.
func TestDSSTransformCreateTransformAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	cases := []struct {
		name       string
		namespace  *common.DSSNamespace
		transforms func() []DSSTransform
	}{
		{"tf-base64", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewBase64Transform()}
		}},
		{"tf-base64-dsig", xadesRefsDsigNamespace, func() []DSSTransform {
			return []DSSTransform{NewBase64TransformWithNamespace(xadesRefsDsigNamespace)}
		}},
		{"tf-base64-default-prefix", xadesRefsDefaultNamespace, func() []DSSTransform {
			return []DSSTransform{NewBase64TransformWithNamespace(xadesRefsDefaultNamespace)}
		}},

		{"tf-c14n-exclusive", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NExclusive)}
		}},
		{"tf-c14n-inclusive", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NInclusive)}
		}},
		{"tf-c14n-11-comments", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14N11Comments)}
		}},
		{"tf-c14n-exclusive-dsig", xadesRefsDsigNamespace, func() []DSSTransform {
			return []DSSTransform{
				NewCanonicalizationTransformWithNamespace(xadesRefsDsigNamespace, xadesRefsC14NExclusive)}
		}},

		{"tf-enveloped", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewEnvelopedSignatureTransform()}
		}},
		{"tf-enveloped-default-prefix", xadesRefsDefaultNamespace, func() []DSSTransform {
			return []DSSTransform{NewEnvelopedSignatureTransformWithNamespace(xadesRefsDefaultNamespace)}
		}},

		{"tf-xpath", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewXPathTransform("//*[local-name()='child']")}
		}},
		{"tf-xpath-dsig", xadesRefsDsigNamespace, func() []DSSTransform {
			return []DSSTransform{
				NewXPathTransformWithNamespace(xadesRefsDsigNamespace, "//*[local-name()='child']")}
		}},
		{"tf-xpath-enveloped", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewXPathEnvelopedSignatureTransform()}
		}},
		{"tf-xpath-enveloped-dsig", xadesRefsDsigNamespace, func() []DSSTransform {
			return []DSSTransform{NewXPathEnvelopedSignatureTransformWithNamespace(xadesRefsDsigNamespace)}
		}},

		{"tf-xpath2filter", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewXPath2FilterTransform("/root/child", "intersect")}
		}},
		{"tf-xpath2filter-enveloped", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewXPath2FilterEnvelopedSignatureTransform()}
		}},
		{"tf-xpath2filter-enveloped-dsig", xadesRefsDsigNamespace, func() []DSSTransform {
			return []DSSTransform{
				NewXPath2FilterEnvelopedSignatureTransformWithNamespace(xadesRefsDsigNamespace)}
		}},
		// The one case that reaches the "xmlns" (no prefix) branch of XPath2FilterTransform, and
		// the one that shows the expression upstream builds from an empty prefix.
		{"tf-xpath2filter-enveloped-default-prefix", xadesRefsDefaultNamespace, func() []DSSTransform {
			return []DSSTransform{
				NewXPath2FilterEnvelopedSignatureTransformWithNamespace(xadesRefsDefaultNamespace)}
		}},

		{"tf-xslt", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewXsltTransform(xadesRefsXSLTDOM(t))}
		}},
		{"tf-spdoc-digest", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{NewSPDocDigestAsInSpecificationTransform()}
		}},

		{"tf-chain-enveloped-c14n", xadesRefsDSNamespace, func() []DSSTransform {
			return []DSSTransform{
				NewXPath2FilterEnvelopedSignatureTransform(),
				NewCanonicalizationTransform(xadesRefsC14NExclusive),
			}
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			transforms := testCase.transforms()

			document := xmlutils.DomUtilsBuildDOMEmpty()
			referenceDom := xmlutils.DomUtilsCreateElementNS(document, testCase.namespace,
				common.XMLDSigElementReference)
			document.AppendChild(referenceDom)
			DSSXMLUtilsIncorporateTransforms(referenceDom, transforms, testCase.namespace)

			serialized, err := xmlutils.DomUtilsSerializeNode(referenceDom)
			if err != nil {
				t.Fatalf("serializing the ds:Reference: %v", err)
			}
			xadesRefsAssertBytes(t, testCase.name, "transforms",
				oracle.expect(t, testCase.name, "transforms"), serialized)

			algorithms := make([]string, 0, len(transforms))
			for _, transform := range transforms {
				algorithms = append(algorithms, transform.Algorithm())
			}
			xadesRefsAssertString(t, testCase.name, "algorithm",
				oracle.expectString(t, testCase.name, "algorithm"), strings.Join(algorithms, "\n"))
		})
	}
}

// ------------------------------------------------------------- 2. the transforms as executed

// TestApplyTransformsAgainstJavaOracle pins the octets a transform chain produces, which is what
// a reference digest is taken over. This is the Santuario path: everything digested goes through
// internal/xmldsig, and these bytes are the proof that the routing is byte-exact.
func TestApplyTransformsAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	cases := []struct {
		name       string
		transforms func() []DSSTransform
	}{
		{"apply-no-transforms", func() []DSSTransform { return nil }},

		{"apply-c14n-exclusive", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NExclusive)}
		}},
		{"apply-c14n-inclusive", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NInclusive)}
		}},
		{"apply-c14n-inclusive-comments", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NInclusiveComments)}
		}},
		{"apply-c14n-11", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14N11)}
		}},
		{"apply-c14n-exclusive-comments", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NExclusiveComments)}
		}},

		{"apply-xpath2filter-enveloped-then-c14n", func() []DSSTransform {
			return []DSSTransform{
				NewXPath2FilterEnvelopedSignatureTransform(),
				NewCanonicalizationTransform(xadesRefsC14NExclusive),
			}
		}},
		{"apply-xpath2filter-enveloped", func() []DSSTransform {
			return []DSSTransform{NewXPath2FilterEnvelopedSignatureTransform()}
		}},
		{"apply-xpath-enveloped-then-c14n", func() []DSSTransform {
			return []DSSTransform{
				NewXPathEnvelopedSignatureTransform(),
				NewCanonicalizationTransform(xadesRefsC14NExclusive),
			}
		}},
		{"apply-xpath-enveloped", func() []DSSTransform {
			return []DSSTransform{NewXPathEnvelopedSignatureTransform()}
		}},
		{"apply-xpath-select-child", func() []DSSTransform {
			return []DSSTransform{
				NewXPathTransform("ancestor-or-self::*[local-name()='child']"),
				NewCanonicalizationTransform(xadesRefsC14NExclusive),
			}
		}},
		{"apply-xpath2filter-intersect", func() []DSSTransform {
			return []DSSTransform{
				NewXPath2FilterTransform("//*[local-name()='child']", "intersect"),
				NewCanonicalizationTransform(xadesRefsC14NExclusive),
			}
		}},

		// Base64Transform#performTransform is the identity, so the octets are the node's.
		{"apply-base64", func() []DSSTransform {
			return []DSSTransform{NewBase64Transform()}
		}},

		// Rejections, both of them upstream's own.
		{"apply-spdoc-digest", func() []DSSTransform {
			return []DSSTransform{NewSPDocDigestAsInSpecificationTransform()}
		}},
		{"apply-c14n-physical", func() []DSSTransform {
			return []DSSTransform{NewCanonicalizationTransform(xadesRefsC14NPhysical)}
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := DSSXMLUtilsApplyTransforms(xadesRefsSignedXMLDOM(t), testCase.transforms())
			if _, produced := oracle[testCase.name]["output"]; !produced {
				xadesRefsAssertError(t, oracle, testCase.name, err)
				return
			}
			if err != nil {
				t.Fatalf("case %s: ApplyTransforms: %v", testCase.name, err)
			}
			xadesRefsAssertBytes(t, testCase.name, "output",
				oracle.expect(t, testCase.name, "output"), output)
		})
	}

	// An identity transform followed by a real one. The chain is the combination
	// EnvelopedSignatureTransform's own javadoc prescribes ("must be followed up by a
	// CanonicalizationTransform"). It used to trip a defect in internal/xmldsig, where
	// Data.IsElement/IsNodeSet consulted the octet cache Data.Bytes leaves behind, so a value
	// that had been read once reported "no usable state" to the next transform; Santuario's
	// isElement()/isNodeSet() look at inputOctetStreamProxy and never at the cached bytes. That
	// was fixed by dropping the !hasOctets term from both predicates (see the FIX note in
	// internal/xmldsig/data.go), and this case asserts the repaired behaviour outright - it must
	// never be softened back into a skip, since the transform chain silently degrading is
	// precisely the regression it exists to catch.
	t.Run("apply-enveloped-then-c14n", func(t *testing.T) {
		const name = "apply-enveloped-then-c14n"
		output, err := DSSXMLUtilsApplyTransforms(xadesRefsSignedXMLDOM(t), []DSSTransform{
			NewEnvelopedSignatureTransform(),
			NewCanonicalizationTransform(xadesRefsC14NExclusive),
		})
		if err != nil {
			t.Fatalf("case %s: ApplyTransforms: %v", name, err)
		}
		xadesRefsAssertBytes(t, name, "output", oracle.expect(t, name, "output"), output)
	})

	// The XSLT transform. internal/xmldsig registers it as a transform that always refuses, so
	// this port does not run a stylesheet that arrives inside a signature - upstream does, and
	// the oracle records its output. The divergence is deliberate and documented in
	// xslt_transform.go and in internal/xmldsig's DefaultRegistry; this case asserts the
	// refusal so that a silent change of behaviour cannot slip through.
	t.Run("apply-xslt", func(t *testing.T) {
		_, err := DSSXMLUtilsApplyTransforms(xadesRefsSignedXMLDOM(t),
			[]DSSTransform{NewXsltTransform(xadesRefsXSLTDOM(t))})
		if err == nil {
			t.Fatal("the XSLT transform ran; internal/xmldsig is meant to refuse it")
		}
		if !errors.Is(err, xmldsig.ErrForbiddenTransform) &&
			!strings.Contains(err.Error(), xmldsig.ErrForbiddenTransform.Error()) {
			t.Errorf("XSLT refused with %v, want the forbidden-transform refusal", err)
		}
		if want := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><out/>"; string(
			oracle.expect(t, "apply-xslt", "output")) != want {
			t.Errorf("the oracle no longer records upstream's XSLT output as %q; re-check the "+
				"divergence recorded in xslt_transform.go", want)
		}
	})
}

// -------------------------------------------------------------------- 3. ReferenceIdProvider

func TestReferenceIdProviderAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	ids := func(provider *ReferenceIdProvider, count int) string {
		generated := make([]string, 0, count)
		for i := 0; i < count; i++ {
			generated = append(generated, provider.ReferenceId())
		}
		return strings.Join(generated, "\n")
	}

	t.Run("refid-no-parameters", func(t *testing.T) {
		xadesRefsAssertString(t, "refid-no-parameters", "ids",
			oracle.expectString(t, "refid-no-parameters", "ids"), ids(NewReferenceIdProvider(), 3))
	})

	t.Run("refid-with-parameters", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		xadesRefsAssertString(t, "refid-with-parameters", "deterministic-id",
			oracle.expectString(t, "refid-with-parameters", "deterministic-id"), params.GetDeterministicId())

		provider := NewReferenceIdProvider()
		provider.SetSignatureParameters(params)
		xadesRefsAssertString(t, "refid-with-parameters", "ids",
			oracle.expectString(t, "refid-with-parameters", "ids"), ids(provider, 3))
	})

	t.Run("refid-custom-prefix", func(t *testing.T) {
		provider := NewReferenceIdProvider()
		provider.SetReferenceIdPrefix("m-id")
		xadesRefsAssertString(t, "refid-custom-prefix", "ids",
			oracle.expectString(t, "refid-custom-prefix", "ids"), ids(provider, 2))
	})

	t.Run("refid-blank-prefix", func(t *testing.T) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				t.Fatal("SetReferenceIdPrefix accepted a blank prefix")
			}
			message, ok := recovered.(string)
			if !ok {
				t.Fatalf("panicked with %T, want the Java message as a string", recovered)
			}
			xadesRefsAssertString(t, "refid-blank-prefix", "error",
				xadesRefsErrorMessage(oracle.expectString(t, "refid-blank-prefix", "error")), message)
		}()
		NewReferenceIdProvider().SetReferenceIdPrefix("  ")
	})
}

// ------------------------------------------------------------------------ 4. ReferenceBuilder

// TestReferenceBuilderAgainstJavaOracle pins the DSSReferences ReferenceBuilder derives from the
// signature parameters, field by field - including whether the URI is Java's null or the empty
// string, which is what makes the difference between writing URI="" and writing no URI at all.
func TestReferenceBuilderAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	assertReferences := func(t *testing.T, name string, references []*DSSReference) {
		t.Helper()
		xadesRefsAssertString(t, name, "count", oracle.expectString(t, name, "count"),
			strconv.Itoa(len(references)))
		for i, reference := range references {
			key := "ref" + strconv.Itoa(i)
			xadesRefsAssertString(t, name, key+"-id",
				oracle.expectString(t, name, key+"-id"), xadesRefsNullable(reference.Id()))
			xadesRefsAssertString(t, name, key+"-uri-present",
				oracle.expectString(t, name, key+"-uri-present"), strconv.FormatBool(reference.HasUri()))
			uri := xadesRefsNull
			if reference.HasUri() {
				uri = reference.Uri()
			}
			xadesRefsAssertString(t, name, key+"-uri", oracle.expectString(t, name, key+"-uri"), uri)
			xadesRefsAssertString(t, name, key+"-type",
				oracle.expectString(t, name, key+"-type"), xadesRefsNullable(reference.Type()))
			xadesRefsAssertString(t, name, key+"-digest-method",
				oracle.expectString(t, name, key+"-digest-method"), string(reference.DigestMethodAlgorithm()))
			algorithms := make([]string, 0, len(reference.Transforms()))
			for _, transform := range reference.Transforms() {
				algorithms = append(algorithms, transform.Algorithm())
			}
			xadesRefsAssertString(t, name, key+"-transforms",
				oracle.expectString(t, name, key+"-transforms"), strings.Join(algorithms, "\n"))
		}
	}

	build := func(t *testing.T, params *SignatureParameters,
		documents []model.DSSDocument) ([]*DSSReference, error) {
		t.Helper()
		provider := NewReferenceIdProvider()
		provider.SetSignatureParameters(params)
		return NewReferenceBuilder(documents, params, provider).Build()
	}

	cases := []struct {
		name      string
		params    func(t *testing.T) *SignatureParameters
		documents func(t *testing.T) []model.DSSDocument
	}{
		{"refbuild-enveloped", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsXMLDocument()} }},

		{"refbuild-enveloping", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},

		{"refbuild-detached", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},

		// The document has no name, so upstream leaves the URI null - not empty.
		{"refbuild-detached-unnamed", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsUnnamedDocument()} }},

		{"refbuild-internally-detached", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingInternallyDetached)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsXMLDocument()} }},

		{"refbuild-enveloping-embed-xml", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetEmbedXML(true)
			return params
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsXMLDocument()} }},

		{"refbuild-enveloping-manifest", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetManifestSignature(true)
			return params
		}, func(t *testing.T) []model.DSSDocument {
			return []model.DSSDocument{xadesRefsManifestDocument(t)}
		}},

		{"refbuild-two-documents", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached)
		}, func(t *testing.T) []model.DSSDocument {
			return []model.DSSDocument{xadesRefsTextDocument(), xadesRefsSecondDocument()}
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			references, err := build(t, testCase.params(t), testCase.documents(t))
			if err != nil {
				t.Fatalf("case %s: Build: %v", testCase.name, err)
			}
			assertReferences(t, testCase.name, references)
		})
	}

	// The ManifestBuilder path: the detached constructor, without signature parameters.
	t.Run("refbuild-no-parameters", func(t *testing.T) {
		references, err := NewReferenceBuilderWithDigestAlgorithm(
			[]model.DSSDocument{xadesRefsTextDocument(), xadesRefsXMLDocument()},
			enumerations.DigestAlgorithmSHA256, NewReferenceIdProvider()).Build()
		if err != nil {
			t.Fatalf("Build: %v", err)
		}
		assertReferences(t, "refbuild-no-parameters", references)
	})

	// Rejections.
	rejections := []struct {
		name      string
		params    func(t *testing.T) *SignatureParameters
		documents func(t *testing.T) []model.DSSDocument
	}{
		{"refbuild-enveloped-not-xml", func(t *testing.T) *SignatureParameters {
			return xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},

		{"refbuild-manifest-without-id", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetManifestSignature(true)
			return params
		}, func(t *testing.T) []model.DSSDocument {
			return []model.DSSDocument{xadesRefsXMLWithoutIDDocument()}
		}},

		{"refbuild-embed-xml-not-xml", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetEmbedXML(true)
			return params
		}, func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},
	}

	for _, testCase := range rejections {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := build(t, testCase.params(t), testCase.documents(t))
			xadesRefsAssertError(t, oracle, testCase.name, err)
		})
	}
}

// xadesRefsNullable renders "" the way RefsOracle.java renders a Java null, for the fields where
// the two coincide (id and type are never assigned the empty string by DSS).
func xadesRefsNullable(value string) string {
	if value == "" {
		return xadesRefsNull
	}
	return value
}

// ---------------------------------------------------------------------- 5. ReferenceProcessor

// TestReferenceProcessorAgainstJavaOracle pins the complete ds:Reference list, digests included.
// This is the chunk's central expectation: the Id/URI/Type attributes and their order, the
// ds:Transforms subtree, the ds:DigestMethod, and a ds:DigestValue computed over the octets the
// transform chain really produced. The enveloped cases are also what proves URI="" is written -
// upstream distinguishes an empty URI from an absent one.
func TestReferenceProcessorAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	assertIncorporated := func(t *testing.T, name string, processor *ReferenceProcessor,
		references []*DSSReference, namespace *common.DSSNamespace) {
		t.Helper()
		document := xmlutils.DomUtilsBuildDOMEmpty()
		signedInfoDom := xmlutils.DomUtilsCreateElementNS(document, namespace, common.XMLDSigElementSignedInfo)
		document.AppendChild(signedInfoDom)
		if err := processor.IncorporateReferences(signedInfoDom, references, namespace); err != nil {
			t.Fatalf("case %s: IncorporateReferences: %v", name, err)
		}
		serialized, err := xmlutils.DomUtilsSerializeNode(signedInfoDom)
		if err != nil {
			t.Fatalf("case %s: serializing ds:SignedInfo: %v", name, err)
		}
		xadesRefsAssertBytes(t, name, "references", oracle.expect(t, name, "references"), serialized)

		for i, reference := range references {
			key := "output" + strconv.Itoa(i)
			output, err := processor.ReferenceOutput(reference)
			if err != nil {
				t.Fatalf("case %s: ReferenceOutput(%d): %v", name, i, err)
			}
			if digest, isDigestOnly := oracle[name][key+"-digest"]; isDigestOnly {
				got, err := output.DigestValue(enumerations.DigestAlgorithmSHA256)
				if err != nil {
					t.Fatalf("case %s: digesting the reference output: %v", name, err)
				}
				xadesRefsAssertBytes(t, name, key+"-digest", digest, got)
				continue
			}
			xadesRefsAssertBytes(t, name, key, oracle.expect(t, name, key), xadesRefsReadAll(t, output))
		}
	}

	incorporate := func(t *testing.T, name string, params *SignatureParameters,
		documents []model.DSSDocument, namespace *common.DSSNamespace) {
		t.Helper()
		provider := NewReferenceIdProvider()
		provider.SetSignatureParameters(params)
		references, err := NewReferenceBuilder(documents, params, provider).Build()
		if err != nil {
			t.Fatalf("case %s: Build: %v", name, err)
		}
		assertIncorporated(t, name, NewReferenceProcessor(params), references, namespace)
	}

	t.Run("incorporate-enveloped", func(t *testing.T) {
		incorporate(t, "incorporate-enveloped",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped),
			[]model.DSSDocument{xadesRefsXMLDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-enveloping", func(t *testing.T) {
		incorporate(t, "incorporate-enveloping",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping),
			[]model.DSSDocument{xadesRefsTextDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-detached", func(t *testing.T) {
		incorporate(t, "incorporate-detached",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached),
			[]model.DSSDocument{xadesRefsTextDocument()}, xadesRefsDSNamespace)
	})
	// The URI is null here, so no URI attribute is written at all.
	t.Run("incorporate-detached-unnamed", func(t *testing.T) {
		incorporate(t, "incorporate-detached-unnamed",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached),
			[]model.DSSDocument{xadesRefsUnnamedDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-internally-detached", func(t *testing.T) {
		incorporate(t, "incorporate-internally-detached",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingInternallyDetached),
			[]model.DSSDocument{xadesRefsXMLDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-two-documents", func(t *testing.T) {
		incorporate(t, "incorporate-two-documents",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached),
			[]model.DSSDocument{xadesRefsTextDocument(), xadesRefsSecondDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-enveloped-dsig-prefix", func(t *testing.T) {
		incorporate(t, "incorporate-enveloped-dsig-prefix",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped),
			[]model.DSSDocument{xadesRefsXMLDocument()}, xadesRefsDsigNamespace)
	})
	t.Run("incorporate-enveloped-default-prefix", func(t *testing.T) {
		incorporate(t, "incorporate-enveloped-default-prefix",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped),
			[]model.DSSDocument{xadesRefsXMLDocument()}, xadesRefsDefaultNamespace)
	})
	t.Run("incorporate-embed-xml", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		params.SetEmbedXML(true)
		incorporate(t, "incorporate-embed-xml", params,
			[]model.DSSDocument{xadesRefsXMLDocument()}, xadesRefsDSNamespace)
	})
	t.Run("incorporate-manifest", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		params.SetManifestSignature(true)
		incorporate(t, "incorporate-manifest", params,
			[]model.DSSDocument{xadesRefsManifestDocument(t)}, xadesRefsDSNamespace)
	})

	// The ManifestBuilder path, with the parameter-less processor.
	t.Run("incorporate-no-parameters", func(t *testing.T) {
		references, err := NewReferenceBuilderWithDigestAlgorithm(
			[]model.DSSDocument{xadesRefsTextDocument(), xadesRefsXMLDocument()},
			enumerations.DigestAlgorithmSHA256, NewReferenceIdProvider()).Build()
		if err != nil {
			t.Fatalf("Build: %v", err)
		}
		assertIncorporated(t, "incorporate-no-parameters", NewReferenceProcessorEmpty(),
			references, xadesRefsDSNamespace)
	})

	// A DigestDocument short-circuits getReferenceOutput: no transform runs and the recorded
	// digest is reused verbatim.
	t.Run("incorporate-digest-document", func(t *testing.T) {
		reference := NewDSSReference()
		reference.SetId("r-digest")
		reference.SetUri("digest.bin")
		reference.SetDigestMethodAlgorithm(enumerations.DigestAlgorithmSHA256)
		reference.SetContents(xadesRefsDigestDocument(t))
		assertIncorporated(t, "incorporate-digest-document", NewReferenceProcessorEmpty(),
			[]*DSSReference{reference}, xadesRefsDSNamespace)
	})

	// An explicit reference over already-signed XML, filtering the existing ds:Signature away:
	// the digest is taken over the filtered node set, not the whole document.
	t.Run("incorporate-filtered-signature", func(t *testing.T) {
		reference := NewDSSReference()
		reference.SetId("r-filtered")
		reference.SetUri("")
		reference.SetDigestMethodAlgorithm(enumerations.DigestAlgorithmSHA256)
		reference.SetContents(xadesRefsSignedXMLDocument())
		reference.SetTransforms([]DSSTransform{
			NewXPath2FilterEnvelopedSignatureTransform(),
			NewCanonicalizationTransform(xadesRefsC14NExclusive),
		})
		assertIncorporated(t, "incorporate-filtered-signature", NewReferenceProcessorEmpty(),
			[]*DSSReference{reference}, xadesRefsDSNamespace)
	})
}

// ----------------------------------------------------------------------- 6. ReferenceVerifier

func TestReferenceVerifierAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	reference := func(id, uri string, contents model.DSSDocument, transforms []DSSTransform) *DSSReference {
		built := NewDSSReference()
		if id != "" {
			built.SetId(id)
		}
		built.SetUri(uri)
		built.SetDigestMethodAlgorithm(enumerations.DigestAlgorithmSHA256)
		built.SetContents(contents)
		built.SetTransforms(transforms)
		return built
	}

	// The three cases where the builder's own references are verified.
	built := []struct {
		name      string
		packaging enumerations.SignaturePackaging
		documents func(t *testing.T) []model.DSSDocument
	}{
		{"verify-enveloped-ok", enumerations.SignaturePackagingEnveloped,
			func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsXMLDocument()} }},
		{"verify-enveloping-ok", enumerations.SignaturePackagingEnveloping,
			func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},
		{"verify-detached-ok", enumerations.SignaturePackagingDetached,
			func(t *testing.T) []model.DSSDocument { return []model.DSSDocument{xadesRefsTextDocument()} }},
	}
	for _, testCase := range built {
		t.Run(testCase.name, func(t *testing.T) {
			params := xadesRefsParamsFor(t, testCase.packaging)
			provider := NewReferenceIdProvider()
			provider.SetSignatureParameters(params)
			references, err := NewReferenceBuilder(testCase.documents(t), params, provider).Build()
			if err != nil {
				t.Fatalf("case %s: Build: %v", testCase.name, err)
			}
			params.SetReferences(references)
			xadesRefsAssertError(t, oracle, testCase.name, NewReferenceVerifier(params).CheckReferencesValidity())
		})
	}

	// The rejections, each with hand-assembled references, exactly as the oracle assembles them.
	explicit := []struct {
		name   string
		params func(t *testing.T) *SignatureParameters
	}{
		{"verify-enveloped-without-transform", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
			params.SetReferences([]*DSSReference{reference("r-1", "", xadesRefsXMLDocument(), nil)})
			return params
		}},
		{"verify-base64-embed-xml", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetEmbedXML(true)
			params.SetReferences([]*DSSReference{reference("r-1", "#o-r-1", xadesRefsXMLDocument(),
				[]DSSTransform{NewBase64Transform()})})
			return params
		}},
		{"verify-base64-manifest", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetManifestSignature(true)
			params.SetReferences([]*DSSReference{reference("r-1", "#o-r-1", xadesRefsXMLDocument(),
				[]DSSTransform{NewBase64Transform()})})
			return params
		}},
		{"verify-base64-detached", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached)
			params.SetReferences([]*DSSReference{reference("r-1", "hello.txt", xadesRefsTextDocument(),
				[]DSSTransform{NewBase64Transform()})})
			return params
		}},
		{"verify-base64-with-other-transform", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
			params.SetReferences([]*DSSReference{reference("r-1", "#o-r-1", xadesRefsTextDocument(),
				[]DSSTransform{
					NewBase64Transform(),
					NewCanonicalizationTransform(xadesRefsC14NExclusive),
				})})
			return params
		}},
		{"verify-missing-element-id", func(t *testing.T) *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingInternallyDetached)
			params.SetReferences([]*DSSReference{reference("r-1", "#absent-id", xadesRefsXMLDocument(),
				[]DSSTransform{NewCanonicalizationTransform(xadesRefsC14NExclusive)})})
			return params
		}},
	}
	for _, testCase := range explicit {
		t.Run(testCase.name, func(t *testing.T) {
			xadesRefsAssertError(t, oracle, testCase.name,
				NewReferenceVerifier(testCase.params(t)).CheckReferencesValidity())
		})
	}

	// A reference without an Id is given a deterministic one IN PLACE.
	t.Run("verify-generates-missing-id", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		withoutId := reference("", "#o-r-1", xadesRefsTextDocument(),
			[]DSSTransform{NewBase64Transform()})
		params.SetReferences([]*DSSReference{withoutId})
		xadesRefsAssertError(t, oracle, "verify-generates-missing-id",
			NewReferenceVerifier(params).CheckReferencesValidity())
		xadesRefsAssertString(t, "verify-generates-missing-id", "assigned-id",
			oracle.expectString(t, "verify-generates-missing-id", "assigned-id"), withoutId.Id())
	})
}

// ------------------------------------------------------------ 7. ds:SignedInfo, end to end

// TestSignedInfoWithReferencesAgainstJavaOracle is the mandated SignedInfo byte KAT: for
// fixed inputs, the ds:SignedInfo element upstream's XAdESSignatureBuilder produces - before and
// after canonicalization - and the signature document built from it, across every reference and
// transform configuration this package is responsible for.
func TestSignedInfoWithReferencesAgainstJavaOracle(t *testing.T) {
	oracle := loadXAdESRefsOracle(t)

	run := func(t *testing.T, name string, params *SignatureParameters, documents []model.DSSDocument) {
		t.Helper()
		builderRef, err := SignatureBuilderGetSignatureBuilderForDocuments(params, documents,
			validation.NewCommonCertificateVerifier())
		if err != nil {
			t.Fatalf("case %s: building the signature builder: %v", name, err)
		}
		canonicalized, err := builderRef.Build()
		if err != nil {
			t.Fatalf("case %s: Build: %v", name, err)
		}
		xadesRefsAssertBytes(t, name, "signedinfo-c14n", oracle.expect(t, name, "signedinfo-c14n"), canonicalized)

		base := xadesSignABuilderBaseOf(t, builderRef)
		raw, err := xmlutils.DomUtilsSerializeNode(base.SignedInfoDom)
		if err != nil {
			t.Fatalf("case %s: serializing ds:SignedInfo: %v", name, err)
		}
		xadesRefsAssertBytes(t, name, "signedinfo-raw", oracle.expect(t, name, "signedinfo-raw"), raw)

		signed, err := builderRef.SignDocument(xadesSignABuilderSignatureValue())
		if err != nil {
			t.Fatalf("case %s: SignDocument: %v", name, err)
		}
		xadesRefsAssertBytes(t, name, "document", oracle.expect(t, name, "document"),
			xadesRefsReadAll(t, signed))
	}

	t.Run("si-enveloped", func(t *testing.T) {
		run(t, "si-enveloped", xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped),
			[]model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-enveloping", func(t *testing.T) {
		run(t, "si-enveloping", xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping),
			[]model.DSSDocument{xadesRefsTextDocument()})
	})
	t.Run("si-detached", func(t *testing.T) {
		run(t, "si-detached", xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached),
			[]model.DSSDocument{xadesRefsTextDocument()})
	})
	t.Run("si-internally-detached", func(t *testing.T) {
		run(t, "si-internally-detached",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingInternallyDetached),
			[]model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-enveloping-embed-xml", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		params.SetEmbedXML(true)
		run(t, "si-enveloping-embed-xml", params, []model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-enveloping-manifest", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloping)
		params.SetManifestSignature(true)
		run(t, "si-enveloping-manifest", params, []model.DSSDocument{xadesRefsManifestDocument(t)})
	})
	t.Run("si-enveloped-xpath-transform", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
		explicit := NewDSSReference()
		explicit.SetId("r-1")
		explicit.SetUri("")
		explicit.SetDigestMethodAlgorithm(enumerations.DigestAlgorithmSHA256)
		explicit.SetContents(xadesRefsXMLDocument())
		explicit.SetTransforms([]DSSTransform{
			NewXPathEnvelopedSignatureTransform(),
			NewCanonicalizationTransform(xadesRefsC14NExclusive),
		})
		params.SetReferences([]*DSSReference{explicit})
		run(t, "si-enveloped-xpath-transform", params, []model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-enveloped-dsig-prefix", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
		params.SetXmldsigNamespace(xadesRefsDsigNamespace)
		run(t, "si-enveloped-dsig-prefix", params, []model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-enveloped-default-prefix", func(t *testing.T) {
		params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
		params.SetXmldsigNamespace(xadesRefsDefaultNamespace)
		run(t, "si-enveloped-default-prefix", params, []model.DSSDocument{xadesRefsXMLDocument()})
	})
	t.Run("si-detached-two-documents", func(t *testing.T) {
		run(t, "si-detached-two-documents",
			xadesRefsParamsFor(t, enumerations.SignaturePackagingDetached),
			[]model.DSSDocument{xadesRefsTextDocument(), xadesRefsSecondDocument()})
	})

	// The EnvelopedSignatureTransform is an identity on the creation path, which used to leave
	// the following transform without usable state; see apply-enveloped-then-c14n above for the
	// diagnosis and the internal/xmldsig fix. Asserted outright, never skipped.
	t.Run("si-enveloped-enveloped-transform", func(t *testing.T) {
		newParams := func() *SignatureParameters {
			params := xadesRefsParamsFor(t, enumerations.SignaturePackagingEnveloped)
			explicit := NewDSSReference()
			explicit.SetId("r-1")
			explicit.SetUri("")
			explicit.SetDigestMethodAlgorithm(enumerations.DigestAlgorithmSHA256)
			explicit.SetContents(xadesRefsXMLDocument())
			explicit.SetTransforms([]DSSTransform{
				NewEnvelopedSignatureTransform(),
				NewCanonicalizationTransform(xadesRefsC14NInclusive),
			})
			params.SetReferences([]*DSSReference{explicit})
			return params
		}

		run(t, "si-enveloped-enveloped-transform", newParams(), []model.DSSDocument{xadesRefsXMLDocument()})
	})
}

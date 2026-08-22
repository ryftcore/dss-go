// Known-answer test for xades_builder.go against a Java oracle.
//
// testdata/xades-builder-oracle.json is produced by testdata/gen/XAdESBuilderOracle.java, which
// drives upstream DSS 6.5.RC1's own XAdESBuilder (see that file's header for the exact command).
// Nothing here is hand-derived: for every (XAdES namespace x en319132 x digest algorithm)
// combination the golden holds the serialized DOM upstream builds, and this test rebuilds the
// same fragment with the Go port and compares the serialization byte for byte. That is the
// byte-compatibility contract this package owes the signature it ends up inside: element order,
// where the xmlns declarations land, the ds:/xades: prefixes, and the Qualifier/Algorithm
// attribute values.
//
// The golden also records the two namespace/element combinations upstream refuses
// (IssuerSerialV2 under XAdES 1.1.1 and 1.2.2, which XAdES11xElement answers with an
// UnsupportedOperationException) so the Go port has to refuse them too, and the deterministic
// XML Id XAdESBuilder#toXmlIdentifier derives from a DSS Identifier.
package xades

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// xadesBuilderKATCase is one row of testdata/xades-builder-oracle.json.
type xadesBuilderKATCase struct {
	Operation           string `json:"operation"`
	XadesNamespaceUri   string `json:"xadesNamespaceUri"`
	XadesNamespacePrefx string `json:"xadesNamespacePrefix"`
	En319132            bool   `json:"en319132"`
	DigestAlgorithm     string `json:"digestAlgorithm"`
	XML                 string `json:"xml"`
	Unsupported         string `json:"unsupported"`
	AsXmlID             string `json:"asXmlId"`
	XmlIdentifier       string `json:"xmlIdentifier"`
}

// xadesBuilderKATBuilder is the concrete Builder the KAT drives, mirroring the anonymous
// subclass the Java oracle uses: XAdESBuilder is abstract only in AlignNodes, and the KAT never
// reaches CreateXmlDocument.
type xadesBuilderKATBuilder struct {
	Builder
}

func (b *xadesBuilderKATBuilder) AlignNodes() {
	// never reached: the KAT never calls CreateXmlDocument
}

func newXAdESBuilderKATBuilder(namespace *common.DSSNamespace, en319132 bool) (*xadesBuilderKATBuilder, *xmldom.Node) {
	params := NewXAdESSignatureParameters()
	params.SetXadesNamespace(namespace)
	params.SetEn319132(en319132)

	document := xmlutils.DomUtilsBuildDOMEmpty()
	root := xmldom.NewElement(xmldom.Name{Space: "urn:dss:go:oracle", Local: "Root", Prefix: "o"})
	root.SetAttr(xmldom.Name{Space: xmldom.XMLNSNamespace, Local: "o", Prefix: "xmlns"}, "urn:dss:go:oracle")
	document.AppendChild(root)

	builder := &xadesBuilderKATBuilder{}
	builder.InitXAdESBuilder(builder)
	builder.Params = params
	builder.DocumentDom = document
	return builder, root
}

func xadesBuilderKATNamespace(t *testing.T, uri string) *common.DSSNamespace {
	t.Helper()
	for _, namespace := range []*common.DSSNamespace{
		definition.XAdESNamespaceXAdES132,
		definition.XAdESNamespaceXAdES122,
		definition.XAdESNamespaceXAdES111,
	} {
		if namespace.Uri() == uri {
			return namespace
		}
	}
	t.Fatalf("unknown XAdES namespace in the oracle: %s", uri)
	return nil
}

func xadesBuilderKATCases(t *testing.T) []xadesBuilderKATCase {
	t.Helper()
	raw, err := os.ReadFile(corpustest.Path(t, "xades-builder-oracle.json"))
	if err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	var cases []xadesBuilderKATCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("parse oracle: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("the oracle is empty")
	}
	return cases
}

func xadesBuilderKATCertificate(t *testing.T) *model.CertificateToken {
	t.Helper()
	certificate, err := spi.DSSUtilsLoadCertificate("testdata/signer.crt")
	if err != nil {
		t.Fatalf("load signer.crt: %v", err)
	}
	return certificate
}

// TestXAdESBuilderDOMMatchesJavaOracle rebuilds every DOM fragment the oracle dumped and
// compares the serialization with upstream's, byte for byte.
func TestXAdESBuilderDOMMatchesJavaOracle(t *testing.T) {
	certificate := xadesBuilderKATCertificate(t)

	for _, testCase := range xadesBuilderKATCases(t) {
		if testCase.Operation == "toXmlIdentifier" {
			continue
		}
		testCase := testCase
		name := testCase.Operation + "/" + testCase.XadesNamespacePrefx + "/" +
			testCase.DigestAlgorithm
		if testCase.En319132 {
			name += "/en319132"
		}

		t.Run(name, func(t *testing.T) {
			namespace := xadesBuilderKATNamespace(t, testCase.XadesNamespaceUri)
			builder, root := newXAdESBuilderKATBuilder(namespace, testCase.En319132)
			digestAlgorithm := enumerations.DigestAlgorithm(testCase.DigestAlgorithm)

			err := xadesBuilderKATRun(builder, root, testCase.Operation, certificate, digestAlgorithm)

			if testCase.Unsupported != "" {
				// Upstream raises UnsupportedOperationException from the XAdES 1.1.1/1.2.2
				// element enums; the Go port must refuse the same combinations.
				if err == nil {
					t.Fatalf("expected the port to refuse this combination, upstream reports %q",
						testCase.Unsupported)
				}
				return
			}
			if err != nil {
				t.Fatalf("build fragment: %v", err)
			}

			serialized, err := xmlutils.DomUtilsSerializeNode(builder.DocumentDom)
			if err != nil {
				t.Fatalf("serialize: %v", err)
			}
			if string(serialized) != testCase.XML {
				t.Errorf("DOM differs from the Java oracle\n  got:  %s\n  want: %s",
					serialized, testCase.XML)
			}
		})
	}
}

// xadesBuilderKATRun dispatches one oracle row onto the method under test.
func xadesBuilderKATRun(builder *xadesBuilderKATBuilder, root *xmldom.Node, operation string,
	certificate *model.CertificateToken, digestAlgorithm enumerations.DigestAlgorithm) (err error) {
	// The DEF chunk's XAdES111Element/XAdES122Element answer an unsupported element the way
	// Java does, with an unchecked exception; recovering here lets the "unsupported" rows of the
	// oracle be asserted as errors regardless of which of the two shapes it lands in.
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = recoveredErr
			} else {
				err = errXAdESBuilderKATUnsupported
			}
		}
	}()

	switch operation {
	case "incorporateCertDigest":
		return builder.IncorporateCertDigest(root, digestAlgorithm, certificate)
	case "incorporateCert":
		_, err := builder.IncorporateCert(root, certificate, digestAlgorithm)
		return err
	case "incorporateSPDocSpecification":
		spDocSpecification := model.NewSpDocSpecification()
		spDocSpecification.SetId("1.2.3.4.5")
		spDocSpecification.SetQualifier(enumerations.ObjectIdentifierQualifierOIDAsURN)
		spDocSpecification.SetDescription("DSS Go port oracle policy")
		spDocSpecification.SetDocumentationReferences("http://nowina.lu/ref1", "http://nowina.lu/ref2")
		return builder.IncorporateSPDocSpecification(root, spDocSpecification)
	}
	return errXAdESBuilderKATUnknownOperation
}

var (
	errXAdESBuilderKATUnsupported      = xadesBuilderKATError("the XAdES element is not supported by this namespace")
	errXAdESBuilderKATUnknownOperation = xadesBuilderKATError("unknown oracle operation")
)

// xadesBuilderKATError is a minimal error type, kept local to this test file.
type xadesBuilderKATError string

func (e xadesBuilderKATError) Error() string { return string(e) }

// TestXAdESBuilderToXmlIdentifierMatchesJavaOracle pins the deterministic XML Id every
// timestamp and validation-data element in this package is labelled with.
func TestXAdESBuilderToXmlIdentifierMatchesJavaOracle(t *testing.T) {
	certificate := xadesBuilderKATCertificate(t)

	for _, testCase := range xadesBuilderKATCases(t) {
		if testCase.Operation != "toXmlIdentifier" {
			continue
		}
		if got := certificate.DSSID().AsXmlID(); got != testCase.AsXmlID {
			t.Fatalf("Identifier.AsXmlID() = %q, oracle says %q", got, testCase.AsXmlID)
		}

		builder, _ := newXAdESBuilderKATBuilder(definition.XAdESNamespaceXAdES132, false)
		got, err := builder.ToXmlIdentifier(certificate.DSSID())
		if err != nil {
			t.Fatalf("ToXmlIdentifier: %v", err)
		}
		if got != testCase.XmlIdentifier {
			t.Errorf("ToXmlIdentifier() = %q, oracle says %q", got, testCase.XmlIdentifier)
		}
		return
	}
	t.Fatal("the oracle carries no toXmlIdentifier row")
}

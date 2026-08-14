// Generates ../expressions.txt: the deduplicated inventory of every XPath expression that
// DSS 6.5.RC1 can hand to its XPath engine.
//
// "Its XPath engine" is eu.europa.esig.dss.xml.utils.xpath.JavaXmlXPathQueryExecutor, which
// compiles a string with javax.xml.xpath and evaluates it as a NODESET. Everything that
// reaches XPath.compile() in DSS goes through it or through the deprecated
// DomUtils.createXPathExpression, which is the same call with the same namespace context.
// internal/xpath10 replaces exactly that, so this inventory bounds exactly that.
//
// Five sources are walked, and each becomes a section of the output:
//
//   A. Reflection over every *Path class. Every static or instance member holding an
//      XPathQuery is asked for getQueryString(). That is the entire output of AbstractPath /
//      XPathQueryBuilder as it appears in compiled constants, which is where all but a
//      handful of DSS's expressions come from.
//   B. The XPathQueryBuilder shapes that only exist once a runtime value is bound - the
//      identifier predicate of XPathQueryIdentifierParameter and the attribute-value
//      predicate of XPathQueryAttributeParameter. Reproduced here against Id values that
//      really occur in ../fixtures, plus one that occurs in none.
//   C. The XPathQueryBuilder shapes exercised by AbstractTestXPathQueryExecutor and by the
//      runtime call sites in DSSXMLUtils, XAdESLevelBaselineT and XAdESCanonicalizationTest:
//      the "*" any-item, the "@name" attribute item and the unprefixed element name.
//   D. Signature-placement expressions. XAdESSignatureParameters.setXPathLocationString feeds
//      XPathPlacementSignatureBuilder, which calls getXPathStringExecutor().getNodeList(),
//      so a caller-written string reaches the same engine. These are the ones upstream
//      writes in its own tests.
//   E. The one expression DomUtils builds by string concatenation, in isNotEmpty().
//
// DELIBERATELY EXCLUDED, and why:
//
//   * XML-DSig transform expressions - XPathTransform, XPath2FilterTransform, and the
//     "not(ancestor-or-self::ds:Signature)" of the enveloped-signature transform. Those are
//     evaluated by Santuario inside the Reference/Transform pipeline, never by
//     JavaXmlXPathQueryExecutor, and XPath Filter 2.0 is a different evaluation model
//     (intersect/subtract over node-sets) on top of a different engine. They belong to the
//     xmldsig package of phase 4b.
//   * XSLT "select" attributes in XSLTTransform test stylesheets, evaluated by the XSLT
//     processor.
//   * Resource paths, URIs and format strings that merely start with "/" or "//".
//
// Build (the same maven-free pattern as ../../../cmscore/testdata/gen):
//
//   U=/home/user/dss-upstream
//   CP=$(find $U -maxdepth 3 -type d -name classes -path '*/target/*' | tr '\n' ':')$(
//       find ~/.m2/repository -name '*.jar' ! -name '*-sources.jar' | tr '\n' ':')
//   javac -nowarn -cp "$CP" -d /tmp/xpinv $(find \
//       $U/dss-xades/src/main/java/eu/europa/esig/dss/xades/definition \
//       $U/dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition \
//       $U/dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/definition \
//       $U/dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/definition -name '*.java')
//   javac -nowarn -cp "$CP:/tmp/xpinv" -d /tmp/xpinv InventoryDump.java
//   java -cp "$CP:/tmp/xpinv" InventoryDump <testdata directory>
//
import eu.europa.esig.dss.xml.common.definition.DSSAttribute;
import eu.europa.esig.dss.xml.common.definition.DSSElement;
import eu.europa.esig.dss.xml.common.definition.DSSNamespace;
import eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigAttribute;
import eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigElement;
import eu.europa.esig.dss.xml.common.xpath.XPathQuery;
import eu.europa.esig.dss.xml.common.xpath.XPathQueryBuilder;

import java.io.BufferedWriter;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;
import java.util.TreeSet;

public class InventoryDump {

    /** Every class that declares XPathQuery constants, across the six modules that have any. */
    private static final String[] PATH_CLASSES = {
        "eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigPath",
        "eu.europa.esig.dss.xades.definition.XAdESPath",
        "eu.europa.esig.dss.xades.definition.xades111.XAdES111Path",
        "eu.europa.esig.dss.xades.definition.xades122.XAdES122Path",
        "eu.europa.esig.dss.xades.definition.xades132.XAdES132Path",
        "eu.europa.esig.dss.xades.definition.tsl.TrustedListPath",
        "eu.europa.esig.dss.asic.common.definition.ASiCManifestPath",
        "eu.europa.esig.dss.asic.xades.definition.ManifestPath",
        "eu.europa.esig.dss.evidencerecord.xml.definition.XMLERSPath",
        "eu.europa.esig.dss.tsl.definition.mra.MRAPath",
    };

    /**
     * Id values the section B predicates are instantiated with. All but the last two really
     * occur in ../fixtures, so the predicate has something to select and a wrong "or" chain or
     * a wrong string comparison shows up as a missing node; "void" and "no-such-id" occur in
     * none, which pins the empty result.
     */
    private static final String[] ID_VALUES = {
        "id-c57e94ffe3f0aefb2a1d3eac29bccb71", // sample-counter-signed.xml, outer signature
        "id-df21797d87f2abeef3653fda4ee84f2a", // sample-counter-signed.xml, counter signature
        "r-id-1",                              // signedXmlXadesLT.xml
        "ZZ-TL-TEST-001",                      // mra-tl-zz.xml, the TrustServiceStatusList
        "ref-enveloped-signature",             // mra-tl-zz.xml, a ds:Reference
        "HS_signature",                        // Signature-X-AT-1.xml
        "reference-data-0",                    // Signature-X-AT-1.xml
        "world",                               // AbstractTestXPathQueryExecutor.xml, idAttributeTest
        "void",                                // idAttributeTest's negative case
        "no-such-id",                          // present in no fixture
    };

    /**
     * Attribute name/value pairs the section B attribute-value predicate is instantiated with.
     * Each value really occurs in a fixture, except the last two, which pin the empty result.
     */
    private static final String[][] ATTRIBUTE_VALUES = {
        { "URI", "#r-id-1" },
        { "URI", "#xades-id-327e0434284baa0c608061420ccbeb85" },
        { "Algorithm", "http://www.w3.org/2001/04/xmlenc#sha256" },
        { "Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature" },
        { "pos", "nested" },     // AbstractTestXPathQueryExecutor.attributeValueTest
        { "pos", "notnested" },  // the same test's negative case
        { "Id", "no-such-id" },
        // Attribute names that collide with a namespace PREFIX declared in a fixture, paired
        // with the URI that prefix is bound to. A namespace declaration is an attribute node
        // in org.w3c.dom and in xmldom, but XPath 1.0 clause 5.3 puts it on the namespace
        // axis, so the attribute axis must not see it and these must select nothing. Without
        // them the "skip xmlns declarations" rule of the attribute axis is unverified: no
        // other pair has a local name that any xmlns declaration shares.
        { "ds", "http://www.w3.org/2000/09/xmldsig#" },
        { "h", "http://www.w3.org/TR/html4/" },
        { "ext", "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2" },
        { "xmlns", "urn:oasis:names:specification:ubl:schema:xsd:GuaranteeCertificate-2" },
    };

    public static void main(String[] args) throws Exception {
        Path testdata = Paths.get(args[0]);
        try (BufferedWriter w = Files.newBufferedWriter(
                testdata.resolve("expressions.txt"), StandardCharsets.UTF_8)) {
            for (String line : HEADER) {
                w.write(line);
                w.write('\n');
            }
            Set<String> seen = new LinkedHashSet<>();
            section(w, seen, "A. path definitions - XPathQuery.getQueryString()", reflectPathClasses());
            section(w, seen, "B. runtime-parameterised predicates - id and attribute value", runtimePredicates());
            section(w, seen, "C. XPathQueryBuilder shapes without an element chain", builderShapes());
            section(w, seen, "D. signature placement - setXPathLocationString", placementExpressions());
            section(w, seen, "E. DomUtils.isNotEmpty, built by concatenation", isNotEmptyForms());
        }
    }

    /**
     * Writes one section, skipping expressions an earlier section already emitted: the file is
     * a deduplicated inventory, and an expression two sources happen to agree on - "/*" is both
     * a builder shape and a placement string - is one expression to evaluate, not two.
     */
    private static void section(BufferedWriter w, Set<String> seen, String title, Set<String> body)
            throws Exception {
        w.write('\n');
        w.write("## " + title + '\n');
        for (String s : body) {
            if (seen.add(s)) {
                w.write(s);
                w.write('\n');
            }
        }
    }

    // ------------------------------------------------------------------ A

    private static Set<String> reflectPathClasses() {
        Set<String> out = new TreeSet<>();
        for (String className : PATH_CLASSES) {
            Class<?> c;
            try {
                c = Class.forName(className);
            } catch (Throwable t) {
                throw new IllegalStateException("cannot load " + className, t);
            }
            Object instance = null;
            try {
                instance = c.getDeclaredConstructor().newInstance();
            } catch (Throwable ignored) {
                // abstract, or no accessible no-arg constructor: static members only
            }
            for (Class<?> k = c; k != null && k != Object.class; k = k.getSuperclass()) {
                for (Field f : k.getDeclaredFields()) {
                    boolean isStatic = Modifier.isStatic(f.getModifiers());
                    if (!isStatic && instance == null) {
                        continue;
                    }
                    try {
                        f.setAccessible(true);
                        Object v = f.get(isStatic ? null : instance);
                        if (v instanceof XPathQuery) {
                            out.add(((XPathQuery) v).getQueryString());
                        }
                    } catch (Throwable ignored) {
                        // inaccessible under the module system: the getters below still reach it
                    }
                }
            }
            if (instance != null) {
                for (Method m : c.getMethods()) {
                    if (m.getParameterCount() != 0 || !XPathQuery.class.isAssignableFrom(m.getReturnType())) {
                        continue;
                    }
                    try {
                        XPathQuery q = (XPathQuery) m.invoke(instance);
                        if (q != null) {
                            out.add(q.getQueryString());
                        }
                    } catch (Throwable ignored) {
                        // a getter needing state we do not have
                    }
                }
            }
        }
        if (out.isEmpty()) {
            throw new IllegalStateException("no XPathQuery constants found - wrong classpath?");
        }
        return out;
    }

    // ------------------------------------------------------------------ B

    /**
     * The predicate shapes that only appear once a runtime value is bound.
     *
     * XPathUtils.getElementById(node, id) builds the first; DSSXMLUtils line 347 builds it
     * too. PathsTest and XAdESCanonicalizationTest build the rest through
     * XPathQueryBuilder.fromXPathQuery(...).idValue(...) and .attribute(attr, value).
     */
    private static Set<String> runtimePredicates() {
        Set<String> out = new LinkedHashSet<>();
        for (String id : ID_VALUES) {
            // XPathUtils.getElementById(Node, String)
            out.add(XPathQueryBuilder.allFromCurrentPosition().idValue(id).build().getQueryString());
            // the same, anchored at the document
            out.add(XPathQueryBuilder.all().idValue(id).build().getQueryString());
            // fromXPathQuery(SIGNATURE_PATH).idValue(id) and fromXPathQuery(OBJECT_PATH).idValue(id)
            out.add(XPathQueryBuilder.all().element(XMLDSigElement.SIGNATURE).idValue(id).build().getQueryString());
            out.add(XPathQueryBuilder.fromCurrentPosition().element(XMLDSigElement.OBJECT).idValue(id).build().getQueryString());
        }
        for (String[] pair : ATTRIBUTE_VALUES) {
            DSSAttribute attribute = DSSAttribute.fromDefinition(pair[0]);
            // XAdESCanonicalizationTest: all().element(REFERENCE).attribute(URI, uri)
            out.add(XPathQueryBuilder.all().element(XMLDSigElement.REFERENCE)
                    .attribute(attribute, pair[1]).build().getQueryString());
            out.add(XPathQueryBuilder.allFromCurrentPosition().element(XMLDSigElement.TRANSFORM)
                    .attribute(attribute, pair[1]).build().getQueryString());
            // AbstractTestXPathQueryExecutor.attributeValueTest: no element, so the predicate
            // lands on the "*" any-item
            out.add(XPathQueryBuilder.all().attribute(attribute, pair[1]).build().getQueryString());
        }
        // AbstractPath.allNotParent, with the counter-signature element of both XAdES versions
        out.add(XPathQueryBuilder.all().element(XMLDSigElement.SIGNATURE)
                .notChildOf(DSSElement.fromDefinition("CounterSignature", XADES_132)).build().getQueryString());
        return out;
    }

    // ------------------------------------------------------------------ C

    /**
     * Shapes with no element chain, or with an unprefixed element name.
     *
     * XPathQueryBuilder emits XPathQueryAnyItem ("*") when no element is given, and appends an
     * XPathQueryAttributeItem ("@name") as an independent step when an attribute is given
     * without a value. DSSXMLUtils line 526 and XAdESLevelBaselineT line 489 do both;
     * AbstractTestXPathQueryExecutor does all of them over unprefixed names.
     */
    private static Set<String> builderShapes() {
        Set<String> out = new LinkedHashSet<>();
        // XAdESLevelBaselineT: fromCurrentPosition() with nothing set
        out.add(XPathQueryBuilder.fromCurrentPosition().build().getQueryString());
        out.add(XPathQueryBuilder.all().build().getQueryString());
        out.add(XPathQueryBuilder.allFromCurrentPosition().build().getQueryString());
        // DSSXMLUtils: all().attribute(ID) - the "@name" step on the "*" any-item
        for (String name : new String[] { "Id", "id", "ID", "URI", "Algorithm", "Target", "pos", "evil" }) {
            out.add(XPathQueryBuilder.all().attribute(DSSAttribute.fromDefinition(name)).build().getQueryString());
            out.add(XPathQueryBuilder.allFromCurrentPosition().attribute(DSSAttribute.fromDefinition(name)).build().getQueryString());
        }
        // the "@name" step after an element chain
        out.add(XPathQueryBuilder.fromCurrentPosition().elements(XMLDSigElement.SIGNED_INFO, XMLDSigElement.REFERENCE)
                .attribute(XMLDSigAttribute.URI).build().getQueryString());
        // unprefixed element names: DSSElement.fromDefinition(local, null)
        for (String[] chain : new String[][] {
                { "a" }, { "b" }, { "c" }, { "d" }, { "e" },
                { "manifest", "file-entry" }, { "placeOfSignature" },
                { "a", "b" }, { "b", "d" }, { "c", "d" }, { "e", "e" } }) {
            out.add(XPathQueryBuilder.all().elements(unprefixed(chain)).build().getQueryString());
            out.add(XPathQueryBuilder.fromCurrentPosition().elements(unprefixed(chain)).build().getQueryString());
            out.add(XPathQueryBuilder.allFromCurrentPosition().elements(unprefixed(chain)).build().getQueryString());
        }
        // notChildOf over unprefixed names, as AbstractTestXPathQueryExecutor writes it
        for (String parent : new String[] { "a", "b", "c" }) {
            out.add(XPathQueryBuilder.all().element(DSSElement.fromDefinition("d", null))
                    .notChildOf(DSSElement.fromDefinition(parent, null)).build().getQueryString());
        }
        // a prefixed element whose namespace is registered under a prefix the document does
        // not use: XAdESLevelBEnvelopedWithXPathPlacementAfterTest's h:tr
        out.add(XPathQueryBuilder.all().element(DSSElement.fromDefinition("tr", HTML4)).build().getQueryString());
        // XAdESEnumsTest / ASiCEnumsTest walk an XSD with xsd:element and xsd:attribute
        out.add(XPathQueryBuilder.all().element(DSSElement.fromDefinition("element", XSD)).build().getQueryString());
        out.add(XPathQueryBuilder.all().element(DSSElement.fromDefinition("attribute", XSD)).build().getQueryString());
        out.add(XPathQueryBuilder.all().element(DSSElement.fromDefinition("element", XSD))
                .attribute(DSSAttribute.fromDefinition("name"), "SignedProperties").build().getQueryString());
        return out;
    }

    private static DSSElement[] unprefixed(String[] locals) {
        List<DSSElement> out = new ArrayList<>();
        for (String local : locals) {
            out.add(DSSElement.fromDefinition(local, null));
        }
        return out.toArray(new DSSElement[0]);
    }

    // ------------------------------------------------------------------ D

    /**
     * Caller-written placement expressions, verbatim from the upstream tests that set them.
     * XPathPlacementSignatureBuilder hands these straight to the string executor.
     */
    private static Set<String> placementExpressions() {
        Set<String> out = new LinkedHashSet<>();
        out.add("/*");                                       // ...AfterRootNodeTest
        out.add("//*[local-name() = 'tr']");                 // ...PlacementAfter/FirstChild/NoneTest
        out.add("//*[local-name() = 'ElementNotExists']");   // ...InvalidXPathPlacementTest
        out.add("//placeOfSignature");                       // XAdESLevelBEnvelopedWithXPathTest
        // XAdESSecondSignatureToParentNodeWithEnvelopedTransformTest et al: String.format
        for (String id : ID_VALUES) {
            out.add("//*[@Id='" + id + "']");
        }
        // XAdESXPathPlacementWithEmptyNamespaceTest, split across two source lines upstream
        out.add("//*[local-name()='GuaranteeCertificate']//*[local-name()='UBLExtensions']"
                + "//*[local-name()='UBLExtension']//*[local-name()='ExtensionContent']");
        // XAdESLevelBInternallyDetachedWithXPathLocationTest: "//" + containerNodeName
        out.add("//internally-detached");
        return out;
    }

    // ------------------------------------------------------------------ E

    /**
     * DomUtils.isNotEmpty(node, s) evaluates s + "/child::node()[not(self::text())]".
     *
     * The method is deprecated in 6.5 and has no caller left in the upstream tree, but it is
     * public API of dss-xml-utils and the concatenation is unconditional, so the form is
     * reachable and is the only place child::, node(), self:: and text() enter the inventory.
     */
    private static Set<String> isNotEmptyForms() {
        Set<String> out = new LinkedHashSet<>();
        String suffix = "/child::node()[not(self::text())]";
        out.add("." + suffix);
        out.add("./ds:SignedInfo" + suffix);
        out.add("./ds:Object/ds:Manifest" + suffix);
        out.add("//ds:Signature" + suffix);
        out.add(".//ds:Reference" + suffix);
        out.add("/*" + suffix);
        return out;
    }

    // ------------------------------------------------------------------ constants

    private static final DSSNamespace XADES_132 =
            new DSSNamespace("http://uri.etsi.org/01903/v1.3.2#", "xades132");
    private static final DSSNamespace HTML4 =
            new DSSNamespace("http://www.w3.org/TR/html4/", "h");
    private static final DSSNamespace XSD =
            new DSSNamespace("http://www.w3.org/2001/XMLSchema", "xsd");

    private static final String[] HEADER = {
        "# XPath expression inventory for internal/xpath10. GENERATED by gen/InventoryDump.java;",
        "# do not edit by hand. Build and run it as documented in that file's header.",
        "#",
        "# Every expression DSS 6.5.RC1 can hand to javax.xml.xpath through",
        "# JavaXmlXPathQueryExecutor (or the deprecated DomUtils.createXPathExpression, which is",
        "# the same call), deduplicated. internal/xpath10 implements exactly the constructs these",
        "# expressions use and rejects every other XPath 1.0 construct by name.",
        "#",
        "# Sections:",
        "#   A  reflection over every *Path class: XPathQuery.getQueryString()",
        "#   B  XPathQueryBuilder predicates that need a runtime value (id, attribute value)",
        "#   C  XPathQueryBuilder shapes with no element chain, or an unprefixed element name",
        "#   D  caller-written placement strings, XAdESSignatureParameters.setXPathLocationString",
        "#   E  the concatenated form DomUtils.isNotEmpty builds",
        "#",
        "# Excluded on purpose: XML-DSig transform expressions (XPathTransform,",
        "# XPath2FilterTransform, enveloped-signature), which Santuario evaluates inside the",
        "# Reference pipeline and not this engine, and XSLT select attributes. See the header of",
        "# gen/InventoryDump.java for the full reasoning.",
        "#",
        "# Blank lines, '#' comments and '## ' section headings are ignored by both readers:",
        "# gen/XPathOracle.java (javax.xml.xpath, produces kat.txt) and TestKnownAnswers (Go).",
        "# Namespace context: namespaces.txt.",
        "#",
        "# This file is the inventory and stays exactly that. The conversion-rule corners the",
        "# accepted grammar can reach but no upstream expression uses live in semantics.txt,",
        "# which the same oracle turns into semantics-kat.txt.",
    };
}

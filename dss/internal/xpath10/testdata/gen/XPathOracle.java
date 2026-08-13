// Generates ../kat.txt: the known-answer vectors for internal/xpath10.
//
// Every expression in ../expressions.txt is evaluated with javax.xml.xpath - the same
// XPathFactory.newInstance() that DSS's JavaXmlXPathQueryExecutor uses, so the answers come
// from the JDK's Xalan and not from a reimplementation - against every context node of every
// fixture in ../fixtures, as a NODESET, which is the only way DSS evaluates anything.
//
// A context node is the document node or any element node, in document order. That is exactly
// the set DSS itself can pass: XPathUtils.getNodeList is called either with a Document (the
// "//" and "/" forms of AllXPathQuery) or with an element (the "./" and ".//" forms of
// FromCurrentPositionXPathQuery), and evaluating from EVERY element rather than from the few
// upstream happens to use is what makes a disagreement anywhere in the subset fail the test.
//
// The namespace context is ../namespaces.txt wrapped in the semantics of
// eu.europa.esig.dss.xml.utils.NamespaceContextMap: an unregistered prefix resolves to
// XMLConstants.NULL_NS_URI rather than raising, so a stray prefix silently means "no
// namespace". internal/xpath10 has to reproduce that, so it is pinned here.
//
// Only contexts with a non-empty result are written; the "X" header line carries the number of
// such contexts and the "F" header line the number of contexts tried, so a Go side that
// enumerates a different context set, or drops a match, cannot produce this file.
//
// Node paths are child indices, not names, so a tree-shape difference cannot hide behind
// matching names:
//
//   document        /
//   any child node  <parent>/<1-based index among ALL of the parent's child nodes>
//   attribute       <owner element>/@<qualified name>
//
// Each fixture is preceded by its node table - one "N <path> <label>" line for every node
// except attributes, in document order, where <label> is the qualified name for an element,
// "#text", "#cdata", "#comment" or "?<target>" otherwise. The table names the index paths for
// a human reading a diff, and it is itself a known answer: reproducing it requires splitting
// text nodes, CDATA sections, comments and processing instructions exactly as Xerces does.
//
// TWO ANSWER SETS are generated, selected by the optional second argument:
//
//   (none)      expressions.txt + fixtures/           -> kat.txt
//               the upstream inventory: what DSS can hand the engine, and nothing else.
//   semantics   semantics.txt   + semantics-fixtures/ -> semantics-kat.txt
//               the XPath 1.0 conversion rules (clause 3.4) and string-values (clause 5) that
//               the accepted grammar can reach but no upstream expression happens to use.
//               They are pinned by this same oracle rather than asserted from memory.
//
// Run it with OpenJDK 21 (no third-party jars needed):
//
//   javac -d /tmp/xpathoracle XPathOracle.java
//   java -cp /tmp/xpathoracle XPathOracle <testdata directory>
//   java -cp /tmp/xpathoracle XPathOracle <testdata directory> semantics
//
import org.w3c.dom.Attr;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import javax.xml.XMLConstants;
import javax.xml.namespace.NamespaceContext;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.xpath.XPath;
import javax.xml.xpath.XPathConstants;
import javax.xml.xpath.XPathExpression;
import javax.xml.xpath.XPathFactory;
import java.io.BufferedWriter;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class XPathOracle {

    public static void main(String[] args) throws Exception {
        Path testdata = Paths.get(args[0]);
        // Optional second argument: the base name of an expression/fixture/answer set other
        // than the default "expressions" + "fixtures" + "kat". Used for the "semantics" set,
        // which pins the XPath 1.0 conversion rules the accepted grammar can reach but no
        // upstream expression happens to use.
        String set = args.length > 1 ? args[1] : null;
        String exprFile = set == null ? "expressions.txt" : set + ".txt";
        String fixtureDir = set == null ? "fixtures" : set + "-fixtures";
        String outFile = set == null ? "kat.txt" : set + "-kat.txt";

        List<String> expressions = readExpressions(testdata.resolve(exprFile));
        Map<String, String> namespaces = readNamespaces(testdata.resolve("namespaces.txt"));
        List<Path> fixtures = readFixtures(testdata.resolve(fixtureDir));

        XPath xpath = XPathFactory.newInstance().newXPath();
        xpath.setNamespaceContext(new PrefixMap(namespaces));

        List<XPathExpression> compiled = new ArrayList<>();
        for (String e : expressions) {
            compiled.add(xpath.compile(e));
        }

        DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
        dbf.setNamespaceAware(true);
        dbf.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
        dbf.setExpandEntityReferences(false);
        dbf.setCoalescing(false);
        dbf.setIgnoringComments(false);
        dbf.setIgnoringElementContentWhitespace(false);

        try (BufferedWriter w = Files.newBufferedWriter(
                testdata.resolve(outFile), StandardCharsets.UTF_8)) {
            w.write("# Known-answer vectors produced by gen/XPathOracle.java "
                    + "(javax.xml.xpath, OpenJDK " + System.getProperty("java.version") + ").\n");
            w.write("# Do not edit by hand. Format: see the header of gen/XPathOracle.java.\n");

            for (Path fixture : fixtures) {
                DocumentBuilder db = dbf.newDocumentBuilder();
                Document doc = db.parse(fixture.toFile());

                List<Node> all = new ArrayList<>();
                collectNodes(doc, all);
                List<Node> contexts = new ArrayList<>();
                contexts.add(doc);
                for (Node n : all) {
                    if (n.getNodeType() == Node.ELEMENT_NODE) {
                        contexts.add(n);
                    }
                }

                String name = testdata.resolve(fixtureDir).relativize(fixture).toString()
                        .replace('\\', '/');
                w.write("F " + name + " nodes=" + all.size() + " contexts=" + contexts.size() + "\n");
                for (Node n : all) {
                    w.write("N " + path(n) + " " + label(n) + "\n");
                }

                for (int i = 0; i < expressions.size(); i++) {
                    StringBuilder body = new StringBuilder();
                    int hits = 0;
                    for (Node ctx : contexts) {
                        String line = evaluate(compiled.get(i), ctx);
                        if (line == null) {
                            continue;
                        }
                        hits++;
                        body.append("= ").append(path(ctx)).append(" -> ").append(line).append('\n');
                    }
                    w.write("X hits=" + hits + " " + expressions.get(i) + "\n");
                    w.write(body.toString());
                }
            }
        }
        System.err.println("fixtures=" + fixtures.size() + " expressions=" + expressions.size());
    }

    /**
     * Evaluates one expression against one context node and renders the matched nodes in the
     * order the NodeList reports them, which for a NODESET result is document order. Returns
     * null when the node-set is empty, so nothing is written for that context.
     */
    private static String evaluate(XPathExpression x, Node ctx) throws Exception {
        NodeList nodes = (NodeList) x.evaluate(ctx, XPathConstants.NODESET);
        if (nodes == null || nodes.getLength() == 0) {
            return null;
        }
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < nodes.getLength(); i++) {
            if (i > 0) {
                sb.append(' ');
            }
            sb.append(path(nodes.item(i)));
        }
        return sb.toString();
    }

    private static void collectNodes(Node parent, List<Node> out) {
        NodeList children = parent.getChildNodes();
        for (int i = 0; i < children.getLength(); i++) {
            Node child = children.item(i);
            out.add(child);
            collectNodes(child, out);
        }
    }

    private static String path(Node n) {
        switch (n.getNodeType()) {
            case Node.DOCUMENT_NODE:
                return "/";
            case Node.ATTRIBUTE_NODE: {
                Element owner = ((Attr) n).getOwnerElement();
                return path(owner) + "/@" + n.getNodeName();
            }
            default: {
                Node parent = n.getParentNode();
                String prefix = parent.getNodeType() == Node.DOCUMENT_NODE ? "" : path(parent);
                return prefix + "/" + childIndex(parent, n);
            }
        }
    }

    private static int childIndex(Node parent, Node child) {
        NodeList children = parent.getChildNodes();
        for (int i = 0; i < children.getLength(); i++) {
            if (children.item(i) == child) {
                return i + 1;
            }
        }
        throw new IllegalStateException("node is not a child of its own parent");
    }

    private static String label(Node n) {
        switch (n.getNodeType()) {
            case Node.ELEMENT_NODE:
                return n.getNodeName();
            case Node.TEXT_NODE:
                return "#text";
            case Node.CDATA_SECTION_NODE:
                return "#cdata";
            case Node.COMMENT_NODE:
                return "#comment";
            case Node.PROCESSING_INSTRUCTION_NODE:
                return "?" + n.getNodeName();
            default:
                throw new IllegalStateException("unexpected node type " + n.getNodeType());
        }
    }

    // ------------------------------------------------------------------ inventory files

    private static List<String> readExpressions(Path file) throws IOException {
        List<String> out = new ArrayList<>();
        for (String line : Files.readAllLines(file, StandardCharsets.UTF_8)) {
            if (line.isBlank() || line.startsWith("#")) {
                continue;
            }
            out.add(line);
        }
        return out;
    }

    private static Map<String, String> readNamespaces(Path file) throws IOException {
        Map<String, String> out = new LinkedHashMap<>();
        for (String line : Files.readAllLines(file, StandardCharsets.UTF_8)) {
            if (line.isBlank() || line.startsWith("#")) {
                continue;
            }
            int tab = line.indexOf('\t');
            out.put(line.substring(0, tab), line.substring(tab + 1));
        }
        return out;
    }

    private static List<Path> readFixtures(Path dir) throws IOException {
        List<Path> out = new ArrayList<>();
        Files.walkFileTree(dir, new SimpleFileVisitor<Path>() {
            @Override
            public FileVisitResult visitFile(Path f, BasicFileAttributes attrs) {
                if (f.getFileName().toString().endsWith(".xml")) {
                    out.add(f);
                }
                return FileVisitResult.CONTINUE;
            }
        });
        out.sort(Path::compareTo);
        return out;
    }

    /** eu.europa.esig.dss.xml.utils.NamespaceContextMap, reduced to what an XPath needs. */
    private static final class PrefixMap implements NamespaceContext {
        private final Map<String, String> prefixMap;

        PrefixMap(Map<String, String> prefixMap) {
            this.prefixMap = prefixMap;
        }

        @Override
        public String getNamespaceURI(String prefix) {
            if (prefix == null) {
                throw new IllegalArgumentException("null");
            }
            String uri = prefixMap.get(prefix);
            return uri == null ? XMLConstants.NULL_NS_URI : uri;
        }

        @Override
        public String getPrefix(String namespaceURI) {
            if (namespaceURI == null) {
                throw new IllegalArgumentException("null");
            }
            for (Map.Entry<String, String> e : prefixMap.entrySet()) {
                if (e.getValue().equals(namespaceURI)) {
                    return e.getKey();
                }
            }
            return null;
        }

        @Override
        public Iterator<String> getPrefixes(String namespaceURI) {
            List<String> out = new ArrayList<>();
            for (Map.Entry<String, String> e : prefixMap.entrySet()) {
                if (e.getValue().equals(namespaceURI)) {
                    out.add(e.getKey());
                }
            }
            return out.iterator();
        }
    }
}

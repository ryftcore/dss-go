// Generates testdata/serialize/goldens.txt: what OpenJDK's identity Transformer - the one
// DomUtils.serializeNode and DomUtils.getNodeBytes run - actually emits for every case in
// testdata/serialize/corpus.txt.
//
// serializeNode(Node), getNodeBytes(Node) and the two factory builders below are copied
// VERBATIM from DSS 6.5.RC1 (dss-xml-utils DomUtils, dss-xml-common
// TransformerFactoryBuilder / DocumentBuilderFactoryBuilder) so that the transformer
// configuration - which output properties are set, and which deliberately are not - is the
// production one. Nothing here may be "tidied": the omission of OMIT_XML_DECLARATION is the
// single most consequential fact this oracle records.
//
// Run it with OpenJDK 21 and commons-codec 1.18.0 from the local Maven repository (the
// version dss-utils-apache-commons pins, whose Base64.decodeBase64 is Utils.fromBase64):
//
//   CC=$HOME/.m2/repository/commons-codec/commons-codec/1.18.0/commons-codec-1.18.0.jar
//   javac -cp "$CC" -d /tmp/serOracle SerializeOracle.java
//   java -cp "$CC:/tmp/serOracle" SerializeOracle ../serialize/corpus.txt ../serialize/goldens.txt
//
// Every value in goldens.txt is base64 of the exact bytes Java produced, so the file is
// insensitive to the console charset and to trailing-whitespace mangling.
import org.w3c.dom.*;

import javax.xml.XMLConstants;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.transform.OutputKeys;
import javax.xml.transform.Result;
import javax.xml.transform.Source;
import javax.xml.transform.Transformer;
import javax.xml.transform.TransformerFactory;
import javax.xml.transform.dom.DOMSource;
import javax.xml.transform.stream.StreamResult;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

public class SerializeOracle {

    // ------------------------------------------------------------------ DSS, verbatim

    /** dss-xml-common DocumentBuilderFactoryBuilder's secure constructor, inlined. */
    static DocumentBuilderFactory secureDocumentBuilderFactory() throws Exception {
        DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
        factory.setFeature("http://xml.org/sax/features/namespaces", true);
        factory.setFeature("http://apache.org/xml/features/dom/create-entity-ref-nodes", true);
        factory.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
        factory.setFeature("http://xml.org/sax/features/external-general-entities", false);
        factory.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
        factory.setFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false);
        factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_DTD, "");
        factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_SCHEMA, "");
        return factory;
    }

    /** DomUtils.getSecureTransformer() over TransformerFactoryBuilder's secure constructor. */
    static final String TRANSFORMER_METHOD_VALUE = "xml";

    static Transformer getSecureTransformer() throws Exception {
        TransformerFactory transformerFactory = TransformerFactory.newInstance();
        transformerFactory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
        transformerFactory.setAttribute(XMLConstants.ACCESS_EXTERNAL_DTD, "");
        transformerFactory.setAttribute(XMLConstants.ACCESS_EXTERNAL_STYLESHEET, "");
        Transformer transformer = transformerFactory.newTransformer();
        transformer.setOutputProperty(OutputKeys.METHOD, TRANSFORMER_METHOD_VALUE);
        return transformer;
    }

    /** DomUtils.serializeNode(Node) - note: no OMIT_XML_DECLARATION anywhere. */
    static byte[] serializeNode(final Node xmlNode) throws Exception {
        try (ByteArrayOutputStream bos = new ByteArrayOutputStream()) {
            Transformer transformer = getSecureTransformer();
            Document document;
            if (Node.DOCUMENT_NODE == xmlNode.getNodeType()) {
                document = (Document) xmlNode;
            } else {
                document = xmlNode.getOwnerDocument();
            }

            if (document != null) {
                String xmlEncoding = document.getXmlEncoding();
                if (isStringNotBlank(xmlEncoding)) {
                    transformer.setOutputProperty(OutputKeys.ENCODING, xmlEncoding);
                }
            }

            StreamResult result = new StreamResult(bos);
            Source source = new DOMSource(xmlNode);
            transformer.transform(source, result);

            return bos.toByteArray();
        }
    }

    /** DomUtils.getNodeBytes(Node). */
    static byte[] getNodeBytes(Node node) throws Exception {
        switch (node.getNodeType()) {
            case Node.ELEMENT_NODE:
            case Node.DOCUMENT_NODE:
            case Node.COMMENT_NODE:
                byte[] bytes = serializeNode(node);
                String str = new String(bytes);
                // TODO: better
                // remove <?xml version="1.0" encoding="UTF-8"?>
                if (str.startsWith("<?")) {
                    str = str.substring(str.indexOf("?>") + 2);
                }
                return str.getBytes();

            case Node.TEXT_NODE:
                String textContent = node.getTextContent();
                // Use try-catch for performance purposes
                try {
                    // Utils.fromBase64 -> ApacheCommonsUtils.fromBase64 -> Base64.decodeBase64
                    return org.apache.commons.codec.binary.Base64.decodeBase64(textContent);
                } catch (Exception e) {
                    return textContent.getBytes();
                }

            default:
                return null;
        }
    }

    /** eu.europa.esig.dss.utils.Utils#isStringNotBlank. */
    static boolean isStringNotBlank(String s) {
        return s != null && !s.trim().isEmpty();
    }

    // -------------------------------------------------------------------- corpus driver

    public static void main(String[] args) throws Exception {
        Path corpus = Paths.get(args[0]);
        Path out = Paths.get(args[1]);

        StringBuilder sb = new StringBuilder();
        sb.append("# GENERATED by gen/SerializeOracle.java on OpenJDK ")
          .append(System.getProperty("java.version"))
          .append(" - do not edit by hand.\n")
          .append("# One record per corpus case. Values are base64 of the exact bytes Java produced;\n")
          .append("# NULL is a null return, and ERROR <text> is the exception DSS would propagate.\n")
          .append("#\n")
          .append("#   tree      base64 of the structural dump of the parsed/built DOM (see below)\n")
          .append("#   node      base64 of the dump line of the selected node\n")
          .append("#   serialize base64 of DomUtils.serializeNode(node)\n")
          .append("#   bytes     base64 of DomUtils.getNodeBytes(node)\n")
          .append("#\n")
          .append("# The tree dump is one line per node in document order: \"<path> <kind> <nodeName>\"\n")
          .append("# followed, for elements, by the attribute node names in NamedNodeMap order.\n")
          .append("# It pins that the Go DOM being serialized is the same DOM Xerces built.\n");

        int count = 0;
        for (Record r : readCorpus(corpus)) {
            count++;
            sb.append("\ncase ").append(r.name).append('\n');
            String serialize;
            String bytes;
            String tree;
            String node;
            try {
                Document doc = r.build(secureDocumentBuilderFactory().newDocumentBuilder());
                Node target = select(doc, r.sel);
                tree = b64(dumpTree(doc).getBytes(StandardCharsets.UTF_8));
                node = b64(dumpNode("@", target).getBytes(StandardCharsets.UTF_8));
                serialize = call(() -> serializeNode(target));
                bytes = call(() -> getNodeBytes(target));
            } catch (Throwable t) {
                tree = node = serialize = bytes = "ERROR " + t.getClass().getName() + ": " + oneLine(t.getMessage());
            }
            sb.append("tree ").append(tree).append('\n');
            sb.append("node ").append(node).append('\n');
            sb.append("serialize ").append(serialize).append('\n');
            sb.append("bytes ").append(bytes).append('\n');
        }
        Files.write(out, sb.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("wrote " + count + " cases to " + out);
    }

    interface Call {
        byte[] run() throws Exception;
    }

    static String call(Call c) {
        try {
            byte[] b = c.run();
            return b == null ? "NULL" : b64(b);
        } catch (Throwable t) {
            Throwable root = t;
            while (root.getCause() != null && root.getCause() != root) root = root.getCause();
            return "ERROR " + root.getClass().getName() + ": " + oneLine(root.getMessage());
        }
    }

    static String oneLine(String s) {
        return s == null ? "" : s.replace('\n', ' ').replace('\r', ' ');
    }

    static String b64(byte[] b) {
        return Base64.getEncoder().encodeToString(b);
    }

    // ------------------------------------------------------------------- DOM structure

    static String dumpTree(Document doc) {
        StringBuilder sb = new StringBuilder();
        dumpTree(sb, "", doc);
        return sb.toString();
    }

    static void dumpTree(StringBuilder sb, String path, Node n) {
        sb.append(dumpNode(path.isEmpty() ? "." : path, n)).append('\n');
        NodeList kids = n.getChildNodes();
        for (int i = 0; i < kids.getLength(); i++) {
            dumpTree(sb, path.isEmpty() ? Integer.toString(i) : path + "." + i, kids.item(i));
        }
    }

    static String dumpNode(String path, Node n) {
        StringBuilder sb = new StringBuilder();
        sb.append(path).append(' ').append(kind(n)).append(' ').append(n.getNodeName());
        if (n.getNodeType() == Node.ELEMENT_NODE) {
            NamedNodeMap m = n.getAttributes();
            for (int i = 0; i < m.getLength(); i++) sb.append(' ').append(m.item(i).getNodeName());
        }
        return sb.toString();
    }

    static String kind(Node n) {
        switch (n.getNodeType()) {
            case Node.DOCUMENT_NODE: return "document";
            case Node.ELEMENT_NODE: return "element";
            case Node.ATTRIBUTE_NODE: return "attribute";
            case Node.TEXT_NODE: return "text";
            case Node.CDATA_SECTION_NODE: return "cdata";
            case Node.COMMENT_NODE: return "comment";
            case Node.PROCESSING_INSTRUCTION_NODE: return "pi";
            default: return "kind" + n.getNodeType();
        }
    }

    static Node select(Document doc, String sel) {
        if (sel == null || sel.isEmpty() || sel.equals(".")) return doc;
        String attr = null;
        int at = sel.indexOf('@');
        if (at >= 0) {
            attr = sel.substring(at + 1);
            sel = sel.substring(0, at);
        }
        Node n = doc;
        if (!sel.isEmpty()) {
            for (String part : sel.split("\\.")) {
                n = n.getChildNodes().item(Integer.parseInt(part));
                if (n == null) throw new IllegalStateException("selector " + sel + " does not resolve");
            }
        }
        if (attr != null) {
            Node a = n.getAttributes().getNamedItem(attr);
            if (a == null) throw new IllegalStateException("no attribute " + attr);
            return a;
        }
        return n;
    }

    // ------------------------------------------------------------------ corpus records

    static class Record {
        String name;
        String parse;      // escaped text source, or null
        byte[] parseBytes; // raw source bytes, or null
        String build;      // DSL, or null
        String sel = ".";
        String divergence; // why this case is known NOT to agree with Java, or null

        Document build(DocumentBuilder db) throws Exception {
            Document doc;
            if (parseBytes != null) {
                doc = db.parse(new ByteArrayInputStream(parseBytes));
            } else if (parse != null) {
                doc = db.parse(new ByteArrayInputStream(parse.getBytes(StandardCharsets.UTF_8)));
            } else {
                doc = db.newDocument();
            }
            if (build != null) {
                // A "build" alongside a "parse" is DOM surgery on the parsed document -
                // what XAdES does to a document before serializing it - and starts at the
                // document element rather than at the document.
                runDSL(doc, build, doc.getDocumentElement() != null ? doc.getDocumentElement() : doc);
            }
            return doc;
        }
    }

    static void runDSL(Document doc, String script, Node start) {
        Node cur = start;
        for (String op : script.split(";")) {
            op = op.trim();
            if (op.isEmpty() || op.equals("~")) continue;
            String[] tok = op.split(" ");
            switch (tok[0]) {
                case "e": {
                    Element e = doc.createElementNS(nullable(unescape(tok[1])), unescape(tok[2]));
                    cur.appendChild(e);
                    cur = e;
                    break;
                }
                case "a":
                    ((Element) cur).setAttributeNS(nullable(unescape(tok[1])), unescape(tok[2]), unescape(tok[3]));
                    break;
                case "ap":
                    ((Element) cur).setAttribute(unescape(tok[1]), unescape(tok[2]));
                    break;
                case "t":
                    cur.appendChild(doc.createTextNode(unescape(tok[1])));
                    break;
                case "c":
                    cur.appendChild(doc.createComment(unescape(tok[1])));
                    break;
                case "d":
                    cur.appendChild(doc.createCDATASection(unescape(tok[1])));
                    break;
                case "p":
                    cur.appendChild(doc.createProcessingInstruction(unescape(tok[1]), unescape(tok[2])));
                    break;
                case "/":
                    cur = cur.getParentNode();
                    break;
                default:
                    throw new IllegalArgumentException("unknown DSL op " + tok[0]);
            }
        }
    }

    static String nullable(String s) {
        return "-".equals(s) ? null : s;
    }

    static List<Record> readCorpus(Path p) throws Exception {
        List<Record> out = new ArrayList<>();
        Record cur = null;
        for (String line : Files.readAllLines(p, StandardCharsets.UTF_8)) {
            String t = line.trim();
            if (t.isEmpty() || t.startsWith("#")) continue;
            int sp = t.indexOf(' ');
            String key = sp < 0 ? t : t.substring(0, sp);
            String val = sp < 0 ? "" : t.substring(sp + 1);
            switch (key) {
                case "case":
                    cur = new Record();
                    cur.name = val;
                    out.add(cur);
                    break;
                case "parse":
                    cur.parse = unescape(val);
                    break;
                case "parseb64":
                    cur.parseBytes = Base64.getDecoder().decode(val);
                    break;
                case "build":
                    cur.build = val;
                    break;
                case "sel":
                    cur.sel = val;
                    break;
                case "divergence":
                    // The corpus records here why a case is known NOT to agree, and the Go
                    // side is required to refuse it rather than compare a golden. Recorded
                    // so that this program can still read the file it generated goldens for;
                    // it changes nothing about what is emitted, because a divergence is a
                    // statement about the Go side and Java's answer is still Java's answer.
                    cur.divergence = cur.divergence == null ? val : cur.divergence + " " + val;
                    break;
                default:
                    throw new IllegalArgumentException("unknown corpus key " + key);
            }
        }
        return out;
    }

    /** Backslash escapes n, r, t, s (space), backslash and u+4 hex digits; ~ alone is "". */
    static String unescape(String s) {
        if (s.equals("~")) return "";
        StringBuilder sb = new StringBuilder(s.length());
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (c != '\\' || i + 1 >= s.length()) {
                sb.append(c);
                continue;
            }
            char n = s.charAt(++i);
            switch (n) {
                case 'n': sb.append('\n'); break;
                case 'r': sb.append('\r'); break;
                case 't': sb.append('\t'); break;
                case 's': sb.append(' '); break;
                case '\\': sb.append('\\'); break;
                case 'u':
                    sb.append((char) Integer.parseInt(s.substring(i + 1, i + 5), 16));
                    i += 4;
                    break;
                default: sb.append('\\').append(n);
            }
        }
        return sb.toString();
    }
}

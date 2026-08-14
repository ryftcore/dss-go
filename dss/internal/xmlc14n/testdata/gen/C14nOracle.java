// Generates testdata/golden/* and testdata/manifest.txt: the canonical bytes Apache Santuario
// 3.0.6 produces, through the exact call path DSS uses, for every (document, scope, algorithm)
// triple listed in gen/cases.txt. The goldens are Java's answers and only Java's; the Go port
// in internal/xmlc14n must reproduce them byte for byte. No golden may ever be produced, fixed
// up or "corrected" from Go output.
//
// Rules this oracle obeys (XML_DESIGN.md section 3.4):
//
//   1. A fresh canonicalizer per KAT. Santuario's inclusive canonicalizers are single-use
//      (SANTUARIO-463: the firstCall field is never reset), so a reused instance silently
//      drops inherited namespaces and xml:* attributes. Every DSS call site constructs one
//      inline, so DSS always gets the first-call behaviour and so must we.
//   2. Document-scope KATs go through XMLCanonicalizer.createInstance(uri).canonicalize(byte[]),
//      the exact DSS call path, which parses with Santuario's own secure parser.
//   3. Apex-scope KATs parse with DomUtils.buildDOM(byte[]) - the DSS secure
//      DocumentBuilderFactory - and then call canonicalize(Node).
//   4. PrefixList KATs need a parameter XMLCanonicalizer does not expose, so they drop to
//      Canonicalizer.getInstance(uri).canonicalizeSubtree(node, prefixList, out).
//      XMLCanonicalizer.canonicalize(Node) is a thin, parameterless wrapper over exactly this
//      call (see XMLCanonicalizer#canonicalize(Node, OutputStream)), so the two paths differ
//      only in that parameter. Only the two exclusive algorithms take a PrefixList;
//      Canonicalizer20010315 throws CanonicalizationException("UnsupportedOperation") for the
//      inclusive ones, and the Go API documents the parameter as ignored there, so those
//      combinations are not KATs and are not emitted.
//   5. A failure is recorded as the single line "!ERROR <SimpleClassName>" in the golden file,
//      never as a missing file, so that a Go-side success where Java fails is a test failure.
//      The class recorded is the root cause, not the DSSException wrapper.
//   6. manifest.txt records the SHA-256 of every golden, so a corrupted regeneration shows up
//      in review as a changed hash next to a changed file.
//   7. This program never reads Go sources, and the Go tests never invoke Java.
//
// Run it with OpenJDK 21, Apache Santuario 3.0.6 and the maven-built DSS 6.5.RC1 jars:
//
//   DSS=/home/user/dss-upstream
//   M2=$HOME/.m2/repository
//   CP=$DSS/dss-xml-utils/target/dss-xml-utils-6.5.RC1.jar\
//   :$DSS/dss-xml-common/target/dss-xml-common-6.5.RC1.jar\
//   :$DSS/dss-model/target/dss-model-6.5.RC1.jar\
//   :$DSS/dss-enumerations/target/dss-enumerations-6.5.RC1.jar\
//   :$DSS/dss-utils/target/dss-utils-6.5.RC1.jar\
//   :$M2/org/apache/santuario/xmlsec/3.0.6/xmlsec-3.0.6.jar\
//   :$DSS/dss-utils-apache-commons/target/dss-utils-apache-commons-6.5.RC1.jar\
//   :$M2/org/apache/commons/commons-lang3/3.20.0/commons-lang3-3.20.0.jar\
//   :$M2/org/apache/commons/commons-collections4/4.5.0/commons-collections4-4.5.0.jar\
//   :$M2/commons-io/commons-io/2.22.0/commons-io-2.22.0.jar\
//   :$M2/commons-codec/commons-codec/1.18.0/commons-codec-1.18.0.jar\
//   :$M2/org/slf4j/slf4j-api/2.0.18/slf4j-api-2.0.18.jar
//   javac -encoding UTF-8 -cp "$CP" -d /tmp/c14noracle C14nOracle.java
//   java -cp "$CP:/tmp/c14noracle" C14nOracle <corpus dir> <golden dir> <manifest file> <cases file>
//
// gen/generate.sh does all of the above.
import eu.europa.esig.dss.xml.utils.DomUtils;
import eu.europa.esig.dss.xml.utils.XMLCanonicalizer;

import org.apache.xml.security.c14n.Canonicalizer;
import org.w3c.dom.Attr;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import java.io.ByteArrayOutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.TreeMap;

public class C14nOracle {

    /** The seven algorithms XMLCanonicalizer.registerDefaultCanonicalizers() registers. */
    private static final String[][] ALGS = {
        {"c14n10",   "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"},
        {"c14n10wc", "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"},
        {"c14n11",   "http://www.w3.org/2006/12/xml-c14n11"},
        {"c14n11wc", "http://www.w3.org/2006/12/xml-c14n11#WithComments"},
        {"excl",     "http://www.w3.org/2001/10/xml-exc-c14n#"},
        {"exclwc",   "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"},
        {"phys",     "http://santuario.apache.org/c14n/physical"},
    };

    private static Path corpusDir;
    private static Path goldenDir;
    private static final StringBuilder manifest = new StringBuilder();

    public static void main(String[] args) throws Exception {
        corpusDir = Paths.get(args[0]);
        goldenDir = Paths.get(args[1]);
        Path manifestFile = Paths.get(args[2]);
        Path casesFile = Paths.get(args[3]);

        Files.createDirectories(goldenDir);
        for (Path p : Files.newDirectoryStream(goldenDir, "*")) {
            if (Files.isRegularFile(p)) {
                Files.delete(p);
            }
        }

        manifest.append("# Known-answer tests for internal/xmlc14n. Generated by gen/C14nOracle.java\n");
        manifest.append("# against Apache Santuario 3.0.6 / DSS 6.5.RC1 / OpenJDK 21. Do not edit by hand,\n");
        manifest.append("# and never regenerate a line to make a Go test pass.\n");
        manifest.append("#\n");
        manifest.append("# document\talgorithm\tscope\tapex\tprefixlist\tgolden\tsha256\tstatus\tdetail\n");

        for (String line : Files.readAllLines(casesFile, StandardCharsets.UTF_8)) {
            if (line.isEmpty() || line.charAt(0) == '#') {
                continue;
            }
            String[] f = split(line, 4);
            runCase(f[0], f[1], f[2], f[3]);
        }

        runNegatives();

        Files.write(manifestFile, manifest.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("wrote " + manifestFile);
    }

    /** One corpus document and scope, canonicalized under every applicable algorithm. */
    private static void runCase(String docName, String scope, String apexId, String prefixList)
        throws Exception {
        byte[] src = Files.readAllBytes(corpusDir.resolve(docName));
        // "#empty" is the cases.txt spelling of "drive the PrefixList overload with the empty
        // string", which is not the same case as having no PrefixList at all: the former goes
        // through Canonicalizer.canonicalizeSubtree(node, prefixList, out) with "", the latter
        // through XMLCanonicalizer.canonicalize(node). InclusiveNamespaces.prefixStr2Set maps
        // both to an empty set, and recording that they agree is the point.
        boolean emptyPrefixList = "#empty".equals(prefixList);
        if (emptyPrefixList) {
            prefixList = "";
        }
        boolean hasPrefixList = emptyPrefixList || !prefixList.isEmpty();

        for (String[] alg : ALGS) {
            String slug = alg[0];
            String uri = alg[1];
            if (hasPrefixList && !slug.startsWith("excl")) {
                // Rule 4: only exclusive c14n accepts a PrefixList.
                continue;
            }
            byte[] out;
            String status = "OK";
            String detail = "";
            try {
                if (apexId.isEmpty()) {
                    // Rule 2: the DSS byte[] call path, parser included.
                    out = XMLCanonicalizer.createInstance(uri).canonicalize(src);
                } else {
                    // Rule 3: the DSS secure DocumentBuilderFactory, then the Node call path.
                    Document doc = DomUtils.buildDOM(src);
                    Node apex = findApex(doc, apexId);
                    if (apex == null) {
                        throw new IllegalStateException(
                            "no element with an Id/ID/xml:id attribute equal to " + apexId
                            + " in " + docName);
                    }
                    if (hasPrefixList) {
                        // Rule 4.
                        ByteArrayOutputStream baos = new ByteArrayOutputStream();
                        Canonicalizer.getInstance(uri).canonicalizeSubtree(apex, prefixList, baos);
                        out = baos.toByteArray();
                    } else {
                        out = XMLCanonicalizer.createInstance(uri).canonicalize(apex);
                    }
                }
            } catch (Throwable t) {
                Throwable root = rootCause(t);
                status = "ERROR";
                detail = oneLine(root.getClass().getSimpleName() + ": " + root.getMessage());
                out = ("!ERROR " + root.getClass().getSimpleName() + "\n")
                    .getBytes(StandardCharsets.UTF_8);
            }
            String golden = base(docName) + "." + slug + "." + scope;
            Files.write(goldenDir.resolve(golden), out);
            manifest.append(docName).append('\t').append(slug).append('\t').append(scope)
                .append('\t').append(apexId).append('\t')
                .append(emptyPrefixList ? "#empty" : prefixList)
                .append('\t').append(golden).append('\t').append(sha256(out))
                .append('\t').append(status).append('\t').append(detail).append('\n');
        }
    }

    /**
     * The negative corpus: documents no conforming parser may accept. Recorded so that the Go
     * parser (internal/xmldom) has a Java answer for each rejection instead of a guess. Both
     * the Santuario parser (the byte[] path) and the DSS DocumentBuilderFactory are exercised,
     * because they are configured separately upstream and can disagree.
     */
    private static void runNegatives() throws Exception {
        Path dir = corpusDir.resolve("negative");
        if (!Files.isDirectory(dir)) {
            return;
        }
        Map<String, byte[]> docs = new TreeMap<>();
        for (Path p : Files.newDirectoryStream(dir, "*.xml")) {
            docs.put(p.getFileName().toString(), Files.readAllBytes(p));
        }
        manifest.append("#\n# negative corpus: parse must fail. \"santuario\" is the byte[] call path,\n");
        manifest.append("# \"dombuilder\" is DomUtils.buildDOM. accepted=... means Java did NOT reject it.\n");
        for (Map.Entry<String, byte[]> e : docs.entrySet()) {
            String name = "negative/" + e.getKey();
            manifest.append(name).append('\t').append("santuario").append('\t')
                .append(negative(e.getValue(), true)).append('\n');
            manifest.append(name).append('\t').append("dombuilder").append('\t')
                .append(negative(e.getValue(), false)).append('\n');
        }
    }

    private static String negative(byte[] src, boolean santuarioPath) {
        try {
            if (santuarioPath) {
                byte[] out = XMLCanonicalizer.createInstance(
                    "http://www.w3.org/TR/2001/REC-xml-c14n-20010315").canonicalize(src);
                return "accepted=" + oneLine(new String(out, StandardCharsets.UTF_8));
            }
            DomUtils.buildDOM(src);
            return "accepted";
        } catch (Throwable t) {
            Throwable root = rootCause(t);
            return root.getClass().getSimpleName() + "\t" + oneLine(root.getMessage());
        }
    }

    /**
     * The apex selection rule, shared verbatim with the Go test: the first element in document
     * order carrying an attribute whose local name equals "id" case-insensitively (so Id, ID,
     * id, xml:id and any prefixed spelling) with the wanted value. Deliberately not XPath and
     * deliberately not Document.getElementById, which needs DTD-declared or explicitly
     * registered ID attributes and would pin a harness artefact.
     */
    private static Node findApex(Document doc, String id) {
        List<Element> stack = new ArrayList<>();
        stack.add(doc.getDocumentElement());
        for (int i = 0; i < stack.size(); i++) {
            Element el = stack.get(i);
            NamedNodeMap attrs = el.getAttributes();
            for (int a = 0; a < attrs.getLength(); a++) {
                Attr attr = (Attr) attrs.item(a);
                String local = attr.getLocalName() != null ? attr.getLocalName() : attr.getName();
                if ("id".equals(local.toLowerCase(Locale.ROOT)) && id.equals(attr.getValue())) {
                    return el;
                }
            }
            List<Element> kids = new ArrayList<>();
            NodeList children = el.getChildNodes();
            for (int c = 0; c < children.getLength(); c++) {
                if (children.item(c).getNodeType() == Node.ELEMENT_NODE) {
                    kids.add((Element) children.item(c));
                }
            }
            stack.addAll(i + 1, kids);
        }
        return null;
    }

    private static Throwable rootCause(Throwable t) {
        Throwable cur = t;
        while (cur.getCause() != null && cur.getCause() != cur) {
            cur = cur.getCause();
        }
        return cur;
    }

    private static String base(String docName) {
        int dot = docName.lastIndexOf('.');
        return dot == -1 ? docName : docName.substring(0, dot);
    }

    private static String sha256(byte[] b) throws Exception {
        byte[] d = MessageDigest.getInstance("SHA-256").digest(b);
        StringBuilder sb = new StringBuilder(64);
        for (byte x : d) {
            sb.append(Character.forDigit((x >> 4) & 0xf, 16)).append(Character.forDigit(x & 0xf, 16));
        }
        return sb.toString();
    }

    /** Keeps the manifest one record per line and free of tabs. */
    private static String oneLine(String s) {
        if (s == null) {
            return "";
        }
        return s.replace("\r", "\\r").replace("\n", "\\n").replace("\t", " ");
    }

    /** Tab-split that keeps trailing empty fields, which String.split drops. */
    private static String[] split(String line, int fields) {
        String[] out = new String[fields];
        int start = 0;
        for (int i = 0; i < fields; i++) {
            int tab = line.indexOf('\t', start);
            if (i == fields - 1 || tab == -1) {
                out[i] = line.substring(Math.min(start, line.length()));
                start = line.length();
            } else {
                out[i] = line.substring(start, tab);
                start = tab + 1;
            }
        }
        return out;
    }
}

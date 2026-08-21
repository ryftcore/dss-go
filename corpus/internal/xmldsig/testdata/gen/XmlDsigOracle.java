// Generates testdata/golden/* and testdata/manifest.txt for internal/xmldsig: for every
// ds:Reference of every ds:Signature in the real upstream XAdES fixtures listed in
// gen/cases.txt, the octets Apache Santuario 3.0.6 feeds to the digest (the transform
// pipeline output, XMLSignatureInput.getBytes()), the digest itself, and the reference's
// verify() verdict; plus, per signature, the canonicalized ds:SignedInfo octets and the
// verdict of checking ds:SignatureValue over them with the ds:KeyInfo public key.
//
// The goldens are Java's answers and only Java's. The Go port in internal/xmldsig must
// reproduce them byte for byte. No golden may ever be produced, fixed up or "corrected"
// from Go output, and this program never reads a Go source file.
//
// What is *not* Santuario here, and why:
//
//   * ID registration. DSS does not let the parser type ID attributes (DOCTYPE is banned),
//     it registers them itself in XAdESDOMDocument.recursiveIdBrowse/setIDIdentifier: per
//     element, in NamedNodeMap order, the FIRST attribute whose local name equals "Id"
//     case-insensitively is declared an ID and the scan for that element stops. registerIds
//     below is a line-for-line copy of that loop. Without it doc.getElementById returns null
//     and every "#foo" reference fails, which is not the behaviour under test.
//   * DetachedSignatureResolver. dss-xades cannot be built offline in this environment (its
//     JAXB codegen needs the network), so detachedResolver below is a reduced port of
//     eu.europa.esig.dss.xades.validation.DetachedSignatureResolver covering the single case
//     the corpus exercises - one detached document, matched by file name or, failing that,
//     used as the sole candidate. The bytes it hands Santuario are the file's bytes; the
//     pipeline that consumes them is Santuario's own.
//   * EnforcedResolverFragment's XPath-injection filter is applied by the same reduced port
//     (checkValueForXpathInjection), because it decides which resolver wins for a "#..." URI.
//
// Every other decision - dereferencing, transform execution, canonicalization, digesting,
// signature checking - is Santuario's, reached through the classes DSS itself calls.
//
// Run it with OpenJDK 21, Apache Santuario 3.0.6 and the maven-built DSS 6.5.RC1 jars; see
// gen/generate.sh, which does all of it.
import org.apache.xml.security.algorithms.SignatureAlgorithm;
import org.apache.xml.security.c14n.Canonicalizer;
import org.apache.xml.security.keys.KeyInfo;
import org.apache.xml.security.signature.Reference;
import org.apache.xml.security.signature.SignedInfo;
import org.apache.xml.security.signature.XMLSignature;
import org.apache.xml.security.signature.XMLSignatureInput;
import org.apache.xml.security.utils.resolver.ResourceResolver;
import org.apache.xml.security.utils.resolver.ResourceResolverContext;
import org.apache.xml.security.utils.resolver.ResourceResolverException;
import org.apache.xml.security.utils.resolver.ResourceResolverSpi;
import org.apache.xml.security.utils.resolver.implementations.ResolverFragment;
import org.apache.xml.security.utils.resolver.implementations.ResolverXPointer;

import org.w3c.dom.Attr;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

public final class XmlDsigOracle {

    private static final String DS = "http://www.w3.org/2000/09/xmldsig#";

    /** Digest algorithm URI -> JCA name, the subset the corpus uses plus the usual suspects. */
    private static String jcaDigest(String uri) {
        switch (uri) {
            case "http://www.w3.org/2000/09/xmldsig#sha1":       return "SHA-1";
            case "http://www.w3.org/2001/04/xmldsig-more#sha224":return "SHA-224";
            case "http://www.w3.org/2001/04/xmlenc#sha256":      return "SHA-256";
            case "http://www.w3.org/2001/04/xmldsig-more#sha384":return "SHA-384";
            case "http://www.w3.org/2001/04/xmlenc#sha512":      return "SHA-512";
            case "http://www.w3.org/2001/04/xmlenc#ripemd160":   return "RIPEMD160";
            default: return null;
        }
    }

    private static final StringBuilder MANIFEST = new StringBuilder();
    private static Path goldenDir;

    public static void main(String[] args) throws Exception {
        Path corpusDir = Paths.get(args[0]);
        goldenDir = Paths.get(args[1]);
        Path manifestFile = Paths.get(args[2]);
        Path casesFile = Paths.get(args[3]);

        org.apache.xml.security.Init.init();
        // XAdESSignature.initDefaultResolvers(): the fragment resolver guarded against XPath
        // injection, then the XPointer resolver. Init.init() already registered the stock
        // ResolverFragment/ResolverXPointer; registering ours in front is what DSS does.
        ResourceResolver.register(new EnforcedFragment(), true);

        Files.createDirectories(goldenDir);
        MANIFEST.append("# Known-answer tests for internal/xmldsig. Generated by gen/XmlDsigOracle.java\n");
        MANIFEST.append("# against Apache Santuario 3.0.6 / OpenJDK 21 over the upstream DSS 6.5.RC1\n");
        MANIFEST.append("# dss-xades test fixtures. Do not edit by hand, and never regenerate a line to\n");
        MANIFEST.append("# make a Go test pass.\n");
        MANIFEST.append("#\n");
        MANIFEST.append("# fixture\tsig\tkind\tindex\tdetail\tgolden\tsha256\tstatus\n");

        for (String line : Files.readAllLines(casesFile, StandardCharsets.UTF_8)) {
            line = line.trim();
            if (line.isEmpty() || line.startsWith("#")) {
                continue;
            }
            String[] parts = line.split("\t", -1);
            String fixture = parts[0];
            List<File> detached = new ArrayList<>();
            if (parts.length > 1 && !parts[1].isEmpty()) {
                for (String d : parts[1].split(",")) {
                    detached.add(corpusDir.resolve(d).toFile());
                }
            }
            try {
                run(fixture, corpusDir.resolve(fixture).toFile(), detached);
            } catch (Throwable t) {
                emit(fixture, "-", "fixture", "-", "-", null, rootCause(t));
            }
        }

        Files.write(manifestFile, MANIFEST.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("wrote " + manifestFile);
    }

    private static void run(String fixture, File file, List<File> detached) throws Exception {
        Document doc = parse(file);
        registerIds(doc.getDocumentElement());

        NodeList sigs = doc.getElementsByTagNameNS(DS, "Signature");
        for (int s = 0; s < sigs.getLength(); s++) {
            Element sigEl = (Element) sigs.item(s);
            // Only top-level signatures: a ds:Signature nested inside another signature's
            // ds:Object is a counter signature and is validated through its own path.
            String sigId = String.valueOf(s);
            try {
                oneSignature(fixture, sigId, sigEl, detached);
            } catch (Throwable t) {
                emit(fixture, sigId, "signature", "-", "-", null, rootCause(t));
            }
        }
    }

    private static void oneSignature(String fixture, String sigId, Element sigEl, List<File> detached)
            throws Exception {
        // DSS: new XMLSignature(signatureElement, "", false) - secure validation OFF.
        XMLSignature sig = new XMLSignature(sigEl, "", false);
        if (!detached.isEmpty()) {
            sig.addResourceResolver(new DetachedResolver(detached));
        }

        SignedInfo si = sig.getSignedInfo();

        // ---- ds:SignedInfo canonical octets
        byte[] c14n = null;
        try {
            c14n = si.getCanonicalizedOctetStream();
            emit(fixture, sigId, "signedinfo", "-", si.getCanonicalizationMethodURI(), c14n, null);
        } catch (Throwable t) {
            emit(fixture, sigId, "signedinfo", "-", "-", null, rootCause(t));
        }

        // ---- ds:SignatureValue over those octets, with the ds:KeyInfo public key
        String sigVerdict;
        try {
            PublicKey pk = keyInfoPublicKey(sigEl);
            if (pk == null) {
                sigVerdict = "nokey";
            } else if (c14n == null) {
                sigVerdict = "noc14n";
            } else {
                SignatureAlgorithm sa = si.getSignatureAlgorithm();
                sa.initVerify(pk);
                sa.update(c14n);
                sigVerdict = sa.verify(sig.getSignatureValue()) ? "true" : "false";
            }
        } catch (Throwable t) {
            sigVerdict = "!" + rootCause(t);
        }
        emitLiteral(fixture, sigId, "sigvalue", "-", si.getSignatureMethodURI(), sigVerdict);

        // ---- every ds:Reference
        int n = si.getLength();
        for (int i = 0; i < n; i++) {
            reference(fixture, sigId, "ref", i, si.item(i), detached);
        }
    }

    private static void reference(String fixture, String sigId, String kind, int i, Reference ref,
                                  List<File> detached) {
        String uri = ref.getURI() == null ? "" : ref.getURI();
        byte[] bytes = null;
        try {
            // getReferencedBytes() == dereferenceURIandPerformTransforms(null).getBytes():
            // exactly the octets calculateDigest streams into the digester, and exactly what
            // DSSXMLUtils.getReferenceOriginalContentBytes reads for an enveloped reference.
            bytes = ref.getReferencedBytes();
            emit(fixture, sigId, kind + "bytes", String.valueOf(i), uri, bytes, null);
        } catch (Throwable t) {
            emit(fixture, sigId, kind + "bytes", String.valueOf(i), uri, null, rootCause(t));
        }

        String digest;
        try {
            String algUri = digestUri(ref);
            String jca = jcaDigest(algUri);
            if (bytes == null) {
                digest = "-";
            } else if (jca == null) {
                digest = "!UnknownDigest:" + algUri;
            } else {
                digest = hex(MessageDigest.getInstance(jca).digest(bytes));
            }
        } catch (Throwable t) {
            digest = "!" + rootCause(t);
        }
        emitLiteral(fixture, sigId, kind + "digest", String.valueOf(i), uri, digest);

        String verdict;
        try {
            verdict = ref.verify() ? "true" : "false";
        } catch (Throwable t) {
            verdict = "!" + rootCause(t);
        }
        emitLiteral(fixture, sigId, kind + "verify", String.valueOf(i), uri, verdict);

        // ---- ds:Manifest referenced by Type: its own references, one level deep
        try {
            if (Reference.MANIFEST_URI.equals(ref.getType())) {
                Element man = manifestElement(ref);
                if (man != null) {
                    org.apache.xml.security.signature.Manifest m =
                        new org.apache.xml.security.signature.Manifest(man, "", false);
                    if (!detached.isEmpty()) {
                        m.addResourceResolver(new DetachedResolver(detached));
                    }
                    for (int j = 0; j < m.getLength(); j++) {
                        reference(fixture, sigId, "man" + i + ".", j, m.item(j), detached);
                    }
                }
            }
        } catch (Throwable t) {
            emitLiteral(fixture, sigId, "manifest", String.valueOf(i), uri, "!" + rootCause(t));
        }
    }

    /** ds:Manifest with the Id the reference's "#id" URI names, searched in the whole document. */
    private static Element manifestElement(Reference ref) {
        String uri = ref.getURI();
        if (uri == null || !uri.startsWith("#")) {
            return null;
        }
        String id = uri.substring(1);
        Document doc = ref.getElement().getOwnerDocument();
        NodeList manifests = doc.getElementsByTagNameNS(DS, "Manifest");
        for (int i = 0; i < manifests.getLength(); i++) {
            Element m = (Element) manifests.item(i);
            if (id.equals(m.getAttribute("Id"))) {
                return m;
            }
        }
        return null;
    }

    private static String digestUri(Reference ref) {
        Element dm = firstChildNS(ref.getElement(), DS, "DigestMethod");
        return dm == null ? "" : dm.getAttributeNS(null, "Algorithm");
    }

    private static Element firstChildNS(Element parent, String ns, String local) {
        for (Node n = parent.getFirstChild(); n != null; n = n.getNextSibling()) {
            if (n.getNodeType() == Node.ELEMENT_NODE
                    && ns.equals(n.getNamespaceURI()) && local.equals(n.getLocalName())) {
                return (Element) n;
            }
        }
        return null;
    }

    private static PublicKey keyInfoPublicKey(Element sigEl) {
        try {
            Element ki = firstChildNS(sigEl, DS, "KeyInfo");
            if (ki == null) {
                return null;
            }
            return new KeyInfo(ki, "").getPublicKey();
        } catch (Throwable t) {
            return null;
        }
    }

    // ------------------------------------------------------------------ parsing / IDs

    private static Document parse(File f) throws Exception {
        // DomUtils.buildDOM's secure DocumentBuilderFactory: namespace aware, no DOCTYPE, no
        // external entities.
        DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
        dbf.setNamespaceAware(true);
        dbf.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
        dbf.setFeature("http://xml.org/sax/features/external-general-entities", false);
        dbf.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
        dbf.setExpandEntityReferences(false);
        dbf.setXIncludeAware(false);
        DocumentBuilder db = dbf.newDocumentBuilder();
        try (InputStream in = new FileInputStream(f)) {
            return db.parse(in);
        }
    }

    /** XAdESDOMDocument.recursiveIdBrowse + setIDIdentifier, verbatim. */
    private static void registerIds(Element element) {
        NamedNodeMap attributes = element.getAttributes();
        for (int jj = 0; jj < attributes.getLength(); jj++) {
            Node item = attributes.item(jj);
            String localName = item.getLocalName();
            String nodeName = item.getNodeName();
            if (localName != null && "Id".equalsIgnoreCase(localName)) {
                element.setIdAttribute(nodeName, true);
                break;
            }
        }
        NodeList children = element.getChildNodes();
        for (int ii = 0; ii < children.getLength(); ii++) {
            Node childNode = children.item(ii);
            if (childNode.getNodeType() == Node.ELEMENT_NODE) {
                registerIds((Element) childNode);
            }
        }
    }

    // ------------------------------------------------------------------ resolvers

    /** EnforcedResolverFragment: ResolverFragment plus the XPath-injection character filter. */
    static final class EnforcedFragment extends ResolverFragment {
        private static final String XPATH_CHAR_FILTER = "()='[]:,*/ ";

        @Override
        public boolean engineCanResolveURI(ResourceResolverContext context) {
            return checkValueForXpathInjection(context.uriToResolve) && super.engineCanResolveURI(context);
        }

        static boolean checkValueForXpathInjection(String xpathString) {
            if (xpathString != null && !xpathString.isEmpty()) {
                String decoded = decodeURI(xpathString);
                for (char c : decoded.toCharArray()) {
                    if (XPATH_CHAR_FILTER.indexOf(c) != -1) {
                        return false;
                    }
                }
            }
            return true;
        }
    }

    /** DSSUtils.decodeURI: URLDecoder.decode(uri, UTF-8), returning the input on failure. */
    static String decodeURI(String uri) {
        try {
            return java.net.URLDecoder.decode(uri, StandardCharsets.UTF_8);
        } catch (Exception e) {
            return uri;
        }
    }

    /**
     * Reduced DetachedSignatureResolver: resolves a non-"#" URI to the one detached file it
     * was built with, matching by name when the URI names it and otherwise accepting it as
     * the sole candidate (the "single detached document reference" branch upstream).
     */
    static final class DetachedResolver extends ResourceResolverSpi {
        private final List<File> files;

        DetachedResolver(List<File> files) {
            this.files = files;
        }

        @Override
        public boolean engineCanResolveURI(ResourceResolverContext context) {
            return context.attr == null || definedFilename(context.attr.getNodeValue());
        }

        private static boolean definedFilename(String uri) {
            return uri != null && !uri.trim().isEmpty() && !uri.startsWith("#");
        }

        private File best(ResourceResolverContext context) {
            if (context.attr == null) {
                return files.size() == 1 ? files.get(0) : null;
            }
            String uri = decodeURI(context.attr.getNodeValue());
            for (File f : files) {
                if (f.getName().equals(uri)) {
                    return f;
                }
            }
            return files.size() == 1 ? files.get(0) : null;
        }

        @Override
        public XMLSignatureInput engineResolveURI(ResourceResolverContext context)
                throws ResourceResolverException {
            File file = best(context);
            if (file == null) {
                throw new ResourceResolverException("generic.EmptyMessage",
                        new Object[]{"Unable to find document (detached signature)"},
                        context.uriToResolve, context.baseUri);
            }
            try {
                XMLSignatureInput in = new XMLSignatureInput(Files.readAllBytes(file.toPath()));
                in.setSourceURI(file.getName());
                return in;
            } catch (Exception e) {
                throw new ResourceResolverException("generic.EmptyMessage",
                        new Object[]{e.getMessage()}, context.uriToResolve, context.baseUri);
            }
        }
    }

    // ------------------------------------------------------------------ output

    private static void emit(String fixture, String sig, String kind, String index, String detail,
                             byte[] payload, String error) {
        if (payload == null) {
            MANIFEST.append(String.join("\t", fixture, sig, kind, index, detail, "", "",
                    "!ERROR " + error)).append('\n');
            return;
        }
        String name = (fixture.replace('/', '_') + "." + sig + "." + kind + "." + index)
                .replace("..", ".");
        try {
            Files.write(goldenDir.resolve(name), payload);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        MANIFEST.append(String.join("\t", fixture, sig, kind, index, detail, name,
                sha256(payload), "OK")).append('\n');
    }

    private static void emitLiteral(String fixture, String sig, String kind, String index,
                                    String detail, String value) {
        MANIFEST.append(String.join("\t", fixture, sig, kind, index, detail, "", "", value))
                .append('\n');
    }

    private static String sha256(byte[] b) {
        try {
            return hex(MessageDigest.getInstance("SHA-256").digest(b));
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private static String hex(byte[] b) {
        StringBuilder sb = new StringBuilder(b.length * 2);
        for (byte x : b) {
            sb.append(Character.forDigit((x >> 4) & 0xf, 16)).append(Character.forDigit(x & 0xf, 16));
        }
        return sb.toString().toLowerCase(Locale.ROOT);
    }

    /** The root cause's simple class name: the wrapper chain is Santuario's, not information. */
    private static String rootCause(Throwable t) {
        Throwable c = t;
        while (c.getCause() != null && c.getCause() != c) {
            c = c.getCause();
        }
        return c.getClass().getSimpleName();
    }

    private XmlDsigOracle() {
    }
}

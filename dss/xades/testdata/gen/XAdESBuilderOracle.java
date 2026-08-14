// Generates testdata/xades-builder-oracle.json: the exact XML that upstream DSS 6.5.RC1's
// XAdESBuilder produces for the DOM fragments this porter chunk owns - xades:CertDigest,
// xades:Cert (with IssuerSerial for en319132=false and IssuerSerialV2 for en319132=true) and
// xades141:SPDocSpecification - plus the deterministic XML Id that XAdESBuilder#toXmlIdentifier
// derives from a DSS Identifier.
//
// These are the byte-compatibility surfaces of xades_builder.go: element order, namespace
// placement, prefixes and attribute values all come out of this file, so the Go KAT
// (xades_builder_kat_test.go) compares its own DOM serialization against these strings rather
// than against anything hand-derived.
//
// The oracle lives in package eu.europa.esig.dss.xades.signature because XAdESBuilder is
// abstract and every member it exercises is protected.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-xades dependency:build-classpath -Dmdep.outputFile=/tmp/cp-xades.txt \
//       -Dmdep.includeScope=test
//   CP="dss-xades/target/classes:$(cat /tmp/cp-xades.txt)"
//   javac -cp "$CP" -d /tmp/xadesoracle XAdESBuilderOracle.java
//   java  -cp "$CP:/tmp/xadesoracle" \
//         eu.europa.esig.dss.xades.signature.XAdESBuilderOracle <testdata directory>
package eu.europa.esig.dss.xades.signature;

import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.ObjectIdentifierQualifier;
import eu.europa.esig.dss.model.SpDocSpecification;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.xades.XAdESSignatureParameters;
import eu.europa.esig.dss.xades.definition.XAdESNamespace;
import eu.europa.esig.dss.xml.common.definition.DSSNamespace;
import eu.europa.esig.dss.xml.utils.DomUtils;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class XAdESBuilderOracle {

    /** The three XAdES namespaces XAdESBuilder#getCurrentXAdESElements accepts. */
    static final DSSNamespace[] NAMESPACES = {
            XAdESNamespace.XADES_132, XAdESNamespace.XADES_122, XAdESNamespace.XADES_111
    };

    /** Digest algorithms exercised for CertDigest / DigestAlgAndValue. */
    static final DigestAlgorithm[] DIGESTS = {
            DigestAlgorithm.SHA1, DigestAlgorithm.SHA256, DigestAlgorithm.SHA512
    };

    /** A concrete XAdESBuilder: the class under test is abstract only in alignNodes(). */
    static final class OracleBuilder extends XAdESBuilder {
        OracleBuilder(XAdESSignatureParameters params, Document documentDom) {
            this.params = params;
            this.documentDom = documentDom;
        }

        @Override
        protected void alignNodes() {
            // never reached: the oracle never calls createXmlDocument()
        }
    }

    public static void main(String[] args) throws Exception {
        Path testdata = Paths.get(args.length > 0 ? args[0] : ".");
        CertificateToken certificate = DSSUtils.loadCertificate(
                Files.readAllBytes(testdata.resolve("signer.crt")));

        List<Map<String, Object>> cases = new ArrayList<>();

        for (DSSNamespace namespace : NAMESPACES) {
            for (boolean en319132 : new boolean[] { false, true }) {
                for (DigestAlgorithm digestAlgorithm : DIGESTS) {
                    cases.add(certDigestCase(certificate, namespace, en319132, digestAlgorithm));
                    cases.add(certCase(certificate, namespace, en319132, digestAlgorithm));
                }
                cases.add(spDocSpecificationCase(namespace, en319132));
            }
        }
        cases.add(xmlIdentifierCase(certificate));

        write(testdata.resolve("xades-builder-oracle.json"), cases);
    }

    /** <CertDigest> built by XAdESBuilder#incorporateCertDigest. */
    static Map<String, Object> certDigestCase(CertificateToken certificate, DSSNamespace namespace,
                                              boolean en319132, DigestAlgorithm digestAlgorithm) {
        Document document = DomUtils.buildDOM();
        Element root = newRoot(document);
        OracleBuilder builder = new OracleBuilder(params(namespace, en319132), document);
        try {
            builder.incorporateCertDigest(root, digestAlgorithm, certificate);
        } catch (RuntimeException e) {
            return failure("incorporateCertDigest", namespace, en319132, digestAlgorithm, e);
        }
        return entry("incorporateCertDigest", namespace, en319132, digestAlgorithm, document);
    }

    /** <Cert> built by XAdESBuilder#incorporateCert (IssuerSerial or IssuerSerialV2). */
    static Map<String, Object> certCase(CertificateToken certificate, DSSNamespace namespace,
                                        boolean en319132, DigestAlgorithm digestAlgorithm) {
        Document document = DomUtils.buildDOM();
        Element root = newRoot(document);
        OracleBuilder builder = new OracleBuilder(params(namespace, en319132), document);
        try {
            builder.incorporateCert(root, certificate, digestAlgorithm);
        } catch (RuntimeException e) {
            return failure("incorporateCert", namespace, en319132, digestAlgorithm, e);
        }
        return entry("incorporateCert", namespace, en319132, digestAlgorithm, document);
    }

    /** <xades141:SPDocSpecification> built by XAdESBuilder#incorporateSPDocSpecification. */
    static Map<String, Object> spDocSpecificationCase(DSSNamespace namespace, boolean en319132) {
        Document document = DomUtils.buildDOM();
        Element root = newRoot(document);
        OracleBuilder builder = new OracleBuilder(params(namespace, en319132), document);

        SpDocSpecification spDocSpecification = new SpDocSpecification();
        spDocSpecification.setId("1.2.3.4.5");
        spDocSpecification.setQualifier(ObjectIdentifierQualifier.OID_AS_URN);
        spDocSpecification.setDescription("DSS Go port oracle policy");
        spDocSpecification.setDocumentationReferences("http://nowina.lu/ref1", "http://nowina.lu/ref2");

        try {
            builder.incorporateSPDocSpecification(root, spDocSpecification);
        } catch (RuntimeException e) {
            return failure("incorporateSPDocSpecification", namespace, en319132, null, e);
        }
        return entry("incorporateSPDocSpecification", namespace, en319132, null, document);
    }

    /** XAdESBuilder#toXmlIdentifier, i.e. "id-" + SHA1 of Identifier#asXmlId(). */
    static Map<String, Object> xmlIdentifierCase(CertificateToken certificate) {
        OracleBuilder builder = new OracleBuilder(params(XAdESNamespace.XADES_132, false),
                DomUtils.buildDOM());
        Map<String, Object> entry = new LinkedHashMap<>();
        entry.put("operation", "toXmlIdentifier");
        entry.put("asXmlId", certificate.getDSSId().asXmlId());
        entry.put("xmlIdentifier", builder.toXmlIdentifier(certificate.getDSSId()));
        return entry;
    }

    static XAdESSignatureParameters params(DSSNamespace namespace, boolean en319132) {
        XAdESSignatureParameters params = new XAdESSignatureParameters();
        params.setXadesNamespace(namespace);
        params.setEn319132(en319132);
        return params;
    }

    /** A namespace-carrying host element, so the dump shows where declarations land. */
    static Element newRoot(Document document) {
        Element root = document.createElementNS("urn:dss:go:oracle", "o:Root");
        document.appendChild(root);
        return root;
    }

    /** An operation upstream refuses for this namespace: the exception is the ground truth. */
    static Map<String, Object> failure(String operation, DSSNamespace namespace, boolean en319132,
                                       DigestAlgorithm digestAlgorithm, RuntimeException exception) {
        Map<String, Object> entry = new LinkedHashMap<>();
        entry.put("operation", operation);
        entry.put("xadesNamespaceUri", namespace.getUri());
        entry.put("xadesNamespacePrefix", namespace.getPrefix());
        entry.put("en319132", en319132);
        entry.put("digestAlgorithm", digestAlgorithm == null ? null : digestAlgorithm.name());
        entry.put("unsupported", exception.getMessage());
        return entry;
    }

    static Map<String, Object> entry(String operation, DSSNamespace namespace, boolean en319132,
                                     DigestAlgorithm digestAlgorithm, Document document) {
        Map<String, Object> entry = new LinkedHashMap<>();
        entry.put("operation", operation);
        entry.put("xadesNamespaceUri", namespace.getUri());
        entry.put("xadesNamespacePrefix", namespace.getPrefix());
        entry.put("en319132", en319132);
        entry.put("digestAlgorithm", digestAlgorithm == null ? null : digestAlgorithm.name());
        entry.put("xml", new String(DomUtils.serializeNode(document), StandardCharsets.UTF_8));
        return entry;
    }

    // ---------------------------------------------------------------- minimal JSON writer

    static void write(Path target, List<Map<String, Object>> cases) throws Exception {
        StringBuilder out = new StringBuilder();
        out.append("[\n");
        for (int i = 0; i < cases.size(); i++) {
            out.append("  {");
            boolean first = true;
            for (Map.Entry<String, Object> field : cases.get(i).entrySet()) {
                if (!first) {
                    out.append(", ");
                }
                first = false;
                out.append(quote(field.getKey())).append(": ").append(value(field.getValue()));
            }
            out.append("}");
            if (i < cases.size() - 1) {
                out.append(",");
            }
            out.append("\n");
        }
        out.append("]\n");
        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(target, StandardCharsets.UTF_8))) {
            writer.print(out);
        }
    }

    static String value(Object object) {
        if (object == null) {
            return "null";
        }
        if (object instanceof Boolean) {
            return object.toString();
        }
        return quote(object.toString());
    }

    static String quote(String text) {
        StringBuilder out = new StringBuilder("\"");
        for (int i = 0; i < text.length(); i++) {
            char c = text.charAt(i);
            switch (c) {
                case '"': out.append("\\\""); break;
                case '\\': out.append("\\\\"); break;
                case '\n': out.append("\\n"); break;
                case '\r': out.append("\\r"); break;
                case '\t': out.append("\\t"); break;
                default:
                    if (c < 0x20) {
                        out.append(String.format("\\u%04x", (int) c));
                    } else {
                        out.append(c);
                    }
            }
        }
        return out.append("\"").toString();
    }
}

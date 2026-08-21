// Broad differential corpus oracle (JAdES), Java side. Walks the ENTIRE upstream
// dss-jades/src/test/resources corpus (every *.json fixture, not just the ~26 sampled into
// ../upstream/), runs each one through upstream DSS 6.5.RC1's own JWSDocumentAnalyzerFactory /
// JAdESSignature, and dumps a rich per-signature record. gobroad_main.go produces the identical
// shape from the Go port; diffing the two JSON files is the parity proof. Same role
// dss/pades/testdata/broadgen/BroadOracle.java plays for PAdES - the sampled suite there missed
// 16 defects that only the broad run surfaced.
//
// Beyond what ../gen/CrossValidationOracle.java dumps (certificate, signing time, level,
// counter-signature structure, reference/signature integrity, sigD mechanism, serialization type,
// structural-error presence), this also dumps certificate/revocation EXTRACTION (certificate
// source sizes and ref counts, CRL/OCSP binary counts), the timestamp inventory (per type, with
// message-imprint found/intact and token signature intactness), signature scopes, and the
// signature algorithm - the columns a silently wrong etsiU/xVals/rVals/tstTokens parse would show
// up in but a signature-count/level dump would not.
//
// Build & run (OpenJDK 21, maven-built upstream):
//
//   cd $DSS_UPSTREAM_HOME  (your built upstream DSS 6.5.RC1 checkout)
//   mvn -q -o -pl dss-jades dependency:build-classpath -Dmdep.outputFile=/tmp/jcp.txt -Dmdep.includeScope=test
//   CP="dss-jades/target/classes:$(cat /tmp/jcp.txt)"
//   javac -cp "$CP" -d /tmp/jbroad BroadOracle.java
//   java  -cp "$CP:/tmp/jbroad" BroadOracle <corpus root> <output json>
//
// where <corpus root> is <your upstream DSS checkout>/dss-jades/src/test/resources.
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.JWSSerializationType;
import eu.europa.esig.dss.enumerations.SigDMechanism;
import eu.europa.esig.dss.jades.validation.AbstractJWSDocumentAnalyzer;
import eu.europa.esig.dss.jades.validation.JAdESSignature;
import eu.europa.esig.dss.jades.validation.JWSDocumentAnalyzerFactory;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.Digest;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.signature.SignatureCryptographicVerification;
import eu.europa.esig.dss.model.scope.SignatureScope;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.signature.AdvancedSignature;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.tsp.TimestampToken;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class BroadOracle {

    /** Detached fixtures and the original content upstream's own tests feed them, transcribed
     *  1:1 from the corresponding upstream test classes (JWSSimpleDetachedTest,
     *  JWSSimpleDetachedWithWrongAlgoValidationTest, JAdESDetachedByUriWithURLEncodedParsTest,
     *  JAdESDetachedByUriByHashWithURLEncodedParsTest,
     *  JAdESLevelBWithObjectIdByUriHashMechanismTest). Value = list of "relativeFile[|overrideName]"
     *  entries; an entry of the form "inline:NAME:CONTENT" is an InMemoryDocument. */
    static final Map<String, String[]> DETACHED_CONTENT = new LinkedHashMap<>();
    static {
        DETACHED_CONTENT.put("validation/simple-detached.json", new String[] { "sample.json" });
        DETACHED_CONTENT.put("validation/simple-detached-wrong-algo.json", new String[] { "sample.json" });
        DETACHED_CONTENT.put("validation/jades-detached-by-uri-encoded-pars.json", new String[] {
                "ObjectIdByURI-1.html|https://nowina.lu/pub/JAdES/ObjectIdByURI-1.html",
                "ObjectIdByURI-2.html|https://nowina.lu/pub/JAdES/ObjectIdByURI-2.html" });
        DETACHED_CONTENT.put("validation/jades-detached-by-uri-hash-encoded-pars.json", new String[] {
                "ObjectIdByURIHash-1.html|https://signature-plugtests.etsi.org/pub/JAdES/ObjectIdByURIHash-1.html",
                "ObjectIdByURIHash-2.html|https://signature-plugtests.etsi.org/pub/JAdES/ObjectIdByURIHash-2.html" });
        DETACHED_CONTENT.put("validation/jades-flattened-BpB-detached-objectByURIHash.json", new String[] {
                "inline:TEST-DOC.txt:TL-039 Test Document" });
    }

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        List<Path> files = new ArrayList<>();
        Files.walk(root).filter(p -> p.toString().toLowerCase().endsWith(".json")).forEach(files::add);
        Collections.sort(files);
        StringBuilder json = new StringBuilder();
        json.append("{\n  \"files\": [\n");
        boolean first = true;
        for (Path p : files) {
            String rel = root.relativize(p).toString().replace('\\', '/');
            StringBuilder one = new StringBuilder();
            try {
                dumpFile(one, root, rel);
            } catch (Throwable t) {
                one.setLength(0);
                one.append("    {\n      \"path\": ").append(str(rel)).append(",\n      \"error\": ")
                        .append(str(messageChain(t))).append("\n    }");
            }
            if (!first) {
                json.append(",\n");
            }
            first = false;
            json.append(one);
        }
        json.append("\n  ]\n}\n");
        Files.write(Paths.get(args[1]), json.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("written " + args[1] + " (" + files.size() + " files)");
    }

    static String messageChain(Throwable t) {
        StringBuilder chain = new StringBuilder();
        for (Throwable c = t; c != null; c = c.getCause()) {
            if (c.getMessage() != null) {
                if (chain.length() > 0) {
                    chain.append(" | ");
                }
                chain.append(c.getMessage());
            }
        }
        return chain.length() == 0 ? t.getClass().getName() : chain.toString();
    }

    static void dumpFile(StringBuilder json, Path root, String relativePath) throws Exception {
        Path filePath = root.resolve(relativePath);
        DSSDocument document = new FileDocument(filePath.toFile());

        JWSDocumentAnalyzerFactory factory = new JWSDocumentAnalyzerFactory();
        AbstractJWSDocumentAnalyzer analyzer = factory.create(document);
        CommonCertificateVerifier certificateVerifier = new CommonCertificateVerifier();
        analyzer.setCertificateVerifier(certificateVerifier);

        String[] detached = DETACHED_CONTENT.get(relativePath);
        if (detached != null) {
            List<DSSDocument> detachedDocuments = new ArrayList<>();
            for (String entry : detached) {
                if (entry.startsWith("inline:")) {
                    String rest = entry.substring("inline:".length());
                    int colon = rest.indexOf(':');
                    detachedDocuments.add(new InMemoryDocument(
                            rest.substring(colon + 1).getBytes(StandardCharsets.UTF_8), rest.substring(0, colon)));
                } else {
                    int bar = entry.indexOf('|');
                    String file = bar < 0 ? entry : entry.substring(0, bar);
                    DSSDocument d = new FileDocument(root.resolve(file).toFile());
                    if (bar >= 0) {
                        d.setName(entry.substring(bar + 1));
                    }
                    detachedDocuments.add(d);
                }
            }
            analyzer.setDetachedContents(detachedDocuments);
        }

        List<AdvancedSignature> signatures = analyzer.getSignatures();

        json.append("    {\n");
        json.append("      \"path\": ").append(str(relativePath)).append(",\n");
        json.append("      \"signatureCount\": ").append(signatures.size()).append(",\n");
        json.append("      \"signatures\": [\n");
        for (int i = 0; i < signatures.size(); i++) {
            dumpSignature(json, (JAdESSignature) signatures.get(i), certificateVerifier);
            json.append(i == signatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("      ]\n");
        json.append("    }");
    }

    static void dumpSignature(StringBuilder json, JAdESSignature signature, CommonCertificateVerifier certificateVerifier) {
        signature.initBaselineRequirementsChecker(certificateVerifier);

        CertificateToken signingCertificate = nullSafe(() -> signature.getSigningCertificateToken());
        Long signingTimeMillis = signature.getSigningTime() != null ? signature.getSigningTime().getTime() : null;

        SignatureCryptographicVerification verification = nullSafe(() -> signature.getSignatureCryptographicVerification());
        Digest dtbsr = nullSafe(() -> signature.getDataToBeSignedRepresentation());
        SigDMechanism sigDMechanism = nullSafe(() -> signature.getSigDMechanism());
        JWSSerializationType serializationType = signature.getJws().getJwsSerializationType();
        List<String> structureErrors = signature.getStructureValidationResult();

        json.append("        {\n");
        json.append("          \"signingCertificateFound\": ").append(signingCertificate != null).append(",\n");
        json.append("          \"signingCertificateSHA256\": ")
                .append(signingCertificate != null ? str(hex(DSSUtils.digest(DigestAlgorithm.SHA256, signingCertificate.getEncoded()))) : "null")
                .append(",\n");
        json.append("          \"claimedSigningTimeMillis\": ").append(signingTimeMillis != null ? signingTimeMillis.toString() : "null").append(",\n");
        json.append("          \"dataFoundUpToLevel\": ").append(str(nameOf(nullSafe(() -> signature.getDataFoundUpToLevel())))).append(",\n");
        json.append("          \"isCounterSignature\": ").append(signature.isCounterSignature()).append(",\n");
        json.append("          \"referenceDataFound\": ").append(verification != null && verification.isReferenceDataFound()).append(",\n");
        json.append("          \"referenceDataIntact\": ").append(verification != null && verification.isReferenceDataIntact()).append(",\n");
        json.append("          \"signatureIntact\": ").append(verification != null && verification.isSignatureIntact()).append(",\n");
        json.append("          \"dtbsrAlgorithm\": ").append(dtbsr != null ? str(dtbsr.getAlgorithm().name()) : "null").append(",\n");
        json.append("          \"dtbsrHex\": ").append(dtbsr != null ? str(hex(dtbsr.getValue())) : "null").append(",\n");
        json.append("          \"sigDMechanism\": ").append(sigDMechanism != null ? str(sigDMechanism.name()) : "null").append(",\n");
        json.append("          \"serializationType\": ").append(str(serializationType.name())).append(",\n");
        json.append("          \"hasStructureErrors\": ").append(!structureErrors.isEmpty()).append(",\n");
        json.append("          \"signatureAlgorithm\": ").append(str(nameOf(nullSafe(() -> signature.getSignatureAlgorithm())))).append(",\n");

        // --- certificate extraction ---
        Integer certCount = nullSafe(() -> signature.getCertificateSource().getCertificates().size());
        Integer signingCertRefCount = nullSafe(() -> signature.getCertificateSource().getSigningCertificateRefs().size());
        Integer certValuesCount = nullSafe(() -> signature.getCertificateSource().getCertificateValues().size());
        Integer attrCertRefCount = nullSafe(() -> signature.getCertificateSource().getAttributeCertificateRefs().size());
        Integer completeCertRefCount = nullSafe(() -> signature.getCertificateSource().getCompleteCertificateRefs().size());
        json.append("          \"certificateCount\": ").append(num(certCount)).append(",\n");
        json.append("          \"signingCertificateRefCount\": ").append(num(signingCertRefCount)).append(",\n");
        json.append("          \"certificateValuesCount\": ").append(num(certValuesCount)).append(",\n");
        json.append("          \"attributeCertificateRefCount\": ").append(num(attrCertRefCount)).append(",\n");
        json.append("          \"completeCertificateRefCount\": ").append(num(completeCertRefCount)).append(",\n");

        // --- revocation extraction ---
        Integer crlBinaryCount = nullSafe(() -> signature.getCRLSource().getAllRevocationBinaries().size());
        Integer crlRefCount = nullSafe(() -> signature.getCRLSource().getAllRevocationReferences().size());
        Integer ocspBinaryCount = nullSafe(() -> signature.getOCSPSource().getAllRevocationBinaries().size());
        Integer ocspRefCount = nullSafe(() -> signature.getOCSPSource().getAllRevocationReferences().size());
        json.append("          \"crlBinaryCount\": ").append(num(crlBinaryCount)).append(",\n");
        json.append("          \"crlRefCount\": ").append(num(crlRefCount)).append(",\n");
        json.append("          \"ocspBinaryCount\": ").append(num(ocspBinaryCount)).append(",\n");
        json.append("          \"ocspRefCount\": ").append(num(ocspRefCount)).append(",\n");

        // --- signature scopes ---
        List<SignatureScope> scopes = nullSafe(() -> signature.getSignatureScopes());
        json.append("          \"signatureScopes\": [\n");
        if (scopes != null) {
            for (int i = 0; i < scopes.size(); i++) {
                SignatureScope s = scopes.get(i);
                // Java's getDocumentName() is null for a scope with no document; Go's
                // DocumentName() returns "" for the same state (no null String in Go), so both
                // sides normalize the absent name to "" rather than diffing on the rendering.
                json.append("            { \"name\": ").append(str(s.getDocumentName() == null ? "" : s.getDocumentName()))
                        .append(", \"type\": ").append(str(nameOf(s.getType()))).append(" }")
                        .append(i == scopes.size() - 1 ? "\n" : ",\n");
            }
        }
        json.append("          ],\n");

        // --- timestamps ---
        dumpTimestamps(json, "contentTimestamps", nullSafe(() -> signature.getContentTimestamps()));
        dumpTimestamps(json, "signatureTimestamps", nullSafe(() -> signature.getSignatureTimestamps()));
        dumpTimestamps(json, "timestampsX1", nullSafe(() -> signature.getTimestampsX1()));
        dumpTimestamps(json, "timestampsX2", nullSafe(() -> signature.getTimestampsX2()));
        dumpTimestamps(json, "archiveTimestamps", nullSafe(() -> signature.getArchiveTimestamps()));

        List<AdvancedSignature> counterSignatures = nullSafe(() -> signature.getCounterSignatures());
        int counterCount = counterSignatures == null ? 0 : counterSignatures.size();
        json.append("          \"counterSignatureCount\": ").append(counterCount).append(",\n");
        json.append("          \"counterSignatures\": [\n");
        for (int i = 0; i < counterCount; i++) {
            dumpSignature(json, (JAdESSignature) counterSignatures.get(i), certificateVerifier);
            json.append(i == counterCount - 1 ? "\n" : ",\n");
        }
        json.append("          ]\n");
        json.append("        }");
    }

    static void dumpTimestamps(StringBuilder json, String field, List<TimestampToken> timestamps) {
        json.append("          ").append(str(field)).append(": [\n");
        if (timestamps != null) {
            for (int i = 0; i < timestamps.size(); i++) {
                TimestampToken t = timestamps.get(i);
                Long production = t.getGenerationTime() != null ? t.getGenerationTime().getTime() : null;
                json.append("            { \"type\": ").append(str(nameOf(t.getTimeStampType())))
                        .append(", \"dssId\": ").append(str(String.valueOf(nullSafe(() -> t.getDSSIdAsString()))))
                        .append(", \"generationTimeMillis\": ").append(production != null ? production.toString() : "null")
                        .append(", \"messageImprintDataFound\": ").append(bool(nullSafe(() -> t.isMessageImprintDataFound())))
                        .append(", \"messageImprintDataIntact\": ").append(bool(nullSafe(() -> t.isMessageImprintDataIntact())))
                        .append(", \"messageImprintHex\": ").append(t.getMessageImprint() != null && t.getMessageImprint().getValue() != null
                                ? str(hex(t.getMessageImprint().getValue())) : "null")
                        .append(", \"signatureIntact\": ").append(bool(nullSafe(() -> t.isSignatureIntact())))
                        // The timestamped-reference set is where a missed virtual dispatch on
                        // getSignatureTimestampReferences()/getArchiveTimestampReferences() shows
                        // up (JAdES's override folds in getKeyInfoReferences()); sorted so the
                        // two sides' insertion orders are not compared, only the set.
                        .append(", \"timestampedReferences\": ").append(str(refs(t)))
                        .append(" }")
                        .append(i == timestamps.size() - 1 ? "\n" : ",\n");
            }
        }
        json.append("          ],\n");
    }

    /** Sorted "CATEGORY:objectId" rendering of a timestamp token's timestamped references. */
    static String refs(TimestampToken t) {
        List<eu.europa.esig.dss.spi.x509.tsp.TimestampedReference> references =
                nullSafe(() -> t.getTimestampedReferences());
        if (references == null) {
            return "";
        }
        List<String> rendered = new ArrayList<>();
        for (eu.europa.esig.dss.spi.x509.tsp.TimestampedReference r : references) {
            rendered.add(r.getCategory().name() + ":" + r.getObjectId());
        }
        Collections.sort(rendered);
        return String.join(",", rendered);
    }

    interface Supplier<T> {
        T get() throws Exception;
    }

    /** Upstream throws from several of these getters on deliberately-malformed fixtures; the Go
     *  port's counterpart panics. Both sides record "absent" rather than aborting the whole file,
     *  so a divergence shows up as a field diff instead of an all-or-nothing error diff. */
    static <T> T nullSafe(Supplier<T> supplier) {
        try {
            return supplier.get();
        } catch (Throwable t) {
            return null;
        }
    }

    static String nameOf(Object enumValue) {
        if (enumValue == null) {
            return "";
        }
        if (enumValue instanceof Enum) {
            return ((Enum<?>) enumValue).name();
        }
        return String.valueOf(enumValue);
    }

    static String num(Integer i) {
        return i == null ? "-1" : i.toString();
    }

    static String bool(Boolean b) {
        return b == null ? "false" : b.toString();
    }

    static String str(String s) {
        StringBuilder sb = new StringBuilder("\"");
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"': sb.append("\\\""); break;
                case '\\': sb.append("\\\\"); break;
                case '\n': sb.append("\\n"); break;
                case '\r': sb.append("\\r"); break;
                case '\t': sb.append("\\t"); break;
                default:
                    if (c < 0x20) {
                        sb.append(String.format("\\u%04x", (int) c));
                    } else {
                        sb.append(c);
                    }
            }
        }
        return sb.append('"').toString();
    }

    static String hex(byte[] bytes) {
        StringBuilder builder = new StringBuilder(bytes.length * 2);
        for (byte b : bytes) {
            builder.append(Character.forDigit((b >> 4) & 0xf, 16));
            builder.append(Character.forDigit(b & 0xf, 16));
        }
        return builder.toString();
    }
}

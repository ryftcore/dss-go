// Generates testdata/upstream-cross-validation.json: ground-truth facts about the JAdES
// signatures in testdata/upstream/, dumped straight from upstream DSS 6.5.RC1's own
// JWSDocumentAnalyzerFactory / JAdESSignature. This is the golden file that
// jades_upstream_cross_validation_test.go compares its own parse of the same files against
// (direction "UPSTREAM -> GO" of the cross-validation harness, task #12, JAdES extension).
//
// JAdESSignature extends DefaultAdvancedSignature directly (not CAdESSignature/CMS - JAdES has
// no ASN.1 SignerInfo layer), so what this dumps mirrors cades/testdata/gen/CrossValidationOracle.java
// in spirit (signing certificate, claimed signing time, level, counter-signature structure,
// reference/signature integrity) but swaps the CMS SignerId identity anchor for JAdES's own:
// the compact/flattened/full JWS serialization type actually parsed, the 'cty'/'typ' header
// values, the SigDMechanism a detached ('sigD') signature uses (or none for enveloping/attached),
// and getDataToBeSignedRepresentation() (JAdES's message-digest analogue: the digest of
// ASCII(BASE64URL(header) || '.' || payload || '.' || BASE64URL(signature)), whichever digest
// algorithm the JWS_SIGNING_INPUT ReferenceValidation reports) in place of CAdES's flat
// getMessageDigestValue(). Also dumped: getStructureValidationResult() (JAdES-specific 'crit'/
// 'b64'/header-conflict structural checks - several fixtures below are the DSS-2620 regression
// set that exercises exactly this), since a port that silently drops a structural defect Java
// flags would be the dangerous failure mode, same rationale as PAdES's modification-detection
// buckets.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-jades dependency:build-classpath -Dmdep.outputFile=/tmp/jcp.txt -Dmdep.includeScope=test
//   CP="dss-jades/target/classes:$(cat /tmp/jcp.txt)"
//   javac -cp "$CP" -d /tmp/jvaloracle CrossValidationOracle.java
//   java  -cp "$CP:/tmp/jvaloracle" CrossValidationOracle <testdata directory>
//
// Unlike dss-cades/dss-pades, dss-jades needs no extra hand-added module: it does not go through
// dss-cms's ServiceLoader (no CMS layer in JAdES at all) and has no pluggable backend the way
// dss-pades-pdfbox is - dependency:build-classpath alone resolves everything (jose4j, the JSON
// library, BouncyCastle).
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.JWSSerializationType;
import eu.europa.esig.dss.enumerations.SigDMechanism;
import eu.europa.esig.dss.jades.validation.AbstractJWSDocumentAnalyzer;
import eu.europa.esig.dss.jades.validation.JAdESSignature;
import eu.europa.esig.dss.jades.validation.JWSDocumentAnalyzerFactory;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.Digest;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.signature.SignatureCryptographicVerification;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.signature.AdvancedSignature;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.tsp.TimestampToken;
import eu.europa.esig.dss.spi.x509.tsp.TimestampedReference;

import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class CrossValidationOracle {

    /** Files dumped, relative to testdata/upstream/. Covers compact serialization (the
     *  validation/dss2620/* DSS-2620 regression set - 'crit'/'b64' header edge cases, several
     *  structurally non-ETSI-conformant on purpose), full/flattened JSON serialization,
     *  detached ('sigD') packaging, B/T/LT/LTA baseline levels, attached counter-signatures
     *  (including a broken one and a replace-attack), and a couple of genuinely malformed/empty
     *  cases (altered-jws.json is a tampered payload, jws-serialization-no-signatures.json
     *  parses but carries zero signatures). */
    static final String[] FILES = {
        "validation/simple-detached.json",
        "validation/jades-b-no-signing-time.json",
        "validation/jades-b-with-xvals.json",
        "validation/jades-double-sig-different-etsiU.json",
        "validation/jades-lta.json",
        "validation/jades-lta-broken-arcTst.json",
        "validation/jades-t-clear-etsiu.json",
        "validation/jades-t-duplicated-sigtst.json",
        "validation/jades-triple-LTA.json",
        "validation/jades-with-asn1policy.json",
        "validation/jades-with-attached-counter-signature-replace-attack.json",
        "validation/jades-with-broken-attached-counter-signature.json",
        "validation/jades-with-certified.json",
        "validation/jades-with-counter-signature.json",
        "validation/jades-with-multiple-sign-cert-refs.json",
        // Content ('adoTst') time-stamps, i.e. time-stamps carried by a SIGNED attribute rather
        // than by 'etsiU'. They are the only fixtures whose time-stamp identifiers exercise
        // getAttributeOrder over the signed signature properties, which is where a port that
        // hands back freshly-allocated attribute objects on every call silently loses the
        // attribute's position (see the timestamps[] dump below).
        "validation/jades-b-copied-cnttst.json",
        "validation/jades-t-copied-sigtst.json",
        // Two 'sigRTst'/'rfsTst' time-stamps on top of a signature time-stamp: the X1/X2 buckets,
        // which take their timestamped references from a different base code path than the
        // signature time-stamp does.
        "validation/jades-with-sigAndRefsTst-with-dot.json",
        "validation/altered-jws.json",
        "validation/jws-serialization-no-signatures.json",
        "validation/dss2620/jades-b-level-with-etsiu-in-crit.json",
        "validation/dss2620/jades-t-level-with-etsiu-in-crit.json",
        "validation/dss2620/jades-with-b64-not-in-crit.json",
        "validation/dss2620/jades-with-crit-string-type.json",
        "validation/dss2620/jades-with-crit-with-duplicated-header.json",
        "validation/dss2620/jades-with-crit-with-wrong-entry-type.json",
        "validation/dss2620/jades-with-critical-headers-in-crit.json",
        "validation/dss2620/jades-with-empty-crit.json",
        "validation/dss2620/jades-with-expiration-time-header-in-crit.json",
        "validation/dss2620/jades-without-non-critical-header-in-crit.json",
    };

    /** The one genuinely DETACHED fixture and its original content, checked in alongside it
     *  (mirrors CAdES/XAdES/PAdES oracles' DETACHED_CONTENT map). */
    static final Map<String, String> DETACHED_CONTENT = new LinkedHashMap<>();
    static {
        DETACHED_CONTENT.put("validation/simple-detached.json", "sample.json");
    }

    /** Files upstream's own JWSCompactDocumentAnalyzer/JWSDocumentAnalyzerFactory.create()
     *  refuses outright (raises IllegalInputException while parsing, before a JAdESSignature
     *  ever exists) - the "structure-error" half of the fixture set that can't go through
     *  dumpSignature() at all. jades-with-double-sigt.json embeds the 'sigT' header twice at the
     *  raw JSON level in the protected header, which jose4j's DupeKeyDisallowingLinkedHashMap
     *  rejects during JSON parsing itself (see JAdESWithDoubleSigTTest.java upstream, which
     *  asserts exactly this exception rather than validating a JAdESSignature at all): file ->
     *  substring expected in the exception's message chain (getMessage() + getCause().getMessage()). */
    static final Map<String, String> EXCEPTION_FILES = new LinkedHashMap<>();
    static {
        EXCEPTION_FILES.put("validation/jades-with-double-sigt.json", "An entry for 'sigT' already exists");
    }

    public static void main(String[] args) throws Exception {
        Path testdata = Paths.get(args[0]).resolve("upstream");
        StringBuilder json = new StringBuilder();
        json.append("{\n");
        json.append("  \"_comment\": \"Ground truth from upstream DSS 6.5.RC1, generated by gen/CrossValidationOracle.java. Do not edit by hand.\",\n");
        json.append("  \"files\": [\n");
        for (int i = 0; i < FILES.length; i++) {
            dumpFile(json, testdata, FILES[i]);
            json.append(i == FILES.length - 1 ? "\n" : ",\n");
        }
        json.append("  ],\n");
        json.append("  \"exceptionFiles\": [\n");
        int ei = 0;
        int en = EXCEPTION_FILES.size();
        for (Map.Entry<String, String> entry : EXCEPTION_FILES.entrySet()) {
            dumpExceptionFile(json, testdata, entry.getKey(), entry.getValue());
            json.append(++ei == en ? "\n" : ",\n");
        }
        json.append("  ]\n");
        json.append("}\n");

        Path out = Paths.get(args[0]).resolve("upstream-cross-validation.json");
        Files.write(out, json.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("written " + out);
    }

    static void dumpExceptionFile(StringBuilder json, Path testdata, String relativePath, String expectedSubstring) throws Exception {
        DSSDocument document = new FileDocument(testdata.resolve(relativePath).toFile());
        JWSDocumentAnalyzerFactory factory = new JWSDocumentAnalyzerFactory();
        String gotMessage;
        try {
            AbstractJWSDocumentAnalyzer analyzer = factory.create(document);
            // Force lazy JWS parsing (create() alone may defer it) the same way getSignatures()
            // does, so a fixture that only fails once parsing actually happens is still caught.
            analyzer.getSignatures();
            throw new IllegalStateException("expected " + relativePath + " to raise an exception, but it parsed cleanly");
        } catch (RuntimeException e) {
            StringBuilder chain = new StringBuilder();
            for (Throwable t = e; t != null; t = t.getCause()) {
                if (t.getMessage() != null) {
                    chain.append(t.getMessage()).append(" | ");
                }
            }
            gotMessage = chain.toString();
            if (!gotMessage.contains(expectedSubstring)) {
                throw new AssertionError("oracle's own expectation for " + relativePath +
                        " is stale: got exception chain [" + gotMessage + "], want it to contain [" + expectedSubstring + "]", e);
            }
        }
        json.append("    {\n");
        json.append("      \"path\": ").append(str(relativePath)).append(",\n");
        json.append("      \"expectedMessageSubstring\": ").append(str(expectedSubstring)).append("\n");
        json.append("    }");
    }

    static void dumpFile(StringBuilder json, Path testdata, String relativePath) throws Exception {
        Path filePath = testdata.resolve(relativePath);
        DSSDocument document = new FileDocument(filePath.toFile());

        JWSDocumentAnalyzerFactory factory = new JWSDocumentAnalyzerFactory();
        AbstractJWSDocumentAnalyzer analyzer = factory.create(document);
        CommonCertificateVerifier certificateVerifier = new CommonCertificateVerifier();
        analyzer.setCertificateVerifier(certificateVerifier);

        String detachedRelative = DETACHED_CONTENT.get(relativePath);
        if (detachedRelative != null) {
            DSSDocument detached = new FileDocument(testdata.resolve(detachedRelative).toFile());
            analyzer.setDetachedContents(Collections.singletonList(detached));
        }

        List<AdvancedSignature> signatures = analyzer.getSignatures();

        json.append("    {\n");
        json.append("      \"path\": ").append(str(relativePath)).append(",\n");
        json.append("      \"signatureCount\": ").append(signatures.size()).append(",\n");
        json.append("      \"signatures\": [\n");
        for (int i = 0; i < signatures.size(); i++) {
            JAdESSignature signature = (JAdESSignature) signatures.get(i);
            dumpSignature(json, signature, certificateVerifier);
            json.append(i == signatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("      ]\n");
        json.append("    }");
    }

    static void dumpSignature(StringBuilder json, JAdESSignature signature, CommonCertificateVerifier certificateVerifier) {
        // Counter signatures are not initialized by AbstractJWSDocumentAnalyzer.buildSignatures()
        // (only the top-level ones are); getDataFoundUpToLevel() requires it.
        signature.initBaselineRequirementsChecker(certificateVerifier);

        CertificateToken signingCertificate = signature.getSigningCertificateToken();
        Long signingTimeMillis = signature.getSigningTime() != null ? signature.getSigningTime().getTime() : null;

        // Force the cryptographic verification the same way DSS's own validation process does.
        SignatureCryptographicVerification verification = signature.getSignatureCryptographicVerification();

        Digest dtbsr = null;
        try {
            dtbsr = signature.getDataToBeSignedRepresentation();
        } catch (Exception e) {
            // JWS_SIGNING_INPUT reference validation genuinely absent in some malformed fixtures
        }

        SigDMechanism sigDMechanism = null;
        try {
            sigDMechanism = signature.getSigDMechanism();
        } catch (Exception e) {
            // no 'sigD' header, or an unrecognized mechanism URI
        }

        JWSSerializationType serializationType = signature.getJws().getJwsSerializationType();

        List<String> structureErrors = signature.getStructureValidationResult();

        json.append("        {\n");
        json.append("          \"signingCertificateFound\": ").append(signingCertificate != null).append(",\n");
        json.append("          \"signingCertificateSHA256\": ")
                .append(signingCertificate != null ? str(hex(DSSUtils.digest(DigestAlgorithm.SHA256, signingCertificate.getEncoded()))) : "null")
                .append(",\n");
        json.append("          \"claimedSigningTimeMillis\": ").append(signingTimeMillis != null ? signingTimeMillis.toString() : "null").append(",\n");
        json.append("          \"dataFoundUpToLevel\": ").append(str(signature.getDataFoundUpToLevel().name())).append(",\n");
        json.append("          \"isCounterSignature\": ").append(signature.isCounterSignature()).append(",\n");
        json.append("          \"referenceDataFound\": ").append(verification.isReferenceDataFound()).append(",\n");
        json.append("          \"referenceDataIntact\": ").append(verification.isReferenceDataIntact()).append(",\n");
        json.append("          \"signatureIntact\": ").append(verification.isSignatureIntact()).append(",\n");
        json.append("          \"dtbsrAlgorithm\": ").append(dtbsr != null ? str(dtbsr.getAlgorithm().name()) : "null").append(",\n");
        json.append("          \"dtbsrHex\": ").append(dtbsr != null ? str(hex(dtbsr.getValue())) : "null").append(",\n");
        json.append("          \"sigDMechanism\": ").append(sigDMechanism != null ? str(sigDMechanism.name()) : "null").append(",\n");
        json.append("          \"serializationType\": ").append(str(serializationType.name())).append(",\n");
        json.append("          \"structureValidationErrorCount\": ").append(structureErrors.size()).append(",\n");

        // Time-stamp inventory. Two columns here are deliberately identity-sensitive rather than
        // merely structural, because both caught a real defect that every other column in this
        // file was blind to (found by testdata/broadgen's whole-corpus differential run):
        //
        //  - "dssId": the DSS token identifier, which SignatureTimestampIdentifierBuilder derives
        //    from the token binaries PLUS the signature id, the carrying attribute's identifier,
        //    the attribute's ORDER among the signature properties, and the token's order within
        //    that attribute. A port whose getAttributeOrder lookup fails to find the attribute
        //    gets a null order and therefore a different id, with no other visible symptom.
        //
        //  - "timestampedReferences": the set of objects a time-stamp covers. JAdESTimestampSource
        //    overrides getSignatureTimestampReferences() to fold in getKeyInfoReferences(), so a
        //    port that misses that override produces signature time-stamps covering one
        //    certificate fewer than upstream's.
        dumpTimestamps(json, "contentTimestamps", signature.getContentTimestamps());
        dumpTimestamps(json, "signatureTimestamps", signature.getSignatureTimestamps());
        dumpTimestamps(json, "timestampsX1", signature.getTimestampsX1());
        dumpTimestamps(json, "timestampsX2", signature.getTimestampsX2());
        dumpTimestamps(json, "archiveTimestamps", signature.getArchiveTimestamps());

        List<AdvancedSignature> counterSignatures = signature.getCounterSignatures();
        json.append("          \"counterSignatureCount\": ").append(counterSignatures.size()).append(",\n");
        json.append("          \"counterSignatures\": [\n");
        for (int i = 0; i < counterSignatures.size(); i++) {
            dumpSignature(json, (JAdESSignature) counterSignatures.get(i), certificateVerifier);
            json.append(i == counterSignatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("          ]\n");
        json.append("        }");
    }

    /** Dumps one time-stamp bucket. See dumpSignature's comment for why "dssId" and
     *  "timestampedReferences" are here. References are sorted so the two implementations'
     *  insertion orders are not compared - only the covered set. */
    static void dumpTimestamps(StringBuilder json, String field, List<TimestampToken> timestamps) {
        json.append("          ").append(str(field)).append(": [\n");
        for (int i = 0; i < timestamps.size(); i++) {
            TimestampToken timestampToken = timestamps.get(i);
            Digest messageImprint = timestampToken.getMessageImprint();
            List<String> references = new ArrayList<>();
            for (TimestampedReference reference : timestampToken.getTimestampedReferences()) {
                references.add(reference.getCategory().name() + ":" + reference.getObjectId());
            }
            Collections.sort(references);
            json.append("            {\n");
            json.append("              \"type\": ").append(str(timestampToken.getTimeStampType().name())).append(",\n");
            json.append("              \"dssId\": ").append(str(timestampToken.getDSSIdAsString())).append(",\n");
            json.append("              \"generationTimeMillis\": ")
                    .append(timestampToken.getGenerationTime() != null ? Long.toString(timestampToken.getGenerationTime().getTime()) : "null").append(",\n");
            json.append("              \"messageImprintDataFound\": ").append(timestampToken.isMessageImprintDataFound()).append(",\n");
            json.append("              \"messageImprintDataIntact\": ").append(timestampToken.isMessageImprintDataIntact()).append(",\n");
            json.append("              \"messageImprintHex\": ")
                    .append(messageImprint != null && messageImprint.getValue() != null ? str(hex(messageImprint.getValue())) : "null").append(",\n");
            json.append("              \"signatureIntact\": ").append(timestampToken.isSignatureIntact()).append(",\n");
            json.append("              \"timestampedReferences\": ").append(str(String.join(",", references))).append("\n");
            json.append("            }").append(i == timestamps.size() - 1 ? "\n" : ",\n");
        }
        json.append("          ],\n");
    }

    static String str(String s) {
        StringBuilder sb = new StringBuilder("\"");
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"': sb.append("\\\""); break;
                case '\\': sb.append("\\\\"); break;
                case '\n': sb.append("\\n"); break;
                default: sb.append(c);
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

// Generates testdata/upstream-cross-validation.json: ground-truth facts about the XAdES
// signatures in testdata/upstream/, dumped straight from upstream DSS 6.5.RC1's own
// XMLDocumentAnalyzer / XAdESSignature. This is the golden file that
// xades_upstream_cross_validation_test.go compares its own parse of the same files against
// (direction "UPSTREAM -> GO" of the cross-validation harness, task #12, XAdES extension).
//
// For every file this dumps, per signature: whether a signing certificate was identified, its
// SHA-256 digest, the claimed signing time (epoch millis, or null), the level upstream detects
// (SignatureLevel#name(), e.g. "XAdES_BASELINE_LT"), whether it is a counter signature, the
// overall reference/signature cryptographic verification verdict, and - unlike the CAdES oracle,
// which anchors on the CMS SignerId - the full per-Reference breakdown from
// XAdESSignature#getReferenceValidations() (type, id, uri, found, intact), since that is XAdES's
// own natural identity anchor and the most direct probe of this port's XML canonicalization /
// digesting pipeline (exclusive vs. inclusive C14N, multiple digest algorithms, enveloped vs.
// enveloping vs. detached references, xpointer/manifest references).
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd $DSS_UPSTREAM_HOME  (your built upstream DSS 6.5.RC1 checkout)
//   mvn -q -o -pl dss-xades dependency:build-classpath -Dmdep.outputFile=/tmp/cp.txt -Dmdep.includeScope=test
//   CP="dss-xades/target/classes:$(cat /tmp/cp.txt)"
//   javac -cp "$CP" -d /tmp/xvaloracle CrossValidationOracle.java
//   java  -cp "$CP:/tmp/xvaloracle" CrossValidationOracle <testdata directory>
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.ReferenceValidation;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.signature.AdvancedSignature;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.model.signature.SignatureCryptographicVerification;
import eu.europa.esig.dss.xades.validation.XAdESSignature;
import eu.europa.esig.dss.xades.validation.XMLDocumentAnalyzer;

import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

public class CrossValidationOracle {

    /** Files dumped, relative to testdata/upstream/. Covers B/T/LT/LTA baseline levels, legacy
     *  (pre-baseline) BES/EPES/C/X/XL profiles, enveloped/enveloping packaging, counter-signatures,
     *  a manifest reference, and several non-default digest/canonicalization combinations (SHA-1,
     *  SHA-256, SHA-512, exclusive vs. inclusive C14N, brainpool ECC).
     *
     *  The second block below is the adversarial half of the corpus, added by the Phase 4d audit:
     *  signatures that upstream itself REJECTS, in every way a XAdES signature can be broken -
     *  tampered reference/KeyInfo/manifest content, forged counter signatures, signing-certificate
     *  references whose digest, serial number or IssuerSerial does not match, XML signature
     *  wrapping (the xsw/ family, including two elements sharing one Id), timestamps of the wrong
     *  type or version - plus structural variety the first block missed: XAdES 1.1.1 with an XPath
     *  Filter 2.0 transform, plain XMLDSig with no XAdES at all, EPES with and without an
     *  SPDocSpecification, multiple archive/signature timestamps on one signature, an embedded
     *  evidence record, and four further national plugtest profiles. These matter more than the
     *  valid ones: a port that quietly ACCEPTS what upstream rejects is the dangerous failure
     *  mode, and only a corpus of rejections can catch it. */
    static final String[] FILES = {
        "valid-xades.xml",
        "XAdESLTA.xml",
        "xades-lta-valid.xml",
        "xades-extended-t.xml",
        "xades-extended-xl.xml",
        "xades-extended-short-a.xml",
        "xades-x-level.xml",
        "xades-x-level-v2.xml",
        "xades-level-b-with-certvals.xml",
        "BaselineBWithCertificateValues.xml",
        "Signature-X-AT-1.xml",
        "Signature-X-BE_ECON-3.xml",
        "Signature-X-BG-1.xml",
        "Signature-X-CY-1.xml",
        "Signature-X-CZ_ICZ-1.xml",
        "Signature-X-ES-100.xml",
        "Signature-X-FR_NOT-3.xml",
        "Signature-X-HR_FIN-1.xml",
        "Signature-X-HU_POL-3.xml",
        "Signature-X-RO_TRA-4.xml",
        "Signature-X-SK_DIT-1.xml",
        "Signature-X-UK_ASC-2.xml",
        "doubleSignedTest.xml",
        "xades-with-counter-sig-type.xml",
        "xades-counter-signature-no-type.xml",
        "xades-500-references.xml",
        "sig_bundle.signed_detached.xml",
        "dss1811-multi-algo.xml",
        "xades-ecc-brainpool.xml",
        "xades-with-spUserNotice-refs.xml",
        "xades-tampered-reference.xml",
        "xades-tampered-keyinfo.xml",
        "xades-tampered-manifest.xml",
        "xades-fake-counter-signature.xml",
        "xades-sign-cert-v2-wrong-digest.xml",
        "xades-sign-cert-v1-wrong-serialnumber.xml",
        "xades-wrong-sign-cert-digest.xml",
        "xades-x-level-v2-wrong-sigAndRefsTst-type.xml",
        "xades-extended-epes.xml",
        "xades-sigPolicy-with-SPDocSpecification.xml",
        "xades-detached-with-object-type-ref.xml",
        "xades-no-keyinfo-sign-cert.xml",
        "xades-with-equivalent-certs.xml",
        "qes-xades111-filter2.xml",
        "xmldsig-only.xml",
        "Signature-X-PT-4.xml",
        "Signature-X-CZ_SEF-5.xml",
        "Signature-X-UK_ELD-4.xml",
        "xades-lta-with-additional-cert-in-keyinfo.xml",
        "xades-counter-signature-injected.xml",
        "signature_property_signed.xml",
        "signing-cert-multiple-refs-sig.xml",
        "XSW-enveloped-fake-content.xml",
        "XSW-enveloped-fake-signedProperties.xml",
        "XSW-enveloping-fake-manifest.xml",
        "XSW-enveloped-fake-content-two-same-id.xml",
        "xades-with-multiple-archivetimestamps.xml",
        "xades-with-multiple-signaturetimestamps.xml",
        "dss1770.xml",
        "X-E-ERS-basic.xml",
    };

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
        json.append("  ]\n");
        json.append("}\n");

        Path out = Paths.get(args[0]).resolve("upstream-cross-validation.json");
        Files.write(out, json.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("written " + out);
    }

    static void dumpFile(StringBuilder json, Path testdata, String relativePath) throws Exception {
        Path filePath = testdata.resolve(relativePath);
        DSSDocument document = new FileDocument(filePath.toFile());
        XMLDocumentAnalyzer analyzer = new XMLDocumentAnalyzer(document);
        CommonCertificateVerifier certificateVerifier = new CommonCertificateVerifier();
        analyzer.setCertificateVerifier(certificateVerifier);

        List<AdvancedSignature> signatures = analyzer.getSignatures();

        json.append("    {\n");
        json.append("      \"path\": ").append(str(relativePath)).append(",\n");
        json.append("      \"signatureCount\": ").append(signatures.size()).append(",\n");
        json.append("      \"signatures\": [\n");
        for (int i = 0; i < signatures.size(); i++) {
            XAdESSignature signature = (XAdESSignature) signatures.get(i);
            dumpSignature(json, signature, certificateVerifier);
            json.append(i == signatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("      ]\n");
        json.append("    }");
    }

    static void dumpSignature(StringBuilder json, XAdESSignature signature, CommonCertificateVerifier certificateVerifier) {
        // Counter signatures are not initialized by DefaultDocumentAnalyzer.buildSignatures() (only
        // the top-level ones are); getDataFoundUpToLevel() requires it.
        signature.initBaselineRequirementsChecker(certificateVerifier);
        CertificateToken signingCertificate = signature.getSigningCertificateToken();
        Long signingTimeMillis = signature.getSigningTime() != null ? signature.getSigningTime().getTime() : null;

        SignatureCryptographicVerification verification = signature.getSignatureCryptographicVerification();

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

        List<ReferenceValidation> referenceValidations = signature.getReferenceValidations();
        json.append("          \"referenceValidationCount\": ").append(referenceValidations.size()).append(",\n");
        json.append("          \"referenceValidations\": [\n");
        for (int i = 0; i < referenceValidations.size(); i++) {
            ReferenceValidation rv = referenceValidations.get(i);
            json.append("            {\n");
            json.append("              \"type\": ").append(str(rv.getType().name())).append(",\n");
            json.append("              \"id\": ").append(str(rv.getId() != null ? rv.getId() : "")).append(",\n");
            json.append("              \"uri\": ").append(str(rv.getUri() != null ? rv.getUri() : "")).append(",\n");
            json.append("              \"found\": ").append(rv.isFound()).append(",\n");
            json.append("              \"intact\": ").append(rv.isIntact()).append(",\n");
            // Dependent validations: the entries of a signed ds:Manifest a
            // ds:Reference[@Type=".../Manifest"] points at. Dumped because a DataObjectFormat
            // qualifying property is allowed to reference a manifest entry rather than a
            // ds:SignedInfo/ds:Reference, so a port that drops them silently downgrades an
            // otherwise valid signature - which is exactly the defect the Phase 4d audit found
            // on Signature-X-CZ_SEF-5.xml while these were NOT being dumped.
            List<ReferenceValidation> dependentValidations = rv.getDependentValidations();
            json.append("              \"dependentValidations\": [\n");
            for (int j = 0; j < dependentValidations.size(); j++) {
                ReferenceValidation dv = dependentValidations.get(j);
                json.append("                {\n");
                json.append("                  \"type\": ").append(str(dv.getType().name())).append(",\n");
                json.append("                  \"id\": ").append(str(dv.getId() != null ? dv.getId() : "")).append(",\n");
                json.append("                  \"uri\": ").append(str(dv.getUri() != null ? dv.getUri() : "")).append(",\n");
                json.append("                  \"found\": ").append(dv.isFound()).append(",\n");
                json.append("                  \"intact\": ").append(dv.isIntact()).append("\n");
                json.append("                }").append(j == dependentValidations.size() - 1 ? "\n" : ",\n");
            }
            json.append("              ]\n");
            json.append("            }").append(i == referenceValidations.size() - 1 ? "\n" : ",\n");
        }
        json.append("          ],\n");

        List<AdvancedSignature> counterSignatures = signature.getCounterSignatures();
        json.append("          \"counterSignatureCount\": ").append(counterSignatures.size()).append(",\n");
        json.append("          \"counterSignatures\": [\n");
        for (int i = 0; i < counterSignatures.size(); i++) {
            dumpSignature(json, (XAdESSignature) counterSignatures.get(i), certificateVerifier);
            json.append(i == counterSignatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("          ]\n");
        json.append("        }");
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

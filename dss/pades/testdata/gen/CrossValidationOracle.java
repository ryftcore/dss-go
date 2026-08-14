// Generates testdata/upstream-cross-validation.json: ground-truth facts about the PAdES
// signatures in testdata/upstream/, dumped straight from upstream DSS 6.5.RC1's own
// PDFDocumentAnalyzer / PAdESSignature. This is the golden file that
// pades_upstream_cross_validation_test.go compares its own parse of the same files against
// (direction "UPSTREAM -> GO" of the cross-validation harness, task #12, PAdES extension).
//
// PAdESSignature extends CAdESSignature (dss-cades -> dss-pades in Java, cades -> pades in this
// port), so most of what this dumps mirrors cades/testdata/gen/CrossValidationOracle.java
// exactly (signing certificate, claimed signing time, level, CMS SignerId, cryptographic
// verification, message-digest value) - the CMS SignerInfo remains PAdES's own natural identity
// anchor, same as CAdES. What is new here is the PDF-specific layer PAdESSignature.getPdfRevision()
// exposes: the /ByteRange actually signed, whether the incremental-update chain covers every byte
// of the previous revision (areAllOriginalBytesCovered - the PAdES analogue of XAdES's per-
// Reference "found/intact"), the /SubFilter and other PDF signature dictionary fields, the
// signature field name(s) the /Sig dictionary is attached to, and - the most PAdES-specific probe
// of all - PdfModificationDetection.areModificationsDetected() plus the secure/formFill/annotation/
// undefined change buckets PdfObjectModifications sorts every object-graph diff between PDF
// revisions into. That last one is what actually exercises this port's native PDF
// parser/incremental-update diffing engine (internal/pdf) against real, sometimes adversarial,
// incremental updates - several fixtures below are PDFs upstream itself flags as modified/spoofed
// after signing, and the golden JSON captures upstream's own (true) verdict for them.
//
// Also dumped at the top level: the document's detached timestamps (PDF doc-timestamps not
// covering any earlier /Sig revision), since a PDF can carry a bare DocTimeStamp with no
// signature at all.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-pades dependency:build-classpath -Dmdep.outputFile=/tmp/cp.txt -Dmdep.includeScope=test
//   PDFBOX=/root/.m2/repository/org/apache/pdfbox
//   CP="dss-pades/target/classes:dss-cms-object/target/classes:dss-pades-pdfbox/target/classes"
//   CP="$CP:$PDFBOX/pdfbox/3.0.7/pdfbox-3.0.7.jar:$PDFBOX/fontbox/3.0.7/fontbox-3.0.7.jar"
//   CP="$CP:$PDFBOX/pdfbox-io/3.0.7/pdfbox-io-3.0.7.jar:$(cat /tmp/cp.txt)"
//   javac -cp "$CP" -d /tmp/pvaloracle CrossValidationOracle.java
//   java  -cp "$CP:/tmp/pvaloracle" CrossValidationOracle <testdata directory>
//
// dss-cms-object has to be added by hand for the same reason the CAdES/XAdES oracles need it: it
// is the runtime CMS implementation dss-cms selects through its service loader, and dss-cades
// (which dss-pades depends on) declares neither of the two implementations.
//
// dss-pades-pdfbox and the three org.apache.pdfbox jars have to be added by hand too, and for a
// closely related reason: dss-pades declares no PDF backend either (upstream ships two,
// dss-pades-pdfbox and dss-pades-openpdf, and ServiceLoaderPdfObjFactory picks whichever is on
// the classpath). Without one, this oracle dies on the very first fixture with "No
// implementation found for IPdfObjFactory in classpath". pdfbox is the backend this port's
// golden file is generated against and the one dss-pades-pdfbox's own tests run on, so it is
// also the reference the native internal/pdf engine's PDF-layer behaviour is compared to.
// Its jars are pulled from the local Maven repository directly rather than through
// `mvn -pl dss-pades-pdfbox dependency:build-classpath`, which cannot run offline here: that
// module declares a test-scope dependency on dss-pdfa, which needs org.verapdf artifacts from
// Maven Central (the same obstacle pades_downstream_cross_validation_test.go documents).
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.signature.SignatureCryptographicVerification;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.pades.validation.ByteRange;
import eu.europa.esig.dss.pades.validation.PAdESSignature;
import eu.europa.esig.dss.pades.validation.PDFDocumentAnalyzer;
import eu.europa.esig.dss.pades.validation.PdfSignatureDictionary;
import eu.europa.esig.dss.pades.validation.PdfSignatureField;
import eu.europa.esig.dss.pdf.PdfCMSRevision;
import eu.europa.esig.dss.pdf.modifications.PdfModificationDetection;
import eu.europa.esig.dss.pdf.modifications.PdfObjectModifications;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.signature.AdvancedSignature;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.tsp.TimestampToken;
import org.bouncycastle.cms.SignerId;
import org.bouncycastle.cms.SignerInformation;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

public class CrossValidationOracle {

    /** Files dumped, relative to testdata/upstream/. Covers PAdES-B/T/LT/LTA baseline levels,
     *  legacy (pre-baseline) BES/EPES/PKCS7 profiles, VRI and document timestamps, multiple
     *  signatures/revisions in one PDF, several national plugtest profiles, and - the adversarial
     *  half of the corpus, deliberately more than half of it - PDFs upstream itself detects as
     *  modified after signing: an added annotation, a spoofed /ByteRange gap, a replaced /Reason,
     *  a byte-range overlap attack, and a fully rewritten page. A port that quietly ACCEPTS what
     *  upstream flags as tampered is the dangerous failure mode, and only a corpus of genuine
     *  incremental-update attacks can catch it. */
    static final String[] FILES = {
        "validation/pades-bes.pdf",
        "validation/pades-epes.pdf",
        "validation/pades-not-epes.pdf",
        "validation/doc-firmado.pdf",
        "validation/doc-firmado-T.pdf",
        "validation/doc-firmado-LT.pdf",
        "validation/PAdES-LT.pdf",
        "validation/PAdES-LTA.pdf",
        "validation/Test.signed_Certipost-2048-SHA512.extended-LTA.pdf",
        "validation/pades-ltv.pdf",
        "validation/pades-ltv-with-reason.pdf",
        "validation/pades-lt-extended-dss.pdf",
        "validation/pades-t-level-extended.pdf",
        "validation/pades-5-signatures-and-1-document-timestamp.pdf",
        "validation/pades-two-sig-copied-tst.pdf",
        "validation/pdf-with-vri-timestamp.pdf",
        "validation/test-with-vri.pdf",
        "validation/timestamped_and_signed.pdf",
        "validation/timestamped-fields.pdf",
        "validation/belgian_pki_multiple_ocsps_lt.pdf",
        "validation/Signature-P-HU_POL-3.pdf",
        "validation/Signature-P-DE_SCI-4.pdf",
        "validation/adbe_crl_signed.pdf",
        "validation/adbe_ocsp_signed.pdf",
        "validation/pkcs7.pdf",
        "validation/pades-signed-annot-added.pdf",
        "validation/pdf-spoofing-attack.pdf",
        "validation/pdf-byterange-overlap.pdf",
        "validation/pades-spoofing-replaced-reason.pdf",
        "validation/pades-multiple-pages-annots-overlap.pdf",
        "validation/pdf-removed-pages.pdf",
        // Regression fixtures for the per-signature scoping of validation data in a multi-
        // revision document, added after a mutation check showed the corpus above could not
        // detect a broken PdfDssDict equality nor a document-timestamp source that never
        // reaches PAdESTimestampSource's own getDocumentTimestamps(). Every one of these
        // carries a signature upstream reports at a LOWER level than the document's latest
        // revision would suggest - PAdES-BASELINE-T next to an LT/LTA sibling - which is
        // exactly what over-sharing a /DSS dictionary, or losing a /DocTimeStamp on the way
        // into the LT-level revocation-presence check, silently turns into a false LT/LTA.
        "validation/Signature-P-SK-6.pdf",
        "validation/pades-lt-extended-abde.pdf",
        "validation/pades-t-duplicated-doctst.pdf",
        "validation/pades3_Baseline_B.pdf",
        // /Name written in PDFDocEncoding's Central-European half ("Martin Petrzela", with a
        // z-caron at byte 0x9E): guards the PDFDocEncoding table against the Latin-1
        // approximation it used to be (internal/pdf/pdfdocencoding.go).
        "validation/pades-ocsp-archiveCutOff-invalid.pdf",
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
        PDFDocumentAnalyzer analyzer = new PDFDocumentAnalyzer(document);
        CommonCertificateVerifier certificateVerifier = new CommonCertificateVerifier();
        analyzer.setCertificateVerifier(certificateVerifier);

        // validate() is what actually runs PDFDocumentAnalyzer's postProcessing/
        // timestampPostProcessing (AnalyzePdfModifications) - getSignatures() alone (protected
        // getAllSignatures() is the one PDFDocumentAnalyzer overrides to call postProcessing, and
        // it is not reachable from outside the class) never populates
        // PdfModificationDetection at all, which would make every "modificationsDetected" dump
        // below silently false regardless of what the fixture actually contains. validate()
        // mutates the same AdvancedSignature objects getSignatures() below then returns
        // (cached), so this is the same object graph a full validateDocument() would produce.
        analyzer.validate();
        List<AdvancedSignature> signatures = analyzer.getSignatures();
        List<TimestampToken> detachedTimestamps = analyzer.getDetachedTimestamps();

        json.append("    {\n");
        json.append("      \"path\": ").append(str(relativePath)).append(",\n");
        json.append("      \"signatureCount\": ").append(signatures.size()).append(",\n");
        json.append("      \"signatures\": [\n");
        for (int i = 0; i < signatures.size(); i++) {
            PAdESSignature signature = (PAdESSignature) signatures.get(i);
            dumpSignature(json, signature, certificateVerifier);
            json.append(i == signatures.size() - 1 ? "\n" : ",\n");
        }
        json.append("      ],\n");

        json.append("      \"detachedTimestampCount\": ").append(detachedTimestamps.size()).append(",\n");
        json.append("      \"detachedTimestamps\": [\n");
        for (int i = 0; i < detachedTimestamps.size(); i++) {
            dumpTimestamp(json, detachedTimestamps.get(i));
            json.append(i == detachedTimestamps.size() - 1 ? "\n" : ",\n");
        }
        json.append("      ]\n");
        json.append("    }");
    }

    static void dumpSignature(StringBuilder json, PAdESSignature signature, CommonCertificateVerifier certificateVerifier) {
        // Not initialized by PDFDocumentAnalyzer.buildSignatures() by default the way a nested
        // counter signature would need it elsewhere - PAdES has no counter signatures at all
        // (getCounterSignatures() is a hard-coded empty list) - but getDataFoundUpToLevel()
        // still requires the baseline requirements checker, same as CAdES/XAdES.
        signature.initBaselineRequirementsChecker(certificateVerifier);
        CertificateToken signingCertificate = signature.getSigningCertificateToken();
        Long signingTimeMillis = signature.getSigningTime() != null ? signature.getSigningTime().getTime() : null;

        SignerInformation signerInformation = signature.getSignerInformation();
        String signerInfoDigest;
        try {
            signerInfoDigest = hex(DSSUtils.digest(DigestAlgorithm.SHA256,
                    signerInformation.toASN1Structure().getEncoded("DER")));
        } catch (Exception e) {
            signerInfoDigest = null;
        }

        SignerId sid = signerInformation.getSID();
        String signerIdIssuerSerial = sid.getIssuer() != null && sid.getSerialNumber() != null
                ? sid.getIssuer().toString() + "#" + sid.getSerialNumber().toString()
                : null;
        String signerIdSki = sid.getSubjectKeyIdentifier() != null ? hex(sid.getSubjectKeyIdentifier()) : null;

        SignatureCryptographicVerification verification = signature.getSignatureCryptographicVerification();

        byte[] messageDigestValue = null;
        try {
            messageDigestValue = signature.getMessageDigestValue();
        } catch (Exception e) {
            // no message-digest attribute at all in some legacy PKCS7 profiles
        }

        json.append("        {\n");
        json.append("          \"signingCertificateFound\": ").append(signingCertificate != null).append(",\n");
        json.append("          \"signingCertificateSHA256\": ")
                .append(signingCertificate != null ? str(hex(DSSUtils.digest(DigestAlgorithm.SHA256, signingCertificate.getEncoded()))) : "null")
                .append(",\n");
        json.append("          \"claimedSigningTimeMillis\": ").append(signingTimeMillis != null ? signingTimeMillis.toString() : "null").append(",\n");
        json.append("          \"dataFoundUpToLevel\": ").append(str(signature.getDataFoundUpToLevel().name())).append(",\n");
        json.append("          \"signerInformationDigestSHA256\": ").append(signerInfoDigest != null ? str(signerInfoDigest) : "null").append(",\n");
        json.append("          \"signerIdIssuerSerial\": ").append(signerIdIssuerSerial != null ? str(signerIdIssuerSerial) : "null").append(",\n");
        json.append("          \"signerIdSubjectKeyIdentifier\": ").append(signerIdSki != null ? str(signerIdSki) : "null").append(",\n");
        json.append("          \"isCounterSignature\": ").append(signature.isCounterSignature()).append(",\n");
        json.append("          \"referenceDataFound\": ").append(verification.isReferenceDataFound()).append(",\n");
        json.append("          \"referenceDataIntact\": ").append(verification.isReferenceDataIntact()).append(",\n");
        json.append("          \"signatureIntact\": ").append(verification.isSignatureIntact()).append(",\n");
        json.append("          \"messageDigestValueHex\": ").append(messageDigestValue != null ? str(hex(messageDigestValue)) : "null").append(",\n");
        json.append("          \"counterSignatureCount\": ").append(signature.getCounterSignatures().size()).append(",\n");

        // --- PDF-specific layer: PAdESSignature.getPdfRevision() ---
        PdfCMSRevision pdfRevision = signature.getPdfRevision();
        PdfSignatureDictionary sigDict = pdfRevision.getPdfSigDictInfo();
        ByteRange byteRange = pdfRevision.getByteRange();

        json.append("          \"subFilter\": ").append(str(sigDict.getSubFilter())).append(",\n");
        json.append("          \"filter\": ").append(sigDict.getFilter() != null ? str(sigDict.getFilter()) : "null").append(",\n");
        json.append("          \"signerName\": ").append(sigDict.getSignerName() != null ? str(sigDict.getSignerName()) : "null").append(",\n");
        json.append("          \"reason\": ").append(sigDict.getReason() != null ? str(sigDict.getReason()) : "null").append(",\n");
        json.append("          \"location\": ").append(sigDict.getLocation() != null ? str(sigDict.getLocation()) : "null").append(",\n");
        json.append("          \"docMDP\": ").append(sigDict.getDocMDP() != null ? str(sigDict.getDocMDP().name()) : "null").append(",\n");
        json.append("          \"byteRange\": [").append(byteRange.getFirstPartStart()).append(", ")
                .append(byteRange.getFirstPartEnd()).append(", ").append(byteRange.getSecondPartStart()).append(", ")
                .append(byteRange.getSecondPartEnd()).append("],\n");
        json.append("          \"areAllOriginalBytesCovered\": ").append(pdfRevision.areAllOriginalBytesCovered()).append(",\n");

        List<PdfSignatureField> fields = pdfRevision.getFields();
        json.append("          \"fieldNames\": [");
        for (int i = 0; i < fields.size(); i++) {
            json.append(str(fields.get(i).getFieldName()));
            if (i != fields.size() - 1) {
                json.append(", ");
            }
        }
        json.append("],\n");

        json.append("          \"documentTimestampCount\": ").append(signature.getDocumentTimestamps().size()).append(",\n");
        json.append("          \"vriTimestampCount\": ").append(signature.getVRITimestamps().size()).append(",\n");

        PdfModificationDetection modificationDetection = pdfRevision.getModificationDetection();
        boolean modificationsDetected = modificationDetection != null && modificationDetection.areModificationsDetected();
        json.append("          \"pdfModificationsDetected\": ").append(modificationsDetected).append(",\n");
        if (modificationDetection != null) {
            PdfObjectModifications objectModifications = modificationDetection.getObjectModifications();
            json.append("          \"secureChangeCount\": ").append(objectModifications.getSecureChanges().size()).append(",\n");
            json.append("          \"formFillChangeCount\": ").append(objectModifications.getFormFillInAndSignatureCreationChanges().size()).append(",\n");
            json.append("          \"annotChangeCount\": ").append(objectModifications.getAnnotCreationChanges().size()).append(",\n");
            json.append("          \"undefinedChangeCount\": ").append(objectModifications.getUndefinedChanges().size()).append(",\n");
        } else {
            json.append("          \"secureChangeCount\": 0,\n");
            json.append("          \"formFillChangeCount\": 0,\n");
            json.append("          \"annotChangeCount\": 0,\n");
            json.append("          \"undefinedChangeCount\": 0,\n");
        }

        // DSS-Id stability is asserted Go-side only (calling ID() twice); nothing to dump here.
        json.append("          \"id\": ").append(str(signature.getId())).append("\n");
        json.append("        }");
    }

    static void dumpTimestamp(StringBuilder json, TimestampToken timestamp) {
        json.append("        {\n");
        json.append("          \"type\": ").append(str(timestamp.getTimeStampType().name())).append(",\n");
        json.append("          \"messageImprintDataFound\": ").append(timestamp.isMessageImprintDataFound()).append(",\n");
        json.append("          \"messageImprintDataIntact\": ").append(timestamp.isMessageImprintDataIntact()).append(",\n");
        json.append("          \"signatureIntact\": ").append(timestamp.isSignatureIntact()).append(",\n");
        json.append("          \"genTimeMillis\": ").append(timestamp.getGenerationTime() != null ? timestamp.getGenerationTime().getTime() : "null").append("\n");
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

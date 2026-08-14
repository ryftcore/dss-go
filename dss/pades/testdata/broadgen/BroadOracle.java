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

public class BroadOracle {

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        java.util.List<Path> files = new java.util.ArrayList<>();
        java.nio.file.Files.walk(root).filter(p -> p.toString().toLowerCase().endsWith(".pdf")).forEach(files::add);
        java.util.Collections.sort(files);
        StringBuilder json = new StringBuilder();
        json.append("{\n  \"files\": [\n");
        boolean first = true;
        for (Path p : files) {
            String rel = root.relativize(p).toString().replace('\\','/');
            StringBuilder one = new StringBuilder();
            try {
                dumpFile(one, root, rel);
            } catch (Throwable t) {
                one.setLength(0);
                one.append("    {\n      \"path\": ").append(str(rel)).append(",\n      \"error\": ")
                   .append(str(t.getClass().getName() + ": " + String.valueOf(t.getMessage()))).append("\n    }");
            }
            if (!first) json.append(",\n");
            first = false;
            json.append(one);
        }
        json.append("\n  ]\n}\n");
        java.nio.file.Files.write(Paths.get(args[1]), json.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("written " + args[1] + " (" + files.size() + " files)");
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

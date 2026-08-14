// UPSTREAM -> GO -> UPSTREAM round-trip, step 1 of 3 (S7_BRIEF.md container semantics: "Java-built
// containers extended ... by Go and re-validated by Java").
//
// The plain GO -> UPSTREAM direction (main.go + CrossGenValidator.java) proves upstream accepts a
// container this port BUILT. It says nothing about a container this port only TOUCHED: extension
// rewrites an existing signed container in place - it re-encodes the CMS signature, copies every
// other entry across, and rebuilds the zip - so a port that silently normalised entry metadata,
// dropped an entry, or re-serialised the signature wrongly would still pass the build-only test.
//
// This program closes that hole by producing the input for it: ASiC-S and ASiC-E containers signed
// at CAdES-BASELINE-B by UPSTREAM DSS 6.5.RC1's own ASiCWithCAdESService, using the same
// signer_rsa.p12 test key the Go generator uses (a currently-valid certificate - the corpus
// fixtures' signing certificates all expired years ago, and augmenting an expired signature is
// refused by both implementations alike). asic_downstream_cross_validation_test.go then has the Go
// port extend each of them to CAdES-BASELINE-T and hands the result back to CrossGenValidator.java.
//
// Usage: java JavaBuiltGenerator <output directory> <signer_rsa.p12 path>
//
// Writes <outdir>/java-asics-cades-b.scs and <outdir>/java-asice-cades-b.sce.
import eu.europa.esig.dss.asic.cades.ASiCWithCAdESSignatureParameters;
import eu.europa.esig.dss.asic.cades.signature.ASiCWithCAdESService;
import eu.europa.esig.dss.enumerations.ASiCContainerType;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.SignatureValue;
import eu.europa.esig.dss.model.ToBeSigned;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
import eu.europa.esig.dss.token.KeyStoreSignatureTokenConnection;
import eu.europa.esig.dss.token.Pkcs12SignatureToken;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.List;

public final class JavaBuiltGenerator {

    public static void main(String[] args) throws Exception {
        File outDir = Paths.get(args[0]).toFile();
        if (!outDir.isDirectory() && !outDir.mkdirs()) {
            throw new IllegalStateException("cannot create " + outDir);
        }
        File keyStore = Paths.get(args[1]).toFile();

        DSSDocument single = new InMemoryDocument(
                ("DSS Go port round-trip sample - built by UPSTREAM DSS, extended by the Go port, "
                        + "re-validated by UPSTREAM DSS (ASiC-S).").getBytes(StandardCharsets.UTF_8),
                "sample.txt");
        DSSDocument multiA = new InMemoryDocument(
                ("DSS Go port round-trip sample - built by UPSTREAM DSS (ASiC-E, entry A).")
                        .getBytes(StandardCharsets.UTF_8), "sample-a.txt");
        DSSDocument multiB = new InMemoryDocument(
                ("DSS Go port round-trip sample - built by UPSTREAM DSS (ASiC-E, entry B).")
                        .getBytes(StandardCharsets.UTF_8), "sample-b.txt");

        sign(keyStore, Arrays.asList(single), ASiCContainerType.ASiC_S,
                new File(outDir, "java-asics-cades-b.scs"));
        sign(keyStore, Arrays.asList(multiA, multiB), ASiCContainerType.ASiC_E,
                new File(outDir, "java-asice-cades-b.sce"));
        System.out.println("GENERATED OK");
    }

    private static void sign(File keyStore, List<DSSDocument> documents,
                             ASiCContainerType containerType, File target) throws Exception {
        try (KeyStoreSignatureTokenConnection token =
                     new Pkcs12SignatureToken(keyStore, new java.security.KeyStore.PasswordProtection(
                             "testpassword".toCharArray()))) {
            DSSPrivateKeyEntry key = token.getKeys().get(0);

            ASiCWithCAdESSignatureParameters parameters = new ASiCWithCAdESSignatureParameters();
            parameters.setSignatureLevel(SignatureLevel.CAdES_BASELINE_B);
            parameters.setDigestAlgorithm(DigestAlgorithm.SHA256);
            parameters.setSigningCertificate(key.getCertificate());
            parameters.setCertificateChain(key.getCertificateChain());
            parameters.aSiC().setContainerType(containerType);

            ASiCWithCAdESService service = new ASiCWithCAdESService(new CommonCertificateVerifier());
            ToBeSigned dataToSign = service.getDataToSign(documents, parameters);
            SignatureValue signatureValue = token.sign(dataToSign, parameters.getDigestAlgorithm(), key);
            DSSDocument signed = service.signDocument(documents, parameters, signatureValue);
            signed.save(target.getAbsolutePath());
        }
    }
}

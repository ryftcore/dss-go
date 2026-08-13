// The upstream half of the CAdES-B byte-exactness harness (cades_byte_exactness_test.go): signs
// the same fixed payload, with the same key store, signing certificate and chain, signing time
// and digest algorithm as main.go (bytecmp), using upstream DSS 6.5.RC1's own CAdESService. The
// Go test then compares the two outputs byte for byte - upstream's bytes are the ground truth,
// nothing here is derived from the port.
//
// Usage: java ByteExactnessFixtures <keystore.p12> <password> <signing year> <outdir> <prefix>
//
// Writes <outdir>/<prefix>-datatosign.bin and <outdir>/<prefix>.p7m, exactly as its Go
// counterpart does.
import eu.europa.esig.dss.cades.CAdESSignatureParameters;
import eu.europa.esig.dss.cades.signature.CAdESService;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.enumerations.SignaturePackaging;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.SignatureValue;
import eu.europa.esig.dss.model.ToBeSigned;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.token.DSSPrivateKeyEntry;
import eu.europa.esig.dss.token.Pkcs12SignatureToken;

import java.io.FileOutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyStore.PasswordProtection;
import java.util.Calendar;
import java.util.List;
import java.util.TimeZone;

public class ByteExactnessFixtures {

    /** The fixed payload, byte-identical to main.go's `content`. */
    static final byte[] CONTENT = "byte-exactness probe content".getBytes(StandardCharsets.UTF_8);

    public static void main(String[] args) throws Exception {
        if (args.length != 5) {
            System.err.println("usage: ByteExactnessFixtures <keystore.p12> <password> <signing year> <outdir> <prefix>");
            System.exit(2);
        }
        String keyStorePath = args[0];
        String password = args[1];
        int year = Integer.parseInt(args[2]);
        Path outDir = Paths.get(args[3]);
        String prefix = args[4];

        try (Pkcs12SignatureToken token = new Pkcs12SignatureToken(Paths.get(keyStorePath).toFile(),
                new PasswordProtection(password.toCharArray()))) {
            List<DSSPrivateKeyEntry> keys = token.getKeys();
            DSSPrivateKeyEntry entry = keys.get(0);

            Calendar calendar = Calendar.getInstance(TimeZone.getTimeZone("UTC"));
            calendar.clear();
            calendar.set(year, Calendar.JANUARY, 15, 10, 30, 45);

            CAdESSignatureParameters parameters = new CAdESSignatureParameters();
            parameters.setSignatureLevel(SignatureLevel.CAdES_BASELINE_B);
            parameters.setSignaturePackaging(SignaturePackaging.ENVELOPING);
            parameters.setDigestAlgorithm(DigestAlgorithm.SHA256);
            parameters.setSigningCertificate(entry.getCertificate());
            parameters.setCertificateChain(entry.getCertificateChain());
            parameters.bLevel().setSigningDate(calendar.getTime());

            CAdESService service = new CAdESService(new CommonCertificateVerifier());
            DSSDocument document = new InMemoryDocument(CONTENT, "probe.bin");

            ToBeSigned dataToSign = service.getDataToSign(document, parameters);
            Files.write(outDir.resolve(prefix + "-datatosign.bin"), dataToSign.getBytes());

            SignatureValue signatureValue = token.sign(dataToSign, parameters.getDigestAlgorithm(), entry);
            DSSDocument signed = service.signDocument(document, parameters, signatureValue);
            try (FileOutputStream os = new FileOutputStream(outDir.resolve(prefix + ".p7m").toFile())) {
                signed.writeTo(os);
            }
            System.out.println("wrote " + prefix + "-datatosign.bin and " + prefix + ".p7m");
        }
    }
}

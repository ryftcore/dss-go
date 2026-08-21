// GO -> UPSTREAM cross-validation (task #12, ASiC-with-CAdES extension): loads each container
// main.go (crossgen) produced with upstream DSS 6.5.RC1's own SignedDocumentValidator/
// SignedDocumentDiagnosticDataBuilder and asserts the container type is recognized, the
// signature is cryptographically intact, the signing certificate was identified, and the
// expected baseline level was recognized. This is the actual compatibility proof of the GO ->
// UPSTREAM direction: upstream DSS accepting what the Go port produced, not another Go-side
// self-check.
//
// Modeled on cades/testdata/crossgen/CrossGenValidator.java and
// xades/testdata/crossgen/CrossGenValidator.java (same three signature-level assertions), with
// one addition: the expected ASiCContainerType (ASiC-S vs ASiC-E) is checked too, since that is
// the property this harness extension is actually about.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1, dss-validation AND
// dss-asic-cades included (dss-asic-cades registers ASiCContainerWithCAdESValidatorFactory as a
// SignedDocumentValidator.fromDocument ServiceLoader provider - without it on the classpath,
// fromDocument() throws "Document format not recognized/handled" for a .scs/.sce file):
//
//   cd $DSS_UPSTREAM_HOME  (your built upstream DSS 6.5.RC1 checkout)
//   mvn -q -o -pl dss-validation dependency:build-classpath -Dmdep.outputFile=/tmp/valcp.txt -Dmdep.includeScope=runtime
//   mvn -q -o -pl dss-asic-cades dependency:build-classpath -Dmdep.outputFile=/tmp/asiccp.txt -Dmdep.includeScope=runtime
//   CP="dss-validation/target/classes:dss-asic-cades/target/classes:dss-asic-common/target/classes:dss-cades/target/classes:dss-cms/target/classes:dss-cms-object/target/classes:dss-document/target/classes:$(cat /tmp/valcp.txt):$(cat /tmp/asiccp.txt)"
//   javac -cp "$CP" -d /tmp/crossgenval CrossGenValidator.java
//   java  -cp "$CP:/tmp/crossgenval" CrossGenValidator <fixtures dir> \
//       asics-cades-b.scs:ASiC_S:CAdES_BASELINE_B \
//       asics-cades-t.scs:ASiC_S:CAdES_BASELINE_T \
//       asice-cades-b.sce:ASiC_E:CAdES_BASELINE_B \
//       asice-cades-t.sce:ASiC_E:CAdES_BASELINE_T
//
// Each fixture argument is "<file>:<expectedContainerType>:<expectedLevel>". Exits 0 and prints
// "ALL OK" when every fixture passes all assertions; otherwise prints the specific failure for
// each fixture and exits 1 - asic_downstream_cross_validation_test.go greps stdout for exactly
// those two markers, so a genuine rejection surfaces as a Go test failure rather than being
// swallowed.
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.enumerations.ASiCContainerType;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.SignedDocumentValidator;
import eu.europa.esig.dss.validation.reports.Reports;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

public class CrossGenValidator {

    public static void main(String[] args) throws Exception {
        if (args.length < 2) {
            System.err.println("usage: CrossGenValidator <fixtures dir> <file>:<expectedContainerType>:<expectedLevel> ...");
            System.exit(2);
        }
        Path dir = Paths.get(args[0]);

        boolean allOk = true;
        for (int i = 1; i < args.length; i++) {
            String[] parts = args[i].split(":");
            String fileName = parts[0];
            ASiCContainerType expectedContainerType = ASiCContainerType.valueOf(parts[1]);
            SignatureLevel expectedLevel = SignatureLevel.valueOf(parts[2]);

            try {
                checkFixture(dir, fileName, expectedContainerType, expectedLevel);
                System.out.println("OK " + fileName);
            } catch (AssertionError | RuntimeException e) {
                allOk = false;
                System.out.println("FAIL " + fileName + ": " + e.getMessage());
            }
        }

        if (allOk) {
            System.out.println("ALL OK");
        } else {
            System.out.println("SOME FAILED");
            System.exit(1);
        }
    }

    static void checkFixture(Path dir, String fileName, ASiCContainerType expectedContainerType,
                              SignatureLevel expectedLevel) throws Exception {
        DSSDocument document = new FileDocument(dir.resolve(fileName).toFile());

        SignedDocumentValidator validator = SignedDocumentValidator.fromDocument(document);
        validator.setCertificateVerifier(new CommonCertificateVerifier());

        Reports reports = validator.validateDocument();
        DiagnosticData diagnosticData = reports.getDiagnosticData();

        // (1) container type recognized
        ASiCContainerType actualContainerType = diagnosticData.getContainerType();
        if (actualContainerType != expectedContainerType) {
            throw new AssertionError("container type = " + actualContainerType + ", want " + expectedContainerType);
        }

        List<String> signatureIds = diagnosticData.getSignatureIdList();
        if (signatureIds.size() != 1) {
            throw new AssertionError("expected exactly 1 signature, found " + signatureIds.size());
        }
        SignatureWrapper signature = diagnosticData.getSignatureById(signatureIds.get(0));

        // (2) signature intact
        if (!signature.isSignatureIntact()) {
            throw new AssertionError("signature is not cryptographically intact");
        }
        if (!signature.isSignatureValid()) {
            throw new AssertionError("signature is not fully valid (reference data found/intact + signature intact)");
        }

        // (3) signing certificate identified
        if (signature.getSigningCertificate() == null) {
            throw new AssertionError("signing certificate was not identified");
        }

        // (4) level recognized
        SignatureLevel actualLevel = signature.getSignatureFormat();
        if (actualLevel != expectedLevel) {
            throw new AssertionError("signature level = " + actualLevel + ", want " + expectedLevel);
        }

        // (5) for a -T signature, the RFC 3161 token the Go port's own TSP source produced must
        // itself be cryptographically verified by upstream, not merely detected.
        if (expectedLevel.name().endsWith("_T")) {
            checkSignatureTimestamp(signature);
        }
    }

    /** Asserts the signature carries exactly one signature-timestamp and that upstream found it
     *  intact: message-imprint present and matching, and the token's own signature valid. */
    static void checkSignatureTimestamp(SignatureWrapper signature) {
        List<TimestampWrapper> timestamps = signature.getTimestampListByType(TimestampType.SIGNATURE_TIMESTAMP);
        if (timestamps.size() != 1) {
            throw new AssertionError("expected exactly 1 signature timestamp, found " + timestamps.size());
        }
        TimestampWrapper timestamp = timestamps.get(0);
        if (!timestamp.isMessageImprintDataFound()) {
            throw new AssertionError("timestamp message-imprint data was not found");
        }
        if (!timestamp.isMessageImprintDataIntact()) {
            throw new AssertionError("timestamp message-imprint does not match the signature it covers");
        }
        if (!timestamp.isSignatureIntact()) {
            throw new AssertionError("timestamp token signature is not cryptographically intact");
        }
        if (!timestamp.isSignatureValid()) {
            throw new AssertionError("timestamp token signature is not valid");
        }
    }
}

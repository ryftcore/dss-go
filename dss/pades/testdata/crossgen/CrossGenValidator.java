// GO -> UPSTREAM cross-validation (task #12, PAdES extension): loads each PAdES file main.go
// (crossgen) produced with upstream DSS 6.5.RC1's own
// SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and asserts the signature is
// cryptographically intact, the signing certificate was identified, and the expected baseline
// level was recognized. This is the actual compatibility proof of the GO -> UPSTREAM direction:
// upstream DSS accepting what the Go port produced, not another Go-side self-check. Mirrors
// dss/cades and dss/xades's own CrossGenValidator.java precedent.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1, dss-validation and a PDF backend
// (dss-pades-pdfbox) included:
//
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-pades-pdfbox compile
//   mvn -q -o -pl dss-validation dependency:build-classpath -Dmdep.outputFile=/tmp/valcp.txt -Dmdep.includeScope=runtime
//   PDFJARS="$(find ~/.m2 -path '*pdfbox*3.0.7*.jar' ! -name '*sources*' ! -name '*javadoc*' | tr '\n' ':')"
//   CP="dss-validation/target/classes:dss-pades/target/classes:dss-pades-pdfbox/target/classes:dss-cades/target/classes:dss-cms-object/target/classes:$(cat /tmp/valcp.txt):$PDFJARS"
//   javac -cp "$CP" -d /tmp/crossgenval CrossGenValidator.java
//   java  -cp "$CP:/tmp/crossgenval" CrossGenValidator <fixtures dir> \
//       EmptyPage-b.pdf:PAdES-BASELINE-B \
//       EmptyPage-t.pdf:PAdES-BASELINE-T \
//       ...
//
// Each fixture argument is "<file>:<expectedLevel>". Exits 0 and prints "ALL OK" when every
// fixture passes all assertions; otherwise prints the specific failure for each fixture and exits
// 1 - pades_downstream_cross_validation_test.go greps stdout for exactly those two markers, so a
// genuine rejection surfaces as a Go test failure rather than being swallowed.
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
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
            System.err.println("usage: CrossGenValidator <fixtures dir> <file>:<expectedLevel> ...");
            System.exit(2);
        }
        Path dir = Paths.get(args[0]);

        boolean allOk = true;
        for (int i = 1; i < args.length; i++) {
            String[] parts = args[i].split(":");
            String fileName = parts[0];
            SignatureLevel expectedLevel = SignatureLevel.valueOf(parts[1]);

            try {
                checkFixture(dir, fileName, expectedLevel);
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

    static void checkFixture(Path dir, String fileName, SignatureLevel expectedLevel) throws Exception {
        DSSDocument document = new FileDocument(dir.resolve(fileName).toFile());

        SignedDocumentValidator validator = SignedDocumentValidator.fromDocument(document);
        validator.setCertificateVerifier(new CommonCertificateVerifier());

        Reports reports = validator.validateDocument();
        DiagnosticData diagnosticData = reports.getDiagnosticData();

        List<String> signatureIds = diagnosticData.getSignatureIdList();
        if (signatureIds.size() != 1) {
            throw new AssertionError("expected exactly 1 signature, found " + signatureIds.size());
        }
        SignatureWrapper signature = diagnosticData.getSignatureById(signatureIds.get(0));

        // (1) signature intact
        if (!signature.isSignatureIntact()) {
            throw new AssertionError("signature is not cryptographically intact");
        }
        if (!signature.isSignatureValid()) {
            throw new AssertionError("signature is not fully valid (reference data found/intact + signature intact)");
        }

        // (2) signing certificate identified
        if (signature.getSigningCertificate() == null) {
            throw new AssertionError("signing certificate was not identified");
        }

        // (3) level recognized
        SignatureLevel actualLevel = signature.getSignatureFormat();
        if (actualLevel != expectedLevel) {
            throw new AssertionError("signature level = " + actualLevel + ", want " + expectedLevel);
        }

        // (4) for a -T signature, the RFC 3161 token the Go port's own TSP source produced must
        // itself be cryptographically verified by upstream, not merely detected - see
        // dss/cades/testdata/crossgen/CrossGenValidator.java's identical check for why (an earlier
        // revision of that harness asserted only the level, which passed even with a broken
        // token).
        if (expectedLevel == SignatureLevel.PAdES_BASELINE_T) {
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

/*
 * Java oracle driver for CertificateDiagnosticDataBuilder (and, through it, every
 * certificate / certificate-extension / QcStatements builder of DiagnosticDataBuilder
 * and XmlQcStatementsBuilder).
 *
 * It runs the UPSTREAM builder
 *
 *   eu.europa.esig.dss.validation.reports.diagnostic.CertificateDiagnosticDataBuilder
 *
 * over REAL inputs - every DER certificate of dss/spi/testdata/certificate_extensions,
 * one XmlDiagnosticData per certificate, at a fixed validation date - and writes the
 * marshalled XML of each. The Go test builds the same XmlDiagnosticData from the same
 * DER and byte-compares the marshalled document.
 *
 * One document per certificate (rather than one holding all of them) keeps the
 * comparison independent of java.util.HashSet iteration order over the used-certificate
 * set, which is not reproducible in Go.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle CertificateDiagnosticDataOracle.java
 *   java  -cp "$CP:/tmp/oracle" CertificateDiagnosticDataOracle <der-dir> <out-dir>
 */

import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.TokenExtractionStrategy;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.validation.reports.diagnostic.CertificateDiagnosticDataBuilder;

import java.io.File;
import java.nio.file.Files;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;

public class CertificateDiagnosticDataOracle {

    /** Fixed validation time: 2024-01-01T00:00:00Z. */
    private static final Date VALIDATION_DATE = new Date(1704067200000L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File outDir = new File(args[1]);
        outDir.mkdirs();

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".der"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        for (File file : files) {
            CertificateToken token;
            try {
                token = DSSUtils.loadCertificate(Files.readAllBytes(file.toPath()));
            } catch (Exception e) {
                System.err.println("SKIPPED " + file.getName() + ": " + e);
                continue;
            }
            try {
                XmlDiagnosticData diagnosticData = new CertificateDiagnosticDataBuilder()
                        .usedCertificates(Collections.singleton(token))
                        .usedRevocations(Collections.emptySet())
                        .defaultDigestAlgorithm(DigestAlgorithm.SHA256)
                        .tokenExtractionStrategy(TokenExtractionStrategy.NONE)
                        .validationDate(VALIDATION_DATE)
                        .build();
                // marshall(Object, boolean) - the StringWriter overload, the same one
                // DiagnosticDataOracle/DiagnosticDataFillOracle use for the marshal-parity
                // corpus. NOTE: the OutputStream overload takes a different JAXB output
                // class (IndentingUTF8XmlOutput) whose printIndent() wraps the indentation
                // back to column 0 every 8 levels, so the two overloads do NOT produce the
                // same bytes for documents deeper than 8 elements. The Go port reproduces
                // the Writer path, which is what the committed corpus pins.
                String xml = DiagnosticDataFacade.newFacade().marshall(diagnosticData, true);
                Files.write(new File(outDir, file.getName().replace(".der", ".xml")).toPath(),
                        xml.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            } catch (Exception e) {
                System.err.println("SKIPPED-BUILD " + file.getName() + ": " + e);
            }
        }
    }

}

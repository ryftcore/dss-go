import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.SignedDocumentValidator;
import eu.europa.esig.dss.validation.reports.Reports;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Calendar;
import java.util.Date;
import java.util.TimeZone;

/**
 * Ground truth for the Go port of the specs-validation-report generated
 * model: runs upstream validation over a fixture (default validation policy)
 * and marshals the resulting ValidationReportType with ValidationReportFacade
 * (via Reports.getXmlValidationReport()), i.e. exactly the bytes the JAXB
 * reference implementation produces for the model the Go structs mirror.
 *
 * Usage: ValidationReportOracle &lt;outDir&gt; &lt;fixture&gt; ...
 *
 * The validation time is pinned so a dump is reproducible.
 */
public class ValidationReportOracle {

    private static Date fixedTime() {
        Calendar c = Calendar.getInstance(TimeZone.getTimeZone("UTC"));
        c.clear();
        c.set(2025, Calendar.JUNE, 15, 12, 0, 0);
        return c.getTime();
    }

    public static void main(String[] args) throws Exception {
        Path outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);
        int ok = 0;
        int ko = 0;
        for (int i = 1; i < args.length; i++) {
            String docPath = args[i];
            File file = new File(docPath);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try {
                DSSDocument doc = new FileDocument(file);
                SignedDocumentValidator validator = SignedDocumentValidator.fromDocument(doc);
                validator.setCertificateVerifier(new CommonCertificateVerifier());
                validator.setValidationTime(fixedTime());
                Reports reports = validator.validateDocument();
                String xml;
                try {
                    xml = reports.getXmlValidationReport();
                } catch (Exception schemaFailure) {
                    // getXmlValidationReport() requests schema validation by
                    // default in some code paths; fall back to marshalling
                    // without it if the XSD rejects the tree (a validity
                    // question the marshal-parity KAT does not depend on).
                    xml = eu.europa.esig.validationreport.ValidationReportFacade.newFacade()
                            .marshall(reports.getEtsiValidationReportJaxb(), false);
                }
                Files.write(outDir.resolve(name + ".xml"), xml.getBytes(StandardCharsets.UTF_8));
                System.out.println("OK   " + name);
                ok++;
            } catch (Throwable t) {
                System.out.println("SKIP " + name + " : " + t.getClass().getSimpleName() + " " + t.getMessage());
                ko++;
            }
        }
        System.out.println("dumped=" + ok + " skipped=" + ko);
    }
}

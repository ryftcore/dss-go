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
 * Ground truth for the Go port of the dss-simple-report-jaxb generated
 * model: runs upstream validation (default validation policy) over a
 * fixture and dumps Reports.getXmlSimpleReport(), i.e. exactly the bytes
 * SimpleReportFacade produces for the same tree.
 *
 * Usage: SimpleReportOracle &lt;outDir&gt; &lt;fixture&gt;...
 *
 * The validation time is pinned so a dump is reproducible.
 */
public class SimpleReportOracle {

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
                String xml = reports.getXmlSimpleReport();
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

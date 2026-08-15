import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.CertificateValidator;
import eu.europa.esig.dss.validation.reports.CertificateReports;

import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.Calendar;
import java.util.Date;
import java.util.TimeZone;

/**
 * Ground truth for the Go port of the dss-simple-certificate-report-jaxb
 * generated model: runs upstream certificate validation (default validation
 * policy) over a certificate fixture and dumps
 * CertificateReports.getXmlSimpleReport(), i.e. exactly the bytes
 * SimpleCertificateReportFacade produces for the same tree.
 *
 * Usage: SimpleCertificateReportOracle &lt;outDir&gt; &lt;certFile&gt;...
 *
 * The validation time is pinned so a dump is reproducible.
 */
public class SimpleCertificateReportOracle {

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
        CertificateFactory cf = CertificateFactory.getInstance("X.509");
        for (int i = 1; i < args.length; i++) {
            String certPath = args[i];
            File file = new File(certPath);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try (FileInputStream fis = new FileInputStream(file)) {
                X509Certificate x509 = (X509Certificate) cf.generateCertificate(fis);
                CertificateToken token = new CertificateToken(x509);
                CertificateValidator validator = CertificateValidator.fromCertificate(token);
                validator.setCertificateVerifier(new CommonCertificateVerifier());
                validator.setValidationTime(fixedTime());
                CertificateReports reports = validator.validate();
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

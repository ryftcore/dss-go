import eu.europa.esig.dss.simplecertificatereport.SimpleCertificateReportFacade;
import eu.europa.esig.dss.simplecertificatereport.jaxb.XmlSimpleCertificateReport;

import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

public class ReserializeSimpleCertificateReport {
    public static void main(String[] args) throws Exception {
        Path outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);
        for (int i = 1; i < args.length; i++) {
            File file = new File(args[i]);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try (FileInputStream fis = new FileInputStream(file)) {
                XmlSimpleCertificateReport report = SimpleCertificateReportFacade.newFacade().unmarshall(fis);
                String xml = SimpleCertificateReportFacade.newFacade().marshall(report, false);
                Files.write(outDir.resolve(name), xml.getBytes(StandardCharsets.UTF_8));
                System.out.println("OK   " + name);
            } catch (Throwable t) {
                System.out.println("SKIP " + name + " : " + t.getClass().getSimpleName() + " " + t.getMessage());
            }
        }
    }
}

import eu.europa.esig.dss.detailedreport.DetailedReportFacade;
import eu.europa.esig.dss.detailedreport.jaxb.XmlDetailedReport;

import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/**
 * Re-marshals a detailed report through DetailedReportFacade, i.e. through
 * exactly the call AbstractReports.getXmlDetailedReport() makes
 * (DetailedReportFacade.newFacade().marshall(jaxb, validateXml), which
 * marshals into a StringWriter). Using the facade matters: marshalling the
 * same tree into an OutputStream instead selects the RI's
 * IndentingUTF8XmlOutput, which writes the root element's xmlns declaration
 * before the other attributes and wraps its indentation every 8 levels -
 * bytes DSS itself never emits.
 *
 * Usage: ReserializeDetailedReport &lt;outDir&gt; &lt;detailed-report.xml&gt;...
 */
public class ReserializeDetailedReport {
    public static void main(String[] args) throws Exception {
        Path outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);
        for (int i = 1; i < args.length; i++) {
            File file = new File(args[i]);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try (FileInputStream fis = new FileInputStream(file)) {
                XmlDetailedReport report = DetailedReportFacade.newFacade().unmarshall(fis, false);
                String xml = DetailedReportFacade.newFacade().marshall(report, false);
                Files.write(outDir.resolve(name), xml.getBytes(StandardCharsets.UTF_8));
                System.out.println("OK   " + name);
            } catch (Throwable t) {
                System.out.println("SKIP " + name + " : " + t.getClass().getSimpleName() + " " + t.getMessage());
            }
        }
    }
}

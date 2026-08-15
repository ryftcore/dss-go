import eu.europa.esig.dss.simplereport.SimpleReportFacade;
import eu.europa.esig.dss.simplereport.jaxb.XmlSimpleReport;

import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/** Re-marshals an existing SimpleReport XML fixture through the RI facade,
 * so hand/IDE-formatted test resources become valid marshal-parity oracle
 * bytes. Usage: ReserializeSimpleReport &lt;outDir&gt; &lt;fixture&gt;... */
public class ReserializeSimpleReport {
    public static void main(String[] args) throws Exception {
        Path outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);
        for (int i = 1; i < args.length; i++) {
            File file = new File(args[i]);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try (FileInputStream fis = new FileInputStream(file)) {
                XmlSimpleReport report = SimpleReportFacade.newFacade().unmarshall(fis);
                String xml = SimpleReportFacade.newFacade().marshall(report, false);
                Files.write(outDir.resolve(name), xml.getBytes(StandardCharsets.UTF_8));
                System.out.println("OK   " + name);
            } catch (Throwable t) {
                System.out.println("SKIP " + name + " : " + t.getClass().getSimpleName() + " " + t.getMessage());
            }
        }
    }
}

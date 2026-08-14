import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.TokenExtractionStrategy;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.SignedDocumentValidator;

import java.io.File;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Calendar;
import java.util.Collections;
import java.util.Date;
import java.util.TimeZone;

/**
 * Ground truth for the Go port of the dss-diagnostic-jaxb generated model: runs
 * upstream validation over a fixture and marshals the resulting
 * XmlDiagnosticData with DiagnosticDataFacade, i.e. exactly the bytes the JAXB
 * reference implementation produces for the model the Go structs mirror.
 *
 * Usage: DiagnosticDataOracle &lt;outDir&gt; [none:]&lt;fixture&gt;[::detachedFile] ...
 *
 * Tokens are extracted into the dump (Base64Encoded bodies) unless the fixture
 * is prefixed with "none:", which keeps the dump of a large fixture small.
 *
 * The validation time is pinned so a dump is reproducible.
 */
public class DiagnosticDataOracle {

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
            String spec = args[i];
            boolean extract = true;
            if (spec.startsWith("none:")) {
                extract = false;
                spec = spec.substring(5);
            }
            String docPath = spec;
            String detached = null;
            int sep = spec.lastIndexOf("::");
            if (sep > 0) {
                docPath = spec.substring(0, sep);
                detached = spec.substring(sep + 2);
            }
            File file = new File(docPath);
            String name = file.getName().replaceAll("[^A-Za-z0-9._-]", "_");
            try {
                DSSDocument doc = new FileDocument(file);
                SignedDocumentValidator validator = SignedDocumentValidator.fromDocument(doc);
                validator.setCertificateVerifier(new CommonCertificateVerifier());
                validator.setValidationTime(fixedTime());
                // Extract everything: the dumps then also carry Base64Encoded
                // certificate/revocation/timestamp bodies.
                validator.setTokenExtractionStrategy(
                        extract ? TokenExtractionStrategy.EXTRACT_ALL : TokenExtractionStrategy.NONE);
                if (detached != null) {
                    validator.setDetachedContents(Collections.singletonList(new FileDocument(new File(detached))));
                }
                XmlDiagnosticData dd = validator.getDiagnosticData();
                String xml;
                try {
                    xml = DiagnosticDataFacade.newFacade().marshall(dd, true);
                } catch (Exception schemaFailure) {
                    xml = DiagnosticDataFacade.newFacade().marshall(dd, false);
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

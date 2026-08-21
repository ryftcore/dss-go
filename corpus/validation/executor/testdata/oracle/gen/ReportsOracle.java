import eu.europa.esig.dss.detailedreport.DetailedReportFacade;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.ValidationLevel;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.simplecertificatereport.SimpleCertificateReportFacade;
import eu.europa.esig.dss.simplereport.SimpleReportFacade;
import eu.europa.esig.dss.validation.executor.certificate.DefaultCertificateProcessExecutor;
import eu.europa.esig.dss.validation.executor.signature.DefaultSignatureProcessExecutor;
import eu.europa.esig.dss.validation.policy.ValidationPolicyLoader;
import eu.europa.esig.dss.validation.reports.CertificateReports;
import eu.europa.esig.dss.validation.reports.Reports;
import eu.europa.esig.validationreport.ValidationReportFacade;

import java.io.File;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;
import java.util.Locale;

/** Dumps the SHA-256 of every report the phase-8f executors produce, plus full XML for a few files. */
public class ReportsOracle {

    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        File inDir = new File(args[0]);
        Path outDir = Paths.get(args[1]);
        Path xmlDir = Paths.get(args[2]);
        Files.createDirectories(outDir);
        Files.createDirectories(xmlDir);
        List<String> keepXml = Arrays.asList(args).subList(3, args.length);
        Date currentTime = new Date(1700000000000L);

        File[] files = inDir.listFiles((d, n) -> n.endsWith(".xml") && !n.endsWith(".javaout"));
        Arrays.sort(files);
        List<String> rows = new ArrayList<>();
        for (File f : files) {
            String name = f.getName();
            if (name.startsWith("model-")) {
                continue; // schema-coverage fixtures, excluded as in the blocks corpus
            }
            String base = name.substring(0, name.length() - 4);
            XmlDiagnosticData dd = DiagnosticDataFacade.newFacade().unmarshall(f, false);

            ValidationPolicy sigPolicy = ValidationPolicyLoader.fromDefaultValidationPolicy().create();
            DefaultSignatureProcessExecutor exec = new DefaultSignatureProcessExecutor();
            exec.setDiagnosticData(dd);
            exec.setValidationPolicy(sigPolicy);
            exec.setCurrentTime(currentTime);
            exec.setValidationLevel(ValidationLevel.ARCHIVAL_DATA);
            exec.setEnableEtsiValidationReport(true);
            exec.setLocale(Locale.ENGLISH);
            Reports reports = exec.execute();
            String sr = SimpleReportFacade.newFacade().marshall(reports.getSimpleReportJaxb(), false);
            String dr = DetailedReportFacade.newFacade().marshall(reports.getDetailedReportJaxb(), false);
            String vr = ValidationReportFacade.newFacade().marshall(reports.getEtsiValidationReportJaxb(), false);

            String certId = "";
            String scr = "";
            String cdr = "";
            dd = DiagnosticDataFacade.newFacade().unmarshall(f, false);
            if (!dd.getUsedCertificates().isEmpty()) {
                XmlCertificate cert = dd.getUsedCertificates().get(0);
                certId = cert.getId();
                try (InputStream is = ReportsOracle.class.getResourceAsStream("/policy/certificate-constraint.xml")) {
                    ValidationPolicy certPolicy = ValidationPolicyLoader.fromValidationPolicy(is).create();
                    DefaultCertificateProcessExecutor cexec = new DefaultCertificateProcessExecutor();
                    cexec.setDiagnosticData(dd);
                    cexec.setValidationPolicy(certPolicy);
                    cexec.setCurrentTime(currentTime);
                    cexec.setCertificateId(certId);
                    cexec.setLocale(Locale.ENGLISH);
                    CertificateReports creports = cexec.execute();
                    scr = SimpleCertificateReportFacade.newFacade().marshall(creports.getSimpleReportJaxb(), false);
                    cdr = DetailedReportFacade.newFacade().marshall(creports.getDetailedReportJaxb(), false);
                }
            }

            rows.add(String.format(
                    "{\"file\":\"%s\",\"simpleReport\":\"%s\",\"detailedReport\":\"%s\",\"etsiValidationReport\":\"%s\","
                            + "\"certificateId\":\"%s\",\"simpleCertificateReport\":\"%s\",\"certificateDetailedReport\":\"%s\"}",
                    base, sha256(sr), sha256(dr), sha256(vr), certId, sha256(scr), sha256(cdr)));

            if (keepXml.contains(base)) {
                Files.write(xmlDir.resolve(base + ".sr.xml"), sr.getBytes(StandardCharsets.UTF_8));
                Files.write(xmlDir.resolve(base + ".dr.xml"), dr.getBytes(StandardCharsets.UTF_8));
                Files.write(xmlDir.resolve(base + ".vr.xml"), vr.getBytes(StandardCharsets.UTF_8));
                if (!certId.isEmpty()) {
                    Files.write(xmlDir.resolve(base + ".scr.xml"), scr.getBytes(StandardCharsets.UTF_8));
                    Files.write(xmlDir.resolve(base + ".cdr.xml"), cdr.getBytes(StandardCharsets.UTF_8));
                }
            }
        }
        Files.write(outDir.resolve("reports.jsonl"), String.join("\n", rows).concat("\n").getBytes(StandardCharsets.UTF_8));
        System.out.println("rows=" + rows.size());
    }

    private static String sha256(String s) throws Exception {
        if (s == null || s.isEmpty()) {
            return "";
        }
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        byte[] d = md.digest(s.getBytes(StandardCharsets.UTF_8));
        StringBuilder sb = new StringBuilder();
        for (byte b : d) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }
}

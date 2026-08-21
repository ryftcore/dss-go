import eu.europa.esig.dss.detailedreport.DetailedReportFacade;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
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
import java.nio.file.*;
import java.security.MessageDigest;
import java.util.*;

/** Byte-parity oracle over the FULL diag-data corpus: sha256 of every marshalled report. */
public class FullCorpusReportsOracle {
    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        File inDir = new File(args[0]);
        Path out = Paths.get(args[1]);
        Date now = new Date(1700000000000L);
        InputStream cps = FullCorpusReportsOracle.class.getResourceAsStream("/policy/certificate-constraint.xml");
        byte[] certPolicy = cps.readAllBytes(); cps.close();

        List<File> files = new ArrayList<>();
        collect(inDir, files);
        files.sort(Comparator.comparing(FullCorpusReportsOracle::rel));
        StringBuilder sb = new StringBuilder();
        int n = 0, ok = 0;
        for (File f : files) {
            String key = rel(f);
            XmlDiagnosticData dd;
            try { dd = DiagnosticDataFacade.newFacade().unmarshall(f, false); }
            catch (Throwable t) { n++; continue; }
            sb.append("{\"file\":\"").append(key).append("\"");
            try {
                DefaultSignatureProcessExecutor e = new DefaultSignatureProcessExecutor();
                e.setDiagnosticData(dd);
                e.setValidationPolicy(ValidationPolicyLoader.fromDefaultValidationPolicy().create());
                e.setCurrentTime(now);
                e.setValidationLevel(ValidationLevel.ARCHIVAL_DATA);
                e.setEnableEtsiValidationReport(true);
                e.setLocale(Locale.ENGLISH);
                Reports r = e.execute();
                sb.append(",\"simpleReport\":\"").append(sha(SimpleReportFacade.newFacade().marshall(r.getSimpleReportJaxb(), false))).append("\"");
                sb.append(",\"detailedReport\":\"").append(sha(DetailedReportFacade.newFacade().marshall(r.getDetailedReportJaxb(), false))).append("\"");
                sb.append(",\"etsiValidationReport\":\"").append(sha(ValidationReportFacade.newFacade().marshall(r.getEtsiValidationReportJaxb(), false))).append("\"");
            } catch (Throwable t) {
                sb.append(",\"signatureError\":\"").append(t.getClass().getSimpleName()).append("\"");
            }
            try {
                XmlDiagnosticData dd2 = DiagnosticDataFacade.newFacade().unmarshall(f, false);
                if (!dd2.getUsedCertificates().isEmpty()) {
                    String certId = dd2.getUsedCertificates().get(0).getId();
                    ValidationPolicy p = ValidationPolicyLoader.fromValidationPolicy(new java.io.ByteArrayInputStream(certPolicy)).create();
                    DefaultCertificateProcessExecutor c = new DefaultCertificateProcessExecutor();
                    c.setDiagnosticData(dd2); c.setValidationPolicy(p); c.setCurrentTime(now);
                    c.setCertificateId(certId); c.setLocale(Locale.ENGLISH);
                    CertificateReports cr = c.execute();
                    sb.append(",\"certificateId\":\"").append(certId).append("\"");
                    sb.append(",\"simpleCertificateReport\":\"").append(sha(SimpleCertificateReportFacade.newFacade().marshall(cr.getSimpleReportJaxb(), false))).append("\"");
                    sb.append(",\"certificateDetailedReport\":\"").append(sha(DetailedReportFacade.newFacade().marshall(cr.getDetailedReportJaxb(), false))).append("\"");
                }
            } catch (Throwable t) {
                sb.append(",\"certificateError\":\"").append(t.getClass().getSimpleName()).append("\"");
            }
            sb.append("}\n"); n++; ok++;
        }
        Files.write(out, sb.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("files="+n+" rows="+ok);
    }
    static void collect(File d, List<File> o){ File[] c=d.listFiles(); if(c==null) return; for(File x:c){ if(x.isDirectory()) collect(x,o); else if(x.getName().endsWith(".xml")) o.add(x);} }
    static String rel(File f){ String a=f.getAbsolutePath().replace('\\','/'); int i=a.indexOf("/diag-data/"); return i>=0?a.substring(i+11):f.getName(); }
    static String sha(String s) throws Exception { if(s==null||s.isEmpty()) return ""; MessageDigest m=MessageDigest.getInstance("SHA-256"); StringBuilder b=new StringBuilder(); for(byte y:m.digest(s.getBytes(StandardCharsets.UTF_8))) b.append(String.format("%02x",y)); return b.toString(); }
}

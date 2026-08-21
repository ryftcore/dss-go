import eu.europa.esig.dss.detailedreport.jaxb.XmlBasicBuildingBlocks;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCertificate;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlDetailedReport;
import eu.europa.esig.dss.detailedreport.jaxb.XmlEvidenceRecord;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSignature;
import eu.europa.esig.dss.detailedreport.jaxb.XmlTimestamp;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.ValidationLevel;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.simplereport.jaxb.XmlToken;
import eu.europa.esig.dss.validation.executor.certificate.DefaultCertificateProcessExecutor;
import eu.europa.esig.dss.validation.executor.signature.DefaultSignatureProcessExecutor;
import eu.europa.esig.dss.validation.policy.ValidationPolicyLoader;
import eu.europa.esig.dss.validation.reports.CertificateReports;
import eu.europa.esig.dss.validation.reports.Reports;

import java.io.File;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.Date;
import java.util.List;
import java.util.Locale;

/**
 * Dumps, for EVERY diagnostic-data XML in the upstream corpus (recursively,
 * every subdirectory, no exclusions), the final Indication/SubIndication of
 * every signature/timestamp/evidence-record the default signature executor
 * produces, every BasicBuildingBlocks conclusion, and (when the document
 * carries at least one used certificate) the same for the certificate
 * executor run against its first certificate. This is the Phase 8f exit
 * criterion's executor-level oracle: one row per input file, no quarantines.
 *
 * Both executors run with the SAME configuration used everywhere else in this
 * port's oracle corpus (default validation policy / certificate-constraint.xml,
 * ValidationLevel.ARCHIVAL_DATA, validation time 1700000000000, en locale) -
 * this is a uniform parity probe, not a reproduction of each fixture's
 * original upstream unit test scenario (several subdirectories - cert-validation,
 * timestamp-validation, qwac-validation, eaa-validation - are consumed by
 * OTHER executors upstream, e.g. DefaultCertificateProcessExecutor directly or
 * DefaultTimestampProcessExecutor). Running one fixed executor pair over every
 * file still exercises the exact same Go/Java code paths on identical input,
 * which is what "verdict parity" means here; it is documented, not silent.
 */
public class FullCorpusOracle {

    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        File inDir = new File(args[0]);
        Path outFile = Paths.get(args[1]);
        Date currentTime = new Date(1700000000000L);

        InputStream certPolicyStream = FullCorpusOracle.class.getResourceAsStream("/policy/certificate-constraint.xml");
        byte[] certPolicyBytes = certPolicyStream.readAllBytes();
        certPolicyStream.close();

        List<File> files = new ArrayList<>();
        collectXml(inDir, files);
        files.sort(Comparator.comparing(FullCorpusOracle::relKey));

        StringBuilder out = new StringBuilder();
        int rows = 0;
        for (File f : files) {
            String key = relKey(f);
            out.append("{\"file\":").append(json(key));
            try {
                XmlDiagnosticData dd = DiagnosticDataFacade.newFacade().unmarshall(f, false);

                ValidationPolicy sigPolicy = ValidationPolicyLoader.fromDefaultValidationPolicy().create();
                DefaultSignatureProcessExecutor exec = new DefaultSignatureProcessExecutor();
                exec.setDiagnosticData(dd);
                exec.setValidationPolicy(sigPolicy);
                exec.setCurrentTime(currentTime);
                exec.setValidationLevel(ValidationLevel.ARCHIVAL_DATA);
                exec.setEnableEtsiValidationReport(false);
                exec.setLocale(Locale.ENGLISH);
                Reports reports = exec.execute();
                XmlDetailedReport dr = reports.getDetailedReportJaxb();

                out.append(",\"signatureExecutor\":{");
                dumpTokens(out, dr);
                out.append(",\"bbb\":");
                dumpBbb(out, dr);
                out.append(",\"qualification\":");
                dumpQualification(out, reports.getSimpleReportJaxb().getSignatureOrTimestampOrEvidenceRecord());
                out.append("}");
            } catch (Throwable t) {
                out.append(",\"signatureExecutorError\":").append(json(errClass(t) + ": " + t.getMessage()));
            }

            try {
                XmlDiagnosticData dd2 = DiagnosticDataFacade.newFacade().unmarshall(f, false);
                if (!dd2.getUsedCertificates().isEmpty()) {
                    String certId = dd2.getUsedCertificates().get(0).getId();
                    ValidationPolicy certPolicy = ValidationPolicyLoader
                            .fromValidationPolicy(new java.io.ByteArrayInputStream(certPolicyBytes)).create();
                    DefaultCertificateProcessExecutor cexec = new DefaultCertificateProcessExecutor();
                    cexec.setDiagnosticData(dd2);
                    cexec.setValidationPolicy(certPolicy);
                    cexec.setCurrentTime(currentTime);
                    cexec.setCertificateId(certId);
                    cexec.setLocale(Locale.ENGLISH);
                    CertificateReports creports = cexec.execute();
                    XmlDetailedReport cdr = creports.getDetailedReportJaxb();

                    out.append(",\"certificateExecutor\":{\"certificateId\":").append(json(certId));
                    out.append(",\"bbb\":");
                    dumpBbb(out, cdr);
                    out.append("}");
                }
            } catch (Throwable t) {
                out.append(",\"certificateExecutorError\":").append(json(errClass(t) + ": " + t.getMessage()));
            }

            out.append("}\n");
            rows++;
        }
        Files.write(outFile, out.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("rows=" + rows);
    }

    private static void collectXml(File dir, List<File> out) {
        File[] children = dir.listFiles();
        if (children == null) {
            return;
        }
        for (File c : children) {
            if (c.isDirectory()) {
                collectXml(c, out);
            } else if (c.getName().endsWith(".xml")) {
                out.add(c);
            }
        }
    }

    // relKey is the path relative to the diag-data root, with forward
    // slashes, so the Go side can address the same row regardless of OS.
    private static String relKey(File f) {
        String abs = f.getAbsolutePath().replace('\\', '/');
        int idx = abs.indexOf("/diag-data/");
        return idx >= 0 ? abs.substring(idx + "/diag-data/".length()) : f.getName();
    }

    private static void dumpTokens(StringBuilder out, XmlDetailedReport dr) {
        out.append("\"tokens\":[");
        boolean first = true;
        for (Object item : dr.getSignatureOrTimestampOrEvidenceRecord()) {
            if (!first) {
                out.append(",");
            }
            first = false;
            String kind;
            String id;
            XmlConclusion c;
            if (item instanceof XmlSignature s) {
                kind = "Signature";
                id = s.getId();
                c = s.getConclusion();
            } else if (item instanceof XmlTimestamp t) {
                kind = "Timestamp";
                id = t.getId();
                c = t.getConclusion();
            } else if (item instanceof XmlEvidenceRecord e) {
                kind = "EvidenceRecord";
                id = e.getId();
                c = e.getConclusion();
            } else {
                kind = item.getClass().getSimpleName();
                id = "";
                c = null;
            }
            out.append("{\"kind\":").append(json(kind));
            out.append(",\"id\":").append(json(id));
            out.append(",\"indication\":").append(json(c == null || c.getIndication() == null ? "" : c.getIndication().toString()));
            out.append(",\"subIndication\":").append(json(c == null || c.getSubIndication() == null ? "" : c.getSubIndication().toString()));
            out.append("}");
        }
        out.append("]");
    }

    private static void dumpBbb(StringBuilder out, XmlDetailedReport dr) {
        out.append("[");
        boolean first = true;
        for (XmlBasicBuildingBlocks bbb : dr.getBasicBuildingBlocks()) {
            if (!first) {
                out.append(",");
            }
            first = false;
            XmlConclusion c = bbb.getConclusion();
            out.append("{\"id\":").append(json(bbb.getId()));
            out.append(",\"type\":").append(json(bbb.getType() == null ? "" : bbb.getType().toString()));
            out.append(",\"indication\":").append(json(c == null || c.getIndication() == null ? "" : c.getIndication().toString()));
            out.append(",\"subIndication\":").append(json(c == null || c.getSubIndication() == null ? "" : c.getSubIndication().toString()));
            out.append("}");
        }
        out.append("]");
    }

    private static void dumpQualification(StringBuilder out, List<XmlToken> tokens) {
        out.append("[");
        boolean first = true;
        for (XmlToken t : tokens) {
            if (!(t instanceof eu.europa.esig.dss.simplereport.jaxb.XmlSignature sig)) {
                continue;
            }
            if (!first) {
                out.append(",");
            }
            first = false;
            String q = sig.getSignatureLevel() == null || sig.getSignatureLevel().getValue() == null
                    ? "" : sig.getSignatureLevel().getValue().name();
            out.append("{\"id\":").append(json(sig.getId()));
            out.append(",\"qualification\":").append(json(q));
            out.append("}");
        }
        out.append("]");
    }

    private static String errClass(Throwable t) {
        Throwable root = t;
        while (root.getCause() != null && root.getCause() != root) {
            root = root.getCause();
        }
        return root.getClass().getSimpleName();
    }

    private static String json(String s) {
        if (s == null) {
            s = "";
        }
        StringBuilder sb = new StringBuilder("\"");
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"': sb.append("\\\""); break;
                case '\\': sb.append("\\\\"); break;
                case '\n': sb.append("\\n"); break;
                case '\r': sb.append("\\r"); break;
                case '\t': sb.append("\\t"); break;
                default:
                    if (c < 0x20) {
                        sb.append(String.format("\\u%04x", (int) c));
                    } else {
                        sb.append(c);
                    }
            }
        }
        sb.append("\"");
        return sb.toString();
    }
}

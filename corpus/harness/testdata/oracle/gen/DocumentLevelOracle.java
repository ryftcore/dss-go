import eu.europa.esig.dss.alert.SilentOnStatusAlert;
import eu.europa.esig.dss.detailedreport.jaxb.XmlBasicBuildingBlocks;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlDetailedReport;
import eu.europa.esig.dss.detailedreport.jaxb.XmlEvidenceRecord;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSignature;
import eu.europa.esig.dss.detailedreport.jaxb.XmlTimestamp;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.ValidationLevel;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.simplereport.jaxb.XmlToken;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.validation.SignedDocumentValidator;
import eu.europa.esig.dss.validation.reports.Reports;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

/**
 * Item (B) of the s8f harness brief: document-level end-to-end parity. For
 * every (format, name, path) row of a manifest TSV, loads the real signed
 * document, dispatches it through SignedDocumentValidator.fromDocument
 * (upstream's format-autodetection registry - the same one
 * dssvalidation.SignedDocumentValidatorFromDocument on the Go side mirrors),
 * validates it with a permissive CertificateVerifier (matching every format
 * package's own smoke-test convention: CommonCertificateVerifier(true) with
 * the same seven alerts silenced - see e.g.
 * cades/cms_document_validator_smoke_test.go's permissiveCertificateVerifier)
 * at ValidationLevel.ARCHIVAL_DATA, and dumps the resulting verdicts plus a
 * diagnostic-data core summary.
 */
public class DocumentLevelOracle {

    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        Path manifestPath = Paths.get(args[0]);
        Path outFile = Paths.get(args[1]);

        StringBuilder out = new StringBuilder();
        int rows = 0;
        try (BufferedReader r = new BufferedReader(new FileReader(manifestPath.toFile(), StandardCharsets.UTF_8))) {
            String line;
            while ((line = r.readLine()) != null) {
                if (line.isBlank()) {
                    continue;
                }
                String[] parts = line.split("\t");
                String format = parts[0];
                String name = parts[1];
                String path = parts[2];

                out.append("{\"format\":").append(json(format));
                out.append(",\"name\":").append(json(name));
                try {
                    DSSDocument doc = new FileDocument(new File(path));
                    SignedDocumentValidator validator = SignedDocumentValidator.fromDocument(doc);
                    validator.setCertificateVerifier(permissiveCertificateVerifier());
                    validator.setValidationLevel(ValidationLevel.ARCHIVAL_DATA);
                    validator.setLocale(Locale.ENGLISH);
                    Reports reports = validator.validateDocument();
                    XmlDetailedReport dr = reports.getDetailedReportJaxb();
                    XmlDiagnosticData dd = reports.getDiagnosticDataJaxb();

                    out.append(",\"core\":{");
                    out.append("\"signaturesCount\":").append(dd.getSignatures().size());
                    out.append(",\"certificatesCount\":").append(dd.getUsedCertificates().size());
                    out.append(",\"containerType\":").append(json(dd.getContainerInfo() == null || dd.getContainerInfo().getContainerType() == null
                            ? "" : dd.getContainerInfo().getContainerType().name()));
                    out.append("}");

                    out.append(",\"tokens\":[");
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

                    out.append(",\"bbb\":[");
                    first = true;
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

                    out.append(",\"qualification\":[");
                    first = true;
                    List<XmlToken> simpleTokens = reports.getSimpleReportJaxb().getSignatureOrTimestampOrEvidenceRecord();
                    for (XmlToken t : simpleTokens) {
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
                } catch (Throwable t) {
                    out.append(",\"error\":").append(json(errClass(t) + ": " + t.getMessage()));
                }
                out.append("}\n");
                rows++;
            }
        }
        Files.write(outFile, out.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("rows=" + rows);
    }

    private static CommonCertificateVerifier permissiveCertificateVerifier() {
        CommonCertificateVerifier v = new CommonCertificateVerifier(true);
        v.setAlertOnMissingRevocationData(new SilentOnStatusAlert());
        v.setAlertOnRevokedCertificate(new SilentOnStatusAlert());
        v.setAlertOnInvalidSignature(new SilentOnStatusAlert());
        v.setAlertOnInvalidTimestamp(new SilentOnStatusAlert());
        v.setAlertOnUncoveredPOE(new SilentOnStatusAlert());
        v.setAlertOnExpiredCertificate(new SilentOnStatusAlert());
        v.setAlertOnNotYetValidCertificate(new SilentOnStatusAlert());
        return v;
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

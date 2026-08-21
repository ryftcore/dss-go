/*
 * Java oracle driver for the block-level conclusions of the phase-8e
 * eu.europa.esig.dss.validation.process.qualification tree.
 *
 * It drives eu.europa.esig.dss.validation.process.qualification.certificate.CertQualificationAtTimeBlock
 * - the block that carries the normative TrustServiceFilter ORDER of
 * TS 119 615, every certificate-qualification leaf check, the per-service
 * consistency warnings and the conflict-detection abort - over a designed
 * cross-product of (certificate, list of trust services, validation time),
 * and dumps the WHOLE XmlValidationCertificateQualification: the ordered
 * Constraint list with each item's name key, status, error/warning/info keys
 * and additional info, the Conclusion (indication, sub-indication, ordered
 * message keys), the resulting CertificateQualification, and the block's
 * getFilteredServices() outcome.
 *
 * The Go replay (../../qual_block_oracle_test.go) rebuilds the identical
 * inputs and asserts the identical tree, so a divergence in filter order, in
 * a check's message tag, or in the number of emitted constraints fails.
 *
 * Rows are written to
 *   dss/validation/process/qualification/testdata/oracle/qual_block.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle QualBlockOracle.java
 *   java  -cp "$CP:/tmp/oracle" QualBlockOracle <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlValidationCertificateQualification;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.TrustServiceWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateExtension;
import eu.europa.esig.dss.diagnostic.jaxb.XmlOID;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcCompliance;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcSSCD;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcStatements;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQualifier;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustedList;
import eu.europa.esig.dss.enumerations.CertificateExtensionEnum;
import eu.europa.esig.dss.enumerations.QCTypeEnum;
import eu.europa.esig.dss.enumerations.ServiceQualification;
import eu.europa.esig.dss.enumerations.ValidationTime;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.validation.process.qualification.certificate.CertQualificationAtTimeBlock;
import eu.europa.esig.dss.validation.process.qualification.trust.ServiceTypeIdentifier;
import eu.europa.esig.dss.validation.process.qualification.trust.TrustServiceStatus;

import java.io.IOException;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Date;
import java.util.List;

public class QualBlockOracle {

    private static final long PRE_EIDAS = 1370044800000L;   // 2013-06-01
    private static final long POST_EIDAS = 1514764800000L;  // 2018-01-01
    private static final long LATER = 1577836800000L;       // 2020-01-01
    private static final long EARLY = 1420070400000L;       // 2015-01-01

    // ---- certificates -----------------------------------------------------

    /** label -> certificate. */
    private static final String[] CERT_LABELS = {
            "plain-post", "qc-esign-post", "qc-esign-qscd-post", "qc-eseal-post",
            "qc-wsa-post", "qc-multi-type-post", "qc-esign-pre", "plain-pre"};

    private static XmlCertificate certificate(String label) {
        boolean post = !label.endsWith("-pre");
        long notBefore = post ? POST_EIDAS : PRE_EIDAS;
        Boolean compliance = label.startsWith("qc-") ? Boolean.TRUE : Boolean.FALSE;
        Boolean sscd = label.contains("qscd") ? Boolean.TRUE : Boolean.FALSE;
        List<String> types;
        if (label.contains("multi-type")) {
            types = Arrays.asList(QCTypeEnum.QCT_ESIGN.getOid(), QCTypeEnum.QCT_ESEAL.getOid());
        } else if (label.contains("eseal")) {
            types = Arrays.asList(QCTypeEnum.QCT_ESEAL.getOid());
        } else if (label.contains("wsa")) {
            types = Arrays.asList(QCTypeEnum.QCT_WEB.getOid());
        } else if (label.contains("esign")) {
            types = Arrays.asList(QCTypeEnum.QCT_ESIGN.getOid());
        } else {
            types = Collections.emptyList();
        }

        XmlCertificate cert = new XmlCertificate();
        cert.setId("C-" + label);
        cert.setNotBefore(new Date(notBefore));
        cert.setNotAfter(new Date(notBefore + 3L * 365 * 86400000L));
        XmlQcStatements qcStatements = new XmlQcStatements();
        qcStatements.setOID(CertificateExtensionEnum.QC_STATEMENTS.getOid());
        XmlQcCompliance xmlCompliance = new XmlQcCompliance();
        xmlCompliance.setPresent(compliance);
        qcStatements.setQcCompliance(xmlCompliance);
        XmlQcSSCD xmlSscd = new XmlQcSSCD();
        xmlSscd.setPresent(sscd);
        qcStatements.setQcSSCD(xmlSscd);
        for (String oid : types) {
            XmlOID xmlOid = new XmlOID();
            xmlOid.setValue(oid);
            qcStatements.getQcTypes().add(xmlOid);
        }
        List<XmlCertificateExtension> extensions = new ArrayList<>();
        extensions.add(qcStatements);
        cert.getCertificateExtensions().addAll(extensions);
        return cert;
    }

    // ---- trust services ---------------------------------------------------

    private static final class Svc {
        final String name;
        final String status;
        final String type;
        final long startDate;
        final Long endDate;
        final List<String> qualifiers;
        final List<String> asis;
        final Boolean mra;

        Svc(String name, String status, String type, long startDate, Long endDate,
            List<String> qualifiers, List<String> asis, Boolean mra) {
            this.name = name;
            this.status = status;
            this.type = type;
            this.startDate = startDate;
            this.endDate = endDate;
            this.qualifiers = qualifiers;
            this.asis = asis;
            this.mra = mra;
        }
    }

    private static final String GRANTED = TrustServiceStatus.GRANTED.getUri();
    private static final String WITHDRAWN = TrustServiceStatus.WITHDRAWN.getUri();
    private static final String UNDER_SUPERVISION = TrustServiceStatus.UNDER_SUPERVISION.getUri();
    private static final String CA_QC = ServiceTypeIdentifier.CA_QC.getUri();
    private static final String CA_PKC = ServiceTypeIdentifier.CA_PKC.getUri();
    private static final String ASI_ESIG = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures";
    private static final String ASI_ESEAL = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals";

    private static final Svc[] SERVICES = {
            new Svc("granted-caqc", GRANTED, CA_QC, POST_EIDAS, null,
                    Collections.<String>emptyList(), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-qcstatement", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_STATEMENT.getUri()), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-notqualified", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.NOT_QUALIFIED.getUri()), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-withqscd", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_WITH_QSCD.getUri()), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-eseal", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_FOR_ESEAL.getUri()), Arrays.asList(ASI_ESEAL), null),
            new Svc("granted-caqc-inconsistent-usage", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_FOR_ESIG.getUri(), ServiceQualification.QC_FOR_ESEAL.getUri()),
                    Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-inconsistent-qscd", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_WITH_QSCD.getUri(), ServiceQualification.QC_NO_QSCD.getUri()),
                    Arrays.asList(ASI_ESIG), null),
            new Svc("withdrawn-caqc", WITHDRAWN, CA_QC, POST_EIDAS, null,
                    Collections.<String>emptyList(), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-capkc", GRANTED, CA_PKC, POST_EIDAS, null,
                    Collections.<String>emptyList(), Arrays.asList(ASI_ESIG), null),
            new Svc("undersupervision-caqc-pre", UNDER_SUPERVISION, CA_QC, PRE_EIDAS, null,
                    Collections.<String>emptyList(), Collections.<String>emptyList(), null),
            new Svc("granted-caqc-expired", GRANTED, CA_QC, PRE_EIDAS, EARLY,
                    Collections.<String>emptyList(), Arrays.asList(ASI_ESIG), null),
            new Svc("granted-caqc-mra", GRANTED, CA_QC, POST_EIDAS, null,
                    Collections.<String>emptyList(), Arrays.asList(ASI_ESIG), Boolean.TRUE),
    };

    private static TrustServiceWrapper service(Svc s) {
        TrustServiceWrapper w = new TrustServiceWrapper();
        w.setServiceNames(Collections.singletonList(s.name));
        w.setStatus(s.status);
        w.setType(s.type);
        w.setStartDate(new Date(s.startDate));
        w.setEndDate(s.endDate == null ? null : new Date(s.endDate));
        List<XmlQualifier> qualifiers = new ArrayList<>();
        for (String uri : s.qualifiers) {
            XmlQualifier q = new XmlQualifier();
            q.setValue(uri);
            qualifiers.add(q);
        }
        w.setCapturedQualifiers(qualifiers);
        w.setAdditionalServiceInfos(new ArrayList<>(s.asis));
        w.setCountryCode("LU");
        XmlTrustedList tl = new XmlTrustedList();
        tl.setUrl("https://tl.example/LU");
        tl.setMra(s.mra);
        w.setTrustedList(tl);
        w.setEnactedMRA(s.mra);
        if (Boolean.TRUE.equals(s.mra)) {
            w.setMraTrustServiceEquivalenceStatusStartingTime(new Date(POST_EIDAS));
        }
        return w;
    }

    /** The service-list shapes the block is driven over: singletons, pairs, and the whole panel. */
    private static List<int[]> serviceSelections() {
        List<int[]> result = new ArrayList<>();
        result.add(new int[0]);
        for (int i = 0; i < SERVICES.length; i++) {
            result.add(new int[]{i});
        }
        for (int i = 0; i < SERVICES.length; i++) {
            for (int j = i + 1; j < SERVICES.length; j++) {
                result.add(new int[]{i, j});
            }
        }
        int[] all = new int[SERVICES.length];
        for (int i = 0; i < SERVICES.length; i++) {
            all[i] = i;
        }
        result.add(all);
        return result;
    }

    private static final Long[] DATES = {null, PRE_EIDAS, POST_EIDAS, LATER, EARLY};

    // ---- JSON -------------------------------------------------------------

    private static String esc(String s) {
        if (s == null) {
            return "null";
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
        return sb.append('"').toString();
    }

    private static String arr(List<String> values) {
        if (values == null) {
            return "null";
        }
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < values.size(); i++) {
            if (i > 0) {
                sb.append(',');
            }
            sb.append(esc(values.get(i)));
        }
        return sb.append(']').toString();
    }

    private static String msgKey(XmlMessage m) {
        return m == null ? "null" : esc(m.getKey());
    }

    private static String msgKeys(List<XmlMessage> messages) {
        List<String> keys = new ArrayList<>();
        if (messages != null) {
            for (XmlMessage m : messages) {
                keys.add(m.getKey());
            }
        }
        return arr(keys);
    }

    private static String constraints(List<XmlConstraint> list) {
        StringBuilder sb = new StringBuilder("[");
        if (list != null) {
            for (int i = 0; i < list.size(); i++) {
                XmlConstraint c = list.get(i);
                if (i > 0) {
                    sb.append(',');
                }
                sb.append("{\"name\":").append(msgKey(c.getName()));
                sb.append(",\"status\":").append(esc(c.getStatus() == null ? null : c.getStatus().name()));
                sb.append(",\"error\":").append(msgKey(c.getError()));
                sb.append(",\"warning\":").append(msgKey(c.getWarning()));
                sb.append(",\"info\":").append(msgKey(c.getInfo()));
                sb.append(",\"additionalInfo\":").append(esc(c.getAdditionalInfo()));
                sb.append(",\"id\":").append(esc(c.getId()));
                sb.append('}');
            }
        }
        return sb.append(']').toString();
    }

    private static String conclusion(XmlConclusion c) {
        if (c == null) {
            return "null";
        }
        return "{\"indication\":" + esc(c.getIndication() == null ? null : c.getIndication().name())
                + ",\"subIndication\":" + esc(c.getSubIndication() == null ? null : c.getSubIndication().name())
                + ",\"errors\":" + msgKeys(c.getErrors())
                + ",\"warnings\":" + msgKeys(c.getWarnings())
                + ",\"infos\":" + msgKeys(c.getInfos()) + "}";
    }

    public static void main(String[] args) throws IOException {
        if (args.length < 1) {
            System.err.println("usage: QualBlockOracle <dss-repo-root>");
            System.exit(2);
        }
        Path out = Paths.get(args[0], "validation", "process", "qualification", "testdata", "oracle", "qual_block.jsonl");
        Files.createDirectories(out.getParent());
        I18nProvider i18nProvider = new I18nProvider();
        int index = 0;
        try (PrintWriter pw = new PrintWriter(Files.newBufferedWriter(out, StandardCharsets.UTF_8))) {
            for (String certLabel : CERT_LABELS) {
                for (int[] selection : serviceSelections()) {
                    for (Long date : DATES) {
                        for (ValidationTime validationTime :
                                new ValidationTime[]{ValidationTime.CERTIFICATE_ISSUANCE_TIME, ValidationTime.BEST_SIGNATURE_TIME}) {
                            CertificateWrapper cert = new CertificateWrapper(certificate(certLabel));
                            List<TrustServiceWrapper> services = new ArrayList<>();
                            List<String> labels = new ArrayList<>();
                            for (int i : selection) {
                                services.add(service(SERVICES[i]));
                                labels.add(SERVICES[i].name);
                            }
                            CertQualificationAtTimeBlock block =
                                    validationTime == ValidationTime.CERTIFICATE_ISSUANCE_TIME
                                            ? new CertQualificationAtTimeBlock(i18nProvider, validationTime, cert, services)
                                            : new CertQualificationAtTimeBlock(i18nProvider, validationTime,
                                                    date == null ? null : new Date(date), cert, services);

                            String outcome;
                            String filtered;
                            try {
                                XmlValidationCertificateQualification result = block.execute();
                                List<String> filteredLabels = new ArrayList<>();
                                for (TrustServiceWrapper w : block.getFilteredServices()) {
                                    List<String> names = w.getServiceNames();
                                    filteredLabels.add(names == null || names.isEmpty() ? null : names.get(0));
                                }
                                filtered = arr(filteredLabels);
                                outcome = "{\"title\":" + esc(result.getTitle())
                                        + ",\"certificateQualification\":"
                                        + esc(result.getCertificateQualification() == null ? null
                                                : result.getCertificateQualification().name())
                                        + ",\"validationTime\":"
                                        + esc(result.getValidationTime() == null ? null : result.getValidationTime().name())
                                        + ",\"constraints\":" + constraints(result.getConstraint())
                                        + ",\"conclusion\":" + conclusion(result.getConclusion()) + "}";
                            } catch (RuntimeException e) {
                                outcome = esc("throws:" + e.getClass().getSimpleName());
                                filtered = "null";
                            }

                            pw.println("{\"kind\":\"certQualAtTime\",\"i\":" + index++
                                    + ",\"cert\":" + esc(certLabel)
                                    + ",\"services\":" + arr(labels)
                                    + ",\"date\":" + (date == null ? "null" : date.toString())
                                    + ",\"validationTime\":" + esc(validationTime.name())
                                    + ",\"result\":" + outcome
                                    + ",\"filteredServices\":" + filtered + "}");
                        }
                    }
                }
            }
        }
        System.out.println("wrote " + out + " (" + index + " rows)");
    }
}

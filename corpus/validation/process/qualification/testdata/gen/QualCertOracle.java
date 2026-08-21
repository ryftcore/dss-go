/*
 * Java oracle driver for the certificate-driven half of the phase-8e
 * eu.europa.esig.dss.validation.process.qualification tree: the six
 * "…ByCertificate…EIDAS" strategies, QCTypeIdentifiers,
 * CertificateQualificationCalculator (certificate x trust service),
 * ServiceByCertificateTypeFilter and UniqueServiceFilter.
 *
 * Each is a pure function of a CertificateWrapper (plus, for the last three, a
 * TrustServiceWrapper), so this driver builds a cross-product of synthetic
 * XmlCertificates - QcCompliance x QcSSCD x QcTypes x QcCClegislation x
 * certificate policies x notBefore - and dumps the upstream answer.  The Go
 * replay (../../qual_cert_oracle_test.go) rebuilds the identical certificates
 * and asserts the identical answers.
 *
 * Rows are written to
 *   dss/validation/process/qualification/testdata/oracle/qual_cert.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle QualCertOracle.java
 *   java  -cp "$CP:/tmp/oracle" QualCertOracle <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.TrustServiceWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateExtension;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificatePolicies;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificatePolicy;
import eu.europa.esig.dss.diagnostic.jaxb.XmlOID;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcCompliance;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcSSCD;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcStatements;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQualifier;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustedList;
import eu.europa.esig.dss.enumerations.CertificateExtensionEnum;
import eu.europa.esig.dss.enumerations.CertificatePolicy;
import eu.europa.esig.dss.enumerations.CertificateQualifiedStatus;
import eu.europa.esig.dss.enumerations.QCTypeEnum;
import eu.europa.esig.dss.enumerations.ServiceQualification;
import eu.europa.esig.dss.validation.process.qualification.certificate.CertificateQualificationCalculator;
import eu.europa.esig.dss.validation.process.qualification.certificate.QCTypeIdentifiers;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qscd.QSCDStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qualified.QualificationStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.type.TypeStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.trust.ServiceTypeIdentifier;
import eu.europa.esig.dss.validation.process.qualification.trust.TrustServiceStatus;
import eu.europa.esig.dss.validation.process.qualification.trust.filter.TrustServiceFilter;
import eu.europa.esig.dss.validation.process.qualification.trust.filter.TrustServicesFilterFactory;

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

public class QualCertOracle {

    private static final long PRE_EIDAS = 1370044800000L;   // 2013-06-01
    private static final long POST_EIDAS = 1514764800000L;  // 2018-01-01

    private static final String QCT_ESIGN = QCTypeEnum.QCT_ESIGN.getOid();
    private static final String QCT_ESEAL = QCTypeEnum.QCT_ESEAL.getOid();
    private static final String QCT_WEB = QCTypeEnum.QCT_WEB.getOid();
    private static final String QCT_UNKNOWN = "0.4.0.1862.1.6.99";

    private static final String POLICY_QCP = CertificatePolicy.QCP_PUBLIC.getOid();
    private static final String POLICY_QCP_PLUS = CertificatePolicy.QCP_PUBLIC_WITH_SSCD.getOid();
    private static final String POLICY_OTHER = "1.2.3.4.5";

    private static final List<List<String>> QC_TYPE_SETS = Arrays.asList(
            null,
            Collections.<String>emptyList(),
            Arrays.asList(QCT_ESIGN),
            Arrays.asList(QCT_ESEAL),
            Arrays.asList(QCT_WEB),
            Arrays.asList(QCT_ESIGN, QCT_ESEAL),
            Arrays.asList(QCT_ESIGN, QCT_WEB),
            Arrays.asList(QCT_ESEAL, QCT_WEB),
            Arrays.asList(QCT_ESIGN, QCT_ESEAL, QCT_WEB),
            Arrays.asList(QCT_UNKNOWN),
            Arrays.asList(QCT_ESIGN, QCT_UNKNOWN));

    private static final List<List<String>> POLICY_SETS = Arrays.asList(
            null,
            Collections.<String>emptyList(),
            Arrays.asList(POLICY_QCP),
            Arrays.asList(POLICY_QCP_PLUS),
            Arrays.asList(POLICY_QCP, POLICY_QCP_PLUS),
            Arrays.asList(POLICY_OTHER));

    private static final List<List<String>> LEGISLATION_SETS = Arrays.asList(
            null,
            Collections.<String>emptyList(),
            Arrays.asList("CH"),
            Arrays.asList("CH", "US"));

    private static final class CertConfig {
        final Boolean qcCompliance;
        final Boolean qcSSCD;
        final List<String> qcTypes;
        final List<String> legislations;
        final List<String> policies;
        final long notBefore;
        final boolean qcStatementsPresent;

        CertConfig(Boolean qcCompliance, Boolean qcSSCD, List<String> qcTypes, List<String> legislations,
                   List<String> policies, long notBefore, boolean qcStatementsPresent) {
            this.qcCompliance = qcCompliance;
            this.qcSSCD = qcSSCD;
            this.qcTypes = qcTypes;
            this.legislations = legislations;
            this.policies = policies;
            this.notBefore = notBefore;
            this.qcStatementsPresent = qcStatementsPresent;
        }
    }

    private static XmlCertificate buildCertificate(CertConfig c) {
        XmlCertificate cert = new XmlCertificate();
        cert.setId("C-SYNTH");
        cert.setNotBefore(new Date(c.notBefore));
        cert.setNotAfter(new Date(c.notBefore + 3L * 365 * 86400000L));

        List<XmlCertificateExtension> extensions = new ArrayList<>();
        if (c.qcStatementsPresent) {
            XmlQcStatements qcStatements = new XmlQcStatements();
            qcStatements.setOID(CertificateExtensionEnum.QC_STATEMENTS.getOid());
            if (c.qcCompliance != null) {
                XmlQcCompliance compliance = new XmlQcCompliance();
                compliance.setPresent(c.qcCompliance);
                qcStatements.setQcCompliance(compliance);
            }
            if (c.qcSSCD != null) {
                XmlQcSSCD sscd = new XmlQcSSCD();
                sscd.setPresent(c.qcSSCD);
                qcStatements.setQcSSCD(sscd);
            }
            if (c.qcTypes != null) {
                for (String oid : c.qcTypes) {
                    XmlOID xmlOid = new XmlOID();
                    xmlOid.setValue(oid);
                    qcStatements.getQcTypes().add(xmlOid);
                }
            }
            if (c.legislations != null) {
                qcStatements.getQcCClegislation().addAll(c.legislations);
            }
            extensions.add(qcStatements);
        }
        if (c.policies != null) {
            XmlCertificatePolicies certificatePolicies = new XmlCertificatePolicies();
            certificatePolicies.setOID(CertificateExtensionEnum.CERTIFICATE_POLICIES.getOid());
            for (String oid : c.policies) {
                XmlCertificatePolicy policy = new XmlCertificatePolicy();
                policy.setValue(oid);
                certificatePolicies.getCertificatePolicy().add(policy);
            }
            extensions.add(certificatePolicies);
        }
        cert.getCertificateExtensions().addAll(extensions);
        return cert;
    }

    private static List<CertConfig> certConfigs() {
        List<CertConfig> result = new ArrayList<>();
        // no QcStatements extension at all
        for (long notBefore : new long[]{PRE_EIDAS, POST_EIDAS}) {
            for (List<String> policies : POLICY_SETS) {
                result.add(new CertConfig(null, null, null, null, policies, notBefore, false));
            }
        }
        // QcStatements present, cross-product of the flags
        for (long notBefore : new long[]{PRE_EIDAS, POST_EIDAS}) {
            for (Boolean compliance : new Boolean[]{null, Boolean.FALSE, Boolean.TRUE}) {
                for (Boolean sscd : new Boolean[]{null, Boolean.FALSE, Boolean.TRUE}) {
                    for (List<String> qcTypes : QC_TYPE_SETS) {
                        result.add(new CertConfig(compliance, sscd, qcTypes, null,
                                Collections.<String>emptyList(), notBefore, true));
                    }
                }
            }
        }
        // legislation x compliance (QualificationByCertificate*'s second conjunct)
        for (long notBefore : new long[]{PRE_EIDAS, POST_EIDAS}) {
            for (Boolean compliance : new Boolean[]{Boolean.FALSE, Boolean.TRUE}) {
                for (List<String> legislations : LEGISLATION_SETS) {
                    for (List<String> policies : POLICY_SETS) {
                        result.add(new CertConfig(compliance, Boolean.FALSE, Arrays.asList(QCT_ESIGN),
                                legislations, policies, notBefore, true));
                    }
                }
            }
        }
        // policies x compliance at both sides of the eIDAS boundary
        for (long notBefore : new long[]{PRE_EIDAS, POST_EIDAS}) {
            for (Boolean compliance : new Boolean[]{Boolean.FALSE, Boolean.TRUE}) {
                for (Boolean sscd : new Boolean[]{Boolean.FALSE, Boolean.TRUE}) {
                    for (List<String> policies : POLICY_SETS) {
                        result.add(new CertConfig(compliance, sscd, Collections.<String>emptyList(),
                                null, policies, notBefore, true));
                    }
                }
            }
        }
        return result;
    }

    /** The trust services each certificate is combined with. */
    private static final class SvcConfig {
        final String label;
        final String status;
        final String type;
        final Long startDate;
        final List<String> qualifiers;
        final List<String> asis;

        SvcConfig(String label, String status, String type, Long startDate, List<String> qualifiers, List<String> asis) {
            this.label = label;
            this.status = status;
            this.type = type;
            this.startDate = startDate;
            this.qualifiers = qualifiers;
            this.asis = asis;
        }
    }

    private static final List<SvcConfig> SERVICES = Arrays.asList(
            new SvcConfig("none", null, null, null, null, null),
            new SvcConfig("granted-post-plain", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Collections.<String>emptyList(), Collections.<String>emptyList()),
            new SvcConfig("granted-post-qcstatement", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_STATEMENT.getUri()), Collections.<String>emptyList()),
            new SvcConfig("granted-post-notqualified", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.NOT_QUALIFIED.getUri()), Collections.<String>emptyList()),
            new SvcConfig("granted-post-withqscd", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_WITH_QSCD.getUri()), Collections.<String>emptyList()),
            new SvcConfig("granted-post-noqscd", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_NO_QSCD.getUri()), Collections.<String>emptyList()),
            new SvcConfig("granted-post-qscdasincert", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_QSCD_STATUS_AS_IN_CERT.getUri()), Collections.<String>emptyList()),
            new SvcConfig("granted-post-foreseal", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_FOR_ESEAL.getUri()),
                    Arrays.asList("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSeals")),
            new SvcConfig("granted-post-forwsa", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_FOR_WSA.getUri()),
                    Arrays.asList("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForWebSiteAuthentication")),
            new SvcConfig("granted-pre-plain", TrustServiceStatus.UNDER_SUPERVISION.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), PRE_EIDAS,
                    Collections.<String>emptyList(), Collections.<String>emptyList()),
            new SvcConfig("granted-pre-withsscd", TrustServiceStatus.UNDER_SUPERVISION.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), PRE_EIDAS,
                    Arrays.asList(ServiceQualification.QC_WITH_SSCD.getUri()), Collections.<String>emptyList()),
            new SvcConfig("withdrawn-post", TrustServiceStatus.WITHDRAWN.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Collections.<String>emptyList(), Collections.<String>emptyList()),
            new SvcConfig("granted-post-esig-asi", TrustServiceStatus.GRANTED.getUri(),
                    ServiceTypeIdentifier.CA_QC.getUri(), POST_EIDAS,
                    Arrays.asList(ServiceQualification.QC_FOR_ESIG.getUri()),
                    Arrays.asList("http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures")));

    private static TrustServiceWrapper buildService(SvcConfig c) {
        if (c.status == null && c.type == null && c.startDate == null) {
            return null;
        }
        TrustServiceWrapper w = new TrustServiceWrapper();
        w.setStatus(c.status);
        w.setType(c.type);
        w.setStartDate(c.startDate == null ? null : new Date(c.startDate));
        if (c.qualifiers != null) {
            List<XmlQualifier> qualifiers = new ArrayList<>();
            for (String uri : c.qualifiers) {
                XmlQualifier q = new XmlQualifier();
                q.setValue(uri);
                qualifiers.add(q);
            }
            w.setCapturedQualifiers(qualifiers);
        }
        if (c.asis != null) {
            w.setAdditionalServiceInfos(new ArrayList<>(c.asis));
        }
        w.setCountryCode("LU");
        XmlTrustedList tl = new XmlTrustedList();
        tl.setUrl("https://tl.example/LU");
        w.setTrustedList(tl);
        return w;
    }

    // ---- JSON helpers (same shape as QualTrustOracle) ---------------------

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

    private static String en(Enum<?> e) {
        return e == null ? "null" : esc(e.name());
    }

    public static void main(String[] args) throws IOException {
        if (args.length < 1) {
            System.err.println("usage: QualCertOracle <dss-repo-root>");
            System.exit(2);
        }
        Path out = Paths.get(args[0], "validation", "process", "qualification", "testdata", "oracle", "qual_cert.jsonl");
        Files.createDirectories(out.getParent());
        int index = 0;
        try (PrintWriter pw = new PrintWriter(Files.newBufferedWriter(out, StandardCharsets.UTF_8))) {
            for (CertConfig cc : certConfigs()) {
                XmlCertificate xml = buildCertificate(cc);
                CertificateWrapper cert = new CertificateWrapper(xml);

                StringBuilder sb = new StringBuilder();
                sb.append("{\"kind\":\"cert\",\"i\":").append(index++);
                sb.append(",\"in\":{\"qcCompliance\":").append(cc.qcCompliance == null ? "null" : cc.qcCompliance.toString());
                sb.append(",\"qcSSCD\":").append(cc.qcSSCD == null ? "null" : cc.qcSSCD.toString());
                sb.append(",\"qcTypes\":").append(arr(cc.qcTypes));
                sb.append(",\"legislations\":").append(arr(cc.legislations));
                sb.append(",\"policies\":").append(arr(cc.policies));
                sb.append(",\"notBefore\":").append(cc.notBefore);
                sb.append(",\"qcStatementsPresent\":").append(cc.qcStatementsPresent);
                sb.append("}");

                sb.append(",\"qcTypeIdentifiers\":{");
                sb.append("\"esign\":").append(QCTypeIdentifiers.isQCTypeEsign(cert));
                sb.append(",\"eseal\":").append(QCTypeIdentifiers.isQCTypeEseal(cert));
                sb.append(",\"web\":").append(QCTypeIdentifiers.isQCTypeWeb(cert));
                sb.append("}");

                sb.append(",\"fromCert\":{");
                sb.append("\"qualification\":").append(en(QualificationStrategyFactory.createQualificationFromCert(cert).getQualifiedStatus()));
                sb.append(",\"type\":").append(en(TypeStrategyFactory.createTypeFromCert(cert).getType()));
                sb.append(",\"qscd\":").append(en(QSCDStrategyFactory.createQSCDFromCert(cert).getQSCDStatus()));
                sb.append("}");

                // full CertificateQualificationCalculator per trust service, plus
                // the certificate-type filter's verdict on that service.
                sb.append(",\"perService\":{");
                boolean first = true;
                for (SvcConfig sc : SERVICES) {
                    TrustServiceWrapper svc = buildService(sc);
                    if (!first) {
                        sb.append(',');
                    }
                    first = false;
                    sb.append(esc(sc.label)).append(":{");
                    sb.append("\"qualification\":").append(en(new CertificateQualificationCalculator(cert, svc).getQualification()));
                    if (svc == null) {
                        sb.append(",\"byCertificateType\":null");
                    } else {
                        TrustServiceFilter filter = TrustServicesFilterFactory.createFilterByCertificateType(cert);
                        List<TrustServiceWrapper> filtered = filter.filter(Collections.singletonList(svc));
                        sb.append(",\"byCertificateType\":").append(filtered != null && !filtered.isEmpty());
                    }
                    // full "from cert and TL" strategy triple
                    CertificateQualifiedStatus qualified =
                            QualificationStrategyFactory.createQualificationFromCertAndTL(cert, svc).getQualifiedStatus();
                    sb.append(",\"qualifiedStatus\":").append(en(qualified));
                    sb.append(",\"certType\":").append(en(TypeStrategyFactory.createTypeFromCertAndTL(cert, svc, qualified).getType()));
                    sb.append(",\"qscdStatus\":").append(en(QSCDStrategyFactory.createQSCDFromCertAndTL(cert, svc, qualified).getQSCDStatus()));
                    sb.append("}");
                }
                sb.append("}");

                // UniqueServiceFilter over the whole service list, and over each
                // adjacent pair, so both the "one conclusion" and the "several
                // conclusions" branches are observed.
                List<TrustServiceWrapper> all = new ArrayList<>();
                List<String> allLabels = new ArrayList<>();
                for (SvcConfig sc : SERVICES) {
                    TrustServiceWrapper svc = buildService(sc);
                    if (svc != null) {
                        svc.setServiceNames(Collections.singletonList(sc.label));
                        all.add(svc);
                        allLabels.add(sc.label);
                    }
                }
                sb.append(",\"uniqueAll\":").append(uniqueLabels(cert, all));
                sb.append(",\"uniquePairs\":[");
                for (int i = 0; i + 1 < all.size(); i++) {
                    if (i > 0) {
                        sb.append(',');
                    }
                    sb.append(uniqueLabels(cert, Arrays.asList(all.get(i), all.get(i + 1))));
                }
                sb.append("]");
                sb.append(",\"uniqueSingleton\":").append(uniqueLabels(cert, Collections.singletonList(all.get(0))));
                sb.append(",\"uniqueEmpty\":").append(uniqueLabels(cert, Collections.<TrustServiceWrapper>emptyList()));
                sb.append(",\"serviceLabels\":").append(arr(allLabels));
                sb.append("}");
                pw.println(sb);
            }
        }
        System.out.println("wrote " + out);
    }

    private static String uniqueLabels(CertificateWrapper cert, List<TrustServiceWrapper> services) {
        TrustServiceFilter filter = TrustServicesFilterFactory.createUniqueServiceFilter(cert);
        List<TrustServiceWrapper> filtered = filter.filter(services);
        List<String> labels = new ArrayList<>();
        for (TrustServiceWrapper w : filtered) {
            List<String> names = w.getServiceNames();
            labels.add(names == null || names.isEmpty() ? null : names.get(0));
        }
        return arr(labels);
    }
}

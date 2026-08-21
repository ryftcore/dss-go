/*
 * Java oracle driver for the trust-service half of the phase-8e
 * eu.europa.esig.dss.validation.process.qualification tree: the ten
 * TrustServiceChecker consistency predicates (and the twelve
 * TrustServiceCondition classes behind them), the thirteen
 * TrustServicesFilterFactory filters, the three "…ByTL" qualification
 * strategies, EIDASUtils, TrustServiceStatus, ServiceTypeIdentifier and the
 * three qualification matrices.
 *
 * Every one of those is a pure function of a TrustServiceWrapper (plus enum
 * inputs), so this driver builds a designed cross-product of
 * TrustServiceWrappers - covering both branches of every predicate - and dumps
 * the upstream answer for each.  The Go replay
 * (../../qual_trust_oracle_test.go) rebuilds the identical wrappers from the
 * same JSONL rows and asserts the identical answers.
 *
 * Rows are written to
 *   dss/validation/process/qualification/testdata/oracle/qual_trust.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle QualTrustOracle.java
 *   java  -cp "$CP:/tmp/oracle" QualTrustOracle <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.TrustServiceWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQualifier;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustedList;
import eu.europa.esig.dss.enumerations.AdditionalServiceInformation;
import eu.europa.esig.dss.enumerations.CertificateQualification;
import eu.europa.esig.dss.enumerations.CertificateQualifiedStatus;
import eu.europa.esig.dss.enumerations.CertificateType;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.QSCDStatus;
import eu.europa.esig.dss.enumerations.ServiceQualification;
import eu.europa.esig.dss.enumerations.SignatureQualification;
import eu.europa.esig.dss.validation.process.qualification.EIDASUtils;
import eu.europa.esig.dss.validation.process.qualification.certificate.CertQualificationMatrix;
import eu.europa.esig.dss.validation.process.qualification.certificate.FinalCertificateQualificationCalculator;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qscd.QSCDStrategy;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qscd.QSCDStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qualified.QualificationStrategy;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.qualified.QualificationStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.type.TypeStrategy;
import eu.europa.esig.dss.validation.process.qualification.certificate.checks.type.TypeStrategyFactory;
import eu.europa.esig.dss.validation.process.qualification.signature.SigQualificationMatrix;
import eu.europa.esig.dss.validation.process.qualification.trust.ServiceTypeIdentifier;
import eu.europa.esig.dss.validation.process.qualification.trust.TrustServiceStatus;
import eu.europa.esig.dss.validation.process.qualification.trust.consistency.TrustServiceChecker;
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
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

public class QualTrustOracle {

    /** 2013-06-01: pre-eIDAS (EIDASUtils.EIDAS_DATE is 2016-07-01). */
    private static final Date PRE_EIDAS = date(1370044800000L);
    /** 2016-06-01: still pre-eIDAS, but after the "grace" boundary. */
    private static final Date PRE_EIDAS_LATE = date(1464739200000L);
    /** 2018-01-01: post-eIDAS. */
    private static final Date POST_EIDAS = date(1514764800000L);
    /** 2020-01-01: post-eIDAS, later. */
    private static final Date POST_EIDAS_LATE = date(1577836800000L);

    private static Date date(long millis) {
        return new Date(millis);
    }

    // ---- the qualifier / ASI / status / type vocabularies -----------------

    private static final String Q_QC_STATEMENT = ServiceQualification.QC_STATEMENT.getUri();
    private static final String Q_NOT_QUALIFIED = ServiceQualification.NOT_QUALIFIED.getUri();
    private static final String Q_WITH_SSCD = ServiceQualification.QC_WITH_SSCD.getUri();
    private static final String Q_WITH_QSCD = ServiceQualification.QC_WITH_QSCD.getUri();
    private static final String Q_NO_SSCD = ServiceQualification.QC_NO_SSCD.getUri();
    private static final String Q_NO_QSCD = ServiceQualification.QC_NO_QSCD.getUri();
    private static final String Q_SSCD_AS_IN_CERT = ServiceQualification.QC_SSCD_STATUS_AS_IN_CERT.getUri();
    private static final String Q_QSCD_AS_IN_CERT = ServiceQualification.QC_QSCD_STATUS_AS_IN_CERT.getUri();
    private static final String Q_QSCD_ON_BEHALF = ServiceQualification.QC_QSCD_MANAGED_ON_BEHALF.getUri();
    private static final String Q_LEGAL_PERSON = ServiceQualification.QC_FOR_LEGAL_PERSON.getUri();
    private static final String Q_FOR_ESIG = ServiceQualification.QC_FOR_ESIG.getUri();
    private static final String Q_FOR_ESEAL = ServiceQualification.QC_FOR_ESEAL.getUri();
    private static final String Q_FOR_WSA = ServiceQualification.QC_FOR_WSA.getUri();
    private static final String Q_UNKNOWN = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/NotAKnownQualifier";

    private static final String ASI_ESIG = AdditionalServiceInformation.FOR_ESIGNATURES.getUri();
    private static final String ASI_ESEAL = AdditionalServiceInformation.FOR_ESEALS.getUri();
    private static final String ASI_WSA = AdditionalServiceInformation.FOR_WEB_AUTHENTICATION.getUri();
    private static final String ASI_UNKNOWN = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForSomethingElse";

    /** The qualifier bundles the cross-product walks. */
    private static final List<List<String>> QUALIFIER_SETS = Arrays.asList(
            Collections.<String>emptyList(),
            Arrays.asList(Q_QC_STATEMENT),
            Arrays.asList(Q_NOT_QUALIFIED),
            Arrays.asList(Q_QC_STATEMENT, Q_NOT_QUALIFIED),
            Arrays.asList(Q_WITH_SSCD),
            Arrays.asList(Q_WITH_QSCD),
            Arrays.asList(Q_NO_SSCD),
            Arrays.asList(Q_NO_QSCD),
            Arrays.asList(Q_WITH_QSCD, Q_NO_QSCD),
            Arrays.asList(Q_WITH_SSCD, Q_NO_SSCD),
            Arrays.asList(Q_SSCD_AS_IN_CERT),
            Arrays.asList(Q_QSCD_AS_IN_CERT),
            Arrays.asList(Q_QSCD_AS_IN_CERT, Q_WITH_QSCD),
            Arrays.asList(Q_QSCD_ON_BEHALF),
            Arrays.asList(Q_QSCD_ON_BEHALF, Q_NO_QSCD),
            Arrays.asList(Q_LEGAL_PERSON),
            Arrays.asList(Q_LEGAL_PERSON, Q_FOR_ESIG),
            Arrays.asList(Q_LEGAL_PERSON, Q_FOR_ESEAL),
            Arrays.asList(Q_FOR_ESIG),
            Arrays.asList(Q_FOR_ESEAL),
            Arrays.asList(Q_FOR_WSA),
            Arrays.asList(Q_FOR_ESIG, Q_FOR_ESEAL),
            Arrays.asList(Q_FOR_ESIG, Q_FOR_WSA),
            Arrays.asList(Q_FOR_ESIG, Q_FOR_ESEAL, Q_FOR_WSA),
            Arrays.asList(Q_UNKNOWN),
            Arrays.asList(Q_QC_STATEMENT, Q_UNKNOWN),
            Arrays.asList(Q_QC_STATEMENT, Q_WITH_QSCD, Q_FOR_ESIG),
            Arrays.asList(Q_NOT_QUALIFIED, Q_WITH_QSCD),
            Arrays.asList(Q_QC_STATEMENT, Q_SSCD_AS_IN_CERT, Q_LEGAL_PERSON));

    /** The additional-service-information bundles the cross-product walks. */
    private static final List<List<String>> ASI_SETS = Arrays.asList(
            null,
            Collections.<String>emptyList(),
            Arrays.asList(ASI_ESIG),
            Arrays.asList(ASI_ESEAL),
            Arrays.asList(ASI_WSA),
            Arrays.asList(ASI_ESIG, ASI_ESEAL),
            Arrays.asList(ASI_ESIG, ASI_WSA),
            Arrays.asList(ASI_ESIG, ASI_ESEAL, ASI_WSA),
            Arrays.asList(ASI_UNKNOWN));

    private static final List<String> STATUSES = buildStatuses();

    private static List<String> buildStatuses() {
        List<String> result = new ArrayList<>();
        for (TrustServiceStatus status : TrustServiceStatus.values()) {
            result.add(status.getUri());
        }
        result.add("http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/notAKnownStatus");
        result.add(null);
        return result;
    }

    private static final List<String> TYPES = buildTypes();

    private static List<String> buildTypes() {
        List<String> result = new ArrayList<>();
        for (ServiceTypeIdentifier type : ServiceTypeIdentifier.values()) {
            result.add(type.getUri());
        }
        result.add("http://uri.etsi.org/TrstSvc/Svctype/NotAKnownType");
        result.add(null);
        return result;
    }

    private static final List<Date> START_DATES = Arrays.asList(null, PRE_EIDAS, PRE_EIDAS_LATE, POST_EIDAS, POST_EIDAS_LATE);

    // ---- stub strategies --------------------------------------------------

    private static QualificationStrategy qcStub(final CertificateQualifiedStatus status) {
        return new QualificationStrategy() {
            @Override
            public CertificateQualifiedStatus getQualifiedStatus() {
                return status;
            }
        };
    }

    private static TypeStrategy typeStub(final CertificateType type) {
        return new TypeStrategy() {
            @Override
            public CertificateType getType() {
                return type;
            }
        };
    }

    private static QSCDStrategy qscdStub(final QSCDStatus status) {
        return new QSCDStrategy() {
            @Override
            public QSCDStatus getQSCDStatus() {
                return status;
            }
        };
    }

    // ---- JSON helpers -----------------------------------------------------

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

    private static String num(Date d) {
        return d == null ? "null" : Long.toString(d.getTime());
    }

    private static String en(Enum<?> e) {
        return e == null ? "null" : esc(e.name());
    }

    // ---- the service under test ------------------------------------------

    /** One designed TrustServiceWrapper configuration. */
    private static final class Config {
        final String status;
        final String type;
        final Date startDate;
        final Date endDate;
        final List<String> qualifiers;
        final List<String> asis;
        final String countryCode;
        final String tlUrl;
        final Boolean enactedMRA;
        final Date mraStart;
        final Date mraEnd;

        Config(String status, String type, Date startDate, Date endDate, List<String> qualifiers,
               List<String> asis, String countryCode, String tlUrl, Boolean enactedMRA, Date mraStart, Date mraEnd) {
            this.status = status;
            this.type = type;
            this.startDate = startDate;
            this.endDate = endDate;
            this.qualifiers = qualifiers;
            this.asis = asis;
            this.countryCode = countryCode;
            this.tlUrl = tlUrl;
            this.enactedMRA = enactedMRA;
            this.mraStart = mraStart;
            this.mraEnd = mraEnd;
        }
    }

    private static TrustServiceWrapper build(Config c) {
        TrustServiceWrapper w = new TrustServiceWrapper();
        w.setStatus(c.status);
        w.setType(c.type);
        w.setStartDate(c.startDate);
        w.setEndDate(c.endDate);
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
        w.setCountryCode(c.countryCode);
        XmlTrustedList tl = new XmlTrustedList();
        tl.setUrl(c.tlUrl);
        w.setTrustedList(tl);
        w.setEnactedMRA(c.enactedMRA);
        w.setMraTrustServiceEquivalenceStatusStartingTime(c.mraStart);
        w.setMraTrustServiceEquivalenceStatusEndingTime(c.mraEnd);
        return w;
    }

    private static List<Config> configs() {
        List<Config> result = new ArrayList<>();
        // (a) status x startDate: GrantedServiceFilter, pre-eIDAS status consistency.
        for (String status : STATUSES) {
            for (Date start : START_DATES) {
                result.add(new Config(status, ServiceTypeIdentifier.CA_QC.getUri(), start, null,
                        Collections.<String>emptyList(), Collections.<String>emptyList(), "LU",
                        "https://tl.example/LU", null, null, null));
            }
        }
        // (b) qualifiers x startDate: every consistency predicate + the ByTL strategies.
        for (List<String> qualifiers : QUALIFIER_SETS) {
            for (Date start : START_DATES) {
                result.add(new Config(TrustServiceStatus.GRANTED.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                        start, null, qualifiers, Collections.<String>emptyList(), "LU",
                        "https://tl.example/LU", null, null, null));
            }
        }
        // (c) qualifiers x ASI: usage / qualifier-vs-ASI consistency.
        for (List<String> qualifiers : QUALIFIER_SETS) {
            for (List<String> asis : ASI_SETS) {
                result.add(new Config(TrustServiceStatus.GRANTED.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                        POST_EIDAS, null, qualifiers, asis, "LU", "https://tl.example/LU", null, null, null));
                result.add(new Config(TrustServiceStatus.UNDER_SUPERVISION.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                        PRE_EIDAS, null, qualifiers, asis, "LU", "https://tl.example/LU", null, null, null));
            }
        }
        // (d) type x status: the CA/QC, QTST and QEAA filters.
        for (String type : TYPES) {
            result.add(new Config(TrustServiceStatus.GRANTED.getUri(), type, POST_EIDAS, null,
                    Collections.<String>emptyList(), Collections.<String>emptyList(), "DE",
                    "https://tl.example/DE", null, null, null));
        }
        // (e) validity window: ServiceByDateFilter.
        Date[] windows = {null, PRE_EIDAS, PRE_EIDAS_LATE, POST_EIDAS, POST_EIDAS_LATE};
        for (Date start : windows) {
            for (Date end : windows) {
                result.add(new Config(TrustServiceStatus.GRANTED.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                        start, end, Collections.<String>emptyList(), Collections.<String>emptyList(), "FR",
                        "https://tl.example/FR", null, null, null));
            }
        }
        // (f) country / URL casing: ServiceByCountryFilter, ServiceByTLUrlFilter.
        for (String country : new String[]{"LU", "lu", "DE", "", null}) {
            for (String url : new String[]{"https://tl.example/LU", "HTTPS://TL.EXAMPLE/lu", "", null}) {
                result.add(new Config(TrustServiceStatus.GRANTED.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                        POST_EIDAS, null, Collections.<String>emptyList(), Collections.<String>emptyList(),
                        country, url, null, null, null));
            }
        }
        // (g) MRA: ServiceByMRAEnactedFilter, ServiceByMRAEquivalenceStartingDateFilter.
        for (Boolean mra : new Boolean[]{null, Boolean.FALSE, Boolean.TRUE}) {
            for (Date mraStart : windows) {
                for (Date mraEnd : new Date[]{null, POST_EIDAS_LATE}) {
                    result.add(new Config(TrustServiceStatus.GRANTED.getUri(), ServiceTypeIdentifier.CA_QC.getUri(),
                            POST_EIDAS, null, Collections.<String>emptyList(), Collections.<String>emptyList(),
                            "LU", "https://tl.example/LU", mra, mraStart, mraEnd));
                }
            }
        }
        return result;
    }

    /** Probe dates the date-sensitive filters are evaluated at. */
    private static final List<Date> PROBE_DATES = Arrays.asList(null, PRE_EIDAS, PRE_EIDAS_LATE, POST_EIDAS, POST_EIDAS_LATE);

    public static void main(String[] args) throws IOException {
        if (args.length < 1) {
            System.err.println("usage: QualTrustOracle <dss-repo-root>");
            System.exit(2);
        }
        Path out = Paths.get(args[0], "validation", "process", "qualification", "testdata", "oracle", "qual_trust.jsonl");
        Files.createDirectories(out.getParent());
        try (PrintWriter pw = new PrintWriter(Files.newBufferedWriter(out, StandardCharsets.UTF_8))) {
            writeServiceRows(pw);
            writeMatrixRows(pw);
            writeEidasRows(pw);
            writeVocabularyRows(pw);
        }
        System.out.println("wrote " + out);
    }

    private static void writeServiceRows(PrintWriter pw) {
        int index = 0;
        for (Config c : configs()) {
            TrustServiceWrapper w = build(c);
            List<TrustServiceWrapper> single = Collections.singletonList(w);

            StringBuilder sb = new StringBuilder();
            sb.append("{\"kind\":\"service\",\"i\":").append(index++);
            sb.append(",\"in\":{\"status\":").append(esc(c.status));
            sb.append(",\"type\":").append(esc(c.type));
            sb.append(",\"startDate\":").append(num(c.startDate));
            sb.append(",\"endDate\":").append(num(c.endDate));
            sb.append(",\"qualifiers\":").append(arr(c.qualifiers));
            sb.append(",\"asis\":").append(arr(c.asis));
            sb.append(",\"countryCode\":").append(esc(c.countryCode));
            sb.append(",\"tlUrl\":").append(esc(c.tlUrl));
            sb.append(",\"enactedMRA\":").append(c.enactedMRA == null ? "null" : c.enactedMRA.toString());
            sb.append(",\"mraStart\":").append(num(c.mraStart));
            sb.append(",\"mraEnd\":").append(num(c.mraEnd));
            sb.append("}");

            // the ten TrustServiceChecker predicates
            sb.append(",\"checker\":{");
            sb.append("\"legalPerson\":").append(TrustServiceChecker.isLegalPersonConsistent(w));
            sb.append(",\"qcStatement\":").append(TrustServiceChecker.isQCStatementConsistent(w));
            sb.append(",\"qscd\":").append(TrustServiceChecker.isQSCDConsistent(w));
            sb.append(",\"qscdStatusAsInCert\":").append(TrustServiceChecker.isQSCDStatusAsInCertConsistent(w));
            sb.append(",\"postEidasQscd\":").append(TrustServiceChecker.isPostEIDASQSCDConsistent(w));
            sb.append(",\"qualifiersListKnown\":").append(TrustServiceChecker.isQualifiersListKnownConsistent(w));
            sb.append(",\"usage\":").append(TrustServiceChecker.isUsageConsistent(w));
            sb.append(",\"preEidasStatus\":").append(TrustServiceChecker.isPreEIDASStatusConsistent(w));
            sb.append(",\"preEidasQualifierAsi\":").append(TrustServiceChecker.isPreEIDASQualifierAndAdditionalServiceInfoConsistent(w));
            sb.append(",\"qualifierAsi\":").append(TrustServiceChecker.isQualifierAndAdditionalServiceInfoConsistent(w));
            sb.append("}");

            // the date-independent filters, as accept/reject over a singleton list
            sb.append(",\"filters\":{");
            sb.append("\"granted\":").append(accepts(TrustServicesFilterFactory.createFilterByGranted(), single));
            sb.append(",\"caQc\":").append(accepts(TrustServicesFilterFactory.createFilterByCaQc(), single));
            sb.append(",\"qtst\":").append(accepts(TrustServicesFilterFactory.createFilterByQTST(), single));
            sb.append(",\"qeaa\":").append(accepts(TrustServicesFilterFactory.createFilterByQEAA(), single));
            sb.append(",\"consistentStatus\":").append(accepts(TrustServicesFilterFactory.createConsistentServiceByStatusFilter(), single));
            sb.append(",\"consistentQC\":").append(accepts(TrustServicesFilterFactory.createConsistentServiceByQCFilter(), single));
            sb.append(",\"consistentQSCD\":").append(accepts(TrustServicesFilterFactory.createConsistentServiceByQSCDFilter(), single));
            sb.append(",\"consistentCertType\":").append(accepts(TrustServicesFilterFactory.createConsistentServiceByCertificateTypeFilter(), single));
            sb.append(",\"mraEnacted\":").append(accepts(TrustServicesFilterFactory.createMRAEnactedFilter(), single));
            sb.append(",\"countryLU\":").append(accepts(TrustServicesFilterFactory.createFilterByCountry("LU"), single));
            sb.append(",\"countriesLUDE\":").append(accepts(TrustServicesFilterFactory.createFilterByCountries(setOf("LU", "DE")), single));
            sb.append(",\"urls\":").append(accepts(TrustServicesFilterFactory.createFilterByUrls(setOf("https://tl.example/LU", "https://tl.example/FR")), single));
            sb.append("}");

            // the date-sensitive filters, probed at each PROBE_DATES entry
            sb.append(",\"byDate\":[");
            for (int i = 0; i < PROBE_DATES.size(); i++) {
                if (i > 0) {
                    sb.append(',');
                }
                sb.append(accepts(TrustServicesFilterFactory.createFilterByDate(PROBE_DATES.get(i)), single));
            }
            sb.append("],\"byMraEquivalenceStartingDate\":[");
            for (int i = 0; i < PROBE_DATES.size(); i++) {
                if (i > 0) {
                    sb.append(',');
                }
                sb.append(accepts(TrustServicesFilterFactory.createFilterByMRAEquivalenceStartingDate(PROBE_DATES.get(i)), single));
            }
            sb.append("]");

            // QualificationByTL over each stub input
            sb.append(",\"qualificationByTL\":{");
            sb.append("\"QC\":").append(en(QualificationStrategyFactory.createQualificationFromTL(w, qcStub(CertificateQualifiedStatus.QC)).getQualifiedStatus()));
            sb.append(",\"NOT_QC\":").append(en(QualificationStrategyFactory.createQualificationFromTL(w, qcStub(CertificateQualifiedStatus.NOT_QC)).getQualifiedStatus()));
            sb.append("}");

            // TypeByTL over each (qualified, typeInCert) stub pair
            sb.append(",\"typeByTL\":{");
            boolean first = true;
            for (CertificateQualifiedStatus qualified : CertificateQualifiedStatus.values()) {
                for (CertificateType inCert : CertificateType.values()) {
                    if (!first) {
                        sb.append(',');
                    }
                    first = false;
                    sb.append(esc(qualified.name() + "|" + inCert.name())).append(':')
                            .append(en(TypeStrategyFactory.createTypeFromTL(w, qualified, typeStub(inCert)).getType()));
                }
            }
            sb.append("}");

            // QSCDByTL over each (qualified, qscdInCert) stub pair
            sb.append(",\"qscdByTL\":{");
            first = true;
            for (CertificateQualifiedStatus qualified : CertificateQualifiedStatus.values()) {
                for (QSCDStatus inCert : QSCDStatus.values()) {
                    if (!first) {
                        sb.append(',');
                    }
                    first = false;
                    sb.append(esc(qualified.name() + "|" + inCert.name())).append(':')
                            .append(en(QSCDStrategyFactory.createQSCDFromTL(w, qualified, qscdStub(inCert)).getQSCDStatus()));
                }
            }
            sb.append("}}");
            pw.println(sb);
        }
    }

    private static Set<String> setOf(String... values) {
        return new LinkedHashSet<>(Arrays.asList(values));
    }

    private static boolean accepts(TrustServiceFilter filter, List<TrustServiceWrapper> single) {
        List<TrustServiceWrapper> filtered = filter.filter(single);
        return filtered != null && !filtered.isEmpty();
    }

    private static void writeMatrixRows(PrintWriter pw) {
        // CertQualificationMatrix: the full 2 x 4+ x 2 cube.
        for (CertificateQualifiedStatus qc : CertificateQualifiedStatus.values()) {
            for (CertificateType type : CertificateType.values()) {
                for (QSCDStatus qscd : QSCDStatus.values()) {
                    pw.println("{\"kind\":\"certQualificationMatrix\",\"qc\":" + en(qc)
                            + ",\"type\":" + en(type) + ",\"qscd\":" + en(qscd)
                            + ",\"out\":" + en(CertQualificationMatrix.getCertQualification(qc, type, qscd)) + "}");
                }
            }
        }
        // SigQualificationMatrix: every (supported Indication) x CertificateQualification cell.
        for (Indication indication : Indication.values()) {
            for (CertificateQualification cq : CertificateQualification.values()) {
                String out;
                try {
                    out = en(SigQualificationMatrix.getSignatureQualification(indication, cq));
                } catch (RuntimeException e) {
                    out = esc("throws:" + e.getClass().getSimpleName());
                }
                pw.println("{\"kind\":\"sigQualificationMatrix\",\"indication\":" + en(indication)
                        + ",\"certQualification\":" + en(cq) + ",\"out\":" + out + "}");
            }
        }
        // FinalCertificateQualificationCalculator: every ordered pair.
        for (CertificateQualification atIssuance : CertificateQualification.values()) {
            for (CertificateQualification atSigning : CertificateQualification.values()) {
                String out;
                try {
                    out = en(new FinalCertificateQualificationCalculator(atIssuance, atSigning).getFinalQualification());
                } catch (RuntimeException e) {
                    out = esc("throws:" + e.getClass().getSimpleName());
                }
                pw.println("{\"kind\":\"finalCertQualification\",\"atIssuance\":" + en(atIssuance)
                        + ",\"atSigning\":" + en(atSigning) + ",\"out\":" + out + "}");
            }
        }
    }

    private static void writeEidasRows(PrintWriter pw) {
        long[] probes = {
                Long.MIN_VALUE / 4,
                0L,
                1467324000000L - 1L,          // just before 2016-06-30T22:00:00Z
                1467324000000L,               // EIDAS_DATE  (2016-07-01 in Brussels)
                1467324000000L + 1L,
                1467324000000L + 86400000L,
                1519858800000L - 1L,          // just before EIDAS_GRACE_DATE
                1519858800000L,               // EIDAS_GRACE_DATE (2018-03-01 in Brussels)
                1519858800000L + 1L,
                Long.MAX_VALUE / 4,
        };
        for (long p : probes) {
            Date d = new Date(p);
            pw.println("{\"kind\":\"eidas\",\"millis\":" + p
                    + ",\"isPostEIDAS\":" + EIDASUtils.isPostEIDAS(d)
                    + ",\"isPreEIDAS\":" + EIDASUtils.isPreEIDAS(d)
                    + ",\"isPostGracePeriod\":" + EIDASUtils.isPostGracePeriod(d) + "}");
        }
        pw.println("{\"kind\":\"eidas\",\"millis\":null"
                + ",\"isPostEIDAS\":" + EIDASUtils.isPostEIDAS(null)
                + ",\"isPreEIDAS\":" + EIDASUtils.isPreEIDAS(null)
                + ",\"isPostGracePeriod\":" + EIDASUtils.isPostGracePeriod(null) + "}");
    }

    private static void writeVocabularyRows(PrintWriter pw) {
        for (TrustServiceStatus status : TrustServiceStatus.values()) {
            pw.println("{\"kind\":\"status\",\"name\":" + en(status)
                    + ",\"uri\":" + esc(status.getUri())
                    + ",\"shortName\":" + esc(status.getShortName())
                    + ",\"isPostEidas\":" + status.isPostEidas()
                    + ",\"isValid\":" + status.isValid() + "}");
        }
        List<String> statusProbes = new ArrayList<>(STATUSES);
        for (String uri : statusProbes) {
            pw.println("{\"kind\":\"statusProbe\",\"uri\":" + esc(uri)
                    + ",\"acceptableAfter\":" + TrustServiceStatus.isAcceptableStatusAfterEIDAS(uri)
                    + ",\"acceptableBefore\":" + TrustServiceStatus.isAcceptableStatusBeforeEIDAS(uri) + "}");
        }
        for (ServiceTypeIdentifier type : ServiceTypeIdentifier.values()) {
            pw.println("{\"kind\":\"serviceType\",\"name\":" + en(type) + ",\"uri\":" + esc(type.getUri()) + "}");
        }
        for (String uri : TYPES) {
            pw.println("{\"kind\":\"serviceTypeProbe\",\"uri\":" + esc(uri)
                    + ",\"isCaQc\":" + ServiceTypeIdentifier.isCaQc(uri)
                    + ",\"isQTST\":" + ServiceTypeIdentifier.isQTST(uri)
                    + ",\"isQEAA\":" + ServiceTypeIdentifier.isQEAA(uri) + "}");
        }
    }
}

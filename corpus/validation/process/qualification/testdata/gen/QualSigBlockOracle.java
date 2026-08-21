/*
 * Java oracle driver for the top-level signature/timestamp qualification
 * blocks of the phase-8e eu.europa.esig.dss.validation.process.qualification
 * tree: SignatureQualificationBlock and TimestampQualificationBlock.
 *
 * Both are driven over a designed cross-product of (signing certificate with
 * attached trust-service providers, list of XmlTLAnalysis, EN 319 102-1
 * conclusion / best-signature-time), and the WHOLE result tree is dumped: the
 * ordered Constraint list with each item's message keys, the Conclusion, the
 * final SignatureQualification / TimestampQualification, and, for signatures,
 * the two nested XmlValidationCertificateQualification blocks.
 *
 * The Go replay (../../qual_sig_block_oracle_test.go) rebuilds the identical
 * inputs and asserts the identical tree.  Among other things this pins the
 * "no acceptable trust service survives the filters" shape, on which upstream
 * reads getFilteredServices() as an EMPTY LIST (never null).
 *
 * Rows are written to
 *   dss/validation/process/qualification/testdata/oracle/qual_sig_block.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle QualSigBlockOracle.java
 *   java  -cp "$CP:/tmp/oracle" QualSigBlockOracle <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusionWithProofOfExistence;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlProofOfExistence;
import eu.europa.esig.dss.detailedreport.jaxb.XmlTLAnalysis;
import eu.europa.esig.dss.detailedreport.jaxb.XmlValidationCertificateQualification;
import eu.europa.esig.dss.detailedreport.jaxb.XmlValidationSignatureQualification;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateExtension;
import eu.europa.esig.dss.diagnostic.jaxb.XmlLangAndValue;
import eu.europa.esig.dss.diagnostic.jaxb.XmlOID;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcCompliance;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcSSCD;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQcStatements;
import eu.europa.esig.dss.diagnostic.jaxb.XmlQualifier;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustService;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustServiceProvider;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrustedList;
import eu.europa.esig.dss.enumerations.CertificateExtensionEnum;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.QCTypeEnum;
import eu.europa.esig.dss.enumerations.ServiceQualification;
import eu.europa.esig.dss.enumerations.SubIndication;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.validation.process.qualification.signature.SignatureQualificationBlock;
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

public class QualSigBlockOracle {

    private static final long PRE_EIDAS = 1370044800000L;
    private static final long POST_EIDAS = 1514764800000L;
    private static final long LATER = 1577836800000L;
    private static final long EARLY = 1420070400000L;

    private static final String TL_URL = "https://tl.example/LU";
    private static final String LOTL_URL = "https://lotl.example/EU";
    private static final String OTHER_TL_URL = "https://tl.example/DE";

    private static final String GRANTED = TrustServiceStatus.GRANTED.getUri();
    private static final String WITHDRAWN = TrustServiceStatus.WITHDRAWN.getUri();
    private static final String CA_QC = ServiceTypeIdentifier.CA_QC.getUri();
    private static final String CA_PKC = ServiceTypeIdentifier.CA_PKC.getUri();
    private static final String ASI_ESIG = "http://uri.etsi.org/TrstSvc/TrustedList/SvcInfoExt/ForeSignatures";

    /** Trust-service-provider shapes attached to the signing certificate. */
    private static final String[] TSP_LABELS = {
            "none",             // no trust service provider at all
            "granted-tl",       // one granted CA/QC service under TL_URL
            "granted-tl-qscd",  // ... with the QCWithQSCD qualifier
            "withdrawn-tl",     // one withdrawn CA/QC service under TL_URL
            "capkc-tl",         // one granted CA/PKC service (not CA/QC)
            "expired-tl",       // granted CA/QC whose validity window ended in 2015
            "granted-lotl",     // granted CA/QC under TL_URL, itself under LOTL_URL
            "granted-other-tl", // granted CA/QC under a TL the analyses do not cover
            "two-conflicting",  // two granted CA/QC services that disagree
    };

    private static List<XmlTrustServiceProvider> tsps(String label) {
        if ("none".equals(label)) {
            return Collections.emptyList();
        }
        XmlTrustServiceProvider tsp = new XmlTrustServiceProvider();
        XmlLangAndValue tspName = new XmlLangAndValue();
        tspName.setLang("en");
        tspName.setValue("TSP " + label);
        tsp.getTSPNames().add(tspName);

        XmlTrustedList tl = new XmlTrustedList();
        tl.setUrl("granted-other-tl".equals(label) ? OTHER_TL_URL : TL_URL);
        tl.setCountryCode("LU");
        tsp.setTL(tl);
        if ("granted-lotl".equals(label)) {
            XmlTrustedList lotl = new XmlTrustedList();
            lotl.setUrl(LOTL_URL);
            lotl.setLOTL(Boolean.TRUE);
            tsp.setLOTL(lotl);
        }

        List<XmlTrustService> services = new ArrayList<>();
        if ("two-conflicting".equals(label)) {
            services.add(trustService("svc-a", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_STATEMENT.getUri())));
            services.add(trustService("svc-b", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.NOT_QUALIFIED.getUri())));
        } else if ("withdrawn-tl".equals(label)) {
            services.add(trustService("svc", WITHDRAWN, CA_QC, POST_EIDAS, null, Collections.<String>emptyList()));
        } else if ("capkc-tl".equals(label)) {
            services.add(trustService("svc", GRANTED, CA_PKC, POST_EIDAS, null, Collections.<String>emptyList()));
        } else if ("expired-tl".equals(label)) {
            services.add(trustService("svc", GRANTED, CA_QC, PRE_EIDAS, EARLY, Collections.<String>emptyList()));
        } else if ("granted-tl-qscd".equals(label)) {
            services.add(trustService("svc", GRANTED, CA_QC, POST_EIDAS, null,
                    Arrays.asList(ServiceQualification.QC_WITH_QSCD.getUri())));
        } else {
            services.add(trustService("svc", GRANTED, CA_QC, POST_EIDAS, null, Collections.<String>emptyList()));
        }
        tsp.setTrustServices(services);
        return Collections.singletonList(tsp);
    }

    private static XmlTrustService trustService(String name, String status, String type,
                                                long start, Long end, List<String> qualifiers) {
        XmlTrustService svc = new XmlTrustService();
        XmlLangAndValue svcName = new XmlLangAndValue();
        svcName.setLang("en");
        svcName.setValue(name);
        svc.getServiceNames().add(svcName);
        svc.setStatus(status);
        svc.setServiceType(type);
        svc.setStartDate(new Date(start));
        if (end != null) {
            svc.setEndDate(new Date(end));
        }
        for (String uri : qualifiers) {
            XmlQualifier q = new XmlQualifier();
            q.setValue(uri);
            svc.getCapturedQualifiers().add(q);
        }
        svc.getAdditionalServiceInfoUris().add(ASI_ESIG);
        XmlCertificate sdi = new XmlCertificate();
        sdi.setId("C-SDI-" + name);
        sdi.setNotBefore(new Date(PRE_EIDAS));
        sdi.setNotAfter(new Date(LATER));
        svc.setServiceDigitalIdentifier(sdi);
        return svc;
    }

    private static final String[] CERT_LABELS = {"plain-post", "qc-esign-post", "qc-esign-qscd-post", "qc-esign-pre"};

    private static XmlCertificate certificate(String label, String tspLabel) {
        boolean post = !label.endsWith("-pre");
        long notBefore = post ? POST_EIDAS : PRE_EIDAS;
        XmlCertificate cert = new XmlCertificate();
        cert.setId("C-" + label);
        cert.setNotBefore(new Date(notBefore));
        cert.setNotAfter(new Date(notBefore + 3L * 365 * 86400000L));
        XmlQcStatements qcStatements = new XmlQcStatements();
        qcStatements.setOID(CertificateExtensionEnum.QC_STATEMENTS.getOid());
        XmlQcCompliance compliance = new XmlQcCompliance();
        compliance.setPresent(label.startsWith("qc-"));
        qcStatements.setQcCompliance(compliance);
        XmlQcSSCD sscd = new XmlQcSSCD();
        sscd.setPresent(label.contains("qscd"));
        qcStatements.setQcSSCD(sscd);
        if (label.contains("esign")) {
            XmlOID oid = new XmlOID();
            oid.setValue(QCTypeEnum.QCT_ESIGN.getOid());
            qcStatements.getQcTypes().add(oid);
        }
        List<XmlCertificateExtension> extensions = new ArrayList<>();
        extensions.add(qcStatements);
        cert.getCertificateExtensions().addAll(extensions);
        cert.setTrustServiceProviders(tsps(tspLabel));
        return cert;
    }

    /** XmlTLAnalysis panels. */
    private static final String[] TLA_LABELS = {"none", "tl-ok", "tl-failed", "tl-warn", "lotl-and-tl-ok", "lotl-failed"};

    private static List<XmlTLAnalysis> tlAnalyses(String label) {
        List<XmlTLAnalysis> result = new ArrayList<>();
        switch (label) {
            case "none":
                break;
            case "tl-ok":
                result.add(tlAnalysis(TL_URL, Indication.PASSED, null, false, false));
                break;
            case "tl-failed":
                result.add(tlAnalysis(TL_URL, Indication.FAILED, null, true, false));
                break;
            case "tl-warn":
                result.add(tlAnalysis(TL_URL, Indication.PASSED, null, false, true));
                break;
            case "lotl-and-tl-ok":
                result.add(tlAnalysis(LOTL_URL, Indication.PASSED, null, false, false));
                result.add(tlAnalysis(TL_URL, Indication.PASSED, null, false, false));
                break;
            case "lotl-failed":
                result.add(tlAnalysis(LOTL_URL, Indication.FAILED, null, true, false));
                result.add(tlAnalysis(TL_URL, Indication.PASSED, null, false, false));
                break;
            default:
                throw new IllegalStateException(label);
        }
        return result;
    }

    private static XmlTLAnalysis tlAnalysis(String url, Indication indication, SubIndication subIndication,
                                            boolean withError, boolean withWarning) {
        XmlTLAnalysis analysis = new XmlTLAnalysis();
        analysis.setURL(url);
        analysis.setCountryCode("LU");
        analysis.setTitle("TL " + url);
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(indication);
        conclusion.setSubIndication(subIndication);
        if (withError) {
            XmlMessage m = new XmlMessage();
            m.setKey("QUAL_TL_EXP_ANS");
            m.setValue("the trusted list has expired");
            conclusion.getErrors().add(m);
        }
        if (withWarning) {
            XmlMessage m = new XmlMessage();
            m.setKey("QUAL_TL_FRESH_ANS");
            m.setValue("the trusted list is not fresh");
            conclusion.getWarnings().add(m);
        }
        analysis.setConclusion(conclusion);
        return analysis;
    }

    /** EN 319 102-1 conclusions the block is fed. */
    private static final String[] ETSI_LABELS = {"passed", "indeterminate", "failed"};

    private static XmlConstraintsConclusionWithProofOfExistence etsiResult(String label, long bestSignatureTime) {
        XmlConstraintsConclusionWithProofOfExistence result = new XmlConstraintsConclusionWithProofOfExistence();
        XmlConclusion conclusion = new XmlConclusion();
        if ("passed".equals(label)) {
            conclusion.setIndication(Indication.PASSED);
        } else if ("indeterminate".equals(label)) {
            conclusion.setIndication(Indication.INDETERMINATE);
            conclusion.setSubIndication(SubIndication.TRY_LATER);
        } else {
            conclusion.setIndication(Indication.TOTAL_FAILED);
            conclusion.setSubIndication(SubIndication.HASH_FAILURE);
        }
        result.setConclusion(conclusion);
        XmlProofOfExistence poe = new XmlProofOfExistence();
        poe.setTime(new Date(bestSignatureTime));
        result.setProofOfExistence(poe);
        return result;
    }

    private static final long[] BEST_SIGNATURE_TIMES = {EARLY, POST_EIDAS, LATER};

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
            System.err.println("usage: QualSigBlockOracle <dss-repo-root>");
            System.exit(2);
        }
        Path out = Paths.get(args[0], "validation", "process", "qualification", "testdata", "oracle", "qual_sig_block.jsonl");
        Files.createDirectories(out.getParent());
        I18nProvider i18nProvider = new I18nProvider();
        int index = 0;
        try (PrintWriter pw = new PrintWriter(Files.newBufferedWriter(out, StandardCharsets.UTF_8))) {
            for (String certLabel : CERT_LABELS) {
                for (String tspLabel : TSP_LABELS) {
                    for (String tlaLabel : TLA_LABELS) {
                        for (String etsiLabel : ETSI_LABELS) {
                            for (long bst : BEST_SIGNATURE_TIMES) {
                                CertificateWrapper cert = new CertificateWrapper(certificate(certLabel, tspLabel));
                                List<XmlTLAnalysis> analyses = tlAnalyses(tlaLabel);
                                SignatureQualificationBlock block = new SignatureQualificationBlock(
                                        i18nProvider, etsiResult(etsiLabel, bst), cert, analyses);
                                String outcome;
                                try {
                                    XmlValidationSignatureQualification result = block.execute();
                                    StringBuilder nested = new StringBuilder("[");
                                    List<XmlValidationCertificateQualification> certQuals =
                                            result.getValidationCertificateQualification();
                                    for (int i = 0; i < certQuals.size(); i++) {
                                        XmlValidationCertificateQualification q = certQuals.get(i);
                                        if (i > 0) {
                                            nested.append(',');
                                        }
                                        nested.append("{\"validationTime\":")
                                                .append(esc(q.getValidationTime() == null ? null : q.getValidationTime().name()))
                                                .append(",\"certificateQualification\":")
                                                .append(esc(q.getCertificateQualification() == null ? null
                                                        : q.getCertificateQualification().name()))
                                                .append(",\"constraints\":").append(constraints(q.getConstraint()))
                                                .append(",\"conclusion\":").append(conclusion(q.getConclusion()))
                                                .append('}');
                                    }
                                    nested.append(']');
                                    outcome = "{\"title\":" + esc(result.getTitle())
                                            + ",\"signatureQualification\":"
                                            + esc(result.getSignatureQualification() == null ? null
                                                    : result.getSignatureQualification().name())
                                            + ",\"constraints\":" + constraints(result.getConstraint())
                                            + ",\"conclusion\":" + conclusion(result.getConclusion())
                                            + ",\"certificateQualifications\":" + nested + "}";
                                } catch (RuntimeException e) {
                                    outcome = esc("throws:" + e.getClass().getSimpleName());
                                }
                                pw.println("{\"kind\":\"sigQual\",\"i\":" + index++
                                        + ",\"cert\":" + esc(certLabel)
                                        + ",\"tsp\":" + esc(tspLabel)
                                        + ",\"tlAnalyses\":" + esc(tlaLabel)
                                        + ",\"etsi\":" + esc(etsiLabel)
                                        + ",\"bestSignatureTime\":" + bst
                                        + ",\"result\":" + outcome + "}");
                            }
                        }
                    }
                }
            }
        }
        System.out.println("wrote " + out + " (" + index + " rows)");
    }
}

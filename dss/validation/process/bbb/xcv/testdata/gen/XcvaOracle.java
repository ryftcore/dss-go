/*
 * Java oracle driver for the three XCVA building blocks of phase 8d:
 *
 *   eu.europa.esig.dss.validation.process.bbb.xcv.X509CertificateValidation
 *   eu.europa.esig.dss.validation.process.bbb.xcv.crs.CertificateRevocationSelector
 *   eu.europa.esig.dss.validation.process.bbb.xcv.rac.RevocationAcceptanceChecker
 *
 * It runs the UPSTREAM classes over REAL inputs - every XmlDiagnosticData dump of
 * the marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle - with the
 * default ETSI validation policy and the fixed validation time 2024-01-01T00:00:00Z,
 * and dumps one JSON object per line.
 *
 * X509CertificateValidation is driven exactly the way BasicBuildingBlocks drives it:
 * over the signing certificate of every signature / time-stamp / revocation token
 * and over every used certificate, with the usage time that dispatcher passes
 * (the certificate's notBefore, the time-stamp's production time, the revocation's
 * production date). CertificateRevocationSelector is driven over every used
 * certificate; RevocationAcceptanceChecker over every (certificate, certificate
 * revocation data) pair, each with its own fresh validated-token set.
 *
 * The rows carry the nested results too - a CRS row its RAC children, a RAC row its
 * CRS child, an XCV row its SubXCV children - so the Go replay compares the whole
 * tree the block produced and not only its conclusion.
 *
 * Rows are written to
 *   dss/validation/process/bbb/xcv/testdata/oracle/xcva_blocks.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle XcvaOracle.java
 *   java  -cp "$CP:/tmp/oracle" XcvaOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlAOV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCRS;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlRAC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSubXCV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlXCV;
import eu.europa.esig.dss.diagnostic.CertificateRevocationWrapper;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;

import java.io.FileInputStream;
import java.io.InputStream;
import eu.europa.esig.dss.validation.process.bbb.xcv.X509CertificateValidation;
import eu.europa.esig.dss.validation.process.bbb.xcv.crs.CertificateRevocationSelector;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.RevocationAcceptanceChecker;

import java.io.File;
import java.io.PrintWriter;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.HashSet;
import java.util.List;
import java.util.TimeZone;

public class XcvaOracle {

    /** Fixed validation time: 2024-01-01T00:00:00Z. */
    private static final Date CURRENT_TIME = new Date(1704067200000L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();
        ValidationPolicy defaultPolicy = new EtsiValidationPolicyFactory().loadDefaultValidationPolicy();

        // A second policy, committed next to the dumps, that switches the validation
        // model to CHAIN and drops the prospective-certificate-chain constraint to
        // WARN. Under the default policy that check is FAIL, so an untrusted chain
        // stops the block on its first item and neither the SubXCV loop nor the
        // per-model lastDate progression is ever reached.
        ValidationPolicy chainPolicy;
        try (InputStream is = new FileInputStream(new File(repoRoot,
                "validation/process/bbb/xcv/testdata/policy/constraint-chain.xml"))) {
            chainPolicy = new EtsiValidationPolicyFactory().loadValidationPolicy(is);
        }

        // The synthetic dumps of XcvaSyntheticDumps, which carry the trust anchors
        // the marshal-parity corpus has none of.
        File[] synthetic = new File(repoRoot, "validation/process/bbb/xcv/testdata/dd")
                .listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(synthetic, Comparator.comparing(File::getName));

        try (PrintWriter out = writer(repoRoot,
                "validation/process/bbb/xcv/testdata/oracle/xcva_blocks.jsonl")) {

            for (File synth : synthetic) {
                DiagnosticData dd = new DiagnosticData(
                        DiagnosticDataFacade.newFacade().unmarshall(synth, false));
                emit(out, i18n, defaultPolicy, "dd/" + synth.getName(), "default", dd);
                emit(out, i18n, chainPolicy, "dd/" + synth.getName(), "chain", dd);
            }

            for (File file : files) {
                String name = file.getName();
                // See ../../../fc/testdata/README.md for why the four model-*.xml
                // schema-coverage fixtures are excluded.
                if (name.startsWith("model-")) {
                    System.err.println("SKIPPED-FIXTURE " + name);
                    continue;
                }
                DiagnosticData dd;
                try {
                    XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                    dd = new DiagnosticData(jaxb);
                    dd.getSignatures();
                    dd.getTimestampList();
                    dd.getUsedCertificates();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED-FILE " + name + ": " + e);
                    continue;
                }

                emit(out, i18n, defaultPolicy, name, "default", dd);
            }
        }
    }

    /** Drives every XCVA block of one dump under one policy. */
    private static void emit(PrintWriter out, I18nProvider i18n, ValidationPolicy policy, String file,
                             String policyName, DiagnosticData dd) {
        for (SignatureWrapper signature : dd.getSignatures()) {
            Context context = signature.isCounterSignature()
                    ? Context.COUNTER_SIGNATURE : Context.SIGNATURE;
            CertificateWrapper signingCertificate = signature.getSigningCertificate();
            if (signingCertificate != null) {
                xcvRow(out, i18n, policy, file, policyName, "SIG|" + signature.getId(), context,
                        signingCertificate, signingCertificate.getNotBefore());
            }
        }
        for (TimestampWrapper timestamp : dd.getTimestampList()) {
            CertificateWrapper signingCertificate = timestamp.getSigningCertificate();
            if (signingCertificate != null) {
                xcvRow(out, i18n, policy, file, policyName, "TST|" + timestamp.getId(), Context.TIMESTAMP,
                        signingCertificate, timestamp.getProductionTime());
            }
        }
        for (RevocationWrapper revocation : sortedById(dd.getAllRevocationData())) {
            CertificateWrapper signingCertificate = revocation.getSigningCertificate();
            if (signingCertificate != null) {
                xcvRow(out, i18n, policy, file, policyName, "REV|" + revocation.getId(), Context.REVOCATION,
                        signingCertificate, revocation.getProductionDate());
            }
        }
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            xcvRow(out, i18n, policy, file, policyName, "CERT|" + certificate.getId(), Context.CERTIFICATE,
                    certificate, certificate.getNotBefore());
            crsRow(out, i18n, policy, file, policyName, certificate);
            for (CertificateRevocationWrapper revocation : certificate.getCertificateRevocationData()) {
                racRow(out, i18n, policy, file, policyName, certificate, revocation);
            }
        }
    }

    private static List<RevocationWrapper> sortedById(java.util.Set<RevocationWrapper> revocations) {
        List<RevocationWrapper> list = new java.util.ArrayList<>(revocations);
        list.sort(Comparator.comparing(RevocationWrapper::getId));
        return list;
    }

    /** The PASSED XmlAOV the blocks are fed, the way the sav corpus does. */
    private static XmlAOV passedAOV() {
        XmlAOV aov = new XmlAOV();
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(Indication.PASSED);
        aov.setConclusion(conclusion);
        return aov;
    }

    private static void xcvRow(PrintWriter out, I18nProvider i18n, ValidationPolicy policy, String file,
                               String policyName, String token, Context context,
                               CertificateWrapper certificate, Date usageTime) {
        try {
            XmlXCV result = new X509CertificateValidation(i18n, certificate, CURRENT_TIME, usageTime,
                    context, passedAOV(), policy).execute();
            StringBuilder row = new StringBuilder("{");
            head(row, file, policyName, token, context, "XCV");
            body(row, result);
            row.append(',');
            key(row, "subXCV").append('[');
            List<XmlSubXCV> subs = result.getSubXCV();
            for (int i = 0; i < subs.size(); i++) {
                if (i > 0) {
                    row.append(',');
                }
                row.append(subXCV(subs.get(i)));
            }
            row.append(']');
            out.println(row.append('}').toString());
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " XCV: " + e);
        }
    }

    private static void crsRow(PrintWriter out, I18nProvider i18n, ValidationPolicy policy, String file,
                               String policyName, CertificateWrapper certificate) {
        try {
            XmlCRS result = new CertificateRevocationSelector(i18n, certificate, CURRENT_TIME, policy).execute();
            StringBuilder row = new StringBuilder("{");
            head(row, file, policyName, "CRS|" + certificate.getId(), Context.REVOCATION, "CRS");
            row.append(crsBody(result).substring(1));
            out.println(row.toString());
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " CRS|" + certificate.getId() + ": " + e);
        }
    }

    private static void racRow(PrintWriter out, I18nProvider i18n, ValidationPolicy policy, String file,
                               String policyName, CertificateWrapper certificate,
                               CertificateRevocationWrapper revocation) {
        try {
            XmlRAC result = new RevocationAcceptanceChecker(i18n, certificate, revocation, CURRENT_TIME,
                    policy, new HashSet<>()).execute();
            StringBuilder row = new StringBuilder("{");
            head(row, file, policyName, "RAC|" + certificate.getId() + "|" + revocation.getId(), Context.REVOCATION, "RAC");
            row.append(racBody(result).substring(1));
            out.println(row.toString());
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " RAC|" + certificate.getId() + "|" + revocation.getId()
                    + ": " + e);
        }
    }

    // --------------------------------------------------------------------- JSON

    private static void head(StringBuilder out, String file, String policyName, String token,
                             Context context, String block) {
        key(out, "file").append(str(file)).append(',');
        key(out, "policy").append(str(policyName)).append(',');
        key(out, "token").append(str(token)).append(',');
        key(out, "context").append(str(context.name())).append(',');
        key(out, "block").append(str(block)).append(',');
    }

    /** title + conclusion + constraints, without the enclosing braces. */
    private static void body(StringBuilder out, XmlConstraintsConclusion result) {
        key(out, "title").append(str(result.getTitle())).append(',');
        key(out, "conclusion").append(conclusion(result.getConclusion())).append(',');
        key(out, "constraints").append('[');
        List<XmlConstraint> constraints = result.getConstraint();
        for (int i = 0; i < constraints.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(constraint(constraints.get(i)));
        }
        out.append(']');
    }

    private static String crsBody(XmlCRS crs) {
        StringBuilder out = new StringBuilder("{");
        key(out, "id").append(str(crs.getId())).append(',');
        key(out, "latestAcceptableRevocationId").append(str(crs.getLatestAcceptableRevocationId())).append(',');
        key(out, "acceptableRevocationId").append(strings(crs.getAcceptableRevocationId())).append(',');
        body(out, crs);
        out.append(',');
        key(out, "rac").append('[');
        List<XmlRAC> racs = crs.getRAC();
        for (int i = 0; i < racs.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(racBody(racs.get(i)));
        }
        out.append(']');
        return out.append('}').toString();
    }

    private static String racBody(XmlRAC rac) {
        StringBuilder out = new StringBuilder("{");
        key(out, "id").append(str(rac.getId())).append(',');
        key(out, "revocationThisUpdate").append(str(iso(rac.getRevocationThisUpdate()))).append(',');
        key(out, "revocationProductionDate").append(str(iso(rac.getRevocationProductionDate()))).append(',');
        body(out, rac);
        out.append(',');
        key(out, "crs").append(rac.getCRS() == null ? "null" : crsBody(rac.getCRS()));
        return out.append('}').toString();
    }

    private static String subXCV(XmlSubXCV sub) {
        StringBuilder out = new StringBuilder("{");
        key(out, "id").append(str(sub.getId())).append(',');
        key(out, "trustAnchor").append(bool(sub.isTrustAnchor())).append(',');
        key(out, "selfSigned").append(bool(sub.isSelfSigned())).append(',');
        body(out, sub);
        return out.append('}').toString();
    }

    private static String bool(Boolean value) {
        return value == null ? "null" : value.toString();
    }

    private static String strings(List<String> values) {
        if (values == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("[");
        for (int i = 0; i < values.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(str(values.get(i)));
        }
        return out.append(']').toString();
    }

    private static String iso(Date date) {
        if (date == null) {
            return null;
        }
        SimpleDateFormat format = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");
        format.setTimeZone(TimeZone.getTimeZone("UTC"));
        return format.format(date);
    }

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static StringBuilder key(StringBuilder out, String name) {
        return out.append(str(name)).append(':');
    }

    private static String conclusion(XmlConclusion conclusion) {
        if (conclusion == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "indication").append(str(conclusion.getIndication() == null ? null : conclusion.getIndication().name())).append(',');
        key(out, "subIndication").append(str(conclusion.getSubIndication() == null ? null : conclusion.getSubIndication().name())).append(',');
        key(out, "errors").append(messages(conclusion.getErrors())).append(',');
        key(out, "warnings").append(messages(conclusion.getWarnings())).append(',');
        key(out, "infos").append(messages(conclusion.getInfos()));
        return out.append('}').toString();
    }

    private static String messages(List<XmlMessage> messages) {
        StringBuilder out = new StringBuilder("[");
        for (int i = 0; i < messages.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(message(messages.get(i)));
        }
        return out.append(']').toString();
    }

    private static String message(XmlMessage message) {
        if (message == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "key").append(str(message.getKey())).append(',');
        key(out, "value").append(str(message.getValue()));
        return out.append('}').toString();
    }

    private static String constraint(XmlConstraint constraint) {
        StringBuilder out = new StringBuilder("{");
        key(out, "name").append(message(constraint.getName())).append(',');
        key(out, "status").append(str(constraint.getStatus() == null ? null : constraint.getStatus().value())).append(',');
        key(out, "error").append(message(constraint.getError())).append(',');
        key(out, "warning").append(message(constraint.getWarning())).append(',');
        key(out, "info").append(message(constraint.getInfo())).append(',');
        key(out, "additionalInfo").append(str(constraint.getAdditionalInfo())).append(',');
        key(out, "id").append(str(constraint.getId())).append(',');
        key(out, "blockType").append(str(constraint.getBlockType() == null ? null : constraint.getBlockType().value()));
        return out.append('}').toString();
    }

    private static String str(String value) {
        if (value == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("\"");
        for (int i = 0; i < value.length(); i++) {
            char c = value.charAt(i);
            switch (c) {
                case '"': out.append("\\\""); break;
                case '\\': out.append("\\\\"); break;
                case '\n': out.append("\\n"); break;
                case '\r': out.append("\\r"); break;
                case '\t': out.append("\\t"); break;
                default:
                    if (c < 0x20) {
                        out.append(String.format("\\u%04x", (int) c));
                    } else {
                        out.append(c);
                    }
            }
        }
        return out.append('"').toString();
    }

}

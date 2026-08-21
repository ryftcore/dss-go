/*
 * Java oracle driver for the twelve
 * eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks classes (the XCVA half
 * of phase 8d).
 *
 * Every check is driven directly through a chain of exactly one item at
 * Level.FAIL, over REAL inputs: each dump's used certificates crossed with each of
 * their certificate revocation data. The inputs are the marshal-parity corpus in
 * dss/diagnostic/jaxb/testdata/oracle plus the committed synthetic dumps of
 * XcvaSyntheticDumps in ../dd, which carry the revocation shapes that corpus has
 * none of (a certHash extension, a resolving responder-id reference, an absent
 * thisUpdate, an expiredCertsOnCRL, an archiveCutoff). Driving every check over
 * every pair is what produces both the happy and the failure row; the Go replay
 * reads the very same files, so neither side gets a private fixture.
 *
 * RevocationAcceptanceCheckerResultCheck is driven over the XmlRAC that the
 * upstream RevocationAcceptanceChecker produces for the same pair, so its input is
 * a real result too.
 *
 * Rows are written to
 *   dss/validation/process/bbb/xcv/testdata/oracle/xcva_direct.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle XcvaDirectOracle.java
 *   java  -cp "$CP:/tmp/oracle" XcvaDirectOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlRAC;
import eu.europa.esig.dss.diagnostic.CertificateRevocationWrapper;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.RevocationAcceptanceChecker;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationAcceptanceCheckerResultCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationAfterCertificateIssuanceCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationCertHashMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationCertHashPresenceCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationDataKnownCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationHasInformationAboutCertificateCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerKnownCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerRevocationDataAvailableCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerValidAtProductionTimeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationResponderIdMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.SelfIssuedOCSPCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.ThisUpdatePresenceCheck;

import java.io.File;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.HashSet;
import java.util.List;
import java.util.function.BiFunction;

public class XcvaDirectOracle {

    private static final LevelRule FAIL = () -> Level.FAIL;

    /** Fixed validation time: 2024-01-01T00:00:00Z, the one XcvaOracle uses. */
    private static final Date CURRENT_TIME = new Date(1704067200000L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));
        File[] synthetic = new File(repoRoot, "validation/process/bbb/xcv/testdata/dd")
                .listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(synthetic, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();
        ValidationPolicy policy = new EtsiValidationPolicyFactory().loadDefaultValidationPolicy();

        try (PrintWriter out = writer(repoRoot,
                "validation/process/bbb/xcv/testdata/oracle/xcva_direct.jsonl")) {

            for (File synth : synthetic) {
                DiagnosticData dd = new DiagnosticData(
                        DiagnosticDataFacade.newFacade().unmarshall(synth, false));
                emitDump(out, i18n, policy, "dd/" + synth.getName(), dd);
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
                    dd.getUsedCertificates();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED-FILE " + name + ": " + e);
                    continue;
                }
                emitDump(out, i18n, policy, name, dd);
            }
        }
    }

    private static void emitDump(PrintWriter out, I18nProvider i18n, ValidationPolicy policy,
                                 String file, DiagnosticData dd) {
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            row(out, i18n, file, certificate.getId(), "RevocationIssuerRevocationDataAvailableCheck",
                    (r, l) -> new RevocationIssuerRevocationDataAvailableCheck(i18n, r, certificate, l));

            for (CertificateRevocationWrapper revocation : certificate.getCertificateRevocationData()) {
                String token = certificate.getId() + "|" + revocation.getId();
                emitAll(out, i18n, file, token, certificate, revocation);

                XmlRAC racResult;
                try {
                    racResult = new RevocationAcceptanceChecker(i18n, certificate, revocation,
                            CURRENT_TIME, policy, new HashSet<>()).execute();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED " + file + " " + token + " RAC: " + e);
                    continue;
                }
                row(out, i18n, file, token, "RevocationAcceptanceCheckerResultCheck",
                        (r, l) -> new RevocationAcceptanceCheckerResultCheck<>(i18n, r, racResult, l));
            }
        }
    }

    /** Drives every rac check that reads a (certificate, revocation) pair. */
    private static void emitAll(PrintWriter out, I18nProvider i18n, String file, String token,
                                CertificateWrapper certificate, CertificateRevocationWrapper revocation) {
        row(out, i18n, file, token, "RevocationDataKnownCheck",
                (r, l) -> new RevocationDataKnownCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationIssuerKnownCheck",
                (r, l) -> new RevocationIssuerKnownCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "ThisUpdatePresenceCheck",
                (r, l) -> new ThisUpdatePresenceCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationIssuerValidAtProductionTimeCheck",
                (r, l) -> new RevocationIssuerValidAtProductionTimeCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationResponderIdMatchCheck",
                (r, l) -> new RevocationResponderIdMatchCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationCertHashPresenceCheck",
                (r, l) -> new RevocationCertHashPresenceCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationCertHashMatchCheck",
                (r, l) -> new RevocationCertHashMatchCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "SelfIssuedOCSPCheck",
                (r, l) -> new SelfIssuedOCSPCheck(i18n, r, certificate, revocation, l));
        row(out, i18n, file, token, "RevocationAfterCertificateIssuanceCheck",
                (r, l) -> new RevocationAfterCertificateIssuanceCheck(i18n, r, certificate, revocation, l));
        row(out, i18n, file, token, "RevocationHasInformationAboutCertificateCheck",
                (r, l) -> new RevocationHasInformationAboutCertificateCheck(i18n, r, certificate, revocation, l));
    }

    // ------------------------------------------------------------------ harness

    /** A chain of exactly one item over an XmlRAC. */
    static class SingleRACChain extends Chain<XmlRAC> {
        private final BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory;

        SingleRACChain(I18nProvider i18nProvider, BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory) {
            super(i18nProvider, new XmlRAC());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    private static void row(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                            BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory) {
        try {
            XmlRAC result = new SingleRACChain(i18n, factory).execute();
            out.println(json(file, token, check, result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    // --------------------------------------------------------------------- JSON

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static String json(String file, String token, String check, XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(token)).append(',');
        key(out, "check").append(str(check)).append(',');
        key(out, "context").append(str(Context.REVOCATION.name())).append(',');
        key(out, "block").append(str("RAC")).append(',');
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
        return out.append('}').toString();
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

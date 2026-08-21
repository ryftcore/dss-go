/*
 * Java oracle driver for the fc (format checking) and sav (signature acceptance
 * validation) building blocks.
 *
 * It runs the UPSTREAM building blocks
 *
 *   eu.europa.esig.dss.validation.process.bbb.fc.SignatureFormatChecking
 *   eu.europa.esig.dss.validation.process.bbb.fc.TimestampFormatChecking
 *   eu.europa.esig.dss.validation.process.bbb.sav.SignatureAcceptanceValidation
 *   eu.europa.esig.dss.validation.process.bbb.sav.TimestampAcceptanceValidation
 *   eu.europa.esig.dss.validation.process.bbb.sav.RevocationAcceptanceValidation
 *
 * over REAL inputs - every XmlDiagnosticData dump of the marshal-parity corpus in
 * dss/diagnostic/jaxb/testdata/oracle - with the default ETSI validation policy,
 * and dumps each produced XmlFC / XmlSAV as one JSON object per line, in the same
 * row shape BbbBlocksOracle uses.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle FcSavOracle.java
 *   java  -cp "$CP:/tmp/oracle" FcSavOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlBasicBuildingBlocks;
import eu.europa.esig.dss.detailedreport.jaxb.XmlAOV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlFC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSAV;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;
import eu.europa.esig.dss.validation.process.bbb.fc.SignatureFormatChecking;
import eu.europa.esig.dss.validation.process.bbb.fc.TimestampFormatChecking;
import eu.europa.esig.dss.validation.process.bbb.sav.RevocationAcceptanceValidation;
import eu.europa.esig.dss.validation.process.bbb.sav.SignatureAcceptanceValidation;
import eu.europa.esig.dss.validation.process.bbb.sav.TimestampAcceptanceValidation;

import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.List;

public class FcSavOracle {

    /** Fixed validation time: 2024-01-01T00:00:00Z. */
    private static final Date CURRENT_TIME = new Date(1704067200000L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        I18nProvider i18nProvider = new I18nProvider();
        ValidationPolicy policy = new EtsiValidationPolicyFactory().loadDefaultValidationPolicy();

        try (PrintWriter fc = writer(repoRoot, "validation/process/bbb/fc/testdata/oracle/fc_blocks.jsonl");
             PrintWriter sav = writer(repoRoot, "validation/process/bbb/sav/testdata/oracle/sav_blocks.jsonl")) {

            for (File file : files) {
                XmlDiagnosticData jaxb;
                DiagnosticData diagnosticData;
                String name = file.getName();
                try {
                    jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                    diagnosticData = new DiagnosticData(jaxb);
                    // touch the wrapper graph the same way BbbBlocksOracle does, so the
                    // four schema-coverage fixtures are skipped wholesale, as there
                    diagnosticData.getSignatures();
                    diagnosticData.getTimestampList();
                    diagnosticData.getAllRevocationData();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED-FILE " + name + ": " + e);
                    continue;
                }

                for (SignatureWrapper signature : diagnosticData.getSignatures()) {
                    Context context = signature.isCounterSignature() ? Context.COUNTER_SIGNATURE : Context.SIGNATURE;
                    try {
                        fc.println(row(name, signature.getId(), context, "FC",
                                new SignatureFormatChecking(i18nProvider, diagnosticData, signature, context, policy).execute()));
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED-FC " + name + " " + signature.getId() + ": " + e);
                    }
                    try {
                        sav.println(row(name, signature.getId(), context, "SAV",
                                new SignatureAcceptanceValidation(i18nProvider, diagnosticData, CURRENT_TIME, signature,
                                        context, Collections.<String, XmlBasicBuildingBlocks>emptyMap(), passedAov(), policy).execute()));
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED-SAV " + name + " " + signature.getId() + ": " + e);
                    }
                }

                for (TimestampWrapper timestamp : diagnosticData.getTimestampList()) {
                    try {
                        fc.println(row(name, timestamp.getId(), Context.TIMESTAMP, "FC",
                                new TimestampFormatChecking(i18nProvider, diagnosticData, timestamp, Context.TIMESTAMP, policy).execute()));
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED-FC " + name + " " + timestamp.getId() + ": " + e);
                    }
                    try {
                        sav.println(row(name, timestamp.getId(), Context.TIMESTAMP, "SAV",
                                new TimestampAcceptanceValidation(i18nProvider, CURRENT_TIME, timestamp, passedAov(), policy).execute()));
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED-SAV " + name + " " + timestamp.getId() + ": " + e);
                    }
                }

                List<RevocationWrapper> revocations = new ArrayList<>(diagnosticData.getAllRevocationData());
                revocations.sort(Comparator.comparing(RevocationWrapper::getId));
                for (RevocationWrapper revocation : revocations) {
                    try {
                        sav.println(row(name, revocation.getId(), Context.REVOCATION, "SAV",
                                new RevocationAcceptanceValidation(i18nProvider, CURRENT_TIME, revocation, passedAov(), policy).execute()));
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED-SAV " + name + " " + revocation.getId() + ": " + e);
                    }
                }
            }
        }
    }

    /**
     * A PASSED Algorithm Obsolescence Validation result with no errors, warnings or
     * infos: AlgorithmObsolescenceValidationCheck reads its conclusion to pick its
     * own Level (FAIL here) and to decide process(). AOV production itself belongs
     * to phase 8d; the sav chains only consume the finished block.
     */
    private static XmlAOV passedAov() {
        XmlAOV aov = new XmlAOV();
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
        aov.setConclusion(conclusion);
        return aov;
    }

    // ------------------------------------------------------------------- JSON

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static String row(String file, String tokenId, Context context, String block,
                              XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(tokenId)).append(',');
        key(out, "context").append(str(context.name())).append(',');
        key(out, "block").append(str(block)).append(',');
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

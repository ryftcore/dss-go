/*
 * Java oracle driver for the Chain / ChainItem / UninterruptedChainItem port.
 *
 * It drives the UPSTREAM classes (eu.europa.esig.dss.validation.process.Chain,
 * ChainItem, UninterruptedChainItem) over a table of synthetic chain
 * specifications and dumps, for each specification, the resulting
 * XmlConstraintsConclusion (title, conclusion, constraint list) as one JSON
 * object per line. The Go table test in ../../chain_test.go replays the same
 * table against the Go port and compares against these rows.
 *
 * The scenarios cover, per the check-execution algorithm: every Level branch
 * (FAIL / WARN / INFORM / IGNORE / undefined constraint), the first-failure
 * short circuit, continueProcessOnFail (UninterruptedChainItem), the
 * conclusion-reuse of an uninterrupted chain, custom success conclusions,
 * previous-error propagation, sub-indication propagation, warning/info
 * collection order, and the population of every XmlConstraint member
 * (Name/Status/Error/Warning/Info/AdditionalInfo/Id/BlockType).
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"      # mvn dependency:build-classpath -pl dss-validation
 *   javac -cp "$CP" -d /tmp/oracle ChainSemanticsOracle.java
 *   java  -cp "$CP:/tmp/oracle" ChainSemanticsOracle > ../oracle/chain_semantics.jsonl
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlBlockType;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.enumerations.SubIndication;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.i18n.MessageTag;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.UninterruptedChainItem;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

public class ChainSemanticsOracle {

    /** Specification of a single chain item. */
    static class Spec {
        Level level;                       // null = constraint not defined
        boolean valid;
        boolean uninterrupted;
        String bbbId;
        MessageTag messageTag = MessageTag.BBB_ICS_ISCI;
        MessageTag errorMessageTag = MessageTag.BBB_ICS_ISCI_ANS;
        MessageTag additionalInfo;
        XmlBlockType blockType;
        Indication failedIndication = Indication.INDETERMINATE;
        SubIndication failedSubIndication = SubIndication.NO_SIGNING_CERTIFICATE_FOUND;
        Indication successIndication;
        SubIndication successSubIndication;
        List<MessageTag> previousErrors = Collections.emptyList();

        Spec(Level level, boolean valid) {
            this.level = level;
            this.valid = valid;
        }

        Spec tags(MessageTag messageTag, MessageTag errorMessageTag) {
            this.messageTag = messageTag;
            this.errorMessageTag = errorMessageTag;
            return this;
        }

        Spec additionalInfo(MessageTag tag) {
            this.additionalInfo = tag;
            return this;
        }

        Spec blockType(XmlBlockType blockType) {
            this.blockType = blockType;
            return this;
        }

        Spec bbbId(String bbbId) {
            this.bbbId = bbbId;
            return this;
        }

        Spec failure(Indication indication, SubIndication subIndication) {
            this.failedIndication = indication;
            this.failedSubIndication = subIndication;
            return this;
        }

        Spec success(Indication indication, SubIndication subIndication) {
            this.successIndication = indication;
            this.successSubIndication = subIndication;
            return this;
        }

        Spec uninterrupted() {
            this.uninterrupted = true;
            return this;
        }

        Spec previousErrors(MessageTag... tags) {
            this.previousErrors = Arrays.asList(tags);
            return this;
        }

        LevelRule rule() {
            if (level == null) {
                return null;
            }
            return () -> level;
        }
    }

    static class SpecItem extends ChainItem<XmlConstraintsConclusion> {
        private final Spec spec;

        SpecItem(I18nProvider i18nProvider, XmlConstraintsConclusion result, Spec spec) {
            super(i18nProvider, result, spec.rule(), spec.bbbId);
            this.spec = spec;
        }

        @Override protected boolean process() { return spec.valid; }
        @Override protected MessageTag getMessageTag() { return spec.messageTag; }
        @Override protected MessageTag getErrorMessageTag() { return spec.errorMessageTag; }
        @Override protected MessageTag getAdditionalInfo() { return spec.additionalInfo; }
        @Override protected XmlBlockType getBlockType() { return spec.blockType; }
        @Override protected Indication getFailedIndicationForConclusion() { return spec.failedIndication; }
        @Override protected SubIndication getFailedSubIndicationForConclusion() { return spec.failedSubIndication; }
        @Override protected Indication getSuccessIndication() { return spec.successIndication; }
        @Override protected SubIndication getSuccessSubIndication() { return spec.successSubIndication; }

        @Override protected List<XmlMessage> getPreviousErrors() {
            List<XmlMessage> messages = new ArrayList<>();
            for (MessageTag tag : spec.previousErrors) {
                messages.add(buildXmlMessage(tag));
            }
            return messages;
        }
    }

    static class UninterruptedSpecItem extends UninterruptedChainItem<XmlConstraintsConclusion> {
        private final Spec spec;

        UninterruptedSpecItem(I18nProvider i18nProvider, XmlConstraintsConclusion result, Spec spec) {
            super(i18nProvider, result, spec.rule(), spec.bbbId);
            this.spec = spec;
        }

        @Override protected boolean process() { return spec.valid; }
        @Override protected MessageTag getMessageTag() { return spec.messageTag; }
        @Override protected MessageTag getErrorMessageTag() { return spec.errorMessageTag; }
        @Override protected MessageTag getAdditionalInfo() { return spec.additionalInfo; }
        @Override protected XmlBlockType getBlockType() { return spec.blockType; }
        @Override protected Indication getFailedIndicationForConclusion() { return spec.failedIndication; }
        @Override protected SubIndication getFailedSubIndicationForConclusion() { return spec.failedSubIndication; }
        @Override protected Indication getSuccessIndication() { return spec.successIndication; }
        @Override protected SubIndication getSuccessSubIndication() { return spec.successSubIndication; }

        @Override protected List<XmlMessage> getPreviousErrors() {
            List<XmlMessage> messages = new ArrayList<>();
            for (MessageTag tag : spec.previousErrors) {
                messages.add(buildXmlMessage(tag));
            }
            return messages;
        }
    }

    static class SpecChain extends Chain<XmlConstraintsConclusion> {
        private final List<Spec> specs;
        private final MessageTag title;

        SpecChain(I18nProvider i18nProvider, MessageTag title, List<Spec> specs) {
            super(i18nProvider, new XmlConstraintsConclusion());
            this.title = title;
            this.specs = specs;
        }

        @Override protected MessageTag getTitle() { return title; }

        @Override protected void initChain() {
            ChainItem<XmlConstraintsConclusion> item = null;
            for (Spec spec : specs) {
                ChainItem<XmlConstraintsConclusion> next = spec.uninterrupted
                        ? new UninterruptedSpecItem(i18nProvider, result, spec)
                        : new SpecItem(i18nProvider, result, spec);
                if (item == null) {
                    item = firstItem = next;
                } else {
                    item = item.setNextItem(next);
                }
            }
        }
    }

    // --------------------------------------------------------------- scenarios

    static class Scenario {
        final String name;
        final MessageTag title;
        final List<Spec> specs;

        Scenario(String name, MessageTag title, Spec... specs) {
            this.name = name;
            this.title = title;
            this.specs = Arrays.asList(specs);
        }
    }

    static List<Scenario> scenarios() {
        List<Scenario> scenarios = new ArrayList<>();

        scenarios.add(new Scenario("fail-level-valid",
                MessageTag.IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE,
                new Spec(Level.FAIL, true)));

        scenarios.add(new Scenario("fail-level-invalid",
                MessageTag.IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE,
                new Spec(Level.FAIL, false)));

        scenarios.add(new Scenario("fail-level-short-circuit",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.FAIL, false),
                new Spec(Level.FAIL, true).tags(MessageTag.BBB_CV_ISI, MessageTag.BBB_CV_ISI_ANS)));

        scenarios.add(new Scenario("fail-level-valid-then-invalid",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.FAIL, true),
                new Spec(Level.FAIL, false).tags(MessageTag.BBB_CV_ISI, MessageTag.BBB_CV_ISI_ANS)
                        .failure(Indication.FAILED, SubIndication.SIG_CRYPTO_FAILURE)));

        scenarios.add(new Scenario("warn-level-invalid",
                MessageTag.VALIDATION_CONTEXT_INITIALIZATION,
                new Spec(Level.WARN, false).tags(MessageTag.BBB_VCI_IZHSP, MessageTag.BBB_VCI_IZHSP_ANS)));

        scenarios.add(new Scenario("warn-level-valid",
                MessageTag.VALIDATION_CONTEXT_INITIALIZATION,
                new Spec(Level.WARN, true).tags(MessageTag.BBB_VCI_IZHSP, MessageTag.BBB_VCI_IZHSP_ANS)));

        scenarios.add(new Scenario("inform-level-invalid",
                MessageTag.VALIDATION_CONTEXT_INITIALIZATION,
                new Spec(Level.INFORM, false).tags(MessageTag.BBB_VCI_ISPSUPP, MessageTag.BBB_VCI_ISPSUPP_ANS)));

        scenarios.add(new Scenario("ignore-level-invalid",
                MessageTag.VALIDATION_CONTEXT_INITIALIZATION,
                new Spec(Level.IGNORE, false).additionalInfo(MessageTag.TOKEN_ID),
                new Spec(Level.FAIL, true).tags(MessageTag.BBB_CV_ISI, MessageTag.BBB_CV_ISI_ANS)));

        scenarios.add(new Scenario("undefined-constraint",
                null,
                new Spec(null, false),
                new Spec(Level.FAIL, true).tags(MessageTag.BBB_CV_ISI, MessageTag.BBB_CV_ISI_ANS)));

        scenarios.add(new Scenario("warn-info-then-fail",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.WARN, false).tags(MessageTag.BBB_VCI_IZHSP, MessageTag.BBB_VCI_IZHSP_ANS),
                new Spec(Level.INFORM, false).tags(MessageTag.BBB_VCI_ISPSUPP, MessageTag.BBB_VCI_ISPSUPP_ANS),
                new Spec(Level.FAIL, false).tags(MessageTag.BBB_CV_IRDOF, MessageTag.BBB_CV_IRDOF_ANS)
                        .failure(Indication.INDETERMINATE, SubIndication.SIGNED_DATA_NOT_FOUND)));

        scenarios.add(new Scenario("uninterrupted-continues-on-fail",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.FAIL, false).uninterrupted()
                        .tags(MessageTag.BBB_CV_IRDOF, MessageTag.BBB_CV_IRDOF_ANS),
                new Spec(Level.FAIL, false).uninterrupted()
                        .tags(MessageTag.BBB_CV_IRDOI, MessageTag.BBB_CV_IRDOI_ANS)
                        .failure(Indication.FAILED, SubIndication.HASH_FAILURE),
                new Spec(Level.WARN, false).tags(MessageTag.BBB_VCI_IZHSP, MessageTag.BBB_VCI_IZHSP_ANS)));

        scenarios.add(new Scenario("custom-success-conclusion",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.FAIL, true).success(Indication.PASSED, SubIndication.NO_POE),
                new Spec(Level.FAIL, false).tags(MessageTag.BBB_CV_ISI, MessageTag.BBB_CV_ISI_ANS)));

        scenarios.add(new Scenario("previous-errors",
                MessageTag.CRYPTOGRAPHIC_VERIFICATION,
                new Spec(Level.FAIL, false).tags(MessageTag.BBB_CV_IRDOF, MessageTag.BBB_CV_IRDOF_ANS)
                        .previousErrors(MessageTag.BBB_CV_ISI_ANS, MessageTag.BBB_ICS_ISCI_ANS)));

        scenarios.add(new Scenario("constraint-members-populated",
                MessageTag.IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE,
                new Spec(Level.FAIL, false).bbbId("S-1234").blockType(XmlBlockType.SIG_BBB)
                        .additionalInfo(MessageTag.EMPTY)));

        scenarios.add(new Scenario("no-title-passed",
                null,
                new Spec(Level.IGNORE, true)));

        return scenarios;
    }

    public static void main(String[] args) {
        I18nProvider i18nProvider = new I18nProvider();
        StringBuilder out = new StringBuilder();
        for (Scenario scenario : scenarios()) {
            SpecChain chain = new SpecChain(i18nProvider, scenario.title, scenario.specs);
            XmlConstraintsConclusion result = chain.execute();
            out.setLength(0);
            out.append('{');
            key(out, "scenario").append(str(scenario.name)).append(',');
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
            out.append(']').append('}');
            System.out.println(out);
        }
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

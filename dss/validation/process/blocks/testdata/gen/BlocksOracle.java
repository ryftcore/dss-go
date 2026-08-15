/*
 * Java oracle driver for the 5.2 building-block dispatcher of phase 8d:
 *
 *   eu.europa.esig.dss.validation.process.bbb.BasicBuildingBlocks
 *
 * (ported here as package blocks - see basic_building_blocks.go's header for why
 * it cannot live in the Go package bbb).
 *
 * It runs the UPSTREAM class over REAL inputs - every XmlDiagnosticData dump of
 * the marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle plus the nine
 * synthetic dumps in ../../bbb/xcv/testdata/dd - with the default ETSI validation
 * policy and the fixed validation time 2024-01-01T00:00:00Z, over every signature
 * (Context.SIGNATURE / COUNTER_SIGNATURE), time-stamp (TIMESTAMP), revocation
 * (REVOCATION) and used certificate (CERTIFICATE) of each dump, and dumps one
 * JSON object per line.
 *
 * A row carries what the dispatcher itself decides: which of the seven sub-blocks
 * it instantiated at all, each one's conclusion, the aggregated final conclusion
 * updateFinalConclusion() built out of them, the certificate chain it copied off
 * the ISC block, and the cross/equivalent certificate lists addAdditionalInfo()
 * hung off each XmlSubXCV.
 *
 * Rows are written to
 *   dss/validation/process/blocks/testdata/oracle/blocks.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle BlocksOracle.java
 *   java  -cp "$CP:/tmp/oracle" BlocksOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlBasicBuildingBlocks;
import eu.europa.esig.dss.detailedreport.jaxb.XmlChainItem;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSubXCV;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.TokenProxy;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;
import eu.europa.esig.dss.validation.process.bbb.BasicBuildingBlocks;

import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class BlocksOracle {

    /** Fixed validation time: 2024-01-01T00:00:00Z. */
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

        try (PrintWriter out = writer(repoRoot, "validation/process/blocks/testdata/oracle/blocks.jsonl")) {
            for (File synth : synthetic) {
                emit(out, i18n, policy, "dd/" + synth.getName(), new DiagnosticData(
                        DiagnosticDataFacade.newFacade().unmarshall(synth, false)));
            }
            for (File file : files) {
                String name = file.getName();
                // See ../../bbb/fc/testdata/README.md for why the four model-*.xml
                // schema-coverage fixtures are excluded.
                if (name.startsWith("model-")) {
                    System.err.println("SKIPPED-FIXTURE " + name);
                    continue;
                }
                XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                emit(out, i18n, policy, name, new DiagnosticData(jaxb));
            }
        }
    }

    private static void emit(PrintWriter out, I18nProvider i18n, ValidationPolicy policy,
                             String file, DiagnosticData dd) {
        setCurrent(dd);
        Map<String, XmlBasicBuildingBlocks> bbbs = new HashMap<>();
        for (SignatureWrapper signature : dd.getSignatures()) {
            Context context = signature.isCounterSignature() ? Context.COUNTER_SIGNATURE : Context.SIGNATURE;
            row(out, i18n, policy, file, "SIG|" + signature.getId(), signature, context, bbbs);
        }
        for (TimestampWrapper timestamp : dd.getTimestampList()) {
            row(out, i18n, policy, file, "TST|" + timestamp.getId(), timestamp, Context.TIMESTAMP, bbbs);
        }
        List<RevocationWrapper> revocations = new ArrayList<>(dd.getAllRevocationData());
        revocations.sort(Comparator.comparing(RevocationWrapper::getId));
        for (RevocationWrapper revocation : revocations) {
            row(out, i18n, policy, file, "REV|" + revocation.getId(), revocation, Context.REVOCATION, bbbs);
        }
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            row(out, i18n, policy, file, "CERT|" + certificate.getId(), certificate, Context.CERTIFICATE, bbbs);
        }
    }

    private static void row(PrintWriter out, I18nProvider i18n, ValidationPolicy policy, String file,
                            String token, TokenProxy proxy, Context context,
                            Map<String, XmlBasicBuildingBlocks> bbbs) {
        XmlBasicBuildingBlocks result;
        try {
            result = new BasicBuildingBlocks(i18n, dataOf(proxy), proxy, CURRENT_TIME, bbbs, policy, context).execute();
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + ": " + e);
            return;
        }
        StringBuilder row = new StringBuilder("{");
        key(row, "file").append(str(file)).append(',');
        key(row, "token").append(str(token)).append(',');
        key(row, "context").append(str(context.name())).append(',');
        key(row, "id").append(str(result.getId())).append(',');
        key(row, "type").append(str(result.getType() == null ? null : result.getType().name())).append(',');
        key(row, "conclusion").append(conclusion(result.getConclusion())).append(',');
        key(row, "fc").append(block(result.getFC())).append(',');
        key(row, "isc").append(block(result.getISC())).append(',');
        key(row, "vci").append(block(result.getVCI())).append(',');
        key(row, "aov").append(block(result.getAOV())).append(',');
        key(row, "xcv").append(block(result.getXCV())).append(',');
        key(row, "cv").append(block(result.getCV())).append(',');
        key(row, "sav").append(block(result.getSAV())).append(',');
        key(row, "certificateChain").append(chain(result)).append(',');
        key(row, "subXCV").append(subXCVs(result));
        out.println(row.append('}').toString());
    }

    /** The DiagnosticData the wrapper belongs to; the wrappers carry no back-reference. */
    private static DiagnosticData currentDiagnosticData;

    private static DiagnosticData dataOf(TokenProxy proxy) {
        return currentDiagnosticData;
    }

    // ------------------------------------------------------------------- JSON

    private static String block(XmlConstraintsConclusion block) {
        if (block == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "title").append(str(block.getTitle())).append(',');
        key(out, "conclusion").append(conclusion(block.getConclusion()));
        return out.append('}').toString();
    }

    private static String chain(XmlBasicBuildingBlocks result) {
        if (result.getCertificateChain() == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("[");
        List<XmlChainItem> items = result.getCertificateChain().getChainItem();
        for (int i = 0; i < items.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(str(items.get(i).getId()));
        }
        return out.append(']').toString();
    }

    private static String subXCVs(XmlBasicBuildingBlocks result) {
        if (result.getXCV() == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("[");
        List<XmlSubXCV> subs = result.getXCV().getSubXCV();
        for (int i = 0; i < subs.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            XmlSubXCV sub = subs.get(i);
            out.append('{');
            key(out, "id").append(str(sub.getId())).append(',');
            key(out, "crossCertificates").append(strings(sub.getCrossCertificates())).append(',');
            key(out, "equivalentCertificates").append(strings(sub.getEquivalentCertificates())).append(',');
            key(out, "conclusion").append(conclusion(sub.getConclusion()));
            out.append('}');
        }
        return out.append(']').toString();
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
            XmlMessage message = messages.get(i);
            out.append('{');
            key(out, "key").append(str(message.getKey())).append(',');
            key(out, "value").append(str(message.getValue()));
            out.append('}');
        }
        return out.append(']').toString();
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

    static void setCurrent(DiagnosticData dd) {
        currentDiagnosticData = dd;
    }
}

/*
 * Java oracle driver for the POE core of EN 319 102-1 5.6.2.3 / 5.6.2.4:
 *
 *   eu.europa.esig.dss.validation.process.vpfswatsp.POE
 *   eu.europa.esig.dss.validation.process.vpfswatsp.TimestampPOE
 *   eu.europa.esig.dss.validation.process.vpfswatsp.EvidenceRecordPOE
 *   eu.europa.esig.dss.validation.process.vpfswatsp.POEComparator
 *   eu.europa.esig.dss.validation.process.vpfswatsp.POEExtraction
 *
 * It runs the UPSTREAM classes over REAL inputs - every XmlDiagnosticData dump of
 * the marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle - plus this
 * package's synthetic dumps in ../dd (see POESyntheticDumps.java), with the fixed
 * control time 2024-01-01T00:00:00Z, and dumps one JSON object per dump.
 *
 * A row carries, in this order:
 *
 *   - "tokens": for every token id init() seeds (in the order init() walks them,
 *     which is the order the row is to be replayed in), the lowest POE after
 *     init(), after collectAllPOE(timestamps), and after extracting every
 *     evidence record - each as {time, providerId, tokenProvided} - plus the
 *     isPOEExists() answer at four probe times and, for used certificates, the
 *     isPOEExistInRange(notBefore, notAfter) answer;
 *   - "compare": POEComparator#compare over every ordered pair of the POEs the
 *     dump can build (the control-time POE, one TimestampPOE per time-stamp whose
 *     production time is set, one EvidenceRecordPOE per evidence record with a
 *     time-stamp), so that all four tie-breakers are observed directly;
 *   - "signaturePOE": the lowest POE of every signature after addSignaturePOE()
 *     has added the control-time POE a second time, which is what
 *     ValidationProcessForSignaturesWithArchivalData step 4) does.
 *
 * Rows are written to
 *   dss/validation/process/vpfswatsp/testdata/oracle/poe.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle POESyntheticDumps.java POEOracle.java
 *   java  -cp "$CP:/tmp/oracle" POESyntheticDumps <dss-repo-root>
 *   java  -cp "$CP:/tmp/oracle" POEOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.EAAWrapper;
import eu.europa.esig.dss.diagnostic.EvidenceRecordWrapper;
import eu.europa.esig.dss.diagnostic.OrphanTokenWrapper;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.SignerDataWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.validation.process.vpfswatsp.EvidenceRecordPOE;
import eu.europa.esig.dss.validation.process.vpfswatsp.POE;
import eu.europa.esig.dss.validation.process.vpfswatsp.POEComparator;
import eu.europa.esig.dss.validation.process.vpfswatsp.POEExtraction;
import eu.europa.esig.dss.validation.process.vpfswatsp.TimestampPOE;

import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class POEOracle {

    /** Fixed control time: 2024-01-01T00:00:00Z. */
    private static final Date CONTROL_TIME = new Date(1704067200000L);

    /** Probe times for isPOEExists: 2000, 2016-06-01, 2020, the control time. */
    private static final long[] PROBES = {
            946684800000L, 1464739200000L, 1577836800000L, 1704067200000L
    };

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        List<File> files = new ArrayList<>();
        files.addAll(sorted(inputDir));
        files.addAll(sorted(new File(repoRoot, "validation/process/vpfswatsp/testdata/dd")));

        File outFile = new File(repoRoot, "validation/process/vpfswatsp/testdata/oracle/poe.jsonl");
        outFile.getParentFile().mkdirs();

        try (PrintWriter out = new PrintWriter(outFile, "UTF-8")) {
            for (File file : files) {
                XmlDiagnosticData xml = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                DiagnosticData dd = new DiagnosticData(xml);
                out.println(row(file.getName(), dd));
            }
        }
    }

    private static List<File> sorted(File dir) {
        File[] files = dir.listFiles((d, name) -> name.endsWith(".xml"));
        if (files == null) {
            return new ArrayList<>();
        }
        Arrays.sort(files, Comparator.comparing(File::getName));
        return Arrays.asList(files);
    }

    private static String row(String name, DiagnosticData dd) {
        StringBuilder sb = new StringBuilder();
        sb.append('{');
        sb.append("\"dump\":").append(quote(name));

        // The token ids init() seeds, in init()'s own walk order.
        List<String> tokenIds = seededTokenIds(dd);

        POEExtraction poe = new POEExtraction();
        poe.init(dd, CONTROL_TIME);
        Map<String, Object[]> afterInit = snapshot(poe, tokenIds);

        poe.collectAllPOE(dd.getTimestampList());
        Map<String, Object[]> afterTimestamps = snapshot(poe, tokenIds);

        // A dump may carry an evidence record whose time-stamp reference was never
        // resolved; upstream's EvidenceRecordPOE then raises a NullPointerException
        // out of getFirstTimestamp(). The outcome is recorded rather than hidden,
        // so that the Go replay asserts the same failure.
        List<String> erExtraction = new ArrayList<>();
        for (EvidenceRecordWrapper er : dd.getEvidenceRecords()) {
            try {
                poe.extractPOE(er);
                erExtraction.add("ok");
            } catch (RuntimeException e) {
                erExtraction.add("throws");
            }
        }
        Map<String, Object[]> afterEvidenceRecords = snapshot(poe, tokenIds);

        sb.append(",\"tokens\":[");
        boolean first = true;
        for (String id : tokenIds) {
            if (!first) {
                sb.append(',');
            }
            first = false;
            sb.append('{');
            sb.append("\"id\":").append(quote(id));
            sb.append(",\"afterInit\":").append(poeJson(afterInit.get(id)));
            sb.append(",\"afterTimestamps\":").append(poeJson(afterTimestamps.get(id)));
            sb.append(",\"afterEvidenceRecords\":").append(poeJson(afterEvidenceRecords.get(id)));
            sb.append(",\"exists\":[");
            for (int i = 0; i < PROBES.length; i++) {
                if (i > 0) {
                    sb.append(',');
                }
                sb.append(poe.isPOEExists(id, new Date(PROBES[i])));
            }
            sb.append(']');
            CertificateWrapper certificate = certificateById(dd, id);
            if (certificate != null && certificate.getNotBefore() != null && certificate.getNotAfter() != null) {
                sb.append(",\"inValidityRange\":")
                        .append(poe.isPOEExistInRange(id, certificate.getNotBefore(), certificate.getNotAfter()));
            }
            sb.append('}');
        }
        sb.append(']');

        // POEComparator over every ordered pair of the POEs this dump can build.
        List<Object[]> poes = comparablePOEs(dd);
        POEComparator comparator = new POEComparator();
        sb.append(",\"erExtraction\":[");
        first = true;
        for (String outcome : erExtraction) {
            if (!first) {
                sb.append(',');
            }
            first = false;
            sb.append(quote(outcome));
        }
        sb.append(']');

        sb.append(",\"compare\":[");
        first = true;
        for (Object[] a : poes) {
            for (Object[] b : poes) {
                if (!first) {
                    sb.append(',');
                }
                first = false;
                sb.append('{');
                sb.append("\"a\":").append(quote((String) a[0]));
                sb.append(",\"b\":").append(quote((String) b[0]));
                sb.append(",\"compare\":").append(sign(comparator.compare((POE) a[1], (POE) b[1])));
                sb.append(",\"before\":").append(comparator.before((POE) a[1], (POE) b[1]));
                sb.append('}');
            }
        }
        sb.append(']');

        // addSignaturePOE, as ValidationProcessForSignaturesWithArchivalData step 4) does.
        POEExtraction sigPoe = new POEExtraction();
        sigPoe.init(dd, CONTROL_TIME);
        sigPoe.collectAllPOE(dd.getTimestampList());
        sb.append(",\"signaturePOE\":[");
        first = true;
        for (SignatureWrapper signature : dd.getAllSignatures()) {
            sigPoe.addSignaturePOE(signature, new POE(new Date(PROBES[1])));
            if (!first) {
                sb.append(',');
            }
            first = false;
            sb.append('{');
            sb.append("\"id\":").append(quote(signature.getId()));
            sb.append(",\"lowest\":").append(poeJson(describe(sigPoe.getLowestPOE(signature.getId()))));
            sb.append('}');
        }
        sb.append(']');

        sb.append('}');
        return sb.toString();
    }

    /** The token ids POEExtraction#init seeds, in the order init() walks them. */
    private static List<String> seededTokenIds(DiagnosticData dd) {
        List<String> ids = new ArrayList<>();
        for (SignatureWrapper w : dd.getAllSignatures()) {
            ids.add(w.getId());
        }
        for (TimestampWrapper w : dd.getTimestampList()) {
            ids.add(w.getId());
        }
        for (EvidenceRecordWrapper w : dd.getEvidenceRecords()) {
            ids.add(w.getId());
        }
        for (EAAWrapper w : dd.getEAAs()) {
            ids.add(w.getId());
        }
        for (CertificateWrapper w : dd.getUsedCertificates()) {
            ids.add(w.getId());
        }
        for (RevocationWrapper w : dd.getAllRevocationData()) {
            ids.add(w.getId());
        }
        for (SignerDataWrapper w : dd.getAllSignerDocuments()) {
            ids.add(w.getId());
        }
        for (OrphanTokenWrapper<?> w : dd.getAllOrphanCertificateObjects()) {
            ids.add(w.getId());
        }
        for (OrphanTokenWrapper<?> w : dd.getAllOrphanCertificateReferences()) {
            ids.add(w.getId());
        }
        for (OrphanTokenWrapper<?> w : dd.getAllOrphanRevocationObjects()) {
            ids.add(w.getId());
        }
        for (OrphanTokenWrapper<?> w : dd.getAllOrphanRevocationReferences()) {
            ids.add(w.getId());
        }
        // de-duplicate, keeping the first occurrence
        List<String> unique = new ArrayList<>();
        for (String id : ids) {
            if (!unique.contains(id)) {
                unique.add(id);
            }
        }
        return unique;
    }

    private static Map<String, Object[]> snapshot(POEExtraction poe, List<String> tokenIds) {
        Map<String, Object[]> result = new LinkedHashMap<>();
        for (String id : tokenIds) {
            result.put(id, describe(poe.getLowestPOE(id)));
        }
        return result;
    }

    private static Object[] describe(POE poe) {
        if (poe == null) {
            return null;
        }
        return new Object[] { poe.getTime(), poe.getPOEProviderId(), poe.isTokenProvided(),
                poe.getPOEObjects() == null ? -1 : poe.getPOEObjects().size() };
    }

    /** Every POE the dump can build, labelled, in a deterministic order. */
    private static List<Object[]> comparablePOEs(DiagnosticData dd) {
        List<Object[]> poes = new ArrayList<>();
        poes.add(new Object[] { "control-time", new POE(CONTROL_TIME) });
        for (TimestampWrapper timestamp : dd.getTimestampList()) {
            if (timestamp.getProductionTime() != null) {
                poes.add(new Object[] { "tst:" + timestamp.getId(), new TimestampPOE(timestamp) });
            }
        }
        for (EvidenceRecordWrapper er : dd.getEvidenceRecords()) {
            try {
                if (er.getFirstTimestamp() != null && er.getFirstTimestamp().getProductionTime() != null) {
                    poes.add(new Object[] { "er:" + er.getId(), new EvidenceRecordPOE(er) });
                }
            } catch (RuntimeException e) {
                // unresolved time-stamp reference; see the erExtraction note
            }
        }
        return poes;
    }

    private static CertificateWrapper certificateById(DiagnosticData dd, String id) {
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            if (certificate.getId().equals(id)) {
                return certificate;
            }
        }
        return null;
    }

    private static int sign(int value) {
        return Integer.compare(value, 0);
    }

    private static String poeJson(Object[] described) {
        if (described == null) {
            return "null";
        }
        StringBuilder sb = new StringBuilder();
        sb.append("{\"time\":").append(((Date) described[0]).getTime());
        sb.append(",\"providerId\":").append(quote((String) described[1]));
        sb.append(",\"tokenProvided\":").append(described[2]);
        sb.append(",\"objects\":").append(described[3]);
        sb.append('}');
        return sb.toString();
    }

    private static String quote(String value) {
        if (value == null) {
            return "null";
        }
        StringBuilder sb = new StringBuilder("\"");
        for (int i = 0; i < value.length(); i++) {
            char c = value.charAt(i);
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

}

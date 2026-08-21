import eu.europa.esig.dss.enumerations.TSLType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.tsl.CertificateContentEquivalence;
import eu.europa.esig.dss.model.tsl.ConditionForQualifiers;
import eu.europa.esig.dss.model.tsl.MRA;
import eu.europa.esig.dss.model.tsl.OtherTSLPointer;
import eu.europa.esig.dss.model.tsl.QCStatementOids;
import eu.europa.esig.dss.model.tsl.ServiceEquivalence;
import eu.europa.esig.dss.model.tsl.ServiceTypeASi;
import eu.europa.esig.dss.model.tsl.TrustService;
import eu.europa.esig.dss.model.tsl.TrustServiceProvider;
import eu.europa.esig.dss.model.tsl.TrustServiceStatusAndInformationExtensions;
import eu.europa.esig.dss.model.timedependent.MutableTimeDependentValues;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.tsl.parsing.LOTLParsingResult;
import eu.europa.esig.dss.tsl.parsing.LOTLParsingTask;
import eu.europa.esig.dss.tsl.parsing.TLParsingResult;
import eu.europa.esig.dss.tsl.parsing.TLParsingTask;
import eu.europa.esig.dss.tsl.source.LOTLSource;
import eu.europa.esig.dss.tsl.source.TLSource;
import eu.europa.esig.dss.tsl.function.XMLOtherTSLPointer;
import eu.europa.esig.dss.tsl.function.PDFOtherTSLPointer;
import eu.europa.esig.dss.tsl.function.OfficialJournalSchemeInformationURI;
import eu.europa.esig.dss.tsl.function.TLPredicateFactory;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.TreeMap;

/**
 * Phase 9 harness contract item (A): LOTL/TL/pivot/MRA parse parity. For every (kind, name,
 * path) row of a manifest TSV, parses the fixture with TLParsingTask/LOTLParsingTask over a
 * default TLSource/LOTLSource (LOTLSource gets setMraSupport(true) for mra-lotl.xml, matching
 * the Go side) and dumps every field of the resulting *ParsingResult as one JSON line.
 *
 * Map-shaped fields are dumped through a TreeMap with each value list sorted, since Java
 * HashMap iteration order is unspecified and unrelated to any real defect - see
 * dss/harness/tsl_parsing_oracle_test.go's header for the full rationale. List-shaped fields
 * that reflect real XML document order (TSP list, per-service time-dependent history,
 * distribution points, pivot URLs, positional certificate lists) are left in encounter order.
 */
public class TSLParsingOracle {

    public static void main(String[] args) throws Exception {
        Locale.setDefault(Locale.ENGLISH);
        Path manifestPath = Paths.get(args[0]);
        Path outFile = Paths.get(args[1]);

        StringBuilder out = new StringBuilder();
        int rows = 0;
        try (BufferedReader r = new BufferedReader(new FileReader(manifestPath.toFile(), StandardCharsets.UTF_8))) {
            String line;
            while ((line = r.readLine()) != null) {
                if (line.isBlank()) {
                    continue;
                }
                String[] parts = line.split("\t");
                String kind = parts[0];
                String name = parts[1];
                String path = parts[2];
                String opts = parts.length > 3 ? parts[3] : "";

                out.append("{\"name\":").append(json(name));
                out.append(",\"kind\":").append(json(kind));

                DSSDocument doc = new FileDocument(new File(path));
                try {
                    if ("LOTL".equals(kind)) {
                        LOTLSource source = new LOTLSource();
                        applyOptions(source, opts);
                        LOTLParsingResult result;
                        try {
                            result = new LOTLParsingTask(doc, source).get();
                        } catch (Exception e) {
                            out.append(",\"parseError\":true}\n");
                            rows++;
                            continue;
                        }
                        out.append(",\"parseError\":false");
                        dumpCommon(out, result.getTSLType(), result.getSequenceNumber(), result.getVersion(),
                                result.getTerritory(), result.getIssueDate(), result.getNextUpdateDate(),
                                result.getDistributionPoints());
                        out.append(",\"tsps\":[]");
                        out.append(",\"lotlPointers\":");
                        dumpPointers(out, result.getLotlPointers());
                        out.append(",\"tlPointers\":");
                        dumpPointers(out, result.getTlPointers());
                        out.append(",\"signingCertAnnouncementURL\":").append(json(result.getSigningCertificateAnnouncementURL()));
                        out.append(",\"pivotURLs\":");
                        dumpStrList(out, result.getPivotURLs());
                    } else {
                        TLSource source = new TLSource();
                        TLParsingResult result;
                        try {
                            result = new TLParsingTask(doc, source).get();
                        } catch (Exception e) {
                            out.append(",\"parseError\":true}\n");
                            rows++;
                            continue;
                        }
                        out.append(",\"parseError\":false");
                        dumpCommon(out, result.getTSLType(), result.getSequenceNumber(), result.getVersion(),
                                result.getTerritory(), result.getIssueDate(), result.getNextUpdateDate(),
                                result.getDistributionPoints());
                        out.append(",\"tsps\":");
                        dumpTSPs(out, result.getTrustServiceProviders());
                        out.append(",\"lotlPointers\":[]");
                        out.append(",\"tlPointers\":[]");
                        out.append(",\"signingCertAnnouncementURL\":\"\"");
                        out.append(",\"pivotURLs\":[]");
                    }
                } catch (Throwable t) {
                    throw new RuntimeException("fixture " + name + " failed", t);
                }
                out.append("}\n");
                rows++;
            }
        }
        Files.write(outFile, out.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("rows=" + rows);
    }

    /**
     * Applies the manifest's option column to a LOTLSource. Options are separated by ';';
     * recognised keys mirror the knobs upstream's own LOTLParsingTaskTest exercises:
     * "mra" (setMraSupport), "pivot" (setPivotSupport), "xmlpred"/"pdfpred"
     * (setTlPredicate with the XML/PDF OtherTSLPointer predicate - the MRA pointer of the
     * mra-* fixtures is a non-EUgeneric TSLType and is filtered out by the DEFAULT
     * EUTLOtherTSLPointer predicate, so without this the whole MRA tree parses as null),
     * "cc=XX,YY" (TLPredicateFactory.createEUTLCountryCodePredicate) and "oj=&lt;url&gt;"
     * (setSigningCertificatesAnnouncementPredicate with OfficialJournalSchemeInformationURI).
     */
    private static void applyOptions(LOTLSource source, String opts) {
        if (opts == null || opts.isBlank()) {
            return;
        }
        for (String opt : opts.split(";")) {
            opt = opt.trim();
            if (opt.isEmpty()) {
                continue;
            }
            if ("mra".equals(opt)) {
                source.setMraSupport(true);
            } else if ("pivot".equals(opt)) {
                source.setPivotSupport(true);
            } else if ("xmlpred".equals(opt)) {
                source.setTlPredicate(new XMLOtherTSLPointer());
            } else if ("pdfpred".equals(opt)) {
                source.setTlPredicate(new PDFOtherTSLPointer());
            } else if (opt.startsWith("cc=")) {
                source.setTlPredicate(TLPredicateFactory.createEUTLCountryCodePredicate(
                        opt.substring(3).split(",")));
            } else if (opt.startsWith("oj=")) {
                source.setSigningCertificatesAnnouncementPredicate(
                        new OfficialJournalSchemeInformationURI(opt.substring(3)));
            } else {
                throw new IllegalArgumentException("Unknown option: " + opt);
            }
        }
    }

    private static void dumpCommon(StringBuilder out, TSLType tslType, Integer sequenceNumber, Integer version,
                                    String territory, Date issueDate, Date nextUpdateDate, List<String> distributionPoints) {
        out.append(",\"tslType\":").append(json(tslType == null ? "" : tslType.getUri()));
        out.append(",\"sequenceNumber\":").append(sequenceNumber == null ? Integer.MIN_VALUE : sequenceNumber);
        out.append(",\"version\":").append(version == null ? Integer.MIN_VALUE : version);
        out.append(",\"territory\":").append(json(territory));
        out.append(",\"issueDate\":").append(epochMillis(issueDate));
        out.append(",\"nextUpdateDate\":").append(epochMillis(nextUpdateDate));
        out.append(",\"distributionPoints\":");
        dumpStrList(out, distributionPoints);
    }

    private static void dumpTSPs(StringBuilder out, List<TrustServiceProvider> tsps) {
        out.append("[");
        if (tsps != null) {
            boolean first = true;
            for (TrustServiceProvider tsp : tsps) {
                if (!first) out.append(",");
                first = false;
                out.append("{\"names\":");
                dumpSortedMultiMap(out, tsp.getNames());
                out.append(",\"tradeNames\":");
                dumpSortedMultiMap(out, tsp.getTradeNames());
                out.append(",\"registrationIdentifiers\":");
                dumpStrList(out, tsp.getRegistrationIdentifiers());
                out.append(",\"postalAddresses\":");
                dumpSortedStrMap(out, tsp.getPostalAddresses());
                out.append(",\"electronicAddresses\":");
                dumpSortedMultiMap(out, tsp.getElectronicAddresses());
                out.append(",\"information\":");
                dumpSortedStrMap(out, tsp.getInformation());
                out.append(",\"territory\":").append(json(tsp.getTerritory()));
                out.append(",\"services\":");
                dumpServices(out, tsp.getServices());
                out.append("}");
            }
        }
        out.append("]");
    }

    private static void dumpServices(StringBuilder out, List<TrustService> services) {
        out.append("[");
        boolean first = true;
        if (services != null) {
            for (TrustService svc : services) {
                List<String> certDigests = certSha256(svc.getCertificates());
                for (TrustServiceStatusAndInformationExtensions e : svc.getStatusAndInformationExtensions()) {
                    if (!first) out.append(",");
                    first = false;
                    out.append("{\"start\":").append(epochMillis(e.getStartDate()));
                    out.append(",\"end\":").append(epochMillis(e.getEndDate()));
                    out.append(",\"names\":");
                    dumpSortedMultiMap(out, e.getNames());
                    out.append(",\"type\":").append(json(e.getType()));
                    out.append(",\"status\":").append(json(e.getStatus()));
                    out.append(",\"additionalServiceInfoUris\":");
                    dumpStrList(out, e.getAdditionalServiceInfoUris());
                    out.append(",\"serviceSupplyPoints\":");
                    dumpStrList(out, e.getServiceSupplyPoints());
                    out.append(",\"expiredCertsRevocationInfo\":").append(epochMillis(e.getExpiredCertsRevocationInfo()));
                    out.append(",\"certificatesSha256\":");
                    dumpStrList(out, certDigests);
                    out.append(",\"conditions\":");
                    dumpConditions(out, e.getConditionsForQualifiers());
                    out.append("}");
                }
            }
        }
        out.append("]");
    }

    private static void dumpConditions(StringBuilder out, List<ConditionForQualifiers> cfqs) {
        out.append("[");
        if (cfqs != null) {
            boolean first = true;
            for (ConditionForQualifiers c : cfqs) {
                if (!first) out.append(",");
                first = false;
                out.append("{\"qualifiers\":");
                dumpStrList(out, c.getQualifiers());
                out.append(",\"critical\":").append(c.isCritical());
                out.append(",\"condition\":").append(json(conditionString(c.getCondition())));
                out.append("}");
            }
        }
        out.append("]");
    }

    /**
     * Renders a Condition through Condition#toString(String) and normalises the two places that
     * tree can print a Java null: QCStatementCondition#toString writes "type: null" /
     * "legislation: null" for an absent QcType / QcCClegislation. Go has no null String, so the
     * port's QCStatementCondition holds "" there and prints "type: " - the values are equivalent
     * everywhere either engine READS them (both are the "absent" case Utils.isStringNotEmpty
     * rejects in check()), and this rendering is debug-only. Nothing else in the condition tree
     * can produce the token "null" at the start of a line: CompositeCondition prints an enum
     * name, KeyUsageCondition a name plus true/false, PolicyIdCondition a non-null OID, and the
     * two list conditions are built from lists this port keeps non-null on both sides.
     */
    private static String conditionString(eu.europa.esig.dss.model.tsl.Condition condition) {
        if (condition == null) {
            return "";
        }
        return condition.toString("").replace(": null\n", ": \n");
    }

    private static void dumpPointers(StringBuilder out, List<OtherTSLPointer> pointers) {
        out.append("[");
        if (pointers != null) {
            boolean first = true;
            for (OtherTSLPointer p : pointers) {
                if (!first) out.append(",");
                first = false;
                out.append("{\"location\":").append(json(p.getLocation()));
                out.append(",\"tslLocation\":").append(json(p.getTSLLocation()));
                out.append(",\"schemeTerritory\":").append(json(p.getSchemeTerritory()));
                out.append(",\"tslType\":").append(json(p.getTslType()));
                out.append(",\"mimeType\":").append(json(p.getMimeType()));
                out.append(",\"schemeOperatorNames\":");
                dumpSortedMultiMap(out, p.getSchemeOperatorNames());
                out.append(",\"schemeTypeCommunityRules\":");
                dumpSortedMultiMap(out, p.getSchemeTypeCommunityRules());
                out.append(",\"sdiCertificatesSha256\":");
                dumpStrList(out, certSha256(p.getSdiCertificates()));
                out.append(",\"mra\":");
                dumpMra(out, p.getMra());
                out.append("}");
            }
        }
        out.append("]");
    }

    private static void dumpMra(StringBuilder out, MRA mra) {
        if (mra == null) {
            out.append("null");
            return;
        }
        out.append("{\"technicalType\":").append(json(mra.getTechnicalType()));
        out.append(",\"version\":").append(json(mra.getVersion()));
        out.append(",\"pointingLegislation\":").append(json(mra.getPointingContractingPartyLegislation()));
        out.append(",\"pointedLegislation\":").append(json(mra.getPointedContractingPartyLegislation()));
        out.append(",\"serviceEquivalences\":[");
        boolean first = true;
        if (mra.getServiceEquivalence() != null) {
            for (MutableTimeDependentValues<ServiceEquivalence> mtdv : mra.getServiceEquivalence()) {
                if (mtdv == null) continue;
                for (ServiceEquivalence se : mtdv.getList()) {
                    if (!first) out.append(",");
                    first = false;
                    dumpServiceEquivalence(out, se);
                }
            }
        }
        out.append("]}");
    }

    private static void dumpServiceEquivalence(StringBuilder out, ServiceEquivalence se) {
        out.append("{\"start\":").append(epochMillis(se.getStartDate()));
        out.append(",\"end\":").append(epochMillis(se.getEndDate()));
        out.append(",\"legalInfoIdentifier\":").append(json(se.getLegalInfoIdentifier()));
        out.append(",\"status\":").append(json(se.getStatus() == null ? "" : se.getStatus().name()));
        out.append(",\"typeAsiEquivalence\":");
        Map<String, String> typeAsi = new TreeMap<>();
        if (se.getTypeAsiEquivalence() != null) {
            for (Map.Entry<ServiceTypeASi, ServiceTypeASi> en : se.getTypeAsiEquivalence().entrySet()) {
                // A null ServiceTypeASi type/asi is "" on the Go side, which has no null String.
                typeAsi.put(nullToEmpty(en.getKey().getType()) + "|" + nullToEmpty(en.getKey().getAsi()),
                        nullToEmpty(en.getValue().getType()) + "|" + nullToEmpty(en.getValue().getAsi()));
            }
        }
        dumpSortedStrMap(out, typeAsi);
        // getStatusEquivalence() is a HashMap; its ENTRY iteration order is unspecified and
        // cannot be a parity target, so the entries are emitted in a canonical order (sorted by
        // their sorted key list). The Go side sorts its own slice with the identical key.
        out.append(",\"statusEquivalence\":[");
        List<String> statusEntries = new ArrayList<>();
        if (se.getStatusEquivalence() != null) {
            for (Map.Entry<List<String>, List<String>> en : se.getStatusEquivalence().entrySet()) {
                StringBuilder entry = new StringBuilder();
                entry.append("{\"pointedStatuses\":");
                dumpSortedStrList(entry, en.getKey());
                entry.append(",\"pointingStatuses\":");
                dumpSortedStrList(entry, en.getValue());
                entry.append("}");
                statusEntries.add(entry.toString());
            }
        }
        statusEntries.sort(String::compareTo);
        boolean first = true;
        for (String entry : statusEntries) {
            if (!first) out.append(",");
            first = false;
            out.append(entry);
        }
        out.append("]");
        out.append(",\"certificateContentEquivalences\":[");
        first = true;
        if (se.getCertificateContentEquivalences() != null) {
            for (CertificateContentEquivalence cce : se.getCertificateContentEquivalences()) {
                if (!first) out.append(",");
                first = false;
                out.append("{\"context\":").append(json(cce.getContext() == null ? "" : cce.getContext().name()));
                out.append(",\"condition\":").append(json(conditionString(cce.getCondition())));
                out.append(",\"contentReplacement\":");
                dumpQcStatementOids(out, cce.getContentReplacement());
                out.append("}");
            }
        }
        out.append("]");
        out.append(",\"qualifierEquivalence\":");
        dumpSortedStrMap(out, se.getQualifierEquivalence());
        out.append("}");
    }

    private static void dumpQcStatementOids(StringBuilder out, QCStatementOids q) {
        if (q == null) {
            out.append("null");
            return;
        }
        out.append("{\"qcStatementIds\":");
        dumpStrList(out, q.getQcStatementIds());
        out.append(",\"qcTypeIds\":");
        dumpStrList(out, q.getQcTypeIds());
        out.append(",\"qcCClegislations\":");
        dumpStrList(out, q.getQcCClegislations());
        out.append(",\"qcStatementIdsToRemove\":");
        dumpStrList(out, q.getQcStatementIdsToRemove());
        out.append(",\"qcTypeIdsToRemove\":");
        dumpStrList(out, q.getQcTypeIdsToRemove());
        out.append(",\"qcCClegislationsToRemove\":");
        dumpStrList(out, q.getQcCClegislationsToRemove());
        out.append("}");
    }

    private static List<String> certSha256(List<CertificateToken> certs) {
        List<String> out = new ArrayList<>();
        if (certs == null) {
            return out;
        }
        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            for (CertificateToken c : certs) {
                byte[] digest = md.digest(c.getEncoded());
                StringBuilder hex = new StringBuilder();
                for (byte b : digest) {
                    hex.append(String.format("%02x", b));
                }
                out.add(hex.toString());
            }
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
        return out;
    }

    private static String nullToEmpty(String s) {
        return s == null ? "" : s;
    }

    private static long epochMillis(Date d) {
        return d == null ? 0 : d.getTime();
    }

    private static void dumpStrList(StringBuilder out, List<String> list) {
        out.append("[");
        if (list != null) {
            boolean first = true;
            for (String s : list) {
                if (!first) out.append(",");
                first = false;
                out.append(json(s));
            }
        }
        out.append("]");
    }

    private static void dumpSortedStrList(StringBuilder out, List<String> list) {
        List<String> sorted = list == null ? new ArrayList<>() : new ArrayList<>(list);
        sorted.sort(String::compareTo);
        dumpStrList(out, sorted);
    }

    private static void dumpSortedStrMap(StringBuilder out, Map<String, String> map) {
        Map<String, String> sorted = new TreeMap<>();
        if (map != null) {
            for (Map.Entry<String, String> en : map.entrySet()) {
                // A null xml:lang key (a TL that omits the attribute - tl-ecdsa-brainpool.xml
                // has one) is the empty string on the Go side, which has no null String.
                sorted.put(en.getKey() == null ? "" : en.getKey(), en.getValue());
            }
        }
        out.append("{");
        boolean first = true;
        for (Map.Entry<String, String> en : sorted.entrySet()) {
            if (!first) out.append(",");
            first = false;
            out.append(json(en.getKey())).append(":").append(json(en.getValue()));
        }
        out.append("}");
    }

    private static void dumpSortedMultiMap(StringBuilder out, Map<String, List<String>> map) {
        Map<String, List<String>> sorted = new TreeMap<>();
        if (map != null) {
            for (Map.Entry<String, List<String>> en : map.entrySet()) {
                // See dumpSortedStrMap: a null xml:lang key is "" on the Go side.
                sorted.put(en.getKey() == null ? "" : en.getKey(), en.getValue());
            }
        }
        out.append("{");
        boolean first = true;
        for (Map.Entry<String, List<String>> en : sorted.entrySet()) {
            if (!first) out.append(",");
            first = false;
            out.append(json(en.getKey())).append(":");
            dumpSortedStrList(out, en.getValue());
        }
        out.append("}");
    }

    private static String json(String s) {
        if (s == null) {
            s = "";
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
        sb.append("\"");
        return sb.toString();
    }
}

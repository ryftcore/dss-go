// Java oracle #3 for the Go dss/asic port: the CONTAINER MERGE differential.
//
// dss-asic-*'s merger hierarchy (ASiCContainerMerger / DefaultContainerMerger /
// ASiC{S,E}With{CAdES,XAdES}ContainerMerger) decides, for a pair of containers, whether they may
// be merged at all and - when they may - what the merged container's entry inventory is. Those
// rules guard against silently corrupting signed containers, so the Go port is diffed against
// this oracle rather than against hand-derived expectations.
//
// For every unordered pair (including each container with itself) drawn from the ASiC containers
// in the fixture corpus, this dumps either the merge result (final container name, container
// type, zip comment, and every merged entry's name + SHA-256) or the exact exception class and
// message upstream raised. Output is NDJSON - one self-contained JSON object per line - so the
// ~16k pairs stay streamable.
//
// Build/run (OpenJDK 21, against a built upstream DSS 6.5.RC1):
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-asic-cades dependency:build-classpath -Dmdep.outputFile=/tmp/cp.txt -Dmdep.includeScope=runtime
//   CP="dss-asic-cades/target/classes:dss-asic-xades/target/classes:dss-asic-common/target/classes:...:$(cat /tmp/cp.txt)"
//   javac -cp "$CP" -d /tmp/mergeoracle MergeOracle.java
//   java  -cp "$CP:/tmp/mergeoracle" MergeOracle <fixtures root> <fixtures.txt> <out.ndjson>
import eu.europa.esig.dss.asic.common.ASiCUtils;
import eu.europa.esig.dss.asic.common.ZipUtils;
import eu.europa.esig.dss.asic.common.merge.ASiCContainerMerger;
import eu.europa.esig.dss.asic.common.merge.DefaultContainerMerger;
import eu.europa.esig.dss.enumerations.ASiCContainerType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;

import java.io.InputStream;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;

public final class MergeOracle {

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        List<String> relPaths = Files.readAllLines(Paths.get(args[1]));

        // Keep only the ASiC containers: the mergers reject anything else up front, so pairing
        // non-ASiC files would only re-test ASiCUtils.isASiC.
        List<String> asicPaths = new ArrayList<>();
        for (String rel : relPaths) {
            if (rel.trim().isEmpty()) {
                continue;
            }
            try {
                if (ASiCUtils.isASiC(new FileDocument(root.resolve(rel).toFile()))) {
                    asicPaths.add(rel);
                }
            } catch (Throwable t) {
                // not a readable container; skipped exactly as the Go side skips it
            }
        }
        System.err.println("asic containers: " + asicPaths.size());

        int pairs = 0;
        try (PrintStream out = new PrintStream(new java.io.FileOutputStream(args[2]), false, "UTF-8")) {
            for (int i = 0; i < asicPaths.size(); i++) {
                for (int j = i; j < asicPaths.size(); j++) {
                    out.println(mergePair(root, asicPaths.get(i), asicPaths.get(j)));
                    pairs++;
                }
                out.flush();
                System.err.println("row " + i + " done, pairs=" + pairs);
            }
        }
        System.err.println("pairs=" + pairs);
    }

    private static String mergePair(Path root, String relA, String relB) {
        StringBuilder sb = new StringBuilder("{");
        sb.append("\"a\": ").append(json(relA)).append(", \"b\": ").append(json(relB));
        try {
            DSSDocument docA = new FileDocument(root.resolve(relA).toFile());
            DSSDocument docB = new FileDocument(root.resolve(relB).toFile());
            ASiCContainerMerger merger = DefaultContainerMerger.fromDocuments(docA, docB);
            DSSDocument merged = merger.merge();

            sb.append(", \"ok\": true");
            sb.append(", \"name\": ").append(json(merged.getName()));
            sb.append(", \"mimeType\": ").append(
                    json(merged.getMimeType() == null ? null : merged.getMimeType().getMimeTypeString()));
            ASiCContainerType type = ASiCUtils.getContainerType(merged);
            sb.append(", \"containerType\": ").append(type == null ? "null" : json(type.name()));
            sb.append(", \"zipComment\": ").append(json(ASiCUtils.getZipComment(merged)));

            List<DSSDocument> entries = ZipUtils.getInstance().extractContainerContent(merged);
            sb.append(", \"entries\": [");
            for (int k = 0; k < entries.size(); k++) {
                if (k > 0) {
                    sb.append(", ");
                }
                DSSDocument entry = entries.get(k);
                sb.append("{\"name\": ").append(json(entry.getName()))
                  .append(", \"sha256\": ").append(json(sha256(entry))).append("}");
            }
            sb.append("]");
        } catch (Throwable t) {
            sb.append(", \"ok\": false");
            sb.append(", \"errorClass\": ").append(json(t.getClass().getSimpleName()));
            sb.append(", \"errorMessage\": ").append(json(t.getMessage()));
        }
        return sb.append("}").toString();
    }

    private static String sha256(DSSDocument document) {
        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            try (InputStream is = document.openStream()) {
                byte[] buf = new byte[8192];
                int n;
                while ((n = is.read(buf)) != -1) {
                    md.update(buf, 0, n);
                }
            }
            StringBuilder sb = new StringBuilder();
            for (byte b : md.digest()) {
                sb.append(String.format("%02x", b));
            }
            return sb.toString();
        } catch (Exception e) {
            return "ERROR:" + e.getClass().getSimpleName();
        }
    }

    private static String json(String s) {
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
                    if (c < 0x20 || c > 0x7e) {
                        sb.append(String.format("\\u%04x", (int) c));
                    } else {
                        sb.append(c);
                    }
            }
        }
        return sb.append("\"").toString();
    }
}

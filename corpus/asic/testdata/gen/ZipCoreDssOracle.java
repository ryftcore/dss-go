// Java oracle #2 for the Go dss/asic ZIPCORE port: runs the REAL upstream DSS 6.5.RC1 classes
// (ASiCUtils, ZipUtils, SecureContainerHandler, DSSZipEntry, AbstractASiCFilenameFactory) over the
// fixture corpus and dumps everything the Go port has to reproduce, including the exception class
// and message on the failure paths.
import eu.europa.esig.dss.asic.common.ASiCUtils;
import eu.europa.esig.dss.asic.common.AbstractASiCFilenameFactory;
import eu.europa.esig.dss.asic.common.DSSZipEntry;
import eu.europa.esig.dss.asic.common.DSSZipEntryDocument;
import eu.europa.esig.dss.asic.common.ZipUtils;
import eu.europa.esig.dss.enumerations.ASiCContainerType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;

import java.io.File;
import java.io.InputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

public final class ZipCoreDssOracle {

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        List<String> relPaths = Files.readAllLines(Paths.get(args[1]));
        try (PrintStream out = new PrintStream(new java.io.FileOutputStream(args[2]), false, "UTF-8")) {
            out.println("{");
            out.println("  \"filenameSuffixes\": [");
            dumpFilenameSuffixes(out);
            out.println("  ],");
            out.println("  \"containers\": [");
            boolean first = true;
            for (String rel : relPaths) {
                if (rel.trim().isEmpty()) {
                    continue;
                }
                if (!first) {
                    out.println(",");
                }
                first = false;
                dumpContainer(out, rel, root.resolve(rel).toFile());
            }
            out.println();
            out.println("  ]");
            out.println("}");
        }
    }

    /** Drives AbstractASiCFilenameFactory#getNextAvailableDocumentName through reflection. */
    private static void dumpFilenameSuffixes(PrintStream out) throws Exception {
        AbstractASiCFilenameFactory factory = new AbstractASiCFilenameFactory() { };
        Method method = AbstractASiCFilenameFactory.class.getDeclaredMethod(
                "getNextAvailableDocumentName", String.class, java.util.Collection.class);
        method.setAccessible(true);

        String[][] cases = {
                { "META-INF/signature001.p7s" },
                { "META-INF/signature001.p7s", "META-INF/signature001.p7s" },
                { "META-INF/signature001.p7s", "META-INF/signature001.p7s", "META-INF/signature002.p7s" },
                { "META-INF/signatures001.xml", "META-INF/signatures002.xml", "META-INF/signatures003.xml" },
                { "META-INF/ASiCManifest001.xml", "META-INF/ASiCManifest003.xml" },
                { "META-INF/signature001.p7s", "META-INF/signature001.p7s", "META-INF/signature001.p7s" },
        };
        for (int i = 0; i < cases.length; i++) {
            String template = cases[i][0].replaceAll("00[0-9]", "001");
            List<String> existing = new ArrayList<>(Arrays.asList(cases[i]).subList(1, cases[i].length));
            String result = (String) method.invoke(factory, template, existing);
            out.print("    {\"template\": " + json(template) + ", \"existing\": " + jsonArray(existing)
                    + ", \"result\": " + json(result) + "}");
            out.println(i == cases.length - 1 ? "" : ",");
        }
    }

    private static void dumpContainer(PrintStream out, String rel, File file) throws Exception {
        out.println("    {");
        out.println("      \"path\": " + json(rel) + ",");

        DSSDocument document = new FileDocument(file);

        out.println("      \"isZip\": " + call(() -> String.valueOf(ASiCUtils.isZip(document))) + ",");
        out.println("      \"isASiC\": " + call(() -> String.valueOf(ASiCUtils.isASiC(document))) + ",");
        out.println("      \"isContainerOpenDocument\": "
                + call(() -> String.valueOf(ASiCUtils.isContainerOpenDocument(document))) + ",");
        out.println("      \"containerType\": " + call(() -> {
            ASiCContainerType type = ASiCUtils.getContainerType(document);
            return type == null ? "null" : json(type.name());
        }) + ",");
        out.println("      \"zipComment\": " + call(() -> json(ASiCUtils.getZipComment(document))) + ",");
        out.println("      \"entryNames\": " + call(() -> {
            List<String> names = ZipUtils.getInstance().extractEntryNames(document);
            List<String> quoted = new ArrayList<>();
            for (String n : names) {
                quoted.add(json(n));
            }
            return quoted.toString().replace("[", "[").replace("]", "]");
        }) + ",");
        out.println("      \"entries\": " + call(() -> {
            List<DSSDocument> documents = ZipUtils.getInstance().extractContainerContent(document);
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < documents.size(); i++) {
                if (i > 0) {
                    sb.append(", ");
                }
                sb.append(entryJson(documents.get(i)));
            }
            return sb.append("]").toString();
        }));
        out.print("    }");
    }

    private interface Producer {
        String get() throws Exception;
    }

    /** Returns the produced JSON fragment, or {"error": "Class: message"} when it throws. */
    private static String call(Producer producer) {
        try {
            return producer.get();
        } catch (Throwable t) {
            return "{\"error\": " + json(t.getClass().getSimpleName() + ": " + t.getMessage()) + "}";
        }
    }

    private static String entryJson(DSSDocument document) throws Exception {
        StringBuilder sb = new StringBuilder("{");
        sb.append("\"name\": ").append(json(document.getName()));
        sb.append(", \"mimeType\": ").append(json(document.getMimeType() == null ? null : document.getMimeType().getMimeTypeString()));
        sb.append(", \"contentSha256\": ").append(json(sha256(document)));
        sb.append(", \"class\": ").append(json(document.getClass().getSimpleName()));
        if (document instanceof DSSZipEntryDocument) {
            DSSZipEntry entry = ((DSSZipEntryDocument) document).getZipEntry();
            sb.append(", \"entryName\": ").append(json(entry.getName()));
            sb.append(", \"comment\": ").append(json(entry.getComment()));
            sb.append(", \"compressionMethod\": ").append(entry.getCompressionMethod());
            sb.append(", \"size\": ").append(entry.getSize());
            sb.append(", \"compressedSize\": ").append(entry.getCompressedSize());
            sb.append(", \"crc\": ").append(entry.getCrc());
            sb.append(", \"extra\": ").append(json(hex(entry.getExtra())));
            sb.append(", \"creationTime\": ").append(entry.getCreationTime() == null ? "null" : entry.getCreationTime().toMillis());
            sb.append(", \"modificationTime\": ").append(dateMillis(entry.getModificationTime()));
            sb.append(", \"lastAccessTime\": ").append(dateMillis(entry.getLastAccessTime()));
        }
        return sb.append("}").toString();
    }

    private static String dateMillis(Date date) {
        return date == null ? "null" : Long.toString(date.getTime());
    }

    private static String sha256(DSSDocument document) throws Exception {
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        try (InputStream is = document.openStream()) {
            byte[] buf = new byte[8192];
            int n;
            while ((n = is.read(buf)) != -1) {
                md.update(buf, 0, n);
            }
        }
        return hex(md.digest());
    }

    private static String hex(byte[] data) {
        if (data == null) {
            return null;
        }
        StringBuilder sb = new StringBuilder(data.length * 2);
        for (byte b : data) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    private static String jsonArray(List<String> values) {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < values.size(); i++) {
            if (i > 0) {
                sb.append(", ");
            }
            sb.append(json(values.get(i)));
        }
        return sb.append("]").toString();
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
        return sb.append('"').toString();
    }

}

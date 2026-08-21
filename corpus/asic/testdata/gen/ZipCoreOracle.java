// Java oracle for the Go dss/asic ZIPCORE port.
//
// Dumps, for every upstream ASiC fixture container, exactly what
// SecureContainerHandler sees through java.util.zip:
//   - the sequential ZipInputStream view (local file headers), which is what
//     SecureContainerHandler.extractZipEntries() / getCurrentEntryDocument() use;
//   - the ZipFile central-directory view, which is what extractComments() and
//     FileArchiveEntry use;
//   - the ASiCUtils.getZipComment(DSSDocument) byte scan (replicated verbatim,
//     including its signed-byte commentLen quirk).
//
// Pure JDK - no DSS classes - so it can run without building the upstream project.
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.attribute.FileTime;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.List;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;
import java.util.zip.ZipInputStream;

public final class ZipCoreOracle {

    /** ASiCUtils.MAGIC_DIR */
    private static final byte[] MAGIC_DIR = { 0x50, 0x4b, 0x05, 0x06 };

    /** ASiCUtils.MAX_TO_READ */
    private static final int MAX_TO_READ = 0xFFFF + 2 + MAGIC_DIR.length;

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        List<String> relPaths = Files.readAllLines(Paths.get(args[1]));
        try (PrintStream out = new PrintStream(new java.io.FileOutputStream(args[2]), false, "UTF-8")) {
            out.println("{");
            out.println("  \"containers\": [");
            boolean firstContainer = true;
            for (String rel : relPaths) {
                if (rel.trim().isEmpty()) {
                    continue;
                }
                File file = root.resolve(rel).toFile();
                if (!firstContainer) {
                    out.println(",");
                }
                firstContainer = false;
                dumpContainer(out, rel, file);
            }
            out.println();
            out.println("  ]");
            out.println("}");
        }
    }

    private static void dumpContainer(PrintStream out, String rel, File file) throws Exception {
        out.println("    {");
        out.println("      \"path\": " + json(rel) + ",");
        out.println("      \"size\": " + file.length() + ",");
        out.println("      \"zipComment\": " + json(getZipComment(file)) + ",");

        // --- sequential (ZipInputStream / local file header) view
        out.println("      \"streamEntries\": [");
        String streamError = null;
        List<String> stream = new ArrayList<>();
        try (InputStream is = new FileInputStream(file); ZipInputStream zis = new ZipInputStream(is)) {
            ZipEntry entry;
            while ((entry = zis.getNextEntry()) != null) {
                ByteArrayOutputStream baos = new ByteArrayOutputStream();
                byte[] buf = new byte[8192];
                int n;
                while ((n = zis.read(buf)) != -1) {
                    baos.write(buf, 0, n);
                }
                stream.add(entryJson(entry, baos.toByteArray()));
            }
        } catch (Exception e) {
            streamError = e.getClass().getSimpleName() + ": " + e.getMessage();
        }
        for (int i = 0; i < stream.size(); i++) {
            out.print("        " + stream.get(i));
            out.println(i == stream.size() - 1 ? "" : ",");
        }
        out.println("      ],");
        out.println("      \"streamError\": " + json(streamError) + ",");

        // --- central-directory (ZipFile) view
        out.println("      \"centralEntries\": [");
        String centralError = null;
        List<String> central = new ArrayList<>();
        try (ZipFile zipFile = new ZipFile(file)) {
            Enumeration<? extends ZipEntry> en = zipFile.entries();
            while (en.hasMoreElements()) {
                ZipEntry entry = en.nextElement();
                byte[] content;
                if (entry.isDirectory()) {
                    content = new byte[0];
                } else {
                    ByteArrayOutputStream baos = new ByteArrayOutputStream();
                    try (InputStream is = zipFile.getInputStream(entry)) {
                        byte[] buf = new byte[8192];
                        int n;
                        while ((n = is.read(buf)) != -1) {
                            baos.write(buf, 0, n);
                        }
                    }
                    content = baos.toByteArray();
                }
                central.add(entryJson(entry, content));
            }
        } catch (Exception e) {
            centralError = e.getClass().getSimpleName() + ": " + e.getMessage();
        }
        for (int i = 0; i < central.size(); i++) {
            out.print("        " + central.get(i));
            out.println(i == central.size() - 1 ? "" : ",");
        }
        out.println("      ],");
        out.println("      \"centralError\": " + json(centralError));
        out.print("    }");
    }

    private static String entryJson(ZipEntry e, byte[] content) throws Exception {
        StringBuilder sb = new StringBuilder();
        sb.append("{");
        sb.append("\"name\": ").append(json(e.getName()));
        sb.append(", \"method\": ").append(e.getMethod());
        sb.append(", \"size\": ").append(e.getSize());
        sb.append(", \"compressedSize\": ").append(e.getCompressedSize());
        sb.append(", \"crc\": ").append(e.getCrc());
        sb.append(", \"time\": ").append(e.getTime());
        sb.append(", \"comment\": ").append(json(e.getComment()));
        sb.append(", \"extra\": ").append(json(hex(e.getExtra())));
        sb.append(", \"creationTime\": ").append(fileTime(e.getCreationTime()));
        sb.append(", \"lastModifiedTime\": ").append(fileTime(e.getLastModifiedTime()));
        sb.append(", \"lastAccessTime\": ").append(fileTime(e.getLastAccessTime()));
        sb.append(", \"isDirectory\": ").append(e.isDirectory());
        sb.append(", \"contentLength\": ").append(content.length);
        sb.append(", \"contentSha256\": ").append(json(sha256(content)));
        sb.append("}");
        return sb.toString();
    }

    private static String fileTime(FileTime t) {
        return t == null ? "null" : Long.toString(t.toMillis());
    }

    private static String sha256(byte[] data) throws Exception {
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        return hex(md.digest(data));
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

    /** Verbatim replica of ASiCUtils.getZipComment(DSSDocument) over a FileDocument. */
    private static String getZipComment(File file) throws IOException {
        long fileLength = file.length();
        try (InputStream is = new FileInputStream(file)) {
            if (fileLength > MAX_TO_READ) {
                long toSkip = fileLength - MAX_TO_READ;
                long skipped = is.skip(toSkip);
                if (skipped != toSkip) {
                    throw new IOException("Different amount of bytes have been skipped!");
                }
            }
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            byte[] tmp = new byte[8192];
            int n;
            while ((n = is.read(tmp)) != -1) {
                baos.write(tmp, 0, n);
            }
            byte[] buffer = baos.toByteArray();
            if (buffer.length == 0) {
                return null;
            }
            final int len = buffer.length;
            for (int ii = len - 22; ii >= 0; ii--) {
                boolean isMagicStart = true;
                for (int jj = 0; jj < MAGIC_DIR.length; jj++) {
                    if (buffer[ii + jj] != MAGIC_DIR[jj]) {
                        isMagicStart = false;
                        break;
                    }
                }
                if (isMagicStart) {
                    int realLen = len - ii - 22;
                    if (realLen == 0) {
                        return null;
                    }
                    return new String(buffer, ii + 22, realLen, StandardCharsets.UTF_8);
                }
            }
        }
        return null;
    }

}

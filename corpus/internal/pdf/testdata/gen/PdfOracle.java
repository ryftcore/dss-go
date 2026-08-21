// PdfOracle — the pdfbox 3.0.7 oracle for internal/pdf (DESIGN.md §6.2).
//
// Run BY HAND, never from `go test`. The goldens it writes are checked in; a
// golden that changes in a PR is a red flag and must be justified there.
//
// Proven command (pdfbox 3.0.7 + its runtime deps, fetched once with
// `mvn dependency:build-classpath` from scratchpad/pdfprobe-pom.xml):
//
//   M2=$HOME/.m2/repository
//   CP=$M2/org/apache/pdfbox/pdfbox/3.0.7/pdfbox-3.0.7.jar\
//   :$M2/org/apache/pdfbox/pdfbox-io/3.0.7/pdfbox-io-3.0.7.jar\
//   :$M2/org/apache/pdfbox/fontbox/3.0.7/fontbox-3.0.7.jar\
//   :$M2/commons-logging/commons-logging/1.3.5/commons-logging-1.3.5.jar
//   javac -nowarn -cp "$CP" -d /tmp/pdforacle PdfOracle.java
//   java -Dorg.apache.commons.logging.Log=org.apache.commons.logging.impl.NoOpLog \
//        -cp "$CP:/tmp/pdforacle" PdfOracle \
//        <corpus dir> <golden file> <manifest file> [file-list]
//
// Output: one tab-separated record per PDF, fields as `key=value`, stable order.
// A Java-side failure is recorded as `error=!ERROR <SimpleClassName>` and the Go
// test asserts a corresponding Go-side failure — a Go success where Java failed
// is a test failure.

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import org.apache.pdfbox.Loader;
import org.apache.pdfbox.cos.COSArray;
import org.apache.pdfbox.cos.COSBase;
import org.apache.pdfbox.cos.COSDictionary;
import org.apache.pdfbox.cos.COSName;
import org.apache.pdfbox.cos.COSObjectKey;
import org.apache.pdfbox.cos.COSString;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.encryption.AccessPermission;
import org.apache.pdfbox.pdmodel.encryption.PDEncryption;
import org.apache.pdfbox.pdmodel.interactive.digitalsignature.PDSignature;
import org.apache.pdfbox.pdmodel.interactive.form.PDAcroForm;
import org.apache.pdfbox.pdmodel.interactive.form.PDField;
import org.apache.pdfbox.pdmodel.interactive.form.PDSignatureField;

public final class PdfOracle {

    private static final byte[] PDF_EOF = "%%EOF".getBytes(StandardCharsets.ISO_8859_1);

    public static void main(String[] args) throws Exception {
        if (args.length < 3) {
            System.err.println("usage: PdfOracle <corpus dir> <golden> <manifest> [list file]");
            System.exit(2);
        }
        File corpus = new File(args[0]);
        List<String> rel = new ArrayList<>();
        if (args.length >= 4) {
            for (String line : Files.readAllLines(new File(args[3]).toPath())) {
                line = line.trim();
                if (!line.isEmpty() && !line.startsWith("#")) {
                    rel.add(line);
                }
            }
        } else {
            collect(corpus, "", rel);
            rel.sort(Comparator.naturalOrder());
        }

        try (PrintWriter out = new PrintWriter(args[1], "UTF-8");
             PrintWriter man = new PrintWriter(args[2], "UTF-8")) {
            for (String path : rel) {
                File f = new File(corpus, path);
                byte[] raw = Files.readAllBytes(f.toPath());
                man.println(sha256Hex(raw) + "  " + path);
                out.println(record(path, f, raw));
            }
            man.flush();
            out.flush();
        }
    }

    private static void collect(File dir, String prefix, List<String> out) {
        File[] kids = dir.listFiles();
        if (kids == null) {
            return;
        }
        Arrays.sort(kids, Comparator.comparing(File::getName));
        for (File k : kids) {
            if (k.isDirectory()) {
                collect(k, prefix + k.getName() + "/", out);
            } else if (k.getName().toLowerCase().endsWith(".pdf")) {
                out.add(prefix + k.getName());
            }
        }
    }

    private static String record(String path, File f, byte[] raw) {
        StringBuilder sb = new StringBuilder();
        sb.append("path=").append(path);
        sb.append("\tsize=").append(raw.length);
        // The %%EOF scan is independent of pdfbox: it is a transcription of
        // eu.europa.esig.dss.pades.PAdESUtils.extractRevisions.
        List<Long> revs = extractRevisions(raw);
        sb.append("\teofRevisions=").append(revs.size());
        sb.append("\trevisionEnds=").append(join(revs));
        sb.append("\txrefChain=").append(xrefChain(raw));

        String password = null;
        PDDocument doc = null;
        try {
            doc = Loader.loadPDF(f);
            password = "";
        } catch (Exception first) {
            try {
                doc = Loader.loadPDF(f, " ");
                password = " ";
            } catch (Exception second) {
                sb.append("\terror=!ERROR ").append(first.getClass().getSimpleName());
                sb.append("\terrorMessage=").append(clean(String.valueOf(first.getMessage())));
                return sb.toString();
            }
        }
        try {
            sb.append("\tpassword=").append(password.isEmpty() ? "-" : "SP");
            sb.append("\theaderVersion=").append(fmt(doc.getDocument().getVersion()));
            sb.append("\tcatalogVersion=").append(fmt(doc.getVersion()));

            COSDictionary trailer = doc.getDocument().getTrailer();
            Set<String> trailerKeys = new TreeSet<>();
            for (COSName n : trailer.keySet()) {
                trailerKeys.add(n.getName());
            }
            sb.append("\ttrailer=").append(String.join(",", trailerKeys));

            Map<COSObjectKey, Long> xref = doc.getDocument().getXrefTable();
            List<COSObjectKey> keys = new ArrayList<>(xref.keySet());
            keys.sort((a, b) -> {
                int c = Long.compare(a.getNumber(), b.getNumber());
                return c != 0 ? c : Integer.compare(a.getGeneration(), b.getGeneration());
            });
            sb.append("\tobjects=").append(keys.size());
            StringBuilder ks = new StringBuilder();
            for (COSObjectKey k : keys) {
                if (ks.length() > 0) {
                    ks.append(' ');
                }
                ks.append(k.getNumber()).append(':').append(k.getGeneration());
            }
            sb.append("\tobjectKeys=").append(ks);

            // Stream inventory: for every object in the snapshot above, the
            // SHA-256 of its raw (still encoded) bytes and of its decoded bytes.
            // This is what exercises the filter chain, the predictors, object
            // streams and — on the protected/*.pdf files — decryption.
            StringBuilder sdig = new StringBuilder();
            int streamCount = 0;
            for (COSObjectKey k : keys) {
                COSBase base;
                try {
                    base = doc.getDocument().getObjectFromPool(k).getObject();
                } catch (Exception e) {
                    continue;
                }
                if (!(base instanceof org.apache.pdfbox.cos.COSStream)) {
                    continue;
                }
                org.apache.pdfbox.cos.COSStream st = (org.apache.pdfbox.cos.COSStream) base;
                streamCount++;
                sdig.append(k.getNumber()).append(':').append(k.getGeneration());
                try (InputStream is = st.createRawInputStream()) {
                    sdig.append(':').append(sha256Hex(readAll(is)));
                } catch (Exception e) {
                    sdig.append(":!raw");
                }
                if (!decodable(st)) {
                    // DESIGN.md §2.5: image filters are never decoded, so the
                    // oracle does not decode them either.
                    sdig.append(":!unsupported");
                } else {
                    try (InputStream is = st.createInputStream()) {
                        sdig.append(':').append(sha256Hex(readAll(is)));
                    } catch (Exception e) {
                        sdig.append(":!dec");
                    }
                }
                sdig.append(' ');
            }
            sb.append("\tstreams=").append(streamCount);
            sb.append("\tstreamDigest=").append(sha256Hex(sdig.toString().getBytes(StandardCharsets.ISO_8859_1)));

            sb.append("\txrefType=").append(doc.getDocument().isXRefStream() ? "stream" : "table");
            sb.append("\thybrid=").append(doc.getDocument().hasHybridXRef());

            sb.append("\tpages=").append(doc.getNumberOfPages());
            StringBuilder boxes = new StringBuilder();
            StringBuilder rots = new StringBuilder();
            for (int i = 0; i < doc.getNumberOfPages(); i++) {
                PDPage p = doc.getPage(i);
                if (i > 0) {
                    boxes.append(' ');
                    rots.append(' ');
                }
                boxes.append(fmt(p.getMediaBox().getLowerLeftX())).append(',')
                     .append(fmt(p.getMediaBox().getLowerLeftY())).append(',')
                     .append(fmt(p.getMediaBox().getUpperRightX())).append(',')
                     .append(fmt(p.getMediaBox().getUpperRightY()));
                rots.append(p.getRotation());
            }
            sb.append("\tpageBoxes=").append(boxes);
            sb.append("\tpageRotations=").append(rots);

            sb.append("\tencrypted=").append(encryption(doc));
            AccessPermission ap = doc.getCurrentAccessPermission();
            sb.append("\tpermissions=")
              .append("modify=").append(ap.canModify())
              .append(",annots=").append(ap.canModifyAnnotations())
              .append(",fill=").append(ap.canFillInForm())
              .append(",owner=").append(ap.isOwnerPermission());

            sb.append("\tsigs=").append(signatures(doc, raw.length));
            sb.append("\tsigFields=").append(signatureFieldGeometry(doc));
            sb.append("\tannots=").append(annotations(doc));
            sb.append("\tdss=").append(dssDictionary(doc));
            sb.append("\tvri=").append(vriDictionary(doc));
            sb.append("\tacroForm=").append(acroForm(doc));
            sb.append("\terror=");
            return sb.toString();
        } catch (Throwable t) {
            sb.append("\terror=!ERROR ").append(t.getClass().getSimpleName());
            sb.append("\terrorMessage=").append(clean(String.valueOf(t.getMessage())));
            return sb.toString();
        } finally {
            try {
                doc.close();
            } catch (IOException ignored) {
                // nothing to do
            }
        }
    }

    // --- DESIGN.md §6.2 fields: annots, dss, vri, sigFields -----------------
    //
    // These four cover the call sites in §0.1 items 5 and 6 — the /DSS and /VRI
    // dictionaries that PAdES LT/LTA validation reads, and the page annotations
    // that DefaultPdfDifferencesFinder compares. They are read through raw COS
    // traversal rather than pdfbox's PD* wrappers, because that is the level the
    // Go reader exposes (Document.GetDict / GetArray / IndexRef / StreamData)
    // and the level DSS's own PdfDict / PdfArray SPI sits at.
    //
    // /VRI key order is deliberately taken from COSDictionary.keySet(), which is
    // a LinkedHashMap: SingleDssDict.extractVRIs enumerates in that order, so it
    // is the insertion-order contract of DESIGN.md §2.2 under test.

    private static String annotations(PDDocument doc) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < doc.getNumberOfPages(); i++) {
            if (i > 0) {
                sb.append(' ');
            }
            COSBase ab = doc.getPage(i).getCOSObject().getDictionaryObject(COSName.ANNOTS);
            if (!(ab instanceof COSArray)) {
                sb.append('-');
                continue;
            }
            COSArray annots = (COSArray) ab;
            StringBuilder per = new StringBuilder();
            for (int j = 0; j < annots.size(); j++) {
                COSBase e = annots.getObject(j);
                if (!(e instanceof COSDictionary)) {
                    continue;
                }
                COSDictionary ad = (COSDictionary) e;
                if (per.length() > 0) {
                    per.append('|');
                }
                long key = -1;
                COSBase rawItem = annots.get(j);
                if (rawItem instanceof org.apache.pdfbox.cos.COSObject) {
                    key = ((org.apache.pdfbox.cos.COSObject) rawItem).getKey().getNumber();
                }
                per.append(key).append(';').append(rect(ad.getDictionaryObject(COSName.RECT)));
                COSBase t = ad.getDictionaryObject(COSName.T);
                per.append(';').append(t instanceof COSString ? clean(((COSString) t).getString()) : "");
                per.append(';').append(ad.getDictionaryObject(COSName.V) != null);
                int f = ad.getInt(COSName.F, 0);
                per.append(';').append((f & 2) != 0).append(';').append((f & 16) != 0);
            }
            sb.append(per.length() == 0 ? "-" : per.toString());
        }
        return sb.toString();
    }

    private static String rect(COSBase b) {
        if (!(b instanceof COSArray)) {
            return "-";
        }
        COSArray a = (COSArray) b;
        if (a.size() < 4) {
            return "-";
        }
        float[] v = new float[4];
        for (int i = 0; i < 4; i++) {
            COSBase e = a.getObject(i);
            if (!(e instanceof org.apache.pdfbox.cos.COSNumber)) {
                return "-";
            }
            v[i] = ((org.apache.pdfbox.cos.COSNumber) e).floatValue();
        }
        // Normalised so min <= max on both axes, matching pdf.RectFromArray.
        return fmt(Math.min(v[0], v[2])) + "," + fmt(Math.min(v[1], v[3])) + ","
             + fmt(Math.max(v[0], v[2])) + "," + fmt(Math.max(v[1], v[3]));
    }

    private static String signatureFieldGeometry(PDDocument doc) {
        PDAcroForm form = doc.getDocumentCatalog().getAcroForm(null);
        if (form == null) {
            return "";
        }
        StringBuilder sb = new StringBuilder();
        for (PDField field : form.getFieldTree()) {
            if (!(field instanceof PDSignatureField)) {
                continue;
            }
            PDSignatureField sf = (PDSignatureField) field;
            if (sb.length() > 0) {
                sb.append(' ');
            }
            COSDictionary fd = sf.getCOSObject();
            int widgets = 0;
            if (fd.getNameAsString(COSName.SUBTYPE) != null
                    && "Widget".equals(fd.getNameAsString(COSName.SUBTYPE))) {
                widgets = 1;
            } else {
                COSBase kids = fd.getDictionaryObject(COSName.KIDS);
                if (kids instanceof COSArray) {
                    for (COSBase k : (COSArray) kids) {
                        if (resolve(k) instanceof COSDictionary) {
                            widgets++;
                        }
                    }
                }
            }
            sb.append(clean(sf.getFullyQualifiedName()))
              .append(";widgets=").append(widgets)
              .append(";lock=").append(fd.getDictionaryObject(COSName.getPDFName("Lock")) != null);
        }
        return sb.toString();
    }

    private static COSBase resolve(COSBase b) {
        return b instanceof org.apache.pdfbox.cos.COSObject
                ? ((org.apache.pdfbox.cos.COSObject) b).getObject() : b;
    }

    private static String dssDictionary(PDDocument doc) {
        COSDictionary dss = doc.getDocumentCatalog().getCOSObject()
                .getCOSDictionary(COSName.getPDFName("DSS"));
        if (dss == null) {
            return "-";
        }
        return "Certs=" + tokenArray(dss, "Certs")
             + ",CRLs=" + tokenArray(dss, "CRLs")
             + ",OCSPs=" + tokenArray(dss, "OCSPs");
    }

    // tokenArray renders "<count>[key:sha256-of-decoded-bytes …]", which is what
    // DSSDictionaryExtractionUtils.getCertsFromArray consumes: the object key of
    // every element plus its decoded stream content.
    private static String tokenArray(COSDictionary parent, String name) {
        COSBase b = parent.getDictionaryObject(COSName.getPDFName(name));
        if (!(b instanceof COSArray)) {
            return "0[]";
        }
        COSArray arr = (COSArray) b;
        StringBuilder sb = new StringBuilder();
        int n = 0;
        for (int i = 0; i < arr.size(); i++) {
            COSBase raw = arr.get(i);
            long key = -1;
            if (raw instanceof org.apache.pdfbox.cos.COSObject) {
                key = ((org.apache.pdfbox.cos.COSObject) raw).getKey().getNumber();
            }
            COSBase e = arr.getObject(i);
            if (!(e instanceof org.apache.pdfbox.cos.COSStream)) {
                continue;
            }
            if (n > 0) {
                sb.append(' ');
            }
            n++;
            sb.append(key).append(':');
            try (InputStream is = ((org.apache.pdfbox.cos.COSStream) e).createInputStream()) {
                sb.append(sha256Hex(readAll(is)));
            } catch (Exception ex) {
                sb.append("!dec");
            }
        }
        return n + "[" + sb + "]";
    }

    private static String vriDictionary(PDDocument doc) {
        COSDictionary dss = doc.getDocumentCatalog().getCOSObject()
                .getCOSDictionary(COSName.getPDFName("DSS"));
        if (dss == null) {
            return "-";
        }
        COSDictionary vri = dss.getCOSDictionary(COSName.getPDFName("VRI"));
        if (vri == null) {
            return "-";
        }
        StringBuilder sb = new StringBuilder();
        // Insertion order: COSDictionary is a LinkedHashMap and that order is
        // what PdfDict.list() hands to SingleDssDict.extractVRIs.
        for (COSName n : vri.keySet()) {
            if (sb.length() > 0) {
                sb.append(' ');
            }
            sb.append(n.getName());
            COSDictionary entry = vri.getCOSDictionary(n);
            if (entry == null) {
                sb.append(":!notdict");
                continue;
            }
            sb.append(':').append(tokenArray(entry, "Cert"))
              .append(',').append(tokenArray(entry, "CRL"))
              .append(',').append(tokenArray(entry, "OCSP"));
            COSBase tu = entry.getDictionaryObject(COSName.getPDFName("TU"));
            sb.append(",TU=").append(tu instanceof COSString
                    ? clean(((COSString) tu).getString()) : "-");
            COSBase ts = entry.getDictionaryObject(COSName.getPDFName("TS"));
            if (ts instanceof org.apache.pdfbox.cos.COSStream) {
                try (InputStream is = ((org.apache.pdfbox.cos.COSStream) ts).createInputStream()) {
                    sb.append(",TS=").append(sha256Hex(readAll(is)));
                } catch (Exception ex) {
                    sb.append(",TS=!dec");
                }
            } else {
                sb.append(",TS=-");
            }
        }
        return sb.length() == 0 ? "" : sb.toString();
    }

    // decodable mirrors DESIGN.md §2.5: FlateDecode, LZWDecode, ASCIIHexDecode,
    // ASCII85Decode, RunLengthDecode and /Crypt /Identity, nothing else.
    private static boolean decodable(org.apache.pdfbox.cos.COSStream st) {
        COSBase f = st.getDictionaryObject(COSName.FILTER);
        List<COSName> names = new ArrayList<>();
        if (f instanceof COSName) {
            names.add((COSName) f);
        } else if (f instanceof COSArray) {
            for (COSBase e : (COSArray) f) {
                if (e instanceof COSName) {
                    names.add((COSName) e);
                }
            }
        }
        for (COSName n : names) {
            switch (n.getName()) {
                case "FlateDecode":
                case "Fl":
                case "LZWDecode":
                case "LZW":
                case "ASCIIHexDecode":
                case "AHx":
                case "ASCII85Decode":
                case "A85":
                case "RunLengthDecode":
                case "RL":
                case "Crypt":
                    break;
                default:
                    return false;
            }
        }
        return true;
    }

    private static byte[] readAll(InputStream is) throws IOException {
        ByteArrayOutputStream bos = new ByteArrayOutputStream();
        byte[] buf = new byte[8192];
        int n;
        while ((n = is.read(buf)) > 0) {
            bos.write(buf, 0, n);
        }
        return bos.toByteArray();
    }

    private static String acroForm(PDDocument doc) {
        try {
            PDAcroForm form = doc.getDocumentCatalog().getAcroForm(null);
            return form != null ? "yes" : "no";
        } catch (Exception e) {
            return "err";
        }
    }

    private static String encryption(PDDocument doc) {
        if (!doc.isEncrypted()) {
            return "-";
        }
        PDEncryption enc = doc.getEncryption();
        StringBuilder sb = new StringBuilder();
        sb.append("Filter=").append(enc.getFilter());
        sb.append(",V=").append(enc.getVersion());
        sb.append(",R=").append(enc.getRevision());
        sb.append(",Length=").append(enc.getLength());
        String cfm = "-";
        try {
            COSDictionary cf = enc.getCOSObject().getCOSDictionary(COSName.CF);
            if (cf != null) {
                StringBuilder m = new StringBuilder();
                for (COSName n : cf.keySet()) {
                    COSDictionary sub = cf.getCOSDictionary(n);
                    if (sub != null) {
                        m.append(n.getName()).append('=')
                         .append(sub.getNameAsString(COSName.CFM)).append(';');
                    }
                }
                cfm = m.toString();
            }
        } catch (Exception e) {
            cfm = "err";
        }
        sb.append(",CFM=").append(cfm);
        sb.append(",StmF=").append(enc.getStreamFilterName().getName());
        sb.append(",StrF=").append(enc.getStringFilterName().getName());
        return sb.toString();
    }

    private static String signatures(PDDocument doc, long fileLength) throws IOException {
        StringBuilder sb = new StringBuilder();
        PDAcroForm form = doc.getDocumentCatalog().getAcroForm(null);
        if (form == null) {
            return "";
        }
        for (PDField field : form.getFieldTree()) {
            if (!(field instanceof PDSignatureField)) {
                continue;
            }
            PDSignatureField sf = (PDSignatureField) field;
            if (sb.length() > 0) {
                sb.append(' ');
            }
            sb.append(clean(sf.getFullyQualifiedName())).append(':');
            PDSignature sig = sf.getSignature();
            COSBase v = sf.getCOSObject().getItem(COSName.V);
            long valueKey = -1;
            if (v instanceof org.apache.pdfbox.cos.COSObject) {
                valueKey = ((org.apache.pdfbox.cos.COSObject) v).getKey().getNumber();
            }
            if (sig == null || sig.getCOSObject() == null || sig.getCOSObject().size() == 0) {
                sb.append("EMPTY");
                continue;
            }
            COSDictionary sd = sig.getCOSObject();
            sb.append("valueKey=").append(valueKey);
            sb.append(",Type=").append(nz(sd.getNameAsString(COSName.TYPE)));
            sb.append(",Filter=").append(nz(sig.getFilter()));
            sb.append(",SubFilter=").append(nz(sig.getSubFilter()));
            int[] br = sig.getByteRange();
            sb.append(",BR=[");
            for (int i = 0; i < br.length; i++) {
                if (i > 0) {
                    sb.append(' ');
                }
                sb.append(br[i]);
            }
            sb.append(']');
            byte[] contents = null;
            COSBase c = sd.getDictionaryObject(COSName.CONTENTS);
            if (c instanceof COSString) {
                contents = ((COSString) c).getBytes();
            }
            sb.append(",ContentsLen=").append(contents == null ? -1 : contents.length);
            sb.append(",ContentsSHA256=").append(contents == null ? "-" : sha256Hex(contents));
            sb.append(",covers=").append(coversWholeDocument(br, fileLength));
        }
        return sb.toString();
    }

    // Transcription of PdfBoxDocumentReader.isSignatureCoversWholeDocument. The
    // arithmetic is deliberately the upstream one, not the obvious one.
    private static boolean coversWholeDocument(int[] br, long fileLength) {
        if (br == null || br.length < 4) {
            return false;
        }
        long before = (long) br[1] - br[0];
        long expectedCMS = (long) br[2] - br[1] - br[0];
        long after = br[3];
        return fileLength == before + expectedCMS + after;
    }

    // Transcription of PAdESUtils.extractRevisions (DSS 6.5.RC1).
    private static List<Long> extractRevisions(byte[] raw) {
        List<Long> out = new ArrayList<>();
        int position = 0;
        ByteArrayOutputStream tempLine = new ByteArrayOutputStream();
        int i = 0;
        while (i < raw.length) {
            int b = raw[i++] & 0xFF;
            ++position;
            tempLine.write(b);
            byte[] cur = tempLine.toByteArray();
            if (Arrays.equals(PDF_EOF, cur)) {
                tempLine = new ByteArrayOutputStream();
                int eofPosition = position;
                int c = i < raw.length ? (raw[i++] & 0xFF) : -1;
                if (c != -1) {
                    ++position;
                }
                if (c == '\n') {
                    ++eofPosition;
                } else if (c == '\r') {
                    ++eofPosition;
                    int d = i < raw.length ? (raw[i++] & 0xFF) : -1;
                    if (d != -1) {
                        ++position;
                    }
                    if (d == '\n') {
                        ++eofPosition;
                    }
                }
                out.add((long) eofPosition);
            } else if (b == '\n' || b == '\r' || cur.length > PDF_EOF.length) {
                tempLine = new ByteArrayOutputStream();
            }
        }
        return out;
    }

    // A naive /Prev walk over the raw bytes: `table`, `stream` or `?` per hop.
    // It applies no repair, so a `?` marks a hop whose offset is broken and the
    // Go side is allowed to differ there (it reproduces pdfbox's X3 repair).
    private static String xrefChain(byte[] raw) {
        StringBuilder sb = new StringBuilder();
        long off = startXref(raw);
        java.util.Set<Long> seen = new java.util.HashSet<>();
        while (off > 0 && off < raw.length && seen.add(off) && seen.size() < 600) {
            int p = (int) off;
            while (p < raw.length && isWhite(raw[p])) {
                p++;
            }
            if (sb.length() > 0) {
                sb.append(',');
            }
            if (startsWith(raw, p, "xref")) {
                sb.append("table");
                int t = indexOf(raw, "trailer", p);
                if (t < 0) {
                    break;
                }
                off = dictLong(raw, t, "/Prev");
            } else {
                int s = indexOf(raw, "stream", p);
                if (s < 0 || s - p > 4096) {
                    sb.append('?');
                    break;
                }
                sb.append("stream");
                off = dictLong(raw, p, "/Prev");
            }
        }
        return sb.toString();
    }

    private static long startXref(byte[] raw) {
        int from = Math.max(0, raw.length - 2048);
        int idx = lastIndexOf(raw, "startxref", raw.length);
        if (idx < from) {
            return -1;
        }
        int p = idx + "startxref".length();
        while (p < raw.length && isWhite(raw[p])) {
            p++;
        }
        long v = 0;
        boolean any = false;
        while (p < raw.length && raw[p] >= '0' && raw[p] <= '9') {
            v = v * 10 + (raw[p] - '0');
            p++;
            any = true;
        }
        return any ? v : -1;
    }

    private static long dictLong(byte[] raw, int from, String key) {
        int end = Math.min(raw.length, from + 4096);
        int idx = indexOf(raw, key, from);
        if (idx < 0 || idx > end) {
            return -1;
        }
        int p = idx + key.length();
        while (p < raw.length && isWhite(raw[p])) {
            p++;
        }
        long v = 0;
        boolean any = false;
        while (p < raw.length && raw[p] >= '0' && raw[p] <= '9') {
            v = v * 10 + (raw[p] - '0');
            p++;
            any = true;
        }
        return any ? v : -1;
    }

    private static boolean isWhite(byte b) {
        return b == ' ' || b == '\r' || b == '\n' || b == '\t' || b == 0 || b == 12;
    }

    private static boolean startsWith(byte[] raw, int at, String s) {
        if (at < 0 || at + s.length() > raw.length) {
            return false;
        }
        for (int i = 0; i < s.length(); i++) {
            if (raw[at + i] != s.charAt(i)) {
                return false;
            }
        }
        return true;
    }

    private static int indexOf(byte[] raw, String s, int from) {
        for (int i = Math.max(0, from); i + s.length() <= raw.length; i++) {
            if (startsWith(raw, i, s)) {
                return i;
            }
        }
        return -1;
    }

    private static int lastIndexOf(byte[] raw, String s, int before) {
        for (int i = Math.min(before, raw.length - s.length()); i >= 0; i--) {
            if (startsWith(raw, i, s)) {
                return i;
            }
        }
        return -1;
    }

    private static String join(List<Long> v) {
        StringBuilder sb = new StringBuilder();
        for (Long l : v) {
            if (sb.length() > 0) {
                sb.append(' ');
            }
            sb.append(l);
        }
        return sb.toString();
    }

    private static String fmt(float f) {
        if (f == Math.rint(f) && !Float.isInfinite(f)) {
            return String.valueOf((long) f) + ".0";
        }
        return String.valueOf(f);
    }

    // nz renders an absent name as the empty string; pdfbox's getters return
    // null there and String.valueOf would print the literal "null".
    private static String nz(String s) {
        return s == null ? "" : s;
    }

    private static String clean(String s) {
        if (s == null) {
            return "";
        }
        return s.replace('\t', ' ').replace('\n', ' ').replace('\r', ' ');
    }

    private static String sha256Hex(byte[] b) {
        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            byte[] h = md.digest(b);
            StringBuilder sb = new StringBuilder();
            for (byte x : h) {
                sb.append(Character.forDigit((x >> 4) & 0xF, 16));
                sb.append(Character.forDigit(x & 0xF, 16));
            }
            return sb.toString();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private PdfOracle() {
        // no instances
    }
}

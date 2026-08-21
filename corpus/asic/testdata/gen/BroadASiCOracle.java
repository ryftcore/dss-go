// Java oracle #4 for the Go dss/asic port: the BROAD-CORPUS differential oracle.
//
// Where ZipCoreDssOracle.java pins the ZIP layer (entry inventory + per-entry metadata) and
// MergeOracle.java pins the merger, this oracle pins the two surfaces between them - the layers
// that turn a bag of zip entries into an ASiC container model - over EVERY fixture in the corpus,
// through BOTH per-format stacks:
//
//   1. Per-format container extraction. ASiCWithCAdESContainerExtractor and
//      ASiCWithXAdESContainerExtractor are run over every fixture and the full ASiCContent
//      bucketing is dumped: container type, zip comment, mimetype document, and the name of every
//      document in each of the ten buckets (signature / manifest / archiveManifest /
//      evidenceRecordManifest / timestamp / evidenceRecord / signed / unsupported / folders /
//      containerDocuments), plus the two derived views (rootLevelSignedDocuments,
//      allManifestDocuments). A single entry landing in the wrong bucket silently changes what a
//      signature is taken to cover, so every bucket is compared by name, in order.
//
//   2. Manifest parsing, both flavors. Every manifest/archive-manifest/evidence-record-manifest
//      document the CAdES extractor found is run through ASiCManifestParser.getManifestFile, and
//      every ODF-style manifest the XAdES extractor found through ASiCEWithXAdESManifestParser;
//      the resulting ManifestFile (signature filename, manifest type) and each ManifestEntry (uri,
//      document name, mime type, digest algorithm + hex digest, rootfile flag) is dumped.
//      ASiCManifestParser.getLinkedManifest is exercised for every signature name in the
//      container as well.
//
// Failure paths are part of the contract, so an exception on any of the above is recorded as its
// class name and message rather than skipping the fixture: the Go port has to fail the same way.
//
// Usage: java BroadASiCOracle <upstream-root> <fixtures.txt> <out.json>
import eu.europa.esig.dss.asic.cades.extract.ASiCWithCAdESContainerExtractor;
import eu.europa.esig.dss.asic.common.ASiCContent;
import eu.europa.esig.dss.asic.common.validation.ASiCManifestParser;
import eu.europa.esig.dss.asic.xades.extract.ASiCWithXAdESContainerExtractor;
import eu.europa.esig.dss.asic.xades.validation.ASiCContainerWithXAdESAnalyzer;
import eu.europa.esig.dss.asic.xades.validation.ASiCEWithXAdESManifestParser;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.ManifestEntry;
import eu.europa.esig.dss.model.ManifestFile;
import eu.europa.esig.dss.spi.signature.AdvancedSignature;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;

import java.io.File;
import java.io.FileOutputStream;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public final class BroadASiCOracle {

    public static void main(String[] args) throws Exception {
        Path root = Paths.get(args[0]);
        List<String> relPaths = Files.readAllLines(Paths.get(args[1]));
        try (PrintStream out = new PrintStream(new FileOutputStream(args[2]), false, "UTF-8")) {
            out.println("[");
            boolean first = true;
            for (String rel : relPaths) {
                if (rel.trim().isEmpty()) {
                    continue;
                }
                if (!first) {
                    out.println(",");
                }
                first = false;
                dumpFixture(out, rel, root.resolve(rel).toFile());
            }
            out.println();
            out.println("]");
        }
    }

    private static void dumpFixture(PrintStream out, String rel, File file) {
        out.println("{");
        out.println("  \"path\": " + str(rel) + ",");
        DSSDocument document = new FileDocument(file);

        out.println("  \"cades\": ");
        ASiCContent cadesContent = dumpExtraction(out, () ->
                new ASiCWithCAdESContainerExtractor(document).extract());
        out.println(",");
        out.println("  \"xades\": ");
        ASiCContent xadesContent = dumpExtraction(out, () ->
                new ASiCWithXAdESContainerExtractor(document).extract());
        out.println(",");

        out.println("  \"cadesManifests\": [");
        if (cadesContent != null) {
            List<DSSDocument> manifests = new ArrayList<>();
            manifests.addAll(cadesContent.getManifestDocuments());
            manifests.addAll(cadesContent.getArchiveManifestDocuments());
            manifests.addAll(cadesContent.getEvidenceRecordManifestDocuments());
            boolean firstManifest = true;
            for (DSSDocument manifest : manifests) {
                if (!firstManifest) {
                    out.println(",");
                }
                firstManifest = false;
                dumpManifest(out, manifest.getName(), () -> ASiCManifestParser.getManifestFile(manifest));
            }
            out.println();
        }
        out.println("  ],");

        out.println("  \"xadesManifests\": [");
        if (xadesContent != null) {
            boolean firstManifest = true;
            for (DSSDocument manifest : xadesContent.getManifestDocuments()) {
                if (!firstManifest) {
                    out.println(",");
                }
                firstManifest = false;
                dumpManifest(out, manifest.getName(),
                        () -> new ASiCEWithXAdESManifestParser(manifest).getManifest());
            }
            out.println();
        }
        out.println("  ],");

        // getLinkedManifest is what binds a signature to its ASiCManifest; probe it with every
        // signature name the CAdES extractor found, so an off-by-one in the linking rule shows up.
        out.println("  \"linkedManifests\": [");
        if (cadesContent != null) {
            List<DSSDocument> manifests = new ArrayList<>(cadesContent.getManifestDocuments());
            manifests.addAll(cadesContent.getArchiveManifestDocuments());
            boolean firstLink = true;
            for (DSSDocument signature : cadesContent.getSignatureDocuments()) {
                if (!firstLink) {
                    out.println(",");
                }
                firstLink = false;
                String linked;
                try {
                    DSSDocument result = ASiCManifestParser.getLinkedManifest(manifests, signature.getName());
                    linked = result == null ? null : result.getName();
                } catch (Exception e) {
                    linked = "!" + e.getClass().getName();
                }
                out.print("    {\"signature\": " + str(signature.getName())
                        + ", \"manifest\": " + str(linked) + "}");
            }
            out.println();
        }
        out.println("  ],");

        out.println("  \"xadesAnalysis\": ");
        dumpXAdESAnalysis(out, document);
        out.println();
        out.print("}");
    }

    // Per-signature analysis, XAdES flavor. ASiCContainerWithXAdESAnalyzer is the layer above
    // extraction: it turns the extracted entries into AdvancedSignature objects, resolves what each
    // signature actually covers, and describes the container's manifests. It needs no part of the
    // Phase-8 validation engine on either side, so the whole surface is comparable today.
    private static void dumpXAdESAnalysis(PrintStream out, DSSDocument document) {
        StringBuilder sb = new StringBuilder();
        try {
            ASiCContainerWithXAdESAnalyzer analyzer = new ASiCContainerWithXAdESAnalyzer(document);
            if (!analyzer.isSupported(document)) {
                out.print("  {\"error\": null, \"supported\": false}");
                return;
            }
            analyzer.setCertificateVerifier(new CommonCertificateVerifier());
            sb.append("  {\"error\": null, \"supported\": true");
            sb.append(", \"containerType\": ").append(
                    str(analyzer.getContainerType() == null ? null : analyzer.getContainerType().name()));
            List<ManifestFile> manifestFiles = analyzer.getManifestFiles();
            sb.append(", \"manifestFiles\": [");
            for (int i = 0; i < manifestFiles.size(); i++) {
                if (i > 0) {
                    sb.append(", ");
                }
                ManifestFile manifestFile = manifestFiles.get(i);
                sb.append("{\"filename\": ").append(str(manifestFile.getFilename()));
                sb.append(", \"signatureFilename\": ").append(str(manifestFile.getSignatureFilename()));
                sb.append(", \"entryCount\": ").append(
                        manifestFile.getEntries() == null ? 0 : manifestFile.getEntries().size());
                sb.append("}");
            }
            sb.append("]");

            List<AdvancedSignature> signatures = analyzer.getSignatures();
            sb.append(",\n   \"signatures\": [");
            for (int i = 0; i < signatures.size(); i++) {
                if (i > 0) {
                    sb.append(", ");
                }
                AdvancedSignature signature = signatures.get(i);
                sb.append("{\"id\": ").append(str(signature.getId()));
                sb.append(", \"filename\": ").append(str(signature.getFilename()));
                sb.append(", \"originalDocuments\": ");
                try {
                    sb.append(names(analyzer.getOriginalDocuments(signature)));
                } catch (Throwable t) {
                    sb.append(str("!" + t.getClass().getName()));
                }
                sb.append("}");
            }
            sb.append("]}");
        } catch (Throwable t) {
            out.print("  {\"error\": " + str(t.getClass().getName())
                    + ", \"errorMessage\": " + str(t.getMessage()) + "}");
            return;
        }
        out.print(sb);
    }

    private interface Extraction {
        ASiCContent extract() throws Exception;
    }

    private static ASiCContent dumpExtraction(PrintStream out, Extraction extraction) {
        ASiCContent content;
        try {
            content = extraction.extract();
        } catch (Throwable t) {
            out.print("  {\"error\": " + str(t.getClass().getName())
                    + ", \"errorMessage\": " + str(t.getMessage()) + "}");
            return null;
        }
        StringBuilder sb = new StringBuilder();
        sb.append("  {\"error\": null");
        sb.append(", \"containerType\": ")
                .append(str(content.getContainerType() == null ? null : content.getContainerType().name()));
        sb.append(", \"zipComment\": ").append(str(content.getZipComment()));
        sb.append(", \"mimeTypeDocument\": ").append(str(
                content.getMimeTypeDocument() == null ? null : content.getMimeTypeDocument().getName()));
        sb.append(",\n   \"signatureDocuments\": ").append(names(content.getSignatureDocuments()));
        sb.append(",\n   \"manifestDocuments\": ").append(names(content.getManifestDocuments()));
        sb.append(",\n   \"archiveManifestDocuments\": ").append(names(content.getArchiveManifestDocuments()));
        sb.append(",\n   \"evidenceRecordManifestDocuments\": ")
                .append(names(content.getEvidenceRecordManifestDocuments()));
        sb.append(",\n   \"timestampDocuments\": ").append(names(content.getTimestampDocuments()));
        sb.append(",\n   \"evidenceRecordDocuments\": ").append(names(content.getEvidenceRecordDocuments()));
        sb.append(",\n   \"signedDocuments\": ").append(names(content.getSignedDocuments()));
        sb.append(",\n   \"unsupportedDocuments\": ").append(names(content.getUnsupportedDocuments()));
        sb.append(",\n   \"folders\": ").append(names(content.getFolders()));
        sb.append(",\n   \"containerDocuments\": ").append(names(content.getContainerDocuments()));
        sb.append(",\n   \"rootLevelSignedDocuments\": ").append(namesSafe(content, "rootLevel"));
        sb.append(",\n   \"allManifestDocuments\": ").append(namesSafe(content, "allManifests"));
        sb.append("}");
        out.print(sb);
        return content;
    }

    private interface ManifestSupplier {
        ManifestFile get() throws Exception;
    }

    private static void dumpManifest(PrintStream out, String name, ManifestSupplier supplier) {
        ManifestFile manifest;
        try {
            manifest = supplier.get();
        } catch (Throwable t) {
            out.print("    {\"name\": " + str(name) + ", \"error\": " + str(t.getClass().getName())
                    + ", \"errorMessage\": " + str(t.getMessage()) + "}");
            return;
        }
        if (manifest == null) {
            out.print("    {\"name\": " + str(name) + ", \"error\": null, \"manifest\": null}");
            return;
        }
        StringBuilder sb = new StringBuilder();
        sb.append("    {\"name\": ").append(str(name)).append(", \"error\": null");
        sb.append(", \"filename\": ").append(str(manifest.getFilename()));
        sb.append(", \"signatureFilename\": ").append(str(manifest.getSignatureFilename()));
        sb.append(", \"manifestType\": ").append(
                str(manifest.getManifestType() == null ? null : manifest.getManifestType().name()));
        sb.append(", \"entries\": [");
        List<ManifestEntry> entries = manifest.getEntries();
        if (entries != null) {
            for (int i = 0; i < entries.size(); i++) {
                ManifestEntry entry = entries.get(i);
                if (i > 0) {
                    sb.append(", ");
                }
                sb.append("{\"uri\": ").append(str(entry.getUri()));
                sb.append(", \"documentName\": ").append(str(
                        entry.getDocument() == null ? null : entry.getDocument().getName()));
                sb.append(", \"mimeType\": ").append(str(
                        entry.getMimeType() == null ? null : entry.getMimeType().getMimeTypeString()));
                sb.append(", \"digestAlgorithm\": ").append(str(
                        entry.getDigest() == null || entry.getDigest().getAlgorithm() == null
                                ? null : entry.getDigest().getAlgorithm().name()));
                sb.append(", \"digestValue\": ").append(str(
                        entry.getDigest() == null ? null : hex(entry.getDigest().getValue())));
                sb.append(", \"rootfile\": ").append(entry.isRootfile());
                sb.append("}");
            }
        }
        sb.append("]}");
        out.print(sb);
    }

    private static String namesSafe(ASiCContent content, String which) {
        try {
            return names("rootLevel".equals(which)
                    ? content.getRootLevelSignedDocuments() : content.getAllManifestDocuments());
        } catch (Throwable t) {
            return str("!" + t.getClass().getName());
        }
    }

    private static String names(List<DSSDocument> documents) {
        StringBuilder sb = new StringBuilder("[");
        if (documents != null) {
            for (int i = 0; i < documents.size(); i++) {
                if (i > 0) {
                    sb.append(", ");
                }
                sb.append(str(documents.get(i).getName()));
            }
        }
        return sb.append("]").toString();
    }

    private static String hex(byte[] value) {
        if (value == null) {
            return null;
        }
        StringBuilder sb = new StringBuilder(value.length * 2);
        for (byte b : value) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    private static String str(String value) {
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

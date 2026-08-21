// Generates testdata/manifest-oracle.json: the exact ASiCManifest XML bytes that upstream DSS
// 6.5.RC1 produces for the four ASiC with CAdES manifest builders this porter chunk owns -
// ASiCWithCAdESSignatureManifestBuilder, ASiCWithCAdESTimestampManifestBuilder and
// ASiCEWithCAdESArchiveManifestBuilder (with and without a Rootfile), plus a non-default digest
// algorithm and URI-encoding-sensitive filenames.
//
// These manifests are SIGNED (the CAdES signature covers the manifest bytes), so element order,
// namespace declarations and placement, attribute order, digest-algorithm URIs, MimeType strings
// and the DSSUtils.encodeURI treatment of entry names are all byte-compatibility surfaces. The Go
// KAT (manifest_kat_test.go) compares its own serialization against these strings rather than
// against anything hand-derived.
//
// The oracle lives in package eu.europa.esig.dss.asic.cades.signature.manifest because
// ASiCEWithCAdESManifestBuilder is abstract and its constructors are protected.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd $DSS_UPSTREAM_HOME  (your built upstream DSS 6.5.RC1 checkout)
//   mvn -q -o -pl dss-asic-cades -am -DskipTests install
//   mvn -q -o -pl dss-asic-cades dependency:build-classpath \
//       -Dmdep.outputFile=/tmp/cp-asic-cades.txt -Dmdep.includeScope=test
//   CP="dss-asic-cades/target/classes:$(cat /tmp/cp-asic-cades.txt)"
//   javac -cp "$CP" -d /tmp/manifestoracle ManifestOracle.java
//   java  -cp "$CP:/tmp/manifestoracle" \
//         eu.europa.esig.dss.asic.cades.signature.manifest.ManifestOracle <testdata directory>
package eu.europa.esig.dss.asic.cades.signature.manifest;

import eu.europa.esig.dss.asic.common.ASiCContent;
import eu.europa.esig.dss.asic.common.AbstractASiCManifestBuilder;
import eu.europa.esig.dss.enumerations.ASiCContainerType;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.MimeTypeEnum;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.utils.Utils;

import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class ManifestOracle {

    private ManifestOracle() {
        // utility
    }

    public static void main(String[] args) throws Exception {
        if (args.length != 1) {
            throw new IllegalArgumentException("Usage: ManifestOracle <testdata directory>");
        }
        Path outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);

        Map<String, Map<String, String>> fixtures = new LinkedHashMap<>();

        // --- signature manifest, SHA-256, plain ASCII names -------------------------------
        fixtures.put("signature-sha256", dump(new ASiCWithCAdESSignatureManifestBuilder(
                asicEContent(), DigestAlgorithm.SHA256, "META-INF/signature001.p7s")));

        // --- signature manifest, SHA-512 --------------------------------------------------
        fixtures.put("signature-sha512", dump(new ASiCWithCAdESSignatureManifestBuilder(
                asicEContent(), DigestAlgorithm.SHA512, "META-INF/signature001.p7s")));

        // --- signature manifest over URI-encoding sensitive filenames ---------------------
        fixtures.put("signature-encoded-uris", dump(new ASiCWithCAdESSignatureManifestBuilder(
                encodingContent(), DigestAlgorithm.SHA256, "META-INF/signature001.p7s")));

        // --- signature manifest over percent-carrying filenames ---------------------------
        // Isolated because DSSUtils.encodeURI percent-encodes a literal '%' to "%25"; see the
        // Go KAT's TestSignatureManifestPercentEncodedURIs.
        fixtures.put("signature-percent-uris", dump(new ASiCWithCAdESSignatureManifestBuilder(
                percentContent(), DigestAlgorithm.SHA256, "META-INF/signature001.p7s")));

        // --- timestamp manifest -----------------------------------------------------------
        fixtures.put("timestamp-sha256", dump(new ASiCWithCAdESTimestampManifestBuilder(
                asicEContent(), DigestAlgorithm.SHA256, "META-INF/timestamp001.tst")));

        // --- archive manifest, no rootfile ------------------------------------------------
        fixtures.put("archive-no-rootfile", dump(new ASiCEWithCAdESArchiveManifestBuilder(
                archiveContent(), null, DigestAlgorithm.SHA256, "META-INF/timestamp002.tst")));

        // --- archive manifest, with rootfile ----------------------------------------------
        ASiCContent contentWithRootfile = archiveContent();
        DSSDocument lastArchiveManifest = contentWithRootfile.getArchiveManifestDocuments().get(0);
        fixtures.put("archive-with-rootfile", dump(new ASiCEWithCAdESArchiveManifestBuilder(
                contentWithRootfile, lastArchiveManifest, DigestAlgorithm.SHA256, "META-INF/timestamp002.tst")));

        // --- archive manifest, SHA-384 ----------------------------------------------------
        fixtures.put("archive-sha384", dump(new ASiCEWithCAdESArchiveManifestBuilder(
                archiveContent(), null, DigestAlgorithm.SHA384, "META-INF/timestamp002.tst")));

        writeJson(outDir.resolve("manifest-oracle.json"), fixtures);
    }

    /** Builds the manifest and captures its filename plus its exact serialized bytes. */
    private static Map<String, String> dump(AbstractASiCManifestBuilder builder) throws Exception {
        DSSDocument manifest = builder.build();
        Map<String, String> entry = new LinkedHashMap<>();
        entry.put("name", manifest.getName() == null ? "" : manifest.getName());
        entry.put("mimeType", manifest.getMimeType() == null ? "" : manifest.getMimeType().getMimeTypeString());
        entry.put("xml", new String(readAll(manifest), StandardCharsets.UTF_8));
        entry.put("base64", Utils.toBase64(readAll(manifest)));
        return entry;
    }

    private static byte[] readAll(DSSDocument document) throws Exception {
        try (InputStream is = document.openStream(); ByteArrayOutputStream baos = new ByteArrayOutputStream()) {
            byte[] buffer = new byte[4096];
            int read;
            while ((read = is.read(buffer)) != -1) {
                baos.write(buffer, 0, read);
            }
            return baos.toByteArray();
        }
    }

    /** ASiC-E container with three signed documents of distinct mime-types. */
    private static ASiCContent asicEContent() {
        ASiCContent asicContent = new ASiCContent();
        asicContent.setContainerType(ASiCContainerType.ASiC_E);
        asicContent.setMimeTypeDocument(new InMemoryDocument(
                MimeTypeEnum.ASICE.getMimeTypeString().getBytes(StandardCharsets.UTF_8),
                "mimetype", MimeTypeEnum.BINARY));
        asicContent.setSignedDocuments(new ArrayList<>(Arrays.asList(
                new InMemoryDocument("Hello World !".getBytes(StandardCharsets.UTF_8), "test.text", MimeTypeEnum.TEXT),
                new InMemoryDocument("<root/>".getBytes(StandardCharsets.UTF_8), "data.xml", MimeTypeEnum.XML),
                new InMemoryDocument(new byte[] { 0x00, 0x01, 0x02, 0x03 }, "binary.bin", MimeTypeEnum.BINARY))));
        return asicContent;
    }

    /** ASiC-E container whose entry names exercise DSSUtils.encodeURI. */
    private static ASiCContent encodingContent() {
        ASiCContent asicContent = new ASiCContent();
        asicContent.setContainerType(ASiCContainerType.ASiC_E);
        asicContent.setSignedDocuments(new ArrayList<>(Arrays.asList(
                new InMemoryDocument("a".getBytes(StandardCharsets.UTF_8), "document 2.txt", MimeTypeEnum.TEXT),
                new InMemoryDocument("b".getBytes(StandardCharsets.UTF_8), "détaché.txt", MimeTypeEnum.TEXT),
                new InMemoryDocument("c".getBytes(StandardCharsets.UTF_8), "dir/sub dir/file&name.txt", MimeTypeEnum.TEXT),
                new InMemoryDocument("d".getBytes(StandardCharsets.UTF_8), "a[b]<c>{d}|e\\f^g`h\"i.txt", MimeTypeEnum.TEXT))));
        return asicContent;
    }

    /** ASiC-E container whose entry names carry a literal '%', which encodeURI escapes to "%25". */
    private static ASiCContent percentContent() {
        ASiCContent asicContent = new ASiCContent();
        asicContent.setContainerType(ASiCContainerType.ASiC_E);
        asicContent.setSignedDocuments(new ArrayList<>(Arrays.asList(
                new InMemoryDocument("a".getBytes(StandardCharsets.UTF_8), "100%_done.txt", MimeTypeEnum.TEXT),
                new InMemoryDocument("b".getBytes(StandardCharsets.UTF_8), "already%20encoded.txt", MimeTypeEnum.TEXT))));
        return asicContent;
    }

    /** ASiC-E container carrying signatures, timestamps and manifests, as an archive manifest sees it. */
    private static ASiCContent archiveContent() {
        ASiCContent asicContent = asicEContent();
        asicContent.setSignatureDocuments(new ArrayList<>(List.of(
                new InMemoryDocument("signature-bytes".getBytes(StandardCharsets.UTF_8),
                        "META-INF/signature001.p7s", MimeTypeEnum.PKCS7))));
        asicContent.setTimestampDocuments(new ArrayList<>(List.of(
                new InMemoryDocument("timestamp-bytes".getBytes(StandardCharsets.UTF_8),
                        "META-INF/timestamp001.tst", MimeTypeEnum.TST))));
        asicContent.setManifestDocuments(new ArrayList<>(List.of(
                new InMemoryDocument("<manifest/>".getBytes(StandardCharsets.UTF_8),
                        "META-INF/ASiCManifest001.xml", MimeTypeEnum.XML))));
        asicContent.setArchiveManifestDocuments(new ArrayList<>(List.of(
                new InMemoryDocument("<archive-manifest/>".getBytes(StandardCharsets.UTF_8),
                        "META-INF/ASiCArchiveManifest.xml", MimeTypeEnum.XML))));
        return asicContent;
    }

    private static void writeJson(Path file, Map<String, Map<String, String>> fixtures) throws Exception {
        StringBuilder sb = new StringBuilder();
        sb.append("{\n");
        int i = 0;
        for (Map.Entry<String, Map<String, String>> fixture : fixtures.entrySet()) {
            sb.append("  ").append(quote(fixture.getKey())).append(": {\n");
            int j = 0;
            for (Map.Entry<String, String> field : fixture.getValue().entrySet()) {
                sb.append("    ").append(quote(field.getKey())).append(": ").append(quote(field.getValue()));
                sb.append(++j == fixture.getValue().size() ? "\n" : ",\n");
            }
            sb.append("  }").append(++i == fixtures.size() ? "\n" : ",\n");
        }
        sb.append("}\n");
        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(file, StandardCharsets.UTF_8))) {
            writer.print(sb);
        }
    }

    private static String quote(String value) {
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

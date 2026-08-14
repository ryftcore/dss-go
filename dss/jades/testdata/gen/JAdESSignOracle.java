// Generates testdata/jades-sign-oracle.json: everything the SIGN chunk of the JAdES port has to
// reproduce byte for byte - the serialized JWS protected header JAdESLevelBaselineB builds
// (insertion order included), the JWS payload each SignaturePackaging / 'sigD' mechanism / 'b64'
// combination yields, the signing input BASE64URL(header) '.' BASE64URL(payload), and the final
// compact, flattened-JSON and complete-JSON serializations JAdESCompactBuilder and
// JAdESSerializationBuilder emit for a FIXED signature value.
//
// Nothing in the Go KAT is hand-derived: every expectation is a string produced here by upstream
// DSS 6.5.RC1 driving its own builders.
//
// The oracle lives in package eu.europa.esig.dss.jades.signature because JAdESLevelBaselineB's
// getSignedProperties()/getPayloadBytes() are public but the builders' assertions are not.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies:
//
//   cd /home/user/dss-upstream
//   mvn -q -o -pl dss-jades -am install -DskipTests
//   mvn -q -o -pl dss-jades dependency:build-classpath -Dmdep.outputFile=/tmp/cp-jades.txt \
//       -Dmdep.includeScope=test
//   CP="dss-jades/target/classes:$(cat /tmp/cp-jades.txt)"
//   javac -cp "$CP" -d /tmp/jadesoracle JAdESSignOracle.java
//   java  -cp "$CP:/tmp/jadesoracle" \
//         eu.europa.esig.dss.jades.signature.JAdESSignOracle <testdata directory>
package eu.europa.esig.dss.jades.signature;

import eu.europa.esig.dss.enumerations.CommitmentTypeEnum;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.JWSSerializationType;
import eu.europa.esig.dss.enumerations.MimeTypeEnum;
import eu.europa.esig.dss.enumerations.ObjectIdentifierQualifier;
import eu.europa.esig.dss.enumerations.SigDMechanism;
import eu.europa.esig.dss.enumerations.SignatureAlgorithm;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.enumerations.SignaturePackaging;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.jades.HTTPHeader;
import eu.europa.esig.dss.jades.HTTPHeaderDigest;
import eu.europa.esig.dss.jades.JAdESSignatureParameters;
import eu.europa.esig.dss.jades.JAdESSigningTimeType;
import eu.europa.esig.dss.model.CommonCommitmentType;
import eu.europa.esig.dss.model.CommitmentQualifier;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.Policy;
import eu.europa.esig.dss.model.SignatureValue;
import eu.europa.esig.dss.model.SignerLocation;
import eu.europa.esig.dss.model.SpDocSpecification;
import eu.europa.esig.dss.model.ToBeSigned;
import eu.europa.esig.dss.model.UserNotice;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.validation.CertificateVerifier;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.tsp.TimestampToken;
import eu.europa.esig.dss.utils.Utils;
import org.jose4j.json.JsonUtil;

import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Date;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public class JAdESSignOracle {

    /** Fixed signing time: 2020-09-13T12:26:40Z. */
    static final Date SIGNING_DATE = new Date(1600000000000L);

    /** Fixed 'exp' (expiration time) claim: 2030-01-01T00:00:00Z. */
    static final Date EXPIRATION_DATE = new Date(1893456000000L);

    /** A fixed, arbitrary "signature value": the builders only copy it through. */
    static final byte[] SIGNATURE_VALUE = new byte[64];
    static {
        for (int i = 0; i < SIGNATURE_VALUE.length; i++) {
            SIGNATURE_VALUE[i] = (byte) (i * 7 + 3);
        }
    }

    /** The enveloped/enveloping payload: URL-safe ASCII, so 'b64'=false compact works too. */
    static final byte[] PAYLOAD = "Hello JAdES!".getBytes(StandardCharsets.US_ASCII);

    static CertificateToken signingCertificate;
    static byte[] contentTimestampBinary;

    public static void main(String[] args) throws Exception {
        Path outDir = Paths.get(args[0]);
        signingCertificate = DSSUtils.loadCertificate(Files.readAllBytes(outDir.resolve("signer.crt")));
        contentTimestampBinary = Files.readAllBytes(outDir.resolve("content-timestamp.tst"));

        List<Map<String, Object>> cases = new ArrayList<>();
        for (Case c : cases()) {
            cases.add(run(c));
        }

        Map<String, Object> root = new LinkedHashMap<>();
        root.put("signingCertificate", Utils.toBase64(signingCertificate.getEncoded()));
        root.put("signatureValue", Utils.toBase64(SIGNATURE_VALUE));
        root.put("payload", new String(PAYLOAD, StandardCharsets.UTF_8));
        root.put("signingDateMillis", SIGNING_DATE.getTime());
        root.put("expirationDateMillis", EXPIRATION_DATE.getTime());
        root.put("cases", cases);
        root.put("httpHeaders", httpHeadersPayloads());

        try (PrintWriter writer = new PrintWriter(
                Files.newBufferedWriter(outDir.resolve("jades-sign-oracle.json"), StandardCharsets.UTF_8))) {
            writer.println(JsonUtil.toJson(root));
        }
        System.out.println("wrote " + cases.size() + " cases");
    }

    /** One oracle row: a name plus the parameter/document configuration it exercises. */
    interface Case {
        String name();
        JAdESSignatureParameters parameters();
        List<DSSDocument> documents();
    }

    static Map<String, Object> run(Case c) {
        Map<String, Object> row = new LinkedHashMap<>();
        row.put("name", c.name());

        JAdESSignatureParameters parameters = c.parameters();
        List<DSSDocument> documents = c.documents();
        CertificateVerifier certificateVerifier = new CommonCertificateVerifier();

        try {
            JAdESLevelBaselineB levelBaselineB = new JAdESLevelBaselineB(certificateVerifier, parameters, documents);
            row.put("protectedHeader", JsonUtil.toJson(levelBaselineB.getSignedProperties()));
            row.put("payloadBase64", Utils.toBase64(levelBaselineB.getPayloadBytes()));
        } catch (Exception e) {
            row.put("levelBError", e.getClass().getSimpleName() + ": " + e.getMessage());
            return row;
        }

        // A fresh JAdESLevelBaselineB per builder: getSignedProperties() is not idempotent
        // (it appends to the same map), exactly as upstream instantiates one builder per use.
        try {
            JAdESBuilder builder = builder(certificateVerifier, c.parameters(), documents);
            ToBeSigned toBeSigned = builder.buildDataToBeSigned();
            row.put("dataToBeSigned", new String(toBeSigned.getBytes(), StandardCharsets.UTF_8));
            row.put("mimeType", builder.getMimeType().getMimeTypeString());
        } catch (Exception e) {
            row.put("dataToBeSignedError", e.getClass().getSimpleName() + ": " + e.getMessage());
        }

        try {
            JAdESBuilder builder = builder(certificateVerifier, c.parameters(), documents);
            SignatureValue signatureValue = new SignatureValue(SignatureAlgorithm.RSA_SHA256, SIGNATURE_VALUE);
            DSSDocument signed = builder.build(signatureValue);
            row.put("signature", new String(DSSUtils.toByteArray(signed), StandardCharsets.UTF_8));
        } catch (Exception e) {
            row.put("signatureError", e.getClass().getSimpleName() + ": " + e.getMessage());
        }
        return row;
    }

    static JAdESBuilder builder(CertificateVerifier certificateVerifier, JAdESSignatureParameters parameters,
                                List<DSSDocument> documents) {
        switch (parameters.getJwsSerializationType()) {
            case COMPACT_SERIALIZATION:
                return new JAdESCompactBuilder(certificateVerifier, parameters, documents);
            default:
                return new JAdESSerializationBuilder(certificateVerifier, parameters, documents);
        }
    }

    // ------------------------------------------------------------------ parameter fixtures

    static JAdESSignatureParameters baseParameters() {
        JAdESSignatureParameters parameters = new JAdESSignatureParameters();
        parameters.setSigningCertificate(signingCertificate);
        parameters.setCertificateChain(signingCertificate);
        parameters.setSignatureLevel(SignatureLevel.JAdES_BASELINE_B);
        parameters.setDigestAlgorithm(DigestAlgorithm.SHA256);
        parameters.setSigningCertificateDigestMethod(DigestAlgorithm.SHA256);
        parameters.bLevel().setSigningDate(SIGNING_DATE);
        return parameters;
    }

    static DSSDocument payloadDocument() {
        InMemoryDocument document = new InMemoryDocument(PAYLOAD, "payload.txt");
        document.setMimeType(MimeTypeEnum.TEXT);
        return document;
    }

    static List<DSSDocument> detachedDocuments() {
        InMemoryDocument first = new InMemoryDocument("first detached".getBytes(StandardCharsets.UTF_8), "first.txt");
        first.setMimeType(MimeTypeEnum.TEXT);
        InMemoryDocument second = new InMemoryDocument("second detached".getBytes(StandardCharsets.UTF_8), "second.bin");
        second.setMimeType(MimeTypeEnum.BINARY);
        return Arrays.asList(first, second);
    }

    static List<DSSDocument> httpHeaderDocuments() {
        List<DSSDocument> documents = new ArrayList<>();
        documents.add(new HTTPHeader("content-type", "application/json"));
        documents.add(new HTTPHeader("x-example", " leading and trailing "));
        documents.add(new HTTPHeader("x-example", "second value"));
        documents.add(new HTTPHeaderDigest(
                new InMemoryDocument("{\"hello\":\"world\"}".getBytes(StandardCharsets.UTF_8)),
                DigestAlgorithm.SHA256));
        return documents;
    }

    /**
     * One case reuses a single parameter object across the three runs below, exactly as a caller
     * would: JAdESLevelBaselineB reads it without mutating it, and a fresh JAdESLevelBaselineB is
     * built per run because getSignedProperties() appends to its own map and is not idempotent.
     */
    static Case simple(String name, JAdESSignatureParameters parameters, List<DSSDocument> documents) {
        return new Case() {
            @Override public String name() { return name; }
            @Override public JAdESSignatureParameters parameters() { return parameters; }
            @Override public List<DSSDocument> documents() { return documents; }
        };
    }

    static List<Case> cases() {
        List<Case> cases = new ArrayList<>();

        // 1. Compact, enveloping, all defaults.
        cases.add(simple("compact-enveloping-default", enveloping(JWSSerializationType.COMPACT_SERIALIZATION),
                Collections.singletonList(payloadDocument())));

        // 2. Compact, enveloping, RFC 7797 unencoded payload.
        JAdESSignatureParameters b64False = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        b64False.setBase64UrlEncodedPayload(false);
        cases.add(simple("compact-enveloping-b64false", b64False, Collections.singletonList(payloadDocument())));

        // 3. Complete JSON serialization, enveloping.
        cases.add(simple("json-enveloping-default", enveloping(JWSSerializationType.JSON_SERIALIZATION),
                Collections.singletonList(payloadDocument())));

        // 4. Flattened JSON serialization, enveloping.
        cases.add(simple("flattened-enveloping-default",
                enveloping(JWSSerializationType.FLATTENED_JSON_SERIALIZATION),
                Collections.singletonList(payloadDocument())));

        // 5. Flattened JSON serialization, enveloping, unencoded payload.
        JAdESSignatureParameters flatB64False = enveloping(JWSSerializationType.FLATTENED_JSON_SERIALIZATION);
        flatB64False.setBase64UrlEncodedPayload(false);
        cases.add(simple("flattened-enveloping-b64false", flatB64False,
                Collections.singletonList(payloadDocument())));

        // 6. Compact, detached, NO_SIG_D.
        JAdESSignatureParameters noSigD = detached(JWSSerializationType.COMPACT_SERIALIZATION, SigDMechanism.NO_SIG_D);
        cases.add(simple("compact-detached-nosigd", noSigD, Collections.singletonList(payloadDocument())));

        // 7. JSON, detached, ObjectIdByURI.
        cases.add(simple("json-detached-objectidbyuri",
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.OBJECT_ID_BY_URI),
                detachedDocuments()));

        // 8. JSON, detached, ObjectIdByURI, unencoded payload.
        JAdESSignatureParameters uriB64False =
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.OBJECT_ID_BY_URI);
        uriB64False.setBase64UrlEncodedPayload(false);
        cases.add(simple("json-detached-objectidbyuri-b64false", uriB64False, detachedDocuments()));

        // 9. JSON, detached, ObjectIdByURIHash.
        cases.add(simple("json-detached-objectidbyurihash",
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.OBJECT_ID_BY_URI_HASH),
                detachedDocuments()));

        // 10. JSON, detached, ObjectIdByURIHash, unencoded payload (digest over raw octets).
        JAdESSignatureParameters hashB64False =
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.OBJECT_ID_BY_URI_HASH);
        hashB64False.setBase64UrlEncodedPayload(false);
        cases.add(simple("json-detached-objectidbyurihash-b64false", hashB64False, detachedDocuments()));

        // 11. JSON, detached, ObjectIdByURIHash with an explicit reference digest algorithm.
        JAdESSignatureParameters hashSha512 =
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.OBJECT_ID_BY_URI_HASH);
        hashSha512.setReferenceDigestAlgorithm(DigestAlgorithm.SHA512);
        cases.add(simple("json-detached-objectidbyurihash-sha512", hashSha512, detachedDocuments()));

        // 12. JSON, detached, HttpHeaders ('b64' must be false).
        JAdESSignatureParameters httpHeaders =
                detached(JWSSerializationType.JSON_SERIALIZATION, SigDMechanism.HTTP_HEADERS);
        httpHeaders.setBase64UrlEncodedPayload(false);
        cases.add(simple("json-detached-httpheaders", httpHeaders, httpHeaderDocuments()));

        // 13. 'sigT' claimed signing time instead of 'iat'.
        JAdESSignatureParameters sigT = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        sigT.setJadesSigningTimeType(JAdESSigningTimeType.SIG_T);
        cases.add(simple("compact-enveloping-sigt", sigT, Collections.singletonList(payloadDocument())));

        // 14. No claimed signing time at all.
        JAdESSignatureParameters noTime = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        noTime.setJadesSigningTimeType(JAdESSigningTimeType.NONE);
        cases.add(simple("compact-enveloping-no-signing-time", noTime, Collections.singletonList(payloadDocument())));

        // 15. 'x5t#o' instead of 'x5t#S256'.
        JAdESSignatureParameters x5to = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        x5to.setSigningCertificateDigestMethod(DigestAlgorithm.SHA512);
        cases.add(simple("compact-enveloping-x5to", x5to, Collections.singletonList(payloadDocument())));

        // 16. Minimal protected header: no 'kid', no 'x5c', no 'typ', explicit 'cty'.
        JAdESSignatureParameters minimal = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        minimal.setIncludeKeyIdentifier(false);
        minimal.setIncludeCertificateChain(false);
        minimal.setIncludeSignatureType(false);
        minimal.setContentType("application/vnd.oracle.test+json");
        cases.add(simple("compact-enveloping-minimal", minimal, Collections.singletonList(payloadDocument())));

        // 17. Explicit 'kid', 'x5u' and a custom 'typ'.
        JAdESSignatureParameters explicitKid = enveloping(JWSSerializationType.COMPACT_SERIALIZATION);
        explicitKid.setKeyIdentifier("my-key-identifier");
        explicitKid.setX509Url("https://example.org/certs/signer.pem");
        explicitKid.setSignatureType("application/jose+json");
        explicitKid.setExpirationTime(EXPIRATION_DATE);
        cases.add(simple("compact-enveloping-kid-x5u-typ-exp", explicitKid,
                Collections.singletonList(payloadDocument())));

        // 18. The full TS 119 182-1 signed-header surface.
        cases.add(simple("json-enveloping-full", full(), Collections.singletonList(payloadDocument())));

        return cases;
    }

    static JAdESSignatureParameters enveloping(JWSSerializationType type) {
        JAdESSignatureParameters parameters = baseParameters();
        parameters.setSignaturePackaging(SignaturePackaging.ENVELOPING);
        parameters.setJwsSerializationType(type);
        return parameters;
    }

    static JAdESSignatureParameters detached(JWSSerializationType type, SigDMechanism mechanism) {
        JAdESSignatureParameters parameters = baseParameters();
        parameters.setSignaturePackaging(SignaturePackaging.DETACHED);
        parameters.setJwsSerializationType(type);
        parameters.setSigDMechanism(mechanism);
        return parameters;
    }

    static JAdESSignatureParameters full() {
        JAdESSignatureParameters parameters = enveloping(JWSSerializationType.JSON_SERIALIZATION);

        SignerLocation signerLocation = new SignerLocation();
        signerLocation.setCountry("LU");
        signerLocation.setLocality("Luxembourg");
        signerLocation.setStateOrProvince("Luxembourg");
        signerLocation.setPostOfficeBoxNumber("PO-42");
        signerLocation.setPostalCode("L-1234");
        signerLocation.setStreetAddress("1 rue de la Gare");
        parameters.bLevel().setSignerLocation(signerLocation);

        parameters.bLevel().setClaimedSignerRoles(Arrays.asList("head-of-department", "auditor"));
        parameters.bLevel().setSignedAssertions(Collections.singletonList("{\"assertion\":\"value\"}"));

        CommonCommitmentType commitment = new CommonCommitmentType();
        commitment.setUri("http://uri.etsi.org/01903/v1.2.2#ProofOfOrigin");
        commitment.setDescription("Proof of origin");
        commitment.setDocumentationReferences("https://example.org/doc1", "https://example.org/doc2");
        CommitmentQualifier jsonQualifier = new CommitmentQualifier();
        jsonQualifier.setContent(new InMemoryDocument(
                "{\"qualifier\":\"json\",\"nested\":{\"a\":1}}".getBytes(StandardCharsets.UTF_8)));
        CommitmentQualifier textQualifier = new CommitmentQualifier();
        textQualifier.setContent(new InMemoryDocument("plain qualifier".getBytes(StandardCharsets.UTF_8)));
        commitment.setCommitmentTypeQualifiers(jsonQualifier, textQualifier);
        parameters.bLevel().setCommitmentTypeIndications(
                Arrays.asList(CommitmentTypeEnum.ProofOfReceipt, commitment));

        Policy policy = new Policy();
        policy.setId("urn:oid:1.2.3.4.5");
        policy.setDescription("Oracle test policy");
        policy.setDocumentationReferences("https://example.org/policy.pdf");
        policy.setDigestAlgorithm(DigestAlgorithm.SHA256);
        policy.setDigestValue(DSSUtils.digest(DigestAlgorithm.SHA256, "policy".getBytes(StandardCharsets.UTF_8)));
        policy.setSpuri("https://example.org/policy");
        UserNotice userNotice = new UserNotice();
        userNotice.setOrganization("DSS Go Port");
        userNotice.setNoticeNumbers(1, 2, 3);
        userNotice.setExplicitText("Explicit notice text");
        policy.setUserNotice(userNotice);
        SpDocSpecification spDocSpecification = new SpDocSpecification();
        spDocSpecification.setId("urn:oid:1.2.3.4.6");
        spDocSpecification.setDescription("Policy document specification");
        spDocSpecification.setDocumentationReferences("https://example.org/spec");
        spDocSpecification.setQualifier(ObjectIdentifierQualifier.OID_AS_URN);
        policy.setSpDocSpecification(spDocSpecification);
        parameters.bLevel().setSignaturePolicy(policy);

        try {
            parameters.setContentTimestamps(Collections.singletonList(
                    new TimestampToken(contentTimestampBinary, TimestampType.CONTENT_TIMESTAMP)));
        } catch (Exception e) {
            throw new IllegalStateException("cannot load the content timestamp", e);
        }
        return parameters;
    }

    // ------------------------------------------------------------- HttpHeadersPayloadBuilder

    static Map<String, Object> httpHeadersPayloads() {
        Map<String, Object> result = new LinkedHashMap<>();
        result.put("signature", Utils.toBase64(
                new HttpHeadersPayloadBuilder(httpHeaderDocuments(), false).build()));
        result.put("timestamp", Utils.toBase64(
                new HttpHeadersPayloadBuilder(httpHeaderDocuments(), true).build()));
        return result;
    }

}

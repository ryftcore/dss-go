import eu.europa.esig.dss.enumerations.CommitmentTypeEnum;
import eu.europa.esig.dss.enumerations.JWSSerializationType;
import eu.europa.esig.dss.jades.DSSJsonUtils;
import eu.europa.esig.dss.jades.JWSCompactSerializationParser;
import eu.europa.esig.dss.jades.JWSConverter;
import eu.europa.esig.dss.jades.JWSJsonSerializationGenerator;
import eu.europa.esig.dss.jades.JWSJsonSerializationObject;
import eu.europa.esig.dss.jades.JWSJsonSerializationParser;
import eu.europa.esig.dss.jades.JsonObject;
import eu.europa.esig.dss.jades.validation.JWS;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.TimestampBinary;
import eu.europa.esig.dss.spi.DSSUtils;
import org.jose4j.json.internal.json_simple.JSONArray;
import org.jose4j.jws.AlgorithmIdentifiers;
import org.jose4j.jws.JsonWebSignature;
import org.jose4j.jwx.HeaderParameterNames;

import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.SecureRandom;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Ground truth for the Go port of DSSJsonUtils, JsonObject, JWSConverter, the two
 * JWSJsonSerialization classes and JWSCompactSerializationParser, taken from dss-jades 6.5.RC1
 * itself. TSV, every payload column hex of UTF-8 bytes.
 */
public class JadesOracle {

    static PrintStream out;

    static String hex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    static String hex(String s) {
        return hex(s.getBytes(StandardCharsets.UTF_8));
    }

    static void row(String... cols) {
        out.println(String.join("\t", cols));
    }

    /** Builds a compact JWS with jose4j, optionally with an unencoded (b64=false) payload. */
    static String compact(KeyPair kp, String payload, boolean unencoded) throws Exception {
        JsonWebSignature jws = new JsonWebSignature();
        jws.setAlgorithmHeaderValue(AlgorithmIdentifiers.RSA_USING_SHA256);
        jws.setHeader(HeaderParameterNames.TYPE, "jose");
        jws.setHeader("sigT", "2026-08-14T10:00:00Z");
        if (unencoded) {
            jws.setHeader(HeaderParameterNames.BASE64URL_ENCODE_PAYLOAD, Boolean.FALSE);
            jws.setCriticalHeaderNames("b64");
        }
        jws.setPayload(payload);
        jws.setKey(kp.getPrivate());
        jws.setDoKeyValidation(false);
        return jws.getCompactSerialization();
    }

    public static void main(String[] args) throws Exception {
        out = new PrintStream(System.out, true, "UTF-8");

        KeyPairGenerator g = KeyPairGenerator.getInstance("RSA");
        g.initialize(2048, new SecureRandom());
        KeyPair kp = g.generateKeyPair();

        String[][] compacts = {
            { "encoded", compact(kp, "Hello JAdES payload", false) },
            { "unencoded", compact(kp, "unencoded payload no period", true) },
            { "utf8", compact(kp, "payload with non-ASCII: e-acute and CJK", false) },
        };

        for (String[] c : compacts) {
            String name = c[0];
            String cs = c[1];
            row("COMPACT", name, "", hex(cs));

            DSSDocument doc = new InMemoryDocument(cs.getBytes(StandardCharsets.UTF_8));
            JWSCompactSerializationParser parser = new JWSCompactSerializationParser(doc);
            row("COMPACT_SUPPORTED", name, hex(cs), String.valueOf(parser.isSupported()));

            JWS jws = parser.parse();
            row("SIGNING_INPUT", name, hex(cs), hex(DSSJsonUtils.getSigningInputBytes(jws)));
            row("SIGNED_PAYLOAD", name, hex(cs), hex(jws.getSignedPayload()));
            row("ENCODED_HEADER", name, hex(cs), hex(jws.getEncodedHeader()));

            JWSJsonSerializationObject obj = DSSJsonUtils.toJWSJsonSerializationObject(jws);
            DSSDocument flat = new JWSJsonSerializationGenerator(obj,
                    JWSSerializationType.FLATTENED_JSON_SERIALIZATION).generate();
            row("GEN_FLATTENED", name, hex(cs), hex(DSSUtils.toByteArray(flat)));

            JWSJsonSerializationObject obj2 = DSSJsonUtils.toJWSJsonSerializationObject(parser.parse());
            DSSDocument full = new JWSJsonSerializationGenerator(obj2,
                    JWSSerializationType.JSON_SERIALIZATION).generate();
            row("GEN_COMPLETE", name, hex(cs), hex(DSSUtils.toByteArray(full)));

            DSSDocument convFlat = JWSConverter.fromJWSCompactToJSONFlattenedSerialization(doc);
            row("CONV_FLATTENED", name, hex(cs), hex(DSSUtils.toByteArray(convFlat)));
            DSSDocument convFull = JWSConverter.fromJWSCompactToJSONSerialization(doc);
            row("CONV_COMPLETE", name, hex(cs), hex(DSSUtils.toByteArray(convFull)));

            // Re-parse the complete serialization and report what the parser made of it.
            JWSJsonSerializationParser p2 = new JWSJsonSerializationParser(convFull);
            JWSJsonSerializationObject reparsed = p2.parse();
            row("REPARSE", name, hex(DSSUtils.toByteArray(convFull)),
                    reparsed.getJWSSerializationType() + "|" + reparsed.getSignatures().size()
                            + "|" + hex(reparsed.getPayload())
                            + "|" + hex(reparsed.getSignatures().get(0).getEncodedHeader())
                            + "|" + hex(DSSJsonUtils.getSigningInputBytes(reparsed.getSignatures().get(0))));
        }

        // ---- etsiU incorporation conversions ----
        {
            String cs = compacts[0][1];
            JWS jws = new JWSCompactSerializationParser(
                    new InMemoryDocument(cs.getBytes(StandardCharsets.UTF_8))).parse();

            // An unprotected header with a base64url-encoded etsiU array of two components.
            Map<String, Object> inner1 = new LinkedHashMap<>();
            inner1.put("zeta", "1");
            inner1.put("alpha", "2");
            Map<String, Object> comp1 = new LinkedHashMap<>();
            comp1.put("xVals", inner1);
            Map<String, Object> comp2 = new LinkedHashMap<>();
            comp2.put("rVals", "value/with slash and \"quote\"");

            List<Object> b64Components = new ArrayList<>();
            b64Components.add(DSSJsonUtils.toBase64Url(new JsonObject(comp1)));
            b64Components.add(DSSJsonUtils.toBase64Url(new JsonObject(comp2)));

            Map<String, Object> unprotected = new LinkedHashMap<>();
            unprotected.put("etsiU", b64Components);
            jws.setUnprotected(unprotected);

            JWSJsonSerializationObject obj = DSSJsonUtils.toJWSJsonSerializationObject(jws);
            DSSDocument serialized = new JWSJsonSerializationGenerator(obj,
                    JWSSerializationType.JSON_SERIALIZATION).generate();
            row("ETSIU_B64_DOC", "doc", "", hex(DSSUtils.toByteArray(serialized)));

            DSSDocument clear = JWSConverter.fromEtsiUWithBase64UrlToClearJsonIncorporation(serialized);
            row("ETSIU_TO_CLEAR", "doc", hex(DSSUtils.toByteArray(serialized)), hex(DSSUtils.toByteArray(clear)));

            DSSDocument back = JWSConverter.fromEtsiUWithClearJsonToBase64UrlIncorporation(clear);
            row("ETSIU_TO_B64", "doc", hex(DSSUtils.toByteArray(clear)), hex(DSSUtils.toByteArray(back)));
        }

        // ---- DSSJsonUtils scalars ----
        row("OID", "commitment", "", hex(DSSJsonUtils.getOidObject(CommitmentTypeEnum.ProofOfOrigin).toJSONString()));
        row("OID", "uri-only", "", hex(DSSJsonUtils.getOidObject("http://x/y", null, null).toJSONString()));
        row("OID", "uri-desc", "", hex(DSSJsonUtils.getOidObject("http://x/y", "a desc", null).toJSONString()));
        row("OID", "uri-desc-refs", "", hex(DSSJsonUtils
                .getOidObject("http://x/y", "a desc", new String[] { "r1", "r2" }).toJSONString()));

        {
            List<TimestampBinary> bins = Arrays.asList(
                    new TimestampBinary(new byte[] { 1, 2, 3 }),
                    new TimestampBinary(new byte[] { (byte) 0xff, 0 }));
            row("TSTCONTAINER", "no-canon", "", hex(DSSJsonUtils.getTstContainer(bins, null).toJSONString()));
            row("TSTCONTAINER", "canon", "",
                    hex(DSSJsonUtils.getTstContainer(bins, "http://www.w3.org/2001/10/xml-exc-c14n#").toJSONString()));
        }

        {
            JSONArray array = new JSONArray(Arrays.asList("a", "b"));
            row("B64URL_OBJECT", "array", "", hex(DSSJsonUtils.toBase64Url(array)));
            Map<String, Object> m = new LinkedHashMap<>();
            m.put("b", "2");
            m.put("a", "1");
            row("B64URL_OBJECT", "object", "", hex(DSSJsonUtils.toBase64Url(new JsonObject(m))));
        }

        String[] b64UrlChecks = { "abcXYZ012-_", "abc=", "a.b", "a b", "", "a+b", "a/b" };
        for (String s : b64UrlChecks) {
            row("IS_B64URL", "c", hex(s), String.valueOf(DSSJsonUtils.isBase64UrlEncoded(s)));
        }
        String[] urlSafeChecks = { "abc", "a.b", "a b", "", "a\tb", "a~b" };
        for (String s : urlSafeChecks) {
            row("IS_URLSAFE_PAYLOAD", "c", hex(s), String.valueOf(DSSJsonUtils.isUrlSafePayload(s)));
        }
        String[] mimeChecks = { "json", "application/json", "", "text/plain", "octet-stream" };
        for (String s : mimeChecks) {
            row("MIMETYPE", "c", hex(s), hex(DSSJsonUtils.getMimeTypeString(s)));
        }
        row("CONCAT", "three", "", hex(DSSJsonUtils.concatenate("a", "b", "c")));
        row("CONCAT", "empty-mid", "", hex(DSSJsonUtils.concatenate("a", "", "c")));

        // ---- JsonObject: the no-argument constructor really is HashMap-ordered ----
        {
            JsonObject o = new JsonObject();
            o.put("crlVals", "c");
            o.put("ocspVals", "o");
            row("JSONOBJECT", "rVals", "", hex(o.toJSONString()));

            JsonObject o2 = new JsonObject();
            o2.put("xVals", "x");
            o2.put("rVals", "r");
            row("JSONOBJECT", "tstVd", "", hex(o2.toJSONString()));
        }

        // ---- isSupported on non-JWS inputs ----
        String[] notCompact = { "{\"payload\":\"x\"}", "a.b.c.d", "abc", "a.b.c\nrest", "a.b.c\n" };
        for (String s : notCompact) {
            DSSDocument d = new InMemoryDocument(s.getBytes(StandardCharsets.UTF_8));
            row("COMPACT_SUPPORTED", "neg", hex(s),
                    String.valueOf(new JWSCompactSerializationParser(d).isSupported()));
            row("JSON_SUPPORTED", "neg", hex(s),
                    String.valueOf(new JWSJsonSerializationParser(d).isSupported()));
        }
    }
}

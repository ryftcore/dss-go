import org.jose4j.jws.AlgorithmIdentifiers;
import org.jose4j.jws.JsonWebSignature;
import org.jose4j.jwx.HeaderParameterNames;
import org.jose4j.lang.JoseException;

import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.SecureRandom;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;

/**
 * Produces signed JWSs with jose4j 0.9.6 so the Go port can be checked against real signatures
 * rather than against its own idea of one. TSV columns:
 *
 *   JWS \t name \t alg \t hex(X.509 SubjectPublicKeyInfo) \t compactSerialization \t hex(signingInput)
 */
public class JwsOracle {

    static PrintStream out;

    static String hex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    static KeyPair rsa(int bits) throws Exception {
        KeyPairGenerator g = KeyPairGenerator.getInstance("RSA");
        g.initialize(bits, new SecureRandom());
        return g.generateKeyPair();
    }

    static KeyPair ec(String curve) throws Exception {
        KeyPairGenerator g = KeyPairGenerator.getInstance("EC");
        g.initialize(new ECGenParameterSpec(curve), new SecureRandom());
        return g.generateKeyPair();
    }

    static KeyPair ed25519() throws Exception {
        KeyPairGenerator g = KeyPairGenerator.getInstance("Ed25519");
        return g.generateKeyPair();
    }

    /** Signs and dumps one case. */
    static void emit(String name, String alg, KeyPair kp, String payload, boolean unencodedPayload,
                     String[] crit) throws JoseException {
        JsonWebSignature jws = new JsonWebSignature();
        jws.setAlgorithmHeaderValue(alg);
        // Insertion order below is the order the header members will be serialized in.
        jws.setHeader(HeaderParameterNames.TYPE, "JOSE");
        jws.setHeader("sigT", "2026-08-14T10:00:00Z");
        if (unencodedPayload) {
            jws.setHeader(HeaderParameterNames.BASE64URL_ENCODE_PAYLOAD, Boolean.FALSE);
        }
        if (crit != null) {
            jws.setHeader(HeaderParameterNames.CRITICAL, Arrays.asList(crit));
            jws.setKnownCriticalHeaders(crit);
        }
        jws.setPayload(payload);
        PrivateKey priv = kp.getPrivate();
        jws.setKey(priv);
        jws.setDoKeyValidation(false);

        String compact;
        if (unencodedPayload) {
            // RFC 7797 forbids a '.' in an unencoded compact payload; the fixtures avoid one.
            compact = jws.getCompactSerialization();
        } else {
            compact = jws.getCompactSerialization();
        }

        // Recover the signing input the same way the verifier does.
        JsonWebSignature check = new JsonWebSignature();
        check.setCompactSerialization(compact);
        check.setDoKeyValidation(false);
        check.setKnownCriticalHeaders(crit == null ? new String[0] : crit);
        PublicKey pub = kp.getPublic();
        check.setKey(pub);
        if (!check.verifySignature()) {
            throw new IllegalStateException("jose4j failed to verify its own signature for " + name);
        }

        String[] parts = compact.split("\\.");
        byte[] signingInput;
        if (unencodedPayload) {
            byte[] header = parts[0].getBytes(StandardCharsets.US_ASCII);
            byte[] body = payload.getBytes(StandardCharsets.UTF_8);
            byte[] all = new byte[header.length + 1 + body.length];
            System.arraycopy(header, 0, all, 0, header.length);
            all[header.length] = (byte) '.';
            System.arraycopy(body, 0, all, header.length + 1, body.length);
            signingInput = all;
        } else {
            signingInput = (parts[0] + "." + parts[1]).getBytes(StandardCharsets.US_ASCII);
        }

        out.println(String.join("\t", "JWS", name, alg, hex(pub.getEncoded()), compact, hex(signingInput)));
    }

    public static void main(String[] args) throws Exception {
        out = new PrintStream(System.out, true, "UTF-8");

        KeyPair rsaKp = rsa(2048);
        emit("rs256", AlgorithmIdentifiers.RSA_USING_SHA256, rsaKp, "Hello JAdES", false, null);
        emit("rs384", AlgorithmIdentifiers.RSA_USING_SHA384, rsaKp, "Hello JAdES", false, null);
        emit("rs512", AlgorithmIdentifiers.RSA_USING_SHA512, rsaKp, "Hello JAdES", false, null);
        emit("ps256", AlgorithmIdentifiers.RSA_PSS_USING_SHA256, rsaKp, "Hello JAdES", false, null);
        emit("ps384", AlgorithmIdentifiers.RSA_PSS_USING_SHA384, rsaKp, "Hello JAdES", false, null);
        emit("ps512", AlgorithmIdentifiers.RSA_PSS_USING_SHA512, rsaKp, "Hello JAdES", false, null);
        emit("rs256-unencoded", AlgorithmIdentifiers.RSA_USING_SHA256, rsaKp,
                "detached payload with no period", true, new String[] { "b64" });
        emit("rs256-utf8", AlgorithmIdentifiers.RSA_USING_SHA256, rsaKp,
                "payload with non-ASCII: é中文", false, null);

        emit("es256", AlgorithmIdentifiers.ECDSA_USING_P256_CURVE_AND_SHA256, ec("secp256r1"),
                "Hello JAdES", false, null);
        emit("es384", AlgorithmIdentifiers.ECDSA_USING_P384_CURVE_AND_SHA384, ec("secp384r1"),
                "Hello JAdES", false, null);
        emit("es512", AlgorithmIdentifiers.ECDSA_USING_P521_CURVE_AND_SHA512, ec("secp521r1"),
                "Hello JAdES", false, null);

        emit("eddsa", AlgorithmIdentifiers.EDDSA, ed25519(), "Hello JAdES", false, null);
    }
}

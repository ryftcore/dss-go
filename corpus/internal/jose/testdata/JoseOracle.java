import org.jose4j.base64url.Base64Url;
import org.jose4j.json.JsonUtil;
import org.jose4j.json.internal.json_simple.JSONArray;
import org.jose4j.json.internal.json_simple.JSONValue;
import org.jose4j.json.internal.json_simple.parser.JSONParser;
import org.jose4j.jwx.CompactSerializer;

import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Dumps jose4j 0.9.6 ground truth for the Go internal/jose port.
 * Output is TSV; every payload column is hex of the UTF-8 bytes so control characters
 * never reach the terminal.
 */
public class JoseOracle {

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

    /** Builds a String from UTF-16 code units, so no control character appears in this source. */
    static String cu(int... units) {
        StringBuilder sb = new StringBuilder();
        for (int u : units) {
            sb.append((char) u);
        }
        return sb.toString();
    }

    static void row(String... cols) {
        out.println(String.join("\t", cols));
    }

    static void escapeCase(String name, String value) {
        Map<String, Object> m = new LinkedHashMap<String, Object>();
        m.put("k", value);
        row("ESC", name, hex(value), hex(JsonUtil.toJson(m)));
    }

    static void keyEscapeCase(String name, String key) {
        Map<String, Object> m = new LinkedHashMap<String, Object>();
        m.put(key, "v");
        row("ESCKEY", name, hex(key), hex(JsonUtil.toJson(m)));
    }

    static void hashOrderCase(String name, String... keys) {
        Map<String, Object> m = new HashMap<String, Object>();
        for (int i = 0; i < keys.length; i++) {
            m.put(keys[i], "v" + i);
        }
        row("HASHORDER", name, String.join(",", keys), hex(JsonUtil.toJson(m)));
    }

    public static void main(String[] args) throws Exception {
        out = new PrintStream(System.out, true, "UTF-8");

        // ===== escaping of string values =====
        escapeCase("plain", "hello");
        escapeCase("quote", cu(97, 34, 98));
        escapeCase("backslash", cu(97, 92, 98));
        escapeCase("slash", cu(97, 47, 98));
        escapeCase("backspace", cu(97, 8, 98));
        escapeCase("formfeed", cu(97, 12, 98));
        escapeCase("newline", cu(97, 10, 98));
        escapeCase("cr", cu(97, 13, 98));
        escapeCase("tab", cu(97, 9, 98));
        escapeCase("nul", cu(97, 0, 98));
        escapeCase("c0_01", cu(97, 1, 98));
        escapeCase("c0_1f", cu(97, 31, 98));
        escapeCase("space", cu(97, 32, 98));
        escapeCase("del7f", cu(97, 0x7f, 98));
        escapeCase("c1_80", cu(97, 0x80, 98));
        escapeCase("c1_9f", cu(97, 0x9f, 98));
        escapeCase("a0_nbsp", cu(97, 0xa0, 98));
        escapeCase("latin_e9", cu(99, 97, 102, 0xe9));
        escapeCase("u1fff", cu(97, 0x1fff, 98));
        escapeCase("u2000", cu(97, 0x2000, 98));
        escapeCase("u2028", cu(97, 0x2028, 98));
        escapeCase("u2029", cu(97, 0x2029, 98));
        escapeCase("u20ff", cu(97, 0x20ff, 98));
        escapeCase("u2100", cu(97, 0x2100, 98));
        escapeCase("cjk", cu(0x4e2d, 0x6587));
        escapeCase("emoji", cu(97, 0xd83d, 0xde00, 98));
        escapeCase("empty", "");
        escapeCase("mixed", cu(120, 34, 92, 47, 8, 12, 10, 13, 9, 0x7f, 0x2000, 0xe9, 121));

        // ===== escaping of keys =====
        keyEscapeCase("k_plain", "hello");
        keyEscapeCase("k_quote", cu(97, 34, 98));
        keyEscapeCase("k_ctrl", cu(97, 1, 98));
        keyEscapeCase("k_u2000", cu(97, 0x2000, 98));
        keyEscapeCase("k_hash", "x5t#S256");

        // ===== structure / value types =====
        {
            Map<String, Object> inner = new LinkedHashMap<String, Object>();
            inner.put("z", "1");
            inner.put("a", "2");

            Map<String, Object> m = new LinkedHashMap<String, Object>();
            m.put("s", "str");
            m.put("n", Long.valueOf(42L));
            m.put("i", Integer.valueOf(7));
            m.put("neg", Long.valueOf(-13L));
            m.put("big", Long.valueOf(9007199254740993L));
            m.put("t", Boolean.TRUE);
            m.put("f", Boolean.FALSE);
            m.put("nul", null);
            m.put("arr", new JSONArray(Arrays.asList("a", Long.valueOf(1L), Boolean.TRUE, null)));
            m.put("emptyArr", new JSONArray());
            m.put("obj", inner);
            m.put("emptyObj", new LinkedHashMap<String, Object>());
            m.put("strArr", new String[] { "a", "b" });
            m.put("listOfMaps", Arrays.asList(inner, inner));
            m.put("arrayList", new ArrayList<Object>(Arrays.asList("p", "q")));
            row("STRUCT", "kitchensink", "", hex(JsonUtil.toJson(m)));
        }
        row("STRUCT", "empty", "", hex(JsonUtil.toJson(new LinkedHashMap<String, Object>())));
        {
            Map<String, Object> m = new LinkedHashMap<String, Object>();
            m.put("zeta", "1");
            m.put("alpha", "2");
            m.put("mid", "3");
            row("STRUCT", "insertionorder", "", hex(JsonUtil.toJson(m)));
        }
        {
            Map<String, Object> m = new LinkedHashMap<String, Object>();
            m.put("a", "1");
            m.put("b", "2");
            m.put("a", "3");
            row("STRUCT", "reput", "", hex(JsonUtil.toJson(m)));
        }

        // ===== top-level non-map values (DSSJsonUtils.toBase64Url(Object)) =====
        row("VALUE", "array", "", hex(JSONValue.toJSONString(new JSONArray(Arrays.asList("a", "b")))));
        row("VALUE", "string", "", hex(JSONValue.toJSONString(cu(97, 34, 98))));
        row("VALUE", "null", "", hex(JSONValue.toJSONString(null)));
        row("VALUE", "long", "", hex(JSONValue.toJSONString(Long.valueOf(123L))));
        row("VALUE", "double", "", hex(JSONValue.toJSONString(Double.valueOf(1.5d))));
        row("VALUE", "bool", "", hex(JSONValue.toJSONString(Boolean.TRUE)));

        // ===== java HashMap iteration order (the no-arg JsonObject constructor) =====
        hashOrderCase("rVals", "crlVals", "ocspVals");
        hashOrderCase("rVals_rev", "ocspVals", "crlVals");
        hashOrderCase("tstVd", "xVals", "rVals");
        hashOrderCase("tstVd_rev", "rVals", "xVals");
        hashOrderCase("tstToken", "val");
        hashOrderCase("etsiU", "etsiU");
        hashOrderCase("three", "a", "b", "c");
        hashOrderCase("headers", "alg", "kid", "x5t#S256", "crit", "sigT", "b64", "typ", "cty");
        hashOrderCase("many12", "k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8", "k9", "k10", "k11", "k12");
        hashOrderCase("many13", "k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8", "k9", "k10", "k11", "k12", "k13");
        hashOrderCase("many20", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "a10",
                "a11", "a12", "a13", "a14", "a15", "a16", "a17", "a18", "a19", "a20");
        hashOrderCase("many40", "b1", "b2", "b3", "b4", "b5", "b6", "b7", "b8", "b9", "b10",
                "b11", "b12", "b13", "b14", "b15", "b16", "b17", "b18", "b19", "b20",
                "b21", "b22", "b23", "b24", "b25", "b26", "b27", "b28", "b29", "b30",
                "b31", "b32", "b33", "b34", "b35", "b36", "b37", "b38", "b39", "b40");
        hashOrderCase("collide", "Aa", "BB", "AaAa", "BBBB", "AaBB", "BBAa");
        hashOrderCase("jades", "sigT", "x5t#S256", "crit", "alg", "cty", "sigD", "b64", "adoTst", "srCms");

        // ===== base64url encode =====
        String[] b64inputs = { "", "f", "fo", "foo", "foob", "fooba", "foobar", "the quick brown fox" };
        for (String s : b64inputs) {
            byte[] b = s.getBytes(StandardCharsets.UTF_8);
            row("B64ENC", "s" + s.length(), hex(b), hex(Base64Url.encode(b)));
        }
        byte[][] rawInputs = {
            new byte[0],
            new byte[] { 0 },
            new byte[] { (byte) 0xff },
            new byte[] { (byte) 0xfb, (byte) 0xff },
            new byte[] { (byte) 0xfb, (byte) 0xef, (byte) 0xff },
            new byte[] { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 },
            new byte[] { (byte) 0xff, (byte) 0xef, (byte) 0xfe },
        };
        for (int i = 0; i < rawInputs.length; i++) {
            row("B64ENC", "raw" + i, hex(rawInputs[i]), hex(Base64Url.encode(rawInputs[i])));
        }

        // ===== base64url decode (including the lenient/odd inputs) =====
        String[] decInputs = { "", "Zg", "Zm8", "Zm9v", "Zg==", "Zm8=", "-_8", "_-8", "AQID", "AQI",
            "AQ", "A", "Zm9vYg", "+/8=", "Zm9vYmFy", "Zm 9v", "Zm9v.", "!!!!" };
        for (String s : decInputs) {
            String res;
            try {
                res = hex(Base64Url.decode(s));
            } catch (Exception e) {
                res = "ERR:" + e.getClass().getSimpleName();
            }
            row("B64DEC", "d", hex(s), res);
        }
        row("B64DEC", "nl", hex(cu(90, 109, 57, 118, 10)), hex(Base64Url.decode(cu(90, 109, 57, 118, 10))));

        // ===== compact serializer =====
        row("CSSER", "3parts", "", hex(CompactSerializer.serialize("a", "b", "c")));
        row("CSSER", "emptymid", "", hex(CompactSerializer.serialize("a", "", "c")));
        row("CSSER", "nullmid", "", hex(CompactSerializer.serialize("a", null, "c")));
        row("CSSER", "2parts", "", hex(CompactSerializer.serialize("a", "b")));
        row("CSSER", "1part", "", hex(CompactSerializer.serialize("a")));
        row("CSSER", "0parts", "", hex(CompactSerializer.serialize()));
        row("CSSER", "trailing", "", hex(CompactSerializer.serialize("a", "b", "")));

        String[] deserInputs = { "a.b.c", "a..c", "a.b.", "a.b", "a", "", ".", "..", "a...", "...a", "a.b.c.d" };
        for (String s : deserInputs) {
            String[] parts = CompactSerializer.deserialize(s);
            List<String> hexParts = new ArrayList<String>();
            for (String p : parts) {
                hexParts.add(hex(p));
            }
            row("CSDESER", "n" + parts.length, hex(s), String.join("|", hexParts));
        }

        // ===== parse then re-serialize (number and format fidelity) =====
        String[] parseCases = {
            "{\"a\":1}", "{\"a\":-0}", "{\"a\":0}", "{\"a\":1.5}", "{\"a\":1e3}", "{\"a\":1E3}",
            "{\"a\":1.0}", "{\"a\":-1.5e-7}", "{\"a\":123456789012345678901234567890}",
            "{\"a\":9223372036854775807}", "{\"a\":-9223372036854775808}",
            "{\"a\":1607000000}", "{\"a\":0.1}", "{\"a\":100000000.0}", "{\"a\":1e7}", "{\"a\":1e-3}",
            "{\"a\":1e-4}", "{\"a\":12345678901234567890.5}", "{\"a\":1e21}",
            "{\"b\":\"x\",\"a\":\"y\"}", "{\"a\":[1,2,{\"c\":true}]}", "{\"a\":null}",
            "{\"a\":\"\\u0041\\n\\/\"}", "{ \"a\" : 1 , \"b\" : 2 }", "{\"a\":\"\\ud83d\\ude00\"}",
            "{\"a\":\"\\u2000\"}", "{\"a\":\"\\u007f\"}"
        };
        for (String s : parseCases) {
            String res;
            try {
                res = hex(JsonUtil.toJson(JsonUtil.parseJson(s)));
            } catch (Exception e) {
                res = "ERR:" + e.getClass().getSimpleName();
            }
            row("PARSE", "p", hex(s), res);
        }
        {
            String s = "{\"a\":1,\"a\":2}";
            String res;
            try {
                res = hex(JsonUtil.toJson(JsonUtil.parseJson(s)));
            } catch (Exception e) {
                res = "ERR:" + e.getClass().getSimpleName();
            }
            row("PARSE", "dupe", hex(s), res);
        }

        String[] badCases = { "", "{", "not json", "[]", "\"x\"", "null", "{\"a\"}", "{'a':1}", "{\"a\":1}}",
            "{\"a\":01}", "{\"a\":+1}", "{\"a\":.5}", "{\"a\":1.}", "{\"a\":Infinity}", "{\"a\":NaN}",
            "{\"a\":\"\\x\"}", "{\"a\":\"\\u00zz\"}", "{\"a\":1,}", "{,}", "{\"a\":1 \"b\":2}",
            "  {\"a\":1}  ", "{\"a\":1}\n", "{\"a\":\"unterminated}", "{\"a\":00}", "{\"a\":1e}",
            "{\"a\":-}", "{\"a\":tru}", "{\"a\":TRUE}", "{\"\":1}", "{\"a\":\"\t\"}", "{\"a\":0.5e+3}" };
        for (String s : badCases) {
            String res;
            try {
                res = "OK:" + hex(JsonUtil.toJson(JsonUtil.parseJson(s)));
            } catch (Exception e) {
                res = "ERR:" + e.getClass().getSimpleName();
            }
            row("PARSEBAD", "b", hex(s), res);
        }

        // ===== Double.toString parity =====
        double[] ds = { 0.0d, -0.0d, 1.0d, 1.5d, 0.1d, 1e3, 1e7, 1e-3, 1e-4, 1.0e20, 1.0e-20,
            123456789.0d, 3.141592653589793d, -2.5e-8, 1e21, 4.9e-324, 1.7976931348623157e308,
            1234567.0d, 12345678.0d, 0.001d, 0.0001d };
        for (double d : ds) {
            row("DBL", "d", Double.toHexString(d), hex(Double.toString(d)));
        }

        // ===== JSONParser without a container factory =====
        String[] anyCases = { "{\"a\":1}", "[1,2,3]", "\"str\"", "123", "true", "null", "{\"b\":1,\"a\":2}" };
        for (String s : anyCases) {
            String res;
            try {
                Object o = new JSONParser().parse(s);
                res = (o == null ? "null" : o.getClass().getSimpleName()) + ":" + hex(JSONValue.toJSONString(o));
            } catch (Exception e) {
                res = "ERR:" + e.getClass().getSimpleName();
            }
            row("ANY", "a", hex(s), res);
        }
    }
}

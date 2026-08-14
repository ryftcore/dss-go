// Ported from org.jose4j.json.internal.json_simple.JSONValue (writeJSONString, escape),
// .JSONObject (writeJSONString(Map, Writer)), .JSONArray (writeJSONString(Collection, Writer))
// and org.jose4j.json.JsonUtil.toJson (jose4j 0.9.6).
package jose

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// JSON renders value as JSON exactly as jose4j does. Port of JsonUtil.toJson(Map) and, for a
// non-map argument, of JSONValue.toJSONString(Object) - DSSJsonUtils.toBase64Url(Object) calls
// the latter on bare JSON arrays.
//
// The output of this function is what gets base64url-encoded and signed. Do not reimplement it
// with encoding/json: the two differ on '<', '>', '&', '/' and on which code points are escaped.
func JSON(value any) string {
	var sb strings.Builder
	writeJSONValue(&sb, value)
	return sb.String()
}

// writeJSONValue is JSONValue.writeJSONString(Object, Writer). The type dispatch order is
// jose4j's, and it ends the same way: anything unrecognised is rendered as the JSON string of
// its Java toString(), never dropped and never an error.
func writeJSONValue(sb *strings.Builder, value any) {
	if value == nil {
		sb.WriteString("null")
		return
	}
	switch v := value.(type) {
	case string:
		writeJSONString(sb, v)
	case *Number:
		writeNumber(sb, v)
	case float64:
		// java.lang.Double: infinities and NaN are written as null, not as a JSON number.
		if math.IsInf(v, 0) || math.IsNaN(v) {
			sb.WriteString("null")
			return
		}
		sb.WriteString(JavaDoubleToString(v))
	case float32:
		f := float64(v)
		if math.IsInf(f, 0) || math.IsNaN(f) {
			sb.WriteString("null")
			return
		}
		sb.WriteString(JavaDoubleToString(f))
	case int:
		sb.WriteString(strconv.FormatInt(int64(v), 10))
	case int32:
		sb.WriteString(strconv.FormatInt(int64(v), 10))
	case int64:
		sb.WriteString(strconv.FormatInt(v, 10))
	case bool:
		if v {
			sb.WriteString("true")
			return
		}
		sb.WriteString("false")
	case *Object:
		writeJSONObject(sb, v)
	case ObjectHolder:
		// The `value instanceof Map` branch: a wrapper type that IS a JSON object renders as
		// one, not through the toString() fallback. eu.europa.esig.dss.jades.JsonObject reaches
		// the writer this way, since it implements java.util.Map upstream.
		writeJSONObject(sb, v.JSONObject())
	case []any:
		writeJSONArray(sb, v)
	case []string:
		// java.util.List<String> and String[] both land in the array branches; DSS builds the
		// 'crit' header and the x5c chain this way.
		items := make([]any, len(v))
		for i, s := range v {
			items[i] = s
		}
		writeJSONArray(sb, items)
	case []*Object:
		items := make([]any, len(v))
		for i, o := range v {
			items[i] = o
		}
		writeJSONArray(sb, items)
	default:
		// JSONValue's fallback: writeJSONString(value.toString(), out).
		writeJSONString(sb, fmt.Sprint(value))
	}
}

// writeNumber renders a *Number, honouring the Double/null rule for infinities and NaN.
func writeNumber(sb *strings.Builder, n *Number) {
	if n == nil {
		sb.WriteString("null")
		return
	}
	if n.Kind() == NumberDouble {
		f := n.Float64()
		if math.IsInf(f, 0) || math.IsNaN(f) {
			sb.WriteString("null")
			return
		}
	}
	sb.WriteString(n.String())
}

// writeJSONObject is JSONObject.writeJSONString(Map, Writer): '{', then each member as
// "key":value separated by ',', then '}'. No whitespace anywhere - a single space would change
// the signature.
func writeJSONObject(sb *strings.Builder, o *Object) {
	if o == nil {
		sb.WriteString("null")
		return
	}
	sb.WriteByte('{')
	for i, k := range o.Keys() {
		if i > 0 {
			sb.WriteByte(',')
		}
		writeJSONString(sb, k)
		sb.WriteByte(':')
		writeJSONValue(sb, o.values[k])
	}
	sb.WriteByte('}')
}

// writeJSONArray is JSONArray.writeJSONString(Collection, Writer).
func writeJSONArray(sb *strings.Builder, items []any) {
	if items == nil {
		sb.WriteString("null")
		return
	}
	sb.WriteByte('[')
	for i, item := range items {
		if i > 0 {
			sb.WriteByte(',')
		}
		writeJSONValue(sb, item)
	}
	sb.WriteByte(']')
}

// writeJSONString writes s as a quoted, escaped JSON string. Port of JSONValue.escape(String,
// StringBuilder) wrapped in the quotes its callers add.
func writeJSONString(sb *strings.Builder, s string) {
	sb.WriteByte('"')
	sb.WriteString(EscapeJSONString(s))
	sb.WriteByte('"')
}

// EscapeJSONString applies json-simple's escaping rules, as jose4j ships them. Every clause below
// is a deliberate difference from encoding/json, and each one is pinned by TestKnownAnswers:
//
//   - '"' and '\' take their two-character escapes, as does each of \b \f \n \r \t.
//   - '/' is NOT escaped. Stock json-simple writes "\/"; jose4j patched that out ("It also will
//     let me to make changes to the JSON processing like not escaping forward slashes").
//   - '<', '>' and '&' are NOT escaped. encoding/json escapes all three unless HTML escaping is
//     switched off, which would silently corrupt any header carrying a URL with a query string.
//   - U+0000-U+001F, U+007F-U+009F and U+2000-U+20FF become \uXXXX with UPPERCASE hex digits,
//     zero-padded to four. The last range is json-simple's own idea (it is the General
//     Punctuation block, escaped so that U+2028/U+2029 cannot break JavaScript eval) and it has
//     no counterpart in encoding/json at all.
//   - Everything else - including U+00A0, U+1FFF, U+2100, CJK and astral characters - is emitted
//     as literal UTF-8.
//
// The ranges are checked against UTF-16 code units in Java. For a Go rune that only matters above
// U+FFFF, where the surrogates (U+D800-U+DFFF) fall outside every escaped range, so a supplementary
// character is emitted literally either way.
func EscapeJSONString(s string) string {
	if !needsJSONEscaping(s) {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + 8)
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\b':
			sb.WriteString(`\b`)
		case '\f':
			sb.WriteString(`\f`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		default:
			if isJSONSimpleEscapedRune(r) {
				sb.WriteString(`\u`)
				sb.WriteString(upperHex4(uint32(r)))
				continue
			}
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// needsJSONEscaping is a fast path so that the common case - an all-ASCII base64url value or a
// header name - is returned unchanged instead of being rebuilt rune by rune.
func needsJSONEscaping(s string) bool {
	for _, r := range s {
		switch r {
		case '"', '\\', '\b', '\f', '\n', '\r', '\t':
			return true
		}
		if isJSONSimpleEscapedRune(r) {
			return true
		}
	}
	return false
}

// isJSONSimpleEscapedRune reports whether r falls in one of json-simple's three \u-escaped
// ranges. The two-character escapes are handled by the caller and are not included here.
func isJSONSimpleEscapedRune(r rune) bool {
	return (r >= 0x0000 && r <= 0x001F) ||
		(r >= 0x007F && r <= 0x009F) ||
		(r >= 0x2000 && r <= 0x20FF)
}

// upperHex4 formats v as exactly four uppercase hex digits, which is what
// Integer.toHexString(ch) padded to four and upper-cased produces.
func upperHex4(v uint32) string {
	const digits = "0123456789ABCDEF"
	return string([]byte{
		digits[(v>>12)&0xf],
		digits[(v>>8)&0xf],
		digits[(v>>4)&0xf],
		digits[v&0xf],
	})
}

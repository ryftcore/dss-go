// Ported from org.jose4j.jwx.CompactSerializer (jose4j 0.9.6).
package jose

import "strings"

// CompactSerialize joins parts with '.'. Port of CompactSerializer.serialize(String...), which
// DSSJsonUtils.concatenate re-exports and which builds both the JWS signing input
// (header + '.' + payload) and the compact serialization (+ '.' + signature).
//
// Go has no null string, so a caller that wants jose4j's `(parts[i] == null) ? "" : parts[i]`
// simply passes "": the two are already the same value here.
func CompactSerialize(parts ...string) string {
	return strings.Join(parts, ".")
}

// CompactDeserialize splits a compact serialization on '.'. Port of
// CompactSerializer.deserialize(String).
//
// The implementation is `compactSerialization.split("\\.")` plus a hand-written fix-up, and both
// halves have to be reproduced because JWSCompactSerializationParser feeds the result straight
// into JWS's three-part constructor, so the part count decides whether a document is accepted:
//
//   - java.lang.String.split with the default limit DROPS trailing empty fields, so "a.b."
//     splits to just ["a", "b"], and "." and ".." split to nothing at all;
//   - jose4j then appends one empty part when the input ends in '.', which restores the
//     detached-payload form "header..signature" to three parts and turns "." into [""].
//
// The two rules do not cancel out: "a..." yields ["a", ""], not ["a", "", "", ""].
func CompactDeserialize(compactSerialization string) []string {
	parts := javaSplitOnPeriod(compactSerialization)
	if strings.HasSuffix(compactSerialization, ".") {
		parts = append(parts, "")
	}
	return parts
}

// javaSplitOnPeriod is java.lang.String.split("\\.") with limit 0: split on every '.', then
// remove the trailing empty fields. An input with no '.' at all yields the input itself, even
// when it is empty - which is why "" deserializes to one part rather than none.
func javaSplitOnPeriod(s string) []string {
	if !strings.Contains(s, ".") {
		return []string{s}
	}
	fields := strings.Split(s, ".")
	end := len(fields)
	for end > 0 && fields[end-1] == "" {
		end--
	}
	return fields[:end]
}

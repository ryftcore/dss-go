// Package jose provides the JOSE (RFC 7515/7518) machinery that upstream DSS obtains from the
// third-party org.jose4j library. Exactly as PORTING.md prescribes for the BouncyCastle
// replacements, there is no Java class here to mirror one-to-one, so the package lives under
// internal/ instead of being duplicated inside dss/jades. It imports the standard library and
// nothing else - not even another DSS package, since dss/jades is what imports it.
//
// # Provenance
//
// Everything here is ported from org.jose4j 0.9.6, the version dss-jades 6.5.RC1 declares
// (sd-dss/pom.xml, <jose4j.version>0.9.6</jose4j.version>). The mapping is:
//
//	org.jose4j.base64url.Base64Url                        -> base64url.go
//	  + org.jose4j.base64url.internal.apache.commons
//	    .codec.binary.Base64 (the bundled decoder)
//	org.jose4j.lang.StringUtil                            -> base64url.go (ASCIIBytes/UTF8String)
//	org.jose4j.json.internal.json_simple.JSONValue        -> writer.go (JSON, WriteJSON, escape)
//	org.jose4j.json.internal.json_simple.JSONObject       -> writer.go (object rendering)
//	org.jose4j.json.internal.json_simple.JSONArray        -> writer.go (array rendering)
//	org.jose4j.json.JsonUtil                              -> parser.go (ParseJSON), writer.go (JSON)
//	org.jose4j.json.internal.json_simple.parser.JSONParser
//	  + .Yylex                                            -> parser.go
//	java.util.LinkedHashMap / java.util.HashMap           -> object.go, javahash.go
//	java.lang.Long/Double/BigInteger toString             -> number.go
//	org.jose4j.jwx.CompactSerializer                      -> compact.go
//	org.jose4j.jwx.HeaderParameterNames                   -> header_parameter_names.go
//	org.jose4j.jwx.Headers                                -> headers.go
//	org.jose4j.jwx.JsonWebStructure
//	  + org.jose4j.jws.JsonWebSignature                   -> jws.go
//	org.jose4j.jws.AlgorithmIdentifiers and the
//	  JsonWebSignatureAlgorithm implementations           -> sigalg.go
//
// Only what dss-jades actually consumes is ported. JWE, JWK, the JWT consumer, the HTTPS JWKS
// client and the signing half of the algorithm factory have no call site in dss-jades and are
// deliberately absent; adding them later means porting them, not un-commenting them.
//
// # Why byte-exactness is the whole point
//
// A JWS signs ASCII(BASE64URL(UTF8(protected header))), i.e. it signs the *serialized bytes* of
// a JSON object. Two JSON documents that are equal as data but differ by one byte - a different
// member order, a differently escaped character - have different signatures. jose4j serializes
// LinkedHashMap-backed headers, so member order is insertion order, which is the order of the
// Java statements that built the header. A Go map cannot reproduce that (its iteration order is
// randomized on purpose), so Object here is insertion-ordered, and the writer reproduces
// json-simple's escaping rules character for character rather than using encoding/json:
//
//   - '/' is NOT escaped (jose4j patched json-simple to stop escaping it);
//   - '<', '>' and '&' are NOT escaped (encoding/json escapes all three by default);
//   - U+0000-U+001F, U+007F-U+009F and U+2000-U+20FF are escaped as \uXXXX with UPPERCASE hex
//     (encoding/json escapes a different set, in lowercase);
//   - everything else, including U+2028/U+2029 outside those ranges and all astral characters,
//     is emitted as literal UTF-8.
//
// Nothing in this package may be re-routed through encoding/json for convenience. TestKnownAnswers
// in writer_test.go pins every rule above against testdata/jose4j_oracle.tsv, which was produced
// by running jose4j 0.9.6 itself (the generator is recorded in testdata/README.md).
package jose

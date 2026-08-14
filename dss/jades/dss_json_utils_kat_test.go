package jades

import (
	"bufio"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/jose"
	"github.com/utain/esig/dss/model"
)

// jadesJSONKATCase is one row of testdata/jades_json_oracle.tsv, which was produced by running
// dss-jades 6.5.RC1 itself (see testdata/JadesOracle.java). Payload columns are hex of UTF-8
// bytes; a few output columns are plain "true"/"false" or a '|'-separated tuple, and the code
// below knows which.
type jadesJSONKATCase struct {
	line    int
	section string
	name    string
	in      []byte
	rawOut  string
	out     []byte
}

func jadesJSONKATLoad(t *testing.T) []jadesJSONKATCase {
	t.Helper()
	f, err := os.Open("testdata/jades_json_oracle.tsv")
	if err != nil {
		t.Fatalf("cannot open the dss-jades oracle: %v", err)
	}
	defer f.Close()

	var cases []jadesJSONKATCase
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<22)
	for line := 1; scanner.Scan(); line++ {
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			continue
		}
		cols := strings.Split(text, "\t")
		if len(cols) != 4 {
			t.Fatalf("line %d: expected 4 columns, got %d", line, len(cols))
		}
		in, _ := hex.DecodeString(cols[2])
		out, _ := hex.DecodeString(cols[3])
		cases = append(cases, jadesJSONKATCase{
			line:    line,
			section: cols[0],
			name:    cols[1],
			in:      in,
			rawOut:  cols[3],
			out:     out,
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading the dss-jades oracle: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("the dss-jades oracle is empty")
	}
	return cases
}

func jadesJSONKATSection(t *testing.T, cases []jadesJSONKATCase, section string) []jadesJSONKATCase {
	t.Helper()
	var out []jadesJSONKATCase
	for _, c := range cases {
		if c.section == section {
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		t.Fatalf("the dss-jades oracle has no %s cases", section)
	}
	return out
}

// jadesJSONKATByName indexes a section by case name, for the sections whose input is not in the
// row (the case is built in Java code instead).
func jadesJSONKATByName(t *testing.T, cases []jadesJSONKATCase, section string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	for _, c := range jadesJSONKATSection(t, cases, section) {
		out[c.name] = c.out
	}
	return out
}

// TestKnownAnswersCompactParsing pins the compact parser against dss-jades: whether a document is
// accepted, and what signing input, signed payload and encoded header it yields. The signing
// input is the byte string a JAdES signature is computed over, so it is the assertion that
// matters most here.
func TestKnownAnswersCompactParsing(t *testing.T) {
	cases := jadesJSONKATLoad(t)
	compactByName := jadesJSONKATByName(t, cases, "COMPACT")

	for _, section := range []string{"SIGNING_INPUT", "SIGNED_PAYLOAD", "ENCODED_HEADER"} {
		for _, c := range jadesJSONKATSection(t, cases, section) {
			t.Run(section+"/"+c.name, func(t *testing.T) {
				jws := jadesJSONKATParseCompact(t, c.in)
				var got []byte
				switch section {
				case "SIGNING_INPUT":
					got = DSSJsonUtilsSigningInputBytes(jws)
				case "SIGNED_PAYLOAD":
					got = []byte(jws.SignedPayload())
				case "ENCODED_HEADER":
					got = []byte(jws.EncodedHeader())
				}
				if string(got) != string(c.out) {
					t.Errorf("\n got %q\nwant %q", got, c.out)
				}
			})
		}
	}
	// Every compact fixture must round-trip through the parser without loss.
	for name, compact := range compactByName {
		t.Run("roundtrip/"+name, func(t *testing.T) {
			jws := jadesJSONKATParseCompact(t, compact)
			rebuilt := DSSJsonUtilsConcatenate(jws.EncodedHeader(), jws.SignedPayload(), jws.EncodedSignature())
			if rebuilt != string(compact) {
				t.Errorf("\n got %q\nwant %q", rebuilt, compact)
			}
		})
	}
}

// TestKnownAnswersCompactIsSupported pins the format sniffer, including the accepted trailing
// line break and the rejected trailing content after it.
func TestKnownAnswersCompactIsSupported(t *testing.T) {
	for _, c := range jadesJSONKATSection(t, jadesJSONKATLoad(t), "COMPACT_SUPPORTED") {
		t.Run(c.name+"/"+string(c.in[:min(len(c.in), 20)]), func(t *testing.T) {
			parser := NewJWSCompactSerializationParser(model.NewInMemoryDocument(c.in))
			got, err := parser.IsSupported()
			if err != nil {
				t.Fatalf("IsSupported: %v", err)
			}
			if want := c.rawOut == "true"; got != want {
				t.Errorf("IsSupported(%q) = %v, want %v", c.in, got, want)
			}
		})
	}
}

// TestKnownAnswersJSONIsSupported pins isJsonDocument, whose leading-'{' test is what keeps a
// compact JWS out of the JSON parser.
func TestKnownAnswersJSONIsSupported(t *testing.T) {
	for _, c := range jadesJSONKATSection(t, jadesJSONKATLoad(t), "JSON_SUPPORTED") {
		t.Run(string(c.in[:min(len(c.in), 20)]), func(t *testing.T) {
			parser := NewJWSJsonSerializationParser(model.NewInMemoryDocument(c.in))
			got, err := parser.IsSupported()
			if err != nil {
				t.Fatalf("IsSupported: %v", err)
			}
			if want := c.rawOut == "true"; got != want {
				t.Errorf("IsSupported(%q) = %v, want %v", c.in, got, want)
			}
		})
	}
}

// TestKnownAnswersSerializationGenerator pins the generated signature documents byte for byte.
// Member order here is not decorative: an archive timestamp is computed over these octets.
func TestKnownAnswersSerializationGenerator(t *testing.T) {
	cases := jadesJSONKATLoad(t)
	for _, section := range []struct {
		name   string
		output enumerations.JWSSerializationType
	}{
		{"GEN_FLATTENED", enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION},
		{"GEN_COMPLETE", enumerations.JWSSerializationType_JSON_SERIALIZATION},
	} {
		for _, c := range jadesJSONKATSection(t, cases, section.name) {
			t.Run(section.name+"/"+c.name, func(t *testing.T) {
				jws := jadesJSONKATParseCompact(t, c.in)
				object := DSSJsonUtilsToJWSJsonSerializationObject(jws)
				document, err := NewJWSJsonSerializationGenerator(object, section.output).Generate()
				if err != nil {
					t.Fatalf("Generate: %v", err)
				}
				jadesJSONKATAssertBytes(t, document, c.out)
			})
		}
	}
}

// TestKnownAnswersConverter pins JWSConverter's four public conversions, including the document
// name and mime type it stamps on the result.
func TestKnownAnswersConverter(t *testing.T) {
	cases := jadesJSONKATLoad(t)

	for _, c := range jadesJSONKATSection(t, cases, "CONV_FLATTENED") {
		t.Run("compactToFlattened/"+c.name, func(t *testing.T) {
			document, err := JWSConverterFromJWSCompactToJSONFlattenedSerialization(model.NewInMemoryDocument(c.in))
			if err != nil {
				t.Fatalf("conversion: %v", err)
			}
			jadesJSONKATAssertBytes(t, document, c.out)
			if document.Name() != "json-flattened-serialization.json" {
				t.Errorf("document name = %q", document.Name())
			}
			if document.MimeType() != enumerations.MimeTypeEnum_JSON {
				t.Errorf("mime type = %v", document.MimeType())
			}
		})
	}
	for _, c := range jadesJSONKATSection(t, cases, "CONV_COMPLETE") {
		t.Run("compactToComplete/"+c.name, func(t *testing.T) {
			document, err := JWSConverterFromJWSCompactToJSONSerialization(model.NewInMemoryDocument(c.in))
			if err != nil {
				t.Fatalf("conversion: %v", err)
			}
			jadesJSONKATAssertBytes(t, document, c.out)
			if document.Name() != "json-serialization.json" {
				t.Errorf("document name = %q", document.Name())
			}
		})
	}
	for _, c := range jadesJSONKATSection(t, cases, "ETSIU_TO_CLEAR") {
		t.Run("etsiUToClear", func(t *testing.T) {
			document, err := JWSConverterFromEtsiUWithBase64UrlToClearJSONIncorporation(model.NewInMemoryDocument(c.in))
			if err != nil {
				t.Fatalf("conversion: %v", err)
			}
			jadesJSONKATAssertBytes(t, document, c.out)
			if document.Name() != "etsiU-clear-incorporation.json" {
				t.Errorf("document name = %q", document.Name())
			}
		})
	}
	for _, c := range jadesJSONKATSection(t, cases, "ETSIU_TO_B64") {
		t.Run("etsiUToBase64Url", func(t *testing.T) {
			document, err := JWSConverterFromEtsiUWithClearJSONToBase64UrlIncorporation(model.NewInMemoryDocument(c.in))
			if err != nil {
				t.Fatalf("conversion: %v", err)
			}
			jadesJSONKATAssertBytes(t, document, c.out)
			if document.Name() != "etsiU-base64url-incorporation.json" {
				t.Errorf("document name = %q", document.Name())
			}
		})
	}
}

// TestKnownAnswersSerializationParser pins what the parser makes of a complete JWS JSON
// Serialization: the serialization type, the signature count, the shared payload, the encoded
// header of the first signature and - crucially - the signing input recovered from it, which
// must be the same string the compact form produced.
func TestKnownAnswersSerializationParser(t *testing.T) {
	for _, c := range jadesJSONKATSection(t, jadesJSONKATLoad(t), "REPARSE") {
		t.Run(c.name, func(t *testing.T) {
			fields := strings.Split(c.rawOut, "|")
			if len(fields) != 5 {
				t.Fatalf("oracle line %d: expected 5 REPARSE fields, got %d", c.line, len(fields))
			}
			wantType := fields[0]
			wantCount := fields[1]
			wantPayload := jadesJSONKATUnhex(t, fields[2])
			wantHeader := jadesJSONKATUnhex(t, fields[3])
			wantSigningInput := jadesJSONKATUnhex(t, fields[4])

			object, err := NewJWSJsonSerializationParser(model.NewInMemoryDocument(c.in)).Parse()
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if string(object.JWSSerializationType()) != wantType {
				t.Errorf("serialization type = %q, want %q", object.JWSSerializationType(), wantType)
			}
			if got := len(object.Signatures()); strconv.Itoa(got) != wantCount {
				t.Fatalf("signature count = %d, want %s", got, wantCount)
			}
			if object.Payload() != string(wantPayload) {
				t.Errorf("payload = %q, want %q", object.Payload(), wantPayload)
			}
			signature := object.Signatures()[0]
			if signature.EncodedHeader() != string(wantHeader) {
				t.Errorf("encoded header = %q, want %q", signature.EncodedHeader(), wantHeader)
			}
			if got := DSSJsonUtilsSigningInputBytes(signature); string(got) != string(wantSigningInput) {
				t.Errorf("signing input\n got %q\nwant %q", got, wantSigningInput)
			}
		})
	}
}

// TestKnownAnswersOidObject pins the 'oid' data type of TS 119 182-1 clause 5.4.1, whose member
// order is the order of the Java statements that build it.
func TestKnownAnswersOidObject(t *testing.T) {
	byName := jadesJSONKATByName(t, jadesJSONKATLoad(t), "OID")

	commitment, err := DSSJsonUtilsOidObject(enumerations.CommitmentTypeEnum_ProofOfOrigin)
	if err != nil {
		t.Fatalf("DSSJsonUtilsOidObject: %v", err)
	}
	tests := map[string]*JsonObject{
		"commitment":    commitment,
		"uri-only":      DSSJsonUtilsOidObjectFromURI("http://x/y", "", nil),
		"uri-desc":      DSSJsonUtilsOidObjectFromURI("http://x/y", "a desc", nil),
		"uri-desc-refs": DSSJsonUtilsOidObjectFromURI("http://x/y", "a desc", []string{"r1", "r2"}),
	}
	for name, object := range tests {
		want, ok := byName[name]
		if !ok {
			t.Fatalf("the oracle has no OID case %q", name)
		}
		t.Run(name, func(t *testing.T) {
			if got := object.ToJSONString(); got != string(want) {
				t.Errorf("\n got %s\nwant %s", got, want)
			}
		})
	}
}

// TestKnownAnswersTstContainer pins the 'tstContainer' type of TS 119 182-1 clause 5.4.3.3,
// including that 'canonAlg' precedes 'tstTokens' when present and is absent otherwise.
func TestKnownAnswersTstContainer(t *testing.T) {
	byName := jadesJSONKATByName(t, jadesJSONKATLoad(t), "TSTCONTAINER")
	binaries := []*model.TimestampBinary{
		model.NewTimestampBinary([]byte{1, 2, 3}),
		model.NewTimestampBinary([]byte{0xff, 0}),
	}

	noCanon, err := DSSJsonUtilsTstContainer(binaries, "")
	if err != nil {
		t.Fatalf("DSSJsonUtilsTstContainer: %v", err)
	}
	if got, want := noCanon.ToJSONString(), string(byName["no-canon"]); got != want {
		t.Errorf("no-canon\n got %s\nwant %s", got, want)
	}
	canon, err := DSSJsonUtilsTstContainer(binaries, "http://www.w3.org/2001/10/xml-exc-c14n#")
	if err != nil {
		t.Fatalf("DSSJsonUtilsTstContainer: %v", err)
	}
	if got, want := canon.ToJSONString(), string(byName["canon"]); got != want {
		t.Errorf("canon\n got %s\nwant %s", got, want)
	}
	if _, err := DSSJsonUtilsTstContainer(nil, ""); err == nil {
		t.Error("an empty timestamp list was accepted")
	}
}

// TestKnownAnswersJsonObjectHashOrder is the one that would catch a JsonObject silently
// switched to insertion order: upstream's no-argument constructor wraps a java.util.HashMap, and
// these two objects are the ones DSS really builds that way ('rVals' and 'tstVd').
func TestKnownAnswersJsonObjectHashOrder(t *testing.T) {
	byName := jadesJSONKATByName(t, jadesJSONKATLoad(t), "JSONOBJECT")

	rVals := NewJsonObject()
	rVals.Put("crlVals", "c")
	rVals.Put("ocspVals", "o")
	if got, want := rVals.ToJSONString(), string(byName["rVals"]); got != want {
		t.Errorf("rVals\n got %s\nwant %s", got, want)
	}

	tstVd := NewJsonObject()
	tstVd.Put("xVals", "x")
	tstVd.Put("rVals", "r")
	if got, want := tstVd.ToJSONString(), string(byName["tstVd"]); got != want {
		t.Errorf("tstVd\n got %s\nwant %s", got, want)
	}
}

// TestKnownAnswersScalarHelpers pins the small predicates and converters, several of which decide
// whether a document is treated as a JAdES signature at all.
func TestKnownAnswersScalarHelpers(t *testing.T) {
	cases := jadesJSONKATLoad(t)

	for _, c := range jadesJSONKATSection(t, cases, "IS_B64URL") {
		t.Run("isBase64UrlEncoded/"+string(c.in), func(t *testing.T) {
			if got, want := DSSJsonUtilsIsBase64UrlEncoded(string(c.in)), c.rawOut == "true"; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
	for _, c := range jadesJSONKATSection(t, cases, "IS_URLSAFE_PAYLOAD") {
		t.Run("isUrlSafePayload/"+string(c.in), func(t *testing.T) {
			if got, want := DSSJsonUtilsIsUrlSafePayload(string(c.in)), c.rawOut == "true"; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
	for _, c := range jadesJSONKATSection(t, cases, "MIMETYPE") {
		t.Run("mimeType/"+string(c.in), func(t *testing.T) {
			if got := DSSJsonUtilsMimeTypeString(string(c.in)); got != string(c.out) {
				t.Errorf("got %q, want %q", got, c.out)
			}
		})
	}
	concat := jadesJSONKATByName(t, cases, "CONCAT")
	if got, want := DSSJsonUtilsConcatenate("a", "b", "c"), string(concat["three"]); got != want {
		t.Errorf("concatenate: got %q, want %q", got, want)
	}
	if got, want := DSSJsonUtilsConcatenate("a", "", "c"), string(concat["empty-mid"]); got != want {
		t.Errorf("concatenate: got %q, want %q", got, want)
	}

	b64Object := jadesJSONKATByName(t, cases, "B64URL_OBJECT")
	if got, want := DSSJsonUtilsToBase64UrlObject([]any{"a", "b"}), string(b64Object["array"]); got != want {
		t.Errorf("toBase64Url(array): got %q, want %q", got, want)
	}
	object := jose.NewObject()
	object.Put("b", "2")
	object.Put("a", "1")
	if got, want := DSSJsonUtilsToBase64UrlObject(NewJsonObjectFromMap(object)), string(b64Object["object"]); got != want {
		t.Errorf("toBase64Url(object): got %q, want %q", got, want)
	}
}

// TestCriticalHeaderSets covers the three static sets, whose membership decides whether a
// baseline requirement is met.
func TestCriticalHeaderSets(t *testing.T) {
	if !DSSJsonUtilsIsCriticalHeaderException(jose.HeaderAlgorithm) {
		t.Error("'alg' is not reported as a critical header exception")
	}
	if DSSJsonUtilsIsCriticalHeaderException(JAdESHeaderParameterNamesSigT) {
		t.Error("'sigT' is reported as a critical header exception")
	}
	if !DSSJsonUtilsIsRequiredCriticalHeader(jose.HeaderBase64URLEncodePayload) {
		t.Error("'b64' is not reported as a required critical header")
	}
	if DSSJsonUtilsIsRequiredCriticalHeader(JAdESHeaderParameterNamesSigT) {
		t.Error("'sigT' is reported as a required critical header")
	}
	supported := DSSJsonUtilsSupportedProtectedCriticalHeaders()
	for _, name := range []string{JAdESHeaderParameterNamesSigT, JAdESHeaderParameterNamesSigD,
		JWTClaimNamesIat, JWTClaimNamesExp, jose.HeaderBase64URLEncodePayload} {
		if _, found := supported[name]; !found {
			t.Errorf("%q is missing from the supported protected critical headers", name)
		}
	}
	// The accessor must hand out a copy, so that a caller cannot widen what this package
	// considers supported.
	supported["invented"] = struct{}{}
	if _, found := DSSJsonUtilsSupportedProtectedCriticalHeaders()["invented"]; found {
		t.Error("the returned set aliases the package-level one")
	}
}

// TestJWSConstantsAndClaimNames pins the string constants; PORTING.md forbids inventing,
// abbreviating or "fixing" any of them.
func TestJWSConstantsAndClaimNames(t *testing.T) {
	constants := map[string]string{
		JWSConstantsPayload:    "payload",
		JWSConstantsSignatures: "signatures",
		JWSConstantsProtected:  "protected",
		JWSConstantsHeader:     "header",
		JWSConstantsSignature:  "signature",
		JWTClaimNamesIss:       "iss",
		JWTClaimNamesSub:       "sub",
		JWTClaimNamesAud:       "aud",
		JWTClaimNamesExp:       "exp",
		JWTClaimNamesNbf:       "nbf",
		JWTClaimNamesIat:       "iat",
		JWTClaimNamesJti:       "jti",
	}
	for got, want := range constants {
		if got != want {
			t.Errorf("constant %q should be %q", got, want)
		}
	}
}

// TestFlattenedSerializationRejectsMultipleSignatures covers the one precondition of the
// flattened form.
func TestFlattenedSerializationRejectsMultipleSignatures(t *testing.T) {
	object := NewJWSJsonSerializationObject()
	object.AddSignature(NewJWS())
	object.AddSignature(NewJWS())
	_, err := NewJWSJsonSerializationGenerator(object,
		enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION).Generate()
	if err == nil {
		t.Fatal("a flattened serialization with two signatures was accepted")
	}

	empty := NewJWSJsonSerializationObject()
	if _, err := NewJWSJsonSerializationGenerator(empty,
		enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION).Generate(); err == nil {
		t.Fatal("a flattened serialization with no signature was accepted")
	}
	if _, err := NewJWSJsonSerializationGenerator(empty,
		enumerations.JWSSerializationType_COMPACT_SERIALIZATION).Generate(); err == nil {
		t.Fatal("the generator accepted the compact serialization type")
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func jadesJSONKATParseCompact(t *testing.T, compact []byte) *JWS {
	t.Helper()
	jws, err := NewJWSCompactSerializationParser(model.NewInMemoryDocument(compact)).Parse()
	if err != nil {
		t.Fatalf("parsing the compact JWS: %v", err)
	}
	return jws
}

func jadesJSONKATAssertBytes(t *testing.T, document model.DSSDocument, want []byte) {
	t.Helper()
	stream, err := document.OpenStream()
	if err != nil {
		t.Fatalf("OpenStream: %v", err)
	}
	defer stream.Close()
	got := make([]byte, 0, len(want)+16)
	buf := make([]byte, 4096)
	for {
		n, err := stream.Read(buf)
		got = append(got, buf[:n]...)
		if err != nil {
			break
		}
	}
	if string(got) != string(want) {
		t.Errorf("document bytes\n got %s\nwant %s", got, want)
	}
}

func jadesJSONKATUnhex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("expected hex, got %q: %v", s, err)
	}
	return b
}

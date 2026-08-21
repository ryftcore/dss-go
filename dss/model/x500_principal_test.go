package model

import (
	"encoding/hex"
	"strings"
	"testing"
)

// x500PrincipalGoldenCases are golden vectors captured from a real JDK 21 JVM:
// each DER encoding was handed to javax.security.auth.x500.X500Principal and the three name
// forms dss-model consumes were recorded verbatim. They pin the JDK behaviours the port has
// to reproduce - the RFC 2253 escaping rules, the CANONICAL pipeline (multi-valued RDN
// sorting, whitespace collapsing, US-locale upper-then-lower casing and NFKD normalization)
// and the hex fallback for the string types canonical form rejects.
var x500PrincipalGoldenCases = []struct {
	name      string
	der       string
	rfc2253   string
	canonical string
	pretty    string
}{
	{
		name:      "simple",
		der:       "3032310b3009060355040613024245310f300d060355040a13064e6f77696e613112301006035504031309546573742055736572",
		rfc2253:   "CN=Test User,O=Nowina,C=BE",
		canonical: "cn=test user,o=nowina,c=be",
		pretty:    "commonName=Test User,organizationName=Nowina,countryName=BE",
	},
	{
		name:      "utf8_accents",
		der:       "3025310b30090603550406130246523116301406035504030c0d4a6f73c3a920416d706c69c3a9",
		rfc2253:   "CN=Jos\u00e9 Ampli\u00e9,C=FR",
		canonical: "cn=jose\u0301 amplie\u0301,c=fr",
		pretty:    "commonName=Jos\u00e9 Ampli\u00e9,countryName=FR",
	},
	{
		name:      "utf8_sharp_s",
		der:       "30183116301406035504030c0d53747261c39f652047726fc39f",
		rfc2253:   "CN=Stra\u00dfe Gro\u00df",
		canonical: "cn=strasse gross",
		pretty:    "commonName=Stra\u00dfe Gro\u00df",
	},
	{
		name:      "utf8_greek",
		der:       "30193117301506035504030c0ece91ceb3ceb3ceb5cebbcebfcf82",
		rfc2253:   "CN=\u0391\u03b3\u03b3\u03b5\u03bb\u03bf\u03c2",
		canonical: "cn=\u03b1\u03b3\u03b3\u03b5\u03bb\u03bf\u03c2",
		pretty:    "commonName=\u0391\u03b3\u03b3\u03b5\u03bb\u03bf\u03c2",
	},
	{
		name:      "utf8_cyrillic",
		der:       "301531133011060355040a0c0ad090d0bad182d0b8d0b2",
		rfc2253:   "O=\u0410\u043a\u0442\u0438\u0432",
		canonical: "o=\u0430\u043a\u0442\u0438\u0432",
		pretty:    "organizationName=\u0410\u043a\u0442\u0438\u0432",
	},
	{
		name:      "spaces_inner",
		der:       "3015311330110603550403130a41202020422020202043",
		rfc2253:   "CN=A   B    C",
		canonical: "cn=a b c",
		pretty:    "commonName=A   B    C",
	},
	{
		name:      "spaces_edges",
		der:       "3017311530130603550403130c202020706164646564202020",
		rfc2253:   "CN=\\ \\ \\ padded\\ \\ \\ ",
		canonical: "cn=padded",
		pretty:    "commonName=\\ \\ \\ padded\\ \\ \\ ",
	},
	{
		name:      "specials",
		der:       "301e311c301a06035504030c13612c622b6322645c653c663e673b6823693d6a",
		rfc2253:   "CN=a\\,b\\+c\\\"d\\\\e\\<f\\>g\\;h\\#i\\=j",
		canonical: "cn=a\\,b\\+c\\\"d\\\\e\\<f\\>g\\;h#i=j",
		pretty:    "commonName=a\\,b\\+c\\\"d\\\\e\\<f\\>g\\;h\\#i\\=j",
	},
	{
		name:      "leading_hash",
		der:       "3010310e300c060355040313052368617368",
		rfc2253:   "CN=\\#hash",
		canonical: "cn=\\#hash",
		pretty:    "commonName=\\#hash",
	},
	{
		name:      "bmpstring",
		der:       "301b3119301706035504031e100042004d00500020004e0061006d0065",
		rfc2253:   "CN=BMP Name",
		canonical: "cn=#1e100042004d00500020004e0061006d0065",
		pretty:    "commonName=BMP Name",
	},
	{
		name:      "t61string",
		der:       "30133111300f06035504031408543631204e616d65",
		rfc2253:   "CN=T61 Name",
		canonical: "cn=#1408543631204e616d65",
		pretty:    "commonName=T61 Name",
	},
	{
		name:      "ia5_email",
		der:       "301c311a301806092a864886f70d010901160b6140622e6578616d706c65",
		rfc2253:   "1.2.840.113549.1.9.1=#160b6140622e6578616d706c65",
		canonical: "1.2.840.113549.1.9.1=#160b6140622e6578616d706c65",
		pretty:    "emailAddress=a@b.example",
	},
	{
		name:      "generalstring",
		der:       "30123110300e06035504031b0747656e6572616c",
		rfc2253:   "CN=General",
		canonical: "cn=#1b0747656e6572616c",
		pretty:    "commonName=General",
	},
	{
		name:      "no_keyword_oid",
		der:       "30133111300f060355040c13084469726563746f72",
		rfc2253:   "2.5.4.12=#13084469726563746f72",
		canonical: "2.5.4.12=#13084469726563746f72",
		pretty:    "title=Director",
	},
	{
		name:      "custom_oid",
		der:       "30173115301306092b06010401868d1f011306437573746f6d",
		rfc2253:   "1.3.6.1.4.1.99999.1=#1306437573746f6d",
		canonical: "1.3.6.1.4.1.99999.1=#1306437573746f6d",
		pretty:    "1.3.6.1.4.1.99999.1=#1306437573746f6d",
	},
	{
		name:      "multivalued",
		der:       "30283126300b060355040313045a657461300c060355040a1305416c7068613009060355040c13024d75",
		rfc2253:   "CN=Zeta+O=Alpha+2.5.4.12=#13024d75",
		canonical: "cn=zeta+o=alpha+2.5.4.12=#13024d75",
		pretty:    "commonName=Zeta+organizationName=Alpha+title=Mu",
	},
	{
		name:      "multivalued2",
		der:       "301a3118300a060355040c13037a7a7a300a06035504031303616161",
		rfc2253:   "2.5.4.12=#13037a7a7a+CN=aaa",
		canonical: "cn=aaa+2.5.4.12=#13037a7a7a",
		pretty:    "title=zzz+commonName=aaa",
	},
	{
		name:      "empty_dn",
		der:       "3000",
		rfc2253:   "",
		canonical: "",
		pretty:    "",
	},
	{
		name:      "empty_value",
		der:       "300b3109300706035504031300",
		rfc2253:   "CN=",
		canonical: "cn=",
		pretty:    "commonName=",
	},
	{
		name:      "latin1_t61",
		der:       "300f310d300b060355040314044a6f73e9",
		rfc2253:   "CN=Jos\u00e9",
		canonical: "cn=#14044a6f73e9",
		pretty:    "commonName=Jos\u00e9",
	},
	{
		name:      "dc_uid",
		der:       "302f31173015060a0992268993f22c64011916076578616d706c6531143012060a0992268993f22c64010113046a646f65",
		rfc2253:   "UID=jdoe,DC=example",
		canonical: "uid=jdoe,dc=#16076578616d706c65",
		pretty:    "UID=jdoe,DC=example",
	},
	{
		name:      "street_st",
		der:       "30363110300e0603550408130742726162616e74310e300c060355040913055275652031311230100603550407130942727578656c6c6573",
		rfc2253:   "L=Bruxelles,STREET=Rue 1,ST=Brabant",
		canonical: "l=bruxelles,street=rue 1,st=brabant",
		pretty:    "L=Bruxelles,STREET=Rue 1,ST=Brabant",
	},
	{
		name:      "serialnumber",
		der:       "30133111300f060355040513083132333435363738",
		rfc2253:   "2.5.4.5=#13083132333435363738",
		canonical: "2.5.4.5=#13083132333435363738",
		pretty:    "2.5.4.5=#13083132333435363738",
	},
	{
		name:      "tab_and_nl",
		der:       "30123110300e06035504030c076109620a630d64",
		rfc2253:   "CN=a\tb\nc\rd",
		canonical: "cn=a\tb\nc\rd",
		pretty:    "commonName=a\tb\nc\rd",
	},
	{
		name:      "ligature_nfkd",
		der:       "30143112301006035504030c096fefac81636520c2bd",
		rfc2253:   "CN=o\ufb01ce \u00bd",
		canonical: "cn=ofice 1\u20442",
		pretty:    "commonName=o\ufb01ce \u00bd",
	},
	{
		name:      "turkish_dotted_i",
		der:       "30143112301006035504030c09c4b07374616e62756c",
		rfc2253:   "CN=\u0130stanbul",
		canonical: "cn=i\u0307stanbul",
		pretty:    "commonName=\u0130stanbul",
	},
	{
		name:      "nbsp",
		der:       "300f310d300b06035504030c0461c2a062",
		rfc2253:   "CN=a\u00a0b",
		canonical: "cn=a b",
		pretty:    "commonName=a\u00a0b",
	},
	{
		name:      "supplementary",
		der:       "3011310f300d06035504030c0678f09d908079",
		rfc2253:   "CN=x\U0001d400y",
		canonical: "cn=xAy",
		pretty:    "commonName=x\U0001d400y",
	},
	{
		name:      "trailing_cr",
		der:       "300e310c300a06035504031303616263",
		rfc2253:   "CN=abc",
		canonical: "cn=abc",
		pretty:    "commonName=abc",
	},
}

// x500PrincipalPrettyOIDMap is the OID/keyword map the golden vectors were captured with.
// It is a subset of enumerations.GetOidDescriptions(), kept explicit so that the expectations
// stay pinned to fixed keywords.
var x500PrincipalPrettyOIDMap = map[string]string{
	"2.5.4.3":              "commonName",
	"2.5.4.6":              "countryName",
	"2.5.4.10":             "organizationName",
	"2.5.4.11":             "organizationUnitName",
	"2.5.4.12":             "title",
	"1.2.840.113549.1.9.1": "emailAddress",
}

func TestX500PrincipalNameFormsMatchTheJDK(t *testing.T) {
	for _, tc := range x500PrincipalGoldenCases {
		t.Run(tc.name, func(t *testing.T) {
			der, err := hex.DecodeString(tc.der)
			if err != nil {
				t.Fatalf("bad fixture: %v", err)
			}
			principal, err := NewX500Principal(der)
			if err != nil {
				t.Fatalf("NewX500Principal: %v", err)
			}
			if got := principal.RFC2253Name(); got != tc.rfc2253 {
				t.Errorf("RFC2253Name()\n got %q\nwant %q", got, tc.rfc2253)
			}
			if got := principal.Canonical(); got != tc.canonical {
				t.Errorf("Canonical()\n got %q\nwant %q", got, tc.canonical)
			}
			pretty, err := principal.RFC2253NameWithOIDMap(x500PrincipalPrettyOIDMap)
			if err != nil {
				t.Fatalf("RFC2253NameWithOIDMap: %v", err)
			}
			if pretty != tc.pretty {
				t.Errorf("RFC2253NameWithOIDMap()\n got %q\nwant %q", pretty, tc.pretty)
			}
		})
	}
}

func TestX500PrincipalCachesAndKeepsTheEncoding(t *testing.T) {
	der, _ := hex.DecodeString(x500PrincipalGoldenCases[0].der)
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatal(err)
	}
	if got := principal.Encoded(); hex.EncodeToString(got) != x500PrincipalGoldenCases[0].der {
		t.Errorf("Encoded() must return the DER unchanged, got %x", got)
	}
	// The second call must come from the cache and produce the same string.
	firstRFC2253Name, secondRFC2253Name := principal.RFC2253Name(), principal.RFC2253Name()
	if firstRFC2253Name != secondRFC2253Name {
		t.Error("RFC2253Name() is not stable")
	}
	firstCanonical, secondCanonical := principal.Canonical(), principal.Canonical()
	if firstCanonical != secondCanonical {
		t.Error("Canonical() is not stable")
	}
	if principal.String() != principal.RFC2253Name() {
		t.Error("String() must return the RFC 2253 form")
	}
}

func TestX500PrincipalEqualsComparesCanonicalForms(t *testing.T) {
	// Same name, but the CN is encoded once as PrintableString and once as UTF8String, and
	// the case differs: X500Principal#equals compares canonical forms, so both are equal.
	printable, err := NewX500Principal(mustHex(t, "3014311230100603550403130954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	utf8Upper, err := NewX500Principal(mustHex(t, "30143112301006035504030c0954455354204e414d45"))
	if err != nil {
		t.Fatal(err)
	}
	other, err := NewX500Principal(mustHex(t, "30153113301106035504030c0a4f74686572204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	if !printable.Equals(utf8Upper) {
		t.Errorf("expected %q to equal %q", printable.Canonical(), utf8Upper.Canonical())
	}
	if printable.Equals(other) {
		t.Error("different names must not be equal")
	}
	if printable.Equals(nil) {
		t.Error("a principal must not equal nil")
	}
}

func TestX500PrincipalRejectsMalformedEncodings(t *testing.T) {
	for name, encoded := range map[string]string{
		"not a sequence":     "310a30080603550403130141",
		"rdn not a set":      "300c300a3008060355040313014121",
		"truncated":          "300a3008060355040313",
		"trailing bytes":     "30003000",
		"ava not a sequence": "300c310a31080603550403130141",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := NewX500Principal(mustHex(t, encoded)); err == nil {
				t.Errorf("expected an error for %s", name)
			}
		})
	}
}

func TestX500PrincipalKeywordValidation(t *testing.T) {
	der := mustHex(t, "3014311230100603550403130954657374204e616d65")
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatal(err)
	}
	for name, oidMap := range map[string]map[string]string{
		"empty keyword":     {"2.5.4.3": ""},
		"leading digit":     {"2.5.4.3": "1bad"},
		"illegal character": {"2.5.4.3": "bad-keyword"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := principal.RFC2253NameWithOIDMap(oidMap); err == nil {
				t.Errorf("expected an error for %s", name)
			}
		})
	}
	// A keyword is trimmed the way java.lang.String#trim does, and wins over the built-in one.
	got, err := principal.RFC2253NameWithOIDMap(map[string]string{"2.5.4.3": "  myCN  "})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "myCN=") {
		t.Errorf("expected the mapped keyword to win, got %q", got)
	}
}

func TestX500PrincipalJavaHelpers(t *testing.T) {
	// Character.isWhitespace excludes the non-breaking spaces unicode.IsSpace accepts.
	for _, r := range []rune{' ', '\t', '\n', '\v', '\f', '\r', 0x1C, 0x1F, 0x2028, 0x2029, 0x2000} {
		if !x500PrincipalIsJavaWhitespace(r) {
			t.Errorf("U+%04X should be Java whitespace", r)
		}
	}
	for _, r := range []rune{0x00A0, 0x2007, 0x202F, 0x0085, 'a'} {
		if x500PrincipalIsJavaWhitespace(r) {
			t.Errorf("U+%04X should not be Java whitespace", r)
		}
	}
	// String#trim strips everything up to and including U+0020, and nothing above it.
	if got := x500PrincipalJavaTrim("\x01\x02 abc \x00"); got != "abc" {
		t.Errorf("javaTrim = %q", got)
	}
	if got := x500PrincipalJavaTrim(" abc "); got != " abc " {
		t.Errorf("javaTrim must not strip U+00A0, got %q", got)
	}
	// String#compareTo orders by UTF-16 code unit, so a supplementary character sorts below
	// U+E000..U+FFFF even though its UTF-8 bytes are larger.
	if x500PrincipalCompareJavaStrings("\U0001D400", "") >= 0 {
		t.Error("expected the supplementary character to sort first, as Java does")
	}
	if x500PrincipalCompareJavaStrings("abc", "abc") != 0 {
		t.Error("equal strings must compare equal")
	}
	if x500PrincipalCompareJavaStrings("ab", "abc") >= 0 {
		t.Error("a prefix must sort first")
	}
}

func TestX500PrincipalDecodeOID(t *testing.T) {
	for encoded, want := range map[string]string{
		"550403":               "2.5.4.3",
		"0992268993f22c640101": "0.9.2342.19200300.100.1.1",
		"2a864886f70d010901":   "1.2.840.113549.1.9.1",
	} {
		got, err := x500PrincipalDecodeOID(mustHex(t, encoded))
		if err != nil {
			t.Fatalf("%s: %v", encoded, err)
		}
		if got != want {
			t.Errorf("decodeOID(%s) = %s, want %s", encoded, got, want)
		}
	}
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("bad fixture %q: %v", s, err)
	}
	return b
}

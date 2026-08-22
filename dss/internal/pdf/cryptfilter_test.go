// The /V 4 and /V 5 crypt-filter acceptance boundary (issue #33), and the /V 0
// rejection that rides along with it.
//
// Before checkCryptFilters existed, useAES was set from /CFM alone inside
// `case 4, 5:` of setupEncryption, with no else and no default: a /V 5 document
// whose crypt filter could not be identified - no /CF, a /CF without the named
// entry, or a /CFM this handler does not implement - left useAES false, and the
// 32 bytes unwrapped from /UE were handed to rc4Apply on the read path
// (decryptBytes) and the write path (encryptForWrite) alike, while Encryption()
// reported /CFM /None. These tests pin the refusal, pin the three shapes that
// must keep opening so the refusal stays exactly as narrow as "a crypt filter
// the document actually selects, that this handler cannot identify", and pin
// that /V 0 - "shall not be used" per ISO 32000-1 Table 20 - is refused too.
//
// The AES-256 (/V 5 /R 6) builder is cryptaes256_test.go's buildAES256Doc;
// buildAESV2Doc below is this file's own /V 4 /R 4 (AES-128) equivalent, built
// the same way - /O, /U and the file key derived with the same crypt.go
// functions the reader checks them with - since no /V 4 builder existed before
// this file. buildRC4Document in crypt_test.go is the /V 1 /R 2 control for the
// /V 0 case.
package pdf

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// stdCFAESV2 is the crypt-filter section a well-formed /V 4 AES-128 document
// carries. Mirrors stdCFAESV3 in cryptaes256_test.go.
const stdCFAESV2 = "/CF << /StdCF << /CFM /AESV2 /Length 16 >> >> /StmF /StdCF /StrF /StdCF"

// buildAESV2Doc produces a real /Filter /Standard /V 4 /R 4 /Length 128
// document that opens with the empty user password and the owner password
// "owner". cfEntry is spliced into /Encrypt verbatim, exactly as
// buildAES256Doc's cfEntry parameter is, so a fixture can name any /CF shape
// (including none at all) without touching key derivation.
//
//	/O = Algorithm 3: the padded user password, RC4-encrypted 20 times with the
//	     owner password's RC4 key, each round XORed by the round number
//	/U  = Algorithm 4/5's 32 bytes (computeUserEntry)
func buildAESV2Doc(t *testing.T, cfEntry string) []byte {
	t.Helper()
	const rev = 4
	const keyLen = 16
	const perm = int32(-1052)
	id := []byte("0123456789abcdef")
	ownerPw := []byte("owner")
	var userPw []byte // the empty user password

	rc4Key := computeRC4Key(ownerPw, rev, keyLen)
	owner := truncateOrPad(userPw)
	iter := make([]byte, len(rc4Key))
	for i := 0; i <= 19; i++ {
		copy(iter, rc4Key)
		for j := range iter {
			iter[j] ^= byte(i)
		}
		owner = rc4Apply(iter, owner)
	}

	fileKey := computeKeyRev234(userPw, owner, perm, id, true, keyLen, rev)
	user := computeUserEntry(userPw, owner, perm, id, rev, keyLen, true)

	enc := aesEncryptorFor(t, fileKey)
	secret := enc(4, []byte("classified"), true)
	streamData := enc(6, []byte("stream payload"), false)

	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"},
		{num: 4, body: "<< /Secret <" + hex.EncodeToString(secret) + "> >>"},
		{num: 5, body: fmt.Sprintf("<< /Filter /Standard /V 4 /R 4 /Length 128 %s /P %d "+
			"/O <%s> /U <%s> >>",
			cfEntry, perm, hex.EncodeToString(owner), hex.EncodeToString(user))},
		{num: 6, body: fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream",
			len(streamData), streamData)},
	}
	return buildReaderPDF("%PDF-1.6\n", objs,
		fmt.Sprintf("/Encrypt 5 0 R\n/ID [<%s> <%s>]\n",
			hex.EncodeToString(id), hex.EncodeToString(id)))
}

// asV4R6 rewrites a document buildAES256Doc produced into the /V 4 /R 6
// crossbreed: the crypt-filter rules of /V 4 (where /CFM /V2 is legitimate) over
// the key derivation of /R 6 (where the file key is the 32 bytes unwrapped from
// /UE or /OE, independent of /V - computeEncryptionKey branches on /R alone).
// "/V 5 /R 6" and "/V 4 /R 6" are the same width, so no xref offset moves.
func asV4R6(t *testing.T, data []byte) []byte {
	t.Helper()
	out := replaceOnce(t, data, []byte("/V 5 /R 6"), []byte("/V 4 /R 6"))
	if out == nil {
		t.Fatal("could not locate a single /V 5 /R 6 in the fixture")
	}
	return out
}

// TestEncryptionRejectsUnidentifiableCryptFilter is the read-side regression for
// issue #33: every crypt filter a /V 4 or /V 5 document selects must resolve
// through /CF to a /CFM this handler implements, and for /V 5 that /CFM must be
// /AESV3.
//
// Before checkCryptFilters existed, every row here opened with useAES false and
// handed the declared cipher's key to rc4Apply; `if err == nil { t.Fatalf(...) }`
// below is the assertion that fires in that world.
func TestEncryptionRejectsUnidentifiableCryptFilter(t *testing.T) {
	for _, tc := range []struct {
		name string
		why  string
		data func(t *testing.T) []byte
	}{
		{"V5 with no /CF at all", "/StmF names /StdCF, nothing defines it", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{cfEntry: "/StmF /StdCF /StrF /StdCF", pValue: "-1052"})
		}},
		{"V5 whose /CF lacks the named filter", "/CF defines /Other, /StmF names /StdCF", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /Other << /CFM /AESV3 /Length 32 >> >> /StmF /StdCF /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V5 whose named filter carries no /CFM", "the entry resolves, the method does not", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /AuthEvent /DocOpen /Length 32 >> >> /StmF /StdCF /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V5 declaring the RC4 /CFM /V2", "AES-256 key material, RC4 requested", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /V2 /Length 32 >> >> /StmF /StdCF /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V5 declaring /CFM /AESV2", "/V 5 is /AESV3 only; AESV2 with a 32-byte key is AES-256 under an AES-128 name", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /AESV2 /Length 16 >> >> /StmF /StdCF /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V5 declaring /CFM /None", "/None is also crypt.go's own sentinel for unresolved; both refuse", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /None /Length 32 >> >> /StmF /StdCF /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V5 with /StmF /Identity and an unidentifiable /StrF", "the string filter alone is enough to select a cipher", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /V2 /Length 32 >> >> /StmF /Identity /StrF /StdCF",
				pValue:  "-1052"})
		}},
		{"V4 with no /CF at all", "/CF is required when /V is 4", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/StmF /StdCF /StrF /StdCF")
		}},
		{"V4 whose /CF lacks the named filter", "", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /Other << /CFM /AESV2 /Length 16 >> >> /StmF /StdCF /StrF /StdCF")
		}},
		{"V4 declaring an unknown /CFM", "/Foo is not a method this handler implements", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /StdCF << /CFM /Foo /Length 16 >> >> /StmF /StdCF /StrF /StdCF")
		}},
		{"V4 declaring /CFM /None", "ISO 32000-1 Table 25 hands /None back to the security handler; we have no such handler", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /StdCF << /CFM /None /Length 16 >> >> /StmF /StdCF /StrF /StdCF")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.why != "" {
				t.Logf("why: %s", tc.why)
			}
			d, err := OpenBytes(tc.data(t), nil)
			if err == nil {
				t.Fatalf("opened a document whose crypt filter cannot be identified: "+
					"Encryption() = %+v, useAES = %v, file key = %d bytes - the cipher "+
					"applied is not the cipher declared", d.Encryption(), d.sec.useAES, len(d.sec.key))
			}
			if !errors.Is(err, ErrUnsupportedSecurityHandler) {
				t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
			}
		})
	}
}

// TestEncryptionRejectsDisagreeingCryptFilters pins the half of issue #33 a
// single-name check misses: setupEncryption derives the handler's /CFM from
// /StmF, falling back to /StrF only when /StmF is /Identity, so a document
// naming two different crypt filters has exactly one of them read. The handler
// carries one key and one useAES flag and cannot honour two methods, so the
// other filter's declaration is silently overruled - in the first row below,
// every string in the document is RC4-decrypted although /StrF /StdCF declares
// AES-128. checkCryptFilters resolves /StmF and /StrF independently and refuses
// when they disagree, which is exactly what a per-cfName-only check cannot do:
// it never looks at the filter setupEncryption did not pick.
func TestEncryptionRejectsDisagreeingCryptFilters(t *testing.T) {
	for _, tc := range []struct {
		name string
		data func(t *testing.T) []byte
	}{
		{"V4 streams /V2, strings /AESV2 (the exact single-name-check trap)", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /StdCF << /CFM /AESV2 /Length 16 >> "+
				"/RC4CF << /CFM /V2 /Length 16 >> >> /StmF /RC4CF /StrF /StdCF")
		}},
		{"V4 streams /AESV2, strings /V2", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /StdCF << /CFM /AESV2 /Length 16 >> "+
				"/RC4CF << /CFM /V2 /Length 16 >> >> /StmF /StdCF /StrF /RC4CF")
		}},
		{"V5 streams /AESV3, strings /V2", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /AESV3 /Length 32 >> " +
					"/RC4CF << /CFM /V2 /Length 32 >> >> /StmF /StdCF /StrF /RC4CF",
				pValue: "-1052"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := OpenBytes(tc.data(t), nil)
			if err == nil {
				t.Fatalf("opened a document whose /StmF and /StrF select different "+
					"crypt filter methods: Encryption() = %+v, useAES = %v - one of the "+
					"two declarations is being applied to both", d.Encryption(), d.sec.useAES)
			}
			if !errors.Is(err, ErrUnsupportedSecurityHandler) {
				t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
			}
		})
	}
}

// TestEncryptionRejectsV5PairedWithWrongRevision is the direct regression for
// F2 (issue #33's last surviving instance): ISO 32000-2 pairs /V 5 with /R 5
// or /R 6 only, but neither computeEncryptionKey nor the postcondition guard
// in setupEncryption looks at /V - both key off /R alone. Before
// checkCryptFilters rejected the pairing, a /V 4 /R 4 document edited to
// /V 5 /R 4, with a conforming /CF resolving to /AESV3, opened as
// {V:5 R:4 KeyLength:128 CFM:AESV3}: checkCryptFilters' per-filter loop saw
// exactly the /CFM /V 5 requires and had nothing to object to, useAES came
// back true (so the RC4-only postcondition guard never fired either), and
// computeEncryptionKey took /R 4's computeKeyRev234 (MD5) path and produced a
// <=16-byte key - AES-128, MD5-derived, salted, in a document declaring
// AES-256.
func TestEncryptionRejectsV5PairedWithWrongRevision(t *testing.T) {
	build := func(t *testing.T) []byte {
		data := buildAESV2Doc(t, "/CF << /StdCF << /CFM /AESV3 /Length 16 >> >> /StmF /StdCF /StrF /StdCF")
		out := replaceOnce(t, data, []byte("/V 4 /R 4"), []byte("/V 5 /R 4"))
		if out == nil {
			t.Fatal("could not locate a single /V 4 /R 4 in the fixture")
		}
		return out
	}

	t.Run("V5 /R4 with a conforming /CF is refused, not downgraded to AES-128", func(t *testing.T) {
		d, err := OpenBytes(build(t), nil)
		if err == nil {
			t.Fatalf("opened a /V 5 document declaring /R 4: Encryption() = %+v, "+
				"useAES = %v, file key = %d bytes - AES-256 was declared, a shorter "+
				"MD5-derived key was applied", d.Encryption(), d.sec.useAES, len(d.sec.key))
		}
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
		if got := err.Error(); !strings.Contains(got, "/V 5") || !strings.Contains(got, "/R 4") {
			t.Errorf("error %q does not name the offending /V 5 /R 4 pairing", got)
		}
	})

	// The converse of the rule above, and the same defect wearing /V 4: an
	// /AESV3 crypt filter under /R 4 promises a 32-byte key that only the
	// /UE //OE unwrap of /R 5 //R 6 ever produces, so before this rule the shape
	// opened as {V:4 R:4 KeyLength:256 CFM:AESV3} and applied AES-128 under an
	// MD5 key on both the read and the write path. Constraining /V 5 alone
	// would have left the downgrade intact one /V value away.
	t.Run("V4 /R4 naming /CFM /AESV3 is refused for the same reason", func(t *testing.T) {
		d, err := OpenBytes(buildAESV2Doc(t,
			"/CF << /StdCF << /CFM /AESV3 /Length 16 >> >> /StmF /StdCF /StrF /StdCF"), nil)
		if err == nil {
			t.Fatalf("opened a /R 4 document declaring /CFM /AESV3: Encryption() = %+v, "+
				"useAES = %v, file key = %d bytes - AES-256 was declared, AES-128 applied",
				d.Encryption(), d.sec.useAES, len(d.sec.key))
		}
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
		if got := err.Error(); !strings.Contains(got, "AESV3") || !strings.Contains(got, "/R 4") {
			t.Errorf("error %q does not name the offending /AESV3 under /R 4", got)
		}
	})

	t.Run("V4 /R4 with /CFM /AESV2 still opens (positive control)", func(t *testing.T) {
		d, err := OpenBytes(buildAESV2Doc(t,
			"/CF << /StdCF << /CFM /AESV2 /Length 16 >> >> /StmF /StdCF /StrF /StdCF"), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.V != 4 || enc.R != 4 || enc.CFM != "AESV2" {
			t.Errorf("encryption = %+v, want /V 4 /R 4 /CFM /AESV2", enc)
		}
	})
}

// TestEncryptionCryptFilterPositiveControls bounds the refusal. Each shape here
// selects a cipher this handler implements, or selects none at all, and must
// keep opening and keep decrypting or being left alone.
func TestEncryptionCryptFilterPositiveControls(t *testing.T) {
	for _, tc := range []struct {
		name   string
		data   func(t *testing.T) []byte
		cfm    Name
		useAES bool
	}{
		{"V5 with a conforming /CF", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052"})
		}, "AESV3", true},
		{"V4 with a conforming /CF", func(t *testing.T) []byte {
			return buildAESV2Doc(t, stdCFAESV2)
		}, "AESV2", true},
		{"V5 with /StmF /Identity, strings still AES-256", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{
				cfEntry: "/CF << /StdCF << /CFM /AESV3 /Length 32 >> >> /StmF /Identity /StrF /StdCF",
				pValue:  "-1052"})
		}, "AESV3", true},
		{"V4 with /StrF /Identity, streams still AES-128", func(t *testing.T) []byte {
			return buildAESV2Doc(t, "/CF << /StdCF << /CFM /AESV2 /Length 16 >> >> /StmF /StdCF /StrF /Identity")
		}, "AESV2", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, err := OpenBytes(tc.data(t), nil)
			if err != nil {
				t.Fatalf("OpenBytes: %v", err)
			}
			if enc := d.Encryption(); enc.CFM != tc.cfm {
				t.Errorf("CFM %q, want %q", enc.CFM, tc.cfm)
			}
			if d.sec.useAES != tc.useAES {
				t.Errorf("useAES = %v, want %v", d.sec.useAES, tc.useAES)
			}
			// The /Secret string is only reachable when /StrF is not /Identity;
			// where it is, the ciphertext must come back untouched instead.
			obj, err := d.Object(ObjectKey{Num: 4})
			if err != nil {
				t.Fatal(err)
			}
			got, _ := d.GetString(obj.(*Dict), "Secret")
			if d.Encryption().StrF == "Identity" {
				if string(got) == "classified" {
					t.Error("a string was decrypted although /StrF is /Identity")
				}
				return
			}
			if string(got) != "classified" {
				t.Errorf("decrypted string = %q, want %q", got, "classified")
			}
		})
	}

	t.Run("V5 naming no crypt filter is /Identity, not a fallback", func(t *testing.T) {
		// ISO 32000-2 Table 20: /StmF and /StrF both default to /Identity. No
		// cipher is selected, so there is nothing to identify and nothing to
		// refuse - and no byte is ever decrypted, so the 32-byte key never
		// reaches rc4Apply.
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{pValue: "-1052"}), nil)
		if err != nil {
			t.Fatalf("a /V 5 document with no /StmF, /StrF or /CF must open as /Identity: %v", err)
		}
		enc := d.Encryption()
		if enc.StmF != "Identity" || enc.StrF != "Identity" {
			t.Errorf("encryption = %+v, want both filters /Identity", enc)
		}
		if d.sec.useAES {
			t.Error("useAES with no crypt filter named")
		}
	})

	t.Run("V4 with both filters /Identity and no /CF", func(t *testing.T) {
		// The legitimately-unencrypted /V 4 shape: an /Encrypt dictionary that
		// selects nothing. An unidentifiable /CF is harmless here because no
		// cipher is ever applied.
		d, err := OpenBytes(buildAESV2Doc(t, "/StmF /Identity /StrF /Identity"), nil)
		if err != nil {
			t.Fatalf("a /V 4 document selecting /Identity for both must open: %v", err)
		}
		if d.sec.useAES {
			t.Error("useAES with no crypt filter named")
		}
	})
}

// buildAESV2DocIndirectFilters is buildAESV2Doc with /StmF and /StrF spliced in
// as whatever raw text the caller supplies, plus one extra object (8, holding
// the crypt filter name /StdCF) that a caller can point an indirect /StmF or
// /StrF at. Kept separate from buildAESV2Doc because that builder has no way to
// add an object of its own.
func buildAESV2DocIndirectFilters(t *testing.T, stmf, strf string) []byte {
	t.Helper()
	const rev = 4
	const keyLen = 16
	const perm = int32(-1052)
	id := []byte("0123456789abcdef")
	ownerPw := []byte("owner")
	var userPw []byte // the empty user password

	rc4Key := computeRC4Key(ownerPw, rev, keyLen)
	owner := truncateOrPad(userPw)
	iter := make([]byte, len(rc4Key))
	for i := 0; i <= 19; i++ {
		copy(iter, rc4Key)
		for j := range iter {
			iter[j] ^= byte(i)
		}
		owner = rc4Apply(iter, owner)
	}

	fileKey := computeKeyRev234(userPw, owner, perm, id, true, keyLen, rev)
	user := computeUserEntry(userPw, owner, perm, id, rev, keyLen, true)

	h := &securityHandler{key: fileKey, useAES: true, encryptMeta: true}
	h.enc.KeyLength = keyLen * 8
	h.enc.StmF, h.enc.StrF = "StdCF", "StdCF"
	d0 := &Document{opts: Options{Random: fixedRandom()}.withDefaults(), sec: h}
	enc := func(num int64, data []byte, isString bool) []byte {
		out, err := d0.encryptForWrite(ObjectKey{Num: num}, data, isString)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	secret := enc(4, []byte("classified"), true)
	streamData := enc(6, []byte("stream payload"), false)

	cfEntry := fmt.Sprintf("/CF << /StdCF << /CFM /AESV2 /Length 16 >> >> %s %s", stmf, strf)
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"},
		{num: 4, body: "<< /Secret <" + hex.EncodeToString(secret) + "> >>"},
		{num: 5, body: fmt.Sprintf("<< /Filter /Standard /V 4 /R 4 /Length 128 %s /P %d "+
			"/O <%s> /U <%s> >>",
			cfEntry, perm, hex.EncodeToString(owner), hex.EncodeToString(user))},
		{num: 6, body: fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream",
			len(streamData), streamData)},
		{num: 8, body: "/StdCF"},
	}
	return buildReaderPDF("%PDF-1.6\n", objs,
		fmt.Sprintf("/Encrypt 5 0 R\n/ID [<%s> <%s>]\n",
			hex.EncodeToString(id), hex.EncodeToString(id)))
}

// TestEncryptionResolvesIndirectCryptFilterNames is the direct regression for
// F3: crypt.go used to read /StmF and /StrF with a bare
// `encDict.GetRaw("StmF").(Name)`, which does not follow a reference the way
// PDEncryption.getStreamFilterName/getStringFilterName's
// `dictionary.getCOSName(...)` -> getDictionaryObject does. An indirect
// /StmF or /StrF therefore failed the type assertion, was silently read as
// absent, fell back to /Identity, and - because checkCryptFilters only looks
// at what setupEncryption resolved - was never even checked: streams or
// strings the document declares encrypted came back as raw ciphertext.
func TestEncryptionResolvesIndirectCryptFilterNames(t *testing.T) {
	t.Run("an indirect /StmF is honoured", func(t *testing.T) {
		d, err := OpenBytes(buildAESV2DocIndirectFilters(t, "/StmF 8 0 R", "/StrF /StdCF"), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.StmF != "StdCF" || enc.CFM != "AESV2" {
			t.Errorf("encryption = %+v, want /StmF /StdCF /CFM /AESV2", enc)
		}
		if !d.sec.useAES {
			t.Error("useAES = false, want true")
		}
		obj, err := d.Object(ObjectKey{Num: 6})
		if err != nil {
			t.Fatal(err)
		}
		raw, err := d.RawStreamData(obj.(*Stream))
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != "stream payload" {
			t.Errorf("decrypted stream = %q, want %q (an indirect /StmF must select AES-128, not /Identity)",
				raw, "stream payload")
		}
	})

	t.Run("an indirect /StrF is honoured", func(t *testing.T) {
		d, err := OpenBytes(buildAESV2DocIndirectFilters(t, "/StmF /StdCF", "/StrF 8 0 R"), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.StrF != "StdCF" || enc.CFM != "AESV2" {
			t.Errorf("encryption = %+v, want /StrF /StdCF /CFM /AESV2", enc)
		}
		obj, err := d.Object(ObjectKey{Num: 4})
		if err != nil {
			t.Fatal(err)
		}
		got, _ := d.GetString(obj.(*Dict), "Secret")
		if string(got) != "classified" {
			t.Errorf("decrypted string = %q, want %q (an indirect /StrF must select AES-128, not /Identity)",
				got, "classified")
		}
	})

	t.Run("an indirect /StmF still forces checkCryptFilters to run on it", func(t *testing.T) {
		// The other half of the bug: because an unresolved indirect /StmF read
		// as /Identity, checkCryptFilters never saw it and had nothing to
		// refuse. With resolution in place, an indirect /StmF naming a crypt
		// filter this handler cannot identify is caught exactly like a direct
		// one would be.
		data := buildAESV2DocIndirectFilters(t, "/StmF 8 0 R", "/StrF /StdCF")
		data = replaceOnce(t, data, []byte("8 0 obj\n/StdCF"), []byte("8 0 obj\n/Other"))
		if data == nil {
			t.Fatal("could not locate the /StdCF object body to rename")
		}
		_, err := OpenBytes(data, nil)
		if err == nil {
			t.Fatal("opened a document whose indirect /StmF names an unidentifiable crypt filter")
		}
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
	})

	t.Run("a non-Name /StmF still degrades to /Identity (unchanged)", func(t *testing.T) {
		// A resolved value that is not a Name - here, a literal string in place
		// of a name - is exactly as unidentifiable as an absent one, and the
		// fallback to /Identity for that case is deliberate (ISO 32000-2
		// Table 20's default), not a bug F3 touches.
		d, err := OpenBytes(buildAESV2DocIndirectFilters(t, "/StmF(StdCF)", "/StrF /StdCF"), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.StmF != "Identity" {
			t.Errorf("encryption = %+v, want /StmF /Identity", enc)
		}
		obj, err := d.Object(ObjectKey{Num: 6})
		if err != nil {
			t.Fatal(err)
		}
		raw, err := d.RawStreamData(obj.(*Stream))
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) == "stream payload" {
			t.Error("a stream was decrypted although /StmF resolves to a non-Name, want /Identity")
		}
	})
}

// TestEncryptionRejectsUndocumentedV0 covers the second half of issue #33.
//
// ISO 32000-1 Table 20 documents /V 0 as "an algorithm that is undocumented ...
// shall not be used", and 0 is also what an /Encrypt dictionary with no /V entry
// reads as (dictInt defaults to 0). crypt.go used to fold it into `case 0, 1:`
// and force a 5-byte (40-bit) key, which is not even what pdfbox guesses for it
// (pdfbox takes 5 bytes only for /V 1; /V 0 gets /Length/8 like every other
// unlisted /V) - so the two implementations silently disagreed on the key length
// of a document neither can name an algorithm for. /V 1 is unaffected: it keeps
// its own case, and its 5-byte key is now exact pdfbox parity instead of an
// incidental match.
func TestEncryptionRejectsUndocumentedV0(t *testing.T) {
	base := buildRC4Document(t, []byte("owner"), []byte(""), -44)

	t.Run("V1 still opens and still decrypts", func(t *testing.T) {
		d, err := OpenBytes(base, nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.V != 1 || enc.KeyLength != 40 {
			t.Errorf("encryption = %+v, want /V 1 with a 40-bit key", enc)
		}
		obj, err := d.Object(ObjectKey{Num: 4})
		if err != nil {
			t.Fatal(err)
		}
		if got, _ := d.GetString(obj.(*Dict), "Secret"); string(got) != "classified" {
			t.Errorf("decrypted string = %q", got)
		}
	})

	t.Run("V0 is refused", func(t *testing.T) {
		// Before the fix: `if err == nil` fires, because case 0 landed in the
		// same arm as case 1 and opened with a 5-byte RC4 key.
		data := replaceOnce(t, base, []byte("/V 1 /R 2"), []byte("/V 0 /R 2"))
		if data == nil {
			t.Fatal("could not locate a single /V 1 /R 2 in the fixture")
		}
		d, err := OpenBytes(data, nil)
		if err == nil {
			t.Fatalf("a document declaring the undocumented /V 0 opened: "+
				"Encryption() = %+v, file key = %d bytes", d.Encryption(), len(d.sec.key))
		}
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
	})

	t.Run("an /Encrypt with no /V at all is /V 0", func(t *testing.T) {
		// /V is optional with default 0 (dictInt(encDict, "V", 0)), so an
		// /Encrypt that omits it lands in the same arm. Renaming the key to /W
		// keeps every offset in place and leaves /O and /U, which make the
		// document open, untouched.
		data := replaceOnce(t, base, []byte("/V 1 /R 2"), []byte("/W 1 /R 2"))
		if data == nil {
			t.Fatal("could not locate a single /V 1 /R 2 in the fixture")
		}
		d, err := OpenBytes(data, nil)
		if err == nil {
			t.Fatalf("a document with no /V opened: Encryption() = %+v, file key = %d bytes",
				d.Encryption(), len(d.sec.key))
		}
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
	})
}

// TestEncryptionWritePathNeverRC4sAnAES256Key is the write-side regression for
// issue #33.
//
// encryptForWrite mirrors decryptBytes exactly, including `if !h.useAES {
// return rc4Apply(...) }`, so every read-side downgrade is also a write-side
// one: a signature added to such a document would be re-encrypted with RC4-128
// keyed from the AES-256 file key, self-consistent with our own reader and wrong
// for every other consumer. The invariant pinned here is the one the
// postcondition guard in setupEncryption asserts explicitly: a 32-byte file key
// - which, after computeEncryptionKey's S3 fix, can only come from the exact
// /R 5//R 6 /UE//OE unwrap, since computeKeyRev234 caps every other revision at
// MD5's 16 bytes - never reaches the RC4 branch of decryptBytes or
// encryptForWrite.
func TestEncryptionWritePathNeverRC4sAnAES256Key(t *testing.T) {
	plain := []byte("payload")

	for _, tc := range []struct {
		name string
		why  string
		data func(t *testing.T) []byte
	}{
		{"V5 with no /CF", "the read-side downgrade, seen from the writer", func(t *testing.T) []byte {
			return buildAES256Doc(t, aes256Spec{cfEntry: "/StmF /StdCF /StrF /StdCF", pValue: "-1052"})
		}},
		{"V4 /R 6 declaring /CFM /V2 (the crossbreed checkCryptFilters cannot see)",
			"/V 4 makes /CFM /V2 legitimate on its own; /R 6 still hands over the honest 32-byte /UE unwrap - " +
				"only the postcondition guard in setupEncryption, not checkCryptFilters, catches this shape",
			func(t *testing.T) []byte {
				return asV4R6(t, buildAES256Doc(t, aes256Spec{
					cfEntry: "/CF << /StdCF << /CFM /V2 /Length 32 >> >> /StmF /StdCF /StrF /StdCF",
					pValue:  "-1052"}))
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("why: %s", tc.why)
			d, err := OpenBytes(tc.data(t), nil)
			if err == nil {
				out, werr := d.encryptForWrite(ObjectKey{Num: 9}, plain, true)
				if werr != nil {
					t.Fatalf("the document opened with useAES = %v and a %d-byte key; "+
						"encryptForWrite: %v", d.sec.useAES, len(d.sec.key), werr)
				}
				rc4 := rc4Apply(d.sec.objectKeyFor(9, 0), plain)
				t.Fatalf("the document opened with useAES = %v and a %d-byte key, and "+
					"encryptForWrite emitted %d bytes; RC4 over that key emits %d, equal = %v "+
					"- an increment would be written under RC4-128 into a file declaring AES-256",
					d.sec.useAES, len(d.sec.key), len(out), len(rc4), bytes.Equal(out, rc4))
			}
			if !errors.Is(err, ErrUnsupportedSecurityHandler) {
				t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
			}
		})
	}

	t.Run("a conforming AES-256 document still writes AES", func(t *testing.T) {
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052"}),
			&Options{Random: fixedRandom()})
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		out, err := d.encryptForWrite(ObjectKey{Num: 9}, plain, true)
		if err != nil {
			t.Fatalf("encryptForWrite: %v", err)
		}
		if bytes.Equal(out, rc4Apply(d.sec.objectKeyFor(9, 0), plain)) {
			t.Fatal("encryptForWrite emitted RC4 over the AES-256 file key")
		}
		if len(out)%16 != 0 || len(out) < 32 {
			t.Errorf("encryptForWrite emitted %d bytes, want a 16-byte IV plus whole AES blocks", len(out))
		}
		if got := d.sec.decryptBytes(out, 9, 0); !bytes.Equal(got, plain) {
			t.Errorf("decrypt(encrypt(x)) = %q, want %q", got, plain)
		}
	})

	t.Run("an increment over an AES-256 document re-opens", func(t *testing.T) {
		// End to end through the Updater: R20 says the handler is never
		// upgraded, downgraded or re-keyed, so the string written here must come
		// back out of a fresh OpenBytes of original||increment.
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052"}),
			&Options{Random: fixedRandom()})
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		u, err := NewUpdater(d)
		if err != nil {
			t.Fatalf("NewUpdater: %v", err)
		}
		added := &Dict{}
		added.Set("Note", String{Bytes: []byte("added under AES-256")})
		k := u.Add(added)
		res, err := u.Write()
		if err != nil {
			t.Fatalf("Write: %v", err)
		}
		re, err := OpenBytes(res.Bytes, nil)
		if err != nil {
			t.Fatalf("re-opening the increment: %v", err)
		}
		obj, err := re.Object(k)
		if err != nil {
			t.Fatal(err)
		}
		if got, _ := re.GetString(obj.(*Dict), "Note"); string(got) != "added under AES-256" {
			t.Errorf("round-tripped string = %q, want %q", got, "added under AES-256")
		}
	})
}

// TestEncryptionRequiresExactUEOELength is the direct regression for S3:
// computeEncryptionKey used to accept /UE or /OE of any length that was a
// multiple of the AES block size, returning a "file key" exactly as long as the
// ciphertext (`out := make([]byte, len(fileKeyEnc))`). AES-CBC decrypts each
// block independent of the ones after it, so padding /UE past 32 bytes with
// arbitrary trailing blocks does not disturb the first two: out[0:32] is still
// the genuine AES-256 file key, sitting inside a longer slice that a
// `len(h.key) == 32` guard - objectKeyFor's fast path, and the postcondition
// guard in setupEncryption - would fail to recognise. ISO 32000-2 8.7.4.1 fixes
// both entries at exactly 32 bytes; requiring that in computeEncryptionKey turns
// the length itself into a reliable signal instead of attacker-controlled input.
//
// Honesty about what each subtest actually pins, end to end through OpenBytes:
// before S3 existed, only the 16-byte /UE case genuinely discriminated. A
// 48- or 64-byte /UE was never rejected by computeEncryptionKey - it decrypted
// to a same-length "file key" that validatePerms then handed to
// aes.NewCipher, which already rejects a non-{16,24,32}-byte key and already
// answers with a message containing "/UE or /OE" (see
// TestValidatePermsRejectsAWrongSizeFileKey in cryptaes256_test.go, which pins
// that fallback path directly now that it is unreachable from here). A
// 16-byte /UE, though, *is* a valid AES-128 key size, so pre-S3 it sailed past
// aes.NewCipher and only failed later at the 'adb' marker, in a message that
// does not mention /UE or /OE at all - the one case a control on message
// content alone would have caught. The "isolates computeEncryptionKey"
// subtest below is what closes that gap for the other lengths: by omitting
// /Perms it removes validatePerms from the path entirely, so nothing but
// computeEncryptionKey's own check can produce a rejection blaming /UE or /OE.
func TestEncryptionRequiresExactUEOELength(t *testing.T) {
	fileKey := bytes.Repeat([]byte{0x5A}, 32)

	for _, n := range []int{16, 48, 64} {
		t.Run(fmt.Sprintf("the user-password /UE branch, %d bytes", n), func(t *testing.T) {
			data := buildAES256Doc(t, aes256Spec{
				cfEntry: stdCFAESV3, pValue: "-1052",
				fileKey:  bytes.Repeat([]byte{0x5A}, n),
				permsRaw: make([]byte, 16), // bypass the builder's own AES-encrypt-of-fileKey step
			})
			_, err := OpenBytes(data, nil)
			if !errors.Is(err, ErrUnsupportedSecurityHandler) {
				t.Fatalf("/UE %d bytes: error is %v, want ErrUnsupportedSecurityHandler", n, err)
			}
			if got := err.Error(); !strings.Contains(got, "/UE or /OE") {
				t.Errorf("/UE %d bytes: error %q does not name /UE or /OE", n, got)
			}
		})
	}

	t.Run("the owner-password /OE branch, 48 bytes", func(t *testing.T) {
		ownerPw := []byte("owner-secret")
		data := buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052",
			fileKey:  bytes.Repeat([]byte{0x5A}, 48),
			ownerPw:  ownerPw,
			permsRaw: make([]byte, 16),
		})
		_, err := OpenBytes(data, &Options{Password: ownerPw})
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Fatalf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
		if got := err.Error(); !strings.Contains(got, "/UE or /OE") {
			t.Errorf("error %q does not name /UE or /OE", got)
		}
	})

	t.Run("computeEncryptionKey alone, isolated from validatePerms by omitting /Perms", func(t *testing.T) {
		// This is the subtest that actually pins the exact-32-byte check. The
		// ones above do not discriminate: validatePerms rejects a non-AES-size
		// file key too, with a message that also names /UE or /OE, so they pass
		// whether or not computeEncryptionKey checks anything.
		//
		// spec.omitPerms leaves /Perms out of the dictionary. validatePerms
		// still runs - its gate is `h.enc.R == 5 || h.enc.R == 6`, nothing
		// more - but computeEncryptionKey runs first, inside the password
		// switch, so its rejection is the one that escapes. Delete the
		// exact-32-byte check and the surviving error becomes validatePerms's
		// own "/Encrypt /Perms is 0 bytes, want 16", which is ErrInvalidPassword
		// and names /Perms rather than /UE or /OE - so both assertions below
		// fail, which is what makes this a real regression pin.
		data := buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052",
			fileKey:   bytes.Repeat([]byte{0x5A}, 48),
			permsRaw:  make([]byte, 16), // bypass the builder's own AES-encrypt-of-fileKey step
			omitPerms: true,             // and then never written: /Perms is absent
		})
		_, err := OpenBytes(data, nil)
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Fatalf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
		if got := err.Error(); !strings.Contains(got, "/UE or /OE") {
			t.Errorf("error %q does not name /UE or /OE", got)
		}
	})

	t.Run("a genuine 32-byte /UE still opens (positive control)", func(t *testing.T) {
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052", fileKey: fileKey}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if len(d.sec.key) != 32 {
			t.Errorf("file key = %d bytes, want 32", len(d.sec.key))
		}
	})
}

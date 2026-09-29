// AES-256 (/V 5 /R 6) fixtures, and the acceptance-boundary regression they exist
// for: the /Encrypt /Perms check of ISO 32000-2 Algorithm 13 (crypt.go
// validatePerms, the port of StandardSecurityHandler.validatePerms).
//
// buildRC4Document in crypt_test.go covers the /R 2 half of DESIGN.md §2.6 end to
// end; nothing covered the /R 6 half with bytes we control, because the only
// AES-256 input in the tree is the vendored corpus file. buildAES256Doc closes
// that gap: it derives /U, /UE, /O, /OE and /Perms with the same algorithms the
// reader checks them with, so a fixture can be tampered with one byte at a time.

package pdf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// stdCFAESV3 is the crypt-filter section a well-formed AES-256 document carries.
const stdCFAESV3 = "/CF << /StdCF << /CFM /AESV3 /Length 32 >> >> /StmF /StdCF /StrF /StdCF"

// aes256Spec parameterises buildAES256Doc. The zero value, with cfEntry and
// pValue filled in, is a conforming document that opens with the empty user
// password; every other field exists so one property at a time can be broken.
type aes256Spec struct {
	rev       int    // /R: 5 or 6; zero means 6
	cfEntry   string // spliced into /Encrypt verbatim
	pValue    string // the literal /P text, so an edit can be length-preserving
	ownerPw   []byte // nil: /O is 48 zero bytes and no owner password can match
	fileKey   []byte // nil: 32 bytes of 0x5A
	permsP    *int32 // nil: bind /Perms to pValue
	permsMeta *bool  // nil: bind /Perms to /EncryptMetadata true
	permsRaw  []byte // non-nil: write these bytes as /Perms instead of Algorithm 13's
	omitPerms bool   // true: leave /Perms out of the dictionary altogether

	// permsLiteral, when non-empty, is written as the /Perms value verbatim, so a
	// fixture can carry a /Perms that is not a string at all.
	permsLiteral string
	// permsIndirect writes /Perms as a reference to object 7 instead of inline.
	permsIndirect bool

	// pIndirect writes /P as a reference to object 8, holding pValue, instead of
	// the literal integer inline.
	pIndirect bool

	// encryptMetaValue, when non-empty ("true" or "false"), splices an
	// /EncryptMetadata entry into /Encrypt; empty leaves it absent, as before.
	// encryptMetaIndirect writes that value as a reference to object 9 instead
	// of inline.
	encryptMetaValue    string
	encryptMetaIndirect bool
}

func int32Ptr(v int32) *int32 { return &v }

// buildAES256Doc produces a real /Filter /Standard /V 5 /R 6 document that opens
// with the empty user password:
//
//	/U     = hash2A(pw, validation salt) || validation salt || key salt   (Algorithm 8)
//	/UE    = the file key under hash2A(pw, key salt), AES-256-CBC, zero IV
//	/O     = hash2A(ownerPw, validation salt, /U) || salts, or 48 zero bytes
//	/OE    = the file key under hash2A(ownerPw, key salt, /U)             (Algorithm 9)
//	/Perms = Algorithm 13's 16 bytes, AES-ECB under the file key
func buildAES256Doc(t *testing.T, spec aes256Spec) []byte {
	t.Helper()
	rev := spec.rev
	if rev == 0 {
		rev = 6
	}
	id := []byte("0123456789abcdef")

	var perm int32
	p64, err := strconv.ParseInt(spec.pValue, 10, 32)
	switch {
	case err == nil:
		perm = int32(p64)
	case spec.permsP != nil:
		// pValue is deliberately not an integer literal (a real, say); the
		// caller states separately what /Perms should be bound to.
		perm = *spec.permsP
	default:
		t.Fatalf("buildAES256Doc: /P %q: %v", spec.pValue, err)
	}

	fileKey := spec.fileKey
	if fileKey == nil {
		fileKey = bytes.Repeat([]byte{0x5A}, 32)
	}

	var userPw []byte // the empty user password
	user := append([]byte{}, hash2AOr256(userPw, []byte("VALIDSLT"), nil, rev)...)
	user = append(user, "VALIDSLT"...)
	user = append(user, "KEY-SALT"...)
	ue := aesCBCZeroIV(t, hash2AOr256(userPw, []byte("KEY-SALT"), nil, rev), fileKey)

	owner, oe := make([]byte, 48), []byte(nil)
	if spec.ownerPw != nil {
		owner = append([]byte{}, hash2AOr256(spec.ownerPw, []byte("OVALIDST"), user, rev)...)
		owner = append(owner, "OVALIDST"...)
		owner = append(owner, "OKEYSALT"...)
		oe = aesCBCZeroIV(t, hash2AOr256(spec.ownerPw, []byte("OKEYSALT"), user, rev), fileKey)
	}

	permsP, permsMeta := perm, true
	if spec.permsP != nil {
		permsP = *spec.permsP
	}
	if spec.permsMeta != nil {
		permsMeta = *spec.permsMeta
	}
	perms := spec.permsRaw
	if perms == nil {
		perms = aes256PermsEntry(t, fileKey, permsP, permsMeta)
	}
	permsEntry := fmt.Sprintf("/Perms <%s> ", hex.EncodeToString(perms))
	switch {
	case spec.omitPerms:
		permsEntry = ""
	case spec.permsLiteral != "":
		permsEntry = "/Perms " + spec.permsLiteral + " "
	case spec.permsIndirect:
		permsEntry = "/Perms 7 0 R "
	}
	oeEntry := ""
	if oe != nil {
		oeEntry = fmt.Sprintf("/OE <%s> ", hex.EncodeToString(oe))
	}

	pEntry := spec.pValue
	if spec.pIndirect {
		pEntry = "8 0 R"
	}
	encryptMetaEntry := ""
	if spec.encryptMetaValue != "" {
		if spec.encryptMetaIndirect {
			encryptMetaEntry = "/EncryptMetadata 9 0 R "
		} else {
			encryptMetaEntry = "/EncryptMetadata " + spec.encryptMetaValue + " "
		}
	}

	enc := aesEncryptorFor(t, fileKey)
	secret := enc(4, []byte("classified"), true)
	streamData := enc(6, []byte("stream payload"), false)

	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"},
		{num: 4, body: "<< /Secret <" + hex.EncodeToString(secret) + "> >>"},
		{num: 5, body: fmt.Sprintf("<< /Filter /Standard /V 5 /R %d /Length 256 %s /P %s %s"+
			"/O <%s> /U <%s> /UE <%s> %s%s>>",
			rev, spec.cfEntry, pEntry, encryptMetaEntry,
			hex.EncodeToString(owner), hex.EncodeToString(user),
			hex.EncodeToString(ue), oeEntry, permsEntry)},
		{num: 6, body: fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream",
			len(streamData), streamData)},
	}
	if spec.permsIndirect {
		objs = append(objs, rdrObj{num: 7, body: "<" + hex.EncodeToString(perms) + ">"})
	}
	if spec.pIndirect {
		objs = append(objs, rdrObj{num: 8, body: spec.pValue})
	}
	if spec.encryptMetaValue != "" && spec.encryptMetaIndirect {
		objs = append(objs, rdrObj{num: 9, body: spec.encryptMetaValue})
	}
	return buildReaderPDF("%PDF-1.7\n", objs,
		fmt.Sprintf("/Encrypt 5 0 R\n/ID [<%s> <%s>]\n",
			hex.EncodeToString(id), hex.EncodeToString(id)))
}

// aesEncryptorFor returns the writer's own AES encryptor bound to fileKey, so a
// fixture's strings and streams are produced by exactly the code decryptBytes
// inverts. The initialisation vectors come from fixedRandom, keeping every
// fixture byte-for-byte reproducible.
func aesEncryptorFor(t *testing.T, fileKey []byte) func(num int64, data []byte, isString bool) []byte {
	t.Helper()
	h := &securityHandler{key: fileKey, useAES: true, encryptMeta: true}
	h.enc.KeyLength = len(fileKey) * 8
	h.enc.StmF, h.enc.StrF = "StdCF", "StdCF"
	d := &Document{opts: Options{Random: fixedRandom()}.withDefaults(), sec: h}
	return func(num int64, data []byte, isString bool) []byte {
		out, err := d.encryptForWrite(ObjectKey{Num: num}, data, isString)
		if err != nil {
			t.Fatalf("encryptForWrite(object %d): %v", num, err)
		}
		return out
	}
}

// aesCBCZeroIV is the wrapping half of Algorithms 8 and 9: AES-CBC over a whole
// number of blocks with a zero initialisation vector and no padding.
func aesCBCZeroIV(t *testing.T, key, data []byte) []byte {
	t.Helper()
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	out := make([]byte, len(data))
	cipher.NewCBCEncrypter(block, make([]byte, aes.BlockSize)).CryptBlocks(out, data)
	return out
}

// aes256PermsEntry builds the /Encrypt /Perms entry of ISO 32000-2 Algorithm 13:
// /P as a little-endian int32, 0xff 0xff 0xff 0xff, 'T' or 'F' for
// /EncryptMetadata, the constant 'a' 'd' 'b', then four filler bytes — the whole
// 16 bytes encrypted with the file key, AES-ECB, no padding.
func aes256PermsEntry(t *testing.T, fileKey []byte, p int32, encryptMetadata bool) []byte {
	t.Helper()
	var plain [aes.BlockSize]byte
	binary.LittleEndian.PutUint32(plain[0:4], uint32(p))
	copy(plain[4:8], []byte{0xff, 0xff, 0xff, 0xff})
	plain[8] = 'F'
	if encryptMetadata {
		plain[8] = 'T'
	}
	copy(plain[9:12], "adb")
	copy(plain[12:16], "Pdf!") // "four bytes of arbitrary data", fixed here
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	out := make([]byte, aes.BlockSize)
	block.Encrypt(out, plain[:])
	return out
}

// canCreateSignatureField mirrors
// pades.NativePdfDocumentReader.CanCreateSignatureField, the predicate
// PdfPermissionsChecker gates signing on. internal/pdf must not import pades, so
// the conjunction is restated here; if that method changes, this must too.
func canCreateSignatureField(p Permissions) bool { return p.CanModify && p.CanModifyAnnots }

// TestEncryptionAES256PermsMustMatchP pins the /Encrypt /Perms check.
//
// For /R 5 and /R 6 the file key comes out of /UE (or /OE) alone: /P is not mixed
// into it, and neither isUserPassword nor computeEncryptionKey looks at it. The
// only thing binding /P to the key is the /Encrypt /Perms entry — 16 bytes of
// AES-ECB under the file key carrying /P in its first four bytes and the constant
// 'a' 'd' 'b' at bytes 9..11. Without that check anyone who can edit bytes
// rewrites /P, flips CanCreateSignatureField from false to true and walks through
// PdfPermissionsChecker.CheckDocumentPermissions.
//
// /P -1052 clears bit 4 (modify), /P -1028 sets it; both are five characters
// wide, so the edit moves no offset in the file.
func TestEncryptionAES256PermsMustMatchP(t *testing.T) {
	const honestP, tamperedP = "/P -1052", "/P -1028"

	t.Run("a conforming document opens and keeps its declared permissions", func(t *testing.T) {
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052"}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.V != 5 || enc.R != 6 || enc.CFM != "AESV3" || enc.KeyLength != 256 {
			t.Errorf("encryption = %+v", enc)
		}
		obj, err := d.Object(ObjectKey{Num: 4})
		if err != nil {
			t.Fatal(err)
		}
		if got, _ := d.GetString(obj.(*Dict), "Secret"); string(got) != "classified" {
			t.Errorf("decrypted string = %q", got)
		}
		sobj, err := d.Object(ObjectKey{Num: 6})
		if err != nil {
			t.Fatal(err)
		}
		if raw, _ := d.RawStreamData(sobj.(*Stream)); string(raw) != "stream payload" {
			t.Errorf("decrypted stream = %q", raw)
		}
		if p := d.Permissions(); p.OwnerAccess || p.CanModify || canCreateSignatureField(p) {
			t.Errorf("permissions = %+v, want user access with bit 4 clear", p)
		}
	})

	t.Run("a byte-edited /P is rejected", func(t *testing.T) {
		data := replaceOnce(t, buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052"}),
			[]byte(honestP), []byte(tamperedP))
		if data == nil {
			t.Fatal("could not locate a single /P -1052 in the fixture")
		}
		d, err := OpenBytes(data, nil)
		if err == nil {
			t.Fatalf("a document whose /P was edited from -1052 to -1028 opened: "+
				"permissions = %+v, CanCreateSignatureField() = %v (was false) "+
				"while /Perms still says -1052",
				d.Permissions(), canCreateSignatureField(d.Permissions()))
		}
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("a flipped /Perms ciphertext byte is rejected", func(t *testing.T) {
		// An attacker who edits /P cannot recompute /Perms without the file key,
		// so the 'a' 'd' 'b' constant is what makes the check unforgeable.
		perms := aes256PermsEntry(t, bytes.Repeat([]byte{0x5A}, 32), -1052, true)
		perms[0] ^= 0x01
		data := buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052", permsRaw: perms})
		if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("a /Perms carrying a different /P is rejected", func(t *testing.T) {
		// The 'adb' marker and the metadata byte are both correct here: only the
		// embedded /P disagrees, which is the comparison pdfbox logs and we refuse.
		other := int32(-1028)
		data := buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052", permsP: &other})
		if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("a /Perms disagreeing with /EncryptMetadata is rejected", func(t *testing.T) {
		// pdfbox compares byte 8 too. /EncryptMetadata decides whether metadata
		// streams are left in the clear, so an unauthenticated copy of it is worth
		// as little as an unauthenticated /P.
		metaFalse := false
		data := buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052", permsMeta: &metaFalse})
		if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("a missing /Perms is rejected", func(t *testing.T) {
		// /Perms is mandatory for /R 5 and /R 6. pdfbox hands the entry to
		// Cipher.doFinal unchecked, so a document without it fails there too — with
		// an IllegalArgumentException out of Loader.loadPDF that upstream DSS does
		// not map. We fail it as a typed error instead.
		data := buildAES256Doc(t, aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052", omitPerms: true})
		if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("a wrong-length /Perms is rejected, not indexed", func(t *testing.T) {
		// Five lengths spanning the four classes pdfbox treats differently: 0
		// bytes throws ArrayIndexOutOfBoundsException (0 is a legal AES block
		// count, so doFinal itself does not object); 12, 15, and 17 (not a
		// multiple of 16) throw IllegalBlockSizeException wrapped into
		// IOException; 32 (a larger multiple of 16) is the one length pdfbox does
		// not reject - validatePerms there reads only the first block - so this
		// port is stricter on that length, not merely at parity with pdfbox.
		for _, n := range []int{0, 12, 15, 17, 32} {
			data := buildAES256Doc(t, aes256Spec{
				cfEntry: stdCFAESV3, pValue: "-1052", permsRaw: make([]byte, n)})
			if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
				t.Errorf("/Perms of %d bytes: error is %v, want ErrInvalidPassword", n, err)
			}
		}
	})

	t.Run("a /Perms that is not a string is rejected, not indexed", func(t *testing.T) {
		for _, lit := range []string{"7", "/Perms", "[ ]", "<< >>", "null", "99 0 R"} {
			data := buildAES256Doc(t, aes256Spec{
				cfEntry: stdCFAESV3, pValue: "-1052", permsLiteral: lit})
			if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
				t.Errorf("/Perms %s: error is %v, want ErrInvalidPassword", lit, err)
			}
		}
	})

	t.Run("an indirect /Perms is resolved, as PDEncryption.getPerms does", func(t *testing.T) {
		// pdfbox reads /Perms with getDictionaryObject, which follows a reference.
		// The entry lives in /Encrypt, which is never itself decrypted, so what a
		// reference resolves to is still raw ciphertext.
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052", permsIndirect: true}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if p := d.Permissions(); p.CanModify {
			t.Errorf("permissions = %+v", p)
		}
	})

	t.Run("an indirect /P is resolved, as PDEncryption.getPermissions does", func(t *testing.T) {
		// getPermissions -> COSDictionary.getInt -> getDictionaryObject follows a
		// reference; validatePerms compares against /P, so /P itself must be read
		// the same way or an honest document with an indirect /P is rejected.
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052", pIndirect: true}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if p := d.Permissions(); p.CanModify {
			t.Errorf("permissions = %+v, want bit 4 (modify) clear", p)
		}
	})

	t.Run("a /P written as a real is read, as COSDictionary.getInt does", func(t *testing.T) {
		// getInt accepts any COSNumber, not just COSInteger. Requiring an
		// Integer here would leave p at 0 and reject an honest document at the
		// /Perms comparison, so a real-valued /P is read the same way.
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052.0", permsP: int32Ptr(-1052)}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if p := d.Permissions(); p.Raw != -1052 || p.CanModify {
			t.Errorf("permissions = %+v, want Raw -1052 with bit 4 (modify) clear", p)
		}
	})

	t.Run("an indirect /EncryptMetadata is resolved, as PDEncryption.isEncryptMetaData does", func(t *testing.T) {
		// isEncryptMetaData -> COSDictionary.getBoolean -> getDictionaryObject
		// follows a reference too. /Perms byte 8 is bound to /EncryptMetadata
		// false here so the comparison only passes if the reference is followed.
		metaFalse := false
		d, err := OpenBytes(buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052",
			encryptMetaValue: "false", encryptMetaIndirect: true, permsMeta: &metaFalse,
		}), nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if p := d.Permissions(); p.CanModify {
			t.Errorf("permissions = %+v, want bit 4 (modify) clear", p)
		}
	})

	t.Run("a file key that is not an AES key size is rejected, not a panic", func(t *testing.T) {
		// For /R 5 and /R 6 the file key is whatever /UE unwraps to.
		// computeEncryptionKey's own exact-32-byte check (issue #33's S3) is what
		// catches a 48-byte /UE now, before the malformed "key" it would have
		// produced ever reaches validatePerms's aes.NewCipher(fileKey) - a 48-byte
		// key is not a valid AES key size, and validatePerms turning that failure
		// into a rejection is exercised separately by
		// TestValidatePermsRejectsAWrongSizeFileKey below, which bypasses
		// computeEncryptionKey's check entirely to reach it. Either way the
		// outcome pinned here is the same: rejected cleanly, not a panic, and the
		// error blames /UE or /OE, the entry that is actually malformed - not
		// /Perms, which never gets read.
		//
		// The error is ErrUnsupportedSecurityHandler, not ErrInvalidPassword: a
		// malformed /UE is not a wrong password (see computeEncryptionKey's
		// // DIVERGENCE, deliberate: note in crypt.go), and pades maps the two
		// onto different exceptions - a caller retrying passwords against
		// ErrInvalidPassword would otherwise spin on a document no password can
		// open.
		data := buildAES256Doc(t, aes256Spec{
			cfEntry: stdCFAESV3, pValue: "-1052", fileKey: bytes.Repeat([]byte{0x5A}, 48),
			permsRaw: make([]byte, aes.BlockSize)})
		_, err := OpenBytes(data, nil)
		if !errors.Is(err, ErrUnsupportedSecurityHandler) {
			t.Errorf("error is %v, want ErrUnsupportedSecurityHandler", err)
		}
		// The bad key size comes from /UE or /OE, not from /Perms; the message
		// must blame the right entry.
		if err == nil || !strings.Contains(err.Error(), "/UE or /OE") || strings.Contains(err.Error(), "/Perms is") {
			t.Errorf("error is %q, want it to blame /UE or /OE, not /Perms", err)
		}
	})

	t.Run("the owner password path is validated too", func(t *testing.T) {
		ownerPw := []byte("owner-secret")
		spec := aes256Spec{cfEntry: stdCFAESV3, pValue: "-1052", ownerPw: ownerPw}

		d, err := OpenBytes(buildAES256Doc(t, spec), &Options{Password: ownerPw})
		if err != nil {
			t.Fatalf("the owner password no longer opens the fixture: %v", err)
		}
		if p := d.Permissions(); !p.OwnerAccess {
			t.Fatalf("permissions = %+v, want the owner branch", p)
		}

		data := replaceOnce(t, buildAES256Doc(t, spec), []byte(honestP), []byte(tamperedP))
		if data == nil {
			t.Fatal("could not locate a single /P -1052 in the fixture")
		}
		if _, err := OpenBytes(data, &Options{Password: ownerPw}); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("owner path: error is %v, want ErrInvalidPassword", err)
		}
	})

	t.Run("the vendored AES-256 fixture: honest opens, byte-edited /P does not", func(t *testing.T) {
		// The same attack on a real Acrobat-produced file, which is also the proof
		// that the check accepts a /Perms this port did not write.
		// corpusDir already skips this test when corpus/ is absent altogether;
		// a present corpus/ missing this specific fixture is drift between the
		// corpus and the code, not something to skip past silently.
		base, err := os.ReadFile(filepath.Join(corpusDir(t), "protected", "restricted_fields.pdf"))
		if err != nil {
			t.Fatalf("restricted_fields.pdf missing from a present corpus/: %v", err)
		}
		d, err := OpenBytes(base, nil)
		if err != nil {
			t.Fatalf("the vendored AES-256 fixture no longer opens: %v", err)
		}
		if p := d.Permissions(); canCreateSignatureField(p) {
			t.Fatalf("fixture permissions = %+v, want CanCreateSignatureField() false", p)
		}
		data := replaceOnce(t, base, []byte(honestP), []byte(tamperedP))
		if data == nil {
			t.Fatal("could not locate a single /P -1052 in the fixture")
		}
		d, err = OpenBytes(data, nil)
		if err == nil {
			t.Fatalf("the fixture opened with a byte-edited /P: permissions = %+v, "+
				"CanCreateSignatureField() = %v (was false)",
				d.Permissions(), canCreateSignatureField(d.Permissions()))
		}
		if !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})
}

// TestValidatePermsRejectsAWrongSizeFileKey pins validatePerms's own defence
// directly, calling it below OpenBytes: computeEncryptionKey's exact-32-byte
// /UE//OE check (issue #33's S3) means a document opened through OpenBytes can
// no longer hand validatePerms anything but a 32-byte file key, so
// validatePerms's own `aes.NewCipher(fileKey)` failure - a leftover from
// before S3 existed, when a malformed /UE could reach validatePerms with
// whatever length AES-CBC happened to decrypt - is unreachable end to end.
// It stays in place as defence in depth (validatePerms is not otherwise
// guaranteed only ever to be called with a 32-byte key), so it keeps its own
// direct test rather than losing coverage silently.
func TestValidatePermsRejectsAWrongSizeFileKey(t *testing.T) {
	fileKey := bytes.Repeat([]byte{0x5A}, 48) // not a valid AES key size
	err := validatePerms(make([]byte, aes.BlockSize), fileKey, -1052, true)
	if !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("error is %v, want ErrInvalidPassword", err)
	}
	if err == nil || !strings.Contains(err.Error(), "/UE or /OE") {
		t.Errorf("error is %q, want it to blame /UE or /OE", err)
	}
}

// TestPermissionWordMatchesJavaNarrowing pins permissionWord against Java's two
// narrowing conversions, which differ: COSInteger.intValue() is (int) of a long
// and keeps the low 32 bits, COSFloat.intValue() is (int) of a float and
// saturates with NaN mapped to zero. The unsigned row is the one that matters in
// practice — a producer that writes /P 4294966244 rather than /P -1052 depends
// on the truncation, and for /R 2 to /R 4 the value is hashed into Algorithm 2,
// so saturating there would derive a different file key than pdfbox.
func TestPermissionWordMatchesJavaNarrowing(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   Object
		want int32
	}{
		{"signed", Integer(-1052), -1052},
		{"the same word written unsigned", Integer(4294966244), -1052},
		{"positive", Integer(100), 100},
		{"minus one", Integer(-1), -1},
		{"above 32 bits keeps the low word", Integer(1<<32 + 5), 5},
		{"int32 max", Integer(math.MaxInt32), math.MaxInt32},
		{"int32 max plus one wraps", Integer(math.MaxInt32 + 1), math.MinInt32},
		{"a real is read, as getInt accepts any COSNumber", Real{Val: -1052}, -1052},
		{"a real saturates rather than wrapping", Real{Val: 1e18}, math.MaxInt32},
		{"a negative real saturates", Real{Val: -1e18}, math.MinInt32},
		{"NaN is zero", Real{Val: math.NaN()}, 0},
		{"absent is zero", nil, 0},
		{"a name is zero", Name("nope"), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := permissionWord(tc.in); got != tc.want {
				t.Errorf("permissionWord(%v) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestEncryptionAES256R5KnownAnswer pins the /R 5 password hash, the one
// AES-256 derivation the /R 6 fixtures above never reach: /R 5 is plain SHA-256
// (Adobe's Extension Level 3, ISO 32000-2 Algorithm 2.A without the Algorithm
// 2.B rounds). The expected /U, /O, /UE and /OE are computed here with
// crypto/sha256 directly, not through hash2AOr256, so a slip in the /R 5 arm (a
// wrong salt slice, the /R 6 hash pasted in, the user key left out of the owner
// hash) cannot be mirrored by the fixture builder and cancel out.
//
//	/U[0:32]  = SHA-256(userPw  || validation salt)
//	/O[0:32]  = SHA-256(ownerPw || validation salt || /U[0:48])
//	/UE key   = SHA-256(userPw  || key salt)
//	/OE key   = SHA-256(ownerPw || key salt || /U[0:48])
func TestEncryptionAES256R5KnownAnswer(t *testing.T) {
	ownerPw := []byte("owner-secret")
	data := buildAES256Doc(t, aes256Spec{rev: 5, cfEntry: stdCFAESV3, pValue: "-1052", ownerPw: ownerPw})

	entry := func(name string) []byte {
		t.Helper()
		i := bytes.Index(data, []byte("/"+name+" <"))
		if i < 0 {
			t.Fatalf("no /%s entry in the fixture", name)
		}
		rest := data[i+len(name)+3:]
		raw, err := hex.DecodeString(string(rest[:bytes.IndexByte(rest, '>')]))
		if err != nil {
			t.Fatalf("/%s: %v", name, err)
		}
		return raw
	}
	sha := func(parts ...[]byte) []byte {
		h := sha256.New()
		for _, p := range parts {
			h.Write(p)
		}
		return h.Sum(nil)
	}

	u, o, ue, oe := entry("U"), entry("O"), entry("UE"), entry("OE")
	if !bytes.Equal(u[:32], sha(nil, []byte("VALIDSLT"))) {
		t.Errorf("/U validation hash is not SHA-256(pw || salt)")
	}
	if !bytes.Equal(o[:32], sha(ownerPw, []byte("OVALIDST"), u[:48])) {
		t.Errorf("/O validation hash is not SHA-256(pw || salt || /U)")
	}
	fileKey := bytes.Repeat([]byte{0x5A}, 32)
	if !bytes.Equal(ue, aesCBCZeroIV(t, sha(nil, []byte("KEY-SALT")), fileKey)) {
		t.Errorf("/UE is not the file key under SHA-256(pw || key salt)")
	}
	if !bytes.Equal(oe, aesCBCZeroIV(t, sha(ownerPw, []byte("OKEYSALT"), u[:48]), fileKey)) {
		t.Errorf("/OE is not the file key under SHA-256(pw || key salt || /U)")
	}

	// The reader's own arm agrees with the independent computation.
	if got := hash2AOr256(ownerPw, []byte("OVALIDST"), u, 5); !bytes.Equal(got, sha(ownerPw, []byte("OVALIDST"), u[:48])) {
		t.Errorf("hash2AOr256 /R 5 (owner) = %x", got)
	}
	if got := hash2AOr256(nil, []byte("VALIDSLT"), nil, 5); !bytes.Equal(got, sha(nil, []byte("VALIDSLT"))) {
		t.Errorf("hash2AOr256 /R 5 (user) = %x", got)
	}

	secret := func(t *testing.T, d *Document) string {
		t.Helper()
		obj, err := d.Object(ObjectKey{Num: 4})
		if err != nil {
			t.Fatal(err)
		}
		got, _ := d.GetString(obj.(*Dict), "Secret")
		return string(got)
	}
	t.Run("opens with the user password", func(t *testing.T) {
		d, err := OpenBytes(data, nil)
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if enc := d.Encryption(); enc.V != 5 || enc.R != 5 || enc.CFM != "AESV3" {
			t.Errorf("encryption = %+v", enc)
		}
		if d.Permissions().OwnerAccess {
			t.Error("the empty user password must not grant owner access")
		}
		if got := secret(t, d); got != "classified" {
			t.Errorf("decrypted string = %q", got)
		}
	})
	t.Run("opens with the owner password", func(t *testing.T) {
		d, err := OpenBytes(data, &Options{Password: ownerPw})
		if err != nil {
			t.Fatalf("OpenBytes: %v", err)
		}
		if !d.Permissions().OwnerAccess {
			t.Error("the owner password must grant owner access")
		}
		if got := secret(t, d); got != "classified" {
			t.Errorf("decrypted string = %q", got)
		}
	})
	t.Run("refuses a wrong password", func(t *testing.T) {
		if _, err := OpenBytes(data, &Options{Password: []byte("nope")}); !errors.Is(err, ErrInvalidPassword) {
			t.Errorf("error is %v, want ErrInvalidPassword", err)
		}
	})
}

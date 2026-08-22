// The Standard security handler. See DESIGN.md §2.6.
//
// Provenance: org.apache.pdfbox.pdmodel.encryption.{SecurityHandler,
// StandardSecurityHandler,PDEncryption,AccessPermission} of pdfbox 3.0.7. The
// key derivations are ISO 32000-1 Algorithms 2/3/4/5 and ISO 32000-2 Algorithm
// 2.A/2.B; the deviations pdfbox takes from the spec (notably computeRC4key
// passing the key length into the 50-round MD5 loop) are copied deliberately —
// documents in the corpus only open with pdfbox's version.

package pdf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Encryption describes the document's security handler as configured.
type Encryption struct {
	Handler   Name // always "Standard" (others are rejected at Open)
	V, R      int
	KeyLength int
	StmF      Name // "StdCF" | "Identity"
	StrF      Name
	CFM       Name // "V2" | "AESV2" | "AESV3" | "None"
}

// clampDictInt narrows an /Encrypt entry to an int, saturating at the int32
// bounds. Every entry it is used on (/V, /R, /Length) has a small spec-defined
// range, so saturation cannot change the reading of a well-formed document; it
// only stops a hostile 64-bit value from wrapping into a plausible one where int
// is 32 bits wide.
func clampDictInt(v int64) int {
	switch {
	case v > math.MaxInt32:
		return math.MaxInt32
	case v < math.MinInt32:
		return math.MinInt32
	default:
		return int(v)
	}
}

// clampKeyLen bounds a /Length-derived key size to what an MD5 digest can
// supply. clampDictInt saturates /Length's magnitude but not its sign, and the
// /V 4 and /V 5 branches divide it by 8 without a floor, so a hostile /Encrypt
// reaches the Algorithm 2/3 digest slicing with an arbitrary int: /Length -8
// gives -1 and /Length 2000000000 gives 250000000, and both panic. Clamping to
// [0, digest] leaves every well-formed /Length untouched — the default is 40
// bits, or 5 bytes — and lets a malformed one fail closed on the password check
// instead of taking the process down.
func clampKeyLen(keyLen, digestLen int) int {
	switch {
	case keyLen < 0:
		return 0
	case keyLen > digestLen:
		return digestLen
	default:
		return keyLen
	}
}

// permissionWord reads a resolved /Encrypt /P. COSDictionary.getInt accepts any
// COSNumber, not just COSInteger, so a /P written as a real is read rather than
// ignored — for /R 5 and /R 6 ignoring it would mean rejecting the document at
// the /Perms check below.
//
// The two branches narrow differently because Java's two narrowing conversions
// do. COSInteger.intValue() is (int) of a long, which keeps the low 32 bits:
// that is not merely a hostile-input concern, since a producer that writes the
// permission word unsigned (/P 4294966244 for -1052) relies on it, and
// saturating instead would derive a different key for /R 2 to /R 4, where /P is
// hashed into Algorithm 2. The truncation is spelled out rather than left to an
// unchecked int64-to-int32 conversion. COSFloat.intValue() is (int) of a float,
// which saturates and maps NaN to zero — in Go that conversion is undefined for
// out-of-range values, so those cases are spelled out too.
func permissionWord(o Object) int32 {
	switch v := o.(type) {
	case Integer:
		low := int64(v) & 0xFFFFFFFF
		if low > math.MaxInt32 {
			low -= 1 << 32 // reinterpret the low 32 bits as two's complement
		}
		return int32(low)
	case Real:
		switch {
		case math.IsNaN(v.Val):
			return 0
		case v.Val >= math.MaxInt32:
			return math.MaxInt32
		case v.Val <= math.MinInt32:
			return math.MinInt32
		default:
			return int32(v.Val)
		}
	default:
		return 0
	}
}

// Permissions is the decoded /P bitfield plus which password matched.
type Permissions struct {
	Raw              int32
	OwnerAccess      bool // the owner password matched
	CanModify        bool // bit 4
	CanModifyAnnots  bool // bit 6
	CanFillInForm    bool // bit 9
	CanPrint         bool // bit 3
	CanExtract       bool // bit 5
	CanAssemble      bool // bit 11
	CanPrintFaithful bool // bit 12
}

func permissionsFrom(p int32, owner bool) Permissions {
	if owner {
		return Permissions{
			Raw: p, OwnerAccess: true, CanModify: true, CanModifyAnnots: true,
			CanFillInForm: true, CanPrint: true, CanExtract: true,
			CanAssemble: true, CanPrintFaithful: true,
		}
	}
	bit := func(n uint) bool { return p&(1<<(n-1)) != 0 }
	return Permissions{
		Raw:              p,
		CanPrint:         bit(3),
		CanModify:        bit(4),
		CanExtract:       bit(5),
		CanModifyAnnots:  bit(6),
		CanFillInForm:    bit(9),
		CanAssemble:      bit(11),
		CanPrintFaithful: bit(12),
	}
}

// encryptPadding is the 32-byte padding string of ISO 32000-1 Algorithm 2.
var encryptPadding = []byte{
	0x28, 0xBF, 0x4E, 0x5E, 0x4E, 0x75, 0x8A, 0x41, 0x64, 0x00, 0x4E, 0x56,
	0xFF, 0xFA, 0x01, 0x08, 0x2E, 0x2E, 0x00, 0xB6, 0xD0, 0x68, 0x3E, 0x80,
	0x2F, 0x0C, 0xA9, 0xFE, 0x64, 0x53, 0x69, 0x7A,
}

// securityHandler holds the derived file key and the per-object rules.
type securityHandler struct {
	enc         Encryption
	key         []byte
	useAES      bool
	perms       Permissions
	encryptMeta bool
	// encryptKey is the object key of the /Encrypt dictionary, which is never
	// itself decrypted.
	encryptKey ObjectKey
}

// setupEncryption prepares decryption when the trailer carries /Encrypt.
func (d *Document) setupEncryption() error {
	raw := d.trailer.GetRaw("Encrypt")
	if raw == nil {
		return nil
	}
	if r, ok := raw.(Ref); ok {
		d.encryptRef = r.Key()
	}
	encDict, _ := d.Resolve(raw).(*Dict)
	if encDict == nil {
		return nil
	}
	h := &securityHandler{encryptKey: d.encryptRef, encryptMeta: true}

	filter, _ := encDict.GetRaw("Filter").(Name)
	if filter != "Standard" {
		return fmt.Errorf("%w: /Filter /%s", ErrUnsupportedSecurityHandler, filter)
	}
	h.enc.Handler = filter
	h.enc.V = clampDictInt(dictInt(encDict, "V", 0))
	h.enc.R = clampDictInt(dictInt(encDict, "R", 0))
	length := clampDictInt(dictInt(encDict, "Length", 40))
	// isEncryptMetaData -> COSDictionary.getBoolean -> getDictionaryObject
	// follows /EncryptMetadata through a reference; d.sec is still nil here, so
	// resolving cannot recurse into decryption.
	if b, ok := d.Resolve(encDict.GetRaw("EncryptMetadata")).(Bool); ok {
		h.encryptMeta = bool(b)
	}
	keyLenBytes := length / 8
	h.enc.StmF, h.enc.StrF, h.enc.CFM = "Identity", "Identity", "None"

	switch h.enc.V {
	case 1:
		keyLenBytes = 5
		h.enc.StmF, h.enc.StrF, h.enc.CFM = "StdCF", "StdCF", "V2"
	case 2:
		if keyLenBytes <= 0 {
			keyLenBytes = 5
		}
		h.enc.StmF, h.enc.StrF, h.enc.CFM = "StdCF", "StdCF", "V2"
	case 4, 5:
		// PDEncryption.getStreamFilterName/getStringFilterName ->
		// COSDictionary.getCOSName -> getDictionaryObject follow a reference;
		// resolve here rather than reading GetRaw directly, so an indirect or
		// otherwise-typed /StmF or /StrF is not silently read as absent (and
		// therefore /Identity) the way a bare type assertion on GetRaw would.
		stmf, _ := d.Resolve(encDict.GetRaw("StmF")).(Name)
		strf, _ := d.Resolve(encDict.GetRaw("StrF")).(Name)
		if stmf == "" {
			stmf = "Identity"
		}
		if strf == "" {
			strf = "Identity"
		}
		h.enc.StmF, h.enc.StrF = stmf, strf
		cfName := stmf
		if cfName == "Identity" {
			cfName = strf
		}
		if cf, ok := d.Resolve(encDict.GetRaw("CF")).(*Dict); ok && cfName != "Identity" {
			if sub, ok := d.Resolve(cf.GetRaw(cfName)).(*Dict); ok {
				cfm, _ := sub.GetRaw("CFM").(Name)
				h.enc.CFM = cfm
				if l := clampDictInt(dictInt(sub, "Length", 0)); l > 0 {
					if l <= 40 {
						keyLenBytes = l // some files give bytes, some bits
					} else {
						keyLenBytes = l / 8
					}
				}
			}
		}
		switch h.enc.CFM {
		case "AESV2":
			h.useAES = true
			if keyLenBytes > 16 || keyLenBytes <= 0 {
				keyLenBytes = 16
			}
		case "AESV3":
			h.useAES = true
			keyLenBytes = 32
			if l := clampDictInt(dictInt(encDict, "Length", 0)); l > 0 && l/8 < 32 {
				keyLenBytes = l / 8
			}
		}
		// Runs before the password branch, so a document whose crypt filter
		// cannot be identified is refused with ErrUnsupportedSecurityHandler
		// rather than ErrInvalidPassword: the password was never the problem.
		if err := d.checkCryptFilters(encDict, h.enc); err != nil {
			return err
		}
	default:
		// DIVERGENCE, deliberate: pdfbox opens /V 0 - and, since /V is optional
		// with default 0 (dictInt(encDict, "V", 0) above), an /Encrypt that
		// omits /V altogether lands here too - as plain RC4 keyed by
		// /Length/8. StandardSecurityHandler.java:173 is
		// `int dicLength = encryptionVersion == 1 ? 5 : encryption.getLength() / 8;`,
		// so /V 0 takes the ternary's else, not the /V 1 special case; and
		// prepareForDecryption's `if (encryptionVersion >= REVISION_4)`
		// (line 157) is false for /V 0, so streamFilterName and
		// stringFilterName are never set away from null.
		// SecurityHandler.decryptStream/decryptString (lines 497, 628) test
		// `COSName.IDENTITY.equals(streamFilterName)` - false for a null
		// name - so both proceed and apply RC4 rather than leaving the bytes
		// alone. ISO 32000-1 Table 20 says of /V 0 only "An algorithm that is
		// undocumented ... shall not be used", so this port refuses it
		// instead of guessing at a key length pdfbox itself does not derive
		// consistently: the `keyLenBytes = 5` this arm used to share with
		// /V 1 under `case 0, 1:` was already a divergence from pdfbox's own
		// ternary, which gives /V 0 /Length/8 like every other unlisted /V,
		// not 5 - rejecting outright retires that silent disagreement too.
		return fmt.Errorf("%w: /V %d", ErrUnsupportedSecurityHandler, h.enc.V)
	}
	h.enc.KeyLength = keyLenBytes * 8

	o := stringBytes(encDict.GetRaw("O"))
	u := stringBytes(encDict.GetRaw("U"))
	oe := stringBytes(encDict.GetRaw("OE"))
	ue := stringBytes(encDict.GetRaw("UE"))
	// PDEncryption.getPermissions -> COSDictionary.getInt -> getDictionaryObject
	// follows /P through a reference; resolve here rather than through the
	// shared dictInt helper, which does not.
	p := permissionWord(d.Resolve(encDict.GetRaw("P")))
	id := d.ID()[0]

	password := d.opts.Password
	// pdfbox tries the OWNER password first (prepareForDecryption); DESIGN.md
	// §2.6 says "user, then owner". We follow pdfbox, because which one matched
	// decides Permissions.OwnerAccess and the KAT compares that field.
	switch {
	case isOwnerPassword(password, u, o, p, id, h.enc.R, keyLenBytes, h.encryptMeta):
		h.perms = permissionsFrom(p, true)
		pw := password
		if h.enc.R != 5 && h.enc.R != 6 {
			pw = userPasswordFromOwner(password, o, h.enc.R, keyLenBytes)
		}
		key, err := computeEncryptionKey(pw, o, u, oe, ue, p, id, h.enc.R, keyLenBytes, h.encryptMeta, true)
		if err != nil {
			return err
		}
		h.key = key
	case isUserPassword(password, u, o, p, id, h.enc.R, keyLenBytes, h.encryptMeta):
		h.perms = permissionsFrom(p, false)
		key, err := computeEncryptionKey(password, o, u, oe, ue, p, id, h.enc.R, keyLenBytes, h.encryptMeta, false)
		if err != nil {
			return err
		}
		h.key = key
	default:
		return ErrInvalidPassword
	}
	if h.enc.R == 4 && len(h.key) < 16 {
		// PDFBOX-5955: pad the RC4 key to 16 bytes.
		k := make([]byte, 16)
		copy(k, h.key)
		h.key = k
	}
	if h.enc.R == 5 || h.enc.R == 6 {
		// Same position and the same gate as pdfbox: after the key is installed,
		// on whichever password branch produced it, for /R 5 and /R 6 only.
		//
		// All three Algorithm 13 inputs - /Perms here, /P and /EncryptMetadata
		// above - follow references, as getDictionaryObject does; d.Resolve, not
		// the bare stringBytes(encDict.GetRaw(...)). /Encrypt is never itself
		// decrypted and d.sec is still nil here, so what comes back is still raw
		// ciphertext.
		if err := validatePerms(stringBytes(d.Resolve(encDict.GetRaw("Perms"))), h.key, p, h.encryptMeta); err != nil {
			return err
		}
	}
	// checkCryptFilters rejects an unidentifiable crypt filter, but it does not
	// - and, being a single-key single-useAES handler, cannot - stop a /V 4
	// document from legitimately declaring /CFM /V2 while /R 6 still hands over
	// the 32-byte file key unwrapped from /UE or /OE (issue #33's /V 4 //R 6
	// crossbreed). computeEncryptionKey now fixes that key at exactly 32 bytes
	// for /R 5 and /R 6 and at <=16 bytes for every other revision
	// (computeKeyRev234 slices an MD5 digest, never more than its 16 bytes), so
	// len(h.key) == 32 is an exact test for "this key came out of the /R 5//R 6
	// /UE or /OE unwrap" - asserted here, rather than left emergent, because
	// rc4Apply on both the read path (decryptBytes) and the write path
	// (encryptForWrite) trusts it never fires on such a key. Both filters
	// /Identity is exempt: decryptValue, decryptStream and encryptForWrite all
	// return before either sink is reached, so no cipher is chosen at all.
	//
	// This guard leaves one shape standing, and it is parity rather than a bug:
	// /V 4 with /R 5 or /R 6 and /CFM /AESV2, with a genuine 32-byte /UE//OE
	// unwrap, opens reporting {CFM: AESV2, KeyLength: 128} - the declaration is
	// honest AES-128 - while objectKeyFor's `useAES && len(h.key) == 32` fast
	// path (Algorithm 1.A) applies the file key directly as AES-256, exactly as
	// SecurityHandler.encryptData does at line 221. CFM and KeyLength describe
	// what the document declares, not always what cipher width is actually
	// applied; useAES and len(h.key) are what settle that, same as pdfbox.
	if !h.useAES && len(h.key) == 32 && (h.enc.StmF != "Identity" || h.enc.StrF != "Identity") {
		return fmt.Errorf("%w: /V %d /R %d recovers a 32-byte AES key but /CFM /%s selects RC4",
			ErrUnsupportedSecurityHandler, h.enc.V, h.enc.R, h.enc.CFM)
	}
	d.sec = h
	return nil
}

// validatePerms is ISO 32000-2 Algorithm 13, the /R 5//R 6 check that binds /P to
// the file key. /Encrypt /Perms is one AES block, ECB with no padding, under the
// file encryption key; it carries /P as a little-endian int32 at bytes 0-3, the
// /EncryptMetadata flag as 'T' or 'F' at byte 8, and the constant 'a' 'd' 'b' at
// bytes 9-11. Ports StandardSecurityHandler.validatePerms, which compares exactly
// those three fields: bytes 4-7 are the reserved 0xFF filler and bytes 12-15 are
// producer-chosen randomness, so neither is checked here either.
//
// DIVERGENCE, deliberate: pdfbox makes all three comparisons and answers every
// failure with LOG.warn("Verification of permissions failed ...") - no else, no
// return, no throw - so the document loads and AccessPermission is driven by the
// unauthenticated dictionary /P. Executed against pdfbox 3.0.7 on
// corpus/internal/pdf/testdata/corpus/protected/restricted_fields.pdf, rewriting
// its /P -1052 to /P -1028 logs the warning, loads, and reports canModify=true.
// For /R 5 and /R 6 /P is not mixed into the file key, so /Perms is the only
// thing that authenticates it, and pades.PdfPermissionsChecker's
// CanCreateSignatureField gate (CanModify && CanModifyAnnots) would be driven by
// that unauthenticated /P: warning and continuing would leave it defeated by a
// one-byte edit of a document under validation, which is hostile input. This
// port therefore rejects. An absent or wrong-length /Perms is rejected too, but
// pdfbox's fate per length is not uniformly "the same verdict, reached by
// crashing": absent, doFinal(null) throws IllegalArgumentException out of
// Loader.loadPDF; 0 bytes, doFinal returns an empty array and perms[9] throws
// ArrayIndexOutOfBoundsException (0 is a legal AES block count, not the
// IllegalBlockSizeException the other bad lengths get); 12, 15, 17 (not a
// multiple of 16), IllegalBlockSizeException wrapped into IOException. A
// larger multiple of 16 (32, 48, ...) is the one length class pdfbox does not
// crash on: doFinal succeeds and validatePerms reads only the first block, so
// a fabricated tail that hides a valid first block loads. This port accepts
// exactly one block and is therefore stricter than pdfbox on that class, not
// merely at parity with it.
func validatePerms(perms, fileKey []byte, p int32, encryptMetadata bool) error {
	if len(perms) != aes.BlockSize {
		return fmt.Errorf("%w: the password was accepted but /Encrypt /Perms is %d bytes, want %d",
			ErrInvalidPassword, len(perms), aes.BlockSize)
	}
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		// A key size error here means /UE or /OE unwrapped to the wrong length,
		// not that /Perms itself is malformed - say so, rather than blaming /Perms.
		return fmt.Errorf("%w: the password was accepted but the file key recovered from /UE or /OE is not an AES key size: %v",
			ErrInvalidPassword, err)
	}
	var out [aes.BlockSize]byte
	block.Decrypt(out[:], perms)
	if out[9] != 'a' || out[10] != 'd' || out[11] != 'b' {
		return fmt.Errorf("%w: the password was accepted but /Encrypt /Perms does not decrypt to the 'adb' marker, so it was not produced with this file key",
			ErrInvalidPassword)
	}
	if got := int32(binary.LittleEndian.Uint32(out[0:4])); got != p {
		return fmt.Errorf("%w: the password was accepted but /Encrypt /Perms says /P %d and the dictionary says %d",
			ErrInvalidPassword, got, p)
	}
	want := byte('F')
	if encryptMetadata {
		want = 'T'
	}
	if out[8] != want {
		return fmt.Errorf("%w: the password was accepted but /Encrypt /Perms says /EncryptMetadata %q and the dictionary says %q",
			ErrInvalidPassword, rune(out[8]), rune(want))
	}
	return nil
}

func stringBytes(o Object) []byte {
	if s, ok := o.(String); ok {
		return s.Bytes
	}
	return nil
}

// checkCryptFilters is the acceptance boundary of the /V 4 and /V 5 crypt-filter
// path: /V 5 must pair with /R 5 or /R 6 (ISO 32000-2 defines no other
// combination), every crypt filter the document actually *selects* (named by
// /StmF or /StrF, and not /Identity) must resolve through /CF to a /CFM this
// handler implements, every selected filter must resolve to the *same* one
// (the handler carries a single key and a single useAES flag, so it cannot
// honour two), for /V 5 that one must be /AESV3, and /AESV3 must itself pair
// with /R 5 or /R 6.
//
// DIVERGENCE, deliberate: pdfbox reads the crypt filter as a hardcoded /StdCF
// (PDEncryption.getStdCryptFilterDictionary) and only calls setAES(true) for
// /CFM /AESV2 or /AESV3, inside an `if (stdCryptFilterDictionary != null)` that
// has no else and no default - so a /V 5 document whose /CF is missing, whose
// /CF has no entry under the selected name, or whose /CFM is /V2, /None or
// anything unknown keeps useAES false, and SecurityHandler.encryptData then RC4s
// the 32 bytes unwrapped from /UE. Executed against real pdfbox 3.0.7 on three
// byte-edits of the corpus AES-256 fixture (/CF -> /XF, /StdCF -> /ZtdCF, /CFM
// /AESV3 -> /CFM /V2): all three cleared the R6 password check, applied RC4, and
// died downstream as "Page tree root must be a dictionary" - undefined,
// layout-dependent garbage that a differently laid out document would not
// produce, never a security-handler verdict. The bytes reaching this reader come
// off a document under signature validation, and this package is also the
// writer (encryptForWrite mirrors decryptBytes), so a silent RC4 fallback would
// both mis-report the cipher through Encryption() and re-emit a signed increment
// under RC4-128 in a file that declares AES-256. This port turns that undefined
// case into a defined rejection, which is the cheapest kind to justify. No
// corpus document is affected: the three /V 4 //V 5 fixtures all name /StdCF and
// resolve to /AESV2 or /AESV3; the fourth vendored encrypted file is /V 2 and
// never reaches this path.
//
// The /V 5-pairs-only-with-/R-5-or-6 check closes a shape none of the above
// catches: computeEncryptionKey and the postcondition guard both key off /R,
// not /V, so a /V 5 /R 4 document whose /CF resolves /StmF and /StrF to
// /AESV3 - a perfectly legitimate /V 5 crypt filter by every rule above -
// still takes computeKeyRev234's MD5 path (r == 4, not 5 or 6) and derives a
// <=16-byte key the way a genuine /V 4 /R 4 document would, while /Encrypt
// /CFM /AESV3 and useAES == true both say AES-256. That is exact pdfbox
// parity (dicRevision, not encryptionVersion, drives StandardSecurityHandler
// past prepareForDecryption:175), not a new divergence: pdfbox derives the
// same MD5-based key for the same bytes and, since useAES is true, also
// applies AES rather than RC4 to it - but at whatever length keyLenBytes came
// out to (16, in the /Length 128 case), not the 32 bytes /AESV3 promises. It
// is nonetheless the last surviving "declares AES-256, applies AES-128" shape
// issue #33 exists to close, so this port refuses it rather than reproduce
// it.
//
// checkCryptFilters is a lookup only: it does not feed back into which name
// setupEncryption derives its own CFM, useAES and key length from - that is
// still /StmF, falling back to /StrF only when /StmF is /Identity, the
// pre-existing divergence from pdfbox's hardcoded /StdCF noted above. This check
// neither widens nor narrows it; it only adds the "every selected filter must
// agree" refusal, which is a limitation of *this* implementation (one key, one
// useAES flag) rather than a pdfbox divergence - ISO 32000-1 permits distinct
// crypt filters for streams and strings, pdfbox just never notices because it
// never resolves the one it did not pick.
func (d *Document) checkCryptFilters(encDict *Dict, enc Encryption) error {
	// ISO 32000-2 pairs /V 5 with /R 5 or /R 6 only. Checked unconditionally,
	// not just when a crypt filter is selected: computeEncryptionKey and the
	// postcondition guard in setupEncryption both key off /R alone, so a /V 5
	// /R 4 document - even one whose /CF resolves to a legitimate /AESV3 -
	// would derive its key and choose its cipher length as if it were /V 4
	// /R 4, applying AES-128 to a document declaring AES-256.
	if enc.V == 5 && enc.R != 5 && enc.R != 6 {
		return fmt.Errorf("%w: /V 5 /R %d, and /V 5 pairs only with /R 5 or /R 6",
			ErrUnsupportedSecurityHandler, enc.R)
	}
	var selected Name
	for _, f := range []struct {
		entry Name // "StmF" or "StrF"
		name  Name // the crypt filter it names
	}{{"StmF", enc.StmF}, {"StrF", enc.StrF}} {
		// ISO 32000-2 Table 20: /Identity is the default for both entries and
		// means "leave the bytes alone". No cipher is selected, so there is
		// nothing to identify.
		if f.name == "Identity" {
			continue
		}
		cfm := d.cryptFilterMethod(encDict, f.name)
		if cfm == "" {
			return fmt.Errorf("%w: /%s names the crypt filter /%s, which does not resolve through /CF to a /CFM",
				ErrUnsupportedSecurityHandler, f.entry, f.name)
		}
		switch cfm {
		case "V2", "AESV2", "AESV3":
		default:
			return fmt.Errorf("%w: /%s /%s declares /CFM /%s",
				ErrUnsupportedSecurityHandler, f.entry, f.name, cfm)
		}
		if enc.V == 5 && cfm != "AESV3" {
			return fmt.Errorf("%w: /V 5 /%s /%s declares /CFM /%s, and /V 5 is /AESV3 only",
				ErrUnsupportedSecurityHandler, f.entry, f.name, cfm)
		}
		// The converse, and for the same reason as the /V 5 /R check above:
		// /AESV3 promises a 32-byte key, and the only derivation that produces
		// one is the /UE //OE unwrap, which computeEncryptionKey selects on /R
		// alone. A /V 4 /R 4 document naming /CFM /AESV3 therefore reports
		// /AESV3 with /KeyLength 256 while applying AES-128 under an MD5 key.
		// Without this rule the /V 5 check above is one-directional and the
		// same "declares AES-256, applies AES-128" shape simply moves to /V 4.
		if cfm == "AESV3" && enc.R != 5 && enc.R != 6 {
			return fmt.Errorf("%w: /%s /%s declares /CFM /AESV3 under /R %d, and /AESV3 needs the 32-byte key only /R 5 or /R 6 derives",
				ErrUnsupportedSecurityHandler, f.entry, f.name, enc.R)
		}
		if selected == "" {
			selected = cfm
		} else if cfm != selected {
			return fmt.Errorf("%w: /StmF and /StrF select different crypt filter methods, /%s and /%s",
				ErrUnsupportedSecurityHandler, selected, cfm)
		}
	}
	return nil
}

// cryptFilterMethod resolves one crypt filter name through /CF to its /CFM,
// returning "" when /CF, the named entry or /CFM is missing. It is a lookup
// only; see checkCryptFilters's doc comment for why that is what keeps the
// pre-existing /StmF-then-/StrF divergence from widening.
func (d *Document) cryptFilterMethod(encDict *Dict, name Name) Name {
	cf, _ := d.Resolve(encDict.GetRaw("CF")).(*Dict)
	if cf == nil {
		return ""
	}
	sub, _ := d.Resolve(cf.GetRaw(name)).(*Dict)
	if sub == nil {
		return ""
	}
	cfm, _ := sub.GetRaw("CFM").(Name)
	return cfm
}

// --- password checks -------------------------------------------------------

func isUserPassword(pw, u, o []byte, p int32, id []byte, r, keyLen int, encMeta bool) bool {
	if r == 5 || r == 6 {
		if len(u) < 48 {
			return false
		}
		hash := hash2AOr256(truncate127(saslPrepMaybe(pw, r)), u[32:40], nil, r)
		return bytes.Equal(hash, u[:32])
	}
	computed := computeUserEntry(pw, o, p, id, r, keyLen, encMeta)
	if r == 2 {
		return bytes.Equal(u, computed)
	}
	if len(u) < 16 || len(computed) < 16 {
		return false
	}
	return bytes.Equal(u[:16], computed[:16])
}

func isOwnerPassword(pw, u, o []byte, p int32, id []byte, r, keyLen int, encMeta bool) bool {
	if r == 5 || r == 6 {
		if len(o) < 48 || len(u) < 48 {
			return false
		}
		hash := hash2AOr256(truncate127(saslPrepMaybe(pw, r)), o[32:40], u, r)
		return bytes.Equal(hash, o[:32])
	}
	user := userPasswordFromOwner(pw, o, r, keyLen)
	return isUserPassword(user, u, o, p, id, r, keyLen, encMeta)
}

// userPasswordFromOwner is StandardSecurityHandler.getUserPassword234.
func userPasswordFromOwner(pw, o []byte, r, keyLen int) []byte {
	rc4Key := computeRC4Key(pw, r, keyLen)
	switch r {
	case 2:
		return rc4Apply(rc4Key, o)
	default:
		buf := append([]byte(nil), o...)
		iter := make([]byte, len(rc4Key))
		for i := 19; i >= 0; i-- {
			copy(iter, rc4Key)
			for j := range iter {
				iter[j] ^= byte(i)
			}
			buf = rc4Apply(iter, buf)
		}
		return buf
	}
}

// computeRC4Key is StandardSecurityHandler.computeRC4key. Note the documented
// deviation from the spec: the 50-round loop hashes only `length` bytes.
func computeRC4Key(pw []byte, r, keyLen int) []byte {
	sum := md5.Sum(truncateOrPad(pw))
	digest := sum[:]
	keyLen = clampKeyLen(keyLen, len(digest))
	if r == 3 || r == 4 {
		for i := 0; i < 50; i++ {
			s := md5.Sum(digest[:keyLen])
			digest = s[:]
		}
	}
	return append([]byte(nil), digest[:keyLen]...)
}

// computeUserEntry is StandardSecurityHandler.computeUserPassword (Algorithm 4/5).
func computeUserEntry(pw, o []byte, p int32, id []byte, r, keyLen int, encMeta bool) []byte {
	encKey := computeKeyRev234(pw, o, p, id, encMeta, keyLen, r)
	if r == 2 {
		return rc4Apply(encKey, encryptPadding)
	}
	h := md5.New()
	h.Write(encryptPadding)
	h.Write(id)
	buf := h.Sum(nil)
	iter := make([]byte, len(encKey))
	for i := 0; i < 20; i++ {
		copy(iter, encKey)
		for j := range iter {
			iter[j] ^= byte(i)
		}
		buf = rc4Apply(iter, buf)
	}
	out := make([]byte, 32)
	copy(out, buf[:16])
	copy(out[16:], encryptPadding[:16])
	return out
}

// computeKeyRev234 is Algorithm 2.
func computeKeyRev234(pw, o []byte, p int32, id []byte, encMeta bool, keyLen, r int) []byte {
	h := md5.New()
	h.Write(truncateOrPad(pw))
	h.Write(o)
	var pb [4]byte
	binary.LittleEndian.PutUint32(pb[:], uint32(p))
	h.Write(pb[:])
	h.Write(id)
	if r == 4 && !encMeta {
		h.Write([]byte{0xff, 0xff, 0xff, 0xff})
	}
	digest := h.Sum(nil)
	keyLen = clampKeyLen(keyLen, len(digest))
	if r == 3 || r == 4 {
		for i := 0; i < 50; i++ {
			s := md5.Sum(digest[:keyLen])
			digest = s[:]
		}
	}
	return append([]byte(nil), digest[:keyLen]...)
}

func computeEncryptionKey(pw, o, u, oe, ue []byte, p int32, id []byte, r, keyLen int, encMeta, ownerPw bool) ([]byte, error) {
	if r != 5 && r != 6 {
		return computeKeyRev234(pw, o, p, id, encMeta, keyLen, r), nil
	}
	var hash, fileKeyEnc []byte
	pw = truncate127(saslPrepMaybe(pw, r))
	if ownerPw {
		if len(oe) == 0 || len(o) < 48 {
			return nil, fmt.Errorf("%w: /Encrypt /OE entry is missing", ErrUnsupportedSecurityHandler)
		}
		hash = hash2AOr256(pw, o[40:48], u, r)
		fileKeyEnc = oe
	} else {
		if len(ue) == 0 || len(u) < 48 {
			return nil, fmt.Errorf("%w: /Encrypt /UE entry is missing", ErrUnsupportedSecurityHandler)
		}
		hash = hash2AOr256(pw, u[40:48], nil, r)
		fileKeyEnc = ue
	}
	block, err := aes.NewCipher(hash)
	if err != nil {
		return nil, err
	}
	// DIVERGENCE, deliberate: ISO 32000-2 8.7.4.1 fixes /UE and /OE at exactly
	// 32 bytes - the wrapped AES-256 file key, two whole AES blocks and nothing
	// else. pdfbox is equally lenient here only over lengths that are already a
	// whole multiple of the AES block size: StandardSecurityHandler hands
	// fileKeyEnc straight to `Cipher.getInstance("AES/CBC/NoPadding").doFinal
	// (fileKeyEnc)` (line 827), and AES/CBC/NoPadding throws
	// IllegalBlockSizeException for anything else, caught at line 829 and
	// rethrown as IOException - so pdfbox does reject a 33-byte /UE, just not a
	// 48-byte one. This port cannot afford even that narrower leniency:
	// AES-CBC decryption of one block never depends on the ciphertext blocks
	// after it, so a /UE padded past 32 bytes with further whole blocks still
	// decrypts its first two blocks unchanged, and out[0:32] is still the
	// genuine 32-byte file key sitting inside a longer slice. Every
	// len(h.key) == 32 fast path downstream - objectKeyFor's AES-256 branch, and
	// the postcondition guard setupEncryption asserts just before d.sec is
	// published - assumes 32 bytes means "the R5/R6 AES unwrap", and only a
	// hard length check here makes that assumption sound rather than bypassable
	// by padding or truncating /UE or /OE by exactly one AES block.
	//
	// The error is ErrUnsupportedSecurityHandler, not ErrInvalidPassword: a
	// malformed /UE or /OE is not a wrong password (validatePerms's own
	// aes.NewCipher failure, above in this file, is the "password was accepted
	// but the recovered key is unusable" case, and stays ErrInvalidPassword for
	// that reason), and pades maps the two onto different exceptions
	// (InvalidPasswordException vs ProtectedDocumentException) - a caller
	// retrying passwords against ErrInvalidPassword would otherwise spin on a
	// document no password can ever open.
	if len(fileKeyEnc) != 32 {
		return nil, fmt.Errorf("%w: /UE or /OE is %d bytes, ISO 32000-2 fixes it at exactly 32",
			ErrUnsupportedSecurityHandler, len(fileKeyEnc))
	}
	out := make([]byte, 32)
	cipher.NewCBCDecrypter(block, make([]byte, 16)).CryptBlocks(out, fileKeyEnc)
	return out, nil
}

// hash2AOr256 dispatches between R5's plain SHA-256 and R6's Algorithm 2.B.
func hash2AOr256(password, salt, u []byte, r int) []byte {
	userKey := adjustUserKey(u)
	if r == 5 {
		h := sha256.New()
		h.Write(password)
		h.Write(salt)
		h.Write(userKey)
		return h.Sum(nil)
	}
	input := make([]byte, 0, len(password)+len(salt)+len(userKey))
	input = append(input, password...)
	input = append(input, salt...)
	input = append(input, userKey...)
	return computeHash2B(input, password, userKey)
}

func adjustUserKey(u []byte) []byte {
	if len(u) == 0 {
		return nil
	}
	if len(u) < 48 {
		return nil
	}
	return u[:48]
}

// computeHash2B is ISO 32000-2 Algorithm 2.B, transcribed from
// StandardSecurityHandler.computeHash2B.
func computeHash2B(input, password, userKey []byte) []byte {
	sum := sha256.Sum256(input)
	k := sum[:]
	var e []byte
	for round := 0; round < 64 || int(e[len(e)-1]) > round-32; round++ {
		unit := len(password) + len(k)
		if len(userKey) >= 48 {
			unit += 48
		}
		k1 := make([]byte, 0, 64*unit)
		for i := 0; i < 64; i++ {
			k1 = append(k1, password...)
			k1 = append(k1, k...)
			if len(userKey) >= 48 {
				k1 = append(k1, userKey[:48]...)
			}
		}
		block, err := aes.NewCipher(k[:16])
		if err != nil {
			return k
		}
		e = make([]byte, len(k1))
		cipher.NewCBCEncrypter(block, k[16:32]).CryptBlocks(e, k1)
		mod := new(big.Int).Mod(new(big.Int).SetBytes(e[:16]), big.NewInt(3)).Int64()
		switch mod {
		case 0:
			s := sha256.Sum256(e)
			k = s[:]
		case 1:
			s := sha512.Sum384(e)
			k = s[:]
		default:
			s := sha512.Sum512(e)
			k = s[:]
		}
	}
	if len(k) > 32 {
		return k[:32]
	}
	return k
}

// saslPrepMaybe applies SASLprep for R6. Our corpus passwords are ASCII, so the
// full stringprep profile is not implemented: ASCII input is returned unchanged,
// which is what saslPrepQuery does for it.
func saslPrepMaybe(pw []byte, r int) []byte {
	if r != 6 {
		return pw
	}
	for _, c := range pw {
		if c >= 0x80 {
			// Non-ASCII password with R6: pass through unmapped. Flagged in
			// DESIGN notes; no corpus document exercises it.
			return pw
		}
	}
	return pw
}

func truncate127(b []byte) []byte {
	if len(b) <= 127 {
		return b
	}
	return b[:127]
}

func truncateOrPad(pw []byte) []byte {
	out := make([]byte, 32)
	n := copy(out, pw)
	copy(out[n:], encryptPadding)
	return out
}

func rc4Apply(key, data []byte) []byte {
	if len(key) == 0 {
		return append([]byte(nil), data...)
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return append([]byte(nil), data...)
	}
	out := make([]byte, len(data))
	c.XORKeyStream(out, data)
	return out
}

// --- per-object decryption -------------------------------------------------

// objectKeyFor is SecurityHandler.calcFinalKey (Algorithm 1). AES-256 uses the
// file key directly (Algorithm 1.A).
func (h *securityHandler) objectKeyFor(num int64, gen uint16) []byte {
	if h.useAES && len(h.key) == 32 {
		return h.key
	}
	newKey := make([]byte, len(h.key)+5)
	copy(newKey, h.key)
	n := len(newKey)
	newKey[n-5] = byte(num)
	newKey[n-4] = byte(num >> 8)
	newKey[n-3] = byte(num >> 16)
	newKey[n-2] = byte(gen)
	newKey[n-1] = byte(gen >> 8)
	md := md5.New()
	md.Write(newKey)
	if h.useAES {
		md.Write([]byte{0x73, 0x41, 0x6c, 0x54}) // "sAlT"
	}
	digest := md.Sum(nil)
	l := len(newKey)
	if l > 16 {
		l = 16
	}
	return digest[:l]
}

// decryptBytes decrypts one string or stream payload.
func (h *securityHandler) decryptBytes(data []byte, num int64, gen uint16) []byte {
	if len(h.key) == 0 {
		return data
	}
	key := h.objectKeyFor(num, gen)
	if !h.useAES {
		// The read half of issue #33's sink. checkCryptFilters and the
		// postcondition guard in setupEncryption are what keep a 32-byte
		// AES-256 file key out of here; nothing on this per-object path
		// re-checks it.
		return rc4Apply(key, data)
	}
	if len(data) < aes.BlockSize {
		return []byte{}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return data
	}
	iv := data[:aes.BlockSize]
	body := data[aes.BlockSize:]
	body = body[:len(body)-len(body)%aes.BlockSize]
	if len(body) == 0 {
		return []byte{}
	}
	out := make([]byte, len(body))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, body)
	// PKCS#5 padding, stripped leniently: pdfbox's CipherInputStream swallows a
	// BadPaddingException and keeps whatever it decoded.
	if n := int(out[len(out)-1]); n >= 1 && n <= aes.BlockSize && n <= len(out) {
		out = out[:len(out)-n]
	}
	return out
}

// decryptObject applies SecurityHandler.decrypt to a freshly parsed indirect
// object, in place. Objects taken out of an object stream are never passed here:
// their container was decrypted as a whole.
func (d *Document) decryptObject(obj Object, key ObjectKey) Object {
	h := d.sec
	if h == nil || key == h.encryptKey {
		return obj
	}
	return d.decryptValue(obj, key, 0)
}

func (d *Document) decryptValue(obj Object, key ObjectKey, depth int) Object {
	if depth > d.opts.MaxDepth {
		return obj
	}
	h := d.sec
	switch v := obj.(type) {
	case String:
		if h.enc.StrF == "Identity" {
			return v
		}
		return String{Bytes: h.decryptBytes(v.Bytes, key.Num, key.Gen), Hex: v.Hex}
	case Array:
		for i, e := range v {
			v[i] = d.decryptValue(e, key, depth+1)
		}
		return v
	case *Dict:
		d.decryptDict(v, key, depth)
		return v
	case *Stream:
		d.decryptStream(v, key, depth)
		return v
	default:
		return obj
	}
}

// decryptDict is SecurityHandler.decryptDictionary, including the signature
// guard: the /Contents of a /Sig or /DocTimeStamp dictionary — and, per
// PDFBOX-4466 (raised as DSS-1538), of any dictionary that pairs a /Contents
// string with a /ByteRange array — is never decrypted.
func (d *Document) decryptDict(dict *Dict, key ObjectKey, depth int) {
	if dict.Has("CF") {
		// PDFBOX-2936: avoid orphan /CF dictionaries in US govt "I-" files.
		return
	}
	typ, _ := dict.GetRaw("Type").(Name)
	_, contentsIsString := dict.GetRaw("Contents").(String)
	_, byteRangeIsArray := dict.GetRaw("ByteRange").(Array)
	isSignature := typ == "Sig" || typ == "DocTimeStamp" || (contentsIsString && byteRangeIsArray)
	for _, k := range dict.Keys() {
		if isSignature && k == "Contents" {
			continue
		}
		switch v := dict.GetRaw(k).(type) {
		case String, Array, *Dict:
			dict.Set(k, d.decryptValue(v, key, depth+1))
		}
	}
}

// decryptStream is SecurityHandler.decryptStream.
func (d *Document) decryptStream(s *Stream, key ObjectKey, depth int) {
	h := d.sec
	if s.decrypted {
		return
	}
	s.decrypted = true
	typ, _ := s.Dict.GetRaw("Type").(Name)
	d.decryptDict(s.Dict, key, depth)
	if h.enc.StmF == "Identity" {
		return
	}
	if typ == "XRef" {
		// "The cross-reference stream shall not be encrypted."
		return
	}
	if typ == "Metadata" {
		if !h.encryptMeta {
			return
		}
		// PDFBOX-3229: metadata that is plainly not encrypted despite the flag.
		if len(s.Raw) >= 10 && string(s.Raw[:10]) == "<?xpacket " {
			return
		}
	}
	names, _ := streamFilters(s.Dict, d.Resolve)
	for _, n := range names {
		if canonicalFilterName(n) == FilterCrypt {
			return // /Crypt /Identity
		}
	}
	s.Raw = h.decryptBytes(s.Raw, key.Num, key.Gen)
	s.Length = int64(len(s.Raw))
}

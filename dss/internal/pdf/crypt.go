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
	h.enc.V = int(dictInt(encDict, "V", 0))
	h.enc.R = int(dictInt(encDict, "R", 0))
	length := int(dictInt(encDict, "Length", 40))
	if b, ok := encDict.GetRaw("EncryptMetadata").(Bool); ok {
		h.encryptMeta = bool(b)
	}
	keyLenBytes := length / 8
	h.enc.StmF, h.enc.StrF, h.enc.CFM = "Identity", "Identity", "None"

	switch h.enc.V {
	case 0, 1:
		keyLenBytes = 5
		h.enc.StmF, h.enc.StrF, h.enc.CFM = "StdCF", "StdCF", "V2"
	case 2:
		if keyLenBytes <= 0 {
			keyLenBytes = 5
		}
		h.enc.StmF, h.enc.StrF, h.enc.CFM = "StdCF", "StdCF", "V2"
	case 4, 5:
		stmf, _ := encDict.GetRaw("StmF").(Name)
		strf, _ := encDict.GetRaw("StrF").(Name)
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
				if l := dictInt(sub, "Length", 0); l > 0 {
					if l <= 40 {
						keyLenBytes = int(l) // some files give bytes, some bits
					} else {
						keyLenBytes = int(l) / 8
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
			if l := dictInt(encDict, "Length", 0); l > 0 && int(l)/8 < 32 {
				keyLenBytes = int(l) / 8
			}
		}
	default:
		return fmt.Errorf("%w: /V %d", ErrUnsupportedSecurityHandler, h.enc.V)
	}
	h.enc.KeyLength = keyLenBytes * 8

	o := stringBytes(encDict.GetRaw("O"))
	u := stringBytes(encDict.GetRaw("U"))
	oe := stringBytes(encDict.GetRaw("OE"))
	ue := stringBytes(encDict.GetRaw("UE"))
	p := int32(dictInt(encDict, "P", 0))
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
	d.sec = h
	return nil
}

func stringBytes(o Object) []byte {
	if s, ok := o.(String); ok {
		return s.Bytes
	}
	return nil
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
	if r == 3 || r == 4 {
		for i := 0; i < 50; i++ {
			s := md5.Sum(digest[:keyLen])
			digest = s[:]
		}
	}
	if keyLen > len(digest) {
		keyLen = len(digest)
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
	if r == 3 || r == 4 {
		if keyLen > len(digest) {
			keyLen = len(digest)
		}
		for i := 0; i < 50; i++ {
			s := md5.Sum(digest[:keyLen])
			digest = s[:]
		}
	}
	if keyLen > len(digest) {
		keyLen = len(digest)
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
	if len(fileKeyEnc)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("%w: /UE or /OE is not a multiple of the AES block size", ErrInvalidPassword)
	}
	out := make([]byte, len(fileKeyEnc))
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

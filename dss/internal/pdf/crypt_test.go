package pdf

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

// buildRC4Document produces a real /Filter /Standard V1 R2 (RC4 40-bit)
// document: /O and /U are computed with the same algorithms the reader uses to
// check them, and every string and stream is encrypted with its per-object key.
// It exercises the whole decryption path end to end — /Encrypt parsing, key
// derivation, per-object keys, string and stream decryption, and the exceptions.
func buildRC4Document(t *testing.T, ownerPw, userPw []byte, perm int32) []byte {
	t.Helper()
	const keyLen = 5
	const rev = 2
	id := []byte("0123456789abcdef")

	owner := rc4Apply(computeRC4Key(ownerPw, rev, keyLen), truncateOrPad(userPw))
	fileKey := computeKeyRev234(userPw, owner, perm, id, true, keyLen, rev)
	user := computeUserEntry(userPw, owner, perm, id, rev, keyLen, true)

	h := &securityHandler{key: fileKey}
	enc := func(num int64, data []byte) []byte {
		return h.decryptBytes(data, num, 0) // RC4 is its own inverse
	}

	secret := []byte("classified")
	streamData := []byte("stream payload")
	objs := []rdrObj{
		{num: 1, body: "<< /Type /Catalog /Pages 2 0 R >>"},
		{num: 2, body: "<< /Type /Pages /Kids [3 0 R] /Count 1 >>"},
		{num: 3, body: "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>"},
		{num: 4, body: "<< /Secret <" + hex.EncodeToString(enc(4, secret)) + "> >>"},
		{num: 5, body: fmt.Sprintf("<< /Filter /Standard /V 1 /R 2 /Length 40 /P %d /O <%s> /U <%s> >>",
			perm, hex.EncodeToString(owner), hex.EncodeToString(user))},
		{num: 6, body: fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream",
			len(streamData), enc(6, streamData))},
		// A signature dictionary: its /Contents is never encrypted, so it is
		// stored in the clear and must come back unchanged.
		{num: 7, body: "<< /Type /Sig /ByteRange [0 1 2 3] /Contents <DEADBEEF> " +
			"/Name <" + hex.EncodeToString(enc(7, []byte("signer"))) + "> >>"},
		// PDFBOX-4466 / DSS-1538: no /Type, but /Contents + /ByteRange are enough.
		{num: 8, body: "<< /ByteRange [0 1 2 3] /Contents <C0FFEE> >>"},
	}
	return buildReaderPDF("%PDF-1.4\n", objs,
		fmt.Sprintf("/Encrypt 5 0 R\n/ID [<%s> <%s>]\n",
			hex.EncodeToString(id), hex.EncodeToString(id)))
}

func TestEncryptionRC4UserPassword(t *testing.T) {
	data := buildRC4Document(t, []byte("owner"), []byte(""), -44)
	d, err := OpenBytes(data, nil)
	if err != nil {
		t.Fatalf("open with the empty user password: %v", err)
	}
	if !d.IsEncrypted() {
		t.Fatal("IsEncrypted = false")
	}
	enc := d.Encryption()
	if enc.Handler != "Standard" || enc.V != 1 || enc.R != 2 || enc.KeyLength != 40 {
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
	raw, _ := d.RawStreamData(sobj.(*Stream))
	if string(raw) != "stream payload" {
		t.Errorf("decrypted stream = %q", raw)
	}
	perms := d.Permissions()
	if perms.OwnerAccess {
		t.Error("the user password must not grant owner access")
	}
	if perms.Raw != -44 {
		t.Errorf("permission bits = %d", perms.Raw)
	}
}

// TestEncryptionSignatureContentsIsNeverDecrypted is the rule whose absence
// silently corrupts every signature in an encrypted document.
func TestEncryptionSignatureContentsIsNeverDecrypted(t *testing.T) {
	data := buildRC4Document(t, []byte("owner"), []byte(""), -1)
	d, err := OpenBytes(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		num  int64
		want []byte
	}{
		{7, []byte{0xDE, 0xAD, 0xBE, 0xEF}},
		{8, []byte{0xC0, 0xFF, 0xEE}}, // /Type is optional (PDFBOX-4466)
	} {
		obj, err := d.Object(ObjectKey{Num: tc.num})
		if err != nil {
			t.Fatal(err)
		}
		got, _ := d.GetString(obj.(*Dict), "Contents")
		if !bytes.Equal(got, tc.want) {
			t.Errorf("object %d /Contents = %x, want %x", tc.num, got, tc.want)
		}
	}
	// …while the other strings of the same dictionary *are* decrypted.
	obj, _ := d.Object(ObjectKey{Num: 7})
	if got, _ := d.GetString(obj.(*Dict), "Name"); string(got) != "signer" {
		t.Errorf("/Name = %q, want the decrypted value", got)
	}
}

func TestEncryptionOwnerPassword(t *testing.T) {
	data := buildRC4Document(t, []byte("owner"), []byte("user"), -44)
	if _, err := OpenBytes(data, nil); !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("empty password err = %v, want ErrInvalidPassword", err)
	}
	d, err := OpenBytes(data, &Options{Password: []byte("owner")})
	if err != nil {
		t.Fatalf("owner password: %v", err)
	}
	perms := d.Permissions()
	if !perms.OwnerAccess || !perms.CanModify || !perms.CanFillInForm {
		t.Errorf("owner permissions = %+v", perms)
	}
	// The owner password must still decrypt content.
	obj, err := d.Object(ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := d.GetString(obj.(*Dict), "Secret"); string(got) != "classified" {
		t.Errorf("decrypted string = %q", got)
	}

	d, err = OpenBytes(data, &Options{Password: []byte("user")})
	if err != nil {
		t.Fatalf("user password: %v", err)
	}
	if d.Permissions().OwnerAccess {
		t.Error("the user password must not grant owner access")
	}
}

func TestEncryptionWrongPassword(t *testing.T) {
	data := buildRC4Document(t, []byte("owner"), []byte("user"), -44)
	_, err := OpenBytes(data, &Options{Password: []byte("nope")})
	if !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("err = %v, want ErrInvalidPassword", err)
	}
}

func TestEncryptionUnsupportedHandler(t *testing.T) {
	objs := catalogObjs(rdrObj{num: 5,
		body: "<< /Filter /Adobe.PubSec /V 4 /R 4 /Length 128 >>"})
	data := buildReaderPDF("%PDF-1.6\n", objs, "/Encrypt 5 0 R\n")
	_, err := OpenBytes(data, nil)
	if !errors.Is(err, ErrUnsupportedSecurityHandler) {
		t.Errorf("err = %v, want ErrUnsupportedSecurityHandler", err)
	}
}

func TestPermissionBits(t *testing.T) {
	// bit numbers are 1-based: 3 print, 4 modify, 5 extract, 6 annots,
	// 9 fill-in, 11 assemble, 12 print faithful.
	all := int32(-1)
	p := permissionsFrom(all, false)
	if !p.CanPrint || !p.CanModify || !p.CanExtract || !p.CanModifyAnnots ||
		!p.CanFillInForm || !p.CanAssemble || !p.CanPrintFaithful {
		t.Errorf("all bits set = %+v", p)
	}
	none := permissionsFrom(0, false)
	if none.CanPrint || none.CanModify || none.CanFillInForm || none.OwnerAccess {
		t.Errorf("no bits set = %+v", none)
	}
	onlyFill := permissionsFrom(1<<8, false)
	if !onlyFill.CanFillInForm || onlyFill.CanModify {
		t.Errorf("bit 9 only = %+v", onlyFill)
	}
	owner := permissionsFrom(0, true)
	if !owner.OwnerAccess || !owner.CanModify || !owner.CanFillInForm {
		t.Errorf("owner access must grant everything: %+v", owner)
	}
}

func TestUnencryptedDocumentGrantsEverything(t *testing.T) {
	d := mustOpen(t, buildReaderPDF("%PDF-1.4\n", catalogObjs(), ""))
	p := d.Permissions()
	if !p.OwnerAccess || !p.CanModify || !p.CanModifyAnnots || !p.CanFillInForm {
		t.Errorf("permissions = %+v", p)
	}
	if d.Encryption() != nil {
		t.Error("Encryption() must be nil when the document is not encrypted")
	}
}

func TestObjectKeyDerivation(t *testing.T) {
	// Algorithm 1: MD5(key || low3(num) || low2(gen)), truncated to
	// min(len(key)+5, 16). With AES the "sAlT" suffix is added.
	h := &securityHandler{key: bytes.Repeat([]byte{0xAB}, 5)}
	k := h.objectKeyFor(1, 0)
	if len(k) != 10 {
		t.Errorf("40-bit object key length = %d, want 10", len(k))
	}
	h = &securityHandler{key: bytes.Repeat([]byte{0xAB}, 16)}
	if got := len(h.objectKeyFor(1, 0)); got != 16 {
		t.Errorf("128-bit object key length = %d, want 16", got)
	}
	// AES-256 uses the file key directly (Algorithm 1.A).
	h = &securityHandler{key: bytes.Repeat([]byte{0xCD}, 32), useAES: true}
	if got := h.objectKeyFor(7, 3); !bytes.Equal(got, h.key) {
		t.Errorf("AES-256 must use the file key unchanged")
	}
	// The object and generation numbers change the key.
	h = &securityHandler{key: bytes.Repeat([]byte{0xAB}, 16)}
	if bytes.Equal(h.objectKeyFor(1, 0), h.objectKeyFor(2, 0)) {
		t.Error("the object number must feed the key")
	}
	if bytes.Equal(h.objectKeyFor(1, 0), h.objectKeyFor(1, 1)) {
		t.Error("the generation must feed the key")
	}
}

func TestTruncateOrPad(t *testing.T) {
	got := truncateOrPad([]byte("abc"))
	if len(got) != 32 || string(got[:3]) != "abc" {
		t.Fatalf("padded = %x", got)
	}
	if !bytes.Equal(got[3:], encryptPadding[:29]) {
		t.Errorf("padding = %x", got[3:])
	}
	long := bytes.Repeat([]byte{'x'}, 40)
	if got := truncateOrPad(long); len(got) != 32 {
		t.Errorf("truncated = %d bytes", len(got))
	}
}

func TestComputeHash2BIsDeterministic(t *testing.T) {
	// Algorithm 2.B is iterative and expensive; the property that matters here is
	// that it terminates, is deterministic, and yields 32 bytes.
	pw := []byte("password")
	salt := []byte("12345678")
	a := computeHash2B(append(append([]byte{}, pw...), salt...), pw, nil)
	b := computeHash2B(append(append([]byte{}, pw...), salt...), pw, nil)
	if len(a) != 32 {
		t.Fatalf("hash length = %d, want 32", len(a))
	}
	if !bytes.Equal(a, b) {
		t.Error("computeHash2B is not deterministic")
	}
	c := computeHash2B(append(append([]byte{}, pw...), []byte("87654321")...), pw, nil)
	if bytes.Equal(a, c) {
		t.Error("the salt must change the hash")
	}
}

func TestRC4KnownAnswer(t *testing.T) {
	// RFC 6229 / classic test vector: key "Key", plaintext "Plaintext".
	got := rc4Apply([]byte("Key"), []byte("Plaintext"))
	if hex.EncodeToString(got) != "bbf316e8d940af0ad3" {
		t.Errorf("RC4(Key, Plaintext) = %x", got)
	}
}

func TestXRefStreamsAreNeverDecrypted(t *testing.T) {
	// /Type /XRef and /Crypt /Identity streams are exempt from decryption.
	h := &securityHandler{key: bytes.Repeat([]byte{1}, 16), encryptMeta: true}
	h.enc.StmF, h.enc.StrF = "StdCF", "StdCF"
	d := &Document{opts: Options{}.withDefaults(), sec: h}
	payload := []byte("must not change")

	xrefStream := NewStream(DictOf(Name("Type"), Name("XRef")), append([]byte{}, payload...))
	d.decryptStream(xrefStream, ObjectKey{Num: 1}, 0)
	if !bytes.Equal(xrefStream.Raw, payload) {
		t.Error("a cross-reference stream was decrypted")
	}

	identity := NewStream(DictOf(
		Name("Filter"), Name("Crypt"),
		Name("DecodeParms"), DictOf(Name("Name"), Name("Identity")),
	), append([]byte{}, payload...))
	d.decryptStream(identity, ObjectKey{Num: 2}, 0)
	if !bytes.Equal(identity.Raw, payload) {
		t.Error("a /Crypt /Identity stream was decrypted")
	}
}

func TestUnencryptedMetadataIsLeftAlone(t *testing.T) {
	h := &securityHandler{key: bytes.Repeat([]byte{1}, 16), encryptMeta: false}
	h.enc.StmF, h.enc.StrF = "StdCF", "StdCF"
	d := &Document{opts: Options{}.withDefaults(), sec: h}
	payload := []byte("<?xpacket begin…")
	meta := NewStream(DictOf(Name("Type"), Name("Metadata")), append([]byte{}, payload...))
	d.decryptStream(meta, ObjectKey{Num: 1}, 0)
	if !bytes.Equal(meta.Raw, payload) {
		t.Error("/EncryptMetadata false was not honoured")
	}
}

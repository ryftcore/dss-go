// The incremental-update serializer: object allocation, the write set, xref
// continuation and the /ByteRange patch. See DESIGN.md §3.1, §3.2 (R15, R16,
// R20, R21), §3.4 and §4.6.
//
// Provenance: COSWriter's incremental path (prepareIncrement, doWriteBody,
// doWriteXRefInc, doWriteIncrement, doWriteSignature) of pdfbox 3.0.7, driven
// the way PdfBoxSignatureService drives PDDocument.saveIncremental.
//
// # The one invariant (§3.1)
//
// The output is the input's bytes, unchanged, followed by an appended
// increment. Not "logically equivalent" — byte-identical prefix. Nothing here
// ever seeks backwards into the original: every offset is computed in the
// concatenation, and the original is copied verbatim exactly once, at the end.

package pdf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sort"
)

// errCatalogNotIndirect is returned when /Root is a direct dictionary, which
// leaves an incremental update no slot to rewrite.
var errCatalogNotIndirect = errors.New("pdf: /Root is not an indirect reference; cannot update incrementally")

// writerEncryptHook encrypts one string or stream payload of the increment with
// the source document's existing handler and key (R20). It is a variable so
// that a test can pin a deterministic transform, and so that crypt.go may take
// the implementation over wholesale later.
//
// DESIGN GAP, flagged for the integrator/auditor: DESIGN.md §4 pins no
// encryption-on-write entry point — §4.3 and §4.4 expose decryption only — yet
// R20 requires the increment to be re-encrypted for an encrypted source.
// encryptForWrite below is the writer-side minimum that closes the gap: it
// derives nothing of its own, reusing crypt.go's securityHandler.objectKeyFor
// (Algorithm 1 / 1.A) and rc4Apply, so the handler is never upgraded,
// downgraded or re-keyed and /Encrypt is carried by reference. If crypt.go
// grows a real encryptBytes, delete encryptForWrite and point the hook at it.
//
// The exemptions of §2.6 (the signature /Contents, /Encrypt itself, the trailer
// /ID, xref streams, /Crypt /Identity) are applied on this side, in writer.go
// and in installEncryption, so the hook only ever sees payloads that must be
// encrypted.
var writerEncryptHook = func(d *Document, k ObjectKey, data []byte, isString bool) ([]byte, error) {
	return d.encryptForWrite(k, data, isString)
}

// ErrEncryptedWriteUnsupported is returned by NewUpdater for an encrypted source
// document that has no usable security handler.
var ErrEncryptedWriteUnsupported = errors.New("pdf: writing an increment into an encrypted document requires the crypt.go encryptor")

// encryptForWrite is the inverse of securityHandler.decryptBytes: RC4 or
// AES-CBC with a fresh initialisation vector, PKCS#5 padded. The IV comes from
// Options.Random, so the pades layer can inject DSS's deterministic
// SecureRandomProvider seed and get reproducible output (R20, R21).
func (d *Document) encryptForWrite(k ObjectKey, data []byte, isString bool) ([]byte, error) {
	h := d.sec
	if h == nil || len(h.key) == 0 {
		return data, nil
	}
	if isString && h.enc.StrF == "Identity" {
		return data, nil
	}
	if !isString && h.enc.StmF == "Identity" {
		return data, nil
	}
	key := h.objectKeyFor(k.Num, k.Gen)
	if !h.useAES {
		return rc4Apply(key, data), nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	src := d.opts.Random
	if src == nil {
		src = rand.Reader
	}
	if _, err := io.ReadFull(src, iv); err != nil {
		return nil, err
	}
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	body := make([]byte, 0, len(data)+padLen)
	body = append(body, data...)
	for i := 0; i < padLen; i++ {
		body = append(body, byte(padLen))
	}
	out := make([]byte, aes.BlockSize+len(body))
	copy(out, iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out[aes.BlockSize:], body)
	return out, nil
}

// Updater accumulates an incremental update over a Document. The Document's
// bytes are never modified; Write emits original||increment.
type Updater struct {
	doc     *Document
	orig    []byte
	origLen int64

	// nextNum is R16's allocator: numbers start at HighestObjectNumber()+1 and
	// increase by one per Alloc/Add call, in call order. Generation is always
	// 0 and free-list reuse is deliberately not implemented — a reused object
	// number in an incremental update is a well-known interoperability hazard.
	nextNum int64

	sched   map[ObjectKey]Object
	updated map[ObjectKey]Object // Update()'s clone cache

	catalogKey ObjectKey
	catalog    *Dict

	sigAdded bool
	sig      *signaturePlan

	documentID []byte
}

// signaturePlan is what AddSignature leaves behind for Write to patch.
type signaturePlan struct {
	dict     *Dict
	key      ObjectKey
	fieldKey ObjectKey
}

// NewUpdater starts an incremental update over d.
func NewUpdater(d *Document) (*Updater, error) {
	if d == nil {
		return nil, errors.New("pdf: NewUpdater(nil)")
	}
	if d.IsEncrypted() && (writerEncryptHook == nil || d.sec == nil) {
		return nil, ErrEncryptedWriteUnsupported
	}
	orig, err := d.Bytes()
	if err != nil {
		return nil, err
	}
	u := &Updater{
		doc:     d,
		orig:    orig,
		origLen: int64(len(orig)),
		nextNum: d.HighestObjectNumber() + 1,
		sched:   make(map[ObjectKey]Object),
		updated: make(map[ObjectKey]Object),
	}
	r, ok := d.Trailer().GetRaw("Root").(Ref)
	if !ok {
		// Every real document has an indirect /Root; a direct one would leave
		// catalog edits with no slot to be written into, which would fail
		// silently rather than loudly.
		return nil, errCatalogNotIndirect
	}
	u.catalogKey = r.Key()
	return u, nil
}

// Alloc reserves the next object number (R16) without writing anything.
func (u *Updater) Alloc() ObjectKey {
	k := ObjectKey{Num: u.nextNum, Gen: 0}
	u.nextNum++
	return k
}

// Add allocates a key and schedules obj to be written under it.
func (u *Updater) Add(obj Object) ObjectKey {
	k := u.Alloc()
	u.sched[k] = obj
	return k
}

// Put schedules obj to be written under an existing key, replacing that object
// in this revision. It never touches the original bytes.
func (u *Updater) Put(k ObjectKey, obj Object) { u.sched[k] = obj }

// Catalog returns a mutable clone of the catalog, already scheduled for
// writing. Repeated calls return the same instance.
func (u *Updater) Catalog() *Dict {
	if u.catalog != nil {
		return u.catalog
	}
	cat, err := u.doc.Catalog()
	if err != nil || cat == nil {
		// A document whose catalog cannot be read never reaches here: Open
		// fails with ErrBrokenCatalog first (§2.7). Degrade to an empty
		// dictionary rather than panicking in a writer.
		cat = NewDict()
	}
	u.catalog = cat.Clone()
	if !u.catalogKey.IsZero() {
		u.sched[u.catalogKey] = u.catalog
		u.updated[u.catalogKey] = u.catalog
	}
	return u.catalog
}

// Update returns a mutable clone of an existing object, already scheduled.
// Repeated calls for one key return the same instance, so two features editing
// the same dictionary compose instead of clobbering each other.
func (u *Updater) Update(k ObjectKey) (Object, error) {
	if o, ok := u.updated[k]; ok {
		return o, nil
	}
	if k == u.catalogKey {
		return u.Catalog(), nil
	}
	o, err := u.doc.Object(k)
	if err != nil {
		return nil, err
	}
	var c Object
	switch v := o.(type) {
	case *Dict:
		c = v.Clone()
	case *Stream:
		if u.doc.IsEncrypted() {
			// The clone would carry bytes that are still encrypted under the
			// document key; re-writing them would encrypt them twice. DSS never
			// rewrites a stream in place, so refuse rather than corrupt.
			return nil, errors.New("pdf: rewriting a stream of an encrypted document is not supported")
		}
		c = &Stream{Dict: v.Dict.Clone(), Raw: append([]byte(nil), v.Raw...)}
	case Array:
		c = append(Array(nil), v...)
	default:
		c = o
	}
	u.updated[k] = c
	u.sched[k] = c
	return c, nil
}

// Scheduled reports the keys that will be written, ascending (R15).
func (u *Updater) Scheduled() []ObjectKey {
	keys := make([]ObjectKey, 0, len(u.sched))
	for k := range u.sched {
		keys = append(keys, k)
	}
	sortKeys(keys)
	return keys
}

func sortKeys(keys []ObjectKey) {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Num != keys[j].Num {
			return keys[i].Num < keys[j].Num
		}
		return keys[i].Gen < keys[j].Gen
	})
}

// Result is the laid-out increment, before the CMS is inserted.
type Result struct {
	Bytes          []byte // original || increment
	OriginalLength int64
	// ByteRange is zero-valued when the update carries no signature.
	ByteRange           [4]int64
	ContentsOffset      int64 // absolute offset of the '<'
	ContentsLength      int64 // reserved bytes including '<' and '>'
	StartXref           int64
	XRefStyle           XRefStyle
	HighestObjectNumber int64
}

// Write lays out the increment and patches /ByteRange (R18). /Contents still
// holds the placeholder.
//
// Object write order is ascending object number. This is a deliberate deviation
// from pdfbox (R15): COSWriter drains an ArrayDeque seeded from a HashSet, so
// its order is not reproducible across runs, and reproducing it would mean
// reproducing a bug. Offsets in the xref come from the actual write positions,
// so nothing about validity changes. We also never emit an object stream on
// write — every new object is a plain "N G obj".
func (u *Updater) Write() (*Result, error) {
	var buf bytes.Buffer
	w := newWriterAt(&buf, u.origLen)
	if u.sig != nil {
		w.watch = &sigWatch{dict: u.sig.dict}
	}

	keys := u.Scheduled()
	entries := make([]xrefEntry, 0, len(keys)+2)
	highest := u.doc.HighestObjectNumber()

	for _, k := range keys {
		if k.Num > highest {
			highest = k.Num
		}
		entries = append(entries, xrefEntry{Key: k, Offset: w.Pos()})
		u.installEncryption(w, k)
		if err := w.WriteIndirect(k, u.sched[k]); err != nil {
			return nil, err
		}
	}
	w.encrypt, w.skipEncrypt = nil, nil

	style := u.xrefStyle()
	var startxref int64
	switch style {
	case XRefStream:
		// The xref stream object is allocated last, exactly as pdfbox does, so
		// that /Size is highest+2 (§3.4).
		xrefKey := u.Alloc()
		startxref = w.Pos()
		entries = append(entries, freeHeadEntry())
		// Unlike pdfbox we list the xref stream object in its own table: the
		// offset is known before it is written, and a self-describing section
		// is what every other producer emits.
		entries = append(entries, xrefEntry{Key: xrefKey, Offset: startxref})
		stm, err := buildXRefStream(entries, u.buildTrailer(highest, false), xrefKey.Num+1, u.doc.StartXref())
		if err != nil {
			return nil, err
		}
		if err := w.WriteIndirect(xrefKey, stm); err != nil {
			return nil, err
		}
		if xrefKey.Num > highest {
			highest = xrefKey.Num
		}
	default:
		style = XRefTable
		entries = append(entries, freeHeadEntry())
		startxref = w.Pos()
		if err := writeXRefTable(w, entries); err != nil {
			return nil, err
		}
		if err := writeTrailer(w, u.buildTrailer(highest, true)); err != nil {
			return nil, err
		}
	}
	if err := writeTail(w, startxref); err != nil {
		return nil, err
	}
	if err := w.Err(); err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(u.orig)+buf.Len())
	out = append(out, u.orig...)
	out = append(out, buf.Bytes()...)

	res := &Result{
		Bytes:               out,
		OriginalLength:      u.origLen,
		StartXref:           startxref,
		XRefStyle:           style,
		HighestObjectNumber: highest,
	}
	if u.sig != nil {
		if err := patchByteRange(res, w.watch); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// installEncryption arms the writer's encryption hooks for object k (R20). The
// handler, the key and every §2.6 exemption stay exactly as the source document
// has them: nothing is upgraded, downgraded or re-keyed, and /Encrypt is
// carried by reference.
func (u *Updater) installEncryption(w *Writer, k ObjectKey) {
	if !u.doc.IsEncrypted() || writerEncryptHook == nil {
		w.encrypt, w.skipEncrypt = nil, nil
		return
	}
	doc, key := u.doc, k
	w.encrypt = func(data []byte, isString bool) ([]byte, error) {
		return writerEncryptHook(doc, key, data, isString)
	}
	// The /Contents of a /Sig or /DocTimeStamp dictionary is never encrypted —
	// pdfbox guards this explicitly and getting it wrong silently corrupts
	// every signature in an encrypted document (§2.6). writer.go additionally
	// exempts the watched signature dictionary; this covers a signature
	// dictionary reached any other way.
	w.skipEncrypt = nil
	if d, ok := u.sched[k].(*Dict); ok {
		if t, ok := d.GetRaw("Type").(Name); ok && (t == "Sig" || t == "DocTimeStamp") {
			w.skipEncrypt = map[Name]bool{"Contents": true}
		}
	}
}

// xrefStyle implements §3.4's style selection. A hybrid-reference file always
// degrades to a table, matching COSWriter.doWriteXRefInc.
func (u *Updater) xrefStyle() XRefStyle {
	if u.doc.HasHybridXRef() {
		return XRefTable
	}
	sections := u.doc.XRefSections()
	if len(sections) == 0 {
		return XRefTable
	}
	if sections[0].Style == XRefStream {
		return XRefStream
	}
	return XRefTable
}

// buildTrailer implements R13. table selects the keys that only make sense in
// the table form: the stream form carries /Prev and /Size itself and must not
// repeat /XRefStm.
func (u *Updater) buildTrailer(highest int64, table bool) *Dict {
	t := u.doc.Trailer().Clone()
	if t == nil {
		t = NewDict()
	}
	if table {
		t.Set("Prev", Integer(u.doc.StartXref()))
		t.Set("Size", Integer(highest+1))
	}
	// COSWriter.doWriteTrailer drops /DocChecksum; 35 corpus documents carry
	// one and it would be stale the moment we append.
	t.Delete("DocChecksum")
	// A hybrid file degraded to a table must not keep pointing at the old
	// /XRefStm section.
	t.Delete("XRefStm")
	if id, ok := u.trailerID(); ok {
		t.Set("ID", id) // forced direct (R13)
	}
	return t
}

// trailerID implements R21's /ID rule: the first element is carried from the
// source, the second is the caller-supplied DocumentID when there is one, else
// the first element repeated. Both are written as hex strings, which is what
// every producer emits and what keeps the increment printable.
func (u *Updater) trailerID() (Array, bool) {
	id := u.doc.ID()
	first, second := id[0], id[1]
	if u.documentID != nil {
		second = u.documentID
		if first == nil {
			first = u.documentID
		}
	} else {
		second = first
	}
	if first == nil && second == nil {
		return nil, false
	}
	if first == nil {
		first = second
	}
	if second == nil {
		second = first
	}
	return Array{
		String{Bytes: first, Hex: true},
		String{Bytes: second, Hex: true},
	}, true
}

// patchByteRange formats the real /ByteRange over the reserved 35 bytes (R18)
// and pads the remainder with spaces. The corpus confirms the shape exactly:
// DSS's own validation/doc-firmado-LT.pdf contains "/ByteRange [0 1383 28009
// 307]" followed by eighteen spaces.
func patchByteRange(res *Result, watch *sigWatch) error {
	if watch == nil || !watch.contentsSeen || !watch.byteRangeSeen {
		return errors.New("pdf: signature placeholders were not laid out")
	}
	total := int64(len(res.Bytes))
	before := watch.contentsOffset
	afterOffset := watch.contentsOffset + watch.contentsLength
	afterLength := total - afterOffset

	res.ByteRange = [4]int64{0, before, afterOffset, afterLength}
	res.ContentsOffset = watch.contentsOffset
	res.ContentsLength = watch.contentsLength

	s := fmt.Sprintf("0 %d %d %d]", before, afterOffset, afterLength)
	if int64(len(s)) > watch.byteRangeLength {
		return fmt.Errorf("%w: %q needs %d bytes, %d reserved",
			ErrByteRangeTooLarge, s, len(s), watch.byteRangeLength)
	}
	for i := int64(0); i < watch.byteRangeLength; i++ {
		if i < int64(len(s)) {
			res.Bytes[watch.byteRangeOffset+i] = s[i]
		} else {
			res.Bytes[watch.byteRangeOffset+i] = 0x20
		}
	}
	return nil
}

// SignedData is the byte stream the CMS must be computed over: exactly the two
// spans named by ByteRange. When the update carries no signature the whole
// output is returned, since there is no excluded span.
func (r *Result) SignedData() io.Reader {
	if r.ByteRange == [4]int64{} {
		return bytes.NewReader(r.Bytes)
	}
	return io.MultiReader(
		bytes.NewReader(r.Bytes[r.ByteRange[0]:r.ByteRange[0]+r.ByteRange[1]]),
		bytes.NewReader(r.Bytes[r.ByteRange[2]:r.ByteRange[2]+r.ByteRange[3]]),
	)
}

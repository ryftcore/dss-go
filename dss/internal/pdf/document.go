// The reader API. See DESIGN.md §4.3, §2.7 (header rules H1–H3) and §2.9.
//
// Provenance: org.apache.pdfbox.Loader / COSDocument / COSParser.parseHeader and
// org.apache.pdfbox.pdmodel.{PDDocument,PDDocumentCatalog,PDPageTree,PDPage} of
// pdfbox 3.0.7, plus the call sites in
// eu.europa.esig.dss.pdf.pdfbox.PdfBoxDocumentReader that define our scope.

package pdf

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

// Options configures Open. The zero value is valid: every limit falls back to
// its default.
type Options struct {
	// Password is the document's password as UTF-8 text - the Go shape of the
	// Java String pdfbox receives - tried as the owner password, then as the
	// user password. It is hashed the way pdfbox hashes it: encoded as
	// ISO-8859-1 for /R 2-4 (a code point above U+00FF becomes '?'), and as
	// UTF-8 after SASLprep for /R 6 (crypt.go's passwordBytes), so the
	// document opens with the same text it opens with in pdfbox. nil and
	// empty are the same password, the empty one pdfbox tries by default.
	Password []byte
	// Random supplies AES initialisation vectors on write. nil means crypto/rand.
	Random io.Reader

	MaxObjects    int   // default 500000
	MaxDepth      int   // default 512
	MaxStreamSize int64 // default 512 << 20
}

func (o Options) withDefaults() Options {
	if o.MaxObjects <= 0 {
		o.MaxObjects = 500000
	}
	if o.MaxDepth <= 0 {
		o.MaxDepth = 512
	}
	if o.MaxStreamSize <= 0 {
		o.MaxStreamSize = 512 << 20
	}
	return o
}

// Document is a parsed PDF file.
type Document struct {
	data []byte
	opts Options

	warnings []Warning

	headerVersion float32
	headerOffset  int64

	xref          map[ObjectKey]xrefRec
	byNum         map[int64]ObjectKey
	sections      []XRefSection
	trailer       *Dict
	startxref     int64
	hybrid        bool
	highestObjNum int64

	cache   map[ObjectKey]Object
	loading map[ObjectKey]bool
	objStms map[int64]*objStm
	brute   bruteForce

	encryptRef ObjectKey
	sec        *securityHandler

	revisions      []Revision
	catalog        *Dict
	trailerRebuilt bool
	pages          []pageRef
	sigFields      []SignatureField
	sigDicts       []SignatureDictionary
	sigsDone       bool
}

type pageRef struct {
	dict *Dict
	key  ObjectKey
}

// Open parses the document behind r. size is the length of the source.
func Open(r io.ReaderAt, size int64, opts *Options) (*Document, error) {
	if size < 0 {
		return nil, ErrNotPDF
	}
	buf := make([]byte, size)
	if size > 0 {
		if _, err := io.ReadFull(io.NewSectionReader(r, 0, size), buf); err != nil {
			return nil, err
		}
	}
	return OpenBytes(buf, opts)
}

// OpenBytes parses b. The slice is retained, not copied: the writer needs the
// original bytes verbatim (R: §3.1) and callers must not mutate it.
func OpenBytes(b []byte, opts *Options) (*Document, error) {
	var o Options
	if opts != nil {
		o = *opts
	}
	d := &Document{
		data:      b,
		opts:      o.withDefaults(),
		cache:     make(map[ObjectKey]Object),
		loading:   make(map[ObjectKey]bool),
		objStms:   make(map[int64]*objStm),
		startxref: -1,
	}
	if err := d.parseHeader(); err != nil {
		return nil, err
	}
	if err := d.buildXRef(); err != nil {
		return nil, err
	}
	if len(d.xref) > d.opts.MaxObjects {
		return nil, fmt.Errorf("%w: %d objects exceeds MaxObjects", ErrLimitExceeded, len(d.xref))
	}
	if err := d.setupEncryption(); err != nil {
		return nil, err
	}
	if d.sec != nil {
		// Objects loaded while bootstrapping the xref and the security handler
		// were cached undecrypted; drop them.
		d.cache = make(map[ObjectKey]Object)
		d.objStms = make(map[int64]*objStm)
	}
	if err := d.checkCatalog(); err != nil {
		return nil, err
	}
	return d, nil
}

// Close releases nothing: Document holds no OS resources. It exists so callers
// can treat it like upstream's PDDocument.
func (d *Document) Close() error { return nil }

// Size is the length of the source document in bytes.
func (d *Document) Size() int64 { return int64(len(d.data)) }

// Bytes returns the whole source. The slice is the document's own; callers must
// not modify it.
func (d *Document) Bytes() ([]byte, error) { return d.data, nil }

func (d *Document) addWarning(code WarningCode, off int64, msg string) {
	addWarning(&d.warnings, code, off, msg)
}

// Warnings returns every recovered defect, in the order they were noticed.
func (d *Document) Warnings() []Warning {
	out := make([]Warning, len(d.warnings))
	copy(out, d.warnings)
	return out
}

func (d *Document) newParser() *parser {
	p := newParser(d.data, &d.warnings, d.opts.MaxDepth, d.opts.MaxStreamSize)
	p.resolveLength = d.lengthResolver()
	return p
}

// lengthResolver resolves an indirect /Length once the xref is usable.
func (d *Document) lengthResolver() func(Ref) (int64, bool) {
	return func(r Ref) (int64, bool) {
		if d.xref == nil {
			return 0, false
		}
		obj, ok := d.object(r.Key())
		if !ok {
			return 0, false
		}
		if v, ok := obj.(Integer); ok {
			return int64(v), true
		}
		return 0, false
	}
}

// --- header ---------------------------------------------------------------

// headerSearchLimit caps the header scan. pdfbox reads lines until one contains
// the marker or one starts with a digit; the cap is ours, so a file made of one
// enormous binary "line" cannot cost more than a page of scanning.
const headerSearchLimit = 4096

// parseHeader implements COSParser.parseHeader (H1–H3).
//
// A file with no %PDF- marker at all is *not* an error: PDFParser only logs
// "Error: Header doesn't contain versioninfo" in lenient mode and leaves
// COSDocument's default version of 1.4 in place. validation/malformed-pades.pdf
// is exactly that file and upstream parses it, so ErrNotPDF is reserved for a
// document from which not even the brute-force scan can recover an object.
func (d *Document) parseHeader() error {
	d.headerVersion = 1.4 // COSDocument's default

	limit := headerSearchLimit
	if len(d.data) < limit {
		limit = len(d.data)
	}
	region := d.data[:limit]

	header := ""
	found := false
	pos := 0
	for lineNo := 0; pos < len(region); lineNo++ {
		lineEnd := pos
		for lineEnd < len(region) && region[lineEnd] != '\r' && region[lineEnd] != '\n' {
			lineEnd++
		}
		line := string(region[pos:lineEnd])
		if idx := strings.Index(line, "%PDF-"); idx >= 0 {
			header = line[idx:]
			d.headerOffset = int64(pos + idx)
			if idx > 0 || lineNo > 0 {
				// H1: garbage before the marker is trimmed. All offsets stay
				// relative to byte 0 of the file, which is what makes the
				// /ByteRange arithmetic work.
				d.addWarning(WarnHeaderGarbage, int64(pos), "garbage before the %PDF- header")
			}
			found = true
			break
		}
		if lineNo > 0 && len(line) > 0 && isDigit(line[0]) {
			// a line starting with a digit has to be the first one with data in it
			break
		}
		pos = lineEnd
		for pos < len(region) && (region[pos] == '\r' || region[pos] == '\n') {
			pos++
			if pos < len(region) && region[pos-1] == '\r' && region[pos] == '\n' {
				continue
			}
			break
		}
	}
	if !found {
		d.addWarning(WarnHeaderVersion, 0, "header doesn't contain version info; defaults to 1.4")
		return nil
	}

	rest := header[len("%PDF-"):]
	if len(rest) < 3 {
		// H2: fewer than three bytes after the marker.
		d.addWarning(WarnHeaderVersion, d.headerOffset, "short %PDF- header; version defaults to 1.4")
		return nil
	}
	if !isVersionTriple(rest) {
		// H3: garbage after %PDF-X.Y on the same line is discarded.
		d.addWarning(WarnHeaderGarbage, d.headerOffset, "garbage after the %PDF- header version")
	}
	v, err := strconv.ParseFloat(rest[:3], 32)
	if err != nil {
		// An unparseable version is 1.7 in lenient mode.
		d.headerVersion = 1.7
		d.addWarning(WarnHeaderVersion, d.headerOffset, "unparseable header version; defaults to 1.7")
		return nil
	}
	d.headerVersion = float32(v)
	return nil
}

// isVersionTriple reports whether s starts with exactly `\d.\d` and nothing else,
// i.e. the header matched pdfbox's "%PDF-\d.\d" regex.
func isVersionTriple(s string) bool {
	return len(s) == 3 && isDigit(s[0]) && s[1] == '.' && isDigit(s[2])
}

// HeaderVersion is the version from %PDF-x.y.
func (d *Document) HeaderVersion() float32 { return d.headerVersion }

// Version is the catalog's /Version when present, else the header version.
func (d *Document) Version() float32 {
	cat, err := d.Catalog()
	if err != nil || cat == nil {
		return d.headerVersion
	}
	if n, ok := d.GetName(cat, "Version"); ok {
		if v, err := strconv.ParseFloat(string(n), 32); err == nil {
			// pdfbox: the catalog version only wins when it is the higher one.
			if float32(v) > d.headerVersion {
				return float32(v)
			}
		}
	}
	return d.headerVersion
}

// Trailer returns the resolved trailer dictionary.
func (d *Document) Trailer() *Dict { return d.trailer }

// StartXref is the offset read from the last startxref, or -1.
func (d *Document) StartXref() int64 { return d.startxref }

// HasHybridXRef reports whether any section carried /XRefStm.
func (d *Document) HasHybridXRef() bool { return d.hybrid }

// XRefSections returns the /Prev chain, newest first.
func (d *Document) XRefSections() []XRefSection {
	out := make([]XRefSection, len(d.sections))
	copy(out, d.sections)
	return out
}

// checkCatalog is the one defect we do not recover from: a /Root that is missing
// or whose /Pages is not a dictionary (pdfbox: "Page tree root must be a
// dictionary").
func (d *Document) checkCatalog() error {
	cat, _ := d.Resolve(d.trailer.GetRaw("Root")).(*Dict)
	if cat == nil {
		return ErrBrokenCatalog
	}
	d.catalog = cat
	// COSParser.checkPages: /Pages is dereferenced once here — which is why a
	// document whose page tree root lives outside the xref grows its xref table
	// by that one brute-force entry at load time — and a non-dictionary is the
	// one defect that fails the load.
	pages, ok := d.Resolve(cat.GetRaw("Pages")).(*Dict)
	if !ok {
		return ErrBrokenCatalog
	}
	if d.trailerRebuilt {
		// checkPagesDictionary runs only when the trailer was rebuilt: it walks
		// /Kids, dropping entries that do not dereference to a dictionary.
		d.checkPagesDictionary(pages, make(map[*Dict]bool), 0)
	}
	return nil
}

// checkPagesDictionary is COSParser.checkPagesDictionary: it removes /Kids
// entries that resolve to null and recurses into nested /Pages nodes.
func (d *Document) checkPagesDictionary(node *Dict, seen map[*Dict]bool, depth int) {
	if node == nil || depth > d.opts.MaxDepth || seen[node] {
		return
	}
	seen[node] = true
	kids, ok := d.GetArray(node, "Kids")
	if !ok {
		return
	}
	kept := make(Array, 0, len(kids))
	for _, kid := range kids {
		resolved := d.Resolve(kid)
		// Resolve maps a missing object and a dangling reference to Null{}, so testing for
		// Null covers every "no such kid" case; it never returns a nil interface.
		if _, isNull := resolved.(Null); isNull {
			d.addWarning(WarnKidRemoved, -1, "removed null object from the pages dictionary")
			continue
		}
		kept = append(kept, kid)
		if kd, isDict := resolved.(*Dict); isDict {
			if n, _ := d.GetName(kd, "Type"); n == "Pages" {
				d.checkPagesDictionary(kd, seen, depth+1)
			}
		}
	}
	if len(kept) != len(kids) {
		node.Set("Kids", kept)
	}
}

// Catalog returns the document catalog.
func (d *Document) Catalog() (*Dict, error) {
	if d.catalog != nil {
		return d.catalog, nil
	}
	cat, _ := d.Resolve(d.trailer.GetRaw("Root")).(*Dict)
	if cat == nil {
		return nil, ErrBrokenCatalog
	}
	d.catalog = cat
	return cat, nil
}

// Info returns the document information dictionary, or (nil, nil) when absent.
func (d *Document) Info() (*Dict, error) {
	v := d.Resolve(d.trailer.GetRaw("Info"))
	if dict, ok := v.(*Dict); ok {
		return dict, nil
	}
	return nil, nil
}

// ID returns the two /ID strings; missing elements are nil.
func (d *Document) ID() [2][]byte {
	var out [2][]byte
	arr, ok := d.Resolve(d.trailer.GetRaw("ID")).(Array)
	if !ok {
		return out
	}
	for i := 0; i < 2 && i < len(arr); i++ {
		if s, ok := d.Resolve(arr[i]).(String); ok {
			out[i] = s.Bytes
		}
	}
	return out
}

// --- object access --------------------------------------------------------

// Resolve follows Ref chains. A dangling reference resolves to Null{}, never an
// error (§2.7 O3) — visitFromDictionary upstream skips nil-valued entries, and a
// hard error here would fail documents pdfbox parses.
func (d *Document) Resolve(o Object) Object {
	seen := 0
	for {
		r, ok := o.(Ref)
		if !ok {
			if o == nil {
				return Null{}
			}
			return o
		}
		seen++
		if seen > d.opts.MaxDepth {
			d.addWarning(WarnDanglingReference, -1, "reference chain too deep")
			return Null{}
		}
		obj, ok := d.object(r.Key())
		if !ok {
			d.addWarning(WarnDanglingReference, -1, "dangling reference "+r.Key().String())
			return Null{}
		}
		o = obj
	}
}

// Object returns the object stored under k.
func (d *Document) Object(k ObjectKey) (Object, error) {
	obj, ok := d.object(k)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchObject, k)
	}
	return obj, nil
}

// object is the internal loader: cache, xref lookup, brute-force fallback.
func (d *Document) object(k ObjectKey) (Object, bool) {
	if obj, ok := d.cache[k]; ok {
		return obj, true
	}
	if d.loading[k] {
		return nil, false // cycle while resolving an indirect /Length
	}
	e, ok := d.xref[k]
	if !ok {
		// pdfbox falls back to the brute-force map, which is keyed by number.
		if alt, ok2 := d.byNum[k.Num]; ok2 {
			if e2, ok3 := d.xref[alt]; ok3 {
				e, ok, k = e2, true, alt
			}
		}
	}
	if !ok {
		// COSParser.getObjectOffset: an object the xref does not know about is
		// looked up in the brute-force map and, when found, *written back into
		// the xref table* ("Set missing offset … for object …"). The write-back
		// is not cosmetic: COSDocument.getXrefTable() therefore grows as objects
		// are dereferenced, and its size is what the oracle reports as `objects`.
		if off, found := d.bf().offsets[k]; found {
			e, ok = xrefRec{typ: 1, offset: off, gen: k.Gen}, true
			d.xref[k] = e
			if _, exists := d.byNum[k.Num]; !exists {
				d.byNum[k.Num] = k
			}
			d.noteObjectNumber(k.Num)
			d.addWarning(WarnXRefBruteForce, off, "missing offset for object "+k.String()+" supplied by brute force")
		}
	}
	if !ok {
		return nil, false
	}
	d.loading[k] = true
	defer delete(d.loading, k)

	var obj Object
	switch e.typ {
	case 1:
		o, good := d.loadAt(e.offset, k)
		if !good {
			return nil, false
		}
		obj = d.decryptObject(o, k)
	case 2:
		o, good := d.objectFromStm(k, e)
		if !good {
			return nil, false
		}
		// Objects inside an object stream are never decrypted individually: the
		// container was decrypted as a whole.
		obj = o
	default:
		return nil, false
	}
	d.cache[k] = obj
	return obj, true
}

// loadAt parses the indirect object at off, expecting key k.
func (d *Document) loadAt(off int64, k ObjectKey) (Object, bool) {
	p := d.newParser()
	obj, _, ok := p.parseIndirectAt(off, k)
	if !ok {
		// O1: the header did not match. pdfbox tries the brute-force offset next.
		if bfOff, found := d.bf().offsets[k]; found && bfOff != off {
			p = d.newParser()
			obj, _, ok = p.parseIndirectAt(bfOff, k)
			if ok {
				d.addWarning(WarnObjectHeaderFixed, off,
					"object "+k.String()+" found by brute force at "+strconv.FormatInt(bfOff, 10))
				return obj, true
			}
		}
		return nil, false
	}
	return obj, true
}

// ObjectKeys returns every key in the resolved xref, ascending by Num then Gen.
// Free entries are not included, matching COSDocument.getXrefTable().
func (d *Document) ObjectKeys() []ObjectKey { return d.sortedXRefKeys() }

// HighestObjectNumber is the maximum object number over all xref sections,
// including free entries and /Size - 1 (R16).
func (d *Document) HighestObjectNumber() int64 { return d.highestObjNum }

// --- typed accessors ------------------------------------------------------

// Get resolves dict[key]. It returns nil when the key is absent.
func (d *Document) Get(dict *Dict, key Name) Object {
	if dict == nil {
		return nil
	}
	raw := dict.GetRaw(key)
	if raw == nil {
		return nil
	}
	return d.Resolve(raw)
}

// GetDict returns dict[key] as a dictionary. A stream counts as a dictionary
// only through GetStream; pdfbox's getCOSDictionary behaves the same way.
func (d *Document) GetDict(dict *Dict, key Name) (*Dict, bool) {
	switch v := d.Get(dict, key).(type) {
	case *Dict:
		return v, true
	case *Stream:
		return v.Dict, true
	}
	return nil, false
}

// GetArray returns dict[key] as an array.
func (d *Document) GetArray(dict *Dict, key Name) (Array, bool) {
	v, ok := d.Get(dict, key).(Array)
	return v, ok
}

// GetStream returns dict[key] as a stream.
func (d *Document) GetStream(dict *Dict, key Name) (*Stream, bool) {
	v, ok := d.Get(dict, key).(*Stream)
	return v, ok
}

// GetName returns dict[key] as a name.
func (d *Document) GetName(dict *Dict, key Name) (Name, bool) {
	v, ok := d.Get(dict, key).(Name)
	return v, ok
}

// GetString returns dict[key]'s bytes.
func (d *Document) GetString(dict *Dict, key Name) ([]byte, bool) {
	v, ok := d.Get(dict, key).(String)
	if !ok {
		return nil, false
	}
	return v.Bytes, true
}

// GetInt returns dict[key] as an integer. A real is truncated, as
// COSNumber.intValue does.
func (d *Document) GetInt(dict *Dict, key Name) (int64, bool) {
	switch v := d.Get(dict, key).(type) {
	case Integer:
		return int64(v), true
	case Real:
		return int64(v.Val), true
	}
	return 0, false
}

// GetReal returns dict[key] as a float.
func (d *Document) GetReal(dict *Dict, key Name) (float64, bool) {
	switch v := d.Get(dict, key).(type) {
	case Integer:
		return float64(v), true
	case Real:
		return v.Val, true
	}
	return 0, false
}

// GetBool returns dict[key] as a boolean.
func (d *Document) GetBool(dict *Dict, key Name) (bool, bool) {
	v, ok := d.Get(dict, key).(Bool)
	return bool(v), ok
}

// GetDate parses a PDF date string "D:YYYYMMDDHHmmSSOHH'mm'" leniently: any
// truncation from the right is accepted, as is a missing "D:" prefix.
func (d *Document) GetDate(dict *Dict, key Name) (time.Time, bool) {
	b, ok := d.GetString(dict, key)
	if !ok {
		return time.Time{}, false
	}
	return parseDate(string(b))
}

// parseDate implements the lenient PDF date parse used by GetDate. It is
// deliberately unexported: §4.3 pins the exported surface and GetDate is the
// only date entry point on it.
func parseDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "D:")
	digits := make([]byte, 0, 14)
	i := 0
	for ; i < len(s) && len(digits) < 14; i++ {
		if isDigit(s[i]) {
			digits = append(digits, s[i])
			continue
		}
		break
	}
	if len(digits) < 4 {
		return time.Time{}, false
	}
	get := func(off, n, def int) int {
		if off+n > len(digits) {
			return def
		}
		v, err := strconv.Atoi(string(digits[off : off+n]))
		if err != nil {
			return def
		}
		return v
	}
	year := get(0, 4, 0)
	month := get(4, 2, 1)
	day := get(6, 2, 1)
	hour := get(8, 2, 0)
	minute := get(10, 2, 0)
	sec := get(12, 2, 0)
	if month < 1 || month > 12 {
		month = 1
	}
	if day < 1 || day > 31 {
		day = 1
	}
	loc := time.UTC
	rest := s[i:]
	if len(rest) > 0 && (rest[0] == '+' || rest[0] == '-' || rest[0] == 'Z') {
		sign := 1
		if rest[0] == '-' {
			sign = -1
		}
		oh, om := 0, 0
		digs := make([]byte, 0, 4)
		for j := 1; j < len(rest) && len(digs) < 4; j++ {
			if isDigit(rest[j]) {
				digs = append(digs, rest[j])
			}
		}
		if len(digs) >= 2 {
			oh, _ = strconv.Atoi(string(digs[:2]))
		}
		if len(digs) >= 4 {
			om, _ = strconv.Atoi(string(digs[2:4]))
		}
		if oh != 0 || om != 0 {
			loc = time.FixedZone("", sign*(oh*3600+om*60))
		}
	}
	return time.Date(year, time.Month(month), day, hour, minute, sec, 0, loc), true
}

// RefAt returns the key of the indirect reference stored at key, if it is one.
// This is upstream's PdfDict.getObjectKey.
func (d *Document) RefAt(dict *Dict, key Name) (ObjectKey, bool) {
	if dict == nil {
		return ObjectKey{}, false
	}
	r, ok := dict.GetRaw(key).(Ref)
	if !ok {
		return ObjectKey{}, false
	}
	return r.Key(), true
}

// IndexRef is the array-element form: upstream's PdfArray.getObjectKey(i).
func (d *Document) IndexRef(a Array, i int) (ObjectKey, bool) {
	if i < 0 || i >= len(a) {
		return ObjectKey{}, false
	}
	r, ok := a[i].(Ref)
	if !ok {
		return ObjectKey{}, false
	}
	return r.Key(), true
}

// --- streams --------------------------------------------------------------

// StreamData returns the decrypted, fully decoded stream bytes.
func (d *Document) StreamData(s *Stream) ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	raw, err := d.RawStreamData(s)
	if err != nil {
		return nil, err
	}
	names, parms := streamFilters(s.Dict, d.Resolve)
	return Decode(raw, names, parms, &d.warnings)
}

// RawStreamData returns the decrypted but still encoded stream bytes. This is
// what DefaultPdfObjectModificationsFinder.compareDictStreams compares.
func (d *Document) RawStreamData(s *Stream) ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return s.Raw, nil
}

// RawStreamSize is the length of the raw bytes, or -1 when s is nil.
func (d *Document) RawStreamSize(s *Stream) int64 {
	if s == nil {
		return -1
	}
	return int64(len(s.Raw))
}

// --- encryption -----------------------------------------------------------

// IsEncrypted reports whether the document carries an /Encrypt dictionary.
func (d *Document) IsEncrypted() bool { return d.sec != nil }

// Encryption describes the security handler, or nil when the document is not
// encrypted.
func (d *Document) Encryption() *Encryption {
	if d.sec == nil {
		return nil
	}
	e := d.sec.enc
	return &e
}

// Permissions returns the decoded /P bits. An unencrypted document grants
// everything, matching AccessPermission's default.
func (d *Document) Permissions() Permissions {
	if d.sec == nil {
		return permissionsFrom(-1, true)
	}
	return d.sec.perms
}

// --- pages ----------------------------------------------------------------

// buildPages walks the page tree once (PDPageTree.PageIterator).
func (d *Document) buildPages() {
	if d.pages != nil {
		return
	}
	d.pages = []pageRef{}
	cat, err := d.Catalog()
	if err != nil {
		return
	}
	rootObj := cat.GetRaw("Pages")
	root, _ := d.Resolve(rootObj).(*Dict)
	if root == nil {
		return
	}
	rootKey, _ := d.RefAt(cat, "Pages")
	visited := make(map[*Dict]bool)
	d.enqueueKids(root, rootKey, visited, 0)
}

func (d *Document) enqueueKids(node *Dict, key ObjectKey, visited map[*Dict]bool, depth int) {
	if node == nil || depth > d.opts.MaxDepth {
		return
	}
	if isPageTreeNode(node) {
		kids, _ := d.GetArray(node, "Kids")
		for i := range kids {
			kidKey, _ := d.IndexRef(kids, i)
			kid, ok := d.Resolve(kids[i]).(*Dict)
			if !ok {
				// A /Kids entry that resolves to null is removed with a warning.
				d.addWarning(WarnKidRemoved, -1, "page tree kid is not a dictionary")
				continue
			}
			if visited[kid] {
				d.addWarning(WarnKidRemoved, -1, "page tree node already visited")
				continue
			}
			if kid.Has("Kids") {
				visited[kid] = true
			}
			d.enqueueKids(kid, kidKey, visited, depth+1)
		}
		return
	}
	if n, _ := d.GetName(node, "Type"); n == "Page" || !node.Has("Type") {
		d.pages = append(d.pages, pageRef{dict: node, key: key})
		return
	}
	d.addWarning(WarnKidRemoved, -1, "page skipped due to an invalid type")
}

// isPageTreeNode is PDPageTree.isPageTreeNode: some files do not set /Type, so
// the presence of /Kids counts too.
func isPageTreeNode(node *Dict) bool {
	if node == nil {
		return false
	}
	if n, ok := node.GetRaw("Type").(Name); ok && n == "Pages" {
		return true
	}
	return node.Has("Kids")
}

// NumberOfPages is the number of leaves in the page tree.
func (d *Document) NumberOfPages() int {
	d.buildPages()
	return len(d.pages)
}

// Page returns the 1-based page dictionary and its object key.
func (d *Document) Page(page int) (*Dict, ObjectKey, error) {
	d.buildPages()
	if page < 1 || page > len(d.pages) {
		return nil, ObjectKey{}, fmt.Errorf("pdf: page %d out of range (1..%d)", page, len(d.pages))
	}
	p := d.pages[page-1]
	return p.dict, p.key, nil
}

// inherited walks /Parent for an inheritable attribute (/MediaBox, /Resources,
// /Rotate, /CropBox), matching PDPageTree.getInheritableAttribute.
func (d *Document) inherited(node *Dict, key Name) Object {
	visited := make(map[*Dict]bool)
	for i := 0; node != nil && i <= d.opts.MaxDepth; i++ {
		if visited[node] {
			return nil
		}
		visited[node] = true
		if v := d.Get(node, key); v != nil {
			if _, isNull := v.(Null); !isNull {
				return v
			}
		}
		parent, ok := d.GetDict(node, "Parent")
		if !ok {
			parent, ok = d.GetDict(node, "P")
			if !ok {
				return nil
			}
		}
		node = parent
	}
	return nil
}

// PageBox returns the page's /MediaBox, inherited through /Parent. The default
// is US Letter, as PDPage.getMediaBox does.
func (d *Document) PageBox(page int) (Rect, error) {
	dict, _, err := d.Page(page)
	if err != nil {
		return Rect{}, err
	}
	arr, _ := d.inherited(dict, "MediaBox").(Array)
	resolved := make(Array, len(arr))
	for i := range arr {
		resolved[i] = d.Resolve(arr[i])
	}
	if r, ok := RectFromArray(resolved); ok {
		return r, nil
	}
	return Rect{MinX: 0, MinY: 0, MaxX: 612, MaxY: 792}, nil
}

// PageRotation returns /Rotate normalised into {0,90,180,270}.
func (d *Document) PageRotation(page int) int {
	dict, _, err := d.Page(page)
	if err != nil {
		return 0
	}
	var r int64
	switch v := d.inherited(dict, "Rotate").(type) {
	case Integer:
		r = int64(v)
	case Real:
		if math.IsNaN(v.Val) || v.Val < math.MinInt32 || v.Val > math.MaxInt32 {
			return 0
		}
		r = int64(v.Val)
	default:
		return 0
	}
	r = ((r % 360) + 360) % 360
	return int(((r + 45) / 90 % 4) * 90)
}

// Annotation is one /Annots entry.
type Annotation struct {
	Key      ObjectKey
	Dict     *Dict
	Rect     Rect
	Name     string // /T
	Signed   bool   // /V present
	Hidden   bool   // /F bit 2
	NoRotate bool   // /F bit 5
}

// Annotations returns the annotations of a 1-based page.
func (d *Document) Annotations(page int) ([]Annotation, error) {
	dict, _, err := d.Page(page)
	if err != nil {
		return nil, err
	}
	arr, ok := d.GetArray(dict, "Annots")
	if !ok {
		return nil, nil
	}
	out := make([]Annotation, 0, len(arr))
	for i := range arr {
		key, _ := d.IndexRef(arr, i)
		ad, ok := d.Resolve(arr[i]).(*Dict)
		if !ok {
			continue
		}
		a := Annotation{Key: key, Dict: ad}
		if rectArr, ok := d.GetArray(ad, "Rect"); ok {
			resolved := make(Array, len(rectArr))
			for j := range rectArr {
				resolved[j] = d.Resolve(rectArr[j])
			}
			a.Rect, _ = RectFromArray(resolved)
		}
		if t, ok := d.GetString(ad, "T"); ok {
			a.Name = decodeTextString(t)
		}
		if v := ad.GetRaw("V"); v != nil {
			a.Signed = true
		}
		if f, ok := d.GetInt(ad, "F"); ok {
			a.Hidden = f&2 != 0
			a.NoRotate = f&16 != 0
		}
		out = append(out, a)
	}
	return out, nil
}

// decodeTextString renders a PDF text string as Go text, mirroring COSString.getString.
// See DecodeTextString in pdfdocencoding.go, which it now simply delegates to: this used to
// approximate PDFDocEncoding with Latin-1, which mangles every code in 0x18-0x1F and 0x80-0xA0.
func decodeTextString(b []byte) string {
	return DecodeTextString(b)
}

// --- AcroForm and signatures ----------------------------------------------

// AcroForm returns /Root /AcroForm.
func (d *Document) AcroForm() (*Dict, bool) {
	cat, err := d.Catalog()
	if err != nil {
		return nil, false
	}
	return d.GetDict(cat, "AcroForm")
}

// SignatureField is one /FT /Sig field of the AcroForm.
type SignatureField struct {
	Key        ObjectKey
	Dict       *Dict
	Name       string    // fully-qualified /T, dot-joined
	Value      *Dict     // /V, nil for an empty field
	ValueKey   ObjectKey // key of the /V reference; zero when /V is direct or absent
	Widgets    []*Dict
	WidgetKeys []ObjectKey
	Page       int // 0 when the widget is not on any page
	Rect       Rect
	Lock       *Dict // /Lock, nil when absent
}

// SignatureDictionary is the raw /V content. All ETSI semantics live in `pades`.
type SignatureDictionary struct {
	Key       ObjectKey
	Dict      *Dict
	Type      Name // "Sig" | "DocTimeStamp" | ""
	Filter    Name
	SubFilter Name
	Contents  []byte  // decoded string bytes, i.e. the DER CMS
	ByteRange []int64 // verbatim, length not forced to 4
	Fields    []int   // indices into the SignatureFields slice that point here
}

// SignatureFields reproduces PDDocument.getSignatureFields().
func (d *Document) SignatureFields() ([]SignatureField, error) {
	if err := d.collectSignatures(); err != nil {
		return nil, err
	}
	out := make([]SignatureField, len(d.sigFields))
	copy(out, d.sigFields)
	return out, nil
}

// SignatureDictionaries returns one entry per distinct /V signature dictionary,
// deduplicated by object key exactly as PdfBoxDocumentReader does with
// sigDictObject.getKey().getNumber().
func (d *Document) SignatureDictionaries() ([]SignatureDictionary, error) {
	if err := d.collectSignatures(); err != nil {
		return nil, err
	}
	out := make([]SignatureDictionary, len(d.sigDicts))
	copy(out, d.sigDicts)
	return out, nil
}

type fieldNode struct {
	dict   *Dict
	key    ObjectKey
	name   string
	ft     Name
	value  Object
	lock   *Dict
	parent *fieldNode
}

func (d *Document) collectSignatures() error {
	if d.sigsDone {
		return nil
	}
	d.sigsDone = true
	d.sigFields = []SignatureField{}
	d.sigDicts = []SignatureDictionary{}

	form, ok := d.AcroForm()
	if !ok {
		return nil
	}
	fields, ok := d.GetArray(form, "Fields")
	if !ok {
		return nil
	}
	pageOf := d.widgetPageIndex()
	visited := make(map[*Dict]bool)
	for i := range fields {
		key, _ := d.IndexRef(fields, i)
		fd, ok := d.Resolve(fields[i]).(*Dict)
		if !ok {
			continue
		}
		d.walkField(fd, key, nil, visited, pageOf, 0)
	}
	return nil
}

// walkField recurses into /Kids when a node has no /FT of its own, propagating
// the inheritable attributes /FT /Ff /V /DA, and keeps the /Sig leaves.
func (d *Document) walkField(dict *Dict, key ObjectKey, parent *fieldNode, visited map[*Dict]bool, pageOf map[*Dict]int, depth int) {
	if dict == nil || depth > d.opts.MaxDepth || visited[dict] {
		return
	}
	visited[dict] = true

	node := &fieldNode{dict: dict, key: key, parent: parent}
	if ft, ok := d.GetName(dict, "FT"); ok {
		node.ft = ft
	} else if parent != nil {
		node.ft = parent.ft
	}
	if v := dict.GetRaw("V"); v != nil {
		node.value = v
	} else if parent != nil {
		node.value = parent.value
	}
	if lock, ok := d.GetDict(dict, "Lock"); ok {
		node.lock = lock
	}
	node.name = fullyQualifiedName(d, dict, parent)

	kids, hasKids := d.GetArray(dict, "Kids")
	// A node whose kids are widgets (no /T, no /FT of their own) is a merged
	// field; a node whose kids are fields is an intermediate node.
	kidsAreFields := false
	if hasKids {
		for i := range kids {
			kd, ok := d.Resolve(kids[i]).(*Dict)
			if !ok {
				continue
			}
			if kd.Has("T") || kd.Has("FT") {
				kidsAreFields = true
				break
			}
		}
	}
	if hasKids && kidsAreFields {
		for i := range kids {
			kkey, _ := d.IndexRef(kids, i)
			kd, ok := d.Resolve(kids[i]).(*Dict)
			if !ok {
				continue
			}
			d.walkField(kd, kkey, node, visited, pageOf, depth+1)
		}
		return
	}

	if node.ft != "Sig" {
		return
	}
	d.appendSignatureField(node, kids, pageOf)
}

func (d *Document) appendSignatureField(node *fieldNode, kids Array, pageOf map[*Dict]int) {
	f := SignatureField{
		Key:  node.key,
		Dict: node.dict,
		Name: node.name,
		Lock: node.lock,
	}
	// Widgets: the field dict itself when it carries /Subtype /Widget (the merged
	// field+widget case), else its /Kids.
	if sub, ok := d.GetName(node.dict, "Subtype"); ok && sub == "Widget" {
		f.Widgets = []*Dict{node.dict}
		f.WidgetKeys = []ObjectKey{node.key}
	} else {
		for i := range kids {
			kkey, _ := d.IndexRef(kids, i)
			kd, ok := d.Resolve(kids[i]).(*Dict)
			if !ok {
				continue
			}
			f.Widgets = append(f.Widgets, kd)
			f.WidgetKeys = append(f.WidgetKeys, kkey)
		}
	}
	for _, w := range f.Widgets {
		if p, ok := pageOf[w]; ok {
			f.Page = p
			break
		}
	}
	if len(f.Widgets) > 0 {
		if rectArr, ok := d.GetArray(f.Widgets[0], "Rect"); ok {
			resolved := make(Array, len(rectArr))
			for j := range rectArr {
				resolved[j] = d.Resolve(rectArr[j])
			}
			f.Rect, _ = RectFromArray(resolved)
		}
	}
	if r, ok := node.value.(Ref); ok {
		f.ValueKey = r.Key()
	}
	if node.value != nil {
		if vd, ok := d.Resolve(node.value).(*Dict); ok {
			f.Value = vd
		}
	}
	idx := len(d.sigFields)
	d.sigFields = append(d.sigFields, f)
	if f.Value == nil {
		return // an empty signature field
	}
	d.appendSignatureDictionary(f, idx)
}

func (d *Document) appendSignatureDictionary(f SignatureField, fieldIndex int) {
	for i := range d.sigDicts {
		sd := &d.sigDicts[i]
		if (!f.ValueKey.IsZero() && sd.Key == f.ValueKey) || sd.Dict == f.Value {
			sd.Fields = append(sd.Fields, fieldIndex)
			return
		}
	}
	sd := SignatureDictionary{
		Key:    f.ValueKey,
		Dict:   f.Value,
		Fields: []int{fieldIndex},
	}
	sd.Type, _ = d.GetName(f.Value, "Type")
	sd.Filter, _ = d.GetName(f.Value, "Filter")
	sd.SubFilter, _ = d.GetName(f.Value, "SubFilter")
	if c, ok := f.Value.GetRaw("Contents").(String); ok {
		sd.Contents = c.Bytes
	} else if c, ok := d.Get(f.Value, "Contents").(String); ok {
		sd.Contents = c.Bytes
	}
	if br, ok := d.GetArray(f.Value, "ByteRange"); ok {
		sd.ByteRange = make([]int64, 0, len(br))
		for i := range br {
			switch v := d.Resolve(br[i]).(type) {
			case Integer:
				sd.ByteRange = append(sd.ByteRange, int64(v))
			case Real:
				sd.ByteRange = append(sd.ByteRange, int64(v.Val))
			default:
				sd.ByteRange = append(sd.ByteRange, 0)
			}
		}
	}
	d.sigDicts = append(d.sigDicts, sd)
}

// fullyQualifiedName joins the /T values of the ancestor chain with '.', skipping
// nodes without /T (PDField.getFullyQualifiedName).
func fullyQualifiedName(d *Document, dict *Dict, parent *fieldNode) string {
	var own string
	if t, ok := d.GetString(dict, "T"); ok {
		own = decodeTextString(t)
	}
	if parent == nil || parent.name == "" {
		return own
	}
	if own == "" {
		return parent.name
	}
	return parent.name + "." + own
}

// widgetPageIndex maps each annotation dictionary to its 1-based page number.
func (d *Document) widgetPageIndex() map[*Dict]int {
	d.buildPages()
	out := make(map[*Dict]int)
	for i := range d.pages {
		arr, ok := d.GetArray(d.pages[i].dict, "Annots")
		if !ok {
			continue
		}
		for j := range arr {
			if ad, ok := d.Resolve(arr[j]).(*Dict); ok {
				if _, seen := out[ad]; !seen {
					out[ad] = i + 1
				}
			}
		}
	}
	return out
}

// bytesReader avoids importing bytes at every call site.
func bytesReader(b []byte) io.Reader { return bytes.NewReader(b) }

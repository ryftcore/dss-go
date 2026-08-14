// Cross-reference resolution: tables, streams, hybrid files and the /Prev chain.
// See DESIGN.md §2.3 and §2.7 T1–T4.
//
// Provenance: org.apache.pdfbox.pdfparser.COSParser (parseXref, parseXrefTable,
// parseTrailer, checkXrefOffsets, validateXrefOffsets), PDFXrefStreamParser and
// XrefTrailerResolver of pdfbox 3.0.7.

package pdf

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

// XRefStyle distinguishes the two section forms.
type XRefStyle uint8

const (
	XRefTable XRefStyle = iota + 1
	XRefStream
)

func (s XRefStyle) String() string {
	switch s {
	case XRefTable:
		return "table"
	case XRefStream:
		return "stream"
	default:
		return "?"
	}
}

// XRefSection is one hop of the /Prev chain.
type XRefSection struct {
	// Offset is the byte offset of `xref`, or of the `N G obj` of the xref stream.
	Offset    int64
	Style     XRefStyle
	Trailer   *Dict
	Prev      int64 // -1 when absent
	XRefStm   int64 // hybrid: /XRefStm offset, -1 when absent
	Entries   int   // entries contributed by this section
	Recovered bool  // true when the offset had to be repaired (§2.7 X3)
}

// xrefRec is one resolved cross-reference entry. Free entries are not stored:
// pdfbox's parseXrefTable drops them and PDFXrefStreamParser skips type 0, and
// COSDocument.getXrefTable()'s size — which the oracle reports as `objects` — is
// therefore a count of in-use entries only.
type xrefRec struct {
	typ    uint8  // 1 = byte offset, 2 = inside an object stream
	offset int64  // typ 1: byte offset; typ 2: container object number
	gen    uint16 // typ 1: generation
	index  int    // typ 2: index within the container
}

// rawSection is one parsed hop before merging.
type rawSection struct {
	sec     XRefSection
	entries []keyedEntry
	// stmEntries are the hybrid /XRefStm entries, merged beneath the table's.
	stmEntries []keyedEntry
}

type keyedEntry struct {
	key ObjectKey
	e   xrefRec
}

// trailBytes is pdfbox's DEFAULT_TRAIL_BYTECOUNT. We pin it and expose no knob.
const trailBytes = 2048

// maxSections caps the /Prev chain; the corpus maximum is 54 revisions.
const maxSections = 512

// buildXRef performs the startup sequence of DESIGN.md §2.3.
func (d *Document) buildXRef() error {
	d.startxref = -1
	off, ok := d.findStartXref()
	var sections []*rawSection
	if ok {
		d.startxref = off
		sections = d.walkChain(off)
	} else {
		// T3: a missing or unparseable startxref triggers full reconstruction.
		d.addWarning(WarnXRefBruteForce, -1, "no usable startxref; rebuilding by brute force")
	}
	if len(sections) == 0 {
		return d.rebuildFromBruteForce()
	}
	d.mergeSections(sections)
	if d.trailer == nil || !d.trailer.Has("Root") {
		// pdfbox: no /Root in the resolved trailer -> rebuild.
		d.addWarning(WarnXRefBruteForce, -1, "no /Root in the resolved trailer; rebuilding")
		return d.rebuildFromBruteForce()
	}
	d.checkXrefOffsets()
	return nil
}

// findStartXref reads the last min(size, 2048) bytes, finds the last %%EOF and
// the last `startxref` preceding it (T1, T2).
func (d *Document) findStartXref() (int64, bool) {
	n := int64(trailBytes)
	if int64(len(d.data)) < n {
		n = int64(len(d.data))
	}
	base := int64(len(d.data)) - n
	tail := d.data[base:]

	end := bytes.LastIndex(tail, []byte("%%EOF"))
	if end < 0 {
		// T1: a missing %%EOF is tolerated; search runs to end of file.
		d.addWarning(WarnMissingEOF, -1, "no %%EOF in the last 2048 bytes")
		end = len(tail)
	}
	idx := bytes.LastIndex(tail[:end], []byte("startxref"))
	if idx < 0 {
		return 0, false
	}
	l := newLexer(d.data, &d.warnings)
	l.seek(base + int64(idx) + int64(len("startxref")))
	t := l.next()
	if t.kind != tokInteger || t.i <= 0 {
		// pdfbox: Math.max(0, parseStartXref()) followed by `while (prev > 0)`,
		// so a non-positive startxref simply yields no sections.
		return 0, false
	}
	// An offset past EOF is *not* rejected here: checkXRefOffset repairs it by
	// brute force (X3), which is what saves validation/pades_infinite_loop.pdf,
	// whose startxref is 35274 in a 29469-byte file.
	return t.i, true
}

// walkChain follows /Prev from off, newest first.
func (d *Document) walkChain(off int64) []*rawSection {
	var out []*rawSection
	seen := make(map[int64]bool)
	// pdfbox's chain loop is `while (prev > 0)`: a /Prev of 0 — which
	// validation/pades-5-signatures-and-1-document-timestamp.pdf really contains,
	// as `/Prev 0` in its newest trailer — ends the chain. Treating 0 as an offset
	// and repairing it would graft an older revision's entries onto the newest
	// section and inflate the object count.
	for off > 0 && len(out) < maxSections {
		if seen[off] {
			break // cycle guard
		}
		seen[off] = true
		sec, err := d.parseSectionAt(off, false)
		if err != nil {
			// X3: repair a bad startxref / /Prev offset.
			fixed, ok := d.findXRefNear(off)
			if !ok {
				d.addWarning(WarnXRefStmSkipped, off,
					"xref section at "+strconv.FormatInt(off, 10)+" is unusable and could not be repaired")
				break
			}
			if seen[fixed] {
				break
			}
			seen[fixed] = true
			sec, err = d.parseSectionAt(fixed, true)
			if err != nil {
				break
			}
			d.addWarning(WarnXRefOffsetRepaired, off,
				"xref offset repaired to "+strconv.FormatInt(fixed, 10))
		}
		out = append(out, sec)
		if sec.sec.XRefStm >= 0 {
			d.hybrid = true
			d.loadXRefStm(sec)
		}
		off = sec.sec.Prev
	}
	return out
}

// loadXRefStm parses the hybrid /XRefStm section, repairing a bad offset by
// brute force and, failing that, skipping it (X4).
func (d *Document) loadXRefStm(sec *rawSection) {
	off := sec.sec.XRefStm
	if off <= 0 {
		d.addWarning(WarnXRefStmSkipped, off, "skipped XRef stream due to a corrupt offset")
		return
	}
	stm, err := d.parseSectionAt(off, false)
	if err != nil {
		fixed, ok := d.findXRefNear(off)
		if ok {
			stm, err = d.parseSectionAt(fixed, true)
		}
		if !ok || err != nil {
			d.addWarning(WarnXRefStmSkipped, off, "skipped XRef stream due to a corrupt offset")
			return
		}
		d.addWarning(WarnXRefOffsetRepaired, off, "/XRefStm offset repaired")
	}
	if stm.sec.Style != XRefStream {
		d.addWarning(WarnXRefStmSkipped, off, "/XRefStm does not point at an xref stream")
		return
	}
	sec.stmEntries = stm.entries
	// The /XRefStm's own trailer keys are not merged: pdfbox parses it into the
	// same XrefTrailerObj, and setTrailer is only called for the table's trailer.
}

// parseSectionAt parses one xref section, table or stream.
func (d *Document) parseSectionAt(off int64, recovered bool) (*rawSection, error) {
	if off < 0 || off >= int64(len(d.data)) {
		return nil, fmt.Errorf("pdf: xref offset %d out of range", off)
	}
	p := d.newParser()
	p.lex.seek(off)
	p.lex.skipSpaces()
	if hasPrefixAt(p.lex.data, p.lex.pos, "xref") {
		return d.parseXrefTableAt(p, off, recovered)
	}
	return d.parseXrefStreamAt(p, off, recovered)
}

// parseXrefTableAt implements COSParser.parseXrefTable + parseTrailer.
func (d *Document) parseXrefTableAt(p *parser, off int64, recovered bool) (*rawSection, error) {
	rs := &rawSection{sec: XRefSection{
		Offset: off, Style: XRefTable, Prev: -1, XRefStm: -1, Recovered: recovered,
	}}
	l := p.lex
	l.pos += len("xref")

	for {
		l.skipSpaces()
		if hasPrefixAt(l.data, l.pos, "trailer") {
			break
		}
		first, count, ok := d.readSubsectionHeader(l)
		if !ok {
			// A header that does not split into exactly two integers aborts the
			// section (warning), exactly as pdfbox returns false.
			d.addWarning(WarnXRefEntryInvalid, int64(l.pos), "unexpected xref subsection header")
			break
		}
		objID := first
		for i := int64(0); i < count; i++ {
			l.skipSpaces()
			if l.eof() || l.data[l.pos] == 't' || isDelimiter(l.data[l.pos]) {
				break // early termination on `trailer` or a delimiter
			}
			line := readLineAt(l)
			fields := splitOnSpace(line)
			if len(fields) < 3 {
				// PDFBOX-474 tolerance.
				d.addWarning(WarnXRefEntryInvalid, int64(l.pos), "invalid xref line: "+line)
				break
			}
			switch {
			case fields[len(fields)-1] == "n":
				offVal, err1 := strconv.ParseInt(fields[0], 10, 64)
				genVal, err2 := strconv.ParseInt(fields[1], 10, 32)
				if err1 == nil && err2 == nil && offVal > 0 {
					rs.entries = append(rs.entries, keyedEntry{
						key: ObjectKey{Num: objID, Gen: uint16(genVal)},
						e:   xrefRec{typ: 1, offset: offVal, gen: uint16(genVal)},
					})
				}
			case fields[2] == "f":
				// free entries are counted but never recorded
				d.noteObjectNumber(objID)
			default:
				d.addWarning(WarnXRefEntryInvalid, int64(l.pos),
					"corrupt xref entry for object "+strconv.FormatInt(objID, 10))
			}
			d.noteObjectNumber(objID)
			objID++
		}
		l.skipSpaces()
		if l.eof() || !isDigit(l.data[l.pos]) {
			break
		}
	}

	// trailer
	l.skipSpaces()
	if !hasPrefixAt(l.data, l.pos, "trailer") {
		d.addWarning(WarnXRefEntryInvalid, int64(l.pos), "xref table without a trailer")
		rs.sec.Trailer = NewDict()
		rs.sec.Entries = len(rs.entries)
		if len(rs.entries) == 0 {
			return nil, fmt.Errorf("pdf: empty xref table at %d", off)
		}
		return rs, nil
	}
	l.pos += len("trailer")
	l.skipSpaces()
	t := l.next()
	if t.kind != tokDictOpen {
		return nil, fmt.Errorf("pdf: trailer at %d is not a dictionary", off)
	}
	tr := p.parseDictBody()
	rs.sec.Trailer = tr
	rs.sec.Prev = dictInt(tr, "Prev", -1)
	rs.sec.XRefStm = dictInt(tr, "XRefStm", -1)
	rs.sec.Entries = len(rs.entries)
	d.noteSize(tr)
	return rs, nil
}

// readSubsectionHeader reads `first count`.
func (d *Document) readSubsectionHeader(l *lexer) (int64, int64, bool) {
	save := l.pos
	line := readLineAt(l)
	fields := splitOnSpace(line)
	if len(fields) != 2 {
		l.pos = save
		return 0, 0, false
	}
	first, err1 := strconv.ParseInt(fields[0], 10, 64)
	count, err2 := strconv.ParseInt(fields[1], 10, 64)
	if err1 != nil || err2 != nil || count < 0 {
		l.pos = save
		return 0, 0, false
	}
	return first, count, true
}

// parseXrefStreamAt parses an `N G obj << /Type /XRef … >> stream` section.
func (d *Document) parseXrefStreamAt(p *parser, off int64, recovered bool) (*rawSection, error) {
	obj, _, ok := p.parseIndirectAt(off, ObjectKey{})
	if !ok {
		return nil, fmt.Errorf("pdf: no object at xref offset %d", off)
	}
	st, isStream := obj.(*Stream)
	if !isStream {
		return nil, fmt.Errorf("pdf: object at xref offset %d is not a stream", off)
	}
	if n, ok := st.Dict.GetRaw("Type").(Name); ok && n != "XRef" {
		return nil, fmt.Errorf("pdf: object at xref offset %d is /Type /%s", off, n)
	}
	rs := &rawSection{sec: XRefSection{
		Offset: off, Style: XRefStream, Trailer: st.Dict, Recovered: recovered,
		Prev: dictInt(st.Dict, "Prev", -1), XRefStm: -1,
	}}
	// A cross-reference stream is never encrypted, so decoding needs no key.
	names, parms := streamFilters(st.Dict, nil)
	data, err := Decode(st.Raw, names, parms, &d.warnings)
	if err != nil {
		return nil, err
	}
	entries, err := d.decodeXRefStream(st.Dict, data)
	if err != nil {
		return nil, err
	}
	rs.entries = entries
	rs.sec.Entries = len(entries)
	d.noteSize(st.Dict)
	return rs, nil
}

// decodeXRefStream reads /W-wide big-endian fields over the /Index ranges.
func (d *Document) decodeXRefStream(dict *Dict, data []byte) ([]keyedEntry, error) {
	wArr, ok := dict.GetRaw("W").(Array)
	if !ok || len(wArr) != 3 {
		return nil, fmt.Errorf("pdf: bad /W array in xref stream")
	}
	var w [3]int
	total := 0
	for i := 0; i < 3; i++ {
		v, ok := wArr[i].(Integer)
		if !ok || v < 0 {
			return nil, fmt.Errorf("pdf: bad /W array in xref stream")
		}
		w[i] = int(v)
		total += w[i]
	}
	if total == 0 || total > 20 { // PDFBOX-6037
		return nil, fmt.Errorf("pdf: bad /W array in xref stream")
	}
	var index []int64
	if arr, ok := dict.GetRaw("Index").(Array); ok && len(arr) > 0 && len(arr)%2 == 0 {
		for _, e := range arr {
			v, _ := e.(Integer)
			index = append(index, int64(v))
		}
	} else {
		size := dictInt(dict, "Size", 0)
		index = []int64{0, size}
	}

	var out []keyedEntry
	pos := 0
	read := func(off, n int) int64 {
		var v int64
		for i := 0; i < n; i++ {
			v = v<<8 | int64(data[pos+off+i])
		}
		return v
	}
	for r := 0; r+1 < len(index); r += 2 {
		first, count := index[r], index[r+1]
		for i := int64(0); i < count; i++ {
			if pos+total > len(data) {
				return out, nil // truncated stream: keep what we have
			}
			num := first + i
			d.noteObjectNumber(num)
			typ := int64(1)
			if w[0] != 0 {
				typ = read(0, w[0])
			}
			second := read(w[0], w[1])
			third := read(w[0]+w[1], w[2])
			pos += total
			switch typ {
			case 0:
				// free: skipped, exactly as PDFXrefStreamParser does
			case 1:
				out = append(out, keyedEntry{
					key: ObjectKey{Num: num, Gen: uint16(third)},
					e:   xrefRec{typ: 1, offset: second, gen: uint16(third)},
				})
			case 2:
				out = append(out, keyedEntry{
					key: ObjectKey{Num: num},
					e:   xrefRec{typ: 2, offset: second, index: int(third)},
				})
			default:
				// unknown type: skipped
			}
		}
	}
	return out, nil
}

// mergeSections resolves the chain exactly as XrefTrailerResolver.setStartxref
// does: walk oldest to newest, letting newer entries and newer trailer values win.
//
// Note the divergence from DESIGN.md §2.3(8), which describes the resolved
// trailer as "the newest section's trailer with /Root and /Info inherited".
// pdfbox actually unions every trailer in the chain (COSDictionary.addAll,
// oldest first), so the resolved trailer's key set is the union — which is what
// the oracle's `trailer` field reports and what the KAT compares against.
func (d *Document) mergeSections(sections []*rawSection) {
	d.xref = make(map[ObjectKey]xrefRec, 64)
	d.byNum = make(map[int64]ObjectKey, 64)
	trailer := NewDict()
	for i := len(sections) - 1; i >= 0; i-- {
		rs := sections[i]
		// Within a section the table wins over the hybrid /XRefStm (PDFBOX-3506).
		local := make(map[ObjectKey]bool, len(rs.entries))
		for _, ke := range rs.entries {
			local[ke.key] = true
			d.putXRef(ke)
		}
		for _, ke := range rs.stmEntries {
			if local[ke.key] {
				continue
			}
			d.putXRef(ke)
		}
		if rs.sec.Trailer != nil {
			for _, k := range rs.sec.Trailer.Keys() {
				trailer.Set(k, rs.sec.Trailer.GetRaw(k))
			}
		}
	}
	d.trailer = trailer
	d.sections = make([]XRefSection, 0, len(sections))
	for _, rs := range sections {
		d.sections = append(d.sections, rs.sec)
	}
}

func (d *Document) putXRef(ke keyedEntry) {
	d.xref[ke.key] = ke.e
	d.byNum[ke.key.Num] = ke.key
	d.noteObjectNumber(ke.key.Num)
}

// noteObjectNumber tracks the highest object number seen anywhere, including free
// entries — R16 needs the maximum over all sections, not just the live entries.
func (d *Document) noteObjectNumber(num int64) {
	if num > d.highestObjNum {
		d.highestObjNum = num
	}
}

func (d *Document) noteSize(trailer *Dict) {
	if size := dictInt(trailer, "Size", 0); size > 0 {
		d.noteObjectNumber(size - 1)
	}
}

// checkXrefOffsets is COSParser.checkXrefOffsets + validateXrefOffsets (X1/X2):
// every type-1 entry is validated by seeking to its offset and reading the
// `N G obj` header.
//
//   - header matches            -> the key is valid;
//   - header names a different  -> the key is *corrected* to what the header
//     object number or a higher    says, and the original key is dropped. In
//     generation                   lenient mode a wrong object number is not a
//     failure: pdfbox adopts the number it found. This is what keeps
//     validation/dss-2236/hide.pdf at 19 objects — two xref entries there point
//     at the same offset and the loser is simply removed;
//   - nothing parses            -> validation fails for the whole table and the
//     brute-force map replaces it wholesale.
func (d *Document) checkXrefOffsets() {
	p := d.newParser()
	type fix struct{ from, to ObjectKey }
	var fixes []fix
	valid := make(map[ObjectKey]bool)
	for _, key := range d.sortedXRefKeys() {
		e := d.xref[key]
		if e.typ != 1 {
			// A type-2 entry holds an object number, not an offset; pdfbox skips
			// negative "offsets" here for the same reason.
			continue
		}
		found, ok := d.findObjectKey(p, key, e.offset)
		if !ok {
			d.rebuildXRefFromBruteForce()
			return
		}
		if found == key {
			valid[key] = true
			continue
		}
		fixes = append(fixes, fix{from: key, to: found})
	}
	// Only replace an entry when the corrected key does not already point at a
	// valid object; then drop every corrected original, as pdfbox does.
	type pointer struct {
		key ObjectKey
		e   xrefRec
	}
	var pointers []pointer
	for _, f := range fixes {
		if !valid[f.to] {
			pointers = append(pointers, pointer{key: f.to, e: d.xref[f.from]})
		}
	}
	for _, f := range fixes {
		delete(d.xref, f.from)
		if d.byNum[f.from.Num] == f.from {
			delete(d.byNum, f.from.Num)
		}
		d.addWarning(WarnXRefOffsetRepaired, d.startxref,
			"xref key "+f.from.String()+" corrected to "+f.to.String())
	}
	for _, ptr := range pointers {
		e := ptr.e
		e.gen = ptr.key.Gen
		d.xref[ptr.key] = e
		d.byNum[ptr.key.Num] = ptr.key
		d.noteObjectNumber(ptr.key.Num)
	}
}

// findObjectKey is COSParser.findObjectKey: it returns the key actually found at
// off — which may differ from the requested one in number, generation or both —
// or ok=false when nothing usable is there.
func (d *Document) findObjectKey(p *parser, key ObjectKey, off int64) (ObjectKey, bool) {
	// There can't be any object at the very beginning of a PDF.
	if off < minimumSearchOffset || off >= int64(len(d.data)) {
		return key, false
	}
	p.lex.seek(off)
	p.lex.skipSpaces()
	if int64(p.lex.pos) == off && off > 0 && isDigit(d.data[off-1]) {
		// No whitespace in front of the object number and a digit right before
		// it: this may be the tail of the *previous* object. Read that object's
		// key and, if the xref already places it within 10 bytes of this offset,
		// something is definitely wrong and the whole table is rejected.
		j := int(off) - 1
		for j >= 0 && isDigit(d.data[j]) {
			j--
		}
		q := d.newParser()
		q.lex.seek(int64(j + 1))
		t1 := q.lex.next()
		t2 := q.lex.next()
		if t1.kind == tokInteger && t2.kind == tokInteger {
			other := ObjectKey{Num: t1.i, Gen: uint16(t2.i)}
			if e, exists := d.xref[other]; exists && e.typ == 1 && e.offset > 0 && absInt64(off-e.offset) < 10 {
				return key, false
			}
		}
		p.lex.seek(off)
	}
	t1 := p.lex.next()
	if t1.kind != tokInteger {
		return key, false
	}
	if t1.i != key.Num {
		// Lenient mode adopts the object number that is actually there.
		key = ObjectKey{Num: t1.i, Gen: key.Gen}
	}
	t2 := p.lex.next()
	if t2.kind != tokInteger {
		return key, false
	}
	t3 := p.lex.next()
	if t3.kind != tokKeyword || t3.raw != "obj" {
		return key, false
	}
	gen := uint16(t2.i)
	switch {
	case gen == key.Gen:
		return key, true
	case gen > key.Gen:
		return ObjectKey{Num: key.Num, Gen: gen}, true
	default:
		return key, false
	}
}

func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// sortedXRefKeys returns the xref keys in ascending (Num, Gen) order, so every
// traversal of the map is deterministic (R21).
func (d *Document) sortedXRefKeys() []ObjectKey {
	keys := make([]ObjectKey, 0, len(d.xref))
	for k := range d.xref {
		keys = append(keys, k)
	}
	sortObjectKeys(keys)
	return keys
}

// sortObjectKeys orders keys ascending by Num, then Gen. Every map traversal in
// this package goes through it so no output can depend on map order (R21).
func sortObjectKeys(keys []ObjectKey) {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Num != keys[j].Num {
			return keys[i].Num < keys[j].Num
		}
		return keys[i].Gen < keys[j].Gen
	})
}

func dictInt(d *Dict, key Name, def int64) int64 {
	if d == nil {
		return def
	}
	if v, ok := d.GetRaw(key).(Integer); ok {
		return int64(v)
	}
	return def
}

// readLineAt reads to the next EOL, consuming it.
func readLineAt(l *lexer) string {
	start := l.pos
	for l.pos < len(l.data) && l.data[l.pos] != '\r' && l.data[l.pos] != '\n' {
		l.pos++
	}
	s := string(l.data[start:l.pos])
	if l.pos < len(l.data) {
		if l.data[l.pos] == '\r' {
			l.pos++
			if l.pos < len(l.data) && l.data[l.pos] == '\n' {
				l.pos++
			}
		} else {
			l.pos++
		}
	}
	return s
}

// splitOnSpace is org.apache.pdfbox.util.StringUtil.splitOnSpace: split on runs
// of 0x20, dropping empty fields.
func splitOnSpace(s string) []string {
	var out []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n' || s[i] == 0 {
			if start >= 0 {
				out = append(out, s[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}

// Brute-force recovery: DESIGN.md §2.7 X2–X4.
//
// Provenance: org.apache.pdfbox.pdfparser.BruteForceParser of pdfbox 3.0.7
// (bfSearchForObjects, bfSearchForXRef, bfSearchForXRefTables,
// bfSearchForXRefStreams, bfSearchForObjStreams, rebuildTrailer,
// searchForTrailerItems).

package pdf

import (
	"bytes"
	"strconv"
)

// bruteForce holds the results of the whole-file object scan, computed at most
// once per document.
type bruteForce struct {
	done    bool
	offsets map[ObjectKey]int64
	// xrefCandidates are the offsets of `xref` keywords and of `/Type /XRef`
	// objects, used to repair a bad startxref, /Prev or /XRefStm (X3, X4).
	xrefCandidates []int64
}

// minimumSearchOffset is BruteForceParser.MINIMUM_SEARCH_OFFSET: nothing useful
// can start before byte 6 (`%PDF-x`).
const minimumSearchOffset = 6

func (d *Document) bf() *bruteForce {
	if d.brute.done {
		return &d.brute
	}
	d.brute.done = true
	d.brute.offsets = make(map[ObjectKey]int64)
	data := d.data

	// pdfbox stops the scan at the last %%EOF marker.
	limit := len(data)
	if i := bytes.LastIndex(data, []byte("%%EOF")); i >= 0 {
		limit = i + len("%%EOF")
	}

	for i := minimumSearchOffset; i+3 < limit; i++ {
		if data[i] != 'o' || !hasPrefixAt(data, i, "obj") {
			continue
		}
		// walk back over whitespace, one generation digit, whitespace, digits
		j := i - 1
		if j < minimumSearchOffset || !isWhitespace(data[j]) {
			continue
		}
		for j >= minimumSearchOffset && isWhitespace(data[j]) {
			j--
		}
		genEnd := j
		for j >= minimumSearchOffset && isDigit(data[j]) {
			j--
		}
		if genEnd == j {
			continue
		}
		genStr := string(data[j+1 : genEnd+1])
		if j < minimumSearchOffset || !isWhitespace(data[j]) {
			continue
		}
		for j >= minimumSearchOffset && isWhitespace(data[j]) {
			j--
		}
		numEnd := j
		for j >= minimumSearchOffset && isDigit(data[j]) {
			j--
		}
		if numEnd == j {
			continue
		}
		numStr := string(data[j+1 : numEnd+1])
		num, err1 := strconv.ParseInt(numStr, 10, 64)
		gen, err2 := strconv.ParseInt(genStr, 10, 32)
		if err1 != nil || err2 != nil || gen > 65535 {
			continue
		}
		// Last occurrence wins: later revisions are later in the file.
		d.brute.offsets[ObjectKey{Num: num, Gen: uint16(gen)}] = int64(j + 1)
	}

	// xref candidates: `xref` keywords …
	for i := 0; i+4 <= len(data); i++ {
		if data[i] == 'x' && hasPrefixAt(data, i, "xref") {
			if i > 0 && isRegular(data[i-1]) {
				continue
			}
			d.brute.xrefCandidates = append(d.brute.xrefCandidates, int64(i))
		}
	}
	// … and objects whose dictionary is /Type /XRef.
	keys := make([]ObjectKey, 0, len(d.brute.offsets))
	for k := range d.brute.offsets {
		keys = append(keys, k)
	}
	sortObjectKeys(keys)
	p := d.newParser()
	for _, k := range keys {
		off := d.brute.offsets[k]
		if !d.looksLikeXRefStream(p, off) {
			continue
		}
		d.brute.xrefCandidates = append(d.brute.xrefCandidates, off)
	}
	return &d.brute
}

// looksLikeXRefStream reports whether the object at off is `N G obj << /Type /XRef`.
func (d *Document) looksLikeXRefStream(p *parser, off int64) bool {
	p.lex.seek(off)
	p.lex.skipSpaces()
	if t := p.lex.next(); t.kind != tokInteger {
		return false
	}
	if t := p.lex.next(); t.kind != tokInteger {
		return false
	}
	if t := p.lex.next(); t.kind != tokKeyword || t.raw != "obj" {
		return false
	}
	p.lex.skipSpaces()
	if t := p.lex.next(); t.kind != tokDictOpen {
		return false
	}
	// Only look at the first few keys: a full parse here would be wasteful and
	// pdfbox's bfSearchForXRefStreams is a byte scan for "/Type" + "/XRef" too.
	for i := 0; i < 32; i++ {
		t := p.lex.next()
		if t.kind != tokName {
			if t.kind == tokDictClose || t.kind == tokEOF {
				return false
			}
			continue
		}
		if t.name != "Type" {
			// skip the value
			p.parseDirObject()
			continue
		}
		v := p.parseDirObject()
		n, ok := v.(Name)
		return ok && n == "XRef"
	}
	return false
}

// findXRefNear is BruteForceParser.bfSearchForXRef: pick the xref candidate
// closest to the claimed offset (X3/X4).
func (d *Document) findXRefNear(off int64) (int64, bool) {
	b := d.bf()
	best := int64(-1)
	var bestDist int64
	for _, c := range b.xrefCandidates {
		dist := c - off
		if dist < 0 {
			dist = -dist
		}
		if best < 0 || dist < bestDist {
			best, bestDist = c, dist
		}
	}
	if best < 0 {
		return 0, false
	}
	return best, true
}

// rebuildXRefFromBruteForce replaces the whole xref map with the scan results
// (X2). The trailer is left alone: pdfbox only rebuilds it when it has no /Root.
func (d *Document) rebuildXRefFromBruteForce() {
	b := d.bf()
	if len(b.offsets) == 0 {
		return
	}
	d.addWarning(WarnXRefBruteForce, -1, "replaced the xref table with a brute-force scan")
	d.xref = make(map[ObjectKey]xrefRec, len(b.offsets))
	d.byNum = make(map[int64]ObjectKey, len(b.offsets))
	keys := make([]ObjectKey, 0, len(b.offsets))
	for k := range b.offsets {
		keys = append(keys, k)
	}
	sortObjectKeys(keys)
	for _, k := range keys {
		d.putXRef(keyedEntry{key: k, e: xrefRec{typ: 1, offset: b.offsets[k], gen: k.Gen}})
	}
	d.registerObjectStreams()
}

// registerObjectStreams is bfSearchForObjStreams: every /Type /ObjStm found by
// the scan contributes type-2 entries for the objects it contains, unless the
// object number is already known.
func (d *Document) registerObjectStreams() {
	keys := d.sortedXRefKeys()
	for _, k := range keys {
		e := d.xref[k]
		if e.typ != 1 {
			continue
		}
		obj, ok := d.loadAt(e.offset, k)
		if !ok {
			continue
		}
		st, isStream := obj.(*Stream)
		if !isStream {
			continue
		}
		if n, _ := st.Dict.GetRaw("Type").(Name); n != "ObjStm" {
			continue
		}
		nums, err := d.objStmNumbers(k, st)
		if err != nil {
			continue
		}
		for i, num := range nums {
			if _, exists := d.byNum[num]; exists {
				continue
			}
			d.putXRef(keyedEntry{
				key: ObjectKey{Num: num},
				e:   xrefRec{typ: 2, offset: k.Num, index: i},
			})
		}
	}
}

// rebuildFromBruteForce is the full reconstruction path (T3/X2): rebuild the
// xref from the scan and then rebuild the trailer from the objects found.
func (d *Document) rebuildFromBruteForce() error {
	b := d.bf()
	if len(b.offsets) == 0 {
		return ErrNotPDF
	}
	d.xref = make(map[ObjectKey]xrefRec, len(b.offsets))
	d.byNum = make(map[int64]ObjectKey, len(b.offsets))
	d.trailerRebuilt = true
	d.addWarning(WarnXRefBruteForce, -1, "document rebuilt by brute-force object scan")
	keys := make([]ObjectKey, 0, len(b.offsets))
	for k := range b.offsets {
		keys = append(keys, k)
	}
	sortObjectKeys(keys)
	for _, k := range keys {
		d.putXRef(keyedEntry{key: k, e: xrefRec{typ: 1, offset: b.offsets[k], gen: k.Gen}})
	}
	d.registerObjectStreams()
	if d.trailer == nil {
		d.trailer = NewDict()
	}
	d.searchForTrailer()
	d.searchForTrailerItems()
	if d.sections == nil {
		d.sections = []XRefSection{{
			Offset: 0, Style: XRefTable, Trailer: d.trailer, Prev: -1, XRefStm: -1,
			Entries: len(d.xref), Recovered: true,
		}}
	}
	return nil
}

// searchForTrailer is bfSearchForTrailer: scan for `trailer` keywords and take
// /Root, /Info, /Encrypt and /ID from the dictionaries that follow.
func (d *Document) searchForTrailer() {
	p := d.newParser()
	for i := 0; i+7 <= len(d.data); i++ {
		if d.data[i] != 't' || !hasPrefixAt(d.data, i, "trailer") {
			continue
		}
		p.lex.seek(int64(i + len("trailer")))
		p.lex.skipSpaces()
		if t := p.lex.next(); t.kind != tokDictOpen {
			continue
		}
		td := p.parseDictBody()
		for _, k := range []Name{"Root", "Info", "Encrypt", "ID"} {
			if v := td.GetRaw(k); v != nil && !d.trailer.Has(k) {
				d.trailer.Set(k, v)
			}
		}
	}
}

// searchForTrailerItems is BruteForceParser.searchForTrailerItems: when the
// trailer still has no /Root, find the object whose /Type is /Catalog (and, for
// /Info, the one that looks like a document information dictionary).
func (d *Document) searchForTrailerItems() {
	if d.trailer.Has("Root") && d.resolveDict(d.trailer.GetRaw("Root")) != nil {
		return
	}
	for _, k := range d.sortedXRefKeys() {
		obj, ok := d.object(k)
		if !ok {
			continue
		}
		dict, isDict := obj.(*Dict)
		if !isDict {
			continue
		}
		if n, _ := dict.GetRaw("Type").(Name); n == "Catalog" {
			d.trailer.Set("Root", Ref{Num: k.Num, Gen: k.Gen})
			return
		}
	}
}

// resolveDict is a nil-safe helper used during recovery, before the public
// accessors are usable.
func (d *Document) resolveDict(o Object) *Dict {
	if o == nil {
		return nil
	}
	v, _ := d.Resolve(o).(*Dict)
	return v
}

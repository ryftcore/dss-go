// Cross-reference emission for an incremental update: both the table form and
// the stream form. See DESIGN.md §3.4 and rules R12, R13, R14.
//
// Provenance: COSWriter.doWriteXRefInc / doWriteXRefTable / doWriteTrailer /
// writeXrefRange / writeXrefEntry / getXRefRanges and
// org.apache.pdfbox.pdfwriter.compress-adjacent PDFXRefStream of pdfbox 3.0.7.
//
// Style is not a free choice (§3.4):
//
//	if previousSectionStyle == XRefTable || document.HasHybridXRef() {
//	        table + trailer          // a hybrid file always degrades to a table
//	} else {
//	        stream                   // /Type /XRef
//	}

package pdf

import (
	"errors"
	"sort"
	"strconv"
)

// xrefEntry is one row of the cross-reference section we are about to emit.
type xrefEntry struct {
	Key ObjectKey
	// Offset is the absolute byte offset of the object's "N G obj" header for
	// an in-use entry, and the next-free object number for a free entry.
	Offset int64
	Free   bool
}

// freeHeadEntry is the object-0 free head every incremental update carries
// (R12). pdfbox spells it FreeXReference.NULL_ENTRY: generation 65535, next
// free object 0.
func freeHeadEntry() xrefEntry {
	return xrefEntry{Key: ObjectKey{Num: 0, Gen: 65535}, Offset: 0, Free: true}
}

// sortEntries orders entries by object number, then generation. pdfbox relies
// on XReferenceEntry's natural ordering here; we sort explicitly because R15
// makes ascending object number the writer's only ordering.
func sortEntries(entries []xrefEntry) []xrefEntry {
	out := make([]xrefEntry, len(entries))
	copy(out, entries)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Key.Num != out[j].Key.Num {
			return out[i].Key.Num < out[j].Key.Num
		}
		return out[i].Key.Gen < out[j].Key.Gen
	})
	return out
}

// xrefRanges groups sorted entries into maximal contiguous runs of object
// numbers and returns them as {first, count} pairs — the subsection headers of
// R12, and the /Index of the stream form. This is COSWriter.getXRefRanges.
func xrefRanges(entries []xrefEntry) [][2]int64 {
	var ranges [][2]int64
	for i := 0; i < len(entries); {
		first := entries[i].Key.Num
		count := int64(1)
		j := i + 1
		for ; j < len(entries); j++ {
			if entries[j].Key.Num != first+count {
				break
			}
			count++
		}
		ranges = append(ranges, [2]int64{first, count})
		i = j
	}
	return ranges
}

// writeXRefTable emits the "xref" section (R12). entries must already include
// the object-0 free head; they are sorted here. Each entry occupies exactly 20
// bytes: %010d SP %05d SP [n|f] CR LF.
func writeXRefTable(w *Writer, entries []xrefEntry) error {
	sorted := sortEntries(entries)
	if err := w.writeString("xref"); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	i := 0
	for _, r := range xrefRanges(sorted) {
		if err := w.writeString(formatInt(r[0])); err != nil {
			return err
		}
		if err := w.writeByte(' '); err != nil {
			return err
		}
		if err := w.writeString(formatInt(r[1])); err != nil {
			return err
		}
		if err := w.WriteEOL(); err != nil {
			return err
		}
		for n := int64(0); n < r[1]; n++ {
			if err := writeXRefEntry(w, sorted[i]); err != nil {
				return err
			}
			i++
		}
	}
	return w.Err()
}

func writeXRefEntry(w *Writer, e xrefEntry) error {
	if err := w.WriteRaw(pad(e.Offset, 10)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	if err := w.WriteRaw(pad(int64(e.Key.Gen), 5)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	kind := byte('n')
	if e.Free {
		kind = 'f'
	}
	if err := w.writeByte(kind); err != nil {
		return err
	}
	// Exactly two bytes, and deliberately not WriteEOL: pdfbox's writeCRLF
	// leaves onNewLine false, which is why "trailer" follows the last entry
	// with no blank line between them.
	return w.writeCRLF()
}

// pad renders v right-aligned in width digits with leading zeros, the
// DecimalFormat("0000000000") / ("00000") of COSWriter. A value too wide for
// the field is emitted in full — an xref table cannot describe a file that
// large anyway, and truncating would produce a silently wrong offset.
func pad(v int64, width int) []byte {
	s := formatInt(v)
	if len(s) >= width {
		return []byte(s)
	}
	out := make([]byte, width)
	for i := range out {
		out[i] = '0'
	}
	copy(out[width-len(s):], s)
	return out
}

// formatInt is R7's rendering: strconv.FormatInt, base 10, nothing else.
func formatInt(v int64) string { return strconv.FormatInt(v, 10) }

// writeTrailer emits "trailer" plus the trailer dictionary (R13). The caller
// owns the dictionary's contents; see incremental.buildTrailer.
func writeTrailer(w *Writer, trailer *Dict) error {
	if err := w.writeString("trailer"); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	return w.writeDict(trailer)
}

// writeTail emits "startxref <offset> %%EOF" (R14).
func writeTail(w *Writer, startxref int64) error {
	if err := w.writeString("startxref"); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	if err := w.writeString(formatInt(startxref)); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	if err := w.writeString("%%EOF"); err != nil {
		return err
	}
	return w.WriteEOL()
}

// errXRefStreamTooWide guards a file too large for the widths we compute.
var errXRefStreamTooWide = errors.New("pdf: cross-reference offset exceeds 8 bytes")

// buildXRefStream builds the /Type /XRef stream object of §3.4. entries must
// include the object-0 free head and the xref stream object's own entry; size
// is the value for /Size (highest object number + 2, the extra one being the
// xref stream object itself, as pdfbox's pdfxRefStream.setSize(number + 2)).
//
// The payload is PNG-Up predicted (/Predictor 12) and then Flate encoded, which
// is what the corpus overwhelmingly contains on the reading side.
func buildXRefStream(entries []xrefEntry, trailer *Dict, size, prev int64) (*Stream, error) {
	sorted := sortEntries(entries)

	var maxOffset int64
	for _, e := range sorted {
		if e.Offset > maxOffset {
			maxOffset = e.Offset
		}
	}
	w2 := 1
	for v := maxOffset; v >= 256; v >>= 8 {
		w2++
	}
	if w2 > 8 {
		return nil, errXRefStreamTooWide
	}
	// /W [1 <n> 2]: one byte of type, n bytes of offset, two of generation.
	w1, w3 := 1, 2
	columns := w1 + w2 + w3

	raw := make([]byte, 0, len(sorted)*columns)
	for _, e := range sorted {
		typ := int64(1)
		if e.Free {
			typ = 0
		}
		raw = appendBE(raw, typ, w1)
		raw = appendBE(raw, e.Offset, w2)
		raw = appendBE(raw, int64(e.Key.Gen), w3)
	}

	d := NewDict()
	d.Set("Type", Name("XRef"))
	d.Set("Size", Integer(size))
	index := Array{}
	for _, r := range xrefRanges(sorted) {
		index = append(index, Integer(r[0]), Integer(r[1]))
	}
	d.Set("Index", index)
	d.Set("W", Array{Integer(w1), Integer(w2), Integer(w3)})
	d.Set("Filter", Name("FlateDecode"))
	d.Set("DecodeParms", DictOf(
		Name("Predictor"), Integer(12),
		Name("Columns"), Integer(columns),
	))
	// Trailer information travels in the stream dictionary itself
	// (PDFXRefStream.addTrailerInfo keeps exactly these keys).
	for _, k := range []Name{"Root", "Info", "Encrypt", "ID"} {
		if v := trailer.GetRaw(k); v != nil {
			d.Set(k, v)
		}
	}
	if prev >= 0 {
		d.Set("Prev", Integer(prev))
	}
	d.Set("Length", Integer(0)) // rewritten by Writer.writeStream (R11)

	return NewStream(d, FlateEncode(pngUpEncode(raw, columns))), nil
}

// appendBE appends v as width big-endian bytes.
func appendBE(dst []byte, v int64, width int) []byte {
	for i := width - 1; i >= 0; i-- {
		dst = append(dst, byte(v>>(8*i)))
	}
	return dst
}

// pngUpEncode applies the PNG "Up" filter (/Predictor 12) row by row, tagging
// each row with its filter type as the format requires. The inverse is the
// reader's ApplyPredictor.
func pngUpEncode(data []byte, columns int) []byte {
	if columns <= 0 || len(data) == 0 {
		return data
	}
	rows := len(data) / columns
	out := make([]byte, 0, rows*(columns+1))
	prev := make([]byte, columns)
	for r := 0; r < rows; r++ {
		row := data[r*columns : (r+1)*columns]
		out = append(out, 2) // PNG filter type 2 = Up
		for i := 0; i < columns; i++ {
			out = append(out, row[i]-prev[i])
		}
		prev = row
	}
	return out
}

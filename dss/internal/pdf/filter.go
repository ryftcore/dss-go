// Stream filters. See DESIGN.md §2.5 and §4.4.
//
// Provenance: org.apache.pdfbox.filter.{FlateFilter,FlateFilterDecoderStream,
// Predictor,LZWFilter,ASCIIHexFilter,ASCII85Filter,RunLengthDecodeFilter,
// CryptFilter} of pdfbox 3.0.7.
//
// We decode exactly the filters DSS ever consumes decoded bytes from: xref
// streams, object streams, /DSS token streams, /VRI /TS, and the non-emptiness
// test in PdfObjectModificationsFilter.isStreamFill. Image filters (DCTDecode,
// CCITTFaxDecode, JPXDecode, JBIG2Decode) are never decoded — stream comparison
// in DefaultPdfObjectModificationsFinder is on raw bytes — and yield a
// *FilterError instead.

package pdf

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"hash/adler32"
	"io"
)

// The filters we decode.
const (
	FilterFlate     Name = "FlateDecode"
	FilterLZW       Name = "LZWDecode"
	FilterASCIIHex  Name = "ASCIIHexDecode"
	FilterASCII85   Name = "ASCII85Decode"
	FilterRunLength Name = "RunLengthDecode"
	FilterCrypt     Name = "Crypt" // /Identity only; anything else -> ErrUnsupportedFilter
)

// filterAbbreviations are the inline-image short names; they appear in the wild
// on ordinary streams too, and pdfbox's FilterFactory accepts them.
var filterAbbreviations = map[Name]Name{
	"Fl":  FilterFlate,
	"LZW": FilterLZW,
	"AHx": FilterASCIIHex,
	"A85": FilterASCII85,
	"RL":  FilterRunLength,
	"CCF": "CCITTFaxDecode",
	"DCT": "DCTDecode",
}

func canonicalFilterName(n Name) Name {
	if full, ok := filterAbbreviations[n]; ok {
		return full
	}
	return n
}

// Decode applies the filter chain named by /Filter with the parameters in
// /DecodeParms to raw. Unsupported filters yield a *FilterError. Any recovered
// defect appends to warn.
func Decode(raw []byte, filters []Name, parms []*Dict, warn *[]Warning) ([]byte, error) {
	out := raw
	for i, f := range filters {
		var parm *Dict
		if i < len(parms) {
			parm = parms[i]
		}
		var err error
		out, err = decodeOne(out, canonicalFilterName(f), parm, warn)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func decodeOne(raw []byte, f Name, parm *Dict, warn *[]Warning) ([]byte, error) {
	switch f {
	case FilterFlate:
		return applyPredictorFromParms(FlateDecode(raw, warn), parm, warn), nil
	case FilterLZW:
		early := 1
		if parm != nil {
			if v, ok := parm.GetRaw("EarlyChange").(Integer); ok {
				early = int(v)
			}
		}
		return applyPredictorFromParms(lzwDecode(raw, early, warn), parm, warn), nil
	case FilterASCIIHex:
		return asciiHexDecode(raw), nil
	case FilterASCII85:
		return ascii85Decode(raw), nil
	case FilterRunLength:
		return runLengthDecode(raw), nil
	case FilterCrypt:
		// Only /Identity is meaningful here: a real /Crypt filter means the stream
		// was decrypted by the security handler already (pdfbox's CryptFilter is a
		// no-op for Identity and unsupported otherwise).
		if parm != nil {
			if n, ok := parm.GetRaw("Name").(Name); ok && n != "Identity" {
				return nil, &FilterError{Filter: f}
			}
		}
		return raw, nil
	case "":
		return raw, nil
	default:
		return nil, &FilterError{Filter: f}
	}
}

func applyPredictorFromParms(data []byte, parm *Dict, warn *[]Warning) []byte {
	if parm == nil {
		return data
	}
	predictor := intFromDict(parm, "Predictor", 1)
	if predictor <= 1 {
		return data
	}
	colors := intFromDict(parm, "Colors", 1)
	bpc := intFromDict(parm, "BitsPerComponent", 8)
	columns := intFromDict(parm, "Columns", 1)
	return ApplyPredictor(data, predictor, colors, bpc, columns, warn)
}

func intFromDict(d *Dict, key Name, def int) int {
	if d == nil {
		return def
	}
	if v, ok := d.GetRaw(key).(Integer); ok {
		return int(v)
	}
	return def
}

// FlateDecode reproduces pdfbox FlateFilterDecoderStream exactly: it discards the
// first two bytes, inflates as raw DEFLATE with no zlib-header and no Adler-32
// validation, and on a corrupt stream returns the bytes decoded so far with a
// WarnFlateTruncated warning and no error.
//
// Do not "fix" this to use compress/zlib: the corpus contains non-conforming
// headers and truncated checksums that zlib.NewReader rejects and pdfbox accepts.
func FlateDecode(raw []byte, warn *[]Warning) []byte {
	if len(raw) < 2 {
		// Fewer than two bytes of input: empty output, no error.
		return []byte{}
	}
	r := flate.NewReader(bytes.NewReader(raw[2:]))
	defer r.Close()
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		addWarning(warn, WarnFlateTruncated, -1,
			"premature end of flate stream: "+err.Error())
	}
	return buf.Bytes()
}

// FlateEncode produces zlib-wrapped DEFLATE at the fixed compression level 6.
// Output is byte-stable for a given toolchain but NOT across Go releases —
// Go 1.27 changed compress/flate's encoder output — which is fine: nothing
// signs or pins these compressed bytes (signatures cover ByteRanges and
// decompressed content), only self-consistency within one produced revision.
// The zlib wrapper is written by hand rather than taken from compress/zlib so the
// header bytes are pinned here and cannot drift.
func FlateEncode(data []byte) []byte {
	var buf bytes.Buffer
	// zlib header: CM=8, CINFO=7 (32K window), FLEVEL=2 for level 6, FCHECK such
	// that the 16-bit big-endian value is a multiple of 31.
	buf.Write([]byte{0x78, 0x9c})
	w, err := flate.NewWriter(&buf, 6)
	if err != nil {
		// level 6 is always valid
		panic("pdf: flate.NewWriter: " + err.Error())
	}
	_, _ = w.Write(data)
	_ = w.Close()
	var sum [4]byte
	binary.BigEndian.PutUint32(sum[:], adler32.Checksum(data))
	buf.Write(sum[:])
	return buf.Bytes()
}

// ApplyPredictor undoes /Predictor: 1 = none, 2 = TIFF, 10..15 = PNG (each row
// carries its own filter-type byte, so a declared 15 still dispatches per row
// over 0..4). Row length is ceil(Colors*BPC*Columns/8); bpp = ceil(Colors*BPC/8).
// A truncated final row is zero-padded and warned about.
func ApplyPredictor(data []byte, predictor, colors, bpc, columns int, warn *[]Warning) []byte {
	if predictor <= 1 {
		return data
	}
	if colors <= 0 || bpc <= 0 || columns <= 0 {
		return data
	}
	rowLen := (colors*bpc*columns + 7) / 8
	if rowLen <= 0 {
		return data
	}
	if predictor == 2 {
		return tiffPredictor(data, colors, bpc, columns, rowLen)
	}
	bpp := (colors*bpc + 7) / 8
	var out []byte
	prev := make([]byte, rowLen)
	cur := make([]byte, rowLen)
	pos := 0
	for pos < len(data) {
		ft := int(data[pos])
		pos++
		n := copy(cur, data[pos:])
		if n < rowLen {
			for i := n; i < rowLen; i++ {
				cur[i] = 0
			}
			addWarning(warn, WarnPredictorTruncated, -1, "truncated predictor row zero-padded")
		}
		pos += n
		switch ft {
		case 0: // None
		case 1: // Sub
			for i := bpp; i < rowLen; i++ {
				cur[i] += cur[i-bpp]
			}
		case 2: // Up
			for i := 0; i < rowLen; i++ {
				cur[i] += prev[i]
			}
		case 3: // Average
			for i := 0; i < rowLen; i++ {
				var left int
				if i >= bpp {
					left = int(cur[i-bpp])
				}
				cur[i] += byte((left + int(prev[i])) / 2)
			}
		case 4: // Paeth
			for i := 0; i < rowLen; i++ {
				var left, upLeft int
				if i >= bpp {
					left = int(cur[i-bpp])
					upLeft = int(prev[i-bpp])
				}
				cur[i] += byte(paeth(left, int(prev[i]), upLeft))
			}
		default:
			addWarning(warn, WarnPredictorTruncated, -1, "unknown PNG predictor row type")
		}
		out = append(out, cur...)
		copy(prev, cur)
		if n < rowLen {
			break
		}
	}
	return out
}

func paeth(a, b, c int) int {
	p := a + b - c
	pa, pb, pc := abs(p-a), abs(p-b), abs(p-c)
	if pa <= pb && pa <= pc {
		return a
	}
	if pb <= pc {
		return b
	}
	return c
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// tiffPredictor undoes TIFF predictor 2. It has no corpus coverage; it is
// implemented for completeness and has its own unit test.
func tiffPredictor(data []byte, colors, bpc, columns, rowLen int) []byte {
	out := make([]byte, len(data))
	copy(out, data)
	rows := len(out) / rowLen
	switch bpc {
	case 8:
		for r := 0; r < rows; r++ {
			row := out[r*rowLen : (r+1)*rowLen]
			for i := colors; i < len(row); i++ {
				row[i] += row[i-colors]
			}
		}
	default:
		// Sub-byte components: operate on the bit stream.
		for r := 0; r < rows; r++ {
			row := out[r*rowLen : (r+1)*rowLen]
			br := newBitReader(row, bpc)
			vals := make([]int, 0, columns*colors)
			for i := 0; i < columns*colors; i++ {
				vals = append(vals, br.read())
			}
			mask := (1 << bpc) - 1
			for i := colors; i < len(vals); i++ {
				vals[i] = (vals[i] + vals[i-colors]) & mask
			}
			bw := newBitWriter(row, bpc)
			for _, v := range vals {
				bw.write(v)
			}
		}
	}
	return out
}

type bitReader struct {
	data []byte
	bits int
	pos  int // bit position
}

func newBitReader(d []byte, bits int) *bitReader { return &bitReader{data: d, bits: bits} }

func (b *bitReader) read() int {
	v := 0
	for i := 0; i < b.bits; i++ {
		byteIdx := b.pos >> 3
		if byteIdx >= len(b.data) {
			return v << (b.bits - i)
		}
		bit := (b.data[byteIdx] >> (7 - uint(b.pos&7))) & 1
		v = v<<1 | int(bit)
		b.pos++
	}
	return v
}

type bitWriter struct {
	data []byte
	bits int
	pos  int
}

func newBitWriter(d []byte, bits int) *bitWriter { return &bitWriter{data: d, bits: bits} }

func (b *bitWriter) write(v int) {
	for i := b.bits - 1; i >= 0; i-- {
		byteIdx := b.pos >> 3
		if byteIdx >= len(b.data) {
			return
		}
		bit := byte((v >> uint(i)) & 1)
		shift := uint(7 - (b.pos & 7))
		b.data[byteIdx] = (b.data[byteIdx] &^ (1 << shift)) | bit<<shift
		b.pos++
	}
}

// asciiHexDecode implements ASCIIHexDecode: hex digits until '>', whitespace
// skipped, an odd trailing digit padded with '0' (ASCIIHexFilter).
func asciiHexDecode(raw []byte) []byte {
	var digits []byte
	for _, c := range raw {
		switch {
		case isHexDigit(c):
			digits = append(digits, c)
		case c == '>':
			return decodeHexDigits(digits)
		default:
			// whitespace and junk alike are skipped
		}
	}
	return decodeHexDigits(digits)
}

// ascii85Decode implements ASCII85Decode (ASCII85InputStream).
func ascii85Decode(raw []byte) []byte {
	var out []byte
	var group [5]byte
	n := 0
	i := 0
	// an optional <~ prefix is tolerated
	if len(raw) >= 2 && raw[0] == '<' && raw[1] == '~' {
		i = 2
	}
	for ; i < len(raw); i++ {
		c := raw[i]
		if isWhitespace(c) {
			continue
		}
		if c == '~' {
			break
		}
		if c == 'z' && n == 0 {
			out = append(out, 0, 0, 0, 0)
			continue
		}
		if c < '!' || c > 'u' {
			continue // junk is skipped
		}
		group[n] = c - '!'
		n++
		if n == 5 {
			v := uint32(0)
			for _, g := range group {
				v = v*85 + uint32(g)
			}
			out = append(out, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			n = 0
		}
	}
	if n > 0 {
		for j := n; j < 5; j++ {
			group[j] = 84
		}
		v := uint32(0)
		for _, g := range group {
			v = v*85 + uint32(g)
		}
		b := [4]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		out = append(out, b[:n-1]...)
	}
	return out
}

// runLengthDecode implements RunLengthDecode (RunLengthDecodeFilter).
func runLengthDecode(raw []byte) []byte {
	var out []byte
	for i := 0; i < len(raw); {
		l := int(raw[i])
		i++
		if l == 128 {
			break
		}
		if l < 128 {
			n := l + 1
			if i+n > len(raw) {
				n = len(raw) - i
			}
			out = append(out, raw[i:i+n]...)
			i += n
		} else {
			if i >= len(raw) {
				break
			}
			b := raw[i]
			i++
			for j := 0; j < 257-l; j++ {
				out = append(out, b)
			}
		}
	}
	return out
}

// lzwDecode implements LZWDecode (LZWFilter). earlyChange defaults to 1.
func lzwDecode(raw []byte, earlyChange int, warn *[]Warning) []byte {
	if earlyChange != 0 {
		earlyChange = 1
	}
	var out []byte
	table := make([][]byte, 0, 4096)
	reset := func() {
		table = table[:0]
		for i := 0; i < 256; i++ {
			table = append(table, []byte{byte(i)})
		}
		table = append(table, nil, nil) // 256 = clear, 257 = EOD
	}
	reset()
	codeLen := 9
	var prev []byte
	bitPos := 0
	readCode := func() int {
		v := 0
		for i := 0; i < codeLen; i++ {
			byteIdx := bitPos >> 3
			if byteIdx >= len(raw) {
				return 257
			}
			bit := (raw[byteIdx] >> (7 - uint(bitPos&7))) & 1
			v = v<<1 | int(bit)
			bitPos++
		}
		return v
	}
	for {
		code := readCode()
		if code == 257 {
			break
		}
		if code == 256 {
			reset()
			codeLen = 9
			prev = nil
			continue
		}
		var entry []byte
		switch {
		case code < len(table) && table[code] != nil:
			entry = table[code]
		case prev != nil:
			entry = append(append([]byte{}, prev...), prev[0])
		default:
			addWarning(warn, WarnFlateTruncated, -1, "invalid LZW code")
			return out
		}
		out = append(out, entry...)
		if prev != nil {
			table = append(table, append(append([]byte{}, prev...), entry[0]))
		}
		prev = entry
		switch len(table) + earlyChange {
		case 512:
			codeLen = 10
		case 1024:
			codeLen = 11
		case 2048:
			codeLen = 12
		}
		if len(table) >= 4096 {
			reset()
			codeLen = 9
			prev = nil
		}
	}
	return out
}

// streamFilters extracts the filter chain and its parameters from a stream
// dictionary, honouring the /F and /DP abbreviations. resolve follows indirect
// references; it may be nil for in-memory dictionaries.
func streamFilters(d *Dict, resolve func(Object) Object) ([]Name, []*Dict) {
	if resolve == nil {
		resolve = func(o Object) Object { return o }
	}
	get := func(primary, abbrev Name) Object {
		if v := d.GetRaw(primary); v != nil {
			return resolve(v)
		}
		if v := d.GetRaw(abbrev); v != nil {
			return resolve(v)
		}
		return nil
	}
	var names []Name
	switch f := get("Filter", "F").(type) {
	case Name:
		names = []Name{f}
	case Array:
		for _, e := range f {
			if n, ok := resolve(e).(Name); ok {
				names = append(names, n)
			}
		}
	}
	var parms []*Dict
	switch p := get("DecodeParms", "DP").(type) {
	case *Dict:
		parms = []*Dict{p}
	case Array:
		for _, e := range p {
			if pd, ok := resolve(e).(*Dict); ok {
				parms = append(parms, pd)
			} else {
				parms = append(parms, nil)
			}
		}
	}
	return names, parms
}

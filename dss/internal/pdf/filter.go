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
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"math"
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
// defect appends to warn. Every intermediate and final output is bounded by
// the default MaxStreamSize; exceeding it yields ErrLimitExceeded.
func Decode(raw []byte, filters []Name, parms []*Dict, warn *[]Warning) ([]byte, error) {
	return decodeLimited(raw, filters, parms, warn, defaultMaxStreamSize)
}

// defaultMaxStreamSize is Options.MaxStreamSize's default.
const defaultMaxStreamSize = 512 << 20

// decodeLimited is Decode with an explicit bound on the size of every decoded
// intermediate (DESIGN.md §2.7 resource guards): a few kilobytes of Flate, LZW
// or RunLength input, or a /DecodeParms /Columns of 2^40, must not be able to
// make us allocate gigabytes. pdfbox has no such bound; it has the JVM heap.
func decodeLimited(raw []byte, filters []Name, parms []*Dict, warn *[]Warning, limit int64) ([]byte, error) {
	if limit <= 0 {
		limit = defaultMaxStreamSize
	}
	out := raw
	for i, f := range filters {
		var parm *Dict
		if i < len(parms) {
			parm = parms[i]
		}
		var err error
		out, err = decodeOne(out, canonicalFilterName(f), parm, warn, limit)
		if err != nil {
			return nil, err
		}
		if int64(len(out)) > limit {
			return nil, errDecodedTooLarge(f, limit)
		}
	}
	if len(out) > 0 && len(raw) > 0 && &out[0] == &raw[0] {
		// No filter changed the bytes. raw may alias the document's own buffer
		// (the parser does not copy stream bodies), so hand out a copy.
		out = append([]byte(nil), out...)
	}
	return out, nil
}

func errDecodedTooLarge(f Name, limit int64) error {
	return fmt.Errorf("%w: /%s output exceeds %d bytes (MaxStreamSize, or a smaller bound for an xref stream)", ErrLimitExceeded, f, limit)
}

func decodeOne(raw []byte, f Name, parm *Dict, warn *[]Warning, limit int64) ([]byte, error) {
	switch f {
	case FilterFlate:
		out, err := flateDecodeLimited(raw, warn, limit)
		if err != nil {
			return nil, err
		}
		return applyPredictorFromParms(out, parm, warn, limit)
	case FilterLZW:
		early := 1
		if parm != nil {
			if v, ok := parm.GetRaw("EarlyChange").(Integer); ok {
				early = int(v)
			}
		}
		out, err := lzwDecodeLimited(raw, early, warn, limit)
		if err != nil {
			return nil, err
		}
		return applyPredictorFromParms(out, parm, warn, limit)
	case FilterASCIIHex:
		return asciiHexDecode(raw), nil
	case FilterASCII85:
		return ascii85Decode(raw), nil
	case FilterRunLength:
		return runLengthDecodeLimited(raw, limit)
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

func applyPredictorFromParms(data []byte, parm *Dict, warn *[]Warning, limit int64) ([]byte, error) {
	if parm == nil {
		return data, nil
	}
	predictor := intFromDict(parm, "Predictor", 1)
	if predictor <= 1 {
		return data, nil
	}
	colors := intFromDict(parm, "Colors", 1)
	bpc := intFromDict(parm, "BitsPerComponent", 8)
	columns := intFromDict(parm, "Columns", 1)
	return applyPredictor(data, predictor, colors, bpc, columns, warn, limit)
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
	out, _ := flateDecodeLimited(raw, warn, math.MaxInt64)
	return out
}

// flateDecodeLimited is FlateDecode with a bound on the inflated size: more
// than limit bytes of output is ErrLimitExceeded, not a truncation.
func flateDecodeLimited(raw []byte, warn *[]Warning, limit int64) ([]byte, error) {
	if len(raw) < 2 {
		// Fewer than two bytes of input: empty output, no error.
		return []byte{}, nil
	}
	r := flate.NewReader(bytes.NewReader(raw[2:]))
	defer r.Close()
	buf := boundedBuffer{limit: limit}
	_, err := buf.ReadFrom(r)
	if err == errBoundExceeded {
		return nil, errDecodedTooLarge(FilterFlate, limit)
	}
	if err != nil {
		addWarning(warn, WarnFlateTruncated, -1,
			"premature end of flate stream: "+err.Error())
	}
	return buf.Bytes(), nil
}

// errBoundExceeded is boundedBuffer's refusal to hold more than its limit.
var errBoundExceeded = errors.New("pdf: bounded buffer full")

// boundedChunk is the largest single allocation boundedBuffer makes while
// filling; past it, output accumulates in chunks of this size.
const boundedChunk = 1 << 20

// boundedBuffer collects a decoder's output and refuses to hold more than
// limit bytes. A bytes.Buffer (or append) doubles its backing array, so output
// that ends up just over the limit briefly holds about twice the limit, all of
// it cleared memory: a 4 MB Flate bomb under the 512 MiB default peaked at
// 2 GiB and took seconds before failing. Filling fixed-size chunks keeps the
// peak at the limit plus one chunk; only a result over one chunk is copied
// once, into an exactly sized slice, by Bytes.
type boundedBuffer struct {
	limit int64
	n     int64    // bytes held
	full  [][]byte // filled chunks, each boundedChunk long
	cur   []byte   // the chunk being filled
}

// room makes the current chunk non-full and returns its free capacity.
func (b *boundedBuffer) room() int {
	if len(b.cur) == cap(b.cur) {
		switch {
		case cap(b.cur) == 0:
			b.cur = make([]byte, 0, 512)
		case cap(b.cur) < boundedChunk:
			// Small outputs stay in one doubling slice, as with bytes.Buffer.
			b.cur = append(make([]byte, 0, min(2*cap(b.cur), boundedChunk)), b.cur...)
		default:
			b.full = append(b.full, b.cur)
			b.cur = make([]byte, 0, boundedChunk)
		}
	}
	return cap(b.cur) - len(b.cur)
}

// ReadFrom reads r to EOF. It stops with errBoundExceeded once more than limit
// bytes have been read (it reads at most one byte past the limit), and returns
// any other read error with the bytes read before it kept.
func (b *boundedBuffer) ReadFrom(r io.Reader) (int64, error) {
	var total int64
	for {
		room := b.room()
		if left := b.limit - b.n; int64(room) > left {
			room = int(left) + 1
		}
		n, err := r.Read(b.cur[len(b.cur) : len(b.cur)+room])
		b.cur = b.cur[:len(b.cur)+n]
		b.n += int64(n)
		total += int64(n)
		if b.n > b.limit {
			return total, errBoundExceeded
		}
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}
	}
}

// Write appends p, or nothing and errBoundExceeded when that would take the
// buffer past its limit.
func (b *boundedBuffer) Write(p []byte) (int, error) {
	if int64(len(p)) > b.limit-b.n {
		return 0, errBoundExceeded
	}
	written := len(p)
	for len(p) > 0 {
		room := b.room() // may replace b.cur: call it before slicing b.cur
		k := copy(b.cur[len(b.cur):len(b.cur)+min(room, len(p))], p)
		b.cur = b.cur[:len(b.cur)+k]
		p = p[k:]
	}
	b.n += int64(written)
	return written, nil
}

// Len is the number of bytes held.
func (b *boundedBuffer) Len() int64 { return b.n }

// Bytes returns everything written, nil when nothing was.
func (b *boundedBuffer) Bytes() []byte {
	if len(b.full) == 0 {
		return b.cur
	}
	out := make([]byte, 0, b.n)
	for _, c := range b.full {
		out = append(out, c...)
	}
	return append(out, b.cur...)
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
	out, err := applyPredictor(data, predictor, colors, bpc, columns, warn, defaultMaxStreamSize)
	if err != nil {
		addWarning(warn, WarnPredictorTruncated, -1, err.Error())
		return data
	}
	return out
}

// predictorRowLen is ceil(colors*bpc*columns/8), or -1 when the product
// overflows int64 or the row is longer than limit.
func predictorRowLen(colors, bpc, columns int, limit int64) int64 {
	bits := int64(colors)
	for _, f := range []int64{int64(bpc), int64(columns)} {
		if bits > math.MaxInt64/f {
			return -1
		}
		bits *= f
	}
	rowLen := bits/8 + boolInt64(bits%8 != 0)
	if rowLen > limit {
		return -1
	}
	return rowLen
}

func boolInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func applyPredictor(data []byte, predictor, colors, bpc, columns int, warn *[]Warning, limit int64) ([]byte, error) {
	if predictor <= 1 {
		return data, nil
	}
	if colors <= 0 || bpc <= 0 || columns <= 0 {
		return data, nil
	}
	// A hostile /Columns (or /Colors, /BitsPerComponent) would otherwise
	// overflow the row length or allocate a row of terabytes before a single
	// input byte is looked at.
	rl := predictorRowLen(colors, bpc, columns, limit)
	if rl < 0 {
		return nil, fmt.Errorf("%w: predictor row of %d x %d x %d bits exceeds MaxStreamSize (%d bytes)",
			ErrLimitExceeded, colors, bpc, columns, limit)
	}
	rowLen := int(rl)
	if rowLen <= 0 {
		return data, nil
	}
	if predictor == 2 {
		return tiffPredictor(data, colors, bpc, columns, rowLen), nil
	}
	bpp := (colors*bpc + 7) / 8
	// Sized up front: appending row by row regrew the output past its final
	// size, holding up to twice it, on top of the input.
	rows := (len(data) + rowLen) / (rowLen + 1)
	out := make([]byte, 0, rows*rowLen)
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
	return out, nil
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
	out, _ := runLengthDecodeLimited(raw, math.MaxInt64)
	return out
}

// runLengthDecodeLimited is runLengthDecode with a bound on the output: each
// two-byte run can expand to 128 bytes.
func runLengthDecodeLimited(raw []byte, limit int64) ([]byte, error) {
	out := boundedBuffer{limit: limit}
	var run [128]byte
	for i := 0; i < len(raw); {
		l := int(raw[i])
		i++
		if l == 128 {
			break
		}
		var err error
		if l < 128 {
			n := l + 1
			if i+n > len(raw) {
				n = len(raw) - i
			}
			_, err = out.Write(raw[i : i+n])
			i += n
		} else {
			if i >= len(raw) {
				break
			}
			b := raw[i]
			i++
			n := 257 - l
			for j := range n {
				run[j] = b
			}
			_, err = out.Write(run[:n])
		}
		if err != nil {
			return nil, errDecodedTooLarge(FilterRunLength, limit)
		}
	}
	return out.Bytes(), nil
}

// lzwDecode implements LZWDecode (LZWFilter). earlyChange defaults to 1.
func lzwDecode(raw []byte, earlyChange int, warn *[]Warning) []byte {
	out, _ := lzwDecodeLimited(raw, earlyChange, warn, math.MaxInt64)
	return out
}

// lzwDecodeLimited is lzwDecode with a bound on the output: a code can expand
// to a 4 KiB table entry, so 12 bits of input can yield thousands of bytes.
func lzwDecodeLimited(raw []byte, earlyChange int, warn *[]Warning, limit int64) ([]byte, error) {
	if earlyChange != 0 {
		earlyChange = 1
	}
	out := boundedBuffer{limit: limit}
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
			return out.Bytes(), nil
		}
		if _, err := out.Write(entry); err != nil {
			return nil, errDecodedTooLarge(FilterLZW, limit)
		}
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
	return out.Bytes(), nil
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

// PDF serialization primitives. See DESIGN.md §4.5 and the byte-determinism
// rules R1..R11 of §3.2.
//
// Provenance: org.apache.pdfbox.pdfwriter.COSWriter and COSStandardOutputStream
// of pdfbox 3.0.7 (visitFromDictionary / visitFromArray / visitFromStream /
// writeString / doWriteObject), plus COSName.writePDF (PDFBOX-2073's restricted
// name charset) and COSFloat.formatString (Float.toString, then
// BigDecimal.stripTrailingZeros().toPlainString() when that used an exponent).
//
// Matching pdfbox's formatting is free — the bytes land inside the signed
// /ByteRange and matching makes divergence visible in a diff — but it is not the
// contract; DESIGN.md §3.6 is. What *is* the contract: determinism. Nothing here
// consults a map iteration order, the clock, or a PRNG.

package pdf

import (
	"io"
	"math"
	"strconv"
	"strings"
)

// Writer emits PDF syntax under the byte-determinism rules R1..R14 of DESIGN.md.
// It tracks line state so WriteEOL is suppressed at a line start (R2).
type Writer struct {
	w   io.Writer
	pos int64
	// onNewLine mirrors COSStandardOutputStream.onNewLine: it is set only by
	// WriteEOL and cleared by every other write, exactly as pdfbox's
	// write(byte[]) does (writeCRLF therefore does *not* set it).
	onNewLine bool
	err       error

	// watch, when non-nil, records the /Contents and /ByteRange spans of one
	// specific signature dictionary as it is serialized (R17, R18). This
	// replaces pdfbox's reachedSignature heuristic in detectPossibleSignature:
	// we know which dictionary is ours, so we compare identity instead of
	// guessing from /Type and /ByteRange[2].
	watch *sigWatch

	// encrypt, when non-nil, is applied to every string and stream payload
	// written from here on (R20). It is set per indirect object by the Updater
	// because the encryption key is derived from the object key. skipEncrypt
	// names dictionary keys of the *current* dictionary whose string value must
	// be left in the clear — in practice the signature /Contents (§2.6).
	encrypt     func(data []byte, isString bool) ([]byte, error)
	skipEncrypt map[Name]bool
}

// sigWatch records where the signature placeholders landed in the output.
type sigWatch struct {
	dict *Dict

	contentsOffset int64 // absolute offset of the '<'
	contentsLength int64 // bytes including '<' and '>'
	contentsSeen   bool

	byteRangeOffset int64 // absolute offset of the '[' plus one (R18)
	byteRangeLength int64 // 35 for the pinned placeholder
	byteRangeSeen   bool
}

// NewWriter returns a Writer positioned at offset 0.
func NewWriter(w io.Writer) *Writer { return &Writer{w: w} }

// newWriterAt returns a Writer whose Pos is offset by base. The incremental
// serializer uses it so that every offset it computes is absolute in
// original||increment without the original ever being re-written (§3.1).
func newWriterAt(w io.Writer, base int64) *Writer { return &Writer{w: w, pos: base} }

// Pos reports the current absolute byte offset.
func (w *Writer) Pos() int64 { return w.pos }

// Err reports the first error seen. Every write is a no-op after it is set.
func (w *Writer) Err() error { return w.err }

// WriteRaw writes b verbatim.
func (w *Writer) WriteRaw(b []byte) error {
	if w.err != nil {
		return w.err
	}
	n, err := w.w.Write(b)
	w.pos += int64(n)
	if n > 0 {
		w.onNewLine = false
	}
	if err != nil {
		w.err = err
	}
	return w.err
}

func (w *Writer) writeString(s string) error { return w.WriteRaw([]byte(s)) }

func (w *Writer) writeByte(b byte) error { return w.WriteRaw([]byte{b}) }

// WriteEOL writes the end-of-line byte \n unless the last byte written was
// already an EOL written by WriteEOL (R1, R2). Without this suppression every
// dictionary gains a blank line and no golden matches.
func (w *Writer) WriteEOL() error {
	if w.err != nil {
		return w.err
	}
	if w.onNewLine {
		return nil
	}
	if err := w.writeByte('\n'); err != nil {
		return err
	}
	w.onNewLine = true
	return nil
}

// writeCRLF writes the two-byte sequence required around stream data and after
// each xref entry (R1). Like pdfbox's writeCRLF it leaves onNewLine false.
func (w *Writer) writeCRLF() error { return w.WriteRaw([]byte{'\r', '\n'}) }

// WriteObject writes o in PDF syntax, dispatching on its kind. Composite values
// are written direct: an indirect value is a Ref, never a *Dict.
func (w *Writer) WriteObject(o Object) error {
	if w.err != nil {
		return w.err
	}
	switch v := o.(type) {
	case nil:
		return w.writeString("null")
	case Null:
		return w.writeString("null")
	case Bool:
		if v {
			return w.writeString("true")
		}
		return w.writeString("false")
	case Integer:
		return w.writeString(strconv.FormatInt(int64(v), 10)) // R7
	case Real:
		if v.Raw != "" { // R8: echo what we parsed
			return w.writeString(v.Raw)
		}
		return w.writeString(FormatReal(v.Val))
	case String:
		return w.writeStringObject(v, false)
	case Name:
		return w.WriteRaw(EncodeName(v)) // R5
	case Array:
		return w.writeArray(v)
	case *Dict:
		return w.writeDict(v)
	case *Stream:
		return w.writeStream(v)
	case Ref:
		return w.writeRef(v) // R9
	default:
		return w.writeString("null")
	}
}

func (w *Writer) writeStringObject(s String, skipEncryption bool) error {
	b := s.Bytes
	if w.encrypt != nil && !skipEncryption {
		enc, err := w.encrypt(b, true)
		if err != nil {
			w.err = err
			return err
		}
		b = enc
	}
	return w.WriteRaw(EncodeString(String{Bytes: b, Hex: s.Hex})) // R6
}

func (w *Writer) writeRef(r Ref) error {
	if err := w.writeString(strconv.FormatInt(r.Num, 10)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	if err := w.writeString(strconv.FormatUint(uint64(r.Gen), 10)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	return w.writeByte('R')
}

// writeArray implements R4: items separated by a space, except that the
// separator after every tenth item is an EOL (COSWriter.visitFromArray's
// count % 10 == 0). Note that pdfbox only writes a separator *between* items,
// so a 10-element array ends "…item10]" with no trailing EOL of its own.
func (w *Writer) writeArray(a Array) error {
	if err := w.writeByte('['); err != nil {
		return err
	}
	for i, item := range a {
		if err := w.WriteObject(item); err != nil {
			return err
		}
		if i == len(a)-1 {
			break
		}
		if (i+1)%10 == 0 {
			if err := w.WriteEOL(); err != nil {
				return err
			}
		} else if err := w.writeByte(' '); err != nil {
			return err
		}
	}
	if err := w.writeByte(']'); err != nil {
		return err
	}
	return w.WriteEOL()
}

// writeDict implements R3. Entries whose value is nil are skipped entirely,
// matching visitFromDictionary's "dangling reference that points to nothing"
// branch.
func (w *Writer) writeDict(d *Dict) error {
	if d == nil {
		return w.writeString("null")
	}
	if err := w.writeString("<<"); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	watched := w.watch != nil && w.watch.dict == d
	for _, k := range d.Keys() {
		v := d.GetRaw(k)
		if v == nil {
			continue
		}
		if err := w.WriteRaw(EncodeName(k)); err != nil {
			return err
		}
		if err := w.writeByte(' '); err != nil {
			return err
		}
		if err := w.writeDictValue(d, k, v, watched); err != nil {
			return err
		}
		if err := w.WriteEOL(); err != nil {
			return err
		}
	}
	if err := w.writeString(">>"); err != nil {
		return err
	}
	return w.WriteEOL()
}

// writeDictValue writes one dictionary value, recording the placeholder spans
// when this is the watched signature dictionary. The offsets are exactly
// pdfbox's: signatureOffset is the position of the '<' and byteRangeOffset is
// the position of the '[' plus one.
func (w *Writer) writeDictValue(d *Dict, k Name, v Object, watched bool) error {
	skipEnc := w.skipEncrypt != nil && w.skipEncrypt[k]
	if watched && k == "Contents" {
		start := w.pos
		s, ok := v.(String)
		if !ok {
			return w.WriteObject(v)
		}
		// The signature /Contents is never encrypted, on read or on write
		// (§2.6, R20) — pdfbox guards this explicitly and getting it wrong
		// silently corrupts every signature in an encrypted document.
		if err := w.writeStringObject(s, true); err != nil {
			return err
		}
		w.watch.contentsOffset = start
		w.watch.contentsLength = w.pos - start
		w.watch.contentsSeen = true
		return nil
	}
	if watched && k == "ByteRange" {
		start := w.pos + 1
		if err := w.WriteObject(v); err != nil {
			return err
		}
		// pdfbox: byteRangeLength = pos - 1 - byteRangeOffset, evaluated after
		// visitFromArray has also written its trailing EOL. That yields 35 for
		// the reserved placeholder: 34 content bytes plus the ']'.
		w.watch.byteRangeOffset = start
		w.watch.byteRangeLength = w.pos - 1 - start
		w.watch.byteRangeSeen = true
		return nil
	}
	if s, ok := v.(String); ok {
		return w.writeStringObject(s, skipEnc)
	}
	return w.WriteObject(v)
}

// writeStream implements R11. /Length is rewritten as a direct integer holding
// the raw byte count — never an indirect reference — because the increment must
// be readable without resolving anything we did not write.
func (w *Writer) writeStream(s *Stream) error {
	raw := s.Raw
	if w.encrypt != nil && !streamNeverEncrypted(s) {
		enc, err := w.encrypt(raw, false)
		if err != nil {
			w.err = err
			return err
		}
		raw = enc
	}
	d := s.Dict.Clone()
	if d == nil {
		d = NewDict()
	}
	d.Set("Length", Integer(len(raw)))
	if err := w.writeDict(d); err != nil {
		return err
	}
	if err := w.writeString("stream"); err != nil {
		return err
	}
	if err := w.writeCRLF(); err != nil {
		return err
	}
	if len(raw) > 0 {
		if err := w.WriteRaw(raw); err != nil {
			return err
		}
	}
	if err := w.writeCRLF(); err != nil {
		return err
	}
	if err := w.writeString("endstream"); err != nil {
		return err
	}
	return w.WriteEOL()
}

// streamNeverEncrypted reports the §2.6 exemptions that apply to writing: an
// xref stream is never encrypted, nor is a stream carrying /Crypt /Identity.
func streamNeverEncrypted(s *Stream) bool {
	if s == nil || s.Dict == nil {
		return false
	}
	if t, ok := s.Dict.GetRaw("Type").(Name); ok && t == "XRef" {
		return true
	}
	switch f := s.Dict.GetRaw("Filter").(type) {
	case Name:
		return f == filterCryptName
	case Array:
		for _, e := range f {
			if n, ok := e.(Name); ok && n == filterCryptName {
				return true
			}
		}
	}
	return false
}

// filterCryptName is the /Crypt filter name. It is spelled out here rather than
// referenced from filter.go's exported FilterCrypt so that writer.go carries no
// dependency on a reader-owned file beyond the API DESIGN.md §4 pins.
const filterCryptName Name = "Crypt"

// WriteIndirect writes "N G obj … endobj" (R10).
func (w *Writer) WriteIndirect(k ObjectKey, o Object) error {
	if err := w.writeString(strconv.FormatInt(k.Num, 10)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	if err := w.writeString(strconv.FormatUint(uint64(k.Gen), 10)); err != nil {
		return err
	}
	if err := w.writeByte(' '); err != nil {
		return err
	}
	if err := w.writeString("obj"); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	if err := w.WriteObject(o); err != nil {
		return err
	}
	if err := w.WriteEOL(); err != nil {
		return err
	}
	if err := w.writeString("endobj"); err != nil {
		return err
	}
	return w.WriteEOL()
}

// EncodeName renders a name under R5: the byte is emitted verbatim only if it
// is in [A-Za-z0-9+\-_@*$;.], else as '#' plus two uppercase hex digits. This
// is deliberately stricter than ISO 32000-1 (PDFBOX-2073) — '!', ',', '~' and
// '\” are escaped even though the spec allows them literally.
func EncodeName(n Name) []byte {
	out := make([]byte, 0, len(n)+1)
	out = append(out, '/')
	for i := 0; i < len(n); i++ {
		c := n[i]
		switch {
		case c >= 'A' && c <= 'Z',
			c >= 'a' && c <= 'z',
			c >= '0' && c <= '9',
			c == '+', c == '-', c == '_', c == '@', c == '*', c == '$', c == ';', c == '.':
			out = append(out, c)
		default:
			out = append(out, '#', hexUpper[c>>4], hexUpper[c&0x0f])
		}
	}
	return out
}

const hexUpper = "0123456789ABCDEF"

// EncodeString renders a string under R6: hex form when any byte is >= 0x80 or
// is CR or LF, or when Hex is set; otherwise literal form escaping only '(',
// ')' and '\'. No octal escape is ever emitted.
func EncodeString(s String) []byte {
	hex := s.Hex
	if !hex {
		for _, b := range s.Bytes {
			if b >= 0x80 || b == 0x0d || b == 0x0a {
				hex = true
				break
			}
		}
	}
	if hex {
		out := make([]byte, 0, 2*len(s.Bytes)+2)
		out = append(out, '<')
		for _, b := range s.Bytes {
			out = append(out, hexUpper[b>>4], hexUpper[b&0x0f])
		}
		return append(out, '>')
	}
	out := make([]byte, 0, len(s.Bytes)+2)
	out = append(out, '(')
	for _, b := range s.Bytes {
		if b == '(' || b == ')' || b == '\\' {
			out = append(out, '\\')
		}
		out = append(out, b)
	}
	return append(out, ')')
}

// FormatReal renders v the way pdfbox writes a COSFloat (R8): Java's
// Float.toString, and — when that produced exponent notation — the
// BigDecimal(s).stripTrailingZeros().toPlainString() expansion of it. NaN and
// both infinities render as "0.0" rather than producing a syntactically invalid
// PDF number.
//
// Note that reals we merely echo never reach this function: Real.Raw is written
// verbatim, which is what makes round-tripping /Rect [0.0 0.0 595.276 841.89]
// byte-exact.
func FormatReal(v float64) string {
	f := float32(v)
	d := float64(f)
	if math.IsNaN(d) || math.IsInf(d, 0) {
		return "0.0"
	}
	s := javaFloatToString(f)
	if !strings.ContainsRune(s, 'E') {
		return s
	}
	return plainFromScientific(s)
}

// javaFloatToString reproduces java.lang.Float.toString: the shortest decimal
// that round-trips as a float32, in plain form when the decimal point position
// is in (-3, 8) and in "d.dddEnn" form otherwise, always with at least one
// digit on each side of the point.
func javaFloatToString(f float32) string {
	neg := math.Signbit(float64(f))
	sign := ""
	if neg {
		sign = "-"
	}
	if f == 0 {
		return sign + "0.0"
	}
	// 'e' with precision -1 gives the shortest round-tripping float32 digits.
	e := strconv.FormatFloat(float64(f), 'e', -1, 32)
	e = strings.TrimPrefix(e, "-")
	mant, expPart, _ := strings.Cut(e, "e")
	exp, err := strconv.Atoi(expPart)
	if err != nil {
		return sign + mant
	}
	digits := strings.Replace(mant, ".", "", 1)
	// pointPos is Java's decExponent: value = 0.<digits> * 10^pointPos.
	pointPos := exp + 1
	if pointPos > -3 && pointPos < 8 {
		return sign + plainDigits(digits, pointPos)
	}
	frac := digits[1:]
	if frac == "" {
		frac = "0"
	}
	return sign + digits[:1] + "." + frac + "E" + strconv.Itoa(exp)
}

// plainDigits places the decimal point pointPos digits into digits, padding so
// that at least one digit stands on each side of it.
func plainDigits(digits string, pointPos int) string {
	switch {
	case pointPos <= 0:
		return "0." + strings.Repeat("0", -pointPos) + digits
	case pointPos >= len(digits):
		return digits + strings.Repeat("0", pointPos-len(digits)) + ".0"
	default:
		return digits[:pointPos] + "." + digits[pointPos:]
	}
}

// plainFromScientific is BigDecimal(s).stripTrailingZeros().toPlainString() for
// the "d.dddEnn" strings javaFloatToString produces.
func plainFromScientific(s string) string {
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}
	mant, expPart, _ := strings.Cut(s, "E")
	exp, err := strconv.Atoi(expPart)
	if err != nil {
		return sign + mant
	}
	digits := strings.Replace(mant, ".", "", 1)
	// value = digits * 10^scale
	scale := exp - (len(digits) - 1)
	for len(digits) > 1 && digits[len(digits)-1] == '0' {
		digits = digits[:len(digits)-1]
		scale++
	}
	if digits == "0" {
		return sign + "0"
	}
	if scale >= 0 {
		return sign + digits + strings.Repeat("0", scale)
	}
	k := -scale
	if k < len(digits) {
		return sign + digits[:len(digits)-k] + "." + digits[len(digits)-k:]
	}
	return sign + "0." + strings.Repeat("0", k-len(digits)) + digits
}

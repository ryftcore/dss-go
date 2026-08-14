// Object and indirect-object parsing, including stream bodies.
// See DESIGN.md §2.1 (one-token lookahead for `N G R`) and §2.7 O1–O3, S1–S5.
//
// Provenance: org.apache.pdfbox.pdfparser.BaseParser (parseDirObject,
// parseCOSDictionary, parseCOSArray, skipWhiteSpaces) and COSParser
// (parseFileObject, parseCOSStream, readUntilEndStream, validateStreamLength)
// of pdfbox 3.0.7, plus EndstreamFilterStream for the recovered-length rules.

package pdf

import (
	"errors"
	"fmt"
)

// parser turns bytes into objects. It carries no document state beyond the
// resolver it needs for an indirect /Length.
type parser struct {
	lex *lexer
	// resolveLength resolves an indirect /Length. It is nil while the xref is
	// still being bootstrapped, which is exactly when pdfbox's getLength returns
	// null and the fallback endstream scan kicks in.
	resolveLength func(Ref) (int64, bool)
	warn          *[]Warning
	maxDepth      int
	maxStreamSize int64
	depth         int
	err           error // set for hard errors (resource guards) only
}

func newParser(data []byte, warn *[]Warning, maxDepth int, maxStreamSize int64) *parser {
	return &parser{
		lex:           newLexer(data, warn),
		warn:          warn,
		maxDepth:      maxDepth,
		maxStreamSize: maxStreamSize,
	}
}

func (p *parser) addWarning(code WarningCode, off int64, msg string) {
	addWarning(p.warn, code, off, msg)
}

// parseDirObject parses one direct object at the cursor, resolving the
// `N G R` three-token form into a Ref. It returns nil at EOF.
func (p *parser) parseDirObject() Object {
	if p.depth >= p.maxDepth {
		p.err = fmt.Errorf("%w: object nesting deeper than %d", ErrLimitExceeded, p.maxDepth)
		return Null{}
	}
	p.depth++
	defer func() { p.depth-- }()

	t := p.lex.next()
	return p.objectFromToken(t)
}

func (p *parser) objectFromToken(t token) Object {
	switch t.kind {
	case tokEOF:
		return nil
	case tokInteger:
		return p.maybeRef(t)
	case tokReal:
		return Real{Val: t.f, Raw: t.raw}
	case tokName:
		return t.name
	case tokString:
		return t.str
	case tokArrayOpen:
		return p.parseArrayBody()
	case tokDictOpen:
		return p.parseDictBody()
	case tokArrayClose, tokDictClose, tokBraceOpen, tokBraceClose:
		// Unbalanced close: pdfbox returns null and lets the caller recover.
		return nil
	case tokKeyword:
		switch t.raw {
		case "true":
			return Bool(true)
		case "false":
			return Bool(false)
		case "null":
			return Null{}
		case "R":
			// A bare R with no preceding integers; pdfbox builds an empty COSObject.
			return Null{}
		default:
			// endobj/endstream/stream/xref/trailer: push back, the caller wants it.
			p.lex.seek(t.pos)
			return nil
		}
	default: // tokJunk
		p.addWarning(WarnLexer, t.pos, "skipped unexpected object '"+t.raw+"'")
		return Null{}
	}
}

// maybeRef implements the one-token-beyond-the-second-integer lookahead: `6 0 R`
// is a reference, `6 0` followed by anything else is two integers.
func (p *parser) maybeRef(first token) Object {
	save := p.lex.pos
	t2 := p.lex.next()
	if t2.kind == tokInteger && t2.i >= 0 && t2.i <= 65535 {
		save2 := p.lex.pos
		t3 := p.lex.next()
		if t3.kind == tokKeyword && t3.raw == "R" {
			return Ref{Num: first.i, Gen: uint16(t2.i)}
		}
		p.lex.pos = save2
		_ = save2
	}
	p.lex.pos = save
	return Integer(first.i)
}

func (p *parser) parseArrayBody() Array {
	if p.depth >= p.maxDepth {
		p.err = fmt.Errorf("%w: array nesting deeper than %d", ErrLimitExceeded, p.maxDepth)
		return nil
	}
	p.depth++
	defer func() { p.depth-- }()

	arr := Array{}
	for {
		if p.err != nil {
			return arr
		}
		t := p.lex.next()
		switch t.kind {
		case tokEOF:
			return arr
		case tokArrayClose:
			return arr
		case tokDictClose:
			// pdfbox's parseCOSArray tolerates a stray '>>' by stopping.
			p.addWarning(WarnLexer, t.pos, "'>>' inside array; array closed early")
			return arr
		case tokKeyword:
			switch t.raw {
			case "true":
				arr = append(arr, Bool(true))
			case "false":
				arr = append(arr, Bool(false))
			case "null":
				arr = append(arr, Null{})
			case "R":
				// stray R, ignore
			default:
				// endobj/stream/…: the array is unterminated.
				p.lex.seek(t.pos)
				p.addWarning(WarnLexer, t.pos, "unterminated array before '"+t.raw+"'")
				return arr
			}
		default:
			o := p.objectFromToken(t)
			if o == nil {
				return arr
			}
			arr = append(arr, o)
		}
	}
}

// parseDictBody parses entries after '<<' has been consumed.
func (p *parser) parseDictBody() *Dict {
	if p.depth >= p.maxDepth {
		p.err = fmt.Errorf("%w: dictionary nesting deeper than %d", ErrLimitExceeded, p.maxDepth)
		return NewDict()
	}
	p.depth++
	defer func() { p.depth-- }()

	d := NewDict()
	for {
		if p.err != nil {
			return d
		}
		t := p.lex.next()
		switch t.kind {
		case tokEOF:
			return d
		case tokDictClose:
			return d
		case tokName:
			v := p.parseDirObject()
			if v == nil {
				// value missing (EOF or a keyword such as `endobj`): pdfbox stores
				// nothing and stops.
				return d
			}
			d.Set(t.name, v)
		case tokKeyword:
			if t.raw == "endobj" || t.raw == "stream" || t.raw == "endstream" ||
				t.raw == "trailer" || t.raw == "startxref" {
				p.lex.seek(t.pos)
				p.addWarning(WarnLexer, t.pos, "unterminated dictionary before '"+t.raw+"'")
				return d
			}
		default:
			// A value where a key was expected: skip it (BaseParser logs
			// "Invalid dictionary, found: ...").
			p.addWarning(WarnLexer, t.pos, "non-name key in dictionary; entry skipped")
		}
	}
}

// skipStreamWhitespace is BaseParser.skipWhiteSpaces: spaces, then CR[LF] or LF;
// a lone CR is accepted with a warning (S4).
func (p *parser) skipStreamWhitespace() {
	l := p.lex
	for l.pos < len(l.data) && l.data[l.pos] == ' ' {
		l.pos++
	}
	if l.pos >= len(l.data) {
		return
	}
	switch l.data[l.pos] {
	case '\r':
		l.pos++
		if l.pos < len(l.data) && l.data[l.pos] == '\n' {
			l.pos++
		} else {
			p.addWarning(WarnStreamEndFixed, int64(l.pos), "'stream' followed by a lone CR")
		}
	case '\n':
		l.pos++
	default:
		p.addWarning(WarnStreamEndFixed, int64(l.pos), "'stream' not followed by an EOL")
	}
}

// parseStreamBody parses a stream whose dictionary has already been read and
// whose `stream` keyword is at the cursor.
func (p *parser) parseStreamBody(dict *Dict) *Stream {
	t := p.lex.next() // 'stream'
	if t.kind != tokKeyword || t.raw != "stream" {
		p.lex.seek(t.pos)
		return nil
	}
	p.skipStreamWhitespace()
	start := int64(p.lex.pos)

	declared, haveDeclared := p.streamLength(dict)
	var length int64
	if haveDeclared && p.validStreamLength(start, declared) {
		length = declared
		p.lex.seek(start + length)
	} else {
		length = p.scanForEndstream(start)
		if !haveDeclared || declared != length {
			// S1: the dictionary's /Length is rewritten in memory.
			dict.Set("Length", Integer(length))
			p.addWarning(WarnStreamLengthFixed, start,
				fmt.Sprintf("stream length recovered as %d", length))
		}
		p.lex.seek(start + length)
	}

	if p.maxStreamSize > 0 && length > p.maxStreamSize {
		p.err = fmt.Errorf("%w: stream of %d bytes exceeds MaxStreamSize", ErrLimitExceeded, length)
		return nil
	}

	// endstream, with S2/S3 tolerances.
	save := p.lex.pos
	p.lex.skipSpaces()
	kw := p.readKeywordRun()
	switch {
	case kw == "endstream":
		// ok
	case kw == "endobj":
		// S2: rewind so the caller sees the endobj.
		p.addWarning(WarnStreamEndFixed, int64(save), "stream ends with 'endobj'")
		p.lex.pos -= len("endobj")
	case len(kw) > 9 && kw[:9] == "endstream":
		// S3: extra bytes glued to the keyword; rewind by the excess.
		p.addWarning(WarnStreamEndFixed, int64(save), "stream ends with '"+kw+"'")
		p.lex.pos -= len(kw) - 9
	default:
		p.addWarning(WarnStreamEndFixed, int64(save), "expected 'endstream', found '"+kw+"'")
		p.lex.pos = save
	}

	raw := make([]byte, length)
	copy(raw, p.lex.data[start:min64(start+length, int64(len(p.lex.data)))])
	return &Stream{Dict: dict, Raw: raw, Offset: start, Length: length}
}

// readKeywordRun reads a run of regular characters at the cursor.
func (p *parser) readKeywordRun() string {
	l := p.lex
	start := l.pos
	for l.pos < len(l.data) && isRegular(l.data[l.pos]) {
		l.pos++
	}
	return string(l.data[start:l.pos])
}

// streamLength reads /Length, resolving an indirect reference when the xref is
// already usable. The bool is false when no usable integer was found.
func (p *parser) streamLength(dict *Dict) (int64, bool) {
	switch v := dict.GetRaw("Length").(type) {
	case Integer:
		return int64(v), true
	case Ref:
		if p.resolveLength == nil {
			return 0, false
		}
		return p.resolveLength(v)
	}
	return 0, false
}

// validStreamLength is COSParser.validateStreamLength: a length of 0 or less, or
// one that runs past EOF, or one that does not land on `endstream`, is rejected.
func (p *parser) validStreamLength(start, length int64) bool {
	if length <= 0 {
		return false // S5 (and PDFBOX-5880)
	}
	end := start + length
	if end > int64(len(p.lex.data)) {
		return false
	}
	save := p.lex.pos
	p.lex.seek(end)
	p.lex.skipSpaces()
	ok := hasPrefixAt(p.lex.data, p.lex.pos, "endstream")
	p.lex.pos = save
	return ok
}

func hasPrefixAt(data []byte, pos int, s string) bool {
	if pos < 0 || pos+len(s) > len(data) {
		return false
	}
	return string(data[pos:pos+len(s)]) == s
}

// scanForEndstream reproduces COSParser.readUntilEndStream + EndstreamFilterStream:
// scan for `endstream` (falling back to `endobj` when the first three characters
// matched `end`), then drop a trailing CRLF or LF — but keep a lone CR, and keep
// the trailing EOL entirely when the first ten bytes look like ASCII text
// (PDFBOX-2120).
func (p *parser) scanForEndstream(start int64) int64 {
	data := p.lex.data
	i := int(start)
	end := -1
	for i < len(data) {
		if data[i] == 'e' {
			if hasPrefixAt(data, i, "endstream") {
				end = i
				break
			}
			if hasPrefixAt(data, i, "endobj") {
				end = i
				break
			}
		}
		i++
	}
	if end < 0 {
		end = len(data)
	}
	content := data[start:end]
	return int64(len(trimStreamTrailer(content)))
}

// trimStreamTrailer applies EndstreamFilterStream's rules to the raw run of bytes
// that precedes the endstream keyword.
func trimStreamTrailer(b []byte) []byte {
	if len(b) > 10 && !looksBinary(b[:10]) {
		// PDFBOX-2120: ASCII content keeps its final CR LF / LF.
		return b
	}
	if n := len(b); n > 0 {
		if b[n-1] == '\n' {
			b = b[:n-1]
			if n := len(b); n > 0 && b[n-1] == '\r' {
				b = b[:n-1]
			}
		}
		// a lone CR is kept (calculateLength writes it back)
	}
	return b
}

// looksBinary is EndstreamFilterStream's heuristic, with Java's signed bytes
// spelled out: a byte >= 0x80 counts as binary.
func looksBinary(b []byte) bool {
	for _, c := range b {
		if c >= 0x80 || c < 0x09 || (c > 0x0a && c < 0x20 && c != 0x0d) {
			return true
		}
	}
	return false
}

// parseIndirectAt parses `N G obj … endobj` at off. want is the key the xref
// claims lives there; a generation mismatch corrects the key (O1), an object
// number mismatch is reported by returning ok=false.
func (p *parser) parseIndirectAt(off int64, want ObjectKey) (Object, ObjectKey, bool) {
	if off < 0 || off >= int64(len(p.lex.data)) {
		return nil, want, false
	}
	p.lex.seek(off)
	p.lex.skipSpaces()
	t1 := p.lex.next()
	if t1.kind != tokInteger {
		return nil, want, false
	}
	t2 := p.lex.next()
	if t2.kind != tokInteger {
		return nil, want, false
	}
	t3 := p.lex.next()
	if t3.kind != tokKeyword || t3.raw != "obj" {
		return nil, want, false
	}
	got := ObjectKey{Num: t1.i, Gen: uint16(t2.i)}
	if !want.IsZero() || want.Num != 0 {
		if got.Num != want.Num {
			p.addWarning(WarnObjectHeaderFixed, off,
				fmt.Sprintf("object header %s does not match xref entry %s", got, want))
			return nil, want, false
		}
		if got.Gen != want.Gen {
			p.addWarning(WarnObjectHeaderFixed, off,
				fmt.Sprintf("object header %s has a different generation than xref entry %s", got, want))
		}
	}

	obj := p.parseDirObject()
	if p.err != nil {
		return nil, got, false
	}
	if d, ok := obj.(*Dict); ok {
		save := p.lex.pos
		p.lex.skipSpaces()
		if hasPrefixAt(p.lex.data, p.lex.pos, "stream") {
			p.lex.pos = save
			s := p.parseStreamBody(d)
			if p.err != nil {
				return nil, got, false
			}
			if s != nil {
				s.key = got
				obj = s
			}
		} else {
			p.lex.pos = save
		}
	}

	// endobj is optional (O2).
	save := p.lex.pos
	p.lex.skipSpaces()
	if hasPrefixAt(p.lex.data, p.lex.pos, "endobj") {
		p.lex.pos += len("endobj")
	} else {
		p.lex.pos = save
	}
	if obj == nil {
		obj = Null{}
	}
	return obj, got, true
}

// parseTrailerDictAt parses the dictionary following the `trailer` keyword.
func (p *parser) parseTrailerDictAt(off int64) (*Dict, error) {
	p.lex.seek(off)
	p.lex.skipSpaces()
	t := p.lex.next()
	if t.kind != tokDictOpen {
		return nil, errors.New("pdf: trailer is not a dictionary")
	}
	d := p.parseDictBody()
	if p.err != nil {
		return nil, p.err
	}
	return d, nil
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

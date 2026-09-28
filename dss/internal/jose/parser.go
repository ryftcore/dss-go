// Ported from org.jose4j.json.internal.json_simple.parser.JSONParser and its JFlex lexer Yylex,
// together with org.jose4j.json.JsonUtil.parseJson and its DupeKeyDisallowingLinkedHashMap
// container factory (jose4j 0.9.6).
package jose

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// ParseError is the counterpart of json_simple's ParseException as JsonUtil re-wraps it in a
// JoseException. Callers in dss-jades only ever check "did this parse", never why, so the
// message is informative rather than contractual.
type ParseError struct {
	// Position is the index, in bytes, at which the offending token starts.
	Position int
	// Message describes what went wrong.
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("jose: parsing error at position %d: %s", e.Position, e.Message)
}

// ErrNotJSONObject reports a document whose root is not a JSON object. JsonUtil.parseJson casts
// the parse result to Map and lets the resulting ClassCastException become
// JoseException("Expecting a JSON object at the root but ..."), so "[1,2]" and "\"x\"" are
// failures there even though they are perfectly good JSON.
var ErrNotJSONObject = errors.New("jose: expecting a JSON object at the root")

// ParseJSON parses a JSON object, insertion-ordered, rejecting duplicate member names. Port of
// JsonUtil.parseJson(String) - the entry point Headers.setFullHeaderAsJsonString,
// DSSJsonUtils.parseJsonStringToMap and JWSJsonSerializationParser all use.
//
// Two behaviours come from the container factory rather than from the grammar and both matter:
// objects are LinkedHashMaps, so a header parsed here and written back out keeps its original
// member order (and therefore its signature); and a repeated member name is an error, not a
// last-one-wins overwrite.
func ParseJSON(s string) (*Object, error) {
	value, err := parseJSONValue(s, false)
	if err != nil {
		return nil, err
	}
	object, ok := value.(*Object)
	if !ok {
		return nil, ErrNotJSONObject
	}
	return object, nil
}

// MaxJSONDepth is how deeply objects and arrays may nest in a document ParseJSON or
// ParseJSONAny accepts. The parser itself keeps its stack on the heap, as json_simple's does,
// but materialize and the writer recurse once per level, and a JAdES header is attacker
// input: 20,000,000 '[' (40 MB) overflowed the goroutine stack in materialize, a fatal error
// no recover() can catch. A real JOSE header nests a handful of levels.
//
// DIVERGENCE, deliberate: json_simple's JSONParser.parse (driven by JsonUtil.parseJson) has
// no nesting limit, so jose4j parses such a document and only a later recursive consumer -
// JSONValue.toJSONString, say - dies of StackOverflowError. Here the parse fails with an
// ordinary *ParseError instead.
const MaxJSONDepth = 1000

// ParseJSONAny parses any JSON value - object, array, string, number, boolean or null. Port of
// `new JSONParser().parse(String)` with no container factory, which DSSJsonUtils.parseJsonString
// and parseBase64UrlEncoded use.
//
// Without the factory json-simple builds plain JSONObjects, which extend java.util.HashMap, so
// the objects returned here iterate in HashMap bucket order and duplicate names are allowed
// (last one wins). That is a real difference from ParseJSON and it is reproduced rather than
// smoothed over, since it decides the member order if such a value is ever re-serialized.
func ParseJSONAny(s string) (any, error) {
	return parseJSONValue(s, true)
}

// parseJSONValue drives json_simple's JSONParser state machine. hashOrdered selects the
// container factory: false is JsonUtil's dupe-disallowing LinkedHashMap, true is the default
// JSONObject.
func parseJSONValue(s string, hashOrdered bool) (any, error) {
	lex := &jsonLexer{input: s}

	newObject := func() *Object {
		if hashOrdered {
			return NewHashObject()
		}
		return NewObject()
	}
	// putMember is where the two container factories differ.
	putMember := func(object *Object, key string, value any) error {
		if !hashOrdered && object.ContainsKey(key) {
			return &ParseError{
				Position: lex.pos,
				Message:  fmt.Sprintf("an entry for '%s' already exists, names must be unique", key),
			}
		}
		object.Put(key, value)
		return nil
	}

	const (
		sInit = iota
		sInFinishedValue
		sInObject
		sInArray
		sPassedPairKey
	)

	var statusStack []int
	var valueStack []any
	status := sInit
	depth := 0 // open objects and arrays

	push := func(st int, v any) {
		statusStack = append(statusStack, st)
		valueStack = append(valueStack, v)
	}
	// open is called for every '{' and '['; see MaxJSONDepth.
	open := func(tok jsonToken) error {
		depth++
		if depth > MaxJSONDepth {
			return &ParseError{Position: tok.pos, Message: fmt.Sprintf("objects and arrays nested more than %d deep", MaxJSONDepth)}
		}
		return nil
	}
	popStatus := func() {
		statusStack = statusStack[:len(statusStack)-1]
	}
	popValue := func() any {
		v := valueStack[len(valueStack)-1]
		valueStack = valueStack[:len(valueStack)-1]
		return v
	}
	peekValue := func() any { return valueStack[len(valueStack)-1] }
	peekStatus := func() int {
		if len(statusStack) == 0 {
			return -1
		}
		return statusStack[len(statusStack)-1]
	}

	unexpected := func(tok jsonToken) error {
		return &ParseError{Position: tok.pos, Message: fmt.Sprintf("unexpected token %s", tok)}
	}

	for {
		tok, err := lex.next()
		if err != nil {
			return nil, err
		}

		switch status {
		case sInit:
			switch tok.kind {
			case tokenValue:
				status = sInFinishedValue
				push(status, tok.value)
			case tokenLeftBrace:
				if err := open(tok); err != nil {
					return nil, err
				}
				status = sInObject
				push(status, newObject())
			case tokenLeftSquare:
				if err := open(tok); err != nil {
					return nil, err
				}
				status = sInArray
				push(status, &jsonArray{})
			default:
				return nil, unexpected(tok)
			}

		case sInFinishedValue:
			if tok.kind == tokenEOF {
				return materialize(popValue()), nil
			}
			return nil, unexpected(tok)

		case sInObject:
			switch tok.kind {
			case tokenComma:
				// json-simple tolerates missing and stray commas alike.
			case tokenValue:
				key, ok := tok.value.(string)
				if !ok {
					return nil, unexpected(tok)
				}
				push(sPassedPairKey, key)
				status = sPassedPairKey
			case tokenRightBrace:
				depth--
				if len(valueStack) > 1 {
					popStatus()
					popValue()
					status = peekStatus()
				} else {
					status = sInFinishedValue
				}
			default:
				return nil, unexpected(tok)
			}

		case sPassedPairKey:
			switch tok.kind {
			case tokenColon:
				// separator, nothing to do
			case tokenValue:
				popStatus()
				key := popValue().(string)
				if err := putMember(peekValue().(*Object), key, tok.value); err != nil {
					return nil, err
				}
				status = peekStatus()
			case tokenLeftSquare:
				if err := open(tok); err != nil {
					return nil, err
				}
				popStatus()
				key := popValue().(string)
				array := &jsonArray{}
				if err := putMember(peekValue().(*Object), key, array); err != nil {
					return nil, err
				}
				status = sInArray
				push(status, array)
			case tokenLeftBrace:
				if err := open(tok); err != nil {
					return nil, err
				}
				popStatus()
				key := popValue().(string)
				object := newObject()
				if err := putMember(peekValue().(*Object), key, object); err != nil {
					return nil, err
				}
				status = sInObject
				push(status, object)
			default:
				return nil, unexpected(tok)
			}

		case sInArray:
			switch tok.kind {
			case tokenComma:
				// as above
			case tokenValue:
				array := peekValue().(*jsonArray)
				array.items = append(array.items, tok.value)
			case tokenRightSquare:
				depth--
				if len(valueStack) > 1 {
					popStatus()
					popValue()
					status = peekStatus()
				} else {
					status = sInFinishedValue
				}
			case tokenLeftBrace:
				if err := open(tok); err != nil {
					return nil, err
				}
				array := peekValue().(*jsonArray)
				object := newObject()
				array.items = append(array.items, object)
				status = sInObject
				push(status, object)
			case tokenLeftSquare:
				if err := open(tok); err != nil {
					return nil, err
				}
				array := peekValue().(*jsonArray)
				inner := &jsonArray{}
				array.items = append(array.items, inner)
				status = sInArray
				push(status, inner)
			default:
				return nil, unexpected(tok)
			}
		}

		if tok.kind == tokenEOF {
			return nil, unexpected(tok)
		}
	}
}

// jsonArray is a growable list used while parsing. It is a distinct type only so that the
// parser can append to a container that is already referenced by its parent; materialize turns
// it into the []any the rest of the package uses.
type jsonArray struct {
	items []any
}

// materialize replaces every *jsonArray in the tree with a plain []any, so that callers of
// ParseJSON never see a parser-internal type.
func materialize(v any) any {
	switch t := v.(type) {
	case *jsonArray:
		out := make([]any, len(t.items))
		for i, item := range t.items {
			out[i] = materialize(item)
		}
		return out
	case *Object:
		for _, k := range t.keys {
			t.values[k] = materialize(t.values[k])
		}
		return t
	default:
		return v
	}
}

// ---------------------------------------------------------------------------
// Lexer - port of Yylex
// ---------------------------------------------------------------------------

type jsonTokenKind int

const (
	tokenEOF jsonTokenKind = iota
	tokenValue
	tokenLeftBrace
	tokenRightBrace
	tokenLeftSquare
	tokenRightSquare
	tokenComma
	tokenColon
)

type jsonToken struct {
	kind  jsonTokenKind
	value any
	pos   int
}

func (t jsonToken) String() string {
	switch t.kind {
	case tokenEOF:
		return "end of input"
	case tokenValue:
		return fmt.Sprintf("value %v", t.value)
	case tokenLeftBrace:
		return "'{'"
	case tokenRightBrace:
		return "'}'"
	case tokenLeftSquare:
		return "'['"
	case tokenRightSquare:
		return "']'"
	case tokenComma:
		return "','"
	default:
		return "':'"
	}
}

type jsonLexer struct {
	input string
	pos   int
}

// next returns the next token. Port of Yylex.yylex().
//
// The lexer is looser than RFC 8259 in ways that are reproduced verbatim, because they decide
// whether a JAdES document parses at all:
//
//	{"a":01}       -> 1        (the integer pattern is -?[0-9]+, leading zeros and all)
//	{"a":"\x"}     -> "\x"     (an unrecognised escape contributes a literal backslash)
//	{"a":"\u00zz"} -> "\u00zz" (likewise, since \u needs four hex digits to match)
//
// and stricter than one might guess in others: "+1", ".5", "1.", "1e" and "-" are all errors.
func (l *jsonLexer) next() (jsonToken, error) {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return jsonToken{kind: tokenEOF, pos: l.pos}, nil
	}
	start := l.pos
	c := l.input[l.pos]
	switch c {
	case '{':
		l.pos++
		return jsonToken{kind: tokenLeftBrace, pos: start}, nil
	case '}':
		l.pos++
		return jsonToken{kind: tokenRightBrace, pos: start}, nil
	case '[':
		l.pos++
		return jsonToken{kind: tokenLeftSquare, pos: start}, nil
	case ']':
		l.pos++
		return jsonToken{kind: tokenRightSquare, pos: start}, nil
	case ',':
		l.pos++
		return jsonToken{kind: tokenComma, pos: start}, nil
	case ':':
		l.pos++
		return jsonToken{kind: tokenColon, pos: start}, nil
	case '"':
		s, err := l.lexString()
		if err != nil {
			return jsonToken{}, err
		}
		return jsonToken{kind: tokenValue, value: s, pos: start}, nil
	}
	if c == '-' || (c >= '0' && c <= '9') {
		return l.lexNumber()
	}
	if strings.HasPrefix(l.input[l.pos:], "true") {
		l.pos += 4
		return jsonToken{kind: tokenValue, value: true, pos: start}, nil
	}
	if strings.HasPrefix(l.input[l.pos:], "false") {
		l.pos += 5
		return jsonToken{kind: tokenValue, value: false, pos: start}, nil
	}
	if strings.HasPrefix(l.input[l.pos:], "null") {
		l.pos += 4
		return jsonToken{kind: tokenValue, value: nil, pos: start}, nil
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return jsonToken{}, &ParseError{Position: start, Message: fmt.Sprintf("unexpected character %q", r)}
}

// skipWhitespace consumes Yylex's WS class, [ \t\r\n]+.
func (l *jsonLexer) skipWhitespace() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\r', '\n':
			l.pos++
		default:
			return
		}
	}
}

// lexNumber implements the INT and DOUBLE patterns with JFlex's longest-match rule: INT is
// -?[0-9]+ and DOUBLE is INT with an optional .[0-9]+ fraction and an optional [eE][-+]?[0-9]+
// exponent, so a token that has neither is an INT (java.lang.Long, or BigInteger on overflow)
// and anything longer is a DOUBLE.
func (l *jsonLexer) lexNumber() (jsonToken, error) {
	start := l.pos
	p := l.pos
	if p < len(l.input) && l.input[p] == '-' {
		p++
	}
	digitsStart := p
	for p < len(l.input) && l.input[p] >= '0' && l.input[p] <= '9' {
		p++
	}
	if p == digitsStart {
		return jsonToken{}, &ParseError{Position: start, Message: "unexpected character '-'"}
	}
	intEnd := p

	// Optional fraction: only a '.' followed by at least one digit extends the match.
	if p < len(l.input) && l.input[p] == '.' {
		q := p + 1
		fracStart := q
		for q < len(l.input) && l.input[q] >= '0' && l.input[q] <= '9' {
			q++
		}
		if q > fracStart {
			p = q
		}
	}
	// Optional exponent, likewise all-or-nothing.
	if p < len(l.input) && (l.input[p] == 'e' || l.input[p] == 'E') {
		q := p + 1
		if q < len(l.input) && (l.input[q] == '+' || l.input[q] == '-') {
			q++
		}
		expStart := q
		for q < len(l.input) && l.input[q] >= '0' && l.input[q] <= '9' {
			q++
		}
		if q > expStart {
			p = q
		}
	}

	text := l.input[start:p]
	l.pos = p
	if p == intEnd {
		// INT: Long.valueOf, falling back to BigInteger on overflow.
		if v, err := strconv.ParseInt(text, 10, 64); err == nil {
			return jsonToken{kind: tokenValue, value: NewLong(v), pos: start}, nil
		}
		b, ok := new(big.Int).SetString(text, 10)
		if !ok {
			return jsonToken{}, &ParseError{Position: start, Message: fmt.Sprintf("malformed number %q", text)}
		}
		return jsonToken{kind: tokenValue, value: NewBigInteger(b), pos: start}, nil
	}
	// DOUBLE: Double.valueOf. Overflow becomes +/-Infinity in Java rather than an error, which
	// is what ParseFloat's ErrRange result already holds.
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		var numErr *strconv.NumError
		if !errors.As(err, &numErr) || !errors.Is(numErr.Err, strconv.ErrRange) {
			return jsonToken{}, &ParseError{Position: start, Message: fmt.Sprintf("malformed number %q", text)}
		}
	}
	return jsonToken{kind: tokenValue, value: NewDouble(f), pos: start}, nil
}

// lexString implements Yylex's STRING_BEGIN state.
func (l *jsonLexer) lexString() (string, error) {
	start := l.pos
	l.pos++ // opening quote
	var sb strings.Builder
	for {
		if l.pos >= len(l.input) {
			return "", &ParseError{Position: start, Message: "unterminated string"}
		}
		c := l.input[l.pos]
		if c == '"' {
			l.pos++
			return sb.String(), nil
		}
		if c != '\\' {
			r, size := utf8.DecodeRuneInString(l.input[l.pos:])
			sb.WriteRune(r)
			l.pos += size
			continue
		}
		// Escape sequence.
		if l.pos+1 >= len(l.input) {
			// A trailing backslash matches Yylex's bare '\\' rule and contributes itself.
			sb.WriteByte('\\')
			l.pos++
			continue
		}
		switch l.input[l.pos+1] {
		case '"':
			sb.WriteByte('"')
			l.pos += 2
		case '\\':
			sb.WriteByte('\\')
			l.pos += 2
		case '/':
			sb.WriteByte('/')
			l.pos += 2
		case 'b':
			sb.WriteByte('\b')
			l.pos += 2
		case 'f':
			sb.WriteByte('\f')
			l.pos += 2
		case 'n':
			sb.WriteByte('\n')
			l.pos += 2
		case 'r':
			sb.WriteByte('\r')
			l.pos += 2
		case 't':
			sb.WriteByte('\t')
			l.pos += 2
		case 'u':
			unit, ok := readHex4(l.input, l.pos+2)
			if !ok {
				// \u without four hex digits does not match the rule, so the bare backslash
				// rule fires instead and the "u..." is copied literally.
				sb.WriteByte('\\')
				l.pos++
				continue
			}
			l.pos += 6
			// Java appends the raw UTF-16 code unit; a surrogate pair only becomes one
			// character once its partner arrives, which is what utf16.DecodeRune models.
			if utf16.IsSurrogate(rune(unit)) {
				if low, ok2 := readEscapedUTF16(l.input, l.pos); ok2 {
					if r := utf16.DecodeRune(rune(unit), rune(low)); r != utf8.RuneError {
						sb.WriteRune(r)
						l.pos += 6
						continue
					}
				}
				// A lone surrogate cannot be represented in a Go string; U+FFFD is what a
				// Java String holding one turns into as soon as it is encoded to UTF-8,
				// which is the only form these bytes ever reach.
				sb.WriteRune(utf8.RuneError)
				continue
			}
			sb.WriteRune(rune(unit))
		default:
			// Any other escape: the bare backslash rule.
			sb.WriteByte('\\')
			l.pos++
		}
	}
}

// readHex4 reads exactly four hexadecimal digits at index i.
func readHex4(s string, i int) (uint16, bool) {
	if i+4 > len(s) {
		return 0, false
	}
	var v uint16
	for k := 0; k < 4; k++ {
		c := s[i+k]
		switch {
		case c >= '0' && c <= '9':
			v = v<<4 | uint16(c-'0')
		case c >= 'a' && c <= 'f':
			v = v<<4 | uint16(c-'a'+10)
		case c >= 'A' && c <= 'F':
			v = v<<4 | uint16(c-'A'+10)
		default:
			return 0, false
		}
	}
	return v, true
}

// readEscapedUTF16 reads a "\uXXXX" sequence at index i, used to pair up surrogates.
func readEscapedUTF16(s string, i int) (uint16, bool) {
	if i+6 > len(s) || s[i] != '\\' || s[i+1] != 'u' {
		return 0, false
	}
	return readHex4(s, i+2)
}

package xpath10

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokKind uint8

const (
	tokEOF tokKind = iota
	tokName
	tokLiteral     // '...' or "..."
	tokNumber      // a number literal: lexed so it can be refused by name, never evaluated
	tokSlash       // /
	tokDoubleSlash // //
	tokLBracket    // [
	tokRBracket    // ]
	tokLParen      // (
	tokRParen      // )
	tokAt          // @
	tokDot         // .
	tokDotDot      // ..
	tokStar        // * in node-test position
	tokColonColon  // ::
	tokColon       // :
	tokComma       // ,
	tokEq          // =
	tokOr          // the operator name "or"

	// tokOperator carries a token that is valid XPath 1.0 but outside this subset - an
	// arithmetic, relational or union operator, "and", "!=", or a variable reference - so the
	// parser can name it in an *UnsupportedError instead of calling it a syntax error.
	tokOperator
)

type token struct {
	kind tokKind
	text string
	pos  int
}

// lex splits an XPath expression into tokens.
//
// The one subtle rule is XPath 1.0 clause 3.7's disambiguation of "*" and of the operator
// names: after a token that can end an operand, "*" is a multiply operator and "and"/"or"/
// "div"/"mod" are operator names, while everywhere else "*" is a wildcard node test and a name
// is a name. The preceding-token test below is that rule verbatim - without it a step named
// "and" and the literal in "[local-name()='or']" would both lex wrongly.
func lex(expr string) ([]token, error) {
	var out []token
	i := 0
	for i < len(expr) {
		c := expr[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
			continue
		case c == '/':
			if i+1 < len(expr) && expr[i+1] == '/' {
				out = append(out, token{tokDoubleSlash, "//", i})
				i += 2
			} else {
				out = append(out, token{tokSlash, "/", i})
				i++
			}
		case c == '[':
			out = append(out, token{tokLBracket, "[", i})
			i++
		case c == ']':
			out = append(out, token{tokRBracket, "]", i})
			i++
		case c == '(':
			out = append(out, token{tokLParen, "(", i})
			i++
		case c == ')':
			out = append(out, token{tokRParen, ")", i})
			i++
		case c == '@':
			out = append(out, token{tokAt, "@", i})
			i++
		case c == ',':
			out = append(out, token{tokComma, ",", i})
			i++
		case c == ':':
			if i+1 < len(expr) && expr[i+1] == ':' {
				out = append(out, token{tokColonColon, "::", i})
				i += 2
			} else {
				out = append(out, token{tokColon, ":", i})
				i++
			}
		case c == '=':
			out = append(out, token{tokEq, "=", i})
			i++
		case c == '!':
			if i+1 < len(expr) && expr[i+1] == '=' {
				out = append(out, token{tokOperator, "!=", i})
				i += 2
			} else {
				return nil, &SyntaxError{expr, i, "'!' must be followed by '='"}
			}
		case c == '\'' || c == '"':
			end := strings.IndexByte(expr[i+1:], c)
			if end < 0 {
				return nil, &SyntaxError{expr, i, "unterminated string literal"}
			}
			out = append(out, token{tokLiteral, expr[i+1 : i+1+end], i})
			i += end + 2
		case c == '*':
			if precedingEndsOperand(out) {
				out = append(out, token{tokOperator, "*", i}) // multiply
			} else {
				out = append(out, token{tokStar, "*", i})
			}
			i++
		case c == '.':
			switch {
			case i+1 < len(expr) && expr[i+1] == '.':
				out = append(out, token{tokDotDot, "..", i})
				i += 2
			case i+1 < len(expr) && isDigit(expr[i+1]):
				j := i + 1
				for j < len(expr) && isDigit(expr[j]) {
					j++
				}
				out = append(out, token{tokNumber, expr[i:j], i})
				i = j
			default:
				out = append(out, token{tokDot, ".", i})
				i++
			}
		case isDigit(c):
			j := i
			for j < len(expr) && (isDigit(expr[j]) || expr[j] == '.') {
				j++
			}
			out = append(out, token{tokNumber, expr[i:j], i})
			i = j
		case c == '|' || c == '+' || c == '-' || c == '<' || c == '>' || c == '$':
			j := i + 1
			if (c == '<' || c == '>') && j < len(expr) && expr[j] == '=' {
				j++
			}
			out = append(out, token{tokOperator, expr[i:j], i})
			i = j
		default:
			name, n := lexName(expr[i:])
			if n == 0 {
				r, _ := utf8.DecodeRuneInString(expr[i:])
				return nil, &SyntaxError{expr, i, "unexpected character " + string(r)}
			}
			kind := tokName
			if precedingEndsOperand(out) {
				switch name {
				case "or":
					kind = tokOr
				case "and", "div", "mod":
					kind = tokOperator
				}
				// A name in operator position that is none of those is a syntax error in
				// XPath, but reporting it here would pre-empt the parser's better message,
				// so it stays a name.
			}
			out = append(out, token{kind, name, i})
			i += n
		}
	}
	out = append(out, token{tokEOF, "", len(expr)})
	return out, nil
}

// precedingEndsOperand reports whether the previous token can end an operand, which is XPath
// 1.0 clause 3.7's condition for reading "*" as a multiply operator and "and"/"or"/"div"/"mod"
// as operator names: true unless there is no preceding token, or it is one of "@", "::", "(",
// "[", "," or an operator. tokColon is in the list because "prefix:*" reaches the lexer as
// three tokens - "*" is not an NCName - so the "*" of a namespace wildcard follows a colon and
// must still lex as a wildcard, to be refused as a node test rather than as a stray operator.
func precedingEndsOperand(out []token) bool {
	if len(out) == 0 {
		return false
	}
	switch out[len(out)-1].kind {
	case tokAt, tokColonColon, tokColon, tokLParen, tokLBracket, tokComma,
		tokSlash, tokDoubleSlash, tokEq, tokOr, tokOperator:
		return false
	}
	return true
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// lexName reads a QName - "prefix:local" or a bare NCName - as a single token, which the
// parser splits on the colon. Joining the two halves here is what keeps "a:b" apart from the
// "::" of an axis specifier: the colon is only swallowed when a second colon does not follow,
// so "self::text()" stops at the axis and "ers:TimeStamp" does not.
//
// "prefix:*" is left as three tokens, since "*" is not an NCName.
func lexName(s string) (string, int) {
	i := readNCName(s, 0)
	if i == 0 {
		return "", 0
	}
	if i < len(s) && s[i] == ':' && i+1 < len(s) && s[i+1] != ':' {
		if j := readNCName(s, i+1); j > i+1 {
			i = j
		}
	}
	return s[:i], i
}

// readNCName returns the end offset of the NCName starting at from, or from when there is
// none.
func readNCName(s string, from int) int {
	i := from
	for i < len(s) {
		r, n := utf8.DecodeRuneInString(s[i:])
		if i == from {
			if !isNameStart(r) {
				return from
			}
		} else if !isNameChar(r) {
			break
		}
		i += n
	}
	return i
}

func isNameStart(r rune) bool { return r == '_' || unicode.IsLetter(r) }

func isNameChar(r rune) bool {
	return r == '_' || r == '-' || r == '.' ||
		unicode.IsLetter(r) || unicode.IsDigit(r) ||
		unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r)
}

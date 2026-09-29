// Ported from the output writers of
// org.apache.xml.security.c14n.implementations.CanonicalizerBase and UtfHelpper
// (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import (
	"bufio"
	"unicode/utf8"
)

// The escaping tables, per context. Uppercase hex digits, lowercase x, no leading zeroes.
//
//	character | attribute value | text/CDATA | comment | PI target | PI data
//	&         | &amp;           | &amp;      | -       | -         | -
//	<         | &lt;            | &lt;       | -       | -         | -
//	>         | -               | &gt;       | -       | -         | -
//	"         | &quot;          | -          | -       | -         | -
//	#x9       | &#x9;           | -          | -       | -         | -
//	#xA       | &#xA;           | -          | -       | -         | -
//	#xD       | &#xD;           | &#xD;      | &#xD;   | &#xD;     | &#xD;
//
// Everything else is written as UTF-8; "'" is never escaped and attribute values are always
// delimited by '"'.
const (
	escAmp  = "&amp;"
	escLT   = "&lt;"
	escGT   = "&gt;"
	escQuot = "&quot;"
	escX9   = "&#x9;"
	escXA   = "&#xA;"
	escXD   = "&#xD;"
)

// writeAttrValueEscaped ports CanonicalizerBase.outputAttrToWriter's value loop.
func writeAttrValueEscaped(w *bufio.Writer, s string) {
	for i := 0; i < len(s); {
		r, size := nextRune(s, i)
		i += size
		switch r {
		case '&':
			w.WriteString(escAmp)
		case '<':
			w.WriteString(escLT)
		case '"':
			w.WriteString(escQuot)
		case 0x09:
			w.WriteString(escX9)
		case 0x0A:
			w.WriteString(escXA)
		case 0x0D:
			w.WriteString(escXD)
		default:
			writeRune(w, r)
		}
	}
}

// writeTextEscaped ports CanonicalizerBase.outputTextToWriter. Text and CDATA nodes go through
// the same writer upstream, so a CDATA section is emitted as escaped text.
func writeTextEscaped(w *bufio.Writer, s string) {
	for i := 0; i < len(s); {
		r, size := nextRune(s, i)
		i += size
		switch r {
		case '&':
			w.WriteString(escAmp)
		case '<':
			w.WriteString(escLT)
		case '>':
			w.WriteString(escGT)
		case 0x0D:
			w.WriteString(escXD)
		default:
			writeRune(w, r)
		}
	}
}

// writeCarriageReturnEscaped ports the comment and processing-instruction loops of
// CanonicalizerBase.outputCommentToWriter and outputPItoWriter, which escape #xD and nothing
// else: comments and PIs recognize neither markup nor references, so "&#13;" inside a comment
// is six literal characters that pass through unchanged. Because XML 2.11 turns a literal CR
// in the source into LF, a CR can only reach this function from DOM construction.
func writeCarriageReturnEscaped(w *bufio.Writer, s string) {
	for i := 0; i < len(s); {
		r, size := nextRune(s, i)
		i += size
		if r == 0x0D {
			w.WriteString(escXD)
			continue
		}
		writeRune(w, r)
	}
}

// invalidRune marks a byte that is not part of a well-formed UTF-8 sequence. It is not a code
// point, so it can never collide with real content - in particular not with a well-formed
// U+FFFD, which must be written out as itself.
const invalidRune rune = -1

// writeRune ports UtfHelpper.writeCodePointToUtf8, which writes '?' for anything that is not a
// valid code point - in Java, an unpaired surrogate. Go strings cannot hold one after a
// successful parse, so this only fires on invalid UTF-8 reaching the tree through DOM
// construction; the substitution keeps that case byte-comparable with Java rather than
// emitting U+FFFD, which Java never emits.
func writeRune(w *bufio.Writer, r rune) {
	if r == invalidRune {
		w.WriteByte('?')
		return
	}
	if r < 0x80 {
		w.WriteByte(byte(r))
		return
	}
	w.WriteRune(r)
}

// nextRune decodes the rune at s[i], reporting an invalid byte as invalidRune so that writeRune
// substitutes '?' exactly where Java's UtfHelpper does. A plain range over the string would
// collapse an invalid byte and a well-formed U+FFFD into the same rune; keeping them apart
// costs one function and removes a silent divergence. It decodes in place: materialising a
// []rune per text node, attribute value, comment and PI is one heap allocation (4 bytes per
// character) on the hottest loop of the canonicalizer.
func nextRune(s string, i int) (rune, int) {
	r, size := utf8.DecodeRuneInString(s[i:])
	if r == utf8.RuneError && size == 1 {
		return invalidRune, 1
	}
	return r, size
}

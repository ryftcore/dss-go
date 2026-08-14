package xmldom

import "fmt"

// SyntaxError reports a malformed or rejected document. Line and Column are 1-based
// and count runes; Offset is the byte offset into the document after transcoding to
// UTF-8, which is the source itself for the usual UTF-8-without-BOM case.
type SyntaxError struct {
	Line, Column int
	Offset       int64
	Msg          string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("xmldom: %s (line %d, column %d)", e.Msg, e.Line, e.Column)
}

// lineCol converts a byte offset into a 1-based line and rune column.
func lineCol(buf []byte, off int64) (int, int) {
	if off < 0 {
		off = 0
	}
	if off > int64(len(buf)) {
		off = int64(len(buf))
	}
	line, col := 1, 1
	for i := int64(0); i < off; i++ {
		switch {
		case buf[i] == '\n':
			line++
			col = 1
		case buf[i]&0xC0 == 0x80:
			// UTF-8 continuation byte: part of the rune already counted.
		default:
			col++
		}
	}
	return line, col
}

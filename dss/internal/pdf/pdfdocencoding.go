package pdf

import "unicode/utf16"

// PDFDocEncoding, the single-byte encoding ISO 32000-1:2008 Annex D.2 defines for PDF text
// strings that carry no UTF-16 byte-order mark.
//
// It is NOT Latin-1: 0x18-0x1F and 0x80-0xA0 carry accents, dashes, quotation marks, ligatures
// and Central-European letters at code points that have nothing to do with Latin-1's C1 control
// block, and 0x7F/0x9F/0xAD are undefined. Approximating it with Latin-1 - as this package and
// pades/native_pdf_dict.go both used to - silently mangles every /Reason, /Location, /Name and
// signature-field name written by a producer using the Central-European half of the table: a
// Czech signer name "Martin Petrzela" (with a z-caron, byte 0x9E) came out as U+009E instead of
// U+017E against upstream's own pades-ocsp-archiveCutOff-invalid.pdf test fixture.
//
// codeToUnicode is PDFBox's own org.apache.pdfbox.cos.PDFDocEncoding CODE_TO_UNI table
// (pdfbox 3.0.7, the PDF backend upstream DSS's dss-pades-pdfbox module runs on, and therefore
// the reference this port cross-validates against), reproduced verbatim including its two
// quirks: the three undefined codes 0x7F and 0x9F map to U+FFFD REPLACEMENT CHARACTER, while
// 0xAD - skipped by PDFBox's initialisation loop and never assigned afterwards - is left at the
// zero value its int[] was allocated with, i.e. U+0000.
var codeToUnicode = func() [256]rune {
	var table [256]rune
	// Initialize with basically ISO-8859-1, skipping the entries with no Unicode column.
	for i := 0; i < 256; i++ {
		if i > 0x17 && i < 0x20 {
			continue
		}
		if i > 0x7E && i < 0xA1 {
			continue
		}
		if i == 0xAD {
			continue
		}
		table[i] = rune(i)
	}
	// Then all deviations, based on the table in ISO 32000-1:2008.
	deviations := map[int]rune{
		0x18: '˘', // BREVE
		0x19: 'ˇ', // CARON
		0x1A: 'ˆ', // MODIFIER LETTER CIRCUMFLEX ACCENT
		0x1B: '˙', // DOT ABOVE
		0x1C: '˝', // DOUBLE ACUTE ACCENT
		0x1D: '˛', // OGONEK
		0x1E: '˚', // RING ABOVE
		0x1F: '˜', // SMALL TILDE
		0x7F: '�', // undefined
		0x80: '•', // BULLET
		0x81: '†', // DAGGER
		0x82: '‡', // DOUBLE DAGGER
		0x83: '…', // HORIZONTAL ELLIPSIS
		0x84: '—', // EM DASH
		0x85: '–', // EN DASH
		0x86: 'ƒ', // LATIN SMALL LETTER SCRIPT F
		0x87: '⁄', // FRACTION SLASH (solidus)
		0x88: '‹', // SINGLE LEFT-POINTING ANGLE QUOTATION MARK
		0x89: '›', // SINGLE RIGHT-POINTING ANGLE QUOTATION MARK
		0x8A: '−', // MINUS SIGN
		0x8B: '‰', // PER MILLE SIGN
		0x8C: '„', // DOUBLE LOW-9 QUOTATION MARK (quotedblbase)
		0x8D: '“', // LEFT DOUBLE QUOTATION MARK (quotedblleft)
		0x8E: '”', // RIGHT DOUBLE QUOTATION MARK (quotedblright)
		0x8F: '‘', // LEFT SINGLE QUOTATION MARK (quoteleft)
		0x90: '’', // RIGHT SINGLE QUOTATION MARK (quoteright)
		0x91: '‚', // SINGLE LOW-9 QUOTATION MARK (quotesinglbase)
		0x92: '™', // TRADE MARK SIGN
		0x93: 'ﬁ', // LATIN SMALL LIGATURE FI
		0x94: 'ﬂ', // LATIN SMALL LIGATURE FL
		0x95: 'Ł', // LATIN CAPITAL LETTER L WITH STROKE
		0x96: 'Œ', // LATIN CAPITAL LIGATURE OE
		0x97: 'Š', // LATIN CAPITAL LETTER S WITH CARON
		0x98: 'Ÿ', // LATIN CAPITAL LETTER Y WITH DIAERESIS
		0x99: 'Ž', // LATIN CAPITAL LETTER Z WITH CARON
		0x9A: 'ı', // LATIN SMALL LETTER DOTLESS I
		0x9B: 'ł', // LATIN SMALL LETTER L WITH STROKE
		0x9C: 'œ', // LATIN SMALL LIGATURE OE
		0x9D: 'š', // LATIN SMALL LETTER S WITH CARON
		0x9E: 'ž', // LATIN SMALL LETTER Z WITH CARON
		0x9F: '�', // undefined
		0xA0: '€', // EURO SIGN
	}
	for code, unicode := range deviations {
		table[code] = unicode
	}
	return table
}()

// DecodeTextString renders a PDF text string as Go text, reproducing PDFBox's
// COSString#getString: a UTF-16BE byte-order mark selects UTF-16BE, a UTF-16LE one selects
// UTF-16LE (PDFBox accepts it although the PDF specification does not define it), and anything
// else is PDFDocEncoding - see codeToUnicode above for why that is not the same as Latin-1.
func DecodeTextString(b []byte) string {
	if len(b) >= 2 {
		if b[0] == 0xFE && b[1] == 0xFF {
			return decodeUTF16(b[2:], true)
		}
		if b[0] == 0xFF && b[1] == 0xFE {
			return decodeUTF16(b[2:], false)
		}
	}
	runes := make([]rune, len(b))
	for i, c := range b {
		runes[i] = codeToUnicode[c]
	}
	return string(runes)
}

// decodeUTF16 decodes b as UTF-16 in the requested byte order, pairing surrogates the way
// java.lang.String's UTF_16BE/UTF_16LE decoders do.
func decodeUTF16(b []byte, bigEndian bool) string {
	units := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		if bigEndian {
			units = append(units, uint16(b[i])<<8|uint16(b[i+1]))
		} else {
			units = append(units, uint16(b[i+1])<<8|uint16(b[i]))
		}
	}
	return string(utf16.Decode(units))
}

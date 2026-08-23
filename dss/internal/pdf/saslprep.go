// saslprep.go is the SASLprep profile (RFC 4013) StandardSecurityHandler
// applies to an /R 6 password before hashing it (PDFBOX-4155).
//
// Provenance: org.apache.pdfbox.pdmodel.encryption.SaslPrep (pdfbox 3.0.7),
// ported table for table - including the places where it departs from a
// by-the-book RFC 3454 implementation - because what has to match is the
// byte string pdfbox hashes, not the RFC. Only the query-string form
// (saslPrepQuery, unassigned code points allowed) is needed: it is the one
// prepareForDecryption calls; saslPrepStored is used by the encryption side
// this package does not implement (DESIGN.md §0.2). What cannot be made
// identical is the Unicode version behind the tables - x/text's against the
// JDK's - so a code point assigned in one and not the other can still be
// classified differently; nothing short of pinning both runtimes fixes that.

package pdf

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/bidi"
	"golang.org/x/text/unicode/norm"
)

// saslPrepQuery is SaslPrep.saslPrepQuery: the RFC 4013 canonical form of
// s for use as a query string. It returns ErrProhibitedPassword - the Go
// shape of the IllegalArgumentException pdfbox throws - when s contains a
// prohibited code point or breaks the bidirectional-text rules.
func saslPrepQuery(s string) (string, error) {
	// 1. Map. pdfbox walks UTF-16 code units here; every code point the two
	// mapping tables name lies in the BMP, so walking code points is the
	// same walk. The order matters: U+200B is both a non-ASCII space and a
	// mapped-to-nothing character, and pdfbox maps it to ' ' first, which
	// is then kept.
	mapped := make([]rune, 0, len(s))
	for _, ch := range s {
		if saslNonASCIISpace(ch) {
			ch = ' '
		}
		if saslMappedToNothing(ch) {
			continue
		}
		mapped = append(mapped, ch)
	}

	// 2. Normalize.
	normalized := norm.NFKC.String(string(mapped))

	var containsRandALCat, containsLCat, initialRandALCat bool
	for i := 0; i < len(normalized); {
		codepoint, size := utf8.DecodeRuneInString(normalized[i:])
		// 3. Prohibit.
		if saslProhibited(codepoint) {
			return "", fmt.Errorf("%w: prohibited character %U at position %d", ErrProhibitedPassword, codepoint, i)
		}
		// 4. Check bidi. Character.getDirectionality's RIGHT_TO_LEFT,
		// RIGHT_TO_LEFT_ARABIC and LEFT_TO_RIGHT are the Bidi_Class values
		// R, AL and L. Java answers DIRECTIONALITY_UNDEFINED for every
		// unassigned code point, whereas x/text's tables give one its
		// block's default class (R across the Hebrew and Arabic blocks, L
		// elsewhere); the query profile lets unassigned code points through,
		// so they have to count as neither RandALCat nor LCat here.
		var isRandALcat, isLCat bool
		if !unicode.Is(unicode.Cn, codepoint) {
			properties, _ := bidi.LookupRune(codepoint)
			class := properties.Class()
			isRandALcat = class == bidi.R || class == bidi.AL
			isLCat = class == bidi.L
		}
		containsRandALCat = containsRandALCat || isRandALcat
		containsLCat = containsLCat || isLCat
		initialRandALCat = initialRandALCat || (i == 0 && isRandALcat)
		i += size
		if initialRandALCat && i >= len(normalized) && !isRandALcat {
			return "", fmt.Errorf("%w: first character is RandALCat, but last character is not", ErrProhibitedPassword)
		}
	}
	if containsRandALCat && containsLCat {
		return "", fmt.Errorf("%w: contains both RandALCat characters and LCat characters", ErrProhibitedPassword)
	}
	return normalized, nil
}

// saslProhibited is SaslPrep.prohibited: RFC 4013 §2.3, as the union of the
// RFC 3454 tables below. pdfbox passes the code point to its two
// char-typed predicates through a (char) cast, which keeps only the low 16
// bits of a supplementary code point; that truncation is reproduced, since
// it decides what pdfbox rejects.
func saslProhibited(codepoint rune) bool {
	truncated := rune(uint16(codepoint))
	return saslNonASCIISpace(truncated) ||
		saslASCIIControl(truncated) ||
		saslNonASCIIControl(codepoint) ||
		saslPrivateUse(codepoint) ||
		saslNonCharacterCodePoint(codepoint) ||
		saslSurrogateCodePoint(codepoint) ||
		saslInappropriateForPlainText(codepoint) ||
		saslInappropriateForCanonical(codepoint) ||
		saslChangeDisplayProperties(codepoint) ||
		saslTagging(codepoint)
}

// saslTagging is RFC 3454 Appendix C.9.
func saslTagging(codepoint rune) bool {
	return codepoint == 0xE0001 ||
		0xE0020 <= codepoint && codepoint <= 0xE007F
}

// saslChangeDisplayProperties is RFC 3454 Appendix C.8.
func saslChangeDisplayProperties(codepoint rune) bool {
	switch codepoint {
	case 0x0340, 0x0341, 0x200E, 0x200F, 0x202A, 0x202B, 0x202C, 0x202D, 0x202E,
		0x206A, 0x206B, 0x206C, 0x206D, 0x206E, 0x206F:
		return true
	}
	return false
}

// saslInappropriateForCanonical is RFC 3454 Appendix C.7.
func saslInappropriateForCanonical(codepoint rune) bool {
	return 0x2FF0 <= codepoint && codepoint <= 0x2FFB
}

// saslInappropriateForPlainText is RFC 3454 Appendix C.6.
func saslInappropriateForPlainText(codepoint rune) bool {
	return 0xFFF9 <= codepoint && codepoint <= 0xFFFD
}

// saslSurrogateCodePoint is RFC 3454 Appendix C.5.
func saslSurrogateCodePoint(codepoint rune) bool {
	return 0xD800 <= codepoint && codepoint <= 0xDFFF
}

// saslNonCharacterCodePoint is RFC 3454 Appendix C.4: U+FDD0..U+FDEF and
// the last two code points of every plane.
func saslNonCharacterCodePoint(codepoint rune) bool {
	if 0xFDD0 <= codepoint && codepoint <= 0xFDEF {
		return true
	}
	return codepoint <= 0x10FFFF && codepoint&0xFFFE == 0xFFFE
}

// saslPrivateUse is RFC 3454 Appendix C.3.
func saslPrivateUse(codepoint rune) bool {
	return 0xE000 <= codepoint && codepoint <= 0xF8FF ||
		0xF0000 <= codepoint && codepoint <= 0xFFFFD ||
		0x100000 <= codepoint && codepoint <= 0x10FFFD
}

// saslNonASCIIControl is RFC 3454 Appendix C.2.2.
func saslNonASCIIControl(codepoint rune) bool {
	switch codepoint {
	case 0x06DD, 0x070F, 0x180E, 0x200C, 0x200D, 0x2028, 0x2029,
		0x2060, 0x2061, 0x2062, 0x2063, 0xFEFF:
		return true
	}
	return 0x0080 <= codepoint && codepoint <= 0x009F ||
		0x206A <= codepoint && codepoint <= 0x206F ||
		0xFFF9 <= codepoint && codepoint <= 0xFFFC ||
		0x1D173 <= codepoint && codepoint <= 0x1D17A
}

// saslASCIIControl is RFC 3454 Appendix C.2.1.
func saslASCIIControl(ch rune) bool {
	return ch <= 0x001F || ch == 0x007F
}

// saslNonASCIISpace is RFC 3454 Appendix C.1.2.
func saslNonASCIISpace(ch rune) bool {
	switch ch {
	case 0x00A0, 0x1680, 0x202F, 0x205F, 0x3000:
		return true
	}
	return 0x2000 <= ch && ch <= 0x200B
}

// saslMappedToNothing is RFC 3454 Appendix B.1, "commonly mapped to nothing".
func saslMappedToNothing(ch rune) bool {
	switch ch {
	case 0x00AD, 0x034F, 0x1806, 0x180B, 0x180C, 0x180D,
		0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF:
		return true
	}
	return 0xFE00 <= ch && ch <= 0xFE0F
}

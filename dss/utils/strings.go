// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// string-related methods, matching org.apache.commons.lang3.StringUtils /
// Strings.CS / Strings.CI / NumberUtils.isDigits semantics.

package utils

import (
	"strings"
	"unicode"
)

// IsStringEmpty checks if the string is empty.
// Ex. "nowina" = false; "" = true
//
// NOTE: a blank string (e.g. "   ") is not empty!
func IsStringEmpty(text string) bool {
	return text == ""
}

// IsStringNotEmpty checks if the string is not empty.
// Ex. "nowina" = true; "" = false
func IsStringNotEmpty(text string) bool {
	return !IsStringEmpty(text)
}

// AreAllStringsEmpty checks if all strings are empty.
func AreAllStringsEmpty(values ...string) bool {
	for _, v := range values {
		if IsStringNotEmpty(v) {
			return false
		}
	}
	return true
}

// IsAtLeastOneStringNotEmpty checks if at least one string is not empty.
func IsAtLeastOneStringNotEmpty(values ...string) bool {
	return !AreAllStringsEmpty(values...)
}

// isJavaWhitespace mirrors java.lang.Character.isWhitespace, which is what
// commons-lang3 StringUtils.isBlank relies on. It is Unicode-category
// based (Zs/Zl/Zp, i.e. SPACE_SEPARATOR/LINE_SEPARATOR/PARAGRAPH_SEPARATOR)
// plus a handful of ASCII control characters, but — unlike Go's
// unicode.IsSpace — it explicitly EXCLUDES the non-breaking space
// variants U+00A0, U+2007 and U+202F, which Java does not consider
// whitespace.
func isJavaWhitespace(r rune) bool {
	switch r {
	case '\t', '\n', 0x0B, '\f', '\r', 0x1C, 0x1D, 0x1E, 0x1F:
		return true
	case 0x00A0, 0x2007, 0x202F:
		return false
	}
	return unicode.Is(unicode.Zs, r) || unicode.Is(unicode.Zl, r) || unicode.Is(unicode.Zp, r)
}

// IsStringBlank checks if the string is blank.
// Ex. "nowina" = false; "   " = true
func IsStringBlank(text string) bool {
	for _, r := range text {
		if !isJavaWhitespace(r) {
			return false
		}
	}
	return true
}

// IsStringNotBlank checks if the string is not blank.
// Ex. "nowina" = true; "   " = false
func IsStringNotBlank(text string) bool {
	return !IsStringBlank(text)
}

// AreStringsEqual checks if the strings are equal (case-sensitive).
// Ex. "nowina" == "nowina" = true; "nowina" == "Nowina" = false
func AreStringsEqual(text1, text2 string) bool {
	return text1 == text2
}

// AreStringsEqualIgnoreCase checks if the strings are equal, ignoring case.
// Ex. "nowina" == "Nowina" = true; "water" == "fire" = false
func AreStringsEqualIgnoreCase(text1, text2 string) bool {
	return strings.EqualFold(text1, text2)
}

// IsStringDigits checks if the string contains only digits.
// Ex. "123" = true; "1a2b" = false; "" = false
//
// Mirrors commons-lang3 NumberUtils.isDigits: unicode digits (category Nd)
// are accepted, an empty string is not.
func IsStringDigits(text string) bool {
	if text == "" {
		return false
	}
	for _, r := range text {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Trim trims the string, removing all leading and trailing whitespace.
// Ex. "   123 " = "123"
//
// Mirrors commons-lang3 StringUtils.trim, which trims characters <= U+0020
// (Java String.trim semantics), not full Unicode whitespace.
func Trim(text string) string {
	return strings.TrimFunc(text, func(r rune) bool {
		return r <= ' '
	})
}

// JoinStrings joins the strings with the given separator.
// Ex. ["Nowina", "123"], "," = "Nowina,123"
func JoinStrings(values []string, separator string) string {
	return strings.Join(values, separator)
}

// SubstringAfter returns the substring after the specified separator.
// Ex. "aaaaa?bbb", "?" = "bbb"
//
// Mirrors commons-lang3 StringUtils.substringAfter:
//   - an empty text returns the text unchanged
//   - a separator not found returns ""
func SubstringAfter(text, after string) string {
	if text == "" {
		return text
	}
	pos := strings.Index(text, after)
	if pos == -1 {
		return ""
	}
	return text[pos+len(after):]
}

// EndsWithIgnoreCase checks if the string ends with the expected suffix,
// ignoring case.
// Ex. "hello", "LO" = true; "hello", "a" = false
func EndsWithIgnoreCase(text, expected string) bool {
	if len(expected) > len(text) {
		return false
	}
	return strings.EqualFold(text[len(text)-len(expected):], expected)
}

// LowerCase converts a string to its lower case representation.
// Ex. "Nowina" = "nowina"
func LowerCase(text string) string {
	return strings.ToLower(text)
}

// UpperCase converts a string to its upper case representation.
// Ex. "Nowina" = "NOWINA"
func UpperCase(text string) string {
	return strings.ToUpper(text)
}

// GetFileNameExtension returns the extension of the given filename.
// Ex. "file.xml" = "xml"; "document.pdf" = "pdf"; "noext" = ""
//
// Mirrors commons-io FilenameUtils.getExtension: the last '.' only counts
// as an extension separator if it occurs after the last path separator
// ('/' or '\'); otherwise there is no extension.
func GetFileNameExtension(filename string) string {
	lastDot := strings.LastIndexByte(filename, '.')
	lastSep := lastIndexAny(filename, "/\\")
	if lastSep > lastDot {
		return ""
	}
	if lastDot == -1 {
		return ""
	}
	return filename[lastDot+1:]
}

func lastIndexAny(s, chars string) int {
	return strings.LastIndexAny(s, chars)
}

// IsTrue checks if the Boolean value is set to true.
//
// NOTE: if nil, returns false!
func IsTrue(b *bool) bool {
	return b != nil && *b
}

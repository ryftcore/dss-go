// Ported from java.text.MessageFormat, as used by
// dss-i18n/.../i18n/I18nProvider.java (DSS 6.5.RC1).
//
// Only the subset of java.text.MessageFormat exercised by DSS message
// patterns is implemented: plain "{index}" argument placeholders (no
// FormatType/FormatStyle, e.g. "{0,number}", appear anywhere in
// dss-messages.properties) and the single-quote quoting/escaping rules
// that apply to the pattern text regardless of whether it contains any
// placeholders. That last part matters for parity: several upstream
// message patterns contain a single, undoubled apostrophe (e.g.
// "...OCSP Responder's certificate path!"). java.text.MessageFormat's
// quoting rule silently drops such a bare quote character and treats the
// remainder of the pattern as quoted (i.e. any subsequent braces would be
// interpreted as literal characters, not placeholders) - it does not
// throw, and it does not restore the apostrophe. This is a known upstream
// quirk (the apostrophe never reaches I18nProvider's callers), and the
// marshal-parity contract requires reproducing it rather than "fixing" it.
package i18n

import (
	"fmt"
	"strconv"
	"strings"
)

// messageFormat ports the relevant subset of
// java.text.MessageFormat#format(String, Object...) (as invoked via the
// static convenience method MessageFormat.format(pattern, args) used by
// I18nProvider).
func messageFormat(pattern string, args []interface{}) string {
	var out strings.Builder
	inQuote := false

	runes := []rune(pattern)
	i := 0
	for i < len(runes) {
		ch := runes[i]

		if ch == '\'' {
			// "''" is an escaped literal single quote, in or out of a
			// quoted section.
			if i+1 < len(runes) && runes[i+1] == '\'' {
				out.WriteRune('\'')
				i += 2
				continue
			}
			// A lone quote toggles quoting; the quote character itself
			// is consumed, not emitted (this is the source of the
			// apostrophe-swallowing quirk described above).
			inQuote = !inQuote
			i++
			continue
		}

		if inQuote {
			out.WriteRune(ch)
			i++
			continue
		}

		if ch == '{' {
			// Scan for the FormatElement's matching close brace.
			// java.text.MessageFormat requires the ArgumentIndex to be a
			// run of ASCII digits; a ',' introduces an optional
			// FormatType which DSS patterns never use (see the file
			// header), so any comma-suffixed content up to the matching
			// '}' is treated as an (unused) format spec and ignored.
			j := i + 1
			for j < len(runes) && runes[j] != '}' {
				j++
			}
			if j >= len(runes) {
				// Unterminated FormatElement: mirrors
				// java.text.MessageFormat's pattern-parse failure mode
				// by emitting the remainder of the pattern verbatim
				// rather than panicking on malformed (never-upstream)
				// input.
				out.WriteString(string(runes[i:]))
				break
			}

			inner := string(runes[i+1 : j])
			indexPart := inner
			if comma := strings.IndexByte(inner, ','); comma >= 0 {
				indexPart = inner[:comma]
			}

			if argIndex, err := strconv.Atoi(indexPart); err == nil && argIndex >= 0 {
				if argIndex < len(args) {
					out.WriteString(messageFormatArgString(args[argIndex]))
				} else {
					// No argument supplied for this index: upstream
					// MessageFormat leaves the placeholder unformatted.
					out.WriteString(string(runes[i : j+1]))
				}
			} else {
				// Not a valid ArgumentIndex: not a real FormatElement,
				// emit verbatim (never occurs in upstream patterns).
				out.WriteString(string(runes[i : j+1]))
			}

			i = j + 1
			continue
		}

		out.WriteRune(ch)
		i++
	}

	return out.String()
}

// messageFormatArgString renders a substitution argument. Java's
// MessageFormat falls back to Object#toString() for arguments without a
// registered Format (which, absent a FormatType in the pattern, is every
// argument DSS ever substitutes - always a String, produced either
// directly by callers or by I18nProvider's nested-MessageTag resolution).
// fmt.Sprint mirrors that default toString() fallback for the argument
// types this port passes through I18nProvider.
func messageFormatArgString(arg interface{}) string {
	if s, ok := arg.(string); ok {
		return s
	}
	return fmt.Sprint(arg)
}

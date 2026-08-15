// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/LinkHeaderParser.java (DSS 6.5.RC1).
//
// This class is used to parse the "Link" HTTP response header value, according to RFC 8288
// requirements. Fully self-contained (pure string parsing), so it carries no forward
// dependencies on the rest of this batch.
package qwac

import (
	"fmt"
	"strings"
)

// LinkHeaderParser parses a "Link" HTTP response header.
type LinkHeaderParser struct{}

// NewLinkHeaderParser is the default constructor.
func NewLinkHeaderParser() *LinkHeaderParser {
	return &LinkHeaderParser{}
}

// Parse parses a full Link header into 1..n LinkHeader entries. Port of parse(String).
//
// Returns an error when headerValueStr is empty or blank, or malformed (Java's thrown
// IllegalArgumentException).
func (p *LinkHeaderParser) Parse(headerValueStr string) ([]*LinkHeader, error) {
	if strings.TrimSpace(headerValueStr) == "" {
		return nil, fmt.Errorf("link header cannot be null or empty")
	}

	// One header can have multiple links, separated by a comma BUT commas inside quotes must
	// not split entries.
	linkValues := splitHeaderValues(headerValueStr)

	result := make([]*LinkHeader, 0, len(linkValues))
	for _, linkValue := range linkValues {
		header, err := parseSingleLink(strings.TrimSpace(linkValue))
		if err != nil {
			return nil, err
		}
		result = append(result, header)
	}
	return result, nil
}

// splitHeaderValues splits a header by commas while respecting quoted values.
func splitHeaderValues(headerValueStr string) []string {
	var values []string

	var current strings.Builder
	inQuotes := false

	for _, c := range headerValueStr {
		if c == '"' {
			inQuotes = !inQuotes
		}

		if c == ',' && !inQuotes {
			values = append(values, current.String())
			current.Reset()
		} else {
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		values = append(values, current.String())
	}

	return values
}

// parseSingleLink parses a single "<url>; param1=...; param2=..." entry.
//
// Returns an error when headerValueStr does not start with '<' or is missing the closing '>'
// (Java's thrown IllegalArgumentException).
func parseSingleLink(headerValueStr string) (*LinkHeader, error) {
	if !strings.HasPrefix(headerValueStr, "<") {
		return nil, fmt.Errorf("link entry must start with '<': %s", headerValueStr)
	}

	urlEnd := strings.IndexByte(headerValueStr, '>')
	if urlEnd == -1 {
		return nil, fmt.Errorf("missing '>' in Link entry: %s", headerValueStr)
	}

	url := strings.TrimSpace(headerValueStr[1:urlEnd])
	attributes := make(map[string]*string)

	// Remaining attributes after '>'.
	attributesPart := strings.TrimSpace(headerValueStr[urlEnd+1:])

	if attributesPart != "" {
		// Split by semicolons.
		for _, part := range strings.Split(attributesPart, ";") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			keyValue := strings.SplitN(part, "=", 2)
			if len(keyValue) == 2 {
				key := strings.TrimSpace(keyValue[0])
				value := stripQuotes(strings.TrimSpace(keyValue[1]))
				attributes[key] = &value
			} else {
				// attribute without "=", allowed by RFC 8288 (flags).
				attributes[part] = nil
			}
		}
	}

	header := &LinkHeader{}
	header.SetURL(url)
	header.SetAttributes(attributes)
	return header, nil
}

// stripQuotes removes surrounding quotes if present.
func stripQuotes(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 1 && strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return value[1 : len(value)-1]
	}
	return value
}

// Ported from dss-i18n/.../i18n/I18nProvider.java (DSS 6.5.RC1).
package i18n

import (
	_ "embed"
	"strings"
)

// dssMessagesProperties is a byte-identical copy of upstream's
// src/main/resources/dss-messages.properties.
//
//go:embed dss-messages.properties
var dssMessagesProperties string

// propertyBundle is a parsed .properties resource, mirroring the subset of
// java.util.ResourceBundle's key/value view that I18nProvider relies on.
type propertyBundle map[string]string

// defaultBundle is the single embedded message bundle. Upstream's
// dss-i18n main resources ship only dss-messages.properties (no
// dss-messages_xx.properties locale variants - see MessageTagTest and
// I18nProviderTest, which rely on java.util.ResourceBundle silently
// falling back to this same default for every locale), so this is the
// bundle every Provider resolves to, regardless of locale.
var defaultBundle = parseProperties(dssMessagesProperties)

// I18nProvider is the internalization provider. Ports I18nProvider.
type Provider struct {
	bundle propertyBundle
}

// NewI18nProvider instantiates a default I18nProvider. Ports the
// zero-argument I18nProvider() constructor, which in Java resolves
// java.util.Locale.getDefault(); Go has no equivalent implicit JVM
// default Locale, and - per the defaultBundle doc comment - the locale
// makes no difference to bundle resolution in this port, so this is
// exactly NewI18nProviderForLocale("").
func NewI18nProvider() *Provider {
	return NewI18nProviderForLocale("")
}

// NewI18nProviderForLocale returns an instance of Provider for the
// given locale (e.g. "fr", "fr_FR", "de"). Ports I18nProvider(Locale).
//
// The locale argument is accepted for API parity with the Java
// constructor but does not currently affect resolution: see the
// defaultBundle doc comment.
func NewI18nProviderForLocale(locale string) *Provider {
	return &Provider{bundle: defaultBundle}
}

// GetMessage extracts a message by its key. Ports
// I18nProvider#getMessage(MessageTag, Object...).
//
// Panics with the Java exception's message if messageTag is the empty
// MessageTag (""), the Go stand-in for Java's null messageTag (a
// requireNonNull-style check on a value type - see PORTING.md).
func (p *Provider) GetMessage(messageTag MessageTag, args ...any) string {
	if messageTag == "" {
		panic("messageTag cannot be null!")
	}

	patternString, ok := p.bundle[messageTag.Id()]
	if !ok {
		// In case a value for the message tag does not exist, Java logs a slf4j warning and
		// continues; slf4j has no Go equivalent and is not ported, but messageTag.Id() is
		// still returned.
		return messageTag.Id()
	}
	return messageFormat(patternString, p.translateArgs(args))
}

// translateArgs allows nested MessageTags. Ports the private
// getArgs(Object[]) helper.
func (p *Provider) translateArgs(args []any) []any {
	if args == nil {
		return nil
	}
	translated := make([]any, len(args))
	copy(translated, args)
	for i, a := range args {
		if nested, ok := a.(MessageTag); ok {
			translated[i] = p.GetMessage(nested)
		}
	}
	return translated
}

// parseProperties parses a java.util.Properties-format resource's
// "key=value" lines into a propertyBundle, following
// java.util.Properties#load(Reader)'s whitespace-around-separator rules:
// whitespace between the key and the '='/':' separator, and whitespace
// immediately after it, is not part of the key or value (some
// dss-messages.properties lines are written as "KEY = value"). Trailing
// value whitespace is preserved, matching Properties#load.
//
// Upstream's dss-messages.properties uses only plain ASCII lines,
// '#'-prefixed comments, and blank lines - no line continuations and no
// '\uXXXX' unicode escapes (verified against the embedded resource) - so
// only that subset of the format is implemented; a byte-exact port of
// the general escaping grammar is unneeded (and untestable against
// upstream, since none of those forms occur in the resource being
// ported).
func parseProperties(data string) propertyBundle {
	bundle := make(propertyBundle)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		key, value, ok := parsePropertyLine(line)
		if !ok {
			continue
		}
		bundle[key] = value
	}
	return bundle
}

// parsePropertyLine splits a single non-comment, non-blank properties
// line into its key and value, per the whitespace/separator rules
// described in parseProperties's doc comment.
func parsePropertyLine(line string) (key, value string, ok bool) {
	n := len(line)
	i := 0
	for i < n && isPropertyWhitespace(line[i]) {
		i++
	}
	keyStart := i
	for i < n && !isPropertyWhitespace(line[i]) && line[i] != '=' && line[i] != ':' {
		i++
	}
	key = line[keyStart:i]
	if key == "" {
		return "", "", false
	}

	for i < n && isPropertyWhitespace(line[i]) {
		i++
	}
	if i < n && (line[i] == '=' || line[i] == ':') {
		i++
		for i < n && isPropertyWhitespace(line[i]) {
			i++
		}
	}
	return key, line[i:], true
}

// isPropertyWhitespace reports whether c is whitespace per
// java.util.Properties#load's definition (space, tab, form feed).
func isPropertyWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\f'
}

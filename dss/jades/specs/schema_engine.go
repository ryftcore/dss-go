// schema_engine.go has no single upstream Java source to port from. Upstream's
// eu.europa.esig.json.JSONSchemaAbstractUtils (dss-json-common) delegates all JSON Schema
// parsing/loading/evaluation to the third-party com.github.erosb.jsonsKema library (an
// "everit"-style draft-07 validator); there is no Go equivalent on the approved dependency list
// (PORTING.md: stdlib-first, golang.org/x/... only). This file is new, sanctioned machinery
// analogous to internal/xpath10's evidence-bounded XPath subset: it implements exactly the JSON
// Schema draft-07 constructs the ETSI TS 119 182-1 (JAdES) schemas and their RFC 7515/7517/7519/
// 7797 dependency schemas under ./schema/ actually use -
// type/required/properties/additionalProperties/items/minItems/maxItems/minProperties/
// maxProperties/enum/pattern/oneOf/anyOf/allOf/not/$ref/definitions - and nothing more (no
// patternProperties, no tuple-style "items" arrays, no "if/then/else", no numeric range
// keywords: the corpus never exercises them).
//
// DEVIATION FROM UPSTREAM (documented per PORTING.md): error message text does NOT match
// jsonsKema's wording or structure. Only the presence/absence of a validation error - and, where
// this port's own tests assert it, a recognizable substring naming the offending property or
// keyword - is a contract; exact message parity with everit/jsonsKema is out of scope.
// "format" (date-time, uri, uri-reference) and "contentEncoding" (base64) are both draft-07
// *annotation* keywords, not validation assertions in the general spec, and this corpus's own
// values regularly hold base64url content under a "contentEncoding":"base64" declaration (e.g.
// 'x5t#S256', 'signature') - so, matching that non-enforcement, neither is checked here.
package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// schemaSet resolves and caches the parsed JSON documents (schemas) reachable from a root schema
// URI through the definitions map, and evaluates instances against them.
type schemaSet struct {
	files fs.FS
	// defs maps a schema URI (as used in $ref targets and as the map keys upstream's
	// getSchemaDefinitions() populates) to its file path under files.
	defs map[string]string
	// docs caches parsed root documents by URI.
	docs map[string]any
}

func newSchemaSet(files fs.FS, defs map[string]string) *schemaSet {
	return &schemaSet{files: files, defs: defs, docs: map[string]any{}}
}

// decodeJSON parses JSON text preserving numbers as json.Number, so integer/number type checks
// and enum equality comparisons are consistent between schema documents and validated instances.
func decodeJSON(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

// sharedSchemaDocs caches, process-wide, the decoded documents of the embedded schema tree
// (jadesSchemaRoot) by file path. Upstream's JSONSchemaAbstractUtils#getValidator caches the
// compiled schema across calls; this is the equivalent, so a validation does not re-read and
// re-decode the ~11 embedded documents every time. The cached values are never mutated: the
// validate* methods only read the decoded trees.
var sharedSchemaDocs sync.Map // map[string]any

// compiledPatterns caches compiled "pattern" keywords (a nil *regexp.Regexp for an invalid one).
// The patterns come from the embedded schemas, so the set is small and bounded.
var compiledPatterns sync.Map // map[string]*regexp.Regexp

// compilePattern returns the compiled form of a "pattern" keyword, or nil when it is invalid.
func compilePattern(pattern string) *regexp.Regexp {
	if cached, ok := compiledPatterns.Load(pattern); ok {
		return cached.(*regexp.Regexp)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		re = nil
	}
	compiledPatterns.Store(pattern, re)
	return re
}

// load fetches and parses (once, cached) the root document registered under uri.
func (s *schemaSet) load(uri string) (any, error) {
	if doc, ok := s.docs[uri]; ok {
		return doc, nil
	}
	relPath, ok := s.defs[uri]
	if !ok {
		return nil, fmt.Errorf("unable to load a schema for URI: %s", uri)
	}
	// Only the embedded tree is shared across calls: its content is fixed, whereas any other
	// fs.FS may change (and is only used by tests).
	shared := s.files == jadesSchemaRoot
	if shared {
		if doc, ok := sharedSchemaDocs.Load(relPath); ok {
			s.docs[uri] = doc
			return doc, nil
		}
	}
	data, err := fs.ReadFile(s.files, relPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load a schema for URI: %s: %w", uri, err)
	}
	doc, err := decodeJSON(data)
	if err != nil {
		return nil, fmt.Errorf("malformed schema for URI: %s: %w", uri, err)
	}
	if shared {
		sharedSchemaDocs.Store(relPath, doc)
	}
	s.docs[uri] = doc
	return doc, nil
}

// resolveRef resolves a $ref value seen while evaluating a schema loaded from baseURI. It
// returns the URI of the document the fragment lives in and the node the fragment points to.
// Relative document references (RFC 3986, no scheme) resolve against baseURI's directory - the
// same resolution the upstream JSON Reference loader performs, verified against every $ref this
// corpus actually contains (see the file header).
func (s *schemaSet) resolveRef(baseURI, ref string) (string, any, error) {
	docPart, fragment, _ := strings.Cut(ref, "#")

	targetURI := baseURI
	if docPart != "" {
		if strings.Contains(docPart, "://") {
			targetURI = docPart
		} else {
			targetURI = path.Join(path.Dir(baseURI), docPart)
		}
	}

	root, err := s.load(targetURI)
	if err != nil {
		return "", nil, err
	}

	node, err := jsonPointerGet(root, fragment)
	if err != nil {
		return "", nil, fmt.Errorf("$ref %q (resolved %s#%s): %w", ref, targetURI, fragment, err)
	}
	return targetURI, node, nil
}

// jsonPointerGet navigates a JSON Pointer (RFC 6901) fragment - as found after '#' in a $ref -
// against root. Each segment is percent-decoded (schemas encode '#' inside property names, e.g.
// "x5t%23o", as %23) then JSON-Pointer-unescaped ("~1" -> "/", "~0" -> "~").
func jsonPointerGet(root any, fragment string) (any, error) {
	fragment = strings.TrimPrefix(fragment, "/")
	if fragment == "" {
		return root, nil
	}
	node := root
	for _, raw := range strings.Split(fragment, "/") {
		seg, err := url.QueryUnescape(raw)
		if err != nil {
			seg = raw
		}
		seg = strings.ReplaceAll(seg, "~1", "/")
		seg = strings.ReplaceAll(seg, "~0", "~")
		obj, ok := node.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot navigate into non-object at segment %q", seg)
		}
		next, ok := obj[seg]
		if !ok {
			return nil, fmt.Errorf("no such definition: %q", seg)
		}
		node = next
	}
	return node, nil
}

// validate evaluates instance against schema (loaded relative to baseURI), appending a
// human-readable, path-prefixed error string for every violation found. It never stops at the
// first error (siblings under allOf/properties/items are all evaluated), matching the upstream
// contract of returning every error, empty when valid.
func (s *schemaSet) validate(schemaNode any, baseURI string, instance any, path string) []string {
	schema, ok := schemaNode.(map[string]any)
	if !ok {
		// A bare schema (not an object, e.g. `true`/`false`) never appears in this corpus;
		// treat as "no constraint" rather than panicking on unexpected shapes.
		return nil
	}

	if ref, ok := schema["$ref"].(string); ok {
		targetURI, node, err := s.resolveRef(baseURI, ref)
		if err != nil {
			return []string{fmt.Sprintf("%s: %v", errPath(path), err)}
		}
		return s.validate(node, targetURI, instance, path)
	}

	var errs []string

	if declared, ok := schema["type"].(string); ok {
		if !typeMatches(declared, instance) {
			errs = append(errs, fmt.Sprintf("%s: expected type %q, got %q", errPath(path), declared, jsonType(instance)))
			return errs
		}
	}

	if enumValues, ok := schema["enum"].([]any); ok {
		matched := false
		for _, v := range enumValues {
			if reflect.DeepEqual(v, instance) {
				matched = true
				break
			}
		}
		if !matched {
			errs = append(errs, fmt.Sprintf("%s: value is not one of the enumerated values %v", errPath(path), enumValues))
		}
	}

	if str, isString := instance.(string); isString {
		if pattern, ok := schema["pattern"].(string); ok {
			if re := compilePattern(pattern); re != nil && !re.MatchString(str) {
				errs = append(errs, fmt.Sprintf("%s: value does not match pattern %q", errPath(path), pattern))
			}
		}
		// "contentEncoding" is a draft-07 annotation, not a validation assertion (its base64
		// values in this corpus are frequently base64url in practice, e.g. 'x5t#S256',
		// 'signature'), so it is intentionally never enforced here - see the file header.
	}

	if obj, isObject := instance.(map[string]any); isObject {
		errs = append(errs, s.validateObject(schema, baseURI, obj, path)...)
	}

	if arr, isArray := instance.([]any); isArray {
		errs = append(errs, s.validateArray(schema, baseURI, arr, path)...)
	}

	errs = append(errs, s.validateCombinators(schema, baseURI, instance, path)...)

	return errs
}

// validateObject applies the object-only keywords (required/properties/additionalProperties/
// minProperties/maxProperties), only reachable when instance is itself a JSON object.
func (s *schemaSet) validateObject(schema map[string]any, baseURI string, obj map[string]any, path string) []string {
	var errs []string

	if required, ok := schema["required"].([]any); ok {
		for _, r := range required {
			name, _ := r.(string)
			if _, present := obj[name]; !present {
				errs = append(errs, fmt.Sprintf("%s: missing required property %q", errPath(path), name))
			}
		}
	}

	if n, ok := schema["minProperties"].(json.Number); ok {
		if min, err := n.Int64(); err == nil && int64(len(obj)) < min {
			errs = append(errs, fmt.Sprintf("%s: expected at least %d properties, got %d", errPath(path), min, len(obj)))
		}
	}
	if n, ok := schema["maxProperties"].(json.Number); ok {
		if max, err := n.Int64(); err == nil && int64(len(obj)) > max {
			errs = append(errs, fmt.Sprintf("%s: expected at most %d properties, got %d", errPath(path), max, len(obj)))
		}
	}

	properties, _ := schema["properties"].(map[string]any)
	for name, value := range obj {
		if sub, declared := properties[name]; declared {
			errs = append(errs, s.validate(sub, baseURI, value, joinPath(path, name))...)
		}
	}

	if additional, ok := schema["additionalProperties"]; ok {
		if allowed, isBool := additional.(bool); isBool && !allowed {
			for name := range obj {
				if _, declared := properties[name]; !declared {
					errs = append(errs, fmt.Sprintf("%s: additionalProperties not allowed: %q", errPath(path), name))
				}
			}
		}
	}

	return errs
}

// validateArray applies the array-only keywords (items/minItems/maxItems), only reachable when
// instance is itself a JSON array. Only the single-schema ("list validation") form of "items" is
// used by this corpus, never the tuple-validation array-of-schemas form.
func (s *schemaSet) validateArray(schema map[string]any, baseURI string, arr []any, path string) []string {
	var errs []string

	if n, ok := schema["minItems"].(json.Number); ok {
		if min, err := n.Int64(); err == nil && int64(len(arr)) < min {
			errs = append(errs, fmt.Sprintf("%s: expected at least %d items, got %d", errPath(path), min, len(arr)))
		}
	}
	if n, ok := schema["maxItems"].(json.Number); ok {
		if max, err := n.Int64(); err == nil && int64(len(arr)) > max {
			errs = append(errs, fmt.Sprintf("%s: expected at most %d items, got %d", errPath(path), max, len(arr)))
		}
	}

	if items, ok := schema["items"]; ok {
		for i, v := range arr {
			errs = append(errs, s.validate(items, baseURI, v, fmt.Sprintf("%s[%d]", path, i))...)
		}
	}

	return errs
}

// validateCombinators applies allOf/anyOf/oneOf/not, which - per JSON Schema draft-07 - apply
// regardless of the instance's JSON type.
func (s *schemaSet) validateCombinators(schema map[string]any, baseURI string, instance any, path string) []string {
	var errs []string

	if allOf, ok := schema["allOf"].([]any); ok {
		for _, sub := range allOf {
			errs = append(errs, s.validate(sub, baseURI, instance, path)...)
		}
	}

	if anyOf, ok := schema["anyOf"].([]any); ok {
		matched := false
		var summaries []string
		var subErrs [][]string
		for _, sub := range anyOf {
			branchErrs := s.validate(sub, baseURI, instance, path)
			if len(branchErrs) == 0 {
				matched = true
				break
			}
			summaries = append(summaries, schemaSummary(sub))
			subErrs = append(subErrs, branchErrs)
		}
		if !matched {
			errs = append(errs, fmt.Sprintf("%s: value must match at least one of the schemas [%s]", errPath(path), strings.Join(summaries, " | ")))
			// Surface every alternative's own errors too - not just the summary - so a caller
			// can still find the substring naming the actual offending property/keyword, the
			// same way an "everit"-style validator's nested ValidationFailure causes would.
			for _, branchErrs := range subErrs {
				errs = append(errs, branchErrs...)
			}
		}
	}

	if oneOf, ok := schema["oneOf"].([]any); ok {
		matches := 0
		var summaries []string
		var subErrs [][]string
		for _, sub := range oneOf {
			branchErrs := s.validate(sub, baseURI, instance, path)
			if len(branchErrs) == 0 {
				matches++
			} else {
				subErrs = append(subErrs, branchErrs)
			}
			summaries = append(summaries, schemaSummary(sub))
		}
		if matches != 1 {
			errs = append(errs, fmt.Sprintf("%s: value must match exactly one of the schemas [%s], matched %d", errPath(path), strings.Join(summaries, " | "), matches))
			if matches == 0 {
				// Only surface branch detail on the "none matched" case: when matches > 1, every
				// branch already validated cleanly, so there is nothing more specific to report.
				for _, branchErrs := range subErrs {
					errs = append(errs, branchErrs...)
				}
			}
		}
	}

	if not, ok := schema["not"]; ok {
		if len(s.validate(not, baseURI, instance, path)) == 0 {
			errs = append(errs, fmt.Sprintf("%s: value must not match schema (%s)", errPath(path), schemaSummary(not)))
		}
	}

	return errs
}

// schemaSummary renders a short, human-readable description of a schema node for use inside
// oneOf/anyOf/not error messages - the "required" property list when present (which is how every
// oneOf/not alternative in this corpus is expressed), else the $ref target, else "schema".
func schemaSummary(schemaNode any) string {
	m, ok := schemaNode.(map[string]any)
	if !ok {
		return "schema"
	}
	if ref, ok := m["$ref"].(string); ok {
		return ref
	}
	if required, ok := m["required"].([]any); ok {
		names := make([]string, 0, len(required))
		for _, r := range required {
			if name, ok := r.(string); ok {
				names = append(names, name)
			}
		}
		return "required=[" + strings.Join(names, ",") + "]"
	}
	return "schema"
}

// jsonType reports the JSON Schema type name of a decodeJSON-produced value.
func jsonType(instance any) string {
	switch v := instance.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return "string"
	case json.Number:
		if isIntegerNumber(string(v)) {
			return "integer"
		}
		return "number"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	default:
		return "unknown"
	}
}

// isIntegerNumber reports whether a JSON number literal (as produced by json.Decoder, so already
// syntactically valid) is typed "integer" by the validator upstream uses, com.github.erosb.
// jsonsKema (0.31.0). Draft-06+ defines an integer as a number with a zero fractional part, so
// 1.0 is one, and that validator agrees - but only up to the way it decides it, which is
// reproduced here (Validator.findActualNumberType, JsonParser.parseNumber):
//
//   - a literal without '.', 'e' or 'E' is an integer;
//   - any other literal is parsed as a Double and its Double.toString() is inspected: integer when
//     nothing but zeros follows the '.'. Double.toString switches to scientific notation
//     ("1.0E7", "1.7E9") from 10^7 on, where an 'E' follows the '.', so 1700000000.0 and 1e7 are
//     NOT integers there while 1.0, 5e0 and 12.0 are;
//   - a literal too large for a Double falls back to BigDecimal.toString(), which is integer when
//     the coefficient is a single digit ("1E+999").
//
// Before, every literal with a '.' or an exponent was a "number".
func isIntegerNumber(literal string) bool {
	if !strings.ContainsAny(literal, ".eE") {
		return true
	}
	value, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		// out of range for a Double (+/-Inf): BigDecimal.toString() of the literal
		mantissa, _, _ := strings.Cut(strings.ToLower(strings.TrimPrefix(literal, "-")), "e")
		coefficient := strings.TrimLeft(strings.Replace(mantissa, ".", "", 1), "0")
		return len(coefficient) <= 1
	}
	magnitude := math.Abs(value)
	if magnitude == 0 {
		return true // "0.0"
	}
	// Double.toString is plain decimal ("N.0" for a whole value) only within [1e-3, 1e7)
	return magnitude >= 1e-3 && magnitude < 1e7 && value == math.Trunc(value)
}

// typeMatches reports whether instance's JSON type satisfies the declared draft-07 "type" value,
// including the "integer" <: "number" relationship.
func typeMatches(declared string, instance any) bool {
	actual := jsonType(instance)
	switch declared {
	case "number":
		return actual == "number" || actual == "integer"
	case "integer":
		return actual == "integer"
	default:
		return actual == declared
	}
}

// joinPath appends a property name to a validation path, dot-separated.
func joinPath(base, name string) string {
	if base == "" {
		return name
	}
	return base + "." + name
}

// errPath renders the empty (root) path the way Java's ValidationFailure#getInstance().
// getLocation() renders it: as "$" (JSON Pointer document root).
func errPath(path string) string {
	if path == "" {
		return "$"
	}
	return "$." + path
}

// validateJSONAgainstSchema parses jsonText and validates it against the schema registered under
// rootURI in defs (resolved against files), returning the list of validation error messages -
// empty when the JSON is valid, per validateAgainstSchema's upstream contract.
func validateJSONAgainstSchema(files fs.FS, defs map[string]string, rootURI, jsonText string) []string {
	instance, err := decodeJSON([]byte(jsonText))
	if err != nil {
		return []string{fmt.Sprintf("unable to parse JSON: %v", err)}
	}

	set := newSchemaSet(files, defs)
	root, err := set.load(rootURI)
	if err != nil {
		return []string{err.Error()}
	}

	return set.validate(root, rootURI, instance, "")
}

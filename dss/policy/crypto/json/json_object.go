// Collapses the slice of eu.europa.esig.json.JsonObjectWrapper and
// eu.europa.esig.json.JSONParser this package's callers use directly, using
// Go's stdlib encoding/json instead of dss-json-common's
// com.github.erosb.jsonsKema parser - see doc.go for why.
package cryptojson

import (
	"encoding/json"
	"errors"
	"io"
)

// jsonObject wraps a parsed JSON object (a map, per encoding/json's decode
// of a JSON object into interface{}) and provides the subset of
// JsonObjectWrapper's typed accessors this package's callers use. Ports
// JsonObjectWrapper, minus the isEmpty()/getJsonObject() methods this
// package never calls.
type jsonObject map[string]interface{}

// parseJSONObject parses r and returns its root value as a jsonObject.
// Ports `new JSONParser().parse(is)` followed by the caller's own check
// that the parsed value is a JSON object (JSONParser.parse itself accepts
// any JSON value; CryptographicSuiteJsonFactory is the only caller in this
// package, and it immediately requires an object - see
// cryptographic_suite_json_factory.go).
func parseJSONObject(r io.Reader) (jsonObject, error) {
	var raw interface{}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, errors.New("parsed JSON cannot be null")
	}
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.New("parsed JSON cannot be null")
	}
	return jsonObject(obj), nil
}

// getAsObject gets a value of the header name as a JSON object. If not
// present, or not able to convert, returns nil. Ports
// JsonObjectWrapper#getAsObject; slf4j warnings are dropped throughout,
// per PORTING.md.
func (o jsonObject) getAsObject(name string) jsonObject {
	v, ok := o[name]
	if !ok {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return jsonObject(m)
}

// getAsString gets a value of the header name as a string. If not present,
// or not able to convert, returns "" - the Go zero value stands in for
// Java's null return, matching every caller in this package, whose
// underlying model fields (dss/model/policy.CryptographicSuiteMetadata et
// al.) are plain (non-pointer) strings with the same "" empty/unset
// default. Ports JsonObjectWrapper#getAsString.
func (o jsonObject) getAsString(name string) string {
	v, ok := o[name]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// getAsNumber gets a value of the header name as a float64 (encoding/json's
// decoding target for every JSON number) and whether it was present and of
// the expected type. Ports JsonObjectWrapper#getAsNumber, whose null return
// this package's only caller (buildParameter, via toInteger) distinguishes
// from a present value - hence the explicit ok return instead of a
// zero-value sentinel.
func (o jsonObject) getAsNumber(name string) (float64, bool) {
	v, ok := o[name]
	if !ok {
		return 0, false
	}
	n, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return n, true
}

// getAsObjectList gets a value of the header name as a list of JSON
// objects. If not present, or not able to convert, returns nil (ranges
// identically to Java's Collections.emptyList()). Ports
// JsonObjectWrapper#getAsObjectList.
func (o jsonObject) getAsObjectList(name string) []jsonObject {
	v, ok := o[name]
	if !ok {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []jsonObject
	for _, elem := range arr {
		m, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, jsonObject(m))
	}
	return result
}

// getAsStringList gets a value of the header name as a list of strings. If
// not present, or not able to convert, returns nil (ranges identically to
// Java's Collections.emptyList()). Ports
// JsonObjectWrapper#getAsStringList.
func (o jsonObject) getAsStringList(name string) []string {
	v, ok := o[name]
	if !ok {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, elem := range arr {
		s, ok := elem.(string)
		if !ok {
			continue
		}
		result = append(result, s)
	}
	return result
}

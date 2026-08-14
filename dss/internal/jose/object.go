// Ported from the java.util.LinkedHashMap / java.util.HashMap that back a JOSE header in
// org.jose4j.jwx.Headers and eu.europa.esig.dss.jades.JsonObject (jose4j 0.9.6, DSS 6.5.RC1).
package jose

// Object is a JSON object that remembers the order of its members.
//
// It exists because a Go map cannot be serialized reproducibly: Go randomizes map iteration order
// deliberately, while every JSON object jose4j writes comes out of a LinkedHashMap (insertion
// order) or a HashMap (bucket order), and those bytes are what a JWS signs. Object therefore
// keeps the keys in a slice next to the value map and renders them through Keys.
//
// Two orderings are supported, chosen at construction and fixed for the object's lifetime:
//
//	NewObject()      insertion order   - java.util.LinkedHashMap, the JOSE header and every
//	                                     DSS site that builds one with new LinkedHashMap<>()
//	NewHashObject()  Java bucket order - java.util.HashMap, which is what DSS's own
//	                                     JsonObject() no-arg constructor wraps
//
// The zero value is not usable; construct with NewObject or NewHashObject.
type Object struct {
	keys      []string
	values    map[string]any
	hashOrder bool
}

// ObjectHolder is implemented by a type that wraps a JSON object and should be serialized as
// one. It stands in for `value instanceof Map` in JSONValue.writeJSONString, which is how
// eu.europa.esig.dss.jades.JsonObject - a Map implementation rather than a JSONObject subclass -
// reaches the writer upstream.
type ObjectHolder interface {
	// JSONObject returns the wrapped object.
	JSONObject() *Object
}

// NewObject returns an empty insertion-ordered object - the Go counterpart of
// `new LinkedHashMap<String, Object>()`.
func NewObject() *Object {
	return &Object{values: make(map[string]any)}
}

// NewHashObject returns an empty object that iterates in java.util.HashMap bucket order - the
// counterpart of `new HashMap<String, Object>()`, and specifically of what DSS's JsonObject()
// no-arg constructor wraps. Use it only where upstream really does use a bare HashMap;
// everywhere else NewObject is both correct and easier to reason about.
func NewHashObject() *Object {
	return &Object{values: make(map[string]any), hashOrder: true}
}

// NewObjectFromPairs builds an insertion-ordered object from alternating key/value arguments.
// It panics on an odd number of arguments or a non-string key, both of which are programming
// errors rather than input errors.
func NewObjectFromPairs(kv ...any) *Object {
	if len(kv)%2 != 0 {
		panic("jose: NewObjectFromPairs needs an even number of arguments")
	}
	o := NewObject()
	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			panic("jose: NewObjectFromPairs key is not a string")
		}
		o.Put(key, kv[i+1])
	}
	return o
}

// Put inserts or replaces the value for key and returns the previous value, mirroring
// Map.put's return. A key that is already present keeps its position in the iteration order,
// which is LinkedHashMap's documented behaviour for a re-put (access order is not enabled).
func (o *Object) Put(key string, value any) any {
	previous, existed := o.values[key]
	if !existed {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
	if !existed {
		return nil
	}
	return previous
}

// PutAll copies every member of other, in other's iteration order. Port of Map.putAll.
func (o *Object) PutAll(other *Object) {
	if other == nil {
		return
	}
	for _, k := range other.Keys() {
		o.Put(k, other.values[k])
	}
}

// Get returns the value stored under key and whether it was present. A nil *Object reads as
// empty, matching a Go map and Java's behaviour of returning null from an absent key.
func (o *Object) Get(key string) (any, bool) {
	if o == nil {
		return nil, false
	}
	v, ok := o.values[key]
	return v, ok
}

// Value returns the value stored under key, or nil if absent - the direct counterpart of
// Map.get, which cannot distinguish "absent" from "null" either.
func (o *Object) Value(key string) any {
	if o == nil {
		return nil
	}
	return o.values[key]
}

// ContainsKey reports whether key is present. Port of Map.containsKey.
func (o *Object) ContainsKey(key string) bool {
	if o == nil {
		return false
	}
	_, ok := o.values[key]
	return ok
}

// Remove deletes key, if present, and returns the value it held. Port of Map.remove.
func (o *Object) Remove(key string) any {
	if o == nil {
		return nil
	}
	previous, existed := o.values[key]
	if !existed {
		return nil
	}
	delete(o.values, key)
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	return previous
}

// Replace sets a new value for key only if key is already present, and returns the previous
// value. Port of Map.replace, which JWSConverter uses to swap the 'etsiU' array in place -
// notably without appending the key if it had been absent.
func (o *Object) Replace(key string, value any) any {
	if o == nil {
		return nil
	}
	previous, existed := o.values[key]
	if !existed {
		return nil
	}
	o.values[key] = value
	return previous
}

// Clear removes every member. Port of Map.clear.
func (o *Object) Clear() {
	o.keys = nil
	o.values = make(map[string]any)
}

// Size returns the number of members. Port of Map.size.
func (o *Object) Size() int {
	if o == nil {
		return 0
	}
	return len(o.values)
}

// IsEmpty reports whether the object has no members. Port of Map.isEmpty.
func (o *Object) IsEmpty() bool { return o.Size() == 0 }

// Keys returns the member names in iteration order: insertion order for an object built with
// NewObject, java.util.HashMap bucket order for one built with NewHashObject.
//
// The returned slice is a fresh copy, so a caller cannot reorder the object by sorting it.
func (o *Object) Keys() []string {
	if o == nil {
		return nil
	}
	if o.hashOrder {
		return JavaHashMapOrder(o.keys)
	}
	out := make([]string, len(o.keys))
	copy(out, o.keys)
	return out
}

// Values returns the member values in iteration order. Port of Map.values().
func (o *Object) Values() []any {
	keys := o.Keys()
	out := make([]any, len(keys))
	for i, k := range keys {
		out[i] = o.values[k]
	}
	return out
}

// IsHashOrdered reports whether the object iterates in java.util.HashMap bucket order.
func (o *Object) IsHashOrdered() bool { return o != nil && o.hashOrder }

// Clone returns a shallow copy with the same ordering mode: the member values are shared, only
// the key order and the map itself are duplicated.
func (o *Object) Clone() *Object {
	if o == nil {
		return nil
	}
	c := &Object{
		keys:      make([]string, len(o.keys)),
		values:    make(map[string]any, len(o.values)),
		hashOrder: o.hashOrder,
	}
	copy(c.keys, o.keys)
	for k, v := range o.values {
		c.values[k] = v
	}
	return c
}

// String renders the object as JSON, so that a %v of an Object reads like the Java toString()
// of a JSONObject rather than like a Go struct dump.
func (o *Object) String() string { return JSON(o) }

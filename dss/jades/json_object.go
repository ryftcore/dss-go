// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JsonObject.java (DSS 6.5.RC1).
package jades

import "github.com/utain/esig/dss/internal/jose"

// JsonObject is a wrapper of a map with JsonObject methods. Port of the class JsonObject, which
// upstream implements java.util.Map<String, Object> and renders through
// org.jose4j.json.JsonUtil.toJson.
//
// Go has no Map interface to implement, so the wrapper keeps the same method surface over an
// internal/jose Object and declares jose.ObjectHolder, which is what puts it down the writer's
// object branch rather than its toString() fallback - the exact effect `implements Map` has in
// Java.
//
// The ordering of the wrapped map is the whole point of the type, and the two constructors
// differ in it:
//
//	NewJsonObject()          wraps a java.util.HashMap, so members come out in Java bucket order
//	NewJsonObjectFromMap(m)  wraps whatever m is, normally a LinkedHashMap built in code order
//
// That asymmetry is upstream's (`map = new HashMap<>()` versus `map = m`) and it is visible in
// signed bytes: JAdESLevelBaselineLT builds 'rVals' and 'tstVd' through the no-argument
// constructor, so their members are serialized in HashMap order and an archive timestamp covers
// exactly those bytes.
type JsonObject struct {
	// map is the wrapped map. Port of the private field of the same name.
	object *jose.Object
}

// NewJsonObject creates an object over an empty HashMap. Port of the JsonObject() constructor.
//
// jose.NewHashObject reproduces java.util.HashMap's iteration order; see its documentation for
// why that is modelled rather than approximated.
func NewJsonObject() *JsonObject {
	return &JsonObject{object: jose.NewHashObject()}
}

// NewJsonObjectFromMap wraps the provided map. Port of the JsonObject(Map) constructor,
// including its Objects.requireNonNull - PORTING.md maps requireNonNull onto a panic carrying
// the Java message.
func NewJsonObjectFromMap(m *jose.Object) *JsonObject {
	if m == nil {
		panic("Map cannot be null!")
	}
	return &JsonObject{object: m}
}

// JSONObject returns the wrapped map, satisfying jose.ObjectHolder so that the writer renders
// this value as a JSON object. It is also how the other jades files reach the members without
// going through the Map facade.
func (o *JsonObject) JSONObject() *jose.Object {
	if o == nil {
		return nil
	}
	return o.object
}

// ToJSONString converts the object to its JSON string representation. Port of toJSONString().
func (o *JsonObject) ToJSONString() string {
	return jose.JSON(o.object)
}

// Size returns the number of members. Port of size().
func (o *JsonObject) Size() int { return o.object.Size() }

// IsEmpty reports whether the object has no members. Port of isEmpty().
func (o *JsonObject) IsEmpty() bool { return o.object.IsEmpty() }

// ContainsKey reports whether key is present. Port of containsKey(Object).
func (o *JsonObject) ContainsKey(key string) bool { return o.object.ContainsKey(key) }

// ContainsValue reports whether any member holds value. Port of containsValue(Object).
//
// Comparison is Go's ==, which panics on an uncomparable dynamic type (a slice, say) exactly
// where Java's equals() would have compared by identity; no DSS call site uses this method, and
// it exists only to keep the ported surface complete.
func (o *JsonObject) ContainsValue(value any) bool {
	for _, v := range o.object.Values() {
		if v == value {
			return true
		}
	}
	return false
}

// Get returns the value stored under key, or nil. Port of get(Object).
func (o *JsonObject) Get(key string) any { return o.object.Value(key) }

// Put stores value under key and returns the previous value. Port of put(String, Object).
func (o *JsonObject) Put(key string, value any) any { return o.object.Put(key, value) }

// Remove deletes key and returns the value it held. Port of remove(Object).
func (o *JsonObject) Remove(key string) any { return o.object.Remove(key) }

// PutAll copies every member of m. Port of putAll(Map).
func (o *JsonObject) PutAll(m *jose.Object) { o.object.PutAll(m) }

// Clear removes every member. Port of clear().
func (o *JsonObject) Clear() { o.object.Clear() }

// KeySet returns the member names in iteration order. Port of keySet(); Java's Set becomes an
// ordered slice, since the iteration order is the reason this type exists.
func (o *JsonObject) KeySet() []string { return o.object.Keys() }

// Values returns the member values in iteration order. Port of values().
func (o *JsonObject) Values() []any { return o.object.Values() }

// String returns the JSON representation. Port of toString(), which delegates to toJSONString().
func (o *JsonObject) String() string { return o.ToJSONString() }

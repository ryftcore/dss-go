// Ported from dss-model/.../claim/ClaimMap.java (DSS 6.5.RC1).
package claim

import (
	"reflect"
	"sort"
	"strings"
)

// ClaimMap represents a Map encoded (selectively) disclosable claim. It is
// designed for embedding by concrete subtypes: Java's abstract
// `getKeyAsString(Object)` and `createClaim(String, Object)` methods are
// ported as the GetKeyAsString and CreateClaim function fields, which
// embedders MUST set (typically from their own constructor) before
// calling MapValue/Keys/Get/ValueAsString and friends.
type ClaimMap struct {
	AbstractClaim

	// value is the map value of the claim (mirrors Java's Map<?,?>).
	value map[any]any

	// GetKeyAsString implements the abstract getKeyAsString(Object)
	// method: returns a map key as a string.
	GetKeyAsString func(key any) string

	// CreateClaim implements the abstract createClaim(String, Object)
	// method: creates a Claim for a map entry.
	CreateClaim func(name string, value any) Claim
}

// NewClaimMap ports the protected default constructor.
func NewClaimMap(value map[any]any, getKeyAsString func(key any) string, createClaim func(name string, value any) Claim) *ClaimMap {
	return NewClaimMapWithParent("", value, false, nil, getKeyAsString, createClaim)
}

// NewClaimMapWithParent ports the constructor with claim name, value,
// selectively disclosable status and parent claim provided.
func NewClaimMapWithParent(name string, value map[any]any, selectivelyDisclosable bool, parent Claim, getKeyAsString func(key any) string, createClaim func(name string, value any) Claim) *ClaimMap {
	return &ClaimMap{
		AbstractClaim:  NewAbstractClaimWithParent(name, selectivelyDisclosable, parent),
		value:          value,
		GetKeyAsString: getKeyAsString,
		CreateClaim:    createClaim,
	}
}

// NewClaimMapFull ports the constructor with claim name, namespace, value,
// selectively disclosable status and parent claim provided.
func NewClaimMapFull(name, namespace string, value map[any]any, selectivelyDisclosable bool, parent Claim, getKeyAsString func(key any) string, createClaim func(name string, value any) Claim) *ClaimMap {
	return &ClaimMap{
		AbstractClaim:  NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:          value,
		GetKeyAsString: getKeyAsString,
		CreateClaim:    createClaim,
	}
}

// MapValue converts every map entry into a Claim, keyed by
// GetKeyAsString(key). Ports ClaimMap#getMapValue.
func (c *ClaimMap) MapValue() map[string]Claim {
	if len(c.value) == 0 {
		return map[string]Claim{}
	}
	result := make(map[string]Claim, len(c.value))
	for k, v := range c.value {
		headerName := c.GetKeyAsString(k)
		result[headerName] = c.CreateClaim(headerName, v)
	}
	return result
}

// Keys gets a set of map keys, as a String. Ports ClaimMap#getKeys.
// NOTE: Java returns a java.util.Set<String> (unordered); this port
// returns an unordered []string built by ranging the underlying map, so
// iteration order is likewise unspecified between calls.
func (c *ClaimMap) Keys() []string {
	keys := make([]string, 0, len(c.value))
	for k := range c.value {
		keys = append(keys, c.GetKeyAsString(k))
	}
	return keys
}

// sortedKeys is an unexported helper for ValueAsString, which needs a
// deterministic key order (Java's own iteration order over its Map<?,?>
// is likewise implementation-defined, so this is a best-effort rendering
// aid, not a fidelity requirement).
func (c *ClaimMap) sortedKeys() []string {
	keys := c.Keys()
	sort.Strings(keys)
	return keys
}

// Get gets the claim for the corresponding header name key. Ports
// ClaimMap#get.
func (c *ClaimMap) Get(headerName string) Claim {
	return c.MapValue()[headerName]
}

// GetAsMap gets the claim value that is a map from the current map using
// headerName as a key. Ports ClaimMap#getAsMap(String).
func (c *ClaimMap) GetAsMap(headerName string) *ClaimMap {
	return c.getAsMap(c.Get(headerName))
}

// getAsMap checks if claim is of map type and returns its value as
// *ClaimMap. Ports ClaimMap#getAsMap(Claim) (protected).
func (c *ClaimMap) getAsMap(claim Claim) *ClaimMap {
	if claim != nil && claim.IsMapValueType() {
		if cm, ok := claim.(*ClaimMap); ok {
			return cm
		}
	}
	return nil
}

// GetAsArray gets the claim value that is an array from the current map
// using headerName as a key. Ports ClaimMap#getAsArray(String).
func (c *ClaimMap) GetAsArray(headerName string) *ClaimArray {
	return c.getAsArray(c.Get(headerName))
}

// getAsArray checks if claim is of array type and returns its value as
// *ClaimArray. Ports ClaimMap#getAsArray(Claim) (protected).
func (c *ClaimMap) getAsArray(claim Claim) *ClaimArray {
	if claim != nil && claim.IsArrayValueType() {
		if ca, ok := claim.(*ClaimArray); ok {
			return ca
		}
	}
	return nil
}

// GetAsNumber gets the claim value that is a number from the current map
// using headerName as a key. Ports ClaimMap#getAsNumber(String).
func (c *ClaimMap) GetAsNumber(headerName string) *ClaimNumber {
	return c.getAsNumber(c.Get(headerName))
}

// getAsNumber checks if claim is of number type and returns its value as
// *ClaimNumber. Ports ClaimMap#getAsNumber(Claim) (protected).
func (c *ClaimMap) getAsNumber(claim Claim) *ClaimNumber {
	if claim != nil && claim.IsNumberValueType() {
		if cn, ok := claim.(*ClaimNumber); ok {
			return cn
		}
	}
	return nil
}

// GetAsString gets the claim value that is a string from the current map
// using headerName as a key. Ports ClaimMap#getAsString(String).
func (c *ClaimMap) GetAsString(headerName string) *ClaimString {
	return c.getAsString(c.Get(headerName))
}

// getAsString checks if claim is of string type and returns its value as
// *ClaimString. Ports ClaimMap#getAsString(Claim) (protected).
func (c *ClaimMap) getAsString(claim Claim) *ClaimString {
	if claim != nil && claim.IsStringValueType() {
		if cs, ok := claim.(*ClaimString); ok {
			return cs
		}
	}
	return nil
}

// GetAsBoolean gets the claim value that is a boolean from the current map
// using headerName as a key. Ports ClaimMap#getAsBoolean(String).
func (c *ClaimMap) GetAsBoolean(headerName string) *ClaimBoolean {
	return c.getAsBoolean(c.Get(headerName))
}

// getAsBoolean checks if claim is of boolean type and returns its value as
// *ClaimBoolean. Ports ClaimMap#getAsBoolean(Claim) (protected).
func (c *ClaimMap) getAsBoolean(claim Claim) *ClaimBoolean {
	if claim != nil && claim.IsBooleanValueType() {
		if cb, ok := claim.(*ClaimBoolean); ok {
			return cb
		}
	}
	return nil
}

// GetAsNull gets the claim value that is null from the current map using
// headerName as a key. Ports ClaimMap#getAsNull(String).
func (c *ClaimMap) GetAsNull(headerName string) *ClaimNull {
	return c.getAsNull(c.Get(headerName))
}

// getAsNull checks if claim is of null type and returns its value as
// *ClaimNull. Ports ClaimMap#getAsNull(Claim) (protected).
func (c *ClaimMap) getAsNull(claim Claim) *ClaimNull {
	if claim != nil && claim.IsNullValueType() {
		if cn, ok := claim.(*ClaimNull); ok {
			return cn
		}
	}
	return nil
}

// IsMapValueType always returns true.
func (c *ClaimMap) IsMapValueType() bool { return true }

// IsNullOrEmpty ports ClaimMap#isNullOrEmpty.
func (c *ClaimMap) IsNullOrEmpty() bool { return len(c.value) == 0 }

// Size gets the number of entries within the map. Ports ClaimMap#getSize.
func (c *ClaimMap) Size() int {
	if c.IsNullOrEmpty() {
		return 0
	}
	return len(c.value)
}

// ValueAsString ports ClaimMap#getValueAsString.
func (c *ClaimMap) ValueAsString() string {
	var sb strings.Builder
	sb.WriteString("{")
	keys := c.sortedKeys()
	for i, key := range keys {
		sb.WriteString("\"")
		sb.WriteString(key)
		sb.WriteString("\": ")
		claimValue := c.Get(key)
		if claimValue.IsStringValueType() {
			sb.WriteString("\"")
		}
		sb.WriteString(claimValue.ValueAsString())
		if claimValue.IsStringValueType() {
			sb.WriteString("\"")
		}
		if i < len(keys)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

// Equals ports ClaimMap#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ClaimMap) Equals(other *ClaimMap) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if !c.AbstractClaim.Equals(&other.AbstractClaim) {
		return false
	}
	return reflect.DeepEqual(c.value, other.value)
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimMap) String() string { return AbstractClaimString(c) }

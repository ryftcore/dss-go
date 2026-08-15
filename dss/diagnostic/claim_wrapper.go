// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/claim/ClaimWrapper.java (DSS 6.5.RC1).
//
// Java's per-claim-shape subtypes (XmlAddressClaim, XmlCredentialSubjectClaim, ...) all extend
// the base XmlClaim class in the generated JAXB model, so a ClaimWrapper built over any of them
// can be viewed as - or built from - the common XmlClaim through an ordinary upcast;
// ClaimWrapper.getWrapped() itself only ever declares XmlClaim. The generated Go jaxb structs are
// not modelled with that inheritance (see jaxb_claim.go): each concrete Xml*Claim type
// independently embeds XmlClaimContent/XmlClaimAttrs rather than an XmlClaim. Every claim subtype
// wrapper in this package (AddressClaimWrapper, BirthdateClaimWrapper, ...) therefore embeds
// ClaimWrapper as its own base - constructed over a synthetic *jaxb.XmlClaim assembled from the
// wrapped type's embedded content/attrs via claimBase (defined in eaa_payload_proxy.go,
// DIAGWRAP_A, and used here as an ordinary same-package function) - and exposes
// AsClaim() *ClaimWrapper to view itself as the Java base type for contexts needing a homogeneous
// []*ClaimWrapper the way Java's covariant upcast provides for free (see the same convention
// documented in eaa_payload_proxy.go and eaa_wrapper.go).
//
// A child claim's Parent() is set to the plain embedded &w.ClaimWrapper (never w.AsClaim()) by
// every overriding subtype's own getters: AsClaim() computes Map()/List() by calling those same
// getters, so a getter that itself called w.AsClaim() to build its child's parent link would
// recurse into AsClaim() again. The one place this could matter - a caller chasing
// child.Parent().Map()/.List() back up and expecting the override - is not exercised by any
// consumer in this port (DIAGWRAP_A's eaa_wrapper.go/eaa_payload_proxy.go read Parent() nowhere);
// Parent() otherwise behaves identically to Java, including for the leaf accessors that are not
// overridden and so never risk this recursion.
//
// A subtype that overrides isList()/getList() or isMap()/getMap() in Java (AddressClaimWrapper,
// CredentialSubjectClaimWrapper, DrivingPrivilege(s|Codes)ClaimWrapper, StatusClaimWrapper, ...)
// cannot rely on Go's static embedding to reach that override once the object is later used
// through its ClaimWrapper view (e.g. placed in a []*ClaimWrapper and walked recursively by
// selectivelyDisclosableClaimsRecursively in eaa_wrapper.go) - Go has no virtual dispatch back
// from an embedded base to the embedding type. Each such subtype's AsClaim() therefore bakes its
// override's current result into the returned ClaimWrapper via the unexported
// mapOverride/listOverride fields below, which IsMap/Map/IsList/List consult before falling back
// to the generic Entry/Item-derived behaviour; this reproduces Java's polymorphism without an
// override-registration mechanism. A nil override means "no override recorded, use the generic
// Entry/Item-derived behaviour"; a non-nil (possibly empty) override, including one produced by a
// conditional Java override that decided not to apply, is used verbatim.
package diagnostic

import (
	"bytes"
	"encoding/base64"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
)

// claimDateTimeFormat is the lexical form of eu.europa.esig.dss.jaxb.parsers.DateParser
// ("yyyy-MM-dd'T'HH:mm:ss'Z'" evaluated in UTC), used by DisplayValue for a DateTime claim.
const claimDateTimeFormat = "2006-01-02T15:04:05Z"

// ClaimWrapper is a user-friendly wrapper of a generic jaxb.XmlClaim object, containing
// information about a single claim extracted from an EAA Payload.
//
// The wrapper may return only one of the following values: Text; or Number; or Boolean; or
// DateTime; or Binary. Should you need to retrieve any value irrespective of its original data
// type, use DisplayValue.
type ClaimWrapper struct {
	// wrapped is the wrapped XmlClaim.
	wrapped *jaxb.XmlClaim
	// parent is the parent claim, when present.
	parent *ClaimWrapper

	// mapOverride, when non-nil (even if empty), is returned by Map instead of the generic
	// Entry-derived behaviour; see the package note above. Set independently of isMapOverride
	// because one claim subtype (AttestedAttributesSubjectClaimWrapper) overrides getMap() in
	// Java without overriding isMap(), so Map's override and IsMap's override can disagree.
	mapOverride map[string]*ClaimWrapper
	// isMapOverride, when non-nil, is returned by IsMap instead of the generic Entry-derived
	// behaviour; see mapOverride above.
	isMapOverride *bool
	// listOverride, when non-nil (even if empty), is returned by IsList/List instead of the
	// generic Item-derived behaviour; see the package note above. No claim subtype in this
	// package overrides one of isList()/getList() without the other, so - unlike the map case -
	// one field can drive both.
	listOverride []*ClaimWrapper

	// overrides points back at the embedding claim subtype, so that the base bodies Java
	// inherits unchanged - isNull(), isEmpty() and getDisplayValue() - reach a subtype's
	// isList()/getList()/isMap()/getMap() override the way Java's virtual dispatch does. nil
	// means "no subtype embeds this wrapper", i.e. dispatch to the base itself. See
	// InitClaimOverrides.
	overrides ClaimWrapperOverrides
}

// ClaimWrapperOverrides declares the four ClaimWrapper operations a claim subtype may override
// in Java and that ClaimWrapper's own inherited bodies call back into: isNull(), isEmpty() and
// getDisplayValue() are declared once on ClaimWrapper and never overridden, yet each of them
// self-calls isList()/getList()/isMap()/getMap(), which Java resolves against the concrete
// subtype. Go promotes those base bodies to the subtype unchanged and would resolve the
// self-calls against the base, so a subtype registers itself with InitClaimOverrides and the
// three inherited bodies route through this interface instead.
//
// This is deliberately narrower than the mapOverride/listOverride fields above: those exist so a
// subtype's override survives the *copy* AsClaim() hands out as a plain *ClaimWrapper, whereas
// this interface restores dispatch on the subtype value itself.
type ClaimWrapperOverrides interface {
	// IsList reports whether the claim value is of a list type. Port of isList().
	IsList() bool
	// List returns the claim value as a list. Port of getList().
	List() []*ClaimWrapper
	// IsMap reports whether the claim value is of a map type. Port of isMap().
	IsMap() bool
	// Map returns the claim value as a map. Port of getMap().
	Map() map[string]*ClaimWrapper
}

// InitClaimOverrides registers the embedding claim subtype with its ClaimWrapper base so the
// base's inherited IsNull/IsEmpty/DisplayValue bodies dispatch to that subtype's
// IsList/List/IsMap/Map, as Java's virtual dispatch does. A subtype that overrides any of the
// four calls it from each of its constructors; a plain ClaimWrapper never does.
func (c *ClaimWrapper) InitClaimOverrides(overrides ClaimWrapperOverrides) {
	c.overrides = overrides
}

// claimOverrides returns the registered subtype, or the wrapper itself when none was registered.
func (c *ClaimWrapper) claimOverrides() ClaimWrapperOverrides {
	if c.overrides == nil {
		return c
	}
	return c.overrides
}

// NewClaimWrapper is the default constructor. Port of ClaimWrapper(XmlClaim); panics per
// Objects.requireNonNull(wrapped, "XmlClaim cannot be null!").
func NewClaimWrapper(wrapped *jaxb.XmlClaim) *ClaimWrapper {
	return NewClaimWrapperWithParent(wrapped, nil)
}

// NewClaimWrapperWithParent is the constructor with a parent claim provided. Port of
// ClaimWrapper(XmlClaim, ClaimWrapper); panics per
// Objects.requireNonNull(wrapped, "XmlClaim cannot be null!").
func NewClaimWrapperWithParent(wrapped *jaxb.XmlClaim, parent *ClaimWrapper) *ClaimWrapper {
	if wrapped == nil {
		panic("XmlClaim cannot be null!")
	}
	return &ClaimWrapper{wrapped: wrapped, parent: parent}
}

// AsClaim views the wrapper as its own ClaimWrapper base type; the identity conversion every
// other claim subtype's AsClaim() bottoms out at. See the package note above.
func (c *ClaimWrapper) AsClaim() *ClaimWrapper { return c }

// Name gets the claim name. Port of getName().
func (c *ClaimWrapper) Name() string {
	if c.wrapped.Name != nil {
		return *c.wrapped.Name
	}
	return ""
}

// Namespace gets the claim's namespace (used for mdoc). Port of getNamespace().
func (c *ClaimWrapper) Namespace() string {
	if c.wrapped.Namespace != nil {
		return *c.wrapped.Namespace
	}
	return ""
}

// IsSelectivelyDisclosable gets whether the claim was made selectively disclosable and its
// value has been obtained from a provided disclosure. Port of isSelectivelyDisclosable().
func (c *ClaimWrapper) IsSelectivelyDisclosable() bool {
	return c.wrapped.Disclosure != nil && *c.wrapped.Disclosure
}

// Text gets the value as a string. If the value is not of a string type, returns "". Port of
// getText().
func (c *ClaimWrapper) Text() string {
	if c.wrapped.Text != nil {
		return *c.wrapped.Text
	}
	return ""
}

// IsText gets whether the claim value is of String type. Port of isText().
func (c *ClaimWrapper) IsText() bool { return c.wrapped.Text != nil }

// Number gets the value as a number. If the value is not of a number type, returns nil. Port
// of getNumber().
func (c *ClaimWrapper) Number() *big.Int { return c.wrapped.Number }

// IsNumber gets whether the claim value is of Number type. Port of isNumber().
func (c *ClaimWrapper) IsNumber() bool { return c.wrapped.Number != nil }

// Boolean gets the value as boolean. If the value is not of a boolean type, returns nil. Port
// of getBoolean().
func (c *ClaimWrapper) Boolean() *bool { return c.wrapped.Boolean }

// IsBoolean gets whether the claim value is of Boolean type. Port of isBoolean().
func (c *ClaimWrapper) IsBoolean() bool { return c.wrapped.Boolean != nil }

// Binary gets the binary value. If the value is not of a binary type, returns nil. Port of
// getBinary().
func (c *ClaimWrapper) Binary() []byte {
	if c.wrapped.Binary != nil {
		return []byte(*c.wrapped.Binary)
	}
	return nil
}

// IsBinary gets whether the claim value is of Binary type. Port of isBinary().
func (c *ClaimWrapper) IsBinary() bool { return c.wrapped.Binary != nil }

// DateTime gets the value as a date. If the value is not of a date type, returns nil. Port of
// getDateTime().
func (c *ClaimWrapper) DateTime() *time.Time {
	if c.wrapped.DateTime != nil {
		t := c.wrapped.DateTime.Time()
		return &t
	}
	return nil
}

// IsDateTime gets whether the claim value is of Date type. Port of isDateTime().
func (c *ClaimWrapper) IsDateTime() bool { return c.wrapped.DateTime != nil }

// List gets the value as a list. If the value is not of a list type, returns nil. Port of
// getList().
func (c *ClaimWrapper) List() []*ClaimWrapper {
	if c.listOverride != nil {
		return c.listOverride
	}
	if !c.IsList() {
		return nil
	}
	result := make([]*ClaimWrapper, 0, len(c.wrapped.Item))
	for _, item := range c.wrapped.Item {
		result = append(result, NewClaimWrapperWithParent(item, c))
	}
	return result
}

// IsList gets whether the claim value is of a list type. Port of isList().
func (c *ClaimWrapper) IsList() bool {
	if c.listOverride != nil {
		return true
	}
	return len(c.wrapped.Item) != 0
}

// Map gets the value as a map. If the value is not of a map type, returns nil. Port of
// getMap().
func (c *ClaimWrapper) Map() map[string]*ClaimWrapper {
	if c.mapOverride != nil {
		return c.mapOverride
	}
	if !c.IsMap() {
		return nil
	}
	result := make(map[string]*ClaimWrapper, len(c.wrapped.Entry))
	for _, entry := range c.wrapped.Entry {
		var name string
		if entry.Name != nil {
			name = *entry.Name
		}
		result[name] = NewClaimWrapperWithParent(entry, c)
	}
	return result
}

// IsMap gets whether the claim value is of a map type. Port of isMap().
func (c *ClaimWrapper) IsMap() bool {
	if c.isMapOverride != nil {
		return *c.isMapOverride
	}
	return len(c.wrapped.Entry) != 0
}

// IsNull gets whether the claim is of no known type. Port of isNull().
func (c *ClaimWrapper) IsNull() bool {
	o := c.claimOverrides()
	return !c.IsText() && !c.IsNumber() && !c.IsBoolean() && !c.IsBinary() && !c.IsDateTime() &&
		o.List() == nil && o.Map() == nil
}

// Wrapped gets the wrapped JAXB claim object. Port of getWrapped().
func (c *ClaimWrapper) Wrapped() *jaxb.XmlClaim { return c.wrapped }

// Parent gets parent claim, when present. Port of getParent().
func (c *ClaimWrapper) Parent() *ClaimWrapper { return c.parent }

// IsEmpty checks whether the claim is null or empty. Port of isEmpty(): the upstream method body
// - "isText() || isNumber() || isBoolean() || isDateTime() || isList() || isMap()" - returns true
// exactly when a value IS present, the reverse of what the method name and its Javadoc promise.
// The Java source itself flags this with a "// TODO : review" comment; reproduced verbatim here,
// inversion included, rather than "fixed".
func (c *ClaimWrapper) IsEmpty() bool {
	o := c.claimOverrides()
	return c.IsText() || c.IsNumber() || c.IsBoolean() || c.IsDateTime() || o.IsList() || o.IsMap()
}

// DisplayValue converts the claim's value to its corresponding string representation. Port of
// getDisplayValue().
func (c *ClaimWrapper) DisplayValue() string {
	o := c.claimOverrides()
	switch {
	case c.IsText():
		return c.Text()
	case c.IsNumber():
		return c.Number().String()
	case c.IsBoolean():
		if *c.Boolean() {
			return "true"
		}
		return "false"
	case c.IsBinary():
		return base64.StdEncoding.EncodeToString(c.Binary())
	case c.IsDateTime():
		return c.DateTime().UTC().Format(claimDateTimeFormat)
	case o.IsList():
		return claimListDisplayValue(o.List())
	case o.IsMap():
		return claimMapDisplayValue(o.Map())
	case c.IsNull():
		return "null"
	default:
		return ""
	}
}

// claimListDisplayValue is the private helper backing DisplayValue for a list-typed claim. Port
// of the private toDisplayValue(List<ClaimWrapper>).
func claimListDisplayValue(items []*ClaimWrapper) string {
	var sb strings.Builder
	for i, item := range items {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(item.DisplayValue())
	}
	return sb.String()
}

// claimMapDisplayValue is the private helper backing DisplayValue for a map-typed claim. Port of
// the private toDisplayValue(Map<String, ClaimWrapper>).
//
// The Java helper walks the java.util.HashMap that getMap() builds, so the order it renders the
// entries in is the HashMap's bucket order - arbitrary, but a deterministic function of the key
// set. Ranging over a Go map instead randomises the order on every call, which would make
// DisplayValue non-reproducible run to run. claimMapKeyOrder reproduces the HashMap's own order
// so this renders the same string Java does.
func claimMapDisplayValue(entries map[string]*ClaimWrapper) string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, key := range claimMapKeyOrder(entries) {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("\"")
		sb.WriteString(key)
		sb.WriteString("\": ")
		claimEmbedValueWithEnvelope(&sb, entries[key])
	}
	sb.WriteString("}")
	return sb.String()
}

// claimMapKeyOrder returns the map's keys in a stable order.
//
// Reproducing Java's own order exactly is not possible here. HashMap iterates by table bucket -
// the spread hash of the key masked to the table size, which a Go port can compute - but two
// keys landing in the same bucket come out in the order they were inserted, and that sequence is
// not recoverable from a Go map. It is not recoverable from the wrapped claim either: a subtype
// override builds its result as new HashMap<>(super.getMap()) plus further puts, so the
// insertion sequence is itself a previous HashMap's bucket order rather than the document order
// of the XmlClaim entries. Collisions do occur at these sizes - "country"/"street_address" share
// a bucket in a 16-entry table, for one - so a bucket-ordered rendering would still differ from
// Java's, only less visibly.
//
// The rendering therefore sorts by key: a documented, stable order that differs from Java's
// arbitrary-but-fixed one. Nothing in the port depends on the order - DisplayValue is a
// diagnostic string and is not part of the marshal-parity surface - whereas ranging over the Go
// map, which is what this replaced, made the string differ from one call to the next.
func claimMapKeyOrder(entries map[string]*ClaimWrapper) []string {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// claimEmbedValueWithEnvelope is the private helper backing claimMapDisplayValue. Port of the
// private embedValueWithEnvelope(StringBuilder, ClaimWrapper).
func claimEmbedValueWithEnvelope(sb *strings.Builder, claim *ClaimWrapper) {
	if claim.IsText() || claim.IsDateTime() {
		sb.WriteString("\"")
	} else if claim.IsList() {
		sb.WriteString("[")
	}
	sb.WriteString(claim.DisplayValue())
	if claim.IsText() || claim.IsDateTime() {
		sb.WriteString("\"")
	} else if claim.IsList() {
		sb.WriteString("]")
	}
}

// Equals reports whether other wraps a claim with the same name, disclosability, and value.
// Port of equals(Object): the initial `this == o` / instanceof checks become the pointer-equal
// and nil checks below; hashCode() has no Go equivalent (nothing here keys a hash-based
// collection on a ClaimWrapper) and is dropped, matching the omission pattern documented
// elsewhere in this port (see model/certificate_token.go).
func (c *ClaimWrapper) Equals(other *ClaimWrapper) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if c.Name() != other.Name() {
		return false
	}
	if c.IsSelectivelyDisclosable() != other.IsSelectivelyDisclosable() {
		return false
	}
	if c.IsText() && other.IsText() {
		return c.Text() == other.Text()
	}
	if c.IsNumber() && other.IsNumber() {
		return claimBigIntEqual(c.Number(), other.Number())
	}
	if c.IsBoolean() && other.IsBoolean() {
		return claimBoolPtrEqual(c.Boolean(), other.Boolean())
	}
	if c.IsBinary() && other.IsBinary() {
		return bytes.Equal(c.Binary(), other.Binary())
	}
	if c.IsDateTime() && other.IsDateTime() {
		return c.DateTime().Equal(*other.DateTime())
	}
	if c.IsList() && other.IsList() {
		return claimListEqual(c.List(), other.List())
	}
	if c.IsMap() && other.IsMap() {
		return claimMapEqual(c.Map(), other.Map())
	}
	// If types differ or both have no value
	return c.IsEmpty() && other.IsEmpty()
}

func claimBigIntEqual(a, b *big.Int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Cmp(b) == 0
}

func claimBoolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func claimListEqual(a, b []*ClaimWrapper) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

func claimMapEqual(a, b map[string]*ClaimWrapper) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || !v.Equals(bv) {
			return false
		}
	}
	return true
}

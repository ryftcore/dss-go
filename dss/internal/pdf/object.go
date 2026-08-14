// The PDF object model. See DESIGN.md §2.2 and §4.1 (frozen with errors.go).
//
// Provenance: org.apache.pdfbox.cos.COS* of pdfbox 3.0.7. Two properties are
// load-bearing and are not negotiable:
//
//   - *Dict preserves insertion order, because COSDictionary is a LinkedHashMap
//     and PdfDict.list() (used by SingleDssDict.extractVRIs) exposes that order
//     to DSS. A bare map[Name]Object is a determinism bug.
//   - Real keeps its source text, so a real we parsed and echo back is written
//     byte-exactly without having to reproduce Java's Float.toString.

package pdf

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Object is the sum type of PDF object kinds: Null, Bool, Integer, Real, String,
// Name, Array, *Dict, *Stream, Ref.
type Object interface{ pdfObject() }

// Null is the PDF null object. A dangling indirect reference resolves to Null{},
// never to an error (DESIGN.md §2.7 O3).
type Null struct{}

// Bool is a PDF boolean.
type Bool bool

// Integer is a PDF integer. Literals exceeding int64 are clamped with a warning.
type Integer int64

// Real carries the literal it was parsed from; Raw is "" for values built in
// memory. See R8: the writer echoes Raw verbatim when it is set.
type Real struct {
	Val float64
	Raw string
}

// String is a PDF string. Hex records the source form and forces hex on output.
type String struct {
	Bytes []byte
	Hex   bool
}

// Name is a PDF name, stored already unescaped (no leading slash).
type Name string

// Array is a PDF array.
type Array []Object

// Ref is an indirect reference appearing as a value inside the object graph.
type Ref struct {
	Num int64
	Gen uint16
}

// ObjectKey identifies an indirect object slot. It is the analogue of upstream's
// PdfObjectKey / PdfBoxObjectKey. The zero value means "no key". It is a distinct
// type from Ref on purpose: Ref is a value inside the object graph, ObjectKey
// identifies a slot.
type ObjectKey struct {
	Num int64
	Gen uint16
}

func (Null) pdfObject()    {}
func (Bool) pdfObject()    {}
func (Integer) pdfObject() {}
func (Real) pdfObject()    {}
func (String) pdfObject()  {}
func (Name) pdfObject()    {}
func (Array) pdfObject()   {}
func (*Dict) pdfObject()   {}
func (*Stream) pdfObject() {}
func (Ref) pdfObject()     {}

// IsZero reports whether k is the zero key.
func (k ObjectKey) IsZero() bool { return k.Num == 0 && k.Gen == 0 }

// String renders the key as upstream's PdfObjectKey.toString does: "12 0".
func (k ObjectKey) String() string {
	return strconv.FormatInt(k.Num, 10) + " " + strconv.FormatUint(uint64(k.Gen), 10)
}

// Key returns the slot the reference points at.
func (r Ref) Key() ObjectKey { return ObjectKey{Num: r.Num, Gen: r.Gen} }

// Dict is an insertion-ordered PDF dictionary. Order is part of the contract.
type Dict struct {
	keys []Name
	m    map[Name]int
	vals []Object
}

// NewDict returns an empty dictionary.
func NewDict() *Dict {
	return &Dict{m: make(map[Name]int)}
}

// DictOf builds a dictionary from alternating Name, Object pairs. It panics on a
// malformed pair — it is a construction helper for tests and for the writer, not
// a parser entry point.
func DictOf(kv ...any) *Dict {
	if len(kv)%2 != 0 {
		panic("pdf: DictOf needs an even number of arguments")
	}
	d := NewDict()
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(Name)
		if !ok {
			panic(fmt.Sprintf("pdf: DictOf key %d is %T, want pdf.Name", i, kv[i]))
		}
		v, ok := kv[i+1].(Object)
		if !ok {
			panic(fmt.Sprintf("pdf: DictOf value for /%s is %T, want pdf.Object", k, kv[i+1]))
		}
		d.Set(k, v)
	}
	return d
}

// Len reports the number of entries.
func (d *Dict) Len() int {
	if d == nil {
		return 0
	}
	return len(d.keys)
}

// Keys returns the keys in insertion order, as a fresh slice.
func (d *Dict) Keys() []Name {
	if d == nil {
		return nil
	}
	out := make([]Name, len(d.keys))
	copy(out, d.keys)
	return out
}

// Has reports whether key is present.
func (d *Dict) Has(key Name) bool {
	if d == nil || d.m == nil {
		return false
	}
	_, ok := d.m[key]
	return ok
}

// GetRaw returns the stored, unresolved value, or nil when the key is absent.
func (d *Dict) GetRaw(key Name) Object {
	if d == nil || d.m == nil {
		return nil
	}
	i, ok := d.m[key]
	if !ok {
		return nil
	}
	return d.vals[i]
}

// Set stores v under key. An existing key is updated in place and keeps its
// position, matching LinkedHashMap.put.
func (d *Dict) Set(key Name, v Object) {
	if d.m == nil {
		d.m = make(map[Name]int)
	}
	if i, ok := d.m[key]; ok {
		d.vals[i] = v
		return
	}
	d.m[key] = len(d.keys)
	d.keys = append(d.keys, key)
	d.vals = append(d.vals, v)
}

// Delete removes key, preserving the order of the remaining entries.
func (d *Dict) Delete(key Name) {
	if d == nil || d.m == nil {
		return
	}
	i, ok := d.m[key]
	if !ok {
		return
	}
	d.keys = append(d.keys[:i], d.keys[i+1:]...)
	d.vals = append(d.vals[:i], d.vals[i+1:]...)
	delete(d.m, key)
	for j := i; j < len(d.keys); j++ {
		d.m[d.keys[j]] = j
	}
}

// Clone returns a shallow copy: the value objects are shared.
func (d *Dict) Clone() *Dict {
	if d == nil {
		return nil
	}
	c := &Dict{
		keys: make([]Name, len(d.keys)),
		vals: make([]Object, len(d.vals)),
		m:    make(map[Name]int, len(d.m)),
	}
	copy(c.keys, d.keys)
	copy(c.vals, d.vals)
	for k, v := range d.m {
		c.m[k] = v
	}
	return c
}

// String renders the dictionary for diagnostics. It is not the writer's output
// format; see writer.go for that.
func (d *Dict) String() string {
	if d == nil {
		return "<<nil>>"
	}
	var b strings.Builder
	b.WriteString("<<")
	for i, k := range d.keys {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteByte('/')
		b.WriteString(string(k))
		b.WriteByte(' ')
		b.WriteString(objectString(d.vals[i]))
	}
	b.WriteString(">>")
	return b.String()
}

func objectString(o Object) string {
	switch v := o.(type) {
	case nil:
		return "<nil>"
	case Null:
		return "null"
	case Bool:
		if v {
			return "true"
		}
		return "false"
	case Integer:
		return strconv.FormatInt(int64(v), 10)
	case Real:
		if v.Raw != "" {
			return v.Raw
		}
		return strconv.FormatFloat(v.Val, 'g', -1, 64)
	case String:
		return "(" + string(v.Bytes) + ")"
	case Name:
		return "/" + string(v)
	case Array:
		parts := make([]string, len(v))
		for i, e := range v {
			parts[i] = objectString(e)
		}
		return "[" + strings.Join(parts, " ") + "]"
	case *Dict:
		return v.String()
	case *Stream:
		return v.Dict.String() + "stream(" + strconv.Itoa(len(v.Raw)) + ")"
	case Ref:
		return v.Key().String() + " R"
	default:
		return fmt.Sprintf("%v", o)
	}
}

// Stream is a PDF stream: its dictionary plus its still-encoded, still-encrypted
// bytes.
type Stream struct {
	Dict *Dict
	Raw  []byte
	// Offset and Length locate Raw in the source file; both are 0 for streams
	// built in memory.
	Offset int64
	Length int64
	// key is the slot the stream was parsed from; it is what the decryption
	// routine needs and is not part of the exported surface.
	key ObjectKey
	// decrypted guards against decrypting one stream twice, the same job
	// SecurityHandler's IdentityHashMap of already-seen COSBase does.
	decrypted bool
}

// NewStream builds an in-memory stream. Offset and Length stay 0.
func NewStream(d *Dict, raw []byte) *Stream {
	if d == nil {
		d = NewDict()
	}
	return &Stream{Dict: d, Raw: raw}
}

// Rect is a PDF rectangle, normalised so Min <= Max on both axes.
type Rect struct{ MinX, MinY, MaxX, MaxY float64 }

// Width is MaxX-MinX.
func (r Rect) Width() float64 { return r.MaxX - r.MinX }

// Height is MaxY-MinY.
func (r Rect) Height() float64 { return r.MaxY - r.MinY }

// Array renders the rectangle as a four-element PDF array of in-memory reals.
func (r Rect) Array() Array {
	return Array{
		Real{Val: r.MinX}, Real{Val: r.MinY},
		Real{Val: r.MaxX}, Real{Val: r.MaxY},
	}
}

// RectFromArray reads a rectangle from a four-number array, normalising the
// corners. The array must already have had its elements resolved.
func RectFromArray(a Array) (Rect, bool) {
	if len(a) < 4 {
		return Rect{}, false
	}
	var v [4]float64
	for i := 0; i < 4; i++ {
		f, ok := numberValue(a[i])
		if !ok {
			return Rect{}, false
		}
		v[i] = f
	}
	r := Rect{MinX: v[0], MinY: v[1], MaxX: v[2], MaxY: v[3]}
	if r.MinX > r.MaxX {
		r.MinX, r.MaxX = r.MaxX, r.MinX
	}
	if r.MinY > r.MaxY {
		r.MinY, r.MaxY = r.MaxY, r.MinY
	}
	return r, true
}

// numberValue extracts a float from Integer or Real.
func numberValue(o Object) (float64, bool) {
	switch v := o.(type) {
	case Integer:
		return float64(v), true
	case Real:
		if math.IsNaN(v.Val) || math.IsInf(v.Val, 0) {
			return 0, true
		}
		return v.Val, true
	default:
		return 0, false
	}
}

// Ported from dss-enumerations/.../COSESignatureType.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// COSESignatureType identifies the context of the signature. E.g. a
// COSE_Sign or COSE_Sign1 signature.
type COSESignatureType string

const (
	// COSESignatureTypeCoseSign is for signatures using the
	// COSE_Signature structure.
	COSESignatureTypeCoseSign COSESignatureType = "COSE_SIGN"
	// COSESignatureTypeCoseSign1 is for signatures using the COSE_Sign1
	// structure.
	COSESignatureTypeCoseSign1 COSESignatureType = "COSE_SIGN1"
	// COSESignatureTypeCoseSignature is for signatures using the
	// COSE_Signature structure.
	COSESignatureTypeCoseSignature COSESignatureType = "COSE_SIGNATURE"
	// COSESignatureTypeCoseCounterSignature is for full
	// counter-signatures.
	COSESignatureTypeCoseCounterSignature COSESignatureType = "COSE_COUNTER_SIGNATURE"
	// COSESignatureTypeCoseCounterSignature0 is for abbreviated
	// counter-signatures0.
	COSESignatureTypeCoseCounterSignature0 COSESignatureType = "COSE_COUNTER_SIGNATURE0"
	// COSESignatureTypeCoseCounterSignatureV2 is for full
	// counter-signatures with other_fields present.
	COSESignatureTypeCoseCounterSignatureV2 COSESignatureType = "COSE_COUNTER_SIGNATURE_V2"
	// COSESignatureTypeCoseCounterSignature0V2 is for abbreviated
	// counter-signatures0 with other_fields present.
	COSESignatureTypeCoseCounterSignature0V2 COSESignatureType = "COSE_COUNTER_SIGNATURE0_V2"
)

type coseSignatureTypeFields struct {
	context                      string
	hasContext                   bool
	label                        string
	hasLabel                     bool
	tag                          int64
	hasTag                       bool
	counterSignature             bool
	counterSignatureV2           bool
	counterSignatureHeaderKey    int64
	hasCounterSignatureHeaderKey bool
}

// coseSignatureTypeData holds the (context, label, tag, counterSignature,
// counterSignatureV2, counterSignatureHeaderKey) tuple for each constant.
// Java's null fields are represented with the corresponding has* boolean set
// to false.
var coseSignatureTypeData = map[COSESignatureType]coseSignatureTypeFields{
	COSESignatureTypeCoseSign: {
		context: "Signature", hasContext: true,
		label: "COSE_Sign", hasLabel: true,
		tag: 98, hasTag: true,
	},
	COSESignatureTypeCoseSign1: {
		context: "Signature1", hasContext: true,
		label: "COSE_Sign1", hasLabel: true,
		tag: 18, hasTag: true,
	},
	COSESignatureTypeCoseSignature: {},
	COSESignatureTypeCoseCounterSignature: {
		context: "CounterSignature", hasContext: true,
		label: "COSE_Countersignature", hasLabel: true,
		tag: 19, hasTag: true,
		counterSignature: true, counterSignatureV2: false,
		counterSignatureHeaderKey: 7, hasCounterSignatureHeaderKey: true,
	},
	COSESignatureTypeCoseCounterSignature0: {
		context: "CounterSignature0", hasContext: true,
		label: "COSE_Countersignature0", hasLabel: true,
		counterSignature: true, counterSignatureV2: false,
		counterSignatureHeaderKey: 9, hasCounterSignatureHeaderKey: true,
	},
	COSESignatureTypeCoseCounterSignatureV2: {
		context: "CounterSignatureV2", hasContext: true,
		label: "COSE_Countersignature_V2", hasLabel: true,
		tag: 19, hasTag: true,
		counterSignature: true, counterSignatureV2: true,
		counterSignatureHeaderKey: 11, hasCounterSignatureHeaderKey: true,
	},
	COSESignatureTypeCoseCounterSignature0V2: {
		context: "CounterSignature0V2", hasContext: true,
		label: "COSE_Countersignature0_V2", hasLabel: true,
		counterSignature: true, counterSignatureV2: true,
		counterSignatureHeaderKey: 12, hasCounterSignatureHeaderKey: true,
	},
}

// COSESignatureTypeValues returns all constants in declaration order.
func COSESignatureTypeValues() []COSESignatureType {
	return []COSESignatureType{
		COSESignatureTypeCoseSign,
		COSESignatureTypeCoseSign1,
		COSESignatureTypeCoseSignature,
		COSESignatureTypeCoseCounterSignature,
		COSESignatureTypeCoseCounterSignature0,
		COSESignatureTypeCoseCounterSignatureV2,
		COSESignatureTypeCoseCounterSignature0V2,
	}
}

// Context returns the context text string identifying the context of the
// signature. The value is used for DTBS computation. Returns "" if not set
// (mirrors Java's null).
func (c COSESignatureType) Context() string {
	return coseSignatureTypeData[c].context
}

// Label returns the RFC definition user-friendly label. Returns "" if not
// set (mirrors Java's null).
func (c COSESignatureType) Label() string {
	return coseSignatureTypeData[c].label
}

// Tag returns the tag of the corresponding signature structure. It panics if
// the tag is not available for this constant, mirroring Java's
// Objects.requireNonNull check in getTag().
func (c COSESignatureType) Tag() int64 {
	fields := coseSignatureTypeData[c]
	if !fields.hasTag {
		panic(fmt.Sprintf("The tag is not available for COSESignatureType '%s'", string(c)))
	}
	return fields.tag
}

// IsCounterSignature returns if the context corresponds to a counter
// signature.
func (c COSESignatureType) IsCounterSignature() bool {
	return coseSignatureTypeData[c].counterSignature
}

// IsCounterSignatureV2 returns if the context corresponds to an RFC 9338
// counter signature V2.
func (c COSESignatureType) IsCounterSignatureV2() bool {
	return coseSignatureTypeData[c].counterSignatureV2
}

// CounterSignatureHeaderKey returns the header key of unsigned property used
// to embed the counter signature in, and whether it is set (mirrors Java's
// nullable Long return).
func (c COSESignatureType) CounterSignatureHeaderKey() (int64, bool) {
	fields := coseSignatureTypeData[c]
	return fields.counterSignatureHeaderKey, fields.hasCounterSignatureHeaderKey
}

// COSESignatureTypeForLabel gets the corresponding COSESignatureType for the
// given label, if found. Panics if label is empty, mirroring Java's
// Objects.requireNonNull.
func COSESignatureTypeForLabel(label string) COSESignatureType {
	for _, v := range COSESignatureTypeValues() {
		fields := coseSignatureTypeData[v]
		if fields.hasLabel && fields.label == label {
			return v
		}
	}
	return ""
}

// COSESignatureTypeCounterSignatureContextByHeaderKey returns a
// corresponding counter signature context based on the used header
// identifier. Returns "" (zero value) when the header key is unknown,
// mirroring Java's null return.
func COSESignatureTypeCounterSignatureContextByHeaderKey(headerKey int64) COSESignatureType {
	for _, v := range COSESignatureTypeValues() {
		fields := coseSignatureTypeData[v]
		if fields.hasCounterSignatureHeaderKey && fields.counterSignatureHeaderKey == headerKey {
			return v
		}
	}
	return ""
}

// Ported from dss-enumerations/.../TSLType.java (DSS 6.5.RC1).
package enumerations

// TSLType defines a TSLType element of a Trusted List.
//
// NOTE: TSLTypeEnum (the concrete enum implementing this interface, with its
// TSLTypeEnumValues() accessor) is defined outside this file's manifest and
// is assumed to exist per the porting brief.
type TSLType interface {
	ListType

	// Label gets label.
	Label() string
}

// tslType is a plain TSLType implementation backing the fallback branch of
// TSLTypeFromURI (Java's TSLType.fromUri anonymous class).
type tslType struct {
	uri string
}

func (t *tslType) URI() string   { return t.uri }
func (t *tslType) Label() string { return "" }

// TSLTypeFromURI returns a TSLType for the given URI. If uri does not match
// any known TSLTypeEnum constant, a TSLType wrapping the given uri with an
// empty Label is returned (matching Java's anonymous-class fallback, whose
// getLabel() returned null).
func TSLTypeFromURI(uri string) TSLType {
	for _, t := range TSLTypeEnumValues() {
		if t.URI() == uri {
			return t
		}
	}
	return &tslType{uri: uri}
}

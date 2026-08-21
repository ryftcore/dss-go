// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/AbstractOtherTSLPointerPredicate.java (DSS 6.5.RC1).
//
// Java models this as an abstract class contributing one protected helper method,
// extractAdditionalInformation. Go has no abstract classes; the helper becomes a package-level
// function, called directly by TypeOtherTSLPointer, MimetypeOtherTSLPointer and
// SchemeTerritoryOtherTSLPointer instead of through embedding - the "one Go file per Java class"
// rule is honored by keeping the function in its own file, matching the file it was ported from.
package tsl

import (
	"encoding/xml"

	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
)

// extractAdditionalInformation extracts the additional information map from an
// OtherTSLPointerType. Port of extractAdditionalInformation(OtherTSLPointerType).
//
// Each entry's key mirrors Java's javax.xml.namespace.QName#toString(): "{namespaceURI}localPart"
// for both the recognized-element branch (Elem, JAXBElement in Java) and the raw-element branch
// (Raw, org.w3c.dom.Element in Java). The value is the element's text content: the recognized
// wildcard elements this port's predicates ever look up (TSLType, SchemeTerritory, MimeType) all
// decode to a Go string (see dss/trustedlist/jaxb's wildcardElements table), matching Java's
// JAXBElement<String>#getValue(); a raw (unrecognized) element's text content is approximated by
// concatenating its own top-level character-data tokens, which is exact for every element these
// predicates target (they carry no nested markup) though not a full DOM getTextContent()
// (recursive descendant-text concatenation).
func extractAdditionalInformation(o *jaxb.OtherTSLPointerType) map[string]any {
	result := make(map[string]any)

	additionalInformation := o.AdditionalInformation
	if additionalInformation == nil {
		return result
	}
	for _, item := range additionalInformation.Items {
		otherInformation := item.OtherInformation
		if otherInformation == nil {
			continue
		}
		for _, content := range otherInformation.Items {
			switch {
			case content.Elem != nil:
				key := "{" + content.ElemName.Space + "}" + content.ElemName.Local
				if s, ok := content.Elem.(*string); ok {
					result[key] = *s
				} else {
					result[key] = content.Elem
				}
			case content.Raw != nil:
				key := "{" + content.Raw.Name.Space + "}" + content.Raw.Name.Local
				result[key] = rawWildcardElementText(content.Raw)
			}
		}
	}
	return result
}

// rawWildcardElementText concatenates a RawWildcardElement's top-level character-data tokens,
// approximating org.w3c.dom.Element#getTextContent() - see extractAdditionalInformation's header.
func rawWildcardElementText(raw *jaxb.RawWildcardElement) string {
	var text string
	for _, tok := range raw.Tokens {
		if cd, ok := tok.(xml.CharData); ok {
			text += string(cd)
		}
	}
	return text
}

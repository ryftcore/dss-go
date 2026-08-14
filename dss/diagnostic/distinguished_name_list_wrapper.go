// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/DistinguishedNameListWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// DistinguishedNameListWrapper wraps a list of jaxb.XmlDistinguishedName.
type DistinguishedNameListWrapper struct {
	// xmlDistinguishedNames are the distinguished names.
	xmlDistinguishedNames []*jaxb.XmlDistinguishedName
}

// NewDistinguishedNameListWrapper is the default constructor.
func NewDistinguishedNameListWrapper(xmlDistinguishedNames []*jaxb.XmlDistinguishedName) *DistinguishedNameListWrapper {
	return &DistinguishedNameListWrapper{xmlDistinguishedNames: xmlDistinguishedNames}
}

// Value returns a value according to the given format. Port of getValue(String).
func (w *DistinguishedNameListWrapper) Value(format string) string {
	if w.xmlDistinguishedNames != nil {
		for _, distinguishedName := range w.xmlDistinguishedNames {
			if distinguishedName.Format == format {
				return distinguishedName.Value
			}
		}
	}
	return ""
}

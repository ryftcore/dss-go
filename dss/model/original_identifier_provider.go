// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/OriginalIdentifierProvider.java (DSS 6.5.RC1).
package model

// OriginalIdentifierProvider returns the original hash-based calculated String identifier
// for the given token.
type OriginalIdentifierProvider struct{}

// NewOriginalIdentifierProvider creates the provider.
func NewOriginalIdentifierProvider() *OriginalIdentifierProvider {
	return &OriginalIdentifierProvider{}
}

// IDAsString returns the XML Id of the object's identifier. Port of getIdAsString().
func (p *OriginalIdentifierProvider) IDAsString(object IdentifierBasedObject) string {
	return object.DSSID().AsXmlID()
}

// compile-time interface assertion.
var _ TokenIdentifierProvider = (*OriginalIdentifierProvider)(nil)

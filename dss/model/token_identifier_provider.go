// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/TokenIdentifierProvider.java (DSS 6.5.RC1).
package model

// TokenIdentifierProvider generates a String identifier for a given token (an
// AdvancedSignature, a CertificateToken, ...). Implementations cache the calculated values
// and take care of duplicates.
type TokenIdentifierProvider interface {
	// IDAsString returns a String identifier for the given object. Port of getIdAsString().
	IDAsString(object IdentifierBasedObject) string
}

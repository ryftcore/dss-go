// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/IdentifierBasedObject.java (DSS 6.5.RC1).
package model

// IdentifierBasedObject defines an object having an identifier (e.g. AdvancedSignature, Token).
type IdentifierBasedObject interface {
	// DSSID returns the Identifier of the object. Port of getDSSId().
	//
	// Java narrows the return type in subclasses (Token#getDSSId returns a TokenIdentifier);
	// Go has no covariant returns, so callers that need the narrower type assert on the
	// result, e.g. tokenIdentifier, ok := token.DSSID().(*TokenIdentifier). Note that a
	// token hands out its embedded *TokenIdentifier, not the concrete subclass value, so an
	// assertion to *CertificateTokenIdentifier does not succeed; the concrete Java class is
	// still reflected in the identifier's String() and Equals().
	DSSID() Identifier
}

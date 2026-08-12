// Ported from dss-crl-parser-x509crl/src/main/java/eu/europa/esig/dss/crl/x509/impl/X509CRLValidity.java (DSS 6.5.RC1).
//
// DEVIATION: upstream declares this as a class X509CRLValidity extending CRLValidity, reached
// polymorphically only when the classpath selects dss-crl-parser-x509crl over
// dss-crl-parser-stream. This port has a single native implementation, so the extra field
// (the parsed CRL) lives directly on CRLValidity (crl_validity.go); this file keeps only the
// accessors, so the Java class this behaviour came from stays traceable.
package crlparser

import "crypto/x509"

// X509CRL returns the parsed CRL this validity was built from, nil until
// CRLUtilsBuildCRLValidity (or SetX509CRL) populates it. Port of getX509CRL().
func (v *CRLValidity) X509CRL() *x509.RevocationList {
	return v.x509CRL
}

// SetX509CRL sets the parsed CRL. Port of setX509CRL(X509CRL).
func (v *CRLValidity) SetX509CRL(x509CRL *x509.RevocationList) {
	v.x509CRL = x509CRL
}

// NOTE: X509CRLValidity#equals()/#hashCode() additionally compare the cached x509CRL field
// after delegating to CRLValidity's own equals()/hashCode(); CRLValidity.Equals (crl_validity.go)
// does not compare it here. java.security.cert.X509CRL has no meaningful equals() of its own
// (reference identity in every JCE provider this project has seen), so upstream's extra
// comparison only ever short-circuits true for `this == that`, a case CRLValidity.Equals
// already covers before comparing any field; folding it in would add a reflect-based
// *x509.RevocationList comparison for no behavioural gain.

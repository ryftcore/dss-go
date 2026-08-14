// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/revocation/Revocation.java (DSS 6.5.RC1).
//
// This package flattens the Java packages eu.europa.esig.dss.model.x509.revocation,
// .revocation.crl and .revocation.ocsp into a single Go package, because Go forbids the
// import cycles their one-to-one mapping would create.
package revocation

// Revocation represents revocation data. It is a marker interface, so every Go type
// satisfies it; it exists to keep the CRL/OCSP type parameters of the Java API meaningful.
type Revocation interface {
}

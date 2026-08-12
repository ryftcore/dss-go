// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateSourceEntity.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
package spi

// CertificateSourceEntity defines items of a certificate source, for instance certificates
// grouped by a public key.
//
// Java declares this interface with no members (it extends only Serializable, which has no
// Go counterpart); an entirely empty Go interface would be satisfied by any value, defeating
// its purpose as a marker for CertificateSource#GetEntities. isCertificateSourceEntity keeps
// it a marker restricted to types defined in this package - today only *equivalentCertificatesEntity.
type CertificateSourceEntity interface {
	isCertificateSourceEntity()
}

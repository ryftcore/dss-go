// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/AIASource.java (DSS 6.5.RC1).
package aia

import "github.com/ryftcore/dss-go/dss/model"

// AIASource allows loading of issuing certificates by defined AIA URI within a
// model.CertificateToken.
type AIASource interface {
	// CertificatesByAIA loads a set of CertificateTokens accessed by AIA URIs from the
	// provided certificateToken. Ports getCertificatesByAIA(CertificateToken); the returned
	// slice stands in for Java's Set<CertificateToken>, order-preserving and de-duplicated by
	// DSSIDAsString() like the LinkedHashSet implementations upstream construct.
	CertificatesByAIA(certificateToken *model.CertificateToken) []*model.CertificateToken
}

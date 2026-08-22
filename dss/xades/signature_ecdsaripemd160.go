// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/SignatureECDSARIPEMD160.java (DSS 6.5.RC1).
//
// Upstream subclasses Apache Santuario's SignatureECDSA solely to override engineGetURI() with
// the ECDSA/RIPEMD160 signature algorithm URI, so that Santuario's algorithm registry can be
// taught a URI it does not know by default. internal/xmldsig's SignedInfo verification does not
// go through a Santuario-style pluggable JCA algorithm registry (see its doc.go table): it looks
// up enumerations.SignatureAlgorithm by URI directly. That lookup already knows
// SignatureAlgorithmECDSARIPEMD160's URI (enumerations is frozen and verbatim per PORTING.md),
// so there is no registration hook for this file to port into - SignatureAlgorithmURI is the Go
// replacement for what SignatureECDSARIPEMD160 exists to teach Santuario.
package xades

import "github.com/ryftcore/dss-go/dss/enumerations"

// SignatureECDSARIPEMD160URI is the URI SignatureECDSARIPEMD160.engineGetURI() returns,
// exposed here for parity with upstream in case a caller still expects a symbol for it.
func SignatureECDSARIPEMD160URI() string {
	return enumerations.SignatureAlgorithmECDSARIPEMD160.URI()
}

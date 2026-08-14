// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/SignatureRSARIPEMD160AT.java
// (DSS 6.5.RC1).
//
// Upstream extends Apache Santuario's SignatureBaseRSA to register a non-standard signature
// algorithm URI ("http://www.w3.org/2001/04/xmldsig-more/rsa-ripemd160", used by some Austrian
// signature providers) with Santuario's JCEMapper/SignatureAlgorithm registries, so that
// XMLSignature#checkSignatureValue can resolve it. internal/xmldsig replaces the whole
// Santuario signature-algorithm-registry mechanism (see its doc.go's "Signature verification is
// not routed through spi.ContentVerifier" section): verifySignature dispatches on the
// ds:SignatureMethod URI directly against crypto/rsa, with no JCE-style provider registry to
// register into. XML_ID is kept here, verbatim, as the URI constant this port's own dispatch
// needs to recognise the same non-standard algorithm; internal/xmldsig's RSA verification path
// (sigalg.go) is expected to treat it as RSA-RIPEMD160 alongside the standard rsa-ripemd160
// registration, the same way JCEMapper.register wired both names to one JCE algorithm name
// upstream. No other member of SignatureRSARIPEMD160AT is ported: engineGetURI() has no
// caller once there is no JCE-style algorithm object to ask, and the inherited SignatureBaseRSA
// engineSign/engineVerify are internal/xmldsig's own RSA verification code.
package xades

// SignatureRSARIPEMD160AT_XML_ID is the non-standard RSA-RIPEMD160 signature algorithm URI used
// by some AT (Austrian) signature providers. Port of the public static final String XML_ID
// field.
const SignatureRSARIPEMD160AT_XML_ID = "http://www.w3.org/2001/04/xmldsig-more/rsa-ripemd160"

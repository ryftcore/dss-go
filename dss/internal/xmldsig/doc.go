// Package xmldsig implements the XML-DSig reference-processing core that upstream DSS obtains
// from Apache Santuario (xmlsec 3.0.6): URI resolution, the transform pipeline,
// canonicalization of the result, digest computation, ds:Manifest validation, and verification
// of ds:SignatureValue over the canonicalized ds:SignedInfo.
//
// # Provenance
//
// Per PORTING.md, machinery that upstream gets from a third party has no Java class to mirror
// one-to-one and therefore lives under internal/. This package replaces exactly the part of
// org.apache.xml.security that dss-xades imports - the 25 import sites under
// dss-xades/src/main, enumerated in the table below - and nothing else. It is layered on
// internal/xmldom, internal/xmlc14n and internal/xpath10, and on model/enumerations for the
// DSSDocument and DigestAlgorithm types that its detached-content resolver and its digests
// are expressed in.
//
// The Go names below are the entry points; the Java column is the Santuario member each one
// replaces, and every exported identifier repeats its own mapping in its doc comment.
//
//	Go                                    Apache Santuario 3.0.6
//	------------------------------------------------------------------------------------
//	Data                                  signature.XMLSignatureInput
//	Data.Bytes                            XMLSignatureInput#getBytes
//	NodeFilter                            signature.NodeFilter (re-exported from xmlc14n)
//	Transform                             transforms.TransformSpi
//	Registry / DefaultRegistry            transforms.Transform's algorithm registry
//	Registry.Perform                      transforms.Transform#performTransform
//	PerformTransforms                     transforms.Transforms#performTransforms
//	URIResolver                           utils.resolver.ResourceResolverSpi
//	ResolverContext                       utils.resolver.ResourceResolverContext
//	Resolve                               utils.resolver.ResourceResolver#resolve
//	ResolverFragment                      utils.resolver.implementations.ResolverFragment
//	EnforcedResolverFragment              dss-xades EnforcedResolverFragment
//	ResolverXPointer                      utils.resolver.implementations.ResolverXPointer
//	DetachedSignatureResolver             dss-xades DetachedSignatureResolver
//	Reference                             signature.Reference
//	Reference.ContentsBeforeTransformation  Reference#getContentsBeforeTransformation
//	Reference.ContentsAfterTransformation   Reference#getContentsAfterTransformation
//	Reference.ReferencedBytes             Reference#getReferencedBytes
//	Reference.CalculateDigest             Reference#calculateDigest
//	Reference.Verify                      Reference#verify
//	Manifest                              signature.Manifest
//	Manifest.VerifyReferences             Manifest#verifyReferences
//	SignedInfo                            signature.SignedInfo
//	SignedInfo.CanonicalizedOctets        SignedInfo#getCanonicalizedOctetStream
//	XMLSignature                          signature.XMLSignature
//	XMLSignature.CheckSignatureValue      XMLSignature#checkSignatureValue
//	IsDescendantOrSelf, NodeSetOf         utils.XMLUtils#isDescendantOrSelf, #getSet
//
// # Two deliberate structural deviations
//
// No streaming. Santuario threads an OutputStream through the last transform of a chain so
// that the digest can be computed without materializing the transform output
// (Reference#calculateDigest passes a DigesterOutputStream; Transforms#performTransforms hands
// it to the last Transform only). Every transform that accepts the stream writes to it exactly
// the octets it would otherwise have returned - TransformC14N and TransformBase64Decode are
// the only two that look at it, and both branches are the same bytes - and a chain whose last
// transform ignores it ends in XMLSignatureInput#updateOutputStream, which canonicalizes with
// Canonicalizer20010315OmitComments, the same canonicalizer getBytes uses. So the streamed and
// the buffered paths are byte-identical, and this package implements only the buffered one -
// which is also the path DSS itself reads through Reference#getReferencedBytes and
// getContentsAfterTransformation().getBytes(). Reference.CalculateDigest digests
// ReferencedBytes.
//
// Signature verification is not routed through spi.ContentVerifier. XML-DSig encodes an ECDSA
// signature as the raw IEEE P1363 pair r||s, not as the DER SEQUENCE crypto/x509 expects
// (Santuario converts in SignatureECDSA), and XMLSignature#checkSignatureValue must
// distinguish "the signature does not verify" (false) from "this algorithm or key cannot be
// used" (an exception), which spi.ContentVerifier collapses into one error. verifySignature
// therefore mirrors Santuario's SignatureBaseRSA/SignatureECDSA/SignatureEDDSA engineVerify
// directly on crypto/rsa, crypto/ecdsa and crypto/ed25519.
//
// # What is out of scope, and where it went instead
//
// Building a signature (Santuario's Transform/Transforms/Reference construction side) belongs
// to the XAdES port: this package reads a signature that exists. dss-xades's
// CounterSignatureResolver is likewise not here - it serializes a node with
// DomUtils.serializeNode and queries it with XPathUtils, both of which live in DSS packages
// that internal/ may not import - but URIResolver is an interface precisely so the XAdES port
// can register it. The XSLT transform is refused, exactly as Santuario refuses it under secure
// validation and as DSS never enables it.
package xmldsig

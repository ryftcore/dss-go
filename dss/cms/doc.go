// Package cms ports dss-cms (eu.europa.esig.dss.cms), the public CMS layer DSS builds CAdES
// (and, in later phases, PAdES/JAdES) SignedData structures on top of.
//
// # Layering
//
// Upstream splits the module in three: dss-cms declares the CMS interface and the builders
// (Builder, SignerInfoGeneratorBuilder, ...), while dss-cms-object (an in-memory
// BouncyCastle-backed CMSSignedData) and dss-cms-stream (a streaming BouncyCastle-backed
// implementation, for large documents) are the two interchangeable ICMSUtils/Generator
// implementations an application picks between via ServiceLoader.
//
// Go has neither BouncyCastle nor a ServiceLoader, and internal/cmscore (see its own package
// doc) already replaces BouncyCastle's CMS/TSP engine with a byte-exact, dependency-free one.
// This package therefore has exactly ONE native implementation, built directly on
// internal/cmscore, rather than two interchangeable BC-backed ones: every exported type below
// (CMS, Builder, the Generator this package registers, SignerInfoGeneratorBuilder, ...)
// wraps cmscore's CMS/SignedData/SignerInfo builders instead of delegating to a pluggable
// ICMSUtils. AbstractCMSGenerator and CMSGenerator keep their Java shape (a settable "recipe"
// object with a Generate method) so a second implementation could still be added later, but
// there is no ServiceLoader indirection to select one - NewCMSGenerator is the only
// constructor, and the CMSGenerator.loadCMSGenerator() Java static factory has no port.
//
// DEVIATION (accepted, revisit later): dss-cms-object builds and holds the whole CMSSignedData
// in memory; dss-cms-stream additionally supports streaming a large to-be-signed document and
// writing the produced CMS through a caller-supplied DSSResourcesHandler rather than fully
// materialising it. This port's single implementation follows the dss-cms-object behaviour
// throughout: DSSResourcesHandlerBuilder is accepted for API parity (UtilsResourcesHandlerBuilder,
// CMSUtilsGetDSSResourcesHandlerBuilder) but is not consulted - every CMS document, signed or
// parsed, is read and built fully in memory, exactly as dss-cms-object's CMSObjectUtils does
// (its own writeToDSSDocument comment: "the 'dss-cms-object' implementation does not require
// using of resourcesHandlerBuilder"). A future phase may add real streaming without changing
// this package's public API.
//
// # BouncyCastle replacements
//
// Every BC type this package's Java sources import has a documented Go replacement, consistent
// with the rest of the port (see PORTING.md and internal/cmscore's doc.go):
//
//   - org.bouncycastle.cms.CMSSignedData / eu.europa.esig.dss.cms.CMS -> *cms.CMS, wrapping
//     *cmscore.CMS.
//   - org.bouncycastle.cms.SignerInformation                          -> *cmscore.SignerInfo.
//   - org.bouncycastle.cms.SignerInformationStore                     -> []*cmscore.SignerInfo.
//   - org.bouncycastle.asn1.cms.AttributeTable / Attribute             -> cmscore.Attributes /
//     *cmscore.Attribute.
//   - org.bouncycastle.util.Store<X509CertificateHolder/X509CRLHolder> -> [][]byte, the DER
//     encoding of each member (Store#getMatches(null) is then iterating the slice).
//   - org.bouncycastle.cms.SignerInfoGenerator                         -> *SignerInfoGenerator.
//   - org.bouncycastle.operator.ContentSigner                         -> the ContentSigner
//     interface below.
//   - org.bouncycastle.operator.DigestCalculatorProvider/DigestCalculator, whose two-level
//     stream-then-read indirection existed so BC could construct the calculator once and only
//     later learn which digest algorithm and bytes it applied to -> the single-method
//     DigestCalculatorProvider interface below: neither DSS provider (CustomMessageDigestCalculatorProvider,
//     PrecomputedDigestCalculatorProvider) actually depends on streamed bytes, so the collapse
//     changes nothing observable.
//   - org.bouncycastle.operator.DefaultSignatureAlgorithmIdentifierFinder, which CustomContentSigner
//     used to turn a JCE algorithm name into a signature AlgorithmIdentifier (OID plus, for
//     RSASSA-PSS, its explicit parameters) -> the signatureAlgorithmIdentifiers table in
//     custom_content_signer.go, generated from a real BouncyCastle 1.84 run (the exact version
//     DSS 6.5.RC1 depends on - see testdata/gen) rather than hand-derived, so that every
//     AlgorithmIdentifier this package writes into a SignerInfo.signatureAlgorithm field is
//     byte-identical to BouncyCastle's.
//
// # Byte-exactness
//
// SignedData/SignerInfo assembly itself is internal/cmscore's job and is byte-exact there
// (including DER SET-OF sorting and CMSVersion computation); this package's own responsibility
// is producing the three attributes BouncyCastle's DefaultAuthenticatedAttributeTableGenerator
// injects irrespective of what the caller's signed-attributes table already holds
// (content-type, message-digest, cms-algorithm-protection - see
// cms_signed_attribute_table_generator.go) and wiring the digest/signature algorithm
// identifiers that feed them.
package cms

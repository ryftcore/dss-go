// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSSignerInfoGeneratorBuilder.java
// (DSS 6.5.RC1).
//
// Java builds an org.bouncycastle.cms.SignerInfoGenerator, a "recipe" object BouncyCastle's
// CMSSignedDataGenerator.generate(content, encapsulate) later calls generate(contentType) on,
// once the real content-type and content bytes are known (see AbstractCMSGenerator's own
// comment and CMSGenerator.java). SignerInfoGenerator below is this port's replacement: a
// struct capturing the same recipe (SignerIdentifier, digest/signature algorithm identifiers,
// the not-yet-content-type/message-digest/algorithm-protection-completed signed attributes,
// the unsigned attributes, and the ContentSigner to ask for the signature), with a Generate
// method standing in for BouncyCastle's deferred generate(ASN1ObjectIdentifier contentType).
//
// # What is eager and what is deferred, and why
//
// Unlike BouncyCastle's generic SignerInfoGeneratorBuilder (reusable across arbitrary content),
// this port's Build already receives the DSSDocument to be signed, so DEVIATION: the
// message-digest value itself (which never depends on contentType) is computed eagerly here,
// at Build time, exactly as Java's own getDigestCalculatorProvider(toSignDocument) already does
// eagerly (toSignDocument.getDigestValue(digestAlgorithm) runs synchronously inside build()).
// Only contentType-dependent work - the content-type attribute and the final DER SET the
// signature is computed over - is deferred to Generate, matching upstream's own deferral
// exactly (CMSSignedAttributeTableGenerator#createStandardAttributeTable(Map) is what BC calls
// from inside generate(contentType), i.e. after contentType is known).
//
// # SignerInfo.digestAlgorithm derivation
//
// BouncyCastle's SignerInfoGeneratorBuilder.build derives the SignerInfo.digestAlgorithm field
// from the ContentSigner's own signature algorithm identifier
// (DefaultDigestAlgorithmIdentifierFinder#find(sigAlgId)), not from the digestAlgorithm this
// builder was configured with - the two are independent inputs in BC's generic API. DEVIATION:
// DSS always keeps them consistent (SignatureParameters derives its SignatureAlgorithm
// from its DigestAlgorithm, see CMSForCAdESBuilderHelper), so this port derives
// SignerInfo.digestAlgorithm from the configured digestAlgorithm field directly when one was
// set, falling back to the digest algorithm the ContentSigner's signature algorithm implies
// (enumerations.SignatureAlgorithm.DigestAlgorithm(), the DSS-canonical mapping BC's own finder
// agrees with for every family DSS produces) only when it was not - reproducing the derive-
// from-signature-algorithm behaviour for the one caller (build(ContentSigner) with no document)
// that can reach it without ever needing a second BouncyCastle-derived lookup table.
package cms

import (
	"encoding/asn1"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// SignerInfoGenerator is the not-yet-finalised recipe for a SignerInfo, replacing
// org.bouncycastle.cms.SignerInfoGenerator. AbstractCMSGenerator's native Generate (see
// abstract_cms_generator.go) is the only caller of Generate for a top-level CMS signature; a
// counter-signature builder (a later phase) calls it directly with a nil contentType.
type SignerInfoGenerator struct {
	// sid identifies the signing certificate.
	sid *cmscore.SignerIdentifier
	// digestAlgorithm is the digest algorithm identifier, i.e. what SignerInfoGenerator#getDigestAlgorithm
	// answers and what Builder unions into SignedData.digestAlgorithms.
	digestAlgorithm *asn1ber.AlgorithmIdentifier
	// messageDigest is the pre-computed message-digest attribute value.
	messageDigest []byte
	// signedAttributes holds the caller-supplied signed attributes, before content-type,
	// message-digest and cms-algorithm-protection are injected by Generate.
	signedAttributes cmscore.Attributes
	// unsignedAttributes holds the caller-supplied unsigned attributes, nil when there are none
	// (RFC 5652 5.3 forbids an empty-but-present unsignedAttrs field).
	unsignedAttributes cmscore.Attributes
	// contentSigner produces the signature over the final signed-attributes DER SET.
	contentSigner ContentSigner
}

// DigestAlgorithm returns the digest algorithm identifier. Port of BouncyCastle's
// SignerInfoGenerator#getDigestAlgorithm, which CMSBuilder needs to populate
// SignedData.digestAlgorithms.
func (g *SignerInfoGenerator) DigestAlgorithm() *asn1ber.AlgorithmIdentifier {
	return g.digestAlgorithm
}

// Generate produces the final SignerInfo for the given content type, nil for a
// counter-signature (Java: "contentType will be null if we're trying to generate a counter
// signature"). Port of BouncyCastle's deferred SignerInfoGenerator#generate(ASN1ObjectIdentifier),
// folding in what CMSSignedAttributeTableGenerator#getAttributes(Map) does once BC calls it
// from here.
//
// DEVIATION, flagged for the integrator (PORTING.md: internal/ is frozen, no direct edits): a
// SignerInfo.signature of zero length is exactly what the "data to sign" half of DSS's two-step
// signing produces (CustomContentSigner built with no pre-computed signature) and is
// syntactically legal CMS - BouncyCastle's ASN1OctetString does not reject an empty one, so
// upstream's CMSSignedDataGenerator.generate() completes normally and hands back a (never
// consulted - the caller only reads contentSigner.getOutputStream() afterward, see
// CAdESService#getDataToSign) CMSSignedData carrying that empty-signature SignerInfo.
// cmscore.SignerInfoBuilder.Build() instead rejects a zero-length Signature outright ("the
// SignerInfo builder needs a signature"). Rather than construct a SignerInfo cmscore refuses to
// build, Generate still writes the signed-attributes DER SET to the ContentSigner's
// OutputStream (the only side effect getDataToSign needs) and returns (nil, nil) - "evaluated,
// but no real SignerInfo exists yet" - which AbstractCMSGenerator.Generate propagates the same
// way (a nil, no-error CMS), matching every known call site's actual usage: CAdES discards the
// CMS this path produces and checks only that no error occurred.
func (g *SignerInfoGenerator) Generate(contentType asn1.ObjectIdentifier) (*cmscore.SignerInfo, error) {
	signedAttributes := cmsSignedAttributeTableGenerate(g.signedAttributes, contentType,
		g.messageDigest, g.digestAlgorithm, g.contentSigner.AlgorithmIdentifier())

	if _, err := g.contentSigner.OutputStream().Write(signedAttributes.DERSetEncoded()); err != nil {
		return nil, err
	}

	if len(g.contentSigner.Signature()) == 0 {
		return nil, nil
	}

	unsignedAttributes := g.unsignedAttributes
	if len(unsignedAttributes) == 0 {
		unsignedAttributes = nil
	}

	builder := &cmscore.SignerInfoBuilder{
		SID:                g.sid,
		DigestAlgorithm:    g.digestAlgorithm,
		SignedAttributes:   signedAttributes,
		SignatureAlgorithm: g.contentSigner.AlgorithmIdentifier(),
		Signature:          g.contentSigner.Signature(),
		UnsignedAttributes: unsignedAttributes,
	}
	return builder.Build()
}

// SignerInfoGeneratorBuilder is used to build a SignerInfoGenerator.
// Port of the CMSSignerInfoGeneratorBuilder class.
type SignerInfoGeneratorBuilder struct {
	// signingCertificate is the signing-certificate of the signer.
	signingCertificate *model.CertificateToken
	// digestAlgorithm is the digest algorithm to be used on message-digest computation.
	digestAlgorithm enumerations.DigestAlgorithm
	// signedAttributes are the attributes to be signed.
	signedAttributes cmscore.Attributes
	// unsignedAttributes are the unsigned attributes.
	unsignedAttributes cmscore.Attributes
}

// NewCMSSignerInfoGeneratorBuilder is the default constructor. Port of the no-arg constructor.
func NewCMSSignerInfoGeneratorBuilder() *SignerInfoGeneratorBuilder {
	return &SignerInfoGeneratorBuilder{}
}

// SetSigningCertificate sets the signing-certificate of the signer. Port of #setSigningCertificate.
func (b *SignerInfoGeneratorBuilder) SetSigningCertificate(signingCertificate *model.CertificateToken) *SignerInfoGeneratorBuilder {
	b.signingCertificate = signingCertificate
	return b
}

// SetDigestAlgorithm sets the Digest Algorithm to be used on message-digest computation.
// Port of #setDigestAlgorithm.
func (b *SignerInfoGeneratorBuilder) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *SignerInfoGeneratorBuilder {
	b.digestAlgorithm = digestAlgorithm
	return b
}

// SetSignedAttributes sets the signed attributes. Port of #setSignedAttributes.
func (b *SignerInfoGeneratorBuilder) SetSignedAttributes(signedAttributes cmscore.Attributes) *SignerInfoGeneratorBuilder {
	b.signedAttributes = signedAttributes
	return b
}

// SetUnsignedAttributes sets the unsigned attributes. Port of #setUnsignedAttributes.
func (b *SignerInfoGeneratorBuilder) SetUnsignedAttributes(unsignedAttributes cmscore.Attributes) *SignerInfoGeneratorBuilder {
	b.unsignedAttributes = unsignedAttributes
	return b
}

// BuildWithoutDocument builds a SignerInfoGenerator with no original document provided. Port of
// #build(ContentSigner), i.e. build(null, contentSigner); Go has no method overloading, so the
// single-argument Java overload gets this distinct name (Build is the two-argument
// #build(DSSDocument, ContentSigner), DSS's only call site).
//
// This overload only behaves usefully when no digestAlgorithm was configured (see the file
// header): with one configured, Java's toSignDocument.getDigestValue(digestAlgorithm) would
// raise a NullPointerException against the null toSignDocument, which this port reproduces as
// a panic, matching Java's own uncaught-exception behaviour for that programmer error.
func (b *SignerInfoGeneratorBuilder) BuildWithoutDocument(contentSigner ContentSigner) (*SignerInfoGenerator, error) {
	return b.build(nil, contentSigner)
}

// Build builds a SignerInfoGenerator for signing toSignDocument.
// Port of #build(DSSDocument, ContentSigner).
func (b *SignerInfoGeneratorBuilder) Build(toSignDocument model.DSSDocument, contentSigner ContentSigner) (*SignerInfoGenerator, error) {
	return b.build(toSignDocument, contentSigner)
}

// build ports the shared body of both build() overloads.
func (b *SignerInfoGeneratorBuilder) build(toSignDocument model.DSSDocument, contentSigner ContentSigner) (*SignerInfoGenerator, error) {
	digestAlgorithm, err := b.resolveDigestAlgorithm(contentSigner)
	if err != nil {
		return nil, err
	}

	messageDigest, err := b.messageDigest(toSignDocument, digestAlgorithm)
	if err != nil {
		return nil, err
	}

	digestAlgorithmIdentifier, err := spi.DSSASN1UtilsAlgorithmIdentifierForDigest(digestAlgorithm)
	if err != nil {
		return nil, err
	}

	signedAttributes := b.signedAttributes
	if len(signedAttributes) == 0 {
		signedAttributes = nil
	}
	unsignedAttributes := b.unsignedAttributes
	if len(unsignedAttributes) == 0 {
		unsignedAttributes = nil
	}

	sid, err := b.signerIdentifier()
	if err != nil {
		return nil, err
	}

	return &SignerInfoGenerator{
		sid:                sid,
		digestAlgorithm:    digestAlgorithmIdentifier,
		messageDigest:      messageDigest,
		signedAttributes:   signedAttributes,
		unsignedAttributes: unsignedAttributes,
		contentSigner:      contentSigner,
	}, nil
}

// resolveDigestAlgorithm returns the digest algorithm the built SignerInfo's digestAlgorithm
// field (and the message-digest computation) uses: b.digestAlgorithm when one was configured,
// otherwise the digest algorithm the content signer's signature algorithm implies. See the
// file header ("SignerInfo.digestAlgorithm derivation").
func (b *SignerInfoGeneratorBuilder) resolveDigestAlgorithm(contentSigner ContentSigner) (enumerations.DigestAlgorithm, error) {
	if b.digestAlgorithm != "" {
		return b.digestAlgorithm, nil
	}
	signatureAlgorithmIdentifier := contentSigner.AlgorithmIdentifier()
	signatureAlgorithm, err := enumerations.SignatureAlgorithmForOIDAndParams(
		signatureAlgorithmIdentifier.Algorithm.String(), signatureAlgorithmIdentifier.Parameters)
	if err != nil {
		return "", err
	}
	return signatureAlgorithm.DigestAlgorithm(), nil
}

// messageDigest computes the message-digest attribute value. Port of the combined effect of
// #getDigestCalculatorProvider(DSSDocument) and the DigestCalculator#getDigest() BouncyCastle
// calls once digestion completes:
//
//   - digestAlgorithm configured (CustomMessageDigestCalculatorProvider): the document's own
//     digest for that algorithm, propagating a DigestDocument's "no such algorithm" error like
//     Java's uncaught IllegalArgumentException does.
//   - not configured, toSignDocument a *model.DigestDocument (PrecomputedDigestCalculatorProvider):
//     the document's digest for the resolved algorithm, an empty slice on any failure (Java
//     logs a warning and falls back to DSSUtils.EMPTY_BYTE_ARRAY).
//   - otherwise (BcDigestCalculatorProvider, a genuine BouncyCastle class with no DSS port):
//     the real digest of the document's bytes for the resolved algorithm.
func (b *SignerInfoGeneratorBuilder) messageDigest(toSignDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	if b.digestAlgorithm != "" {
		return toSignDocument.DigestValue(b.digestAlgorithm)
	}
	if digestDocument, ok := toSignDocument.(*model.DigestDocument); ok {
		provider := NewPrecomputedDigestCalculatorProvider(digestDocument)
		digestAlgorithmIdentifier, err := spi.DSSASN1UtilsAlgorithmIdentifierForDigest(digestAlgorithm)
		if err != nil {
			return []byte{}, nil
		}
		return provider.Digest(digestAlgorithmIdentifier)
	}
	return toSignDocument.DigestValue(digestAlgorithm)
}

// signerIdentifier builds the SignerIdentifier the SignerInfo carries. Port of the SignerId
// selection inside the private #getSignerInfoGenerator: an issuerAndSerialNumber built from
// signingCertificate when one is set, otherwise the empty subjectKeyIdentifier of
// `new SignerId(DSSUtils.EMPTY_BYTE_ARRAY)` - used to generate data-to-be-signed without a
// signing certificate (CAdESSignatureParameters#isGenerateTBSWithoutCertificate).
func (b *SignerInfoGeneratorBuilder) signerIdentifier() (*cmscore.SignerIdentifier, error) {
	if b.signingCertificate == nil {
		return cmscore.NewSubjectKeyIdentifierSID([]byte{}), nil
	}
	certificate := b.signingCertificate.Certificate()
	return cmscore.NewIssuerAndSerialNumberSID(certificate.RawIssuer, b.signingCertificate.SerialNumber()), nil
}

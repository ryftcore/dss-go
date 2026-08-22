// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/AbstractCMSGenerator.java
// (DSS 6.5.RC1), plus - since this package has one native CMSGenerator rather than the
// dss-cms-object/dss-cms-stream pair Java's ServiceLoader chooses between (see doc.go) - the
// generation logic Java splits out into dss-cms-object's CMSObjectGenerator and CMSObjectUtils'
// populateDigestAlgorithmSet, folded into Generate below.
//
// # What Generate replaces
//
// org.bouncycastle.cms.CMSSignedDataGenerator#generate(CMSTypedData, boolean) - "sign the
// content, encapsulating it when asked to" - together with the addDigestAlgorithm
// post-processing CMSObjectGenerator#generate() applies via CMSUtils#populateDigestAlgorithmSet
// afterwards. Two structural choices this port makes explicitly rather than by construction:
//
//   - SignedData.digestAlgorithms is the union of every (old and new) SignerInfo's digest
//     algorithm - what cmscore.SignedDataBuilder derives automatically when left unset, exactly
//     mirroring BC's own generate(), which populates digestAlgs purely from
//     signerGen.getDigestAlgorithm() (new signers) and each old SignerInformation's
//     digestAlgorithmID - unioned with digestAlgorithmIDs (CMSBuilder's own explicit
//     addition, covering an original CMS's digest algorithms that no current SignerInfo
//     implies, e.g. after a signer was dropped). CMSObjectGenerator does this as two separate
//     steps (generate(), then CMSUtils.populateDigestAlgorithmSet); Generate does it in one.
//   - SignedData.certificates/crls ordering is DER-sorted (internal/cmscore's CertificateSet/
//     RevocationInfoChoices always sort), where BC's generator preserves insertion order via
//     BERSet. This is a known, deliberate cmscore deviation for CAdES signing:
//     the SignerInfo bytes it does not affect are what has to be byte-exact.
package cms

import (
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// AbstractCMSGenerator is this package's (only) CMSGenerator implementation, containing the set
// variable values. Port of the abstract AbstractCMSGenerator class plus - see the file header -
// the native generation logic dss-cms-object's CMSObjectGenerator supplies upstream.
type AbstractCMSGenerator struct {
	// signerInfoGenerator is the new signer to be generated.
	signerInfoGenerator *SignerInfoGenerator
	// certificateStore is the collection of certificates to be encapsulated within
	// SignedData.certificates field, each as its own DER encoding.
	certificateStore [][]byte
	// signers is the collection of existing signers to be added.
	signers []*cmscore.SignerInfo
	// attributeCertificates is the collection of attribute certificates, each already the
	// complete [2] IMPLICIT AttributeCertificateV2 encoding (see CMS#AttributeCertificates).
	attributeCertificates [][]byte
	// crls is the collection of CRLs to be encapsulated within SignedData.crls field, each as
	// its own DER encoding.
	crls [][]byte
	// ocspBasicStore is the collection of OCSP basic responses (id-pkix-ocsp-basic), each as
	// its own DER encoding.
	ocspBasicStore [][]byte
	// ocspResponsesStore is the collection of OCSP responses (id-ri-ocsp-response), each as its
	// own DER encoding.
	ocspResponsesStore [][]byte
	// digestAlgorithmIDs is the collection of digest algorithm IDs to be included.
	digestAlgorithmIDs []*asn1ber.AlgorithmIdentifier
	// toBeSignedDocument is the document to be signed.
	toBeSignedDocument model.DSSDocument
	// encapsulate defines whether the signed document shall be encapsulated within the CMS.
	encapsulate bool
}

// NewAbstractCMSGenerator is the default constructor. Port of the protected no-arg constructor;
// exported since Go has no ServiceLoader to hide it behind (see doc.go).
func NewAbstractCMSGenerator() *AbstractCMSGenerator {
	return &AbstractCMSGenerator{}
}

// SetSignerInfoGenerator ports #setSignerInfoGenerator.
func (g *AbstractCMSGenerator) SetSignerInfoGenerator(signerInfoGenerator *SignerInfoGenerator) {
	g.signerInfoGenerator = signerInfoGenerator
}

// SetCertificates ports #setCertificates.
func (g *AbstractCMSGenerator) SetCertificates(certificateStore [][]byte) {
	g.certificateStore = certificateStore
}

// SetSigners ports #setSigners.
func (g *AbstractCMSGenerator) SetSigners(signers []*cmscore.SignerInfo) { g.signers = signers }

// SetAttributeCertificates ports #setAttributeCertificates.
func (g *AbstractCMSGenerator) SetAttributeCertificates(attributeCertificates [][]byte) {
	g.attributeCertificates = attributeCertificates
}

// SetCRLs ports #setCRLs.
func (g *AbstractCMSGenerator) SetCRLs(crls [][]byte) { g.crls = crls }

// SetOcspBasicStore ports #setOcspBasicStore.
func (g *AbstractCMSGenerator) SetOcspBasicStore(ocspBasicStore [][]byte) {
	g.ocspBasicStore = ocspBasicStore
}

// SetOcspResponsesStore ports #setOcspResponsesStore.
func (g *AbstractCMSGenerator) SetOcspResponsesStore(ocspResponsesStore [][]byte) {
	g.ocspResponsesStore = ocspResponsesStore
}

// SetDigestAlgorithmIDs ports #setDigestAlgorithmIDs.
func (g *AbstractCMSGenerator) SetDigestAlgorithmIDs(digestAlgorithmIDs []*asn1ber.AlgorithmIdentifier) {
	g.digestAlgorithmIDs = digestAlgorithmIDs
}

// SetToBeSignedDocument ports #setToBeSignedDocument.
func (g *AbstractCMSGenerator) SetToBeSignedDocument(document model.DSSDocument) {
	g.toBeSignedDocument = document
}

// SetEncapsulate ports #setEncapsulate.
func (g *AbstractCMSGenerator) SetEncapsulate(encapsulate bool) { g.encapsulate = encapsulate }

// Generate generates the CMS. See the file header for what this replaces.
//
// Panics with the Java message when toBeSignedDocument is nil (DSSUtils.toByteArray/
// toCMSEncapsulatedContent's Objects.requireNonNull); a *model.DigestDocument asked to
// encapsulate is a returned error standing in for BouncyCastle's
// CMSAbsentContent#write UnsupportedOperationException, which - unlike the requireNonNull
// above - is a data-dependent condition (a caller-chosen SignaturePackaging combined with a
// DigestDocument) rather than a programmer error.
func (g *AbstractCMSGenerator) Generate() (*CMS, error) {
	if g.toBeSignedDocument == nil {
		panic("Document to be signed is missing")
	}

	contentBytes, err := g.encapsulatedContentBytes()
	if err != nil {
		return nil, err
	}
	encapContentInfo := cmscore.NewEncapsulatedContentInfo(cmscore.OIDData, contentBytes)

	newSignerInfo, err := g.signerInfoGenerator.Generate(cmscore.OIDData)
	if err != nil {
		return nil, err
	}
	if newSignerInfo == nil {
		// The signer-info recipe carries no signature yet (DSS's "data to sign" half of
		// two-step signing, see SignerInfoGenerator.Generate) - the bytes to be signed were
		// still written to the ContentSigner's OutputStream, but there is no real CMS to
		// build. Every known caller of this path (CAdESService#getDataToSign) discards the
		// returned CMS and checks only that no error occurred.
		return nil, nil
	}
	signerInfos := make([]*cmscore.SignerInfo, 0, len(g.signers)+1)
	signerInfos = append(signerInfos, newSignerInfo)
	signerInfos = append(signerInfos, g.signers...)

	builder := &cmscore.SignedDataBuilder{
		EncapContentInfo: encapContentInfo,
		Certificates:     certificateChoices(g.certificateStore, g.attributeCertificates),
		CRLs:             revocationInfoChoices(g.crls, g.ocspResponsesStore, g.ocspBasicStore),
		SignerInfos:      signerInfos,
	}
	signedData, err := builder.Build()
	if err != nil {
		return nil, err
	}
	signedData.DigestAlgorithms = mergeAlgorithmIdentifiers(signedData.DigestAlgorithms, g.digestAlgorithmIDs)

	return newCMS(cmscore.NewCMS(signedData)), nil
}

// encapsulatedContentBytes resolves the eContent octets, nil for a detached signature. Port of
// the content-handling half of CMSUtils#toCMSEncapsulatedContent(DSSDocument) (the
// DigestDocument/other split) together with CMSSignedDataGenerator#generate's own
// "encapsulate ? full bytes : null" choice.
func (g *AbstractCMSGenerator) encapsulatedContentBytes() ([]byte, error) {
	if !g.encapsulate {
		return nil, nil
	}
	if _, isDigestDocument := g.toBeSignedDocument.(*model.DigestDocument); isDigestDocument {
		return nil, model.NewDSSError("Unable to encapsulate a DigestDocument: its content is not available")
	}
	return spi.DSSUtilsToByteArrayOfDocument(g.toBeSignedDocument)
}

// certificateChoices builds the SignedData.certificates members: the plain certificates of
// certificateStore, followed by the already-tagged [2] IMPLICIT AttributeCertificateV2 members
// of attributeCertificates. Shared with CMSUtilsReplaceCertificatesAndCRLs, which rebuilds the
// same field from a different source.
func certificateChoices(certificateStore, attributeCertificates [][]byte) []cmscore.CertificateChoice {
	choices := make([]cmscore.CertificateChoice, 0, len(certificateStore)+len(attributeCertificates))
	for _, certificate := range certificateStore {
		choices = append(choices, cmscore.NewCertificateChoice(certificate))
	}
	for _, attributeCertificate := range attributeCertificates {
		choices = append(choices, cmscore.NewTaggedCertificateChoice(cmscore.CertificateChoiceV2AttrCert, attributeCertificate))
	}
	return choices
}

// revocationInfoChoices builds the SignedData.crls members: the CRLs of crls, the OCSP
// responses of ocspResponsesStore under id-ri-ocsp-response, and the OCSP basic responses of
// ocspBasicStore under id-pkix-ocsp-basic. Port of CMSObjectUtils#toCRLsStore. Shared with
// CMSUtilsReplaceCertificatesAndCRLs.
func revocationInfoChoices(crls, ocspResponsesStore, ocspBasicStore [][]byte) []cmscore.RevocationInfoChoice {
	choices := make([]cmscore.RevocationInfoChoice, 0, len(crls)+len(ocspResponsesStore)+len(ocspBasicStore))
	for _, crl := range crls {
		choices = append(choices, cmscore.NewCRLRevocationInfoChoice(crl))
	}
	for _, ocspResponse := range ocspResponsesStore {
		choices = append(choices, cmscore.NewOtherRevocationInfoChoice(
			cmscore.NewOtherRevocationInfoFormat(cmscore.OIDRIOCSPResponse, ocspResponse)))
	}
	for _, ocspBasicResponse := range ocspBasicStore {
		choices = append(choices, cmscore.NewOtherRevocationInfoChoice(
			cmscore.NewOtherRevocationInfoFormat(cmscore.OIDPKIXOCSPBasic, ocspBasicResponse)))
	}
	return choices
}

// mergeAlgorithmIdentifiers returns base with every member of extra that is not already present
// (by AlgorithmIdentifier.Equals) appended; the slice order does not affect the eventual DER,
// since SignedData.DER() always writes digestAlgorithms as a DER SET OF (sorted by encoding).
// Port of CMSObjectUtils#populateDigestAlgorithmSet's "if !contains: add" loop.
func mergeAlgorithmIdentifiers(base, extra []*asn1ber.AlgorithmIdentifier) []*asn1ber.AlgorithmIdentifier {
	merged := make([]*asn1ber.AlgorithmIdentifier, len(base), len(base)+len(extra))
	copy(merged, base)
	for _, candidate := range extra {
		if !containsAlgorithmIdentifier(merged, candidate) {
			merged = append(merged, candidate)
		}
	}
	return merged
}

// containsAlgorithmIdentifier reports whether identifiers holds an entry equal to candidate.
func containsAlgorithmIdentifier(identifiers []*asn1ber.AlgorithmIdentifier, candidate *asn1ber.AlgorithmIdentifier) bool {
	for _, identifier := range identifiers {
		if identifier.Equals(candidate) {
			return true
		}
	}
	return false
}

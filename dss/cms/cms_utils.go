// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSUtils.java and
// dss-cms/src/main/java/eu/europa/esig/dss/cms/ICMSUtils.java (DSS 6.5.RC1).
//
// Java splits CMSUtils (a static facade) from ICMSUtils (the ServiceLoader-selected interface
// it delegates every call to, implemented once per dss-cms-object/dss-cms-stream backend) so
// that application code never has to know which backend is on the classpath. This port has one
// native implementation and no ServiceLoader (see doc.go), so the two collapse into the single
// set of flattened CMSUtils<MethodName> functions below, following PORTING.md's convention for
// static-utility classes; ICMSUtils itself has no separate port; func literals implementing it
// would have added indirection with no second implementation to justify it.
//
// DEVIATION (in-memory only, accepted for now - see doc.go): every function that takes a
// DSSResourcesHandlerBuilder accepts it for API parity but never consults it - CMSUtilsWriteToDSSDocument
// always builds its result fully in memory, exactly as dss-cms-object's CMSObjectUtils does
// ("the 'dss-cms-object' implementation does not require using of resourcesHandlerBuilder");
// unlike dss-cms-object, this port does not reject a caller-supplied builder either
// (CMSUtilsResourcesHandlerBuilder is a pass-through, not dss-cms-object's
// UnsupportedOperationException), since accepting one costs nothing and keeps the door open for
// a later streaming implementation to honour it.
package cms

import (
	"bytes"
	"io"
	"sort"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
)

// CMSUtilsParseToCMS parses the given DSSDocument to a CMS object. Port of #parseToCMS(DSSDocument),
// specialised to CMSSignedDataObject's behaviour: a *CMSSignedDocument hands back the CMS it
// already wraps rather than being re-parsed from its bytes.
//
// A malformed document is an *exception.IllegalInputException, matching Java's own
// IllegalInputException("Not a valid CAdES file. ..."); an I/O failure reading document is a
// *model.DSSError, matching Java's DSSException.
func CMSUtilsParseToCMS(document model.DSSDocument) (*CMS, error) {
	if signedDocument, ok := document.(*CMSSignedDocument); ok {
		return signedDocument.CMSSignedData(), nil
	}
	binaries, err := spi.DSSUtilsToByteArrayOfDocument(document)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to read a document.", err)
	}
	return CMSUtilsParseToCMSBinaries(binaries)
}

// CMSUtilsParseToCMSBinaries parses the given byte array to a CMS object. Port of
// #parseToCMS(byte[]).
//
// This is a whole-document parse - binaries is everything a caller read from a file or
// InMemoryDocument, not a value already isolated to its own field - so it goes through
// cmscore.ParseCMSTolerateTrailingBytes rather than the plain, trailing-bytes-rejecting
// cmscore.ParseCMS: see that function's doc comment for why the two differ, and for the real
// fixture (DSS-1188) that depends on the tolerance.
func CMSUtilsParseToCMSBinaries(binaries []byte) (*CMS, error) {
	core, err := cmscore.ParseCMSTolerateTrailingBytes(binaries)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause("Not a valid CAdES file.", err)
	}
	return newCMS(core), nil
}

// CMSUtilsWriteToDSSDocument creates a DSSDocument from the given CMS. Port of
// #writeToDSSDocument(CMS, DSSResourcesHandlerBuilder); see the file header for
// resourcesHandlerBuilder.
func CMSUtilsWriteToDSSDocument(cms *CMS, resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) (model.DSSDocument, error) {
	return NewCMSSignedDocument(cms), nil
}

// cmsUtilsCertificatesAndCRLs returns signedData's certificates/CRLs as the plain slices
// cmscore.SignedDataBuilder takes, nil for either that is absent rather than panicking: both
// Certificates and CRLs are themselves pointers (SignedData.Certificates *CertificateSet,
// SignedData.CRLs *RevocationInfoChoices), nil exactly when RFC 5652's OPTIONAL field was not
// present, which every rebuild of a SignedData (CMSUtilsReplaceSigners,
// CMSUtilsPopulateDigestAlgorithmSet) has to carry forward without dereferencing.
func cmsUtilsCertificatesAndCRLs(signedData *cmscore.SignedData) (certificates []cmscore.CertificateChoice, crls []cmscore.RevocationInfoChoice) {
	if signedData.Certificates != nil {
		certificates = signedData.Certificates.Choices
	}
	if signedData.CRLs != nil {
		crls = signedData.CRLs.Choices
	}
	return certificates, crls
}

// CMSUtilsReplaceSigners replaces the signers within cms with newSignerStore. Port of
// #replaceSigners(CMS, SignerInformationStore).
//
// DEVIATION: BouncyCastle's CMSSignedData#replaceSigners keeps the original SignedData.version
// untouched and derives SignedData.digestAlgorithms purely from newSignerStore, dropping any
// digest algorithm the previous signerInfos alone implied. This port instead keeps the union of
// the two (see mergeAlgorithmIdentifiers): digestAlgorithms is not itself signed data, so
// widening rather than narrowing it cannot invalidate anything RFC 5652 requires, and CAdES's
// baseline usage (a single signer being augmented in place) never observes the difference.
//
// original.Certificates/CRLs are nil, not merely empty, when the corresponding SignedData field
// is absent - the certificates field in particular is OPTIONAL in RFC 5652 and a CAdES-T
// enveloping signature carrying no CRLs at all is exactly the ordinary case (revocation data
// only shows up from -LT on). Both are guarded here (via cmsUtilsCertificatesAndCRLs) rather
// than dereferenced directly.
func CMSUtilsReplaceSigners(cms *CMS, newSignerStore []*cmscore.SignerInfo) (*CMS, error) {
	original := cms.Core().SignedData()
	certificates, crls := cmsUtilsCertificatesAndCRLs(original)
	builder := &cmscore.SignedDataBuilder{
		DigestAlgorithms: mergeAlgorithmIdentifiers(original.DigestAlgorithms, digestAlgorithmsOfSigners(newSignerStore)),
		EncapContentInfo: original.EncapContentInfo,
		Certificates:     certificates,
		CRLs:             crls,
		SignerInfos:      newSignerStore,
	}
	signedData, err := builder.Build()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to replace signerInfo of CMS SignedData. Corrupted content has been provided.", err)
	}
	signedData.Version = original.Version
	return newCMS(cmscore.NewCMS(signedData)), nil
}

// CMSUtilsRecomputeSignerInformation looks up the SignerInfo of cmsDoc identified by signerId.
// Port of #recomputeSignerInformation(CMS, SignerId, DigestCalculatorProvider,
// DSSResourcesHandlerBuilder), specialised to dss-cms-object's implementation:
// `new CMSSignedDataParser(digestCalculatorProvider, cms.getDEREncoded()).getSignerInfos().get(signerId)`.
//
// DEVIATION: BC's CMSSignedDataParser re-parses the CMS bytes while streaming the
// digestCalculatorProvider's content through a digest calculator, so the returned
// SignerInformation carries a cached content-digest computed against whatever document the
// caller's digestCalculatorProvider wraps (CAdESSignature#recreateSignerInformation feeds it a
// PrecomputedDigestCalculatorProvider(detachedContents.get(0)) for exactly this reason).
// internal/cmscore.SignerInfo has no such cache (see its doc comment) - every caller in this
// port that needs to check a signature against real content (CAdESSignature.
// signedContentForIntegrityCheck / ReferenceValidationsForSignerInformation) computes and
// compares the digest itself rather than reading it back off the SignerInformation, so a
// recomputed content digest has nothing to attach to. This function therefore reduces to the
// SignerId lookup alone; digestCalculatorProvider and resourcesHandlerBuilder are accepted for
// API parity (the same "accepted but not consulted" convention CMSUtilsWriteToDSSDocument's
// resourcesHandlerBuilder already uses, see this file's header) and never consulted.
//
// A signerId with no match returns an error - matching BC's Store#get, which upstream never
// observes returning null because getSignerId() is always derived from the very cms being
// searched.
func CMSUtilsRecomputeSignerInformation(cmsDoc *CMS, signerId *cmscore.SignerIdentifier,
	digestCalculatorProvider DigestCalculatorProvider, resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) (*cmscore.SignerInfo, error) {
	target := signerId.DER()
	for _, signerInformation := range cmsDoc.SignerInfos() {
		if bytes.Equal(signerInformation.SID.DER(), target) {
			return signerInformation, nil
		}
	}
	return nil, model.NewDSSError("Unable to find a SignerInformation matching the provided SignerId")
}

// digestAlgorithmsOfSigners returns the digest algorithm identifier of every SignerInfo, one
// per signer (duplicates removed by mergeAlgorithmIdentifiers at the call site).
func digestAlgorithmsOfSigners(signers []*cmscore.SignerInfo) []*asn1ber.AlgorithmIdentifier {
	identifiers := make([]*asn1ber.AlgorithmIdentifier, len(signers))
	for index, signer := range signers {
		identifiers[index] = signer.DigestAlgorithm
	}
	return identifiers
}

// CMSUtilsReplaceCertificatesAndCRLs replaces SignedData content within cms with the provided
// values. Port of #replaceCertificatesAndCRLs(CMS, Store<X509CertificateHolder>,
// Store<X509AttributeCertificateHolder>, Store<X509CRLHolder>, Store<?>, Store<?>); the last
// two Store<?> parameters are, per doc.go's convention, the id-ri-ocsp-response and
// id-pkix-ocsp-basic members of SignedData.crls respectively.
//
// DEVIATION: as CMSUtilsReplaceSigners, SignedData.version is kept from the original rather
// than recomputed, matching BouncyCastle's CMSSignedData#replaceCertificatesAndCRLs (a
// low-level field swap, not a full rebuild).
func CMSUtilsReplaceCertificatesAndCRLs(cms *CMS, certificates, attributeCertificates, crls,
	ocspResponses, ocspBasicResponses [][]byte) (*CMS, error) {
	original := cms.Core().SignedData()
	builder := &cmscore.SignedDataBuilder{
		DigestAlgorithms: original.DigestAlgorithms,
		EncapContentInfo: original.EncapContentInfo,
		Certificates:     certificateChoices(certificates, attributeCertificates),
		CRLs:             revocationInfoChoices(crls, ocspResponses, ocspBasicResponses),
		SignerInfos:      original.SignerInfos,
	}
	signedData, err := builder.Build()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to replace validation content of CMS SignedData. Corrupted content has been provided.", err)
	}
	signedData.Version = original.Version
	return newCMS(cmscore.NewCMS(signedData)), nil
}

// CMSUtilsPopulateDigestAlgorithmSet adds digest algorithms to cms. Port of
// #populateDigestAlgorithmSet(CMS, Collection<AlgorithmIdentifier>).
func CMSUtilsPopulateDigestAlgorithmSet(cms *CMS, digestAlgorithmsToAdd []*asn1ber.AlgorithmIdentifier) (*CMS, error) {
	original := cms.Core().SignedData()
	certificates, crls := cmsUtilsCertificatesAndCRLs(original)
	builder := &cmscore.SignedDataBuilder{
		DigestAlgorithms: mergeAlgorithmIdentifiers(original.DigestAlgorithms, digestAlgorithmsToAdd),
		EncapContentInfo: original.EncapContentInfo,
		Certificates:     certificates,
		CRLs:             crls,
		SignerInfos:      original.SignerInfos,
	}
	signedData, err := builder.Build()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Unable to populate digest algorithms within CMS SignedData. Corrupted content has been provided.", err)
	}
	signedData.Version = original.Version
	return newCMS(cmscore.NewCMS(signedData)), nil
}

// CMSUtilsToCMS converts a TimeStampToken to a CMS. Port of #toCMS(TimeStampToken); the
// argument is *cmscore.TimeStampToken (internal/cmscore's RFC 3161 replacement of
// BouncyCastle's org.bouncycastle.tsp.TimeStampToken - see PORTING.md), whose own #CMS() answers
// exactly the *cmscore.CMS BouncyCastle's timeStampToken.toCMSSignedData() would build.
func CMSUtilsToCMS(timeStampToken *cmscore.TimeStampToken) (*CMS, error) {
	return newCMS(timeStampToken.CMS()), nil
}

// CMSUtilsContentInfoEncoding gets encoding of the ContentInfo of cms: "DL" for
// definite-length, "BER" otherwise. Port of #getContentInfoEncoding.
func CMSUtilsContentInfoEncoding(cms *CMS) string {
	if cms.core.IsDefiniteLength() {
		return "DL"
	}
	return "BER"
}

// CMSUtilsWriteSignedDataDigestAlgorithmsEncoded writes the encoded binaries of the
// SignedData.digestAlgorithms field to w. Port of #writeSignedDataDigestAlgorithmsEncoded.
//
// NOTE: used for evidence record hash computation (a later phase); the DER SET OF encoding
// written here is internal/cmscore's, i.e. it normalises a BER length to DER but otherwise
// preserves the original ordering when cms was parsed.
func CMSUtilsWriteSignedDataDigestAlgorithmsEncoded(cms *CMS, w io.Writer) error {
	element := cms.core.SignedData().DigestAlgorithmsElement()
	if element != nil {
		_, err := w.Write(element.DEREncoded())
		return err
	}
	algorithms := cms.core.SignedData().DigestAlgorithms
	encodings := make([][]byte, len(algorithms))
	for index, algorithm := range algorithms {
		encodings[index] = algorithm.DER()
	}
	_, err := w.Write(cmscoreDERSetOf(encodings))
	return err
}

// CMSUtilsWriteContentInfoEncoded writes the encoded binaries of the ContentInfo element to w.
// Port of #writeContentInfoEncoded.
//
// NOTE: used for archive-time-stamp-v2 message-imprint computation.
//
// The "ContentInfo" of upstream's method name is SignedData.encapContentInfo - the eContentType
// and eContent of the SIGNED content (CMSObjectUtils#writeContentInfoEncoded:
// `signedData.getEncapContentInfo()`) - NOT the outer CMS ContentInfo that wraps the whole
// SignedData. Writing the outer one instead (this port's first reading of the name) fed the
// entire signature - certificates, SignerInfos and all - into the archive-timestamp-v2 message
// imprint, so every CAdES archive-timestamp-v2 verified as FAILED/HASH_FAILURE; found by the
// phase-8f document-level harness on Signature-C-B-LTA-10.p7m.
//
// Upstream picks the encoding from the eContent's own form - BER when the OCTET STRING is
// constructed (BouncyCastle hands such an eContent back as a BEROctetString), DER otherwise -
// which is what the two branches below reproduce.
func CMSUtilsWriteContentInfoEncoded(cms *CMS, w io.Writer) error {
	encapsulatedContentInfo := cms.core.SignedData().EncapContentInfo
	contentElement := encapsulatedContentInfo.ContentElement()
	var encoded []byte
	if contentElement != nil && contentElement.IsConstructed() {
		// DSSASN1Utils.getBEREncoded(ContentInfo): ContentInfo#toASN1Primitive builds a
		// BERSequence holding the eContentType and a BERTaggedObject(0, eContent), and a BER
		// ASN1OutputStream writes every one of those with the indefinite-length form.
		encoded = asn1ber.WriteIndefiniteTLV([]byte{asn1ber.ClassUniversal | asn1ber.Constructed | asn1ber.TagSequence},
			append(asn1ber.EncodeOID(encapsulatedContentInfo.EContentType),
				asn1ber.WriteIndefiniteTLV([]byte{asn1ber.ClassContextSpecific | asn1ber.Constructed}, contentElement.BEREncoded())...))
	} else {
		encoded = encapsulatedContentInfo.DER()
	}
	_, err := w.Write(encoded)
	return err
}

// CMSUtilsWriteSignedDataCertificatesEncoded writes the encoded binaries of the
// SignedData.certificates field to w. Port of #writeSignedDataCertificatesEncoded.
//
// NOTE: used for archive-time-stamp-v2 message-imprint computation. Absent when
// SignedData.certificates is absent, matching Java's "Certificates are not present" no-op.
func CMSUtilsWriteSignedDataCertificatesEncoded(cms *CMS, w io.Writer) error {
	encoded := cms.core.SignedData().Certificates.Encoded()
	if encoded == nil {
		return nil
	}
	_, err := w.Write(encoded)
	return err
}

// CMSUtilsWriteSignedDataCRLsEncoded writes the encoded binaries of the SignedData.crls field to
// w. Port of #writeSignedDataCRLsEncoded.
//
// NOTE: used for archive-time-stamp-v2 message-imprint computation. Absent when
// SignedData.crls is absent, matching Java's "CRLs are not present" no-op.
func CMSUtilsWriteSignedDataCRLsEncoded(cms *CMS, w io.Writer) error {
	encoded := cms.core.SignedData().CRLs.Encoded()
	if encoded == nil {
		return nil
	}
	_, err := w.Write(encoded)
	return err
}

// CMSUtilsWriteSignedDataSignerInfosEncoded writes the encoded binaries of the
// SignedData.signerInfos field to w. Port of #writeSignedDataSignerInfosEncoded.
//
// NOTE: used for evidence record hash computation.
func CMSUtilsWriteSignedDataSignerInfosEncoded(cms *CMS, w io.Writer) error {
	element := cms.core.SignedData().SignerInfosElement()
	if element != nil {
		_, err := w.Write(element.DEREncoded())
		return err
	}
	signers := cms.core.SignerInfos()
	encodings := make([][]byte, len(signers))
	for index, signer := range signers {
		encodings[index] = signer.DER()
	}
	_, err := w.Write(cmscoreDERSetOf(encodings))
	return err
}

// CMSUtilsToCMSEncapsulatedContent converts a DSSDocument to the octets its SignedData.encapContentInfo.eContent
// would carry, nil for a *model.DigestDocument (whose content is not available - BouncyCastle's
// CMSAbsentContent). Port of #toCMSEncapsulatedContent(DSSDocument), minus the CMSTypedData
// wrapping: that BC type has no Go replacement, callers here only ever need the octets
// themselves (see AbstractCMSGenerator#encapsulatedContentBytes, which additionally applies the
// "only if encapsulating" choice this function - like Java's - does not make).
//
// Panics with the Java message when document is nil (Objects.requireNonNull).
func CMSUtilsToCMSEncapsulatedContent(document model.DSSDocument) ([]byte, error) {
	if document == nil {
		panic("Document to be signed is missing")
	}
	if _, isDigestDocument := document.(*model.DigestDocument); isDigestDocument {
		return nil, nil
	}
	return spi.DSSUtilsToByteArrayOfDocument(document)
}

// CMSUtilsResourcesHandlerBuilder verifies whether dssResourcesHandlerBuilder is supported by
// the current implementation, returning it unchanged on success. Port of
// #getDSSResourcesHandlerBuilder; see the file header for why this never rejects one.
func CMSUtilsResourcesHandlerBuilder(dssResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) resources.DSSResourcesHandlerBuilder {
	return dssResourcesHandlerBuilder
}

// CMSUtilsReplaceUnsignedAttributes replaces unsignedAttributes within the given signerInformation.
// Port of #replaceUnsignedAttributes(SignerInformation, AttributeTable); BouncyCastle's
// SignerInformation#replaceUnsignedAttributes keeps every other field of the SignerInfo (its
// version, sid, digestAlgorithm, signedAttrs and signature) untouched.
func CMSUtilsReplaceUnsignedAttributes(signerInformation *cmscore.SignerInfo, unsignedAttributes cmscore.Attributes) (*cmscore.SignerInfo, error) {
	builder := &cmscore.SignerInfoBuilder{
		SID:                signerInformation.SID,
		DigestAlgorithm:    signerInformation.DigestAlgorithm,
		SignedAttributes:   signerInformation.SignedAttributes,
		SignatureAlgorithm: signerInformation.SignatureAlgorithm,
		Signature:          signerInformation.Signature,
		UnsignedAttributes: unsignedAttributes,
	}
	if !signerInformation.HasSignedAttributes() {
		builder.SignedAttributes = nil
	}
	updated, err := builder.Build()
	if err != nil {
		return nil, err
	}
	updated.Version = signerInformation.Version
	return updated, nil
}

// CMSUtilsAssertATSv2AugmentationSupported reports whether the augmentation of signatures with
// an archive-time-stamp-v2 is supported by the current implementation. Port of
// #assertATSv2AugmentationSupported; this native implementation always supports it (Java's
// dss-cms-object answers "supported, do nothing" too - only a future streaming implementation
// might not).
func CMSUtilsAssertATSv2AugmentationSupported() error { return nil }

// CMSUtilsAssertEvidenceRecordEmbeddingSupported reports whether the embedding of existing
// Evidence Records within CMS is supported by the current implementation. Port of
// #assertEvidenceRecordEmbeddingSupported; always supported, as CMSUtilsAssertATSv2AugmentationSupported.
func CMSUtilsAssertEvidenceRecordEmbeddingSupported() error { return nil }

// cmscoreDERSetOf writes a universal SET OF from its already DER-encoded members, X.690
// clause 11.6 order (members compared as unsigned octet strings). A local copy of
// internal/cmscore's unexported derSetOf(setIdentifier, members), needed here because
// CMSUtilsWriteSignedDataDigestAlgorithmsEncoded/-SignerInfosEncoded are the only two write-*
// functions whose field cmscore does not already keep pre-encoded as a SET (see
// SignedData.DigestAlgorithmsElement/SignerInfosElement, populated only when cms was parsed).
func cmscoreDERSetOf(members [][]byte) []byte {
	ordered := make([][]byte, len(members))
	copy(ordered, members)
	sort.SliceStable(ordered, func(a, b int) bool {
		return bytes.Compare(ordered[a], ordered[b]) < 0
	})
	var body []byte
	for _, member := range ordered {
		body = append(body, member...)
	}
	return asn1ber.WriteTLV(asn1ber.TagSet|asn1ber.Constructed, body)
}

// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/timestamp/CAdESTimestampMessageDigestBuilder.java (DSS 6.5.RC1).
//
// This file owns the byte-exact message-imprint computation for every CAdES timestamp type
// (content, signature, X1/X2, archive-v2/v3): the data a TSA digests and a verifier re-derives
// to check a timestamp's message-imprint. See archive-timestamp-v3's counterpart in
// cades_level_baseline_lta_timestamp_extractor.go, which this file drives for both the
// ats-hash-index-v3 re-verification and the V3 message-imprint itself; archive-timestamp-v2 is
// entirely owned here (there is no upstream V2-building counterpart to lean on: TS 101 733's
// Annex K legacy hash covers the CMS SignedData shape directly, not an ats-hash-index).
//
// FORWARD DEPENDENCY: *CAdESSignature (Java eu.europa.esig.dss.cades.validation.CAdESSignature)
// is assigned to a sibling chunk of this manifest and is not (re)declared here; every other
// landed file in this package already forward-references it the same way (see e.g.
// cades_level_baseline_lta_timestamp_extractor.go's NewCadesLevelBaselineLTATimestampExtractor).
// The subset of its API this file calls, inferred from the Java constructor bodies:
//
//	func (s *CAdESSignature) CertificateSource() *spi.SignatureCertificateSource // getCertificateSource()
//	func (s *CAdESSignature) CMS() *cms.CMS                                     // getCMS()
//	func (s *CAdESSignature) SignerInformation() *cmscore.SignerInfo           // getSignerInformation()
//	func (s *CAdESSignature) DetachedContents() []model.DSSDocument            // getDetachedContents()
//
// all four already assumed by cadesLTASignature/NewCadesLevelBaselineLTATimestampExtractor in
// this same package, so *CAdESSignature satisfying them is not a new requirement this file adds.
//
// FORWARD DEPENDENCY: the dss-cms package (Java eu.europa.esig.dss.cms), assigned to the
// dependency-closed CMSAPI chunk, is imported here exactly as cades_utils.go and
// cades_level_baseline_lta_timestamp_extractor.go already do (cms.CMS, cms.CMS#IsDetachedSignature).
// This file additionally needs three CMSUtils statics flattened per PORTING.md's
// "<JavaClass><MethodName>" convention for static utility classes (the same convention
// CAdESUtils already follows in this very package):
//
//	func CMSUtilsWriteContentInfoEncoded(cmsDocument *cms.CMS, w io.Writer) error
//	func CMSUtilsWriteSignedDataCertificatesEncoded(cmsDocument *cms.CMS, w io.Writer) error
//	func CMSUtilsWriteSignedDataCRLsEncoded(cmsDocument *cms.CMS, w io.Writer) error
//
// ports of ICMSUtils#writeContentInfoEncoded/writeSignedDataCertificatesEncoded/
// writeSignedDataCRLsEncoded(CMS, OutputStream), each of which "Writes the encoded binaries of
// the [...] field to the given OutputStream" (their Javadoc) - i.e. an io.Writer per this port's
// stream convention (PORTING.md "Streams & documents"), with the declared IOException becoming
// an error return per PORTING.md's "Errors, exceptions, alerts".
//
// DEVIATIONS:
//
//   - getSignedDataCertificateReferences/getSignatureSignedDataReferences in
//     cades_timestamp_source.go and, transitively, nothing in *this* file itself needs the
//     `instanceof CMSCertificateSource`/`instanceof CMSCRLSource`/`instanceof CMSOCSPSource`
//     guards Java's CAdESTimestampSource uses: Go's spi.SignatureCertificateSource already
//     exposes SignedDataCertificates() unconditionally on the base type (mirroring Java's own
//     SignatureCertificateSource#getSignedDataCertificates(), which - contrary to what the
//     instanceof guard suggests - is *also* already declared on the abstract base class, not
//     added by CMSCertificateSource; the guard is defensive/redundant in Java too). This file
//     does not read certificate/revocation sources at all, so the deviation is noted here for
//     the sibling file that does.
//   - writeSignerInfoBytes reads signerInformation.UnsignedAttributes directly rather than
//     rejecting a SignerInfo with no unsignedAttrs field at all (Java: NullPointerException from
//     unauthenticatedAttributes.size()); an absent field yields an empty attribute list instead.
//     This mirrors the identical, already-documented judgment call in
//     cades_level_baseline_lta_timestamp_extractor.go's unsignedAttributesHashIndex, for the same
//     reason: unreachable in the augmentation flow that ever calls this (an archive-timestamp is
//     only ever computed for a signature that already carries at least a signature-time-stamp
//     unsigned attribute), so no byte of a reachable output changes.
//   - the outer try/catch of getTimestampX1MessageDigest/getTimestampX2MessageDigest/
//     getArchiveTimestampDataV2, which Java uses to turn any failure into a
//     log-and-return-null/empty-digest, is reproduced with ordinary Go error returns from
//     unexported helpers, converted to the null/empty digest at the two exported boundary
//     methods and at getArchiveTimestampDataV2 - see the MESSAGE_IMPRINT_ERROR/errorMessage
//     upstream log format strings, carried here only as comments (slf4j dropped per PORTING.md).
//     A failure that Java lets escape uncaught (getContentTimestampMessageDigest's
//     originalDocument.getDigestValue, getSignatureTimestampMessageDigest's DSSUtils.digest,
//     getArchiveTimestampDataV3's timestampExtractor calls) panics instead, matching the
//     unchecked-DSSException-propagates convention this port uses throughout (see e.g.
//     spi/validation/timestamp/abstract_timestamp_source.go's `must` helper).
package cades

import (
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/timestamp"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESTimestampMessageDigestBuilder builds timestamped data binaries for a CAdES signature.
type CAdESTimestampMessageDigestBuilder struct {
	// cms is the CMS.
	cms *cms.CMS

	// signerInformation is the SignerInformation of the related signature.
	signerInformation *cmscore.SignerInfo

	// detachedDocuments is the list of detached documents.
	detachedDocuments []model.DSSDocument

	// timestampExtractor is the instance of CadesLevelBaselineLTATimestampExtractor.
	timestampExtractor *CadesLevelBaselineLTATimestampExtractor

	// digestAlgorithm is the digest algorithm to be used for message-imprint digest computation.
	digestAlgorithm enumerations.DigestAlgorithm

	// timestampToken is the timestamp token to compute message-digest for.
	timestampToken *validation.TimestampToken
}

// NewCAdESTimestampMessageDigestBuilder is the constructor to compute message-imprint for
// timestamps related to the signature, to be used on timestamp creation.
// Port of the (CAdESSignature, DigestAlgorithm) constructor.
//
// Panics with the Java message when digestAlgorithm is empty (Objects.requireNonNull).
func NewCAdESTimestampMessageDigestBuilder(signature *CAdESSignature,
	digestAlgorithm enumerations.DigestAlgorithm) *CAdESTimestampMessageDigestBuilder {
	builder := newCAdESTimestampMessageDigestBuilder(signature, signature.CertificateSource().SignedDataCertificates())
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	builder.digestAlgorithm = digestAlgorithm
	return builder
}

// NewCAdESTimestampMessageDigestBuilderForToken is the constructor to compute message-imprint
// for timestamps related to the signature. This constructor uses the provided certificateSource
// to validate the ats-v3-hash-table. Port of the (CAdESSignature, ListCertificateSource,
// TimestampToken) constructor.
//
// Panics with the Java message when timestampToken is nil (Objects.requireNonNull).
func NewCAdESTimestampMessageDigestBuilderForToken(signature *CAdESSignature,
	certificateSource *spi.ListCertificateSource, timestampToken *validation.TimestampToken) *CAdESTimestampMessageDigestBuilder {
	builder := newCAdESTimestampMessageDigestBuilder(signature, certificateSource.Certificates())
	if timestampToken == nil {
		panic("TimestampToken cannot be null!")
	}
	builder.timestampToken = timestampToken
	builder.digestAlgorithm = timestampToken.DigestAlgorithm()
	return builder
}

// newCAdESTimestampMessageDigestBuilder is the default (private) constructor.
// Port of the (CAdESSignature, List<CertificateToken>) constructor.
//
// certificates is required non-nil (Objects.requireNonNull, like signature), matching Java,
// even though - as upstream - it is otherwise unused by this constructor: the two public
// constructors above only ever pass it for the requireNonNull side effect.
//
// Panics with the Java messages when signature or certificates is nil.
func newCAdESTimestampMessageDigestBuilder(signature *CAdESSignature,
	certificates []*model.CertificateToken) *CAdESTimestampMessageDigestBuilder {
	if signature == nil {
		panic("Signature cannot be null!")
	}
	if certificates == nil {
		panic("List of CertificateToken's cannot be null!")
	}
	return &CAdESTimestampMessageDigestBuilder{
		cms:                signature.CMS(),
		signerInformation:  signature.SignerInformation(),
		detachedDocuments:  signature.DetachedContents(),
		timestampExtractor: NewCadesLevelBaselineLTATimestampExtractor(signature),
	}
}

// ContentTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getContentTimestampMessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) ContentTimestampMessageDigest() model.DSSMessageDigest {
	return b.originalDocumentDigest()
}

// SignatureTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getSignatureTimestampMessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) SignatureTimestampMessageDigest() model.DSSMessageDigest {
	signature := b.signerInformation.Signature
	digest, err := spi.DSSUtilsDigest(b.digestAlgorithm, signature)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return model.NewDSSMessageDigestWithValue(b.digestAlgorithm, digest)
}

// TimestampX1MessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getTimestampX1MessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) TimestampX1MessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.timestampX1MessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR ("Unable to compute message-imprint for
		// TimestampToken with Id '{}'. Reason : {}") and returns null.
		return model.NewDSSMessageDigest()
	}
	return messageDigest
}

// timestampX1MessageDigest is the fallible core of TimestampX1MessageDigest.
// Port of the try block of getTimestampX1MessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) timestampX1MessageDigest() (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	digestCalculator.Update(b.signerInformation.Signature)
	// We don't include the outer SEQUENCE, only the attrType and attrValues as stated by the TS
	// §6.3.5, NOTE 2.

	attributes := CAdESUtilsUnsignedAttributesOfType(b.signerInformation, OID_id_aa_signatureTimeStampToken)
	if utils.IsArrayNotEmpty(attributes) {
		for _, attribute := range attributes {
			typeDER, valuesDER, err := cadesTMDBAttrTypeAndValuesDER(attribute)
			if err != nil {
				return model.DSSMessageDigest{}, err
			}
			digestCalculator.Update(typeDER)
			digestCalculator.Update(valuesDER)
		}
	}
	// Method is common to Type 1 and Type 2.
	if err := b.writeTimestampX2MessageDigest(digestCalculator); err != nil {
		return model.DSSMessageDigest{}, err
	}
	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// TimestampX2MessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getTimestampX2MessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) TimestampX2MessageDigest() model.DSSMessageDigest {
	messageDigest, err := b.timestampX2MessageDigest()
	if err != nil {
		// Upstream logs MESSAGE_IMPRINT_ERROR and returns null.
		return model.NewDSSMessageDigest()
	}
	return messageDigest
}

// timestampX2MessageDigest is the fallible core of TimestampX2MessageDigest.
// Port of the try block of getTimestampX2MessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) timestampX2MessageDigest() (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	if err := b.writeTimestampX2MessageDigest(digestCalculator); err != nil {
		return model.DSSMessageDigest{}, err
	}
	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// writeTimestampX2MessageDigest feeds digestCalculator with the certificate-refs and
// revocation-refs unsigned attributes, common to Type 1 (X1) and Type 2 (X2) message-imprints.
// Port of the private writeTimestampX2MessageDigest(DSSMessageDigestCalculator).
func (b *CAdESTimestampMessageDigestBuilder) writeTimestampX2MessageDigest(digestCalculator *spi.DSSMessageDigestCalculator) error {
	certAttributes := CAdESUtilsUnsignedAttributesOfType(b.signerInformation, spi.OID_id_aa_ets_certificateRefs)
	if utils.IsArrayNotEmpty(certAttributes) {
		for _, attribute := range certAttributes {
			typeDER, valuesDER, err := cadesTMDBAttrTypeAndValuesDER(attribute)
			if err != nil {
				return err
			}
			digestCalculator.Update(typeDER)
			digestCalculator.Update(valuesDER)
		}
	}
	revAttributes := CAdESUtilsUnsignedAttributesOfType(b.signerInformation, spi.OID_id_aa_ets_revocationRefs)
	if utils.IsArrayNotEmpty(revAttributes) {
		for _, attribute := range revAttributes {
			typeDER, valuesDER, err := cadesTMDBAttrTypeAndValuesDER(attribute)
			if err != nil {
				return err
			}
			digestCalculator.Update(typeDER)
			digestCalculator.Update(valuesDER)
		}
	}
	return nil
}

// cadesTMDBAttrTypeAndValuesDER returns the DER encoding of an attribute's attrType, and,
// separately, of its attrValues (the SET OF AttributeValue) - the two pieces
// writeTimestampX2MessageDigest (and, through it, TimestampX1MessageDigest) each feed to a
// digest calculator via two consecutive update() calls, deliberately WITHOUT the outer Attribute
// SEQUENCE wrapper: see the "We don't include the outer SEQUENCE" comment this ports.
//
// attribute is assumed parsed (attribute.Element() non-nil): every attribute this file reaches
// through CAdESUtilsUnsignedAttributesOfType comes from a *cmscore.SignerInfo obtained by
// parsing an existing CMS signature (this builder computes message-imprints for validation, not
// construction), which cmscore always parses from bytes. asn1ber.Element#DEREncoded() sorts a
// SET's members by their encoding on the way out, reproducing BouncyCastle's DERSet ordering for
// the (rare, but possible for e.g. multi-valued attrValues) case of more than one member.
func cadesTMDBAttrTypeAndValuesDER(attribute *cmscore.Attribute) (typeDER, valuesDER []byte, err error) {
	element := attribute.Element()
	if element == nil || len(element.Children()) < 2 {
		return nil, nil, model.NewDSSError(
			fmt.Sprintf("Unable to extract attrValues from attribute with OID '%s'", attribute.Type))
	}
	return asn1ber.EncodeOID(attribute.Type), element.Children()[1].DEREncoded(), nil
}

// ArchiveTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getArchiveTimestampMessageDigest().
func (b *CAdESTimestampMessageDigestBuilder) ArchiveTimestampMessageDigest() model.DSSMessageDigest {
	// V3 is used by default.
	archiveTimestampType := enumerations.ArchiveTimestampType_CAdES_V3
	if b.timestampToken != nil {
		archiveTimestampType = b.timestampToken.ArchiveTimestampType()
	}

	switch archiveTimestampType {
	case enumerations.ArchiveTimestampType_CAdES_V2:
		/*
		 * There is a difference between message imprint calculation in ETSI TS 101 733 version
		 * 1.8.3 and version 2.2.1. So we first check the message imprint according to 2.2.1
		 * version and then if it fails get the message imprint data for the 1.8.3 version
		 * message imprint calculation.
		 */
		messageDigest := b.archiveTimestampDataV2(true)
		if !b.timestampToken.MatchDataMessageDigestSuppressingWarnings(messageDigest, true) {
			// Upstream logs "Unable to match message imprint for an Archive TimestampToken V2
			// with Id '{}' by including unsigned attribute tags and length, try to compute the
			// data without...".
			messageDigest = b.archiveTimestampDataV2(false)
		}
		return messageDigest
	case enumerations.ArchiveTimestampType_CAdES_V3:
		return b.archiveTimestampDataV3()
	default:
		panic(model.NewDSSError(fmt.Sprintf("Unsupported ArchiveTimestampType %s", archiveTimestampType)))
	}
}

// archiveTimestampDataV3 ports the private getArchiveTimestampDataV3(), which declares `throws
// DSSException`: nothing catches it at either of its two call sites upstream (this method,
// itself uncaught in ArchiveTimestampMessageDigest), so a failure here panics, matching Java's
// unchecked propagation.
func (b *CAdESTimestampMessageDigestBuilder) archiveTimestampDataV3() model.DSSMessageDigest {
	atsHashIndexAttribute, err := b.timestampExtractor.VerifiedAtsHashIndex(b.signerInformation, b.timestampToken)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	originalDocument := b.originalDocument()
	if originalDocument == nil {
		// Upstream logs "The original document is not found for TimestampToken with Id '{}'!
		// Unable to compute message imprint.".
		return model.CreateEmptyDSSMessageDigest()
	}
	messageDigest, err := b.timestampExtractor.ArchiveTimestampV3MessageImprint(
		b.signerInformation, atsHashIndexAttribute, originalDocument, b.digestAlgorithm)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return messageDigest
}

// originalDocumentDigest ports the private getOriginalDocumentDigest().
func (b *CAdESTimestampMessageDigestBuilder) originalDocumentDigest() model.DSSMessageDigest {
	originalDocument := b.originalDocument()
	if originalDocument == nil {
		// Upstream logs "The original document is not found for TimestampToken with Id '{}'!
		// Unable to compute message imprint.".
		return model.CreateEmptyDSSMessageDigest()
	}
	digest, err := originalDocument.DigestValue(b.digestAlgorithm)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return model.NewDSSMessageDigestWithValue(b.digestAlgorithm, digest)
}

// archiveTimestampDataV2 ports the private getArchiveTimestampDataV2(boolean), whose declared
// `throws DSSException` never happens: every failure inside the try block is caught by its own
// catch-all and converted to model.CreateEmptyDSSMessageDigest(), reproduced here by
// archiveTimestampDataV2Try's error return.
//
// There is a difference in ETSI TS 101 733 version 1.8.3 and version 2.2.1 in archive-timestamp-v2
// hash calculation. In the 1.8.3 version the calculation did not include the tag and the length
// octets of the unsigned attributes set. The hash calculation is described in Annex K in both
// versions of ETSI TS 101 733. The differences are in TableK.3: Signed Data in rows 22 and 23.
// However, there is a note in 2.2.1 version (Annex K, Table K.3: SignedData, Note 3) that says:
// "A previous version of CAdES did not include the tag and length octets of this SET OF type of
// unsignedAttrs element in this annex, which contradicted the normative section. To maximize
// interoperability, it is recommended to simultaneously compute the two hash values (including
// and not including the tag and length octets of SET OF type) and to test the value of the
// timestamp against both." The includeUnsignedAttrsTagAndLength parameter decides whether the
// tag and length octets are included.
//
// According to RFC 5652 it is possible to use DER or BER encoding for SignedData structure. The
// exception is the signed attributes attribute and authenticated attributes which have to be DER
// encoded.
func (b *CAdESTimestampMessageDigestBuilder) archiveTimestampDataV2(includeUnsignedAttrsTagAndLength bool) model.DSSMessageDigest {
	messageDigest, err := b.archiveTimestampDataV2Try(includeUnsignedAttrsTagAndLength)
	if err != nil {
		// When error in computing or in format the algorithm just continues. Upstream logs
		// "An error in computing of message-imprint for a TimestampToken with Id : %s. Reason :
		// %s".
		return model.CreateEmptyDSSMessageDigest()
	}
	return messageDigest
}

// archiveTimestampDataV2Try is the fallible core of archiveTimestampDataV2.
// Port of the try block of getArchiveTimestampDataV2(boolean).
func (b *CAdESTimestampMessageDigestBuilder) archiveTimestampDataV2Try(includeUnsignedAttrsTagAndLength bool) (model.DSSMessageDigest, error) {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, err
	}

	dos := digestCalculator.Writer()
	defer func() { _ = dos.Close() }()

	if err := b.writeContentInfoBytes(dos); err != nil {
		return model.DSSMessageDigest{}, err
	}

	if b.cms.IsDetachedSignature() {
		if err := b.writeOriginalDocumentBinaries(dos); err != nil {
			return model.DSSMessageDigest{}, err
		}
	}

	if err := b.writeCertificateDataBytes(dos); err != nil {
		return model.DSSMessageDigest{}, err
	}

	if err := b.writeCRLDataBytes(dos); err != nil {
		return model.DSSMessageDigest{}, err
	}

	if err := b.writeSignerInfoBytes(dos, includeUnsignedAttrsTagAndLength); err != nil {
		return model.DSSMessageDigest{}, err
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm), nil
}

// writeContentInfoBytes ports the private writeContentInfoBytes(OutputStream).
func (b *CAdESTimestampMessageDigestBuilder) writeContentInfoBytes(w io.Writer) error {
	return cms.CMSUtilsWriteContentInfoEncoded(b.cms, w)
}

// writeOriginalDocumentBinaries ports the private writeOriginalDocumentBinaries(OutputStream).
//
// Detached signatures have either no encapContentInfo in signedData, or it exists but has no
// eContent.
func (b *CAdESTimestampMessageDigestBuilder) writeOriginalDocumentBinaries(w io.Writer) error {
	originalDocument := b.originalDocument()
	if originalDocument == nil {
		return model.NewDSSError(fmt.Sprintf(
			"The detached content is not provided for a TimestampToken with Id '%s'. Not possible to compute message imprint!",
			b.timestampToken.DSSIDAsString()))
	}
	bytes, err := spi.DSSUtilsToByteArrayOfDocument(originalDocument)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

// writeCertificateDataBytes ports the private writeCertificateDataBytes(OutputStream).
func (b *CAdESTimestampMessageDigestBuilder) writeCertificateDataBytes(w io.Writer) error {
	return cms.CMSUtilsWriteSignedDataCertificatesEncoded(b.cms, w)
}

// writeCRLDataBytes ports the private writeCRLDataBytes(OutputStream).
func (b *CAdESTimestampMessageDigestBuilder) writeCRLDataBytes(w io.Writer) error {
	return cms.CMSUtilsWriteSignedDataCRLsEncoded(b.cms, w)
}

// writeSignerInfoBytes ports the private writeSignerInfoBytes(OutputStream, boolean), together
// with getSignerInfoEncoded(SignerInfo, ASN1Sequence, boolean) - "Copied from
// org.bouncycastle.asn1.cms.SignerInfo#toASN1Object() and adapted to be able to use the custom
// unauthenticatedAttributes" - inlined here since Go has no BouncyCastle SignerInfo to copy from
// and every field is written directly to w as soon as it is DER-encoded, rather than assembled
// into an intermediate ASN1Sequence first: writing DER(field) for field in [version, sid,
// digestAlgorithm, signedAttrs?, signatureAlgorithm, signature, unsignedAttrs-tagged-or-flattened]
// one at a time produces the identical concatenation of octets that iterating
// getSignerInfoEncoded(...)'s resulting ASN1Sequence and DER-encoding each member individually
// does upstream.
func (b *CAdESTimestampMessageDigestBuilder) writeSignerInfoBytes(w io.Writer, includeUnsignedAttrsTagAndLength bool) error {
	signerInfo := b.signerInformation

	filteredUnauthenticatedAttributes, err := b.filterUnauthenticatedAttributes()
	if err != nil {
		return err
	}

	write := func(member []byte) error {
		_, writeErr := w.Write(member)
		return writeErr
	}

	if err := write(cadesLTAEncodedVersion(signerInfo.Version)); err != nil {
		return err
	}
	if err := write(signerInfo.SID.DER()); err != nil {
		return err
	}
	if err := write(signerInfo.DigestAlgorithm.DER()); err != nil {
		return err
	}

	signedAttributes, err := CAdESUtilsDERSignedAttributes(signerInfo)
	if err != nil {
		return err
	}
	if signedAttributes != nil {
		if err := write(signedAttributes); err != nil {
			return err
		}
	}

	if err := write(signerInfo.SignatureAlgorithm.DER()); err != nil {
		return err
	}
	if err := write(cadesLTADEROctetString(signerInfo.Signature)); err != nil {
		return err
	}

	// unauthenticatedAttributes (the filtered result) is never nil upstream:
	// filterUnauthenticatedAttributes always returns a (possibly empty) DERSequence, so
	// `if (unauthenticatedAttributes != null)` is unconditionally true and this branch always
	// runs - including when filteredUnauthenticatedAttributes is empty, e.g. when re-verifying
	// the message-imprint of the sole archive-timestamp of a signature (which the filter always
	// excludes from its own covered unsigned attributes).
	if includeUnsignedAttrsTagAndLength {
		var content []byte
		for _, member := range filteredUnauthenticatedAttributes {
			content = append(content, member...)
		}
		if err := write(asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, content)); err != nil {
			return err
		}
	} else {
		for _, member := range filteredUnauthenticatedAttributes {
			if err := write(member); err != nil {
				return err
			}
		}
	}
	return nil
}

// filterUnauthenticatedAttributes returns the DER encoding of every unsigned attribute of
// signerInformation that survives the archive-timestamp-v2/v3 filter: an archive-timestamp
// attribute added at or after b.timestampToken's own generation time is excluded, since an
// archive-timestamp only ever seals the state of the signature as it stood before it was itself
// added. Port of the private filterUnauthenticatedAttributes(ASN1Set, TimestampToken), returning
// the DER-encoded members directly (see writeSignerInfoBytes) rather than an ASN1Sequence, and
// dropping its `new DERSequence(result)` wrapper accordingly.
//
// DEVIATION: reads signerInformation.UnsignedAttributes directly instead of raising on a
// SignerInfo carrying no unsignedAttrs field at all (Java: NullPointerException from
// unauthenticatedAttributes.size()); see the file header.
func (b *CAdESTimestampMessageDigestBuilder) filterUnauthenticatedAttributes() ([][]byte, error) {
	var kept [][]byte
	for _, attribute := range b.signerInformation.UnsignedAttributes {
		if spi.OID_id_aa_ets_archiveTimestampV2.Equal(attribute.Type) || spi.OID_id_aa_ets_archiveTimestampV3.Equal(attribute.Type) {
			token := CAdESUtilsTimeStampToken(attribute)
			if token == nil || !token.TSTInfo().GenTime.Before(b.timestampToken.GenerationTime()) {
				continue
			}
		}
		element := attribute.Element()
		if element == nil {
			return nil, model.NewDSSError("Unexpected error occurred on reading unsigned properties")
		}
		kept = append(kept, element.DEREncoded())
	}
	return kept, nil
}

// originalDocument ports the private getOriginalDocument(), which never lets
// CAdESUtils.getOriginalDocument's DSSException escape: it is caught, logged, and answered with
// null.
func (b *CAdESTimestampMessageDigestBuilder) originalDocument() model.DSSDocument {
	document, err := CAdESUtilsOriginalDocument(b.cms, b.detachedDocuments)
	if err != nil {
		// Upstream logs "Cannot extract original document! Reason : {}".
		return nil
	}
	return document
}

// compile-time assertion: *CAdESTimestampMessageDigestBuilder satisfies
// timestamp.TimestampMessageDigestBuilder, matching Java's "implements TimestampMessageDigestBuilder".
var _ timestamp.TimestampMessageDigestBuilder = (*CAdESTimestampMessageDigestBuilder)(nil)

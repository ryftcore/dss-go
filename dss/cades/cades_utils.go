// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/CAdESUtils.java (DSS 6.5.RC1).
//
// # BouncyCastle replacements
//
//   - org.bouncycastle.cms.SignerInformation      -> *cmscore.SignerInfo
//   - org.bouncycastle.asn1.cms.AttributeTable    -> cmscore.Attributes
//   - org.bouncycastle.asn1.cms.Attribute         -> *cmscore.Attribute
//   - org.bouncycastle.cms.CMSSignedData          -> *cmscore.CMS
//   - org.bouncycastle.tsp.TimeStampToken         -> *cmscore.TimeStampToken
//   - org.bouncycastle.asn1.ASN1ObjectIdentifier  -> encoding/asn1.ObjectIdentifier
//   - every other ASN.1 object (ASN1Sequence, ASN1Set, ASN1Encodable, DERTaggedObject, the ESS
//     SigningCertificate family) -> its encoding, []byte, the representation dss-spi already
//     chose for an ASN.1 value; the BER/DER engine that produces and re-encodes it is
//     internal/asn1ber. Bytes that a digest or an identity depends on are carried through
//     untouched and never rebuilt.
//
// # The ESS structures built here
//
// AddSigningCertificateAttribute is the only place in this port that *writes* an ESS
// SigningCertificate / SigningCertificateV2; spi/cms_certificate_source.go models the read
// side. Rather than grow the frozen dss-spi types a build side, the three SEQUENCEs are
// assembled here from asn1ber primitives, next to their single consumer - the same choice
// dss-spi made when it put the read side next to *its* consumer. The encodings reproduce
// BouncyCastle's, ESSCertIDv2's DEFAULT hashAlgorithm included (id-sha256 is omitted).
//
// # Deviations
//
//   - getEvidenceRecordGenerationTime takes the DER encoding of an RFC 4998 EvidenceRecord
//     rather than a parsed object: there is no Go counterpart of
//     org.bouncycastle.asn1.tsp.EvidenceRecord yet (dss-evidence-record-asn1 is a later phase),
//     and the method only needs to reach the first ArchiveTimeStamp's ContentInfo.
//   - Methods that Java lets fail with an unchecked exception on malformed input return an
//     error instead, per PORTING.md; slf4j warnings are dropped and noted in comments.
package cades

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESUtilsDefaultArchiveTimestampHashAlgo is the default DigestAlgorithm for an
// ArchiveTimestamp. Port of DEFAULT_ARCHIVE_TIMESTAMP_HASH_ALGO.
const CAdESUtilsDefaultArchiveTimestampHashAlgo = enumerations.DigestAlgorithmSHA256

// CAdESUtilsDefaultResourcesHandlerBuilder is the default resources handler builder to be used
// across the code. Port of DEFAULT_RESOURCES_HANDLER_BUILDER.
var CAdESUtilsDefaultResourcesHandlerBuilder = document.NewInMemoryResourcesHandlerBuilder()

// 01-01-1950 and 01-01-2050, the bounds RFC 3852 sets on UTCTime; port of the private
// JANUARY_1950 / JANUARY_2050 constants (DSSUtils#getUtcDate takes a zero-based month).
var (
	cadesUtilsJanuary1950 = spi.DSSUtilsUTCDate(1950, 0, 1)
	cadesUtilsJanuary2050 = spi.DSSUtilsUTCDate(2050, 0, 1)
)

// The PKCS#9 attribute types this file recognises. They are
// org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers constants upstream and keep their exact Java
// field name behind the "OID_" prefix, the convention spi/oid.go established.
var (
	// OIDIdAaEtsContentTimestamp is 1.2.840.113549.1.9.16.2.20.
	OIDIdAaEtsContentTimestamp = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 20}

	// OIDIdAaSignatureTimeStampToken is 1.2.840.113549.1.9.16.2.14.
	OIDIdAaSignatureTimeStampToken = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 14}

	// OIDIdAaEtsCertCRLTimestamp is 1.2.840.113549.1.9.16.2.26.
	OIDIdAaEtsCertCRLTimestamp = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 26}

	// OIDIdAaEtsEscTimeStamp is 1.2.840.113549.1.9.16.2.25.
	OIDIdAaEtsEscTimeStamp = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 25}
)

// cadesUtilsTimestampOids contains a list of all CAdES timestamp OIDs, in the order the static
// initialiser fills it. Port of the private timestampOids field.
var cadesUtilsTimestampOids = []asn1.ObjectIdentifier{
	OIDIdAaEtsContentTimestamp,
	spi.OIDIdAaEtsArchiveTimestampV2,
	spi.OIDIdAaEtsArchiveTimestampV3,
	OIDIdAaEtsCertCRLTimestamp,
	OIDIdAaEtsEscTimeStamp,
	OIDIdAaSignatureTimeStampToken,
}

// cadesUtilsEvidenceRecordOids contains a list of all CAdES evidence record OIDs. Port of the
// private evidenceRecordOids field.
var cadesUtilsEvidenceRecordOids = []asn1.ObjectIdentifier{
	spi.OIDIdAaErInternal,
	spi.OIDIdAaErExternal,
}

// UtilsDERSignedAttributes gets the DER SignedAttributes table from the given
// SignerInformation, as the [0] IMPLICIT field of a SignerInfo.
// Port of getDERSignedAttributes(SignerInformation).
//
// Java re-parses SignerInformation#getEncodedSignedAttributes - the DER SET OF that RFC 5652
// clause 5.4 makes the signature input - and wraps it in an implicitly tagged [0]; cmscore
// produces exactly that encoding from the parsed attributes. Returns nil when the SignerInfo
// carries no signed attributes, as Java returns null. The error return is Java's
// DSSException("Unable to extract SignedAttributes. Reason : ..."), which wraps the IOException
// getEncodedSignedAttributes may raise; re-encoding parsed attributes cannot fail here, so it is
// nil today.
func UtilsDERSignedAttributes(signerInformation *cmscore.SignerInfo) ([]byte, error) {
	if signerInformation == nil || !signerInformation.HasSignedAttributes() {
		return nil, nil
	}
	return signerInformation.SignedAttributes.DERImplicitTagged(0), nil
}

// UtilsAddAttribute returns a new attribute table holding an additional attribute of the
// given type and value. It has no CAdESUtils counterpart upstream: it is
// org.bouncycastle.asn1.cms.AttributeTable#add(ASN1ObjectIdentifier, ASN1Encodable), which every
// CAdES augmentation calls and which cmscore.Attributes - a slice, not a Hashtable - does not
// carry a method for.
//
// AttributeTable#add appends a *separate* Attribute of that type rather than merging the value
// into an existing one, and leaves the receiver untouched; both are reproduced here.
func UtilsAddAttribute(attributeTable cmscore.Attributes, attrType asn1.ObjectIdentifier,
	attrValue []byte) cmscore.Attributes {
	extended := make(cmscore.Attributes, 0, len(attributeTable)+1)
	extended = append(extended, attributeTable...)
	return append(extended, cmscore.NewAttribute(attrType, attrValue))
}

// UtilsUnsignedAttributes returns the existing unsigned attributes, or an empty table.
// Port of getUnsignedAttributes(SignerInformation).
func UtilsUnsignedAttributes(signerInformation *cmscore.SignerInfo) cmscore.Attributes {
	if signerInformation == nil {
		return spi.DSSASN1UtilsEmptyIfNull(nil)
	}
	return spi.DSSASN1UtilsEmptyIfNull(signerInformation.UnsignedAttributes)
}

// UtilsSignedAttributes returns the existing signed attributes, or an empty table.
// Port of getSignedAttributes(SignerInformation).
func UtilsSignedAttributes(signerInformation *cmscore.SignerInfo) cmscore.Attributes {
	if signerInformation == nil {
		return spi.DSSASN1UtilsEmptyIfNull(nil)
	}
	return spi.DSSASN1UtilsEmptyIfNull(signerInformation.SignedAttributes)
}

// UtilsAttributesFromByteArray returns an attribute table parsed from its ASN.1 encoded
// representation. Port of getAttributesFromByteArray(byte[]).
//
// Java casts the parsed object to a DLSet and builds an AttributeTable, whose Hashtable loses
// the order of the members; the Go table is a slice and keeps it. The difference is invisible
// downstream, since every re-encoding of an attribute set is a DER SET OF and therefore sorted.
func UtilsAttributesFromByteArray(encodedAttributes []byte) (cmscore.Attributes, error) {
	element, _, err := asn1ber.Parse(encodedAttributes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Error while reading ASN.1 encoded attributes", err)
	}
	if !element.IsConstructed() {
		return nil, model.NewDSSError("Error while reading ASN.1 encoded attributes : not a SET")
	}
	attributes := make(cmscore.Attributes, 0, len(element.Children()))
	for _, child := range element.Children() {
		attribute, err := cmscore.AttributeFromElement(child)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Error while reading ASN.1 encoded attributes", err)
		}
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

// UtilsAddSigningCertificateAttribute appends the signing certificate to the ASN.1 DER
// encoded signed attributes. The certificate is added as either a signing-certificate or a
// signing-certificate-v2 attribute, depending on the digest algorithm being used.
// Port of addSigningCertificateAttribute(ASN1EncodableVector, DigestAlgorithm, CertificateToken).
//
// Java mutates the caller's ASN1EncodableVector; the Go port appends to and returns the
// attribute slice, which is the same operation on the representation this port uses.
func UtilsAddSigningCertificateAttribute(signedAttributes cmscore.Attributes,
	digestAlgorithm enumerations.DigestAlgorithm, signingToken *model.CertificateToken) (cmscore.Attributes, error) {
	issuerSerial := spi.DSSASN1UtilsIssuerSerialForCertificate(signingToken)

	certHash, err := signingToken.Digest(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	// Upstream logs "Adding Certificate Hash {} with algorithm {}".

	var attribute *cmscore.Attribute
	if digestAlgorithm == enumerations.DigestAlgorithmSHA1 {
		essCertID := cadesUtilsESSCertID(certHash, issuerSerial)
		signingCertificate := asn1ber.WriteSequence(asn1ber.WriteSequence(essCertID))
		attribute = cmscore.NewAttribute(spi.OIDIdAaSigningCertificate, signingCertificate)
	} else {
		var hashAlgorithm *spi.AlgorithmIdentifier
		if enumerations.DigestAlgorithmSHA256 != digestAlgorithm {
			// SHA-256 is the ESSCertIDv2 DEFAULT and is left out; see cadesUtilsESSCertIDv2.
			hashAlgorithm, err = spi.DSSASN1UtilsAlgorithmIdentifierForDigest(digestAlgorithm)
			if err != nil {
				return nil, err
			}
		}
		essCertIDv2 := cadesUtilsESSCertIDv2(hashAlgorithm, certHash, issuerSerial)
		signingCertificateV2 := asn1ber.WriteSequence(asn1ber.WriteSequence(essCertIDv2))
		attribute = cmscore.NewAttribute(spi.OIDIdAaSigningCertificateV2, signingCertificateV2)
	}
	return append(signedAttributes, attribute), nil
}

// cadesUtilsESSCertID builds
//
//	ESSCertID ::= SEQUENCE { certHash Hash, issuerSerial IssuerSerial OPTIONAL }
//
// i.e. org.bouncycastle.asn1.ess.ESSCertID(byte[], IssuerSerial).
func cadesUtilsESSCertID(certHash []byte, issuerSerial *spi.IssuerSerial) []byte {
	body := asn1ber.WriteTLV(asn1ber.TagOctetString, certHash)
	if issuerSerial != nil {
		body = append(body, issuerSerial.DER()...)
	}
	return asn1ber.WriteSequence(body)
}

// cadesUtilsESSCertIDv2 builds
//
//	ESSCertIDv2 ::= SEQUENCE {
//	    hashAlgorithm AlgorithmIdentifier DEFAULT {algorithm id-sha256},
//	    certHash      Hash,
//	    issuerSerial  IssuerSerial OPTIONAL }
//
// i.e. org.bouncycastle.asn1.ess.ESSCertIDv2(AlgorithmIdentifier, byte[], IssuerSerial), whose
// toASN1Primitive omits a hashAlgorithm equal to the DEFAULT. A nil hashAlgorithm is the null
// upstream passes for SHA-256.
func cadesUtilsESSCertIDv2(hashAlgorithm *spi.AlgorithmIdentifier, certHash []byte, issuerSerial *spi.IssuerSerial) []byte {
	var body []byte
	if hashAlgorithm != nil && !hashAlgorithm.Equals(cadesUtilsESSCertIDv2DefaultAlgorithm) {
		body = append(body, hashAlgorithm.DER()...)
	}
	body = append(body, asn1ber.WriteTLV(asn1ber.TagOctetString, certHash)...)
	if issuerSerial != nil {
		body = append(body, issuerSerial.DER()...)
	}
	return asn1ber.WriteSequence(body)
}

// cadesUtilsESSCertIDv2DefaultAlgorithm is the DEFAULT of ESSCertIDv2.hashAlgorithm,
// {algorithm id-sha256} without parameters - BouncyCastle's ESSCertIDv2.DEFAULT_ALG_ID.
var cadesUtilsESSCertIDv2DefaultAlgorithm = spi.NewAlgorithmIdentifier(asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1})

// UtilsIsCMSSignedDataEqual compares two CMS documents by their encoded binaries.
// Port of isCMSSignedDataEqual(CMSSignedData, CMSSignedData), whose declared IOException cannot
// be raised here: the encoding is the one the document was parsed from or built with.
func UtilsIsCMSSignedDataEqual(signedData, signedDataToCompare *cmscore.CMS) bool {
	return bytes.Equal(signedData.Encoded(), signedDataToCompare.Encoded())
}

// UtilsSignedAttribute returns the signed attribute with the given oid, when present and
// unique. Port of getSignedAttribute(SignerInformation, ASN1ObjectIdentifier).
//
// Use UtilsSignedAttributesOfType when several attributes are expected: this method
// answers nil - as Java does - when the table carries more than one.
func UtilsSignedAttribute(signerInformation *cmscore.SignerInfo, oid asn1.ObjectIdentifier) *cmscore.Attribute {
	attributes := UtilsSignedAttributesOfType(signerInformation, oid)
	if utils.IsArrayEmpty(attributes) {
		return nil
	}
	if utils.ArraySize(attributes) > 1 {
		// Upstream logs "More than attribute with OID '{}' found in signed attributes table!
		// Value is skipped.".
		return nil
	}
	return attributes[0]
}

// UtilsSignedAttributesOfType returns the signed attributes matching the given oid, or an
// empty slice. Port of the (SignerInformation, ASN1ObjectIdentifier) overload of
// getSignedAttributes.
func UtilsSignedAttributesOfType(signerInformation *cmscore.SignerInfo, oid asn1.ObjectIdentifier) []*cmscore.Attribute {
	return spi.DSSASN1UtilsAsn1Attributes(UtilsSignedAttributes(signerInformation), oid)
}

// UtilsUnsignedAttribute returns the unsigned attribute with the given oid, when present
// and unique. Port of getUnsignedAttribute(SignerInformation, ASN1ObjectIdentifier).
func UtilsUnsignedAttribute(signerInformation *cmscore.SignerInfo, oid asn1.ObjectIdentifier) *cmscore.Attribute {
	attributes := UtilsUnsignedAttributesOfType(signerInformation, oid)
	if utils.IsArrayEmpty(attributes) {
		return nil
	}
	if utils.ArraySize(attributes) > 1 {
		// Upstream logs "More than attribute with OID '{}' found in unsigned attributes table!
		// Value is skipped.".
		return nil
	}
	return attributes[0]
}

// UtilsUnsignedAttributesOfType returns the unsigned attributes matching the given oid, or
// an empty slice. Port of the (SignerInformation, ASN1ObjectIdentifier) overload of
// getUnsignedAttributes.
func UtilsUnsignedAttributesOfType(signerInformation *cmscore.SignerInfo, oid asn1.ObjectIdentifier) []*cmscore.Attribute {
	return spi.DSSASN1UtilsAsn1Attributes(UtilsUnsignedAttributes(signerInformation), oid)
}

// UtilsOriginalDocument returns the original document of the provided CMS.
// Port of getOriginalDocument(CMS, List<DSSDocument>).
func UtilsOriginalDocument(cmsDocument *cms.CMS, detachedDocuments []model.DSSDocument) (model.DSSDocument, error) {
	if cmsDocument == nil {
		panic("CMS shall be provided!")
	}

	if !cmsDocument.IsDetachedSignature() {
		signedContent := cmsDocument.SignedContent()
		if signedContent == nil {
			return nil, model.NewDSSError("No signed content found within enveloping CMS signature!")
		}
		return signedContent, nil

	} else if utils.CollectionSize(detachedDocuments) == 1 {
		return detachedDocuments[0], nil
	}
	return nil, model.NewDSSError("Detached content is not provided or cannot be identified (only one document shall be provided)!")
}

// UtilsContainsATSTv2 reports whether the SignerInformation's unsigned properties contain
// an archive-time-stamp (ATSv2) element.
// Port of containsATSTv2(SignerInformation).
func UtilsContainsATSTv2(signerInformation *cmscore.SignerInfo) bool {
	unsignedAttributes := UtilsUnsignedAttributes(signerInformation)
	for _, attribute := range unsignedAttributes {
		if spi.DSSASN1UtilsIsAttributeOfType(attribute, spi.OIDIdAaEtsArchiveTimestampV2) {
			return true
		}
	}
	return false
}

// UtilsReadSigningDate reads the signing date with respect to RFC 3852, returning the zero
// time.Time (Java: null) when the value is not a date, or is a date encoded against the rule
// below. Port of readSigningDate(ASN1Encodable), whose argument is here the encoding of the
// attribute value.
func UtilsReadSigningDate(attrValue []byte) time.Time {
	if attrValue != nil {
		signingDate := spi.DSSASN1UtilsDate(attrValue)
		if !signingDate.IsZero() {
			/*
			 * RFC 3852 [4] states that "dates between January 1, 1950 and
			 * December 31, 2049 (inclusive) MUST be encoded as UTCTime. Any
			 * dates with year values before 1950 or after 2049 MUST be encoded
			 * as GeneralizedTime".
			 */
			if !signingDate.Before(cadesUtilsJanuary1950) && signingDate.Before(cadesUtilsJanuary2050) &&
				!cadesUtilsIsUTCTime(attrValue) { // must be ASN1UTCTime
				// Upstream logs the RFC 3852 rule and the offending encoding.
				return time.Time{}
			}
			return signingDate
		}
		// Upstream logs "Error when reading signing time. Unrecognized {}".
	}
	return time.Time{}
}

// cadesUtilsIsUTCTime reports whether the encoding is that of an ASN1UTCTime.
func cadesUtilsIsUTCTime(encoded []byte) bool {
	element, rest, err := asn1ber.Parse(encoded)
	return err == nil && len(rest) == 0 && element.IsUniversal(asn1ber.TagUTCTime)
}

// UtilsFindArchiveTimeStampTokens finds the archive time-stamp tokens of an unsigned
// attribute table. Port of findArchiveTimeStampTokens(AttributeTable).
//
// A token that cannot be built is logged and skipped upstream, and skipped here too; the error
// return mirrors the one every other CAdESUtils accessor has and is nil today.
func UtilsFindArchiveTimeStampTokens(unsignedAttributes cmscore.Attributes) ([]*cmscore.TimeStampToken, error) {
	timeStamps := make([]*cmscore.TimeStampToken, 0)
	for _, attribute := range unsignedAttributes {
		if UtilsIsArchiveTimeStampToken(attribute) {
			timeStampToken := UtilsTimeStampToken(attribute)
			if timeStampToken != nil {
				timeStamps = append(timeStamps, timeStampToken)
			}
		}
	}
	return timeStamps, nil
}

// UtilsTimestampOids returns a list of all CMS timestamp identifiers.
// Port of getTimestampOids().
func UtilsTimestampOids() []asn1.ObjectIdentifier {
	return cadesUtilsTimestampOids
}

// UtilsIsArchiveTimeStampToken reports whether the attribute is of an allowed archive
// timestamp type. Port of isArchiveTimeStampToken(Attribute).
func UtilsIsArchiveTimeStampToken(attribute *cmscore.Attribute) bool {
	if attribute != nil && attribute.Type != nil {
		return enumerations.TimestampTypeArchiveTimestamp == UtilsTimestampTypeByOid(attribute.Type)
	}
	return false
}

// UtilsTimestampTypeByOid returns the TimestampType matching the given CMS attribute oid,
// or the empty TimestampType (Java: null) when the OID is not recognised.
// Port of getTimestampTypeByOid(ASN1ObjectIdentifier).
func UtilsTimestampTypeByOid(oid asn1.ObjectIdentifier) enumerations.TimestampType {
	switch {
	case OIDIdAaEtsContentTimestamp.Equal(oid):
		return enumerations.TimestampTypeContentTimestamp
	case OIDIdAaSignatureTimeStampToken.Equal(oid):
		return enumerations.TimestampTypeSignatureTimestamp
	case OIDIdAaEtsCertCRLTimestamp.Equal(oid):
		return enumerations.TimestampTypeValidationDataRefsOnlyTimestamp
	case OIDIdAaEtsEscTimeStamp.Equal(oid):
		return enumerations.TimestampTypeValidationDataTimestamp
	case spi.OIDIdAaEtsArchiveTimestampV2.Equal(oid) || spi.OIDIdAaEtsArchiveTimestampV3.Equal(oid):
		return enumerations.TimestampTypeArchiveTimestamp
	}
	return ""
}

// UtilsAtsHashIndex returns the ats-hash-index table found in the timestamp's unsigned
// properties, whichever version it carries. Port of getAtsHashIndex(AttributeTable).
func UtilsAtsHashIndex(timestampUnsignedAttributes cmscore.Attributes) []byte {
	atsHashIndexVersionIdentifier := UtilsAtsHashIndexVersionIdentifier(timestampUnsignedAttributes)
	return UtilsAtsHashIndexByVersion(timestampUnsignedAttributes, atsHashIndexVersionIdentifier)
}

// UtilsCertificatesHashIndex extracts the certificates hash index of an ats-hash-index
// value. Port of getCertificatesHashIndex(ASN1Sequence).
func UtilsCertificatesHashIndex(atsHashIndexValue []byte) []byte {
	return cadesUtilsHashIndexMember(atsHashIndexValue, 0)
}

// UtilsCRLHashIndex extracts the CRL hash index of an ats-hash-index value.
// Port of getCRLHashIndex(ASN1Sequence).
func UtilsCRLHashIndex(atsHashIndexValue []byte) []byte {
	return cadesUtilsHashIndexMember(atsHashIndexValue, 1)
}

// UtilsUnsignedAttributesHashIndex extracts the unsigned-attributes hash index of an
// ats-hash-index value. Port of getUnsignedAttributesHashIndex(ASN1Sequence).
func UtilsUnsignedAttributesHashIndex(atsHashIndexValue []byte) []byte {
	return cadesUtilsHashIndexMember(atsHashIndexValue, 2)
}

// cadesUtilsHashIndexMember returns the member of an ats-hash-index table at the given base
// index, shifted by one when the table opens with the optional hashIndAlgorithm - the "size > 3"
// test the three Java accessors share. It returns nil where Java returns null, and also where
// Java would raise a ClassCastException or an ArrayIndexOutOfBoundsException on a malformed
// table (an unchecked exception the Go port has no channel for and no caller relies on).
func cadesUtilsHashIndexMember(atsHashIndexValue []byte, index int) []byte {
	if atsHashIndexValue == nil {
		return nil
	}
	element, _, err := asn1ber.Parse(atsHashIndexValue)
	if err != nil || !element.IsConstructed() {
		return nil
	}
	if len(element.Children()) > 3 {
		index++
	}
	if index >= len(element.Children()) {
		return nil
	}
	member := element.Children()[index]
	if !member.IsUniversal(asn1ber.TagSequence) {
		return nil
	}
	return member.Encoded()
}

// UtilsAtsHashIndexByVersion returns the ats-hash-index table of the requested version
// found in the timestamp's unsigned properties, nil when there is none.
// Port of getAtsHashIndexByVersion(AttributeTable, ASN1ObjectIdentifier).
func UtilsAtsHashIndexByVersion(timestampUnsignedAttributes cmscore.Attributes,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier) []byte {
	if timestampUnsignedAttributes != nil && atsHashIndexVersionIdentifier != nil {
		attributes := spi.DSSASN1UtilsAsn1Attributes(timestampUnsignedAttributes, atsHashIndexVersionIdentifier)
		if utils.ArraySize(attributes) == 1 {
			atsHashIndexAttribute := attributes[0]
			attrValue := spi.DSSASN1UtilsAsn1Encodable(atsHashIndexAttribute)
			if attrValue != nil {
				return attrValue.Encoded()
			}
		}
	}
	return nil
}

// UtilsAtsHashIndexVersionIdentifier returns the OID of the AtsHashIndex found in the
// timestamp's unsigned attributes, nil when there is none.
// Port of getAtsHashIndexVersionIdentifier(AttributeTable).
func UtilsAtsHashIndexVersionIdentifier(timestampUnsignedAttributes cmscore.Attributes) asn1.ObjectIdentifier {
	if timestampUnsignedAttributes != nil {
		for _, attribute := range timestampUnsignedAttributes {
			attrType := attribute.Type
			if spi.OIDIdAaATSHashIndex.Equal(attrType) || spi.OIDIdAaATSHashIndexV2.Equal(attrType) ||
				spi.OIDIdAaATSHashIndexV3.Equal(attrType) {
				// Upstream logs "Unsigned attribute of type [{}] found in the timestamp.".
				return attrType
			}
		}
		// Upstream logs "The timestamp unsignedAttributes does not contain ATSHashIndex!".
	}
	return nil
}

// UtilsOctetStringForAtsHashIndex returns the octets of the given attribute, as the
// ats-hash-index version in use defines them.
// Port of getOctetStringForAtsHashIndex(Attribute, ASN1ObjectIdentifier).
func UtilsOctetStringForAtsHashIndex(attribute *cmscore.Attribute,
	atsHashIndexVersionIdentifier asn1.ObjectIdentifier) ([][]byte, error) {
	/*
	 *  id_aa_ATSHashIndexV3 (EN 319 122-1 v1.1.1) -> Each one shall contain the hash
	 *  value of the octets resulting from concatenating the Attribute.attrType field and one of the instances of
	 *  AttributeValue within the Attribute.attrValues within the unsignedAttrs field. One concatenation
	 *  operation shall be performed as indicated above, and the hash value of the obtained result included in
	 *  unsignedAttrsHashIndex
	 */
	if spi.OIDIdAaATSHashIndexV3.Equal(atsHashIndexVersionIdentifier) {
		return UtilsATSHashIndexV3OctetString(attribute.Type, attribute.ValueEncodings())
	}
	/*
	 * id_aa_ATSHashIndex (TS 101 733 v2.2.1) and id_aa_ATSHashIndexV2 (EN 319 122-1 v1.0.0) ->
	 * The field unsignedAttrsHashIndex shall be a sequence of octet strings. Each one shall contain the hash value of
	 * one instance of Attribute within the unsignedAttrs field of the SignerInfo.
	 */
	return [][]byte{attribute.DER()}, nil
}

// UtilsATSHashIndexV3OctetString returns the octets of the given attribute for an
// ATS-Hash-Index-v3 table, i.e. the DER of the attribute type concatenated with the DER of each
// attribute value. Port of getATSHashIndexV3OctetString(ASN1ObjectIdentifier, ASN1Set), whose
// ASN1Set argument is here the encodings of the set's members.
func UtilsATSHashIndexV3OctetString(attributeIdentifier asn1.ObjectIdentifier,
	attributeValues [][]byte) ([][]byte, error) {
	octets := make([][]byte, 0, len(attributeValues))
	attrType := asn1ber.EncodeOID(attributeIdentifier)
	if attrType == nil {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to encode the attribute type '%s'", attributeIdentifier))
	}
	for _, asn1Encodable := range attributeValues {
		derEncoded, err := spi.DSSASN1UtilsDEREncoded(asn1Encodable)
		if err != nil {
			return nil, err
		}
		octets = append(octets, utils.Concat(attrType, derEncoded))
	}
	return octets, nil
}

// UtilsEvidenceRecordOids returns a list of all CMS evidence record identifiers.
// Port of getEvidenceRecordOids().
func UtilsEvidenceRecordOids() []asn1.ObjectIdentifier {
	return cadesUtilsEvidenceRecordOids
}

// UtilsEvidenceRecordIncorporationType returns the evidence record incorporation type of
// the given unsigned attribute OID.
// Port of getEvidenceRecordIncorporationType(ASN1ObjectIdentifier), whose
// UnsupportedOperationException - unchecked, and raised on an OID the caller chose - becomes a
// panic carrying the same message.
func UtilsEvidenceRecordIncorporationType(unsignedAttributeOID asn1.ObjectIdentifier) enumerations.EvidenceRecordIncorporationType {
	switch {
	case spi.OIDIdAaErInternal.Equal(unsignedAttributeOID):
		return enumerations.EvidenceRecordIncorporationTypeInternalEvidenceRecord
	case spi.OIDIdAaErExternal.Equal(unsignedAttributeOID):
		return enumerations.EvidenceRecordIncorporationTypeExternalEvidenceRecord
	}
	panic(fmt.Sprintf("The unsigned attribute with OID '%s' is not supported "+
		"for the evidence record incorporation!", unsignedAttributeOID))
}

// UtilsContainsEvidenceRecord reports whether the signer carries an evidence record
// unsigned attribute. Port of containsEvidenceRecord(SignerInformation).
func UtilsContainsEvidenceRecord(signerInformation *cmscore.SignerInfo) bool {
	if signerInformation != nil && signerInformation.HasUnsignedAttributes() {
		return signerInformation.UnsignedAttributes.Get(spi.OIDIdAaErInternal) != nil ||
			signerInformation.UnsignedAttributes.Get(spi.OIDIdAaErExternal) != nil
	}
	return false
}

// UtilsEvidenceRecordGenerationTime returns the generation time of an evidence record, as
// indicated by the generation time of its first archive time-stamp; the zero time.Time (Java:
// null) when the record carries none.
// Port of getEvidenceRecordGenerationTime(org.bouncycastle.asn1.tsp.EvidenceRecord).
//
// DEVIATION: the argument is the DER encoding of the RFC 4998 EvidenceRecord rather than a
// parsed object - there is no Go counterpart of BouncyCastle's EvidenceRecord yet (the ASN.1
// evidence-record structures belong to a later phase), and this method only reaches
// evidenceRecord.archiveTimeStampSequence[0][0].timeStamp, which the walk below finds:
//
//	EvidenceRecord           ::= SEQUENCE { version INTEGER, digestAlgorithms ..., ...,
//	                                        archiveTimeStampSequence ArchiveTimeStampSequence }
//	ArchiveTimeStampSequence ::= SEQUENCE OF ArchiveTimeStampChain
//	ArchiveTimeStampChain    ::= SEQUENCE OF ArchiveTimeStamp
//	ArchiveTimeStamp         ::= SEQUENCE { ..., timeStamp ContentInfo }
func UtilsEvidenceRecordGenerationTime(evidenceRecord []byte) (time.Time, error) {
	if evidenceRecord == nil {
		return time.Time{}, nil
	}
	element, _, err := asn1ber.Parse(evidenceRecord)
	if err != nil || !element.IsConstructed() || len(element.Children()) == 0 {
		return time.Time{}, nil
	}
	// archiveTimeStampSequence is the last component; cryptoInfos [0] and encryptionInfo [1]
	// are optional and precede it.
	archiveTimeStampSequence := element.Children()[len(element.Children())-1]
	if !archiveTimeStampSequence.IsConstructed() || len(archiveTimeStampSequence.Children()) == 0 {
		return time.Time{}, nil
	}
	archiveTimeStampChain := archiveTimeStampSequence.Children()[0]
	if !archiveTimeStampChain.IsConstructed() || len(archiveTimeStampChain.Children()) == 0 {
		return time.Time{}, nil
	}
	archiveTimestamp := archiveTimeStampChain.Children()[0]
	if !archiveTimestamp.IsConstructed() || len(archiveTimestamp.Children()) == 0 {
		return time.Time{}, nil
	}
	contentInfo := archiveTimestamp.Children()[len(archiveTimestamp.Children())-1]
	timeStampToken, err := cadesUtilsToTimeStampToken(contentInfo)
	if err != nil {
		return time.Time{}, err
	}
	if timeStampToken == nil {
		return time.Time{}, nil
	}
	return timeStampToken.TSTInfo().GenTime, nil
}

// cadesUtilsToTimeStampToken ports the private toTimeStampToken(ContentInfo).
func cadesUtilsToTimeStampToken(contentInfo *asn1ber.Element) (*cmscore.TimeStampToken, error) {
	if contentInfo == nil {
		return nil, nil
	}
	timeStampToken, err := cmscore.TimeStampTokenFromElement(contentInfo)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to build a timestamp token : %s", err), err)
	}
	return timeStampToken, nil
}

// UtilsTimeStampToken creates a TimeStampToken from the provided attribute, nil when it
// cannot be built. Port of getTimeStampToken(Attribute).
func UtilsTimeStampToken(attribute *cmscore.Attribute) *cmscore.TimeStampToken {
	signedData, err := UtilsCMSSignedData(attribute)
	if err == nil && signedData != nil {
		timeStampToken, tokenErr := cmscore.TimeStampTokenFromCMS(signedData)
		if tokenErr == nil {
			return timeStampToken
		}
	}
	// Upstream logs "The given TimeStampToken cannot be created! Reason: [{}]".
	return nil
}

// UtilsCMSSignedData creates a CMS document from the provided attribute, nil when the
// attribute value is an OCTET STRING (which CMS forbids here).
// Port of getCMSSignedData(Attribute).
func UtilsCMSSignedData(attribute *cmscore.Attribute) (*cmscore.CMS, error) {
	value := spi.DSSASN1UtilsAsn1Encodable(attribute)
	if value == nil {
		return nil, nil
	}
	if value.IsUniversal(asn1ber.TagOctetString) {
		// Upstream logs "Illegal content for CMSSignedData (OID : {}) : OCTET STRING is not
		// allowed !".
		return nil, nil
	}
	return cmscore.ParseCMS(value.Encoded())
}

// UtilsEncodedValue returns the encoded value of the attribute.
// Port of getEncodedValue(Attribute).
//
// Java re-encodes the parsed value; the Go port hands out the octets the value arrived in, which
// is what an identity or a digest computed over it depends on and what BouncyCastle's default
// (BER) re-encoding of a parsed object reproduces.
func UtilsEncodedValue(attribute *cmscore.Attribute) ([]byte, error) {
	value := spi.DSSASN1UtilsAsn1Encodable(attribute)
	if value == nil {
		return nil, model.NewDSSError("The attribute carries no value!")
	}
	return value.Encoded(), nil
}

// UtilsSignedDataEncodedOCSPResponse returns the encoded binaries used for an OCSP token
// incorporation within a SignedData.crls attribute.
// Port of getSignedDataEncodedOCSPResponse(byte[], ASN1ObjectIdentifier).
func UtilsSignedDataEncodedOCSPResponse(binaries []byte, objectIdentifier asn1.ObjectIdentifier) ([]byte, error) {
	// Compute the tagged object with the same algorithm BouncyCastle used to create it,
	// see org.bouncycastle.cms.CMSUtils getOthersFromStore().
	//
	//	OtherRevocationInfoFormat ::= SEQUENCE {
	//	    otherRevInfoFormat OBJECT IDENTIFIER,
	//	    otherRevInfo       ANY DEFINED BY otherRevInfoFormat }
	info, err := spi.DSSASN1UtilsDEREncoded(binaries)
	if err != nil {
		return nil, err
	}
	format := asn1ber.EncodeOID(objectIdentifier)
	if format == nil {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to encode the revocation info format '%s'", objectIdentifier))
	}
	// false value specifies an implicit encoding method: the SEQUENCE's identifier octet is
	// replaced by the context-specific [1] one, its content staying as is.
	return asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, append(format, info...)), nil
}

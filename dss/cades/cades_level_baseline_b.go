// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESLevelBaselineB.java (DSS 6.5.RC1).
//
// This is the byte-exactness core of the CAdES port: the DER of every attribute produced here
// is an input of the signature computation, so it must equal, octet for octet, what upstream
// DSS and BouncyCastle write for the same parameters. testdata/baseline-b-attributes.txt is
// upstream's own output for 30 parameter combinations (see testdata/gen/BaselineBFixtures.java)
// and cades_level_baseline_b_kat_test.go checks every one of them.
//
// # ASN.1 representation
//
// BouncyCastle's ASN1Encodable objects become their DER encodings ([]byte), the representation
// dss-spi already chose; internal/asn1ber writes them. The structures assembled by hand here are
// the ones BouncyCastle supplies upstream and that therefore have no Java class to mirror:
// org.bouncycastle.asn1.esf.{SignerAttribute, SignerLocation, CommitmentTypeIndication,
// CommitmentTypeQualifier, SignaturePolicyIdentifier, SignaturePolicyId, OtherHashAlgAndValue,
// SigPolicyQualifiers, SigPolicyQualifierInfo}, org.bouncycastle.asn1.ess.{ContentHints,
// ContentIdentifier}, org.bouncycastle.asn1.x509.{Attribute, Time, UserNotice, NoticeReference,
// DisplayText}. Each is written next to its single consumer, with the ASN.1 module it encodes
// and the BouncyCastle quirks it has to reproduce spelled out in the comment above it. The three
// DSS-specific structures (SignerAttributeV2, SignedAssertion, SignedAssertions) do have Java
// classes, in dss-cms, so they come from the cms package.
//
// # Attribute order
//
// Java collects the attributes in an ASN1EncodableVector and hands it to
// new AttributeTable(...), whose Hashtable loses the order; the vector is then re-expanded and
// DER-encoded as a SET OF, which sorts by encoding. The order this file appends in is therefore
// invisible in the output, and the Go port keeps the upstream append order anyway.
//
// # Deviations
//
//   - A UserNotice carrying only a notice reference (organization + notice numbers, no explicit
//     text) makes upstream call org.bouncycastle.asn1.x509.UserNotice(NoticeReference, (String)
//     null), which raises a NullPointerException inside DisplayText. There is no encoding to
//     reproduce, so BuildSigPolicyQualifiers reports it as an error instead.
//   - A CommitmentTypeQualifier whose content is not ASN.1 is incorporated as a UTF8String built
//     from Java's new String(byte[]); the Go port passes the bytes through unchanged. The two
//     agree for every well-formed UTF-8 content and differ only in how many U+FFFD a malformed
//     byte sequence collapses to.
//   - "parameters instanceof CAdESCounterSignatureParameters", which suppresses the mime-type
//     attribute, has no Go equivalent across an embedded struct; see the counterSignature field.
package cades

import (
	"crypto/rand"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"unicode/utf16"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// The PKCS#9 signed-attribute types this file writes, i.e. the
// org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers constants upstream static-imports. They keep
// their exact Java field name behind the "OID_" prefix.
var (
	// OIDPkcs9AtSigningTime is 1.2.840.113549.1.9.5.
	OIDPkcs9AtSigningTime = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}

	// OIDIdAaContentHint is 1.2.840.113549.1.9.16.2.4.
	OIDIdAaContentHint = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 4}

	// OIDIdAaContentIdentifier is 1.2.840.113549.1.9.16.2.7.
	OIDIdAaContentIdentifier = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 7}

	// OIDIdAaEtsSigPolicyId is 1.2.840.113549.1.9.16.2.15.
	OIDIdAaEtsSigPolicyId = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 15}

	// OIDIdAaEtsCommitmentType is 1.2.840.113549.1.9.16.2.16.
	OIDIdAaEtsCommitmentType = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 16}

	// OIDIdAaEtsSignerLocation is 1.2.840.113549.1.9.16.2.17.
	OIDIdAaEtsSignerLocation = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 17}

	// OIDIdAaEtsSignerAttr is 1.2.840.113549.1.9.16.2.18.
	OIDIdAaEtsSignerAttr = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 2, 18}

	// OIDIdSpqEtsUri is 1.2.840.113549.1.9.16.5.1.
	OIDIdSpqEtsUri = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 5, 1}

	// OIDIdSpqEtsUnotice is 1.2.840.113549.1.9.16.5.2.
	OIDIdSpqEtsUnotice = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 5, 2}
)

// CAdESLevelBaselineB holds the CAdES-B signature profile; it supports the inclusion of the
// mandatory signed id_aa_ets_sigPolicyId attribute as specified in ETSI TS 101 733 V1.8.1,
// clause 5.8.1.
type CAdESLevelBaselineB struct {
	// documentToSign is the document to be signed.
	documentToSign model.DSSDocument

	// counterSignature stands in for Java's "parameters instanceof
	// CAdESCounterSignatureParameters" test in addMimeType. Go's embedding does not preserve
	// the dynamic type of an embedded *CAdESSignatureParameters, so the fact that the
	// parameters describe a counter-signature travels with the profile instead: the
	// counter-signature builder asks CMSForCAdESBuilderHelper for a counter-signature profile
	// (SetCounterSignature), which is the single upstream call site of the instanceof test.
	counterSignature bool
}

// NewCAdESLevelBaselineB is the default constructor for CAdESLevelBaselineB.
// Port of the no-arg constructor, i.e. of CAdESLevelBaselineB(null).
func NewCAdESLevelBaselineB() *CAdESLevelBaselineB {
	return NewCAdESLevelBaselineBWithDocument(nil)
}

// NewCAdESLevelBaselineBWithDocument is the constructor for CAdESLevelBaselineB with a
// documentToSign. Port of CAdESLevelBaselineB(DSSDocument).
func NewCAdESLevelBaselineBWithDocument(documentToSign model.DSSDocument) *CAdESLevelBaselineB {
	return &CAdESLevelBaselineB{documentToSign: documentToSign}
}

// SetCounterSignature marks the profile as building the signed attributes of a counter
// signature, which suppresses the mime-type attribute. See the counterSignature field: this
// stands in for the "parameters instanceof CAdESCounterSignatureParameters" test of
// #addMimeType, which Go's embedding cannot express.
func (b *CAdESLevelBaselineB) SetCounterSignature(counterSignature bool) {
	b.counterSignature = counterSignature
}

// UnsignedAttributes returns the table of unsigned properties, which is empty at level B.
// Port of getUnsignedAttributes().
func (b *CAdESLevelBaselineB) UnsignedAttributes() cmscore.Attributes {
	return cmscore.Attributes{}
}

// SignedAttributes generates and returns the signed attributes table.
// Port of getSignedAttributes(CAdESSignatureParameters).
func (b *CAdESLevelBaselineB) SignedAttributes(parameters *CAdESSignatureParameters) (cmscore.Attributes, error) {
	if utils.IsArrayNotEmpty(parameters.SignedData()) {
		// Upstream logs "Using explicit SignedAttributes from parameter".
		return CAdESUtilsAttributesFromByteArray(parameters.SignedData())
	}

	signedAttributes := cmscore.Attributes{}
	return b.AddSignedAttributes(parameters, signedAttributes)
}

// AddSignedAttributes adds the signed attributes to the signedAttributes table.
// Port of the protected #addSignedAttributes(CAdESSignatureParameters, ASN1EncodableVector);
// Java mutates the vector in place, the Go port appends to and returns the slice.
func (b *CAdESLevelBaselineB) AddSignedAttributes(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	var err error
	if signedAttributes, err = b.AddSigningCertificateAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSigningTimeAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignerAttribute(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignaturePolicyId(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddContentHints(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddMimeType(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddContentIdentifier(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddCommitmentType(parameters, signedAttributes); err != nil {
		return nil, err
	}
	if signedAttributes, err = b.AddSignerLocation(parameters, signedAttributes); err != nil {
		return nil, err
	}
	return b.AddContentTimestamps(parameters, signedAttributes)
}

// AddSignerAttribute adds the signer-attributes attribute.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
// 5.11.3 signer-attributes Attribute
// NOTE 1: Only a single signer-attributes can be used.
//
// The signer-attributes attribute specifies additional attributes of the signer (e.g. role).
// It may be either:
//   - claimed attributes of the signer; or
//   - certified attributes of the signer.
//
// The signer-attributes attribute shall be a signed attribute.
//
// Port of the protected #addSignerAttribute.
func (b *CAdESLevelBaselineB) AddSignerAttribute(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	claimedSignerRoles := parameters.BLevel().ClaimedSignerRoles()
	if claimedSignerRoles != nil {

		claimedAttributes := make([][]byte, 0, len(claimedSignerRoles))
		for _, claimedSignerRole := range claimedSignerRoles {
			roles := cadesLevelBaselineBUTF8String(claimedSignerRole)
			idAaEtsSignerAttr := cadesLevelBaselineBX509Attribute(spi.OIDIdAtRole, roles)
			claimedAttributes = append(claimedAttributes, idAaEtsSignerAttr)
		}
		var signerAttributes *cmscore.Attribute
		if !parameters.IsEn319122() {
			signerAttributes = cmscore.NewAttribute(OIDIdAaEtsSignerAttr,
				cadesLevelBaselineBSignerAttribute(claimedAttributes))
		} else {
			signerAttributes = cmscore.NewAttribute(spi.OIDIdAaEtsSignerAttrV2,
				cms.NewSignerAttributeV2FromClaimedAttributes(claimedAttributes).DER())
		}
		return append(signedAttributes, signerAttributes), nil
	}

	signedAssertions := parameters.BLevel().SignedAssertions()
	if signedAssertions != nil && parameters.IsEn319122() {
		assertionsToAdd := make([]*cms.SignedAssertion, 0, len(signedAssertions))
		for _, signedAssertion := range signedAssertions {
			assertionsToAdd = append(assertionsToAdd, cms.NewSignedAssertion(signedAssertion))
		}

		if len(assertionsToAdd) != 0 {
			signerAttributes := cmscore.NewAttribute(spi.OIDIdAaEtsSignerAttrV2,
				cms.NewSignerAttributeV2FromSignedAssertions(cms.NewSignedAssertions(assertionsToAdd)).DER())
			signedAttributes = append(signedAttributes, signerAttributes)
		}
	}
	return signedAttributes, nil
}

// AddSigningTimeAttribute adds a signing time attribute.
// Port of the protected #addSigningTimeAttribute.
func (b *CAdESLevelBaselineB) AddSigningTimeAttribute(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	signingDate := parameters.BLevel().SigningDate()
	if signingDate != nil {
		attribute := cmscore.NewAttribute(OIDPkcs9AtSigningTime, cadesLevelBaselineBTime(*signingDate))
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// AddSignerLocation adds the signer-location attribute.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
// 5.11.2 signer-location Attribute
// The signer-location attribute specifies a mnemonic for an address associated with the signer
// at a particular geographical (e.g. city) location. The mnemonic is registered in the country
// in which the signer is located and is used in the provision of the Public Telegram Service
// (according to Recommendation ITU-T F.1 [11]). The signer-location attribute shall be a signed
// attribute.
//
// Port of the protected #addSignerLocation.
func (b *CAdESLevelBaselineB) AddSignerLocation(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	signerLocationParameter := parameters.BLevel().SignerLocation()
	if signerLocationParameter != nil && !signerLocationParameter.IsEmpty() {

		var country, locality []byte
		if signerLocationParameter.Country() != "" {
			country = cadesLevelBaselineBUTF8String(signerLocationParameter.Country())
		}
		if signerLocationParameter.Locality() != "" {
			locality = cadesLevelBaselineBUTF8String(signerLocationParameter.Locality())
		}
		derSequencePostalAddress := cadesLevelBaselineBPostalAddressSequence(signerLocationParameter.PostalAddress())
		signerLocation, err := cadesLevelBaselineBSignerLocation(country, locality,
			derSequencePostalAddress, len(signerLocationParameter.PostalAddress()))
		if err != nil {
			return nil, err
		}
		attribute := cmscore.NewAttribute(OIDIdAaEtsSignerLocation, signerLocation)
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// cadesLevelBaselineBPostalAddressSequence ports the private getPostalAddressSequence.
func cadesLevelBaselineBPostalAddressSequence(postalAddressParameter []string) []byte {
	var derSequencePostalAddress []byte
	if utils.IsCollectionNotEmpty(postalAddressParameter) {
		var postalAddress []byte
		for _, addressLine := range postalAddressParameter {
			postalAddress = append(postalAddress, cadesLevelBaselineBUTF8String(addressLine)...)
		}
		derSequencePostalAddress = asn1ber.WriteSequence(postalAddress)
	}
	return derSequencePostalAddress
}

// AddCommitmentType adds the commitment-type-indication attribute.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
//
// 5.11.1 commitment-type-indication Attribute
// There may be situations where a signer wants to explicitly indicate to a verifier that by
// signing the data, it illustrates a type of commitment on behalf of the signer. The
// commitment-type-indication attribute conveys such information.
//
// Port of the protected #addCommitmentType.
func (b *CAdESLevelBaselineB) AddCommitmentType(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	commitmentTypeIndications := parameters.BLevel().CommitmentTypeIndications()
	if utils.IsCollectionNotEmpty(commitmentTypeIndications) {
		size := len(commitmentTypeIndications)
		asn1Encodables := make([][]byte, size)
		for ii := 0; ii < size; ii++ {
			commitmentType := commitmentTypeIndications[ii]
			if utils.IsStringEmpty(commitmentType.OID()) {
				return nil, model.NewDSSError("The commitmentTypeIndication OID must be defined for CAdES creation!")
			}
			objectIdentifier, err := asn1ber.OIDFromString(commitmentType.OID())
			if err != nil {
				return nil, err
			}
			qualifiers, err := b.CommitmentQualifiers(commitmentType)
			if err != nil {
				return nil, err
			}
			encodedIdentifier := asn1ber.EncodeOID(objectIdentifier)
			if encodedIdentifier == nil {
				return nil, model.NewDSSError(fmt.Sprintf("string %s not a valid OID", commitmentType.OID()))
			}
			asn1Encodables[ii] = cadesLevelBaselineBCommitmentTypeIndication(encodedIdentifier, qualifiers)
		}
		attribute := cmscore.NewAttribute(OIDIdAaEtsCommitmentType, asn1Encodables...)
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// CommitmentQualifiers creates the SEQUENCE of CommitmentQualifiers of a commitment type, nil
// when it carries none.
//
//	CommitmentTypeQualifier ::= SEQUENCE {
//	 commitmentQualifierId COMMITMENT-QUALIFIER.&id,
//	 qualifier COMMITMENT-QUALIFIER.&Qualifier OPTIONAL
//	}
//	COMMITMENT-QUALIFIER ::= CLASS {
//	 &id OBJECT IDENTIFIER UNIQUE,
//	 &Qualifier OPTIONAL }
//	WITH SYNTAX {
//	 COMMITMENT-QUALIFIER-ID &id
//	 [COMMITMENT-TYPE &Qualifier] }
//
// Port of the protected #getCommitmentQualifiers(CommitmentType).
func (b *CAdESLevelBaselineB) CommitmentQualifiers(commitmentType enumerations.CommitmentType) ([]byte, error) {
	var qualifiers []byte
	if commonCommitmentType, ok := commitmentType.(*model.CommonCommitmentType); ok {
		commitmentTypeQualifiers := commonCommitmentType.CommitmentTypeQualifiers()
		if utils.IsArrayNotEmpty(commitmentTypeQualifiers) {
			var vector []byte
			for _, commitmentQualifier := range commitmentTypeQualifiers {
				if commitmentQualifier == nil {
					panic("CommitmentTypeQualifier cannot be null!")
				}

				if utils.IsStringEmpty(commitmentQualifier.Oid()) {
					return nil, model.NewDSSError("CommitmentTypeQualifier OID cannot be null for CAdES!")
				}
				commitmentIdentifier, err := asn1ber.OIDFromString(commitmentQualifier.Oid())
				if err != nil {
					return nil, err
				}
				encodedIdentifier := asn1ber.EncodeOID(commitmentIdentifier)
				if encodedIdentifier == nil {
					return nil, model.NewDSSError(fmt.Sprintf("string %s not a valid OID", commitmentQualifier.Oid()))
				}

				content := commitmentQualifier.Content()
				if content == nil {
					return nil, model.NewDSSError("CommitmentTypeQualifier content cannot be null!")
				}
				var qualifier []byte
				binaries, err := spi.DSSUtilsToByteArrayOfDocument(content)
				if err != nil {
					return nil, err
				}
				if spi.DSSASN1UtilsIsAsn1Encoded(binaries) {
					if qualifier, err = spi.DSSASN1UtilsToASN1Primitive(binaries); err != nil {
						return nil, err
					}
				} else {
					// Upstream logs "None ASN.1 encoded CommitmentTypeQualifier has been
					// provided. Incorporate as DERUTF8String.".
					qualifier = cadesLevelBaselineBUTF8String(string(binaries))
				}
				vector = append(vector, cadesLevelBaselineBCommitmentTypeQualifier(encodedIdentifier, qualifier)...)
			}
			qualifiers = asn1ber.WriteSequence(vector)
		}
	}
	return qualifiers, nil
}

// AddContentTimestamps adds the content-timestamp attributes.
//
// A content time-stamp allows a time-stamp token of the data to be signed to be incorporated
// into the signed information. It provides proof of the existence of the data before the
// signature was created.
//
// A content time-stamp attribute is the time-stamp token of the signed data content before it
// is signed. This attribute is a signed attribute. Its object identifier is:
// id-aa-ets-contentTimestamp OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
// rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) id-aa(2) 20}
//
// Content time-stamp attribute values have ASN.1 type ContentTimestamp:
// ContentTimestamp ::= TimeStampToken
//
// The value of messageImprint of TimeStampToken (as described in RFC 3161) is the hash of the
// message digest as defined in ETSI standard 101733 v.2.2.1, clause 5.6.1.
//
// NOTE: content-time-stamp indicates that the signed information was formed before the date
// included in the content-time-stamp.
// NOTE (bis): There is a small difference in treatment between the content-time-stamp and the
// archive-timestamp (ATSv2) when the signature is attached. In that case, the content-time-stamp
// is computed on the raw data (without ASN.1 tag and length) whereas the archive-timestamp is
// computed on data as read.
//
// Port of the protected #addContentTimestamps.
func (b *CAdESLevelBaselineB) AddContentTimestamps(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {

	if utils.IsCollectionNotEmpty(parameters.ContentTimestamps()) {

		contentTimestamps := parameters.ContentTimestamps()
		for _, contentTimestamp := range contentTimestamps {

			asn1Object, err := spi.DSSASN1UtilsToASN1Primitive(contentTimestamp.Encoded())
			if err != nil {
				return nil, err
			}
			attribute := cmscore.NewAttribute(OIDIdAaEtsContentTimestamp, asn1Object)
			signedAttributes = append(signedAttributes, attribute)
		}
	}
	return signedAttributes, nil
}

// AddContentHints adds the content-hints attribute.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
//
// 5.10.3 content-hints Attribute
// The content-hints attribute provides information on the innermost signed content of a
// multi-layer message where one content is encapsulated in another.
// The syntax of the content-hints attribute type of the ES is as defined in ESS (RFC 2634 [5]).
// When used to indicate the precise format of the data to be presented to the user, the
// following rules apply:
//   - the contentType indicates the type of the associated content. It is an object identifier
//     (i.e. a unique string of integers) assigned by an authority that defines the content type;
//     and
//   - when the contentType is id-data the contentDescription shall define the presentation
//     format; the format may be defined by MIME types.
//
// When the format of the content is defined by MIME types, the following rules apply:
//   - the contentType shall be id-data as defined in CMS (RFC 3852 [4]);
//   - the contentDescription shall be used to indicate the encoding of the data, in accordance
//     with the rules defined RFC 2045 [6]; see annex F for an example of structured contents and
//     MIME.
//
// NOTE 1: id-data OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs7(7) 1 }.
// NOTE 2: contentDescription is optional in ESS (RFC 2634 [5]). It may be used to complement
// contentTypes defined elsewhere; such definitions are outside the scope of the present document.
//
// Port of the protected #addContentHints.
func (b *CAdESLevelBaselineB) AddContentHints(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	if utils.IsStringNotBlank(parameters.ContentHintsType()) {

		contentHintsType, err := asn1ber.OIDFromString(parameters.ContentHintsType())
		if err != nil {
			return nil, err
		}
		encodedType := asn1ber.EncodeOID(contentHintsType)
		if encodedType == nil {
			return nil, model.NewDSSError(fmt.Sprintf("string %s not a valid OID", parameters.ContentHintsType()))
		}
		contentHintsDescriptionString := parameters.ContentHintsDescription()
		var contentHintsDescription []byte
		if !utils.IsStringBlank(contentHintsDescriptionString) {
			contentHintsDescription = cadesLevelBaselineBUTF8String(contentHintsDescriptionString)
		}
		// "text/plain";
		// "1.2.840.113549.1.7.1";

		contentHints := cadesLevelBaselineBContentHints(encodedType, contentHintsDescription)
		attribute := cmscore.NewAttribute(OIDIdAaContentHint, contentHints)
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// AddContentIdentifier adds the content-identifier attribute.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
//
// 5.10.2 content-identifier Attribute
// The content-identifier attribute provides an identifier for the signed content, for use when a
// reference may be later required to that content; for example, in the content-reference
// attribute in other signed data sent later. The content-identifier shall be a signed attribute.
// content-identifier attribute type values for the ES have an ASN.1 type ContentIdentifier, as
// defined in ESS (RFC 2634 [5]).
//
// The minimal content-identifier attribute should contain a concatenation of user-specific
// identification information (such as a user name or public keying material identification
// information), a GeneralizedTime string, and a random number.
//
// Port of the protected #addContentIdentifier.
func (b *CAdESLevelBaselineB) AddContentIdentifier(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	contentIdentifierPrefix := parameters.ContentIdentifierPrefix()
	if utils.IsStringNotBlank(contentIdentifierPrefix) {
		if utils.IsStringBlank(parameters.ContentIdentifierSuffix()) {
			suffix, err := cadesLevelBaselineBContentIdentifierSuffix(time.Now())
			if err != nil {
				return nil, err
			}
			parameters.SetContentIdentifierSuffix(suffix)
		}
		contentIdentifierString := contentIdentifierPrefix + parameters.ContentIdentifierSuffix()
		// ContentIdentifier ::= OCTET STRING
		contentIdentifier := asn1ber.WriteTLV(asn1ber.TagOctetString, []byte(contentIdentifierString))
		attribute := cmscore.NewAttribute(OIDIdAaContentIdentifier, contentIdentifier)
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// cadesLevelBaselineBContentIdentifierSuffix builds the default content-identifier suffix, i.e.
// new ASN1GeneralizedTime(new Date()).getTimeString() followed by new SecureRandom().nextLong().
//
// BouncyCastle formats the date as "yyyyMMddHHmmss'Z'" in UTC; Java appends the random long in
// its signed decimal form, which is what strconv.FormatInt writes.
func cadesLevelBaselineBContentIdentifierSuffix(now time.Time) (string, error) {
	var buffer [8]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", model.NewDSSErrorMessageCause("Unable to generate a content identifier suffix", err)
	}
	random := int64(binary.BigEndian.Uint64(buffer[:]))
	return now.UTC().Format("20060102150405") + "Z" + strconv.FormatInt(random, 10), nil
}

// AddSignaturePolicyId adds a signature policy identifier.
// Port of the protected #addSignaturePolicyId.
func (b *CAdESLevelBaselineB) AddSignaturePolicyId(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	policy := parameters.BLevel().SignaturePolicy()
	if policy != nil {

		policyId := policy.Id()
		var sigPolicy []byte

		if utils.IsStringEmpty(policyId) { // implicit
			// SignaturePolicyIdentifier is a CHOICE whose signaturePolicyImplied alternative
			// BouncyCastle encodes as a NULL.
			sigPolicy = asn1ber.DERNull

		} else { // explicit
			derOIPolicyId, err := asn1ber.OIDFromString(policyId)
			if err != nil {
				return nil, err
			}
			encodedPolicyId := asn1ber.EncodeOID(derOIPolicyId)
			if encodedPolicyId == nil {
				return nil, model.NewDSSError(fmt.Sprintf("string %s not a valid OID", policyId))
			}
			oid, err := asn1ber.OIDFromString(policy.DigestAlgorithm().OID())
			if err != nil {
				return nil, err
			}
			algorithmIdentifier := spi.NewAlgorithmIdentifier(oid)
			otherHashAlgAndValue := cadesLevelBaselineBOtherHashAlgAndValue(algorithmIdentifier, policy.DigestValue())

			if policy.IsSPQualifierPresent() {
				sigPolicyQualifiers, err := b.BuildSigPolicyQualifiers(policy)
				if err != nil {
					return nil, err
				}
				sigPolicy = cadesLevelBaselineBSignaturePolicyId(encodedPolicyId, otherHashAlgAndValue, sigPolicyQualifiers)
			} else {
				sigPolicy = cadesLevelBaselineBSignaturePolicyId(encodedPolicyId, otherHashAlgAndValue, nil)
			}
		}

		attribute := cmscore.NewAttribute(OIDIdAaEtsSigPolicyId, sigPolicy)
		signedAttributes = append(signedAttributes, attribute)
	}
	return signedAttributes, nil
}

// BuildSigPolicyQualifiers builds
//
//	SigPolicyQualifiers ::= SEQUENCE SIZE (1..MAX) OF SigPolicyQualifierInfo
//
// Port of the private buildSigPolicyQualifiers(Policy). Exported because AddSignaturePolicyId,
// which upstream may override, is the only caller and Go has no package-private visibility
// between a type's own methods.
func (b *CAdESLevelBaselineB) BuildSigPolicyQualifiers(policy *model.Policy) ([]byte, error) {
	var qualifierInfos []byte

	spUri := policy.Spuri()
	if utils.IsStringNotEmpty(spUri) {
		spuriQualifier := cadesLevelBaselineBSigPolicyQualifierInfo(OIDIdSpqEtsUri,
			asn1ber.WriteTLV(asn1ber.TagIA5String, []byte(policy.Spuri())))
		qualifierInfos = append(qualifierInfos, spuriQualifier...)
	}

	userNotice := policy.UserNotice()
	if userNotice != nil && !userNotice.IsEmpty() {
		if err := spi.DSSUtilsAssertSPUserNoticeConfigurationValid(userNotice); err != nil {
			return nil, err
		}

		var noticeReference []byte
		var explicitText string

		organization := userNotice.Organization()
		noticeNumbers := userNotice.NoticeNumbers()
		if utils.IsStringNotEmpty(organization) && len(noticeNumbers) > 0 {
			numbers := spi.DSSUtilsToBigIntegerList(noticeNumbers)
			noticeReference = cadesLevelBaselineBNoticeReference(organization, numbers)
		}
		if utils.IsStringNotEmpty(userNotice.ExplicitText()) {
			explicitText = userNotice.ExplicitText()
		}

		if explicitText == "" {
			// UPSTREAM DEFECT, reported rather than reproduced: this configuration passes
			// DSSUtils#assertSPUserNoticeConfigurationValid and then reaches
			// org.bouncycastle.asn1.x509.UserNotice(NoticeReference, (String) null), which
			// hands the null to DisplayText(String) and raises a NullPointerException. There
			// are no ground-truth bytes for it.
			return nil, model.NewDSSError("The UserNotice explicitText shall be defined " +
				"for a CAdES signature policy qualifier!")
		}
		asn1UserNotice := cadesLevelBaselineBUserNotice(noticeReference, explicitText)
		userNoticeQualifier := cadesLevelBaselineBSigPolicyQualifierInfo(OIDIdSpqEtsUnotice, asn1UserNotice)
		qualifierInfos = append(qualifierInfos, userNoticeQualifier...)
	}

	spDocSpecification := policy.SpDocSpecification()
	if spDocSpecification != nil && utils.IsStringNotEmpty(spDocSpecification.Id()) {
		spDocSpecificationId, err := spi.DSSASN1UtilsBuildSPDocSpecificationID(spDocSpecification.Id())
		if err != nil {
			return nil, err
		}
		spDocSpecificationQualifier := cadesLevelBaselineBSigPolicyQualifierInfo(
			spi.OIDIdSpDocSpecification, spDocSpecificationId)
		qualifierInfos = append(qualifierInfos, spDocSpecificationQualifier...)
	}

	return asn1ber.WriteSequence(qualifierInfos), nil
}

// AddSigningCertificateAttribute adds a signing-certificate attribute.
// Port of the protected #addSigningCertificateAttribute.
func (b *CAdESLevelBaselineB) AddSigningCertificateAttribute(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	if parameters.SigningCertificate() == nil && parameters.GenerateTBSWithoutCertificate() {
		// Upstream logs "Signing certificate not available and must be added to signed
		// attributes later".
		return signedAttributes, nil
	}

	return CAdESUtilsAddSigningCertificateAttribute(signedAttributes, parameters.DigestAlgorithm(),
		parameters.SigningCertificate())
}

// AddMimeType adds a MimeType attribute.
// Port of the protected #addMimeType.
func (b *CAdESLevelBaselineB) AddMimeType(parameters *CAdESSignatureParameters,
	signedAttributes cmscore.Attributes) (cmscore.Attributes, error) {
	if b.counterSignature {
		// Java: parameters instanceof CAdESCounterSignatureParameters.
		return signedAttributes, nil
	}
	if utils.IsStringNotBlank(parameters.ContentHintsType()) {
		return signedAttributes, nil
	}

	var mimeType enumerations.MimeType = enumerations.MimeTypeEnumBinary
	if b.documentToSign != nil && b.documentToSign.MimeType() != nil {
		mimeType = b.documentToSign.MimeType()
	}
	mimeTypeDerString := cadesLevelBaselineBUTF8String(mimeType.MimeTypeString())
	attribute := cmscore.NewAttribute(spi.OIDIdAaEtsMimeType, mimeTypeDerString)
	return append(signedAttributes, attribute), nil
}

// -----------------------------------------------------------------------------
// The BouncyCastle ASN.1 structures this file writes.
//
// Each one reproduces the toASN1Primitive of the named class, DER throughout - which is what
// BouncyCastle produces here too, since every value ends up inside a DERSet.
// -----------------------------------------------------------------------------

// cadesLevelBaselineBUTF8String encodes a DERUTF8String.
func cadesLevelBaselineBUTF8String(value string) []byte {
	return asn1ber.WriteTLV(asn1ber.TagUTF8String, []byte(value))
}

// cadesLevelBaselineBExplicit wraps the value in an explicitly tagged context-specific element,
// i.e. new DERTaggedObject(tagNo, value) / new DERTaggedObject(true, tagNo, value).
func cadesLevelBaselineBExplicit(tagNo int, value []byte) []byte {
	return asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|byte(tagNo), value)
}

// cadesLevelBaselineBTime encodes
//
//	Time ::= CHOICE { utcTime UTCTime, generalTime GeneralizedTime }
//
// i.e. org.bouncycastle.asn1.x509.Time(Date): the date is formatted as "yyyyMMddHHmmss" in UTC
// with a trailing "Z", and a year outside [1950, 2049] switches from UTCTime (which drops the
// century) to GeneralizedTime. Sub-second precision is lost, as SimpleDateFormat drops it.
func cadesLevelBaselineBTime(date time.Time) []byte {
	formatted := date.UTC().Format("20060102150405") + "Z"
	year := date.UTC().Year()
	if year < 1950 || year > 2049 {
		return asn1ber.WriteTLV(asn1ber.TagGeneralizedTime, []byte(formatted))
	}
	return asn1ber.WriteTLV(asn1ber.TagUTCTime, []byte(formatted[2:]))
}

// cadesLevelBaselineBX509Attribute encodes
//
//	Attribute ::= SEQUENCE { attrType OBJECT IDENTIFIER, attrValues SET OF AttributeValue }
//
// i.e. org.bouncycastle.asn1.x509.Attribute with a single value. The attrType OIDs used here are
// DSS's own constants, which are always encodable.
func cadesLevelBaselineBX509Attribute(attrType asn1.ObjectIdentifier, value []byte) []byte {
	body := asn1ber.EncodeOID(attrType)
	body = append(body, asn1ber.WriteTLV(asn1ber.TagSet|asn1ber.Constructed, value)...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBSignerAttribute encodes
//
//	SignerAttribute ::= SEQUENCE OF CHOICE {
//	    claimedAttributes   [0] ClaimedAttributes,
//	    certifiedAttributes [1] CertifiedAttributes }
//	ClaimedAttributes ::= SEQUENCE OF Attribute
//
// i.e. org.bouncycastle.asn1.esf.SignerAttribute(Attribute[]), whose toASN1Primitive wraps the
// claimed attributes in an explicitly tagged [0].
func cadesLevelBaselineBSignerAttribute(claimedAttributes [][]byte) []byte {
	var attributes []byte
	for _, claimedAttribute := range claimedAttributes {
		attributes = append(attributes, claimedAttribute...)
	}
	return asn1ber.WriteSequence(cadesLevelBaselineBExplicit(0, asn1ber.WriteSequence(attributes)))
}

// cadesLevelBaselineBSignerLocation encodes
//
//	SignerLocation ::= SEQUENCE {
//	    countryName   [0] DirectoryString OPTIONAL,
//	    localityName  [1] DirectoryString OPTIONAL,
//	    postalAddress [2] PostalAddress OPTIONAL }
//	PostalAddress ::= SEQUENCE SIZE(1..6) OF DirectoryString
//
// i.e. org.bouncycastle.asn1.esf.SignerLocation(ASN1UTF8String, ASN1UTF8String, ASN1Sequence),
// whose constructor rejects a postal address of more than 6 lines and whose toASN1Primitive
// tags the three fields explicitly.
func cadesLevelBaselineBSignerLocation(country, locality, postalAddress []byte, postalAddressSize int) ([]byte, error) {
	if postalAddress != nil && postalAddressSize > 6 {
		return nil, model.NewDSSError("postal address sequence no more than 6 strings")
	}
	var body []byte
	if country != nil {
		body = append(body, cadesLevelBaselineBExplicit(0, country)...)
	}
	if locality != nil {
		body = append(body, cadesLevelBaselineBExplicit(1, locality)...)
	}
	if postalAddress != nil {
		body = append(body, cadesLevelBaselineBExplicit(2, postalAddress)...)
	}
	return asn1ber.WriteSequence(body), nil
}

// cadesLevelBaselineBCommitmentTypeIndication encodes
//
//	CommitmentTypeIndication ::= SEQUENCE {
//	    commitmentTypeId         CommitmentTypeIdentifier,
//	    commitmentTypeQualifier  SEQUENCE SIZE (1..MAX) OF CommitmentTypeQualifier OPTIONAL }
//
// i.e. org.bouncycastle.asn1.esf.CommitmentTypeIndication.
func cadesLevelBaselineBCommitmentTypeIndication(commitmentTypeID, qualifiers []byte) []byte {
	body := append([]byte{}, commitmentTypeID...)
	if qualifiers != nil {
		body = append(body, qualifiers...)
	}
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBCommitmentTypeQualifier encodes
//
//	CommitmentTypeQualifier ::= SEQUENCE {
//	    commitmentQualifierId COMMITMENT-QUALIFIER.&id,
//	    qualifier             COMMITMENT-QUALIFIER.&Qualifier OPTIONAL }
//
// i.e. org.bouncycastle.asn1.esf.CommitmentTypeQualifier.
func cadesLevelBaselineBCommitmentTypeQualifier(commitmentQualifierID, qualifier []byte) []byte {
	body := append([]byte{}, commitmentQualifierID...)
	if qualifier != nil {
		body = append(body, qualifier...)
	}
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBContentHints encodes
//
//	ContentHints ::= SEQUENCE {
//	    contentDescription UTF8String (SIZE (1..MAX)) OPTIONAL,
//	    contentType        ContentType }
//
// i.e. org.bouncycastle.asn1.ess.ContentHints: the OPTIONAL description comes FIRST, per
// RFC 2634 clause 2.9.
func cadesLevelBaselineBContentHints(contentType, contentDescription []byte) []byte {
	var body []byte
	if contentDescription != nil {
		body = append(body, contentDescription...)
	}
	body = append(body, contentType...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBOtherHashAlgAndValue encodes
//
//	OtherHashAlgAndValue ::= SEQUENCE {
//	    hashAlgorithm AlgorithmIdentifier,
//	    hashValue     OtherHashValue }
//	OtherHashValue ::= OCTET STRING
//
// i.e. org.bouncycastle.asn1.esf.OtherHashAlgAndValue.
func cadesLevelBaselineBOtherHashAlgAndValue(hashAlgorithm *spi.AlgorithmIdentifier, hashValue []byte) []byte {
	body := hashAlgorithm.DER()
	body = append(body, asn1ber.WriteTLV(asn1ber.TagOctetString, hashValue)...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBSignaturePolicyId encodes
//
//	SignaturePolicyId ::= SEQUENCE {
//	    sigPolicyId         SigPolicyId,
//	    sigPolicyHash       SigPolicyHash,
//	    sigPolicyQualifiers SEQUENCE SIZE (1..MAX) OF SigPolicyQualifierInfo OPTIONAL }
//
// i.e. org.bouncycastle.asn1.esf.SignaturePolicyId. The SignaturePolicyIdentifier CHOICE around
// it is transparent: its signaturePolicyId alternative encodes as the SEQUENCE itself.
func cadesLevelBaselineBSignaturePolicyId(sigPolicyID, sigPolicyHash, sigPolicyQualifiers []byte) []byte {
	body := append([]byte{}, sigPolicyID...)
	body = append(body, sigPolicyHash...)
	if sigPolicyQualifiers != nil {
		body = append(body, sigPolicyQualifiers...)
	}
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBSigPolicyQualifierInfo encodes
//
//	SigPolicyQualifierInfo ::= SEQUENCE {
//	    sigPolicyQualifierId SigPolicyQualifierId,
//	    sigQualifier         ANY DEFINED BY sigPolicyQualifierId }
//
// i.e. org.bouncycastle.asn1.esf.SigPolicyQualifierInfo. The three qualifier OIDs used here are
// constants, hence always encodable.
func cadesLevelBaselineBSigPolicyQualifierInfo(qualifierID asn1.ObjectIdentifier, qualifier []byte) []byte {
	body := asn1ber.EncodeOID(qualifierID)
	body = append(body, qualifier...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBUserNotice encodes
//
//	UserNotice ::= SEQUENCE {
//	    noticeRef    NoticeReference OPTIONAL,
//	    explicitText DisplayText OPTIONAL }
//
// i.e. org.bouncycastle.asn1.x509.UserNotice(NoticeReference, String).
func cadesLevelBaselineBUserNotice(noticeRef []byte, explicitText string) []byte {
	var body []byte
	if noticeRef != nil {
		body = append(body, noticeRef...)
	}
	body = append(body, cadesLevelBaselineBDisplayText(explicitText)...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBNoticeReference encodes
//
//	NoticeReference ::= SEQUENCE {
//	    organization  DisplayText,
//	    noticeNumbers SEQUENCE OF INTEGER }
//
// i.e. org.bouncycastle.asn1.x509.NoticeReference(String, Vector).
func cadesLevelBaselineBNoticeReference(organization string, noticeNumbers []*big.Int) []byte {
	body := cadesLevelBaselineBDisplayText(organization)
	var numbers []byte
	for _, noticeNumber := range noticeNumbers {
		numbers = append(numbers, asn1ber.EncodeInteger(noticeNumber)...)
	}
	body = append(body, asn1ber.WriteSequence(numbers)...)
	return asn1ber.WriteSequence(body)
}

// cadesLevelBaselineBDisplayTextMaximumSize is DisplayText.DISPLAY_TEXT_MAXIMUM_SIZE, the 200
// characters RFC 3280 allows and beyond which BouncyCastle silently truncates.
const cadesLevelBaselineBDisplayTextMaximumSize = 200

// cadesLevelBaselineBDisplayText encodes
//
//	DisplayText ::= CHOICE { ia5String IA5String, visibleString VisibleString,
//	                         bmpString BMPString, utf8String UTF8String }
//
// i.e. org.bouncycastle.asn1.x509.DisplayText(String), which picks the utf8String alternative
// and truncates the text to DISPLAY_TEXT_MAXIMUM_SIZE characters first. Java counts UTF-16 code
// units, which is what the conversion below reproduces.
func cadesLevelBaselineBDisplayText(text string) []byte {
	units := utf16.Encode([]rune(text))
	if len(units) > cadesLevelBaselineBDisplayTextMaximumSize {
		units = units[:cadesLevelBaselineBDisplayTextMaximumSize]
		text = string(utf16.Decode(units))
	}
	return cadesLevelBaselineBUTF8String(text)
}

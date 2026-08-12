// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSASN1Utils.java (DSS 6.5.RC1).
//
// This file replaces BouncyCastle's ASN.1 layer. Two conventions follow from that:
//
//   - Java passes around ASN1Encodable/ASN1Primitive objects; the Go port represents an
//     ASN.1 value by its encoding ([]byte), because that is what the surrounding code
//     digests, compares and serialises. Wherever identity or a digest depends on the
//     bytes, the original bytes are carried through untouched and never re-encoded.
//   - The BER/DER/DL re-encoding BouncyCastle performs in getEncoded(String) is
//     implemented here (see dssASN1UtilsElement): DER output must be byte-exact, since it
//     feeds signature verification and archive-timestamp hashing.
//
// NOT PORTED (deferred to the CMS/timestamp phase, which owns the types they take):
// getEncoded(TimeStampToken), getEncoded(CMSSignedData), getDEREncoded(TimeStampToken),
// getDEREncoded(CMSSignedData), getAsn1Encodable(Attribute), getAsn1Attributes,
// isEmpty(AttributeTable), emptyIfNull(AttributeTable), isAttributeOfType,
// getTimeStampTokenGenerationTime, getRevocationValues, getCertificateRef(OtherCertID),
// getFirstSignerInformation, toSignerIdentifier(SignerId), getX509CertificateHolder and
// getCertificate(X509CertificateHolder). See the TODO markers next to their Go neighbours.
package spi

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// Universal ASN.1 tag numbers used by this file.
const (
	dssASN1UtilsTagBoolean         = 0x01
	dssASN1UtilsTagInteger         = 0x02
	dssASN1UtilsTagBitString       = 0x03
	dssASN1UtilsTagOctetString     = 0x04
	dssASN1UtilsTagNull            = 0x05
	dssASN1UtilsTagOID             = 0x06
	dssASN1UtilsTagObjectDesc      = 0x07
	dssASN1UtilsTagUTF8String      = 0x0C
	dssASN1UtilsTagSequence        = 0x10
	dssASN1UtilsTagSet             = 0x11
	dssASN1UtilsTagNumericString   = 0x12
	dssASN1UtilsTagPrintableString = 0x13
	dssASN1UtilsTagT61String       = 0x14
	dssASN1UtilsTagVideotexString  = 0x15
	dssASN1UtilsTagIA5String       = 0x16
	dssASN1UtilsTagUTCTime         = 0x17
	dssASN1UtilsTagGeneralizedTime = 0x18
	dssASN1UtilsTagGraphicString   = 0x19
	dssASN1UtilsTagVisibleString   = 0x1A
	dssASN1UtilsTagGeneralString   = 0x1B
	dssASN1UtilsTagUniversalString = 0x1C
	dssASN1UtilsTagBMPString       = 0x1E
)

// ASN.1 identifier octet bit masks.
const (
	dssASN1UtilsClassMask   = 0xC0
	dssASN1UtilsConstructed = 0x20
	dssASN1UtilsTagMask     = 0x1F
)

// dssASN1UtilsDERNull is the DER encoding of ASN.1 NULL, i.e. BouncyCastle's DERNull.INSTANCE.
var dssASN1UtilsDERNull = []byte{dssASN1UtilsTagNull, 0x00}

// X.520 attribute type OIDs the human-readable-name helpers look for; they mirror the
// org.bouncycastle.asn1.x500.style.BCStyle constants DSSASN1Utils names.
var (
	dssASN1UtilsOIDCN        = asn1.ObjectIdentifier{2, 5, 4, 3}
	dssASN1UtilsOIDSurname   = asn1.ObjectIdentifier{2, 5, 4, 4}
	dssASN1UtilsOIDO         = asn1.ObjectIdentifier{2, 5, 4, 10}
	dssASN1UtilsOIDOU        = asn1.ObjectIdentifier{2, 5, 4, 11}
	dssASN1UtilsOIDName      = asn1.ObjectIdentifier{2, 5, 4, 41}
	dssASN1UtilsOIDGivenName = asn1.ObjectIdentifier{2, 5, 4, 42}
	dssASN1UtilsOIDPseudonym = asn1.ObjectIdentifier{2, 5, 4, 65}
)

// -----------------------------------------------------------------------------
// ASN.1 structures the ported methods produce and consume.
// -----------------------------------------------------------------------------

// AlgorithmIdentifier is the X.509 structure
//
//	AlgorithmIdentifier ::= SEQUENCE {
//	    algorithm   OBJECT IDENTIFIER,
//	    parameters  ANY DEFINED BY algorithm OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.AlgorithmIdentifier. Parameters holds the complete
// DER encoding of the parameters element, or nil when they are absent - the distinction
// matters, since an explicit NULL and an absent parameter are different encodings.
type AlgorithmIdentifier struct {
	// Algorithm is the algorithm OID.
	Algorithm asn1.ObjectIdentifier
	// Parameters is the DER encoding of the parameters, nil when absent.
	Parameters []byte
}

// NewAlgorithmIdentifier builds an AlgorithmIdentifier without parameters.
func NewAlgorithmIdentifier(algorithm asn1.ObjectIdentifier) *AlgorithmIdentifier {
	return &AlgorithmIdentifier{Algorithm: algorithm}
}

// NewAlgorithmIdentifierWithParameters builds an AlgorithmIdentifier carrying the given
// DER-encoded parameters.
func NewAlgorithmIdentifierWithParameters(algorithm asn1.ObjectIdentifier, parameters []byte) *AlgorithmIdentifier {
	return &AlgorithmIdentifier{Algorithm: algorithm, Parameters: parameters}
}

// ParseAlgorithmIdentifier decodes an AlgorithmIdentifier from its DER encoding.
// Port of AlgorithmIdentifier.getInstance(Object).
func ParseAlgorithmIdentifier(der []byte) (*AlgorithmIdentifier, error) {
	element, rest, err := dssASN1UtilsParse(der)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("extra data found after the AlgorithmIdentifier")
	}
	return dssASN1UtilsAlgorithmIdentifierFromElement(element)
}

// dssASN1UtilsAlgorithmIdentifierFromElement decodes an already parsed SEQUENCE.
func dssASN1UtilsAlgorithmIdentifierFromElement(element *dssASN1UtilsElement) (*AlgorithmIdentifier, error) {
	if !element.isUniversal(dssASN1UtilsTagSequence) || !element.constructed {
		return nil, errors.New("AlgorithmIdentifier is not a SEQUENCE")
	}
	if len(element.children) == 0 || len(element.children) > 2 {
		return nil, fmt.Errorf("AlgorithmIdentifier: bad sequence size: %d", len(element.children))
	}
	oid, err := element.children[0].objectIdentifier()
	if err != nil {
		return nil, err
	}
	identifier := &AlgorithmIdentifier{Algorithm: oid}
	if len(element.children) == 2 {
		identifier.Parameters = element.children[1].derEncoded()
	}
	return identifier, nil
}

// DER returns the DER encoding of the AlgorithmIdentifier.
func (a *AlgorithmIdentifier) DER() []byte {
	body := dssASN1UtilsEncodeOID(a.Algorithm)
	if a.Parameters != nil {
		body = append(body, a.Parameters...)
	}
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, body)
}

// Equals compares two AlgorithmIdentifiers by their DER encoding.
func (a *AlgorithmIdentifier) Equals(other *AlgorithmIdentifier) bool {
	if a == other {
		return true
	}
	if a == nil || other == nil {
		return false
	}
	return bytes.Equal(a.DER(), other.DER())
}

// GeneralName is one alternative of the X.509 GeneralName CHOICE, reduced to what
// IssuerSerial needs. TagNo is the context-specific tag number (4 for directoryName) and
// Name is the DER encoding of the chosen alternative's value.
type GeneralName struct {
	// TagNo is the CHOICE alternative, e.g. 4 for directoryName.
	TagNo int
	// Name is the DER encoding of the value, e.g. an X.501 Name SEQUENCE for tag 4.
	Name []byte
}

// DER returns the DER encoding of the GeneralName.
//
// directoryName (tag 4) wraps a CHOICE and is therefore explicitly tagged, as
// org.bouncycastle.asn1.x509.GeneralName encodes it; every other alternative is implicitly
// tagged, i.e. the value's identifier octet is replaced by the context-specific tag.
func (g *GeneralName) DER() []byte {
	if g.TagNo == 4 {
		return dssASN1UtilsWriteTLV(0xA0|byte(g.TagNo), g.Name)
	}
	element, _, err := dssASN1UtilsParse(g.Name)
	if err != nil {
		return nil
	}
	identifier := byte(0x80) | byte(g.TagNo)
	if element.constructed {
		identifier |= dssASN1UtilsConstructed
	}
	return dssASN1UtilsWriteTLV(identifier, element.derContent())
}

// IssuerSerial is the X.509 attribute-certificate structure
//
//	IssuerSerial ::= SEQUENCE {
//	    issuer   GeneralNames,
//	    serial   CertificateSerialNumber,
//	    issuerUID UniqueIdentifier OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.IssuerSerial. Only the two fields DSS reads are
// modelled; a third field found while parsing is preserved in IssuerUID.
type IssuerSerial struct {
	// Issuer is the GeneralNames sequence identifying the issuer.
	Issuer []GeneralName
	// Serial is the certificate serial number.
	Serial *big.Int
	// IssuerUID is the DER encoding of the optional issuerUID BIT STRING, nil when absent.
	IssuerUID []byte
}

// DER returns the DER encoding of the IssuerSerial.
func (i *IssuerSerial) DER() []byte {
	var names []byte
	for index := range i.Issuer {
		names = append(names, i.Issuer[index].DER()...)
	}
	body := dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, names)
	body = append(body, dssASN1UtilsEncodeInteger(i.Serial)...)
	if i.IssuerUID != nil {
		body = append(body, i.IssuerUID...)
	}
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, body)
}

// -----------------------------------------------------------------------------
// The ported DSSASN1Utils methods.
// -----------------------------------------------------------------------------

// DSSASN1UtilsToASN1Primitive validates that the given bytes hold exactly one ASN.1 element
// and returns that element's encoding. Port of toASN1Primitive(byte[]); an ASN1Primitive is
// represented by its bytes in this port, so the method degenerates to a parse-and-return.
func DSSASN1UtilsToASN1Primitive(bytes []byte) ([]byte, error) {
	element, rest, err := dssASN1UtilsParse(bytes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot convert binaries to ASN1Primitive", err)
	}
	if len(rest) != 0 {
		return nil, model.NewDSSError("Cannot convert binaries to ASN1Primitive : extra data found after the object")
	}
	return element.encoded, nil
}

// DSSASN1UtilsIsDEROctetStringNull reports whether the given DER OCTET STRING encapsulates
// an ASN.1 NULL. Port of isDEROctetStringNull(DEROctetString), whose argument is here the
// complete encoding of the OCTET STRING rather than a parsed object.
func DSSASN1UtilsIsDEROctetStringNull(derOctetString []byte) bool {
	element, rest, err := dssASN1UtilsParse(derOctetString)
	if err != nil || len(rest) != 0 || !element.isUniversal(dssASN1UtilsTagOctetString) {
		return false
	}
	asn1Null, err := DSSASN1UtilsToASN1Primitive(element.octets())
	if err != nil {
		return false
	}
	return bytes.Equal(dssASN1UtilsDERNull, asn1Null)
}

// DSSASN1UtilsDEREncoded returns the DER encoding of the given ASN.1 value.
// Port of getDEREncoded(ASN1Encodable) and getDEREncoded(byte[]), which coincide here.
//
// The conversion is a real BER-to-DER normalisation: indefinite lengths become definite,
// constructed OCTET/BIT STRINGs are collapsed into their primitive form, non-minimal length
// encodings are rewritten, BOOLEAN TRUE becomes 0xFF and the components of every SET are
// sorted by their encoding.
func DSSASN1UtilsDEREncoded(asn1Encodable []byte) ([]byte, error) {
	element, err := dssASN1UtilsParseOne(asn1Encodable, "DER")
	if err != nil {
		return nil, err
	}
	return element.derEncoded(), nil
}

// DSSASN1UtilsDLEncoded returns the DL encoding of the given ASN.1 value: every length is
// definite, but the structure is otherwise preserved - constructed strings stay constructed
// and SET components keep their order. Port of getDLEncoded(ASN1Encodable)/getDLEncoded(byte[]).
func DSSASN1UtilsDLEncoded(asn1Encodable []byte) ([]byte, error) {
	element, err := dssASN1UtilsParseOne(asn1Encodable, "DL")
	if err != nil {
		return nil, err
	}
	return element.dlEncoded(), nil
}

// DSSASN1UtilsBEREncoded returns the BER encoding of the given ASN.1 value, preserving the
// indefinite lengths the input used. Port of getBEREncoded(ASN1Encodable).
func DSSASN1UtilsBEREncoded(asn1Encodable []byte) ([]byte, error) {
	element, err := dssASN1UtilsParseOne(asn1Encodable, "BER")
	if err != nil {
		return nil, err
	}
	return element.berEncoded(), nil
}

// dssASN1UtilsParseOne parses exactly one element, reporting failures the way upstream's
// private getEncoded(ASN1Encodable, String) does.
func dssASN1UtilsParseOne(binaries []byte, encoding string) (*dssASN1UtilsElement, error) {
	element, rest, err := dssASN1UtilsParse(binaries)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to encode to %s", encoding), err)
	}
	if len(rest) != 0 {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to encode to %s. Reason : extra data found after the object", encoding))
	}
	return element, nil
}

// DSSASN1UtilsDEREncodedTimestampBinary returns the DER encoding of a TimestampBinary.
// Port of getDEREncoded(TimestampBinary).
func DSSASN1UtilsDEREncodedTimestampBinary(timestampBinary *model.TimestampBinary) ([]byte, error) {
	return DSSASN1UtilsDEREncoded(timestampBinary.Bytes())
}

// DSSASN1UtilsToDate converts a DER GENERALIZED TIME to a time.Time.
// Port of toDate(ASN1GeneralizedTime).
func DSSASN1UtilsToDate(asn1Date []byte) (time.Time, error) {
	element, rest, err := dssASN1UtilsParse(asn1Date)
	if err != nil || len(rest) != 0 || !element.isUniversal(dssASN1UtilsTagGeneralizedTime) {
		return time.Time{}, model.NewDSSError("Cannot parse Date : not an ASN1GeneralizedTime")
	}
	date, err := dssASN1UtilsParseGeneralizedTime(string(element.content))
	if err != nil {
		return time.Time{}, model.NewDSSErrorMessageCause("Cannot parse Date", err)
	}
	return date, nil
}

// DSSASN1UtilsToString reads the value of a DER OCTET STRING as a string.
// Port of toString(ASN1OctetString), i.e. new String(value.getOctets()).
func DSSASN1UtilsToString(value []byte) (string, error) {
	element, rest, err := dssASN1UtilsParse(value)
	if err != nil {
		return "", err
	}
	if len(rest) != 0 || !element.isUniversal(dssASN1UtilsTagOctetString) {
		return "", errors.New("not an ASN1OctetString")
	}
	return string(element.octets()), nil
}

// DSSASN1UtilsAsn1SequenceFromDerOctetString returns the DER encoding of the SEQUENCE
// encapsulated in the given DER OCTET STRING.
// Port of getAsn1SequenceFromDerOctetString(byte[]).
func DSSASN1UtilsAsn1SequenceFromDerOctetString(binaries []byte) ([]byte, error) {
	content, err := dssASN1UtilsDEROctetStringContent(binaries)
	if err != nil {
		return nil, err
	}
	element, _, err := dssASN1UtilsParse(content)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the ASN1Sequence", err)
	}
	if !element.isUniversal(dssASN1UtilsTagSequence) {
		return nil, model.NewDSSError("Unable to retrieve the ASN1Sequence : the encapsulated object is not a SEQUENCE")
	}
	return element.encoded, nil
}

// DSSASN1UtilsAsn1IntegerFromDerOctetString returns the INTEGER encapsulated in the given
// DER OCTET STRING. Port of getAsn1IntegerFromDerOctetString(byte[]).
func DSSASN1UtilsAsn1IntegerFromDerOctetString(binaries []byte) (*big.Int, error) {
	content, err := dssASN1UtilsDEROctetStringContent(binaries)
	if err != nil {
		return nil, err
	}
	element, _, err := dssASN1UtilsParse(content)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the ASN1Integer", err)
	}
	if !element.isUniversal(dssASN1UtilsTagInteger) {
		return nil, model.NewDSSError("Unable to retrieve the ASN1Integer : the encapsulated object is not an INTEGER")
	}
	return element.integer(), nil
}

// dssASN1UtilsDEROctetStringContent ports the private getDEROctetStringContent(byte[]).
func dssASN1UtilsDEROctetStringContent(binaries []byte) ([]byte, error) {
	element, _, err := dssASN1UtilsParse(binaries)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the DEROctetString content", err)
	}
	if !element.isUniversal(dssASN1UtilsTagOctetString) {
		return nil, model.NewDSSError("Unable to retrieve the DEROctetString content : not an OCTET STRING")
	}
	return element.octets(), nil
}

// TODO(phase-3): getAsn1Encodable(Attribute), getAsn1Attributes(AttributeTable,
// ASN1ObjectIdentifier), isEmpty(AttributeTable), emptyIfNull(AttributeTable) and
// isAttributeOfType(Attribute, ASN1ObjectIdentifier) operate on the CMS attribute types.
// They are only called from the CMS*Source classes, which are deferred with the CMS phase.

// DSSASN1UtilsAsn1SignaturePolicyDigest computes the digest of an ASN.1 signature policy
// (used in CAdES). Port of getAsn1SignaturePolicyDigest(DigestAlgorithm, byte[]).
//
// TS 101 733 5.8.1: if the signature policy is defined using ASN.1, the hash is calculated
// on the value without the outer type and length fields.
func DSSASN1UtilsAsn1SignaturePolicyDigest(digestAlgorithm enumerations.DigestAlgorithm, policyBytes []byte) ([]byte, error) {
	element, _, err := dssASN1UtilsParse(policyBytes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot convert binaries to ASN1Primitive", err)
	}
	if !element.constructed || len(element.children) < 2 {
		return nil, model.NewDSSError("The signature policy shall be a SEQUENCE of at least 2 elements")
	}
	signPolicyHashAlgIdentifier, err := dssASN1UtilsAlgorithmIdentifierFromElement(element.children[0])
	if err != nil {
		return nil, err
	}
	signPolicyInfo := element.children[1]

	hashAlgorithmDEREncoded := signPolicyHashAlgIdentifier.DER()
	signPolicyInfoDEREncoded := signPolicyInfo.derEncoded()
	return dssASN1UtilsDigest(digestAlgorithm, hashAlgorithmDEREncoded, signPolicyInfoDEREncoded)
}

// DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex returns the algorithm identifier found in
// the provided ats-hash-index table, or nil when the table carries none.
// Port of getAlgorithmIdentifier(ASN1Sequence).
func DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(atsHashIndexValue []byte) *AlgorithmIdentifier {
	element, _, err := dssASN1UtilsParse(atsHashIndexValue)
	if err != nil || !element.constructed {
		return nil
	}
	if len(element.children) <= 3 {
		return nil
	}
	const algorithmIndex = 0
	candidate := element.children[algorithmIndex]
	switch {
	case candidate.isUniversal(dssASN1UtilsTagSequence):
		identifier, err := dssASN1UtilsAlgorithmIdentifierFromElement(candidate)
		if err != nil {
			return nil
		}
		return identifier
	case candidate.isUniversal(dssASN1UtilsTagOID):
		// TODO (upstream, 16/11/2014): the relevance and usefulness of this test case must
		// be checked (do signatures like this exist?)
		oid, err := candidate.objectIdentifier()
		if err != nil {
			return nil
		}
		return NewAlgorithmIdentifier(oid)
	}
	return nil
}

// dssASN1UtilsDigestAlgorithmIdentifierParameters records, per digest algorithm OID, whether
// BouncyCastle's DefaultDigestAlgorithmIdentifierFinder emits an explicit NULL parameter.
// The table is a known-answer capture of that finder (bcpkix 1.78.1): the NIST OIDs (SHA-2,
// SHA-3, SHAKE) and WHIRLPOOL carry no parameters, the older OIDs carry DERNull.
var dssASN1UtilsDigestAlgorithmIdentifierParameters = map[enumerations.DigestAlgorithm]bool{
	enumerations.DigestAlgorithm_SHA1:         true,
	enumerations.DigestAlgorithm_SHA224:       false,
	enumerations.DigestAlgorithm_SHA256:       false,
	enumerations.DigestAlgorithm_SHA384:       false,
	enumerations.DigestAlgorithm_SHA512:       false,
	enumerations.DigestAlgorithm_SHA3_224:     false,
	enumerations.DigestAlgorithm_SHA3_256:     false,
	enumerations.DigestAlgorithm_SHA3_384:     false,
	enumerations.DigestAlgorithm_SHA3_512:     false,
	enumerations.DigestAlgorithm_SHAKE128:     false,
	enumerations.DigestAlgorithm_SHAKE256:     false,
	enumerations.DigestAlgorithm_SHAKE256_512: false,
	enumerations.DigestAlgorithm_RIPEMD160:    true,
	enumerations.DigestAlgorithm_MD2:          true,
	enumerations.DigestAlgorithm_MD5:          true,
	enumerations.DigestAlgorithm_WHIRLPOOL:    false,
}

// DSSASN1UtilsAlgorithmIdentifierForDigest returns the ASN.1 algorithm identifier structure
// of a digest algorithm. Port of getAlgorithmIdentifier(DigestAlgorithm).
//
// SHAKE256-512 is the documented special case (DSS-3651): it needs the output length as an
// INTEGER parameter. Every other algorithm reproduces
// DefaultDigestAlgorithmIdentifierFinder#find, see the table above.
func DSSASN1UtilsAlgorithmIdentifierForDigest(digestAlgorithm enumerations.DigestAlgorithm) (*AlgorithmIdentifier, error) {
	oid, err := dssASN1UtilsObjectIdentifier(digestAlgorithm.OID())
	if err != nil {
		return nil, err
	}
	if enumerations.DigestAlgorithm_SHAKE256_512 == digestAlgorithm {
		// Special case, requiring the parameter definition.
		return NewAlgorithmIdentifierWithParameters(oid, dssASN1UtilsEncodeInteger(big.NewInt(512))), nil
	}
	withNull, known := dssASN1UtilsDigestAlgorithmIdentifierParameters[digestAlgorithm]
	if withNull || !known {
		// An algorithm outside the table gets the DERNull parameters the finder defaults to.
		return NewAlgorithmIdentifierWithParameters(oid, dssASN1UtilsDERNull), nil
	}
	return NewAlgorithmIdentifier(oid), nil
}

// DSSASN1UtilsDEROctetStrings returns the content octets of every OCTET STRING found in the
// given DER SEQUENCE, which is how the ats-hash-index tables carry lists of hash values.
// Port of getDEROctetStrings(ASN1Sequence); Java returns the DEROctetString objects, whose
// only use at the call sites is getOctets().
func DSSASN1UtilsDEROctetStrings(asn1Sequence []byte) ([][]byte, error) {
	derOctetStrings := make([][]byte, 0)
	if asn1Sequence == nil {
		return derOctetStrings, nil
	}
	element, _, err := dssASN1UtilsParse(asn1Sequence)
	if err != nil {
		return nil, err
	}
	if !element.constructed {
		return nil, errors.New("not an ASN1Sequence")
	}
	for _, child := range element.children {
		if !child.isUniversal(dssASN1UtilsTagOctetString) {
			return nil, errors.New("the sequence holds an element that is not an OCTET STRING")
		}
		derOctetStrings = append(derOctetStrings, child.octets())
	}
	return derOctetStrings, nil
}

// DSSASN1UtilsComputeSkiFromCert computes the SHA-1 hash of the certificate's public key.
// Port of computeSkiFromCert(CertificateToken).
func DSSASN1UtilsComputeSkiFromCert(certificateToken *model.CertificateToken) ([]byte, error) {
	return DSSASN1UtilsComputeSkiFromCertPublicKey(certificateToken.PublicKey())
}

// DSSASN1UtilsComputeSkiFromCertPublicKey computes the SHA-1 hash of the given public key,
// i.e. of the subjectPublicKey BIT STRING of its SubjectPublicKeyInfo.
// Port of computeSkiFromCertPublicKey(PublicKey).
func DSSASN1UtilsComputeSkiFromCertPublicKey(publicKey *model.PublicKey) ([]byte, error) {
	element, _, err := dssASN1UtilsParse(publicKey.Encoded())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute ski from public key", err)
	}
	if !element.constructed || len(element.children) < 2 {
		return nil, model.NewDSSError("Unable to compute ski from public key : malformed SubjectPublicKeyInfo")
	}
	item := element.children[1]
	if !item.isUniversal(dssASN1UtilsTagBitString) {
		return nil, model.NewDSSError("Unable to compute ski from public key : subjectPublicKey is not a BIT STRING")
	}
	return dssASN1UtilsDigest(enumerations.DigestAlgorithm_SHA1, item.bitStringOctets())
}

// DSSASN1UtilsIsSkiEqual reports whether the provided ski matches the one computed from the
// certificate's public key. Port of isSkiEqual(byte[], CertificateToken).
func DSSASN1UtilsIsSkiEqual(ski []byte, certificateToken *model.CertificateToken) bool {
	certSki, err := DSSASN1UtilsComputeSkiFromCert(certificateToken)
	if err != nil {
		return false
	}
	return bytes.Equal(certSki, ski)
}

// TODO(phase-3): getX509CertificateHolder(CertificateToken) and
// getCertificate(X509CertificateHolder) only convert to and from BouncyCastle's
// X509CertificateHolder, which the CMS layer needs and which has no Go counterpart: a
// CertificateToken already carries the certificate's DER (CertificateToken.Encoded()).

// TODO(phase-3): toSignerIdentifier(SignerId) reads a CMS SignerIdentifier.

// DSSASN1UtilsToX500Principal converts a DER-encoded X.501 Name into an X500Principal.
// Port of toX500Principal(X500Name).
func DSSASN1UtilsToX500Principal(x500Name []byte) (*model.X500Principal, error) {
	if x500Name == nil {
		return nil, nil
	}
	principal, err := model.NewX500Principal(x500Name)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot extract X500Principal!", err)
	}
	return principal, nil
}

// NOT PORTED: getX500PrincipalOrNull(String) builds an X500Principal from its RFC 2253
// string form, using the X520Attributes keyword map. It has no call site in DSS 6.5.RC1
// outside of tests, and it would require an RFC 2253 parser and DN encoder that the port
// does not otherwise need; model.X500Principal is deliberately construction-from-DER only,
// so that the encoding it reports is always the one it was given.

// DSSASN1UtilsToSignerIdentifier builds a SignerIdentifier from the issuer, the serial
// number and the subject key identifier.
// Port of toSignerIdentifier(X500Principal, BigInteger, byte[]).
func DSSASN1UtilsToSignerIdentifier(issuerX500Principal *model.X500Principal, serialNumber *big.Int, ski []byte) *SignerIdentifier {
	signerIdentifier := NewSignerIdentifier()
	signerIdentifier.SetIssuerName(issuerX500Principal)
	signerIdentifier.SetSerialNumber(serialNumber)
	signerIdentifier.SetSki(ski)
	return signerIdentifier
}

// DSSASN1UtilsIssuerSerialForCertificate returns a new IssuerSerial based on the
// certificate token. Port of getIssuerSerial(CertificateToken).
func DSSASN1UtilsIssuerSerialForCertificate(certToken *model.CertificateToken) *IssuerSerial {
	issuerX500Name := certToken.Certificate().RawIssuer
	generalName := GeneralName{TagNo: 4, Name: issuerX500Name}
	serialNumber := certToken.Certificate().SerialNumber
	return &IssuerSerial{Issuer: []GeneralName{generalName}, Serial: serialNumber}
}

// DSSASN1UtilsIssuerSerial parses an IssuerSerial, returning nil when the binaries cannot be
// decoded. Port of getIssuerSerial(byte[]).
func DSSASN1UtilsIssuerSerial(binaries []byte) *IssuerSerial {
	element, _, err := dssASN1UtilsParse(binaries)
	if err != nil || !element.constructed || len(element.children) < 2 {
		// Upstream logs "Unable to decode IssuerSerialV2 textContent '{}' : {}".
		return nil
	}
	names := element.children[0]
	if !names.constructed {
		return nil
	}
	issuerSerial := &IssuerSerial{}
	for _, name := range names.children {
		if name.class != 0x80 {
			return nil
		}
		tagNo := int(name.tagNumber)
		value := name.derContent()
		if tagNo == 4 {
			// directoryName is explicitly tagged: the content is the Name itself.
			inner, _, err := dssASN1UtilsParse(value)
			if err != nil {
				return nil
			}
			value = inner.derEncoded()
		}
		issuerSerial.Issuer = append(issuerSerial.Issuer, GeneralName{TagNo: tagNo, Name: value})
	}
	serial := element.children[1]
	if !serial.isUniversal(dssASN1UtilsTagInteger) {
		return nil
	}
	issuerSerial.Serial = serial.integer()
	if len(element.children) > 2 {
		issuerSerial.IssuerUID = element.children[2].derEncoded()
	}
	return issuerSerial
}

// DSSASN1UtilsToSignerIdentifierFromIssuerSerial transforms an IssuerSerial into a
// SignerIdentifier, returning nil when it cannot be read.
// Port of toSignerIdentifier(IssuerSerial).
func DSSASN1UtilsToSignerIdentifierFromIssuerSerial(issuerAndSerial *IssuerSerial) *SignerIdentifier {
	if issuerAndSerial == nil {
		return nil
	}
	signerIdentifier := NewSignerIdentifier()
	if len(issuerAndSerial.Issuer) == 1 {
		principal, err := DSSASN1UtilsToX500Principal(issuerAndSerial.Issuer[0].Name)
		if err != nil {
			// Upstream logs "Unable to read the IssuerSerial object" and returns null.
			return nil
		}
		signerIdentifier.SetIssuerName(principal)
	}
	// Upstream logs "More than one GeneralName" and leaves the issuer unset when the
	// GeneralNames carries more than one name.
	if issuerAndSerial.Serial != nil {
		signerIdentifier.SetSerialNumber(issuerAndSerial.Serial)
	}
	return signerIdentifier
}

// DSSASN1UtilsX500PrincipalAreEquals compares two X500Principals, first by the JDK's own
// canonical comparison and then by their attribute maps.
// Port of x500PrincipalAreEquals(X500Principal, X500Principal).
func DSSASN1UtilsX500PrincipalAreEquals(firstX500Principal, secondX500Principal *model.X500Principal) bool {
	if firstX500Principal == nil || secondX500Principal == nil {
		return false
	}
	if firstX500Principal.Equals(secondX500Principal) {
		return true
	}
	firstMap := DSSASN1UtilsAttributeMap(firstX500Principal)
	secondMap := DSSASN1UtilsAttributeMap(secondX500Principal)
	if firstMap == nil || secondMap == nil {
		// Java raises a DSSException on a malformed name rather than reporting equality.
		return false
	}
	if len(firstMap) != len(secondMap) {
		return false
	}
	for key, value := range firstMap {
		other, ok := secondMap[key]
		if !ok || other != value {
			return false
		}
	}
	return true
}

// DSSASN1UtilsAttributeMap returns a map of the principal's X.500 attribute types and
// values. Port of get(X500Principal); the Java method name has no meaning of its own, so
// the Go name says what it returns.
//
// UPSTREAM QUIRK, reproduced: both the key and the value go through
// IETFUtils#valueToString, and an OBJECT IDENTIFIER is not an ASN1String, so the keys come
// out as "#" followed by the hex of the attribute type's DER encoding rather than as a
// dotted OID. The map is only ever compared against another map built the same way, so the
// quirk is invisible to callers - but it must be kept for the comparison to stay faithful.
//
// Java throws a DSSException on a malformed name; the Go port returns nil instead, which
// makes DSSASN1UtilsX500PrincipalAreEquals answer false for two different malformed names.
func DSSASN1UtilsAttributeMap(x500Principal *model.X500Principal) map[string]string {
	treeMap := make(map[string]string)
	element, _, err := dssASN1UtilsParse(x500Principal.Encoded())
	if err != nil || !element.constructed {
		return nil
	}
	for _, set := range element.children {
		if !set.constructed {
			return nil
		}
		for _, sequence := range set.children {
			if !sequence.constructed || len(sequence.children) != 2 {
				// Java: "The DLSequence must contains exactly 2 elements."
				return nil
			}
			stringAttributeType := DSSASN1UtilsString(sequence.children[0].encoded)
			stringAttributeValue := DSSASN1UtilsString(sequence.children[1].encoded)
			treeMap[stringAttributeType] = stringAttributeValue
		}
	}
	return treeMap
}

// DSSASN1UtilsString converts the DER encoding of an ASN.1 value to a string, preserving the
// object class and structure: a string type yields its text, anything else yields "#"
// followed by the hex of its DER encoding. Whitespace is trimmed, per RFC 4518 "2.6.1.
// Insignificant Space Handling". Port of getString(ASN1Encodable), i.e. of BouncyCastle's
// IETFUtils#valueToString followed by String#trim.
//
// Returns the empty string where Java returns null (a nil argument, or a value that cannot
// be handled).
func DSSASN1UtilsString(attributeValue []byte) string {
	if attributeValue == nil {
		// Upstream logs "Null attribute has been provided!".
		return ""
	}
	element, rest, err := dssASN1UtilsParse(attributeValue)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to handle attribute : {}".
		return ""
	}
	return dssASN1UtilsJavaTrim(dssASN1UtilsValueToString(element))
}

// DSSASN1UtilsExtractAttributeFromX500Principal returns the value of the first attribute
// with the given type found in the principal, or the empty string (Java: null) when the
// principal carries none.
// Port of extractAttributeFromX500Principal(ASN1ObjectIdentifier, X500PrincipalHelper).
func DSSASN1UtilsExtractAttributeFromX500Principal(identifier asn1.ObjectIdentifier, principal *model.X500PrincipalHelper) string {
	element, _, err := dssASN1UtilsParse(principal.Encoded())
	if err != nil || !element.constructed {
		return ""
	}
	for _, set := range element.children {
		if !set.constructed {
			continue
		}
		for _, sequence := range set.children {
			if !sequence.constructed || len(sequence.children) != 2 {
				continue
			}
			oid, err := sequence.children[0].objectIdentifier()
			if err != nil || !identifier.Equal(oid) {
				continue
			}
			return dssASN1UtilsASN1ToString(sequence.children[1])
		}
	}
	return ""
}

// DSSASN1UtilsSubjectCommonName extracts the subject common name from the certificate token.
// Port of getSubjectCommonName(CertificateToken).
func DSSASN1UtilsSubjectCommonName(cert *model.CertificateToken) string {
	return DSSASN1UtilsExtractAttributeFromX500Principal(dssASN1UtilsOIDCN, cert.Subject())
}

// DSSASN1UtilsHumanReadableName extracts the pretty-printed name of the certificate token.
// Port of getHumanReadableName(CertificateToken).
func DSSASN1UtilsHumanReadableName(cert *model.CertificateToken) string {
	return DSSASN1UtilsHumanReadableNameOfPrincipal(cert.Subject())
}

// DSSASN1UtilsHumanReadableNameOfPrincipal extracts the pretty-printed name from the
// principal: the first of CN, GIVENNAME, SURNAME, NAME, PSEUDONYM, O, OU that is present.
// Port of getHumanReadableName(X500PrincipalHelper).
func DSSASN1UtilsHumanReadableNameOfPrincipal(x500PrincipalHelper *model.X500PrincipalHelper) string {
	return dssASN1UtilsFirstNotNull(x500PrincipalHelper, dssASN1UtilsOIDCN, dssASN1UtilsOIDGivenName,
		dssASN1UtilsOIDSurname, dssASN1UtilsOIDName, dssASN1UtilsOIDPseudonym, dssASN1UtilsOIDO, dssASN1UtilsOIDOU)
}

// dssASN1UtilsFirstNotNull ports the private firstNotNull(X500PrincipalHelper, oids...).
func dssASN1UtilsFirstNotNull(x500PrincipalHelper *model.X500PrincipalHelper, oids ...asn1.ObjectIdentifier) string {
	for _, oid := range oids {
		value := DSSASN1UtilsExtractAttributeFromX500Principal(oid, x500PrincipalHelper)
		if value != "" {
			return value
		}
	}
	return ""
}

// TODO(phase-3): getFirstSignerInformation(SignerInformationStore) belongs to the CMS layer.

// DSSASN1UtilsIsASN1SequenceTag reports whether the byte is the identifier octet of an ASN.1
// SEQUENCE. Port of isASN1SequenceTag(byte).
func DSSASN1UtilsIsASN1SequenceTag(tagByte byte) bool {
	// BERTags.SEQUENCE | BERTags.CONSTRUCTED = 0x30
	return (dssASN1UtilsTagSequence | dssASN1UtilsConstructed) == tagByte
}

// DSSASN1UtilsDate reads a DER-encoded X.509 Time (a UTCTime or a GeneralizedTime) and
// returns it, or the zero time.Time when it cannot be read.
// Port of getDate(ASN1Encodable), which returns null on failure.
func DSSASN1UtilsDate(encodable []byte) time.Time {
	element, rest, err := dssASN1UtilsParse(encodable)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to retrieve the date {}".
		return time.Time{}
	}
	var date time.Time
	switch {
	case element.isUniversal(dssASN1UtilsTagUTCTime):
		date, err = dssASN1UtilsParseUTCTime(string(element.content))
	case element.isUniversal(dssASN1UtilsTagGeneralizedTime):
		date, err = dssASN1UtilsParseGeneralizedTime(string(element.content))
	default:
		return time.Time{}
	}
	if err != nil {
		return time.Time{}
	}
	return date
}

// TODO(phase-3): getTimeStampTokenGenerationTime(TimeStampToken),
// getRevocationValues(ASN1Encodable) and getCertificateRef(OtherCertID) read the RFC 3161
// and CAdES structures the CMS phase introduces.

// DSSASN1UtilsIsAsn1Encoded reports whether the binaries are ASN.1 encoded.
// Port of isAsn1Encoded(byte[]).
func DSSASN1UtilsIsAsn1Encoded(binaries []byte) bool {
	if len(binaries) == 0 {
		return false
	}
	_, _, err := dssASN1UtilsParse(binaries)
	return err == nil
}

// DSSASN1UtilsIsAsn1EncodedSignatureValue reports whether the signature-value binaries are
// an ASN.1 SEQUENCE of two elements. Port of isAsn1EncodedSignatureValue(byte[]).
func DSSASN1UtilsIsAsn1EncodedSignatureValue(binaries []byte) bool {
	element, _, err := dssASN1UtilsParse(binaries)
	if err != nil || !element.isUniversal(dssASN1UtilsTagSequence) || !element.constructed {
		return false
	}
	return len(element.children) == 2
}

// DSSASN1UtilsEnsurePlainSignatureValue converts an ASN.1 signature value to the
// concatenated (plain) R || S format when the encryption algorithm calls for it.
// NOTE: used in XAdES and JAdES. Port of ensurePlainSignatureValue(EncryptionAlgorithm, byte[]).
func DSSASN1UtilsEnsurePlainSignatureValue(algorithm enumerations.EncryptionAlgorithm, signatureValue []byte) ([]byte, error) {
	if (enumerations.EncryptionAlgorithm_ECDSA == algorithm || enumerations.EncryptionAlgorithm_PLAIN_ECDSA == algorithm ||
		enumerations.EncryptionAlgorithm_DSA == algorithm) && DSSASN1UtilsIsAsn1EncodedSignatureValue(signatureValue) {
		return DSSASN1UtilsToPlainDSASignatureValue(signatureValue)
	}
	return signatureValue, nil
}

// DSSASN1UtilsToPlainDSASignatureValue converts an ASN.1 (r, s) pair into the concatenation
// of R and S used by ECDSA/DSA in XML and JOSE signatures.
// Port of toPlainDSASignatureValue(byte[]), i.e. of BouncyCastle's PlainDSAEncoding#encode.
//
// See http://www.w3.org/TR/xmldsig-core/#dsa-sha1 and RFC 4050 "3.3. ECDSA Signatures".
func DSSASN1UtilsToPlainDSASignatureValue(asn1SignatureValue []byte) ([]byte, error) {
	order, err := DSSASN1UtilsOrderFromSignatureValue(asn1SignatureValue)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to convert to plain", err)
	}
	element, _, err := dssASN1UtilsParse(asn1SignatureValue)
	if err != nil || !element.constructed || len(element.children) != 2 {
		return nil, model.NewDSSError("Unable to convert to plain : the signature value is not a SEQUENCE of 2 elements")
	}
	if !element.children[0].isUniversal(dssASN1UtilsTagInteger) || !element.children[1].isUniversal(dssASN1UtilsTagInteger) {
		return nil, model.NewDSSError("Unable to convert to plain : the signature value components are not INTEGERs")
	}
	r := element.children[0].integer()
	s := element.children[1].integer()

	// PlainDSAEncoding#encode: each value is left-padded to the unsigned byte length of the
	// order, and a value outside [0, order) is rejected.
	valueLength := dssASN1UtilsUnsignedByteLength(order)
	result := make([]byte, valueLength*2)
	if err := dssASN1UtilsEncodeDSAValue(order, r, result[:valueLength]); err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to convert to plain", err)
	}
	if err := dssASN1UtilsEncodeDSAValue(order, s, result[valueLength:]); err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to convert to plain", err)
	}
	return result, nil
}

// dssASN1UtilsEncodeDSAValue ports PlainDSAEncoding#encodeValue.
func dssASN1UtilsEncodeDSAValue(order, value *big.Int, buffer []byte) error {
	if value.Sign() < 0 || value.Cmp(order) >= 0 {
		return errors.New("Value out of range")
	}
	valueBytes := value.Bytes()
	if len(valueBytes) > len(buffer) {
		return errors.New("Value out of range")
	}
	copy(buffer[len(buffer)-len(valueBytes):], valueBytes)
	return nil
}

// DSSASN1UtilsToStandardDSASignatureValue converts a plain signature value into its ASN.1
// SEQUENCE { r INTEGER, s INTEGER } form.
// Port of toStandardDSASignatureValue(byte[]).
func DSSASN1UtilsToStandardDSASignatureValue(signatureValue []byte) ([]byte, error) {
	if len(signatureValue) == 0 {
		return nil, model.NewDSSError("Unable to convert to standard DSA : the signature value is empty")
	}
	signatureValuePartLength := len(signatureValue) / 2
	r := new(big.Int).SetBytes(signatureValue[0:signatureValuePartLength])
	s := new(big.Int).SetBytes(signatureValue[signatureValuePartLength : signatureValuePartLength*2])
	body := append(dssASN1UtilsEncodeInteger(r), dssASN1UtilsEncodeInteger(s)...)
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagSequence|dssASN1UtilsConstructed, body), nil
}

// DSSASN1UtilsOrderFromSignatureValue returns the order parameter corresponding to the given
// signature value, i.e. max(r, s) + 1. Port of getOrderFromSignatureValue(byte[]).
func DSSASN1UtilsOrderFromSignatureValue(signatureValue []byte) (*big.Int, error) {
	var rValue, sValue *big.Int
	if DSSASN1UtilsIsAsn1EncodedSignatureValue(signatureValue) {
		element, _, err := dssASN1UtilsParse(signatureValue)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to extract order from a signature value", err)
		}
		if len(element.children) != 2 {
			return nil, model.NewDSSError("Unable to extract order from a signature value : ASN1 Sequence size should be 2!")
		}
		if !element.children[0].isUniversal(dssASN1UtilsTagInteger) || !element.children[1].isUniversal(dssASN1UtilsTagInteger) {
			return nil, model.NewDSSError("Unable to extract order from a signature value : the components are not INTEGERs")
		}
		rValue = element.children[0].integer()
		sValue = element.children[1].integer()
	} else {
		if len(signatureValue)%2 != 0 {
			return nil, model.NewDSSError("Unable to extract order from a signature value : signatureValue binaries length shall be dividable by 2!")
		}
		valueLength := len(signatureValue) / 2
		rValue = new(big.Int).SetBytes(signatureValue[0:valueLength])
		sValue = new(big.Int).SetBytes(signatureValue[valueLength : valueLength*2])
	}
	max := rValue
	if sValue.Cmp(rValue) > 0 {
		max = sValue
	}
	return new(big.Int).Add(max, big.NewInt(1)), nil
}

// DSSASN1UtilsSignatureValueBitLength returns the bit length of the provided signature
// value. Port of getSignatureValueBitLength(byte[]).
func DSSASN1UtilsSignatureValueBitLength(signatureValue []byte) (int, error) {
	order, err := DSSASN1UtilsOrderFromSignatureValue(signatureValue)
	if err != nil {
		return 0, model.NewDSSErrorMessageCause("Unable to extract a signature value bit length", err)
	}
	return dssASN1UtilsUnsignedByteLength(order) * 8, nil // convert to bits
}

// dssASN1UtilsUnsignedByteLength ports org.bouncycastle.util.BigIntegers#getUnsignedByteLength.
func dssASN1UtilsUnsignedByteLength(value *big.Int) int {
	return (value.BitLen() + 7) / 8
}

// DSSASN1UtilsDirectoryStringValue returns the value of an ASN.1 DirectoryString, or the
// empty string (Java: null) when the value is not one.
// Port of getDirectoryStringValue(ASN1Encodable).
func DSSASN1UtilsDirectoryStringValue(directoryStringInstance []byte) string {
	element, rest, err := dssASN1UtilsParse(directoryStringInstance)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to build a DirectoryString instance. Reason : {}".
		return ""
	}
	switch {
	case element.isUniversal(dssASN1UtilsTagPrintableString), element.isUniversal(dssASN1UtilsTagT61String),
		element.isUniversal(dssASN1UtilsTagUTF8String), element.isUniversal(dssASN1UtilsTagUniversalString),
		element.isUniversal(dssASN1UtilsTagBMPString):
		return element.asString()
	}
	return ""
}

// DSSASN1UtilsBuildSPDocSpecificationID builds the SPDocSpecification attribute value from
// the given OID or URI.
//
//	SPDocSpecification ::= CHOICE {
//	    oid OBJECT IDENTIFIER,
//	    uri IA5String }
//
// Port of buildSPDocSpecificationId(String).
func DSSASN1UtilsBuildSPDocSpecificationID(oidOrURI string) ([]byte, error) {
	if dssASN1UtilsIsOidCode(oidOrURI) {
		oid, err := dssASN1UtilsObjectIdentifier(oidOrURI)
		if err != nil {
			return nil, err
		}
		encoded := dssASN1UtilsEncodeOID(oid)
		if encoded == nil {
			// Java: new ASN1ObjectIdentifier("0") raises
			// IllegalArgumentException("string 0 not a valid OID"). DSSUtils#isOidCode
			// accepts a single arc, ASN.1 does not.
			return nil, fmt.Errorf("string %s not a valid OID", oidOrURI)
		}
		return encoded, nil
	}
	return dssASN1UtilsWriteTLV(dssASN1UtilsTagIA5String, []byte(oidOrURI)), nil
}

// dssASN1UtilsOidCodePattern is DSSUtils#isOidCode's regular expression, verbatim.
var dssASN1UtilsOidCodePattern = regexp.MustCompile(`^([0-2])((\.0)|(\.[1-9][0-9]*))*$`)

// dssASN1UtilsIsOidCode reproduces DSSUtils#isOidCode, which decides whether a string is to
// be read as an OBJECT IDENTIFIER rather than as a URI.
func dssASN1UtilsIsOidCode(id string) bool {
	return id != "" && dssASN1UtilsOidCodePattern.MatchString(id)
}

// dssASN1UtilsDigest computes the digest of the concatenated inputs, standing in for
// DSSUtils#digest(DigestAlgorithm, byte[]...).
func dssASN1UtilsDigest(digestAlgorithm enumerations.DigestAlgorithm, data ...[]byte) ([]byte, error) {
	calculator, err := NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		return nil, err
	}
	for _, chunk := range data {
		calculator.Update(chunk)
	}
	return calculator.MessageDigest(digestAlgorithm).Value(), nil
}

// -----------------------------------------------------------------------------
// String conversion (BouncyCastle IETFUtils / ASN1String).
// -----------------------------------------------------------------------------

// dssASN1UtilsValueToString ports org.bouncycastle.asn1.x500.style.IETFUtils#valueToString.
func dssASN1UtilsValueToString(element *dssASN1UtilsElement) string {
	var buffer []rune
	// IETFUtils#valueToString takes the string branch for every ASN1String EXCEPT
	// ASN1UniversalString, which it explicitly excludes so that it is hash-encoded like a
	// non-string value.
	if element.isASN1String() && !element.isUniversal(dssASN1UtilsTagUniversalString) {
		value := element.asString()
		if len(value) > 0 && value[0] == '#' {
			buffer = append(buffer, '\\')
		}
		buffer = append(buffer, []rune(value)...)
	} else {
		buffer = append(buffer, '#')
		buffer = append(buffer, []rune(dssASN1UtilsHexLower(element.derEncoded()))...)
	}

	index := 0
	if len(buffer) >= 2 && buffer[0] == '\\' && buffer[1] == '#' {
		index += 2
	}
	for ; index < len(buffer); index++ {
		switch buffer[index] {
		case ',', '"', '\\', '+', '=', '<', '>', ';':
			buffer = append(buffer[:index], append([]rune{'\\'}, buffer[index:]...)...)
			index++
		}
	}

	start := 0
	for start < len(buffer) && buffer[start] == ' ' {
		buffer = append(buffer[:start], append([]rune{'\\'}, buffer[start:]...)...)
		start += 2
	}
	end := len(buffer) - 1
	for end >= 0 && buffer[end] == ' ' {
		buffer = append(buffer[:end], append([]rune{'\\'}, buffer[end:]...)...)
		end--
	}
	return string(buffer)
}

// dssASN1UtilsASN1ToString reproduces ASN1Primitive#toString for the value types an X.500
// attribute can carry: a string type yields its text (a BIT STRING and a UniversalString
// their "#"+UPPER-case-hex form, see asString), an OBJECT IDENTIFIER its dotted form,
// anything else "#" followed by the hex of its DER encoding.
//
// DEVIATION: BouncyCastle renders a constructed value (an attribute whose value is a SEQUENCE
// or a SET, which no X.520 attribute type defines) as its ASN1Dump-style "[a, b]" listing;
// this port hash-encodes it like any other non-string value.
func dssASN1UtilsASN1ToString(element *dssASN1UtilsElement) string {
	if element.isASN1String() {
		return element.asString()
	}
	switch {
	case element.isUniversal(dssASN1UtilsTagOID):
		if oid, err := element.objectIdentifier(); err == nil {
			return oid.String()
		}
	case element.isUniversal(dssASN1UtilsTagInteger):
		return element.integer().String()
	case element.isUniversal(dssASN1UtilsTagOctetString):
		// ASN1OctetString#toString hexes the CONTENT, not the whole encoding.
		return "#" + dssASN1UtilsHexLower(element.octets())
	case element.isUniversal(dssASN1UtilsTagBoolean):
		if len(element.content) > 0 && element.content[0] != 0x00 {
			return "TRUE"
		}
		return "FALSE"
	case element.isUniversal(dssASN1UtilsTagNull):
		return "NULL"
	}
	return "#" + dssASN1UtilsHexLower(element.derEncoded())
}

// dssASN1UtilsJavaTrim ports java.lang.String#trim, which strips every leading and trailing
// character whose code point is not greater than U+0020.
func dssASN1UtilsJavaTrim(value string) string {
	return strings.TrimFunc(value, func(r rune) bool { return r <= ' ' })
}

// dssASN1UtilsHexUpper renders bytes as upper-case hex without a separator, the way the
// private hex table of ASN1BitString#getString / ASN1UniversalString#getString does.
func dssASN1UtilsHexUpper(data []byte) string {
	const digits = "0123456789ABCDEF"
	var builder strings.Builder
	builder.Grow(2 * len(data))
	for _, b := range data {
		builder.WriteByte(digits[b>>4])
		builder.WriteByte(digits[b&0x0F])
	}
	return builder.String()
}

// dssASN1UtilsHexLower renders bytes as lower-case hex without a separator, the way
// org.bouncycastle.util.encoders.Hex#encode does.
func dssASN1UtilsHexLower(data []byte) string {
	const digits = "0123456789abcdef"
	var builder strings.Builder
	builder.Grow(2 * len(data))
	for _, b := range data {
		builder.WriteByte(digits[b>>4])
		builder.WriteByte(digits[b&0x0F])
	}
	return builder.String()
}

// -----------------------------------------------------------------------------
// Time parsing (BouncyCastle ASN1UTCTime / ASN1GeneralizedTime).
// -----------------------------------------------------------------------------

// dssASN1UtilsParseUTCTime ports ASN1UTCTime#getDate, including getAdjustedTime's two-digit
// year rule: a year of 50 or more is 19xx, below 50 it is 20xx.
func dssASN1UtilsParseUTCTime(value string) (time.Time, error) {
	if len(value) < 10 {
		return time.Time{}, fmt.Errorf("invalid UTCTime: %q", value)
	}
	year, err := strconv.Atoi(value[0:2])
	if err != nil {
		return time.Time{}, err
	}
	century := "20"
	if year >= 50 {
		century = "19"
	}
	return dssASN1UtilsParseTimeDigits(century+value[0:2], value[2:])
}

// dssASN1UtilsParseGeneralizedTime ports ASN1GeneralizedTime#getDate for the
// YYYYMMDDHHMM[SS[.f+]] forms, with an optional Z or +/-HH[MM] offset.
func dssASN1UtilsParseGeneralizedTime(value string) (time.Time, error) {
	if len(value) < 10 {
		return time.Time{}, fmt.Errorf("invalid GeneralizedTime: %q", value)
	}
	return dssASN1UtilsParseTimeDigits(value[0:4], value[4:])
}

// dssASN1UtilsParseTimeDigits decodes MMDDHHMM[SS[.f+]][zone] against an already resolved
// four-digit year.
//
// DEVIATION: a time without a zone is read as UTC. BouncyCastle interprets it in the JVM's
// default time zone, which makes the parsed instant depend on the host configuration;
// certificates and revocation data always carry the "Z" suffix, so the two agree in practice.
func dssASN1UtilsParseTimeDigits(year string, rest string) (time.Time, error) {
	offset := time.Duration(0)
	location := time.UTC
	switch {
	case strings.HasSuffix(rest, "Z"):
		rest = rest[:len(rest)-1]
	default:
		if index := strings.LastIndexAny(rest, "+-"); index > 0 {
			zone := rest[index:]
			rest = rest[:index]
			sign := time.Duration(1)
			if zone[0] == '-' {
				sign = -1
			}
			zone = zone[1:]
			if len(zone) != 2 && len(zone) != 4 && !(len(zone) == 5 && zone[2] == ':') {
				return time.Time{}, fmt.Errorf("invalid time zone: %q", zone)
			}
			hours, err := strconv.Atoi(zone[0:2])
			if err != nil {
				return time.Time{}, err
			}
			minutes := 0
			if len(zone) >= 4 {
				minutes, err = strconv.Atoi(zone[len(zone)-2:])
				if err != nil {
					return time.Time{}, err
				}
			}
			offset = sign * (time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
		}
	}

	fraction := time.Duration(0)
	if index := strings.IndexAny(rest, ".,"); index >= 0 {
		digits := rest[index+1:]
		rest = rest[:index]
		// Java keeps millisecond precision; a longer fraction is truncated.
		for len(digits) < 3 {
			digits += "0"
		}
		milliseconds, err := strconv.Atoi(digits[0:3])
		if err != nil {
			return time.Time{}, err
		}
		fraction = time.Duration(milliseconds) * time.Millisecond
	}

	if len(rest) != 6 && len(rest) != 8 && len(rest) != 10 {
		return time.Time{}, fmt.Errorf("invalid time digits: %q", year+rest)
	}
	layoutSource := year + rest
	layout := "2006010215"
	switch len(rest) {
	case 8:
		layout = "200601021504"
	case 10:
		layout = "20060102150405"
	}
	parsed, err := time.ParseInLocation(layout, layoutSource, location)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.Add(fraction).Add(-offset), nil
}

// -----------------------------------------------------------------------------
// The ASN.1 element model: BER parsing, DER/DL/BER writing.
// -----------------------------------------------------------------------------

// dssASN1UtilsElement is one parsed ASN.1 element. It replaces BouncyCastle's ASN1Primitive
// hierarchy with exactly what the ported code needs: the identifier, the content and, for a
// constructed element, its children.
type dssASN1UtilsElement struct {
	// class is the identifier's class bits: 0x00 universal, 0x40 application,
	// 0x80 context-specific, 0xC0 private.
	class byte
	// constructed reports whether the element is constructed.
	constructed bool
	// tagNumber is the tag number, decoded from the high-tag-number form when needed.
	tagNumber uint64
	// indefinite reports whether the element used the indefinite-length form.
	indefinite bool
	// content holds the content octets of a primitive element.
	content []byte
	// children holds the components of a constructed element.
	children []*dssASN1UtilsElement
	// encoded is the element's original encoding, tag and length included.
	encoded []byte
}

// isUniversal reports whether the element carries the given universal tag number.
func (e *dssASN1UtilsElement) isUniversal(tagNumber uint64) bool {
	return e.class == 0 && e.tagNumber == tagNumber
}

// isASN1String reports whether the element is one of the types implementing BouncyCastle's
// ASN1String interface.
//
// ASN1BitString and ASN1UniversalString implement it too - their getString() answers "#"
// followed by the upper-case hex of the whole encoding, see asString - while
// ASN1ObjectDescriptor does NOT (it wraps an ASN1GraphicString without implementing the
// interface), so tag 7 falls through to the hash-encoded branch of the callers.
func (e *dssASN1UtilsElement) isASN1String() bool {
	if e.class != 0 {
		return false
	}
	switch e.tagNumber {
	case dssASN1UtilsTagBitString, dssASN1UtilsTagUTF8String, dssASN1UtilsTagNumericString,
		dssASN1UtilsTagPrintableString, dssASN1UtilsTagT61String, dssASN1UtilsTagVideotexString,
		dssASN1UtilsTagIA5String, dssASN1UtilsTagGraphicString, dssASN1UtilsTagVisibleString,
		dssASN1UtilsTagGeneralString, dssASN1UtilsTagUniversalString, dssASN1UtilsTagBMPString:
		return true
	}
	return false
}

// asString decodes the element's content the way the matching BouncyCastle ASN1String does:
// UTF8String as UTF-8, BMPString as UTF-16BE, BIT STRING and UniversalString as "#" followed
// by the UPPER-case hex of the complete encoding (ASN1BitString#getString and
// ASN1UniversalString#getString are both spelled that way, and neither decodes its content as
// text), and every other string type byte-per-character (ISO-8859-1), which is what
// org.bouncycastle.util.Strings#fromByteArray produces.
func (e *dssASN1UtilsElement) asString() string {
	switch e.tagNumber {
	case dssASN1UtilsTagUTF8String:
		return string(e.content)
	case dssASN1UtilsTagBitString, dssASN1UtilsTagUniversalString:
		return "#" + dssASN1UtilsHexUpper(e.derEncoded())
	case dssASN1UtilsTagBMPString:
		if len(e.content)%2 != 0 {
			return ""
		}
		units := make([]uint16, len(e.content)/2)
		for index := range units {
			units[index] = uint16(e.content[2*index])<<8 | uint16(e.content[2*index+1])
		}
		return string(utf16.Decode(units))
	}
	runes := make([]rune, len(e.content))
	for index, b := range e.content {
		runes[index] = rune(b)
	}
	return string(runes)
}

// octets returns the content of an OCTET STRING, joining the segments of a constructed one.
func (e *dssASN1UtilsElement) octets() []byte {
	if !e.constructed {
		return e.content
	}
	var joined []byte
	for _, child := range e.children {
		joined = append(joined, child.octets()...)
	}
	return joined
}

// bitStringOctets returns the value bits of a BIT STRING, i.e. its content without the
// leading count of unused bits, which is what BouncyCastle's ASN1BitString#getOctets returns.
func (e *dssASN1UtilsElement) bitStringOctets() []byte {
	content := e.derContent()
	if len(content) == 0 {
		return nil
	}
	return content[1:]
}

// integer returns the value of an INTEGER, decoded as a two's-complement big-endian number.
func (e *dssASN1UtilsElement) integer() *big.Int {
	value := new(big.Int)
	if len(e.content) == 0 {
		return value
	}
	if e.content[0]&0x80 != 0 {
		// Negative: subtract 2^(8*len) from the unsigned interpretation.
		value.SetBytes(e.content)
		value.Sub(value, new(big.Int).Lsh(big.NewInt(1), uint(8*len(e.content))))
		return value
	}
	return value.SetBytes(e.content)
}

// objectIdentifier decodes an OBJECT IDENTIFIER.
func (e *dssASN1UtilsElement) objectIdentifier() (asn1.ObjectIdentifier, error) {
	if !e.isUniversal(dssASN1UtilsTagOID) {
		return nil, errors.New("not an OBJECT IDENTIFIER")
	}
	var oid asn1.ObjectIdentifier
	if _, err := asn1.Unmarshal(e.derEncoded(), &oid); err != nil {
		return nil, err
	}
	return oid, nil
}

// identifierOctets rebuilds the identifier octets of the element.
func (e *dssASN1UtilsElement) identifierOctets(constructed bool) []byte {
	first := e.class
	if constructed {
		first |= dssASN1UtilsConstructed
	}
	if e.tagNumber < 0x1F {
		return []byte{first | byte(e.tagNumber)}
	}
	first |= dssASN1UtilsTagMask
	var trailer []byte
	value := e.tagNumber
	trailer = append(trailer, byte(value&0x7F))
	for value >>= 7; value > 0; value >>= 7 {
		trailer = append([]byte{byte(value&0x7F) | 0x80}, trailer...)
	}
	return append([]byte{first}, trailer...)
}

// derContent returns the element's content in DER form, collapsing a constructed OCTET or
// BIT STRING into its primitive content.
func (e *dssASN1UtilsElement) derContent() []byte {
	if !e.constructed {
		return e.content
	}
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagOctetString {
		return e.octets()
	}
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagBitString {
		// A segmented BIT STRING collapses to the concatenation of the segments' value
		// bits, prefixed by the unused-bit count of the final segment (only the final
		// segment may declare unused bits).
		var bits []byte
		unused := byte(0)
		for index, child := range e.children {
			segment := child.derContent()
			if len(segment) == 0 {
				continue
			}
			bits = append(bits, segment[1:]...)
			if index == len(e.children)-1 {
				unused = segment[0]
			}
		}
		return append([]byte{unused}, bits...)
	}
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagSet {
		return dssASN1UtilsSortSet(e.children)
	}
	var body []byte
	for _, child := range e.children {
		body = append(body, child.derEncoded()...)
	}
	return body
}

// derEncoded returns the DER encoding of the element.
func (e *dssASN1UtilsElement) derEncoded() []byte {
	if !e.constructed {
		content := e.content
		if e.class == 0 && e.tagNumber == dssASN1UtilsTagBoolean && len(content) == 1 && content[0] != 0x00 {
			// DER requires TRUE to be all-ones.
			content = []byte{0xFF}
		}
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), content)
	}
	if e.class == 0 && (e.tagNumber == dssASN1UtilsTagOctetString || e.tagNumber == dssASN1UtilsTagBitString) {
		// A constructed string is primitive in DER.
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), e.derContent())
	}
	var body []byte
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagSet {
		body = dssASN1UtilsSortSet(e.children)
	} else {
		for _, child := range e.children {
			body = append(body, child.derEncoded()...)
		}
	}
	return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(true), body)
}

// dlEncoded returns the DL encoding of the element: definite lengths, and the component
// order of the input preserved (a SET is NOT sorted). A constructed OCTET or BIT STRING
// still collapses to its primitive form, as BouncyCastle's DL encoding of a BEROctetString
// or a BERBitString does.
func (e *dssASN1UtilsElement) dlEncoded() []byte {
	if !e.constructed {
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), e.content)
	}
	if e.class == 0 && (e.tagNumber == dssASN1UtilsTagOctetString || e.tagNumber == dssASN1UtilsTagBitString) {
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), e.derContent())
	}
	var body []byte
	for _, child := range e.children {
		body = append(body, child.dlEncoded()...)
	}
	return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(true), body)
}

// berEncoded returns the BER encoding of the element, reproducing what BouncyCastle emits
// for an object read back from the same input:
//
//   - a constructed BIT STRING collapses to its primitive form;
//   - a constructed OCTET STRING always becomes an indefinite-length constructed string,
//     holding one segment with the joined content when the input was indefinite (the
//     BEROctetStringParser path) and the input's segments otherwise;
//   - every other constructed element keeps the length form the input used, so a
//     definite-length input encodes exactly as its DL form.
func (e *dssASN1UtilsElement) berEncoded() []byte {
	if !e.constructed {
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), e.content)
	}
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagBitString {
		return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(false), e.derContent())
	}
	var body []byte
	if e.class == 0 && e.tagNumber == dssASN1UtilsTagOctetString {
		if e.indefinite {
			body = dssASN1UtilsWriteTLV(dssASN1UtilsTagOctetString, e.octets())
		} else {
			for _, child := range e.children {
				body = append(body, child.berEncoded()...)
			}
		}
		return dssASN1UtilsWriteIndefiniteTLV(e.identifierOctets(true), body)
	}
	for _, child := range e.children {
		body = append(body, child.berEncoded()...)
	}
	if e.indefinite {
		return dssASN1UtilsWriteIndefiniteTLV(e.identifierOctets(true), body)
	}
	return dssASN1UtilsWriteIdentifiedTLV(e.identifierOctets(true), body)
}

// dssASN1UtilsSortSet returns the DER encodings of a SET's components in ascending order, as
// DER requires and as BouncyCastle's ASN1Set#sortElements implements.
func dssASN1UtilsSortSet(children []*dssASN1UtilsElement) []byte {
	encodings := make([][]byte, len(children))
	for index, child := range children {
		encodings[index] = child.derEncoded()
	}
	sort.SliceStable(encodings, func(a, b int) bool {
		return bytes.Compare(encodings[a], encodings[b]) < 0
	})
	var body []byte
	for _, encoding := range encodings {
		body = append(body, encoding...)
	}
	return body
}

// dssASN1UtilsWriteTLV builds an element from a single-byte identifier and its content.
func dssASN1UtilsWriteTLV(identifier byte, content []byte) []byte {
	return dssASN1UtilsWriteIdentifiedTLV([]byte{identifier}, content)
}

// dssASN1UtilsWriteIdentifiedTLV builds an element from its identifier octets and content,
// using the shortest definite-length encoding.
func dssASN1UtilsWriteIdentifiedTLV(identifier []byte, content []byte) []byte {
	out := make([]byte, 0, len(identifier)+4+len(content))
	out = append(out, identifier...)
	length := len(content)
	if length < 0x80 {
		out = append(out, byte(length))
	} else {
		var lengthBytes []byte
		for value := length; value > 0; value >>= 8 {
			lengthBytes = append([]byte{byte(value & 0xFF)}, lengthBytes...)
		}
		out = append(out, byte(0x80|len(lengthBytes)))
		out = append(out, lengthBytes...)
	}
	return append(out, content...)
}

// dssASN1UtilsWriteIndefiniteTLV builds a constructed element using the indefinite-length
// form, i.e. the 0x80 length octet followed by the content and the end-of-contents octets.
func dssASN1UtilsWriteIndefiniteTLV(identifier []byte, content []byte) []byte {
	out := make([]byte, 0, len(identifier)+3+len(content))
	out = append(out, identifier...)
	out = append(out, 0x80)
	out = append(out, content...)
	return append(out, 0x00, 0x00)
}

// dssASN1UtilsEncodeOID returns the DER encoding of an OBJECT IDENTIFIER.
func dssASN1UtilsEncodeOID(oid asn1.ObjectIdentifier) []byte {
	encoded, err := asn1.Marshal(oid)
	if err != nil {
		return nil
	}
	return encoded
}

// dssASN1UtilsObjectIdentifier parses a dotted-decimal OID string.
func dssASN1UtilsObjectIdentifier(value string) (asn1.ObjectIdentifier, error) {
	if value == "" {
		return nil, errors.New("the OID is not defined")
	}
	parts := strings.Split(value, ".")
	oid := make(asn1.ObjectIdentifier, 0, len(parts))
	for _, part := range parts {
		component, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("string %s not an OID", value)
		}
		oid = append(oid, component)
	}
	return oid, nil
}

// dssASN1UtilsEncodeInteger returns the DER encoding of an INTEGER.
func dssASN1UtilsEncodeInteger(value *big.Int) []byte {
	encoded, err := asn1.Marshal(value)
	if err != nil {
		return nil
	}
	return encoded
}

// dssASN1UtilsParse reads one ASN.1 element from the front of the input, supporting both the
// definite and the indefinite length forms, and returns it together with the remaining bytes.
func dssASN1UtilsParse(input []byte) (*dssASN1UtilsElement, []byte, error) {
	if len(input) < 2 {
		return nil, nil, errors.New("truncated ASN.1 element")
	}
	cursor := 0
	first := input[cursor]
	cursor++
	element := &dssASN1UtilsElement{
		class:       first & dssASN1UtilsClassMask,
		constructed: first&dssASN1UtilsConstructed != 0,
	}
	if first&dssASN1UtilsTagMask == dssASN1UtilsTagMask {
		var tagNumber uint64
		for {
			if cursor >= len(input) {
				return nil, nil, errors.New("truncated ASN.1 tag")
			}
			b := input[cursor]
			cursor++
			if tagNumber > (1<<56)-1 {
				return nil, nil, errors.New("ASN.1 tag number overflow")
			}
			tagNumber = tagNumber<<7 | uint64(b&0x7F)
			if b&0x80 == 0 {
				break
			}
		}
		element.tagNumber = tagNumber
	} else {
		element.tagNumber = uint64(first & dssASN1UtilsTagMask)
	}

	if cursor >= len(input) {
		return nil, nil, errors.New("truncated ASN.1 length")
	}
	lengthByte := input[cursor]
	cursor++

	if lengthByte == 0x80 {
		if !element.constructed {
			return nil, nil, errors.New("indefinite length on a primitive ASN.1 element")
		}
		element.indefinite = true
		rest := input[cursor:]
		for {
			if len(rest) >= 2 && rest[0] == 0x00 && rest[1] == 0x00 {
				rest = rest[2:]
				break
			}
			child, remaining, err := dssASN1UtilsParse(rest)
			if err != nil {
				return nil, nil, err
			}
			element.children = append(element.children, child)
			rest = remaining
			if len(rest) < 2 {
				return nil, nil, errors.New("missing end-of-contents octets")
			}
		}
		element.encoded = input[:len(input)-len(rest)]
		return element, rest, nil
	}

	length := 0
	if lengthByte&0x80 == 0 {
		length = int(lengthByte)
	} else {
		count := int(lengthByte & 0x7F)
		if count > 4 || cursor+count > len(input) {
			return nil, nil, errors.New("unsupported or truncated ASN.1 length")
		}
		for _, b := range input[cursor : cursor+count] {
			length = length<<8 | int(b)
		}
		cursor += count
		if length < 0 {
			return nil, nil, errors.New("invalid ASN.1 length")
		}
	}
	if cursor+length > len(input) {
		return nil, nil, errors.New("truncated ASN.1 element")
	}
	body := input[cursor : cursor+length]
	element.encoded = input[:cursor+length]
	if element.constructed {
		rest := body
		for len(rest) > 0 {
			child, remaining, err := dssASN1UtilsParse(rest)
			if err != nil {
				return nil, nil, err
			}
			element.children = append(element.children, child)
			rest = remaining
		}
	} else {
		element.content = body
	}
	return element, input[cursor+length:], nil
}

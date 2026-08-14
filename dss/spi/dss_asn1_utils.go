// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/DSSASN1Utils.java (DSS 6.5.RC1).
//
// This file replaces BouncyCastle's ASN.1 layer. Two conventions follow from that:
//
//   - Java passes around ASN1Encodable/ASN1Primitive objects; the Go port represents an
//     ASN.1 value by its encoding ([]byte), because that is what the surrounding code
//     digests, compares and serialises. Wherever identity or a digest depends on the
//     bytes, the original bytes are carried through untouched and never re-encoded.
//   - The BER/DER/DL re-encoding BouncyCastle performs in getEncoded(String) lives in
//     internal/asn1ber (see asn1ber.Element): DER output must be byte-exact, since it feeds
//     signature verification and archive-timestamp hashing.
//
// The BER/DER engine and the generic X.509 structures it produces (AlgorithmIdentifier,
// IssuerSerial, GeneralName) were extracted into github.com/utain/esig/dss/internal/asn1ber so
// that the CMS, XAdES and PAdES phases can share them without depending on dss-spi. The types
// stay reachable under their original names through the aliases below; only the DSSASN1Utils
// methods themselves are defined here.
//
// NOT PORTED (deferred to the CMS/timestamp phase, which owns the types they take):
// getEncoded(TimeStampToken), getEncoded(CMSSignedData), getDEREncoded(TimeStampToken),
// getDEREncoded(CMSSignedData), getTimeStampTokenGenerationTime, getX509CertificateHolder and
// getCertificate(X509CertificateHolder). See the TODO markers next to their Go neighbours.
//
// The CMS-typed methods the CAdES phase needed have landed: isEmpty(AttributeTable),
// emptyIfNull(AttributeTable), isAttributeOfType and getFirstSignerInformation are below;
// getAsn1Encodable(Attribute), getAsn1Attributes, toSignerIdentifier(SignerId),
// getCertificate(X509CertificateHolder) and getCertificateRef(OtherCertID) are in
// cms_certificate_source.go, and getRevocationValues in cms_crl_source.go.
package spi

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
)

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
//
// They live in internal/asn1ber; the aliases keep dss-spi's API unchanged.
// -----------------------------------------------------------------------------

// AlgorithmIdentifier is the X.509 structure
//
//	AlgorithmIdentifier ::= SEQUENCE {
//	    algorithm   OBJECT IDENTIFIER,
//	    parameters  ANY DEFINED BY algorithm OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.AlgorithmIdentifier.
type AlgorithmIdentifier = asn1ber.AlgorithmIdentifier

// GeneralName is one alternative of the X.509 GeneralName CHOICE, reduced to what
// IssuerSerial needs.
type GeneralName = asn1ber.GeneralName

// IssuerSerial is the X.509 attribute-certificate structure
//
//	IssuerSerial ::= SEQUENCE {
//	    issuer   GeneralNames,
//	    serial   CertificateSerialNumber,
//	    issuerUID UniqueIdentifier OPTIONAL }
//
// replacing org.bouncycastle.asn1.x509.IssuerSerial.
type IssuerSerial = asn1ber.IssuerSerial

// NewAlgorithmIdentifier builds an AlgorithmIdentifier without parameters.
func NewAlgorithmIdentifier(algorithm asn1.ObjectIdentifier) *AlgorithmIdentifier {
	return asn1ber.NewAlgorithmIdentifier(algorithm)
}

// NewAlgorithmIdentifierWithParameters builds an AlgorithmIdentifier carrying the given
// DER-encoded parameters.
func NewAlgorithmIdentifierWithParameters(algorithm asn1.ObjectIdentifier, parameters []byte) *AlgorithmIdentifier {
	return asn1ber.NewAlgorithmIdentifierWithParameters(algorithm, parameters)
}

// ParseAlgorithmIdentifier decodes an AlgorithmIdentifier from its DER encoding.
// Port of AlgorithmIdentifier.getInstance(Object).
func ParseAlgorithmIdentifier(der []byte) (*AlgorithmIdentifier, error) {
	return asn1ber.ParseAlgorithmIdentifier(der)
}

// -----------------------------------------------------------------------------
// The ported DSSASN1Utils methods.
// -----------------------------------------------------------------------------

// DSSASN1UtilsToASN1Primitive validates that the given bytes hold exactly one ASN.1 element
// and returns that element's encoding. Port of toASN1Primitive(byte[]); an ASN1Primitive is
// represented by its bytes in this port, so the method degenerates to a parse-and-return.
func DSSASN1UtilsToASN1Primitive(bytes []byte) ([]byte, error) {
	element, rest, err := asn1ber.Parse(bytes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot convert binaries to ASN1Primitive", err)
	}
	if len(rest) != 0 {
		return nil, model.NewDSSError("Cannot convert binaries to ASN1Primitive : extra data found after the object")
	}
	return element.Encoded(), nil
}

// DSSASN1UtilsIsDEROctetStringNull reports whether the given DER OCTET STRING encapsulates
// an ASN.1 NULL. Port of isDEROctetStringNull(DEROctetString), whose argument is here the
// complete encoding of the OCTET STRING rather than a parsed object.
func DSSASN1UtilsIsDEROctetStringNull(derOctetString []byte) bool {
	element, rest, err := asn1ber.Parse(derOctetString)
	if err != nil || len(rest) != 0 || !element.IsUniversal(asn1ber.TagOctetString) {
		return false
	}
	asn1Null, err := DSSASN1UtilsToASN1Primitive(element.Octets())
	if err != nil {
		return false
	}
	return bytes.Equal(asn1ber.DERNull, asn1Null)
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
	return element.DEREncoded(), nil
}

// DSSASN1UtilsDLEncoded returns the DL encoding of the given ASN.1 value: every length is
// definite, but the structure is otherwise preserved - constructed strings stay constructed
// and SET components keep their order. Port of getDLEncoded(ASN1Encodable)/getDLEncoded(byte[]).
func DSSASN1UtilsDLEncoded(asn1Encodable []byte) ([]byte, error) {
	element, err := dssASN1UtilsParseOne(asn1Encodable, "DL")
	if err != nil {
		return nil, err
	}
	return element.DLEncoded(), nil
}

// DSSASN1UtilsBEREncoded returns the BER encoding of the given ASN.1 value, preserving the
// indefinite lengths the input used. Port of getBEREncoded(ASN1Encodable).
func DSSASN1UtilsBEREncoded(asn1Encodable []byte) ([]byte, error) {
	element, err := dssASN1UtilsParseOne(asn1Encodable, "BER")
	if err != nil {
		return nil, err
	}
	return element.BEREncoded(), nil
}

// dssASN1UtilsParseOne parses exactly one element, reporting failures the way upstream's
// private getEncoded(ASN1Encodable, String) does.
func dssASN1UtilsParseOne(binaries []byte, encoding string) (*asn1ber.Element, error) {
	element, rest, err := asn1ber.Parse(binaries)
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
	element, rest, err := asn1ber.Parse(asn1Date)
	if err != nil || len(rest) != 0 || !element.IsUniversal(asn1ber.TagGeneralizedTime) {
		return time.Time{}, model.NewDSSError("Cannot parse Date : not an ASN1GeneralizedTime")
	}
	date, err := asn1ber.ParseGeneralizedTime(string(element.Content()))
	if err != nil {
		return time.Time{}, model.NewDSSErrorMessageCause("Cannot parse Date", err)
	}
	return date, nil
}

// DSSASN1UtilsToString reads the value of a DER OCTET STRING as a string.
// Port of toString(ASN1OctetString), i.e. new String(value.getOctets()).
func DSSASN1UtilsToString(value []byte) (string, error) {
	element, rest, err := asn1ber.Parse(value)
	if err != nil {
		return "", err
	}
	if len(rest) != 0 || !element.IsUniversal(asn1ber.TagOctetString) {
		return "", errors.New("not an ASN1OctetString")
	}
	return string(element.Octets()), nil
}

// DSSASN1UtilsAsn1SequenceFromDerOctetString returns the DER encoding of the SEQUENCE
// encapsulated in the given DER OCTET STRING.
// Port of getAsn1SequenceFromDerOctetString(byte[]).
func DSSASN1UtilsAsn1SequenceFromDerOctetString(binaries []byte) ([]byte, error) {
	content, err := dssASN1UtilsDEROctetStringContent(binaries)
	if err != nil {
		return nil, err
	}
	element, _, err := asn1ber.Parse(content)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the ASN1Sequence", err)
	}
	if !element.IsUniversal(asn1ber.TagSequence) {
		return nil, model.NewDSSError("Unable to retrieve the ASN1Sequence : the encapsulated object is not a SEQUENCE")
	}
	return element.Encoded(), nil
}

// DSSASN1UtilsAsn1IntegerFromDerOctetString returns the INTEGER encapsulated in the given
// DER OCTET STRING. Port of getAsn1IntegerFromDerOctetString(byte[]).
func DSSASN1UtilsAsn1IntegerFromDerOctetString(binaries []byte) (*big.Int, error) {
	content, err := dssASN1UtilsDEROctetStringContent(binaries)
	if err != nil {
		return nil, err
	}
	element, _, err := asn1ber.Parse(content)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the ASN1Integer", err)
	}
	if !element.IsUniversal(asn1ber.TagInteger) {
		return nil, model.NewDSSError("Unable to retrieve the ASN1Integer : the encapsulated object is not an INTEGER")
	}
	return element.Integer(), nil
}

// dssASN1UtilsDEROctetStringContent ports the private getDEROctetStringContent(byte[]).
func dssASN1UtilsDEROctetStringContent(binaries []byte) ([]byte, error) {
	element, _, err := asn1ber.Parse(binaries)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve the DEROctetString content", err)
	}
	if !element.IsUniversal(asn1ber.TagOctetString) {
		return nil, model.NewDSSError("Unable to retrieve the DEROctetString content : not an OCTET STRING")
	}
	return element.Octets(), nil
}

// getAsn1Encodable(Attribute) and getAsn1Attributes(AttributeTable, ASN1ObjectIdentifier) live
// in cms_certificate_source.go, next to the ESS structures that made the CMS attribute types
// necessary there; the three attribute-table predicates below complete the group.

// DSSASN1UtilsIsEmpty reports whether the attribute table is absent or empty.
// Port of isEmpty(AttributeTable).
func DSSASN1UtilsIsEmpty(attributeTable cmscore.Attributes) bool {
	return len(attributeTable) == 0
}

// DSSASN1UtilsEmptyIfNull returns the given attribute table, or an empty one when it is nil.
// Port of emptyIfNull(AttributeTable).
//
// Java hands back a fresh AttributeTable wrapping an empty Hashtable; a nil cmscore.Attributes
// already behaves like one for every read operation, but callers append to the result, so an
// empty non-nil slice is returned to keep "the table exists" distinguishable.
func DSSASN1UtilsEmptyIfNull(originalAttributeTable cmscore.Attributes) cmscore.Attributes {
	if originalAttributeTable != nil {
		return originalAttributeTable
	}
	return cmscore.Attributes{}
}

// DSSASN1UtilsIsAttributeOfType reports whether the attribute carries the given type.
// Port of isAttributeOfType(Attribute, ASN1ObjectIdentifier).
func DSSASN1UtilsIsAttributeOfType(attribute *cmscore.Attribute, asn1ObjectIdentifier asn1.ObjectIdentifier) bool {
	if attribute == nil {
		return false
	}
	return asn1ObjectIdentifier.Equal(attribute.Type)
}

// DSSASN1UtilsAsn1SignaturePolicyDigest computes the digest of an ASN.1 signature policy
// (used in CAdES). Port of getAsn1SignaturePolicyDigest(DigestAlgorithm, byte[]).
//
// TS 101 733 5.8.1: if the signature policy is defined using ASN.1, the hash is calculated
// on the value without the outer type and length fields.
func DSSASN1UtilsAsn1SignaturePolicyDigest(digestAlgorithm enumerations.DigestAlgorithm, policyBytes []byte) ([]byte, error) {
	element, _, err := asn1ber.Parse(policyBytes)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Cannot convert binaries to ASN1Primitive", err)
	}
	if !element.IsConstructed() || len(element.Children()) < 2 {
		return nil, model.NewDSSError("The signature policy shall be a SEQUENCE of at least 2 elements")
	}
	signPolicyHashAlgIdentifier, err := asn1ber.AlgorithmIdentifierFromElement(element.Children()[0])
	if err != nil {
		return nil, err
	}
	signPolicyInfo := element.Children()[1]

	hashAlgorithmDEREncoded := signPolicyHashAlgIdentifier.DER()
	signPolicyInfoDEREncoded := signPolicyInfo.DEREncoded()
	return dssASN1UtilsDigest(digestAlgorithm, hashAlgorithmDEREncoded, signPolicyInfoDEREncoded)
}

// DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex returns the algorithm identifier found in
// the provided ats-hash-index table, or nil when the table carries none.
// Port of getAlgorithmIdentifier(ASN1Sequence).
func DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(atsHashIndexValue []byte) *AlgorithmIdentifier {
	element, _, err := asn1ber.Parse(atsHashIndexValue)
	if err != nil || !element.IsConstructed() {
		return nil
	}
	if len(element.Children()) <= 3 {
		return nil
	}
	const algorithmIndex = 0
	candidate := element.Children()[algorithmIndex]
	switch {
	case candidate.IsUniversal(asn1ber.TagSequence):
		identifier, err := asn1ber.AlgorithmIdentifierFromElement(candidate)
		if err != nil {
			return nil
		}
		return identifier
	case candidate.IsUniversal(asn1ber.TagOID):
		// TODO (upstream, 16/11/2014): the relevance and usefulness of this test case must
		// be checked (do signatures like this exist?)
		oid, err := candidate.ObjectIdentifier()
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
	oid, err := asn1ber.OIDFromString(digestAlgorithm.OID())
	if err != nil {
		return nil, err
	}
	if enumerations.DigestAlgorithm_SHAKE256_512 == digestAlgorithm {
		// Special case, requiring the parameter definition.
		return NewAlgorithmIdentifierWithParameters(oid, asn1ber.EncodeInteger(big.NewInt(512))), nil
	}
	withNull, known := dssASN1UtilsDigestAlgorithmIdentifierParameters[digestAlgorithm]
	if withNull || !known {
		// An algorithm outside the table gets the DERNull parameters the finder defaults to.
		return NewAlgorithmIdentifierWithParameters(oid, asn1ber.DERNull), nil
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
	element, _, err := asn1ber.Parse(asn1Sequence)
	if err != nil {
		return nil, err
	}
	if !element.IsConstructed() {
		return nil, errors.New("not an ASN1Sequence")
	}
	for _, child := range element.Children() {
		if !child.IsUniversal(asn1ber.TagOctetString) {
			return nil, errors.New("the sequence holds an element that is not an OCTET STRING")
		}
		derOctetStrings = append(derOctetStrings, child.Octets())
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
	element, _, err := asn1ber.Parse(publicKey.Encoded())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute ski from public key", err)
	}
	if !element.IsConstructed() || len(element.Children()) < 2 {
		return nil, model.NewDSSError("Unable to compute ski from public key : malformed SubjectPublicKeyInfo")
	}
	item := element.Children()[1]
	if !item.IsUniversal(asn1ber.TagBitString) {
		return nil, model.NewDSSError("Unable to compute ski from public key : subjectPublicKey is not a BIT STRING")
	}
	return dssASN1UtilsDigest(enumerations.DigestAlgorithm_SHA1, item.BitStringOctets())
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
	element, _, err := asn1ber.Parse(binaries)
	if err != nil || !element.IsConstructed() || len(element.Children()) < 2 {
		// Upstream logs "Unable to decode IssuerSerialV2 textContent '{}' : {}".
		return nil
	}
	names := element.Children()[0]
	if !names.IsConstructed() {
		return nil
	}
	issuerSerial := &IssuerSerial{}
	for _, name := range names.Children() {
		if name.Class() != asn1ber.ClassContextSpecific {
			return nil
		}
		tagNo := int(name.TagNumber())
		value := name.DERContent()
		if tagNo == 4 {
			// directoryName is explicitly tagged: the content is the Name itself.
			inner, _, err := asn1ber.Parse(value)
			if err != nil {
				return nil
			}
			value = inner.DEREncoded()
		}
		issuerSerial.Issuer = append(issuerSerial.Issuer, GeneralName{TagNo: tagNo, Name: value})
	}
	serial := element.Children()[1]
	if !serial.IsUniversal(asn1ber.TagInteger) {
		return nil
	}
	issuerSerial.Serial = serial.Integer()
	if len(element.Children()) > 2 {
		issuerSerial.IssuerUID = element.Children()[2].DEREncoded()
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
	element, _, err := asn1ber.Parse(x500Principal.Encoded())
	if err != nil || !element.IsConstructed() {
		return nil
	}
	for _, set := range element.Children() {
		if !set.IsConstructed() {
			return nil
		}
		for _, sequence := range set.Children() {
			if !sequence.IsConstructed() || len(sequence.Children()) != 2 {
				// Java: "The DLSequence must contains exactly 2 elements."
				return nil
			}
			stringAttributeType := DSSASN1UtilsString(sequence.Children()[0].Encoded())
			stringAttributeValue := DSSASN1UtilsString(sequence.Children()[1].Encoded())
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
	element, rest, err := asn1ber.Parse(attributeValue)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to handle attribute : {}".
		return ""
	}
	return asn1ber.JavaTrim(asn1ber.ValueToString(element))
}

// DSSASN1UtilsExtractAttributeFromX500Principal returns the value of the first attribute
// with the given type found in the principal, or the empty string (Java: null) when the
// principal carries none.
// Port of extractAttributeFromX500Principal(ASN1ObjectIdentifier, X500PrincipalHelper).
func DSSASN1UtilsExtractAttributeFromX500Principal(identifier asn1.ObjectIdentifier, principal *model.X500PrincipalHelper) string {
	element, _, err := asn1ber.Parse(principal.Encoded())
	if err != nil || !element.IsConstructed() {
		return ""
	}
	for _, set := range element.Children() {
		if !set.IsConstructed() {
			continue
		}
		for _, sequence := range set.Children() {
			if !sequence.IsConstructed() || len(sequence.Children()) != 2 {
				continue
			}
			oid, err := sequence.Children()[0].ObjectIdentifier()
			if err != nil || !identifier.Equal(oid) {
				continue
			}
			return asn1ber.ASN1ToString(sequence.Children()[1])
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

// DSSASN1UtilsFirstSignerInformation returns the first signer of a SignerInformationStore,
// warning when the store holds more than one. Port of
// getFirstSignerInformation(SignerInformationStore).
//
// Returns nil for a store without a signer, where Java's iterator().next() raises a
// NoSuchElementException: every call site has already established that there is a signer, and a
// nil result keeps the CMS-typed helper free of an error channel Java does not have either.
func DSSASN1UtilsFirstSignerInformation(signerInformationStore []*cmscore.SignerInfo) *cmscore.SignerInfo {
	if len(signerInformationStore) == 0 {
		return nil
	}
	// Upstream logs "!!! The framework handles only one signer (SignerInformation) !!!" when
	// the store holds more than one.
	return signerInformationStore[0]
}

// DSSASN1UtilsIsASN1SequenceTag reports whether the byte is the identifier octet of an ASN.1
// SEQUENCE. Port of isASN1SequenceTag(byte).
func DSSASN1UtilsIsASN1SequenceTag(tagByte byte) bool {
	// BERTags.SEQUENCE | BERTags.CONSTRUCTED = 0x30
	return (asn1ber.TagSequence | asn1ber.Constructed) == tagByte
}

// DSSASN1UtilsDate reads a DER-encoded X.509 Time (a UTCTime or a GeneralizedTime) and
// returns it, or the zero time.Time when it cannot be read.
// Port of getDate(ASN1Encodable), which returns null on failure.
func DSSASN1UtilsDate(encodable []byte) time.Time {
	element, rest, err := asn1ber.Parse(encodable)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to retrieve the date {}".
		return time.Time{}
	}
	var date time.Time
	switch {
	case element.IsUniversal(asn1ber.TagUTCTime):
		date, err = asn1ber.ParseUTCTime(string(element.Content()))
	case element.IsUniversal(asn1ber.TagGeneralizedTime):
		date, err = asn1ber.ParseGeneralizedTime(string(element.Content()))
	default:
		return time.Time{}
	}
	if err != nil {
		return time.Time{}
	}
	return date
}

// DSSASN1UtilsTimeStampTokenGenerationTime returns the generation time of a timestamp token,
// i.e. TSTInfo.genTime. Port of getTimeStampTokenGenerationTime(TimeStampToken); the
// TSPValidationException the Java overload catches around TimeStampToken construction has no
// counterpart here, since ParseTimeStampToken already rejected an unreadable token before this
// is ever reached.
//
// Completes the TODO(phase-3) marker this used to carry (getRevocationValues(ASN1Encodable)
// and getCertificateRef(OtherCertID), the other two methods it named, were completed earlier
// in cms_crl_source.go and cms_certificate_source.go as DSSASN1UtilsRevocationValues and
// DSSASN1UtilsCertificateRef).
func DSSASN1UtilsTimeStampTokenGenerationTime(timeStampToken *cmscore.TimeStampToken) time.Time {
	return timeStampToken.TSTInfo().GenTime
}

// DSSASN1UtilsIsAsn1Encoded reports whether the binaries are ASN.1 encoded.
// Port of isAsn1Encoded(byte[]).
func DSSASN1UtilsIsAsn1Encoded(binaries []byte) bool {
	if len(binaries) == 0 {
		return false
	}
	_, _, err := asn1ber.Parse(binaries)
	return err == nil
}

// DSSASN1UtilsIsAsn1EncodedSignatureValue reports whether the signature-value binaries are
// an ASN.1 SEQUENCE of two elements. Port of isAsn1EncodedSignatureValue(byte[]).
func DSSASN1UtilsIsAsn1EncodedSignatureValue(binaries []byte) bool {
	element, _, err := asn1ber.Parse(binaries)
	if err != nil || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return false
	}
	return len(element.Children()) == 2
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
	element, _, err := asn1ber.Parse(asn1SignatureValue)
	if err != nil || !element.IsConstructed() || len(element.Children()) != 2 {
		return nil, model.NewDSSError("Unable to convert to plain : the signature value is not a SEQUENCE of 2 elements")
	}
	if !element.Children()[0].IsUniversal(asn1ber.TagInteger) || !element.Children()[1].IsUniversal(asn1ber.TagInteger) {
		return nil, model.NewDSSError("Unable to convert to plain : the signature value components are not INTEGERs")
	}
	r := element.Children()[0].Integer()
	s := element.Children()[1].Integer()

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
	body := append(asn1ber.EncodeInteger(r), asn1ber.EncodeInteger(s)...)
	return asn1ber.WriteSequence(body), nil
}

// DSSASN1UtilsOrderFromSignatureValue returns the order parameter corresponding to the given
// signature value, i.e. max(r, s) + 1. Port of getOrderFromSignatureValue(byte[]).
func DSSASN1UtilsOrderFromSignatureValue(signatureValue []byte) (*big.Int, error) {
	var rValue, sValue *big.Int
	if DSSASN1UtilsIsAsn1EncodedSignatureValue(signatureValue) {
		element, _, err := asn1ber.Parse(signatureValue)
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to extract order from a signature value", err)
		}
		if len(element.Children()) != 2 {
			return nil, model.NewDSSError("Unable to extract order from a signature value : ASN1 Sequence size should be 2!")
		}
		if !element.Children()[0].IsUniversal(asn1ber.TagInteger) || !element.Children()[1].IsUniversal(asn1ber.TagInteger) {
			return nil, model.NewDSSError("Unable to extract order from a signature value : the components are not INTEGERs")
		}
		rValue = element.Children()[0].Integer()
		sValue = element.Children()[1].Integer()
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
	element, rest, err := asn1ber.Parse(directoryStringInstance)
	if err != nil || len(rest) != 0 {
		// Upstream logs "Unable to build a DirectoryString instance. Reason : {}".
		return ""
	}
	switch {
	case element.IsUniversal(asn1ber.TagPrintableString), element.IsUniversal(asn1ber.TagT61String),
		element.IsUniversal(asn1ber.TagUTF8String), element.IsUniversal(asn1ber.TagUniversalString),
		element.IsUniversal(asn1ber.TagBMPString):
		return element.AsString()
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
		oid, err := asn1ber.OIDFromString(oidOrURI)
		if err != nil {
			return nil, err
		}
		encoded := asn1ber.EncodeOID(oid)
		if encoded == nil {
			// Java: new ASN1ObjectIdentifier("0") raises
			// IllegalArgumentException("string 0 not a valid OID"). DSSUtils#isOidCode
			// accepts a single arc, ASN.1 does not.
			return nil, fmt.Errorf("string %s not a valid OID", oidOrURI)
		}
		return encoded, nil
	}
	return asn1ber.WriteTLV(asn1ber.TagIA5String, []byte(oidOrURI)), nil
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

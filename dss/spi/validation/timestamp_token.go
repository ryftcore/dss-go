// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampToken.java (DSS 6.5.RC1).
//
// # BouncyCastle replacements
//
// Upstream wraps org.bouncycastle.tsp.TimeStampToken. internal/cmscore already carries the RFC
// 3161 structures (TimeStampToken, TSTInfo, MessageImprint) and the RFC 5652 CMS around them,
// but it is deliberately parse-only and models no *behaviour*: the two BouncyCastle operations
// this class reaches through - TimeStampToken#validate(SignerInformationVerifier), which
// upstream calls from checkIsSignedBy via isValidTimestamp, and
// SignerInformation#verify(SignerInformationVerifier), which it calls from isValidCMSSignedData -
// therefore live here, in the file of the class that needs them, as the unexported
// timestampTokenValidate / timestampTokenVerifySignerInfo pair. So do the two checks
// BouncyCastle's TimeStampToken(ContentInfo) constructor performs (a single signer, and an ESS
// signing-certificate attribute from which the CertID is read), which upstream gets for free
// when it builds the BouncyCastle token and which the Go constructors below perform instead.
//
// The signature check itself is spi.SignerInformationVerifier, i.e. exactly the object upstream
// obtains from DSSSignerInformationVerifierSecurityFactory.CERTIFICATE_TOKEN_INSTANCE.
//
// # Deviations
//
//   - BouncyCastle's SignerId#match(X509CertificateHolder) becomes DSS's own
//     SignerIdentifier#isRelatedToCertificate: both compare issuer and serial when the
//     SignerIdentifier carries them, but for the subjectKeyIdentifier alternative BouncyCastle
//     falls back to the "MS Outlook" key-id calculation for a certificate without a
//     subjectKeyIdentifier extension, where DSS simply finds no match.
//   - DigestAlgorithm.forOID throws an unchecked IllegalArgumentException for a message-imprint
//     algorithm DSS does not know; DigestAlgorithm() panics with the same message rather than
//     growing an error return that would spread through every matchData variant.
//   - The suppressMatchWarnings flag of the matchData family only silences log output, and this
//     port drops slf4j; the flag is kept for API fidelity but has no observable effect.
//   - synchronized is dropped from isSignedBy, as model.TokenBase already drops it from its own.
package validation

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// timestampTokenOIDKPTimeStamping is
// id-kp-timeStamping OBJECT IDENTIFIER ::= { id-kp 8 }, i.e. 1.3.6.1.5.5.7.3.8, the only key
// purpose a TSA certificate may bear. Port of KeyPurposeId.id_kp_timeStamping.
var timestampTokenOIDKPTimeStamping = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 8}

// timestampTokenOIDSHA1 is id-sha1, the algorithm an ESSCertID (signing-certificate v1) fixes
// for its certHash. Port of OIWObjectIdentifiers.idSHA1.
var timestampTokenOIDSHA1 = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}

// timestampTokenExtendedKeyUsageOID is the extendedKeyUsage extension OID, 2.5.29.37.
var timestampTokenExtendedKeyUsageOID = asn1.ObjectIdentifier{2, 5, 29, 37}

// timestampTokenJavaError carries the exception class name alongside the message, so that the
// invalidity reason upstream builds as `e.getClass().getSimpleName() + " : " + e.getMessage()`
// is simply the error's own rendering.
type timestampTokenJavaError struct {
	// className is the Java exception's getSimpleName().
	className string
	// message is the Java exception's getMessage().
	message string
}

// Error renders the error as upstream renders the exception it stands for.
func (e *timestampTokenJavaError) Error() string {
	return e.className + " : " + e.message
}

// timestampTokenTSPValidationError stands for org.bouncycastle.tsp.TSPValidationException.
func timestampTokenTSPValidationError(format string, args ...any) error {
	return &timestampTokenJavaError{className: "TSPValidationException", message: fmt.Sprintf(format, args...)}
}

// timestampTokenTSPError stands for org.bouncycastle.tsp.TSPException.
func timestampTokenTSPError(format string, args ...any) error {
	return &timestampTokenJavaError{className: "TSPException", message: fmt.Sprintf(format, args...)}
}

// timestampTokenCMSError stands for org.bouncycastle.cms.CMSException and its subclasses, whose
// simple name upstream shows in the invalidity reason.
func timestampTokenCMSError(className, format string, args ...any) error {
	return &timestampTokenJavaError{className: className, message: fmt.Sprintf(format, args...)}
}

// timestampTokenCertID is the reference to the TSA certificate an RFC 5035 signing-certificate
// or signing-certificate-v2 signed attribute carries. Port of the private
// org.bouncycastle.tsp.TimeStampToken.CertID.
type timestampTokenCertID struct {
	// hashAlgorithm is the algorithm certHash was computed with: id-sha1 for an ESSCertID, the
	// ESSCertIDv2 hashAlgorithm (id-sha256 by default) otherwise.
	hashAlgorithm *asn1ber.AlgorithmIdentifier
	// certHash is the digest of the TSA certificate.
	certHash []byte
	// issuerSerial identifies the TSA certificate, nil when the optional field is absent.
	issuerSerial *asn1ber.IssuerSerial
}

// timestampTokenParseCertID reads the CertID of the TSA signer, i.e. the first value of its
// signing-certificate or signing-certificate-v2 signed attribute. Port of the tail of the
// BouncyCastle TimeStampToken(ContentInfo) constructor.
func timestampTokenParseCertID(signerInfo *cmscore.SignerInfo) (*timestampTokenCertID, error) {
	if attribute := signerInfo.SignedAttributes.Get(spi.OIDIdAaSigningCertificate); attribute != nil {
		values := attribute.ValueEncodings()
		if len(values) == 0 {
			return nil, timestampTokenTSPValidationError("no signing certificate attribute found, time stamp invalid.")
		}
		signingCertificate, err := spi.ParseSigningCertificate(values[0])
		if err != nil {
			return nil, timestampTokenTSPError("unable to read signing certificate attribute: %s", err.Error())
		}
		certs, err := signingCertificate.Certs()
		if err != nil || len(certs) == 0 {
			return nil, timestampTokenTSPError("unable to read signing certificate attribute: no ESSCertID")
		}
		return &timestampTokenCertID{
			hashAlgorithm: asn1ber.NewAlgorithmIdentifier(timestampTokenOIDSHA1),
			certHash:      certs[0].CertHash,
			issuerSerial:  certs[0].IssuerSerial,
		}, nil
	}

	attribute := signerInfo.SignedAttributes.Get(spi.OIDIdAaSigningCertificateV2)
	if attribute == nil {
		return nil, timestampTokenTSPValidationError("no signing certificate attribute found, time stamp invalid.")
	}
	values := attribute.ValueEncodings()
	if len(values) == 0 {
		return nil, timestampTokenTSPValidationError("no signing certificate attribute found, time stamp invalid.")
	}
	signingCertificate, err := spi.ParseSigningCertificateV2(values[0])
	if err != nil {
		return nil, timestampTokenTSPError("unable to read signing certificate attribute: %s", err.Error())
	}
	certs, err := signingCertificate.Certs()
	if err != nil || len(certs) == 0 {
		return nil, timestampTokenTSPError("unable to read signing certificate attribute: no ESSCertIDv2")
	}
	return &timestampTokenCertID{
		hashAlgorithm: certs[0].HashAlgorithm,
		certHash:      certs[0].CertHash,
		issuerSerial:  certs[0].IssuerSerial,
	}, nil
}

// timestampTokenValidateCertificate checks that the given certificate may issue time-stamps.
// Port of org.bouncycastle.tsp.TSPUtil#validateCertificate(X509CertificateHolder).
func timestampTokenValidateCertificate(certificate *model.CertificateToken) error {
	if certificate.Certificate().Version != 3 {
		// Java raises an IllegalArgumentException here, not a TSPValidationException.
		return &timestampTokenJavaError{
			className: "IllegalArgumentException",
			message:   "Certificate must have an ExtendedKeyUsage extension.",
		}
	}
	present := false
	for _, certificateExtension := range certificate.Certificate().Extensions {
		if certificateExtension.Id.Equal(timestampTokenExtendedKeyUsageOID) {
			present = true
			break
		}
	}
	if !present {
		return timestampTokenTSPValidationError("Certificate must have an ExtendedKeyUsage extension.")
	}
	extendedKeyUsage := spi.CertificateExtensionsUtilsExtendedKeyUsage(certificate)
	if extendedKeyUsage == nil {
		return timestampTokenTSPValidationError("Certificate must have an ExtendedKeyUsage extension.")
	}
	if !extendedKeyUsage.IsCritical() {
		return timestampTokenTSPValidationError("Certificate must have an ExtendedKeyUsage extension marked as critical.")
	}
	oids := extendedKeyUsage.Oids()
	if len(oids) != 1 || oids[0] != timestampTokenOIDKPTimeStamping.String() {
		return timestampTokenTSPValidationError("ExtendedKeyUsage not solely time stamping.")
	}
	return nil
}

// timestampTokenSignedAttributeValue returns the single value of the given signed attribute, or
// nil when the attribute is absent. It fails when the attribute occurs more than once or holds
// more than one value, which is what BouncyCastle's
// SignerInformation#getSingleValuedSignedAttribute does.
func timestampTokenSignedAttributeValue(signerInfo *cmscore.SignerInfo, attributeType asn1.ObjectIdentifier,
	name string) (*asn1ber.Element, error) {
	attributes := signerInfo.SignedAttributes.GetAll(attributeType)
	if len(attributes) == 0 {
		return nil, nil
	}
	if len(attributes) > 1 {
		return nil, timestampTokenCMSError("CMSException",
			"The SignedAttributes in a signerInfo MUST NOT include multiple instances of the %s attribute", name)
	}
	values := attributes[0].Values
	if len(values) != 1 {
		return nil, timestampTokenCMSError("CMSException",
			"A %s attribute MUST have a single attribute value", name)
	}
	return values[0], nil
}

// timestampTokenVerifySignerInfo checks the cryptographic validity of a CMS SignerInfo the way
// org.bouncycastle.cms.SignerInformation#verify(SignerInformationVerifier) does: the signing
// certificate has to be valid at the claimed signing time, the content-type and message-digest
// signed attributes have to agree with the content, and the signature has to verify over the
// DER re-encoding of the signed attributes (over the content itself when there are none).
//
// The two failure modes of the Java method are kept apart, because their callers treat them
// differently: a signature that simply does not verify returns (false, nil), where every
// structural problem - Java's CMSException and its subclasses - returns an error carrying the
// exception's simple name.
func timestampTokenVerifySignerInfo(signerInfo *cmscore.SignerInfo, cms *cmscore.CMS,
	candidate *model.CertificateToken, verifier *spi.SignerInformationVerifier,
	signatureAlgorithm enumerations.SignatureAlgorithm) (bool, error) {
	// SignerInformation#verify: the verifier's certificate has to be valid at signingTime.
	signingTime, err := timestampTokenSignedAttributeValue(signerInfo, cmscore.OIDSigningTime, "signing-time")
	if err != nil {
		return false, err
	}
	if signingTime != nil {
		if date := spi.DSSASN1UtilsDate(signingTime.Encoded()); !date.IsZero() && !candidate.IsValidOn(date) {
			return false, timestampTokenCMSError("CMSVerifierCertificateNotValidException",
				"verifier not valid at signingTime")
		}
	}

	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(signerInfo.DigestAlgorithm.Algorithm.String())
	if err != nil {
		return false, timestampTokenCMSError("CMSException", "can't create content verifier: %s", err.Error())
	}
	content := cms.SignedContent()
	if content == nil && !signerInfo.HasSignedAttributes() {
		return false, timestampTokenCMSError("CMSException",
			"data not encapsulated in signature - use detached constructor.")
	}
	var resultDigest []byte
	if content != nil {
		resultDigest, err = spi.DSSUtilsDigest(digestAlgorithm, content)
		if err != nil {
			return false, timestampTokenCMSError("CMSException", "can't create digest calculator: %s", err.Error())
		}
	}

	// RFC 5652 clause 11.1: check the content-type attribute is correct.
	contentType, err := timestampTokenSignedAttributeValue(signerInfo, cmscore.OIDContentType, "content-type")
	if err != nil {
		return false, err
	}
	if contentType == nil {
		if signerInfo.HasSignedAttributes() {
			return false, timestampTokenCMSError("CMSException",
				"The content-type attribute type MUST be present whenever signed attributes are present in signed-data")
		}
	} else {
		signedContentType, err := contentType.ObjectIdentifier()
		if err != nil {
			return false, timestampTokenCMSError("CMSException",
				"content-type attribute value not of ASN.1 type 'OBJECT IDENTIFIER'")
		}
		if !signedContentType.Equal(cms.SignedContentType()) {
			return false, timestampTokenCMSError("CMSException",
				"content-type attribute value does not match eContentType")
		}
	}

	// RFC 5652 clause 11.2: check the message-digest attribute is correct.
	messageDigest, err := timestampTokenSignedAttributeValue(signerInfo, cmscore.OIDMessageDigest, "message-digest")
	if err != nil {
		return false, err
	}
	if messageDigest == nil {
		if signerInfo.HasSignedAttributes() {
			return false, timestampTokenCMSError("CMSException",
				"the message-digest signed attribute type MUST be present when there are any signed attributes present")
		}
	} else {
		if !messageDigest.IsUniversal(asn1ber.TagOctetString) {
			return false, timestampTokenCMSError("CMSException",
				"message-digest attribute value not of ASN.1 type 'OCTET STRING'")
		}
		if !bytes.Equal(resultDigest, messageDigest.Octets()) {
			return false, timestampTokenCMSError("CMSSignerDigestMismatchException",
				"message-digest attribute value does not match calculated value")
		}
	}

	// The signature covers the DER re-encoding of the signed attributes when they are present,
	// and the content itself otherwise.
	signedContent := content
	if signerInfo.HasSignedAttributes() {
		signedContent = signerInfo.SignedAttributesDER()
	}
	// ContentVerifier#verify returns a boolean; an unverifiable algorithm is the
	// OperatorCreationException the Java code turns into "can't create content verifier".
	if err := verifier.Verify(signatureAlgorithm, signedContent, signerInfo.Signature); err != nil {
		if strings.HasPrefix(err.Error(), "NoSuchAlgorithmException") {
			return false, timestampTokenCMSError("CMSException", "can't create content verifier: %s", err.Error())
		}
		return false, nil
	}
	return true, nil
}

// timestampTokenValidate validates the time-stamp token against the given TSA certificate the
// way org.bouncycastle.tsp.TimeStampToken#validate(SignerInformationVerifier) does.
func timestampTokenValidate(timeStamp *cmscore.TimeStampToken, signerInfo *cmscore.SignerInfo,
	certID *timestampTokenCertID, candidate *model.CertificateToken, verifier *spi.SignerInformationVerifier,
	signatureAlgorithm enumerations.SignatureAlgorithm) error {
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(certID.hashAlgorithm.Algorithm.String())
	if err != nil {
		return timestampTokenTSPError("unable to calculate certificate hash: %s", err.Error())
	}
	certificateHash, err := spi.DSSUtilsDigest(digestAlgorithm, candidate.Encoded())
	if err != nil {
		return timestampTokenTSPError("unable to calculate certificate hash: %s", err.Error())
	}
	if !bytes.Equal(certID.certHash, certificateHash) {
		return timestampTokenTSPValidationError("certificate hash does not match certID hash.")
	}

	if certID.issuerSerial != nil {
		if certID.issuerSerial.Serial == nil || candidate.SerialNumber() == nil ||
			certID.issuerSerial.Serial.Cmp(candidate.SerialNumber()) != 0 {
			return timestampTokenTSPValidationError("certificate serial number does not match certID for signature.")
		}
		found := false
		for index := range certID.issuerSerial.Issuer {
			generalName := certID.issuerSerial.Issuer[index]
			if generalName.TagNo != 4 {
				continue
			}
			issuer, err := spi.DSSASN1UtilsToX500Principal(generalName.Name)
			if err != nil {
				continue
			}
			if spi.DSSASN1UtilsX500PrincipalAreEquals(issuer, candidate.IssuerX500Principal()) {
				found = true
				break
			}
		}
		if !found {
			return timestampTokenTSPValidationError("certificate name does not match certID for signature. ")
		}
	}

	if err := timestampTokenValidateCertificate(candidate); err != nil {
		return err
	}
	if !candidate.IsValidOn(timeStamp.TSTInfo().GenTime) {
		return timestampTokenTSPValidationError("certificate not valid when time stamp created.")
	}
	verified, err := timestampTokenVerifySignerInfo(signerInfo, timeStamp.CMS(), candidate, verifier,
		signatureAlgorithm)
	if err != nil {
		// BouncyCastle wraps the CMSException in a TSPException("unable to verify signature: ").
		var message string
		if javaError, ok := err.(*timestampTokenJavaError); ok {
			message = javaError.message
		} else {
			message = err.Error()
		}
		return timestampTokenTSPError("unable to verify signature: %s", message)
	}
	if !verified {
		return timestampTokenTSPValidationError("signature not created by certificate.")
	}
	return nil
}

// TimestampToken is a SignedToken containing a TimeStamp.
type TimestampToken struct {
	model.TokenBase

	// timeStamp is the parsed representation of a TimeStamp Token, replacing BouncyCastle's.
	timeStamp *cmscore.TimeStampToken
	// tsaSignerInfo is the single signer of the time-stamp, i.e. the TSA signature. It stands
	// in for BouncyCastle's TimeStampToken.tsaSignerInfo.
	tsaSignerInfo *cmscore.SignerInfo
	// certID is the reference to the TSA certificate the signing-certificate(-v2) signed
	// attribute carries. It stands in for BouncyCastle's TimeStampToken.certID.
	certID *timestampTokenCertID

	// timeStampType is the type of the timestamp relatively to the signature.
	timeStampType enumerations.TimestampType

	// certificateSource is the certificate source extracted from the timestamp.
	certificateSource *TimestampCertificateSource
	// crlSource is the CRL source extracted from the timestamp.
	crlSource *TimestampCRLSource
	// ocspSource is the OCSP source extracted from the timestamp.
	ocspSource *TimestampOCSPSource

	// timestampedReferences is the list of references to tokens covered (protected) by the
	// timestamp.
	timestampedReferences []*TimestampedReference

	// identifierBuilder builds the time-stamp token unique identifier.
	identifierBuilder TimestampTokenIdentifierBuilder

	// processed defines whether the timestamp has been validated.
	processed bool

	// messageImprint is the computed message-imprint; a nil value means "not computed yet".
	messageImprint *model.Digest

	// messageImprintData defines whether the message-imprint has been found.
	messageImprintData bool

	// messageImprintIntact defines whether the computed message-imprint is intact; nil stands
	// for Java's null Boolean, i.e. "never checked".
	messageImprintIntact *bool

	// filename is set in case of a detached timestamp.
	filename string

	// timestampScopes is only present for detached timestamps.
	timestampScopes []scope.SignatureScope

	// manifestFile is the timestamped manifest file, when applicable (ASiC with CAdES).
	manifestFile *model.ManifestFile

	// timestampIncludes holds the Includes of a XAdES IndividualDataObjectsTimeStamp.
	timestampIncludes []*TimestampInclude

	// referenceValidations holds, for an evidence record time-stamp, the references it covers.
	referenceValidations []*model.ReferenceValidation

	// detachedEvidenceRecords is the list of detached evidence records covering the time-stamp,
	// when applicable.
	detachedEvidenceRecords []EvidenceRecord

	// archiveTimestampType defines, for an archive timestamp, its type.
	archiveTimestampType enumerations.ArchiveTimestampType

	// evidenceRecordTimestampType defines, for an evidence record archive timestamp, its type.
	evidenceRecordTimestampType enumerations.EvidenceRecordTimestampType

	// canonicalizationMethod is used for XAdES timestamps: it indicates the canonicalization
	// method used for message-imprint computation.
	//
	// NOTE: Used for XAdES/JAdES only.
	canonicalizationMethod string

	// tsaX500Principal identifies the TSA that issued the timestamp token.
	//
	// NOTE: Takes a value only for a successfully validated token.
	tsaX500Principal *model.X500Principal

	// candidatesForSigningCertificate is the cached list of signing certificate candidates.
	candidatesForSigningCertificate *spi.CandidatesForSigningCertificate

	// hashIndexStatus contains the validation status of the ats-hash-index(-v3) attribute.
	//
	// NOTE: applicable only for CMS archive-time-stamp-v3 timestamps.
	hashIndexStatus *ArchiveTimestampHashIndexStatus
}

// NewTimestampToken builds a time-stamp token over the given binaries.
// Port of the TimestampToken(byte[], TimestampType) constructor.
func NewTimestampToken(binaries []byte, timestampType enumerations.TimestampType) (*TimestampToken, error) {
	return NewTimestampTokenWithReferences(binaries, timestampType, []*TimestampedReference{})
}

// NewTimestampTokenWithReferences builds a time-stamp token over the given binaries and the
// references it covers.
// Port of the TimestampToken(byte[], TimestampType, List<TimestampedReference>) constructor.
func NewTimestampTokenWithReferences(binaries []byte, timestampType enumerations.TimestampType,
	timestampedReferences []*TimestampedReference) (*TimestampToken, error) {
	return NewTimestampTokenWithIdentifierBuilder(binaries, timestampType, timestampedReferences, nil)
}

// NewTimestampTokenWithIdentifierBuilder builds a time-stamp token over the given binaries, the
// references it covers and the builder of its identifier.
// Port of the TimestampToken(byte[], TimestampType, List<TimestampedReference>,
// TimestampIdentifierBuilder) constructor.
func NewTimestampTokenWithIdentifierBuilder(binaries []byte, timestampType enumerations.TimestampType,
	timestampedReferences []*TimestampedReference,
	identifierBuilder TimestampTokenIdentifierBuilder) (*TimestampToken, error) {
	cms, err := cmscore.ParseCMS(binaries)
	if err != nil {
		return nil, err
	}
	return NewTimestampTokenFromCMSWithIdentifierBuilder(cms, timestampType, timestampedReferences, identifierBuilder)
}

// NewTimestampTokenFromCMS builds a time-stamp token over the given CMS and the references it
// covers. Port of the TimestampToken(CMSSignedData, TimestampType, List<TimestampedReference>)
// constructor.
func NewTimestampTokenFromCMS(cms *cmscore.CMS, timestampType enumerations.TimestampType,
	timestampedReferences []*TimestampedReference) (*TimestampToken, error) {
	return NewTimestampTokenFromCMSWithIdentifierBuilder(cms, timestampType, timestampedReferences, nil)
}

// NewTimestampTokenFromCMSWithIdentifierBuilder builds a time-stamp token over the given CMS,
// the references it covers and the builder of its identifier.
// Port of the TimestampToken(CMSSignedData, TimestampType, List<TimestampedReference>,
// TimestampIdentifierBuilder) constructor.
func NewTimestampTokenFromCMSWithIdentifierBuilder(cms *cmscore.CMS, timestampType enumerations.TimestampType,
	timestampedReferences []*TimestampedReference,
	identifierBuilder TimestampTokenIdentifierBuilder) (*TimestampToken, error) {
	timeStamp, err := cmscore.TimeStampTokenFromCMS(cms)
	if err != nil {
		return nil, err
	}
	return NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder(timeStamp, timestampType,
		timestampedReferences, identifierBuilder)
}

// NewTimestampTokenFromTimeStampToken builds a time-stamp token over the given RFC 3161 token
// and the references it covers.
// Port of the TimestampToken(TimeStampToken, TimestampType, List<TimestampedReference>)
// constructor.
func NewTimestampTokenFromTimeStampToken(timeStamp *cmscore.TimeStampToken,
	timestampType enumerations.TimestampType,
	timestampedReferences []*TimestampedReference) (*TimestampToken, error) {
	return NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder(timeStamp, timestampType,
		timestampedReferences, nil)
}

// NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder builds a time-stamp token over the
// given RFC 3161 token, the references it covers and the builder of its identifier. It is the
// constructor every other one funnels into.
// Port of the TimestampToken(TimeStampToken, TimestampType, List<TimestampedReference>,
// TimestampIdentifierBuilder) constructor.
//
// The two checks BouncyCastle performs when it builds its own TimeStampToken - exactly one
// signer, and a signing-certificate(-v2) signed attribute to read the CertID from - are
// performed here, since internal/cmscore models the structures only. Upstream reports them by
// letting the BouncyCastle constructor throw; here they are returned as errors.
func NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder(timeStamp *cmscore.TimeStampToken,
	timestampType enumerations.TimestampType, timestampedReferences []*TimestampedReference,
	identifierBuilder TimestampTokenIdentifierBuilder) (*TimestampToken, error) {
	signerInfos := timeStamp.CMS().SignerInfos()
	if len(signerInfos) != 1 {
		return nil, &timestampTokenJavaError{
			className: "IllegalArgumentException",
			message: fmt.Sprintf("Time-stamp token signed by %d signers, but it must contain just the TSA signature.",
				len(signerInfos)),
		}
	}
	tsaSignerInfo := signerInfos[0]
	certID, err := timestampTokenParseCertID(tsaSignerInfo)
	if err != nil {
		return nil, err
	}

	certificateSource, err := NewTimestampCertificateSource(timeStamp)
	if err != nil {
		return nil, err
	}
	ocspSource, err := newTimestampOCSPSource(timeStamp)
	if err != nil {
		return nil, err
	}
	crlSource, err := newTimestampCRLSource(timeStamp)
	if err != nil {
		return nil, err
	}

	token := &TimestampToken{
		TokenBase:             model.NewTokenBase(),
		timeStamp:             timeStamp,
		tsaSignerInfo:         tsaSignerInfo,
		certID:                certID,
		timeStampType:         timestampType,
		certificateSource:     certificateSource,
		ocspSource:            ocspSource,
		crlSource:             crlSource,
		timestampedReferences: timestampedReferences,
		identifierBuilder:     identifierBuilder,
	}
	// Registers the token with model.TokenBase, so that the identifier and the signature
	// dispatch reach this type.
	token.InitToken(token)
	return token, nil
}

// IssuerX500Principal returns the TSA that issued the token, which takes a value only after a
// successful IsSignedByToken. Port of the getIssuerX500Principal() override.
func (t *TimestampToken) IssuerX500Principal() *model.X500Principal {
	return t.tsaX500Principal
}

// Abbreviation returns the DSS abbreviation of the token, used for debugging.
// Port of the getAbbreviation() override.
func (t *TimestampToken) Abbreviation() string {
	return string(t.timeStampType) + ": " + t.DSSIDAsString() + ": " +
		spi.DSSUtilsFormatDateToRFC(t.timeStamp.TSTInfo().GenTime)
}

// CertificateSource returns the TimestampCertificateSource of the timestamp.
// Port of getCertificateSource().
func (t *TimestampToken) CertificateSource() *TimestampCertificateSource {
	return t.certificateSource
}

// CRLSource returns the TimestampCRLSource of the timestamp. Port of getCRLSource().
func (t *TimestampToken) CRLSource() *TimestampCRLSource {
	return t.crlSource
}

// OCSPSource returns the TimestampOCSPSource of the timestamp. Port of getOCSPSource().
func (t *TimestampToken) OCSPSource() *TimestampOCSPSource {
	return t.ocspSource
}

// IsValid reports whether the signature is intact and the message-imprint matches the computed
// message-imprint. Port of the isValid() override.
//
// NOTE: IsSignedByToken must be called before calling this method; and, as in Java, the call
// fails when no matchData variant has run yet - see IsMessageImprintDataIntact.
func (t *TimestampToken) IsValid() bool {
	return t.IsSignatureIntact() && t.IsMessageImprintDataFound() && t.IsMessageImprintDataIntact() &&
		t.AreReferenceValidationsValid()
}

// IsSignedByToken checks whether the timestamp token is signed by the given certificate.
// Port of the isSignedBy(CertificateToken) override.
//
// Upstream declares the method synchronized; the Go port is not goroutine-safe, matching
// model.TokenBase.
func (t *TimestampToken) IsSignedByToken(certificateToken *model.CertificateToken) bool {
	if publicKeyOfTheSigner := t.PublicKeyOfTheSigner(); publicKeyOfTheSigner != nil {
		return publicKeyOfTheSigner.Equals(certificateToken.PublicKey())
	} else if enumerations.SignatureValidityValid == t.CheckIsSignedByToken(certificateToken) {
		if !t.IsSelfSigned() {
			t.SetPublicKeyOfTheSigner(certificateToken.PublicKey())
		}
		return true
	}
	return false
}

// IsSignedBy is not supported for a timestamp token. Port of the isSignedBy(PublicKey) override,
// whose UnsupportedOperationException becomes a panic carrying the same message.
func (t *TimestampToken) IsSignedBy(publicKey *model.PublicKey) bool {
	panic("Use method isSignedBy(certificateToken) for a TimestampToken validation!")
}

// CheckIsSignedByToken checks whether the timestamp is signed by the given certificate, keeping
// the token's validity state, TSA name and signature algorithm up to date.
// Port of the protected checkIsSignedBy(CertificateToken).
//
// The timestamp is validated first as an RFC 3161 token and, if that fails, as a plain CMS
// SignedData, exactly as upstream does.
func (t *TimestampToken) CheckIsSignedByToken(candidate *model.CertificateToken) enumerations.SignatureValidity {
	signerIdentifier, err := spi.DSSASN1UtilsToSignerIdentifierFromSignerID(t.tsaSignerInfo.SID)
	if err != nil {
		return enumerations.SignatureValidityInvalid
	}
	related, err := signerIdentifier.IsRelatedToCertificate(candidate)
	if err != nil || !related {
		return enumerations.SignatureValidityInvalid
	}

	verifier, err := spi.DSSSignerInformationVerifierSecurityFactoryCertificateTokenInstance.Build(candidate)
	if err != nil {
		t.SetSignatureValidity(enumerations.SignatureValidityInvalid)
		t.SetInvalidityReason(err.Error())
		return t.SignatureValidity()
	}
	signatureAlgorithm, algorithmError := t.signatureAlgorithm(candidate)

	// Try firstly to validate as a Timestamp and if that fails try to validate the timestamp
	// as a CMSSignedData.
	if algorithmError == nil &&
		(t.isValidTimestamp(candidate, verifier, signatureAlgorithm) ||
			t.isValidCMSSignedData(candidate, verifier, signatureAlgorithm)) {
		t.SetSignatureValidity(enumerations.SignatureValidityValid)
		t.tsaX500Principal = candidate.Subject().Principal()
		t.SetSignatureAlgorithm(signatureAlgorithm)
	} else {
		if algorithmError != nil {
			t.SetInvalidityReason(algorithmError.Error())
		}
		t.SetSignatureValidity(enumerations.SignatureValidityInvalid)
	}
	return t.SignatureValidity()
}

// signatureAlgorithm computes the algorithm the TSA signed with, the way upstream does once the
// verification succeeded: from the signature algorithm identifier and its parameters for
// RSASSA-PSS, and from the candidate's key type combined with the SignerInfo digest algorithm
// otherwise.
func (t *TimestampToken) signatureAlgorithm(candidate *model.CertificateToken) (enumerations.SignatureAlgorithm, error) {
	signatureAlgorithmOID := t.tsaSignerInfo.SignatureAlgorithm.Algorithm.String()
	if enumerations.EncryptionAlgorithmRSASSAPSS.OID() == signatureAlgorithmOID {
		return enumerations.SignatureAlgorithmForOIDAndParams(signatureAlgorithmOID,
			t.tsaSignerInfo.SignatureAlgorithm.Parameters)
	}
	encryptionAlgorithm, err := enumerations.EncryptionAlgorithmForName(candidate.PublicKey().Algorithm())
	if err != nil {
		return "", err
	}
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(t.tsaSignerInfo.DigestAlgorithm.Algorithm.String())
	if err != nil {
		return "", err
	}
	return enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm), nil
}

// isValidTimestamp validates the timestamp, the signing certificate and the message imprint.
// Port of the private isValidTimestamp(SignerInformationVerifier).
func (t *TimestampToken) isValidTimestamp(candidate *model.CertificateToken,
	verifier *spi.SignerInformationVerifier, signatureAlgorithm enumerations.SignatureAlgorithm) bool {
	if err := timestampTokenValidate(t.timeStamp, t.tsaSignerInfo, t.certID, candidate, verifier,
		signatureAlgorithm); err != nil {
		// Upstream logs "Unable to validate timestamp token : {}".
		t.SetInvalidityReason(err.Error())
		return false
	}
	return true
}

// isValidCMSSignedData validates only the cryptographic validity of the CMS.
// Port of the private isValidCMSSignedData(SignerInformationVerifier).
//
// Only the CMSException branch sets the invalidity reason, exactly as upstream: a
// SignerInformation#verify that merely returns false leaves the reason untouched.
func (t *TimestampToken) isValidCMSSignedData(candidate *model.CertificateToken,
	verifier *spi.SignerInformationVerifier, signatureAlgorithm enumerations.SignatureAlgorithm) bool {
	verified, err := timestampTokenVerifySignerInfo(t.tsaSignerInfo, t.timeStamp.CMS(), candidate, verifier,
		signatureAlgorithm)
	if err != nil {
		// Upstream logs "Unable to validate the related CMSSignedData : {}".
		t.SetInvalidityReason(err.Error())
		return false
	}
	return verified
}

// CheckIsSignedBy is not supported for a timestamp token. Port of the protected
// checkIsSignedBy(PublicKey) override, whose UnsupportedOperationException becomes a panic
// carrying the same message. It is the model.TokenOverrides member of this type.
func (t *TimestampToken) CheckIsSignedBy(publicKey *model.PublicKey) enumerations.SignatureValidity {
	panic("Use method checkIsSignedBy(certificateToken) for a TimestampToken validation!")
}

// MatchDataDocument checks whether the TimeStampToken matches the signed data.
// Port of matchData(DSSDocument).
func (t *TimestampToken) MatchDataDocument(timestampedData model.DSSDocument) (bool, error) {
	return t.MatchDataDocumentSuppressingWarnings(timestampedData, false)
}

// MatchDataDocumentSuppressingWarnings checks whether the TimeStampToken matches the signed
// data. Port of matchData(DSSDocument, boolean).
//
// Java's DSSException on an unreadable document becomes the returned error; the returned bool
// is then the value Java's matchData would have produced.
func (t *TimestampToken) MatchDataDocumentSuppressingWarnings(timestampedData model.DSSDocument,
	suppressMatchWarnings bool) (bool, error) {
	t.processed = true

	t.messageImprintData = timestampedData != nil
	t.setMessageImprintIntact(false)

	if !t.messageImprintData {
		// Upstream logs "Timestamped data not found !".
		return false, nil
	}

	currentMessageImprint := t.MessageImprint()
	computedDigest, err := timestampedData.DigestValue(currentMessageImprint.Algorithm())
	if err != nil {
		return false, err
	}
	return t.MatchDataSuppressingWarnings(computedDigest, suppressMatchWarnings), nil
}

// MatchDataMessageDigest checks whether the TimeStampToken matches the message-imprint digest,
// with warning enabled. Port of matchData(DSSMessageDigest).
func (t *TimestampToken) MatchDataMessageDigest(messageDigest model.DSSMessageDigest) bool {
	return t.MatchDataMessageDigestSuppressingWarnings(messageDigest, false)
}

// MatchDataMessageDigestSuppressingWarnings checks whether the TimeStampToken matches the
// message-imprint digest. Port of matchData(DSSMessageDigest, boolean).
//
// Java's `messageDigest == null` and `messageDigest.isEmpty()` share a branch, and a Go
// DSSMessageDigest value cannot be null, so the zero value stands for both.
func (t *TimestampToken) MatchDataMessageDigestSuppressingWarnings(messageDigest model.DSSMessageDigest,
	suppressMatchWarnings bool) bool {
	t.processed = true

	if messageDigest.IsEmpty() {
		// Upstream logs "Invalid or incomplete message-digest has been provided for timestamp
		// verification!" unless the warnings are suppressed.
		t.setMessageImprintIntact(false)

	} else if t.DigestAlgorithm() != messageDigest.Algorithm() {
		// Upstream logs "DigestAlgorithm '{}' used in the provided message-digest does not
		// match the one used in the timestamp token '{}'!" unless the warnings are suppressed.
		t.setMessageImprintIntact(false)

	} else {
		t.setMessageImprintIntact(t.MatchDataSuppressingWarnings(messageDigest.Value(), suppressMatchWarnings))
	}
	return utils.IsTrue(t.messageImprintIntact)
}

// MatchData checks whether the TimeStampToken matches the given expected message-imprint digest
// value. Port of matchData(byte[]).
func (t *TimestampToken) MatchData(expectedMessageImprintDigest []byte) bool {
	return t.MatchDataSuppressingWarnings(expectedMessageImprintDigest, false)
}

// MatchDataSuppressingWarnings checks whether the TimeStampToken matches the given expected
// message-imprint digest value. Port of matchData(byte[], boolean).
func (t *TimestampToken) MatchDataSuppressingWarnings(expectedMessageImprintDigest []byte,
	suppressMatchWarnings bool) bool {
	t.processed = true

	t.messageImprintData = expectedMessageImprintDigest != nil
	t.setMessageImprintIntact(false)

	if t.messageImprintData {
		currentMessageImprint := t.MessageImprint()
		t.setMessageImprintIntact(bytes.Equal(expectedMessageImprintDigest, currentMessageImprint.Value()))
		// Upstream logs the provided and the embedded digest, and the outcome, when the
		// message-imprint does not match and the warnings are not suppressed.
	}
	// Upstream logs "Timestamped data not found !" otherwise.

	return utils.IsTrue(t.messageImprintIntact)
}

// setMessageImprintIntact assigns the nullable messageImprintIntact field.
func (t *TimestampToken) setMessageImprintIntact(intact bool) {
	t.messageImprintIntact = &intact
}

// IsProcessed reports whether the timestamp's signature has been validated.
// Port of isProcessed().
func (t *TimestampToken) IsProcessed() bool {
	return t.processed
}

// TimeStampType retrieves the type of the timestamp token. Port of getTimeStampType().
func (t *TimestampToken) TimeStampType() enumerations.TimestampType {
	return t.timeStampType
}

// GenerationTime retrieves the timestamp generation time. Port of getGenerationTime().
func (t *TimestampToken) GenerationTime() time.Time {
	return t.timeStamp.TSTInfo().GenTime
}

// CreationDate returns the timestamp generation time. Port of the getCreationDate() override.
func (t *TimestampToken) CreationDate() time.Time {
	return t.GenerationTime()
}

// MessageImprint returns the embedded message-imprint value, i.e. the algorithm and the value.
// Port of getMessageImprint().
func (t *TimestampToken) MessageImprint() model.Digest {
	if t.messageImprint == nil {
		messageImprint := model.NewDigest(t.DigestAlgorithm(), t.timeStamp.TSTInfo().MessageImprint.HashedMessage)
		t.messageImprint = &messageImprint
	}
	return *t.messageImprint
}

// DigestAlgorithm returns the DigestAlgorithm used for the message-imprint computation of the
// timestamp token. Port of getDigestAlgorithm().
//
// DigestAlgorithm.forOID throws an unchecked IllegalArgumentException for an algorithm DSS does
// not know; the Go counterpart of that unchecked throw is a panic carrying the same message.
func (t *TimestampToken) DigestAlgorithm() enumerations.DigestAlgorithm {
	oid := t.timeStamp.TSTInfo().MessageImprint.HashAlgorithm.Algorithm.String()
	digestAlgorithm, err := enumerations.DigestAlgorithmForOID(oid)
	if err != nil {
		panic(err.Error())
	}
	return digestAlgorithm
}

// IsMessageImprintDataFound reports whether the data for the message-imprint computation has
// been found. Port of isMessageImprintDataFound().
func (t *TimestampToken) IsMessageImprintDataFound() bool {
	return t.messageImprintData
}

// IsMessageImprintDataIntact reports whether the message-imprint data is intact. A matchData
// variant must have been invoked before. Port of isMessageImprintDataIntact(), whose
// IllegalStateException becomes a panic carrying the same message.
func (t *TimestampToken) IsMessageImprintDataIntact() bool {
	if !t.processed {
		panic("Invoke matchData(byte[] data) method before!")
	}
	return utils.IsTrue(t.messageImprintIntact)
}

// Filename returns the file name of a detached timestamp. Port of getFilename().
func (t *TimestampToken) Filename() string {
	return t.filename
}

// SetFilename sets the filename of a detached timestamp. Port of setFilename(String).
func (t *TimestampToken) SetFilename(filename string) {
	t.filename = filename
}

// ManifestFile returns the covered manifest file.
// NOTE: applicable only for ASiC-E CAdES. Port of getManifestFile().
func (t *TimestampToken) ManifestFile() *model.ManifestFile {
	return t.manifestFile
}

// SetManifestFile sets the manifest file covered by the current timestamp.
// NOTE: applicable only for ASiC-E CAdES. Port of setManifestFile(ManifestFile).
func (t *TimestampToken) SetManifestFile(manifestFile *model.ManifestFile) {
	t.manifestFile = manifestFile
}

// TimestampedReferences gets the list of TimestampedReferences covered by the current timestamp.
// Port of getTimestampedReferences().
func (t *TimestampToken) TimestampedReferences() []*TimestampedReference {
	return t.timestampedReferences
}

// SetTimestampedReferences sets the list of TimestampedReferences covered by the current
// timestamp. Java has no equivalent method: it mutates the List<TimestampedReference>
// getTimestampedReferences() returns in place (List reference semantics) at several call sites
// in AbstractTimestampSource (incorporateArchiveTimestampReferences, processExternalTimestamp,
// ...). Go slices returned by value do not alias the field they came from, so this setter is
// the only way to make such a mutation observable to later callers.
func (t *TimestampToken) SetTimestampedReferences(timestampedReferences []*TimestampedReference) {
	t.timestampedReferences = timestampedReferences
}

// ArchiveTimestampType gets the ArchiveTimestampType, when applicable; the empty value otherwise.
// Port of getArchiveTimestampType().
func (t *TimestampToken) ArchiveTimestampType() enumerations.ArchiveTimestampType {
	return t.archiveTimestampType
}

// SetArchiveTimestampType sets the subtype of an archive timestamp.
// Port of setArchiveTimestampType(ArchiveTimestampType).
func (t *TimestampToken) SetArchiveTimestampType(archiveTimestampType enumerations.ArchiveTimestampType) {
	t.archiveTimestampType = archiveTimestampType
}

// EvidenceRecordTimestampType gets the EvidenceRecordTimestampType, when applicable; the empty
// value otherwise. Port of getEvidenceRecordTimestampType().
func (t *TimestampToken) EvidenceRecordTimestampType() enumerations.EvidenceRecordTimestampType {
	return t.evidenceRecordTimestampType
}

// SetEvidenceRecordTimestampType sets the EvidenceRecordTimestampType for an evidence record's
// time-stamp. Port of setEvidenceRecordTimestampType(EvidenceRecordTimestampType).
func (t *TimestampToken) SetEvidenceRecordTimestampType(
	evidenceRecordTimestampType enumerations.EvidenceRecordTimestampType) {
	t.evidenceRecordTimestampType = evidenceRecordTimestampType
}

// CanonicalizationMethod returns the canonicalization method used by the timestamp. Applies only
// to XAdES timestamps. Port of getCanonicalizationMethod().
func (t *TimestampToken) CanonicalizationMethod() string {
	return t.canonicalizationMethod
}

// SetCanonicalizationMethod sets the canonicalization method used by the timestamp. Applies only
// to XAdES timestamps. Port of setCanonicalizationMethod(String).
func (t *TimestampToken) SetCanonicalizationMethod(canonicalizationMethod string) {
	t.canonicalizationMethod = canonicalizationMethod
}

// Encoded returns the DER encoding of the time-stamp token.
// Port of the getEncoded() override, i.e. DSSASN1Utils.getDEREncoded(timeStamp).
func (t *TimestampToken) Encoded() []byte {
	return t.timeStamp.CMS().DEREncoded()
}

// TimestampIncludes returns the references covered by the current timestamp (XAdES
// IndividualDataObjectsTimeStamp). Port of getTimestampIncludes().
func (t *TimestampToken) TimestampIncludes() []*TimestampInclude {
	return t.timestampIncludes
}

// SetTimestampIncludes sets the references covered by the current timestamp (XAdES
// IndividualDataObjectsTimeStamp). Port of setTimestampIncludes(List<TimestampInclude>).
func (t *TimestampToken) SetTimestampIncludes(timestampIncludes []*TimestampInclude) {
	t.timestampIncludes = timestampIncludes
}

// ReferenceValidations returns the timestamped data reference validations (used for evidence
// record timestamps). Port of getReferenceValidations().
func (t *TimestampToken) ReferenceValidations() []*model.ReferenceValidation {
	return t.referenceValidations
}

// SetReferenceValidations sets the timestamped data reference validations (used for evidence
// record timestamps). Port of setReferenceValidations(List<ReferenceValidation>).
func (t *TimestampToken) SetReferenceValidations(referenceValidations []*model.ReferenceValidation) {
	t.referenceValidations = referenceValidations
}

// AreReferenceValidationsValid verifies whether the corresponding reference validations are
// valid. Port of areReferenceValidationsValid().
func (t *TimestampToken) AreReferenceValidationsValid() bool {
	if utils.IsCollectionNotEmpty(t.referenceValidations) {
		for _, referenceValidation := range t.referenceValidations {
			if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference != referenceValidation.Type() &&
				(!referenceValidation.IsFound() || !referenceValidation.IsIntact()) {
				return false
			}
		}
	}
	return true
}

// DetachedEvidenceRecords gets the detached evidence records covering the time-stamp, when
// applicable. Port of getDetachedEvidenceRecords(), which lazily creates the list.
func (t *TimestampToken) DetachedEvidenceRecords() []EvidenceRecord {
	if t.detachedEvidenceRecords == nil {
		t.detachedEvidenceRecords = []EvidenceRecord{}
	}
	return t.detachedEvidenceRecords
}

// AddDetachedEvidenceRecord adds an evidence record to the time-stamp's list.
// Port of addDetachedEvidenceRecord(EvidenceRecord).
func (t *TimestampToken) AddDetachedEvidenceRecord(evidenceRecord EvidenceRecord) {
	t.detachedEvidenceRecords = append(t.DetachedEvidenceRecords(), evidenceRecord)
}

// TimestampScopes returns the scope of the current timestamp (detached timestamps only).
// Port of getTimestampScopes().
func (t *TimestampToken) TimestampScopes() []scope.SignatureScope {
	return t.timestampScopes
}

// SetTimestampScopes sets the timestamp's signature scopes.
// Port of setTimestampScopes(List<SignatureScope>).
func (t *TimestampToken) SetTimestampScopes(timestampScopes []scope.SignatureScope) {
	t.timestampScopes = timestampScopes
}

// Certificates returns the list of wrapped certificates. Port of getCertificates().
func (t *TimestampToken) Certificates() []*model.CertificateToken {
	return t.certificateSource.Certificates()
}

// CertificateRefs returns the contained certificate references. Port of getCertificateRefs().
func (t *TimestampToken) CertificateRefs() []*spi.CertificateRef {
	return t.certificateSource.AllCertificateRefs()
}

// UnsignedAttributes gets the unsigned attribute table of the TSA signer.
// Port of getUnsignedAttributes(), which returns BouncyCastle's AttributeTable.
func (t *TimestampToken) UnsignedAttributes() cmscore.Attributes {
	return t.tsaSignerInfo.UnsignedAttributes
}

// TSTInfoTsa returns the TSTInfo.tsa attribute identifying the timestamp issuer, when the
// attribute is present and decodable, and nil otherwise. Port of getTSTInfoTsa().
func (t *TimestampToken) TSTInfoTsa() *model.X500Principal {
	tsaGeneralName := t.timeStamp.TSTInfo().TSA
	if tsaGeneralName != nil {
		principal, err := spi.DSSASN1UtilsToX500Principal(tsaGeneralName.Name)
		if err == nil {
			return principal
		}
		// Upstream logs "Unable to decode TSTInfo.tsa attribute value to X500Principal.
		// Reason : {}".
	}
	return nil
}

// TimeStamp gets the parsed RFC 3161 representation of the time-stamp token.
// Port of getTimeStamp(), which hands out the BouncyCastle TimeStampToken.
func (t *TimestampToken) TimeStamp() *cmscore.TimeStampToken {
	return t.timeStamp
}

// ToString renders the token with the given indentation. Port of the toString(String) override.
//
// Java wraps the whole body in a try/catch that falls back to getClass().getName(); nothing in
// the Go body can fail, so the fallback has no counterpart.
func (t *TimestampToken) ToString(indentStr string) string {
	var out strings.Builder
	out.WriteString(indentStr)
	out.WriteString("TimestampToken[signedBy=")
	out.WriteString(timestampTokenPrincipalString(t.IssuerX500Principal()))
	out.WriteString(", generated: ")
	out.WriteString(spi.DSSUtilsFormatDateToRFC(t.timeStamp.TSTInfo().GenTime))
	out.WriteString(" / ")
	out.WriteString(string(t.timeStampType))
	out.WriteString("\n")
	if t.IsSignatureIntact() {
		indentStr += "\t"
		out.WriteString(indentStr)
		out.WriteString("Timestamp's signature validity: VALID\n")
		indentStr = indentStr[1:]
	} else if t.InvalidityReason() != "" {
		indentStr += "\t"
		out.WriteString(indentStr)
		out.WriteString("Timestamp's signature validity: INVALID - ")
		out.WriteString(t.InvalidityReason())
		out.WriteString("\n")
		indentStr = indentStr[1:]
	}
	indentStr += "\t"
	if t.messageImprintIntact != nil {
		out.WriteString(indentStr)
		if *t.messageImprintIntact {
			out.WriteString("Timestamp MATCHES the signed data.\n")
		} else {
			out.WriteString("Timestamp DOES NOT MATCH the signed data.\n")
		}
	}
	out.WriteString("]")
	return out.String()
}

// String returns ToString(""). Port of Token#toString().
func (t *TimestampToken) String() string {
	return t.ToString("")
}

// timestampTokenPrincipalString renders a possibly missing principal the way Java's string
// concatenation does.
func timestampTokenPrincipalString(principal *model.X500Principal) string {
	if principal == nil {
		return "null"
	}
	return principal.String()
}

// SignerInformationStoreInfos returns the SignerIdentifiers found in the SignerInformationStore.
// Port of getSignerInformationStoreInfos().
func (t *TimestampToken) SignerInformationStoreInfos() []*spi.SignerIdentifier {
	return t.CertificateSource().AllCertificateIdentifiers()
}

// CandidatesForSigningCertificate returns the object with the signing certificate candidates,
// computing it on first use. Port of getCandidatesForSigningCertificate().
func (t *TimestampToken) CandidatesForSigningCertificate() *spi.CandidatesForSigningCertificate {
	if t.candidatesForSigningCertificate == nil {
		t.candidatesForSigningCertificate = t.CertificateSource().CandidatesForSigningCertificate(nil)
	}
	return t.candidatesForSigningCertificate
}

// SignerInformation returns the signer information used from the CMS SignedData object.
// Port of getSignerInformation().
func (t *TimestampToken) SignerInformation() *cmscore.SignerInfo {
	return t.tsaSignerInfo
}

// AtsHashIndexStatus gets the validation status of the ats-hash-index(-v3) attribute, when
// applicable, and nil otherwise.
// NOTE: supports only the archive-time-stamp-v3 timestamp type. Port of getAtsHashIndexStatus().
func (t *TimestampToken) AtsHashIndexStatus() *ArchiveTimestampHashIndexStatus {
	return t.hashIndexStatus
}

// SetAtsHashIndexStatus sets the validation status of the ats-hash-index(-v3) attribute, when
// applicable. Port of setAtsHashIndexStatus(ArchiveTimestampHashIndexStatus).
func (t *TimestampToken) SetAtsHashIndexStatus(hashIndexStatus *ArchiveTimestampHashIndexStatus) {
	t.hashIndexStatus = hashIndexStatus
}

// BuildTokenIdentifier builds the token's unique identifier.
// Port of the protected buildTokenIdentifier() override.
func (t *TimestampToken) BuildTokenIdentifier() *model.TokenIdentifier {
	return &t.TimestampIdentifierBuilder().BuildTimestampTokenIdentifier().TokenIdentifier
}

// TimestampIdentifierBuilder returns the TimestampIdentifierBuilder of the token, building the
// format-independent one on first use. Port of the protected getTimestampIdentifierBuilder().
func (t *TimestampToken) TimestampIdentifierBuilder() TimestampTokenIdentifierBuilder {
	if t.identifierBuilder == nil {
		t.identifierBuilder = NewTimestampIdentifierBuilder(t.Encoded()).SetFilename(t.filename)
	}
	return t.identifierBuilder
}

// Digest returns the digest value of the token's DER encoding for the requested algorithm.
// Port of the getDigest(DigestAlgorithm) override, i.e. DSSUtils.digest(digestAlgorithm,
// getEncoded()); Java's DSSException for an unavailable algorithm becomes the returned error.
func (t *TimestampToken) Digest(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	return spi.DSSUtilsDigest(digestAlgorithm, t.Encoded())
}

// compile-time assertions: a TimestampToken is a model.Token and supplies the operations
// model.TokenBase dispatches back into.
var (
	_ model.Token          = (*TimestampToken)(nil)
	_ model.TokenOverrides = (*TimestampToken)(nil)
)

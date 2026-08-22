// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/KeyEntityTSPSource.java (DSS 6.5.RC1).
//
// # BouncyCastle replacements
//
// Upstream delegates the whole issuance to BouncyCastle: TimeStampRequestGenerator,
// SignerInfoGeneratorBuilder with a DefaultSignedAttributeTableGenerator, TimeStampTokenGenerator
// and TimeStampResponseGenerator. internal/cmscore builds the CMS SignedData, but it is
// deliberately parse-only for RFC 3161, so the TSTInfo encoding and the signed-attribute set a
// TSA has to produce are assembled here. They reproduce what BouncyCastle 1.84 emits for the
// same inputs, which testdata/gen/TspFixtures.java pins down:
//
//   - TSTInfo: version 1, the TSA policy, the echoed MessageImprint (an AlgorithmIdentifier
//     without parameters), the serial number and a seconds-resolution GeneralizedTime; no
//     accuracy, ordering, nonce, tsa or extensions, which upstream never sets.
//   - Signed attributes: content-type, signing-time (the production time, pinned by
//     KeyEntityTSPSource#getSignedAttributeGenerator), the RFC 6211 CMS algorithm protection
//     attribute and the message-digest, all added by DefaultSignedAttributeTableGenerator, plus
//     the signing-certificate-v2 attribute TimeStampTokenGenerator contributes - an ESSCertIDv2
//     over the message-imprint digest algorithm, without issuerSerial, since the three-argument
//     TimeStampTokenGenerator constructor passes isIssuerSerialIncluded=false.
//
// # Deviations
//
//   - The protected methods whose parameters are BouncyCastle types (createRequest,
//     generateResponse, initResponseGenerator, getSignedAttributeGenerator, buildResponse,
//     getTimestampBinary) have no Go counterpart: no DSS class overrides them - PKITSPSource,
//     the only subclass, only reuses the constructor - and they exist solely to shape the
//     BouncyCastle objects. The three protected hooks that carry no BouncyCastle type
//     (getProductionTime, getSignatureAlgorithm, getTimeStampSerialNumber) stay overridable
//     through KeyEntityTSPSourceOverrides.
//   - The rejection branch of TimeStampResponseGenerator is unreachable: upstream builds the
//     request itself, from a digest algorithm it has already checked against
//     acceptedDigestAlgorithms, so the generator never refuses it. Only the accepted-algorithm
//     check survives, with its own DSSException message.
//   - java.security.KeyStore has no Go counterpart. The key-store constructors read a PKCS#12
//     store through golang.org/x/crypto/pkcs12 and refuse any other store type with a
//     documented error, the same treatment JKSSignatureToken gets in the dss-token port.
//   - The certificates field of the produced SignedData is a DER SET OF, i.e. its members are
//     ordered by their encoding, because that is how cmscore encodes every SET OF. BouncyCastle
//     emits the certificate chain in insertion order instead. Both are valid CMS; the tokens
//     differ byte for byte only when the chain holds more than one certificate.
package validation

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/internal/eccurve"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"

	"golang.org/x/crypto/pkcs12"
)

// keyEntityTSPSourceCMSAlgorithmProtect is
// id-aa-CMSAlgorithmProtection OBJECT IDENTIFIER ::= { iso(1) member-body(2) us(840)
// rsadsi(113549) pkcs(1) pkcs9(9) 52 }, the RFC 6211 attribute BouncyCastle's
// DefaultSignedAttributeTableGenerator adds to every signature it builds.
// Port of CMSAttributes.cmsAlgorithmProtect.
var keyEntityTSPSourceCMSAlgorithmProtect = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 52}

// keyEntityTSPSourceOIDSHA256 is id-sha256, the DEFAULT hashAlgorithm of an ESSCertIDv2, which
// DER omits.
var keyEntityTSPSourceOIDSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

// keyEntityTSPSourceHashes maps the digest algorithms this source can sign or digest with onto
// their crypto.Hash counterparts. An algorithm absent from the table cannot be used, which is
// the Go counterpart of the JCA refusing to provide a MessageDigest for it.
var keyEntityTSPSourceHashes = map[enumerations.DigestAlgorithm]crypto.Hash{
	enumerations.DigestAlgorithmSHA1:   crypto.SHA1,
	enumerations.DigestAlgorithmSHA224: crypto.SHA224,
	enumerations.DigestAlgorithmSHA256: crypto.SHA256,
	enumerations.DigestAlgorithmSHA384: crypto.SHA384,
	enumerations.DigestAlgorithmSHA512: crypto.SHA512,
}

// KeyEntityTSPSourceOverrides is the contract a subclass of KeyEntityTSPSource may implement.
// Port of the three protected methods that carry no BouncyCastle type; dispatched the way
// model.TokenBase dispatches to model.TokenOverrides via InitToken.
type KeyEntityTSPSourceOverrides interface {
	// ProductionTime gets the production time of the time-stamp.
	// Port of the protected getProductionTime().
	ProductionTime() time.Time
	// SignatureAlgorithm returns the target signature algorithm to be used for the time-stamp
	// generation. Port of the protected getSignatureAlgorithm(), whose IllegalArgumentException
	// becomes the returned error.
	SignatureAlgorithm() (enumerations.SignatureAlgorithm, error)
	// TimeStampSerialNumber generates the serial number of the produced timestamp token.
	// Port of the protected getTimeStampSerialNumber().
	TimeStampSerialNumber() (*big.Int, error)
}

// KeyEntityTSPSource is a TSPSource implementation allowing the issuance of a time-stamp using a
// local key entry.
type KeyEntityTSPSource struct {
	// overrides points back at the concrete source; see InitKeyEntityTSPSource. It is nil for a
	// bare KeyEntityTSPSource, which then uses this type's own implementations.
	overrides KeyEntityTSPSourceOverrides

	// signer is the private key to be used to sign a time-stamp token. Java holds a
	// java.security.PrivateKey and hands it to the JCA; Go signs through crypto.Signer, the
	// interface a key usable for signing satisfies.
	signer crypto.Signer

	// certificate is the certificate representing the time-stamp issuer.
	certificate *x509.Certificate

	// certificateChain is the certificate chain associated with the certificate.
	certificateChain []*x509.Certificate

	// acceptedDigestAlgorithms is the collection of digest algorithms accepted by the current
	// TSP source in the request.
	acceptedDigestAlgorithms []enumerations.DigestAlgorithm

	// tsaPolicy is the TSA policy.
	tsaPolicy string

	// productionTime is the static production date of the timestamp. Java declares the field
	// protected; ProductionTime and SetProductionTime are its Go accessors.
	productionTime time.Time

	// digestAlgorithm is the digest algorithm of the signature of the created time-stamp token.
	digestAlgorithm enumerations.DigestAlgorithm

	// encryptionAlgorithm is the encryption algorithm of the signature of the time-stamp.
	encryptionAlgorithm enumerations.EncryptionAlgorithm
}

// keyEntityTSPSourceDefaults returns the state Java's field initialisers set: the accepted digest
// algorithms and the SHA-512 signature digest algorithm.
func keyEntityTSPSourceDefaults() KeyEntityTSPSource {
	return KeyEntityTSPSource{
		acceptedDigestAlgorithms: []enumerations.DigestAlgorithm{
			enumerations.DigestAlgorithmSHA224, enumerations.DigestAlgorithmSHA256,
			enumerations.DigestAlgorithmSHA384, enumerations.DigestAlgorithmSHA512,
		},
		digestAlgorithm: enumerations.DigestAlgorithmSHA512,
	}
}

// NewKeyEntityTSPSourceBase instantiates the empty configuration a subclass embeds. Port of the
// protected KeyEntityTSPSource() constructor; the subclass constructor must follow it with
// InitKeyEntityTSPSource.
func NewKeyEntityTSPSourceBase() KeyEntityTSPSource {
	return keyEntityTSPSourceDefaults()
}

// InitKeyEntityTSPSource registers the concrete source with its base so that the base can
// dispatch the three overridable hooks. Every subclass constructor must call this once.
func (s *KeyEntityTSPSource) InitKeyEntityTSPSource(overrides KeyEntityTSPSourceOverrides) {
	s.overrides = overrides
}

// NewKeyEntityTSPSource instantiates a KeyEntityTSPSource with the given signing key, the
// corresponding time-stamp issuer certificate and the certificate chain to embed in the
// time-stamp. Port of the KeyEntityTSPSource(PrivateKey, X509Certificate, List<X509Certificate>)
// constructor.
//
// Panics with the Java messages when an argument is missing (Objects.requireNonNull).
func NewKeyEntityTSPSource(signer crypto.Signer, certificate *x509.Certificate,
	certificateChain []*x509.Certificate) *KeyEntityTSPSource {
	if signer == nil {
		panic("PrivateKey is not defined!")
	}
	if certificate == nil {
		panic("Certificate is not defined!")
	}
	if certificateChain == nil {
		panic("Certificate chain is not defined!")
	}
	source := keyEntityTSPSourceDefaults()
	source.signer = signer
	source.certificate = certificate
	source.certificateChain = certificateChain
	return &source
}

// NewKeyEntityTSPSourceFromTokens instantiates a KeyEntityTSPSource with the given signing key,
// the corresponding CertificateToken and the certificate chain to embed in the time-stamp.
// Port of the KeyEntityTSPSource(PrivateKey, CertificateToken, List<CertificateToken>)
// constructor.
func NewKeyEntityTSPSourceFromTokens(signer crypto.Signer, certificateToken *model.CertificateToken,
	certificateChain []*model.CertificateToken) *KeyEntityTSPSource {
	certificates := make([]*x509.Certificate, 0, len(certificateChain))
	for _, token := range certificateChain {
		certificates = append(certificates, token.Certificate())
	}
	var certificate *x509.Certificate
	if certificateToken != nil {
		certificate = certificateToken.Certificate()
	}
	return NewKeyEntityTSPSource(signer, certificate, certificates)
}

// NewKeyEntityTSPSourceFromKeyStore instantiates a KeyEntityTSPSource from the content of a key
// store. Port of the KeyEntityTSPSource(byte[], String, char[], String, char[]) constructor.
//
// ksType names the store format: only "PKCS12" is supported, java.security.KeyStore having no Go
// counterpart. Java's char[] passwords become strings, which is what the PKCS#12 reader takes.
// The alias selects the key entry; a PKCS#12 store read through golang.org/x/crypto/pkcs12
// exposes no friendly names, so only the empty alias - "the single key entry of the store" - is
// accepted.
func NewKeyEntityTSPSourceFromKeyStore(ksContent []byte, ksType string, ksPassword string,
	alias string, keyEntryPassword string) (*KeyEntityTSPSource, error) {
	if ksContent == nil {
		panic("KeyStore is not defined!")
	}
	if ksType != "PKCS12" && ksType != "pkcs12" {
		return nil, model.NewDSSError(fmt.Sprintf(
			"Unable to instantiate KeyStore: the key store type '%s' is not supported in the Go port, only PKCS12 is",
			ksType))
	}
	if alias != "" {
		return nil, model.NewDSSError(fmt.Sprintf(
			"Unable to recover the key entry with alias '%s'. Reason : a PKCS#12 key store carries no alias in the Go port",
			alias))
	}
	if keyEntryPassword != ksPassword {
		return nil, model.NewDSSError(
			"Unable to recover the key entry. Reason : a PKCS#12 key store protects its key entry with the store password")
	}
	// ToPEM, rather than Decode, because a time-stamp embeds the whole chain and Decode hands
	// out a single certificate. The blocks come out in the order the store holds them.
	blocks, err := pkcs12.ToPEM(ksContent, ksPassword)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to instantiate KeyStore", err)
	}
	var signer crypto.Signer
	var chain []*x509.Certificate
	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			certificate, err := eccurve.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, model.NewDSSErrorMessageCause("Unable to instantiate KeyStore", err)
			}
			chain = append(chain, certificate)
		case "PRIVATE KEY":
			key, err := keyEntityTSPSourcePrivateKey(block.Bytes)
			if err != nil {
				return nil, model.NewDSSErrorMessageCause(fmt.Sprintf(
					"Unable to recover the key entry with alias '%s'", alias), err)
			}
			signer = key
		}
	}
	if signer == nil || len(chain) == 0 {
		return nil, model.NewDSSError(fmt.Sprintf(
			"No related/supported key entry found for alias '%s'!", alias))
	}
	// x/crypto/pkcs12 does not tell which certificate the key belongs to; the store lists the
	// key entry's own certificate first, which is also the order a time-stamp chain needs.
	return NewKeyEntityTSPSource(signer, chain[0], chain), nil
}

// keyEntityTSPSourcePrivateKey reads the key of a PKCS#12 "PRIVATE KEY" PEM block, which
// x/crypto/pkcs12 writes in the algorithm's own format rather than PKCS#8.
func keyEntityTSPSourcePrivateKey(der []byte) (crypto.Signer, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("the key entry cannot sign")
	}
	return signer, nil
}

// NewKeyEntityTSPSourceFromKeyStorePath instantiates a KeyEntityTSPSource from the key store file
// at the given path. Port of the KeyEntityTSPSource(String, String, char[], String, char[]) and
// KeyEntityTSPSource(File, ...) constructors, whose IOException becomes the returned error.
func NewKeyEntityTSPSourceFromKeyStorePath(ksPath string, ksType string, ksPassword string,
	alias string, keyEntryPassword string) (*KeyEntityTSPSource, error) {
	content, err := os.ReadFile(ksPath)
	if err != nil {
		return nil, err
	}
	return NewKeyEntityTSPSourceFromKeyStore(content, ksType, ksPassword, alias, keyEntryPassword)
}

// SetPrivateKey sets the private key used to sign the time-stamp token.
// Port of setPrivateKey(PrivateKey).
func (s *KeyEntityTSPSource) SetPrivateKey(signer crypto.Signer) {
	s.signer = signer
}

// SetCertificate sets the time-stamp issuer certificate. Port of setCertificate(X509Certificate).
func (s *KeyEntityTSPSource) SetCertificate(certificate *x509.Certificate) {
	s.certificate = certificate
}

// SetCertificateChain sets the certificate chain to be embedded within the time-stamp token.
// Port of setCertificateChain(List<X509Certificate>).
func (s *KeyEntityTSPSource) SetCertificateChain(certificateChain []*x509.Certificate) {
	s.certificateChain = certificateChain
}

// SetTsaPolicy sets the TSA policy.
// NOTE: the property is mandatory for TimeStampToken generation. Port of setTsaPolicy(String).
func (s *KeyEntityTSPSource) SetTsaPolicy(tsaPolicy string) {
	s.tsaPolicy = tsaPolicy
}

// SetAcceptedDigestAlgorithms sets the digest algorithms to be accepted within a timestamp
// request. Default: SHA-224, SHA-256, SHA-384, SHA-512.
// Port of setAcceptedDigestAlgorithms(Collection<DigestAlgorithm>).
func (s *KeyEntityTSPSource) SetAcceptedDigestAlgorithms(digestAlgorithms []enumerations.DigestAlgorithm) {
	s.acceptedDigestAlgorithms = digestAlgorithms
}

// ProductionTime gets the production time of the time-stamp, defaulting to the current time.
// Port of the protected getProductionTime().
func (s *KeyEntityTSPSource) ProductionTime() time.Time {
	if s.productionTime.IsZero() {
		return time.Now()
	}
	return s.productionTime
}

// SetProductionTime sets the production time of the timestamp.
// NOTE: if not defined, the current time will be used. Port of setProductionTime(Date).
func (s *KeyEntityTSPSource) SetProductionTime(productionTime time.Time) {
	s.productionTime = productionTime
}

// SetDigestAlgorithm sets the digest algorithm of the signature of the generated time-stamp
// token. Default: DigestAlgorithm.SHA512. Port of setDigestAlgorithm(DigestAlgorithm).
func (s *KeyEntityTSPSource) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	s.digestAlgorithm = digestAlgorithm
}

// SetEncryptionAlgorithm sets the encryption algorithm to be used for the time-stamp signature
// generation.
// NOTE: the encryptionAlgorithm, when defined, shall be compatible with the encryption algorithm
// used by the target key. Port of setEncryptionAlgorithm(EncryptionAlgorithm).
func (s *KeyEntityTSPSource) SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm) {
	s.encryptionAlgorithm = encryptionAlgorithm
}

// SignatureAlgorithm returns the target signature algorithm to be used for the time-stamp
// generation. Port of the protected getSignatureAlgorithm(), whose IllegalArgumentException for
// an encryption algorithm that does not match the key becomes the returned error.
func (s *KeyEntityTSPSource) SignatureAlgorithm() (enumerations.SignatureAlgorithm, error) {
	keyAlgorithm, err := keyEntityTSPSourceEncryptionAlgorithmForKey(s.signer)
	if err != nil {
		return "", err
	}
	if s.encryptionAlgorithm != "" {
		if !s.encryptionAlgorithm.IsEquivalent(keyAlgorithm) {
			return "", fmt.Errorf(
				"Defined EncryptionAlgorithm '%s' is not equivalent to the one returned by time-stamp issuer '%s'",
				s.encryptionAlgorithm, keyAlgorithm)
		}
		keyAlgorithm = s.encryptionAlgorithm
	}
	return enumerations.SignatureAlgorithmGetAlgorithm(keyAlgorithm, s.digestAlgorithm), nil
}

// TimeStampSerialNumber generates the serial number of the produced timestamp token: a random
// 128-bit non-negative integer. Port of the protected getTimeStampSerialNumber(), i.e.
// new BigInteger(128, secureRandom).
func (s *KeyEntityTSPSource) TimeStampSerialNumber() (*big.Int, error) {
	// BigInteger(numBits, Random) draws a uniformly distributed value in [0, 2^128).
	return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
}

// keyEntityTSPSourceOverrides returns the registered overrides, falling back to this type's own
// implementations for a source that was not subclassed.
func (s *KeyEntityTSPSource) keyEntityTSPSourceOverrides() KeyEntityTSPSourceOverrides {
	if s.overrides == nil {
		return s
	}
	return s.overrides
}

// TimeStampResponse issues a time-stamp token over the given digest.
// Port of getTimeStampResponse(DigestAlgorithm, byte[]).
//
// Panics with the Java messages when a mandatory property is missing
// (Objects.requireNonNull); the DSSExceptions become returned errors.
func (s *KeyEntityTSPSource) TimeStampResponse(digestAlgorithm enumerations.DigestAlgorithm,
	digest []byte) (*model.TimestampBinary, error) {
	if s.signer == nil {
		panic("PrivateKey is not defined! Use #setPrivateKey method.")
	}
	if s.certificate == nil {
		panic("Certificate is not defined! Use #setCertificate method.")
	}
	if s.certificateChain == nil {
		panic("Certificate chain is not defined! Use #setCertificateChain method.")
	}
	if digestAlgorithm == "" {
		panic("DigestAlgorithm is not defined!")
	}
	if digest == nil {
		panic("digest is not defined!")
	}
	if s.tsaPolicy == "" {
		panic("TSAPolicy OID is not defined! Use #setTsaPolicy method.")
	}

	if !keyEntityTSPSourceContains(s.acceptedDigestAlgorithms, digestAlgorithm) {
		return nil, model.NewDSSError(fmt.Sprintf(
			"DigestAlgorithm '%s' is not supported by the KeyEntityTSPSource implementation!", digestAlgorithm))
	}

	token, err := s.generate(digestAlgorithm, digest)
	if err != nil {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to generate a timestamp. Reason : %s", err.Error()))
	}
	return model.NewTimestampBinary(token), nil
}

// generate assembles and signs the time-stamp token, i.e. what upstream reaches through
// createRequest, generateResponse and getTimestampBinary.
func (s *KeyEntityTSPSource) generate(messageImprintAlgorithm enumerations.DigestAlgorithm,
	digest []byte) ([]byte, error) {
	overrides := s.keyEntityTSPSourceOverrides()
	genTime := overrides.ProductionTime()
	signatureAlgorithm, err := overrides.SignatureAlgorithm()
	if err != nil {
		return nil, err
	}
	serialNumber, err := overrides.TimeStampSerialNumber()
	if err != nil {
		return nil, err
	}

	policy, err := keyEntityTSPSourceParseOID(s.tsaPolicy)
	if err != nil {
		return nil, err
	}
	messageImprintOID, err := keyEntityTSPSourceParseOID(messageImprintAlgorithm.OID())
	if err != nil {
		return nil, err
	}
	tstInfo := keyEntityTSPSourceTSTInfo(policy, messageImprintOID, digest, serialNumber, genTime)

	signedAttributes, err := s.signedAttributes(tstInfo, messageImprintAlgorithm, signatureAlgorithm, genTime)
	if err != nil {
		return nil, err
	}
	signerInfoBuilder := &cmscore.SignerInfoBuilder{
		SID: cmscore.NewIssuerAndSerialNumberSID(s.certificate.RawIssuer, s.certificate.SerialNumber),
		DigestAlgorithm: asn1ber.NewAlgorithmIdentifier(
			mustParseOID(signatureAlgorithm.DigestAlgorithm().OID())),
		SignedAttributes:   signedAttributes,
		SignatureAlgorithm: keyEntityTSPSourceSignatureAlgorithmIdentifier(signatureAlgorithm),
	}
	signature, err := keyEntityTSPSourceSign(s.signer, signatureAlgorithm, signerInfoBuilder.SignedAttributesDER())
	if err != nil {
		return nil, err
	}
	signerInfoBuilder.Signature = signature
	signerInfo, err := signerInfoBuilder.Build()
	if err != nil {
		return nil, err
	}

	certificates := make([]cmscore.CertificateChoice, 0, len(s.certificateChain))
	for _, certificate := range s.certificateChain {
		certificates = append(certificates, cmscore.NewCertificateChoice(certificate.Raw))
	}
	signedDataBuilder := &cmscore.SignedDataBuilder{
		EncapContentInfo: cmscore.NewEncapsulatedContentInfo(cmscore.OIDCTTSTInfo, tstInfo),
		Certificates:     certificates,
		SignerInfos:      []*cmscore.SignerInfo{signerInfo},
	}
	cms, err := signedDataBuilder.BuildCMS()
	if err != nil {
		return nil, err
	}
	return cms.DEREncoded(), nil
}

// signedAttributes builds the signed attributes of the TSA signature: the four
// DefaultSignedAttributeTableGenerator ones and the signing-certificate-v2 attribute
// TimeStampTokenGenerator adds.
func (s *KeyEntityTSPSource) signedAttributes(tstInfo []byte,
	messageImprintAlgorithm enumerations.DigestAlgorithm,
	signatureAlgorithm enumerations.SignatureAlgorithm, genTime time.Time) (cmscore.Attributes, error) {
	messageDigest, err := spi.DSSUtilsDigest(signatureAlgorithm.DigestAlgorithm(), tstInfo)
	if err != nil {
		return nil, err
	}
	certificateHash, err := spi.DSSUtilsDigest(messageImprintAlgorithm, s.certificate.Raw)
	if err != nil {
		return nil, err
	}
	messageImprintOID, err := keyEntityTSPSourceParseOID(messageImprintAlgorithm.OID())
	if err != nil {
		return nil, err
	}

	return cmscore.Attributes{
		cmscore.NewAttribute(cmscore.OIDContentType, asn1ber.EncodeOID(cmscore.OIDCTTSTInfo)),
		cmscore.NewAttribute(cmscore.OIDSigningTime, keyEntityTSPSourceTime(genTime)),
		cmscore.NewAttribute(keyEntityTSPSourceCMSAlgorithmProtect,
			keyEntityTSPSourceAlgorithmProtection(signatureAlgorithm)),
		cmscore.NewAttribute(spi.OIDIdAaSigningCertificateV2,
			keyEntityTSPSourceSigningCertificateV2(messageImprintOID, certificateHash)),
		cmscore.NewAttribute(cmscore.OIDMessageDigest, asn1ber.WriteTLV(asn1ber.TagOctetString, messageDigest)),
	}, nil
}

// keyEntityTSPSourceTSTInfo encodes the TSTInfo a TSA signs, with the fields upstream never sets
// left out.
func keyEntityTSPSourceTSTInfo(policy, messageImprintAlgorithm asn1.ObjectIdentifier, digest []byte,
	serialNumber *big.Int, genTime time.Time) []byte {
	messageImprint := asn1ber.WriteSequence(asn1ber.EncodeOID(messageImprintAlgorithm))
	messageImprint = append(messageImprint, asn1ber.WriteTLV(asn1ber.TagOctetString, digest)...)

	body := asn1ber.EncodeInteger(big.NewInt(1))
	body = append(body, asn1ber.EncodeOID(policy)...)
	body = append(body, asn1ber.WriteSequence(messageImprint)...)
	body = append(body, asn1ber.EncodeInteger(serialNumber)...)
	// TimeStampTokenGenerator#createGeneralizedTime writes the production time with the
	// resolution it was configured with, seconds by default, and no fractional part.
	body = append(body, asn1ber.WriteTLV(asn1ber.TagGeneralizedTime,
		[]byte(genTime.UTC().Format("20060102150405")+"Z"))...)
	return asn1ber.WriteSequence(body)
}

// keyEntityTSPSourceTime encodes a CMS Time, i.e. the UTCTime alternative BouncyCastle's
// org.bouncycastle.asn1.cms.Time picks for a year between 1950 and 2049, and the GeneralizedTime
// alternative outside that range.
func keyEntityTSPSourceTime(value time.Time) []byte {
	utc := value.UTC()
	if year := utc.Year(); year < 1950 || year > 2049 {
		return asn1ber.WriteTLV(asn1ber.TagGeneralizedTime, []byte(utc.Format("20060102150405")+"Z"))
	}
	return asn1ber.WriteTLV(asn1ber.TagUTCTime, []byte(utc.Format("060102150405")+"Z"))
}

// keyEntityTSPSourceAlgorithmProtection encodes the RFC 6211 CMSAlgorithmProtection value:
//
//	CMSAlgorithmProtection ::= SEQUENCE {
//	    digestAlgorithm        DigestAlgorithmIdentifier,
//	    signatureAlgorithm [1] SignatureAlgorithmIdentifier OPTIONAL,
//	    macAlgorithm       [2] MessageAuthenticationCodeAlgorithm OPTIONAL }
func keyEntityTSPSourceAlgorithmProtection(signatureAlgorithm enumerations.SignatureAlgorithm) []byte {
	body := asn1ber.WriteSequence(asn1ber.EncodeOID(mustParseOID(signatureAlgorithm.DigestAlgorithm().OID())))
	identifier := keyEntityTSPSourceSignatureAlgorithmIdentifier(signatureAlgorithm)
	// The signatureAlgorithm field is [1] IMPLICIT, i.e. the SEQUENCE identifier octet is
	// replaced by the context-specific one.
	signatureAlgorithmContent := asn1ber.EncodeOID(identifier.Algorithm)
	signatureAlgorithmContent = append(signatureAlgorithmContent, identifier.Parameters...)
	body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1,
		signatureAlgorithmContent)...)
	return asn1ber.WriteSequence(body)
}

// keyEntityTSPSourceSigningCertificateV2 encodes the RFC 5035 SigningCertificateV2 value the
// three-argument TimeStampTokenGenerator constructor contributes: one ESSCertIDv2 holding the
// certificate hash, with the DEFAULT hashAlgorithm omitted and no issuerSerial.
func keyEntityTSPSourceSigningCertificateV2(hashAlgorithm asn1.ObjectIdentifier, certificateHash []byte) []byte {
	var essCertID []byte
	if !hashAlgorithm.Equal(keyEntityTSPSourceOIDSHA256) {
		essCertID = asn1ber.WriteSequence(asn1ber.EncodeOID(hashAlgorithm))
	}
	essCertID = append(essCertID, asn1ber.WriteTLV(asn1ber.TagOctetString, certificateHash)...)
	return asn1ber.WriteSequence(asn1ber.WriteSequence(asn1ber.WriteSequence(essCertID)))
}

// keyEntityTSPSourceSignatureAlgorithmIdentifier returns the AlgorithmIdentifier of a signature
// algorithm, i.e. what BouncyCastle's DefaultSignatureAlgorithmIdentifierFinder produces: the
// RSA algorithms carry an explicit NULL, the others no parameters at all.
//
// RSASSA-PSS is deliberately absent: its parameters are a structure of their own, and upstream
// only reaches it when a caller sets the encryption algorithm explicitly. SignatureAlgorithm
// refuses it before this point.
func keyEntityTSPSourceSignatureAlgorithmIdentifier(
	signatureAlgorithm enumerations.SignatureAlgorithm) *asn1ber.AlgorithmIdentifier {
	oid := mustParseOID(signatureAlgorithm.OID())
	if signatureAlgorithm.EncryptionAlgorithm() == enumerations.EncryptionAlgorithmRSA {
		return asn1ber.NewAlgorithmIdentifierWithParameters(oid, asn1ber.DERNull)
	}
	return asn1ber.NewAlgorithmIdentifier(oid)
}

// keyEntityTSPSourceSign signs the given bytes with the source's key, applying the digest and
// padding the signature algorithm calls for. It stands in for the JcaContentSignerBuilder of
// upstream's initResponseGenerator.
func keyEntityTSPSourceSign(signer crypto.Signer, signatureAlgorithm enumerations.SignatureAlgorithm,
	toBeSigned []byte) ([]byte, error) {
	if signatureAlgorithm.EncryptionAlgorithm() == enumerations.EncryptionAlgorithmEDDSA {
		// EdDSA is a pure signature scheme: the whole message is signed, not a digest of it.
		return signer.Sign(rand.Reader, toBeSigned, crypto.Hash(0))
	}
	hash, supported := keyEntityTSPSourceHashes[signatureAlgorithm.DigestAlgorithm()]
	if !supported {
		return nil, fmt.Errorf("NoSuchAlgorithmException : %s cannot be used to sign a time-stamp",
			signatureAlgorithm.DigestAlgorithm())
	}
	hasher := hash.New()
	hasher.Write(toBeSigned)
	return signer.Sign(rand.Reader, hasher.Sum(nil), hash)
}

// keyEntityTSPSourceEncryptionAlgorithmForKey returns the encryption algorithm of the signing key.
// Port of EncryptionAlgorithm.forKey(Key), which reads java.security.Key#getAlgorithm().
//
// RSASSA-PSS is not offered: a Go RSA key does not record the padding it is meant to be used
// with, so it always reports RSA, exactly as a JCA RSA key does.
func keyEntityTSPSourceEncryptionAlgorithmForKey(signer crypto.Signer) (enumerations.EncryptionAlgorithm, error) {
	if signer == nil {
		return "", fmt.Errorf("unsupported algorithm: the signing key is not defined")
	}
	switch signer.Public().(type) {
	case *rsa.PublicKey:
		return enumerations.EncryptionAlgorithmRSA, nil
	case *ecdsa.PublicKey:
		return enumerations.EncryptionAlgorithmECDSA, nil
	case ed25519.PublicKey:
		return enumerations.EncryptionAlgorithmEDDSA, nil
	}
	return "", fmt.Errorf("unsupported algorithm: %T", signer.Public())
}

// keyEntityTSPSourceContains reports whether the collection holds the given digest algorithm.
func keyEntityTSPSourceContains(digestAlgorithms []enumerations.DigestAlgorithm,
	digestAlgorithm enumerations.DigestAlgorithm) bool {
	for _, candidate := range digestAlgorithms {
		if candidate == digestAlgorithm {
			return true
		}
	}
	return false
}

// keyEntityTSPSourceParseOID reads a dotted OID, the way `new ASN1ObjectIdentifier(oid)` does.
func keyEntityTSPSourceParseOID(oid string) (asn1.ObjectIdentifier, error) {
	if oid == "" {
		return nil, fmt.Errorf("IllegalArgumentException : string  not an OID")
	}
	var identifier asn1.ObjectIdentifier
	for _, arc := range splitOIDArcs(oid) {
		// encoding/asn1 stores an arc in an int, so an arc that does not fit one is refused -
		// the same restriction internal/cmscore documents for the OIDs it reads.
		value, err := strconv.Atoi(arc)
		if err != nil || value < 0 {
			return nil, fmt.Errorf("IllegalArgumentException : string %s not an OID", oid)
		}
		identifier = append(identifier, value)
	}
	if len(identifier) < 2 {
		return nil, fmt.Errorf("IllegalArgumentException : string %s not an OID", oid)
	}
	return identifier, nil
}

// splitOIDArcs splits a dotted OID into its arcs.
func splitOIDArcs(oid string) []string {
	var arcs []string
	start := 0
	for index := 0; index <= len(oid); index++ {
		if index == len(oid) || oid[index] == '.' {
			arcs = append(arcs, oid[start:index])
			start = index + 1
		}
	}
	return arcs
}

// mustParseOID reads an OID that comes from an enumeration and can therefore only be well formed.
func mustParseOID(oid string) asn1.ObjectIdentifier {
	identifier, err := keyEntityTSPSourceParseOID(oid)
	if err != nil {
		panic(err.Error())
	}
	return identifier
}

// compile-time assertions: a KeyEntityTSPSource is a TSPSource and supplies its own overridable
// hooks.
var (
	_ TSPSource                   = (*KeyEntityTSPSource)(nil)
	_ KeyEntityTSPSourceOverrides = (*KeyEntityTSPSource)(nil)
)

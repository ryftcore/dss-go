package pfx

import (
	"crypto"
	"crypto/hmac"
	"encoding/asn1"
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// Certificate is one certBag SafeBag extracted from a PFX PDU.
type Certificate struct {
	// Raw is the DER encoding of the X.509 certificate.
	Raw []byte
	// LocalKeyID is the bag's localKeyId attribute, nil when absent. A certificate and the
	// (shrouded) private key it belongs to share the same LocalKeyID.
	LocalKeyID []byte
	// FriendlyName is the bag's friendlyName attribute, "" when absent.
	FriendlyName string
}

// PrivateKey is one keyBag or pkcs8ShroudedKeyBag SafeBag extracted from a PFX PDU, already
// decrypted (for a shrouded bag) and parsed.
type PrivateKey struct {
	// Key is the parsed private key: *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey or
	// *dsa.PrivateKey. It never implements crypto.Signer for the DSA case - crypto/dsa predates
	// that interface - so a caller that needs a crypto.Signer wraps Key itself.
	Key crypto.PrivateKey
	// LocalKeyID is the bag's localKeyId attribute, nil when absent.
	LocalKeyID []byte
	// FriendlyName is the bag's friendlyName attribute, "" when absent.
	FriendlyName string
}

// Store is the decrypted, parsed content of a PFX PDU, in encounter order. Correlating a
// PrivateKey with its Certificate chain (by LocalKeyID, then by issuer/subject linkage) is left
// to the caller - java.security.KeyStore's PKCS12 provider does the same internal bookkeeping,
// which this package does not reproduce since dss/token already implements it for the
// certificate-chain values this package now feeds it instead of golang.org/x/crypto/pkcs12's.
type Store struct {
	Certificates []Certificate
	PrivateKeys  []PrivateKey
}

// Load decrypts and parses a PKCS#12 (PFX) file, verifying its password-integrity MAC first.
// See this package's doc comment for the subset of RFC 7292 it supports.
func Load(der []byte, password string) (*Store, error) {
	root, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, fmt.Errorf("pfx: %w", err)
	}
	if len(rest) != 0 {
		return nil, errors.New("pfx: trailing data after the PFX PDU")
	}
	if !root.IsUniversal(asn1ber.TagSequence) || !root.IsConstructed() {
		return nil, errors.New("pfx: not a PFX SEQUENCE")
	}
	children := root.Children()
	if len(children) < 2 {
		return nil, fmt.Errorf("pfx: PFX holds %d components, at least 2 expected", len(children))
	}
	if !children[0].IsUniversal(asn1ber.TagInteger) || children[0].Integer().Int64() != 3 {
		return nil, errors.New("pfx: unsupported PFX version (only v3 PDUs are supported)")
	}
	if len(children) < 3 {
		return nil, errors.New("pfx: PFX carries no MacData (public-key/signed integrity PFX PDUs are not supported)")
	}

	authSafeContentType, authSafeContent, err := parseContentInfoHeader(children[1])
	if err != nil {
		return nil, err
	}
	if !authSafeContentType.Equal(oidData) || authSafeContent == nil {
		return nil, fmt.Errorf("pfx: unsupported PFX.authSafe content type %s (only id-data is supported)", authSafeContentType)
	}
	if !authSafeContent.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: PFX.authSafe.content is not an OCTET STRING")
	}
	authenticatedSafeDER := authSafeContent.Octets()

	md, err := parseMacData(children[2])
	if err != nil {
		return nil, err
	}
	bmpPassword, err := bmpStringPassword(password)
	if err != nil {
		return nil, err
	}
	if err := verifyMAC(md, authenticatedSafeDER, bmpPassword); err != nil {
		return nil, err
	}

	authenticatedSafe, rest, err := asn1ber.Parse(authenticatedSafeDER)
	if err != nil {
		return nil, fmt.Errorf("pfx: invalid AuthenticatedSafe: %w", err)
	}
	if len(rest) != 0 || !authenticatedSafe.IsUniversal(asn1ber.TagSequence) || !authenticatedSafe.IsConstructed() {
		return nil, errors.New("pfx: AuthenticatedSafe is not a SEQUENCE")
	}

	store := &Store{}
	for _, contentInfo := range authenticatedSafe.Children() {
		safeContentsDER, err := decodeAuthenticatedSafeEntry(contentInfo, bmpPassword, []byte(password))
		if err != nil {
			return nil, err
		}
		bags, err := parseSafeContents(safeContentsDER)
		if err != nil {
			return nil, err
		}
		if err := store.absorb(bags, bmpPassword, []byte(password)); err != nil {
			return nil, err
		}
	}
	if len(store.PrivateKeys) == 0 {
		return nil, errors.New("pfx: no private key entry found (wrong password, or an unsupported key type)")
	}
	return store, nil
}

// decodeAuthenticatedSafeEntry returns the DER of the SafeContents one AuthenticatedSafe
// ContentInfo carries, decrypting it first when its content type is id-encryptedData.
func decodeAuthenticatedSafeEntry(contentInfo *asn1ber.Element, bmpPassword, rawPassword []byte) ([]byte, error) {
	contentType, content, err := parseContentInfoHeader(contentInfo)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, errors.New("pfx: AuthenticatedSafe entry carries no content")
	}
	switch {
	case contentType.Equal(oidData):
		if !content.IsUniversal(asn1ber.TagOctetString) {
			return nil, errors.New("pfx: AuthenticatedSafe data entry is not an OCTET STRING")
		}
		return content.Octets(), nil
	case contentType.Equal(oidEncryptedData):
		return decryptEncryptedData(content, bmpPassword, rawPassword)
	default:
		return nil, fmt.Errorf("pfx: unsupported AuthenticatedSafe content type %s", contentType)
	}
}

// absorb decodes each SafeBag into store's Certificates/PrivateKeys, decrypting a shrouded key
// bag as it goes. Bag types this package does not use (secretBag, crlBag, ...) are ignored.
func (store *Store) absorb(bags []safeBag, bmpPassword, rawPassword []byte) error {
	for _, bag := range bags {
		switch {
		case bag.id.Equal(oidCertBag):
			raw, err := decodeCertBag(bag.value)
			if err != nil {
				return err
			}
			store.Certificates = append(store.Certificates, Certificate{
				Raw: raw, LocalKeyID: bag.localKeyID, FriendlyName: bag.friendlyName,
			})
		case bag.id.Equal(oidKeyBag):
			key, err := parsePrivateKeyInfo(bag.value.DEREncoded())
			if err != nil {
				return fmt.Errorf("pfx: keyBag: %w", err)
			}
			store.PrivateKeys = append(store.PrivateKeys, PrivateKey{
				Key: key, LocalKeyID: bag.localKeyID, FriendlyName: bag.friendlyName,
			})
		case bag.id.Equal(oidPKCS8ShroudedKeyBag):
			plaintext, err := decryptShroudedKeyBag(bag.value, bmpPassword, rawPassword)
			if err != nil {
				return fmt.Errorf("pfx: pkcs8ShroudedKeyBag: %w", err)
			}
			key, err := parsePrivateKeyInfo(plaintext)
			if err != nil {
				return fmt.Errorf("pfx: pkcs8ShroudedKeyBag: %w", err)
			}
			store.PrivateKeys = append(store.PrivateKeys, PrivateKey{
				Key: key, LocalKeyID: bag.localKeyID, FriendlyName: bag.friendlyName,
			})
		}
	}
	return nil
}

// decryptWithAlgorithm decrypts ciphertext under algorithm, dispatching to PBES2 (modern
// openssl/JDK default) or one of the RFC 7292 Appendix B legacy PBE schemes (openssl's
// "-legacy" mode, every JDK before the PBES2 default).
func decryptWithAlgorithm(algorithm *asn1ber.AlgorithmIdentifier, bmpPassword, rawPassword, ciphertext []byte) ([]byte, error) {
	if algorithm.Algorithm.Equal(oidPBES2) {
		return decryptPBES2(algorithm, rawPassword, ciphertext)
	}
	return decryptLegacyPBE(algorithm, bmpPassword, ciphertext)
}

// decodeCertBag decodes a
//
//	CertBag ::= SEQUENCE {
//	    certId    OBJECT IDENTIFIER,
//	    certValue [0] EXPLICIT ANY DEFINED BY certId }
//
// returning the DER of the wrapped certificate for the one certId RFC 7292 Appendix D and every
// producer this package reads use: x509Certificate.
func decodeCertBag(bagValue *asn1ber.Element) ([]byte, error) {
	if !bagValue.IsUniversal(asn1ber.TagSequence) || !bagValue.IsConstructed() || len(bagValue.Children()) != 2 {
		return nil, errors.New("pfx: CertBag is not a two-component SEQUENCE")
	}
	certID, err := bagValue.Children()[0].ObjectIdentifier()
	if err != nil {
		return nil, err
	}
	if !certID.Equal(oidCertTypeX509Certificate) {
		return nil, fmt.Errorf("pfx: unsupported CertBag certificate type %s (only X.509 is supported)", certID)
	}
	if !bagValue.Children()[1].IsContextSpecific(0) {
		return nil, errors.New("pfx: CertBag.certValue is not [0] EXPLICIT")
	}
	certValue, err := explicitContent(bagValue.Children()[1])
	if err != nil {
		return nil, err
	}
	if !certValue.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: CertBag.certValue does not wrap an OCTET STRING")
	}
	return certValue.Octets(), nil
}

// decryptShroudedKeyBag decrypts a
//
//	PKCS8ShroudedKeyBag ::= EncryptedPrivateKeyInfo ::= SEQUENCE {
//	    encryptionAlgorithm AlgorithmIdentifier,
//	    encryptedData       OCTET STRING }
//
// returning the DER of the PrivateKeyInfo it wraps.
func decryptShroudedKeyBag(bagValue *asn1ber.Element, bmpPassword, rawPassword []byte) ([]byte, error) {
	if !bagValue.IsUniversal(asn1ber.TagSequence) || !bagValue.IsConstructed() || len(bagValue.Children()) != 2 {
		return nil, errors.New("pfx: EncryptedPrivateKeyInfo is not a two-component SEQUENCE")
	}
	algorithm, err := asn1ber.AlgorithmIdentifierFromElement(bagValue.Children()[0])
	if err != nil {
		return nil, err
	}
	encryptedData := bagValue.Children()[1]
	if !encryptedData.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: EncryptedPrivateKeyInfo.encryptedData is not an OCTET STRING")
	}
	return decryptWithAlgorithm(algorithm, bmpPassword, rawPassword, encryptedData.Octets())
}

// decryptEncryptedData decrypts a
//
//	EncryptedData ::= SEQUENCE {
//	    version              INTEGER,
//	    encryptedContentInfo EncryptedContentInfo }
//	EncryptedContentInfo ::= SEQUENCE {
//	    contentType                OBJECT IDENTIFIER,
//	    contentEncryptionAlgorithm AlgorithmIdentifier,
//	    encryptedContent           [0] IMPLICIT OCTET STRING OPTIONAL }
//
// content is the already EXPLICIT-unwrapped EncryptedData value from the owning ContentInfo.
func decryptEncryptedData(content *asn1ber.Element, bmpPassword, rawPassword []byte) ([]byte, error) {
	if !content.IsUniversal(asn1ber.TagSequence) || !content.IsConstructed() || len(content.Children()) != 2 {
		return nil, errors.New("pfx: EncryptedData is not a two-component SEQUENCE")
	}
	version := content.Children()[0]
	if !version.IsUniversal(asn1ber.TagInteger) || version.Integer().Int64() != 0 {
		return nil, errors.New("pfx: unsupported EncryptedData version")
	}
	eci := content.Children()[1]
	if !eci.IsUniversal(asn1ber.TagSequence) || !eci.IsConstructed() || len(eci.Children()) < 3 {
		return nil, errors.New("pfx: EncryptedContentInfo is missing its encryptedContent")
	}
	algorithm, err := asn1ber.AlgorithmIdentifierFromElement(eci.Children()[1])
	if err != nil {
		return nil, err
	}
	encryptedContent := eci.Children()[2]
	if !encryptedContent.IsContextSpecific(0) {
		return nil, errors.New("pfx: EncryptedContentInfo.encryptedContent is not [0] IMPLICIT")
	}
	return decryptWithAlgorithm(algorithm, bmpPassword, rawPassword, encryptedContent.Octets())
}

// parseContentInfoHeader decodes a
//
//	ContentInfo ::= SEQUENCE {
//	    contentType OBJECT IDENTIFIER,
//	    content     [0] EXPLICIT ANY DEFINED BY contentType OPTIONAL }
//
// content is nil when the optional field is absent.
func parseContentInfoHeader(element *asn1ber.Element) (asn1.ObjectIdentifier, *asn1ber.Element, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, nil, errors.New("pfx: ContentInfo is not a SEQUENCE")
	}
	children := element.Children()
	if len(children) == 0 {
		return nil, nil, errors.New("pfx: ContentInfo holds no components")
	}
	contentType, err := children[0].ObjectIdentifier()
	if err != nil {
		return nil, nil, err
	}
	if len(children) < 2 {
		return contentType, nil, nil
	}
	if !children[1].IsContextSpecific(0) {
		return nil, nil, errors.New("pfx: ContentInfo.content is not [0] EXPLICIT")
	}
	content, err := explicitContent(children[1])
	if err != nil {
		return nil, nil, err
	}
	return contentType, content, nil
}

// macData is RFC 7292 section 4's MacData.
type macData struct {
	digestAlgorithm asn1.ObjectIdentifier
	digest          []byte
	salt            []byte
	iterations      int
}

// parseMacData decodes a
//
//	MacData ::= SEQUENCE {
//	    mac        DigestInfo,
//	    macSalt    OCTET STRING,
//	    iterations INTEGER DEFAULT 1 }
//	DigestInfo ::= SEQUENCE {
//	    digestAlgorithm AlgorithmIdentifier,
//	    digest          OCTET STRING }
func parseMacData(element *asn1ber.Element) (*macData, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() || len(element.Children()) < 2 {
		return nil, errors.New("pfx: MacData is not a SEQUENCE")
	}
	children := element.Children()
	digestInfo := children[0]
	if !digestInfo.IsUniversal(asn1ber.TagSequence) || !digestInfo.IsConstructed() || len(digestInfo.Children()) != 2 {
		return nil, errors.New("pfx: MacData.mac is not a DigestInfo SEQUENCE")
	}
	algorithm, err := asn1ber.AlgorithmIdentifierFromElement(digestInfo.Children()[0])
	if err != nil {
		return nil, err
	}
	digestValue := digestInfo.Children()[1]
	if !digestValue.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: MacData.mac.digest is not an OCTET STRING")
	}
	salt := children[1]
	if !salt.IsUniversal(asn1ber.TagOctetString) {
		return nil, errors.New("pfx: MacData.macSalt is not an OCTET STRING")
	}
	iterations := 1
	if len(children) > 2 {
		if !children[2].IsUniversal(asn1ber.TagInteger) {
			return nil, errors.New("pfx: MacData.iterations is not an INTEGER")
		}
		iterations = int(children[2].Integer().Int64())
	}
	return &macData{
		digestAlgorithm: algorithm.Algorithm,
		digest:          digestValue.Octets(),
		salt:            salt.Octets(),
		iterations:      iterations,
	}, nil
}

// verifyMAC recomputes the RFC 7292 section 4 integrity MAC over content (the AuthenticatedSafe
// DER) and checks it against md, in constant time.
func verifyMAC(md *macData, content, bmpPassword []byte) error {
	h := pfxHashByOID(md.digestAlgorithm)
	if h == nil {
		return fmt.Errorf("pfx: unsupported MacData digest algorithm %s", md.digestAlgorithm)
	}
	key := deriveKeyMaterial(*h, 3, md.salt, bmpPassword, md.iterations, h.u)
	mac := hmac.New(h.new, key)
	mac.Write(content)
	if !hmac.Equal(mac.Sum(nil), md.digest) {
		return errors.New("pfx: MAC verification failed (wrong password, or corrupt PFX)")
	}
	return nil
}

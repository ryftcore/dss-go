// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/KeyStoreCertificateSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
//
// DEVIATION: java.security.KeyStore is a generic, provider-backed abstraction covering many
// formats (JKS, PKCS12, ...) addressed by alias. Go has no equivalent generic keystore type.
// This port narrows KeyStoreCertificateSourceType to the two formats the brief approves:
// KeyStoreCertificateSourceType_PKCS12 (read-only: golang.org/x/crypto/pkcs12 exposes no
// encoder in the vendored version, so Store on a PKCS12-typed source returns an error) and
//
// LIMITATION (PKCS12 reading): golang.org/x/crypto/pkcs12 only understands the legacy
// PKCS#12 profile - a SHA-1 MAC with PBE-SHA1-3DES / PBE-SHA1-RC2-40 encryption. Files
// produced with the modern defaults (OpenSSL 3.x, `keytool` on recent JDKs: SHA-256 MAC,
// PBES2/AES-CBC encryption) are REJECTED with "pkcs12: unknown digest algorithm ..." or
// "pkcs12: unknown algorithm identifier ...", where java.security.KeyStore reads them all.
// Loading such a keystore therefore needs it re-exported with `openssl pkcs12 -legacy`
// (or an alternative PKCS#12 reader) until the dependency grows PBES2 support.
//
// KeyStoreCertificateSourceType_PEM (a plain concatenated PEM/DER certificate collection,
// which Java's KeyStore SPI has no counterpart for - it stands in for "PEM/DER cert
// collections" the brief calls for). JKS - and any other Java keystore type string - is
// rejected with a clear "unsupported" error; this is a documented deviation, flagged in the
// porter's notes for integrator review.
//
// Since neither format reliably carries Java-style KeyStore aliases (PKCS12 bags may carry a
// friendlyName attribute, but golang.org/x/crypto/pkcs12's ToPEM does not surface it; a PEM
// collection has none at all), certificates loaded from a stream are assigned synthetic
// aliases "cert-<n>" in encounter order. Certificates added at runtime via
// AddCertificateToKeyStore keep Java's own convention of using the certificate's DSSIdAsString
// as its alias.
package spi

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/pkcs12"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// KeyStoreCertificateSourceType identifies the keystore format a KeyStoreCertificateSource
// was built with. Port of the free-form Java String ksType, narrowed to the formats this Go
// port supports (see the file header for the documented deviation from Java's generic
// KeyStore SPI).
type KeyStoreCertificateSourceType string

const (
	// KeyStoreCertificateSourceType_PKCS12 is a PKCS#12 keystore ("PKCS12" in Java), read-only
	// in this port.
	KeyStoreCertificateSourceType_PKCS12 KeyStoreCertificateSourceType = "PKCS12"
	// KeyStoreCertificateSourceType_PEM is a plain concatenated PEM (or raw concatenated DER)
	// certificate collection; not a Java KeyStore type, added to satisfy the "PEM/DER cert
	// collections" requirement without a real generic keystore abstraction.
	KeyStoreCertificateSourceType_PEM KeyStoreCertificateSourceType = "PEM"
)

// KeyStoreCertificateSource implements a CertificateSource using a keystore (PKCS12 or a
// plain PEM/DER certificate collection; JKS is documented as unsupported, see the file header).
type KeyStoreCertificateSource struct {
	CommonCertificateSource

	// ksType is the keystore format.
	ksType KeyStoreCertificateSourceType

	// passwordProtection is the keystore password.
	passwordProtection []byte

	// entries maps an alias to the certificate token stored under it, mirroring Java's
	// alias-addressed KeyStore. Java's KeyStore aliases enumerate in provider-defined (but
	// stable-per-run) order; a bare Go map is randomized on every run instead, and Store()
	// writes certificates in this map's iteration order, so it is kept insertion-ordered
	// (slice + index map, PORTING.md's Collections rule) rather than a bare map.
	entries *utils.OrderedMap[string, *model.CertificateToken]
}

// NewKeyStoreCertificateSource creates a new, empty keystore of the given type.
// Port of the KeyStoreCertificateSource(String, char[]) constructor.
func NewKeyStoreCertificateSource(ksType KeyStoreCertificateSourceType, ksPassword []byte) (*KeyStoreCertificateSource, error) {
	return NewKeyStoreCertificateSourceFromReader(nil, ksType, ksPassword)
}

// NewKeyStoreCertificateSourceFromFilePath loads a keystore from the given file path.
// Port of the KeyStoreCertificateSource(String, String, char[]) constructor.
func NewKeyStoreCertificateSourceFromFilePath(ksFilePath string, ksType KeyStoreCertificateSourceType, ksPassword []byte) (*KeyStoreCertificateSource, error) {
	f, err := os.Open(ksFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewKeyStoreCertificateSourceFromReader(f, ksType, ksPassword)
}

// NewKeyStoreCertificateSourceFromReader loads a keystore from ksStream, or creates a new
// empty one when ksStream is nil. Port of the KeyStoreCertificateSource(InputStream, String,
// char[]) constructor, the "default" constructor upstream documents.
func NewKeyStoreCertificateSourceFromReader(ksStream io.Reader, ksType KeyStoreCertificateSourceType, ksPassword []byte) (*KeyStoreCertificateSource, error) {
	source := &KeyStoreCertificateSource{
		CommonCertificateSource: NewCommonCertificateSource(),
		ksType:                  ksType,
		passwordProtection:      ksPassword,
		entries:                 utils.NewOrderedMap[string, *model.CertificateToken](),
	}
	if err := source.initKeystore(ksStream); err != nil {
		return nil, err
	}
	return source, nil
}

// initKeystore ports the private initKeystore(InputStream, String, char[]).
func (k *KeyStoreCertificateSource) initKeystore(ksStream io.Reader) error {
	switch k.ksType {
	case KeyStoreCertificateSourceType_PKCS12, KeyStoreCertificateSourceType_PEM:
		// supported, handled below
	case "JKS":
		return model.NewDSSError("Unable to initialize the keystore: JKS keystores are not supported by this Go port (documented deviation; use PKCS12 or PEM)")
	default:
		return model.NewDSSError(fmt.Sprintf("Unable to initialize the keystore: unsupported keystore type %q", k.ksType))
	}

	if ksStream == nil {
		return nil
	}
	data, err := io.ReadAll(ksStream)
	if err != nil {
		return model.NewDSSErrorMessageCause("Unable to initialize the keystore", err)
	}
	if len(data) == 0 {
		return nil
	}

	var certificates []*model.CertificateToken
	switch k.ksType {
	case KeyStoreCertificateSourceType_PKCS12:
		certificates, err = keyStoreCertificateSourceParsePKCS12(data, string(k.passwordProtection))
	case KeyStoreCertificateSourceType_PEM:
		certificates, err = keyStoreCertificateSourceParsePEMOrDER(data)
	}
	if err != nil {
		return model.NewDSSErrorMessageCause("Unable to retrieve certificates from the keystore", err)
	}

	for i, certificateToken := range certificates {
		alias := k.getKey("cert-" + strconv.Itoa(i))
		k.entries.Set(alias, certificateToken)
		k.CommonCertificateSource.AddCertificate(certificateToken)
	}
	return nil
}

// keyStoreCertificateSourceParsePKCS12 extracts the certificates (not the private keys) from
// a PKCS#12 file. golang.org/x/crypto/pkcs12 (the vendored version) exposes no alias/friendly
// name, so certificates are returned in encounter order.
func keyStoreCertificateSourceParsePKCS12(data []byte, password string) ([]*model.CertificateToken, error) {
	blocks, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return nil, err
	}
	var result []*model.CertificateToken
	for _, block := range blocks {
		if block.Type != "CERTIFICATE" {
			continue
		}
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		certificateToken, err := model.NewCertificateToken(certificate)
		if err != nil {
			return nil, err
		}
		result = append(result, certificateToken)
	}
	return result, nil
}

// keyStoreCertificateSourceParsePEMOrDER extracts certificates from a byte sequence that is
// either concatenated PEM blocks or concatenated DER certificates.
func keyStoreCertificateSourceParsePEMOrDER(data []byte) ([]*model.CertificateToken, error) {
	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if bytes.HasPrefix(trimmed, []byte("-----BEGIN")) {
		var result []*model.CertificateToken
		rest := data
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				continue
			}
			certificate, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certificateToken, err := model.NewCertificateToken(certificate)
			if err != nil {
				return nil, err
			}
			result = append(result, certificateToken)
		}
		return result, nil
	}

	certificates, err := x509.ParseCertificates(data)
	if err != nil {
		return nil, err
	}
	result := make([]*model.CertificateToken, 0, len(certificates))
	for _, certificate := range certificates {
		certificateToken, err := model.NewCertificateToken(certificate)
		if err != nil {
			return nil, err
		}
		result = append(result, certificateToken)
	}
	return result, nil
}

// Certificate retrieves a certificate by its alias. Port of getCertificate(String).
//
// Returns nil when the alias is unknown, matching Java's null return (a WARN log is dropped,
// slf4j is not load-bearing here).
func (k *KeyStoreCertificateSource) Certificate(alias string) *model.CertificateToken {
	certificateToken, _ := k.entries.Get(k.getKey(alias))
	return certificateToken
}

// AddAllCertificatesToKeyStore adds a list of certificates to the keystore.
// Port of addAllCertificatesToKeyStore(List<CertificateToken>).
func (k *KeyStoreCertificateSource) AddAllCertificatesToKeyStore(certificates []*model.CertificateToken) {
	for _, certificateToken := range certificates {
		k.AddCertificateToKeyStore(certificateToken)
	}
}

// AddCertificateToKeyStore adds a certificate to the keystore. The generated alias is the
// certificate's DSS Id. Port of addCertificateToKeyStore(CertificateToken).
func (k *KeyStoreCertificateSource) AddCertificateToKeyStore(certificateToken *model.CertificateToken) {
	alias := k.getKey(certificateToken.DSSIDAsString())
	k.entries.Set(alias, certificateToken)
	k.CommonCertificateSource.AddCertificate(certificateToken)
}

// AddCertificate always panics: use AddCertificateToKeyStore instead.
// Port of addCertificate(CertificateToken), which throws UnsupportedOperationException.
func (k *KeyStoreCertificateSource) AddCertificate(certificateToAdd *model.CertificateToken) *model.CertificateToken {
	panic("Use AddCertificateToKeyStore(CertificateToken) method to add a certificate to keyStore!")
}

// DeleteCertificateFromKeyStore removes a certificate from the keystore by its alias.
// Port of deleteCertificateFromKeyStore(String).
//
// A missing alias is a no-op, matching Java's WARN-and-return (the log is dropped, slf4j is
// not load-bearing here).
func (k *KeyStoreCertificateSource) DeleteCertificateFromKeyStore(alias string) {
	key := k.getKey(alias)
	certificate, found := k.entries.Get(key)
	if !found {
		return
	}
	k.removeCertificate(certificate)
	k.entries.Delete(key)
}

// ClearAllCertificates removes all certificates from the keystore.
// Port of clearAllCertificates().
func (k *KeyStoreCertificateSource) ClearAllCertificates() {
	for _, alias := range k.entries.Keys() {
		k.DeleteCertificateFromKeyStore(alias)
	}
	k.CommonCertificateSource.reset()
}

// Store writes the keystore to w. Port of store(OutputStream).
//
// PKCS12 has no counterpart to Java's KeyStore#store in this port (golang.org/x/crypto/pkcs12
// only decodes in the vendored version) and returns an error; a PEM-typed source is written as
// concatenated PEM CERTIFICATE blocks.
func (k *KeyStoreCertificateSource) Store(w io.Writer) error {
	switch k.ksType {
	case KeyStoreCertificateSourceType_PEM:
		for _, certificateToken := range k.entries.Values() {
			block := &pem.Block{Type: "CERTIFICATE", Bytes: certificateToken.Encoded()}
			if err := pem.Encode(w, block); err != nil {
				return model.NewDSSErrorMessageCause("Unable to store the keystore", err)
			}
		}
		return nil
	default:
		return model.NewDSSError("Unable to store the keystore: writing a PKCS12 keystore is not supported by this Go port (documented deviation; golang.org/x/crypto/pkcs12 does not encode)")
	}
}

// getKey ports the private getKey(String): a PKCS12 workaround for
// https://bugs.openjdk.java.net/browse/JDK-8079616, lower-casing the alias.
func (k *KeyStoreCertificateSource) getKey(inputKey string) string {
	if k.ksType == KeyStoreCertificateSourceType_PKCS12 {
		return strings.ToLower(inputKey)
	}
	return inputKey
}

// compile-time assertion: a KeyStoreCertificateSource is a CertificateSource.
var _ CertificateSource = (*KeyStoreCertificateSource)(nil)

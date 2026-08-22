package dss

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/token"
)

// ErrNoKeyEntry is returned when a key store holds no usable private key
// entry.
var ErrNoKeyEntry = errors.New("dss: the key store holds no private key entry")

// Signer pairs a key store connection with the one key entry a signature is
// produced with. It is what [Sign] and [SignMultiple] take instead of a raw
// private key: DSS never hands the private key to the signature builder, it
// asks the token to sign the computed data-to-be-signed, so a smart card, an
// HSM or a remote signing service fits the same interface.
//
// Use [OpenPKCS12] for the common case, or [NewSigner] to drive any
// [token.SignatureTokenConnection] the port provides (see the
// [github.com/ryftcore/dss-go/dss/token] package) or one of your own.
type Signer struct {
	conn  token.SignatureTokenConnection
	key   token.DSSPrivateKeyEntry
	owned bool
}

// NewSigner pairs an already-open token connection with the key entry to sign
// with. Close is a no-op for such a Signer: the caller keeps ownership of the
// connection and closes it itself.
func NewSigner(conn token.SignatureTokenConnection, key token.DSSPrivateKeyEntry) (*Signer, error) {
	if conn == nil {
		return nil, errors.New("dss: token connection cannot be nil")
	}
	if key == nil {
		return nil, ErrNoKeyEntry
	}
	return &Signer{conn: conn, key: key}, nil
}

// OpenPKCS12 opens the PKCS#12 (.p12/.pfx) key store at path with the given
// password and selects its first key entry. The returned Signer owns the
// connection, so [Signer.Close] closes it.
//
// RSA, EC, Ed25519 and DSA key stores all load, through the port's own
// RFC 7292 reader. Note the port's accepted gaps in the neighbouring key
// store types: JKS, PKCS#11, the Windows certificate store and the macOS
// Keychain are not supported at all, since none of them has a Go counterpart
// under this port's dependency policy.
func OpenPKCS12(path, password string) (*Signer, error) {
	conn, err := token.NewPkcs12SignatureTokenFromFilepath(path, token.NewPasswordProtection([]byte(password)))
	if err != nil {
		return nil, fmt.Errorf("dss: opening PKCS#12 key store %q: %w", path, err)
	}
	signer, err := signerFromToken(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return signer, nil
}

// OpenPKCS12Bytes is [OpenPKCS12] for a key store already held in memory.
func OpenPKCS12Bytes(store []byte, password string) (*Signer, error) {
	conn, err := token.NewPkcs12SignatureTokenFromBytes(store, token.NewPasswordProtection([]byte(password)))
	if err != nil {
		return nil, fmt.Errorf("dss: opening PKCS#12 key store: %w", err)
	}
	signer, err := signerFromToken(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return signer, nil
}

// signerFromToken picks the first key entry of an opened key store.
func signerFromToken(conn token.SignatureTokenConnection) (*Signer, error) {
	keys, err := conn.Keys()
	if err != nil {
		return nil, fmt.Errorf("dss: listing key entries: %w", err)
	}
	if len(keys) == 0 {
		return nil, ErrNoKeyEntry
	}
	return &Signer{conn: conn, key: keys[0], owned: true}, nil
}

// Certificate returns the signing certificate of the selected key entry.
func (s *Signer) Certificate() *CertificateToken {
	return s.key.Certificate()
}

// CertificateChain returns the certificate chain of the selected key entry,
// as the key store carries it. It is embedded in the signature so that a
// verifier can build the path.
func (s *Signer) CertificateChain() []*CertificateToken {
	return s.key.CertificateChain()
}

// KeyEntry returns the underlying key entry, for code that drives the ported
// services directly.
func (s *Signer) KeyEntry() token.DSSPrivateKeyEntry { return s.key }

// Token returns the underlying token connection, for code that drives the
// ported services directly.
func (s *Signer) Token() token.SignatureTokenConnection { return s.conn }

// Close releases the key store connection when this Signer opened it. It is a
// no-op for a Signer built with [NewSigner], whose connection the caller owns.
func (s *Signer) Close() {
	if s.owned && s.conn != nil {
		s.conn.Close()
	}
}

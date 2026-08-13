// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/AbstractKeyStoreTokenConnection.java (DSS 6.5.RC1).
package token

import (
	"crypto"
	"crypto/x509"

	"github.com/utain/esig/dss/model"
)

// keyStoreEntry is a single alias's private-key-entry contents, standing in for what
// java.security.KeyStore.getEntry(alias, protection) exposes as a KeyStore.PrivateKeyEntry.
type keyStoreEntry struct {
	alias       string
	certificate *x509.Certificate
	chain       []*x509.Certificate
	privateKey  crypto.Signer
}

// keyStore is a minimal, already-decrypted stand-in for java.security.KeyStore (a JDK class, not
// a DSS source file): the set of key entries a concrete token connection's KeyStore() override
// has parsed. There is no lazy per-alias decryption to mirror, since Go has no provider-based
// KeyStore SPI to defer to; every implementation of AbstractKeyStoreTokenConnectionOverrides
// parses everything it can up front.
type keyStore struct {
	entries []keyStoreEntry
}

// aliases lists this keyStore's aliases in encounter order, standing in for
// KeyStore#aliases().
func (k *keyStore) aliases() []string {
	names := make([]string, len(k.entries))
	for i, entry := range k.entries {
		names[i] = entry.alias
	}
	return names
}

// entry looks up the entry for the given alias, standing in for
// KeyStore#isKeyEntry(alias) + KeyStore#getEntry(alias, protection).
func (k *keyStore) entry(alias string) (*keyStoreEntry, bool) {
	for i := range k.entries {
		if k.entries[i].alias == alias {
			return &k.entries[i], true
		}
	}
	return nil, false
}

// AbstractKeyStoreTokenConnectionOverrides is the contract a concrete key-store token connection
// implements, standing in for the two protected abstract methods AbstractKeyStoreTokenConnection
// declares.
type AbstractKeyStoreTokenConnectionOverrides interface {
	// KeyStore gets the key store. Port of the protected abstract getKeyStore().
	KeyStore() (*keyStore, error)

	// KeyProtectionParameter gets the password protection. Port of the protected abstract
	// getKeyProtectionParameter().
	KeyProtectionParameter() *PasswordProtection
}

// AbstractKeyStoreTokenConnection is the keyStore token connection, embedded by every concrete
// key-store-backed SignatureTokenConnection (AppleSignatureToken, KeyStoreSignatureTokenConnection,
// MSCAPISignatureToken, Pkcs11SignatureToken).
type AbstractKeyStoreTokenConnection struct {
	AbstractSignatureTokenConnection

	// overrides points back at the concrete connection; see InitAbstractKeyStoreTokenConnection.
	overrides AbstractKeyStoreTokenConnectionOverrides

	// keyEntryPredicate filters keys obtained from token connection to be returned within the
	// Keys() method. Default: AllKeyEntryPredicate - returns all keys extracted from the token
	// connection.
	keyEntryPredicate DSSKeyEntryPredicate
}

// InitAbstractKeyStoreTokenConnection wires overrides into the base struct and sets the default
// key entry predicate, standing in for the default constructor plus field initializer, following
// the model.TokenBase.InitToken dispatch pattern (see PORTING.md).
func (a *AbstractKeyStoreTokenConnection) InitAbstractKeyStoreTokenConnection(overrides AbstractKeyStoreTokenConnectionOverrides) {
	a.overrides = overrides
	a.keyEntryPredicate = NewAllKeyEntryPredicate()
}

// SetKeyEntryPredicate sets a predicate to filter keys to be returned by the Keys() method.
// Default: NewAllKeyEntryPredicate() - returns all keys extracted from the token connection.
// Port of setKeyEntryPredicate(Predicate<DSSPrivateKeyEntry>).
//
// Panics with the Java message if keyEntryPredicate is nil (Objects.requireNonNull).
func (a *AbstractKeyStoreTokenConnection) SetKeyEntryPredicate(keyEntryPredicate DSSKeyEntryPredicate) {
	if keyEntryPredicate == nil {
		panic("Key entry predicate cannot be null!")
	}
	a.keyEntryPredicate = keyEntryPredicate
}

// Keys implements SignatureTokenConnection. Port of getKeys().
func (a *AbstractKeyStoreTokenConnection) Keys() ([]DSSPrivateKeyEntry, error) {
	list := make([]DSSPrivateKeyEntry, 0)
	ks, err := a.overrides.KeyStore()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to retrieve keys from keystore", err)
	}
	for _, alias := range ks.aliases() {
		dssPrivateKeyEntry, err := a.dssPrivateKeyEntry(ks, alias, a.overrides.KeyProtectionParameter())
		if err != nil {
			return nil, err
		}
		if dssPrivateKeyEntry != nil && a.keyEntryPredicate(dssPrivateKeyEntry) {
			list = append(list, dssPrivateKeyEntry)
		}
	}
	return list, nil
}

// Key allows retrieval of a DSSPrivateKeyEntry by alias. Returns nil if the alias does not
// exist. Port of getKey(String).
func (a *AbstractKeyStoreTokenConnection) Key(alias string) (DSSPrivateKeyEntry, error) {
	return a.KeyWithPassword(alias, a.overrides.KeyProtectionParameter())
}

// KeyWithPassword allows retrieval of a DSSPrivateKeyEntry by alias, unlocked with the given
// passwordProtection. Returns nil if the alias does not exist. Port of
// getKey(String, PasswordProtection).
func (a *AbstractKeyStoreTokenConnection) KeyWithPassword(alias string, passwordProtection *PasswordProtection) (DSSPrivateKeyEntry, error) {
	ks, err := a.overrides.KeyStore()
	if err != nil {
		return nil, err
	}
	return a.dssPrivateKeyEntry(ks, alias, passwordProtection)
}

// dssPrivateKeyEntry ports the private getDSSPrivateKeyEntry(KeyStore, String, PasswordProtection).
//
// DEVIATION: Java logs and returns null both for "no key entry at this alias" and for "the entry
// is not a PrivateKeyEntry" (a class this Go port has no counterpart for, since keyStore only
// ever holds private-key entries); both collapse to the same (nil, nil) result here.
func (a *AbstractKeyStoreTokenConnection) dssPrivateKeyEntry(ks *keyStore, alias string,
	passwordProtection *PasswordProtection) (DSSPrivateKeyEntry, error) {
	entry, ok := ks.entry(alias)
	if !ok {
		return nil, nil
	}
	return NewKSPrivateKeyEntry(entry.alias, entry.certificate, entry.chain, entry.privateKey)
}

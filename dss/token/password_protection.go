// PasswordProtection stands in for java.security.KeyStore.PasswordProtection, a JDK class (not
// a DSS source file) that every dss-token constructor threads through: KeyStoreSignatureTokenConnection
// takes it to unlock the store and, by default, every entry within it; AbstractKeyStoreTokenConnection
// exposes it as the getKeyProtectionParameter() contract; Pkcs11SignatureToken wraps it in a
// PrefilledPasswordCallback.
package token

// PasswordProtection is a password holder that can be destroyed (zeroed) once it is no longer
// needed, mirroring java.security.KeyStore.PasswordProtection implementing
// javax.security.auth.Destroyable.
//
// DEVIATION: Java stores a char[]; Go has no char type, so the password is held as []byte.
type PasswordProtection struct {
	password  []byte
	destroyed bool
}

// NewPasswordProtection creates a PasswordProtection wrapping the given password. Port of
// PasswordProtection(char[] password).
func NewPasswordProtection(password []byte) *PasswordProtection {
	return &PasswordProtection{password: password}
}

// Password gets the password. Port of getPassword().
func (p *PasswordProtection) Password() []byte {
	return p.password
}

// Destroy destroys (zeroes) this object's sensitive information, i.e. clears the password.
// Port of destroy().
func (p *PasswordProtection) Destroy() {
	for i := range p.password {
		p.password[i] = 0
	}
	p.password = nil
	p.destroyed = true
}

// IsDestroyed reports whether this object has been destroyed. Port of isDestroyed().
func (p *PasswordProtection) IsDestroyed() bool {
	return p.destroyed
}

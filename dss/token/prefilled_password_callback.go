// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/PrefilledPasswordCallback.java (DSS 6.5.RC1).
package token

// PrefilledPasswordCallback holds a simple, pre-known password.
type PrefilledPasswordCallback struct {
	// password is the password.
	password *PasswordProtection
}

// NewPrefilledPasswordCallback creates a PrefilledPasswordCallback wrapping the given password.
// Port of the default constructor PrefilledPasswordCallback(PasswordProtection).
func NewPrefilledPasswordCallback(password *PasswordProtection) *PrefilledPasswordCallback {
	return &PrefilledPasswordCallback{password: password}
}

// Password implements PasswordInputCallback. Port of getPassword().
func (c *PrefilledPasswordCallback) Password() []byte {
	return c.password.Password()
}

// Destroy destroys this object's sensitive information, i.e. destroys the wrapped password.
// Port of destroy(), implementing javax.security.auth.Destroyable.
func (c *PrefilledPasswordCallback) Destroy() {
	c.password.Destroy()
}

// IsDestroyed reports whether this object has been destroyed. Port of isDestroyed().
func (c *PrefilledPasswordCallback) IsDestroyed() bool {
	return c.password.IsDestroyed()
}

// compile-time interface assertion.
var _ PasswordInputCallback = (*PrefilledPasswordCallback)(nil)

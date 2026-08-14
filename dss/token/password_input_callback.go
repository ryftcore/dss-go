// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/PasswordInputCallback.java (DSS 6.5.RC1).
package token

// PasswordInputCallback is created in order to retrieve the password from the user when
// accessing the Key Store.
//
// DEVIATION: Java returns char[] (mutable, zeroable in place); Go strings are immutable, so
// this returns []byte instead, which callers may zero after use the same way.
type PasswordInputCallback interface {
	// Password retrieves the password (from the user). Port of getPassword().
	Password() []byte
}

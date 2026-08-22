// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/SunPKCS11Initializer.java (DSS 6.5.RC1).
package token

import "github.com/ryftcore/dss-go/dss/model"

// SunPKCS11InitializerGetProvider initializes the JCA SunPKCS11 provider from the given
// configuration. Port of the static getProvider(String).
//
// DEVIATION: upstream instantiates sun.security.pkcs11.SunPKCS11 through reflection
// (Class.forName + a Constructor<Provider>(InputStream)) precisely because there is no
// compile-time API for it; Go has neither that class nor a PKCS11 provider mechanism to reflect
// into (no cgo, no unmaintained third-party crypto - see PORTING.md), so this always reports the
// gap instead of attempting anything. See Pkcs11SignatureToken, whose KeyStore() override is the
// only caller.
func SunPKCS11InitializerGetProvider(configString string) (any, error) {
	return nil, model.NewDSSErrorMessageCause("Unable to instantiate PKCS11 (JDK < 9) ",
		model.NewDSSError("no counterpart to the JCA SunPKCS11 provider in the Go port"))
}

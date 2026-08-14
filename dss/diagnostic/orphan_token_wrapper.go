// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanTokenWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"fmt"

	"github.com/utain/esig/dss/diagnostic/jaxb"
)

// OrphanTokenWrapperOverrides declares the operations Java's abstract generic class
// OrphanTokenWrapper<T extends XmlOrphanToken> leaves abstract (getBinaries(),
// getDigestAlgoAndValue()) plus the Id() needed by the base's own toString/equals, routed
// through the same override-interface pattern used by AbstractTokenProxyBase. Go has no
// generics-friendly way to express "any type that embeds XmlOrphanToken" once concrete
// subtypes (XmlOrphanCertificateToken, XmlOrphanRevocationToken) add their own extra fields, so
// each concrete wrapper (OrphanCertificateTokenWrapper, OrphanRevocationTokenWrapper) keeps its
// own strongly-typed wrapped field and satisfies this interface directly.
type OrphanTokenWrapperOverrides interface {
	// Id returns the identifier of the orphan token. Port of getId() (delegates to
	// orphanToken.getId() in Java).
	Id() string
	// Binaries returns base64-encoded byte array of the token. Port of the abstract
	// getBinaries().
	Binaries() []byte
	// DigestAlgoAndValue returns digest of the token. Port of the abstract
	// getDigestAlgoAndValue().
	DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue
}

// OrphanTokenWrapperBase carries the state and the concrete behaviour of the Java abstract
// generic class OrphanTokenWrapper<T extends XmlOrphanToken>. Concrete wrappers embed it and
// register themselves with InitOrphanTokenWrapper.
type OrphanTokenWrapperBase struct {
	overrides OrphanTokenWrapperOverrides
}

// InitOrphanTokenWrapper registers the concrete wrapper with its base. Port of the constructor's
// Objects.requireNonNull(orphanToken, "XmlOrphanToken cannot be null!") plus the field
// assignment; must be called exactly once, by the concrete wrapper's constructor.
func (o *OrphanTokenWrapperBase) InitOrphanTokenWrapper(overrides OrphanTokenWrapperOverrides) {
	if overrides == nil {
		panic("XmlOrphanToken cannot be null!")
	}
	o.overrides = overrides
}

func (o *OrphanTokenWrapperBase) orphanTokenOverrides() OrphanTokenWrapperOverrides {
	if o.overrides == nil {
		panic("OrphanTokenWrapper was not initialised: the concrete wrapper must call InitOrphanTokenWrapper in its constructor")
	}
	return o.overrides
}

// Id returns identifier of the orphan token. Port of getId().
func (o *OrphanTokenWrapperBase) Id() string {
	return o.orphanTokenOverrides().Id()
}

// String returns a string representation of the wrapper. Port of toString().
func (o *OrphanTokenWrapperBase) String() string {
	overrides := o.orphanTokenOverrides()
	return fmt.Sprintf("OrphanTokenWrapper Class='%T', Id='%s'", overrides, overrides.Id())
}

// Equals reports whether other wraps an orphan token with the same Id. Port of equals(Object
// obj): Java's `!(obj instanceof OrphanTokenWrapper)` check accepts any OrphanTokenWrapper
// subtype (unlike AbstractTokenProxy.equals(), it does NOT additionally compare getClass()), so
// this compares only by Id, faithfully - no reflect.TypeOf check here.
func (o *OrphanTokenWrapperBase) Equals(other OrphanTokenWrapperOverrides) bool {
	if other == nil {
		return false
	}
	return o.orphanTokenOverrides().Id() == other.Id()
}

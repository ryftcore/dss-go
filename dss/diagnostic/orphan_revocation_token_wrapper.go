// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/OrphanRevocationTokenWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// OrphanRevocationTokenWrapper wraps XML orphan revocation data. Port of
// OrphanTokenWrapper<XmlOrphanRevocationToken>'s revocation-token specialization.
type OrphanRevocationTokenWrapper struct {
	OrphanTokenWrapperBase

	// orphanToken is the wrapped XmlOrphanRevocationToken.
	orphanToken *jaxb.XmlOrphanRevocationToken
}

// NewOrphanRevocationTokenWrapper is the default constructor. Port of
// OrphanRevocationTokenWrapper(XmlOrphanRevocationToken); protected in Java (only subclassed
// by OrphanRevocationWrapper in this package).
func NewOrphanRevocationTokenWrapper(orphanToken *jaxb.XmlOrphanRevocationToken) *OrphanRevocationTokenWrapper {
	w := &OrphanRevocationTokenWrapper{orphanToken: orphanToken}
	w.InitOrphanTokenWrapper(w)
	return w
}

// Id returns identifier of the orphan token. Port of getId() (inherited from
// OrphanTokenWrapper).
func (w *OrphanRevocationTokenWrapper) Id() string {
	return w.orphanToken.Id
}

// RevocationType returns a revocation data type (CRL or OCSP). Port of getRevocationType().
func (w *OrphanRevocationTokenWrapper) RevocationType() enumerations.RevocationType {
	return w.orphanToken.RevocationType
}

// Binaries returns base64-encoded byte array of the token. Port of getBinaries().
func (w *OrphanRevocationTokenWrapper) Binaries() []byte {
	return w.orphanToken.Base64Encoded
}

// DigestAlgoAndValue returns digest of the token. Port of getDigestAlgoAndValue().
func (w *OrphanRevocationTokenWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.orphanToken.DigestAlgoAndValue
}

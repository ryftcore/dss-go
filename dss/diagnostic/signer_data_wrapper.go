// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/SignerDataWrapper.java (DSS 6.5.RC1).
package diagnostic

import "github.com/utain/esig/dss/diagnostic/jaxb"

// SignerDataWrapper wraps the jaxb.XmlSignerData.
type SignerDataWrapper struct {
	// signerData is the wrapped Signed data.
	signerData *jaxb.XmlSignerData
}

// NewSignerDataWrapper is the default constructor. Port of SignerDataWrapper(XmlSignerData).
func NewSignerDataWrapper(signerData *jaxb.XmlSignerData) *SignerDataWrapper {
	return &SignerDataWrapper{signerData: signerData}
}

// Id gets identifier of the signer data. Port of getId().
func (w *SignerDataWrapper) Id() string {
	if w.signerData.Id != nil {
		return string(*w.signerData.Id)
	}
	return ""
}

// ReferencedName gets referenced name of the signer data. Port of getReferencedName().
func (w *SignerDataWrapper) ReferencedName() string {
	if w.signerData.ReferencedName != nil {
		return *w.signerData.ReferencedName
	}
	return ""
}

// DigestAlgoAndValue gets digest algo and value of the signer data. Port of
// getDigestAlgoAndValue().
func (w *SignerDataWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.signerData.DigestAlgoAndValue
}

// Equals reports whether other wraps a signer data with the same Id. Port of equals(Object):
// hashCode() has no Go equivalent (nothing here keys a hash-based collection on a
// SignerDataWrapper) and is dropped, matching the omission pattern documented elsewhere in this
// port (see model/certificate_token.go).
func (w *SignerDataWrapper) Equals(other *SignerDataWrapper) bool {
	if w == other {
		return true
	}
	if other == nil {
		return false
	}
	return w.Id() == other.Id()
}

// String returns a string representation of the wrapper. Port of toString().
func (w *SignerDataWrapper) String() string {
	return "SignerData Id='" + w.Id() + "'"
}

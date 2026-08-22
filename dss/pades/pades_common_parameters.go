// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/PAdESCommonParameters.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart).
package pades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CommonParameters defines a list of common PAdES parameters between signature and
// timestamps.
type CommonParameters interface {
	// SigningDate returns a claimed signing time. Port of #getSigningDate.
	SigningDate() *time.Time

	// Filter returns Filter value. Port of #getFilter.
	Filter() string

	// SubFilter returns SubFilter value. Port of #getSubFilter.
	SubFilter() string

	// ImageParameters returns SignatureImageParameters for field's visual representation. Port
	// of #getImageParameters.
	ImageParameters() *SignatureImageParameters

	// ContentSize returns a length of the reserved /Contents attribute. Port of
	// #getContentSize.
	ContentSize() int

	// DigestAlgorithm returns a DigestAlgorithm to be used to hash the signed/timestamped data.
	// Port of #getDigestAlgorithm.
	DigestAlgorithm() enumerations.DigestAlgorithm

	// EncryptionAlgorithm returns an EncryptionAlgorithm to be used to hash the
	// signed/timestamped data. Port of #getEncryptionAlgorithm.
	EncryptionAlgorithm() enumerations.EncryptionAlgorithm

	// PasswordProtection returns a password used to encrypt a document. Port of
	// #getPasswordProtection.
	PasswordProtection() []byte

	// AppName returns name of an application used to create a signature/timestamp. Port of
	// #getAppName.
	AppName() string

	// DeterministicId returns the deterministic identifier to be used to define a documentId on
	// signing/timestamping, when necessary. Port of #getDeterministicId.
	DeterministicId() string

	// PdfSignatureCache returns an internal variable, used to cache data in order to accelerate
	// the signing process. Port of #getPdfSignatureCache.
	PdfSignatureCache() *PdfSignatureCache

	// Reinit re-inits signature parameters to clean temporary settings. Port of #reinit.
	Reinit()
}

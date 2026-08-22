// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAA.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/eaa"
	"github.com/ryftcore/dss-go/dss/spi"
)

// EAA represents a presentation of Electronic Attestation of Attributes.
type EAA interface {
	model.IdentifierBasedObject

	// Filename returns a name of the presentation of Electronic Attestation of Attributes
	// document, when present. Port of getFilename().
	Filename() string

	// Signatures gets a list of signatures used to issue the Electronic Attestation of
	// Attributes. Port of getSignatures().
	Signatures() []AdvancedSignature

	// EAAType gets the type of the Electronic Attestation of Attributes. Port of getEAAType().
	EAAType() enumerations.EAAType

	// DisclosureValidations gets a list of validation results performed on the selectively
	// disclosable claims. Port of getDisclosureValidations().
	DisclosureValidations() []eaa.DisclosureValidation

	// KeyBindingSignature gets key binding signature, when present. Port of
	// getKeyBindingSignature().
	KeyBindingSignature() AdvancedSignature

	// KeyBindingSignaturePayload gets key binding payload, when present. Port of
	// getKeyBindingSignaturePayload().
	KeyBindingSignaturePayload() EAAKeyBindingPayload

	// DeviceKeyCertificateSource gets the certificate source containing a public key or
	// certificate representation of the device holder. Port of
	// getDeviceKeyCertificateSource().
	DeviceKeyCertificateSource() spi.CertificateSource

	// Payload gets a clear payload of the Electronic Attestation of Attributes. Port of
	// getPayload().
	Payload() EAAPayload

	// SelectiveDisclosuresDigestAlgorithm gets the DigestAlgorithm used for selective
	// disclosures hashes computation. Port of getSelectiveDisclosuresDigestAlgorithm().
	SelectiveDisclosuresDigestAlgorithm() enumerations.DigestAlgorithm

	// ID returns the DSS unique id. It allows to unambiguously identify each token. Port of
	// getId().
	ID() string
}

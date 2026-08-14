// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/ValidationDataContainer.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: ValidationData (Java spi.validation.ValidationData) is owned by sibling
// chunk VAL-D (per signature_validation_context.go's "ValidationData (VAL-D): NewValidationData()
// *ValidationData, AddToken(token model.Token) bool" header entry, already landed in this
// package). This file additionally needs the members below, inferred from every ValidationData
// call this Java source makes:
//
//	func NewValidationData() *ValidationData
//	func (d *ValidationData) AddValidationData(other *ValidationData)
//	func (d *ValidationData) ExcludeCertificateTokens(certificateTokens []*model.CertificateToken)
//	func (d *ValidationData) ExcludeCRLTokens(crlIdentifiers []model.Identifier)
//	func (d *ValidationData) ExcludeOCSPTokens(ocspIdentifiers []model.Identifier)
//	func (d *ValidationData) IsEmpty() bool
//
// Integration note: this header originally assumed ExcludeCRLTokens/ExcludeOCSPTokens took
// []spi.EncapsulatedRevocationTokenIdentifier[R] directly (mirroring
// AllRevocationBinaries()'s return type); the landed ValidationData instead takes the more
// literal translation of Java's Collection<? extends Identifier>, []model.Identifier. Since Go
// slices are not covariant, excludePresentValidationData below converts via the small
// revocationBinaryIdentifiers helper (each EncapsulatedRevocationTokenIdentifier[R] embeds
// model.IdentifierBasedObject, so .DSSID() recovers the model.Identifier that
// excludeCRLTokens(Collection<? extends Identifier>) receives directly in Java).
package validation

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// revocationBinaryIdentifiers converts a slice of EncapsulatedRevocationTokenIdentifier[R]
// (the return type of OfflineRevocationSource[R].AllRevocationBinaries) into the
// []model.Identifier ValidationData.ExcludeCRLTokens/ExcludeOCSPTokens expect, standing in for
// Java's covariant Collection<? extends Identifier> argument passing.
func revocationBinaryIdentifiers[R revocation.Revocation](binaries []spi.EncapsulatedRevocationTokenIdentifier[R]) []model.Identifier {
	identifiers := make([]model.Identifier, 0, len(binaries))
	for _, binary := range binaries {
		identifiers = append(identifiers, binary.DSSID())
	}
	return identifiers
}

// ValidationDataContainer contains a ValidationData for a list of signatures/timestamps.
//
// Judgment call: upstream keys signatureValidationDataMap/timestampValidationDataMap in
// java.util.HashMaps, whose iteration order is arbitrary but stable within a JVM run;
// AdvancedSignature/TimestampToken equality in Java is default (identity-based), which a Go
// map keyed on the same interface/pointer values reproduces directly (unlike the
// value-equality tokens elsewhere in this package that need the sorted-slice workaround; see
// token_status.go). A bare Go map's iteration order, unlike Java's, is randomized on every
// run; since Signatures() and DetachedTimestamps() return in this map's iteration order, both
// maps are kept insertion-ordered (slice + index map, PORTING.md's Collections rule) instead.
type ValidationDataContainer struct {
	// signatureValidationDataMap maps signatures to their corresponding ValidationData.
	signatureValidationDataMap *utils.OrderedMap[AdvancedSignature, *ValidationData]

	// timestampValidationDataMap maps timestamps to their corresponding ValidationData.
	timestampValidationDataMap *utils.OrderedMap[*TimestampToken, *ValidationData]
}

// NewValidationDataContainer instantiates empty maps of tokens and validation data
// relationships. Ports the default constructor.
func NewValidationDataContainer() *ValidationDataContainer {
	return &ValidationDataContainer{
		signatureValidationDataMap: utils.NewOrderedMap[AdvancedSignature, *ValidationData](),
		timestampValidationDataMap: utils.NewOrderedMap[*TimestampToken, *ValidationData](),
	}
}

// AddValidationDataForSignature adds validation data to the container. Port of the
// addValidationData(AdvancedSignature, ValidationData) overload; Go has no overloading, so the
// signature and timestamp overloads carry different names.
func (c *ValidationDataContainer) AddValidationDataForSignature(signature AdvancedSignature, validationData *ValidationData) {
	c.signatureValidationDataMap.Set(signature, validationData)
}

// AddValidationDataForTimestamp adds validation data to the container. Port of the
// addValidationData(TimestampToken, ValidationData) overload.
func (c *ValidationDataContainer) AddValidationDataForTimestamp(timestampToken *TimestampToken, validationData *ValidationData) {
	c.timestampValidationDataMap.Set(timestampToken, validationData)
}

// ValidationDataForSignature returns the related ValidationData for the given signature. Port
// of the getValidationData(AdvancedSignature) overload.
func (c *ValidationDataContainer) ValidationDataForSignature(signature AdvancedSignature) *ValidationData {
	validationData, _ := c.signatureValidationDataMap.Get(signature)
	return validationData
}

// ValidationDataForTimestamp returns the related ValidationData for the given timestamp token.
// Port of the getValidationData(TimestampToken) overload.
func (c *ValidationDataContainer) ValidationDataForTimestamp(timestampToken *TimestampToken) *ValidationData {
	validationData, _ := c.timestampValidationDataMap.Get(timestampToken)
	return validationData
}

// AllValidationData returns a combined validation data for all tokens. Port of
// getAllValidationData().
func (c *ValidationDataContainer) AllValidationData() *ValidationData {
	result := NewValidationData()
	for _, validationData := range c.signatureValidationDataMap.Values() {
		result.AddValidationData(validationData)
	}
	for _, validationData := range c.timestampValidationDataMap.Values() {
		result.AddValidationData(validationData)
	}
	return result
}

// Signatures returns a collection of AdvancedSignatures. Port of getSignatures().
func (c *ValidationDataContainer) Signatures() []AdvancedSignature {
	return c.signatureValidationDataMap.Keys()
}

// DetachedTimestamps returns a collection of TimestampTokens. Port of getDetachedTimestamps().
func (c *ValidationDataContainer) DetachedTimestamps() []*TimestampToken {
	return c.timestampValidationDataMap.Keys()
}

// IsEmpty checks if the validation data for inclusion is empty. Port of isEmpty().
func (c *ValidationDataContainer) IsEmpty() bool {
	for _, validationData := range c.signatureValidationDataMap.Values() {
		if !validationData.IsEmpty() {
			return false
		}
	}
	for _, validationData := range c.timestampValidationDataMap.Values() {
		if !validationData.IsEmpty() {
			return false
		}
	}
	return true
}

// AllValidationDataForSignature returns a complete validation data for a signature, including
// the data for incorporated timestamps and/or counter-signatures. Port of
// getAllValidationDataForSignature(AdvancedSignature).
func (c *ValidationDataContainer) AllValidationDataForSignature(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	validationDataForInclusion.AddValidationData(c.validationDataForSignature(signature))
	validationDataForInclusion.AddValidationData(c.validationDataForSignatureTimestamps(signature))
	validationDataForInclusion.AddValidationData(c.validationDataForCounterSignatures(signature))
	validationDataForInclusion.AddValidationData(c.validationDataForCounterSignatureTimestamps(signature))

	return validationDataForInclusion
}

// AllValidationDataForSignatureForInclusion returns a complete validation data for a
// signature, including the data for incorporated timestamps and/or counter-signatures, but
// excluding the tokens already incorporated within the signature. Port of
// getAllValidationDataForSignatureForInclusion(AdvancedSignature).
func (c *ValidationDataContainer) AllValidationDataForSignatureForInclusion(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	validationDataForInclusion.AddValidationData(c.ValidationDataForSignatureForInclusion(signature))
	validationDataForInclusion.AddValidationData(c.ValidationDataForSignatureTimestampsForInclusion(signature))
	validationDataForInclusion.AddValidationData(c.ValidationDataForCounterSignaturesForInclusion(signature))
	validationDataForInclusion.AddValidationData(c.ValidationDataForCounterSignatureTimestampsForInclusion(signature))

	return validationDataForInclusion
}

// excludePresentValidationData is the port of the private excludePresentValidationData(ValidationData, AdvancedSignature).
func excludePresentValidationData(validationData *ValidationData, signature AdvancedSignature) {
	validationData.ExcludeCertificateTokens(signature.CertificateSource().Certificates())
	validationData.ExcludeCRLTokens(revocationBinaryIdentifiers(signature.CRLSource().AllRevocationBinaries()))
	validationData.ExcludeOCSPTokens(revocationBinaryIdentifiers(signature.OCSPSource().AllRevocationBinaries()))
}

// validationDataForSignature returns all validation data for the signature. Port of the
// protected getValidationDataForSignature(AdvancedSignature).
func (c *ValidationDataContainer) validationDataForSignature(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	signatureValidationData := c.ValidationDataForSignature(signature)
	validationDataForInclusion.AddValidationData(signatureValidationData)

	return validationDataForInclusion
}

// ValidationDataForSignatureForInclusion returns all validation data for a signature, but
// excluding the tokens already incorporated within the signature. Port of
// getValidationDataForSignatureForInclusion(AdvancedSignature).
func (c *ValidationDataContainer) ValidationDataForSignatureForInclusion(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := c.validationDataForSignature(signature)
	excludePresentValidationData(validationDataForInclusion, signature)
	return validationDataForInclusion
}

// validationDataForCounterSignatures returns all validation data for the incorporated
// counter-signatures. Port of the protected getValidationDataForCounterSignatures(AdvancedSignature).
func (c *ValidationDataContainer) validationDataForCounterSignatures(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	for _, counterSignature := range signature.CounterSignatures() {
		counterSignatureValidationData := c.ValidationDataForSignature(counterSignature)
		validationDataForInclusion.AddValidationData(counterSignatureValidationData)
	}

	return validationDataForInclusion
}

// ValidationDataForCounterSignaturesForInclusion returns all validation data for incorporated
// counter-signatures, but excluding the tokens already incorporated within the signature or
// counter-signatures. Port of getValidationDataForCounterSignaturesForInclusion(AdvancedSignature).
func (c *ValidationDataContainer) ValidationDataForCounterSignaturesForInclusion(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := c.validationDataForCounterSignatures(signature)
	excludePresentValidationData(validationDataForInclusion, signature)
	for _, counterSignature := range signature.CounterSignatures() {
		excludePresentValidationData(validationDataForInclusion, counterSignature)
	}
	return validationDataForInclusion
}

// validationDataForSignatureTimestamps returns all validation data for the timestamps
// incorporated within the signature. Port of the protected
// getValidationDataForSignatureTimestamps(AdvancedSignature).
func (c *ValidationDataContainer) validationDataForSignatureTimestamps(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	for _, timestampToken := range signature.AllTimestamps() {
		timestampValidationData := c.ValidationDataForTimestamp(timestampToken)
		validationDataForInclusion.AddValidationData(timestampValidationData)
	}

	return validationDataForInclusion
}

// ValidationDataForSignatureTimestampsForInclusion returns all validation data for the
// timestamps incorporated within the signature, but excluding the tokens already incorporated
// within the signature. Port of getValidationDataForSignatureTimestampsForInclusion(AdvancedSignature).
func (c *ValidationDataContainer) ValidationDataForSignatureTimestampsForInclusion(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := c.validationDataForSignatureTimestamps(signature)
	excludePresentValidationData(validationDataForInclusion, signature)
	return validationDataForInclusion
}

// validationDataForCounterSignatureTimestamps returns all validation data for the timestamps
// incorporated within counter signatures of the current signature. Port of the protected
// getValidationDataForCounterSignatureTimestamps(AdvancedSignature).
func (c *ValidationDataContainer) validationDataForCounterSignatureTimestamps(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := NewValidationData()

	for _, counterSignature := range signature.CounterSignatures() {
		for _, timestampToken := range counterSignature.AllTimestamps() {
			timestampValidationData := c.ValidationDataForTimestamp(timestampToken)
			validationDataForInclusion.AddValidationData(timestampValidationData)
		}
	}

	return validationDataForInclusion
}

// ValidationDataForCounterSignatureTimestampsForInclusion returns all validation data for the
// timestamps incorporated within counter signatures of the current signature, but excluding
// the tokens already incorporated within the signature. Port of
// getValidationDataForCounterSignatureTimestampsForInclusion(AdvancedSignature).
func (c *ValidationDataContainer) ValidationDataForCounterSignatureTimestampsForInclusion(signature AdvancedSignature) *ValidationData {
	validationDataForInclusion := c.validationDataForCounterSignatureTimestamps(signature)
	excludePresentValidationData(validationDataForInclusion, signature)
	for _, counterSignature := range signature.CounterSignatures() {
		excludePresentValidationData(validationDataForInclusion, counterSignature)
	}
	return validationDataForInclusion
}

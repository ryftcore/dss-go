// Ported from dss-model/.../AbstractSerializableSignatureParameters.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// AbstractSerializableSignatureParameters holds parameters for a signature
// creation/extension, generic over the SerializableTimestampParameters
// implementation TP used for content/signature/archive timestamps.
//
// SerializableSignatureParameters (the implemented Java interface),
// SerializableTimestampParameters and BLevelParameters are outside this
// manifest; assumed to already exist in this package.
type AbstractSerializableSignatureParameters[TP SerializableTimestampParameters] struct {
	// checkCertificateRevocation indicates whether a signing certificate
	// revocation shall be checked. Default: false.
	checkCertificateRevocation bool

	// generateTBSWithoutCertificate indicates if it is possible to generate
	// ToBeSigned data without the signing certificate. Default: false.
	generateTBSWithoutCertificate bool

	// signatureLevel is the expected signature level.
	signatureLevel enumerations.SignatureLevel

	// signaturePackaging is the expected signature packaging.
	signaturePackaging enumerations.SignaturePackaging

	// signatureAlgorithm indicates the algorithms used to sign
	// ds:SignedInfo (XAdES).
	signatureAlgorithm enumerations.SignatureAlgorithm

	// encryptionAlgorithm is normally extracted automatically from the
	// signing token.
	encryptionAlgorithm enumerations.EncryptionAlgorithm

	// digestAlgorithm is the digest algorithm used to hash ds:SignedInfo
	// (XAdES).
	digestAlgorithm enumerations.DigestAlgorithm

	// referenceDigestAlgorithm is the digest algorithm used to hash
	// ds:Reference (XAdES).
	referenceDigestAlgorithm enumerations.DigestAlgorithm

	// bLevelParams holds the parameters related to B-level (signed
	// properties).
	bLevelParams *BLevelParameters

	// validationDataEncapsulationStrategy defines the validation data
	// encapsulation mechanism on -LT and -LTA level augmentation.
	validationDataEncapsulationStrategy enumerations.ValidationDataEncapsulationStrategy

	// ContentTimestampParameters holds parameters related to the content
	// timestamp (Baseline-B). Exported: Java exposes it as `protected` for
	// direct subclass access, and concrete parameter subclasses live in
	// other Go packages (cades, xades, ...).
	ContentTimestampParameters TP

	// SignatureTimestampParameters holds parameters related to the
	// signature timestamp (Baseline-T).
	SignatureTimestampParameters TP

	// ArchiveTimestampParameters holds parameters related to the archive
	// timestamp (Baseline-LTA).
	ArchiveTimestampParameters TP
}

// NewAbstractSerializableSignatureParameters instantiates the object with
// default values. Ports the protected no-arg constructor.
func NewAbstractSerializableSignatureParameters[TP SerializableTimestampParameters]() AbstractSerializableSignatureParameters[TP] {
	signatureAlgorithm := enumerations.SignatureAlgorithmRSASSAPSSSHA512MGF1
	return AbstractSerializableSignatureParameters[TP]{
		signatureAlgorithm:                  signatureAlgorithm,
		encryptionAlgorithm:                 signatureAlgorithm.EncryptionAlgorithm(),
		digestAlgorithm:                     signatureAlgorithm.DigestAlgorithm(),
		bLevelParams:                        NewBLevelParameters(),
		validationDataEncapsulationStrategy: enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData,
	}
}

// CheckCertificateRevocation ports
// AbstractSerializableSignatureParameters#isCheckCertificateRevocation.
func (p *AbstractSerializableSignatureParameters[TP]) CheckCertificateRevocation() bool {
	return p.checkCertificateRevocation
}

// SetCheckCertificateRevocation allows setting whether a revocation status
// for a signing certificate should be checked on signature creation or
// T-level extension.
//
// NOTE: in order to specify a behavior for this check, the relevant alerts
// should be specified within a CertificateVerifier instance, used in a
// service for signing/extension.
//
// Default: false (do not perform revocation data check on signature
// creation/T-level extension).
func (p *AbstractSerializableSignatureParameters[TP]) SetCheckCertificateRevocation(checkCertificateRevocation bool) {
	p.checkCertificateRevocation = checkCertificateRevocation
}

// GenerateTBSWithoutCertificate indicates if it is possible to generate
// ToBeSigned data without the signing certificate. Default: false. Ports
// AbstractSerializableSignatureParameters#isGenerateTBSWithoutCertificate.
func (p *AbstractSerializableSignatureParameters[TP]) GenerateTBSWithoutCertificate() bool {
	return p.generateTBSWithoutCertificate
}

// SetGenerateTBSWithoutCertificate changes the default behaviour regarding
// the requirement of a signing certificate to generate ToBeSigned data.
// NOTE: when using this method, ensure the same EncryptionAlgorithm is
// provided via SetEncryptionAlgorithm as the one used on signature value
// creation.
func (p *AbstractSerializableSignatureParameters[TP]) SetGenerateTBSWithoutCertificate(generateTBSWithoutCertificate bool) {
	p.generateTBSWithoutCertificate = generateTBSWithoutCertificate
}

// SignatureLevel gets the expected signature level: XAdES_BASELINE_T,
// CAdES_BASELINE_LTA...
func (p *AbstractSerializableSignatureParameters[TP]) SignatureLevel() enumerations.SignatureLevel {
	return p.signatureLevel
}

// SetSignatureLevel sets the signature level. Panics if signatureLevel is
// the zero value (Java Objects.requireNonNull("Signature Level cannot be
// null")).
func (p *AbstractSerializableSignatureParameters[TP]) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) {
	if signatureLevel == "" {
		panic("Signature Level cannot be null")
	}
	p.signatureLevel = signatureLevel
}

// SignaturePackaging gets the expected signature packaging.
func (p *AbstractSerializableSignatureParameters[TP]) SignaturePackaging() enumerations.SignaturePackaging {
	return p.signaturePackaging
}

// SetSignaturePackaging sets the expected signature packaging.
func (p *AbstractSerializableSignatureParameters[TP]) SetSignaturePackaging(signaturePackaging enumerations.SignaturePackaging) {
	p.signaturePackaging = signaturePackaging
}

// DigestAlgorithm ports
// AbstractSerializableSignatureParameters#getDigestAlgorithm.
func (p *AbstractSerializableSignatureParameters[TP]) DigestAlgorithm() enumerations.DigestAlgorithm {
	return p.digestAlgorithm
}

// SetDigestAlgorithm sets the digest algorithm, and recomputes
// signatureAlgorithm from the current encryptionAlgorithm when defined.
// Panics if digestAlgorithm is the zero value (Java
// Objects.requireNonNull("DigestAlgorithm cannot be null!")).
func (p *AbstractSerializableSignatureParameters[TP]) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	p.digestAlgorithm = digestAlgorithm
	if p.encryptionAlgorithm != "" {
		p.signatureAlgorithm = enumerations.SignatureAlgorithmGetAlgorithm(p.encryptionAlgorithm, p.digestAlgorithm)
	}
}

// EncryptionAlgorithm ports
// AbstractSerializableSignatureParameters#getEncryptionAlgorithm.
func (p *AbstractSerializableSignatureParameters[TP]) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return p.encryptionAlgorithm
}

// SetEncryptionAlgorithm sets the encryption algorithm to be used on
// signature creation, useful when a specific encryption algorithm is
// expected; it becomes the algorithm used to create the SignatureValue.
// NOTE: the encryption algorithm is automatically extracted from the
// certificate's key when a signing certificate is set elsewhere.
func (p *AbstractSerializableSignatureParameters[TP]) SetEncryptionAlgorithm(encryptionAlgorithm enumerations.EncryptionAlgorithm) {
	p.encryptionAlgorithm = encryptionAlgorithm
	if p.digestAlgorithm != "" {
		p.signatureAlgorithm = enumerations.SignatureAlgorithmGetAlgorithm(p.encryptionAlgorithm, p.digestAlgorithm)
	}
}

// SignatureAlgorithm ports
// AbstractSerializableSignatureParameters#getSignatureAlgorithm.
func (p *AbstractSerializableSignatureParameters[TP]) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	return p.signatureAlgorithm
}

// ReferenceDigestAlgorithm gets the digest algorithm for ds:Reference or
// message-digest attribute.
func (p *AbstractSerializableSignatureParameters[TP]) ReferenceDigestAlgorithm() enumerations.DigestAlgorithm {
	return p.referenceDigestAlgorithm
}

// SetReferenceDigestAlgorithm sets the DigestAlgorithm to be used for
// reference digest calculation.
func (p *AbstractSerializableSignatureParameters[TP]) SetReferenceDigestAlgorithm(referenceDigestAlgorithm enumerations.DigestAlgorithm) {
	p.referenceDigestAlgorithm = referenceDigestAlgorithm
}

// BLevel gets the Baseline B parameters (signed properties). Ports
// AbstractSerializableSignatureParameters#bLevel.
func (p *AbstractSerializableSignatureParameters[TP]) BLevel() *BLevelParameters {
	return p.bLevelParams
}

// SetBLevelParams sets the Baseline B parameters (signed properties).
// Panics if bLevelParams is nil (Java Objects.requireNonNull(
// "bLevelParams cannot be null!")).
func (p *AbstractSerializableSignatureParameters[TP]) SetBLevelParams(bLevelParams *BLevelParameters) {
	if bLevelParams == nil {
		panic("bLevelParams cannot be null!")
	}
	p.bLevelParams = bLevelParams
}

// ValidationDataEncapsulationStrategy gets the validation data
// encapsulation mechanism to be used on -LT and -LTA level augmentation.
func (p *AbstractSerializableSignatureParameters[TP]) ValidationDataEncapsulationStrategy() enumerations.ValidationDataEncapsulationStrategy {
	return p.validationDataEncapsulationStrategy
}

// SetValidationDataEncapsulationStrategy sets the validation data
// encapsulation mechanism to be used on -LT and -LTA level augmentation.
//
// Default: CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_AND_ANY_VALIDATION_DATA
// (the validation data for a signature's certificate chain is included
// within CertificateValues and RevocationValues elements on LT-level
// augmentation, and within AnyValidationData element on LTA-level
// augmentation. The validation data for all present timestamps will be
// included within TimeStampValidationData element.)
//
// NOTE: Applicable only for XAdES and JAdES signatures. Panics if
// validationDataEncapsulationStrategy is the zero value.
func (p *AbstractSerializableSignatureParameters[TP]) SetValidationDataEncapsulationStrategy(validationDataEncapsulationStrategy enumerations.ValidationDataEncapsulationStrategy) {
	if validationDataEncapsulationStrategy == "" {
		panic("ValidationDataEncapsulationStrategy cannot be null!")
	}
	p.validationDataEncapsulationStrategy = validationDataEncapsulationStrategy
}

// GetContentTimestampParameters gets the parameters for content timestamp
// (Baseline-B). Ports
// AbstractSerializableSignatureParameters#getContentTimestampParameters:
// the Java base implementation always throws UnsupportedOperationException
// ("Cannot extract ContentTimestampParameters! Not implemented by
// default."); concrete parameter subclasses override it. Direct field
// access via ContentTimestampParameters is available to embedding types.
func (p *AbstractSerializableSignatureParameters[TP]) GetContentTimestampParameters() TP {
	panic("Cannot extract ContentTimestampParameters! Not implemented by default.")
}

// SetContentTimestampParameters sets the parameters to produce the content
// timestamp (Baseline-B).
func (p *AbstractSerializableSignatureParameters[TP]) SetContentTimestampParameters(contentTimestampParameters TP) {
	p.ContentTimestampParameters = contentTimestampParameters
}

// GetSignatureTimestampParameters gets the parameters for signature
// timestamp (Baseline-T). See GetContentTimestampParameters doc: the base
// implementation always panics.
func (p *AbstractSerializableSignatureParameters[TP]) GetSignatureTimestampParameters() TP {
	panic("Cannot extract SignatureTimestampParameters! Not implemented by default.")
}

// SetSignatureTimestampParameters sets the parameters to produce the
// signature timestamp (Baseline-T).
func (p *AbstractSerializableSignatureParameters[TP]) SetSignatureTimestampParameters(signatureTimestampParameters TP) {
	p.SignatureTimestampParameters = signatureTimestampParameters
}

// GetArchiveTimestampParameters gets the parameters for archive timestamp
// (Baseline-LTA). See GetContentTimestampParameters doc: the base
// implementation always panics.
func (p *AbstractSerializableSignatureParameters[TP]) GetArchiveTimestampParameters() TP {
	panic("Cannot extract ArchiveTimestampParameters! Not implemented by default.")
}

// SetArchiveTimestampParameters sets the parameters to produce the archive
// timestamp (Baseline-LTA).
func (p *AbstractSerializableSignatureParameters[TP]) SetArchiveTimestampParameters(archiveTimestampParameters TP) {
	p.ArchiveTimestampParameters = archiveTimestampParameters
}

// Equals ports AbstractSerializableSignatureParameters#equals.
//
// Java's getClass() check belongs to the concrete subclass, so this method compares only
// the fields the abstract class owns and is meant to be called from a subclass's own
// Equals, the way Java subclasses call super.equals(o).
//
// NOTE: upstream deliberately leaves validationDataEncapsulationStrategy out of equals(),
// hashCode() and toString(); the omission is reproduced here rather than "fixed".
func (p *AbstractSerializableSignatureParameters[TP]) Equals(other *AbstractSerializableSignatureParameters[TP]) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.checkCertificateRevocation != other.checkCertificateRevocation ||
		p.generateTBSWithoutCertificate != other.generateTBSWithoutCertificate ||
		p.signatureLevel != other.signatureLevel ||
		p.signaturePackaging != other.signaturePackaging ||
		p.signatureAlgorithm != other.signatureAlgorithm ||
		p.encryptionAlgorithm != other.encryptionAlgorithm ||
		p.digestAlgorithm != other.digestAlgorithm ||
		p.referenceDigestAlgorithm != other.referenceDigestAlgorithm {
		return false
	}
	if !abstractSerializableSignatureParametersBLevelEquals(p.bLevelParams, other.bLevelParams) {
		return false
	}
	return abstractSerializableSignatureParametersTimestampEquals(p.ContentTimestampParameters, other.ContentTimestampParameters) &&
		abstractSerializableSignatureParametersTimestampEquals(p.SignatureTimestampParameters, other.SignatureTimestampParameters) &&
		abstractSerializableSignatureParametersTimestampEquals(p.ArchiveTimestampParameters, other.ArchiveTimestampParameters)
}

// abstractSerializableSignatureParametersBLevelEquals ports Objects.equals(bLevelParams, ...),
// which treats two nulls as equal.
func abstractSerializableSignatureParametersBLevelEquals(a, b *BLevelParameters) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// abstractSerializableSignatureParametersTimestampEquals ports
// Objects.equals(contentTimestampParameters, ...) for the TP type parameter, whose concrete
// type is only known to the subclass. It prefers the value's own Equals when it has one
// (every SerializableTimestampParameters implementation in dss-model does), and otherwise
// falls back to reflect.DeepEqual.
func abstractSerializableSignatureParametersTimestampEquals[TP SerializableTimestampParameters](a, b TP) bool {
	if equatable, ok := any(a).(interface{ Equals(TP) bool }); ok {
		return equatable.Equals(b)
	}
	return reflect.DeepEqual(a, b)
}

// String ports AbstractSerializableSignatureParameters#toString.
func (p *AbstractSerializableSignatureParameters[TP]) String() string {
	return fmt.Sprintf("AbstractSerializableSignatureParameters [checkCertificateRevocation=%v, generateTBSWithoutCertificate=%v, signatureLevel=%v, signaturePackaging=%v, signatureAlgorithm=%v, encryptionAlgorithm=%v, digestAlgorithm=%v, referenceDigestAlgorithm=%v, bLevelParams=%v, contentTimestampParameters=%v, signatureTimestampParameters=%v, archiveTimestampParameters=%v]",
		p.checkCertificateRevocation, p.generateTBSWithoutCertificate, p.signatureLevel, p.signaturePackaging,
		p.signatureAlgorithm, p.encryptionAlgorithm, p.digestAlgorithm, p.referenceDigestAlgorithm,
		p.bLevelParams, p.ContentTimestampParameters, p.SignatureTimestampParameters, p.ArchiveTimestampParameters)
}

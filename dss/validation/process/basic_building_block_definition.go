// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/BasicBuildingBlockDefinition.java (DSS 6.5.RC1).
package process

import "github.com/ryftcore/dss-go/dss/enumerations"

// BasicBuildingBlockDefinition is the definition of signature Basic Building
// Blocks as per EN 319 102-1. Implements enumerations.UriBasedEnum.
type BasicBuildingBlockDefinition string

const (
	// BasicBuildingBlockDefinitionFormatChecking: 5.2.2 Format Checking.
	BasicBuildingBlockDefinitionFormatChecking BasicBuildingBlockDefinition = "FORMAT_CHECKING"
	// BasicBuildingBlockDefinitionIdentificationOfTheSigningCertificate:
	// 5.2.3 Identification of the signing certificate.
	BasicBuildingBlockDefinitionIdentificationOfTheSigningCertificate BasicBuildingBlockDefinition = "IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE"
	// BasicBuildingBlockDefinitionValidationContextInitialization: 5.2.4
	// Validation context initialization.
	BasicBuildingBlockDefinitionValidationContextInitialization BasicBuildingBlockDefinition = "VALIDATION_CONTEXT_INITIALIZATION"
	// BasicBuildingBlockDefinitionRevocationFreshnessChecker: 5.2.5
	// Revocation freshness checker.
	BasicBuildingBlockDefinitionRevocationFreshnessChecker BasicBuildingBlockDefinition = "REVOCATION_FRESHNESS_CHECKER"
	// BasicBuildingBlockDefinitionX509CertificateValidation: 5.2.6 X.509
	// certificate validation.
	BasicBuildingBlockDefinitionX509CertificateValidation BasicBuildingBlockDefinition = "X509_CERTIFICATE_VALIDATION"
	// BasicBuildingBlockDefinitionCryptographicVerification: 5.2.7
	// Cryptographic verification.
	BasicBuildingBlockDefinitionCryptographicVerification BasicBuildingBlockDefinition = "CRYPTOGRAPHIC_VERIFICATION"
	// BasicBuildingBlockDefinitionSignatureAcceptanceValidation: 5.2.8
	// Signature Acceptance Validation (SAV).
	BasicBuildingBlockDefinitionSignatureAcceptanceValidation BasicBuildingBlockDefinition = "SIGNATURE_ACCEPTANCE_VALIDATION"
	// BasicBuildingBlockDefinitionPastCertificateValidation: 5.6.2.1 Past
	// certificate validation.
	BasicBuildingBlockDefinitionPastCertificateValidation BasicBuildingBlockDefinition = "PAST_CERTIFICATE_VALIDATION"
	// BasicBuildingBlockDefinitionValidationTimeSliding: 5.6.2.2 Validation
	// time sliding process.
	BasicBuildingBlockDefinitionValidationTimeSliding BasicBuildingBlockDefinition = "VALIDATION_TIME_SLIDING"
	// BasicBuildingBlockDefinitionPastSignatureValidation: 5.6.2.4 Past
	// signature validation building block.
	BasicBuildingBlockDefinitionPastSignatureValidation BasicBuildingBlockDefinition = "PAST_SIGNATURE_VALIDATION"
)

// basicBuildingBlockDefinitionURIs holds the URI identifying each
// BasicBuildingBlock.
var basicBuildingBlockDefinitionURIs = map[BasicBuildingBlockDefinition]string{
	BasicBuildingBlockDefinitionFormatChecking:                        "urn:cef:dss:bbb:formatChecking",
	BasicBuildingBlockDefinitionIdentificationOfTheSigningCertificate: "urn:cef:dss:bbb:identificationOfTheSigningCertificate",
	BasicBuildingBlockDefinitionValidationContextInitialization:       "urn:cef:dss:bbb:validationContextInitialization",
	BasicBuildingBlockDefinitionRevocationFreshnessChecker:            "urn:cef:dss:bbb:revocationFreshnessChecker",
	BasicBuildingBlockDefinitionX509CertificateValidation:             "urn:cef:dss:bbb:x509CertificateValidation",
	BasicBuildingBlockDefinitionCryptographicVerification:             "urn:cef:dss:bbb:cryptographicVerification",
	BasicBuildingBlockDefinitionSignatureAcceptanceValidation:         "urn:cef:dss:bbb:signatureAcceptanceValidation",
	BasicBuildingBlockDefinitionPastCertificateValidation:             "urn:cef:dss:bbb:pastCertificateValidation",
	BasicBuildingBlockDefinitionValidationTimeSliding:                 "urn:cef:dss:bbb:validationTimeSliding",
	BasicBuildingBlockDefinitionPastSignatureValidation:               "urn:cef:dss:bbb:pastSignatureValidation",
}

// BasicBuildingBlockDefinitionValues returns all constants in declaration order.
func BasicBuildingBlockDefinitionValues() []BasicBuildingBlockDefinition {
	return []BasicBuildingBlockDefinition{
		BasicBuildingBlockDefinitionFormatChecking,
		BasicBuildingBlockDefinitionIdentificationOfTheSigningCertificate,
		BasicBuildingBlockDefinitionValidationContextInitialization,
		BasicBuildingBlockDefinitionRevocationFreshnessChecker,
		BasicBuildingBlockDefinitionX509CertificateValidation,
		BasicBuildingBlockDefinitionCryptographicVerification,
		BasicBuildingBlockDefinitionSignatureAcceptanceValidation,
		BasicBuildingBlockDefinitionPastCertificateValidation,
		BasicBuildingBlockDefinitionValidationTimeSliding,
		BasicBuildingBlockDefinitionPastSignatureValidation,
	}
}

// URI returns the URI identifying the BasicBuildingBlock. Port of getUri().
func (b BasicBuildingBlockDefinition) URI() string {
	return basicBuildingBlockDefinitionURIs[b]
}

// compile-time assertion that the enum implements the marker interface Java
// declares on it.
var _ enumerations.UriBasedEnum = BasicBuildingBlockDefinition("")

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/BasicBuildingBlockDefinition.java (DSS 6.5.RC1).
package process

import "github.com/utain/esig/dss/enumerations"

// BasicBuildingBlockDefinition is the definition of signature Basic Building
// Blocks as per EN 319 102-1. Implements enumerations.UriBasedEnum.
type BasicBuildingBlockDefinition string

const (
	// BasicBuildingBlockDefinition_FORMAT_CHECKING: 5.2.2 Format Checking.
	BasicBuildingBlockDefinition_FORMAT_CHECKING BasicBuildingBlockDefinition = "FORMAT_CHECKING"
	// BasicBuildingBlockDefinition_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE:
	// 5.2.3 Identification of the signing certificate.
	BasicBuildingBlockDefinition_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE BasicBuildingBlockDefinition = "IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE"
	// BasicBuildingBlockDefinition_VALIDATION_CONTEXT_INITIALIZATION: 5.2.4
	// Validation context initialization.
	BasicBuildingBlockDefinition_VALIDATION_CONTEXT_INITIALIZATION BasicBuildingBlockDefinition = "VALIDATION_CONTEXT_INITIALIZATION"
	// BasicBuildingBlockDefinition_REVOCATION_FRESHNESS_CHECKER: 5.2.5
	// Revocation freshness checker.
	BasicBuildingBlockDefinition_REVOCATION_FRESHNESS_CHECKER BasicBuildingBlockDefinition = "REVOCATION_FRESHNESS_CHECKER"
	// BasicBuildingBlockDefinition_X509_CERTIFICATE_VALIDATION: 5.2.6 X.509
	// certificate validation.
	BasicBuildingBlockDefinition_X509_CERTIFICATE_VALIDATION BasicBuildingBlockDefinition = "X509_CERTIFICATE_VALIDATION"
	// BasicBuildingBlockDefinition_CRYPTOGRAPHIC_VERIFICATION: 5.2.7
	// Cryptographic verification.
	BasicBuildingBlockDefinition_CRYPTOGRAPHIC_VERIFICATION BasicBuildingBlockDefinition = "CRYPTOGRAPHIC_VERIFICATION"
	// BasicBuildingBlockDefinition_SIGNATURE_ACCEPTANCE_VALIDATION: 5.2.8
	// Signature Acceptance Validation (SAV).
	BasicBuildingBlockDefinition_SIGNATURE_ACCEPTANCE_VALIDATION BasicBuildingBlockDefinition = "SIGNATURE_ACCEPTANCE_VALIDATION"
	// BasicBuildingBlockDefinition_PAST_CERTIFICATE_VALIDATION: 5.6.2.1 Past
	// certificate validation.
	BasicBuildingBlockDefinition_PAST_CERTIFICATE_VALIDATION BasicBuildingBlockDefinition = "PAST_CERTIFICATE_VALIDATION"
	// BasicBuildingBlockDefinition_VALIDATION_TIME_SLIDING: 5.6.2.2 Validation
	// time sliding process.
	BasicBuildingBlockDefinition_VALIDATION_TIME_SLIDING BasicBuildingBlockDefinition = "VALIDATION_TIME_SLIDING"
	// BasicBuildingBlockDefinition_PAST_SIGNATURE_VALIDATION: 5.6.2.4 Past
	// signature validation building block.
	BasicBuildingBlockDefinition_PAST_SIGNATURE_VALIDATION BasicBuildingBlockDefinition = "PAST_SIGNATURE_VALIDATION"
)

// basicBuildingBlockDefinitionURIs holds the URI identifying each
// BasicBuildingBlock.
var basicBuildingBlockDefinitionURIs = map[BasicBuildingBlockDefinition]string{
	BasicBuildingBlockDefinition_FORMAT_CHECKING:                           "urn:cef:dss:bbb:formatChecking",
	BasicBuildingBlockDefinition_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE: "urn:cef:dss:bbb:identificationOfTheSigningCertificate",
	BasicBuildingBlockDefinition_VALIDATION_CONTEXT_INITIALIZATION:         "urn:cef:dss:bbb:validationContextInitialization",
	BasicBuildingBlockDefinition_REVOCATION_FRESHNESS_CHECKER:              "urn:cef:dss:bbb:revocationFreshnessChecker",
	BasicBuildingBlockDefinition_X509_CERTIFICATE_VALIDATION:               "urn:cef:dss:bbb:x509CertificateValidation",
	BasicBuildingBlockDefinition_CRYPTOGRAPHIC_VERIFICATION:                "urn:cef:dss:bbb:cryptographicVerification",
	BasicBuildingBlockDefinition_SIGNATURE_ACCEPTANCE_VALIDATION:           "urn:cef:dss:bbb:signatureAcceptanceValidation",
	BasicBuildingBlockDefinition_PAST_CERTIFICATE_VALIDATION:               "urn:cef:dss:bbb:pastCertificateValidation",
	BasicBuildingBlockDefinition_VALIDATION_TIME_SLIDING:                   "urn:cef:dss:bbb:validationTimeSliding",
	BasicBuildingBlockDefinition_PAST_SIGNATURE_VALIDATION:                 "urn:cef:dss:bbb:pastSignatureValidation",
}

// BasicBuildingBlockDefinitionValues returns all constants in declaration order.
func BasicBuildingBlockDefinitionValues() []BasicBuildingBlockDefinition {
	return []BasicBuildingBlockDefinition{
		BasicBuildingBlockDefinition_FORMAT_CHECKING,
		BasicBuildingBlockDefinition_IDENTIFICATION_OF_THE_SIGNING_CERTIFICATE,
		BasicBuildingBlockDefinition_VALIDATION_CONTEXT_INITIALIZATION,
		BasicBuildingBlockDefinition_REVOCATION_FRESHNESS_CHECKER,
		BasicBuildingBlockDefinition_X509_CERTIFICATE_VALIDATION,
		BasicBuildingBlockDefinition_CRYPTOGRAPHIC_VERIFICATION,
		BasicBuildingBlockDefinition_SIGNATURE_ACCEPTANCE_VALIDATION,
		BasicBuildingBlockDefinition_PAST_CERTIFICATE_VALIDATION,
		BasicBuildingBlockDefinition_VALIDATION_TIME_SLIDING,
		BasicBuildingBlockDefinition_PAST_SIGNATURE_VALIDATION,
	}
}

// URI returns the URI identifying the BasicBuildingBlock. Port of getUri().
func (b BasicBuildingBlockDefinition) URI() string {
	return basicBuildingBlockDefinitionURIs[b]
}

// compile-time assertion that the enum implements the marker interface Java
// declares on it.
var _ enumerations.UriBasedEnum = BasicBuildingBlockDefinition("")

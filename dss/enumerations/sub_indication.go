// Ported from dss-enumerations/.../SubIndication.java (DSS 6.5.RC1).
package enumerations

// SubIndication holds sub indication values. Source ETSI EN 319 102-1.
type SubIndication string

const (
	// SubIndicationFormatFailure: the signature is not conformant to one
	// of the base standards to the extent that the cryptographic
	// verification building block is unable to process it.
	SubIndicationFormatFailure SubIndication = "FORMAT_FAILURE"
	// SubIndicationHashFailure: the signature validation process results
	// into TOTAL-FAILED because at least one hash of a signed data
	// object(s) that has been included in the signing process does not
	// match the corresponding hash value in the signature.
	SubIndicationHashFailure SubIndication = "HASH_FAILURE"
	// SubIndicationSigCryptoFailure: the signature validation process
	// results into TOTAL-FAILED because the signature value in the
	// signature could not be verified using the signer's public key in the
	// signing certificate.
	SubIndicationSigCryptoFailure SubIndication = "SIG_CRYPTO_FAILURE"
	// SubIndicationRevoked: the signature validation process results into
	// TOTAL-FAILED because: the signing certificate has been revoked; and
	// there is proof that the signature has been created after the
	// revocation time.
	SubIndicationRevoked SubIndication = "REVOKED"
	// SubIndicationExpired: the signature validation process results into
	// TOTAL-FAILED because there is proof that the signature has been
	// created after the expiration date (notAfter) of the signing
	// certificate.
	SubIndicationExpired SubIndication = "EXPIRED"
	// SubIndicationNotYetValid: the signature validation process
	// results into TOTAL-FAILED because there is proof that the signature
	// was created before the issuance date (notBefore) of the signing
	// certificate.
	SubIndicationNotYetValid SubIndication = "NOT_YET_VALID"
	// SubIndicationSigConstraintsFailure: the signature validation
	// process results into INDETERMINATE because one or more attributes of
	// the signature do not match the validation constraints.
	SubIndicationSigConstraintsFailure SubIndication = "SIG_CONSTRAINTS_FAILURE"
	// SubIndicationChainConstraintsFailure: the signature validation
	// process results into INDETERMINATE because the certificate chain
	// used in the validation process does not match the validation
	// constraints related to the certificate.
	SubIndicationChainConstraintsFailure SubIndication = "CHAIN_CONSTRAINTS_FAILURE"
	// SubIndicationCertificateChainGeneralFailure: the signature
	// validation process results into INDETERMINATE because the set of
	// certificates available for chain validation produced an error for an
	// unspecified reason.
	SubIndicationCertificateChainGeneralFailure SubIndication = "CERTIFICATE_CHAIN_GENERAL_FAILURE"
	// SubIndicationCryptoConstraintsFailure: the signature validation
	// process results into INDETERMINATE because at least one of the
	// algorithms that have been used in material (e.g. the signature
	// value, a certificate...) involved in validating the signature, or
	// the size of a key used with such an algorithm, is below the
	// required cryptographic security level, and: this material was
	// produced after the time up to which this algorithm/key was
	// considered secure (if such a time is known); and the material is not
	// protected by a sufficiently strong time-stamp applied before the
	// time up to which the algorithm/key was considered secure (if such a
	// time is known).
	SubIndicationCryptoConstraintsFailure SubIndication = "CRYPTO_CONSTRAINTS_FAILURE"
	// SubIndicationPolicyProcessingError: the signature validation
	// process results into INDETERMINATE because a given formal policy
	// file could not be processed for any reason (e.g. not accessible, not
	// parseable, digest mismatch, etc.).
	SubIndicationPolicyProcessingError SubIndication = "POLICY_PROCESSING_ERROR"
	// SubIndicationSignaturePolicyNotAvailable: the signature
	// validation process results into INDETERMINATE because the
	// electronic document containing the details of the policy is not
	// available.
	SubIndicationSignaturePolicyNotAvailable SubIndication = "SIGNATURE_POLICY_NOT_AVAILABLE"
	// SubIndicationTimestampOrderFailure: the signature validation
	// process results into INDETERMINATE because some constraints on the
	// order of signature time-stamps and/or signed data object(s)
	// time-stamps are not respected.
	SubIndicationTimestampOrderFailure SubIndication = "TIMESTAMP_ORDER_FAILURE"
	// SubIndicationNoSigningCertificateFound: the signature validation
	// process results into INDETERMINATE because the signing certificate
	// cannot be identified.
	SubIndicationNoSigningCertificateFound SubIndication = "NO_SIGNING_CERTIFICATE_FOUND"
	// SubIndicationNoCertificateChainFound: the signature validation
	// process results into INDETERMINATE because no certificate chain has
	// been found for the identified signing certificate.
	SubIndicationNoCertificateChainFound SubIndication = "NO_CERTIFICATE_CHAIN_FOUND"
	// SubIndicationNoCertificateChainFoundNoPOE: the signature
	// validation process results into INDETERMINATE because no
	// certificate chain has been found for the identified signing
	// certificate due to the trust anchor not being trusted at the
	// validation date/time by the validation policy in use. However the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies before or after a time when the trust anchor was trusted
	// by the validation policy in use.
	SubIndicationNoCertificateChainFoundNoPOE SubIndication = "NO_CERTIFICATE_CHAIN_FOUND_NO_POE"
	// SubIndicationRevokedNoPOE: the signature validation process
	// results into INDETERMINATE because the signing certificate was
	// revoked at the validation date/time. However, the Signature
	// Validation Algorithm cannot ascertain that the signing time lies
	// before or after the revocation time.
	SubIndicationRevokedNoPOE SubIndication = "REVOKED_NO_POE"
	// SubIndicationRevokedCANoPOE: the signature validation process
	// results into INDETERMINATE because at least one certificate chain
	// was found but an intermediate CA certificate is revoked.
	SubIndicationRevokedCANoPOE SubIndication = "REVOKED_CA_NO_POE"
	// SubIndicationOutOfBoundsNotRevoked: the signature validation
	// process results into INDETERMINATE because the signing certificate
	// is expired or not yet valid at the validation date/time and the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies within the validity interval of the signing certificate.
	// The certificate is known not to be revoked.
	SubIndicationOutOfBoundsNotRevoked SubIndication = "OUT_OF_BOUNDS_NOT_REVOKED"
	// SubIndicationOutOfBoundsNoPOE: the signature validation process
	// results into INDETERMINATE because the signing certificate is
	// expired or not yet valid at the validation date/time and the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies within the validity interval of the signing certificate.
	SubIndicationOutOfBoundsNoPOE SubIndication = "OUT_OF_BOUNDS_NO_POE"
	// SubIndicationRevocationOutOfBoundsNoPOE: the signature
	// validation process results into INDETERMINATE because the signing
	// certificate of the revocation information of the signature signing
	// certificate is expired or not yet valid at the validation date/time
	// and the Signature Validation Algorithm cannot ascertain that the
	// revocation information issuance time lies within the validity
	// interval of the signing certificate of that revocation information.
	SubIndicationRevocationOutOfBoundsNoPOE SubIndication = "REVOCATION_OUT_OF_BOUNDS_NO_POE"
	// SubIndicationCryptoConstraintsFailureNoPOE: the signature
	// validation process results into INDETERMINATE because at least one
	// of the algorithms that have been used in objects (e.g. the signature
	// value, a certificate, etc.) involved in validating the signature, or
	// the size of a key used with such an algorithm, is below the
	// required cryptographic security level, and there is no proof that
	// this material was produced before the time up to which this
	// algorithm/key was considered secure.
	SubIndicationCryptoConstraintsFailureNoPOE SubIndication = "CRYPTO_CONSTRAINTS_FAILURE_NO_POE"
	// SubIndicationNoPOE: the signature validation process results into
	// INDETERMINATE because a proof of existence is missing to ascertain
	// that a signed object has been produced before some compromising
	// event (e.g. broken algorithm).
	SubIndicationNoPOE SubIndication = "NO_POE"
	// SubIndicationTryLater: the signature validation process results
	// into INDETERMINATE because not all constraints can be fulfilled
	// using available information. However, it may be possible to do so
	// using additional revocation information that will be available at a
	// later point of time.
	SubIndicationTryLater SubIndication = "TRY_LATER"
	// SubIndicationSignedDataNotFound: the signature validation
	// process results into INDETERMINATE because signed data cannot be
	// obtained.
	SubIndicationSignedDataNotFound SubIndication = "SIGNED_DATA_NOT_FOUND"
	// SubIndicationEAAConstraintsFailure: the EAA validation process
	// results into INDETERMINATE because there was some failure in the
	// validation constraints.
	SubIndicationEAAConstraintsFailure SubIndication = "EAA_CONSTRAINTS_FAILURE"
)

// subIndicationURI holds the VR URI for each constant.
var subIndicationURI = map[SubIndication]string{
	SubIndicationFormatFailure:                  "urn:etsi:019102:subindication:FORMAT_FAILURE",
	SubIndicationHashFailure:                    "urn:etsi:019102:subindication:HASH_FAILURE",
	SubIndicationSigCryptoFailure:               "urn:etsi:019102:subindication:SIG_CRYPTO_FAILURE",
	SubIndicationRevoked:                        "urn:etsi:019102:subindication:REVOKED",
	SubIndicationExpired:                        "urn:etsi:019102:subindication:EXPIRED",
	SubIndicationNotYetValid:                    "urn:etsi:019102:subindication:NOT_YET_VALID",
	SubIndicationSigConstraintsFailure:          "urn:etsi:019102:subindication:SIG_CONSTRAINTS_FAILURE",
	SubIndicationChainConstraintsFailure:        "urn:etsi:019102:subindication:CHAIN_CONSTRAINTS_FAILURE",
	SubIndicationCertificateChainGeneralFailure: "urn:etsi:019102:subindication:CERTIFICATE_CHAIN_GENERAL_FAILURE",
	SubIndicationCryptoConstraintsFailure:       "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE",
	SubIndicationPolicyProcessingError:          "urn:etsi:019102:subindication:POLICY_PROCESSING_ERROR",
	SubIndicationSignaturePolicyNotAvailable:    "urn:etsi:019102:subindication:SIGNATURE_POLICY_NOT_AVAILABLE",
	SubIndicationTimestampOrderFailure:          "urn:etsi:019102:subindication:TIMESTAMP_ORDER_FAILURE",
	SubIndicationNoSigningCertificateFound:      "urn:etsi:019102:subindication:NO_SIGNING_CERTIFICATE_FOUND",
	SubIndicationNoCertificateChainFound:        "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND",
	SubIndicationNoCertificateChainFoundNoPOE:   "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND_NO_POE",
	SubIndicationRevokedNoPOE:                   "urn:etsi:019102:subindication:REVOKED_NO_POE",
	SubIndicationRevokedCANoPOE:                 "urn:etsi:019102:subindication:REVOKED_CA_NO_POE",
	SubIndicationOutOfBoundsNotRevoked:          "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NOT_REVOKED",
	SubIndicationOutOfBoundsNoPOE:               "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NO_POE",
	SubIndicationRevocationOutOfBoundsNoPOE:     "urn:etsi:019102:subindication:REVOCATION_OUT_OF_BOUNDS_NO_POE",
	SubIndicationCryptoConstraintsFailureNoPOE:  "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE_NO_POE",
	SubIndicationNoPOE:                          "urn:etsi:019102:subindication:NO_POE",
	SubIndicationTryLater:                       "urn:etsi:019102:subindication:TRY_LATER",
	SubIndicationSignedDataNotFound:             "urn:etsi:019102:subindication:SIGNED_DATA_NOT_FOUND",
	SubIndicationEAAConstraintsFailure:          "urn:cef:dss:subindication:EAA_CONSTRAINTS_FAILURE",
}

// SubIndicationValues returns all constants in declaration order.
func SubIndicationValues() []SubIndication {
	return []SubIndication{
		SubIndicationFormatFailure,
		SubIndicationHashFailure,
		SubIndicationSigCryptoFailure,
		SubIndicationRevoked,
		SubIndicationExpired,
		SubIndicationNotYetValid,
		SubIndicationSigConstraintsFailure,
		SubIndicationChainConstraintsFailure,
		SubIndicationCertificateChainGeneralFailure,
		SubIndicationCryptoConstraintsFailure,
		SubIndicationPolicyProcessingError,
		SubIndicationSignaturePolicyNotAvailable,
		SubIndicationTimestampOrderFailure,
		SubIndicationNoSigningCertificateFound,
		SubIndicationNoCertificateChainFound,
		SubIndicationNoCertificateChainFoundNoPOE,
		SubIndicationRevokedNoPOE,
		SubIndicationRevokedCANoPOE,
		SubIndicationOutOfBoundsNotRevoked,
		SubIndicationOutOfBoundsNoPOE,
		SubIndicationRevocationOutOfBoundsNoPOE,
		SubIndicationCryptoConstraintsFailureNoPOE,
		SubIndicationNoPOE,
		SubIndicationTryLater,
		SubIndicationSignedDataNotFound,
		SubIndicationEAAConstraintsFailure,
	}
}

// URI returns the VR URI of the SubIndication.
func (s SubIndication) URI() string {
	return subIndicationURI[s]
}

// SubIndicationForName converts the given value to the related
// SubIndication. Returns "" (zero value) and no error if value is empty,
// mirroring Java's null return; an unknown value is an error, mirroring
// Java's valueOf IllegalArgumentException.
func SubIndicationForName(value string) (SubIndication, error) {
	if value == "" {
		return "", nil
	}
	return SubIndicationValueOf(value)
}

// SubIndicationValueOf returns the constant matching the given Java enum
// name.
func SubIndicationValueOf(name string) (SubIndication, error) {
	for _, v := range SubIndicationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &subIndicationInvalidValueError{name}
}

type subIndicationInvalidValueError struct {
	name string
}

func (e *subIndicationInvalidValueError) Error() string {
	return "no enum constant SubIndication." + e.name
}

// Compile-time assertion that SubIndication implements UriBasedEnum.
var _ UriBasedEnum = SubIndicationFormatFailure

// Ported from dss-enumerations/.../SubIndication.java (DSS 6.5.RC1).
package enumerations

// SubIndication holds sub indication values. Source ETSI EN 319 102-1.
type SubIndication string

const (
	// SubIndication_FORMAT_FAILURE: the signature is not conformant to one
	// of the base standards to the extent that the cryptographic
	// verification building block is unable to process it.
	SubIndication_FORMAT_FAILURE SubIndication = "FORMAT_FAILURE"
	// SubIndication_HASH_FAILURE: the signature validation process results
	// into TOTAL-FAILED because at least one hash of a signed data
	// object(s) that has been included in the signing process does not
	// match the corresponding hash value in the signature.
	SubIndication_HASH_FAILURE SubIndication = "HASH_FAILURE"
	// SubIndication_SIG_CRYPTO_FAILURE: the signature validation process
	// results into TOTAL-FAILED because the signature value in the
	// signature could not be verified using the signer's public key in the
	// signing certificate.
	SubIndication_SIG_CRYPTO_FAILURE SubIndication = "SIG_CRYPTO_FAILURE"
	// SubIndication_REVOKED: the signature validation process results into
	// TOTAL-FAILED because: the signing certificate has been revoked; and
	// there is proof that the signature has been created after the
	// revocation time.
	SubIndication_REVOKED SubIndication = "REVOKED"
	// SubIndication_EXPIRED: the signature validation process results into
	// TOTAL-FAILED because there is proof that the signature has been
	// created after the expiration date (notAfter) of the signing
	// certificate.
	SubIndication_EXPIRED SubIndication = "EXPIRED"
	// SubIndication_NOT_YET_VALID: the signature validation process
	// results into TOTAL-FAILED because there is proof that the signature
	// was created before the issuance date (notBefore) of the signing
	// certificate.
	SubIndication_NOT_YET_VALID SubIndication = "NOT_YET_VALID"
	// SubIndication_SIG_CONSTRAINTS_FAILURE: the signature validation
	// process results into INDETERMINATE because one or more attributes of
	// the signature do not match the validation constraints.
	SubIndication_SIG_CONSTRAINTS_FAILURE SubIndication = "SIG_CONSTRAINTS_FAILURE"
	// SubIndication_CHAIN_CONSTRAINTS_FAILURE: the signature validation
	// process results into INDETERMINATE because the certificate chain
	// used in the validation process does not match the validation
	// constraints related to the certificate.
	SubIndication_CHAIN_CONSTRAINTS_FAILURE SubIndication = "CHAIN_CONSTRAINTS_FAILURE"
	// SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE: the signature
	// validation process results into INDETERMINATE because the set of
	// certificates available for chain validation produced an error for an
	// unspecified reason.
	SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE SubIndication = "CERTIFICATE_CHAIN_GENERAL_FAILURE"
	// SubIndication_CRYPTO_CONSTRAINTS_FAILURE: the signature validation
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
	SubIndication_CRYPTO_CONSTRAINTS_FAILURE SubIndication = "CRYPTO_CONSTRAINTS_FAILURE"
	// SubIndication_POLICY_PROCESSING_ERROR: the signature validation
	// process results into INDETERMINATE because a given formal policy
	// file could not be processed for any reason (e.g. not accessible, not
	// parseable, digest mismatch, etc.).
	SubIndication_POLICY_PROCESSING_ERROR SubIndication = "POLICY_PROCESSING_ERROR"
	// SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE: the signature
	// validation process results into INDETERMINATE because the
	// electronic document containing the details of the policy is not
	// available.
	SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE SubIndication = "SIGNATURE_POLICY_NOT_AVAILABLE"
	// SubIndication_TIMESTAMP_ORDER_FAILURE: the signature validation
	// process results into INDETERMINATE because some constraints on the
	// order of signature time-stamps and/or signed data object(s)
	// time-stamps are not respected.
	SubIndication_TIMESTAMP_ORDER_FAILURE SubIndication = "TIMESTAMP_ORDER_FAILURE"
	// SubIndication_NO_SIGNING_CERTIFICATE_FOUND: the signature validation
	// process results into INDETERMINATE because the signing certificate
	// cannot be identified.
	SubIndication_NO_SIGNING_CERTIFICATE_FOUND SubIndication = "NO_SIGNING_CERTIFICATE_FOUND"
	// SubIndication_NO_CERTIFICATE_CHAIN_FOUND: the signature validation
	// process results into INDETERMINATE because no certificate chain has
	// been found for the identified signing certificate.
	SubIndication_NO_CERTIFICATE_CHAIN_FOUND SubIndication = "NO_CERTIFICATE_CHAIN_FOUND"
	// SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE: the signature
	// validation process results into INDETERMINATE because no
	// certificate chain has been found for the identified signing
	// certificate due to the trust anchor not being trusted at the
	// validation date/time by the validation policy in use. However the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies before or after a time when the trust anchor was trusted
	// by the validation policy in use.
	SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE SubIndication = "NO_CERTIFICATE_CHAIN_FOUND_NO_POE"
	// SubIndication_REVOKED_NO_POE: the signature validation process
	// results into INDETERMINATE because the signing certificate was
	// revoked at the validation date/time. However, the Signature
	// Validation Algorithm cannot ascertain that the signing time lies
	// before or after the revocation time.
	SubIndication_REVOKED_NO_POE SubIndication = "REVOKED_NO_POE"
	// SubIndication_REVOKED_CA_NO_POE: the signature validation process
	// results into INDETERMINATE because at least one certificate chain
	// was found but an intermediate CA certificate is revoked.
	SubIndication_REVOKED_CA_NO_POE SubIndication = "REVOKED_CA_NO_POE"
	// SubIndication_OUT_OF_BOUNDS_NOT_REVOKED: the signature validation
	// process results into INDETERMINATE because the signing certificate
	// is expired or not yet valid at the validation date/time and the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies within the validity interval of the signing certificate.
	// The certificate is known not to be revoked.
	SubIndication_OUT_OF_BOUNDS_NOT_REVOKED SubIndication = "OUT_OF_BOUNDS_NOT_REVOKED"
	// SubIndication_OUT_OF_BOUNDS_NO_POE: the signature validation process
	// results into INDETERMINATE because the signing certificate is
	// expired or not yet valid at the validation date/time and the
	// Signature Validation Algorithm cannot ascertain that the signing
	// time lies within the validity interval of the signing certificate.
	SubIndication_OUT_OF_BOUNDS_NO_POE SubIndication = "OUT_OF_BOUNDS_NO_POE"
	// SubIndication_REVOCATION_OUT_OF_BOUNDS_NO_POE: the signature
	// validation process results into INDETERMINATE because the signing
	// certificate of the revocation information of the signature signing
	// certificate is expired or not yet valid at the validation date/time
	// and the Signature Validation Algorithm cannot ascertain that the
	// revocation information issuance time lies within the validity
	// interval of the signing certificate of that revocation information.
	SubIndication_REVOCATION_OUT_OF_BOUNDS_NO_POE SubIndication = "REVOCATION_OUT_OF_BOUNDS_NO_POE"
	// SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE: the signature
	// validation process results into INDETERMINATE because at least one
	// of the algorithms that have been used in objects (e.g. the signature
	// value, a certificate, etc.) involved in validating the signature, or
	// the size of a key used with such an algorithm, is below the
	// required cryptographic security level, and there is no proof that
	// this material was produced before the time up to which this
	// algorithm/key was considered secure.
	SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE SubIndication = "CRYPTO_CONSTRAINTS_FAILURE_NO_POE"
	// SubIndication_NO_POE: the signature validation process results into
	// INDETERMINATE because a proof of existence is missing to ascertain
	// that a signed object has been produced before some compromising
	// event (e.g. broken algorithm).
	SubIndication_NO_POE SubIndication = "NO_POE"
	// SubIndication_TRY_LATER: the signature validation process results
	// into INDETERMINATE because not all constraints can be fulfilled
	// using available information. However, it may be possible to do so
	// using additional revocation information that will be available at a
	// later point of time.
	SubIndication_TRY_LATER SubIndication = "TRY_LATER"
	// SubIndication_SIGNED_DATA_NOT_FOUND: the signature validation
	// process results into INDETERMINATE because signed data cannot be
	// obtained.
	SubIndication_SIGNED_DATA_NOT_FOUND SubIndication = "SIGNED_DATA_NOT_FOUND"
	// SubIndication_EAA_CONSTRAINTS_FAILURE: the EAA validation process
	// results into INDETERMINATE because there was some failure in the
	// validation constraints.
	SubIndication_EAA_CONSTRAINTS_FAILURE SubIndication = "EAA_CONSTRAINTS_FAILURE"
)

// subIndicationURI holds the VR URI for each constant.
var subIndicationURI = map[SubIndication]string{
	SubIndication_FORMAT_FAILURE:                    "urn:etsi:019102:subindication:FORMAT_FAILURE",
	SubIndication_HASH_FAILURE:                      "urn:etsi:019102:subindication:HASH_FAILURE",
	SubIndication_SIG_CRYPTO_FAILURE:                "urn:etsi:019102:subindication:SIG_CRYPTO_FAILURE",
	SubIndication_REVOKED:                           "urn:etsi:019102:subindication:REVOKED",
	SubIndication_EXPIRED:                           "urn:etsi:019102:subindication:EXPIRED",
	SubIndication_NOT_YET_VALID:                     "urn:etsi:019102:subindication:NOT_YET_VALID",
	SubIndication_SIG_CONSTRAINTS_FAILURE:           "urn:etsi:019102:subindication:SIG_CONSTRAINTS_FAILURE",
	SubIndication_CHAIN_CONSTRAINTS_FAILURE:         "urn:etsi:019102:subindication:CHAIN_CONSTRAINTS_FAILURE",
	SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE: "urn:etsi:019102:subindication:CERTIFICATE_CHAIN_GENERAL_FAILURE",
	SubIndication_CRYPTO_CONSTRAINTS_FAILURE:        "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE",
	SubIndication_POLICY_PROCESSING_ERROR:           "urn:etsi:019102:subindication:POLICY_PROCESSING_ERROR",
	SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE:    "urn:etsi:019102:subindication:SIGNATURE_POLICY_NOT_AVAILABLE",
	SubIndication_TIMESTAMP_ORDER_FAILURE:           "urn:etsi:019102:subindication:TIMESTAMP_ORDER_FAILURE",
	SubIndication_NO_SIGNING_CERTIFICATE_FOUND:      "urn:etsi:019102:subindication:NO_SIGNING_CERTIFICATE_FOUND",
	SubIndication_NO_CERTIFICATE_CHAIN_FOUND:        "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND",
	SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE: "urn:etsi:019102:subindication:NO_CERTIFICATE_CHAIN_FOUND_NO_POE",
	SubIndication_REVOKED_NO_POE:                    "urn:etsi:019102:subindication:REVOKED_NO_POE",
	SubIndication_REVOKED_CA_NO_POE:                 "urn:etsi:019102:subindication:REVOKED_CA_NO_POE",
	SubIndication_OUT_OF_BOUNDS_NOT_REVOKED:         "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NOT_REVOKED",
	SubIndication_OUT_OF_BOUNDS_NO_POE:              "urn:etsi:019102:subindication:OUT_OF_BOUNDS_NO_POE",
	SubIndication_REVOCATION_OUT_OF_BOUNDS_NO_POE:   "urn:etsi:019102:subindication:REVOCATION_OUT_OF_BOUNDS_NO_POE",
	SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE: "urn:etsi:019102:subindication:CRYPTO_CONSTRAINTS_FAILURE_NO_POE",
	SubIndication_NO_POE:                            "urn:etsi:019102:subindication:NO_POE",
	SubIndication_TRY_LATER:                         "urn:etsi:019102:subindication:TRY_LATER",
	SubIndication_SIGNED_DATA_NOT_FOUND:             "urn:etsi:019102:subindication:SIGNED_DATA_NOT_FOUND",
	SubIndication_EAA_CONSTRAINTS_FAILURE:           "urn:cef:dss:subindication:EAA_CONSTRAINTS_FAILURE",
}

// SubIndicationValues returns all constants in declaration order.
func SubIndicationValues() []SubIndication {
	return []SubIndication{
		SubIndication_FORMAT_FAILURE,
		SubIndication_HASH_FAILURE,
		SubIndication_SIG_CRYPTO_FAILURE,
		SubIndication_REVOKED,
		SubIndication_EXPIRED,
		SubIndication_NOT_YET_VALID,
		SubIndication_SIG_CONSTRAINTS_FAILURE,
		SubIndication_CHAIN_CONSTRAINTS_FAILURE,
		SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE,
		SubIndication_CRYPTO_CONSTRAINTS_FAILURE,
		SubIndication_POLICY_PROCESSING_ERROR,
		SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE,
		SubIndication_TIMESTAMP_ORDER_FAILURE,
		SubIndication_NO_SIGNING_CERTIFICATE_FOUND,
		SubIndication_NO_CERTIFICATE_CHAIN_FOUND,
		SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE,
		SubIndication_REVOKED_NO_POE,
		SubIndication_REVOKED_CA_NO_POE,
		SubIndication_OUT_OF_BOUNDS_NOT_REVOKED,
		SubIndication_OUT_OF_BOUNDS_NO_POE,
		SubIndication_REVOCATION_OUT_OF_BOUNDS_NO_POE,
		SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE,
		SubIndication_NO_POE,
		SubIndication_TRY_LATER,
		SubIndication_SIGNED_DATA_NOT_FOUND,
		SubIndication_EAA_CONSTRAINTS_FAILURE,
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
var _ UriBasedEnum = SubIndication_FORMAT_FAILURE

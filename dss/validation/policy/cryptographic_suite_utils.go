// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/policy/CryptographicSuiteUtils.java (DSS 6.5.RC1).
//
// Out-of-manifest dependency (phase 8d AOV porter note - flagged per porter
// brief rather than silently added): this static Java utility class is not
// itself listed in the phase 8d AOV manifest, but every file in
// eu.europa.esig.dss.validation.process.bbb.aov.cc.checks (all manifest
// files) calls it directly, so it cannot be left unported without leaving the
// aov package uncompilable. Its Go home is dss/validation/policy per the
// package-mapping table in /PORTING_PLAN.md ("dss-validation (+ ...) ->
// validation, validation/diagnostic, validation/reports, validation/policy").
//
// Java's static methods become package-level functions; the two overloaded
// pairs Go cannot express by argument type alone
// (isSignatureAlgorithmWithKeySizeReliable/getExpirationDate, each with a
// String-keyLength and an Integer-keySize overload) are split into
// differently-named functions, keeping the String-taking entry point (the
// only one any upstream caller outside this file uses) at the overload's
// original name.
//
// The package is named "policy", matching the directory-name convention
// dss/model/policy already established; callers that also import
// dss/model/policy alias one of the two (see aov/cc's own imports).
package policy

import (
	"strconv"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
)

// Key size parameter names used by RSA/DSA algorithms. Ports the public
// static final MODULES_LENGTH_PARAMETER / PLENGTH_PARAMETER / QLENGTH_PARAMETER
// String constants.
const (
	modulesLengthParameter = "moduluslength"
	pLengthParameter       = "plength"
	qLengthParameter       = "qlength"
)

// IsSignatureAlgorithmReliable checks if the given SignatureAlgorithm is
// reliable (acceptable). Port of isSignatureAlgorithmReliable(CryptographicSuite, SignatureAlgorithm).
func IsSignatureAlgorithmReliable(cryptographicSuite policy.CryptographicSuite, signatureAlgorithm enumerations.SignatureAlgorithm) bool {
	if cryptographicSuite == nil {
		return true
	}
	return signatureAlgorithm != "" && containsSignatureAlgorithm(cryptographicSuite, signatureAlgorithm)
}

// IsDigestAlgorithmReliable checks if the given DigestAlgorithm is reliable
// (acceptable). Port of isDigestAlgorithmReliable(CryptographicSuite, DigestAlgorithm).
func IsDigestAlgorithmReliable(cryptographicSuite policy.CryptographicSuite, digestAlgorithm enumerations.DigestAlgorithm) bool {
	if cryptographicSuite == nil {
		return true
	}
	return digestAlgorithm != "" && containsDigestAlgorithm(cryptographicSuite, digestAlgorithm)
}

// IsSignatureAlgorithmWithKeySizeReliable checks if the keyLength for
// SignatureAlgorithm is reliable (acceptable). Port of the
// isSignatureAlgorithmWithKeySizeReliable(CryptographicSuite, SignatureAlgorithm, String) overload.
func IsSignatureAlgorithmWithKeySizeReliable(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyLength string) bool {
	return isSignatureAlgorithmWithKeySizeReliable(cryptographicSuite, signatureAlgorithm, parseKeySize(keyLength))
}

// isSignatureAlgorithmWithKeySizeReliable ports the
// isSignatureAlgorithmWithKeySizeReliable(CryptographicSuite, SignatureAlgorithm, Integer) overload.
func isSignatureAlgorithmWithKeySizeReliable(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keySize int) bool {
	if cryptographicSuite == nil {
		return true
	}

	if signatureAlgorithm != "" {
		acceptableSignatureAlgorithms := cryptographicSuite.AcceptableSignatureAlgorithms()
		evaluations, found := acceptableSignatureAlgorithms[signatureAlgorithm]
		if !found {
			return false
		}

		if utils.IsCollectionNotEmpty(evaluations) {
			for _, evaluation := range evaluations {
				if isEvaluationApplicable(signatureAlgorithm.EncryptionAlgorithm(), &keySize, evaluation) {
					return true
				}
			}
		} else {
			// no evaluations -> return true
			return true
		}
	}
	return false
}

// parseKeySize ports the private static parseKeySize(String).
func parseKeySize(keyLength string) int {
	if utils.IsStringDigits(keyLength) {
		size, err := strconv.Atoi(keyLength)
		if err == nil {
			return size
		}
	}
	return 0
}

// IsSignatureAlgorithmKeyLengthBigEnough verifies whether the given keyLength
// of the signatureAlgorithm is big enough. NOTE: this method only ensures
// that the key length is bigger than the minimal accepted key size. It does
// not consider the maximum requirements. Port of
// isSignatureAlgorithmKeyLengthBigEnough(CryptographicSuite, SignatureAlgorithm, String).
func IsSignatureAlgorithmKeyLengthBigEnough(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyLength string) bool {
	if cryptographicSuite == nil {
		return false
	}

	keySize := parseKeySize(keyLength)
	if signatureAlgorithm != "" {
		acceptableSignatureAlgorithms := cryptographicSuite.AcceptableSignatureAlgorithms()
		evaluations, found := acceptableSignatureAlgorithms[signatureAlgorithm]
		if !found {
			return false
		}

		if utils.IsCollectionNotEmpty(evaluations) {
			for _, evaluation := range evaluations {
				parameterList := evaluation.ParameterList()
				if utils.IsCollectionNotEmpty(parameterList) {
					for _, parameter := range parameterList {
						if isSupported(signatureAlgorithm.EncryptionAlgorithm(), parameter) {
							parameterMin := parameter.Min()
							if parameterMin == nil || *parameterMin < keySize {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// GetExpirationDateForSignatureAlgorithm gets an expiration date for the
// encryption algorithm with name signatureAlgorithm and keyLength. Returns
// nil if the expiration date is not defined for the algorithm. Port of the
// getExpirationDate(CryptographicSuite, SignatureAlgorithm, String) overload.
func GetExpirationDateForSignatureAlgorithm(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyLength string) *time.Time {
	return getExpirationDateForSignatureAlgorithm(cryptographicSuite, signatureAlgorithm, parseKeySize(keyLength))
}

// getExpirationDateForSignatureAlgorithm ports the
// getExpirationDate(CryptographicSuite, SignatureAlgorithm, Integer) overload.
func getExpirationDateForSignatureAlgorithm(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keySize int) *time.Time {
	if cryptographicSuite == nil {
		return nil
	}

	var expirationDate *time.Time
	if signatureAlgorithm != "" {
		acceptableSignatureAlgorithms := cryptographicSuite.AcceptableSignatureAlgorithms()
		evaluations := acceptableSignatureAlgorithms[signatureAlgorithm]
		if utils.IsCollectionNotEmpty(evaluations) {
			for _, evaluation := range evaluations {
				if isEvaluationApplicable(signatureAlgorithm.EncryptionAlgorithm(), &keySize, evaluation) {
					// return the last expiration date (at least one evaluation shall match)
					validityEnd := evaluation.ValidityEnd()
					if validityEnd == nil {
						return nil
					}
					if expirationDate == nil || validityEnd.After(*expirationDate) {
						expirationDate = validityEnd
					}
				}
			}
		}
	}
	return expirationDate
}

// GetExpirationDateForDigestAlgorithm gets an expiration date for the
// digestAlgorithm. Returns nil if the expiration date is not defined for the
// algorithm. Port of getExpirationDate(CryptographicSuite, DigestAlgorithm).
func GetExpirationDateForDigestAlgorithm(cryptographicSuite policy.CryptographicSuite,
	digestAlgorithm enumerations.DigestAlgorithm) *time.Time {
	if cryptographicSuite == nil {
		return nil
	}

	var expirationDate *time.Time
	if digestAlgorithm != "" {
		acceptableDigestAlgorithms := cryptographicSuite.AcceptableDigestAlgorithms()
		evaluations := acceptableDigestAlgorithms[digestAlgorithm]
		if utils.IsCollectionNotEmpty(evaluations) {
			for _, evaluation := range evaluations {
				// return the last expiration date (at least one evaluation shall match)
				validityEnd := evaluation.ValidityEnd()
				if validityEnd == nil {
					return nil
				}
				if expirationDate == nil || validityEnd.After(*expirationDate) {
					expirationDate = validityEnd
				}
			}
		}
	}
	return expirationDate
}

// IsDigestAlgorithmReliableAtTime verifies whether the digestAlgorithm is
// reliable at the validationTime. Port of
// isDigestAlgorithmReliableAtTime(CryptographicSuite, DigestAlgorithm, Date).
func IsDigestAlgorithmReliableAtTime(cryptographicSuite policy.CryptographicSuite,
	digestAlgorithm enumerations.DigestAlgorithm, validationTime time.Time) bool {
	evaluations := cryptographicSuite.AcceptableDigestAlgorithms()[digestAlgorithm]
	return reliableEvaluationExistsAtTimeSimple(evaluations, validationTime)
}

// GetReliableDigestAlgorithmsAtTime returns the list of reliable
// DigestAlgorithms according to the current validation policy at the given
// validation time. Port of getReliableDigestAlgorithmsAtTime(CryptographicSuite, Date):
// the Java HashSet becomes a deduplicated slice built in the map's key order,
// which - like every unordered Java Set the rest of this port encounters - has
// no defined iteration order to preserve; callers requiring determinism sort
// or otherwise post-process the result.
func GetReliableDigestAlgorithmsAtTime(cryptographicSuite policy.CryptographicSuite, validationTime time.Time) []enumerations.DigestAlgorithm {
	var reliableDigestAlgorithms []enumerations.DigestAlgorithm
	for digestAlgorithm, evaluations := range cryptographicSuite.AcceptableDigestAlgorithms() {
		if reliableEvaluationExistsAtTimeSimple(evaluations, validationTime) {
			reliableDigestAlgorithms = append(reliableDigestAlgorithms, digestAlgorithm)
		}
	}
	return reliableDigestAlgorithms
}

// GetReliableSignatureAlgorithmsWithMinimalKeyLengthAtTime returns the list of
// reliable SignatureAlgorithmWithMinKeySize according to the current
// validation policy and at the given time. Port of
// getReliableSignatureAlgorithmsWithMinimalKeyLengthAtTime(CryptographicSuite, Date);
// see GetReliableDigestAlgorithmsAtTime for the Set->slice iteration-order note.
func GetReliableSignatureAlgorithmsWithMinimalKeyLengthAtTime(cryptographicSuite policy.CryptographicSuite,
	validationTime time.Time) []*policy.SignatureAlgorithmWithMinKeySize {
	var result []*policy.SignatureAlgorithmWithMinKeySize
	for signatureAlgorithm, evaluations := range cryptographicSuite.AcceptableSignatureAlgorithms() {
		var minKeyLength *int
		if utils.IsCollectionNotEmpty(evaluations) {
			for _, evaluation := range evaluations {
				if isEvaluationApplicable(signatureAlgorithm.EncryptionAlgorithm(), nil, evaluation) {
					validityStart := evaluation.ValidityStart()
					validityEnd := evaluation.ValidityEnd()
					if (validityStart == nil || validityStart.Before(validationTime)) &&
						(validityEnd == nil || validityEnd.After(validationTime)) {
						keyLength := getMinKeyLength(signatureAlgorithm.EncryptionAlgorithm(), evaluation)
						if minKeyLength == nil || *minKeyLength > keyLength {
							minKeyLength = &keyLength
						}
					}
				}
			}
		}
		if minKeyLength != nil {
			result = append(result, policy.NewSignatureAlgorithmWithMinKeySize(signatureAlgorithm, *minKeyLength))
		}
	}
	return result
}

// IsSignatureAlgorithmReliableAtTime verifies whether the signatureAlgorithm
// with the keyLength is reliable at the validationTime. Port of
// isSignatureAlgorithmReliableAtTime(CryptographicSuite, SignatureAlgorithm, String, Date).
func IsSignatureAlgorithmReliableAtTime(cryptographicSuite policy.CryptographicSuite,
	signatureAlgorithm enumerations.SignatureAlgorithm, keyLength string, validationTime time.Time) bool {
	evaluations := cryptographicSuite.AcceptableSignatureAlgorithms()[signatureAlgorithm]
	keySize := parseKeySize(keyLength)
	return reliableEvaluationExistsAtTime(signatureAlgorithm.EncryptionAlgorithm(), &keySize, evaluations, validationTime)
}

// reliableEvaluationExistsAtTimeSimple ports the private
// reliableEvaluationExistsAtTime(Collection, Date) overload.
func reliableEvaluationExistsAtTimeSimple(evaluations []*policy.CryptographicSuiteEvaluation, validationTime time.Time) bool {
	return reliableEvaluationExistsAtTime("", nil, evaluations, validationTime)
}

// reliableEvaluationExistsAtTime ports the private
// reliableEvaluationExistsAtTime(EncryptionAlgorithm, Integer, Collection, Date)
// overload; an empty encryptionAlgorithm mirrors Java's null, under which every
// evaluation is applicable regardless of its parameters.
func reliableEvaluationExistsAtTime(encryptionAlgorithm enumerations.EncryptionAlgorithm, keySize *int,
	evaluations []*policy.CryptographicSuiteEvaluation, validationTime time.Time) bool {
	if utils.IsCollectionNotEmpty(evaluations) {
		for _, evaluation := range evaluations {
			if encryptionAlgorithm == "" || isEvaluationApplicable(encryptionAlgorithm, keySize, evaluation) {
				validityStart := evaluation.ValidityStart()
				validityEnd := evaluation.ValidityEnd()
				if (validityStart == nil || validityStart.Before(validationTime)) &&
					(validityEnd == nil || validityEnd.After(validationTime)) {
					return true
				}
			}
		}
	}
	return false
}

// isEvaluationApplicable ports the private static isEvaluationApplicable(EncryptionAlgorithm, Integer, CryptographicSuiteEvaluation).
func isEvaluationApplicable(algorithm enumerations.EncryptionAlgorithm, keySize *int, evaluation *policy.CryptographicSuiteEvaluation) bool {
	parameterList := evaluation.ParameterList()
	if utils.IsCollectionEmpty(parameterList) {
		return true
	}
	for _, parameter := range parameterList {
		if !isSupported(algorithm, parameter) {
			continue
		}
		if keySize != nil && ((parameter.Min() != nil && *keySize < *parameter.Min()) ||
			(parameter.Max() != nil && *keySize > *parameter.Max())) {
			continue
		}
		return true
	}
	return false
}

// getMinKeyLength ports the private static getMinKeyLength(EncryptionAlgorithm, CryptographicSuiteEvaluation).
func getMinKeyLength(algorithm enumerations.EncryptionAlgorithm, evaluation *policy.CryptographicSuiteEvaluation) int {
	var minKeyLength *int
	parameterList := evaluation.ParameterList()
	if utils.IsCollectionNotEmpty(parameterList) {
		for _, parameter := range parameterList {
			if !isSupported(algorithm, parameter) {
				continue
			}
			if parameter.Min() == nil {
				return 0
			}
			if minKeyLength == nil || *minKeyLength > *parameter.Min() {
				min := *parameter.Min()
				minKeyLength = &min
			}
		}
	}
	if minKeyLength != nil {
		return *minKeyLength
	}
	return 0
}

// isSupported ports the private static isSupported(EncryptionAlgorithm, CryptographicSuiteParameter).
// slf4j logging (the "Unknown Algorithms Parameter type" debug record) is
// dropped per PORTING.md.
func isSupported(encryptionAlgorithm enumerations.EncryptionAlgorithm, parameter *policy.CryptographicSuiteParameter) bool {
	parameterName := parameter.Name()
	// first come, first served logic
	switch parameterName {
	case modulesLengthParameter:
		if enumerations.EncryptionAlgorithm_RSA.IsEquivalent(encryptionAlgorithm) {
			return true
		}
	case pLengthParameter:
		if enumerations.EncryptionAlgorithm_DSA.IsEquivalent(encryptionAlgorithm) ||
			enumerations.EncryptionAlgorithm_ECDSA.IsEquivalent(encryptionAlgorithm) ||
			enumerations.EncryptionAlgorithm_EDDSA.IsEquivalent(encryptionAlgorithm) {
			return true
		}
	case qLengthParameter:
		// process silently (not supported)
	}
	return false
}

// containsSignatureAlgorithm reports whether the cryptographic suite's
// acceptable-signature-algorithms map has an entry for signatureAlgorithm.
// Port of the Map#containsKey call inlined in isSignatureAlgorithmReliable.
func containsSignatureAlgorithm(cryptographicSuite policy.CryptographicSuite, signatureAlgorithm enumerations.SignatureAlgorithm) bool {
	_, found := cryptographicSuite.AcceptableSignatureAlgorithms()[signatureAlgorithm]
	return found
}

// containsDigestAlgorithm reports whether the cryptographic suite's
// acceptable-digest-algorithms map has an entry for digestAlgorithm. Port of
// the Map#containsKey call inlined in isDigestAlgorithmReliable.
func containsDigestAlgorithm(cryptographicSuite policy.CryptographicSuite, digestAlgorithm enumerations.DigestAlgorithm) bool {
	_, found := cryptographicSuite.AcceptableDigestAlgorithms()[digestAlgorithm]
	return found
}

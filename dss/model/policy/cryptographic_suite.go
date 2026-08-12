// Ported from dss-model/.../model/policy/CryptographicSuite.java (DSS 6.5.RC1).
package policy

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
)

// CryptographicSuite contains policy for validation of cryptographic
// suites used on the signature or certificates.
//
// Java's Map<DigestAlgorithm/SignatureAlgorithm, Set<CryptographicSuiteEvaluation>>
// is represented as map[...][]* CryptographicSuiteEvaluation: the value
// type is not Go-map-key-comparable (it embeds slices), so the "Set" is
// represented as a slice deduplicated by CryptographicSuiteEvaluation#Equals
// by producers of this interface.
type CryptographicSuite interface {
	LevelRule

	// PolicyName gets a cryptographic suite name.
	PolicyName() string

	// AcceptableDigestAlgorithms gets a map of DigestAlgorithm's extracted
	// from a cryptographic suite and their corresponding
	// CryptographicSuiteEvaluation rules.
	AcceptableDigestAlgorithms() map[enumerations.DigestAlgorithm][]*CryptographicSuiteEvaluation

	// AcceptableSignatureAlgorithms gets a map of SignatureAlgorithm's
	// extracted from a cryptographic suite and their corresponding
	// CryptographicSuiteEvaluation rules.
	AcceptableSignatureAlgorithms() map[enumerations.SignatureAlgorithm][]*CryptographicSuiteEvaluation

	// SetLevel sets the global execution level for the cryptographic suite
	// constraints.
	SetLevel(level enumerations.Level)

	// AcceptableDigestAlgorithmsLevel returns a level constraint for
	// AcceptableDigestAlgo constraint if present, the global Level()
	// otherwise.
	AcceptableDigestAlgorithmsLevel() enumerations.Level

	// SetAcceptableDigestAlgorithmsLevel sets the execution level for the
	// acceptable digest algorithms check.
	SetAcceptableDigestAlgorithmsLevel(acceptableDigestAlgorithmsLevel enumerations.Level)

	// AcceptableSignatureAlgorithmsLevel returns a level constraint for
	// AcceptableEncryptionAlgo constraint if present, the global Level()
	// otherwise.
	AcceptableSignatureAlgorithmsLevel() enumerations.Level

	// SetAcceptableSignatureAlgorithmsLevel sets the execution level for
	// the acceptable signature algorithms check.
	SetAcceptableSignatureAlgorithmsLevel(acceptableSignatureAlgorithmsLevel enumerations.Level)

	// AcceptableSignatureAlgorithmsMiniKeySizeLevel returns a level
	// constraint for MiniPublicKeySize constraint if present, the global
	// Level() otherwise.
	AcceptableSignatureAlgorithmsMiniKeySizeLevel() enumerations.Level

	// SetAcceptableSignatureAlgorithmsMiniKeySizeLevel sets the execution
	// level for the acceptable minimum key sizes of signature algorithms
	// check.
	SetAcceptableSignatureAlgorithmsMiniKeySizeLevel(acceptableSignatureAlgorithmsMiniKeySizeLevel enumerations.Level)

	// AlgorithmsExpirationDateLevel returns a level constraint for
	// AlgoExpirationDate constraint if present, the global Level()
	// otherwise.
	AlgorithmsExpirationDateLevel() enumerations.Level

	// SetAlgorithmsExpirationDateLevel sets the execution level for
	// checking algorithms expiration.
	SetAlgorithmsExpirationDateLevel(algorithmsExpirationDateLevel enumerations.Level)

	// AlgorithmsExpirationDateAfterUpdateLevel returns a level constraint
	// for AlgoExpirationDate constraint if present, the global Level()
	// otherwise.
	AlgorithmsExpirationDateAfterUpdateLevel() enumerations.Level

	// SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel sets the execution
	// level for checking algorithms expiration after the validation
	// policy update.
	//
	// Default: Level_WARN (warning message is returned in case of
	// expiration of the used cryptographic constraints after the policy
	// update date).
	SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel(algorithmsExpirationTimeAfterPolicyUpdateLevel enumerations.Level)

	// CryptographicSuiteUpdateDate returns a date of the update of the
	// cryptographic suites within the validation policy.
	CryptographicSuiteUpdateDate() *time.Time
}

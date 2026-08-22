// Ported from dss-model/.../model/policy/crypto/CryptographicSuite19322.java (DSS 6.5.RC1).
package policy

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// CryptographicSuite19322 contains common methods for processing XML and
// JSON TS 119 322 schemas.
type CryptographicSuite19322 struct {
	// metadata is the metadata of the cryptographic suite.
	metadata *CryptographicSuiteMetadata

	// algorithmList is the list of applicable algorithm constraints.
	algorithmList []*CryptographicSuiteAlgorithm

	// globalLevel defines global execution level of the cryptographic
	// rules.
	globalLevel enumerations.Level

	// acceptableSignatureAlgorithmsLevel defines execution level of the
	// acceptability of signature algorithms check.
	acceptableSignatureAlgorithmsLevel enumerations.Level

	// acceptableSignatureAlgorithmsMinKeySizeLevel defines execution
	// level of the acceptability of the signature algorithms' key length
	// check.
	acceptableSignatureAlgorithmsMinKeySizeLevel enumerations.Level

	// acceptableDigestAlgorithmsLevel defines execution level of the
	// acceptability of digest algorithms check.
	acceptableDigestAlgorithmsLevel enumerations.Level

	// algorithmsExpirationDateLevel defines execution level of the
	// algorithms expiration check.
	algorithmsExpirationDateLevel enumerations.Level

	// algorithmsExpirationTimeAfterPolicyUpdateLevel defines execution
	// level of the algorithms expiration check with expiration occurred
	// after the update of the cryptographic suite.
	algorithmsExpirationTimeAfterPolicyUpdateLevel enumerations.Level

	// acceptableDigestAlgorithms is a cached map of acceptable digest
	// algorithms and their corresponding validation requirements.
	acceptableDigestAlgorithms map[enumerations.DigestAlgorithm][]*CryptographicSuiteEvaluation

	// acceptableSignatureAlgorithms is a cached map of acceptable
	// signature algorithms and their corresponding validation
	// requirements.
	acceptableSignatureAlgorithms map[enumerations.SignatureAlgorithm][]*CryptographicSuiteEvaluation
}

// NewCryptographicSuite19322 is the default constructor.
//
// Panics if metadata or algorithmList is nil, mirroring
// Objects.requireNonNull in the Java constructor (programmer error).
func NewCryptographicSuite19322(metadata *CryptographicSuiteMetadata, algorithmList []*CryptographicSuiteAlgorithm) *CryptographicSuite19322 {
	if metadata == nil {
		panic("metadata cannot be null!")
	}
	if algorithmList == nil {
		panic("algorithmList cannot be null!")
	}
	return &CryptographicSuite19322{
		metadata:      metadata,
		algorithmList: algorithmList,
		globalLevel:   enumerations.LevelFail,
		algorithmsExpirationTimeAfterPolicyUpdateLevel: enumerations.LevelWarn,
	}
}

// PolicyName gets a cryptographic suite name.
func (c *CryptographicSuite19322) PolicyName() string {
	return c.metadata.PolicyName()
}

// CryptographicSuiteUpdateDate returns a date of the update of the
// cryptographic suites within the validation policy.
func (c *CryptographicSuite19322) CryptographicSuiteUpdateDate() *time.Time {
	return c.metadata.PolicyIssueDate()
}

// Level gets the constraint execution level.
func (c *CryptographicSuite19322) Level() enumerations.Level {
	return c.globalLevel
}

// SetLevel sets the global execution level for the cryptographic suite
// constraints.
func (c *CryptographicSuite19322) SetLevel(level enumerations.Level) {
	c.globalLevel = level
}

// AcceptableDigestAlgorithmsLevel returns a level constraint for
// AcceptableDigestAlgo constraint if present, the global Level()
// otherwise.
func (c *CryptographicSuite19322) AcceptableDigestAlgorithmsLevel() enumerations.Level {
	return c.levelOrGlobal(c.acceptableDigestAlgorithmsLevel)
}

// SetAcceptableDigestAlgorithmsLevel sets the execution level for the
// acceptable digest algorithms check.
func (c *CryptographicSuite19322) SetAcceptableDigestAlgorithmsLevel(acceptableDigestAlgorithmsLevel enumerations.Level) {
	c.acceptableDigestAlgorithmsLevel = acceptableDigestAlgorithmsLevel
}

// AcceptableSignatureAlgorithmsLevel returns a level constraint for
// AcceptableEncryptionAlgo constraint if present, the global Level()
// otherwise.
func (c *CryptographicSuite19322) AcceptableSignatureAlgorithmsLevel() enumerations.Level {
	return c.levelOrGlobal(c.acceptableSignatureAlgorithmsLevel)
}

// SetAcceptableSignatureAlgorithmsLevel sets the execution level for the
// acceptable signature algorithms check.
func (c *CryptographicSuite19322) SetAcceptableSignatureAlgorithmsLevel(acceptableSignatureAlgorithmsLevel enumerations.Level) {
	c.acceptableSignatureAlgorithmsLevel = acceptableSignatureAlgorithmsLevel
}

// AcceptableSignatureAlgorithmsMiniKeySizeLevel returns a level
// constraint for MiniPublicKeySize constraint if present, the global
// Level() otherwise.
func (c *CryptographicSuite19322) AcceptableSignatureAlgorithmsMiniKeySizeLevel() enumerations.Level {
	return c.levelOrGlobal(c.acceptableSignatureAlgorithmsMinKeySizeLevel)
}

// SetAcceptableSignatureAlgorithmsMiniKeySizeLevel sets the execution
// level for the acceptable minimum key sizes of signature algorithms
// check.
func (c *CryptographicSuite19322) SetAcceptableSignatureAlgorithmsMiniKeySizeLevel(acceptableSignatureAlgorithmsMinKeySizeLevel enumerations.Level) {
	c.acceptableSignatureAlgorithmsMinKeySizeLevel = acceptableSignatureAlgorithmsMinKeySizeLevel
}

// AlgorithmsExpirationDateLevel returns a level constraint for
// AlgoExpirationDate constraint if present, the global Level()
// otherwise.
func (c *CryptographicSuite19322) AlgorithmsExpirationDateLevel() enumerations.Level {
	return c.levelOrGlobal(c.algorithmsExpirationDateLevel)
}

// SetAlgorithmsExpirationDateLevel sets the execution level for checking
// algorithms expiration.
func (c *CryptographicSuite19322) SetAlgorithmsExpirationDateLevel(algorithmsExpirationDateLevel enumerations.Level) {
	c.algorithmsExpirationDateLevel = algorithmsExpirationDateLevel
}

// AlgorithmsExpirationDateAfterUpdateLevel returns a level constraint for
// AlgoExpirationDate constraint if present, the global Level() otherwise.
func (c *CryptographicSuite19322) AlgorithmsExpirationDateAfterUpdateLevel() enumerations.Level {
	return c.algorithmsExpirationTimeAfterPolicyUpdateLevel
}

// SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel sets the execution
// level for checking algorithms expiration after the validation policy
// update.
func (c *CryptographicSuite19322) SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel(algorithmsExpirationTimeAfterPolicyUpdateLevel enumerations.Level) {
	c.algorithmsExpirationTimeAfterPolicyUpdateLevel = algorithmsExpirationTimeAfterPolicyUpdateLevel
}

// levelOrGlobal returns the global level in case of failure. Ports the
// private getLevel(Level) helper.
func (c *CryptographicSuite19322) levelOrGlobal(level enumerations.Level) enumerations.Level {
	if level != "" {
		return level
	}
	return c.globalLevel
}

// AcceptableDigestAlgorithms gets a map of DigestAlgorithm's extracted
// from a cryptographic suite and their corresponding
// CryptographicSuiteEvaluation rules.
func (c *CryptographicSuite19322) AcceptableDigestAlgorithms() map[enumerations.DigestAlgorithm][]*CryptographicSuiteEvaluation {
	if c.acceptableDigestAlgorithms == nil {
		acceptable := make(map[enumerations.DigestAlgorithm][]*CryptographicSuiteEvaluation)
		for _, algorithm := range c.algorithmList {
			digestAlgorithm := cryptographicSuite19322DigestAlgorithm(algorithm)
			if digestAlgorithm == "" {
				continue
			}
			acceptable[digestAlgorithm] = cryptographicSuite19322AddAllUnique(acceptable[digestAlgorithm], algorithm.EvaluationList())
		}
		c.acceptableDigestAlgorithms = acceptable
	}
	return c.acceptableDigestAlgorithms
}

// cryptographicSuite19322DigestAlgorithm ports the private
// getDigestAlgorithm(CryptographicSuiteAlgorithm) helper.
func cryptographicSuite19322DigestAlgorithm(algorithm *CryptographicSuiteAlgorithm) enumerations.DigestAlgorithm {
	if algorithm == nil {
		return ""
	}

	// NOTE: Name is not evaluated, it is not supposed to be
	// machine-processable.
	for _, oid := range algorithm.AlgorithmIdentifierOIDs() {
		// first come, first served policy
		if digestAlgorithm, err := enumerations.DigestAlgorithmForOID(oid); err == nil {
			return digestAlgorithm
		}
	}
	// optional
	for _, uri := range algorithm.AlgorithmIdentifierURIs() {
		if digestAlgorithm, err := enumerations.DigestAlgorithmForXML(uri); err == nil {
			return digestAlgorithm
		}
	}
	return ""
}

// AcceptableSignatureAlgorithms gets a map of SignatureAlgorithm's
// extracted from a cryptographic suite and their corresponding
// CryptographicSuiteEvaluation rules.
func (c *CryptographicSuite19322) AcceptableSignatureAlgorithms() map[enumerations.SignatureAlgorithm][]*CryptographicSuiteEvaluation {
	if c.acceptableSignatureAlgorithms == nil {
		acceptable := make(map[enumerations.SignatureAlgorithm][]*CryptographicSuiteEvaluation)

		// Step 1. Find all entries matching the SignatureAlgorithm
		// definition.
		for _, algorithm := range c.algorithmList {
			signatureAlgorithm := cryptographicSuite19322SignatureAlgorithm(algorithm)
			if signatureAlgorithm == "" {
				continue
			}
			acceptable[signatureAlgorithm] = cryptographicSuite19322AddAllUnique(acceptable[signatureAlgorithm], algorithm.EvaluationList())
		}

		// Step 2a. Extract supported digest algorithms for mapping.
		digestAlgorithmsMap := c.AcceptableDigestAlgorithms()

		// Step 2b. Extract supported encryption algorithms for mapping.
		// NOTE: we build a temp map to avoid conflict with
		// acceptableSignatureAlgorithms during map building.
		tempMap := make(map[enumerations.SignatureAlgorithm][]*CryptographicSuiteEvaluation)
		for _, algorithm := range c.algorithmList {
			encryptionAlgorithm := cryptographicSuite19322EncryptionAlgorithm(algorithm)
			if encryptionAlgorithm == "" {
				continue
			}
			for digestAlgorithm, digestEvaluations := range digestAlgorithmsMap {
				signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
				if signatureAlgorithm == "" {
					continue
				}
				if _, ok := acceptable[signatureAlgorithm]; ok {
					// if the SignatureAlgorithm is already present, prefer
					// the explicit definition
					continue
				}

				// Apply stricter requirements, if applicable.
				var finalEvaluationList []*CryptographicSuiteEvaluation
				for _, evaluation := range algorithm.EvaluationList() {
					for _, digestAlgoEvaluation := range digestEvaluations {
						evaluationCopy := CryptographicSuiteEvaluationCopy(evaluation)
						digestAlgoValidityStart := digestAlgoEvaluation.ValidityStart()
						if digestAlgoValidityStart != nil && (evaluationCopy.ValidityStart() == nil || digestAlgoValidityStart.After(*evaluationCopy.ValidityStart())) {
							evaluationCopy.SetValidityStart(digestAlgoValidityStart)
						}
						digestAlgoValidityEnd := digestAlgoEvaluation.ValidityEnd()
						if digestAlgoValidityEnd != nil && (evaluationCopy.ValidityEnd() == nil || digestAlgoValidityEnd.Before(*evaluationCopy.ValidityEnd())) {
							evaluationCopy.SetValidityEnd(digestAlgoValidityEnd)
						}

						// avoid dates misconfiguration
						if evaluationCopy.ValidityStart() == nil || evaluationCopy.ValidityEnd() == nil ||
							evaluationCopy.ValidityStart().Before(*evaluationCopy.ValidityEnd()) {
							finalEvaluationList = cryptographicSuite19322AddUnique(finalEvaluationList, evaluationCopy)
						}
					}
				}

				tempMap[signatureAlgorithm] = cryptographicSuite19322AddAllUnique(tempMap[signatureAlgorithm], finalEvaluationList)
			}
		}

		// Step 2c. Populate the main map.
		for signatureAlgorithm, evaluations := range tempMap {
			acceptable[signatureAlgorithm] = cryptographicSuite19322AddAllUnique(acceptable[signatureAlgorithm], evaluations)
		}

		c.acceptableSignatureAlgorithms = acceptable
	}
	return c.acceptableSignatureAlgorithms
}

// cryptographicSuite19322SignatureAlgorithm ports the private
// getSignatureAlgorithm(CryptographicSuiteAlgorithm) helper.
func cryptographicSuite19322SignatureAlgorithm(algorithm *CryptographicSuiteAlgorithm) enumerations.SignatureAlgorithm {
	if algorithm == nil {
		return ""
	}
	for _, oid := range algorithm.AlgorithmIdentifierOIDs() {
		signatureAlgorithm, err := enumerations.SignatureAlgorithmForOID(oid)
		if err == nil && signatureAlgorithm != "" {
			/*
			 * Here we check for a potential conflict with the
			 * EncryptionAlgorithm definition. If matching
			 * EncryptionAlgorithm is found as well, then we continue with
			 * the URIs check.
			 * Example: RSASSA-PSS using the same OID
			 * ("1.2.840.113549.1.1.10") as the RSA_SSA_PSS_SHA1_MGF1
			 * Signature Algorithm.
			 */
			encryptionAlgorithm := cryptographicSuite19322EncryptionAlgorithm(algorithm)
			if encryptionAlgorithm == "" {
				return signatureAlgorithm
			}
		}
	}
	// optional
	for _, uri := range algorithm.AlgorithmIdentifierURIs() {
		signatureAlgorithm, err := enumerations.SignatureAlgorithmForXML(uri)
		if err == nil && signatureAlgorithm != "" {
			return signatureAlgorithm
		}
	}
	return ""
}

// cryptographicSuite19322EncryptionAlgorithm ports the private
// getEncryptionAlgorithm(CryptographicSuiteAlgorithm) helper.
func cryptographicSuite19322EncryptionAlgorithm(algorithm *CryptographicSuiteAlgorithm) enumerations.EncryptionAlgorithm {
	if algorithm == nil {
		return ""
	}
	for _, oid := range algorithm.AlgorithmIdentifierOIDs() {
		// first come, first served policy
		if encryptionAlgorithm, err := enumerations.EncryptionAlgorithmForOID(oid); err == nil {
			return encryptionAlgorithm
		}
	}
	return ""
}

// cryptographicSuite19322AddUnique appends evaluation to list if an equal
// (by CryptographicSuiteEvaluation#Equals) entry is not already present,
// mirroring java.util.HashSet#add semantics for the Set<CryptographicSuiteEvaluation>
// Java uses (see the CryptographicSuite type doc for why a slice, not a
// map, backs the "Set" in this port).
func cryptographicSuite19322AddUnique(list []*CryptographicSuiteEvaluation, evaluation *CryptographicSuiteEvaluation) []*CryptographicSuiteEvaluation {
	for _, e := range list {
		if e.Equals(evaluation) {
			return list
		}
	}
	return append(list, evaluation)
}

// cryptographicSuite19322AddAllUnique appends every element of src into
// dst that is not already present (by CryptographicSuiteEvaluation#Equals).
func cryptographicSuite19322AddAllUnique(dst, src []*CryptographicSuiteEvaluation) []*CryptographicSuiteEvaluation {
	for _, e := range src {
		dst = cryptographicSuite19322AddUnique(dst, e)
	}
	return dst
}

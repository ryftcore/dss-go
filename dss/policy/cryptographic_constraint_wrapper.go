// Ported from dss-policy-jaxb/.../policy/CryptographicConstraintWrapper.java
// (DSS 6.5.RC1).
//
// Java's Set<CryptographicSuiteEvaluation> is represented as
// []*modelpolicy.CryptographicSuiteEvaluation deduplicated by
// CryptographicSuiteEvaluation#Equals - see dss/model/policy's
// CryptographicSuite doc comment, which documents and mandates this
// representation for every producer of the interface (including this one).
//
// Unlike TimeConstraintWrapper (RuleUtils) and other policy wrappers, this
// class does not delegate to DateUtils for date parsing: upstream gives it
// its own private getUsedDateFormat/getDate pair with different semantics
// (java.text.SimpleDateFormat's default lenient=true, forced UTC, a caught
// ParseException logged and swallowed rather than rethrown) - so this file
// carries its own local, similarly-scoped date-layout translation rather
// than reusing date_utils.go's (which mirrors DateUtils.parseDate's
// strict/throwing behavior instead). slf4j warnings are dropped throughout,
// per PORTING.md.
package policy

import (
	"strings"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// Parameter names used by RSA vs. DSA/ECDSA/EDDSA key-size constraints.
// Ports CryptographicConstraintWrapper#MODULES_LENGTH_PARAMETER and
// #PLENGTH_PARAMETER.
const (
	cryptographicConstraintWrapperModulesLengthParameter = "moduluslength"
	cryptographicConstraintWrapperPLengthParameter       = "plength"
)

// cryptographicConstraintWrapperDefaultDateFormat ports
// CryptographicConstraintWrapper#DEFAULT_DATE_FORMAT.
const cryptographicConstraintWrapperDefaultDateFormat = "yyyy-MM-dd"

// CryptographicConstraintWrapper wraps a CryptographicConstraint of the DSS
// JAXB validation policy implementation into a
// eu.europa.esig.dss.model.policy.CryptographicSuite.
type CryptographicConstraintWrapper struct {
	LevelConstraintWrapper
	constraint *jaxb.CryptographicConstraint

	acceptableDigestAlgorithms    map[enumerations.DigestAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation
	acceptableSignatureAlgorithms map[enumerations.SignatureAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation
}

var _ modelpolicy.CryptographicSuite = (*CryptographicConstraintWrapper)(nil)

// NewCryptographicConstraintWrapper is the default constructor. Passing nil
// ports the zero-argument CryptographicConstraintWrapper() constructor
// ("Constructor to create an empty instance of Cryptographic constraints"),
// which in Java calls super(null).
func NewCryptographicConstraintWrapper(constraint *jaxb.CryptographicConstraint) *CryptographicConstraintWrapper {
	var base *jaxb.LevelConstraint
	if constraint != nil {
		base = &constraint.LevelConstraint
	}
	return &CryptographicConstraintWrapper{
		LevelConstraintWrapper: LevelConstraintWrapper{constraint: base},
		constraint:             constraint,
	}
}

// PolicyName gets a cryptographic suite name. Ports
// CryptographicConstraintWrapper#getPolicyName.
func (w *CryptographicConstraintWrapper) PolicyName() string {
	return "DSS Cryptographic Constraint"
}

// AcceptableDigestAlgorithms gets a map of DigestAlgorithm's extracted from a
// cryptographic suite and their corresponding CryptographicSuiteEvaluation
// rules. Ports CryptographicConstraintWrapper#getAcceptableDigestAlgorithms.
func (w *CryptographicConstraintWrapper) AcceptableDigestAlgorithms() map[enumerations.DigestAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation {
	if w.acceptableDigestAlgorithms == nil {
		acceptable := make(map[enumerations.DigestAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation)

		if w.constraint != nil {
			// Step 1. Build evaluations based on acceptable algo list.
			if digestAlgo := w.constraint.AcceptableDigestAlgo; digestAlgo != nil {
				for _, algo := range digestAlgo.Algos {
					if digestAlgorithm, ok := cryptographicConstraintWrapperToDigestAlgorithm(algo.Value); ok {
						if _, exists := acceptable[digestAlgorithm]; !exists {
							acceptable[digestAlgorithm] = nil
						}
					}
				}
			}
			// Step 2. Build evaluations based on expiration dates (for
			// acceptable digest algos only).
			if algoExpirationDate := w.constraint.AlgoExpirationDate; algoExpirationDate != nil {
				layout := cryptographicConstraintWrapperUsedDateLayout(algoExpirationDate)
				for _, algo := range algoExpirationDate.Algos {
					digestAlgorithm, ok := cryptographicConstraintWrapperToDigestAlgorithm(algo.Value)
					if !ok {
						continue
					}
					if _, exists := acceptable[digestAlgorithm]; !exists {
						continue
					}
					evaluation := cryptographicConstraintWrapperBuildEvaluation("", algo, layout, nil)
					acceptable[digestAlgorithm] = cryptographicConstraintWrapperAddUniqueEvaluation(acceptable[digestAlgorithm], evaluation)
				}
			}
		}

		// Step 3. For acceptable digest algos without expiration date, add
		// an empty evaluation (does not expire).
		for digestAlgorithm, evaluationList := range acceptable {
			if len(evaluationList) == 0 {
				acceptable[digestAlgorithm] = append(evaluationList, modelpolicy.NewCryptographicSuiteEvaluation())
			}
		}

		w.acceptableDigestAlgorithms = acceptable
	}
	return w.acceptableDigestAlgorithms
}

// AcceptableSignatureAlgorithms gets a map of SignatureAlgorithm's extracted
// from a cryptographic suite and their corresponding
// CryptographicSuiteEvaluation rules. Ports
// CryptographicConstraintWrapper#getAcceptableSignatureAlgorithms.
func (w *CryptographicConstraintWrapper) AcceptableSignatureAlgorithms() map[enumerations.SignatureAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation {
	if w.acceptableSignatureAlgorithms == nil {
		acceptable := make(map[enumerations.SignatureAlgorithm][]*modelpolicy.CryptographicSuiteEvaluation)
		acceptableDigestAlgorithmsMap := w.AcceptableDigestAlgorithms()

		if w.constraint != nil {
			// Step 1. Build evaluations based on acceptable algo list.
			if encryptionAlgo := w.constraint.AcceptableEncryptionAlgo; encryptionAlgo != nil {
				for _, algo := range encryptionAlgo.Algos {
					encryptionAlgorithm, ok := cryptographicConstraintWrapperToEncryptionAlgorithm(algo.Value)
					if !ok {
						continue
					}
					for digestAlgorithm := range acceptableDigestAlgorithmsMap {
						signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
						if signatureAlgorithm == "" {
							continue
						}
						if _, exists := acceptable[signatureAlgorithm]; !exists {
							acceptable[signatureAlgorithm] = nil
						}
					}
				}
			}

			// Step 2a. Build evaluations based on expiration dates (for
			// acceptable signature algos only).
			if algoExpirationDate := w.constraint.AlgoExpirationDate; algoExpirationDate != nil {
				layout := cryptographicConstraintWrapperUsedDateLayout(algoExpirationDate)
				for _, algo := range algoExpirationDate.Algos {
					encryptionAlgorithm, ok := cryptographicConstraintWrapperToEncryptionAlgorithm(algo.Value)
					if !ok {
						continue
					}
					for digestAlgorithm, digestEvaluations := range acceptableDigestAlgorithmsMap {
						signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
						if signatureAlgorithm == "" {
							continue
						}
						if _, exists := acceptable[signatureAlgorithm]; !exists {
							continue
						}
						digestAlgoValidityEnd := cryptographicConstraintWrapperAlgorithmExpirationDate(digestEvaluations)
						evaluation := cryptographicConstraintWrapperBuildEvaluation(encryptionAlgorithm, algo, layout, digestAlgoValidityEnd)
						acceptable[signatureAlgorithm] = cryptographicConstraintWrapperAddUniqueEvaluation(acceptable[signatureAlgorithm], evaluation)
					}
				}
			}

			// Step 2b. Build evaluations based on min key sizes (for
			// acceptable signature algos only).
			if miniPublicKeySize := w.constraint.MiniPublicKeySize; miniPublicKeySize != nil {
				for _, algo := range miniPublicKeySize.Algos {
					encryptionAlgorithm, ok := cryptographicConstraintWrapperToEncryptionAlgorithm(algo.Value)
					if !ok {
						continue
					}
					for digestAlgorithm, digestEvaluations := range acceptableDigestAlgorithmsMap {
						signatureAlgorithm := enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
						if signatureAlgorithm == "" {
							continue
						}
						if _, exists := acceptable[signatureAlgorithm]; !exists {
							continue
						}
						digestAlgoValidityEnd := cryptographicConstraintWrapperAlgorithmExpirationDate(digestEvaluations)
						acceptable[signatureAlgorithm] = cryptographicConstraintWrapperFloorEvaluations(
							acceptable[signatureAlgorithm], encryptionAlgorithm, algo, digestAlgoValidityEnd)
					}
				}
			}
		}

		// Step 3. For acceptable signature algos without expiration date,
		// add an empty evaluation (does not expire).
		for signatureAlgorithm, evaluationList := range acceptable {
			if len(evaluationList) == 0 {
				evaluation := modelpolicy.NewCryptographicSuiteEvaluation()
				if digestAlgoEvaluations, ok := acceptableDigestAlgorithmsMap[signatureAlgorithm.DigestAlgorithm()]; ok {
					evaluation.SetValidityEnd(cryptographicConstraintWrapperAlgorithmExpirationDate(digestAlgoEvaluations))
				}
				acceptable[signatureAlgorithm] = append(evaluationList, evaluation)
			}
		}

		w.acceptableSignatureAlgorithms = acceptable
	}
	return w.acceptableSignatureAlgorithms
}

// cryptographicConstraintWrapperAlgorithmExpirationDate ports the private
// getAlgorithmExpirationDate(Set<CryptographicSuiteEvaluation>) helper: the
// latest validityEnd across evaluations, or nil if any evaluation has no
// validityEnd (an unbounded evaluation makes the whole set unbounded).
func cryptographicConstraintWrapperAlgorithmExpirationDate(evaluations []*modelpolicy.CryptographicSuiteEvaluation) *time.Time {
	if len(evaluations) == 0 {
		return nil
	}
	var expirationDate *time.Time
	for _, evaluation := range evaluations {
		validityEnd := evaluation.ValidityEnd()
		if validityEnd == nil {
			return nil
		}
		if expirationDate == nil || expirationDate.Before(*validityEnd) {
			expirationDate = validityEnd
		}
	}
	return expirationDate
}

// cryptographicConstraintWrapperFloorEvaluations ports the private
// getFloorEvaluations(Set<CryptographicSuiteEvaluation>, EncryptionAlgorithm,
// Algo, Date) helper.
func cryptographicConstraintWrapperFloorEvaluations(existingEvaluations []*modelpolicy.CryptographicSuiteEvaluation,
	encryptionAlgorithm enumerations.EncryptionAlgorithm, algo *jaxb.Algo, forcedValidityEnd *time.Time) []*modelpolicy.CryptographicSuiteEvaluation {
	if len(existingEvaluations) != 0 {
		minSize := algo.Size
		for _, evaluation := range existingEvaluations {
			for _, parameter := range evaluation.ParameterList() {
				if minSize != nil && (parameter.Min() == nil || *minSize > *parameter.Min()) {
					parameter.SetMin(minSize)
				}
			}
		}
		return existingEvaluations
	}
	evaluation := cryptographicConstraintWrapperBuildEvaluation(encryptionAlgorithm, algo, "", forcedValidityEnd)
	return append(existingEvaluations, evaluation)
}

// cryptographicConstraintWrapperBuildEvaluation ports the private
// buildEvaluation(Algo, SimpleDateFormat) and buildEvaluation(EncryptionAlgorithm,
// Algo, SimpleDateFormat, Date) overloads, unified: pass encryptionAlgorithm=""
// for the Algo-only overload, and layout="" to skip the validity-end
// computation entirely (Java only calls the Algo-only overload with a
// non-null SimpleDateFormat, so layout is always non-empty at that call site
// too - see AcceptableDigestAlgorithms above).
func cryptographicConstraintWrapperBuildEvaluation(encryptionAlgorithm enumerations.EncryptionAlgorithm, algo *jaxb.Algo, layout string, forcedValidityEnd *time.Time) *modelpolicy.CryptographicSuiteEvaluation {
	evaluation := modelpolicy.NewCryptographicSuiteEvaluation()
	evaluation.SetParameterList(cryptographicConstraintWrapperBuildParameters(encryptionAlgorithm, algo))
	if layout != "" {
		validityEnd := cryptographicConstraintWrapperGetDate(algo, layout)
		if validityEnd == nil || (forcedValidityEnd != nil && validityEnd.After(*forcedValidityEnd)) {
			validityEnd = forcedValidityEnd
		}
		evaluation.SetValidityEnd(validityEnd)
	}
	evaluation.SetAlgorithmUsage(cryptographicConstraintWrapperBuildUsages())
	return evaluation
}

// cryptographicConstraintWrapperBuildParameters ports the private
// buildParameters(EncryptionAlgorithm, Algo) helper.
func cryptographicConstraintWrapperBuildParameters(encryptionAlgorithm enumerations.EncryptionAlgorithm, algo *jaxb.Algo) []*modelpolicy.CryptographicSuiteParameter {
	var parameters []*modelpolicy.CryptographicSuiteParameter
	if algo.Size != nil {
		parameter := modelpolicy.NewCryptographicSuiteParameter()
		parameter.SetName(cryptographicConstraintWrapperGetParameterName(encryptionAlgorithm))
		parameter.SetMin(algo.Size)
		parameters = append(parameters, parameter)
	}
	return parameters
}

// cryptographicConstraintWrapperGetParameterName ports the private
// getParameterName(EncryptionAlgorithm) helper. Java's null return (no
// applicable parameter name) becomes the empty string, since
// CryptographicSuiteParameter#SetName takes a plain (non-pointer) string.
func cryptographicConstraintWrapperGetParameterName(encryptionAlgorithm enumerations.EncryptionAlgorithm) string {
	if encryptionAlgorithm == "" {
		return ""
	} else if enumerations.EncryptionAlgorithmRSA.IsEquivalent(encryptionAlgorithm) {
		return cryptographicConstraintWrapperModulesLengthParameter
	} else if enumerations.EncryptionAlgorithmDSA.IsEquivalent(encryptionAlgorithm) ||
		enumerations.EncryptionAlgorithmECDSA.IsEquivalent(encryptionAlgorithm) ||
		enumerations.EncryptionAlgorithmEDDSA.IsEquivalent(encryptionAlgorithm) {
		return cryptographicConstraintWrapperPLengthParameter
	}
	return ""
}

// cryptographicConstraintWrapperBuildUsages ports the private
// buildUsages() helper: only global usage is supported so far.
func cryptographicConstraintWrapperBuildUsages() []enumerations.CryptographicSuiteAlgorithmUsage {
	return nil
}

// cryptographicConstraintWrapperToDigestAlgorithm ports the private
// toDigestAlgorithm(String) helper; Java's caught IllegalArgumentException
// (silently continue) becomes ok=false.
func cryptographicConstraintWrapperToDigestAlgorithm(algorithmName string) (enumerations.DigestAlgorithm, bool) {
	digestAlgorithm, err := enumerations.DigestAlgorithmForName(algorithmName)
	if err != nil {
		return "", false
	}
	return digestAlgorithm, true
}

// cryptographicConstraintWrapperToEncryptionAlgorithm ports the private
// toEncryptionAlgorithm(String) helper; Java's caught IllegalArgumentException
// (silently continue) becomes ok=false.
func cryptographicConstraintWrapperToEncryptionAlgorithm(algorithmName string) (enumerations.EncryptionAlgorithm, bool) {
	encryptionAlgorithm, err := enumerations.EncryptionAlgorithmForName(algorithmName)
	if err != nil {
		return "", false
	}
	return encryptionAlgorithm, true
}

// cryptographicConstraintWrapperAddUniqueEvaluation appends evaluation to
// list if an equal (by CryptographicSuiteEvaluation#Equals) entry is not
// already present, mirroring java.util.HashSet#add semantics for the
// Set<CryptographicSuiteEvaluation> Java uses (see this file's header for
// why a slice, not a map, backs the "Set" in this port).
func cryptographicConstraintWrapperAddUniqueEvaluation(list []*modelpolicy.CryptographicSuiteEvaluation, evaluation *modelpolicy.CryptographicSuiteEvaluation) []*modelpolicy.CryptographicSuiteEvaluation {
	for _, e := range list {
		if e.Equals(evaluation) {
			return list
		}
	}
	return append(list, evaluation)
}

// cryptographicConstraintWrapperUsedDateLayout ports the private
// getUsedDateFormat(AlgoExpirationDate) helper, returning a Go reference-time
// layout instead of a java.text.SimpleDateFormat (the UTC timezone this
// applies is handled by cryptographicConstraintWrapperGetDate/ParseDate
// below, not by the layout string).
func cryptographicConstraintWrapperUsedDateLayout(expirations *jaxb.AlgoExpirationDate) string {
	format := cryptographicConstraintWrapperDefaultDateFormat
	if expirations.Format != nil {
		format = *expirations.Format
	}
	return cryptographicConstraintWrapperDateLayout(format)
}

// cryptographicConstraintWrapperDateLayout translates a
// java.text.SimpleDateFormat pattern into a Go reference-time layout. This
// mirrors javaSimpleDateFormatToGoLayout in date_utils.go but is kept as a
// separate, file-local copy: unlike DateUtils#parseDate (strict, throwing),
// CryptographicConstraintWrapper's own private getUsedDateFormat/getDate
// pair uses default (lenient) SimpleDateFormat parsing and swallows
// ParseException - genuinely different upstream code, not a shared
// dependency, so it is not factored into a cross-file helper (see
// PORTING.md's "no cross-file shared helpers" rule).
func cryptographicConstraintWrapperDateLayout(format string) string {
	letters := map[byte]string{
		'y': "2006",
		'M': "01",
		'd': "02",
		'H': "15",
		'm': "04",
		's': "05",
	}
	var out strings.Builder
	for i := 0; i < len(format); {
		c := format[i]
		layout, known := letters[c]
		if !known {
			out.WriteByte(c)
			i++
			continue
		}
		j := i
		for j < len(format) && format[j] == c {
			j++
		}
		out.WriteString(layout)
		i = j
	}
	return out.String()
}

// cryptographicConstraintWrapperGetDate ports the private
// getDate(Algo, SimpleDateFormat) helper.
func cryptographicConstraintWrapperGetDate(algo *jaxb.Algo, layout string) *time.Time {
	if algo == nil {
		return nil
	}
	return cryptographicConstraintWrapperParseDate(algo.Date, layout)
}

// cryptographicConstraintWrapperParseDate ports the private
// getDate(String, SimpleDateFormat) helper: a parse failure is logged (slf4j
// dropped, per PORTING.md) and swallowed, returning nil - unlike
// DateUtilsParseDate, which returns an error.
func cryptographicConstraintWrapperParseDate(dateString *string, layout string) *time.Time {
	if dateString == nil {
		return nil
	}
	t, err := time.Parse(layout, *dateString)
	if err != nil {
		return nil
	}
	utc := t.UTC()
	return &utc
}

// Level gets the constraint execution level. Ports
// CryptographicConstraintWrapper#getLevel (inherited from
// LevelConstraintWrapper - re-declared here only to document that
// CryptographicSuite embeds LevelRule; the implementation is the promoted
// LevelConstraintWrapper.Level method).

// SetLevel sets the global execution level for the cryptographic suite
// constraints. Ports CryptographicConstraintWrapper#setLevel.
func (w *CryptographicConstraintWrapper) SetLevel(level enumerations.Level) {
	if w.constraint != nil {
		w.constraint.Level = jaxb.LevelValue(level)
	}
}

// AcceptableSignatureAlgorithmsLevel returns a level constraint for
// AcceptableEncryptionAlgo constraint if present, the global Level()
// otherwise. Ports
// CryptographicConstraintWrapper#getAcceptableSignatureAlgorithmsLevel.
func (w *CryptographicConstraintWrapper) AcceptableSignatureAlgorithmsLevel() enumerations.Level {
	if w.constraint == nil {
		return ""
	}
	return w.getCryptographicLevel(cryptographicConstraintWrapperListAlgoBase(w.constraint.AcceptableEncryptionAlgo))
}

// SetAcceptableSignatureAlgorithmsLevel sets the execution level for the
// acceptable signature algorithms check. Ports
// CryptographicConstraintWrapper#setAcceptableSignatureAlgorithmsLevel.
func (w *CryptographicConstraintWrapper) SetAcceptableSignatureAlgorithmsLevel(level enumerations.Level) {
	if w.constraint != nil && w.constraint.AcceptableEncryptionAlgo != nil {
		w.constraint.AcceptableEncryptionAlgo.Level = jaxb.LevelValue(level)
	}
}

// AcceptableSignatureAlgorithmsMiniKeySizeLevel returns a level constraint
// for MiniPublicKeySize constraint if present, the global Level() otherwise.
// Ports CryptographicConstraintWrapper#getAcceptableSignatureAlgorithmsMiniKeySizeLevel.
func (w *CryptographicConstraintWrapper) AcceptableSignatureAlgorithmsMiniKeySizeLevel() enumerations.Level {
	if w.constraint == nil {
		return ""
	}
	return w.getCryptographicLevel(cryptographicConstraintWrapperListAlgoBase(w.constraint.MiniPublicKeySize))
}

// SetAcceptableSignatureAlgorithmsMiniKeySizeLevel sets the execution level
// for the acceptable minimum key sizes of signature algorithms check. Ports
// CryptographicConstraintWrapper#setAcceptableSignatureAlgorithmsMiniKeySizeLevel.
func (w *CryptographicConstraintWrapper) SetAcceptableSignatureAlgorithmsMiniKeySizeLevel(level enumerations.Level) {
	if w.constraint != nil && w.constraint.MiniPublicKeySize != nil {
		w.constraint.MiniPublicKeySize.Level = jaxb.LevelValue(level)
	}
}

// AcceptableDigestAlgorithmsLevel returns a level constraint for
// AcceptableDigestAlgo constraint if present, the global Level() otherwise.
// Ports CryptographicConstraintWrapper#getAcceptableDigestAlgorithmsLevel.
func (w *CryptographicConstraintWrapper) AcceptableDigestAlgorithmsLevel() enumerations.Level {
	if w.constraint == nil {
		return ""
	}
	return w.getCryptographicLevel(cryptographicConstraintWrapperListAlgoBase(w.constraint.AcceptableDigestAlgo))
}

// SetAcceptableDigestAlgorithmsLevel sets the execution level for the
// acceptable digest algorithms check. Ports
// CryptographicConstraintWrapper#setAcceptableDigestAlgorithmsLevel.
func (w *CryptographicConstraintWrapper) SetAcceptableDigestAlgorithmsLevel(level enumerations.Level) {
	if w.constraint != nil && w.constraint.AcceptableDigestAlgo != nil {
		w.constraint.AcceptableDigestAlgo.Level = jaxb.LevelValue(level)
	}
}

// AlgorithmsExpirationDateLevel returns a level constraint for
// AlgoExpirationDate constraint if present, the global Level() otherwise.
// Ports CryptographicConstraintWrapper#getAlgorithmsExpirationDateLevel.
func (w *CryptographicConstraintWrapper) AlgorithmsExpirationDateLevel() enumerations.Level {
	if w.constraint == nil {
		return ""
	}
	var base *jaxb.LevelConstraint
	if w.constraint.AlgoExpirationDate != nil {
		base = &w.constraint.AlgoExpirationDate.LevelConstraint
	}
	return w.getCryptographicLevel(base)
}

// SetAlgorithmsExpirationDateLevel sets the execution level for checking
// algorithms expiration. Ports
// CryptographicConstraintWrapper#setAlgorithmsExpirationDateLevel.
func (w *CryptographicConstraintWrapper) SetAlgorithmsExpirationDateLevel(level enumerations.Level) {
	if w.constraint != nil && w.constraint.AlgoExpirationDate != nil {
		w.constraint.AlgoExpirationDate.Level = jaxb.LevelValue(level)
	}
}

// AlgorithmsExpirationDateAfterUpdateLevel returns a level constraint for
// AlgoExpirationDate constraint if present, the global Level() otherwise.
// Ports CryptographicConstraintWrapper#getAlgorithmsExpirationDateAfterUpdateLevel.
func (w *CryptographicConstraintWrapper) AlgorithmsExpirationDateAfterUpdateLevel() enumerations.Level {
	if w.constraint == nil {
		return ""
	}
	aed := w.constraint.AlgoExpirationDate
	if aed != nil && aed.LevelAfterUpdate.Level() != "" {
		return aed.LevelAfterUpdate.Level()
	}
	var base *jaxb.LevelConstraint
	if aed != nil {
		base = &aed.LevelConstraint
	}
	return w.getCryptographicLevel(base)
}

// SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel sets the execution level
// for checking algorithms expiration after the validation policy update.
// Ports CryptographicConstraintWrapper#setAlgorithmsExpirationTimeAfterPolicyUpdateLevel.
func (w *CryptographicConstraintWrapper) SetAlgorithmsExpirationTimeAfterPolicyUpdateLevel(level enumerations.Level) {
	if w.constraint != nil && w.constraint.AlgoExpirationDate != nil {
		w.constraint.AlgoExpirationDate.LevelAfterUpdate = jaxb.LevelValue(level)
	}
}

// CryptographicSuiteUpdateDate returns a date of the update of the
// cryptographic suites within the validation policy. Ports
// CryptographicConstraintWrapper#getCryptographicSuiteUpdateDate.
func (w *CryptographicConstraintWrapper) CryptographicSuiteUpdateDate() *time.Time {
	if w.constraint != nil && w.constraint.AlgoExpirationDate != nil {
		layout := cryptographicConstraintWrapperUsedDateLayout(w.constraint.AlgoExpirationDate)
		return cryptographicConstraintWrapperParseDate(w.constraint.AlgoExpirationDate.UpdateDate, layout)
	}
	return nil
}

// cryptographicConstraintWrapperListAlgoBase extracts the embedded
// LevelConstraint of a *jaxb.ListAlgo, or nil if listAlgo is nil - the Go
// stand-in for Java's implicit upcast of ListAlgo (and its AlgoExpirationDate
// subtype) to LevelConstraint when calling getCryptographicLevel.
func cryptographicConstraintWrapperListAlgoBase(listAlgo *jaxb.ListAlgo) *jaxb.LevelConstraint {
	if listAlgo == nil {
		return nil
	}
	return &listAlgo.LevelConstraint
}

// getCryptographicLevel ports the private getCryptographicLevel(LevelConstraint)
// helper.
func (w *CryptographicConstraintWrapper) getCryptographicLevel(cryptoConstraint *jaxb.LevelConstraint) enumerations.Level {
	if cryptoConstraint != nil && cryptoConstraint.Level.Level() != "" {
		return cryptoConstraint.Level.Level()
	}
	// return global Level if target level is not present
	return w.Level()
}

// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/DigestAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Enables validation of a digest matcher chain. T is the validation token
// wrapper type.
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// digestMatcherTypesToIgnore contains a list of digest matcher types to be
// ignored in case of unknown digest algorithm. Port of the private static
// List<DigestMatcherType> digestMatcherTypesToIgnore.
var digestMatcherTypesToIgnore = []enumerations.DigestMatcherType{
	enumerations.DigestMatcherTypeCounterSignedSignatureValue,
	enumerations.DigestMatcherTypeEvidenceRecordArchiveObject,
	enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStamp,
	enumerations.DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence,
	enumerations.DigestMatcherTypeEvidenceRecordMasterSignature,
	enumerations.DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim,
}

// DigestAlgorithmObsolescenceValidation is the Go form of the abstract Java
// class DigestAlgorithmObsolescenceValidation<T>. Java's constructor is a pure
// pass-through to the parent, so no additional wiring method is needed here:
// leaf constructors call InitAlgorithmObsolescenceValidation directly, reached
// through promotion.
type DigestAlgorithmObsolescenceValidation[T any] struct {
	AlgorithmObsolescenceValidation[T]
}

// buildDigestMatchersValidationChain builds a chain of crypto checks to be
// executed on a signature's references (digest matchers). Port of
// buildDigestMatchersValidationChain(ChainItem, List, String).
func (c *DigestAlgorithmObsolescenceValidation[T]) buildDigestMatchersValidationChain(item process.ChainItem[*jaxb.XmlAOV],
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher, tokenId string) process.ChainItem[*jaxb.XmlAOV] {
	if utils.IsCollectionEmpty(digestMatchers) {
		return item
	}

	var cryptographicValidation *jaxb.XmlCryptographicValidation

	digestMatchersToProcess := getDigestMatchersToProcess(digestMatchers)
	usedDigestAlgorithms := getUsedDigestAlgorithms(digestMatchersToProcess)
	usedPositions := getUsedPositions(digestMatchersToProcess)
	for _, digestAlgorithm := range usedDigestAlgorithms {
		for _, position := range usedPositions {
			digestMatchersGroup := getDigestMatchersByAlgorithmAndPosition(digestMatchersToProcess, digestAlgorithm, position)
			if utils.IsCollectionNotEmpty(digestMatchersGroup) {
				dac := NewDigestAlgorithmCryptographicChecker(c.I18nProvider, digestAlgorithm, c.validationDate, position, c.cryptographicSuite)
				dacResult := dac.Execute()

				if item == nil {
					item = c.digestAlgorithmCheckResult(digestMatchersGroup, dacResult, c.cryptographicSuite)
					c.FirstItem = item
				} else {
					item = item.SetNextItem(c.digestAlgorithmCheckResult(digestMatchersGroup, dacResult, c.cryptographicSuite))
				}

				if cryptographicValidation == nil || (c.isValid(cryptographicValidation) &&
					enumerations.IndicationPassed != dacResult.Conclusion.Indication.Indication()) {
					cryptographicValidation = dacResult.CryptographicValidation
					cryptographicValidation.ConcernedMaterialDescription = c.materialDescription(digestMatchersGroup)
				}
			}
		}
	}

	c.digestMatchersCryptographicValidation = cryptographicValidation
	if c.digestMatchersCryptographicValidation != nil {
		c.digestMatchersCryptographicValidation.TokenId = &tokenId
	}

	return item
}

// getDigestMatchersToProcess omits digest matchers that were created for
// validation purposes but are not originally present in a signature. Port of
// the private getDigestMatchersToProcess(List).
func getDigestMatchersToProcess(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) []*diagnosticjaxb.XmlDigestMatcher {
	var digestMatchersToProcess []*diagnosticjaxb.XmlDigestMatcher
	for _, digestMatcher := range digestMatchers {
		if digestMatcher.DigestMethod != nil || !containsDigestMatcherType(digestMatcherTypesToIgnore, digestMatcherType(digestMatcher)) {
			digestMatchersToProcess = append(digestMatchersToProcess, digestMatcher)
		}
	}
	if utils.IsCollectionEmpty(digestMatchersToProcess) {
		return digestMatchers // return original values if no matching entries found
	}
	return digestMatchersToProcess
}

// getUsedDigestAlgorithms ports the private getUsedDigestAlgorithms(List): the
// Java LinkedHashSet becomes a deduplicated slice, preserving insertion order
// for deterministic output.
func getUsedDigestAlgorithms(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) []enumerations.DigestAlgorithm {
	seen := make(map[enumerations.DigestAlgorithm]struct{})
	var result []enumerations.DigestAlgorithm
	for _, digestMatcher := range digestMatchers {
		algorithm := digestMatcherDigestAlgorithm(digestMatcher)
		if _, ok := seen[algorithm]; !ok {
			seen[algorithm] = struct{}{}
			result = append(result, algorithm)
		}
	}
	return result
}

// getUsedPositions ports the private getUsedPositions(List): the Java
// LinkedHashSet becomes a deduplicated slice, preserving insertion order for
// deterministic output.
func getUsedPositions(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) []i18n.MessageTag {
	seen := make(map[i18n.MessageTag]struct{})
	var result []i18n.MessageTag
	for _, digestMatcher := range digestMatchers {
		position, err := process.GetDigestMatcherCryptoPosition(digestMatcher)
		if err != nil {
			panic(err)
		}
		if _, ok := seen[position]; !ok {
			seen[position] = struct{}{}
			result = append(result, position)
		}
	}
	return result
}

// getDigestMatchersByAlgorithmAndPosition ports the private
// getDigestMatchersByAlgorithmAndPosition(List, DigestAlgorithm, MessageTag).
func getDigestMatchersByAlgorithmAndPosition(digestMatchers []*diagnosticjaxb.XmlDigestMatcher,
	digestAlgorithm enumerations.DigestAlgorithm, position i18n.MessageTag) []*diagnosticjaxb.XmlDigestMatcher {
	if position == "" {
		return nil
	}
	var result []*diagnosticjaxb.XmlDigestMatcher
	for _, d := range digestMatchers {
		matcherPosition, err := process.GetDigestMatcherCryptoPosition(d)
		if err != nil {
			panic(err)
		}
		if digestAlgorithm == digestMatcherDigestAlgorithm(d) && position == matcherPosition &&
			// COUNTER_SIGNED_SIGNATURE_VALUE is an internal variable
			enumerations.DigestMatcherTypeCounterSignedSignatureValue != digestMatcherType(d) {
			result = append(result, d)
		}
	}
	return result
}

// digestAlgorithmCheckResult ports the private digestAlgorithmCheckResult(List, XmlCC, CryptographicSuite).
func (c *DigestAlgorithmObsolescenceValidation[T]) digestAlgorithmCheckResult(digestMatchers []*diagnosticjaxb.XmlDigestMatcher,
	ccResult *jaxb.XmlCC, constraint policy.CryptographicSuite) process.ChainItem[*jaxb.XmlAOV] {
	// Java calls the getDigestMatcherCryptoPosition(Collection) overload, which
	// resolves the PLURAL MessageTag when the group holds more than one matcher.
	position, err := process.GetDigestMatchersCryptoPosition(digestMatchers)
	if err != nil {
		panic(err)
	}
	return NewDigestMatcherCryptographicCheckerResultCheck(c.I18nProvider, c.Result, c.validationDate, position,
		getReferenceNames(digestMatchers), ccResult, constraint)
}

// materialDescription ports the private getMaterialDescription(List).
func (c *DigestAlgorithmObsolescenceValidation[T]) materialDescription(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) *string {
	referenceNames := getReferenceNames(digestMatchers)
	var message string
	if utils.IsCollectionNotEmpty(referenceNames) {
		message = c.I18nProvider.GetMessage(i18n.MessageTagACCMDescWithName, c.position, utils.JoinStrings(referenceNames, ", "))
	} else {
		message = c.I18nProvider.GetMessage(c.position)
	}
	return &message
}

// getReferenceNames ports the private getReferenceNames(List).
func getReferenceNames(digestMatchers []*diagnosticjaxb.XmlDigestMatcher) []string {
	var referenceNames []string
	for _, d := range digestMatchers {
		var name *string
		if d.Id != nil {
			name = d.Id
		} else if d.Uri != nil {
			name = d.Uri
		} else if d.DocumentName != nil {
			name = d.DocumentName
		}
		if name != nil {
			referenceNames = append(referenceNames, *name)
		}
	}
	return referenceNames
}

// digestMatcherDigestAlgorithm reads XmlDigestMatcher#getDigestMethod(): the
// generated member is a pointer, whose nil is Java's null.
func digestMatcherDigestAlgorithm(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestAlgorithm {
	if digestMatcher.DigestMethod == nil {
		return ""
	}
	return digestMatcher.DigestMethod.DigestAlgorithm()
}

// digestMatcherType reads XmlDigestMatcher#getType(): the generated member is
// a pointer, whose nil is Java's null.
func digestMatcherType(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}

// containsDigestMatcherType ports List#contains(Object) over a
// DigestMatcherType slice.
func containsDigestMatcherType(types []enumerations.DigestMatcherType, t enumerations.DigestMatcherType) bool {
	for _, v := range types {
		if v == t {
			return true
		}
	}
	return false
}

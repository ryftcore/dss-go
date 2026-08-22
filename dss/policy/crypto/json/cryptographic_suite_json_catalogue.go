// Ported from dss-policy-crypto-json/.../json/CryptographicSuiteJsonCatalogue.java (DSS 6.5.RC1).
package cryptojson

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// cryptographicSuiteJsonCatalogueDefaultVersion is the default value of the
// "version" parameter. Ports CryptographicSuiteJsonCatalogue#DEFAULT_VERSION.
const cryptographicSuiteJsonCatalogueDefaultVersion = "1"

// cryptographicSuiteJsonCatalogueDefaultLang is the default value of the
// "lang" parameter. Ports CryptographicSuiteJsonCatalogue#DEFAULT_LANG.
const cryptographicSuiteJsonCatalogueDefaultLang = "en"

// NewCryptographicSuiteJsonCatalogue parses an ETSI TS 119 322 JSON
// cryptographic suite catalog document and returns the extracted values.
// Ports CryptographicSuiteJsonCatalogue's constructor, collapsed with the
// class body itself: Java's protected abstract
// buildMetadata()/buildAlgorithmList() (see dss/model/policy's
// CryptographicSuiteCatalogue doc comment on why they become
// constructor-supplied closures rather than a subclass with virtual
// dispatch) close directly over securitySuitabilityPolicy.
//
// Java's Objects.requireNonNull(securitySuitabilityPolicy, "...") becomes a
// panic with the same message, per PORTING.md's requireNonNull rule.
func newCryptographicSuiteJsonCatalogue(securitySuitabilityPolicy jsonObject) *modelpolicy.CryptographicSuiteCatalogue {
	if securitySuitabilityPolicy == nil {
		panic("securitySuitabilityPolicy cannot be null!")
	}
	return modelpolicy.NewCryptographicSuiteCatalogue(
		func() *modelpolicy.CryptographicSuiteMetadata {
			return buildJSONMetadata(securitySuitabilityPolicy)
		},
		func() []*modelpolicy.CryptographicSuiteAlgorithm {
			return buildJSONAlgorithmList(securitySuitabilityPolicy)
		},
	)
}

// buildJSONMetadata ports the protected buildMetadata() method.
func buildJSONMetadata(securitySuitabilityPolicy jsonObject) *modelpolicy.CryptographicSuiteMetadata {
	metadata := modelpolicy.NewCryptographicSuiteMetadata()

	if policyName := securitySuitabilityPolicy.asObject(jsonConstraintPolicyName); policyName != nil {
		metadata.SetPolicyName(policyName.asString(jsonConstraintNameC))
		metadata.SetPolicyOID(policyName.asString(jsonConstraintObjectIdentifier))
		metadata.SetPolicyURI(policyName.asString(jsonConstraintURI))
	}

	if publisher := securitySuitabilityPolicy.asObject(jsonConstraintPublisher); publisher != nil {
		metadata.SetPublisherName(publisher.asString(jsonConstraintNameC))
		metadata.SetPublisherAddress(publisher.asString(jsonConstraintAddress))
		metadata.SetPublisherURI(publisher.asString(jsonConstraintURI))
	}

	metadata.SetPolicyIssueDate(asDateTime(securitySuitabilityPolicy, jsonConstraintPolicyIssueDate))
	metadata.SetNextUpdate(asDateTime(securitySuitabilityPolicy, jsonConstraintNextUpdate))
	metadata.SetUsage(securitySuitabilityPolicy.asString(jsonConstraintUsage))

	metadata.SetVersion(jsonCatalogueVersion(securitySuitabilityPolicy))
	metadata.SetLang(jsonCatalogueLang(securitySuitabilityPolicy))
	metadata.SetId(securitySuitabilityPolicy.asString(jsonConstraintID))

	return metadata
}

// buildJSONAlgorithmList ports the protected buildAlgorithmList() method.
func buildJSONAlgorithmList(securitySuitabilityPolicy jsonObject) []*modelpolicy.CryptographicSuiteAlgorithm {
	var algorithmList []*modelpolicy.CryptographicSuiteAlgorithm
	for _, algorithmType := range securitySuitabilityPolicy.asObjectList(jsonConstraintAlgorithm) {
		if algorithm := buildJSONAlgorithm(algorithmType); algorithm != nil {
			algorithmList = append(algorithmList, algorithm)
		}
	}
	return algorithmList
}

// buildJSONAlgorithm ports the private buildAlgorithm(JsonObjectWrapper)
// method.
//
// Java catches any exception raised while processing a single algorithm
// entry, logs it (slf4j dropped, per PORTING.md), and skips the entry
// (returning null). The accessors here never panic on malformed JSON
// shapes (see json_object.go: every getter degrades to a zero value rather
// than erroring), so there is nothing to recover from here - the function
// simply cannot fail, unlike its Java counterpart.
func buildJSONAlgorithm(algorithmType jsonObject) *modelpolicy.CryptographicSuiteAlgorithm {
	algorithm := modelpolicy.NewCryptographicSuiteAlgorithm()

	if algorithmIdentifier := algorithmType.asObject(jsonConstraintAlgorithmIdentifier); algorithmIdentifier != nil {
		algorithm.SetAlgorithmIdentifierName(algorithmIdentifier.asString(jsonConstraintNameC))
		algorithm.SetAlgorithmIdentifierOIDs(algorithmIdentifierOIDs(algorithmIdentifier))
		algorithm.SetAlgorithmIdentifierURIs(algorithmIdentifierURIs(algorithmIdentifier))
	}

	algorithm.SetEvaluationList(buildJSONEvaluationList(algorithmType.asObjectList(jsonConstraintEvaluation)))
	algorithm.SetInformationTextList(jsonInformationText(algorithmType))

	return algorithm
}

// algorithmIdentifierOIDs ports the private
// getAlgorithmIdentifierOIDs(JsonObjectWrapper) helper.
func algorithmIdentifierOIDs(algorithmIdentifier jsonObject) []string {
	if algorithmOID := algorithmIdentifier.asString(jsonConstraintObjectIdentifier); algorithmOID != "" {
		return []string{algorithmOID}
	}
	return nil
}

// algorithmIdentifierURIs ports the private
// getAlgorithmIdentifierURIs(JsonObjectWrapper) helper.
func algorithmIdentifierURIs(algorithmIdentifier jsonObject) []string {
	if algorithmURI := algorithmIdentifier.asString(jsonConstraintURI); algorithmURI != "" {
		return []string{algorithmURI}
	}
	return nil
}

// buildJSONEvaluationList ports the private
// buildEvaluationList(List<JsonObjectWrapper>) helper.
func buildJSONEvaluationList(evaluations []jsonObject) []*modelpolicy.CryptographicSuiteEvaluation {
	var evaluationList []*modelpolicy.CryptographicSuiteEvaluation
	for _, evaluationType := range evaluations {
		evaluationList = append(evaluationList, buildJSONEvaluation(evaluationType))
	}
	return evaluationList
}

// jsonInformationText ports the private
// getInformationText(JsonObjectWrapper) helper.
func jsonInformationText(algorithmType jsonObject) []string {
	information := algorithmType.asObject(jsonConstraintInformation)
	if information == nil {
		return nil
	}
	return information.asStringList(jsonConstraintText)
}

// buildJSONEvaluation ports the private
// buildEvaluation(JsonObjectWrapper) helper.
func buildJSONEvaluation(evaluationType jsonObject) *modelpolicy.CryptographicSuiteEvaluation {
	evaluation := modelpolicy.NewCryptographicSuiteEvaluation()
	evaluation.SetParameterList(buildJSONParameterList(evaluationType.asObjectList(jsonConstraintParameter)))

	if validity := evaluationType.asObject(jsonConstraintValidity); validity != nil {
		evaluation.SetValidityStart(asDate(validity, jsonConstraintStart))
		evaluation.SetValidityEnd(asDate(validity, jsonConstraintEnd))
	}

	evaluation.SetAlgorithmUsage(jsonAlgorithmUsage(evaluationType))
	evaluation.SetRecommendation(jsonRecommendation(evaluationType))

	return evaluation
}

// buildJSONParameterList ports the private
// buildParameterList(List<JsonObjectWrapper>) helper.
func buildJSONParameterList(parameters []jsonObject) []*modelpolicy.CryptographicSuiteParameter {
	if len(parameters) == 0 {
		return nil
	}
	parameterList := make([]*modelpolicy.CryptographicSuiteParameter, 0, len(parameters))
	for _, parameterType := range parameters {
		parameterList = append(parameterList, buildJSONParameter(parameterType))
	}
	return parameterList
}

// buildJSONParameter ports the private buildParameter(JsonObjectWrapper)
// helper.
func buildJSONParameter(parameterType jsonObject) *modelpolicy.CryptographicSuiteParameter {
	parameter := modelpolicy.NewCryptographicSuiteParameter()
	parameter.SetName(parameterType.asString(jsonConstraintName))
	parameter.SetMin(jsonToInteger(parameterType, jsonConstraintMin))
	parameter.SetMax(jsonToInteger(parameterType, jsonConstraintMax))
	return parameter
}

// jsonAlgorithmUsage ports the private
// getAlgorithmUsage(JsonObjectWrapper) helper.
func jsonAlgorithmUsage(evaluationType jsonObject) []enumerations.CryptographicSuiteAlgorithmUsage {
	algorithmUsageStr := evaluationType.asString(jsonConstraintAlgorithmUsage)
	if algorithmUsageStr == "" {
		return nil
	}
	algorithmUsage := enumerations.CryptographicSuiteAlgorithmUsageFromURI(algorithmUsageStr)
	if algorithmUsage == "" {
		return nil
	}
	return []enumerations.CryptographicSuiteAlgorithmUsage{algorithmUsage}
}

// jsonRecommendation ports the private
// getRecommendation(JsonObjectWrapper) helper.
func jsonRecommendation(evaluationType jsonObject) enumerations.CryptographicSuiteRecommendation {
	recommendation := evaluationType.asString(jsonConstraintRecommendation)
	if recommendation == "" {
		return ""
	}
	return enumerations.CryptographicSuiteRecommendationFromValue(recommendation)
}

// jsonCatalogueVersion ports the private getVersion() helper.
func jsonCatalogueVersion(securitySuitabilityPolicy jsonObject) string {
	if version := securitySuitabilityPolicy.asString(jsonConstraintVersion); version != "" {
		return version
	}
	return cryptographicSuiteJsonCatalogueDefaultVersion
}

// jsonCatalogueLang ports the private getLang() helper.
func jsonCatalogueLang(securitySuitabilityPolicy jsonObject) string {
	if lang := securitySuitabilityPolicy.asString(jsonConstraintLang); lang != "" {
		return lang
	}
	return cryptographicSuiteJsonCatalogueDefaultLang
}

// jsonToInteger ports the private toInteger(Number) helper, folded
// together with its asNumber(name) call site (both buildParameter call
// sites pass the immediate result of asNumber straight to toInteger).
func jsonToInteger(parameterType jsonObject, name string) *int {
	n, ok := parameterType.asNumber(name)
	if !ok {
		return nil
	}
	v := int(n)
	return &v
}

// asDate gets a value of the header name as a time.Time (date only). If
// not present, or not able to convert, returns nil. Ports the private
// getAsDate(JsonObjectWrapper, String) helper.
//
// Java's IllegalArgumentException from a malformed (but present) date
// string propagates uncaught out of buildEvaluation, through
// buildEvaluationList, and is only caught by buildAlgorithm's
// try/catch (skipping that whole algorithm entry - see buildJSONAlgorithm's
// doc comment on why there is no equivalent catch here). This function
// instead treats an unparseable date the same as an absent one
// (returns nil), which is more permissive for a single malformed date but
// avoids introducing panic/recover control flow for a case the upstream
// resources this package ports (dss-crypto-suite.json) never exercise.
func asDate(jsonObj jsonObject, name string) *time.Time {
	dateString := jsonObj.asString(name)
	if dateString == "" {
		return nil
	}
	t, err := rfc3339GetDate(dateString)
	if err != nil {
		return nil
	}
	return &t
}

// asDateTime gets a value of the header name as a time.Time (with
// time). If not present, or not able to convert, returns nil. Ports the
// private getAsDateTime(JsonObjectWrapper, String) helper; see asDate's
// doc comment for the same malformed-input handling note (buildMetadata,
// this method's only caller, has no enclosing try/catch in Java either -
// an unparseable PolicyIssueDate/NextUpdate propagates all the way out of
// getCryptographicSuite in Java, but callers here never observe
// that path since dss-crypto-suite.json's dates are well-formed).
func asDateTime(jsonObj jsonObject, name string) *time.Time {
	dateString := jsonObj.asString(name)
	if dateString == "" {
		return nil
	}
	t, err := rfc3339GetDateTime(dateString)
	if err != nil {
		return nil
	}
	return &t
}

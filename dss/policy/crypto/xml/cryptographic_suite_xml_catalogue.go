// Ported from dss-policy-crypto-xml/.../xml/CryptographicSuiteXmlCatalogue.java (DSS 6.5.RC1).
package cryptoxml

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
)

// xmlCatalogueDefaultVersion/xmlCatalogueDefaultLang are the xs:attribute
// defaults SecuritySuitabilityPolicyType.version/.lang carry in
// rfc5698.xsd (default="1"/default="en"); the JAXB RI applies these during
// unmarshalling when the attribute is absent from the document, a step
// Go's encoding/xml does not perform on its own - see xml_types.go's
// SecuritySuitabilityPolicyType doc comment.
const (
	xmlCatalogueDefaultVersion = "1"
	xmlCatalogueDefaultLang    = "en"
)

// newCryptographicSuiteXmlCatalogue parses an ETSI TS 119 322 XML
// cryptographic suite catalog document and returns the extracted values.
// Ports CryptographicSuiteXmlCatalogue's constructor, collapsed with the
// class body itself - see cryptojson's newCryptographicSuiteJsonCatalogue
// for why buildMetadata/buildAlgorithmList become constructor-supplied
// closures rather than a subclass.
//
// Java's Objects.requireNonNull(securitySuitabilityPolicy, "...") becomes a
// panic with the same message, per PORTING.md's requireNonNull rule.
func newCryptographicSuiteXmlCatalogue(securitySuitabilityPolicy *SecuritySuitabilityPolicyType) *modelpolicy.CryptographicSuiteCatalogue {
	if securitySuitabilityPolicy == nil {
		panic("securitySuitabilityPolicy cannot be null!")
	}
	return modelpolicy.NewCryptographicSuiteCatalogue(
		func() *modelpolicy.CryptographicSuiteMetadata {
			return buildXMLMetadata(securitySuitabilityPolicy)
		},
		func() []*modelpolicy.CryptographicSuiteAlgorithm {
			return buildXMLAlgorithmList(securitySuitabilityPolicy)
		},
	)
}

// buildXMLMetadata ports the protected buildMetadata() method.
func buildXMLMetadata(securitySuitabilityPolicy *SecuritySuitabilityPolicyType) *modelpolicy.CryptographicSuiteMetadata {
	metadata := modelpolicy.NewCryptographicSuiteMetadata()

	// Java's PolicyName/Publisher are non-optional child elements
	// (minOccurs defaults to 1 in the XSD), so JAXB never leaves them
	// null and the constructor's own `policyName != null`/`publisher !=
	// null` guards are unreachable in practice; ported faithfully
	// anyway since Go's value-typed (non-pointer) struct fields make
	// "always present" the natural representation.
	metadata.SetPolicyName(securitySuitabilityPolicy.PolicyName.Name)
	metadata.SetPolicyOID(securitySuitabilityPolicy.PolicyName.ObjectIdentifier)
	metadata.SetPolicyURI(securitySuitabilityPolicy.PolicyName.URI)

	metadata.SetPublisherName(securitySuitabilityPolicy.Publisher.Name)
	metadata.SetPublisherAddress(securitySuitabilityPolicy.Publisher.Address)
	metadata.SetPublisherURI(securitySuitabilityPolicy.Publisher.URI)

	if t, ok := xsdDateTimeToTime(securitySuitabilityPolicy.PolicyIssueDate); ok {
		metadata.SetPolicyIssueDate(&t)
	}
	if t, ok := xsdDateTimeToTime(securitySuitabilityPolicy.NextUpdate); ok {
		metadata.SetNextUpdate(&t)
	}
	metadata.SetUsage(securitySuitabilityPolicy.Usage)

	metadata.SetVersion(xmlCatalogueVersionOrDefault(securitySuitabilityPolicy.Version))
	metadata.SetLang(xmlCatalogueLangOrDefault(securitySuitabilityPolicy.Lang))
	metadata.SetId(securitySuitabilityPolicy.Id)

	return metadata
}

// xmlCatalogueVersionOrDefault/xmlCatalogueLangOrDefault apply the
// xs:attribute defaults described at this file's top.
func xmlCatalogueVersionOrDefault(version string) string {
	if version != "" {
		return version
	}
	return xmlCatalogueDefaultVersion
}

func xmlCatalogueLangOrDefault(lang string) string {
	if lang != "" {
		return lang
	}
	return xmlCatalogueDefaultLang
}

// buildXMLAlgorithmList ports the protected buildAlgorithmList() method.
func buildXMLAlgorithmList(securitySuitabilityPolicy *SecuritySuitabilityPolicyType) []*modelpolicy.CryptographicSuiteAlgorithm {
	var algorithmList []*modelpolicy.CryptographicSuiteAlgorithm
	for i := range securitySuitabilityPolicy.Algorithm {
		if algorithm := buildXMLAlgorithm(&securitySuitabilityPolicy.Algorithm[i]); algorithm != nil {
			algorithmList = append(algorithmList, algorithm)
		}
	}
	return algorithmList
}

// buildXMLAlgorithm ports the private buildAlgorithm(AlgorithmType)
// method.
//
// Java catches any exception raised while processing a single algorithm
// entry, logs it (slf4j dropped, per PORTING.md), and skips the entry
// (returning null). This port's accessors never panic on the XSD-shaped
// input this package's Go structs already constrain (see xml_types.go),
// so there is nothing to recover from here - unlike buildJSONAlgorithm's
// counterpart, whose input shape (a generic jsonObject) is not statically
// constrained the same way.
func buildXMLAlgorithm(algorithmType *AlgorithmType) *modelpolicy.CryptographicSuiteAlgorithm {
	algorithm := modelpolicy.NewCryptographicSuiteAlgorithm()

	algorithm.SetAlgorithmIdentifierName(algorithmType.AlgorithmIdentifier.Name)
	algorithm.SetAlgorithmIdentifierOIDs(algorithmType.AlgorithmIdentifier.ObjectIdentifier)
	algorithm.SetAlgorithmIdentifierURIs(algorithmType.AlgorithmIdentifier.URI)

	algorithm.SetEvaluationList(buildXMLEvaluationList(algorithmType.Evaluation))
	algorithm.SetInformationTextList(xmlInformationText(algorithmType))

	return algorithm
}

// buildXMLEvaluationList ports the private
// buildEvaluationList(List<EvaluationType>) helper.
func buildXMLEvaluationList(evaluations []EvaluationType) []*modelpolicy.CryptographicSuiteEvaluation {
	var evaluationList []*modelpolicy.CryptographicSuiteEvaluation
	for i := range evaluations {
		evaluationList = append(evaluationList, buildXMLEvaluation(&evaluations[i]))
	}
	return evaluationList
}

// xmlInformationText ports the private getInformationText(AlgorithmType)
// helper.
func xmlInformationText(algorithmType *AlgorithmType) []string {
	if algorithmType.Information == nil {
		return nil
	}
	return algorithmType.Information.Text
}

// buildXMLEvaluation ports the private buildEvaluation(EvaluationType)
// helper.
func buildXMLEvaluation(evaluationType *EvaluationType) *modelpolicy.CryptographicSuiteEvaluation {
	evaluation := modelpolicy.NewCryptographicSuiteEvaluation()
	evaluation.SetParameterList(buildXMLParameterList(evaluationType.Parameter))

	if t, ok := xsdDateToTime(evaluationType.Validity.Start); ok {
		evaluation.SetValidityStart(&t)
	}
	if t, ok := xsdDateToTime(evaluationType.Validity.End); ok {
		evaluation.SetValidityEnd(&t)
	}

	evaluation.SetAlgorithmUsage(xmlAlgorithmUsageList(evaluationType.MoreDetails))
	evaluation.SetRecommendation(xmlRecommendation(evaluationType.MoreDetails))

	return evaluation
}

// buildXMLParameterList ports the private
// buildParameterList(List<ParameterType>) helper.
func buildXMLParameterList(parameters []ParameterType) []*modelpolicy.CryptographicSuiteParameter {
	if len(parameters) == 0 {
		return nil
	}
	parameterList := make([]*modelpolicy.CryptographicSuiteParameter, 0, len(parameters))
	for i := range parameters {
		parameterList = append(parameterList, buildXMLParameter(&parameters[i]))
	}
	return parameterList
}

// buildXMLParameter ports the private buildParameter(ParameterType)
// helper.
func buildXMLParameter(parameterType *ParameterType) *modelpolicy.CryptographicSuiteParameter {
	parameter := modelpolicy.NewCryptographicSuiteParameter()
	parameter.SetName(parameterType.Name)
	parameter.SetMin(parameterType.Min)
	parameter.SetMax(parameterType.Max)
	return parameter
}

// xmlAlgorithmUsageList ports the private
// getAlgorithmUsageList(ExtensionType) helper.
func xmlAlgorithmUsageList(extensionType *ExtensionType) []enumerations.CryptographicSuiteAlgorithmUsage {
	if extensionType == nil {
		return nil
	}
	var algorithmUsageList []enumerations.CryptographicSuiteAlgorithmUsage
	for _, algorithmUsageURI := range extensionType.AlgorithmUsage {
		if algorithmUsage := enumerations.CryptographicSuiteAlgorithmUsageFromURI(algorithmUsageURI); algorithmUsage != "" {
			algorithmUsageList = append(algorithmUsageList, algorithmUsage)
		}
	}
	return algorithmUsageList
}

// xmlRecommendation ports the private getRecommendation(ExtensionType)
// helper.
func xmlRecommendation(extensionType *ExtensionType) enumerations.CryptographicSuiteRecommendation {
	if extensionType == nil || extensionType.Recommendation == "" {
		return ""
	}
	return enumerations.CryptographicSuiteRecommendationFromValue(extensionType.Recommendation)
}

// Ported from dss-policy-crypto-json/.../json/CryptographicSuiteJsonConstraints.java (DSS 6.5.RC1).
package cryptojson

// Header name constants for a JSON cryptographic suite as per ETSI TS 119
// 322. Ports CryptographicSuiteJsonConstraints; unexported (package-private
// in spirit) since only this package's own catalogue/factory reference
// them - upstream's public visibility exists for cross-module reuse this
// port has no equivalent of within the manifest.
const (
	jsonConstraintAddress                   = "Address"
	jsonConstraintAlgorithm                 = "Algorithm"
	jsonConstraintAlgorithmIdentifier       = "AlgorithmIdentifier"
	jsonConstraintAlgorithmUsage            = "AlgorithmUsage"
	jsonConstraintAny                       = "Any"
	jsonConstraintEnd                       = "End"
	jsonConstraintEvaluation                = "Evaluation"
	jsonConstraintID                        = "id"
	jsonConstraintInformation               = "Information"
	jsonConstraintLang                      = "lang"
	jsonConstraintMax                       = "Max"
	jsonConstraintMin                       = "Min"
	jsonConstraintName                      = "name"
	jsonConstraintNameC                     = "Name"
	jsonConstraintNextUpdate                = "NextUpdate"
	jsonConstraintObjectIdentifier          = "ObjectIdentifier"
	jsonConstraintParameter                 = "Parameter"
	jsonConstraintPolicyIssueDate           = "PolicyIssueDate"
	jsonConstraintPolicyName                = "PolicyName"
	jsonConstraintPublisher                 = "Publisher"
	jsonConstraintRecommendation            = "Recommendation"
	jsonConstraintSecuritySuitabilityPolicy = "SecuritySuitabilityPolicy"
	jsonConstraintStart                     = "Start"
	jsonConstraintText                      = "Text"
	jsonConstraintURI                       = "URI"
	jsonConstraintUsage                     = "Usage"
	jsonConstraintValidity                  = "Validity"
	jsonConstraintVersion                   = "version"
)
